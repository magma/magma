"""
Copyright (c) 2018-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""

import unittest
from concurrent.futures import Future
from unittest.mock import MagicMock

import warnings
from lte.protos.mconfig.mconfigs_pb2 import PipelineD
from lte.protos.policydb_pb2 import FlowDescription, FlowMatch, PolicyRule, \
    RedirectInformation
from magma.pipelined.app.enforcement_stats import EnforcementStatsController
from magma.pipelined.bridge_util import BridgeTools
from magma.pipelined.tests.app.packet_builder import IPPacketBuilder, \
    TCPPacketBuilder
from magma.pipelined.tests.app.packet_injector import ScapyPacketInjector
from magma.pipelined.tests.app.start_pipelined import PipelinedController, \
    TestSetup
from magma.pipelined.tests.app.subscriber import RyuDirectSubscriberContext
from magma.pipelined.tests.app.table_isolation import RyuDirectTableIsolator, \
    RyuForwardFlowArgsBuilder
from magma.pipelined.tests.pipelined_test_util import FlowVerifier, \
    create_service_manager, get_enforcement_stats, start_ryu_app_thread, \
    stop_ryu_app_thread, wait_after_send, wait_for_enforcement_stats
from scapy.all import IP


class EnforcementStatsTest(unittest.TestCase):
    BRIDGE = 'testing_br'
    IFACE = 'testing_br'
    MAC_DEST = "5e:cc:cc:b1:49:4b"

    def setUp(self):
        """
        Starts the thread which launches ryu apps

        Create a testing bridge, add a port, setup the port interfaces. Then
        launch the ryu apps for testing pipelined. Gets the references
        to apps launched by using futures.

        Mocks the redis policy_dictionary of enforcement_controller.
        Mocks the loop for testing EnforcementStatsController
        """
        super(EnforcementStatsTest, self).setUpClass()
        warnings.simplefilter('ignore')
        self._static_rule_dict = {}
        service_manager = create_service_manager([PipelineD.ENFORCEMENT])
        self._tbl_num = service_manager.get_table_num(
            EnforcementStatsController.APP_NAME)

        enforcement_controller_reference = Future()
        testing_controller_reference = Future()
        enf_stat_ref = Future()

        """
        Enforcement_stats reports data by using loop.call_soon_threadsafe, but
        as we don't have an eventloop in testing, just directly call the stats
        handling function

        Here is how the mocked function is used in EnforcementStatsController:
        self.loop.call_soon_threadsafe(self._handle_flow_stats, ev.msg.body)
        """
        def mock_thread_safe(cmd, body):
            cmd(body)
        loop_mock = MagicMock()
        loop_mock.call_soon_threadsafe = mock_thread_safe

        test_setup = TestSetup(
            apps=[PipelinedController.Enforcement,
                  PipelinedController.Testing,
                  PipelinedController.Enforcement_stats],
            references={
                PipelinedController.Enforcement:
                    enforcement_controller_reference,
                PipelinedController.Testing:
                    testing_controller_reference,
                PipelinedController.Enforcement_stats:
                    enf_stat_ref
            },
            config={
                'bridge_name': self.BRIDGE,
                'bridge_ip_address': '192.168.128.1',
                'enforcement': {'poll_interval': 5},
                'nat_iface': 'eth2',
                'enodeb_iface': 'eth1',
                'enable_queue_pgm': False,
            },
            mconfig=PipelineD(
                relay_enabled=True,
            ),
            loop=loop_mock,
            service_manager=service_manager,
            integ_test=False,
            rpc_stubs={'sessiond': MagicMock()}
        )

        BridgeTools.create_bridge(self.BRIDGE, self.IFACE)

        self.thread = start_ryu_app_thread(test_setup)

        self.enforcement_stats_controller = enf_stat_ref.result()
        self.enforcement_controller = enforcement_controller_reference.result()
        self.testing_controller = testing_controller_reference.result()

        self.enforcement_stats_controller._report_usage = MagicMock()

        self.enforcement_controller._policy_dict = self._static_rule_dict
        self.enforcement_controller._redirect_manager._save_redirect_entry = \
            MagicMock()

    def tearDown(self):
        stop_ryu_app_thread(self.thread)
        BridgeTools.destroy_bridge(self.BRIDGE)

    def _wait_func(self, stat_names):
        def func():
            wait_after_send(self.testing_controller)
            wait_for_enforcement_stats(self.enforcement_stats_controller,
                                       stat_names)

        return func

    def test_subscriber_policy(self):
        """
        Adds 2 policies to subscriber, verifies that EnforcementStatsController
        reports correct stats to sessiond

        Assert:
            UPLINK policy matches 128 packets (*34 = 4352 bytes)
            DOWNLINK policy matches 256 packets (*34 = 8704 bytes)
            No other stats are reported
        """
        imsi = 'IMSI001010000000013'
        sub_ip = '192.168.128.74'
        num_pkts_tx_match = 128
        num_pkts_rx_match = 256

        """ Create 2 policy rules for the subscriber """
        flow_list1 = [FlowDescription(
            match=FlowMatch(
                ipv4_dst='45.10.0.0/25', direction=FlowMatch.UPLINK),
            action=FlowDescription.PERMIT)
        ]
        flow_list2 = [FlowDescription(
            match=FlowMatch(
                ipv4_src='45.10.0.0/24', direction=FlowMatch.DOWNLINK),
            action=FlowDescription.PERMIT)
        ]
        policies = [
            PolicyRule(id='tx_match', priority=3, flow_list=flow_list1),
            PolicyRule(id='rx_match', priority=5, flow_list=flow_list2)
        ]
        enf_stat_name = [imsi + '|tx_match', imsi + '|rx_match']

        """ Setup subscriber, setup table_isolation to fwd pkts """
        self._static_rule_dict[policies[0].id] = policies[0]
        self._static_rule_dict[policies[1].id] = policies[1]
        sub_context = RyuDirectSubscriberContext(
            imsi, sub_ip, self.enforcement_controller, self._tbl_num
        ).add_static_rule(policies[0].id).add_static_rule(policies[1].id)
        isolator = RyuDirectTableIsolator(
            RyuForwardFlowArgsBuilder.from_subscriber(sub_context.cfg)
                                     .build_requests(),
            self.testing_controller
        )

        """ Create 2 sets of packets, for policry rule1&2 """
        pkt_sender = ScapyPacketInjector(self.IFACE)
        packet1 = IPPacketBuilder()\
            .set_ip_layer('45.10.0.0/20', sub_ip)\
            .set_ether_layer(self.MAC_DEST, "00:00:00:00:00:00")\
            .build()
        packet2 = IPPacketBuilder()\
            .set_ip_layer(sub_ip, '45.10.0.0/20')\
            .set_ether_layer(self.MAC_DEST, "00:00:00:00:00:00")\
            .build()

        # =========================== Verification ===========================
        flow_verifier = FlowVerifier([], self._wait_func(enf_stat_name))
        """ Send packets, wait until pkts are received by ovs and enf stats """
        with isolator, sub_context, flow_verifier:
            pkt_sender.send(packet1)
            pkt_sender.send(packet2)

        stats = get_enforcement_stats(
            self.enforcement_stats_controller._report_usage.call_args_list)

        self.assertEqual(stats[enf_stat_name[0]].sid, imsi)
        self.assertEqual(stats[enf_stat_name[0]].rule_id, "tx_match")
        self.assertEqual(stats[enf_stat_name[0]].bytes_rx, 0)
        self.assertEqual(stats[enf_stat_name[0]].bytes_tx,
                         num_pkts_tx_match * len(packet1))

        self.assertEqual(stats[enf_stat_name[1]].sid, imsi)
        self.assertEqual(stats[enf_stat_name[1]].rule_id, "rx_match")
        self.assertEqual(stats[enf_stat_name[1]].bytes_tx, 0)

        # downlink packets will discount ethernet header by default
        # so, only count the IP portion
        total_bytes_pkt2 = num_pkts_rx_match * len(packet2[IP])
        self.assertEqual(stats[enf_stat_name[1]].bytes_rx, total_bytes_pkt2)

        self.assertEqual(len(stats), 2)

    def test_redirect_policy(self):
        """
        Add a redirect policy, verifies that EnforcementStatsController reports
        correct stats to sessiond

        Assert:
            1 Packet is matched and reported
        """
        redirect_ips = ["185.128.101.5", "185.128.121.4"]
        self.enforcement_controller._redirect_manager._dns_cache.get(
            "about.sha.ddih.org", lambda: redirect_ips, max_age=42
        )
        imsi = 'IMSI010000000088888'
        sub_ip = '192.168.128.74'
        flow_list = [FlowDescription(match=FlowMatch())]
        policy = PolicyRule(
            id='redir_test', priority=3, flow_list=flow_list,
            redirect=RedirectInformation(
                support=1,
                address_type=2,
                server_address="http://about.sha.ddih.org/"
            )
        )
        stat_name = imsi + '|redir_test'

        """ Setup subscriber, setup table_isolation to fwd pkts """
        self._static_rule_dict[policy.id] = policy
        sub_context = RyuDirectSubscriberContext(
            imsi, sub_ip, self.enforcement_controller, self._tbl_num
        ).add_dynamic_rule(policy)
        isolator = RyuDirectTableIsolator(
            RyuForwardFlowArgsBuilder.from_subscriber(sub_context.cfg)
                                     .build_requests(),
            self.testing_controller
        )
        pkt_sender = ScapyPacketInjector(self.IFACE)
        packet = TCPPacketBuilder() \
            .set_tcp_layer(42132, 80, 321) \
            .set_tcp_flags("S") \
            .set_ip_layer('151.42.41.122', sub_ip) \
            .set_ether_layer(self.MAC_DEST, "00:00:00:00:00:00") \
            .build()

        # =========================== Verification ===========================
        flow_verifier = FlowVerifier([], self._wait_func([stat_name]))
        """ Send packet, wait until pkts are received by ovs and enf stats """
        with isolator, sub_context, flow_verifier:
            pkt_sender.send(packet)

        stats = get_enforcement_stats(
            self.enforcement_stats_controller._report_usage.call_args_list)

        self.assertEqual(stats[stat_name].sid, imsi)
        self.assertEqual(stats[stat_name].rule_id, "redir_test")
        self.assertEqual(stats[stat_name].bytes_rx, 0)
        self.assertEqual(stats[stat_name].bytes_tx, len(packet))

if __name__ == "__main__":
    unittest.main()

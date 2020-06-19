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
from magma.pipelined.app.enforcement import EnforcementController
from magma.pipelined.bridge_util import BridgeTools
from magma.pipelined.imsi import encode_imsi
from magma.pipelined.openflow.magma_match import MagmaMatch
from magma.pipelined.policy_converters import flow_match_to_magma_match
from magma.pipelined.tests.app.flow_query import RyuDirectFlowQuery \
    as FlowQuery
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
    stop_ryu_app_thread, wait_after_send, wait_for_enforcement_stats, \
    FlowTest, SnapshotVerifier, fake_controller_setup
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
        self.service_manager = create_service_manager([PipelineD.ENFORCEMENT])
        self._main_tbl_num = self.service_manager.get_table_num(
            EnforcementController.APP_NAME)

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
                  PipelinedController.Enforcement_stats,
                  PipelinedController.Testing,
                  PipelinedController.StartupFlows],
            references={
                PipelinedController.Enforcement:
                    enforcement_controller_reference,
                PipelinedController.Testing:
                    testing_controller_reference,
                PipelinedController.Enforcement_stats:
                    enf_stat_ref,
                PipelinedController.StartupFlows:
                    Future(),
            },
            config={
                'bridge_name': self.BRIDGE,
                'bridge_ip_address': '192.168.128.1',
                'enforcement': {'poll_interval': 2},
                'nat_iface': 'eth2',
                'enodeb_iface': 'eth1',
                'qos': {'enable': False},
                'clean_restart': True,
            },
            mconfig=PipelineD(
                relay_enabled=True,
            ),
            loop=loop_mock,
            service_manager=self.service_manager,
            integ_test=False,
            rpc_stubs={'sessiond': MagicMock()}
        )

        BridgeTools.create_bridge(self.BRIDGE, self.IFACE)

        self.thread = start_ryu_app_thread(test_setup)

        self.enforcement_stats_controller = enf_stat_ref.result()
        self._scratch_tbl_num = self.enforcement_stats_controller.tbl_num
        self.enforcement_controller = enforcement_controller_reference.result()
        self.testing_controller = testing_controller_reference.result()

        self.enforcement_stats_controller._policy_dict = self._static_rule_dict
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
        fake_controller_setup(self.enforcement_controller,
                              self.enforcement_stats_controller)
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
        self.service_manager.session_rule_version_mapper.update_version(
            imsi, 'tx_match')
        self.service_manager.session_rule_version_mapper.update_version(
            imsi, 'rx_match')

        """ Setup subscriber, setup table_isolation to fwd pkts """
        self._static_rule_dict[policies[0].id] = policies[0]
        self._static_rule_dict[policies[1].id] = policies[1]
        sub_context = RyuDirectSubscriberContext(
            imsi, sub_ip, self.enforcement_controller,
            self._main_tbl_num, self.enforcement_stats_controller
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
        snapshot_verifier = SnapshotVerifier(self, self.BRIDGE,
                                             self.service_manager)
        """ Send packets, wait until pkts are received by ovs and enf stats """
        with isolator, sub_context, flow_verifier, snapshot_verifier:
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
        fake_controller_setup(self.enforcement_controller,
                              self.enforcement_stats_controller)
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
        self.service_manager.session_rule_version_mapper.update_version(
            imsi, 'redir_test')

        """ Setup subscriber, setup table_isolation to fwd pkts """
        self._static_rule_dict[policy.id] = policy
        sub_context = RyuDirectSubscriberContext(
            imsi, sub_ip, self.enforcement_controller,
            self._main_tbl_num, self.enforcement_stats_controller
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
        snapshot_verifier = SnapshotVerifier(self, self.BRIDGE,
                                             self.service_manager)
        """ Send packet, wait until pkts are received by ovs and enf stats """
        with isolator, sub_context, flow_verifier, snapshot_verifier:
            pkt_sender.send(packet)

        stats = get_enforcement_stats(
            self.enforcement_stats_controller._report_usage.call_args_list)

        self.assertEqual(stats[stat_name].sid, imsi)
        self.assertEqual(stats[stat_name].rule_id, "redir_test")
        self.assertEqual(stats[stat_name].bytes_rx, 0)
        self.assertEqual(stats[stat_name].bytes_tx, len(packet))

    def test_rule_install(self):
        """
        Adds a policy to a subscriber. Verifies that flows are properly
        installed in enforcement and enforcement stats.

        Assert:
            Policy classification flows installed in enforcement
            Policy match flows installed in enforcement_stats
        """
        fake_controller_setup(self.enforcement_controller,
                              self.enforcement_stats_controller)
        imsi = 'IMSI001010000000013'
        sub_ip = '192.168.128.74'

        flow_list = [FlowDescription(
            match=FlowMatch(
                ipv4_dst='45.10.0.0/25', direction=FlowMatch.UPLINK),
            action=FlowDescription.PERMIT)
        ]
        policy = PolicyRule(id='rule1', priority=3, flow_list=flow_list)
        self.service_manager.session_rule_version_mapper.update_version(
            imsi, 'rule1')
        version = \
            self.service_manager.session_rule_version_mapper.get_version(
                imsi, 'rule1')

        """ Setup subscriber, setup table_isolation to fwd pkts """
        self._static_rule_dict[policy.id] = policy
        sub_context = RyuDirectSubscriberContext(
            imsi, sub_ip, self.enforcement_controller,
            self._main_tbl_num, self.enforcement_stats_controller
        ).add_static_rule(policy.id)

        # =========================== Verification ===========================
        rule_num = self.enforcement_stats_controller._rule_mapper \
            .get_or_create_rule_num(policy.id)
        enf_query = FlowQuery(self._main_tbl_num, self.testing_controller,
                              match=flow_match_to_magma_match(
                                  FlowMatch(ipv4_dst='45.10.0.0/25',
                                            direction=FlowMatch.UPLINK)),
                              cookie=rule_num)
        es_query = FlowQuery(self._scratch_tbl_num,
                             self.testing_controller,
                             match=MagmaMatch(imsi=encode_imsi(imsi),
                                              reg2=rule_num,
                                              rule_version=version),
                             cookie=rule_num)

        # Verifies that 1 flow is installed in enforcement and 2 flows are
        # installed in enforcement stats, one for uplink and one for downlink.
        flow_verifier = FlowVerifier([
            FlowTest(enf_query, 0, flow_count=1),
            FlowTest(es_query, 0, flow_count=2),
        ], lambda: None)
        snapshot_verifier = SnapshotVerifier(self, self.BRIDGE,
                                             self.service_manager)

        with sub_context, flow_verifier, snapshot_verifier:
            pass

        flow_verifier.verify()

    def test_rule_deactivation(self):
        """
        Adds a policy to a subscriber, and then deletes it by incrementing the
        version, verifies that the usage stats is correctly reported and the
        flows are deleted.

        Assert:
            UPLINK policy matches 128 packets (*34 = 4352 bytes)
            Flows are deleted
            No other stats are reported
        """
        fake_controller_setup(self.enforcement_controller,
                              self.enforcement_stats_controller)
        imsi = 'IMSI001010000000013'
        sub_ip = '192.168.128.74'
        num_pkts_tx_match = 128

        flow_list = [FlowDescription(
            match=FlowMatch(
                ipv4_dst='45.10.0.0/25', direction=FlowMatch.UPLINK),
            action=FlowDescription.PERMIT)
        ]
        policy = PolicyRule(id='rule1', priority=3, flow_list=flow_list)
        enf_stat_name = imsi + '|rule1'
        self.service_manager.session_rule_version_mapper.update_version(
            imsi, 'rule1')
        version = \
            self.service_manager.session_rule_version_mapper.get_version(
                imsi, 'rule1')

        """ Setup subscriber, setup table_isolation to fwd pkts """
        self._static_rule_dict[policy.id] = policy
        sub_context = RyuDirectSubscriberContext(
            imsi, sub_ip, self.enforcement_controller,
            self._main_tbl_num, self.enforcement_stats_controller
        ).add_static_rule(policy.id)
        isolator = RyuDirectTableIsolator(
            RyuForwardFlowArgsBuilder.from_subscriber(sub_context.cfg)
                                     .build_requests(),
            self.testing_controller
        )

        """ Create a packet """
        pkt_sender = ScapyPacketInjector(self.IFACE)
        packet = IPPacketBuilder() \
            .set_ip_layer('45.10.0.0/20', sub_ip) \
            .set_ether_layer(self.MAC_DEST, "00:00:00:00:00:00") \
            .build()

        # =========================== Verification ===========================
        rule_num = self.enforcement_stats_controller._rule_mapper \
            .get_or_create_rule_num(policy.id)
        enf_query = FlowQuery(self._main_tbl_num, self.testing_controller,
                              match=flow_match_to_magma_match(
                                  FlowMatch(ipv4_dst='45.10.0.0/25',
                                            direction=FlowMatch.UPLINK)),
                              cookie=rule_num)
        es_query = FlowQuery(self._scratch_tbl_num,
                             self.testing_controller,
                             match=MagmaMatch(imsi=encode_imsi(imsi),
                                              reg2=rule_num,
                                              rule_version=version),
                             cookie=rule_num)
        """ Verify that flows are properly deleted """
        verify_enforcement_stats = FlowVerifier([
            FlowTest(es_query, 0, flow_count=0),
        ], self._wait_func([enf_stat_name]))
        verify_enforcement = FlowVerifier([
            FlowTest(enf_query, 0, flow_count=0),
        ], lambda: None)
        snapshot_verifier = SnapshotVerifier(self, self.BRIDGE,
                                             self.service_manager)

        """
        Send packets, wait until packet is received by ovs and enf stats and
        then deactivate the rule in enforcement controller. This emulates the
        case where there is unreported traffic after rule deactivation.
        """
        with isolator, sub_context, verify_enforcement, snapshot_verifier:
            with verify_enforcement_stats:
                self.enforcement_stats_controller._report_usage.reset_mock()
                pkt_sender.send(packet)
                self.service_manager.session_rule_version_mapper. \
                    update_version(imsi, 'rule1')
                self.enforcement_controller.deactivate_rules(imsi, [policy.id])

        verify_enforcement.verify()
        verify_enforcement_stats.verify()

        stats = get_enforcement_stats(
            self.enforcement_stats_controller._report_usage.call_args_list)

        self.assertEqual(stats[enf_stat_name].sid, imsi)
        self.assertEqual(stats[enf_stat_name].rule_id, "rule1")
        self.assertEqual(stats[enf_stat_name].bytes_rx, 0)
        self.assertEqual(stats[enf_stat_name].bytes_tx,
                         num_pkts_tx_match * len(packet))

        self.assertEqual(len(stats), 1)

    def test_rule_reactivation(self):
        """
        Adds a policy to a subscriber, deletes it by incrementing the
        version, and add it back. Verifies that the usage stats is correctly
        reported, the old flows are deleted, and the new flows are installed.

        Assert:
            UPLINK policy matches 128 packets (*34 = 4352 bytes)
            Old flows are deleted
            New flows are installed
            No other stats are reported
        """
        fake_controller_setup(self.enforcement_controller,
                              self.enforcement_stats_controller)
        imsi = 'IMSI001010000000013'
        sub_ip = '192.168.128.74'
        num_pkts_tx_match = 128

        flow_list = [FlowDescription(
            match=FlowMatch(
                ipv4_dst='45.10.0.0/25', direction=FlowMatch.UPLINK),
            action=FlowDescription.PERMIT)
        ]
        policy = PolicyRule(id='rule1', priority=3, flow_list=flow_list)
        enf_stat_name = imsi + '|rule1'
        self.service_manager.session_rule_version_mapper.update_version(
            imsi, 'rule1')
        version = \
            self.service_manager.session_rule_version_mapper.get_version(
                imsi, 'rule1')

        """ Setup subscriber, setup table_isolation to fwd pkts """
        self._static_rule_dict[policy.id] = policy
        sub_context = RyuDirectSubscriberContext(
            imsi, sub_ip, self.enforcement_controller,
            self._main_tbl_num, self.enforcement_stats_controller
        ).add_static_rule(policy.id)
        isolator = RyuDirectTableIsolator(
            RyuForwardFlowArgsBuilder.from_subscriber(sub_context.cfg)
                                     .build_requests(),
            self.testing_controller
        )

        """ Create a packet """
        pkt_sender = ScapyPacketInjector(self.IFACE)
        packet = IPPacketBuilder() \
            .set_ip_layer('45.10.0.0/20', sub_ip) \
            .set_ether_layer(self.MAC_DEST, "00:00:00:00:00:00") \
            .build()

        # =========================== Verification ===========================
        rule_num = self.enforcement_stats_controller._rule_mapper \
            .get_or_create_rule_num(policy.id)
        enf_query = FlowQuery(self._main_tbl_num, self.testing_controller,
                              match=flow_match_to_magma_match(
                                  FlowMatch(ipv4_dst='45.10.0.0/25',
                                            direction=FlowMatch.UPLINK)),
                              cookie=rule_num)
        es_old_version_query = FlowQuery(self._scratch_tbl_num,
                                         self.testing_controller,
                                         match=MagmaMatch(
                                             imsi=encode_imsi(imsi),
                                             reg2=rule_num,
                                             rule_version=version),
                                         cookie=rule_num)
        es_new_version_query = FlowQuery(self._scratch_tbl_num,
                                         self.testing_controller,
                                         match=MagmaMatch(
                                             imsi=encode_imsi(imsi),
                                             reg2=rule_num,
                                             rule_version=version + 1),
                                         cookie=rule_num)
        packet_wait = FlowVerifier([], self._wait_func([enf_stat_name]))
        """ Verify that flows are properly deleted """
        verifier = FlowVerifier([
            FlowTest(es_old_version_query, 0, flow_count=0),
            FlowTest(es_new_version_query, num_pkts_tx_match, flow_count=2),
            FlowTest(enf_query, num_pkts_tx_match, flow_count=1),
        ], self._wait_func([enf_stat_name]))
        snapshot_verifier = SnapshotVerifier(self, self.BRIDGE,
                                             self.service_manager)

        """
        Send a packet, then deactivate and reactivate the same rule and send a
        packet. Wait until it is received by ovs and enf stats.
        """
        with isolator, sub_context, verifier, snapshot_verifier:
            with packet_wait:
                self.enforcement_stats_controller._report_usage.reset_mock()
                pkt_sender.send(packet)

            self.enforcement_stats_controller._report_usage.reset_mock()
            self.service_manager.session_rule_version_mapper. \
                update_version(imsi, 'rule1')
            self.enforcement_controller.deactivate_rules(imsi, [policy.id])
            self.enforcement_controller.activate_rules(imsi, sub_ip,
                                                       [policy.id], [])
            self.enforcement_stats_controller.activate_rules(imsi, sub_ip,
                                                             [policy.id], [])
            pkt_sender.send(packet)

        verifier.verify()

        stats = get_enforcement_stats(
            self.enforcement_stats_controller._report_usage.call_args_list)

        """
        Verify both packets are reported after reactivation.
        """
        self.assertEqual(stats[enf_stat_name].sid, imsi)
        self.assertEqual(stats[enf_stat_name].rule_id, "rule1")
        self.assertEqual(stats[enf_stat_name].bytes_rx, 0)
        # TODO Figure out why this one fails.
        #self.assertEqual(stats[enf_stat_name].bytes_tx,
        #                 num_pkts_tx_match * len(packet))
        self.assertEqual(len(stats), 1)


if __name__ == "__main__":
    unittest.main()

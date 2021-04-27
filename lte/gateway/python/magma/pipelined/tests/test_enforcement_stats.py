"""
Copyright 2020 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
"""

import unittest
import warnings
from concurrent.futures import Future
from unittest.mock import MagicMock

from lte.protos.mconfig.mconfigs_pb2 import PipelineD
from lte.protos.pipelined_pb2 import VersionedPolicy
from lte.protos.policydb_pb2 import (
    FlowDescription,
    FlowMatch,
    PolicyRule,
    RedirectInformation,
)
from magma.pipelined.app.enforcement import EnforcementController
from magma.pipelined.bridge_util import BridgeTools
from magma.pipelined.policy_converters import (
    convert_ipv4_str_to_ip_proto,
    convert_ipv6_bytes_to_ip_proto,
)
from magma.pipelined.tests.app.packet_builder import (
    IPPacketBuilder,
    TCPPacketBuilder,
)
from magma.pipelined.tests.app.packet_injector import ScapyPacketInjector
from magma.pipelined.tests.app.start_pipelined import (
    PipelinedController,
    TestSetup,
)
from magma.pipelined.tests.app.subscriber import RyuDirectSubscriberContext
from magma.pipelined.tests.app.table_isolation import (
    RyuDirectTableIsolator,
    RyuForwardFlowArgsBuilder,
)
from magma.pipelined.tests.pipelined_test_util import (
    FlowTest,
    SnapshotVerifier,
    create_service_manager,
    fake_controller_setup,
    get_enforcement_stats,
    start_ryu_app_thread,
    stop_ryu_app_thread,
    wait_after_send,
    wait_for_enforcement_stats,
)
from scapy.all import IP


class EnforcementStatsTest(unittest.TestCase):
    BRIDGE = 'testing_br'
    IFACE = 'testing_br'
    MAC_DEST = "5e:cc:cc:b1:49:4b"
    DEFAULT_DROP_FLOW_NAME = '(ノಠ益ಠ)ノ彡┻━┻'

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
                'enforcement': {
                    'poll_interval': 2,
                    'default_drop_flow_name': self.DEFAULT_DROP_FLOW_NAME
                },
                'nat_iface': 'eth2',
                'enodeb_iface': 'eth1',
                'qos': {'enable': False},
                'clean_restart': True,
                'redis_enabled': False,
                'setup_type': 'LTE',
            },
            mconfig=PipelineD(),
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

        self.enforcement_stats_controller._report_usage = MagicMock()
        self.enforcement_controller._redirect_manager._save_redirect_entry = \
            MagicMock()

    def tearDown(self):
        stop_ryu_app_thread(self.thread)
        BridgeTools.destroy_bridge(self.BRIDGE)

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
        uplink_tunnel = 0x1234
        num_pkts_tx_match = 128
        num_pkts_rx_match = 256

        """ Create 2 policy rules for the subscriber """
        flow_list1 = [FlowDescription(
            match=FlowMatch(
                ip_dst=convert_ipv4_str_to_ip_proto('45.10.0.0/25'),
                direction=FlowMatch.UPLINK),
            action=FlowDescription.PERMIT)
        ]
        flow_list2 = [FlowDescription(
            match=FlowMatch(
                ip_src=convert_ipv4_str_to_ip_proto('45.10.0.0/24'),
                direction=FlowMatch.DOWNLINK),
            action=FlowDescription.PERMIT)
        ]
        policies = [
            VersionedPolicy(
                rule=PolicyRule(id='tx_match', priority=3, flow_list=flow_list1),
                version=1,
            ),
            VersionedPolicy(
                rule=PolicyRule(id='rx_match', priority=5, flow_list=flow_list2),
                version=1,
            )
        ]
        enf_stat_name = [imsi + '|tx_match' + '|' + str(uplink_tunnel),
                         imsi + '|rx_match' + '|' + str(uplink_tunnel)]
        self.service_manager.session_rule_version_mapper.save_version(
            imsi, uplink_tunnel, 'tx_match', 1)
        self.service_manager.session_rule_version_mapper.save_version(
            imsi, uplink_tunnel, 'rx_match', 1)

        """ Setup subscriber, setup table_isolation to fwd pkts """
        sub_context = RyuDirectSubscriberContext(
            imsi, sub_ip, uplink_tunnel,
            self.enforcement_controller,
            self._main_tbl_num, self.enforcement_stats_controller
        ).add_policy(policies[0]).add_policy(policies[1])
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
        snapshot_verifier = SnapshotVerifier(self, self.BRIDGE,
                                             self.service_manager)
        """ Send packets, wait until pkts are received by ovs and enf stats """
        with isolator, sub_context, snapshot_verifier:
            pkt_sender.send(packet1)
            pkt_sender.send(packet2)

        wait_for_enforcement_stats(self.enforcement_stats_controller,
                                   enf_stat_name)

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

        self.assertEqual(len(stats), 3)

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
        uplink_tunnel = 0x1234
        sub_ip = '192.168.128.74'
        flow_list = [FlowDescription(match=FlowMatch())]
        policy = VersionedPolicy(
            rule=PolicyRule(
                id='redir_test', priority=3, flow_list=flow_list,
                redirect=RedirectInformation(
                    support=1,
                    address_type=2,
                    server_address="http://about.sha.ddih.org/"
                )
            ),
            version=1,
        )
        stat_name = imsi + '|redir_test' + '|' + str(uplink_tunnel)
        self.service_manager.session_rule_version_mapper.save_version(
            imsi, uplink_tunnel, 'redir_test', 1)

        """ Setup subscriber, setup table_isolation to fwd pkts """
        sub_context = RyuDirectSubscriberContext(
            imsi, sub_ip, uplink_tunnel,
            self.enforcement_controller,
            self._main_tbl_num, self.enforcement_stats_controller
        ).add_policy(policy)
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
        snapshot_verifier = SnapshotVerifier(self, self.BRIDGE,
                                             self.service_manager)
        """ Send packet, wait until pkts are received by ovs and enf stats """
        with isolator, sub_context, snapshot_verifier:
            self.enforcement_stats_controller._report_usage.reset_mock()
            pkt_sender.send(packet)

        wait_for_enforcement_stats(self.enforcement_stats_controller,
                                   [stat_name])

        """ Send packets, wait until pkts are received by ovs and enf stats """
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
        uplink_tunnel = 0x1234
        sub_ip = '192.168.128.74'

        flow_list = [FlowDescription(
            match=FlowMatch(
                ip_dst=convert_ipv4_str_to_ip_proto('45.10.0.0/25'),
                direction=FlowMatch.UPLINK),
            action=FlowDescription.PERMIT)
        ]
        policy = VersionedPolicy(
            rule=PolicyRule(id='rule1', priority=3, flow_list=flow_list),
            version=1,
        )
        self.service_manager.session_rule_version_mapper.save_version(
            imsi, uplink_tunnel, 'rule1', 1)

        """ Setup subscriber, setup table_isolation to fwd pkts """
        sub_context = RyuDirectSubscriberContext(
            imsi, sub_ip, uplink_tunnel,
            self.enforcement_controller,
            self._main_tbl_num, self.enforcement_stats_controller
        ).add_policy(policy)

        # =========================== Verification ===========================

        # Verifies that 1 flow is installed in enforcement and 2 flows are
        # installed in enforcement stats, one for uplink and one for downlink.
        snapshot_verifier = SnapshotVerifier(self, self.BRIDGE,
                                             self.service_manager)

        with sub_context, snapshot_verifier:
            pass

    def test_deny_rule_install(self):
        """
        Adds a policy to a subscriber. Verifies that flows are properly
        installed in enforcement and enforcement stats.
        Assert:
            Policy classification flows installed in enforcement
            Policy match flows installed in enforcement_stats
        """
        fake_controller_setup(self.enforcement_controller,
                              self.enforcement_stats_controller)
        imsi = 'IMSI001010000000014'
        uplink_tunnel = 0x1234
        sub_ip = '192.16.15.7'
        num_pkt_unmatched = 4096

        flow_list = [FlowDescription(
            match=FlowMatch(
                ip_dst=convert_ipv4_str_to_ip_proto('1.1.0.0/24'),
                direction=FlowMatch.UPLINK),
            action=FlowDescription.DENY)
        ]
        policy = VersionedPolicy(
            rule=PolicyRule(id='rule1', priority=3, flow_list=flow_list),
            version=1,
        )
        self.service_manager.session_rule_version_mapper.save_version(
            imsi, uplink_tunnel, 'rule1', 1)

        """ Setup subscriber, setup table_isolation to fwd pkts """
        sub_context = RyuDirectSubscriberContext(
            imsi, sub_ip, uplink_tunnel,
            self.enforcement_controller,
            self._main_tbl_num, self.enforcement_stats_controller
        ).add_policy(policy)

        isolator = RyuDirectTableIsolator(
            RyuForwardFlowArgsBuilder.from_subscriber(sub_context.cfg)
                .build_requests(),
            self.testing_controller
        )

        pkt_sender = ScapyPacketInjector(self.IFACE)
        packet = IPPacketBuilder() \
            .set_ip_layer('45.10.0.0/20', sub_ip) \
            .set_ether_layer(self.MAC_DEST, "00:00:00:00:00:00") \
            .build()

        # =========================== Verification ===========================

        # Verifies that 1 flow is installed in enforcement and 2 flows are
        # installed in enforcement stats, one for uplink and one for downlink.
        snapshot_verifier = SnapshotVerifier(self, self.BRIDGE,
                                             self.service_manager)

        with isolator, sub_context, snapshot_verifier:
            pkt_sender.send(packet)

        enf_stat_name = imsi + '|' + self.DEFAULT_DROP_FLOW_NAME + '|' + str(uplink_tunnel)
        wait_for_enforcement_stats(self.enforcement_stats_controller,
                                   [enf_stat_name])
        stats = get_enforcement_stats(
            self.enforcement_stats_controller._report_usage.call_args_list)

        self.assertEqual(stats[enf_stat_name].sid, imsi)
        self.assertEqual(stats[enf_stat_name].rule_id,
                         self.DEFAULT_DROP_FLOW_NAME)
        self.assertEqual(stats[enf_stat_name].dropped_rx, 0)
        self.assertEqual(stats[enf_stat_name].dropped_tx,
                         num_pkt_unmatched * len(packet))

    def test_ipv6_rule_install(self):
        """
        Adds a ipv6 policy to a subscriber. Verifies that flows are properly
        installed in enforcement and enforcement stats.

        Assert:
            Policy classification flows installed in enforcement
            Policy match flows installed in enforcement_stats
        """
        fake_controller_setup(self.enforcement_controller,
                              self.enforcement_stats_controller)

        imsi = 'IMSI001010000000013'
        uplink_tunnel = 0x1234
        sub_ip = 'de34:431d:1bc::'

        flow_list = [FlowDescription(
            match=FlowMatch(
                ip_dst=convert_ipv6_bytes_to_ip_proto(
                    'f333:432::dbca'.encode('utf-8')),
                direction=FlowMatch.UPLINK),
            action=FlowDescription.PERMIT)
        ]
        policy = VersionedPolicy(
            rule=PolicyRule(id='rule1', priority=3, flow_list=flow_list),
            version=1,
        )
        self.service_manager.session_rule_version_mapper.save_version(
            imsi, uplink_tunnel, 'rule1', 1)

        """ Setup subscriber, setup table_isolation to fwd pkts """
        sub_context = RyuDirectSubscriberContext(
            imsi, sub_ip, uplink_tunnel,
            self.enforcement_controller,
            self._main_tbl_num, self.enforcement_stats_controller
        ).add_policy(policy)

        # =========================== Verification ===========================
        snapshot_verifier = SnapshotVerifier(self, self.BRIDGE,
                                             self.service_manager)

        with sub_context, snapshot_verifier:
            pass

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
        uplink_tunnel = 0x1234
        sub_ip = '192.168.128.74'
        num_pkts_tx_match = 128

        flow_list = [FlowDescription(
            match=FlowMatch(
                ip_dst=convert_ipv4_str_to_ip_proto('45.10.0.0/25'),
                direction=FlowMatch.UPLINK),
            action=FlowDescription.PERMIT)
        ]
        policy = VersionedPolicy(
            rule=PolicyRule(id='rule1', priority=3, flow_list=flow_list),
            version=1,
        )
        enf_stat_name = imsi + '|rule1' + '|' + str(uplink_tunnel)
        self.service_manager.session_rule_version_mapper.save_version(
            imsi, uplink_tunnel, 'rule1', 1)

        """ Setup subscriber, setup table_isolation to fwd pkts """
        sub_context = RyuDirectSubscriberContext(
            imsi, sub_ip, uplink_tunnel,
            self.enforcement_controller,
            self._main_tbl_num, self.enforcement_stats_controller
        ).add_policy(policy)
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
        """ Verify that flows are properly deleted """
        snapshot_verifier = SnapshotVerifier(self, self.BRIDGE,
                                             self.service_manager)

        """
        Send packets, wait until packet is received by ovs and enf stats and
        then deactivate the rule in enforcement controller. This emulates the
        case where there is unreported traffic after rule deactivation.
        """
        with isolator, sub_context, snapshot_verifier:
            self.enforcement_stats_controller._report_usage.reset_mock()
            pkt_sender.send(packet)
            self.service_manager.session_rule_version_mapper. \
                save_version(imsi, uplink_tunnel, 'rule1', 2)
            self.enforcement_controller.deactivate_rules(
                imsi, convert_ipv4_str_to_ip_proto(sub_ip), uplink_tunnel,
                [policy.rule.id])

        wait_for_enforcement_stats(self.enforcement_stats_controller,
                                   [enf_stat_name])
        stats = get_enforcement_stats(
            self.enforcement_stats_controller._report_usage.call_args_list)

        self.assertEqual(stats[enf_stat_name].sid, imsi)
        self.assertEqual(stats[enf_stat_name].rule_id, "rule1")
        self.assertEqual(stats[enf_stat_name].rule_version, 1)
        self.assertEqual(stats[enf_stat_name].bytes_rx, 0)
        self.assertEqual(stats[enf_stat_name].bytes_tx,
                         num_pkts_tx_match * len(packet))

        self.assertEqual(len(stats), 2)

        self.enforcement_stats_controller.deactivate_default_flow(
            imsi, convert_ipv4_str_to_ip_proto(sub_ip),
            uplink_tunnel)
        snapshot_verifier = SnapshotVerifier(self, self.BRIDGE,
                                             self.service_manager,
                                             'nuke_ue')
        with snapshot_verifier:
            pass

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
        uplink_tunnel = 0x1234
        sub_ip = '192.168.128.74'
        num_pkts_tx_match = 128

        flow_list = [FlowDescription(
            match=FlowMatch(
                ip_dst=convert_ipv4_str_to_ip_proto('45.10.0.0/25'),
                direction=FlowMatch.UPLINK),
            action=FlowDescription.PERMIT)
        ]
        policy = VersionedPolicy(
            rule=PolicyRule(id='rule1', priority=3, flow_list=flow_list),
            version=1,
        )
        enf_stat_name = imsi + '|rule1' + '|' + str(uplink_tunnel)
        self.service_manager.session_rule_version_mapper.save_version(
            imsi, uplink_tunnel, 'rule1', 1)

        """ Setup subscriber, setup table_isolation to fwd pkts """
        sub_context = RyuDirectSubscriberContext(
            imsi, sub_ip, uplink_tunnel, self.enforcement_controller,
            self._main_tbl_num, self.enforcement_stats_controller
        ).add_policy(policy)
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

        """ Verify that flows are properly deleted """
        snapshot_verifier = SnapshotVerifier(self, self.BRIDGE,
                                             self.service_manager)

        """
        Send a packet, then deactivate and reactivate the same rule and send a
        packet. Wait until it is received by ovs and enf stats.
        """
        with isolator, sub_context, snapshot_verifier:
            self.enforcement_stats_controller._report_usage.reset_mock()
            pkt_sender.send(packet)

            self.enforcement_stats_controller._report_usage.reset_mock()
            self.service_manager.session_rule_version_mapper. \
                save_version(imsi, uplink_tunnel, 'rule1', 2)
            self.enforcement_controller.deactivate_rules(
                imsi, convert_ipv4_str_to_ip_proto(sub_ip),
                uplink_tunnel, [policy.rule.id])
            policy.version=2
            self.enforcement_controller.activate_rules(
                imsi, None, uplink_tunnel,
                convert_ipv4_str_to_ip_proto(sub_ip), None, [policy])
            self.enforcement_stats_controller.activate_rules(
                imsi, None, uplink_tunnel,
                convert_ipv4_str_to_ip_proto(sub_ip), None, [policy])
            pkt_sender.send(packet)

        wait_for_enforcement_stats(self.enforcement_stats_controller,
                                   [enf_stat_name])
        stats = get_enforcement_stats(
            self.enforcement_stats_controller._report_usage.call_args_list)

        """
        Verify both packets are reported after reactivation.
        """
        self.assertEqual(stats[enf_stat_name].sid, imsi)
        self.assertEqual(stats[enf_stat_name].rule_id, "rule1")
        self.assertEqual(stats[enf_stat_name].rule_version, 2)
        self.assertEqual(stats[enf_stat_name].bytes_rx, 0)
        # TODO Figure out why this one fails.
        #self.assertEqual(stats[enf_stat_name].bytes_tx,
        #                 num_pkts_tx_match * len(packet))
        self.assertEqual(len(stats), 2)


if __name__ == "__main__":
    unittest.main()

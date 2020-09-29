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
from concurrent.futures import Future
from unittest.mock import MagicMock

import warnings
from lte.protos.mconfig.mconfigs_pb2 import PipelineD
from lte.protos.policydb_pb2 import FlowDescription, FlowMatch, PolicyRule, \
    RedirectInformation
from magma.pipelined.app.enforcement import EnforcementController
from magma.pipelined.bridge_util import BridgeTools
from magma.pipelined.policy_converters import convert_ipv4_str_to_ip_proto, \
    convert_ipv6_bytes_to_ip_proto
from magma.pipelined.tests.app.packet_builder import IPPacketBuilder, \
    TCPPacketBuilder, GTPUHeaderPacketBuilder, UDPPacketBuilder
from magma.pipelined.tests.app.packet_injector import ScapyPacketInjector
from magma.pipelined.tests.app.start_pipelined import PipelinedController, \
    TestSetup
from magma.pipelined.tests.app.subscriber import RyuDirectSubscriberContext
from magma.pipelined.tests.app.table_isolation import RyuDirectTableIsolator, \
    RyuForwardFlowArgsBuilder
from magma.pipelined.tests.pipelined_test_util import create_service_manager, \
    get_enforcement_stats, start_ryu_app_thread,  stop_ryu_app_thread, \
    wait_after_send, wait_for_enforcement_stats, FlowTest, SnapshotVerifier, \
    fake_controller_setup
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
                'enforcement': {
                    'poll_interval': 2,
                    'default_drop_flow_name': self.DEFAULT_DROP_FLOW_NAME
                },
                'nat_iface': 'eth2',
                'enodeb_iface': 'eth1',
                'qos': {'enable': False},
                'clean_restart': True,
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

        self.enforcement_stats_controller._policy_dict = self._static_rule_dict
        self.enforcement_stats_controller._report_usage = MagicMock()

        self.enforcement_controller._policy_dict = self._static_rule_dict
        self.enforcement_controller._redirect_manager._save_redirect_entry = \
            MagicMock()

    def tearDown(self):
        stop_ryu_app_thread(self.thread)
        BridgeTools.destroy_bridge(self.BRIDGE)

    def test_ng_subscriber_policy(self):
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
        session_version = 2

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
            PolicyRule(id='tx_match', priority=3, flow_list=flow_list1),
            PolicyRule(id='rx_match', priority=5, flow_list=flow_list2)
        ]
        enf_stat_name = [imsi + '|tx_match' + '|' + sub_ip,
                         imsi + '|rx_match' + '|' + sub_ip]

        #Only one update is required
        self.service_manager.session_rule_version_mapper.\
                   ng_update_rules_version(imsi, session_version)

        """ Setup subscriber, setup table_isolation to fwd pkts """
        self._static_rule_dict[policies[0].id] = policies[0]
        self._static_rule_dict[policies[1].id] = policies[1]

        sub_context = RyuDirectSubscriberContext(
            imsi, sub_ip, self.enforcement_controller,
            self._main_tbl_num, self.enforcement_stats_controller,
            match_teid=100, action_teid=200
        ).add_static_rule(policies[0].id).add_static_rule(policies[1].id)
        isolator = RyuDirectTableIsolator(
            RyuForwardFlowArgsBuilder.from_subscriber(sub_context.cfg)
                                     .build_requests(),
            self.testing_controller
        )

        """ Create 2 sets of packets, for policry rule1&2 """
        pkt_sender = ScapyPacketInjector(self.IFACE)
        packet1 = GTPUHeaderPacketBuilder() \
            .set_ether_layer(self.MAC_DEST, "00:00:00:00:00:00") \
            .set_ip_layer('45.10.0.0/20', sub_ip) \
            .set_udp_layer(2152, 2152) \
            .set_gtp_u_header_layer(100, 20, 255)\
            .build(IP(src='1.1.1.1', dst='2.2.2.2'))

        packet2 = UDPPacketBuilder() \
            .set_ether_layer(self.MAC_DEST, "00:00:00:00:00:00") \
            .set_ip_layer(sub_ip, '45.10.0.0/20') \
            .set_udp_layer(1111, 2222) \
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

        self.assertEqual(len(stats), 2)


if __name__ == "__main__":
    unittest.main()

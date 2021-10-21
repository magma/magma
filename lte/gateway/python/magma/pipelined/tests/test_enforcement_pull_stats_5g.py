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
from lte.protos.policydb_pb2 import FlowDescription, FlowMatch, PolicyRule
from magma.pipelined.app.enforcement import EnforcementController
from magma.pipelined.bridge_util import BridgeTools
from magma.pipelined.policy_converters import convert_ipv4_str_to_ip_proto
from magma.pipelined.tests.app.packet_builder import IPPacketBuilder
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
    SnapshotVerifier,
    create_service_manager,
    fake_controller_setup,
    get_enforcement_stats,
    start_ryu_app_thread,
    stop_ryu_app_thread,
    wait_for_enforcement_stats,
)
from scapy.all import IP


class EnforcementPullStatsTest(unittest.TestCase):
    BRIDGE = 'testing_br'
    IFACE = 'testing_br'
    MAC_DEST = "5e:cc:cc:b1:49:4b"
    MAC_2 = "0a:00:27:00:00:02"
    DEFAULT_DROP_FLOW_NAME = 'internal_default_drop_flow_rule'

    def setUp(self):
        """
        Starts the thread which launches ryu apps
        Create a testing bridge, add a port, setup the port interfaces. Then
        launch the ryu apps for testing pipelined. Gets the references
        to apps launched by using futures, mocks the redis policy_dictionary
        of enforcement_controller
        """
        super(EnforcementPullStatsTest, self).setUpClass()
        warnings.simplefilter('ignore')
        self._static_rule_dict = {}
        self.service_manager = create_service_manager([PipelineD.ENFORCEMENT])
        self._tbl_num = self.service_manager.get_table_num(
            EnforcementController.APP_NAME,
        )

        enforcement_controller_reference = Future()
        enforcement_controller_stats = Future()
        testing_controller_reference = Future()

        def mock_thread_safe(cmd, body):
            cmd(body)
        loop_mock = MagicMock()
        loop_mock.call_soon_threadsafe = mock_thread_safe

        test_setup = TestSetup(
            apps=[
                PipelinedController.Enforcement,
                PipelinedController.Enforcement_stats,
                PipelinedController.Testing,
                PipelinedController.StartupFlows,
            ],
            references={
                PipelinedController.Enforcement:
                    enforcement_controller_reference,
                PipelinedController.Enforcement_stats:
                    enforcement_controller_stats,
                PipelinedController.Testing:
                    testing_controller_reference,
                PipelinedController.StartupFlows:
                    Future(),
            },
            config={
                'bridge_name': self.BRIDGE,
                'bridge_ip_address': '192.168.128.1',
                'enforcement': {
                    'poll_interval': 1,
                    'default_drop_flow_name': self.DEFAULT_DROP_FLOW_NAME,
                },
                'nat_iface': 'eth2',
                'enodeb_iface': 'eth1',
                'enable5g_features': True,
                'qos': {'enable': False},
                'clean_restart': True,
                'uplink_port': 20,
                'enable_nat': True,
                'ovs_gtp_port_number': 10,
                'setup_type': 'LTE',
            },
            mconfig=PipelineD(),
            loop=loop_mock,
            service_manager=self.service_manager,
            integ_test=False,
            rpc_stubs={'sessiond': MagicMock()},
        )

        BridgeTools.create_bridge(self.BRIDGE, self.IFACE)

        self.thread = start_ryu_app_thread(test_setup)

        self.enforcement_controller = enforcement_controller_reference.result()
        self.enforcement_stats_controller = enforcement_controller_stats.result()
        self.testing_controller = testing_controller_reference.result()

        self.enforcement_stats_controller._prepare_ruleRecord_report = MagicMock()

    def tearDown(self):
        stop_ryu_app_thread(self.thread)
        BridgeTools.destroy_bridge(self.BRIDGE)

    def test_enforcemnet_stats_rule(self):
        """
        Add QOS policy to enforcement table into OVS.
        """
        fake_controller_setup(
            self.enforcement_controller,
            self.enforcement_stats_controller,
        )
        imsi = 'IMSI001010000000013'
        sub_ip = '192.168.128.30'
        flow_list1 = [
            FlowDescription(
                match=FlowMatch(
                    direction=FlowMatch.UPLINK,
                ),
                action=FlowDescription.PERMIT,
            ),
            FlowDescription(
                match=FlowMatch(
                    ip_dst=convert_ipv4_str_to_ip_proto("192.168.0.0/24"),
                    direction=FlowMatch.DOWNLINK,
                ),
                action=FlowDescription.PERMIT,
            ),
        ]
        self.service_manager.session_rule_version_mapper.save_version(
            imsi, convert_ipv4_str_to_ip_proto(sub_ip), "rule1", 1,
        )
        policy = VersionedPolicy(
            rule=PolicyRule(id='rule1', priority=65530, flow_list=flow_list1),
            version=1,
        )
        sub_context = RyuDirectSubscriberContext(
            imsi, sub_ip, self.enforcement_controller,
            self._tbl_num, self.enforcement_stats_controller, local_f_teid_ng=100,
        ).add_policy(policy)
        isolator = RyuDirectTableIsolator(
            RyuForwardFlowArgsBuilder.from_subscriber(sub_context.cfg)
                                     .build_requests(),
            self.testing_controller,
        )

        snapshot_verifier = SnapshotVerifier(
            self, self.BRIDGE,
            self.service_manager,
        )

        """ Create 2 sets of packets, for policry rule1&2 """
        pkt_sender = ScapyPacketInjector(self.BRIDGE)

        packet1 = IPPacketBuilder() \
            .set_ip_layer('192.168.60.141/24', sub_ip) \
            .set_ether_layer(self.MAC_DEST, self.MAC_2) \
            .build()

        packet2 = IPPacketBuilder()\
            .set_ip_layer(sub_ip, '192.168.60.141/24')\
            .set_ether_layer(self.MAC_DEST, self.MAC_2)\
            .build()

        """ Send packet, wait until pkts are received by ovs and enf stats """
        with isolator, sub_context, snapshot_verifier:
            pkt_sender.send(packet1)
            pkt_sender.send(packet2)

        enf_stat_name = imsi + '|rule1' + '|' + sub_ip + '|' + "1"

        wait_for_enforcement_stats(
            self.enforcement_stats_controller,
            [enf_stat_name], flag5g=1,
        )
        stats = get_enforcement_stats(
            self.enforcement_stats_controller._prepare_ruleRecord_report.call_args_list,
        )

        self.assertEqual(stats[enf_stat_name].sid, imsi)
        self.assertEqual(stats[enf_stat_name].rule_id, "rule1")
        self.assertEqual(stats[enf_stat_name].teid, 100)
        self.assertEqual(stats[enf_stat_name].rule_version, 1)
        self.assertEqual(stats[enf_stat_name].ue_ipv4, sub_ip)
        self.assertEqual(stats[enf_stat_name].bytes_rx, 5120)
        self.assertEqual(
            stats[enf_stat_name].bytes_tx,
            256 * len(packet1),
        )
        self.assertEqual(len(stats), 2)

    def test_enforcemnet_stats_multiple_rules(self):
        """
        Add QOS policy to enforcement table into OVS.
        """
        fake_controller_setup(
            self.enforcement_controller,
            self.enforcement_stats_controller,
        )
        imsi = 'IMSI001010000000032'
        sub_ip = '192.168.128.45'
        flow_list1 = [
            FlowDescription(
                match=FlowMatch(
                    direction=FlowMatch.UPLINK,
                ),
                action=FlowDescription.PERMIT,
            ),
            FlowDescription(
                match=FlowMatch(
                    ip_dst=convert_ipv4_str_to_ip_proto("192.168.0.0/24"),
                    direction=FlowMatch.DOWNLINK,
                ),
                action=FlowDescription.PERMIT,
            ),
        ]
        flow_list2 = [
            FlowDescription(
                match=FlowMatch(
                    direction=FlowMatch.UPLINK,
                ),
                action=FlowDescription.PERMIT,
            ),
            FlowDescription(
                match=FlowMatch(
                    ip_dst=convert_ipv4_str_to_ip_proto("192.168.0.0/16"),
                    direction=FlowMatch.DOWNLINK,
                ),
                action=FlowDescription.PERMIT,
            ),
        ]

        self.service_manager.session_rule_version_mapper.save_version(
            imsi, convert_ipv4_str_to_ip_proto(sub_ip), "rule1", 1,
        )
        self.service_manager.session_rule_version_mapper.save_version(
            imsi, convert_ipv4_str_to_ip_proto(sub_ip), 'rule2', 2,
        )
        policies = [
            VersionedPolicy(
                rule=PolicyRule(id='rule1', priority=65530, flow_list=flow_list1),
                version=1,
            ),
            VersionedPolicy(
                rule=PolicyRule(id='rule2', priority=65530, flow_list=flow_list2),
                version=2,
            ),
        ]

        sub_context = RyuDirectSubscriberContext(
            imsi, sub_ip, self.enforcement_controller,
            self._tbl_num, self.enforcement_stats_controller, local_f_teid_ng=40,
        ).add_policy(policies[0]) \
         .add_policy(policies[1])

        isolator = RyuDirectTableIsolator(
            RyuForwardFlowArgsBuilder.from_subscriber(sub_context.cfg)
                                     .build_requests(),
            self.testing_controller,
        )

        snapshot_verifier = SnapshotVerifier(
            self, self.BRIDGE,
            self.service_manager,
        )

        """ Create 2 sets of packets, for policry rule1&2 """
        pkt_sender = ScapyPacketInjector(self.BRIDGE)

        packet1 = IPPacketBuilder() \
            .set_ip_layer('192.168.60.14/24', sub_ip) \
            .set_ether_layer(self.MAC_DEST, self.MAC_2) \
            .build()

        packet2 = IPPacketBuilder()\
            .set_ip_layer(sub_ip, '192.168.60.14/24')\
            .set_ether_layer(self.MAC_DEST, self.MAC_2)\
            .build()

        """ Send packet, wait until pkts are received by ovs and enf stats """
        with isolator, sub_context, snapshot_verifier:
            pkt_sender.send(packet1)
            pkt_sender.send(packet2)

        enf_stat_name = [
            imsi + '|rule1' + '|' + sub_ip + '|' + "1",
            imsi + '|rule2' + '|' + sub_ip + '|' + "2",
        ]

        wait_for_enforcement_stats(
            self.enforcement_stats_controller,
            enf_stat_name, flag5g=1,
        )
        stats = get_enforcement_stats(
            self.enforcement_stats_controller._prepare_ruleRecord_report.call_args_list,
        )

        self.assertEqual(stats[enf_stat_name[0]].sid, imsi)
        self.assertEqual(stats[enf_stat_name[0]].rule_id, "rule1")
        self.assertEqual(stats[enf_stat_name[0]].teid, 40)
        self.assertEqual(stats[enf_stat_name[0]].rule_version, 1)
        self.assertEqual(stats[enf_stat_name[0]].ue_ipv4, sub_ip)
        self.assertEqual(stats[enf_stat_name[0]].bytes_rx, 0)

        self.assertEqual(stats[enf_stat_name[1]].sid, imsi)
        self.assertEqual(stats[enf_stat_name[1]].rule_id, "rule2")
        self.assertEqual(stats[enf_stat_name[1]].teid, 40)
        self.assertEqual(stats[enf_stat_name[1]].rule_version, 2)
        self.assertEqual(stats[enf_stat_name[1]].bytes_tx, 8704)

        total_bytes_pkt2 = 256 * len(packet2[IP])
        self.assertEqual(stats[enf_stat_name[1]].bytes_rx, total_bytes_pkt2)

        self.assertEqual(len(stats), 3)


if __name__ == "__main__":
    unittest.main()

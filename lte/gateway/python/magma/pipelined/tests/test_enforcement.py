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
from typing import List

from lte.protos.mconfig.mconfigs_pb2 import PipelineD
from lte.protos.mobilityd_pb2 import IPAddress
from lte.protos.pipelined_pb2 import VersionedPolicy
from lte.protos.policydb_pb2 import (
    FlowDescription,
    FlowMatch,
    HeaderEnrichment,
    PolicyRule,
)
from magma.pipelined.app import he
from magma.pipelined.app.enforcement import EnforcementController
from magma.pipelined.bridge_util import BridgeTools
from magma.pipelined.policy_converters import (
    convert_ipv4_str_to_ip_proto,
    convert_ipv6_bytes_to_ip_proto,
    flow_match_to_magma_match,
)
from magma.pipelined.tests.app.flow_query import RyuDirectFlowQuery as FlowQuery
from magma.pipelined.tests.app.packet_builder import (
    IPPacketBuilder,
    IPv6PacketBuilder,
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
    FlowVerifier,
    PktsToSend,
    SnapshotVerifier,
    SubTest,
    create_service_manager,
    fake_controller_setup,
    start_ryu_app_thread,
    stop_ryu_app_thread,
    wait_after_send,
)


def mocked_activate_he_urls_for_ue(ip: IPAddress, rule_id, urls: List[str], imsi: str, msisdn: str):
    return True


def mocked_deactivate_he_urls_for_ue(ip: IPAddress, rule_id):
    pass


class EnforcementTableTest(unittest.TestCase):
    BRIDGE = 'testing_br'
    IFACE = 'testing_br'
    MAC_DEST = "5e:cc:cc:b1:49:4b"
    he_controller_reference = Future()
    VETH = 'tveth1'
    VETH_NS = 'tveth1_ns'
    PROXY_PORT = '16'

    @classmethod
    def setUpClass(cls):
        """
        Starts the thread which launches ryu apps

        Create a testing bridge, add a port, setup the port interfaces. Then
        launch the ryu apps for testing pipelined. Gets the references
        to apps launched by using futures, mocks the redis policy_dictionary
        of enforcement_controller
        """
        super(EnforcementTableTest, cls).setUpClass()
        warnings.simplefilter('ignore')
        cls.service_manager = create_service_manager([PipelineD.ENFORCEMENT], ['proxy'])
        cls._tbl_num = cls.service_manager.get_table_num(
            EnforcementController.APP_NAME)
        BridgeTools.create_bridge(cls.BRIDGE, cls.IFACE)

        BridgeTools.create_veth_pair(cls.VETH, cls.VETH_NS)
        BridgeTools.add_ovs_port(cls.BRIDGE, cls.VETH, cls.PROXY_PORT)

        enforcement_controller_reference = Future()
        testing_controller_reference = Future()
        he.activate_he_urls_for_ue = mocked_activate_he_urls_for_ue
        he.deactivate_he_urls_for_ue = mocked_deactivate_he_urls_for_ue

        test_setup = TestSetup(
            apps=[PipelinedController.Enforcement,
                  PipelinedController.HeaderEnrichment,
                  PipelinedController.Testing,
                  PipelinedController.StartupFlows],
            references={
                PipelinedController.Enforcement:
                    enforcement_controller_reference,
                PipelinedController.HeaderEnrichment:
                    cls.he_controller_reference,
                PipelinedController.Testing:
                    testing_controller_reference,
                PipelinedController.StartupFlows:
                    Future(),
            },
            config={
                'bridge_name': cls.BRIDGE,
                'bridge_ip_address': '192.168.128.1',
                'nat_iface': 'eth2',
                'enodeb_iface': 'eth1',
                'qos': {'enable': False},
                'clean_restart': True,
                'uplink_port': 20,
                'proxy_port_name': cls.VETH,
                'enable_nat': True,
                'ovs_gtp_port_number': 10,
                'setup_type': 'LTE',
            },
            mconfig=PipelineD(),
            loop=None,
            service_manager=cls.service_manager,
            integ_test=False
        )


        cls.thread = start_ryu_app_thread(test_setup)

        cls.enforcement_controller = enforcement_controller_reference.result()
        cls.testing_controller = testing_controller_reference.result()

    @classmethod
    def tearDownClass(cls):
        stop_ryu_app_thread(cls.thread)
        BridgeTools.destroy_bridge(cls.BRIDGE)

    def test_subscriber_policy(self):
        """
        Add policy to subscriber, send 4096 packets

        Assert:
            Packets are properly matched with the 'simple_match' policy
            Send /20 (4096) packets, match /16 (256) packets
        """
        fake_controller_setup(self.enforcement_controller)
        imsi = 'IMSI010000000088888'
        sub_ip = '192.168.128.74'
        uplink_tunnel = 0x1234
        flow_list1 = [FlowDescription(
            match=FlowMatch(
                ip_dst=convert_ipv4_str_to_ip_proto('45.10.0.0/24'),
                direction=FlowMatch.UPLINK),
            action=FlowDescription.PERMIT)
        ]
        policies = [
            VersionedPolicy(
                rule=PolicyRule(id='simple_match', priority=2,flow_list=flow_list1),
                version=1,
            ),
        ]
        pkts_matched = 256
        pkts_sent = 4096

        # ============================ Subscriber ============================
        sub_context = RyuDirectSubscriberContext(
            imsi, sub_ip, uplink_tunnel,
            self.enforcement_controller, self._tbl_num
        ).add_policy(policies[0])
        isolator = RyuDirectTableIsolator(
            RyuForwardFlowArgsBuilder.from_subscriber(sub_context.cfg)
                                     .build_requests(),
            self.testing_controller
        )
        pkt_sender = ScapyPacketInjector(self.IFACE)
        packet = IPPacketBuilder()\
            .set_ip_layer('45.10.0.0/20', sub_ip)\
            .set_ether_layer(self.MAC_DEST, "00:00:00:00:00:00")\
            .build()
        flow_query = FlowQuery(
            self._tbl_num, self.testing_controller,
            match=flow_match_to_magma_match(flow_list1[0].match)
        )

        # =========================== Verification ===========================
        # Verify aggregate table stats, subscriber 1 'simple_match' pkt count
        flow_verifier = FlowVerifier([
            FlowTest(FlowQuery(self._tbl_num, self.testing_controller),
                     pkts_sent),
            FlowTest(flow_query, pkts_matched)
        ], lambda: wait_after_send(self.testing_controller))
        snapshot_verifier = SnapshotVerifier(self, self.BRIDGE,
                                             self.service_manager)

        with isolator, sub_context, flow_verifier, snapshot_verifier:
            pkt_sender.send(packet)

        flow_verifier.verify()

    def test_subscriber_ipv6_policy(self):
        """
        Add policy to subscriber, send 4096 packets

        Assert:
            Packets are properly matched with the 'simple_match' policy
            Send /20 (4096) packets, match /16 (256) packets
        """
        fake_controller_setup(self.enforcement_controller)
        imsi = 'IMSI010000000088888'
        sub_ip = 'de34:431d:1bc::'
        uplink_tunnel = 0x1234
        flow_list1 = [FlowDescription(
            match=FlowMatch(
                ip_dst=convert_ipv6_bytes_to_ip_proto(
                    'f333:432::dbca'.encode('utf-8')),
                direction=FlowMatch.UPLINK),
            action=FlowDescription.PERMIT)
        ]
        policies = [
            VersionedPolicy(
                rule=PolicyRule(id='simple_match', priority=2, flow_list=flow_list1),
                version=1,
            ),
        ]

        # ============================ Subscriber ============================
        sub_context = RyuDirectSubscriberContext(
            imsi, sub_ip, uplink_tunnel,
            self.enforcement_controller, self._tbl_num
        ).add_policy(policies[0])
        isolator = RyuDirectTableIsolator(
            RyuForwardFlowArgsBuilder.from_subscriber(sub_context.cfg)
                .build_requests(),
            self.testing_controller
        )
        pkt_sender = ScapyPacketInjector(self.IFACE)
        packet = IPv6PacketBuilder() \
            .set_ip_layer('f333:432::dbca', sub_ip) \
            .set_ether_layer(self.MAC_DEST, "00:00:00:00:00:00") \
            .build()

        # =========================== Verification ===========================
        snapshot_verifier = SnapshotVerifier(self, self.BRIDGE,
                                             self.service_manager)

        with isolator, sub_context, snapshot_verifier:
            pkt_sender.send(packet)

    def test_invalid_subscriber(self):
        """
        Try to apply an invalid policy to a subscriber, should log and error

        Assert:
            Only 1 flow gets added to the table (drop flow)
        """
        fake_controller_setup(self.enforcement_controller)
        imsi = 'IMSI000000000000001'
        sub_ip = '192.168.128.45'
        uplink_tunnel = 0x1234
        flow_list = [FlowDescription(
            match=FlowMatch(
                ip_src=convert_ipv4_str_to_ip_proto('9999.0.0.0/24')),
            action=FlowDescription.DENY
        )]
        policy = \
            VersionedPolicy(
                rule=PolicyRule(id='invalid', priority=2, flow_list=flow_list),
                version=1,
            )
        invalid_sub_context = RyuDirectSubscriberContext(
            imsi, sub_ip,  uplink_tunnel,
            self.enforcement_controller, self._tbl_num).add_policy(policy)
        isolator = RyuDirectTableIsolator(
            RyuForwardFlowArgsBuilder.from_subscriber(invalid_sub_context.cfg)
                                     .build_requests(),
            self.testing_controller
        )
        flow_query = FlowQuery(self._tbl_num, self.testing_controller)
        num_flows_start = len(flow_query.lookup())
        snapshot_verifier = SnapshotVerifier(self, self.BRIDGE,
                                             self.service_manager)

        with isolator, invalid_sub_context, snapshot_verifier:
            wait_after_send(self.testing_controller)
            num_flows_final = len(flow_query.lookup())

        self.assertEqual(num_flows_final - num_flows_start, 0)

    def test_subscriber_two_policies(self):
        """
        Add 2 policies to subscriber

        Assert:
            Packets are properly matched with the 'match' policy
            The total packet delta in the table is from the above match
        """
        fake_controller_setup(self.enforcement_controller)
        imsi = 'IMSI208950000000001'
        sub_ip = '192.168.128.74'
        uplink_tunnel = 0x1234
        flow_list1 = [FlowDescription(
            match=FlowMatch(
                ip_src=convert_ipv4_str_to_ip_proto('15.0.0.0/24'),
                direction=FlowMatch.DOWNLINK),
            action=FlowDescription.DENY)
        ]
        flow_list2 = [FlowDescription(
            match=FlowMatch(ip_proto=6, direction=FlowMatch.UPLINK),
            action=FlowDescription.PERMIT)
        ]

        policies = [
            VersionedPolicy(
                rule=PolicyRule(id='match', priority=2, flow_list=flow_list1),
                version=1,
            ),
            VersionedPolicy(
                rule=PolicyRule(id='no_match', priority=2, flow_list=flow_list2),
                version=1,
            ),
        ]
        pkts_sent = 42

        # ============================ Subscriber ============================
        sub_context = RyuDirectSubscriberContext(imsi, sub_ip,uplink_tunnel,
            self.enforcement_controller, self._tbl_num) \
            .add_policy(policies[0]).add_policy(policies[1])
        isolator = RyuDirectTableIsolator(
            RyuForwardFlowArgsBuilder.from_subscriber(sub_context.cfg)
                                     .build_requests(),
            self.testing_controller
        )
        pkt_sender = ScapyPacketInjector(self.IFACE)
        packet = IPPacketBuilder()\
            .set_ip_layer(sub_ip, '15.0.0.8')\
            .set_ether_layer(self.MAC_DEST, "00:00:00:00:00:00")\
            .build()
        flow_query = FlowQuery(
            self._tbl_num, self.testing_controller,
            match=flow_match_to_magma_match(flow_list1[0].match)
        )

        # =========================== Verification ===========================
        # Verify aggregate table stats, subscriber 1 'match' rule pkt count
        flow_verifier = FlowVerifier([
            FlowTest(FlowQuery(self._tbl_num, self.testing_controller),
                     pkts_sent),
            FlowTest(flow_query, pkts_sent)
        ], lambda: wait_after_send(self.testing_controller))
        snapshot_verifier = SnapshotVerifier(self, self.BRIDGE,
                                             self.service_manager)

        with isolator, sub_context, flow_verifier, snapshot_verifier:
            pkt_sender.send(packet, pkts_sent)

        flow_verifier.verify()

    def test_two_subscribers(self):
        """
        Add 2 subscribers at the same time

        Assert:
            For subcriber1 the packets are matched to the proper policy
            For subcriber2 the packets are matched to the proper policy
            The total packet delta in the table is from the above matches
        """
        uplink_tunnel = 0x1234
        fake_controller_setup(self.enforcement_controller)
        pkt_sender = ScapyPacketInjector(self.IFACE)
        ip_match = [FlowDescription(
            match=FlowMatch(ip_src=convert_ipv4_str_to_ip_proto('8.8.8.0/24'),
                            direction=1),
            action=1)
        ]
        tcp_match = [FlowDescription(
            match=FlowMatch(ip_proto=6, direction=FlowMatch.DOWNLINK),
            action=FlowDescription.DENY)
        ]

        policy = \
            VersionedPolicy(
                rule=PolicyRule(id='t', priority=2, flow_list=ip_match),
                version=1,
            )
        # =========================== Subscriber 1 ===========================
        sub_context1 = RyuDirectSubscriberContext(
            'IMSI208950001111111', '192.168.128.5',  uplink_tunnel,
            self.enforcement_controller, self._tbl_num
        ).add_policy(policy)
        isolator1 = RyuDirectTableIsolator(
            RyuForwardFlowArgsBuilder.from_subscriber(sub_context1.cfg)
                                     .build_requests(),
            self.testing_controller
        )
        packet_ip = IPPacketBuilder()\
            .set_ether_layer(self.MAC_DEST, "00:00:00:00:00:00")\
            .set_ip_layer(sub_context1.cfg.ip, '8.8.8.8')\
            .build()
        s1_pkts_sent = 29
        pkts_to_send = [PktsToSend(packet_ip, s1_pkts_sent)]
        flow_query1 = FlowQuery(
            self._tbl_num, self.testing_controller,
            match=flow_match_to_magma_match(ip_match[0].match)
        )
        s1 = SubTest(
            sub_context1, isolator1, FlowTest(flow_query1, s1_pkts_sent)
        )

        # =========================== Subscriber 2 ===========================
        sub_context2 = RyuDirectSubscriberContext(
            'IMSI911500451242001', '192.168.128.100', uplink_tunnel,
            self.enforcement_controller, self._tbl_num
        ).add_policy(
            VersionedPolicy(
                rule=PolicyRule(id='qqq', priority=2, flow_list=tcp_match),
                version=1,
            )
        )
        isolator2 = RyuDirectTableIsolator(
            RyuForwardFlowArgsBuilder.from_subscriber(sub_context2.cfg)
                                     .build_requests(),
            self.testing_controller
        )
        packet_tcp = TCPPacketBuilder()\
            .set_ether_layer(self.MAC_DEST, "00:00:00:00:00:00")\
            .set_ip_layer(sub_context2.cfg.ip, '15.0.0.8')\
            .build()
        s2_pkts_sent = 18
        pkts_to_send.append(PktsToSend(packet_tcp, s2_pkts_sent))
        flow_query2 = FlowQuery(
            self._tbl_num, self.testing_controller,
            match=flow_match_to_magma_match(tcp_match[0].match)
        )
        s2 = SubTest(
            sub_context2, isolator2, FlowTest(flow_query2, s2_pkts_sent)
        )

        # =========================== Verification ===========================
        # Verify aggregate table stats, subscriber 1 & 2 flows packet matches
        pkts = s1_pkts_sent + s2_pkts_sent
        flow_verifier = FlowVerifier([
            FlowTest(FlowQuery(self._tbl_num, self.testing_controller), pkts),
            s1.flowtest_list,
            s2.flowtest_list
        ], lambda: wait_after_send(self.testing_controller))
        snapshot_verifier = SnapshotVerifier(self, self.BRIDGE,
                                             self.service_manager)

        with s1.isolator, s1.context, s2.isolator, s2.context, flow_verifier, \
             snapshot_verifier:
            for pkt in pkts_to_send:
                pkt_sender.send(pkt.pkt, pkt.num)

        flow_verifier.verify()


if __name__ == "__main__":
    unittest.main()

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

from lte.protos.mconfig.mconfigs_pb2 import PipelineD
from magma.pipelined.app.access_control import AccessControlController
from magma.pipelined.bridge_util import BridgeTools
from magma.pipelined.openflow.magma_match import MagmaMatch
from magma.pipelined.openflow.registers import Direction
from magma.pipelined.tests.app.flow_query import RyuDirectFlowQuery as FlowQuery
from magma.pipelined.tests.app.packet_builder import IPPacketBuilder
from magma.pipelined.tests.app.packet_injector import ScapyPacketInjector
from magma.pipelined.tests.app.start_pipelined import (
    PipelinedController,
    TestSetup,
)
from magma.pipelined.tests.app.subscriber import (
    SubContextConfig,
    default_ambr_config,
)
from magma.pipelined.tests.app.table_isolation import (
    RyuDirectTableIsolator,
    RyuForwardFlowArgsBuilder,
)
from magma.pipelined.tests.pipelined_test_util import (
    FlowTest,
    FlowVerifier,
    assert_bridge_snapshot_match,
    create_service_manager,
    start_ryu_app_thread,
    stop_ryu_app_thread,
    wait_after_send,
)
from ryu.lib.packet import ether_types


class AccessControlTestLTE(unittest.TestCase):
    BRIDGE = 'testing_br'
    IFACE = 'testing_br'
    MAC_DEST = "5e:cc:cc:b1:49:4b"
    BRIDGE_IP = '192.168.128.1'
    INBOUND_TEST_IP = '127.0.0.1'
    OUTBOUND_TEST_IP = '127.1.0.1'
    BOTH_DIR_TEST_IP = '127.2.0.1'

    @classmethod
    def setUpClass(cls):
        """
        Starts the thread which launches ryu apps

        Create a testing bridge, add a port, setup the port interfaces. Then
        launch the ryu apps for testing pipelined. Gets the references
        to apps launched by using futures.
        """
        super(AccessControlTestLTE, cls).setUpClass()
        warnings.simplefilter('ignore')
        cls.service_manager = create_service_manager([],
            ['access_control'])
        cls._tbl_num = cls.service_manager.get_table_num(
            AccessControlController.APP_NAME)

        access_control_controller_reference = Future()
        testing_controller_reference = Future()
        test_setup = TestSetup(
            apps=[PipelinedController.AccessControl,
                  PipelinedController.Testing,
                  PipelinedController.StartupFlows],
            references={
                PipelinedController.AccessControl:
                    access_control_controller_reference,
                PipelinedController.Testing:
                    testing_controller_reference,
                PipelinedController.StartupFlows:
                    Future(),
            },
            config={
                'setup_type': 'LTE',
                'allow_unknown_arps': False,
                'bridge_name': cls.BRIDGE,
                'bridge_ip_address': cls.BRIDGE_IP,
                'nat_iface': 'eth2',
                'enodeb_iface': 'eth1',
                'qos': {'enable': False},
                'access_control': {
                    'ip_blocklist': [
                        {
                            'ip': cls.INBOUND_TEST_IP,
                            'direction': 'inbound',
                        },
                        {
                            'ip': cls.OUTBOUND_TEST_IP,
                            'direction': 'outbound',
                        },
                        {
                            'ip': cls.BOTH_DIR_TEST_IP,
                        },
                    ]
                },
                'clean_restart': True,
            },
            mconfig=PipelineD(
                allowed_gre_peers=[{'ip': '1.2.3.4/24', 'key': 123}],
            ),
            loop=None,
            service_manager=cls.service_manager,
            integ_test=False,
        )

        BridgeTools.create_bridge(cls.BRIDGE, cls.IFACE)

        cls.thread = start_ryu_app_thread(test_setup)
        cls.access_control_controller = \
            access_control_controller_reference.result()
        cls.testing_controller = testing_controller_reference.result()

    @classmethod
    def tearDownClass(cls):
        stop_ryu_app_thread(cls.thread)
        BridgeTools.destroy_bridge(cls.BRIDGE)

    def test_inbound_ip_match(self):
        """
        Inbound ip match test, checks that packets are properly matched when
        the inbound traffic matches an ip in the blocklist.

        Assert:
            Both packets are matched
            Ip match flows are added
        """
        # Set up subscribers
        sub = SubContextConfig('IMSI001010000000013', '192.168.128.74', 0x1234,
                               default_ambr_config, self._tbl_num)

        isolator = RyuDirectTableIsolator(
            RyuForwardFlowArgsBuilder.from_subscriber(sub).build_requests(),
            self.testing_controller,
        )

        # Set up packets
        pkt_sender = ScapyPacketInjector(self.BRIDGE)
        packets = [
            self._build_default_ip_packet(self.INBOUND_TEST_IP, sub.ip),
            self._build_default_ip_packet(self.BOTH_DIR_TEST_IP, sub.ip),
        ]

        # Check if these flows were added (queries should return flows)
        inbound_flow_queries = [
            FlowQuery(self._tbl_num, self.testing_controller,
                      match=MagmaMatch(eth_type=ether_types.ETH_TYPE_IP,
                                       direction=Direction.OUT,
                                       ipv4_dst=self.INBOUND_TEST_IP)),
            FlowQuery(self._tbl_num, self.testing_controller,
                      match=MagmaMatch(eth_type=ether_types.ETH_TYPE_IP,
                                       direction=Direction.OUT,
                                       ipv4_dst=self.BOTH_DIR_TEST_IP)),
        ]

        # =========================== Verification ===========================
        # packets matched, ip match flows installed
        flow_verifier = FlowVerifier(
            [
                FlowTest(
                    FlowQuery(self._tbl_num, self.testing_controller), 2),
            ] + [FlowTest(query, 1, flow_count=1) for query in
                 inbound_flow_queries],
            lambda: wait_after_send(
                self.testing_controller)
        )

        with isolator, flow_verifier:
            for packet in packets:
                pkt_sender.send(packet)

        flow_verifier.verify()
        assert_bridge_snapshot_match(self,
                                     self.BRIDGE,
                                     self.service_manager)

    def test_outbound_ip_match(self):
        """
        Outbound ip match test, checks that packets are properly matched when
        the outbound traffic matches an ip in the blocklist.

        Assert:
            Both packets are matched
            Ip match flows are added
        """
        # Set up subscribers
        sub = SubContextConfig('IMSI001010000000013', '192.168.128.74', 0x1234,
                               default_ambr_config, self._tbl_num)

        isolator = RyuDirectTableIsolator(
            RyuForwardFlowArgsBuilder.from_subscriber(sub).build_requests(),
            self.testing_controller,
        )

        # Set up packets
        pkt_sender = ScapyPacketInjector(self.BRIDGE)
        packets = [
            self._build_default_ip_packet(sub.ip, self.OUTBOUND_TEST_IP),
            self._build_default_ip_packet(sub.ip, self.BOTH_DIR_TEST_IP),
        ]

        # Check if these flows were added (queries should return flows)
        outbound_flow_queries = [
            FlowQuery(self._tbl_num, self.testing_controller,
                      match=MagmaMatch(eth_type=ether_types.ETH_TYPE_IP,
                                       direction=Direction.IN,
                                       ipv4_src=self.OUTBOUND_TEST_IP)),
            FlowQuery(self._tbl_num, self.testing_controller,
                      match=MagmaMatch(eth_type=ether_types.ETH_TYPE_IP,
                                       direction=Direction.IN,
                                       ipv4_src=self.BOTH_DIR_TEST_IP)),
        ]

        # =========================== Verification ===========================
        # packets matched, ip match flows installed
        flow_verifier = FlowVerifier(
            [
                FlowTest(
                    FlowQuery(self._tbl_num, self.testing_controller), 2),
            ] + [FlowTest(query, 1, flow_count=1) for query in
                 outbound_flow_queries],
            lambda: wait_after_send(
                self.testing_controller)
        )

        with isolator, flow_verifier:
            for packet in packets:
                pkt_sender.send(packet)

        flow_verifier.verify()
        assert_bridge_snapshot_match(self,
                                     self.BRIDGE,
                                     self.service_manager)

    def test_no_match(self):
        """
        No match test, checks that packets are not matched when
        the there is no match to the ip and direction in the blocklist.

        Assert:
            Both packets are not matched
            Ip match flows are added
        """
        # Set up subscribers
        sub = SubContextConfig('IMSI001010000000013', '192.168.128.74', 0x1234,
                               default_ambr_config, self._tbl_num)

        isolator = RyuDirectTableIsolator(
            RyuForwardFlowArgsBuilder.from_subscriber(sub).build_requests(),
            self.testing_controller,
        )

        # Set up packets. The directions of the packets are opposite of the
        # installed match flow, so there should not matches.
        pkt_sender = ScapyPacketInjector(self.BRIDGE)
        packets = [
            self._build_default_ip_packet(self.OUTBOUND_TEST_IP, sub.ip),
            self._build_default_ip_packet(sub.ip, self.INBOUND_TEST_IP),
        ]

        # Check if these flows were added (queries should return flows)
        outbound_flow_queries = [
            FlowQuery(self._tbl_num, self.testing_controller,
                      match=MagmaMatch(eth_type=ether_types.ETH_TYPE_IP,
                                       direction=Direction.OUT,
                                       ipv4_dst=self.INBOUND_TEST_IP)),
            FlowQuery(self._tbl_num, self.testing_controller,
                      match=MagmaMatch(eth_type=ether_types.ETH_TYPE_IP,
                                       direction=Direction.IN,
                                       ipv4_src=self.OUTBOUND_TEST_IP)),
        ]

        # =========================== Verification ===========================
        # packets are not matched, ip match flows installed
        flow_verifier = FlowVerifier(
            [
                FlowTest(
                    FlowQuery(self._tbl_num, self.testing_controller), 2),
            ] + [FlowTest(query, 0, flow_count=1) for query in
                 outbound_flow_queries],
            lambda: wait_after_send(
                self.testing_controller)
        )

        with isolator, flow_verifier:
            for packet in packets:
                pkt_sender.send(packet)

        flow_verifier.verify()
        assert_bridge_snapshot_match(self,
                                     self.BRIDGE,
                                     self.service_manager)

    def _build_default_ip_packet(self, dst, src):
        return IPPacketBuilder() \
            .set_ip_layer(dst, src) \
            .set_ether_layer(self.MAC_DEST, "00:00:00:00:00:00") \
            .build()


class AccessControlTestCWF(unittest.TestCase):
    BRIDGE = 'testing_br'
    IFACE = 'testing_br'
    MAC_DEST = "5e:cc:cc:b1:49:4b"
    BRIDGE_IP = '192.168.128.1'
    INBOUND_TEST_IP = '127.0.0.1'
    OUTBOUND_TEST_IP = '127.1.0.1'
    BOTH_DIR_TEST_IP = '127.2.0.1'

    @classmethod
    def setUpClass(cls):
        """
        Starts the thread which launches ryu apps

        Create a testing bridge, add a port, setup the port interfaces. Then
        launch the ryu apps for testing pipelined. Gets the references
        to apps launched by using futures.
        """
        super(AccessControlTestCWF, cls).setUpClass()
        warnings.simplefilter('ignore')
        cls.service_manager = create_service_manager([], ['access_control'])
        cls._tbl_num = cls.service_manager.get_table_num(
            AccessControlController.APP_NAME)

        access_control_controller_reference = Future()
        testing_controller_reference = Future()
        test_setup = TestSetup(
            apps=[PipelinedController.AccessControl,
                  PipelinedController.Testing,
                  PipelinedController.StartupFlows],
            references={
                PipelinedController.AccessControl:
                    access_control_controller_reference,
                PipelinedController.Testing:
                    testing_controller_reference,
                PipelinedController.StartupFlows:
                    Future(),
            },
            config={
                'setup_type': 'CWF',
                'allow_unknown_arps': False,
                'bridge_name': cls.BRIDGE,
                'bridge_ip_address': cls.BRIDGE_IP,
                'internal_ip_subnet': '192.168.0.0/16',
                'nat_iface': 'eth2',
                'enodeb_iface': 'eth1',
                'enable_queue_pgm': False,
                'clean_restart': True,
                'access_control': {
                    'ip_blocklist': [
                        {
                            'ip': cls.INBOUND_TEST_IP,
                            'direction': 'inbound',
                        },
                        {
                            'ip': cls.OUTBOUND_TEST_IP,
                            'direction': 'outbound',
                        },
                        {
                            'ip': cls.BOTH_DIR_TEST_IP,
                        },
                    ]
                }
            },
            mconfig=PipelineD(
                allowed_gre_peers=[{'ip': '2.2.2.2/24'},
                                   {'ip': '1.2.3.4/24', 'key': 123}],
            ),
            loop=None,
            service_manager=cls.service_manager,
            integ_test=False,
        )

        BridgeTools.create_bridge(cls.BRIDGE, cls.IFACE)

        cls.thread = start_ryu_app_thread(test_setup)
        cls.access_control_controller = \
            access_control_controller_reference.result()
        cls.testing_controller = testing_controller_reference.result()

    @classmethod
    def tearDownClass(cls):
        stop_ryu_app_thread(cls.thread)
        BridgeTools.destroy_bridge(cls.BRIDGE)

    def test_gre_peer_rules(self):
        """
        Inbound ip match test, checks that packets are properly matched when
        the inbound traffic matches an ip in the blocklist.

        Assert:
            Both packets are matched
            Ip match flows are added
        """
        assert_bridge_snapshot_match(self,
                                     self.BRIDGE,
                                     self.service_manager)


if __name__ == "__main__":
    unittest.main()

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

from magma.pipelined.app.ue_mac import UEMacAddressController
from magma.pipelined.bridge_util import BridgeTools
from magma.pipelined.openflow.magma_match import MagmaMatch
from magma.pipelined.tests.app.flow_query import RyuDirectFlowQuery as FlowQuery
from magma.pipelined.tests.app.packet_builder import EtherPacketBuilder
from magma.pipelined.tests.app.packet_injector import ScapyPacketInjector
from magma.pipelined.tests.app.start_pipelined import (
    PipelinedController,
    TestSetup,
)
from magma.pipelined.tests.pipelined_test_util import (
    FlowTest,
    FlowVerifier,
    SnapshotVerifier,
    create_service_manager,
    start_ryu_app_thread,
    stop_ryu_app_thread,
    wait_after_send,
)


class UEMacAddressTest(unittest.TestCase):
    BRIDGE = 'testing_br'
    IFACE = 'testing_br'
    UE_MAC_1 = '5e:cc:cc:b1:49:4b'
    UE_MAC_2 = '5e:cc:cc:aa:aa:fe'
    BRIDGE_IP = '192.168.130.1'
    DPI_PORT = 'mon1'
    DPI_IP = '1.1.1.1'

    @classmethod
    def setUpClass(cls):
        """
        Starts the thread which launches ryu apps

        Create a testing bridge, add a port, setup the port interfaces. Then
        launch the ryu apps for testing pipelined. Gets the references
        to apps launched by using futures.
        """
        super(UEMacAddressTest, cls).setUpClass()
        warnings.simplefilter('ignore')
        cls.service_manager = create_service_manager([], ['ue_mac'])
        cls._tbl_num = cls.service_manager.get_table_num(
            UEMacAddressController.APP_NAME)
        ue_mac_controller_reference = Future()
        testing_controller_reference = Future()
        test_setup = TestSetup(
            apps=[PipelinedController.UEMac,
                  PipelinedController.Testing,
                  PipelinedController.StartupFlows],
            references={
                PipelinedController.UEMac:
                    ue_mac_controller_reference,
                PipelinedController.Testing:
                    testing_controller_reference,
                PipelinedController.Arp:
                    Future(),
                PipelinedController.StartupFlows:
                    Future(),
            },
            config={
                'setup_type': 'CWF',
                'allow_unknown_arps': False,
                'bridge_name': cls.BRIDGE,
                'bridge_ip_address': cls.BRIDGE_IP,
                'internal_ip_subnet': '192.168.0.0/16',
                'ovs_gtp_port_number': 32768,
                'clean_restart': True,
                'dpi': {
                    'enabled': False,
                    'mon_port': 'mon1',
                    'mon_port_number': 32769,
                    'idle_timeout': 42,
                },
            },
            mconfig=None,
            loop=None,
            service_manager=cls.service_manager,
            integ_test=False,
        )

        BridgeTools.create_bridge(cls.BRIDGE, cls.IFACE)
        BridgeTools.create_internal_iface(cls.BRIDGE, cls.DPI_PORT,
                                          cls.DPI_IP)

        cls.thread = start_ryu_app_thread(test_setup)
        cls.ue_mac_controller = ue_mac_controller_reference.result()
        cls.testing_controller = testing_controller_reference.result()

    @classmethod
    def tearDownClass(cls):
        stop_ryu_app_thread(cls.thread)
        BridgeTools.destroy_bridge(cls.BRIDGE)

    def test_add_two_subscribers(self):
        """
           Add UE MAC flows for two subscribers
        """
        imsi_1 = 'IMSI010000000088888'
        imsi_2 = 'IMSI010000000099999'
        other_mac = '5e:cc:cc:b1:aa:aa'

        # Add subscriber with UE MAC address """
        self.ue_mac_controller.add_ue_mac_flow(imsi_1, self.UE_MAC_1)
        self.ue_mac_controller.add_ue_mac_flow(imsi_2, self.UE_MAC_2)

        # Create a set of packets
        pkt_sender = ScapyPacketInjector(self.BRIDGE)

        # Only send downlink as the pkt_sender sends pkts from in_port=LOCAL
        downlink_packet1 = EtherPacketBuilder() \
            .set_ether_layer(self.UE_MAC_1, other_mac) \
            .build()
        downlink_packet2 = EtherPacketBuilder() \
            .set_ether_layer(self.UE_MAC_2, other_mac) \
            .build()

        # Check if these flows were added (queries should return flows)
        flow_queries = [
            FlowQuery(self._tbl_num, self.testing_controller,
                      match=MagmaMatch(eth_dst=self.UE_MAC_1)),
            FlowQuery(self._tbl_num, self.testing_controller,
                      match=MagmaMatch(eth_dst=self.UE_MAC_2))
        ]

        # =========================== Verification ===========================
        # Verify 5 flows installed and 2 total pkts matched (one for each UE)
        flow_verifier = FlowVerifier(
            [FlowTest(FlowQuery(self._tbl_num, self.testing_controller), 2, 5)]
            + [FlowTest(query, 1, 1) for query in flow_queries],
            lambda: wait_after_send(self.testing_controller))

        snapshot_verifier = SnapshotVerifier(self, self.BRIDGE,
                                             self.service_manager)

        with flow_verifier, snapshot_verifier:
            pkt_sender.send(downlink_packet1)
            pkt_sender.send(downlink_packet2)

        flow_verifier.verify()

    def test_delete_one_subscriber(self):
        """
            Delete one of the existing subscribers
        """

        imsi_1 = 'IMSI010000000088888'
        other_mac = '5e:cc:cc:b1:aa:aa'

        # Delete subscriber with UE MAC_1 address
        self.ue_mac_controller.delete_ue_mac_flow(imsi_1, self.UE_MAC_1)

        # Create a set of packets
        pkt_sender = ScapyPacketInjector(self.BRIDGE)

        # Only send downlink as the pkt_sender sends pkts from in_port=LOCAL
        removed_ue_packet = EtherPacketBuilder() \
            .set_ether_layer(self.UE_MAC_1, other_mac) \
            .build()
        remaining_ue_packet = EtherPacketBuilder() \
            .set_ether_layer(self.UE_MAC_2, other_mac) \
            .build()

        # Ensure the first query doesn't match anything
        # And the second query still does
        flow_queries = [
            FlowQuery(self._tbl_num, self.testing_controller,
                      match=MagmaMatch(eth_dst=self.UE_MAC_1)),
            FlowQuery(self._tbl_num, self.testing_controller,
                      match=MagmaMatch(eth_dst=self.UE_MAC_2))
        ]

        # =========================== Verification ===========================
        # Verify flows installed and 1 total pkt matched
        flow_verifier = FlowVerifier(
            [
                FlowTest(FlowQuery(self._tbl_num, self.testing_controller), 1,
                         3),
                FlowTest(flow_queries[0], 0, 0),
                FlowTest(flow_queries[1], 1, 1),
            ], lambda: wait_after_send(self.testing_controller))

        snapshot_verifier = SnapshotVerifier(self, self.BRIDGE,
                                             self.service_manager)

        with flow_verifier, snapshot_verifier:
            pkt_sender.send(removed_ue_packet)
            pkt_sender.send(remaining_ue_packet)

        flow_verifier.verify()


if __name__ == "__main__":
    unittest.main()

"""
Copyright (c) 2019-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""

import unittest
import warnings
from concurrent.futures import Future

from magma.pipelined.app.ue_mac import UEMacAddressController
from magma.pipelined.app.inout import INGRESS, EGRESS
from magma.pipelined.app.ue_mac import UEMacAddressController
from magma.pipelined.tests.app.packet_builder import EtherPacketBuilder, \
    UDPPacketBuilder, ARPPacketBuilder
from magma.pipelined.tests.app.packet_injector import ScapyPacketInjector
from magma.pipelined.tests.app.start_pipelined import (
    TestSetup,
    PipelinedController,
)
from magma.pipelined.openflow.magma_match import MagmaMatch
from magma.pipelined.tests.app.flow_query import RyuDirectFlowQuery \
    as FlowQuery
from magma.pipelined.bridge_util import BridgeTools
from magma.pipelined.tests.pipelined_test_util import (
    start_ryu_app_thread,
    stop_ryu_app_thread,
    create_service_manager,
    wait_after_send,
    FlowVerifier,
    FlowTest,
    SnapshotVerifier,
)


class UEMacAddressTest(unittest.TestCase):
    BRIDGE = 'testing_br'
    IFACE = 'testing_br'
    UE_MAC_1 = '5e:cc:cc:b1:49:4b'
    UE_MAC_2 = '5e:cc:cc:aa:aa:fe'
    BRIDGE_IP = '192.168.130.1'

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
        cls.service_manager = create_service_manager([], include_ue_mac=True)
        cls._tbl_num = cls.service_manager.get_table_num(
            UEMacAddressController.APP_NAME)
        cls._ingress_tbl_num = cls.service_manager.get_table_num(INGRESS)
        cls._egress_tbl_num = cls.service_manager.get_table_num(EGRESS)

        inout_controller_reference = Future()
        ue_mac_controller_reference = Future()
        testing_controller_reference = Future()
        test_setup = TestSetup(
            apps=[PipelinedController.InOut,
                  PipelinedController.UEMac,
                  PipelinedController.Testing],
            references={
                PipelinedController.InOut:
                    inout_controller_reference,
                PipelinedController.UEMac:
                    ue_mac_controller_reference,
                PipelinedController.Testing:
                    testing_controller_reference
            },
            config={
                'bridge_name': cls.BRIDGE,
                'bridge_ip_address': cls.BRIDGE_IP,
                'ovs_gtp_port_number': 32768,
            },
            mconfig=None,
            loop=None,
            service_manager=cls.service_manager,
            integ_test=False,
        )

        BridgeTools.create_bridge(cls.BRIDGE, cls.IFACE)

        cls.thread = start_ryu_app_thread(test_setup)
        cls.ue_mac_controller = ue_mac_controller_reference.result()
        cls.testing_controller = testing_controller_reference.result()

    @classmethod
    def tearDownClass(cls):
        stop_ryu_app_thread(cls.thread)
        BridgeTools.destroy_bridge(cls.BRIDGE)

    def test_passthrough_rules(self):
        """
           Add UE MAC flows for two subscribers
        """
        imsi_1 = 'IMSI010000000088888'
        other_mac = '5e:cc:cc:b1:aa:aa'

        # Add subscriber with UE MAC address """
        self.ue_mac_controller.add_ue_mac_flow(imsi_1, self.UE_MAC_1)

        # Create a set of packets
        pkt_sender = ScapyPacketInjector(self.BRIDGE)

        # Only send downlink as the pkt_sender sends pkts from in_port=LOCAL
        downlink_packet1 = EtherPacketBuilder() \
            .set_ether_layer(self.UE_MAC_1, other_mac) \
            .build()
        dhcp_packet = UDPPacketBuilder() \
            .set_ether_layer(self.UE_MAC_1, other_mac) \
            .set_ip_layer('151.42.41.122', '1.1.1.1') \
            .set_udp_layer(67, 68) \
            .build()
        dns_packet = UDPPacketBuilder() \
            .set_ether_layer(self.UE_MAC_1, other_mac) \
            .set_ip_layer('151.42.41.122', '1.1.1.1') \
            .set_udp_layer(53, 32795) \
            .build()
        arp_packet = ARPPacketBuilder() \
            .set_arp_layer('151.42.41.122') \
            .set_arp_hwdst(self.UE_MAC_1) \
            .set_arp_src(other_mac, '1.1.1.1') \
            .build()


        # Check if these flows were added (queries should return flows)
        flow_queries = [
            FlowQuery(self._tbl_num, self.testing_controller,
                      match=MagmaMatch(eth_dst=self.UE_MAC_1))
        ]

        # =========================== Verification ===========================
        # Verify 9 flows installed for ue_mac table (3 pkts matched)
        #        4 flows installed for inout (3 pkts matched)
        #        2 flows installed (2 pkts matches)
        flow_verifier = FlowVerifier(
            [
                FlowTest(FlowQuery(self._tbl_num,
                                   self.testing_controller), 4, 9),
                FlowTest(FlowQuery(self._ingress_tbl_num,
                                   self.testing_controller), 4, 4),
                FlowTest(FlowQuery(self._egress_tbl_num,
                                   self.testing_controller), 3, 2),
                FlowTest(flow_queries[0], 3, 4),
            ], lambda: wait_after_send(self.testing_controller))

        snapshot_verifier = SnapshotVerifier(self, self.BRIDGE,
                                             self.service_manager)

        with flow_verifier, snapshot_verifier:
            pkt_sender.send(downlink_packet1)
            pkt_sender.send(dhcp_packet)
            pkt_sender.send(dns_packet)
            pkt_sender.send(arp_packet)

        flow_verifier.verify()


if __name__ == "__main__":
    unittest.main()

"""
Copyright (c) 2020-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""
import time
import unittest
import warnings
from concurrent.futures import Future

from lte.protos.mconfig.mconfigs_pb2 import PipelineD
from magma.pipelined.app.arp import ArpController
from magma.pipelined.tests.app.packet_injector import ScapyPacketInjector
from magma.pipelined.bridge_util import BridgeTools
from magma.pipelined.tests.app.packet_builder import ARPPacketBuilder
from magma.pipelined.tests.app.start_pipelined import TestSetup, \
    PipelinedController
from magma.pipelined.openflow.registers import IMSI_REG, DIRECTION_REG,\
    Direction
from magma.pipelined.tests.app.table_isolation import RyuDirectTableIsolator, \
    RyuForwardFlowArgsBuilder
from magma.pipelined.tests.pipelined_test_util import start_ryu_app_thread, \
    stop_ryu_app_thread, create_service_manager, wait_after_send, \
    SnapshotVerifier


def _pkt_total(stats):
    return sum(n.packets for n in stats)


class ArpTableTest(unittest.TestCase):
    BRIDGE = 'testing_br'
    IFACE = 'testing_br'
    MAC_DEST = "5e:cc:cc:b1:49:4b"
    BRIDGE_IP = '192.168.128.1'
    UE_BLOCK = '192.168.128.0/24'
    UE_MAC = '5e:cc:cc:b1:49:4b'
    UE_IP = '192.168.128.22'
    OTHER_MAC = '0a:00:27:00:00:02'
    OTHER_IP = '1.2.3.4'
    VIRTUAL_MAC = '0a:00:FF:00:00:FF'
    MTR_IP = '5.6.7.8'
    MTR_MAC = "FF:EE:DD:CC:49:4b"

    @unittest.mock.patch('netifaces.ifaddresses',
                return_value=[[{'addr': '00:11:22:33:44:55'}]])
    @unittest.mock.patch('netifaces.AF_LINK', 0)
    def setUp(self, *_):
        """
        Starts the thread which launches ryu apps

        Create a testing bridge, add a port, setup the port interfaces. Then
        launch the ryu apps for testing pipelined. Gets the references
        to apps launched by using futures.
        """
        cls = ArpTableTest
        #super(ArpTableTest, cls).setUpClass()
        warnings.simplefilter('ignore')
        cls.service_manager = create_service_manager([], ['arpd'])
        cls._tbl_num = cls.service_manager.get_table_num(ArpController.APP_NAME)

        arp_controller_reference = Future()
        testing_controller_reference = Future()
        test_setup = TestSetup(
            apps=[
                PipelinedController.Arp,
                PipelinedController.Testing,
                PipelinedController.StartupFlows
            ],
            references={
                PipelinedController.Arp:
                    arp_controller_reference,
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
                'ovs_gtp_port_number': 32768,
                'virtual_mac': cls.VIRTUAL_MAC,
                'local_ue_eth_addr': True,
                'quota_check_ip': '1.2.3.4',
                'clean_restart': True,
                'non_nat': True,
                'mtr_ip': cls.MTR_IP,
                'mtr_mac': cls.MTR_MAC,
            },
            mconfig=PipelineD(
                ue_ip_block=cls.UE_BLOCK,
            ),
            loop=None,
            service_manager=cls.service_manager,
            integ_test=False,
        )

        BridgeTools.create_bridge(cls.BRIDGE, cls.IFACE)

        cls.thread = start_ryu_app_thread(test_setup)
        cls.arp_controller = arp_controller_reference.result()
        cls.testing_controller = testing_controller_reference.result()

    def tearDown(self):
        cls = ArpTableTest
        stop_ryu_app_thread(cls.thread)
        BridgeTools.destroy_bridge(cls.BRIDGE)

    def test_uplink_to_ue_arp(self):
        """
        Verify that a UPLINK->UE arp request with IMSI set is properly matched
        """
        pkt_sender = ScapyPacketInjector(self.IFACE)
        arp_packet = ARPPacketBuilder() \
            .set_ether_layer(self.UE_MAC, self.OTHER_MAC) \
            .set_arp_layer(self.UE_IP) \
            .set_arp_hwdst(self.UE_MAC) \
            .set_arp_src(self.OTHER_MAC, self.OTHER_IP) \
            .build()

        dlink_args = RyuForwardFlowArgsBuilder(self._tbl_num) \
            .set_eth_match(eth_dst=self.UE_MAC, eth_src=self.OTHER_MAC) \
            .set_reg_value(DIRECTION_REG, Direction.IN) \
            .set_reg_value(IMSI_REG, 123) \
            .build_requests()
        isolator = RyuDirectTableIsolator(dlink_args, self.testing_controller)

        snapshot_verifier = SnapshotVerifier(self, self.BRIDGE,
                                             self.service_manager)

        with isolator, snapshot_verifier:
            pkt_sender.send(arp_packet, count=4)
            wait_after_send(self.testing_controller)

    def test_ue_to_uplink_arp(self):
        """
        Verify that a UE->UPLINK arp request is dropped
        MME would not set IMSI reg
        """
        pkt_sender = ScapyPacketInjector(self.IFACE)
        arp_packet = ARPPacketBuilder() \
            .set_ether_layer(self.OTHER_MAC, self.UE_MAC) \
            .set_arp_layer(self.OTHER_IP) \
            .set_arp_hwdst(self.OTHER_MAC) \
            .set_arp_src(self.UE_MAC, self.UE_IP) \
            .build()

        uplink_args = RyuForwardFlowArgsBuilder(self._tbl_num) \
            .set_eth_match(eth_src=self.UE_MAC, eth_dst=self.OTHER_MAC) \
            .set_reg_value(DIRECTION_REG, Direction.OUT) \
            .build_requests()
        isolator = RyuDirectTableIsolator(uplink_args, self.testing_controller)

        snapshot_verifier = SnapshotVerifier(self, self.BRIDGE,
                                             self.service_manager)

        with isolator, snapshot_verifier:
            pkt_sender.send(arp_packet, count=4)
            wait_after_send(self.testing_controller)

    def test_stray_arp_drop(self):
        """
        Verify that an arp that neither UE->UPLINK nor UPLINK->UE is dropped
        """
        pkt_sender = ScapyPacketInjector(self.IFACE)
        arp_packet = ARPPacketBuilder() \
            .set_ether_layer('11:11:11:11:11:1', self.OTHER_MAC) \
            .set_arp_layer(self.OTHER_IP) \
            .set_arp_hwdst(self.OTHER_MAC) \
            .set_arp_src('22:22:22:22:22:22', '1.1.1.1') \
            .build()

        uplink_args = RyuForwardFlowArgsBuilder(self._tbl_num) \
            .set_eth_match(eth_dst='11:11:11:11:11:1', eth_src=self.OTHER_MAC) \
            .set_reg_value(DIRECTION_REG, Direction.OUT) \
            .build_requests()
        isolator = RyuDirectTableIsolator(uplink_args, self.testing_controller)

        snapshot_verifier = SnapshotVerifier(self, self.BRIDGE,
                                             self.service_manager)

        with isolator, snapshot_verifier:
            pkt_sender.send(arp_packet, count=4)
            wait_after_send(self.testing_controller)

    def test_mtr_arp(self):
        """
        Verify that a MTR arp request is handled in ARP responder
        """
        pkt_sender = ScapyPacketInjector(self.IFACE)
        arp_packet = ARPPacketBuilder() \
            .set_ether_layer(self.OTHER_MAC, self.MAC_DEST) \
            .set_arp_layer(self.MTR_IP) \
            .set_arp_hwdst(self.OTHER_MAC) \
            .set_arp_src(self.UE_MAC, self.UE_IP) \
            .build()

        uplink_args = RyuForwardFlowArgsBuilder(self._tbl_num) \
            .set_eth_type_arp() \
            .set_reg_value(DIRECTION_REG, Direction.IN) \
            .build_requests()
        isolator = RyuDirectTableIsolator(uplink_args, self.testing_controller)
        time.sleep(1)
        snapshot_verifier = SnapshotVerifier(self, self.BRIDGE,
                                             self.service_manager)

        with isolator, snapshot_verifier:
            pkt_sender.send(arp_packet, count=4)
            wait_after_send(self.testing_controller)


if __name__ == "__main__":
    unittest.main()

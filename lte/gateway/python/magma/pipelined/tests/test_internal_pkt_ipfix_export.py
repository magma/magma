"""
Copyright (c) 2018-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""

import unittest
from concurrent.futures import Future

import warnings

from lte.protos.mconfig.mconfigs_pb2 import PipelineD
from lte.protos.policydb_pb2 import FlowMatch
from magma.pipelined.app.dpi import DPIController
from magma.pipelined.bridge_util import BridgeTools
from magma.pipelined.tests.app.start_pipelined import PipelinedController, \
    TestSetup
from magma.pipelined.tests.pipelined_test_util import create_service_manager, \
    start_ryu_app_thread, stop_ryu_app_thread, SnapshotVerifier


class InternalPktIpfixExportTest(unittest.TestCase):
    BRIDGE = 'testing_br'
    IFACE = 'testing_br'
    MAC_DEST = "5e:cc:cc:b1:49:4b"
    BRIDGE_IP = '192.168.128.1'

    @classmethod
    def setUpClass(cls):
        """
        Starts the thread which launches ryu apps

        Create a testing bridge, add a port, setup the port interfaces. Then
        launch the ryu apps for testing pipelined. Gets the references
        to apps launched by using futures, mocks the redis policy_dictionary
        of dpi_controller
        """
        super(InternalPktIpfixExportTest, cls).setUpClass()
        warnings.simplefilter('ignore')
        cls._static_rule_dict = {}
        cls.service_manager = create_service_manager(
            [PipelineD.DPI], ['ue_mac', 'ipfix'])
        cls._tbl_num = cls.service_manager.get_table_num(
            DPIController.APP_NAME)

        ue_mac_controller_reference = Future()
        dpi_controller_reference = Future()
        testing_controller_reference = Future()
        test_setup = TestSetup(
            apps=[PipelinedController.UEMac,
                  PipelinedController.DPI,
                  PipelinedController.Testing,
                  PipelinedController.StartupFlows],
            references={
                PipelinedController.UEMac:
                    ue_mac_controller_reference,
                PipelinedController.DPI:
                    dpi_controller_reference,
                PipelinedController.Arp:
                    Future(),
                PipelinedController.Testing:
                    testing_controller_reference,
                PipelinedController.StartupFlows:
                    Future(),
            },
            config={
                'bridge_name': cls.BRIDGE,
                'bridge_ip_address': '192.168.128.1',
                'internal_ip_subnet': '192.168.0.0/16',
                'nat_iface': 'eth2',
                'enodeb_iface': 'eth1',
                'enable_queue_pgm': False,
                'clean_restart': True,
                'setup_type': 'CWF',
                'dpi': {
                    'enabled': False,
                    'mon_port': 'mon1',
                    'mon_port_number': 32769,
                    'idle_timeout': 42,
                },
            },
            mconfig=PipelineD(
                relay_enabled=True
            ),
            loop=None,
            service_manager=cls.service_manager,
            integ_test=False
        )

        BridgeTools.create_bridge(cls.BRIDGE, cls.IFACE)
        cls.thread = start_ryu_app_thread(test_setup)

        cls.ue_mac_controller = ue_mac_controller_reference.result()
        cls.dpi_controller = dpi_controller_reference.result()
        cls.testing_controller = testing_controller_reference.result()

        cls.dpi_controller._policy_dict = cls._static_rule_dict

    @classmethod
    def tearDownClass(cls):
        stop_ryu_app_thread(cls.thread)
        BridgeTools.destroy_bridge(cls.BRIDGE)

    def test_subscriber_policy(self):
        """
        Classify DPI flow, verify internal packet is generated

        Assert:
            snapshots math
        """
        imsi = 'IMSI010000000088888'
        ue_mac = '5e:cc:cc:b1:49:4b'

        self.ue_mac_controller.add_ue_mac_flow(imsi, ue_mac)

        dst_mac = "5e:cc:cc:b1:49:4b"
        flow_match = FlowMatch(
            ip_proto=FlowMatch.IPPROTO_TCP, ipv4_dst='45.10.0.1',
            ipv4_src='1.2.3.0', tcp_dst=80, tcp_src=51115,
            direction=FlowMatch.UPLINK
        )
        self.dpi_controller.add_classify_flow(
            flow_match, 'base.ip.http.facebook', 'tbd', ue_mac, dst_mac)

        snapshot_verifier = SnapshotVerifier(self, self.BRIDGE,
                                             self.service_manager)

        with snapshot_verifier:
            pass


if __name__ == "__main__":
    unittest.main()

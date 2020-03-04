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

from ryu.lib import hub
from lte.protos.mconfig.mconfigs_pb2 import PipelineD
from lte.protos.policydb_pb2 import FlowMatch


from magma.pipelined.tests.app.start_pipelined import (
    TestSetup,
    PipelinedController,
)
from magma.pipelined.bridge_util import BridgeTools
from magma.pipelined.tests.pipelined_test_util import (
    start_ryu_app_thread,
    stop_ryu_app_thread,
    create_service_manager,
    assert_bridge_snapshot_match,
)


class DPITest(unittest.TestCase):
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
        to apps launched by using futures.
        """
        super(DPITest, cls).setUpClass()
        warnings.simplefilter('ignore')
        cls.service_manager = create_service_manager([PipelineD.DPI], [])

        dpi_controller_reference = Future()
        testing_controller_reference = Future()
        test_setup = TestSetup(
            apps=[PipelinedController.DPI,
                  PipelinedController.Testing,
                  PipelinedController.StartupFlows],
            references={
                PipelinedController.DPI:
                    dpi_controller_reference,
                PipelinedController.Testing:
                    testing_controller_reference,
                PipelinedController.StartupFlows:
                    Future(),
            },
            config={
                'bridge_name': cls.BRIDGE,
                'bridge_ip_address': cls.BRIDGE_IP,
                'ovs_gtp_port_number': 32768,
                'clean_restart': True,
                'setup_type': 'LTE',
                'dpi': {
                    'enabled': False,
                    'mon_port': 'mon1',
                    'mon_port_number': 32769,
                },
            },
            mconfig=None,
            loop=None,
            service_manager=cls.service_manager,
            integ_test=False,
        )

        BridgeTools.create_bridge(cls.BRIDGE, cls.IFACE)

        cls.thread = start_ryu_app_thread(test_setup)
        cls.dpi_controller = dpi_controller_reference.result()
        cls.testing_controller = testing_controller_reference.result()

    @classmethod
    def tearDownClass(cls):
        stop_ryu_app_thread(cls.thread)
        BridgeTools.destroy_bridge(cls.BRIDGE)

    def test_add_app_rules(self):
        """
        Test DPI classifier flows are properly added

        Assert:
            2 App types are matched on (`facebook` and `instagram`)
        """
        flow_match1 = FlowMatch(
            ip_proto=FlowMatch.IPPROTO_TCP, ipv4_dst='45.10.0.0/24',
            ipv4_src='1.2.3.0/24', tcp_dst=80, tcp_src=51115,
            direction=FlowMatch.UPLINK
        )
        flow_match2 = FlowMatch(
            ip_proto=FlowMatch.IPPROTO_UDP, ipv4_dst='45.10.0.0/24',
            ipv4_src='1.2.3.0/24', udp_src=111, udp_dst=222,
            direction=FlowMatch.UPLINK
        )
        self.dpi_controller.add_classify_flow(flow_match1, 'facebook')
        self.dpi_controller.add_classify_flow(flow_match2, 'instagram')
        hub.sleep(5)
        assert_bridge_snapshot_match(self, self.BRIDGE, self.service_manager)

    def test_remove_app_rules(self):
        """
        Test DPI classifier flows are properly removed

        Assert:
            Initial state has (`facebook` and `instagram`), remove (`facebook`)
            App type matches on (`instagram`)
        """
        flow_match1 = FlowMatch(
            ip_proto=FlowMatch.IPPROTO_TCP, ipv4_dst='45.10.0.0/24',
            ipv4_src='1.2.3.0/24', tcp_dst=80, tcp_src=51115,
            direction=FlowMatch.UPLINK
        )
        self.dpi_controller.remove_classify_flow(flow_match1)
        hub.sleep(5)
        assert_bridge_snapshot_match(self, self.BRIDGE, self.service_manager)


if __name__ == "__main__":
    unittest.main()

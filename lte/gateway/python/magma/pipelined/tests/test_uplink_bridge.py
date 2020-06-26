"""
Copyright (c) 2020-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""
import unittest
import warnings
from concurrent.futures import Future

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


class UplinkBridgeTest(unittest.TestCase):
    BRIDGE = 'testing_br'
    MAC_DEST = "5e:cc:cc:b1:49:4b"
    BRIDGE_IP = '192.168.128.1'

    UPLINK_BRIDGE = 'up_br0'

    @classmethod
    def setUpClass(cls):
        """
        Starts the thread which launches ryu apps

        Create a testing bridge, add a port, setup the port interfaces. Then
        launch the ryu apps for testing pipelined. Gets the references
        to apps launched by using futures.
        """
        super(UplinkBridgeTest, cls).setUpClass()
        warnings.simplefilter('ignore')
        cls.service_manager = create_service_manager([])

        uplink_bridge_controller_reference = Future()
        testing_controller_reference = Future()
        test_setup = TestSetup(
            apps=[PipelinedController.UplinkBridge,
                  PipelinedController.Testing,
                  PipelinedController.StartupFlows],
            references={
                PipelinedController.UplinkBridge:
                    uplink_bridge_controller_reference,
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
            },
            mconfig=None,
            loop=None,
            service_manager=cls.service_manager,
            integ_test=False,
        )

        BridgeTools.create_bridge(cls.BRIDGE, cls.BRIDGE)
        BridgeTools.create_bridge(cls.UPLINK_BRIDGE, cls.UPLINK_BRIDGE)

        cls.thread = start_ryu_app_thread(test_setup)
        cls.uplink_br_controller = uplink_bridge_controller_reference.result()
        cls.testing_controller = testing_controller_reference.result()

    @classmethod
    def tearDownClass(cls):
        stop_ryu_app_thread(cls.thread)
        BridgeTools.destroy_bridge(cls.BRIDGE)
        BridgeTools.destroy_bridge(cls.UPLINK_BRIDGE)

    def testFlowSnapshotMatch(self):
        assert_bridge_snapshot_match(self, self.UPLINK_BRIDGE, self.service_manager)


class UplinkBridgeWithNonNATTest(unittest.TestCase):
    BRIDGE = 'testing_br'
    MAC_DEST = "5e:cc:cc:b1:49:4b"
    BRIDGE_IP = '192.168.128.1'

    UPLINK_BRIDGE = 'up_br0'
    UPLINK_DHCP = 'test_dhcp0'
    UPLINK_PATCH = 'test_patch_p2'
    UPLINK_ETH_PORT = 'test_eth3'

    @classmethod
    def setUpClass(cls):
        """
        Starts the thread which launches ryu apps

        Create a testing bridge, add a port, setup the port interfaces. Then
        launch the ryu apps for testing pipelined. Gets the references
        to apps launched by using futures.
        """
        super(UplinkBridgeWithNonNATTest, cls).setUpClass()
        warnings.simplefilter('ignore')
        cls.service_manager = create_service_manager([])

        uplink_bridge_controller_reference = Future()
        testing_controller_reference = Future()
        test_setup = TestSetup(
            apps=[PipelinedController.UplinkBridge,
                  PipelinedController.Testing,
                  PipelinedController.StartupFlows],
            references={
                PipelinedController.UplinkBridge:
                    uplink_bridge_controller_reference,
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
                'non_nat': True,
                'uplink_bridge': cls.UPLINK_BRIDGE,
                'uplink_eth_port_name': cls.UPLINK_ETH_PORT,
                'virtual_mac': '02:bb:5e:36:06:4b',
                'uplink_patch': cls.UPLINK_PATCH,
                'uplink_dhcp_port': cls.UPLINK_DHCP,
            },
            mconfig=None,
            loop=None,
            service_manager=cls.service_manager,
            integ_test=False,
        )

        BridgeTools.create_bridge(cls.BRIDGE, cls.BRIDGE)

        # dummy uplink interface
        BridgeTools.create_bridge(cls.UPLINK_BRIDGE, cls.UPLINK_BRIDGE)
        BridgeTools.create_internal_iface(cls.UPLINK_BRIDGE,
                                          cls.UPLINK_DHCP, None)
        BridgeTools.create_internal_iface(cls.UPLINK_BRIDGE,
                                          cls.UPLINK_PATCH, None)
        BridgeTools.create_internal_iface(cls.UPLINK_BRIDGE,
                                          cls.UPLINK_ETH_PORT, None)

        cls.thread = start_ryu_app_thread(test_setup)
        cls.uplink_br_controller = uplink_bridge_controller_reference.result()

        cls.testing_controller = testing_controller_reference.result()

    @classmethod
    def tearDownClass(cls):
        stop_ryu_app_thread(cls.thread)
        BridgeTools.destroy_bridge(cls.BRIDGE)
        BridgeTools.destroy_bridge(cls.UPLINK_BRIDGE)

    def testFlowSnapshotMatch(self):
        # time.sleep(100)
        assert_bridge_snapshot_match(self, self.UPLINK_BRIDGE, self.service_manager,
                                     include_stats=False)


if __name__ == "__main__":
    unittest.main()

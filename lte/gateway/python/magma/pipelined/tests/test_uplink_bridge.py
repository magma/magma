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
import subprocess
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
    get_ovsdb_port_tag,
    get_iface_ipv4,
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
                'enable_nat': True,
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
                'enable_nat': False,
                'uplink_bridge': cls.UPLINK_BRIDGE,
                'uplink_eth_port_name': cls.UPLINK_ETH_PORT,
                'virtual_mac': '02:bb:5e:36:06:4b',
                'uplink_patch': cls.UPLINK_PATCH,
                'uplink_dhcp_port': cls.UPLINK_DHCP,
                'sgi_management_iface_vlan': "",
            },
            mconfig=None,
            loop=None,
            service_manager=cls.service_manager,
            integ_test=False,
        )

        BridgeTools.create_bridge(cls.BRIDGE, cls.BRIDGE)

        # dummy uplink interface
        BridgeTools.create_bridge(cls.UPLINK_BRIDGE, cls.UPLINK_BRIDGE)
        vlan = "10"
        BridgeTools.create_bridge(cls.UPLINK_BRIDGE, cls.UPLINK_BRIDGE)
        subprocess.Popen(["ovs-vsctl", "set", "port", cls.UPLINK_BRIDGE,
                          "tag=" + vlan]).wait()
        assert get_ovsdb_port_tag(cls.UPLINK_BRIDGE) == vlan

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
        cls = self.__class__
        assert_bridge_snapshot_match(self, self.UPLINK_BRIDGE, self.service_manager,
                                     include_stats=False)
        self.assertEqual(get_ovsdb_port_tag(cls.UPLINK_BRIDGE), '[]')


class UplinkBridgeWithNonNATTestVlan(unittest.TestCase):
    BRIDGE = 'testing_br'
    MAC_DEST = "5e:cc:cc:b1:49:4b"
    BRIDGE_IP = '192.168.128.1'

    UPLINK_BRIDGE = 'ut_up_br0'
    UPLINK_DHCP = 'test_dhcp0'
    UPLINK_PATCH = 'test_patch_p2'
    UPLINK_ETH_PORT = 'test_eth3'
    VLAN_TAG='100'

    @classmethod
    def setUpClass(cls):
        """
        Starts the thread which launches ryu apps

        Create a testing bridge, add a port, setup the port interfaces. Then
        launch the ryu apps for testing pipelined. Gets the references
        to apps launched by using futures.
        """
        super(UplinkBridgeWithNonNATTestVlan, cls).setUpClass()
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
                'enable_nat': False,
                'uplink_bridge': cls.UPLINK_BRIDGE,
                'uplink_eth_port_name': cls.UPLINK_ETH_PORT,
                'virtual_mac': '02:bb:5e:36:06:4b',
                'uplink_patch': cls.UPLINK_PATCH,
                'uplink_dhcp_port': cls.UPLINK_DHCP,
                'sgi_management_iface_vlan': cls.VLAN_TAG
            },
            mconfig=None,
            loop=None,
            service_manager=cls.service_manager,
            integ_test=False,
        )

        BridgeTools.create_bridge(cls.BRIDGE, cls.BRIDGE)
        # validate vlan id set
        vlan = "10"
        BridgeTools.create_bridge(cls.UPLINK_BRIDGE, cls.UPLINK_BRIDGE)
        subprocess.Popen(["ovs-vsctl", "set", "port", cls.UPLINK_BRIDGE,
                          "tag=" + vlan]).wait()
        assert get_ovsdb_port_tag(cls.UPLINK_BRIDGE) == vlan

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
        cls = self.__class__
        assert_bridge_snapshot_match(self, self.UPLINK_BRIDGE, self.service_manager,
                                     include_stats=False)

        self.assertEqual(get_ovsdb_port_tag(cls.UPLINK_BRIDGE), cls.VLAN_TAG)

class UplinkBridgeWithNonNATTest_IP_VLAN(unittest.TestCase):
    BRIDGE = 'testing_br'
    MAC_DEST = "5e:cc:cc:b1:49:4b"
    BRIDGE_IP = '192.168.128.1'

    UPLINK_BRIDGE = 'ut_up_br0'
    UPLINK_DHCP = 'test_dhcp0'
    UPLINK_PATCH = 'test_patch_p2'
    UPLINK_ETH_PORT = 'test_eth3'
    VLAN_TAG='100'
    SGi_IP="1.6.5.7"

    @classmethod
    def setUpClass(cls):
        """
        Starts the thread which launches ryu apps

        Create a testing bridge, add a port, setup the port interfaces. Then
        launch the ryu apps for testing pipelined. Gets the references
        to apps launched by using futures.
        """
        super(UplinkBridgeWithNonNATTest_IP_VLAN, cls).setUpClass()
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
                'enable_nat': False,
                'uplink_bridge': cls.UPLINK_BRIDGE,
                'uplink_eth_port_name': cls.UPLINK_ETH_PORT,
                'virtual_mac': '02:bb:5e:36:06:4b',
                'uplink_patch': cls.UPLINK_PATCH,
                'uplink_dhcp_port': cls.UPLINK_DHCP,
                'sgi_management_iface_vlan': cls.VLAN_TAG,
                'sgi_management_iface_ip_addr': cls.SGi_IP
            },
            mconfig=None,
            loop=None,
            service_manager=cls.service_manager,
            integ_test=False,
        )

        BridgeTools.create_bridge(cls.BRIDGE, cls.BRIDGE)
        # validate vlan id set
        vlan = "10"
        BridgeTools.create_bridge(cls.UPLINK_BRIDGE, cls.UPLINK_BRIDGE)
        subprocess.Popen(["ovs-vsctl", "set", "port", cls.UPLINK_BRIDGE,
                          "tag=" + vlan]).wait()
        assert get_ovsdb_port_tag(cls.UPLINK_BRIDGE) == vlan

        set_ip_cmd = ["ip",
                      "addr", "replace",
                      "2.33.44.6",
                      "dev",
                      cls.UPLINK_BRIDGE]
        subprocess.check_call(set_ip_cmd)

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
        cls = self.__class__
        assert_bridge_snapshot_match(self, self.UPLINK_BRIDGE, self.service_manager,
                                     include_stats=False)
        self.assertEqual(get_ovsdb_port_tag(cls.UPLINK_BRIDGE), cls.VLAN_TAG)

        self.assertIn(cls.SGi_IP, get_iface_ipv4(cls.UPLINK_BRIDGE), "ip not found")


if __name__ == "__main__":
    unittest.main()

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
import logging
import subprocess
import unittest
import warnings
from concurrent.futures import Future

from magma.pipelined.bridge_util import BridgeTools
from magma.pipelined.tests.app.start_pipelined import (
    PipelinedController,
    TestSetup,
)
from magma.pipelined.tests.pipelined_test_util import (
    assert_bridge_snapshot_match,
    create_service_manager,
    get_iface_gw_ipv4,
    get_iface_ipv4,
    get_ovsdb_port_tag,
    start_ryu_app_thread,
    stop_ryu_app_thread,
)
from ryu.lib import hub


class UplinkBridgeTest(unittest.TestCase):
    BRIDGE = 'testing_br'
    MAC_DEST = "5e:cc:cc:b1:49:4b"
    BRIDGE_IP = '192.168.128.1'

    UPLINK_BRIDGE = 'upt_br0'

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

    UPLINK_BRIDGE = 'upt_br0'
    UPLINK_DHCP = 'test_dhcp0'
    UPLINK_PATCH = 'test_patch_p2'
    UPLINK_ETH_PORT = 'test_eth3'
    VLAN_DEV_IN = "test_v_in"
    VLAN_DEV_OUT = "test_v_out"

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
                'dev_vlan_in': cls.VLAN_DEV_IN,
                'dev_vlan_out': cls.VLAN_DEV_OUT,
                'ovs_vlan_workaround': False,
                'sgi_management_iface_ip_addr': '1.1.11.1',
            },
            mconfig=None,
            loop=None,
            service_manager=cls.service_manager,
            integ_test=False,
        )

        BridgeTools.create_bridge(cls.BRIDGE, cls.BRIDGE)
        BridgeTools.create_bridge(cls.UPLINK_BRIDGE, cls.UPLINK_BRIDGE)

        BridgeTools.create_veth_pair(cls.VLAN_DEV_IN,
                                     cls.VLAN_DEV_OUT)
        # Add to OVS,
        BridgeTools.add_ovs_port(cls.UPLINK_BRIDGE,
                                 cls.VLAN_DEV_IN, "70")
        BridgeTools.add_ovs_port(cls.UPLINK_BRIDGE,
                                 cls.VLAN_DEV_OUT, "71")

        # dummy uplink interface
        vlan = "10"

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


class UplinkBridgeWithNonNATTestVlan(unittest.TestCase):
    BRIDGE = 'testing_br'
    MAC_DEST = "5e:cc:cc:b1:49:4b"
    BRIDGE_IP = '192.168.128.1'

    UPLINK_BRIDGE = 'upt_br0'
    UPLINK_DHCP = 'test_dhcp0'
    UPLINK_PATCH = 'test_patch_p2'
    UPLINK_ETH_PORT = 'test_eth3'
    VLAN_TAG='100'
    VLAN_DEV_IN = "test_v_in"
    VLAN_DEV_OUT = "test_v_out"

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
                'sgi_management_iface_vlan': cls.VLAN_TAG,
                'dev_vlan_in': cls.VLAN_DEV_IN,
                'dev_vlan_out': cls.VLAN_DEV_OUT,
                'sgi_management_iface_ip_addr': '1.1.11.1',
            },
            mconfig=None,
            loop=None,
            service_manager=cls.service_manager,
            integ_test=False,
        )

        BridgeTools.create_bridge(cls.BRIDGE, cls.BRIDGE)
        BridgeTools.create_bridge(cls.UPLINK_BRIDGE, cls.UPLINK_BRIDGE)

        BridgeTools.create_veth_pair(cls.VLAN_DEV_IN,
                                     cls.VLAN_DEV_OUT)
        # Add to OVS,
        BridgeTools.add_ovs_port(cls.UPLINK_BRIDGE,
                                 cls.VLAN_DEV_IN, "70")
        BridgeTools.add_ovs_port(cls.UPLINK_BRIDGE,
                                 cls.VLAN_DEV_OUT, "71")

        # validate vlan id set
        vlan = "10"
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
        cls = self.__class__
        assert_bridge_snapshot_match(self, self.UPLINK_BRIDGE, self.service_manager,
                                     include_stats=False)


@unittest.skip
# this reset default GW
class UplinkBridgeWithNonNATTest_IP_VLAN(unittest.TestCase):
    BRIDGE = 'testing_br'
    MAC_DEST = "5e:cc:cc:b1:49:4b"
    BRIDGE_IP = '192.168.128.1'

    UPLINK_BRIDGE = 'upt_br0'
    UPLINK_DHCP = 'test_dhcp0'
    UPLINK_PATCH = 'test_patch_p2'
    UPLINK_ETH_PORT = 'test_eth3'
    VLAN_TAG='500'
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
                'sgi_management_iface_ip_addr': cls.SGi_IP,
                'dev_vlan_in': "test_v_in",
                'dev_vlan_out': "test_v_out",
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

        self.assertIn(cls.SGi_IP, get_iface_ipv4(cls.UPLINK_BRIDGE), "ip not found")

@unittest.skip
# this reset default GW
class UplinkBridgeWithNonNATTest_IP_VLAN_GW(unittest.TestCase):
    BRIDGE = 'testing_br'
    MAC_DEST = "5e:cc:cc:b1:49:4b"
    BRIDGE_IP = '192.168.128.1'

    UPLINK_BRIDGE = 'upt_br0'
    UPLINK_DHCP = 'test_dhcp0'
    UPLINK_PATCH = 'test_patch_p2'
    UPLINK_ETH_PORT = 'test_eth3'
    VLAN_TAG = '100'
    SGi_IP = "1.6.5.7/24"
    SGi_GW = "1.6.5.1"

    @classmethod
    def setUpClass(cls):
        """
        Starts the thread which launches ryu apps

        Create a testing bridge, add a port, setup the port interfaces. Then
        launch the ryu apps for testing pipelined. Gets the references
        to apps launched by using futures.
        """
        super(UplinkBridgeWithNonNATTest_IP_VLAN_GW, cls).setUpClass()
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
                'sgi_management_iface_ip_addr': cls.SGi_IP,
                'sgi_management_iface_gw': cls.SGi_GW,
                'dev_vlan_in': "test_v_in",
                'dev_vlan_out': "test_v_out",
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
        assert_bridge_snapshot_match(self, self.UPLINK_BRIDGE,
                                     self.service_manager,
                                     include_stats=False)

        self.assertIn(cls.SGi_GW, get_iface_gw_ipv4(cls.UPLINK_BRIDGE),
                      "gw not found")

@unittest.skip
class UplinkBridgeWithNonNatUplinkConnect_Test(unittest.TestCase):
    BRIDGE = 'testing_br'
    IFACE = 'testing_br'
    MAC_DEST = "5e:cc:cc:b1:49:4b"
    BRIDGE_IP = '192.168.128.1'
    SCRIPT_PATH = "/home/vagrant/magma/lte/gateway/python/magma/mobilityd/"
    NET_SW_BR = "net_sw_up1"
    UPLINK_DHCP = "tino_dhcp"
    SCRIPT_PATH = "/home/vagrant/magma/lte/gateway/python/magma/mobilityd/"
    UPLINK_ETH_PORT = "upb_ul_0"
    UPLINK_BRIDGE = 'upt_br0'
    UPLINK_PATCH = 'test_patch_p2'
    ROUTER_IP = "10.55.0.211"


    @classmethod
    def _setup_vlan_network(cls, vlan: str):
        setup_vlan_switch = cls.SCRIPT_PATH + "scripts/setup-uplink-vlan-sw.sh"
        subprocess.check_call([setup_vlan_switch, cls.NET_SW_BR, "upb"])
        cls._setup_vlan(vlan)

    @classmethod
    def _setup_vlan(cls, vlan):
        setup_vlan_switch = cls.SCRIPT_PATH + "scripts/setup-uplink-vlan-srv.sh"
        subprocess.check_call([setup_vlan_switch, cls.NET_SW_BR, vlan, "55"])

    @classmethod
    def setUpClass(cls):
        """
        Starts the thread which launches ryu apps

        Create a testing bridge, add a port, setup the port interfaces. Then
        launch the ryu apps for testing pipelined. Gets the references
        to apps launched by using futures.
        """
        super(UplinkBridgeWithNonNatUplinkConnect_Test, cls).setUpClass()
        warnings.simplefilter('ignore')
        cls.service_manager = create_service_manager([])

        cls._setup_vlan_network("0")

        BridgeTools.create_bridge(cls.UPLINK_BRIDGE, cls.UPLINK_BRIDGE)
        BridgeTools.create_internal_iface(cls.UPLINK_BRIDGE,
                                          cls.UPLINK_DHCP, None)
        BridgeTools.create_internal_iface(cls.UPLINK_BRIDGE,
                                          cls.UPLINK_PATCH, None)

        check_connectivity(cls.ROUTER_IP, cls.UPLINK_ETH_PORT)

        BridgeTools.add_ovs_port(cls.UPLINK_BRIDGE, cls.UPLINK_ETH_PORT, "200")

        # this is setup after AGW boot up in NATed mode.
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
                'ovs_vlan_workaround': True,
                'dev_vlan_in': "testv1_in",
                'dev_vlan_out': "testv1_out",
            },
            mconfig=None,
            loop=None,
            service_manager=cls.service_manager,
            integ_test=False,
        )

        BridgeTools.create_bridge(cls.BRIDGE, cls.BRIDGE)

        cls.thread = start_ryu_app_thread(test_setup)
        cls.uplink_br_controller = uplink_bridge_controller_reference.result()

        cls.testing_controller = testing_controller_reference.result()


    @classmethod
    def tearDownClass(cls):
        stop_ryu_app_thread(cls.thread)
        BridgeTools.destroy_bridge(cls.BRIDGE)
        BridgeTools.destroy_bridge(cls.UPLINK_BRIDGE)
        BridgeTools.destroy_bridge(cls.NET_SW_BR)

    def testFlowSnapshotMatch(self):
        cls = self.__class__
        assert_bridge_snapshot_match(self, self.UPLINK_BRIDGE, self.service_manager,
                                     include_stats=False)
        self.assertEqual(get_ovsdb_port_tag(cls.UPLINK_BRIDGE), '[]')
        # after Non NAT init, router shld be accessible.
        # manually start DHCP client on up-br
        check_connectivity(cls.ROUTER_IP, cls.UPLINK_BRIDGE)


class UplinkBridgeTestNatIPAddr(unittest.TestCase):
    BRIDGE = 'testing_br'
    MAC_DEST = "5e:cc:cc:b1:49:4b"
    BRIDGE_IP = '192.168.128.1'
    BRIDGE_ETH_PORT = "eth_t1"
    UPLINK_BRIDGE = 'upt_br0'
    SGi_IP="1.6.5.77"

    @classmethod
    def setUpClass(cls):
        """
        Starts the thread which launches ryu apps

        Create a testing bridge, add a port, setup the port interfaces. Then
        launch the ryu apps for testing pipelined. Gets the references
        to apps launched by using futures.
        """
        super(UplinkBridgeTestNatIPAddr, cls).setUpClass()
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
                'uplink_bridge': cls.UPLINK_BRIDGE,
                'sgi_management_iface_ip_addr': cls.SGi_IP,
                'uplink_eth_port_name': cls.BRIDGE_ETH_PORT,
            },
            mconfig=None,
            loop=None,
            service_manager=cls.service_manager,
            integ_test=False,
        )

        BridgeTools.create_bridge(cls.BRIDGE, cls.BRIDGE)
        BridgeTools.create_bridge(cls.UPLINK_BRIDGE, cls.UPLINK_BRIDGE)
        BridgeTools.create_internal_iface(cls.BRIDGE,
                                          cls.BRIDGE_ETH_PORT, '2.2.2.2')
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

        assert_bridge_snapshot_match(self, self.UPLINK_BRIDGE, self.service_manager)
        self.assertIn(cls.SGi_IP, get_iface_ipv4(cls.BRIDGE_ETH_PORT), "ip not found")


if __name__ == "__main__":
    unittest.main()


def check_connectivity(dst: str, dev_name: str):
    try:
        ifdown_if = ["dhclient", dev_name]
        subprocess.check_call(ifdown_if)
    except subprocess.SubprocessError as e:
        logging.warning("Error while setting dhcl IP: %s: %s",
                        dev_name, e)
        return
    hub.sleep(1)

    try:
        ping_cmd = ["ping", "-c", "3", dst]
        subprocess.check_call(ping_cmd)
    except subprocess.SubprocessError as e:
        logging.warning("Error while ping: %s", e)
        # for now dont assert here.

    validate_routing_table(dst, dev_name)


def validate_routing_table(dst: str, dev_name: str) -> str:
    dump1 = subprocess.Popen(["ip", "r", "get", dst],
                             stdout=subprocess.PIPE)
    for line in dump1.stdout.readlines():
        if "dev" not in str(line):
            continue
        try:
            if dev_name in str(line):
                return
        except ValueError:
            pass
    logging.error("could not find route to %s via %s", dst, dev_name)
    dump1 = subprocess.Popen(["ovs-ofctl", "dump-flows", cls.UPLINK_BRIDGE],
                             stdout=subprocess.PIPE)
    for line in dump1.stdout.readlines():
        print("pbs: %s", line)
    assert 0


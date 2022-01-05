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
import ipaddress
import logging
import subprocess
import threading
import time
import unittest
import warnings
from concurrent.futures import Future
from os import pipe
from typing import List

from lte.protos.mobilityd_pb2 import GWInfo, IPAddress, IPBlock
from magma.pipelined.app import inout
from magma.pipelined.bridge_util import BridgeTools
from magma.pipelined.tests.app.start_pipelined import (
    PipelinedController,
    TestSetup,
)
from magma.pipelined.tests.pipelined_test_util import (
    SnapshotVerifier,
    create_service_manager,
    fake_inout_setup,
    start_ryu_app_thread,
    stop_ryu_app_thread,
)
from ryu.lib import hub
from ryu.ofproto.ofproto_v1_4 import OFPP_LOCAL

gw_info_map = {}
gw_info_lock = threading.RLock()  # re-entrant locks


def mocked_get_mobilityd_gw_info() -> List[GWInfo]:
    global gw_info_map
    global gw_info_lock

    with gw_info_lock:
        return gw_info_map.values()


def mocked_set_mobilityd_gw_info(ip: IPAddress, mac: str, vlan: str):
    global gw_info_map
    global gw_info_lock

    with gw_info_lock:
        gw_info = GWInfo(ip=ip, mac=mac, vlan=vlan)
        gw_info_map[vlan] = gw_info


def clear_gw_info_map():
    global gw_info_map
    global gw_info_lock

    with gw_info_lock:
        gw_info_map = {}


def check_GW_rec(vlan):
    global gw_info_map
    global gw_info_lock

    with gw_info_lock:
        while gw_info_map[vlan].mac is None or gw_info_map[vlan].mac == '':
            threading.Event().wait(.5)
            print("waiting for mac on vlan [%s]" % vlan)


def assert_GW_mac(tc, vlan, mac):
    global gw_info_map
    global gw_info_lock

    with gw_info_lock:
        tc.assertEqual(gw_info_map[vlan].mac, mac)


class InOutNonNatTest(unittest.TestCase):
    BRIDGE = 'testing_br'
    IFACE = 'testing_br'
    MAC_DEST = "5e:cc:cc:b1:49:4b"
    BRIDGE_IP = '192.168.128.1'
    SCRIPT_PATH = "/home/vagrant/magma/lte/gateway/python/magma/mobilityd/"
    NON_NAT_ARP_EGRESS_PORT = "tinouplink_p0"
    UPLINK_BR = "up_inout_br0"
    DHCP_PORT = "tino_dhcp"
    UPLINK_VLAN_SW = "vlan_inout"

    @classmethod
    def setup_uplink_br(cls):
        setup_dhcp_server = cls.SCRIPT_PATH + "scripts/setup-test-dhcp-srv.sh"
        subprocess.check_call([setup_dhcp_server, "tino"])

        BridgeTools.destroy_bridge(cls.UPLINK_BR)
        setup_uplink_br = [
            cls.SCRIPT_PATH + "scripts/setup-uplink-br.sh",
            cls.UPLINK_BR,
            cls.NON_NAT_ARP_EGRESS_PORT,
            cls.DHCP_PORT,
        ]
        subprocess.check_call(setup_uplink_br)

    @classmethod
    def _setup_vlan_network(cls, vlan: str):
        setup_vlan_switch = cls.SCRIPT_PATH + "scripts/setup-uplink-vlan-sw.sh"
        subprocess.check_call([setup_vlan_switch, cls.UPLINK_VLAN_SW, "inout"])

        setup_uplink_br = [
            cls.SCRIPT_PATH + "scripts/setup-uplink-br.sh",
            cls.UPLINK_BR,
            "inout_ul_0",
            cls.DHCP_PORT,
        ]

        subprocess.check_call(setup_uplink_br)
        cls._setup_vlan(vlan)

    @classmethod
    def _setup_vlan(cls, vlan):
        setup_vlan_switch = cls.SCRIPT_PATH + "scripts/setup-uplink-vlan-srv.sh"
        subprocess.check_call([setup_vlan_switch, cls.UPLINK_VLAN_SW, vlan])

    def setUpNetworkAndController(
        self, vlan: str = "",
        non_nat_arp_egress_port: str = None,
        gw_mac_addr="ff:ff:ff:ff:ff:ff",
    ):
        """
        Starts the thread which launches ryu apps

        Create a testing bridge, add a port, setup the port interfaces. Then
        launch the ryu apps for testing pipelined. Gets the references
        to apps launched by using futures.
        """

        cls = self.__class__
        super(InOutNonNatTest, cls).setUpClass()
        inout.get_mobilityd_gw_info = mocked_get_mobilityd_gw_info
        inout.set_mobilityd_gw_info = mocked_set_mobilityd_gw_info

        warnings.simplefilter('ignore')
        cls.setup_uplink_br()

        if vlan != "":
            cls._setup_vlan_network(vlan)

        cls.service_manager = create_service_manager([])

        inout_controller_reference = Future()
        testing_controller_reference = Future()

        if non_nat_arp_egress_port is None:
            non_nat_arp_egress_port = cls.DHCP_PORT

        patch_up_port_no = BridgeTools.get_ofport('patch-up')
        test_setup = TestSetup(
            apps=[
                PipelinedController.InOut,
                PipelinedController.Testing,
                PipelinedController.StartupFlows,
            ],
            references={
                PipelinedController.InOut:
                    inout_controller_reference,
                PipelinedController.Testing:
                    testing_controller_reference,
                PipelinedController.StartupFlows:
                    Future(),
            },
            config={
                'setup_type': 'LTE',
                'bridge_name': cls.BRIDGE,
                'bridge_ip_address': cls.BRIDGE_IP,
                'ovs_gtp_port_number': 32768,
                'clean_restart': True,
                'enable_nat': False,
                'non_nat_gw_probe_frequency': 0.5,
                'non_nat_arp_egress_port': non_nat_arp_egress_port,
                'uplink_port': patch_up_port_no,
                'uplink_gw_mac': gw_mac_addr,
            },
            mconfig=None,
            loop=None,
            service_manager=cls.service_manager,
            integ_test=False,
        )

        BridgeTools.create_bridge(cls.BRIDGE, cls.IFACE)
        subprocess.Popen(["ifconfig", cls.UPLINK_BR, "192.168.128.41"]).wait()
        cls.thread = start_ryu_app_thread(test_setup)
        cls.inout_controller = inout_controller_reference.result()

        cls.testing_controller = testing_controller_reference.result()

    def tearDown(self):
        cls = self.__class__
        cls.inout_controller._stop_gw_mac_monitor()
        stop_ryu_app_thread(cls.thread)
        BridgeTools.destroy_bridge(cls.BRIDGE)
        BridgeTools.destroy_bridge(cls.UPLINK_BR)
        BridgeTools.destroy_bridge(cls.UPLINK_VLAN_SW)

        time.sleep(1)
        clear_gw_info_map()

    def setUp(self):
        clear_gw_info_map()

    def testFlowSnapshotMatch(self):
        cls = self.__class__
        self.setUpNetworkAndController(
            non_nat_arp_egress_port=cls.UPLINK_BR,
            gw_mac_addr="33:44:55:ff:ff:ff",
        )
        fake_inout_setup(cls.inout_controller)
        # wait for atleast one iteration of the ARP probe.

        ip_addr = ipaddress.ip_address("192.168.128.211")
        vlan = ""
        mocked_set_mobilityd_gw_info(
            IPAddress(
                address=ip_addr.packed,
                version=IPBlock.IPV4,
            ),
            "b2:a0:cc:85:80:7a", vlan,
        )
        check_GW_rec(vlan)

        snapshot_verifier = SnapshotVerifier(
            self, self.BRIDGE,
            self.service_manager,
            max_sleep_time=40,
            datapath=cls.inout_controller._datapath,
        )
        with snapshot_verifier:
            pass
        assert_GW_mac(self, vlan, 'b2:a0:cc:85:80:7a')

    def testFlowVlanSnapshotMatch(self):
        cls = self.__class__
        vlan = "11"
        self.setUpNetworkAndController(vlan)
        # wait for atleast one iteration of the ARP probe.
        fake_inout_setup(cls.inout_controller)

        ip_addr = ipaddress.ip_address("10.200.11.211")
        mocked_set_mobilityd_gw_info(
            IPAddress(
                address=ip_addr.packed,
                version=IPBlock.IPV4,
            ),
            'b2:a0:cc:85:80:11', vlan,
        )
        hub.sleep(2)

        check_GW_rec(vlan)

        logging.info("done waiting for vlan: %s", vlan)
        snapshot_verifier = SnapshotVerifier(
            self, self.BRIDGE,
            self.service_manager,
            max_sleep_time=40,
            datapath=cls.inout_controller._datapath,
        )

        with snapshot_verifier:
            pass
        assert_GW_mac(self, vlan, 'b2:a0:cc:85:80:11')

    def testFlowVlanSnapshotMatch2(self):
        cls = self.__class__
        vlan1 = "21"
        self.setUpNetworkAndController(vlan1)
        vlan2 = "22"
        cls._setup_vlan(vlan2)
        fake_inout_setup(cls.inout_controller)

        ip_addr = ipaddress.ip_address("10.200.21.211")
        mocked_set_mobilityd_gw_info(
            IPAddress(
                address=ip_addr.packed,
                version=IPBlock.IPV4,
            ),
            "b2:a0:cc:85:80:21", vlan1,
        )

        ip_addr = ipaddress.ip_address("10.200.22.211")
        mocked_set_mobilityd_gw_info(
            IPAddress(
                address=ip_addr.packed,
                version=IPBlock.IPV4,
            ),
            "b2:a0:cc:85:80:22", vlan2,
        )
        check_GW_rec(vlan1)
        check_GW_rec(vlan2)

        logging.info("done waiting for vlan: %s", vlan1)
        snapshot_verifier = SnapshotVerifier(
            self, self.BRIDGE,
            self.service_manager,
            max_sleep_time=40,
            datapath=cls.inout_controller._datapath,
        )

        with snapshot_verifier:
            pass

    def testFlowVlanSnapshotMatch_static1(self):
        cls = self.__class__
        # setup network on unused vlan.
        vlan1 = "21"
        self.setUpNetworkAndController(vlan1)
        # statically configured config
        vlan2 = "22"
        fake_inout_setup(cls.inout_controller)

        ip_addr = ipaddress.ip_address("10.200.21.211")
        mocked_set_mobilityd_gw_info(
            IPAddress(
                address=ip_addr.packed,
                version=IPBlock.IPV4,
            ),
            "11:33:44:55:66:77", vlan1,
        )

        ip_addr = ipaddress.ip_address("10.200.22.211")
        mocked_set_mobilityd_gw_info(
            IPAddress(
                address=ip_addr.packed,
                version=IPBlock.IPV4,
            ),
            "22:33:44:55:66:77", vlan2,
        )

        hub.sleep(2)
        check_GW_rec(vlan2)

        logging.info("done waiting for vlan: %s", vlan1)
        snapshot_verifier = SnapshotVerifier(
            self, self.BRIDGE,
            self.service_manager,
            max_sleep_time=20,
            datapath=cls.inout_controller._datapath,
        )

        with snapshot_verifier:
            pass
        assert_GW_mac(self, vlan1, 'b2:a0:cc:85:80:21')
        assert_GW_mac(self, vlan2, '22:33:44:55:66:77')

    def testFlowVlanSnapshotMatch_static2(self):
        cls = self.__class__
        # setup network on unused vlan.
        self.setUpNetworkAndController("34")
        # statically configured config
        fake_inout_setup(cls.inout_controller)

        vlan1 = "31"
        ip_addr = ipaddress.ip_address("10.200.21.100")
        mocked_set_mobilityd_gw_info(
            IPAddress(
                address=ip_addr.packed,
                version=IPBlock.IPV4,
            ),
            "11:33:44:55:66:77", vlan1,
        )

        vlan2 = "32"
        ip_addr = ipaddress.ip_address("10.200.22.200")
        mocked_set_mobilityd_gw_info(
            IPAddress(
                address=ip_addr.packed,
                version=IPBlock.IPV4,
            ),
            "22:33:44:55:66:77", vlan2,
        )

        ip_addr = ipaddress.ip_address("10.200.22.10")
        mocked_set_mobilityd_gw_info(
            IPAddress(
                address=ip_addr.packed,
                version=IPBlock.IPV4,
            ),
            "00:33:44:55:66:77", "",
        )

        check_GW_rec(vlan1)
        check_GW_rec(vlan2)
        check_GW_rec("")

        assert_GW_mac(self, vlan1, '11:33:44:55:66:77')
        assert_GW_mac(self, vlan2, '22:33:44:55:66:77')
        assert_GW_mac(self, "", '00:33:44:55:66:77')

        hub.sleep(2)
        snapshot_verifier = SnapshotVerifier(
            self, self.BRIDGE,
            self.service_manager,
            max_sleep_time=20,
            datapath=cls.inout_controller._datapath,
            try_snapshot=True,
        )

        with snapshot_verifier:
            pass


class InOutTestNonNATBasicFlows(unittest.TestCase):
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
        super(InOutTestNonNATBasicFlows, cls).setUpClass()
        warnings.simplefilter('ignore')
        cls.service_manager = create_service_manager([])

        inout_controller_reference = Future()
        testing_controller_reference = Future()
        test_setup = TestSetup(
            apps=[
                PipelinedController.InOut,
                PipelinedController.Testing,
                PipelinedController.StartupFlows,
            ],
            references={
                PipelinedController.InOut:
                    inout_controller_reference,
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
                'uplink_gw_mac': '11:22:33:44:55:66',
                'uplink_port': OFPP_LOCAL,
            },
            mconfig=None,
            loop=None,
            service_manager=cls.service_manager,
            integ_test=False,
        )

        BridgeTools.create_bridge(cls.BRIDGE, cls.IFACE)

        cls.thread = start_ryu_app_thread(test_setup)
        cls.inout_controller = inout_controller_reference.result()
        cls.testing_controller = testing_controller_reference.result()

    @classmethod
    def tearDownClass(cls):
        cls.inout_controller._stop_gw_mac_monitor()
        stop_ryu_app_thread(cls.thread)
        BridgeTools.destroy_bridge(cls.BRIDGE)

        time.sleep(1)
        clear_gw_info_map()

    def testFlowSnapshotMatch(self):
        fake_inout_setup(self.inout_controller)
        snapshot_verifier = SnapshotVerifier(
            self,
            self.BRIDGE,
            self.service_manager,
            max_sleep_time=20,
            datapath=InOutTestNonNATBasicFlows.inout_controller._datapath,
        )

        with snapshot_verifier:
            pass


ipv6_mac_table = {}


def mocked_setmacbyip6(ipv6_addr: str, mac: str):
    global ipv6_mac_table
    with gw_info_lock:
        ipv6_mac_table[ipv6_addr] = mac


def mocked_getmacbyip6(ipv6_addr: str) -> str:
    global ipv6_mac_table
    with gw_info_lock:
        return ipv6_mac_table.get(ipv6_addr, None)


class InOutTestNonNATBasicFlowsIPv6(unittest.TestCase):
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
        inout.getmacbyip6 = mocked_getmacbyip6
        inout.get_mobilityd_gw_info = mocked_get_mobilityd_gw_info
        inout.set_mobilityd_gw_info = mocked_set_mobilityd_gw_info

        super(InOutTestNonNATBasicFlowsIPv6, cls).setUpClass()
        warnings.simplefilter('ignore')
        cls.service_manager = create_service_manager([])

        inout_controller_reference = Future()
        testing_controller_reference = Future()
        test_setup = TestSetup(
            apps=[
                PipelinedController.InOut,
                PipelinedController.Testing,
                PipelinedController.StartupFlows,
            ],
            references={
                PipelinedController.InOut:
                    inout_controller_reference,
                PipelinedController.Testing:
                    testing_controller_reference,
                PipelinedController.StartupFlows:
                    Future(),
            },
            config={
                'setup_type': 'LTE',
                'bridge_name': cls.BRIDGE,
                'bridge_ip_address': cls.BRIDGE_IP,
                'ovs_gtp_port_number': 32768,
                'clean_restart': True,
                'enable_nat': False,
                'uplink_gw_mac': '11:22:33:44:55:11',
                'uplink_port': OFPP_LOCAL,
                'non_nat_gw_probe_frequency': 0.5,
            },
            mconfig=None,
            loop=None,
            service_manager=cls.service_manager,
            integ_test=False,
        )

        BridgeTools.create_bridge(cls.BRIDGE, cls.IFACE)

        cls.thread = start_ryu_app_thread(test_setup)
        cls.inout_controller = inout_controller_reference.result()
        cls.testing_controller = testing_controller_reference.result()

    @classmethod
    def tearDownClass(cls):
        cls.inout_controller._stop_gw_mac_monitor()
        stop_ryu_app_thread(cls.thread)
        BridgeTools.destroy_bridge(cls.BRIDGE)
        time.sleep(1)
        clear_gw_info_map()

    def testFlowSnapshotMatch(self):
        ipv6_addr1 = "2002::22"
        mac_addr1 = "11:22:33:44:55:88"
        mocked_setmacbyip6(ipv6_addr1, mac_addr1)

        ip_addr = ipaddress.ip_address(ipv6_addr1)
        vlan = "NO_VLAN"
        mocked_set_mobilityd_gw_info(
            IPAddress(
                address=ip_addr.packed,
                version=IPBlock.IPV6,
            ),
            "", vlan,
        )
        fake_inout_setup(self.inout_controller)

        snapshot_verifier = SnapshotVerifier(
            self,
            self.BRIDGE,
            self.service_manager,
            max_sleep_time=20,
            datapath=InOutTestNonNATBasicFlowsIPv6.inout_controller._datapath,
        )

        with snapshot_verifier:
            pass
        assert_GW_mac(self, vlan, mac_addr1)


if __name__ == "__main__":
    unittest.main()

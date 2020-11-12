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
import subprocess
import time
import threading
import unittest
import warnings
from concurrent.futures import Future
import logging
from typing import List
from ryu.lib import hub

from lte.protos.mobilityd_pb2 import IPAddress, GWInfo, IPBlock

from magma.pipelined.tests.app.start_pipelined import (
    TestSetup,
    PipelinedController,
)
from magma.pipelined.bridge_util import BridgeTools
from magma.pipelined.tests.pipelined_test_util import (
    start_ryu_app_thread,
    stop_ryu_app_thread,
    create_service_manager,
    SnapshotVerifier
)
from ryu.ofproto.ofproto_v1_4 import OFPP_LOCAL

from magma.pipelined.app import inout

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


@unittest.skip("needs more investigation.")
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
        setup_uplink_br = [cls.SCRIPT_PATH + "scripts/setup-uplink-br.sh",
                           cls.UPLINK_BR,
                           cls.NON_NAT_ARP_EGRESS_PORT,
                           cls.DHCP_PORT]
        subprocess.check_call(setup_uplink_br)

    @classmethod
    def _setup_vlan_network(cls, vlan: str):
        setup_vlan_switch = cls.SCRIPT_PATH + "scripts/setup-uplink-vlan-sw.sh"
        subprocess.check_call([setup_vlan_switch, cls.UPLINK_VLAN_SW, "inout"])

        setup_uplink_br = [cls.SCRIPT_PATH + "scripts/setup-uplink-br.sh",
                           cls.UPLINK_BR,
                           "inout_ul_0",
                           cls.DHCP_PORT]

        subprocess.check_call(setup_uplink_br)
        cls._setup_vlan(vlan)

    @classmethod
    def _setup_vlan(cls, vlan):
        setup_vlan_switch = cls.SCRIPT_PATH + "scripts/setup-uplink-vlan-srv.sh"
        subprocess.check_call([setup_vlan_switch, cls.UPLINK_VLAN_SW, vlan])

    def setUpNetworkAndController(self, vlan: str = "",
                                  non_nat_arp_egress_port: str = None,
                                  gw_mac_addr="ff:ff:ff:ff:ff:ff"):
        """
        Starts the thread which launches ryu apps

        Create a testing bridge, add a port, setup the port interfaces. Then
        launch the ryu apps for testing pipelined. Gets the references
        to apps launched by using futures.
        """
        global gw_info_map
        gw_info_map.clear()
        hub.sleep(2)

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
            apps=[PipelinedController.InOut,
                  PipelinedController.Testing,
                  PipelinedController.StartupFlows],
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
            integ_test=False
        )

        BridgeTools.create_bridge(cls.BRIDGE, cls.IFACE)
        subprocess.Popen(["ifconfig", cls.UPLINK_BR, "192.168.128.41"]).wait()
        cls.thread = start_ryu_app_thread(test_setup)
        cls.inout_controller = inout_controller_reference.result()

        cls.testing_controller = testing_controller_reference.result()

    def tearDown(self):
        cls = self.__class__
        stop_ryu_app_thread(cls.thread)
        subprocess.Popen(["ovs-ofctl", "del-flows", cls.BRIDGE]).wait()
        subprocess.Popen(["ovs-vsctl", "del-br", cls.UPLINK_BR]).wait()

        hub.sleep(1)

    def testFlowSnapshotMatch(self):
        cls = self.__class__
        self.setUpNetworkAndController(non_nat_arp_egress_port=cls.UPLINK_BR,
                                       gw_mac_addr="33:44:55:ff:ff:ff")
        # wait for atleast one iteration of the ARP probe.

        ip_addr = ipaddress.ip_address("192.168.128.211")
        vlan = ""
        mocked_set_mobilityd_gw_info(IPAddress(address=ip_addr.packed,
                                               version=IPBlock.IPV4),
                                     "", vlan)
        hub.sleep(2)
        while gw_info_map[vlan].mac is None or gw_info_map[vlan].mac == '':
            threading.Event().wait(.5)

        snapshot_verifier = SnapshotVerifier(self, self.BRIDGE,
                                             self.service_manager,
                                             max_sleep_time=40,
                                             datapath=cls.inout_controller._datapath,
                                             try_snapshot=True)
        with snapshot_verifier:
            pass
        self.assertEqual(gw_info_map[vlan].mac, 'b2:a0:cc:85:80:7a')

    def testFlowVlanSnapshotMatch(self):
        cls = self.__class__
        vlan = "11"
        self.setUpNetworkAndController(vlan)
        # wait for atleast one iteration of the ARP probe.

        ip_addr = ipaddress.ip_address("10.200.11.211")
        mocked_set_mobilityd_gw_info(IPAddress(address=ip_addr.packed,
                                               version=IPBlock.IPV4),
                                     "", vlan)
        hub.sleep(2)
        while gw_info_map[vlan].mac is None or gw_info_map[vlan].mac == '':
            threading.Event().wait(.5)

        logging.info("done waiting for vlan: %s", vlan)
        snapshot_verifier = SnapshotVerifier(self, self.BRIDGE,
                                             self.service_manager,
                                             max_sleep_time=40,
                                             datapath=cls.inout_controller._datapath,
                                             try_snapshot=True)

        with snapshot_verifier:
            pass
        self.assertEqual(gw_info_map[vlan].mac, 'b2:a0:cc:85:80:11')

    def testFlowVlanSnapshotMatch2(self):
        cls = self.__class__
        vlan1 = "21"
        self.setUpNetworkAndController(vlan1)
        vlan2 = "22"
        cls._setup_vlan(vlan2)

        ip_addr = ipaddress.ip_address("10.200.21.211")
        mocked_set_mobilityd_gw_info(IPAddress(address=ip_addr.packed,
                                               version=IPBlock.IPV4),
                                     "", vlan1)

        ip_addr = ipaddress.ip_address("10.200.22.211")
        mocked_set_mobilityd_gw_info(IPAddress(address=ip_addr.packed,
                                               version=IPBlock.IPV4),
                                     "", vlan2)

        hub.sleep(2)
        while gw_info_map[vlan2].mac is None or gw_info_map[vlan2].mac == '':
            threading.Event().wait(.5)

        logging.info("done waiting for vlan: %s", vlan1)
        snapshot_verifier = SnapshotVerifier(self, self.BRIDGE,
                                             self.service_manager,
                                             max_sleep_time=40,
                                             datapath=cls.inout_controller._datapath,
                                             try_snapshot=True)

        with snapshot_verifier:
            pass
        self.assertEqual(gw_info_map[vlan1].mac, 'b2:a0:cc:85:80:21')
        self.assertEqual(gw_info_map[vlan2].mac, 'b2:a0:cc:85:80:22')

    def testFlowVlanSnapshotMatch_static1(self):
        cls = self.__class__
        # setup network on unused vlan.
        vlan1 = "21"
        self.setUpNetworkAndController(vlan1)
        # statically configured config
        vlan2 = "22"

        ip_addr = ipaddress.ip_address("10.200.21.211")
        mocked_set_mobilityd_gw_info(IPAddress(address=ip_addr.packed,
                                               version=IPBlock.IPV4),
                                     "11:33:44:55:66:77", vlan1)

        ip_addr = ipaddress.ip_address("10.200.22.211")
        mocked_set_mobilityd_gw_info(IPAddress(address=ip_addr.packed,
                                               version=IPBlock.IPV4),
                                     "22:33:44:55:66:77", vlan2)

        hub.sleep(2)
        while gw_info_map[vlan2].mac is None or gw_info_map[vlan2].mac == '':
            threading.Event().wait(.5)

        logging.info("done waiting for vlan: %s", vlan1)
        snapshot_verifier = SnapshotVerifier(self, self.BRIDGE,
                                             self.service_manager,
                                             max_sleep_time=20,
                                             datapath=cls.inout_controller._datapath,
                                             try_snapshot=True)

        with snapshot_verifier:
            pass
        self.assertEqual(gw_info_map[vlan1].mac, 'b2:a0:cc:85:80:21')
        self.assertEqual(gw_info_map[vlan2].mac, '22:33:44:55:66:77')

    def testFlowVlanSnapshotMatch_static2(self):
        cls = self.__class__
        # setup network on unused vlan.
        self.setUpNetworkAndController("34")
        # statically configured config

        vlan1 = "31"
        ip_addr = ipaddress.ip_address("10.200.21.100")
        mocked_set_mobilityd_gw_info(IPAddress(address=ip_addr.packed,
                                               version=IPBlock.IPV4),
                                     "11:33:44:55:66:77", vlan1)

        vlan2 = "32"
        ip_addr = ipaddress.ip_address("10.200.22.200")
        mocked_set_mobilityd_gw_info(IPAddress(address=ip_addr.packed,
                                               version=IPBlock.IPV4),
                                     "22:33:44:55:66:77", vlan2)

        ip_addr = ipaddress.ip_address("10.200.22.10")
        mocked_set_mobilityd_gw_info(IPAddress(address=ip_addr.packed,
                                               version=IPBlock.IPV4),
                                     "00:33:44:55:66:77", "")

        hub.sleep(2)
        while gw_info_map[vlan2].mac is None or gw_info_map[vlan2].mac == '':
            threading.Event().wait(.5)

        logging.info("done waiting for vlan: %s", vlan1)
        snapshot_verifier = SnapshotVerifier(self, self.BRIDGE,
                                             self.service_manager,
                                             max_sleep_time=20,
                                             datapath=cls.inout_controller._datapath)

        with snapshot_verifier:
            pass
        self.assertEqual(gw_info_map[vlan1].mac, '11:33:44:55:66:77')
        self.assertEqual(gw_info_map[vlan2].mac, '22:33:44:55:66:77')
        self.assertEqual(gw_info_map[""].mac, '00:33:44:55:66:77')


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
            apps=[PipelinedController.InOut,
                  PipelinedController.Testing,
                  PipelinedController.StartupFlows],
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
                'uplink_port': OFPP_LOCAL
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
        stop_ryu_app_thread(cls.thread)
        BridgeTools.destroy_bridge(cls.BRIDGE)

    def testFlowSnapshotMatch(self):
        snapshot_verifier = SnapshotVerifier(self,
                                             self.BRIDGE,
                                             self.service_manager,
                                             max_sleep_time=20,
                                             datapath=InOutTestNonNATBasicFlows.inout_controller._datapath)

        with snapshot_verifier:
            pass


if __name__ == "__main__":
    unittest.main()

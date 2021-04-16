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
import datetime
import ipaddress
import logging
import os
import subprocess
import sys
import threading
import time
import unittest

from freezegun import freeze_time
from magma.mobilityd.dhcp_client import DHCPClient
from magma.mobilityd.dhcp_desc import DHCPDescriptor, DHCPState
from magma.mobilityd.mac import MacAddress
from magma.mobilityd.uplink_gw import UplinkGatewayInfo
from magma.pipelined.bridge_util import BridgeTools
from scapy.layers.dhcp import DHCP
from scapy.layers.l2 import Dot1Q, Ether
from scapy.sendrecv import AsyncSniffer

LOG = logging.getLogger('mobilityd.dhcp.test')
LOG.isEnabledFor(logging.DEBUG)

logging.basicConfig(stream=sys.stderr, level=logging.DEBUG)
SCRIPT_PATH = "/home/vagrant/magma/lte/gateway/python/magma/mobilityd/"
DHCP_IFACE = "cl1_dhcp0"
PKT_CAPTURE_WAIT = 2
RETRY_LIMIT = 10

"""
Test dhclient class independent of IP allocator.
"""


class DhcpClient(unittest.TestCase):
    def setUp(self):
        self._br = "dh_br0"
        self.vlan_sw = "vlan_sw"
        self.up_link_port = ""

        try:
            subprocess.check_call(["pkill", "dnsmasq"])
        except subprocess.CalledProcessError:
            pass

        self.dhcp_wait = threading.Condition()
        self.pkt_list_lock = threading.Condition()
        self.dhcp_store = {}
        self.gw_info_map = {}
        self.gw_info = UplinkGatewayInfo(self.gw_info_map)
        self._sniffer = None
        self._last_xid = -1

    def tearDown(self):
        self._dhcp_client.stop()
        BridgeTools.destroy_bridge(self._br)

    @unittest.skipIf(os.getuid(), reason="needs root user")
    def test_dhcp_lease1(self):
        self._setup_dhcp_vlan_off()

        mac1 = MacAddress("11:22:33:44:55:66")
        self._validate_dhcp_alloc_renew(mac1, None)

    """
    Network diagram of vlan test setup                                    VETH pair        +---------------------+
    ==================================                                                     |  Namespace 1 with   |
                                                                      +--------------------+  DHCP server        |
                                                                      |                    |                     |
                                                                      |                    +---------------------+
                                                                      |
    +----------------+  Patch   +--------------+  uplink-    +--------+-------+            +---------------------+
    |                |  Port    |              |  iface      |                | VETH Pair  |  Namespace 2 with   |
    |  GTP_BR0       +----------+  UPLINK_BR0  +-------------+    VLAN_SW     +------------+  DHCP server        |
    |                |          |              |             |                |            |                     |
    +----------------+          +--------------+             +----------------+            +---------------------+
    """

    @unittest.skipIf(os.getuid(), reason="needs root user")
    def test_dhcp_vlan(self):
        vlan1 = "2"
        mac1 = MacAddress("11:22:33:44:55:66")

        self._setup_vlan_network()
        self._setup_dhcp_on_vlan(vlan1)

        self._validate_dhcp_alloc_renew(mac1, vlan1)
        self._validate_ip_subnet(mac1, vlan1)
        self._release_ip(mac1, vlan1)

    @unittest.skip("needs more investigation.")
    def test_dhcp_vlan_multi(self):
        self._setup_vlan_network()

        vlan1 = "51"
        mac1 = MacAddress("11:22:33:44:55:66")
        vlan2 = "52"
        mac2 = MacAddress("22:22:33:44:55:66")
        vlan3 = "53"
        mac3 = MacAddress("11:22:33:44:55:66")

        self._setup_dhcp_on_vlan(vlan1)
        self._setup_dhcp_on_vlan(vlan2)
        self._setup_dhcp_on_vlan(vlan3)

        self._validate_dhcp_alloc_renew(mac1, vlan1)
        self._validate_dhcp_alloc_renew(mac2, vlan2)
        self._validate_dhcp_alloc_renew(mac3, vlan3)

        self._validate_ip_subnet(mac1, vlan1)
        self._validate_ip_subnet(mac2, vlan2)
        self._validate_ip_subnet(mac3, vlan3)

        self._release_ip(mac1, vlan1)
        self._release_ip(mac2, vlan2)
        self._release_ip(mac3, vlan3)

    def _validate_dhcp_alloc_renew(self, mac1: MacAddress, vlan: str):
        self._alloc_ip_address_from_dhcp(mac1, vlan)
        self._validate_req_state(mac1, DHCPState.REQUEST, vlan)
        self._validate_state_as_current(mac1, vlan)

        # trigger lease reneval before deadline
        self._last_xid = self._get_state_xid(mac1, vlan)
        LOG.debug("time: %s", datetime.datetime.now())
        time1 = datetime.datetime.now() + datetime.timedelta(seconds=100)
        self._start_sniffer()
        with freeze_time(time1):
            LOG.debug("check req packets time: %s", datetime.datetime.now())
            self._stop_sniffer_and_check(DHCPState.REQUEST, mac1, vlan)
            self._validate_req_state(mac1, DHCPState.REQUEST, vlan)
            self._validate_state_as_current(mac1, vlan)

            # trigger lease after deadline
            self._last_xid = self._get_state_xid(mac1, vlan)
            time2 = datetime.datetime.now() + datetime.timedelta(seconds=2000)
            self._start_sniffer()
            LOG.debug(
                "check discover packets time: %s",
                datetime.datetime.now(),
            )
            with freeze_time(time2):
                LOG.debug("check discover after lease loss")
                self._stop_sniffer_and_check(DHCPState.DISCOVER, mac1, vlan)
                self._validate_req_state(mac1, DHCPState.REQUEST, vlan)
                self._validate_state_as_current(mac1, vlan)

    def _setup_dhcp_vlan_off(self):
        self.up_link_port = "cl1uplink_p0"
        setup_dhcp_server = SCRIPT_PATH + "scripts/setup-test-dhcp-srv.sh"
        subprocess.check_call([setup_dhcp_server, "cl1"])

        setup_uplink_br = [
            SCRIPT_PATH + "scripts/setup-uplink-br.sh",
            self._br,
            self.up_link_port,
            DHCP_IFACE,
        ]
        subprocess.check_call(setup_uplink_br)
        self._setup_dhclp_client()

    def _setup_vlan_network(self):
        self.up_link_port = "v_ul_0"
        setup_vlan_switch = SCRIPT_PATH + "scripts/setup-uplink-vlan-sw.sh"
        subprocess.check_call([setup_vlan_switch, self.vlan_sw, "v"])

        setup_uplink_br = [
            SCRIPT_PATH + "scripts/setup-uplink-br.sh",
            self._br,
            self.up_link_port,
            DHCP_IFACE,
        ]
        subprocess.check_call(setup_uplink_br)
        self._setup_dhclp_client()

    def _setup_dhclp_client(self):
        self._dhcp_client = DHCPClient(
            dhcp_wait=self.dhcp_wait,
            dhcp_store=self.dhcp_store,
            gw_info=self.gw_info,
            iface=DHCP_IFACE,
            lease_renew_wait_min=1,
        )
        self._dhcp_client.run()

    def _setup_dhcp_on_vlan(self, vlan: str):
        setup_vlan_switch = SCRIPT_PATH + "scripts/setup-uplink-vlan-srv.sh"
        subprocess.check_call([setup_vlan_switch, self.vlan_sw, vlan])

    def _validate_req_state(
        self, mac: MacAddress,
        state: DHCPState, vlan: str,
    ):
        for x in range(RETRY_LIMIT):
            LOG.debug("wait for state: %d" % x)
            with self.dhcp_wait:
                dhcp1 = self.dhcp_store.get(mac.as_redis_key(vlan))
                if state == DHCPState.RELEASE and dhcp1 is None:
                    return
                if dhcp1.state_requested == state:
                    return
            time.sleep(PKT_CAPTURE_WAIT)

        assert 0

    def _get_state_xid(self, mac: MacAddress, vlan: str):
        with self.dhcp_wait:
            dhcp1 = self.dhcp_store.get(mac.as_redis_key(vlan))
            return dhcp1.xid

    def _validate_state_as_current(self, mac: MacAddress, vlan: str):
        with self.dhcp_wait:
            dhcp1 = self.dhcp_store.get(mac.as_redis_key(vlan))
            self.assert_(
                dhcp1.state == DHCPState.OFFER or dhcp1.state == DHCPState.ACK,
            )

    def _alloc_ip_address_from_dhcp(
            self, mac: MacAddress, vlan: str,
    ) -> DHCPDescriptor:
        retry_count = 0
        with self.dhcp_wait:
            dhcp_desc = None
            while (
                retry_count < 60 and (
                    dhcp_desc is None
                    or dhcp_desc.ip_is_allocated() is not True
                )
            ):
                if retry_count % 5 == 0:
                    self._dhcp_client.send_dhcp_packet(
                        mac, vlan,
                        DHCPState.DISCOVER,
                    )

                self.dhcp_wait.wait(timeout=1)
                dhcp_desc = self._dhcp_client.get_dhcp_desc(mac, vlan)
                retry_count = retry_count + 1

            return dhcp_desc

    def _validate_ip_subnet(self, mac: MacAddress, vlan: str):
        # vlan is configured with subnet : 10.200.x.1
        # router IP is 10.200.x.211
        # x is vlan id
        exptected_subnet = ipaddress.ip_network("10.200.%s.0/24" % vlan)
        exptected_router_ip = ipaddress.ip_address("10.200.%s.211" % vlan)
        with self.dhcp_wait:
            dhcp1 = self.dhcp_store.get(mac.as_redis_key(vlan))
            self.assertEqual(dhcp1.subnet, str(exptected_subnet))
            self.assertEqual(dhcp1.router_ip, exptected_router_ip)
            self.assert_(ipaddress.ip_address(dhcp1.ip) in exptected_subnet)

    def _release_ip(self, mac: MacAddress, vlan: str):
        self._dhcp_client.release_ip_address(mac, vlan)
        time.sleep(PKT_CAPTURE_WAIT)
        self._validate_req_state(mac, DHCPState.RELEASE, vlan)

    def _handle_dhcp_req_packet(self, packet):
        if DHCP not in packet:
            return
        with self.pkt_list_lock:
            self.pkt_list.append(packet)

    def _start_sniffer(self):
        # drop dhclient requests, this would avoid lease
        # renewal after freezing time.
        subprocess.check_call([
            "ovs-ofctl", "add-flow", self._br,
            "priority=100,in_port=" + self.up_link_port + ",action=drop",
        ])
        time.sleep(PKT_CAPTURE_WAIT)

        self._sniffer = AsyncSniffer(
            iface=self._br,
            filter="udp and (port 67 or 68)",
            store=False,
            prn=self._handle_dhcp_req_packet,
        )

        self.pkt_list = []
        self._sniffer.start()
        LOG.debug("sniffer started")
        time.sleep(PKT_CAPTURE_WAIT)

    def _stop_sniffer_and_check(self, state: DHCPState, mac: MacAddress, vlan):
        LOG.debug("delete drop flow")
        subprocess.check_call([
            "ovs-ofctl", "del-flows", self._br,
            "in_port=" + self.up_link_port,
        ])

        for x in range(RETRY_LIMIT):
            LOG.debug("wait for pkt: %d" % x)
            time.sleep(PKT_CAPTURE_WAIT)

            with self.pkt_list_lock:
                for pkt in self.pkt_list:
                    if DHCP in pkt:
                        if vlan and (
                            Dot1Q not in pkt
                            or vlan != str(pkt[Dot1Q].vlan)
                        ):
                            continue

                        if pkt[DHCP].options[0][1] == int(state) and \
                                pkt[Ether].src == str(mac):
                            self._sniffer.stop()
                            return

        LOG.debug("Failed check for dhcp packet: %s", state)
        with self.pkt_list_lock:
            for pkt in self.pkt_list:
                LOG.debug("DHCP pkt %s", pkt.show(dump=True))

        # validate if any dhcp packet was sent.
        if state == DHCPState.DISCOVER:
            self.assertNotEqual(
                self._last_xid,
                self._get_state_xid(mac, vlan),
            )

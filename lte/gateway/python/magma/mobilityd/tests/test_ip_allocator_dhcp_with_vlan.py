"""
Copyright 2020 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

Allocates IP address as per DHCP server in the uplink network.
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

import fakeredis
from freezegun import freeze_time
from magma.mobilityd.dhcp_desc import DHCPDescriptor, DHCPState
from magma.mobilityd.ip_allocator_dhcp import IPAllocatorDHCP
from magma.mobilityd.ip_descriptor import IPState
from magma.mobilityd.mac import MacAddress, create_mac_from_sid
from magma.mobilityd.mobility_store import MobilityStore
from magma.pipelined.bridge_util import BridgeTools

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


class IpAllocatorDhcp(unittest.TestCase):
    def setUp(self) -> None:
        self._br = "dh_br0"
        self.vlan_sw = "vlan_sw"
        self.up_link_port = ""

        try:
            subprocess.check_call(["pkill", "dnsmasq"])
        except subprocess.CalledProcessError:
            pass

        self.dhcp_wait = threading.Condition()
        self.pkt_list_lock = threading.Condition()

    def tearDown(self) -> None:
        BridgeTools.destroy_bridge(self._br)

    @unittest.skipIf(os.getuid(), reason="needs root user")
    def test_dhcp_lease1(self) -> None:
        self._setup_dhcp_vlan_off()
        sid1 = "IMSI001010000000001"
        vlan = 0
        self._validate_dhcp_alloc_renew(sid1, vlan)

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
    def test_dhcp_vlan(self) -> None:
        vlan1: int = 2
        sid1 = "IMSI001010000000001"

        self._setup_vlan_network()
        self._setup_dhcp_on_vlan(vlan1)

        dhcp_desc = self._validate_dhcp_alloc_renew(sid1, vlan1)
        self._validate_ip_subnet(sid1, vlan1)
        self._release_ip(sid1, vlan1, dhcp_desc)

    @unittest.skip("needs more investigation.")
    def test_dhcp_vlan_multi(self) -> None:
        self._setup_vlan_network()

        vlan1 = 51
        sid1 = "IMSI001010000000001"
        vlan2 = 52
        sid2 = "IMSI001010000000002"
        vlan3 = 53
        sid3 = "IMSI001010000000001"

        self._setup_dhcp_on_vlan(vlan1)
        self._setup_dhcp_on_vlan(vlan2)
        self._setup_dhcp_on_vlan(vlan3)

        dhcp_desc1 = self._validate_dhcp_alloc_renew(sid1, vlan1)
        dhcp_desc2 = self._validate_dhcp_alloc_renew(sid2, vlan2)
        dhcp_desc3 = self._validate_dhcp_alloc_renew(sid3, vlan3)

        self._validate_ip_subnet(sid1, vlan1)
        self._validate_ip_subnet(sid2, vlan2)
        self._validate_ip_subnet(sid3, vlan3)

        self._release_ip(sid1, vlan1, dhcp_desc1)
        self._release_ip(sid2, vlan2, dhcp_desc2)
        self._release_ip(sid3, vlan3, dhcp_desc3)

    def _validate_dhcp_alloc_renew(self, sid: str, vlan: int) -> DHCPDescriptor:
        dhcp_desc = self._alloc_ip_address_from_dhcp(sid, vlan)
        mac = create_mac_from_sid(sid)
        self._validate_req_state(mac, DHCPState.ACK, vlan)
        self._validate_state_as_current(mac, vlan)

        # trigger lease renewal after deadline but before expiry
        LOG.debug("time: %s", datetime.datetime.now())
        time_after_renew_deadline = datetime.datetime.now() + datetime.timedelta(seconds=100)
        with freeze_time(time_after_renew_deadline):
            LOG.debug("check req packets time: %s", datetime.datetime.now())
            self._ip_allocator.stop_monitor_thread(join=True, reset=True)
            self._validate_req_state(mac, DHCPState.ACK, vlan)
            self._validate_state_as_current(mac, vlan)

            # trigger lease after expiry
            time_after_lease_expiry = datetime.datetime.now() + datetime.timedelta(seconds=2000)
            self._ip_allocator.start_monitor_thread()
            LOG.debug(
                "check discover packets time: %s",
                datetime.datetime.now(),
            )
            with freeze_time(time_after_lease_expiry):
                LOG.debug("check discover after lease loss")
                self._ip_allocator.stop_monitor_thread(join=True, reset=True)
                self._validate_req_state(mac, DHCPState.ACK, vlan)
                self._validate_state_as_current(mac, vlan)

        return dhcp_desc

    def _setup_dhcp_vlan_off(self) -> None:
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
        self._setup_ip_allocator_dhcp()

    def _setup_vlan_network(self) -> None:
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
        self._setup_ip_allocator_dhcp()

    def _setup_ip_allocator_dhcp(self) -> None:
        self._ip_allocator = IPAllocatorDHCP(
            store=MobilityStore(fakeredis.FakeStrictRedis()),
            iface=DHCP_IFACE,
            lease_renew_wait_min=1,
        )

    def _setup_dhcp_on_vlan(self, vlan: int) -> None:
        setup_vlan_switch = SCRIPT_PATH + "scripts/setup-uplink-vlan-srv.sh"
        subprocess.check_call([setup_vlan_switch, self.vlan_sw, str(vlan)])

    def _validate_req_state(
        self, mac: MacAddress,
        state: DHCPState, vlan: int,
    ) -> None:
        for x in range(RETRY_LIMIT):
            LOG.debug("wait for state: %d" % x)
            with self.dhcp_wait:
                dhcp_desc = self._ip_allocator.get_dhcp_desc_from_store(mac, vlan)
                if state == DHCPState.RELEASE and dhcp_desc is None:
                    return
                if dhcp_desc.state_requested == state:
                    return
            time.sleep(PKT_CAPTURE_WAIT)

        assert 0

    def _validate_state_as_current(self, mac: MacAddress, vlan: int) -> None:
        with self.dhcp_wait:
            dhcp_desc = self._ip_allocator._store.dhcp_store.get(mac.as_redis_key(vlan))
            self.assertTrue(
                dhcp_desc.state == DHCPState.OFFER or dhcp_desc.state == DHCPState.ACK,
            )

    def _alloc_ip_address_from_dhcp(self, sid: str, vlan: int) -> DHCPDescriptor:
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
                    self._ip_allocator.alloc_ip_address(sid, vlan)

                self.dhcp_wait.wait(timeout=1)
                mac = create_mac_from_sid(sid)
                dhcp_desc = self._ip_allocator.get_dhcp_desc_from_store(mac, vlan)
                retry_count = retry_count + 1

            return dhcp_desc

    def _validate_ip_subnet(self, sid: str, vlan: int) -> None:
        # vlan is configured with subnet : 10.200.x.1
        # router IP is 10.200.x.211
        exptected_subnet = ipaddress.ip_network(f"10.200.{vlan}.0/24")
        exptected_router_ip = ipaddress.ip_address(f"10.200.{vlan}.211")
        mac = create_mac_from_sid(sid)
        with self.dhcp_wait:
            dhcp1 = self._ip_allocator._store.dhcp_store.get(mac.as_redis_key(vlan))
            self.assertEqual(dhcp1.subnet, str(exptected_subnet))
            self.assertEqual(dhcp1.router_ip, exptected_router_ip)
            self.assertTrue(ipaddress.ip_address(dhcp1.ip) in exptected_subnet)

    def _release_ip(self, sid: str, vlan: int, dhcp_desc: DHCPDescriptor) -> None:
        mac = create_mac_from_sid(sid)
        key = dhcp_desc.ip
        ip_desc = self._ip_allocator._store.ip_state_map.ip_states[IPState.ALLOCATED][key]
        self._ip_allocator.release_ip(ip_desc)
        time.sleep(PKT_CAPTURE_WAIT)
        self._validate_req_state(mac, DHCPState.RELEASE, vlan)

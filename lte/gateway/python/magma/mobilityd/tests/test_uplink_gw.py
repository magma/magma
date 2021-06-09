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
import unittest
from collections import defaultdict

from lte.protos.mobilityd_pb2 import GWInfo, IPAddress
from magma.mobilityd.uplink_gw import NO_VLAN, UplinkGatewayInfo

LOG = logging.getLogger('mobilityd.def_gw.test')
LOG.isEnabledFor(logging.DEBUG)


def _get_gw_info(ip, mac="", vlan=NO_VLAN):
    ip_addr = ipaddress.ip_address(ip)
    gw_ip = IPAddress(
        version=IPAddress.IPV4,
        address=ip_addr.packed,
    )
    return GWInfo(ip=gw_ip, mac=mac, vlan=vlan)


def gw_list_to_set(list1):
    s1 = set()
    for i1 in list1:
        ip = ipaddress.ip_address(i1.ip.address)
        ip_str = str(ipaddress.ip_address(ip))
        s1.add(ip_str + "_" + i1.vlan + "_" + i1.mac)

    return s1


class DefGwTest(unittest.TestCase):
    """
    Validate default router setting.
    """

    def setUp(self):
        self.gw_store = defaultdict(str)
        self.dhcp_gw_info = UplinkGatewayInfo(self.gw_store)

    def test_gw_ip_for_DHCP(self):
        self.assertEqual(self.dhcp_gw_info.get_gw_ip(), None)
        self.assertEqual(self.dhcp_gw_info.get_gw_mac(), None)

    def test_gw_ip_for_Ip_pool(self):
        self.dhcp_gw_info.read_default_gw()

        def_gw_cmd = "ip route show |grep default| awk '{print $3}'"
        p = subprocess.Popen(
            [def_gw_cmd],
            stdout=subprocess.PIPE,
            shell=True,
        )
        def_ip = p.stdout.read().decode("utf-8").strip()

        self.assertEqual(self.dhcp_gw_info.get_gw_ip(), str(def_ip))
        self.assertEqual(self.dhcp_gw_info.get_gw_mac(), '')
        mac1 = "11:22:33:44:55:66"
        self.dhcp_gw_info.update_mac(def_ip, mac1)
        self.assertEqual(self.dhcp_gw_info.get_gw_ip(), str(def_ip))
        self.assertEqual(self.dhcp_gw_info.get_gw_mac(), mac1)

        # updating IP with same address shld keep mac
        self.dhcp_gw_info.update_ip(def_ip)
        self.assertEqual(self.dhcp_gw_info.get_gw_ip(), str(def_ip))
        self.assertEqual(self.dhcp_gw_info.get_gw_mac(), mac1)

        ip1 = "1.2.3.4"
        self.dhcp_gw_info.update_ip(ip1)
        self.assertEqual(self.dhcp_gw_info.get_gw_mac(), '')

    def test_vlan_gw_info(self):
        ip1 = "1.2.3.4"
        vlan1 = "1"
        self.dhcp_gw_info.update_ip(ip1, vlan1)
        self.assertEqual(self.dhcp_gw_info.get_gw_ip(vlan1), str(ip1))
        mac1 = "11:22:33:44:55:66"
        self.dhcp_gw_info.update_mac(ip1, mac1, vlan1)
        self.assertEqual(self.dhcp_gw_info.get_gw_ip(vlan1), str(ip1))
        self.assertEqual(self.dhcp_gw_info.get_gw_mac(vlan1), mac1)
        # updating IP with same address shld keep mac
        self.dhcp_gw_info.update_mac(ip1, mac1, vlan1)
        self.assertEqual(self.dhcp_gw_info.get_gw_ip(vlan1), str(ip1))
        self.assertEqual(self.dhcp_gw_info.get_gw_mac(vlan1), mac1)

        # any IP update shld invalidate MAC address.
        ip2 = "2.2.3.4"
        self.dhcp_gw_info.update_ip(ip2, vlan1)
        self.assertEqual(self.dhcp_gw_info.get_gw_ip(vlan1), str(ip2))
        self.assertEqual(self.dhcp_gw_info.get_gw_mac(vlan1), '')

    def test_vlan_gw_info_none(self):
        ip1 = "1.2.3.4"
        vlan1 = "1"
        mac1 = "11:22:33:44:55:66"

        self.dhcp_gw_info.update_mac(ip1, mac1, vlan1)
        self.assertEqual(self.dhcp_gw_info.get_gw_ip(vlan1), str(ip1))
        self.assertEqual(self.dhcp_gw_info.get_gw_mac(vlan1), mac1)

        # check None IP or mac updates
        self.dhcp_gw_info.update_ip(None, vlan1)
        self.assertEqual(self.dhcp_gw_info.get_gw_ip(vlan1), str(ip1))
        self.assertEqual(self.dhcp_gw_info.get_gw_mac(vlan1), mac1)

        self.dhcp_gw_info.update_mac(ip1, None, vlan1)
        self.assertEqual(self.dhcp_gw_info.get_gw_ip(vlan1), str(ip1))
        self.assertEqual(self.dhcp_gw_info.get_gw_mac(vlan1), mac1)

        self.dhcp_gw_info.update_mac(None, mac1, vlan1)
        self.assertEqual(self.dhcp_gw_info.get_gw_ip(vlan1), str(ip1))
        self.assertEqual(self.dhcp_gw_info.get_gw_mac(vlan1), mac1)

        self.dhcp_gw_info.update_mac(ip1, '', vlan1)
        self.assertEqual(self.dhcp_gw_info.get_gw_ip(vlan1), str(ip1))
        self.assertEqual(self.dhcp_gw_info.get_gw_mac(vlan1), mac1)

        vlan2 = None
        mac2 = "22:22:33:44:55:66"

        self.dhcp_gw_info.update_mac(ip1, mac2, vlan2)
        self.assertEqual(self.dhcp_gw_info.get_gw_ip(vlan2), str(ip1))
        self.assertEqual(self.dhcp_gw_info.get_gw_mac(vlan2), mac2)

    def test_vlan_gw_info_list(self):
        ip1 = "1.2.3.4"
        vlan1 = "1"
        gw1 = _get_gw_info(ip1, '', vlan1)
        self.dhcp_gw_info.update_ip(ip1, vlan1)

        ip2 = "2.2.3.4"
        vlan2 = "2"
        gw2 = _get_gw_info(ip2, '', vlan2)
        self.dhcp_gw_info.update_ip(ip2, vlan2)

        ip3 = "3.2.3.4"
        gw3 = _get_gw_info(ip3)
        self.dhcp_gw_info.update_ip(ip3)

        ip4 = "4.2.3.4"
        vlan4 = "4"
        mac4 = "11:22:33:44:55:66"
        gw4 = _get_gw_info(ip4, mac4, vlan4)
        self.dhcp_gw_info.update_mac(ip4, mac4, vlan4)

        gw_list = self.dhcp_gw_info.get_all_router_ips()

        expeected = gw_list_to_set([gw1, gw2, gw3, gw4])
        self.assertEqual(gw_list_to_set(gw_list), expeected)

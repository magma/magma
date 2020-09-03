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
import netifaces

from typing import MutableMapping, Optional, List
from lte.protos.mobilityd_pb2 import GWInfo, IPAddress

NO_VLAN = "NO_VLAN"


def _get_vlan_key(vlan: Optional[str]) -> str:
    if vlan is None or vlan == '' or vlan == NO_VLAN or vlan == "0":
        return NO_VLAN
    if int(vlan) < 0 or int(vlan) > 4095:
        raise InvalidVlanId("invalid vlan: " + vlan)

    return vlan


# TODO: move helper class to separate directory.
class UplinkGatewayInfo:
    def __init__(self, gw_info_map: MutableMapping[str, GWInfo]):
        """
            This maintains uptodate information about upstream GW.

        Args:
            gw_info_map: map to store GW info.
        """
        self._backing_map = gw_info_map

    # TODO: change vlan_id type to int
    def get_gw_ip(self, vlan_id: Optional[str] = "") -> Optional[str]:
        vlan_key = _get_vlan_key(vlan_id)
        if vlan_key in self._backing_map:
            gw_info = self._backing_map.get(vlan_key)
            ip = ipaddress.ip_address(gw_info.ip.address)
            return str(ip)

    def read_default_gw(self):
        gws = netifaces.gateways()
        logging.info("Using GW info: %s", gws)
        if gws is not None:
            default_gw = gws['default']
            if default_gw is not None and \
                    default_gw[netifaces.AF_INET] is not None:
                self.update_ip(default_gw[netifaces.AF_INET][0])

    def update_ip(self, ip: str, vlan_id: Optional[str] = ""):
        vlan_key = _get_vlan_key(vlan_id)

        logging.info("GW IP[%s]: %s" % (vlan_key, ip))
        ip_addr = ipaddress.ip_address(ip)
        gw_ip = IPAddress(version=IPAddress.IPV4,
                          address=ip_addr.packed)
        # keep mac address same if its same GW IP
        if vlan_key in self._backing_map:
            gw_info = self._backing_map[vlan_key]
            if gw_info and gw_info.ip == gw_ip:
                logging.debug("IP update: no change %s", ip)
                return

        updated_info = GWInfo(ip=gw_ip, mac="", vlan=vlan_id)
        self._backing_map[vlan_key] = updated_info

    def get_gw_mac(self, vlan_id: Optional[str] = "") -> Optional[str]:
        vlan_key = _get_vlan_key(vlan_id)

        if vlan_key in self._backing_map:
            return self._backing_map.get(vlan_key).mac
        else:
            return None

    def update_mac(self, ip: str, mac: Optional[str], vlan_id: Optional[str] = ""):
        vlan_key = _get_vlan_key(vlan_id)

        # TODO: enhance check for MAC address sanity.
        if mac is None or ':' not in mac:
            logging.error("Incorrect mac format: %s for IP %s (vlan_key %s)",
                          mac, ip, vlan_id)
            return
        ip_addr = ipaddress.ip_address(ip)
        gw_ip = IPAddress(version=IPAddress.IPV4,
                          address=ip_addr.packed)
        updated_info = GWInfo(ip=gw_ip, mac=mac, vlan=vlan_id)
        self._backing_map[vlan_key] = updated_info

    def get_all_router_ips(self) -> List[GWInfo]:
        return list(self._backing_map.values())


class InvalidVlanId(Exception):
    pass

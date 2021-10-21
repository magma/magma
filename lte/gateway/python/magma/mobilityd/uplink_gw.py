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
import threading
from typing import List, MutableMapping, Optional

import netifaces
from lte.protos.mobilityd_pb2 import GWInfo, IPAddress

NO_VLAN = "NO_VLAN"


def _get_vlan_key(vlan) -> str:
    # Validate vlan id is valid VLAN
    vlan_id_parsed = 0
    try:
        if vlan:
            vlan_id_parsed = int(vlan)
    except ValueError:
        logging.debug("invalid vlan id: %s", vlan)

    if vlan_id_parsed == 0:
        return NO_VLAN
    if vlan_id_parsed < 0 or vlan_id_parsed > 4095:
        raise InvalidVlanId("invalid vlan: " + str(vlan))

    return str(vlan)


# TODO: move helper class to separate directory.
class UplinkGatewayInfo:
    def __init__(self, gw_info_map: MutableMapping[str, GWInfo]):
        """
            This maintains uptodate information about upstream GW.
            The GW table is keyed bt vlan-id.

        Args:
            gw_info_map: map to store GW info.
        """
        self._backing_map = gw_info_map
        self._read_default_gw_timer = None
        self._read_default_gw_interval_seconds = 20

    def get_gw_ip(self, vlan_id: Optional[str] = "") -> Optional[str]:
        """
        Retrieve gw IP address
        Args:
            vlan_id: vlan if of the GW.
        """
        vlan_key = _get_vlan_key(vlan_id)
        if vlan_key in self._backing_map:
            gw_info = self._backing_map.get(vlan_key)
            ip = ipaddress.ip_address(gw_info.ip.address)
            return str(ip)

    def read_default_gw(self):
        self._do_read_default_gw()

    def _do_read_default_gw(self):
        gws = netifaces.gateways()
        logging.info("Using GW info: %s", gws)
        if gws is not None:
            default_gw = gws.get('default', None)
            gw_ip_addr = None
            if default_gw is not None:
                gw_ip_addr = default_gw.get(netifaces.AF_INET, None)
            if gw_ip_addr is not None:
                self.update_ip(gw_ip_addr[0])
                logging.info("GW probe: timer stopped")
                self._read_default_gw_timer = None
                return

        self._read_default_gw_timer = threading.Timer(
            self._read_default_gw_interval_seconds,
            self._do_read_default_gw,
        )
        self._read_default_gw_timer.start()
        logging.info("GW probe: timer started")

    def update_ip(self, ip: Optional[str], vlan_id=None):
        """
        Update IP address of the GW in mobilityD GW table.
        Args:
            ip: gw ip address
            vlan_id: vlan of the GW, None in case of no vlan used.
        """
        try:
            ip_addr = ipaddress.ip_address(ip)
        except ValueError:
            logging.debug("could not parse GW IP: %s", ip)
            return

        gw_ip = IPAddress(
            version=IPAddress.IPV4,
            address=ip_addr.packed,
        )
        # keep mac address same if its same GW IP
        vlan_key = _get_vlan_key(vlan_id)
        if vlan_key in self._backing_map:
            gw_info = self._backing_map[vlan_key]
            if gw_info and gw_info.ip == gw_ip:
                logging.debug("GW update: no change %s", ip)
                return

        updated_info = GWInfo(ip=gw_ip, mac="", vlan=vlan_key)
        self._backing_map[vlan_key] = updated_info
        logging.info("GW update: GW IP[%s]: %s" % (vlan_key, ip))

    def get_gw_mac(self, vlan_id: Optional[str] = None) -> Optional[str]:
        """
        Retrieve Mac address of default gw.
        Args:
            vlan_id: vlan of the gw, None if GW is not in a vlan.
        """
        vlan_key = _get_vlan_key(vlan_id)
        if vlan_key in self._backing_map:
            return self._backing_map.get(vlan_key).mac
        else:
            return None

    def update_mac(self, ip: Optional[str], mac: Optional[str], vlan_id=None):
        """
        Update mac address of GW in mobilityD GW table
        Args:
            ip: gw ip address.
            vlan_id: Vlan of the gw.
            mac: mac address of the GW.
        """
        try:
            ip_addr = ipaddress.ip_address(ip)
        except ValueError:
            logging.debug("could not parse GW IP: %s", ip)
            return
        vlan_key = _get_vlan_key(vlan_id)

        # TODO: enhance check for MAC address sanity.
        if mac is None or ':' not in mac:
            logging.error(
                "Incorrect mac format: %s for IP %s (vlan_key %s)",
                mac, ip, vlan_id,
            )
            return
        gw_ip = IPAddress(
            version=IPAddress.IPV4,
            address=ip_addr.packed,
        )
        updated_info = GWInfo(ip=gw_ip, mac=mac, vlan=vlan_key)
        self._backing_map[vlan_key] = updated_info
        logging.info("GW update: GW IP[%s]: %s : mac %s" % (vlan_key, ip, mac))

    def get_all_router_ips(self) -> List[GWInfo]:
        return list(self._backing_map.values())


class InvalidVlanId(Exception):
    pass

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
import netifaces

from typing import MutableMapping, Optional

DHCP_Router_key = "DHCPRouterKey"
DHCP_Router_Mac_key = "DHCPRouterMacKey"

# TODO: move helper class to separate directory.
class UplinkGatewayInfo:
    def __init__(self, gw_info_map: MutableMapping[str, str]):
        """
            This maintains uptodate information about upstream GW.

        Args:
            gw_info_map: map to store GW info.
        """
        self._backing_map = gw_info_map
        self._current_ip = None
        self._current_mac = None

    def getIP(self) -> Optional[str]:
        if DHCP_Router_key in self._backing_map:
            self._current_ip = self._backing_map.get(DHCP_Router_key)
            return self._current_ip

    def read_default_gw(self):
        gws = netifaces.gateways()
        logging.info("Using GW info: %s", gws)
        if gws is not None:
            default_gw = gws['default']
            if default_gw is not None and \
                    default_gw[netifaces.AF_INET] is not None:
                self.update_ip(default_gw[netifaces.AF_INET][0])

    def ip_exists(self) -> bool:
        return DHCP_Router_key in self._backing_map

    def update_ip(self, value: str):
        logging.info("GW IP: %s" % value)
        if self._current_ip != value:
            self._backing_map[DHCP_Router_key] = value
            self._current_ip = value

    def getMac(self) -> Optional[str]:
        if DHCP_Router_Mac_key in self._backing_map:
            self._current_mac = self._backing_map.get(DHCP_Router_Mac_key)
            return self._current_mac

    def mac_exists(self) -> str:
        return DHCP_Router_Mac_key in self._backing_map

    def update_mac(self, value: str):
        # TODO: enhance check for MAC address sanity.
        if ':' not in value:
            return

        if self._current_mac != value:
            self._backing_map[DHCP_Router_Mac_key] = value
            self._current_mac = value

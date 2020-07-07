"""
Copyright (c) 2020-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.

"""
import logging
from typing import MutableMapping, Optional

DHCP_Router_key = "DHCPRouterKey"
DHCP_Router_Mac_key = "DHCPRouterMacKey"


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
        if self._current_mac != value:
            self._backing_map[DHCP_Router_Mac_key] = value
            self._current_mac = value

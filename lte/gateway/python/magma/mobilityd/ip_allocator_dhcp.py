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

from __future__ import (
    absolute_import,
    division,
    print_function,
    unicode_literals,
)

import logging
from copy import deepcopy
from ipaddress import ip_address, ip_network
from threading import Condition
from typing import List

from magma.mobilityd.ip_descriptor import IPDesc, IPState, IPType

from .dhcp_client import DHCPClient
from .dhcp_desc import DHCPDescriptor, DHCPState
from .ip_allocator_base import IPAllocator, NoAvailableIPError
from .mac import MacAddress, create_mac_from_sid
from .mobility_store import MobilityStore

DEFAULT_DHCP_REQUEST_RETRY_FREQUENCY = 10
DEFAULT_DHCP_REQUEST_RETRY_DELAY = 1

LOG = logging.getLogger('mobilityd.dhcp.alloc')


class IPAllocatorDHCP(IPAllocator):
    def __init__(
        self, store: MobilityStore, retry_limit: int = 300,
        iface: str = "dhcp0",
    ):
        """
        Allocate IP address for SID using DHCP server.
        SID is mapped to MAC address using function defined in mac.py
        then this mac address used in DHCP request to allocate new IP
        from DHCP server.
        This IP is also cached to improve performance in case of
        reallocation for same SID in short period of time.

        Args:
            store: Moblityd storage instance
            retry_limit: try DHCP request
            iface: DHCP interface.
        """
        self._store = store
        self.dhcp_wait = Condition()
        self._dhcp_client = DHCPClient(
            dhcp_wait=self.dhcp_wait,
            dhcp_store=store.dhcp_store,
            gw_info=store.dhcp_gw_info,
            iface=iface,
        )
        self._retry_limit = retry_limit  # default wait for two minutes
        self._dhcp_client.run()

    def add_ip_block(self, ipblock: ip_network):
        logging.warning(
            "No need to allocate block for DHCP allocator: %s",
            ipblock,
        )

    def remove_ip_blocks(
        self, *ipblocks: List[ip_network],
        force: bool = False
    ) -> List[ip_network]:
        logging.warning(
            "Trying to delete ipblock from DHCP allocator: %s",
            ipblocks,
        )
        return []

    def list_added_ip_blocks(self) -> List[ip_network]:
        return list(deepcopy(self._store.assigned_ip_blocks))

    def list_allocated_ips(self, ipblock: ip_network) -> List[ip_address]:
        """ List IP addresses allocated from a given IP block

        Args:
            ipblock (ipaddress.ip_network): ip network to add
            e.g. ipaddress.ip_network("10.0.0.0/24")

        Return:
            list of IP addresses (ipaddress.ip_address)

        """
        return [
            ip for ip in
            self._store.ip_state_map.list_ips(IPState.ALLOCATED)
            if ip in ipblock
        ]

    def alloc_ip_address(self, sid: str, vlan: int) -> IPDesc:
        """
        Assumption: one-to-one mappings between SID and IP.

        Args:
            sid (string): universal subscriber id
            vlan: vlan of the APN

        Returns:
            ipaddress.ip_address: IP address allocated

        Raises:
            NoAvailableIPError: if run out of available IP addresses
        """
        mac = create_mac_from_sid(sid)

        dhcp_desc = self._dhcp_client.get_dhcp_desc(mac, vlan)
        LOG.debug(
            "allocate IP for %s mac %s dhcp_desc %s", sid, mac,
            dhcp_desc,
        )

        if dhcp_allocated_ip(dhcp_desc) is not True:
            dhcp_desc = self._alloc_ip_address_from_dhcp(mac, vlan)

        if dhcp_allocated_ip(dhcp_desc):
            ip_block = ip_network(dhcp_desc.subnet)
            ip_desc = IPDesc(
                ip=ip_address(dhcp_desc.ip),
                state=IPState.ALLOCATED,
                sid=sid,
                ip_block=ip_block,
                ip_type=IPType.DHCP,
                vlan_id=vlan,
            )
            self._store.assigned_ip_blocks.add(ip_block)

            return ip_desc
        else:
            msg = "No available IP addresses From DHCP for SID: {} MAC {}".format(
                sid, mac,
            )
            raise NoAvailableIPError(msg)

    def release_ip(self, ip_desc: IPDesc):
        """
        Release IP address, this involves following steps.
        1. send DHCP protocol packet to release the IP.
        2. update IP block list.
        3. update IP from ip-state.

        Args:
            ip_desc: release needs following info from IPDesc.
                SID used to get mac address, IP assigned to this SID,
                IP block of the IP address, vlan id of the APN.
        Returns: None
        """
        self._dhcp_client.release_ip_address(
            create_mac_from_sid(ip_desc.sid),
            ip_desc.vlan_id,
        )
        # Remove the IP from free IP list, since DHCP is the
        # owner of this IP
        self._store.ip_state_map.remove_ip_from_state(ip_desc.ip, IPState.FREE)

        list_allocated_ips = self._store.ip_state_map.list_ips(
            IPState.ALLOCATED,
        )
        for ipaddr in list_allocated_ips:
            if ipaddr in ip_desc.ip_block:
                # found the IP, do not remove this ip_block
                return

        ip_block_network = ip_network(ip_desc.ip_block)
        if ip_block_network in self._store.assigned_ip_blocks:
            self._store.assigned_ip_blocks.remove(ip_block_network)
        logging.debug(
            "del: _assigned_ip_blocks %s ipblock %s",
            self._store.assigned_ip_blocks, ip_desc.ip_block,
        )

    def stop_dhcp_sniffer(self):
        self._dhcp_client.stop()

    def _alloc_ip_address_from_dhcp(
        self, mac: MacAddress,
        vlan: int,
    ) -> DHCPDescriptor:
        retry_count = 0
        with self.dhcp_wait:
            dhcp_desc = None
            while (
                retry_count < self._retry_limit
                and dhcp_allocated_ip(dhcp_desc) is not True
            ):

                if retry_count % DEFAULT_DHCP_REQUEST_RETRY_FREQUENCY == 0:
                    self._dhcp_client.send_dhcp_packet(
                        mac, vlan,
                        DHCPState.DISCOVER,
                    )
                self.dhcp_wait.wait(timeout=DEFAULT_DHCP_REQUEST_RETRY_DELAY)

                dhcp_desc = self._dhcp_client.get_dhcp_desc(mac, vlan)

                retry_count = retry_count + 1

            return dhcp_desc


def dhcp_allocated_ip(dhcp_desc) -> bool:
    return dhcp_desc is not None and dhcp_desc.ip_is_allocated()

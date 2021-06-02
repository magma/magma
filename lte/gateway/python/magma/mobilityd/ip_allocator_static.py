"""
Copyright 2020 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

This is one of ip allocator for ip address manager.
The IP allocator accepts IP blocks (range of IP addresses), and supports
allocating and releasing IP addresses from the assigned IP blocks.
"""

from __future__ import (
    absolute_import,
    division,
    print_function,
    unicode_literals,
)

import logging
from ipaddress import ip_address, ip_network
from typing import List, Optional

from lte.protos.subscriberdb_pb2_grpc import SubscriberDBStub
from magma.mobilityd.ip_allocator_base import (
    DuplicateIPAssignmentError,
    IPAllocator,
)
from magma.mobilityd.ip_descriptor import IPDesc, IPState, IPType
from magma.mobilityd.mobility_store import MobilityStore
from magma.mobilityd.subscriberdb_client import SubscriberDbClient

DEFAULT_IP_RECYCLE_INTERVAL = 15


class IPAllocatorStaticWrapper(IPAllocator):

    def __init__(
        self, store: MobilityStore,
        subscriberdb_rpc_stub: SubscriberDBStub,
        ip_allocator: IPAllocator,
    ):
        """ Initializes a static IP allocator
            This is wrapper around other configured Ip allocator. If subscriber
            does have static IP, it uses underlying IP allocator to allocate IP
            for the subscriber.
        """
        self._store = store
        self._subscriber_client = SubscriberDbClient(subscriberdb_rpc_stub)
        self._ip_allocator = ip_allocator

    def add_ip_block(self, ipblock: ip_network):
        """ Add a block of IP addresses to the free IP list
        """
        self._ip_allocator.add_ip_block(ipblock)

    def remove_ip_blocks(
        self, ipblocks: List[ip_network],
        force: bool = False,
    ) -> List[ip_network]:
        """ Remove allocated IP blocks.
        """
        return self._ip_allocator.remove_ip_blocks(ipblocks, force)

    def list_added_ip_blocks(self) -> List[ip_network]:
        """ List IP blocks added to the IP allocator
        Return:
             copy of the list of assigned IP blocks
        """
        return self._ip_allocator.list_added_ip_blocks()

    def list_allocated_ips(self, ipblock: ip_network) -> List[ip_address]:
        """ List IP addresses allocated from a given IP block
        """
        return self._ip_allocator.list_allocated_ips(ipblock)

    def alloc_ip_address(self, sid: str, vlan: int) -> IPDesc:
        """ Check if subscriber has static IP assigned.
        If it is not allocated use IP allocator to assign an IP.
        """
        ip_desc = self._allocate_static_ip(sid)
        if ip_desc is None:
            ip_desc = self._ip_allocator.alloc_ip_address(sid, vlan)
        return ip_desc

    def release_ip(self, ip_desc: IPDesc):
        """
        Statically allocated IPs do not need to do any update on
        ip release
        """
        if ip_desc.type == IPType.STATIC:
            self._store.ip_state_map.remove_ip_from_state(
                ip_desc.ip,
                IPState.FREE,
            )
            ip_block_network = ip_network(ip_desc.ip_block)
            if ip_block_network in self._store.assigned_ip_blocks:
                self._store.assigned_ip_blocks.remove(ip_block_network)
        else:
            self._ip_allocator.release_ip(ip_desc)

    def _allocate_static_ip(self, sid: str) -> Optional[IPDesc]:
        """
        Check if static IP allocation is enabled and then check
        subscriber DB for assigned static IP for the SID
        """
        ip_addr_info = self._subscriber_client.get_subscriber_ip(sid)
        if ip_addr_info is None:
            return None
        logging.debug(
            "Found static IP: sid: %s ip_addr_info: %s",
            sid, str(ip_addr_info),
        )
        # Validate static IP is not in any of IP pool.
        for ip_pool in self._store.assigned_ip_blocks:
            if ip_addr_info.ip in ip_pool:
                error_msg = "Static Ip {} Overlap with IP-POOL: {}".format(
                    ip_addr_info.ip, ip_pool,
                )
                logging.error(error_msg)
                raise DuplicateIPAssignmentError(error_msg)

        # update gw info if available.
        if ip_addr_info.net_info.gw_ip:
            self._store.dhcp_gw_info.update_ip(
                ip_addr_info.net_info.gw_ip,
                ip_addr_info.net_info.vlan,
            )
            # update mac if IP is present.
            if ip_addr_info.net_info.gw_mac != "":
                self._store.dhcp_gw_info.update_mac(
                    ip_addr_info.net_info.gw_ip,
                    ip_addr_info.net_info.gw_mac,
                    ip_addr_info.net_info.vlan,
                )
        ip_block = ip_network(ip_addr_info.ip)
        self._store.assigned_ip_blocks.add(ip_block)
        return IPDesc(
            ip=ip_addr_info.ip, state=IPState.ALLOCATED,
            sid=sid, ip_block=ip_block, ip_type=IPType.STATIC,
            vlan_id=ip_addr_info.net_info.vlan,
        )

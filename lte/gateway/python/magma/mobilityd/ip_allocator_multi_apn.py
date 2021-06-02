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

from ipaddress import ip_address, ip_network
from typing import List

from magma.mobilityd.ip_allocator_base import IPAllocator
from magma.mobilityd.ip_descriptor import IPDesc
from magma.mobilityd.mobility_store import MobilityStore
from magma.mobilityd.subscriberdb_client import SubscriberDbClient

DEFAULT_IP_RECYCLE_INTERVAL = 15


class IPAllocatorMultiAPNWrapper(IPAllocator):

    def __init__(
        self, store: MobilityStore,
        subscriberdb_rpc_stub,
        ip_allocator: IPAllocator,
    ):
        """ Initializes a Multi APN IP allocator
            This is wrapper around other configured Ip allocator. If subscriber
            has vlan configured in APN config, it would be used for allocating
            IP address by underlying IP allocator. For DHCP it means using
            vlan tag for DHCP request, for IP pool allocator it does not change
            behaviour.
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

    def alloc_ip_address(self, sid: str, vlan_id: int) -> IPDesc:
        """ Check if subscriber has APN configuration and vlan.
        once we have APN specific info use IP allocator to assign an IP.
        """
        network_info = self._subscriber_client.get_subscriber_apn_network_info(
            sid,
        )
        # Update GW info from subscriber DB data.
        # This could be overwritten by DHCP response.
        self._store.dhcp_gw_info.update_mac(
            network_info.gw_ip,
            network_info.gw_mac,
            network_info.vlan,
        )
        return self._ip_allocator.alloc_ip_address(sid, network_info.vlan)

    def release_ip(self, ip_desc: IPDesc):
        """
        Multi APN allocated IPs do not need to do any update on
        ip release, let actual IP allocator release the IP.
        """
        self._ip_allocator.release_ip(ip_desc)

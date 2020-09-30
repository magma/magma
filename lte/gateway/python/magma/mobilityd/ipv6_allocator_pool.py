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
import random

from ipaddress import ip_address, ip_network
from typing import List, Optional

import threading

from .ip_address_man import IPNotInUseError
from .ip_descriptor import IPDesc
from .ip_allocator_base import IPAllocator, OverlappedIPBlocksError
from .ip_descriptor import IPv6SessionAllocType

IPV6_PREFIX_PART_LEN = 64
IID_PART_LEN = 64
MAX_RAND_IID_CALC_TRIES = 5


class IPv6AllocatorPool(IPAllocator):
    def __init__(self, config):
        super().__init__()
        self._lock = threading.RLock()  # Re-entrant lock
        self._config = config
        self._assigned_ip_block = None
        self._allocated_iid = set()
        self._sid_session_prefix_allocated = {}
        self._sid_ips_map = {}
        self._ipv6_session_prefix_alloc_mode = self._config[
            'ipv6_session_prefix_alloc_mode']

    def add_ip_block(self, ipblock: ip_network) -> None:
        """
        Adds IP block to the assigned IP block of the IPv6 allocator
        :param ipblock: IPv6 ip_network object
        """
        with self._lock:
            if self._assigned_ip_block and ipblock.overlaps(
                    self._assigned_ip_block):
                logging.error("Overlapped IP block: %s", ipblock)
                raise OverlappedIPBlocksError(ipblock)

            # For now only one IPv6 network is supported
            self._assigned_ip_block = ipblock

    def remove_ip_blocks(self, *ipblocks: List[ip_network],
                         force: bool = False) -> List[ip_network]:
        """
        Removes assigned IP block (as it only supports one for now)
        :param ipblocks: List of ip_network objects to remove
        :param force:
        :return: Removed list of IP blocks
        """
        with self._lock:
            removed_blocks = []
            if self._assigned_ip_block in ipblocks:
                # Clear allocated session prefix and IID store
                self._allocated_iid.clear()
                self._sid_session_prefix_allocated.clear()

                removed_blocks.append(self._assigned_ip_block)
                self._assigned_ip_block = None
        return removed_blocks

    def alloc_ip_address(self, sid: str, _) -> ip_address:
        """
        Calculates full IPv6 Address from configured prefix part
        (assigned_ip_block) + unique session prefix part + interface
        identifier

        :param sid: composite SID (IMSI + apn)
        :param _:
        :return: Full IPv6 address
        """
        with self._lock:

            # Take available ipv6 host from network
            ipv6_addr_part = next(self._assigned_ip_block.hosts())

            # Calculate session part from rest of 64 prefix bits
            session_prefix_part = self._get_session_prefix(sid)

            # Get interface identifier from 64 bits fixed length
            iid_part = self._get_ipv6_iid_part(IID_PART_LEN)
            if not iid_part:
                logging.error('Could not get IPv6 IID for sid: %s', sid)
                raise MaxRandIIDCalculationError(
                    'Could not get IPv6 IID for sid: %s', sid)

            ipv6_addr = ipv6_addr_part + session_prefix_part + iid_part
            return ipv6_addr

    def release_ip(self, ip_desc: IPDesc):
        """
        Deallocates session prefix part from identifier and allocated IID
        :param ip_desc: IP Desc to remove

        """
        with self._lock:
            sid = ip_desc.sid
            ip_addr = ip_desc.ip
            ipv6_addr_part = int(next(self._assigned_ip_block.hosts()))

            session_prefix = self._sid_session_prefix_allocated.get(sid)
            if not session_prefix:
                raise IPNotInUseError('IP %s not allocated', ip_addr)

            if ip_addr in self._assigned_ip_block and session_prefix:
                iid_part = int(ip_addr) - ipv6_addr_part - int(session_prefix)

                if iid_part in self._allocated_iid:
                    del self._sid_session_prefix_allocated[sid]
                    self._allocated_iid.remove(iid_part)
                else:
                    raise IPNotInUseError('IP %s not allocated', ip_addr)

    def _get_ipv6_iid_part(self, length: int = 64) -> Optional[int]:
        """
        Calculates IPv6 Interface identifier using random calculation,
        if it exceeds MAX_RAND calculation tries, it returns None
        :param length: length for ipv6 IID
        :return: Randomized IID N bits
        """
        for i in range(MAX_RAND_IID_CALC_TRIES):
            rand_iid_bits = random.getrandbits(length)
            if rand_iid_bits not in self._allocated_iid:
                self._allocated_iid.add(rand_iid_bits)
                return rand_iid_bits
        return None

    def _get_session_prefix(self, sid: str) -> int:
        """
        Returns unique session prefix using configurable allocation mode

        :param sid:
        :return: Session prefix N bits
        """
        session_prefix_len = IPV6_PREFIX_PART_LEN - self._assigned_ip_block.prefixlen
        session_prefix_allocated = self._sid_session_prefix_allocated.get(sid)
        # TODO: Support multiple alloc modes
        if self._ipv6_session_prefix_alloc_mode == IPv6SessionAllocType.RANDOM:
            session_prefix_part = random.getrandbits(session_prefix_len)
            if session_prefix_part != session_prefix_allocated:
                self._sid_session_prefix_allocated[sid] = session_prefix_part
            else:
                session_prefix_part = session_prefix_allocated
            return session_prefix_part

    def list_added_ip_blocks(self) -> List[ip_network]:
        """
        Returns assigned IP blocks on the allocator
        :return: list of allocated ip_networks
        """
        return [self._assigned_ip_block]

    def list_allocated_ips(self, ipblock: ip_network) -> List[ip_address]:
        raise NotImplementedError


class MaxRandIIDCalculationError(Exception):
    """
    Exception thrown when calculation of random IID reaches maximum tries
    """
    pass

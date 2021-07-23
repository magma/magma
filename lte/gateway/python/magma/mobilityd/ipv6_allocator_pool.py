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

from .ip_allocator_base import (
    IPAllocator,
    IPNotInUseError,
    NoAvailableIPError,
    OverlappedIPBlocksError,
)
from .ip_descriptor import IPDesc, IPState, IPType, IPv6SessionAllocType
from .mobility_store import MobilityStore

IPV6_PREFIX_PART_LEN = 64
IID_PART_LEN = 64
MAX_IPV6_CONF_PREFIX_LEN = 48
MAX_CALC_TRIES = 5


class IPv6AllocatorPool(IPAllocator):
    def __init__(self, store: MobilityStore, session_prefix_alloc_mode: str):
        super().__init__()
        self._store = store
        self._assigned_ip_block = None
        self._ipv6_session_prefix_alloc_mode = session_prefix_alloc_mode

    def add_ip_block(self, ipblock: ip_network):
        """
        Adds IP block to the assigned IP block of the IPv6 allocator

        Args:
            ipblock: ipv6 ip_network object

        Returns:
        """
        if self._assigned_ip_block and ipblock.overlaps(
                self._assigned_ip_block,
        ):
            raise OverlappedIPBlocksError(ipblock)

        if ipblock.prefixlen > MAX_IPV6_CONF_PREFIX_LEN:
            msg = "IPv6 block exceeds maximum allowed prefix length"
            logging.error(msg)
            raise InvalidIPv6NetworkError(msg)

        # For now only one IPv6 network is supported
        self._assigned_ip_block = ipblock

    def remove_ip_blocks(
        self, ipblocks: List[ip_network],
        force: bool = False,
    ) -> List[ip_network]:
        """
        Removes assigned IP block (as it only supports one for now)

        If force is False, blocks that have any addresses currently allocated
        will not be removed. Otherwise, if force is True, the indicated blocks
        will be removed regardless of whether any addresses have been allocated
        and any allocated addresses will no longer be served.

        Args:
            ipblocks: List of ip_network objects to remove
            force:

        Returns: Removed list of IP blocks
        """
        if self._assigned_ip_block not in ipblocks:
            return []

        if not force:
            allocated_ip_block_set = self._store.ipv6_state_map.get_allocated_ip_block_set()
            if allocated_ip_block_set:
                return []

        removed_blocks = []
        # Clear allocated session prefix and IID store
        self._store.allocated_iid.clear()
        self._store.sid_session_prefix_allocated.clear()

        for sid in list(self._store.sid_ips_map):
            ip_desc = self._store.sid_ips_map[sid]
            if ip_desc.ip.version == 6:
                self._store.ipv6_state_map.remove_ip_from_state(
                    ip_desc.ip,
                    IPState.FREE,
                )
                if force:
                    self._store.ipv6_state_map.remove_ip_from_state(
                        ip_desc.ip, IPState.ALLOCATED,
                    )
                self._store.sid_ips_map.pop(sid)

        removed_blocks.append(self._assigned_ip_block)
        logging.info(
            'Removed IP block %s from IPv6 address pool',
            self._assigned_ip_block,
        )
        self._assigned_ip_block = None
        return removed_blocks

    def alloc_ip_address(self, sid: str, _) -> IPDesc:
        """
        Calculates full IPv6 Address from configured prefix part
        (assigned_ip_block) + unique session prefix part + interface
        identifier

        Args:
            sid: composite SID (IMSI + apn)
            _:

        Returns: full ipv6 address
        """
        if not self._assigned_ip_block:
            raise NoAvailableIPError('No IP block assigned to the allocator')
        # Take available ipv6 host from network
        ipv6_addr_part = next(self._assigned_ip_block.hosts())

        # Calculate session part from rest of 64 prefix bits
        session_prefix_part = self._get_session_prefix(sid)
        if not session_prefix_part:
            logging.error('Could not get IPv6 session prefix for sid: %s', sid)
            raise MaxCalculationError(
                'Could not get IPv6 session prefix for sid: %s', sid,
            )

        # Get interface identifier from 64 bits fixed length
        iid_part = self._get_ipv6_iid_part(sid, IID_PART_LEN)
        if not iid_part:
            logging.error('Could not get IPv6 IID for sid: %s', sid)
            raise MaxCalculationError(
                'Could not get IPv6 IID for sid: %s', sid,
            )

        ipv6_addr = ipv6_addr_part + (session_prefix_part * iid_part)
        ip_desc = IPDesc(
            ip=ipv6_addr, state=IPState.ALLOCATED, sid=sid,
            ip_block=self._assigned_ip_block,
            ip_type=IPType.IP_POOL,
        )
        return ip_desc

    def release_ip(self, ip_desc: IPDesc):
        """
        Deallocates session prefix part from identifier and allocated IID

        Args:
            ip_desc: IPDesc to remove

        Returns:
        """
        sid = ip_desc.sid
        ip_addr = ip_desc.ip
        ipv6_addr_part = int(next(self._assigned_ip_block.hosts()))

        session_prefix = self._store.sid_session_prefix_allocated.get(sid)
        if not session_prefix:
            raise IPNotInUseError('IP %s not allocated', ip_addr)

        if ip_addr in self._assigned_ip_block and session_prefix:
            # Extract IID part of the given IPv6 prefix and session prefix
            iid_part = float(
                (int(ip_addr) - ipv6_addr_part) / int(session_prefix),
            )

            if iid_part in self._store.allocated_iid.values():
                del self._store.sid_session_prefix_allocated[sid]
                del self._store.allocated_iid[sid]
            else:
                raise IPNotInUseError('IP %s not allocated', ip_addr)

    def _get_ipv6_iid_part(self, sid: str, length: int = 64) -> Optional[int]:
        """
        Calculates IPv6 Interface identifier using random calculation,
        if it exceeds MAX_RAND calculation tries, it returns None

        Args:
            sid: composite SID
            length: length for ipv6 IID

        Returns: IID N bits
        """
        for i in range(MAX_CALC_TRIES):
            rand_iid_bits = random.getrandbits(length)
            if rand_iid_bits not in self._store.allocated_iid.values():
                self._store.allocated_iid[sid] = float(rand_iid_bits)
                return rand_iid_bits
        return None

    def _get_session_prefix(self, sid: str) -> Optional[int]:
        """
        Returns unique session prefix using configurable allocation mode,
        it return None if it exceeds MAX calculation tries

        Args:
            sid:

        Returns: session prefix N bits
        """
        session_prefix_len = IPV6_PREFIX_PART_LEN - self._assigned_ip_block.prefixlen
        session_prefix_allocated = self._store.sid_session_prefix_allocated.get(
            sid,
        )
        # TODO: Support multiple alloc modes
        if self._ipv6_session_prefix_alloc_mode == IPv6SessionAllocType.RANDOM:
            for i in range(MAX_CALC_TRIES):
                session_prefix_part = random.getrandbits(session_prefix_len)
                if session_prefix_part != session_prefix_allocated:
                    self._store.sid_session_prefix_allocated[
                        sid
                    ] = session_prefix_part
                    return session_prefix_part
        return None

    def list_added_ip_blocks(self) -> List[ip_network]:
        """
        Returns: assigned IP blocks on the allocator
        """
        return [self._assigned_ip_block]

    def list_allocated_ips(self, ipblock: ip_network) -> List[ip_address]:
        raise NotImplementedError


class InvalidIPv6NetworkError(Exception):
    """
    Exception thrown if allocated IPv6 network is invalid
    """
    pass


class MaxCalculationError(Exception):
    """
    Exception thrown when calculation of session prefix / IID reaches
    maximum tries
    """
    pass

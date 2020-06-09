"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.


This is one of ip allocator for ip address manager.
The IP allocator accepts IP blocks (range of IP addresses), and supports
allocating and releasing IP addresses from the assigned IP blocks.
"""

from __future__ import absolute_import, division, print_function, \
    unicode_literals

import logging
from collections import defaultdict
from ipaddress import ip_address, ip_network
from typing import List, Set

from copy import deepcopy
from magma.mobilityd.ip_descriptor import IPDesc, IPState
from .ip_descriptor_map import IpDescriptorMap
from .ip_allocator_base import IPAllocator, NoAvailableIPError, \
    IPBlockNotFoundError, OverlappedIPBlocksError


DEFAULT_IP_RECYCLE_INTERVAL = 15


class IpAllocatorStatic(IPAllocator):

    def __init__(self,
                 assigned_ip_blocks: Set[ip_network],
                 ip_state_map: IpDescriptorMap,
                 sid_ips_map: defaultdict):
        """ Initializes a new IP allocator
        """
        self._assigned_ip_blocks = assigned_ip_blocks  # {ip_block}
        self._ip_state_map = ip_state_map  # {state=>{ip=>ip_desc}}
        self._sid_ips_map = sid_ips_map  # {SID=>IPDesc}

    def add_ip_block(self, ipblock: ip_network):
        """ Add a block of IP addresses to the free IP list

        IP blocks should not overlap.

        Args:
            ipblock (ipaddress.ip_network): ip network to add
            e.g. ipaddress.ip_network("10.0.0.0/24")

        Raises:
            OverlappedIPBlocksError: if the given IP block overlaps with
            existing ones
        """
        for blk in self._assigned_ip_blocks:
            if ipblock.overlaps(blk):
                logging.error("Overlapped IP block: %s", ipblock)
                raise OverlappedIPBlocksError(ipblock)

        self._assigned_ip_blocks.add(ipblock)
        # TODO(oramadan) t23793559 HACK reserve the GW address for
        #  gtp_br0 iface and test VM
        num_reserved_addresses = 11
        for ip in ipblock.hosts():
            state = IPState.RESERVED if num_reserved_addresses > 0 \
                else IPState.FREE
            ip_desc = IPDesc(ip=ip, state=state,
                             ip_block=ipblock, sid=None)
            self._ip_state_map.add_ip_to_state(ip, ip_desc, state)
            if num_reserved_addresses > 0:
                num_reserved_addresses -= 1

    def remove_ip_blocks(self, ipblocks: List[ip_network],
                         _force: bool = False) -> List[ip_network]:
        """ Makes the indicated block(s) unavailable for allocation

        If force is False, blocks that have any addresses currently allocated
        will not be removed. Otherwise, if force is True, the indicated blocks
        will be removed regardless of whether any addresses have been allocated
        and any allocated addresses will no longer be served.

        Removing a block entails removing the IP addresses within that block
        from the internal state machine.

        Args:
            ipblocks (ipaddress.ip_network): variable number of objects of type
                ipaddress.ip_network, representing the blocks that are intended
                to be removed. The blocks should have been explicitly added and
                not yet removed. Any blocks that are not active in the IP
                allocator will be ignored with a warning.
            _force (bool): whether to forcibly remove the blocks indicated. If
                False, will only remove a block if no addresses from within the
                block have been allocated. If True, will remove all blocks
                regardless of whether any addresses have been allocated from
                them.

        Returns a set of the blocks that have been successfully removed.
        """

        remove_blocks = set(ipblocks) & self._assigned_ip_blocks
        logging.warning("_assigned_ip_blocks %s", self._assigned_ip_blocks)
        logging.warning("arg ipblocks %s", ipblocks)

        extraneous_blocks = set(ipblocks) ^ remove_blocks
        # check unknown ip blocks
        if extraneous_blocks:
            logging.warning("Cannot remove unknown IP block(s): %s",
                            extraneous_blocks)
        del extraneous_blocks

        # "soft" removal does not remove blocks have IPs allocated
        if not _force:
            allocated_ip_block_set = self._ip_state_map.get_allocated_ip_block_set()
            remove_blocks -= allocated_ip_block_set
            del allocated_ip_block_set

        # Remove the associated IP addresses
        remove_ips = \
            (ip for block in remove_blocks for ip in block.hosts())
        for ip in remove_ips:
            for state in (IPState.FREE, IPState.RELEASED, IPState.REAPED):
                self._ip_state_map.remove_ip_from_state(ip, state)
            if _force:
                self._ip_state_map.remove_ip_from_state(ip, IPState.ALLOCATED)
            else:
                assert not self._ip_state_map.test_ip_state(ip, IPState.ALLOCATED), \
                    "Unexpected ALLOCATED IP %s from a soft IP block " \
                    "removal "

            # Clean up SID maps
            for sid in list(self._sid_ips_map):
                self._sid_ips_map.pop(sid)

        # Remove the IP blocks
        self._assigned_ip_blocks -= remove_blocks

        # Can't use generators here
        remove_sids = tuple(sid for sid in self._sid_ips_map
                            if not self._sid_ips_map[sid])
        for sid in remove_sids:
            self._sid_ips_map.pop(sid)

        for block in remove_blocks:
            logging.info('Removed IP block %s from IPv4 address pool', block)
        return remove_blocks

    def list_added_ip_blocks(self) -> List[ip_network]:
        """ List IP blocks added to the IP allocator

        Return:
             copy of the list of assigned IP blocks
        """
        return list(deepcopy(self._assigned_ip_blocks))

    def list_allocated_ips(self, ipblock: ip_network) -> List[ip_address]:
        """ List IP addresses allocated from a given IP block

        Args:
            ipblock (ipaddress.ip_network): ip network to add
            e.g. ipaddress.ip_network("10.0.0.0/24")

        Return:
            list of IP addresses (ipaddress.ip_address)

        Raises:
          IPBlockNotFoundError: if the given IP block is not found in the
          internal list
        """
        if ipblock not in self._assigned_ip_blocks:
            logging.error("Listing an unknown IP block: %s", ipblock)
            raise IPBlockNotFoundError(ipblock)

        res = [ip for ip in ipblock \
               if self._ip_state_map.test_ip_state(ip, IPState.ALLOCATED)]
        return res

    def alloc_ip_address(self, sid: str) -> IPDesc:
        """ Allocate an IP address from the free list

        Assumption: one-to-one mappings between SID and IP.

        Args:
            sid (string): universal subscriber id

        Returns:
            ipaddress.ip_address: IP address allocated

        Raises:
            NoAvailableIPError: if run out of available IP addresses
            DuplicatedIPAllocationError: if an IP has been allocated to a UE
                with the same IMSI
        """
        # if an IP is not yet allocated for the UE, allocate a new IP
        if self._ip_state_map.get_ip_count(IPState.FREE):
            return self._ip_state_map.pop_ip_from_state(IPState.FREE)
        else:
            logging.error("Run out of available IP addresses")
            raise NoAvailableIPError("No available IP addresses")

    def release_ip(self, sid):
        pass

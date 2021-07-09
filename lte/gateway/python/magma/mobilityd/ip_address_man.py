"""
Copyright 2020 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

The IP allocator maintains the life cycle of assigned IP addresses.

The IP address manager supports allocating and releasing IP addresses
Note that an IP address is not immediately made available for allocation right
after release: it is "reserved" for the same client for a certain period of
time to ensure that 1) an observer, e.g. pipelined, that caches IP states has
enough time to pull the updated IP states; 2) IP packets intended for the
old client will not be unintentionally routed to a new client until the old
TCP connection expires.

We have plug in design for IP address allocator. today we have two allocator
1. IP_POOL allocator, which allocates IP address from block of IP address
2. DHCP IP allocator. This allocates IP address assigned by DHCP server
   in network. This is typically used in Non NAT GW environment.

IP allocator supports static IP allocation for sub set of subscribers.
When static Ip mode is enabled, mobilityD would check subscriberDB for
assigned IP address before allocating Ip from configured allocator.

To support this semantic, an IP address can have the following states
during it's life cycle in the IP allocator:
    FREE: IP is available for allocation
    ALLOCATED: IP is allocated for a client.
    RELEASED: IP is released, but still reserved for the client
    REAPED: IPs are periodically reaped from the RELEASED state to the
        REAPED state, and at the same time a timer is set. All REAPED state
        IPs are freed once the time goes off. The purpose of this state is
        to age IPs for a certain period of time before freeing.
"""

from __future__ import (
    absolute_import,
    division,
    print_function,
    unicode_literals,
)

import ipaddress
import logging
import threading
from ipaddress import ip_address, ip_network
from typing import List, Optional, Tuple

from lte.protos.mobilityd_pb2 import GWInfo, IPAddress
from magma.mobilityd.ip_descriptor import IPState
from magma.mobilityd.metrics import IP_ALLOCATED_TOTAL, IP_RELEASED_TOTAL

from .ip_allocator_base import (
    DuplicateIPAssignmentError,
    IPAllocator,
    IPNotInUseError,
    MappingNotFoundError,
)
from .mobility_store import MobilityStore

DEFAULT_IP_RECYCLE_INTERVAL = 15


class IPAddressManager:
    """ A thread-safe IP allocator: all mutating functions are protected by a
    re-entrant lock.


    The IPAllocator maintains IP life cycle, as well as the mapping between
    SubscriberID (SID) and IPs. For now, only one-to-one mapping is supported
    between SID and IP.

    IP Recycling:
        An IP address is periodically recycled for allocation after releasing.

        Constraints:
        1. Define a maturity time T, in seconds. Any IPs that were released
           within the past T seconds cannot be freed.
        2. All released IPs must eventually be freed.

        To achieve these constraints, whenever we release an IP address we try
        to set a recycle timer if it's not already set. When setting a timer,
        we move all IPs in the RELEASED state to the REAPED state. Those IPs
        will be freed once the timer goes off, at which time those IPs are
        guaranteed to be "matured". Concurrently, newly released IPs are added
        to the RELEASED state.

        The time between when an IP is released and when it is freed is
        guaranteed to be between T and 2T seconds. The worst-case occurs when
        the release happens right after the previous timer has been initiated.
        This leads to an additional T seconds wait for the existing timer to
        time out and trigger a callback which initiate a new timer, on top of
        the T seconds of timeout for the next timer.

    Persistence to Redis:
        The IP allocator now by default persists its state to a redis-server
        process denoted in mobilityd.yml. The allocator's persisted properties
        and associated structure:
            - self._assigned_ip_blocks: {ip_block}
            - self.ip_state_map: {state=>{ip=>ip_desc}}
            - self.sid_ips_map: {SID=>[IPDesc]}

        The utilized redis_containers store a cache of state in local memory,
        so reads are the same speed as without persistence. For writes, state
        is serialized with provided serializers and written through to the
        Redis process. Used as expected, writes will be small (a few bytes).
        Redis's performance can be measured with the redis-benchmark tool,
        but we expect almost 100% of writes to take less than 1 millisecond.

    """

    def __init__(
        self,
        ipv4_allocator: IPAllocator,
        ipv6_allocator: IPAllocator,
        store: MobilityStore,
        recycling_interval: int = DEFAULT_IP_RECYCLE_INTERVAL,
    ):
        """ Initializes a new IP address manager

        Args:
            recycling_interval (number): minimum time, in seconds, before a
                released IP is recycled and freed into the pool of available
                IPs.

                Default: None, no recycling will occur automatically.
        """
        self._lock = threading.RLock()  # re-entrant locks
        self._recycle_timer = None  # reference to recycle timer
        self._recycling_interval_seconds = recycling_interval

        # Load store dependencies
        self._store = store

        self.ip_allocator = ipv4_allocator
        self.ipv6_allocator = ipv6_allocator

    def add_ip_block(self, ipblock: ip_network):
        """ Add a block of IP addresses to the free IP list

        IP blocks should not overlap.

        Args:
            ipblock (ipaddress.ip_network): ip network to add
            e.g. ipaddress.ip_network("10.0.0.0/24")

        Raises:
            OverlappedIPBlocksError: if the given IP block overlaps with
            existing ones
            InvalidIPv6NetworkError: if IPv6 block is invalid
        """
        with self._lock:
            if ipblock.version == 4:
                self.ip_allocator.add_ip_block(ipblock)
                logging.info(
                    "Added block %s to the IPv4 address pool",
                    ipblock,
                )
            elif ipblock.version == 6:
                self.ipv6_allocator.add_ip_block(ipblock)
                logging.info(
                    "Added block %s to the IPv6 address pool",
                    ipblock,
                )
            else:
                logging.warning("Failing to add IPBlock as is invalid")

    def remove_ip_blocks(
        self, *ipblocks: List[ip_network],
        force: bool = False
    ) -> List[ip_network]:
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
            force (bool): whether to forcibly remove the blocks indicated. If
                False, will only remove a block if no addresses from within the
                block have been allocated. If True, will remove all blocks
                regardless of whether any addresses have been allocated from
                them.

        Returns a list of the blocks that have been successfully removed.
        """

        with self._lock:
            ipv4_blocks, ipv6_blocks = [], []
            ipv4_blocks_deleted, ipv6_blocks_deleted = [], []
            for b in ipblocks:
                if b.version == 4:
                    ipv4_blocks.append(b)
                elif b.version == 6:
                    ipv6_blocks.append(b)

            if ipv4_blocks:
                ipv4_blocks_deleted.extend(
                    self.ip_allocator.remove_ip_blocks(
                    ipv4_blocks, force=force,
                    ),
                )
            if ipv6_blocks:
                ipv6_blocks_deleted.extend(
                    self.ipv6_allocator.remove_ip_blocks(
                        ipv6_blocks, force=force,
                    ),
                )

        return ipv4_blocks_deleted + ipv6_blocks_deleted

    def list_added_ip_blocks(self) -> List[ip_network]:
        """ List IP blocks added to the IP allocator

        Return:
             copy of the list of assigned IP blocks
        """
        with self._lock:
            ip_blocks = self.ip_allocator.list_added_ip_blocks()
        return ip_blocks

    def get_assigned_ipv6_block(self) -> ip_network:
        """
        Returns currently assigned block to IPv6 allocator
        :return: ip_network object for assigned block
        """
        with self._lock:
            ip_block = self.ipv6_allocator.list_added_ip_blocks()[0]
        return ip_block

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

        with self._lock:
            res = self.ip_allocator.list_allocated_ips(ipblock)
        return res

    def alloc_ip_address(self, sid: str, version: int = IPAddress.IPV4) -> \
            Tuple[ip_address, int]:
        """ Allocate an IP address from the free list

        Assumption: one-to-one mappings between SID and IP.

        Args:
            sid (string): universal subscriber id
            version (int): version of IP to allocate

        Returns:
            ipaddress.ip_address: IP address allocated

        Raises:
            NoAvailableIPError: if run out of available IP addresses
            DuplicatedIPAllocationError: if an IP has been allocated to a UE
                with the same IMSI
        """

        with self._lock:
            # if an IP is reserved for the UE, this IP could be in the state of
            # ALLOCATED, RELEASED or REAPED.
            if sid in self._store.sid_ips_map:
                old_ip_desc = self._store.sid_ips_map[sid]
                if self.is_ip_in_state(old_ip_desc.ip, IPState.ALLOCATED):
                    # MME state went out of sync with mobilityd!
                    # Recover gracefully by allocating the same IP
                    logging.warning(
                        "Re-allocate IP %s for sid %s without "
                        "MME releasing it first ip-state-map",
                        old_ip_desc.ip,
                        sid,
                    )

                    # TODO: enable strict checking after root causing the
                    # issue in MME
                    # raise DuplicatedIPAllocationError(
                    #     "An IP has been allocated for this IMSI")
                elif self.is_ip_in_state(old_ip_desc.ip, IPState.RELEASED):
                    ip_desc = self._store.ip_state_map.mark_ip_state(
                        old_ip_desc.ip,
                        IPState.ALLOCATED,
                    )
                    ip_desc.sid = sid
                    logging.debug(
                        "SID %s IP %s RELEASED => ALLOCATED",
                        sid, old_ip_desc.ip,
                    )
                elif self.is_ip_in_state(old_ip_desc.ip, IPState.REAPED):
                    ip_desc = self._store.ip_state_map.mark_ip_state(
                        old_ip_desc.ip,
                        IPState.ALLOCATED,
                    )
                    ip_desc.sid = sid
                    logging.debug(
                        "SID %s IP %s REAPED => ALLOCATED",
                        sid, old_ip_desc.ip,
                    )
                else:
                    raise AssertionError("Unexpected internal state")

                logging.info(
                    "Allocating the same IP %s for sid %s",
                    old_ip_desc.ip, sid,
                )
                IP_ALLOCATED_TOTAL.inc()
                return old_ip_desc.ip, old_ip_desc.vlan_id

            # Now try to allocate it from underlying allocator.
            allocator = self.ip_allocator if version == IPAddress.IPV4 \
                else self.ipv6_allocator
            ip_desc = allocator.alloc_ip_address(sid, 0)
            existing_sid = self.get_sid_for_ip(ip_desc.ip)
            if existing_sid:
                error_msg = "Dup IP: {} for SID: {}, which already is " \
                            "assigned to SID: {}".format(
                                ip_desc.ip,
                                sid,
                                existing_sid,
                            )
                logging.error(error_msg)
                raise DuplicateIPAssignmentError(error_msg)

            if version == IPAddress.IPV4:
                self._store.ip_state_map.add_ip_to_state(
                    ip_desc.ip, ip_desc,
                    IPState.ALLOCATED,
                )
            elif version == IPAddress.IPV6:
                self._store.ipv6_state_map.add_ip_to_state(
                    ip_desc.ip, ip_desc,
                    IPState.ALLOCATED,
                )

            self._store.sid_ips_map[sid] = ip_desc

            logging.debug("Allocating New IP: %s", str(ip_desc))
            IP_ALLOCATED_TOTAL.inc()
            return ip_desc.ip, ip_desc.vlan_id

    def get_sid_ip_table(self) -> List[Tuple[str, ip_address]]:
        """ Return list of tuples (sid, ip) """
        with self._lock:
            res = [
                (sid, ip_desc.ip) for sid, ip_desc in
                self._store.sid_ips_map.items()
            ]
            return res

    def get_ip_for_sid(self, sid: str) -> Optional[ip_address]:
        """ if ip is mapped to sid, return it, else return None """
        if not self._store.sid_ips_map.get(sid, None):
            return None
        else:
            return self._store.sid_ips_map[sid].ip

    def get_sid_for_ip(self, requested_ip: ip_address) -> Optional[str]:
        """ If ip is associated with an sid, return the sid, else None """
        for sid, ip_desc in self._store.sid_ips_map.items():
            if requested_ip == ip_desc.ip:
                return sid
        return None

    def is_ip_in_state(self, ip_addr: ip_address, state: IPState):
        """
            Check if IP address is on a given state
        """
        ip_state_map = self._store.ip_state_map if ip_addr.version == 4 \
            else self._store.ipv6_state_map
        return ip_state_map.test_ip_state(ip_addr, state)

    def release_ip_address(
        self, sid: str, ip: ip_address,
        version: int = IPAddress.IPV4,
    ):
        """ Release an IP address.

        A released IP is moved to a released list. Released IPs are recycled
        periodically to the free list. SID IP mappings are removed at the
        recycling time.

        Args:
            sid (string): universal subscriber id
            ip (ipaddress.ip_address): IP address to release
            version (int): version of IP address to release

        Raises:
            MappingNotFoundError: if the given sid-ip mapping is not found
            IPNotInUseError: if the given IP is not found in the used list
        """
        with self._lock:
            if not (
                sid in self._store.sid_ips_map and ip
                == self._store.sid_ips_map[sid].ip
            ):
                logging.error(
                    "Releasing unknown <SID, IP> pair: <%s, %s> "
                    "sid_ips_map[%s]: %s",
                    sid, ip, sid, self._store.sid_ips_map.get(sid),
                )
                raise MappingNotFoundError(
                    "(%s, %s) pair is not found", sid, str(ip),
                )
            if not self.is_ip_in_state(ip, IPState.ALLOCATED):
                logging.error(
                    "IP not found in used list, check if IP is "
                    "already released: <%s, %s>", sid, ip,
                )
                raise IPNotInUseError("IP not found in used list: %s", str(ip))

            IP_RELEASED_TOTAL.inc()

            if version == IPAddress.IPV4:
                self._store.ip_state_map.mark_ip_state(ip, IPState.RELEASED)
                self._try_set_recycle_timer()  # start the timer to recycle
            elif version == IPAddress.IPV6:
                # For IPv6, no recycling logic
                ip_desc = self._store.ipv6_state_map.mark_ip_state(
                    ip,
                    IPState.FREE,
                )
                self.ipv6_allocator.release_ip(ip_desc)
                del self._store.sid_ips_map[ip_desc.sid]

    def list_gateway_info(self) -> List[GWInfo]:
        with self._lock:
            return self._store.dhcp_gw_info.get_all_router_ips()

    def set_gateway_info(self, info: GWInfo):
        ip = str(ipaddress.ip_address(info.ip.address))
        with self._lock:
            self._store.dhcp_gw_info.update_mac(ip, info.mac, info.vlan)

    def _recycle_reaped_ips(self):
        """ Periodically called to recycle the given IPs

        *** It is highly not recommended to call this function directly, even
        in tests. ***

        Recycling depends on the period, T = self._recycling_interval_seconds,
        which is set at construction time.

        """
        with self._lock:
            for ip in self._store.ip_state_map.list_ips(IPState.REAPED):
                ip_desc = self._store.ip_state_map.mark_ip_state(
                    ip,
                    IPState.FREE,
                )
                logging.debug("Release Reaped IP: %s", ip_desc)

                self.ip_allocator.release_ip(ip_desc)
                # update SID-IP map
                del self._store.sid_ips_map[ip_desc.sid]

            # Set timer for the next round of recycling
            self._recycle_timer = None
            if not self._store.ip_state_map.is_ip_state_map_empty(
                    IPState.RELEASED,
            ):
                self._try_set_recycle_timer()

    def _try_set_recycle_timer(self):
        """ Try set the recycle timer and move RELEASED IPs to the REAPED state

        self._try_set_recycle_timer is called in two places:
            1) at the end of self.release_ip_address, we are guaranteed that
            some IPs exist in RELEASED state, so we attempt to initiate a timer
            then.
            2) at the end of self._recycle_reaped_ips, the call to
            self._try_set_recycle_timer serves as a callback for setting the
            next timer, if any IPs have been released since the current timer
            was initiated.
        """
        with self._lock:
            # check if auto recycling is enabled and no timer has been set
            if self._recycling_interval_seconds is not None \
                    and not self._recycle_timer:
                for ip in self._store.ip_state_map.list_ips(IPState.RELEASED):
                    self._store.ip_state_map.mark_ip_state(ip, IPState.REAPED)
                if self._recycling_interval_seconds:
                    self._recycle_timer = threading.Timer(
                        self._recycling_interval_seconds,
                        self._recycle_reaped_ips,
                    )
                    self._recycle_timer.start()
                else:
                    self._recycle_reaped_ips()

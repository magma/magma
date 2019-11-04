"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.


The IP allocator maintains the life cycle of assigned IP addresses.

The IP allocator accepts IP blocks (range of IP addresses), and supports
allocating and releasing IP addresses from the assigned IP blocks. Note
that an IP address is not immediately made available for allocation right
after release: it is "reserved" for the same client for a certain period of
time to ensure that 1) an observer, e.g. pipelined, that caches IP states has
enough time to pull the updated IP states; 2) IP packets intended for the
old client will not be unintentionally routed to a new client until the old
TCP connection expires.

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

from __future__ import absolute_import
from __future__ import division
from __future__ import print_function
from __future__ import unicode_literals

import logging
import threading
from collections import defaultdict
from copy import deepcopy
from ipaddress import ip_address

import redis

import magma.mobilityd.serialize_utils as serialize_utils
from magma.common.redis.containers import RedisHashDict, RedisSet
from magma.mobilityd.ip_descriptor import IPDesc, IPState
from magma.mobilityd.metrics import (IP_ALLOCATED_TOTAL, IP_RELEASED_TOTAL)


DEFAULT_IP_RECYCLE_INTERVAL = 15


class IPAllocator():
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
            - self._ip_states: {state=>{ip=>ip_desc}}
            - self._sid_ips_map: {SID=>[IPDesc]}

        The utilized redis_containers store a cache of state in local memory,
        so reads are the same speed as without persistence. For writes, state
        is serialized with provided serializers and written through to the
        Redis process. Used as expected, writes will be small (a few bytes).
        Redis's performance can be measured with the redis-benchmark tool,
        but we expect almost 100% of writes to take less than 1 millisecond.

    TODO (ouyangj):
        1) use persistent store for internal states. #13576041
            1.1) handle the case when the PGW restarts/crashes: mobilityd needs
            to reset the set of allocated ips. The pool configuration
            should persist across pGW crashes.
            1.2) handle the case when mobilityd restarts/crashes: mobilityd
            needs to reconstruct the states of assigned ip blocks, as well as
            allocated/released IPs.
    """

    def __init__(self,
                 *,
                 recycling_interval=DEFAULT_IP_RECYCLE_INTERVAL,
                 persist_to_redis=True,
                 redis_port=6379):
        """ Initializes a new IP allocator

        Args:
            recycling_interval (number): minimum time, in seconds, before a
                released IP is recycled and freed into the pool of available
                IPs.

                Default: None, no recycling will occur automatically.
            persist_to_redis (bool): store all state in local process if falsy,
                else write state to Redis service
        """
        logging.debug('Persist to Redis: %s', persist_to_redis)
        self._lock = threading.RLock()  # re-entrant locks

        self._recycle_timer = None  # reference to recycle timer
        self._recycling_interval_seconds = recycling_interval

        if not persist_to_redis:
            self._assigned_ip_blocks = set()  # {ip_block}
            self._ip_states = defaultdict(dict)  # {state=>{ip=>ip_desc}}
            self._sid_ips_map = defaultdict(list)  # {SID=>[IPDesc]}
        else:
            if not redis_port:
                raise ValueError(
                    'Must specify a redis_port in mobilityd config.')
            client = redis.Redis(host='localhost', port=redis_port)
            self._assigned_ip_blocks = RedisSet(
                client,
                'mobilityd:assigned_ip_blocks',
                serialize_utils.serialize_ip_block,
                serialize_utils.deserialize_ip_block)
            self._ip_states = defaultdict_key(
                lambda key: cleared_ip_states(client, key))
            self._sid_ips_map = RedisHashDict(
                client,
                'mobilityd:sid_to_descriptors',
                serialize_utils.serialize_ip_descs,
                serialize_utils.deserialize_ip_descs,
                default_factory=list)

        # NOTE: clearing below (and _ip_states above) is done to sync
        # a mobilityd restart with restarts from other services. Once
        # all mobilityd-related services store persistent state (e.g.
        # using the Redis service), we can remove these sync-necessitated
        # clears on initialization.

        # TODO: once we are no longer clearing all values on initialization,
        # we should add unit tests to ensure values are written to Redis
        # correctly.
        self._assigned_ip_blocks.clear()
        self._sid_ips_map.clear()

    def add_ip_block(self, ipblock):
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

        with self._lock:
            self._assigned_ip_blocks.add(ipblock)
            # TODO(oramadan) t23793559 HACK reserve the GW address for gtp_br0 iface and test VM
            num_reserved_addresses = 11
            for ip in ipblock.hosts():
                state = IPState.RESERVED if num_reserved_addresses > 0 else IPState.FREE
                ip_desc = IPDesc(ip=ip, state=state,
                                 ip_block=ipblock, sid=None)
                self._add_ip_to_state(ip, ip_desc, state)
                if num_reserved_addresses > 0:
                    num_reserved_addresses -= 1


    def remove_ip_blocks(self, *ipblocks, force=False):
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

        Returns a set of the blocks that have been successfully removed.
        """

        with self._lock:
            remove_blocks = set(ipblocks) & self._assigned_ip_blocks

            extraneous_blocks = set(ipblocks) ^ remove_blocks
            # check unknown ip blocks
            if extraneous_blocks:
                logging.warning("Cannot remove unknown IP block(s): %s",
                                extraneous_blocks)
            del extraneous_blocks

            # "soft" removal does not remove blocks have IPs allocated
            if not force:
                allocated_ip_block_set = self._get_allocated_ip_block_set()
                remove_blocks -= allocated_ip_block_set
                del allocated_ip_block_set

            # Remove the associated IP addresses
            remove_ips = \
                (ip for block in remove_blocks for ip in block.hosts())
            for ip in remove_ips:
                self._remove_ip_from_state(ip, IPState.FREE)
                self._remove_ip_from_state(ip, IPState.RELEASED)
                self._remove_ip_from_state(ip, IPState.REAPED)
                if force:
                    self._remove_ip_from_state(ip, IPState.ALLOCATED)
                else:
                    assert not self._test_ip_state(ip, IPState.ALLOCATED), \
                    "Unexpected ALLOCATED IP %s from a soft IP block removal"

            # Remove the IP blocks
            self._assigned_ip_blocks -= remove_blocks

            # Clean up SID maps
            remove_addresses = \
                (ip for block in remove_blocks for ip in block.hosts())
            for address in remove_addresses:
                for sid in self._sid_ips_map:
                    if address in self._sid_ips_map[sid]:
                        self._sid_ips_map[sid].remove(address)

            # Can't use generators here
            remove_sids = tuple(sid for sid in self._sid_ips_map
                                if not self._sid_ips_map[sid])
            for sid in remove_sids:
                self._sid_ips_map.pop(sid)

            return remove_blocks

    def list_added_ip_blocks(self):
        """ List IP blocks added to the IP allocator

        Return:
             copy of the list of assigned IP blocks
        """
        with self._lock:
            ip_blocks = list(deepcopy(self._assigned_ip_blocks))
        return ip_blocks

    def list_allocated_ips(self, ipblock):
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

        with self._lock:
            res = [ip for ip in ipblock \
                   if self._test_ip_state(ip, IPState.ALLOCATED)]
        return res

    def alloc_ip_address(self, sid):
        """ Allocate an IP address from the free list

        Assumption: one IP per UE and one-to-one mappings between SID and IP.
        At most one IP address can be allocated for a UE at any given time.

        Args:
            sid (string): universal subscriber id

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
            if sid in self._sid_ips_map:
                old_ip = self._sid_ips_map[sid][0]
                if self._test_ip_state(old_ip, IPState.ALLOCATED):
                    # MME state went out of sync with mobilityd!
                    # Recover gracefully by allocating the same IP
                    logging.warn("Re-allocate IP %s for sid %s without "
                                 "MME releasing it first", old_ip, sid)
                    # TODO: enable strict checking after root causing the
                    # issue in MME
                    # raise DuplicatedIPAllocationError(
                    #     "An IP has been allocated for this IMSI")
                elif self._test_ip_state(old_ip, IPState.RELEASED):
                    ip_desc = self._mark_ip_state(old_ip, IPState.ALLOCATED)
                    ip_desc.sid = sid
                    logging.debug("SID %s IP %s RELEASED => ALLOCATED",
                                  sid, old_ip)
                elif self._test_ip_state(old_ip, IPState.REAPED):
                    ip_desc = self._mark_ip_state(old_ip, IPState.ALLOCATED)
                    ip_desc.sid = sid
                    logging.debug("SID %s IP %s REAPED => ALLOCATED",
                                  sid, old_ip)
                else:
                    raise AssertionError("Unexpected internal state")
                logging.info("Allocating the same IP %s for sid %s",
                             old_ip, sid)

                IP_ALLOCATED_TOTAL.inc()
                return old_ip

            # if an IP is not yet allocated for the UE, allocate a new IP
            if self._get_ip_count(IPState.FREE):
                ip_desc = self._pop_ip_from_state(IPState.FREE)
                ip_desc.sid = sid
                ip_desc.state = IPState.ALLOCATED
                self._add_ip_to_state(ip_desc.ip, ip_desc, IPState.ALLOCATED)
                self._sid_ips_map[sid].append(ip_desc.ip)
                assert len(self._sid_ips_map[sid]) == 1, \
                        "Only one IP per SID is supported"

                IP_ALLOCATED_TOTAL.inc()
                return ip_desc.ip
            else:
                logging.error("Run out of available IP addresses")
                raise NoAvailableIPError("No available IP addresses")

    def get_sid_ip_table(self):
        """ Return list of tuples (sid, ip) """
        res = []
        with self._lock:
            res = [(sid, ip) for sid, ips in self._sid_ips_map.items()
                   for ip in ips]
        return res

    def get_ip_for_sid(self, sid):
        """ if ip is mapped to sid, return it, else return None """
        with self._lock:
            if sid in self._sid_ips_map:
                if not self._sid_ips_map[sid]:
                    raise AssertionError("Unexpected internal state")
                else:
                    return self._sid_ips_map[sid][0]
            return None

    def get_sid_for_ip(self, requested_ip):
        """ If ip is associated with an sid, return the sid, else None """
        with self._lock:
            for sid, ips in self._sid_ips_map.items():
                if requested_ip in ips:
                    return sid
            return None

    def release_ip_address(self, sid, ip):
        """ Release an IP address.

        A released IP is moved to a released list. Released IPs are recycled
        periodically to the free list. SID IP mappings are removed at the
        recycling time.

        Args:
            sid (string): universal subscriber id
            ip (ipaddress.ip_address): IP address to release

        Raises:
            MappingNotFoundError: if the given sid-ip mapping is not found
            IPNotInUseError: if the given IP is not found in the used list
        """
        with self._lock:
            if not (sid in self._sid_ips_map and ip in self._sid_ips_map[sid]):
                logging.error(
                    "Releasing unknown <SID, IP> pair: <%s, %s>", sid, ip)
                raise MappingNotFoundError(
                    "(%s, %s) pair is not found", sid, str(ip))
            if not self._test_ip_state(ip, IPState.ALLOCATED):
                logging.error("IP not found in used list, check if IP is "
                              "already released: <%s, %s>", sid, ip)
                raise IPNotInUseError("IP not found in used list: %s", str(ip))

            self._mark_ip_state(ip, IPState.RELEASED)
            IP_RELEASED_TOTAL.inc()

            self._try_set_recycle_timer()  # start the timer to recycle

    def reset(self):
        """ Nuke all internal states

            Handle the case when only PGW crashes/restarts
        """
        raise NotImplementedError()

    def reload(self):
        """ Reconstruct states from the persistent data

            Handle the case when only Mobilityd crashes/restarts
        """
        raise NotImplementedError()

    def _recycle_reaped_ips(self):
        """ Periodically called to recycle the given IPs

        *** It is highly not recommended to call this function directly, even
        in tests. ***

        Recycling depends on the period, T = self._recycling_interval_seconds,
        which is set at construction time.

        """
        with self._lock:
            for ip in self._list_ips(IPState.REAPED):
                ip_desc = self._mark_ip_state(ip, IPState.FREE)
                sid = ip_desc.sid

                # update SID-IP map
                ips = self._sid_ips_map[sid]
                ips.remove(ip)
                if not ips:
                    del self._sid_ips_map[sid]

                ip_desc.sid = None

            # Set timer for the next round of recycling
            self._recycle_timer = None
            if self._get_ip_count(IPState.RELEASED):
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
                for ip in self._list_ips(IPState.RELEASED):
                    self._mark_ip_state(ip, IPState.REAPED)
                if self._recycling_interval_seconds:
                    self._recycle_timer = threading.Timer(
                        self._recycling_interval_seconds,
                        self._recycle_reaped_ips)
                    self._recycle_timer.start()
                else:
                    self._recycle_reaped_ips()

    def _add_ip_to_state(self, ip, ip_desc, state):
        """ Add ip=>ip_desc pairs to a internal dict """
        assert ip_desc.state == state, \
                "ip_desc.state %s does not match with state %s" \
                % (ip_desc.state, state)
        assert state in IPState, "unknown state %s" % state

        with self._lock:
            self._ip_states[state][ip.exploded] = ip_desc

    def _remove_ip_from_state(self, ip, state):
        """ Remove an IP from a internal dict """
        assert state in IPState, "unknown state %s" % state

        with self._lock:
            ip_desc = self._ip_states[state].pop(ip.exploded, None)
        return ip_desc

    def _pop_ip_from_state(self, state):
        """ Pop an IP from a internal dict """
        assert state in IPState, "unknown state %s" % state

        with self._lock:
            _, ip_desc = self._ip_states[state].popitem()
        return ip_desc

    def _get_ip_count(self, state):
        """ Return number of IPs in a state """
        assert state in IPState, "unknown state %s" % state

        with self._lock:
            return len(self._ip_states[state])

    def _test_ip_state(self, ip, state):
        """ check if IP is in state X """
        assert state in IPState, "unknown state %s" % state

        with self._lock:
            return ip.exploded in self._ip_states[state]

    def _get_ip_state(self, ip):
        """ return the state of an IP """
        for state in IPState:
            if self._test_ip_state(ip, state):
                return state
        raise AssertionError("IP %s not found in any states" % ip)

    def _list_ips(self, state):
        """ return a list of IPs in state X """
        assert state in IPState, "unknown state %s" % state

        with self._lock:
            return [ip_address(ip) for ip in self._ip_states[state]]

    def _mark_ip_state(self, ip, state):
        """ Remove, mark, add: move IP to a new state """
        assert state in IPState, "unknown state %s" % state

        old_state = self._get_ip_state(ip)
        with self._lock:
            ip_desc = self._ip_states[old_state][ip.exploded]

            # some internal checks
            assert ip_desc.state != state, \
                    "move IP to the same state %s" % state
            assert ip == ip_desc.ip, "Unmatching ip_desc for %s" % ip
            if ip_desc.state == IPState.FREE:
                assert ip_desc.sid is None, "Unexpected sid in a freed IPDesc"
            else:
                assert ip_desc.sid is not None, \
                        "Missing sid in state %s IPDesc" % state

            # remove, mark, add
            self._remove_ip_from_state(ip, old_state)
            ip_desc.state = state
            self._add_ip_to_state(ip, ip_desc, state)
        return ip_desc

    def _get_allocated_ip_block_set(self):
        """ A IP block is allocated if ANY IP is allocated from it """
        with self._lock:
            allocated_ips = self._ip_states[IPState.ALLOCATED]
        return set([ip_desc.ip_block for ip_desc in allocated_ips.values()])


class OverlappedIPBlocksError(Exception):
    """ Exception thrown when a given IP block overlaps with existing ones
    """
    pass


class IPBlockNotFoundError(Exception):
    """ Exception thrown when listing an IP block that is not found in the ip
    block list
    """
    pass


class NoAvailableIPError(Exception):
    """ Exception thrown when no IP is available in the free list for an ip
    allocation request
    """
    pass


class DuplicatedIPAllocationError(Exception):
    """ Exception thrown when an IP has already been allocated to a UE """
    pass


class IPNotInUseError(Exception):
    """ Exception thrown when releasing an IP address that is not found in the
    used list
    """
    pass


class MappingNotFoundError(Exception):
    """ Exception thrown when releasing a non-exising SID-IP mapping """
    pass


class SubscriberNotFoundError(Exception):
    """ Exception thrown when subscriber ID is not found in SID-IP mapping """


def cleared_ip_states(client, key):
    """ Get cleared Redis view of IP states. """
    redis_dict = RedisHashDict(
        client,
        'mobilityd:ip_to_descriptor:{}'.format(key),
        serialize_utils.serialize_ip_desc,
        serialize_utils.deserialize_ip_desc)
    redis_dict.clear()
    return redis_dict


class defaultdict_key(defaultdict):
    """
    Same as standard lib's defaultdict, but takes the key as a parameter.
    """
    def __missing__(self, key):
        # Follow defaultdict pattern in raising KeyError
        if self.default_factory is None:
            raise KeyError(key)
        self[key] = self.default_factory(key)
        return self[key]

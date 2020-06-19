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

from __future__ import absolute_import, division, print_function, \
    unicode_literals

from collections import defaultdict
from ipaddress import ip_address, ip_network
from typing import List, Set

import redis
from magma.mobilityd import mobility_store as store
from magma.mobilityd.ip_descriptor import IPDesc, IPState
from random import choice

DEFAULT_IP_RECYCLE_INTERVAL = 15


class IpDescriptorMap:

    def __init__(self,
                 persist_to_redis: bool = True,
                 redis_port: int = 6379):
        """

        Args:
            persist_to_redis (bool): store all state in local process if falsy,
                else write state to Redis service
            redis_port (int): redis server port number.
        """
        if not persist_to_redis:
            self.ip_states = defaultdict(dict)  # {state=>{ip=>ip_desc}}
        else:
            if not redis_port:
                raise ValueError(
                    'Must specify a redis_port in mobilityd config.')
            client = redis.Redis(host='localhost', port=redis_port)
            self.ip_states = store.defaultdict_key(
                lambda key: store.ip_states(client, key))

    def add_ip_to_state(self, ip: ip_address, ip_desc: IPDesc,
                         state: IPState):
        """ Add ip=>ip_desc pairs to a internal dict """
        assert ip_desc.state == state, \
            "ip_desc.state %s does not match with state %s" \
            % (ip_desc.state, state)
        assert state in IPState, "unknown state %s" % state

        self.ip_states[state][ip.exploded] = ip_desc

    def remove_ip_from_state(self, ip: ip_address, state: IPState) -> IPDesc:
        """ Remove an IP from a internal dict """
        assert state in IPState, "unknown state %s" % state

        ip_desc = self.ip_states[state].pop(ip.exploded, None)
        return ip_desc

    def pop_ip_from_state(self, state: IPState) -> IPDesc:
        """ Pop an IP from a internal dict """
        assert state in IPState, "unknown state %s" % state

        ip_state_key = choice(list(self.ip_states[state].keys()))
        ip_desc = self.ip_states[state].pop(ip_state_key)
        return ip_desc

    def get_ip_count(self, state: IPState) -> int:
        """ Return number of IPs in a state """
        assert state in IPState, "unknown state %s" % state

        return len(self.ip_states[state])

    def test_ip_state(self, ip: ip_address, state: IPState) -> bool:
        """ check if IP is in state X """
        assert state in IPState, "unknown state %s" % state

        return ip.exploded in self.ip_states[state]

    def get_ip_state(self, ip: ip_address) -> IPState:
        """ return the state of an IP """
        for state in IPState:
            if self.test_ip_state(ip, state):
                return state
        raise AssertionError("IP %s not found in any states" % ip)

    def list_ips(self, state: IPState) -> List[ip_address]:
        """ return a list of IPs in state X """
        assert state in IPState, "unknown state %s" % state

        return [ip_address(ip) for ip in self.ip_states[state]]

    def mark_ip_state(self, ip: ip_address, state: IPState) -> IPDesc:
        """ Remove, mark, add: move IP to a new state """
        assert state in IPState, "unknown state %s" % state

        old_state = self.get_ip_state(ip)
        ip_desc = self.ip_states[old_state][ip.exploded]

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
        self.remove_ip_from_state(ip, old_state)
        ip_desc.state = state
        self.add_ip_to_state(ip, ip_desc, state)
        return ip_desc

    def get_allocated_ip_block_set(self) -> Set[ip_network]:
        """ A IP block is allocated if ANY IP is allocated from it """
        allocated_ips = self.ip_states[IPState.ALLOCATED]
        return {ip_desc.ip_block for ip_desc in allocated_ips.values()}
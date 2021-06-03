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

from __future__ import (
    absolute_import,
    division,
    print_function,
    unicode_literals,
)

from ipaddress import ip_address, ip_network
from typing import Dict, List, MutableMapping, Optional, Set

from magma.mobilityd.ip_descriptor import IPDesc, IPState

DEFAULT_IP_RECYCLE_INTERVAL = 15


class IpDescriptorMap:

    def __init__(self, ip_states: MutableMapping[IPState, Dict[str, IPDesc]]):
        """

        Args:
            ip_states: Dictionary containing IPDesc keyed by current state
        """
        self.ip_states = ip_states

    def add_ip_to_state(
        self, ip: ip_address, ip_desc: IPDesc,
        state: IPState,
    ):
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

    def pop_ip_from_state(self, state: IPState) -> Optional[IPDesc]:
        """ Pop an IP from a internal dict """
        assert state in IPState, "unknown state %s" % state

        try:
            _, ip_desc = self.ip_states[state].popitem()
            return ip_desc
        except KeyError:
            return None

    def is_ip_state_map_empty(self, state: IPState) -> bool:
        """
        Args:
            state: IP State to check for

        Returns: True if IPs map is empty for a given state, else otherwise
        """
        assert state in IPState, "unknown state %s" % state

        return bool(self.ip_states[state]) == False

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
            assert ip_desc.sid is None, \
                "Unexpected sid in a freed IPDesc {}".format(ip_desc)
        else:
            assert ip_desc.sid is not None, \
                "Missing sid in state %s IPDesc {}".format(ip_desc)

        # remove, mark, add
        self.remove_ip_from_state(ip, old_state)
        ip_desc.state = state
        self.add_ip_to_state(ip, ip_desc, state)
        return ip_desc

    def get_allocated_ip_block_set(self) -> Set[ip_network]:
        """ A IP block is allocated if ANY IP is allocated from it """
        allocated_ips = self.ip_states[IPState.ALLOCATED]
        return {ip_desc.ip_block for ip_desc in allocated_ips.values()}

    def __str__(self) -> str:
        """ return the state of an IP """
        ret_str = "{}:".format(self.__class__.__name__)
        for state in IPState:
            ret_str = ret_str + "\n{}".format(state)
            for _ip, ip_desc in self.ip_states[state].items():
                ret_str = ret_str + "\n{}".format(str(ip_desc))
        return ret_str

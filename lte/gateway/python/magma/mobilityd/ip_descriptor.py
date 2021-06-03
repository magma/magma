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
import ipaddress
from enum import Enum


class IPState(Enum):
    FREE = 1
    ALLOCATED = 2
    RELEASED = 3
    REAPED = 4
    RESERVED = 5


class IPType(Enum):
    STATIC = 1
    IP_POOL = 2
    DHCP = 3


class IPv6SessionAllocType(str, Enum):
    RANDOM = "RANDOM"
    HASH = "HASH"


class IPDesc:
    """
    IP descriptor.

    Properties:
        ip (ipaddress.ip_address)
        state (IPState)
        sid (str)
        ip_block (ipaddress.ip_network)
        type (IPType)
        vlan_id (int)
    """

    def __init__(
        self, ip: ipaddress.ip_address = None, state: IPState = None,
        sid: str = None, ip_block: ipaddress.ip_network = None,
        ip_type: IPType = None, vlan_id: int = 0,
    ):
        self.ip = ip
        self.ip_block = ip_block
        self.state = state
        self.sid = sid
        self.type = ip_type
        self.vlan_id = 0
        if 0 < vlan_id < 4096:
            self.vlan_id = vlan_id

    def __str__(self):
        as_str = '<mobilityd.IPDesc ' + \
                 '{{ip: {}, ip_block: {}, state: {}, sid: {}, type: {}'.format(
                     self.ip,
                     self.ip_block,
                     self.state,
                     self.sid,
                     self.type,
                 )

        if self.vlan_id != 0:
            as_str = as_str + " vlan_is: {}".format(self.vlan_id)

        as_str = as_str + " }}>"
        return as_str

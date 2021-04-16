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
from datetime import datetime, timedelta
from enum import IntEnum
from typing import Optional

from .mac import MacAddress


# map the enum to actual protocol values.
class DHCPState(IntEnum):
    """
    DHCP protocol states.
    """
    UNKNOWN = 0
    DISCOVER = 1
    OFFER = 2
    REQUEST = 3
    DECLINE = 4
    ACK = 5
    NAK = 6
    RELEASE = 7
    FORCE_RENEW = 8


class DHCPDescriptor:
    def __init__(
        self, mac: MacAddress, ip: str,
        state_requested: DHCPState, vlan: str,
        state: DHCPState = DHCPState.UNKNOWN,
        subnet: str = None, server_ip: str = None,
        router_ip: str = None, lease_expiration_time: int = 0,
        xid: str = None,
    ):
        """
        DHCP descriptor. This object maintains all information for
        given DHCP protocol transactions.

        Args:
            mac: Mac address of request
            ip: Allocated IP if IP is assigned by DHCP server
            state: Last known protocol state on server
            state_requested: Last requested state by client
            vlan: track vlan id if MAC is allocated in a VLAN
            subnet: subnet of IP from DHCP offer or ACK
            server_ip: DHCP server IP address
            router_ip: GW IP address
            lease_expiration_time: DHCP lease time.
            xid: XID used in DHCP protocol.
        """
        self.mac = mac
        self.ip = ip
        self.vlan = vlan
        self.subnet = subnet
        self.state = state
        self.state_requested = state_requested
        self.server_ip = server_ip
        self.xid = xid
        self.lease_expiration_time = datetime.now(
        ) + timedelta(seconds=lease_expiration_time)
        self.router_ip = router_ip
        if self.state == DHCPState.ACK:
            new_deadline = datetime.now() + timedelta(seconds=(lease_expiration_time / 2))
            self.lease_renew_deadline = new_deadline
        else:
            self.lease_renew_deadline = datetime.now()

    def __str__(self):
        return "state: {:8s} requested state: {:8s} mac {} ip {} subnet {} " \
               "DHCP: {} router {} " \
               "lease time {}, renew {} xid {} vlan {}" \
            .format(
                self.state.name, self.state_requested.name,
                str(self.mac), self.ip, self.subnet,
                self.server_ip, self.router_ip, self.lease_expiration_time,
                self.lease_renew_deadline, self.xid, self.vlan,
            )

    def get_ip_address(self) -> Optional[str]:
        """
        Return valid IP address as per DHCP state.
        :return: IP address
        """
        if self.state == DHCPState.OFFER or self.state == DHCPState.REQUEST \
                or self.state == DHCPState.ACK:
            return self.ip

    def ip_is_allocated(self) -> bool:
        """
        Check if validate IP address is allocated by DHCP server.

        :return: True or False
        """
        return (
            self.state == DHCPState.OFFER
            or self.state == DHCPState.REQUEST
            or self.state == DHCPState.ACK
        )

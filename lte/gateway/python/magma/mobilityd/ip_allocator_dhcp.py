"""
Copyright 2020 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

Allocates IP address as per DHCP server in the uplink network.
"""

from __future__ import (
    absolute_import,
    division,
    print_function,
    unicode_literals,
)

import json
import logging
import subprocess
import threading
from copy import deepcopy
from datetime import datetime
from ipaddress import IPv4Network, ip_address, ip_network
from json import JSONDecodeError
from os import environ
from threading import Condition
from typing import List, Optional, cast

from magma.mobilityd.ip_descriptor import IPDesc, IPState, IPType

from .dhcp_desc import DHCPDescriptor, DHCPState
from .ip_allocator_base import IPAllocator, NoAvailableIPError
from .mac import MacAddress, create_mac_from_sid
from .mobility_store import MobilityStore
from .utils import IPAddress, IPNetwork

DEFAULT_DHCP_REQUEST_RETRY_FREQUENCY = 10
DEFAULT_DHCP_REQUEST_RETRY_DELAY = 1
LEASE_RENEW_WAIT_MIN = 200

DHCP_CLI_HELPER_PATH= f"{environ.get('MAGMA_ROOT')}/lte/gateway/python/scripts/dhcp_helper_cli.py"
LOG = logging.getLogger('mobilityd.dhcp.alloc')

DHCP_ACTIVE_STATES = [DHCPState.ACK, DHCPState.OFFER]

class IPAllocatorDHCP(IPAllocator):
    def __init__(
        self, store: MobilityStore, retry_limit: int = 300,
        iface: str = "eth2", #TODO read this from config file
    ):
        """
        Allocate IP address for SID using DHCP server.
        SID is mapped to MAC address using function defined in mac.py
        then this mac address used in DHCP request to allocate new IP
        from DHCP server.
        This IP is also cached to improve performance in case of
        reallocation for same SID in short period of time.

        Args:
            store: Moblityd storage instance
            retry_limit: try DHCP request
            iface: DHCP interface.
        """
        self._store = store
        self.dhcp_wait = Condition()
        self._retry_limit = retry_limit  # default wait for two minutes
        self._iface = iface
        self._monitor_thread = threading.Thread(
            target=self._monitor_dhcp_state,
        )

    def _monitor_dhcp_state(self):
        """
        monitor DHCP client state.
        """
        while True:
            wait_time = LEASE_RENEW_WAIT_MIN
            with self.dhcp_wait:
                dhcp_desc: DHCPDescriptor

                for dhcp_desc in self._store.dhcp_store.values():
                    logging.debug("monitor: %s", dhcp_desc)
                    # Only process active records.
                    if dhcp_desc.state not in DHCP_ACTIVE_STATES:
                        continue

                    now = datetime.now()
                    logging.debug("monitor time: %s", now)
                    request_state = DHCPState.REQUEST
                    # in case of lost DHCP lease rediscover it.
                    if now >= dhcp_desc.lease_expiration_time:
                        logging.debug("sending lease renewal")
                        dhcp_cli_response = subprocess.run([
                            DHCP_CLI_HELPER_PATH,
                            "--mac", str(dhcp_desc.mac),
                            "--vlan", str(dhcp_desc.vlan),
                            "--interface", self._iface,
                            "allocate"],
                            capture_output=True
                        )

                        if dhcp_cli_response.returncode != 0:
                            logging.error(f"Could not decode '{dhcp_cli_response.stdout}' received '{dhcp_cli_response.stderr}' from {DHCP_CLI_HELPER_PATH} called with parameters '{dhcp_cli_response.args}'")
                            raise NoAvailableIPError(f'Failed to call dhcp_helper_cli.')
                    if now >= dhcp_desc.lease_renew_deadline:
                        logging.debug("sending lease renewal")
                        dhcp_cli_response = subprocess.run([
                            DHCP_CLI_HELPER_PATH,
                            "--mac", str(dhcp_desc.mac),
                            "--vlan", str(dhcp_desc.vlan),
                            "--interface", self._iface,
                            "renew"],
                            "--ip", str(dhcp_desc.ip),
                            "--server-ip", str(dhcp_desc.server_ip),
                            capture_output=True
                        )

                        if dhcp_cli_response.returncode != 0:
                            logging.error(f"Could not decode '{dhcp_cli_response.stdout}' received '{dhcp_cli_response.stderr}' from {DHCP_CLI_HELPER_PATH} called with parameters '{dhcp_cli_response.args}'")
                            raise NoAvailableIPError(f'Failed to call dhcp_helper_cli.')
                    else:
                        # Find next renewal wait time.
                        time_to_renew = dhcp_desc.lease_renew_deadline - now
                        wait_time = min(
                            wait_time, time_to_renew.total_seconds(),
                        )

            # default in wait is 30 sec
            wait_time = max(wait_time, self._lease_renew_wait_min)
            logging.debug("lease renewal check after: %s sec", wait_time)
            self._monitor_thread_event.wait(wait_time)
            if self._monitor_thread_event.is_set():
                break

    def add_ip_block(self, ipblock: IPNetwork):
        logging.warning(
            "No need to allocate block for DHCP allocator: %s",
            ipblock,
        )

    def remove_ip_blocks(
        self,
        ipblocks: List[IPNetwork],
        force: bool = False,
    ) -> List[IPNetwork]:
        logging.warning(
            "Trying to delete ipblock from DHCP allocator: %s",
            ipblocks,
        )
        return []

    def list_added_ip_blocks(self) -> List[IPNetwork]:
        return list(deepcopy(self._store.assigned_ip_blocks))

    def list_allocated_ips(self, ipblock: IPNetwork) -> List[IPAddress]:
        """ List IP addresses allocated from a given IP block

        Args:
            ipblock (ipaddress.ip_network): ip network to add
            e.g. ipaddress.ip_network("10.0.0.0/24")

        Return:
            list of IP addresses (ipaddress.ip_address)

        """
        return [
            ip for ip in
            self._store.ip_state_map.list_ips(IPState.ALLOCATED)
            if ip in ipblock
        ]

    def get_dhcp_desc_from_store(
        self, mac: MacAddress,
        vlan: int,
    ) -> Optional[DHCPDescriptor]:
        """
                Get DHCP description for given MAC.
        Args:
            mac: Mac address of the client
            vlan: vlan id if the IP allocated in a VLAN

        Returns: Current DHCP info.
        """
        key = mac.as_redis_key(vlan)
        if key in self._store.dhcp_store:
            return self._store.dhcp_store[key]

        LOG.debug("lookup error for %s", str(key))
        return None

    def alloc_ip_address(self, sid: str, vlan: int) -> IPDesc:
        """
        Assumption: one-to-one mappings between SID and IP.

        Args:
            sid (string): universal subscriber id
            vlan: vlan of the APN

        Returns:
            ipaddress.ip_address: IP address allocated

        Raises:
            NoAvailableIPError: if run out of available IP addresses
        """
        mac = create_mac_from_sid(sid)

        dhcp_desc = self.get_dhcp_desc_from_store(mac, vlan)
        LOG.debug(
            "allocate IP for %s mac %s dhcp_desc %s", sid, mac,
            dhcp_desc,
        )

        if not dhcp_desc or not dhcp_allocated_ip(dhcp_desc):
            dhcp_response = subprocess.run([
                DHCP_CLI_HELPER_PATH,
                "--mac", str(mac),
                "--vlan", str(vlan),
                "--interface", self._iface,
                "--json",
                "allocate"],
                capture_output=True
            )

            if dhcp_response.returncode != 0:
                logging.error(f"Could not decode '{dhcp_response.stdout}' received '{dhcp_response.stderr}' from {DHCP_CLI_HELPER_PATH} called with parameters '{dhcp_response.args}'")
                raise NoAvailableIPError(f'Failed to call dhcp_helper_cli.')

            try:
                dhcp_json = json.loads(dhcp_response.stdout)
            except JSONDecodeError as e:
                logging.error(f"Could not decode '{dhcp_response.stdout}' received '{dhcp_response.stderr}' from dhcp_helper_cli called with parameters '{dhcp_response.args}'")
                raise NoAvailableIPError(f'Failed to json parse message returned from dhcp_helper_cli.')

            if dhcp_json:
                dhcp_desc = DHCPDescriptor(
                    mac=mac,
                    ip=ip_address(dhcp_json["ip"]),
                    vlan=vlan,
                    state_requested=DHCPState.ACK,
                    state=DHCPState.ACK,
                    subnet=IPv4Network(dhcp_json["subnet"], strict=False),
                    server_ip=ip_address(dhcp_json["server_ip"]) if dhcp_json["server_ip"] else None,
                    router_ip=ip_address(dhcp_json["server_ip"]) if dhcp_json["server_ip"] else None,  # TODO extract router ip from dhcp pkt
                    lease_expiration_time=int(dhcp_json["lease_expiration_time"]),
                )

                with self.dhcp_wait:
                    self._store.dhcp_store[mac.as_redis_key(vlan)] = dhcp_desc

        if dhcp_desc and dhcp_desc.ip and dhcp_desc.subnet:
            ip_block = ip_network(dhcp_desc.subnet)
            ip_desc = IPDesc(
                ip=ip_address(dhcp_desc.ip),
                state=IPState.ALLOCATED,
                sid=sid,
                ip_block=ip_block,
                ip_type=IPType.DHCP,
                vlan_id=vlan,
            )
            self._store.assigned_ip_blocks.add(ip_block)
            logging.info("returning IPDesc")
            return ip_desc
        else:
            msg = f"No available IP addresses From DHCP for SID: {sid} MAC {mac}"
            raise NoAvailableIPError(msg)

    def release_ip(self, ip_desc: IPDesc):
        """
        Release IP address, this involves following steps.
        1. send DHCP protocol packet to release the IP.
        2. update IP block list.
        3. update IP from ip-state.

        Args:
            ip_desc: release needs following info from IPDesc.
                SID used to get mac address, IP assigned to this SID,
                IP block of the IP address, vlan id of the APN.
        Returns: None
        """
        mac = create_mac_from_sid(ip_desc.sid)
        vlan = ip_desc.vlan_id
        dhcp_desc = self.get_dhcp_desc_from_store(mac, vlan)

        if dhcp_desc:
            dhcp_cli_response = subprocess.run([
                DHCP_CLI_HELPER_PATH,
                "--mac", str(mac),
                "--vlan", str(vlan),
                "--interface", self._iface,
                "release"],
                "--ip", str(ip_desc.ip),
                "--server-ip", str(dhcp_desc.server_ip),
                capture_output=True
            )

            if dhcp_cli_response.returncode != 0:
                logging.error(f"Could not decode '{dhcp_cli_response.stdout}' received '{dhcp_cli_response.stderr}' from {DHCP_CLI_HELPER_PATH} called with parameters '{dhcp_cli_response.args}'")
                raise NoAvailableIPError(f'Failed to call dhcp_helper_cli.')

            with self.dhcp_wait:
                key = mac.as_redis_key(vlan)
                del self._store.dhcp_store[key]
        else:
            LOG.error("Unallocated DHCP release for MAC: %s", mac)


        # Remove the IP from free IP list, since DHCP is the
        # owner of this IP
        self._store.ip_state_map.remove_ip_from_state(ip_desc.ip, IPState.FREE)

        list_allocated_ips = self._store.ip_state_map.list_ips(
            IPState.ALLOCATED,
        )
        for ipaddr in list_allocated_ips:
            if ipaddr in ip_desc.ip_block:
                # found the IP, do not remove this ip_block
                return

        ip_block_network = ip_network(ip_desc.ip_block)
        if ip_block_network in self._store.assigned_ip_blocks:
            self._store.assigned_ip_blocks.remove(ip_block_network)
        logging.debug(
            "del: _assigned_ip_blocks %s ipblock %s",
            self._store.assigned_ip_blocks, ip_desc.ip_block,
        )

def dhcp_allocated_ip(dhcp_desc: DHCPDescriptor) -> bool:
    return dhcp_desc.ip_is_allocated()

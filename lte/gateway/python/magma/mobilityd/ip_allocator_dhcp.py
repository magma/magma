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
import tempfile
import threading
from copy import deepcopy
from datetime import datetime
from ipaddress import IPv4Network, ip_address, ip_network
from threading import Condition
from typing import Any, Dict, List, Optional

from magma.mobilityd.ip_descriptor import IPDesc, IPState, IPType

from .dhcp_desc import DHCPDescriptor, DHCPState
from .ip_allocator_base import IPAllocator, NoAvailableIPError
from .mac import MacAddress, create_mac_from_sid
from .mobility_store import MobilityStore
from .utils import IPAddress, IPNetwork

DEFAULT_DHCP_REQUEST_RETRY_FREQUENCY = 10
DEFAULT_DHCP_REQUEST_RETRY_DELAY = 1
LEASE_RENEW_WAIT_MIN = 200

DHCP_HELPER_CLI = "dhcp_helper_cli.py"
LOG = logging.getLogger('mobilityd.dhcp.alloc')

DHCP_ACTIVE_STATES = [DHCPState.ACK, DHCPState.OFFER]


class IPAllocatorDHCP(IPAllocator):

    MAX_DHCP_PROCS: int = 3
    dhcp_helper_procs: List[subprocess.Popen] = []

    def __init__(
        self, store: MobilityStore, retry_limit: int = 300, start: bool = True,
            iface: str = "eth2", lease_renew_wait_min: float = LEASE_RENEW_WAIT_MIN,  # TODO read this from config file
    ) -> None:
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
        self._lease_renew_wait_min = lease_renew_wait_min
        self._monitor_thread: Optional[threading.Thread] = threading.Thread(
            target=self._monitor_dhcp_state,
            daemon=True,
        )
        self._monitor_thread_event = threading.Event()
        if start:
            self.start_monitor_thread()

    def start_monitor_thread(self) -> None:
        if self._monitor_thread is None:
            self._monitor_thread = threading.Thread(
                target=self._monitor_dhcp_state,
                daemon=True,
            )
        self._monitor_thread.start()

    def stop_monitor_thread(self, join: bool = False, reset: bool = False) -> None:
        self._monitor_thread_event.set()
        if join and self._monitor_thread:
            self._monitor_thread.join()
        if reset:
            self._monitor_thread = None
            self._monitor_thread_event.clear()

    def _monitor_dhcp_state(self) -> None:
        """
        monitor DHCP client state.
        """
        while True:
            wait_time = self._lease_renew_wait_min
            with self.dhcp_wait:
                for _, dhcp_desc in self._store.dhcp_store.items():
                    logging.debug("monitor: %s", dhcp_desc)
                    # Only process active records.
                    if dhcp_desc.state not in DHCP_ACTIVE_STATES:
                        continue
                    now = datetime.now()
                    logging.debug("monitor time: %s", now)

                    if now >= dhcp_desc.lease_expiration_time:
                        logging.debug("sending lease allocate")
                        with tempfile.NamedTemporaryFile() as tmpfile:
                            call_args = [
                                DHCP_HELPER_CLI,
                                "--mac", str(dhcp_desc.mac),
                                "--vlan", str(dhcp_desc.vlan),
                                "--interface", self._iface,
                                "--save-file", tmpfile.name,
                                "--json",
                                "allocate",
                            ]
                            dhcp_cli_response = self._get_dhcp_helper_cli_response(call_args, tmpfile.name)
                        self._parse_dhcp_helper_cli_response_to_store(
                            dhcp_desc=dhcp_desc, dhcp_response=dhcp_cli_response,
                            mac=dhcp_desc.mac, vlan=dhcp_desc.vlan,
                        )
                    elif now >= dhcp_desc.lease_renew_deadline:
                        logging.debug("sending lease renewal")
                        with tempfile.NamedTemporaryFile() as tmpfile:
                            call_args = [
                                DHCP_HELPER_CLI,
                                "--mac", str(dhcp_desc.mac),
                                "--vlan", str(dhcp_desc.vlan),
                                "--interface", self._iface,
                                "--save-file", tmpfile.name,
                                "--json",
                                "renew",
                                "--ip", str(dhcp_desc.ip),
                                "--server-ip", str(dhcp_desc.server_ip),
                            ]
                            dhcp_cli_response = self._get_dhcp_helper_cli_response(call_args, tmpfile.name)
                        self._parse_dhcp_helper_cli_response_to_store(
                            dhcp_desc=dhcp_desc, dhcp_response=dhcp_cli_response,
                            mac=dhcp_desc.mac, vlan=dhcp_desc.vlan,
                        )
                    else:
                        time_to_renew = dhcp_desc.lease_renew_deadline - now
                        wait_time = min(
                            wait_time, time_to_renew.total_seconds(),
                        )

            # default in wait is 30 sec
            logging.debug("lease renewal check after: %s sec", wait_time)
            self._monitor_thread_event.wait(wait_time)
            if self._monitor_thread_event.is_set():
                break

    def add_ip_block(self, ipblock: IPNetwork) -> None:
        logging.warning(
            "No need to allocate block for DHCP allocator: %s",
            ipblock,
        )

    def remove_ip_blocks(
        self, ipblocks: List[IPNetwork],
        force: bool = False,
    ) -> List[IPNetwork]:
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

        remove_blocks = set(ipblocks) & self._store.assigned_ip_blocks
        logging.debug(
            "Current assigned IP blocks: %s",
            self._store.assigned_ip_blocks,
        )
        logging.debug("IP blocks to remove: %s", ipblocks)

        extraneous_blocks = set(ipblocks) ^ remove_blocks
        # check unknown ip blocks
        if extraneous_blocks:
            logging.warning(
                "Cannot remove unknown IP block(s): %s",
                extraneous_blocks,
            )
        del extraneous_blocks

        # "soft" removal does not remove blocks which have IPs allocated
        if not force:
            allocated_ip_block_set = self._store.ip_state_map.get_allocated_ip_block_set()
            remove_blocks -= allocated_ip_block_set
            del allocated_ip_block_set

        # Remove the associated IP addresses
        remove_ips = (ip for block in remove_blocks for ip in block.hosts())
        for ip in remove_ips:
            for state in (IPState.FREE, IPState.RELEASED, IPState.REAPED):
                ip_desc = self._store.ip_state_map.remove_ip_from_state(ip, state)
                if ip_desc:
                    self._release_dhcp_ip(ip_desc)
            if force:
                ip_desc = self._store.ip_state_map.remove_ip_from_state(
                    ip,
                    IPState.ALLOCATED,
                )
                if ip_desc:
                    self._release_dhcp_ip(ip_desc)
            else:
                assert not self._store.ip_state_map.test_ip_state(
                    ip,
                    IPState.ALLOCATED,
                ), \
                    "Unexpected ALLOCATED IP %s from a soft IP block " \
                    "removal "

            # Clean up SID maps
            for sid in list(self._store.sid_ips_map):
                self._store.sid_ips_map.pop(sid)

        # Remove the IP blocks
        self._store.assigned_ip_blocks -= remove_blocks

        # Can't use generators here
        remove_sids = tuple(
            sid for sid in self._store.sid_ips_map
            if not self._store.sid_ips_map[sid]
        )
        for sid in remove_sids:
            self._store.sid_ips_map.pop(sid)

        for block in remove_blocks:
            logging.info('Removed IP block %s from IPv4 address pool', block)
        return list(remove_blocks)

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

    def alloc_ip_address(self, sid: str, vlan_id: int) -> IPDesc:
        """
        Assumption: one-to-one mappings between SID and IP.

        Args:
            sid (string): universal subscriber id
            vlan_id: vlan of the APN

        Returns:
            ipaddress.ip_address: IP address allocated

        Raises:
            NoAvailableIPError: if run out of available IP addresses
        """
        mac = create_mac_from_sid(sid)

        dhcp_desc = self.get_dhcp_desc_from_store(mac, vlan_id)
        LOG.debug(
            "allocate IP for %s mac %s dhcp_desc %s", sid, mac,
            dhcp_desc,
        )

        if not dhcp_desc or not dhcp_allocated_ip(dhcp_desc):
            with tempfile.NamedTemporaryFile() as tmpfile:
                call_args = [
                    DHCP_HELPER_CLI,
                    "--mac", str(mac),
                    "--vlan", str(vlan_id),
                    "--interface", self._iface,
                    "--save-file", tmpfile.name,
                    "--json",
                    "allocate",
                ]
                dhcp_response = self._get_dhcp_helper_cli_response(call_args, tmpfile.name)
            with self.dhcp_wait:
                dhcp_desc = self._parse_dhcp_helper_cli_response_to_store(
                    dhcp_desc=dhcp_desc, dhcp_response=dhcp_response,
                    mac=mac, vlan=vlan_id,
                )

        if dhcp_desc and dhcp_desc.ip and dhcp_desc.subnet:
            ip_block = ip_network(dhcp_desc.subnet)
            ip_desc = IPDesc(
                ip=ip_address(dhcp_desc.ip),
                state=IPState.ALLOCATED,
                sid=sid,
                ip_block=ip_block,
                ip_type=IPType.DHCP,
                vlan_id=vlan_id,
            )
            self._store.assigned_ip_blocks.add(ip_block)
            self._store.ip_state_map.add_ip_to_state(
                ip=ip_address(dhcp_desc.ip),
                ip_desc=ip_desc,
                state=IPState.ALLOCATED,
            )
            return ip_desc
        else:
            msg = f"No available IP addresses From DHCP for SID: {sid} MAC {mac}"
            raise NoAvailableIPError(msg)

    def _parse_dhcp_helper_cli_response_to_store(
            self, dhcp_desc: Optional[DHCPDescriptor], dhcp_response: Dict[str, Any],
            mac: MacAddress, vlan: int,
    ) -> Optional[DHCPDescriptor]:
        if dhcp_response:
            dhcp_desc = DHCPDescriptor(
                mac=mac,
                ip=dhcp_response["ip"],
                vlan=vlan,
                state_requested=DHCPState.ACK,
                state=DHCPState.ACK,
                subnet=str(IPv4Network(dhcp_response["subnet"], strict=False)),
                server_ip=ip_address(dhcp_response["server_ip"]) if dhcp_response["server_ip"] else None,
                router_ip=ip_address(dhcp_response["router_ip"]) if dhcp_response["router_ip"] else None,
                lease_expiration_time=int(dhcp_response["lease_expiration_time"]),
            )

            self._store.dhcp_store[mac.as_redis_key(vlan)] = dhcp_desc
            self._store.dhcp_gw_info.update_ip(dhcp_desc.router_ip, vlan)
        return dhcp_desc

    @staticmethod
    def _get_dhcp_helper_cli_response(call_args: List[str], save_file: str) -> Dict[str, Any]:
        ret = subprocess.run(
            call_args,
            capture_output=True,
            check=False,
        )
        if ret.returncode != 0:
            call_str = " ".join(call_args)
            logging.error(
                "CLI call '%s' failed with return code %s and error %s",
                call_str, ret.returncode, ret.stderr,
            )
            raise NoAvailableIPError('Failed to call dhcp_helper_cli.')

        with open(save_file, 'r', encoding="utf-8") as f:
            dhcp_response = json.load(f)
        return dhcp_response

    def release_ip(self, ip_desc: IPDesc) -> None:
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
        self._release_dhcp_ip(ip_desc)

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

    def _release_dhcp_ip(self, ip_desc: IPDesc) -> None:
        logging.info("Releasing: %s", ip_desc)
        mac = create_mac_from_sid(ip_desc.sid)
        vlan = ip_desc.vlan_id
        dhcp_desc = self.get_dhcp_desc_from_store(mac, vlan)
        logging.info("Releasing dhcp desc: %s", dhcp_desc)
        if dhcp_desc:
            self._deque_old_process(ip_desc)
            proc = subprocess.Popen(
                [
                    DHCP_HELPER_CLI,
                    "--mac", str(mac),
                    "--vlan", str(vlan),
                    "--interface", self._iface,
                    "release",
                    "--ip", str(ip_desc.ip),
                    "--server-ip", str(dhcp_desc.server_ip),
                ],
            )
            self.dhcp_helper_procs.insert(0, proc)
            key = mac.as_redis_key(vlan)
            with self.dhcp_wait:
                del self._store.dhcp_store[key]
        else:
            LOG.error("Unallocated DHCP release for MAC: %s", mac)

    def _deque_old_process(self, ip_desc: IPDesc) -> None:
        while len(self.dhcp_helper_procs) >= self.MAX_DHCP_PROCS:
            oldest_proc = self.dhcp_helper_procs.pop()
            if oldest_proc.poll() is None:
                try:
                    oldest_proc.wait(timeout=10)
                except subprocess.TimeoutExpired:
                    logging.warning(
                        "Unable to release IP %s."
                        "Killing release process "
                        "and moving on.", ip_desc.ip,
                    )
                    oldest_proc.kill()


def dhcp_allocated_ip(dhcp_desc: DHCPDescriptor) -> bool:
    return dhcp_desc.ip_is_allocated()

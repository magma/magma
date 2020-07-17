"""
Copyright (c) 2020-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.


Allocates IP address as per DHCP server in the uplink network.

"""

from __future__ import absolute_import, division, print_function, \
    unicode_literals

import logging
from copy import deepcopy
from ipaddress import ip_address, ip_network
from typing import List, Set, MutableMapping
from threading import Condition

from magma.mobilityd.ip_descriptor import IPState, IPDesc

from .ip_descriptor_map import IpDescriptorMap
from .ip_allocator_base import IPAllocator, NoAvailableIPError
from .dhcp_client import DHCPClient
from .mac import MacAddress, create_mac_from_sid
from .dhcp_desc import DHCPState, DHCPDescriptor
from .uplink_gw import UplinkGatewayInfo

DEFAULT_DHCP_REQUEST_RETRY_FREQUENCY = 10
DEFAULT_DHCP_REQUEST_RETRY_DELAY = 1

LOG = logging.getLogger('mobilityd.dhcp.alloc')


class IPAllocatorDHCP(IPAllocator):
    def __init__(self,
                 assigned_ip_blocks: Set[ip_network],
                 ip_state_map: IpDescriptorMap,
                 dhcp_store: MutableMapping[str, DHCPDescriptor],
                 gw_info: UplinkGatewayInfo,
                 retry_limit: int = 300,
                 iface: str = "dhcp0"):
        """
        Allocate IP address for SID using DHCP server.
        SID is mapped to MAC address using function defined in mac.py
        then this mac address used in DHCP request to allocate new IP
        from DHCP server.
        This IP is also cached to improve performance in case of
        reallocation for same SID in short period of time.

        Args:
            assigned_ip_blocks: set of IP blocks, populated from DHCP.
            ip_state_map: maintains state of IP allocation to UE.
            dhcp_store: maintains DHCP transaction for each active MAC address
            gw_info_map: maintains uplink GW info
            retry_limit: try DHCP request
            iface: DHCP interface.
        """
        self._ip_state_map = ip_state_map  # {state=>{ip=>ip_desc}}
        self._assigned_ip_blocks = assigned_ip_blocks
        self.dhcp_wait = Condition()
        self._dhcp_client = DHCPClient(dhcp_wait=self.dhcp_wait,
                                       dhcp_store=dhcp_store,
                                       gw_info=gw_info,
                                       iface=iface)
        self._retry_limit = retry_limit  # default wait for two minutes
        self._dhcp_client.run()

    def add_ip_block(self, ipblock: ip_network):
        logging.warning("No need to allocate block for DHCP allocator: %s", ipblock)

    def remove_ip_blocks(self, *ipblocks: List[ip_network],
                         _force: bool = False) -> List[ip_network]:
        logging.warning("trying to delete ipblock from DHCP allocator: %s", ipblocks)
        return []

    def list_added_ip_blocks(self) -> List[ip_network]:
        return list(deepcopy(self._assigned_ip_blocks))

    def list_allocated_ips(self, ipblock: ip_network) -> List[ip_address]:
        """ List IP addresses allocated from a given IP block

        Args:
            ipblock (ipaddress.ip_network): ip network to add
            e.g. ipaddress.ip_network("10.0.0.0/24")

        Return:
            list of IP addresses (ipaddress.ip_address)

        """
        return [ip for ip in self._ip_state_map.list_ips(IPState.ALLOCATED)
                if ip in ipblock]

    def alloc_ip_address(self, sid: str) -> IPDesc:
        """
        Assumption: one-to-one mappings between SID and IP.

        Args:
            sid (string): universal subscriber id

        Returns:
            ipaddress.ip_address: IP address allocated

        Raises:
            NoAvailableIPError: if run out of available IP addresses
        """
        mac = create_mac_from_sid(sid)
        LOG.debug("allocate IP for %s mac %s", sid, mac)

        dhcp_desc = self._dhcp_client.get_dhcp_desc(mac)
        LOG.debug("got IP from redis: %s", dhcp_desc)

        if dhcp_allocated_ip(dhcp_desc) is not True:
            dhcp_desc = self._alloc_ip_address_from_dhcp(mac)

        if dhcp_allocated_ip(dhcp_desc):
            ip_block = ip_network(dhcp_desc.subnet)
            ip_desc = IPDesc(ip_address(dhcp_desc.ip), IPState.ALLOCATED, sid, ip_block)
            LOG.debug("Got IP after sending DHCP requests: %s", ip_desc)
            self._assigned_ip_blocks.add(ip_block)

            return ip_desc
        else:
            raise NoAvailableIPError("No available IP addresses From DHCP")

    def release_ip(self, sid: str, deleted_ip: ip_address, ip_block: ip_network):
        """
        Release IP address, this involves following steps.
        1. send DHCP protocol packet to release the IP.
        2. update IP block list.
        3. update IP from ip-state.

        Args:
            sid: SID, used to get mac address.
            deleted_ip: IP assigned to this SID
            ip_block: IP block of the IP address.

        Returns: None
        """
        self._dhcp_client.release_ip_address(create_mac_from_sid(sid))
        # Remove the IP from free IP list, since DHCP is the
        # owner of this IP
        self._ip_state_map.remove_ip_from_state(deleted_ip, IPState.FREE)

        list_allocated_ips = self._ip_state_map.list_ips(IPState.ALLOCATED)
        for ipaddr in list_allocated_ips:
            if ipaddr in ip_block:
                # found the IP, do not remove this ip_block
                return

        ip_block_network = ip_network(ip_block)
        if ip_block_network in self._assigned_ip_blocks:
            self._assigned_ip_blocks.remove(ip_block_network)
        logging.debug("del: _assigned_ip_blocks %s ipblock %s",
                      self._assigned_ip_blocks, ip_block)

    def stop_dhcp_sniffer(self):
        self._dhcp_client.stop()

    def _alloc_ip_address_from_dhcp(self, mac: MacAddress) -> DHCPDescriptor:
        retry_count = 0
        with self.dhcp_wait:
            dhcp_desc = None
            while (retry_count < self._retry_limit and
                   dhcp_allocated_ip(dhcp_desc) is not True):

                if retry_count % DEFAULT_DHCP_REQUEST_RETRY_FREQUENCY == 0:
                    self._dhcp_client.send_dhcp_packet(mac, DHCPState.DISCOVER)
                self.dhcp_wait.wait(timeout=DEFAULT_DHCP_REQUEST_RETRY_DELAY)

                dhcp_desc = self._dhcp_client.get_dhcp_desc(mac)

                retry_count = retry_count + 1

            return dhcp_desc


def dhcp_allocated_ip(dhcp_desc) -> bool:
    return dhcp_desc is not None and dhcp_desc.ip_is_allocated()

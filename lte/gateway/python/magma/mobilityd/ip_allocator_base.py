"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.

"""

from __future__ import absolute_import, division, print_function, \
    unicode_literals

from abc import ABC, abstractmethod

from ipaddress import ip_address, ip_network
from typing import List
from enum import Enum

from magma.mobilityd.ip_descriptor import IPDesc

DEFAULT_IP_RECYCLE_INTERVAL = 15


class IPAllocator(ABC):

    @abstractmethod
    def add_ip_block(self, ipblock: ip_network):
        ...

    @abstractmethod
    def remove_ip_blocks(self, *ipblocks: List[ip_network],
                         force: bool) -> List[ip_network]:
        ...

    @abstractmethod
    def list_added_ip_blocks(self) -> List[ip_network]:
        ...

    @abstractmethod
    def list_allocated_ips(self, ipblock: ip_network) -> List[ip_address]:
        ...

    @abstractmethod
    def alloc_ip_address(self, sid: str) -> IPDesc:
        ...

    @abstractmethod
    def release_ip(self, sid):
        ...

class IPAllocatorType(Enum):
    IP_POOL = 1


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
    """ Exception thrown when an IP has already been allocated to a UE
    """
    pass
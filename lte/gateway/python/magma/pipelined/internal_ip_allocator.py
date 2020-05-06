"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""
import ipaddress
import threading
from typing import List


class InternalIPAllocator:
    """
    Internal IP Allocator

    This class assigns internal IP from the same subnet as the bridge to the UE,
    this is used primarily for redirection to catch ARPs as they'll be directed
    to the bridge(where we auto respond to them)
    """

    def __init__(self):
        self.internal_ip_network = ipaddress.ip_network('192.168.0.0/16')
        self.internal_ip_iterator = self.internal_ip_network.hosts()
        self._lock = threading.Lock()  # write lock

    def next_ip(self, invalid_ips: List[str] = None):
        with self._lock:
            ip = next(self.internal_ip_iterator, None)
            while str(ip) in invalid_ips:
                ip = next(self.internal_ip_iterator, None)
                if ip is None:
                    self.internal_ip_iterator = self.internal_ip_network.hosts()
            return ip

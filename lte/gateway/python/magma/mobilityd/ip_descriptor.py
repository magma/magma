"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""
import ipaddress
from enum import Enum


class IPState(Enum):
    FREE = 1
    ALLOCATED = 2
    RELEASED = 3
    REAPED = 4
    RESERVED = 5


class IPDesc():
    """
    IP descriptor.

    Properties:
        ip (ipaddress.ip_address)
        state (IPState)
        sid (str)
        ip_block (ipaddress.ip_network)
    """

    def __init__(self, ip: ipaddress.ip_address = None, state: IPState = None,
                 sid: str = None, ip_block: ipaddress.ip_network = None):
        self.ip = ip
        self.ip_block = ip_block
        self.state = state
        self.sid = sid

    def __str__(self):
        as_str = '<mobilityd.IPDesc ' + \
                 '{{ip: {}, ip_block: {}, state: {}, sid: {}}}>'.format(
                     self.ip,
                     self.ip_block,
                     self.state,
                     self.sid)
        return as_str

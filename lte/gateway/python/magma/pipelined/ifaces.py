"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""
import asyncio
import netifaces

from magma.pipelined.metrics import NETWORK_IFACE_STATUS

POLL_INTERVAL_SECONDS = 3


@asyncio.coroutine
def monitor_ifaces(iface_names, loop):
    """
    Call to poll the network interfaces and set the corresponding metric
    """
    while True:
        active = set(netifaces.interfaces())
        for iface in iface_names:
            status = 1 if iface in active else 0
            NETWORK_IFACE_STATUS.labels(iface_name=iface).set(status)
        yield from asyncio.sleep(POLL_INTERVAL_SECONDS, loop=loop)

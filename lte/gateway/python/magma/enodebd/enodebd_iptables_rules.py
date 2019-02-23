#!/usr/bin/env python3

"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""


import asyncio
import logging
import shlex

from magma.common.misc_utils import IpPreference, get_ip_from_if
from magma.configuration.service_configs import load_service_config

IPTABLES_RULE_FMT = """sudo iptables -t nat
    -{add} PREROUTING
    -d {public_ip}
    -p tcp
    --dport {port}
    -j DNAT --to-destination {private_ip}"""


def get_iptables_rule(port, enodebd_public_ip, private_ip, add=True):
    return IPTABLES_RULE_FMT.format(
        add='A' if add else 'D',
        public_ip=enodebd_public_ip,
        port=port,
        private_ip=private_ip,
    )


async def run(cmd):
    """Fork shell and run command NOTE: Popen is non-blocking"""
    cmd = shlex.split(cmd)
    proc = await asyncio.create_subprocess_shell(" ".join(cmd))
    await proc.communicate()
    if proc.returncode != 0:
        # This can happen because the NAT prerouting rule didn't exist
        logging.info('Possible error running async subprocess: %s exited with '
                     'return code [%d].', cmd, proc.returncode)
    return proc.returncode


@asyncio.coroutine
def set_enodebd_iptables_rule():
    """
    Remove & Set iptable rules for exposing public IP
    for enobeb instead of private IP..
    """
    # Remove & Set iptable rules for exposing public ip
    # for enobeb instead of private
    cfg = load_service_config('enodebd')
    port, interface = cfg['tr069']['port'], cfg['tr069']['interface']
    enodebd_public_ip = cfg['tr069']['public_ip']
    # IPv4 only as iptables only works for IPv4. TODO: Investigate ip6tables?
    enodebd_ip = get_ip_from_if(interface, preference=IpPreference.IPV4_ONLY)
    # Incoming data from 192.88.99.142 -> enodebd address (eg 192.168.60.142)
    yield from run(get_iptables_rule(
        port, enodebd_public_ip, enodebd_ip, add=False))
    yield from run(get_iptables_rule(
        port, enodebd_public_ip, enodebd_ip, add=True))


if __name__ == '__main__':
    set_enodebd_iptables_rule()

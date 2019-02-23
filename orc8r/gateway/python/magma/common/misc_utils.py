"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""

import asyncio
from enum import Enum
import ipaddress
import os

import netifaces


class IpPreference(Enum):
    IPV4_ONLY = 1
    IPV4_PREFERRED = 2
    IPV6_PREFERRED = 3
    IPV6_ONLY = 4


def get_if_ip_with_netmask(interface, preference=IpPreference.IPV4_PREFERRED):
    """
    Get IP address and netmask (in form /255.255.255.0)
    from interface name and return as tuple (ip, netmask).
    Note: If multiple v4/v6 addresses exist, the first is chosen

    Raise ValueError if unable to get requested IP address.
    """
    # Raises ValueError if interface is unavailable
    ip_addresses = netifaces.ifaddresses(interface)

    try:
        ipv4_address = (ip_addresses[netifaces.AF_INET][0]['addr'],
                        ip_addresses[netifaces.AF_INET][0]['netmask'])
    except KeyError:
        ipv4_address = None

    try:
        ipv6_address = (ip_addresses[netifaces.AF_INET6][0]['addr'],
                        ip_addresses[netifaces.AF_INET6][0]['netmask'])
    except KeyError:
        ipv6_address = None

    if preference == IpPreference.IPV4_ONLY:
        if ipv4_address is not None:
            return ipv4_address
        else:
            raise ValueError('Error getting IPv4 address for %s' % interface)

    elif preference == IpPreference.IPV4_PREFERRED:
        if ipv4_address is not None:
            return ipv4_address
        elif ipv6_address is not None:
            return ipv6_address
        else:
            raise ValueError('Error getting IPv4/6 address for %s' % interface)

    elif preference == IpPreference.IPV6_PREFERRED:
        if ipv6_address is not None:
            return ipv6_address
        elif ipv4_address is not None:
            return ipv4_address
        else:
            raise ValueError('Error getting IPv6/4 address for %s' % interface)

    elif preference == IpPreference.IPV6_ONLY:
        if ipv6_address is not None:
            return ipv6_address
        else:
            raise ValueError('Error getting IPv6 address for %s' % interface)

    else:
        raise ValueError('Unknown IP preference %s' % preference)


def get_ip_from_if(iface_name, preference=IpPreference.IPV4_PREFERRED):
    """
    Get ip address from interface name and return as string.
    Extract only ip address from (ip, netmask)
    """
    return get_if_ip_with_netmask(iface_name, preference)[0]


def get_ip_from_if_cidr(iface_name, preference=IpPreference.IPV4_PREFERRED):
    """
    Get IPAddress with netmask from interface name and
    transform into CIDR (eth1 -> 192.168.60.142/24)
    notation return as string.
    """
    ip, netmask = get_if_ip_with_netmask(iface_name, preference)
    ip = '%s/%s' % (ip, netmask)
    interface = ipaddress.ip_interface(ip).with_prefixlen  # Set CIDR notation
    return interface


def cidr_to_ip_netmask_tuple(cidr_network):
    """
    Convert CIDR-format IP network string (e.g. 10.0.0.1/24) to a tuple
    (ip, netmask) where netmask is in the form (n.n.n.n).

    Args:
        cidr_network (str): IPv4 network in CIDR notation

    Returns:
        (str, str): 2-tuple of IP address and netmask
    """
    network = ipaddress.ip_network(cidr_network)
    return '{}'.format(network.network_address), '{}'.format(network.netmask)


def call_process(cmd, callback, loop):
    loop = loop or asyncio.get_event_loop()
    loop.create_task(loop.subprocess_shell(
        lambda: SubprocessProtocol(callback), "nohup " + cmd,
        preexec_fn=os.setsid))


class SubprocessProtocol(asyncio.SubprocessProtocol):
    def __init__(self, callback):
        self._callback = callback
        self._transport = None

    def connection_made(self, transport):
        self._transport = transport

    def process_exited(self):
        self._callback(self._transport.get_returncode())

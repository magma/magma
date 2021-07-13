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

import asyncio
import ipaddress
import os
from enum import Enum

import netifaces
import snowflake


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
        ipv4_address = (
            ip_addresses[netifaces.AF_INET][0]['addr'],
            ip_addresses[netifaces.AF_INET][0]['netmask'],
        )
    except KeyError:
        ipv4_address = None

    try:
        ipv6_address = (
            ip_addresses[netifaces.AF_INET6][0]["addr"].split("%")[0],
            ip_addresses[netifaces.AF_INET6][0]["netmask"],
        )

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


def get_all_if_ips_with_netmask(
    interface,
    preference=IpPreference.IPV4_PREFERRED,
):
    """
    Get all IP addresses and netmasks (in form /255.255.255.0)
    from interface name and return as a list of tuple (ip, netmask).

    Raise ValueError if unable to get requested IP addresses.
    """
    # Raises ValueError if interface is unavailable
    ip_addresses = netifaces.ifaddresses(interface)

    try:
        ipv4_addresses = [(ip_address['addr'], ip_address['netmask']) for
                          ip_address in ip_addresses[netifaces.AF_INET]]
    except KeyError:
        ipv4_addresses = None

    try:
        ipv6_addresses = [(ip_address['addr'], ip_address['netmask']) for
                          ip_address in ip_addresses[netifaces.AF_INET6]]
    except KeyError:
        ipv6_addresses = None

    if preference == IpPreference.IPV4_ONLY:
        if ipv4_addresses is not None:
            return ipv4_addresses
        else:
            raise ValueError('Error getting IPv4 addresses for %s' % interface)

    elif preference == IpPreference.IPV4_PREFERRED:
        if ipv4_addresses is not None:
            return ipv4_addresses
        elif ipv6_addresses is not None:
            return ipv6_addresses
        else:
            raise ValueError(
                'Error getting IPv4/6 addresses for %s' % interface,
            )

    elif preference == IpPreference.IPV6_PREFERRED:
        if ipv6_addresses is not None:
            return ipv6_addresses
        elif ipv4_addresses is not None:
            return ipv4_addresses
        else:
            raise ValueError(
                'Error getting IPv6/4 addresses for %s' % interface,
            )

    elif preference == IpPreference.IPV6_ONLY:
        if ipv6_addresses is not None:
            return ipv6_addresses
        else:
            raise ValueError('Error getting IPv6 addresses for %s' % interface)

    else:
        raise ValueError('Unknown IP preference %s' % preference)


def get_ip_from_if(iface_name, preference=IpPreference.IPV4_PREFERRED):
    """
    Get ip address from interface name and return as string.
    Extract only ip address from (ip, netmask)
    """
    return get_if_ip_with_netmask(iface_name, preference)[0]


def get_all_ips_from_if(iface_name, preference=IpPreference.IPV4_PREFERRED):
    """
    Get all ip addresses from interface name and return as a list of string.
    Extract only ip address from (ip, netmask)
    """
    return [
        ip[0] for ip in
        get_all_if_ips_with_netmask(iface_name, preference)
    ]


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


def get_all_ips_from_if_cidr(
    iface_name,
    preference=IpPreference.IPV4_PREFERRED,
):
    """
    Get all IPAddresses with netmask from interface name and
    transform into CIDR (eth1 -> 192.168.60.142/24) notation
    return as a list of string.
    """

    def ip_cidr_gen():
        for ip, netmask in get_all_if_ips_with_netmask(iface_name, preference):
            ip = '%s/%s' % (ip, netmask)
            # Set CIDR notation
            ip_cidr = ipaddress.ip_interface(ip).with_prefixlen
            yield ip_cidr

    return [ip_cidr for ip_cidr in ip_cidr_gen()]


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


def get_if_mac_address(interface):
    """
    Returns the MAC address of an interface.
    Note: If multiple MAC addresses exist, the first one is chosen.

    Raise ValueError if unable to get requested IP address.
    """
    addr = netifaces.ifaddresses(interface)
    try:
        return addr[netifaces.AF_LINK][0]['addr']
    except KeyError:
        raise ValueError('Error getting MAC address for %s' % interface)


def get_gateway_hwid() -> str:
    """
    Returns the HWID of the gateway
    Note: Currently this uses the snowflake at /etc/snowflake
    """
    return snowflake.snowflake()


def is_interface_up(interface):
    """
    Returns whether an interface is up.
    """
    try:
        addr = netifaces.ifaddresses(interface)
    except ValueError:
        return False
    return netifaces.AF_INET in addr


def call_process(cmd, callback, loop):
    loop = loop or asyncio.get_event_loop()
    loop.create_task(
        loop.subprocess_shell(
        lambda: SubprocessProtocol(callback), "nohup " + cmd,
        preexec_fn=os.setsid,
        ),
    )


class SubprocessProtocol(asyncio.SubprocessProtocol):
    def __init__(self, callback):
        self._callback = callback
        self._transport = None

    def connection_made(self, transport):
        self._transport = transport

    def process_exited(self):
        self._callback(self._transport.get_returncode())

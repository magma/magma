#!/usr/bin/env python3

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

"""
CLI for dhclient
"""


import argparse
import ipaddress
import random
import sys
from ipaddress import ip_address, ip_network

from magma.mobilityd import mobility_store as store
from magma.mobilityd.dhcp_desc import DHCPDescriptor, DHCPState
from magma.mobilityd.mac import MacAddress
from magma.mobilityd.uplink_gw import UplinkGatewayInfo


class DhcpClientCLI:
    def __init__(self):
        """
        Inspect and manipulate DHCP client state.
        """
        self.dhcp_client_state = store.MacToIP()  # mac => DHCP_State

    def list_all_record(self, args):
        """
                Print list DHCP and GW records.
        Args:
            args: None

        Returns: None
        """

        for k, v in self.dhcp_client_state.items():
            print("mac: %s DHCP state: %s" % (k.replace('_', ':'), str(v)))

        gw_info = UplinkGatewayInfo(store.GatewayInfoMap())
        print("Default GW: %s" % gw_info.get_gw_ip())
        print("GW Mac address: %s" % gw_info.get_gw_mac())

    def add_record(self, args):
        """
        Add DHCP record.
        Args:
            args: All data required for DHCP state.

        Returns:

        """

        state = DHCPState(args.state)
        ipaddr = ip_address(args.ip)
        subnet = ip_network(args.subnet, strict=False)
        dhcp_ip = ip_address(args.dhcp)
        desc = DHCPDescriptor(
            args.mac, str(ipaddr), state, subnet, dhcp_ip, None,
            args.lease, random.randint(0, 50000),
        )
        mac = MacAddress(args.mac)
        self.dhcp_client_state[mac.as_redis_key()] = desc
        print("Added mac %s with DHCP rec %s" % str(mac), desc)

    def del_record(self, args):
        """
        Delete DHCP state record from the redis map.
        Args:
            args: Mac address.

        Returns: None
        """

        mac = MacAddress(args.mac)
        desc = self.dhcp_client_state[mac.as_redis_key()]
        print("Deleted mac %s with DHCP rec %s" % (str(mac), desc))
        self.dhcp_client_state[mac.as_redis_key()] = None

    def set_deafult_gw(self, args):
        """
        Set GW for given uplink network
        Args:
            args: IP address of GW.

        Returns:
        """

        gw_ip = ip_address(args.ip)
        gw_info = UplinkGatewayInfo()
        gw_info.update_ip(str(gw_ip))
        print("set Default gw IP to %s" % gw_info.get_gw_ip())

    def set_deafult_gw_mac(self, args):
        """
        Set mac address of the GW
        Args:
            args: mac address
        Returns:
        """

        gw_mac = ip_address(args.ip)
        gw_info = UplinkGatewayInfo()
        gw_info.update_mac(str(gw_mac))

        print("set Default gw mac to %s" % gw_info.get_gw_mac())


def main():
    """
    main function for cli.
    Returns: None
    """
    cli = DhcpClientCLI()

    parser = argparse.ArgumentParser(
        description='Management CLI for Mobility DHCP Client',
        formatter_class=argparse.ArgumentDefaultsHelpFormatter,
    )

    # Add sub commands
    subparsers = parser.add_subparsers(title='subcommands', dest='cmd')

    # List
    subparser = subparsers.add_parser(
        'list_dhcp_records',
        help='Lists all records from Redis',
    )
    subparser.set_defaults(func=cli.list_all_record)

    # Add
    subparser = subparsers.add_parser(
        'add_rec',
        help='Add ip allocation record',
    )
    subparser.add_argument(
        'mac', help='Mac address, e.g. "8a:00:00:00:0b:11"',
        type=str,
    )
    subparser.add_argument(
        'ip', help='IP address, e.g. "1.1.1.1"',
        type=ip_address,
    )

    subparser.add_argument(
        'state',
        help='DHCP protocol state 1 to 7, e.g. "1"',
        type=int,
    )
    subparser.add_argument(
        'subnet',
        help='IP address subnet, e.g. "1.1.1.0/24"',
        type=ipaddress.ip_network,
    )

    subparser.add_argument('dhcp', help='DHCP IP address, e.g. "1.1.1.100"')
    subparser.add_argument('lease', help='Lease time in seconds, e.g. "100"')
    subparser.set_defaults(func=cli.add_record)

    # del
    subparser = subparsers.add_parser(
        'del_rec',
        help='Add ip allocation record',
    )
    subparser.add_argument('mac', help='Mac address, e.g. "8a:00:00:00:0b:11"')
    subparser.set_defaults(func=cli.del_record)

    # set default gw
    subparser = subparsers.add_parser(
        'set_default_gw',
        help='Set default GW',
    )
    subparser.add_argument('ip', help='IP address, e.g. "1.1.1.1"')

    subparser.set_defaults(func=cli.set_deafult_gw)

    # set gw mac
    subparser = subparsers.add_parser(
        'set_gw_mac',
        help='Set GW Mac address',
    )
    subparser.add_argument('mac', help='Mac address, e.g. "8a:00:00:00:0b:11"')

    subparser.set_defaults(func=cli.set_deafult_gw)

    # Parse the args
    args = parser.parse_args()
    if not args.cmd:
        parser.print_usage()
        sys.exit(1)

    # Execute the sub-command function
    args.func(args)


if __name__ == "__main__":
    main()

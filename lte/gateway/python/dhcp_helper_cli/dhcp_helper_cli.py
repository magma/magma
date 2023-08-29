#!/usr/bin/env python3

"""
Copyright (C) 2022  The Magma Authors

This program is free software; you can redistribute it and/or
modify it under the terms of the GNU General Public License
as published by the Free Software Foundation; either version 2
of the License, or (at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program; If not, see <http://www.gnu.org/licenses/>.
"""
from __future__ import annotations

import argparse
import json
import random
import time
from datetime import datetime, timedelta
from enum import IntEnum
from queue import Empty, Queue
from typing import Any, Callable, Dict, List, Optional, Tuple

import scapy.packet
from scapy.all import AsyncSniffer
from scapy.layers.dhcp import BOOTP, DHCP
from scapy.layers.inet import IP, UDP
from scapy.layers.l2 import Dot1Q, Ether
from scapy.sendrecv import sendp


class DHCPState(IntEnum):
    UNKNOWN = 0
    DISCOVER = 1
    OFFER = 2
    REQUEST = 3
    DECLINE = 4
    ACK = 5
    NAK = 6
    RELEASE = 7
    FORCE_RENEW = 8


class MacAddress:
    def __init__(self, mac: str) -> None:
        self.mac_address = mac.lower()

    def __eq__(self, other: Any) -> bool:
        return hasattr(other, 'mac_address') and self.mac_address == other.mac_address

    def as_hex(self) -> bytes:
        return bytes.fromhex(self.mac_address.replace(':', ''))

    def as_redis_key(self, vlan: int) -> str:
        key = str(self.mac_address).replace(':', '_').lower()
        if vlan:
            return "v{}.{}".format(vlan, key)
        else:
            return key

    @classmethod
    def from_hex(cls, sid: str) -> MacAddress:
        return cls(':'.join(''.join(x) for x in zip(*[iter(sid)] * 2)))

    def __str__(self) -> str:
        return self.mac_address


class DhcpHelperCli:
    _SNIFFER_FILTER = "udp and (port 67 or 68)"
    _SNIFFER_STARTUP_WAIT = 0.5
    _TIMEOUT = 50
    _state = DHCPState.DISCOVER
    _vlan = None

    def __init__(
        self, mac: MacAddress, vlan: int, iface: str, ip: Optional[str] = None,
            server_ip: Optional[str] = None, router_ip: Optional[str] = None,
    ) -> None:
        self._lease_expiration_time: Optional[str] = None
        self._iface = iface
        self._mac = mac
        self._vlan = vlan
        self._ip = ip
        self._server_ip = server_ip
        self._router_ip = router_ip
        self._pkt_queue: Queue = Queue()
        self._ip_subnet = ""

        self._sniffer = AsyncSniffer(
            iface=iface,
            filter=(self._SNIFFER_FILTER),
            store=False,
            prn=self._receive_answer,
        )

        self._sniffer.start()

        time.sleep(self._SNIFFER_STARTUP_WAIT)

    def allocate(self) -> None:
        self.send_dhcp_discover()
        self.wait_for(self.receive_dhcp_offer)
        self.send_dhcp_request()
        self.wait_for(self.receive_dhcp_ack)

    def release(self) -> None:
        self.send_dhcp_release()
        # Receiving release ack is not mandatory

    def renew(self) -> None:
        self._state = DHCPState.OFFER
        self.send_dhcp_request()
        self.wait_for(self.receive_dhcp_ack)

    def _receive_answer(self, pkt: scapy.packet.Packet) -> None:
        if DHCP in pkt:
            self._pkt_queue.put(pkt)

    @staticmethod
    def _get_new_xid() -> int:
        return random.randint(0, 2 ** 32 - 1)

    @staticmethod
    def _get_option(packet: scapy.packet.Packet, name: str) -> Optional[str]:
        for opt in packet[DHCP].options:
            if opt[0] == name:
                return opt[1]
        return None

    def send_dhcp_discover(self) -> None:
        if self._state == DHCPState.DISCOVER:
            dhcp_opts = [
                ("message-type", "discover"),
                "end",
            ]
        else:
            print(f"Wrong previous state {DHCPState(self._state).name} != DISCOVER")
            return

        self.send_dhcp_pkt(dhcp_opts)

    def send_dhcp_request(self) -> None:
        if self._state == DHCPState.OFFER:
            self._state = DHCPState.REQUEST
            dhcp_opts = [
                ("message-type", "request"),
                ("requested_addr", self._ip),
                ("server_id", self._server_ip),
                "end",
            ]
        else:
            print(f"Wrong previous state {DHCPState(self._state).name} != OFFER")
            return

        self.send_dhcp_pkt(dhcp_opts)

    def send_dhcp_release(self) -> None:
        self._state = DHCPState.RELEASE
        dhcp_opts = [
            ("message-type", "release"),
            ("server_id", self._server_ip),
            "end",
        ]
        ciaddr = self._ip

        self.send_dhcp_pkt(dhcp_opts, ciaddr)

    def send_dhcp_pkt(self, dhcp_opts: List[Any], ciaddr: Optional[str] = None) -> None:
        pkt = Ether(src=str(self._mac), dst="ff:ff:ff:ff:ff:ff")
        if self._vlan and self._vlan != 0:
            pkt /= Dot1Q(vlan=self._vlan)
        pkt /= IP(src="0.0.0.0", dst="255.255.255.255")
        pkt /= UDP(sport=68, dport=67)
        pkt /= BOOTP(op=1, chaddr=self._mac.as_hex(), xid=self._get_new_xid(), ciaddr=ciaddr)
        pkt /= DHCP(options=dhcp_opts)
        sendp(pkt, iface=self._iface, verbose=0)

    def wait_for(self, handler: Callable) -> None:
        start_time = datetime.now()

        while datetime.now() - start_time < timedelta(seconds=self._TIMEOUT):
            try:
                pkt = self._pkt_queue.get_nowait()
            except Empty:
                continue

            dhcp_state_code = int(pkt[DHCP].options[0][1])

            if handler(dhcp_state_code, pkt):
                self._pkt_queue.task_done()
                return
            else:
                self._pkt_queue.task_done()

        raise TimeoutError(f"Timed out while waiting for {handler} after {start_time}.")

    def receive_dhcp_offer(
            self, dhcp_state_code: int, pkt: scapy.packet.Packet,
    ) -> bool:
        return self.receive_dhcp_packet(dhcp_state_code, pkt, DHCPState.OFFER)

    def receive_dhcp_ack(
            self, dhcp_state_code: int, pkt: scapy.packet.Packet,
    ) -> bool:
        return self.receive_dhcp_packet(dhcp_state_code, pkt, DHCPState.ACK)

    def receive_dhcp_packet(
            self,
            dhcp_state_code: int,
            pkt: scapy.packet.Packet,
            dhcp_state_expected: DHCPState,
    ) -> bool:
        mac_addr, vlan = self.parse_reply_header(pkt)
        if not (
            mac_addr == self._mac and vlan == self._vlan
            and dhcp_state_code == dhcp_state_expected
        ):
            return False
        if BOOTP not in pkt or pkt[BOOTP].yiaddr is None:
            return False
        self._state = dhcp_state_expected
        self.update_dhcp_state(pkt)
        return True

    def update_dhcp_state(self, pkt: scapy.packet.Packet) -> None:
        self._ip = pkt[BOOTP].yiaddr
        self._router_ip = self._get_option(pkt, "router")

        subnet_mask = self._get_option(pkt, "subnet_mask") or "32"
        self._ip_subnet = str(self._ip) + "/" + subnet_mask
        if IP in pkt:
            self._server_ip = pkt[IP].src
        self._lease_expiration_time = self._get_option(pkt, "lease_time")

    @staticmethod
    def parse_reply_header(pkt: scapy.packet.Packet) -> Tuple[MacAddress, int]:
        mac_addr = MacAddress.from_hex(pkt[BOOTP].chaddr.hex()[0:12])
        vlan: int = 0
        if Dot1Q in pkt:
            vlan = pkt[Dot1Q].vlan
        return mac_addr, vlan

    def get_info(self) -> Dict[str, Optional[str]]:
        return {
            "ip": self._ip,
            "subnet": self._ip_subnet,
            "lease_expiration_time": self._lease_expiration_time,
            "server_ip": self._server_ip,
            "router_ip": self._router_ip,
        }


def print_info(info: Dict, print_json: bool) -> None:
    if print_json:
        print(json.dumps(info))
    else:
        print(f"ip: {info['ip']}")
        print(f"subnet: {info['subnet']}")
        print(f"lease_expiration_time: {info['lease_expiration_time']}")
        print(f"server_ip: {info['server_ip']}")
        print(f"router_ip: {info['router_ip']}")


def save_to_file(info: Dict, filename: str) -> None:
    if filename:
        with open(filename, "w") as f:
            f.write(json.dumps(info))


def allocate_arg_handler(opts: argparse.Namespace) -> None:
    mac = MacAddress(opts.mac)
    vlan = int(opts.vlan)
    interface = opts.interface

    cli = DhcpHelperCli(mac, vlan, interface)
    cli.allocate()

    print_info(cli.get_info(), opts.json)
    save_to_file(cli.get_info(), opts.save_file)


def release_arg_handler(opts: argparse.Namespace) -> None:
    mac = MacAddress(opts.mac)
    vlan = int(opts.vlan)
    interface = opts.interface
    ip = opts.ip
    server_ip = opts.server_ip

    cli = DhcpHelperCli(mac, vlan, interface, ip, server_ip)
    cli.release()

    print_info(cli.get_info(), opts.json)
    save_to_file(cli.get_info(), opts.save_file)


def renew_arg_handler(opts: argparse.Namespace) -> None:
    mac = MacAddress(opts.mac)
    vlan = int(opts.vlan)
    interface = opts.interface
    ip = opts.ip
    server_ip = opts.server_ip

    cli = DhcpHelperCli(mac, vlan, interface, ip, server_ip)
    cli.renew()

    print_info(cli.get_info(), opts.json)
    save_to_file(cli.get_info(), opts.save_file)


def create_parser() -> argparse.ArgumentParser:
    parser = argparse.ArgumentParser(
        description='DHCP helper to get IPs for a DHCP IP allocator.',
        formatter_class=argparse.ArgumentDefaultsHelpFormatter,
    )

    parser.add_argument('--mac', help='MAC address to allocate/release', required=True)
    parser.add_argument('--json', help='Print the allocation/release information in json format', default=False, action='store_true')
    parser.add_argument('--save-file', help='Save to the specified file', default="")
    parser.add_argument('--vlan', help='Whether to use VLAN (0 means no VLAN)', default=0)
    parser.add_argument('--interface', help='The network interface to send the request to', default='eth0')

    subparsers = parser.add_subparsers(title='subcommands', dest='cmd')
    parser_allocate = subparsers.add_parser('allocate', help='Allocate an IP for a given MAC address')

    parser_release = subparsers.add_parser('release', help='Release the specified IP for a given MAC address')
    parser_release.add_argument('--ip', help='The IP to release', required=True)
    parser_release.add_argument('--server-ip', help='The server IP to release the IP from', required=True)

    parser_renew = subparsers.add_parser('renew', help='release ip')
    parser_renew.add_argument('--ip', help='The IP to renew', required=True)
    parser_renew.add_argument('--server-ip', help='The server IP for which to renew the IP', required=True)

    parser_allocate.set_defaults(func=allocate_arg_handler)
    parser_release.set_defaults(func=release_arg_handler)
    parser_renew.set_defaults(func=renew_arg_handler)
    return parser


def main() -> None:
    parser = create_parser()
    # Parse the args
    args = parser.parse_args()
    if not args.cmd:
        parser.print_usage()
        exit(1)

    # Execute the subcommand function
    args.func(args)


if __name__ == "__main__":
    main()

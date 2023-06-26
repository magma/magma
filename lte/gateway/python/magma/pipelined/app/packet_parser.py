"""
Copyright 2022 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
"""

import binascii
from typing import Optional, Tuple

from dpkt import arp, ethernet
from magma.mobilityd.mac import MacAddress

ETH_P_TYPES = {
    ethernet.ETH_TYPE_UNKNOWN: 'unknown',
    ethernet.ETH_TYPE_EDP: 'EDP',
    ethernet.ETH_TYPE_PUP: 'PUP',
    ethernet.ETH_TYPE_IP: 'IPv4',
    ethernet.ETH_TYPE_ARP: 'ARP',
    ethernet.ETH_TYPE_AOE: 'AOE',
    ethernet.ETH_TYPE_CDP: 'CDP',
    ethernet.ETH_TYPE_DTP: 'DTP',
    ethernet.ETH_TYPE_REVARP: 'RARP',
    ethernet.ETH_TYPE_8021Q: 'n_802_1Q',
    ethernet.ETH_TYPE_8021AD: 'n_802_1AD',
    ethernet.ETH_TYPE_QINQ1: 'n_QINQ1',
    ethernet.ETH_TYPE_QINQ2: 'n_QINQ2',
    ethernet.ETH_TYPE_IPX: 'IPX',
    ethernet.ETH_TYPE_IP6: 'IPv6',
    ethernet.ETH_TYPE_PPP: 'PPP',
    ethernet.ETH_TYPE_MPLS: 'MPLS',
    ethernet.ETH_TYPE_MPLS_MCAST: 'MPLS_multicast',
    ethernet.ETH_TYPE_PPPoE_DISC: 'PPPOE_discovery',
    ethernet.ETH_TYPE_PPPoE: 'PPPOE_session',
    ethernet.ETH_TYPE_LLDP: 'LLDP',
    ethernet.ETH_TYPE_TEB: 'TEB',
    ethernet.ETH_TYPE_PROFINET: 'profinet',
}

ARP_OP_TYPES = {
    arp.ARP_OP_REQUEST: 'ARP_REQUEST',
    arp.ARP_OP_REPLY: 'ARP_REPLY',
    arp.ARP_OP_REVREQUEST: 'RARP_REQUEST',
    arp.ARP_OP_REVREPLY: 'RARP_REPLY',
}


def _getptype(ptype: bytes):
    return ETH_P_TYPES.get(ptype[0] << 8 | ptype[1], 'unknown')


def _getop(op: bytes):
    return ARP_OP_TYPES.get(op[0] << 8 | op[1], 'unknown')


class ParseSocketPacket:
    def __init__(self, packet: bytes):
        self.eth, self.dot1q, self.arp = _parse_socket_packet(packet)

    def __str__(self):
        parsed = f"{self.eth}\n"
        if self.dot1q:
            parsed += f"{self.dot1q}\n"
        parsed += f"{self.arp}"
        return parsed


class Ethernet:
    def __init__(self, dst_mac: MacAddress, src_mac: MacAddress, ptype: bytes):
        self.dst_mac = dst_mac
        self.src_mac = src_mac
        self.ptype = _getptype(ptype)

    def __str__(self):
        return f"Ethernet header\n" \
               f"  dst_mac : {self.dst_mac}\n" \
               f"  src_mac : {self.src_mac}\n" \
               f"  ptype   : {self.ptype}"


class Dot1Q:
    def __init__(self, prio: int, dei: int, vlan: int, ptype: bytes):
        self.prio = prio
        self.dei = dei
        self.vlan = vlan
        self.ptype = _getptype(ptype)

    def __str__(self):
        return f"VLAN tag\n" \
               f"  prio  : {self.prio}\n" \
               f"  dei   : {self.dei}\n" \
               f"  vlan  : {self.vlan}\n" \
               f"  ptype : {self.ptype}"


class Arp:
    def __init__(
            self,
            hwtype: str,
            ptype: bytes,
            hwlen: int,
            plen: int,
            op: bytes,
            hwsrc: MacAddress,
            psrc: str,
            hwdst: MacAddress,
            pdst: str,
    ):
        self.hwtype = hwtype
        self.ptype = _getptype(ptype)
        self.hwlen = hwlen
        self.plen = plen
        self.op = _getop(op)
        self.hwsrc = hwsrc
        self.psrc = psrc
        self.hwdst = hwdst
        self.pdst = pdst

    def __str__(self):
        return f"ARP packet\n" \
               f"  hwtype : {self.hwtype}\n" \
               f"  ptype  : {self.ptype}\n" \
               f"  hwlen  : {self.hwlen}\n" \
               f"  plen   : {self.plen}\n" \
               f"  op     : {self.op}\n" \
               f"  hwsrc  : {self.hwsrc}\n" \
               f"  psrc   : {self.psrc}\n" \
               f"  hwdst  : {self.hwdst}\n" \
               f"  pdst   : {self.pdst}"


def _parse_socket_packet(packet: bytes) -> Tuple[Ethernet, Optional[Dot1Q], Arp]:
    eth, packet = _parse_ethernet_header(packet)
    if eth.ptype == 'n_802_1Q':
        dot1q, packet = _parse_vlan_tag(packet)
    else:
        dot1q = None
    arp_packet = _parse_arp_packet(packet)
    return eth, dot1q, arp_packet


def _parse_ethernet_header(packet: bytes) -> Tuple[Ethernet, bytes]:
    dst_mac = MacAddress(_bytes_to_mac(packet[0:6]))
    src_mac = MacAddress(_bytes_to_mac(packet[6:12]))
    ptype = packet[12:14]
    return Ethernet(dst_mac, src_mac, ptype), packet[14:]


def _parse_vlan_tag(packet: bytes) -> Tuple[Dot1Q, bytes]:
    prio = packet[0] >> 5
    dei = packet[0] & 0b00010000 >> 4
    vlan = packet[0] & 0b00001111 | packet[1]
    ptype = packet[2:4]
    return Dot1Q(prio, dei, vlan, ptype), packet[4:]


def _parse_arp_packet(packet: bytes) -> Arp:
    hwtype = hex(int(binascii.hexlify(packet[0:2])))
    ptype = packet[2:4]
    hwlen = packet[4]
    plen = packet[5]
    op = packet[6:8]
    hwsrc = MacAddress(_bytes_to_mac(packet[8:14]))
    psrc = _bytes_to_ip(packet[14:18])
    hwdst = MacAddress(_bytes_to_mac(packet[18:24]))
    pdst = _bytes_to_ip(packet[24:28])
    return Arp(hwtype, ptype, hwlen, plen, op, hwsrc, psrc, hwdst, pdst)


def _bytes_to_mac(packet: bytes) -> str:
    return ":".join(["%02x" % b for b in packet])


def _bytes_to_ip(packet: bytes) -> str:
    return ".".join([str(b) for b in packet])

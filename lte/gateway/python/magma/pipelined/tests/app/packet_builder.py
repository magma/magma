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

import abc
import logging
logging.getLogger("scapy.runtime").setLevel(logging.ERROR)
from scapy.all import Ether, IP, IPv6, ARP, TCP, UDP, ICMP, DHCP, BOOTP, \
    wrpcap, rdpcap
from scapy.contrib.gtp import GTP_U_Header

'''

ScapyPacketBuilder uses **kwargs to set a specific layer of the packet, the
names of the fields can be found in the description bellow, these are taken
from the Scapy protocol definitions. If in the future we need to set additional
fields simply pass another value and use the field names from below.
Here is the list of currently supported protocols:

    NAME       : TYPE                 = DEFAULT VALUE

Ether
    dst        : DestMACField         = (None)
    src        : SourceMACField       = (None)
    type       : XShortEnumField      = (36864)

ARP
    hwtype     : XShortField          = (1)
    ptype      : XShortEnumField      = (2048)
    hwlen      : ByteField            = (6)
    plen       : ByteField            = (4)
    op         : ShortEnumField       = (1)
    hwsrc      : ARPSourceMACField    = (None)
    psrc       : SourceIPField        = (None)
    hwdst      : MACField             = ('00:00:00:00:00:00')
    pdst       : IPField              = ('0.0.0.0')

IP
    version    : BitField             = (4)
    ihl        : BitField             = (None)
    tos        : XByteField           = (0)
    len        : ShortField           = (None)
    id         : ShortField           = (1)
    flags      : FlagsField           = (0)
    frag       : BitField             = (0)
    ttl        : ByteField            = (64)
    proto      : ByteEnumField        = (0)
    chksum     : XShortField          = (None)
    src        : Emph                 = (None)
    dst        : Emph                 = ('127.0.0.1')
    options    : PacketListField      = ([])

ICMP
    type       : ByteEnumField        = (8)
    code       : MultiEnumField       = (0)
    chksum     : XShortField          = (None)
    id         : ConditionalField     = (0)
    seq        : ConditionalField     = (0)
    ts_ori     : ConditionalField     = (84535061)
    ts_rx      : ConditionalField     = (84535061)
    ts_tx      : ConditionalField     = (84535061)
    gw         : ConditionalField     = ('0.0.0.0')
    ptr        : ConditionalField     = (0)
    reserved   : ConditionalField     = (0)
    addr_mask  : ConditionalField     = ('0.0.0.0')
    unused     : ConditionalField     = (0)

TCP
    sport      : ShortEnumField       = (20)
    dport      : ShortEnumField       = (80)
    seq        : IntField             = (0)
    ack        : IntField             = (0)
    dataofs    : BitField             = (None)
    reserved   : BitField             = (0)
    flags      : FlagsField           = (2)
    window     : ShortField           = (8192)
    chksum     : XShortField          = (None)
    urgptr     : ShortField           = (0)
    options    : TCPOptionsField      = ({})

UDP
    sport      : ShortEnumField       = (53)
    dport      : ShortEnumField       = (53)
    len        : ShortField           = (None)
    chksum     : XShortField          = (None)

GTP_U_Header
    teid       : IntField             = (None)
    length     : ShortField           = (None)

Example usage:

Getting a default packet
    default_pkt = UDPPacketBuilder().build()

Building an ARP packet
    arpb = ARPPacketBuilder()
    packet = arpb.set_arp_layer("192.168.184.1")\
                 .set_arp_src("21:12:4e:8c:98:25", "10.0.0.4")\
                 .build()

Building a TCP packet
    tcpb = TCPPacketBuilder()
    packet = tcpb.set_ether_layer("12:99:cc:97:47:4e", "00:00:00:00:00:01")\
                 .set_ip_layer("12.12.12.1", "12.12.12.4")\
                 .set_tcp_layer(89,22,8417)\
                 .set_tcp_flags("S")\
                 .build()
'''


class ScapyPacket:
    def __init__(self):
        self.Ether = Ether()
        self.IP = IP()
        self.IPv6 = IPv6()
        self.ARP = ARP()
        self.TCP = TCP()
        self.UDP = UDP()
        self.GTP_U_Header = GTP_U_Header()
        self.ICMP = ICMP()
        self.BOOTP = BOOTP()
        self.DHCP = DHCP()


class PacketBuilder(abc.ABC):
    """Interface for packet Builder"""
    def __init__(self, packet):
        self.packet = packet

    @abc.abstractmethod
    def _set_ether(self, **kwargs):
        """
        Set the Ether layer of the packet
        Args:
            fields (dict): ether layer fields
        """
        raise NotImplementedError()

    @abc.abstractmethod
    def _set_arp(self, **kwargs):
        """
        Set the ARP layer of the packet
        Args:
            fields (dict): ARP layer fields
        """
        raise NotImplementedError()

    @abc.abstractmethod
    def _set_ip(self, **kwargs):
        """
        Set the IP layer of the packet
        Args:
            fields (dict): IP layer fields
        """
        raise NotImplementedError()

    @abc.abstractmethod
    def _set_icmp(self, **kwargs):
        """
        Set the ICMP layer of the packet
        Args:
            fields (dict): ICMP layer fields
        """
        raise NotImplementedError()

    @abc.abstractmethod
    def _set_tcp(self, **kwargs):
        """
        Set the TCP layer of the packet
        Args:
            fields (dict): TCP layer fields
        """
        raise NotImplementedError()

    @abc.abstractmethod
    def _set_udp(self, **kwargs):
        """
        Set the UDP layer of the packet
        Args:
            fields (dict): UDP layer fields
        """
        raise NotImplementedError()

    @abc.abstractmethod
    def _set_gtp_u_header(self, **kwargs):
        """
        Set the GTPU layer of the packet
        Args:
            fields (dict): GTPU layer fields
        """
        raise NotImplementedError()

    @abc.abstractmethod
    def _set_bootp(self, **kwargs):
        """
        Set the BOOTP layer of the packet
        Args:
            fields (dict): BOOTP layer fields
        """
        raise NotImplementedError()

    @abc.abstractmethod
    def _set_dhcp(self, **kwargs):
        """
        Set the DHCP layer of the packet
        Args:
            fields (dict): DHCP layer fields
        """
        raise NotImplementedError()

    @abc.abstractmethod
    def build(self):
        """Build the final packet"""
        raise NotImplementedError()


class ScapyPacketBuilder(PacketBuilder):
    """Scapy packet builder implementation of PacketBuilder"""
    def __init__(self):
        super().__init__(ScapyPacket())

    def _set_ether(self, **kwargs):
        for key, value in kwargs.items():
            setattr(self.packet.Ether, key, value)

    def _set_arp(self, **kwargs):
        for key, value in kwargs.items():
            setattr(self.packet.ARP, key, value)

    def _set_ip(self, **kwargs):
        for key, value in kwargs.items():
            setattr(self.packet.IP, key, value)

    def _set_ipv6(self, **kwargs):
        for key, value in kwargs.items():
            setattr(self.packet.IPv6, key, value)

    def _set_icmp(self, **kwargs):
        for key, value in kwargs.items():
            setattr(self.packet.ICMP, key, value)

    def _set_tcp(self, **kwargs):
        for key, value in kwargs.items():
            setattr(self.packet.TCP, key, value)

    def _set_udp(self, **kwargs):
        for key, value in kwargs.items():
            setattr(self.packet.UDP, key, value)

    def _set_gtp_u_header(self, **kwargs):
        for key, value in kwargs.items():
            setattr(self.packet.GTP_U_Header, key, value)

    def _set_bootp(self, **kwargs):
        for key, value in kwargs.items():
            setattr(self.packet.BOOTP, key, value)

    def _set_dhcp(self, **kwargs):
        for key, value in kwargs.items():
            setattr(self.packet.DHCP, key, value)

    def build(self):
        pass

    @staticmethod
    def construct_from_pcap(filename):
        packets = rdpcap(filename)
        return packets

    @staticmethod
    def save_to_pcap(pkts, filename):
        for pkt in pkts:
            wrpcap(filename, pkt, append=True)


class EtherPacketBuilder(ScapyPacketBuilder):
    def set_ether_layer(self, mac_dst, mac_src):
        self._set_ether(dst=mac_dst, src=mac_src)
        return self

    def build(self):
        return self.packet.Ether


class ARPPacketBuilder(EtherPacketBuilder):
    def set_arp_layer(self, pdst):
        self._set_arp(pdst=pdst)
        return self

    def set_arp_op(self, op):
        self._set_arp(op=op)
        return self

    def set_arp_hwdst(self, hwdst):
        self._set_arp(hwdst=hwdst)
        return self

    def set_arp_src(self, hwsrc, psrc):
        self._set_arp(hwsrc=hwsrc, psrc=psrc)
        return self

    def build(self):
        return self.packet.Ether / self.packet.ARP


class IPPacketBuilder(EtherPacketBuilder):
    def set_ip_layer(self, dst, src):
        self._set_ip(dst=dst, src=src)
        return self

    def set_ttl(self, ttl):
        self._set_ip(ttl=ttl)
        return self

    def set_ip_flags(self, flags):
        self._set_ip(flags=flags)
        return self

    def build(self):
        return self.packet.Ether / self.packet.IP


class IPv6PacketBuilder(EtherPacketBuilder):
    def set_ip_layer(self, dst, src):
        self._set_ipv6(dst=dst, src=src)
        return self

    def set_ip_flags(self, flags):
        self._set_ipv6(flags=flags)
        return self

    def build(self):
        return self.packet.Ether / self.packet.IPv6


class ICMPPacketBuilder(IPPacketBuilder):
    def set_icmp_layer(self, type, code):
        self._set_icmp(type=type, code=code)
        return self

    def build(self):
        return self.packet.Ether / self.packet.IP / self.packet.ICMP


class TCPPacketBuilder(IPPacketBuilder):
    def set_tcp_layer(self, sport, dport, seq):
        self._set_tcp(dport=dport, sport=sport, seq=seq)
        return self

    def set_tcp_flags(self, flags):
        self._set_tcp(flags=flags)
        return self

    def set_ack(self, ack):
        self._set_tcp(ack=ack)
        return self

    def set_urgptr(self, urgptr):
        self.set_tcp(urgptr=urgptr)
        return self

    def build(self):
        return self.packet.Ether / self.packet.IP / self.packet.TCP


class UDPPacketBuilder(IPPacketBuilder):
    def set_udp_layer(self, sport, dport):
        self._set_udp(sport=sport, dport=dport)
        return self

    def build(self):
        return self.packet.Ether / self.packet.IP / self.packet.UDP

class GTPUHeaderPacketBuilder(UDPPacketBuilder):
    def set_gtp_u_header_layer(self, teid, length, gtp_type):
        self._set_gtp_u_header(teid=teid, length=length, gtp_type=gtp_type)
        return self

    def build(self, EncapIP):
        return self.packet.Ether / self.packet.IP / self.packet.UDP / self.packet.GTP_U_Header \
               / EncapIP

class BOOTPPacketBuilder(UDPPacketBuilder):
    def set_bootp_layer(self, op, yiaddr, siaddr, chaddr):
        self._set_bootp(op=op, yiaddr=yiaddr, siaddr=siaddr, chaddr=chaddr)
        return self

    def build(self):
        return self.packet.Ether / self.packet.IP / self.packet.UDP \
            / self.packet.BOOTP


class DHCPPacketBuilder(BOOTPPacketBuilder):
    def set_dhcp_layer(self, options):
        self._set_dhcp(options=options)
        return self

    def build(self):
        return self.packet.Ether / self.packet.IP / self.packet.UDP \
            / self.packet.BOOTP / self.packet.DHCP

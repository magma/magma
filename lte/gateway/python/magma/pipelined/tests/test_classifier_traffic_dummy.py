from magma.pipelined.tests.app.packet_injector import ScapyPacketInjector
from magma.pipelined.tests.app.start_pipelined import (
    TestSetup,
    PipelinedController,
)
from scapy.contrib.gtp import GTP_U_Header
from scapy.all import *
from magma.pipelined.app.classifier import Classifier
from scapy.all import Ether, IP, UDP, ARP

class testgtptraffic(object):
    BRIDGE = 'gtp_br0'
    IFACE = 'gtp_br0'
    MAC_1 = '5e:cc:cc:b1:49:4b'
    MAC_2 = '0a:00:27:00:00:02'
    BRIDGE_IP = '192.168.128.1'
    EnodeB_IP = '192.168.60.141'
    MTR_IP = "10.0.2.10"
    Dst_nat = '192.168.129.42'

   # Create a set of packets
    def test_traffic_flows(self):
	    
        pkt_sender = ScapyPacketInjector(self.BRIDGE)
        eth = Ether(dst=self.MAC_1, src=self.MAC_2)
        ip = IP(src=self.Dst_nat, dst='192.168.128.30')
        o_udp = UDP(sport=2152, dport=2152)
        i_udp = UDP(sport=1111, dport=2222)
        i_tcp = TCP(seq=1, sport=1111, dport=2222)
        i_ip = IP(src='192.168.60.142', dst=self.EnodeB_IP)

        arp = ARP(hwdst=self.MAC_1,hwsrc=self.MAC_2, psrc=self.Dst_nat, pdst='192.168.128.30')

        gtp_packet_udp = eth / ip / o_udp / GTP_U_Header(teid=100, length=28,gtp_type=255) / i_ip / i_udp
        gtp_packet_tcp = eth / ip / o_udp / GTP_U_Header(teid=100, length=68, gtp_type=255) / i_ip / i_tcp
        arp_packet = eth / arp
        print(gtp_packet_udp.show())
        print(gtp_packet_tcp.show())
        print(arp_packet.show())
        pkt_sender.send(gtp_packet_udp)
        pkt_sender.send(gtp_packet_tcp)
        pkt_sender.send(arp_packet)

if __name__ == "__main__":
    unittest.main()

#!/usr/bin/env python3

# sudo /home/vagrant/build/python/bin/python gtp-cmd.py 192.168.60.141 192.168.60.142 192.168.128.12 192.168.129.42 eth1

import sys
import time

from scapy.all import IP, UDP, Ether, sendp
from scapy.contrib.gtp import GTP_U_Header
from scapy.layers.l2 import getmacbyip

ip_src = sys.argv[1]
ip_dst = sys.argv[2]
i_ip_src = sys.argv[3]
i_ip_dst = sys.argv[4]
egress_dev = sys.argv[5]

dst_mac = getmacbyip(ip_dst)

eth = Ether(src='08:00:27:d3:52:d1', dst=dst_mac)
ip = IP(src=ip_src, dst=ip_dst)
udp = UDP(sport=2152, dport=2152)

i_ip = IP(src=i_ip_src, dst=i_ip_dst)
i_udp = UDP(sport=56531, dport=5001)

gtp_packet_tcp = eth / ip / udp / GTP_U_Header(teid=104) / i_ip / i_udp

sendp(gtp_packet_tcp, iface=egress_dev, count=1)
time.sleep(.5)

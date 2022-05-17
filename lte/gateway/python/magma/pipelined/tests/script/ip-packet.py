#!/usr/bin/env python3

# sudo /home/vagrant/build/python/bin/python gtp-cmd.py 192.168.128.12 192.168.129.42 eth1

import sys
import time

from scapy.all import IP, Ether, sendp
from scapy.layers.l2 import getmacbyip

ip_src = sys.argv[1]
ip_dst = sys.argv[2]
egress_dev = sys.argv[3]

dst_mac = getmacbyip(ip_dst)

eth = Ether(src='08:00:27:d3:52:d1', dst=dst_mac)
ip = IP(src=ip_src, dst=ip_dst)


ip_packet = eth / ip

sendp(ip_packet, iface=egress_dev, count=1)
time.sleep(.5)

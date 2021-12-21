#!/usr/bin/env python3

from scapy.contrib.gtp import GTP_U_Header, GTPPDUSessionContainer

from scapy.all import *
import time
import sys

dst = sys.argv[1]

print ("dst: %s" % dst)


eth = Ether()

i=IPv6()
i.dst=dst   # UE_IP
i.src='2001::2'

q=ICMPv6EchoRequest()

packet_icmp = eth / i / q
for x in range(1):
    sendp(packet_icmp, iface='gtp_br0', count=1)
    time.sleep(.5)

""" Script to send icmpv6 data"""
import sys
import time

from scapy.all import Ether, ICMPv6EchoRequest, IPv6, sendp

dst = sys.argv[1]

print("dst: %s" % dst)


eth = Ether()

i = IPv6()
i.dst = dst  # UE_IP
i.src = "2001::2"

q = ICMPv6EchoRequest()

packet_icmp = eth / i / q
for _ in range(1):
    sendp(packet_icmp, iface="gtp_br0", count=1)
    time.sleep(0.5)

#!/usr/bin/env python3
from scapy.all import *
import datetime
#NOTES: make sure to pass the right interface in iface

def gtp_echo_monitor_callback(pkt):
    d = datetime.datetime.utcnow().strftime("%H:%M:%S")
    byte_array = list(Raw(pkt[UDP].payload).load)
    print(byte_array)
    byte_array[1] = 2
    resp = IP(src=pkt[IP].dst, dst=pkt[IP].src)/UDP(sport=pkt[UDP].dport,dport=pkt[UDP].sport)/bytes(byte_array)
    send(resp)

# sniff src & dst port == 2152 and GTP message type == 1 (echo request)
sniff( iface="eth1", prn=gtp_echo_monitor_callback, filter="udp[0:2]==0x0868 and udp[2:2]==0x0868 and udp[8:2]==0x3201", store=0)

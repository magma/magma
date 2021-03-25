---
id: dp_probe_cli
title: Data Path Probe Cli Command
hide_title: true
---
# dp_probe_cli

## Overview:

This command helps an operator to probe the Datapath to the UE.


### Usage

```
# On your GW host, run the following command
$ dp_probe_cli.py -i <IMSI>
$ dp_probe_cli.py -i 1010000051011
IMSI: 1010000051011, IP: 192.168.128.15
Running: sudo ovs-appctl ofproto/trace gtp_br0 tcp,in_port=local,ip_dst=192.168.128.15,ip_src=8.8.8.8,tcp_src=80,tcp_dst=3372
Datapath Actions: set(tunnel(tun_id=0x1000308,dst=10.0.2.240,ttl=64,tp_dst=2152,flags(df|key))),pop_eth,2

# You can also supply the followin options
# -I <External IP Address>
# -P <External Port>
# -UP <UE Port>
# -p <Protocol tcp/udp/icmp>

$ dp_probe_cli.py -i 1010000051016 -I 4.2.2.2 -p tcp -P 8080 -UP 3172
IMSI: 1010000051016, IP: 192.168.128.14
Running: sudo ovs-appctl ofproto/trace gtp_br0 tcp,in_port=local,ip_dst=192.168.128.14,ip_src=4.2.2.2,tcp_src=8080,tcp_dst=3172
Datapath Actions: set(tunnel(tun_id=0x1000208,dst=10.0.2.240,ttl=64,tp_dst=2152,flags(df|key))),pop_eth,2

```



### What is this command doing?

This command is a wrapper around *ovs-appctl ofproto/trace* which simulates the Datapath of the packet.


### How to read the Output?

Based on the *IMSI* supplied to this command:
- If the UE is not connected, the command will output ***UE is not connected***
- If the UE is connected, it prints the *IP Address* of the UE and the Datapath Actions based on the ovs-appctl ofproto/trace

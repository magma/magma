---
id: version-1.5.0-debug_dp_probe
title: Datapath Probe CLI
hide_title: true
original_id: debug_dp_probe
---

# Datapath Probe CLI

The datapath probe CLI helps an operator to probe the Datapath to the UE.

## Usage

On your GW host, run the following command

```$ dp_probe_cli.py -i <IMSI> --direction <DL/UL> list_rules```

For example:
```sh
$ dp_probe_cli.py -i 1010000051013 --direction DL
IMSI: 1010000051013, IP: 192.168.128.13
Running: sudo ovs-appctl ofproto/trace gtp_br0 tcp,in_port=local,ip_dst=192.168.128.13,ip_src=8.8.8.8,tcp_src=80,tcp_dst=3372
Datapath Actions: set(tunnel(tun_id=0x1000008,dst=10.0.2.241,ttl=64,tp_dst=2152,flags(df|key))),pop_eth,2
Downlink rules: allowlist_sid-IMSI1010000051013-magma.ipv4
```

```sh
dp_probe_cli.py -i 1010000051013 --direction UL
IMSI: 1010000051013, IP: 192.168.128.13
Running: sudo ovs-appctl ofproto/trace gtp_br0 tcp,in_port=4,tun_id=0x2,ip_dst=8.8.8.8,ip_src=192.168.128.13,tcp_src=3372,tcp_dst=80
Datapath Actions: set(eth(src=02:00:00:00:00:01,dst=92:9d:a2:1f:ea:44)),1
Uplink rules: allowlist_sid-IMSI1010000051013-magma.ipv4
```

You can also supply the followin options
```
# list_rules
# -I <External IP Address>
# -P <External Port>
# -UP <UE Port>
# -p <Protocol tcp/udp/icmp>
```

```sh
dp_probe_cli.py -i 1010000051013 --direction UL -p tcp -I 4.2.2.2 -P 8080 -UP 3172
IMSI: 1010000051013, IP: 192.168.128.13
Running: sudo ovs-appctl ofproto/trace gtp_br0 tcp,in_port=4,tun_id=0x2,ip_dst=4.2.2.2,ip_src=192.168.128.13,tcp_src=3172,tcp_dst=8080
Datapath Actions: set(eth(src=02:00:00:00:00:01,dst=92:9d:a2:1f:ea:44)),1
Uplink rules: allowlist_sid-IMSI1010000051013-magma.ipv4
```

```sh
dp_probe_cli.py -i 1010000051013 --direction UL list_rules
IMSI: 1010000051013, IP: 192.168.128.13
Uplink rules: allowlist_sid-IMSI1010000051013-magma.ipv4
```

## What is this command doing?

This command is a wrapper around *ovs-appctl ofproto/trace* which simulates the Datapath of the packet based on the direction specified in the command. 'DL' - Incoming Packet, 'UL' - Outgoing packet.

## How to read the output?

Based on the *IMSI* supplied to this command
- If the UE is connected it prints:
    - The *IP Address* of the UE and the Datapath Actions based on the ovs-appctl ofproto/trace
    - The enforced uplink or downlink rules in pipelined obtained from the script ```pipelined_cli.py enforcement display_flows```
    - If the parameter `list_rules` was set, only the enforced rules in pipelined will be shown
- If the UE is not connected, the command will output ***UE is not connected***

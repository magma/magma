---
id: datapath_connectivity
title: Debugging AGW datapath issues
hide_title: true
---

## AGW datapath debugging
Following is a step by step guide for debugging datapath connectivity issues
of a UE.

AGW datapath is based on OVS but there are multiple components that handle the
packet in uplink and downlink direction. Any one of the components can result
in connectivity issues.
### Major components of datapath:
1. S1 (GTP) tunnel
2. OVS datapath
3. NAT/NonNAT forwarding plane.

You need to check if any of this component dropping packet to root cause packet
drop issues. Following steps guides through the debugging process.

### Datapath debugging when 100% packets are dropped
Debugging datapath issues is much easier when you have traffic running. This
is specially important in case of LTE to avoid UE getting into inactive state.
Inactive state changes state of datapath flows for a UE, so its hard to debug
issues when there are such state changes.
It is recommended to have `ping` or other traffic generating utility running
on UE or the server (on SGi side of the network) while debugging the issue.

1. Check magma services are up and running:
   `service magma@* status`. For datapath health mme, sessions and pipelineD are
   important services to look at. Check syslog for ERRORs from services.
   If All looks good continue to next step.
2. Check for OVS services:
```
service openvswitch-switch status
```
3. Check OVS Bridge status: gtp ports might vary depending on number of eNB
   connected sessions. but `ovs-vsctl show` should not show any port with
   any errors. If you see GTP related error run `/usr/local/bin/ovs-kmod-upgrade.sh`.
   After running this command you need to reattach UEs.

```
ovs-vsctl show
...-...-...-...-.....
    Manager "ptcp:6640"
    Bridge gtp_br0
        Controller "tcp:127.0.0.1:6633"
            is_connected: true
        Controller "tcp:127.0.0.1:6654"
            is_connected: true
        fail_mode: secure
        Port mtr0
            Interface mtr0
                type: internal
        Port g_563160a
            Interface g_563160a
                type: gtpu
                options: {key=flow, remote_ip="w.z.y.z"}
        Port ipfix0
            Interface ipfix0
                type: internal
        Port patch-up
            Interface patch-up
                type: patch
                options: {peer=patch-agw}
        Port gtp0
            Interface gtp0
                type: gtpu
                options: {key=flow, remote_ip=flow}
        Port g_963160a
            Interface g_963160a
                type: gtpu
                options: {key=flow, remote_ip="a.b.c.d"}
        Port li_port
            Interface li_port
                type: internal
        Port gtp_br0
            Interface gtp_br0
                type: internal
        Port proxy_port
            Interface proxy_port
...
    Bridge uplink_br0
        Port uplink_br0
            Interface uplink_br0
                type: internal
        Port dhcp0
            Interface dhcp0
                type: internal
        Port patch-agw
            Interface patch-agw
                type: patch
                options: {peer=patch-up}
    ovs_version: "2.14.3"
```

4. Check if UE is actually connected to datapath using:
   `mobility_cli.py get_subscriber_table`. In case the IMSI is missing in this
   table, you need to debug issue in control plane. UE is not attached to the
   AGW, you need to inspect MME logs for control plane issues.
   If UE is connection continue to next step.
5. From here onwards you are going to debug OVS datapath, so you need to select
   a UE and identify which traffic direction is broken. You can do so by
   - Generating uplink traffic in UE
   - Capturing packets on gtp_br0 in NATed datapath case and SGi_dev interface for
     NonNAT datapath: `tcpdump -eni gtp_br0 host $UE_IP`.
   - For NATed datapath you also need to check if packet are egressing on
     the SGi port. You can do so by running tcpdump on SGi port
     `tcpdump -eni $SGi_dev dst $SERVER_IP`
   - In case the packet is missing on SGi port, you have issue with the routing.
     check routing table on the AGW.
   - In case uplink packets are reaching SGi port, you need to debug issues in
     downlink direction.
6. Check if you are receiving packets from server by capturing return traffic
   packet: `tcpdump -eni $SGi_dev src $SERVER_IP`. If you do not see these packets
   you need to debug SGi network configuration.
7. Check traffic stats from UE in OVS. `dp_probe_cli.py --imsi 1234 -D UL stats`
   In case the stats show packets reaching OVS. This should be non-zero.
   For downlink traffic, Check stats for DL.
8. If all looks good so far, you need to trace packet in OVS pipeline, This command
   would show datapath action that OVS would apply to incoming packets. If it
   shows ‘drop’ it means OVS is dropping the packet,
   For tracing packets in UL direction:
   - If there is action to forward the traffic to egress port check connectivity
    between SGi interface and destination host.
   - For NonNat (Bridged mode) you might need vlan action for handling MultiAPN.
```
dp_probe_cli.py -i 414200000000029 -d UL -I 114.114.114.114 -P 80 -p tcp`.
IMSI: 414200000000029, IP: 192.168.128.12
Running: sudo ovs-appctl ofproto/trace gtp_br0 tcp,in_port=3,tun_id=0x1,ip_dst=114.114.114.114,ip_src=192.168.128.12,tcp_src=3372,tcp_dst=80
Datapath Actions: set(eth(src=02:00:00:00:00:01,dst=5e:5b:d1:8a:1a:42)),set(skb_mark(0x5)),1
Uplink rules: allowlist_sid-IMSI414200000000029-APNNAME1
```

For DL traffic: check if action show tunnel set action.
```
Dp_probe_cli.py -i 414200000000029 -d DL -I 114.114.114.114 -P 80 -p tcp
IMSI: 414200000000029, IP: 192.168.128.12
Running: sudo ovs-appctl ofproto/trace gtp_br0 tcp,in_port=local,ip_dst=192.168.128.12,ip_src=114.114.114.114,tcp_src=80,tcp_dst=3372
Datapath Actions: set(tunnel(tun_id=0xc400003f,dst=10.0.2.208,ttl=64,tp_dst=2152,flags(df|key))),pop_eth,set(skb_mark(0x4)),2
```
9. In case of DL traffic, if you see datapath action, check if the dst ip address in tunnel()
   action is the right eNB for the UE.
   - Check routing table for this IP address `ip route get $dst_ip`
   - Check if the eNB is reachable from the AGW. there could be FW rules dropping
      the packets.

10. In case probe command shows drop you need to check which table is dropping
    the packet. Manually run the OVS trace command from above output shown on line
    starting with `Running`. For above DL example `sudo ovs-appctl ofproto/trace
    gtp_br0 tcp,in_port=local,ip_dst=192.168.128.12,ip_src=114.114.114.114,tcp_src=80,tcp_dst=3372`
11. The trace command shows which table is dropping the packet. To map the numberical
    tble number to AGW pipeline table use pipelined-cli.

```
root@magma:~# pipelined_cli.py debug table_assignment
App                      Main Table          Scratch Tables
----------------------------------------------------------------------
mme                      0                   []
ingress                  1                   []
arpd                     2                   []
access_control           3                   [21]
proxy                    4                   []
middle                   10                  []
gy                       11                  [22, 23]
enforcement              12                  [24]
enforcement_stats        13                  []
egress                   20                  []
```
12. In case enforcement or gy table is dropping the packet, it means there is
    no rule for traffic or there is blocking rule for the traffic, that drops
    the packet.
    - You can check rules in datapath using dp-probe command:
    `dp_probe_cli.py -i 414200000000029 --direction UL list_rules`
    - To validate rules pushed from orc8r, you can use stat-cli: `state_cli.py
      parse "policydb:rules"`, This command would dump all rules, you need
      to check which rule are applicable to the UE.
13. Packet drops in access_control means there is static config in pipelineD
    which does not allow this connection.
14. AGW should not be dropping packet in any other table. File a bug report with
    the trace output in a github issue.
15. If this document does not help to debug the issue, please post output of
    all steps in new github issue.

### Intermittent packets drop
Intermittent packets loss is harder to debug than previous case. In this case the
services and flow tables are configured currently but still some packets are dropped.
Following are usual suspects:
1. TC queue is dropping packets due to rate limiting, command
   `pipelined_cli.py debug qos` shows stats for all dropped packets. Run the
   test case and observe if you see any dropped packets
```
root@agw:~# pipelined_cli.py debug qos
/usr/local/lib/python3.5/dist-packages/scapy/config.py:411: CryptographyDeprecationWarning: Python 3.5 support will be dropped in the next release of cryptography. Please upgrade your Python.
  import cryptography
Root stats for:  eth0
qdisc htb 1: root refcnt 2 r2q 10 default 0 direct_packets_stat 5487 ver 3.17 direct_qlen 1000
 Sent 1082274 bytes 7036 pkt (dropped 846, overlimits 4244 requeues 0)
 backlog 0b 0p requeues 0

Root stats for:  eth1
qdisc htb 1: root refcnt 2 r2q 10 default 0 direct_packets_stat 41140 ver 3.17 direct_qlen 1000
 Sent 3603343 bytes 41337 pkt (dropped 0, overlimits 0 requeues 0)
 backlog 0b 0p requeues 0
```

2. NAT could be dropping packets. This can be due to no ports available in NAT
   table due to large number of open connections. AGW has default setting for
   the max connections `sysctl net.netfilter.nf_conntrack_max` and default
   range of source port `sysctl net.ipv4.ip_local_port_range`. If you see
   higher number of simultaneous connections, you need to tune these parameters.

If none of this works file detailed bug report on github.

#!/bin/bash

# This script setup OVS bridges for userspace datapath
# Such setup is used for running integ tests only.
ovs-vsctl add-br dl_br0

ovs-vsctl set bridge uplink_br0 datapath_type=netdev
ovs-vsctl set bridge gtp_br0 datapath_type=netdev
ovs-vsctl set Interface gtp0  type=gtpu
ovs-vsctl set bridge dl_br0 datapath_type=netdev
ovs-vsctl set interface gtp0  options:csum=true

ip a f eth1
ovs-vsctl add-port dl_br0 eth1
ethtool -K eth1 gso off
ethtool -K eth1 gro off
ethtool -K dl_br0 gso off
ethtool -K dl_br0 gro off

ethtool -K gtp_br0 gso off
ethtool -K gtp_br0 gro off

ifconfig  uplink_br0 up
ifconfig  dl_br0   192.168.60.142/24  up
ifconfig  gtp_br0  192.168.128.0/24   up
ifconfig  mtr0     10.1.0.0/24        up
ifconfig  ipfix0   1.2.3.4/24         up
ifconfig  li_port  127.1.0.20/24      up

# ping 192.168.60.141  -c 3
# sudo  ovs-appctl tnl/arp/set dl_br0 192.168.60.141 `arp -n |grep  192.168.60.141 |awk '{print $3}'`
# ovs-vsctl show

ovs-ofctl add-flow dl_br0 "priority=100, in_port=eth1,sctp actions=output:dl_br0"

disable-tcp-checksumming

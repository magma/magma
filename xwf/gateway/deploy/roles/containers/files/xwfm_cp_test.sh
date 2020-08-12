#!/bin/bash
ip link add veth5 type veth peer name veth6
ifconfig veth5 up;ifconfig veth6 up
ip netns add ns5
ip link set veth5 netns ns5
ip netns exec ns5 ifconfig veth5 up
ovs-vsctl add-port uplink_br0 veth6
echo "requesting dhcp please wait..."
ip netns exec ns5 dhclient -1 veth5
echo "the ip we got is:"
ip netns exec ns5 ifconfig veth5
echo "requesting url \n"
ip netns exec ns5 curl -I www.google.com
echo "cleaning up \n"
ip netns exec ns5 ip link del veth5
ps -ef |grep dhclient | grep veth | awk '{ print $2 }' | xargs kill -9
ovs-vsctl del-port uplink_br0 veth6
ip netns del ns5
ip link del veth5 type veth peer name veth6


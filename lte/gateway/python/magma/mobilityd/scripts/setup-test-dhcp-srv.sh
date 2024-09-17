#!/bin/bash -x
prefix=$1

ip link del "$prefix"uplink_p0
sleep 0.1

ip link add "$prefix"uplink_p0 type veth peer name  "$prefix"uplink_p1
ifconfig "$prefix"uplink_p0 up

mkdir -p /etc/netns/"$prefix"dhcp_srv
touch /etc/netns/"$prefix"dhcp_srv/resolv.conf

ip netns add "$prefix"dhcp_srv

ip link set dev  "$prefix"uplink_p1 netns "$prefix"dhcp_srv
ip netns exec    "$prefix"dhcp_srv   ifconfig  "$prefix"uplink_p1 192.168.128.100/24
ip netns exec    "$prefix"dhcp_srv   ifconfig  "$prefix"uplink_p1 hw ether b2:a0:cc:85:80:7a
ip netns exec    "$prefix"dhcp_srv   ifconfig  "$prefix"uplink_p1 up
ip netns exec    "$prefix"dhcp_srv   ip addr add 192.168.128.211  dev "$prefix"uplink_p1
PID=$(pgrep -f "dnsmasq.*mobilityd.*$prefix")
if [[ -n "$PID" ]]
then
	kill "$PID"
fi

sleep 1
ip netns exec "$prefix"dhcp_srv  /usr/sbin/dnsmasq -q --conf-file=/home/vagrant/magma/lte/gateway/python/magma/mobilityd/scripts/dnsd.conf  --log-queries --log-facility=/var/log/"$prefix"dnsmasq.test.log &

logger "DHCP server started"

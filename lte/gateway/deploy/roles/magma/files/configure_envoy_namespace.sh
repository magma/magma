#!/bin/bash

envoy_ip_cntr_root="10.5.0.2"
envoy_ip_cntr="10.5.0.3"
envoylink_cntr="envoy_cntr"

Router_IP="10.6.0.1"
Router_mac=$(ip l sh gtp_br0|grep link|awk '{print $2}')

envoy_dp_dev_ip="10.6.0.2"
envoy_dp_dev="proxy_port"
mac_addr='e6:8f:a2:80:80:80'

envoy_ns="envoy_ns1"

function setup {
  # root nameapce config
  /sbin/ip link add $envoylink_cntr type veth peer name  "$envoylink_cntr"_ns

  # envoy controller IP
  /sbin/ifconfig $envoylink_cntr "$envoy_ip_cntr_root"/24 up

  # add namespace
  /sbin/ip netns add $envoy_ns

  # move devices
  /sbin/ip link set dev "$envoy_dp_dev"_ns    netns $envoy_ns
  /sbin/ip link set dev "$envoylink_cntr"_ns  netns $envoy_ns

  /sbin/ip netns exec  $envoy_ns /sbin/ip link set dev "$envoy_dp_dev"_ns address $mac_addr
  # namespace configi
  /sbin/ip netns exec  $envoy_ns /sbin/ifconfig "$envoylink_cntr"_ns  "$envoy_ip_cntr"/24 up

  /sbin/ip netns exec  $envoy_ns /sbin/ifconfig "$envoy_dp_dev"_ns    "$envoy_dp_dev_ip"/24 up

  /sbin/ip netns exec  $envoy_ns /sbin/ifconfig lo up

  /sbin/ip netns exec  $envoy_ns /sbin/ip route add default via $Router_IP

  /sbin/ip netns exec  $envoy_ns /sbin/ip neigh replace $Router_IP  lladdr $Router_mac dev "$envoy_dp_dev"_ns

  /sbin/ip netns exec  $envoy_ns /sbin/iptables -t mangle -I PREROUTING -p tcp --dport 80 -j MARK --set-mark 1
  /sbin/ip netns exec  $envoy_ns /sbin/iptables -t mangle -I PREROUTING -p tcp --sport 80 -j MARK --set-mark 1

  /sbin/ip netns exec  $envoy_ns /sbin/ip rule add fwmark 1 lookup 100
  /sbin/ip netns exec  $envoy_ns /sbin/ip route add local 0.0.0.0/0 dev lo table 100

  /sbin/ip netns exec  $envoy_ns /sbin/sysctl -w net.ipv4.conf.all.rp_filter=0
  /sbin/ip netns exec  $envoy_ns /sbin/sysctl -w net.ipv4.conf.all.route_localnet=1
}

function destroy {
  /sbin/ip netns exec  $envoy_ns /sbin/ip link set dev "$envoy_dp_dev"_ns    netns 1
  /sbin/ip netns exec  $envoy_ns /sbin/ip link set dev "$envoylink_cntr"_ns  netns 1

  /bin/sleep 1
  /sbin/ip link del $envoylink_cntr
  /bin/sleep 1
  /sbin/ip netns delete $envoy_ns
}

$1
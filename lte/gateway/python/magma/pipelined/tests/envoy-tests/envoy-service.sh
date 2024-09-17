#!/bin/bash -ex

conf="/home/vagrant/magma/lte/gateway/python/magma/pipelined/tests/envoy-tests/envoy-cntr-prod.yaml"
bin="/usr/bin/envoy"

envoy_ip_cntr_root="10.5.0.2"
envoy_ip_cntr="10.5.0.3"
envoylink_cntr="envoy_cntr"

Router_IP="10.6.0.1"
Router_mac=$(ip l sh gtp_br0|grep link|awk '{print $2}')

envoy_dp_dev_ip="10.6.0.2"
envoy_dp_dev="proxy_port"
mac_addr='e6:8f:a2:80:80:80'

envoy_ns="envoy_ns1"

function init {
  # root nameapce config
  ip link add $envoylink_cntr type veth peer name  "$envoylink_cntr"_ns

  # envoy controller IP
  ifconfig $envoylink_cntr "$envoy_ip_cntr_root"/24 up

  # add namespace
  ip netns add $envoy_ns

  # move devices
  ip link set dev "$envoy_dp_dev"_ns    netns $envoy_ns
  ip link set dev "$envoylink_cntr"_ns  netns $envoy_ns

  ip netns exec  $envoy_ns ip link set dev "$envoy_dp_dev"_ns address $mac_addr
  # namespace configi
  ip netns exec  $envoy_ns ifconfig "$envoylink_cntr"_ns  "$envoy_ip_cntr"/24 up

  ip netns exec  $envoy_ns ifconfig "$envoy_dp_dev"_ns    "$envoy_dp_dev_ip"/24 up

  ip netns exec  $envoy_ns ifconfig lo up

  ip netns exec  $envoy_ns ip route add default via $Router_IP

  ip netns exec  $envoy_ns ip neigh replace $Router_IP  lladdr $Router_mac dev "$envoy_dp_dev"_ns
  
  ip netns exec  $envoy_ns iptables -t mangle -I PREROUTING -p tcp --dport 80 -j MARK --set-mark 1
  ip netns exec  $envoy_ns iptables -t mangle -I PREROUTING -p tcp --sport 80 -j MARK --set-mark 1

  ip netns exec  $envoy_ns ip rule add fwmark 1 lookup 100
  ip netns exec  $envoy_ns ip route add local 0.0.0.0/0 dev lo table 100

  ip netns exec  $envoy_ns sysctl -w net.ipv4.conf.all.rp_filter=0
  ip netns exec  $envoy_ns sysctl -w net.ipv4.conf.all.route_localnet=1
}

function start {
  nohup /home/vagrant/magma/feg/gateway/services/envoy_controller/envoy_controller &
  sleep 1
  nohup ip netns exec  $envoy_ns $bin -c $conf -l debug&
}

function stop {
  killall envoy
  killall envoy_controller
}

function destroy {
  ip netns exec  $envoy_ns ip link set dev "$envoy_dp_dev"_ns    netns 1
  ip netns exec  $envoy_ns ip link set dev "$envoylink_cntr"_ns  netns 1

  sleep 1
  ip link del $envoylink_cntr
  sleep 1
  ip netns delete $envoy_ns
}

$1

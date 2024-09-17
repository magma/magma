#!/bin/bash -ex

envoy_ns="envoy_ns"
conf="/home/vagrant/magma/lte/gateway/python/magma/pipelined/tests/envoy-tests/envoy-cntr-prod.yaml"
bin="/usr/bin/envoy"

/home/vagrant/magma/feg/gateway/services/envoy_controller/envoy_controller&
sleep 5

ip netns exec envoy_ns bash -x ./pkt-routing-setup.sh
sleep 1
ip netns exec envoy_ns $bin -c $conf -l debug&



#!/bin/bash -ex

envoy_ns="envoy_ns"
conf="/home/vagrant/magma/lte/gateway/python/magma/pipelined/tests/envoy-tests/envoy.yaml"
bin="/usr/bin/envoy"

bash -x ./pkt-routing-setup.sh
sleep 1
$bin -c $conf -l debug

Setup OVS flow table using S1ap test. Once flow table is up stop services and follow the steps:

1. sudo bash -x sim-ue.sh  s 1 192.168.128.72 0xa000128
2. sudo python http-serve.py&
3. sudo bash -x envoy-service.sh init
4. sudo ip netns exec envoy_ns1 bash
   -> /usr/bin/envoy -c /home/vagrant/magma/lte/gateway/python/magma/pipelined/tests/envoy-tests/envoy.yaml  -l debug

Validate ping:
sudo ip netns exec ue_ns_1 ping 192.168.128.1

Validate http:
sudo ip netns exec ue_ns_1 curl   192.168.128.1:80/index

Destroy:
sudo bash -x sim-ue.sh  d 1
sudo bash -x envoy-service.sh destroy


DEV="eth3"
NS_NAME="trf_ns"
LINK_NAME="trf"
MAGMA_DEV=192.168.60.142

function setup()
{
  sshpass -p vagrant ssh vagrant@$MAGMA_DEV -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no sudo ip r r 192.168.60.144 via 192.168.129.42 dev eth2

  ethtool --offload eth1 rx off tx off
  ethtool --offload eth2 rx off tx off
  ethtool --offload eth3 rx off tx off

  echo "1" > /proc/sys/net/ipv4/ip_forward
  ip netns add $NS_NAME

  ip link add "$LINK_NAME"1  type veth peer name  "$LINK_NAME"2

  ip link set dev  $DEV  netns $NS_NAME
  ip link set dev  "$LINK_NAME"2  netns $NS_NAME
  ifconfig "$LINK_NAME"1 up

  ip link set dev  "$LINK_NAME"1 address "08:00:27:62:75:8b"
  ip r r 192.168.60.144 dev "$LINK_NAME"1

  ip netns exec $NS_NAME ifconfig lo up
  ip netns exec $NS_NAME ifconfig "$LINK_NAME"2 0.0.0.0 up
  ip netns exec $NS_NAME ip link set dev "$LINK_NAME"2 address "08:00:27:62:75:8b"
  ip netns exec $NS_NAME ifconfig  $DEV up
  ip netns exec $NS_NAME ip a add 192.168.60.144/24 dev $DEV
  ip netns exec $NS_NAME ip a add 192.168.129.42/24 dev $DEV
  ip netns exec $NS_NAME ip r add 10.0.2.15       dev  "$LINK_NAME"2
  ip netns exec $NS_NAME ip r add 192.168.60.141  dev  "$LINK_NAME"2
  ip netns exec $NS_NAME ip r add 192.168.128.11  dev  "$LINK_NAME"2
  ip netns exec $NS_NAME ip r add default  via 192.168.129.1 dev $DEV

  sleep 2
  ip netns exec $NS_NAME /usr/sbin/sshd
  sleep 1
  nohup ip netns exec $NS_NAME /home/vagrant/magma/lte/gateway/deploy/roles/trfserver/files/traffic_server.py 192.168.60.144 62462 &
}

function destroy()
{
  ip netns exec $NS_NAME ip link set $DEV  netns 1
  sleep 1
  ip netns del $NS_NAME
  ip link  del "$LINK_NAME"1

  ifconfig  $DEV up
}

function reset()
{
  destroy
  setup
}

$1

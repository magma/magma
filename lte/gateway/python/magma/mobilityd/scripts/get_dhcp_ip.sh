#!/bin/bash -e

PORT_MAC=$1
PORT=$2
BR=$3
NS="dhcp_cl"

ip netns add "$NS"
ovs-vsctl add-port "$BR" "$PORT" -- set interface "$PORT" type=internal -- set interface "$PORT" mac=\""$PORT_MAC"\"

ip l set dev "$PORT" netns "$NS"
nohup ip netns exec "$NS" dhclient "$PORT"

logger "dhclient started: $PORT_MAC on $PORT"

#!/bin/bash

PORT_MAC=$1
BR=$2

NS="dhcp_cl"

PORT=$(ip netns exec "$NS" ip l sh |grep "$PORT_MAC" -B1 | head -1 | grep -v link|cut -d: -f 2| xargs)
if [ -z "$PORT" ]
then
    PORT=$(ip l sh |grep "$PORT_MAC" -B1 |head -1 |grep -v link|cut -d: -f 2| xargs)
fi
PID=$(pgrep -f "dhclient.*$PORT")
if [[ -n $PID ]]
then
    kill "$PID"
fi
ip netns exec "$NS" ip link set "$PORT" netns 1
sleep .2

ovs-vsctl del-port "$BR" "$PORT"
logger "IP released for $PORT with mac: $PORT_MAC"

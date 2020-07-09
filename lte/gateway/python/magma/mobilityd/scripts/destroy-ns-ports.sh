#!/bin/bash

BR=$1
NS="dhcp_cl"
DEVS=$(ip netns exec "$NS" ip l sh|grep -v lo|grep -v link|cut -d: -f 2| xargs)

for DEV in $DEVS
do
	echo "$DEV"
        MAC=$(ip netns exec "$NS" ip l sh |grep "$DEV" -A1|grep link|xargs|cut -d' ' -f2)
        bash ./python/magma/mobilityd/scripts/release_dhcp_ip.sh "$MAC" "$BR"
done

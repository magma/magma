#!/bin/bash -ex

br="etest"

envoy_port="10000"
dst_port="80"

up1_ip="3.3.3.3"
up2_ip="4.4.4.4"
ue1_ip="2.2.2.2"
ue2_ip="2.2.2.222"
envoy_ip="2.2.2.20"
envoy_ip_cntr="10.5.0.1"
envoy_ip_cntr_root="10.5.0.2"

Router_IP="2.2.2.1"
MAC="11:22:33:44:55:66"
mac_hex="112233445566"

Router_IP_up1_net="3.3.3.1"
Router_IP_up2_net="4.4.4.1"
MAC_up1_net="11:22:33:44:55:77"
mac_hex_up1_net="112233445577"

ue1link="ue1_dev"
ue2link="ue2_dev"
envoylink="envoy_dev"
envoylink_cntr="envoy_devc"
up1link="up1_dev"
up2link="up2_dev"

ue1_ns="ue1_ns"
ue2_ns="ue2_ns"
up1_ns="up1_ns"
up2_ns="up2_ns"
envoy_ns="envoy_ns"

ue1_ip_net="2.2.2.0"

function setup {
	ip link add $ue1link  type veth peer name  "$ue1link"_ns
	ip link add $ue2link  type veth peer name  "$ue2link"_ns
	ip link add $envoylink type veth peer name  "$envoylink"_ns
	ip link add $envoylink_cntr type veth peer name  "$envoylink_cntr"_ns
	ip link add $up1link    type veth peer name  "$up1link"_ns
	ip link add $up2link    type veth peer name  "$up2link"_ns

	ip netns add $up1_ns
	ip netns add $up2_ns
	ip netns add $ue1_ns
	ip netns add $ue2_ns
	ip netns add $envoy_ns

	ip link set dev "$up1link"_ns  netns $up1_ns
	ip link set dev "$up2link"_ns  netns $up2_ns
	ip link set dev "$ue1link"_ns  netns $ue1_ns
	ip link set dev "$ue2link"_ns  netns $ue2_ns
	ip link set dev "$envoylink"_ns  netns $envoy_ns
	ip link set dev "$envoylink_cntr"_ns  netns $envoy_ns

	ip netns exec  $up1_ns      ifconfig "$up1link"_ns        "$up1_ip"/24    up 
	ip netns exec  $up2_ns      ifconfig "$up2link"_ns        "$up2_ip"/24    up 
	ip netns exec  $ue1_ns      ifconfig "$ue1link"_ns      "$ue1_ip"/24    up
	ip netns exec  $ue2_ns      ifconfig "$ue2link"_ns      "$ue2_ip"/24    up
	ip netns exec  $envoy_ns   ifconfig "$envoylink"_ns       "$envoy_ip"/24 up
	ip netns exec  $envoy_ns   ifconfig "$envoylink_cntr"_ns  "$envoy_ip_cntr"/24 up

	ip netns exec  $up1_ns     ifconfig lo up
	ip netns exec  $up2_ns     ifconfig lo up
	ip netns exec  $ue1_ns     ifconfig lo up
	ip netns exec  $ue2_ns     ifconfig lo up
	ip netns exec  $envoy_ns   ifconfig lo up

	ip netns exec  $ue1_ns      ip route add default via $Router_IP
	ip netns exec  $ue2_ns      ip route add default via $Router_IP
	ip netns exec  $envoy_ns   ip route add default via $Router_IP
	ip netns exec  $up1_ns      ip route add default via $Router_IP_up1_net
	ip netns exec  $up2_ns      ip route add default via $Router_IP_up2_net

	# ip netns exec  $envoy_ns   sysctl net.ipv4.ip_nonlocal_bind=1
	# setup bridge
	ovs-vsctl --may-exist add-br "$br"
	ovs-vsctl --may-exist add-port "$br" "$up1link"
	ovs-vsctl --may-exist add-port "$br" "$up2link"
	ovs-vsctl --may-exist add-port "$br" "$ue1link"
	ovs-vsctl --may-exist add-port "$br" "$ue2link"
	ovs-vsctl --may-exist add-port "$br" "$envoylink"

        ifconfig "$envoylink_cntr"  "$envoy_ip_cntr_root"/24 up

	ifconfig "$br" up
	ifconfig "$up1link" up
	ifconfig "$up2link" up
	ifconfig "$ue1link" up
	ifconfig "$ue2link" up
	ifconfig "$envoylink" up
	ip route add "$ue1_ip_net"/24 dev "$br"

	up1_ip_mac=$(ip netns exec $up1_ns  ip l sh "$up1link"_ns | grep ether| awk '{print $2}')
	up2_ip_mac=$(ip netns exec $up2_ns  ip l sh "$up2link"_ns | grep ether| awk '{print $2}')
	ue1_ip_mac=$(ip netns exec $ue1_ns  ip l sh "$ue1link"_ns | grep ether| awk '{print $2}')
	ue2_ip_mac=$(ip netns exec $ue2_ns  ip l sh "$ue2link"_ns | grep ether| awk '{print $2}')
	envoy_ip_mac=$(ip netns exec $envoy_ns  ip l sh "$envoylink"_ns | grep ether| awk '{print $2}')

	ovs-ofctl del-flows $br 
	# from internet
	ovs-ofctl add-flow $br "table=0, priority=200,in_port=$up1link,ip,tcp,tp_src=$dst_port actions=mod_dl_dst:$envoy_ip_mac,output:$envoylink"
	ovs-ofctl add-flow $br "table=0, priority=200,in_port=$up2link,ip,tcp,tp_src=$dst_port actions=mod_dl_dst:$envoy_ip_mac,output:$envoylink"

	ovs-ofctl add-flow $br "table=0, priority=200,in_port=$envoylink,ip,ip_dst=$ue1_ip,tcp actions=mod_dl_dst:$ue1_ip_mac,output:$ue1link"
	ovs-ofctl add-flow $br "table=0, priority=200,in_port=$envoylink,ip,ip_dst=$ue2_ip,tcp actions=mod_dl_dst:$ue2_ip_mac,output:$ue2link"

	# to internet
	ovs-ofctl add-flow $br "table=0, priority=100,in_port=$ue1link,ip,tcp,tp_dst=$dst_port actions=mod_dl_dst:$envoy_ip_mac,output:$envoylink"
	ovs-ofctl add-flow $br "table=0, priority=100,in_port=$ue2link,ip,tcp,tp_dst=$dst_port actions=mod_dl_dst:$envoy_ip_mac,output:$envoylink"

	ovs-ofctl add-flow $br "table=0, priority=100,in_port=$envoylink,ip,ip_dst=$up1_ip,tcp,tp_dst=$dst_port actions=mod_dl_dst:$up1_ip_mac,output:$up1link"
	ovs-ofctl add-flow $br "table=0, priority=100,in_port=$envoylink,ip,ip_dst=$up2_ip,tcp,tp_dst=$dst_port actions=mod_dl_dst:$up2_ip_mac,output:$up2link"

	ovs-ofctl add-flow $br "table=0, priority=123,arp,arp_tpa=$Router_IP/24,arp_op=1 actions=move:NXM_OF_ETH_SRC[]->NXM_OF_ETH_DST[],mod_dl_src:$MAC,load:0x2->NXM_OF_ARP_OP[],move:NXM_NX_ARP_SHA[]->NXM_NX_ARP_THA[],load:0x$mac_hex->NXM_NX_ARP_SHA[],move:NXM_OF_ARP_TPA[]->NXM_NX_REG0[],move:NXM_OF_ARP_SPA[]->NXM_OF_ARP_TPA[],move:NXM_NX_REG0[]->NXM_OF_ARP_SPA[],IN_PORT"
	ovs-ofctl add-flow $br "table=0, priority=123,arp,arp_tpa=$Router_IP_up1_net,arp_op=1 actions=move:NXM_OF_ETH_SRC[]->NXM_OF_ETH_DST[],mod_dl_src:$MAC_up1_net,load:0x2->NXM_OF_ARP_OP[],move:NXM_NX_ARP_SHA[]->NXM_NX_ARP_THA[],load:0x$mac_hex_up1_net->NXM_NX_ARP_SHA[],move:NXM_OF_ARP_TPA[]->NXM_NX_REG0[],move:NXM_OF_ARP_SPA[]->NXM_OF_ARP_TPA[],move:NXM_NX_REG0[]->NXM_OF_ARP_SPA[],IN_PORT"
	ovs-ofctl add-flow $br "table=0, priority=123,arp,arp_tpa=$Router_IP_up2_net,arp_op=1 actions=move:NXM_OF_ETH_SRC[]->NXM_OF_ETH_DST[],mod_dl_src:$MAC_up1_net,load:0x2->NXM_OF_ARP_OP[],move:NXM_NX_ARP_SHA[]->NXM_NX_ARP_THA[],load:0x$mac_hex_up1_net->NXM_NX_ARP_SHA[],move:NXM_OF_ARP_TPA[]->NXM_NX_REG0[],move:NXM_OF_ARP_SPA[]->NXM_OF_ARP_TPA[],move:NXM_NX_REG0[]->NXM_OF_ARP_SPA[],IN_PORT"
	
	ovs-ofctl add-flow $br "table=0, priority=001,arp  actions=NORMAL"

	ip netns exec  $up1_ns python http-serve.py &
	ip netns exec  $up2_ns python http-serve.py &
}


function destroy {

	ip netns exec  $ue1_ns     ip link set dev "$ue1link"_ns   netns 1
	ip netns exec  $ue2_ns     ip link set dev "$ue2link"_ns   netns 1
	ip netns exec  $envoy_ns  ip link set dev "$envoylink"_ns  netns 1
	ip netns exec  $envoy_ns  ip link set dev "$envoylink_cntr"_ns  netns 1

	ip link del "$ue1link"
	ip link del "$ue2link"
	ip link del "$envoylink"
	ip link del "$envoylink_cntr"
	ip link del "$up1link"
	ip link del "$up2link"

	sleep 1
	ip netns delete $up1_ns
	ip netns delete $up2_ns
	ip netns delete $ue1_ns
	ip netns delete $ue2_ns
	ip netns delete $envoy_ns

	sleep 1

	ovs-vsctl --if-exist del-br "$br"

	PIDs=$(pgrep -f "python.*http-serve.py")
	for PID in $PIDs
	do
		    kill "$PID"
        done
	PIDs=$(pgrep -f "envoy.*magma")
	for PID in $PIDs
	do
		    kill "$PID"
        done
	killall envoy_controller
}

$1

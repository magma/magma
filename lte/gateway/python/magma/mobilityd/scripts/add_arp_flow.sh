#!/bin/bash

cmd=$1
IP=$2
MAC=$3
port=$4
mac_hex=${MAC//:/}


if [[ "$cmd" == "add" ]]
then
	ovs-ofctl add-flow up_br0 "table=0, priority=123,in_port=$port,arp,arp_tpa=$IP,arp_op=1 actions=move:NXM_OF_ETH_SRC[]->NXM_OF_ETH_DST[],mod_dl_src:$MAC,load:0x2->NXM_OF_ARP_OP[],move:NXM_NX_ARP_SHA[]->NXM_NX_ARP_THA[],load:0x$mac_hex->NXM_NX_ARP_SHA[],move:NXM_OF_ARP_TPA[]->NXM_NX_REG0[],move:NXM_OF_ARP_SPA[]->NXM_OF_ARP_TPA[],move:NXM_NX_REG0[]->NXM_OF_ARP_SPA[],IN_PORT"
else
	ovs-ofctl del-flows up_br0 "table=0, priority=123,in_port=$port,arp,arp_tpa=$IP,arp_op=1"
fi

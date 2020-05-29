#!/usr/bin/env bash
#
# Copyright (c) 2018-present, Facebook, Inc.
# All rights reserved.
#
# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree. An additional grant
# of patent rights can be found in the PATENTS file in the same directory.

# gw0 configuration, TODO might want to avoid NORMAL for perf reasons
ip_addr=$(/sbin/ifconfig gw0 | sed -En 's/127.0.0.1//;s/.*inet (addr:)?(([0-9]*\.){3}[0-9]*).*/\2/p')
ovs-ofctl add-flow uplink_br0 "table=0, priority=999, arp, nw_src=$ip_addr, actions=NORMAL"
ovs-ofctl add-flow uplink_br0 "table=0, priority=999, arp, nw_dst=$ip_addr, actions=NORMAL"
ovs-ofctl add-flow uplink_br0 "table=0, priority=999, ip, nw_src=$ip_addr, actions=NORMAL"
ovs-ofctl add-flow uplink_br0 "table=0, priority=999, ip, nw_dst=$ip_addr, actions=NORMAL"

# Some setups might not have 2 nics. In that case just use eth2
if [ -d "/sys/class/net/eth2" ] && [ -d "/sys/class/net/eth3" ]
then
  ovs-vsctl --may-exist add-port uplink_br0 eth2
  ovs-vsctl --may-exist add-port uplink_br0 eth3
  ovs-ofctl -O OpenFlow14 add-group uplink_br0 "group_id=42, type=select, selection_method=dp_hash, bucket=actions=output:eth2, bucket=actions=output:eth3"
  ovs-ofctl add-flow uplink_br0 "table=0, priority=10, in_port=eth2, actions=output:uplink_patch"
  ovs-ofctl add-flow uplink_br0 "table=0, priority=10, in_port=eth3, actions=output:uplink_patch"
  ovs-ofctl add-flow uplink_br0 "table=0, priority=10, in_port=uplink_patch, actions=group:42"
elif [ -d "/sys/class/net/eth2" ]
then
  ovs-vsctl --may-exist add-port uplink_br0 eth2
  ovs-ofctl add-flow uplink_br0 "table=0, priority=10, in_port=eth2, actions=output:uplink_patch"
  ovs-ofctl add-flow uplink_br0 "table=0, priority=10, in_port=uplink_patch, actions=output:eth2"
else
  ovs-ofctl add-flow uplink_br0 "table=0, priority=10, actions=LOCAL"
fi

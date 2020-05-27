#!/usr/bin/env bash
#
# Copyright (c) 2018-present, Facebook, Inc.
# All rights reserved.
#
# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree. An additional grant
# of patent rights can be found in the PATENTS file in the same directory.

sudo ovs-vsctl add-br uplink_br0
sudo ovs-vsctl set-fail-mode uplink_br0 secure
sudo ovs-ofctl del-flows uplink_br0

sudo ovs-vsctl --may-exist add-port uplink_br0 gw0 \
  -- set Interface gw0 type=internal \
  -- set interface gw0 ofport=1
sudo ifconfig gw0 up

sudo ovs-vsctl --may-exist add-port uplink_br0 uplink_patch \
  -- set Interface uplink_patch type=patch options:peer=cwag_patch \
  -- --may-exist add-port cwag_br0 cwag_patch \
  -- set Interface cwag_patch type=patch  options:peer=uplink_patch

# gw0 configuration, TODO might want to avoid NORMAL for perf reasons
ip_addr=$(sudo /sbin/ifconfig gw0 | grep -m 1 inet | awk '{ print $2}')
sudo ovs-ofctl add-flow uplink_br0 "table=0, priority=999, arp, nw_src=$ip_addr, actions=NORMAL"
sudo ovs-ofctl add-flow uplink_br0 "table=0, priority=999, arp, nw_dst=$ip_addr, actions=NORMAL"
sudo ovs-ofctl add-flow uplink_br0 "table=0, priority=999, ip, nw_src=$ip_addr, actions=NORMAL"
sudo ovs-ofctl add-flow uplink_br0 "table=0, priority=999, ip, nw_dst=$ip_addr, actions=NORMAL"

# Some setups might not have 2 nics. In that case just use eth2
if [ -d "/sys/class/net/eth2" ] && [ -d "/sys/class/net/eth3" ]
then
  sudo ovs-vsctl --may-exist add-port uplink_br0 eth2
  sudo ovs-vsctl --may-exist add-port uplink_br0 eth3
  sudo ovs-ofctl -O OpenFlow14 add-group uplink_br0 "group_id=42, type=select, selection_method=dp_hash, bucket=actions=output:eth2, bucket=actions=output:eth3"
  sudo ovs-ofctl add-flow uplink_br0 "table=0, priority=10, in_port=eth2, actions=output:uplink_patch"
  sudo ovs-ofctl add-flow uplink_br0 "table=0, priority=10, in_port=eth3, actions=output:uplink_patch"
  sudo ovs-ofctl add-flow uplink_br0 "table=0, priority=10, in_port=uplink_patch, actions=group:42"
elif [ -d "/sys/class/net/eth2" ]
then
  sudo ovs-vsctl --may-exist add-port uplink_br0 eth2
  sudo ovs-ofctl add-flow uplink_br0 "table=0, priority=10, in_port=eth2, actions=output:uplink_patch"
  sudo ovs-ofctl add-flow uplink_br0 "table=0, priority=10, in_port=uplink_patch, actions=output:eth2"
else
  sudo ovs-ofctl add-flow uplink_br0 "table=0, priority=10, actions=LOCAL"
fi

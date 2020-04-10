#!/usr/bin/env bash
#
# Copyright (c) 2018-present, Facebook, Inc.
# All rights reserved.
#
# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree. An additional grant
# of patent rights can be found in the PATENTS file in the same directory.

sudo ovs-vsctl add-br uplink_br0
sudo ovs-vsctl --may-exist add-port uplink_br0 gw0 \
  -- set Interface gw0 type=internal \
  -- set interface gw0 ofport=1

sudo ovs-vsctl --may-exist add-port uplink_br0 uplink_patch \
  -- set Interface uplink_patch type=patch options:peer=cwag_patch \
  -- --may-exist add-port cwag_br0 cwag_patch \
  -- set Interface cwag_patch type=patch  options:peer=uplink_patch

# Some setups might not have 2 nics. In that case just use eth2
if [ ! -d "/sys/class/net/eth3" ]
then
  sudo ovs-vsctl --may-exist add-port uplink_br0 eth2
else
  sudo ovs-vsctl --may-exist add-bond uplink_br0 uplink_bond eth2 eth3
fi

#!/usr/bin/env bash
#
# Copyright (c) 2018-present, Facebook, Inc.
# All rights reserved.
#
# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree. An additional grant
# of patent rights can be found in the PATENTS file in the same directory.

sudo ovs-vsctl add-br uplink_br0

sudo ovs-vsctl --may-exist add-port uplink_br0 uplink_patch
sudo ovs-vsctl --may-exist add-port cwag_br0 cwag_patch
sudo ovs-vsctl set interface uplink_patch type=patch options:peer=cwag_patch
sudo ovs-vsctl set interface cwag_patch type=patch options:peer=uplink_patch

sudo ovs-vsctl --may-exist add-bond uplink_br0 uplink_bond eth2 eth3

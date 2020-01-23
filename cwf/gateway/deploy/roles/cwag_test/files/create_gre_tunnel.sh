#!/usr/bin/env bash
#
# Copyright (c) 2018-present, Facebook, Inc.
# All rights reserved.
#
# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree. An additional grant
# of patent rights can be found in the PATENTS file in the same directory.

sudo ovs-vsctl --may-exist add-br cwag_test_br0
sudo ovs-vsctl --may-exist add-port cwag_test_br0 gre0 -- set interface gre0 type=gre options:remote_ip=192.168.70.101
sudo ovs-ofctl add-flow cwag_test_br0 in_port=cwag_test_br0,actions=gre0
sudo ovs-ofctl add-flow cwag_test_br0 in_port=gre0,actions=cwag_test_br0

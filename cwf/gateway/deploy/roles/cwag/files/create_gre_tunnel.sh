#!/usr/bin/env bash
#
# Copyright (c) 2018-present, Facebook, Inc.
# All rights reserved.
#
# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree. An additional grant
# of patent rights can be found in the PATENTS file in the same directory.

sudo ovs-vsctl --may-exist add-br cwag_br0
sudo sysctl net.ipv4.ip_forward=1
sudo iptables -t mangle -A FORWARD -i cwag_br0 -p tcp --tcp-flags SYN,RST SYN -j TCPMSS --set-mss 1400
sudo iptables -t mangle -A FORWARD -o cwag_br0 -p tcp --tcp-flags SYN,RST SYN -j TCPMSS --set-mss 1400
sudo ovs-vsctl --may-exist add-port cwag_br0 gre0 -- set interface gre0 ofport_request=32768 type=gre options:remote_ip=flow

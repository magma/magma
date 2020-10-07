#!/bin/bash
#
# Copyright 2020 The Magma Authors.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License

ip link add veth5 type veth peer name veth6
ifconfig veth5 up;ifconfig veth6 up
ip netns add ns5
ip link set veth5 netns ns5
ip netns exec ns5 ifconfig veth5 up
ovs-vsctl add-port uplink_br0 veth6
echo "requesting dhcp please wait..."
ip netns exec ns5 dhclient -1 veth5
echo "the ip we got is:"
ip netns exec ns5 ifconfig veth5
echo "requesting url \n"
ip netns exec ns5 curl -s -I www.google.com
echo "cleaning up \n"
ip netns exec ns5 ip link del veth5
ps -ef |grep dhclient | grep veth | awk '{ print $2 }' | xargs kill -9
ovs-vsctl del-port uplink_br0 veth6
ip netns del ns5


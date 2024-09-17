#!/usr/bin/env bash
#
# Copyright 2020 The Magma Authors.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# gw0 configuration, TODO might want to avoid NORMAL for perf reasons
sudo ifconfig gw0 up
ip_addr=$(sudo /sbin/ifconfig gw0 | grep -m 1 inet | awk '{ print $2}')
sudo ovs-ofctl add-flow uplink_br0 "table=0, priority=999, arp, nw_src=$ip_addr, actions=NORMAL"
sudo ovs-ofctl add-flow uplink_br0 "table=0, priority=999, arp, nw_dst=$ip_addr, actions=NORMAL"
sudo ovs-ofctl add-flow uplink_br0 "table=0, priority=999, ip, nw_src=$ip_addr, actions=NORMAL"
sudo ovs-ofctl add-flow uplink_br0 "table=0, priority=999, ip, nw_dst=$ip_addr, actions=NORMAL"

# Some setups might not have 2 nics. In that case just use $1
if [ -n "$1" ] && [ -n "$2" ]
then
  sudo ovs-vsctl --may-exist add-port uplink_br0 "$1"
  sudo ovs-vsctl --may-exist add-port uplink_br0 "$2"
  sudo ovs-ofctl -O OpenFlow14 add-group uplink_br0 "group_id=42, type=select, selection_method=dp_hash, bucket=actions=output:$1, bucket=actions=output:$2"
  sudo ovs-ofctl add-flow uplink_br0 "table=0, priority=10, in_port=$1, actions=output:uplink_patch"
  sudo ovs-ofctl add-flow uplink_br0 "table=0, priority=10, in_port=$2, actions=output:uplink_patch"
  sudo ovs-ofctl add-flow uplink_br0 "table=0, priority=10, in_port=uplink_patch, actions=group:42"
elif [ -n "$1" ]
then
  sudo ovs-vsctl --may-exist add-port uplink_br0 "$1"
  sudo ovs-ofctl add-flow uplink_br0 "table=0, priority=10, in_port=$1, actions=output:uplink_patch"
  sudo ovs-ofctl add-flow uplink_br0 "table=0, priority=10, in_port=uplink_patch, actions=output:$1"
else
  sudo ovs-ofctl add-flow uplink_br0 "table=0, priority=10, actions=LOCAL"
fi

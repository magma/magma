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

sudo ovs-vsctl --may-exist add-port cwag_br0 li_port -- set Interface li_port type=internal
sudo ifconfig li_port up

# LI might not be enabled for all setups
if [ -d "/sys/class/net/$2" ]
then
  sudo ovs-vsctl --may-exist add-port cwag_br0 "$2"

  # Setup tc rules to mirror traffic to li mirror bridge
  ip_addr=$(sudo /sbin/ifconfig "$1" | grep -m 1 inet | awk '{ print $2}')
  sudo tc qdisc add dev "$1" handle ffff: ingress
  sudo tc filter add dev "$1" parent ffff: protocol ip u32 match ip protocol 17 0xff  match ip dport 1812 0xfffe match ip dst "$ip_addr" action mirred egress mirror dev li_port

  sudo tc qdisc add dev "$1" handle 1: root prio
  sudo tc filter add dev "$1" parent 1: protocol ip u32 match ip protocol 17 0xff match ip sport 1812 0xfffe match ip src "$ip_addr" action mirred egress mirror dev li_port
fi

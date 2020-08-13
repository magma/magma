#!/bin/bash -x
#
# Copyright 2020 The Magma Authors.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
# of patent rights can be found in the PATENTS file in the same directory.
#

br_name=$1
tag=$2

prefix="vt$tag"
ip_addr="10.200.$tag.111/24"
router_ip="10.200.$tag.211"

ip link add "$prefix"_port0 type veth peer name  "$prefix"_port1
ifconfig "$prefix"_port0 up

mkdir -p /etc/netns/"$prefix"_dhcp_srv
touch /etc/netns/"$prefix"_dhcp_srv/resolv.conf

ip netns add "$prefix"_dhcp_srv

ip link set dev  "$prefix"_port1 netns "$prefix"_dhcp_srv
ip netns exec    "$prefix"_dhcp_srv   ifconfig  "$prefix"_port1 "$ip_addr"
ip netns exec    "$prefix"_dhcp_srv   ifconfig  "$prefix"_port1 hw ether b2:a0:cc:85:80:7a
ip netns exec    "$prefix"_dhcp_srv   ip addr add $router_ip  dev "$prefix"_port1

ip netns exec    "$prefix"_dhcp_srv   ifconfig  "$prefix"_port1 up

PID=$(pgrep -f "dnsmasq.*mobilityd.*$prefix")
if [[ -n "$PID" ]]
then
	kill "$PID"
fi

sed "s/.x./."$tag"./g" /home/vagrant/magma/lte/gateway/python/magma/mobilityd/scripts/dnsd.x.conf > /tmp/dns."$tag".conf
sleep 1
ip netns exec "$prefix"_dhcp_srv  /usr/sbin/dnsmasq -q --conf-file=/tmp/dns."$tag".conf --log-queries --log-facility=/var/log/"$prefix"dnsmasq."$tag".test.log &

logger "DHCP server started"

ovs-vsctl --may-exist add-port "$br_name" "$prefix"_port0 tag="$tag"

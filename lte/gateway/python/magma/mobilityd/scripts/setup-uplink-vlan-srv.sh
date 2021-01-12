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
network_id=$3

if [[ -z $network_id ]]; then
        network_id="200"
fi

prefix="vt$tag"
ip_addr="10.$network_id.$tag.111/24"
router_ip="10.$network_id.$tag.211"

ip link del "$prefix"_port
ip link add "$prefix"_port type veth peer name  "$prefix"_port1
ifconfig "$prefix"_port up

mkdir -p /etc/netns/"$prefix"_dhcp_srv
touch /etc/netns/"$prefix"_dhcp_srv/resolv.conf

ip netns add "$prefix"_dhcp_srv

ip link set dev  "$prefix"_port1 netns "$prefix"_dhcp_srv
ip netns exec    "$prefix"_dhcp_srv   ifconfig  "$prefix"_port1 "$ip_addr"
ip netns exec    "$prefix"_dhcp_srv   ifconfig  "$prefix"_port1 hw ether b2:a0:cc:85:80:$tag
ip netns exec    "$prefix"_dhcp_srv   ip addr add $router_ip  dev "$prefix"_port1

ip netns exec    "$prefix"_dhcp_srv   ifconfig  "$prefix"_port1 up

PID=$(pgrep -f "dnsmasq.*mobilityd.*$prefix")
if [[ -n "$PID" ]]
then
	kill "$PID"
fi

sed "s/.x./."$network_id.$tag"./g" /home/vagrant/magma/lte/gateway/python/magma/mobilityd/scripts/dnsd.x.conf > /tmp/dns."$tag".conf
sleep 1
ip netns exec "$prefix"_dhcp_srv  /usr/sbin/dnsmasq -q --conf-file=/tmp/dns."$tag".conf --log-queries --log-facility=/var/log/"$prefix"dnsmasq."$tag".test.log &

logger "DHCP server started"

existing_br=$(sudo ovs-vsctl iface-to-br vt1_port)
if [[ ! -z $existing_br ]]; then
	ovs-vsctl del-port "$existing_br" "$prefix"_port
fi

ovs-vsctl --may-exist add-port "$br_name" "$prefix"_port
if [[ "$tag" -ne "0" ]]
then
	ovs-vsctl set port "$prefix"_port tag="$tag"
fi

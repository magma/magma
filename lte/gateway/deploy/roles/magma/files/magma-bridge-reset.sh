#!/bin/bash
# Copyright 2021 The Magma Authors.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
GTP_BR="gtp_br0"

# ARG 1
# -f: force reload kmod and OVS restart
# -y: reload kmod and restart OVS in case uplink-br is not configured
# -n: no kmod reload / no OVS restart
# any other value: restart OVS in case uplink-br is not configured

KMOD_RELOAD=${1:-""}

# ARG 2
# uplink bridge name
SGI_BR=${2:-"uplink_br0"}

# ARG 3: used to detect Non-NAT
# sgi port name.
SGI_PORT=$(ovs-ofctl show "$SGI_BR"|grep eth|cut -d'(' -f2|cut -d')' -f1)
if [[ -z $SGI_PORT ]];
then
  # This is Non-NAT case, SGi-port is egress interface
  SGI_PORT=$3
fi

# ARG 4
# static SGi IP
SGI_IP=$(ip a sh "$SGI_BR"|grep 'inet ' |awk '{print $2}')
if [[ -z $SGI_IP ]];
then
  SGI_IP=$4
else
  SGI_IP_SET="true"
fi

# ARG 5
# SGi network GW IP.
SGI_DEF_GW=$(ip r sh|grep default|grep "$SGI_BR"| awk '{print $3}')
if [[ -z $SGI_DEF_GW ]];
then
  SGI_DEF_GW=$5
fi


if [[ "$KMOD_RELOAD" != '-f' ]];
then
  if [[ $SGI_IP_SET == "true" ]];
  then
    echo "Uplink-br has IP address. dont reset the bridge"
    ifdown "$GTP_BR"
    ifdown patch-up
    sleep 1
    ifup "$GTP_BR"
    ifup patch-up
    exit 0
  fi
fi
# local variables
FLOW_DUMP="$(mktemp)"

#check DHCP client
DHCP_PID=$(pgrep -a 'dhclient' | grep "$SGI_BR" | awk '{print $1}')
if [[ -n $DHCP_PID ]];
then
  for pid in $DHCP_PID
  do
    kill "$pid"
  done
fi

# start reset procedure:
# save flows
ovs-ofctl dump-flows --no-names --no-stats "$SGI_BR" | \
            sed -e '/NXST_FLOW/d' \
                -e '/OFPST_FLOW/d' \
                -e 's/\(idle\|hard\)_age=[^,]*,//g' > "$FLOW_DUMP"

# remove OVS objects
ovs-vsctl --all destroy Flow_Sample_Collector_Set

ifdown "$SGI_BR"
ifdown "$GTP_BR"
ifdown patch-up
if [ "$KMOD_RELOAD" == '-n' ];
then
  :
elif [ "$KMOD_RELOAD" != '-y' ] && [ "$KMOD_RELOAD" != '-f' ];
then
  service openvswitch-switch restart
else
  /etc/init.d/openvswitch-switch  force-reload-kmod
fi

# create OVS objects
sleep 1
ifup "$SGI_BR"
ifup "$GTP_BR"
ifup patch-up
sleep 1

if [[ -n $SGI_PORT ]];
then
  ovs-vsctl --may-exist add-port "$SGI_BR" "$SGI_PORT"
  ip link set dev "$SGI_PORT"  up
fi

# restore OVS flows
ovs-ofctl del-flows "$SGI_BR"
ovs-ofctl add-flows "$SGI_BR" "$FLOW_DUMP"
rm "$FLOW_DUMP"

# start DHCP client if needed
if [[ -n $DHCP_PID ]];
then
  dhclient "$SGI_BR" &
else
  #restore IP config
  if [[ -n $SGI_IP ]];
  then
    ip a add "$SGI_IP" dev "$SGI_BR"
  fi
  if [[ -n $SGI_DEF_GW ]];
  then
    ip r add default via "$SGI_DEF_GW"
  fi
fi

#!/bin/bash
# Copyright 2021 The Magma Authors.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

key_dir=$1
allowed_ips=$2

#Optional ARGs
local_ip=$3
peer_pub_key=$4
peer_ip=$5

dev=magma_wg0
portno=9333
priv_key_file="$key_dir/wg_privatekey"

mkdir -p "$key_dir"

if [[ ! -s $priv_key_file ]]; then
    wg genkey | tee "$key_dir"/wg_privatekey | wg pubkey > "$key_dir"/wg_publickey

    chmod 600 "$key_dir"/wg_privatekey
    chmod 600 "$key_dir"/wg_publickey
fi

if ! ip l sh $dev  &> /dev/null;
then
    echo "create $dev";
    ip link add dev $dev type wireguard

    ip address flush dev $dev
    ip address replace dev $dev "$local_ip"

    ip link set $dev up
fi

curren_allowed_ips=$(wg show $dev | grep allowed | cut -d: -f2 | xargs)
if [[ -n $curren_allowed_ips ]]; then
    allowed_ips="$allowed_ips,$curren_allowed_ips"
fi
if [[ -n $local_ip ]]; then
  allowed_ips="$allowed_ips,$local_ip"
fi

if [[ -n $peer_ip ]]; then
  end_point_arg="$peer_ip:$portno"
else
  end_point_arg="$(wg show magma_wg0 | grep endpoint | cut -d: -f2 | xargs):$portno"
fi

if [[ -z $peer_pub_key ]]; then
  peer_pub_key=$(wg show |grep 'peer:' | cut -d: -f2 | xargs)
fi

wg set $dev listen-port $portno  private-key "$priv_key_file" peer "$peer_pub_key"  allowed-ips "$allowed_ips" endpoint "$end_point_arg"  persistent-keepalive 15

cp /etc/wireguard/$dev.conf /etc/wireguard/$dev.conf.bk
touch /etc/wireguard/$dev.conf
wg-quick save $dev
if ! diff -q /etc/wireguard/$dev.conf /etc/wireguard/$dev.conf.bk; then
  ip link del $dev
  wg-quick up $dev
fi
systemctl enable wg-quick@$dev

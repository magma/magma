#!/bin/bash
# Copyright 2021 The Magma Authors.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

SKIP_CHECK=${1:-""}

echo "This upgrade would result in datapath level downtime for few minutes!"

if [[ "$SKIP_CHECK" != '-y' ]];
then
    while true; do
        read -p "Do you want to proceed with upgrade ?(y/n)" yn
        case $yn in
            [Yy]* ) break;;
            [Nn]* ) exit;;
            * ) echo "Please answer yes or no.";;
        esac
    done
fi

apt update
apt install -y  openvswitch-datapath-dkms libopenvswitch openvswitch-common openvswitch-switch python3-openvswitch

dkms autoinstall
service magma@* stop
sleep 5
ifdown gtp_br0
ifdown uplink_br0
sleep 5
/etc/init.d/openvswitch-switch  force-reload-kmod
sleep 5
ifup uplink_br0
ifup gtp_br0
sleep 5
service magma@magmad start

kernel_ver=$(cat /sys/module/openvswitch/srcversion)
mod_ver=$(modinfo /lib/modules/"$(uname -r)"/updates/dkms/openvswitch.ko |grep srcversion|awk '{print $2}')

if [[ "$kernel_ver" == "$mod_ver" ]]; then
        echo "update successful"
else
        echo "update failed"
fi

echo "To cleanup control plane state run: config_stateless_agw.py flushall_redis"

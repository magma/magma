#!/bin/bash
# Copyright 2021 The Magma Authors.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

port_name=$1
enb_addr=$2
gtp_echo=$3
gtp_csum=$4
enable_wg_tuneling=$5

wg_setup_utility="/usr/local/bin/magma-setup-wg.sh"
wg_key_dir="/var/opt/magma/sgi-tunnel"
bfd_time=5000

sudo ovs-vsctl --may-exist add-port gtp_br0 $port_name -- set Interface $port_name type=gtpu options:remote_ip=$enb_addr options:key=flow bfd:enable=$gtp_echo options:csum=$gtp_csum bfd:min_tx=$bfd_time bfd:min_rx=$bfd_time

if [[ $enable_wg_tuneling == "true" ]]; then
  $wg_setup_utility $wg_key_dir $enb_addr
fi

#!/bin/bash -ex
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

br=$1
prefix=$2
for pid in $(pgrep -f  "dnsmasq.*vt.*\..*\.")
do
        kill $pid
done
ip -all netns  delete

set +e
ip link delete "$prefix"_ul_0
set -e

ip link add "$prefix"_ul_0 type veth peer name  "$prefix"_ul_1
ifconfig "$prefix"_ul_0 up
ifconfig "$prefix"_ul_1 up


# setup bridge

ovs-vsctl --if-exist del-br "$br"
ovs-vsctl add-br "$br"
ovs-vsctl --may-exist add-port "$br" "$prefix"_ul_1

ifconfig "$br" up

#ovs-vsctl set Bridge "$br" fail_mode=secure


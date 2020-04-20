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

sudo ovs-vsctl add-br uplink_br0
sudo ovs-vsctl set-fail-mode uplink_br0 secure
sudo ovs-ofctl del-flows uplink_br0

sudo ovs-vsctl --may-exist add-port uplink_br0 gw0 \
  -- set Interface gw0 type=internal \
  -- set interface gw0 ofport=1

sudo ovs-vsctl --may-exist add-port uplink_br0 uplink_patch \
  -- set Interface uplink_patch type=patch options:peer=cwag_patch \
  -- --may-exist add-port cwag_br0 cwag_patch \
  -- set Interface cwag_patch type=patch  options:peer=uplink_patch

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

sudo ovs-vsctl --may-exist add-br cwag_test_br0
sudo ovs-vsctl --may-exist add-port cwag_test_br0 gre0 -- set interface gre0 type=gre options:remote_ip=192.168.70.101
sudo ovs-ofctl add-flow cwag_test_br0 in_port=cwag_test_br0,actions=gre0
sudo ovs-ofctl add-flow cwag_test_br0 in_port=gre0,actions=cwag_test_br0

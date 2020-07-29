#!/bin/bash

# Copyright 2020 The Magma Authors.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#
cd /usr/share/yang/models || exit

validate_yang.sh fbc-symphony-device.yang

# The order of models matters here for some strange reason so just be aware.
yanglint --strict \
  fbc-symphony-device.yang \
  openconfig-access-points.yang \
  openconfig-ap-manager.yang \
  openconfig-extensions.yang \
  openconfig-wifi-mac.yang \
  openconfig-wifi-phy.yang \
  openconfig-wifi-types.yang \
  openconfig-interfaces.yang \
  openconfig-if-ip.yang \
  iana-if-type.yang \
  ietf-interfaces.yang \
  ietf-system.yang \
  /validate.json

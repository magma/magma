#!/bin/bash
# Copyright 2024 The Magma Authors.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

  ln -sf /var/run/systemd/resolve/resolv.conf /etc/resolv.conf
  sed -i 's/#DNS=/DNS=8.8.8.8 208.67.222.222/' /etc/systemd/resolved.conf
  service systemd-resolved restart

  echo 'debconf debconf/frontend select Noninteractive' | debconf-set-selections

  cat > /etc/apt/apt.conf.d/20auto-upgrades << EOF
  APT::Periodic::Update-Package-Lists "0";
  APT::Periodic::Download-Upgradeable-Packages "0";
  APT::Periodic::AutocleanInterval "0";
  APT::Periodic::Unattended-Upgrade "0";
EOF
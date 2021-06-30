#!/bin/bash
################################################################################
# Copyright 2020 The Magma Authors.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
################################################################################

set -e

# Packer may ssh into the box too early since SSH is ready before Debian
# actually is
sleep 30

# Adding the snapshot to retrieve 4.9.0-9-amd64, install the kernel, then
# remove this snapshot
echo "deb http://snapshot.debian.org/archive/debian/20190801T025637Z stretch main non-free contrib" >> /etc/apt/sources.list
apt-get update
apt install -y linux-image-4.9.0-9-amd64 linux-headers-4.9.0-9-amd64
sed -i '/20190801T025637Z/d' /etc/apt/sources.list

# Install some packages
apt-get update
apt-get install -y openssh-server gcc rsync dirmngr

# Add the Etagecom magma repo
bash -c 'echo -e "deb https://artifactory.magmacore.org/artifactory/debian stretch-1.5.0 main" > /etc/apt/sources.list.d/packages_magma_etagecom_io.list'

# Create the preferences file for backports
bash -c 'cat <<EOF > /etc/apt/preferences.d/magma-preferences
Package: *
Pin: origin artifactory.magmacore.org
Pin-Priority: 900
EOF'

# Add the Etagecom key
wget https://artifactory.magmacore.org:443/artifactory/api/gpg/key/public -O /tmp/public
apt-key add /tmp/public
apt-get update

# Disable daily auto updates, so that vagrant ansible scripts can
# acquire apt lock immediately on startup
systemctl stop apt-daily.timer
systemctl disable apt-daily.timer
systemctl disable apt-daily.service
systemctl daemon-reload

echo "Done"

#!/bin/bash -eux
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

# Add ubuntu user to sudoers.
echo "ubuntu        ALL=(ALL)       NOPASSWD: ALL" >> /etc/sudoers
echo "vagrant        ALL=(ALL)       NOPASSWD: ALL" >> /etc/sudoers
sed -i "s/^.*requiretty/#Defaults requiretty/" /etc/sudoers

# Disable daily apt unattended updates.
echo 'APT::Periodic::Enable "0";' >> /etc/apt/apt.conf.d/10periodic

apt update
apt install -y ansible

# Install some packages
apt-get update
apt-get install -y openssh-server gcc rsync dirmngr
apt-get install -y apt-transport-https ca-certificates

# Add the artifactory magma repo
bash -c 'echo -e "deb https://artifactory.magmacore.org/artifactory/debian focal-1.6.0 main" > /etc/apt/sources.list.d/packages_magma_etagecom_io.list'


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

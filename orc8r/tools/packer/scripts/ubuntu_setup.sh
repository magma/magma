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

# Mount the guest additions iso and run the install script
mkdir -p /mnt/iso
mount -t iso9660 -o loop /home/vagrant/VBoxGuestAdditions.iso /mnt/iso
/mnt/iso/VBoxLinuxAdditions.run || :

umount /mnt/iso
rm -rf /mnt/iso /home/vagrant/VBoxGuestAdditions.iso

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

# Install the dependencies
apt-get install -y build-essential linux-headers-"$(uname -r)"

# Mount the guest additions iso and run the install script
mkdir -p /mnt/iso
mount -t iso9660 -o loop /home/vagrant/VBoxGuestAdditions.iso /mnt/iso

/mnt/iso/VBoxLinuxAdditions.run

umount /mnt/iso
rm -rf /mnt/iso /home/vagrant/VBoxGuestAdditions.iso

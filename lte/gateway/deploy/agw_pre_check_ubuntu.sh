#!/bin/bash
# Copyright 2021 The Magma Authors.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# Check for system changes before magma deploy
# Setting up env variable, user and project path
MAGMA_USER="magma"

echo "- Check if Ubuntu is installed"
if ! grep -q 'Ubuntu' /etc/issue; then
  echo "Ubuntu is not installed"
else
  echo "Ubuntu is installed"
fi

echo "- Check for magma user"
if ! (getent passwd | grep -q 'magma'); then
    echo "magma user is not Installed"
elif  ! grep -q "$MAGMA_USER ALL=(ALL) NOPASSWD:ALL" /etc/sudoers; then
    echo "magma will be added to sudoers"
else
    echo "magma user configured"
fi

echo "- Check if both interfaces are named eth0 and eth1"
INTERFACES=$(ip -br a)
if [[ ! $INTERFACES == *'eth0'*  ]] || [[ ! $INTERFACES == *'eth1'* ]] || ! grep -q 'GRUB_CMDLINE_LINUX="net.ifnames=0 biosdevname=0"' /etc/default/grub; then
  echo "Interfaces will be renamed to eth0 and eth1"
  echo "eth0 will be set to dhcp and eth1 10.0.2.1"
else
  echo "eth0 will be set to dhcp and eth1 10.0.2.1"
fi

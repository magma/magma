#!/bin/bash
# Copyright 2021 The Magma Authors.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# Setting up env variable, user and project path
set -x

addr1="$1"
gw_addr="$2"

MAGMA_USER="vagrant"
AGW_INSTALL_CONFIG_LINK="/etc/systemd/system/multi-user.target.wants/agw_installation.service"
AGW_INSTALL_CONFIG="/lib/systemd/system/agw_installation.service"
AGW_SCRIPT_PATH="/root/agw_install_ubuntu.sh"
DEPLOY_PATH="/home/$MAGMA_USER/magma/lte/gateway/deploy"
SUCCESS_MESSAGE="ok"
NEED_REBOOT=0
WHOAMI=$(whoami)
MAGMA_VERSION="${MAGMA_VERSION:-v1.7}"
CLOUD_INSTALL="cloud"
GIT_URL="${GIT_URL:-https://github.com/magma/magma.git}"
INTERFACE_DIR="/etc/network/interfaces.d"

echo "Checking if the script has been executed by root user"
if [ "$WHOAMI" != "root" ]; then
  echo "You're executing the script as $WHOAMI instead of root.. exiting"
  exit 1
fi

echo "Checking if Ubuntu is installed"
if ! grep -q 'Ubuntu' /etc/issue; then
  echo "Ubuntu is not installed"
  exit 1
fi

apt-get update

echo "Need to check if both interfaces are named eth0 and eth1"
INTERFACES=$(ip -br a)
if [[ $1 != "$CLOUD_INSTALL" ]] && ( [[ ! $INTERFACES == *'eth0'*  ]] || [[ ! $INTERFACES == *'eth1'* ]] || ! grep -q 'GRUB_CMDLINE_LINUX="net.ifnames=0 biosdevname=0"' /etc/default/grub); then
  # changing intefaces name
  sed -i 's/GRUB_CMDLINE_LINUX=""/GRUB_CMDLINE_LINUX="net.ifnames=0 biosdevname=0"/g' /etc/default/grub
  sed -i 's/enp0s3/eth0/g' /etc/netplan/50-cloud-init.yaml
  # changing interface name
  grub-mkconfig -o /boot/grub/grub.cfg

  # name server config
  ln -sf /var/run/systemd/resolve/resolv.conf /etc/resolv.conf
  sed -i 's/#DNS=/DNS=8.8.8.8 208.67.222.222/' /etc/systemd/resolved.conf
  service systemd-resolved restart

  # interface config
  apt install -y ifupdown net-tools ipcalc
  mkdir -p "$INTERFACE_DIR"
  echo "source-directory $INTERFACE_DIR" > /etc/network/interfaces

  if [ -z "$addr1" ] || [ -z "$gw_addr" ]
  then
    # DHCP allocated interface IP
    echo "auto eth0
    iface eth0 inet dhcp" > "$INTERFACE_DIR"/eth0
  else
    # Statically allocated interface IP
    if ipcalc -c "$addr1" | grep INVALID
    then
      echo "Interface ip is not valid IP"
      exit 1
    fi

    if ipcalc -c "$gw_addr" | grep INVALID
    then
      echo "Upstream Router ip is not valid IP"
      exit 1
    fi

    addr=$(   ipcalc -n "$addr1"  | grep Address | awk '{print $2}')
    netmask=$(ipcalc -n "$addr1"  | grep Netmask | awk '{print $2}')
    gw_addr=$(ipcalc -n "$gw_addr"| grep Address | awk '{print $2}')

    echo "auto eth0
  iface eth0 inet static
  address $addr
  netmask $netmask
  gateway $gw_addr" > "$INTERFACE_DIR"/eth0
  fi

  # configuring eth1
  echo "auto eth1
  iface eth1 inet static
  address 192.168.60.142
  netmask 255.255.255.0" > "$INTERFACE_DIR"/eth1

  # get rid of netplan
  systemctl unmask networking
  systemctl enable networking

  apt-get --assume-yes purge nplan netplan.i

  # Setting REBOOT flag to 1 because we need to reload new interface and network services.
  NEED_REBOOT=1
else
  echo "Interfaces name are correct, let's check if network and DNS are up"
  while ! nslookup google.com; do
    echo "DNS not reachable"
    sleep 1
  done

  while ! ping -c 1 -W 1 -I eth0 8.8.8.8; do
    echo "Network not ready yet"
    sleep 1
  done
fi

echo "Making sure $MAGMA_USER user is sudoers"
if ! grep -q "$MAGMA_USER ALL=(ALL) NOPASSWD:ALL" /etc/sudoers; then
  apt install -y sudo
  adduser --disabled-password --gecos "" $MAGMA_USER
  adduser $MAGMA_USER sudo
  echo "$MAGMA_USER ALL=(ALL) NOPASSWD:ALL" >> /etc/sudoers
fi

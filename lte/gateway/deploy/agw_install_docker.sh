#!/bin/bash
# Copyright 2021 The Magma Authors.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

MODE=$1
RERUN=0    # Set to 1 to skip network configuration and run ansible playbook only
WHOAMI=$(whoami)
MAGMA_USER="ubuntu"
MAGMA_VERSION="${MAGMA_VERSION:-v1.8}"
GIT_URL="${GIT_URL:-https://github.com/magma/magma.git}"
DEPLOY_PATH="/opt/magma/lte/gateway/deploy"

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


ROOTCA="/var/opt/magma/certs/rootCA.pem"
if [ ! -f "$ROOTCA" ]; then
  echo "Upload rootCA to $ROOTCA"
  exit 1
fi

if [ $RERUN -eq 0 ]; then
  # Update DNS resolvers
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

  apt purge --auto-remove unattended-upgrades -y
  apt-mark hold "$(uname -r)" linux-aws linux-headers-aws linux-image-aws

  # interface config
  INTERFACE_DIR="/etc/network/interfaces.d"
  mkdir -p "$INTERFACE_DIR"
  echo "source-directory $INTERFACE_DIR" > /etc/network/interfaces

  # get rid of netplan
  systemctl unmask networking
  systemctl enable networking

  echo "Install Magma"
  apt-get update -y
  apt-get upgrade -y
  apt-get install curl zip python3-pip net-tools sudo ca-certificates gnupg lsb-release -y

  mkdir -p /etc/apt/keyrings
  curl -fsSL https://download.docker.com/linux/ubuntu/gpg | gpg --dearmor -o /etc/apt/keyrings/docker.gpg
  echo \
  "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu \
  $(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null

  apt-get update -y
  apt-get install docker-ce docker-ce-cli containerd.io docker-compose-plugin -y


  echo "Making sure $MAGMA_USER user is sudoers"
  if ! grep -q "$MAGMA_USER ALL=(ALL) NOPASSWD:ALL" /etc/sudoers; then
    adduser --disabled-password --gecos "" $MAGMA_USER
    adduser $MAGMA_USER sudo
    echo "$MAGMA_USER ALL=(ALL) NOPASSWD:ALL" >> /etc/sudoers

    adduser $MAGMA_USER docker
  fi

  alias python=python3
  # TODO GH13915 pinned for now because of breaking change in ansible-core 2.13.4
  pip3 install ansible==5.0.1

  rm -rf /opt/magma/
  git clone "${GIT_URL}" /opt/magma
  cd /opt/magma || exit
  git checkout "$MAGMA_VERSION"

  # changing intefaces name
  sed -i 's/GRUB_CMDLINE_LINUX=""/GRUB_CMDLINE_LINUX="net.ifnames=0 biosdevname=0"/g' /etc/default/grub
  sed -i 's/ens5/eth0/g; s/ens6/eth1/g' /etc/netplan/50-cloud-init.yaml
  # changing interface name
  update-grub2

  if ! grep -q "eth1" /etc/netplan/50-cloud-init.yaml; then
    cat > /etc/netplan/70-secondary-itf.yaml << EOF
    network:
        ethernets:
            eth1:
                dhcp4: true
                dhcp6: false
                dhcp4-overrides:
                  route-metric: 200
        version: 2
EOF
  fi
  netplan apply
fi

echo "Generating localhost hostfile for Ansible"
echo "[agw_docker]
127.0.0.1 ansible_connection=local" > $DEPLOY_PATH/agw_hosts

if [ "$MODE" == "base" ]; then
  su - $MAGMA_USER -c "sudo ansible-playbook -v -e \"MAGMA_ROOT='/opt/magma' OUTPUT_DIR='/tmp'\" -i $DEPLOY_PATH/agw_hosts --tags base $DEPLOY_PATH/magma_docker.yml"
else
  # install magma and its dependencies including OVS.
  su - $MAGMA_USER -c "sudo ansible-playbook -v -e \"MAGMA_ROOT='/opt/magma' OUTPUT_DIR='/tmp'\" -i $DEPLOY_PATH/agw_hosts --tags agwc $DEPLOY_PATH/magma_docker.yml"
fi

[[ $RERUN -eq 1 ]] || echo "Reboot this VM to apply kernel settings"

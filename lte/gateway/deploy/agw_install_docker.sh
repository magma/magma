#!/bin/bash
# Copyright 2021 The Magma Authors.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

MAGMA_VERSION="master"
WHOAMI=$(whoami)
MAGMA_USER="ubuntu"
MAGMA_VERSION="${MAGMA_VERSION:-v1.6}"
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

echo "Making sure $MAGMA_USER user is sudoers"
if ! grep -q "$MAGMA_USER ALL=(ALL) NOPASSWD:ALL" /etc/sudoers; then
  apt install -y sudo
  adduser --disabled-password --gecos "" $MAGMA_USER
  adduser $MAGMA_USER sudo
  echo "$MAGMA_USER ALL=(ALL) NOPASSWD:ALL" >> /etc/sudoers
fi

# unlink resolv.conf to avoid overrides by systemd-resolved
cp /etc/resolv.conf /etc/resolv.conf.bak
rm -f /etc/resolv.conf && mv /etc/resolv.conf.bak /etc/resolv.conf

ROOTCA="/var/opt/magma/certs/rootCA.pem"
if [ ! -f "$ROOTCA" ]; then

  echo "Upload rootCA to $ROOTCA"
  exit 1
fi

echo "Install Magma"
apt-get update
apt-get -y install curl make virtualenv zip rsync git software-properties-common python3-pip python-dev apt-transport-https

alias python=python3
pip3 install ansible

rm -rf /opt/magma/
git clone "${GIT_URL}" /opt/magma
cd /opt/magma || exit
git checkout "$MAGMA_VERSION"

cp -f $DEPLOY_PATH/roles/magma/files/magma_ifaces_gtp /etc/network/interfaces.d/gtp
cp -f $DEPLOY_PATH/roles/magma/files/ovs-kmod-upgrade.sh /usr/local/bin/
cp -f $DEPLOY_PATH/roles/magma/files/magma_modules_load /etc/modules-load.d/magma.conf

echo "Generating localhost hostfile for Ansible"
echo "[agw_docker]
127.0.0.1 ansible_connection=local
[magma_deploy]
127.0.0.1 ansible_connection=local" > $DEPLOY_PATH/agw_hosts

# install magma and its dependencies including OVS.
su - $MAGMA_USER -c "sudo ansible-playbook -e \"MAGMA_ROOT='/opt/magma' OUTPUT_DIR='/tmp'\" -i $DEPLOY_PATH/agw_hosts --tags agwc $DEPLOY_PATH/magma_deploy.yml"
su - $MAGMA_USER -c "sudo ansible-playbook -e \"MAGMA_ROOT='/opt/magma' OUTPUT_DIR='/tmp'\" -i $DEPLOY_PATH/agw_hosts $DEPLOY_PATH/magma_docker.yml"

echo "Cleanup temp files"
cd /root || exit

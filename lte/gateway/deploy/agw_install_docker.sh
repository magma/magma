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
WHOAMI=$(whoami)
MAGMA_USER="ubuntu"
MAGMA_VERSION="${MAGMA_VERSION:-master}"
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

ROOTCA="/var/opt/magma/certs/rootCA.pem"
if [ ! -f "$ROOTCA" ]; then
  echo "Upload rootCA to $ROOTCA"
  exit 1
fi

echo "Install Magma"
apt-get update -y
apt-get install curl zip python3-pip -y

alias python=python3
pip3 install ansible

rm -rf /opt/magma/
git clone "${GIT_URL}" /opt/magma
cd /opt/magma || exit
git checkout "$MAGMA_VERSION"


echo "Generating localhost hostfile for Ansible"
echo "[agw_docker]
127.0.0.1 ansible_connection=local" > $DEPLOY_PATH/agw_hosts

if [ "$MODE" == "base" ]; then
  su - $MAGMA_USER -c "sudo ansible-playbook -e \"MAGMA_ROOT='/opt/magma' OUTPUT_DIR='/tmp'\" -i $DEPLOY_PATH/agw_hosts --tags base $DEPLOY_PATH/magma_docker.yml"
else
  # install magma and its dependencies including OVS.
  su - $MAGMA_USER -c "sudo ansible-playbook -e \"MAGMA_ROOT='/opt/magma' OUTPUT_DIR='/tmp'\" -i $DEPLOY_PATH/agw_hosts --tags agwc $DEPLOY_PATH/magma_docker.yml"
fi
cd /root || exit

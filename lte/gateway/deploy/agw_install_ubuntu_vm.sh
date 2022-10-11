#!/bin/bash
# Copyright 2022 The Magma Authors.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# Setting up env variable, user and project path
set -x

MAGMA_USER="vagrant"
AGW_INSTALL_CONFIG="/lib/systemd/system/agw_installation.service"
DEPLOY_PATH="/home/$MAGMA_USER/magma/lte/gateway/deploy"
SUCCESS_MESSAGE="ok"
WHOAMI=$(whoami)

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

echo "Checking if magma has been installed"
MAGMA_INSTALLED=$(apt-cache show magma >  /dev/null 2>&1 echo "$SUCCESS_MESSAGE")
if [ "$MAGMA_INSTALLED" != "$SUCCESS_MESSAGE" ]; then
  echo "Magma not installed, processing installation"
  apt-get -y install curl make virtualenv zip rsync git software-properties-common python3-pip python-dev apt-transport-https

  alias python=python3
  pip3 install ansible==5.10.0

  echo "Generating localhost hostfile for Ansible"
  echo "[magma_deploy]
  127.0.0.1 ansible_connection=local" > $DEPLOY_PATH/agw_hosts

  # install magma and its dependencies including OVS.
  su - $MAGMA_USER -c "ansible-playbook -e \"MAGMA_ROOT='/home/$MAGMA_USER/magma' OUTPUT_DIR='/tmp'\" -i $DEPLOY_PATH/agw_hosts -e \"use_master=True\" $DEPLOY_PATH/magma_deploy.yml --extra-vars \"MAGMA_PACKAGE=\"$MAGMA_PACKAGE\"\""

  echo "Cleanup temp files"
  cd /root || exit
  rm -rf $AGW_INSTALL_CONFIG
  rm -f $DEPLOY_PATH/agw_hosts

else
  echo "Magma already installed, skipping.."
fi

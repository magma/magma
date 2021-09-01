#!/bin/bash
#
# Copyright 2020 The Magma Authors.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# This script is intended to install a docker-based gateway deployment

set -e

CWAG="cwag"
FEG="feg"
XWF="xwf"
INSTALL_DIR="/tmp/magmagw_install"
cd /opt/magma/

# TODO: Update docker-compose to stable version

DOCKER_COMPOSE_VERSION=1.29.1

DIR="."
echo "Setting working directory as: $DIR"
cd "$DIR"

if [ -z $1 ]; then
  echo "Please supply a gateway type to install. Valid types are: ['$FEG', '$CWAG', '$XWF']"
  exit
fi

GW_TYPE=$1
echo "Setting gateway type as: '$GW_TYPE'"

if [ "$GW_TYPE" != "$FEG" ] && [ "$GW_TYPE" != "$CWAG" ] && [ "$GW_TYPE" != "$XWF" ]; then
  echo "Gateway type '$GW_TYPE' is not valid. Valid types are: ['$FEG', '$CWAG', '$XWF']"
  exit
fi

# Ensure necessary files are in place
if [ ! -f .env ]; then
    echo ".env file is missing! Please add this file to the directory that you are running this command and re-try."
    exit
fi

if [ ! -f rootCA.pem ]; then
    echo "rootCA.pem file is missing! Please add this file to the directory that you are running this command and re-try."
    exit
fi

# TODO: Remove this once .env is used for control_proxy
if [ ! -f control_proxy.yml ]; then
    echo "control_proxy.yml file is missing! Please add this file to the directory that you are running this command and re-try."
    exit
fi

# Fetch files from github repo
rm -rf "$INSTALL_DIR"
mkdir -p "$INSTALL_DIR"

MAGMA_GITHUB_URL="{{ .Values.feg.repo.url }}"
git -C "$INSTALL_DIR" clone "$MAGMA_GITHUB_URL"

source .env
if [[ $IMAGE_VERSION == *"|"* ]]; then
  GIT_HASH=$(cut -d'|' -f2 <<< "$IMAGE_VERSION")
  IMAGE_VERSION=$(cut -d'|' -f1 <<< "$IMAGE_VERSION")
fi

if [ "$IMAGE_VERSION" != "latest" ]; then
    git -C $INSTALL_DIR/magma checkout "{{ .Values.feg.repo.branch }}"
fi

if [ "$GW_TYPE" == "$CWAG" ] || [ "$GW_TYPE" == "$XWF" ]; then
  MODULE_DIR="cwf"

  # Run CWAG ansible role to setup OVS
  echo "Copying and running ansible..."
  apt-add-repository -y ppa:ansible/ansible
  apt-get update -y
  apt-get -y install ansible
  ANSIBLE_CONFIG="$INSTALL_DIR"/magma/"$MODULE_DIR"/gateway/ansible.cfg ansible-playbook "$INSTALL_DIR"/magma/"$MODULE_DIR"/gateway/deploy/cwag.yml -i "localhost," -c local -v -e ingress_port="${INGRESS_PORT:-eth1}" -e uplink_ports="${UPLINK_PORTS:-eth2 eth3}" -e li_port="${LI_PORT:-eth4}"
fi

if [ "$GW_TYPE" == "$XWF" ]; then
  MODULE_DIR="xwf"
  CONNECTION_MODE=${MODE:=tcp}
  ANSIBLE_CONFIG="$INSTALL_DIR"/magma/"$MODULE_DIR"/gateway/ansible.cfg \
    ansible-playbook -e "xwf_ctrl_ip=$XWF_CTRL connection_mode=$CONNECTION_MODE" \
    "$INSTALL_DIR"/magma/"$MODULE_DIR"/gateway/deploy/xwf.yml -i "localhost," -c local -v
fi

if [ "$GW_TYPE" == "$FEG" ]; then
  MODULE_DIR="$GW_TYPE"

  # Load kernel module necessary for docker SCTP support
  sudo modprobe nf_conntrack_proto_sctp
  sudo tee -a /etc/modules <<< nf_conntrack_proto_sctp
fi

cp "$INSTALL_DIR"/magma/"$MODULE_DIR"/gateway/docker/docker-compose.yml .
cp "$INSTALL_DIR"/magma/orc8r/tools/docker/recreate_services.sh .
cp "$INSTALL_DIR"/magma/orc8r/tools/docker/recreate_services_cron .

# Install Docker
sudo apt-get update
sudo apt-get install -y \
    apt-transport-https \
    ca-certificates \
    curl \
    gnupg-agent \
    software-properties-common
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -
sudo add-apt-repository \
   "deb [arch=amd64] https://download.docker.com/linux/ubuntu \
   $(lsb_release -cs) \
   stable"
sudo apt-get update
sudo apt-get install -y docker-ce docker-ce-cli containerd.io

# Install Docker-Compose
sudo curl -L "https://github.com/docker/compose/releases/download/$DOCKER_COMPOSE_VERSION/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose

# Create snowflake to be mounted into containers
touch /etc/snowflake

echo "Placing configs in the appropriate place..."
mkdir -p /var/opt/magma
mkdir -p /var/opt/magma/configs
mkdir -p /var/opt/magma/certs
mkdir -p /etc/magma
mkdir -p /var/opt/magma/docker

#If this XWF installation copy the cwf config files as well
if [ "$GW_TYPE" == "$XWF" ]; then
  cp -TR "$INSTALL_DIR"/magma/cwf/gateway/configs /etc/magma
fi

# Copy default configs directory
cp -TR "$INSTALL_DIR"/magma/"$MODULE_DIR"/gateway/configs /etc/magma

# Copy config templates
cp -R "$INSTALL_DIR"/magma/orc8r/gateway/configs/templates /etc/magma

# Copy certs
cp rootCA.pem /var/opt/magma/certs/

# Copy control_proxy override
cp control_proxy.yml /var/opt/magma/configs/

# Copy docker files
cp docker-compose.yml /var/opt/magma/docker/
cp .env /var/opt/magma/docker/

# Copy recreate_services scripts to complete auto-upgrades
cp recreate_services.sh /var/opt/magma/docker/
cp recreate_services_cron /etc/cron.d/

# Copy DPI docker files
if [ "$GW_TYPE" == "$CWAG" ] && [ -f "$DPI_LICENSE_NAME" ]; then
  MODULE_DIR="cwf"
  mkdir -p "$SECRETS_VOLUME"
  cp "$INSTALL_DIR"/magma/"$MODULE_DIR"/gateway/docker/docker-compose-dpi.override.yml /var/opt/magma/docker/
  cp "$DPI_LICENSE_NAME" "$SECRETS_VOLUME"
fi

cd /var/opt/magma/docker

if [ ! -z "$DOCKER_USERNAME" ] && [ ! -z "$DOCKER_PASSWORD" ] && [ ! -z "$DOCKER_REGISTRY" ]; then
echo "Logging into docker registry at $DOCKER_REGISTRY"
docker login "$DOCKER_REGISTRY" --username "$DOCKER_USERNAME" --password "$DOCKER_PASSWORD"
fi
docker-compose pull
docker-compose -f docker-compose.yml up -d

# Pull and Run DPI container
if [ "$GW_TYPE" == "$CWAG" ] && [ -f "$DPI_LICENSE_NAME" ]; then
  cd /var/opt/magma/docker
  docker-compose -f docker-compose-dpi.override.yml pull
  docker-compose -f docker-compose-dpi.override.yml up -d
fi

echo "Installed successfully!!"

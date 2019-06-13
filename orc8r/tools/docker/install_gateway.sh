#!/bin/bash
#
# Copyright (c) 2016-present, Facebook, Inc.
# All rights reserved.
#
# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree. An additional grant
# of patent rights can be found in the PATENTS file in the same directory.

# This script is intended to install a docker-based gateway deployment

set -e

INSTALL_DIR="/tmp/magmagw_install"

DIR=$(dirname "$0")
if [ "$#" -eq 1 ]; then
  DIR=$1
fi
echo "Setting working directory as: $DIR"
cd "$DIR"

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

MAGMA_GITHUB_URL="https://github.com/facebookincubator/magma.git"
git -C "$INSTALL_DIR" clone "$MAGMA_GITHUB_URL"

# TODO: Add this back once this code is included in a github version
#TAG=$(git -C $INSTALL_DIR/magma tag | tail -1)
#git -C $INSTALL_DIR/magma checkout "tags/$TAG"

# Ensure this script hasn't changed
if ! cmp "$INSTALL_DIR"/magma/orc8r/tools/docker/install_gateway.sh install_gateway.sh; then
    echo "This 'install_gateway.sh' script has changed..."
    echo "Please copy this file from $INSTALL_DIR/magma/orc8r/tools/docker/install_gateway.sh and re-run"
    exit
fi

cp "$INSTALL_DIR"/magma/feg/gateway/docker/docker-compose.yml .
cp "$INSTALL_DIR"/magma/orc8r/tools/docker/upgrade_gateway.sh .

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
sudo usermod -aG docker "$SUDO_USER"

# Install Docker-Compose
sudo curl -L "https://github.com/docker/compose/releases/download/1.24.0/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose

# Create snowflake to be mounted into containers
touch /etc/snowflake

echo "Placing configs in the appropriate place..."
mkdir -p /var/opt/magma
mkdir -p /var/opt/magma/configs
mkdir -p /var/opt/magma/certs
mkdir -p /etc/magma
mkdir -p /var/opt/magma/docker

# Copy certs and configs
cp rootCA.pem /var/opt/magma/certs/
cp control_proxy.yml /etc/magma/
cp -r "$INSTALL_DIR"/magma/feg/gateway/configs /var/opt/magma/

# Copy docker files
cp docker-compose.yml /var/opt/magma/docker/
cp .env /var/opt/magma/docker/

# Copy upgrade script for future usage
cp upgrade_gateway.sh /var/opt/magma/docker/

cd /var/opt/magma/docker
source .env

echo "Logging into docker registry at $DOCKER_REGISTRY"
docker login "$DOCKER_REGISTRY"
docker-compose pull
docker-compose -f docker-compose.yml up -d

echo "Installed successfully!!"


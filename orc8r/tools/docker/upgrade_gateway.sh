#!/bin/bash
#
# Copyright (c) 2016-present, Facebook, Inc.
# All rights reserved.
#
# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree. An additional grant
# of patent rights can be found in the PATENTS file in the same directory.

# This script is intended to upgrade a docker-based gateway deployment
set -e

CWAG="cwag"
FEG="feg"
UPGRADE_DIR="/tmp/magmagw_install"

DIR="."
echo "Setting working directory as: $DIR"
cd "$DIR"

if [ -z $1 ]; then
  echo "Please supply a gateway type to upgrade. Valid types are: ['$FEG', '$CWAG']"
  exit
fi

GW_TYPE=$1
echo "Setting gateway type as: '$GW_TYPE'"

if [ "$GW_TYPE" != "$FEG" ] && [ "$GW_TYPE" != "$CWAG" ]; then
  echo "Gateway type '$GW_TYPE' is not valid. Valid types are: ['$FEG', '$CWAG']"
  exit
fi

if [ "$GW_TYPE" == "$CWAG" ]; then
  MODULE_DIR="cwf"
else
  MODULE_DIR=$GW_TYPE
fi

TAG=$2
echo "Using github tag: $TAG"

if [ -z $2 ]; then
    echo "Please provide a github tag as an argument"
    exit 1
fi
cd /var/opt/magma/docker

# Fetch files from github repo
rm -rf "$UPGRADE_DIR"
mkdir -p "$UPGRADE_DIR"

MAGMA_GITHUB_URL="https://github.com/facebookincubator/magma.git"
git -C "$UPGRADE_DIR" clone "$MAGMA_GITHUB_URL"
git -C "$UPGRADE_DIR"/magma checkout "tags/$TAG"

# Preserve control_proxy override
cp /etc/magma/control_proxy.yml ./

# Copy default configs directory
cp -TR "$UPGRADE_DIR"/magma/"$MODULE_DIR"/gateway/configs /etc/magma

# Copy config templates
cp -R "$UPGRADE_DIR"/magma/orc8r/gateway/configs/templates /etc/magma

# Move control_proxy override back
cp control_proxy.yml /etc/magma

# Update docker-compose file
cp "$UPGRADE_DIR"/magma/"$MODULE_DIR"/gateway/docker/docker-compose.yml /var/opt/magma/docker
source .env

echo "Logging into docker registry at $DOCKER_REGISTRY"
docker login "$DOCKER_REGISTRY"
docker-compose pull
docker-compose -f docker-compose.yml up --force-recreate -d

# Remove all stopped containers and dangling images
docker system prune -af

echo "Upgraded successfully!!"

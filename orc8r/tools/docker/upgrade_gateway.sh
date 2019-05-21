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

TAG=$1
if [ $# -eq 0 ]; then
    echo "Please provide a github tag as an argument"
    exit 1
fi
cd /var/opt/magma/docker

UPGRADE_DIR="/tmp/magmagw_install"

# Fetch files from github repo
rm -rf "$UPGRADE_DIR"
mkdir -p "$UPGRADE_DIR"

MAGMA_GITHUB_URL="https://github.com/facebookincubator/magma.git"
git -C "$UPGRADE_DIR" clone "$MAGMA_GITHUB_URL"
git -C "$UPGRADE_DIR"/magma checkout "tags/$TAG"

# Update docker-compose file
cp "$UPGRADE_DIR"/magma/feg/gateway/docker/docker-compose.yml /var/opt/magma/docker

source .env

echo "Logging into docker registry at $DOCKER_REGISTRY"
docker login "$DOCKER_REGISTRY"
docker-compose pull
docker-compose -f docker-compose.yml up --force-recreate -d

# Remove all stopped containers and dangling images
docker system prune -af

echo "Upgraded successfully!!"

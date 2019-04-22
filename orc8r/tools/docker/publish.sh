#!/bin/bash
#
# Copyright (c) 2016-present, Facebook, Inc.
# All rights reserved.
#
# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree. An additional grant
# of patent rights can be found in the PATENTS file in the same directory.

# This script pushes the docker images for orc8r/feg to a private registry.
# NOTE: Ensure that the image is built before this script is run.

set -e

usage() {
  echo "Usage: $0 -r REGISTRY -i IMAGE [-v VERSION] [-u USERNAME -p PASSFILE]"
  exit 2
}

# Parse the args and declare defaults
VERSION="latest"
while getopts 'r:i:v:u:p:h' OPT; do
  case "${OPT}" in
    r) REGISTRY=${OPTARG} ;;
    i) IMAGE=${OPTARG} ;;
    v) VERSION=${OPTARG} ;;
    u) USERNAME=${OPTARG} ;;
    p) PASSFILE=${OPTARG} ;;
    h|*) usage ;;
  esac
done

# Check if the required args are present
[ -z "${REGISTRY}" ] || [ -z "${IMAGE}" ] || [ -z "${VERSION}" ] && usage

export $(grep -v "#" .env | xargs)
PROJECT=${COMPOSE_PROJECT_NAME}

# Find the image ID for the latest build
IMAGE_ID=$(docker images "${PROJECT}_${IMAGE}:latest" --format "{{.ID}}")
if [ -z "${IMAGE_ID}" ]; then
  echo "Error! Missing image! Please build the image"
  exit 1
fi

echo "Pushing docker images for ${PROJECT}.. ${IMAGE}:${IMAGE_ID}"
echo "Logging into the docker registry..."
if [ -z "${USERNAME}" ]; then
  docker login "${REGISTRY}"
else
  [ -z "${USERNAME}" ] || [ -z "${PASSFILE}" ] && usage
  docker login "${REGISTRY}" -u "${USERNAME}" --password-stdin < "${PASSFILE}"
fi

# Tag and push the image
docker tag "${IMAGE_ID}" "${REGISTRY}/${IMAGE}:${VERSION}"
docker push "${REGISTRY}/${IMAGE}:${VERSION}"

echo "Image pushed succesfully!!"

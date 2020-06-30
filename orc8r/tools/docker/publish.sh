#!/usr/bin/env bash

# Copyright (c) 2016-present, Facebook, Inc.
# All rights reserved.
#
# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree. An additional grant
# of patent rights can be found in the PATENTS file in the same directory.

# publish.sh pushes a Docker image to a private registry.
# NOTE: ensure the image is built before running this script.

set -e -o pipefail

usage() {
  echo "Usage: $0 -r REGISTRY -i IMAGE [-v VERSION] [-u USERNAME -p PASSFILE]"
  exit 2
}

exitmsg() {
  echo "$1"
  exit 1
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
[[ -z "${REGISTRY}" ]] || [[ -z "${IMAGE}" ]] || [[ -z "${VERSION}" ]] && usage

# Find COMPOSE_PROJECT_NAME from environment or .env file
if [[ ${COMPOSE_PROJECT_NAME} == "" ]] && [[ -f .env ]]; then
  export $(grep -v "#" .env | xargs)
fi
# Exit if project name still empty
if [[ ${COMPOSE_PROJECT_NAME} == "" ]]; then
  exitmsg "[Error] project name cannot be empty: \
  set COMPOSE_PROJECT_NAME or add relevant entry to a .env file in the \
  working directory"
fi
PROJECT=${COMPOSE_PROJECT_NAME}

# Find the image ID for the latest build
DESIRED_IMAGE="${PROJECT}_${IMAGE}"
IMAGE_ID=$(docker images "${DESIRED_IMAGE}:latest" --format "{{.ID}}")
if [[ -z "${IMAGE_ID}" ]]; then
  exitmsg "[Error] project ${PROJECT} missing image ${DESIRED_IMAGE}: please build the image"
fi

echo "Pushing docker images for ${PROJECT}... ${IMAGE}:${IMAGE_ID}"
echo "Logging into the docker registry..."
if [[ -z "${USERNAME}" ]]; then
  docker login "${REGISTRY}"
else
  [[ -z "${USERNAME}" ]] || [[ -z "${PASSFILE}" ]] && usage
  docker login "${REGISTRY}" -u "${USERNAME}" --password-stdin < "${PASSFILE}"
fi

# Tag and push the image
docker tag "${IMAGE_ID}" "${REGISTRY}/${IMAGE}:${VERSION}"
docker push "${REGISTRY}/${IMAGE}:${VERSION}"

echo "Image pushed successfully"

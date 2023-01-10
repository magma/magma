#!/usr/bin/env bash

# Copyright 2020 The Magma Authors.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# tag-push-docker.sh tags and push specified images to a specific docker repository

set -e -o pipefail

function tag_and_push {
  docker tag "$IMAGE_ID" "${DOCKER_REGISTRY}/$IMAGE:$1"
  echo "Pushing ${DOCKER_REGISTRY}/$IMAGE:$1"
  docker push "${DOCKER_REGISTRY}/$IMAGE:$1"
}

usage() {
  echo "Usage: $0 [-i|--images 'image1|image2|image3'] [-t|--tag TAG] [-tl|--tag-latest BOOLEAN] [-p|--project PROJECT]"
  exit 2
}

exitmsg() {
  echo "$1"
  exit 1
}

# Parse the args
while [[ $# -gt 0 ]]
do
key="$1"
case $key in
    -i|--images)
    IMAGES="$2"
    shift  # pass argument or value
    ;;
    -t|--tag)
    TAG="$2"
    shift
    ;;
    -p|--project)
    PROJECT="$2"
    shift
    ;;
    -tl|--tag-latest)
    TAG_LATEST="$2"
    shift
    ;;
    -h|--help)
    usage
    shift
    ;;
    *)
    echo "Error: unknown cmdline option: $key"
    usage
    ;;
esac
shift  # past argument or value
done

# Define default values
TAG_LATEST="${TAG_LATEST:-true}"

if [[ -z $DOCKER_REGISTRY ]]; then
  exitmsg "Environment variable DOCKER_REGISTRY must be set"
fi

if [[ -z $DOCKER_USER ]]; then
  exitmsg "Environment variable DOCKER_USER must be set"
fi

if [[ -z $DOCKER_PASSWORD ]]; then
  exitmsg "Environment variable DOCKER_PASSWORD must be set"
fi

docker login "${DOCKER_REGISTRY}" -u "${DOCKER_USER}" -p "${DOCKER_PASSWORD}"
# shellcheck disable=SC2207
# Docker images does not contains special characters so we can skip the check
IMAGES_ARRAY=($(echo "$IMAGES" | tr "|" "\n"))
for IMAGE in "${IMAGES_ARRAY[@]}"; do
  IMAGE_TOSEARCH=$IMAGE
  if [ -n "${PROJECT}" ]; then
    IMAGE_TOSEARCH="${PROJECT}_${IMAGE}"
  fi
  IMAGE_ID=$(docker images "$IMAGE_TOSEARCH:latest" --format "{{.ID}}")
  tag_and_push "$TAG"
  if [ "$TAG_LATEST" = true ]; then
    tag_and_push "latest"
  fi
done

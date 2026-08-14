#!/bin/bash
#
# Copyright 2020 The Magma Authors.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License

DOCKER_DIR=${MAGMA_DOCKER_DIR:-/var/opt/magma/docker}

RUNNING_TAG=$(docker ps --filter name=magmad --format "{{.Image}}" | cut -d ":" -f 2)

source "$DOCKER_DIR/.env"

# If tag running is equal to .env, then do nothing
if [ "$RUNNING_TAG" == "$IMAGE_VERSION" ]; then
  exit
fi

if pidof -o %PPID -x $0 >/dev/null; then
  echo "Upgrade process already running"
  exit
fi

# Otherwise recreate containers with the new image
cd "$DOCKER_DIR" || exit

COMPOSE=(docker compose --compatibility -f docker-compose.yaml)

# Validate docker-compose file
CONFIG=$("${COMPOSE[@]}" config)
if [ -z "$CONFIG" ]; then
  echo "docker-compose.yaml is not valid"
  exit
fi

# Pull all images
OLD_IMAGES=$("${COMPOSE[@]}" images -q | sort -u)
[[ -z "$DOCKER_REGISTRY" ]] || "${COMPOSE[@]}" pull

# Stop and remove only containers that belong to this Compose project. Other
# workloads may share the Docker host and must not be interrupted by an AGW
# upgrade.
"${COMPOSE[@]}" down --remove-orphans

# Bring containers up
"${COMPOSE[@]}" up -d

# Remove old AGW images only when no container still references them. Docker
# refuses to remove an image that is in use, but checking first avoids noisy
# failures when an image is shared with another workload.
while IFS= read -r image; do
  [[ -z "$image" ]] && continue
  if ! docker ps -a -q --filter "ancestor=$image" | grep -q .; then
    docker image rm "$image" || true
  fi
done <<< "$OLD_IMAGES"

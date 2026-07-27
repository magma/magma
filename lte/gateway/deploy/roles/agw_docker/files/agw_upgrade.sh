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

readonly DOCKER_DIR=/var/opt/magma/docker

upgrade_agw() {
  local docker_dir=$1
  local docker_registry=''
  local image_version=''
  local key
  local value

  if [[ ! -r "$docker_dir/.env" ]]; then
    echo "Unable to read $docker_dir/.env" >&2
    return 1
  fi

  # Read only the values needed by this script. Docker Compose loads the same
  # .env file itself; executing it as shell code is unnecessary.
  while IFS='=' read -r key value; do
    case "$key" in
      DOCKER_REGISTRY) docker_registry=$value ;;
      IMAGE_VERSION) image_version=$value ;;
    esac
  done < "$docker_dir/.env"

  local running_tag
  running_tag=$(docker ps --filter name=magmad --format "{{.Image}}" | cut -d ":" -f 2)

  # If tag running is equal to .env, then do nothing
  if [ "$running_tag" == "$image_version" ]; then
    return 0
  fi

  if pidof -o %PPID -x "$0" >/dev/null; then
    echo "Upgrade process already running"
    return 0
  fi

  # Otherwise recreate containers with the new image
  cd "$docker_dir" || return 1

  local -a compose=(docker compose --compatibility -f docker-compose.yaml)

  # Validate docker-compose file
  local config
  config=$("${compose[@]}" config)
  if [ -z "$config" ]; then
    echo "docker-compose.yaml is not valid"
    return 1
  fi

  # Pull all images
  local old_images
  old_images=$("${compose[@]}" images -q | sort -u)
  [[ -z "$docker_registry" ]] || "${compose[@]}" pull

  # Stop and remove only containers that belong to this Compose project. Other
  # workloads may share the Docker host and must not be interrupted by an AGW
  # upgrade.
  "${compose[@]}" down --remove-orphans

  # Bring containers up
  "${compose[@]}" up -d

  # Remove old AGW images only when no container still references them. Docker
  # refuses to remove an image that is in use, but checking first avoids noisy
  # failures when an image is shared with another workload.
  local image
  while IFS= read -r image; do
    [[ -z "$image" ]] && continue
    if ! docker ps -a -q --filter "ancestor=$image" | grep -q .; then
      docker image rm "$image" || true
    fi
  done <<< "$old_images"
}

if [[ "${BASH_SOURCE[0]}" == "$0" ]]; then
  upgrade_agw "$DOCKER_DIR"
fi

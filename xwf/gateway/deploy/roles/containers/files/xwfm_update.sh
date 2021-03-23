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
RUNNING_TAG=$(docker ps --filter name=magmad --format "{{.Image}}" | cut -d ":" -f 2)

source /var/opt/magma/docker/.env

# If tag running is equal to .env, then do nothing
if [ "$RUNNING_TAG" == "$IMAGE_VERSION" ]; then
  exit
fi

if pidof -o %PPID -x $0 >/dev/null; then
  echo "Upgrade process already running"
  exit
fi

# Otherwise recreate containers with the new image
cd /var/opt/magma/docker || exit

# Validate docker-compose file
CONFIG=$(docker-compose -f docker-compose.yml config)
if [ -z "$CONFIG" ]; then
  echo "docker-compose.yml is not valid"
  exit
fi

# Pull all images
/usr/local/bin/docker-compose pull

# Stop all containers
# NOTE: add docker restart to work around GoRadius failures to stop
/usr/local/bin/docker-compose down ||
  systemctl restart docker.socket docker.service

CONTAINERS=$(docker ps -a -q)
[[ -z "$CONTAINERS" ]] || docker stop $CONTAINERS

# Bring containers up
/usr/local/bin/docker-compose up -d

# Remove all stopped containers and dangling images
docker system prune -af

# update image version in xwfmwhoami
sed -i "s/image_.*/image_version: $IMAGE_VERSION/g" /etc/xwfwhoami

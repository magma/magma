#!/bin/bash
# Copyright 2020 The Magma Authors.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# Create docker networks
docker network inspect s1&>/dev/null || {
docker network create --driver=bridge --subnet=192.168.60.0/24 --gateway=192.168.60.142 sgi && \
docker network create --driver=bridge --subnet=192.168.128.0/23 --gateway=192.168.129.1 s1 && \
docker inspect s1aptester&>/dev/null && docker rm s1aptester && \
DOCKER_BRIDGE_ID=`docker network inspect -f "{{ slice .Id 0 12 }}" s1`
ip r add 192.168.128.11/32 dev br-$DOCKER_BRIDGE_ID
ip r add 192.168.129.42/32 dev br-$DOCKER_BRIDGE_ID
}

# Temporarily create password for ubuntu user
echo "ubuntu:ubuntu" | chpasswd
sed -i 's/PasswordAuthentication no/PasswordAuthentication yes/' /etc/ssh/sshd_config
systemctl reload ssh

# Create containers and attach to s1aptester container
docker-compose --env-file .env -f docker-compose.yaml -f s1ap/docker-compose.s1ap.yaml up -d
docker attach s1aptester

# Remove temporary ubuntu user password
sed -i 's/PasswordAuthentication yes/PasswordAuthentication no/' /etc/ssh/sshd_config
systemctl reload ssh
passwd -d ubuntu

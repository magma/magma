#!/bin/bash
# Copyright 2020 The Magma Authors.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
SGI_INTERFACE=eth1
S1_INTERFACE=eth0
# Create docker networks
docker network inspect c_s1&>/dev/null || {
docker network create --driver=bridge a_bridge && \
docker network create --driver=bridge --subnet=192.168.60.0/24 --gateway=192.168.60.142 --opt com.docker.network.bridge.name=br-sgi b_sgi && \
docker network create --driver=bridge --subnet=192.168.128.0/23 --gateway=192.168.129.1 --opt com.docker.network.bridge.name=br-s1 c_s1 && \
#S1_DOCKER_BRIDGE_ID=`docker network inspect -f "{{ slice .Id 0 12 }}" c_s1` && \
ip r add 192.168.128.11/32 dev br-s1 && \
ip r add 192.168.129.42/32 dev br-s1
}
docker inspect s1aptester&>/dev/null && \
docker rm s1aptester

for i in $(grep $SGI_INTERFACE -rl /etc/magma); do sed -i s/$SGI_INTERFACE/br-sgi/ $i; done && \
for i in $(grep $S1_INTERFACE -rl /etc/magma); do sed -i s/$S1_INTERFACE/br-s1/ $i; done

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

for i in $(grep br-sgi -rl /etc/magma); do sed -i s/br-sgi/$SGI_INTERFACE/ $i; done && \
for i in $(grep br-s1 -rl /etc/magma); do sed -i s/br-s1/$S1_INTERFACE/ $i; done

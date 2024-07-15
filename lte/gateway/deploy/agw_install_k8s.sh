#!/bin/bash
# Copyright 2021 The Magma Authors.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

MAGMA_USER=magma

echo "Generate key for SSH access"
ssh-keygen -f -t rsa

apt-get install curl zip python3-pip net-tools sudo ca-certificates gnupg lsb-release sshpass -y

alias python=python3
# TODO GH13915 pinned for now because of breaking change in ansible-core 2.13.4
pip3 install ansible==5.0.1

echo "Installing kubernetes and configuring the nodes"
# su - $MAGMA_USER -c "ansible-playbook -v -e \"MAGMA_ROOT='/opt/magma' OUTPUT_DIR='/tmp'\" -i $DEPLOY_PATH/agw_k8s_hosts.ini $DEPLOY_PATH/magma_k8s.yml"
su $MAGMA_USER -c "ansible-playbook -v -e \"MAGMA_ROOT='/opt/magma' OUTPUT_DIR='/tmp'\" -i agw_k8s_hosts.ini magma_k8s.yml"

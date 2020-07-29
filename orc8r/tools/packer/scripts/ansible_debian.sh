#!/bin/bash -eux
################################################################################
# Copyright 2020 The Magma Authors.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
################################################################################

# Installs ansible on debian
bash -c 'echo -e "deb http://ppa.launchpad.net/ansible/ansible/ubuntu trusty main" > /etc/apt/sources.list.d/packages_ansible.list'

sudo apt-key adv --keyserver keyserver.ubuntu.com --recv-keys 93C4A3FD7BB9C367
sudo apt-get update
sudo apt-get install -y ansible

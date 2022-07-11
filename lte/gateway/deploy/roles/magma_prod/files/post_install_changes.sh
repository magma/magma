#!/bin/bash
# Copyright 2022 The Magma Authors.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

service magma@* stop
mkdir /var/opt/magma/configs
ln -s home/vagrant/magma/lte/gateway/configs/control_proxy.yml /var/opt/magma/configs/control_proxy.yml
cp /home/vagrant/magma/.cache/test_certs/rootCA.* /var/opt/magma/certs/

mv /etc/magma/pipelined.yml /etc/magma/pipelined.yml_BKP
mv /etc/magma/sessiond.yml /etc/magma/sessiond.yml_BKP
ln -s /home/vagrant/magma/lte/gateway/configs/pipelined.yml /etc/magma/pipelined.yml
ln -s /home/vagrant/magma/lte/gateway/configs/sessiond.yml /etc/magma/sessiond.yml

cp /home/vagrant/magma/lte/gateway/deploy/roles/dev_common/files/sshd_config /etc/ssh
service ssh reload

mkdir -p /home/vagrant/build/python/bin
touch /home/vagrant/build/python/bin/activate
ln -s /usr/bin/python3 /home/vagrant/build/python/bin/python3

service magma@magmad start


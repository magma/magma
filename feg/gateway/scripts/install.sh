#!/bin/bash
#
# Copyright 2020 The Magma Authors.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

set -e

DIR=$(dirname "$0")
if [ "$#" -eq 1 ]; then
  DIR=$1
fi
echo "Setting working directory as: $DIR"
cd "$DIR"

echo "Copying and running ansible..."
apt-add-repository -y ppa:ansible/ansible
apt-get update -y
apt-get -y install ansible
ansible-playbook ansible/main.yml -i "localhost," -c local -v

echo "Running apt autoremove..."
apt-get autoremove -y

echo "Installing python dependencies..."
pip3 install -r python/requires.txt

# Stop the services to avoid the textfile busy error
# NOTE: Please make sure the logic that follows the stop
# doesn't fail, since the failure aren't handled gracefully
echo "Stopping services..."
systemctl stop magma@*

echo "Copying python code..."
cp -TR python/lib/ /usr/local/lib/python3.5/dist-packages/
cp -TR python/scripts/ /usr/local/bin/

echo "Copying Go binaries and configs..."
mkdir -p /var/opt/magma/envdir
cp -TR bin/ /var/opt/magma/bin/
cp -TR certs/ /var/opt/magma/certs/
cp -TR config/ /etc/magma/

echo "Sarting magmad..."
systemctl daemon-reload
systemctl start magma@magmad
systemctl enable magma@magmad

echo "Installed successfully!!"

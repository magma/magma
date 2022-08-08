#!/bin/bash
################################################################################
# Copyright 2022 The Magma Authors.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
################################################################################

set -euo pipefail
set -x

export version="2.15.4-9-magma"
export DEBIAN_FRONTEND="noninteractive"

sudo rm -rf /etc/apt/sources.list.d/*.list
sudo rm -rf /var/cache/apt/

echo "Installing certificates ..."
sudo apt-get update
sudo apt-get install -y ca-certificates

echo "Adding artifactory ..."
echo 'deb https://artifactory.magmacore.org/artifactory/debian-test focal-ci main' | sudo tee /etc/apt/sources.list.d/magma.list

echo "Upgrading OVS packages ..."
sudo apt-get update
sudo apt-get install -y \
    libopenvswitch=$version \
    libopenvswitch-dev=$version \
    openvswitch-common=$version \
    openvswitch-datapath-dkms=$version \
    openvswitch-switch=$version \
    openvswitch-test=$version \
    python3-openvswitch=$version

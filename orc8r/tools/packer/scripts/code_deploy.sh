#!/bin/bash

# Copyright 2020 The Magma Authors.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# Don't error out if dpkg lock is held by someone else
function wait_for_lock() {
    while sudo fuser /var/lib/dpkg/lock >/dev/null 2>&1 ; do
        echo "\rWaiting for other software managers to finish...\n"
        sleep 1
    done
}

# Set up Code Deploy so we can do installs
wait_for_lock
apt-get -y update
wait_for_lock
apt-get -y install awscli libssl-dev ruby
wget https://aws-codedeploy-us-west-1.s3.amazonaws.com/latest/install
chmod +x ./install
./install auto

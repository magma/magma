#!/bin/bash
#
# Copyright 2021 The Magma Authors.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
PKGNAME=getenvoy-envoy-dev
WORK_DIR=/tmp/build-${PKGNAME}
set -e
SCRIPT_DIR="$(dirname "$(realpath "$0")")"
source "$SCRIPT_DIR/../lib/util.sh"

if_subcommand_exec

if [ -d "$WORK_DIR" ]; then
  rm -rf "$WORK_DIR"
fi
mkdir "$WORK_DIR"
cd "$WORK_DIR"

wget https://linuxfoundation.jfrog.io/artifactory/magma-packages-test/pool/focal-1.8.0/getenvoy-envoy_1.16.2.p0.ge98e41a-1p71.gbe6132a_amd64.deb
cp getenvoy-envoy_1.16.2.p0.ge98e41a-1p71.gbe6132a_amd64.deb $SCRIPT_DIR/../

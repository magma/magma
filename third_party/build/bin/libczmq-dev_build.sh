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

set -e
SCRIPT_DIR="$(dirname "$(realpath "$0")")"
source "$SCRIPT_DIR/../lib/util.sh"
PKGNAME=libczmq-dev
WORK_DIR=/tmp/build-${PKGNAME}

function buildrequires() {
  # Replace or remove packages that are unavailable in Ubuntu Focal
  echo  comerr-dev cppzmq-dev fontconfig-config fonts-dejavu-core icu-devtools krb5-multidev libaom3 libbrotli1 \
        libbsd-dev libc-dev-bin libc-devtools libc6-dev libcrypt-dev libde265-0 libdeflate0 libexpat1 libfontconfig1 libfreetype6 \
        libgd3 libgssrpc4 libicu-dev libicu66 libjpeg62-turbo libkrb5-dev libmd-dev libnorm-dev libnorm1 libnspr4 libnspr4-dev \
        libnss3 libnss3-dev libnuma1 libpgm-dev libpng16-16 libsodium-dev libsodium23 libsqlite3-0 libtiff5 libtirpc-dev libtirpc3 \
        libx11-6 libx11-data libxau6 libxcb1 libxdmcp6 libxml2 libxml2-dev libzmq3-dev libzmq5 \
        linux-libc-dev manpages manpages-dev rpcsvc-proto uuid-dev
}

if_subcommand_exec

# Ensure the package list is updated
apt-get update || echo "Warning: apt-get update failed, continuing with cached package list..."

if [ -d "$WORK_DIR" ]; then
  rm -rf "$WORK_DIR"
fi
mkdir "$WORK_DIR"
cd "$WORK_DIR"

# Download necessary .deb files for libczmq if not available in apt repositories
wget https://linuxfoundation.jfrog.io/artifactory/magma-packages-test/pool/focal-1.8.0/libczmq-dev_4.2.0-2_amd64.deb
wget http://ftp.us.debian.org/debian/pool/main/c/czmq/libczmq4_4.2.0-2_amd64.deb

cp libczmq-dev_4.2.0-2_amd64.deb $SCRIPT_DIR/../

# Install downloaded packages and handle any missing dependencies automatically
dpkg -i libczmq4_4.2.0-2_amd64.deb libczmq-dev_4.2.0-2_amd64.deb || \
  (apt-get install -f -y && dpkg -i libczmq4_4.2.0-2_amd64.deb libczmq-dev_4.2.0-2_amd64.deb)

# Re-run apt-get install to fix any other dependency issues if needed
apt-get install -y libnspr4-dev libnss3-dev || echo "Some dependencies might still be unavailable, check specific requirements."

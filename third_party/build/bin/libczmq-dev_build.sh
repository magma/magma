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


function buildrequires() {
  echo  comerr-dev cppzmq-dev fontconfig-config fonts-dejavu-core icu-devtools krb5-multidev libabsl20220623 libaom3 libavif15 libbrotli1 \ 
        libbsd-dev libbsd0 libc-dev-bin libc-devtools libc6-dev libcrypt-dev libdav1d6 libde265-0 libdeflate0 libexpat1 libfontconfig1 libfreetype6 \ 
        libgav1-1 libgd3 libgssrpc4 libheif1 libicu-dev libicu72 libjbig0 libjpeg62-turbo libkadm5clnt-mit12 libkadm5srv-mit12 \ 
        libkdb5-10 libkrb5-dev liblerc4 libmd-dev libnorm-dev libnorm1 libnsl-dev libnsl2 libnspr4 libnspr4-dev libnss3 libnss3-dev libnuma1 \ 
        libpgm-5.3-0 libpgm-dev libpng16-16 librav1e0 libsodium-dev libsodium23 libsqlite3-0 libsvtav1enc1 libtiff6 libtirpc-common libtirpc-dev \ 
        libtirpc3 libwebp7 libx11-6 libx11-data libx265-199 libxau6 libxcb1 libxdmcp6 libxml2 libxml2-dev libxpm4 libyuv0 libzmq3-dev libzmq5 \ 
        linux-libc-dev manpages manpages-dev rpcsvc-proto uuid-dev
}

if_subcommand_exec

if [ -d "$WORK_DIR" ]; then
  rm -rf "$WORK_DIR"
fi
mkdir "$WORK_DIR"
cd "$WORK_DIR"

wget https://ftp.debian.org/debian/pool/main/c/czmq/libczmq-dev_4.2.1-1_amd64.deb
wget http://ftp.us.debian.org/debian/pool/main/c/czmq/libczmq4_4.2.1-1_amd64.deb

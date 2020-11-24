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

# Generate the debian package from source for libfluid msg/base
# Example output:
#   magma-libfluid_0.1.0-1_amd64.deb

set -e

SCRIPT_DIR="$(dirname "$(realpath "$0")")"
source "${SCRIPT_DIR}"/../lib/util.sh

GIT_VERSION=0.1.0
ITERATION=1
PKGVERSION=${GIT_VERSION}.5
VERSION="${PKGVERSION}"-"${ITERATION}"
PKGNAME=magma-libfluid

function buildrequires() {
    echo g++ make libtool pkg-config libevent-dev libssl-dev
}

if_subcommand_exec

WORK_DIR=/tmp/build-${PKGNAME}

# Commit on the origin/0.2 branch, which has a lot of fixes
LIBFLUID_BASE_COMMIT=56df5e20c49387ab8e6b5cd363c6c10d309f263e
# Latest master commit with fixes passed v0.1.0
LIBFLUID_MSG_COMMIT=71a4fccdedfabece730082fbe87ef8ae5f92059f

# The resulting package is placed in $OUTPUT_DIR
# or in the cwd.
if [ -z "$1" ]; then
  OUTPUT_DIR=`pwd`
else
  OUTPUT_DIR=$1
  if [ ! -d "$OUTPUT_DIR" ]; then
    echo "error: $OUTPUT_DIR is not a valid directory. Exiting..."
    exit 1
  fi
fi

# build from source
if [ -d ${WORK_DIR} ]; then
  rm -rf ${WORK_DIR}
fi
mkdir ${WORK_DIR}
cd ${WORK_DIR}

# Clone repos and checkout latest commit
git clone https://github.com/OpenNetworkingFoundation/libfluid_base.git
git -C libfluid_base checkout $LIBFLUID_BASE_COMMIT

pushd libfluid_base
git apply "${PATCH_DIR}"/libfluid_base_patches/ExternalEventPatch.patch
popd

git clone https://github.com/OpenNetworkingFoundation/libfluid_msg.git
git -C libfluid_msg checkout $LIBFLUID_MSG_COMMIT

pushd libfluid_msg
git apply "${PATCH_DIR}"/libfluid_msg_patches/TunnelDstPatch.patch
git apply "${PATCH_DIR}"/libfluid_msg_patches/Add-support-for-setting-OVS-reg8.patch
popd

for repo in libfluid_base libfluid_msg
do
  cd $repo
  # Configure and compile
  ./autogen.sh
  ./configure --prefix=/usr
  make
  make install DESTDIR=${WORK_DIR}/install
  cd ../
done
# packaging
PKGFILE="$(pkgfilename)"
BUILD_PATH=${OUTPUT_DIR}/${PKGFILE}

# remove old packages
if [ -f ${BUILD_PATH} ]; then
  rm ${BUILD_PATH}
fi

fpm \
    -s dir \
    -t ${PKGFMT} \
    -a ${ARCH} \
    -n ${PKGNAME} \
    -v ${PKGVERSION} \
    --iteration ${ITERATION} \
    --provides ${PKGNAME} \
    --conflicts ${PKGNAME} \
    --replaces ${PKGNAME} \
    --package ${BUILD_PATH} \
    --depends "libevent-dev" \
    --depends "libssl-dev" \
    --description 'Libfluid Openflow Controller' \
    -C ${WORK_DIR}/install

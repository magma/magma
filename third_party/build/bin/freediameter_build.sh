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
SCRIPT_DIR="$(dirname "$(realpath "$0")")"
source "${SCRIPT_DIR}"/../lib/util.sh
GIT_URL=https://github.com/lionelgo/opencord.org.freeDiameter.git
GIT_COMMIT=13b0e7de0d66906d50e074a339f890d6e59813ad
PKGVERSION=0.0.1
ITERATION=1
VERSION="${PKGVERSION}"-"${ITERATION}"
PKGNAME=oai-freediameter
WORK_DIR=/tmp/build-${PKGNAME}

function buildrequires() {
    if [ "${OS_RELEASE}" == 'centos8' ]; then
        echo autoconf automake gcc git libgcrypt-devel libidn-devel bison-devel bison flex lksctp-tools-devel
    else
        echo autoconf automake build-essential git libgcrypt20-dev cmake libsctp-dev bison flex
    fi
}

function buildafter() {
    echo gnutls
}

if_subcommand_exec

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

git clone "${GIT_URL}" freediameter
cd freediameter
git checkout "${GIT_COMMIT}"

patch -p1 < $MAGMA_ROOT/lte/gateway/c/oai/patches/0001-opencoord.org.freeDiameter.patch

awk '{if (/^DISABLE_SCTP/) gsub(/OFF/, "ON"); print}' CMakeCache.txt > tmp && mv tmp CMakeCache.txt

mkdir build
cd build || exit 1

cmake -DCMAKE_C_FLAGS=-DOLD_SCTP_SOCKET_API=1 ..

make -j$(nproc)

make DESTDIR=${WORK_DIR}/install install

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
    --description 'freediameter' \
    -C ${WORK_DIR}/install

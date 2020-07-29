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

# Generate the debian package from source for nettle
# Example output:
#  oai-nettle_2.5-1_amd64.deb

set -e
SCRIPT_DIR="$(dirname "$(realpath "$0")")"
source "${SCRIPT_DIR}"/../lib/util.sh

PKGVERSION=2.5
ITERATION=1
VERSION="${PKGVERSION}"-"${ITERATION}"
PKGNAME=oai-nettle
WORK_DIR=/tmp/build-${PKGNAME}

function buildrequires() {
    echo autoconf automake build-essential libgmp-dev
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

wget https://ftp.gnu.org/gnu/nettle/nettle-$PKGVERSION.tar.gz
tar xf nettle-$PKGVERSION.tar.gz
cd nettle-$PKGVERSION/
./configure --disable-openssl --enable-shared --prefix=/usr --build=arm-linux-gnu
make -j`nproc`
make install DESTDIR=${WORK_DIR}/install/

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
    --description 'A low-level cryptographic library' \
    -C ${WORK_DIR}/install

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

# Generate the debian package from source for gnutls
# Example output:
#   oai-gnutls_3.1.23-1_amd64.deb

set -e
SCRIPT_DIR="$(dirname "$(realpath "$0")")"
source "${SCRIPT_DIR}"/../lib/util.sh

PKGVERSION=3.1.23
ITERATION=1
VERSION="${PKGVERSION}"-"${ITERATION}"
PKGNAME=oai-gnutls

function buildafter() {
    echo nettle
}

function buildrequires() {
    echo libtasn1-6-dev libp11-kit-dev \
         libtspi-dev libtspi1 libidn2-0-dev libidn11-dev
}

if_subcommand_exec

# continuing with main script

WORK_DIR=/tmp/build-${PKGNAME}

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
mkdir -p ${WORK_DIR}/install
cd ${WORK_DIR}

wget http://mirrors.dotsrc.org/gcrypt/gnutls/v3.1/gnutls-$PKGVERSION.tar.xz
tar xf gnutls-$PKGVERSION.tar.xz
cd gnutls-$PKGVERSION/
./configure --prefix=/usr
make -j`nproc`
make install DESTDIR=${WORK_DIR}/install/

# hotfix: this file conflicts with the nettle 2.5 package
rm -f ${WORK_DIR}/install/usr/share/info/dir

# packaging
BUILD_PATH=${OUTPUT_DIR}/"$(pkgfilename)"

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
    --depends "libtspi1" \
    --description 'GnuTLS is a secure communications library implementing the SSL, TLS and DTLS protocols and technologies around them.' \
    -C ${WORK_DIR}/install

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

# This script builds gRPC packages from upstream source code on github
#
# NOTE: build before installing protobuf packages
#
# example output:
#    grpc_1.0.0-2_amd64.deb

set -e
SCRIPT_DIR="$(dirname "$(realpath "$0")")"
source "${SCRIPT_DIR}"/../lib/util.sh

GIT_VERSION=1.15.0
ITERATION=3
VERSION="${GIT_VERSION}"-"${ITERATION}"
PKGNAME=grpc-dev

function buildrequires() {
    if [ "${PKGFMT}" == 'deb' ]; then
        echo build-essential autoconf libtool
    else
        echo autoconf gcc gcc-c++ libtool golang cmake protobuf-c-devel protobuf-c-compiler protobuf-devel
    fi
}
function installrequires() {
    if [ "${PKGFMT}" == 'deb' ]; then
        echo libgoogle-perftools4
    else
        echo gperftools
    fi
}

if_subcommand_exec

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

if [ -d ${WORK_DIR} ]; then
  rm -rf ${WORK_DIR}
fi

mkdir ${WORK_DIR}
echo /sbin/ldconfig > "${WORK_DIR}"/after_install.sh

cd ${WORK_DIR}
git clone https://github.com/grpc/grpc
cd grpc
git checkout -b v"${GIT_VERSION}" tags/v"${GIT_VERSION}"
git submodule update --init

# IMPORTANT: update prefix in Makefile
# change default prefix from /usr/local to /tmp/build-grpc-dev/install/usr/local
sed -i 's/.usr.local$/\/tmp\/build-grpc-dev\/install\/usr\/local/' Makefile

# build and install grpc
cd "${WORK_DIR}"/grpc
rm -rf "${WORK_DIR}"/grpc/build
mkdir build
cd build
cmake \
    -DgRPC_INSTALL=ON \
    -DBUILD_SHARED_LIBS=OFF \
    -DgRPC_BUILD_TESTS=OFF \
    -DgRPC_GFLAGS_PROVIDER=package \
    -DgRPC_PROTOBUF_PROVIDER=package \
    -DgRPC_ZLIB_PROVIDER=package \
    -DgRPC_SSL_PROVIDER=package \
    -DCMAKE_BUILD_TYPE=Release \
    ..

make -j$(nproc)
make DESTDIR="${WORK_DIR}"/install install

cp "${WORK_DIR}"/grpc/build/*.a "${WORK_DIR}"/install/usr/local/lib
cp "${WORK_DIR}"/grpc/build/grpc_* "${WORK_DIR}"/install/usr/local/bin

# HACK see https://github.com/grpc/grpc/issues/11868
# package still links to libgrpc++.so.1 even though libgrpc++.so.6 is needed
ln -sf ${WORK_DIR}/install/usr/local/lib/libgrpc++.so."${GIT_VERSION}" ${WORK_DIR}/install/usr/local/lib/libgrpc++.so.1

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
    -v ${GIT_VERSION} \
    --iteration ${ITERATION} \
    $(fpminstallrequires) \
    --provides ${PKGNAME} \
    --conflicts ${PKGNAME} \
    --replaces ${PKGNAME} \
    --package ${BUILD_PATH} \
    --after-install ${WORK_DIR}/after_install.sh \
    --description 'gRPC library' \
    -C ${WORK_DIR}/install

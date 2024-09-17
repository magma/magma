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

# Build prometheus-cpp packages from source
#
# Since there is no official release of prometheus-cpp, we use the unix
# timestamp as the version number of the lastest commit
# example output:
#    prometheus-cpp-dev_1.0.2-1485901529-d8326b2-1_amd64.deb

set -e
SCRIPT_DIR="$(dirname "$(realpath "$0")")"
source "${SCRIPT_DIR}"/../lib/util.sh

VERSION=1.0.2

GITCOMMIT=d8326b2bba945a435f299e7526c403d7a1f68c1f
VERSION="${VERSION}-$(echo "${GITCOMMIT}" | head -c 7)"

ITERATION=1
PKGNAME=prometheus-cpp-dev

function buildrequires() {
    echo build-essential autoconf libtool cmake libprotobuf-dev protobuf-compiler
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
cd ${WORK_DIR}
git clone https://github.com/jupp0r/prometheus-cpp
cd prometheus-cpp

# Master doesn't include all headers..
# use branch till its merged
# https://github.com/jupp0r/prometheus-cpp/pull/34
# 2019-08-14 pull request merged but deprecated last year...
git checkout "${GITCOMMIT}"

git submodule update --init

# Gen Makefile
cmake -DCMAKE_INSTALL_PREFIX=${WORK_DIR}/install/usr/local .

# build and install grpc
make -j$(nproc)
make install

# unix timestamp of the commit
COMMIT_TIMESTAMP=$(git show -s --format=%ct HEAD)
COMMIT_HASH=$(git rev-parse --short HEAD)

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
    -v ${VERSION} \
    --iteration ${ITERATION} \
    --provides ${PKGNAME} \
    --conflicts ${PKGNAME} \
    --replaces ${PKGNAME} \
    --package ${BUILD_PATH} \
    --description 'Prometheus C++ Library' \
    -C ${WORK_DIR}/install

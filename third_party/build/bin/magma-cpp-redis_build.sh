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

# Generate the debian package from source for cpp_redis
# Example output:
#   magma-cpp_redis.1.0-1_amd64.deb

set -e
SCRIPT_DIR="$(dirname "$(realpath "$0")")"
source "${SCRIPT_DIR}"/../lib/util.sh

GIT_VERSION=4.3.1
ITERATION=2
PKGVERSION=${GIT_VERSION}.1
VERSION="${PKGVERSION}"-"${ITERATION}"
PKGNAME=magma-cpp-redis

function buildrequires() {
    echo g++ make cmake libtool pkg-config
}

function buildafter() {
    echo magma-libtacopie
}

if_subcommand_exec

WORK_DIR=/tmp/build-${PKGNAME}

# The current release branch isn't building, so tag this commit
GIT_COMMIT=e7ef1f30eef8a073b0ff4ef90ad7e254167afc18

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

git clone https://github.com/Cylix/cpp_redis.git
cd cpp_redis
git checkout ${GIT_COMMIT}
git submodule init && git submodule update
mkdir build && cd build
cmake .. -DCMAKE_BUILD_TYPE=Release
make
make install DESTDIR=${WORK_DIR}/install

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
    --depends magma-libtacopie \
    --description 'Redis C++ client library' \
    -C ${WORK_DIR}/install

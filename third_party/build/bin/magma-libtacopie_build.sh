#!/bin/bash
################################################################################
# Copyright 2020 The Magma Authors.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
################################################################################
set -e
SCRIPT_DIR="$(dirname "$(realpath "$0")")"
source "${SCRIPT_DIR}"/../lib/util.sh

GIT_VERSION=3.2.0
ITERATION=1
PKGVERSION=${GIT_VERSION}.1
VERSION="${PKGVERSION}"-"${ITERATION}"
PKGNAME=magma-libtacopie

function buildrequires() {
    echo g++ make cmake libtool pkg-config
}

if_subcommand_exec

WORK_DIR=/tmp/build-${PKGNAME}

# The commit used in cpp_redis tacopie submodule at the version magma uses
GIT_COMMIT=8714fcec4ba9694fb00e83e788aadc099bd0fd5d

# Build time dependencies
BUILD_DEPS="g++ make libtool pkg-config"

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

git clone https://github.com/Cylix/tacopie.git
cd tacopie
git checkout ${GIT_COMMIT}

mkdir build && cd build
cmake .. -DCMAKE_BUILD_TYPE=Release -DCMAKE_INSTALL_PREFIX=${WORK_DIR}/install/usr/local
make
make install

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
    --description 'TacoPie C++ library' \
    -C ${WORK_DIR}/install

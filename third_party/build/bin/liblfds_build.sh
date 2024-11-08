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

# Generate the debian package from source for liblfds7.1.0
# Example output:
#   liblfds710_7.1.0-0_amd64.deb
set -ex

SCRIPT_DIR="$(dirname "$(realpath "$0")")"
source "${SCRIPT_DIR}"/../lib/util.sh

ITERATION=0
PKGVERSION=7.1.0
VERSION="${PKGVERSION}"-"${ITERATION}"
PKGNAME=liblfds710

if_subcommand_exec

WORK_DIR=/tmp/build-${PKGNAME}

# The resulting package is placed in $OUTPUT_DIR or in the cwd.
if [ -z "$1" ]; then
  OUTPUT_DIR=`pwd`
else
  OUTPUT_DIR=$1
  if [ ! -d "$OUTPUT_DIR" ]; then
    echo "error: $OUTPUT_DIR is not a valid directory. Exiting..."
    exit 1
  fi
fi

# Build from source
if [ -d ${WORK_DIR} ]; then
  rm -rf ${WORK_DIR}
fi
mkdir ${WORK_DIR}
cd ${WORK_DIR}

echo "Cloning the repository..."
git clone https://liblfds.org/git/liblfds
if [ $? -ne 0 ]; then
  echo "Error: Git clone failed. Exiting..."
  exit 1
fi
echo "Git clone successful"

cd liblfds/liblfds/liblfds7.1.0/liblfds710/build/gcc_gnumake

# Ensure that the build script exists and is executable
if [ ! -f ./bin/liblfds_build.sh ]; then
  echo "Error: liblfds_build.sh script not found. Exiting..."
  exit 1
fi

chmod +x ./bin/liblfds_build.sh

# Attempt to run the build script manually for debugging
echo "Running build script manually for debugging..."
./bin/liblfds_build.sh
if [ $? -ne 0 ]; then
  echo "Error: 'liblfds_build.sh' failed. Exiting..."
  exit 1
fi

echo "Running 'make so_vanilla'..."
make so_vanilla
if [ $? -ne 0 ]; then
  echo "Error: 'make so_vanilla' failed. Exiting..."
  exit 1
fi

LIB_DIR=/usr/local/lib
INC_DIR=/usr/local/include
mkdir -p ${WORK_DIR}/install$INC_DIR
mkdir -p ${WORK_DIR}/install$LIB_DIR

echo "Running 'make so_install'..."
make INSINCDIR="${WORK_DIR}/install/${INC_DIR}" INSLIBDIR="${WORK_DIR}/install/${LIB_DIR}" so_install
if [ $? -ne 0 ]; then
  echo "Error: 'make so_install' failed. Exiting..."
  exit 1
fi

# Packaging
PKGFILE="$(pkgfilename)"
BUILD_PATH=${OUTPUT_DIR}/${PKGFILE}

# Remove old packages
if [ -f ${BUILD_PATH} ]; then
  rm ${BUILD_PATH}
fi

echo "Running 'fpm' to create the package..."
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
    --description 'Lock-free data structure library' \
    -C ${WORK_DIR}/install
if [ $? -ne 0 ]; then
  echo "Error: 'fpm' command failed. Exiting..."
  exit 1
fi

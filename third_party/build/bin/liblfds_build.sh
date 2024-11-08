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
set -e
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
  OUTPUT_DIR=$(pwd)
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

# Attempt to clone the repository with error handling
REPO_URL="https://liblfds.org/git/liblfds"
RETRIES=3
for attempt in $(seq 1 $RETRIES); do
  echo "Attempting to clone repository (attempt $attempt/$RETRIES)..."
  if git clone "$REPO_URL"; then
    echo "Repository cloned successfully."
    break
  else
    echo "Error: Failed to clone repository. Attempt $attempt of $RETRIES."
    if [ $attempt -eq $RETRIES ]; then
      echo "Error: Could not clone repository after $RETRIES attempts. Exiting..."
      exit 1
    fi
    sleep 2  # Wait before retrying
  fi
done

cd liblfds/liblfds/liblfds7.1.0/liblfds710/build/gcc_gnumake

# Check if liblfds_build.sh exists and make it executable
if [ ! -f ./bin/liblfds_build.sh ]; then
  echo "Error: liblfds_build.sh script not found in ./bin directory."
  exit 1
fi
chmod +x ./bin/liblfds_build.sh

# Run the build script and capture output in a log
./bin/liblfds_build.sh > liblfds_build_output.log 2>&1
if [ $? -ne 0 ]; then
  echo "Error: liblfds_build.sh failed. Check liblfds_build_output.log for details."
  exit 1
fi

LIB_DIR=/usr/local/lib
INC_DIR=/usr/local/include
mkdir -p ${WORK_DIR}/install$INC_DIR
mkdir -p ${WORK_DIR}/install$LIB_DIR

make INSINCDIR="${WORK_DIR}/install/${INC_DIR}" INSLIBDIR="${WORK_DIR}/install/${LIB_DIR}" so_install 

# Packaging
PKGFILE="$(pkgfilename)"
BUILD_PATH=${OUTPUT_DIR}/${PKGFILE}

# Remove old packages
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
    --description 'Lock-free data structure library' \
    -C ${WORK_DIR}/install

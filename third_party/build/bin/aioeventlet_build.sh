#!/bin/bash
#
# Copyright 2021 The Magma Authors.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# build aioeventlet from locally patched aioeventlet lib
set -e
SCRIPT_DIR="$(dirname "$(realpath "$0")")"
source "$SCRIPT_DIR"/../lib/util.sh
PKGVERSION=0.5.1
ITERATION=2
VERSION="$PKGVERSION-$ITERATION"
PKGNAME=python3-aioeventlet
REPO="https://github.com/openstack-archive/deb-python-aioeventlet.git"
WORK_DIR=/tmp/build-${PKGNAME}

if [ -z "$1" ]; then
  OUTPUT_DIR=$(pwd)
else
  OUTPUT_DIR=$1
  if [ ! -d "$OUTPUT_DIR" ]; then
    echo "error: $OUTPUT_DIR is not a valid directory. Exiting..."
    exit 1
  fi
fi

# building from source
if [ -d ${WORK_DIR} ]; then
  rm -rf ${WORK_DIR}
fi
mkdir ${WORK_DIR}
cd ${WORK_DIR}

git clone ${REPO}
cd deb-python-aioeventlet
git apply ${PATCH_DIR}/*.patch

# packaging
PKGFILE="$(pkgfilename)"
BUILD_PATH="$OUTPUT_DIR"/"$PKGFILE"

# remove old packages
if [ -f "$BUILD_PATH" ]; then
  rm "$BUILD_PATH"
fi

fpm \
    -s python \
    -t "$PKGFMT" \
    -a "$ARCH" \
    -n "$PKGNAME" \
    -v "$PKGVERSION" \
    --iteration "$ITERATION" \
    --provides "$PKGNAME" \
    --conflicts "$PKGNAME" \
    --replaces "$PKGNAME" \
    --package "$BUILD_PATH" \
    --description 'patched aioeventlet' \
    --python-bin /usr/bin/python3 \
    --python-package-name-prefix 'python3' \
    ${WORK_DIR}/deb-python-aioeventlet/setup.py

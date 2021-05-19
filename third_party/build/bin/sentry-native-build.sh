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

set -e
SCRIPT_DIR="$(dirname "$(realpath "$0")")"
source "$SCRIPT_DIR/../lib/util.sh"
PKGVERSION=0.4.9
ITERATION=1
PKGNAME=sentry-native
WORK_DIR=/tmp/build-"$PKGNAME"
VERSION="$PKGVERSION-$ITERATION"

function buildrequires() {
    if [ "$OS_RELEASE" == 'centos8' ]; then
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
  OUTPUT_DIR=$(pwd)
else
  OUTPUT_DIR="$1"
  if [ ! -d "$OUTPUT_DIR" ]; then
    echo "error: $OUTPUT_DIR is not a valid directory. Exiting..."
    exit 1
  fi
fi

# build from source
if [ -d "$WORK_DIR" ]; then
  rm -rf "$WORK_DIR"
fi
mkdir "$WORK_DIR"
cd "$WORK_DIR"

wget https://github.com/getsentry/sentry-native/releases/download/"$PKGVERSION"/sentry-native.zip

unzip sentry-native.zip

cmake -B build -DCMAKE_BUILD_TYPE=RelWithDebInfo
cmake --build build --parallel
cmake --install build --prefix install --config RelWithDebInfo

cd build || exit 1

make DESTDIR="$WORK_DIR"/install install

# packaging
PKGFILE="$(pkgfilename)"
BUILD_PATH="$OUTPUT_DIR"/"$PKGFILE"

# remove old packages
if [ -f "$BUILD_PATH" ]; then
  rm "$BUILD_PATH"
fi

fpm \
    -s dir \
    -t "$PKGFMT" \
    -a "$ARCH" \
    -n "$PKGNAME" \
    -v "$PKGVERSION" \
    --iteration "$ITERATION" \
    --provides "$PKGNAME" \
    --conflicts "$PKGNAME" \
    --replaces "$PKGNAME" \
    --package "$BUILD_PATH" \
    --description 'sentry-native' \
    -C "$WORK_DIR"/install

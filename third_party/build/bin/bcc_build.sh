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
BCC_VER=23
SCRIPT_DIR="$(dirname "$(realpath "$0")")"
source "$SCRIPT_DIR/../lib/util.sh"
PKGVERSION="0.$BCC_VER"
ITERATION=1
PKGNAME=bcc-tools
WORK_DIR=/tmp/build-"$PKGNAME"
VERSION="$PKGVERSION-$ITERATION"

function buildrequires() {
    echo \
      bison \
      build-essential \
      cmake \
      flex \
      git \
      libedit-dev \
      libllvm7 \
      llvm-7-dev \
      libclang-7-dev \
      python3 \
      zlib1g-dev \
      libelf-dev \
      libfl-dev \
      python3-distutils \
      luajit \
      libluajit-5.1-dev
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

# https://github.com/iovisor/bcc/blob/master/INSTALL.md#ubuntu---source
git clone https://github.com/iovisor/bcc.git
mkdir bcc/build; cd bcc/build
git checkout "v0.$BCC_VER.0"
cmake ..
make
make install DESTDIR="$WORK_DIR"/install
cmake -DPYTHON_CMD=python3 .. # build python3 binding
pushd src/python/
make
make DESTDIR="$WORK_DIR"/install install
popd


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
    --description 'bcc tools' \
    -C "$WORK_DIR"/install

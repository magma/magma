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
ITERATION=1
VERSION="$PKGVERSION-$ITERATION"
PKGNAME=python3-aioeventlet

# packaging
OUTPUT_DIR=$(pwd)
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
    --description 'patched aioeventlet' \
    /usr/local/lib/python3.8/dist-packages/aioeventlet.py=/usr/local/lib/python3.8/dist-packages/aioeventlet.py \
    /usr/local/lib/python3.8/dist-packages/aioeventlet-0.5.1-py3.8.egg-info=/usr/local/lib/python3.8/dist-packages/aioeventlet-0.5.1-py3.8.egg-info

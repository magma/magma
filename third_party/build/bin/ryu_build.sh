#!/bin/bash
#
# Copyright 2022 The Magma Authors.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# package ryu with magma patches

set -e
SCRIPT_DIR="$(dirname "$(realpath "$0")")"
PATCHES_DIR="${SCRIPT_DIR}/../../../lte/gateway/deploy/roles/magma/files/patches"
# shellcheck source=/dev/null # relative path unknown and option -x not given
source "$SCRIPT_DIR"/../lib/util.sh
PKGVERSION=4.34
ITERATION=1
# shellcheck disable=SC2034 # variable is used in third_party/build/lib/util.sh
VERSION="$PKGVERSION-$ITERATION"
PKGNAME=python3-ryu
REPO="https://github.com/faucetsdn/ryu.git"
GIT_TAG="v${PKGVERSION}"
WORK_DIR=/tmp/build-${PKGNAME}

if_subcommand_exec

if [ -z "$1" ]; then
  OUTPUT_DIR=$(pwd)
else
  OUTPUT_DIR=$1
  if [ ! -d "$OUTPUT_DIR" ]; then
    echo "error: $OUTPUT_DIR is not a valid directory. Exiting..."
    exit 1
  fi
fi

# remove working directory from previous run
if [ -d ${WORK_DIR} ]; then
  rm -rf ${WORK_DIR}
fi
mkdir ${WORK_DIR}
cd ${WORK_DIR}

# cloning and patching
git clone -c advice.detachedHead=false --depth 1 --branch ${GIT_TAG} ${REPO}

patch -N -s -f ryu/ryu/ofproto/nx_actions.py <"${PATCHES_DIR}"/ryu_ipfix_args.patch 
patch -N -s -f ryu/ryu/app/ofctl/service.py <"${PATCHES_DIR}"/0001-Set-unknown-dpid-ofctl-log-to-debug.patch
patch -N -s -f ryu/ryu/ofproto/nicira_ext.py <"${PATCHES_DIR}"/0002-QFI-value-set-in-Openflow-controller-using-RYU.patch
patch -N -s -f ryu/ryu/ofproto/nx_match.py <"${PATCHES_DIR}"/0003-QFI-value-set-in-Openflow-controller-using-RYU.patch

# packaging
PKGFILE="$(pkgfilename)"
BUILD_PATH="$OUTPUT_DIR"/"$PKGFILE"

# remove old packages
if [[ -f "$BUILD_PATH" ]]; then
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
    --description 'patched ryu' \
    --python-bin /usr/bin/python3 \
    --python-package-name-prefix 'python3' \
    ${WORK_DIR}/ryu/setup.py

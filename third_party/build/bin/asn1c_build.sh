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

# Generate the debian package from source for asn1c
# The source code is a forked version hosted on the openair repo, which comes
# with OpenAirInterface (OAI) specific changes.
#
# Example output:
#   oai-asn1c_0~20160721+c3~r43c4a295-1_amd64.deb

set -e
SCRIPT_DIR="$(dirname "$(realpath "$0")")"
source "${SCRIPT_DIR}"/../lib/util.sh

COMMIT_DATE=20190423
# index of the commit from a particular date, start from 0
COMMIT_INDEX=0
COMMIT=f12568d6
# 0~ makes the version compatible with real version numbers
# 0~20160721+c3~r43c4a295 < 0~20160721+c5~r43c4a295 < 0~20160722+c0~r43c4a295 < 0.1
ITERATION=0
PKGVERSION=0~${COMMIT_DATE}+c${COMMIT_INDEX}~r${COMMIT}
VERSION="${PKGVERSION}"-"${ITERATION}"

PKGNAME=oai-asn1c

if_subcommand_exec

function configureopts() {
    if [ "${ARCH}" = "arm64" ]; then
        echo --build=arm-linux-gnu
    else
        echo -n
    fi
}

WORK_DIR=/tmp/build-${PKGNAME}

# The resulting package is placed in $OUTPUT_DIR
# or in the cwd.
if [ -z "$1" ]; then
  OUTPUT_DIR=$(pwd)
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
git clone https://gitlab.eurecom.fr/oai/asn1c.git
cd asn1c
git checkout ${COMMIT} .

autoreconf -iv

./configure "$(configureopts)"
make -j"$(nproc)"
make install DESTDIR=${WORK_DIR}/install/

# packaging
PKGFILE="$(pkgfilename)"
BUILD_PATH=${OUTPUT_DIR}/${PKGFILE}

# remove old packages
if [ -f "${BUILD_PATH}" ]; then
  rm "${BUILD_PATH}"
fi

fpm \
    -s dir \
    -t "${PKGFMT}" \
    -a "${ARCH}" \
    -n "${PKGNAME}" \
    -v "${PKGVERSION}" \
    --iteration "${ITERATION}" \
    --provides "${PKGNAME}" \
    --conflicts "${PKGNAME}" \
    --replaces "${PKGNAME}" \
    --package "${BUILD_PATH}" \
    -C ${WORK_DIR}/install \
    --description "ASN.1 (Release 15) compiler with OpenAirInterface (OAI)
    specific changes. ASN.1 to C compiler takes the ASN.1 module files (example)
    and generates the C++ compatible C source code. That code can be used to
    serialize the native C structures into compact and unambiguous
    BER/XER/PER-based data files, and deserialize the files back."

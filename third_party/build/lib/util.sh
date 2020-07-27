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

BUILDTS="$(date -u +"%Y%m%d%H%M%S")"

# IMPORTANT! set variables
# * PKGNAME
# * VERSION
# * OS_RELEASE (optional)
# * ARCH
# * PKGFMT
# to the correct values before calling pkgfilename()

ARCH="${ARCH:-amd64}"
PKGFMT="${PKGFMT:-deb}"
SCRIPTKEY="$(basename "$0" | sed 's/_build.sh//g')"
OUTPUT_DIR="$(pwd)"
PATCH_DIR="$(realpath "${SCRIPT_DIR}"/../patches)"/"${SCRIPTKEY}"
SUBCMD=
VERBOSE=
DEBUG=



while true; do
  case "$1" in
    -d | --debug ) DEBUG=true; shift ;;
    -v | --verbose ) VERBOSE=true; shift ;;

    -A | --build-after ) SUBCMD=buildafter; shift ;;
    -B | --build-requires ) SUBCMD=buildrequires; shift ;;
    -F | --package-file ) SUBCMD=pkgfilename; shift ;;

    -a | --arch ) ARCH="$2"; shift 2 ;;
    -f | --format ) PKGFMT="$2"; shift 2 ;;
    -r | --os-release ) OS_RELEASE="$2"; shift 2 ;;
    -o | --output-dir ) OUTPUT_DIR="$2"; shift 2 ;;
    -p | --patch-dir ) PATCH_DIR="$2"; shift 2 ;;
    -T | --timestamp ) BUILDTS="$2"; shift 2 ;;
    -- ) shift; break ;;
    * ) break ;;
  esac
done

if [ ! -d "$OUTPUT_DIR" ]; then
    echo "error: $OUTPUT_DIR is not a valid directory. Exiting..."
    exit 1
fi

function if_subcommand_exec() {
    # invoke this in build scripts after defining functions referenced here
    case "${SUBCMD}" in
        buildafter | buildrequires | pkgfilename) "${SUBCMD}" ; exit $? ;;
        * ) ;;
    esac

}

function pkgfilename() {
    if [ "${PKGFMT}" == 'deb' ]; then
        echo "${PKGNAME}_${VERSION}${OS_RELEASE}_${ARCH}.${PKGFMT}"
    elif [ "${PKGFMT}" == 'rpm' ]; then
        if [ "x${OS_RELEASE}" == 'x' ]; then
            echo "${PKGNAME}-${VERSION}.${ARCH}.${PKGFMT}"
        else
            echo "${PKGNAME}-${VERSION}.${OS_RELEASE}.${ARCH}.${PKGFMT}"
        fi
    fi
}

function buildafter() {
    # print a list of packages that must be built and installed before this one
    echo -n
}

function buildrequires() {
    # print a list of packages that must be installed from configured sources before build
    echo -n
}

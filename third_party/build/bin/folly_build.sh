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

# Generate the debian package from source for folly
# Example output:
#   folly_0.1.0-1_amd64.deb

set -e
SCRIPT_DIR="$(dirname "$(realpath "$0")")"
source "${SCRIPT_DIR}"/../lib/util.sh

GIT_VERSION=2018.02.26.00
ITERATION=6
PKGNAME=libfolly-dev
VERSION="${GIT_VERSION}-${ITERATION}"

LIBGFLAGS=libgflags2v5
BOOST_VERSION=1.62.0
LIBEVENT_VERSION=2.0-5
SSL_VERSION=1.1
LIBICU=libicu57
DOUBLE_CONVERSION_VERSION=1
JEMALLOC_VERSION=1

if [ "${OS_RELEASE}" = 'ubuntu18.04' ]; then
    LIBGFLAGS=libgflags2.2
    LIBEVENT_VERSION=2.1-6
    LIBICU=libicu60
elif [ "${OS_RELEASE}" = 'ubuntu20.04' ]; then
    BOOST_VERSION=1.71.0
    LIBGFLAGS=libgflags2.2
    LIBEVENT_VERSION=2.1-7
    LIBICU=libicu66
    DOUBLE_CONVERSION_VERSION=3
    JEMALLOC_VERSION=2
fi

function buildrequires() {
    echo \
        libboost-all-dev \
        libevent-dev \
        libdouble-conversion-dev \
        libgoogle-glog-dev \
        libgflags-dev \
        libiberty-dev \
        liblz4-dev \
        liblzma-dev \
        libsnappy-dev \
        zlib1g-dev \
        binutils-dev \
        libjemalloc-dev \
        libssl-dev \
        pkg-config
}

if_subcommand_exec

WORK_DIR=/tmp/build-${PKGNAME}

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
mkdir ${WORK_DIR}

# post-install script
echo sudo /sbin/ldconfig > "${WORK_DIR}"/after_install.sh

cd ${WORK_DIR}
DESTDIR=${WORK_DIR}/install
DEBUGDIR=${DESTDIR}/usr/local/lib/debug

# Get the code
git clone https://github.com/facebook/folly.git
cd folly
git checkout v${GIT_VERSION}

# build folly
cd folly

if [ "${ARCH}" = "arm64" ]; then
patch Makefile.am <<'EOF'
@@ -733,7 +733,12 @@ libfolly_la_SOURCES += \
 endif

 libfollybasesse42_la_LDFLAGS = $(AM_LDFLAGS) -version-info $(LT_VERSION)
+
+if HAVE_X86_64
 libfollybasesse42_la_CXXFLAGS = -msse4.2 -mpclmul
+else
+libfollybasesse42_la_CXXFLAGS =
+endif

 libfollybase_la_LIBADD = libfollybasesse42.la
 libfollybase_la_LDFLAGS = $(AM_LDFLAGS) -version-info $(LT_VERSION)735a736,737
EOF
fi

autoreconf -ivf
./configure
make -j $(nproc)
make install DESTDIR=${DESTDIR}

# Only include stripped .so files
[ -d "$DEBUGDIR/usr/local/lib" ] || mkdir -p "$DEBUGDIR/usr/local/lib"
find "$DESTDIR/usr/local/lib" -maxdepth 1 -iname "lib*.so.*" -type f \
  -execdir objcopy --only-keep-debug {} "$DEBUGDIR/usr/local/lib/{}.debug" \; \
  -execdir strip --strip-debug --strip-unneeded {} \; \
  -execdir objcopy --add-gnu-debuglink "$DEBUGDIR/usr/local/lib/{}.debug" {} \;

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
    -v ${GIT_VERSION} \
    -C ${DESTDIR} \
    --iteration ${ITERATION} \
    --description "Facebook Folly C++ Library" \
    --provides ${PKGNAME} \
    --conflicts ${PKGNAME} \
    --replaces ${PKGNAME} \
    --package ${BUILD_PATH} \
    --after-install "${WORK_DIR}"/after_install.sh \
    --depends libc6 \
    --depends libstdc++6 \
    --depends libboost-context"$BOOST_VERSION" \
    --depends libboost-filesystem"$BOOST_VERSION" \
    --depends libboost-program-options"$BOOST_VERSION" \
    --depends libboost-regex"$BOOST_VERSION" \
    --depends libboost-system"$BOOST_VERSION" \
    --depends libboost-thread"$BOOST_VERSION" \
    --depends libdouble-conversion"${DOUBLE_CONVERSION_VERSION}" \
    --depends libevent-"$LIBEVENT_VERSION" \
    --depends "${LIBGFLAGS}" \
    --depends libgoogle-glog0v5 \
    --depends "${LIBICU}" \
    --depends libjemalloc"${JEMALLOC_VERSION}" \
    --depends liblz4-1 \
    --depends liblzma5 \
    --depends libsnappy1v5 \
    --depends libssl"$SSL_VERSION" \
    --depends zlib1g \
    --exclude usr/local/lib/debug \
    --exclude usr/local/lib/*.a \
    --exclude usr/local/lib/*.la \
    --exclude usr/local/lib/*init* \
    --exclude usr/local/lib/*logging* \
    --exclude usr/local/lib/*benchmark* \
    usr/local/lib \
    usr/local/include

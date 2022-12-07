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

# Generate the debian package from source for gnutls
# Example output:
#   oai-gnutls_3.1.23-1_amd64.deb

set -e
SCRIPT_DIR="$(dirname "$(realpath "$0")")"
source "${SCRIPT_DIR}"/../lib/util.sh

PKGVERSION=3.1.23
ITERATION=1
VERSION="${PKGVERSION}"-"${ITERATION}"
PKGNAME=oai-gnutls

function buildafter() {
    echo nettle
}

function buildrequires() {
    echo \
        libtasn1-6-dev \
        libp11-kit-dev \
        libtspi-dev \
        libtspi1 \
        libidn2-0-dev \
        libidn11-dev
}

if_subcommand_exec

# continuing with main script

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
mkdir -p ${WORK_DIR}/install
cd ${WORK_DIR}

wget http://mirrors.dotsrc.org/gcrypt/gnutls/v3.1/gnutls-$PKGVERSION.tar.xz
tar xf gnutls-$PKGVERSION.tar.xz
cd gnutls-$PKGVERSION/

# These patches account for the bug reported in 
# https://lists.gnu.org/r/bug-gnulib/2018-03/msg00000.html.
# This bug is still present in the gnutls version used here.
# The underlying gnulib library was patched in 
# https://github.com/coreutils/gnulib/commit/4af4a4a71827c0bc5e0ec67af23edef4f15cee8e.
patch gl/stdio-impl.h <<'EOF'
@@ -18,6 +18,9 @@
    the same implementation of stdio extension API, except that some fields
    have different naming conventions, or their access requires some casts.  */
 
+#if !defined _IO_IN_BACKUP && defined _IO_EOF_SEEN
+# define _IO_IN_BACKUP 0x100
+#endif
 
 /* BSD stdio derived implementations.  */
 
EOF

patch gl/fseterr.c <<'EOF'
@@ -29,7 +29,7 @@
   /* Most systems provide FILE as a struct and the necessary bitmask in
      <stdio.h>, because they need it for implementing getc() and putc() as
      fast macros.  */
-#if defined _IO_ftrylockfile || __GNU_LIBRARY__ == 1 /* GNU libc, BeOS, Haiku, Linux libc5 */
+#if defined _IO_EOF_SEEN || __GNU_LIBRARY__ == 1 /* GNU libc, BeOS, Haiku, Linux libc5 */
   fp->_flags |= _IO_ERR_SEEN;
 #elif defined __sferror || defined __DragonFly__ /* FreeBSD, NetBSD, OpenBSD, DragonFly, Mac OS X, Cygwin */
   fp_->_flags |= __SERR;
EOF

./configure --prefix=/usr
make -j`nproc`
make install DESTDIR=${WORK_DIR}/install/

# hotfix: this file conflicts with the nettle 2.5 package
rm -f ${WORK_DIR}/install/usr/share/info/dir

# packaging
BUILD_PATH=${OUTPUT_DIR}/"$(pkgfilename)"

# remove old packages
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
    --depends "libtspi1" \
    --description 'GnuTLS is a secure communications library implementing the SSL, TLS and DTLS protocols and technologies around them.' \
    -C ${WORK_DIR}/install

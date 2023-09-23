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

# This script builds gRPC packages from upstream source code on github
#
# NOTE: build before installing protobuf packages
#
# example output:
#    grpc_1.0.0-2_amd64.deb

set -e
SCRIPT_DIR="$(dirname "$(realpath "$0")")"
source "${SCRIPT_DIR}"/../lib/util.sh

GIT_VERSION=1.15.0
ITERATION=3
VERSION="${GIT_VERSION}"-"${ITERATION}"
PKGNAME=grpc-dev

function buildrequires() {
    echo \
        build-essential \
        autoconf \
        libtool
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

if [ -d ${WORK_DIR} ]; then
  rm -rf ${WORK_DIR}
fi

mkdir ${WORK_DIR}
echo /sbin/ldconfig > "${WORK_DIR}"/after_install.sh

cd ${WORK_DIR}
git clone https://github.com/grpc/grpc
cd grpc
git checkout -b v"${GIT_VERSION}" tags/v"${GIT_VERSION}"
git submodule update --init

# Patches to avoid ambiguous declaration of `getid()` in grpc and glibc, see
# https://sourceware.org/git/?p=glibc.git;a=commit;h=1d0fc213824eaa2a8f8c4385daaa698ee8fb7c92
# and the respective commit to the grpc repository
# https://github.com/grpc/grpc/commit/57586a1ca7f17b1916aed3dea4ff8de872dbf853
patch src/core/lib/gpr/log_linux.cc <<'EOF'
@@ -40,7 +40,7 @@
 #include <time.h>
 #include <unistd.h>

-static long gettid(void) { return syscall(__NR_gettid); }
+static long sys_gettid(void) { return syscall(__NR_gettid); }

 void gpr_log(const char* file, int line, gpr_log_severity severity,
              const char* format, ...) {
@@ -70,7 +70,7 @@ void gpr_default_log(gpr_log_func_args* args) {
   gpr_timespec now = gpr_now(GPR_CLOCK_REALTIME);
   struct tm tm;
   static __thread long tid = 0;
-  if (tid == 0) tid = gettid();
+  if (tid == 0) tid = sys_gettid();

   timer = static_cast<time_t>(now.tv_sec);
   final_slash = strrchr(args->file, '/');
EOF

patch src/core/lib/gpr/log_posix.cc <<'EOF'
@@ -30,7 +30,7 @@
 #include <string.h>
 #include <time.h>

-static intptr_t gettid(void) { return (intptr_t)pthread_self(); }
+static intptr_t sys_gettid(void) { return (intptr_t)pthread_self(); }

 void gpr_log(const char* file, int line, gpr_log_severity severity,
              const char* format, ...) {
@@ -86,7 +86,7 @@ void gpr_default_log(gpr_log_func_args* args) {
   char* prefix;
   gpr_asprintf(&prefix, "%s%s.%09d %7" PRIdPTR " %s:%d]",
                gpr_log_severity_string(args->severity), time_buffer,
-               (int)(now.tv_nsec), gettid(), display_file, args->line);
+               (int)(now.tv_nsec), sys_gettid(), display_file, args->line);

   fprintf(stderr, "%-70s %s\n", prefix, args->message);
   gpr_free(prefix);
EOF

patch src/core/lib/iomgr/ev_epollex_linux.cc <<'EOF'
@@ -1146,7 +1146,7 @@
 }

 #ifndef NDEBUG
-static long gettid(void) { return syscall(__NR_gettid); }
+static long sys_gettid(void) { return syscall(__NR_gettid); }
 #endif

 /* pollset->mu lock must be held by the caller before calling this.
@@ -1166,7 +1166,7 @@
 #define WORKER_PTR (&worker)
 #endif
 #ifndef NDEBUG
-  WORKER_PTR->originator = gettid();
+  WORKER_PTR->originator = sys_gettid();
 #endif
   if (GRPC_TRACE_FLAG_ENABLED(grpc_polling_trace)) {
     gpr_log(GPR_INFO,
EOF

# Disable flag that treats warnings as errors
# see https://github.com/grpc/grpc/issues/17287#issuecomment-444876748
sed -i 's/-Werror/ /g' Makefile

# IMPORTANT: update prefix in Makefile
# change default prefix from /usr/local to /tmp/build-grpc-dev/install/usr/local
sed -i 's/.usr.local$/\/tmp\/build-grpc-dev\/install\/usr\/local/' Makefile

# build and install grpc
make -j$(nproc)
make install

# HACK see https://github.com/grpc/grpc/issues/11868
# package still links to libgrpc++.so.1 even though libgrpc++.so.6 is needed
ln -sf ${WORK_DIR}/install/usr/local/lib/libgrpc++.so."${GIT_VERSION}" ${WORK_DIR}/install/usr/local/lib/libgrpc++.so.1

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
    --iteration ${ITERATION} \
    --depends "libgoogle-perftools4" \
    --provides ${PKGNAME} \
    --conflicts ${PKGNAME} \
    --replaces ${PKGNAME} \
    --package ${BUILD_PATH} \
    --after-install ${WORK_DIR}/after_install.sh \
    --description 'gRPC library' \
    -C ${WORK_DIR}/install

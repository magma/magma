#!/bin/bash
#
# Copyright (c) 2016-present, Facebook, Inc.
# All rights reserved.
#
# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree. An additional grant
# of patent rights can be found in the PATENTS file in the same directory.

# Generate the debian package from source for nghttpx
# Example output:
#   magma-nghttpx_1.0.1-1_amd64.deb


set -e
BUILD_DATE=`date -u +"%Y%m%d%H%M%S"`
GIT_VERSION=1.31.1
VERSION=${GIT_VERSION}
ITERATION=1
ARCH=amd64
PKGFMT=deb
PKGNAME=magma-nghttpx
WORK_DIR=/tmp/build-${PKGNAME}
SCRIPT_DIR=$(dirname `realpath $0`)
POSTINST=${SCRIPT_DIR}/nghttpx-postinst

# Build time dependencies
BUILD_DEPS="g++ make binutils autoconf automake autotools-dev \
    libtool pkg-config zlib1g-dev libcunit1-dev libssl-dev libxml2-dev \
    libev-dev libevent-dev libjansson-dev libjemalloc-dev cython python3-dev \
    python-setuptools libc-ares-dev git ruby bison"

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

# install build time dependency
sudo apt-get install -y $BUILD_DEPS

# build from source
if [ -d ${WORK_DIR} ]; then
  rm -rf ${WORK_DIR}
fi
mkdir ${WORK_DIR}
cd ${WORK_DIR}

# Get the code
git clone https://github.com/nghttp2/nghttp2.git
cd nghttp2
git checkout v${GIT_VERSION}

# Apply the nghttpx patches
# Command to generate patch after creating a commit:
# git format-patch HEAD~ > out.patch
patch_dir="${SCRIPT_DIR}/nghttpx_patches"
git apply ${patch_dir}/${GIT_VERSION}/0001-patch.patch

# Compile
git submodule update --init
autoreconf -i
automake
autoconf
./configure --with-mruby --enable-app --disable-examples --disable-python-bindings
make

# Install the binaries
make install DESTDIR=${WORK_DIR}/install/

# packaging
PKGFILE=${PKGNAME}_${VERSION}-${ITERATION}_${ARCH}.${PKGFMT}
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
    -v ${VERSION} \
    --iteration ${ITERATION} \
    --provides ${PKGNAME} \
    --conflicts ${PKGNAME} \
    --replaces ${PKGNAME} \
    --package ${BUILD_PATH} \
    --depends "libjansson-dev" \
    --depends "libjemalloc-dev" \
    --depends "libssl-dev" \
    --depends "libevent-dev" \
    --depends "libev-dev" \
    --depends "libc-ares-dev" \
    --description 'Nghttp2 http/2 library and nghttpx proxy' \
    --after-install ${POSTINST} \
    -C ${WORK_DIR}/install

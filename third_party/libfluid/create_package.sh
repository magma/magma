#!/bin/bash
#
# Copyright (c) 2016-present, Facebook, Inc.
# All rights reserved.
#
# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree. An additional grant
# of patent rights can be found in the PATENTS file in the same directory.

# Generate the debian package from source for libfluid msg/base
# Example output:
#   magma-libfluid_0.1.0-1_amd64.deb
set -e
BUILD_DATE=`date -u +"%Y%m%d%H%M%S"`
GIT_VERSION=0.1.0
VERSION=${GIT_VERSION}.4
ITERATION=1
ARCH=amd64
PKGFMT=deb
PKGNAME=magma-libfluid
WORK_DIR=/tmp/build-${PKGNAME}
SCRIPT_DIR=$(dirname `realpath $0`)
# Commit on the origin/0.2 branch, which has a lot of fixes
LIBFLUID_BASE_COMMIT=56df5e20c49387ab8e6b5cd363c6c10d309f263e
# Latest master commit with fixes passed v0.1.0
LIBFLUID_MSG_COMMIT=71a4fccdedfabece730082fbe87ef8ae5f92059f

# Build time dependencies
BUILD_DEPS="g++ make libtool pkg-config libevent-dev libssl-dev"

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

# Clone repos and checkout latest commit
git clone https://github.com/OpenNetworkingFoundation/libfluid_base.git
git -C libfluid_base checkout $LIBFLUID_BASE_COMMIT

git clone https://github.com/OpenNetworkingFoundation/libfluid_msg.git
git -C libfluid_msg checkout $LIBFLUID_MSG_COMMIT

for repo in libfluid_base libfluid_msg
do
  cd $repo
  patch_files=${SCRIPT_DIR}/${repo}_patches/*
  for patch in $patch_files
  do
    git apply $patch
  done
  # Configure and compile
  ./autogen.sh
  ./configure --prefix=/usr
  make
  sudo make install DESTDIR=${WORK_DIR}/install
  cd ../
done
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
    --depends "libevent-dev" \
    --depends "libssl-dev" \
    --description 'Libfluid Openflow Controller' \
    -C ${WORK_DIR}/install

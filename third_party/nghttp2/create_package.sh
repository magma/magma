#!/bin/bash
#
# Copyright (c) 2016-present, Facebook, Inc.
# All rights reserved.
#
# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree. An additional grant
# of patent rights can be found in the PATENTS file in the same directory.

# Generate nghttp2 deb packages based on the debian sid release branch
# Example output:
#   libnghttp2-14_1.17.0-1_amd64.deb
#   libnghttp2-dev_1.17.0-1_amd64.deb
#   libnghttp2-doc_1.17.0-1_all.deb
#   nghttp2_1.17.0-1_all.deb
#   nghttp2-client_1.17.0-1_amd64.deb
#   nghttp2-proxy_1.17.0-1_amd64.deb
#   nghttp2-server_1.17.0-1_amd64.deb

set -e
PKGNAME=nghttp2
VERSION=1.17.0 # more recent 1.18 requires debhelper >= 10 that is not available in jessie
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
cd ${WORK_DIR}

# Build time requirements reported by 'dpkg-buildpackage -b -us -uc'.
sudo apt-get install debhelper libc-ares-dev python-sphinx dh-systemd 

# packages not included by jessie. They are pulled from debian unstable and hosted by us
sudo apt-get install libspdylay7 libspdylay-dev

wget http://http.debian.net/debian/pool/main/n/nghttp2/${PKGNAME}_${VERSION}.orig.tar.bz2
wget http://http.debian.net/debian/pool/main/n/nghttp2/${PKGNAME}_${VERSION}-1.debian.tar.xz
tar xf ${PKGNAME}_${VERSION}.orig.tar.bz2
tar xf ${PKGNAME}_${VERSION}-1.debian.tar.xz
mv debian ${PKGNAME}-${VERSION}
cd ${PKGNAME}-${VERSION}/

# packaging
DEB_BUILD_OPTIONS=nocheck dpkg-buildpackage -b -us -uc

# move packages to output dir
mv ../*.deb ${OUTPUT_DIR}

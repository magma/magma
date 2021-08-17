#! /bin/bash

set -ex
OVS_VER='2.15'

git clone --depth 1 https://github.com/magma/magma ovs-build
cd "ovs-build/third_party/gtp_ovs/ovs-gtp-patches/$OVS_VER/"
git clone --depth 1 --branch "branch-$OVS_VER"  https://github.com/openvswitch/ovs
cd ovs/
git am ../00*
DEB_BUILD_OPTIONS='parallel=8 nocheck' fakeroot debian/rules binary
cd ..
ls *.deb

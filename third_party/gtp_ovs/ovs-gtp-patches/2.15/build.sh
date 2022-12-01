#! /bin/bash

set -ex
OVS_VER='2.15'
DIR="ovs-build"
DEST=$1

sudo apt update
sudo apt install -y build-essential linux-headers-generic
sudo apt install -y dh-make debhelper dh-python devscripts python3-dev
sudo apt install -y graphviz libssl-dev python3-all python3-sphinx libunbound-dev libunwind-dev

rm -rf ~/$DIR
mkdir ~/$DIR
cd ~/$DIR

git clone  https://github.com/openvswitch/ovs
cd ovs/
git checkout 31288dc725be6bc8eaa4e8641ee28895c9d0fd7a
git apply "$MAGMA_ROOT/third_party/gtp_ovs/ovs-gtp-patches/$OVS_VER"/00*
DEB_BUILD_OPTIONS='parallel=8 nocheck' fakeroot debian/rules binary
cd ..
ls ./*.deb
if [[ -n "$DEST" ]] && [[ -d "$DEST" ]];
then
        mv ./*.deb "$DEST"
fi

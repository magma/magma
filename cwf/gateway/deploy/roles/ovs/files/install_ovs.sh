#!/bin/bash

apt-get -y update
apt-get -y install git
apt-get remove -y openvswitch-common
apt-get remove -y openvswitch-switch
apt-get -y install automake
apt-get -y install gcc
apt-get -y install libtool
apt-get -y libcap-ng-dev
apt-get -y install linux-headers-"$(uname -r)"
apt-get -y update
git clone https://github.com/openvswitch/ovs.git
cd ovs/ || exit
git checkout v2.12.0
git apply /tmp/0001-Add-custom-IPDR-fields-for-IPFIX-export.patch
git apply /tmp/0002-ovs-Handle-spaces-in-ovs-arguments.patch
git apply /tmp/0003-Add-pdp_start_epoch-custom-field-to-IPFIX-export.patch
./boot.sh
./configure --prefix=/usr --localstatedir=/var --sysconfdir=/etc --with-linux=/lib/modules/"$(uname -r)"/build
make
make install

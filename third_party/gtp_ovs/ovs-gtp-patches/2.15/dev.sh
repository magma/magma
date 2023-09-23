#! /bin/bash

set -x
OVS_VER='2.15'
DIR="/root/ovs-build"

MAGMA_ROOT="/home/vagrant/magma"

function setup() {

  set -e
  apt install -y build-essential linux-headers-generic
  apt install -y dh-make debhelper dh-python devscripts python3-dev
  apt install -y graphviz libssl-dev python3-all python3-sphinx libunbound-dev libunwind-dev
  apt install -y linux-modules-extra-`uname -r`
  pip3 install scapy

  rm -rf $DIR
  mkdir $DIR
  cd $DIR

  git clone  https://github.com/openvswitch/ovs
  cd ovs/
  git checkout 31288dc725be6bc8eaa4e8641ee28895c9d0fd7a
  git am "$MAGMA_ROOT/third_party/gtp_ovs/ovs-gtp-patches/$OVS_VER"/00*

  cd ../
  git clone https://gitea.osmocom.org/cellular-infrastructure/libgtpnl
  cd libgtpnl
  autoreconf -fi && ./configure && make && make install
  cp ./src/.libs/libgtpnl.so.0.1.2 /lib/libgtpnl.so.0
}


function  build_test() {
  service magma@* stop
  # sometimes this package is auto removed.

  apt install -y linux-modules-extra-`uname -r`

  sleep 1
  ifdown gtp_br0
  ifdown uplink_br0
  sleep 1
  ovs-dpctl del-dp ovs-system
  
  set -e
  cd $DIR/ovs

  ./boot.sh
  ./configure  --with-linux=/lib/modules/`uname -r`/build
  make clean
  make -j 10
  make -j 10 install

  sudo sysctl -w kernel.core_pattern=core.%u.%p.%t

  # insert required modules.
  modprobe udp_tunnel
  modprobe udp_tunnel
  modprobe nf_nat
  modprobe nf_defrag_ipv6
  modprobe nf_defrag_ipv4
  modprobe nf_conntrack
  modprobe tunnel6
  modprobe nf_defrag_ipv6
  modprobe nft_connlimit
  modprobe nft_counter
  modprobe gtp
  rmmod vport_gtp || true
  rmmod openvswitch || true

  cp datapath/linux/*.ko /lib/modules/`uname -r`/kernel/net/openvswitch/
  cp datapath/linux/*.ko /lib/modules/`uname -r`/updates/dkms/
  sync
  depmod -a

  make check-kmod TESTSUITEFLAGS="144-157"
  RET=$?
  exit $RET
}

$1


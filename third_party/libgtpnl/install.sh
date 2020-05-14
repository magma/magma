#!/usr/bin/env bash

function check_success {
  ret=$?
  if [[ $ret == 0 ]]; then
    return 0
  fi
  echo  "$1 failed with return code $ret"
  popd
  exit 1
}

set -e
HOME_DIR=/home/vagrant/
sudo apt-get install -y  autoconf automake build-essential libmnl-dev

rm -rf $HOME_DIR/libgtpnl
pushd $HOME_DIR

git clone https://git.osmocom.org/libgtpnl
check_success "Cloning libgtpnl git repo"

cd libgtpnl
autoreconf -fi
./configure
make -j`nproc`
check_success "Make"

sudo make install
check_success "Installing libgtpnl"

sudo ldconfig
check_success "Linking libgtpnl"

popd

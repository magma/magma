#!/usr/bin/env bash

# Exit on error
set -e
HOME_DIR=/home/vagrant/
sudo apt-get install -y  autoconf automake build-essential libmnl-dev

rm -rf $HOME_DIR/libgtpnl
pushd $HOME_DIR

git clone https://git.osmocom.org/libgtpnl
cd libgtpnl
autoreconf -fi
./configure
make -j`nproc`
sudo make install
sudo ldconfig
popd

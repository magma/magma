#!/bin/bash
################################################################################
# Copyright (c) Facebook, Inc. and its affiliates.
# All rights reserved.
#
# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.
################################################################################

set -e

# Packer may ssh into the box too early since SSH is ready before Debian
# actually is
sleep 30

apt-get update

# Install a deprecated linux kernel for openvswitch
if [ "$(uname -r)" != "4.9.0-9-amd64" ]; then
  # Adding the snapshot to retrieve 4.9.0-9-amd64
  if ! grep -q "deb http://snapshot.debian.org/archive/debian/20190801T025637Z" /etc/apt/sources.list; then
    echo "deb http://snapshot.debian.org/archive/debian/20190801T025637Z stretch main non-free contrib" >> /etc/apt/sources.list
  fi
  apt update
  # Installing prerequesites, Kvers, headers
  apt install -y sudo python-minimal aptitude linux-image-4.9.0-9-amd64 linux-headers-4.9.0-9-amd64
  # Removing dev repository snapshot from source.list
  sed -i '/20190801T025637Z/d' /etc/apt/sources.list
  # Removing incompatible Kernel version
  DEBIAN_FRONTEND=noninteractive apt remove -y linux-image-4.9.0-11-amd64
fi

# Install some packages
apt-get install -y openssh-server gcc rsync dirmngr

# Add the Etagecom key
apt-key adv --fetch-keys http://packages.magma.etagecom.io/pubkey.gpg

# Add the Etagecom magma repo
bash -c 'echo -e "deb http://packages.magma.etagecom.io magma-custom main" > /etc/apt/sources.list.d/packages_magma_etagecom_io.list'

# Create the preferences file for backports
bash -c 'cat <<EOF > /etc/apt/preferences.d/magma-preferences
Package: *
Pin: origin packages.magma.etagecom.io
Pin-Priority: 900
EOF'

apt-get update

# Disable daily auto updates, so that vagrant ansible scripts can
# acquire apt lock immediately on startup
systemctl stop apt-daily.timer
systemctl disable apt-daily.timer
systemctl disable apt-daily.service
systemctl daemon-reload

echo "Done"

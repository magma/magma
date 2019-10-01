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

apt-get install -y linux-headers-4.9.0-8-amd64 linux-image-4.9.0-8-amd64

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

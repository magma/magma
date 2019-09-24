#!/bin/bash
################################################################################
# Copyright (c) Facebook, Inc. and its affiliates.
# All rights reserved.
#
# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.
################################################################################

# Install the dependencies
apt-get install -y build-essential linux-headers-"$(uname -r)"

# Mount the guest additions iso and run the install script
mkdir -p /mnt/iso
mount -t iso9660 -o loop /home/vagrant/VBoxGuestAdditions.iso /mnt/iso

/mnt/iso/VBoxLinuxAdditions.run

umount /mnt/iso
rm -rf /mnt/iso /home/vagrant/VBoxGuestAdditions.iso

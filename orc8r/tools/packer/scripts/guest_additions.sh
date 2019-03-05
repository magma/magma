#!/bin/bash

# Install the dependencies
apt-get install -y build-essential linux-headers-"$(uname -r)"

# Mount the guest additions iso and run the install script
mkdir -p /mnt/iso
mount -t iso9660 -o loop /home/vagrant/VBoxGuestAdditions.iso /mnt/iso

/mnt/iso/VBoxLinuxAdditions.run

umount /mnt/iso
rm -rf /mnt/iso /home/vagrant/VBoxGuestAdditions.iso

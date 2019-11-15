#!/bin/sh

# changing intefaces name
sed -i 's/GRUB_CMDLINE_LINUX=""/GRUB_CMDLINE_LINUX="net.ifnames=0 biosdevname=0"/g' /etc/default/grub

# changing interface name
grub-mkconfig -o /boot/grub/grub.cfg
sed -i 's/enp1s0/eth0/g' /etc/network/interfaces

# configuring eth1
echo "auto eth1
iface eth1 inet static
address 10.0.2.1
netmask 255.255.255.0" > /etc/network/interfaces.d/eth1

# As 4.9.0-9-amd64 has been removed from the current deb repo we're temporary using a snapshot
if ! grep -q "deb http://snapshot.debian.org/archive/debian/20190801T025637Z" /etc/apt/sources.list; then
echo "deb http://snapshot.debian.org/archive/debian/20190801T025637Z stretch main non-free contrib" >> /etc/apt/sources.list
fi

# Update apt
apt update

# Installing prerequesites
apt install -y sudo python-minimal aptitude linux-image-4.9.0-9-amd64 linux-headers-4.9.0-9-amd64

# Removing dev repository snapshot from source.list
sed -i '/20190801T025637Z/d' /etc/apt/sources.list

# Making magma a sudoer
adduser magma sudo
if ! grep -q "magma ALL=(ALL) NOPASSWD:ALL" /etc/sudoers; then
    echo "magma ALL=(ALL) NOPASSWD:ALL" >> /etc/sudoers
fi


# Making sure .ssh is created in magma user
mkdir -p /home/magma/.ssh
chown magma:magma /home/magma/.ssh
# Removing incompatible Kernel version
DEBIAN_FRONTEND=noninteractive apt remove -y linux-image-4.9.0-11-amd64

reboot

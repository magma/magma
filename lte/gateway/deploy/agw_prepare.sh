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

# Installing prerequesites
apt install -y sudo python-minimal aptitude

# Making magma a sudoer
adduser magma sudo
if ! grep -q "magma ALL=(ALL) NOPASSWD:ALL" /etc/sudoers; then
    echo "magma ALL=(ALL) NOPASSWD:ALL" >> /etc/sudoers
fi


# Making sure .ssh is created in magma user
mkdir -p /home/magma/.ssh
chown magma:magma /home/magma/.ssh
# Removing incompatible Kernel version
apt remove -y linux-image-4.9.0-11-amd64

reboot

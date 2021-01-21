#!/bin/bash
# Check for system changes before magma deploy
# Setting up env variable, user and project path
KVERS=$(uname -r)
MAGMA_USER="magma"

echo "- Check if Debian is installed"
if ! grep -q 'Debian' /etc/issue; then
  echo "Debian is not installed"
else
  echo "Debian is installed"
fi

echo "- Check for correct Linux Headers"
if [ "$KVERS" != "4.9.0-9-amd64" ]; then
    echo "New Linux Headers will be Installed"
fi

echo "- Check for magma user"
if ! (getent passwd | grep -q 'magma'); then
    echo "magma user is not Installed"
elif  ! grep -q "$MAGMA_USER ALL=(ALL) NOPASSWD:ALL" /etc/sudoers; then
    echo "magma will be added to sudoers"
else
    echo "magma user configured"
fi

echo "- Check if both interfaces are named eth0 and eth1"
INTERFACES=$(ip -br a)
if [[ ! $INTERFACES == *'eth0'*  ]] || [[ ! $INTERFACES == *'eth1'* ]] || ! grep -q 'GRUB_CMDLINE_LINUX="net.ifnames=0 biosdevname=0"' /etc/default/grub; then
  echo "Interfaces will be renamed to eth0 and eth1"
  echo "eth0 will be set to dhcp and eth1 10.0.2.1"
else
  echo "eth0 will be set to dhcp and eth1 10.0.2.1"
fi

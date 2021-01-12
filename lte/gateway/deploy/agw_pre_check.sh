#!/bin/bash
# Check for system changes before magma deploy
# Setting up env variable, user and project path
RED='\033[0;31m'
WHITE='\033[1;37m'
GREEN='\033[0;32m'
NC='\033[0m' # No Color
KVERS=$(uname -r)

echo -e "${WHITE}Checking if Debian is installed"
if ! grep -q 'Debian' /etc/issue; then
  echo -e "${RED}Debian is not installed"
else
  echo -e "${GREEN}Debian is installed"
fi

echo -e "${WHITE}Check for correct Linux Headers"
if [ "$KVERS" != "4.9.0-9-amd64" ]; then
    echo -e "${RED}New Linux Headers will be Installed"
fi

echo -e "${WHITE}Check for magma User"
if ! (getent passwd | grep -q 'magma'); then
    echo -e "${RED}magma User is not Installed"
elif  ! grep -q "$MAGMA_USER ALL=(ALL) NOPASSWD:ALL" /etc/sudoers; then
    echo -e "${RED}magma will be added to sudoers"
fi

echo -e "${WHITE}Need to check if both interfaces are named eth0 and eth1"
INTERFACES=$(ip -br a)
if [[ ! $INTERFACES == *'eth0'*  ]] || [[ ! $INTERFACES == *'eth1'* ]] || ! grep -q 'GRUB_CMDLINE_LINUX="net.ifnames=$
  echo -e "${RED}Interfaces will be renamed to eth0 and eth1"
  echo -e "${RED}eth0 will be set to dhcp and eth1 10.0.2.1"
else
  echo -e "${RED}eth0 will be set to dhcp and eth1 10.0.2.1"
fi

echo -e "$NC"

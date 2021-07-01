#!/bin/bash

WHOAMI=$(whoami)
MAGMA_VERSION="1.6.0"
# Default is focal
OS_VERSION="focal"

echo "Checking if the script has been executed by root user"
if [ "$WHOAMI" != "root" ]; then
  echo "You're executing the script as $WHOAMI instead of root.. exiting"
  exit 1
fi

while true; do
    read -p "You're about to upgrade magma to $MAGMA_VERSION, are you sure?(y/n)" yn
    case $yn in
        [Yy]* ) break;;
        [Nn]* ) exit;;
        * ) echo "Please answer yes or no.";;
    esac
done

if grep -q 'Debian' /etc/issue; then
  OS_VERSION="stretch"
  # We don't support 1.6.0 for debian so bumping down to 1.5.2
  if [ "$MAGMA_VERSION" == "1.6.0" ]; then
    MAGMA_VERSION="1.5.2"
  fi
fi

apt update
apt install -y apt-transport-https gnupg2

# We have changed the name too many time we have to wipe all versions
rm -rf /etc/apt/sources.list.d/*

wget https://artifactory.magmacore.org:443/artifactory/api/gpg/key/public -O /tmp/public
apt-key add /tmp/public

echo "deb https://artifactory.magmacore.org/artifactory/debian $OS_VERSION-$MAGMA_VERSION main" >> /etc/apt/sources.list.d/magma.list

apt update
apt install -y magma -o Dpkg::Options::="--force-overwrite"

#Upgrade OVS
ovs-kmod-upgrade.sh -y

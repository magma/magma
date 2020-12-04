#!/bin/bash
if [ -z "$1" ]; then
  echo "Please specify a version"
  exit
fi

VERSION="$1"
echo -e "deb http://packages.magma.etagecom.io stretch-${VERSION} main" > /etc/apt/sources.list.d/packages_magma_etagecom_io.list
apt update

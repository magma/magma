#!/bin/bash

set -euo pipefail

BOX_FILE=$1

USER=magmacore
BOX=$(basename $BOX_FILE | cut -d_ -f1-2 | cut -d. -f1)
VERSION="1.2.$(date +"%Y%m%d")"
BOX_PROVIDER=virtualbox
if echo $BOX_FILE | grep -q libvirt; then
  BOX_PROVIDER=libvirt
fi

#source vagrant-cloud-token
if [ -z "$VAGRANT_CLOUD_TOKEN" ]; then
  echo "VAGRANT_CLOUD_TOKEN variable is unset. Cannot continue." 1>&2
  exit 1
fi

vagrant cloud auth login --token "$VAGRANT_CLOUD_TOKEN"
vagrant cloud publish \
    "${USER}/${BOX}"
    "$VERSION" \
    "$BOX_PROVIDER" \
    "$BOX_FILE" \
    --release

#!/bin/bash

set -euo pipefail

BOX_FILE=$1
FORCE="--no-force"

USER=magmacore
BOX=$(basename $BOX_FILE | cut -d_ -f1-2 | cut -d. -f1)
VERSION="1.3.$(date +"%Y%m%d")"
BOX_PROVIDER=virtualbox

#source vagrant-cloud-token
if [ -z "$VAGRANT_CLOUD_TOKEN" ]; then
  echo "VAGRANT_CLOUD_TOKEN variable is unset. Cannot continue." 1>&2
  exit 1
fi

shift
while getopts ":f" opt; do
  case $opt in
    f)
      FORCE="--force"
      ;;
    \?)
      echo "Invalid option: -$OPTARG" >&2
      ;;
  esac
done

vagrant cloud auth login --token "$VAGRANT_CLOUD_TOKEN"
vagrant cloud publish \
    --release \
    "${FORCE}" \
    "${USER}/${BOX}" \
    "$VERSION" \
    "$BOX_PROVIDER" \
    "$BOX_FILE"

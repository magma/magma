#!/bin/bash
# Syntax: ./make_magma_pvcs.sh NAMESPACE STORAGE_CLASS_NAME ACCESS_MODE NAME SIZE
# Example: ./make_magma_pvcs.sh magma nfs ReadWriteMany openvpn 2Mi

set -e
set +x
set +v
set -o pipefail

function usage() {
  echo
  echo "Syntax: $0 NAMESPACE STORAGE_CLASS_NAME ACCESS_MODE NAME SIZE"
  echo "Example: $0 magma nfs ReadWriteMany openvpn 2Mi"
}

export namespace=$1
export storageclass=$2
export accessmode=$3
export pvc=$4
export size=$5

if [ -z "$1" ]; then
  echo "Error: Namespace required." 1>&2
  usage
  exit 1
fi

if [ -z "$2" ]; then
  echo "Error: storageClassName required." 1>&2
  usage
  exit 1
fi

if [ -z "$3" ]; then
  echo "Error: access mode required. " 1>&2
  usage
  exit 1
fi

if [ -z "$4" ]; then
  echo "Error: PVC name required. " 1>&2
  usage
  exit 1
fi

if [ -z "$5" ]; then
  echo "Error: size required. " 1>&2
  usage
  exit 1
fi

cat << EOF | kubectl -n $namespace apply -f -
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: ${pvc}
spec:
  accessModes:
  - ${accessmode}
  resources:
    requests:
      storage: ${size}
  storageClassName: ${storageclass}
  volumeMode: Filesystem
EOF

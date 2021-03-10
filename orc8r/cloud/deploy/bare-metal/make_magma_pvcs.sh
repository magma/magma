#!/usr/bin/env bash

# Syntax: ./make_magma_pvcs.sh NAMESPACE STORAGE_CLASS_NAME ACCESS_MODE
# Example: ./make_magma_pvcs.sh magma nfs ReadWriteMany

set -e
set +x
set +v
set -o pipefail

function usage() {
  echo 
  echo "Syntax: $0 NAMESPACE STORAGE_CLASS_NAME ACCESS_MODE"
  echo "Example: $0 magma nfs ReadWriteMany"
}

#PVC names and respective sizes
declare -A pvcs
pvcs[grafanadashboards]=2Gi
pvcs[grafanadata]=2Gi
pvcs[grafanadatasources]=100M
pvcs[grafanaproviders]=100M
pvcs[openvpn]=2M
pvcs[promcfg]=1Gi
pvcs[promdata]=64Gi

export namespace=$1
export storageclass=$2
export accessmode=$3

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
  export accessmode="ReadWriteMany"
fi

for pvc in "${!pvcs[@]}"; do
  export pvcname=$pvc
  export pvsize="${pvcs[$pvc]}"
  echo "Creating pvc $pvcname size $pvsize in namespace $namespace..."
  envsubst < pvc.yaml.template | kubectl -n $namespace apply -f -
done


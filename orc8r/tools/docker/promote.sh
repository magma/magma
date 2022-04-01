#!/usr/bin/env bash

set -ex

MAGMA_TAG=1.7
NEW_MAGMA_TAG=1.7.0
MAGMA_ARTIFACTORY=artifactory.magmacore.org

declare -A repositories=(
  [feg]="gateway_go gateway_python"
  [orc8r]="controller magmalte nginx"
  [agw]="agw_gateway_c agw_gateway_python ghz_gateway_c ghz_gateway_python"
  [cwf]="cwag_go gateway_go gateway_pipelined gateway_python gateway_sessiond" # operator
)

for repo in ${!repositories[@]}; do
  for image in ${repositories[${repo}]}; do
    # Pull docker image from test registry
    docker pull ${repo}-test.${MAGMA_ARTIFACTORY}/${image}:${MAGMA_TAG}

    # Tag docker image with new tag
    docker tag ${repo}-test.${MAGMA_ARTIFACTORY}/${image}:${MAGMA_TAG} ${repo}-prod.${MAGMA_ARTIFACTORY}/${image}:${NEW_MAGMA_TAG} 

    # Push docker image to prod registry
    docker push ${repo}-prod.${MAGMA_ARTIFACTORY}/${image}:${NEW_MAGMA_TAG} 
  done
done

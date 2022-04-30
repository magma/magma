#!/usr/bin/env bash

set -ex

MAGMA_TAG=1.7
NEW_MAGMA_TAG=1.7.1
MAGMA_ARTIFACTORY=artifactory.magmacore.org

declare -A repositories=(
  [orc8r]="controller magmalte nginx"
  [feg]="gateway_go gateway_python"
  [agw]="agw_gateway_c agw_gateway_python ghz_gateway_c ghz_gateway_python"
  [cwf]="cwag_go gateway_go gateway_pipelined gateway_python gateway_sessiond operator"
)

for repo in ${!repositories[@]}; do
  for image in ${repositories[${repo}]}; do

    # Change docker URL to Artifactory
    sed -i "s/docker/${repo}-prod/g" ~/.docker/config.json

    # Pull docker image from test registry
    docker pull ${repo}-test.${MAGMA_ARTIFACTORY}/${image}:${MAGMA_TAG}

    # Tag docker image with new tag
    docker tag ${repo}-test.${MAGMA_ARTIFACTORY}/${image}:${MAGMA_TAG} ${repo}-prod.${MAGMA_ARTIFACTORY}/${image}:${NEW_MAGMA_TAG}
    docker tag ${repo}-test.${MAGMA_ARTIFACTORY}/${image}:${MAGMA_TAG} ${repo}-prod.${MAGMA_ARTIFACTORY}/${image}:latest

    # Push docker image to prod registry
    docker push ${repo}-prod.${MAGMA_ARTIFACTORY}/${image}:${NEW_MAGMA_TAG}
    docker push ${repo}-prod.${MAGMA_ARTIFACTORY}/${image}:latest

    # Remove uploaded image
    docker rmi ${repo}-test.${MAGMA_ARTIFACTORY}/${image}:${MAGMA_TAG}
    docker rmi ${repo}-prod.${MAGMA_ARTIFACTORY}/${image}:${NEW_MAGMA_TAG}
    docker rmi ${repo}-prod.${MAGMA_ARTIFACTORY}/${image}:latest

    # Change docker URL back to docker
    sed -i "s/${repo}-prod/docker/g" ~/.docker/config.json

  done
done

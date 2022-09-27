#!/usr/bin/env bash

set -eou pipefail

if [[ -z $MAGMA_ARTIFACTORY ]]; then
  exitmsg "Environment variable MAGMA_ARTIFACTORY must be set."
fi

if [[ -z $BRANCH_TAG ]]; then
  exitmsg "Environment variable BRANCH_TAG must be set."
fi

if [[ -z $RELEASE_TAG ]]; then
  exitmsg "Environment variable RELEASE_TAG must be set."
fi

declare -A repositories=(
  [orc8r]="controller magmalte nginx active-mode-controller configuration-controller radio-controller db-service"
  [feg]="gateway_go gateway_python"
  [agw]="agw_gateway_c agw_gateway_python ghz_gateway_c ghz_gateway_python agw_gateway_c_arm agw_gateway_python_arm"
  [cwf]="cwag_go gateway_go gateway_pipelined gateway_python gateway_sessiond operator"
)

# shellcheck disable=SC2068
for repo in ${!repositories[@]}; do
  for image in ${repositories[${repo}]}; do

    # Change docker URL to Artifactory
    sed -i "s/docker/${repo}-prod/g" ~/.docker/config.json

    # Pull docker image from test registry
    docker pull "${repo}-test.${MAGMA_ARTIFACTORY}/${image}:${BRANCH_TAG}"

    # Tag docker image with new tag
    docker tag "${repo}-test.${MAGMA_ARTIFACTORY}/${image}:${BRANCH_TAG} ${repo}-prod.${MAGMA_ARTIFACTORY}/${image}:${RELEASE_TAG}"
    docker tag "${repo}-test.${MAGMA_ARTIFACTORY}/${image}:${BRANCH_TAG} ${repo}-prod.${MAGMA_ARTIFACTORY}/${image}:latest"

    # Push docker image to prod registry
    docker push "${repo}-prod.${MAGMA_ARTIFACTORY}/${image}:${RELEASE_TAG}"
    docker push "${repo}-prod.${MAGMA_ARTIFACTORY}/${image}:latest"

    # Remove uploaded image
    docker rmi "${repo}-test.${MAGMA_ARTIFACTORY}/${image}:${BRANCH_TAG}"
    docker rmi "${repo}-prod.${MAGMA_ARTIFACTORY}/${image}:${RELEASE_TAG}"
    docker rmi "${repo}-prod.${MAGMA_ARTIFACTORY}/${image}:latest"

    # Change docker URL back to docker
    sed -i "s/${repo}-prod/docker/g" ~/.docker/config.json
    echo "Promoted docker image artifact ${image} from ${repo}-test to ${repo}-prod registry successfully."
  done
done

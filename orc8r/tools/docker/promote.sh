#!/usr/bin/env bash

# Copyright 2022 The Magma Authors.
#
# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

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
  [cwag]="cwag_go gateway_go gateway_pipelined gateway_python gateway_sessiond operator"
)

# shellcheck disable=SC2068
for repo in ${!repositories[@]}; do
  for image in ${repositories[${repo}]}; do

    # Change docker URL to Artifactory
    sed -i "s/docker/magma-docker-${repo}-prod/g" ~/.docker/config.json

    # Pull docker image from test registry
    docker pull "${MAGMA_ARTIFACTORY}/magma-docker-${repo}-test/${image}:${BRANCH_TAG}"

    # Tag docker image with new tag
    docker tag "${MAGMA_ARTIFACTORY}/magma-docker-${repo}-test/${image}:${BRANCH_TAG}" "${MAGMA_ARTIFACTORY}/magma-docker-${repo}-prod/${image}:${RELEASE_TAG}"
    docker tag "${MAGMA_ARTIFACTORY}/magma-docker-${repo}-test/${image}:${BRANCH_TAG}" "${MAGMA_ARTIFACTORY}/magma-docker-${repo}-prod/${image}:latest"

    # Push docker image to prod registry
    docker push "${MAGMA_ARTIFACTORY}/magma-docker-${repo}-prod/${image}:${RELEASE_TAG}"
    docker push "${MAGMA_ARTIFACTORY}/magma-docker-${repo}-prod/${image}:latest"

    # Remove uploaded image
    docker rmi "${MAGMA_ARTIFACTORY}/magma-docker-${repo}-test/${image}:${BRANCH_TAG}"
    docker rmi "${MAGMA_ARTIFACTORY}/magma-docker-${repo}-prod/${image}:${RELEASE_TAG}"
    docker rmi "${MAGMA_ARTIFACTORY}/magma-docker-${repo}-prod/${image}:latest"

    # Change docker URL back to docker
    sed -i "s/magma-docker-${repo}-prod/docker/g" ~/.docker/config.json
    echo "Promoted docker image artifact ${image} from magma-docker-${repo}-test to magma-docker-${repo}-prod registry successfully."
  done
done

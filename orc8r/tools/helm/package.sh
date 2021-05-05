#!/usr/bin/env bash

# Copyright 2020 The Magma Authors.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# package.sh packages and publishes orc8r helm charts to a private git repo

set -e -o pipefail

# Valid deployment types
FWA="fwa"
FFWA="federated_fwa"
ALL="all"
ORC8R_VERSION="1.4"

# package chart, update index.yaml and push it to artifactory
update_and_send_to_artifactory () {
  CHART_PATH=$1
  helm dependency update $CHART_PATH
  ARTIFACT_PATH=$(helm package $CHART_PATH | awk '{print $8}')
  helm repo index .
  MD5_CHECKSUM=$(md5sum $ARTIFACT_PATH | awk '{print $1}')
  SHA1_CHECKSUM=$(shasum -a 1 $ARTIFACT_PATH | awk '{ print $1 }')
  SHA256_CHECKSUM=$(shasum -a 256 $ARTIFACT_PATH | awk '{ print $1 }')
  curl -u $HELM_CHART_MUSEUM_USERNAME:$HELM_CHART_MUSEUM_TOKEN \
              --header "X-Checksum-MD5:${MD5_CHECKSUM}" \
              --header "X-Checksum-Sha1:${SHA1_CHECKSUM}" \
              --header "X-Checksum-Sha256:${SHA256_CHECKSUM}" \
               -T $ARTIFACT_PATH $HELM_CHART_MUSEUM_URL/$(basename $ARTIFACT_PATH)

} 


usage() {
  echo "Usage: $0 -d DEPLOYMENT_TYPE"
  exit 2
}

exitmsg() {
  echo "$1"
  exit 1
}

# Parse the args and declare defaults
while getopts 'd:h' OPT; do
  case "${OPT}" in
    d) DEPLOYMENT_TYPE=${OPTARG} ;;
    h|*) usage ;;
  esac
done

# Check if the required args and env-vars present
[[ -z "${DEPLOYMENT_TYPE}" ]] && usage

if [ "$DEPLOYMENT_TYPE" != "$FWA" ] && [ "$DEPLOYMENT_TYPE" != "$FFWA" ] && [ "$DEPLOYMENT_TYPE" != "$ALL" ]; then
  exitmsg "Deployment type '$DEPLOYMENT_TYPE' is not valid. Valid types are: \
  ['$FWA', '$FFWA', '$ALL']"
fi

# Check for artifactory URL presence
if [[ -z $HELM_CHART_MUSEUM_URL ]]; then

  if [[ -z $GITHUB_REPO ]]; then
    exitmsg "Environment variable GITHUB_REPO must be set"
  fi

  if [[ -z $GITHUB_REPO_URL ]]; then
    exitmsg "Environment variable GITHUB_REPO_URL must be set"
  fi

  if [[ -z $GITHUB_USERNAME ]]; then
    exitmsg "Environment variable GITHUB_USERNAME must be set"
  fi

  if [[ -z $GITHUB_ACCESS_TOKEN ]]; then
    exitmsg "Environment variable GITHUB_ACCESS_TOKEN must be set"
  fi

  if [[ -z $MAGMA_ROOT ]]; then
    exitmsg "Environment variable MAGMA_ROOT must be set"
  fi

  # Set up repo for charts
  mkdir -p ~/magma-charts && cd ~/magma-charts
  git init

  # Begin packaging necessary helm charts
  helm dependency update $MAGMA_ROOT/orc8r/cloud/helm/orc8r/
  helm package $MAGMA_ROOT/orc8r/cloud/helm/orc8r/ && helm repo index .

  if [ "$DEPLOYMENT_TYPE" == "$FWA" ]; then
    helm dependency update $MAGMA_ROOT/lte/cloud/helm/lte-orc8r/
    helm package $MAGMA_ROOT/lte/cloud/helm/lte-orc8r/ && helm repo index .
  fi

  if [ "$DEPLOYMENT_TYPE" == "$FFWA" ]; then
    helm dependency update $MAGMA_ROOT/lte/cloud/helm/lte-orc8r/
    helm package $MAGMA_ROOT/lte/cloud/helm/lte-orc8r/ && helm repo index .

    helm dependency update $MAGMA_ROOT/feg/cloud/helm/feg-orc8r/
    helm package $MAGMA_ROOT/feg/cloud/helm/feg-orc8r/ && helm repo index .
  fi

  if  [ "$DEPLOYMENT_TYPE" == "$ALL" ]; then
    helm dependency update $MAGMA_ROOT/cwf/cloud/helm/cwf-orc8r/
    helm package $MAGMA_ROOT/cwf/cloud/helm/cwf-orc8r/ && helm repo index .

    helm dependency update $MAGMA_ROOT/lte/cloud/helm/lte-orc8r/
    helm package $MAGMA_ROOT/lte/cloud/helm/lte-orc8r/ && helm repo index .

    helm dependency update $MAGMA_ROOT/feg/cloud/helm/feg-orc8r/
    helm package $MAGMA_ROOT/feg/cloud/helm/feg-orc8r/ && helm repo index .

    helm dependency update $MAGMA_ROOT/fbinternal/cloud/helm/fbinternal-orc8r/
    helm package $MAGMA_ROOT/fbinternal/cloud/helm/fbinternal-orc8r/ && helm repo index .

    helm dependency update $MAGMA_ROOT/wifi/cloud/helm/wifi-orc8r/
    helm package $MAGMA_ROOT/wifi/cloud/helm/wifi-orc8r/ && helm repo index .
  fi

  # Push charts
  git add . && git commit -m "orc8r charts commit for version $ORC8R_VERSION"
  git config remote.origin.url >&- || git remote add origin $GITHUB_REPO_URL
  git push -u origin master

  # Ensure push was successful
  helm repo add $GITHUB_REPO --username $GITHUB_USERNAME --password $GITHUB_ACCESS_TOKEN \
        "https://raw.githubusercontent.com/$GITHUB_USERNAME/$GITHUB_REPO/master/"
  helm repo update

  # The helm command returns 0 even when no results are found. Search for err str
  # instead
  HELM_SEARCH_RESULTS=$(helm search repo $GITHUB_REPO) # should list the uploaded charts
  if [ "$HELM_SEARCH_RESULTS" == "No results found" ]; then
    exitmsg "Error! Unable to find uploaded orc8r charts"
  fi
else

  if [[ -z $HELM_CHART_MUSEUM_USERNAME ]]; then
    exitmsg "Environment variable HELM_CHART_MUSEUM_USERNAME must be set"
  fi

  if [[ -z $HELM_CHART_MUSEUM_TOKEN ]]; then
    exitmsg "Environment variable HELM_CHART_MUSEUM_TOKEN must be set"
  fi

  if [[ -z $MAGMA_ROOT ]]; then
    exitmsg "Environment variable MAGMA_ROOT must be set"
  fi

  # Begin packaging necessary helm charts
  update_and_send_to_artifactory $MAGMA_ROOT/orc8r/cloud/helm/orc8r/

  if [ "$DEPLOYMENT_TYPE" == "$FWA" ]; then
    update_and_send_to_artifactory $MAGMA_ROOT/lte/cloud/helm/lte-orc8r/
  fi

  if [ "$DEPLOYMENT_TYPE" == "$FFWA" ]; then
    update_and_send_to_artifactory $MAGMA_ROOT/lte/cloud/helm/lte-orc8r/

    update_and_send_to_artifactory $MAGMA_ROOT/feg/cloud/helm/feg-orc8r/
  fi

  if  [ "$DEPLOYMENT_TYPE" == "$ALL" ]; then
    update_and_send_to_artifactory $MAGMA_ROOT/cwf/cloud/helm/cwf-orc8r/

    update_and_send_to_artifactory $MAGMA_ROOT/lte/cloud/helm/lte-orc8r/
    
    update_and_send_to_artifactory $MAGMA_ROOT/feg/cloud/helm/feg-orc8r/

    update_and_send_to_artifactory $MAGMA_ROOT/fbinternal/cloud/helm/fbinternal-orc8r/
    
    update_and_send_to_artifactory $MAGMA_ROOT/wifi/cloud/helm/wifi-orc8r/
  fi

  # Push index.yaml
  INDEX_MD5_CHECKSUM=$(md5sum $MAGMA_ROOT/index.yaml | awk '{print $1}')
  INDEX_SHA1_CHECKSUM=$(shasum -a 1 $MAGMA_ROOT/index.yaml | awk '{ print $1 }')
  INDEX_SHA256_CHECKSUM=$(shasum -a 256 $MAGMA_ROOT/index.yaml | awk '{ print $1 }')
  curl -u $HELM_CHART_MUSEUM_USERNAME:$HELM_CHART_MUSEUM_TOKEN \
              --header "X-Checksum-MD5:${INDEX_MD5_CHECKSUM}" \
              --header "X-Checksum-Sha1:${INDEX_SHA1_CHECKSUM}" \
              --header "X-Checksum-Sha256:${INDEX_SHA256_CHECKSUM}" \
              -T $MAGMA_ROOT/index.yaml $HELM_CHART_MUSEUM_URL/index.yaml

  # Ensure push was successful
  helm repo add $HELM_CHART_MUSEUM_URL $HELM_CHART_MUSEUM_URL --username $HELM_CHART_MUSEUM_USERNAME --password $HELM_CHART_MUSEUM_TOKEN 
  helm repo update

  # The helm command returns 0 even when no results are found. Search for err str
  # instead
  HELM_SEARCH_RESULTS=$(helm search repo $HELM_CHART_MUSEUM_URL) # should list the uploaded charts
  if [ "$HELM_SEARCH_RESULTS" == "No results found" ]; then
    exitmsg "Error! Unable to find uploaded orc8r charts"
  fi
fi


echo "Uploaded orc8r charts successfully."

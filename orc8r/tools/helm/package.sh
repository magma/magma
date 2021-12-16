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
  CHART_PATH="$1"
  helm dependency update "$CHART_PATH"
  # shellcheck disable=SC2086
  # We want $VERSION to be split as this is an option to the helm command
  ARTIFACT_PATH="$(helm package "$CHART_PATH" $VERSION | awk '{print $8}')"
  helm repo index .

  if [[ $ONLY_PACKAGE = true ]]; then
    mv "$ARTIFACT_PATH" "$MAGMA_ROOT/charts"
  else
    MD5_CHECKSUM="$(md5sum "$ARTIFACT_PATH" | awk '{print $1}')"
    SHA1_CHECKSUM="$(shasum -a 1 "$ARTIFACT_PATH" | awk '{ print $1 }')"
    SHA256_CHECKSUM="$(shasum -a 256 "$ARTIFACT_PATH" | awk '{ print $1 }')"
    curl --user "$HELM_CHART_MUSEUM_USERNAME":"$HELM_CHART_MUSEUM_TOKEN" --fail \
                --header "X-Checksum-MD5:${MD5_CHECKSUM}" \
                --header "X-Checksum-Sha1:${SHA1_CHECKSUM}" \
                --header "X-Checksum-Sha256:${SHA256_CHECKSUM}" \
                --upload-file "$ARTIFACT_PATH" "$HELM_CHART_MUSEUM_URL/$(basename "$ARTIFACT_PATH")"
  fi
}

usage() {
  echo "Usage: $0 [-v|--version V] [-d|--deployment-type $FWA|$FFWA|$ALL] [-p|--only-package]"
  exit 2
}

exitmsg() {
  echo "$1"
  exit 1
}

ONLY_PACKAGE=false
# Parse the args
while [[ $# -gt 0 ]]
do
key="$1"
case $key in
    -v|--version)
    VERSION="--version $2"
    shift  # pass argument or value
    ;;
    -d|--deployment-type)
    DEPLOYMENT_TYPE="$2"
    shift
    ;;
    -p|--only-package)
    ONLY_PACKAGE=true
    ;;
    -h|--help)
    usage
    shift
    ;;
    *)
    echo "Error: unknown cmdline option: $key"
    usage
    ;;
esac
shift  # past argument or value
done

# Check if the required args and env-vars present
[[ -z "${DEPLOYMENT_TYPE}" ]] && usage

if [ "$DEPLOYMENT_TYPE" != "$FWA" ] && [ "$DEPLOYMENT_TYPE" != "$FFWA" ] && [ "$DEPLOYMENT_TYPE" != "$ALL" ]; then
  exitmsg "Deployment type '$DEPLOYMENT_TYPE' is not valid. Valid types are: \
  ['$FWA', '$FFWA', '$ALL']"
fi

# Check for artifactory URL presence
if [[ $ONLY_PACKAGE = false && -z $HELM_CHART_ARTIFACTORY_URL ]]; then
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

  # Begin packaging necessary Helm charts
  helm dependency update "$MAGMA_ROOT/orc8r/cloud/helm/orc8r/"
  # shellcheck disable=SC2086
  # We want $VERSION to be splitted as this is an option to the helm command
  helm package "$MAGMA_ROOT/orc8r/cloud/helm/orc8r/" $VERSION && helm repo index .

  if [ "$DEPLOYMENT_TYPE" == "$FWA" ]; then
    helm dependency update "$MAGMA_ROOT/lte/cloud/helm/lte-orc8r/"
    # shellcheck disable=SC2086
    helm package "$MAGMA_ROOT/lte/cloud/helm/lte-orc8r/" $VERSION && helm repo index .
  fi

  if [ "$DEPLOYMENT_TYPE" == "$FFWA" ]; then
    helm dependency update "$MAGMA_ROOT/lte/cloud/helm/lte-orc8r/"
    # shellcheck disable=SC2086
    helm package "$MAGMA_ROOT/lte/cloud/helm/lte-orc8r/" $VERSION && helm repo index .

    helm dependency update "$MAGMA_ROOT/feg/cloud/helm/feg-orc8r/"
    # shellcheck disable=SC2086
    helm package "$MAGMA_ROOT/feg/cloud/helm/feg-orc8r/" $VERSION && helm repo index .
  fi

  if  [ "$DEPLOYMENT_TYPE" == "$ALL" ]; then
    helm dependency update "$MAGMA_ROOT/cwf/cloud/helm/cwf-orc8r/"
    # shellcheck disable=SC2086
    helm package "$MAGMA_ROOT/cwf/cloud/helm/cwf-orc8r/" $VERSION && helm repo index .

    helm dependency update "$MAGMA_ROOT/lte/cloud/helm/lte-orc8r/"
    # shellcheck disable=SC2086
    helm package "$MAGMA_ROOT/lte/cloud/helm/lte-orc8r/" $VERSION && helm repo index .

    helm dependency update "$MAGMA_ROOT/feg/cloud/helm/feg-orc8r/"
    # shellcheck disable=SC2086
    helm package "$MAGMA_ROOT/feg/cloud/helm/feg-orc8r/" $VERSION && helm repo index .

    helm dependency update "$MAGMA_ROOT/fbinternal/cloud/helm/fbinternal-orc8r/"
    # shellcheck disable=SC2086
    helm package "$MAGMA_ROOT/fbinternal/cloud/helm/fbinternal-orc8r/" $VERSION && helm repo index .
  fi

  # Push charts
  git add . && git commit -m "orc8r charts commit for version $ORC8R_VERSION"
  git config remote.origin.url >&- || git remote add origin "$GITHUB_REPO_URL"
  git push -u origin master

  # Ensure push was successful
  helm repo add "$GITHUB_REPO" --username "$GITHUB_USERNAME" --password "$GITHUB_ACCESS_TOKEN" \
        "https://raw.githubusercontent.com/$GITHUB_USERNAME/$GITHUB_REPO/master/"
  helm repo update

  # The Helm command returns 0 even when no results are found. Search for err str
  # instead
  HELM_SEARCH_RESULTS="$(helm search repo "$GITHUB_REPO")" # should list the uploaded charts
  if [ "$HELM_SEARCH_RESULTS" == "No results found" ]; then
    exitmsg "Error! Unable to find uploaded orc8r charts"
  fi
else
  if [[ -z $MAGMA_ROOT ]]; then
    exitmsg "Environment variable MAGMA_ROOT must be set"
  fi

  if [[ $ONLY_PACKAGE = false ]]; then
      if [[ -z $HELM_CHART_MUSEUM_REPO ]]; then
    exitmsg "Environment variable $HELM_CHART_MUSEUM_REPO must be set"
    fi

    if [[ -z $HELM_CHART_MUSEUM_USERNAME ]]; then
      exitmsg "Environment variable HELM_CHART_MUSEUM_USERNAME must be set"
    fi

    if [[ -z $HELM_CHART_MUSEUM_TOKEN ]]; then
      exitmsg "Environment variable HELM_CHART_MUSEUM_TOKEN must be set"
    fi
    # Trim last backslash if exists
    # shellcheck disable=SC2001
    HELM_CHART_ARTIFACTORY_URL="$(echo "$HELM_CHART_ARTIFACTORY_URL" | sed 's:/$::')"
    # Verify existence of the helm repo
    RESPONSE_CODE_REPO="$(curl --output /dev/null --stderr /dev/null --silent --write-out "%{http_code}"  "$HELM_CHART_ARTIFACTORY_URL/$HELM_CHART_MUSEUM_REPO/" || :)"
    if [ $ONLY_PACKAGE != "True" ] && [ "$RESPONSE_CODE_REPO" != "200" ]; then
      exitmsg "There was an error connecting to the artifactory repository $HELM_CHART_MUSEUM_ORIGIN_REPO, the http error code was $RESPONSE_CODE_REPO"
    fi

    HELM_CHART_MUSEUM_URL="$HELM_CHART_ARTIFACTORY_URL/$HELM_CHART_MUSEUM_REPO"
    # Form API URL
    HELM_CHART_MUSEUM_API_URL="$HELM_CHART_ARTIFACTORY_URL/api"
  fi

  # Begin packaging necessary Helm charts
  update_and_send_to_artifactory "$MAGMA_ROOT/orc8r/cloud/helm/orc8r/"

  if [ "$DEPLOYMENT_TYPE" == "$FWA" ]; then
    update_and_send_to_artifactory "$MAGMA_ROOT/lte/cloud/helm/lte-orc8r/"
  fi

  if [ "$DEPLOYMENT_TYPE" == "$FFWA" ]; then
    update_and_send_to_artifactory "$MAGMA_ROOT/lte/cloud/helm/lte-orc8r/"
    update_and_send_to_artifactory "$MAGMA_ROOT/feg/cloud/helm/feg-orc8r/"
  fi

  if  [ "$DEPLOYMENT_TYPE" == "$ALL" ]; then
    update_and_send_to_artifactory "$MAGMA_ROOT/cwf/cloud/helm/cwf-orc8r/"
    update_and_send_to_artifactory "$MAGMA_ROOT/lte/cloud/helm/lte-orc8r/"
    update_and_send_to_artifactory "$MAGMA_ROOT/feg/cloud/helm/feg-orc8r/"
    update_and_send_to_artifactory "$MAGMA_ROOT/fbinternal/cloud/helm/fbinternal-orc8r/"
  fi

  if [[ $ONLY_PACKAGE = false ]]; then
    # Refresh index.yaml
    curl --request POST --user "$HELM_CHART_MUSEUM_USERNAME":"$HELM_CHART_MUSEUM_TOKEN" \
                "$HELM_CHART_MUSEUM_API_URL/helm/$HELM_CHART_MUSEUM_REPO/reindex"

    # Ensure push was successful
    helm repo add "$(basename "$HELM_CHART_MUSEUM_URL")" "$HELM_CHART_MUSEUM_URL" --username "$HELM_CHART_MUSEUM_USERNAME" --password "$HELM_CHART_MUSEUM_TOKEN"
    helm repo update

    # The Helm command returns 0 even when no results are found. Search for err str
    # instead
    HELM_SEARCH_RESULTS="$(helm search repo "$(basename "$HELM_CHART_MUSEUM_URL")")" # should list the uploaded charts
    if [ "$HELM_SEARCH_RESULTS" == "No results found" ]; then
      exitmsg "Error! Unable to find uploaded orc8r charts"
    fi
  fi
fi


echo "Uploaded orc8r charts successfully."

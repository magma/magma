#!/usr/bin/env bash

# Copyright 2020 The Magma Authors.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# promote.sh copies specified artifacts from one repo to another.

set -e -o pipefail

promote_artifact () {
  ARTIFACT="$1"
  curl --request POST --user "$HELM_CHART_MUSEUM_USERNAME":"$HELM_CHART_MUSEUM_TOKEN" --fail \
   "$HELM_CHART_MUSEUM_API_URL/copy/$HELM_CHART_MUSEUM_ORIGIN_REPO/$ARTIFACT?to=/$HELM_CHART_MUSEUM_DEST_REPO/$ARTIFACT"
}

get_artifact () {
  ARTIFACT="$1"
  curl --output /dev/null --silent  --write-out "%{http_code}" \
   "$HELM_CHART_ARTIFACTORY_URL/$HELM_CHART_MUSEUM_ORIGIN_REPO/$ARTIFACT" || :
}

usage() {
  echo "Supply at least one artifact to promote: $0 ARTIFACT_PATH"
  exit 2
}

exitmsg() {
  echo "$1"
  exit 1
}

# Check if the required args and env-vars present
[ $# -eq 0 ] && usage

# Define default values
HELM_CHART_ARTIFACTORY_URL="${HELM_CHART_ARTIFACTORY_URL:-https://artifactory.magmacore.org:443/artifactory}"
HELM_CHART_MUSEUM_ORIGIN_REPO="${HELM_CHART_MUSEUM_ORIGIN_REPO:-helm-test}"
HELM_CHART_MUSEUM_DEST_REPO="${HELM_CHART_MUSEUM_DEST_REPO:-helm-prod}"

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
RESPONSE_CODE_REPO="$(curl --output /dev/null --stderr /dev/null --silent --write-out "%{http_code}"  "$HELM_CHART_ARTIFACTORY_URL/$HELM_CHART_MUSEUM_ORIGIN_REPO/" || :)"
if [ "$RESPONSE_CODE_REPO" != "200" ]; then
  exitmsg "There was an error connecting to the artifactory repository $HELM_CHART_MUSEUM_ORIGIN_REPO, the http error code was $RESPONSE_CODE_REPO"
fi

# Form API URL
HELM_CHART_MUSEUM_API_URL="$HELM_CHART_ARTIFACTORY_URL/api"

# iterate through artifacts to promote
for artifact in "$@"
do
    RESPONSE_CODE_ARTIFACT="$(get_artifact "$artifact")"
    if [ "$RESPONSE_CODE_ARTIFACT" == "200" ]; then
      promote_artifact "$artifact"
    elif [ "$RESPONSE_CODE_ARTIFACT" == "404" ]; then
      exitmsg "The artifact $artifact was not found in repository $HELM_CHART_MUSEUM_ORIGIN_REPO"
    else
      exitmsg "There was an error retrieving $artifact from repository $HELM_CHART_MUSEUM_ORIGIN_REPO, the http error code was $RESPONSE_CODE_ARTIFACT"
    fi
done

printf '\n'
echo "Promoted Orc8r chart artifacts $* from $HELM_CHART_MUSEUM_ORIGIN_REPO to $HELM_CHART_MUSEUM_DEST_REPO successfully."

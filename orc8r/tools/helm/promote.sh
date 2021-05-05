#!/usr/bin/env bash

# Copyright 2020 The Magma Authors.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# promote.sh copies specified artifacts from one repo to another

set -e -o pipefail


promote_artifact () {
  ARTIFACT="$1"
  curl -X POST -u "$HELM_CHART_MUSEUM_USERNAME":"$HELM_CHART_MUSEUM_TOKEN" --fail \
   "$HELM_CHART_MUSEUM_API_URL/copy/$HELM_CHART_MUSEUM_ORIGIN_REPO/$ARTIFACT?to=/$HELM_CHART_MUSEUM_DEST_REPO/$ARTIFACT"
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

if [[ -z $HELM_CHART_MUSEUM_API_URL ]]; then
  exitmsg "Environment variable HELM_CHART_MUSEUM_API_URL must be set"
fi

if [[ -z $HELM_CHART_MUSEUM_ORIGIN_REPO ]]; then
  exitmsg "Environment variable HELM_CHART_MUSEUM_ORIGIN_REPO must be set"
fi

if [[ -z $HELM_CHART_MUSEUM_DEST_REPO ]]; then
  exitmsg "Environment variable HELM_CHART_MUSEUM_DEST_REPO must be set"
fi

if [[ -z $HELM_CHART_MUSEUM_USERNAME ]]; then
  exitmsg "Environment variable HELM_CHART_MUSEUM_USERNAME must be set"
fi

if [[ -z $HELM_CHART_MUSEUM_TOKEN ]]; then
  exitmsg "Environment variable HELM_CHART_MUSEUM_TOKEN must be set"
fi

# iterate through artifacts to promote
for artifact in "$@"
do
    promote_artifact "$artifact"
done

printf '\n'
echo "Promoted orc8r chart artifacts $* from $HELM_CHART_MUSEUM_ORIGIN_REPO to $HELM_CHART_MUSEUM_DEST_REPO successfully."

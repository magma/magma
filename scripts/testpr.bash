#!/usr/bin/env bash
# Copyright 2020 The Magma Authors.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# codeowners.bash prints information on the current set of Magma codeowners.

function usage() {
    echo ${1}
    exit 1
}

token=${GITHUB_TOKEN}
username=${1}

[[ ${token} == '' ]] && usage 'GITHUB_TOKEN environment variable must be set to a GitHub personal access token'
[[ ${username} == '' ]] && usage 'Usage: codeowners.bash YOUR_GITHUB_USERNAME'

codeowners=$(curl --silent --show-error -u ${username}:${token} \
    -H "Accept: application/vnd.github.v3+json" \
    https://api.github.com/organizations/66266171/team/4477370/members \
    | jq '.[].login' \
    | tr -d '"' \
    | sort
)
n_codeowners=$(echo ${codeowners} | wc -w)
n_majority=$(python3 -c "import math; print(math.floor((${n_codeowners}/2)+1))")

echo ''
echo 'This script pulls the set of codeowners as defined by the magma/magma-maintainers GitHub team.'
echo 'That team is expected to be kept manually updated to reflect the magma/magma CODEOWNERS file.'
echo 'To double-check, consider looking through'
echo '    - CODEOWNERS file: https://github.com/magma/magma/blob/master/CODEOWNERS'
echo '    - All Magma teams: https://github.com/orgs/magma/teams'
echo ''
echo "${codeowners}"
echo ''
echo ${n_codeowners} total codeowners
echo ${n_majority} constitutes majority

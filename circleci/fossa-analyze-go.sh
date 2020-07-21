#!/bin/bash
################################################################################
# Copyright 2020 The Magma Authors.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
################################################################################

if [ "x${MAGMA_ROOT}" = x ]; then
    MAGMA_ROOT=$(pwd)
fi
cd "${MAGMA_ROOT}" || exit 1

declare -a FOSSA_OPTS
FOSSA_OPTS=()
if [ "x${FOSSA_API_KEY}" = x ]; then
    FOSSA_OPTS=("-o")
fi

# break MAGMA_MODULES (str) into an array split on spaces
declare -a module_array
IFS=" " read -r -a module_array <<< "${MAGMA_MODULES}"

for mod in "${module_array[@]}"; do
    if [ -d "${mod}/cloud/go" ]; then
        pushd "${mod}/cloud/go" || { echo error "${mod}"; continue; }
        modfile="$(pwd)/go.mod"
        src=$(dirname "${modfile}")
        proj=$(echo "${src}" | sed -e 's|^\./||g' -e 's|/|_|g')
        fossa init -t "${proj}"
        fossa analyze "${FOSSA_OPTS[@]}"
        popd || echo "should never happen ${src}"
    fi
done

exit 0

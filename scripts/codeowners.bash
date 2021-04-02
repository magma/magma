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

codeowners=$(curl --silent --show-error https://raw.githubusercontent.com/magma/magma/master/CODEOWNERS \
    | egrep --only-matching '@\S+' \
    | grep --invert-match '@magma/' \
    | sort \
    | uniq
)
n_codeowners=$(echo ${codeowners} | wc -w)
n_majority=$(python3 -c "import math; print(math.floor((${n_codeowners}/2)+1))")

echo ${codeowners}
echo ''
echo ${n_codeowners} total codeowners
echo ${n_majority} constitutes majority

#!/bin/bash
# Copyright 2023 The Magma Authors.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

set -eou pipefail

CLANG_FILE=/usr/bin/clang-format-11
if test ! -f "$CLANG_FILE"; then
  echo "clang-format-11 not found. You should execute the script on the dev VM
    or install clang-format-11 on your host machine."
  exit 1
fi

for dir in /orc8r/gateway/c/ /lte/gateway/c/ /lte/gateway/python/;
do
  find "${MAGMA_ROOT}/${dir}" \( -iname "*.c" -o -iname "*.cpp" -o -iname "*.h" -o -iname "*.hpp" \) -exec $CLANG_FILE -i {} \;
done

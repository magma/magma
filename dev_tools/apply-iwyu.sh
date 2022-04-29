#!/bin/bash
################################################################################
# Copyright 2022 The Magma Authors.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
################################################################################

# How to use: 
# 1. $MAGMA_ROOT/dev_tools/apply-iwyu.sh $PATH_TO_APPLY_IWYU
#    Ex. $MAGMA_ROOT/dev_tools/apply-iwyu.sh lte/gateway/c/core/common
SOURCE_CODE_PATH=$1

IWYU_OUTPUT_FILE=/tmp/iwyu-output.txt

echo ""
echo "Generating /tmp/compile_commands.json ..."
echo ""

dev_tools/gen_compilation_database.py --output_dir=/tmp

echo ""
echo "Piping IWYU output into /tmp/iwyu-output.txt ..."
echo ""

iwyu_tool.py -j "$(nproc --all)" -p /tmp/compile_commands.json "$SOURCE_CODE_PATH" -- -Xiwyu --mapping_file="$MAGMA_ROOT"/dev_tools/iwyu.imp | tee "$IWYU_OUTPUT_FILE"

echo ""
echo "Applying /tmp/iwyu-output.txt ..."
echo ""

fix_includes.py < "$IWYU_OUTPUT_FILE" -b --reorder

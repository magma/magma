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
# 1. Run `Bzl: Bazel/C++: Generate Compilation Database` in the VSCode command palette (CMD+Shift+P) to generate the compilation database
# 2. $MAGMA_ROOT/dev_tools/apply-iwyu.sh $PATH_TO_APPLY_IWYU
#    Ex. $MAGMA_ROOT/dev_tools/apply-iwyu.sh lte/gateway/c/core/common
SOURCE_CODE_PATH=$1

IWYU_OUTPUT_FILE=/tmp/iwyu-output.txt

echo ""
echo "[IMPORTANT] Make sure you've run 'Bzl: Bazel/C++: Generate Compilation Database' in VSCode's command palette (CMD+Shift+P) to generate an up-to-date compile_commands.json!"
echo ""

iwyu_tool.py -j "$(nproc --all)" -p "$MAGMA_ROOT"/compile_commands.json "$SOURCE_CODE_PATH" -- -Xiwyu --mapping_file="$MAGMA_ROOT"/dev_tools/iwyu.imp | tee "$IWYU_OUTPUT_FILE"

fix_includes.py < "$IWYU_OUTPUT_FILE" -b --reorder

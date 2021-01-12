#!/bin/bash
# Copyright 2020 The Magma Authors.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# Usage: update_package.sh [pull request id]
set -e
SCRIPT_DIR=$MAGMA_ROOT/lte/gateway/deploy/roles/magma_test/files

# Clone the code
$SCRIPT_DIR/clone_s1_tester.sh

# If any pull request is provided, apply them, else use master
pushd "$S1AP_TESTER_SRC" || exit
git checkout master
git branch | grep -v '^*' | xargs --no-run-if-empty git branch -D

if [ "$#" == 1 ]; then
  pull_req=$1
  git fetch origin pull/$pull_req/head:branch_$pull_req
  git checkout branch_$pull_req
fi

popd

# Link the build binaries again
$SCRIPT_DIR/build_s1_tester.sh

echo "Install pyparsing"
# Re-generate the s1ap_types.py
pip3 show pyparsing 1>/dev/null
if [ $? != 0 ]; then
   echo "Installing pyparsing"
   pip3 install pyparsing
fi

echo "Regenerating s1ap_types.py"
/usr/bin/python3.5 $SCRIPT_DIR/c_parser.py

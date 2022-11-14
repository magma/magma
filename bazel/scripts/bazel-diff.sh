#!/usr/bin/env bash

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

set -euo pipefail

###############################################################################
# FUNCTION DECLARATIONS
###############################################################################

generate_hashes() {
    local GIT_COMMIT=$1
    local HASHES_JSON=$2
    local CURRENT_BRANCH=$(git branch --show-current)
    # TODO: Alternatively use     local CURRENT_COMMIT=$(git rev-parse HEAD)

    echo "Checking out the commit $GIT_COMMIT." 1>&2
    git -C $MAGMA_ROOT checkout $GIT_COMMIT --quiet
    echo "Generating the hashes for the commit $GIT_COMMIT" 1>&2
    $BAZEL_DIFF generate-hashes -w $MAGMA_ROOT -b $BAZEL_PATH $HASHES_JSON

    # Undo the git checkout
    git -C $MAGMA_ROOT checkout $CURRENT_BRANCH --quiet
}

###############################################################################
# SCRIPT SECTION
###############################################################################

BAZEL_PATH=$(which bazel)
BAZEL_DIFF="/tmp/bazel_diff"
HASHES_PRE_JSON="/tmp/hashes_pre.json"
HASHES_POST_JSON="/tmp/hashes_post.json"
IMPACTED_TARGETS_FILE="/tmp/impacted_targets.txt"

# Git SHA of the previous commit
GIT_SHA_PRE=$1
# Git SHA of the commit with the changes
GIT_SHA_POST=$2

echo "Generating the bazel-diff script at $BAZEL_DIFF." 1>&2
$BAZEL_PATH run :bazel-diff --script_path="$BAZEL_DIFF"

# Generate the JSON files containing the hashes 
generate_hashes $GIT_SHA_PRE $HASHES_PRE_JSON
generate_hashes $GIT_SHA_POST $HASHES_POST_JSON

echo "Determining the impacted targets" 1>&2
$BAZEL_DIFF get-impacted-targets -sh $HASHES_PRE_JSON -fh $HASHES_POST_JSON -o $IMPACTED_TARGETS_FILE

echo "The list of affected targets of the changes between commit $GIT_SHA_PRE and $GIT_SHA_POST:" 1>&2
cat $IMPACTED_TARGETS_FILE

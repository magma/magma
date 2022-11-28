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
    local GIT_COMMIT
    GIT_COMMIT=$1
    local HASHES_JSON
    HASHES_JSON=$2

    echo -e "${LOG_HEADER} Creating temporary directory to check out the commit $GIT_COMMIT ..." 1>&2
    local WORKDIR
    WORKDIR=$(mktemp --directory)
    git -c advice.detachedHead=false clone "$MAGMA_ROOT" "$WORKDIR" --quiet

    echo -e "${LOG_HEADER} Checking out the commit $GIT_COMMIT." 1>&2
    git -C "$WORKDIR" -c advice.detachedHead=false checkout "$GIT_COMMIT" --quiet

    echo -e "${LOG_HEADER} Generating the hashes for the commit $GIT_COMMIT" 1>&2
    "$BAZEL_DIFF" generate-hashes --workspacePath "$WORKDIR" --bazelPath "$BAZEL_PATH" "$HASHES_JSON"

    rm -rf "$WORKDIR"
}

###############################################################################
# SCRIPT SECTION
###############################################################################

# ANSI escape sequences for log formatting.
YELLOW='\033[0;33m'
NO_FORMAT='\033[0m'
LOG_HEADER="${YELLOW}BAZEL-DIFF INFO:${NO_FORMAT}"

# BAZEL_PATH is needed for the generate-hashes --bazelPath argument.
BAZEL_PATH=$(command -v bazel)
BAZEL_DIFF="/tmp/bazel_diff"
HASHES_PRE_JSON="/tmp/hashes_pre.json"
HASHES_POST_JSON="/tmp/hashes_post.json"
IMPACTED_TARGETS_FILE="/tmp/impacted_targets_bazel_diff.txt"

# Git SHA of the previous commit.
GIT_SHA_PRE=$1
# Git SHA of the commit with the changes.
GIT_SHA_POST=$2

echo -e "${LOG_HEADER} Generating the bazel-diff script at $BAZEL_DIFF." 1>&2
bazel run :bazel-diff --script_path="$BAZEL_DIFF"

# Generate the JSON files containing the hashes.
generate_hashes "$GIT_SHA_PRE" "$HASHES_PRE_JSON"
generate_hashes "$GIT_SHA_POST" "$HASHES_POST_JSON"

echo -e "${LOG_HEADER} Determining the impacted targets" 1>&2
"$BAZEL_DIFF" get-impacted-targets --startingHashes "$HASHES_PRE_JSON" --finalHashes "$HASHES_POST_JSON" --output "$IMPACTED_TARGETS_FILE"

# Check if the file containing the affected targets is empty.
if [[ -s "$IMPACTED_TARGETS_FILE" ]];
then
    echo -e "${LOG_HEADER} The list of targets impacted by the changes between commits $GIT_SHA_PRE and $GIT_SHA_POST is:" 1>&2
    cat "$IMPACTED_TARGETS_FILE"
else
    echo -e "${LOG_HEADER} No targets are impacted by the changes between commits $GIT_SHA_PRE and $GIT_SHA_POST." 1>&2
fi

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

create_bazel_test_query() {
    # Write beginning of the bazel test query.
    printf '%s' "attr(tags, $BAZEL_FILTER_TAG, kind(\"$BAZEL_FILTER_RULE\"," > "$BAZEL_QUERY_FILE"

    # Append union of targets to the bazel test query.
    while IFS="" read -r impacted_target
    do
        printf '%s ' "$impacted_target union" >> "$BAZEL_QUERY_FILE"
    done < "$IMPACTED_TARGETS_FILE"

    # Remove the last 'union ' in the list of targets and
    # append the end of the bazel test query, which filters out manual test targets.
    sed -i 's/union $/except attr(tags, manual, kind(.*_test, \/\/...))))/' "$BAZEL_QUERY_FILE"
}

run_bazel_query() {
    # A query file is needed to avoid the bazel command line max args limit.
    bazel query --query_file="$BAZEL_QUERY_FILE" > "$IMPACTED_TEST_TARGETS_FILE"
}

###############################################################################
# SCRIPT SECTION
###############################################################################

# ANSI escape sequences for log formatting.
YELLOW='\033[0;33m'
NO_FORMAT='\033[0m'
LOG_HEADER="${YELLOW}INFO:${NO_FORMAT}"

IMPACTED_TARGETS_FILE="/tmp/impacted_targets_pre_filter.txt"
rm -f "$IMPACTED_TARGETS_FILE"
touch "$IMPACTED_TARGETS_FILE"
BAZEL_QUERY_FILE="/tmp/bazel_query_impacted_targets.txt"
IMPACTED_TEST_TARGETS_FILE="/tmp/impacted_test_targets.txt"

# Read stdin into IMPACTED_TARGETS_FILE.
while IFS= read -r impacted_target
do
    echo "$impacted_target" >> "$IMPACTED_TARGETS_FILE"
done

if [[ ! -s "$IMPACTED_TARGETS_FILE" ]];
then
    echo -e "${LOG_HEADER} Input is empty, no targets to be filtered." 1>&2
    exit 0
fi

# Optional script argument, e.g. 'cc_test' or 'py_test', default is all test rules.
BAZEL_FILTER_RULE=${1:-".*_test"}
# Optional script argument for filtering tags, e.g. 'service', default is no filter
BAZEL_FILTER_TAG=${2:-".*"}

create_bazel_test_query
run_bazel_query

# Check if the file containing the affected test targets is empty.
if [[ -s "$IMPACTED_TEST_TARGETS_FILE" ]];
then
    echo -e "${LOG_HEADER} The impacted $BAZEL_FILTER_RULE targets are:" 1>&2
    cat $IMPACTED_TEST_TARGETS_FILE
else
    echo -e "${LOG_HEADER} No $BAZEL_FILTER_RULE targets are impacted by the changes." 1>&2
fi

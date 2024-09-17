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

help() {
    echo "Run the AGW load tests with bazel."
    echo "Usage:"
    echo "   $(basename "$0") --help"
    echo "      Display this help message."
    echo "   $(basename "$0")" 
    echo "      Executes all bazel tests that are tagged as load_test." 
}

build_load_tests() {
    echo "Building load tests..."
    bazel build //lte/gateway/python/load_tests:all
    echo "Finished building load tests."
    echo "#############################"
}

run_load_tests() {
    echo "Running load tests..."
    for LOAD_TEST in "${LOAD_TEST_LIST[@]}";
    do  
        echo "#############################"
        echo "Running load test: ${LOAD_TEST}"
        # shellcheck disable=SC2086
        sudo -E PATH="${PATH}" "${MAGMA_ROOT}/bazel-bin/lte/gateway/python/load_tests/"${LOAD_TEST}
    done
    echo "#############################"
    echo "Finished running load tests."
    echo "The result JSON files (result_<name_of_grpc_function>.json) can be found in /var/tmp."
}

###############################################################################
# SCRIPT SECTION
###############################################################################

if [[ "${1:-}" != "" ]];
then
    help
    exit 0
fi

declare -a LOAD_TEST_LIST=("loadtest_mobilityd allocate" \
                           "loadtest_mobilityd release" \
                           "loadtest_pipelined activate_flows" \
                           "loadtest_pipelined deactivate_flows" \
                           "loadtest_sessiond create" \
                           "loadtest_sessiond end" \
                           "loadtest_subscriberdb add" \
                           "loadtest_subscriberdb list" \
                           "loadtest_subscriberdb delete" \
                           "loadtest_subscriberdb get" \
                           "loadtest_subscriberdb update" \
                           "loadtest_policydb enable_static_rules" \
                           "loadtest_policydb disable_static_rules" \
                           "loadtest_directoryd update_record" \
                           "loadtest_directoryd delete_record" \
                           "loadtest_directoryd get_record" \
                           "loadtest_directoryd get_all_records")

build_load_tests
run_load_tests

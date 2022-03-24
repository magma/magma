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
    echo "Executes all bazel tests that are tagged as sudo_test inside the"
    echo "specified directory (recursively), or one single test if a test name"
    echo "is provided."
    echo "Usage:"
    echo "   $(basename "$0")  # execute all tests in the magma repository" 
    echo "   $(basename "$0") path_to_tests_directory/"
    echo "   $(basename "$0") path_to_tests_directory:test_name"
    exit 1
}

create_test_targets() {
    if [[ "${TARGET_PATH}" == *":"* ]];
    then
        echo "Single target specified - running test:"
        TEST_TARGETS=( "${TARGET_PATH}" )
    else
        echo "Multiple targets specified - running tests:"
        mapfile -t TEST_TARGETS < <(bazel query "attr(tags, sudo_test, kind(py_test, //${TARGET_PATH}...))")
    fi
    if [[ "${#TEST_TARGETS[@]}" -eq 0 ]];
    then
        echo "ERROR: No test found."
        help
        exit 1
    fi
    for TARGET in "${TEST_TARGETS[@]}"
    do
        echo "${TARGET}"
    done
}

run_test() {
    local TARGET=$1
    local TARGET_PATH=${TARGET%:*}
    local SHORT_TARGET=${TARGET#*:}
    (
        set -x
        bazel build "${TARGET}"
        sudo "bazel-bin/${TARGET_PATH}/${SHORT_TARGET}"
    )
}

print_summary() {
    local NUM_SUCCESS=$1
    local TOTAL_TESTS=$2
    echo "SUMMARY: ${NUM_SUCCESS}/${TOTAL_TESTS} tests were successful."
    for TARGET in "${!TEST_RESULTS[@]}"
    do
        echo "  ${TARGET}: ${TEST_RESULTS[${TARGET}]}"
    done
}

###############################################################################
# SCRIPT SECTION
###############################################################################

TARGET_PATH="${1:-}"

declare -a TEST_TARGETS
declare -A TEST_RESULTS
NUM_SUCCESS=0
NUM_RUN=1

cd "${MAGMA_ROOT}"

create_test_targets

TOTAL_TESTS=${#TEST_TARGETS[@]}

for TARGET in "${TEST_TARGETS[@]}"
do
    echo "Starting test ${NUM_RUN}/${TOTAL_TESTS}: ${TARGET}"

    if run_test "${TARGET}";
    then
        NUM_SUCCESS=$((NUM_SUCCESS + 1))
        TEST_RESULTS["${TARGET}"]="PASSED"
    else
        TEST_RESULTS["${TARGET}"]="FAILED"
    fi
    NUM_RUN=$((NUM_RUN + 1))
done

print_summary "${NUM_SUCCESS}" "${TOTAL_TESTS}"

[[ "${TOTAL_TESTS}" == "${NUM_SUCCESS}" ]]

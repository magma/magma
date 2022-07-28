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
    echo "Executes all integration tests."
    echo "Usage:"
    echo "   $(basename "$0")" 
    echo "      Execute all integration tests in the magma repository." 
    echo "   $(basename "$0") path_to_tests_directory:bazel_test_target_name"
    echo "      Execute the specified test."
    echo "   $(basename "$0") --list"
    echo "      List all integration tests."
    echo "   $(basename "$0") --list-extended"
    echo "      List the extended integration tests."
    echo "   $(basename "$0") --list-traffic-server"
    echo "      List all integration tests that use the traffic server."
    echo "   $(basename "$0") --setup-extended"
    echo "      Execute the setup test for the extended tests."
    echo "   $(basename "$0") --teardown-extended"
    echo "      Execute the teardown test for the extended tests."
    echo "   $(basename "$0") --skip-setup-teardown-extended"
    echo "      Execute all integration tests in the magma repository," 
    echo "      except the setup and teardown for extended tests." 
    echo "   $(basename "$0") --skip-setup-teardown-extended path_to_tests_directory:bazel_test_target_name"
    echo "      Execute the specified test, without executing"
    echo "      the setup and teardown for extended tests." 
    echo "   $(basename "$0") --help"
    echo "      Display this help message."
}

categorize_test() {
    local TARGET=$1
    if bazel query "attr(tags, extended_test, kind(py_test, ${TARGET}))";
    then
        EXTENDED_TEST_TARGETS=( "${TARGET}" )
    else
        echo "ERROR: Could not categorize the provided test."
        exit 1
    fi
    # TODO: other tests
}

create_test_targets() {
    if [[ "${TARGET_PATH}" == *":"* ]];
    then
        echo "Single target specified - running test:"
        categorize_test "${TARGET_PATH}"
    else
        echo "Multiple targets specified - running tests:"
        create_extended_test_targets
        # TODO: other tests
    fi
    if [[ "${#EXTENDED_TEST_TARGETS[@]}" -eq 0 ]]; # TODO other tests
    then
        echo "ERROR: No test found."
        help
        exit 1
    fi
    for TARGET in "${EXTENDED_TEST_TARGETS[@]}"
    do
        echo "${TARGET}"
    done
}

create_extended_test_targets() {
    mapfile -t EXTENDED_TEST_TARGETS < <(bazel query "attr(tags, extended_test, kind(py_test, //lte/gateway/python/integ_tests/s1aptests/...))")
}

list_all_tests() {
    echo "All integration tests:"
    list_extended_tests
    # TODO: other tests
    exit 0
}

list_extended_tests() {
    echo "Extended tests:"
    create_extended_test_targets
    for TARGET in "${EXTENDED_TEST_TARGETS[@]}"
    do
        echo "${TARGET}"
    done
}

list_traffic_server_tests() {
    echo "Tests that require the traffic server:"
    bazel query "attr(tags, traffic_server_test, kind(py_test, //lte/gateway/python/integ_tests/s1aptests/...))"
    exit 0
}

setup_extended_tests() {
    echo "Setting up the environment for the extended tests."
    echo "Building..."
    bazel build "//lte/gateway/python/integ_tests/s1aptests:test_modify_mme_config_for_sanity" --define=on_magma_test=1
    echo "Executing..."
    sudo "${MAGMA_ROOT}/bazel-bin/lte/gateway/python/integ_tests/s1aptests/test_modify_mme_config_for_sanity"
    echo "Setup finished successfully."
}

teardown_extended_tests() {
    if [[ -f "${EXTENDED_TEST_CLEANUP_FILE_NAME}" ]];
    then
        echo "Cleaning up the environment after the extended tests."
        echo "Building..."
        bazel build "//lte/gateway/python/integ_tests/s1aptests:test_restore_mme_config_after_sanity" --define=on_magma_test=1
        echo "Executing..."
        sudo "${MAGMA_ROOT}/bazel-bin/lte/gateway/python/integ_tests/s1aptests/test_restore_mme_config_after_sanity"
        echo "Cleanup finished successfully."
    else
        echo "No backup file found, skipping cleanup."
    fi
}

run_test() {
    local TARGET=$1
    local TARGET_PATH=${TARGET%:*}
    local SHORT_TARGET=${TARGET#*:}
    (
        echo "BUILDING TEST: ${TARGET}"
        set -x
        bazel build "${TARGET}" --define=on_magma_test=1
        set +x
        echo "RUNNING TEST: ${TARGET}"
        set -x
        sudo "${MAGMA_ROOT}/bazel-bin/${TARGET_PATH}/${SHORT_TARGET}"
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

declare -a EXTENDED_TEST_TARGETS
declare -A TEST_RESULTS
NUM_SUCCESS=0
NUM_RUN=1

cd "${MAGMA_ROOT}"

EXTENDED_TEST_SETUP="lte/gateway/python/integ_tests/s1aptests:test_modify_mme_config_for_sanity"
EXTENDED_TEST_TEARDOWN="lte/gateway/python/integ_tests/s1aptests:test_restore_mme_config_after_sanity"
SKIP_EXTENDED_SETUP_AND_TEARDOWN="false"
EXTENDED_TEST_CLEANUP_FILE_NAME="${MAGMA_ROOT}/lte/gateway/configs/templates/mme.conf.template.bak"

declare -a POSITIONAL_ARGS

while [[ $# -gt 0 ]]; do
  case $1 in
    --list)
      list_all_tests
      ;;
    --list-extended)
      list_extended_tests
      exit 0
      ;;
    --list-traffic-server)
      list_traffic_server_tests
      ;;
    --setup-extended)
      setup_extended_tests
      exit 0
      ;;
    --teardown-extended)
      teardown_extended_tests
      exit 0
      ;;
    --skip-setup-teardown-extended)
      SKIP_EXTENDED_SETUP_AND_TEARDOWN="true"
      shift
      ;;
    --help)
      help
      exit 0
      ;;
    --*|-*)
      echo "Unknown option $1"
      exit 1
      ;;
    *)
      POSITIONAL_ARGS+=("$1")
      shift
      ;;
  esac
done

set -- "${POSITIONAL_ARGS[@]}"

TARGET_PATH="${1:-}"

if [[ "${TARGET_PATH}" == *"${EXTENDED_TEST_SETUP}" ]];
then
    setup_extended_tests
    exit 0
fi

if [[ "${TARGET_PATH}" == *"${EXTENDED_TEST_TEARDOWN}" ]];
then
    teardown_extended_tests
    exit 0
fi

create_test_targets

TOTAL_TESTS=${#EXTENDED_TEST_TARGETS[@]}
# TODO TOTAL_TESTS=$((TOTAL_TESTS + PRECOMMIT))

if [[ "${SKIP_EXTENDED_SETUP_AND_TEARDOWN}" == "false" ]];
then
    setup_extended_tests
fi

for TARGET in "${EXTENDED_TEST_TARGETS[@]}"
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

if [[ "${SKIP_EXTENDED_SETUP_AND_TEARDOWN}" == "false" ]];
then
    teardown_extended_tests
fi
# TODO other tests

print_summary "${NUM_SUCCESS}" "${TOTAL_TESTS}"

[[ "${TOTAL_TESTS}" == "${NUM_SUCCESS}" ]]

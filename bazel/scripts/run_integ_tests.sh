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
    echo -e "${BOLD}Executes all integration tests."
    echo -e "Usage:${NO_FORMATTING}"
    echo "   $(basename "$0") --help"
    echo "      Display this help message."
    echo "   $(basename "$0")"
    echo "      Execute all precommit and extended integration tests in the magma repository."
    echo "   $(basename "$0") path_to_tests_directory:bazel_test_target_name"
    echo "      Execute the specified test."
    echo "   --retry-on-failure"
    echo "      Retry twice for every test in case of failure."
    echo "   --retry-attempts N"
    echo "      Retry N times for every test in case of failure."
    echo "      Should be used together with --retry-on-failure."
    echo "   --rerun-previously-failed"
    echo "      Rerun all tests that have failed during the previous run."
    echo "   --profile PROFILE_FILE_PATH"
    echo "      Create a bazel profile during build. The profile will be stored in the given file."
    echo -e "${BOLD}List tests:${NO_FORMATTING}"
    echo "   $(basename "$0") --list"
    echo "      List all integration tests."
    echo "   $(basename "$0") --list-precommit"
    echo "      List the precommit integration tests."
    echo "   $(basename "$0") --list-extended"
    echo "      List the precommit integration tests."
    echo "   $(basename "$0") --list-nonsanity"
    echo "      List the nonsanity integration tests."
    echo "   $(basename "$0") --list-traffic-server"
    echo "      List all integration tests that use the traffic server."
    echo -e "${BOLD}Execute precommit:${NO_FORMATTING}"
    echo "   $(basename "$0") --precommit"
    echo "      Run all precommit integration tests."
    echo -e "${BOLD}Execute extended:${NO_FORMATTING}"
    echo "   $(basename "$0") --extended"
    echo "      Run all extended integration tests."
    echo "   $(basename "$0") --setup-extended"
    echo "      Execute the setup test for the extended tests."
    echo "   $(basename "$0") --teardown-extended"
    echo "      Execute the teardown test for the extended tests."
    echo "   $(basename "$0") --skip-setup-teardown-extended"
    echo "      Execute all precommit and extended integration tests in the magma repository," 
    echo "      except the setup and teardown for extended tests."
    echo "      Should be combined with an extended test target or --extended."
    echo "   $(basename "$0") --skip-setup-teardown-extended path_to_tests_directory:bazel_test_target_name"
    echo "      Execute the specified test, without executing"
    echo "      the setup and teardown for extended tests."
    echo -e "${BOLD}Execute nonsanity:${NO_FORMATTING}"
    echo "   $(basename "$0") --nonsanity"
    echo "      Run all nonsanity integration tests."
    echo "   $(basename "$0") --setup-nonsanity"
    echo "      Execute the setup test for the nonsanity tests."
    echo "   $(basename "$0") --teardown-nonsanity"
    echo "      Execute the teardown test for the nonsanity tests."
    echo "   $(basename "$0") --skip-setup-teardown-nonsanity"
    echo "      Execute all precommit and extended integration tests in the magma repository,"
    echo "      except the setup and teardown for nonsanity tests."
    echo "      Should be combined with a nonsanity test target or --nonsanity."
    echo "   $(basename "$0") --skip-setup-teardown-nonsanity --nonsanity"
    echo "      Execute the nonsanity tests, without executing"
    echo "      the setup and teardown for them."
    echo "   $(basename "$0") --skip-setup-teardown-nonsanity path_to_tests_directory:bazel_test_target_name"
    echo "      Execute the specified test, without executing"
    echo "      the setup and teardown for nonsanity tests."
}

categorize_test() {
    local TARGET=$1
    if [[ $(bazel query attr\(tags, precommit_test, kind\(py_test, "${TARGET}"\)\)) == *"${TARGET}" ]];
    then
        PRECOMMIT_TEST_TARGETS+=( "${TARGET}" )
    elif [[ $(bazel query attr\(tags, extended_test, kind\(py_test, "${TARGET}"\)\)) == *"${TARGET}" ]];
    then
        EXTENDED_TEST_TARGETS+=( "${TARGET}" )
    elif [[ $(bazel query attr\(tags, nonsanity_test, kind\(py_test, "${TARGET}"\)\)) == *"${TARGET}" ]];
    then
        NONSANITY_TEST_TARGETS+=( "${TARGET}" )
    else
        echo "ERROR: Could not categorize the provided test."
        exit 1
    fi
}

create_test_targets() {
    local ONLY_FOR_LISTING=${1:-"false"}
    if [[ "${RERUN_PREVIOUSLY_FAILED}" == "true" ]];
    then
        # shellcheck disable=SC2013
        for TARGET_PATH in $(cat "${FAILED_LIST_FILE}");
        do
            categorize_test "${TARGET_PATH}"
        done
    elif (( ${#POSITIONAL_ARGS[@]} > 0 ));
    then
        echo "Custom target list provided - running tests:"
        for TARGET_PATH in "${POSITIONAL_ARGS[@]}"
        do
            if [[ "${TARGET_PATH}" == *":"* ]];
            then
                categorize_test "${TARGET_PATH}"
            else
                echo "ERROR: Invalid test target name."
                exit 1
            fi
        done
    else
        if [[ "${ONLY_FOR_LISTING}" == "false" ]];
        then
            echo "Multiple targets specified - running tests:"
        fi
        if [[ "${RUN_ALL}" == "true" || "${PRECOMMIT}" == "true" ]];
        then
            create_precommit_test_targets
        fi
        if [[ "${RUN_ALL}" == "true" || "${EXTENDED}" == "true" ]];
        then
            create_extended_test_targets
        fi
        if [[ "${RUN_ALL}" == "false" && "${NONSANITY}" == "true" ]];
        then
            create_nonsanity_test_targets
        fi
    fi
    ALL_REQUESTED_TARGETS=( "${PRECOMMIT_TEST_TARGETS[@]}" "${EXTENDED_TEST_TARGETS[@]}" "${NONSANITY_TEST_TARGETS[@]}" )
    for TARGET in "${ALL_REQUESTED_TARGETS[@]}"
    do
        echo "${TARGET}"
    done
    if [[ "${#ALL_REQUESTED_TARGETS[@]}" -eq 0 ]];
    then
        echo "ERROR: No targets found with the given options!"
        exit 1
    fi
}

create_precommit_test_targets() {
    mapfile -t PRECOMMIT_TEST_TARGETS < <(bazel query "attr(tags, precommit_test, kind(py_test, //lte/gateway/python/integ_tests/s1aptests/...))")
}

create_extended_test_targets() {
    mapfile -t EXTENDED_TEST_TARGETS < <(bazel query "attr(tags, extended_test, kind(py_test, //lte/gateway/python/integ_tests/s1aptests/...))")
}

create_nonsanity_test_targets() {
    mapfile -t NONSANITY_TEST_TARGETS < <(bazel query "attr(tags, nonsanity_test, kind(py_test, //lte/gateway/python/integ_tests/s1aptests/...))")
}

list_all_tests() {
    TARGET_PATH=""
    echo "All integration tests:"
    create_test_targets "true"
    exit 0
}

list_precommit_tests() {
    echo "Precommit tests:"
    create_precommit_test_targets
    for TARGET in "${PRECOMMIT_TEST_TARGETS[@]}"
    do
        echo "${TARGET}"
    done
}

list_extended_tests() {
    echo "Extended tests:"
    create_extended_test_targets
    for TARGET in "${EXTENDED_TEST_TARGETS[@]}"
    do
        echo "${TARGET}"
    done
}

list_nonsanity_tests() {
    echo "Nonsanity tests:"
    create_nonsanity_test_targets
    for TARGET in "${NONSANITY_TEST_TARGETS[@]}"
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
    if sudo "${MAGMA_ROOT}/bazel-bin/lte/gateway/python/integ_tests/s1aptests/test_modify_mme_config_for_sanity";
    then
        echo "Setup finished successfully."
    else
        return 1
    fi
}

teardown_extended_tests() {
    if [[ -f "${TEST_CLEANUP_FILE_NAME}" ]];
    then
        echo "Cleaning up the environment after the extended tests."
        echo "Building..."
        bazel build "//lte/gateway/python/integ_tests/s1aptests:test_restore_mme_config_after_sanity" --define=on_magma_test=1
        echo "Executing..."
        if sudo "${MAGMA_ROOT}/bazel-bin/lte/gateway/python/integ_tests/s1aptests/test_restore_mme_config_after_sanity";
        then
            echo "Cleanup finished successfully."
        else
            return 1
        fi
    else
        echo "No backup file found, skipping cleanup."
    fi
}

setup_nonsanity_tests() {
    echo "Setting up the environment for the nonsanity tests."
    echo "Building..."
    bazel build "//lte/gateway/python/integ_tests/s1aptests:test_modify_config_for_non_sanity" --define=on_magma_test=1
    echo "Executing..."
    if sudo "${MAGMA_ROOT}/bazel-bin/lte/gateway/python/integ_tests/s1aptests/test_modify_config_for_non_sanity";
    then
        echo "Setup finished successfully."
    else
        return 1
    fi
}

teardown_nonsanity_tests() {
    if [[ -f "${TEST_CLEANUP_FILE_NAME}" ]];
    then
        echo "Cleaning up the environment after the nonsanity tests."
        echo "Building..."
        bazel build "//lte/gateway/python/integ_tests/s1aptests:test_restore_config_after_non_sanity" --define=on_magma_test=1
        echo "Executing..."
        if sudo "${MAGMA_ROOT}/bazel-bin/lte/gateway/python/integ_tests/s1aptests/test_restore_config_after_non_sanity";
        then
            echo "Cleanup finished successfully."
        else
            return 1
        fi
    else
        echo "No backup file found, skipping cleanup."
    fi
}

run_test_batch() {
    for TARGET in "${TEST_BATCH_TO_RUN[@]}"
    do
        echo "Starting test ${NUM_RUN}/${TOTAL_TESTS}: ${TARGET}"
        if run_test "${TARGET}";
        then
            NUM_SUCCESS=$((NUM_SUCCESS + 1))
            TEST_RESULTS["${TARGET}"]="${GREEN}PASSED${NO_FORMATTING}"
        else
            TEST_RESULTS["${TARGET}"]="${RED}FAILED${NO_FORMATTING}"
            FAILED_LIST+=( "${TARGET}" )
        fi
        NUM_RUN=$((NUM_RUN + 1))
    done
}

build_tests() {
    (
        set -x
        bazel build "${ALL_REQUESTED_TARGETS[@]}" --define=on_magma_test=1 "$PROFILE"
    )
}

run_test() {
    local TARGET=$1
    local TARGET_PATH=${TARGET%:*}
    local SHORT_TARGET=${TARGET#*:}
    (
        echo "RUNNING TEST: ${TARGET}"
        set -x
        sudo "${MAGMA_ROOT}/bazel-bin/${TARGET_PATH}/${SHORT_TARGET}" "${FLAKY_ARGS[@]}" \
            --junit-xml="${INTEGTEST_REPORT_FOLDER}/${SHORT_TARGET}_report.xml" \
            -o "junit_suite_name=${SHORT_TARGET}" -o "junit_logging=no";
    )
}

create_xml_report() {
    local MERGED_REPORT_XML="integtests_report.xml"
    rm -f "${MERGED_REPORT_FOLDER}/${MERGED_REPORT_XML}"
    mkdir -p "${INTEGTEST_REPORT_FOLDER}"
    python3 lte/gateway/python/scripts/runtime_report.py -i "[^\/]+\.xml" -w "${INTEGTEST_REPORT_FOLDER}" -o "${MERGED_REPORT_FOLDER}/${MERGED_REPORT_XML}" 
    sudo rm -f "${INTEGTEST_REPORT_FOLDER}/"*.xml
}

print_summary() {
    local NUM_SUCCESS=$1
    local TOTAL_TESTS=$2
    echo "The bazel profile was saved to the following location: ${PROFILE#*=}"
    echo "SUMMARY: ${NUM_SUCCESS}/${TOTAL_TESTS} tests were successful."
    for TARGET in "${!TEST_RESULTS[@]}"
    do
        echo -e "  ${TARGET}: ${TEST_RESULTS[${TARGET}]}"
    done
}

###############################################################################
# SCRIPT SECTION
###############################################################################

PRECOMMIT_TEST_TARGETS=()
EXTENDED_TEST_TARGETS=()
NONSANITY_TEST_TARGETS=()
FLAKY_ARGS=()
declare -A TEST_RESULTS
NUM_SUCCESS=0
NUM_RUN=1
PRECOMMIT="false"
EXTENDED="false"
NONSANITY="false"
RUN_ALL="true"

cd "${MAGMA_ROOT}"

EXTENDED_TEST_SETUP="lte/gateway/python/integ_tests/s1aptests:test_modify_mme_config_for_sanity"
EXTENDED_TEST_TEARDOWN="lte/gateway/python/integ_tests/s1aptests:test_restore_mme_config_after_sanity"
SKIP_EXTENDED_SETUP_AND_TEARDOWN="false"
NONSANITY_TEST_SETUP="lte/gateway/python/integ_tests/s1aptests:test_modify_config_for_non_sanity"
NONSANITY_TEST_TEARDOWN="lte/gateway/python/integ_tests/s1aptests:test_restore_config_after_non_sanity"
SKIP_NONSANITY_SETUP_AND_TEARDOWN="false"
TEST_CLEANUP_FILE_NAME="${MAGMA_ROOT}/lte/gateway/configs/templates/mme.conf.template.bak"
RETRY_ON_FAILURE="false"
RETRY_ATTEMPTS=2
RERUN_PREVIOUSLY_FAILED="false"
FAILED_LIST=()
FAILED_LIST_FILE="/tmp/last_failed_integration_tests.txt"
MERGED_REPORT_FOLDER="/var/tmp/test_results"
INTEGTEST_REPORT_FOLDER="${MERGED_REPORT_FOLDER}/integtest_reports"
PROFILE="--profile=/var/tmp/bazel_profile_lte_integ_tests"

BOLD='\033[1m'
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
NO_FORMATTING='\033[0m'

declare -a POSITIONAL_ARGS=()

while [[ $# -gt 0 ]]; do
  case $1 in
    --list)
      list_all_tests
      ;;
    --list-precommit)
      list_precommit_tests
      exit 0
      ;;
    --list-extended)
      list_extended_tests
      exit 0
      ;;
    --list-nonsanity)
      list_nonsanity_tests
      exit 0
      ;;
    --list-traffic-server)
      list_traffic_server_tests
      ;;
    --precommit)
      RUN_ALL="false"
      PRECOMMIT="true"
      shift
      ;;
    --extended)
      RUN_ALL="false"
      EXTENDED="true"
      shift
      ;;
    --nonsanity)
      RUN_ALL="false"
      NONSANITY="true"
      shift
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
    --setup-nonsanity)
      setup_nonsanity_tests
      exit 0
      ;;
    --teardown-nonsanity)
      teardown_nonsanity_tests
      exit 0
      ;;
    --skip-setup-teardown-nonsanity)
      SKIP_NONSANITY_SETUP_AND_TEARDOWN="true"
      shift
      ;;
    --retry-on-failure)
      RETRY_ON_FAILURE="true"
      shift
      ;;
    --retry-attempts)
      shift
      RETRY_ATTEMPTS="$1"
      shift
      ;;
    --rerun-previously-failed)
      RERUN_PREVIOUSLY_FAILED="true"
      break
      ;;
    --profile)
      shift
      PROFILE="--profile=$1"
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

for TARGET_PATH in "${POSITIONAL_ARGS[@]}"
do
    if [[ "${TARGET_PATH}" == *"${EXTENDED_TEST_SETUP}" ]];
    then
        setup_extended_tests
        exit 0
    elif [[ "${TARGET_PATH}" == *"${EXTENDED_TEST_TEARDOWN}" ]];
    then
        teardown_extended_tests
        exit 0
    elif [[ "${TARGET_PATH}" == *"${NONSANITY_TEST_SETUP}" ]];
    then
        setup_nonsanity_tests
        exit 0
    elif [[ "${TARGET_PATH}" == *"${NONSANITY_TEST_TEARDOWN}" ]];
    then
        teardown_nonsanity_tests
        exit 0
    fi
done

create_test_targets

TOTAL_TESTS=${#PRECOMMIT_TEST_TARGETS[@]}
TOTAL_TESTS=$((TOTAL_TESTS + ${#EXTENDED_TEST_TARGETS[@]}))
TOTAL_TESTS=$((TOTAL_TESTS + ${#NONSANITY_TEST_TARGETS[@]}))

if [[ "${RETRY_ON_FAILURE}" == "true" ]];
then
    FLAKY_ARGS=( --force-flaky --no-flaky-report "--max-runs=$((RETRY_ATTEMPTS + 1))" "--min-passes=1" )
fi

echo "#######################################"
echo "BUILDING REQUESTED TESTS"
echo "#######################################"
build_tests

declare -a TEST_BATCH_TO_RUN

if [[ ${#PRECOMMIT_TEST_TARGETS[@]} -gt 0 ]];
then
    echo "#######################################"
    echo "PRECOMMIT TESTS"
    echo "#######################################"
    TEST_BATCH_TO_RUN=( "${PRECOMMIT_TEST_TARGETS[@]}" )
    run_test_batch
fi

if [[ ${#EXTENDED_TEST_TARGETS[@]} -gt 0 ]];
then
    echo "#######################################"
    echo "EXTENDED TESTS"
    echo "#######################################"
    EXTENDED_SETUP="success"
    if [[ "${SKIP_EXTENDED_SETUP_AND_TEARDOWN}" == "false" ]];
    then
        if setup_extended_tests;
        then
            EXTENDED_SETUP="success"
        else
            EXTENDED_SETUP="failed"
            echo -e "The setup ${RED}failed${NO_FORMATTING}. Skipping extended tests."
            TEST_RESULTS["SETUP FOR THE EXTENDED TESTS"]="${RED}FAILED ${YELLOW}(EXTENDED TESTS SKIPPED)${NO_FORMATTING}"
        fi
    fi

    if [[ "${EXTENDED_SETUP}" == "success" ]];
    then
        TEST_BATCH_TO_RUN=( "${EXTENDED_TEST_TARGETS[@]}" )
        run_test_batch
    fi

    if [[ "${SKIP_EXTENDED_SETUP_AND_TEARDOWN}" == "false" ]];
    then
        if ! teardown_extended_tests;
        then
            echo "Teardown ${RED}failed${NO_FORMATTING}."
        fi
    fi
fi

if [[ ${#NONSANITY_TEST_TARGETS[@]} -gt 0 ]];
then
    echo "#######################################"
    echo "NONSANITY TESTS"
    echo "#######################################"
    NONSANITY_SETUP="success"
    if [[ "${SKIP_NONSANITY_SETUP_AND_TEARDOWN}" == "false" ]];
    then
        if setup_nonsanity_tests;
        then
            NONSANITY_SETUP="success"
        else
            NONSANITY_SETUP="failed"
            echo -e "The setup ${RED}failed${NO_FORMATTING}. Skipping nonsanity tests."
            TEST_RESULTS["SETUP FOR THE NONSANITY TESTS"]="${RED}FAILED ${YELLOW}(NONSANITY TESTS SKIPPED)${NO_FORMATTING}"
        fi
    fi

    if [[ "${NONSANITY_SETUP}" == "success" ]];
    then
        TEST_BATCH_TO_RUN=( "${NONSANITY_TEST_TARGETS[@]}" )
        run_test_batch
    fi
    
    if [[ "${SKIP_NONSANITY_SETUP_AND_TEARDOWN}" == "false" ]];
    then
        if ! teardown_nonsanity_tests;
        then
            echo "Teardown ${RED}failed${NO_FORMATTING}."
        fi
    fi
fi

echo "${FAILED_LIST[@]}" > "${FAILED_LIST_FILE}"

create_xml_report

print_summary "${NUM_SUCCESS}" "${TOTAL_TESTS}"

[[ "${TOTAL_TESTS}" == "${NUM_SUCCESS}" ]]

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
    echo -e "${BOLD}Run sudo tests with Bazel."
    echo -e "Usage:${NO_FORMATTING}"
    echo "   $(basename "$0") --help"
    echo "      Display this help message."
    echo "   $(basename "$0")"
    echo "      Executes all bazel tests that are tagged as sudo_test."
    echo "   $(basename "$0") path_to_tests_directory/"
    echo "      Executes all bazel tests that are tagged as sudo_test"
    echo "      inside the specified directory (recursively)."
    echo "   $(basename "$0") path_to_tests_directory:test_name"
    echo "      Executes the specified sudo test."
    echo "   --list"
    echo "      List all sudo test targets."
    echo "   --retry-on-failure"
    echo "      Retry twice for every test in case of failure."
    echo "   --retry-attempts N"
    echo "      Retry N times for every test in case of failure."
    echo "      Should be used together with --retry-on-failure."
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
        sudo "bazel-bin/${TARGET_PATH}/${SHORT_TARGET}" "${FLAKY_ARGS[@]}" \
            --junit-xml="${SUDO_TEST_REPORT_FOLDER}/${SHORT_TARGET}_report.xml" \
            -o "junit_suite_name=${SHORT_TARGET}" -o "junit_logging=no";
    )
}

create_xml_report() {
    local MERGED_REPORT_XML="sudotests_report.xml"
    rm -f "${MERGED_REPORT_FOLDER}/${MERGED_REPORT_XML}"
    mkdir -p "${SUDO_TEST_REPORT_FOLDER}"
    python3 lte/gateway/python/scripts/runtime_report.py -i "[^\/]+\.xml" -w "${SUDO_TEST_REPORT_FOLDER}" -o "${MERGED_REPORT_FOLDER}/${MERGED_REPORT_XML}"
    sudo rm -f "${SUDO_TEST_REPORT_FOLDER}/"*.xml
}

print_summary() {
    local NUM_SUCCESS=$1
    local TOTAL_TESTS=$2
    echo "SUMMARY: ${NUM_SUCCESS}/${TOTAL_TESTS} tests were successful."
    for TARGET in "${!TEST_RESULTS[@]}"
    do
        echo -e "  ${TARGET}: ${TEST_RESULTS[${TARGET}]}"
    done
}

###############################################################################
# SCRIPT SECTION
###############################################################################

TARGET_PATH=""
declare -a TEST_TARGETS
declare -A TEST_RESULTS
FLAKY_ARGS=()
NUM_SUCCESS=0
NUM_RUN=1
RETRY_ON_FAILURE="false"
RETRY_ATTEMPTS=2
MERGED_REPORT_FOLDER="/var/tmp/test_results"
SUDO_TEST_REPORT_FOLDER="${MERGED_REPORT_FOLDER}/sudotest_reports"

BOLD='\033[1m'
RED='\033[0;31m'
GREEN='\033[0;32m'
NO_FORMATTING='\033[0m'

while [[ $# -gt 0 ]]; do
  case $1 in
    --list)
      bazel query "attr(tags, sudo_test, kind(py_test, //...))" 2>/dev/null
      exit 0
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
    --help)
      help
      exit 0
      ;;
    --*|-*)
      echo "Unknown option $1"
      exit 1
      ;;
    *)
      TARGET_PATH="$1"
      shift
      ;;
  esac
done

if [[ "${RETRY_ON_FAILURE}" == "true" ]];
then
    FLAKY_ARGS=( --force-flaky --no-flaky-report "--max-runs=$((RETRY_ATTEMPTS + 1))" "--min-passes=1" )
fi

cd "${MAGMA_ROOT}"

# Build the dhcp_helper_cli script and create a symlink for its intended usage.
# This is needed for several sudo tests.
if [[ ! -L "/usr/local/bin/dhcp_helper_cli.py" ]];
then
  bazel build //lte/gateway/python/dhcp_helper_cli:dhcp_helper_cli
  sudo ln -s "${MAGMA_ROOT}"/bazel-bin/lte/gateway/python/dhcp_helper_cli/dhcp_helper_cli.py /usr/local/bin/dhcp_helper_cli.py
fi

create_test_targets

TOTAL_TESTS=${#TEST_TARGETS[@]}

for TARGET in "${TEST_TARGETS[@]}"
do
    echo "Starting test ${NUM_RUN}/${TOTAL_TESTS}: ${TARGET}"

    if run_test "${TARGET}";
    then
        NUM_SUCCESS=$((NUM_SUCCESS + 1))
        TEST_RESULTS["${TARGET}"]="${GREEN}PASSED${NO_FORMATTING}"
    else
        TEST_RESULTS["${TARGET}"]="${RED}FAILED${NO_FORMATTING}"
    fi
    NUM_RUN=$((NUM_RUN + 1))
done

create_xml_report

print_summary "${NUM_SUCCESS}" "${TOTAL_TESTS}"

[[ "${TOTAL_TESTS}" == "${NUM_SUCCESS}" ]]

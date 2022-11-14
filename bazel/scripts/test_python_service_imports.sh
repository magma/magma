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
    echo "Tests all python services for missing imports or the ones found in"
    echo "the specified directory (recursively), or one single service if a"
    echo "service name is provided."
    echo "Usage:"
    echo "   $(basename "$0")  # test all python services in the magma repository"
    echo "   $(basename "$0") path_to_services_directory/"
    echo "   $(basename "$0") path_to_services_directory:service_name"
    exit 1
}

collect_services() {
    if [[ "${SERVICE_PATH}" == *":"* ]];
    then
        echo "Single service specified:"
        SERVICES=( "${SERVICE_PATH}" )
    else
        echo "Multiple services specified:"
        mapfile -t SERVICES < <(bazel query "attr(tags, service, kind(py_binary, //${SERVICE_PATH}...))")
    fi
    if [[ "${#SERVICES[@]}" -eq 0 ]];
    then
        echo "ERROR: No services found."
        help
        exit 1
    fi
    for SERVICE in "${SERVICES[@]}"
    do
        echo "${SERVICE}"
    done
}

print_summary() {
    local NUM_SUCCESS=$1
    local TOTAL_TESTS=$2
    echo "SUMMARY: ${NUM_SUCCESS}/${TOTAL_TESTS} succeeded."
    for SERVICE in "${!TEST_RESULTS[@]}"
    do
        RESULT=${TEST_RESULTS[${SERVICE}]}
        RESULT="${RESULT#ModuleNotFoundError: }"
        echo "  ${SERVICE}: ${RESULT}"
    done
}

###############################################################################
# SCRIPT SECTION
###############################################################################

SERVICE_PATH="${1:-}"

echo "Collecting services"
declare -a SERVICES
declare -A TEST_RESULTS
LOGGING_FILE="/tmp/test_import_logging.log"
NUM_SUCCESS=0

collect_services

for SERVICE in "${SERVICES[@]}"
do
    echo "Building service: ${SERVICE}"
    bazel build --remote_download_toplevel "${SERVICE}"
# 'bazel build' needs the --remote_download_toplevel option here in order to have
# the same options as the 'bazel run' below. For 'bazel run' the option is set
# in the bazelrc, due to https://github.com/bazelbuild/bazel/issues/11920
    echo "Testing service: ${SERVICE}"
    if timeout --signal 9 --preserve-status 5 bazel run "${SERVICE}" 2>&1  | tee "${LOGGING_FILE}";
    then
        echo "Service successfully started."
        NUM_SUCCESS=$((NUM_SUCCESS + 1))
        TEST_RESULTS["${SERVICE}"]="PASSED"
    else
        echo "Checking if ModuleNotFoundError is present in logs."
        RESULT=$(grep -m 1 "ModuleNotFoundError" "${LOGGING_FILE}" || [[ $? == 1 ]])
        if [[ "${RESULT}" != "" ]];
        then
            echo "ModuleNotFoundError found in logs from ${SERVICE}:"
            echo "${RESULT}"
            TEST_RESULTS["${SERVICE}"]="${RESULT}"
        else
            echo "ModuleNotFoundError not found in logs from ${SERVICE}:"
            NUM_SUCCESS=$((NUM_SUCCESS + 1))
            TEST_RESULTS["${SERVICE}"]="PASSED"
        fi
    fi
    echo -e "\n################################################################################\n"
done

print_summary "${NUM_SUCCESS}" "${#SERVICES[@]}"

[[ "${NUM_SUCCESS}" == "${#SERVICES[@]}" ]]

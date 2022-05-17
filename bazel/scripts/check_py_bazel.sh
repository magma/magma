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
# VARIABLES SECTION
###############################################################################

# Folders and files that should not be relevant for bazel.
DENY_LIST_NOT_RELEVANT=(
  "./ci-scripts"
  "./cn/deploy/scripts"
  "./cwf/gateway/deploy"
  "./dp/tools"
  "./dev_tools"
  "./example"
  "./feg/gateway/docker"
  "./lte/gateway/dev_tools.py"
  "./lte/gateway/deploy"
  "./lte/gateway/python/magma/tests/pylint_wrapper.py"
  "./lte/gateway/python/precommit.py"
  "./orc8r/cloud/deploy"
  "./orc8r/cloud/docker"
  "./orc8r/tools"
  "./protos"
  "./show-tech"
  "./third_party"
  "./xwf/gateway/deploy"
)

# Folders and files that are relevant for building with bazel.
# This list needs to be updated if respected structures are bazelified.
DENY_LIST_NOT_YET_BAZELIFIED=(
  "./dp/cloud/python/magma"
  "./lte/gateway/c/core/oai/tasks/s1ap/messages/asn1"
  "./lte/gateway/python/integ_tests"
  "./lte/gateway/python/load_tests"
  "./lte/gateway/python/scripts"
  "./orc8r/gateway/python/scripts"
  "./orc8r/gateway/python/magma/magmad/tests/dummy_service.py"
  "./orc8r/gateway/python/magma/magmad/upgrade/docker_upgrader.py"
  "./orc8r/gateway/python/magma/common/health/docker_health_service.py"
  "./orc8r/gateway/python/magma/common/health/health_service.py"
  "./orc8r/gateway/python/magma/common/health/entities.py"
  "./lte/gateway/python/magma/health/health_service.py"
  "./lte/gateway/python/magma/health/entities.py"
  "./lte/gateway/python/magma/mobilityd/tests/rpc_servicer_tests.py"
  "./lte/gateway/python/magma/pkt_tester/tests/test_topology_builder.py"
  "./lte/gateway/python/magma/pkt_tester/tests/test_ovs_gtp.py"
  "./lte/gateway/python/magma/pkt_tester/topology_builder.py"
  "./lte/gateway/python/magma/pipelined/tests/script/gtp-packet.py"
  "./lte/gateway/python/magma/pipelined/tests/script/ip-packet.py"
  "./lte/gateway/python/magma/pipelined/tests/envoy-tests/http-serve.py"
  "./lte/gateway/python/magma/pipelined/openflow/events.py"
  "./lte/gateway/python/magma/pipelined/pg_set_session_msg.py"
  "./lte/gateway/python/magma/enodebd/tr069/tests/models_tests.py"
  "./lte/gateway/python/magma/enodebd/state_machines/enb_acs_pointer.py"
)

DENY_LIST=( "${DENY_LIST_NOT_RELEVANT[@]}" "${DENY_LIST_NOT_YET_BAZELIFIED[@]}" )

# Files that exists in multiple folders but ar generally not covered by bazel.
NOT_RELEVANT_FILES=(
  "__init__.py"
  "setup.py"
  "pylint_tests.py"
  "fabfile.py"
)

###############################################################################
# FUNCTIONS SECTION
###############################################################################

get_all_py_files() {
  DENY=()
  FIRST_ITERATION=true
  for entry in "${DENY_LIST[@]}"
  do
    if [[ "${FIRST_ITERATION}" = true ]]
    then
      DENY+=( "-path" "$entry" )
      FIRST_ITERATION=false
    else
      DENY+=( "-o" "-path" "$entry" )
    fi
  done
  find . \( "${DENY[@]}" \) -prune -o -iname "*.py" -print0
}

is_file_not_relevant() {
  local file=$1
  if [[ " ${NOT_RELEVANT_FILES[*]} " =~ $file ]]
  then
    return 0
  fi
  return 1
}

check_py_file() {
  local file=$1
  PY_PATH=$(dirname "$file")
  PY_FILE=$(basename "$file")

  if is_file_not_relevant "${PY_FILE}"
  then
    return
  fi

  BUILD_FILE="${PY_PATH}/BUILD.bazel"
  if [[ -f "${BUILD_FILE}" ]]
  then
    if ! grep -q "${PY_FILE}" "${BUILD_FILE}"
    then
      echo "$file"
    fi
  else
    echo "$file"
  fi
}

check_py_files() {
  while IFS= read -r -d '' file
  do
    check_py_file "$file"
  done 
}

report_problematic_files() {
  local files
  files="$(cat)"

  if [[ -z "$files" ]]
  then
    echo "All python files are either covered by a BUILD.bazel file or excluded from this check."
    exit 0
  else
    cat <<EOF
The following files are not covered by a BUILD.bazel files:

$files

Either add the files to the bazel build system or to the deny list in $0.
Feel free to get support in slack #bazel.
EOF
    exit 1
  fi
}

###############################################################################
# SCRIPT SECTION
###############################################################################

get_all_py_files | check_py_files | report_problematic_files

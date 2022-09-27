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
  "./dp/cloud/python/magma"
  "./dp/tools"
  "./dev_tools"
  "./example"
  "./feg/gateway/docker"
  "./lte/gateway/dev_tools.py"
  "./lte/gateway/deploy"
  # manually executed script
  "./lte/gateway/c/core/oai/tasks/s1ap/messages/asn1/asn1tostruct.py"
  "./lte/gateway/python/build"
  # used for manual testing
  "./lte/gateway/python/magma/pipelined/tests/envoy-tests/http-serve.py"
  "./lte/gateway/python/magma/tests/pylint_wrapper.py"
  "./lte/gateway/python/precommit.py"
  "./orc8r/cloud/deploy"
  "./orc8r/cloud/docker"
  # is not relevant for AGW
  "./orc8r/gateway/python/magma/magmad/upgrade/docker_upgrader.py"
  "./orc8r/tools"
  "./protos"
  "./show-tech"
  "./third_party"
  "./hil_testing"
  "./.git"
)

# Folders and files that are relevant for building with bazel.
# This list needs to be updated if respected structures are bazelified.
DENY_LIST_NOT_YET_BAZELIFIED=(
  # TODO: GH12752 tests should be bazelified
  "./lte/gateway/python/integ_tests/cloud"
  "./lte/gateway/python/integ_tests/cloud_tests"
  "./lte/gateway/python/integ_tests/federated_tests"
  "./lte/gateway/python/integ_tests/s1aptests/workflow"
  "./lte/gateway/python/integ_tests/s1aptests/test_enable_ipv6_iface.py"
  "./lte/gateway/python/integ_tests/s1aptests/test_ipv6_non_nat_dp_ul_tcp.py"
  "./lte/gateway/python/integ_tests/s1aptests/test_disable_ipv6_iface.py"
  #  as of 2022.09.19 commented out in make
  "./lte/gateway/python/integ_tests/s1aptests/test_agw_offload_mixed_idle_active_multiue.py"
  "./lte/gateway/python/integ_tests/s1aptests/test_attach_detach_two_pdns_with_tcptraffic.py"
  "./lte/gateway/python/integ_tests/s1aptests/test_attach_dl_tcp_data_multi_ue.py"
  "./lte/gateway/python/integ_tests/s1aptests/test_attach_dl_udp_data_multi_ue.py"
  "./lte/gateway/python/integ_tests/s1aptests/test_attach_dl_ul_tcp_data_multi_ue.py"
  "./lte/gateway/python/integ_tests/s1aptests/test_attach_standalone_act_dflt_ber_ctxt_rej_ded_bearer_activation.py"
  "./lte/gateway/python/integ_tests/s1aptests/test_attach_ul_tcp_data_multi_ue.py"
  "./lte/gateway/python/integ_tests/s1aptests/test_attach_ul_udp_data_multi_ue.py"
  "./lte/gateway/python/integ_tests/s1aptests/test_attach_with_multiple_mme_restarts.py"
  "./lte/gateway/python/integ_tests/s1aptests/test_data_flow_after_service_request.py"
  "./lte/gateway/python/integ_tests/s1aptests/test_ipv4v6_non_nat_ded_bearer_dl_tcp.py"
  "./lte/gateway/python/integ_tests/s1aptests/test_ipv4v6_non_nat_ded_bearer_ul_tcp.py"
  "./lte/gateway/python/integ_tests/s1aptests/test_ipv4v6_non_nat_ul_tcp.py"
  "./lte/gateway/python/integ_tests/s1aptests/test_ipv6_non_nat_ded_bearer_dl_tcp.py"
  "./lte/gateway/python/integ_tests/s1aptests/test_ipv6_non_nat_ded_bearer_ul_tcp.py"
  "./lte/gateway/python/integ_tests/s1aptests/test_ipv6_non_nat_dp_dl_tcp.py"
  "./lte/gateway/python/integ_tests/s1aptests/test_ipv6_non_nat_dp_dl_udp.py"
  "./lte/gateway/python/integ_tests/s1aptests/test_outoforder_erab_setup_rsp_default_bearer.py"
  "./lte/gateway/python/integ_tests/s1aptests/test_scalability_attach_detach_multi_ue.py"
  "./lte/gateway/python/integ_tests/s1aptests/test_stateless_multi_ue_mixedstate_mme_restart.py"
  # TODO: GH12754 move to (lte|orc8r)/gateway/python/scripts/
  "./orc8r/gateway/python/magma/common/health/docker_health_service.py"
  "./orc8r/gateway/python/magma/common/health/health_service.py"
  "./orc8r/gateway/python/magma/common/health/entities.py"
  "./lte/gateway/python/magma/health/health_service.py"
  "./lte/gateway/python/magma/health/entities.py"
  "./lte/gateway/python/magma/pipelined/pg_set_session_msg.py"
  # TODO: GH12755 access via absolut path on the VM,
  # this needs to be refactored when make is not used anymore
  "./lte/gateway/python/magma/pipelined/tests/script/gtp-packet.py"
  "./lte/gateway/python/magma/pipelined/tests/script/ip-packet.py"
  # TODO: GH9878 needs to be further analyzed
  "./lte/gateway/python/load_tests"
  "./lte/gateway/python/scripts"
  "./orc8r/gateway/python/scripts"
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

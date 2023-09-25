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
  "./build"
  # the following files are test suites that are not to be modeled by bazel
  "./lte/gateway/c/core/oai/test/amf/test_amf.cpp"
  "./lte/gateway/c/core/oai/test/sgw_s8_task/sgw_s8_test.cpp"
  "./lte/gateway/c/core/oai/test/ngap/ngap_test.cpp"
  "./lte/gateway/c/core/oai/test/s1ap_task/s1ap_test.cpp"
  "./lte/gateway/c/core/oai/test/mme_app_task/mme_app_test.cpp"
)

# Folders and files that are relevant for building with bazel.
# This list needs to be updated if respected structures are bazelified.
DENY_LIST_NOT_YET_BAZELIFIED=(
  # TODO: GH12755 access via absolut path on the VM,
  # this needs to be refactored when make is not used anymore
  "./lte/gateway/python/magma/pipelined/ebpf/ebpf_ul_handler.c"
  "./lte/gateway/python/magma/pipelined/ebpf/ebpf_dl_handler.c"
)

DENY_LIST=( "${DENY_LIST_NOT_RELEVANT[@]}" "${DENY_LIST_NOT_YET_BAZELIFIED[@]}" )

C_ROOT="."

declare -a BUILD_FILES
mapfile -t BUILD_FILES < <(find "${C_ROOT}" -name "BUILD.bazel")

###############################################################################
# FUNCTIONS SECTION
###############################################################################

get_all_c_cpp_files() {
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
  find "${C_ROOT}" \( "${DENY[@]}" \) -prune -o \( -iname "*.cpp" -or -iname "*.c" -or -iname "*.h" -or -iname "*.hpp" \) -print0
}

check_c_cpp_files() {
  while IFS= read -r -d '' file
  do
    check_c_cpp_file "$file"
  done 
}

check_c_cpp_file() {
  local file=$1
  SRC_FILE=$(basename "$file")

  if ! (grep -F -q "$SRC_FILE" "${BUILD_FILES[@]}")
  then
    echo "$file"
  fi
}

report_problematic_files() {
  local files
  files="$(cat)"

  if [[ -z "$files" ]]
  then
    echo "All c and c++ files are either covered by a BUILD.bazel file or excluded from this check."
    exit 0
  else
    cat <<EOF
The following files are not covered by a BUILD.bazel file:

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

get_all_c_cpp_files | check_c_cpp_files | report_problematic_files

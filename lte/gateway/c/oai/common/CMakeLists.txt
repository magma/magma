# Copyright 2020 The Magma Authors.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

cmake_minimum_required(VERSION 3.7.2)

set(CMAKE_CXX_STANDARD 14)
set(CMAKE_CXX_STANDARD_REQUIRED ON)
set(CMAKE_CXX_EXTENSIONS OFF)

if (${LINK_GCOV})
  set(GCOV_LIB gcov)
else ()
  set(GCOV_LIB "")
endif ()

include($ENV{MAGMA_ROOT}/orc8r/gateway/c/common/CMakeProtoMacros.txt)

# compile the needed protos
set(COMMON_CPP_PROTOS  common_types)
set(REDIS_CPP_PROTOS redis)
list(APPEND PROTO_SRCS "")
list(APPEND PROTO_HDRS "")

set(LTE_OUT_DIR "${CMAKE_CURRENT_BINARY_DIR}/lte/protos/oai")
generate_cpp_protos("${COMMON_CPP_PROTOS}" "${PROTO_SRCS}"
  "${PROTO_HDRS}" ${STATE_PROTO_DIR} ${LTE_OUT_DIR})

create_proto_dir("orc8r" ORC8R_CPP_OUT_DIR)
generate_cpp_protos("${REDIS_CPP_PROTOS}" "${PROTO_SRCS}"
    "${PROTO_HDRS}" ${ORC8R_PROTO_DIR} ${ORC8R_CPP_OUT_DIR})

message("Proto_srcs are ${PROTO_SRCS}")
include_directories(${CMAKE_CURRENT_BINARY_DIR})
include_directories(${LTE_OUT_DIR})
include_directories(${ORC8R_CPP_OUT_DIR})


###############################################################################
# Core Network Utils
###############################################################################

set(S1AP_C_DIR ${PROJECT_BINARY_DIR}/s1ap/r10.5)
include_directories("${S1AP_C_DIR}")

set(COMMON_SRC
  TLVDecoder.c
  TLVEncoder.c
  async_system.c
  backtrace.c
  common_types.c
  conversions.c
  digest.c
  dynamic_memory_check.c
  enum_string.c
  itti_free_defined_msg.c
  mcc_mnc_itu.c
  pid_file.c
  shared_ts_log.c
  log.c
  state_converter.cpp
  ${PROTO_SRCS}
  ${PROTO_HDRS}
)

add_subdirectory(glogwrapper)
add_subdirectory(redis_utils)

add_library(COMMON ${COMMON_SRC})
target_link_libraries (COMMON
  lfds710
  LIB_3GPP LIB_BSTR LIB_HASHTABLE LIB_ITTI glogwrapper
)
target_include_directories(COMMON PUBLIC
  ${CMAKE_CURRENT_SOURCE_DIR}
)

# Copyright 2020 The Magma Authors.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

cmake_minimum_required(VERSION 3.0.2)

PROJECT(LIAgent)

include($ENV{MAGMA_ROOT}/orc8r/gateway/c/common/CMakeProtoMacros.txt)

if (NOT BUILD_TESTS)
  # Add AddressSanitizer (asan) support for debug builds of LiAgentD
  set (CMAKE_CXX_FLAGS_DEBUG
      "${CMAKE_CXX_FLAGS_DEBUG} -fno-omit-frame-pointer -fsanitize=address")
  set (CMAKE_LINKER_FLAGS_DEBUG
      "${CMAKE_LINKER_FLAGS_DEBUG} -fno-omit-frame-pointer -fsanitize=address")
  # Add LeakSanitizer (lsan) support to the release build of LIAgentD so that
  # we can find memory leaks in production.
  set(CMAKE_C_FLAGS_RELWITHDEBINFO
     "${CMAKE_C_FLAGS_RELWITHDEBINFO} -fsanitize=leak -fno-omit-frame-pointer")
  set(CMAKE_EXE_LINKER_FLAGS_RELWITHDEBINFO
     "${CMAKE_EXE_LINKER_FLAGS_RELWITHDEBINFO} -fsanitize=leak")
endif ()

set(OUTPUT_DIR "${CMAKE_CURRENT_BINARY_DIR}")

set(MAGMA_LIB_DIR $ENV{C_BUILD}/magma_common)
add_definitions(-DLOG_WITH_GLOG)

include_directories("${OUTPUT_DIR}")
include_directories("${MAGMA_ROOT}/orc8r/gateway/c/common/logging")

include_directories(${MAGMA_LIB_DIR}/service303)
include_directories(${MAGMA_LIB_DIR}/service_registry)
include_directories(${MAGMA_LIB_DIR}/async_grpc)
include_directories(${MAGMA_LIB_DIR}/config)

include($ENV{MAGMA_ROOT}/orc8r/gateway/c/common/CMakeProtoMacros.txt)

# compile the needed macros
create_proto_dir("lte" LTE_CPP_OUT_DIR)
create_proto_dir("orc8r" ORC8R_CPP_OUT_DIR)
create_proto_dir("feg/gateway/services/aaa" CWF_CPP_OUT_DIR)

list(APPEND PROTO_SRCS "")
list(APPEND PROTO_HDRS "")

message("Proto_srcs are ${PROTO_SRCS}")

find_library(SERVICE303_LIB SERVICE303_LIB HINTS ${MAGMA_LIB_DIR}/service303)

# compile the needed macros
create_proto_dir("lte" LTE_CPP_OUT_DIR)
create_proto_dir("orc8r" ORC8R_CPP_OUT_DIR)
create_proto_dir("feg/gateway/services/aaa" CWF_CPP_OUT_DIR)

list(APPEND PROTO_SRCS "")
list(APPEND PROTO_HDRS "")

set(SMGR_ORC8R_CPP_PROTOS directoryd)
generate_cpp_protos("${SMGR_ORC8R_CPP_PROTOS}" "${PROTO_SRCS}"
  "${PROTO_HDRS}" ${ORC8R_PROTO_DIR} ${ORC8R_CPP_OUT_DIR})

set(SMGR_LTE_CPP_PROTOS mconfig/mconfigs)
generate_cpp_protos("${SMGR_LTE_CPP_PROTOS}" "${PROTO_SRCS}"
  "${PROTO_HDRS}" ${LTE_PROTO_DIR} ${LTE_CPP_OUT_DIR})

set(SMGR_CWF_CPP_PROTOS accounting context)
generate_cwf_cpp_protos("${SMGR_CWF_CPP_PROTOS}" "${PROTO_SRCS}"
  "${PROTO_HDRS}" ${CWF_PROTO_DIR} ${CWF_CPP_OUT_DIR})

set(SMGR_ORC8R_GRPC_PROTOS directoryd)
generate_grpc_protos("${SMGR_ORC8R_GRPC_PROTOS}" "${PROTO_SRCS}"
  "${PROTO_HDRS}" ${ORC8R_PROTO_DIR} ${ORC8R_CPP_OUT_DIR})

message("Proto_srcs are ${PROTO_SRCS}")

link_directories(
  ${MAGMA_LIB_DIR}/async_grpc
  ${MAGMA_LIB_DIR}/config
  ${MAGMA_LIB_DIR}/service303
  ${MAGMA_LIB_DIR}/service_registry)

add_library(LI_AGENT
    PDUGenerator.cpp
    PDUGenerator.h
    InterfaceMonitor.cpp
    InterfaceMonitor.h
    DirectorydClient.cpp
    DirectorydClient.h
    ProxyConnector.cpp
    ProxyConnector.h
    ${PROTO_SRCS}
    ${PROTO_HDRS})

target_link_libraries(LI_AGENT
  SERVICE303_LIB SERVICE_REGISTRY CONFIG ASYNC_GRPC
  glog gflags folly pthread ${GCOV_LIB} grpc++ grpc yaml-cpp protobuf cpp_redis
  prometheus-cpp tacopie mnl tins pcap ssl crypto
  )

add_executable(liagentd ${PROJECT_SOURCE_DIR}/main.cpp)

target_link_libraries (liagentd LI_AGENT)

if (BUILD_TESTS)
  ENABLE_TESTING()
  ADD_SUBDIRECTORY(test)
endif (BUILD_TESTS)

# Copyright 2020 The Magma Authors.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

cmake_minimum_required(VERSION 3.0.2)

PROJECT(ConnectionTracker)

add_compile_options(-std=c++14)

include($ENV{MAGMA_ROOT}/orc8r/gateway/c/common/CMakeProtoMacros.txt)

if (NOT BUILD_TESTS)
  # Add AddressSanitizer (asan) support for debug builds of ConnectionD
  set (CMAKE_CXX_FLAGS_DEBUG
      "${CMAKE_CXX_FLAGS_DEBUG} -fno-omit-frame-pointer -fsanitize=address")
  set (CMAKE_LINKER_FLAGS_DEBUG
      "${CMAKE_LINKER_FLAGS_DEBUG} -fno-omit-frame-pointer -fsanitize=address")
  # Add LeakSanitizer (lsan) support to the release build of SessionD so that
  # we can find memory leaks in production.
  set(CMAKE_C_FLAGS_RELWITHDEBINFO
     "${CMAKE_C_FLAGS_RELWITHDEBINFO} -fsanitize=leak -fno-omit-frame-pointer")
  set(CMAKE_EXE_LINKER_FLAGS_RELWITHDEBINFO
     "${CMAKE_EXE_LINKER_FLAGS_RELWITHDEBINFO} -fsanitize=leak")
else() # BUILD_TESTS
  message("Adding code coverage build and linker flags for BUILD_TESTS")
  set (CMAKE_CXX_FLAGS_DEBUG
      "${CMAKE_CXX_FLAGS_DEBUG} -g -O0 -fprofile-arcs -ftest-coverage")
  set(CMAKE_LINKER_FLAGS_DEBUG
      "${CMAKE_LINKER_FLAGS_DEBUG} -fprofile-arcs -ftest-coverage")
endif ()

set(OUTPUT_DIR "${CMAKE_CURRENT_BINARY_DIR}")

set(MAGMA_LIB_DIR $ENV{C_BUILD}/magma_common)
add_definitions(-DLOG_WITH_GLOG)

message("Build type is ${CMAKE_BUILD_TYPE}")
if (CMAKE_BUILD_TYPE STREQUAL "Debug")
  set(CMAKE_CXX_FLAGS "${CMAKE_CXX_FLAGS} -Wall")
endif ()

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


set(SMGR_LTE_CPP_PROTOS session_manager subscriberdb
  pipelined spgw_service mconfig/mconfigs)
generate_cpp_protos("${SMGR_LTE_CPP_PROTOS}" "${PROTO_SRCS}"
  "${PROTO_HDRS}" ${LTE_PROTO_DIR} ${LTE_CPP_OUT_DIR})

set(SMGR_CWF_CPP_PROTOS accounting context)
generate_cwf_cpp_protos("${SMGR_CWF_CPP_PROTOS}" "${PROTO_SRCS}"
  "${PROTO_HDRS}" ${CWF_PROTO_DIR} ${CWF_CPP_OUT_DIR})

link_directories(
  ${MAGMA_LIB_DIR}/config
  ${MAGMA_LIB_DIR}/service303
  ${MAGMA_LIB_DIR}/service_registry)

add_library(CONNECTION_TRACKER
    PacketGenerator.cpp
    PacketGenerator.h
    EventTracker.cpp
    EventTracker.h
    ${PROTO_SRCS}
    ${PROTO_HDRS})

target_link_libraries(CONNECTION_TRACKER
  SERVICE303_LIB SERVICE_REGISTRY CONFIG
  glog gflags folly pthread ${GCOV_LIB} grpc++ grpc yaml-cpp protobuf cpp_redis
  prometheus-cpp tacopie mnl tins
  )

if (CLANG_FORMAT)
  add_custom_target(
          clang-format
          COMMAND ${CLANG_FORMAT}
          -style=file
          -i
          ${ALL_CXX_SOURCE_FILES}
  )
endif ()

add_executable(connectiond ${PROJECT_SOURCE_DIR}/main.cpp)

target_link_libraries (connectiond CONNECTION_TRACKER)

if (BUILD_TESTS)
  ENABLE_TESTING()
  ADD_SUBDIRECTORY(test)
endif (BUILD_TESTS)

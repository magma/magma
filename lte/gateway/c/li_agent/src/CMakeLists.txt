# Copyright 2020 The Magma Authors.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

cmake_minimum_required(VERSION 3.7.2)

PROJECT(LIAgent)

find_package(SERVICE303_LIB REQUIRED)
find_package(SERVICE_REGISTRY REQUIRED)
find_package(MAGMA_CONFIG REQUIRED)
find_package(ASYNC_GRPC REQUIRED)
find_package(MAGMA_LOGGING REQUIRED)

include($ENV{MAGMA_ROOT}/orc8r/gateway/c/common/CMakeProtoMacros.txt)

if (BUILD_TESTS)
  message("Adding code coverage build and linker flags for BUILD_TESTS")
  set(CMAKE_CXX_FLAGS_DEBUG
      "${CMAKE_CXX_FLAGS_DEBUG} -g -O0 -fprofile-arcs -ftest-coverage")
  set(CMAKE_LINKER_FLAGS_DEBUG
      "${CMAKE_LINKER_FLAGS_DEBUG} -fprofile-arcs -ftest-coverage")
endif ()

# Add AddressSanitizer (asan) support for debug builds of LiAgentD
set(CMAKE_CXX_FLAGS_DEBUG
    "${CMAKE_CXX_FLAGS_DEBUG} -fno-omit-frame-pointer -fsanitize=address")
set(CMAKE_LINKER_FLAGS_DEBUG
    "${CMAKE_LINKER_FLAGS_DEBUG} -fno-omit-frame-pointer -fsanitize=address")
# Add LeakSanitizer (lsan) support to the release build of LIAgentD so that
# we can find memory leaks in production.
set(CMAKE_C_FLAGS_RELWITHDEBINFO
    "${CMAKE_C_FLAGS_RELWITHDEBINFO} -fsanitize=leak -fno-omit-frame-pointer")
set(CMAKE_EXE_LINKER_FLAGS_RELWITHDEBINFO
    "${CMAKE_EXE_LINKER_FLAGS_RELWITHDEBINFO} -fsanitize=leak")

add_definitions(-DLOG_WITH_GLOG)

# compile the needed macros
create_proto_dir("lte" LTE_CPP_OUT_DIR)
list(APPEND PROTO_SRCS "")
list(APPEND PROTO_HDRS "")
message("Proto_srcs are ${PROTO_SRCS}")
set(LI_AGENT_LTE_CPP_PROTOS mconfig/mconfigs mobilityd
    mobilityd subscriberdb apn)
set(RPC_LTE_GRPC_PROTOS mobilityd)
generate_cpp_protos("${LI_AGENT_LTE_CPP_PROTOS}" "${PROTO_SRCS}"
    "${PROTO_HDRS}" ${LTE_PROTO_DIR} ${LTE_CPP_OUT_DIR})
generate_grpc_protos("${RPC_LTE_GRPC_PROTOS}" "${PROTO_SRCS}" "${PROTO_HDRS}"
    ${LTE_PROTO_DIR} ${LTE_CPP_OUT_DIR})

add_library(LI_AGENT
    Utilities.h
    Utilities.cpp
    MobilitydClient.cpp
    MobilitydClient.h
    PDUGenerator.cpp
    PDUGenerator.h
    InterfaceMonitor.cpp
    InterfaceMonitor.h
    ProxyConnector.cpp
    ProxyConnector.h
    ${PROTO_SRCS}
    ${PROTO_HDRS})

target_link_libraries(LI_AGENT
    SERVICE303_LIB SERVICE_REGISTRY MAGMA_CONFIG ASYNC_GRPC MAGMA_LOGGING
    glog gflags folly pthread ${GCOV_LIB} grpc++ grpc yaml-cpp protobuf cpp_redis
    prometheus-cpp tacopie mnl tins pcap ssl crypto uuid
    )

target_include_directories(LI_AGENT PUBLIC
    $<TARGET_FILE_DIR:LI_AGENT>)

add_executable(liagentd main.cpp)

target_link_libraries(liagentd LI_AGENT)

if (BUILD_TESTS)
  ENABLE_TESTING()
  ADD_SUBDIRECTORY(test)
endif (BUILD_TESTS)

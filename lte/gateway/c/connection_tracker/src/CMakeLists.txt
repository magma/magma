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
  set(CMAKE_CXX_FLAGS_DEBUG
      "${CMAKE_CXX_FLAGS_DEBUG} -fno-omit-frame-pointer -fsanitize=address")
  set(CMAKE_LINKER_FLAGS_DEBUG
      "${CMAKE_LINKER_FLAGS_DEBUG} -fno-omit-frame-pointer -fsanitize=address")
  # Add LeakSanitizer (lsan) support to the release build of SessionD so that
  # we can find memory leaks in production.
  set(CMAKE_C_FLAGS_RELWITHDEBINFO
      "${CMAKE_C_FLAGS_RELWITHDEBINFO} -fsanitize=leak -fno-omit-frame-pointer")
  set(CMAKE_EXE_LINKER_FLAGS_RELWITHDEBINFO
      "${CMAKE_EXE_LINKER_FLAGS_RELWITHDEBINFO} -fsanitize=leak")
else () # BUILD_TESTS
  message("Adding code coverage build and linker flags for BUILD_TESTS")
  set(CMAKE_CXX_FLAGS_DEBUG
      "${CMAKE_CXX_FLAGS_DEBUG} -g -O0 -fprofile-arcs -ftest-coverage")
  set(CMAKE_LINKER_FLAGS_DEBUG
      "${CMAKE_LINKER_FLAGS_DEBUG} -fprofile-arcs -ftest-coverage")
endif ()

add_definitions(-DLOG_WITH_GLOG)

message("Build type is ${CMAKE_BUILD_TYPE}")
if (CMAKE_BUILD_TYPE STREQUAL "Debug")
  set(CMAKE_CXX_FLAGS "${CMAKE_CXX_FLAGS} -Wall")
endif ()

find_package(MAGMA_LOGGING REQUIRED)
find_package(SERVICE303_LIB REQUIRED)
find_package(MAGMA_CONFIG REQUIRED)
find_package(SERVICE_REGISTRY REQUIRED)

# compile the needed protos
set(CT_ORC8R_CPP_PROTOS "")
set(CT_LTE_CPP_PROTOS mconfig/mconfigs)
set(CT_ORC8R_GRPC_PROTOS "")
set(CT_LTE_GRPC_PROTOS "")

generate_all_protos("${CT_LTE_CPP_PROTOS}" "${CT_ORC8R_CPP_PROTOS}"
    "${CT_LTE_GRPC_PROTOS}" "${CT_ORC8R_GRPC_PROTOS}" "${PROTO_SRCS}" "${PROTO_HDRS}")
message("Proto_srcs are ${PROTO_SRCS}")
message("Proto_hdr are ${PROTO_HDRS}")

add_library(CONNECTION_TRACKER
    PacketGenerator.cpp
    PacketGenerator.h
    EventTracker.cpp
    EventTracker.h
    ${PROTO_SRCS}
    ${PROTO_HDRS})

target_link_libraries(CONNECTION_TRACKER
    SERVICE303_LIB SERVICE_REGISTRY MAGMA_LOGGING MAGMA_CONFIG
    glog gflags folly pthread ${GCOV_LIB} grpc++ grpc yaml-cpp protobuf cpp_redis
    prometheus-cpp tacopie mnl tins
    )

target_include_directories(CONNECTION_TRACKER PUBLIC
    "$<BUILD_INTERFACE:${CMAKE_CURRENT_SOURCE_DIR}>"
    $<TARGET_FILE_DIR:CONNECTION_TRACKER>)

if (CLANG_FORMAT)
  add_custom_target(
      clang-format
      COMMAND ${CLANG_FORMAT}
      -style=file
      -i
      ${ALL_CXX_SOURCE_FILES}
  )
endif ()

add_executable(connectiond
    main.cpp)

target_link_libraries(connectiond CONNECTION_TRACKER)

if (BUILD_TESTS)
  ENABLE_TESTING()
  ADD_SUBDIRECTORY(test)
endif (BUILD_TESTS)
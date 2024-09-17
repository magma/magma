# Copyright 2020 The Magma Authors.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

cmake_minimum_required(VERSION 3.0.2)

PROJECT(SessionManager)

add_compile_options(-std=c++14)

set(CMAKE_MODULE_PATH "${PROJECT_SOURCE_DIR}/modules")
include(FindClangFormat)

include($ENV{MAGMA_ROOT}/orc8r/gateway/c/common/CMakeProtoMacros.txt)

# Add AddressSanitizer (asan) support for debug builds of SessionD. Disable adding
# DISABLE_ASAN flag
if (DISABLE_ASAN)
  message("-- ASAN DISABLED")
else()
  message("++ ASAN ENABLED")
  set (CMAKE_CXX_FLAGS_DEBUG
          "${CMAKE_CXX_FLAGS_DEBUG} -fno-omit-frame-pointer -fsanitize=address")
  set (CMAKE_LINKER_FLAGS_DEBUG
          "${CMAKE_LINKER_FLAGS_DEBUG} -fno-omit-frame-pointer -fsanitize=address")
endif ()

if (NOT BUILD_TESTS)
  # Add LeakSanitizer (lsan) support to the release build of SessionD so that
  # we can find memory leaks in production.
  # Disabling LSAN on Prod until we figure out how to get gcore to work
  #set(CMAKE_C_FLAGS_RELWITHDEBINFO
  #   "${CMAKE_C_FLAGS_RELWITHDEBINFO} -fsanitize=leak -fno-omit-frame-pointer")
  #set(CMAKE_EXE_LINKER_FLAGS_RELWITHDEBINFO
  #   "${CMAKE_EXE_LINKER_FLAGS_RELWITHDEBINFO} -fsanitize=leak")
else() # BUILD_TESTS
endif ()

set(OUTPUT_DIR "${CMAKE_CURRENT_BINARY_DIR}")

set(MAGMA_LIB_DIR $ENV{C_BUILD}/magma_common)
add_definitions(-DLOG_WITH_GLOG)

message("Build type is ${CMAKE_BUILD_TYPE}")
if (CMAKE_BUILD_TYPE STREQUAL "Debug")
  set(CMAKE_CXX_FLAGS "${CMAKE_CXX_FLAGS} -Wall -Wno-unused-function")
endif ()

include_directories("${MAGMA_ROOT}/orc8r/gateway/c/common/logging")
include_directories(${MAGMA_ROOT}/orc8r/gateway/c/common/service303)
include_directories(${MAGMA_ROOT}/orc8r/gateway/c/common/service_registry)
include_directories(${MAGMA_ROOT}/orc8r/gateway/c/common/async_grpc)
include_directories(${MAGMA_ROOT}/orc8r/gateway/c/common/async_grpc/includes)
include_directories(${MAGMA_ROOT}/orc8r/gateway/c/common/config)
include_directories(${MAGMA_ROOT}/orc8r/gateway/c/common/eventd)
include_directories("${OUTPUT_DIR}")

# TODO: Temp workaround until packages are imported by these cmakefile
# Will be removed in subsequent patch
set(MAGMA_COMMON_INCL_DIR $ENV{C_BUILD}/common)
include_directories(${MAGMA_COMMON_INCL_DIR})

include($ENV{MAGMA_ROOT}/orc8r/gateway/c/common/CMakeProtoMacros.txt)

# compile the needed macros
create_proto_dir("lte" LTE_CPP_OUT_DIR)
create_proto_dir("orc8r" ORC8R_CPP_OUT_DIR)
create_proto_dir("feg/gateway/services/aaa" CWF_CPP_OUT_DIR)

list(APPEND PROTO_SRCS "")
list(APPEND PROTO_HDRS "")

set(SMGR_ORC8R_CPP_PROTOS common digest directoryd eventd metricsd redis service303 mconfig/mconfigs)
generate_cpp_protos("${SMGR_ORC8R_CPP_PROTOS}" "${PROTO_SRCS}"
  "${PROTO_HDRS}" ${ORC8R_PROTO_DIR} ${ORC8R_CPP_OUT_DIR})

set(SMGR_LTE_CPP_PROTOS apn policydb mobilityd subscriberdb session_manager
  pipelined spgw_service abort_session mconfig/mconfigs)
generate_cpp_protos("${SMGR_LTE_CPP_PROTOS}" "${PROTO_SRCS}"
  "${PROTO_HDRS}" ${LTE_PROTO_DIR} ${LTE_CPP_OUT_DIR})

set(SMGR_CWF_CPP_PROTOS accounting context)
generate_cwf_cpp_protos("${SMGR_CWF_CPP_PROTOS}" "${PROTO_SRCS}"
  "${PROTO_HDRS}" ${CWF_PROTO_DIR} ${CWF_CPP_OUT_DIR})

set(SMGR_ORC8R_GRPC_PROTOS directoryd eventd service303)
generate_grpc_protos("${SMGR_ORC8R_GRPC_PROTOS}" "${PROTO_SRCS}"
  "${PROTO_HDRS}" ${ORC8R_PROTO_DIR} ${ORC8R_CPP_OUT_DIR})

set(SMGR_GRPC_PROTOS abort_session session_manager pipelined spgw_service mobilityd)
generate_grpc_protos("${SMGR_GRPC_PROTOS}" "${PROTO_SRCS}"
  "${PROTO_HDRS}" ${LTE_PROTO_DIR} ${LTE_CPP_OUT_DIR})

set(SMGR_CWF_GRPC_PROTOS accounting)
generate_grpc_protos("${SMGR_CWF_GRPC_PROTOS}" "${PROTO_SRCS}"
  "${PROTO_HDRS}" ${CWF_PROTO_DIR} ${CWF_CPP_OUT_DIR})

message("SessionManager proto_srcs are ${PROTO_SRCS}")

link_directories(
  ${MAGMA_LIB_DIR}/async_grpc
  ${MAGMA_LIB_DIR}/config
  ${MAGMA_LIB_DIR}/service303
  ${MAGMA_LIB_DIR}/eventd
  ${MAGMA_LIB_DIR}/service_registry)

add_library(SESSION_MANAGER
    AAAClient.cpp
    AAAClient.hpp
    ChargingGrant.cpp
    ChargingGrant.hpp
    EnumToString.cpp
    EnumToString.hpp
    SessionManagerServer.cpp
    SessionManagerServer.hpp
    LocalSessionManagerHandler.cpp
    LocalSessionManagerHandler.hpp
    LocalEnforcer.cpp
    LocalEnforcer.hpp
    SessionState.cpp
    SessionState.hpp
    SessionCredit.cpp
    SessionCredit.hpp
    RestartHandler.cpp
    RestartHandler.hpp
    RuleStore.cpp
    RuleStore.hpp
    SessionReporter.cpp
    SessionReporter.hpp
    SessionID.cpp
    SessionID.hpp
    ServiceAction.hpp
    SpgwServiceClient.cpp
    SpgwServiceClient.hpp
    PipelinedClient.cpp
    PipelinedClient.hpp
    DirectorydClient.cpp
    DirectorydClient.hpp
    SessionEvents.cpp
    SessionEvents.hpp
    DiameterCodes.cpp
    DiameterCodes.hpp
    CreditKey.hpp
    StatsPoller.cpp
    StatsPoller.hpp
    ShardTracker.cpp
    ShardTracker.hpp
    SessionProxyResponderHandler.cpp
    SessionProxyResponderHandler.hpp
    StoredState.cpp
    StoredState.hpp
    SessionStore.cpp
    SessionStore.hpp
    MemoryStoreClient.cpp
    MemoryStoreClient.hpp
    Monitor.hpp
    RedisStoreClient.cpp
    RedisStoreClient.hpp
    StoreClient.hpp
    MeteringReporter.cpp
    MeteringReporter.hpp
    GrpcMagmaUtils.cpp
    GrpcMagmaUtils.hpp
    SetMessageManagerHandler.hpp
    SetMessageManagerHandler.cpp
    UpfMsgManageHandler.hpp
    UpfMsgManageHandler.cpp
    SessionStateEnforcer.hpp
    SessionStateEnforcer.cpp
    AmfServiceClient.hpp
    AmfServiceClient.cpp
    Utilities.cpp
    Utilities.hpp
    OperationalStatesHandler.cpp
    OperationalStatesHandler.hpp
    Types.hpp
    MobilitydClient.cpp
    MobilitydClient.hpp
    ObjectMap.hpp
    PolicyLoader.cpp
    PolicyLoader.hpp
    ObjectMap.hpp
    RedisMap.hpp
    Serializers.cpp
    Serializers.hpp
    ${PROTO_SRCS}
    ${PROTO_HDRS})

find_package(MAGMA_LOGGING REQUIRED)
find_package(MAGMA_SENTRY REQUIRED)

target_link_libraries(SESSION_MANAGER
  SERVICE303_LIB SERVICE_REGISTRY ASYNC_GRPC MAGMA_CONFIG MAGMA_LOGGING EVENTD MAGMA_SENTRY
  glog gflags folly pthread ${GCOV_LIB} gpr grpc++ grpc yaml-cpp protobuf cpp_redis
  prometheus-cpp tacopie systemd
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

add_executable(sessiond ${PROJECT_SOURCE_DIR}/sessiond_main.cpp)
target_link_libraries(sessiond SESSION_MANAGER)


if (BUILD_TESTS)
  ENABLE_TESTING()
  ADD_SUBDIRECTORY(test)
endif (BUILD_TESTS)

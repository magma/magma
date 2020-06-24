# Copyright (c) 2016-present, Facebook, Inc.
# All rights reserved.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree. An additional grant
# of patent rights can be found in the PATENTS file in the same directory.

cmake_minimum_required(VERSION 3.0.2)

PROJECT(SessionManager)

add_compile_options(-std=c++14)

include($ENV{MAGMA_ROOT}/orc8r/gateway/c/common/CMakeProtoMacros.txt)

if (NOT BUILD_TESTS)
  set (CMAKE_CXX_FLAGS_DEBUG "${CMAKE_CXX_FLAGS_DEBUG} -fno-omit-frame-pointer -fsanitize=address")
  set (CMAKE_LINKER_FLAGS_DEBUG "${CMAKE_LINKER_FLAGS_DEBUG} -fno-omit-frame-pointer -fsanitize=address")
endif ()

set(OUTPUT_DIR "${CMAKE_CURRENT_BINARY_DIR}")

set(MAGMA_LIB_DIR $ENV{C_BUILD}/magma_common)
add_definitions(-DLOG_WITH_GLOG)

message("Build type is ${CMAKE_BUILD_TYPE}")
if (CMAKE_BUILD_TYPE STREQUAL "Debug")
  set(CMAKE_CXX_FLAGS
    "${CMAKE_CXX_FLAGS} -Wall -Werror -Wno-unused-function"
  )
endif ()

include_directories("${OUTPUT_DIR}")
include_directories("${MAGMA_ROOT}/orc8r/gateway/c/common/logging")

include_directories(${MAGMA_LIB_DIR}/service303)
include_directories(${MAGMA_LIB_DIR}/service_registry)
include_directories(${MAGMA_LIB_DIR}/async_grpc)
include_directories(${MAGMA_LIB_DIR}/config)
include_directories(${MAGMA_LIB_DIR}/datastore)
include_directories(${MAGMA_LIB_DIR}/policydb)
include_directories(${MAGMA_LIB_DIR}/eventd)

include($ENV{MAGMA_ROOT}/orc8r/gateway/c/common/CMakeProtoMacros.txt)

# compile the needed macros
create_proto_dir("lte" LTE_CPP_OUT_DIR)
create_proto_dir("orc8r" ORC8R_CPP_OUT_DIR)
create_proto_dir("feg/gateway/services/aaa" CWF_CPP_OUT_DIR)

list(APPEND PROTO_SRCS "")
list(APPEND PROTO_HDRS "")

set(SMGR_ORC8R_CPP_PROTOS common directoryd)
generate_cpp_protos("${SMGR_ORC8R_CPP_PROTOS}" "${PROTO_SRCS}"
  "${PROTO_HDRS}" ${ORC8R_PROTO_DIR} ${ORC8R_CPP_OUT_DIR})

set(SMGR_LTE_CPP_PROTOS session_manager subscriberdb policydb
  pipelined spgw_service mconfig/mconfigs)
generate_cpp_protos("${SMGR_LTE_CPP_PROTOS}" "${PROTO_SRCS}"
  "${PROTO_HDRS}" ${LTE_PROTO_DIR} ${LTE_CPP_OUT_DIR})

set(SMGR_CWF_CPP_PROTOS accounting context)
generate_cwf_cpp_protos("${SMGR_CWF_CPP_PROTOS}" "${PROTO_SRCS}"
  "${PROTO_HDRS}" ${CWF_PROTO_DIR} ${CWF_CPP_OUT_DIR})

set(SMGR_ORC8R_GRPC_PROTOS directoryd)
generate_grpc_protos("${SMGR_ORC8R_GRPC_PROTOS}" "${PROTO_SRCS}"
  "${PROTO_HDRS}" ${ORC8R_PROTO_DIR} ${ORC8R_CPP_OUT_DIR})

set(SMGR_GRPC_PROTOS session_manager pipelined spgw_service)
generate_grpc_protos("${SMGR_GRPC_PROTOS}" "${PROTO_SRCS}"
  "${PROTO_HDRS}" ${LTE_PROTO_DIR} ${LTE_CPP_OUT_DIR})

set(SMGR_CWF_GRPC_PROTOS accounting)
generate_grpc_protos("${SMGR_CWF_GRPC_PROTOS}" "${PROTO_SRCS}"
  "${PROTO_HDRS}" ${CWF_PROTO_DIR} ${CWF_CPP_OUT_DIR})

message("Proto_srcs are ${PROTO_SRCS}")

link_directories(
  ${MAGMA_LIB_DIR}/async_grpc
  ${MAGMA_LIB_DIR}/config
  ${MAGMA_LIB_DIR}/service303
  ${MAGMA_LIB_DIR}/datastore
  ${MAGMA_LIB_DIR}/policydb
  ${MAGMA_LIB_DIR}/eventd
  ${MAGMA_LIB_DIR}/service_registry)

add_library(SESSION_MANAGER
    AAAClient.cpp
    AAAClient.h
    EnumToString.cpp
    EnumToString.h
    SessionManagerServer.cpp
    SessionManagerServer.h
    LocalSessionManagerHandler.cpp
    LocalSessionManagerHandler.h
    LocalEnforcer.cpp
    LocalEnforcer.h
    SessionState.cpp
    SessionState.h
    SessionCredit.cpp
    SessionCredit.h
    RestartHandler.cpp
    RestartHandler.h
    RuleStore.cpp
    RuleStore.h
    SessionReporter.cpp
    SessionReporter.h
    SessionID.cpp
    SessionID.h
    ServiceAction.h
    SpgwServiceClient.cpp
    SpgwServiceClient.h
    PipelinedClient.cpp
    PipelinedClient.h
    DirectorydClient.cpp
    DirectorydClient.h
    SessionEvents.cpp
    SessionEvents.h
    DiameterCodes.cpp
    DiameterCodes.h
    CreditKey.h
    SessionProxyResponderHandler.cpp
    SessionProxyResponderHandler.h
    StoredState.cpp
    StoredState.h
    SessionStore.cpp
    SessionStore.h
    MemoryStoreClient.cpp
    MemoryStoreClient.h
    Monitor.h
    RedisStoreClient.cpp
    RedisStoreClient.h
    StoreClient.h
    MeteringReporter.cpp
    MeteringReporter.h
    ${PROTO_SRCS}
    ${PROTO_HDRS}
    )

target_link_libraries(SESSION_MANAGER
  SERVICE303_LIB SERVICE_REGISTRY ASYNC_GRPC CONFIG POLICYDB EVENTD
  DATASTORE
  glog gflags folly pthread ${GCOV_LIB} grpc++ grpc yaml-cpp protobuf cpp_redis
  prometheus-cpp tacopie
  )


add_executable(sessiond ${PROJECT_SOURCE_DIR}/sessiond_main.cpp)

target_link_libraries (sessiond SESSION_MANAGER)

if (BUILD_TESTS)
  ENABLE_TESTING()
  ADD_SUBDIRECTORY(test)
endif (BUILD_TESTS)

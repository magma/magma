# Copyright (c) 2016-present, Facebook, Inc.
# All rights reserved.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree. An additional grant
# of patent rights can be found in the PATENTS file in the same directory.

set (CMAKE_CXX_FLAGS "-Wno-write-strings -Wno-literal-suffix")

# compile the needed protos
set(S6ASRV_FEG_CPP_PROTOS s6a_proxy)
set(S6ASRV_FEG_GRPC_PROTOS s6a_proxy)
set(S6ASRV_LTE_CPP_PROTOS s6a_service)
set(S6ASRV_LTE_GRPC_PROTOS s6a_service)

list(APPEND PROTO_SRCS "")
list(APPEND PROTO_HDRS "")

create_proto_dir("feg" FEG_OUT_DIR)
create_proto_dir("lte" LTE_OUT_DIR)

generate_cpp_protos("${S6ASRV_FEG_CPP_PROTOS}" "${PROTO_SRCS}"
  "${PROTO_HDRS}" ${FEG_PROTO_DIR} ${FEG_OUT_DIR})
generate_grpc_protos("${S6ASRV_FEG_GRPC_PROTOS}" "${PROTO_SRCS}"
  "${PROTO_HDRS}" ${FEG_PROTO_DIR} ${FEG_OUT_DIR})

generate_cpp_protos("${S6ASRV_LTE_CPP_PROTOS}" "${PROTO_SRCS}"
  "${PROTO_HDRS}" ${LTE_PROTO_DIR} ${LTE_OUT_DIR})
generate_grpc_protos("${S6ASRV_LTE_GRPC_PROTOS}" "${PROTO_SRCS}"
  "${PROTO_HDRS}" ${LTE_PROTO_DIR} ${LTE_OUT_DIR})

include_directories(${CMAKE_CURRENT_BINARY_DIR})

add_library(TASK_S6A_SERVICE
    S6aProxyImpl.cpp
    S6aServiceImpl.cpp
    s6a_service.cpp
    ${PROTO_SRCS}
    ${PROTO_HDRS}
    )

target_link_libraries(TASK_S6A_SERVICE
    COMMON
    LIB_BSTR protobuf
)
target_include_directories(TASK_S6A_SERVICE PUBLIC
    ${CMAKE_CURRENT_SOURCE_DIR}
    )

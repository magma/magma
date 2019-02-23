# Copyright (c) 2016-present, Facebook, Inc.
# All rights reserved.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree. An additional grant
# of patent rights can be found in the PATENTS file in the same directory.

find_package(Protobuf REQUIRED)
add_compile_options(-std=c++14)

# compile the needed protos
set(SGSSRV_FEG_CPP_PROTOS csfb)
set(SGSSRV_FEG_GRPC_PROTOS csfb)

list(APPEND PROTO_SRCS "")
list(APPEND PROTO_HDRS "")

create_proto_dir("feg" FEG_OUT_DIR)

generate_cpp_protos("${SGSSRV_FEG_CPP_PROTOS}" "${PROTO_SRCS}"
  "${PROTO_HDRS}" ${FEG_PROTO_DIR} ${FEG_OUT_DIR})
generate_grpc_protos("${SGSSRV_FEG_GRPC_PROTOS}" "${PROTO_SRCS}"
  "${PROTO_HDRS}" ${FEG_PROTO_DIR} ${FEG_OUT_DIR})

include_directories(${CMAKE_CURRENT_BINARY_DIR})

add_library(TASK_SGS_SERVICE
    CSFBGatewayServiceImpl.cpp
    proto_msg_to_itti_msg.cpp
    sgs_service.cpp
    "${PROTO_SRCS}"
    "${PROTO_HDRS}"
    )

target_link_libraries(TASK_SGS_SERVICE
    COMMON
    LIB_BSTR LIB_HASHTABLE
    ${PROTOBUF_LIBRARIES}
    grpc++
)
target_include_directories(TASK_SGS_SERVICE PUBLIC
    ${CMAKE_CURRENT_SOURCE_DIR}
    )

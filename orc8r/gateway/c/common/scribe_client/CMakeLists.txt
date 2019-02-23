# Copyright (c) 2016-present, Facebook, Inc.
# All rights reserved.
# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree. An additional grant
# of patent rights can be found in the PATENTS file in the same directory.

add_compile_options(-std=c++11)
include_directories(${CMAKE_CURRENT_BINARY_DIR})
#compile the relevant protos
create_proto_dir("orc8r" CPP_OUT_DIR)

list(APPEND PROTO_SRCS "")
list(APPEND PROTO_HDRS "")

set(SCRIBE_PROTOS common logging_service)
generate_cpp_protos("${SCRIBE_PROTOS}" "${PROTO_SRCS}" "${PROTO_HDRS}"
  ${ORC8R_PROTO_DIR} ${CPP_OUT_DIR})

set(SCRIBE_GRPC_PROTOS logging_service)
generate_grpc_protos("${SCRIBE_GRPC_PROTOS}" "${PROTO_SRCS}" "${PROTO_HDRS}"
  ${ORC8R_PROTO_DIR} ${CPP_OUT_DIR})

add_library(SCRIBE_CLIENT
    ScribeClient.cpp
    ScribeRpcClient.cpp
    ${PROTO_SRCS}
    ${PROTO_HDRS}
)

target_link_libraries(SCRIBE_CLIENT SERVICE_REGISTRY ASYNC_GRPC)

target_include_directories(SCRIBE_CLIENT PUBLIC
    ${CMAKE_CURRENT_SOURCE_DIR}
)

# Copyright (c) 2016-present, Facebook, Inc.
# All rights reserved.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree. An additional grant
# of patent rights can be found in the PATENTS file in the same directory.

include_directories("${PROJECT_SOURCE_DIR}/../common/logging")

#compile the relevant protos

list(APPEND PROTO_SRCS "")
list(APPEND PROTO_HDRS "")

set(ASYNC_ORC8R_CPP_PROTOS common eventd)
set(ASYNC_LTE_CPP_PROTOS session_manager policydb subscriberdb mobilityd)
set(ASYNC_LTE_GRPC_PROTOS session_manager mobilityd)
set(ASYNC_ORC8R_GRPC_PROTOS eventd)

generate_all_protos("${ASYNC_LTE_CPP_PROTOS}" "${ASYNC_ORC8R_CPP_PROTOS}"
  "${ASYNC_LTE_GRPC_PROTOS}"
  "${ASYNC_ORC8R_GRPC_PROTOS}" "${PROTO_SRCS}" "${PROTO_HDRS}")

message("Async Proto_srcs are ${PROTO_SRCS}")

add_library(ASYNC_GRPC
    GRPCReceiver.cpp
    ${PROTO_SRCS}
    ${PROTO_HDRS}
)

# copy headers to build directory so they can be shared with OAI,
# session_manager, etc.
add_custom_command(TARGET ASYNC_GRPC POST_BUILD
                   COMMAND ${CMAKE_COMMAND} -E copy
                   ${PROJECT_SOURCE_DIR}/async_grpc/*.h
                   $<TARGET_FILE_DIR:ASYNC_GRPC>)

target_include_directories(ASYNC_GRPC PUBLIC
                   $<TARGET_FILE_DIR:ASYNC_GRPC>
)

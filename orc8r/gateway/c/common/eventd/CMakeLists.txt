# Copyright (c) 2016-present, Facebook, Inc.
# All rights reserved.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree. An additional grant
# of patent rights can be found in the PATENTS file in the same directory.

include($ENV{MAGMA_ROOT}/orc8r/gateway/c/common/CMakeProtoMacros.txt)

add_compile_options(-std=c++11)

set(OUTPUT_DIR "${CMAKE_CURRENT_BINARY_DIR}")

include_directories("${OUTPUT_DIR}")
include_directories(${MAGMA_LIB_DIR}/async_grpc)
include_directories(${MAGMA_LIB_DIR}/service_registry)
include_directories("${PROJECT_SOURCE_DIR}/../common/eventd")


list(APPEND PROTO_SRCS "")
list(APPEND PROTO_HDRS "")

create_proto_dir("orc8r" CPP_OUT_DIR)
set(EVENTD_PROTOS common eventd)
generate_cpp_protos("${EVENTD_PROTOS}" "${PROTO_SRCS}" "${PROTO_HDRS}"
  ${ORC8R_PROTO_DIR} ${CPP_OUT_DIR})
generate_grpc_protos("${EVENTD_PROTOS}" "${PROTO_SRCS}" "${PROTO_HDRS}"
  ${ORC8R_PROTO_DIR} ${CPP_OUT_DIR})


add_library(EVENTD
        ${PROTO_SRCS}
        ${PROTO_HDRS}
        EventdClient.cpp
        EventdClient.h
        )

link_directories(
  ${MAGMA_LIB_DIR}/async_grpc
  ${MAGMA_LIB_DIR}/service_registry
  )


target_link_libraries(
  EVENTD
  SERVICE_REGISTRY ASYNC_GRPC
  grpc++ grpc
  )

# copy headers to build directory so they can be shared with OAI,
# session_manager, etc.
add_custom_command(TARGET EVENTD POST_BUILD
                   COMMAND ${CMAKE_COMMAND} -E copy
                   ${PROJECT_SOURCE_DIR}/eventd/*.h
                   $<TARGET_FILE_DIR:EVENTD>)

target_include_directories(EVENTD PUBLIC
                   $<TARGET_FILE_DIR:EVENTD>
)

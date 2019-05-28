# Copyright (c) 2017-present, Facebook, Inc.
# All rights reserved.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree. An additional grant
# of patent rights can be found in the PATENTS file in the same directory.

# compile the needed protos
set(DRD_ORC8R_CPP_PROTOS common directoryd)
set(DRD_LTE_CPP_PROTOS "")
set(DRD_ORC8R_GRPC_PROTOS directoryd)
set(DRD_LTE_GRPC_PROTOS "")

list(APPEND PROTO_SRCS "")
list(APPEND PROTO_HDRS "")

generate_all_protos("${DRD_LTE_CPP_PROTOS}" "${DRD_ORC8R_CPP_PROTOS}"
  "${DRD_LTE_GRPC_PROTOS}" "${DRD_ORC8R_GRPC_PROTOS}" "${PROTO_SRCS}" "${PROTO_HDRS}")

include_directories(${CMAKE_CURRENT_BINARY_DIR})

add_library(LIB_DIRECTORYD
  directoryd.cpp
  DirectorydClient.cpp
  ${PROTO_SRCS}
  ${PROTO_HDRS}
  )

target_link_libraries(LIB_DIRECTORYD
  LIB_RPC_CLIENT ${ASYNC_GRPC} ${SERVICE_REGISTRY}
  ${CONFIG}
)

target_include_directories(LIB_DIRECTORYD PUBLIC
  ${MAGMA_LIB_DIR}/async_grpc
  ${MAGMA_LIB_DIR}/config
  ${MAGMA_LIB_DIR}/service_registry
  ${CMAKE_CURRENT_SOURCE_DIR}
)

# Copyright (c) 2016-present, Facebook, Inc.
# All rights reserved.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree. An additional grant
# of patent rights can be found in the PATENTS file in the same directory.

#compile the relevant protos
create_proto_dir("orc8r" CPP_OUT_DIR)

list(APPEND PROTO_SRCS "")
list(APPEND PROTO_HDRS "")

set(SERVICE303_PROTOS common metricsd service303)
generate_cpp_protos("${SERVICE303_PROTOS}" "${PROTO_SRCS}" "${PROTO_HDRS}"
  ${ORC8R_PROTO_DIR} ${CPP_OUT_DIR})

set(SERVICE303_GRPC_PROTOS service303)
generate_grpc_protos("${SERVICE303_GRPC_PROTOS}" "${PROTO_SRCS}"
  "${PROTO_HDRS}" ${ORC8R_PROTO_DIR} ${CPP_OUT_DIR})

generate_prometheus_proto("${PROTO_SRCS}" "${PROTO_HDRS}" ${CPP_OUT_DIR})

message("Proto_srcs are ${PROTO_SRCS}")

add_library(SERVICE303_LIB
  MagmaService.cpp
  MetricsSingleton.cpp
  MetricsSingleton.cpp
  ProcFileUtils.cpp
  ${PROTO_SRCS}
  ${PROTO_HDRS}
)

target_link_libraries(SERVICE303_LIB
  prometheus-cpp protobuf grpc grpc++
  SERVICE_REGISTRY LOGGING
)

# copy headers to build directory so they can be shared with OAI,
# session_manager, etc.
add_custom_command(TARGET SERVICE303_LIB POST_BUILD
                   COMMAND ${CMAKE_COMMAND} -E copy
                   ${PROJECT_SOURCE_DIR}/service303/*.h
                   $<TARGET_FILE_DIR:SERVICE303_LIB>)

target_include_directories(SERVICE303_LIB PUBLIC
  $<TARGET_FILE_DIR:SERVICE303_LIB>
  "/usr/local/include/prometheus"
)

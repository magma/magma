# Copyright (c) 2016-present, Facebook, Inc.
# All rights reserved.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree. An additional grant
# of patent rights can be found in the PATENTS file in the same directory.

add_compile_options(-std=c++11)

file(GLOB PROTO_SRCS "orc8r/protos/*.pb.*")
file(GLOB MCONFIG_SRCS "orc8r/protos/mconfig/*.pb.*")
file(GLOB PROMETHEUS_SRCS "orc8r/protos/prometheus/*.pb.*")

include_directories(".")

add_library(MAGMA_PROTOBUF
  ${PROTO_SRCS}
  ${MCONFIG_SRCS}
  ${PROMETHEUS_SRCS}
)
target_link_libraries(MAGMA_PROTOBUF
  protobuf grpc++ grpc dl prometheus-cpp
)
target_include_directories(MAGMA_PROTOBUF SYSTEM PUBLIC
  ${CMAKE_CURRENT_SOURCE_DIR}
)

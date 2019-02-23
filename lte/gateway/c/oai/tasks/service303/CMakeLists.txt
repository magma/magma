# Copyright (c) 2016-present, Facebook, Inc.
# All rights reserved.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree. An additional grant
# of patent rights can be found in the PATENTS file in the same directory.

find_package(Protobuf REQUIRED)
add_compile_options(-std=c++11)

set(OUTPUT_DIR "${CMAKE_CURRENT_BINARY_DIR}")

add_library(TASK_SERVICE303
  service303.cpp
)

target_link_libraries(TASK_SERVICE303
  COMMON
  SERVICE303_LIB
  LIB_BSTR
  ${PROTOBUF_LIBRARIES}
  )

target_include_directories(TASK_SERVICE303 PUBLIC
  ${CMAKE_CURRENT_SOURCE_DIR}
)

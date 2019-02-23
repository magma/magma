# Copyright (c) 2016-present, Facebook, Inc.
# All rights reserved.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree. An additional grant
# of patent rights can be found in the PATENTS file in the same directory.

add_compile_options(-std=c++14)
set(OUTPUT_DIR "${CMAKE_CURRENT_BINARY_DIR}")

include_directories("${PROJECT_SOURCE_DIR}/../common/config")
include_directories("${PROJECT_SOURCE_DIR}/../common/service303")

include_directories("${PROJECT_SOURCE_DIR}/../common/protobuf")

add_library(COMMON_TEST_LIB)

target_link_libraries(COMMON_TEST_LIB
  CONFIG
  gmock_main
  yaml-cpp
  pthread
  rt
  ${GCOV_LIB}
  SERVICE303_LIB
)

foreach(common_test yaml_utils magma_service)
  add_executable(${common_test}_test test_${common_test}.cpp)
  target_link_libraries(${common_test}_test COMMON_TEST_LIB)
  add_test(test_${common_test} ${common_test}_test)
endforeach(common_test)

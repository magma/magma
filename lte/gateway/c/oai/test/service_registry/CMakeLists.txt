# Copyright (c) 2016-present, Facebook, Inc.
# All rights reserved.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree. An additional grant
# of patent rights can be found in the PATENTS file in the same directory.

add_executable(service_registry_test test_service_registry.cpp)

target_link_libraries(service_registry_test
  SERVICE_REGISTRY CONFIG  gtest gtest_main pthread grpc++ dl yaml-cpp
    )

add_test(test_service_registry service_registry_test)

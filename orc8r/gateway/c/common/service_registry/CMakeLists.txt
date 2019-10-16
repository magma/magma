# Copyright (c) 2016-present, Facebook, Inc.
# All rights reserved.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree. An additional grant
# of patent rights can be found in the PATENTS file in the same directory.

add_compile_options(-std=c++11)

include_directories("${PROJECT_SOURCE_DIR}/../common/config")

add_library(SERVICE_REGISTRY
    ServiceRegistrySingleton.h
    ServiceRegistrySingleton.cpp
    )

target_link_libraries(SERVICE_REGISTRY CONFIG
   grpc++ dl yaml-cpp
  )

# copy headers to build directory so they can be shared with OAI,
# session_manager, etc.
add_custom_command(TARGET SERVICE_REGISTRY POST_BUILD
                   COMMAND ${CMAKE_COMMAND} -E copy
                   ${PROJECT_SOURCE_DIR}/service_registry/*.h
                   $<TARGET_FILE_DIR:SERVICE_REGISTRY>)

target_include_directories(SERVICE_REGISTRY PUBLIC
                   $<TARGET_FILE_DIR:SERVICE_REGISTRY>
)

# Copyright (c) 2016-present, Facebook, Inc.
# All rights reserved.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree. An additional grant
# of patent rights can be found in the PATENTS file in the same directory.

add_compile_options(-std=c++14)

include_directories("${PROJECT_SOURCE_DIR}/../common/logging")

add_library(DATASTORE
    RedisMap.hpp
    ObjectMap.h
    Serializers.cpp
    Serializers.h
    )

target_link_libraries(DATASTORE
   tacopie cpp_redis pthread
   )

# copy headers to build directory so they can be shared with OAI,
# session_manager, etc.
add_custom_command(TARGET DATASTORE POST_BUILD
                   COMMAND ${CMAKE_COMMAND} -E copy
                   ${PROJECT_SOURCE_DIR}/datastore/*.h*
                   $<TARGET_FILE_DIR:DATASTORE>)

target_include_directories(DATASTORE PUBLIC
                  $<TARGET_FILE_DIR:DATASTORE>
)

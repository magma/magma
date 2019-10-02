# Copyright (c) 2016-present, Facebook, Inc.
# All rights reserved.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree. An additional grant
# of patent rights can be found in the PATENTS file in the same directory.

add_compile_options(-std=c++11)

include_directories("${PROJECT_SOURCE_DIR}/../common/logging")

add_library(CONFIG
    MConfigLoader.cpp
    MConfigLoader.h
    ServiceConfigLoader.cpp
    ServiceConfigLoader.h
    YAMLUtils.cpp
    YAMLUtils.h
    )

target_link_libraries(CONFIG glog)

# copy headers to build directory so they can be shared with OAI,
# session_manager, etc.
add_custom_command(TARGET CONFIG POST_BUILD
                   COMMAND ${CMAKE_COMMAND} -E copy
                   ${PROJECT_SOURCE_DIR}/config/*.h $<TARGET_FILE_DIR:CONFIG>)

target_include_directories(CONFIG PUBLIC
                  $<TARGET_FILE_DIR:CONFIG>
)

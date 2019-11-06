# Copyright (c) 2016-present, Facebook, Inc.
# All rights reserved.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree. An additional grant
# of patent rights can be found in the PATENTS file in the same directory.

add_compile_options(-std=c++14)

include_directories("${PROJECT_SOURCE_DIR}/../common/logging")

list(APPEND PROTO_SRCS "")
list(APPEND PROTO_HDRS "")

set(ASYNC_ORC8R_CPP_PROTOS common redis)
set(ASYNC_LTE_CPP_PROTOS session_manager policydb subscriberdb)
set(ASYNC_LTE_GRPC_PROTOS session_manager)
set(ASYNC_ORC8R_GRPC_PROTOS "")

generate_all_protos("${ASYNC_LTE_CPP_PROTOS}" "${ASYNC_ORC8R_CPP_PROTOS}"
        "${ASYNC_LTE_GRPC_PROTOS}"
        "${ASYNC_ORC8R_GRPC_PROTOS}" "${PROTO_SRCS}" "${PROTO_HDRS}")

message("Async Proto_srcs are ${PROTO_SRCS}")

add_library(DATASTORE
    RedisMap.hpp
    ObjectMap.h
    Serializers.cpp
    Serializers.h
    ${PROTO_SRCS}
    ${PROTO_HDRS}
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

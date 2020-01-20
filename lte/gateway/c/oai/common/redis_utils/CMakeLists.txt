cmake_minimum_required(VERSION 3.7)
project(redis_utils)

include_directories(${CMAKE_CURRENT_SOURCE_DIR})
include_directories(${CMAKE_CURRENT_BINARY_DIR})

add_library(redis_utils redis_client.cpp)
target_link_libraries(redis_utils ${CONFIG} COMMON cpp_redis tacopie protobuf)


target_include_directories(redis_utils PUBLIC
    ${CMAKE_CURRENT_SOURCE_DIR}
    ${MAGMA_LIB_DIR}/config
)

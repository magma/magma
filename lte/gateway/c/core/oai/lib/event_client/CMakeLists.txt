add_compile_options(-std=c++14)

include($ENV{MAGMA_ROOT}/orc8r/gateway/c/common/CMakeProtoMacros.txt)

set(OUTPUT_DIR "${CMAKE_CURRENT_BINARY_DIR}")

add_library(LIB_EVENT_CLIENT
    EventClientAPI.cpp
    )

target_link_libraries(LIB_EVENT_CLIENT
    COMMON
    LIB_BSTR EVENTD
    ASYNC_GRPC
    )

target_include_directories(LIB_EVENT_CLIENT PUBLIC
    "$<BUILD_INTERFACE:${CMAKE_CURRENT_SOURCE_DIR}>"
    $<TARGET_FILE_DIR:LIB_EVENT_CLIENT>
    ASYNC_GRPC COMMON SERVICE_REGISTRY EVENTD
    )
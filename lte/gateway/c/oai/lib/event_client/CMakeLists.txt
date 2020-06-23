add_compile_options(-std=c++14)

include($ENV{MAGMA_ROOT}/orc8r/gateway/c/common/CMakeProtoMacros.txt)

set(OUTPUT_DIR "${CMAKE_CURRENT_BINARY_DIR}")

set(MAGMA_LIB_DIR $ENV{C_BUILD}/magma_common)
include_directories("${OUTPUT_DIR}")

include_directories(${CMAKE_CURRENT_BINARY_DIR})
include_directories(${MAGMA_LIB_DIR}/async_grpc)
include_directories(${MAGMA_LIB_DIR}/eventd)

link_directories(
        ${MAGMA_LIB_DIR}/async_grpc
        ${MAGMA_LIB_DIR}/eventd
)

add_library(LIB_EVENT_CLIENT
        EventClientAPI.cpp
        )

target_link_libraries(LIB_EVENT_CLIENT
        COMMON
        LIB_BSTR ${EVENTD}
        ${ASYNC_GRPC}
        )

target_include_directories(LIB_EVENT_CLIENT PUBLIC
        ${CMAKE_CURRENT_SOURCE_DIR}
        ASYNC_GRPC COMMON SERVICE_REGISTRY EVENTD
        )

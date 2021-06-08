find_package(Protobuf REQUIRED)
add_compile_options(-std=c++11)

set(OUTPUT_DIR "${CMAKE_CURRENT_BINARY_DIR}")

add_library(TASK_SERVICE303
    service303.cpp
    )

target_link_libraries(TASK_SERVICE303
    COMMON
    SERVICE303_LIB
    ASYNC_GRPC
    SERVICE_REGISTRY
    MAGMA_CONFIG
    LIB_BSTR
    ${PROTOBUF_LIBRARIES}
    )

target_include_directories(TASK_SERVICE303 PUBLIC
    $<TARGET_FILE_DIR:TASK_SERVICE303>
    )

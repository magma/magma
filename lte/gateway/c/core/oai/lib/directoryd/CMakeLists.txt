add_compile_options(-std=c++11)

# compile the needed protos
set(DRD_ORC8R_CPP_PROTOS common directoryd)
set(DRD_LTE_CPP_PROTOS "")
set(DRD_ORC8R_GRPC_PROTOS directoryd)
set(DRD_LTE_GRPC_PROTOS "")

list(APPEND PROTO_SRCS "")
list(APPEND PROTO_HDRS "")

generate_all_protos("${DRD_LTE_CPP_PROTOS}" "${DRD_ORC8R_CPP_PROTOS}"
    "${DRD_LTE_GRPC_PROTOS}" "${DRD_ORC8R_GRPC_PROTOS}" "${PROTO_SRCS}" "${PROTO_HDRS}")

add_library(LIB_DIRECTORYD
    directoryd.cpp
    GatewayDirectorydClient.cpp
    ${PROTO_SRCS}
    ${PROTO_HDRS}
    )

target_link_libraries(LIB_DIRECTORYD
    grpc grpc++
    LIB_MOBILITY_CLIENT ASYNC_GRPC SERVICE_REGISTRY
    MAGMA_CONFIG
    )

target_include_directories(LIB_DIRECTORYD PUBLIC
    ${CMAKE_CURRENT_SOURCE_DIR}
    $<TARGET_FILE_DIR:LIB_DIRECTORYD>
    )

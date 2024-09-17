add_compile_options(-std=c++14)

set(CMAKE_CXX_FLAGS "-Wno-write-strings -Wno-literal-suffix")

# compile the needed protos
set(PCEF_ORC8R_CPP_PROTOS mconfig/mconfigs)
set(PCEF_LTE_CPP_PROTOS policydb session_manager oai/spgw_state oai/sgw_state oai/common_types oai/std_3gpp_types)
set(PCEF_ORC8R_GRPC_PROTOS "")
set(PCEF_LTE_GRPC_PROTOS session_manager)

list(APPEND PROTO_SRCS "")
list(APPEND PROTO_HDRS "")

generate_all_protos("${PCEF_LTE_CPP_PROTOS}" "${PCEF_ORC8R_CPP_PROTOS}"
    "${PCEF_LTE_GRPC_PROTOS}" "${PCEF_ORC8R_GRPC_PROTOS}" "${PROTO_SRCS}" "${PROTO_HDRS}")

add_library(LIB_PCEF
    pcef_handlers.cpp
    PCEFClient.cpp
    ${PROTO_SRCS}
    ${PROTO_HDRS}
    )

target_link_libraries(LIB_PCEF
    COMMON
    ${ITTI_LIB} ASYNC_GRPC SERVICE_REGISTRY
    MAGMA_CONFIG LIB_BSTR LIB_HASHTABLE LIB_MOBILITY_CLIENT yaml-cpp
    )
target_include_directories(LIB_PCEF PUBLIC
    $ENV{MAGMA_ROOT}
    ${CMAKE_CURRENT_SOURCE_DIR}
    $<TARGET_FILE_DIR:LIB_PCEF>
    )

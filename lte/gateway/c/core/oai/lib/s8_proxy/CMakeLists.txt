add_compile_options(-std=c++14)

# compile the needed protos
set(S8C_FEG_CPP_PROTOS s8_proxy)
set(S8C_FEG_GRPC_PROTOS s8_proxy)
set(S8C_LTE_CPP_PROTOS mconfig/mconfigs oai/spgw_state oai/sgw_state oai/common_types oai/std_3gpp_types)

list(APPEND PROTO_SRCS "")
list(APPEND PROTO_HDRS "")

create_proto_dir("feg" FEG_OUT_DIR)
create_proto_dir("lte" LTE_OUT_DIR)

generate_cpp_protos("${S8C_FEG_CPP_PROTOS}" "${PROTO_SRCS}"
    "${PROTO_HDRS}" ${FEG_PROTO_DIR} ${FEG_OUT_DIR})

generate_grpc_protos("${S8C_FEG_GRPC_PROTOS}" "${PROTO_SRCS}"
    "${PROTO_HDRS}" ${FEG_PROTO_DIR} ${FEG_OUT_DIR})

generate_cpp_protos("${S8C_LTE_CPP_PROTOS}" "${PROTO_SRCS}"
    "${PROTO_HDRS}" ${LTE_PROTO_DIR} ${LTE_OUT_DIR})

add_library(LIB_S8_PROXY
    S8Client.cpp
    s8_client_api.cpp
    ${PROTO_SRCS}
    ${PROTO_HDRS}
    )

target_link_libraries(LIB_S8_PROXY
    ASYNC_GRPC SERVICE_REGISTRY MAGMA_CONFIG
    LIB_BSTR LIB_HASHTABLE LIB_PCEF
    )

target_include_directories(LIB_S8_PROXY PUBLIC
    ${CMAKE_CURRENT_SOURCE_DIR}
    $<TARGET_FILE_DIR:LIB_S8_PROXY>
    )

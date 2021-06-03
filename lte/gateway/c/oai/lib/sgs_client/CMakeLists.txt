add_compile_options(-std=c++14)

# compile the needed protos
set(SGSC_FEG_CPP_PROTOS csfb)
set(SGSC_FEG_GRPC_PROTOS csfb)
set(SGSC_LTE_CPP_PROTOS mconfig/mconfigs)

list(APPEND PROTO_SRCS "")
list(APPEND PROTO_HDRS "")

create_proto_dir("feg" FEG_OUT_DIR)
create_proto_dir("lte" LTE_OUT_DIR)
generate_cpp_protos("${SGSC_FEG_CPP_PROTOS}" "${PROTO_SRCS}"
    "${PROTO_HDRS}" ${FEG_PROTO_DIR} ${FEG_OUT_DIR})

generate_grpc_protos("${SGSC_FEG_GRPC_PROTOS}" "${PROTO_SRCS}"
    "${PROTO_HDRS}" ${FEG_PROTO_DIR} ${FEG_OUT_DIR})

generate_cpp_protos("${SGSC_LTE_CPP_PROTOS}" "${PROTO_SRCS}"
    "${PROTO_HDRS}" ${LTE_PROTO_DIR} ${LTE_OUT_DIR})
message("Proto_srcs are ${PROTO_SRCS}")

include_directories(${CMAKE_CURRENT_BINARY_DIR})

add_library(LIB_SGS_CLIENT
    csfb_client_api.cpp
    CSFBClient.cpp
    itti_msg_to_proto_msg.cpp
    ${PROTO_SRCS}
    ${PROTO_HDRS}
    )

target_link_libraries(LIB_SGS_CLIENT
    COMMON
    ${ASYNC_GRPC} ${SERVICE_REGISTRY}
    LIB_BSTR LIB_HASHTABLE ${CONFIG}
    )
target_include_directories(LIB_SGS_CLIENT PUBLIC
    ${MAGMA_LIB_DIR}/async_grpc
    ${MAGMA_LIB_DIR}/config
    ${MAGMA_LIB_DIR}/service_registry
    ${CMAKE_CURRENT_SOURCE_DIR}
    )

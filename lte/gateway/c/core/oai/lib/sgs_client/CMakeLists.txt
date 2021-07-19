add_compile_options(-std=c++14)
set(CMAKE_EXPORT_COMPILE_COMMANDS ON)

# compile the needed protos
set(SGSC_FEG_CPP_PROTOS csfb)
set(SGSC_FEG_GRPC_PROTOS csfb)
set(SGSC_LTE_CPP_PROTOS mconfig/mconfigs)

list(APPEND PROTO_SRCS "")
list(APPEND PROTO_HDRS "")

create_proto_dir("feg" FEG_OUT_DIR)
create_proto_dir("lte" LTE_OUT_DIR)
create_proto_dir("orc8r" ORC8R_CPP_OUT_DIR)

generate_cpp_protos("${SGSC_FEG_CPP_PROTOS}" "${PROTO_SRCS}"
    "${PROTO_HDRS}" ${FEG_PROTO_DIR} ${FEG_OUT_DIR})

generate_grpc_protos("${SGSC_FEG_GRPC_PROTOS}" "${PROTO_SRCS}"
    "${PROTO_HDRS}" ${FEG_PROTO_DIR} ${FEG_OUT_DIR})

generate_cpp_protos("${SGSC_LTE_CPP_PROTOS}" "${PROTO_SRCS}"
    "${PROTO_HDRS}" ${LTE_PROTO_DIR} ${LTE_OUT_DIR})
message("Proto_srcs are ${PROTO_SRCS}")

set(SGS_ROC_ORC8R_CPP_PROTOS common)
generate_cpp_protos("${SGS_ROC_ORC8R_CPP_PROTOS}" "${PROTO_SRCS}"
    "${PROTO_HDRS}" ${ORC8R_PROTO_DIR} ${ORC8R_CPP_OUT_DIR})

message("Proto_srcs are ${PROTO_SRCS}")

add_library(LIB_SGS_CLIENT
    csfb_client_api.cpp
    CSFBClient.cpp
    itti_msg_to_proto_msg.cpp
    ${PROTO_SRCS}
    ${PROTO_HDRS}
    )

target_link_libraries(LIB_SGS_CLIENT
    COMMON
    ASYNC_GRPC SERVICE_REGISTRY MAGMA_CONFIG
    LIB_BSTR LIB_HASHTABLE
    )
target_include_directories(LIB_SGS_CLIENT PUBLIC
    $<TARGET_FILE_DIR:LIB_SGS_CLIENT>
    ${CMAKE_CURRENT_SOURCE_DIR})

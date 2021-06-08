add_compile_options(-std=c++14)

# compile the needed protos
set(SMOC_LTE_CPP_PROTOS sms_orc8r)
set(SMOC_MCONFIG_CPP_PROTOS mconfig/mconfigs)
set(SMOC_LTE_GRPC_PROTOS sms_orc8r)

list(APPEND PROTO_SRCS "")
list(APPEND PROTO_HDRS "")

create_proto_dir("lte" LTE_OUT_DIR)
generate_cpp_protos("${SMOC_LTE_CPP_PROTOS}" "${PROTO_SRCS}"
    "${PROTO_HDRS}" ${LTE_PROTO_DIR} ${LTE_OUT_DIR})

generate_grpc_protos("${SMOC_LTE_GRPC_PROTOS}" "${PROTO_SRCS}"
    "${PROTO_HDRS}" ${LTE_PROTO_DIR} ${LTE_OUT_DIR})

generate_cpp_protos("${SMOC_MCONFIG_CPP_PROTOS}" "${PROTO_SRCS}"
    "${PROTO_HDRS}" ${LTE_PROTO_DIR} ${LTE_OUT_DIR})

message("Proto_srcs are ${PROTO_SRCS}")

add_library(LIB_SMS_ORC8R_CLIENT
    sms_orc8r_client_api.cpp
    SMSOrc8rClient.cpp
    itti_msg_to_proto_msg.cpp
    ${PROTO_SRCS}
    ${PROTO_HDRS}
    )

target_link_libraries(LIB_SMS_ORC8R_CLIENT
    COMMON
    ASYNC_GRPC SERVICE_REGISTRY MAGMA_CONFIG
    LIB_BSTR LIB_HASHTABLE)

target_include_directories(LIB_SMS_ORC8R_CLIENT PUBLIC
    ${CMAKE_CURRENT_SOURCE_DIR}
    $<TARGET_FILE_DIR:LIB_SMS_ORC8R_CLIENT>
    )

find_package(Protobuf REQUIRED)
set(CMAKE_CXX_FLAGS "-Wno-write-strings -Wno-literal-suffix")

# compile the needed protos
list(APPEND PROTO_SRCS "")
list(APPEND PROTO_HDRS "")

create_proto_dir("lte" LTE_OUT_DIR)
create_proto_dir("orc8r" ORC8R_OUT_DIR)
create_proto_dir("feg" FEG_OUT_DIR)

# Common protos
set(ORC8R_CPP_PROTOS common digest)
generate_cpp_protos("${ORC8R_CPP_PROTOS}" "${PROTO_SRCS}"
    "${PROTO_HDRS}" ${ORC8R_PROTO_DIR} ${ORC8R_OUT_DIR})

# SPGW
set(SPGWSRV_LTE_CPP_PROTOS spgw_service apn subscriberdb policydb mobilityd session_manager)
set(SPGWSRV_LTE_GRPC_PROTOS spgw_service session_manager)
generate_cpp_protos("${SPGWSRV_LTE_CPP_PROTOS}" "${PROTO_SRCS}"
    "${PROTO_HDRS}" ${LTE_PROTO_DIR} ${LTE_OUT_DIR})
generate_grpc_protos("${SPGWSRV_LTE_GRPC_PROTOS}" "${PROTO_SRCS}"
    "${PROTO_HDRS}" ${LTE_PROTO_DIR} ${LTE_OUT_DIR})

# S6a
set(S6ASRV_FEG_CPP_PROTOS s6a_proxy)
set(S6ASRV_FEG_GRPC_PROTOS s6a_proxy)
set(S6ASRV_LTE_CPP_PROTOS s6a_service)
set(S6ASRV_LTE_GRPC_PROTOS s6a_service)
generate_cpp_protos("${S6ASRV_FEG_CPP_PROTOS}" "${PROTO_SRCS}"
    "${PROTO_HDRS}" ${FEG_PROTO_DIR} ${FEG_OUT_DIR})
generate_grpc_protos("${S6ASRV_FEG_GRPC_PROTOS}" "${PROTO_SRCS}"
    "${PROTO_HDRS}" ${FEG_PROTO_DIR} ${FEG_OUT_DIR})
generate_cpp_protos("${S6ASRV_LTE_CPP_PROTOS}" "${PROTO_SRCS}"
    "${PROTO_HDRS}" ${LTE_PROTO_DIR} ${LTE_OUT_DIR})
generate_grpc_protos("${S6ASRV_LTE_GRPC_PROTOS}" "${PROTO_SRCS}"
    "${PROTO_HDRS}" ${LTE_PROTO_DIR} ${LTE_OUT_DIR})

# SGs
set(SGSSRV_FEG_CPP_PROTOS csfb)
set(SGSSRV_FEG_GRPC_PROTOS csfb)
generate_cpp_protos("${SGSSRV_FEG_CPP_PROTOS}" "${PROTO_SRCS}"
    "${PROTO_HDRS}" ${FEG_PROTO_DIR} ${FEG_OUT_DIR})
generate_grpc_protos("${SGSSRV_FEG_GRPC_PROTOS}" "${PROTO_SRCS}"
    "${PROTO_HDRS}" ${FEG_PROTO_DIR} ${FEG_OUT_DIR})

# SMS Orc8r
set(SMOSRV_LTE_CPP_PROTOS sms_orc8r)
set(SMOSRV_LTE_GRPC_PROTOS sms_orc8r)
generate_cpp_protos("${SMOSRV_LTE_CPP_PROTOS}" "${PROTO_SRCS}"
    "${PROTO_HDRS}" ${LTE_PROTO_DIR} ${LTE_OUT_DIR})
generate_grpc_protos("${SMOSRV_LTE_GRPC_PROTOS}" "${PROTO_SRCS}"
    "${PROTO_HDRS}" ${LTE_PROTO_DIR} ${LTE_OUT_DIR})

# S1ap Service
set(S1APSRV_LTE_CPP_PROTOS s1ap_service oai/s1ap_state)
set(S1APSRV_LTE_GRPC_PROTOS s1ap_service oai/s1ap_state)
generate_cpp_protos("${S1APSRV_LTE_CPP_PROTOS}" "${PROTO_SRCS}"
    "${PROTO_HDRS}" ${LTE_PROTO_DIR} ${LTE_OUT_DIR})
generate_grpc_protos("${S1APSRV_LTE_GRPC_PROTOS}" "${PROTO_SRCS}"
    "${PROTO_HDRS}" ${LTE_PROTO_DIR} ${LTE_OUT_DIR})

# HA Service
set(HASRV_LTE_CPP_PROTOS ha_service)
set(HASRV_LTE_GRPC_PROTOS ha_service)
generate_cpp_protos("${HASRV_LTE_CPP_PROTOS}" "${PROTO_SRCS}"
    "${PROTO_HDRS}" ${LTE_PROTO_DIR} ${LTE_OUT_DIR})
generate_grpc_protos("${HASRV_LTE_GRPC_PROTOS}" "${PROTO_SRCS}"
    "${PROTO_HDRS}" ${LTE_PROTO_DIR} ${LTE_OUT_DIR})

# AMF
set(AMFSRV_M5G_CPP_PROTOS apn subscriberdb session_manager diam_errors subscriberauth)
set(AMFSRV_M5G_GRPC_PROTOS session_manager subscriberauth)
generate_cpp_protos("${AMFSRV_M5G_CPP_PROTOS}" "${PROTO_SRCS}"
    "${PROTO_HDRS}" ${LTE_PROTO_DIR} ${LTE_OUT_DIR})
generate_grpc_protos("${AMFSRV_M5G_GRPC_PROTOS}" "${PROTO_SRCS}"
    "${PROTO_HDRS}" ${LTE_PROTO_DIR} ${LTE_OUT_DIR})

# S8a
set(S8SRV_FEG_CPP_PROTOS s8_proxy)
set(S8SRV_FEG_GRPC_PROTOS s8_proxy)
set(S8SRV_LTE_CPP_PROTOS oai/spgw_state  oai/common_types oai/std_3gpp_types)
set(S8SRV_LTE_GRPC_PROTOS oai/spgw_state  oai/common_types oai/std_3gpp_types)
generate_cpp_protos("${S8SRV_FEG_CPP_PROTOS}" "${PROTO_SRCS}"
  "${PROTO_HDRS}" ${FEG_PROTO_DIR} ${FEG_OUT_DIR})
generate_grpc_protos("${S8SRV_FEG_GRPC_PROTOS}" "${PROTO_SRCS}"
  "${PROTO_HDRS}" ${FEG_PROTO_DIR} ${FEG_OUT_DIR})
generate_cpp_protos("${S8SRV_LTE_CPP_PROTOS}" "${PROTO_SRCS}"
  "${PROTO_HDRS}" ${LTE_PROTO_DIR} ${LTE_OUT_DIR})
generate_grpc_protos("${S8SRV_LTE_GRPC_PROTOS}" "${PROTO_SRCS}"
  "${PROTO_HDRS}" ${LTE_PROTO_DIR} ${LTE_OUT_DIR})

include_directories(${CMAKE_CURRENT_BINARY_DIR})
include_directories("${MAGMA_ROOT}/orc8r/gateway/c/common/logging")

if (EMBEDDED_SGW)
add_library(TASK_GRPC_SERVICE
    AmfServiceImpl.cpp
    SpgwServiceImpl.cpp
    S6aServiceImpl.cpp
    CSFBGatewayServiceImpl.cpp
    SMSOrc8rGatewayServiceImpl.cpp
    S1apServiceImpl.cpp
    HaServiceImpl.cpp
    AmfServiceImpl.cpp
    spgw_service_handler.c
    amf_service_handler.c
    proto_msg_to_itti_msg.cpp
    grpc_service.cpp
    grpc_service_task.c
    S8ServiceImpl.cpp
    ${PROTO_SRCS}
    ${PROTO_HDRS}
    )
else (EMBEDDED_SGW)
add_library(TASK_GRPC_SERVICE
    AmfServiceImpl.cpp
    SpgwServiceImpl.cpp
    S6aServiceImpl.cpp
    CSFBGatewayServiceImpl.cpp
    SMSOrc8rGatewayServiceImpl.cpp
    S1apServiceImpl.cpp
    HaServiceImpl.cpp
    AmfServiceImpl.cpp
    spgw_service_handler.c
    amf_service_handler.c
    proto_msg_to_itti_msg.cpp
    grpc_service.cpp
    grpc_service_task.c
    ${PROTO_SRCS}
    ${PROTO_HDRS}
    )
endif (EMBEDDED_SGW)

target_link_libraries(TASK_GRPC_SERVICE
    COMMON
    LIB_BSTR protobuf
    LIB_HASHTABLE
    ${PROTOBUF_LIBRARIES}
    grpc++
    TASK_SGS
    TASK_S6A
    TASK_SMS_ORC8R
    )
target_include_directories(TASK_GRPC_SERVICE PUBLIC
    ${CMAKE_CURRENT_SOURCE_DIR}
    ${CMAKE_CURRENT_BINARY_DIR}
    $<TARGET_FILE_DIR:TASK_GRPC_SERVICE>
    )

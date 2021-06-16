set(NGAP_C_DIR ${PROJECT_BINARY_DIR}/ngap/r16.0)
set(NAS5G_C_DIR ${PROJECT_SOURCE_DIR}/tasks/nas5g)
set(NAS_DIR ${PROJECT_SOURCE_DIR}/tasks/nas)
list(APPEND PROTO_SRCS "")
list(APPEND PROTO_HDRS "")

create_proto_dir("feg" FEG_OUT_DIR)

# S6a
set(S6ASRV_FEG_CPP_PROTOS s6a_proxy)
set(S6ASRV_FEG_GRPC_PROTOS s6a_proxy)
generate_cpp_protos("${S6ASRV_FEG_CPP_PROTOS}" "${PROTO_SRCS}"
  "${PROTO_HDRS}" ${FEG_PROTO_DIR} ${FEG_OUT_DIR})
generate_grpc_protos("${S6ASRV_FEG_GRPC_PROTOS}" "${PROTO_SRCS}"
  "${PROTO_HDRS}" ${FEG_PROTO_DIR} ${FEG_OUT_DIR})

include_directories(${CMAKE_CURRENT_BINARY_DIR})
include_directories("${NAS5G_C_DIR}")
include_directories("${NAS5G_C_DIR}/include")
include_directories("${NAS5G_C_DIR}/include/ies")
include_directories("${PROJECT_SOURCE_DIR}/include/nas")
include_directories(${CMAKE_CURRENT_SOURCE_DIR}/../../../../../../orc8r/gateway/c/common/logging)
include_directories(${CMAKE_CURRENT_SOURCE_DIR}/../../common/glogwrapper)

set(MAGMA_LIB_DIR $ENV{C_BUILD}/magma_common)
link_directories(
        ${MAGMA_LIB_DIR}/async_grpc
        ${MAGMA_LIB_DIR}/eventd
)

add_library(TASK_AMF_APP
    amf_fsm.cpp
    amf_app_main.cpp
    amf_config.c
    amf_app_handler.cpp
    amf_app_state_manager.cpp
    amf_app_msg.cpp
    amf_app_transport.cpp
    amf_app_ue_context.cpp
    amf_as.cpp
    amf_authentication.cpp
    amf_data.cpp
    amf_cn.cpp
    amf_sap.cpp
    nas_proc.cpp
    Registration.cpp
    amf_recv.cpp
    amf_smf_send.cpp
    prepare_request_for_smf.cpp
    amf_identity.cpp
    amf_security_mode_control.cpp
    amf_Security_Mode.cpp
    deregistration_request.cpp
    amf_app_pdu_resource_setup_req_rel.cpp
    nas5g_message.cpp
    ${PROTO_SRCS}
    ${PROTO_HDRS}
)

target_compile_definitions(TASK_AMF_APP PRIVATE
  PACKAGE_NAME=\"AMF\"
  PACKAGE_VERSION=\"0.1\"
  PACKAGE_BUGREPORT=\"TBD\"
)

target_link_libraries(TASK_AMF_APP 
  ${CONFIG_LIBRARIES}
  COMMON
  LIB_BSTR LIB_HASHTABLE LIB_DIRECTORYD LIB_SECU LIB_EVENT_CLIENT LIB_NAS5G LIB_N11
  TASK_NGAP TASK_SERVICE303 LIB_S6A_PROXY
  protobuf cpp_redis yaml-cpp redis_utils
)

target_include_directories(TASK_AMF_APP PUBLIC
  ${CMAKE_CURRENT_SOURCE_DIR}
    )

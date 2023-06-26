set(S1AP_C_DIR ${PROJECT_BINARY_DIR}/s1ap/r15)

# compile the needed protos
set(MMEAPP_STATE_CPP_PROTOS mme_nas_state nas_state spgw_state common_types
    std_3gpp_types)
list(APPEND PROTO_SRCS "")
list(APPEND PROTO_HDRS "")

set(LTE_OUT_DIR "${CMAKE_CURRENT_BINARY_DIR}/lte/protos/oai")
generate_cpp_protos("${MMEAPP_STATE_CPP_PROTOS}" "${PROTO_SRCS}"
    "${PROTO_HDRS}" ${STATE_PROTO_DIR} ${LTE_OUT_DIR})

if (EMBEDDED_SGW)
  set(S11_RELATED_SRCS mme_app_embedded_spgw.cpp)
endif (EMBEDDED_SGW)

add_library(TASK_MME_APP
    mme_app_capabilities.cpp
    mme_app_context.cpp
    mme_app_if_nas_transport.cpp
    mme_app_main.cpp
    mme_app_bearer.cpp
    mme_app_detach.cpp
    mme_app_location.cpp
    mme_app_transport.cpp
    mme_app_ue_context.cpp
    mme_app_statistics.cpp
    mme_config.cpp
    s6a_2_nas_cause.cpp
    mme_app_purge_ue.cpp
    mme_app_hss_reset.cpp
    mme_app_sgs_detach.cpp
    mme_app_sgsap_location_update.cpp
    mme_app_sgs_fsm.cpp
    mme_app_sgs_null.cpp
    mme_app_sgs_la_updated.cpp
    mme_app_sgs_associated.cpp
    mme_app_sgs_reset.cpp
    mme_app_sgs_alert.cpp
    mme_app_sgsap_service_abort.cpp
    mme_app_sgs_paging.cpp
    mme_app_apn_selection.cpp
    mme_app_bearer_context.cpp
    mme_app_pdn_context.cpp
    mme_app_procedures.cpp
    mme_app_itti_messaging.cpp
    mme_app_edns_emulation.cpp
    mme_app_sgw_selection.cpp
    mme_app_sgs_status.cpp
    mme_app_state_converter.cpp
    mme_app_state_manager.cpp
    mme_app_state.cpp
    mme_events.cpp
    mme_app_ha.cpp
    mme_app_timer_management.cpp
    mme_app_ip_imsi.cpp
    ${PROTO_SRCS}
    ${PROTO_HDRS}
    ${S11_RELATED_SRCS}
    )

target_compile_definitions(TASK_MME_APP PRIVATE
    PACKAGE_NAME=\"MME\"
    PACKAGE_VERSION=\"0.1\"
    PACKAGE_BUGREPORT=\"TBD\"
    )
target_link_libraries(TASK_MME_APP
    ${CONFIG_LIBRARIES}
    COMMON
    LIB_BSTR LIB_HASHTABLE LIB_DIRECTORYD LIB_SECU LIB_EVENT_CLIENT
    TASK_NAS TASK_S1AP TASK_SERVICE303 TASK_SCTP_SERVER TASK_S6A TASK_SGS
    protobuf cpp_redis yaml-cpp redis_utils
    )

if (EMBEDDED_SGW)
  target_link_libraries(TASK_MME_APP TASK_SGW)
else (EMBEDDED_SGW)
  target_link_libraries(TASK_MME_APP TASK_S11 TASK_UDP)
endif (EMBEDDED_SGW)

target_include_directories(TASK_MME_APP PUBLIC
    ${CMAKE_CURRENT_SOURCE_DIR}
    ${S1AP_C_DIR}
    $<TARGET_FILE_DIR:TASK_MME_APP>
    )

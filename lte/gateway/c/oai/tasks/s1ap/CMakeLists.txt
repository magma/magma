add_list1_option(S1AP_VERSION R10 "S1AP Asn.1 grammar version" R8 R9 R10)
add_definitions(-DASN1_MINIMUM_VERSION=924)


set(ASN1RELDIR r10.5)
set(S1AP_DIR ${CMAKE_CURRENT_SOURCE_DIR})
set(asn1_generated_dir ${CMAKE_CURRENT_BINARY_DIR})
set(asn1c_call ${CMAKE_CURRENT_SOURCE_DIR}/generate_asn1)

set(S1AP_ASN_DIR ${CMAKE_CURRENT_SOURCE_DIR}/messages/asn1/${ASN1RELDIR})
set(S1AP_ASN_FILES
    ${S1AP_ASN_DIR}/S1AP-CommonDataTypes.asn
    ${S1AP_ASN_DIR}/S1AP-Constants.asn
    ${S1AP_ASN_DIR}/S1AP-IEs.asn
    ${S1AP_ASN_DIR}/S1AP-PDU.asn
    )
set(S1AP_C_DIR ${asn1_generated_dir}/${ASN1RELDIR})

include_directories(${S1AP_C_DIR})
include_directories(${S1AP_DIR})

execute_process(COMMAND ${asn1c_call} ${S1AP_C_DIR} ${S1AP_ASN_FILES})
file(GLOB S1AP_source ${S1AP_C_DIR}/*.c)
list(REMOVE_ITEM S1AP_source ${S1AP_C_DIR}/converter-sample.c)

set(S1AP_OAI_generated
    ${S1AP_C_DIR}/s1ap_decoder.c
    ${S1AP_C_DIR}/s1ap_encoder.c
    ${S1AP_C_DIR}/s1ap_xer_print.c
    ${S1AP_C_DIR}/s1ap_ies_defs.h
    )

file(GLOB s1ap_h ${S1AP_C_DIR}/*.h)
set(s1ap_h ${s1ap_h})

add_custom_command(
    OUTPUT ${S1AP_OAI_generated}
    COMMAND ${asn1c_call} ${S1AP_C_DIR} ${S1AP_ASN_FILES}
    COMMAND python ${S1AP_DIR}/messages/asn1/asn1tostruct.py -f${S1AP_ASN_DIR}/S1AP-PDU-Contents.asn -o${S1AP_C_DIR}
    DEPENDS ${S1AP_ASN_FILES}
)

add_library(LIB_S1AP
    ${S1AP_OAI_generated}
    ${S1AP_source}
    s1ap_common.c
)
target_link_libraries(LIB_S1AP
    LIB_BSTR LIB_HASHTABLE
)
target_include_directories(LIB_S1AP PUBLIC
    ${S1AP_C_DIR}
    ${CMAKE_CURRENT_SOURCE_DIR}/r10.5
)

set(PCEF_LTE_CPP_PROTOS s1ap_state common_types)

list(APPEND PROTO_SRCS "")
list(APPEND PROTO_HDRS "")

set(STATE_OUT_DIR "${CMAKE_CURRENT_BINARY_DIR}/lte/protos/oai")
generate_cpp_protos("${PCEF_LTE_CPP_PROTOS}" "${PROTO_SRCS}" "${PROTO_HDRS}"
        "${STATE_PROTO_DIR}" "${STATE_OUT_DIR}")

add_library(TASK_S1AP
    "${PROTO_SRCS}"
    "${PROTO_HDRS}"
    ${S1AP_C_DIR}/s1ap_ies_defs.h
    ${S1AP_DIR}/s1ap_mme_encoder.c
    ${S1AP_DIR}/s1ap_mme_decoder.c
    ${S1AP_DIR}/s1ap_mme_handlers.c
    ${S1AP_DIR}/s1ap_mme_nas_procedures.c
    ${S1AP_DIR}/s1ap_mme.c
    ${S1AP_DIR}/s1ap_mme_itti_messaging.c
    ${S1AP_DIR}/s1ap_mme_ta.c
    ${S1AP_DIR}/s1ap_state.cpp
    ${S1AP_DIR}/s1ap_state_manager.cpp
    ${S1AP_DIR}/s1ap_state_converter.cpp
    )
target_link_libraries(TASK_S1AP
    ${CONFIG}
    COMMON
    LIB_BSTR LIB_HASHTABLE
    TASK_SERVICE303 TASK_MME_APP
    cpp_redis tacopie
)
target_include_directories(TASK_S1AP PUBLIC
    ${CMAKE_CURRENT_BINARY_DIR}
    ${CMAKE_CURRENT_SOURCE_DIR}
    ${CMAKE_CURRENT_SOURCE_DIR}/messages
    ${CMAKE_CURRENT_SOURCE_DIR}/messages/asn1
    ${CMAKE_CURRENT_SOURCE_DIR}/messages/asn1/r10.5
    ${CMAKE_CURRENT_SOURCE_DIR}/messages/asn1/r8.10
    ${CMAKE_CURRENT_SOURCE_DIR}/messages/asn1/r9.8
)

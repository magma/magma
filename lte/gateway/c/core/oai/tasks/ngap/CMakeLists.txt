add_list1_option(NGAP_VERSION R16 "NGAP Asn.1 grammar version" R16)
add_definitions(-DASN1_MINIMUM_VERSION=923)
set(ASN1RELDIR r16)
set(NGAP_DIR ${CMAKE_CURRENT_SOURCE_DIR})
set(asn1_generated_dir ${CMAKE_CURRENT_BINARY_DIR})
set(NGAP_ASN_DIR ${CMAKE_CURRENT_SOURCE_DIR}/messages/asn1/${ASN1RELDIR})
set(NGAP_ASN_FILES "r16.asn1")

set(NGAP_C_DIR ${CMAKE_CURRENT_BINARY_DIR}/${ASN1RELDIR})
include_directories(${NGAP_C_DIR})
include_directories(${NGAP_DIR})
set(GENERATED_FULL_DIR ${NGAP_C_DIR})
set(ASN1_SOURCE_DIR ${NGAP_ASN_DIR}/${NGAP_ASN_FILES})
set(ENV{ASN1C_PREFIX} "Ngap_")
set(ngap_generate_code_done_flag ${GENERATED_FULL_DIR}/NGAP_GENERATE_CODE_DONE)

file(MAKE_DIRECTORY ${GENERATED_FULL_DIR})
set(ngap_generate_code_done_flag ${GENERATED_FULL_DIR}/NGAP_GENERATE_CODE_DONE)

if (${ASN1_SOURCE_DIR} IS_NEWER_THAN ${ngap_generate_code_done_flag})
  file(REMOVE ${GENERATED_FULL_DIR}/${ASN1C_PREFIX}*.c ${GENERATED_FULL_DIR}/${ASN1C_PREFIX}*.h)
  file(MAKE_DIRECTORY ${GENERATED_FULL_DIR})

  #Create NGAP_PDU structure without pointer references
  execute_process(COMMAND asn1c -pdu=all -fcompound-names -fno-include-deps -gen-PER -D ${GENERATED_FULL_DIR} ${ASN1_SOURCE_DIR} RESULT_VARIABLE ret)
  if (NOT ${ret} STREQUAL 0)
    message(FATAL_ERROR "${ret}")
  endif (NOT ${ret} STREQUAL 0)
  execute_process(COMMAND bash "-c" "egrep -lRZ \"18446744073709551615\" ${GENERATED_FULL_DIR} | xargs -0 -l sed -i -e \"s/18446744073709551615/18446744073709551615u/g\"")
endif ()
# TOUCH not in cmake 3.10
file(WRITE ${ngap_generate_code_done_flag})

file(GLOB NGAP_source ${NGAP_C_DIR}/*.c)
list(REMOVE_ITEM NGAP_source ${NGAP_C_DIR}/converter-example.c)

file(GLOB ngap_h ${NGAP_C_DIR}/*.c)
set(ngap_h ${ngap_h})

add_library(LIB_NGAP
    ${NGAP_source}
    ngap_common.c
    )
target_link_libraries(LIB_NGAP
    LIB_BSTR LIB_HASHTABLE
    )
target_include_directories(LIB_NGAP PUBLIC
    ${NGAP_C_DIR}
    ${CMAKE_CURRENT_SOURCE_DIR}/r16
    )

set(NGAP_STATE_CPP_PROTOS ngap_state common_types)

list(APPEND PROTO_SRCS "")
list(APPEND PROTO_HDRS "")

set(STATE_OUT_DIR "${CMAKE_CURRENT_BINARY_DIR}/lte/protos/oai")
generate_cpp_protos("${NGAP_STATE_CPP_PROTOS}" "${PROTO_SRCS}" "${PROTO_HDRS}"
    "${STATE_PROTO_DIR}" "${STATE_OUT_DIR}")

add_library(TASK_NGAP
    "${PROTO_SRCS}"
    "${PROTO_HDRS}"
    ${NGAP_DIR}/ngap_amf.c
    ${NGAP_DIR}/ngap_amf_encoder.c
    ${NGAP_DIR}/ngap_amf_decoder.c
    ${NGAP_DIR}/ngap_amf_handlers.c
    ${NGAP_DIR}/ngap_amf_nas_procedures.c
    ${NGAP_DIR}/ngap_amf.c
    ${NGAP_DIR}/ngap_amf_itti_messaging.c
    ${NGAP_DIR}/ngap_amf_ta.c
    ${NGAP_DIR}/ngap_state.cpp
    ${NGAP_DIR}/ngap_state_manager.cpp
    ${NGAP_DIR}/ngap_state_converter.cpp
    )
target_link_libraries(TASK_NGAP
    MAGMA_CONFIG
    COMMON
    LIB_BSTR LIB_HASHTABLE
    TASK_SERVICE303 TASK_MME_APP
    cpp_redis tacopie
    )
target_include_directories(TASK_NGAP PUBLIC
    ${CMAKE_CURRENT_BINARY_DIR}
    ${CMAKE_CURRENT_SOURCE_DIR}
    ${CMAKE_CURRENT_SOURCE_DIR}/messages
    ${CMAKE_CURRENT_SOURCE_DIR}/messages/asn1
    ${CMAKE_CURRENT_SOURCE_DIR}/messages/asn1/r16
    )

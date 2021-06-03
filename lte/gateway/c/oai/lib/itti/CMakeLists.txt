# TODO: ITTI_Messages shouldn't depend on S1_Types
set(ASN1RELDIR r15)
set(asn1_generated_dir ${PROJECT_BINARY_DIR}/s1ap)
set(asn1c_call ${CMAKE_CURRENT_SOURCE_DIR}/generate_asn1)
set(S1AP_C_DIR ${asn1_generated_dir}/${ASN1RELDIR})
include_directories("${S1AP_C_DIR}")

set(ITTI_FILES
    intertask_interface.c
    signals.c
    timer.c
    )
add_library(LIB_ITTI ${ITTI_FILES})
target_link_libraries(LIB_ITTI
    rt pthread czmq
    COMMON
    LIB_BSTR LIB_HASHTABLE LIB_MESSAGE_UTILS
    )
target_include_directories(LIB_ITTI PUBLIC
    ${CMAKE_CURRENT_SOURCE_DIR}
    )

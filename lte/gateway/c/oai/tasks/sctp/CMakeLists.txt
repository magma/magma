
set(S1AP_C_DIR ${PROJECT_BINARY_DIR}/s1ap/r10.5)
include_directories(${S1AP_C_DIR})

add_library(TASK_SCTP_SERVER
    sctp_common.c
    sctp_itti_messaging.c
    sctp_primitives_server.c
    )
target_link_libraries(TASK_SCTP_SERVER
    COMMON
    LIB_BSTR LIB_HASHTABLE
)
target_include_directories(TASK_SCTP_SERVER PUBLIC
    ${CMAKE_CURRENT_SOURCE_DIR}
)

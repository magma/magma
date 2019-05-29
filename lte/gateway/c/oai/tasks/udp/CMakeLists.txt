set(S1AP_C_DIR ${PROJECT_BINARY_DIR}/s1ap/r10.5)
include_directories(${S1AP_C_DIR})

add_library(TASK_UDP_SERVER
    udp_primitives_server.c
)
target_link_libraries(TASK_UDP_SERVER
    COMMON
    LIB_BSTR LIB_HASHTABLE
)
target_include_directories(TASK_UDP_SERVER PUBLIC
    ${CMAKE_CURRENT_SOURCE_DIR}
)

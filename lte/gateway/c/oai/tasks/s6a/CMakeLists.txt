set(S1AP_C_DIR ${PROJECT_BINARY_DIR}/s1ap/r10.5)
include_directories(${S1AP_C_DIR})

add_library(TASK_S6A
    s6a_auth_info.c
    s6a_dict.c
    s6a_error.c
    s6a_peer.c
    s6a_subscription_data.c
    s6a_task.c
    s6a_up_loc.c
    s6a_cancel_loc.c
    s6a_service_handler.c
    s6a_purge_ue.c
    s6a_hss_reset.c
    s6a_fd_iface.cpp
    s6a_grpc_iface.cpp
    s6a_iface.cpp
    )
target_link_libraries(TASK_S6A
    COMMON
    LIB_BSTR LIB_HASHTABLE LIB_S6A_PROXY
)
target_include_directories(TASK_S6A PUBLIC
    ${CMAKE_CURRENT_SOURCE_DIR}
)

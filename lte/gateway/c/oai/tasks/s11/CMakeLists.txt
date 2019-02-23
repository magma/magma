add_library(TASK_S11_MME
    s11_common.c
    s11_ie_formatter.c
    s11_mme_task.c
    s11_mme_bearer_manager.c
    s11_mme_session_manager.c
    )
target_link_libraries(TASK_S11_MME
    COMMON
    LIB_BSTR LIB_HASHTABLE LIB_GTPV2C
    TASK_NAS TASK_MME_APP TASK_SERVICE303 TASK_SGW
)
target_include_directories(TASK_S11_MME PUBLIC ${CMAKE_CURRENT_SOURCE_DIR})

add_library(TASK_S11_SGW
    s11_common.c
    s11_ie_formatter.c
    s11_sgw.c
    s11_sgw_session_manager.c
    s11_sgw_bearer_manager.c
    )
target_link_libraries(TASK_S11_SGW
    COMMON
    LIB_BSTR LIB_HASHTABLE LIB_GTPV2C
    TASK_NAS TASK_MME_APP TASK_SERVICE303 TASK_SGW
)
target_include_directories(TASK_S11_SGW PUBLIC ${CMAKE_CURRENT_SOURCE_DIR})

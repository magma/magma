
include_directories(${S1AP_C_DIR})


add_library(TASK_S11
s11_mme_bearer_manager.c 
s11_mme_session_manager.c
s11_tasks.c
)

target_link_libraries(TASK_S11
    COMMON
    LIB_BSTR LIB_HASHTABLE LIB_SGS_CLIENT
    ${PROTOBUF_LIBRARIES}
)
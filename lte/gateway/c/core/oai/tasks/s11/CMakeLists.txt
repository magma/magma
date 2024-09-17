
include_directories(${S1AP_C_DIR})
include_directories(${PROJECT_SOURCE_DIR}/lib/gtpv2-c/gtpv2c_ie_formatter/shared)
include_directories(converged-mme/lte/gateway/c/oai/lib/gtpv2-c/nwgtpv2c-0.11/include)
include_directories(converged-mme/lte/gateway/c/oai/lib/gtpv2-c/nwgtpv2c-0.11/shared)

add_library(TASK_S11
    s11_mme_bearer_manager.cpp
    s11_mme_session_manager.cpp
    s11_tasks.cpp
    s11_common.cpp
    s11_ie_formatter.cpp
    )

target_link_libraries(TASK_S11
    COMMON
    LIB_BSTR LIB_HASHTABLE LIB_SGS_CLIENT LIB_GTPV2C TASK_MME_APP TASK_UDP
    ${PROTOBUF_LIBRARIES}
    )

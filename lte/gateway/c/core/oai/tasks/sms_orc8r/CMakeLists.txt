find_package(Protobuf REQUIRED)

add_library(TASK_SMS_ORC8R
    sms_orc8r_task.c
    sms_orc8r_service_handler.c
    )

target_link_libraries(TASK_SMS_ORC8R
    COMMON
    LIB_BSTR LIB_HASHTABLE LIB_SMS_ORC8R_CLIENT
    ${PROTOBUF_LIBRARIES}
    )
target_include_directories(TASK_SMS_ORC8R PUBLIC
    ${CMAKE_CURRENT_SOURCE_DIR}
    )

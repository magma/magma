set (SGW_S8_SRC
    sgw_s8_task.c
    )

add_library(TASK_SGW_S8 ${SGW_S8_SRC})

target_link_libraries(TASK_SGW_S8 COMMON LIB_BSTR LIB_HASHTABLE)

target_include_directories(TASK_SGW_S8 PUBLIC
    ${CMAKE_CURRENT_SOURCE_DIR}
)

find_package(Check REQUIRED)
find_package(Threads REQUIRED)

set(MME_APP_UE_CONTEXT_IMSI_SRC
    test_mme_app_ue_context.c
)

add_executable(test_mme_app_ue_context_imsi ${MME_APP_UE_CONTEXT_IMSI_SRC})
target_link_libraries(test_mme_app_ue_context_imsi
    TASK_MME_APP ${CHECK_LIBRARIES} ${CMAKE_THREAD_LIBS_INIT}
    LIB_BSTR LIB_HASHTABLE
)
target_include_directories(test_mme_app_ue_context_imsi PUBLIC
    ${CMAKE_CURRENT_SOURCE_DIR}
    ${CHECK_INCLUDE_DIRS}
)

add_test(NAME test_mme_app_ue_context COMMAND test_mme_app_ue_context_imsi)

add_subdirectory(mobility_client)
add_subdirectory(openflow)
# Currently broken due to include error.
# add_subdirectory(service303)
# add_subdirectory(service_registry)

add_compile_options(-std=c++11)

set(OUTPUT_DIR "${CMAKE_CURRENT_BINARY_DIR}")

include_directories("${PROJECT_SOURCE_DIR}/openflow/controller")
add_executable(openflow_controller_test test_openflow_controller.cpp)
add_executable(imsi_encoder_test test_imsi_encoder.cpp)
add_executable(gtp_app_test test_gtp_app.cpp)

add_library(OPENFLOW_TEST openflow_mocks.h)
target_link_libraries(OPENFLOW_TEST
    COMMON
    lfds710
    LIB_OPENFLOW_CONTROLLER  LIB_BSTR LIB_HASHTABLE LIB_ITTI LIB_S1AP TASK_S1AP
    gmock_main pthread rt)

set_target_properties(OPENFLOW_TEST PROPERTIES LINKER_LANGUAGE CXX)

target_link_libraries(openflow_controller_test OPENFLOW_TEST)
target_link_libraries(imsi_encoder_test OPENFLOW_TEST)
target_link_libraries(gtp_app_test OPENFLOW_TEST)

add_test(test_openflow_controller openflow_controller_test)
add_test(test_imsi_encoder imsi_encoder_test)
add_test(test_gtp_app gtp_app_test)

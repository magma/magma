add_compile_options(-std=c++14)

set(OUTPUT_DIR "${CMAKE_CURRENT_BINARY_DIR}")

include_directories("${PROJECT_SOURCE_DIR}")

add_library(SESSIOND_TEST_LIB
    Consts.h
    Matchers.h
    ProtobufCreators.cpp
    ProtobufCreators.h
    SessiondMocks.h
    SessionStateTester.h
    )

target_link_libraries(SESSIOND_TEST_LIB SESSION_MANAGER gmock_main pthread rt)

foreach (session_test session_credit cloud_reporter
    sessiond_integ 
    session_store store_client stored_state 
    metering_reporter charging_grant
    usage_monitor)
  add_executable(${session_test}_test test_${session_test}.cpp)
  target_link_libraries(${session_test}_test SESSIOND_TEST_LIB)
  add_test(test_${session_test} ${session_test}_test)
endforeach (session_test)

include_directories("${PROJECT_BINARY_DIR}")
include_directories("${PROJECT_SOURCE_DIR}/pipelined_client")

add_executable(pipelined_client_test test_pipelined_client.cpp)

target_link_libraries(pipelined_client_test
    LIB_PIPELINED_CLIENT gtest gtest_main pthread rt yaml-cpp
    )

add_test(test_pipelined_client pipelined_client_test)

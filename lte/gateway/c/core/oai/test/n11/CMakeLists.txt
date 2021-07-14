include_directories("${PROJECT_BINARY_DIR}")
include_directories("${PROJECT_SOURCE_DIR}/n11")

add_executable(n11_test test_n11.cpp)

target_link_libraries(n11_test
	LIB_N11 gtest gtest_main pthread protobuf rt yaml-cpp
    )

add_test(test_n11 n11_test)

include_directories("${PROJECT_BINARY_DIR}")
include_directories("${PROJECT_SOURCE_DIR}/n11")

add_executable(smf_service_client test_smf_service_client.cpp)
add_executable(auth_service_client test_auth_service_client.cpp)

target_link_libraries(smf_service_client
	LIB_N11 gtest gtest_main pthread protobuf rt yaml-cpp
    )
target_link_libraries(auth_service_client
	LIB_N11 TASK_AMF_APP gtest gtest_main pthread protobuf grpc++ dl m rt yaml-cpp
    )

add_test(test_smf_service_client smf_service_client)
add_test(test_auth_service_client auth_service_client)

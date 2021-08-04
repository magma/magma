include_directories("${PROJECT_BINARY_DIR}")
include_directories("${PROJECT_SOURCE_DIR}/n11")
include_directories("/usr/src/googletest/googlemock/include/")

add_library(AMF_TEST_LIB
    Consts.h
    M5GAmfMock.h
    test_smf_service_client.cpp
    test_auth_service_client.cpp
    test_amf_client_servicer.cpp
    )

link_directories(/usr/src/googletest/googlemock/lib/)
target_link_libraries(AMF_TEST_LIB
	TASK_AMF_APP LIB_N11 gtest gtest_main gmock_main gmock pthread protobuf grpc++ dl m rt yaml-cpp
    )

foreach (amf_client_test smf_service_client auth_service_client
         amf_client_servicer)
  add_executable(${amf_client_test}_test test_${amf_client_test}.cpp)
  target_link_libraries(${amf_client_test}_test AMF_TEST_LIB)
  add_test(test_${amf_client_test} ${amf_client_test}_test)
endforeach (amf_client_test)

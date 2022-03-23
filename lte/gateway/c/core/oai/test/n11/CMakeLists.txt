include_directories("${PROJECT_BINARY_DIR}")
include_directories("${PROJECT_SOURCE_DIR}/n11")

add_executable(smf_service_client test_smf_service_client.cpp)
add_executable(auth_service_client test_auth_service_client.cpp)

file(GLOB ABSL_LIBS "/usr/local/lib/libabsl*.a")

target_link_libraries(smf_service_client
    LIB_N11 gtest gtest_main pthread rt yaml-cpp
    /usr/local/lib/libprotobuf.a
    )
target_link_libraries(auth_service_client
    LIB_N11 TASK_AMF_APP -Wl,--start-group gtest gtest_main pthread  grpc++ grpc dl m rt yaml-cpp
    ${ABSL_LIBS}
    ${ABSL_LIBS}
    absl_debugging_internal
    absl_demangle_internal
    z
    gpr
    upb
    re2
    cares
    ssl
    crypto
    address_sorting
    /usr/local/lib/libprotobuf.a
    -Wl,--end-group
    )

add_test(test_smf_service_client smf_service_client)
add_test(test_auth_service_client auth_service_client)

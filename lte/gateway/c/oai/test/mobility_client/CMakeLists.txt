include_directories("${PROJECT_BINARY_DIR}")
include_directories("${PROJECT_SOURCE_DIR}/mobility_client")

add_executable(mobility_client_test test_mobility_client.cpp)

target_link_libraries(mobility_client_test
    LIB_MOBILITY_CLIENT protobuf grpc++ dl stdc++ m
    )

# TODO add support for integration tests
# add_test(test_rpc_client_integration rpc_client_test)

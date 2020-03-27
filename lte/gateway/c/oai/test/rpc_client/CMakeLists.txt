include_directories("${PROJECT_BINARY_DIR}")
include_directories("${PROJECT_SOURCE_DIR}/rpc_client")

add_executable(rpc_client_test test_rpc_client.cpp)

target_link_libraries(rpc_client_test
    LIB_RPC_CLIENT protobuf grpc++ dl stdc++ m
    )

# TODO add support for integration tests
# add_test(test_rpc_client_integration rpc_client_test)

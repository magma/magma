include_directories("${PROJECT_BINARY_DIR}")
include_directories("${PROJECT_SOURCE_DIR}/service303")

find_package(Protobuf REQUIRED)
add_compile_options(-std=c++11)

# compile the needed protos
set(S303_ORC8R_CPP_PROTOS service303 common metricsd)
set(S303_ORC8R_GRPC_PROTOS service303)

list(APPEND PROTO_SRCS "")
list(APPEND PROTO_HDRS "")

create_proto_dir("orc8r" ORC8R_OUT_DIR)

generate_cpp_protos("${S303_ORC8R_CPP_PROTOS}" "${PROTO_SRCS}"
  "${PROTO_HDRS}" ${ORC8R_PROTO_DIR} ${ORC8R_OUT_DIR})
generate_grpc_protos("${S303_ORC8R_GRPC_PROTOS}" "${PROTO_SRCS}"
  "${PROTO_HDRS}" ${ORC8R_PROTO_DIR} ${ORC8R_OUT_DIR})


add_library(SERVICE303_TEST
    Service303Client.h
    Service303Client.cpp
    )

add_executable(service303_test test_service303.cpp)
add_executable(metrics_test test_metrics.cpp)

target_link_libraries(service303_test
    SERVICE303_LIB
    SERVICE303_TEST
    LIB_BSTR
    TASK_SERVICE303
    SERVICE_REGISTRY
    CONFIG
    gtest gtest_main pthread grpc
    grpc++ prometheus-cpp yaml-cpp
)

target_link_libraries(metrics_test
    SERVICE303_LIB
    SERVICE303_TEST
    LIB_BSTR
    TASK_SERVICE303
    gtest gtest_main pthread prometheus-cpp
)

add_test(test_metrics metrics_test)
add_test(test_service303_integration service303_test)

cmake_minimum_required(VERSION 3.7)
set(CMAKE_EXPORT_COMPILE_COMMANDS ON)

# generate sctpd protos
include($ENV{MAGMA_ROOT}/orc8r/gateway/c/common/CMakeProtoMacros.txt)

create_proto_dir("lte" LTE_CPP_OUT_DIR)

list(APPEND PROTO_SRCS "")
list(APPEND PROTO_HDRS "")

set(SCTPD_LTE_CPP_PROTOS sctpd)
generate_cpp_protos("${SCTPD_LTE_CPP_PROTOS}" "${PROTO_SRCS}"
  "${PROTO_HDRS}" ${LTE_PROTO_DIR} ${LTE_CPP_OUT_DIR})

set(SCTPD_GRPC_PROTOS sctpd)
generate_grpc_protos("${SCTPD_GRPC_PROTOS}" "${PROTO_SRCS}"
  "${PROTO_HDRS}" ${LTE_PROTO_DIR} ${LTE_CPP_OUT_DIR})

# find needed libraries
set(MAGMA_LIB_DIR $ENV{MAGMA_ROOT}/orc8r/gateway/c/common)
add_subdirectory(${MAGMA_LIB_DIR}/logging /tmp)


# add sctpd executable
add_executable(sctpd
    sctp_assoc.cpp
    sctp_connection.cpp
    sctp_desc.cpp
    sctpd_downlink_impl.cpp
    sctpd_event_handler.cpp
    sctpd_uplink_client.cpp
    sctpd.cpp
    util.cpp
    ${PROTO_SRCS}
    ${PROTO_HDRS}
)

target_compile_definitions(sctpd PUBLIC LOG_WITH_GLOG)

target_link_libraries(sctpd
    sctp pthread grpc++ grpc protobuf glog LOGGING
)

target_include_directories(sctpd PUBLIC
    ${PROJECT_BINARY_DIR}
)

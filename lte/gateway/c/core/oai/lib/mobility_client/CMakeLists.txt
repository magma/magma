add_compile_options(-std=c++14)

include($ENV{MAGMA_ROOT}/orc8r/gateway/c/common/CMakeProtoMacros.txt)

set(OUTPUT_DIR "${CMAKE_CURRENT_BINARY_DIR}")

create_proto_dir("lte" LTE_CPP_OUT_DIR)
create_proto_dir("orc8r" ORC8R_CPP_OUT_DIR)


list(APPEND PROTO_SRCS "")
list(APPEND PROTO_HDRS "")

# compile the needed protos
set(RPC_LTE_CPP_PROTOS mobilityd apn subscriberdb)
generate_cpp_protos("${RPC_LTE_CPP_PROTOS}" "${PROTO_SRCS}"
    "${PROTO_HDRS}" ${LTE_PROTO_DIR} ${LTE_CPP_OUT_DIR})

set(RPC_ORC8R_CPP_PROTOS common)
generate_cpp_protos("${RPC_ORC8R_CPP_PROTOS}" "${PROTO_SRCS}"
    "${PROTO_HDRS}" ${ORC8R_PROTO_DIR} ${ORC8R_CPP_OUT_DIR})

set(RPC_LTE_GRPC_PROTOS mobilityd apn subscriberdb)
generate_grpc_protos("${RPC_LTE_GRPC_PROTOS}" "${PROTO_SRCS}" "${PROTO_HDRS}" ${LTE_PROTO_DIR} ${LTE_CPP_OUT_DIR})

set(RPC_ORC8R_GRPC_PROTOS "")
generate_grpc_protos("${RPC_ORC8R_GRPC_PROTOS}" "${PROTO_SRCS}" "${PROTO_HDRS}" ${ORC8R_PROTO_DIR} ${ORC8R_CPP_OUT_DIR})

add_library(LIB_MOBILITY_CLIENT
    MobilityServiceClient.cpp
    MobilityClientAPI.cpp
    ${PROTO_SRCS}
    ${PROTO_HDRS}
    )

target_link_libraries(LIB_MOBILITY_CLIENT
    COMMON
    LIB_BSTR
    LIB_PCEF
    ASYNC_GRPC SERVICE_REGISTRY
    )

target_include_directories(LIB_MOBILITY_CLIENT PUBLIC
    ${CMAKE_CURRENT_SOURCE_DIR}
    COMMON SERVICE_REGISTRY_LIBRARIES
    $<TARGET_FILE_DIR:LIB_MOBILITY_CLIENT>
    )
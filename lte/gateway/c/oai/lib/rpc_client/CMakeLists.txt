add_compile_options(-std=c++11)

# compile the needed protos
set(RPC_ORC8R_CPP_PROTOS common)
set(RPC_LTE_CPP_PROTOS mobilityd subscriberdb)
set(RPC_ORC8R_GRPC_PROTOS "")
set(RPC_LTE_GRPC_PROTOS mobilityd subscriberdb)

list(APPEND PROTO_SRCS "")
list(APPEND PROTO_HDRS "")

generate_all_protos("${RPC_LTE_CPP_PROTOS}" "${RPC_ORC8R_CPP_PROTOS}"
  "${RPC_LTE_GRPC_PROTOS}" "${RPC_ORC8R_GRPC_PROTOS}" "${PROTO_SRCS}" "${PROTO_HDRS}")

include_directories(${CMAKE_CURRENT_BINARY_DIR})

add_library(LIB_RPC_CLIENT
    MobilityClient.cpp
    RpcClient.cpp
    ${PROTO_SRCS}
    ${PROTO_HDRS}
    )

target_include_directories(LIB_RPC_CLIENT PUBLIC
    ${CMAKE_CURRENT_SOURCE_DIR}
)

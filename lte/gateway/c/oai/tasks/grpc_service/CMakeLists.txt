find_package(Protobuf REQUIRED)
set (CMAKE_CXX_FLAGS "-Wno-write-strings -Wno-literal-suffix")

# compile the needed protos
list(APPEND PROTO_SRCS "")
list(APPEND PROTO_HDRS "")


create_proto_dir("lte" LTE_OUT_DIR)
create_proto_dir("orc8r" ORC8R_OUT_DIR)
create_proto_dir("feg" FEG_OUT_DIR)

# Common protos
set(ORC8R_CPP_PROTOS common)
generate_cpp_protos("${ORC8R_CPP_PROTOS}" "${PROTO_SRCS}"
  "${PROTO_HDRS}" ${ORC8R_PROTO_DIR} ${ORC8R_OUT_DIR})

# SPGW
set(SPGWSRV_LTE_CPP_PROTOS spgw_service subscriberdb policydb)
set(SPGWSRV_LTE_GRPC_PROTOS spgw_service)
generate_cpp_protos("${SPGWSRV_LTE_CPP_PROTOS}" "${PROTO_SRCS}"
  "${PROTO_HDRS}" ${LTE_PROTO_DIR} ${LTE_OUT_DIR})
generate_grpc_protos("${SPGWSRV_LTE_GRPC_PROTOS}" "${PROTO_SRCS}"
  "${PROTO_HDRS}" ${LTE_PROTO_DIR} ${LTE_OUT_DIR})

# S6a
set(S6ASRV_FEG_CPP_PROTOS s6a_proxy)
set(S6ASRV_FEG_GRPC_PROTOS s6a_proxy)
set(S6ASRV_LTE_CPP_PROTOS s6a_service)
set(S6ASRV_LTE_GRPC_PROTOS s6a_service)
generate_cpp_protos("${S6ASRV_FEG_CPP_PROTOS}" "${PROTO_SRCS}"
  "${PROTO_HDRS}" ${FEG_PROTO_DIR} ${FEG_OUT_DIR})
generate_grpc_protos("${S6ASRV_FEG_GRPC_PROTOS}" "${PROTO_SRCS}"
  "${PROTO_HDRS}" ${FEG_PROTO_DIR} ${FEG_OUT_DIR})
generate_cpp_protos("${S6ASRV_LTE_CPP_PROTOS}" "${PROTO_SRCS}"
  "${PROTO_HDRS}" ${LTE_PROTO_DIR} ${LTE_OUT_DIR})
generate_grpc_protos("${S6ASRV_LTE_GRPC_PROTOS}" "${PROTO_SRCS}"
  "${PROTO_HDRS}" ${LTE_PROTO_DIR} ${LTE_OUT_DIR})

# SGs
set(SGSSRV_FEG_CPP_PROTOS csfb)
set(SGSSRV_FEG_GRPC_PROTOS csfb)
generate_cpp_protos("${SGSSRV_FEG_CPP_PROTOS}" "${PROTO_SRCS}"
  "${PROTO_HDRS}" ${FEG_PROTO_DIR} ${FEG_OUT_DIR})
generate_grpc_protos("${SGSSRV_FEG_GRPC_PROTOS}" "${PROTO_SRCS}"
  "${PROTO_HDRS}" ${FEG_PROTO_DIR} ${FEG_OUT_DIR})

include_directories(${CMAKE_CURRENT_BINARY_DIR})

add_library(TASK_GRPC_SERVICE
    SpgwServiceImpl.cpp
    S6aGatewayImpl.cpp
    S6aServiceImpl.cpp
    CSFBGatewayServiceImpl.cpp
    spgw_service_handler.c
    proto_msg_to_itti_msg.cpp
    grpc_service.cpp
    ${PROTO_SRCS}
    ${PROTO_HDRS}
    )

target_link_libraries(TASK_GRPC_SERVICE
    COMMON
    LIB_BSTR protobuf
    LIB_HASHTABLE
    ${PROTOBUF_LIBRARIES}
    grpc++
    folly
    )
target_include_directories(TASK_GRPC_SERVICE PUBLIC
    ${CMAKE_CURRENT_SOURCE_DIR}
    )

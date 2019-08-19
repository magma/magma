find_package(Protobuf REQUIRED)
add_compile_options(-std=c++14)

# compile the needed protos
set(SGSSRV_FEG_CPP_PROTOS csfb)
set(SGSSRV_FEG_GRPC_PROTOS csfb)
set(SGSSRV_ORC8R_CPP_PROTOS common)
set(SGSSRV_ORC8R_GRPC_PROTOS common)

list(APPEND PROTO_SRCS "")
list(APPEND PROTO_HDRS "")

create_proto_dir("feg" FEG_OUT_DIR)
create_proto_dir("orc8r" ORC8R_OUT_DIR)

generate_cpp_protos("${SGSSRV_FEG_CPP_PROTOS}" "${PROTO_SRCS}"
  "${PROTO_HDRS}" ${FEG_PROTO_DIR} ${FEG_OUT_DIR})
generate_grpc_protos("${SGSSRV_FEG_GRPC_PROTOS}" "${PROTO_SRCS}"
  "${PROTO_HDRS}" ${FEG_PROTO_DIR} ${FEG_OUT_DIR})

generate_cpp_protos("${SGSSRV_ORC8R_CPP_PROTOS}" "${PROTO_SRCS}"
  "${PROTO_HDRS}" ${ORC8R_PROTO_DIR} ${ORC8R_OUT_DIR})

include_directories(${CMAKE_CURRENT_BINARY_DIR})

add_library(TASK_SGS_SERVICE
    CSFBGatewayServiceImpl.cpp
    proto_msg_to_itti_msg.cpp
    sgs_service.cpp
    "${PROTO_SRCS}"
    "${PROTO_HDRS}"
    )

target_link_libraries(TASK_SGS_SERVICE
    COMMON
    LIB_BSTR LIB_HASHTABLE
    ${PROTOBUF_LIBRARIES}
    grpc++
)
target_include_directories(TASK_SGS_SERVICE PUBLIC
    ${CMAKE_CURRENT_SOURCE_DIR}
    )

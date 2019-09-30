add_compile_options(-std=c++11)
set (CMAKE_CXX_FLAGS "-Wno-write-strings -Wno-literal-suffix")

# compile the needed protos
list(APPEND PROTO_SRCS "")
list(APPEND PROTO_HDRS "")

create_proto_dir("lte" LTE_OUT_DIR)
create_proto_dir("orc8r" ORC8R_CPP_OUT_DIR)

set(SPGWSRV_ORC8R_CPP_PROTOS common)
generate_cpp_protos("${SPGWSRV_ORC8R_CPP_PROTOS}" "${PROTO_SRCS}"
  "${PROTO_HDRS}" ${ORC8R_PROTO_DIR} ${ORC8R_CPP_OUT_DIR})

set(SPGWSRV_LTE_CPP_PROTOS spgw_service subscriberdb policydb)
generate_cpp_protos("${SPGWSRV_LTE_CPP_PROTOS}" "${PROTO_SRCS}"
  "${PROTO_HDRS}" ${LTE_PROTO_DIR} ${LTE_OUT_DIR})

set(SPGWSRV_LTE_GRPC_PROTOS spgw_service)
generate_grpc_protos("${SPGWSRV_LTE_GRPC_PROTOS}" "${PROTO_SRCS}"
  "${PROTO_HDRS}" ${LTE_PROTO_DIR} ${LTE_OUT_DIR})

include_directories(${CMAKE_CURRENT_BINARY_DIR})

add_library(TASK_SPGW_SERVICE
    SpgwServiceImpl.cpp
    spgw_service.cpp
    ${PROTO_SRCS}
    ${PROTO_HDRS}
    )

target_link_libraries(TASK_SPGW_SERVICE
    COMMON
    LIB_BSTR protobuf
)
target_include_directories(TASK_SPGW_SERVICE PUBLIC
    ${CMAKE_CURRENT_SOURCE_DIR}
    )

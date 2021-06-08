add_compile_options(-std=c++14)

set(CMAKE_CXX_FLAGS "-Wno-write-strings -Wno-literal-suffix")
add_definitions(-DDEBUG_IS_ON=1)

# compile the needed protos
set(S6A_FEG_CPP_PROTOS s6a_proxy)
set(S6A_FEG_GRPC_PROTOS s6a_proxy)
set(S6A_LTE_CPP_PROTOS mconfig/mconfigs)

list(APPEND PROTO_SRCS "")
list(APPEND PROTO_HDRS "")

create_proto_dir("feg" FEG_OUT_DIR)
create_proto_dir("lte" LTE_OUT_DIR)

generate_cpp_protos("${S6A_FEG_CPP_PROTOS}" "${PROTO_SRCS}"
    "${PROTO_HDRS}" ${FEG_PROTO_DIR} ${FEG_OUT_DIR})
generate_grpc_protos("${S6A_FEG_GRPC_PROTOS}" "${PROTO_SRCS}"
    "${PROTO_HDRS}" ${FEG_PROTO_DIR} ${FEG_OUT_DIR})

generate_cpp_protos("${S6A_LTE_CPP_PROTOS}" "${PROTO_SRCS}"
    "${PROTO_HDRS}" ${LTE_PROTO_DIR} ${LTE_OUT_DIR})

add_library(LIB_S6A_PROXY
    S6aClient.cpp
    s6a_client_api.cpp
    itti_msg_to_proto_msg.cpp
    proto_msg_to_itti_msg.cpp
    ${PROTO_SRCS}
    ${PROTO_HDRS}
    )

target_link_libraries(LIB_S6A_PROXY
    COMMON
    LIB_MOBILITY_CLIENT ASYNC_GRPC SERVICE_REGISTRY
    MAGMA_CONFIG LIB_BSTR LIB_HASHTABLE
    )
target_include_directories(LIB_S6A_PROXY PUBLIC
    ${CMAKE_CURRENT_SOURCE_DIR}
    $<TARGET_FILE_DIR:LIB_S6A_PROXY>
    )
find_package(Protobuf REQUIRED)
set(CMAKE_CXX_FLAGS "-Wno-write-strings -Wno-literal-suffix")

# compile the needed protos
list(APPEND PROTO_SRCS "")
list(APPEND PROTO_HDRS "")

create_proto_dir("lte" LTE_OUT_DIR)
create_proto_dir("orc8r" ORC8R_OUT_DIR)
create_proto_dir("feg" FEG_OUT_DIR)

# Common protos
set(ORC8R_CPP_PROTOS common digest)
generate_cpp_protos("${ORC8R_CPP_PROTOS}" "${PROTO_SRCS}"
    "${PROTO_HDRS}" ${ORC8R_PROTO_DIR} ${ORC8R_OUT_DIR})

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

add_library(TASK_ASYNC_GRPC_SERVICE
    grpc_async_service_task.cpp
    ${PROTO_SRCS}
    ${PROTO_HDRS}
    )

file(GLOB ABSL_LIBS "/usr/local/lib/libabsl*.a")

target_link_libraries(TASK_ASYNC_GRPC_SERVICE
    -Wl,--start-group
    COMMON
    LIB_BSTR /usr/local/lib/libprotobuf.a
    LIB_HASHTABLE
    ${PROTOBUF_LIBRARIES}
    pthread
    gpr
    grpc++
    grpc
    TASK_S6A
    ${ABSL_LIBS}
    ${ABSL_LIBS}
    absl_debugging_internal
    absl_demangle_internal
    z
    gpr
    upb
    re2
    cares
    ssl
    crypto
    address_sorting
    /usr/local/lib/libprotobuf.a
	-Wl,--end-group
    )
target_include_directories(TASK_ASYNC_GRPC_SERVICE PUBLIC
    ${CMAKE_CURRENT_SOURCE_DIR}
    $<TARGET_FILE_DIR:TASK_ASYNC_GRPC_SERVICE>
    )

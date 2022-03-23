
set(S1AP_C_DIR ${PROJECT_BINARY_DIR}/s1ap/r15)
include_directories(${S1AP_C_DIR})

list(APPEND PROTO_SRCS "")
list(APPEND PROTO_HDRS "")

set(SCTPD_LTE_CPP_PROTOS sctpd)
set(SCTPD_GRPC_PROTOS sctpd)

generate_all_protos(
    "${SCTPD_LTE_CPP_PROTOS}" "" "${SCTPD_GRPC_PROTOS}" ""
    "${PROTO_SRCS}" "${PROTO_HDRS}"
)

add_library(TASK_SCTP_SERVER
    sctpd_downlink_client.cpp
    sctpd_uplink_server.cpp
    sctp_itti_messaging.c
    sctp_primitives_server.c
    ${PROTO_SRCS}
    ${PROTO_HDRS}
    )

file(GLOB ABSL_LIBS "/usr/local/lib/libabsl*.a")

target_link_libraries(TASK_SCTP_SERVER
    -Wl,--start-group
    COMMON
    LIB_BSTR LIB_HASHTABLE
    grpc++ grpc
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
target_include_directories(TASK_SCTP_SERVER PUBLIC
    ${CMAKE_CURRENT_SOURCE_DIR}
    ${CMAKE_CURRENT_BINARY_DIR}
    ${PROJECT_BINARY_DIR}
    )

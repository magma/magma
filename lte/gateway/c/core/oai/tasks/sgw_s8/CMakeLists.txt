# compile required protos
set(SGW_LTE_CPP_PROTOS spgw_state sgw_state common_types std_3gpp_types)
list(APPEND PROTO_SRCS "")
list(APPEND PROTO_HDRS "")

set(STATE_OUT_DIR "${CMAKE_CURRENT_BINARY_DIR}/lte/protos/oai")
generate_cpp_protos("${SGW_LTE_CPP_PROTOS}" "${PROTO_SRCS}" "${PROTO_HDRS}"
    "${STATE_PROTO_DIR}" "${STATE_OUT_DIR}")

message("s8 Proto_srcs are ${PROTO_SRCS}")
message("s8 Proto_includes are ${PROTO_HDRS}")
set(SGW_S8_SRC
    sgw_s8_task.c
    sgw_s8_handlers.c
    sgw_s8_state.cpp
    sgw_s8_state_manager.cpp
    sgw_s8_state_converter.cpp
    ${PROTO_SRCS}
    ${PROTO_HDRS}
    )
add_library(TASK_SGW_S8 ${SGW_S8_SRC})

target_link_libraries(TASK_SGW_S8
    COMMON
    LIB_BSTR LIB_HASHTABLE LIB_S8_PROXY
    LIB_OPENFLOW_CONTROLLER TASK_SERVICE303
    cpp_redis tacopie protobuf
    )

target_include_directories(TASK_SGW_S8 PUBLIC
    ${CMAKE_CURRENT_SOURCE_DIR}
    )

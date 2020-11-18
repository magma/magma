find_package(Protobuf REQUIRED)
# compile the needed protos
list(APPEND PROTO_SRCS "")
list(APPEND PROTO_HDRS "")
# compile MME APP STATE PROTOS
set(MMEAPP_STATE_CPP_PROTOS mme_nas_state nas_state spgw_state common_types
  std_3gpp_types)

set(LTE_OUT_MME_STATE_DIR "${CMAKE_CURRENT_BINARY_DIR}/lte/protos/oai")
generate_cpp_protos("${MMEAPP_STATE_CPP_PROTOS}" "${PROTO_SRCS}"
  "${PROTO_HDRS}" ${STATE_PROTO_DIR} ${LTE_OUT_MME_STATE_DIR})

include_directories(${CMAKE_CURRENT_BINARY_DIR})
include_directories(${LTE_OUT_MME_STATE_DIR})

add_library(TASK_HA
    ha_task.c
    ha_service_handler.cpp
    ${PROTO_SRCS}
    ${PROTO_HDRS}
    )

target_link_libraries(TASK_HA
    COMMON
    LIB_BSTR LIB_HASHTABLE protobuf
    ${PROTOBUF_LIBRARIES}
    TASK_MME_APP
    TASK_S1AP
)

target_include_directories(TASK_HA PUBLIC
    ${CMAKE_CURRENT_SOURCE_DIR}
)

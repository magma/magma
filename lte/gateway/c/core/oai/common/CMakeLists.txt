if (${LINK_GCOV})
  set(GCOV_LIB gcov)
else ()
  set(GCOV_LIB "")
endif ()

include($ENV{MAGMA_ROOT}/orc8r/gateway/c/common/CMakeProtoMacros.txt)

# compile the needed protos
set(COMMON_CPP_PROTOS common_types)
set(REDIS_CPP_PROTOS redis)
set(MCONFIG_PROTOS mconfig/mconfigs)
list(APPEND PROTO_SRCS "")
list(APPEND PROTO_HDRS "")

create_proto_dir("orc8r" ORC8R_CPP_OUT_DIR)
set(LTE_OUT_DIR "${CMAKE_CURRENT_BINARY_DIR}/lte/protos/oai")
generate_cpp_protos("${COMMON_CPP_PROTOS}" "${PROTO_SRCS}"
    "${PROTO_HDRS}" ${STATE_PROTO_DIR} ${LTE_OUT_DIR})

create_proto_dir("orc8r" ORC8R_CPP_OUT_DIR)
set(RPC_ORC8R_CPP_PROTOS common)
generate_cpp_protos("${RPC_ORC8R_CPP_PROTOS}" "${PROTO_SRCS}"
    "${PROTO_HDRS}" ${ORC8R_PROTO_DIR} ${ORC8R_CPP_OUT_DIR})
generate_cpp_protos("${REDIS_CPP_PROTOS}" "${PROTO_SRCS}"
    "${PROTO_HDRS}" ${ORC8R_PROTO_DIR} ${ORC8R_CPP_OUT_DIR})

set(MCONFIG_OUT_DIR "${CMAKE_CURRENT_BINARY_DIR}/lte/protos")
generate_cpp_protos("${MCONFIG_PROTOS}" "${PROTO_SRCS}"
    "${PROTO_HDRS}" ${LTE_PROTO_DIR} ${MCONFIG_OUT_DIR})

message("Proto_srcs are ${PROTO_SRCS}")

###############################################################################
# Core Network Utils
###############################################################################

set(S1AP_C_DIR ${PROJECT_BINARY_DIR}/s1ap/r10.5)
include_directories("${S1AP_C_DIR}")

set(COMMON_SRC
    TLVDecoder.c
    TLVEncoder.c
    async_system.c
    backtrace.c
    common_types.c
    conversions.c
    digest.c
    dynamic_memory_check.c
    enum_string.c
    itti_free_defined_msg.c
    pid_file.c
    shared_ts_log.c
    log.c
    state_converter.cpp
    common_utility_funs.cpp
    sentry_wrapper.cpp
    ${PROTO_SRCS}
    ${PROTO_HDRS}
    )

add_subdirectory(glogwrapper)
add_subdirectory(redis_utils)

add_library(COMMON ${COMMON_SRC})
target_link_libraries(COMMON
    lfds710
    LIB_3GPP LIB_BSTR LIB_HASHTABLE LIB_ITTI glogwrapper
    )

## Find + link sentry if it exists
## TODO we can enable by default after we migrate to ubuntu 20.04
find_library(SENTRY_LIB NAMES sentry sentry-native)
if (SENTRY_LIB)
  # Include dir
  find_path(SENTRY_INCLUDE_DIR
      NAMES
      sentry.h
      )
  target_link_libraries(COMMON ${SENTRY_LIB} MAGMA_CONFIG)
  target_include_directories(COMMON PUBLIC
      ${CMAKE_CURRENT_SOURCE_DIR}
      $<TARGET_FILE_DIR:COMMON>
      )
endif ()

target_include_directories(COMMON PUBLIC
    ${CMAKE_CURRENT_SOURCE_DIR}
    $<TARGET_FILE_DIR:COMMON>
    )

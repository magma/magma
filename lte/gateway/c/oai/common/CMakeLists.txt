if (${LINK_GCOV})
  set(GCOV_LIB gcov)
else ()
  set(GCOV_LIB "")
endif ()

# compile the needed protos
set(COMMON_CPP_PROTOS  common_types)
list(APPEND PROTO_SRCS "")
list(APPEND PROTO_HDRS "")

set(LTE_OUT_DIR "${CMAKE_CURRENT_BINARY_DIR}/lte/gateway/c/oai/protos")
generate_cpp_protos("${COMMON_CPP_PROTOS}" "${PROTO_SRCS}"
  "${PROTO_HDRS}" ${STATE_PROTO_DIR} ${LTE_OUT_DIR})

message("Proto_srcs are ${PROTO_SRCS}")
include_directories(${CMAKE_CURRENT_BINARY_DIR})
include_directories(${LTE_OUT_DIR})


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
  mcc_mnc_itu.c
  pid_file.c
  shared_ts_log.c
  log.c
  state_converter.cpp
  ${PROTO_SRCS}
  ${PROTO_HDRS}
)

add_subdirectory(glogwrapper)
add_subdirectory(redis_utils)

add_library(COMMON ${COMMON_SRC})
target_link_libraries (COMMON
  lfds710
  LIB_3GPP LIB_BSTR LIB_HASHTABLE LIB_ITTI glogwrapper redis_utils
)
target_include_directories(COMMON PUBLIC
  ${CMAKE_CURRENT_SOURCE_DIR}
)

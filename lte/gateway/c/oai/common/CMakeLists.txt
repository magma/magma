if (${LINK_GCOV})
  set(GCOV_LIB gcov)
else ()
  set(GCOV_LIB "")
endif ()

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
)

add_library(COMMON ${COMMON_SRC})
target_link_libraries (COMMON
  lfds710
  LIB_3GPP LIB_BSTR LIB_HASHTABLE LIB_ITTI
)
target_include_directories(COMMON PUBLIC
  ${CMAKE_CURRENT_SOURCE_DIR}
)

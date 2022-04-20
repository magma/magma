if (${LINK_GCOV})
  set(GCOV_LIB gcov)
else ()
  set(GCOV_LIB "")
endif ()

include($ENV{MAGMA_ROOT}/orc8r/gateway/c/common/CMakeProtoMacros.txt)

pkg_search_module(OPENSSL openssl REQUIRED)
include_directories(${OPENSSL_INCLUDE_DIRS})

pkg_search_module(CRYPTO libcrypto REQUIRED)
include_directories(${CRYPTO_INCLUDE_DIRS})

pkg_search_module(NETTLE nettle REQUIRED)
include_directories(${NETTLE_INCLUDE_DIRS})

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
    conversions.c
    digest.c
    enum_string.c
    itti_free_defined_msg.c
    pid_file.c
    shared_ts_log.c
    log.c
    sentry_log.cpp
    state_converter.cpp
    common_utility_funs.cpp
    ${PROTO_SRCS}
    ${PROTO_HDRS}
    )

add_subdirectory(glogwrapper)
add_subdirectory(redis_utils)

add_library(COMMON ${COMMON_SRC})

if (EXISTS "/etc/os-release")
  set(os_file_path "/etc/os-release")
  file(STRINGS "${os_file_path}" data_list REGEX "^(ID|VERSION_ID|VERSION_CODENAME)=")
  foreach (_var ${data_list})
    if ("${_var}" MATCHES "^(ID)=(.*)$")
      set(os_name "${CMAKE_MATCH_2}")
    endif ()
  endforeach ()
  if ("${os_name}" STREQUAL "\"rhel\"")
  endif ()
endif ()

find_package(LTE_GATEWAY_C_COMMON REQUIRED)

message("os_name" ${os_name})
if ("${os_name}" STREQUAL "\"rhel\"")
  target_link_libraries(COMMON
      LTE_GATEWAY_C_COMMON
      lfds710
      LIB_3GPP LIB_BSTR LIB_HASHTABLE LIB_ITTI glogwrapper MAGMA_SENTRY SERVICE303_LIB
      )
else()
target_link_libraries(COMMON
      LTE_GATEWAY_C_COMMON
      lfds710
      ${OPENSSL_LIBRARIES}
      LIB_3GPP LIB_BSTR LIB_HASHTABLE LIB_ITTI glogwrapper MAGMA_SENTRY SERVICE303_LIB
      )
endif()

target_include_directories(COMMON PUBLIC
    $ENV{MAGMA_ROOT}
    ${CMAKE_CURRENT_SOURCE_DIR}
    $<TARGET_FILE_DIR:COMMON>
    )

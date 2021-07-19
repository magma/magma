# Override options for MME
set(BYTE_ORDER=LITTLE_ENDIAN)
set(S1AP_DEBUG_LIST False)
set(SCTP_DUMP_LIST False)
set(TRACE_HASHTABLE False)
set(TRACE_3GPP_SPEC False)


# Dependencies
pkg_search_module(LIBXML2 libxml-2.0 REQUIRED)
include_directories(${LIBXML2_INCLUDE_DIRS})

pkg_search_module(OPENSSL openssl REQUIRED)
include_directories(${OPENSSL_INCLUDE_DIRS})

pkg_search_module(CRYPTO libcrypto REQUIRED)
include_directories(${CRYPTO_INCLUDE_DIRS})

pkg_search_module(NETTLE nettle REQUIRED)
include_directories(${NETTLE_INCLUDE_DIRS})

if (NOT "${NETTLE_VERSION}" STREQUAL "2.5")
  message(FATAL_ERROR "Nettle version incorrect -- 2.5 required")
endif ()

find_library(LFDS lfds710 PATHS /usr/local/lib /usr/lib)

if ("${LFDS}" STREQUAL "LFDS-NOTFOUND")
  message(FATAL_ERROR "LIB LFDS not found, please install it")
endif ()

add_executable(mme
    ${PROJECT_SOURCE_DIR}/oai_mme/oai_mme.c
    ${PROJECT_SOURCE_DIR}/common/common_types.c
    ${PROJECT_SOURCE_DIR}/common/itti_free_defined_msg.c
    ${PROJECT_SOURCE_DIR}/tasks/service303/service303_task.c
    ${PROJECT_SOURCE_DIR}/tasks/service303/service303_mme_stats.c
    ${PROJECT_SOURCE_DIR}/tasks/grpc_service/grpc_service_task.c
    )

target_link_libraries(mme
    -Wl,--start-group
    COMMON
    LIB_3GPP LIB_S1AP LIB_NGAP LIB_SECU LIB_DIRECTORYD LIB_SGS_CLIENT LIB_BSTR
    LIB_HASHTABLE LIB_S6A_PROXY
    TASK_S1AP TASK_NGAP TASK_SCTP_SERVER TASK_SGS TASK_SMS_ORC8R
    TASK_S6A TASK_MME_APP TASK_GRPC_SERVICE TASK_NAS TASK_HA
    ${ITTI_LIB} ${GCOV_LIB}
    -Wl,--end-group
    ${LFDS} pthread m sctp rt crypt ${CRYPTO_LIBRARIES} ${OPENSSL_LIBRARIES}
    ${NETTLE_LIBRARIES} ${CONFIG_LIBRARIES} gnutls
    prometheus-cpp grpc grpc++ yaml-cpp
    )

if (NOT EMBEDDED_SGW)
  target_link_libraries(mme
      LIB_GTPV2C TASK_UDP)
else (EMBEDDED_SGW)
  target_link_libraries(mme TASK_SGW TASK_SGW_S8)
endif (NOT EMBEDDED_SGW)

if (NOT S6A_OVER_GRPC)
  target_link_libraries(mme fdproto fdcore)
endif (NOT S6A_OVER_GRPC)

if (NOT EMBEDDED_SGW AND NOT S6A_OVER_GRPC)
  if (EXISTS "/etc/os-release")
    set(os_file_path "/etc/os-release")
    file(STRINGS "${os_file_path}" data_list REGEX "^(ID|VERSION_ID|VERSION_CODENAME)=")
    foreach (_var ${data_list})
      if ("${_var}" MATCHES "^(ID)=(.*)$")
        set(os_name "${CMAKE_MATCH_2}")
      endif ()
    endforeach ()
    if ("${os_name}" STREQUAL "\"rhel\"")
      target_link_libraries(mme
          gpr yaml z cares boost_program_options boost_filesystem boost_system boost_regex double-conversion address_sorting gflags unwind iberty dl)
    endif ()
  endif ()
endif (NOT EMBEDDED_SGW AND NOT S6A_OVER_GRPC)

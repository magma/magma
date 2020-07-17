# Override options for MME
set (  BYTE_ORDER=LITTLE_ENDIAN )
set (  S1AP_DEBUG_LIST                 False )
set (  SCTP_DUMP_LIST                  False )
set (  TRACE_HASHTABLE                 False )
set (  TRACE_3GPP_SPEC                 False )

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

find_library(LFDS lfds710 PATHS /usr/local/lib /usr/lib )

if ("${LFDS}" STREQUAL "LFDS-NOTFOUND")
  message(FATAL_ERROR "LIB LFDS not found, please install it" )
endif()


add_executable(mme
    ${PROJECT_SOURCE_DIR}/oai_mme/oai_mme_log.c
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
        LIB_3GPP LIB_S1AP LIB_SECU LIB_DIRECTORYD LIB_SGS_CLIENT LIB_BSTR
        LIB_HASHTABLE LIB_S6A_PROXY
        TASK_S1AP TASK_SCTP_SERVER TASK_SGS
        TASK_S6A TASK_MME_APP TASK_GRPC_SERVICE TASK_NAS
        ${ITTI_LIB} ${GCOV_LIB}
    -Wl,--end-group
    ${LFDS} pthread m sctp  rt crypt ${CRYPTO_LIBRARIES} ${OPENSSL_LIBRARIES}
    ${NETTLE_LIBRARIES} ${CONFIG_LIBRARIES} gnutls fdproto fdcore ${SERVICE303_LIB}
    prometheus-cpp grpc grpc++
)

if ( NOT EMBEDDED_SGW )
    target_link_libraries(mme GTPV2C S11_MME)
elseif ( EMBEDDED_SGW )
    target_link_libraries(mme TASK_SGW)
endif ( NOT EMBEDDED_SGW )

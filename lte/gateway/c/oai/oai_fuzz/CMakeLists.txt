set (  BYTE_ORDER=LITTLE_ENDIAN )

pkg_search_module(LIBXML2 libxml-2.0 REQUIRED)
include_directories(${LIBXML2_INCLUDE_DIRS})

pkg_search_module(OPENSSL openssl REQUIRED)
include_directories(${OPENSSL_INCLUDE_DIRS})

pkg_search_module(CRYPTO libcrypto REQUIRED)
include_directories(${CRYPTO_INCLUDE_DIRS})

pkg_search_module(NETTLE nettle REQUIRED)
include_directories(${NETTLE_INCLUDE_DIRS})

find_library(LFDS lfds710 PATHS /usr/local/lib /usr/lib )

add_executable(oai_fuzz
    mock_task.c
    fuzz_s1ap.c
)

target_link_libraries(oai_fuzz
    -Wl,--start-group
        COMMON
        LIB_3GPP LIB_S1AP LIB_SECU LIB_DIRECTORYD LIB_SGS_CLIENT LIB_BSTR
        LIB_HASHTABLE LIB_S6A_PROXY
        TASK_S1AP TASK_SCTP_SERVER TASK_SGS
        TASK_S6A TASK_MME_APP TASK_GRPC_SERVICE
        TASK_NAS TASK_SGW
        ${MSC_LIB} ${ITTI_LIB} ${GCOV_LIB}
    -Wl,--end-group
    ${LFDS} pthread m sctp rt crypt ${CRYPTO_LIBRARIES} ${OPENSSL_LIBRARIES}
    ${NETTLE_LIBRARIES} ${CONFIG_LIBRARIES} gnutls fdproto fdcore
    ${SERVICE303_LIB} ${SERVICE_REGISTRY}
    prometheus-cpp grpc grpc++
)

target_include_directories(oai_fuzz PUBLIC
    ${CMAKE_CURRENT_SOURCE_DIR}
)

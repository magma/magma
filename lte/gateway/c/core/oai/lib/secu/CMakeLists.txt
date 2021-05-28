###############################################################################
# Security library
###############################################################################

add_library(LIB_SECU
    kdf.c
    key_nas_deriver.c
    key_nas_encryption.c
    nas_stream_eea1.c
    nas_stream_eea2.c
    nas_stream_eia1.c
    nas_stream_eia2.c
    rijndael.c
    snow3g.c
    )
target_link_libraries(LIB_SECU
    ${NETTLE_LIBRARIES}
    COMMON
    LIB_BSTR
    )
target_include_directories(LIB_SECU PUBLIC
    ${CMAKE_CURRENT_SOURCE_DIR}
    )

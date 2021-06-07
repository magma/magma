add_library(LIB_3GPP
    3gpp_24.008_cc_ies.c
    3gpp_24.008_common_ies.c
    3gpp_24.008_gmm_ies.c
    3gpp_24.008_gprs_common_ies.c
    3gpp_24.008_mm_ies.c
    3gpp_24.008_sm_ies.c
    )
target_link_libraries(LIB_3GPP
    LIB_BSTR LIB_S1AP LIB_NGAP
    )
target_include_directories(LIB_3GPP PUBLIC
    ${CMAKE_CURRENT_SOURCE_DIR}
    )

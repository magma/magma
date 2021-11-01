add_library(LIB_BSTR
        bstrlib.c
        )
target_include_directories(LIB_BSTR PUBLIC $ENV{MAGMA_ROOT} ${CMAKE_CURRENT_SOURCE_DIR})

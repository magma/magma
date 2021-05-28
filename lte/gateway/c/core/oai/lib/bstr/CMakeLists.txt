add_library(LIB_BSTR
    bstraux.c
    bstrlib.c
    buniutil.c
    utf8util.c
    )
target_include_directories(LIB_BSTR PUBLIC ${CMAKE_CURRENT_SOURCE_DIR})

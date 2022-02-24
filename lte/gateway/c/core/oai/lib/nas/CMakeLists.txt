add_compile_options(-std=c++14)

add_library(LIB_NAS 
	NasBuffer.cpp
        )
target_include_directories(LIB_NAS PUBLIC $ENV{MAGMA_ROOT} ${CMAKE_CURRENT_SOURCE_DIR})

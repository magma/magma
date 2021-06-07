add_library(LIB_HASHTABLE
    hashtable.c
    obj_hashtable.c
    hashtable_uint64.c
    obj_hashtable_uint64.c
    )
target_link_libraries(LIB_HASHTABLE
    COMMON
    LIB_BSTR
    )
target_include_directories(LIB_HASHTABLE PUBLIC
    ${CMAKE_CURRENT_SOURCE_DIR}
    )

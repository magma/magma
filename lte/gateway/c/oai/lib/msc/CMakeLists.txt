if (MESSAGE_CHART_GENERATOR)
    add_library(LIB_MSC
        msc.c
    )
    target_include_directories(LIB_MSC
        PUBLIC
        ${CMAKE_CURRENT_SOURCE_DIR}
        PRIVATE
        ${LIB_DIR}/itti
    )
    target_link_libraries(LIB_MSC
        LIB_BSTR LIB_HASHTABLE
    )
else ()
    add_library(LIB_MSC INTERFACE)
    target_include_directories(LIB_MSC INTERFACE ${CMAKE_CURRENT_SOURCE_DIR})
endif ()

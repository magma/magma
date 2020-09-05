find_package(Check REQUIRED)
find_package(Threads REQUIRED)

set(Reg_Req_Decode
    reg_req_decode.cpp
)

add_executable(reg_req_decode ${Reg_Req_Decode})
target_link_libraries(reg_req_decode
    TASK_NAS5G 
    #    ${CHECK_LIBRARIES} ${CMAKE_THREAD_LIBS_INIT}
    # LIB_BSTR LIB_HASHTABLE
)

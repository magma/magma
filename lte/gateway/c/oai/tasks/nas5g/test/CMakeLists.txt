find_package(Check REQUIRED)
find_package(Threads REQUIRED)

#Executable for Registration Request Message Decode

set(Reg_Req_Decode reg_req_decode.cpp)

add_executable(reg_req_decode ${Reg_Req_Decode})

target_link_libraries(reg_req_decode TASK_NAS5G)


#Executable for Registration Accept  Message encode

set(Reg_Acpt_Encode reg_acpt_encode.cpp)

add_executable(reg_acpt_encode ${Reg_Acpt_Encode})

target_link_libraries(reg_acpt_encode TASK_NAS5G)


#Executable for De-Registration Request(UE)  Message decode

set(De_Reg_Req_Decode de_reg_req_decode.cpp)

add_executable(de_reg_req_decode ${De_Reg_Req_Decode})

target_link_libraries(de_reg_req_decode TASK_NAS5G)


#Executable for De-Registration Accept(UE)  Message encode

set(De_Reg_Acpt_Encode de_reg_acpt_encode.cpp)

add_executable(de_reg_acpt_encode ${De_Reg_Acpt_Encode})

target_link_libraries(de_reg_acpt_encode TASK_NAS5G)


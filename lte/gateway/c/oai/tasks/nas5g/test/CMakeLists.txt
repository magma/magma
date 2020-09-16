find_package(Check REQUIRED)
find_package(Threads REQUIRED)

#Executable for Registration Accept  Message encode

add_executable(RegAcptEncode RegAcptEncode.cpp)

target_link_libraries(RegAcptEncode TASK_NAS5G)


#Executable for De-Registration Request(UE)  Message decode

add_executable(DeRegReqDecode DeRegReqDecode.cpp)

target_link_libraries(DeRegReqDecode TASK_NAS5G)


#Executable for De-Registration Accept(UE)  Message encode

add_executable(DeRegAcptEncode DeRegAcptEncode.cpp)

target_link_libraries(DeRegAcptEncode TASK_NAS5G)

#Executable for Authentication Request Message encode

add_executable(AuthReqEncode AuthReqEncode.cpp)

target_link_libraries(AuthReqEncode TASK_NAS5G)


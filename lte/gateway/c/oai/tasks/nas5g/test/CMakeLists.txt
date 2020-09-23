find_package(Check REQUIRED)
find_package(Threads REQUIRED)

#Executable for Registration Request  Message Decode
add_executable(RegReqDecode RegReqDecode.cpp)

target_link_libraries(RegReqDecode TASK_NAS5G)

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

#Executable for Authentication Response Message Decode
add_executable(AuthResDecode AuthResDecode.cpp)

target_link_libraries(AuthResDecode TASK_NAS5G)

#Executable for Security mode command encode
add_executable(SecModeCommandEncode SecModeCommandEncode.cpp)

target_link_libraries(SecModeCommandEncode TASK_NAS5G)

#Executable for Security mode complete decode
add_executable(SecurityModeCompleteDecode SecurityModeCompleteDecode.cpp)

target_link_libraries(SecurityModeCompleteDecode TASK_NAS5G)


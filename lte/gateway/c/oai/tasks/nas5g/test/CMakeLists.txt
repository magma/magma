find_package(Check REQUIRED) 
find_package(Threads REQUIRED)

#Executable for Registration Request Message Decode
add_executable(RegReqDecode RegReqDecode.cpp)

target_link_libraries(RegReqDecode LIB_NAS5G)

#Executable for Registration Accept Message encode
add_executable(RegAcptEncode RegAcptEncode.cpp)

target_link_libraries(RegAcptEncode LIB_NAS5G)

#Executable for Registration Complete Message Decode
add_executable(RegComDecode RegComDecode.cpp)

target_link_libraries(RegComDecode LIB_NAS5G)

#Executable for Registration Reject Message encode
add_executable(RegRejEncode RegRejEncode.cpp)

target_link_libraries(RegRejEncode LIB_NAS5G)

#Executable for De - Registration Request(UE) Message decode
add_executable(DeRegReqDecode DeRegReqDecode.cpp)

target_link_libraries(DeRegReqDecode LIB_NAS5G)

#Executable for De - Registration Accept(UE) Message encode
add_executable(DeRegAcptEncode DeRegAcptEncode.cpp)

target_link_libraries(DeRegAcptEncode LIB_NAS5G)

#Executable for Authentication Request Message encode
add_executable(AuthReqEncode AuthReqEncode.cpp)

target_link_libraries(AuthReqEncode LIB_NAS5G)

#Executable for Authentication Response Message Decode
add_executable(AuthResDecode AuthResDecode.cpp)

target_link_libraries(AuthResDecode LIB_NAS5G)

#Executable for Authentication Result Message encode
add_executable(AuthResultEncode AuthResultEncode.cpp)

target_link_libraries(AuthResultEncode LIB_NAS5G)

#Executable for Authentication Failure Message Decode
add_executable(AuthFailDecode AuthFailDecode.cpp)

target_link_libraries(AuthFailDecode LIB_NAS5G)

#Executable for Authentication Reject Message encode
add_executable(AuthRejEncode AuthRejEncode.cpp)

target_link_libraries(AuthRejEncode LIB_NAS5G)

#Executable for Security mode command encode
add_executable(SecModeCommandEncode SecModeCommandEncode.cpp)

target_link_libraries(SecModeCommandEncode LIB_NAS5G)

#Executable for Security mode complete decode
add_executable(SecurityModeCompleteDecode SecurityModeCompleteDecode.cpp)

target_link_libraries(SecurityModeCompleteDecode LIB_NAS5G)

#Executable for Security mode Reject decode
add_executable(SecurityModeRejectDecode SecurityModeRejectDecode.cpp)

target_link_libraries(SecurityModeRejectDecode LIB_NAS5G)

#Executable for Identity Request encode
add_executable(IdentityRequestEncode IdentityRequestEncode.cpp)

target_link_libraries(IdentityRequestEncode LIB_NAS5G)

#Executable for Identity Response decode
add_executable(IdentityResponseDecode IdentityResponseDecode.cpp)

target_link_libraries(IdentityResponseDecode LIB_NAS5G)

#Executable for UL NAS Tansport decode
add_executable(ULNASTransportDecode ULNASTransportDecode.cpp)

target_link_libraries(ULNASTransportDecode LIB_NAS5G)

#Executable for DL NAS Tansport encode
add_executable(DLNASTransportEncode DLNASTransportEncode.cpp)

target_link_libraries(DLNASTransportEncode LIB_NAS5G)

#Executable for UL NAS Transport PDU Release Request Decode
add_executable(ULNASTransportPDUReleaseRequestDecode ULNASTransportPDUReleaseRequestDecode.cpp)

target_link_libraries(ULNASTransportPDUReleaseRequestDecode LIB_NAS5G)

#Executable for PDU Session Est. Request Decode
add_executable(PDUSessionEstReqDecode PDUSessionEstReqDecode.cpp)

target_link_libraries(PDUSessionEstReqDecode LIB_NAS5G)

#Executable for PDU Session Est. Accept encode
add_executable(PDUSessionEstAcptEncode PDUSessionEstAcptEncode.cpp)

target_link_libraries(PDUSessionEstAcptEncode LIB_NAS5G)

#Executable for PDU Release Request Decode
add_executable(PDUSessionReleaseReqDecode PDUSessionReleaseReqDecode.cpp)

target_link_libraries(PDUSessionReleaseReqDecode LIB_NAS5G)

#Executable for PDU Session Est Reject Encode
add_executable(PDUSessionEstRejEncode PDUSessionEstRejEncode.cpp)

target_link_libraries(PDUSessionEstRejEncode LIB_NAS5G)
 
#Executable for PDU Session Modification Request Decode
add_executable(PDUSessionModReqDecode PDUSessionModReqDecode.cpp)

target_link_libraries(PDUSessionModReqDecode LIB_NAS5G)

#Executable for PDU Session Modification Reject Encode
add_executable(PDUSessionModRejectEncode PDUSessionModRejectEncode.cpp)

target_link_libraries(PDUSessionModRejectEncode LIB_NAS5G)

#Executable for PDU Session Release Reject Encode
add_executable(PDUSessionReleaseRejEncode PDUSessionReleaseRejEncode.cpp)

target_link_libraries(PDUSessionReleaseRejEncode LIB_NAS5G)


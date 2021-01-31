set(NGAP_C_DIR ${PROJECT_BINARY_DIR}/ngap/r15.8)

# compile the needed protos
set(NAS_STATE_CPP_PROTOS nas_state spgw_state common_types std_3gpp_types)
list(APPEND PROTO_SRCS "")
list(APPEND PROTO_HDRS "")

set(LTE_OUT_DIR "${CMAKE_CURRENT_BINARY_DIR}/lte/protos/oai")
generate_cpp_protos("${NAS_STATE_CPP_PROTOS}" "${PROTO_SRCS}"
  "${PROTO_HDRS}" ${STATE_PROTO_DIR} ${LTE_OUT_DIR})

message("Proto_srcs are ${PROTO_SRCS}")

include_directories(${CMAKE_CURRENT_BINARY_DIR})
include_directories(${LTE_OUT_DIR})
include_directories("${NGAP_C_DIR}")
include_directories(${CMAKE_CURRENT_SOURCE_DIR})
include_directories(${CMAKE_CURRENT_SOURCE_DIR}/src)
include_directories(${CMAKE_CURRENT_SOURCE_DIR}/include)
include_directories(${CMAKE_CURRENT_SOURCE_DIR}/src/ies)
include_directories(${CMAKE_CURRENT_SOURCE_DIR}/include/ies)
include_directories(${CMAKE_CURRENT_SOURCE_DIR}/test)
include_directories(${CMAKE_CURRENT_SOURCE_DIR}/../../../../../../orc8r/gateway/c/common/logging)
include_directories(${CMAKE_CURRENT_SOURCE_DIR}/../../common/glogwrapper)
include_directories(${CMAKE_CURRENT_SOURCE_DIR}/../../common)

set(libnas_src_OBJS
	${CMAKE_CURRENT_SOURCE_DIR}/src/AmfMessage.cpp
	${CMAKE_CURRENT_SOURCE_DIR}/src/M5GRegistrationRequest.cpp
	${CMAKE_CURRENT_SOURCE_DIR}/src/M5GRegistrationAccept.cpp
	${CMAKE_CURRENT_SOURCE_DIR}/src/M5GRegistrationComplete.cpp
	${CMAKE_CURRENT_SOURCE_DIR}/src/M5GRegistrationReject.cpp
	${CMAKE_CURRENT_SOURCE_DIR}/src/M5GDeRegistrationRequestUEInit.cpp
	${CMAKE_CURRENT_SOURCE_DIR}/src/M5GDeRegistrationAcceptUEInit.cpp
	${CMAKE_CURRENT_SOURCE_DIR}/src/M5GAuthenticationRequest.cpp
	${CMAKE_CURRENT_SOURCE_DIR}/src/M5GAuthenticationResponse.cpp
	${CMAKE_CURRENT_SOURCE_DIR}/src/M5GAuthenticationResult.cpp
	${CMAKE_CURRENT_SOURCE_DIR}/src/M5GAuthenticationFailure.cpp
	${CMAKE_CURRENT_SOURCE_DIR}/src/M5GAuthenticationReject.cpp
	${CMAKE_CURRENT_SOURCE_DIR}/src/M5GSecurityModeCommand.cpp
	${CMAKE_CURRENT_SOURCE_DIR}/src/M5GSecurityModeComplete.cpp
	${CMAKE_CURRENT_SOURCE_DIR}/src/M5GSecurityModeReject.cpp
	${CMAKE_CURRENT_SOURCE_DIR}/src/M5GIdentityRequest.cpp
	${CMAKE_CURRENT_SOURCE_DIR}/src/M5GIdentityResponse.cpp
	${CMAKE_CURRENT_SOURCE_DIR}/src/M5GPDUSessionEstablishmentRequest.cpp
	${CMAKE_CURRENT_SOURCE_DIR}/src/M5GPDUSessionEstablishmentAccept.cpp
	${CMAKE_CURRENT_SOURCE_DIR}/src/M5GPDUSessionEstablishmentReject.cpp
	${CMAKE_CURRENT_SOURCE_DIR}/src/M5GPDUSessionModificationRequest.cpp
	${CMAKE_CURRENT_SOURCE_DIR}/src/M5GPDUSessionModificationReject.cpp
	${CMAKE_CURRENT_SOURCE_DIR}/src/M5GPDUSessionReleaseRequest.cpp
	${CMAKE_CURRENT_SOURCE_DIR}/src/M5GPDUSessionReleaseReject.cpp
	${CMAKE_CURRENT_SOURCE_DIR}/src/SmfMessage.cpp
	${CMAKE_CURRENT_SOURCE_DIR}/src/M5GULNASTransport.cpp
	${CMAKE_CURRENT_SOURCE_DIR}/src/M5GDLNASTransport.cpp
    )

set(libnas_ies_OBJS
	${CMAKE_CURRENT_SOURCE_DIR}/src/ies/M5GExtendedProtocolDiscriminator.cpp
	${CMAKE_CURRENT_SOURCE_DIR}/src/ies/M5GNasKeySetIdentifier.cpp
	${CMAKE_CURRENT_SOURCE_DIR}/src/ies/M5GSMobileIdentity.cpp
	${CMAKE_CURRENT_SOURCE_DIR}/src/ies/M5GMessageType.cpp
	${CMAKE_CURRENT_SOURCE_DIR}/src/ies/M5GSecurityHeaderType.cpp
	${CMAKE_CURRENT_SOURCE_DIR}/src/ies/M5GSRegistrationType.cpp
	${CMAKE_CURRENT_SOURCE_DIR}/src/ies/M5GSRegistrationResult.cpp
	${CMAKE_CURRENT_SOURCE_DIR}/src/ies/M5GSDeRegistrationType.cpp
	${CMAKE_CURRENT_SOURCE_DIR}/src/ies/M5GSpareHalfOctet.cpp
	${CMAKE_CURRENT_SOURCE_DIR}/src/ies/M5GABBA.cpp
	${CMAKE_CURRENT_SOURCE_DIR}/src/ies/M5GAuthenticationParameterAUTN.cpp
	${CMAKE_CURRENT_SOURCE_DIR}/src/ies/M5GAuthenticationParameterRAND.cpp
	${CMAKE_CURRENT_SOURCE_DIR}/src/ies/M5GAuthenticationResponseParameter.cpp
	${CMAKE_CURRENT_SOURCE_DIR}/src/ies/M5GEAPMessage.cpp
	${CMAKE_CURRENT_SOURCE_DIR}/src/ies/M5GNASSecurityAlgorithms.cpp
	${CMAKE_CURRENT_SOURCE_DIR}/src/ies/M5GUESecurityCapability.cpp
	${CMAKE_CURRENT_SOURCE_DIR}/src/ies/M5GSIdentityType.cpp
	${CMAKE_CURRENT_SOURCE_DIR}/src/ies/M5GPayloadContainer.cpp
	${CMAKE_CURRENT_SOURCE_DIR}/src/ies/M5GPayloadContainerType.cpp
	${CMAKE_CURRENT_SOURCE_DIR}/src/ies/M5GPTI.cpp
	${CMAKE_CURRENT_SOURCE_DIR}/src/ies/M5GSSCMode.cpp
	${CMAKE_CURRENT_SOURCE_DIR}/src/ies/M5GIntegrityProtMaxDataRate.cpp
	${CMAKE_CURRENT_SOURCE_DIR}/src/ies/M5GPDUSessionType.cpp
	${CMAKE_CURRENT_SOURCE_DIR}/src/ies/M5GPDUSessionIdentity.cpp
	${CMAKE_CURRENT_SOURCE_DIR}/src/ies/M5GPDUAddress.cpp
	${CMAKE_CURRENT_SOURCE_DIR}/src/ies/M5GQOSRules.cpp
	${CMAKE_CURRENT_SOURCE_DIR}/src/ies/M5GSessionAMBR.cpp        
	${CMAKE_CURRENT_SOURCE_DIR}/src/ies/M5GSMCause.cpp        
	${CMAKE_CURRENT_SOURCE_DIR}/src/ies/M5GMMCause.cpp        
	${CMAKE_CURRENT_SOURCE_DIR}/src/ies/M5GIMEISVRequest.cpp   
	${CMAKE_CURRENT_SOURCE_DIR}/src/ies/M5GDNN.cpp   
    )

add_library(LIB_NAS5G
    ${libnas_src_OBJS}
    ${libnas_ies_OBJS}
    ${PROTO_SRCS}
    ${PROTO_HDRS}
    )

target_link_libraries(LIB_NAS5G
    COMMON
    LIB_BSTR LIB_SECU LIB_3GPP
    )
#if (EMBEDDED_SGW)
#  target_link_libraries(TASK_NAS5G TASK_SGW)
#endif(EMBEDDED_SGW)

# TODO leak baby leak
target_include_directories(LIB_NAS5G PUBLIC
    ${CMAKE_CURRENT_SOURCE_DIR}
    ${CMAKE_CURRENT_SOURCE_DIR}/src
    ${CMAKE_CURRENT_SOURCE_DIR}/src/ies
    ${CMAKE_CURRENT_SOURCE_DIR}/include
    ${CMAKE_CURRENT_SOURCE_DIR}/include/ies
  # ${CMAKE_CURRENT_SOURCE_DIR}/test
)

add_subdirectory(test)


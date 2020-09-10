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

set(libnas_src_OBJS
	${CMAKE_CURRENT_SOURCE_DIR}/src/AmfMessage.cpp
	${CMAKE_CURRENT_SOURCE_DIR}/src/RegistrationRequest.cpp
	${CMAKE_CURRENT_SOURCE_DIR}/src/RegistrationAccept.cpp
	${CMAKE_CURRENT_SOURCE_DIR}/src/DeRegistrationRequestUEInit.cpp
	${CMAKE_CURRENT_SOURCE_DIR}/src/DeRegistrationAcceptUEInit.cpp
    )

set(libnas_ies_OBJS
	${CMAKE_CURRENT_SOURCE_DIR}/src/ies/ExtendedProtocolDiscriminator.cpp
	${CMAKE_CURRENT_SOURCE_DIR}/src/ies/NasKeySetIdentifier.cpp
	${CMAKE_CURRENT_SOURCE_DIR}/src/ies/M5GSMobileIdentity.cpp
	${CMAKE_CURRENT_SOURCE_DIR}/src/ies/MessageType.cpp
	${CMAKE_CURRENT_SOURCE_DIR}/src/ies/SecurityHeaderType.cpp
	${CMAKE_CURRENT_SOURCE_DIR}/src/ies/M5GSRegistrationType.cpp
	${CMAKE_CURRENT_SOURCE_DIR}/src/ies/M5GSRegistrationResult.cpp
	${CMAKE_CURRENT_SOURCE_DIR}/src/ies/M5GSDeRegistrationType.cpp
	${CMAKE_CURRENT_SOURCE_DIR}/src/ies/SpareHalfOctet.cpp
    )


add_library(TASK_NAS5G
    ${libnas_src_OBJS}
    ${libnas_ies_OBJS}
    ${PROTO_SRCS}
    ${PROTO_HDRS}
    )

target_link_libraries(TASK_NAS5G
    COMMON
    LIB_BSTR LIB_SECU LIB_3GPP
    # TASK_AMF 
    )
if (EMBEDDED_SGW)
  target_link_libraries(TASK_NAS5G TASK_SGW)
endif(EMBEDDED_SGW)

# TODO leak baby leak
target_include_directories(TASK_NAS5G PUBLIC
    ${CMAKE_CURRENT_SOURCE_DIR}
    ${CMAKE_CURRENT_SOURCE_DIR}/src
    ${CMAKE_CURRENT_SOURCE_DIR}/src/ies
    ${CMAKE_CURRENT_SOURCE_DIR}/include
    ${CMAKE_CURRENT_SOURCE_DIR}/include/ies
    ${CMAKE_CURRENT_SOURCE_DIR}/test
)

add_subdirectory(test)


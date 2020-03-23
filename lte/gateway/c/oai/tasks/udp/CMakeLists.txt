set(ASN1RELDIR r10.5)

include_directories(${CMAKE_CURRENT_SOURCE_DIR}/../../common)
include_directories(${CMAKE_CURRENT_SOURCE_DIR}/../../lib/3gpp)
include_directories(${CMAKE_CURRENT_SOURCE_DIR}/../../lib/itti)
include_directories(${PROJECT_BINARY_DIR}/s1ap/${ASN1RELDIR})

include_directories(${CMAKE_CURRENT_SOURCE_DIR}/../../lib/bstr)
include_directories(${CMAKE_CURRENT_SOURCE_DIR}/../../lib/hashtable)

set(S1AP_C_DIR ${PROJECT_BINARY_DIR}/tasks/s1ap/r10.5)
include_directories(${S1AP_C_DIR})

Message(${S1AP_C_DIR} is normally )


# TODO (amar) fix include leak

include_directories("${CMAKE_CURRENT_SOURCE_DIR}/../nas")
include_directories("${CMAKE_CURRENT_SOURCE_DIR}/../nas/emm")
include_directories("${CMAKE_CURRENT_SOURCE_DIR}/../nas/emm/msg")
include_directories("${CMAKE_CURRENT_SOURCE_DIR}/../nas/emm/sap")
include_directories("${CMAKE_CURRENT_SOURCE_DIR}/../nas/ies")
include_directories("${CMAKE_CURRENT_SOURCE_DIR}/../nas/util")
include_directories("${CMAKE_CURRENT_SOURCE_DIR}/../nas/esm")
include_directories("${CMAKE_CURRENT_SOURCE_DIR}/../nas/esm/msg")
include_directories("${CMAKE_CURRENT_SOURCE_DIR}/../nas/api/network")
include_directories("${CMAKE_CURRENT_SOURCE_DIR}/../nas/api/mme")
include_directories("${CMAKE_CURRENT_SOURCE_DIR}/../mme_app")
include_directories("${CMAKE_CURRENT_SOURCE_DIR}/../s1ap/messages/asn1/${ASN1RELDIR}")
include_directories("${CMAKE_CURRENT_SOURCE_DIR}/../s1ap")


add_library(TASK_UDP udp_primitives_server.c)

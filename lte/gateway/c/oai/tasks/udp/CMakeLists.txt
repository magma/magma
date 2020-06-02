set(ASN1RELDIR r10.5)

include_directories(${PROJECT_SOURCE_DIR}/common)
include_directories(${CMAKE_CURRENT_SOURCE_DIR}/../../lib/3gpp)
include_directories(${CMAKE_CURRENT_SOURCE_DIR}/../../lib/itti)
include_directories(${PROJECT_BINARY_DIR}/s1ap/${ASN1RELDIR})

include_directories(${CMAKE_CURRENT_SOURCE_DIR}/../../lib/bstr)
include_directories(${CMAKE_CURRENT_SOURCE_DIR}/../../lib/hashtable)

set(S1AP_C_DIR ${PROJECT_BINARY_DIR}/tasks/s1ap/r10.5)
include_directories(${S1AP_C_DIR})

add_library(TASK_UDP udp_primitives_server.c)

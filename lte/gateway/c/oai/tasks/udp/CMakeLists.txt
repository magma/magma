include_directories(${PROJECT_SOURCE_DIR}/common)
include_directories(${CMAKE_CURRENT_SOURCE_DIR}/../../lib/3gpp)
include_directories(${CMAKE_CURRENT_SOURCE_DIR}/../../lib/itti)

include_directories(${CMAKE_CURRENT_SOURCE_DIR}/../../lib/bstr)
include_directories(${CMAKE_CURRENT_SOURCE_DIR}/../../lib/hashtable)

add_library(TASK_UDP udp_primitives_server.c)

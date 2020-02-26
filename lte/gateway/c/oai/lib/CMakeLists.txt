set(LIB_DIR ${CMAKE_CURRENT_SOURCE_DIR})

add_subdirectory(3gpp) # LIB_3GPP
add_subdirectory(bstr) # LIB_BSTR
add_subdirectory(directoryd) # LIB_DIRECTORYD
add_subdirectory(hashtable) # LIB_HASHTABLE
add_subdirectory(itti) # LIB_ITTI
add_subdirectory(message_utils) # LIB_MESSAGE_UTILS
add_subdirectory(pcef) # LIB_PCEF
add_subdirectory(rpc_client) # LIB_RPC_CLIENT
add_subdirectory(s6a_proxy) # LIB_S6A_PROXY
add_subdirectory(secu) # LIB_SECU
add_subdirectory(sgs_client) # LIB_SGS_CLIENT

if (EMBEDDED_SGW)
if (ENABLE_OPENFLOW)
add_subdirectory(openflow) # LIB_OPENFLOW
add_subdirectory(openflow/controller) # LIB_OPENFLOW_CONTROLLER
endif (ENABLE_OPENFLOW)
else (EMBEDDED_SGW)
  add_subdirectory(gtpv2-c) # LIB_GTPV2C
endif (NOT EMBEDDED_SGW)

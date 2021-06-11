set(S1AP_C_DIR ${PROJECT_BINARY_DIR}/s1ap/r15)
include_directories(${S1AP_C_DIR})

set(S6A_SRC
    s6a_auth_info.c
    s6a_error.c
    s6a_peer.c
    s6a_subscription_data.c
    s6a_task.c
    s6a_up_loc.c
    s6a_cancel_loc.c
    s6a_service_handler.c
    s6a_purge_ue.c
    s6a_hss_reset.c
    s6a_iface.cpp
    )

set(FD_LIBS)

if (S6A_OVER_GRPC)
  set(S6A_SRC ${S6A_SRC} s6a_grpc_iface.cpp)
  add_library(TASK_S6A ${S6A_SRC})
  target_link_libraries(TASK_S6A
      COMMON
      LIB_BSTR LIB_HASHTABLE LIB_S6A_PROXY
      )
else (S6A_OVER_GRPC)  # Use freeDiameter
  set(S6A_SRC ${S6A_SRC}
      s6a_dict.c
      s6a_fd_iface.cpp)
  add_library(TASK_S6A ${S6A_SRC})
  target_link_libraries(TASK_S6A
      COMMON
      LIB_BSTR LIB_HASHTABLE LIB_S6A_PROXY gnutls fdproto fdcore
      )
endif (S6A_OVER_GRPC)

target_include_directories(TASK_S6A PUBLIC
    ${CMAKE_CURRENT_SOURCE_DIR}
    )

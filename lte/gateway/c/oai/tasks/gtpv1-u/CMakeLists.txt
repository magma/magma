
set (GTPV1U_SRC
    gtpv1u_task.c
)

if (ENABLE_OPENFLOW)  # Use openflow
  set (GTPV1U_SRC ${GTPV1U_SRC} gtp_tunnel_openflow.c)
else ()  # Use libgtpnl
  pkg_search_module(GTPNL libgtpnl REQUIRED)
  include_directories(${GTPNL_INCLUDE_DIRS})
  set (GTPV1U_SRC ${GTPV1U_SRC} gtp_tunnel_libgtpnl.c)
endif ()

set(S1AP_C_DIR ${PROJECT_BINARY_DIR}/s1ap/r10.5)
include_directories(${S1AP_C_DIR})

add_library(TASK_GTPV1U ${GTPV1U_SRC})
target_include_directories(TASK_GTPV1U PUBLIC ${CMAKE_CURRENT_SOURCE_DIR})

target_link_libraries(TASK_GTPV1U
  COMMON
  LIB_BSTR LIB_HASHTABLE LIB_OPENFLOW_CONTROLLER LIB_MOBILITY_CLIENT
  TASK_NAS TASK_MME_APP TASK_SERVICE303 TASK_SGW
)

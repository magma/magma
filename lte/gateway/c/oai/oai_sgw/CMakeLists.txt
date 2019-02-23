# Override options for S/P GW
set(BYTE_ORDER=LITTLE_ENDIAN)
set(DISPLAY_LICENCE_INFO            True)
set(GTPV1U_LINEAR_TEID_ALLOCATION   False)
set(LOG_OAI                         True)
set(MESSAGE_CHART_GENERATOR         True)
set(MEMORY_CHECK                    False)
set(NW_GTPV2C_DISPLAY_LICENCE_INFO  True)
set(S1AP_DEBUG_LIST                 True)
set(SECU_DEBUG                      False)
set(SPGW_BUILD                      True)
set(TRACE_HASHTABLE                 False)
set(DISABLE_EXECUTE_SHELL_COMMAND   False)


if (NOT EMBEDDED_SGW)
  add_executable(spgw
    oai_sgw.c
    ${PROJECT_SOURCE_DIR}/common/3gpp_24.008.c
    ${PROJECT_SOURCE_DIR}/service303/service303_task.c
    ${PROJECT_SOURCE_DIR}/service303/service303_spgw_stats.c
  )

  target_link_libraries(spgw
    -Wl,--start-group
      SERVICE303_LIB lfds710 COMMON
      LIB_BSTR LIB_HASHTABLE LIB_OPENFLOW LIB_GTPV2C LIB_RPC_CLIENT
      TASK_SGW TASK_S11_SGW TASK_UDP_SERVER TASK_GTPV1U
      ${MSC_LIB} ${ITTI_LIB} ${GCOV_LIB}
    -Wl,--end-group
    pthread m rt ${GTPNL_LIBRARIES} ${CONFIG_LIBRARIES}
  )

  if (ENABLE_OPENFLOW)
    target_link_libraries(spgw openvswitch)
  endif (ENABLE_OPENFLOW)
endif (NOT EMBEDDED_SGW)

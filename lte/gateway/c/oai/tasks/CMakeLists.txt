
add_subdirectory(gtpv1-u) # TASK_GTPV1U
add_subdirectory(mme_app) # TASK_MME_APP
add_subdirectory(nas) # TASK_NAS
add_subdirectory(s1ap) # TASK_S1AP, LIB_S1AP
add_subdirectory(s6a) # TASK_S6A
add_subdirectory(s6a_service) # TASK_S6A_SERVICE
add_subdirectory(sctp) # TASK_SCTP_SERVER
add_subdirectory(service303) # TASK_SERVICE303
add_subdirectory(sgs) # TASK_SGS
add_subdirectory(sgs_service) # TASK_SGS_SERVICE
add_subdirectory(sgw) # TASK_SGW
add_subdirectory(udp) # TASK_UDP_SERVER

if (NOT EMBEDDED_SGW)
  add_subdirectory(s11) # TASK_S11_MME, TASK_S11_SGW
endif (NOT EMBEDDED_SGW)

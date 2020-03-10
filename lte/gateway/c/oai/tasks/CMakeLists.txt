add_subdirectory(mme_app) # TASK_MME_APP
add_subdirectory(nas) # TASK_NAS
add_subdirectory(s1ap) # TASK_S1AP, LIB_S1AP
add_subdirectory(s6a) # TASK_S6A
add_subdirectory(sctp) # TASK_SCTP_SERVER
add_subdirectory(service303) # TASK_SERVICE303
add_subdirectory(sgs) # TASK_SGS
add_subdirectory(grpc_service) # TASK_GRPC_SERVICE

if (EMBEDDED_SGW)
add_subdirectory(gtpv1-u) # TASK_GTPV1U
add_subdirectory(sgw) # TASK_SGW
else (EMBEDDED_SGW)
add_subdirectory(s11) # TASK_S11
endif (EMBEDDED_SGW)


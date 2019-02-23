# TODO: ITTI_Messages shouldn't depend on S1_Types
set(ASN1RELDIR r10.5)
set(asn1_generated_dir ${PROJECT_BINARY_DIR}/s1ap)
set(asn1c_call "${CMAKE_CURRENT_SOURCE_DIR}/generate_asn1")
set(S1AP_C_DIR ${asn1_generated_dir}/${ASN1RELDIR})

add_library(TASK_SGW
    spgw_config.c
    pgw_config.c
    sgw_config.c
    s11_causes.c
    sgw_task.c
    sgw_handlers.c
    pgw_task.c
    pgw_handlers.c
    sgw_context_manager.c
    pgw_pco.c
    mobilityd_ue_ip_address_alloc.c
    sgw_paging.c
    pgw_pcef_emulation.c
    pgw_procedures.c
    )
target_compile_definitions(TASK_SGW PRIVATE
    PACKAGE_NAME=\"S/P-GW\"
    PACKAGE_VERSION=\"0.1\"
    PACKAGE_BUGREPORT=\"TBD\"
)
target_link_libraries(TASK_SGW
    COMMON
    ${GTPNL_LIBRARIES}
    LIB_BSTR LIB_HASHTABLE LIB_RPC_CLIENT LIB_PCEF
    TASK_GTPV1U
)
target_include_directories(TASK_SGW PUBLIC
    ${CMAKE_CURRENT_SOURCE_DIR}
    ${S1AP_C_DIR}
)

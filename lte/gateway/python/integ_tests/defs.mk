# Copyright (c) Facebook, Inc. and its affiliates.
# All rights reserved.
#
# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.
#

PROTO_LIST:=orc8r_protos lte_protos feg_protos

# Add the s1aptester integration tests
MANDATORY_TESTS = s1aptests/test_attach_detach.py \
s1aptests/test_gateway_metrics_attach_detach.py \
s1aptests/test_attach_detach_multi_ue.py \
s1aptests/test_attach_detach_looped.py  \
s1aptests/test_attach_emergency.py \
s1aptests/test_attach_combined_eps_imsi.py \
s1aptests/test_attach_via_guti.py \
s1aptests/test_attach_without_ips_available.py \
s1aptests/test_attach_detach_after_ue_context_release.py \
s1aptests/test_attach_detach_duplicate_nas_resp_messages.py \
s1aptests/test_attach_security_mode_reject.py \
s1aptests/test_attach_esm_information.py \
s1aptests/test_attach_esm_information_wrong_apn.py \
s1aptests/test_attach_ue_ctxt_release_cmp_delay.py \
s1aptests/test_attach_auth_failure.py \
s1aptests/test_nas_non_delivery_for_smc.py \
s1aptests/test_nas_non_delivery_for_identity_req.py \
s1aptests/test_attach_no_initial_context_resp.py \
s1aptests/test_attach_detach_no_ueContext_release_comp.py \
s1aptests/test_no_attach_complete.py \
s1aptests/test_no_auth_response.py \
s1aptests/test_no_security_mode_complete.py \
s1aptests/test_tau_periodic_inactive.py \
s1aptests/test_tau_periodic_active.py \
s1aptests/test_attach_service.py \
s1aptests/test_attach_detach_service.py \
s1aptests/test_attach_service_ue_radio_capability.py \
s1aptests/test_attach_service_multi_ue.py \
s1aptests/test_attach_ipv4v6_pdn_type.py \
s1aptests/test_service_info.py \
s1aptests/test_attach_detach_with_ovs.py \
s1aptests/test_resync.py \
s1aptests/test_standalone_pdn_conn_req.py \
s1aptests/test_attach_act_dflt_ber_ctxt_rej.py \
s1aptests/test_attach_detach_EEA1.py \
s1aptests/test_attach_detach_EEA2.py \
s1aptests/test_attach_detach_EIA1.py \
s1aptests/test_attach_detach_emm_status.py \
s1aptests/test_attach_detach_enb_rlf_initial_ue_msg.py \
s1aptests/test_attach_detach_ICS_Failure.py \
s1aptests/test_attach_detach_ps_service_not_available.py \
s1aptests/test_attach_missing_imsi.py \
s1aptests/test_duplicate_attach.py \
s1aptests/test_enb_partial_reset_con_dereg.py \
s1aptests/test_enb_partial_reset.py \
s1aptests/test_enb_complete_reset.py \
s1aptests/test_nas_non_delivery_for_auth.py \
s1aptests/test_outoforder_attach_complete_ICSR.py \
s1aptests/test_s1setup_incorrect_plmn.py \
s1aptests/test_s1setup_incorrect_tac.py \
s1aptests/test_sctp_abort_after_auth_req.py \
s1aptests/test_sctp_abort_after_identity_req.py \
s1aptests/test_sctp_abort_after_smc.py \
s1aptests/test_sctp_shutdown_after_auth_req.py \
s1aptests/test_sctp_shutdown_after_identity_req.py \
s1aptests/test_sctp_shutdown_after_smc.py \
s1aptests/test_sctp_shutdown_after_multi_ue_attach.py \
s1aptests/test_attach_detach_dedicated.py \
s1aptests/test_attach_detach_dedicated_qci_0.py \
s1aptests/test_attach_detach_dedicated_multi_ue.py \
s1aptests/test_attach_detach_dedicated_looped.py \
s1aptests/test_attach_detach_dedicated_deactivation_timer_expiry.py \
s1aptests/test_attach_detach_dedicated_bearer_deactivation_invalid_lbi.py \
s1aptests/test_attach_detach_dedicated_bearer_deactivation_invalid_imsi.py \
s1aptests/test_attach_detach_dedicated_bearer_deactivation_invalid_ebi.py \
s1aptests/test_attach_detach_dedicated_bearer_activation_invalid_lbi.py \
s1aptests/test_attach_detach_dedicated_bearer_activation_invalid_imsi.py \
s1aptests/test_attach_detach_dedicated_activation_timer_expiry.py \
s1aptests/test_attach_detach_dedicated_activation_reject.py \
s1aptests/test_attach_detach_multiple_dedicated.py \
s1aptests/test_attach_detach_secondary_pdn.py \
s1aptests/test_attach_detach_secondary_pdn_multi_ue.py \
s1aptests/test_attach_detach_secondary_pdn_looped.py \
s1aptests/test_attach_detach_secondary_pdn_invalid_apn.py \
s1aptests/test_attach_detach_secondary_pdn_disconnect_dedicated_bearer.py \
s1aptests/test_attach_detach_secondary_pdn_disconnect_invalid_bearer.py \
s1aptests/test_attach_detach_secondary_pdn_no_disconnect.py \
s1aptests/test_attach_detach_secondary_pdn_with_dedicated_bearer_looped.py \
s1aptests/test_attach_detach_secondary_pdn_with_dedicated_bearer_multi_ue.py \
s1aptests/test_attach_detach_secondary_pdn_with_dedicated_bearer.py \
s1aptests/test_attach_detach_secondary_pdn_with_dedicated_bearer_deactivate.py \
s1aptests/test_attach_detach_disconnect_default_pdn.py \
s1aptests/test_attach_detach_maxbearers_twopdns.py \
s1aptests/test_attach_detach_multiple_secondary_pdn.py \
s1aptests/test_attach_detach_nw_triggered_delete_secondary_pdn.py \
s1aptests/test_attach_detach_nw_triggered_delete_last_pdn.py \
s1aptests/test_different_enb_s1ap_id_same_ue.py \
s1aptests/test_attach_detach_with_pcscf_address.py \
s1aptests/test_attach_detach_secondary_pdn_with_pcscf_address.py \
s1aptests/test_multi_enb_multi_ue.py \
s1aptests/test_multi_enb_multi_ue_diff_enbtype.py \
s1aptests/test_multi_enb_partial_reset.py \
s1aptests/test_multi_enb_complete_reset.py \
s1aptests/test_multi_enb_sctp_shutdown.py \
s1aptests/test_attach_ul_udp_data.py \
s1aptests/test_attach_ul_tcp_data.py \
s1aptests/test_attach_detach_rar_tcp_data.py \
s1aptests/test_attach_detach_multiple_rar_tcp_data.py \
s1aptests/test_attach_detach_attach_ul_tcp_data.py

# These test cases pass without memory leaks, but needs DL-route in TRF server
# sudo /sbin/route add -net 192.168.128.0 gw 192.168.60.142
#     netmask 255.255.255.0 dev eth1
# s1aptests/test_attach_dl_udp_data.py \
# s1aptests/test_attach_dl_tcp_data.py \
# s1aptests/test_attach_detach_attach_dl_tcp_data.py

# TODO flaky tests we should look at
# s1aptests/test_attach_detach_multi_ue_looped.py \

CLOUD_TESTS = cloud_tests/checkin_test.py \
cloud_tests/metrics_export_test.py \
cloud_tests/config_test.py

S1AP_TESTER_CFG=$(MAGMA_ROOT)/lte/gateway/python/integ_tests/data/s1ap_tester_cfg
S1AP_TESTER_PYTHON_PATH=$(S1AP_TESTER_ROOT)/bin

# Local integ tests are run on the magma access gateway, not the test VM
LOCAL_INTEG_TESTS = gxgy_tests

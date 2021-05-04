# Copyright 2020 The Magma Authors.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#

PROTO_LIST:=orc8r_protos lte_protos feg_protos

# Add the s1aptester integration tests
PRECOMMIT_TESTS = s1aptests/test_attach_detach.py \
s1aptests/test_gateway_metrics_attach_detach.py \
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
s1aptests/test_attach_detach_security_algo_eea0_eia0.py \
s1aptests/test_attach_detach_security_algo_eea1_eia1.py \
s1aptests/test_attach_detach_security_algo_eea2_eia2.py \
s1aptests/test_attach_detach_emm_status.py \
s1aptests/test_attach_detach_enb_rlf_initial_ue_msg.py \
s1aptests/test_attach_detach_ICS_Failure.py \
s1aptests/test_attach_missing_imsi.py \
s1aptests/test_duplicate_attach.py \
s1aptests/test_enb_partial_reset_con_dereg.py \
s1aptests/test_enb_partial_reset.py \
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
s1aptests/test_attach_detach_multi_ue.py \
s1aptests/test_attach_detach_dedicated.py \
s1aptests/test_attach_detach_dedicated_qci_0.py \
s1aptests/test_attach_detach_dedicated_multi_ue.py \
s1aptests/test_attach_detach_dedicated_looped.py \
s1aptests/test_attach_detach_dedicated_bearer_deactivation_invalid_lbi.py \
s1aptests/test_attach_detach_dedicated_bearer_deactivation_invalid_imsi.py \
s1aptests/test_attach_detach_dedicated_bearer_deactivation_invalid_ebi.py \
s1aptests/test_attach_detach_dedicated_bearer_activation_invalid_lbi.py \
s1aptests/test_attach_detach_dedicated_bearer_activation_invalid_imsi.py \
s1aptests/test_attach_detach_dedicated_activation_timer_expiry.py \
s1aptests/test_attach_detach_dedicated_activation_reject.py \
s1aptests/test_attach_detach_multiple_dedicated.py \
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
s1aptests/test_secondary_pdn_reject_multiple_sessions_not_allowed_per_apn.py \
s1aptests/test_secondary_pdn_reject_unknown_pdn_type.py \
s1aptests/test_attach_standalone_act_dflt_ber_ctxt_rej.py \
s1aptests/test_attach_standalone_act_dflt_ber_ctxt_rej_ded_bearer_activation.py \
s1aptests/test_ics_timer_expiry_ue_registered.py \
s1aptests/test_ics_timer_expiry_ue_unregistered.py \
s1aptests/test_attach_service_with_multi_pdns_and_bearers_looped.py \
s1aptests/test_attach_service_with_multi_pdns_and_bearers_multi_ue.py \
s1aptests/test_attach_service_with_multi_pdns_and_bearers_failure.py \
s1aptests/test_dedicated_bearer_activation_idle_mode.py \
s1aptests/test_dedicated_bearer_activation_idle_mode_multi_ue.py \
s1aptests/test_dedicated_bearer_activation_idle_mode_paging_timer_expiry.py \
s1aptests/test_multi_enb_multi_ue.py \
s1aptests/test_multi_enb_multi_ue_diff_enbtype.py \
s1aptests/test_multi_enb_partial_reset.py \
s1aptests/test_multi_enb_complete_reset.py \
s1aptests/test_multi_enb_sctp_shutdown.py \
s1aptests/test_attach_ul_udp_data.py \
s1aptests/test_attach_ul_tcp_data.py \
s1aptests/test_attach_detach_attach_ul_tcp_data.py \
s1aptests/test_attach_detach_multiple_rar_tcp_data.py \
s1aptests/test_attach_service_with_multi_pdns_and_bearers_mt_data.py \
s1aptests/test_attach_asr.py \
s1aptests/test_attach_detach_with_sctpd_restart.py \
s1aptests/test_attach_nw_initiated_detach_with_mme_restart.py \
s1aptests/test_attach_detach_multiple_ip_blocks_mobilityd_restart.py \
s1aptests/test_attach_ul_udp_data_with_mme_restart.py \
s1aptests/test_attach_ul_udp_data_with_mobilityd_restart.py \
s1aptests/test_attach_ul_udp_data_with_multiple_service_restart.py \
s1aptests/test_attach_ul_udp_data_with_pipelined_restart.py \
s1aptests/test_attach_ul_udp_data_with_sessiond_restart.py \
s1aptests/test_service_req_ul_udp_data_with_mme_restart.py \
s1aptests/test_attach_detach_setsessionrules_tcp_data.py

EXTENDED_TESTS = s1aptests/test_modify_mme_config_for_sanity.py \
s1aptests/test_attach_detach_dedicated_deactivation_timer_expiry.py \
s1aptests/test_attach_detach_secondary_pdn.py \
s1aptests/test_attach_service_with_multi_pdns_and_bearers.py \
s1aptests/test_multi_enb_multi_ue_diff_plmn.py \
s1aptests/test_multi_enb_multi_ue_diff_tac.py \
s1aptests/test_x2_handover.py \
s1aptests/test_x2_handover_ping_pong.py \
s1aptests/test_attach_mobile_reachability_timer_expiry.py \
s1aptests/test_attach_implicit_detach_timer_expiry.py \
s1aptests/test_attach_detach_rar_tcp_data.py \
s1aptests/test_attach_detach_with_mme_restart.py \
s1aptests/test_attach_detach_with_mobilityd_restart.py \
s1aptests/test_idle_mode_with_mme_restart.py \
s1aptests/test_restore_mme_config_after_sanity.py

# Enable these tests once the CI job time-out has increased
# s1aptests/test_mobile_reachability_timer_with_mme_restart.py \
# s1aptests/test_implicit_detach_timer_with_mme_restart.py \
# s1aptests/test_ics_timer_expiry_with_mme_restart.py \

# These test cases pass without memory leaks, but needs DL-route in TRF server
# sudo /sbin/route add -net 192.168.128.0 gw 192.168.60.142
#     netmask 255.255.255.0 dev eth1
# s1aptests/test_attach_dl_udp_data.py \
# s1aptests/test_attach_dl_tcp_data.py \
# s1aptests/test_attach_detach_attach_dl_tcp_data.py

# TODO flaky tests we should look at
# s1aptests/test_attach_detach_ps_service_not_available.py \
# s1aptests/test_enb_complete_reset.py \
# s1aptests/test_attach_detach_multi_ue_looped.py \
# s1aptests/test_attach_detach_rar_tcp_he.py \ GitHubIssue 6254

CLOUD_TESTS = cloud_tests/checkin_test.py \
cloud_tests/metrics_export_test.py \
cloud_tests/config_test.py

S1AP_TESTER_CFG=$(MAGMA_ROOT)/lte/gateway/python/integ_tests/data/s1ap_tester_cfg
S1AP_TESTER_PYTHON_PATH=$(S1AP_TESTER_ROOT)/bin

# Local integ tests are run on the magma access gateway, not the test VM
LOCAL_INTEG_TESTS = gxgy_tests

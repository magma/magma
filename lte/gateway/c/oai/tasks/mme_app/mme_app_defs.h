/*
 * Licensed to the OpenAirInterface (OAI) Software Alliance under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The OpenAirInterface Software Alliance licenses this file to You under
 * the Apache License, Version 2.0  (the "License"); you may not use this file
 * except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *-------------------------------------------------------------------------------
 * For more information about the OpenAirInterface (OAI) Software Alliance:
 *      contact@openairinterface.org
 */

/*! \file mme_app_defs.h
  \brief
  \author Sebastien ROUX, Lionel Gauthier
  \company Eurecom
  \email: lionel.gauthier@eurecom.fr
*/

/* This file contains definitions related to mme applicative layer and should
 * not be included within other layers.
 * Use mme_app_extern.h to expose mme applicative layer procedures/data.
 */

#ifndef FILE_MME_APP_DEFS_SEEN
#define FILE_MME_APP_DEFS_SEEN

#include "intertask_interface.h"

#include "mme_app_desc.h"
#include "mme_app_ue_context.h"
#include "mme_app_sgs_fsm.h"

int mme_app_handle_s1ap_ue_capabilities_ind(
  const itti_s1ap_ue_cap_ind_t const *s1ap_ue_cap_ind_pP);

void mme_app_handle_s1ap_ue_context_release_complete(
  const itti_s1ap_ue_context_release_complete_t const
    *s1ap_ue_context_release_complete);

int mme_app_send_s6a_update_location_req(
  struct ue_mm_context_s *const ue_context_pP);

int mme_app_handle_s6a_update_location_ans(
  const s6a_update_location_ans_t *const ula_pP);

int mme_app_handle_s6a_cancel_location_req(
  const s6a_cancel_location_req_t *const clr_pP);

int mme_app_handle_nas_pdn_connectivity_req(
  itti_nas_pdn_connectivity_req_t *const nas_pdn_connectivity_req_p);

int mme_app_handle_nas_extended_service_req(
  itti_nas_extended_service_req_t *const nas_extended_service_req_pP);

void mme_app_handle_detach_req(const itti_nas_detach_req_t *const detach_req_p);

void mme_app_handle_sgs_detach_req(
  const itti_nas_sgs_detach_req_t *const sgs_detach_req_p);

int mme_app_handle_sgs_eps_detach_ack(
  const const itti_sgsap_eps_detach_ack_t *const eps_detach_ack_p);

int mme_app_handle_sgs_imsi_detach_ack(
  const const itti_sgsap_imsi_detach_ack_t *const imsi_detach_ack_p);

void mme_app_handle_conn_est_cnf(
  itti_nas_conn_est_cnf_t *const nas_conn_est_cnf_pP);

void mme_app_handle_initial_ue_message(
  itti_s1ap_initial_ue_message_t *const conn_est_ind_pP);

int mme_app_handle_create_sess_resp(
  itti_s11_create_session_response_t *const
    create_sess_resp_pP); //not const because we need to free internal stucts

void mme_app_handle_delete_session_rsp(
  const itti_s11_delete_session_response_t *const delete_sess_respP);

void mme_app_handle_erab_setup_req(
  itti_erab_setup_req_t *const itti_erab_setup_req);

int mme_app_handle_establish_ind(
  const nas_establish_ind_t *const nas_establish_ind_pP);

int mme_app_handle_authentication_info_answer(
  const s6a_auth_info_ans_t *const s6a_auth_info_ans_pP);

void mme_app_handle_release_access_bearers_resp(
  const itti_s11_release_access_bearers_response_t
    *const rel_access_bearers_rsp_pP);

void mme_app_handle_s11_create_bearer_req(
  const itti_s11_create_bearer_request_t *const create_bearer_request_pP);

void mme_app_handle_nas_auth_param_req(
  const itti_nas_auth_param_req_t *const nas_auth_param_req_pP);

void mme_app_handle_initial_context_setup_rsp(
  itti_mme_app_initial_context_setup_rsp_t *const initial_ctxt_setup_rsp_pP);

void mme_app_handle_initial_context_setup_failure(
  const itti_mme_app_initial_context_setup_failure_t
    *const initial_ctxt_setup_failure_pP);

bool mme_app_dump_ue_context(
  const hash_key_t keyP,
  void *const ue_context_pP,
  void *unused_param_pP,
  void **unused_result_pP);

int mme_app_handle_nas_dl_req(itti_nas_dl_data_req_t *const nas_dl_req_pP);

void mme_app_handle_e_rab_setup_rsp(
  itti_s1ap_e_rab_setup_rsp_t *const e_rab_setup_rsp);

void mme_app_handle_create_dedicated_bearer_rsp(
  itti_mme_app_create_dedicated_bearer_rsp_t
    *const create_dedicated_bearer_rsp);

void mme_app_handle_create_dedicated_bearer_rej(
  itti_mme_app_create_dedicated_bearer_rej_t
    *const create_dedicated_bearer_rej);

void mme_ue_context_update_ue_sig_connection_state(
  mme_ue_context_t *const mme_ue_context_p,
  struct ue_mm_context_s *ue_context_p,
  ecm_state_t new_ecm_state);

void mme_app_handle_mobile_reachability_timer_expiry(
  struct ue_mm_context_s *ue_context_p);

void mme_app_handle_implicit_detach_timer_expiry(
  struct ue_mm_context_s *ue_context_p);

void mme_app_handle_initial_context_setup_rsp_timer_expiry(
  struct ue_mm_context_s *ue_context_p);

void mme_app_handle_ue_context_modification_timer_expiry(
  struct ue_mm_context_s *ue_context_p);

void mme_app_handle_enb_reset_req(
  const itti_s1ap_enb_initiated_reset_req_t const *enb_reset_req);

int mme_app_handle_initial_paging_request(const char *imsi);

int mme_app_handle_paging_timer_expiry(ue_mm_context_t *ue_context_p);
void mme_app_handle_ulr_timer_expiry(ue_mm_context_t *ue_context_p);

void mme_app_handle_sgs_eps_detach_timer_expiry(ue_mm_context_t *ue_context_p);
void mme_app_handle_sgs_imsi_detach_timer_expiry(ue_mm_context_t *ue_context_p);
void mme_app_handle_sgs_implicit_imsi_detach_timer_expiry(
  ue_mm_context_t *ue_context_p);
void mme_app_handle_sgs_implicit_eps_detach_timer_expiry(
  ue_mm_context_t *ue_context_p);

int mme_app_send_s6a_cancel_location_ans(
  int cla_result,
  const char *imsi,
  uint8_t imsi_length,
  void *msg_cla_p);

int mme_app_send_s6a_purge_ue_req(struct ue_mm_context_s *const ue_context_pP);

int mme_app_handle_s6a_purge_ue_ans(const s6a_purge_ue_ans_t *const pua_pP);

void mme_app_handle_suspend_acknowledge(
  const itti_s11_suspend_acknowledge_t *const suspend_acknowledge);

int mme_app_send_s11_suspend_notification(
  struct ue_mm_context_s *const ue_context_pP, const pdn_cid_t cid);

int mme_app_handle_s6a_reset_req(const s6a_reset_req_t *const rsr_pP);

int mme_app_send_s6a_reset_ans(int rsa_result, void *msg_rsa_p);

int mme_app_send_nas_cs_service_notification(
  mme_ue_s1ap_id_t ue_id,
  uint8_t paging_id,
  bstring cli);

int mme_app_send_sgsap_service_request(
  uint8_t service_indicator,
  struct ue_mm_context_s *ue_context_p);

int mme_app_send_nas_detach_request(
  mme_ue_s1ap_id_t ue_id,
  uint8_t detach_type);

int mme_app_handle_nas_cs_domain_location_update_req(
  itti_nas_cs_domain_location_update_req_t *const itti_nas_location_update_req);

int mme_app_handle_sgsap_location_update_acc(
  itti_sgsap_location_update_acc_t *const itti_sgsap_location_update_acc);

int send_itti_sgsap_location_update_req(ue_mm_context_t *ue_context);

int mme_app_handle_sgsap_location_update_rej(
  itti_sgsap_location_update_rej_t *const itti_sgsap_location_update_rej);

int send_cs_domain_loc_updt_fail_to_nas(
  SgsRejectCause_t cause,
  lai_t *lai,
  mme_ue_s1ap_id_t mme_ue_s1ap_id);
void mme_app_handle_ts6_1_timer_expiry(struct ue_mm_context_s *ue_context_p);

int mme_app_handle_sgsap_reset_indication(
  itti_sgsap_vlr_reset_indication_t *const reset_indication_pP);

int sgs_fsm_associated_reset_indication(const sgs_fsm_t *fsm_evt);

bool mme_app_handle_reset_indication(
  const hash_key_t keyP,
  void *const ue_context_pP,
  void *unused_param_pP,
  void **unused_result_pP);

int mme_app_handle_sgsap_alert_request(
  itti_sgsap_alert_request_t *const sgsap_alert_req_pP);

int mme_app_paging_request_helper(
  ue_mm_context_t *ue_context_p,
  bool set_timer,
  uint8_t paging_id_imsi,
  s1ap_cn_domain_t domain_indicator);

int mme_app_handle_sgsap_paging_request(
  itti_sgsap_paging_request_t *const sgsap_paging_req_pP);

int mme_app_send_sgsap_paging_reject(
  struct ue_mm_context_s *ue_context_p,
  imsi64_t imsi,
  uint8_t imsi_len,
  SgsCause_t sgs_cause);

int mme_app_notify_service_reject_to_nas(
  mme_ue_s1ap_id_t ue_id,
  uint8_t emm_cause,
  uint8_t failed_procedure);

int handle_csfb_s1ap_procedure_failure(
  ue_mm_context_t *ue_context_p,
  char *failed_statement,
  uint8_t failed_procedure);

void mme_app_handle_nas_tau_complete(
  itti_nas_tau_complete_t *itti_nas_tau_complete_p);

int mme_app_handle_sgsap_service_abort_request(
  itti_sgsap_service_abort_req_t *const itti_sgsap_service_abort_req_p);

void mme_app_handle_modify_ue_ambr_request(
  const itti_s11_modify_ue_ambr_request_t *const modify_ue_ambr_request_p);

void mme_app_handle_nw_init_ded_bearer_actv_req(
  const itti_s11_nw_init_actv_bearer_request_t
  *const nw_init_bearer_actv_req_p);

int mme_app_handle_sgs_status_message(
    itti_sgsap_status_t *const sgsap_status_pP);

void mme_app_handle_delete_dedicated_bearer_rsp(
  itti_mme_app_delete_dedicated_bearer_rsp_t
  *const delete_dedicated_bearer_rsp);

void mme_app_handle_erab_rel_cmd(
  itti_erab_rel_cmd_t *const itti_erab_rel_cmd);

void mme_app_handle_e_rab_rel_rsp(
  itti_s1ap_e_rab_rel_rsp_t *const e_rab_rel_rsp);

void mme_app_handle_nw_init_bearer_deactv_req(
  itti_s11_nw_init_deactv_bearer_request_t
  *const nw_init_bearer_deactv_req_p);

void mme_app_handle_delete_dedicated_bearer_rej(
  itti_mme_app_delete_dedicated_bearer_rej_t
  *const delete_dedicated_bearer_rej);

void mme_app_handle_path_switch_request(
    itti_s1ap_path_switch_request_t *const path_switch_req_p);

bool is_e_rab_id_present(
    e_rab_to_be_switched_in_downlink_list_t e_rab_to_be_switched_dl_list,
    ebi_t bearer_id);

void mme_app_handle_path_switch_req_ack(
    itti_s11_modify_bearer_response_t  *const s11_modify_bearer_response,
    struct ue_mm_context_s *ue_context_p);

void mme_app_handle_path_switch_req_failure(
    struct ue_mm_context_s *ue_context_p);

#define mme_stats_read_lock(mMEsTATS)                                          \
  pthread_rwlock_rdlock(&(mMEsTATS)->rw_lock)
#define mme_stats_write_lock(mMEsTATS)                                         \
  pthread_rwlock_wrlock(&(mMEsTATS)->rw_lock)
#define mme_stats_unlock(mMEsTATS) pthread_rwlock_unlock(&(mMEsTATS)->rw_lock)

#endif /* MME_APP_DEFS_H_ */

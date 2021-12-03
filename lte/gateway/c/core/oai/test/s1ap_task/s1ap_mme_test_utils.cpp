/**
 * Copyright 2020 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
#include "lte/gateway/c/core/oai/test/s1ap_task/s1ap_mme_test_utils.h"

#include <cstdlib>
#include "S1ap_Cause.h"

extern "C" {
#include "lte/gateway/c/core/oai/lib/3gpp/3gpp_36.401.h"
#include "lte/gateway/c/core/oai/lib/hashtable/hashtable.h"
#include "lte/gateway/c/core/oai/common/common_types.h"
#include "lte/gateway/c/core/oai/common/conversions.h"
#include "lte/gateway/c/core/oai/common/itti_free_defined_msg.h"
#include "lte/gateway/c/core/oai/include/s1ap_messages_types.h"
#include "lte/gateway/c/core/oai/tasks/nas/ies/TrackingAreaIdentityList.h"
#include "lte/gateway/c/core/oai/lib/itti/intertask_interface.h"
}

#include "lte/gateway/c/core/oai/tasks/s1ap/s1ap_state_manager.h"

namespace magma {
namespace lte {

extern task_zmq_ctx_t task_zmq_ctx_main_s1ap;

status_code_e setup_new_association(
    s1ap_state_t* state, sctp_assoc_id_t assoc_id) {
  bstring ran_cp_ipaddr = bfromcstr("\xc0\xa8\x3c\x8d");
  sctp_new_peer_t p     = {
      .instreams     = 1,
      .outstreams    = 2,
      .assoc_id      = assoc_id,
      .ran_cp_ipaddr = ran_cp_ipaddr,
  };
  status_code_e rc = s1ap_handle_new_association(state, &p);
  bdestroy(ran_cp_ipaddr);
  return rc;
}

status_code_e send_s1ap_close_sctp_association(sctp_assoc_id_t assoc_id) {
  MessageDef* msg = itti_alloc_new_message(TASK_SCTP, SCTP_CLOSE_ASSOCIATION);

  SCTP_CLOSE_ASSOCIATION(msg).assoc_id = assoc_id;
  SCTP_CLOSE_ASSOCIATION(msg).reset    = false;

  return send_msg_to_task(&task_zmq_ctx_main_s1ap, TASK_S1AP, msg);
}

status_code_e generate_s1_setup_request_pdu(S1ap_S1AP_PDU_t* pdu_s1) {
  uint8_t packet_bytes[] = {
      0x00, 0x11, 0x00, 0x2f, 0x00, 0x00, 0x04, 0x00, 0x3b, 0x00, 0x09,
      0x00, 0x00, 0xf1, 0x10, 0x40, 0x00, 0x00, 0x00, 0xa0, 0x00, 0x3c,
      0x40, 0x0b, 0x80, 0x09, 0x22, 0x52, 0x41, 0x44, 0x49, 0x53, 0x59,
      0x53, 0x22, 0x00, 0x40, 0x00, 0x07, 0x00, 0x00, 0x00, 0x40, 0x00,
      0xf1, 0x10, 0x00, 0x89, 0x40, 0x01, 0x00};

  bstring payload_s1_setup;
  payload_s1_setup = blk2bstr(&packet_bytes, sizeof(packet_bytes));

  status_code_e pdu_rc = s1ap_mme_decode_pdu(pdu_s1, payload_s1_setup);
  bdestroy_wrapper(&payload_s1_setup);
  return pdu_rc;
}

void handle_mme_ue_id_notification(s1ap_state_t* s, sctp_assoc_id_t assoc_id) {
  MessageDef* message_p =
      itti_alloc_new_message(TASK_MME_APP, MME_APP_S1AP_MME_UE_ID_NOTIFICATION);
  itti_mme_app_s1ap_mme_ue_id_notification_t* notification_p =
      &message_p->ittiMsg.mme_app_s1ap_mme_ue_id_notification;
  memset(notification_p, 0, sizeof(itti_mme_app_s1ap_mme_ue_id_notification_t));
  notification_p->enb_ue_s1ap_id = 1;
  notification_p->mme_ue_s1ap_id = 7;
  notification_p->sctp_assoc_id  = assoc_id;
  s1ap_handle_mme_ue_id_notification(s, notification_p);
  free(message_p);
}

status_code_e send_s1ap_erab_rel_cmd(
    s1ap_state_t* state, mme_ue_s1ap_id_t ue_id, enb_ue_s1ap_id_t enb_ue_id) {
  MessageDef* message_p;
  message_p = itti_alloc_new_message(TASK_MME_APP, S1AP_E_RAB_REL_CMD);
  itti_s1ap_e_rab_rel_cmd_t* s1ap_e_rab_rel_cmd =
      &message_p->ittiMsg.s1ap_e_rab_rel_cmd;
  s1ap_e_rab_rel_cmd->mme_ue_s1ap_id = ue_id;
  s1ap_e_rab_rel_cmd->enb_ue_s1ap_id = enb_ue_id;

  s1ap_e_rab_rel_cmd->e_rab_to_be_rel_list.no_of_items      = 1;
  s1ap_e_rab_rel_cmd->e_rab_to_be_rel_list.item[0].e_rab_id = 5;

  status_code_e rc =
      s1ap_generate_s1ap_e_rab_rel_cmd(state, s1ap_e_rab_rel_cmd);

  itti_free_msg_content(message_p);
  free(message_p);
  return rc;
}

status_code_e send_conn_establishment_cnf(
    mme_ue_s1ap_id_t ue_id, bool sec_capabilities_present,
    bool ue_radio_capability) {
  MessageDef* message_p;
  message_p = itti_alloc_new_message(
      TASK_MME_APP, MME_APP_CONNECTION_ESTABLISHMENT_CNF);
  itti_mme_app_connection_establishment_cnf_t* establishment_cnf_p = NULL;
  establishment_cnf_p =
      &message_p->ittiMsg.mme_app_connection_establishment_cnf;
  establishment_cnf_p->ue_id        = ue_id;
  establishment_cnf_p->presencemask = 1;

  establishment_cnf_p->no_of_e_rabs = 1;

  establishment_cnf_p->e_rab_id[0] = 1;  //+ EPS_BEARER_IDENTITY_FIRST;
  establishment_cnf_p->e_rab_level_qos_qci[0]            = 1;
  establishment_cnf_p->e_rab_level_qos_priority_level[0] = 1;
  establishment_cnf_p->transport_layer_address[0]        = bfromcstr("test");
  establishment_cnf_p->gtp_teid[0]                       = 1;

  establishment_cnf_p->ue_ambr.br_ul = 1000;
  establishment_cnf_p->ue_ambr.br_dl = 1000;

  apn_ambr_bitrate_unit_t br_unit                                     = BPS;
  establishment_cnf_p->ue_ambr.br_unit                                = br_unit;
  establishment_cnf_p->ue_security_capabilities_encryption_algorithms = 1;
  establishment_cnf_p->ue_security_capabilities_integrity_algorithms  = 1;

  establishment_cnf_p->nr_ue_security_capabilities_present =
      sec_capabilities_present;
  if (ue_radio_capability) {
    establishment_cnf_p->ue_radio_capability = bfromcstr("test");
  }

  return send_msg_to_task(&task_zmq_ctx_main_s1ap, TASK_S1AP, message_p);
}

status_code_e send_s1ap_erab_setup_req(
    s1ap_state_t* state, mme_ue_s1ap_id_t ue_id, enb_ue_s1ap_id_t enb_ue_id,
    ebi_t ebi) {
  MessageDef* message_p =
      itti_alloc_new_message(TASK_MME_APP, S1AP_E_RAB_SETUP_REQ);
  itti_s1ap_e_rab_setup_req_t* s1ap_e_rab_setup_req =
      &message_p->ittiMsg.s1ap_e_rab_setup_req;

  s1ap_e_rab_setup_req->mme_ue_s1ap_id = ue_id;
  s1ap_e_rab_setup_req->enb_ue_s1ap_id = enb_ue_id;

  // E-RAB to Be Setup List
  s1ap_e_rab_setup_req->e_rab_to_be_setup_list.no_of_items      = 1;
  s1ap_e_rab_setup_req->e_rab_to_be_setup_list.item[0].e_rab_id = ebi;
  s1ap_e_rab_setup_req->e_rab_to_be_setup_list.item[0]
      .e_rab_level_qos_parameters.allocation_and_retention_priority
      .pre_emption_capability =
      (pre_emption_capability_t) PRE_EMPTION_CAPABILITY_ENABLED;
  s1ap_e_rab_setup_req->e_rab_to_be_setup_list.item[0]
      .e_rab_level_qos_parameters.allocation_and_retention_priority
      .pre_emption_vulnerability =
      (pre_emption_vulnerability_t) PRE_EMPTION_VULNERABILITY_ENABLED;
  s1ap_e_rab_setup_req->e_rab_to_be_setup_list.item[0]
      .e_rab_level_qos_parameters.allocation_and_retention_priority
      .priority_level = 9;
  s1ap_e_rab_setup_req->e_rab_to_be_setup_list.item[0]
      .e_rab_level_qos_parameters.gbr_qos_information
      .e_rab_maximum_bit_rate_downlink = 2000;
  s1ap_e_rab_setup_req->e_rab_to_be_setup_list.item[0]
      .e_rab_level_qos_parameters.gbr_qos_information
      .e_rab_maximum_bit_rate_uplink = 2000;
  s1ap_e_rab_setup_req->e_rab_to_be_setup_list.item[0]
      .e_rab_level_qos_parameters.gbr_qos_information
      .e_rab_guaranteed_bit_rate_downlink = 10000;
  s1ap_e_rab_setup_req->e_rab_to_be_setup_list.item[0]
      .e_rab_level_qos_parameters.gbr_qos_information
      .e_rab_guaranteed_bit_rate_uplink = 10000;
  s1ap_e_rab_setup_req->e_rab_to_be_setup_list.item[0]
      .e_rab_level_qos_parameters.qci = 1;

  s1ap_e_rab_setup_req->e_rab_to_be_setup_list.item[0].gtp_teid = 1;
  s1ap_e_rab_setup_req->e_rab_to_be_setup_list.item[0].transport_layer_address =
      bfromcstr("127.0.0.1");

  s1ap_e_rab_setup_req->e_rab_to_be_setup_list.item[0].nas_pdu =
      bfromcstr("test");
  status_code_e rc =
      s1ap_generate_s1ap_e_rab_setup_req(state, s1ap_e_rab_setup_req);

  itti_free_msg_content(message_p);
  free(message_p);
  return rc;
}

status_code_e send_s1ap_erab_reset_req(
    sctp_assoc_id_t assoc_id, sctp_stream_id_t stream_id,
    enb_ue_s1ap_id_t enb_ue_id, mme_ue_s1ap_id_t ue_id) {
  MessageDef* msg = DEPRECATEDitti_alloc_new_message_fatal(
      TASK_MME_APP, S1AP_ENB_INITIATED_RESET_ACK);

  itti_s1ap_enb_initiated_reset_ack_t* reset_ack =
      &msg->ittiMsg.s1ap_enb_initiated_reset_ack;

  s1_sig_conn_id_t* list =
      (s1_sig_conn_id_t*) (calloc(1, sizeof(s1_sig_conn_id_t)));
  list->enb_ue_s1ap_id = enb_ue_id;
  list->mme_ue_s1ap_id = ue_id;
  // ue_to_reset_list needs to be freed by S1AP module
  reset_ack->ue_to_reset_list = list;
  reset_ack->s1ap_reset_type  = RESET_PARTIAL;
  reset_ack->sctp_assoc_id    = assoc_id;
  reset_ack->sctp_stream_id   = stream_id;
  reset_ack->num_ue           = 1;

  // Send Reset Ack to S1AP module
  return send_msg_to_task(&task_zmq_ctx_main_s1ap, TASK_S1AP, msg);
}

status_code_e send_s1ap_ue_ctxt_mod(
    enb_ue_s1ap_id_t enb_ue_id, mme_ue_s1ap_id_t ue_id) {
  MessageDef* message_p = itti_alloc_new_message(
      TASK_MME_APP, S1AP_UE_CONTEXT_MODIFICATION_REQUEST);
  memset(
      (void*) &message_p->ittiMsg.s1ap_ue_context_mod_request, 0,
      sizeof(itti_s1ap_ue_context_mod_req_t));
  S1AP_UE_CONTEXT_MODIFICATION_REQUEST(message_p).mme_ue_s1ap_id = ue_id;
  S1AP_UE_CONTEXT_MODIFICATION_REQUEST(message_p).enb_ue_s1ap_id = enb_ue_id;
  S1AP_UE_CONTEXT_MODIFICATION_REQUEST(message_p).presencemask =
      S1AP_UE_CONTEXT_MOD_UE_AMBR_INDICATOR_PRESENT;
  S1AP_UE_CONTEXT_MODIFICATION_REQUEST(message_p).ue_ambr.br_ul = 10000;
  S1AP_UE_CONTEXT_MODIFICATION_REQUEST(message_p).ue_ambr.br_dl = 1000;

  return send_msg_to_task(&task_zmq_ctx_main_s1ap, TASK_S1AP, message_p);
}

status_code_e send_s1ap_paging_request(sctp_assoc_id_t assoc_id) {
  MessageDef* message_p =
      itti_alloc_new_message(TASK_MME_APP, S1AP_PAGING_REQUEST);

  itti_s1ap_paging_request_t* paging_request =
      &message_p->ittiMsg.s1ap_paging_request;
  memset(paging_request, 0, sizeof(itti_s1ap_paging_request_t));

  strncpy(paging_request->imsi, std::string("IMSI00000000001002").c_str(), 15);
  paging_request->imsi_length      = 15;
  paging_request->mme_code         = 1;
  paging_request->m_tmsi           = 1;
  paging_request->sctp_assoc_id    = assoc_id;
  paging_request->paging_id        = S1AP_PAGING_ID_IMSI;
  paging_request->domain_indicator = CN_DOMAIN_CS;

  // Send TAI List
  paging_request->tai_list_count                                 = 1;
  paging_request->paging_tai_list[0].numoftac                    = 1;
  paging_request->paging_tai_list[0].tai_list[0].plmn.mcc_digit1 = 1;
  paging_request->paging_tai_list[0].tai_list[0].plmn.mcc_digit2 = 0;
  paging_request->paging_tai_list[0].tai_list[0].plmn.mcc_digit3 = 0;
  paging_request->paging_tai_list[0].tai_list[0].plmn.mnc_digit1 = 1;
  paging_request->paging_tai_list[0].tai_list[0].plmn.mnc_digit2 = 0;
  paging_request->paging_tai_list[0].tai_list[0].plmn.mnc_digit3 = 1;
  paging_request->paging_tai_list[0].tai_list[0].tac             = 1;

  return send_msg_to_task(&task_zmq_ctx_main_s1ap, TASK_S1AP, message_p);
}

status_code_e send_s1ap_path_switch_failure(
    sctp_assoc_id_t assoc_id, enb_ue_s1ap_id_t enb_ue_id,
    mme_ue_s1ap_id_t ue_id) {
  MessageDef* message_p =
      itti_alloc_new_message(TASK_MME_APP, S1AP_PATH_SWITCH_REQUEST_FAILURE);
  itti_s1ap_path_switch_request_failure_t* s1ap_path_switch_req_failure =
      &message_p->ittiMsg.s1ap_path_switch_request_failure;

  s1ap_path_switch_req_failure->sctp_assoc_id  = assoc_id;
  s1ap_path_switch_req_failure->enb_ue_s1ap_id = enb_ue_id;
  s1ap_path_switch_req_failure->mme_ue_s1ap_id = ue_id;

  return send_msg_to_task(&task_zmq_ctx_main_s1ap, TASK_S1AP, message_p);
}

status_code_e send_s1ap_path_switch_req(
    sctp_assoc_id_t assoc_id, enb_ue_s1ap_id_t enb_ue_id,
    mme_ue_s1ap_id_t ue_id) {
  MessageDef* message_p =
      itti_alloc_new_message(TASK_MME_APP, S1AP_PATH_SWITCH_REQUEST_ACK);

  itti_s1ap_path_switch_request_ack_t* s1ap_path_switch_req_ack =
      &message_p->ittiMsg.s1ap_path_switch_request_ack;

  s1ap_path_switch_req_ack->sctp_assoc_id  = assoc_id;
  s1ap_path_switch_req_ack->enb_ue_s1ap_id = enb_ue_id;
  s1ap_path_switch_req_ack->mme_ue_s1ap_id = ue_id;
  for (auto i = 0; i < AUTH_NEXT_HOP_SIZE; i++) {
    s1ap_path_switch_req_ack->nh[i] = 1;
  }
  s1ap_path_switch_req_ack->ncc = 1;

  return send_msg_to_task(&task_zmq_ctx_main_s1ap, TASK_S1AP, message_p);
}

status_code_e send_s1ap_mme_handover_request(
    sctp_assoc_id_t assoc_id, mme_ue_s1ap_id_t ue_id, uint32_t target_enb_id) {
  MessageDef* message_p =
      itti_alloc_new_message(TASK_MME_APP, MME_APP_HANDOVER_REQUEST);
  itti_mme_app_handover_request_t* ho_request_p =
      &message_p->ittiMsg.mme_app_handover_request;

  char arbitrary_src_tgt_container[20]            = "Arbitrary target";
  ho_request_p->encryption_algorithm_capabilities = 1;
  ho_request_p->integrity_algorithm_capabilities  = 2;
  ho_request_p->mme_ue_s1ap_id                    = ue_id;
  ho_request_p->target_sctp_assoc_id              = assoc_id;
  ho_request_p->target_enb_id                     = target_enb_id;
  ho_request_p->cause.present                     = S1ap_Cause_PR_radioNetwork;
  ho_request_p->cause.choice.radioNetwork =
      S1ap_CauseRadioNetwork_handover_desirable_for_radio_reason;
  ho_request_p->handover_type     = S1ap_HandoverType_intralte;
  ho_request_p->src_tgt_container = blk2bstr(arbitrary_src_tgt_container, 10);
  ho_request_p->ue_ambr.br_unit   = KBPS;
  ho_request_p->ue_ambr.br_ul     = 1000;
  ho_request_p->ue_ambr.br_dl     = 10000;
  ho_request_p->e_rab_list.no_of_items = 2;
  fteid_t s_gw_fteid_s1u               = {1};

  for (int i = 0; i < ho_request_p->e_rab_list.no_of_items; ++i) {
    ho_request_p->e_rab_list.item[i].e_rab_id = 1;
    ho_request_p->e_rab_list.item[i].transport_layer_address =
        fteid_ip_address_to_bstring(&s_gw_fteid_s1u);
    ho_request_p->e_rab_list.item[i].gtp_teid                       = 1;
    ho_request_p->e_rab_list.item[i].e_rab_level_qos_parameters.qci = 9;
    ho_request_p->e_rab_list.item[i]
        .e_rab_level_qos_parameters.allocation_and_retention_priority
        .priority_level = 0;
    ho_request_p->e_rab_list.item[i]
        .e_rab_level_qos_parameters.allocation_and_retention_priority
        .pre_emption_capability =
        (pre_emption_capability_t) PRE_EMPTION_CAPABILITY_ENABLED;
    ho_request_p->e_rab_list.item[i]
        .e_rab_level_qos_parameters.allocation_and_retention_priority
        .pre_emption_vulnerability =
        (pre_emption_vulnerability_t) PRE_EMPTION_VULNERABILITY_DISABLED;
  }
  for (int i = 0; i < AUTH_NEXT_HOP_SIZE; ++i) {
    ho_request_p->nh[i] = 0x11;
  }
  ho_request_p->ncc = 2;
  return send_msg_to_task(&task_zmq_ctx_main_s1ap, TASK_S1AP, message_p);
}

status_code_e send_s1ap_mme_handover_command(
    sctp_assoc_id_t assoc_id, mme_ue_s1ap_id_t ue_id,
    enb_ue_s1ap_id_t src_enb_ue_id, enb_ue_s1ap_id_t tgt_enb_ue_id,
    uint32_t source_enb_id, uint32_t target_enb_id) {
  char arbitrary_src_tgt_container[20] = "Arbitrary target";
  MessageDef* message_p =
      itti_alloc_new_message(TASK_MME_APP, MME_APP_HANDOVER_COMMAND);
  itti_mme_app_handover_command_t* ho_command_p =
      &message_p->ittiMsg.mme_app_handover_command;

  ho_command_p->source_assoc_id    = assoc_id;
  ho_command_p->mme_ue_s1ap_id     = ue_id;
  ho_command_p->src_enb_ue_s1ap_id = src_enb_ue_id;
  ho_command_p->tgt_enb_ue_s1ap_id = tgt_enb_ue_id;
  ho_command_p->source_enb_id      = source_enb_id;
  ho_command_p->target_enb_id      = target_enb_id;
  ho_command_p->handover_type      = S1ap_HandoverType_intralte;
  ho_command_p->tgt_src_container  = blk2bstr(arbitrary_src_tgt_container, 10);

  return send_msg_to_task(&task_zmq_ctx_main_s1ap, TASK_S1AP, message_p);
}

status_code_e send_s1ap_erab_mod_confirm(
    enb_ue_s1ap_id_t enb_ue_id, mme_ue_s1ap_id_t ue_id) {
  MessageDef* message_p =
      itti_alloc_new_message(TASK_MME_APP, S1AP_E_RAB_MODIFICATION_CNF);

  itti_s1ap_e_rab_modification_cnf_t* s1ap_e_rab_modification_cnf_p =
      &message_p->ittiMsg.s1ap_e_rab_modification_cnf;

  /** Set the identifiers. */
  s1ap_e_rab_modification_cnf_p->mme_ue_s1ap_id = ue_id;
  s1ap_e_rab_modification_cnf_p->enb_ue_s1ap_id = enb_ue_id;

  s1ap_e_rab_modification_cnf_p->e_rab_modify_list.e_rab_id[0] = 1;
  s1ap_e_rab_modification_cnf_p->e_rab_modify_list.no_of_items = 1;
  s1ap_e_rab_modification_cnf_p->e_rab_failed_to_modify_list.item[0].e_rab_id =
      1;
  s1ap_e_rab_modification_cnf_p->e_rab_failed_to_modify_list.item[0]
      .cause.present = S1ap_Cause_PR_misc;
  s1ap_e_rab_modification_cnf_p->e_rab_failed_to_modify_list.item[0]
      .cause.present = S1ap_Cause_PR_misc;
  s1ap_e_rab_modification_cnf_p->e_rab_failed_to_modify_list.no_of_items = 1;

  return send_msg_to_task(&task_zmq_ctx_main_s1ap, TASK_S1AP, message_p);
}

bool is_enb_state_valid(
    s1ap_state_t* state, sctp_assoc_id_t assoc_id,
    mme_s1_enb_state_s expected_state, uint32_t expected_num_ues) {
  enb_description_t* enb_associated = nullptr;
  hashtable_ts_get(
      &state->enbs, (const hash_key_t) assoc_id,
      reinterpret_cast<void**>(&enb_associated));
  if (enb_associated->nb_ue_associated == expected_num_ues &&
      enb_associated->s1_state == expected_state) {
    return true;
  }
  return false;
}

bool is_num_enbs_valid(s1ap_state_t* state, uint32_t expected_num_enbs) {
  hash_size_t num_enb_elements = state->enbs.num_elements;
  if ((num_enb_elements == expected_num_enbs) &&
      (state->num_enbs == expected_num_enbs)) {
    return true;
  }
  return false;
}

bool is_ue_state_valid(
    sctp_assoc_id_t assoc_id, enb_ue_s1ap_id_t enb_ue_id,
    enum s1_ue_state_s expected_ue_state) {
  ue_description_t* ue   = nullptr;
  hash_table_ts_t* ue_ht = S1apStateManager::getInstance().get_ue_state_ht();
  uint64_t comp_s1ap_id  = S1AP_GENERATE_COMP_S1AP_ID(assoc_id, enb_ue_id);
  hashtable_rc_t ht_rc   = hashtable_ts_get(
      ue_ht, (const hash_key_t) comp_s1ap_id, reinterpret_cast<void**>(&ue));
  if (ht_rc != HASH_TABLE_OK) {
    return false;
  }
  return ue->s1_ue_state == expected_ue_state ? true : false;
}

status_code_e simulate_pdu_s1_message(
    uint8_t* bytes, long bytes_len, s1ap_state_t* state,
    sctp_assoc_id_t assoc_id, sctp_stream_id_t stream_id) {
  status_code_e rc;
  bstring payload;
  payload = blk2bstr(bytes, bytes_len);
  S1ap_S1AP_PDU_t pdu;
  memset(&pdu, 0, sizeof(pdu));

  rc = s1ap_mme_decode_pdu(&pdu, payload);
  bdestroy_wrapper(&payload);
  rc = s1ap_mme_handle_message(state, assoc_id, stream_id, &pdu);
  ASN_STRUCT_FREE_CONTENTS_ONLY(asn_DEF_S1ap_S1AP_PDU, &pdu);
  return rc;
}

}  // namespace lte
}  // namespace magma

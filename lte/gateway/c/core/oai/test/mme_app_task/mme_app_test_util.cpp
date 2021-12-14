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
#include "lte/gateway/c/core/oai/test/mme_app_task/mme_app_test_util.h"

#include <chrono>
#include <gtest/gtest.h>
#include <cstdint>
#include <thread>

#include "feg/protos/s6a_proxy.pb.h"
#include "lte/gateway/c/core/oai/lib/s6a_proxy/proto_msg_to_itti_msg.h"

extern "C" {
#include "lte/gateway/c/core/oai/lib/itti/intertask_interface.h"
#include "lte/gateway/c/core/oai/lib/3gpp/3gpp_23.003.h"
#include "lte/gateway/c/core/oai/lib/bstr/bstrlib.h"
}

namespace magma {
namespace lte {

extern task_zmq_ctx_t task_zmq_ctx_main;

#define DEFAULT_LBI 5
#define DEFAULT_TEID 1
#define DEFAULT_MME_S1AP_UE_ID 1

#define DEFAULT_UE_IPv4 1000

void nas_config_timer_reinit(nas_config_t* nas_conf, uint32_t timeout_msec) {
  nas_conf->t3402_min = 1;
  nas_conf->t3412_min = 1;
  nas_conf->t3412_msec =
      50 * timeout_msec;  // implicit detach after 2x of this value
  nas_conf->t3422_msec   = timeout_msec;
  nas_conf->t3450_msec   = timeout_msec;
  nas_conf->t3460_msec   = timeout_msec;
  nas_conf->t3470_msec   = timeout_msec;
  nas_conf->t3485_msec   = timeout_msec;
  nas_conf->t3486_msec   = timeout_msec;
  nas_conf->t3489_msec   = timeout_msec;
  nas_conf->t3495_msec   = timeout_msec;
  nas_conf->ts6a_msec    = timeout_msec;
  nas_conf->tics_msec    = timeout_msec;
  nas_conf->tpaging_msec = timeout_msec;
  return;
}

void send_sctp_mme_server_initialized() {
  MessageDef* message_p =
      itti_alloc_new_message(TASK_S1AP, SCTP_MME_SERVER_INITIALIZED);
  SCTP_MME_SERVER_INITIALIZED(message_p).successful = true;
  send_msg_to_task(&task_zmq_ctx_main, TASK_MME_APP, message_p);
  return;
}

void send_activate_message_to_mme_app() {
  MessageDef* message_p = itti_alloc_new_message(TASK_MAIN, ACTIVATE_MESSAGE);
  send_msg_to_task(&task_zmq_ctx_main, TASK_MME_APP, message_p);
  return;
}

void send_mme_app_initial_ue_msg(
    const uint8_t* nas_msg, uint8_t nas_msg_length, const plmn_t& plmn,
    guti_eps_mobile_identity_t& guti, tac_t tac) {
  MessageDef* message_p =
      itti_alloc_new_message(TASK_S1AP, S1AP_INITIAL_UE_MESSAGE);
  ITTI_MSG_LASTHOP_LATENCY(message_p)               = 0;
  S1AP_INITIAL_UE_MESSAGE(message_p).sctp_assoc_id  = DEFAULT_SCTP_ASSOC_ID;
  S1AP_INITIAL_UE_MESSAGE(message_p).enb_ue_s1ap_id = DEFAULT_eNB_S1AP_UE_ID;
  S1AP_INITIAL_UE_MESSAGE(message_p).enb_id         = DEFAULT_ENB_ID;
  S1AP_INITIAL_UE_MESSAGE(message_p).nas = blk2bstr(nas_msg, nas_msg_length);
  S1AP_INITIAL_UE_MESSAGE(message_p).tai.plmn           = plmn;
  S1AP_INITIAL_UE_MESSAGE(message_p).tai.tac            = tac;
  S1AP_INITIAL_UE_MESSAGE(message_p).ecgi.plmn          = plmn;
  S1AP_INITIAL_UE_MESSAGE(message_p).ecgi.cell_identity = {0, 0, 0};
  if (guti.m_tmsi) {
    S1AP_INITIAL_UE_MESSAGE(message_p).is_s_tmsi_valid     = true;
    S1AP_INITIAL_UE_MESSAGE(message_p).opt_s_tmsi.m_tmsi   = guti.m_tmsi;
    S1AP_INITIAL_UE_MESSAGE(message_p).opt_s_tmsi.mme_code = guti.mme_code;
  }
  send_msg_to_task(&task_zmq_ctx_main, TASK_MME_APP, message_p);
  return;
}

void send_mme_app_uplink_data_ind(
    const uint8_t* nas_msg, uint8_t nas_msg_length, const plmn_t& plmn) {
  MessageDef* message_p =
      itti_alloc_new_message(TASK_S1AP, MME_APP_UPLINK_DATA_IND);
  ITTI_MSG_LASTHOP_LATENCY(message_p)     = 0;
  MME_APP_UL_DATA_IND(message_p).ue_id    = 1;
  MME_APP_UL_DATA_IND(message_p).nas_msg  = blk2bstr(nas_msg, nas_msg_length);
  MME_APP_UL_DATA_IND(message_p).tai.plmn = plmn;
  MME_APP_UL_DATA_IND(message_p).tai.tac  = 1;
  MME_APP_UL_DATA_IND(message_p).cgi.plmn = plmn;
  MME_APP_UL_DATA_IND(message_p).cgi.cell_identity = {0, 0, 0};
  send_msg_to_task(&task_zmq_ctx_main, TASK_MME_APP, message_p);
  return;
}

void send_authentication_info_resp(const std::string& imsi, bool success) {
  MessageDef* message_p = itti_alloc_new_message(TASK_S6A, S6A_AUTH_INFO_ANS);
  s6a_auth_info_ans_t* itti_msg = &message_p->ittiMsg.s6a_auth_info_ans;
  strncpy(itti_msg->imsi, imsi.c_str(), imsi.size());
  itti_msg->imsi_length    = imsi.size();
  itti_msg->result.present = S6A_RESULT_BASE;
  if (success) {
    itti_msg->result.choice.base = DIAMETER_SUCCESS;
    magma::feg::AuthenticationInformationAnswer aia;
    magma::feg::AuthenticationInformationAnswer::EUTRANVector eutran_vector;
    uint8_t xres_buf[XRES_LENGTH_MAX]    = {0x66, 0xff, 0x47, 0x2d, 0xd4, 0x93,
                                         0xf1, 0x5a, 0x00, 0x00, 0x00, 0x00,
                                         0x00, 0x00, 0x00, 0x00};
    uint8_t rand_buf[RAND_LENGTH_OCTETS] = {0x68, 0x16, 0xa1, 0x0c, 0x0f, 0xeb,
                                            0x44, 0xa5, 0x00, 0x5c, 0x9c, 0x9c,
                                            0x3c, 0x6f, 0xd6, 0x15};
    uint8_t autn_buf[AUTN_LENGTH_OCTETS] = {0x4a, 0xe4, 0xe0, 0xd9, 0xaa, 0x4b,
                                            0x80, 0x00, 0xc4, 0x80, 0xa1, 0x97,
                                            0x70, 0x4b, 0x7b, 0x8f};
    uint8_t kasme_buf[KASME_LENGTH_OCTETS] = {
        0xc3, 0x5f, 0x03, 0x8f, 0x5f, 0xbe, 0xcc, 0x23, 0xc4, 0xd1, 0xa7,
        0xd6, 0x8a, 0xf7, 0x05, 0x32, 0xf2, 0x37, 0xf6, 0x40, 0x47, 0xdd,
        0x29, 0x6e, 0x7d, 0x0e, 0xf6, 0xe9, 0x26, 0x5f, 0x24, 0x39};
    eutran_vector.set_rand((const void*) rand_buf, RAND_LENGTH_OCTETS);
    eutran_vector.set_xres((const void*) xres_buf, XRES_LENGTH_MAX);
    eutran_vector.set_autn((const void*) autn_buf, AUTN_LENGTH_OCTETS);
    eutran_vector.set_kasme((const void*) kasme_buf, KASME_LENGTH_OCTETS);
    aia.set_error_code(magma::feg::ErrorCode::SUCCESS);
    auto eutran_vectors = aia.mutable_eutran_vectors();
    eutran_vectors->Add()->CopyFrom(eutran_vector);
    magma::convert_proto_msg_to_itti_s6a_auth_info_ans(aia, itti_msg);
  } else {
    itti_msg->result.choice.base = DIAMETER_UNABLE_TO_COMPLY;
  }
  send_msg_to_task(&task_zmq_ctx_main, TASK_MME_APP, message_p);
  return;
}

void send_s6a_ula(const std::string& imsi, bool success) {
  MessageDef* message_p =
      itti_alloc_new_message(TASK_S6A, S6A_UPDATE_LOCATION_ANS);
  s6a_update_location_ans_t* itti_msg =
      &message_p->ittiMsg.s6a_update_location_ans;
  strncpy(itti_msg->imsi, imsi.c_str(), imsi.size());
  itti_msg->imsi_length    = imsi.size();
  itti_msg->result.present = S6A_RESULT_BASE;
  if (success) {
    itti_msg->result.choice.base = DIAMETER_SUCCESS;
    magma::feg::UpdateLocationAnswer ula;
    ula.set_default_context_id(0);
    auto total_ambr = ula.mutable_total_ambr();
    total_ambr->set_max_bandwidth_ul(100000000);
    total_ambr->set_max_bandwidth_dl(200000000);
    ula.set_all_apns_included(false);
    magma::feg::UpdateLocationAnswer::APNConfiguration apnconfig;
    apnconfig.set_context_id(0);
    apnconfig.set_service_selection("magma.ipv4");
    auto apn_qosprofile = apnconfig.mutable_qos_profile();
    apn_qosprofile->set_class_id(9);
    apn_qosprofile->set_priority_level(15);
    auto apn_ambr = apnconfig.mutable_ambr();
    apn_ambr->set_max_bandwidth_ul(10000000);
    apn_ambr->set_max_bandwidth_dl(75000000);
    apnconfig.set_pdn(magma::feg::UpdateLocationAnswer::APNConfiguration::IPV4);
    auto apns = ula.mutable_apn();
    apns->Add()->CopyFrom(apnconfig);
    apnconfig.set_service_selection("ims");
    apns->Add()->CopyFrom(apnconfig);
    magma::convert_proto_msg_to_itti_s6a_update_location_ans(ula, itti_msg);
  } else {
    itti_msg->result.choice.base = DIAMETER_UNABLE_TO_COMPLY;
  }
  send_msg_to_task(&task_zmq_ctx_main, TASK_MME_APP, message_p);
  return;
}

void send_create_session_resp(gtpv2c_cause_value_t cause_value) {
  MessageDef* message_p =
      itti_alloc_new_message(TASK_SPGW_APP, S11_CREATE_SESSION_RESPONSE);
  itti_s11_create_session_response_t* create_session_response_p =
      &message_p->ittiMsg.s11_create_session_response;

  create_session_response_p->teid              = 1;
  create_session_response_p->cause.cause_value = cause_value;
  create_session_response_p->bearer_contexts_created.bearer_contexts[0]
      .cause.cause_value = cause_value;
  create_session_response_p->bearer_contexts_created.num_bearer_context = 1;

  if (cause_value == REQUEST_ACCEPTED) {
    create_session_response_p->paa.pdn_type            = IPv4;
    create_session_response_p->paa.ipv4_address.s_addr = DEFAULT_UE_IPv4;
    create_session_response_p->bearer_contexts_created.bearer_contexts[0]
        .s1u_sgw_fteid.teid = 1000;
    create_session_response_p->bearer_contexts_created.bearer_contexts[0]
        .s1u_sgw_fteid.interface_type = S1_U_SGW_GTP_U;
    create_session_response_p->bearer_contexts_created.bearer_contexts[0]
        .s1u_sgw_fteid.ipv4 = 1;
    create_session_response_p->bearer_contexts_created.bearer_contexts[0]
        .s1u_sgw_fteid.ipv4_address.s_addr = 100;
    create_session_response_p->bearer_contexts_created.bearer_contexts[0]
        .eps_bearer_id = 5;
    create_session_response_p->bearer_contexts_created.bearer_contexts[0]
        .s1u_sgw_fteid.ipv6 = 1;
  }

  send_msg_to_task(&task_zmq_ctx_main, TASK_MME_APP, message_p);
  return;
}

void send_delete_session_resp() {
  MessageDef* message_p =
      itti_alloc_new_message(TASK_SPGW_APP, S11_DELETE_SESSION_RESPONSE);
  itti_s11_delete_session_response_t* delete_session_resp_p =
      &message_p->ittiMsg.s11_delete_session_response;
  delete_session_resp_p->cause.cause_value = REQUEST_ACCEPTED;
  delete_session_resp_p->teid              = 1;
  delete_session_resp_p->peer_ip.s_addr    = 100;
  delete_session_resp_p->lbi               = 5;
  send_msg_to_task(&task_zmq_ctx_main, TASK_MME_APP, message_p);
  return;
}

void send_ics_response() {
  MessageDef* message_p =
      itti_alloc_new_message(TASK_S1AP, MME_APP_INITIAL_CONTEXT_SETUP_RSP);
  MME_APP_INITIAL_CONTEXT_SETUP_RSP(message_p).ue_id                        = 1;
  MME_APP_INITIAL_CONTEXT_SETUP_RSP(message_p).e_rab_setup_list.no_of_items = 1;
  MME_APP_INITIAL_CONTEXT_SETUP_RSP(message_p)
      .e_rab_setup_list.item[0]
      .e_rab_id = 5;
  MME_APP_INITIAL_CONTEXT_SETUP_RSP(message_p)
      .e_rab_setup_list.item[0]
      .gtp_teid                     = 0;
  uint8_t transport_address_buff[4] = {192, 168, 60, 141};
  MME_APP_INITIAL_CONTEXT_SETUP_RSP(message_p)
      .e_rab_setup_list.item[0]
      .transport_layer_address = blk2bstr(transport_address_buff, 4);
  send_msg_to_task(&task_zmq_ctx_main, TASK_MME_APP, message_p);
  return;
}

void send_ics_failure() {
  MessageDef* message_p =
      itti_alloc_new_message(TASK_S1AP, MME_APP_INITIAL_CONTEXT_SETUP_FAILURE);
  MME_APP_INITIAL_CONTEXT_SETUP_FAILURE(message_p).mme_ue_s1ap_id = 1;
  send_msg_to_task(&task_zmq_ctx_main, TASK_MME_APP, message_p);
  return;
}

void send_ue_ctx_release_complete() {
  MessageDef* message_p =
      itti_alloc_new_message(TASK_S1AP, S1AP_UE_CONTEXT_RELEASE_COMPLETE);
  S1AP_UE_CONTEXT_RELEASE_COMPLETE(message_p).mme_ue_s1ap_id = 1;
  send_msg_to_task(&task_zmq_ctx_main, TASK_MME_APP, message_p);
  return;
}

void send_ue_capabilities_ind() {
  MessageDef* message_p =
      itti_alloc_new_message(TASK_S1AP, S1AP_UE_CAPABILITIES_IND);
  itti_s1ap_ue_cap_ind_t* ue_cap_ind_p    = &message_p->ittiMsg.s1ap_ue_cap_ind;
  ue_cap_ind_p->enb_ue_s1ap_id            = 0;
  ue_cap_ind_p->mme_ue_s1ap_id            = 1;
  ue_cap_ind_p->radio_capabilities_length = 200;
  // using malloc to create uninitialized buffer
  ue_cap_ind_p->radio_capabilities =
      (uint8_t*) malloc(ue_cap_ind_p->radio_capabilities_length);
  send_msg_to_task(&task_zmq_ctx_main, TASK_MME_APP, message_p);
  return;
}

void send_context_release_req(s1cause rel_cause, task_id_t TASK_ID) {
  MessageDef* message_p =
      itti_alloc_new_message(TASK_ID, S1AP_UE_CONTEXT_RELEASE_REQ);
  S1AP_UE_CONTEXT_RELEASE_REQ(message_p).mme_ue_s1ap_id =
      DEFAULT_MME_S1AP_UE_ID;
  S1AP_UE_CONTEXT_RELEASE_REQ(message_p).enb_ue_s1ap_id =
      DEFAULT_eNB_S1AP_UE_ID;
  S1AP_UE_CONTEXT_RELEASE_REQ(message_p).enb_id   = DEFAULT_ENB_ID;
  S1AP_UE_CONTEXT_RELEASE_REQ(message_p).relCause = rel_cause;
  send_msg_to_task(&task_zmq_ctx_main, TASK_MME_APP, message_p);
  return;
}

void send_modify_bearer_resp(
    const std::vector<int>& bearer_to_modify,
    const std::vector<int>& bearer_to_remove) {
  MessageDef* message_p =
      itti_alloc_new_message(TASK_SPGW_APP, S11_MODIFY_BEARER_RESPONSE);
  itti_s11_modify_bearer_response_t* modify_response_p =
      &message_p->ittiMsg.s11_modify_bearer_response;
  modify_response_p->teid              = 1;
  modify_response_p->cause.cause_value = REQUEST_ACCEPTED;
  for (int i = 0; i < bearer_to_modify.size(); ++i) {
    modify_response_p->bearer_contexts_modified.bearer_contexts[i]
        .eps_bearer_id = bearer_to_modify[i];
    modify_response_p->bearer_contexts_modified.bearer_contexts[i]
        .cause.cause_value = REQUEST_ACCEPTED;
  }
  modify_response_p->bearer_contexts_modified.num_bearer_context =
      bearer_to_modify.size();
  for (int i = 0; i < bearer_to_remove.size(); ++i) {
    modify_response_p->bearer_contexts_marked_for_removal.bearer_contexts[i]
        .eps_bearer_id = bearer_to_remove[i];
    modify_response_p->bearer_contexts_marked_for_removal.bearer_contexts[i]
        .cause.cause_value = REQUEST_ACCEPTED;
  }
  modify_response_p->bearer_contexts_marked_for_removal.num_bearer_context =
      bearer_to_modify.size();
  send_msg_to_task(&task_zmq_ctx_main, TASK_MME_APP, message_p);
  return;
}

void sgw_send_release_access_bearer_response(gtpv2c_cause_value_t cause) {
  MessageDef* message_p = itti_alloc_new_message(
      TASK_SPGW_APP, S11_RELEASE_ACCESS_BEARERS_RESPONSE);
  itti_s11_release_access_bearers_response_t* release_access_bearers_resp_p =
      &message_p->ittiMsg.s11_release_access_bearers_response;
  release_access_bearers_resp_p->cause.cause_value = cause;
  release_access_bearers_resp_p->teid              = 1;
  send_msg_to_task(&task_zmq_ctx_main, TASK_MME_APP, message_p);
  return;
}

void send_s11_deactivate_bearer_req(
    uint8_t no_of_bearers_to_be_deact, uint8_t* ebi_to_be_deactivated,
    bool delete_default_bearer) {
  MessageDef* message_p = itti_alloc_new_message(
      TASK_SPGW_APP, S11_NW_INITIATED_DEACTIVATE_BEARER_REQUEST);
  itti_s11_nw_init_deactv_bearer_request_t* s11_bearer_deactv_request =
      &message_p->ittiMsg.s11_nw_init_deactv_bearer_request;

  s11_bearer_deactv_request->s11_mme_teid          = DEFAULT_TEID;
  s11_bearer_deactv_request->delete_default_bearer = delete_default_bearer;
  s11_bearer_deactv_request->no_of_bearers         = no_of_bearers_to_be_deact;
  memcpy(
      s11_bearer_deactv_request->ebi, ebi_to_be_deactivated,
      (sizeof(ebi_t) * no_of_bearers_to_be_deact));
  send_msg_to_task(&task_zmq_ctx_main, TASK_MME_APP, message_p);
  return;
}

void send_s11_create_bearer_req() {
  MessageDef* message_p = itti_alloc_new_message(
      TASK_SPGW_APP, S11_NW_INITIATED_ACTIVATE_BEARER_REQUEST);
  itti_s11_nw_init_actv_bearer_request_t* s11_actv_bearer_request =
      &message_p->ittiMsg.s11_nw_init_actv_bearer_request;
  s11_actv_bearer_request->s11_mme_teid = DEFAULT_TEID;
  s11_actv_bearer_request->lbi          = DEFAULT_LBI;
  s11_actv_bearer_request->eps_bearer_qos.gbr.br_dl =
      10000;  // arbitrary number
  s11_actv_bearer_request->eps_bearer_qos.gbr.br_ul =
      10000;  // arbitrary number
  s11_actv_bearer_request->eps_bearer_qos.mbr.br_dl =
      100000;  // arbitrary number
  s11_actv_bearer_request->eps_bearer_qos.mbr.br_ul =
      100000;                                        // arbitrary number
  s11_actv_bearer_request->eps_bearer_qos.pci = 1;   // 0 or 1
  s11_actv_bearer_request->eps_bearer_qos.pl  = 10;  // arbitrary number
  s11_actv_bearer_request->eps_bearer_qos.pvi = 0;   // 0 or 1
  s11_actv_bearer_request->context_teid       = DEFAULT_TEID;
  s11_actv_bearer_request->tft.ebit =
      TRAFFIC_FLOW_TEMPLATE_PARAMETER_LIST_IS_NOT_INCLUDED;
  s11_actv_bearer_request->tft.numberofpacketfilters = 1;
  s11_actv_bearer_request->tft.packetfilterlist.createnewtft[0].direction =
      TRAFFIC_FLOW_TEMPLATE_UPLINK_ONLY;
  s11_actv_bearer_request->tft.packetfilterlist.createnewtft[0]
      .eval_precedence = 250;
  s11_actv_bearer_request->tft.packetfilterlist.createnewtft[0]
      .packetfiltercontents.protocolidentifier_nextheader = 6;
  s11_actv_bearer_request->tft.packetfilterlist.createnewtft[0]
      .packetfiltercontents.flags |=
      (TRAFFIC_FLOW_TEMPLATE_IPV4_REMOTE_ADDR_FLAG |
       TRAFFIC_FLOW_TEMPLATE_SINGLE_REMOTE_PORT_FLAG);
  for (int i = 0; i < TRAFFIC_FLOW_TEMPLATE_IPV4_ADDR_SIZE; ++i) {
    s11_actv_bearer_request->tft.packetfilterlist.createnewtft[0]
        .packetfiltercontents.ipv4remoteaddr[i]
        .addr = 8;
    s11_actv_bearer_request->tft.packetfilterlist.createnewtft[0]
        .packetfiltercontents.ipv4remoteaddr[i]
        .mask = 255;
  }
  s11_actv_bearer_request->tft.packetfilterlist.createnewtft[0]
      .packetfiltercontents.singleremoteport = 80;

  s11_actv_bearer_request->s1_u_sgw_fteid.ipv4                = true;
  s11_actv_bearer_request->s1_u_sgw_fteid.ipv4_address.s_addr = 100;

  s11_actv_bearer_request->s1_u_sgw_fteid.teid           = DEFAULT_TEID;
  s11_actv_bearer_request->s1_u_sgw_fteid.interface_type = S1_U_SGW_GTP_U;
  send_msg_to_task(&task_zmq_ctx_main, TASK_MME_APP, message_p);
  return;
}

void send_erab_setup_rsp() {
  MessageDef* message_p =
      itti_alloc_new_message(TASK_S1AP, S1AP_E_RAB_SETUP_RSP);

  S1AP_E_RAB_SETUP_RSP(message_p).mme_ue_s1ap_id = DEFAULT_MME_S1AP_UE_ID;
  S1AP_E_RAB_SETUP_RSP(message_p).enb_ue_s1ap_id = DEFAULT_eNB_S1AP_UE_ID;
  S1AP_E_RAB_SETUP_RSP(message_p).e_rab_setup_list.no_of_items           = 0;
  S1AP_E_RAB_SETUP_RSP(message_p).e_rab_failed_to_setup_list.no_of_items = 0;
  S1AP_E_RAB_SETUP_RSP(message_p).e_rab_setup_list.item[0].e_rab_id      = 6;
  uint8_t transport_address_buff[4] = {192, 168, 60, 141};
  S1AP_E_RAB_SETUP_RSP(message_p)
      .e_rab_setup_list.item[0]
      .transport_layer_address = blk2bstr(transport_address_buff, 4);
  S1AP_E_RAB_SETUP_RSP(message_p).e_rab_setup_list.item[0].gtp_teid =
      DEFAULT_TEID;
  S1AP_E_RAB_SETUP_RSP(message_p).e_rab_setup_list.no_of_items = 1;
  send_msg_to_task(&task_zmq_ctx_main, TASK_MME_APP, message_p);
  return;
}

void send_erab_release_rsp() {
  MessageDef* message_p = itti_alloc_new_message(TASK_S1AP, S1AP_E_RAB_REL_RSP);
  S1AP_E_RAB_REL_RSP(message_p).mme_ue_s1ap_id = DEFAULT_MME_S1AP_UE_ID;
  S1AP_E_RAB_REL_RSP(message_p).enb_ue_s1ap_id = DEFAULT_eNB_S1AP_UE_ID;
  S1AP_E_RAB_REL_RSP(message_p).e_rab_rel_list.no_of_items           = 1;
  S1AP_E_RAB_REL_RSP(message_p).e_rab_failed_to_rel_list.no_of_items = 0;
  S1AP_E_RAB_REL_RSP(message_p).e_rab_rel_list.item[0].e_rab_id      = 6;
  send_msg_to_task(&task_zmq_ctx_main, TASK_MME_APP, message_p);
  return;
}

void send_paging_request() {
  MessageDef* message_p =
      itti_alloc_new_message(TASK_SPGW_APP, S11_PAGING_REQUEST);
  itti_s11_paging_request_t* paging_request_p =
      &message_p->ittiMsg.s11_paging_request;
  paging_request_p->ipv4_addr.s_addr = DEFAULT_UE_IPv4;
  send_msg_to_task(&task_zmq_ctx_main, TASK_MME_APP, message_p);
  return;
}

void send_s1ap_path_switch_req(
    const uint32_t sctp_assoc_id, const uint32_t enb_id,
    const uint32_t enb_ue_s1ap_id, const plmn_t& plmn) {
  MessageDef* message_p =
      itti_alloc_new_message(TASK_S1AP, S1AP_PATH_SWITCH_REQUEST);

  S1AP_PATH_SWITCH_REQUEST(message_p).sctp_assoc_id  = sctp_assoc_id;
  S1AP_PATH_SWITCH_REQUEST(message_p).enb_id         = enb_id;
  S1AP_PATH_SWITCH_REQUEST(message_p).enb_ue_s1ap_id = enb_ue_s1ap_id;
  S1AP_PATH_SWITCH_REQUEST(message_p).mme_ue_s1ap_id = 1;

  S1AP_PATH_SWITCH_REQUEST(message_p).ecgi.plmn                  = plmn;
  S1AP_PATH_SWITCH_REQUEST(message_p).ecgi.cell_identity.enb_id  = enb_id;
  S1AP_PATH_SWITCH_REQUEST(message_p).ecgi.cell_identity.cell_id = 2;
  S1AP_PATH_SWITCH_REQUEST(message_p).ecgi.cell_identity.empty   = 0;

  S1AP_PATH_SWITCH_REQUEST(message_p).tai.plmn = plmn;
  S1AP_PATH_SWITCH_REQUEST(message_p).tai.tac  = 2;

  S1AP_PATH_SWITCH_REQUEST(message_p).encryption_algorithm_capabilities =
      0xc000;
  S1AP_PATH_SWITCH_REQUEST(message_p).integrity_algorithm_capabilities = 0xc000;

  e_rab_to_be_switched_in_downlink_list_t* erab_to_switch_dl_list =
      &S1AP_PATH_SWITCH_REQUEST(message_p).e_rab_to_be_switched_dl_list;

  erab_to_switch_dl_list->no_of_items      = 1;
  erab_to_switch_dl_list->item[0].e_rab_id = 5;  // default bearer id
  erab_to_switch_dl_list->item[0].gtp_teid = 2;
  uint32_t enb_transport_addr              = 0xc0a83c8d;  // 192.168.60.141
  erab_to_switch_dl_list->item[0].transport_layer_address =
      blk2bstr(&enb_transport_addr, 4);

  send_msg_to_task(&task_zmq_ctx_main, TASK_MME_APP, message_p);
  return;
}

}  // namespace lte
}  // namespace magma

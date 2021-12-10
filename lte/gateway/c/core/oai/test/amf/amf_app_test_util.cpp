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
#include "lte/gateway/c/core/oai/test/amf/amf_app_test_util.h"
#include "lte/gateway/c/core/oai/common/conversions.h"
#include <gtest/gtest.h>

namespace magma5g {

/* Create initial ue message without TMSI */
imsi64_t send_initial_ue_message_no_tmsi(
    amf_app_desc_t* amf_app_desc_p, sctp_assoc_id_t sctp_assoc_id,
    uint32_t gnb_id, gnb_ue_ngap_id_t gnb_ue_ngap_id,
    amf_ue_ngap_id_t amf_ue_ngap_id, const plmn_t& plmn, const uint8_t* nas_msg,
    uint8_t nas_msg_length) {
  itti_ngap_initial_ue_message_t initial_ue_message = {};

  initial_ue_message.sctp_assoc_id  = sctp_assoc_id;
  initial_ue_message.gnb_id         = gnb_id;
  initial_ue_message.gnb_ue_ngap_id = gnb_ue_ngap_id;
  initial_ue_message.amf_ue_ngap_id = amf_ue_ngap_id;

  initial_ue_message.nas = blk2bstr(nas_msg, nas_msg_length);

  initial_ue_message.tai.plmn                    = plmn;
  initial_ue_message.tai.tac                     = 1;
  initial_ue_message.ecgi.plmn                   = plmn;
  initial_ue_message.ecgi.cell_identity          = {0, 0, 0};
  initial_ue_message.m5g_rrc_establishment_cause = M5G_MO_SIGNALLING;
  initial_ue_message.ue_context_request = M5G_UEContextRequest_requested;
  initial_ue_message.is_s_tmsi_valid    = false;

  imsi64_t imsi64 = 0;

  imsi64 =
      amf_app_handle_initial_ue_message(amf_app_desc_p, &initial_ue_message);

  return imsi64;
}

/* Create authentication answer from subscriberdb */
int send_proc_authentication_info_answer(
    const std::string& imsi, amf_ue_ngap_id_t ue_id, bool success) {
  itti_amf_subs_auth_info_ans_t aia_itti_msg = {};

  strncpy(aia_itti_msg.imsi, imsi.c_str(), imsi.size());
  aia_itti_msg.imsi_length = imsi.size();
  if (success) {
    aia_itti_msg.result                  = DIAMETER_SUCCESS;
    aia_itti_msg.ue_id                   = ue_id;
    m5g_authentication_info_t* auth_info = &(aia_itti_msg.auth_info);
    auth_info->nb_of_vectors             = 1;

    uint8_t rand_buff[RAND_LENGTH_OCTETS] = {0x12, 0x12, 0xb2, 0x7b, 0xb5, 0xfa,
                                             0x98, 0x85, 0x81, 0x7a, 0x9a, 0x48,
                                             0x56, 0x43, 0x46, 0x3};

    uint8_t xres_data_buff[XRES_LENGTH_MAX] = {
        0x25, 0x70, 0x6f, 0x9a, 0x5b, 0x90, 0xb6, 0xc9,
        0x57, 0x50, 0x6c, 0x88, 0x3d, 0x76, 0xcc, 0x63};

    uint8_t autn_buff[AUTN_LENGTH_OCTETS] = {0xf2, 0x17, 0x69, 0xbe, 0x78, 0xca,
                                             0x80, 0x0,  0xc2, 0xba, 0x59, 0x32,
                                             0x95, 0xf9, 0x72, 0x1e};

    uint8_t kseaf_buff[KSEAF_LENGTH_OCTETS] = {
        0x75, 0xf7, 0xda, 0xaa, 0x9,  0x56, 0x80, 0x3c, 0x1d, 0x66, 0x5f,
        0xf8, 0x74, 0x49, 0xd,  0x22, 0x7a, 0xfb, 0x5e, 0x7d, 0x98, 0xab,
        0xc6, 0x93, 0xe6, 0x2e, 0xc4, 0xa6, 0x89, 0xb2, 0x95, 0x62};

    memcpy(auth_info->m5gauth_vector[0].rand, rand_buff, RAND_LENGTH_OCTETS);
    auth_info->m5gauth_vector[0].xres_star.size = 0x10;
    memcpy(
        auth_info->m5gauth_vector[0].xres_star.data, xres_data_buff,
        XRES_LENGTH_MAX);
    memcpy(auth_info->m5gauth_vector[0].autn, autn_buff, AUTN_LENGTH_OCTETS);
    memcpy(auth_info->m5gauth_vector[0].kseaf, kseaf_buff, KSEAF_LENGTH_OCTETS);
  } else {
    aia_itti_msg.result = DIAMETER_UNABLE_TO_COMPLY;
  }

  int rc = RETURNerror;
  rc     = amf_nas_proc_authentication_info_answer(&aia_itti_msg);

  return (rc);
}

/* Create authentication response from ue */
int send_uplink_nas_message_ue_auth_response(
    amf_app_desc_t* amf_app_desc_p, amf_ue_ngap_id_t ue_id, const plmn_t& plmn,
    const uint8_t* nas_msg, uint8_t nas_msg_length) {
  bstring uplink_nas_auth_response;
  tai_t originating_tai = {};

  uplink_nas_auth_response = blk2bstr(nas_msg, nas_msg_length);

  originating_tai.plmn = plmn;
  originating_tai.tac  = 1;

  int rc = RETURNerror;
  rc     = amf_app_handle_uplink_nas_message(
      amf_app_desc_p, uplink_nas_auth_response, ue_id, originating_tai);

  return (rc);
}

/* Create security mode complete response from ue */
int send_uplink_nas_message_ue_smc_response(
    amf_app_desc_t* amf_app_desc_p, amf_ue_ngap_id_t ue_id, const plmn_t& plmn,
    const uint8_t* nas_msg, uint8_t nas_msg_length) {
  bstring uplink_nas_smc_response;
  tai_t originating_tai = {};

  uplink_nas_smc_response = blk2bstr(nas_msg, nas_msg_length);

  originating_tai.plmn = plmn;
  originating_tai.tac  = 1;

  int rc = RETURNerror;
  rc     = amf_app_handle_uplink_nas_message(
      amf_app_desc_p, uplink_nas_smc_response, ue_id, originating_tai);

  return (rc);
}

void send_initial_context_response(
    amf_app_desc_t* amf_app_desc_p, amf_ue_ngap_id_t ue_id) {
  itti_amf_app_initial_context_setup_rsp_t ics_resp = {};

  // apn profile received from subscriberd during location update
  ue_m5gmm_context_s* ue_context_p =
      amf_ue_context_exists_amf_ue_ngap_id(ue_id);

  ASSERT_NE(nullptr, ue_context_p);

  apn_config_profile_t& profile = ue_context_p->amf_context.apn_config_profile;
  profile.nb_apns               = 1;
  strncpy(profile.apn_configuration[0].service_selection, "internet", 8);

  ics_resp.ue_id = ue_id;

  amf_app_handle_initial_context_setup_rsp(amf_app_desc_p, &ics_resp);
}

/* Create registration mode complete response from ue */
int send_uplink_nas_registration_complete(
    amf_app_desc_t* amf_app_desc_p, amf_ue_ngap_id_t ue_id, const plmn_t& plmn,
    const uint8_t* nas_msg, uint8_t nas_msg_length) {
  bstring ue_registration_complete;
  tai_t originating_tai = {};

  ue_registration_complete = blk2bstr(nas_msg, nas_msg_length);

  originating_tai.plmn = plmn;
  originating_tai.tac  = 1;

  int rc = RETURNerror;
  rc     = amf_app_handle_uplink_nas_message(
      amf_app_desc_p, ue_registration_complete, ue_id, originating_tai);

  return (rc);
}

/* Create pdu session establishment request from ue */
int send_uplink_nas_pdu_session_establishment_request(
    amf_app_desc_t* amf_app_desc_p, amf_ue_ngap_id_t ue_id, const plmn_t& plmn,
    const uint8_t* nas_msg, uint8_t nas_msg_length) {
  bstring pdu_session_est_req;
  tai_t originating_tai = {};

  if ((!amf_app_desc_p) || (!nas_msg) || (nas_msg_length == 0)) {
    return RETURNerror;
  }

  pdu_session_est_req = blk2bstr(nas_msg, nas_msg_length);

  originating_tai.plmn = plmn;
  originating_tai.tac  = 1;

  int rc = RETURNerror;
  rc     = amf_app_handle_uplink_nas_message(
      amf_app_desc_p, pdu_session_est_req, ue_id, originating_tai);

  return rc;
}

void create_ip_address_response_itti(
    itti_amf_ip_allocation_response_t* response) {
  if (!response) return;
  std::string apn = "internet";
  std::copy(apn.begin(), apn.end(), std::begin(response->apn));
  response->default_ambr.br_unit    = BPS;
  response->default_ambr.br_dl      = 200000000;
  response->default_ambr.br_ul      = 100000000;
  response->gnb_gtp_teid            = 0;
  response->gnb_gtp_teid_ip_addr[0] = 0xc0;
  response->gnb_gtp_teid_ip_addr[1] = 0xa8;
  response->gnb_gtp_teid_ip_addr[2] = 0x3c;
  response->gnb_gtp_teid_ip_addr[3] = 0x96;
  std::string imsi                  = "222456000000001";
  std::copy(imsi.begin(), imsi.end(), std::begin(response->imsi));
  response->imsi_length  = 15;
  response->paa.pdn_type = IPv4;
  inet_pton(AF_INET, "192.168.128.254", &(response->paa.ipv4_address));
  response->paa.vlan         = 1;
  response->pdu_session_id   = 1;
  response->pdu_session_type = IPv4;
  response->pti              = 0x01;
  response->result           = 0;
}

int send_ip_address_response_itti() {
  int rc = RETURNerror;

  itti_amf_ip_allocation_response_t response = {};
  create_ip_address_response_itti(&response);

  rc = amf_smf_handle_ip_address_response(&response);

  return rc;
}

void create_pdu_session_response_ipv4_itti(
    itti_n11_create_pdu_session_response_t* response) {
  if (!response) return;
  std::string imsi = "222456000000001";
  std::copy(imsi.begin(), imsi.end(), std::begin(response->imsi));

  response->sm_session_fsm_state = sm_session_fsm_state_t::CREATING;
  response->sm_session_version   = 0;
  response->pdu_session_id       = 1;
  response->pdu_session_type     = IPV4;
  response->selected_ssc_mode    = SSC_MODE_3;
  response->m5gsm_cause          = M5GSM_OPERATION_SUCCESS;

  response->session_ambr.uplink_unit_type   = 0;
  response->session_ambr.uplink_units       = 100000000;
  response->session_ambr.downlink_unit_type = 0;
  response->session_ambr.downlink_units     = 100000000;

  response->qos_list.qos_flow_req_item.qos_flow_identifier = 9;
  response->qos_list.qos_flow_req_item.qos_flow_level_qos_param
      .qos_characteristic.non_dynamic_5QI_desc.fiveQI = 9;
  response->qos_list.qos_flow_req_item.qos_flow_level_qos_param
      .alloc_reten_priority.priority_level = 1;
  response->qos_list.qos_flow_req_item.qos_flow_level_qos_param
      .alloc_reten_priority.pre_emption_cap = SHALL_NOT_TRIGGER_PRE_EMPTION;
  response->qos_list.qos_flow_req_item.qos_flow_level_qos_param
      .alloc_reten_priority.pre_emption_vul = NOT_PREEMPTABLE;
  response->upf_endpoint.teid[0]            = 0x7f;
  response->upf_endpoint.teid[1]            = 0xff;
  response->upf_endpoint.teid[2]            = 0xff;
  response->upf_endpoint.teid[3]            = 0xff;
  inet_pton(AF_INET, "192.168.128.200", response->upf_endpoint.end_ipv4_addr);

  response->always_on_pdu_session_indication     = false;
  response->allowed_ssc_mode                     = SSC_MODE_3;
  response->m5gsm_congetion_re_attempt_indicator = true;
  response->pdu_address.redirect_address_type    = IPV4_1;
  inet_pton(
      AF_INET, "192.168.128.200",
      response->pdu_address.redirect_server_address);
}

int send_pdu_session_response_itti() {
  int rc                                          = RETURNerror;
  itti_n11_create_pdu_session_response_t response = {};
  create_pdu_session_response_ipv4_itti(&response);

  rc = amf_app_handle_pdu_session_response(&response);

  return rc;
}

void create_pdu_resource_setup_response_itti(
    itti_ngap_pdusessionresource_setup_rsp_t* response,
    amf_ue_ngap_id_t ue_id) {
  if (!response) return;
  response->amf_ue_ngap_id                                       = ue_id;
  response->gnb_ue_ngap_id                                       = 1;
  response->pduSessionResource_setup_list.item[0].Pdu_Session_ID = 1;
  response->pduSessionResource_setup_list.no_of_items            = 1;
  response_gtp_tunnel_t* tunnel =
      &response->pduSessionResource_setup_list.item[0]
           .PDU_Session_Resource_Setup_Response_Transfer.tunnel;
  tunnel->gTP_TEID[0] = 0x0;
  tunnel->gTP_TEID[1] = 0x0;
  tunnel->gTP_TEID[2] = 0x0;
  tunnel->gTP_TEID[3] = 0x1;

  tunnel->transportLayerAddress[0] = 0xc0;
  tunnel->transportLayerAddress[1] = 0xa8;
  tunnel->transportLayerAddress[2] = 0x3c;
  tunnel->transportLayerAddress[3] = 0x96;

  AssociatedQosFlowList_t* qosFlow =
      &response->pduSessionResource_setup_list.item[0]
           .PDU_Session_Resource_Setup_Response_Transfer.associatedQosFlowList;
  qosFlow->items                = 1;
  qosFlow->QosFlowIdentifier[0] = 9;
}
int send_pdu_resource_setup_response(amf_ue_ngap_id_t ue_id) {
  int rc                                            = RETURNok;
  itti_ngap_pdusessionresource_setup_rsp_t response = {};
  create_pdu_resource_setup_response_itti(&response, ue_id);

  amf_app_handle_resource_setup_response(response);

  return rc;
}

void create_pdu_notification_response_itti(
    itti_n11_received_notification_t* response) {
  if (!response) return;
  std::string imsi = "222456000000001";
  std::copy(imsi.begin(), imsi.end(), std::begin(response->imsi));
  response->sm_session_fsm_state = CREATING_0;
  response->sm_session_version   = 0;
  response->pdu_session_id       = 1;
  response->request_type         = EXISTING_PDU_SESSION;
  response->m5g_sm_capability.multi_homed_ipv6_pdu_session = false;
  response->m5g_sm_capability.reflective_qos               = false;
  response->m5gsm_cause      = M5GSM_OPERATION_SUCCESS;
  response->pdu_session_type = IPV4;
  response->notify_ue_evnt   = PDU_SESSION_STATE_NOTIFY;
}

int send_pdu_notification_response() {
  int rc                                    = RETURNerror;
  itti_n11_received_notification_t response = {};
  create_pdu_notification_response_itti(&response);

  rc = amf_app_handle_notification_received(&response);

  return rc;
}

/* Create pdu session release message from ue */
int send_uplink_nas_pdu_session_release_message(
    amf_app_desc_t* amf_app_desc_p, amf_ue_ngap_id_t ue_id, const plmn_t& plmn,
    const uint8_t* nas_msg, uint8_t nas_msg_length) {
  bstring pdu_session_req;
  tai_t originating_tai = {};

  if ((!amf_app_desc_p) || (!nas_msg) || (nas_msg_length == 0)) {
    return RETURNerror;
  }

  pdu_session_req = blk2bstr(nas_msg, nas_msg_length);

  originating_tai.plmn = plmn;
  originating_tai.tac  = 1;

  int rc = RETURNerror;
  rc     = amf_app_handle_uplink_nas_message(
      amf_app_desc_p, pdu_session_req, ue_id, originating_tai);

  return (rc);
}

/* Create security mode complete response from ue */
int send_uplink_nas_ue_deregistration_request(
    amf_app_desc_t* amf_app_desc_p, amf_ue_ngap_id_t ue_id, const plmn_t& plmn,
    uint8_t* nas_msg, uint8_t nas_msg_length) {
  bstring uplink_nas_ue_dereg_req;
  tai_t originating_tai = {};
  int rc                = RETURNerror;

  tmsi_t ue_tmsi = htonl(amf_lookup_guti_by_ueid(ue_id));
  if (ue_tmsi == 0) {
    return (rc);
  }

  uplink_nas_ue_dereg_req = blk2bstr(nas_msg, nas_msg_length);

  /* Last 4 bytes are 0. Replace with TMSI */
  memcpy(
      &(uplink_nas_ue_dereg_req->data[nas_msg_length - sizeof(tmsi_t)]),
      &ue_tmsi, sizeof(tmsi_t));

  originating_tai.plmn = plmn;
  originating_tai.tac  = 1;

  rc = amf_app_handle_uplink_nas_message(
      amf_app_desc_p, uplink_nas_ue_dereg_req, ue_id, originating_tai);

  return (rc);
}

/* Get the ue id from IMSI */
bool get_ue_id_from_imsi(
    amf_app_desc_t* amf_app_desc_p, imsi64_t imsi64, amf_ue_ngap_id_t* ue_id) {
  amf_ue_context_t* amf_ue_context_p = &amf_app_desc_p->amf_ue_contexts;
  magma::map_rc_t rc_map             = magma::MAP_OK;
  rc_map = amf_ue_context_p->imsi_amf_ue_id_htbl.get(imsi64, ue_id);
  if (rc_map != magma::MAP_OK) {
    return (false);
  }
  return true;
}

}  // namespace magma5g

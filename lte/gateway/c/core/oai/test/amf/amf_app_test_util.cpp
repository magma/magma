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
#include "lte/gateway/c/core/common/dynamic_memory_check.h"
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

  initial_ue_message.sctp_assoc_id = sctp_assoc_id;
  initial_ue_message.gnb_id = gnb_id;
  initial_ue_message.gnb_ue_ngap_id = gnb_ue_ngap_id;
  initial_ue_message.amf_ue_ngap_id = amf_ue_ngap_id;

  initial_ue_message.nas = blk2bstr(nas_msg, nas_msg_length);

  initial_ue_message.tai.plmn = plmn;
  initial_ue_message.tai.tac = 1;
  initial_ue_message.ecgi.plmn = plmn;
  initial_ue_message.ecgi.cell_identity = {0, 0, 0};
  initial_ue_message.m5g_rrc_establishment_cause = M5G_MO_SIGNALING;
  initial_ue_message.ue_context_request = M5G_UEContextRequest_requested;
  initial_ue_message.is_s_tmsi_valid = false;

  imsi64_t imsi64 = 0;

  imsi64 =
      amf_app_handle_initial_ue_message(amf_app_desc_p, &initial_ue_message);

  return imsi64;
}

// Create initial ue message without TMSI and replace tmsi in hexbuf
imsi64_t send_initial_ue_message_no_tmsi_replace_mtmsi(
    amf_app_desc_t* amf_app_desc_p, sctp_assoc_id_t sctp_assoc_id,
    uint32_t gnb_id, gnb_ue_ngap_id_t gnb_ue_ngap_id,
    amf_ue_ngap_id_t amf_ue_ngap_id, const plmn_t& plmn, const uint8_t* nas_msg,
    uint8_t nas_msg_length, amf_ue_ngap_id_t ue_id, uint8_t tmsi_offset) {
  itti_ngap_initial_ue_message_t initial_ue_message = {};

  initial_ue_message.sctp_assoc_id = sctp_assoc_id;
  initial_ue_message.gnb_id = gnb_id;
  initial_ue_message.gnb_ue_ngap_id = gnb_ue_ngap_id;
  initial_ue_message.amf_ue_ngap_id = amf_ue_ngap_id;

  initial_ue_message.nas = blk2bstr(nas_msg, nas_msg_length);

  initial_ue_message.tai.plmn = plmn;
  initial_ue_message.tai.tac = 1;
  initial_ue_message.ecgi.plmn = plmn;
  initial_ue_message.ecgi.cell_identity = {0, 0, 0};
  initial_ue_message.m5g_rrc_establishment_cause = M5G_MO_SIGNALING;
  initial_ue_message.ue_context_request = M5G_UEContextRequest_requested;
  initial_ue_message.is_s_tmsi_valid = false;
  tmsi_t ue_tmsi = amf_lookup_guti_by_ueid(ue_id);

  ue_tmsi = htonl(ue_tmsi);
  memcpy(&(initial_ue_message.nas
               ->data[nas_msg_length - sizeof(tmsi_t) - tmsi_offset]),
         &(ue_tmsi), sizeof(tmsi_t));

  imsi64_t imsi64 = 0;

  imsi64 =
      amf_app_handle_initial_ue_message(amf_app_desc_p, &initial_ue_message);

  return imsi64;
}

/* Create initial ue message without TMSI no context*/
imsi64_t send_initial_ue_message_no_tmsi_no_ctx_req(
    amf_app_desc_t* amf_app_desc_p, sctp_assoc_id_t sctp_assoc_id,
    uint32_t gnb_id, gnb_ue_ngap_id_t gnb_ue_ngap_id,
    amf_ue_ngap_id_t amf_ue_ngap_id, const plmn_t& plmn, const uint8_t* nas_msg,
    uint8_t nas_msg_length) {
  itti_ngap_initial_ue_message_t initial_ue_message = {};

  initial_ue_message.sctp_assoc_id = sctp_assoc_id;
  initial_ue_message.gnb_id = gnb_id;
  initial_ue_message.gnb_ue_ngap_id = gnb_ue_ngap_id;
  initial_ue_message.amf_ue_ngap_id = amf_ue_ngap_id;

  initial_ue_message.nas = blk2bstr(nas_msg, nas_msg_length);

  initial_ue_message.tai.plmn = plmn;
  initial_ue_message.tai.tac = 1;
  initial_ue_message.ecgi.plmn = plmn;
  initial_ue_message.ecgi.cell_identity = {0, 0, 0};
  initial_ue_message.m5g_rrc_establishment_cause = M5G_MO_SIGNALING;
  initial_ue_message.is_s_tmsi_valid = false;

  imsi64_t imsi64 = 0;

  imsi64 =
      amf_app_handle_initial_ue_message(amf_app_desc_p, &initial_ue_message);

  return imsi64;
}

/* For guti based registration */
uint64_t send_initial_ue_message_with_tmsi(
    amf_app_desc_t* amf_app_desc_p, sctp_assoc_id_t sctp_assoc_id,
    uint32_t gnb_id, gnb_ue_ngap_id_t gnb_ue_ngap_id,
    amf_ue_ngap_id_t amf_ue_ngap_id, const plmn_t& plmn, uint32_t m_tmsi,
    const uint8_t* nas_msg, uint8_t nas_msg_length) {
  tai_t originating_tai = {};
  int rc = RETURNerror;
  tai_t tai = {.plmn = plmn, .tac = 1};

  itti_ngap_initial_ue_message_t initial_ue_message = {};

  initial_ue_message.nas = blk2bstr(nas_msg, nas_msg_length);

  initial_ue_message.sctp_assoc_id = sctp_assoc_id;
  initial_ue_message.gnb_ue_ngap_id = gnb_ue_ngap_id;
  initial_ue_message.gnb_id = gnb_id;
  initial_ue_message.m5g_rrc_establishment_cause = M5G_MO_SIGNALING;
  initial_ue_message.is_s_tmsi_valid = true;
  initial_ue_message.opt_s_tmsi.amf_set_id = 1;
  initial_ue_message.opt_s_tmsi.amf_pointer = 0;
  initial_ue_message.opt_s_tmsi.m_tmsi = m_tmsi;
  initial_ue_message.tai = tai;
  initial_ue_message.ue_context_request = M5G_UEContextRequest_requested;

  amf_app_handle_initial_ue_message(amf_app_desc_p, &initial_ue_message);

  guti_m5_t guti_search;
  memset(&guti_search, 0, sizeof(guti_m5_t));

  guti_search.guamfi.amf_regionid = 1;
  guti_search.guamfi.amf_set_id = 1;
  guti_search.guamfi.amf_pointer = 0;
  guti_search.guamfi.plmn = plmn;
  guti_search.m_tmsi = m_tmsi;

  uint64_t amf_ue_ngap_id64 = 0;
  magma::map_rc_t m_rc = magma::MAP_KEY_NOT_EXISTS;
  rc = amf_app_desc_p->amf_ue_contexts.guti_ue_context_htbl.get(
      guti_search, &amf_ue_ngap_id64);

  return amf_ue_ngap_id64;
}

/* Create the identiy response message */
status_code_e send_uplink_nas_identity_response_message(
    amf_app_desc_t* amf_app_desc_p, amf_ue_ngap_id_t ue_id, const plmn_t& plmn,
    const uint8_t* nas_msg, uint8_t nas_msg_length) {
  bstring pdu_session_req;
  tai_t originating_tai = {};

  if ((!amf_app_desc_p) || (!nas_msg) || (nas_msg_length == 0)) {
    return RETURNerror;
  }

  pdu_session_req = blk2bstr(nas_msg, nas_msg_length);

  originating_tai.plmn = plmn;
  originating_tai.tac = 1;

  status_code_e rc = RETURNerror;
  rc = amf_app_handle_uplink_nas_message(amf_app_desc_p, pdu_session_req, ue_id,
                                         originating_tai);

  return rc;
}

/* Create authentication answer from subscriberdb */
status_code_e send_proc_authentication_info_answer(const std::string& imsi,
                                                   amf_ue_ngap_id_t ue_id,
                                                   bool success) {
  itti_amf_subs_auth_info_ans_t aia_itti_msg = {};

  strncpy(aia_itti_msg.imsi, imsi.c_str(), imsi.size());
  aia_itti_msg.imsi_length = imsi.size();
  if (success) {
    aia_itti_msg.result = DIAMETER_SUCCESS;
    aia_itti_msg.ue_id = ue_id;
    m5g_authentication_info_t* auth_info = &(aia_itti_msg.auth_info);
    auth_info->nb_of_vectors = 1;

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
    memcpy(auth_info->m5gauth_vector[0].xres_star.data, xres_data_buff,
           XRES_LENGTH_MAX);
    memcpy(auth_info->m5gauth_vector[0].autn, autn_buff, AUTN_LENGTH_OCTETS);
    memcpy(auth_info->m5gauth_vector[0].kseaf, kseaf_buff, KSEAF_LENGTH_OCTETS);
  } else {
    aia_itti_msg.result = DIAMETER_UNABLE_TO_COMPLY;
  }

  status_code_e rc = RETURNerror;
  rc = amf_nas_proc_authentication_info_answer(&aia_itti_msg);

  return (rc);
}

/* Create authentication response from ue */
status_code_e send_uplink_nas_message_ue_auth_response(
    amf_app_desc_t* amf_app_desc_p, amf_ue_ngap_id_t ue_id, const plmn_t& plmn,
    const uint8_t* nas_msg, uint8_t nas_msg_length) {
  bstring uplink_nas_auth_response;
  tai_t originating_tai = {};

  uplink_nas_auth_response = blk2bstr(nas_msg, nas_msg_length);

  originating_tai.plmn = plmn;
  originating_tai.tac = 1;

  status_code_e rc = RETURNerror;
  rc = amf_app_handle_uplink_nas_message(
      amf_app_desc_p, uplink_nas_auth_response, ue_id, originating_tai);

  return (rc);
}

/* Create security mode complete response from ue */
status_code_e send_uplink_nas_message_ue_smc_response(
    amf_app_desc_t* amf_app_desc_p, amf_ue_ngap_id_t ue_id, const plmn_t& plmn,
    const uint8_t* nas_msg, uint8_t nas_msg_length) {
  bstring uplink_nas_smc_response;
  tai_t originating_tai = {};

  uplink_nas_smc_response = blk2bstr(nas_msg, nas_msg_length);

  originating_tai.plmn = plmn;
  originating_tai.tac = 1;

  status_code_e rc = RETURNerror;
  rc = amf_app_handle_uplink_nas_message(
      amf_app_desc_p, uplink_nas_smc_response, ue_id, originating_tai);

  return (rc);
}

void send_initial_context_response(amf_app_desc_t* amf_app_desc_p,
                                   amf_ue_ngap_id_t ue_id) {
  itti_amf_app_initial_context_setup_rsp_t ics_resp = {};

  // apn profile received from subscriberd during location update
  ue_m5gmm_context_s* ue_context_p =
      amf_ue_context_exists_amf_ue_ngap_id(ue_id);

  ASSERT_NE(nullptr, ue_context_p);

  apn_config_profile_t& profile = ue_context_p->amf_context.apn_config_profile;
  profile.nb_apns = 1;
  strncpy(profile.apn_configuration[0].service_selection, "internet", 8);
  profile.apn_configuration[0].service_selection_length = 8;

  ics_resp.ue_id = ue_id;

  amf_app_handle_initial_context_setup_rsp(amf_app_desc_p, &ics_resp);
}

/* Create registration mode complete response from ue */
status_code_e send_uplink_nas_registration_complete(
    amf_app_desc_t* amf_app_desc_p, amf_ue_ngap_id_t ue_id, const plmn_t& plmn,
    const uint8_t* nas_msg, uint8_t nas_msg_length) {
  bstring ue_registration_complete;
  tai_t originating_tai = {};

  ue_registration_complete = blk2bstr(nas_msg, nas_msg_length);

  originating_tai.plmn = plmn;
  originating_tai.tac = 1;

  status_code_e rc = RETURNerror;
  rc = amf_app_handle_uplink_nas_message(
      amf_app_desc_p, ue_registration_complete, ue_id, originating_tai);

  return rc;
}

/* Create pdu session establishment request from ue */
status_code_e send_uplink_nas_pdu_session_establishment_request(
    amf_app_desc_t* amf_app_desc_p, amf_ue_ngap_id_t ue_id, const plmn_t& plmn,
    const uint8_t* nas_msg, uint8_t nas_msg_length) {
  bstring pdu_session_est_req;
  tai_t originating_tai = {};

  if ((!amf_app_desc_p) || (!nas_msg) || (nas_msg_length == 0)) {
    return RETURNerror;
  }

  pdu_session_est_req = blk2bstr(nas_msg, nas_msg_length);

  originating_tai.plmn = plmn;
  originating_tai.tac = 1;

  status_code_e rc = RETURNerror;
  rc = amf_app_handle_uplink_nas_message(amf_app_desc_p, pdu_session_est_req,
                                         ue_id, originating_tai);

  return rc;
}
/* Create pdu session modification complete from ue */
int send_uplink_nas_pdu_session_modification_complete(
    amf_app_desc_t* amf_app_desc_p, amf_ue_ngap_id_t ue_id, const plmn_t& plmn,
    const uint8_t* nas_msg, uint8_t nas_msg_length) {
  bstring pdu_session_est_req;
  tai_t originating_tai = {};

  if ((!amf_app_desc_p) || (!nas_msg) || (nas_msg_length == 0)) {
    return RETURNerror;
  }

  pdu_session_est_req = blk2bstr(nas_msg, nas_msg_length);

  originating_tai.plmn = plmn;
  originating_tai.tac = 1;

  ue_m5gmm_context_s* ue_context_p =
      amf_ue_context_exists_amf_ue_ngap_id(ue_id);

  if (!ue_context_p) {
    return RETURNerror;
  }
  apn_config_profile_t& profile = ue_context_p->amf_context.apn_config_profile;
  profile.nb_apns = 1;
  strncpy(profile.apn_configuration[0].service_selection, "internet",
          SERVICE_SELECTION_MAX_LENGTH - 1);

  int rc = RETURNerror;
  rc = amf_app_handle_uplink_nas_message(amf_app_desc_p, pdu_session_est_req,
                                         ue_id, originating_tai);

  return rc;
}

void create_ip_address_response_itti(
    pdn_type_value_t type, itti_amf_ip_allocation_response_t* response) {
  if (!response) return;
  std::string apn = "internet";
  std::copy(apn.begin(), apn.end(), std::begin(response->apn));
  response->default_ambr.br_unit = BPS;
  response->default_ambr.br_dl = 200000000;
  response->default_ambr.br_ul = 100000000;
  response->gnb_gtp_teid = 0;
  response->gnb_gtp_teid_ip_addr[0] = 0xc0;
  response->gnb_gtp_teid_ip_addr[1] = 0xa8;
  response->gnb_gtp_teid_ip_addr[2] = 0x3c;
  response->gnb_gtp_teid_ip_addr[3] = 0x96;
  std::string imsi = "222456000000001";
  std::copy(imsi.begin(), imsi.end(), std::begin(response->imsi));
  response->imsi_length = 15;

  if (type == IPv4) {
    response->pdu_session_type = IPV4;
  } else if (type == IPv6) {
    response->pdu_session_type = IPV6;
  } else {
    response->pdu_session_type = IPV4IPV6;
  }

  response->paa.pdn_type = type;
  if ((type == IPv4) || (type == IPv4_AND_v6)) {
    inet_pton(AF_INET, "192.168.128.20", &(response->paa.ipv4_address));
  }

  if ((type == IPv6) || (type == IPv4_AND_v6)) {
    inet_pton(AF_INET6, "2001:db8:3a:dd2:0253:a1ff:fe2c:831f",
              &(response->paa.ipv6_address));
    response->paa.ipv6_prefix_length = IPV6_PREFIX_LEN;
  }

  response->paa.vlan = 1;
  response->pdu_session_id = 1;
  response->pti = 0x01;
  response->result = SGI_STATUS_OK;
}

status_code_e send_ip_address_response_itti(pdn_type_value_t type) {
  status_code_e rc = RETURNerror;

  itti_amf_ip_allocation_response_t response = {};
  create_ip_address_response_itti(type, &response);

  rc = amf_smf_handle_ip_address_response(&response);

  return rc;
}

void create_pdu_session_response_itti(
    pdn_type_value_t type, itti_n11_create_pdu_session_response_t* response) {
  if (!response) return;
  std::string imsi = "222456000000001";
  std::copy(imsi.begin(), imsi.end(), std::begin(response->imsi));

  response->sm_session_fsm_state = sm_session_fsm_state_t::CREATING;
  response->sm_session_version = 0;
  response->pdu_session_id = 1;
  response->selected_ssc_mode = SSC_MODE_3;
  response->m5gsm_cause = M5GSM_OPERATION_SUCCESS;

  response->session_ambr.uplink_unit_type = 0;
  response->session_ambr.uplink_units = 100000000;
  response->session_ambr.downlink_unit_type = 0;
  response->session_ambr.downlink_units = 100000000;

  response->qos_flow_list.maxNumOfQosFlows = 1;
  response->qos_flow_list.item[0].qos_flow_req_item.qos_flow_identifier = 9;
  response->qos_flow_list.item[0]
      .qos_flow_req_item.qos_flow_level_qos_param.qos_characteristic
      .non_dynamic_5QI_desc.fiveQI = 9;

  // Setting default flow descriptor
  response->qos_flow_list.item[0]
      .qos_flow_req_item.qos_flow_descriptor.qos_flow_identifier = 9;
  response->qos_flow_list.item[0].qos_flow_req_item.qos_flow_descriptor.fiveQi =
      9;

  response->qos_flow_list.item[0]
      .qos_flow_req_item.qos_flow_level_qos_param.alloc_reten_priority
      .priority_level = 1;
  response->qos_flow_list.item[0]
      .qos_flow_req_item.qos_flow_level_qos_param.alloc_reten_priority
      .pre_emption_cap = SHALL_NOT_TRIGGER_PRE_EMPTION;
  response->qos_flow_list.item[0]
      .qos_flow_req_item.qos_flow_level_qos_param.alloc_reten_priority
      .pre_emption_vul = NOT_PREEMPTABLE;
  response->qos_flow_list.item[0]
      .qos_flow_req_item.qos_flow_descriptor.qos_flow_identifier = 9;
  response->qos_flow_list.item[0].qos_flow_req_item.qos_flow_descriptor.fiveQi =
      9;

  response->upf_endpoint.teid[0] = 0x7f;
  response->upf_endpoint.teid[1] = 0xff;
  response->upf_endpoint.teid[2] = 0xff;
  response->upf_endpoint.teid[3] = 0xff;
  inet_pton(AF_INET, "192.168.128.200", response->upf_endpoint.end_ipv4_addr);

  response->always_on_pdu_session_indication = false;
  response->allowed_ssc_mode = SSC_MODE_3;
  response->m5gsm_congetion_re_attempt_indicator = true;

  response->pdu_address.pdn_type = type;
  // response->pdu_session_type     = type;
  if ((type == IPv4) || (type == IPv4_AND_v6)) {
    std::string ue_ipv4("192.168.128.20");
    inet_pton(AF_INET, ue_ipv4.c_str(), &(response->pdu_address.ipv4_address));
  }

  if ((type == IPv6) || (type == IPv4_AND_v6)) {
    std::string ue_ipv6("2001:db8:3a:dd2:0253:a1ff:fe2c:831f");
    inet_pton(AF_INET6, ue_ipv6.c_str(), &(response->pdu_address.ipv6_address));
    response->pdu_address.ipv6_prefix_length = IPV6_PREFIX_LEN;
  }
}

status_code_e send_pdu_session_response_itti(pdn_type_value_t type) {
  status_code_e rc = RETURNerror;
  itti_n11_create_pdu_session_response_t response = {};
  create_pdu_session_response_itti(type, &response);

  rc = amf_app_handle_pdu_session_response(&response);

  return rc;
}

int send_pdu_session_modification_itti() {
  int rc = RETURNerror;
  itti_n11_create_pdu_session_response_t response = {};
  create_pdu_session_modify_request_itti(&response);

  rc = amf_app_handle_pdu_session_response(&response);

  return rc;
}

int send_pdu_session_modification_deletion_itti() {
  int rc = RETURNerror;
  itti_n11_create_pdu_session_response_t response = {};
  create_pdu_session_modify_deletion_request_itti(&response);

  rc = amf_app_handle_pdu_session_response(&response);

  return rc;
}

void create_pdu_resource_setup_response_itti(
    itti_ngap_pdusessionresource_setup_rsp_t* response,
    amf_ue_ngap_id_t ue_id) {
  if (!response) return;
  response->amf_ue_ngap_id = ue_id;
  response->gnb_ue_ngap_id = 1;
  response->pduSessionResource_setup_list.item[0].Pdu_Session_ID = 1;
  response->pduSessionResource_setup_list.no_of_items = 1;
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
  qosFlow->items = 1;
  qosFlow->QosFlowIdentifier[0] = 9;
}

void create_pdu_resource_modify_response_itti(
    itti_ngap_pdu_session_resource_modify_response_t* response,
    amf_ue_ngap_id_t ue_id) {
  if (!response) return;
  response->amf_ue_ngap_id = ue_id;
  response->gnb_ue_ngap_id = 1;
  response->pduSessResourceModRespList.item[0].Pdu_Session_ID = 1;
  response->pduSessResourceModRespList.no_of_items = 1;

  pdusession_modify_response_item_t* pduSessModifyResp =
      &response->pduSessResourceModRespList.item[0];
  pduSessModifyResp->PDU_Session_Resource_Mpdify_Response_Transfer
      .qos_flow_add_or_modify_response_list.maxNumOfQosFlows = 1;
  pduSessModifyResp->PDU_Session_Resource_Mpdify_Response_Transfer
      .qos_flow_add_or_modify_response_list.item[0]
      .qos_flow_identifier = 3;
}
status_code_e send_pdu_resource_setup_response(amf_ue_ngap_id_t ue_id) {
  status_code_e rc = RETURNok;
  itti_ngap_pdusessionresource_setup_rsp_t response = {};
  create_pdu_resource_setup_response_itti(&response, ue_id);

  amf_app_handle_resource_setup_response(response);

  return rc;
}

void create_pdu_session_modify_deletion_request_itti(
    itti_n11_create_pdu_session_response_t* response) {
  if (!response) return;
  std::string imsi = "222456000000001";
  std::copy(imsi.begin(), imsi.end(), std::begin(response->imsi));

  response->sm_session_fsm_state = sm_session_fsm_state_t::ACTIVE;
  response->sm_session_version = 0;
  response->pdu_session_id = 1;
  response->pdu_session_type = IPV4;
  response->selected_ssc_mode = SSC_MODE_3;
  response->m5gsm_cause = M5GSM_OPERATION_SUCCESS;

  // dedicated qos flow
  response->qos_flow_list.maxNumOfQosFlows = 1;
  response->qos_flow_list.item[0].qos_flow_req_item.qos_flow_identifier = 3;
  response->qos_flow_list.item[0]
      .qos_flow_req_item.qos_flow_level_qos_param.qos_characteristic
      .non_dynamic_5QI_desc.fiveQI = 3;
  response->qos_flow_list.item[0]
      .qos_flow_req_item.qos_flow_level_qos_param.alloc_reten_priority
      .priority_level = 1;
  response->qos_flow_list.item[0]
      .qos_flow_req_item.qos_flow_level_qos_param.alloc_reten_priority
      .pre_emption_cap = SHALL_NOT_TRIGGER_PRE_EMPTION;
  response->qos_flow_list.item[0]
      .qos_flow_req_item.qos_flow_level_qos_param.alloc_reten_priority
      .pre_emption_vul = NOT_PREEMPTABLE;

  // flow action
  response->qos_flow_list.item[0].qos_flow_req_item.qos_flow_action =
      policy_action_del;
  response->qos_flow_list.item[0].qos_flow_req_item.qos_flow_version = 2;
  // flow description
  response->qos_flow_list.item[0]
      .qos_flow_req_item.qos_flow_descriptor.qos_flow_identifier =
      response->qos_flow_list.item[0].qos_flow_req_item.qos_flow_identifier;
  response->qos_flow_list.item[0].qos_flow_req_item.qos_flow_descriptor.fiveQi =
      response->qos_flow_list.item[0].qos_flow_req_item.qos_flow_identifier;
  response->qos_flow_list.item[0].qos_flow_req_item.qos_flow_descriptor.mbr_ul =
      100000;
  response->qos_flow_list.item[0].qos_flow_req_item.qos_flow_descriptor.mbr_dl =
      100000;
  response->qos_flow_list.item[0].qos_flow_req_item.qos_flow_descriptor.gbr_ul =
      100000;
  response->qos_flow_list.item[0].qos_flow_req_item.qos_flow_descriptor.gbr_dl =
      100000;

  // traffic flow template
  response->qos_flow_list.item[0].qos_flow_req_item.ul_tft.tftoperationcode =
      TRAFFIC_FLOW_TEMPLATE_OPCODE_DELETE_EXISTING_TFT;
  response->qos_flow_list.item[0]
      .qos_flow_req_item.ul_tft.numberofpacketfilters = 1;

  // rule id
  strncpy(reinterpret_cast<char*>(
              response->qos_flow_list.item[0].qos_flow_req_item.rule_id),
          "rule2", 6);
}

int send_pdu_resource_modify_response(amf_ue_ngap_id_t ue_id) {
  int rc = RETURNok;
  itti_ngap_pdu_session_resource_modify_response_t response = {};
  create_pdu_resource_modify_response_itti(&response, ue_id);

  amf_app_handle_resource_modify_response(response);

  return rc;
}

void create_pdu_notification_response_itti(
    itti_n11_received_notification_t* response) {
  if (!response) return;
  std::string imsi = "222456000000001";
  std::copy(imsi.begin(), imsi.end(), std::begin(response->imsi));
  response->sm_session_fsm_state = CREATING_0;
  response->sm_session_version = 0;
  response->pdu_session_id = 1;
  response->request_type = EXISTING_PDU_SESSION;
  response->m5g_sm_capability.multi_homed_ipv6_pdu_session = false;
  response->m5g_sm_capability.reflective_qos = false;
  response->m5gsm_cause = M5GSM_OPERATION_SUCCESS;
  response->pdu_session_type = IPV4;
  response->notify_ue_evnt = PDU_SESSION_STATE_NOTIFY;
}

status_code_e send_pdu_notification_response() {
  status_code_e rc = RETURNerror;
  itti_n11_received_notification_t response = {};
  create_pdu_notification_response_itti(&response);

  rc = amf_app_handle_notification_received(&response);

  return rc;
}

void create_pdu_session_modify_request_itti(
    itti_n11_create_pdu_session_response_t* response) {
  if (!response) return;
  std::string imsi = "222456000000001";
  std::copy(imsi.begin(), imsi.end(), std::begin(response->imsi));

  response->sm_session_fsm_state = sm_session_fsm_state_t::ACTIVE;
  response->sm_session_version = 0;
  response->pdu_session_id = 1;
  response->pdu_session_type = IPV4;
  response->selected_ssc_mode = SSC_MODE_3;
  response->m5gsm_cause = M5GSM_OPERATION_SUCCESS;

  // dedicated qos flow
  response->qos_flow_list.maxNumOfQosFlows = 1;
  response->qos_flow_list.item[0].qos_flow_req_item.qos_flow_identifier = 3;
  response->qos_flow_list.item[0]
      .qos_flow_req_item.qos_flow_level_qos_param.qos_characteristic
      .non_dynamic_5QI_desc.fiveQI = 3;
  response->qos_flow_list.item[0]
      .qos_flow_req_item.qos_flow_level_qos_param.alloc_reten_priority
      .priority_level = 1;
  response->qos_flow_list.item[0]
      .qos_flow_req_item.qos_flow_level_qos_param.alloc_reten_priority
      .pre_emption_cap = SHALL_NOT_TRIGGER_PRE_EMPTION;
  response->qos_flow_list.item[0]
      .qos_flow_req_item.qos_flow_level_qos_param.alloc_reten_priority
      .pre_emption_vul = NOT_PREEMPTABLE;

  // flow action
  response->qos_flow_list.item[0].qos_flow_req_item.qos_flow_action =
      policy_action_add;
  response->qos_flow_list.item[0].qos_flow_req_item.qos_flow_version = 2;

  // flow description
  response->qos_flow_list.item[0]
      .qos_flow_req_item.qos_flow_descriptor.qos_flow_identifier = 3;

  response->qos_flow_list.item[0].qos_flow_req_item.qos_flow_descriptor.fiveQi =
      3;

  response->qos_flow_list.item[0]
      .qos_flow_req_item.qos_flow_descriptor.qos_flow_identifier =
      response->qos_flow_list.item[0].qos_flow_req_item.qos_flow_identifier;
  response->qos_flow_list.item[0].qos_flow_req_item.qos_flow_descriptor.fiveQi =
      response->qos_flow_list.item[0].qos_flow_req_item.qos_flow_identifier;
  response->qos_flow_list.item[0].qos_flow_req_item.qos_flow_descriptor.mbr_ul =
      100000;
  response->qos_flow_list.item[0].qos_flow_req_item.qos_flow_descriptor.mbr_dl =
      100000;
  response->qos_flow_list.item[0].qos_flow_req_item.qos_flow_descriptor.gbr_ul =
      100000;
  response->qos_flow_list.item[0].qos_flow_req_item.qos_flow_descriptor.gbr_dl =
      100000;

  // traffic flow template
  response->qos_flow_list.item[0].qos_flow_req_item.ul_tft.tftoperationcode =
      TRAFFIC_FLOW_TEMPLATE_OPCODE_CREATE_NEW_TFT;
  response->qos_flow_list.item[0]
      .qos_flow_req_item.ul_tft.numberofpacketfilters = 1;
  // create new tft
  create_new_tft_t* new_tft =
      &response->qos_flow_list.item[0]
           .qos_flow_req_item.ul_tft.packetfilterlist.createnewtft[0];
  response->qos_flow_list.item[0]
      .qos_flow_req_item.ul_tft.parameterslist.num_parameters = 5;
  new_tft->direction = TRAFFIC_FLOW_TEMPLATE_UPLINK_ONLY;
  new_tft->identifier = 1;
  new_tft->length = 14;
  new_tft->eval_precedence = 254;
  new_tft->packetfiltercontents.flags =
      TRAFFIC_FLOW_TEMPLATE_IPV4_REMOTE_ADDR_FLAG |
      TRAFFIC_FLOW_TEMPLATE_SINGLE_REMOTE_PORT;
  new_tft->packetfiltercontents.ipv4remoteaddr[0].addr = 10;
  new_tft->packetfiltercontents.ipv4remoteaddr[1].addr = 10;
  new_tft->packetfiltercontents.ipv4remoteaddr[2].addr = 2;
  new_tft->packetfiltercontents.ipv4remoteaddr[3].addr = 2;
  new_tft->packetfiltercontents.ipv4remoteaddr[0].mask = 0xff;
  new_tft->packetfiltercontents.ipv4remoteaddr[1].mask = 0xff;
  new_tft->packetfiltercontents.ipv4remoteaddr[2].mask = 0xff;
  new_tft->packetfiltercontents.ipv4remoteaddr[3].mask = 0xff;
  new_tft->packetfiltercontents.singlelocalport = 22334;
  // rule id
  strncpy(reinterpret_cast<char*>(
              response->qos_flow_list.item[0].qos_flow_req_item.rule_id),
          "rule2", 6);
}

/* Create pdu session release message from ue */
status_code_e send_uplink_nas_pdu_session_release_message(
    amf_app_desc_t* amf_app_desc_p, amf_ue_ngap_id_t ue_id, const plmn_t& plmn,
    const uint8_t* nas_msg, uint8_t nas_msg_length) {
  bstring pdu_session_req;
  tai_t originating_tai = {};

  if ((!amf_app_desc_p) || (!nas_msg) || (nas_msg_length == 0)) {
    return RETURNerror;
  }

  pdu_session_req = blk2bstr(nas_msg, nas_msg_length);

  originating_tai.plmn = plmn;
  originating_tai.tac = 1;

  status_code_e rc = RETURNerror;
  rc = amf_app_handle_uplink_nas_message(amf_app_desc_p, pdu_session_req, ue_id,
                                         originating_tai);

  return rc;
}

/* Create ue deregistration request */
status_code_e send_uplink_nas_ue_deregistration_request(
    amf_app_desc_t* amf_app_desc_p, amf_ue_ngap_id_t ue_id, const plmn_t& plmn,
    uint8_t* nas_msg, uint8_t nas_msg_length) {
  bstring uplink_nas_ue_dereg_req;
  tai_t originating_tai = {};
  status_code_e rc = RETURNerror;

  tmsi_t ue_tmsi = htonl(amf_lookup_guti_by_ueid(ue_id));
  if (ue_tmsi == 0) {
    return (rc);
  }

  uplink_nas_ue_dereg_req = blk2bstr(nas_msg, nas_msg_length);

  /* Last 4 bytes are 0. Replace with TMSI */
  memcpy(&(uplink_nas_ue_dereg_req->data[nas_msg_length - sizeof(tmsi_t)]),
         &ue_tmsi, sizeof(tmsi_t));

  originating_tai.plmn = plmn;
  originating_tai.tac = 1;

  rc = amf_app_handle_uplink_nas_message(
      amf_app_desc_p, uplink_nas_ue_dereg_req, ue_id, originating_tai);

  return (rc);
}

/* Get the ue id from IMSI */
bool get_ue_id_from_imsi(amf_app_desc_t* amf_app_desc_p, imsi64_t imsi64,
                         amf_ue_ngap_id_t* ue_id) {
  amf_ue_context_t* amf_ue_context_p = &amf_app_desc_p->amf_ue_contexts;
  magma::map_rc_t rc_map = magma::MAP_OK;
  rc_map = amf_ue_context_p->imsi_amf_ue_id_htbl.get(imsi64, ue_id);
  if (rc_map != magma::MAP_OK) {
    return (false);
  }
  return true;
}

/* Create context release request */
void send_ue_context_release_request_message(amf_app_desc_t* amf_app_desc_p,
                                             uint32_t gnb_id,
                                             gnb_ue_ngap_id_t gnb_ue_ngap_id,
                                             amf_ue_ngap_id_t amf_ue_ngap_id) {
  itti_ngap_ue_context_release_req_t ue_context_release_request = {};

  ue_context_release_request.gnb_id = gnb_id;
  ue_context_release_request.gnb_ue_ngap_id = gnb_ue_ngap_id;
  ue_context_release_request.amf_ue_ngap_id = amf_ue_ngap_id;
  ue_context_release_request.relCause = NGAP_RADIO_NR_GENERATED_REASON;

  amf_app_handle_ngap_ue_context_release_req(&ue_context_release_request);
}

/* Create context release request */
void send_ue_context_release_complete_message(amf_app_desc_t* amf_app_desc_p,
                                              uint32_t gnb_id,
                                              gnb_ue_ngap_id_t gnb_ue_ngap_id,
                                              amf_ue_ngap_id_t amf_ue_ngap_id) {
  itti_ngap_ue_context_release_complete_t ue_context_release_complete = {};

  ue_context_release_complete.gnb_id = gnb_id;
  ue_context_release_complete.gnb_ue_ngap_id = gnb_ue_ngap_id;
  ue_context_release_complete.amf_ue_ngap_id = amf_ue_ngap_id;

  amf_app_handle_ngap_ue_context_release_complete(amf_app_desc_p,
                                                  &ue_context_release_complete);
}

imsi64_t send_initial_ue_message_service_request(
    amf_app_desc_t* amf_app_desc_p, sctp_assoc_id_t sctp_assoc_id,
    uint32_t gnb_id, gnb_ue_ngap_id_t gnb_ue_ngap_id,
    amf_ue_ngap_id_t amf_ue_ngap_id, const plmn_t& plmn, const uint8_t* nas_msg,
    uint8_t nas_msg_length, uint8_t tmsi_offset) {
  tai_t originating_tai = {};
  int rc = RETURNerror;
  imsi64_t imsi64 = 0;
  tai_t tai = {.plmn = plmn, .tac = 1};

  tmsi_t ue_tmsi = amf_lookup_guti_by_ueid(amf_ue_ngap_id);
  if (ue_tmsi == 0) {
    return (rc);
  }

  itti_ngap_initial_ue_message_t initial_ue_message = {};

  initial_ue_message.nas = blk2bstr(nas_msg, nas_msg_length);

  /* Replace TMSI value in nas message
   * message has uplink data status(4 bytess)
   * and PDU session status(4 bytes) after TMSI*/
  ue_tmsi = htonl(ue_tmsi);
  memcpy(&(initial_ue_message.nas
               ->data[nas_msg_length - sizeof(tmsi_t) - tmsi_offset]),
         &(ue_tmsi), sizeof(tmsi_t));
  initial_ue_message.sctp_assoc_id = sctp_assoc_id;
  initial_ue_message.gnb_ue_ngap_id = gnb_ue_ngap_id;
  initial_ue_message.gnb_id = gnb_id;
  initial_ue_message.m5g_rrc_establishment_cause = M5G_MO_SIGNALING;
  initial_ue_message.is_s_tmsi_valid = true;
  initial_ue_message.opt_s_tmsi.amf_set_id = 1;
  initial_ue_message.opt_s_tmsi.amf_pointer = 0;
  initial_ue_message.opt_s_tmsi.m_tmsi = ntohl(ue_tmsi);
  initial_ue_message.tai = tai;
  initial_ue_message.ue_context_request = M5G_UEContextRequest_requested;

  imsi64 =
      amf_app_handle_initial_ue_message(amf_app_desc_p, &initial_ue_message);

  return imsi64;
}

status_code_e send_uplink_nas_message_service_request_with_pdu(
    amf_app_desc_t* amf_app_desc_p, amf_ue_ngap_id_t amf_ue_ngap_id,
    const plmn_t& plmn, const uint8_t* nas_msg, uint8_t nas_msg_length) {
  bstring uplink_nas_service_request;
  tai_t originating_tai = {.plmn = plmn, .tac = 1};
  status_code_e rc = RETURNerror;
  tmsi_t ue_tmsi = amf_lookup_guti_by_ueid(amf_ue_ngap_id);
  if (ue_tmsi == 0) {
    return rc;
  }

  uplink_nas_service_request = blk2bstr(nas_msg, nas_msg_length);
  /* Replace TMSI value in nas message
   * message has uplink data status(4 bytess)
   * and PDU session status(4 bytes) after TMSI*/
  ue_tmsi = htonl(ue_tmsi);
  memcpy(
      &(uplink_nas_service_request->data[nas_msg_length - sizeof(tmsi_t) - 8]),
      &(ue_tmsi), sizeof(tmsi_t));

  rc = amf_app_handle_uplink_nas_message(amf_app_desc_p,
                                         uplink_nas_service_request,
                                         amf_ue_ngap_id, originating_tai);
  return rc;
}

// Check the ue context state
int check_ue_context_state(amf_ue_ngap_id_t ue_id,
                           m5gmm_state_t expected_mm_state,
                           m5gcm_state_t expected_cm_state) {
  OAILOG_FUNC_IN(LOG_AMF_APP);
  m5gmm_state_t mm_state;
  if (amf_get_ue_context_mm_state(ue_id, &mm_state) != RETURNok) {
    OAILOG_ERROR(LOG_AMF_APP,
                 "Error: amf_ue_context_mm_context does not exist, "
                 "ue_id: " AMF_UE_NGAP_ID_FMT "\n",
                 ue_id);
    OAILOG_FUNC_RETURN(LOG_AMF_APP, RETURNerror);
  }

  if (mm_state != expected_mm_state) {
    OAILOG_FUNC_RETURN(LOG_AMF_APP, RETURNerror);
  }

  m5gcm_state_t cm_state;
  if (amf_get_ue_context_cm_state(ue_id, &cm_state) != RETURNok) {
    OAILOG_ERROR(LOG_AMF_APP,
                 "Error: amf_ue_context_cm_context does not exist, "
                 "ue_id: " AMF_UE_NGAP_ID_FMT "\n",
                 ue_id);
    OAILOG_FUNC_RETURN(LOG_AMF_APP, RETURNerror);
  }

  if (cm_state != expected_cm_state) {
    OAILOG_FUNC_RETURN(LOG_AMF_APP, RETURNerror);
  }

  n2cause_e ue_context_rel_cause;
  if (amf_get_ue_context_rel_cause(ue_id, &ue_context_rel_cause) != RETURNok) {
    OAILOG_FUNC_RETURN(LOG_AMF_APP, RETURNerror);
  }

  if (ue_context_rel_cause != NGAP_RADIO_NR_GENERATED_REASON) {
    OAILOG_FUNC_RETURN(LOG_AMF_APP, RETURNerror);
  }

  OAILOG_FUNC_RETURN(LOG_AMF_APP, RETURNok);
}

// 5th expiry of t3550 during registration complete from UE
// mimicing registration_accept_t3550_handler
int unit_test_registration_accept_t3550(amf_ue_ngap_id_t ue_id) {
  int rc = RETURNerror;
  ue_m5gmm_context_s* ue_amf_context = NULL;

  // assuming 5 times expiry of T3550 timer for registration accept
  // Get the UE context
  ue_amf_context = amf_ue_context_exists_amf_ue_ngap_id(ue_id);
  if (ue_amf_context == NULL) {
    return RETURNerror;
  }

  // 5.5.1.2.8 abnormal case on network side
  // at 5th expiry of timer, amf enters into REGISTERED state
  rc = ue_state_handle_message_initial(
      COMMON_PROCEDURE_INITIATED2, STATE_EVENT_REG_COMPLETE, SESSION_NULL,
      ue_amf_context, &ue_amf_context->amf_context);

  return (rc);
}

// Send GNB Reset Request
void send_gnb_reset_req() {
  itti_ngap_gnb_initiated_reset_req_t reset_req_msg = {};
  reset_req_msg.ngap_reset_type = M5G_RESET_ALL;
  reset_req_msg.gnb_id = 1;
  reset_req_msg.sctp_assoc_id = 1;
  reset_req_msg.sctp_stream_id = 1;
  reset_req_msg.num_ue = 1;
  reset_req_msg.ue_to_reset_list =
      reinterpret_cast<ng_sig_conn_id_t*>(calloc(1, sizeof(ng_sig_conn_id_t)));
  reset_req_msg.ue_to_reset_list[0].amf_ue_ngap_id = 1;
  reset_req_msg.ue_to_reset_list[0].gnb_ue_ngap_id = 1;

  amf_app_handle_gnb_reset_req(&reset_req_msg);

  free(reset_req_msg.ue_to_reset_list);
}
}  // namespace magma5g

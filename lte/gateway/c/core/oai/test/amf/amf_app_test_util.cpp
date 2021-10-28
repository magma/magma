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
#include "amf_app_test_util.h"

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
  memcpy(initial_ue_message.nas->data, nas_msg, nas_msg_length);

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
  memcpy(uplink_nas_auth_response->data, nas_msg, nas_msg_length);

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
  memcpy(uplink_nas_smc_response->data, nas_msg, nas_msg_length);

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
  memcpy(ue_registration_complete->data, nas_msg, nas_msg_length);

  originating_tai.plmn = plmn;
  originating_tai.tac  = 1;

  int rc = RETURNerror;
  rc     = amf_app_handle_uplink_nas_message(
      amf_app_desc_p, ue_registration_complete, ue_id, originating_tai);

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
  memcpy(uplink_nas_ue_dereg_req->data, nas_msg, nas_msg_length);

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
  hashtable_rc_t rc_hash             = HASH_TABLE_OK;
  rc_hash                            = hashtable_uint64_ts_get(
      amf_ue_context_p->imsi_amf_ue_id_htbl, imsi64, ue_id);
  if (rc_hash != HASH_TABLE_OK) {
    return (false);
  }

  return true;
}

}  // namespace magma5g

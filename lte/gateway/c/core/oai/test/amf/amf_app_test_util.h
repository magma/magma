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
#include <string>
#include <gmock/gmock.h>
#include <gtest/gtest.h>

extern "C" {
#include "lte/gateway/c/core/oai/include/amf_config.hpp"
}
#include "lte/gateway/c/core/oai/include/amf_app_messages_types.h"
#include "lte/gateway/c/core/oai/tasks/amf/amf_authentication.hpp"
#include "lte/gateway/c/core/oai/tasks/amf/amf_app_defs.hpp"
#include "lte/gateway/c/core/oai/include/map.h"

using ::testing::_;
using ::testing::Return;

namespace magma5g {

/* Utility : Get ue_id from imsi */
bool get_ue_id_from_imsi(amf_app_desc_t* amf_app_desc_p, imsi64_t imsi64,
                         amf_ue_ngap_id_t* ue_id);

/* API for creating intial UE message without TMSI */
imsi64_t send_initial_ue_message_no_tmsi(
    amf_app_desc_t* amf_app_desc_p, sctp_assoc_id_t sctp_assoc_id,
    uint32_t gnb_id, gnb_ue_ngap_id_t gnb_ue_ngap_id,
    amf_ue_ngap_id_t amf_ue_ngap_id, const plmn_t& plmn, const uint8_t* nas_msg,
    uint8_t nas_msg_length);

imsi64_t send_initial_ue_message_no_tmsi_replace_mtmsi(
    amf_app_desc_t* amf_app_desc_p, sctp_assoc_id_t sctp_assoc_id,
    uint32_t gnb_id, gnb_ue_ngap_id_t gnb_ue_ngap_id,
    amf_ue_ngap_id_t amf_ue_ngap_id, const plmn_t& plmn, const uint8_t* nas_msg,
    uint8_t nas_msg_length, amf_ue_ngap_id_t ue_id, uint8_t tmsi_offset);

/* API for creating initial UE message without TMSI */
imsi64_t send_initial_ue_message_no_tmsi_no_ctx_req(
    amf_app_desc_t* amf_app_desc_p, sctp_assoc_id_t sctp_assoc_id,
    uint32_t gnb_id, gnb_ue_ngap_id_t gnb_ue_ngap_id,
    amf_ue_ngap_id_t amf_ue_ngap_id, const plmn_t& plmn, const uint8_t* nas_msg,
    uint8_t nas_msg_length);

/* For guti based registration */
uint64_t send_initial_ue_message_with_tmsi(
    amf_app_desc_t* amf_app_desc_p, sctp_assoc_id_t sctp_assoc_id,
    uint32_t gnb_id, gnb_ue_ngap_id_t gnb_ue_ngap_id,
    amf_ue_ngap_id_t amf_ue_ngap_id, const plmn_t& plmn, uint32_t m_tmsi,
    const uint8_t* nas_msg, uint8_t nas_msg_length);

/* For generating the identity response message */
status_code_e send_uplink_nas_identity_response_message(
    amf_app_desc_t* amf_app_desc_p, amf_ue_ngap_id_t ue_id, const plmn_t& plmn,
    const uint8_t* nas_msg, uint8_t nas_msg_length);

imsi64_t send_initial_ue_message_service_request(
    amf_app_desc_t* amf_app_desc_p, sctp_assoc_id_t sctp_assoc_id,
    uint32_t gnb_id, gnb_ue_ngap_id_t gnb_ue_ngap_id,
    amf_ue_ngap_id_t amf_ue_ngap_id, const plmn_t& plmn, const uint8_t* nas_msg,
    uint8_t nas_msg_length, uint8_t tmsi_offset);

status_code_e send_uplink_nas_message_service_request_with_pdu(
    amf_app_desc_t* amf_app_desc_p, amf_ue_ngap_id_t amf_ue_ngap_id,
    const plmn_t& plmn, const uint8_t* nas_msg, uint8_t nas_msg_length);

/* API for creating subscriberdb auth answer */
status_code_e send_proc_authentication_info_answer(const std::string& imsi,
                                                   amf_ue_ngap_id_t ue_id,
                                                   bool success);

/* API for creating uplink nas message for Auth Response */
status_code_e send_uplink_nas_message_ue_auth_response(
    amf_app_desc_t* amf_app_desc_p, amf_ue_ngap_id_t amf_ue_ngap_id,
    const plmn_t& plmn, const uint8_t* nas_msg, uint8_t nas_msg_length);

/* API for creating uplink nas message for security mode complete response */
status_code_e send_uplink_nas_message_ue_smc_response(
    amf_app_desc_t* amf_app_desc_p, amf_ue_ngap_id_t ue_id, const plmn_t& plmn,
    const uint8_t* nas_msg, uint8_t nas_msg_length);

/* API for sending initial context setup response */
void send_initial_context_response(amf_app_desc_t* amf_app_desc_p,
                                   amf_ue_ngap_id_t ue_id);

/* API for creating uplink nas message for registration complete response */
status_code_e send_uplink_nas_registration_complete(
    amf_app_desc_t* amf_app_desc_p, amf_ue_ngap_id_t ue_id, const plmn_t& plmn,
    const uint8_t* nas_msg, uint8_t nas_msg_length);

/* Create pdu session establishment  request from ue */
status_code_e send_uplink_nas_pdu_session_establishment_request(
    amf_app_desc_t* amf_app_desc_p, amf_ue_ngap_id_t ue_id, const plmn_t& plmn,
    const uint8_t* nas_msg, uint8_t nas_msg_length);

void create_ip_address_response_itti(
    itti_amf_ip_allocation_response_t* response);

status_code_e send_ip_address_response_itti(pdn_type_value_t type);

void create_pdu_session_response_ipv4_itti(
    itti_n11_create_pdu_session_response_t* response);

status_code_e send_pdu_session_response_itti(pdn_type_value_t type);

void create_pdu_resource_setup_response_itti(
    itti_ngap_pdusessionresource_setup_rsp_t* response, amf_ue_ngap_id_t ue_id);

void create_pdu_session_modify_request_itti(
    itti_n11_create_pdu_session_response_t* response);

void create_pdu_session_modify_deletion_request_itti(
    itti_n11_create_pdu_session_response_t* response);

status_code_e send_pdu_resource_setup_response(amf_ue_ngap_id_t ue_id);

void create_pdu_notification_response_itti(
    itti_n11_received_notification_t* response);

status_code_e send_pdu_notification_response();

/* Create pdu session  release from ue */
status_code_e send_uplink_nas_pdu_session_release_message(
    amf_app_desc_t* amf_app_desc_p, amf_ue_ngap_id_t ue_id, const plmn_t& plmn,
    const uint8_t* nas_msg, uint8_t nas_msg_length);

status_code_e send_uplink_nas_ue_deregistration_request(
    amf_app_desc_t* amf_app_desc_p, amf_ue_ngap_id_t ue_id, const plmn_t& plmn,
    uint8_t* nas_msg, uint8_t nas_msg_length);

/* Send ue context release message */
void send_ue_context_release_request_message(amf_app_desc_t* amf_app_desc_p,
                                             uint32_t gnb_id,
                                             gnb_ue_ngap_id_t gnb_ue_ngap_id,
                                             amf_ue_ngap_id_t amf_ue_ngap_id);

/* Send ue context release complete message */
void send_ue_context_release_complete_message(amf_app_desc_t* amf_app_desc_p,
                                              uint32_t gnb_id,
                                              gnb_ue_ngap_id_t gnb_ue_ngap_id,
                                              amf_ue_ngap_id_t amf_ue_ngap_id);

// Check the ue context state
int check_ue_context_state(amf_ue_ngap_id_t ue_id,
                           m5gmm_state_t expected_mm_state,
                           m5gcm_state_t expected_cm_state);

// mimicing registration_accept_t3550_handler
int unit_test_registration_accept_t3550(amf_ue_ngap_id_t ue_id);

void create_pdu_resource_modify_response_itti(
    itti_ngap_pdu_session_resource_modify_response_t* response,
    amf_ue_ngap_id_t ue_id);

int send_uplink_nas_pdu_session_modification_complete(
    amf_app_desc_t* amf_app_desc_p, amf_ue_ngap_id_t ue_id, const plmn_t& plmn,
    const uint8_t* nas_msg, uint8_t nas_msg_length);

int send_pdu_resource_modify_response(amf_ue_ngap_id_t ue_id);

int send_pdu_session_modification_itti();

int send_pdu_session_modification_deletion_itti();

// Send GNB Reset Request
void send_gnb_reset_req();
}  // namespace magma5g

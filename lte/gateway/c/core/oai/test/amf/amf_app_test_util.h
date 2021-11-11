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
#include "lte/gateway/c/core/oai/include/amf_config.h"
}
#include "lte/gateway/c/core/oai/include/amf_app_messages_types.h"
#include "lte/gateway/c/core/oai/tasks/amf/amf_authentication.h"
#include "lte/gateway/c/core/oai/tasks/amf/amf_app_defs.h"

using ::testing::_;
using ::testing::Return;

namespace magma5g {

/* Utility : Get ue_id from imsi */
bool get_ue_id_from_imsi(
    amf_app_desc_t* amf_app_desc_p, imsi64_t imsi64, amf_ue_ngap_id_t* ue_id);

/* API for creating intial UE message without TMSI */
imsi64_t send_initial_ue_message_no_tmsi(
    amf_app_desc_t* amf_app_desc_p, sctp_assoc_id_t sctp_assoc_id,
    uint32_t gnb_id, gnb_ue_ngap_id_t gnb_ue_ngap_id,
    amf_ue_ngap_id_t amf_ue_ngap_id, const plmn_t& plmn, const uint8_t* nas_msg,
    uint8_t nas_msg_length);

/* API for creating subscriberdb auth answer */
int send_proc_authentication_info_answer(
    const std::string& imsi, amf_ue_ngap_id_t ue_id, bool success);

/* API for creating uplink nas message for Auth Response */
int send_uplink_nas_message_ue_auth_response(
    amf_app_desc_t* amf_app_desc_p, amf_ue_ngap_id_t amf_ue_ngap_id,
    const plmn_t& plmn, const uint8_t* nas_msg, uint8_t nas_msg_length);

/* API for creating uplink nas message for security mode complete response */
int send_uplink_nas_message_ue_smc_response(
    amf_app_desc_t* amf_app_desc_p, amf_ue_ngap_id_t ue_id, const plmn_t& plmn,
    const uint8_t* nas_msg, uint8_t nas_msg_length);

/* API for sending initial context setup response */
void send_initial_context_response(
    amf_app_desc_t* amf_app_desc_p, amf_ue_ngap_id_t ue_id);

/* API for creating uplink nas message for registration complete response */
int send_uplink_nas_registration_complete(
    amf_app_desc_t* amf_app_desc_p, amf_ue_ngap_id_t ue_id, const plmn_t& plmn,
    const uint8_t* nas_msg, uint8_t nas_msg_length);

/* Create pdu session establishment  request from ue */
int send_uplink_nas_pdu_session_establishment_request(
    amf_app_desc_t* amf_app_desc_p, amf_ue_ngap_id_t ue_id, const plmn_t& plmn,
    const uint8_t* nas_msg, uint8_t nas_msg_length);

void create_ip_address_response_itti(
    itti_amf_ip_allocation_response_t& response);

int send_ip_address_response_itti();

void create_pdu_session_response_ipv4_itti(
    itti_n11_create_pdu_session_response_t& response);

int send_pdu_session_response_itti();

void create_pdu_resource_setup_response_itti(
    itti_ngap_pdusessionresource_setup_rsp_t& response, amf_ue_ngap_id_t ue_id);

int send_pdu_resource_setup_response(amf_ue_ngap_id_t ue_id);

void create_pdu_notification_response_itti(
    itti_n11_received_notification_t& response);

int send_pdu_notification_response();

/* Create pdu session establishment  release from ue */
int send_uplink_nas_pdu_session_release_request(
    amf_app_desc_t* amf_app_desc_p, amf_ue_ngap_id_t ue_id, const plmn_t& plmn,
    const uint8_t* nas_msg, uint8_t nas_msg_length);

int send_uplink_nas_ue_deregistration_request(
    amf_app_desc_t* amf_app_desc_p, amf_ue_ngap_id_t ue_id, const plmn_t& plmn,
    uint8_t* nas_msg, uint8_t nas_msg_length);

}  // namespace magma5g

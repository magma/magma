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

#pragma once
#include <sstream>
#include "lte/gateway/c/core/oai/include/amf_as_message.h"
#include "lte/gateway/c/core/oai/include/ngap_messages_types.h"
#include "lte/gateway/c/core/oai/include/n11_messages_types.h"
#include "lte/gateway/c/core/oai/tasks/amf/amf_app_ue_context_and_proc.hpp"

namespace magma5g {
typedef struct amf_app_desc_s {
  // UE contexts
  amf_ue_context_t amf_ue_contexts;
  amf_ue_ngap_id_t amf_app_ue_ngap_id_generator;
  long m5_statistic_timer_id;
  uint32_t m5_statistic_timer_period;

  /****************Statistics**************/
  // Number of connected UEs
  uint32_t nb_ue_connected;
  // Number of Attached UEs
  uint32_t nb_ue_attached;
  // Number of Active Pdu session
  uint32_t nb_pdu_sessions;
  // Number of Idle UEs
  uint32_t nb_ue_idle;
} amf_app_desc_t;

// UL and DL routines.
imsi64_t amf_app_handle_initial_ue_message(
    amf_app_desc_t* amf_app_desc_p,
    itti_ngap_initial_ue_message_t* conn_est_ind_pP);
status_code_e amf_app_handle_nas_dl_req(amf_ue_ngap_id_t ue_id, bstring nas_msg,
                                        nas5g_error_code_t transaction_status);
status_code_e amf_app_handle_uplink_nas_message(amf_app_desc_t* amf_app_desc_p,
                                                bstring msg,
                                                amf_ue_ngap_id_t ue_id,
                                                const tai_t originating_tai);
status_code_e amf_app_handle_pdu_session_response(
    itti_n11_create_pdu_session_response_t* pdu_session_resp);
status_code_e amf_app_handle_pdu_session_failure(
    itti_n11_create_pdu_session_failure_t* pdu_session_failure);
status_code_e amf_app_handle_notification_received(
    itti_n11_received_notification_t* notification);
status_code_e amf_app_handle_pdu_session_accept(
    itti_n11_create_pdu_session_response_t* pdu_session_resp, uint64_t ue_id);
void convert_ambr(const uint32_t* pdu_ambr_response_unit,
                  const uint32_t* pdu_ambr_response_value,
                  M5GSessionAmbrUnit* ambr_unit, uint16_t* ambr_value);
status_code_e amf_smf_handle_ip_address_response(
    itti_amf_ip_allocation_response_t* response_p);
void amf_app_handle_initial_context_setup_rsp(
    amf_app_desc_t* amf_app_desc_p,
    itti_amf_app_initial_context_setup_rsp_t* initial_context_rsp);
status_code_e amf_send_n11_update_location_req(amf_ue_ngap_id_t ue_id);
int amf_app_pdu_session_modification_request(
    itti_n11_create_pdu_session_response_t* pdu_sess_mod_req,
    amf_ue_ngap_id_t ue_id);
int amf_app_pdu_session_modification_complete(amf_smf_establish_t* message,
                                              char* imsi, uint32_t version);
int amf_app_pdu_session_modification_command_reject(
    amf_smf_establish_t* message, char* imsi, uint32_t version);
std::string get_message_type_str(uint8_t type);

// Handling ue context release complete
void amf_app_handle_ngap_ue_context_release_complete(
    amf_app_desc_t* amf_app_desc_p,
    const itti_ngap_ue_context_release_complete_t* const
        ngap_ue_context_release_complete);

// Handling of SCTP shutdown
void amf_app_handle_gnb_deregister_ind(
    const itti_ngap_gNB_deregistered_ind_t* const gNB_deregistered_ind);

void amf_app_handle_gnb_reset_req(
    const itti_ngap_gnb_initiated_reset_req_t* const gnb_reset_req);
}  // namespace magma5g

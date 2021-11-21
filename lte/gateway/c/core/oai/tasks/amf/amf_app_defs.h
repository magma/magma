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
#include "lte/gateway/c/core/oai/tasks/amf/amf_app_ue_context_and_proc.h"

namespace magma5g {
typedef struct amf_app_desc_s {
  // UE contexts
  amf_ue_context_t amf_ue_contexts;
  amf_ue_ngap_id_t amf_app_ue_ngap_id_generator;
  long m5_statistic_timer_id;
  uint32_t m5_statistic_timer_period;
} amf_app_desc_t;

// UL and DL routines.
imsi64_t amf_app_handle_initial_ue_message(
    amf_app_desc_t* amf_app_desc_p,
    itti_ngap_initial_ue_message_t* conn_est_ind_pP);
int amf_app_handle_nas_dl_req(
    amf_ue_ngap_id_t ue_id, bstring nas_msg,
    nas5g_error_code_t transaction_status);
int amf_app_handle_uplink_nas_message(
    amf_app_desc_t* amf_app_desc_p, bstring msg, amf_ue_ngap_id_t ue_id,
    const tai_t originating_tai);
int amf_app_handle_pdu_session_response(
    itti_n11_create_pdu_session_response_t* pdu_session_resp);
int amf_app_handle_notification_received(
    itti_n11_received_notification_t* notification);
int amf_app_handle_pdu_session_accept(
    itti_n11_create_pdu_session_response_t* pdu_session_resp, uint64_t ue_id);
void convert_ambr(
    const uint32_t* pdu_ambr_response_unit,
    const uint32_t* pdu_ambr_response_value, uint8_t* ambr_unit,
    uint16_t* ambr_value);
int amf_smf_handle_ip_address_response(
    itti_amf_ip_allocation_response_t* response_p);
void amf_app_handle_initial_context_setup_rsp(
    amf_app_desc_t* amf_app_desc_p,
    itti_amf_app_initial_context_setup_rsp_t* initial_context_rsp);
int amf_send_n11_update_location_req(amf_ue_ngap_id_t ue_id);
}  // namespace magma5g

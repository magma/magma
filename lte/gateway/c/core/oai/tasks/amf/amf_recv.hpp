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
#include "lte/gateway/c/core/oai/tasks/amf/amf_app_ue_context_and_proc.hpp"
#include "lte/gateway/c/core/oai/tasks/nas5g/include/M5GRegistrationAccept.hpp"
#include "lte/gateway/c/core/oai/tasks/amf/amf_asDefs.h"
#include "lte/gateway/c/core/oai/tasks/amf/amf_identity.hpp"

namespace magma5g {
// AMF registration procedures
status_code_e amf_handle_registration_request(
    amf_ue_ngap_id_t ue_id, tai_t* originating_tai, ecgi_t* ecgi,
    RegistrationRequestMsg* msg, const bool is_initial,
    const bool is_amf_ctx_new, int amf_cause,
    amf_nas_message_decode_status_t decode_status);
status_code_e amf_handle_service_request(
    amf_ue_ngap_id_t ue_id, ServiceRequestMsg* msg,
    const amf_nas_message_decode_status_t decode_status);
status_code_e amf_registration_run_procedure(amf_context_t* amf_context);
status_code_e amf_handle_identity_response(
    amf_ue_ngap_id_t ue_id, M5GSMobileIdentityMsg* msg, int amf_cause,
    amf_nas_message_decode_status_t decode_status);
status_code_e amf_handle_authentication_response(
    amf_ue_ngap_id_t ue_id, AuthenticationResponseMsg* msg, int amf_cause,
    amf_nas_message_decode_status_t status);
status_code_e amf_handle_authentication_failure(
    amf_ue_ngap_id_t ue_id, AuthenticationFailureMsg* msg, int amf_cause,
    amf_nas_message_decode_status_t status);
status_code_e amf_handle_security_complete_response(
    amf_ue_ngap_id_t ue_id, amf_nas_message_decode_status_t decode_status);
status_code_e amf_handle_security_mode_reject(
    const amf_ue_ngap_id_t ueid, SecurityModeRejectMsg* msg,
    int const amf_cause, const amf_nas_message_decode_status_t decode_status);
status_code_e amf_handle_registration_complete_response(
    amf_ue_ngap_id_t ue_id, RegistrationCompleteMsg* msg, int amf_cause,
    amf_nas_message_decode_status_t decode_status);
status_code_e amf_handle_deregistration_ue_origin_req(
    amf_ue_ngap_id_t ue_id, DeRegistrationRequestUEInitMsg* msg, int amf_cause,
    amf_nas_message_decode_status_t decode_status);
status_code_e amf_validate_dnn(const amf_context_s* amf_ctxt_p,
                               std::string dnn_string, int* index,
                               bool ue_sent_dnn);
void smf_dnn_ambr_select(const std::shared_ptr<smf_context_t>& smf_ctx,
                         ue_m5gmm_context_s* ue_context, int index_dnn);
status_code_e amf_smf_process_pdu_session_packet(amf_ue_ngap_id_t ueid,
                                                 ULNASTransportMsg* msg,
                                                 int amf_cause);
status_code_e amf_smf_notification_send(amf_ue_ngap_id_t ueid,
                                        ue_m5gmm_context_s* ue_context,
                                        notify_ue_event notify_event_type,
                                        uint16_t session_id);
status_code_e amf_proc_registration_request(
    amf_ue_ngap_id_t ue_id, const bool is_mm_ctx_new,
    amf_registration_request_ies_t* ies);
status_code_e amf_registration_success_identification_cb(
    amf_context_t* amf_context);
status_code_e amf_registration_failure_identification_cb(
    amf_context_t* amf_context);
status_code_e amf_registration_success_authentication_cb(
    amf_context_t* amf_context);
status_code_e amf_registration_success_security_cb(amf_context_t* amf_context);
status_code_e amf_proc_registration_reject(const amf_ue_ngap_id_t ue_id,
                                           amf_cause_t amf_cause);

status_code_e get_decrypt_imsi_suci_extension(amf_context_t* amf_context,
                                              uint8_t ue_pubkey_identifier,
                                              const std::string& ue_pubkey,
                                              const std::string& ciphertext,
                                              const std::string& mac_tag);
int amf_decrypt_msin_info_answer(itti_amf_decrypted_msin_info_ans_t* aia);
void amf_copy_plmn_to_supi(const ImsiM5GSMobileIdentity& imsi,
                           supi_as_imsi_t& supi_imsi);
status_code_e amf_copy_plmn_to_context(const ImsiM5GSMobileIdentity& imsi,
                                       ue_m5gmm_context_s* ue_context);
}  // namespace magma5g

/*
 * Licensed to the OpenAirInterface (OAI) Software Alliance under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The OpenAirInterface Software Alliance licenses this file to You under
 * the terms found in the LICENSE file in the root of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *-------------------------------------------------------------------------------
 * For more information about the OpenAirInterface (OAI) Software Alliance:
 *      contact@openairinterface.org
 */

#pragma once

#include <stdbool.h>

#ifdef __cplusplus
extern "C" {
#endif
#include "lte/gateway/c/core/oai/lib/itti/intertask_interface.h"
#ifdef __cplusplus
}
#endif

#include "S1ap_Cause.h"
#include "S1ap_S1AP-PDU.h"
#include "lte/gateway/c/core/common/common_defs.h"
#include "lte/gateway/c/core/oai/common/common_types.h"
#include "lte/gateway/c/core/oai/include/s1ap_messages_types.h"
#include "lte/gateway/c/core/oai/include/sctp_messages_types.hpp"
#include "lte/gateway/c/core/oai/tasks/s1ap/s1ap_mme.hpp"

namespace magma {
namespace lte {

#define MAX_NUM_PARTIAL_S1_CONN_RESET 256

const char* s1_enb_state2str(enum magma::lte::oai::S1apEnbState state);
const char* s1ap_direction2str(uint8_t dir);

/** \brief Handle decoded incoming messages from SCTP
 * \param assoc_id SCTP association ID
 * \param stream Stream number
 * \param message_p The message decoded by the ASN1C decoder
 * @returns int
 **/
status_code_e s1ap_mme_handle_message(oai::S1apState* state,
                                      const sctp_assoc_id_t assoc_id,
                                      const sctp_stream_id_t stream,
                                      S1ap_S1AP_PDU_t* message_p);

status_code_e s1ap_mme_handle_ue_cap_indication(oai::S1apState* state,
                                                const sctp_assoc_id_t assoc_id,
                                                const sctp_stream_id_t stream,
                                                S1ap_S1AP_PDU_t* message);

/** \brief Handle an S1 Setup request message.
 * Typically add the eNB in the list of served eNB if not present, simply reset
 * UEs association otherwise. S1SetupResponse message is sent in case of success
 * or S1SetupFailure if the MME cannot accept the configuration received.
 * \param assoc_id SCTP association ID
 * \param stream Stream number
 * \param message_p The message decoded by the ASN1C decoder
 * @returns int
 **/
status_code_e s1ap_mme_handle_s1_setup_request(oai::S1apState* state,
                                               const sctp_assoc_id_t assoc_id,
                                               const sctp_stream_id_t stream,
                                               S1ap_S1AP_PDU_t* message_p);

status_code_e s1ap_mme_handle_handover_required(oai::S1apState* state,
                                                const sctp_assoc_id_t assoc_id,
                                                const sctp_stream_id_t stream,
                                                S1ap_S1AP_PDU_t* message_p);

status_code_e s1ap_mme_handle_handover_command(
    oai::S1apState* state, const itti_mme_app_handover_command_t* ho_command_p);

status_code_e s1ap_mme_handle_handover_request(
    oai::S1apState* state, const itti_mme_app_handover_request_t* ho_request_p);

status_code_e s1ap_mme_handle_handover_request_ack(
    oai::S1apState* state, const sctp_assoc_id_t assoc_id,
    const sctp_stream_id_t stream, S1ap_S1AP_PDU_t* message_p);

status_code_e s1ap_mme_handle_handover_cancel(oai::S1apState* state,
                                              const sctp_assoc_id_t assoc_id,
                                              const sctp_stream_id_t stream,
                                              S1ap_S1AP_PDU_t* message_p);

status_code_e s1ap_mme_handle_handover_failure(oai::S1apState* state,
                                               const sctp_assoc_id_t assoc_id,
                                               const sctp_stream_id_t stream,
                                               S1ap_S1AP_PDU_t* message_p);

status_code_e s1ap_mme_handle_enb_status_transfer(
    oai::S1apState* state, const sctp_assoc_id_t assoc_id,
    const sctp_stream_id_t stream, S1ap_S1AP_PDU_t* pdu);

status_code_e s1ap_mme_handle_handover_notify(oai::S1apState* state,
                                              const sctp_assoc_id_t assoc_id,
                                              const sctp_stream_id_t stream,
                                              S1ap_S1AP_PDU_t* message_p);

status_code_e s1ap_mme_handle_path_switch_request(
    oai::S1apState* state, const sctp_assoc_id_t assoc_id,
    const sctp_stream_id_t stream, S1ap_S1AP_PDU_t* message_p);

status_code_e s1ap_mme_handle_ue_context_release_request(
    oai::S1apState* state, const sctp_assoc_id_t assoc_id,
    const sctp_stream_id_t stream, S1ap_S1AP_PDU_t* message_p);

status_code_e s1ap_handle_ue_context_release_command(
    oai::S1apState* state,
    const itti_s1ap_ue_context_release_command_t* const
        ue_context_release_command_pP,
    imsi64_t imsi64);

status_code_e s1ap_mme_handle_ue_context_release_complete(
    oai::S1apState* state, const sctp_assoc_id_t assoc_id,
    const sctp_stream_id_t stream, S1ap_S1AP_PDU_t* message_p);

status_code_e s1ap_handle_ue_context_mod_req(
    oai::S1apState* state,
    const itti_s1ap_ue_context_mod_req_t* const ue_context_mod_req_pP,
    imsi64_t imsi64);

status_code_e s1ap_mme_handle_initial_context_setup_failure(
    oai::S1apState* state, const sctp_assoc_id_t assoc_id,
    const sctp_stream_id_t stream, S1ap_S1AP_PDU_t* message_p);

status_code_e s1ap_mme_handle_initial_context_setup_response(
    oai::S1apState* state, const sctp_assoc_id_t assoc_id,
    const sctp_stream_id_t stream, S1ap_S1AP_PDU_t* message_p);

status_code_e s1ap_handle_sctp_disconnection(oai::S1apState* state,
                                             const sctp_assoc_id_t assoc_id,
                                             bool reset);

status_code_e s1ap_handle_new_association(oai::S1apState* state,
                                          sctp_new_peer_t* sctp_new_peer_p);

status_code_e s1ap_mme_set_cause(S1ap_Cause_t* cause_p,
                                 const S1ap_Cause_PR cause_type,
                                 const long cause_value);

long s1ap_mme_get_cause_value(S1ap_Cause_t* cause);

status_code_e s1ap_mme_generate_s1_setup_failure(const sctp_assoc_id_t assoc_id,
                                                 const S1ap_Cause_PR cause_type,
                                                 const long cause_value,
                                                 const long time_to_wait);

status_code_e s1ap_mme_handle_erab_setup_response(
    oai::S1apState* state, const sctp_assoc_id_t assoc_id,
    const sctp_stream_id_t stream, S1ap_S1AP_PDU_t* message);

status_code_e s1ap_mme_handle_erab_setup_failure(oai::S1apState* state,
                                                 const sctp_assoc_id_t assoc_id,
                                                 const sctp_stream_id_t stream,
                                                 S1ap_S1AP_PDU_t* message);

void s1ap_mme_release_ue_context(oai::S1apState* state,
                                 oai::UeDescription* ue_ref_p, imsi64_t imsi64);

status_code_e s1ap_mme_handle_error_ind_message(oai::S1apState* state,
                                                const sctp_assoc_id_t assoc_id,
                                                const sctp_stream_id_t stream,
                                                S1ap_S1AP_PDU_t* message);

status_code_e s1ap_mme_handle_enb_reset(oai::S1apState* state,
                                        const sctp_assoc_id_t assoc_id,
                                        const sctp_stream_id_t stream,
                                        S1ap_S1AP_PDU_t* message);

status_code_e s1ap_handle_enb_initiated_reset_ack(
    const itti_s1ap_enb_initiated_reset_ack_t* const enb_reset_ack_p,
    imsi64_t imsi64);

status_code_e s1ap_handle_paging_request(
    oai::S1apState* state, const itti_s1ap_paging_request_t* paging_request,
    imsi64_t imsi64);

status_code_e s1ap_mme_handle_ue_context_modification_response(
    oai::S1apState* state, const sctp_assoc_id_t assoc_id,
    const sctp_stream_id_t stream, S1ap_S1AP_PDU_t* message_p);

status_code_e s1ap_mme_handle_ue_context_modification_failure(
    oai::S1apState* state, const sctp_assoc_id_t assoc_id,
    const sctp_stream_id_t stream, S1ap_S1AP_PDU_t* message_p);

status_code_e s1ap_mme_handle_erab_rel_response(oai::S1apState* state,
                                                const sctp_assoc_id_t assoc_id,
                                                const sctp_stream_id_t stream,
                                                S1ap_S1AP_PDU_t* message);

status_code_e s1ap_mme_handle_enb_configuration_transfer(
    oai::S1apState* state, const sctp_assoc_id_t assoc_id,
    const sctp_stream_id_t stream, S1ap_S1AP_PDU_t* message_p);

status_code_e s1ap_handle_path_switch_req_ack(
    oai::S1apState* state,
    const itti_s1ap_path_switch_request_ack_t* path_switch_req_ack_p,
    imsi64_t imsi64);

status_code_e s1ap_handle_path_switch_req_failure(
    const itti_s1ap_path_switch_request_failure_t* path_switch_req_failure_p,
    imsi64_t imsi64);

status_code_e s1ap_mme_handle_erab_modification_indication(
    oai::S1apState* state, const sctp_assoc_id_t assoc_id,
    const sctp_stream_id_t stream, S1ap_S1AP_PDU_t* pdu);

void s1ap_mme_generate_erab_modification_confirm(
    oai::S1apState* state,
    const itti_s1ap_e_rab_modification_cnf_t* const conf);

status_code_e s1ap_mme_generate_ue_context_release_command(
    oai::S1apState* state, oai::UeDescription* ue_ref_p, enum s1cause,
    imsi64_t imsi64, sctp_assoc_id_t assoc_id, sctp_stream_id_t stream,
    mme_ue_s1ap_id_t mme_ue_s1ap_id, enb_ue_s1ap_id_t enb_ue_s1ap_id);

status_code_e s1ap_mme_generate_ue_context_modification(
    oai::UeDescription* ue_ref_p,
    const itti_s1ap_ue_context_mod_req_t* const ue_context_mod_req_pP,
    imsi64_t imsi64);

status_code_e s1ap_mme_remove_stale_ue_context(enb_ue_s1ap_id_t enb_ue_s1ap_id,
                                               uint32_t enb_id);

status_code_e s1ap_send_mme_ue_context_release(oai::S1apState* state,
                                               oai::UeDescription* ue_ref_p,
                                               enum s1cause s1_release_cause,
                                               S1ap_Cause_t ie_cause,
                                               imsi64_t imsi64);

}  // namespace lte
}  // namespace magma

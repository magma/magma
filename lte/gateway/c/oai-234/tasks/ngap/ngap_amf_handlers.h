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

#ifndef FILE_NGAP_MME_HANDLERS_SEEN
#define FILE_NGAP_MME_HANDLERS_SEEN
#include <stdbool.h>

#include "ngap_amf.h"
#include "intertask_interface.h"
#include "Ngap_Cause.h"
#include "common_types.h"
#include "ngap_messages_types.h"
#include "sctp_messages_types.h"

#define MAX_NUM_PARTIAL_NG_CONN_RESET 256

const char* ng_gnb_state2str(enum amf_ng_gnb_state_s state);
const char* ngap_direction2str(uint8_t dir);

/** \brief Handle decoded incoming messages from SCTP
 * \param assoc_id SCTP association ID
 * \param stream Stream number
 * \param message_p The message decoded by the ASN1C decoder
 * @returns int
 **/
int ngap_amf_handle_message(
    ngap_state_t* state, const sctp_assoc_id_t assoc_id,
    const sctp_stream_id_t stream, Ngap_NGAP_PDU_t* message_p);

int ngap_amf_handle_ue_cap_indication(
    ngap_state_t* state, const sctp_assoc_id_t assoc_id,
    const sctp_stream_id_t stream, Ngap_NGAP_PDU_t* message);

/** \brief Handle an S1 Setup request message.
 * Typically add the eNB in the list of served eNB if not present, simply reset
 * UEs association otherwise. S1SetupResponse message is sent in case of success
 *or S1SetupFailure if the MME cannot accept the configuration received. \param
 *assoc_id SCTP association ID \param stream Stream number \param message_p The
 *message decoded by the ASN1C decoder
 * @returns int
 **/
int ngap_amf_handle_ng_setup_request(
    ngap_state_t* state, const sctp_assoc_id_t assoc_id,
    const sctp_stream_id_t stream, Ngap_NGAP_PDU_t* message_p);

int ngap_amf_handle_path_switch_request(
    ngap_state_t* state, const sctp_assoc_id_t assoc_id,
    const sctp_stream_id_t stream, Ngap_NGAP_PDU_t* message_p);

int ngap_amf_handle_ue_context_release_request(
    ngap_state_t* state, const sctp_assoc_id_t assoc_id,
    const sctp_stream_id_t stream, Ngap_NGAP_PDU_t* message_p);

int ngap_handle_ue_context_release_command(
    ngap_state_t* state,
    const itti_ngap_ue_context_release_command_t* const
        ue_context_release_command_pP,
    imsi64_t imsi64);

int ngap_amf_handle_ue_context_release_complete(
    ngap_state_t* state, const sctp_assoc_id_t assoc_id,
    const sctp_stream_id_t stream, Ngap_NGAP_PDU_t* message_p);

int ngap_handle_ue_context_mod_req(
    ngap_state_t* state,
    const itti_ngap_ue_context_mod_req_t* const ue_context_mod_req_pP,
    imsi64_t imsi64);

int ngap_amf_handle_initial_context_setup_failure(
    ngap_state_t* state, const sctp_assoc_id_t assoc_id,
    const sctp_stream_id_t stream, Ngap_NGAP_PDU_t* message_p);

int ngap_amf_handle_initial_context_setup_response(
    ngap_state_t* state, const sctp_assoc_id_t assoc_id,
    const sctp_stream_id_t stream, Ngap_NGAP_PDU_t* message_p);

int ngap_handle_sctp_disconnection(
    ngap_state_t* state, const sctp_assoc_id_t assoc_id, bool reset);

int ngap_handle_new_association(
    ngap_state_t* state, sctp_new_peer_t* sctp_new_peer_p);

int ngap_amf_set_cause(
    Ngap_Cause_t* cause_p, const Ngap_Cause_PR cause_type,
    const long cause_value);

int ngap_amf_set_Rel_cause(
    Ngap_Cause_t* cause_p,
    pdu_session_resource_release_command_transfer* Rel_Cause);

int ngap_amf_generate_ng_setup_failure(
    const sctp_assoc_id_t assoc_id, const Ngap_Cause_PR cause_type,
    const long cause_value, const long time_to_wait);

int ngap_amf_handle_erab_setup_response(
    ngap_state_t* state, const sctp_assoc_id_t assoc_id,
    const sctp_stream_id_t stream, Ngap_NGAP_PDU_t* message);

int ngap_amf_handle_erab_setup_failure(
    ngap_state_t* state, const sctp_assoc_id_t assoc_id,
    const sctp_stream_id_t stream, Ngap_NGAP_PDU_t* message);

void ngap_amf_handle_ue_context_rel_comp_timer_expiry(
    ngap_state_t* state, m5g_ue_description_t* ue_ref_p);

void ngap_amf_release_ue_context(
    ngap_state_t* state, m5g_ue_description_t* ue_ref_p, imsi64_t imsi64);

//<<<<<<< HEAD
//=======
int ngap_amf_handle_error_ind_message(
    ngap_state_t* state, const sctp_assoc_id_t assoc_id,
    const sctp_stream_id_t stream, Ngap_NGAP_PDU_t* message);

//>>>>>>> Ngap_handlers_P2
int ngap_amf_handle_gnb_reset(
    ngap_state_t* state, const sctp_assoc_id_t assoc_id,
    const sctp_stream_id_t stream, Ngap_NGAP_PDU_t* message);

int ngap_handle_gnb_initiated_reset_ack(
    const itti_ngap_gnb_initiated_reset_ack_t* const gnb_reset_ack_p,
    imsi64_t imsi64);

int ngap_handle_paging_request(
    ngap_state_t* state, const itti_ngap_paging_request_t* paging_request,
    imsi64_t imsi64);

//<<<<<<< HEAD
//=======
int ngap_amf_handle_ue_context_modification_response(
    ngap_state_t* state, const sctp_assoc_id_t assoc_id,
    const sctp_stream_id_t stream, Ngap_NGAP_PDU_t* message_p);

int ngap_amf_handle_ue_context_modification_failure(
    ngap_state_t* state, const sctp_assoc_id_t assoc_id,
    const sctp_stream_id_t stream, Ngap_NGAP_PDU_t* message_p);

//>>>>>>> Ngap_handlers_P2
int ngap_amf_handle_erab_release_response(
    ngap_state_t* state, const sctp_assoc_id_t assoc_id,
    const sctp_stream_id_t stream, Ngap_NGAP_PDU_t* message);

int ngap_amf_handle_gnb_configuration_transfer(
    ngap_state_t* state, const sctp_assoc_id_t assoc_id,
    const sctp_stream_id_t stream, Ngap_NGAP_PDU_t* message_p);

int ngap_handle_path_switch_req_ack(
    ngap_state_t* state,
    const itti_ngap_path_switch_request_ack_t* path_switch_req_ack_p,
    imsi64_t imsi64);

int ngap_handle_path_switch_req_failure(
    ngap_state_t* state,
    const itti_ngap_path_switch_request_failure_t* path_switch_req_failure_p,
    imsi64_t imsi64);

int ngap_amf_handle_pduSession_setup_response(
    ngap_state_t* state, const sctp_assoc_id_t assoc_id,
    const sctp_stream_id_t stream, Ngap_NGAP_PDU_t* message_p);

int ngap_amf_handle_pduSession_setup_failure(
    ngap_state_t* state, const sctp_assoc_id_t assoc_id,
    const sctp_stream_id_t stream, Ngap_NGAP_PDU_t* message_p);

int ngap_amf_handle_pduSession_release_response(
    ngap_state_t* state, const sctp_assoc_id_t assoc_id,
    const sctp_stream_id_t stream, Ngap_NGAP_PDU_t* message_p);

#endif /* FILE_NGAP_MME_HANDLERS_SEEN */

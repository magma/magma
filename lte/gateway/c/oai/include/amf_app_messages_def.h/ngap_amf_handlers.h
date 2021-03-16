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
/****************************************************************************
  Source      ngap_amf_handlers.h
  Date        2020/07/28
  Subsystem   Access and Mobility Management Function
  Author      Ashish Prajapati
  Description Defines NG Application Protocol Messages Handlers

*****************************************************************************/
#pragma once

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

/** \brief Handle an N2 Setup request message.
 * Typically add the gNB in the list of served gNB if not present, simply reset
 * UEs association otherwise. N2SetupResponse message is sent in case of success
 *or N2SetupFailure if the AMF cannot accept the configuration received. \param
 *assoc_id SCTP association ID \param stream Stream number \param message_p The
 *message decoded by the ASN1C decoder
 * @returns int
 **/
int ngap_amf_handle_ng_setup_request(
    ngap_state_t* state, const sctp_assoc_id_t assoc_id,
    const sctp_stream_id_t stream, Ngap_NGAP_PDU_t* message_p);

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

int ngap_amf_generate_ng_setup_failure(
    const sctp_assoc_id_t assoc_id, const Ngap_Cause_PR cause_type,
    const long cause_value, const long time_to_wait);

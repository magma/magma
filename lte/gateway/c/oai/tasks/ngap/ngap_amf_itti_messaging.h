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

#include <stdbool.h>
#include <stddef.h>
#include <stdint.h>
#include <czmq.h>

#include "common_defs.h"
#include "3gpp_23.003.h"
#include "3gpp_38.401.h"
#include "3gpp_38.413.h"
#include "Ngap_Cause.h"
#include "TrackingAreaIdentity.h"
#include "bstrlib.h"
#include "common_types.h"
#include "intertask_interface.h"
#include "intertask_interface.h"
#include "ngap_state.h"

task_zmq_ctx_t ngap_task_zmq_ctx;

/** \brief pass msg to SCTP for transmit
 * \param payload msg to transmit
 * \param sctp_assoc_id_t SCTP association ID
 * \param stream Stream number
 * \param amf_ue_ngap_id_t amf_ue_ngap_id
 * @returns int
 **/
int ngap_amf_itti_send_sctp_request(
    STOLEN_REF bstring* payload, const uint32_t sctp_assoc_id_t,
    const sctp_stream_id_t stream, const amf_ue_ngap_id_t ue_id);

/** \brief pass NAS msg to AMF
 * \param amf_ue_ngap_id_t amf_ue_ngap_id
 * \param payload msg to transmit
 * \param tai Tracking Area Identifier
 * \param cgi E-UTRAN Cell Global Identification
 * @returns int
 **/
int ngap_amf_itti_nas_uplink_ind(
    const amf_ue_ngap_id_t ue_id, STOLEN_REF bstring* payload,
    const tai_t* const tai, const ecgi_t* const cgi);

/** \brief Handle initial_ue_message
 * \param assoc_id SCTP association ID
 * \param gnb_id gNB ID
 * \param gnb_ue_ngap_id gnb_ue_ngap_id
 * \param nas_msg NAS Msg
 * \param nas_msg_length NAS Msg Length
 * \param tai Tracking Area Identifier
 * \param cgi E-UTRAN Cell Global Identification
 * \param rrc_cause establishment cause
 * \param opt_s_tmsi shortened TMSI
 * \param opt_guamfi GUAMF Id
 * \param opt_cell_access_mode CELL ACCESS MODE
 * \param opt_cell_gw_transport_address GW Transport Layer Address
 * \param opt_relay_node_indicator Relay Node Indicator
 * @returns nothing
 **/
void ngap_amf_itti_ngap_initial_ue_message(
    const sctp_assoc_id_t assoc_id, const uint32_t gnb_id,
    const gnb_ue_ngap_id_t gnb_ue_ngap_id, const uint8_t* const nas_msg,
    const size_t nas_msg_length, const tai_t* const tai,
    const ecgi_t* const ecgi, const long rrc_cause,
    const s_tmsi_m5_t* const opt_s_tmsi, const csg_id_t* const opt_csg_id,
    const guamfi_t* const opt_guamfi,
    const void* opt_cell_access_mode,          /* unused*/
    const void* opt_cell_gw_transport_address, /* unused*/
    const void* opt_relay_node_indicator       /* unused*/
);

void ngap_amf_itti_nas_non_delivery_ind(
    const amf_ue_ngap_id_t ue_id, uint8_t* const nas_msg,
    const size_t nas_msg_length, const Ngap_Cause_t* const cause,
    const imsi64_t imsi64);

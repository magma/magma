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

/*! \file s1ap_mme_nas_procedures.h
  \brief
  \author Sebastien ROUX, Lionel Gauthier
  \company Eurecom
  \email: lionel.gauthier@eurecom.fr
*/

#ifndef FILE_S1AP_MME_NAS_PROCEDURES_SEEN
#define FILE_S1AP_MME_NAS_PROCEDURES_SEEN

#include "common_defs.h"
#include "3gpp_36.401.h"
#include "bstrlib.h"
#include "common_types.h"
#include "mme_app_messages_types.h"
#include "s1ap_messages_types.h"
#include "s1ap_state.h"

/** \brief Handle an Initial UE message.
 * \param assocId lower layer assoc id (SCTP)
 * \param stream SCTP stream on which data had been received
 * \param message The message as decoded by the ASN.1 codec
 * @returns -1 on failure, 0 otherwise
 **/
int s1ap_mme_handle_initial_ue_message(
    s1ap_state_t* state, const sctp_assoc_id_t assocId,
    const sctp_stream_id_t stream, S1ap_S1AP_PDU_t* message);

/** \brief Handle an Uplink NAS transport message.
 * Process the RRC transparent container and forward it to NAS entity.
 * \param assocId lower layer assoc id (SCTP)
 * \param stream SCTP stream on which data had been received
 * \param message The message as decoded by the ASN.1 codec
 * @returns -1 on failure, 0 otherwise
 **/
int s1ap_mme_handle_uplink_nas_transport(
    s1ap_state_t* state, const sctp_assoc_id_t assocId,
    const sctp_stream_id_t stream, S1ap_S1AP_PDU_t* message);

/** \brief Handle a NAS non delivery indication message from eNB
 * \param assocId lower layer assoc id (SCTP)
 * \param stream SCTP stream on which data had been received
 * \param message The message as decoded by the ASN.1 codec
 * @returns -1 on failure, 0 otherwise
 **/
int s1ap_mme_handle_nas_non_delivery(
    s1ap_state_t* state, const sctp_assoc_id_t assocId,
    const sctp_stream_id_t stream, S1ap_S1AP_PDU_t* message);

void s1ap_handle_conn_est_cnf(
    s1ap_state_t* state,
    const itti_mme_app_connection_establishment_cnf_t* const conn_est_cnf_p);

int s1ap_generate_downlink_nas_transport(
    s1ap_state_t* state, const enb_ue_s1ap_id_t enb_ue_s1ap_id,
    const mme_ue_s1ap_id_t ue_id, STOLEN_REF bstring* payload, imsi64_t imsi64,
    bool* is_state_same);

void s1ap_handle_mme_ue_id_notification(
    s1ap_state_t* state,
    const itti_mme_app_s1ap_mme_ue_id_notification_t* const notification_p);

int s1ap_generate_s1ap_e_rab_setup_req(
    s1ap_state_t* state, itti_s1ap_e_rab_setup_req_t* const e_rab_setup_req);

int s1ap_generate_s1ap_e_rab_rel_cmd(
    s1ap_state_t* state, itti_s1ap_e_rab_rel_cmd_t* const e_rab_rel_cmd);

#endif /* FILE_S1AP_MME_NAS_PROCEDURES_SEEN */

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

/*! \file s1ap_mme_itti_messaging.h
  \brief
  \author Sebastien ROUX, Lionel Gauthier
  \company Eurecom
  \email: lionel.gauthier@eurecom.fr
*/

#ifndef FILE_S1AP_MME_ITTI_MESSAGING_SEEN
#define FILE_S1AP_MME_ITTI_MESSAGING_SEEN

#include <stdbool.h>
#include <stddef.h>
#include <stdint.h>
#include <czmq.h>

#include "common_defs.h"
#include "3gpp_23.003.h"
#include "3gpp_36.401.h"
#include "S1ap_Cause.h"
#include "TrackingAreaIdentity.h"
#include "bstrlib.h"
#include "common_types.h"
#include "intertask_interface.h"

#include "s1ap_state.h"

extern task_zmq_ctx_t s1ap_task_zmq_ctx;
extern long s1ap_last_msg_latency;

int s1ap_mme_itti_send_sctp_request(
    STOLEN_REF bstring* payload, const uint32_t sctp_assoc_id_t,
    const sctp_stream_id_t stream, const mme_ue_s1ap_id_t ue_id);

int s1ap_mme_itti_nas_uplink_ind(
    const mme_ue_s1ap_id_t ue_id, STOLEN_REF bstring* payload,
    const tai_t* const tai, const ecgi_t* const cgi);

int s1ap_mme_itti_nas_downlink_cnf(
    const mme_ue_s1ap_id_t ue_id, const bool is_success);

void s1ap_mme_itti_s1ap_initial_ue_message(
    const sctp_assoc_id_t assoc_id, const uint32_t enb_id,
    const enb_ue_s1ap_id_t enb_ue_s1ap_id, const uint8_t* const nas_msg,
    const size_t nas_msg_length, const tai_t* const tai,
    const ecgi_t* const ecgi, const long rrc_cause,
    const s_tmsi_t* const opt_s_tmsi, const csg_id_t* const opt_csg_id,
    const gummei_t* const opt_gummei,
    const void* const opt_cell_access_mode,           // unused
    const void* const opt_cell_gw_transport_address,  // unused
    const void* const opt_relay_node_indicator);      // unused

void s1ap_mme_itti_nas_non_delivery_ind(
    const mme_ue_s1ap_id_t ue_id, uint8_t* const nas_msg,
    const size_t nas_msg_length, const S1ap_Cause_t* const cause,
    imsi64_t imsi64);

int s1ap_mme_itti_s1ap_path_switch_request(
    const sctp_assoc_id_t assoc_id, const uint32_t enb_id,
    const enb_ue_s1ap_id_t enb_ue_s1ap_id,
    const e_rab_to_be_switched_in_downlink_list_t* const
        e_rab_to_be_switched_dl_list,
    const mme_ue_s1ap_id_t mme_ue_s1ap_id, const ecgi_t* const ecgi,
    const tai_t* const tai, const uint16_t encryption_algorithm_capabilitie,
    uint16_t integrity_algorithm_capabilities, imsi64_t imsi64);

int s1ap_mme_itti_s1ap_handover_required(
    const sctp_assoc_id_t assoc_id, uint32_t enb_id, const S1ap_Cause_t cause,
    const S1ap_HandoverType_t handover_type,
    const mme_ue_s1ap_id_t mme_ue_s1ap_id, const bstring src_tgt_container,
    imsi64_t imsi64);

int s1ap_mme_itti_s1ap_handover_request_ack(
    const mme_ue_s1ap_id_t mme_ue_s1ap_id,
    const enb_ue_s1ap_id_t src_enb_ue_s1ap_id,
    const enb_ue_s1ap_id_t tgt_enb_ue_s1ap_id,
    const S1ap_HandoverType_t handover_type,
    const sctp_assoc_id_t source_assoc_id, const bstring tgt_src_container,
    const uint32_t source_enb_id, const uint32_t target_enb_id,
    imsi64_t imsi64);

int s1ap_mme_itti_s1ap_handover_notify(
    const mme_ue_s1ap_id_t mme_ue_s1ap_id,
    const s1ap_handover_state_t handover_state,
    const enb_ue_s1ap_id_t target_ue_s1ap_id,
    const sctp_assoc_id_t target_sctp_assoc_id, const ecgi_t ecgi,
    imsi64_t imsi64);
#endif /* FILE_S1AP_MME_ITTI_MESSAGING_SEEN */

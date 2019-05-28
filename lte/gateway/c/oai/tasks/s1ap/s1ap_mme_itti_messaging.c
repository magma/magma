/*
 * Licensed to the OpenAirInterface (OAI) Software Alliance under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The OpenAirInterface Software Alliance licenses this file to You under
 * the Apache License, Version 2.0  (the "License"); you may not use this file
 * except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
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

/*! \file s1ap_mme_itti_messaging.c
  \brief
  \author Sebastien ROUX, Lionel Gauthier
  \company Eurecom
  \email: lionel.gauthier@eurecom.fr
*/

#include <stdio.h>
#include <stdbool.h>
#include <stdint.h>

#include "bstrlib.h"
#include "log.h"
#include "assertions.h"
#include "intertask_interface.h"
#include "s1ap_mme_itti_messaging.h"
#include "S1ap-CauseRadioNetwork.h"
#include "nas/as_message.h"
#include "intertask_interface_types.h"
#include "itti_types.h"
#include "nas_messages_types.h"
#include "s1ap_messages_types.h"
#include "sctp_messages_types.h"

//------------------------------------------------------------------------------
int s1ap_mme_itti_send_sctp_request(
  STOLEN_REF bstring *payload,
  const sctp_assoc_id_t assoc_id,
  const sctp_stream_id_t stream,
  const mme_ue_s1ap_id_t ue_id)
{
  MessageDef *message_p = NULL;

  message_p = itti_alloc_new_message(TASK_S1AP, SCTP_DATA_REQ);
  SCTP_DATA_REQ(message_p).payload = *payload;
  *payload = NULL;
  SCTP_DATA_REQ(message_p).assoc_id = assoc_id;
  SCTP_DATA_REQ(message_p).stream = stream;
  SCTP_DATA_REQ(message_p).mme_ue_s1ap_id = ue_id;
  return itti_send_msg_to_task(TASK_SCTP, message_p);
}

//------------------------------------------------------------------------------
int s1ap_mme_itti_nas_uplink_ind(
  const mme_ue_s1ap_id_t ue_id,
  STOLEN_REF bstring *payload,
  const tai_t const *tai,
  const ecgi_t const *cgi)
{
  MessageDef *message_p = NULL;

  OAILOG_INFO(
    LOG_S1AP,
    "Sending NAS Uplink indication to MME_APP, mme_ue_s1ap_id = (%u) \n",
    ue_id);
  message_p = itti_alloc_new_message(TASK_S1AP, NAS_UPLINK_DATA_IND);
  NAS_UL_DATA_IND(message_p).ue_id = ue_id;
  NAS_UL_DATA_IND(message_p).nas_msg = *payload;
  *payload = NULL;
  NAS_UL_DATA_IND(message_p).tai = *tai;
  NAS_UL_DATA_IND(message_p).cgi = *cgi;

  return itti_send_msg_to_task(TASK_NAS_MME, message_p);
}

//------------------------------------------------------------------------------
int s1ap_mme_itti_nas_downlink_cnf(
  const mme_ue_s1ap_id_t ue_id,
  const bool is_success)
{
  MessageDef *message_p = NULL;

  if (ue_id == INVALID_MME_UE_S1AP_ID) {
    if (!is_success) {
      OAILOG_ERROR(
        LOG_S1AP,
        "ERROR: Failed to send connection less S1AP message to eNB. "
        "mme_ue_s1ap_id =  %d \n",
        ue_id);
    }
    // Drop this cnf message here since this is related to connection less S1AP message hence no need to send it to NAS module
    OAILOG_FUNC_RETURN(LOG_S1AP, RETURNok);
  }
  message_p = itti_alloc_new_message(TASK_S1AP, NAS_DOWNLINK_DATA_CNF);
  NAS_DL_DATA_CNF(message_p).ue_id = ue_id;
  if (is_success) {
    NAS_DL_DATA_CNF(message_p).err_code = AS_SUCCESS;
  } else {
    NAS_DL_DATA_CNF(message_p).err_code = AS_FAILURE;
    OAILOG_ERROR(
      LOG_S1AP,
      "ERROR: Failed to send S1AP message to eNB. mme_ue_s1ap_id =  %d \n",
      ue_id);
  }
  return itti_send_msg_to_task(TASK_NAS_MME, message_p);
}

//------------------------------------------------------------------------------

void s1ap_mme_itti_s1ap_initial_ue_message(
  const sctp_assoc_id_t assoc_id,
  const uint32_t enb_id,
  const enb_ue_s1ap_id_t enb_ue_s1ap_id,
  const uint8_t *const nas_msg,
  const size_t nas_msg_length,
  const tai_t const *tai,
  const ecgi_t const *ecgi,
  const long rrc_cause,
  const s_tmsi_t const *opt_s_tmsi,
  const csg_id_t const *opt_csg_id,
  const gummei_t const *opt_gummei,
  const void const *opt_cell_access_mode,          // unused
  const void const *opt_cell_gw_transport_address, // unused
  const void const *opt_relay_node_indicator)      // unused
{
  MessageDef *message_p = NULL;

  OAILOG_FUNC_IN(LOG_S1AP);
  AssertFatal(
    (nas_msg_length < 1000), "Bad length for NAS message %lu", nas_msg_length);
  message_p = itti_alloc_new_message(TASK_S1AP, S1AP_INITIAL_UE_MESSAGE);

  OAILOG_INFO(
    LOG_S1AP,
    "Sending Initial UE Message to MME_APP, enb_ue_s1ap_id "
    ": " ENB_UE_S1AP_ID_FMT "\n",
    enb_ue_s1ap_id);

  S1AP_INITIAL_UE_MESSAGE(message_p).sctp_assoc_id = assoc_id;
  S1AP_INITIAL_UE_MESSAGE(message_p).enb_ue_s1ap_id = enb_ue_s1ap_id;
  S1AP_INITIAL_UE_MESSAGE(message_p).enb_id = enb_id;

  S1AP_INITIAL_UE_MESSAGE(message_p).nas = blk2bstr(nas_msg, nas_msg_length);

  S1AP_INITIAL_UE_MESSAGE(message_p).tai = *tai;
  S1AP_INITIAL_UE_MESSAGE(message_p).ecgi = *ecgi;
  S1AP_INITIAL_UE_MESSAGE(message_p).rrc_establishment_cause = rrc_cause + 1;

  if (opt_s_tmsi) {
    S1AP_INITIAL_UE_MESSAGE(message_p).is_s_tmsi_valid = true;
    S1AP_INITIAL_UE_MESSAGE(message_p).opt_s_tmsi = *opt_s_tmsi;
  } else {
    S1AP_INITIAL_UE_MESSAGE(message_p).is_s_tmsi_valid = false;
  }
  if (opt_csg_id) {
    S1AP_INITIAL_UE_MESSAGE(message_p).is_csg_id_valid = true;
    S1AP_INITIAL_UE_MESSAGE(message_p).opt_csg_id = *opt_csg_id;
  } else {
    S1AP_INITIAL_UE_MESSAGE(message_p).is_csg_id_valid = false;
  }
  if (opt_gummei) {
    S1AP_INITIAL_UE_MESSAGE(message_p).is_gummei_valid = true;
    S1AP_INITIAL_UE_MESSAGE(message_p).opt_gummei = *opt_gummei;
  } else {
    S1AP_INITIAL_UE_MESSAGE(message_p).is_gummei_valid = false;
  }

  S1AP_INITIAL_UE_MESSAGE(message_p).transparent.enb_ue_s1ap_id =
    enb_ue_s1ap_id;
  S1AP_INITIAL_UE_MESSAGE(message_p).transparent.e_utran_cgi = *ecgi;

  itti_send_msg_to_task(TASK_MME_APP, message_p);
  OAILOG_FUNC_OUT(LOG_S1AP);
}

//------------------------------------------------------------------------------
static int s1ap_mme_non_delivery_cause_2_nas_data_rej_cause(
  const S1ap_Cause_t *const cause)
{
  switch (cause->present) {
    case S1ap_Cause_PR_radioNetwork:
      switch (cause->choice.radioNetwork) {
        case S1ap_CauseRadioNetwork_handover_cancelled:
        case S1ap_CauseRadioNetwork_partial_handover:
        case S1ap_CauseRadioNetwork_successful_handover:
        case S1ap_CauseRadioNetwork_ho_failure_in_target_EPC_eNB_or_target_system:
        case S1ap_CauseRadioNetwork_ho_target_not_allowed:
        case S1ap_CauseRadioNetwork_handover_desirable_for_radio_reason: /// ?
        case S1ap_CauseRadioNetwork_time_critical_handover:
        case S1ap_CauseRadioNetwork_resource_optimisation_handover:
        case S1ap_CauseRadioNetwork_s1_intra_system_handover_triggered:
        case S1ap_CauseRadioNetwork_s1_inter_system_handover_triggered:
        case S1ap_CauseRadioNetwork_x2_handover_triggered:
          return AS_NON_DELIVERED_DUE_HO;
          break;

        default: return AS_FAILURE;
      }
      break;

    default: return AS_FAILURE;
  }
  return AS_FAILURE;
}
//------------------------------------------------------------------------------
void s1ap_mme_itti_nas_non_delivery_ind(
  const mme_ue_s1ap_id_t ue_id,
  uint8_t *const nas_msg,
  const size_t nas_msg_length,
  const S1ap_Cause_t *const cause)
{
  MessageDef *message_p = NULL;
  // TODO translate, insert, cause in message
  OAILOG_FUNC_IN(LOG_S1AP);
  message_p = itti_alloc_new_message(TASK_S1AP, NAS_DOWNLINK_DATA_REJ);

  NAS_DL_DATA_REJ(message_p).ue_id = ue_id;
  /* Mapping between asn1 definition and NAS definition */
  NAS_DL_DATA_REJ(message_p).err_code =
    s1ap_mme_non_delivery_cause_2_nas_data_rej_cause(cause);
  NAS_DL_DATA_REJ(message_p).nas_msg = blk2bstr(nas_msg, nas_msg_length);

  // should be sent to MME_APP, but this one would forward it to NAS_MME, so send it directly to NAS_MME
  // but let's see
  itti_send_msg_to_task(TASK_NAS_MME, message_p);
  OAILOG_FUNC_OUT(LOG_S1AP);
}

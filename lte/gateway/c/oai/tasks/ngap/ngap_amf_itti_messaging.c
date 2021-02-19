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

/*! \file ngap_amf_itti_messaging.c
 */

#include <stdio.h>
#include <stdbool.h>
#include <stdint.h>
#include "common_defs.h"
#include "bstrlib.h"
#include "log.h"
#include "assertions.h"
#include "intertask_interface.h"
#include "Ngap_CauseRadioNetwork.h"
#include "nas/as_message.h"
#include "intertask_interface_types.h"
#include "itti_types.h"
#include "ngap_types.h"
#include "sctp_messages_types.h"
#include "amf_app_messages_types.h"
#include "amf_default_values.h"
task_zmq_ctx_t ngap_task_zmq_ctx;
//------------------------------------------------------------------------------
int ngap_amf_itti_send_sctp_request(
    STOLEN_REF bstring* payload, const sctp_assoc_id_t assoc_id,
    const sctp_stream_id_t stream, const amf_ue_ngap_id_t ue_id) {
  MessageDef* message_p = NULL;

  message_p = itti_alloc_new_message(TASK_NGAP, SCTP_DATA_REQ);
  if (message_p == NULL) {
    OAILOG_ERROR(
        LOG_NGAP,
        "itti_alloc_new_message Failed for"
        " SCTP_DATA_REQ \n");
    OAILOG_FUNC_RETURN(LOG_NGAP, RETURNerror);
  }
  SCTP_DATA_REQ(message_p).payload        = *payload;
  *payload                                = NULL;
  SCTP_DATA_REQ(message_p).assoc_id       = assoc_id;
  SCTP_DATA_REQ(message_p).stream         = stream;
  SCTP_DATA_REQ(message_p).amf_ue_ngap_id = ue_id;
  SCTP_DATA_REQ(message_p).ppid           = NGAP_SCTP_PPID;
  OAILOG_ERROR(LOG_NGAP, "######ACL_TAG: %s, %d ", __func__, __LINE__);
  return send_msg_to_task(&ngap_task_zmq_ctx, TASK_SCTP, message_p);
}
//------------------------------------------------------------------------------

typedef uint32_t amf_ue_ngap_id_t;

int ngap_amf_itti_nas_uplink_ind(
    // const amf_ue_ngap_id_t ue_id,
    const amf_ue_ngap_id_t ue_id, STOLEN_REF bstring* payload,
    const tai_t const* tai, const ecgi_t const* cgi) {
  MessageDef* message_p = NULL;
  imsi64_t imsi64       = INVALID_IMSI64;
#if 0
TODO: laterthis maps to state
  ngap_imsi_map_t* imsi_map = get_ngap_imsi_map();

  hashtable_uint64_ts_get(
    imsi_map->amf_ue_id_imsi_htbl, (const hash_key_t) ue_id, &imsi64);
#endif
  OAILOG_INFO_UE(
      LOG_NGAP, imsi64,
      "Sending NAS Uplink indication to NAS_AMF_APP, amf_ue_ngap_id = (%u) \n",
      ue_id);
  message_p = itti_alloc_new_message(TASK_NGAP, AMF_APP_UPLINK_DATA_IND);
  if (message_p == NULL) {
    OAILOG_ERROR_UE(
        LOG_NGAP, imsi64,
        "itti_alloc_new_message Failed for"
        " AMF_APP_UPLINK_DATA_IND \n");
    OAILOG_FUNC_RETURN(LOG_NGAP, RETURNerror);
  }
  AMF_APP_UL_DATA_IND(message_p).ue_id   = ue_id;
  AMF_APP_UL_DATA_IND(message_p).nas_msg = *payload;
  *payload                               = NULL;
  AMF_APP_UL_DATA_IND(message_p).tai     = *tai;
  AMF_APP_UL_DATA_IND(message_p).cgi     = *cgi;

  message_p->ittiMsgHeader.imsi = imsi64;
  OAILOG_ERROR(LOG_NGAP, "%s, ########send to AMF :%d", __func__, __LINE__);
  return send_msg_to_task(&ngap_task_zmq_ctx, TASK_AMF_APP, message_p);
}

//------------------------------------------------------------------------------
static int ngap_amf_non_delivery_cause_2_nas_data_rej_cause(
    const Ngap_Cause_t* const cause) {
  switch (cause->present) {
    case Ngap_Cause_PR_radioNetwork:
      switch (cause->choice.radioNetwork) {
        case Ngap_CauseRadioNetwork_handover_cancelled:
        case Ngap_CauseRadioNetwork_partial_handover:
        case Ngap_CauseRadioNetwork_successful_handover:
        case Ngap_CauseRadioNetwork_ho_failure_in_target_5GC_ngran_node_or_target_system:
        case Ngap_CauseRadioNetwork_ho_target_not_allowed:
        case Ngap_CauseRadioNetwork_handover_desirable_for_radio_reason:  /// ?
        case Ngap_CauseRadioNetwork_time_critical_handover:
        case Ngap_CauseRadioNetwork_resource_optimisation_handover:
        case Ngap_CauseRadioNetwork_ng_intra_system_handover_triggered:
        case Ngap_CauseRadioNetwork_ng_inter_system_handover_triggered:
        case Ngap_CauseRadioNetwork_xn_handover_triggered:
          return AS_NON_DELIVERED_DUE_HO;
          break;

        default:
          return AS_FAILURE;
      }
      break;

    default:
      return AS_FAILURE;
  }
  return AS_FAILURE;
}

//------------------------------------------------------------------------------

void ngap_amf_itti_ngap_initial_ue_message(
    const sctp_assoc_id_t assoc_id, const uint32_t gnb_id,
    const gnb_ue_ngap_id_t gnb_ue_ngap_id, const uint8_t* const nas_msg,
    const size_t nas_msg_length, const tai_t const* tai,
    const ecgi_t const* ecgi, const long rrc_cause,
    const s_tmsi_m5_t const* opt_s_tmsi, const csg_id_t const* opt_csg_id,
    const guamfi_t const* opt_guamfi,
    const void const* opt_cell_access_mode,           // unused
    const void const* opt_cell_gw_transport_address,  // unused
    const void const* opt_relay_node_indicator)       // unused
{
  MessageDef* message_p = NULL;
  // InitialUEMessage_IEs_t *initialUEMessage_IEs = NULL;

  OAILOG_FUNC_IN(LOG_NGAP);
  /*
  AssertFatal(
    (nas_msg_length < 1000), "Bad length for NAS message %lu", nas_msg_length);
  */
  message_p = itti_alloc_new_message(TASK_NGAP, NGAP_INITIAL_UE_MESSAGE);

  if (message_p == NULL) {
    OAILOG_ERROR(
        LOG_NGAP,
        "itti_alloc_new_message Failed for"
        " NGAP_INITIAL_UE_MESSAGE \n");
    OAILOG_FUNC_OUT(LOG_NGAP);
  }

  OAILOG_INFO(
      LOG_NGAP,
      "Sending Initial UE Message to AMF_APP: ID: %d, NGAP_INITIAL_UE_MESSAGE: "
      "%d \n",
      ITTI_MSG_ID(message_p), NGAP_INITIAL_UE_MESSAGE);

  OAILOG_INFO(
      LOG_NGAP,
      "Sending Initial UE Message to AMF_APP \n");  // Need change

  NGAP_INITIAL_UE_MESSAGE(message_p).sctp_assoc_id  = assoc_id;
  NGAP_INITIAL_UE_MESSAGE(message_p).gnb_ue_ngap_id = gnb_ue_ngap_id;
  NGAP_INITIAL_UE_MESSAGE(message_p).gnb_id         = gnb_id;
  // NGAP_INITIAL_UE_MESSAGE(message_p).ran_ue_ngap_id =
  // initialUEMessage_IEs->ran_ue_ngap_id;
  // NGAP_INITIAL_UE_MESSAGE(message_p).gnb_ue_ngap_id =
  // initialUEMessage_IEs->ran_ue_ngap_id;
  NGAP_INITIAL_UE_MESSAGE(message_p).nas = blk2bstr(nas_msg, nas_msg_length);
  // NGAP_INITIAL_UE_MESSAGE(message_p).nas = initialUEMessage_IEs->nas_pdu;
  // NGAP_INITIAL_UE_MESSAGE(message_p).userLocationInformation =
  // initialUEMessage_IEs->userLocationInformation;
  // NGAP_INITIAL_UE_MESSAGE(message_p).rrcEstablishmentCause =
  // initialUEMessage_IEs->rrcEstablishmentCause;
  // NGAP_INITIAL_UE_MESSAGE(message_p).ueContextRequest =
  // initialUEMessage_IEs->ueContextRequest;

#if 0  
  NGAP_INITIAL_UE_MESSAGE(message_p).tai = *tai;
  NGAP_INITIAL_UE_MESSAGE(message_p).ecgi = *ecgi;
#endif
  // NGAP_INITIAL_UE_MESSAGE(message_p).rrc_establishment_cause = rrc_cause + 1;
  NGAP_INITIAL_UE_MESSAGE(message_p).m5g_rrc_establishment_cause =
      rrc_cause + 1;

  if (opt_s_tmsi) {
    NGAP_INITIAL_UE_MESSAGE(message_p).is_s_tmsi_valid = true;
    NGAP_INITIAL_UE_MESSAGE(message_p).opt_s_tmsi      = *opt_s_tmsi;
  } else {
    NGAP_INITIAL_UE_MESSAGE(message_p).is_s_tmsi_valid = false;
  }
#if 0  
  if (opt_csg_id) {
    NGAP_INITIAL_UE_MESSAGE(message_p).is_csg_id_valid = true;
    NGAP_INITIAL_UE_MESSAGE(message_p).opt_csg_id = *opt_csg_id;
  } else {
    NGAP_INITIAL_UE_MESSAGE(message_p).is_csg_id_valid = false;
  }
  if (opt_guamfi) {
    NGAP_INITIAL_UE_MESSAGE(message_p).is_guamfi_valid = true;
    NGAP_INITIAL_UE_MESSAGE(message_p).opt_guamfi = *opt_guamfi;
  } else {
    NGAP_INITIAL_UE_MESSAGE(message_p).is_guamfi_valid = false;
  }
#endif
  NGAP_INITIAL_UE_MESSAGE(message_p).gnb_ue_ngap_id = gnb_ue_ngap_id;
  // NGAP_INITIAL_UE_MESSAGE(message_p).transparent.e_utran_cgi = *ecgi;

  send_msg_to_task(&ngap_task_zmq_ctx, TASK_AMF_APP, message_p);
  OAILOG_ERROR(LOG_NGAP, "####ACL_TAG iniUEmsg sent to TASK_AMF_APP");
  OAILOG_FUNC_OUT(LOG_NGAP);
}

//------------------------------------------------------------------------------
void ngap_amf_itti_nas_non_delivery_ind(
    const amf_ue_ngap_id_t ue_id, uint8_t* const nas_msg,
    const size_t nas_msg_length, const Ngap_Cause_t* const cause,
    const imsi64_t imsi64) {
  MessageDef* message_p = NULL;
  // TODO translate, insert, cause in message
  OAILOG_FUNC_IN(LOG_NGAP);
  message_p = itti_alloc_new_message(TASK_NGAP, AMF_APP_DOWNLINK_DATA_REJ);
  if (message_p == NULL) {
    OAILOG_ERROR_UE(
        LOG_NGAP, imsi64,
        "itti_alloc_new_message Failed for"
        " AMF_APP_DOWNLINK_DATA_REJ \n");
    OAILOG_FUNC_OUT(LOG_NGAP);
  }

  AMF_APP_DL_DATA_REJ(message_p).ue_id = ue_id;
  /* Mapping between asn1 definition and NAS definition */
  AMF_APP_DL_DATA_REJ(message_p).err_code =
      ngap_amf_non_delivery_cause_2_nas_data_rej_cause(cause);
  AMF_APP_DL_DATA_REJ(message_p).nas_msg = blk2bstr(nas_msg, nas_msg_length);
  // should be sent to AMF_APP, but this one would forward it to NAS_AMF, so
  // send it directly to NAS_AMF but let's see

  message_p->ittiMsgHeader.imsi = imsi64;
  send_msg_to_task(&ngap_task_zmq_ctx, TASK_AMF_APP, message_p);
  OAILOG_FUNC_OUT(LOG_NGAP);
}

//------------------------------------------------------------------------------
#if 0
Info: Later to do
int ngap_amf_itti_ngap_path_switch_request(
    const sctp_assoc_id_t assoc_id, const uint32_t enb_id,
    const enb_ue_ngap_id_t enb_ue_ngap_id,
    const e_rab_to_be_switched_in_downlink_list_t const*
        e_rab_to_be_switched_dl_list,
    const amf_ue_ngap_id_t amf_ue_ngap_id, const ecgi_t const* ecgi,
    const tai_t const* tai, const uint16_t encryption_algorithm_capabilities,
    const uint16_t integrity_algorithm_capabilities, const imsi64_t imsi64) {
  MessageDef* message_p = NULL;
  message_p = itti_alloc_new_message(TASK_NGAP, NGAP_PATH_SWITCH_REQUEST);
  if (message_p == NULL) {
    OAILOG_ERROR_UE(LOG_NGAP, imsi64, "itti_alloc_new_message Failed");
    OAILOG_FUNC_RETURN(LOG_NGAP, RETURNerror);
  }
  NGAP_PATH_SWITCH_REQUEST(message_p).sctp_assoc_id  = assoc_id;
  NGAP_PATH_SWITCH_REQUEST(message_p).enb_id         = enb_id;
  NGAP_PATH_SWITCH_REQUEST(message_p).enb_ue_ngap_id = enb_ue_ngap_id;
  NGAP_PATH_SWITCH_REQUEST(message_p).e_rab_to_be_switched_dl_list =
      *e_rab_to_be_switched_dl_list;
  NGAP_PATH_SWITCH_REQUEST(message_p).amf_ue_ngap_id = amf_ue_ngap_id;
  NGAP_PATH_SWITCH_REQUEST(message_p).tai            = *tai;
  NGAP_PATH_SWITCH_REQUEST(message_p).ecgi           = *ecgi;
  NGAP_PATH_SWITCH_REQUEST(message_p).encryption_algorithm_capabilities =
      encryption_algorithm_capabilities;
  NGAP_PATH_SWITCH_REQUEST(message_p).integrity_algorithm_capabilities =
      integrity_algorithm_capabilities;

  OAILOG_DEBUG_UE(
      LOG_NGAP, imsi64,
      "sending Path Switch Request to AMF_APP for source amf_ue_ngap_id %d\n",
      amf_ue_ngap_id);

  message_p->ittiMsgHeader.imsi = imsi64;
  send_msg_to_task(&ngap_task_zmq_ctx, TASK_AMF_APP, message_p);
  OAILOG_FUNC_RETURN(LOG_NGAP, RETURNok);
}
#endif

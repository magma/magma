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

/*! \file s1ap_mme_itti_messaging.c
  \brief
  \author Sebastien ROUX, Lionel Gauthier
  \company Eurecom
  \email: lionel.gauthier@eurecom.fr
*/

#include "lte/gateway/c/core/oai/tasks/s1ap/s1ap_mme_itti_messaging.hpp"

#ifdef __cplusplus
extern "C" {
#endif
#include "lte/gateway/c/core/common/assertions.h"
#include "lte/gateway/c/core/oai/common/log.h"
#include "lte/gateway/c/core/oai/lib/bstr/bstrlib.h"
#include "lte/gateway/c/core/oai/lib/itti/intertask_interface.h"
#include "lte/gateway/c/core/oai/lib/itti/intertask_interface_types.h"
#ifdef __cplusplus
}
#endif

#include "S1ap_CauseRadioNetwork.h"
#include "lte/gateway/c/core/oai/include/mme_app_ue_context.hpp"
#include "lte/gateway/c/core/oai/include/nas/as_message.h"
#include "lte/gateway/c/core/oai/include/s1ap_messages_types.h"
#include "lte/gateway/c/core/oai/include/sctp_messages_types.hpp"

namespace magma {
namespace lte {

//------------------------------------------------------------------------------
status_code_e s1ap_mme_itti_send_sctp_request(STOLEN_REF bstring* payload,
                                              const sctp_assoc_id_t assoc_id,
                                              const sctp_stream_id_t stream,
                                              const mme_ue_s1ap_id_t ue_id) {
  MessageDef* message_p = NULL;

  message_p = itti_alloc_new_message(TASK_S1AP, SCTP_DATA_REQ);
  if (message_p == NULL) {
    OAILOG_ERROR(LOG_S1AP,
                 "itti_alloc_new_message Failed for"
                 " SCTP_DATA_REQ \n");
    OAILOG_FUNC_RETURN(LOG_S1AP, RETURNerror);
  }
  SCTP_DATA_REQ(message_p).payload = *payload;
  *payload = NULL;
  SCTP_DATA_REQ(message_p).assoc_id = assoc_id;
  SCTP_DATA_REQ(message_p).stream = stream;
  SCTP_DATA_REQ(message_p).agw_ue_xap_id = ue_id;
  SCTP_DATA_REQ(message_p).ppid = S1AP_SCTP_PPID;
  return send_msg_to_task(&s1ap_task_zmq_ctx, TASK_SCTP, message_p);
}

//------------------------------------------------------------------------------
status_code_e s1ap_mme_itti_nas_uplink_ind(const mme_ue_s1ap_id_t ue_id,
                                           STOLEN_REF bstring* payload,
                                           const tai_t* const tai,
                                           const ecgi_t* const cgi) {
  MessageDef* message_p = NULL;
  imsi64_t imsi64 = INVALID_IMSI64;

  magma::proto_map_uint32_uint64_t ueid_imsi_map;
  get_s1ap_ueid_imsi_map(&ueid_imsi_map);
  ueid_imsi_map.get(ue_id, &imsi64);

  OAILOG_INFO_UE(
      LOG_S1AP, imsi64,
      "Sending NAS Uplink indication to NAS_MME_APP, mme_ue_s1ap_id = (%u) \n",
      ue_id);
  message_p = itti_alloc_new_message(TASK_S1AP, MME_APP_UPLINK_DATA_IND);
  if (message_p == NULL) {
    OAILOG_ERROR_UE(LOG_S1AP, imsi64,
                    "itti_alloc_new_message Failed for"
                    " MME_APP_UPLINK_DATA_IND \n");
    OAILOG_FUNC_RETURN(LOG_S1AP, RETURNerror);
  }
  ITTI_MSG_LASTHOP_LATENCY(message_p) = s1ap_last_msg_latency;
  MME_APP_UL_DATA_IND(message_p).ue_id = ue_id;
  MME_APP_UL_DATA_IND(message_p).nas_msg = *payload;
  *payload = NULL;
  MME_APP_UL_DATA_IND(message_p).tai = *tai;
  MME_APP_UL_DATA_IND(message_p).cgi = *cgi;

  message_p->ittiMsgHeader.imsi = imsi64;
  return send_msg_to_task(&s1ap_task_zmq_ctx, TASK_MME_APP, message_p);
}

//------------------------------------------------------------------------------
status_code_e s1ap_mme_itti_nas_downlink_cnf(const mme_ue_s1ap_id_t ue_id,
                                             const bool is_success) {
  MessageDef* message_p = NULL;
  imsi64_t imsi64 = INVALID_IMSI64;

  if (ue_id == INVALID_MME_UE_S1AP_ID) {
    if (!is_success) {
      OAILOG_ERROR(LOG_S1AP,
                   "ERROR: Failed to send connectionless S1AP message to eNB. "
                   "mme_ue_s1ap_id =  %d \n",
                   ue_id);
    }
    // Drop this cnf message here since this is related to connectionless S1AP
    // message hence no need to send it to NAS module
    OAILOG_FUNC_RETURN(LOG_S1AP, RETURNok);
  }

  magma::proto_map_uint32_uint64_t ueid_imsi_map;
  get_s1ap_ueid_imsi_map(&ueid_imsi_map);
  ueid_imsi_map.get(ue_id, &imsi64);

  message_p = itti_alloc_new_message(TASK_S1AP, MME_APP_DOWNLINK_DATA_CNF);
  if (message_p == NULL) {
    OAILOG_ERROR_UE(LOG_S1AP, imsi64,
                    "itti_alloc_new_message Failed for"
                    " MME_APP_DOWNLINK_DATA_CNF \n");
    OAILOG_FUNC_RETURN(LOG_S1AP, RETURNerror);
  }
  MME_APP_DL_DATA_CNF(message_p).ue_id = ue_id;
  if (is_success) {
    MME_APP_DL_DATA_CNF(message_p).err_code = AS_SUCCESS;
  } else {
    MME_APP_DL_DATA_CNF(message_p).err_code = AS_FAILURE;
    OAILOG_ERROR_UE(
        LOG_S1AP, imsi64,
        "ERROR: Failed to send S1AP message to eNB. mme_ue_s1ap_id =  %d \n",
        ue_id);
  }
  message_p->ittiMsgHeader.imsi = imsi64;
  return send_msg_to_task(&s1ap_task_zmq_ctx, TASK_MME_APP, message_p);
}

//------------------------------------------------------------------------------

status_code_e s1ap_mme_itti_s1ap_initial_ue_message(
    const sctp_assoc_id_t assoc_id, const uint32_t enb_id,
    const enb_ue_s1ap_id_t enb_ue_s1ap_id, const uint8_t* const nas_msg,
    const size_t nas_msg_length, const tai_t* const tai,
    const ecgi_t* const ecgi, const long rrc_cause,
    const s_tmsi_t* const opt_s_tmsi, const csg_id_t* const opt_csg_id,
    const gummei_t* const opt_gummei,
    const void* const opt_cell_access_mode,           // unused
    const void* const opt_cell_gw_transport_address,  // unused
    const void* const opt_relay_node_indicator)       // unused
{
  MessageDef* message_p = NULL;

  OAILOG_FUNC_IN(LOG_S1AP);

  if (nas_msg_length >= 1000) {
    OAILOG_ERROR(LOG_S1AP,
                 "Bad length for NAS message greater then 1000 "
                 " S1AP_INITIAL_UE_MESSAGE \n");
    OAILOG_FUNC_RETURN(LOG_S1AP, RETURNerror);
  }

  message_p = itti_alloc_new_message(TASK_S1AP, S1AP_INITIAL_UE_MESSAGE);
  if (message_p == NULL) {
    OAILOG_ERROR(LOG_S1AP,
                 "itti_alloc_new_message Failed for"
                 " S1AP_INITIAL_UE_MESSAGE \n");
    OAILOG_FUNC_RETURN(LOG_S1AP, RETURNerror);
  }

  OAILOG_INFO(LOG_S1AP,
              "Sending Initial UE Message to MME_APP, enb_ue_s1ap_id "
              ": " ENB_UE_S1AP_ID_FMT "\n",
              enb_ue_s1ap_id);

  ITTI_MSG_LASTHOP_LATENCY(message_p) = s1ap_last_msg_latency;
  S1AP_INITIAL_UE_MESSAGE(message_p).sctp_assoc_id = assoc_id;
  S1AP_INITIAL_UE_MESSAGE(message_p).enb_ue_s1ap_id = enb_ue_s1ap_id;
  S1AP_INITIAL_UE_MESSAGE(message_p).enb_id = enb_id;

  S1AP_INITIAL_UE_MESSAGE(message_p).nas = blk2bstr(nas_msg, nas_msg_length);

  S1AP_INITIAL_UE_MESSAGE(message_p).tai = *tai;

  S1AP_INITIAL_UE_MESSAGE(message_p).ecgi = *ecgi;
  S1AP_INITIAL_UE_MESSAGE(message_p).rrc_establishment_cause =
      (rrc_establishment_cause_t)(rrc_cause + 1);

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

  send_msg_to_task(&s1ap_task_zmq_ctx, TASK_MME_APP, message_p);
  OAILOG_FUNC_RETURN(LOG_S1AP, RETURNok);
}

//------------------------------------------------------------------------------
static int s1ap_mme_non_delivery_cause_2_nas_data_rej_cause(
    const S1ap_Cause_t* const cause) {
  switch (cause->present) {
    case S1ap_Cause_PR_radioNetwork:
      switch (cause->choice.radioNetwork) {
        case S1ap_CauseRadioNetwork_handover_cancelled:
        case S1ap_CauseRadioNetwork_partial_handover:
        case S1ap_CauseRadioNetwork_successful_handover:
        case S1ap_CauseRadioNetwork_ho_failure_in_target_EPC_eNB_or_target_system:
        case S1ap_CauseRadioNetwork_ho_target_not_allowed:
        case S1ap_CauseRadioNetwork_handover_desirable_for_radio_reason:  /// ?
        case S1ap_CauseRadioNetwork_time_critical_handover:
        case S1ap_CauseRadioNetwork_resource_optimisation_handover:
        case S1ap_CauseRadioNetwork_s1_intra_system_handover_triggered:
        case S1ap_CauseRadioNetwork_s1_inter_system_handover_triggered:
        case S1ap_CauseRadioNetwork_x2_handover_triggered:
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
void s1ap_mme_itti_nas_non_delivery_ind(const mme_ue_s1ap_id_t ue_id,
                                        uint8_t* const nas_msg,
                                        const size_t nas_msg_length,
                                        const S1ap_Cause_t* const cause,
                                        imsi64_t imsi64) {
  MessageDef* message_p = NULL;
  // TODO translate, insert, cause in message
  OAILOG_FUNC_IN(LOG_S1AP);
  message_p = itti_alloc_new_message(TASK_S1AP, MME_APP_DOWNLINK_DATA_REJ);
  if (message_p == NULL) {
    OAILOG_ERROR_UE(LOG_S1AP, imsi64,
                    "itti_alloc_new_message Failed for"
                    " MME_APP_DOWNLINK_DATA_REJ \n");
    OAILOG_FUNC_OUT(LOG_S1AP);
  }

  MME_APP_DL_DATA_REJ(message_p).ue_id = ue_id;
  /* Mapping between asn1 definition and NAS definition */
  MME_APP_DL_DATA_REJ(message_p).err_code =
      s1ap_mme_non_delivery_cause_2_nas_data_rej_cause(cause);
  MME_APP_DL_DATA_REJ(message_p).nas_msg = blk2bstr(nas_msg, nas_msg_length);

  message_p->ittiMsgHeader.imsi = imsi64;
  send_msg_to_task(&s1ap_task_zmq_ctx, TASK_MME_APP, message_p);
  OAILOG_FUNC_OUT(LOG_S1AP);
}

//------------------------------------------------------------------------------
status_code_e s1ap_mme_itti_s1ap_path_switch_request(
    const sctp_assoc_id_t assoc_id, uint32_t enb_id,
    const enb_ue_s1ap_id_t enb_ue_s1ap_id,
    const e_rab_to_be_switched_in_downlink_list_t* const
        e_rab_to_be_switched_dl_list,
    const mme_ue_s1ap_id_t mme_ue_s1ap_id, const ecgi_t* const ecgi,
    const tai_t* const tai, uint16_t encryption_algorithm_capabilities,
    uint16_t integrity_algorithm_capabilities, imsi64_t imsi64) {
  MessageDef* message_p = NULL;
  message_p = itti_alloc_new_message(TASK_S1AP, S1AP_PATH_SWITCH_REQUEST);
  if (message_p == NULL) {
    OAILOG_ERROR_UE(LOG_S1AP, imsi64, "itti_alloc_new_message Failed");
    OAILOG_FUNC_RETURN(LOG_S1AP, RETURNerror);
  }
  S1AP_PATH_SWITCH_REQUEST(message_p).sctp_assoc_id = assoc_id;
  S1AP_PATH_SWITCH_REQUEST(message_p).enb_id = enb_id;
  S1AP_PATH_SWITCH_REQUEST(message_p).enb_ue_s1ap_id = enb_ue_s1ap_id;
  S1AP_PATH_SWITCH_REQUEST(message_p).e_rab_to_be_switched_dl_list =
      *e_rab_to_be_switched_dl_list;
  S1AP_PATH_SWITCH_REQUEST(message_p).mme_ue_s1ap_id = mme_ue_s1ap_id;
  S1AP_PATH_SWITCH_REQUEST(message_p).tai = *tai;
  S1AP_PATH_SWITCH_REQUEST(message_p).ecgi = *ecgi;
  S1AP_PATH_SWITCH_REQUEST(message_p).encryption_algorithm_capabilities =
      encryption_algorithm_capabilities;
  S1AP_PATH_SWITCH_REQUEST(message_p).integrity_algorithm_capabilities =
      integrity_algorithm_capabilities;

  OAILOG_DEBUG_UE(
      LOG_S1AP, imsi64,
      "sending Path Switch Request to MME_APP for source mme_ue_s1ap_id %d\n",
      mme_ue_s1ap_id);

  message_p->ittiMsgHeader.imsi = imsi64;
  send_msg_to_task(&s1ap_task_zmq_ctx, TASK_MME_APP, message_p);
  OAILOG_FUNC_RETURN(LOG_S1AP, RETURNok);
}

//------------------------------------------------------------------------------
status_code_e s1ap_mme_itti_s1ap_handover_required(
    const sctp_assoc_id_t assoc_id, uint32_t enb_id, const S1ap_Cause_t cause,
    const S1ap_HandoverType_t handover_type,
    const mme_ue_s1ap_id_t mme_ue_s1ap_id, const bstring src_tgt_container,
    imsi64_t imsi64) {
  MessageDef* message_p = NULL;
  message_p = itti_alloc_new_message(TASK_S1AP, S1AP_HANDOVER_REQUIRED);
  if (message_p == NULL) {
    OAILOG_ERROR_UE(LOG_S1AP, imsi64, "itti_alloc_new_message Failed");
    OAILOG_FUNC_RETURN(LOG_S1AP, RETURNerror);
  }
  S1AP_HANDOVER_REQUIRED(message_p).sctp_assoc_id = assoc_id;
  S1AP_HANDOVER_REQUIRED(message_p).enb_id = enb_id;
  S1AP_HANDOVER_REQUIRED(message_p).cause = cause;
  S1AP_HANDOVER_REQUIRED(message_p).handover_type = handover_type;
  S1AP_HANDOVER_REQUIRED(message_p).mme_ue_s1ap_id = mme_ue_s1ap_id;
  S1AP_HANDOVER_REQUIRED(message_p).src_tgt_container = src_tgt_container;

  OAILOG_DEBUG_UE(
      LOG_S1AP, imsi64,
      "sending Handover Required to MME_APP for mme_ue_s1ap_id %d\n",
      mme_ue_s1ap_id);

  message_p->ittiMsgHeader.imsi = imsi64;
  send_msg_to_task(&s1ap_task_zmq_ctx, TASK_MME_APP, message_p);
  OAILOG_FUNC_RETURN(LOG_S1AP, RETURNok);
}

status_code_e s1ap_mme_itti_s1ap_handover_request_ack(
    const mme_ue_s1ap_id_t mme_ue_s1ap_id,
    const enb_ue_s1ap_id_t src_enb_ue_s1ap_id,
    const enb_ue_s1ap_id_t tgt_enb_ue_s1ap_id,
    const S1ap_HandoverType_t handover_type,
    const sctp_assoc_id_t source_assoc_id, const bstring tgt_src_container,
    const uint32_t source_enb_id, const uint32_t target_enb_id,
    imsi64_t imsi64) {
  MessageDef* message_p = NULL;
  message_p = itti_alloc_new_message(TASK_S1AP, S1AP_HANDOVER_REQUEST_ACK);
  if (message_p == NULL) {
    OAILOG_ERROR_UE(LOG_S1AP, imsi64, "itti_alloc_new_message Failed");
    OAILOG_FUNC_RETURN(LOG_S1AP, RETURNerror);
  }
  S1AP_HANDOVER_REQUEST_ACK(message_p).mme_ue_s1ap_id = mme_ue_s1ap_id;
  S1AP_HANDOVER_REQUEST_ACK(message_p).src_enb_ue_s1ap_id = src_enb_ue_s1ap_id;
  S1AP_HANDOVER_REQUEST_ACK(message_p).tgt_enb_ue_s1ap_id = tgt_enb_ue_s1ap_id;
  S1AP_HANDOVER_REQUEST_ACK(message_p).source_assoc_id = source_assoc_id;
  S1AP_HANDOVER_REQUEST_ACK(message_p).source_enb_id = source_enb_id;
  S1AP_HANDOVER_REQUEST_ACK(message_p).target_enb_id = target_enb_id;
  S1AP_HANDOVER_REQUEST_ACK(message_p).handover_type = handover_type;
  S1AP_HANDOVER_REQUEST_ACK(message_p).tgt_src_container = tgt_src_container;

  OAILOG_DEBUG_UE(
      LOG_S1AP, imsi64,
      "sending Handover Request Ack to MME_APP for mme_ue_s1ap_id %d\n",
      mme_ue_s1ap_id);

  message_p->ittiMsgHeader.imsi = imsi64;
  send_msg_to_task(&s1ap_task_zmq_ctx, TASK_MME_APP, message_p);
  OAILOG_FUNC_RETURN(LOG_S1AP, RETURNok);
}

status_code_e s1ap_mme_itti_s1ap_handover_notify(
    const mme_ue_s1ap_id_t mme_ue_s1ap_id,
    const oai::S1apHandoverState handover_state,
    const enb_ue_s1ap_id_t target_enb_ue_s1ap_id,
    const sctp_assoc_id_t target_sctp_assoc_id, const ecgi_t ecgi,
    imsi64_t imsi64) {
  MessageDef* message_p = NULL;
  message_p = itti_alloc_new_message(TASK_S1AP, S1AP_HANDOVER_NOTIFY);
  if (message_p == NULL) {
    OAILOG_ERROR_UE(LOG_S1AP, imsi64, "itti_alloc_new_message Failed");
    OAILOG_FUNC_RETURN(LOG_S1AP, RETURNerror);
  }

  S1AP_HANDOVER_NOTIFY(message_p).mme_ue_s1ap_id = mme_ue_s1ap_id;
  S1AP_HANDOVER_NOTIFY(message_p).target_enb_id =
      handover_state.target_enb_id();
  S1AP_HANDOVER_NOTIFY(message_p).target_sctp_assoc_id = target_sctp_assoc_id;
  S1AP_HANDOVER_NOTIFY(message_p).ecgi = ecgi;
  S1AP_HANDOVER_NOTIFY(message_p).target_enb_ue_s1ap_id = target_enb_ue_s1ap_id;
  e_rab_admitted_list_t* e_rab_admitted_list =
      &S1AP_HANDOVER_NOTIFY(message_p).e_rab_admitted_list;
  e_rab_admitted_list->no_of_items =
      handover_state.e_rab_admitted_list().no_of_items();
  for (uint8_t idx = 0; idx < e_rab_admitted_list->no_of_items; idx++) {
    const oai::ERabAdmittedItem& proto_e_rab_item =
        handover_state.e_rab_admitted_list().item(idx);
    e_rab_admitted_list->item[idx].e_rab_id = proto_e_rab_item.e_rab_id();
    e_rab_admitted_list->item[idx].transport_layer_address =
        blk2bstr(proto_e_rab_item.transport_layer_address().c_str(),
                 proto_e_rab_item.transport_layer_address().length());
    e_rab_admitted_list->item[idx].gtp_teid = proto_e_rab_item.gtp_teid();
  }

  message_p->ittiMsgHeader.imsi = imsi64;
  send_msg_to_task(&s1ap_task_zmq_ctx, TASK_MME_APP, message_p);
  OAILOG_FUNC_RETURN(LOG_S1AP, RETURNok);
}

}  // namespace lte
}  // namespace magma

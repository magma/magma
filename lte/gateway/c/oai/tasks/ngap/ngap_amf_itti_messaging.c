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
  SCTP_DATA_REQ(message_p).payload       = *payload;
  *payload                               = NULL;
  SCTP_DATA_REQ(message_p).assoc_id      = assoc_id;
  SCTP_DATA_REQ(message_p).stream        = stream;
  SCTP_DATA_REQ(message_p).agw_ue_xap_id = ue_id;
  SCTP_DATA_REQ(message_p).ppid          = NGAP_SCTP_PPID;
  return send_msg_to_task(&ngap_task_zmq_ctx, TASK_SCTP, message_p);
}

int ngap_amf_itti_nas_uplink_ind(
    const amf_ue_ngap_id_t ue_id, STOLEN_REF bstring* payload,
    const tai_t* const tai, const ecgi_t* const cgi) {
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
  return send_msg_to_task(&ngap_task_zmq_ctx, TASK_AMF_APP, message_p);
}

//------------------------------------------------------------------------------
void ngap_amf_itti_ngap_initial_ue_message(
    const sctp_assoc_id_t assoc_id, const uint32_t gnb_id,
    const gnb_ue_ngap_id_t gnb_ue_ngap_id, const uint8_t* const nas_msg,
    const size_t nas_msg_length, const tai_t* const tai,
    const ecgi_t* const ecgi, const long rrc_cause,
    const s_tmsi_m5_t* const opt_s_tmsi, const csg_id_t* const opt_csg_id,
    const guamfi_t* const opt_guamfi,
    const void* opt_cell_access_mode,           // unused
    const void* opt_cell_gw_transport_address,  // unused
    const void* opt_relay_node_indicator)       // unused
{
  MessageDef* message_p = NULL;

  OAILOG_FUNC_IN(LOG_NGAP);

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

  NGAP_INITIAL_UE_MESSAGE(message_p).sctp_assoc_id  = assoc_id;
  NGAP_INITIAL_UE_MESSAGE(message_p).gnb_ue_ngap_id = gnb_ue_ngap_id;
  NGAP_INITIAL_UE_MESSAGE(message_p).gnb_id         = gnb_id;
  NGAP_INITIAL_UE_MESSAGE(message_p).nas = blk2bstr(nas_msg, nas_msg_length);
  NGAP_INITIAL_UE_MESSAGE(message_p).m5g_rrc_establishment_cause =
      rrc_cause + 1;

  if (opt_s_tmsi) {
    NGAP_INITIAL_UE_MESSAGE(message_p).is_s_tmsi_valid = true;
    NGAP_INITIAL_UE_MESSAGE(message_p).opt_s_tmsi      = *opt_s_tmsi;
  } else {
    NGAP_INITIAL_UE_MESSAGE(message_p).is_s_tmsi_valid = false;
  }
#if 0
  if (opt_guamfi) {
    NGAP_INITIAL_UE_MESSAGE(message_p).is_guamfi_valid = true;
    NGAP_INITIAL_UE_MESSAGE(message_p).opt_guamfi = *opt_guamfi;
  } else {
    NGAP_INITIAL_UE_MESSAGE(message_p).is_guamfi_valid = false;
  }
#endif
  NGAP_INITIAL_UE_MESSAGE(message_p).gnb_ue_ngap_id = gnb_ue_ngap_id;

  send_msg_to_task(&ngap_task_zmq_ctx, TASK_AMF_APP, message_p);
  OAILOG_INFO(LOG_NGAP, "iniUEmsg sent to TASK_AMF_APP");
  OAILOG_FUNC_OUT(LOG_NGAP);
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

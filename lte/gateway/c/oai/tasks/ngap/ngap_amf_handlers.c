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
  Source      ngap_amf_handlers.c
  Version     0.1
  Date        2020/07/28
  Product     NGAP stack
  Subsystem   Access and Mobility Management Function
  Author      Ashish Prajapati
  Description Defines NG Application Protocol Messages Handlers

*****************************************************************************/

#include <netinet/in.h>
#include <stdbool.h>
#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/types.h>

#include "3gpp_23.003.h"
#include "3gpp_24.008.h"
#include "3gpp_38.401.h"
#include "3gpp_38.413.h"
#include "BIT_STRING.h"
#include "INTEGER.h"
#include "Ngap_AMF-UE-NGAP-ID.h"
#include "Ngap_CauseMisc.h"
#include "Ngap_CauseNas.h"
#include "Ngap_CauseProtocol.h"
#include "Ngap_CauseRadioNetwork.h"
#include "Ngap_CauseTransport.h"
#include "Ngap_FiveG-S-TMSI.h"
#include "Ngap_GNB-ID.h"
#include "Ngap_GTP-TEID.h"
#include "Ngap_GlobalGNB-ID.h"
#include "Ngap_NGAP-PDU.h"
#include "Ngap_PLMNIdentity.h"
#include "Ngap_ProcedureCode.h"
#include "Ngap_RAN-UE-NGAP-ID.h"
#include "Ngap_RANNodeName.h"
#include "Ngap_ResetType.h"
#include "Ngap_ServedGUAMIItem.h"
#include "Ngap_TAI.h"
#include "Ngap_TimeToWait.h"
#include "Ngap_TransportLayerAddress.h"
#include "Ngap_UE-NGAP-ID-pair.h"
#include "Ngap_UE-NGAP-IDs.h"
#include "Ngap_UE-associatedLogicalNG-connectionItem.h"
#include "Ngap_UE-associatedLogicalNG-connectionList.h"
#include "Ngap_UEAggregateMaximumBitRate.h"
#include "Ngap_UEPagingIdentity.c"
#include "Ngap_UERadioCapability.h"
#include "amf_app_messages_types.h"
#include "amf_config.h"
#include "asn_SEQUENCE_OF.h"
#include "assertions.h"
#include "bstrlib.h"
#include "common_defs.h"
#include "conversions.h"
#include "dynamic_memory_check.h"
#include "hashtable.h"
#include "intertask_interface.h"
#include "intertask_interface_types.h"
#include "itti_types.h"
#include "log.h"
#include "ngap_amf.h"
#include "ngap_amf_encoder.h"
#include "ngap_amf_handlers.h"
#include "ngap_amf_itti_messaging.h"
#include "ngap_amf_nas_procedures.h"
#include "ngap_amf_ta.h"
#include "ngap_common.h"
#include "ngap_state.h"
#include "service303.h"
#include "timer.h"

/* Handlers matrix. Only amf related procedures present here.
 */
ngap_message_handler_t ngap_message_handlers[][3] = {

    {ngap_amf_handle_ue_context_release_request, 0, 0}, /* UEContextRelease*/

};

int ngap_amf_handle_message(
    ngap_state_t* state, const sctp_assoc_id_t assoc_id,
    const sctp_stream_id_t stream, Ngap_NGAP_PDU_t* pdu) {
  /*
   * Checking procedure Code and direction of pdu
   */
  if (pdu->choice.initiatingMessage.procedureCode >=
          COUNT_OF(ngap_message_handlers) ||
      pdu->present > Ngap_NGAP_PDU_PR_unsuccessfulOutcome) {
    OAILOG_DEBUG(
        LOG_NGAP,
        "[SCTP %d] Either procedureCode %d or direction %d exceed expected\n",
        assoc_id, (int) pdu->choice.initiatingMessage.procedureCode,
        (int) pdu->present);
    return -1;
  }

  ngap_message_handler_t handler =
      ngap_message_handlers[pdu->choice.initiatingMessage.procedureCode]
                           [pdu->present - 1];

  if (handler == NULL) {
    // not implemented or no procedure for gNB (wrong message)
    OAILOG_DEBUG(
        LOG_NGAP, "[SCTP %d] No handler for procedureCode %d in %s\n", assoc_id,
        (int) pdu->choice.initiatingMessage.procedureCode,
        ngap_direction2str(pdu->present));
    return -2;
  }

  return handler(state, assoc_id, stream, pdu);
}

//------------------------------------------------------------------------------
int ngap_amf_handle_ue_context_release_request(
    ngap_state_t* state, __attribute__((unused)) const sctp_assoc_id_t assoc_id,
    __attribute__((unused)) const sctp_stream_id_t stream,
    Ngap_NGAP_PDU_t* pdu) {
  Ngap_UEContextReleaseRequest_t* container;
  Ngap_UEContextReleaseRequest_IEs_t* ie = NULL;
  m5g_ue_description_t* ue_ref_p         = NULL;
  MessageDef* message_p                  = NULL;
  Ngap_Cause_PR cause_type;
  long cause_value;
  enum Ngcause ng_release_cause   = NGAP_RADIO_NR_GENERATED_REASON;
  int rc                          = RETURNok;
  amf_ue_ngap_id_t amf_ue_ngap_id = 0;
  gnb_ue_ngap_id_t gnb_ue_ngap_id = 0;
  imsi64_t imsi64                 = INVALID_IMSI64;

  message_p = itti_alloc_new_message(TASK_NGAP, NGAP_UE_CONTEXT_RELEASE_REQ);
  AssertFatal(message_p != NULL, "itti_alloc_new_message Failed");

  OAILOG_FUNC_IN(LOG_NGAP);
  container =
      &pdu->choice.initiatingMessage.value.choice.UEContextReleaseRequest;
  // Log the Cause Type and Cause value
  NGAP_FIND_PROTOCOLIE_BY_ID(
      Ngap_UEContextReleaseRequest_IEs_t, ie, container,
      Ngap_ProtocolIE_ID_id_AMF_UE_NGAP_ID, true);
  if (ie) {
    amf_ue_ngap_id = ie->value.choice.AMF_UE_NGAP_ID;
  } else {
    OAILOG_FUNC_RETURN(LOG_NGAP, RETURNerror);
  }

  NGAP_FIND_PROTOCOLIE_BY_ID(
      Ngap_UEContextReleaseRequest_IEs_t, ie, container,
      Ngap_ProtocolIE_ID_id_RAN_UE_NGAP_ID, true);
  if (ie) {
    gnb_ue_ngap_id = (gnb_ue_ngap_id_t)(
        ie->value.choice.AMF_UE_NGAP_ID & GNB_UE_NGAP_ID_MASK);
  } else {
    OAILOG_FUNC_RETURN(LOG_NGAP, RETURNerror);
  }

  /*Request list to release*/
  NGAP_FIND_PROTOCOLIE_BY_ID(
      Ngap_UEContextReleaseRequest_IEs_t, ie, container,
      Ngap_ProtocolIE_ID_id_PDUSessionResourceListCxtRelReq, false);
  if (ie) {
    NGAP_UE_CONTEXT_RELEASE_REQ(message_p).pduSession.pduSessionItemCount =
        ie->value.choice.PDUSessionResourceListCxtRelReq.list.count;

    for (int i = 0;
         i < ie->value.choice.PDUSessionResourceListCxtRelReq.list.count; i++) {
      struct Ngap_PDUSessionResourceItemCxtRelReq* RelItem =
          ie->value.choice.PDUSessionResourceListCxtRelReq.list.array[i];

      NGAP_UE_CONTEXT_RELEASE_REQ(message_p).pduSession.pduSessionIDs[i] =
          RelItem->pDUSessionID;
    }
  }

  // Log the Cause Type and Cause value
  NGAP_FIND_PROTOCOLIE_BY_ID(
      Ngap_UEContextReleaseRequest_IEs_t, ie, container,
      Ngap_ProtocolIE_ID_id_Cause, true);
  if (ie) {
    cause_type = ie->value.choice.Cause.present;
  } else {
    OAILOG_FUNC_RETURN(LOG_NGAP, RETURNerror);
  }
  switch (cause_type) {
    case Ngap_Cause_PR_radioNetwork:
      cause_value = ie->value.choice.Cause.choice.radioNetwork;
      OAILOG_INFO(
          LOG_NGAP,
          "UE CONTEXT RELEASE REQUEST with Cause_Type = Radio Network and "
          "Cause_Value = %ld\n",
          cause_value);
      if (cause_value == Ngap_CauseRadioNetwork_user_inactivity) {
        increment_counter(
            "ue_context_release_req", 1, 1, "cause", "user_inactivity");
      } else if (
          cause_value == Ngap_CauseRadioNetwork_radio_connection_with_ue_lost) {
        increment_counter(
            "ue_context_release_req", 1, 1, "cause", "radio_link_failure");
      }
      break;

    case Ngap_Cause_PR_transport:
      cause_value = ie->value.choice.Cause.choice.transport;
      OAILOG_INFO(
          LOG_NGAP,
          "UE CONTEXT RELEASE REQUEST with Cause_Type = Transport and "
          "Cause_Value = %ld\n",
          cause_value);
      break;

    case Ngap_Cause_PR_nas:
      cause_value = ie->value.choice.Cause.choice.nas;
      OAILOG_INFO(
          LOG_NGAP,
          "UE CONTEXT RELEASE REQUEST with Cause_Type = NAS and Cause_Value = "
          "%ld\n",
          cause_value);
      break;

    case Ngap_Cause_PR_protocol:
      cause_value = ie->value.choice.Cause.choice.protocol;
      OAILOG_INFO(
          LOG_NGAP,
          "UE CONTEXT RELEASE REQUEST with Cause_Type = Transport and "
          "Cause_Value = %ld\n",
          cause_value);
      break;

    case Ngap_Cause_PR_misc:
      cause_value = ie->value.choice.Cause.choice.misc;
      OAILOG_INFO(
          LOG_NGAP,
          "UE CONTEXT RELEASE REQUEST with Cause_Type = MISC and Cause_Value = "
          "%ld\n",
          cause_value);
      break;

    default:
      OAILOG_ERROR(
          LOG_NGAP, "UE CONTEXT RELEASE REQUEST with Invalid Cause_Type = %d\n",
          cause_type);
      OAILOG_FUNC_RETURN(LOG_NGAP, RETURNerror);
  }

  if ((ue_ref_p = ngap_state_get_ue_amfid(amf_ue_ngap_id)) == NULL) {
    /*
     * AMF doesn't know the AMF UE NGAP ID provided.
     * No need to do anything. Ignore the message
     */
    OAILOG_DEBUG(
        LOG_NGAP,
        "UE_CONTEXT_RELEASE_REQUEST ignored cause could not get context with "
        "amf_ue_ngap_id " AMF_UE_NGAP_ID_FMT
        " gnb_ue_ngap_id " GNB_UE_NGAP_ID_FMT " ",
        (uint32_t) amf_ue_ngap_id, (uint32_t) gnb_ue_ngap_id);
    OAILOG_FUNC_RETURN(LOG_NGAP, RETURNerror);
  } else {
    if (ue_ref_p->gnb_ue_ngap_id == gnb_ue_ngap_id) {
      /*
       * Both gNB UE NGAP ID and AMF UE NGAP ID match.
       * Send a UE context Release Command to gNB after releasing Ng-U bearer
       * tunnel mapping for all the bearers.
       */
      ngap_imsi_map_t* imsi_map = get_ngap_imsi_map();
      hashtable_uint64_ts_get(
          imsi_map->amf_ue_id_imsi_htbl, (const hash_key_t) amf_ue_ngap_id,
          &imsi64);

      gnb_ref_p = ngap_state_get_gnb(state, ue_ref_p->sctp_assoc_id);

      NGAP_UE_CONTEXT_RELEASE_REQ(message_p).amf_ue_ngap_id =
          ue_ref_p->amf_ue_ngap_id;
      NGAP_UE_CONTEXT_RELEASE_REQ(message_p).gnb_ue_ngap_id =
          ue_ref_p->gnb_ue_ngap_id;
      NGAP_UE_CONTEXT_RELEASE_REQ(message_p).relCause = ng_release_cause;
      NGAP_UE_CONTEXT_RELEASE_REQ(message_p).cause    = ie->value.choice.Cause;

      message_p->ittiMsgHeader.imsi = imsi64;
      rc = send_msg_to_task(&ngap_task_zmq_ctx, TASK_AMF_APP, message_p);
      OAILOG_FUNC_RETURN(LOG_NGAP, rc);
    } else {
      // abnormal case. No need to do anything. Ignore the message
      OAILOG_DEBUG_UE(
          LOG_NGAP, imsi64,
          "UE_CONTEXT_RELEASE_REQUEST ignored, cause mismatch gnb_ue_ngap_id: "
          "ctxt " GNB_UE_NGAP_ID_FMT " != request " GNB_UE_NGAP_ID_FMT " ",
          (uint32_t) ue_ref_p->gnb_ue_ngap_id, (uint32_t) gnb_ue_ngap_id);
      OAILOG_FUNC_RETURN(LOG_NGAP, RETURNerror);
    }
  }
  OAILOG_FUNC_RETURN(LOG_NGAP, RETURNok);
}

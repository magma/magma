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

/*! \file s1ap_mme_handlers.c
  \brief
  \author Sebastien ROUX, Lionel Gauthier
  \company Eurecom
  \email: lionel.gauthier@eurecom.fr
*/

#include <stdlib.h>
#include <stdio.h>
#include <stdbool.h>
#include <stdint.h>
#include <netinet/in.h>
#include <string.h>

#include "S1ap_ProtocolIE-Field.h"
#include "lte/gateway/c/core/oai/lib/bstr/bstrlib.h"
#include "lte/gateway/c/core/oai/lib/hashtable/hashtable.h"
#include "lte/gateway/c/core/oai/common/log.h"
#include "lte/gateway/c/core/oai/common/assertions.h"
#include "lte/gateway/c/core/oai/common/conversions.h"
#include "lte/gateway/c/core/oai/lib/itti/intertask_interface.h"
#include "lte/gateway/c/core/oai/common/dynamic_memory_check.h"
#include "lte/gateway/c/core/oai/include/mme_config.h"
#include "lte/gateway/c/core/oai/tasks/s1ap/s1ap_common.h"
#include "lte/gateway/c/core/oai/tasks/s1ap/s1ap_mme_encoder.h"
#include "lte/gateway/c/core/oai/tasks/s1ap/s1ap_mme_nas_procedures.h"
#include "lte/gateway/c/core/oai/tasks/s1ap/s1ap_mme_itti_messaging.h"
#include "lte/gateway/c/core/oai/tasks/s1ap/s1ap_mme.h"
#include "lte/gateway/c/core/oai/tasks/s1ap/s1ap_mme_ta.h"
#include "lte/gateway/c/core/oai/tasks/s1ap/s1ap_mme_handlers.h"
#include "lte/gateway/c/core/oai/include/mme_events.h"
#include "lte/gateway/c/core/oai/lib/3gpp/3gpp_23.003.h"
#include "lte/gateway/c/core/oai/lib/3gpp/3gpp_36.401.h"
#include "lte/gateway/c/core/oai/lib/3gpp/3gpp_36.413.h"
#include "BIT_STRING.h"
#include "INTEGER.h"
#include "S1ap_S1AP-PDU.h"
#include "S1ap_CNDomain.h"
#include "S1ap_CauseMisc.h"
#include "S1ap_CauseNas.h"
#include "S1ap_CauseProtocol.h"
#include "S1ap_CauseRadioNetwork.h"
#include "S1ap_CauseTransport.h"
#include "S1ap_E-RABItem.h"
#include "S1ap_E-RABSetupItemBearerSURes.h"
#include "S1ap_E-RABSetupItemCtxtSURes.h"
#include "S1ap_ENB-ID.h"
#include "S1ap_ENB-UE-S1AP-ID.h"
#include "S1ap_ENBname.h"
#include "S1ap_GTP-TEID.h"
#include "S1ap_Global-ENB-ID.h"
#include "S1ap_LAI.h"
#include "S1ap_MME-Code.h"
#include "S1ap_MME-Group-ID.h"
#include "S1ap_MME-UE-S1AP-ID.h"
#include "S1ap_PLMNidentity.h"
#include "S1ap_ProcedureCode.h"
#include "S1ap_ResetType.h"
#include "S1ap_S-TMSI.h"
#include "S1ap_ServedGUMMEIsItem.h"
#include "S1ap_ServedGroupIDs.h"
#include "S1ap_ServedMMECs.h"
#include "S1ap_ServedPLMNs.h"
#include "S1ap_TAI.h"
#include "S1ap_TAIItem.h"
#include "S1ap_TimeToWait.h"
#include "S1ap_TransportLayerAddress.h"
#include "S1ap_UE-S1AP-ID-pair.h"
#include "S1ap_UE-S1AP-IDs.h"
#include "S1ap_UE-associatedLogicalS1-ConnectionItem.h"
#include "S1ap_UE-associatedLogicalS1-ConnectionListRes.h"
#include "S1ap_UEAggregateMaximumBitrate.h"
#include "S1ap_UEPagingID.h"
#include "S1ap_UERadioCapability.h"
#include "asn_SEQUENCE_OF.h"
#include "lte/gateway/c/core/oai/common/common_defs.h"
#include "lte/gateway/c/core/oai/lib/itti/intertask_interface_types.h"
#include "lte/gateway/c/core/oai/include/mme_app_messages_types.h"
#include "orc8r/gateway/c/common/service303/includes/MetricsHelpers.h"
#include "lte/gateway/c/core/oai/include/s1ap_state.h"

struct S1ap_E_RABItem_s;
struct S1ap_E_RABSetupItemBearerSURes_s;
struct S1ap_E_RABSetupItemCtxtSURes_s;
struct S1ap_IE;

status_code_e s1ap_generate_s1_setup_response(
    s1ap_state_t* state, enb_description_t* enb_association);

bool is_all_erabId_same(S1ap_PathSwitchRequest_t* container);

/* Handlers matrix. Only mme related procedures present here.
 */
s1ap_message_handler_t message_handlers[][3] = {
    {s1ap_mme_handle_handover_required, 0, 0}, /* HandoverPreparation */
    {0, s1ap_mme_handle_handover_request_ack,
     s1ap_mme_handle_handover_failure},      /* HandoverResourceAllocation */
    {s1ap_mme_handle_handover_notify, 0, 0}, /* HandoverNotification */
    {s1ap_mme_handle_path_switch_request, 0, 0}, /* PathSwitchRequest */
    {s1ap_mme_handle_handover_cancel, 0, 0},     /* HandoverCancel */
    {0, s1ap_mme_handle_erab_setup_response,
     s1ap_mme_handle_erab_setup_failure},      /* E_RABSetup */
    {0, 0, 0},                                 /* E_RABModify */
    {0, s1ap_mme_handle_erab_rel_response, 0}, /* E_RABRelease */
    {0, 0, 0},                                 /* E_RABReleaseIndication */
    {0, s1ap_mme_handle_initial_context_setup_response,
     s1ap_mme_handle_initial_context_setup_failure}, /* InitialContextSetup */
    {0, 0, 0},                                       /* Paging */
    {0, 0, 0},                                       /* downlinkNASTransport */
    {s1ap_mme_handle_initial_ue_message, 0, 0},      /* initialUEMessage */
    {s1ap_mme_handle_uplink_nas_transport, 0, 0},    /* uplinkNASTransport */
    {s1ap_mme_handle_enb_reset, 0, 0},               /* Reset */
    {s1ap_mme_handle_error_ind_message, 0, 0},       /* ErrorIndication */
    {s1ap_mme_handle_nas_non_delivery, 0, 0}, /* NASNonDeliveryIndication */
    {s1ap_mme_handle_s1_setup_request, 0, 0}, /* S1Setup */
    {s1ap_mme_handle_ue_context_release_request, 0,
     0},       /* UEContextReleaseRequest */
    {0, 0, 0}, /* DownlinkS1cdma2000tunneling */
    {0, 0, 0}, /* UplinkS1cdma2000tunneling */
    {0, s1ap_mme_handle_ue_context_modification_response,
     s1ap_mme_handle_ue_context_modification_failure}, /* UEContextModification
                                                        */
    {s1ap_mme_handle_ue_cap_indication, 0, 0}, /* UECapabilityInfoIndication */
    {s1ap_mme_handle_ue_context_release_request,
     s1ap_mme_handle_ue_context_release_complete, 0}, /* UEContextRelease */
    {s1ap_mme_handle_enb_status_transfer, 0, 0},      /* eNBStatusTransfer */
    {0, 0, 0},                                        /* MMEStatusTransfer */
    {0, 0, 0},                                        /* DeactivateTrace */
    {0, 0, 0},                                        /* TraceStart */
    {0, 0, 0}, /* TraceFailureIndication */
    {0, 0, 0}, /* ENBConfigurationUpdate */
    {0, 0, 0}, /* MMEConfigurationUpdate */
    {0, 0, 0}, /* LocationReportingControl */
    {0, 0, 0}, /* LocationReportingFailureIndication */
    {0, 0, 0}, /* LocationReport */
    {0, 0, 0}, /* OverloadStart */
    {0, 0, 0}, /* OverloadStop */
    {0, 0, 0}, /* WriteReplaceWarning */
    {0, 0, 0}, /* eNBDirectInformationTransfer */
    {0, 0, 0}, /* MMEDirectInformationTransfer */
    {0, 0, 0}, /* PrivateMessage */
    {s1ap_mme_handle_enb_configuration_transfer, 0,
     0},       /* eNBConfigurationTransfer */
    {0, 0, 0}, /* MMEConfigurationTransfer */
    // HERE END UPDATES REL 9, 10, 11
    // UPDATE RELEASE 12 :
    {0, 0, 0}, /* CellTrafficTrace */
    {0, 0, 0}, /* 43 Kill */
    {0, 0, 0}, /* 44 DownlinkUEAssociatedLPPaTransport  */
    {0, 0, 0}, /* 45 UplinkUEAssociatedLPPaTransport */
    {0, 0, 0}, /* 46 DownlinkNonUEAssociatedLPPaTransport */
    {0, 0, 0}, /* 47 UplinkNonUEAssociatedLPPaTransport */
    {0, 0, 0}, /* 48 UERadioCapabilityMatch */
    {0, 0, 0}, /* 49 PWSRestartIndication */
    // UPDATE RELEASE 13 :
    {s1ap_mme_handle_erab_modification_indication, 0,
     0},       /* 50 E-RABModificationIndication */
    {0, 0, 0}, /* 51 PWSFailureIndication */
    {0, 0, 0}, /* 52 RerouteNASRequest */
    {0, 0, 0}, /* 53 UEContextModificationIndication */
    {0, 0, 0}, /* 54 ConnectionEstablishmentIndication */
    // UPDATE RELEASE 14 :
    {0, 0, 0}, /* 55 UEContextSuspend */
    {0, 0, 0}, /* 56 UEContextResume */
    {0, 0, 0}, /* 57 NASDeliveryIndication */
    {0, 0, 0}, /* 58 RetrieveUEInformation */
    {0, 0, 0}, /* 59 UEInformationTransfer */
    // UPDATE RELEASE 15 :
    {0, 0, 0}, /* 60 eNBCPRelocationIndication */
    {0, 0, 0}, /* 61 MMECPRelocationIndication */
    {0, 0, 0}, /* 62 SecondaryRATDataUsageReport */
};

status_code_e s1ap_mme_handle_message(
    s1ap_state_t* state, const sctp_assoc_id_t assoc_id,
    const sctp_stream_id_t stream, S1ap_S1AP_PDU_t* pdu) {
  /*
   * Checking procedure Code and present flag of message
   */
  if (pdu->choice.initiatingMessage.procedureCode >=
          COUNT_OF(message_handlers) ||
      pdu->present > S1ap_S1AP_PDU_PR_unsuccessfulOutcome) {
    OAILOG_DEBUG(
        LOG_S1AP,
        "[SCTP %d] Either procedureCode %d or present flag %d exceed "
        "expected\n",
        assoc_id, (int) pdu->choice.initiatingMessage.procedureCode,
        (int) pdu->present);
    return RETURNerror;
  }

  s1ap_message_handler_t handler =
      message_handlers[pdu->choice.initiatingMessage.procedureCode]
                      [pdu->present - 1];

  if (handler == NULL) {
    // not implemented or no procedure for eNB (wrong message)
    OAILOG_DEBUG(
        LOG_S1AP, "[SCTP %d] No handler for procedureCode %d in %s\n", assoc_id,
        (int) pdu->choice.initiatingMessage.procedureCode,
        s1ap_direction2str(pdu->present));
    return RETURNerror;
  }

  return handler(state, assoc_id, stream, pdu);
}

//------------------------------------------------------------------------------
status_code_e s1ap_mme_set_cause(
    S1ap_Cause_t* cause_p, const S1ap_Cause_PR cause_type,
    const long cause_value) {
  if (cause_p == NULL) {
    OAILOG_ERROR(LOG_S1AP, "Cause value is NULL\n");
    return RETURNerror;
  }
  cause_p->present = cause_type;

  switch (cause_type) {
    case S1ap_Cause_PR_radioNetwork:
      cause_p->choice.misc = cause_value;
      break;

    case S1ap_Cause_PR_transport:
      cause_p->choice.transport = cause_value;
      break;

    case S1ap_Cause_PR_nas:
      cause_p->choice.nas = cause_value;
      break;

    case S1ap_Cause_PR_protocol:
      cause_p->choice.protocol = cause_value;
      break;

    case S1ap_Cause_PR_misc:
      cause_p->choice.misc = cause_value;
      break;

    default:
      return RETURNerror;
  }

  return RETURNok;
}

long s1ap_mme_get_cause_value(S1ap_Cause_t* cause) {
  S1ap_Cause_PR cause_type = {0};
  long cause_value         = RETURNerror;

  OAILOG_FUNC_IN(LOG_S1AP);
  if (cause == NULL) {
    OAILOG_ERROR(LOG_S1AP, "Cause is NULL\n");
    OAILOG_FUNC_RETURN(LOG_S1AP, RETURNerror);
  }

  cause_type = cause->present;

  switch (cause_type) {
    case S1ap_Cause_PR_radioNetwork:
      cause_value = cause->choice.radioNetwork;
      break;

    case S1ap_Cause_PR_transport:
      cause_value = cause->choice.transport;
      break;

    case S1ap_Cause_PR_nas:
      cause_value = cause->choice.nas;
      break;

    case S1ap_Cause_PR_protocol:
      cause_value = cause->choice.protocol;
      break;

    case S1ap_Cause_PR_misc:
      cause_value = cause->choice.misc;
      break;

    default:
      OAILOG_ERROR(LOG_S1AP, "Invalid Cause_Type = %d\n", cause_type);
      OAILOG_FUNC_RETURN(LOG_S1AP, RETURNerror);
  }

  OAILOG_FUNC_RETURN(LOG_S1AP, cause_value);
}

//------------------------------------------------------------------------------
status_code_e s1ap_mme_generate_s1_setup_failure(
    const sctp_assoc_id_t assoc_id, const S1ap_Cause_PR cause_type,
    const long cause_value, const long time_to_wait) {
  uint8_t* buffer_p = 0;
  uint32_t length   = 0;
  S1ap_S1AP_PDU_t pdu;
  S1ap_S1SetupFailure_t* out;
  S1ap_S1SetupFailureIEs_t* ie = NULL;
  int rc                       = RETURNok;

  OAILOG_FUNC_IN(LOG_S1AP);

  memset(&pdu, 0, sizeof(pdu));
  pdu.present = S1ap_S1AP_PDU_PR_unsuccessfulOutcome;
  pdu.choice.unsuccessfulOutcome.procedureCode = S1ap_ProcedureCode_id_S1Setup;
  pdu.choice.unsuccessfulOutcome.criticality   = S1ap_Criticality_reject;
  pdu.choice.unsuccessfulOutcome.value.present =
      S1ap_UnsuccessfulOutcome__value_PR_S1SetupFailure;
  out = &pdu.choice.unsuccessfulOutcome.value.choice.S1SetupFailure;

  ie = (S1ap_S1SetupFailureIEs_t*) calloc(1, sizeof(S1ap_S1SetupFailureIEs_t));
  ie->id            = S1ap_ProtocolIE_ID_id_Cause;
  ie->criticality   = S1ap_Criticality_ignore;
  ie->value.present = S1ap_S1SetupFailureIEs__value_PR_Cause;
  s1ap_mme_set_cause(&ie->value.choice.Cause, cause_type, cause_value);
  ASN_SEQUENCE_ADD(&out->protocolIEs.list, ie);

  /*
   * Include the optional field time to wait only if the value is > -1
   */
  if (time_to_wait > -1) {
    ie =
        (S1ap_S1SetupFailureIEs_t*) calloc(1, sizeof(S1ap_S1SetupFailureIEs_t));
    ie->id                      = S1ap_ProtocolIE_ID_id_TimeToWait;
    ie->criticality             = S1ap_Criticality_ignore;
    ie->value.present           = S1ap_S1SetupFailureIEs__value_PR_TimeToWait;
    ie->value.choice.TimeToWait = time_to_wait;
    ASN_SEQUENCE_ADD(&out->protocolIEs.list, ie);
  }

  if (s1ap_mme_encode_pdu(&pdu, &buffer_p, &length) < 0) {
    OAILOG_ERROR(LOG_S1AP, "Failed to encode s1 setup failure\n");
    OAILOG_FUNC_RETURN(LOG_S1AP, RETURNerror);
  }

  bstring b = blk2bstr(buffer_p, (int) length);
  free(buffer_p);
  rc = s1ap_mme_itti_send_sctp_request(&b, assoc_id, 0, INVALID_MME_UE_S1AP_ID);
  OAILOG_FUNC_RETURN(LOG_S1AP, rc);
}

////////////////////////////////////////////////////////////////////////////////
//************************** Management procedures ***************************//
////////////////////////////////////////////////////////////////////////////////

//------------------------------------------------------------------------------
bool get_stale_enb_connection_with_enb_id(
    __attribute__((unused)) const hash_key_t keyP, void* const elementP,
    void* parameterP, void** resultP) {
  enb_description_t* new_enb_association = (enb_description_t*) parameterP;
  enb_description_t* ht_enb_association  = (enb_description_t*) elementP;

  // No need to clean the newly created eNB association
  if (ht_enb_association == new_enb_association) {
    return false;
  }

  // Match old and new association with respect to eNB id
  if (ht_enb_association->enb_id == new_enb_association->enb_id) {
    *resultP = elementP;
    return true;
  }

  return false;
}

void clean_stale_enb_state(
    s1ap_state_t* state, enb_description_t* new_enb_association) {
  enb_description_t* stale_enb_association = NULL;

  hashtable_ts_apply_callback_on_elements(
      &state->enbs, get_stale_enb_connection_with_enb_id, new_enb_association,
      (void**) &stale_enb_association);
  if (stale_enb_association == NULL) {
    // No stale eNB connection found;
    return;
  }

  OAILOG_INFO(
      LOG_S1AP, "Found stale eNB at association id %d",
      stale_enb_association->sctp_assoc_id);
  // Remove the S1 context for UEs associated with old eNB association
  hashtable_key_array_t* keys =
      hashtable_uint64_ts_get_keys(&stale_enb_association->ue_id_coll);
  if (keys != NULL) {
    ue_description_t* ue_ref = NULL;
    for (int i = 0; i < keys->num_keys; i++) {
      ue_ref = s1ap_state_get_ue_mmeid((mme_ue_s1ap_id_t) keys->keys[i]);
      /* The function s1ap_remove_ue will take care of removing the enb also,
       * when the last UE is removed
       */
      s1ap_remove_ue(state, ue_ref);
    }
    FREE_HASHTABLE_KEY_ARRAY(keys);
  } else {
    // Remove the old eNB association
    OAILOG_INFO(
        LOG_S1AP, "Deleting eNB: %s (Sctp_assoc_id = %u)",
        stale_enb_association->enb_name, stale_enb_association->sctp_assoc_id);
    s1ap_remove_enb(state, stale_enb_association);
  }

  OAILOG_DEBUG(LOG_S1AP, "Removed stale eNB and all associated UEs.");
}

status_code_e s1ap_mme_handle_s1_setup_request(
    s1ap_state_t* state, const sctp_assoc_id_t assoc_id,
    const sctp_stream_id_t stream, S1ap_S1AP_PDU_t* pdu) {
  int rc = RETURNok;

  S1ap_S1SetupRequest_t* container                = NULL;
  S1ap_S1SetupRequestIEs_t* ie                    = NULL;
  S1ap_S1SetupRequestIEs_t* ie_enb_name           = NULL;
  S1ap_S1SetupRequestIEs_t* ie_supported_tas      = NULL;
  S1ap_S1SetupRequestIEs_t* ie_default_paging_drx = NULL;

  enb_description_t* enb_association = NULL;
  uint32_t enb_id                    = 0;
  char* enb_name                     = NULL;
  int ta_ret                         = 0;
  uint8_t bplmn_list_count           = 0;  // Broadcast PLMN list count

  OAILOG_FUNC_IN(LOG_S1AP);
  increment_counter("s1_setup", 1, NO_LABELS);
  if (!hss_associated) {
    /*
     * Can not process the request, MME is not connected to HSS
     */
    OAILOG_ERROR(
        LOG_S1AP,
        "Rejecting s1 setup request Can not process the request, MME is not "
        "connected to HSS\n");
    rc = s1ap_mme_generate_s1_setup_failure(
        assoc_id, S1ap_Cause_PR_misc, S1ap_CauseMisc_unspecified, -1);
    increment_counter(
        "s1_setup", 1, 2, "result", "failure", "cause", "s6a_interface_not_up");
    OAILOG_FUNC_RETURN(LOG_S1AP, rc);
  }

  if (pdu == NULL) {
    OAILOG_ERROR(LOG_S1AP, "PDU is NULL\n");
    return RETURNerror;
  }
  container = &pdu->choice.initiatingMessage.value.choice.S1SetupRequest;
  /*
   * We received a new valid S1 Setup Request on a stream != 0.
   * This should not happen -> reject eNB s1 setup request.
   */

  if (stream != 0) {
    OAILOG_ERROR(LOG_S1AP, "Received new s1 setup request on stream != 0\n");
    /*
     * Send a s1 setup failure with protocol cause unspecified
     */
    rc = s1ap_mme_generate_s1_setup_failure(
        assoc_id, S1ap_Cause_PR_protocol, S1ap_CauseProtocol_unspecified, -1);
    increment_counter(
        "s1_setup", 1, 2, "result", "failure", "cause",
        "sctp_stream_id_non_zero");
    OAILOG_FUNC_RETURN(LOG_S1AP, rc);
  }

  /* Handling of s1setup cases as follows.
   * If we don't know about the association, we haven't processed the new
   * association yet, so hope the eNB will retry the s1 setup. Ignore and
   * return. If we get this message when the S1 interface of the MME state is in
   * READY state then it is protocol error or out of sync state. Ignore it and
   * return. Assume MME would detect SCTP association failure and would S1
   * interface state to accept S1setup from eNB. If we get this message when the
   * s1 interface of the MME is in SHUTDOWN stage, we just hope the eNB will
   * retry and that will result in a new association getting established
   * followed by a subsequent s1 setup, return S1ap_TimeToWait_v20s. If we get
   * this message when the s1 interface of the MME is in RESETTING stage then we
   * return S1ap_TimeToWait_v20s.
   */
  if ((enb_association = s1ap_state_get_enb(state, assoc_id)) == NULL) {
    /*
     *
     * This should not happen as the thread processing new associations is the
     * one that reads data from the socket. Promote to an assert once we have
     * more confidence.
     */
    OAILOG_ERROR(LOG_S1AP, "Ignoring s1 setup from unknown assoc %u", assoc_id);
    OAILOG_FUNC_RETURN(LOG_S1AP, RETURNok);
  }

  if (enb_association->s1_state == S1AP_RESETING ||
      enb_association->s1_state == S1AP_SHUTDOWN) {
    OAILOG_WARNING(
        LOG_S1AP, "Ignoring s1setup from eNB in state %s on assoc id %u",
        s1_enb_state2str(enb_association->s1_state), assoc_id);
    OAILOG_DEBUG(
        LOG_S1AP, "Num UEs associated %u num ue_id_coll %zu",
        enb_association->nb_ue_associated,
        enb_association->ue_id_coll.num_elements);
    rc = s1ap_mme_generate_s1_setup_failure(
        assoc_id, S1ap_Cause_PR_transport,
        S1ap_CauseTransport_transport_resource_unavailable,
        S1ap_TimeToWait_v20s);
    increment_counter(
        "s1_setup", 1, 2, "result", "failure", "cause", "invalid_state");
    // Check if the UE counters for eNB are equal.
    // If not, the eNB will never switch to INIT state, particularly in
    // stateless mode.
    // Exit the process so that health checker can clean-up all Redis
    // state and restart all stateless services.
    AssertFatal(
        enb_association->nb_ue_associated ==
            enb_association->ue_id_coll.num_elements,
        "Num UEs associated with eNB (%u) is more than the UEs with valid "
        "mme_ue_s1ap_id (%zu). This is a deadlock state potentially caused by "
        "misbehaving eNB; restarting MME. In stateless mode, health management "
        "service will eventually detect multiple MME restarts due to this "
        "deadlock state and force sctpd and hence all services to restart.",
        enb_association->nb_ue_associated,
        enb_association->ue_id_coll.num_elements);
    OAILOG_FUNC_RETURN(LOG_S1AP, rc);
  }
  log_queue_item_t* context = NULL;
  OAILOG_MESSAGE_START_SYNC(
      OAILOG_LEVEL_DEBUG, LOG_S1AP, (&context),
      "New s1 setup request incoming from ");
  // shared_log_queue_item_t *context = NULL;
  // OAILOG_MESSAGE_START_ASYNC (OAILOG_LEVEL_DEBUG, LOG_S1AP, (&context), "New
  // s1 setup request incoming from ");

  S1AP_FIND_PROTOCOLIE_BY_ID(
      S1ap_S1SetupRequestIEs_t, ie_enb_name, container,
      S1ap_ProtocolIE_ID_id_eNBname, false);
  if (ie_enb_name) {
    OAILOG_MESSAGE_ADD_SYNC(
        context, "%*s ", (int) ie_enb_name->value.choice.ENBname.size,
        ie_enb_name->value.choice.ENBname.buf);
    enb_name = (char*) ie_enb_name->value.choice.ENBname.buf;
  }

  S1AP_FIND_PROTOCOLIE_BY_ID(
      S1ap_S1SetupRequestIEs_t, ie, container,
      S1ap_ProtocolIE_ID_id_Global_ENB_ID, true);
  if (ie->value.choice.Global_ENB_ID.eNB_ID.present ==
      S1ap_ENB_ID_PR_homeENB_ID) {
    // Home eNB ID = 28 bits
    uint8_t* enb_id_buf =
        ie->value.choice.Global_ENB_ID.eNB_ID.choice.homeENB_ID.buf;

    if (ie->value.choice.Global_ENB_ID.eNB_ID.choice.homeENB_ID.size != 28) {
      // TODO: handle case were size != 28 -> notify ? reject ?
    }

    enb_id = (enb_id_buf[0] << 20) + (enb_id_buf[1] << 12) +
             (enb_id_buf[2] << 4) + ((enb_id_buf[3] & 0xf0) >> 4);
    OAILOG_MESSAGE_ADD_SYNC(context, "home eNB id: %07x", enb_id);
  } else {
    // Macro eNB = 20 bits
    uint8_t* enb_id_buf =
        ie->value.choice.Global_ENB_ID.eNB_ID.choice.macroENB_ID.buf;

    if (ie->value.choice.Global_ENB_ID.eNB_ID.choice.macroENB_ID.size != 20) {
      // TODO: handle case were size != 20 -> notify ? reject ?
    }

    enb_id = (enb_id_buf[0] << 12) + (enb_id_buf[1] << 4) +
             ((enb_id_buf[2] & 0xf0) >> 4);
    OAILOG_MESSAGE_ADD_SYNC(context, "macro eNB id: %05x", enb_id);
  }

  OAILOG_MESSAGE_FINISH((void*) context);

  /* Requirement MME36.413R10_8.7.3.4 Abnormal Conditions
   * If the eNB initiates the procedure by sending a S1 SETUP REQUEST message
   * including the PLMN Identity IEs and none of the PLMNs provided by the eNB
   * is identified by the MME, then the MME shall reject the eNB S1 Setup
   * Request procedure with the appropriate cause value, e.g, Unknown PLMN.
   */
  S1AP_FIND_PROTOCOLIE_BY_ID(
      S1ap_S1SetupRequestIEs_t, ie_supported_tas, container,
      S1ap_ProtocolIE_ID_id_SupportedTAs, true);

  ta_ret =
      s1ap_mme_compare_ta_lists(&ie_supported_tas->value.choice.SupportedTAs);

  /*
   * eNB and MME have no common PLMN
   */
  if (ta_ret != TA_LIST_RET_OK) {
    OAILOG_ERROR(
        LOG_S1AP, "No Common PLMN with eNB, generate_s1_setup_failure\n");
    rc = s1ap_mme_generate_s1_setup_failure(
        assoc_id, S1ap_Cause_PR_misc, S1ap_CauseMisc_unknown_PLMN,
        S1ap_TimeToWait_v20s);

    increment_counter(
        "s1_setup", 1, 2, "result", "failure", "cause",
        "plmnid_or_tac_mismatch");
    OAILOG_FUNC_RETURN(LOG_S1AP, rc);
  }

  S1ap_SupportedTAs_t* ta_list = &ie_supported_tas->value.choice.SupportedTAs;
  supported_ta_list_t* supp_ta_list = &enb_association->supported_ta_list;
  supp_ta_list->list_count          = ta_list->list.count;

  /* Storing supported TAI lists received in S1 SETUP REQUEST message */
  for (int tai_idx = 0; tai_idx < supp_ta_list->list_count; tai_idx++) {
    S1ap_SupportedTAs_Item_t* tai = NULL;
    tai                           = ta_list->list.array[tai_idx];
    OCTET_STRING_TO_TAC(
        &tai->tAC, supp_ta_list->supported_tai_items[tai_idx].tac);

    bplmn_list_count = tai->broadcastPLMNs.list.count;
    if (bplmn_list_count > S1AP_MAX_BROADCAST_PLMNS) {
      OAILOG_ERROR(
          LOG_S1AP, "Maximum Broadcast PLMN list count exceeded, count = %d\n",
          bplmn_list_count);
    }
    supp_ta_list->supported_tai_items[tai_idx].bplmnlist_count =
        bplmn_list_count;
    for (int plmn_idx = 0; plmn_idx < bplmn_list_count; plmn_idx++) {
      TBCD_TO_PLMN_T(
          tai->broadcastPLMNs.list.array[plmn_idx],
          &supp_ta_list->supported_tai_items[tai_idx].bplmns[plmn_idx]);
    }
  }
  OAILOG_DEBUG(
      LOG_S1AP, "Adding eNB with enb_id :%d to the list of served eNBs \n",
      enb_id);

  enb_association->enb_id = enb_id;

  S1AP_FIND_PROTOCOLIE_BY_ID(
      S1ap_S1SetupRequestIEs_t, ie_default_paging_drx, container,
      S1ap_ProtocolIE_ID_id_DefaultPagingDRX, true);

  enb_association->default_paging_drx =
      ie_default_paging_drx->value.choice.PagingDRX;

  if (enb_name != NULL) {
    memcpy(
        enb_association->enb_name, ie_enb_name->value.choice.ENBname.buf,
        ie_enb_name->value.choice.ENBname.size);
    enb_association->enb_name[ie_enb_name->value.choice.ENBname.size] = '\0';
  }

  // Clean any stale eNB association (from Redis) for this enb_id
  clean_stale_enb_state(state, enb_association);

  s1ap_dump_enb(enb_association);
  rc = s1ap_generate_s1_setup_response(state, enb_association);
  if (rc == RETURNok) {
    state->num_enbs++;
    set_gauge("s1_connection", 1, 1, "enb_name", enb_association->enb_name);
    increment_counter("s1_setup", 1, 1, "result", "success");
    s1_setup_success_event(enb_name, enb_id);
  }
  OAILOG_FUNC_RETURN(LOG_S1AP, rc);
}

//------------------------------------------------------------------------------
status_code_e s1ap_generate_s1_setup_response(
    s1ap_state_t* state, enb_description_t* enb_association) {
  S1ap_S1AP_PDU_t pdu;
  S1ap_S1SetupResponse_t* out;
  S1ap_S1SetupResponseIEs_t* ie          = NULL;
  S1ap_ServedGUMMEIsItem_t* servedGUMMEI = NULL;
  int i, j;
  int enc_rval    = 0;
  uint8_t* buffer = NULL;
  uint32_t length = 0;
  int rc          = RETURNok;

  OAILOG_FUNC_IN(LOG_S1AP);
  if (enb_association == NULL) {
    OAILOG_ERROR(LOG_S1AP, "enb_association is NULL\n");
    return RETURNerror;
  }
  memset(&pdu, 0, sizeof(pdu));
  pdu.present = S1ap_S1AP_PDU_PR_successfulOutcome;
  pdu.choice.successfulOutcome.procedureCode = S1ap_ProcedureCode_id_S1Setup;
  pdu.choice.successfulOutcome.criticality   = S1ap_Criticality_reject;
  pdu.choice.successfulOutcome.value.present =
      S1ap_SuccessfulOutcome__value_PR_S1SetupResponse;
  out = &pdu.choice.successfulOutcome.value.choice.S1SetupResponse;

  // Generating response
  ie =
      (S1ap_S1SetupResponseIEs_t*) calloc(1, sizeof(S1ap_S1SetupResponseIEs_t));
  ie->id            = S1ap_ProtocolIE_ID_id_ServedGUMMEIs;
  ie->criticality   = S1ap_Criticality_reject;
  ie->value.present = S1ap_S1SetupResponseIEs__value_PR_ServedGUMMEIs;

  // memset for gcc 4.8.4 instead of {0}, servedGUMMEI.servedPLMNs
  servedGUMMEI = calloc(1, sizeof *servedGUMMEI);

  mme_config_read_lock(&mme_config);
  /*
   * Use the gummei parameters provided by configuration
   * that should be sorted
   */
  for (i = 0; i < mme_config.served_tai.nb_tai; i++) {
    bool plmn_added = false;
    for (j = 0; j < i; j++) {
      if ((mme_config.served_tai.plmn_mcc[j] ==
           mme_config.served_tai.plmn_mcc[i]) &&
          (mme_config.served_tai.plmn_mnc[j] ==
           mme_config.served_tai.plmn_mnc[i]) &&
          (mme_config.served_tai.plmn_mnc_len[j] ==
           mme_config.served_tai.plmn_mnc_len[i])) {
        plmn_added = true;
        break;
      }
    }
    if (false == plmn_added) {
      S1ap_PLMNidentity_t* plmn = NULL;
      plmn                      = calloc(1, sizeof(*plmn));
      MCC_MNC_TO_PLMNID(
          mme_config.served_tai.plmn_mcc[i], mme_config.served_tai.plmn_mnc[i],
          mme_config.served_tai.plmn_mnc_len[i], plmn);
      ASN_SEQUENCE_ADD(&servedGUMMEI->servedPLMNs.list, plmn);
    }
  }

  for (i = 0; i < mme_config.gummei.nb; i++) {
    S1ap_MME_Group_ID_t* mme_gid = NULL;
    S1ap_MME_Code_t* mmec        = NULL;

    mme_gid = calloc(1, sizeof(*mme_gid));
    INT16_TO_OCTET_STRING(mme_config.gummei.gummei[i].mme_gid, mme_gid);
    ASN_SEQUENCE_ADD(&servedGUMMEI->servedGroupIDs.list, mme_gid);

    mmec = calloc(1, sizeof(*mmec));
    INT8_TO_OCTET_STRING(mme_config.gummei.gummei[i].mme_code, mmec);
    ASN_SEQUENCE_ADD(&servedGUMMEI->servedMMECs.list, mmec);
  }
  ASN_SEQUENCE_ADD(&ie->value.choice.ServedGUMMEIs.list, servedGUMMEI);
  ASN_SEQUENCE_ADD(&out->protocolIEs.list, ie);

  ie =
      (S1ap_S1SetupResponseIEs_t*) calloc(1, sizeof(S1ap_S1SetupResponseIEs_t));
  ie->id            = S1ap_ProtocolIE_ID_id_RelativeMMECapacity;
  ie->criticality   = S1ap_Criticality_ignore;
  ie->value.present = S1ap_S1SetupResponseIEs__value_PR_RelativeMMECapacity;
  ie->value.choice.RelativeMMECapacity = mme_config.relative_capacity;
  ASN_SEQUENCE_ADD(&out->protocolIEs.list, ie);

  mme_config_unlock(&mme_config);
  /*
   * The MME is only serving E-UTRAN RAT, so the list contains only one element
   */
  enc_rval = s1ap_mme_encode_pdu(&pdu, &buffer, &length);

  /*
   * Failed to encode s1 setup response...
   */
  if (enc_rval < 0) {
    OAILOG_DEBUG(LOG_S1AP, "Removed eNB %d\n", enb_association->sctp_assoc_id);
    s1ap_remove_enb(state, enb_association);
  } else {
    /*
     * Consider the response as sent. S1AP is ready to accept UE contexts
     */
    enb_association->s1_state = S1AP_READY;
  }

  /*
   * Non-UE signalling -> stream 0
   */
  bstring b = blk2bstr(buffer, length);
  free(buffer);
  rc = s1ap_mme_itti_send_sctp_request(
      &b, enb_association->sctp_assoc_id, 0, INVALID_MME_UE_S1AP_ID);

  OAILOG_FUNC_RETURN(LOG_S1AP, rc);
}

//------------------------------------------------------------------------------
status_code_e s1ap_mme_handle_ue_cap_indication(
    s1ap_state_t* state, __attribute__((unused)) const sctp_assoc_id_t assoc_id,
    const sctp_stream_id_t stream, S1ap_S1AP_PDU_t* pdu) {
  ue_description_t* ue_ref_p = NULL;
  S1ap_UECapabilityInfoIndication_t* container;
  S1ap_UECapabilityInfoIndicationIEs_t* ie = NULL;
  int rc                                   = RETURNok;
  mme_ue_s1ap_id_t mme_ue_s1ap_id          = INVALID_MME_UE_S1AP_ID;
  enb_ue_s1ap_id_t enb_ue_s1ap_id          = INVALID_ENB_UE_S1AP_ID;
  imsi64_t imsi64                          = INVALID_IMSI64;

  OAILOG_FUNC_IN(LOG_S1AP);
  if (pdu == NULL) {
    OAILOG_DEBUG(LOG_S1AP, "PDU is NULL\n");
    return RETURNerror;
  }

  container =
      &pdu->choice.initiatingMessage.value.choice.UECapabilityInfoIndication;

  S1AP_FIND_PROTOCOLIE_BY_ID(
      S1ap_UECapabilityInfoIndicationIEs_t, ie, container,
      S1ap_ProtocolIE_ID_id_MME_UE_S1AP_ID, true);

  if (ie) {
    mme_ue_s1ap_id = ie->value.choice.MME_UE_S1AP_ID;
    if ((ue_ref_p = s1ap_state_get_ue_mmeid(mme_ue_s1ap_id)) == NULL) {
      OAILOG_DEBUG(
          LOG_S1AP,
          "No UE is attached to this mme UE s1ap id: " MME_UE_S1AP_ID_FMT "\n",
          mme_ue_s1ap_id);
      OAILOG_FUNC_RETURN(LOG_S1AP, RETURNerror);
    }
  } else {
    OAILOG_FUNC_RETURN(LOG_S1AP, RETURNerror);
  }

  s1ap_imsi_map_t* s1ap_imsi_map = get_s1ap_imsi_map();
  hashtable_uint64_ts_get(
      s1ap_imsi_map->mme_ue_id_imsi_htbl, (const hash_key_t) mme_ue_s1ap_id,
      &imsi64);

  S1AP_FIND_PROTOCOLIE_BY_ID(
      S1ap_UECapabilityInfoIndicationIEs_t, ie, container,
      S1ap_ProtocolIE_ID_id_eNB_UE_S1AP_ID, true);

  if (ie) {
    enb_ue_s1ap_id = (enb_ue_s1ap_id_t)(
        ie->value.choice.ENB_UE_S1AP_ID & ENB_UE_S1AP_ID_MASK);
  } else {
    OAILOG_FUNC_RETURN(LOG_S1AP, RETURNerror);
  }

  if (ue_ref_p->enb_ue_s1ap_id != enb_ue_s1ap_id) {
    OAILOG_DEBUG_UE(
        LOG_S1AP, imsi64,
        "Mismatch in eNB UE S1AP ID, known: " ENB_UE_S1AP_ID_FMT
        ", received: " ENB_UE_S1AP_ID_FMT "\n",
        ue_ref_p->enb_ue_s1ap_id, (uint32_t) enb_ue_s1ap_id);
    OAILOG_FUNC_RETURN(LOG_S1AP, RETURNerror);
  }

  /*
   * Just display a warning when message received over wrong stream
   */
  if (ue_ref_p->sctp_stream_recv != stream) {
    OAILOG_ERROR_UE(
        LOG_S1AP, imsi64,
        "Received ue capability indication for "
        "(MME UE S1AP ID/eNB UE S1AP ID) (" MME_UE_S1AP_ID_FMT
        "/" ENB_UE_S1AP_ID_FMT
        ") over wrong stream "
        "expecting %u, received on %u\n",
        (uint32_t) mme_ue_s1ap_id, ue_ref_p->enb_ue_s1ap_id,
        ue_ref_p->sctp_stream_recv, stream);
  }

  /*
   * Forward the ue capabilities to MME application layer
   */
  S1AP_FIND_PROTOCOLIE_BY_ID(
      S1ap_UECapabilityInfoIndicationIEs_t, ie, container,
      S1ap_ProtocolIE_ID_id_UERadioCapability, true);

  if (ie) {
    MessageDef* message_p                = NULL;
    itti_s1ap_ue_cap_ind_t* ue_cap_ind_p = NULL;

    message_p = itti_alloc_new_message(TASK_S1AP, S1AP_UE_CAPABILITIES_IND);

    if (message_p == NULL) {
      OAILOG_ERROR(LOG_S1AP, "message_p is NULL\n");
      return RETURNerror;
    }
    ue_cap_ind_p                 = &message_p->ittiMsg.s1ap_ue_cap_ind;
    ue_cap_ind_p->enb_ue_s1ap_id = ue_ref_p->enb_ue_s1ap_id;
    ue_cap_ind_p->mme_ue_s1ap_id = ue_ref_p->mme_ue_s1ap_id;
    ue_cap_ind_p->radio_capabilities_length =
        ie->value.choice.UERadioCapability.size;
    ue_cap_ind_p->radio_capabilities = calloc(
        ue_cap_ind_p->radio_capabilities_length,
        sizeof(*ue_cap_ind_p->radio_capabilities));
    memcpy(
        ue_cap_ind_p->radio_capabilities,
        ie->value.choice.UERadioCapability.buf,
        ue_cap_ind_p->radio_capabilities_length);

    message_p->ittiMsgHeader.imsi = imsi64;
    rc = send_msg_to_task(&s1ap_task_zmq_ctx, TASK_MME_APP, message_p);
    OAILOG_FUNC_RETURN(LOG_S1AP, rc);
  } else {
    OAILOG_FUNC_RETURN(LOG_S1AP, RETURNerror);
  }
  OAILOG_FUNC_RETURN(LOG_S1AP, RETURNok);
}

////////////////////////////////////////////////////////////////////////////////
//******************* Context Management procedures **************************//
////////////////////////////////////////////////////////////////////////////////

//------------------------------------------------------------------------------
status_code_e s1ap_mme_handle_initial_context_setup_response(
    s1ap_state_t* state, __attribute__((unused)) const sctp_assoc_id_t assoc_id,
    __attribute__((unused)) const sctp_stream_id_t stream,
    S1ap_S1AP_PDU_t* pdu) {
  S1ap_InitialContextSetupResponse_t* container;
  S1ap_InitialContextSetupResponseIEs_t* ie                   = NULL;
  S1ap_E_RABSetupItemCtxtSUResIEs_t* eRABSetupItemCtxtSURes_p = NULL;
  ue_description_t* ue_ref_p                                  = NULL;
  MessageDef* message_p                                       = NULL;
  int rc                                                      = RETURNok;
  mme_ue_s1ap_id_t mme_ue_s1ap_id = INVALID_MME_UE_S1AP_ID;
  enb_ue_s1ap_id_t enb_ue_s1ap_id = INVALID_ENB_UE_S1AP_ID;
  imsi64_t imsi64;

  OAILOG_FUNC_IN(LOG_S1AP);
  container =
      &pdu->choice.successfulOutcome.value.choice.InitialContextSetupResponse;
  S1AP_FIND_PROTOCOLIE_BY_ID(
      S1ap_InitialContextSetupResponseIEs_t, ie, container,
      S1ap_ProtocolIE_ID_id_MME_UE_S1AP_ID, true);

  if (ie) {
    mme_ue_s1ap_id = ie->value.choice.MME_UE_S1AP_ID;
    if ((ue_ref_p = s1ap_state_get_ue_mmeid((uint32_t) mme_ue_s1ap_id)) ==
        NULL) {
      OAILOG_DEBUG(
          LOG_S1AP,
          "No UE is attached to this mme UE s1ap id: " MME_UE_S1AP_ID_FMT
          " %u(10)\n",
          (uint32_t) mme_ue_s1ap_id, (uint32_t) mme_ue_s1ap_id);
      OAILOG_FUNC_RETURN(LOG_S1AP, RETURNerror);
    }
  } else {
    OAILOG_FUNC_RETURN(LOG_S1AP, RETURNerror);
  }
  s1ap_imsi_map_t* s1ap_imsi_map = get_s1ap_imsi_map();
  hashtable_uint64_ts_get(
      s1ap_imsi_map->mme_ue_id_imsi_htbl, (const hash_key_t) mme_ue_s1ap_id,
      &imsi64);

  S1AP_FIND_PROTOCOLIE_BY_ID(
      S1ap_InitialContextSetupResponseIEs_t, ie, container,
      S1ap_ProtocolIE_ID_id_eNB_UE_S1AP_ID, true);
  if (ie) {
    enb_ue_s1ap_id = (enb_ue_s1ap_id_t)(
        ie->value.choice.ENB_UE_S1AP_ID & ENB_UE_S1AP_ID_MASK);
  } else {
    OAILOG_FUNC_RETURN(LOG_S1AP, RETURNerror);
  }

  if (ue_ref_p->enb_ue_s1ap_id != enb_ue_s1ap_id) {
    OAILOG_DEBUG_UE(
        LOG_S1AP, imsi64,
        "Mismatch in eNB UE S1AP ID, known: " ENB_UE_S1AP_ID_FMT
        " %u(10), received: 0x%06x %u(10)\n",
        ue_ref_p->enb_ue_s1ap_id, ue_ref_p->enb_ue_s1ap_id,
        (uint32_t) enb_ue_s1ap_id, (uint32_t) enb_ue_s1ap_id);
    OAILOG_FUNC_RETURN(LOG_S1AP, RETURNerror);
  }

  S1AP_FIND_PROTOCOLIE_BY_ID(
      S1ap_InitialContextSetupResponseIEs_t, ie, container,
      S1ap_ProtocolIE_ID_id_E_RABSetupListCtxtSURes, true);

  if (ie) {
    if (ie->value.choice.E_RABSetupListCtxtSURes.list.count < 1) {
      OAILOG_WARNING_UE(LOG_S1AP, imsi64, "E-RAB creation has failed\n");
      OAILOG_FUNC_RETURN(LOG_S1AP, RETURNerror);
    }
  } else {
    OAILOG_FUNC_RETURN(LOG_S1AP, RETURNerror);
  }
  ue_ref_p->s1_ue_state = S1AP_UE_CONNECTED;
  message_p             = DEPRECATEDitti_alloc_new_message_fatal(
      TASK_S1AP, MME_APP_INITIAL_CONTEXT_SETUP_RSP);
  MME_APP_INITIAL_CONTEXT_SETUP_RSP(message_p).ue_id = ue_ref_p->mme_ue_s1ap_id;
  MME_APP_INITIAL_CONTEXT_SETUP_RSP(message_p).e_rab_setup_list.no_of_items =
      ie->value.choice.E_RABSetupListCtxtSURes.list.count;
  for (int item = 0; item < ie->value.choice.E_RABSetupListCtxtSURes.list.count;
       item++) {
    /*
     * Bad, very bad cast...
     */
    eRABSetupItemCtxtSURes_p =
        (S1ap_E_RABSetupItemCtxtSUResIEs_t*)
            ie->value.choice.E_RABSetupListCtxtSURes.list.array[item];
    MME_APP_INITIAL_CONTEXT_SETUP_RSP(message_p)
        .e_rab_setup_list.item[item]
        .e_rab_id =
        eRABSetupItemCtxtSURes_p->value.choice.E_RABSetupItemCtxtSURes.e_RAB_ID;
    MME_APP_INITIAL_CONTEXT_SETUP_RSP(message_p)
        .e_rab_setup_list.item[item]
        .gtp_teid = htonl(*((uint32_t*) eRABSetupItemCtxtSURes_p->value.choice
                                .E_RABSetupItemCtxtSURes.gTP_TEID.buf));

    // When Magma AGW runs in a cloud and RAN at the edge (hence eNB is
    // behind NAT), eNB signals its private IP address in ICS Response. By
    // setting "enable_gtpu_private_ip_correction" true in mme.yml file,
    // we can correct that private IP address with the public IP address of
    // the eNB as its public IP address is observed during SCTP link is set up.
    // We store this "control plane IP address" as part of the eNB context
    // information. This feature can be safely used only when NAT uses the same
    // public IP address for both the CP and UP communication to/from the eNB,
    // which typically is the situation.
    enb_description_t* enb_association = s1ap_state_get_enb(state, assoc_id);
    if (mme_config.enable_gtpu_private_ip_correction) {
      OAILOG_INFO(
          LOG_S1AP,
          "Overwriting eNB GTP-U IP ADDRESS with SCTP eNB IP address");
      MME_APP_INITIAL_CONTEXT_SETUP_RSP(message_p)
          .e_rab_setup_list.item[item]
          .transport_layer_address = blk2bstr(
          enb_association->ran_cp_ipaddr, enb_association->ran_cp_ipaddr_sz);
    } else {
      // Print a warning message if CP and UP plane eNB IPs are different
      if (memcmp(
              enb_association->ran_cp_ipaddr,
              eRABSetupItemCtxtSURes_p->value.choice.E_RABSetupItemCtxtSURes
                  .transportLayerAddress.buf,
              enb_association->ran_cp_ipaddr_sz)) {
        OAILOG_WARNING(
            LOG_S1AP,
            "GTP-U eNB IP addr is different than SCTP eNB IP addr. "
            "This can be due to eNB behind a NAT. Consider setting "
            "enable_gtpu_private_ip_correction as true in mme.yml file.");
      }
      MME_APP_INITIAL_CONTEXT_SETUP_RSP(message_p)
          .e_rab_setup_list.item[item]
          .transport_layer_address = blk2bstr(
          eRABSetupItemCtxtSURes_p->value.choice.E_RABSetupItemCtxtSURes
              .transportLayerAddress.buf,
          eRABSetupItemCtxtSURes_p->value.choice.E_RABSetupItemCtxtSURes
              .transportLayerAddress.size);
    }
  }

  // Failed bearers
  itti_mme_app_initial_context_setup_rsp_t* initial_context_setup_rsp =
      &(MME_APP_INITIAL_CONTEXT_SETUP_RSP(message_p));
  initial_context_setup_rsp->e_rab_failed_to_setup_list.no_of_items = 0;
  S1AP_FIND_PROTOCOLIE_BY_ID(
      S1ap_InitialContextSetupResponseIEs_t, ie, container,
      S1ap_ProtocolIE_ID_id_E_RABFailedToSetupListBearerSURes, false);
  if (ie) {
    S1ap_E_RABList_t* s1ap_e_rab_list = &ie->value.choice.E_RABList;
    for (int index = 0; index < s1ap_e_rab_list->list.count; index++) {
      S1ap_E_RABItem_t* erab_item =
          (S1ap_E_RABItem_t*) s1ap_e_rab_list->list.array[index];
      initial_context_setup_rsp->e_rab_failed_to_setup_list.item[index]
          .e_rab_id = erab_item->e_RAB_ID;
      initial_context_setup_rsp->e_rab_failed_to_setup_list.item[index].cause =
          erab_item->cause;
    }
    initial_context_setup_rsp->e_rab_failed_to_setup_list.no_of_items =
        s1ap_e_rab_list->list.count;
  }
  message_p->ittiMsgHeader.imsi = imsi64;
  rc = send_msg_to_task(&s1ap_task_zmq_ctx, TASK_MME_APP, message_p);
  OAILOG_FUNC_RETURN(LOG_S1AP, rc);
}

//------------------------------------------------------------------------------
status_code_e s1ap_mme_handle_ue_context_release_request(
    s1ap_state_t* state, __attribute__((unused)) const sctp_assoc_id_t assoc_id,
    __attribute__((unused)) const sctp_stream_id_t stream,
    S1ap_S1AP_PDU_t* pdu) {
  S1ap_UEContextReleaseRequest_t* container;
  S1ap_UEContextReleaseRequest_IEs_t* ie = NULL;
  ue_description_t* ue_ref_p             = NULL;
  enb_description_t* enb_ref_p           = NULL;
  S1ap_Cause_PR cause_type;
  long cause_value;
  enum s1cause s1_release_cause   = S1AP_RADIO_EUTRAN_GENERATED_REASON;
  int rc                          = RETURNok;
  mme_ue_s1ap_id_t mme_ue_s1ap_id = INVALID_MME_UE_S1AP_ID;
  enb_ue_s1ap_id_t enb_ue_s1ap_id = INVALID_ENB_UE_S1AP_ID;
  imsi64_t imsi64                 = INVALID_IMSI64;

  OAILOG_FUNC_IN(LOG_S1AP);
  if ((enb_ref_p = s1ap_state_get_enb(state, assoc_id)) == NULL) {
    OAILOG_ERROR(
        LOG_S1AP, "Ignoring context release request from unknown assoc %u",
        assoc_id);
    OAILOG_FUNC_RETURN(LOG_S1AP, RETURNerror);
  }

  container =
      &pdu->choice.initiatingMessage.value.choice.UEContextReleaseRequest;
  // Log the Cause Type and Cause value
  S1AP_FIND_PROTOCOLIE_BY_ID(
      S1ap_UEContextReleaseRequest_IEs_t, ie, container,
      S1ap_ProtocolIE_ID_id_MME_UE_S1AP_ID, true);
  if (ie) {
    mme_ue_s1ap_id = ie->value.choice.MME_UE_S1AP_ID;
  } else {
    OAILOG_FUNC_RETURN(LOG_S1AP, RETURNerror);
  }

  S1AP_FIND_PROTOCOLIE_BY_ID(
      S1ap_UEContextReleaseRequest_IEs_t, ie, container,
      S1ap_ProtocolIE_ID_id_eNB_UE_S1AP_ID, true);
  if (ie) {
    enb_ue_s1ap_id = (enb_ue_s1ap_id_t)(
        ie->value.choice.ENB_UE_S1AP_ID & ENB_UE_S1AP_ID_MASK);
  } else {
    OAILOG_FUNC_RETURN(LOG_S1AP, RETURNerror);
  }

  // Log the Cause Type and Cause value
  S1AP_FIND_PROTOCOLIE_BY_ID(
      S1ap_UEContextReleaseRequest_IEs_t, ie, container,
      S1ap_ProtocolIE_ID_id_Cause, true);
  if (ie) {
    cause_type = ie->value.choice.Cause.present;
  } else {
    OAILOG_FUNC_RETURN(LOG_S1AP, RETURNerror);
  }
  switch (cause_type) {
    case S1ap_Cause_PR_radioNetwork:
      cause_value = ie->value.choice.Cause.choice.radioNetwork;
      OAILOG_INFO(
          LOG_S1AP,
          "UE CONTEXT RELEASE REQUEST with Cause_Type = Radio Network and "
          "Cause_Value = %ld\n",
          cause_value);
      if (cause_value == S1ap_CauseRadioNetwork_user_inactivity) {
        increment_counter(
            "ue_context_release_req", 1, 1, "cause", "user_inactivity");
      } else if (
          cause_value == S1ap_CauseRadioNetwork_radio_connection_with_ue_lost) {
        increment_counter(
            "ue_context_release_req", 1, 1, "cause", "radio_link_failure");
      } else if (
          cause_value ==
          S1ap_CauseRadioNetwork_ue_not_available_for_ps_service) {
        increment_counter(
            "ue_context_release_req", 1, 1, "cause",
            "ue_not_available_for_ps_service");
        s1_release_cause = S1AP_NAS_UE_NOT_AVAILABLE_FOR_PS;
      } else if (cause_value == S1ap_CauseRadioNetwork_cs_fallback_triggered) {
        increment_counter(
            "ue_context_release_req", 1, 1, "cause", "cs_fallback_triggered");
        s1_release_cause = S1AP_CSFB_TRIGGERED;
      }
      break;

    case S1ap_Cause_PR_transport:
      cause_value = ie->value.choice.Cause.choice.transport;
      OAILOG_INFO(
          LOG_S1AP,
          "UE CONTEXT RELEASE REQUEST with Cause_Type = Transport and "
          "Cause_Value = %ld\n",
          cause_value);
      break;

    case S1ap_Cause_PR_nas:
      cause_value = ie->value.choice.Cause.choice.nas;
      OAILOG_INFO(
          LOG_S1AP,
          "UE CONTEXT RELEASE REQUEST with Cause_Type = NAS and Cause_Value = "
          "%ld\n",
          cause_value);
      break;

    case S1ap_Cause_PR_protocol:
      cause_value = ie->value.choice.Cause.choice.protocol;
      OAILOG_INFO(
          LOG_S1AP,
          "UE CONTEXT RELEASE REQUEST with Cause_Type = Transport and "
          "Cause_Value = %ld\n",
          cause_value);
      break;

    case S1ap_Cause_PR_misc:
      cause_value = ie->value.choice.Cause.choice.misc;
      OAILOG_INFO(
          LOG_S1AP,
          "UE CONTEXT RELEASE REQUEST with Cause_Type = MISC and Cause_Value = "
          "%ld\n",
          cause_value);
      break;

    default:
      OAILOG_ERROR(
          LOG_S1AP, "UE CONTEXT RELEASE REQUEST with Invalid Cause_Type = %d\n",
          cause_type);
      OAILOG_FUNC_RETURN(LOG_S1AP, RETURNerror);
  }

  /* Fix - MME shall handle UE Context Release received from the eNB
  irrespective of the cause. And MME should release the S1-U bearers for the UE
  and move the UE to ECM idle mode. Cause can influence whether to preserve GBR
  bearers or not.Since, as of now EPC doesn't support dedicated bearers, it is
  don't care scenario till we add support for dedicated bearers.
  */

  if ((ue_ref_p = s1ap_state_get_ue_mmeid(mme_ue_s1ap_id)) == NULL) {
    /*
     * MME doesn't know the MME UE S1AP ID provided.
     * No need to do anything. Ignore the message
     */
    OAILOG_DEBUG(
        LOG_S1AP,
        "UE_CONTEXT_RELEASE_REQUEST ignored cause could not get context with "
        "mme_ue_s1ap_id " MME_UE_S1AP_ID_FMT
        " enb_ue_s1ap_id " ENB_UE_S1AP_ID_FMT " ",
        (uint32_t) mme_ue_s1ap_id, (uint32_t) enb_ue_s1ap_id);
    OAILOG_FUNC_RETURN(LOG_S1AP, RETURNerror);
  } else {
    s1ap_imsi_map_t* imsi_map = get_s1ap_imsi_map();
    hashtable_uint64_ts_get(
        imsi_map->mme_ue_id_imsi_htbl, (const hash_key_t) mme_ue_s1ap_id,
        &imsi64);
    if (ue_ref_p->sctp_assoc_id == assoc_id &&
        ue_ref_p->enb_ue_s1ap_id == enb_ue_s1ap_id) {
      /*
       * Both eNB UE S1AP ID and MME UE S1AP ID match.
       * Send a UE context Release Command to eNB after releasing S1-U bearer
       * tunnel mapping for all the bearers.
       */
      rc = s1ap_send_mme_ue_context_release(
          state, ue_ref_p, s1_release_cause, ie->value.choice.Cause, imsi64);

      OAILOG_FUNC_RETURN(LOG_S1AP, rc);
    } else if (
        enb_ref_p->enb_id == ue_ref_p->s1ap_handover_state.source_enb_id &&
        ue_ref_p->s1ap_handover_state.source_enb_ue_s1ap_id == enb_ue_s1ap_id) {
      /*
       * We just handed over from this eNB.
       * Send a UE context Release Command to eNB. S1-U bearer already released
       */
      rc = s1ap_mme_generate_ue_context_release_command(
          state, ue_ref_p, S1AP_RADIO_EUTRAN_GENERATED_REASON, imsi64, assoc_id,
          ue_ref_p->s1ap_handover_state.source_sctp_stream_send, mme_ue_s1ap_id,
          enb_ue_s1ap_id);
      OAILOG_FUNC_RETURN(LOG_S1AP, rc);
    } else {
      // abnormal case. No need to do anything. Ignore the message
      OAILOG_DEBUG_UE(
          LOG_S1AP, imsi64,
          "UE_CONTEXT_RELEASE_REQUEST ignored, cause mismatch enb_ue_s1ap_id: "
          "ctxt " ENB_UE_S1AP_ID_FMT " != request " ENB_UE_S1AP_ID_FMT " ",
          (uint32_t) ue_ref_p->enb_ue_s1ap_id, (uint32_t) enb_ue_s1ap_id);
      OAILOG_FUNC_RETURN(LOG_S1AP, RETURNerror);
    }
  }
  OAILOG_FUNC_RETURN(LOG_S1AP, RETURNok);
}

//------------------------------------------------------------------------------
status_code_e s1ap_mme_generate_ue_context_release_command(
    s1ap_state_t* state, ue_description_t* ue_ref_p, enum s1cause cause,
    imsi64_t imsi64, const sctp_assoc_id_t assoc_id,
    const sctp_stream_id_t stream, mme_ue_s1ap_id_t mme_ue_s1ap_id,
    enb_ue_s1ap_id_t enb_ue_s1ap_id) {
  uint8_t* buffer = NULL;
  uint32_t length = 0;
  S1ap_S1AP_PDU_t pdu;
  S1ap_UEContextReleaseCommand_t* out;
  S1ap_UEContextReleaseCommand_IEs_t* ie = NULL;
  int rc                                 = RETURNok;
  S1ap_Cause_PR cause_type;
  long cause_value;

  OAILOG_FUNC_IN(LOG_S1AP);
  memset(&pdu, 0, sizeof(pdu));
  pdu.present = S1ap_S1AP_PDU_PR_initiatingMessage;
  pdu.choice.initiatingMessage.procedureCode =
      S1ap_ProcedureCode_id_UEContextRelease;
  pdu.choice.initiatingMessage.criticality = S1ap_Criticality_reject;
  pdu.choice.initiatingMessage.value.present =
      S1ap_InitiatingMessage__value_PR_UEContextReleaseCommand;
  out = &pdu.choice.initiatingMessage.value.choice.UEContextReleaseCommand;
  /*
   * Fill in ID pair
   */
  ie = (S1ap_UEContextReleaseCommand_IEs_t*) calloc(
      1, sizeof(S1ap_UEContextReleaseCommand_IEs_t));
  ie->id            = S1ap_ProtocolIE_ID_id_UE_S1AP_IDs;
  ie->criticality   = S1ap_Criticality_reject;
  ie->value.present = S1ap_UEContextReleaseCommand_IEs__value_PR_UE_S1AP_IDs;
  ie->value.choice.UE_S1AP_IDs.present = S1ap_UE_S1AP_IDs_PR_uE_S1AP_ID_pair;
  ie->value.choice.UE_S1AP_IDs.choice.uE_S1AP_ID_pair.mME_UE_S1AP_ID =
      mme_ue_s1ap_id;
  ie->value.choice.UE_S1AP_IDs.choice.uE_S1AP_ID_pair.eNB_UE_S1AP_ID =
      enb_ue_s1ap_id;
  ie->value.choice.UE_S1AP_IDs.choice.uE_S1AP_ID_pair.iE_Extensions = NULL;
  ASN_SEQUENCE_ADD(&out->protocolIEs.list, ie);

  ie = (S1ap_UEContextReleaseCommand_IEs_t*) calloc(
      1, sizeof(S1ap_UEContextReleaseCommand_IEs_t));
  ie->id            = S1ap_ProtocolIE_ID_id_Cause;
  ie->criticality   = S1ap_Criticality_ignore;
  ie->value.present = S1ap_UEContextReleaseCommand_IEs__value_PR_Cause;
  switch (cause) {
    case S1AP_NAS_DETACH:
      cause_type  = S1ap_Cause_PR_nas;
      cause_value = S1ap_CauseNas_detach;
      break;
    case S1AP_NAS_NORMAL_RELEASE:
      cause_type  = S1ap_Cause_PR_nas;
      cause_value = S1ap_CauseNas_unspecified;
      break;
    case S1AP_RADIO_EUTRAN_GENERATED_REASON:
      cause_type = S1ap_Cause_PR_radioNetwork;
      cause_value =
          S1ap_CauseRadioNetwork_release_due_to_eutran_generated_reason;
      break;
    case S1AP_INITIAL_CONTEXT_SETUP_FAILED:
      cause_type  = S1ap_Cause_PR_radioNetwork;
      cause_value = S1ap_CauseRadioNetwork_unspecified;
      break;
    case S1AP_CSFB_TRIGGERED:
      cause_type  = S1ap_Cause_PR_radioNetwork;
      cause_value = S1ap_CauseRadioNetwork_cs_fallback_triggered;
      break;
    case S1AP_NAS_UE_NOT_AVAILABLE_FOR_PS:
      cause_type  = S1ap_Cause_PR_radioNetwork;
      cause_value = S1ap_CauseRadioNetwork_ue_not_available_for_ps_service;
      break;
    case S1AP_RADIO_MULTIPLE_E_RAB_ID:
      cause_type  = S1ap_Cause_PR_radioNetwork;
      cause_value = S1ap_CauseRadioNetwork_multiple_E_RAB_ID_instances;
    case S1AP_INVALID_MME_UE_S1AP_ID:
      cause_type  = S1ap_Cause_PR_radioNetwork;
      cause_value = S1ap_CauseRadioNetwork_unknown_mme_ue_s1ap_id;
      break;
    case S1AP_NAS_MME_OFFLOADING:
      cause_type  = S1ap_Cause_PR_radioNetwork;
      cause_value = S1ap_CauseRadioNetwork_load_balancing_tau_required;
      break;
    default:
      // Freeing ie and pdu data since it will not be encoded
      free_wrapper((void**) &ie);
      ASN_STRUCT_FREE_CONTENTS_ONLY(asn_DEF_S1ap_S1AP_PDU, &pdu);
      OAILOG_ERROR_UE(LOG_S1AP, imsi64, "Unknown cause for context release");
      OAILOG_FUNC_RETURN(LOG_S1AP, RETURNerror);
  }
  s1ap_mme_set_cause(&ie->value.choice.Cause, cause_type, cause_value);
  ASN_SEQUENCE_ADD(&out->protocolIEs.list, ie);

  if (s1ap_mme_encode_pdu(&pdu, &buffer, &length) < 0) {
    OAILOG_FUNC_RETURN(LOG_S1AP, RETURNerror);
  }

  bstring b = blk2bstr(buffer, length);
  free(buffer);
  rc = s1ap_mme_itti_send_sctp_request(&b, assoc_id, stream, mme_ue_s1ap_id);
  if (ue_ref_p != NULL && ue_ref_p->sctp_assoc_id == assoc_id) {
    ue_ref_p->s1_ue_state = S1AP_UE_WAITING_CRR;
    // We can safely remove UE context now, no need for timer
    s1ap_mme_release_ue_context(state, ue_ref_p, imsi64);
  }
  OAILOG_FUNC_RETURN(LOG_S1AP, rc);
}

//------------------------------------------------------------------------------

//------------------------------------------------------------------------------
status_code_e s1ap_mme_generate_ue_context_modification(
    ue_description_t* ue_ref_p,
    const itti_s1ap_ue_context_mod_req_t* const ue_context_mod_req_pP,
    imsi64_t imsi64) {
  uint8_t* buffer                                = NULL;
  uint32_t length                                = 0;
  S1ap_S1AP_PDU_t pdu                            = {0};
  S1ap_UEContextModificationRequest_t* container = NULL;
  S1ap_UEContextModificationRequestIEs_t* ie     = NULL;
  int rc                                         = RETURNok;

  OAILOG_FUNC_IN(LOG_S1AP);
  if (ue_ref_p == NULL) {
    OAILOG_FUNC_RETURN(LOG_S1AP, RETURNerror);
  }
  pdu.present = S1ap_S1AP_PDU_PR_initiatingMessage;
  pdu.choice.initiatingMessage.procedureCode =
      S1ap_ProcedureCode_id_UEContextModification;
  pdu.choice.initiatingMessage.criticality = S1ap_Criticality_reject;
  pdu.choice.initiatingMessage.value.present =
      S1ap_InitiatingMessage__value_PR_UEContextModificationRequest;
  container =
      &pdu.choice.initiatingMessage.value.choice.UEContextModificationRequest;

  /*
   * Fill in ID pair
   */
  ie = (S1ap_UEContextModificationRequestIEs_t*) calloc(
      1, sizeof(S1ap_UEContextModificationRequestIEs_t));
  ie->id          = S1ap_ProtocolIE_ID_id_UE_S1AP_IDs;
  ie->criticality = S1ap_Criticality_reject;
  ie->value.present =
      S1ap_UEContextModificationRequestIEs__value_PR_MME_UE_S1AP_ID;
  ie->value.choice.MME_UE_S1AP_ID = ue_ref_p->mme_ue_s1ap_id;
  ASN_SEQUENCE_ADD(&container->protocolIEs.list, ie);

  ie = (S1ap_UEContextModificationRequestIEs_t*) calloc(
      1, sizeof(S1ap_UEContextModificationRequestIEs_t));
  ie->id          = S1ap_ProtocolIE_ID_id_UE_S1AP_IDs;
  ie->criticality = S1ap_Criticality_reject;
  ie->value.present =
      S1ap_UEContextModificationRequestIEs__value_PR_ENB_UE_S1AP_ID;
  ie->value.choice.ENB_UE_S1AP_ID = ue_ref_p->enb_ue_s1ap_id;
  ASN_SEQUENCE_ADD(&container->protocolIEs.list, ie);

  if ((ue_context_mod_req_pP->presencemask & S1AP_UE_CONTEXT_MOD_LAI_PRESENT) ==
      S1AP_UE_CONTEXT_MOD_LAI_PRESENT) {
    ie = (S1ap_UEContextModificationRequestIEs_t*) calloc(
        1, sizeof(S1ap_UEContextModificationRequestIEs_t));
    ie->id            = S1ap_ProtocolIE_ID_id_RegisteredLAI;
    ie->criticality   = S1ap_Criticality_reject;
    ie->value.present = S1ap_UEContextModificationRequestIEs__value_PR_LAI;
#define PLMN_SIZE 3
    S1ap_LAI_t* lai_item        = &ie->value.choice.LAI;
    lai_item->pLMNidentity.size = PLMN_SIZE;
    lai_item->pLMNidentity.buf  = calloc(PLMN_SIZE, sizeof(uint8_t));
    uint8_t mnc_length          = mme_config_find_mnc_length(
        ue_context_mod_req_pP->lai.mccdigit1,
        ue_context_mod_req_pP->lai.mccdigit2,
        ue_context_mod_req_pP->lai.mccdigit3,
        ue_context_mod_req_pP->lai.mncdigit1,
        ue_context_mod_req_pP->lai.mncdigit2,
        ue_context_mod_req_pP->lai.mncdigit3);
    if (mnc_length != 2 && mnc_length != 3) {
      free_wrapper((void**) &ie);
      ASN_STRUCT_FREE_CONTENTS_ONLY(
          asn_DEF_S1ap_UEContextModificationRequest, container);
      OAILOG_FUNC_RETURN(LOG_S1AP, RETURNerror);
    }
    LAI_T_TO_TBCD(
        ue_context_mod_req_pP->lai, lai_item->pLMNidentity.buf, mnc_length);

    TAC_TO_ASN1(ue_context_mod_req_pP->lai.lac, &lai_item->lAC);
    lai_item->iE_Extensions = NULL;
    ASN_SEQUENCE_ADD(&container->protocolIEs.list, ie);
  }

  if ((ue_context_mod_req_pP->presencemask &
       S1AP_UE_CONTEXT_MOD_CSFB_INDICATOR_PRESENT) ==
      S1AP_UE_CONTEXT_MOD_CSFB_INDICATOR_PRESENT) {
    ie = (S1ap_UEContextModificationRequestIEs_t*) calloc(
        1, sizeof(S1ap_UEContextModificationRequestIEs_t));
    ie->id          = S1ap_ProtocolIE_ID_id_CSFallbackIndicator;
    ie->criticality = S1ap_Criticality_reject;
    ie->value.present =
        S1ap_UEContextModificationRequestIEs__value_PR_CSFallbackIndicator;
    ie->value.choice.CSFallbackIndicator =
        ue_context_mod_req_pP->cs_fallback_indicator;
    ASN_SEQUENCE_ADD(&container->protocolIEs.list, ie);
  }

  if ((ue_context_mod_req_pP->presencemask &
       S1AP_UE_CONTEXT_MOD_UE_AMBR_INDICATOR_PRESENT) ==
      S1AP_UE_CONTEXT_MOD_UE_AMBR_INDICATOR_PRESENT) {
    ie = (S1ap_UEContextModificationRequestIEs_t*) calloc(
        1, sizeof(S1ap_UEContextModificationRequestIEs_t));
    ie->id          = S1ap_ProtocolIE_ID_id_uEaggregateMaximumBitrate;
    ie->criticality = S1ap_Criticality_ignore;
    ie->value.present =
        S1ap_UEContextModificationRequestIEs__value_PR_UEAggregateMaximumBitrate;
    asn_uint642INTEGER(
        &ie->value.choice.UEAggregateMaximumBitrate.uEaggregateMaximumBitRateDL,
        ue_context_mod_req_pP->ue_ambr.br_dl);
    asn_uint642INTEGER(
        &ie->value.choice.UEAggregateMaximumBitrate.uEaggregateMaximumBitRateUL,
        ue_context_mod_req_pP->ue_ambr.br_ul);
    ASN_SEQUENCE_ADD(&container->protocolIEs.list, ie);
  }

  if (s1ap_mme_encode_pdu(&pdu, &buffer, &length) < 0) {
    OAILOG_FUNC_RETURN(LOG_S1AP, RETURNerror);
  }

  bstring b = blk2bstr(buffer, length);
  free(buffer);
  rc = s1ap_mme_itti_send_sctp_request(
      &b, ue_ref_p->sctp_assoc_id, ue_ref_p->sctp_stream_send,
      ue_ref_p->mme_ue_s1ap_id);

  OAILOG_FUNC_RETURN(LOG_S1AP, rc);
}

//------------------------------------------------------------------------------
status_code_e s1ap_handle_ue_context_release_command(
    s1ap_state_t* state,
    const itti_s1ap_ue_context_release_command_t* const
        ue_context_release_command_pP,
    imsi64_t imsi64) {
  ue_description_t* ue_ref_p = NULL;
  int rc                     = RETURNok;

  OAILOG_FUNC_IN(LOG_S1AP);
  if ((ue_ref_p = s1ap_state_get_ue_mmeid(
           ue_context_release_command_pP->mme_ue_s1ap_id)) == NULL) {
    OAILOG_DEBUG_UE(
        LOG_S1AP, imsi64,
        "Ignoring UE with mme_ue_s1ap_id " MME_UE_S1AP_ID_FMT " %u(10)\n",
        ue_context_release_command_pP->mme_ue_s1ap_id,
        ue_context_release_command_pP->mme_ue_s1ap_id);
    rc = RETURNok;
  } else {
    /*
     * Check the cause. If it is implicit detach or sctp reset/shutdown no need
     * to send UE context release command to eNB. Free UE context locally.
     */

    if (ue_context_release_command_pP->cause == S1AP_IMPLICIT_CONTEXT_RELEASE ||
        ue_context_release_command_pP->cause == S1AP_SCTP_SHUTDOWN_OR_RESET ||
        ue_context_release_command_pP->cause == S1AP_INVALID_ENB_ID) {
      s1ap_remove_ue(state, ue_ref_p);
    } else {
      rc = s1ap_mme_generate_ue_context_release_command(
          state, ue_ref_p, ue_context_release_command_pP->cause, imsi64,
          ue_ref_p->sctp_assoc_id, ue_ref_p->sctp_stream_send,
          ue_ref_p->mme_ue_s1ap_id, ue_ref_p->enb_ue_s1ap_id);
    }
  }

  OAILOG_FUNC_RETURN(LOG_S1AP, rc);
}

//------------------------------------------------------------------------------

//------------------------------------------------------------------------------
status_code_e s1ap_handle_ue_context_mod_req(
    s1ap_state_t* state,
    const itti_s1ap_ue_context_mod_req_t* const ue_context_mod_req_pP,
    imsi64_t imsi64) {
  ue_description_t* ue_ref_p = NULL;
  int rc                     = RETURNok;

  OAILOG_FUNC_IN(LOG_S1AP);
  if (ue_context_mod_req_pP == NULL) {
    OAILOG_ERROR(LOG_S1AP, "ue_context_mod_req_pP is NULL\n");
    return RETURNerror;
  }
  if ((ue_ref_p = s1ap_state_get_ue_mmeid(
           ue_context_mod_req_pP->mme_ue_s1ap_id)) == NULL) {
    OAILOG_DEBUG_UE(
        LOG_S1AP, imsi64,
        "Ignoring UE with mme_ue_s1ap_id " MME_UE_S1AP_ID_FMT " %u(10)\n",
        ue_context_mod_req_pP->mme_ue_s1ap_id,
        ue_context_mod_req_pP->mme_ue_s1ap_id);
    rc = RETURNok;
  } else {
    rc = s1ap_mme_generate_ue_context_modification(
        ue_ref_p, ue_context_mod_req_pP, imsi64);
  }

  OAILOG_FUNC_RETURN(LOG_S1AP, rc);
}

//------------------------------------------------------------------------------
status_code_e s1ap_mme_handle_ue_context_release_complete(
    s1ap_state_t* state, __attribute__((unused)) const sctp_assoc_id_t assoc_id,
    __attribute__((unused)) const sctp_stream_id_t stream,
    S1ap_S1AP_PDU_t* pdu) {
  S1ap_UEContextReleaseComplete_t* container;
  S1ap_UEContextReleaseComplete_IEs_t* ie = NULL;
  ue_description_t* ue_ref_p              = NULL;
  mme_ue_s1ap_id_t mme_ue_s1ap_id         = 0;

  OAILOG_FUNC_IN(LOG_S1AP);
  container =
      &pdu->choice.successfulOutcome.value.choice.UEContextReleaseComplete;

  S1AP_FIND_PROTOCOLIE_BY_ID(
      S1ap_UEContextReleaseComplete_IEs_t, ie, container,
      S1ap_ProtocolIE_ID_id_MME_UE_S1AP_ID, true);

  if (ie) {
    mme_ue_s1ap_id = ie->value.choice.MME_UE_S1AP_ID;
  } else {
    OAILOG_FUNC_RETURN(LOG_S1AP, RETURNok);
  }

  if ((ue_ref_p = s1ap_state_get_ue_mmeid(mme_ue_s1ap_id)) == NULL) {
    /*
     * The UE context has already been deleted when the UE context release
     * command was sent
     * Ignore this message.
     */
    OAILOG_DEBUG(
        LOG_S1AP,
        " UE Context Release commplete: S1 context cleared. Ignore message for "
        "ueid " MME_UE_S1AP_ID_FMT "\n",
        (uint32_t) mme_ue_s1ap_id);
    OAILOG_FUNC_RETURN(LOG_S1AP, RETURNok);
  } else {
    if (ue_ref_p->sctp_assoc_id == assoc_id) {
      /* This is an error scenario, the S1 UE context should have been deleted
       * when UE context release command was sent
       */
      OAILOG_ERROR(
          LOG_S1AP,
          " UE Context Release commplete: S1 context should have been cleared "
          "for ueid " MME_UE_S1AP_ID_FMT "\n",
          (uint32_t) mme_ue_s1ap_id);
      OAILOG_FUNC_RETURN(LOG_S1AP, RETURNerror);
    } else {
      /*
       * UE Context Release commplete received from a different eNB. This could
       * be coming from the source eNB after a successful handover
       */
      OAILOG_DEBUG(
          LOG_S1AP,
          " UE Context Release commplete received from a different eNB."
          " Ignore message for ueid " MME_UE_S1AP_ID_FMT "\n",
          (uint32_t) mme_ue_s1ap_id);
      OAILOG_FUNC_RETURN(LOG_S1AP, RETURNok);
    }
  }
}

//------------------------------------------------------------------------------
status_code_e s1ap_mme_handle_initial_context_setup_failure(
    s1ap_state_t* state, __attribute__((unused)) const sctp_assoc_id_t assoc_id,
    __attribute__((unused)) const sctp_stream_id_t stream,
    S1ap_S1AP_PDU_t* pdu) {
  S1ap_InitialContextSetupFailure_t* container;
  S1ap_InitialContextSetupFailureIEs_t* ie = NULL;
  ue_description_t* ue_ref_p               = NULL;
  MessageDef* message_p                    = NULL;
  S1ap_Cause_PR cause_type;
  long cause_value;
  int rc                          = RETURNok;
  imsi64_t imsi64                 = INVALID_IMSI64;
  mme_ue_s1ap_id_t mme_ue_s1ap_id = INVALID_MME_UE_S1AP_ID;
  enb_ue_s1ap_id_t enb_ue_s1ap_id = INVALID_ENB_UE_S1AP_ID;

  OAILOG_FUNC_IN(LOG_S1AP);
  container =
      &pdu->choice.unsuccessfulOutcome.value.choice.InitialContextSetupFailure;

  S1AP_FIND_PROTOCOLIE_BY_ID(
      S1ap_InitialContextSetupFailureIEs_t, ie, container,
      S1ap_ProtocolIE_ID_id_MME_UE_S1AP_ID, true);
  if (!ie) {
    OAILOG_FUNC_RETURN(LOG_S1AP, RETURNerror);
  }
  mme_ue_s1ap_id = ie->value.choice.MME_UE_S1AP_ID;

  S1AP_FIND_PROTOCOLIE_BY_ID(
      S1ap_InitialContextSetupFailureIEs_t, ie, container,
      S1ap_ProtocolIE_ID_id_eNB_UE_S1AP_ID, true);

  if (!ie) {
    OAILOG_FUNC_RETURN(LOG_S1AP, RETURNerror);
  }
  enb_ue_s1ap_id =
      (enb_ue_s1ap_id_t)(ie->value.choice.ENB_UE_S1AP_ID & ENB_UE_S1AP_ID_MASK);

  if ((ue_ref_p = s1ap_state_get_ue_mmeid(mme_ue_s1ap_id)) == NULL) {
    /*
     * MME doesn't know the MME UE S1AP ID provided.
     */
    OAILOG_INFO(
        LOG_S1AP,
        "INITIAL_CONTEXT_SETUP_FAILURE ignored. No context with "
        "mme_ue_s1ap_id " MME_UE_S1AP_ID_FMT
        " enb_ue_s1ap_id " ENB_UE_S1AP_ID_FMT " ",
        (uint32_t) mme_ue_s1ap_id, (uint32_t) enb_ue_s1ap_id);
    OAILOG_FUNC_RETURN(LOG_S1AP, RETURNerror);
  }

  if (ue_ref_p->enb_ue_s1ap_id != enb_ue_s1ap_id) {
    // abnormal case. No need to do anything. Ignore the message
    OAILOG_DEBUG(
        LOG_S1AP,
        "INITIAL_CONTEXT_SETUP_FAILURE ignored, mismatch enb_ue_s1ap_id: "
        "ctxt " ENB_UE_S1AP_ID_FMT " != received " ENB_UE_S1AP_ID_FMT " ",
        (uint32_t) ue_ref_p->enb_ue_s1ap_id, (uint32_t) enb_ue_s1ap_id);
    OAILOG_FUNC_RETURN(LOG_S1AP, RETURNerror);
  }

  s1ap_imsi_map_t* imsi_map = get_s1ap_imsi_map();
  hashtable_uint64_ts_get(
      imsi_map->mme_ue_id_imsi_htbl, (const hash_key_t) mme_ue_s1ap_id,
      &imsi64);

  // Pass this message to MME APP for necessary handling
  // Log the Cause Type and Cause value
  S1AP_FIND_PROTOCOLIE_BY_ID(
      S1ap_InitialContextSetupFailureIEs_t, ie, container,
      S1ap_ProtocolIE_ID_id_Cause, true);
  if (ie) {
    cause_type = ie->value.choice.Cause.present;
  } else {
    OAILOG_FUNC_RETURN(LOG_S1AP, RETURNerror);
  }

  switch (cause_type) {
    case S1ap_Cause_PR_radioNetwork:
      cause_value = ie->value.choice.Cause.choice.radioNetwork;
      OAILOG_DEBUG_UE(
          LOG_S1AP, imsi64,
          "INITIAL_CONTEXT_SETUP_FAILURE with Cause_Type = Radio Network and "
          "Cause_Value = %ld\n",
          cause_value);
      break;

    case S1ap_Cause_PR_transport:
      cause_value = ie->value.choice.Cause.choice.transport;
      OAILOG_DEBUG_UE(
          LOG_S1AP, imsi64,
          "INITIAL_CONTEXT_SETUP_FAILURE with Cause_Type = Transport and "
          "Cause_Value = %ld\n",
          cause_value);
      break;

    case S1ap_Cause_PR_nas:
      cause_value = ie->value.choice.Cause.choice.nas;
      OAILOG_DEBUG_UE(
          LOG_S1AP, imsi64,
          "INITIAL_CONTEXT_SETUP_FAILURE with Cause_Type = NAS and Cause_Value "
          "= "
          "%ld\n",
          cause_value);
      break;

    case S1ap_Cause_PR_protocol:
      cause_value = ie->value.choice.Cause.choice.protocol;
      OAILOG_DEBUG_UE(
          LOG_S1AP, imsi64,
          "INITIAL_CONTEXT_SETUP_FAILURE with Cause_Type = Protocol and "
          "Cause_Value = %ld\n",
          cause_value);
      break;

    case S1ap_Cause_PR_misc:
      cause_value = ie->value.choice.Cause.choice.misc;
      OAILOG_DEBUG_UE(
          LOG_S1AP, imsi64,
          "INITIAL_CONTEXT_SETUP_FAILURE with Cause_Type = MISC and "
          "Cause_Value "
          "= %ld\n",
          cause_value);
      break;

    default:
      OAILOG_DEBUG_UE(
          LOG_S1AP, imsi64,
          "INITIAL_CONTEXT_SETUP_FAILURE with Invalid Cause_Type = %d\n",
          cause_type);
      OAILOG_FUNC_RETURN(LOG_S1AP, RETURNerror);
  }
  message_p = DEPRECATEDitti_alloc_new_message_fatal(
      TASK_S1AP, MME_APP_INITIAL_CONTEXT_SETUP_FAILURE);
  memset(
      (void*) &message_p->ittiMsg.mme_app_initial_context_setup_failure, 0,
      sizeof(itti_mme_app_initial_context_setup_failure_t));
  MME_APP_INITIAL_CONTEXT_SETUP_FAILURE(message_p).mme_ue_s1ap_id =
      ue_ref_p->mme_ue_s1ap_id;

  message_p->ittiMsgHeader.imsi = imsi64;
  rc = send_msg_to_task(&s1ap_task_zmq_ctx, TASK_MME_APP, message_p);
  OAILOG_FUNC_RETURN(LOG_S1AP, rc);
}

status_code_e s1ap_mme_handle_ue_context_modification_response(
    s1ap_state_t* state, __attribute__((unused)) const sctp_assoc_id_t assoc_id,
    __attribute__((unused)) const sctp_stream_id_t stream,
    S1ap_S1AP_PDU_t* pdu) {
  S1ap_UEContextModificationResponseIEs_t *ie, *ie_enb = NULL;
  S1ap_UEContextModificationResponse_t* container = NULL;
  ue_description_t* ue_ref_p                      = NULL;
  MessageDef* message_p                           = NULL;
  int rc                                          = RETURNok;
  imsi64_t imsi64                                 = INVALID_IMSI64;

  OAILOG_FUNC_IN(LOG_S1AP);
  container =
      &pdu->choice.successfulOutcome.value.choice.UEContextModificationResponse;

  S1AP_FIND_PROTOCOLIE_BY_ID(
      S1ap_UEContextModificationResponseIEs_t, ie, container,
      S1ap_ProtocolIE_ID_id_MME_UE_S1AP_ID, true);
  S1AP_FIND_PROTOCOLIE_BY_ID(
      S1ap_UEContextModificationResponseIEs_t, ie_enb, container,
      S1ap_ProtocolIE_ID_id_eNB_UE_S1AP_ID, true);

  if (!ie) {
    OAILOG_ERROR(LOG_S1AP, "Missing mandatory MME UE S1AP ID ");
    return RETURNerror;
  }
  if (!ie_enb) {
    OAILOG_ERROR(LOG_S1AP, "Missing mandatory ENB UE S1AP ID ");
    return RETURNerror;
  }
  if ((ie) && (ue_ref_p = s1ap_state_get_ue_mmeid(
                   ie->value.choice.MME_UE_S1AP_ID)) == NULL) {
    /*
     * MME doesn't know the MME UE S1AP ID provided.
     * No need to do anything. Ignore the message
     */
    OAILOG_DEBUG(
        LOG_S1AP,
        "UE_CONTEXT_MODIFICATION_RESPONSE ignored cause could not get context "
        "with mme_ue_s1ap_id " MME_UE_S1AP_ID_FMT
        " enb_ue_s1ap_id " ENB_UE_S1AP_ID_FMT " ",
        (uint32_t) ie->value.choice.MME_UE_S1AP_ID,
        (uint32_t) ie_enb->value.choice.ENB_UE_S1AP_ID);
    OAILOG_FUNC_RETURN(LOG_S1AP, RETURNerror);
  } else {
    if ((ie_enb) &&
        (ue_ref_p->enb_ue_s1ap_id ==
         (ie_enb->value.choice.ENB_UE_S1AP_ID & ENB_UE_S1AP_ID_MASK))) {
      /*
       * Both eNB UE S1AP ID and MME UE S1AP ID match.
       * Send a UE context Release Command to eNB after releasing S1-U bearer
       * tunnel mapping for all the bearers.
       */

      s1ap_imsi_map_t* imsi_map = get_s1ap_imsi_map();
      hashtable_uint64_ts_get(
          imsi_map->mme_ue_id_imsi_htbl,
          (const hash_key_t) ie->value.choice.MME_UE_S1AP_ID, &imsi64);

      message_p = DEPRECATEDitti_alloc_new_message_fatal(
          TASK_S1AP, S1AP_UE_CONTEXT_MODIFICATION_RESPONSE);
      memset(
          (void*) &message_p->ittiMsg.s1ap_ue_context_mod_response, 0,
          sizeof(itti_s1ap_ue_context_mod_resp_t));
      S1AP_UE_CONTEXT_MODIFICATION_RESPONSE(message_p).mme_ue_s1ap_id =
          ue_ref_p->mme_ue_s1ap_id;
      S1AP_UE_CONTEXT_MODIFICATION_RESPONSE(message_p).enb_ue_s1ap_id =
          ue_ref_p->enb_ue_s1ap_id;

      message_p->ittiMsgHeader.imsi = imsi64;
      rc = send_msg_to_task(&s1ap_task_zmq_ctx, TASK_MME_APP, message_p);
      OAILOG_FUNC_RETURN(LOG_S1AP, rc);
    } else {
      // abnormal case. No need to do anything. Ignore the message
      OAILOG_DEBUG_UE(
          LOG_S1AP, imsi64,
          "S1AP_UE_CONTEXT_MODIFICATION_RESPONSE ignored, cause mismatch "
          "enb_ue_s1ap_id: ctxt" ENB_UE_S1AP_ID_FMT
          " != request " ENB_UE_S1AP_ID_FMT " ",
          (uint32_t) ue_ref_p->enb_ue_s1ap_id,
          (uint32_t) ie_enb->value.choice.ENB_UE_S1AP_ID);
      OAILOG_FUNC_RETURN(LOG_S1AP, RETURNerror);
    }
  }

  OAILOG_FUNC_RETURN(LOG_S1AP, RETURNok);
}

status_code_e s1ap_mme_handle_ue_context_modification_failure(
    s1ap_state_t* state, __attribute__((unused)) const sctp_assoc_id_t assoc_id,
    __attribute__((unused)) const sctp_stream_id_t stream,
    S1ap_S1AP_PDU_t* pdu) {
  S1ap_UEContextModificationFailureIEs_t *ie, *ie_enb = NULL;
  S1ap_UEContextModificationFailure_t* container = NULL;
  ue_description_t* ue_ref_p                     = NULL;
  MessageDef* message_p                          = NULL;
  int rc                                         = RETURNok;
  S1ap_Cause_PR cause_type                       = {0};
  int64_t cause_value;
  imsi64_t imsi64 = INVALID_IMSI64;

  OAILOG_FUNC_IN(LOG_S1AP);
  container = &pdu->choice.unsuccessfulOutcome.value.choice
                   .UEContextModificationFailure;

  S1AP_FIND_PROTOCOLIE_BY_ID(
      S1ap_UEContextModificationFailureIEs_t, ie, container,
      S1ap_ProtocolIE_ID_id_MME_UE_S1AP_ID, true);
  S1AP_FIND_PROTOCOLIE_BY_ID(
      S1ap_UEContextModificationFailureIEs_t, ie_enb, container,
      S1ap_ProtocolIE_ID_id_eNB_UE_S1AP_ID, true);

  if (!ie) {
    OAILOG_ERROR(LOG_S1AP, "Missing mandatory MME UE S1AP ID ");
    return RETURNerror;
  }
  if (!ie_enb) {
    OAILOG_ERROR(LOG_S1AP, "Missing mandatory ENB UE S1AP ID ");
    return RETURNerror;
  }

  if ((ie) && (ue_ref_p = s1ap_state_get_ue_mmeid(
                   ie->value.choice.MME_UE_S1AP_ID)) == NULL) {
    /*
     * MME doesn't know the MME UE S1AP ID provided.
     * No need to do anything. Ignore the message
     */
    OAILOG_DEBUG(
        LOG_S1AP,
        "UE_CONTEXT_MODIFICATION_FAILURE ignored cause could not get context "
        "with mme_ue_s1ap_id " MME_UE_S1AP_ID_FMT
        " enb_ue_s1ap_id " ENB_UE_S1AP_ID_FMT " ",
        (uint32_t) ie->value.choice.MME_UE_S1AP_ID,
        (uint32_t) ie_enb->value.choice.ENB_UE_S1AP_ID);
    OAILOG_FUNC_RETURN(LOG_S1AP, RETURNerror);
  } else {
    if ((ie_enb) &&
        (ue_ref_p->enb_ue_s1ap_id ==
         (ie_enb->value.choice.ENB_UE_S1AP_ID & ENB_UE_S1AP_ID_MASK))) {
      s1ap_imsi_map_t* imsi_map = get_s1ap_imsi_map();
      hashtable_uint64_ts_get(
          imsi_map->mme_ue_id_imsi_htbl,
          (const hash_key_t) ie->value.choice.MME_UE_S1AP_ID, &imsi64);

      // Pass this message to MME APP for necessary handling
      // Log the Cause Type and Cause value
      S1AP_FIND_PROTOCOLIE_BY_ID(
          S1ap_UEContextModificationFailureIEs_t, ie, container,
          S1ap_ProtocolIE_ID_id_Cause, true);
      if (ie) {
        cause_type = ie->value.choice.Cause.present;
      } else {
        OAILOG_FUNC_RETURN(LOG_S1AP, RETURNerror);
      }

      switch (cause_type) {
        case S1ap_Cause_PR_radioNetwork:
          cause_value = ie->value.choice.Cause.choice.radioNetwork;
          OAILOG_DEBUG_UE(
              LOG_S1AP, imsi64,
              "UE_CONTEXT_MODIFICATION_FAILURE with Cause_Type = Radio Network "
              "and Cause_Value = %ld\n",
              cause_value);
          break;

        case S1ap_Cause_PR_transport:
          cause_value = ie->value.choice.Cause.choice.transport;
          OAILOG_DEBUG_UE(
              LOG_S1AP, imsi64,
              "UE_CONTEXT_MODIFICATION_FAILURE with Cause_Type = Transport and "
              "Cause_Value = %ld\n",
              cause_value);
          break;

        case S1ap_Cause_PR_nas:
          cause_value = ie->value.choice.Cause.choice.nas;
          OAILOG_DEBUG_UE(
              LOG_S1AP, imsi64,
              "UE_CONTEXT_MODIFICATION_FAILURE with Cause_Type = NAS and "
              "Cause_Value = %ld\n",
              cause_value);
          break;

        case S1ap_Cause_PR_protocol:
          cause_value = ie->value.choice.Cause.choice.protocol;
          OAILOG_DEBUG_UE(
              LOG_S1AP, imsi64,
              "UE_CONTEXT_MODIFICATION_FAILURE with Cause_Type = Protocol and "
              "Cause_Value = %ld\n",
              cause_value);
          break;

        case S1ap_Cause_PR_misc:
          cause_value = ie->value.choice.Cause.choice.misc;
          OAILOG_DEBUG_UE(
              LOG_S1AP, imsi64,
              "UE_CONTEXT_MODIFICATION_FAILURE with Cause_Type = MISC and "
              "Cause_Value = %ld\n",
              cause_value);
          break;

        default:
          OAILOG_ERROR_UE(
              LOG_S1AP, imsi64,
              "UE_CONTEXT_MODIFICATION_FAILURE with Invalid Cause_Type = %d\n",
              cause_type);
          OAILOG_FUNC_RETURN(LOG_S1AP, RETURNerror);
      }
      message_p = DEPRECATEDitti_alloc_new_message_fatal(
          TASK_S1AP, S1AP_UE_CONTEXT_MODIFICATION_FAILURE);
      memset(
          (void*) &message_p->ittiMsg.s1ap_ue_context_mod_response, 0,
          sizeof(itti_s1ap_ue_context_mod_resp_fail_t));
      S1AP_UE_CONTEXT_MODIFICATION_FAILURE(message_p).mme_ue_s1ap_id =
          ue_ref_p->mme_ue_s1ap_id;
      S1AP_UE_CONTEXT_MODIFICATION_FAILURE(message_p).enb_ue_s1ap_id =
          ue_ref_p->enb_ue_s1ap_id;
      S1AP_UE_CONTEXT_MODIFICATION_FAILURE(message_p).cause = cause_value;

      message_p->ittiMsgHeader.imsi = imsi64;
      rc = send_msg_to_task(&s1ap_task_zmq_ctx, TASK_MME_APP, message_p);
      OAILOG_FUNC_RETURN(LOG_S1AP, rc);
    } else {
      // abnormal case. No need to do anything. Ignore the message
      OAILOG_DEBUG_UE(
          LOG_S1AP, imsi64,
          "S1AP_UE_CONTEXT_MODIFICATION_FAILURE ignored, cause mismatch "
          "enb_ue_s1ap_id: ctxt " ENB_UE_S1AP_ID_FMT
          " != request " ENB_UE_S1AP_ID_FMT " ",
          (uint32_t) ue_ref_p->enb_ue_s1ap_id,
          (uint32_t) ie_enb->value.choice.ENB_UE_S1AP_ID);
      OAILOG_FUNC_RETURN(LOG_S1AP, RETURNerror);
    }
  }

  OAILOG_FUNC_RETURN(LOG_S1AP, RETURNok);
}
////////////////////////////////////////////////////////////////////////////////
//************************ Handover signalling *******************************//
////////////////////////////////////////////////////////////////////////////////

//------------------------------------------------------------------------------

status_code_e s1ap_mme_handle_handover_request_ack(
    s1ap_state_t* state, __attribute__((unused)) const sctp_assoc_id_t assoc_id,
    __attribute__((unused)) const sctp_stream_id_t stream,
    S1ap_S1AP_PDU_t* pdu) {
  S1ap_HandoverRequestAcknowledge_t* container = NULL;
  S1ap_HandoverRequestAcknowledgeIEs_t* ie     = NULL;
  enb_description_t* source_enb                = NULL;
  enb_description_t* target_enb                = NULL;
  hashtable_element_array_t* enb_array         = NULL;
  uint32_t idx                                 = 0;
  ue_description_t* ue_ref_p                   = NULL;
  mme_ue_s1ap_id_t mme_ue_s1ap_id              = INVALID_MME_UE_S1AP_ID;
  enb_ue_s1ap_id_t tgt_enb_ue_s1ap_id          = INVALID_ENB_UE_S1AP_ID;
  S1ap_HandoverType_t handover_type            = -1;
  bstring tgt_src_container                    = {0};
  e_rab_admitted_list_t e_rab_list             = {0};
  imsi64_t imsi64                              = INVALID_IMSI64;
  s1ap_imsi_map_t* imsi_map                    = get_s1ap_imsi_map();

  OAILOG_FUNC_IN(LOG_S1AP);
  OAILOG_INFO(LOG_S1AP, "handover request ack received");

  container =
      &pdu->choice.successfulOutcome.value.choice.HandoverRequestAcknowledge;

  // MME_UE_S1AP_ID: mandatory
  S1AP_FIND_PROTOCOLIE_BY_ID(
      S1ap_HandoverRequestAcknowledgeIEs_t, ie, container,
      S1ap_ProtocolIE_ID_id_MME_UE_S1AP_ID, true);
  if (ie) {
    mme_ue_s1ap_id = ie->value.choice.MME_UE_S1AP_ID;
  } else {
    OAILOG_FUNC_RETURN(LOG_S1AP, RETURNerror);
  }

  // eNB_UE_S1AP_ID: mandatory
  S1AP_FIND_PROTOCOLIE_BY_ID(
      S1ap_HandoverRequestAcknowledgeIEs_t, ie, container,
      S1ap_ProtocolIE_ID_id_eNB_UE_S1AP_ID, true);
  // eNB UE S1AP ID is limited to 24 bits
  if (ie) {
    tgt_enb_ue_s1ap_id = (enb_ue_s1ap_id_t)(
        ie->value.choice.ENB_UE_S1AP_ID & ENB_UE_S1AP_ID_MASK);
  } else {
    OAILOG_FUNC_RETURN(LOG_S1AP, RETURNerror);
  }

  // E-RABAdmittedList: mandatory
  S1AP_FIND_PROTOCOLIE_BY_ID(
      S1ap_HandoverRequestAcknowledgeIEs_t, ie, container,
      S1ap_ProtocolIE_ID_id_E_RABAdmittedList, true);
  if (ie) {
    S1ap_E_RABAdmittedList_t* erab_admitted_list_req =
        &ie->value.choice.E_RABAdmittedList;
    int num_erab = ie->value.choice.E_RABAdmittedList.list.count;
    for (int i = 0; i < num_erab; i++) {
      S1ap_E_RABAdmittedItemIEs_t* erab_admitted_item_ies =
          (S1ap_E_RABAdmittedItemIEs_t*) erab_admitted_list_req->list.array[i];
      S1ap_E_RABAdmittedItem_t* erab_admitted_item_req =
          (S1ap_E_RABAdmittedItem_t*) &erab_admitted_item_ies->value.choice
              .E_RABAdmittedItem;
      e_rab_list.item[i].e_rab_id = erab_admitted_item_req->e_RAB_ID;
      e_rab_list.item[i].transport_layer_address = blk2bstr(
          erab_admitted_item_req->transportLayerAddress.buf,
          erab_admitted_item_req->transportLayerAddress.size);
      e_rab_list.item[i].gtp_teid =
          htonl(*((uint32_t*) erab_admitted_item_req->gTP_TEID.buf));
      // TODO: Add support for indirect data forwarding. Note that the DL and UL
      // transport address and GTP-TEID are optional, and only used if data
      // forwarding will take place. Since we don't currently support data
      // forwarding, these are ignored.
      e_rab_list.no_of_items += 1;
    }
  } else {
    OAILOG_FUNC_RETURN(LOG_S1AP, RETURNerror);
  }

  // Target To Source Transparent Container: mandatory
  S1AP_FIND_PROTOCOLIE_BY_ID(
      S1ap_HandoverRequestAcknowledgeIEs_t, ie, container,
      S1ap_ProtocolIE_ID_id_Target_ToSource_TransparentContainer, true);
  if (ie) {
    // note: ownership of tgt_src_container transferred to receiver
    tgt_src_container = blk2bstr(
        ie->value.choice.Target_ToSource_TransparentContainer.buf,
        ie->value.choice.Target_ToSource_TransparentContainer.size);
  } else {
    OAILOG_FUNC_RETURN(LOG_S1AP, RETURNerror);
  }

  // get imsi for logging
  hashtable_uint64_ts_get(
      imsi_map->mme_ue_id_imsi_htbl, (const hash_key_t) mme_ue_s1ap_id,
      &imsi64);

  // Retrieve the association ID for the eNB that UE is currently connected
  // (i.e., Source eNB) and pull the Source eNB record from s1ap state using
  // this association
  if ((ue_ref_p = s1ap_state_get_ue_mmeid(mme_ue_s1ap_id)) == NULL) {
    OAILOG_ERROR(
        LOG_S1AP,
        "MME_UE_S1AP_ID (" MME_UE_S1AP_ID_FMT
        ") does not point to any valid UE\n",
        mme_ue_s1ap_id);
    OAILOG_FUNC_RETURN(LOG_S1AP, RETURNerror);
  }
  if ((enb_array = hashtable_ts_get_elements(&state->enbs)) != NULL) {
    for (idx = 0; idx < enb_array->num_elements; idx++) {
      source_enb = (enb_description_t*) (uintptr_t) enb_array->elements[idx];
      if (source_enb->sctp_assoc_id == ue_ref_p->sctp_assoc_id) {
        break;
      }
    }
    free_wrapper((void**) &enb_array->elements);
    free_wrapper((void**) &enb_array);
    if (source_enb->sctp_assoc_id != ue_ref_p->sctp_assoc_id) {
      OAILOG_ERROR_UE(LOG_S1AP, imsi64, "No source eNB found for UE\n");
      OAILOG_FUNC_RETURN(LOG_S1AP, RETURNerror);
    }
  }

  OAILOG_INFO_UE(
      LOG_S1AP, imsi64, "Source enb is %u (association id %u)\n",
      source_enb->enb_id, source_enb->sctp_assoc_id);

  // get the target eNB -- the one that sent this message, target of the
  // handover
  target_enb = s1ap_state_get_enb(state, assoc_id);

  // handover type -- we only support intralte today, and reject all other
  // handover types when we receive HandoverRequired, so we can always assume
  // it's an intralte handover.
  handover_type = S1ap_HandoverType_intralte;

  // Add the e_rab_list to the UE's handover state -- we'll modify the bearers
  // if and when we receive the HANDOVER NOTIFY later in the procedure, so we
  // need to keep track of this.
  if (e_rab_list.no_of_items) {
    ue_ref_p->s1ap_handover_state.e_rab_admitted_list = e_rab_list;
  }

  s1ap_mme_itti_s1ap_handover_request_ack(
      mme_ue_s1ap_id, ue_ref_p->enb_ue_s1ap_id, tgt_enb_ue_s1ap_id,
      handover_type, source_enb->sctp_assoc_id, tgt_src_container,
      source_enb->enb_id, target_enb->enb_id, imsi64);

  OAILOG_FUNC_RETURN(LOG_S1AP, RETURNok);
}

status_code_e s1ap_mme_handle_handover_failure(
    s1ap_state_t* state, const sctp_assoc_id_t assoc_id,
    const sctp_stream_id_t stream, S1ap_S1AP_PDU_t* pdu) {
  S1ap_HandoverFailure_t* container = NULL;
  S1ap_HandoverFailureIEs_t* ie     = NULL;
  S1ap_S1AP_PDU_t out_pdu           = {0};
  S1ap_HandoverPreparationFailure_t* out;
  S1ap_HandoverPreparationFailureIEs_t* hpf_ie = NULL;
  ue_description_t* ue_ref_p                   = NULL;
  mme_ue_s1ap_id_t mme_ue_s1ap_id              = INVALID_MME_UE_S1AP_ID;
  S1ap_Cause_PR cause_type;
  long cause_value;
  uint8_t* buffer_p = NULL;
  uint8_t err       = 0;
  uint32_t length   = 0;

  OAILOG_FUNC_IN(LOG_S1AP);

  container = &pdu->choice.unsuccessfulOutcome.value.choice.HandoverFailure;

  // get the mandantory IEs
  // MME_UE_S1AP_ID
  S1AP_FIND_PROTOCOLIE_BY_ID(
      S1ap_HandoverFailureIEs_t, ie, container,
      S1ap_ProtocolIE_ID_id_MME_UE_S1AP_ID, true);
  if (ie) {
    mme_ue_s1ap_id = ie->value.choice.MME_UE_S1AP_ID;
  } else {
    OAILOG_FUNC_RETURN(LOG_S1AP, RETURNerror);
  }

  // Grab the Cause Type and Cause Value
  S1AP_FIND_PROTOCOLIE_BY_ID(
      S1ap_HandoverFailureIEs_t, ie, container, S1ap_ProtocolIE_ID_id_Cause,
      true);
  if (ie) {
    cause_type  = ie->value.choice.Cause.present;
    cause_value = s1ap_mme_get_cause_value(&ie->value.choice.Cause);
    if (cause_value == RETURNerror) {
      OAILOG_ERROR(
          LOG_S1AP, "HANDOVER FAILURE with Invalid Cause_Type = %d\n",
          cause_type);
      OAILOG_FUNC_RETURN(LOG_S1AP, RETURNerror);
    }
  } else {
    OAILOG_FUNC_RETURN(LOG_S1AP, RETURNerror);
  }

  // A failure means the target eNB was unable to prepare for handover. We need
  // to get rid of UE handover state and send a failure message with cause back
  // to the source eNB.

  // get UE context
  if ((ue_ref_p = s1ap_state_get_ue_mmeid(mme_ue_s1ap_id)) == NULL) {
    OAILOG_ERROR(
        LOG_S1AP,
        "could not get ue context for mme_ue_s1ap_id " MME_UE_S1AP_ID_FMT
        ", failing!\n",
        (uint32_t) mme_ue_s1ap_id);
    OAILOG_FUNC_RETURN(LOG_S1AP, RETURNerror);
  }

  if (ue_ref_p->s1_ue_state == S1AP_UE_HANDOVER) {
    // this effectively cancels the HandoverPreparation proecedure as we
    // only send a HandoverCommand if the UE is in the S1AP_UE_HANDOVER
    // state.
    ue_ref_p->s1_ue_state         = S1AP_UE_CONNECTED;
    ue_ref_p->s1ap_handover_state = (struct s1ap_handover_state_s){0};
  } else {
    // Not a failure, but nothing for us to do.
    OAILOG_INFO(
        LOG_S1AP,
        "Received HANDOVER FAILURE for UE not in handover state, leaving UE "
        "state unmodified for MME_S1AP_UE_ID " MME_UE_S1AP_ID_FMT ".\n",
        (uint32_t) mme_ue_s1ap_id);
  }

  // generate HandoverPreparationFailure
  out_pdu.present = S1ap_S1AP_PDU_PR_unsuccessfulOutcome;
  out_pdu.choice.unsuccessfulOutcome.procedureCode =
      S1ap_ProcedureCode_id_HandoverPreparation;
  out_pdu.choice.unsuccessfulOutcome.value.present =
      S1ap_UnsuccessfulOutcome__value_PR_HandoverPreparationFailure;
  out_pdu.choice.unsuccessfulOutcome.criticality = S1ap_Criticality_ignore;
  out = &out_pdu.choice.unsuccessfulOutcome.value.choice
             .HandoverPreparationFailure;

  // mme_ue_s1ap_id (mandatory)
  hpf_ie = (S1ap_HandoverPreparationFailureIEs_t*) calloc(
      1, sizeof(S1ap_HandoverPreparationFailureIEs_t));
  hpf_ie->id          = S1ap_ProtocolIE_ID_id_MME_UE_S1AP_ID;
  hpf_ie->criticality = S1ap_Criticality_ignore;
  hpf_ie->value.present =
      S1ap_HandoverPreparationFailureIEs__value_PR_MME_UE_S1AP_ID;
  hpf_ie->value.choice.MME_UE_S1AP_ID = mme_ue_s1ap_id;
  ASN_SEQUENCE_ADD(&out->protocolIEs.list, hpf_ie);

  // source enb_ue_s1ap_id (mandatory)
  hpf_ie = (S1ap_HandoverPreparationFailureIEs_t*) calloc(
      1, sizeof(S1ap_HandoverPreparationFailureIEs_t));
  hpf_ie->id          = S1ap_ProtocolIE_ID_id_eNB_UE_S1AP_ID;
  hpf_ie->criticality = S1ap_Criticality_ignore;
  hpf_ie->value.present =
      S1ap_HandoverPreparationFailureIEs__value_PR_ENB_UE_S1AP_ID;
  hpf_ie->value.choice.ENB_UE_S1AP_ID = ue_ref_p->enb_ue_s1ap_id;
  ASN_SEQUENCE_ADD(&out->protocolIEs.list, hpf_ie);

  // cause (mandatory)
  hpf_ie = (S1ap_HandoverPreparationFailureIEs_t*) calloc(
      1, sizeof(S1ap_HandoverPreparationFailureIEs_t));
  hpf_ie->id            = S1ap_ProtocolIE_ID_id_Cause;
  hpf_ie->criticality   = S1ap_Criticality_ignore;
  hpf_ie->value.present = S1ap_HandoverPreparationFailureIEs__value_PR_Cause;
  s1ap_mme_set_cause(&hpf_ie->value.choice.Cause, cause_type, cause_value);
  ASN_SEQUENCE_ADD(&out->protocolIEs.list, hpf_ie);

  // Construct the PDU and send message
  if (s1ap_mme_encode_pdu(&out_pdu, &buffer_p, &length) < 0) {
    err = 1;
  }
  ASN_STRUCT_FREE_CONTENTS_ONLY(asn_DEF_S1ap_HandoverPreparationFailure, out);
  if (err) {
    OAILOG_FUNC_RETURN(LOG_S1AP, RETURNerror);
  }
  bstring b = blk2bstr(buffer_p, length);
  free(buffer_p);

  OAILOG_DEBUG(
      LOG_S1AP,
      "send HANDOVER PREPARATION FAILURE for mme_ue_s1ap_id " MME_UE_S1AP_ID_FMT
      "\n",
      (uint32_t) mme_ue_s1ap_id);

  s1ap_mme_itti_send_sctp_request(
      &b, ue_ref_p->sctp_assoc_id, ue_ref_p->sctp_stream_send, mme_ue_s1ap_id);

  OAILOG_FUNC_RETURN(LOG_S1AP, RETURNok);
}

status_code_e s1ap_mme_handle_handover_cancel(
    s1ap_state_t* state, const sctp_assoc_id_t assoc_id,
    const sctp_stream_id_t stream, S1ap_S1AP_PDU_t* pdu) {
  S1ap_HandoverCancel_t* container = NULL;
  S1ap_HandoverCancelIEs_t* ie     = NULL;
  S1ap_S1AP_PDU_t out_pdu          = {0};
  S1ap_HandoverCancelAcknowledge_t* out;
  S1ap_HandoverCancelAcknowledgeIEs_t* hca_ie = NULL;
  ue_description_t* ue_ref_p                  = NULL;
  e_rab_admitted_list_t e_rab_admitted_list   = {0};
  mme_ue_s1ap_id_t mme_ue_s1ap_id             = INVALID_MME_UE_S1AP_ID;
  enb_ue_s1ap_id_t enb_ue_s1ap_id             = INVALID_ENB_UE_S1AP_ID;
  S1ap_Cause_PR cause_type;
  long cause_value;
  uint8_t* buffer_p = NULL;
  uint8_t err       = 0;
  uint32_t length   = 0;

  OAILOG_FUNC_IN(LOG_S1AP);

  container = &pdu->choice.initiatingMessage.value.choice.HandoverCancel;

  // get the mandantory IEs
  // MME_UE_S1AP_ID
  S1AP_FIND_PROTOCOLIE_BY_ID(
      S1ap_HandoverCancelIEs_t, ie, container,
      S1ap_ProtocolIE_ID_id_MME_UE_S1AP_ID, true);
  if (ie) {
    mme_ue_s1ap_id = ie->value.choice.MME_UE_S1AP_ID;
  } else {
    OAILOG_FUNC_RETURN(LOG_S1AP, RETURNerror);
  }

  // eNB_UE_S1AP_ID
  S1AP_FIND_PROTOCOLIE_BY_ID(
      S1ap_HandoverCancelIEs_t, ie, container,
      S1ap_ProtocolIE_ID_id_eNB_UE_S1AP_ID, true);
  // eNB UE S1AP ID is limited to 24 bits
  if (ie) {
    enb_ue_s1ap_id = (enb_ue_s1ap_id_t)(
        ie->value.choice.ENB_UE_S1AP_ID & ENB_UE_S1AP_ID_MASK);
  } else {
    OAILOG_FUNC_RETURN(LOG_S1AP, RETURNerror);
  }

  // Grab the Cause Type and Cause Value
  S1AP_FIND_PROTOCOLIE_BY_ID(
      S1ap_HandoverCancelIEs_t, ie, container, S1ap_ProtocolIE_ID_id_Cause,
      true);
  if (ie) {
    cause_type  = ie->value.choice.Cause.present;
    cause_value = s1ap_mme_get_cause_value(&ie->value.choice.Cause);
    if (cause_value == RETURNerror) {
      OAILOG_ERROR(
          LOG_S1AP, "HANDOVER CANCEL with Invalid Cause_Type = %d\n",
          cause_type);
      OAILOG_FUNC_RETURN(LOG_S1AP, RETURNerror);
    }
  } else {
    OAILOG_FUNC_RETURN(LOG_S1AP, RETURNerror);
  }

  OAILOG_INFO(
      LOG_S1AP,
      "HANDOVER CANCEL from association ID %u for MME UE S1AP ID "
      "(" MME_UE_S1AP_ID_FMT "), CauseType= %u CauseValue = %ld\n",
      assoc_id, mme_ue_s1ap_id, cause_type, cause_value);

  // make sure any handover state in the UE is reset, move the UE back to
  // connected state, and generate a cancel acknowledgement (immediately).

  // get UE context
  if ((ue_ref_p = s1ap_state_get_ue_mmeid(mme_ue_s1ap_id)) == NULL) {
    OAILOG_ERROR(
        LOG_S1AP,
        "could not get ue context for mme_ue_s1ap_id " MME_UE_S1AP_ID_FMT
        ", failing!\n",
        (uint32_t) mme_ue_s1ap_id);
    OAILOG_FUNC_RETURN(LOG_S1AP, RETURNerror);
  }

  if (ue_ref_p->s1_ue_state == S1AP_UE_HANDOVER) {
    // this effectively cancels the HandoverPreparation proecedure as we
    // only send a HandoverCommand if the UE is in the S1AP_UE_HANDOVER
    // state.
    ue_ref_p->s1_ue_state = S1AP_UE_CONNECTED;
    /* Free all the transport layer address pointers in ERAB admitted list
     * before actually resetting the S1AP handover state
     */
    e_rab_admitted_list = ue_ref_p->s1ap_handover_state.e_rab_admitted_list;
    for (int i = 0; i < e_rab_admitted_list.no_of_items; i++) {
      bdestroy_wrapper(&e_rab_admitted_list.item[i].transport_layer_address);
    }
    ue_ref_p->s1ap_handover_state = (struct s1ap_handover_state_s){0};
  } else {
    // Not a failure, but nothing for us to do.
    OAILOG_INFO(
        LOG_S1AP,
        "Received HANDOVER CANCEL for UE not in handover state, leaving UE "
        "state unmodified for MME_S1AP_UE_ID " MME_UE_S1AP_ID_FMT ".\n",
        (uint32_t) mme_ue_s1ap_id);
  }

  // generate the cancel acknowledge
  out_pdu.present = S1ap_S1AP_PDU_PR_successfulOutcome;
  out_pdu.choice.successfulOutcome.procedureCode =
      S1ap_ProcedureCode_id_HandoverCancel;
  out_pdu.choice.successfulOutcome.value.present =
      S1ap_SuccessfulOutcome__value_PR_HandoverCancelAcknowledge;
  out_pdu.choice.successfulOutcome.criticality = S1ap_Criticality_ignore;
  out =
      &out_pdu.choice.successfulOutcome.value.choice.HandoverCancelAcknowledge;

  /* MME-UE-ID: mandatory */
  hca_ie = (S1ap_HandoverCancelAcknowledgeIEs_t*) calloc(
      1, sizeof(S1ap_HandoverCancelAcknowledgeIEs_t));
  hca_ie->id          = S1ap_ProtocolIE_ID_id_MME_UE_S1AP_ID;
  hca_ie->criticality = S1ap_Criticality_ignore;
  hca_ie->value.present =
      S1ap_HandoverCancelAcknowledgeIEs__value_PR_MME_UE_S1AP_ID;
  hca_ie->value.choice.MME_UE_S1AP_ID = mme_ue_s1ap_id;
  ASN_SEQUENCE_ADD(&out->protocolIEs.list, hca_ie);

  /* eNB-UE-ID: mandatory */
  hca_ie = (S1ap_HandoverCancelAcknowledgeIEs_t*) calloc(
      1, sizeof(S1ap_HandoverCancelAcknowledgeIEs_t));
  hca_ie->id          = S1ap_ProtocolIE_ID_id_eNB_UE_S1AP_ID;
  hca_ie->criticality = S1ap_Criticality_ignore;
  hca_ie->value.present =
      S1ap_HandoverCancelAcknowledgeIEs__value_PR_ENB_UE_S1AP_ID;
  hca_ie->value.choice.ENB_UE_S1AP_ID = enb_ue_s1ap_id;
  ASN_SEQUENCE_ADD(&out->protocolIEs.list, hca_ie);

  // Construct the PDU and send message
  if (s1ap_mme_encode_pdu(&out_pdu, &buffer_p, &length) < 0) {
    err = 1;
  }
  ASN_STRUCT_FREE_CONTENTS_ONLY(asn_DEF_S1ap_HandoverCancelAcknowledge, out);
  if (err) {
    OAILOG_FUNC_RETURN(LOG_S1AP, RETURNerror);
  }
  bstring b = blk2bstr(buffer_p, length);
  free(buffer_p);

  s1ap_mme_itti_send_sctp_request(&b, assoc_id, stream, mme_ue_s1ap_id);

  OAILOG_FUNC_RETURN(LOG_S1AP, RETURNok);
}

status_code_e s1ap_mme_handle_handover_request(
    s1ap_state_t* state, const itti_mme_app_handover_request_t* ho_request_p) {
  uint8_t* buffer_p   = NULL;
  uint8_t err         = 0;
  uint32_t length     = 0;
  S1ap_S1AP_PDU_t pdu = {0};
  S1ap_HandoverRequest_t* out;
  S1ap_HandoverRequestIEs_t* ie = NULL;
  enb_description_t* target_enb = NULL;
  sctp_stream_id_t stream       = 0x0;
  ue_description_t* ue_ref_p    = NULL;

  OAILOG_FUNC_IN(LOG_S1AP);
  if (ho_request_p == NULL) {
    OAILOG_ERROR(LOG_S1AP, "Handover Request is null\n");
    OAILOG_FUNC_RETURN(LOG_S1AP, RETURNerror);
  }

  OAILOG_INFO(LOG_S1AP, "Handover Request received");

  // get the ue description
  if ((ue_ref_p = s1ap_state_get_ue_mmeid(ho_request_p->mme_ue_s1ap_id)) ==
      NULL) {
    OAILOG_ERROR(
        LOG_S1AP,
        "could not get ue context for mme_ue_s1ap_id " MME_UE_S1AP_ID_FMT
        ", failing!\n",
        (uint32_t) ho_request_p->mme_ue_s1ap_id);
    OAILOG_FUNC_RETURN(LOG_S1AP, RETURNerror);
  }

  if ((target_enb = s1ap_state_get_enb(
           state, ho_request_p->target_sctp_assoc_id)) == NULL) {
    OAILOG_ERROR(
        LOG_S1AP, "Could not get enb description for assoc_id %u\n",
        ho_request_p->target_sctp_assoc_id);
    OAILOG_FUNC_RETURN(LOG_S1AP, RETURNerror);
  }

  // set the recv and send streams for UE on the target.
  stream = target_enb->next_sctp_stream;
  ue_ref_p->s1ap_handover_state.target_sctp_stream_recv = stream;
  ue_ref_p->s1ap_handover_state.source_sctp_stream_recv =
      ue_ref_p->sctp_stream_recv;
  target_enb->next_sctp_stream += 1;
  if (target_enb->next_sctp_stream >= target_enb->instreams) {
    target_enb->next_sctp_stream = 1;
  }
  ue_ref_p->s1ap_handover_state.target_sctp_stream_send =
      target_enb->next_sctp_stream;
  ue_ref_p->s1ap_handover_state.source_sctp_stream_send =
      ue_ref_p->sctp_stream_send;

  // Build and send PDU
  pdu.present = S1ap_S1AP_PDU_PR_initiatingMessage;
  pdu.choice.initiatingMessage.procedureCode =
      S1ap_ProcedureCode_id_HandoverResourceAllocation;
  pdu.choice.initiatingMessage.value.present =
      S1ap_InitiatingMessage__value_PR_HandoverRequest;
  pdu.choice.initiatingMessage.criticality = S1ap_Criticality_reject;
  out = &pdu.choice.initiatingMessage.value.choice.HandoverRequest;

  /* MME-UE-ID: mandatory */
  ie =
      (S1ap_HandoverRequestIEs_t*) calloc(1, sizeof(S1ap_HandoverRequestIEs_t));
  ie->id            = S1ap_ProtocolIE_ID_id_MME_UE_S1AP_ID;
  ie->criticality   = S1ap_Criticality_reject;
  ie->value.present = S1ap_HandoverRequestIEs__value_PR_MME_UE_S1AP_ID;
  ie->value.choice.MME_UE_S1AP_ID = ho_request_p->mme_ue_s1ap_id;
  ASN_SEQUENCE_ADD(&out->protocolIEs.list, ie);

  /* HandoverType: mandatory */
  ie =
      (S1ap_HandoverRequestIEs_t*) calloc(1, sizeof(S1ap_HandoverRequestIEs_t));
  ie->id            = S1ap_ProtocolIE_ID_id_HandoverType;
  ie->criticality   = S1ap_Criticality_reject;
  ie->value.present = S1ap_HandoverRequestIEs__value_PR_HandoverType;
  ie->value.choice.HandoverType = ho_request_p->handover_type;
  ASN_SEQUENCE_ADD(&out->protocolIEs.list, ie);

  /* Cause: mandatory */
  ie =
      (S1ap_HandoverRequestIEs_t*) calloc(1, sizeof(S1ap_HandoverRequestIEs_t));
  ie->id                 = S1ap_ProtocolIE_ID_id_Cause;
  ie->criticality        = S1ap_Criticality_ignore;
  ie->value.present      = S1ap_HandoverRequestIEs__value_PR_Cause;
  ie->value.choice.Cause = ho_request_p->cause;
  ASN_SEQUENCE_ADD(&out->protocolIEs.list, ie);

  /* ambr: mandatory */
  ie =
      (S1ap_HandoverRequestIEs_t*) calloc(1, sizeof(S1ap_HandoverRequestIEs_t));
  ie->id          = S1ap_ProtocolIE_ID_id_uEaggregateMaximumBitrate;
  ie->criticality = S1ap_Criticality_reject;
  ie->value.present =
      S1ap_HandoverRequestIEs__value_PR_UEAggregateMaximumBitrate;
  asn_uint642INTEGER(
      &ie->value.choice.UEAggregateMaximumBitrate.uEaggregateMaximumBitRateDL,
      ho_request_p->ue_ambr.br_dl);
  asn_uint642INTEGER(
      &ie->value.choice.UEAggregateMaximumBitrate.uEaggregateMaximumBitRateUL,
      ho_request_p->ue_ambr.br_ul);
  ASN_SEQUENCE_ADD(&out->protocolIEs.list, ie);

  /* e-rab to be setup list: mandatory */
  ie =
      (S1ap_HandoverRequestIEs_t*) calloc(1, sizeof(S1ap_HandoverRequestIEs_t));
  ie->id            = S1ap_ProtocolIE_ID_id_E_RABToBeSetupListHOReq;
  ie->criticality   = S1ap_Criticality_reject;
  ie->value.present = S1ap_HandoverRequestIEs__value_PR_E_RABToBeSetupListHOReq;
  ASN_SEQUENCE_ADD(&out->protocolIEs.list, ie);
  S1ap_E_RABToBeSetupListHOReq_t* const e_rab_to_be_setup_list =
      &ie->value.choice.E_RABToBeSetupListHOReq;

  for (int i = 0; i < ho_request_p->e_rab_list.no_of_items; i++) {
    S1ap_E_RABToBeSetupItemHOReqIEs_t* e_rab_tobesetup_item =
        (S1ap_E_RABToBeSetupItemHOReqIEs_t*) calloc(
            1, sizeof(S1ap_E_RABToBeSetupItemHOReqIEs_t));

    e_rab_tobesetup_item->id = S1ap_ProtocolIE_ID_id_E_RABToBeSetupItemHOReq;
    e_rab_tobesetup_item->criticality = S1ap_Criticality_reject;
    e_rab_tobesetup_item->value.present =
        S1ap_E_RABToBeSetupItemHOReqIEs__value_PR_E_RABToBeSetupItemHOReq;
    S1ap_E_RABToBeSetupItemHOReq_t* e_RABToBeSetup =
        &e_rab_tobesetup_item->value.choice.E_RABToBeSetupItemHOReq;

    // e_rab_id
    e_RABToBeSetup->e_RAB_ID = ho_request_p->e_rab_list.item[i].e_rab_id;

    // transportLayerAddress
    e_RABToBeSetup->transportLayerAddress.buf = calloc(
        blength(ho_request_p->e_rab_list.item[i].transport_layer_address),
        sizeof(uint8_t));
    memcpy(
        e_RABToBeSetup->transportLayerAddress.buf,
        ho_request_p->e_rab_list.item[i].transport_layer_address->data,
        blength(ho_request_p->e_rab_list.item[i].transport_layer_address));
    e_RABToBeSetup->transportLayerAddress.size =
        blength(ho_request_p->e_rab_list.item[i].transport_layer_address);
    e_RABToBeSetup->transportLayerAddress.bits_unused = 0;

    // gtp-teid
    INT32_TO_OCTET_STRING(
        ho_request_p->e_rab_list.item[i].gtp_teid, &e_RABToBeSetup->gTP_TEID);

    // qos params
    e_RABToBeSetup->e_RABlevelQosParameters.qCI =
        ho_request_p->e_rab_list.item[i].e_rab_level_qos_parameters.qci;
    e_RABToBeSetup->e_RABlevelQosParameters.allocationRetentionPriority
        .priorityLevel = ho_request_p->e_rab_list.item[i]
                             .e_rab_level_qos_parameters
                             .allocation_and_retention_priority.priority_level;
    e_RABToBeSetup->e_RABlevelQosParameters.allocationRetentionPriority
        .pre_emptionCapability =
        ho_request_p->e_rab_list.item[i]
            .e_rab_level_qos_parameters.allocation_and_retention_priority
            .pre_emption_capability;
    e_RABToBeSetup->e_RABlevelQosParameters.allocationRetentionPriority
        .pre_emptionVulnerability =
        ho_request_p->e_rab_list.item[i]
            .e_rab_level_qos_parameters.allocation_and_retention_priority
            .pre_emption_vulnerability;

    // data forwarding not supported
    OAILOG_INFO(LOG_S1AP, "Note: data forwarding unsupported\n");
    S1ap_E_RABToBeSetupItemHOReq_ExtIEs_t* exts =
        (S1ap_E_RABToBeSetupItemHOReq_ExtIEs_t*) calloc(
            1, sizeof(S1ap_E_RABToBeSetupItemHOReq_ExtIEs_t));
    exts->id          = S1ap_ProtocolIE_ID_id_Data_Forwarding_Not_Possible;
    exts->criticality = S1ap_Criticality_ignore;
    exts->extensionValue.present =
        S1ap_E_RABToBeSetupItemHOReq_ExtIEs__extensionValue_PR_Data_Forwarding_Not_Possible;
    exts->extensionValue.choice.Data_Forwarding_Not_Possible =
        S1ap_Data_Forwarding_Not_Possible_data_Forwarding_not_Possible;

    S1ap_ProtocolExtensionContainer_7327P1_t* xc =
        (S1ap_ProtocolExtensionContainer_7327P1_t*) calloc(
            1, sizeof(S1ap_ProtocolExtensionContainer_7327P1_t));
    int asn_ret = 0;
    asn_ret     = ASN_SEQUENCE_ADD(&xc->list, exts);
    if (asn_ret) {
      OAILOG_ERROR(LOG_S1AP, "ASN_SEQUENCE_ADD ret = %d\n", asn_ret);
    }

    // Bad cast...
    e_RABToBeSetup->iE_Extensions =
        (struct S1ap_ProtocolExtensionContainer*) xc;

    ASN_SEQUENCE_ADD(&e_rab_to_be_setup_list->list, e_rab_tobesetup_item);
  }

  /* Source-ToTarget-TransparentContainer: mandatory */
  ie =
      (S1ap_HandoverRequestIEs_t*) calloc(1, sizeof(S1ap_HandoverRequestIEs_t));
  ie->id          = S1ap_ProtocolIE_ID_id_Source_ToTarget_TransparentContainer;
  ie->criticality = S1ap_Criticality_reject;
  ie->value.present =
      S1ap_HandoverRequestIEs__value_PR_Source_ToTarget_TransparentContainer;
  OCTET_STRING_fromBuf(
      &ie->value.choice.Source_ToTarget_TransparentContainer,
      (char*) bdata(ho_request_p->src_tgt_container),
      blength(ho_request_p->src_tgt_container));
  ASN_SEQUENCE_ADD(&out->protocolIEs.list, ie);

  /* UESecurityCapabilities: mandatory */
  ie =
      (S1ap_HandoverRequestIEs_t*) calloc(1, sizeof(S1ap_HandoverRequestIEs_t));
  ie->id            = S1ap_ProtocolIE_ID_id_UESecurityCapabilities;
  ie->criticality   = S1ap_Criticality_reject;
  ie->value.present = S1ap_HandoverRequestIEs__value_PR_UESecurityCapabilities;

  S1ap_UESecurityCapabilities_t* const ue_security_capabilities =
      &ie->value.choice.UESecurityCapabilities;

  ue_security_capabilities->encryptionAlgorithms.buf =
      calloc(1, sizeof(uint16_t));
  memcpy(
      ue_security_capabilities->encryptionAlgorithms.buf,
      &ho_request_p->encryption_algorithm_capabilities, sizeof(uint16_t));
  ue_security_capabilities->encryptionAlgorithms.size        = 2;
  ue_security_capabilities->encryptionAlgorithms.bits_unused = 0;
  OAILOG_DEBUG(
      LOG_S1AP, "security_capabilities_encryption_algorithms 0x%04X\n",
      ho_request_p->encryption_algorithm_capabilities);

  ue_security_capabilities->integrityProtectionAlgorithms.buf =
      calloc(1, sizeof(uint16_t));
  memcpy(
      ue_security_capabilities->integrityProtectionAlgorithms.buf,
      &ho_request_p->integrity_algorithm_capabilities, sizeof(uint16_t));
  ue_security_capabilities->integrityProtectionAlgorithms.size        = 2;
  ue_security_capabilities->integrityProtectionAlgorithms.bits_unused = 0;
  OAILOG_DEBUG(
      LOG_S1AP, "security_capabilities_integrity_algorithms 0x%04X\n",
      ho_request_p->integrity_algorithm_capabilities);
  ASN_SEQUENCE_ADD(&out->protocolIEs.list, ie);

  /* SecurityContext: mandatory */
  ie =
      (S1ap_HandoverRequestIEs_t*) calloc(1, sizeof(S1ap_HandoverRequestIEs_t));
  ie->id            = S1ap_ProtocolIE_ID_id_SecurityContext;
  ie->criticality   = S1ap_Criticality_reject;
  ie->value.present = S1ap_HandoverRequestIEs__value_PR_SecurityContext;

  S1ap_SecurityContext_t* const security_context =
      &ie->value.choice.SecurityContext;
  security_context->nextHopChainingCount = ho_request_p->ncc;
  security_context->nextHopParameter.buf =
      calloc(AUTH_NEXT_HOP_SIZE, sizeof(uint8_t));
  memcpy(
      security_context->nextHopParameter.buf, &ho_request_p->nh,
      AUTH_NEXT_HOP_SIZE);
  security_context->nextHopParameter.size        = AUTH_NEXT_HOP_SIZE;
  security_context->nextHopParameter.bits_unused = 0;
  ASN_SEQUENCE_ADD(&out->protocolIEs.list, ie);

  // Construct the PDU and send message
  if (s1ap_mme_encode_pdu(&pdu, &buffer_p, &length) < 0) {
    err = 1;
  }
  ASN_STRUCT_FREE_CONTENTS_ONLY(asn_DEF_S1ap_HandoverRequest, out);
  if (err) {
    OAILOG_FUNC_RETURN(LOG_S1AP, RETURNerror);
  }
  bstring b = blk2bstr(buffer_p, length);
  free(buffer_p);

  s1ap_mme_itti_send_sctp_request(
      &b, ho_request_p->target_sctp_assoc_id, stream,
      ho_request_p->mme_ue_s1ap_id);

  OAILOG_FUNC_RETURN(LOG_S1AP, RETURNok);
}

status_code_e s1ap_mme_handle_handover_required(
    s1ap_state_t* state, __attribute__((unused)) const sctp_assoc_id_t assoc_id,
    __attribute__((unused)) const sctp_stream_id_t stream,
    S1ap_S1AP_PDU_t* pdu) {
  S1ap_HandoverRequired_t* container = NULL;
  S1ap_HandoverRequiredIEs_t* ie     = NULL;
  enb_description_t* enb_association = NULL;
  mme_ue_s1ap_id_t mme_ue_s1ap_id    = INVALID_MME_UE_S1AP_ID;
  enb_ue_s1ap_id_t enb_ue_s1ap_id    = INVALID_ENB_UE_S1AP_ID;
  S1ap_HandoverType_t handover_type  = -1;
  S1ap_Cause_t cause                 = {0};
  S1ap_Cause_PR cause_type;
  long cause_value;
  S1ap_TargeteNB_ID_t* targeteNB_ID         = NULL;
  bstring src_tgt_container                 = {0};
  uint8_t* enb_id_buf                       = NULL;
  enb_description_t* target_enb_association = NULL;
  hashtable_element_array_t* enb_array      = NULL;
  uint32_t target_enb_id                    = 0;
  uint32_t idx                              = 0;
  imsi64_t imsi64                           = INVALID_IMSI64;
  s1ap_imsi_map_t* imsi_map                 = get_s1ap_imsi_map();

  OAILOG_FUNC_IN(LOG_S1AP);

  enb_association = s1ap_state_get_enb(state, assoc_id);
  if (enb_association == NULL) {
    OAILOG_ERROR(
        LOG_S1AP,
        "Ignore Handover Required from unknown assoc "
        "%u\n",
        assoc_id);
    OAILOG_FUNC_RETURN(LOG_S1AP, RETURNerror);
  }

  OAILOG_INFO(
      LOG_S1AP,
      "Handover Required from association id %u, "
      "Connected UEs = %u Num elements = %zu\n",
      assoc_id, enb_association->nb_ue_associated,
      enb_association->ue_id_coll.num_elements);

  container = &pdu->choice.initiatingMessage.value.choice.HandoverRequired;

  // MME_UE_S1AP_ID
  S1AP_FIND_PROTOCOLIE_BY_ID(
      S1ap_HandoverRequiredIEs_t, ie, container,
      S1ap_ProtocolIE_ID_id_MME_UE_S1AP_ID, true);
  if (ie) {
    mme_ue_s1ap_id = ie->value.choice.MME_UE_S1AP_ID;
  } else {
    OAILOG_FUNC_RETURN(LOG_S1AP, RETURNerror);
  }

  // eNB_UE_S1AP_ID
  S1AP_FIND_PROTOCOLIE_BY_ID(
      S1ap_HandoverRequiredIEs_t, ie, container,
      S1ap_ProtocolIE_ID_id_eNB_UE_S1AP_ID, true);
  // eNB UE S1AP ID is limited to 24 bits
  if (ie) {
    enb_ue_s1ap_id = (enb_ue_s1ap_id_t)(
        ie->value.choice.ENB_UE_S1AP_ID & ENB_UE_S1AP_ID_MASK);
  } else {
    OAILOG_FUNC_RETURN(LOG_S1AP, RETURNerror);
  }

  // Handover Type
  S1AP_FIND_PROTOCOLIE_BY_ID(
      S1ap_HandoverRequiredIEs_t, ie, container,
      S1ap_ProtocolIE_ID_id_HandoverType, true);
  if (ie) {
    handover_type = ie->value.choice.HandoverType;
  } else {
    OAILOG_FUNC_RETURN(LOG_S1AP, RETURNerror);
  }

  // Only support intra LTE handovers today.
  if (handover_type != S1ap_HandoverType_intralte) {
    OAILOG_ERROR(
        LOG_S1AP,
        "Unsupported handover type "
        "%ld\n",
        handover_type);

    // TODO: Process a failure message
    OAILOG_FUNC_RETURN(LOG_S1AP, RETURNerror);
  }

  // Grab the Cause Type and Cause Value
  S1AP_FIND_PROTOCOLIE_BY_ID(
      S1ap_HandoverRequiredIEs_t, ie, container, S1ap_ProtocolIE_ID_id_Cause,
      true);
  if (ie) {
    cause_type = ie->value.choice.Cause.present;
    cause      = ie->value.choice.Cause;
  } else {
    OAILOG_FUNC_RETURN(LOG_S1AP, RETURNerror);
  }
  switch (cause_type) {
    case S1ap_Cause_PR_radioNetwork:
      cause_value = ie->value.choice.Cause.choice.radioNetwork;
      break;

    case S1ap_Cause_PR_transport:
      cause_value = ie->value.choice.Cause.choice.transport;
      break;

    case S1ap_Cause_PR_nas:
      cause_value = ie->value.choice.Cause.choice.nas;
      break;

    case S1ap_Cause_PR_protocol:
      cause_value = ie->value.choice.Cause.choice.protocol;
      break;

    case S1ap_Cause_PR_misc:
      cause_value = ie->value.choice.Cause.choice.misc;
      break;

    default:
      OAILOG_ERROR(
          LOG_S1AP, "HANDOVER REQUIRED with Invalid Cause_Type = %d\n",
          cause_type);
      OAILOG_FUNC_RETURN(LOG_S1AP, RETURNerror);
  }

  // Target ID
  S1AP_FIND_PROTOCOLIE_BY_ID(
      S1ap_HandoverRequiredIEs_t, ie, container, S1ap_ProtocolIE_ID_id_TargetID,
      true);
  if (ie) {
    if (ie->value.choice.TargetID.present == S1ap_TargetID_PR_targeteNB_ID) {
      targeteNB_ID = &ie->value.choice.TargetID.choice.targeteNB_ID;
      if (targeteNB_ID->global_ENB_ID.eNB_ID.present ==
          S1ap_ENB_ID_PR_homeENB_ID) {
        // Home eNB ID = 28 bits
        enb_id_buf = targeteNB_ID->global_ENB_ID.eNB_ID.choice.homeENB_ID.buf;

        target_enb_id = (enb_id_buf[0] << 20) + (enb_id_buf[1] << 12) +
                        (enb_id_buf[2] << 4) + ((enb_id_buf[3] & 0xf0) >> 4);
        OAILOG_INFO(LOG_S1AP, "home eNB id: %u\n", target_enb_id);
      } else {
        // Macro eNB = 20 bits
        enb_id_buf = targeteNB_ID->global_ENB_ID.eNB_ID.choice.macroENB_ID.buf;

        target_enb_id = (enb_id_buf[0] << 12) + (enb_id_buf[1] << 4) +
                        ((enb_id_buf[2] & 0xf0) >> 4);
        OAILOG_INFO(LOG_S1AP, "macro eNB id: %u\n", target_enb_id);
      }
    } else {
      OAILOG_ERROR(LOG_S1AP, "Invalid target, only intra LTE HO supported");
      OAILOG_FUNC_RETURN(LOG_S1AP, RETURNerror);
    }
  } else {
    OAILOG_FUNC_RETURN(LOG_S1AP, RETURNerror);
  }

  // Source to Target Transparent Container
  S1AP_FIND_PROTOCOLIE_BY_ID(
      S1ap_HandoverRequiredIEs_t, ie, container,
      S1ap_ProtocolIE_ID_id_Source_ToTarget_TransparentContainer, true);
  if (ie) {
    // note: ownership of src_tgt_container transferred to receiver
    src_tgt_container = blk2bstr(
        ie->value.choice.Source_ToTarget_TransparentContainer.buf,
        ie->value.choice.Source_ToTarget_TransparentContainer.size);

  } else {
    OAILOG_FUNC_RETURN(LOG_S1AP, RETURNerror);
  }

  OAILOG_INFO(
      LOG_S1AP,
      "Handover Required from association id %u, "
      "MME UE S1AP ID (" MME_UE_S1AP_ID_FMT
      ") ENB UE S1AP ID (" ENB_UE_S1AP_ID_FMT
      ") "
      "HandoverType = %ld CauseType = %u CauseValue = %ld Target ID = %u",
      assoc_id, mme_ue_s1ap_id, enb_ue_s1ap_id, handover_type, cause_type,
      cause_value, target_enb_id);

  // get imsi for logging
  hashtable_uint64_ts_get(
      imsi_map->mme_ue_id_imsi_htbl, (const hash_key_t) mme_ue_s1ap_id,
      &imsi64);

  // retrieve enb_description using hash table and match target_enb_id
  if ((enb_array = hashtable_ts_get_elements(&state->enbs)) != NULL) {
    for (idx = 0; idx < enb_array->num_elements; idx++) {
      target_enb_association =
          (enb_description_t*) (uintptr_t) enb_array->elements[idx];
      if (target_enb_association->enb_id == target_enb_id) {
        break;
      }
    }
    free_wrapper((void**) &enb_array->elements);
    free_wrapper((void**) &enb_array);
    if (target_enb_association->enb_id != target_enb_id) {
      bdestroy_wrapper(&src_tgt_container);
      OAILOG_ERROR(LOG_S1AP, "No eNB for enb_id %d\n", target_enb_id);
      OAILOG_FUNC_RETURN(LOG_S1AP, RETURNerror);
    }
  }

  OAILOG_INFO_UE(
      LOG_S1AP, imsi64, "Handing over to enb_id %d (sctp assoc %d)\n",
      target_enb_id, target_enb_association->sctp_assoc_id);

  s1ap_mme_itti_s1ap_handover_required(
      target_enb_association->sctp_assoc_id, target_enb_id, cause,
      handover_type, mme_ue_s1ap_id, src_tgt_container, imsi64);

  OAILOG_FUNC_RETURN(LOG_S1AP, RETURNok);
}

status_code_e s1ap_mme_handle_handover_command(
    s1ap_state_t* state, const itti_mme_app_handover_command_t* ho_command_p) {
  uint8_t* buffer_p   = NULL;
  uint8_t err         = 0;
  uint32_t length     = 0;
  S1ap_S1AP_PDU_t pdu = {0};
  S1ap_HandoverCommand_t* out;
  S1ap_HandoverCommandIEs_t* ie = NULL;
  ue_description_t* ue_ref_p    = NULL;
  sctp_stream_id_t stream       = 0x0;

  OAILOG_FUNC_IN(LOG_S1AP);
  if (ho_command_p == NULL) {
    OAILOG_ERROR(LOG_S1AP, "Handover Command is null\n");
    OAILOG_FUNC_RETURN(LOG_S1AP, RETURNerror);
  }

  if ((ue_ref_p = s1ap_state_get_ue_mmeid(ho_command_p->mme_ue_s1ap_id)) ==
      NULL) {
    OAILOG_ERROR(
        LOG_S1AP,
        "could not get ue context for mme_ue_s1ap_id " MME_UE_S1AP_ID_FMT
        ", failing!\n",
        (uint32_t) ho_command_p->mme_ue_s1ap_id);
    OAILOG_FUNC_RETURN(LOG_S1AP, RETURNerror);
  } else {
    stream = ue_ref_p->sctp_stream_send;
  }

  // we're doing handover, update the ue state
  ue_ref_p->s1_ue_state                        = S1AP_UE_HANDOVER;
  ue_ref_p->s1ap_handover_state.mme_ue_s1ap_id = ho_command_p->mme_ue_s1ap_id;
  ue_ref_p->s1ap_handover_state.source_enb_id  = ho_command_p->source_enb_id;
  ue_ref_p->s1ap_handover_state.target_enb_id  = ho_command_p->target_enb_id;
  ue_ref_p->s1ap_handover_state.target_enb_ue_s1ap_id =
      ho_command_p->tgt_enb_ue_s1ap_id;
  ue_ref_p->s1ap_handover_state.source_enb_ue_s1ap_id =
      ue_ref_p->enb_ue_s1ap_id;

  OAILOG_INFO(LOG_S1AP, "Handover Command received");
  pdu.present = S1ap_S1AP_PDU_PR_successfulOutcome;
  pdu.choice.successfulOutcome.procedureCode =
      S1ap_ProcedureCode_id_HandoverPreparation;
  pdu.choice.successfulOutcome.value.present =
      S1ap_SuccessfulOutcome__value_PR_HandoverCommand;
  pdu.choice.successfulOutcome.criticality = S1ap_Criticality_ignore;
  out = &pdu.choice.successfulOutcome.value.choice.HandoverCommand;

  /* MME-UE-ID: mandatory */
  ie =
      (S1ap_HandoverCommandIEs_t*) calloc(1, sizeof(S1ap_HandoverCommandIEs_t));
  ie->id            = S1ap_ProtocolIE_ID_id_MME_UE_S1AP_ID;
  ie->criticality   = S1ap_Criticality_reject;
  ie->value.present = S1ap_HandoverCommandIEs__value_PR_MME_UE_S1AP_ID;
  ie->value.choice.MME_UE_S1AP_ID = ho_command_p->mme_ue_s1ap_id;
  ASN_SEQUENCE_ADD(&out->protocolIEs.list, ie);

  /* eNB-UE-ID: mandatory */
  ie =
      (S1ap_HandoverCommandIEs_t*) calloc(1, sizeof(S1ap_HandoverCommandIEs_t));
  ie->id            = S1ap_ProtocolIE_ID_id_eNB_UE_S1AP_ID;
  ie->criticality   = S1ap_Criticality_reject;
  ie->value.present = S1ap_HandoverCommandIEs__value_PR_ENB_UE_S1AP_ID;
  ie->value.choice.ENB_UE_S1AP_ID = ho_command_p->src_enb_ue_s1ap_id;
  ASN_SEQUENCE_ADD(&out->protocolIEs.list, ie);

  /* HandoverType: mandatory */
  ie =
      (S1ap_HandoverCommandIEs_t*) calloc(1, sizeof(S1ap_HandoverCommandIEs_t));
  ie->id            = S1ap_ProtocolIE_ID_id_HandoverType;
  ie->criticality   = S1ap_Criticality_reject;
  ie->value.present = S1ap_HandoverCommandIEs__value_PR_HandoverType;
  ie->value.choice.HandoverType = ho_command_p->handover_type;
  ASN_SEQUENCE_ADD(&out->protocolIEs.list, ie);

  /* Target-ToSource-TransparentContainer: mandatory */
  ie =
      (S1ap_HandoverCommandIEs_t*) calloc(1, sizeof(S1ap_HandoverRequestIEs_t));
  ie->id          = S1ap_ProtocolIE_ID_id_Target_ToSource_TransparentContainer;
  ie->criticality = S1ap_Criticality_reject;
  ie->value.present =
      S1ap_HandoverCommandIEs__value_PR_Target_ToSource_TransparentContainer;
  OCTET_STRING_fromBuf(
      &ie->value.choice.Target_ToSource_TransparentContainer,
      (char*) bdata(ho_command_p->tgt_src_container),
      blength(ho_command_p->tgt_src_container));
  ASN_SEQUENCE_ADD(&out->protocolIEs.list, ie);

  // Construct the PDU and send message
  if (s1ap_mme_encode_pdu(&pdu, &buffer_p, &length) < 0) {
    err = 1;
  }
  ASN_STRUCT_FREE_CONTENTS_ONLY(asn_DEF_S1ap_HandoverRequest, out);
  if (err) {
    OAILOG_FUNC_RETURN(LOG_S1AP, RETURNerror);
  }
  bstring b = blk2bstr(buffer_p, length);
  free(buffer_p);

  s1ap_mme_itti_send_sctp_request(
      &b, ho_command_p->source_assoc_id, stream, ho_command_p->mme_ue_s1ap_id);

  OAILOG_FUNC_RETURN(LOG_S1AP, RETURNok);
}

status_code_e s1ap_mme_handle_handover_notify(
    s1ap_state_t* state, const sctp_assoc_id_t assoc_id,
    const sctp_stream_id_t stream, S1ap_S1AP_PDU_t* pdu) {
  S1ap_HandoverNotify_t* container    = NULL;
  S1ap_HandoverNotifyIEs_t* ie        = NULL;
  enb_description_t* target_enb       = NULL;
  ue_description_t* src_ue_ref_p      = NULL;
  ue_description_t* new_ue_ref_p      = NULL;
  mme_ue_s1ap_id_t mme_ue_s1ap_id     = INVALID_MME_UE_S1AP_ID;
  enb_ue_s1ap_id_t tgt_enb_ue_s1ap_id = INVALID_ENB_UE_S1AP_ID;
  ecgi_t ecgi                         = {.plmn = {0}, .cell_identity = {0}};
  tai_t tai                           = {0};
  imsi64_t imsi64                     = INVALID_IMSI64;
  s1ap_imsi_map_t* imsi_map           = get_s1ap_imsi_map();

  OAILOG_FUNC_IN(LOG_S1AP);

  target_enb = s1ap_state_get_enb(state, assoc_id);
  if (target_enb == NULL) {
    OAILOG_ERROR(
        LOG_S1AP, "Ignore HandoverNotify from unknown assoc %u\n", assoc_id);
    OAILOG_FUNC_RETURN(LOG_S1AP, RETURNerror);
  }

  container = &pdu->choice.initiatingMessage.value.choice.HandoverNotify;

  // HandoverNotify means the handover has completed successfully. We can
  // remove the UE context from the old eNB, tear down indirect forwarding
  // tunnels, modify the DL bearer, and create the new UE context on the new
  // eNB.

  // get the mandantory IEs
  // MME_UE_S1AP_ID
  S1AP_FIND_PROTOCOLIE_BY_ID(
      S1ap_HandoverNotifyIEs_t, ie, container,
      S1ap_ProtocolIE_ID_id_MME_UE_S1AP_ID, true);
  if (ie) {
    mme_ue_s1ap_id = ie->value.choice.MME_UE_S1AP_ID;
  } else {
    OAILOG_FUNC_RETURN(LOG_S1AP, RETURNerror);
  }

  // eNB_UE_S1AP_ID
  S1AP_FIND_PROTOCOLIE_BY_ID(
      S1ap_HandoverNotifyIEs_t, ie, container,
      S1ap_ProtocolIE_ID_id_eNB_UE_S1AP_ID, true);
  // eNB UE S1AP ID is limited to 24 bits
  if (ie) {
    tgt_enb_ue_s1ap_id = (enb_ue_s1ap_id_t)(
        ie->value.choice.ENB_UE_S1AP_ID & ENB_UE_S1AP_ID_MASK);
  } else {
    OAILOG_FUNC_RETURN(LOG_S1AP, RETURNerror);
  }

  // CGI
  S1AP_FIND_PROTOCOLIE_BY_ID(
      S1ap_HandoverNotifyIEs_t, ie, container, S1ap_ProtocolIE_ID_id_EUTRAN_CGI,
      true);

  if (!ie) {
    OAILOG_ERROR(LOG_S1AP, "Incorrect IE \n");
    return RETURNerror;
  }

  if (!(ie->value.choice.EUTRAN_CGI.pLMNidentity.size == 3)) {
    OAILOG_ERROR(LOG_S1AP, "Incorrect PLMN size \n");
    return RETURNerror;
  }
  TBCD_TO_PLMN_T(&ie->value.choice.EUTRAN_CGI.pLMNidentity, &ecgi.plmn);
  BIT_STRING_TO_CELL_IDENTITY(
      &ie->value.choice.EUTRAN_CGI.cell_ID, ecgi.cell_identity);

  // TAI
  S1AP_FIND_PROTOCOLIE_BY_ID(
      S1ap_HandoverNotifyIEs_t, ie, container, S1ap_ProtocolIE_ID_id_TAI, true);
  if (!ie) {
    OAILOG_ERROR(LOG_S1AP, "Incorrect IE \n");
    return RETURNerror;
  }
  OCTET_STRING_TO_TAC(&ie->value.choice.TAI.tAC, tai.tac);
  if (!(ie->value.choice.EUTRAN_CGI.pLMNidentity.size == 3)) {
    OAILOG_ERROR(LOG_S1AP, "Incorrect PLMN size \n");
    return RETURNerror;
  }
  TBCD_TO_PLMN_T(&ie->value.choice.TAI.pLMNidentity, &tai.plmn);

  // imsi for logging
  hashtable_uint64_ts_get(
      imsi_map->mme_ue_id_imsi_htbl, (const hash_key_t) mme_ue_s1ap_id,
      &imsi64);

  // get existing UE context
  if ((src_ue_ref_p = s1ap_state_get_ue_mmeid(mme_ue_s1ap_id)) == NULL) {
    OAILOG_ERROR_UE(
        LOG_S1AP, imsi64,
        "source MME_UE_S1AP_ID (" MME_UE_S1AP_ID_FMT
        ") does not point to any valid UE\n",
        mme_ue_s1ap_id);
    OAILOG_FUNC_RETURN(LOG_S1AP, RETURNerror);
  } else {
    // create new UE context, remove the old one.
    new_ue_ref_p =
        s1ap_state_get_ue_enbid(target_enb->sctp_assoc_id, tgt_enb_ue_s1ap_id);
    if (new_ue_ref_p != NULL) {
      OAILOG_ERROR_UE(
          LOG_S1AP, imsi64,
          "S1AP:Handover Notify- Received ENB_UE_S1AP_ID is not Unique "
          "Drop Handover Notify for eNBUeS1APId:" ENB_UE_S1AP_ID_FMT "\n",
          tgt_enb_ue_s1ap_id);
      OAILOG_FUNC_RETURN(LOG_S1AP, RETURNerror);
    }
    if ((new_ue_ref_p = s1ap_new_ue(state, assoc_id, tgt_enb_ue_s1ap_id)) ==
        NULL) {
      // If we failed to allocate a new UE return -1
      OAILOG_ERROR_UE(
          LOG_S1AP, imsi64,
          "S1AP:Handover Notify- Failed to allocate S1AP UE Context, "
          "eNBUeS1APId:" ENB_UE_S1AP_ID_FMT "\n",
          tgt_enb_ue_s1ap_id);
      OAILOG_FUNC_RETURN(LOG_S1AP, RETURNerror);
    }
    new_ue_ref_p->s1_ue_state    = S1AP_UE_CONNECTED;  // handover has completed
    new_ue_ref_p->enb_ue_s1ap_id = tgt_enb_ue_s1ap_id;
    // Will be allocated by NAS
    new_ue_ref_p->mme_ue_s1ap_id = mme_ue_s1ap_id;

    new_ue_ref_p->s1ap_ue_context_rel_timer.id =
        src_ue_ref_p->s1ap_ue_context_rel_timer.id;
    new_ue_ref_p->s1ap_ue_context_rel_timer.msec =
        src_ue_ref_p->s1ap_ue_context_rel_timer.msec;
    new_ue_ref_p->sctp_stream_recv =
        src_ue_ref_p->s1ap_handover_state.target_sctp_stream_recv;
    new_ue_ref_p->sctp_stream_send =
        src_ue_ref_p->s1ap_handover_state.target_sctp_stream_send;
    new_ue_ref_p->s1ap_handover_state = src_ue_ref_p->s1ap_handover_state;

    // generate a message to update bearers
    s1ap_mme_itti_s1ap_handover_notify(
        mme_ue_s1ap_id, src_ue_ref_p->s1ap_handover_state, tgt_enb_ue_s1ap_id,
        assoc_id, ecgi, imsi64);

    /* Remove ue description from source eNB */
    s1ap_remove_ue(state, src_ue_ref_p);

    /* Mapping between mme_ue_s1ap_id, assoc_id and enb_ue_s1ap_id */
    hashtable_rc_t h_rc = hashtable_ts_insert(
        &state->mmeid2associd, (const hash_key_t) new_ue_ref_p->mme_ue_s1ap_id,
        (void*) (uintptr_t) assoc_id);

    hashtable_uint64_ts_insert(
        &target_enb->ue_id_coll,
        (const hash_key_t) new_ue_ref_p->mme_ue_s1ap_id,
        new_ue_ref_p->comp_s1ap_id);

    OAILOG_DEBUG_UE(
        LOG_S1AP, imsi64,
        "Associated sctp_assoc_id %d, enb_ue_s1ap_id " ENB_UE_S1AP_ID_FMT
        ", mme_ue_s1ap_id " MME_UE_S1AP_ID_FMT ":%s \n",
        assoc_id, new_ue_ref_p->enb_ue_s1ap_id, new_ue_ref_p->mme_ue_s1ap_id,
        hashtable_rc_code2string(h_rc));

    s1ap_dump_enb(target_enb);
  }

  OAILOG_FUNC_RETURN(LOG_S1AP, RETURNok);
}

status_code_e s1ap_mme_handle_enb_status_transfer(
    s1ap_state_t* state, const sctp_assoc_id_t assoc_id,
    const sctp_stream_id_t stream, S1ap_S1AP_PDU_t* pdu) {
  S1ap_ENBStatusTransfer_t* container       = NULL;
  S1ap_ENBStatusTransferIEs_t* ie           = NULL;
  ue_description_t* ue_ref_p                = NULL;
  mme_ue_s1ap_id_t mme_ue_s1ap_id           = INVALID_MME_UE_S1AP_ID;
  hashtable_element_array_t* enb_array      = NULL;
  enb_description_t* target_enb_association = NULL;
  uint8_t* buffer                           = NULL;
  uint32_t length                           = 0;
  uint32_t idx                              = 0;

  OAILOG_FUNC_IN(LOG_S1AP);
  container = &pdu->choice.initiatingMessage.value.choice.ENBStatusTransfer;

  // similar to enb_configuration_transfer, we immediately generate the new
  // message by changing type and updating the enb_ue_s1ap_id to match that of
  // the target enb.

  // MME_UE_S1AP_ID
  S1AP_FIND_PROTOCOLIE_BY_ID(
      S1ap_ENBStatusTransferIEs_t, ie, container,
      S1ap_ProtocolIE_ID_id_MME_UE_S1AP_ID, true);
  if (ie) {
    mme_ue_s1ap_id = ie->value.choice.MME_UE_S1AP_ID;
  } else {
    OAILOG_FUNC_RETURN(LOG_S1AP, RETURNerror);
  }

  // get the UE and handover state
  if ((ue_ref_p = s1ap_state_get_ue_mmeid(mme_ue_s1ap_id)) == NULL) {
    OAILOG_ERROR(
        LOG_S1AP,
        "could not get ue context for mme_ue_s1ap_id " MME_UE_S1AP_ID_FMT
        ", failing!\n",
        (uint32_t) mme_ue_s1ap_id);
    OAILOG_FUNC_RETURN(LOG_S1AP, RETURNerror);
  }

  OAILOG_INFO(
      LOG_S1AP,
      "Received eNBStatusTransfer from source enb_id assoc %u for "
      "ue " MME_UE_S1AP_ID_FMT " to target enb_id %u\n",
      assoc_id, mme_ue_s1ap_id,
      ue_ref_p->s1ap_handover_state.target_enb_ue_s1ap_id);

  // set the target eNB_UE_S1AP_ID
  S1AP_FIND_PROTOCOLIE_BY_ID(
      S1ap_ENBStatusTransferIEs_t, ie, container,
      S1ap_ProtocolIE_ID_id_eNB_UE_S1AP_ID, true);
  if (ie) {
    ie->value.choice.ENB_UE_S1AP_ID =
        ue_ref_p->s1ap_handover_state.target_enb_ue_s1ap_id;
  } else {
    OAILOG_FUNC_RETURN(LOG_S1AP, RETURNerror);
  }

  // get the enb_description matching the target_enb_id
  // retrieve enb_description using hash table and match target_enb_id
  if ((enb_array = hashtable_ts_get_elements(&state->enbs)) != NULL) {
    for (idx = 0; idx < enb_array->num_elements; idx++) {
      target_enb_association =
          (enb_description_t*) (uintptr_t) enb_array->elements[idx];
      if (target_enb_association->enb_id ==
          ue_ref_p->s1ap_handover_state.target_enb_id) {
        break;
      }
    }
    free_wrapper((void**) &enb_array->elements);
    free_wrapper((void**) &enb_array);
    if (target_enb_association->enb_id !=
        ue_ref_p->s1ap_handover_state.target_enb_id) {
      OAILOG_ERROR(
          LOG_S1AP, "No eNB for enb_id %d\n",
          ue_ref_p->s1ap_handover_state.target_enb_id);
      OAILOG_FUNC_RETURN(LOG_S1AP, RETURNerror);
    }
  }

  // change the message type and enb_ue_s1_id to the target eNB's ID
  pdu->choice.initiatingMessage.procedureCode =
      S1ap_ProcedureCode_id_MMEStatusTransfer;
  pdu->present = S1ap_S1AP_PDU_PR_initiatingMessage;

  // Encode message
  if (s1ap_mme_encode_pdu(pdu, &buffer, &length) < 0) {
    OAILOG_ERROR(
        LOG_S1AP,
        "Failed to encode MME Configuration Transfer message for enb_id %u\n",
        ue_ref_p->s1ap_handover_state.target_enb_id);
    OAILOG_FUNC_RETURN(LOG_S1AP, RETURNerror);
  }

  bstring b = blk2bstr(buffer, length);
  free(buffer);

  s1ap_mme_itti_send_sctp_request(
      &b, target_enb_association->sctp_assoc_id,
      ue_ref_p->s1ap_handover_state.target_sctp_stream_recv,
      ue_ref_p->mme_ue_s1ap_id);

  OAILOG_FUNC_RETURN(LOG_S1AP, RETURNok);
}

status_code_e s1ap_mme_handle_path_switch_request(
    s1ap_state_t* state, __attribute__((unused)) const sctp_assoc_id_t assoc_id,
    __attribute__((unused)) const sctp_stream_id_t stream,
    S1ap_S1AP_PDU_t* pdu) {
  S1ap_PathSwitchRequest_t* container                            = NULL;
  S1ap_PathSwitchRequestIEs_t* ie                                = NULL;
  S1ap_E_RABToBeSwitchedDLItemIEs_t* eRABToBeSwitchedDlItemIEs_p = NULL;
  enb_description_t* enb_association                             = NULL;
  ue_description_t* ue_ref_p                                     = NULL;
  ue_description_t* new_ue_ref_p                                 = NULL;
  mme_ue_s1ap_id_t mme_ue_s1ap_id = INVALID_MME_UE_S1AP_ID;
  enb_ue_s1ap_id_t enb_ue_s1ap_id = INVALID_ENB_UE_S1AP_ID;
  ecgi_t ecgi                     = {.plmn = {0}, .cell_identity = {0}};
  tai_t tai                       = {0};
  uint16_t encryption_algorithm_capabilities                           = 0;
  uint16_t integrity_algorithm_capabilities                            = 0;
  e_rab_to_be_switched_in_downlink_list_t e_rab_to_be_switched_dl_list = {0};
  uint32_t num_erab                                                    = 0;
  uint16_t index                                                       = 0;
  itti_s1ap_path_switch_request_failure_t path_switch_req_failure      = {0};
  imsi64_t imsi64           = INVALID_IMSI64;
  s1ap_imsi_map_t* imsi_map = get_s1ap_imsi_map();

  OAILOG_FUNC_IN(LOG_S1AP);

  enb_association = s1ap_state_get_enb(state, assoc_id);
  if (enb_association == NULL) {
    OAILOG_ERROR(
        LOG_S1AP,
        "Ignore Path Switch Request from unknown assoc "
        "%u\n",
        assoc_id);
    OAILOG_FUNC_RETURN(LOG_S1AP, RETURNerror);
  }

  container = &pdu->choice.initiatingMessage.value.choice.PathSwitchRequest;

  S1AP_FIND_PROTOCOLIE_BY_ID(
      S1ap_PathSwitchRequestIEs_t, ie, container,
      S1ap_ProtocolIE_ID_id_SourceMME_UE_S1AP_ID, true);
  if (ie) {
    mme_ue_s1ap_id = ie->value.choice.MME_UE_S1AP_ID;
  } else {
    OAILOG_FUNC_RETURN(LOG_S1AP, RETURNerror);
  }

  S1AP_FIND_PROTOCOLIE_BY_ID(
      S1ap_PathSwitchRequestIEs_t, ie, container,
      S1ap_ProtocolIE_ID_id_eNB_UE_S1AP_ID, true);
  // eNB UE S1AP ID is limited to 24 bits
  if (ie) {
    enb_ue_s1ap_id = (enb_ue_s1ap_id_t)(
        ie->value.choice.ENB_UE_S1AP_ID & ENB_UE_S1AP_ID_MASK);
  } else {
    OAILOG_FUNC_RETURN(LOG_S1AP, RETURNerror);
  }

  OAILOG_DEBUG_UE(
      LOG_S1AP, imsi64,
      "Path Switch Request message received from eNB UE S1AP "
      "ID: " ENB_UE_S1AP_ID_FMT "\n",
      enb_ue_s1ap_id);

  hashtable_uint64_ts_get(
      imsi_map->mme_ue_id_imsi_htbl, (const hash_key_t) mme_ue_s1ap_id,
      &imsi64);

  /* If all the E-RAB ID IEs in E-RABToBeSwitchedDLList is set to the
   * same value, send PATH SWITCH REQUEST FAILURE message to eNB */
  if (true == is_all_erabId_same(container)) {
    /*send PATH SWITCH REQUEST FAILURE message to eNB*/
    path_switch_req_failure.sctp_assoc_id  = assoc_id;
    path_switch_req_failure.mme_ue_s1ap_id = mme_ue_s1ap_id;
    path_switch_req_failure.enb_ue_s1ap_id = enb_ue_s1ap_id;
    s1ap_handle_path_switch_req_failure(&path_switch_req_failure, imsi64);
    OAILOG_FUNC_RETURN(LOG_S1AP, RETURNerror);
  }

  if ((ue_ref_p = s1ap_state_get_ue_mmeid(mme_ue_s1ap_id)) == NULL) {
    /*
     * The MME UE S1AP ID provided by eNB doesn't point to any valid UE.
     * MME ignore this PATH SWITCH REQUEST.
     */
    OAILOG_ERROR_UE(
        LOG_S1AP, imsi64,
        "source MME_UE_S1AP_ID (" MME_UE_S1AP_ID_FMT
        ") does not point to any valid UE\n",
        mme_ue_s1ap_id);
    OAILOG_FUNC_RETURN(LOG_S1AP, RETURNerror);
  } else {
    new_ue_ref_p =
        s1ap_state_get_ue_enbid(enb_association->sctp_assoc_id, enb_ue_s1ap_id);
    if (new_ue_ref_p != NULL) {
      OAILOG_ERROR_UE(
          LOG_S1AP, imsi64,
          "S1AP:Path Switch Request- Recieved ENB_UE_S1AP_ID is not Unique "
          "Drop Path Switch Request for eNBUeS1APId:" ENB_UE_S1AP_ID_FMT "\n",
          enb_ue_s1ap_id);
      OAILOG_FUNC_RETURN(LOG_S1AP, RETURNerror);
    }
    /*
     * Creat New UE Context with target eNB and delete Old UE Context
     * from source eNB.
     */
    if ((new_ue_ref_p = s1ap_new_ue(state, assoc_id, enb_ue_s1ap_id)) == NULL) {
      // If we failed to allocate a new UE return -1
      OAILOG_ERROR_UE(
          LOG_S1AP, imsi64,
          "S1AP:Path Switch Request- Failed to allocate S1AP UE Context, "
          "eNBUeS1APId:" ENB_UE_S1AP_ID_FMT "\n",
          enb_ue_s1ap_id);
      OAILOG_FUNC_RETURN(LOG_S1AP, RETURNerror);
    }
    new_ue_ref_p->s1_ue_state    = ue_ref_p->s1_ue_state;
    new_ue_ref_p->enb_ue_s1ap_id = enb_ue_s1ap_id;
    // Will be allocated by NAS
    new_ue_ref_p->mme_ue_s1ap_id = mme_ue_s1ap_id;

    new_ue_ref_p->s1ap_ue_context_rel_timer.id =
        ue_ref_p->s1ap_ue_context_rel_timer.id;
    new_ue_ref_p->s1ap_ue_context_rel_timer.msec =
        ue_ref_p->s1ap_ue_context_rel_timer.msec;
    // On which stream we received the message
    new_ue_ref_p->sctp_stream_recv = stream;
    new_ue_ref_p->sctp_stream_send = enb_association->next_sctp_stream;
    enb_association->next_sctp_stream += 1;
    if (enb_association->next_sctp_stream >= enb_association->instreams) {
      enb_association->next_sctp_stream = 1;
    }
    /* Remove ue description from source eNB */
    s1ap_remove_ue(state, ue_ref_p);

    /* Mapping between mme_ue_s1ap_id, assoc_id and enb_ue_s1ap_id */
    hashtable_rc_t h_rc = hashtable_ts_insert(
        &state->mmeid2associd, (const hash_key_t) new_ue_ref_p->mme_ue_s1ap_id,
        (void*) (uintptr_t) assoc_id);

    hashtable_uint64_ts_insert(
        &enb_association->ue_id_coll,
        (const hash_key_t) new_ue_ref_p->mme_ue_s1ap_id,
        new_ue_ref_p->comp_s1ap_id);

    OAILOG_DEBUG_UE(
        LOG_S1AP, imsi64,
        "Associated sctp_assoc_id %d, enb_ue_s1ap_id " ENB_UE_S1AP_ID_FMT
        ", mme_ue_s1ap_id " MME_UE_S1AP_ID_FMT ":%s \n",
        assoc_id, new_ue_ref_p->enb_ue_s1ap_id, new_ue_ref_p->mme_ue_s1ap_id,
        hashtable_rc_code2string(h_rc));

    s1ap_dump_enb(enb_association);

    S1AP_FIND_PROTOCOLIE_BY_ID(
        S1ap_PathSwitchRequestIEs_t, ie, container,
        S1ap_ProtocolIE_ID_id_E_RABToBeSwitchedDLList, true);

    if (!ie) {
      OAILOG_ERROR(LOG_S1AP, "Incorrect IE \n");
      return RETURNerror;
    }

    S1ap_E_RABToBeSwitchedDLList_t* e_rab_to_be_switched_dl_list_req =
        &ie->value.choice.E_RABToBeSwitchedDLList;

    // E-RAB To Be Switched in Downlink List mandatory IE
    num_erab = e_rab_to_be_switched_dl_list_req->list.count;
    for (index = 0; index < num_erab; ++index) {
      eRABToBeSwitchedDlItemIEs_p =
          (S1ap_E_RABToBeSwitchedDLItemIEs_t*)
              e_rab_to_be_switched_dl_list_req->list.array[index];
      S1ap_E_RABToBeSwitchedDLItem_t* eRab_ToBeSwitchedDLItem =
          &eRABToBeSwitchedDlItemIEs_p->value.choice.E_RABToBeSwitchedDLItem;

      e_rab_to_be_switched_dl_list.item[index].e_rab_id =
          eRab_ToBeSwitchedDLItem->e_RAB_ID;
      e_rab_to_be_switched_dl_list.item[index].transport_layer_address =
          blk2bstr(
              eRab_ToBeSwitchedDLItem->transportLayerAddress.buf,
              eRab_ToBeSwitchedDLItem->transportLayerAddress.size);
      e_rab_to_be_switched_dl_list.item[index].gtp_teid =
          htonl(*((uint32_t*) eRab_ToBeSwitchedDLItem->gTP_TEID.buf));
      e_rab_to_be_switched_dl_list.no_of_items += 1;
    }

    // CGI mandatory IE
    S1AP_FIND_PROTOCOLIE_BY_ID(
        S1ap_PathSwitchRequestIEs_t, ie, container,
        S1ap_ProtocolIE_ID_id_EUTRAN_CGI, true);

    if (!ie) {
      OAILOG_ERROR(LOG_S1AP, "Incorrect IE \n");
      return RETURNerror;
    }

    if (!(ie->value.choice.EUTRAN_CGI.pLMNidentity.size == 3)) {
      OAILOG_ERROR(LOG_S1AP, "Incorrect PLMN size \n");
      return RETURNerror;
    }
    TBCD_TO_PLMN_T(&ie->value.choice.EUTRAN_CGI.pLMNidentity, &ecgi.plmn);
    BIT_STRING_TO_CELL_IDENTITY(
        &ie->value.choice.EUTRAN_CGI.cell_ID, ecgi.cell_identity);

    // TAI mandatory IE
    S1AP_FIND_PROTOCOLIE_BY_ID(
        S1ap_PathSwitchRequestIEs_t, ie, container, S1ap_ProtocolIE_ID_id_TAI,
        true);
    if (!ie) {
      OAILOG_ERROR(LOG_S1AP, "Incorrect IE \n");
      return RETURNerror;
    }
    OCTET_STRING_TO_TAC(&ie->value.choice.TAI.tAC, tai.tac);
    if (!(ie->value.choice.EUTRAN_CGI.pLMNidentity.size == 3)) {
      OAILOG_ERROR(LOG_S1AP, "Incorrect PLMN size \n");
      return RETURNerror;
    }
    TBCD_TO_PLMN_T(&ie->value.choice.TAI.pLMNidentity, &tai.plmn);

    // UE Security Capabilities mandatory IE
    S1AP_FIND_PROTOCOLIE_BY_ID(
        S1ap_PathSwitchRequestIEs_t, ie, container,
        S1ap_ProtocolIE_ID_id_UESecurityCapabilities, true);
    BIT_STRING_TO_INT16(
        &ie->value.choice.UESecurityCapabilities.encryptionAlgorithms,
        encryption_algorithm_capabilities);
    BIT_STRING_TO_INT16(
        &ie->value.choice.UESecurityCapabilities.integrityProtectionAlgorithms,
        integrity_algorithm_capabilities);
  }

  s1ap_mme_itti_s1ap_path_switch_request(
      assoc_id, enb_association->enb_id, new_ue_ref_p->enb_ue_s1ap_id,
      &e_rab_to_be_switched_dl_list, new_ue_ref_p->mme_ue_s1ap_id, &ecgi, &tai,
      encryption_algorithm_capabilities, integrity_algorithm_capabilities,
      imsi64);

  OAILOG_FUNC_RETURN(LOG_S1AP, RETURNok);
}

//------------------------------------------------------------------------------
typedef struct arg_s1ap_send_enb_dereg_ind_s {
  uint8_t current_ue_index;
  uint32_t handled_ues;
  MessageDef* message_p;
  uint32_t associated_enb_id;
  uint32_t deregister_ue_count;
} arg_s1ap_send_enb_dereg_ind_t;

//------------------------------------------------------------------------------
bool s1ap_send_enb_deregistered_ind(
    __attribute__((unused)) const hash_key_t keyP, uint64_t const dataP,
    void* argP, void** resultP) {
  arg_s1ap_send_enb_dereg_ind_t* arg = (arg_s1ap_send_enb_dereg_ind_t*) argP;
  ue_description_t* ue_ref_p         = NULL;

  // Ask for the release of each UE context associated to the eNB
  hash_table_ts_t* s1ap_ue_state = get_s1ap_ue_state();
  hashtable_ts_get(s1ap_ue_state, (const hash_key_t) dataP, (void**) &ue_ref_p);
  if (ue_ref_p) {
    if (arg->current_ue_index == 0) {
      arg->message_p = DEPRECATEDitti_alloc_new_message_fatal(
          TASK_S1AP, S1AP_ENB_DEREGISTERED_IND);
      OAILOG_DEBUG(LOG_S1AP, "eNB Deregesteration");
    }
    if (ue_ref_p->mme_ue_s1ap_id == INVALID_MME_UE_S1AP_ID) {
      /*
       * Send deregistered ind for this also and let MMEAPP find the context
       * using enb_ue_s1ap_id_key
       */
      OAILOG_WARNING(LOG_S1AP, "UE with invalid MME s1ap id found");
    }

    AssertFatal(
        arg->current_ue_index < S1AP_ITTI_UE_PER_DEREGISTER_MESSAGE,
        "Too many deregistered UEs reported in S1AP_ENB_DEREGISTERED_IND "
        "message ");
    S1AP_ENB_DEREGISTERED_IND(arg->message_p)
        .mme_ue_s1ap_id[arg->current_ue_index] = ue_ref_p->mme_ue_s1ap_id;
    S1AP_ENB_DEREGISTERED_IND(arg->message_p)
        .enb_ue_s1ap_id[arg->current_ue_index] = ue_ref_p->enb_ue_s1ap_id;

    arg->handled_ues++;
    arg->current_ue_index++;

    if (arg->handled_ues == arg->deregister_ue_count ||
        arg->current_ue_index == S1AP_ITTI_UE_PER_DEREGISTER_MESSAGE) {
      // Sending INVALID_IMSI64 because message is not specific to any UE/IMSI
      arg->message_p->ittiMsgHeader.imsi               = INVALID_IMSI64;
      S1AP_ENB_DEREGISTERED_IND(arg->message_p).enb_id = arg->associated_enb_id;
      S1AP_ENB_DEREGISTERED_IND(arg->message_p).nb_ue_to_deregister =
          (uint8_t) arg->current_ue_index;

      // Max UEs reached for this ITTI message, send message to MME App
      OAILOG_DEBUG(
          LOG_S1AP,
          "Reached maximum UE count for this ITTI message. Sending "
          "deregistered indication to MME App for UE count = %u\n",
          S1AP_ENB_DEREGISTERED_IND(arg->message_p).nb_ue_to_deregister);

      if (arg->current_ue_index == S1AP_ITTI_UE_PER_DEREGISTER_MESSAGE) {
        arg->current_ue_index = 0;
      }
      send_msg_to_task(&s1ap_task_zmq_ctx, TASK_MME_APP, arg->message_p);
      arg->message_p = NULL;
    }

    *resultP = arg->message_p;
  } else {
    OAILOG_TRACE(LOG_S1AP, "No valid UE provided in callback: %p\n", ue_ref_p);
  }
  return false;
}

typedef struct arg_s1ap_construct_enb_reset_req_s {
  uint8_t current_ue_index;
  MessageDef* msg;
} arg_s1ap_construct_enb_reset_req_t;

bool construct_s1ap_mme_full_reset_req(
    const hash_key_t keyP, const uint64_t dataP, void* argP, void** resultP) {
  arg_s1ap_construct_enb_reset_req_t* arg = argP;
  ue_description_t* ue_ref                = (ue_description_t*) dataP;

  hash_table_ts_t* s1ap_ue_state = get_s1ap_ue_state();
  hashtable_ts_get(s1ap_ue_state, (const hash_key_t) dataP, (void**) &ue_ref);
  uint32_t i = arg->current_ue_index;
  if (ue_ref) {
    S1AP_ENB_INITIATED_RESET_REQ(arg->msg).ue_to_reset_list[i].mme_ue_s1ap_id =
        ue_ref->mme_ue_s1ap_id;
    S1AP_ENB_INITIATED_RESET_REQ(arg->msg).ue_to_reset_list[i].enb_ue_s1ap_id =
        ue_ref->enb_ue_s1ap_id;
  } else {
    OAILOG_TRACE(LOG_S1AP, "No valid UE provided in callback: %p\n", ue_ref);
    S1AP_ENB_INITIATED_RESET_REQ(arg->msg).ue_to_reset_list[i].mme_ue_s1ap_id =
        INVALID_MME_UE_S1AP_ID;
  }
  arg->current_ue_index++;

  return false;
}

//------------------------------------------------------------------------------
status_code_e s1ap_handle_sctp_disconnection(
    s1ap_state_t* state, const sctp_assoc_id_t assoc_id, bool reset) {
  arg_s1ap_send_enb_dereg_ind_t arg  = {0};
  MessageDef* message_p              = NULL;
  enb_description_t* enb_association = NULL;

  OAILOG_FUNC_IN(LOG_S1AP);

  // Checking if the assoc id has a valid eNB attached to it
  enb_association = s1ap_state_get_enb(state, assoc_id);
  if (enb_association == NULL) {
    OAILOG_ERROR(LOG_S1AP, "No eNB attached to this assoc_id: %d\n", assoc_id);
    OAILOG_FUNC_RETURN(LOG_S1AP, RETURNerror);
  }

  OAILOG_INFO(
      LOG_S1AP,
      "SCTP disconnection request for association id %u, Reset Flag = "
      "%u. Connected UEs = %u Num elements = %zu\n",
      assoc_id, reset, enb_association->nb_ue_associated,
      enb_association->ue_id_coll.num_elements);

  // First check if we can just reset the eNB state if there are no UEs
  if (!enb_association->nb_ue_associated) {
    if (reset) {
      OAILOG_INFO(
          LOG_S1AP,
          "SCTP reset request for association id %u. No Connected UEs. "
          "Reset Flag = %u\n",
          assoc_id, reset);

      OAILOG_INFO(
          LOG_S1AP, "Moving eNB with assoc_id %u to INIT state\n", assoc_id);
      enb_association->s1_state = S1AP_INIT;
      state->num_enbs--;
    } else {
      OAILOG_INFO(
          LOG_S1AP,
          "SCTP Shutdown request for association id %u. No Connected UEs. "
          "Reset Flag = %u\n",
          assoc_id, reset);

      OAILOG_INFO(LOG_S1AP, "Removing eNB with association id %u \n", assoc_id);
      s1ap_remove_enb(state, enb_association);
    }
    OAILOG_FUNC_RETURN(LOG_S1AP, RETURNok);
  }

  if (reset) {
    // Check if the UE counters for eNB are equal.
    // If not, the eNB will never switch to INIT state, particularly in
    // stateless mode.
    // Exit the process so that health checker can clean-up all Redis
    // state and restart all stateless services.
    AssertFatal(
        enb_association->nb_ue_associated ==
            enb_association->ue_id_coll.num_elements,
        "Num UEs associated with eNB (%u) is more than the UEs with valid "
        "mme_ue_s1ap_id (%zu). This is a deadlock state potentially caused by "
        "misbehaving eNB; restarting MME. In stateless mode, health management "
        "service will eventually detect multiple MME restarts due to this "
        "deadlock state and force sctpd and hence all services to restart.",
        enb_association->nb_ue_associated,
        enb_association->ue_id_coll.num_elements);
  }
  /*
   * Send S1ap deregister indication to MME app in batches of UEs where
   * UE count in each batch <= S1AP_ITTI_UE_PER_DEREGISTER_MESSAGE
   */

  arg.associated_enb_id   = enb_association->enb_id;
  arg.deregister_ue_count = enb_association->ue_id_coll.num_elements;
  hashtable_uint64_ts_apply_callback_on_elements(
      &enb_association->ue_id_coll, s1ap_send_enb_deregistered_ind,
      (void*) &arg, (void**) &message_p);

  /*
   * Mark the eNB's s1 state as appropriate, the eNB will be deleted or
   * moved to init state when the last UE's s1 state is cleaned up
   */
  enb_association->s1_state = reset ? S1AP_RESETING : S1AP_SHUTDOWN;
  OAILOG_INFO(
      LOG_S1AP, "Marked enb s1 status to %s, attached to assoc_id: %d\n",
      reset ? "Reset" : "Shutdown", assoc_id);

  OAILOG_FUNC_RETURN(LOG_S1AP, RETURNok);
}

//------------------------------------------------------------------------------
status_code_e s1ap_handle_new_association(
    s1ap_state_t* state, sctp_new_peer_t* sctp_new_peer_p) {
  enb_description_t* enb_association = NULL;

  OAILOG_FUNC_IN(LOG_S1AP);

  if (sctp_new_peer_p == NULL) {
    OAILOG_ERROR(LOG_S1AP, "sctp_new_peer_p is NULL\n");
    return RETURNerror;
  }

  /*
   * Checking that the assoc id has a valid eNB attached to.
   */
  enb_association = s1ap_state_get_enb(state, sctp_new_peer_p->assoc_id);
  if (enb_association == NULL) {
    OAILOG_DEBUG(
        LOG_S1AP, "Create eNB context for assoc_id: %d\n",
        sctp_new_peer_p->assoc_id);
    /*
     * Create new context
     */
    enb_association = s1ap_new_enb();

    if (enb_association == NULL) {
      /*
       * We failed to allocate memory
       */
      /*
       * TODO: send reject there
       */
      OAILOG_ERROR(
          LOG_S1AP, "Failed to allocate eNB context for assoc_id: %d\n",
          sctp_new_peer_p->assoc_id);
      OAILOG_FUNC_RETURN(LOG_S1AP, RETURNok);
    }
    enb_association->sctp_assoc_id = sctp_new_peer_p->assoc_id;
    enb_association->enb_id = 0xFFFFFFFF;  // home or macro eNB is 28 or 20bits.
    hashtable_rc_t hash_rc  = hashtable_ts_insert(
        &state->enbs, (const hash_key_t) enb_association->sctp_assoc_id,
        (void*) enb_association);
    if (HASH_TABLE_OK != hash_rc) {
      OAILOG_FUNC_RETURN(LOG_S1AP, RETURNerror);
    }
  } else if (
      (enb_association->s1_state == S1AP_SHUTDOWN) ||
      (enb_association->s1_state == S1AP_RESETING)) {
    OAILOG_WARNING(
        LOG_S1AP,
        "Received new association request on an association that is being %s, "
        "ignoring",
        s1_enb_state2str(enb_association->s1_state));
    OAILOG_FUNC_RETURN(LOG_S1AP, RETURNerror);
  } else {
    OAILOG_DEBUG(
        LOG_S1AP, "eNB context already exists for assoc_id: %d, update it\n",
        sctp_new_peer_p->assoc_id);
  }

  enb_association->sctp_assoc_id = sctp_new_peer_p->assoc_id;
  /*
   * Fill in in and out number of streams available on SCTP connection.
   */
  enb_association->instreams  = (sctp_stream_id_t) sctp_new_peer_p->instreams;
  enb_association->outstreams = (sctp_stream_id_t) sctp_new_peer_p->outstreams;
  /*
   * Fill in control plane IP address of RAN end point for this association
   */
  if (sctp_new_peer_p->ran_cp_ipaddr) {
    memcpy(
        enb_association->ran_cp_ipaddr, sctp_new_peer_p->ran_cp_ipaddr->data,
        sctp_new_peer_p->ran_cp_ipaddr->slen);
    enb_association->ran_cp_ipaddr_sz = sctp_new_peer_p->ran_cp_ipaddr->slen;
  }
  /*
   * initialize the next sctp stream to 1 as 0 is reserved for non
   * * * * ue associated signalling.
   */
  enb_association->next_sctp_stream = 1;
  enb_association->s1_state         = S1AP_INIT;
  OAILOG_FUNC_RETURN(LOG_S1AP, RETURNok);
}

//------------------------------------------------------------------------------
void s1ap_mme_release_ue_context(
    s1ap_state_t* state, ue_description_t* ue_ref_p, imsi64_t imsi64) {
  MessageDef* message_p = NULL;
  OAILOG_FUNC_IN(LOG_S1AP);

  if (ue_ref_p == NULL) {
    OAILOG_ERROR(LOG_S1AP, "ue_ref_p is NULL\n");
  }

  OAILOG_DEBUG_UE(
      LOG_S1AP, imsi64, "Releasing UE Context for UE id  %d \n",
      ue_ref_p->mme_ue_s1ap_id);
  /*
   * Remove UE context and inform MME_APP.
   */
  message_p = DEPRECATEDitti_alloc_new_message_fatal(
      TASK_S1AP, S1AP_UE_CONTEXT_RELEASE_COMPLETE);
  memset(
      (void*) &message_p->ittiMsg.s1ap_ue_context_release_complete, 0,
      sizeof(itti_s1ap_ue_context_release_complete_t));
  S1AP_UE_CONTEXT_RELEASE_COMPLETE(message_p).mme_ue_s1ap_id =
      ue_ref_p->mme_ue_s1ap_id;

  message_p->ittiMsgHeader.imsi = imsi64;
  send_msg_to_task(&s1ap_task_zmq_ctx, TASK_MME_APP, message_p);

  if (!(ue_ref_p->s1_ue_state == S1AP_UE_WAITING_CRR)) {
    OAILOG_ERROR(LOG_S1AP, "Incorrect S1AP UE state\n");
  }
  OAILOG_DEBUG_UE(
      LOG_S1AP, imsi64, "Removed S1AP UE " MME_UE_S1AP_ID_FMT "\n",
      (uint32_t) ue_ref_p->mme_ue_s1ap_id);

  s1ap_remove_ue(state, ue_ref_p);
  OAILOG_FUNC_OUT(LOG_S1AP);
}

//------------------------------------------------------------------------------
status_code_e s1ap_mme_handle_error_ind_message(
    s1ap_state_t* state, const sctp_assoc_id_t assoc_id,
    const sctp_stream_id_t stream, S1ap_S1AP_PDU_t* message) {
  OAILOG_FUNC_IN(LOG_S1AP);
  OAILOG_WARNING(LOG_S1AP, "ERROR IND RCVD on Stream id %d \n", stream);
  increment_counter("s1ap_error_ind_rcvd", 1, NO_LABELS);
  S1ap_ErrorIndication_t* container = NULL;
  S1ap_ErrorIndicationIEs_t* ie     = NULL;
  ue_description_t* ue_ref_p        = NULL;
  enb_ue_s1ap_id_t enb_ue_s1ap_id   = INVALID_ENB_UE_S1AP_ID;
  mme_ue_s1ap_id_t mme_ue_s1ap_id   = INVALID_MME_UE_S1AP_ID;
  S1ap_Cause_PR cause_type;
  long cause_value;

  container = &message->choice.initiatingMessage.value.choice.ErrorIndication;
  S1AP_FIND_PROTOCOLIE_BY_ID(
      S1ap_ErrorIndicationIEs_t, ie, container,
      S1ap_ProtocolIE_ID_id_MME_UE_S1AP_ID, true);
  if (ie) {
    mme_ue_s1ap_id = ie->value.choice.MME_UE_S1AP_ID;
  } else {
    OAILOG_FUNC_RETURN(LOG_S1AP, RETURNerror);
  }

  S1AP_FIND_PROTOCOLIE_BY_ID(
      S1ap_ErrorIndicationIEs_t, ie, container,
      S1ap_ProtocolIE_ID_id_eNB_UE_S1AP_ID, true);
  if (ie) {
    // eNB UE S1AP ID is limited to 24 bits
    enb_ue_s1ap_id = (enb_ue_s1ap_id_t)(
        ie->value.choice.ENB_UE_S1AP_ID & ENB_UE_S1AP_ID_MASK);
  } else {
    OAILOG_FUNC_RETURN(LOG_S1AP, RETURNerror);
  }

  OAILOG_INFO(
      LOG_S1AP,
      "ERROR IND RCVD with mme UE s1ap id " MME_UE_S1AP_ID_FMT
      " and enb UE s1ap id " ENB_UE_S1AP_ID_FMT "\n",
      mme_ue_s1ap_id, enb_ue_s1ap_id);
  S1AP_FIND_PROTOCOLIE_BY_ID(
      S1ap_ErrorIndicationIEs_t, ie, container, S1ap_ProtocolIE_ID_id_Cause,
      true);
  if (ie) {
    cause_type = ie->value.choice.Cause.present;
  } else {
    OAILOG_FUNC_RETURN(LOG_S1AP, RETURNerror);
  }

  if ((ue_ref_p = s1ap_state_get_ue_mmeid((uint32_t) mme_ue_s1ap_id)) == NULL) {
    OAILOG_WARNING(
        LOG_S1AP,
        "No UE is attached to this mme UE s1ap id: " MME_UE_S1AP_ID_FMT
        " and eNB UE s1ap id: \n" ENB_UE_S1AP_ID_FMT,
        mme_ue_s1ap_id, enb_ue_s1ap_id);
    OAILOG_FUNC_RETURN(LOG_S1AP, RETURNerror);
  }

  imsi64_t imsi64           = INVALID_IMSI64;
  s1ap_imsi_map_t* imsi_map = get_s1ap_imsi_map();
  hashtable_uint64_ts_get(
      imsi_map->mme_ue_id_imsi_htbl, (const hash_key_t) mme_ue_s1ap_id,
      &imsi64);

  switch (cause_type) {
    case S1ap_Cause_PR_radioNetwork:
      cause_value = ie->value.choice.Cause.choice.radioNetwork;
      OAILOG_DEBUG_UE(
          LOG_S1AP, imsi64,
          "Error Indication with Cause_Type = Radio Network "
          "and Cause_Value = %ld\n",
          cause_value);
      s1ap_send_mme_ue_context_release(
          state, ue_ref_p, S1AP_RADIO_EUTRAN_GENERATED_REASON,
          ie->value.choice.Cause, imsi64);
      break;

    case S1ap_Cause_PR_transport:
      cause_value = ie->value.choice.Cause.choice.transport;
      OAILOG_DEBUG_UE(
          LOG_S1AP, imsi64,
          "Error Indication with Cause_Type = Transport and "
          "Cause_Value = %ld\n",
          cause_value);
      break;

    case S1ap_Cause_PR_nas:
      cause_value = ie->value.choice.Cause.choice.nas;
      OAILOG_DEBUG_UE(
          LOG_S1AP, imsi64,
          "Error Indication with Cause_Type = NAS and "
          "Cause_Value = %ld\n",
          cause_value);
      break;

    case S1ap_Cause_PR_protocol:
      cause_value = ie->value.choice.Cause.choice.protocol;
      OAILOG_DEBUG_UE(
          LOG_S1AP, imsi64,
          "Error Indication with Cause_Type = Protocol and "
          "Cause_Value = %ld\n",
          cause_value);
      break;

    case S1ap_Cause_PR_misc:
      cause_value = ie->value.choice.Cause.choice.misc;
      OAILOG_DEBUG_UE(
          LOG_S1AP, imsi64,
          "Error Indication with Cause_Type = MISC and "
          "Cause_Value = %ld\n",
          cause_value);
      break;

    default:
      OAILOG_ERROR_UE(
          LOG_S1AP, imsi64, "Error Indication with Invalid Cause_Type = %d\n",
          cause_type);
      OAILOG_FUNC_RETURN(LOG_S1AP, RETURNerror);
  }

  OAILOG_FUNC_RETURN(LOG_S1AP, RETURNok);
}

//------------------------------------------------------------------------------
status_code_e s1ap_mme_handle_erab_setup_response(
    s1ap_state_t* state, const sctp_assoc_id_t assoc_id,
    const sctp_stream_id_t stream, S1ap_S1AP_PDU_t* pdu) {
  OAILOG_FUNC_IN(LOG_S1AP);
  S1ap_E_RABSetupResponse_t* container = NULL;
  S1ap_E_RABSetupResponseIEs_t* ie     = NULL;
  ue_description_t* ue_ref_p           = NULL;
  MessageDef* message_p                = NULL;
  enb_ue_s1ap_id_t enb_ue_s1ap_id      = INVALID_ENB_UE_S1AP_ID;
  mme_ue_s1ap_id_t mme_ue_s1ap_id      = INVALID_MME_UE_S1AP_ID;
  int rc                               = RETURNok;
  imsi64_t imsi64                      = INVALID_IMSI64;

  container = &pdu->choice.successfulOutcome.value.choice.E_RABSetupResponse;

  S1AP_FIND_PROTOCOLIE_BY_ID(
      S1ap_E_RABSetupResponseIEs_t, ie, container,
      S1ap_ProtocolIE_ID_id_MME_UE_S1AP_ID, true);
  if (ie) {
    mme_ue_s1ap_id = ie->value.choice.MME_UE_S1AP_ID;
  } else {
    OAILOG_FUNC_RETURN(LOG_S1AP, RETURNerror);
  }
  S1AP_FIND_PROTOCOLIE_BY_ID(
      S1ap_E_RABSetupResponseIEs_t, ie, container,
      S1ap_ProtocolIE_ID_id_eNB_UE_S1AP_ID, true);
  if (ie) {
    // eNB UE S1AP ID is limited to 24 bits
    enb_ue_s1ap_id = (enb_ue_s1ap_id_t)(
        ie->value.choice.ENB_UE_S1AP_ID & ENB_UE_S1AP_ID_MASK);
  } else {
    OAILOG_FUNC_RETURN(LOG_S1AP, RETURNerror);
  }
  if ((ue_ref_p = s1ap_state_get_ue_mmeid((uint32_t) mme_ue_s1ap_id)) == NULL) {
    OAILOG_DEBUG(
        LOG_S1AP,
        "No UE is attached to this mme UE s1ap id: " MME_UE_S1AP_ID_FMT "\n",
        mme_ue_s1ap_id);
    OAILOG_FUNC_RETURN(LOG_S1AP, RETURNerror);
  }

  if (ue_ref_p->enb_ue_s1ap_id != enb_ue_s1ap_id) {
    OAILOG_DEBUG(
        LOG_S1AP,
        "Mismatch in eNB UE S1AP ID, known: " ENB_UE_S1AP_ID_FMT
        ", received: " ENB_UE_S1AP_ID_FMT "\n",
        ue_ref_p->enb_ue_s1ap_id, enb_ue_s1ap_id);
    OAILOG_FUNC_RETURN(LOG_S1AP, RETURNerror);
  }

  s1ap_imsi_map_t* imsi_map = get_s1ap_imsi_map();
  hashtable_uint64_ts_get(
      imsi_map->mme_ue_id_imsi_htbl,
      (const hash_key_t) ue_ref_p->mme_ue_s1ap_id, &imsi64);

  message_p =
      DEPRECATEDitti_alloc_new_message_fatal(TASK_S1AP, S1AP_E_RAB_SETUP_RSP);
  S1AP_E_RAB_SETUP_RSP(message_p).mme_ue_s1ap_id = ue_ref_p->mme_ue_s1ap_id;
  S1AP_E_RAB_SETUP_RSP(message_p).enb_ue_s1ap_id = ue_ref_p->enb_ue_s1ap_id;
  S1AP_E_RAB_SETUP_RSP(message_p).e_rab_setup_list.no_of_items           = 0;
  S1AP_E_RAB_SETUP_RSP(message_p).e_rab_failed_to_setup_list.no_of_items = 0;

  S1AP_FIND_PROTOCOLIE_BY_ID(
      S1ap_E_RABSetupResponseIEs_t, ie, container,
      S1ap_ProtocolIE_ID_id_E_RABSetupListBearerSURes, false);

  if (ie) {
    int num_erab = ie->value.choice.E_RABSetupListBearerSURes.list.count;
    for (int index = 0; index < num_erab; index++) {
      S1ap_E_RABSetupItemBearerSUResIEs_t* erab_setup_item =
          (S1ap_E_RABSetupItemBearerSUResIEs_t*)
              ie->value.choice.E_RABSetupListBearerSURes.list.array[index];
      S1ap_E_RABSetupItemBearerSURes_t* e_rab_setup_item_bearer_su_res =
          &erab_setup_item->value.choice.E_RABSetupItemBearerSURes;
      S1AP_E_RAB_SETUP_RSP(message_p).e_rab_setup_list.item[index].e_rab_id =
          e_rab_setup_item_bearer_su_res->e_RAB_ID;
      S1AP_E_RAB_SETUP_RSP(message_p)
          .e_rab_setup_list.item[index]
          .transport_layer_address = blk2bstr(
          e_rab_setup_item_bearer_su_res->transportLayerAddress.buf,
          e_rab_setup_item_bearer_su_res->transportLayerAddress.size);
      S1AP_E_RAB_SETUP_RSP(message_p).e_rab_setup_list.item[index].gtp_teid =
          htonl(*((uint32_t*) e_rab_setup_item_bearer_su_res->gTP_TEID.buf));
      S1AP_E_RAB_SETUP_RSP(message_p).e_rab_setup_list.no_of_items += 1;
    }
  }

  S1AP_FIND_PROTOCOLIE_BY_ID(
      S1ap_E_RABSetupResponseIEs_t, ie, container,
      S1ap_ProtocolIE_ID_id_E_RABFailedToSetupListBearerSURes, false);
  if (ie) {
    const S1ap_E_RABList_t* const e_rab_list = &ie->value.choice.E_RABList;
    int num_erab = ie->value.choice.E_RABList.list.count;
    for (int index = 0; index < num_erab; index++) {
      const S1ap_E_RABItemIEs_t* const erab_item_ies =
          (S1ap_E_RABItemIEs_t*) e_rab_list->list.array[index];
      const S1ap_E_RABItem_t* const erab_item =
          (S1ap_E_RABItem_t*) &erab_item_ies->value.choice.E_RABItem;
      S1AP_E_RAB_SETUP_RSP(message_p)
          .e_rab_failed_to_setup_list.item[index]
          .e_rab_id = erab_item->e_RAB_ID;

      S1AP_E_RAB_SETUP_RSP(message_p)
          .e_rab_failed_to_setup_list.item[index]
          .cause = erab_item->cause;
      S1AP_E_RAB_SETUP_RSP(message_p).e_rab_failed_to_setup_list.no_of_items +=
          1;
    }
  }

  message_p->ittiMsgHeader.imsi = imsi64;
  rc = send_msg_to_task(&s1ap_task_zmq_ctx, TASK_MME_APP, message_p);
  OAILOG_FUNC_RETURN(LOG_S1AP, rc);
}

//------------------------------------------------------------------------------
status_code_e s1ap_mme_handle_erab_setup_failure(
    s1ap_state_t* state, const sctp_assoc_id_t assoc_id,
    const sctp_stream_id_t stream, S1ap_S1AP_PDU_t* message) {
  Fatal("TODO Implement s1ap_mme_handle_erab_setup_failure");
}

//------------------------------------------------------------------------------
status_code_e s1ap_mme_handle_enb_reset(
    s1ap_state_t* state, const sctp_assoc_id_t assoc_id,
    const sctp_stream_id_t stream, S1ap_S1AP_PDU_t* pdu) {
  MessageDef* msg                                = NULL;
  itti_s1ap_enb_initiated_reset_req_t* reset_req = NULL;
  ue_description_t* ue_ref_p                     = NULL;
  enb_description_t* enb_association             = NULL;
  s1ap_reset_type_t s1ap_reset_type;
  S1ap_Reset_t* container                                        = NULL;
  S1ap_ResetIEs_t* ie                                            = NULL;
  S1ap_UE_associatedLogicalS1_ConnectionItem_t* s1_sig_conn_id_p = NULL;
  mme_ue_s1ap_id_t mme_ue_s1ap_id;
  enb_ue_s1ap_id_t enb_ue_s1ap_id;
  imsi64_t imsi64                        = INVALID_IMSI64;
  arg_s1ap_construct_enb_reset_req_t arg = {0};
  uint32_t i                             = 0;
  int rc                                 = RETURNok;

  OAILOG_FUNC_IN(LOG_S1AP);

  enb_association = s1ap_state_get_enb(state, assoc_id);

  if (enb_association == NULL) {
    OAILOG_ERROR(LOG_S1AP, "No eNB attached to this assoc_id: %d\n", assoc_id);
    OAILOG_FUNC_RETURN(LOG_S1AP, RETURNerror);
  }

  if (enb_association->s1_state != S1AP_READY) {
    // ignore the message if s1 not ready
    OAILOG_INFO(
        LOG_S1AP,
        "S1 setup is not done.Invalid state.Ignoring ENB Initiated Reset.eNB "
        "Id "
        "= %d , S1AP state = %d \n",
        enb_association->enb_id, enb_association->s1_state);
    OAILOG_FUNC_RETURN(LOG_S1AP, RETURNok);
  }

  if (enb_association->nb_ue_associated == 0) {
    // Even if there are no UEs connected, we proceed -- this can happen if we
    // receive a reset during a handover procedure, for example.
    OAILOG_INFO(
        LOG_S1AP,
        "No UEs connected, still proceeding with ENB Initiated Reset. eNB Id = "
        "%d\n",
        enb_association->enb_id);
  }

  // Check the reset type - partial_reset OR reset_all
  container = &pdu->choice.initiatingMessage.value.choice.Reset;

  S1AP_FIND_PROTOCOLIE_BY_ID(
      S1ap_ResetIEs_t, ie, container, S1ap_ProtocolIE_ID_id_ResetType, true);

  S1ap_ResetType_t* resetType = &ie->value.choice.ResetType;

  switch (resetType->present) {
    case S1ap_ResetType_PR_s1_Interface:
      s1ap_reset_type = RESET_ALL;
      break;
    case S1ap_ResetType_PR_partOfS1_Interface:
      s1ap_reset_type = RESET_PARTIAL;
      break;
    default:
      OAILOG_ERROR(
          LOG_S1AP, "Reset Request from eNB  with Invalid reset_type = %d\n",
          resetType->present);
      // TBD - Here MME should send Error Indication as it is abnormal scenario.
      OAILOG_FUNC_RETURN(LOG_S1AP, RETURNerror);
  }

  if (s1ap_reset_type == RESET_PARTIAL) {
    int reset_count = resetType->choice.partOfS1_Interface.list.count;
    if (reset_count == 0) {
      OAILOG_ERROR(
          LOG_S1AP,
          "Partial Reset Request without any S1 signaling connection. Ignoring "
          "it \n");
      // TBD - Here MME should send Error Indication as it is abnormal scenario.
      OAILOG_FUNC_RETURN(LOG_S1AP, RETURNerror);
    }
    if (reset_count > enb_association->nb_ue_associated) {
      // We proceed here since we could encounter this situation when we
      // receive a reset from the target eNB during a handover procedure.
      OAILOG_WARNING(
          LOG_S1AP,
          "Partial Reset Request. Requested number of UEs %d to be reset is "
          "more "
          "than connected UEs %d \n",
          reset_count, enb_association->nb_ue_associated);
    }
  }
  msg = DEPRECATEDitti_alloc_new_message_fatal(
      TASK_S1AP, S1AP_ENB_INITIATED_RESET_REQ);
  reset_req = &S1AP_ENB_INITIATED_RESET_REQ(msg);

  reset_req->s1ap_reset_type = s1ap_reset_type;
  reset_req->enb_id          = enb_association->enb_id;
  reset_req->sctp_assoc_id   = assoc_id;
  reset_req->sctp_stream_id  = stream;

  switch (s1ap_reset_type) {
    case RESET_ALL:
      increment_counter("s1_reset_from_enb", 1, 1, "type", "reset_all");

      reset_req->num_ue = enb_association->nb_ue_associated;

      reset_req->ue_to_reset_list = calloc(
          enb_association->nb_ue_associated,
          sizeof(*reset_req->ue_to_reset_list));

      if (reset_req->ue_to_reset_list == NULL) {
        OAILOG_ERROR(LOG_S1AP, "ue_to_reset_list is NULL\n");
        return RETURNerror;
      }
      arg.msg              = msg;
      arg.current_ue_index = 0;
      hashtable_uint64_ts_apply_callback_on_elements(
          &enb_association->ue_id_coll, construct_s1ap_mme_full_reset_req, &arg,
          NULL);
      // EURECOM LG 2020-07-16 added break here
      break;
    case RESET_PARTIAL:
      // Partial Reset
      increment_counter("s1_reset_from_enb", 1, 1, "type", "reset_partial");
      reset_req->num_ue = resetType->choice.partOfS1_Interface.list.count;
      reset_req->ue_to_reset_list = calloc(
          resetType->choice.partOfS1_Interface.list.count,
          sizeof(*(reset_req->ue_to_reset_list)));
      // Careful! This struct allocated will be re-used in another itti message.
      if (reset_req->ue_to_reset_list == NULL) {
        OAILOG_ERROR(LOG_S1AP, "ue_to_reset_list is NULL\n");
        return RETURNerror;
      }
      for (i = 0; i < resetType->choice.partOfS1_Interface.list.count; i++) {
        s1_sig_conn_id_p =
            (S1ap_UE_associatedLogicalS1_ConnectionItem_t*)
                resetType->choice.partOfS1_Interface.list.array[i];
        if (s1_sig_conn_id_p == NULL) {
          OAILOG_ERROR(LOG_S1AP, "s1_sig_conn_id_p is NULL\n");
          return RETURNerror;
        }
        S1ap_UE_associatedLogicalS1_ConnectionItemResAck_t* s1_sig_conn_p =
            (S1ap_UE_associatedLogicalS1_ConnectionItemResAck_t*) ie->value
                .choice.ResetType.choice.partOfS1_Interface.list.array[i];
        if (!s1_sig_conn_p) {
          OAILOG_ERROR(
              LOG_S1AP,
              "No logical S1 connection item could be found for the "
              "partial connection. "
              "Ignoring the received partial reset request. \n");
          OAILOG_FUNC_RETURN(LOG_S1AP, RETURNerror);
        }
        S1ap_UE_associatedLogicalS1_ConnectionItem_t* s1_sig_conn_id_p =
            &s1_sig_conn_p->value.choice.UE_associatedLogicalS1_ConnectionItem;

        if (s1_sig_conn_id_p->mME_UE_S1AP_ID != NULL) {
          mme_ue_s1ap_id =
              (mme_ue_s1ap_id_t) * (s1_sig_conn_id_p->mME_UE_S1AP_ID);
          s1ap_imsi_map_t* imsi_map = get_s1ap_imsi_map();
          hashtable_uint64_ts_get(
              imsi_map->mme_ue_id_imsi_htbl, (const hash_key_t) mme_ue_s1ap_id,
              &imsi64);
          if ((ue_ref_p = s1ap_state_get_ue_mmeid(mme_ue_s1ap_id)) != NULL) {
            if (s1_sig_conn_id_p->eNB_UE_S1AP_ID != NULL) {
              enb_ue_s1ap_id_t enb_ue_s1ap_id =
                  (enb_ue_s1ap_id_t) * (s1_sig_conn_id_p->eNB_UE_S1AP_ID);
              if (ue_ref_p->enb_ue_s1ap_id ==
                  (enb_ue_s1ap_id & ENB_UE_S1AP_ID_MASK)) {
                reset_req->ue_to_reset_list[i].mme_ue_s1ap_id =
                    ue_ref_p->mme_ue_s1ap_id;
                enb_ue_s1ap_id &= ENB_UE_S1AP_ID_MASK;
                reset_req->ue_to_reset_list[i].enb_ue_s1ap_id = enb_ue_s1ap_id;
              } else {
                // mismatch in enb_ue_s1ap_id sent by eNB and stored in S1AP ue
                // context in EPC. Abnormal case.
                reset_req->ue_to_reset_list[i].mme_ue_s1ap_id =
                    ue_ref_p->mme_ue_s1ap_id;
                reset_req->ue_to_reset_list[i].enb_ue_s1ap_id =
                    (enb_ue_s1ap_id_t) * (s1_sig_conn_id_p->eNB_UE_S1AP_ID);
                OAILOG_ERROR_UE(
                    LOG_S1AP, imsi64,
                    "Partial Reset Request:enb_ue_s1ap_id mismatch between id "
                    "%d "
                    "sent by eNB and id %d stored in epc for mme_ue_s1ap_id %d "
                    "\n",
                    enb_ue_s1ap_id, ue_ref_p->enb_ue_s1ap_id, mme_ue_s1ap_id);
              }
            } else {
              reset_req->ue_to_reset_list[i].mme_ue_s1ap_id =
                  ue_ref_p->mme_ue_s1ap_id;
              reset_req->ue_to_reset_list[i].enb_ue_s1ap_id =
                  INVALID_ENB_UE_S1AP_ID;
            }
          } else {
            OAILOG_ERROR_UE(
                LOG_S1AP, imsi64,
                "Partial Reset Request - No UE context found for "
                "mme_ue_s1ap_id "
                "%d "
                "\n",
                mme_ue_s1ap_id);
            reset_req->ue_to_reset_list[i].mme_ue_s1ap_id =
                (mme_ue_s1ap_id_t) * (s1_sig_conn_id_p->mME_UE_S1AP_ID);
            if (s1_sig_conn_id_p->eNB_UE_S1AP_ID != NULL) {
              reset_req->ue_to_reset_list[i].enb_ue_s1ap_id =
                  (enb_ue_s1ap_id_t) * (s1_sig_conn_id_p->eNB_UE_S1AP_ID);
            } else {
              reset_req->ue_to_reset_list[i].enb_ue_s1ap_id =
                  INVALID_ENB_UE_S1AP_ID;
            }
          }
          free_wrapper((void**) &s1_sig_conn_id_p->mME_UE_S1AP_ID);
          if (s1_sig_conn_id_p->eNB_UE_S1AP_ID != NULL) {
            free_wrapper((void**) &s1_sig_conn_id_p->eNB_UE_S1AP_ID);
          }
        } else {
          if (s1_sig_conn_id_p->eNB_UE_S1AP_ID != NULL) {
            enb_ue_s1ap_id =
                (enb_ue_s1ap_id_t) * (s1_sig_conn_id_p->eNB_UE_S1AP_ID);
            if ((ue_ref_p = s1ap_state_get_ue_enbid(
                     enb_association->sctp_assoc_id, enb_ue_s1ap_id)) != NULL) {
              enb_ue_s1ap_id &= ENB_UE_S1AP_ID_MASK;
              reset_req->ue_to_reset_list[i].enb_ue_s1ap_id = enb_ue_s1ap_id;
            } else {
              OAILOG_ERROR_UE(
                  LOG_S1AP, imsi64,
                  "Partial Reset Request without any valid S1 signaling "
                  "connection.Sending Reset Ack with received signalling "
                  "connection IDs \n");
              reset_req->ue_to_reset_list[i].enb_ue_s1ap_id =
                  (enb_ue_s1ap_id_t) * (s1_sig_conn_id_p->eNB_UE_S1AP_ID);
            }
            reset_req->ue_to_reset_list[i].mme_ue_s1ap_id =
                INVALID_MME_UE_S1AP_ID;
            free_wrapper((void**) &s1_sig_conn_id_p->eNB_UE_S1AP_ID);
          } else {
            OAILOG_ERROR_UE(
                LOG_S1AP, imsi64,
                "Partial Reset Request without any valid S1 signaling "
                "connection.Sending Reset Ack with received signalling "
                "connection IDs \n");
            reset_req->ue_to_reset_list[i].mme_ue_s1ap_id =
                INVALID_MME_UE_S1AP_ID;
            reset_req->ue_to_reset_list[i].enb_ue_s1ap_id =
                INVALID_ENB_UE_S1AP_ID;
          }
        }
      }
  }

  msg->ittiMsgHeader.imsi = imsi64;
  rc = send_msg_to_task(&s1ap_task_zmq_ctx, TASK_MME_APP, msg);
  OAILOG_FUNC_RETURN(LOG_S1AP, rc);
}
//------------------------------------------------------------------------------
status_code_e s1ap_handle_enb_initiated_reset_ack(
    const itti_s1ap_enb_initiated_reset_ack_t* const enb_reset_ack_p,
    imsi64_t imsi64) {
  uint8_t* buffer = NULL;
  uint32_t length = 0;
  S1ap_S1AP_PDU_t pdu;
  /** Reset Acknowledgment. */
  S1ap_ResetAcknowledge_t* out;
  S1ap_ResetAcknowledgeIEs_t* ie = NULL;
  int rc                         = RETURNok;

  OAILOG_FUNC_IN(LOG_S1AP);

  memset(&pdu, 0, sizeof(pdu));
  pdu.present = S1ap_S1AP_PDU_PR_successfulOutcome;
  pdu.choice.successfulOutcome.procedureCode = S1ap_ProcedureCode_id_Reset;
  pdu.choice.successfulOutcome.criticality   = S1ap_Criticality_ignore;
  pdu.choice.successfulOutcome.value.present =
      S1ap_SuccessfulOutcome__value_PR_ResetAcknowledge;
  out = &pdu.choice.successfulOutcome.value.choice.ResetAcknowledge;

  if (enb_reset_ack_p->s1ap_reset_type == RESET_PARTIAL) {
    if (!(enb_reset_ack_p->num_ue > 0)) {
      OAILOG_ERROR(LOG_S1AP, "Incorrect number of UEs\n");
      return RETURNerror;
    }
    ie = (S1ap_ResetAcknowledgeIEs_t*) calloc(
        1, sizeof(S1ap_ResetAcknowledgeIEs_t));
    ie->id = S1ap_ProtocolIE_ID_id_UE_associatedLogicalS1_ConnectionListResAck;
    ie->criticality = S1ap_Criticality_ignore;
    ie->value.present =
        S1ap_ResetAcknowledgeIEs__value_PR_UE_associatedLogicalS1_ConnectionListResAck;
    ASN_SEQUENCE_ADD(&out->protocolIEs.list, ie);
    /** MME UE S1AP ID. */
    S1ap_UE_associatedLogicalS1_ConnectionListResAck_t* ie_p =
        &ie->value.choice.UE_associatedLogicalS1_ConnectionListResAck;
    for (uint32_t i = 0; i < enb_reset_ack_p->num_ue; i++) {
      S1ap_UE_associatedLogicalS1_ConnectionItemResAck_t* sig_conn_item =
          calloc(1, sizeof(S1ap_UE_associatedLogicalS1_ConnectionItemResAck_t));
      sig_conn_item->id =
          S1ap_ProtocolIE_ID_id_UE_associatedLogicalS1_ConnectionItem;
      sig_conn_item->criticality = S1ap_Criticality_ignore;
      sig_conn_item->value.present =
          S1ap_UE_associatedLogicalS1_ConnectionItemResAck__value_PR_UE_associatedLogicalS1_ConnectionItem;
      S1ap_UE_associatedLogicalS1_ConnectionItem_t* item =
          &sig_conn_item->value.choice.UE_associatedLogicalS1_ConnectionItem;
      if (enb_reset_ack_p->ue_to_reset_list[i].mme_ue_s1ap_id !=
          INVALID_MME_UE_S1AP_ID) {
        item->mME_UE_S1AP_ID = calloc(1, sizeof(S1ap_MME_UE_S1AP_ID_t));
        *item->mME_UE_S1AP_ID =
            enb_reset_ack_p->ue_to_reset_list[i].mme_ue_s1ap_id;
      } else {
        item->mME_UE_S1AP_ID = NULL;
      }
      if (enb_reset_ack_p->ue_to_reset_list[i].enb_ue_s1ap_id !=
          INVALID_ENB_UE_S1AP_ID) {
        item->eNB_UE_S1AP_ID = calloc(1, sizeof(S1ap_ENB_UE_S1AP_ID_t));
        *item->eNB_UE_S1AP_ID =
            enb_reset_ack_p->ue_to_reset_list[i].enb_ue_s1ap_id;
      } else {
        item->eNB_UE_S1AP_ID = NULL;
      }
      ASN_SEQUENCE_ADD(&ie_p->list, sig_conn_item);
    }
  }
  if (s1ap_mme_encode_pdu(&pdu, &buffer, &length) < 0) {
    OAILOG_ERROR(LOG_S1AP, "Failed to S1 Reset command \n");
    /** We rely on the handover_notify timeout to remove the UE context. */
    DevAssert(!buffer);
    OAILOG_FUNC_RETURN(LOG_S1AP, RETURNerror);
  }
  increment_counter("s1_reset_from_enb", 1, 1, "action", "reset_ack_sent");
  if (buffer) {
    bstring b = blk2bstr(buffer, length);
    free_wrapper((void**) &buffer);
    rc = s1ap_mme_itti_send_sctp_request(
        &b, enb_reset_ack_p->sctp_assoc_id, enb_reset_ack_p->sctp_stream_id,
        INVALID_MME_UE_S1AP_ID);
  }
  OAILOG_FUNC_RETURN(LOG_S1AP, rc);
}

//-------------------------------------------------------------------------------
status_code_e s1ap_handle_paging_request(
    s1ap_state_t* state, const itti_s1ap_paging_request_t* paging_request,
    imsi64_t imsi64) {
  OAILOG_FUNC_IN(LOG_S1AP);

  if (paging_request == NULL) {
    OAILOG_ERROR(LOG_S1AP, "paging_request is NULL\n");
    return RETURNerror;
  }
  int rc                  = RETURNok;
  uint8_t num_of_tac      = 0;
  uint16_t tai_list_count = paging_request->tai_list_count;

  bool is_tai_found    = false;
  uint32_t idx         = 0;
  uint8_t* buffer_p    = NULL;
  uint32_t length      = 0;
  S1ap_S1AP_PDU_t pdu  = {0};
  S1ap_Paging_t* out   = NULL;
  S1ap_PagingIEs_t* ie = NULL;

  memset(&pdu, 0, sizeof(pdu));
  pdu.present = S1ap_S1AP_PDU_PR_initiatingMessage;
  pdu.choice.initiatingMessage.procedureCode = S1ap_ProcedureCode_id_Paging;
  pdu.choice.initiatingMessage.criticality   = S1ap_Criticality_ignore;
  pdu.choice.initiatingMessage.value.present =
      S1ap_InitiatingMessage__value_PR_Paging;
  out = &pdu.choice.initiatingMessage.value.choice.Paging;

  // Encode and set the UE Identity Index Value.
  ie                = (S1ap_PagingIEs_t*) calloc(1, sizeof(S1ap_PagingIEs_t));
  ie->id            = S1ap_ProtocolIE_ID_id_UEIdentityIndexValue;
  ie->criticality   = S1ap_Criticality_ignore;
  ie->value.present = S1ap_PagingIEs__value_PR_UEIdentityIndexValue;
  UE_ID_INDEX_TO_BIT_STRING(
      (uint16_t)(imsi64 % 1024), &ie->value.choice.UEIdentityIndexValue);
  ASN_SEQUENCE_ADD(&out->protocolIEs.list, ie);

  // Set UE Paging Identity
  ie                = (S1ap_PagingIEs_t*) calloc(1, sizeof(S1ap_PagingIEs_t));
  ie->id            = S1ap_ProtocolIE_ID_id_UEPagingID;
  ie->criticality   = S1ap_Criticality_ignore;
  ie->value.present = S1ap_PagingIEs__value_PR_UEPagingID;
  if (paging_request->paging_id == S1AP_PAGING_ID_STMSI) {
    ie->value.choice.UEPagingID.present = S1ap_UEPagingID_PR_s_TMSI;
    M_TMSI_TO_OCTET_STRING(
        paging_request->m_tmsi,
        &ie->value.choice.UEPagingID.choice.s_TMSI.m_TMSI);
    // todo: chose the right gummei or get it from the request!
    MME_CODE_TO_OCTET_STRING(
        paging_request->mme_code,
        &ie->value.choice.UEPagingID.choice.s_TMSI.mMEC);
  } else if (paging_request->paging_id == S1AP_PAGING_ID_IMSI) {
    ie->value.choice.UEPagingID.present = S1ap_UEPagingID_PR_iMSI;
    IMSI_TO_OCTET_STRING(
        paging_request->imsi, paging_request->imsi_length,
        &ie->value.choice.UEPagingID.choice.iMSI);
  }
  ASN_SEQUENCE_ADD(&out->protocolIEs.list, ie);

  // Encode the CN Domain.
  ie                = (S1ap_PagingIEs_t*) calloc(1, sizeof(S1ap_PagingIEs_t));
  ie->id            = S1ap_ProtocolIE_ID_id_CNDomain;
  ie->criticality   = S1ap_Criticality_ignore;
  ie->value.present = S1ap_PagingIEs__value_PR_CNDomain;
  if (paging_request->domain_indicator == CN_DOMAIN_PS) {
    ie->value.choice.CNDomain = S1ap_CNDomain_ps;
  } else if (paging_request->domain_indicator == CN_DOMAIN_CS) {
    ie->value.choice.CNDomain = S1ap_CNDomain_cs;
  }
  ASN_SEQUENCE_ADD(&out->protocolIEs.list, ie);

  // Set TAI list
  ie                = (S1ap_PagingIEs_t*) calloc(1, sizeof(S1ap_PagingIEs_t));
  ie->id            = S1ap_ProtocolIE_ID_id_TAIList;
  ie->criticality   = S1ap_Criticality_ignore;
  ie->value.present = S1ap_PagingIEs__value_PR_TAIList;
  ASN_SEQUENCE_ADD(&out->protocolIEs.list, ie);
  S1ap_TAIList_t* const tai_list = &ie->value.choice.TAIList;

  mme_config_read_lock(&mme_config);
  for (int tai_idx = 0; tai_idx < tai_list_count; tai_idx++) {
    num_of_tac = paging_request->paging_tai_list[tai_idx].numoftac;
    // Total number of TACs = number of tac + current ENB's tac(1)
    for (int idx = 0; idx < (num_of_tac + 1); idx++) {
      S1ap_TAIItemIEs_t* tai_item_ies = calloc(1, sizeof(S1ap_TAIItemIEs_t));
      if (tai_item_ies == NULL) {
        OAILOG_ERROR_UE(LOG_S1AP, imsi64, "Failed to allocate memory\n");
        OAILOG_FUNC_RETURN(LOG_S1AP, RETURNerror);
      }

      tai_item_ies->id            = S1ap_ProtocolIE_ID_id_TAIItem;
      tai_item_ies->criticality   = S1ap_Criticality_ignore;
      tai_item_ies->value.present = S1ap_TAIItemIEs__value_PR_TAIItem;
      S1ap_TAIItem_t* tai_item    = &tai_item_ies->value.choice.TAIItem;

      PLMN_T_TO_PLMNID(
          paging_request->paging_tai_list[tai_idx].tai_list[idx].plmn,
          &tai_item->tAI.pLMNidentity);
      TAC_TO_ASN1(
          paging_request->paging_tai_list[tai_idx].tai_list[idx].tac,
          &tai_item->tAI.tAC);
      ASN_SEQUENCE_ADD(&tai_list->list, tai_item_ies);
    }
  }

  mme_config_unlock(&mme_config);

  // Encoding without allocating, buffer_p is allocated by asn.1c
  int err = 0;
  if (s1ap_mme_encode_pdu(&pdu, &buffer_p, &length) < 0) {
    err = 1;
  }
  // TODO look why called proc s1ap_mme_encode_pdu do not return value < 0
  if (length <= 0) {
    err = 1;
  }
  if (err) {
    OAILOG_ERROR_UE(
        LOG_S1AP, imsi64, "Failed to encode paging message for IMSI %s\n",
        paging_request->imsi);
    OAILOG_FUNC_RETURN(LOG_S1AP, RETURNerror);
  }

  /*Fetching eNB list to send paging request message*/
  hashtable_element_array_t* enb_array = NULL;
  enb_description_t* enb_ref_p         = NULL;
  if (state == NULL) {
    OAILOG_ERROR(LOG_S1AP, "eNB Information is NULL!\n");
    OAILOG_FUNC_RETURN(LOG_S1AP, RETURNerror);
  }
  enb_array = hashtable_ts_get_elements(&state->enbs);
  if (enb_array == NULL) {
    OAILOG_ERROR(LOG_S1AP, "Could not find eNB hashlist!\n");
    OAILOG_FUNC_RETURN(LOG_S1AP, RETURNerror);
  }
  const paging_tai_list_t* p_tai_list = paging_request->paging_tai_list;
  for (idx = 0; idx < enb_array->num_elements; idx++) {
    enb_ref_p = (enb_description_t*) enb_array->elements[idx];
    if (enb_ref_p->s1_state == S1AP_READY) {
      supported_ta_list_t* enb_ta_list = &enb_ref_p->supported_ta_list;

      if ((is_tai_found = s1ap_paging_compare_ta_lists(
               enb_ta_list, p_tai_list, paging_request->tai_list_count))) {
        bstring paging_msg_buffer = blk2bstr(buffer_p, length);
        rc                        = s1ap_mme_itti_send_sctp_request(
            &paging_msg_buffer, enb_ref_p->sctp_assoc_id,
            0,   // Stream id 0 for non UE related
                 // S1AP message
            0);  // mme_ue_s1ap_id 0 because UE in idle
      }
    }
  }
  free_wrapper((void**) &enb_array->elements);
  free_wrapper((void**) &enb_array);
  free(buffer_p);
  if (rc != RETURNok) {
    OAILOG_ERROR(
        LOG_S1AP, "Failed to send paging message over sctp for IMSI %s\n",
        paging_request->imsi);
  } else {
    OAILOG_INFO(
        LOG_S1AP, "Sent paging message over sctp for IMSI %s\n",
        paging_request->imsi);
  }

  OAILOG_FUNC_RETURN(LOG_S1AP, rc);
}

//------------------------------------------------------------------------------
status_code_e s1ap_mme_handle_erab_modification_indication(
    s1ap_state_t* state, const sctp_assoc_id_t assoc_id,
    const sctp_stream_id_t stream, S1ap_S1AP_PDU_t* pdu) {
  OAILOG_FUNC_IN(LOG_S1AP);
  enb_ue_s1ap_id_t enb_ue_s1ap_id               = 0;
  mme_ue_s1ap_id_t mme_ue_s1ap_id               = 0;
  int rc                                        = RETURNok;
  S1ap_E_RABModificationIndication_t* container = NULL;
  S1ap_E_RABModificationIndicationIEs_t* ie     = NULL;
  ue_description_t* ue_ref_p                    = NULL;
  MessageDef* message_p                         = NULL;

  container =
      &pdu->choice.initiatingMessage.value.choice.E_RABModificationIndication;

  S1AP_FIND_PROTOCOLIE_BY_ID(
      S1ap_E_RABModificationIndicationIEs_t, ie, container,
      S1ap_ProtocolIE_ID_id_MME_UE_S1AP_ID, true);
  mme_ue_s1ap_id = ie->value.choice.MME_UE_S1AP_ID;

  S1AP_FIND_PROTOCOLIE_BY_ID(
      S1ap_E_RABModificationIndicationIEs_t, ie, container,
      S1ap_ProtocolIE_ID_id_eNB_UE_S1AP_ID, true);
  // eNB UE S1AP ID is limited to 24 bits
  enb_ue_s1ap_id =
      (enb_ue_s1ap_id_t)(ie->value.choice.ENB_UE_S1AP_ID & ENB_UE_S1AP_ID_MASK);

  if ((ue_ref_p = s1ap_state_get_ue_mmeid((uint32_t) mme_ue_s1ap_id)) == NULL) {
    OAILOG_DEBUG(
        LOG_S1AP,
        "No UE is attached to this mme UE s1ap id: " MME_UE_S1AP_ID_FMT
        " %u(10)\n",
        (uint32_t) mme_ue_s1ap_id, (uint32_t) mme_ue_s1ap_id);
    OAILOG_FUNC_RETURN(LOG_S1AP, RETURNerror);
  }

  if (ue_ref_p->enb_ue_s1ap_id != enb_ue_s1ap_id) {
    OAILOG_DEBUG(
        LOG_S1AP,
        "Mismatch in eNB UE S1AP ID, known: " ENB_UE_S1AP_ID_FMT
        ", received: " ENB_UE_S1AP_ID_FMT "\n",
        ue_ref_p->enb_ue_s1ap_id, enb_ue_s1ap_id);
    OAILOG_FUNC_RETURN(LOG_S1AP, RETURNerror);
  }

  message_p = DEPRECATEDitti_alloc_new_message_fatal(
      TASK_S1AP, S1AP_E_RAB_MODIFICATION_IND);
  S1AP_E_RAB_MODIFICATION_IND(message_p).mme_ue_s1ap_id =
      ue_ref_p->mme_ue_s1ap_id;
  S1AP_E_RAB_MODIFICATION_IND(message_p).enb_ue_s1ap_id =
      ue_ref_p->enb_ue_s1ap_id;

  /** Get the bearers to be modified. */
  S1AP_FIND_PROTOCOLIE_BY_ID(
      S1ap_E_RABModificationIndicationIEs_t, ie, container,
      S1ap_ProtocolIE_ID_id_E_RABToBeModifiedListBearerModInd, true);
  const S1ap_E_RABToBeModifiedListBearerModInd_t* const e_rab_list =
      &ie->value.choice.E_RABToBeModifiedListBearerModInd;
  int num_erab = e_rab_list->list.count;
  for (int index = 0; index < num_erab; index++) {
    const S1ap_E_RABToBeModifiedItemBearerModIndIEs_t* const erab_item_ies =
        (S1ap_E_RABToBeModifiedItemBearerModIndIEs_t*)
            e_rab_list->list.array[index];
    const S1ap_E_RABToBeModifiedItemBearerModInd_t* const erab_item =
        (S1ap_E_RABToBeModifiedItemBearerModInd_t*) &erab_item_ies->value.choice
            .E_RABToBeModifiedItemBearerModInd;
    S1AP_E_RAB_MODIFICATION_IND(message_p)
        .e_rab_to_be_modified_list.item[index]
        .e_rab_id = erab_item->e_RAB_ID;

    bstring transport_layer_address = blk2bstr(
        erab_item->transportLayerAddress.buf,
        erab_item->transportLayerAddress.size);

    S1AP_E_RAB_MODIFICATION_IND(message_p)
        .e_rab_to_be_modified_list.item[index]
        .s1_xNB_fteid.teid = htonl(*((uint32_t*) erab_item->dL_GTP_TEID.buf));

    /** Set the IP address from the FTEID. */
    if (4 == blength(transport_layer_address)) {
      S1AP_E_RAB_MODIFICATION_IND(message_p)
          .e_rab_to_be_modified_list.item[index]
          .s1_xNB_fteid.ipv4 = 1;
      memcpy(
          &S1AP_E_RAB_MODIFICATION_IND(message_p)
               .e_rab_to_be_modified_list.item[index]
               .s1_xNB_fteid.ipv4_address,
          transport_layer_address->data, blength(transport_layer_address));
    } else if (16 == blength(transport_layer_address)) {
      S1AP_E_RAB_MODIFICATION_IND(message_p)
          .e_rab_to_be_modified_list.item[index]
          .s1_xNB_fteid.ipv6 = 1;
      memcpy(
          &S1AP_E_RAB_MODIFICATION_IND(message_p)
               .e_rab_to_be_modified_list.item[index]
               .s1_xNB_fteid.ipv6_address,
          transport_layer_address->data, blength(transport_layer_address));
    } else {
      Fatal("TODO IP address %d bytes", blength(transport_layer_address));
    }
    bdestroy_wrapper(&transport_layer_address);

    S1AP_E_RAB_MODIFICATION_IND(message_p)
        .e_rab_to_be_modified_list.no_of_items++;
  }

  /** Get the bearers not to be modified. */
  S1AP_FIND_PROTOCOLIE_BY_ID(
      S1ap_E_RABModificationIndicationIEs_t, ie, container,
      S1ap_ProtocolIE_ID_id_E_RABNotToBeModifiedListBearerModInd, false);
  if (ie) {
    const S1ap_E_RABNotToBeModifiedListBearerModInd_t* const
        e_rab_not_mod_list =
            &ie->value.choice.E_RABNotToBeModifiedListBearerModInd;
    num_erab = e_rab_not_mod_list->list.count;
    for (int index = 0; index < num_erab; index++) {
      const S1ap_E_RABNotToBeModifiedItemBearerModIndIEs_t* const
          erab_item_ies = (S1ap_E_RABNotToBeModifiedItemBearerModIndIEs_t*)
                              e_rab_not_mod_list->list.array[index];
      const S1ap_E_RABNotToBeModifiedItemBearerModInd_t* const erab_item =
          (S1ap_E_RABNotToBeModifiedItemBearerModInd_t*) &erab_item_ies->value
              .choice.E_RABNotToBeModifiedItemBearerModInd;
      S1AP_E_RAB_MODIFICATION_IND(message_p)
          .e_rab_not_to_be_modified_list.item[index]
          .e_rab_id = erab_item->e_RAB_ID;

      bstring transport_layer_address = blk2bstr(
          erab_item->transportLayerAddress.buf,
          erab_item->transportLayerAddress.size);

      S1AP_E_RAB_MODIFICATION_IND(message_p)
          .e_rab_not_to_be_modified_list.item[index]
          .s1_xNB_fteid.teid = htonl(*((uint32_t*) erab_item->dL_GTP_TEID.buf));

      /** Set the IP address from the FTEID. */
      if (blength(transport_layer_address) == 4) {
        S1AP_E_RAB_MODIFICATION_IND(message_p)
            .e_rab_not_to_be_modified_list.item[index]
            .s1_xNB_fteid.ipv4 = 1;
        memcpy(
            &S1AP_E_RAB_MODIFICATION_IND(message_p)
                 .e_rab_not_to_be_modified_list.item[index]
                 .s1_xNB_fteid.ipv4_address,
            transport_layer_address->data, blength(transport_layer_address));
      } else if (blength(transport_layer_address) == 16) {
        S1AP_E_RAB_MODIFICATION_IND(message_p)
            .e_rab_not_to_be_modified_list.item[index]
            .s1_xNB_fteid.ipv6 = 1;
        memcpy(
            &S1AP_E_RAB_MODIFICATION_IND(message_p)
                 .e_rab_not_to_be_modified_list.item[index]
                 .s1_xNB_fteid.ipv6_address,
            transport_layer_address->data, blength(transport_layer_address));
      } else {
        Fatal("TODO IP address %d bytes", blength(transport_layer_address));
      }
      bdestroy_wrapper(&transport_layer_address);

      S1AP_E_RAB_MODIFICATION_IND(message_p)
          .e_rab_not_to_be_modified_list.no_of_items++;
    }
  }
  rc = send_msg_to_task(&s1ap_task_zmq_ctx, TASK_MME_APP, message_p);
  OAILOG_FUNC_RETURN(LOG_S1AP, rc);
}

//------------------------------------------------------------------------------
void s1ap_mme_generate_erab_modification_confirm(
    s1ap_state_t* state, const itti_s1ap_e_rab_modification_cnf_t* const conf) {
  uint8_t* buffer_p        = NULL;
  uint32_t length          = 0;
  ue_description_t* ue_ref = NULL;
  S1ap_S1AP_PDU_t pdu      = {0};
  S1ap_E_RABModificationConfirm_t* out;
  S1ap_E_RABModificationConfirmIEs_t* ie = NULL;

  OAILOG_FUNC_IN(LOG_S1AP);
  DevAssert(conf != NULL);

  if ((ue_ref = s1ap_state_get_ue_mmeid(conf->mme_ue_s1ap_id)) == NULL) {
    OAILOG_ERROR(
        LOG_S1AP,
        "This mme ue s1ap id (" MME_UE_S1AP_ID_FMT
        ") is not attached to any UE context\n",
        conf->mme_ue_s1ap_id);
    OAILOG_FUNC_OUT(LOG_S1AP);
  }

  memset(&pdu, 0, sizeof(pdu));
  pdu.present = S1ap_S1AP_PDU_PR_successfulOutcome;
  pdu.choice.successfulOutcome.procedureCode =
      S1ap_ProcedureCode_id_E_RABModificationIndication;
  pdu.choice.successfulOutcome.criticality = S1ap_Criticality_reject;
  pdu.choice.successfulOutcome.value.present =
      S1ap_SuccessfulOutcome__value_PR_E_RABModificationConfirm;
  out = &pdu.choice.successfulOutcome.value.choice.E_RABModificationConfirm;

  /* mandatory */
  ie = (S1ap_E_RABModificationConfirmIEs_t*) calloc(
      1, sizeof(S1ap_E_RABModificationConfirmIEs_t));
  ie->id            = S1ap_ProtocolIE_ID_id_MME_UE_S1AP_ID;
  ie->criticality   = S1ap_Criticality_ignore;
  ie->value.present = S1ap_E_RABModificationConfirmIEs__value_PR_MME_UE_S1AP_ID;
  ie->value.choice.MME_UE_S1AP_ID = conf->mme_ue_s1ap_id;
  ASN_SEQUENCE_ADD(&out->protocolIEs.list, ie);

  /* mandatory */
  ie = (S1ap_E_RABModificationConfirmIEs_t*) calloc(
      1, sizeof(S1ap_E_RABModificationConfirmIEs_t));
  ie->id            = S1ap_ProtocolIE_ID_id_eNB_UE_S1AP_ID;
  ie->criticality   = S1ap_Criticality_ignore;
  ie->value.present = S1ap_E_RABModificationConfirmIEs__value_PR_ENB_UE_S1AP_ID;
  ie->value.choice.ENB_UE_S1AP_ID = conf->enb_ue_s1ap_id;
  ASN_SEQUENCE_ADD(&out->protocolIEs.list, ie);

  if (conf->e_rab_modify_list.no_of_items) {
    ie = (S1ap_E_RABModificationConfirmIEs_t*) calloc(
        1, sizeof(S1ap_E_RABModificationConfirmIEs_t));
    ie->id          = S1ap_ProtocolIE_ID_id_E_RABModifyListBearerModConf;
    ie->criticality = S1ap_Criticality_reject;
    ie->value.present =
        S1ap_E_RABModificationConfirmIEs__value_PR_E_RABModifyListBearerModConf;
    ASN_SEQUENCE_ADD(&out->protocolIEs.list, ie);

    S1ap_E_RABModifyListBearerModConf_t* e_rabmodifylistbearermodconf =
        &ie->value.choice.E_RABModifyListBearerModConf;

    for (int i = 0; i < conf->e_rab_modify_list.no_of_items; i++) {
      S1ap_E_RABModifyItemBearerModConfIEs_t* item =
          calloc(1, sizeof(S1ap_E_RABModifyItemBearerModConfIEs_t));

      item->id          = S1ap_ProtocolIE_ID_id_E_RABModifyItemBearerModConf;
      item->criticality = S1ap_Criticality_reject;
      item->value.present =
          S1ap_E_RABModifyItemBearerModConfIEs__value_PR_E_RABModifyItemBearerModConf;

      S1ap_E_RABModifyItemBearerModConf_t* bearer =
          &item->value.choice.E_RABModifyItemBearerModConf;
      bearer->e_RAB_ID = conf->e_rab_modify_list.e_rab_id[i];

      ASN_SEQUENCE_ADD(&e_rabmodifylistbearermodconf->list, item);
    }
  }

  if (s1ap_mme_encode_pdu(&pdu, &buffer_p, &length) < 0) {
    OAILOG_ERROR(
        LOG_S1AP, "Encoding of S1ap_E_RABModificationConfirmIEs_t failed \n");
    OAILOG_FUNC_OUT(LOG_S1AP);
  }

  OAILOG_NOTICE(
      LOG_S1AP,
      "Send S1AP E_RAB_MODIFICATION_CONFIRM Command message MME_UE_S1AP_ID "
      "= " MME_UE_S1AP_ID_FMT " eNB_UE_S1AP_ID = " ENB_UE_S1AP_ID_FMT "\n",
      (mme_ue_s1ap_id_t) ue_ref->mme_ue_s1ap_id,
      (enb_ue_s1ap_id_t) ue_ref->enb_ue_s1ap_id);
  bstring b = blk2bstr(buffer_p, length);
  free(buffer_p);
  s1ap_mme_itti_send_sctp_request(
      &b, ue_ref->sctp_assoc_id, ue_ref->sctp_stream_send,
      ue_ref->mme_ue_s1ap_id);
  OAILOG_FUNC_OUT(LOG_S1AP);
}

//----------------------------------------------------------------
status_code_e s1ap_mme_handle_enb_configuration_transfer(
    s1ap_state_t* state, const sctp_assoc_id_t assoc_id,
    const sctp_stream_id_t stream, S1ap_S1AP_PDU_t* pdu) {
  S1ap_ENBConfigurationTransfer_t* container = NULL;
  S1ap_ENBConfigurationTransferIEs_t* ie     = NULL;
  S1ap_TargeteNB_ID_t* targeteNB_ID          = NULL;
  uint8_t* enb_id_buf                        = NULL;
  enb_description_t* enb_association         = NULL;
  enb_description_t* target_enb_association  = NULL;
  hashtable_element_array_t* enb_array       = NULL;
  uint32_t target_enb_id                     = 0;
  uint8_t* buffer                            = NULL;
  uint32_t length                            = 0;
  uint32_t idx                               = 0;
  int rc                                     = RETURNok;

  // Not done according to Rel-15 (Target TAI and Source TAI)
  OAILOG_FUNC_IN(LOG_S1AP);
  container =
      &pdu->choice.initiatingMessage.value.choice.ENBConfigurationTransfer;

  S1AP_FIND_PROTOCOLIE_BY_ID(
      S1ap_ENBConfigurationTransferIEs_t, ie, container,
      S1ap_ProtocolIE_ID_id_SONConfigurationTransferECT, false);

  OAILOG_DEBUG(
      LOG_S1AP, "Recieved eNB Confiuration Request from assoc_id %u\n",
      assoc_id);
  enb_association = s1ap_state_get_enb(state, assoc_id);
  if (enb_association == NULL) {
    OAILOG_ERROR(
        LOG_S1AP, "Ignoring eNB Confiuration Request from unknown assoc %u\n",
        assoc_id);
    OAILOG_FUNC_RETURN(LOG_S1AP, RETURNerror);
  }

  if (enb_association->s1_state != S1AP_READY) {
    // ignore the message if s1 not ready
    OAILOG_INFO(
        LOG_S1AP,
        "S1 setup is not done.Invalid state.Ignoring eNB Configuration Request "
        "eNB Id = %d , S1AP state = %d \n",
        enb_association->enb_id, enb_association->s1_state);
    OAILOG_FUNC_RETURN(LOG_S1AP, RETURNok);
  }

  targeteNB_ID = &ie->value.choice.SONConfigurationTransfer.targeteNB_ID;

  if (targeteNB_ID->global_ENB_ID.eNB_ID.present == S1ap_ENB_ID_PR_homeENB_ID) {
    // Home eNB ID = 28 bits
    enb_id_buf = targeteNB_ID->global_ENB_ID.eNB_ID.choice.homeENB_ID.buf;

    target_enb_id = (enb_id_buf[0] << 20) + (enb_id_buf[1] << 12) +
                    (enb_id_buf[2] << 4) + ((enb_id_buf[3] & 0xf0) >> 4);
    OAILOG_INFO(LOG_S1AP, "home eNB id: %u\n", target_enb_id);
  } else {
    // Macro eNB = 20 bits
    enb_id_buf = targeteNB_ID->global_ENB_ID.eNB_ID.choice.macroENB_ID.buf;

    target_enb_id = (enb_id_buf[0] << 12) + (enb_id_buf[1] << 4) +
                    ((enb_id_buf[2] & 0xf0) >> 4);
    OAILOG_INFO(LOG_S1AP, "macro eNB id: %u\n", target_enb_id);
  }

  // retrieve enb_description using hash table and match target_enb_id
  if ((enb_array = hashtable_ts_get_elements(&state->enbs)) != NULL) {
    for (idx = 0; idx < enb_array->num_elements; idx++) {
      target_enb_association =
          (enb_description_t*) (uintptr_t) enb_array->elements[idx];
      if (target_enb_association->enb_id == target_enb_id) {
        break;
      }
    }
    free_wrapper((void**) &enb_array->elements);
    free_wrapper((void**) &enb_array);
    if (target_enb_association->enb_id != target_enb_id) {
      OAILOG_ERROR(LOG_S1AP, "No eNB for enb_id %d\n", target_enb_id);
      OAILOG_FUNC_RETURN(LOG_S1AP, RETURNerror);
    }
  }

  pdu->choice.initiatingMessage.procedureCode =
      S1ap_ProcedureCode_id_MMEConfigurationTransfer;
  pdu->present = S1ap_S1AP_PDU_PR_initiatingMessage;
  // Message is received and immediately sent back by changing only the IE type
  // which is different from the usual approach of creating a new message.
  ie->id = S1ap_ProtocolIE_ID_id_SONConfigurationTransferMCT;
  // Encode message
  int enc_rval = s1ap_mme_encode_pdu(pdu, &buffer, &length);
  if (enc_rval < 0) {
    OAILOG_ERROR(
        LOG_S1AP,
        "Failed to encode MME Configuration Transfer message for enb_id %u\n",
        target_enb_id);
    OAILOG_FUNC_RETURN(LOG_S1AP, RETURNerror);
  }

  bstring b = blk2bstr(buffer, length);
  free(buffer);

  // Send message
  rc = s1ap_mme_itti_send_sctp_request(
      &b, target_enb_association->sctp_assoc_id,
      0,   // Stream id 0 for non UE related S1AP message
      0);  // mme_ue_s1ap_id 0 because UE in idle

  if (rc != RETURNok) {
    OAILOG_ERROR(
        LOG_S1AP,
        "Failed to send MME Configuration Transfer message over sctp for"
        "enb_id %u\n",
        target_enb_id);
  } else {
    OAILOG_INFO(
        LOG_S1AP,
        "Sent MME Configuration Transfer message over sctp for "
        "target_enb_id %u\n",
        target_enb_id);
  }
  OAILOG_FUNC_RETURN(LOG_S1AP, rc);
}

//------------------------------------------------------------------------------
bool is_all_erabId_same(S1ap_PathSwitchRequest_t* container) {
  S1ap_PathSwitchRequestIEs_t* ie                                = NULL;
  S1ap_E_RABToBeSwitchedDLItemIEs_t* eRABToBeSwitchedDlItemIEs_p = NULL;
  uint8_t item                                                   = 0;
  uint8_t firstItem                                              = 0;
  uint8_t rc                                                     = true;

  S1AP_FIND_PROTOCOLIE_BY_ID(
      S1ap_PathSwitchRequestIEs_t, ie, container,
      S1ap_ProtocolIE_ID_id_E_RABToBeSwitchedDLList, true);
  if (!ie) {
    OAILOG_ERROR(LOG_S1AP, "Incorrect IE \n");
    return RETURNerror;
  }
  S1ap_E_RABToBeSwitchedDLList_t* e_rab_to_be_switched_dl_list =
      &ie->value.choice.E_RABToBeSwitchedDLList;

  if (1 == e_rab_to_be_switched_dl_list->list.count) {
    rc = false;
    OAILOG_FUNC_RETURN(LOG_S1AP, rc);
  }

  eRABToBeSwitchedDlItemIEs_p = (S1ap_E_RABToBeSwitchedDLItemIEs_t*)
                                    e_rab_to_be_switched_dl_list->list.array[0];
  firstItem = eRABToBeSwitchedDlItemIEs_p->value.choice.E_RABToBeSwitchedDLItem
                  .e_RAB_ID;

  for (item = 1; item < e_rab_to_be_switched_dl_list->list.count; ++item) {
    eRABToBeSwitchedDlItemIEs_p =
        (S1ap_E_RABToBeSwitchedDLItemIEs_t*)
            e_rab_to_be_switched_dl_list->list.array[item];
    if (firstItem == eRABToBeSwitchedDlItemIEs_p->value.choice
                         .E_RABToBeSwitchedDLItem.e_RAB_ID) {
      continue;
    } else {
      rc = false;
      break;
    }
  }
  OAILOG_FUNC_RETURN(LOG_S1AP, rc);
}
//------------------------------------------------------------------------------
status_code_e s1ap_handle_path_switch_req_ack(
    s1ap_state_t* state,
    const itti_s1ap_path_switch_request_ack_t* path_switch_req_ack_p,
    imsi64_t imsi64) {
  OAILOG_FUNC_IN(LOG_S1AP);

  uint8_t* buffer                            = NULL;
  uint32_t length                            = 0;
  ue_description_t* ue_ref_p                 = NULL;
  S1ap_S1AP_PDU_t pdu                        = {0};
  S1ap_PathSwitchRequestAcknowledge_t* out   = NULL;
  S1ap_PathSwitchRequestAcknowledgeIEs_t* ie = NULL;
  int rc                                     = RETURNok;

  if ((ue_ref_p = s1ap_state_get_ue_mmeid(
           path_switch_req_ack_p->mme_ue_s1ap_id)) == NULL) {
    OAILOG_DEBUG_UE(
        LOG_S1AP, imsi64,
        "could not get ue context for mme_ue_s1ap_id " MME_UE_S1AP_ID_FMT "\n",
        (uint32_t) path_switch_req_ack_p->mme_ue_s1ap_id);
    OAILOG_FUNC_RETURN(LOG_S1AP, RETURNerror);
  }

  memset(&pdu, 0, sizeof(pdu));
  pdu.present = S1ap_S1AP_PDU_PR_successfulOutcome;
  pdu.choice.initiatingMessage.procedureCode =
      S1ap_ProcedureCode_id_PathSwitchRequest;
  pdu.choice.successfulOutcome.criticality = S1ap_Criticality_ignore;
  pdu.choice.successfulOutcome.value.present =
      S1ap_SuccessfulOutcome__value_PR_PathSwitchRequestAcknowledge;
  out = &pdu.choice.successfulOutcome.value.choice.PathSwitchRequestAcknowledge;

  ie = (S1ap_PathSwitchRequestAcknowledgeIEs_t*) calloc(
      1, sizeof(S1ap_PathSwitchRequestAcknowledgeIEs_t));
  ie->id          = S1ap_ProtocolIE_ID_id_MME_UE_S1AP_ID;
  ie->criticality = S1ap_Criticality_reject;
  ie->value.present =
      S1ap_PathSwitchRequestAcknowledgeIEs__value_PR_MME_UE_S1AP_ID;
  ie->value.choice.MME_UE_S1AP_ID = ue_ref_p->mme_ue_s1ap_id;
  ASN_SEQUENCE_ADD(&out->protocolIEs.list, ie);

  /* mandatory */
  ie = (S1ap_PathSwitchRequestAcknowledgeIEs_t*) calloc(
      1, sizeof(S1ap_PathSwitchRequestAcknowledgeIEs_t));
  ie->id          = S1ap_ProtocolIE_ID_id_eNB_UE_S1AP_ID;
  ie->criticality = S1ap_Criticality_reject;
  ie->value.present =
      S1ap_PathSwitchRequestAcknowledgeIEs__value_PR_ENB_UE_S1AP_ID;
  ie->value.choice.ENB_UE_S1AP_ID = ue_ref_p->enb_ue_s1ap_id;
  ASN_SEQUENCE_ADD(&out->protocolIEs.list, ie);

  /** Add the security context. */
  ie = (S1ap_PathSwitchRequestAcknowledgeIEs_t*) calloc(
      1, sizeof(S1ap_PathSwitchRequestAcknowledgeIEs_t));
  ie->id          = S1ap_ProtocolIE_ID_id_SecurityContext;
  ie->criticality = S1ap_Criticality_reject;
  ie->value.present =
      S1ap_PathSwitchRequestAcknowledgeIEs__value_PR_SecurityContext;
  if (path_switch_req_ack_p->nh) {
    ie->value.choice.SecurityContext.nextHopParameter.buf =
        calloc(AUTH_NEXT_HOP_SIZE, sizeof(uint8_t));
    memcpy(
        ie->value.choice.SecurityContext.nextHopParameter.buf,
        path_switch_req_ack_p->nh, AUTH_NEXT_HOP_SIZE);
    ie->value.choice.SecurityContext.nextHopParameter.size = AUTH_NEXT_HOP_SIZE;
  } else {
    OAILOG_WARNING(LOG_S1AP, "No nh for PSReqAck.\n");
    ie->value.choice.SecurityContext.nextHopParameter.buf  = NULL;
    ie->value.choice.SecurityContext.nextHopParameter.size = 0;
  }
  ie->value.choice.SecurityContext.nextHopParameter.bits_unused = 0;
  ie->value.choice.SecurityContext.nextHopChainingCount =
      path_switch_req_ack_p->ncc;
  ASN_SEQUENCE_ADD(&out->protocolIEs.list, ie);

  if (s1ap_mme_encode_pdu(&pdu, &buffer, &length) < 0) {
    OAILOG_ERROR_UE(
        LOG_S1AP, imsi64, "Path Switch Request Ack encoding failed \n");
    OAILOG_FUNC_RETURN(LOG_S1AP, RETURNerror);
  }
  bstring b = blk2bstr(buffer, length);
  free_wrapper((void**) &buffer);
  OAILOG_DEBUG_UE(
      LOG_S1AP, imsi64,
      "Send PATH_SWITCH_REQUEST_ACK, mme_ue_s1ap_id " MME_UE_S1AP_ID_FMT "\n",
      (uint32_t) path_switch_req_ack_p->mme_ue_s1ap_id);

  rc = s1ap_mme_itti_send_sctp_request(
      &b, path_switch_req_ack_p->sctp_assoc_id, ue_ref_p->sctp_stream_send,
      path_switch_req_ack_p->mme_ue_s1ap_id);

  OAILOG_FUNC_RETURN(LOG_S1AP, rc);
}
//------------------------------------------------------------------------------
status_code_e s1ap_handle_path_switch_req_failure(
    const itti_s1ap_path_switch_request_failure_t* path_switch_req_failure_p,
    imsi64_t imsi64) {
  S1ap_PathSwitchRequestFailure_t* container = NULL;
  uint8_t* buffer                            = NULL;
  uint32_t length                            = 0;
  ue_description_t* ue_ref_p                 = NULL;
  S1ap_S1AP_PDU_t pdu                        = {0};
  S1ap_PathSwitchRequestFailureIEs_t* ie     = NULL;
  int rc                                     = RETURNok;
  mme_ue_s1ap_id_t mme_ue_s1ap_id            = 0;
  OAILOG_FUNC_IN(LOG_S1AP);

  mme_ue_s1ap_id = path_switch_req_failure_p->mme_ue_s1ap_id;
  ue_ref_p       = s1ap_state_get_ue_mmeid(mme_ue_s1ap_id);
  if (ue_ref_p == NULL) {
    OAILOG_DEBUG_UE(
        LOG_S1AP, imsi64,
        "could not get ue context for mme_ue_s1ap_id " MME_UE_S1AP_ID_FMT "\n",
        mme_ue_s1ap_id);
    OAILOG_FUNC_RETURN(LOG_S1AP, RETURNerror);
  }

  memset(&pdu, 0, sizeof(pdu));
  pdu.choice.unsuccessfulOutcome.procedureCode =
      S1ap_ProcedureCode_id_PathSwitchRequest;
  pdu.present = S1ap_S1AP_PDU_PR_unsuccessfulOutcome;
  pdu.choice.unsuccessfulOutcome.criticality = S1ap_Criticality_reject;
  pdu.choice.unsuccessfulOutcome.value.present =
      S1ap_UnsuccessfulOutcome__value_PR_PathSwitchRequestFailure;
  container =
      &pdu.choice.unsuccessfulOutcome.value.choice.PathSwitchRequestFailure;

  ie = (S1ap_PathSwitchRequestFailureIEs_t*) calloc(
      1, sizeof(S1ap_PathSwitchRequestFailureIEs_t));
  ie->id            = S1ap_ProtocolIE_ID_id_MME_UE_S1AP_ID;
  ie->criticality   = S1ap_Criticality_reject;
  ie->value.present = S1ap_PathSwitchRequestFailureIEs__value_PR_MME_UE_S1AP_ID;
  ie->value.choice.MME_UE_S1AP_ID = path_switch_req_failure_p->mme_ue_s1ap_id;
  s1ap_mme_set_cause(
      &ie->value.choice.Cause, S1ap_Cause_PR_radioNetwork,
      S1ap_CauseRadioNetwork_ho_failure_in_target_EPC_eNB_or_target_system);
  ASN_SEQUENCE_ADD(&container->protocolIEs.list, ie);

  ie = (S1ap_PathSwitchRequestFailureIEs_t*) calloc(
      1, sizeof(S1ap_PathSwitchRequestFailureIEs_t));
  ie->id            = S1ap_ProtocolIE_ID_id_eNB_UE_S1AP_ID;
  ie->criticality   = S1ap_Criticality_reject;
  ie->value.present = S1ap_PathSwitchRequestFailureIEs__value_PR_ENB_UE_S1AP_ID;
  ie->value.choice.ENB_UE_S1AP_ID = path_switch_req_failure_p->enb_ue_s1ap_id;
  s1ap_mme_set_cause(
      &ie->value.choice.Cause, S1ap_Cause_PR_radioNetwork,
      S1ap_CauseRadioNetwork_ho_failure_in_target_EPC_eNB_or_target_system);
  ASN_SEQUENCE_ADD(&container->protocolIEs.list, ie);

  if (s1ap_mme_encode_pdu(&pdu, &buffer, &length) < 0) {
    OAILOG_ERROR(LOG_S1AP, "Path Switch Request Failure encoding failed \n");
    OAILOG_FUNC_RETURN(LOG_S1AP, RETURNerror);
  }
  bstring b = blk2bstr(buffer, length);
  free_wrapper((void**) &buffer);
  OAILOG_DEBUG_UE(
      LOG_S1AP, imsi64,
      "send PATH_SWITCH_REQUEST_Failure for mme_ue_s1ap_id " MME_UE_S1AP_ID_FMT
      "\n",
      (uint32_t) path_switch_req_failure_p->mme_ue_s1ap_id);

  rc = s1ap_mme_itti_send_sctp_request(
      &b, path_switch_req_failure_p->sctp_assoc_id, ue_ref_p->sctp_stream_send,
      path_switch_req_failure_p->mme_ue_s1ap_id);

  OAILOG_FUNC_RETURN(LOG_S1AP, rc);
}

const char* s1_enb_state2str(enum mme_s1_enb_state_s state) {
  switch (state) {
    case S1AP_INIT:
      return "S1AP_INIT";
    case S1AP_RESETING:
      return "S1AP_RESETING";
    case S1AP_READY:
      return "S1AP_READY";
    case S1AP_SHUTDOWN:
      return "S1AP_SHUTDOWN";
    default:
      return "unknown s1ap_enb_state";
  }
}

const char* s1ap_direction2str(uint8_t dir) {
  switch (dir) {
    case S1ap_S1AP_PDU_PR_NOTHING:
      return "<nothing>";
    case S1ap_S1AP_PDU_PR_initiatingMessage:
      return "originating message";
    case S1ap_S1AP_PDU_PR_successfulOutcome:
      return "successful outcome";
    case S1ap_S1AP_PDU_PR_unsuccessfulOutcome:
      return "unsuccessful outcome";
    default:
      return "unknown direction";
  }
}

//------------------------------------------------------------------------------
status_code_e s1ap_mme_handle_erab_rel_response(
    s1ap_state_t* state, const sctp_assoc_id_t assoc_id,
    const sctp_stream_id_t stream, S1ap_S1AP_PDU_t* pdu) {
  OAILOG_FUNC_IN(LOG_S1AP);
  S1ap_E_RABReleaseResponseIEs_t* ie     = NULL;
  S1ap_E_RABReleaseResponse_t* container = NULL;
  ue_description_t* ue_ref_p             = NULL;
  MessageDef* message_p                  = NULL;
  int rc                                 = RETURNok;
  imsi64_t imsi64                        = INVALID_IMSI64;
  enb_ue_s1ap_id_t enb_ue_s1ap_id        = INVALID_ENB_UE_S1AP_ID;
  mme_ue_s1ap_id_t mme_ue_s1ap_id        = INVALID_MME_UE_S1AP_ID;

  container = &pdu->choice.successfulOutcome.value.choice.E_RABReleaseResponse;

  S1AP_FIND_PROTOCOLIE_BY_ID(
      S1ap_E_RABReleaseResponseIEs_t, ie, container,
      S1ap_ProtocolIE_ID_id_MME_UE_S1AP_ID, true);
  mme_ue_s1ap_id = ie->value.choice.MME_UE_S1AP_ID;

  if ((ie) &&
      (ue_ref_p = s1ap_state_get_ue_mmeid((uint32_t) mme_ue_s1ap_id)) == NULL) {
    OAILOG_ERROR(
        LOG_S1AP,
        "No UE is attached to this mme UE s1ap id: " MME_UE_S1AP_ID_FMT "\n",
        (mme_ue_s1ap_id_t) mme_ue_s1ap_id);
    OAILOG_FUNC_RETURN(LOG_S1AP, RETURNerror);
  }

  S1AP_FIND_PROTOCOLIE_BY_ID(
      S1ap_E_RABReleaseResponseIEs_t, ie, container,
      S1ap_ProtocolIE_ID_id_eNB_UE_S1AP_ID, true);
  // eNB UE S1AP ID is limited to 24 bits
  enb_ue_s1ap_id =
      (enb_ue_s1ap_id_t)(ie->value.choice.ENB_UE_S1AP_ID & ENB_UE_S1AP_ID_MASK);

  if ((ie) && ue_ref_p->enb_ue_s1ap_id != enb_ue_s1ap_id) {
    OAILOG_ERROR(
        LOG_S1AP,
        "Mismatch in eNB UE S1AP ID, known: " ENB_UE_S1AP_ID_FMT
        ", received: " ENB_UE_S1AP_ID_FMT "\n",
        ue_ref_p->enb_ue_s1ap_id, (enb_ue_s1ap_id_t) enb_ue_s1ap_id);
    OAILOG_FUNC_RETURN(LOG_S1AP, RETURNerror);
  }

  s1ap_imsi_map_t* imsi_map = get_s1ap_imsi_map();
  hashtable_uint64_ts_get(
      imsi_map->mme_ue_id_imsi_htbl,
      (const hash_key_t) ie->value.choice.MME_UE_S1AP_ID, &imsi64);

  message_p = itti_alloc_new_message(TASK_S1AP, S1AP_E_RAB_REL_RSP);
  if (message_p == NULL) {
    OAILOG_ERROR(LOG_S1AP, "itti_alloc_new_message Failed\n");
    OAILOG_FUNC_RETURN(LOG_S1AP, RETURNerror);
  }
  S1AP_E_RAB_REL_RSP(message_p).mme_ue_s1ap_id = ue_ref_p->mme_ue_s1ap_id;
  S1AP_E_RAB_REL_RSP(message_p).enb_ue_s1ap_id = ue_ref_p->enb_ue_s1ap_id;
  S1AP_E_RAB_REL_RSP(message_p).e_rab_rel_list.no_of_items           = 0;
  S1AP_E_RAB_REL_RSP(message_p).e_rab_failed_to_rel_list.no_of_items = 0;

  const S1ap_E_RABList_t* const e_rab_list = &ie->value.choice.E_RABList;
  int num_erab                             = e_rab_list->list.count;
  if (ie) {
    for (int index = 0; index < num_erab; index++) {
      const S1ap_E_RABItemIEs_t* const erab_item_ies =
          (S1ap_E_RABItemIEs_t*) e_rab_list->list.array[index];
      const S1ap_E_RABItem_t* const erab_item =
          (S1ap_E_RABItem_t*) &erab_item_ies->value.choice.E_RABItem;
      S1AP_E_RAB_REL_RSP(message_p).e_rab_rel_list.item[index].e_rab_id =
          erab_item->e_RAB_ID;
      S1AP_E_RAB_REL_RSP(message_p).e_rab_rel_list.item[index].cause =
          erab_item->cause;
      S1AP_E_RAB_REL_RSP(message_p).e_rab_rel_list.no_of_items++;
    }
  }
  if (ie) {
    for (int index = 0; index < num_erab; index++) {
      const S1ap_E_RABItemIEs_t* const erab_item_ies =
          (S1ap_E_RABItemIEs_t*) e_rab_list->list.array[index];
      const S1ap_E_RABItem_t* const erab_item =
          (S1ap_E_RABItem_t*) &erab_item_ies->value.choice.E_RABItem;
      S1AP_E_RAB_REL_RSP(message_p)
          .e_rab_failed_to_rel_list.item[index]
          .e_rab_id = erab_item->e_RAB_ID;
      S1AP_E_RAB_REL_RSP(message_p).e_rab_failed_to_rel_list.item[index].cause =
          erab_item->cause;
      S1AP_E_RAB_REL_RSP(message_p).e_rab_failed_to_rel_list.no_of_items++;
    }
  }
  message_p->ittiMsgHeader.imsi = imsi64;
  rc = send_msg_to_task(&s1ap_task_zmq_ctx, TASK_MME_APP, message_p);
  OAILOG_FUNC_RETURN(LOG_S1AP, rc);
}

status_code_e s1ap_mme_remove_stale_ue_context(
    enb_ue_s1ap_id_t enb_ue_s1ap_id, uint32_t enb_id) {
  OAILOG_FUNC_IN(LOG_S1AP);
  MessageDef* message_p = NULL;
  message_p = itti_alloc_new_message(TASK_S1AP, S1AP_REMOVE_STALE_UE_CONTEXT);
  if (!message_p) {
    OAILOG_ERROR(
        LOG_S1AP,
        "Failed to allocate memory for S1AP_REMOVE_STALE_UE_CONTEXT \n");
    OAILOG_FUNC_RETURN(LOG_S1AP, RETURNerror);
  }
  S1AP_REMOVE_STALE_UE_CONTEXT(message_p).enb_ue_s1ap_id = enb_ue_s1ap_id;
  S1AP_REMOVE_STALE_UE_CONTEXT(message_p).enb_id         = enb_id;
  OAILOG_INFO(
      LOG_S1AP,
      "sent S1AP_REMOVE_STALE_UE_CONTEXT for enb_ue_s1ap_id " ENB_UE_S1AP_ID_FMT
      "\n",
      enb_ue_s1ap_id);
  send_msg_to_task(&s1ap_task_zmq_ctx, TASK_MME_APP, message_p);
  OAILOG_FUNC_RETURN(LOG_S1AP, RETURNok);
}

status_code_e s1ap_send_mme_ue_context_release(
    s1ap_state_t* state, ue_description_t* ue_ref_p,
    enum s1cause s1_release_cause, S1ap_Cause_t ie_cause, imsi64_t imsi64) {
  MessageDef* message_p = NULL;
  message_p = itti_alloc_new_message(TASK_S1AP, S1AP_UE_CONTEXT_RELEASE_REQ);
  if (!message_p) {
    OAILOG_ERROR(
        LOG_S1AP,
        "Failed to allocate memory for S1AP_REMOVE_STALE_UE_CONTEXT \n");
    OAILOG_FUNC_RETURN(LOG_S1AP, RETURNerror);
  }

  enb_description_t* enb_ref_p =
      s1ap_state_get_enb(state, ue_ref_p->sctp_assoc_id);

  S1AP_UE_CONTEXT_RELEASE_REQ(message_p).mme_ue_s1ap_id =
      ue_ref_p->mme_ue_s1ap_id;
  S1AP_UE_CONTEXT_RELEASE_REQ(message_p).enb_ue_s1ap_id =
      ue_ref_p->enb_ue_s1ap_id;
  S1AP_UE_CONTEXT_RELEASE_REQ(message_p).enb_id   = enb_ref_p->enb_id;
  S1AP_UE_CONTEXT_RELEASE_REQ(message_p).relCause = s1_release_cause;
  S1AP_UE_CONTEXT_RELEASE_REQ(message_p).cause    = ie_cause;

  message_p->ittiMsgHeader.imsi = imsi64;
  return send_msg_to_task(&s1ap_task_zmq_ctx, TASK_MME_APP, message_p);
}

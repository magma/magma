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

#include <stdlib.h>
#include <stdio.h>
#include <stdbool.h>
#include <stdint.h>
#include <netinet/in.h>
#include <string.h>
#include <sys/types.h>

#include "bstrlib.h"
#include "hashtable.h"
#include "log.h"
#include "assertions.h"
#include "conversions.h"
#include "intertask_interface.h"
#include "timer.h"
#include "dynamic_memory_check.h"
#include "amf_config.h"
#include "ngap_common.h"
#include "ngap_amf_encoder.h"
#include "ngap_amf_nas_procedures.h"
#include "ngap_amf_itti_messaging.h"
#include "ngap_amf.h"
#include "ngap_amf_ta.h"
#include "ngap_amf_handlers.h"
#include "3gpp_23.003.h"
#include "3gpp_24.008.h"
#include "3gpp_38.401.h"
#include "3gpp_38.413.h"
#include "BIT_STRING.h"
#include "INTEGER.h"
#include "Ngap_NGAP-PDU.h"
#include "Ngap_CauseMisc.h"
#include "Ngap_CauseNas.h"
#include "Ngap_CauseProtocol.h"
#include "Ngap_CauseRadioNetwork.h"
#include "Ngap_CauseTransport.h"
#include "Ngap_GNB-ID.h"
#include "Ngap_RAN-UE-NGAP-ID.h"
#include "Ngap_RANNodeName.h"
#include "Ngap_GTP-TEID.h"
#include "Ngap_GlobalGNB-ID.h"
#include "Ngap_AMF-UE-NGAP-ID.h"
#include "Ngap_PLMNIdentity.h"
#include "Ngap_ProcedureCode.h"
#include "Ngap_ResetType.h"
#include "Ngap_FiveG-S-TMSI.h"
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
#include "asn_SEQUENCE_OF.h"
#include "common_defs.h"
#include "intertask_interface_types.h"
#include "itti_types.h"
#include "amf_app_messages_types.h"
#include "service303.h"
#include "ngap_state.h"

struct Ngap_IE;

int ngap_generate_ng_setup_response(
    ngap_state_t* state, gnb_description_t* gnb_association);

int ngap_amf_generate_ue_context_release_command(
    ngap_state_t* state, m5g_ue_description_t* ue_ref_p, enum Ngcause,
    imsi64_t imsi64);

/* Handlers matrix. Only amf related procedures present here.
 */
ngap_message_handler_t ngap_message_handlers[][3] = {
    {0, 0, 0}, /* HandoverPreparation */
    {0, 0, 0}, /* HandoverResourceAllocation */
    {0, 0, 0}, /* HandoverNotification */
    {/*ngap_amf_handle_path_switch_request*/ 0, 0,
     0},       /* TODO PathSwitchRequest */
    {0, 0, 0}, /* HandoverCancel */
    {0, 0 /*ngap_amf_handle_pduSession_setup_response*/,
     0 /*ngap_amf_handle_pduSession_setup_failure*/}, /* TODO PduSession*/
    {0, 0, 0},                                        /*Padding */
    {0, 0, 0},                                        /*Padding */
    {0, 0, 0},                                        /* PduSessionModify */
    {0, /*ngap_amf_handle_pduSession_release_response*/ 0,
     0},                                       /*TODO  PduSessionRelease */
    {0, 0, 0},                                 /* PduSessionReleaseIndication */
    {0, 0, 0},                                 /* Paging */
    {/*ngap_amf_handle_gnb_reset*/ 0, 0, 0},   /* TODO Reset */
    {ngap_amf_handle_error_ind_message, 0, 0}, /*  ErrorIndication */
    {0, ngap_amf_handle_initial_context_setup_response,
     ngap_amf_handle_initial_context_setup_failure}, /* InitialContextSetup */
    {ngap_amf_handle_initial_ue_message, 0, 0},      /* initialUEMessage */
    {ngap_amf_handle_nas_non_delivery, 0, 0}, /* NASNonDeliveryIndication */
    {/*ngap_amf_handle_ue_context_release_request*/ 0, 0,
     0},       /* UEContextReleaseRequest */
    {0, 0, 0}, /* DownlinkNgcdma2000tunneling */
    {0, 0, 0}, /* UplinkNgcdma2000tunneling */
    {0, /*ngap_amf_handle_ue_context_modification_response*/ 0,
     /*ngap_amf_handle_ue_context_modification_failure*/
     0},                                      /* UEContextModification
                                               */
    {ngap_amf_handle_ng_setup_request, 0, 0}, /* NGSetup */
    {/*ngap_amf_handle_ue_cap_indication*/ 0, 0,
     0},       /* TODO UECapabilityInfoIndication */
    {0, 0, 0}, /* gNBStatusTransfer */
    {0, 0, 0}, /* AMFStatusTransfer */
    {0, 0, 0}, /* DeactivateTrace */
    {0, 0, 0}, /* TraceStart */
    {0, 0, 0}, /* TraceFailureIndication */
    {0, 0, 0}, /* GNBConfigurationUpdate */
    {0, 0, 0}, /* AMFConfigurationUpdate */
    {0, 0, 0}, /* LocationReportingControl */
    {0, 0, 0}, /* LocationReportingFailureIndication */
    {0, 0, 0}, /* LocationReport */
    {0, 0, 0}, /* OverloadStart */
    {0, 0, 0}, /* OverloadStop */
    {0, 0, 0}, /* WriteReplaceWarning */
    {0, 0, 0}, /* gNBDirectInformationTransfer */
    {0, 0, 0}, /* AMFDirectInformationTransfer */
    {0, 0, 0}, /* PrivateMessage */
    {/*ngap_amf_handle_gnb_configuration_transfer*/ 0, 0,
     0},       /* TODO gNBConfigurationTransfer */
    {0, 0, 0}, /* AMFConfigurationTransfer */
    {/*ngap_amf_handle_ue_context_release_request*/ 0,
     /*ngap_amf_handle_ue_context_release_complete*/ 0,
     0},       /* UEContextRelease */
    {0, 0, 0}, /* CellTrafficTrace */
               // UPDATE RELEASE 9
    {0, 0, 0}, /* Kill */
    {0, 0, 0}, /* DownlinkUEAssociatedLPPaTransport  */
    {0, 0, 0}, /* UplinkUEAssociatedLPPaTransport */
    {ngap_amf_handle_uplink_nas_transport, 0, 0}, /* uplinkNASTransport */
    {0, 0, 0},                                    /* downlinkNASTransport */
    {0, 0, 0}, /* DownlinkNonUEAssociatedLPPaTransport */
    {0, 0, 0}, /* UplinkNonUEAssociatedLPPaTransport */
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
int ngap_amf_set_cause(
    Ngap_Cause_t* cause_p, const Ngap_Cause_PR cause_type,
    const long cause_value) {
  DevAssert(cause_p != NULL);
  cause_p->present = cause_type;

  switch (cause_type) {
    case Ngap_Cause_PR_radioNetwork:
      cause_p->choice.misc = cause_value;
      break;

    case Ngap_Cause_PR_transport:
      cause_p->choice.transport = cause_value;
      break;

    case Ngap_Cause_PR_nas:
      cause_p->choice.nas = cause_value;
      break;

    case Ngap_Cause_PR_protocol:
      cause_p->choice.protocol = cause_value;
      break;

    case Ngap_Cause_PR_misc:
      cause_p->choice.misc = cause_value;
      break;

    default:
      return -1;
  }

  return 0;
}

//------------------------------------------------------------------------------
int ngap_amf_generate_ng_setup_failure(
    const sctp_assoc_id_t assoc_id, const Ngap_Cause_PR cause_type,
    const long cause_value, const long time_to_wait) {
  uint8_t* buffer_p = 0;
  uint32_t length   = 0;
  Ngap_NGAP_PDU_t pdu;
  Ngap_NGSetupFailure_t* out;
  Ngap_NGSetupFailureIEs_t* ie = NULL;
  int rc                       = RETURNok;

  OAILOG_FUNC_IN(LOG_NGAP);

  memset(&pdu, 0, sizeof(pdu));
  pdu.present = Ngap_NGAP_PDU_PR_unsuccessfulOutcome;
  pdu.choice.unsuccessfulOutcome.procedureCode = Ngap_ProcedureCode_id_NGSetup;
  pdu.choice.unsuccessfulOutcome.criticality   = Ngap_Criticality_reject;
  pdu.choice.unsuccessfulOutcome.value.present =
      Ngap_UnsuccessfulOutcome__value_PR_NGSetupFailure;
  out = &pdu.choice.unsuccessfulOutcome.value.choice.NGSetupFailure;

  ie = (Ngap_NGSetupFailureIEs_t*) calloc(1, sizeof(Ngap_NGSetupFailureIEs_t));
  ie->id            = Ngap_ProtocolIE_ID_id_Cause;
  ie->criticality   = Ngap_Criticality_ignore;
  ie->value.present = Ngap_NGSetupFailureIEs__value_PR_Cause;
  ngap_amf_set_cause(&ie->value.choice.Cause, cause_type, cause_value);
  ASN_SEQUENCE_ADD(&out->protocolIEs.list, ie);

  /*
   * Include the optional field time to wait only if the value is > -1
   */
  if (time_to_wait > -1) {
    ie =
        (Ngap_NGSetupFailureIEs_t*) calloc(1, sizeof(Ngap_NGSetupFailureIEs_t));
    ie->id                      = Ngap_ProtocolIE_ID_id_TimeToWait;
    ie->criticality             = Ngap_Criticality_ignore;
    ie->value.present           = Ngap_NGSetupFailureIEs__value_PR_TimeToWait;
    ie->value.choice.TimeToWait = time_to_wait;
    ASN_SEQUENCE_ADD(&out->protocolIEs.list, ie);
  }

  if (ngap_amf_encode_pdu(&pdu, &buffer_p, &length) < 0) {
    OAILOG_ERROR(LOG_NGAP, "Failed to encode ng setup failure\n");
    OAILOG_FUNC_RETURN(LOG_NGAP, RETURNerror);
  }

  bstring b = blk2bstr(buffer_p, length);
  free(buffer_p);
  rc = ngap_amf_itti_send_sctp_request(&b, assoc_id, 0, INVALID_AMF_UE_NGAP_ID);
  OAILOG_FUNC_RETURN(LOG_NGAP, rc);
}

////////////////////////////////////////////////////////////////////////////////
//************************** Management procedures ***************************//
////////////////////////////////////////////////////////////////////////////////

//------------------------------------------------------------------------------
int ngap_amf_handle_ng_setup_request(
    ngap_state_t* state, const sctp_assoc_id_t assoc_id,
    const sctp_stream_id_t stream, Ngap_NGAP_PDU_t* pdu) {
  int rc = RETURNok;

  Ngap_NGSetupRequest_t* container                = NULL;
  Ngap_NGSetupRequestIEs_t* ie                    = NULL;
  Ngap_NGSetupRequestIEs_t* ie_gnb_name           = NULL;
  Ngap_NGSetupRequestIEs_t* ie_supported_tas      = NULL;
  Ngap_NGSetupRequestIEs_t* ie_default_paging_drx = NULL;

  gnb_description_t* gnb_association = NULL;
  uint32_t gnb_id                    = 0;
  char* gnb_name                     = NULL;
  int ta_ret                         = 0;
  uint8_t bplmn_list_count           = 0;  // Broadcast PLMN list count

  OAILOG_FUNC_IN(LOG_NGAP);
  increment_counter("ng_setup", 1, NO_LABELS);
  if (!hss_associated) {
    /*
     * Can not process the request, AMF is not connected to HSS
     */
    OAILOG_ERROR(
        LOG_NGAP,
        "Rejecting ng setup request Can not process the request, AMF is not "
        "connected to HSS\n");
    rc = ngap_amf_generate_ng_setup_failure(
        assoc_id, Ngap_Cause_PR_misc, Ngap_CauseMisc_unspecified, -1);
    increment_counter(
        "ng_setup", 1, 2, "result", "failure", "cause", "s6a_interface_not_up");
    OAILOG_FUNC_RETURN(LOG_NGAP, rc);
  }

  DevAssert(pdu != NULL);
  container = &pdu->choice.initiatingMessage.value.choice.NGSetupRequest;
  /*
   * We received a new valid Ng Setup Request on a stream != 0.
   * This should not happen -> reject gNB ng setup request.
   */

  if (stream != 0) {
    OAILOG_ERROR(LOG_NGAP, "Received new ng setup request on stream != 0\n");
    /*
     * Send a ng setup failure with protocol cause unspecified
     */
    rc = ngap_amf_generate_ng_setup_failure(
        assoc_id, Ngap_Cause_PR_protocol, Ngap_CauseProtocol_unspecified, -1);
    increment_counter(
        "ng_setup", 1, 2, "result", "failure", "cause",
        "sctp_stream_id_non_zero");
    OAILOG_FUNC_RETURN(LOG_NGAP, rc);
  }

  /* Handling of ngsetup cases as follows.
   * If we don't know about the association, we haven't processed the new
   * association yet, so hope the gNB will retry the ng setup. Ignore and
   * return. If we get this message when the Ng interface of the AMF state is in
   * READY state then it is protocol error or out of sync state. Ignore it and
   * return. Assume AMF would detect SCTP association failure and would Ng
   * interface state to accept Ngsetup from gNB. If we get this message when the
   * ng interface of the AMF is in SHUTDOWN stage, we just hope the gNB will
   * retry and that will result in a new association getting established
   * followed by a subsequent ng setup, return Ngap_TimeToWait_v20s. If we get
   * this message when the ng interface of the AMF is in RESETTING stage then we
   * return Ngap_TimeToWait_v20s.
   */
  if ((gnb_association = ngap_state_get_gnb(state, assoc_id)) == NULL) {
    /*
     *
     * This should not happen as the thread processing new associations is the
     * one that reads data from the socket. Promote to an assert once we have
     * more confidence.
     */
    OAILOG_ERROR(LOG_NGAP, "Ignoring ng setup from unknown assoc %u", assoc_id);
    OAILOG_FUNC_RETURN(LOG_NGAP, RETURNok);
  }

  if (gnb_association->ng_state == NGAP_RESETING ||
      gnb_association->ng_state == NGAP_SHUTDOWN) {
    OAILOG_WARNING(
        LOG_NGAP, "Ignoring ngsetup from gNB in state %s on assoc id %u",
        ng_gnb_state2str(gnb_association->ng_state), assoc_id);
    rc = ngap_amf_generate_ng_setup_failure(
        assoc_id, Ngap_Cause_PR_transport,
        Ngap_CauseTransport_transport_resource_unavailable,
        Ngap_TimeToWait_v20s);
    increment_counter(
        "ng_setup", 1, 2, "result", "failure", "cause", "invalid_state");
    OAILOG_FUNC_RETURN(LOG_NGAP, rc);
  }
  log_queue_item_t* context = NULL;
  OAILOG_MESSAGE_START_SYNC(
      OAILOG_LEVEL_DEBUG, LOG_NGAP, (&context),
      "New ng setup request incoming from ");
  // shared_log_queue_item_t *context = NULL;
  // OAILOG_MESSAGE_START_ASYNC (OAILOG_LEVEL_DEBUG, LOG_NGAP, (&context), "New
  // ng setup request incoming from ");

  NGAP_FIND_PROTOCOLIE_BY_ID(
      Ngap_NGSetupRequestIEs_t, ie_gnb_name, container,
      Ngap_ProtocolIE_ID_id_RANNodeName, false);
  if (ie_gnb_name) {
    OAILOG_MESSAGE_ADD_SYNC(
        context, "%*s ", (int) ie_gnb_name->value.choice.RANNodeName.size,
        ie_gnb_name->value.choice.RANNodeName.buf);
    gnb_name = (char*) ie_gnb_name->value.choice.RANNodeName.buf;
  }

  NGAP_FIND_PROTOCOLIE_BY_ID(
      Ngap_NGSetupRequestIEs_t, ie, container,
      Ngap_ProtocolIE_ID_id_GlobalRANNodeID, true);
  if (ie->value.choice.GlobalRANNodeID.choice.globalGNB_ID.gNB_ID.present ==
      Ngap_GNB_ID_PR_gNB_ID) {
    // Home gNB ID = 28 bits
    uint8_t* gnb_id_buf = ie->value.choice.GlobalRANNodeID.choice.globalGNB_ID
                              .gNB_ID.choice.gNB_ID.buf;

    if (ie->value.choice.GlobalRANNodeID.choice.globalGNB_ID.gNB_ID.choice
            .gNB_ID.size != 28) {
      // TODO: handle case were size != 28 -> notify ? reject ?
    }

    gnb_id = (gnb_id_buf[0] << 20) + (gnb_id_buf[1] << 12) +
             (gnb_id_buf[2] << 4) + ((gnb_id_buf[3] & 0xf0) >> 4);
    OAILOG_MESSAGE_ADD_SYNC(context, "home gNB id: %07x", gnb_id);
  }

  OAILOG_MESSAGE_FINISH((void*) context);

  /* Requirement AMF36.413R10_8.7.3.4 Abnormal Conditions
   * If the gNB initiates the procedure by sending a Ng SETUP REQUEST message
   * including the PLMN Identity IEs and none of the PLMNs provided by the gNB
   * is identified by the AMF, then the AMF shall reject the gNB Ng Setup
   * Request procedure with the appropriate cause value, e.g, Unknown PLMN.
   */
  NGAP_FIND_PROTOCOLIE_BY_ID(
      Ngap_NGSetupRequestIEs_t, ie_supported_tas, container,
      Ngap_ProtocolIE_ID_id_SupportedTAList, true);

  ta_ret = ngap_amf_compare_ta_lists(
      &ie_supported_tas->value.choice.SupportedTAList);

  /*
   * gNB and AMF have no common PLMN
   */
  if (ta_ret != TA_LIST_RET_OK) {
    OAILOG_ERROR(
        LOG_NGAP, "No Common PLMN with gNB, generate_ng_setup_failure\n");
    rc = ngap_amf_generate_ng_setup_failure(
        assoc_id, Ngap_Cause_PR_misc, Ngap_CauseMisc_unknown_PLMN,
        Ngap_TimeToWait_v20s);

    increment_counter(
        "ng_setup", 1, 2, "result", "failure", "cause",
        "plmnid_or_tac_mismatch");
    OAILOG_FUNC_RETURN(LOG_NGAP, rc);
  }

  Ngap_SupportedTAList_t* ta_list =
      &ie_supported_tas->value.choice.SupportedTAList;
  m5g_supported_ta_list_t* supp_ta_list = &gnb_association->supported_ta_list;
  supp_ta_list->list_count              = ta_list->list.count;

  /* Storing supported TAI lists received in Ng SETUP REQUEST message */
  for (int tai_idx = 0; tai_idx < supp_ta_list->list_count; tai_idx++) {
    Ngap_SupportedTAItem_t* tai = NULL;
    tai                         = ta_list->list.array[tai_idx];
    tai->tAC.size               = 2;  // ACL_TAG temp to test remove later
    OCTET_STRING_TO_TAC(
        &tai->tAC, supp_ta_list->supported_tai_items[tai_idx].tac);

    bplmn_list_count = tai->broadcastPLMNList.list.count;
    if (bplmn_list_count > NGAP_MAX_BROADCAST_PLMNS) {
      OAILOG_ERROR(
          LOG_NGAP, "Maximum Broadcast PLMN list count exceeded, count = %d\n",
          bplmn_list_count);
    }
    supp_ta_list->supported_tai_items[tai_idx].bplmnlist_count =
        bplmn_list_count;
    for (int plmn_idx = 0; plmn_idx < bplmn_list_count; plmn_idx++) {
      TBCD_TO_PLMN_T(
          &tai->broadcastPLMNList.list.array[plmn_idx]->pLMNIdentity,
          &supp_ta_list->supported_tai_items[tai_idx].bplmns[plmn_idx]);
    }
  }
  OAILOG_DEBUG(LOG_NGAP, "Adding gNB to the list of served gNBs\n");

  gnb_association->gnb_id = gnb_id;

  NGAP_FIND_PROTOCOLIE_BY_ID(
      Ngap_NGSetupRequestIEs_t, ie_default_paging_drx, container,
      Ngap_ProtocolIE_ID_id_DefaultPagingDRX, true);

  gnb_association->default_paging_drx =
      ie_default_paging_drx->value.choice.PagingDRX;

  if (gnb_name != NULL) {
    memcpy(
        gnb_association->gnb_name, ie_gnb_name->value.choice.RANNodeName.buf,
        ie_gnb_name->value.choice.RANNodeName.size);
    gnb_association->gnb_name[ie_gnb_name->value.choice.RANNodeName.size] =
        '\0';
  }

  // ngap_dump_gnb(gnb_association);
  rc = ngap_generate_ng_setup_response(state, gnb_association);
  if (rc == RETURNok) {
    // update_amf_app_stats_connected_gnb_add();
    set_gauge("ng_connection", 1, 1, "gnb_name", gnb_association->gnb_name);
    increment_counter("ng_setup", 1, 1, "result", "success");
    // ng_setup_success_event(gnb_name, gnb_id);
  }
  OAILOG_FUNC_RETURN(LOG_NGAP, rc);
}

//------------------------------------------------------------------------------
int ngap_generate_ng_setup_response(
    ngap_state_t* state, gnb_description_t* gnb_association) {
  Ngap_NGAP_PDU_t pdu;
  Ngap_NGSetupResponse_t* out;
  Ngap_NGSetupResponseIEs_t* ie        = NULL;
  Ngap_ServedGUAMIItem_t* servedGUAMFI = NULL;
  int enc_rval                         = 0;
  uint8_t* buffer                      = NULL;
  uint32_t length                      = 0;
  int rc                               = RETURNok;

  DevAssert(gnb_association != NULL);
  memset(&pdu, 0, sizeof(pdu));

  pdu.present = Ngap_NGAP_PDU_PR_successfulOutcome;

  pdu.choice.successfulOutcome.procedureCode = Ngap_ProcedureCode_id_NGSetup;
  pdu.choice.successfulOutcome.criticality   = Ngap_Criticality_reject;
  pdu.choice.successfulOutcome.value.present =
      Ngap_SuccessfulOutcome__value_PR_NGSetupResponse;
  out = &pdu.choice.successfulOutcome.value.choice.NGSetupResponse;

  ie =
      (Ngap_NGSetupResponseIEs_t*) calloc(1, sizeof(Ngap_NGSetupResponseIEs_t));
  ie->id            = Ngap_ProtocolIE_ID_id_AMFName;
  ie->criticality   = Ngap_Criticality_reject;
  ie->value.present = Ngap_NGSetupResponseIEs__value_PR_AMFName;

  char* amf_name = "AMF1";

  OCTET_STRING_fromBuf(&ie->value.choice.AMFName, amf_name, strlen(amf_name));

  ASN_SEQUENCE_ADD(&out->protocolIEs.list, ie);

  // Generating response
  ie =
      (Ngap_NGSetupResponseIEs_t*) calloc(1, sizeof(Ngap_NGSetupResponseIEs_t));
  ie->id            = Ngap_ProtocolIE_ID_id_ServedGUAMIList;
  ie->criticality   = Ngap_Criticality_reject;
  ie->value.present = Ngap_NGSetupResponseIEs__value_PR_ServedGUAMIList;

  // memset for gcc 4.8.4 instead of {0}, servedGUAMFI.servedPLMNs
  servedGUAMFI = calloc(1, sizeof *servedGUAMFI);

#if 0
amf_config_read_lock(&amf_config);
  /*
   * Use the guamfi parameters provided by configuration
   * that should be sorted
   */
  for (i = 0; i < amf_config.served_tai.nb_tai; i++) {
    bool plmn_added = false;
    for (j = 0; j < i; j++) {
      if ((amf_config.served_tai.plmn_mcc[j] ==
           amf_config.served_tai.plmn_mcc[i]) &&
          (amf_config.served_tai.plmn_mnc[j] ==
           amf_config.served_tai.plmn_mnc[i]) &&
          (amf_config.served_tai.plmn_mnc_len[j] ==
           amf_config.served_tai.plmn_mnc_len[i])) {
        plmn_added = true;
        break;
      }
    }
    if (false == plmn_added) {
      Ngap_PLMNIdentity_t* plmn = NULL;
      plmn                      = &servedGUAMFI->gUAMI.pLMNIdentity;
      MCC_MNC_TO_PLMNID(
          amf_config.served_tai.plmn_mcc[i], amf_config.served_tai.plmn_mnc[i],
          amf_config.served_tai.plmn_mnc_len[i], plmn);
    }
  }

  for (i = 0; i < amf_config.guamfi.nb; i++) {
    Ngap_AMFRegionID_t* amf_gid = NULL;
    Ngap_AMFSetID_t* amfc       = NULL;

    amf_gid = &servedGUAMFI->gUAMI.aMFRegionID;
    INT16_TO_OCTET_STRING(amf_config.guamfi.guamfi[i].amf_gid, amf_gid);

    amfc = &servedGUAMFI->gUAMI.aMFSetID;
    INT8_TO_OCTET_STRING(amf_config.guamfi.guamfi[i].amf_code, amfc);
  }
#endif
  /*************************Temp code******************************/
  Ngap_PLMNIdentity_t* plmn = NULL;
  plmn                      = &servedGUAMFI->gUAMI.pLMNIdentity;

  OCTET_STRING_fromBuf(plmn, buf_plmn, sizeof(buf_plmn) /*3bytes*/);
  Ngap_AMFRegionID_t* amf_gid = NULL;
  Ngap_AMFSetID_t* amfc       = NULL;
  Ngap_AMFPointer_t* aMFP     = NULL;

  amf_gid = &servedGUAMFI->gUAMI.aMFRegionID;
  INT8_TO_OCTET_STRING(1, amf_gid);  // 8

  amfc = &servedGUAMFI->gUAMI.aMFSetID;
  UE_ID_INDEX_TO_BIT_STRING(1, amfc);  // 10

  aMFP = &servedGUAMFI->gUAMI.aMFPointer;
  AMF_POINTER_TO_BIT_STRING(1, aMFP);  // 6

  /*************************Temp code******************************/

  ASN_SEQUENCE_ADD(&ie->value.choice.ServedGUAMIList.list, servedGUAMFI);
  ASN_SEQUENCE_ADD(&out->protocolIEs.list, ie);

  ie =
      (Ngap_NGSetupResponseIEs_t*) calloc(1, sizeof(Ngap_NGSetupResponseIEs_t));
  ie->id            = Ngap_ProtocolIE_ID_id_RelativeAMFCapacity;
  ie->criticality   = Ngap_Criticality_ignore;
  ie->value.present = Ngap_NGSetupResponseIEs__value_PR_RelativeAMFCapacity;
  ie->value.choice.RelativeAMFCapacity = amf_config.relative_capacity;
  ASN_SEQUENCE_ADD(&out->protocolIEs.list, ie);

  //  amf_config_unlock(&amf_config);
  /*
   * The AMF is only serving E-UTRAN RAT, so the list contains only one element
   */

  /*PLMNList*/
  ie =
      (Ngap_NGSetupResponseIEs_t*) calloc(1, sizeof(Ngap_NGSetupResponseIEs_t));
  ie->id            = Ngap_ProtocolIE_ID_id_PLMNSupportList;
  ie->criticality   = Ngap_Criticality_reject;
  ie->value.present = Ngap_NGSetupResponseIEs__value_PR_PLMNSupportList;

  Ngap_PLMNSupportItem_t* plmnItem =
      (Ngap_PLMNSupportItem_t*) calloc(1, sizeof(Ngap_PLMNSupportItem_t));

  OCTET_STRING_fromBuf(
      &plmnItem->pLMNIdentity, buf_plmn, sizeof(buf_plmn) /*3bytes*/);

  /* Ngap_SliceSupportList */

  Ngap_SliceSupportItem_t* SliceItem =
      (Ngap_SliceSupportItem_t*) calloc(1, sizeof(Ngap_SliceSupportItem_t));

  char* from_buf = "0x11";

  OCTET_STRING_fromBuf(&SliceItem->s_NSSAI.sST, from_buf, 1);

  ASN_SEQUENCE_ADD(
      &plmnItem->sliceSupportList.list,
      SliceItem);  // adding slice item to slice list

  ASN_SEQUENCE_ADD(
      &ie->value.choice.PLMNSupportList.list,
      plmnItem);  // adding plmn item to plmn list

  ASN_SEQUENCE_ADD(&out->protocolIEs.list, ie);

  enc_rval = ngap_amf_encode_pdu(&pdu, &buffer, &length);

  /*
   * Failed to encode ng setup response...
   */
  if (enc_rval < 0) {
    OAILOG_DEBUG(LOG_NGAP, "Removed gNB %d\n", gnb_association->sctp_assoc_id);
    ngap_remove_gnb(state, gnb_association);
  } else {
    /*
     * Consider the response as sent. NGAP is ready to accept UE contexts
     */
    gnb_association->ng_state = NGAP_READY;
  }

  /*
   * Non-UE signalling -> stream 0
   */
  bstring b = blk2bstr(buffer, length);
  free(buffer);
  rc = ngap_amf_itti_send_sctp_request(
      &b, gnb_association->sctp_assoc_id, 0, INVALID_AMF_UE_NGAP_ID);

  OAILOG_FUNC_RETURN(LOG_NGAP, rc);
}

////////////////////////////////////////////////////////////////////////////////
//******************* Context Management procedures **************************//
////////////////////////////////////////////////////////////////////////////////

//------------------------------------------------------------------------------
int ngap_amf_handle_initial_context_setup_response(
    ngap_state_t* state, __attribute__((unused)) const sctp_assoc_id_t assoc_id,
    __attribute__((unused)) const sctp_stream_id_t stream,
    Ngap_NGAP_PDU_t* pdu) {
  Ngap_InitialContextSetupResponse_t* container;
  Ngap_InitialContextSetupResponseIEs_t* ie                          = NULL;
  Ngap_PDUSessionResourceSetupItemSURes_t* pduSessionSetupListCtxRes = NULL;
  m5g_ue_description_t* ue_ref_p                                     = NULL;
  MessageDef* message_p                                              = NULL;
  int rc                                                             = RETURNok;
  amf_ue_ngap_id_t amf_ue_ngap_id                                    = 0;
  gnb_ue_ngap_id_t gnb_ue_ngap_id                                    = 0;
  imsi64_t imsi64;

  OAILOG_FUNC_IN(LOG_NGAP);
  container =
      &pdu->choice.successfulOutcome.value.choice.InitialContextSetupResponse;
  NGAP_FIND_PROTOCOLIE_BY_ID(
      Ngap_InitialContextSetupResponseIEs_t, ie, container,
      Ngap_ProtocolIE_ID_id_AMF_UE_NGAP_ID, true);

  if (ie) {
    amf_ue_ngap_id = (uint32_t) ie->value.choice.AMF_UE_NGAP_ID;
    if ((ue_ref_p = ngap_state_get_ue_amfid((uint32_t) amf_ue_ngap_id)) ==
        NULL) {
      OAILOG_DEBUG(
          LOG_NGAP,
          "No UE is attached to this amf UE ngap id: " AMF_UE_NGAP_ID_FMT
          " %u(10)\n",
          (uint32_t) amf_ue_ngap_id, (uint32_t) amf_ue_ngap_id);
      OAILOG_FUNC_RETURN(LOG_NGAP, RETURNerror);
    }
  } else {
    OAILOG_FUNC_RETURN(LOG_NGAP, RETURNerror);
  }
  ngap_imsi_map_t* ngap_imsi_map = get_ngap_imsi_map();
  hashtable_uint64_ts_get(
      ngap_imsi_map->amf_ue_id_imsi_htbl, (const hash_key_t) amf_ue_ngap_id,
      &imsi64);

  NGAP_FIND_PROTOCOLIE_BY_ID(
      Ngap_InitialContextSetupResponseIEs_t, ie, container,
      Ngap_ProtocolIE_ID_id_RAN_UE_NGAP_ID, true);
  if (ie) {
    gnb_ue_ngap_id = (gnb_ue_ngap_id_t)(
        ie->value.choice.RAN_UE_NGAP_ID & GNB_UE_NGAP_ID_MASK);
  } else {
    OAILOG_FUNC_RETURN(LOG_NGAP, RETURNerror);
  }

  if (ue_ref_p->gnb_ue_ngap_id != gnb_ue_ngap_id) {
    OAILOG_DEBUG_UE(
        LOG_NGAP, imsi64,
        "Mismatch in gNB UE NGAP ID, known: " GNB_UE_NGAP_ID_FMT
        " %u(10), received: 0x%06x %u(10)\n",
        ue_ref_p->gnb_ue_ngap_id, ue_ref_p->gnb_ue_ngap_id,
        (uint32_t) gnb_ue_ngap_id, (uint32_t) gnb_ue_ngap_id);
    OAILOG_FUNC_RETURN(LOG_NGAP, RETURNerror);
  }

  NGAP_FIND_PROTOCOLIE_BY_ID(
      Ngap_InitialContextSetupResponseIEs_t, ie, container,
      Ngap_ProtocolIE_ID_id_PDUSessionResourceSetupListSURes, false);

  if (ie) {
    if (ie->value.choice.PDUSessionResourceSetupListCxtRes.list.count < 1) {
      OAILOG_WARNING_UE(LOG_NGAP, imsi64, "PDUSession creation has failed\n");
      OAILOG_FUNC_RETURN(LOG_NGAP, RETURNerror);
    }
  }
  ue_ref_p->ng_ue_state = NGAP_UE_CONNECTED;
  message_p =
      itti_alloc_new_message(TASK_NGAP, AMF_APP_INITIAL_CONTEXT_SETUP_RSP);
  AMF_APP_INITIAL_CONTEXT_SETUP_RSP(message_p).ue_id = ue_ref_p->amf_ue_ngap_id;

  if (ie) {
    AMF_APP_INITIAL_CONTEXT_SETUP_RSP(message_p)
        .pdusesssion_setup_list.no_of_items =
        ie->value.choice.PDUSessionResourceSetupListCxtRes.list.count;

    for (int item = 0;
         item < ie->value.choice.PDUSessionResourceSetupListCxtRes.list.count;
         item++) {
      /*
       * Bad, very bad cast...
       */
      // Info: need to update 38413 for this

      pduSessionSetupListCtxRes =
          (Ngap_PDUSessionResourceSetupItemSURes_t*) ie->value.choice
              .PDUSessionResourceSetupListCxtRes.list.array[item];
      AMF_APP_INITIAL_CONTEXT_SETUP_RSP(message_p)
          .pdusesssion_setup_list.item[item]
          .Pdu_Session_ID = pduSessionSetupListCtxRes->pDUSessionID;
    }
  } /*if(ie)*/

  NGAP_FIND_PROTOCOLIE_BY_ID(
      Ngap_InitialContextSetupResponseIEs_t, ie, container,
      Ngap_ProtocolIE_ID_id_PDUSessionResourceFailedToSetupListSURes, false);

  message_p->ittiMsgHeader.imsi = imsi64;
  rc = send_msg_to_task(&ngap_task_zmq_ctx, TASK_AMF_APP, message_p);
  OAILOG_FUNC_RETURN(LOG_NGAP, rc);
}

//------------------------------------------------------------------------------
int ngap_handle_new_association(
    ngap_state_t* state, sctp_new_peer_t* sctp_new_peer_p) {
  gnb_description_t* gnb_association = NULL;

  OAILOG_FUNC_IN(LOG_NGAP);
  DevAssert(sctp_new_peer_p != NULL);

  /*
   * Checking that the assoc id has a valid gNB attached to.
   */
  gnb_association = ngap_state_get_gnb(state, sctp_new_peer_p->assoc_id);
  if (gnb_association == NULL) {
    OAILOG_DEBUG(
        LOG_NGAP, "Create gNB context for assoc_id: %d\n",
        sctp_new_peer_p->assoc_id);
    /*
     * Create new context
     */
    gnb_association = ngap_new_gnb(state);

    if (gnb_association == NULL) {
      /*
       * We failed to allocate memory
       */
      /*
       * TODO: send reject there
       */
      OAILOG_ERROR(
          LOG_NGAP, "Failed to allocate gNB context for assoc_id: %d\n",
          sctp_new_peer_p->assoc_id);
      OAILOG_FUNC_RETURN(LOG_NGAP, RETURNok);
    }
    gnb_association->sctp_assoc_id = sctp_new_peer_p->assoc_id;
    hashtable_rc_t hash_rc         = hashtable_ts_insert(
        &state->gnbs, (const hash_key_t) gnb_association->sctp_assoc_id,
        (void*) gnb_association);
    if (HASH_TABLE_OK != hash_rc) {
      OAILOG_FUNC_RETURN(LOG_NGAP, RETURNerror);
    }
  } else if (
      (gnb_association->ng_state == NGAP_SHUTDOWN) ||
      (gnb_association->ng_state == NGAP_RESETING)) {
    OAILOG_WARNING(
        LOG_NGAP,
        "Received new association request on an association that is being %s, "
        "ignoring",
        ng_gnb_state2str(gnb_association->ng_state));
    OAILOG_FUNC_RETURN(LOG_NGAP, RETURNerror);
  } else {
    OAILOG_DEBUG(
        LOG_NGAP, "gNB context already exists for assoc_id: %d, update it\n",
        sctp_new_peer_p->assoc_id);
  }

  gnb_association->sctp_assoc_id = sctp_new_peer_p->assoc_id;
  /*
   * Fill in in and out number of streams available on SCTP connection.
   */
  gnb_association->instreams  = (sctp_stream_id_t) sctp_new_peer_p->instreams;
  gnb_association->outstreams = (sctp_stream_id_t) sctp_new_peer_p->outstreams;
  /*
   * initialize the next sctp stream to 1 as 0 is reserved for non
   * * * * ue associated signalling.
   */
  gnb_association->next_sctp_stream = 1;
  gnb_association->ng_state         = NGAP_INIT;
  OAILOG_FUNC_RETURN(LOG_NGAP, RETURNok);
}

int ngap_amf_handle_ue_context_release_request(
    __attribute__((unused)) ngap_state_t* state,
    __attribute__((unused)) const sctp_assoc_id_t assoc_id,
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

  /* Fix - AMF shall handle UE Context Release received from the gNB
  irrespective of the cause. And AMF should release the Ng-U bearers for the UE
  and move the UE to ECM idle mode. Cause can influence whether to preserve GBR
  bearers or not.Since, as of now EPC doesn't support dedicated bearers, it is
  don't care scenario till we add support for dedicated bearers.
  */

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

      message_p =
          itti_alloc_new_message(TASK_NGAP, NGAP_UE_CONTEXT_RELEASE_REQ);

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

//------------------------------------------------------------------------------
int ngap_amf_generate_ue_context_release_command(
    ngap_state_t* state, m5g_ue_description_t* ue_ref_p, enum Ngcause cause,
    imsi64_t imsi64) {
  uint8_t* buffer = NULL;
  uint32_t length = 0;
  Ngap_NGAP_PDU_t pdu;
  Ngap_UEContextReleaseCommand_t* out;
  Ngap_UEContextReleaseCommand_IEs_t* ie = NULL;
  int rc                                 = RETURNok;
  Ngap_Cause_PR cause_type;
  long cause_value;

  OAILOG_FUNC_IN(LOG_NGAP);
  if (ue_ref_p == NULL) {
    OAILOG_FUNC_RETURN(LOG_NGAP, RETURNerror);
  }
  memset(&pdu, 0, sizeof(pdu));
  pdu.present = Ngap_NGAP_PDU_PR_initiatingMessage;
  pdu.choice.initiatingMessage.procedureCode =
      Ngap_ProcedureCode_id_UEContextRelease;
  pdu.choice.initiatingMessage.criticality = Ngap_Criticality_reject;
  pdu.choice.initiatingMessage.value.present =
      Ngap_InitiatingMessage__value_PR_UEContextReleaseCommand;
  out = &pdu.choice.initiatingMessage.value.choice.UEContextReleaseCommand;
  /*
   * Fill in ID pair
   */
  ie = (Ngap_UEContextReleaseCommand_IEs_t*) calloc(
      1, sizeof(Ngap_UEContextReleaseCommand_IEs_t));
  ie->id            = Ngap_ProtocolIE_ID_id_UE_NGAP_IDs;
  ie->criticality   = Ngap_Criticality_reject;
  ie->value.present = Ngap_UEContextReleaseCommand_IEs__value_PR_UE_NGAP_IDs;
  ie->value.choice.UE_NGAP_IDs.present = Ngap_UE_NGAP_IDs_PR_uE_NGAP_ID_pair;
  ie->value.choice.UE_NGAP_IDs.choice.uE_NGAP_ID_pair.aMF_UE_NGAP_ID =
      ue_ref_p->amf_ue_ngap_id;
  ie->value.choice.UE_NGAP_IDs.choice.uE_NGAP_ID_pair.rAN_UE_NGAP_ID =
      ue_ref_p->gnb_ue_ngap_id;
  ie->value.choice.UE_NGAP_IDs.choice.uE_NGAP_ID_pair.iE_Extensions = NULL;
  ASN_SEQUENCE_ADD(&out->protocolIEs.list, ie);

  ie = (Ngap_UEContextReleaseCommand_IEs_t*) calloc(
      1, sizeof(Ngap_UEContextReleaseCommand_IEs_t));
  ie->id            = Ngap_ProtocolIE_ID_id_Cause;
  ie->criticality   = Ngap_Criticality_ignore;
  ie->value.present = Ngap_UEContextReleaseCommand_IEs__value_PR_Cause;
  switch (cause) {
    case NGAP_NAS_DEREGISTER:
      cause_type  = Ngap_Cause_PR_nas;
      cause_value = Ngap_CauseNas_deregister;
      break;
    case NGAP_NAS_NORMAL_RELEASE:
      cause_type  = Ngap_Cause_PR_nas;
      cause_value = Ngap_CauseNas_unspecified;
      break;
    case NGAP_RADIO_NR_GENERATED_REASON:
      cause_type  = Ngap_Cause_PR_radioNetwork;
      cause_value = Ngap_CauseRadioNetwork_release_due_to_5gc_generated_reason;
      break;
    case NGAP_INITIAL_CONTEXT_SETUP_FAILED:
      cause_type  = Ngap_Cause_PR_radioNetwork;
      cause_value = Ngap_CauseRadioNetwork_unspecified;
      break;
    default:
      OAILOG_ERROR_UE(LOG_NGAP, imsi64, "Unknown cause for context release");
      OAILOG_FUNC_RETURN(LOG_NGAP, RETURNerror);
  }
  ngap_amf_set_cause(&ie->value.choice.Cause, cause_type, cause_value);
  ASN_SEQUENCE_ADD(&out->protocolIEs.list, ie);

  if (ngap_amf_encode_pdu(&pdu, &buffer, &length) < 0) {
    OAILOG_FUNC_RETURN(LOG_NGAP, RETURNerror);
  }

  bstring b = blk2bstr(buffer, length);
  free(buffer);
  rc = ngap_amf_itti_send_sctp_request(
      &b, ue_ref_p->sctp_assoc_id, ue_ref_p->sctp_stream_send,
      ue_ref_p->amf_ue_ngap_id);
  ue_ref_p->ng_ue_state = NGAP_UE_WAITING_CRR;

  // Start timer to track UE context release complete from gNB

  // We can safely remove UE context now, no need for timer
  ngap_amf_release_ue_context(state, ue_ref_p, imsi64);

  OAILOG_FUNC_RETURN(LOG_NGAP, rc);
}

//------------------------------------------------------------------------------
int ngap_handle_ue_context_release_command(
    ngap_state_t* state,
    const itti_ngap_ue_context_release_command_t* const
        ue_context_release_command_pP,
    imsi64_t imsi64) {
  m5g_ue_description_t* ue_ref_p = NULL;
  int rc                         = RETURNok;

  OAILOG_FUNC_IN(LOG_NGAP);
  if ((ue_ref_p = ngap_state_get_ue_amfid(
           ue_context_release_command_pP->amf_ue_ngap_id)) == NULL) {
    OAILOG_DEBUG_UE(
        LOG_NGAP, imsi64,
        "Ignoring UE with amf_ue_ngap_id " AMF_UE_NGAP_ID_FMT " %u(10)\n",
        ue_context_release_command_pP->amf_ue_ngap_id,
        ue_context_release_command_pP->amf_ue_ngap_id);
    rc = RETURNok;
  } else {
    /*
     * Check the cause. If it is implicit detach or sctp reset/shutdown no need
     * to send UE context release command to gNB. Free UE context locally.
     */

    if (ue_context_release_command_pP->cause == NGAP_IMPLICIT_CONTEXT_RELEASE ||
        ue_context_release_command_pP->cause == NGAP_SCTP_SHUTDOWN_OR_RESET ||
        ue_context_release_command_pP->cause ==
            NGAP_INITIAL_CONTEXT_SETUP_TMR_EXPRD ||
        ue_context_release_command_pP->cause == NGAP_INVALID_GNB_ID) {
      // ngap_remove_ue(state, ue_ref_p);
    } else {
      rc = ngap_amf_generate_ue_context_release_command(
          state, ue_ref_p, ue_context_release_command_pP->cause, imsi64);
    }
  }

  OAILOG_FUNC_RETURN(LOG_NGAP, rc);
}

//------------------------------------------------------------------------------
int ngap_amf_handle_ue_context_release_complete(
    ngap_state_t* state, __attribute__((unused)) const sctp_assoc_id_t assoc_id,
    __attribute__((unused)) const sctp_stream_id_t stream,
    Ngap_NGAP_PDU_t* pdu) {
  Ngap_UEContextReleaseComplete_t* container;
  Ngap_UEContextReleaseComplete_IEs_t* ie = NULL;
  m5g_ue_description_t* ue_ref_p          = NULL;
  amf_ue_ngap_id_t amf_ue_ngap_id         = 0;

  OAILOG_FUNC_IN(LOG_NGAP);
  container =
      &pdu->choice.successfulOutcome.value.choice.UEContextReleaseComplete;

  NGAP_FIND_PROTOCOLIE_BY_ID(
      Ngap_UEContextReleaseComplete_IEs_t, ie, container,
      Ngap_ProtocolIE_ID_id_AMF_UE_NGAP_ID, true);

  if (ie) {
    amf_ue_ngap_id = ie->value.choice.AMF_UE_NGAP_ID;
  } else {
    OAILOG_FUNC_RETURN(LOG_NGAP, RETURNok);
  }

  if ((ue_ref_p = ngap_state_get_ue_amfid(amf_ue_ngap_id)) == NULL) {
    /*
     * The UE context has already been deleted when the UE context release
     * command was sent
     * Ignore this message.
     */
    OAILOG_DEBUG(
        LOG_NGAP,
        " UE Context Release commplete: Ng context cleared. Ignore message for "
        "ueid " AMF_UE_NGAP_ID_FMT "\n",
        (uint32_t) amf_ue_ngap_id);
    OAILOG_FUNC_RETURN(LOG_NGAP, RETURNok);
  } else {
    /* This is an error scenario, the Ng UE context should have been deleted
     * when UE context release command was sent
     */
    OAILOG_ERROR(
        LOG_NGAP,
        " UE Context Release commplete: Ng context should have been cleared "
        "for "
        "ueid " AMF_UE_NGAP_ID_FMT "\n",
        (uint32_t) amf_ue_ngap_id);
    OAILOG_FUNC_RETURN(LOG_NGAP, RETURNerror);
  }
}

//------------------------------------------------------------------------------
int ngap_amf_handle_initial_context_setup_failure(
    ngap_state_t* state, __attribute__((unused)) const sctp_assoc_id_t assoc_id,
    __attribute__((unused)) const sctp_stream_id_t stream,
    Ngap_NGAP_PDU_t* pdu) {
  Ngap_InitialContextSetupFailure_t* container;
  Ngap_InitialContextSetupFailureIEs_t* ie = NULL;
  m5g_ue_description_t* ue_ref_p           = NULL;
  MessageDef* message_p                    = NULL;
  Ngap_Cause_PR cause_type;
  long cause_value;
  int rc                          = RETURNok;
  imsi64_t imsi64                 = INVALID_IMSI64;
  amf_ue_ngap_id_t amf_ue_ngap_id = 0;
  gnb_ue_ngap_id_t gnb_ue_ngap_id = 0;

  OAILOG_FUNC_IN(LOG_NGAP);
  container =
      &pdu->choice.unsuccessfulOutcome.value.choice.InitialContextSetupFailure;

  NGAP_FIND_PROTOCOLIE_BY_ID(
      Ngap_InitialContextSetupFailureIEs_t, ie, container,
      Ngap_ProtocolIE_ID_id_AMF_UE_NGAP_ID, true);
  if (ie) {
    amf_ue_ngap_id = ie->value.choice.AMF_UE_NGAP_ID;
  } else {
    OAILOG_FUNC_RETURN(LOG_NGAP, RETURNok);
  }

  NGAP_FIND_PROTOCOLIE_BY_ID(
      Ngap_InitialContextSetupFailureIEs_t, ie, container,
      Ngap_ProtocolIE_ID_id_RAN_UE_NGAP_ID, true);
  if (ie) {
    gnb_ue_ngap_id = (gnb_ue_ngap_id_t)(
        ie->value.choice.RAN_UE_NGAP_ID & GNB_UE_NGAP_ID_MASK);
  } else {
    OAILOG_FUNC_RETURN(LOG_NGAP, RETURNok);
  }

  if ((ue_ref_p = ngap_state_get_ue_amfid(amf_ue_ngap_id)) == NULL) {
    /*
     * AMF doesn't know the AMF UE NGAP ID provided.
     */
    OAILOG_INFO(
        LOG_NGAP,
        "INITIAL_CONTEXT_SETUP_FAILURE ignored. No context with "
        "amf_ue_ngap_id " AMF_UE_NGAP_ID_FMT
        " gnb_ue_ngap_id " GNB_UE_NGAP_ID_FMT " ",
        (uint32_t) amf_ue_ngap_id, (uint32_t) gnb_ue_ngap_id);
    OAILOG_FUNC_RETURN(LOG_NGAP, RETURNerror);
  }

  if (ue_ref_p->gnb_ue_ngap_id != gnb_ue_ngap_id) {
    // abnormal case. No need to do anything. Ignore the message
    OAILOG_DEBUG(
        LOG_NGAP,
        "INITIAL_CONTEXT_SETUP_FAILURE ignored, mismatch gnb_ue_ngap_id: "
        "ctxt " GNB_UE_NGAP_ID_FMT " != received " GNB_UE_NGAP_ID_FMT " ",
        (uint32_t) ue_ref_p->gnb_ue_ngap_id, (uint32_t) gnb_ue_ngap_id);
    OAILOG_FUNC_RETURN(LOG_NGAP, RETURNerror);
  }

  ngap_imsi_map_t* imsi_map = get_ngap_imsi_map();
  hashtable_uint64_ts_get(
      imsi_map->amf_ue_id_imsi_htbl, (const hash_key_t) amf_ue_ngap_id,
      &imsi64);

  // Pass this message to AMF APP for necessary handling
  // Log the Cause Type and Cause value
  NGAP_FIND_PROTOCOLIE_BY_ID(
      Ngap_InitialContextSetupFailureIEs_t, ie, container,
      Ngap_ProtocolIE_ID_id_Cause, true);
  if (ie) {
    cause_type = ie->value.choice.Cause.present;
  } else {
    OAILOG_FUNC_RETURN(LOG_NGAP, RETURNerror);
  }

  switch (cause_type) {
    case Ngap_Cause_PR_radioNetwork:
      cause_value = ie->value.choice.Cause.choice.radioNetwork;
      OAILOG_DEBUG_UE(
          LOG_NGAP, imsi64,
          "INITIAL_CONTEXT_SETUP_FAILURE with Cause_Type = Radio Network and "
          "Cause_Value = %ld\n",
          cause_value);
      break;

    case Ngap_Cause_PR_transport:
      cause_value = ie->value.choice.Cause.choice.transport;
      OAILOG_DEBUG_UE(
          LOG_NGAP, imsi64,
          "INITIAL_CONTEXT_SETUP_FAILURE with Cause_Type = Transport and "
          "Cause_Value = %ld\n",
          cause_value);
      break;

    case Ngap_Cause_PR_nas:
      cause_value = ie->value.choice.Cause.choice.nas;
      OAILOG_DEBUG_UE(
          LOG_NGAP, imsi64,
          "INITIAL_CONTEXT_SETUP_FAILURE with Cause_Type = NAS and Cause_Value "
          "= "
          "%ld\n",
          cause_value);
      break;

    case Ngap_Cause_PR_protocol:
      cause_value = ie->value.choice.Cause.choice.protocol;
      OAILOG_DEBUG_UE(
          LOG_NGAP, imsi64,
          "INITIAL_CONTEXT_SETUP_FAILURE with Cause_Type = Protocol and "
          "Cause_Value = %ld\n",
          cause_value);
      break;

    case Ngap_Cause_PR_misc:
      cause_value = ie->value.choice.Cause.choice.misc;
      OAILOG_DEBUG_UE(
          LOG_NGAP, imsi64,
          "INITIAL_CONTEXT_SETUP_FAILURE with Cause_Type = MISC and "
          "Cause_Value "
          "= %ld\n",
          cause_value);
      break;

    default:
      OAILOG_DEBUG_UE(
          LOG_NGAP, imsi64,
          "INITIAL_CONTEXT_SETUP_FAILURE with Invalid Cause_Type = %d\n",
          cause_type);
      OAILOG_FUNC_RETURN(LOG_NGAP, RETURNerror);
  }
  message_p =
      itti_alloc_new_message(TASK_NGAP, AMF_APP_INITIAL_CONTEXT_SETUP_FAILURE);
  memset(
      (void*) &message_p->ittiMsg.amf_app_initial_context_setup_failure, 0,
      sizeof(itti_amf_app_initial_context_setup_failure_t));
  AMF_APP_INITIAL_CONTEXT_SETUP_FAILURE(message_p).amf_ue_ngap_id =
      ue_ref_p->amf_ue_ngap_id;

  message_p->ittiMsgHeader.imsi = imsi64;
  rc = send_msg_to_task(&ngap_task_zmq_ctx, TASK_AMF_APP, message_p);
  OAILOG_FUNC_RETURN(LOG_NGAP, rc);
}

//------------------------------------------------------------------------------
const char* ng_gnb_state2str(enum amf_ng_gnb_state_s state) {
  switch (state) {
    case NGAP_INIT:
      return "NGAP_INIT";
    case NGAP_RESETING:
      return "NGAP_RESETING";
    case NGAP_READY:
      return "NGAP_READY";
    case NGAP_SHUTDOWN:
      return "NGAP_SHUTDOWN";
    default:
      return "unknown ngap_gnb_state";
  }
}

//------------------------------------------------------------------------------
typedef struct arg_ngap_send_gnb_dereg_ind_s {
  uint8_t current_ue_index;
  uint32_t handled_ues;
  MessageDef* message_p;
  uint32_t associated_gnb_id;
  uint32_t deregister_ue_count;
} arg_ngap_send_gnb_dereg_ind_t;

//------------------------------------------------------------------------------
bool ngap_send_gnb_deregistered_ind(
    __attribute__((unused)) const hash_key_t keyP, uint64_t const dataP,
    void* argP, void** resultP) {
  arg_ngap_send_gnb_dereg_ind_t* arg = (arg_ngap_send_gnb_dereg_ind_t*) argP;
  m5g_ue_description_t* ue_ref_p     = (m5g_ue_description_t*) dataP;
  imsi64_t imsi64                    = INVALID_IMSI64;

  // Ask for the release of each UE context associated to the gNB
  hash_table_ts_t* ngap_ue_state = get_ngap_ue_state();
  hashtable_ts_get(ngap_ue_state, (const hash_key_t) dataP, (void**) &ue_ref_p);
  if (ue_ref_p) {
    if (arg->current_ue_index == 0) {
      arg->message_p =
          itti_alloc_new_message(TASK_NGAP, NGAP_GNB_DEREGISTERED_IND);
    }
    if (ue_ref_p->amf_ue_ngap_id == INVALID_AMF_UE_NGAP_ID) {
      /*
       * Send deregistered ind for this also and let AMFAPP find the context
       * using gnb_ue_ngap_id_key
       */
      OAILOG_WARNING(LOG_NGAP, "UE with invalid AMF ngap id found");
    }
    ngap_imsi_map_t* imsi_map = get_ngap_imsi_map();
    hashtable_uint64_ts_get(
        imsi_map->amf_ue_id_imsi_htbl,
        (const hash_key_t) ue_ref_p->amf_ue_ngap_id, &imsi64);

    AssertFatal(
        arg->current_ue_index < NGAP_ITTI_UE_PER_DEREGISTER_MESSAGE,
        "Too many deregistered UEs reported in NGAP_GNB_DEREGISTERED_IND "
        "message ");
    NGAP_GNB_DEREGISTERED_IND(arg->message_p)
        .amf_ue_ngap_id[arg->current_ue_index] = ue_ref_p->amf_ue_ngap_id;
    NGAP_GNB_DEREGISTERED_IND(arg->message_p)
        .gnb_ue_ngap_id[arg->current_ue_index] = ue_ref_p->gnb_ue_ngap_id;

    arg->handled_ues++;
    arg->current_ue_index++;

    if (arg->handled_ues == arg->deregister_ue_count ||
        arg->current_ue_index == NGAP_ITTI_UE_PER_DEREGISTER_MESSAGE) {
      // Sending INVALID_IMSI64 because message is not specific to any UE/IMSI
      arg->message_p->ittiMsgHeader.imsi               = INVALID_IMSI64;
      NGAP_GNB_DEREGISTERED_IND(arg->message_p).gnb_id = arg->associated_gnb_id;
      NGAP_GNB_DEREGISTERED_IND(arg->message_p).nb_ue_to_deregister =
          (uint8_t) arg->current_ue_index;

      // Max UEs reached for this ITTI message, send message to AMF App
      OAILOG_DEBUG(
          LOG_NGAP,
          "Reached maximum UE count for this ITTI message. Sending "
          "deregistered indication to AMF App for UE count = %u\n",
          NGAP_GNB_DEREGISTERED_IND(arg->message_p).nb_ue_to_deregister);

      if (arg->current_ue_index == NGAP_ITTI_UE_PER_DEREGISTER_MESSAGE) {
        arg->current_ue_index = 0;
      }
      send_msg_to_task(&ngap_task_zmq_ctx, TASK_AMF_APP, arg->message_p);
      arg->message_p = NULL;
    }

    *resultP = arg->message_p;
  } else {
    OAILOG_TRACE(LOG_NGAP, "No valid UE provided in callback: %p\n", ue_ref_p);
  }
  return false;
}

typedef struct arg_ngap_construct_gnb_reset_req_s {
  uint8_t current_ue_index;
  MessageDef* msg;
} arg_ngap_construct_gnb_reset_req_t;

//------------------------------------------------------------------------------
int ngap_handle_sctp_disconnection(
    ngap_state_t* state, const sctp_assoc_id_t assoc_id, bool reset) {
  arg_ngap_send_gnb_dereg_ind_t arg  = {0};
  int i                              = 0;
  MessageDef* message_p              = NULL;
  gnb_description_t* gnb_association = NULL;

  OAILOG_FUNC_IN(LOG_NGAP);

  // Checking if the assoc id has a valid gNB attached to it
  gnb_association = ngap_state_get_gnb(state, assoc_id);
  if (gnb_association == NULL) {
    OAILOG_ERROR(LOG_NGAP, "No gNB attached to this assoc_id: %d\n", assoc_id);
    OAILOG_FUNC_RETURN(LOG_NGAP, RETURNerror);
  }

  OAILOG_INFO(
      LOG_NGAP,
      "SCTP disconnection request for association id %u, Reset Flag = "
      "%u. Connected UEs = %u \n",
      assoc_id, reset, gnb_association->nb_ue_associated);

  // First check if we can just reset the gNB state if there are no UEs
  if (!gnb_association->nb_ue_associated) {
    if (reset) {
      OAILOG_INFO(
          LOG_NGAP,
          "SCTP reset request for association id %u. No Connected UEs. "
          "Reset Flag = %u\n",
          assoc_id, reset);

      OAILOG_INFO(
          LOG_NGAP, "Moving gNB with assoc_id %u to INIT state\n", assoc_id);
      gnb_association->ng_state = NGAP_INIT;
      // update_amf_app_stats_connected_gnb_sub(); TODO : part of stats
    } else {
      OAILOG_INFO(
          LOG_NGAP,
          "SCTP Shutdown request for association id %u. No Connected UEs. "
          "Reset Flag = %u\n",
          assoc_id, reset);

      OAILOG_INFO(LOG_NGAP, "Removing gNB with association id %u \n", assoc_id);
      ngap_remove_gnb(state, gnb_association);
      // update_amf_app_stats_connected_gnb_sub(); TODO : part of stats
    }
    OAILOG_FUNC_RETURN(LOG_NGAP, RETURNok);
  }

  /*
   * Send Ngap deregister indication to AMF app in batches of UEs where
   * UE count in each batch <= NGAP_ITTI_UE_PER_DEREGISTER_MESSAGE
   */
  arg.associated_gnb_id   = gnb_association->gnb_id;
  arg.deregister_ue_count = gnb_association->ue_id_coll.num_elements;
  hashtable_uint64_ts_apply_callback_on_elements(
      &gnb_association->ue_id_coll, ngap_send_gnb_deregistered_ind,
      (void*) &arg, (void**) &message_p);

  for (i = arg.current_ue_index; i < NGAP_ITTI_UE_PER_DEREGISTER_MESSAGE; i++) {
    NGAP_GNB_DEREGISTERED_IND(message_p).amf_ue_ngap_id[arg.current_ue_index] =
        0;
    NGAP_GNB_DEREGISTERED_IND(message_p).gnb_ue_ngap_id[arg.current_ue_index] =
        0;
  }
  NGAP_GNB_DEREGISTERED_IND(message_p).gnb_id = gnb_association->gnb_id;
  message_p                                   = NULL;

  OAILOG_FUNC_RETURN(LOG_NGAP, RETURNok);
}

//------------------------------------------------------------------------------
void ngap_amf_handle_ue_context_rel_comp_timer_expiry(
    ngap_state_t* state, m5g_ue_description_t* ue_ref_p) {
  MessageDef* message_p = NULL;
  OAILOG_FUNC_IN(LOG_NGAP);
  DevAssert(ue_ref_p != NULL);
  ue_ref_p->ngap_ue_context_rel_timer.id = NGAP_TIMER_INACTIVE_ID;
  imsi64_t imsi64                        = INVALID_IMSI64;

  ngap_imsi_map_t* imsi_map = get_ngap_imsi_map();
  hashtable_uint64_ts_get(
      imsi_map->amf_ue_id_imsi_htbl,
      (const hash_key_t) ue_ref_p->amf_ue_ngap_id, &imsi64);

  OAILOG_DEBUG_UE(
      LOG_NGAP, imsi64, "Expired- UE Context Release Timer for UE id  %d \n",
      ue_ref_p->amf_ue_ngap_id);
  /*
   * Remove UE context and inform AMF_APP.
   */
  message_p =
      itti_alloc_new_message(TASK_NGAP, NGAP_UE_CONTEXT_RELEASE_COMPLETE);
  memset(
      (void*) &message_p->ittiMsg.ngap_ue_context_release_complete, 0,
      sizeof(itti_ngap_ue_context_release_complete_t));
  NGAP_UE_CONTEXT_RELEASE_COMPLETE(message_p).amf_ue_ngap_id =
      ue_ref_p->amf_ue_ngap_id;

  message_p->ittiMsgHeader.imsi = imsi64;
  send_msg_to_task(&ngap_task_zmq_ctx, TASK_AMF_APP, message_p);
  DevAssert(ue_ref_p->ng_ue_state == NGAP_UE_WAITING_CRR);

  OAILOG_DEBUG_UE(
      LOG_NGAP, imsi64, "Removed NGAP UE " AMF_UE_NGAP_ID_FMT "\n",
      (uint32_t) ue_ref_p->amf_ue_ngap_id);
  // ngap_remove_ue(state, ue_ref_p);

  hashtable_uint64_ts_remove(
      imsi_map->amf_ue_id_imsi_htbl,
      (const hash_key_t) ue_ref_p->amf_ue_ngap_id);

  OAILOG_FUNC_OUT(LOG_NGAP);
}

//------------------------------------------------------------------------------
void ngap_amf_release_ue_context(
    ngap_state_t* state, m5g_ue_description_t* ue_ref_p, imsi64_t imsi64) {
  MessageDef* message_p = NULL;
  OAILOG_FUNC_IN(LOG_NGAP);
  DevAssert(ue_ref_p != NULL);
  OAILOG_DEBUG_UE(
      LOG_NGAP, imsi64, "Releasing UE Context for UE id  %d \n",
      ue_ref_p->amf_ue_ngap_id);
  /*
   * Remove UE context and inform AMF_APP.
   */
  message_p =
      itti_alloc_new_message(TASK_NGAP, NGAP_UE_CONTEXT_RELEASE_COMPLETE);
  memset(
      (void*) &message_p->ittiMsg.ngap_ue_context_release_complete, 0,
      sizeof(itti_ngap_ue_context_release_complete_t));
  NGAP_UE_CONTEXT_RELEASE_COMPLETE(message_p).amf_ue_ngap_id =
      ue_ref_p->amf_ue_ngap_id;

  message_p->ittiMsgHeader.imsi = imsi64;
  send_msg_to_task(&ngap_task_zmq_ctx, TASK_AMF_APP, message_p);
  DevAssert(ue_ref_p->ng_ue_state == NGAP_UE_WAITING_CRR);
  OAILOG_DEBUG_UE(
      LOG_NGAP, imsi64, "Removed NGAP UE " AMF_UE_NGAP_ID_FMT "\n",
      (uint32_t) ue_ref_p->amf_ue_ngap_id);

  // ngap_remove_ue(state, ue_ref_p);
  OAILOG_FUNC_OUT(LOG_NGAP);
}

//------------------------------------------------------------------------------
int ngap_amf_handle_error_ind_message(
    ngap_state_t* state, const sctp_assoc_id_t assoc_id,
    const sctp_stream_id_t stream, Ngap_NGAP_PDU_t* message) {
  OAILOG_FUNC_IN(LOG_NGAP);
  OAILOG_WARNING(
      LOG_NGAP, "ERROR IND RCVD on Stream id %d, ignoring it\n", stream);
  increment_counter("ngap_error_ind_rcvd", 1, NO_LABELS);
  OAILOG_FUNC_RETURN(LOG_NGAP, RETURNok);
}

//------------------------------------------------------------------------------
const char* ngap_direction2str(uint8_t dir) {
  switch (dir) {
    case Ngap_NGAP_PDU_PR_NOTHING:
      return "<nothing>";
    case Ngap_NGAP_PDU_PR_initiatingMessage:
      return "originating message";
    case Ngap_NGAP_PDU_PR_successfulOutcome:
      return "successful outcome";
    case Ngap_NGAP_PDU_PR_unsuccessfulOutcome:
      return "unsuccessful outcome";
    default:
      return "unknown direction";
  }
}

//------------------------------------------------------------------------------
int ngap_amf_handle_pduSession_release_response(
    ngap_state_t* state, const sctp_assoc_id_t assoc_id,
    const sctp_stream_id_t stream, Ngap_NGAP_PDU_t* pdu) {
  OAILOG_FUNC_IN(LOG_NGAP);
  Ngap_PDUSessionResourceReleaseResponseIEs_t* ie     = NULL;
  Ngap_PDUSessionResourceReleaseResponse_t* container = NULL;
  m5g_ue_description_t* ue_ref_p                      = NULL;
  MessageDef* message_p                               = NULL;
  int rc                                              = RETURNok;
  imsi64_t imsi64                                     = INVALID_IMSI64;
  gnb_ue_ngap_id_t gnb_ue_ngap_id                     = 0;
  amf_ue_ngap_id_t amf_ue_ngap_id                     = 0;

  container = &pdu->choice.successfulOutcome.value.choice
                   .PDUSessionResourceReleaseResponse;

  NGAP_FIND_PROTOCOLIE_BY_ID(
      Ngap_PDUSessionResourceReleaseResponseIEs_t, ie, container,
      Ngap_ProtocolIE_ID_id_AMF_UE_NGAP_ID, true);
  amf_ue_ngap_id = ie->value.choice.AMF_UE_NGAP_ID;

  if ((ie) &&
      (ue_ref_p = ngap_state_get_ue_amfid((uint32_t) amf_ue_ngap_id)) == NULL) {
    OAILOG_ERROR(
        LOG_NGAP,
        "No UE is attached to this amf UE ngap id: " AMF_UE_NGAP_ID_FMT "\n",
        (amf_ue_ngap_id_t) amf_ue_ngap_id);
    OAILOG_FUNC_RETURN(LOG_NGAP, RETURNerror);
  }

  NGAP_FIND_PROTOCOLIE_BY_ID(
      Ngap_PDUSessionResourceReleaseResponseIEs_t, ie, container,
      Ngap_ProtocolIE_ID_id_RAN_UE_NGAP_ID, true);
  // gNB UE NGAP ID is limited to 24 bits
  gnb_ue_ngap_id =
      (gnb_ue_ngap_id_t)(ie->value.choice.RAN_UE_NGAP_ID & GNB_UE_NGAP_ID_MASK);

  if ((ie) && ue_ref_p->gnb_ue_ngap_id != gnb_ue_ngap_id) {
    OAILOG_ERROR(
        LOG_NGAP,
        "Mismatch in gNB UE NGAP ID, known: " GNB_UE_NGAP_ID_FMT
        ", received: " GNB_UE_NGAP_ID_FMT "\n",
        ue_ref_p->gnb_ue_ngap_id, (gnb_ue_ngap_id_t) gnb_ue_ngap_id);
    OAILOG_FUNC_RETURN(LOG_NGAP, RETURNerror);
  }

  ngap_imsi_map_t* imsi_map = get_ngap_imsi_map();
  hashtable_uint64_ts_get(
      imsi_map->amf_ue_id_imsi_htbl,
      (const hash_key_t) ie->value.choice.AMF_UE_NGAP_ID, &imsi64);

  message_p =
      itti_alloc_new_message(TASK_NGAP, NGAP_PDUSESSIONRESOURCE_REL_RSP);
  if (message_p == NULL) {
    OAILOG_ERROR(LOG_NGAP, "itti_alloc_new_message Failed\n");
    OAILOG_FUNC_RETURN(LOG_NGAP, RETURNerror);
  }
  NGAP_PDUSESSIONRESOURCE_REL_RSP(message_p).amf_ue_ngap_id =
      ue_ref_p->amf_ue_ngap_id;
  NGAP_PDUSESSIONRESOURCE_REL_RSP(message_p).gnb_ue_ngap_id =
      ue_ref_p->gnb_ue_ngap_id;
  NGAP_PDUSESSIONRESOURCE_REL_RSP(message_p)
      .pduSessionResourceReleasedRspList.no_of_items = 0;

  NGAP_FIND_PROTOCOLIE_BY_ID(
      Ngap_PDUSessionResourceReleaseResponseIEs_t, ie, container,
      Ngap_ProtocolIE_ID_id_PDUSessionResourceReleasedListRelRes, true);

  if (ie) {
    const Ngap_PDUSessionResourceReleasedListRelRes_t* const
        pduSessionResourceList =
            &ie->value.choice.PDUSessionResourceReleasedListRelRes;

    int num_pduSessionResource = pduSessionResourceList->list.count;

    for (int index = 0; index < num_pduSessionResource; index++) {
      const Ngap_PDUSessionResourceReleasedItemRelRes_t* const
          pduSessionResource_item =
              (Ngap_PDUSessionResourceReleasedItemRelRes_t*)
                  pduSessionResourceList->list.array[index];
      NGAP_PDUSESSIONRESOURCE_REL_RSP(message_p)
          .pduSessionResourceReleasedRspList.item[index]
          .Pdu_Session_ID = pduSessionResource_item->pDUSessionID;
      NGAP_PDUSESSIONRESOURCE_REL_RSP(message_p)
          .pduSessionResourceReleasedRspList.no_of_items++;
    }
  }
  message_p->ittiMsgHeader.imsi = imsi64;
  rc = send_msg_to_task(&ngap_task_zmq_ctx, TASK_AMF_APP, message_p);
  OAILOG_INFO(LOG_NGAP, " PDU RELEASE msg sent to amf\n");
  OAILOG_FUNC_RETURN(LOG_NGAP, rc);
}

int ngap_amf_handle_pduSession_setup_response(
    ngap_state_t* state, const sctp_assoc_id_t assoc_id,
    const sctp_stream_id_t stream, Ngap_NGAP_PDU_t* pdu) {
  OAILOG_FUNC_IN(LOG_NGAP);
  Ngap_PDUSessionResourceSetupResponse_t* container = NULL;
  Ngap_PDUSessionResourceSetupResponseIEs_t* ie     = NULL;
  QosFlowPerTNLInformation_t response_transfer      = {0};
  m5g_ue_description_t* ue_ref_p                    = NULL;
  MessageDef* message_p                             = NULL;
  gnb_ue_ngap_id_t gnb_ue_ngap_id                   = 0;
  amf_ue_ngap_id_t amf_ue_ngap_id                   = 0;
  int rc                                            = RETURNok;
  imsi64_t imsi64                                   = INVALID_IMSI64;
  container = &pdu->choice.successfulOutcome.value.choice
                   .PDUSessionResourceSetupResponse;
  NGAP_FIND_PROTOCOLIE_BY_ID(
      Ngap_PDUSessionResourceSetupResponseIEs_t, ie, container,
      Ngap_ProtocolIE_ID_id_AMF_UE_NGAP_ID, true);
  if (ie) {
    amf_ue_ngap_id = ie->value.choice.AMF_UE_NGAP_ID;
  } else {
    OAILOG_FUNC_RETURN(LOG_NGAP, RETURNerror);
  }
  NGAP_FIND_PROTOCOLIE_BY_ID(
      Ngap_PDUSessionResourceSetupResponseIEs_t, ie, container,
      Ngap_ProtocolIE_ID_id_RAN_UE_NGAP_ID, true);
  if (ie) {
    // gNB UE NGAP ID is limited to 24 bits
    gnb_ue_ngap_id = (gnb_ue_ngap_id_t)(
        ie->value.choice.RAN_UE_NGAP_ID & GNB_UE_NGAP_ID_MASK);
  } else {
    OAILOG_FUNC_RETURN(LOG_NGAP, RETURNerror);
  }
  if ((ue_ref_p = ngap_state_get_ue_amfid((uint32_t) amf_ue_ngap_id)) == NULL) {
    OAILOG_DEBUG(
        LOG_NGAP,
        "No UE is attached to this amf UE ngap id: " AMF_UE_NGAP_ID_FMT "\n",
        amf_ue_ngap_id);
    OAILOG_FUNC_RETURN(LOG_NGAP, RETURNerror);
  }

  if (ue_ref_p->gnb_ue_ngap_id != gnb_ue_ngap_id) {
    OAILOG_DEBUG(
        LOG_NGAP,
        "Mismatch in gNB UE NGAP ID, known: " GNB_UE_NGAP_ID_FMT
        ", received: " GNB_UE_NGAP_ID_FMT "\n",
        ue_ref_p->gnb_ue_ngap_id, gnb_ue_ngap_id);
    OAILOG_FUNC_RETURN(LOG_NGAP, RETURNerror);
  }
  ngap_imsi_map_t* imsi_map = get_ngap_imsi_map();
  hashtable_uint64_ts_get(
      imsi_map->amf_ue_id_imsi_htbl,
      (const hash_key_t) ue_ref_p->amf_ue_ngap_id, &imsi64);

  message_p =
      itti_alloc_new_message(TASK_NGAP, NGAP_PDUSESSIONRESOURCE_SETUP_RSP);
  NGAP_PDUSESSIONRESOURCE_SETUP_RSP(message_p).amf_ue_ngap_id =
      ue_ref_p->amf_ue_ngap_id;
  NGAP_PDUSESSIONRESOURCE_SETUP_RSP(message_p).gnb_ue_ngap_id =
      ue_ref_p->gnb_ue_ngap_id;

  NGAP_PDUSESSIONRESOURCE_SETUP_RSP(message_p)
      .pduSessionResource_setup_list.no_of_items = 0;
  NGAP_PDUSESSIONRESOURCE_SETUP_RSP(message_p)
      .pduSessionResourceFailedToSetupList.no_of_items = 0;

  NGAP_FIND_PROTOCOLIE_BY_ID(
      Ngap_PDUSessionResourceSetupResponseIEs_t, ie, container,
      Ngap_ProtocolIE_ID_id_PDUSessionResourceSetupListSURes, false);

  if (ie) {
    int pduSessionResource =
        ie->value.choice.PDUSessionResourceSetupListSURes.list.count;

    for (int index = 0; index < pduSessionResource; index++) {
      Ngap_PDUSessionResourceSetupItemSURes_t* pduSession_setup_item =
          (Ngap_PDUSessionResourceSetupItemSURes_t*) ie->value.choice
              .PDUSessionResourceSetupListSURes.list.array[index];

      NGAP_PDUSESSIONRESOURCE_SETUP_RSP(message_p)
          .pduSessionResource_setup_list.item[index]
          .Pdu_Session_ID = pduSession_setup_item->pDUSessionID;

      Ngap_PDUSessionResourceSetupResponseTransfer_t*
          pDUSessionResourceSetupResponseTransfer = NULL;
      asn_dec_rval_t decode_result;

      decode_result = aper_decode_complete(
          NULL, &asn_DEF_Ngap_PDUSessionResourceSetupResponseTransfer,
          (void**) &pDUSessionResourceSetupResponseTransfer,
          pduSession_setup_item->pDUSessionResourceSetupResponseTransfer.buf,
          pduSession_setup_item->pDUSessionResourceSetupResponseTransfer.size);

      if (decode_result.code == RC_OK) {
        OAILOG_DEBUG(LOG_NGAP, " Decode Successful ");
      } else {
        OAILOG_ERROR(LOG_NGAP, " Decode Failed ");
      }

      memcpy(
          response_transfer.tunnel.gTP_TEID,
          pDUSessionResourceSetupResponseTransfer->dLQosFlowPerTNLInformation
              .uPTransportLayerInformation.choice.gTPTunnel.gTP_TEID.buf,
          4);

      memcpy(
          response_transfer.tunnel.transportLayerAddress,
          pDUSessionResourceSetupResponseTransfer->dLQosFlowPerTNLInformation
              .uPTransportLayerInformation.choice.gTPTunnel
              .transportLayerAddress.buf,
          4);

      response_transfer.associatedQosFlowList.items = 1;
      response_transfer.associatedQosFlowList.QosFlowIdentifier[0] =
          pDUSessionResourceSetupResponseTransfer->dLQosFlowPerTNLInformation
              .associatedQosFlowList.list.array[0]
              ->qosFlowIdentifier;

      NGAP_PDUSESSIONRESOURCE_SETUP_RSP(message_p)
          .pduSessionResource_setup_list.item[index]
          .PDU_Session_Resource_Setup_Response_Transfer = response_transfer;

      NGAP_PDUSESSIONRESOURCE_SETUP_RSP(message_p)
          .pduSessionResource_setup_list.no_of_items += 1;
    }
  }

  NGAP_FIND_PROTOCOLIE_BY_ID(
      Ngap_PDUSessionResourceSetupResponseIEs_t, ie, container,
      Ngap_ProtocolIE_ID_id_PDUSessionResourceFailedToSetupListSURes, false);
  if (ie) {
    int pduSessionResource =
        ie->value.choice.PDUSessionResourceFailedToSetupListSURes.list.count;
    for (int index = 0; index < pduSessionResource; index++) {
      Ngap_PDUSessionResourceFailedToSetupItemSURes_t* pduSessionResource_item =
          (Ngap_PDUSessionResourceFailedToSetupItemSURes_t*) ie->value.choice
              .PDUSessionResourceFailedToSetupListSURes.list.array[index];

      NGAP_PDUSESSIONRESOURCE_SETUP_RSP(message_p)
          .pduSessionResourceFailedToSetupList.item[index]
          .Pdu_Session_ID = pduSessionResource_item->pDUSessionID;
      NGAP_PDUSESSIONRESOURCE_SETUP_RSP(message_p)
          .pduSessionResourceFailedToSetupList.item[index]
          .PDU_Session_Resource_Setup_Unsuccessful_Transfer = blk2bstr(
          pduSessionResource_item->pDUSessionResourceSetupUnsuccessfulTransfer
              .buf,
          pduSessionResource_item->pDUSessionResourceSetupUnsuccessfulTransfer
              .size);

      NGAP_PDUSESSIONRESOURCE_SETUP_RSP(message_p)
          .pduSessionResourceFailedToSetupList.no_of_items += 1;
    }
  }

  message_p->ittiMsgHeader.imsi = imsi64;
  rc = send_msg_to_task(&ngap_task_zmq_ctx, TASK_AMF_APP, message_p);
  OAILOG_FUNC_RETURN(LOG_NGAP, rc);
}

//-------------------------------------------------------------------------------
int ngap_handle_paging_request(
    ngap_state_t* state, const itti_ngap_paging_request_t* paging_request,
    imsi64_t imsi64) {
  OAILOG_FUNC_IN(LOG_NGAP);
  DevAssert(paging_request != NULL);
  int rc                  = RETURNok;
  uint16_t tai_list_count = paging_request->TAIListForPaging.no_of_items;
  uint8_t* buffer_p       = NULL;
  uint32_t length         = 0;
  Ngap_NGAP_PDU_t pdu     = {0};
  Ngap_Paging_t* out      = NULL;
  Ngap_PagingIEs_t* ie    = NULL;

  memset(&pdu, 0, sizeof(pdu));
  pdu.present = Ngap_NGAP_PDU_PR_initiatingMessage;
  pdu.choice.initiatingMessage.procedureCode = Ngap_ProcedureCode_id_Paging;
  pdu.choice.initiatingMessage.criticality   = Ngap_Criticality_ignore;
  pdu.choice.initiatingMessage.value.present =
      Ngap_InitiatingMessage__value_PR_Paging;
  out = &pdu.choice.initiatingMessage.value.choice.Paging;

  // UEPagingIdentity
  ie                = (Ngap_PagingIEs_t*) calloc(1, sizeof(Ngap_PagingIEs_t));
  ie->id            = Ngap_ProtocolIE_ID_id_UEPagingIdentity;
  ie->criticality   = Ngap_Criticality_ignore;
  ie->value.present = Ngap_PagingIEs__value_PR_UEPagingIdentity;
  ie->value.choice.UEPagingIdentity.present =
      Ngap_UEPagingIdentity_PR_fiveG_S_TMSI;

  /*BIT_STRING_t*/
  UE_ID_INDEX_TO_BIT_STRING(
      paging_request->UEPagingIdentity.amf_set_id,
      &ie->value.choice.UEPagingIdentity.choice.fiveG_S_TMSI
           .aMFSetID);  // 10  bits

  /*BIT_STRING_t*/
  AMF_POINTER_TO_BIT_STRING(
      paging_request->UEPagingIdentity.amf_pointer,
      &ie->value.choice.UEPagingIdentity.choice.fiveG_S_TMSI
           .aMFPointer);  // 6 bits

  // OCTET_STRING_t
  INT32_TO_OCTET_STRING(
      paging_request->UEPagingIdentity.m_tmsi,
      &ie->value.choice.UEPagingIdentity.choice.fiveG_S_TMSI.fiveG_TMSI);

  ASN_SEQUENCE_ADD(&out->protocolIEs.list, ie);

  // PagingDRX
  ie                = (Ngap_PagingIEs_t*) calloc(1, sizeof(Ngap_PagingIEs_t));
  ie->id            = Ngap_ProtocolIE_ID_id_PagingDRX;
  ie->criticality   = Ngap_Criticality_ignore;
  ie->value.present = Ngap_PagingIEs__value_PR_PagingDRX;
  ie->value.choice.PagingDRX = paging_request->default_paging_drx;
  ASN_SEQUENCE_ADD(&out->protocolIEs.list, ie);

  // Set TAI list
  ie                = (Ngap_PagingIEs_t*) calloc(1, sizeof(Ngap_PagingIEs_t));
  ie->id            = Ngap_ProtocolIE_ID_id_TAIListForPaging;
  ie->criticality   = Ngap_Criticality_ignore;
  ie->value.present = Ngap_PagingIEs__value_PR_TAIListForPaging;

  Ngap_TAIListForPaging_t* tai_list = &ie->value.choice.TAIListForPaging;

  for (int tai_idx = 0; tai_idx < tai_list_count; tai_idx++) {
    Ngap_TAIListForPagingItem_t* tai_item_ies =
        calloc(1, sizeof(Ngap_TAIListForPagingItem_t));

    if (tai_item_ies == NULL) {
      OAILOG_ERROR_UE(LOG_NGAP, imsi64, "Failed to allocate memory\n");
      OAILOG_FUNC_RETURN(LOG_NGAP, RETURNerror);
    }
    PLMN_T_TO_PLMNID(
        paging_request->TAIListForPaging.tai_list[tai_idx].plmn,
        &tai_item_ies->tAI.pLMNIdentity);

    TAC_TO_ASN1_5G(
        paging_request->TAIListForPaging.tai_list[tai_idx].tac,
        &tai_item_ies->tAI.tAC);

    ASN_SEQUENCE_ADD(&tai_list->list, tai_item_ies);
  }
  ASN_SEQUENCE_ADD(&out->protocolIEs.list, ie);

  // PagingPriority
  ie                = (Ngap_PagingIEs_t*) calloc(1, sizeof(Ngap_PagingIEs_t));
  ie->id            = Ngap_ProtocolIE_ID_id_PagingPriority;
  ie->criticality   = Ngap_Criticality_ignore;
  ie->value.present = Ngap_PagingIEs__value_PR_PagingPriority;
  ie->value.choice.PagingPriority = paging_request->PagingPriority;
  ASN_SEQUENCE_ADD(&out->protocolIEs.list, ie);

  int err = 0;
  if (ngap_amf_encode_pdu(&pdu, &buffer_p, &length) < 0) {
    OAILOG_ERROR(LOG_NGAP, "Failed to encode \n");
    err = 1;
  }

  if (length <= 0) {
    err = 1;
  }
  if (err) {
    OAILOG_FUNC_RETURN(LOG_NGAP, RETURNerror);
  }

  /*Fetching gNB list to send paging request message*/
  hashtable_element_array_t* gnb_array = NULL;
  gnb_description_t* gnb_ref_p         = NULL;
  if (state == NULL) {
    OAILOG_ERROR(LOG_NGAP, "gNB Information is NULL!\n");
    OAILOG_FUNC_RETURN(LOG_NGAP, RETURNerror);
  }
  gnb_array = hashtable_ts_get_elements(&state->gnbs);
  if (gnb_array == NULL) {
    OAILOG_ERROR(LOG_NGAP, "Could not find gNB hashlist!\n");
    OAILOG_FUNC_RETURN(LOG_NGAP, RETURNerror);
  }

  for (int idx = 0; idx < gnb_array->num_elements; idx++) {
    gnb_ref_p = (gnb_description_t*) gnb_array->elements[idx];
    if (gnb_ref_p) {
      bstring paging_msg_buffer = blk2bstr(buffer_p, length);

      rc = ngap_amf_itti_send_sctp_request(
          &paging_msg_buffer, gnb_ref_p->sctp_assoc_id,
          0,   // Stream id 0 for non UE related
               // NGAP message
          0);  // amf_ue_ngap_id 0 because UE in idl
    }
  }

  free(buffer_p);
  if (rc != RETURNok) {
    OAILOG_ERROR(LOG_NGAP, "Failed to send paging message over sctp \n");
  } else {
    OAILOG_INFO(LOG_NGAP, "Sent paging message over sctp  \n");
  }

  OAILOG_FUNC_RETURN(LOG_NGAP, rc);
}

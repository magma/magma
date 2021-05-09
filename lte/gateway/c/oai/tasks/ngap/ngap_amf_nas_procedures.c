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
#include <stdint.h>
#include <stdlib.h>
#include <string.h>

#include "ngap_amf_encoder.h"
#include "bstrlib.h"
#include "dynamic_memory_check.h"
#include "assertions.h"
#include "hashtable.h"
#include "log.h"
#include "conversions.h"
#include "asn1_conversions.h"
#include "ngap_amf.h"
#include "ngap_amf_nas_procedures.h"
#include "ngap_amf_itti_messaging.h"
#include "service303.h"
#include "3gpp_23.003.h"
#include "3gpp_24.007.h"
#include "3gpp_38.413.h"
#include "INTEGER.h"
#include "OCTET_STRING.h"
#include "Ngap_NGAP-PDU.h"
#include "Ngap_EUTRA-CGI.h"
#include "Ngap_GBR-QosInformation.h"
#include "Ngap_GUAMI.h"
#include "Ngap_NAS-PDU.h"
#include "Ngap_PLMNIdentity.h"
#include "Ngap_ProcedureCode.h"
#include "Ngap_ProtocolIE-Field.h"
#include "Ngap_SecurityKey.h"
#include "Ngap_TAI.h"
#include "Ngap_TransportLayerAddress.h"
#include "Ngap_UEAggregateMaximumBitRate.h"
#include "Ngap_UESecurityCapabilities.h"
#include "TrackingAreaIdentity.h"
#include "asn_SEQUENCE_OF.h"
#include "ngap_state.h"
#include "Ngap_CauseMisc.h"
#include "Ngap_CauseNas.h"
#include "Ngap_CauseProtocol.h"
#include "Ngap_CauseRadioNetwork.h"
#include "Ngap_CauseTransport.h"
#include "Ngap_InitialUEMessage.h"
#include "ngap_amf_handlers.h"
#include "ngap_common.h"

//------------------------------------------------------------------------------
int ngap_amf_handle_initial_ue_message(
    ngap_state_t* state, const sctp_assoc_id_t assoc_id,
    const sctp_stream_id_t stream, Ngap_NGAP_PDU_t* pdu) {
  Ngap_InitialUEMessage_t* container;
  Ngap_InitialUEMessage_IEs_t *ie = NULL, *ie_e_tmsi, *ie_csg_id = NULL,
                              *ie_guamfi, *ie_cause;
  m5g_ue_description_t* ue_ref    = NULL;
  gnb_description_t* gNB_ref      = NULL;
  gnb_ue_ngap_id_t gnb_ue_ngap_id = 0;

  OAILOG_FUNC_IN(LOG_NGAP);
  container = &pdu->choice.initiatingMessage.value.choice.InitialUEMessage;

  NGAP_FIND_PROTOCOLIE_BY_ID(
      Ngap_InitialUEMessage_IEs_t, ie, container,
      Ngap_ProtocolIE_ID_id_RAN_UE_NGAP_ID, true);

  OAILOG_INFO(
      LOG_NGAP,
      "Received NGAP INITIAL_UE_MESSAGE GNB_UE_NGAP_ID " GNB_UE_NGAP_ID_FMT
      "\n",
      (gnb_ue_ngap_id_t) ie->value.choice.RAN_UE_NGAP_ID);

  if ((gNB_ref = ngap_state_get_gnb(state, assoc_id)) == NULL) {
    OAILOG_ERROR(LOG_NGAP, "Unknown gNB on assoc_id %d\n", assoc_id);
    OAILOG_FUNC_RETURN(LOG_NGAP, RETURNerror);
  }

  // gNB UE NGAP ID is limited to 24 bits
  gnb_ue_ngap_id =
      (gnb_ue_ngap_id_t)(ie->value.choice.RAN_UE_NGAP_ID & GNB_UE_NGAP_ID_MASK);
  OAILOG_ERROR(
      LOG_NGAP,
      "New Initial UE message received with gNB UE NGAP ID: " GNB_UE_NGAP_ID_FMT
      "\n",
      gnb_ue_ngap_id);
  ue_ref = ngap_state_get_ue_gnbid(gNB_ref->sctp_assoc_id, gnb_ue_ngap_id);

  if (ue_ref == NULL) {
    tai_t tai       = {0};
    guamfi_t guamfi = {
        .plmn         = {0},
        .amf_set_id   = 0,
        .amf_regionid = 0};  // initialized after
                             //.plmn = {0}, .amf_code = 0, .amf_gid = 0};  //
                             // initialized after
    s_tmsi_m5_t s_tmsi = {.amf_set_id = 0, .m_tmsi = INVALID_M_TMSI};
    ecgi_t ecgi        = {.plmn = {0}, .cell_identity = {0}};
    csg_id_t csg_id    = 0;

    /*
     * This UE gNB Id has currently no known ng association.
     * * * * Create new UE context by associating new amf_ue_ngap_id.
     * * * * Update gNB UE list.
     * * * * Forward message to NAS.
     */
    if ((ue_ref = ngap_new_ue(state, assoc_id, gnb_ue_ngap_id)) == NULL) {
      // If we failed to allocate a new UE return -1
      OAILOG_ERROR(
          LOG_NGAP,
          "Initial UE Message- Failed to allocate NGAP UE Context, "
          "gNB UE NGAP ID:" GNB_UE_NGAP_ID_FMT "\n",
          gnb_ue_ngap_id);
      OAILOG_FUNC_RETURN(LOG_NGAP, RETURNerror);
    }
    ue_ref->ng_ue_state = NGAP_UE_WAITING_CSR;

    ue_ref->gnb_ue_ngap_id = gnb_ue_ngap_id;
    // Will be allocated by NAS
    ue_ref->amf_ue_ngap_id = INVALID_AMF_UE_NGAP_ID;

    ue_ref->ngap_ue_context_rel_timer.id  = NGAP_TIMER_INACTIVE_ID;
    ue_ref->ngap_ue_context_rel_timer.sec = NGAP_UE_CONTEXT_REL_COMP_TIMER;

    // On which stream we received the message
    ue_ref->sctp_stream_recv = stream;
    ue_ref->sctp_stream_send = gNB_ref->next_sctp_stream;

    /*
     * Increment the sctp stream for the gNB association.
     * If the next sctp stream is >= instream negotiated between gNB and AMF,
     * wrap to first stream.
     * TODO: search for the first available stream instead.
     */

    /*
     * TODO task#15456359.
     * Below logic seems to be incorrect , revisit it.
     */
    gNB_ref->next_sctp_stream += 1;
    if (gNB_ref->next_sctp_stream >= gNB_ref->instreams) {
      gNB_ref->next_sctp_stream = 1;
    }
    // ngap_dump_gnb(gNB_ref); //TODO implement later
    NGAP_FIND_PROTOCOLIE_BY_ID(
        Ngap_InitialUEMessage_IEs_t, ie, container,
        Ngap_ProtocolIE_ID_id_UserLocationInformation, true);

    OCTET_STRING_TO_TAC_5G(
        &ie->value.choice.UserLocationInformation.choice
             .userLocationInformationNR.tAI.tAC,
        tai.tac);
    DevAssert(
        ie->value.choice.UserLocationInformation.choice
            .userLocationInformationNR.tAI.pLMNIdentity.size == 3);
    NGAP_FIND_PROTOCOLIE_BY_ID(
        Ngap_InitialUEMessage_IEs_t, ie, container,
        Ngap_ProtocolIE_ID_id_UserLocationInformation, true);
    DevAssert(
        ie->value.choice.UserLocationInformation.choice
            .userLocationInformationEUTRA.eUTRA_CGI.pLMNIdentity.size == 3);
    TBCD_TO_PLMN_T(
        &ie->value.choice.UserLocationInformation.choice
             .userLocationInformationEUTRA.eUTRA_CGI.pLMNIdentity,
        &ecgi.plmn);
    BIT_STRING_TO_CELL_IDENTITY(
        &ie->value.choice.UserLocationInformation.choice
             .userLocationInformationEUTRA.eUTRA_CGI.eUTRACellIdentity,
        ecgi.cell_identity);

    /** Set the GNB Id. */
    ecgi.cell_identity.enb_id = gNB_ref->gnb_id;

    NGAP_FIND_PROTOCOLIE_BY_ID(
        Ngap_InitialUEMessage_IEs_t, ie_e_tmsi, container,
        Ngap_ProtocolIE_ID_id_FiveG_S_TMSI, false);
    if (ie_e_tmsi) {
      OCTET_STRING_TO_AMF_CODE(
          &ie_e_tmsi->value.choice.FiveG_S_TMSI.aMFSetID, s_tmsi.amf_set_id);
      OCTET_STRING_TO_M_TMSI(
          &ie_e_tmsi->value.choice.FiveG_S_TMSI.fiveG_TMSI, s_tmsi.m_tmsi);
    }

    NGAP_FIND_PROTOCOLIE_BY_ID(
        Ngap_InitialUEMessage_IEs_t, ie_guamfi, container,
        Ngap_ProtocolIE_ID_id_GUAMI, false);
    memset(&guamfi, 0, sizeof(guamfi));

    /*
     * We received the first NAS transport message: initial UE message.
     * * * * Send a NAS ESTAgNBBLISH IND to NAS layer
     */
    NGAP_FIND_PROTOCOLIE_BY_ID(
        Ngap_InitialUEMessage_IEs_t, ie, container,
        Ngap_ProtocolIE_ID_id_NAS_PDU, true);
    NGAP_FIND_PROTOCOLIE_BY_ID(
        Ngap_InitialUEMessage_IEs_t, ie_cause, container,
        Ngap_ProtocolIE_ID_id_RRCEstablishmentCause, true);
    ngap_amf_itti_ngap_initial_ue_message(
        assoc_id, gNB_ref->gnb_id, ue_ref->gnb_ue_ngap_id,
        ie->value.choice.NAS_PDU.buf, ie->value.choice.NAS_PDU.size, &tai,
        &ecgi, ie_cause->value.choice.RRCEstablishmentCause,
        ie_e_tmsi ? &s_tmsi : NULL, ie_csg_id ? &csg_id : NULL,
        ie_guamfi ? &guamfi : NULL,
        NULL,  // CELL ACCESS MODE
        NULL,  // GW Transport Layer Address
        NULL   // Relay Node Indicator
    );

  } else {
    OAILOG_ERROR(
        LOG_NGAP,
        "Initial UE Message- Duplicate GNB_UE_NGAP_ID. Ignoring the "
        "message, gNB UE NGAP ID:" GNB_UE_NGAP_ID_FMT "\n",
        gnb_ue_ngap_id);
  }

  OAILOG_FUNC_RETURN(LOG_NGAP, RETURNok);
}

//------------------------------------------------------------------------------

int ngap_amf_handle_uplink_nas_transport(
    ngap_state_t* state, const sctp_assoc_id_t assoc_id,
    __attribute__((unused)) const sctp_stream_id_t stream,
    Ngap_NGAP_PDU_t* pdu) {
  Ngap_UplinkNASTransport_t* container = NULL;
  Ngap_UplinkNASTransport_IEs_t *ie, *ie_nas_pdu = NULL;
  m5g_ue_description_t* ue_ref    = NULL;
  gnb_description_t* gnb_ref      = NULL;
  tai_t tai                       = {0};
  ecgi_t ecgi                     = {.plmn = {0}, .cell_identity = {0}};
  amf_ue_ngap_id_t amf_ue_ngap_id = 0;
  gnb_ue_ngap_id_t gnb_ue_ngap_id = 0;

  OAILOG_FUNC_IN(LOG_NGAP);
  container = &pdu->choice.initiatingMessage.value.choice.UplinkNASTransport;

  NGAP_FIND_PROTOCOLIE_BY_ID(
      Ngap_UplinkNASTransport_IEs_t, ie, container,
      Ngap_ProtocolIE_ID_id_RAN_UE_NGAP_ID, true);
  gnb_ue_ngap_id = ie->value.choice.RAN_UE_NGAP_ID;

  NGAP_FIND_PROTOCOLIE_BY_ID(
      Ngap_UplinkNASTransport_IEs_t, ie, container,
      Ngap_ProtocolIE_ID_id_AMF_UE_NGAP_ID, true);
  amf_ue_ngap_id = ie->value.choice.AMF_UE_NGAP_ID;

  if (INVALID_AMF_UE_NGAP_ID == ie->value.choice.AMF_UE_NGAP_ID) {
    OAILOG_WARNING(
        LOG_NGAP,
        "Received NGAP UPLINK_NAS_TRANSPORT message AMF_UE_NGAP_ID unknown\n");

    gnb_ref = ngap_state_get_gnb(state, assoc_id);

    if (!(ue_ref = ngap_state_get_ue_gnbid(
              gnb_ref->sctp_assoc_id, (gnb_ue_ngap_id_t) gnb_ue_ngap_id))) {
      OAILOG_WARNING(
          LOG_NGAP,
          "Received NGAP UPLINK_NAS_TRANSPORT No UE is attached to this "
          "gnb_ue_ngap_id: " GNB_UE_NGAP_ID_FMT "\n",
          (gnb_ue_ngap_id_t) gnb_ue_ngap_id);
      OAILOG_FUNC_RETURN(LOG_NGAP, RETURNerror);
    }
  } else {
    OAILOG_INFO(
        LOG_NGAP,
        "Received NGAP UPLINK_NAS_TRANSPORT message "
        "AMF_UE_NGAP_ID " AMF_UE_NGAP_ID_FMT "\n",
        (amf_ue_ngap_id_t) amf_ue_ngap_id);

    if (!(ue_ref = ngap_state_get_ue_amfid(amf_ue_ngap_id))) {
      OAILOG_WARNING(
          LOG_NGAP,
          "Received NGAP UPLINK_NAS_TRANSPORT No UE is attached to this "
          "amf_ue_ngap_id: " AMF_UE_NGAP_ID_FMT "\n",
          (amf_ue_ngap_id_t) amf_ue_ngap_id);
      OAILOG_FUNC_RETURN(LOG_NGAP, RETURNerror);
    }
  }

  if (NGAP_UE_CONNECTED != ue_ref->ng_ue_state) {
    OAILOG_WARNING(
        LOG_NGAP,
        "Received NGAP UPLINK_NAS_TRANSPORT while UE in state != "
        "NGAP_UE_CONNECTED\n");

    OAILOG_FUNC_RETURN(LOG_NGAP, RETURNerror);
  }

  NGAP_FIND_PROTOCOLIE_BY_ID(
      Ngap_UplinkNASTransport_IEs_t, ie_nas_pdu, container,
      Ngap_ProtocolIE_ID_id_NAS_PDU, true);
  // TAI mandatory IE
  NGAP_FIND_PROTOCOLIE_BY_ID(
      Ngap_UplinkNASTransport_IEs_t, ie, container,
      Ngap_ProtocolIE_ID_id_SupportedTAList, true);
  OCTET_STRING_TO_TAC(
      &ie->value.choice.UserLocationInformation.choice.userLocationInformationNR
           .tAI.tAC,
      tai.tac);
  DevAssert(
      ie->value.choice.UserLocationInformation.choice.userLocationInformationNR
          .tAI.pLMNIdentity.size == 3);
  TBCD_TO_PLMN_T(
      &ie->value.choice.UserLocationInformation.choice.userLocationInformationNR
           .tAI.pLMNIdentity,
      &tai.plmn);

  // CGI mandatory IE
  NGAP_FIND_PROTOCOLIE_BY_ID(
      Ngap_UplinkNASTransport_IEs_t, ie, container,
      Ngap_ProtocolIE_ID_id_EUTRA_CGI, true);
  DevAssert(
      ie->value.choice.UserLocationInformation.choice
          .userLocationInformationEUTRA.eUTRA_CGI.pLMNIdentity.size == 3);
  TBCD_TO_PLMN_T(
      &ie->value.choice.UserLocationInformation.choice
           .userLocationInformationEUTRA.eUTRA_CGI.pLMNIdentity,
      &ecgi.plmn);
  BIT_STRING_TO_CELL_IDENTITY(
      &ie->value.choice.UserLocationInformation.choice
           .userLocationInformationEUTRA.eUTRA_CGI.eUTRACellIdentity,
      ecgi.cell_identity);

  // TODO optional GW Transport Layer Address

  bstring b = blk2bstr(
      ie_nas_pdu->value.choice.NAS_PDU.buf,
      ie_nas_pdu->value.choice.NAS_PDU.size);
  ngap_amf_itti_nas_uplink_ind(amf_ue_ngap_id, &b, &tai, &ecgi);
  OAILOG_FUNC_RETURN(LOG_NGAP, RETURNok);
}

//------------------------------------------------------------------------------
int ngap_amf_handle_nas_non_delivery(
    ngap_state_t* state, __attribute__((unused)) sctp_assoc_id_t assoc_id,
    sctp_stream_id_t stream, Ngap_NGAP_PDU_t* pdu) {
  Ngap_NASNonDeliveryIndication_t* container;
  Ngap_NASNonDeliveryIndication_IEs_t *ie = NULL, *ie_nas_pdu;
  m5g_ue_description_t* ue_ref            = NULL;
  imsi64_t imsi64                         = INVALID_IMSI64;
  amf_ue_ngap_id_t amf_ue_ngap_id         = 0;
  gnb_ue_ngap_id_t gnb_ue_ngap_id         = 0;

  OAILOG_FUNC_IN(LOG_NGAP);
  increment_counter("nas_non_delivery_indication_received", 1, NO_LABELS);

  container =
      &pdu->choice.initiatingMessage.value.choice.NASNonDeliveryIndication;
  /*
   * UE associated signaling on stream == 0 is not valid.
   */
  if (stream == 0) {
    OAILOG_NOTICE(
        LOG_NGAP,
        "Received NGAP NAS_NON_DELIVERY_INDICATION message on invalid sctp "
        "stream 0\n");
    OAILOG_FUNC_RETURN(LOG_NGAP, RETURNerror);
  }
  NGAP_FIND_PROTOCOLIE_BY_ID(
      Ngap_NASNonDeliveryIndication_IEs_t, ie, container,
      Ngap_ProtocolIE_ID_id_AMF_UE_NGAP_ID, true);
  amf_ue_ngap_id = ie->value.choice.AMF_UE_NGAP_ID;

  NGAP_FIND_PROTOCOLIE_BY_ID(
      Ngap_NASNonDeliveryIndication_IEs_t, ie, container,
      Ngap_ProtocolIE_ID_id_RAN_UE_NGAP_ID, true);
  gnb_ue_ngap_id = ie->value.choice.RAN_UE_NGAP_ID;

  NGAP_FIND_PROTOCOLIE_BY_ID(
      Ngap_NASNonDeliveryIndication_IEs_t, ie_nas_pdu, container,
      Ngap_ProtocolIE_ID_id_NAS_PDU, true);

  NGAP_FIND_PROTOCOLIE_BY_ID(
      Ngap_NASNonDeliveryIndication_IEs_t, ie, container,
      Ngap_ProtocolIE_ID_id_Cause, true);
  OAILOG_NOTICE(
      LOG_NGAP,
      "Received NGAP NAS_NON_DELIVERY_INDICATION message "
      "AMF_UE_NGAP_ID " AMF_UE_NGAP_ID_FMT " gnb_ue_ngap_id " GNB_UE_NGAP_ID_FMT
      "\n",
      amf_ue_ngap_id, gnb_ue_ngap_id);

  if ((ue_ref = ngap_state_get_ue_amfid(amf_ue_ngap_id)) == NULL) {
    OAILOG_DEBUG(
        LOG_NGAP,
        "No UE is attached to this amf UE ngap id: " AMF_UE_NGAP_ID_FMT "\n",
        amf_ue_ngap_id);
    OAILOG_FUNC_RETURN(LOG_NGAP, RETURNerror);
  }

  ngap_imsi_map_t* imsi_map = get_ngap_imsi_map();
  hashtable_uint64_ts_get(
      imsi_map->amf_ue_id_imsi_htbl, amf_ue_ngap_id, &imsi64);

  if (ue_ref->ng_ue_state != NGAP_UE_CONNECTED) {
    OAILOG_DEBUG_UE(
        LOG_NGAP, imsi64,
        "Received NGAP NAS_NON_DELIVERY_INDICATION while UE in state != "
        "NGAP_UE_CONNECTED\n");
    OAILOG_FUNC_RETURN(LOG_NGAP, RETURNerror);
  }

  // TODO: forward NAS PDU to NAS
  ngap_amf_itti_nas_non_delivery_ind(
      ie->value.choice.AMF_UE_NGAP_ID, ie->value.choice.NAS_PDU.buf,
      ie->value.choice.NAS_PDU.size, &ie->value.choice.Cause, imsi64);
  OAILOG_FUNC_RETURN(LOG_NGAP, RETURNok);
}

//------------------------------------------------------------------------------
int ngap_generate_downlink_nas_transport(
    ngap_state_t* state, const gnb_ue_ngap_id_t gnb_ue_ngap_id,
    const amf_ue_ngap_id_t ue_id, STOLEN_REF bstring* payload,
    const imsi64_t imsi64) {
  m5g_ue_description_t* ue_ref = NULL;
  uint8_t* buffer_p            = NULL;
  uint32_t length              = 0;
  void* id                     = NULL;

  OAILOG_FUNC_IN(LOG_NGAP);

  // Try to retrieve SCTP association id using amf_ue_ngap_id
  if (HASH_TABLE_OK ==
      hashtable_ts_get(
          &state->amfid2associd, (const hash_key_t) ue_id, (void**) &id)) {
    sctp_assoc_id_t sctp_assoc_id = (sctp_assoc_id_t)(uintptr_t) id;
    gnb_description_t* gnb_ref    = ngap_state_get_gnb(state, sctp_assoc_id);
    if (gnb_ref) {
      ue_ref = ngap_state_get_ue_gnbid(gnb_ref->sctp_assoc_id, gnb_ue_ngap_id);
    } else {
      OAILOG_ERROR(
          LOG_NGAP, "No gNB for SCTP association id %d \n", sctp_assoc_id);
      OAILOG_FUNC_RETURN(LOG_NGAP, RETURNerror);
    }
  }
  // TODO remove soon:
  if (!ue_ref) {
    ue_ref = ngap_state_get_ue_amfid(ue_id);
  }
  // finally!
  if (!ue_ref) {
    /*
     * If the UE-associated logical NG-connection is not established,
     * * * * the AMF shall allocate a unique AMF UE NGAP ID to be used for the
     * UE.
     */
    OAILOG_WARNING(
        LOG_NGAP,
        "Unknown UE AMF ID " AMF_UE_NGAP_ID_FMT
        ", This case is not handled right now\n",
        ue_id);
    OAILOG_FUNC_RETURN(LOG_NGAP, RETURNerror);
  } else {
    /*
     * We have fount the UE in the list.
     * * * * Create new IE list message and encode it.
     */
    ngap_imsi_map_t* imsi_map = get_ngap_imsi_map();
    hashtable_uint64_ts_insert(
        imsi_map->amf_ue_id_imsi_htbl, (const hash_key_t) ue_id, imsi64);

    Ngap_DownlinkNASTransport_IEs_t* ie = NULL;
    Ngap_DownlinkNASTransport_t* out;
    Ngap_NGAP_PDU_t pdu = {0};

    memset(&pdu, 0, sizeof(pdu));
    pdu.present = Ngap_NGAP_PDU_PR_initiatingMessage;
    pdu.choice.initiatingMessage.procedureCode =
        Ngap_ProcedureCode_id_DownlinkNASTransport;
    pdu.choice.initiatingMessage.criticality = Ngap_Criticality_ignore;
    pdu.choice.initiatingMessage.value.present =
        Ngap_InitiatingMessage__value_PR_DownlinkNASTransport;

    out = &pdu.choice.initiatingMessage.value.choice.DownlinkNASTransport;
    if (ue_ref->ng_ue_state == NGAP_UE_WAITING_CRR) {
      OAILOG_ERROR_UE(
          LOG_NGAP, imsi64,
          "Already triggred UE Context Release Command and UE is"
          "in NGAP_UE_WAITING_CRR, so dropping the DownlinkNASTransport \n");
      OAILOG_FUNC_RETURN(LOG_NGAP, RETURNerror);
    } else {
      ue_ref->ng_ue_state = NGAP_UE_CONNECTED;
    }
    /*
     * Setting UE informations with the ones found in ue_ref
     */
    ie = (Ngap_DownlinkNASTransport_IEs_t*) calloc(
        1, sizeof(Ngap_DownlinkNASTransport_IEs_t));
    ie->id            = Ngap_ProtocolIE_ID_id_AMF_UE_NGAP_ID;
    ie->criticality   = Ngap_Criticality_reject;
    ie->value.present = Ngap_DownlinkNASTransport_IEs__value_PR_AMF_UE_NGAP_ID;
    ie->value.choice.AMF_UE_NGAP_ID = ue_ref->amf_ue_ngap_id;
    ASN_SEQUENCE_ADD(&out->protocolIEs.list, ie);

    /* mandatory */
    ie = (Ngap_DownlinkNASTransport_IEs_t*) calloc(
        1, sizeof(Ngap_DownlinkNASTransport_IEs_t));
    ie->id            = Ngap_ProtocolIE_ID_id_RAN_UE_NGAP_ID;
    ie->criticality   = Ngap_Criticality_reject;
    ie->value.present = Ngap_DownlinkNASTransport_IEs__value_PR_RAN_UE_NGAP_ID;
    ie->value.choice.RAN_UE_NGAP_ID = ue_ref->gnb_ue_ngap_id;
    ASN_SEQUENCE_ADD(&out->protocolIEs.list, ie);
    /* mandatory */
    ie = (Ngap_DownlinkNASTransport_IEs_t*) calloc(
        1, sizeof(Ngap_DownlinkNASTransport_IEs_t));
    ie->id            = Ngap_ProtocolIE_ID_id_NAS_PDU;
    ie->criticality   = Ngap_Criticality_reject;
    ie->value.present = Ngap_DownlinkNASTransport_IEs__value_PR_NAS_PDU;
    /*gNB
     * Fill in the NAS pdu
     */
    OCTET_STRING_fromBuf(
        &ie->value.choice.NAS_PDU, (char*) bdata(*payload), blength(*payload));
    ASN_SEQUENCE_ADD(&out->protocolIEs.list, ie);

    if (ngap_amf_encode_pdu(&pdu, &buffer_p, &length) < 0) {
      OAILOG_FUNC_RETURN(LOG_NGAP, RETURNerror);
    }

    OAILOG_NOTICE_UE(
        LOG_NGAP, imsi64,
        "Send NGAP DOWNLINK_NAS_TRANSPORT message ue_id = " AMF_UE_NGAP_ID_FMT
        " AMF_UE_NGAP_ID = " AMF_UE_NGAP_ID_FMT
        " gNB_UE_NGAP_ID = " GNB_UE_NGAP_ID_FMT "\n",
        ue_id, ue_ref->amf_ue_ngap_id, gnb_ue_ngap_id);
    bstring b = blk2bstr(buffer_p, length);
    free(buffer_p);
    ngap_amf_itti_send_sctp_request(
        &b, ue_ref->sctp_assoc_id, ue_ref->sctp_stream_send,
        ue_ref->amf_ue_ngap_id);
  }

  OAILOG_FUNC_RETURN(LOG_NGAP, RETURNok);
}

//------------------------------------------------------------------------------
void ngap_handle_amf_ue_id_notification(
    ngap_state_t* state,
    const itti_amf_app_ngap_amf_ue_id_notification_t* const notification_p) {
  OAILOG_FUNC_IN(LOG_NGAP);

  DevAssert(notification_p != NULL);

  sctp_assoc_id_t sctp_assoc_id   = notification_p->sctp_assoc_id;
  gnb_ue_ngap_id_t gnb_ue_ngap_id = notification_p->gnb_ue_ngap_id;
  amf_ue_ngap_id_t amf_ue_ngap_id = notification_p->amf_ue_ngap_id;

  gnb_description_t* gnb_ref = ngap_state_get_gnb(state, sctp_assoc_id);
  if (gnb_ref) {
    m5g_ue_description_t* ue_ref =
        ngap_state_get_ue_gnbid(gnb_ref->sctp_assoc_id, gnb_ue_ngap_id);
    if (ue_ref) {
      ue_ref->amf_ue_ngap_id = amf_ue_ngap_id;
      hashtable_rc_t h_rc    = hashtable_ts_insert(
          &state->amfid2associd, (const hash_key_t) amf_ue_ngap_id,
          (void*) (uintptr_t) sctp_assoc_id);

      hashtable_uint64_ts_insert(
          &gnb_ref->ue_id_coll, (const hash_key_t) amf_ue_ngap_id,
          ue_ref->comp_ngap_id);

      OAILOG_DEBUG(
          LOG_NGAP,
          "Associated sctp_assoc_id %d, gnb_ue_ngap_id " GNB_UE_NGAP_ID_FMT
          ", amf_ue_ngap_id " AMF_UE_NGAP_ID_FMT ":%s \n",
          sctp_assoc_id, gnb_ue_ngap_id, amf_ue_ngap_id,
          hashtable_rc_code2string(h_rc));
      return;
    }
    OAILOG_DEBUG(
        LOG_NGAP,
        "Could not find  ue  with gnb_ue_ngap_id " GNB_UE_NGAP_ID_FMT "\n",
        gnb_ue_ngap_id);
    return;
  }
  OAILOG_DEBUG(
      LOG_NGAP, "Could not find  gNB with sctp_assoc_id %d \n", sctp_assoc_id);

  OAILOG_FUNC_OUT(LOG_NGAP);
}

//------------------------------------------------------------------------------
int ngap_generate_ngap_pdusession_resource_setup_req(
    ngap_state_t* state, itti_ngap_pdusession_resource_setup_req_t* const
                             pdusession_resource_setup_req) {
  OAILOG_FUNC_IN(LOG_NGAP);
  m5g_ue_description_t* ue_ref = NULL;
  uint8_t* buffer_p            = NULL;
  uint32_t length              = 0;
  void* id                     = NULL;
  const gnb_ue_ngap_id_t gnb_ue_ngap_id =
      pdusession_resource_setup_req->gnb_ue_ngap_id;
  const amf_ue_ngap_id_t amf_ue_ngap_id =
      pdusession_resource_setup_req->amf_ue_ngap_id;
  pdu_session_resource_setup_request_transfer_t amf_pdu_ses_setup_transfer_req;
  hashtable_ts_get(
      &state->amfid2associd, (const hash_key_t) amf_ue_ngap_id, (void**) &id);
  if (id) {
    sctp_assoc_id_t sctp_assoc_id = (sctp_assoc_id_t)(uintptr_t) id;
    gnb_description_t* gnb_ref    = ngap_state_get_gnb(state, sctp_assoc_id);
    if (gnb_ref) {
      ue_ref = ngap_state_get_ue_gnbid(gnb_ref->sctp_assoc_id, gnb_ue_ngap_id);
    }
  }
  if (!ue_ref) {
    ue_ref = ngap_state_get_ue_amfid(amf_ue_ngap_id);
  }
  if (!ue_ref) {  // tmp for testing
                  /*
                   * If the UE-associated logical Ng-connection is not established,
                   * * * * the AMF shall allocate a unique AMF UE NGAP ID to be used for the
                   * UE.
                   */
    OAILOG_ERROR(
        LOG_NGAP,
        "Unknown UE AMF ID " AMF_UE_NGAP_ID_FMT
        ", This case is not handled right now\n",
        amf_ue_ngap_id);
    OAILOG_FUNC_RETURN(LOG_NGAP, RETURNerror);
  } else {
    /*
     * We have found the UE in the list.
     * Create new IE list message and encode it.
     */
    Ngap_NGAP_PDU_t pdu                                     = {0};
    Ngap_PDUSessionResourceSetupRequest_t* out              = NULL;
    Ngap_PDUSessionResourceSetupRequestIEs_t* ie            = NULL;
    Ngap_PDUSessionResourceSetupRequestTransferIEs_t* tx_ie = NULL;
    memset(&pdu, 0, sizeof(pdu));

    pdu.choice.initiatingMessage.procedureCode =
        Ngap_ProcedureCode_id_PDUSessionResourceSetup;
    pdu.choice.initiatingMessage.criticality = Ngap_Criticality_reject;
    pdu.present = Ngap_NGAP_PDU_PR_initiatingMessage;
    pdu.choice.initiatingMessage.value.present =
        Ngap_InitiatingMessage__value_PR_PDUSessionResourceSetupRequest;
    out = &pdu.choice.initiatingMessage.value.choice
               .PDUSessionResourceSetupRequest;
    // ue_ref->ng_ue_state = NGAP_UE_CONNECTED;
    /*
     * Setting UE information with the ones found in ue_ref
     */
    ie = (Ngap_PDUSessionResourceSetupRequestIEs_t*) calloc(
        1, sizeof(Ngap_PDUSessionResourceSetupRequestIEs_t));
    ie->id          = Ngap_ProtocolIE_ID_id_AMF_UE_NGAP_ID;
    ie->criticality = Ngap_Criticality_reject;
    ie->value.present =
        Ngap_PDUSessionResourceSetupRequestIEs__value_PR_AMF_UE_NGAP_ID;
    ie->value.choice.AMF_UE_NGAP_ID = ue_ref->amf_ue_ngap_id;
    ASN_SEQUENCE_ADD(&out->protocolIEs.list, ie);

    /* mandatory */
    ie = (Ngap_PDUSessionResourceSetupRequestIEs_t*) calloc(
        1, sizeof(Ngap_PDUSessionResourceSetupRequestIEs_t));
    ie->id          = Ngap_ProtocolIE_ID_id_RAN_UE_NGAP_ID;
    ie->criticality = Ngap_Criticality_reject;
    ie->value.present =
        Ngap_PDUSessionResourceSetupRequestIEs__value_PR_RAN_UE_NGAP_ID;
    ie->value.choice.RAN_UE_NGAP_ID = ue_ref->gnb_ue_ngap_id;
    ASN_SEQUENCE_ADD(&out->protocolIEs.list, ie);

    /* mandatory */
    ie = (Ngap_PDUSessionResourceSetupRequestIEs_t*) calloc(
        1, sizeof(Ngap_PDUSessionResourceSetupRequestIEs_t));
    ie->id          = Ngap_ProtocolIE_ID_id_PDUSessionResourceSetupListSUReq;
    ie->criticality = Ngap_Criticality_reject;
    ie->value.present =
        Ngap_PDUSessionResourceSetupRequestIEs__value_PR_PDUSessionResourceSetupListSUReq;
    ASN_SEQUENCE_ADD(&out->protocolIEs.list, ie);

    for (int i = 0; i < pdusession_resource_setup_req
                            ->pduSessionResource_setup_list.no_of_items;
         i++) {
      Ngap_PDUSessionResourceSetupItemSUReq_t* ngap_pdusession_setup_item_ies =
          calloc(1, sizeof(Ngap_PDUSessionResourceSetupItemSUReq_t));

      ngap_pdusession_setup_item_ies->pDUSessionID =
          pdusession_resource_setup_req->pduSessionResource_setup_list.item[i]
              .Pdu_Session_ID;
      /*NSSAI*/
      ngap_pdusession_setup_item_ies->s_NSSAI.sST.size = 1;
      ngap_pdusession_setup_item_ies->s_NSSAI.sST.buf =
          (uint8_t*) calloc(1, sizeof(uint8_t));
      ngap_pdusession_setup_item_ies->s_NSSAI.sST.buf[0] = 0x11;

      // filling PDU TX Structure
      amf_pdu_ses_setup_transfer_req =
          pdusession_resource_setup_req->pduSessionResource_setup_list.item[i]
              .PDU_Session_Resource_Setup_Request_Transfer;

      /*tx_out*/
      Ngap_PDUSessionResourceSetupRequestTransfer_t*
          pduSessionResourceSetupRequestTransferIEs =
              (Ngap_PDUSessionResourceSetupRequestTransfer_t*) calloc(
                  1, sizeof(Ngap_PDUSessionResourceSetupRequestTransfer_t));

      tx_ie = (Ngap_PDUSessionResourceSetupRequestTransferIEs_t*) calloc(
          1, sizeof(Ngap_PDUSessionResourceSetupRequestTransferIEs_t));
      tx_ie->id          = Ngap_ProtocolIE_ID_id_UL_NGU_UP_TNLInformation;
      tx_ie->criticality = Ngap_Criticality_reject;
      tx_ie->value.present =
          Ngap_PDUSessionResourceSetupRequestTransferIEs__value_PR_UPTransportLayerInformation;

      tx_ie->value.choice.UPTransportLayerInformation.present =
          Ngap_UPTransportLayerInformation_PR_gTPTunnel;

      /*transportLayerAddress*/
      tx_ie->value.choice.UPTransportLayerInformation.choice.gTPTunnel
          .transportLayerAddress.buf = calloc(
          blength(amf_pdu_ses_setup_transfer_req.up_transport_layer_info.gtp_tnl
                      .endpoint_ip_address),
          sizeof(uint8_t));

      memcpy(
          tx_ie->value.choice.UPTransportLayerInformation.choice.gTPTunnel
              .transportLayerAddress.buf,
          amf_pdu_ses_setup_transfer_req.up_transport_layer_info.gtp_tnl
              .endpoint_ip_address->data,
          blength(amf_pdu_ses_setup_transfer_req.up_transport_layer_info.gtp_tnl
                      .endpoint_ip_address));

      tx_ie->value.choice.UPTransportLayerInformation.choice.gTPTunnel
          .transportLayerAddress.size =
          blength(amf_pdu_ses_setup_transfer_req.up_transport_layer_info.gtp_tnl
                      .endpoint_ip_address);

      tx_ie->value.choice.UPTransportLayerInformation.choice.gTPTunnel
          .transportLayerAddress.bits_unused = 0;
      OAILOG_DEBUG(
          LOG_NGAP, " blength :%d\n",
          blength(amf_pdu_ses_setup_transfer_req.up_transport_layer_info.gtp_tnl
                      .endpoint_ip_address));

      OCTET_STRING_fromBuf(
          &tx_ie->value.choice.UPTransportLayerInformation.choice.gTPTunnel
               .gTP_TEID,
          (char*) amf_pdu_ses_setup_transfer_req.up_transport_layer_info.gtp_tnl
              .gtp_tied,
          4);

      int ret = ASN_SEQUENCE_ADD(
          &pduSessionResourceSetupRequestTransferIEs->protocolIEs.list, tx_ie);
      if (ret != 0) OAILOG_ERROR(LOG_NGAP, "  encode error ");

      /*PDUSessionType*/
      tx_ie = (Ngap_PDUSessionResourceSetupRequestTransferIEs_t*) calloc(
          1, sizeof(Ngap_PDUSessionResourceSetupRequestTransferIEs_t));
      tx_ie->id          = Ngap_ProtocolIE_ID_id_PDUSessionType;
      tx_ie->criticality = Ngap_Criticality_reject;
      tx_ie->value.present =
          Ngap_PDUSessionResourceSetupRequestTransferIEs__value_PR_PDUSessionType;

      tx_ie->value.choice.PDUSessionType = /*e_Ngap_PDUSessionType*/
          amf_pdu_ses_setup_transfer_req.pdu_ip_type.pdn_type;

      ret = ASN_SEQUENCE_ADD(
          &pduSessionResourceSetupRequestTransferIEs->protocolIEs.list, tx_ie);
      if (ret != 0) OAILOG_ERROR(LOG_NGAP, "  encode  error ");

      /*Qos*/
      tx_ie = (Ngap_PDUSessionResourceSetupRequestTransferIEs_t*) calloc(
          1, sizeof(Ngap_PDUSessionResourceSetupRequestTransferIEs_t));
      tx_ie->id          = Ngap_ProtocolIE_ID_id_QosFlowSetupRequestList;
      tx_ie->criticality = Ngap_Criticality_reject;
      tx_ie->value.present =
          Ngap_PDUSessionResourceSetupRequestTransferIEs__value_PR_QosFlowSetupRequestList;

      for (int i = 0; i < /*no_of_qos_items*/ 1; i++) {
        Ngap_QosFlowSetupRequestItem_t* qos_item =
            (Ngap_QosFlowSetupRequestItem_t*) calloc(
                1, sizeof(Ngap_QosFlowSetupRequestItem_t));

        qos_item->qosFlowIdentifier =
            amf_pdu_ses_setup_transfer_req.qos_flow_setup_request_list
                .qos_flow_req_item.qos_flow_identifier;

        /* Ngap_QosCharacteristics */
        {
          qos_item->qosFlowLevelQosParameters.qosCharacteristics.present =
              Ngap_QosCharacteristics_PR_nonDynamic5QI;
          qos_item->qosFlowLevelQosParameters.qosCharacteristics.choice
              .nonDynamic5QI.fiveQI =
              amf_pdu_ses_setup_transfer_req.qos_flow_setup_request_list
                  .qos_flow_req_item.qos_flow_level_qos_param.qos_characteristic
                  .non_dynamic_5QI_desc.fiveQI;
        }
        /* Ngap_AllocationAndRetentionPriority */
        qos_item->qosFlowLevelQosParameters.allocationAndRetentionPriority
            .priorityLevelARP =
            amf_pdu_ses_setup_transfer_req.qos_flow_setup_request_list
                .qos_flow_req_item.qos_flow_level_qos_param.alloc_reten_priority
                .priority_level;

        qos_item->qosFlowLevelQosParameters.allocationAndRetentionPriority
            .pre_emptionCapability =
            amf_pdu_ses_setup_transfer_req.qos_flow_setup_request_list
                .qos_flow_req_item.qos_flow_level_qos_param.alloc_reten_priority
                .pre_emption_cap;

        qos_item->qosFlowLevelQosParameters.allocationAndRetentionPriority
            .pre_emptionVulnerability =
            amf_pdu_ses_setup_transfer_req.qos_flow_setup_request_list
                .qos_flow_req_item.qos_flow_level_qos_param.alloc_reten_priority
                .pre_emption_vul;

        asn_set_empty(&tx_ie->value.choice.QosFlowSetupRequestList.list);
        ASN_SEQUENCE_ADD(
            &tx_ie->value.choice.QosFlowSetupRequestList.list, qos_item);
      }

      ret = ASN_SEQUENCE_ADD(
          &pduSessionResourceSetupRequestTransferIEs->protocolIEs.list, tx_ie);
      if (ret != 0) OAILOG_ERROR(LOG_NGAP, "  encode error \n");

      uint32_t buffer_size = 512;
      char* buffer         = (char*) calloc(1, buffer_size);

      asn_enc_rval_t er = aper_encode_to_buffer(
          &asn_DEF_Ngap_PDUSessionResourceSetupRequestTransfer, NULL,
          pduSessionResourceSetupRequestTransferIEs, buffer, buffer_size);

      asn_fprint(
          stderr, &asn_DEF_Ngap_PDUSessionResourceSetupRequestTransfer,
          pduSessionResourceSetupRequestTransferIEs);

      bstring transfer = blk2bstr(buffer, er.encoded);
      ngap_pdusession_setup_item_ies->pDUSessionResourceSetupRequestTransfer
          .buf = (uint8_t*) calloc(er.encoded, sizeof(uint8_t));

      memcpy(
          (void*) ngap_pdusession_setup_item_ies
              ->pDUSessionResourceSetupRequestTransfer.buf,
          (void*) transfer->data, er.encoded);

      ngap_pdusession_setup_item_ies->pDUSessionResourceSetupRequestTransfer
          .size = blength(transfer);

      ASN_SEQUENCE_ADD(
          &ie->value.choice.PDUSessionResourceSetupListSUReq.list,
          ngap_pdusession_setup_item_ies);

    } /*for loop*/

    if (ngap_amf_encode_pdu(&pdu, &buffer_p, &length) < 0) {
      OAILOG_ERROR(LOG_NGAP, "Encoding of IEs failed \n");
      OAILOG_FUNC_RETURN(LOG_NGAP, RETURNerror);
    }

    bstring b = blk2bstr(buffer_p, length);
    free(buffer_p);
    ngap_amf_itti_send_sctp_request(
        &b, ue_ref->sctp_assoc_id, ue_ref->sctp_stream_send,
        ue_ref->amf_ue_ngap_id);
  }

  OAILOG_FUNC_RETURN(LOG_NGAP, RETURNok);
}

//------------------------------------------------------------------------------
int ngap_generate_ngap_pdusession_resource_rel_cmd(
    ngap_state_t* state,
    itti_ngap_pdusessionresource_rel_req_t* const pdusessionresource_rel_cmd) {
  OAILOG_FUNC_IN(LOG_NGAP);

  m5g_ue_description_t* ue_ref = NULL;
  uint8_t* buffer_p            = NULL;
  uint32_t length              = 0;
  void* id                     = NULL;
  const gnb_ue_ngap_id_t gnb_ue_ngap_id =
      pdusessionresource_rel_cmd->gnb_ue_ngap_id;
  const amf_ue_ngap_id_t amf_ue_ngap_id =
      pdusessionresource_rel_cmd->amf_ue_ngap_id;

  hashtable_ts_get(
      &state->amfid2associd, (const hash_key_t) amf_ue_ngap_id, (void**) &id);
  if (id) {
    sctp_assoc_id_t sctp_assoc_id = (sctp_assoc_id_t)(uintptr_t) id;
    gnb_description_t* gnb_ref    = ngap_state_get_gnb(state, sctp_assoc_id);
    if (gnb_ref) {
      ue_ref = ngap_state_get_ue_gnbid(gnb_ref->sctp_assoc_id, gnb_ue_ngap_id);
    }
  }
  if (!ue_ref) {
    ue_ref = ngap_state_get_ue_amfid(amf_ue_ngap_id);
  }

  if (!ue_ref) {
    /*
     * If the UE-associated logical Ng-connection is not established,
     * the AMF shall allocate a unique AMF UE NGAP ID to be used for the UE.
     */
    OAILOG_ERROR(
        LOG_NGAP,
        "Unknown UE AMF ID " AMF_UE_NGAP_ID_FMT
        ", This case is not handled right now\n",
        amf_ue_ngap_id);
    OAILOG_FUNC_RETURN(LOG_NGAP, RETURNerror);
  } else {
    /*
     * We have found the UE in the list.
     * Create new IE list message and encode it.
     */
    Ngap_NGAP_PDU_t pdu                            = {0};
    Ngap_PDUSessionResourceReleaseCommand_t* out   = NULL;
    Ngap_PDUSessionResourceReleaseCommandIEs_t* ie = NULL;

    memset(&pdu, 0, sizeof(pdu));
    pdu.present = Ngap_NGAP_PDU_PR_initiatingMessage;
    pdu.choice.initiatingMessage.procedureCode =
        Ngap_ProcedureCode_id_PDUSessionResourceRelease;
    pdu.choice.initiatingMessage.criticality = Ngap_Criticality_ignore;
    pdu.choice.initiatingMessage.value.present =
        Ngap_InitiatingMessage__value_PR_PDUSessionResourceReleaseCommand;
    out = &pdu.choice.initiatingMessage.value.choice
               .PDUSessionResourceReleaseCommand;

    /*
     * Setting UE information with the ones found in ue_ref
     */
    /* mandatory */
    ie = (Ngap_PDUSessionResourceReleaseCommandIEs_t*) calloc(
        1, sizeof(Ngap_PDUSessionResourceReleaseCommandIEs_t));
    ie->id          = Ngap_ProtocolIE_ID_id_AMF_UE_NGAP_ID;
    ie->criticality = Ngap_Criticality_reject;
    ie->value.present =
        Ngap_PDUSessionResourceReleaseCommandIEs__value_PR_AMF_UE_NGAP_ID;
    ie->value.choice.AMF_UE_NGAP_ID = ue_ref->amf_ue_ngap_id;
    ASN_SEQUENCE_ADD(&out->protocolIEs.list, ie);

    /* mandatory */
    ie = (Ngap_PDUSessionResourceReleaseCommandIEs_t*) calloc(
        1, sizeof(Ngap_PDUSessionResourceReleaseCommandIEs_t));
    ie->id          = Ngap_ProtocolIE_ID_id_RAN_UE_NGAP_ID;
    ie->criticality = Ngap_Criticality_reject;
    ie->value.present =
        Ngap_PDUSessionResourceReleaseCommandIEs__value_PR_RAN_UE_NGAP_ID;
    ie->value.choice.RAN_UE_NGAP_ID = ue_ref->gnb_ue_ngap_id;
    ASN_SEQUENCE_ADD(&out->protocolIEs.list, ie);

    //    ue_ref->ng_ue_state = NGAP_UE_CONNECTED;

    /* optional NAS pdu */
    if (pdusessionresource_rel_cmd->nas_msg) {
      ie = (Ngap_PDUSessionResourceReleaseCommandIEs_t*) calloc(
          1, sizeof(Ngap_PDUSessionResourceReleaseCommandIEs_t));
      ie->id          = Ngap_ProtocolIE_ID_id_NAS_PDU;
      ie->criticality = Ngap_Criticality_reject;
      ie->value.present =
          Ngap_PDUSessionResourceReleaseCommandIEs__value_PR_NAS_PDU;
      OCTET_STRING_fromBuf(
          &ie->value.choice.NAS_PDU,
          (char*) bdata(pdusessionresource_rel_cmd->nas_msg),
          blength(pdusessionresource_rel_cmd->nas_msg));
      ASN_SEQUENCE_ADD(&out->protocolIEs.list, ie);
    }

    ie = (Ngap_PDUSessionResourceReleaseCommandIEs_t*) calloc(
        1, sizeof(Ngap_PDUSessionResourceReleaseCommandIEs_t));
    ie->id = Ngap_ProtocolIE_ID_id_PDUSessionResourceToReleaseListRelCmd;
    ie->criticality = Ngap_Criticality_reject;
    ie->value.present =
        Ngap_PDUSessionResourceReleaseCommandIEs__value_PR_PDUSessionResourceToReleaseListRelCmd;
    ASN_SEQUENCE_ADD(&out->protocolIEs.list, ie);

    for (int i = 0;
         i <
         pdusessionresource_rel_cmd->pduSessionResourceToRelReqList.no_of_items;
         i++) {
      Ngap_PDUSessionResourceToReleaseItemRelCmd_t* RelItem =
          calloc(1, sizeof(Ngap_PDUSessionResourceToReleaseItemRelCmd_t));
      RelItem->pDUSessionID =
          pdusessionresource_rel_cmd->pduSessionResourceToRelReqList.item[i]
              .Pdu_Session_ID;

      Ngap_PDUSessionResourceReleaseCommandTransfer_t*
          PDUSessionResourceReleaseCommandTransferIEs = calloc(
              1, sizeof(Ngap_PDUSessionResourceReleaseCommandTransfer_t));

      PDUSessionResourceReleaseCommandTransferIEs->cause.present =
          Ngap_Cause_PR_nas;
      PDUSessionResourceReleaseCommandTransferIEs->cause.choice.nas =
          Ngap_CauseNas_normal_release;

      asn_fprint(stderr, &asn_DEF_Ngap_NGAP_PDU, &pdu);

      uint32_t buffer_size = 512;
      char* buffer         = (char*) calloc(1, buffer_size);

      ssize_t encoded_size = aper_encode_to_new_buffer(
          &asn_DEF_Ngap_PDUSessionResourceReleaseCommandTransfer, NULL,
          PDUSessionResourceReleaseCommandTransferIEs, (void**) &buffer);

      RelItem->pDUSessionResourceReleaseCommandTransfer.buf =
          (uint8_t*) calloc(encoded_size, sizeof(uint8_t));
      memcpy(
          (void*) RelItem->pDUSessionResourceReleaseCommandTransfer.buf,
          (void*) buffer, encoded_size);
      RelItem->pDUSessionResourceReleaseCommandTransfer.size = encoded_size;

      ASN_SEQUENCE_ADD(
          &ie->value.choice.PDUSessionResourceToReleaseListRelCmd.list,
          RelItem);
    }

    if (ngap_amf_encode_pdu(&pdu, &buffer_p, &length) < 0) {
      OAILOG_ERROR(LOG_NGAP, "Encoding of ngap_pdusessionCommandIEs failed \n");
      OAILOG_FUNC_RETURN(LOG_NGAP, RETURNerror);
    }
    bstring b = blk2bstr(buffer_p, length);
    free(buffer_p);
    ngap_amf_itti_send_sctp_request(
        &b, ue_ref->sctp_assoc_id, ue_ref->sctp_stream_send,
        ue_ref->amf_ue_ngap_id);
  }

  OAILOG_FUNC_RETURN(LOG_NGAP, RETURNok);
}

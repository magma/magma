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

/*! \file s1ap_mme_nas_procedures.c
  \brief
  \author Sebastien ROUX, Lionel Gauthier
  \company Eurecom
  \email: lionel.gauthier@eurecom.fr
*/

#include <stdio.h>
#include <stdint.h>
#include <stdlib.h>
#include <string.h>

#include "bstrlib.h"
#include "dynamic_memory_check.h"
#include "assertions.h"
#include "hashtable.h"
#include "log.h"
#include "conversions.h"
#include "asn1_conversions.h"
#include "s1ap_mme_encoder.h"
#include "s1ap_mme.h"
#include "s1ap_mme_nas_procedures.h"
#include "s1ap_mme_itti_messaging.h"
#include "service303.h"
#include "3gpp_23.003.h"
#include "3gpp_24.007.h"
#include "3gpp_36.413.h"
#include "INTEGER.h"
#include "OCTET_STRING.h"
#include "S1ap_S1AP-PDU.h"
#include "S1ap_AllocationAndRetentionPriority.h"
#include "S1ap_E-RABLevelQoSParameters.h"
#include "S1ap_E-RABToBeSetupItemBearerSUReq.h"
#include "S1ap_E-RABToBeSetupItemCtxtSUReq.h"
#include "S1ap_EUTRAN-CGI.h"
#include "S1ap_EncryptionAlgorithms.h"
#include "S1ap_GBR-QosInformation.h"
#include "S1ap_GUMMEI.h"
#include "S1ap_IntegrityProtectionAlgorithms.h"
#include "S1ap_NAS-PDU.h"
#include "S1ap_PLMNidentity.h"
#include "S1ap_ProcedureCode.h"
#include "S1ap_S-TMSI.h"
#include "S1ap_SecurityKey.h"
#include "S1ap_TAI.h"
#include "S1ap_TransportLayerAddress.h"
#include "S1ap_UEAggregateMaximumBitrate.h"
#include "S1ap_UESecurityCapabilities.h"
#include "TrackingAreaIdentity.h"
#include "asn_SEQUENCE_OF.h"
#include "nas/securityDef.h"
#include "s1ap_state.h"
#include "S1ap_CauseMisc.h"
#include "S1ap_CauseNas.h"
#include "S1ap_CauseProtocol.h"
#include "S1ap_CauseRadioNetwork.h"
#include "S1ap_CauseTransport.h"
#include "S1ap_E-RABItem.h"
#include "s1ap_mme_handlers.h"
#include "S1ap_ProtocolIE-Field.h"
#include "s1ap_common.h"

extern bool s1ap_congestion_control_enabled;
extern long s1ap_last_msg_latency;
extern long s1ap_zmq_th;

//------------------------------------------------------------------------------
int s1ap_mme_handle_initial_ue_message(
    s1ap_state_t* state, const sctp_assoc_id_t assoc_id,
    const sctp_stream_id_t stream, S1ap_S1AP_PDU_t* pdu) {
  S1ap_InitialUEMessage_t* container = NULL;
  S1ap_InitialUEMessage_IEs_t *ie = NULL, *ie_e_tmsi = NULL, *ie_csg_id = NULL,
                              *ie_gummei = NULL, *ie_cause = NULL;
  ue_description_t* ue_ref        = NULL;
  enb_description_t* eNB_ref      = NULL;
  enb_ue_s1ap_id_t enb_ue_s1ap_id = INVALID_ENB_UE_S1AP_ID;

  OAILOG_FUNC_IN(LOG_S1AP);
  container = &pdu->choice.initiatingMessage.value.choice.InitialUEMessage;

  S1AP_FIND_PROTOCOLIE_BY_ID(
      S1ap_InitialUEMessage_IEs_t, ie, container,
      S1ap_ProtocolIE_ID_id_eNB_UE_S1AP_ID, true);

  OAILOG_INFO(
      LOG_S1AP,
      "Received S1AP INITIAL_UE_MESSAGE ENB_UE_S1AP_ID " ENB_UE_S1AP_ID_FMT
      " assoc-id:%d \n",
      (enb_ue_s1ap_id_t) ie->value.choice.ENB_UE_S1AP_ID, assoc_id);

  if (s1ap_congestion_control_enabled &&
      (s1ap_last_msg_latency > S1AP_ZMQ_LATENCY_TH)) {
    OAILOG_WARNING(
        LOG_S1AP,
        "Discarding S1AP INITIAL_UE_MESSAGE for "
        "ENB_UE_S1AP_ID: " ENB_UE_S1AP_ID_FMT " ZMQ latency: %ld",
        (enb_ue_s1ap_id_t) ie->value.choice.ENB_UE_S1AP_ID,
        s1ap_last_msg_latency);
    OAILOG_FUNC_RETURN(LOG_S1AP, RETURNerror);
  }

  if ((eNB_ref = s1ap_state_get_enb(state, assoc_id)) == NULL) {
    OAILOG_ERROR(LOG_S1AP, "Unknown eNB on assoc_id %d\n", assoc_id);
    OAILOG_FUNC_RETURN(LOG_S1AP, RETURNerror);
  }
  // eNB UE S1AP ID is limited to 24 bits
  enb_ue_s1ap_id =
      (enb_ue_s1ap_id_t)(ie->value.choice.ENB_UE_S1AP_ID & 0x00ffffff);
  OAILOG_INFO(
      LOG_S1AP,
      "New Initial UE message received with eNB UE S1AP ID: " ENB_UE_S1AP_ID_FMT
      " assoc-id :%d \n",
      enb_ue_s1ap_id, eNB_ref->sctp_assoc_id);
  ue_ref = s1ap_state_get_ue_enbid(eNB_ref->sctp_assoc_id, enb_ue_s1ap_id);

  if (ue_ref == NULL) {
    tai_t tai       = {0};
    gummei_t gummei = {
        .plmn = {0}, .mme_code = 0, .mme_gid = 0};  // initialized after
    s_tmsi_t s_tmsi = {.mme_code = 0, .m_tmsi = INVALID_M_TMSI};
    ecgi_t ecgi     = {.plmn = {0}, .cell_identity = {0}};
    csg_id_t csg_id = 0;

    /*
     * This UE eNB Id has currently no known s1 association.
     * * * * Create new UE context by associating new mme_ue_s1ap_id.
     * * * * Update eNB UE list.
     * * * * Forward message to NAS.
     */
    if ((ue_ref = s1ap_new_ue(state, assoc_id, enb_ue_s1ap_id)) == NULL) {
      // If we failed to allocate a new UE return -1
      OAILOG_ERROR(
          LOG_S1AP,
          "Initial UE Message- Failed to allocate S1AP UE Context, "
          "eNB UE S1AP ID:" ENB_UE_S1AP_ID_FMT "\n",
          enb_ue_s1ap_id);
      OAILOG_FUNC_RETURN(LOG_S1AP, RETURNerror);
    }

    OAILOG_DEBUG(LOG_S1AP, "Creating new UE Ref on S1ap");

    ue_ref->s1_ue_state = S1AP_UE_WAITING_CSR;

    ue_ref->enb_ue_s1ap_id = enb_ue_s1ap_id;
    // Will be allocated by NAS
    ue_ref->mme_ue_s1ap_id = INVALID_MME_UE_S1AP_ID;

    ue_ref->s1ap_ue_context_rel_timer.id  = S1AP_TIMER_INACTIVE_ID;
    ue_ref->s1ap_ue_context_rel_timer.sec = S1AP_UE_CONTEXT_REL_COMP_TIMER;

    // On which stream we received the message
    ue_ref->sctp_stream_recv = stream;
    ue_ref->sctp_stream_send = eNB_ref->next_sctp_stream;

    /*
     * Increment the sctp stream for the eNB association.
     * If the next sctp stream is >= instream negociated between eNB and MME,
     * wrap to first stream.
     * TODO: search for the first available stream instead.
     */

    /*
     * TODO task#15456359.
     * Below logic seems to be incorrect , revisit it.
     */
    eNB_ref->next_sctp_stream += 1;
    if (eNB_ref->next_sctp_stream >= eNB_ref->instreams) {
      eNB_ref->next_sctp_stream = 1;
    }
    s1ap_dump_enb(eNB_ref);
    // TAI mandatory IE
    S1AP_FIND_PROTOCOLIE_BY_ID(
        S1ap_InitialUEMessage_IEs_t, ie, container, S1ap_ProtocolIE_ID_id_TAI,
        true);
    OCTET_STRING_TO_TAC(&ie->value.choice.TAI.tAC, tai.tac);
    if (!(ie->value.choice.TAI.pLMNidentity.size == 3)) {
      OAILOG_ERROR(LOG_S1AP, "Incorrect PLMN size \n");
      return RETURNerror;
    }
    TBCD_TO_PLMN_T(&ie->value.choice.TAI.pLMNidentity, &tai.plmn);

    // CGI mandatory IE
    S1AP_FIND_PROTOCOLIE_BY_ID(
        S1ap_InitialUEMessage_IEs_t, ie, container,
        S1ap_ProtocolIE_ID_id_EUTRAN_CGI, true);
    if (!(ie->value.choice.EUTRAN_CGI.pLMNidentity.size == 3)) {
      OAILOG_ERROR(LOG_S1AP, "Incorrect PLMN size \n");
      return RETURNerror;
    }
    TBCD_TO_PLMN_T(&ie->value.choice.EUTRAN_CGI.pLMNidentity, &ecgi.plmn);
    BIT_STRING_TO_CELL_IDENTITY(
        &ie->value.choice.EUTRAN_CGI.cell_ID, ecgi.cell_identity);

    /** Set the ENB Id. */
    ecgi.cell_identity.enb_id = eNB_ref->enb_id;

    S1AP_FIND_PROTOCOLIE_BY_ID(
        S1ap_InitialUEMessage_IEs_t, ie_e_tmsi, container,
        S1ap_ProtocolIE_ID_id_S_TMSI, false);
    if (ie_e_tmsi) {
      OCTET_STRING_TO_MME_CODE(
          &ie_e_tmsi->value.choice.S_TMSI.mMEC, s_tmsi.mme_code);
      OCTET_STRING_TO_M_TMSI(
          &ie_e_tmsi->value.choice.S_TMSI.m_TMSI, s_tmsi.m_tmsi);
    }

    S1AP_FIND_PROTOCOLIE_BY_ID(
        S1ap_InitialUEMessage_IEs_t, ie_csg_id, container,
        S1ap_ProtocolIE_ID_id_CSG_Id, false);
    if (ie_csg_id) {
      csg_id = BIT_STRING_to_uint32(&ie_csg_id->value.choice.CSG_Id);
    }

    S1AP_FIND_PROTOCOLIE_BY_ID(
        S1ap_InitialUEMessage_IEs_t, ie_gummei, container,
        S1ap_ProtocolIE_ID_id_GUMMEI_ID, false);
    memset(&gummei, 0, sizeof(gummei));
    if (ie_gummei) {
      OCTET_STRING_TO_MME_GID(
          &ie_gummei->value.choice.GUMMEI.mME_Group_ID, gummei.mme_gid);
      OCTET_STRING_TO_MME_CODE(
          &ie_gummei->value.choice.GUMMEI.mME_Code, gummei.mme_code);
    }
    /*
     * We received the first NAS transport message: initial UE message.
     * * * * Send a NAS ESTAeNBBLISH IND to NAS layer
     */
    S1AP_FIND_PROTOCOLIE_BY_ID(
        S1ap_InitialUEMessage_IEs_t, ie, container,
        S1ap_ProtocolIE_ID_id_NAS_PDU, true);
    S1AP_FIND_PROTOCOLIE_BY_ID(
        S1ap_InitialUEMessage_IEs_t, ie_cause, container,
        S1ap_ProtocolIE_ID_id_RRC_Establishment_Cause, true);
    s1ap_mme_itti_s1ap_initial_ue_message(
        assoc_id, eNB_ref->enb_id, ue_ref->enb_ue_s1ap_id,
        ie->value.choice.NAS_PDU.buf, ie->value.choice.NAS_PDU.size, &tai,
        &ecgi, ie_cause->value.choice.RRC_Establishment_Cause,
        ie_e_tmsi ? &s_tmsi : NULL, ie_csg_id ? &csg_id : NULL,
        ie_gummei ? &gummei : NULL,
        NULL,  // CELL ACCESS MODE
        NULL,  // GW Transport Layer Address
        NULL   // Relay Node Indicator
    );

  } else {
    imsi64_t imsi64                = INVALID_IMSI64;
    s1ap_imsi_map_t* s1ap_imsi_map = get_s1ap_imsi_map();
    hashtable_uint64_ts_get(
        s1ap_imsi_map->mme_ue_id_imsi_htbl,
        (const hash_key_t) ue_ref->mme_ue_s1ap_id, &imsi64);

    OAILOG_ERROR_UE(
        LOG_S1AP, imsi64,
        "Initial UE Message- Duplicate ENB_UE_S1AP_ID. Ignoring the "
        "message, eNB UE S1AP ID:" ENB_UE_S1AP_ID_FMT
        "\n, mme UE s1ap ID: " MME_UE_S1AP_ID_FMT "UE state %u",
        enb_ue_s1ap_id, ue_ref->mme_ue_s1ap_id, ue_ref->s1_ue_state);
  }

  OAILOG_FUNC_RETURN(LOG_S1AP, RETURNok);
}

//------------------------------------------------------------------------------
int s1ap_mme_handle_uplink_nas_transport(
    s1ap_state_t* state, const sctp_assoc_id_t assoc_id,
    __attribute__((unused)) const sctp_stream_id_t stream,
    S1ap_S1AP_PDU_t* pdu) {
  S1ap_UplinkNASTransport_t* container = NULL;
  S1ap_UplinkNASTransport_IEs_t *ie, *ie_nas_pdu = NULL;
  ue_description_t* ue_ref        = NULL;
  enb_description_t* enb_ref      = NULL;
  tai_t tai                       = {0};
  ecgi_t ecgi                     = {.plmn = {0}, .cell_identity = {0}};
  mme_ue_s1ap_id_t mme_ue_s1ap_id = INVALID_MME_UE_S1AP_ID;
  enb_ue_s1ap_id_t enb_ue_s1ap_id = INVALID_ENB_UE_S1AP_ID;

  OAILOG_FUNC_IN(LOG_S1AP);
  container = &pdu->choice.initiatingMessage.value.choice.UplinkNASTransport;

  S1AP_FIND_PROTOCOLIE_BY_ID(
      S1ap_UplinkNASTransport_IEs_t, ie, container,
      S1ap_ProtocolIE_ID_id_eNB_UE_S1AP_ID, true);
  enb_ue_s1ap_id = (enb_ue_s1ap_id_t) ie->value.choice.ENB_UE_S1AP_ID;

  S1AP_FIND_PROTOCOLIE_BY_ID(
      S1ap_UplinkNASTransport_IEs_t, ie, container,
      S1ap_ProtocolIE_ID_id_MME_UE_S1AP_ID, true);
  mme_ue_s1ap_id = (mme_ue_s1ap_id_t) ie->value.choice.MME_UE_S1AP_ID;

  enb_ref = s1ap_state_get_enb(state, assoc_id);
  if (enb_ref == NULL) {
    OAILOG_ERROR(
        LOG_S1AP, "No eNB reference exists for association id %d\n", assoc_id);
    return RETURNerror;
  }

  if (mme_ue_s1ap_id == INVALID_MME_UE_S1AP_ID) {
    OAILOG_WARNING(
        LOG_S1AP,
        "Received S1AP UPLINK_NAS_TRANSPORT message MME_UE_S1AP_ID unknown\n");

    if (!(ue_ref = s1ap_state_get_ue_enbid(
              enb_ref->sctp_assoc_id, enb_ue_s1ap_id))) {
      OAILOG_WARNING(
          LOG_S1AP,
          "Received S1AP UPLINK_NAS_TRANSPORT No UE is attached to this "
          "enb_ue_s1ap_id: " ENB_UE_S1AP_ID_FMT "\n",
          enb_ue_s1ap_id);
      OAILOG_FUNC_RETURN(LOG_S1AP, RETURNerror);
    }
  } else {
    OAILOG_INFO(
        LOG_S1AP,
        "Received S1AP UPLINK_NAS_TRANSPORT message "
        "MME_UE_S1AP_ID " MME_UE_S1AP_ID_FMT "\n",
        mme_ue_s1ap_id);

    if (!(ue_ref = s1ap_state_get_ue_mmeid(mme_ue_s1ap_id))) {
      OAILOG_WARNING(
          LOG_S1AP,
          "Received S1AP UPLINK_NAS_TRANSPORT No UE is attached to this "
          "mme_ue_s1ap_id: " MME_UE_S1AP_ID_FMT "\n",
          mme_ue_s1ap_id);
      imsi64_t imsi64                = INVALID_IMSI64;
      s1ap_imsi_map_t* s1ap_imsi_map = get_s1ap_imsi_map();
      hashtable_uint64_ts_get(
          s1ap_imsi_map->mme_ue_id_imsi_htbl, (const hash_key_t) mme_ue_s1ap_id,
          &imsi64);

      s1ap_mme_generate_ue_context_release_command(
          state, ue_ref, S1AP_INVALID_MME_UE_S1AP_ID, imsi64, assoc_id, stream,
          mme_ue_s1ap_id, enb_ue_s1ap_id);
      /* If UE context doesn't exist for received mme_ue_s1ap_id
       * remove the corresponding enb_ue_s1ap_id_key entry in mme_app
       */
      s1ap_mme_remove_stale_ue_context(enb_ue_s1ap_id, enb_ref->enb_id);
      OAILOG_FUNC_RETURN(LOG_S1AP, RETURNerror);
    }
  }

  if (S1AP_UE_CONNECTED != ue_ref->s1_ue_state) {
    OAILOG_WARNING(
        LOG_S1AP,
        "Received S1AP UPLINK_NAS_TRANSPORT while UE in state != "
        "S1AP_UE_CONNECTED\n");

    OAILOG_FUNC_RETURN(LOG_S1AP, RETURNerror);
  }

  S1AP_FIND_PROTOCOLIE_BY_ID(
      S1ap_UplinkNASTransport_IEs_t, ie_nas_pdu, container,
      S1ap_ProtocolIE_ID_id_NAS_PDU, true);
  // TAI mandatory IE
  S1AP_FIND_PROTOCOLIE_BY_ID(
      S1ap_UplinkNASTransport_IEs_t, ie, container, S1ap_ProtocolIE_ID_id_TAI,
      true);
  OCTET_STRING_TO_TAC(&ie->value.choice.TAI.tAC, tai.tac);
  if (!(ie->value.choice.TAI.pLMNidentity.size == 3)) {
    OAILOG_ERROR(LOG_S1AP, "Incorrect PLMN size \n");
    return RETURNerror;
  }
  TBCD_TO_PLMN_T(&ie->value.choice.TAI.pLMNidentity, &tai.plmn);

  // CGI mandatory IE
  S1AP_FIND_PROTOCOLIE_BY_ID(
      S1ap_UplinkNASTransport_IEs_t, ie, container,
      S1ap_ProtocolIE_ID_id_EUTRAN_CGI, true);
  if (!(ie->value.choice.EUTRAN_CGI.pLMNidentity.size == 3)) {
    OAILOG_ERROR(LOG_S1AP, "Incorrect PLMN size \n");
    return RETURNerror;
  }
  TBCD_TO_PLMN_T(&ie->value.choice.EUTRAN_CGI.pLMNidentity, &ecgi.plmn);
  BIT_STRING_TO_CELL_IDENTITY(
      &ie->value.choice.EUTRAN_CGI.cell_ID, ecgi.cell_identity);
  // set the eNB ID
  ecgi.cell_identity.enb_id = enb_ref->enb_id;
  // TODO optional GW Transport Layer Address

  bstring b = blk2bstr(
      ie_nas_pdu->value.choice.NAS_PDU.buf,
      ie_nas_pdu->value.choice.NAS_PDU.size);
  s1ap_mme_itti_nas_uplink_ind(mme_ue_s1ap_id, &b, &tai, &ecgi);
  OAILOG_FUNC_RETURN(LOG_S1AP, RETURNok);
}

//------------------------------------------------------------------------------
int s1ap_mme_handle_nas_non_delivery(
    s1ap_state_t* state, __attribute__((unused)) sctp_assoc_id_t assoc_id,
    sctp_stream_id_t stream, S1ap_S1AP_PDU_t* pdu) {
  S1ap_NASNonDeliveryIndication_t* container;
  S1ap_NASNonDeliveryIndication_IEs_t *ie = NULL, *ie_nas_pdu = NULL;
  ue_description_t* ue_ref        = NULL;
  imsi64_t imsi64                 = INVALID_IMSI64;
  mme_ue_s1ap_id_t mme_ue_s1ap_id = INVALID_MME_UE_S1AP_ID;
  enb_ue_s1ap_id_t enb_ue_s1ap_id = INVALID_ENB_UE_S1AP_ID;

  OAILOG_FUNC_IN(LOG_S1AP);
  increment_counter("nas_non_delivery_indication_received", 1, NO_LABELS);

  container =
      &pdu->choice.initiatingMessage.value.choice.NASNonDeliveryIndication;
  /*
   * UE associated signalling on stream == 0 is not valid.
   */
  if (stream == 0) {
    OAILOG_NOTICE(
        LOG_S1AP,
        "Received S1AP NAS_NON_DELIVERY_INDICATION message on invalid sctp "
        "stream 0\n");
    OAILOG_FUNC_RETURN(LOG_S1AP, RETURNerror);
  }
  S1AP_FIND_PROTOCOLIE_BY_ID(
      S1ap_NASNonDeliveryIndication_IEs_t, ie, container,
      S1ap_ProtocolIE_ID_id_MME_UE_S1AP_ID, true);
  mme_ue_s1ap_id = ie->value.choice.MME_UE_S1AP_ID;

  S1AP_FIND_PROTOCOLIE_BY_ID(
      S1ap_NASNonDeliveryIndication_IEs_t, ie, container,
      S1ap_ProtocolIE_ID_id_eNB_UE_S1AP_ID, true);
  enb_ue_s1ap_id = ie->value.choice.ENB_UE_S1AP_ID;

  S1AP_FIND_PROTOCOLIE_BY_ID(
      S1ap_NASNonDeliveryIndication_IEs_t, ie_nas_pdu, container,
      S1ap_ProtocolIE_ID_id_NAS_PDU, true);

  S1AP_FIND_PROTOCOLIE_BY_ID(
      S1ap_NASNonDeliveryIndication_IEs_t, ie, container,
      S1ap_ProtocolIE_ID_id_Cause, true);

  /*
   * UE associated signalling on stream == 0 is not valid.
   */
  if (stream == 0) {
    OAILOG_NOTICE(
        LOG_S1AP,
        "Received S1AP NAS_NON_DELIVERY_INDICATION message on "
        "invalid sctp stream 0\n");
    OAILOG_FUNC_RETURN(LOG_S1AP, RETURNerror);
  }

  OAILOG_NOTICE(
      LOG_S1AP,
      "Received S1AP NAS_NON_DELIVERY_INDICATION message "
      "MME_UE_S1AP_ID " MME_UE_S1AP_ID_FMT " enb_ue_s1ap_id " ENB_UE_S1AP_ID_FMT
      "\n",
      mme_ue_s1ap_id, enb_ue_s1ap_id);

  if ((ue_ref = s1ap_state_get_ue_mmeid(mme_ue_s1ap_id)) == NULL) {
    OAILOG_DEBUG(
        LOG_S1AP,
        "No UE is attached to this mme UE s1ap id: " MME_UE_S1AP_ID_FMT "\n",
        mme_ue_s1ap_id);
    OAILOG_FUNC_RETURN(LOG_S1AP, RETURNerror);
  }

  s1ap_imsi_map_t* imsi_map = get_s1ap_imsi_map();
  hashtable_uint64_ts_get(
      imsi_map->mme_ue_id_imsi_htbl, mme_ue_s1ap_id, &imsi64);

  if (ue_ref->s1_ue_state != S1AP_UE_CONNECTED) {
    OAILOG_DEBUG_UE(
        LOG_S1AP, imsi64,
        "Received S1AP NAS_NON_DELIVERY_INDICATION while UE in state != "
        "S1AP_UE_CONNECTED\n");
    OAILOG_FUNC_RETURN(LOG_S1AP, RETURNerror);
  }

  // TODO: forward NAS PDU to NAS
  s1ap_mme_itti_nas_non_delivery_ind(
      ie->value.choice.MME_UE_S1AP_ID, ie_nas_pdu->value.choice.NAS_PDU.buf,
      ie_nas_pdu->value.choice.NAS_PDU.size, &ie->value.choice.Cause, imsi64);
  OAILOG_FUNC_RETURN(LOG_S1AP, RETURNok);
}

//------------------------------------------------------------------------------
int s1ap_generate_downlink_nas_transport(
    s1ap_state_t* state, const enb_ue_s1ap_id_t enb_ue_s1ap_id,
    const mme_ue_s1ap_id_t ue_id, STOLEN_REF bstring* payload,
    const imsi64_t imsi64, bool* is_state_same) {
  ue_description_t* ue_ref = NULL;
  uint8_t* buffer_p        = NULL;
  uint32_t length          = 0;
  void* id                 = NULL;
  uint8_t err              = 0;

  OAILOG_FUNC_IN(LOG_S1AP);

  // Try to retrieve SCTP association id using mme_ue_s1ap_id
  if (HASH_TABLE_OK ==
      hashtable_ts_get(
          &state->mmeid2associd, (const hash_key_t) ue_id, (void**) &id)) {
    sctp_assoc_id_t sctp_assoc_id = (sctp_assoc_id_t)(uintptr_t) id;
    enb_description_t* enb_ref    = s1ap_state_get_enb(state, sctp_assoc_id);
    if (enb_ref) {
      ue_ref = s1ap_state_get_ue_enbid(enb_ref->sctp_assoc_id, enb_ue_s1ap_id);
    } else {
      OAILOG_ERROR(
          LOG_S1AP, "No eNB for SCTP association id %d \n", sctp_assoc_id);
      OAILOG_FUNC_RETURN(LOG_S1AP, RETURNerror);
    }
  }
  // TODO remove soon:
  if (!ue_ref) {
    ue_ref = s1ap_state_get_ue_mmeid(ue_id);
  }
  // finally!
  if (!ue_ref) {
    /*
     * If the UE-associated logical S1-connection is not established,
     * * * * the MME shall allocate a unique MME UE S1AP ID to be used for the
     * UE.
     */
    OAILOG_WARNING(
        LOG_S1AP,
        "Unknown UE MME ID " MME_UE_S1AP_ID_FMT
        ", This case is not handled right now\n",
        ue_id);
    OAILOG_FUNC_RETURN(LOG_S1AP, RETURNerror);
  } else {
    /*
     * We have found the UE in the list.
     * * * * Create new IE list message and encode it.
     */
    s1ap_imsi_map_t* imsi_map = get_s1ap_imsi_map();
    if (hashtable_uint64_ts_insert(
            imsi_map->mme_ue_id_imsi_htbl, (const hash_key_t) ue_id, imsi64) ==
        HASH_TABLE_SAME_KEY_VALUE_EXISTS) {
      *is_state_same = true;
    }

    S1ap_DownlinkNASTransport_IEs_t* ie = NULL;
    S1ap_DownlinkNASTransport_t* out    = NULL;
    S1ap_S1AP_PDU_t pdu                 = {0};

    pdu.present = S1ap_S1AP_PDU_PR_initiatingMessage;
    pdu.choice.initiatingMessage.procedureCode =
        S1ap_ProcedureCode_id_downlinkNASTransport;
    pdu.choice.initiatingMessage.criticality = S1ap_Criticality_ignore;
    pdu.choice.initiatingMessage.value.present =
        S1ap_InitiatingMessage__value_PR_DownlinkNASTransport;

    out = &pdu.choice.initiatingMessage.value.choice.DownlinkNASTransport;

    if (ue_ref->s1_ue_state == S1AP_UE_WAITING_CRR) {
      OAILOG_ERROR_UE(
          LOG_S1AP, imsi64,
          "Already triggered UE Context Release Command and UE is"
          "in S1AP_UE_WAITING_CRR, so dropping the DownlinkNASTransport \n");
      OAILOG_FUNC_RETURN(LOG_S1AP, RETURNerror);
    } else {
      ue_ref->s1_ue_state = S1AP_UE_CONNECTED;
    }
    /*
     * Setting UE informations with the ones found in ue_ref
     */
    ie = (S1ap_DownlinkNASTransport_IEs_t*) calloc(
        1, sizeof(S1ap_DownlinkNASTransport_IEs_t));
    ie->id            = S1ap_ProtocolIE_ID_id_MME_UE_S1AP_ID;
    ie->criticality   = S1ap_Criticality_reject;
    ie->value.present = S1ap_DownlinkNASTransport_IEs__value_PR_MME_UE_S1AP_ID;
    ie->value.choice.MME_UE_S1AP_ID = ue_ref->mme_ue_s1ap_id;
    ASN_SEQUENCE_ADD(&out->protocolIEs.list, ie);

    /* mandatory */
    ie = (S1ap_DownlinkNASTransport_IEs_t*) calloc(
        1, sizeof(S1ap_DownlinkNASTransport_IEs_t));
    ie->id            = S1ap_ProtocolIE_ID_id_eNB_UE_S1AP_ID;
    ie->criticality   = S1ap_Criticality_reject;
    ie->value.present = S1ap_DownlinkNASTransport_IEs__value_PR_ENB_UE_S1AP_ID;
    ie->value.choice.ENB_UE_S1AP_ID = ue_ref->enb_ue_s1ap_id;
    ASN_SEQUENCE_ADD(&out->protocolIEs.list, ie);
    /* mandatory */
    ie = (S1ap_DownlinkNASTransport_IEs_t*) calloc(
        1, sizeof(S1ap_DownlinkNASTransport_IEs_t));
    ie->id            = S1ap_ProtocolIE_ID_id_NAS_PDU;
    ie->criticality   = S1ap_Criticality_reject;
    ie->value.present = S1ap_DownlinkNASTransport_IEs__value_PR_NAS_PDU;
    /*eNB
     * Fill in the NAS pdu
     */
    OCTET_STRING_fromBuf(
        &ie->value.choice.NAS_PDU, (char*) bdata(*payload), blength(*payload));
    ASN_SEQUENCE_ADD(&out->protocolIEs.list, ie);

    if (s1ap_mme_encode_pdu(&pdu, &buffer_p, &length) < 0) {
      err = 1;
    }
    ASN_STRUCT_FREE_CONTENTS_ONLY(asn_DEF_S1ap_DownlinkNASTransport, out);
    if (err) {
      OAILOG_FUNC_RETURN(LOG_S1AP, RETURNerror);
    }

    OAILOG_NOTICE_UE(
        LOG_S1AP, imsi64,
        "Send S1AP DOWNLINK_NAS_TRANSPORT message ue_id = " MME_UE_S1AP_ID_FMT
        " MME_UE_S1AP_ID = " MME_UE_S1AP_ID_FMT
        " eNB_UE_S1AP_ID = " ENB_UE_S1AP_ID_FMT "\n",
        ue_id, ue_ref->mme_ue_s1ap_id, enb_ue_s1ap_id);
    bstring b = blk2bstr(buffer_p, length);
    free(buffer_p);
    s1ap_mme_itti_send_sctp_request(
        &b, ue_ref->sctp_assoc_id, ue_ref->sctp_stream_send,
        ue_ref->mme_ue_s1ap_id);
  }

  OAILOG_FUNC_RETURN(LOG_S1AP, RETURNok);
}

//------------------------------------------------------------------------------
int s1ap_generate_s1ap_e_rab_setup_req(
    s1ap_state_t* state, itti_s1ap_e_rab_setup_req_t* const e_rab_setup_req) {
  OAILOG_FUNC_IN(LOG_S1AP);
  ue_description_t* ue_ref              = NULL;
  uint8_t* buffer_p                     = NULL;
  uint32_t length                       = 0;
  void* id                              = NULL;
  const enb_ue_s1ap_id_t enb_ue_s1ap_id = e_rab_setup_req->enb_ue_s1ap_id;
  const mme_ue_s1ap_id_t ue_id          = e_rab_setup_req->mme_ue_s1ap_id;

  hashtable_ts_get(
      &state->mmeid2associd, (const hash_key_t) ue_id, (void**) &id);
  if (id) {
    sctp_assoc_id_t sctp_assoc_id = (sctp_assoc_id_t)(uintptr_t) id;
    enb_description_t* enb_ref    = s1ap_state_get_enb(state, sctp_assoc_id);
    if (enb_ref) {
      ue_ref = s1ap_state_get_ue_enbid(enb_ref->sctp_assoc_id, enb_ue_s1ap_id);
    }
  }
  // TODO remove soon:
  if (!ue_ref) {
    ue_ref = s1ap_state_get_ue_mmeid(ue_id);
  }
  // finally!
  if (!ue_ref) {
    /*
     * If the UE-associated logical S1-connection is not established,
     * * * * the MME shall allocate a unique MME UE S1AP ID to be used for the
     * UE.
     */
    OAILOG_ERROR(
        LOG_S1AP,
        "Unknown UE MME ID " MME_UE_S1AP_ID_FMT
        ", This case is not handled right now\n",
        ue_id);
    OAILOG_FUNC_RETURN(LOG_S1AP, RETURNerror);
  } else {
    /*
     * We have found the UE in the list.
     * Create new IE list message and encode it.
     */
    S1ap_S1AP_PDU_t pdu             = {0};
    S1ap_E_RABSetupRequest_t* out   = NULL;
    S1ap_E_RABSetupRequestIEs_t* ie = NULL;
    pdu.choice.initiatingMessage.procedureCode =
        S1ap_ProcedureCode_id_E_RABSetup;
    pdu.choice.initiatingMessage.criticality = S1ap_Criticality_reject;
    pdu.present = S1ap_S1AP_PDU_PR_initiatingMessage;
    pdu.choice.initiatingMessage.value.present =
        S1ap_InitiatingMessage__value_PR_E_RABSetupRequest;
    out = &pdu.choice.initiatingMessage.value.choice.E_RABSetupRequest;
    ue_ref->s1_ue_state = S1AP_UE_CONNECTED;
    /*
     * Setting UE information with the ones found in ue_ref
     */
    ie = (S1ap_E_RABSetupRequestIEs_t*) calloc(
        1, sizeof(S1ap_E_RABSetupRequestIEs_t));
    ie->id            = S1ap_ProtocolIE_ID_id_MME_UE_S1AP_ID;
    ie->criticality   = S1ap_Criticality_reject;
    ie->value.present = S1ap_E_RABSetupRequestIEs__value_PR_MME_UE_S1AP_ID;
    ie->value.choice.MME_UE_S1AP_ID = ue_ref->mme_ue_s1ap_id;
    ASN_SEQUENCE_ADD(&out->protocolIEs.list, ie);

    /* mandatory */
    ie = (S1ap_E_RABSetupRequestIEs_t*) calloc(
        1, sizeof(S1ap_E_RABSetupRequestIEs_t));
    ie->id            = S1ap_ProtocolIE_ID_id_eNB_UE_S1AP_ID;
    ie->criticality   = S1ap_Criticality_reject;
    ie->value.present = S1ap_E_RABSetupRequestIEs__value_PR_ENB_UE_S1AP_ID;
    ie->value.choice.ENB_UE_S1AP_ID = ue_ref->enb_ue_s1ap_id;
    ASN_SEQUENCE_ADD(&out->protocolIEs.list, ie);
    /*eNB
     * Fill in the NAS pdu
     */
    if (e_rab_setup_req->ue_aggregate_maximum_bit_rate_present) {
      ie = (S1ap_E_RABSetupRequestIEs_t*) calloc(
          1, sizeof(S1ap_E_RABSetupRequestIEs_t));
      ie->id          = S1ap_ProtocolIE_ID_id_uEaggregateMaximumBitrate;
      ie->criticality = S1ap_Criticality_reject;
      ie->value.present =
          S1ap_E_RABSetupRequestIEs__value_PR_UEAggregateMaximumBitrate;
      asn_uint642INTEGER(
          &ie->value.choice.UEAggregateMaximumBitrate
               .uEaggregateMaximumBitRateDL,
          e_rab_setup_req->ue_aggregate_maximum_bit_rate.dl);
      asn_uint642INTEGER(
          &ie->value.choice.UEAggregateMaximumBitrate
               .uEaggregateMaximumBitRateUL,
          e_rab_setup_req->ue_aggregate_maximum_bit_rate.ul);
      ASN_SEQUENCE_ADD(&out->protocolIEs.list, ie);
    }

    /* mandatory */
    ie = (S1ap_E_RABSetupRequestIEs_t*) calloc(
        1, sizeof(S1ap_E_RABSetupRequestIEs_t));
    ie->id          = S1ap_ProtocolIE_ID_id_E_RABToBeSetupListBearerSUReq;
    ie->criticality = S1ap_Criticality_reject;
    ie->value.present =
        S1ap_E_RABSetupRequestIEs__value_PR_E_RABToBeSetupListBearerSUReq;
    ASN_SEQUENCE_ADD(&out->protocolIEs.list, ie);

    S1ap_E_RABToBeSetupListBearerSUReq_t* e_rabtobesetuplistbearersureq =
        &ie->value.choice.E_RABToBeSetupListBearerSUReq;

    for (int i = 0; i < e_rab_setup_req->e_rab_to_be_setup_list.no_of_items;
         i++) {
      S1ap_E_RABToBeSetupItemBearerSUReqIEs_t* s1ap_e_rab_to_be_setup_item_ies =
          calloc(1, sizeof(S1ap_E_RABToBeSetupItemBearerSUReqIEs_t));

      s1ap_e_rab_to_be_setup_item_ies->id =
          S1ap_ProtocolIE_ID_id_E_RABToBeSetupItemBearerSUReq;
      s1ap_e_rab_to_be_setup_item_ies->criticality = S1ap_Criticality_reject;
      s1ap_e_rab_to_be_setup_item_ies->value.present =
          S1ap_E_RABToBeSetupItemBearerSUReqIEs__value_PR_E_RABToBeSetupItemBearerSUReq;

      S1ap_E_RABToBeSetupItemBearerSUReq_t* e_rab_to_be_set_up_item =
          &s1ap_e_rab_to_be_setup_item_ies->value.choice
               .E_RABToBeSetupItemBearerSUReq;

      e_rab_to_be_set_up_item->e_RAB_ID =
          e_rab_setup_req->e_rab_to_be_setup_list.item[i].e_rab_id;
      e_rab_to_be_set_up_item->e_RABlevelQoSParameters.qCI =
          e_rab_setup_req->e_rab_to_be_setup_list.item[i]
              .e_rab_level_qos_parameters.qci;

      e_rab_to_be_set_up_item->e_RABlevelQoSParameters
          .allocationRetentionPriority.priorityLevel =
          e_rab_setup_req->e_rab_to_be_setup_list.item[i]
              .e_rab_level_qos_parameters.allocation_and_retention_priority
              .priority_level;

      e_rab_to_be_set_up_item->e_RABlevelQoSParameters
          .allocationRetentionPriority.pre_emptionCapability =
          e_rab_setup_req->e_rab_to_be_setup_list.item[i]
              .e_rab_level_qos_parameters.allocation_and_retention_priority
              .pre_emption_capability;

      e_rab_to_be_set_up_item->e_RABlevelQoSParameters
          .allocationRetentionPriority.pre_emptionVulnerability =
          e_rab_setup_req->e_rab_to_be_setup_list.item[i]
              .e_rab_level_qos_parameters.allocation_and_retention_priority
              .pre_emption_vulnerability;
      /* OPTIONAL */
      gbr_qos_information_t* gbr_qos_information =
          &e_rab_setup_req->e_rab_to_be_setup_list.item[i]
               .e_rab_level_qos_parameters.gbr_qos_information;
      if ((gbr_qos_information->e_rab_maximum_bit_rate_downlink) ||
          (gbr_qos_information->e_rab_maximum_bit_rate_uplink) ||
          (gbr_qos_information->e_rab_guaranteed_bit_rate_downlink) ||
          (gbr_qos_information->e_rab_guaranteed_bit_rate_uplink)) {
        OAILOG_NOTICE(
            LOG_S1AP,
            "Encoding of e_RABlevelQoSParameters.gbrQosInformation\n");

        e_rab_to_be_set_up_item->e_RABlevelQoSParameters.gbrQosInformation =
            calloc(1, sizeof(struct S1ap_GBR_QosInformation));

        if (e_rab_to_be_set_up_item->e_RABlevelQoSParameters
                .gbrQosInformation) {
          asn_uint642INTEGER(
              &e_rab_to_be_set_up_item->e_RABlevelQoSParameters
                   .gbrQosInformation->e_RAB_MaximumBitrateDL,
              gbr_qos_information->e_rab_maximum_bit_rate_downlink);

          asn_uint642INTEGER(
              &e_rab_to_be_set_up_item->e_RABlevelQoSParameters
                   .gbrQosInformation->e_RAB_MaximumBitrateUL,
              gbr_qos_information->e_rab_maximum_bit_rate_uplink);

          asn_uint642INTEGER(
              &e_rab_to_be_set_up_item->e_RABlevelQoSParameters
                   .gbrQosInformation->e_RAB_GuaranteedBitrateDL,
              gbr_qos_information->e_rab_guaranteed_bit_rate_downlink);

          asn_uint642INTEGER(
              &e_rab_to_be_set_up_item->e_RABlevelQoSParameters
                   .gbrQosInformation->e_RAB_GuaranteedBitrateUL,
              gbr_qos_information->e_rab_guaranteed_bit_rate_uplink);
        }
      } else {
        OAILOG_NOTICE(
            LOG_S1AP,
            "NOT Encoding of e_RABlevelQoSParameters.gbrQosInformation\n");
      }
      INT32_TO_OCTET_STRING(
          e_rab_setup_req->e_rab_to_be_setup_list.item[i].gtp_teid,
          &e_rab_to_be_set_up_item->gTP_TEID);

      e_rab_to_be_set_up_item->transportLayerAddress.buf = calloc(
          blength(e_rab_setup_req->e_rab_to_be_setup_list.item[i]
                      .transport_layer_address),
          sizeof(uint8_t));
      memcpy(
          e_rab_to_be_set_up_item->transportLayerAddress.buf,
          e_rab_setup_req->e_rab_to_be_setup_list.item[i]
              .transport_layer_address->data,
          blength(e_rab_setup_req->e_rab_to_be_setup_list.item[i]
                      .transport_layer_address));

      e_rab_to_be_set_up_item->transportLayerAddress.size =
          blength(e_rab_setup_req->e_rab_to_be_setup_list.item[i]
                      .transport_layer_address);
      e_rab_to_be_set_up_item->transportLayerAddress.bits_unused = 0;

      OCTET_STRING_fromBuf(
          &e_rab_to_be_set_up_item->nAS_PDU,
          (char*) bdata(
              e_rab_setup_req->e_rab_to_be_setup_list.item[i].nas_pdu),
          blength(e_rab_setup_req->e_rab_to_be_setup_list.item[i].nas_pdu));

      ASN_SEQUENCE_ADD(
          &e_rabtobesetuplistbearersureq->list,
          s1ap_e_rab_to_be_setup_item_ies);
    }

    if (s1ap_mme_encode_pdu(&pdu, &buffer_p, &length) < 0) {
      // TODO: handle something
      OAILOG_ERROR(LOG_S1AP, "Encoding of s1ap_E_RABSetupRequestIEs failed \n");
      OAILOG_FUNC_RETURN(LOG_S1AP, RETURNerror);
    }

    OAILOG_NOTICE(
        LOG_S1AP,
        "Send S1AP E_RABSetup message MME_UE_S1AP_ID = " MME_UE_S1AP_ID_FMT
        " eNB_UE_S1AP_ID = " ENB_UE_S1AP_ID_FMT "\n",
        (mme_ue_s1ap_id_t) ue_ref->mme_ue_s1ap_id,
        (enb_ue_s1ap_id_t) ue_ref->enb_ue_s1ap_id);
    bstring b = blk2bstr(buffer_p, length);
    free(buffer_p);
    s1ap_mme_itti_send_sctp_request(
        &b, ue_ref->sctp_assoc_id, ue_ref->sctp_stream_send,
        ue_ref->mme_ue_s1ap_id);
  }

  OAILOG_FUNC_RETURN(LOG_S1AP, RETURNok);
}

//------------------------------------------------------------------------------
void s1ap_handle_conn_est_cnf(
    s1ap_state_t* state,
    const itti_mme_app_connection_establishment_cnf_t* const conn_est_cnf_pP) {
  /*
   * We received create session response from S-GW on S11 interface abstraction.
   * At least one bearer has been established. We can now send s1ap initial
   * context setup request message to eNB.
   */
  uint8_t* buffer_p                           = NULL;
  uint8_t err                                 = 0;
  uint32_t length                             = 0;
  ue_network_capability_t uenetworkcapability = {0};
  ue_description_t* ue_ref                    = NULL;
  S1ap_InitialContextSetupRequest_t* out;
  S1ap_InitialContextSetupRequestIEs_t* ie = NULL;
  S1ap_S1AP_PDU_t pdu                      = {0};  // yes, alloc on stack

  OAILOG_FUNC_IN(LOG_S1AP);
  if (conn_est_cnf_pP == NULL) {
    OAILOG_DEBUG(LOG_S1AP, "conn_est_cnf_pP is NULL\n");
    return;
  }

  OAILOG_INFO(
      LOG_S1AP,
      "Received Connection Establishment Confirm from MME_APP for ue_id = %u\n",
      conn_est_cnf_pP->ue_id);
  ue_ref = s1ap_state_get_ue_mmeid(conn_est_cnf_pP->ue_id);
  if (!ue_ref) {
    OAILOG_ERROR(
        LOG_S1AP,
        "This mme ue s1ap id (" MME_UE_S1AP_ID_FMT
        ") is not attached to any UE context\n",
        conn_est_cnf_pP->ue_id);
    // There are some race conditions were NAS T3450 timer is stopped and
    // removed at same time
    OAILOG_FUNC_OUT(LOG_S1AP);
  }

  imsi64_t imsi64;
  s1ap_imsi_map_t* imsi_map = get_s1ap_imsi_map();
  hashtable_uint64_ts_get(
      imsi_map->mme_ue_id_imsi_htbl, (const hash_key_t) conn_est_cnf_pP->ue_id,
      &imsi64);

  pdu.present = S1ap_S1AP_PDU_PR_initiatingMessage;
  pdu.choice.initiatingMessage.procedureCode =
      S1ap_ProcedureCode_id_InitialContextSetup;
  pdu.choice.initiatingMessage.value.present =
      S1ap_InitiatingMessage__value_PR_InitialContextSetupRequest;
  pdu.choice.initiatingMessage.criticality = S1ap_Criticality_reject;
  out = &pdu.choice.initiatingMessage.value.choice.InitialContextSetupRequest;

  /* mandatory */
  ie = (S1ap_InitialContextSetupRequestIEs_t*) calloc(
      1, sizeof(S1ap_InitialContextSetupRequestIEs_t));
  ie->id          = S1ap_ProtocolIE_ID_id_MME_UE_S1AP_ID;
  ie->criticality = S1ap_Criticality_reject;
  ie->value.present =
      S1ap_InitialContextSetupRequestIEs__value_PR_MME_UE_S1AP_ID;
  ie->value.choice.MME_UE_S1AP_ID = ue_ref->mme_ue_s1ap_id;
  ASN_SEQUENCE_ADD(&out->protocolIEs.list, ie);

  /* mandatory */
  ie = (S1ap_InitialContextSetupRequestIEs_t*) calloc(
      1, sizeof(S1ap_InitialContextSetupRequestIEs_t));
  ie->id          = S1ap_ProtocolIE_ID_id_eNB_UE_S1AP_ID;
  ie->criticality = S1ap_Criticality_reject;
  ie->value.present =
      S1ap_InitialContextSetupRequestIEs__value_PR_ENB_UE_S1AP_ID;
  ie->value.choice.ENB_UE_S1AP_ID = ue_ref->enb_ue_s1ap_id;
  ASN_SEQUENCE_ADD(&out->protocolIEs.list, ie);

  /* mandatory */
  ie = (S1ap_InitialContextSetupRequestIEs_t*) calloc(
      1, sizeof(S1ap_InitialContextSetupRequestIEs_t));
  ie->id          = S1ap_ProtocolIE_ID_id_uEaggregateMaximumBitrate;
  ie->criticality = S1ap_Criticality_reject;
  ie->value.present =
      S1ap_InitialContextSetupRequestIEs__value_PR_UEAggregateMaximumBitrate;
  asn_uint642INTEGER(
      &ie->value.choice.UEAggregateMaximumBitrate.uEaggregateMaximumBitRateDL,
      conn_est_cnf_pP->ue_ambr.br_dl);
  asn_uint642INTEGER(
      &ie->value.choice.UEAggregateMaximumBitrate.uEaggregateMaximumBitRateUL,
      conn_est_cnf_pP->ue_ambr.br_ul);
  ASN_SEQUENCE_ADD(&out->protocolIEs.list, ie);
  /* mandatory */
  ie = (S1ap_InitialContextSetupRequestIEs_t*) calloc(
      1, sizeof(S1ap_InitialContextSetupRequestIEs_t));
  ie->id          = S1ap_ProtocolIE_ID_id_E_RABToBeSetupListCtxtSUReq;
  ie->criticality = S1ap_Criticality_reject;
  ie->value.present =
      S1ap_InitialContextSetupRequestIEs__value_PR_E_RABToBeSetupListCtxtSUReq;

  ASN_SEQUENCE_ADD(&out->protocolIEs.list, ie);
  S1ap_E_RABToBeSetupListCtxtSUReq_t* const e_rab_to_be_setup_list =
      &ie->value.choice.E_RABToBeSetupListCtxtSUReq;

  for (int item = 0; item < conn_est_cnf_pP->no_of_e_rabs; item++) {
    S1ap_E_RABToBeSetupItemCtxtSUReqIEs_t* e_rab_tobesetup_item =
        (S1ap_E_RABToBeSetupItemCtxtSUReqIEs_t*) calloc(
            1, sizeof(S1ap_E_RABToBeSetupItemCtxtSUReqIEs_t));

    e_rab_tobesetup_item->id =
        S1ap_ProtocolIE_ID_id_E_RABToBeSetupItemCtxtSUReq;
    e_rab_tobesetup_item->criticality = S1ap_Criticality_reject;
    e_rab_tobesetup_item->value.present =
        S1ap_E_RABToBeSetupItemCtxtSUReqIEs__value_PR_E_RABToBeSetupItemCtxtSUReq;
    S1ap_E_RABToBeSetupItemCtxtSUReq_t* e_RABToBeSetup =
        &e_rab_tobesetup_item->value.choice.E_RABToBeSetupItemCtxtSUReq;

    e_RABToBeSetup->e_RAB_ID = conn_est_cnf_pP->e_rab_id[item];  // 5;
    e_RABToBeSetup->e_RABlevelQoSParameters.qCI =
        conn_est_cnf_pP->e_rab_level_qos_qci[item];
    e_RABToBeSetup->e_RABlevelQoSParameters.allocationRetentionPriority
        .priorityLevel = conn_est_cnf_pP->e_rab_level_qos_priority_level[item];
    e_RABToBeSetup->e_RABlevelQoSParameters.allocationRetentionPriority
        .pre_emptionCapability =
        conn_est_cnf_pP->e_rab_level_qos_preemption_capability[item];
    e_RABToBeSetup->e_RABlevelQoSParameters.allocationRetentionPriority
        .pre_emptionVulnerability =
        conn_est_cnf_pP->e_rab_level_qos_preemption_vulnerability[item];

    if (conn_est_cnf_pP->nas_pdu[item]) {
      S1ap_NAS_PDU_t* nas_pdu = calloc(1, sizeof(S1ap_NAS_PDU_t));
      nas_pdu->size           = blength(conn_est_cnf_pP->nas_pdu[item]);
      nas_pdu->buf            = malloc(blength(conn_est_cnf_pP->nas_pdu[item]));
      memcpy(nas_pdu->buf, conn_est_cnf_pP->nas_pdu[item]->data, nas_pdu->size);
      e_RABToBeSetup->nAS_PDU = nas_pdu;
    }
    // Set the GTP-TEID. This is the S1-U S-GW TEID
    INT32_TO_OCTET_STRING(
        conn_est_cnf_pP->gtp_teid[item], &e_RABToBeSetup->gTP_TEID);
    // S-GW IP address(es) for user-plane
    e_RABToBeSetup->transportLayerAddress.buf = calloc(
        blength(conn_est_cnf_pP->transport_layer_address[item]),
        sizeof(uint8_t));
    memcpy(
        e_RABToBeSetup->transportLayerAddress.buf,
        conn_est_cnf_pP->transport_layer_address[item]->data,
        blength(conn_est_cnf_pP->transport_layer_address[item]));
    e_RABToBeSetup->transportLayerAddress.size =
        blength(conn_est_cnf_pP->transport_layer_address[item]);
    e_RABToBeSetup->transportLayerAddress.bits_unused = 0;
    ASN_SEQUENCE_ADD(&e_rab_to_be_setup_list->list, e_rab_tobesetup_item);
  }
  {
    ie = (S1ap_InitialContextSetupRequestIEs_t*) calloc(
        1, sizeof(S1ap_InitialContextSetupRequestIEs_t));
    ie->id          = S1ap_ProtocolIE_ID_id_UESecurityCapabilities;
    ie->criticality = S1ap_Criticality_reject;
    ie->value.present =
        S1ap_InitialContextSetupRequestIEs__value_PR_UESecurityCapabilities;

    S1ap_UESecurityCapabilities_t* const ue_security_capabilities =
        &ie->value.choice.UESecurityCapabilities;

    ue_security_capabilities->encryptionAlgorithms.buf =
        calloc(1, sizeof(uint16_t));
    memcpy(
        ue_security_capabilities->encryptionAlgorithms.buf,
        &conn_est_cnf_pP->ue_security_capabilities_encryption_algorithms,
        sizeof(uint16_t));
    ue_security_capabilities->encryptionAlgorithms.size        = 2;
    ue_security_capabilities->encryptionAlgorithms.bits_unused = 0;
    OAILOG_DEBUG(
        LOG_S1AP, "security_capabilities_encryption_algorithms 0x%04X\n",
        conn_est_cnf_pP->ue_security_capabilities_encryption_algorithms);

    ue_security_capabilities->integrityProtectionAlgorithms.buf =
        calloc(1, sizeof(uint16_t));
    memcpy(
        ue_security_capabilities->integrityProtectionAlgorithms.buf,
        &conn_est_cnf_pP->ue_security_capabilities_integrity_algorithms,
        sizeof(uint16_t));
    ue_security_capabilities->integrityProtectionAlgorithms.size        = 2;
    ue_security_capabilities->integrityProtectionAlgorithms.bits_unused = 0;
    OAILOG_DEBUG(
        LOG_S1AP, "security_capabilities_integrity_algorithms 0x%04X\n",
        conn_est_cnf_pP->ue_security_capabilities_integrity_algorithms);
    ASN_SEQUENCE_ADD(&out->protocolIEs.list, ie);
  }
  /* mandatory */
  ie = (S1ap_InitialContextSetupRequestIEs_t*) calloc(
      1, sizeof(S1ap_InitialContextSetupRequestIEs_t));
  ie->id            = S1ap_ProtocolIE_ID_id_SecurityKey;
  ie->criticality   = S1ap_Criticality_reject;
  ie->value.present = S1ap_InitialContextSetupRequestIEs__value_PR_SecurityKey;
  if (conn_est_cnf_pP->kenb) {
    ie->value.choice.SecurityKey.buf = calloc(AUTH_KENB_SIZE, sizeof(uint8_t));
    memcpy(
        ie->value.choice.SecurityKey.buf, conn_est_cnf_pP->kenb,
        AUTH_KENB_SIZE);
    ie->value.choice.SecurityKey.size = AUTH_KENB_SIZE;
  } else {
    OAILOG_DEBUG(LOG_S1AP, "No kenb\n");
    ie->value.choice.SecurityKey.buf  = NULL;
    ie->value.choice.SecurityKey.size = 0;
  }
  ie->value.choice.SecurityKey.bits_unused = 0;
  ASN_SEQUENCE_ADD(&out->protocolIEs.list, ie);

  /* optional */
  /*
   * Only add capability information if it's not empty.
   */
  if (conn_est_cnf_pP->ue_radio_capability) {
    OAILOG_DEBUG(LOG_S1AP, "UE radio capability found, adding to message\n");

    ie = (S1ap_InitialContextSetupRequestIEs_t*) calloc(
        1, sizeof(S1ap_InitialContextSetupRequestIEs_t));
    ie->id          = S1ap_ProtocolIE_ID_id_UERadioCapability;
    ie->criticality = S1ap_Criticality_ignore;
    ie->value.present =
        S1ap_InitialContextSetupRequestIEs__value_PR_UERadioCapability;
    OCTET_STRING_fromBuf(
        &ie->value.choice.UERadioCapability,
        (const char*) conn_est_cnf_pP->ue_radio_capability->data,
        conn_est_cnf_pP->ue_radio_capability->slen);
    ASN_SEQUENCE_ADD(&out->protocolIEs.list, ie);
  }

  /* optional */
  if (conn_est_cnf_pP->nr_ue_security_capabilities_present) {
    {
      ie = (S1ap_InitialContextSetupRequestIEs_t*) calloc(
          1, sizeof(S1ap_InitialContextSetupRequestIEs_t));
      ie->id          = S1ap_ProtocolIE_ID_id_NRUESecurityCapabilities;
      ie->criticality = S1ap_Criticality_ignore;
      ie->value.present =
          S1ap_InitialContextSetupRequestIEs__value_PR_NRUESecurityCapabilities;

      S1ap_NRUESecurityCapabilities_t* const nr_ue_security_capabilities =
          &ie->value.choice.NRUESecurityCapabilities;

      nr_ue_security_capabilities->nRencryptionAlgorithms.buf =
          calloc(1, sizeof(uint16_t));
      uint16_t ahtobe16 = htobe16(
          conn_est_cnf_pP->nr_ue_security_capabilities_encryption_algorithms);
      memcpy(
          nr_ue_security_capabilities->nRencryptionAlgorithms.buf, &ahtobe16,
          sizeof(uint16_t));
      nr_ue_security_capabilities->nRencryptionAlgorithms.size        = 2;
      nr_ue_security_capabilities->nRencryptionAlgorithms.bits_unused = 0;
      OAILOG_DEBUG(
          LOG_S1AP,
          "NR ue security_capabilities_encryption_algorithms 0x%04X\n",
          conn_est_cnf_pP->nr_ue_security_capabilities_encryption_algorithms);

      nr_ue_security_capabilities->nRintegrityProtectionAlgorithms.buf =
          calloc(1, sizeof(uint16_t));
      ahtobe16 = htobe16(
          conn_est_cnf_pP->nr_ue_security_capabilities_integrity_algorithms);
      memcpy(
          nr_ue_security_capabilities->nRintegrityProtectionAlgorithms.buf,
          &ahtobe16, sizeof(uint16_t));
      nr_ue_security_capabilities->nRintegrityProtectionAlgorithms.size = 2;
      nr_ue_security_capabilities->nRintegrityProtectionAlgorithms.bits_unused =
          0;
      OAILOG_DEBUG(
          LOG_S1AP, "NR ue security_capabilities_integrity_algorithms 0x%04X\n",
          conn_est_cnf_pP->nr_ue_security_capabilities_integrity_algorithms);
      ASN_SEQUENCE_ADD(&out->protocolIEs.list, ie);
    }
  }

  if (s1ap_mme_encode_pdu(&pdu, &buffer_p, &length) < 0) {
    err = 1;
  }
  ASN_STRUCT_FREE_CONTENTS_ONLY(asn_DEF_S1ap_InitialContextSetupRequest, out);
  if (err) {
    OAILOG_FUNC_OUT(LOG_S1AP);
  }

  OAILOG_NOTICE_UE(
      LOG_S1AP, imsi64,
      "Send S1AP_INITIAL_CONTEXT_SETUP_REQUEST message MME_UE_S1AP_ID "
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
//------------------------------------------------------------------------------
void s1ap_handle_mme_ue_id_notification(
    s1ap_state_t* state,
    const itti_mme_app_s1ap_mme_ue_id_notification_t* const notification_p) {
  OAILOG_FUNC_IN(LOG_S1AP);

  if (notification_p == NULL) {
    OAILOG_DEBUG(LOG_S1AP, "notification_p is NULL\n");
    return;
  }
  sctp_assoc_id_t sctp_assoc_id   = notification_p->sctp_assoc_id;
  enb_ue_s1ap_id_t enb_ue_s1ap_id = notification_p->enb_ue_s1ap_id;
  mme_ue_s1ap_id_t mme_ue_s1ap_id = notification_p->mme_ue_s1ap_id;

  enb_description_t* enb_ref = s1ap_state_get_enb(state, sctp_assoc_id);
  if (enb_ref) {
    ue_description_t* ue_ref =
        s1ap_state_get_ue_enbid(enb_ref->sctp_assoc_id, enb_ue_s1ap_id);
    if (ue_ref) {
      ue_ref->mme_ue_s1ap_id = mme_ue_s1ap_id;
      hashtable_rc_t h_rc    = hashtable_ts_insert(
          &state->mmeid2associd, (const hash_key_t) mme_ue_s1ap_id,
          (void*) (uintptr_t) sctp_assoc_id);

      hashtable_uint64_ts_insert(
          &enb_ref->ue_id_coll, (const hash_key_t) mme_ue_s1ap_id,
          ue_ref->comp_s1ap_id);

      OAILOG_DEBUG(
          LOG_S1AP, "Num elements in ue_id_coll %zu and num ue associated %u",
          enb_ref->ue_id_coll.num_elements, enb_ref->nb_ue_associated);

      OAILOG_DEBUG(
          LOG_S1AP,
          "Associated sctp_assoc_id %d, enb_ue_s1ap_id " ENB_UE_S1AP_ID_FMT
          ", mme_ue_s1ap_id " MME_UE_S1AP_ID_FMT ":%s \n",
          sctp_assoc_id, enb_ue_s1ap_id, mme_ue_s1ap_id,
          hashtable_rc_code2string(h_rc));
      return;
    }
    OAILOG_DEBUG(
        LOG_S1AP,
        "Could not find  ue  with enb_ue_s1ap_id " ENB_UE_S1AP_ID_FMT "\n",
        enb_ue_s1ap_id);
    return;
  }
  OAILOG_DEBUG(
      LOG_S1AP, "Could not find  eNB with sctp_assoc_id %d \n", sctp_assoc_id);

  OAILOG_FUNC_OUT(LOG_S1AP);
}

//------------------------------------------------------------------------------
int s1ap_generate_s1ap_e_rab_rel_cmd(
    s1ap_state_t* state, itti_s1ap_e_rab_rel_cmd_t* const e_rab_rel_cmd) {
  OAILOG_FUNC_IN(LOG_S1AP);

  ue_description_t* ue_ref              = NULL;
  uint8_t* buffer_p                     = NULL;
  uint32_t length                       = 0;
  void* id                              = NULL;
  const enb_ue_s1ap_id_t enb_ue_s1ap_id = e_rab_rel_cmd->enb_ue_s1ap_id;
  const mme_ue_s1ap_id_t ue_id          = e_rab_rel_cmd->mme_ue_s1ap_id;

  hashtable_ts_get(
      &state->mmeid2associd, (const hash_key_t) ue_id, (void**) &id);
  if (id) {
    sctp_assoc_id_t sctp_assoc_id = (sctp_assoc_id_t)(uintptr_t) id;
    enb_description_t* enb_ref    = s1ap_state_get_enb(state, sctp_assoc_id);
    if (enb_ref) {
      ue_ref = s1ap_state_get_ue_enbid(enb_ref->sctp_assoc_id, enb_ue_s1ap_id);
    }
  }
  if (!ue_ref) {
    ue_ref = s1ap_state_get_ue_mmeid(ue_id);
  }
  if (!ue_ref) {
    /*
     * If the UE-associated logical S1-connection is not established,
     * the MME shall allocate a unique MME UE S1AP ID to be used for the UE.
     */
    OAILOG_ERROR(
        LOG_S1AP,
        "Unknown UE MME ID " MME_UE_S1AP_ID_FMT
        ", This case is not handled right now\n",
        ue_id);
    OAILOG_FUNC_RETURN(LOG_S1AP, RETURNerror);
  } else {
    /*
     * We have found the UE in the list.
     * Create new IE list message and encode it.
     */
    S1ap_S1AP_PDU_t pdu               = {0};
    S1ap_E_RABReleaseCommand_t* out   = NULL;
    S1ap_E_RABReleaseCommandIEs_t* ie = NULL;

    memset(&pdu, 0, sizeof(pdu));
    pdu.present = S1ap_S1AP_PDU_PR_initiatingMessage;
    pdu.choice.initiatingMessage.procedureCode =
        S1ap_ProcedureCode_id_E_RABRelease;
    pdu.choice.initiatingMessage.criticality = S1ap_Criticality_ignore;
    pdu.choice.initiatingMessage.value.present =
        S1ap_InitiatingMessage__value_PR_E_RABReleaseCommand;
    out = &pdu.choice.initiatingMessage.value.choice.E_RABReleaseCommand;
    /*
     * Setting UE information with the ones found in ue_ref
     */
    /* mandatory */
    ie = (S1ap_E_RABReleaseCommandIEs_t*) calloc(
        1, sizeof(S1ap_E_RABReleaseCommandIEs_t));
    ie->id            = S1ap_ProtocolIE_ID_id_MME_UE_S1AP_ID;
    ie->criticality   = S1ap_Criticality_reject;
    ie->value.present = S1ap_E_RABReleaseCommandIEs__value_PR_MME_UE_S1AP_ID;
    ie->value.choice.MME_UE_S1AP_ID = ue_ref->mme_ue_s1ap_id;
    ASN_SEQUENCE_ADD(&out->protocolIEs.list, ie);
    /* mandatory */
    ie = (S1ap_E_RABReleaseCommandIEs_t*) calloc(
        1, sizeof(S1ap_E_RABReleaseCommandIEs_t));
    ie->id            = S1ap_ProtocolIE_ID_id_eNB_UE_S1AP_ID;
    ie->criticality   = S1ap_Criticality_reject;
    ie->value.present = S1ap_E_RABReleaseCommandIEs__value_PR_ENB_UE_S1AP_ID;
    ie->value.choice.ENB_UE_S1AP_ID = ue_ref->enb_ue_s1ap_id;
    ASN_SEQUENCE_ADD(&out->protocolIEs.list, ie);
    ue_ref->s1_ue_state = S1AP_UE_CONNECTED;

    ie = (S1ap_E_RABReleaseCommandIEs_t*) calloc(
        1, sizeof(S1ap_E_RABReleaseCommandIEs_t));
    ie->id            = S1ap_ProtocolIE_ID_id_E_RABToBeReleasedList;
    ie->criticality   = S1ap_Criticality_ignore;
    ie->value.present = S1ap_E_RABReleaseCommandIEs__value_PR_E_RABList;
    ASN_SEQUENCE_ADD(&out->protocolIEs.list, ie);
    S1ap_E_RABList_t* const e_rab_list = &ie->value.choice.E_RABList;

    for (int i = 0; i < e_rab_rel_cmd->e_rab_to_be_rel_list.no_of_items; i++) {
      S1ap_E_RABItemIEs_t* s1ap_e_rab_item_ies =
          calloc(1, sizeof(S1ap_E_RABItemIEs_t));
      s1ap_e_rab_item_ies->id          = S1ap_ProtocolIE_ID_id_E_RABItem;
      s1ap_e_rab_item_ies->criticality = S1ap_Criticality_ignore;
      s1ap_e_rab_item_ies->value.present =
          S1ap_E_RABItemIEs__value_PR_E_RABItem;

      S1ap_E_RABItem_t* s1ap_e_rab_item =
          &s1ap_e_rab_item_ies->value.choice.E_RABItem;

      s1ap_e_rab_item->e_RAB_ID =
          e_rab_rel_cmd->e_rab_to_be_rel_list.item[i].e_rab_id;
      s1ap_mme_set_cause(
          &s1ap_e_rab_item->cause, S1ap_Cause_PR_radioNetwork,
          S1ap_CauseRadioNetwork_unspecified);

      ASN_SEQUENCE_ADD(&e_rab_list->list, s1ap_e_rab_item_ies);
    }

    /*
     * Fill in the NAS pdu
     */
    if (e_rab_rel_cmd->nas_pdu) {
      ie = (S1ap_E_RABReleaseCommandIEs_t*) calloc(
          1, sizeof(S1ap_E_RABReleaseCommandIEs_t));
      ie->id            = S1ap_ProtocolIE_ID_id_NAS_PDU;
      ie->criticality   = S1ap_Criticality_ignore;
      ie->value.present = S1ap_E_RABReleaseCommandIEs__value_PR_NAS_PDU;

      S1ap_NAS_PDU_t* nas_pdu = &ie->value.choice.NAS_PDU;
      OCTET_STRING_fromBuf(
          nas_pdu, (char*) bdata(e_rab_rel_cmd->nas_pdu),
          blength(e_rab_rel_cmd->nas_pdu));
      ASN_SEQUENCE_ADD(&out->protocolIEs.list, ie);
    } else {
      OAILOG_INFO(
          LOG_S1AP,
          "No NAS message received for S1AP E-RAB release command for "
          "ueId " MME_UE_S1AP_ID_FMT " .\n",
          e_rab_rel_cmd->mme_ue_s1ap_id);
    }

    if (s1ap_mme_encode_pdu(&pdu, &buffer_p, &length) < 0) {
      OAILOG_ERROR(
          LOG_S1AP, "Encoding of s1ap_E_RABReleaseCommandIEs failed \n");
      OAILOG_FUNC_RETURN(LOG_S1AP, RETURNerror);
    }

    OAILOG_NOTICE(
        LOG_S1AP,
        "Send S1AP E_RABRelease Command message MME_UE_S1AP_ID "
        "= " MME_UE_S1AP_ID_FMT " eNB_UE_S1AP_ID = " ENB_UE_S1AP_ID_FMT "\n",
        (mme_ue_s1ap_id_t) ue_ref->mme_ue_s1ap_id,
        (enb_ue_s1ap_id_t) ue_ref->enb_ue_s1ap_id);
    bstring b = blk2bstr(buffer_p, length);
    free(buffer_p);
    s1ap_mme_itti_send_sctp_request(
        &b, ue_ref->sctp_assoc_id, ue_ref->sctp_stream_send,
        ue_ref->mme_ue_s1ap_id);
  }

  OAILOG_FUNC_RETURN(LOG_S1AP, RETURNok);
}

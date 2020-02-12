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
#include "s1ap_ies_defs.h"
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
#include "S1AP-PDU.h"
#include "S1ap-AllocationAndRetentionPriority.h"
#include "S1ap-E-RABLevelQoSParameters.h"
#include "S1ap-E-RABToBeSetupItemBearerSUReq.h"
#include "S1ap-E-RABToBeSetupItemCtxtSUReq.h"
#include "S1ap-EUTRAN-CGI.h"
#include "S1ap-EncryptionAlgorithms.h"
#include "S1ap-GBR-QosInformation.h"
#include "S1ap-GUMMEI.h"
#include "S1ap-IntegrityProtectionAlgorithms.h"
#include "S1ap-NAS-PDU.h"
#include "S1ap-PLMNidentity.h"
#include "S1ap-ProcedureCode.h"
#include "S1ap-S-TMSI.h"
#include "S1ap-SecurityKey.h"
#include "S1ap-TAI.h"
#include "S1ap-TransportLayerAddress.h"
#include "S1ap-UEAggregateMaximumBitrate.h"
#include "S1ap-UESecurityCapabilities.h"
#include "TrackingAreaIdentity.h"
#include "asn_SEQUENCE_OF.h"
#include "nas/securityDef.h"
#include "s1ap_state.h"
#include "S1ap-CauseMisc.h"
#include "S1ap-CauseNas.h"
#include "S1ap-CauseProtocol.h"
#include "S1ap-CauseRadioNetwork.h"
#include "S1ap-CauseTransport.h"
#include "S1ap-E-RABItem.h"
#include "s1ap_mme_handlers.h"

//------------------------------------------------------------------------------
int s1ap_mme_handle_initial_ue_message(
  s1ap_state_t *state,
  const sctp_assoc_id_t assoc_id,
  const sctp_stream_id_t stream,
  struct s1ap_message_s *message)
{
  S1ap_InitialUEMessageIEs_t *initialUEMessage_p = NULL;
  ue_description_t *ue_ref = NULL;
  enb_description_t *eNB_ref = NULL;
  enb_ue_s1ap_id_t enb_ue_s1ap_id = 0;

  OAILOG_FUNC_IN(LOG_S1AP);
  initialUEMessage_p = &message->msg.s1ap_InitialUEMessageIEs;

  OAILOG_INFO(
    LOG_S1AP,
    "Received S1AP INITIAL_UE_MESSAGE eNB_UE_S1AP_ID " ENB_UE_S1AP_ID_FMT "\n",
    (enb_ue_s1ap_id_t) initialUEMessage_p->eNB_UE_S1AP_ID);

  if ((eNB_ref = s1ap_state_get_enb(state, assoc_id)) == NULL) {
    OAILOG_ERROR(LOG_S1AP, "Unknown eNB on assoc_id %d\n", assoc_id);
    OAILOG_FUNC_RETURN(LOG_S1AP, RETURNerror);
  }
  // eNB UE S1AP ID is limited to 24 bits
  enb_ue_s1ap_id =
    (enb_ue_s1ap_id_t)(initialUEMessage_p->eNB_UE_S1AP_ID & 0x00ffffff);
  OAILOG_INFO(
    LOG_S1AP,
    "New Initial UE message received with eNB UE S1AP ID: " ENB_UE_S1AP_ID_FMT
    "\n",
    enb_ue_s1ap_id);
  ue_ref = s1ap_state_get_ue_enbid(eNB_ref, enb_ue_s1ap_id);

  if (ue_ref == NULL) {
    tai_t tai = {0};
    gummei_t gummei = {
      .plmn = {0}, .mme_code = 0, .mme_gid = 0}; // initialized after
    s_tmsi_t s_tmsi = {.mme_code = 0, .m_tmsi = INVALID_M_TMSI};
    ecgi_t ecgi = {.plmn = {0}, .cell_identity = {0}};
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

    ue_ref->s1_ue_state = S1AP_UE_WAITING_CSR;

    ue_ref->enb_ue_s1ap_id = enb_ue_s1ap_id;
    // Will be allocated by NAS
    ue_ref->mme_ue_s1ap_id = INVALID_MME_UE_S1AP_ID;

    ue_ref->s1ap_ue_context_rel_timer.id = S1AP_TIMER_INACTIVE_ID;
    ue_ref->s1ap_ue_context_rel_timer.sec = S1AP_UE_CONTEXT_REL_COMP_TIMER;

    // On which stream we received the message
    ue_ref->sctp_stream_recv = stream;
    ue_ref->sctp_stream_send = ue_ref->enb->next_sctp_stream;

    /*
     * Increment the sctp stream for the eNB association.
     * If the next sctp stream is >= instream negociated between eNB and MME, wrap to first stream.
     * TODO: search for the first available stream instead.
     */

    /*
     * TODO task#15456359.
     * Below logic seems to be incorrect , revisit it.
     */
    ue_ref->enb->next_sctp_stream += 1;
    if (ue_ref->enb->next_sctp_stream >= ue_ref->enb->instreams) {
      ue_ref->enb->next_sctp_stream = 1;
    }
    s1ap_dump_enb(ue_ref->enb);
    // TAI mandatory IE
    OCTET_STRING_TO_TAC(&initialUEMessage_p->tai.tAC, tai.tac);
    DevAssert(initialUEMessage_p->tai.pLMNidentity.size == 3);
    TBCD_TO_PLMN_T(&initialUEMessage_p->tai.pLMNidentity, &tai);

    // CGI mandatory IE
    DevAssert(initialUEMessage_p->eutran_cgi.pLMNidentity.size == 3);
    TBCD_TO_PLMN_T(&initialUEMessage_p->eutran_cgi.pLMNidentity, &ecgi.plmn);
    BIT_STRING_TO_CELL_IDENTITY(
      &initialUEMessage_p->eutran_cgi.cell_ID, ecgi.cell_identity);

    if (
      initialUEMessage_p->presenceMask &
      S1AP_INITIALUEMESSAGEIES_S_TMSI_PRESENT) {
      OCTET_STRING_TO_MME_CODE(
        &initialUEMessage_p->s_tmsi.mMEC, s_tmsi.mme_code);
      OCTET_STRING_TO_M_TMSI(&initialUEMessage_p->s_tmsi.m_TMSI, s_tmsi.m_tmsi);
    }

    if (
      initialUEMessage_p->presenceMask &
      S1AP_INITIALUEMESSAGEIES_CSG_ID_PRESENT) {
      csg_id = BIT_STRING_to_uint32(&initialUEMessage_p->csG_Id);
    }

    memset(&gummei, 0, sizeof(gummei));
    if (
      initialUEMessage_p->presenceMask &
      S1AP_INITIALUEMESSAGEIES_GUMMEI_ID_PRESENT) {
      TBCD_TO_PLMN_T(
        &initialUEMessage_p->gummei_id.pLMN_Identity, &gummei.plmn);
      OCTET_STRING_TO_MME_GID(
        &initialUEMessage_p->gummei_id.mME_Group_ID, gummei.mme_gid);
      OCTET_STRING_TO_MME_CODE(
        &initialUEMessage_p->gummei_id.mME_Code, gummei.mme_code);
    }
    /*
     * We received the first NAS transport message: initial UE message.
     * * * * Send a NAS ESTAeNBBLISH IND to NAS layer
     */
#if ORIGINAL_CODE
    s1ap_mme_itti_nas_establish_ind(
      ue_ref->mme_ue_s1ap_id,
      initialUEMessage_p->nas_pdu.buf,
      initialUEMessage_p->nas_pdu.size,
      initialUEMessage_p->rrC_Establishment_Cause,
      tai_tac);
#else
    s1ap_mme_itti_s1ap_initial_ue_message(
      assoc_id,
      ue_ref->enb->enb_id,
      ue_ref->enb_ue_s1ap_id,
      initialUEMessage_p->nas_pdu.buf,
      initialUEMessage_p->nas_pdu.size,
      &tai,
      &ecgi,
      initialUEMessage_p->rrC_Establishment_Cause,
      (initialUEMessage_p->presenceMask &
       S1AP_INITIALUEMESSAGEIES_S_TMSI_PRESENT) ?
        &s_tmsi :
        NULL,
      (initialUEMessage_p->presenceMask &
       S1AP_INITIALUEMESSAGEIES_CSG_ID_PRESENT) ?
        &csg_id :
        NULL,
      (initialUEMessage_p->presenceMask &
       S1AP_INITIALUEMESSAGEIES_GUMMEI_ID_PRESENT) ?
        &gummei :
        NULL,
      NULL, // CELL ACCESS MODE
      NULL, // GW Transport Layer Address
      NULL  //Relay Node Indicator
    );
#endif
  } else {
    OAILOG_ERROR(
      LOG_S1AP,
      "Initial UE Message- Duplicate ENB_UE_S1AP_ID. Ignoring the "
      "message, eNB UE S1AP ID:" ENB_UE_S1AP_ID_FMT "\n",
      enb_ue_s1ap_id);
  }

  OAILOG_FUNC_RETURN(LOG_S1AP, RETURNok);
}

//------------------------------------------------------------------------------
int s1ap_mme_handle_uplink_nas_transport(
  s1ap_state_t *state,
  const sctp_assoc_id_t assoc_id,
  __attribute__((unused)) const sctp_stream_id_t stream,
  struct s1ap_message_s *message)
{
  S1ap_UplinkNASTransportIEs_t *uplinkNASTransport_p = NULL;
  ue_description_t *ue_ref = NULL;
  enb_description_t *enb_ref = NULL;
  tai_t tai = {0};
  ecgi_t ecgi = {.plmn = {0}, .cell_identity = {0}};

  OAILOG_FUNC_IN(LOG_S1AP);
  uplinkNASTransport_p = &message->msg.s1ap_UplinkNASTransportIEs;

  if (INVALID_MME_UE_S1AP_ID == uplinkNASTransport_p->mme_ue_s1ap_id) {
    OAILOG_WARNING(
      LOG_S1AP,
      "Received S1AP UPLINK_NAS_TRANSPORT message MME_UE_S1AP_ID unknown\n");

    enb_ref = s1ap_state_get_enb(state, assoc_id);

    if (!(ue_ref = s1ap_state_get_ue_enbid(
            enb_ref,
            (enb_ue_s1ap_id_t) uplinkNASTransport_p->eNB_UE_S1AP_ID))) {
      OAILOG_WARNING(
        LOG_S1AP,
        "Received S1AP UPLINK_NAS_TRANSPORT No UE is attached to this "
        "enb_ue_s1ap_id: " ENB_UE_S1AP_ID_FMT "\n",
        (enb_ue_s1ap_id_t) uplinkNASTransport_p->eNB_UE_S1AP_ID);
      OAILOG_FUNC_RETURN(LOG_S1AP, RETURNerror);
    }
  } else {
    OAILOG_INFO(
      LOG_S1AP,
      "Received S1AP UPLINK_NAS_TRANSPORT message "
      "MME_UE_S1AP_ID " MME_UE_S1AP_ID_FMT "\n",
      (mme_ue_s1ap_id_t) uplinkNASTransport_p->mme_ue_s1ap_id);

    if (!(ue_ref = s1ap_state_get_ue_mmeid(
            state, uplinkNASTransport_p->mme_ue_s1ap_id))) {
      OAILOG_WARNING(
        LOG_S1AP,
        "Received S1AP UPLINK_NAS_TRANSPORT No UE is attached to this "
        "mme_ue_s1ap_id: " MME_UE_S1AP_ID_FMT "\n",
        (mme_ue_s1ap_id_t) uplinkNASTransport_p->mme_ue_s1ap_id);
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

  // TAI mandatory IE
  OCTET_STRING_TO_TAC(&uplinkNASTransport_p->tai.tAC, tai.tac);
  DevAssert(uplinkNASTransport_p->tai.pLMNidentity.size == 3);
  TBCD_TO_PLMN_T(&uplinkNASTransport_p->tai.pLMNidentity, &tai);

  // CGI mandatory IE
  DevAssert(uplinkNASTransport_p->eutran_cgi.pLMNidentity.size == 3);
  TBCD_TO_PLMN_T(&uplinkNASTransport_p->eutran_cgi.pLMNidentity, &ecgi.plmn);
  BIT_STRING_TO_CELL_IDENTITY(
    &uplinkNASTransport_p->eutran_cgi.cell_ID, ecgi.cell_identity);

  // TODO optional GW Transport Layer Address

  bstring b = blk2bstr(
    uplinkNASTransport_p->nas_pdu.buf, uplinkNASTransport_p->nas_pdu.size);
  s1ap_mme_itti_nas_uplink_ind(
    uplinkNASTransport_p->mme_ue_s1ap_id, &b, &tai, &ecgi);
  OAILOG_FUNC_RETURN(LOG_S1AP, RETURNok);
}

//------------------------------------------------------------------------------
int s1ap_mme_handle_nas_non_delivery(
  s1ap_state_t *state,
  __attribute__((unused)) sctp_assoc_id_t assoc_id,
  sctp_stream_id_t stream,
  struct s1ap_message_s *message)
{
  S1ap_NASNonDeliveryIndication_IEs_t *nasNonDeliveryIndication_p = NULL;
  ue_description_t *ue_ref = NULL;
  imsi64_t imsi64 = INVALID_IMSI64;

  OAILOG_FUNC_IN(LOG_S1AP);
  increment_counter("nas_non_delivery_indication_received", 1, NO_LABELS);
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

  nasNonDeliveryIndication_p = &message->msg.s1ap_NASNonDeliveryIndication_IEs;

  OAILOG_NOTICE(
    LOG_S1AP,
    "Received S1AP NAS_NON_DELIVERY_INDICATION message "
    "MME_UE_S1AP_ID " MME_UE_S1AP_ID_FMT " enb_ue_s1ap_id " ENB_UE_S1AP_ID_FMT
    "\n",
    (mme_ue_s1ap_id_t) nasNonDeliveryIndication_p->mme_ue_s1ap_id,
    (enb_ue_s1ap_id_t) nasNonDeliveryIndication_p->eNB_UE_S1AP_ID);

  if (
    (ue_ref = s1ap_state_get_ue_mmeid(
       state, nasNonDeliveryIndication_p->mme_ue_s1ap_id)) == NULL) {
    OAILOG_DEBUG(
      LOG_S1AP,
      "No UE is attached to this mme UE s1ap id: " MME_UE_S1AP_ID_FMT "\n",
      (mme_ue_s1ap_id_t) nasNonDeliveryIndication_p->mme_ue_s1ap_id);
    OAILOG_FUNC_RETURN(LOG_S1AP, RETURNerror);
  }

  if (ue_ref->s1_ue_state != S1AP_UE_CONNECTED) {
    OAILOG_DEBUG(
      LOG_S1AP,
      "Received S1AP NAS_NON_DELIVERY_INDICATION while UE in state != "
      "S1AP_UE_CONNECTED\n");
    OAILOG_FUNC_RETURN(LOG_S1AP, RETURNerror);
  }

  s1ap_imsi_map_t* imsi_map = get_s1ap_imsi_map();
  hashtable_uint64_ts_get(
    imsi_map->mme_ue_id_imsi_htbl,
    (const hash_key_t) nasNonDeliveryIndication_p->mme_ue_s1ap_id,
    &imsi64);

  //TODO: forward NAS PDU to NAS
  s1ap_mme_itti_nas_non_delivery_ind(
    nasNonDeliveryIndication_p->mme_ue_s1ap_id,
    nasNonDeliveryIndication_p->nas_pdu.buf,
    nasNonDeliveryIndication_p->nas_pdu.size,
    &nasNonDeliveryIndication_p->cause,
    imsi64);
  OAILOG_FUNC_RETURN(LOG_S1AP, RETURNok);
}

//------------------------------------------------------------------------------
int s1ap_generate_downlink_nas_transport(
  s1ap_state_t *state,
  const enb_ue_s1ap_id_t enb_ue_s1ap_id,
  const mme_ue_s1ap_id_t ue_id,
  STOLEN_REF bstring *payload,
  const imsi64_t imsi64)
{
  ue_description_t *ue_ref = NULL;
  uint8_t *buffer_p = NULL;
  uint32_t length = 0;
  void *id = NULL;

  OAILOG_FUNC_IN(LOG_S1AP);

  // Try to retrieve SCTP assoication id using mme_ue_s1ap_id
  if (
    HASH_TABLE_OK ==
    hashtable_ts_get(
      &state->mmeid2associd, (const hash_key_t) ue_id, (void **) &id)) {
    sctp_assoc_id_t sctp_assoc_id = (sctp_assoc_id_t)(uintptr_t) id;
    enb_description_t *enb_ref = s1ap_state_get_enb(state, sctp_assoc_id);
    if (enb_ref) {
      ue_ref = s1ap_state_get_ue_enbid(enb_ref, enb_ue_s1ap_id);
    } else {
      OAILOG_ERROR(
        LOG_S1AP, "No eNB for SCTP association id %d \n", sctp_assoc_id);
      OAILOG_FUNC_RETURN(LOG_S1AP, RETURNerror);
    }
  }
  // TODO remove soon:
  if (!ue_ref) {
    ue_ref = s1ap_state_get_ue_mmeid(state, ue_id);
  }
  // finally!
  if (!ue_ref) {
    /*
     * If the UE-associated logical S1-connection is not established,
     * * * * the MME shall allocate a unique MME UE S1AP ID to be used for the UE.
     */
    OAILOG_WARNING(
      LOG_S1AP,
      "Unknown UE MME ID " MME_UE_S1AP_ID_FMT
      ", This case is not handled right now\n",
      ue_id);
    OAILOG_FUNC_RETURN(LOG_S1AP, RETURNerror);
  } else {
    /*
     * We have fount the UE in the list.
     * * * * Create new IE list message and encode it.
     */
    s1ap_imsi_map_t* imsi_map = get_s1ap_imsi_map();
    hashtable_uint64_ts_insert(
      imsi_map->enb_s1ap_mme_ue_id_htbl,
      (const hash_key_t) enb_ue_s1ap_id,
      ue_id);
    hashtable_uint64_ts_insert(
      imsi_map->mme_ue_id_imsi_htbl,
      (const hash_key_t) ue_id,
      imsi64);

    S1ap_DownlinkNASTransportIEs_t *downlinkNasTransport = NULL;
    s1ap_message message = {0};

    message.procedureCode = S1ap_ProcedureCode_id_downlinkNASTransport;
    message.direction = S1AP_PDU_PR_initiatingMessage;
    if (ue_ref->s1_ue_state == S1AP_UE_WAITING_CRR) {
      OAILOG_ERROR(
        LOG_S1AP, "Already triggred UE Context Release Command and UE is"
        "in S1AP_UE_WAITING_CRR, so dropping the DownlinkNASTransport \n");
      OAILOG_FUNC_RETURN(LOG_S1AP, RETURNerror);
    } else {
    ue_ref->s1_ue_state = S1AP_UE_CONNECTED;
    }
    downlinkNasTransport = &message.msg.s1ap_DownlinkNASTransportIEs;
    /*
     * Setting UE informations with the ones fount in ue_ref
     */
    downlinkNasTransport->mme_ue_s1ap_id = ue_ref->mme_ue_s1ap_id;
    downlinkNasTransport->eNB_UE_S1AP_ID = ue_ref->enb_ue_s1ap_id;
    /*eNB
     * Fill in the NAS pdu
     */
    OCTET_STRING_fromBuf(
      &downlinkNasTransport->nas_pdu,
      (char *) bdata(*payload),
      blength(*payload));
    bdestroy_wrapper(payload);

    if (s1ap_mme_encode_pdu(&message, &buffer_p, &length) < 0) {
      free_s1ap_downlinknastransport(downlinkNasTransport);
      OAILOG_FUNC_RETURN(LOG_S1AP, RETURNerror);
    }

    OAILOG_NOTICE(
      LOG_S1AP,
      "Send S1AP DOWNLINK_NAS_TRANSPORT message ue_id = " MME_UE_S1AP_ID_FMT
      " MME_UE_S1AP_ID = " MME_UE_S1AP_ID_FMT
      " eNB_UE_S1AP_ID = " ENB_UE_S1AP_ID_FMT "\n",
      ue_id,
      (mme_ue_s1ap_id_t) downlinkNasTransport->mme_ue_s1ap_id,
      (enb_ue_s1ap_id_t) downlinkNasTransport->eNB_UE_S1AP_ID);
    bstring b = blk2bstr(buffer_p, length);
    free(buffer_p);
    s1ap_mme_itti_send_sctp_request(
      &b,
      ue_ref->enb->sctp_assoc_id,
      ue_ref->sctp_stream_send,
      ue_ref->mme_ue_s1ap_id);
    free_s1ap_downlinknastransport(downlinkNasTransport);
  }

  OAILOG_FUNC_RETURN(LOG_S1AP, RETURNok);
}

//------------------------------------------------------------------------------
int s1ap_generate_s1ap_e_rab_setup_req(
  s1ap_state_t *state,
  itti_s1ap_e_rab_setup_req_t *const e_rab_setup_req)
{
  OAILOG_FUNC_IN(LOG_S1AP);
  ue_description_t *ue_ref = NULL;
  uint8_t *buffer_p = NULL;
  uint32_t length = 0;
  void *id = NULL;
  const enb_ue_s1ap_id_t enb_ue_s1ap_id = e_rab_setup_req->enb_ue_s1ap_id;
  const mme_ue_s1ap_id_t ue_id = e_rab_setup_req->mme_ue_s1ap_id;

  hashtable_ts_get(
    &state->mmeid2associd, (const hash_key_t) ue_id, (void **) &id);
  if (id) {
    sctp_assoc_id_t sctp_assoc_id = (sctp_assoc_id_t)(uintptr_t) id;
    enb_description_t *enb_ref = s1ap_state_get_enb(state, sctp_assoc_id);
    if (enb_ref) {
      ue_ref = s1ap_state_get_ue_enbid(enb_ref, enb_ue_s1ap_id);
    }
  }
  // TODO remove soon:
  if (!ue_ref) {
    ue_ref = s1ap_state_get_ue_mmeid(state, ue_id);
  }
  // finally!
  if (!ue_ref) {
    /*
     * If the UE-associated logical S1-connection is not established,
     * * * * the MME shall allocate a unique MME UE S1AP ID to be used for the UE.
     */
    OAILOG_DEBUG(
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
    S1ap_E_RABSetupRequestIEs_t *e_rabsetuprequesties = NULL;
    s1ap_message message = {0};

    message.procedureCode = S1ap_ProcedureCode_id_E_RABSetup;
    message.direction = S1AP_PDU_PR_initiatingMessage;
    ue_ref->s1_ue_state = S1AP_UE_CONNECTED;
    e_rabsetuprequesties = &message.msg.s1ap_E_RABSetupRequestIEs;
    /*
     * Setting UE information with the ones found in ue_ref
     */
    e_rabsetuprequesties->mme_ue_s1ap_id = ue_ref->mme_ue_s1ap_id;
    e_rabsetuprequesties->eNB_UE_S1AP_ID = ue_ref->enb_ue_s1ap_id;
    /*eNB
     * Fill in the NAS pdu
     */
    e_rabsetuprequesties->presenceMask = 0;
    //    if (e_rab_setup_req->ue_aggregate_maximum_bit_rate_present) {
    //      e_rabsetuprequesties->presenceMask |= S1AP_E_RABSETUPREQUESTIES_UEAGGREGATEMAXIMUMBITRATE_PRESENT;
    //      TO DO e_rabsetuprequesties->uEaggregateMaximumBitrate.uEaggregateMaximumBitRateDL.buf
    //    }

    S1ap_E_RABToBeSetupItemBearerSUReq_t s1ap_E_RABToBeSetupItemBearerSUReq
      [e_rab_setup_req->e_rab_to_be_setup_list.no_of_items];
    struct S1ap_GBR_QosInformation
      gbrQosInformation[e_rab_setup_req->e_rab_to_be_setup_list.no_of_items];

    for (int i = 0; i < e_rab_setup_req->e_rab_to_be_setup_list.no_of_items;
         i++) {
      memset(
        &s1ap_E_RABToBeSetupItemBearerSUReq[i],
        0,
        sizeof(S1ap_E_RABToBeSetupItemBearerSUReq_t));

      s1ap_E_RABToBeSetupItemBearerSUReq[i].e_RAB_ID =
        e_rab_setup_req->e_rab_to_be_setup_list.item[i].e_rab_id;
      s1ap_E_RABToBeSetupItemBearerSUReq[i].e_RABlevelQoSParameters.qCI =
        e_rab_setup_req->e_rab_to_be_setup_list.item[i]
          .e_rab_level_qos_parameters.qci;

      s1ap_E_RABToBeSetupItemBearerSUReq[i]
        .e_RABlevelQoSParameters.allocationRetentionPriority.priorityLevel =
        e_rab_setup_req->e_rab_to_be_setup_list.item[i]
          .e_rab_level_qos_parameters.allocation_and_retention_priority
          .priority_level;

      s1ap_E_RABToBeSetupItemBearerSUReq[i]
        .e_RABlevelQoSParameters.allocationRetentionPriority
        .pre_emptionCapability =
        e_rab_setup_req->e_rab_to_be_setup_list.item[i]
          .e_rab_level_qos_parameters.allocation_and_retention_priority
          .pre_emption_capability;

      s1ap_E_RABToBeSetupItemBearerSUReq[i]
        .e_RABlevelQoSParameters.allocationRetentionPriority
        .pre_emptionVulnerability =
        e_rab_setup_req->e_rab_to_be_setup_list.item[i]
          .e_rab_level_qos_parameters.allocation_and_retention_priority
          .pre_emption_vulnerability;
      /* OPTIONAL */
      gbr_qos_information_t *gbr_qos_information =
        &e_rab_setup_req->e_rab_to_be_setup_list.item[i]
           .e_rab_level_qos_parameters.gbr_qos_information;
      if (
        (gbr_qos_information->e_rab_maximum_bit_rate_downlink) ||
        (gbr_qos_information->e_rab_maximum_bit_rate_uplink) ||
        (gbr_qos_information->e_rab_guaranteed_bit_rate_downlink) ||
        (gbr_qos_information->e_rab_guaranteed_bit_rate_uplink)) {
        OAILOG_NOTICE(
          LOG_S1AP, "Encoding of e_RABlevelQoSParameters.gbrQosInformation\n");

        //s1ap_E_RABToBeSetupItemBearerSUReq[i].e_RABlevelQoSParameters.gbrQosInformation = calloc(1, sizeof(struct S1ap_GBR_QosInformation));
        s1ap_E_RABToBeSetupItemBearerSUReq[i]
          .e_RABlevelQoSParameters.gbrQosInformation = &gbrQosInformation[i];
        memset(&gbrQosInformation[i], 0, sizeof(gbrQosInformation[i]));
        if (s1ap_E_RABToBeSetupItemBearerSUReq[i]
              .e_RABlevelQoSParameters.gbrQosInformation) {
          asn_uint642INTEGER(
            &s1ap_E_RABToBeSetupItemBearerSUReq[i]
               .e_RABlevelQoSParameters.gbrQosInformation
               ->e_RAB_MaximumBitrateDL,
            gbr_qos_information->e_rab_maximum_bit_rate_downlink);

          asn_uint642INTEGER(
            &s1ap_E_RABToBeSetupItemBearerSUReq[i]
               .e_RABlevelQoSParameters.gbrQosInformation
               ->e_RAB_MaximumBitrateUL,
            gbr_qos_information->e_rab_maximum_bit_rate_uplink);

          asn_uint642INTEGER(
            &s1ap_E_RABToBeSetupItemBearerSUReq[i]
               .e_RABlevelQoSParameters.gbrQosInformation
               ->e_RAB_GuaranteedBitrateDL,
            gbr_qos_information->e_rab_guaranteed_bit_rate_downlink);

          asn_uint642INTEGER(
            &s1ap_E_RABToBeSetupItemBearerSUReq[i]
               .e_RABlevelQoSParameters.gbrQosInformation
               ->e_RAB_GuaranteedBitrateUL,
            gbr_qos_information->e_rab_guaranteed_bit_rate_uplink);
        }
      } else {
        OAILOG_NOTICE(
          LOG_S1AP,
          "NOT Encoding of e_RABlevelQoSParameters.gbrQosInformation\n");
      }

      INT32_TO_OCTET_STRING(
        e_rab_setup_req->e_rab_to_be_setup_list.item[i].gtp_teid,
        &s1ap_E_RABToBeSetupItemBearerSUReq[i].gTP_TEID);

      s1ap_E_RABToBeSetupItemBearerSUReq[i].transportLayerAddress.buf = calloc(
        blength(e_rab_setup_req->e_rab_to_be_setup_list.item[i]
                  .transport_layer_address),
        sizeof(uint8_t));
      memcpy(
        s1ap_E_RABToBeSetupItemBearerSUReq[i].transportLayerAddress.buf,
        e_rab_setup_req->e_rab_to_be_setup_list.item[i]
          .transport_layer_address->data,
        blength(e_rab_setup_req->e_rab_to_be_setup_list.item[i]
                  .transport_layer_address));

      s1ap_E_RABToBeSetupItemBearerSUReq[i].transportLayerAddress.size =
        blength(e_rab_setup_req->e_rab_to_be_setup_list.item[i]
                  .transport_layer_address);
      s1ap_E_RABToBeSetupItemBearerSUReq[i].transportLayerAddress.bits_unused =
        0;

      OCTET_STRING_fromBuf(
        &s1ap_E_RABToBeSetupItemBearerSUReq[i].nAS_PDU,
        (char *) bdata(e_rab_setup_req->e_rab_to_be_setup_list.item[i].nas_pdu),
        blength(e_rab_setup_req->e_rab_to_be_setup_list.item[i].nas_pdu));

      ASN_SEQUENCE_ADD(
        &e_rabsetuprequesties->e_RABToBeSetupListBearerSUReq,
        &s1ap_E_RABToBeSetupItemBearerSUReq[i]);
    }

    if (s1ap_mme_encode_pdu(&message, &buffer_p, &length) < 0) {
      // TODO: handle something
      OAILOG_ERROR(LOG_S1AP, "Encoding of s1ap_E_RABSetupRequestIEs failed \n");
      OAILOG_FUNC_RETURN(LOG_S1AP, RETURNerror);
    }

    OAILOG_NOTICE(
      LOG_S1AP,
      "Send S1AP E_RABSetup message MME_UE_S1AP_ID = " MME_UE_S1AP_ID_FMT
      " eNB_UE_S1AP_ID = " ENB_UE_S1AP_ID_FMT "\n",
      (mme_ue_s1ap_id_t) e_rabsetuprequesties->mme_ue_s1ap_id,
      (enb_ue_s1ap_id_t) e_rabsetuprequesties->eNB_UE_S1AP_ID);
    bstring b = blk2bstr(buffer_p, length);
    s1ap_mme_itti_send_sctp_request(
      &b,
      ue_ref->enb->sctp_assoc_id,
      ue_ref->sctp_stream_send,
      ue_ref->mme_ue_s1ap_id);
  }

  OAILOG_FUNC_RETURN(LOG_S1AP, RETURNok);
}

//------------------------------------------------------------------------------
void s1ap_handle_conn_est_cnf(
  s1ap_state_t *state,
  const itti_mme_app_connection_establishment_cnf_t *const conn_est_cnf_pP)
{
  /*
   * We received create session response from S-GW on S11 interface abstraction.
   * At least one bearer has been established. We can now send s1ap initial context setup request
   * message to eNB.
   */
  uint8_t *buffer_p = NULL;
  uint32_t length = 0;
  ue_description_t *ue_ref = NULL;
  S1ap_InitialContextSetupRequestIEs_t *initialContextSetupRequest_p = NULL;
  S1ap_E_RABToBeSetupItemCtxtSUReq_t
    *e_RABToBeSetup[conn_est_cnf_pP->no_of_e_rabs];
  S1ap_NAS_PDU_t *nas_pdu = NULL;
  s1ap_message message = {0}; // yes, alloc on stack

  OAILOG_FUNC_IN(LOG_S1AP);
  DevAssert(conn_est_cnf_pP != NULL);

  OAILOG_INFO(
    LOG_S1AP,
    "Received Connection Establishment Confirm from MME_APP for ue_id = %u\n",
    conn_est_cnf_pP->ue_id);
  ue_ref = s1ap_state_get_ue_mmeid(state, conn_est_cnf_pP->ue_id);
  if (!ue_ref) {
    OAILOG_ERROR(
      LOG_S1AP,
      "This mme ue s1ap id (" MME_UE_S1AP_ID_FMT
      ") is not attached to any UE context\n",
      conn_est_cnf_pP->ue_id);
    // There are some race conditions were NAS T3450 timer is stopped and removed at same time
    OAILOG_FUNC_OUT(LOG_S1AP);
  }

  /*
   * Start the outcome response timer.
   * * * * When time is reached, MME consider that procedure outcome has failed.
   */
  //     timer_setup(mme_config.s1ap_config.outcome_drop_timer_sec, 0, TASK_S1AP, INSTANCE_DEFAULT,
  //                 TIMER_ONE_SHOT,
  //                 NULL,
  //                 &ue_ref->outcome_response_timer_id);
  /*
   * Insert the timer in the MAP of mme_ue_s1ap_id <-> timer_id
   */
  //     s1ap_timer_insert(ue_ref->mme_ue_s1ap_id, ue_ref->outcome_response_timer_id);
  message.procedureCode = S1ap_ProcedureCode_id_InitialContextSetup;
  message.direction = S1AP_PDU_PR_initiatingMessage;
  initialContextSetupRequest_p =
    &message.msg.s1ap_InitialContextSetupRequestIEs;
  initialContextSetupRequest_p->mme_ue_s1ap_id =
    (unsigned long) ue_ref->mme_ue_s1ap_id;
  initialContextSetupRequest_p->eNB_UE_S1AP_ID =
    (unsigned long) ue_ref->enb_ue_s1ap_id;

  /*
   * Only add capability information if it's not empty.
   */
  if (conn_est_cnf_pP->ue_radio_capability) {
    OAILOG_DEBUG(LOG_S1AP, "UE radio capability found, adding to message\n");
    initialContextSetupRequest_p->presenceMask |=
      S1AP_INITIALCONTEXTSETUPREQUESTIES_UERADIOCAPABILITY_PRESENT;
    OCTET_STRING_fromBuf(
      &initialContextSetupRequest_p->ueRadioCapability,
      (const char *) conn_est_cnf_pP->ue_radio_capability->data,
      conn_est_cnf_pP->ue_radio_capability->slen);
  }

  /*
   * uEaggregateMaximumBitrateDL and uEaggregateMaximumBitrateUL expressed in term of bits/sec
   */
  asn_uint642INTEGER(
    &initialContextSetupRequest_p->uEaggregateMaximumBitrate
       .uEaggregateMaximumBitRateDL,
    conn_est_cnf_pP->ue_ambr.br_dl);
  asn_uint642INTEGER(
    &initialContextSetupRequest_p->uEaggregateMaximumBitrate
       .uEaggregateMaximumBitRateUL,
    conn_est_cnf_pP->ue_ambr.br_ul);

  for (int item = 0; item < conn_est_cnf_pP->no_of_e_rabs; item++) {
    // Free happens in free_s1ap_initialcontextsetuprequest
    e_RABToBeSetup[item] = calloc(1, sizeof *e_RABToBeSetup[item]);
    memset((void *) e_RABToBeSetup[item], 0, sizeof(*e_RABToBeSetup[item]));
    e_RABToBeSetup[item]->e_RAB_ID = conn_est_cnf_pP->e_rab_id[item]; //5;
    e_RABToBeSetup[item]->e_RABlevelQoSParameters.qCI =
      conn_est_cnf_pP->e_rab_level_qos_qci[item];

    if (conn_est_cnf_pP->nas_pdu[item] != NULL) {
      // NAS PDU is optional in rab_setup
      nas_pdu = calloc(1, sizeof *nas_pdu);
      nas_pdu->size = blength(conn_est_cnf_pP->nas_pdu[item]);
      nas_pdu->buf = malloc(blength(conn_est_cnf_pP->nas_pdu[item]));
      memcpy(
        nas_pdu->buf,
        (void *) conn_est_cnf_pP->nas_pdu[item]->data,
        blength(conn_est_cnf_pP->nas_pdu[item]));
      e_RABToBeSetup[item]->nAS_PDU = nas_pdu;
    }

    e_RABToBeSetup[item]
      ->e_RABlevelQoSParameters.allocationRetentionPriority.priorityLevel =
      conn_est_cnf_pP->e_rab_level_qos_priority_level[item];
    e_RABToBeSetup[item]
      ->e_RABlevelQoSParameters.allocationRetentionPriority
      .pre_emptionCapability =
      conn_est_cnf_pP->e_rab_level_qos_preemption_capability[item];
    e_RABToBeSetup[item]
      ->e_RABlevelQoSParameters.allocationRetentionPriority
      .pre_emptionVulnerability =
      conn_est_cnf_pP->e_rab_level_qos_preemption_vulnerability[item];
    /*
     * Set the GTP-TEID. This is the S1-U S-GW TEID
     */
    INT32_TO_OCTET_STRING(
      conn_est_cnf_pP->gtp_teid[item], &(e_RABToBeSetup[item]->gTP_TEID));
    // S-GW IP address(es) for user-plane
    e_RABToBeSetup[item]->transportLayerAddress.buf = calloc(
      blength(conn_est_cnf_pP->transport_layer_address[item]), sizeof(uint8_t));
    memcpy(
      e_RABToBeSetup[item]->transportLayerAddress.buf,
      conn_est_cnf_pP->transport_layer_address[item]->data,
      blength(conn_est_cnf_pP->transport_layer_address[item]));
    e_RABToBeSetup[item]->transportLayerAddress.size =
      blength(conn_est_cnf_pP->transport_layer_address[item]);
    e_RABToBeSetup[item]->transportLayerAddress.bits_unused = 0;
    ASN_SEQUENCE_ADD(
      &initialContextSetupRequest_p->e_RABToBeSetupListCtxtSUReq,
      e_RABToBeSetup[item]);
  }

  if (
    (conn_est_cnf_pP->presencemask & S1AP_CSFB_INDICATOR_PRESENT) ==
    S1AP_CSFB_INDICATOR_PRESENT) {
    initialContextSetupRequest_p->presenceMask |=
      S1AP_INITIALCONTEXTSETUPREQUESTIES_CSFALLBACKINDICATOR_PRESENT;
    initialContextSetupRequest_p->csFallbackIndicator =
      conn_est_cnf_pP->cs_fallback_indicator;
  }

  initialContextSetupRequest_p->ueSecurityCapabilities.encryptionAlgorithms
    .buf = calloc(2, sizeof(uint8_t));
  memcpy(
    initialContextSetupRequest_p->ueSecurityCapabilities.encryptionAlgorithms
      .buf,
    (uint8_t *) &conn_est_cnf_pP
      ->ue_security_capabilities_encryption_algorithms,
    2);
  initialContextSetupRequest_p->ueSecurityCapabilities.encryptionAlgorithms
    .size = 2;
  initialContextSetupRequest_p->ueSecurityCapabilities.encryptionAlgorithms
    .bits_unused = 0;
  initialContextSetupRequest_p->ueSecurityCapabilities
    .integrityProtectionAlgorithms.buf = calloc(2, sizeof(uint8_t));
  memcpy(
    initialContextSetupRequest_p->ueSecurityCapabilities
      .integrityProtectionAlgorithms.buf,
    (uint8_t *) &conn_est_cnf_pP->ue_security_capabilities_integrity_algorithms,
    2);
  initialContextSetupRequest_p->ueSecurityCapabilities
    .integrityProtectionAlgorithms.size = 2;
  initialContextSetupRequest_p->ueSecurityCapabilities
    .integrityProtectionAlgorithms.bits_unused = 0;
  OAILOG_DEBUG(
    LOG_S1AP,
    "security_capabilities_encryption_algorithms 0x%04X\n",
    conn_est_cnf_pP->ue_security_capabilities_encryption_algorithms);
  OAILOG_DEBUG(
    LOG_S1AP,
    "security_capabilities_integrity_algorithms 0x%04X\n",
    conn_est_cnf_pP->ue_security_capabilities_integrity_algorithms);

  if (conn_est_cnf_pP->kenb) {
    initialContextSetupRequest_p->securityKey.buf =
      calloc(AUTH_KENB_SIZE, sizeof(uint8_t));
    memcpy(
      initialContextSetupRequest_p->securityKey.buf,
      conn_est_cnf_pP->kenb,
      AUTH_KENB_SIZE);
    initialContextSetupRequest_p->securityKey.size = AUTH_KENB_SIZE;
  } else {
    OAILOG_DEBUG(LOG_S1AP, "No kenb\n");
    initialContextSetupRequest_p->securityKey.buf = NULL;
    initialContextSetupRequest_p->securityKey.size = 0;
  }

  initialContextSetupRequest_p->securityKey.bits_unused = 0;

  if (s1ap_mme_encode_pdu(&message, &buffer_p, &length) < 0) {
    free_s1ap_initialcontextsetuprequest(initialContextSetupRequest_p);
    // TODO: handle something
    OAILOG_ERROR(
      LOG_S1AP, "Failed to encode initial context setup request message for "
      "ue_id " MME_UE_S1AP_ID_FMT "\n",
      ue_ref->mme_ue_s1ap_id);
    OAILOG_FUNC_OUT(LOG_S1AP);
  }

  OAILOG_NOTICE(
    LOG_S1AP,
    "Send S1AP_INITIAL_CONTEXT_SETUP_REQUEST message MME_UE_S1AP_ID "
    "= " MME_UE_S1AP_ID_FMT " eNB_UE_S1AP_ID = " ENB_UE_S1AP_ID_FMT "\n",
    (mme_ue_s1ap_id_t) initialContextSetupRequest_p->mme_ue_s1ap_id,
    (enb_ue_s1ap_id_t) initialContextSetupRequest_p->eNB_UE_S1AP_ID);
  bstring b = blk2bstr(buffer_p, length);
  free(buffer_p);
  free_s1ap_initialcontextsetuprequest(initialContextSetupRequest_p);
  s1ap_mme_itti_send_sctp_request(
    &b,
    ue_ref->enb->sctp_assoc_id,
    ue_ref->sctp_stream_send,
    ue_ref->mme_ue_s1ap_id);
  OAILOG_FUNC_OUT(LOG_S1AP);
}
//------------------------------------------------------------------------------
void s1ap_handle_mme_ue_id_notification(
  s1ap_state_t *state,
  const itti_mme_app_s1ap_mme_ue_id_notification_t *const notification_p)
{
  OAILOG_FUNC_IN(LOG_S1AP);
  DevAssert(notification_p != NULL);

  sctp_assoc_id_t sctp_assoc_id = notification_p->sctp_assoc_id;
  enb_ue_s1ap_id_t enb_ue_s1ap_id = notification_p->enb_ue_s1ap_id;
  mme_ue_s1ap_id_t mme_ue_s1ap_id = notification_p->mme_ue_s1ap_id;

  enb_description_t *enb_ref = s1ap_state_get_enb(state, sctp_assoc_id);
  if (enb_ref) {
    ue_description_t *ue_ref =
      s1ap_state_get_ue_enbid(enb_ref, enb_ue_s1ap_id);
    if (ue_ref) {
      ue_ref->mme_ue_s1ap_id = mme_ue_s1ap_id;
      hashtable_rc_t h_rc = hashtable_ts_insert(
        &state->mmeid2associd,
        (const hash_key_t) mme_ue_s1ap_id,
        (void *) (uintptr_t) sctp_assoc_id);

      s1ap_imsi_map_t* imsi_map = get_s1ap_imsi_map();
      hashtable_uint64_ts_insert(
        imsi_map->enb_s1ap_mme_ue_id_htbl,
        (const hash_key_t) enb_ue_s1ap_id,
        mme_ue_s1ap_id);
      OAILOG_DEBUG(
        LOG_S1AP,
        "Associated  sctp_assoc_id %d, enb_ue_s1ap_id " ENB_UE_S1AP_ID_FMT
        ", mme_ue_s1ap_id " MME_UE_S1AP_ID_FMT ":%s \n",
        sctp_assoc_id,
        enb_ue_s1ap_id,
        mme_ue_s1ap_id,
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
  s1ap_state_t *state,
  itti_s1ap_e_rab_rel_cmd_t *const e_rab_rel_cmd)
{
  OAILOG_FUNC_IN(LOG_S1AP);
  ue_description_t *ue_ref = NULL;
  uint8_t *buffer_p = NULL;
  uint32_t length = 0;
  void *id = NULL;
  const enb_ue_s1ap_id_t enb_ue_s1ap_id = e_rab_rel_cmd->enb_ue_s1ap_id;
  const mme_ue_s1ap_id_t ue_id = e_rab_rel_cmd->mme_ue_s1ap_id;

  hashtable_ts_get(
    &state->mmeid2associd, (const hash_key_t) ue_id, (void **) &id);
  if (id) {
    sctp_assoc_id_t sctp_assoc_id = (sctp_assoc_id_t)(uintptr_t) id;
    enb_description_t *enb_ref = s1ap_state_get_enb(state, sctp_assoc_id);
    if (enb_ref) {
      ue_ref = s1ap_state_get_ue_enbid(enb_ref, enb_ue_s1ap_id);
    }
  }
  if (!ue_ref) {
    ue_ref = s1ap_state_get_ue_mmeid(state, ue_id);
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
    S1ap_E_RABReleaseCommandIEs_t *e_rabreleasecmdies = NULL;
    s1ap_message message = {0};

    message.procedureCode = S1ap_ProcedureCode_id_E_RABRelease;
    message.direction = S1AP_PDU_PR_initiatingMessage;
    ue_ref->s1_ue_state = S1AP_UE_CONNECTED;
    e_rabreleasecmdies = &message.msg.s1ap_E_RABReleaseCommandIEs;
    /*
     * Setting UE information with the ones found in ue_ref
     */
    e_rabreleasecmdies->mme_ue_s1ap_id = ue_ref->mme_ue_s1ap_id;
    e_rabreleasecmdies->eNB_UE_S1AP_ID = ue_ref->enb_ue_s1ap_id;
    // e_rabreleasecmdies->uEaggregateMaximumBitrate = NULL;
    /*
     * Fill in the NAS pdu
     */
    e_rabreleasecmdies->presenceMask |=
      S1AP_E_RABRELEASECOMMANDIES_NAS_PDU_PRESENT;

    OCTET_STRING_fromBuf(
      &e_rabreleasecmdies->nas_pdu,
      (char *) bdata(e_rab_rel_cmd->nas_pdu),
      blength(e_rab_rel_cmd->nas_pdu));


    S1ap_E_RABItem_t s1ap_E_RABItemIEs
            [e_rab_rel_cmd->e_rab_to_be_rel_list.no_of_items];
    for (int i = 0; i < e_rab_rel_cmd->e_rab_to_be_rel_list.no_of_items;
         i++) {
      memset(
        &s1ap_E_RABItemIEs[i],
        0,
        sizeof(S1ap_E_RABItem_t));

      s1ap_E_RABItemIEs[i].e_RAB_ID =
        e_rab_rel_cmd->e_rab_to_be_rel_list.item[i].e_rab_id;
      s1ap_mme_set_cause(
        &s1ap_E_RABItemIEs[i].cause,
        S1ap_Cause_PR_radioNetwork,
        S1ap_CauseRadioNetwork_unspecified);

      ASN_SEQUENCE_ADD(
        &e_rabreleasecmdies->e_RABToBeReleasedList,
        &s1ap_E_RABItemIEs[i]);
    }

    if (s1ap_mme_encode_pdu(&message, &buffer_p, &length) < 0) {
      OAILOG_ERROR(LOG_S1AP, "Encoding of  failed \n");
      OAILOG_FUNC_RETURN(LOG_S1AP, RETURNerror);
    }

    OAILOG_NOTICE(
      LOG_S1AP,
      "Send S1AP E_RABRelease Command message MME_UE_S1AP_ID = "
      MME_UE_S1AP_ID_FMT
      " eNB_UE_S1AP_ID = " ENB_UE_S1AP_ID_FMT "\n",
      (mme_ue_s1ap_id_t) e_rabreleasecmdies->mme_ue_s1ap_id,
      (enb_ue_s1ap_id_t) e_rabreleasecmdies->eNB_UE_S1AP_ID);
    bstring b = blk2bstr(buffer_p, length);
    s1ap_mme_itti_send_sctp_request(
      &b,
      ue_ref->enb->sctp_assoc_id,
      ue_ref->sctp_stream_send,
      ue_ref->mme_ue_s1ap_id);
  }

  OAILOG_FUNC_RETURN(LOG_S1AP, RETURNok);
}

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

/*! \file mme_app_bearer.c
  \brief
  \author Sebastien ROUX, Lionel Gauthier
  \company Eurecom
  \email: lionel.gauthier@eurecom.fr
*/
#include <stdio.h>
#include <string.h>
#include <stdlib.h>
#include <stdbool.h>
#include <stdint.h>
#include <3gpp_29.274.h>
#include <inttypes.h>
#include <netinet/in.h>

#include "bstrlib.h"
#include "dynamic_memory_check.h"
#include "log.h"
#include "conversions.h"
#include "common_types.h"
#include "intertask_interface.h"
#include "mme_config.h"
#include "mme_app_ue_context.h"
#include "mme_app_defs.h"
#include "mme_app_bearer_context.h"
#include "sgw_ie_defs.h"
#include "common_defs.h"
#include "mme_app_itti_messaging.h"
#include "mme_app_procedures.h"
#include "mme_app_statistics.h"
#include "timer.h"
#include "nas_proc.h"
#include "3gpp_23.003.h"
#include "3gpp_24.007.h"
#include "3gpp_24.008.h"
#include "3gpp_24.301.h"
#include "3gpp_36.401.h"
#include "3gpp_36.413.h"
#include "CsfbResponse.h"
#include "ServiceType.h"
#include "TrackingAreaIdentity.h"
#include "nas/as_message.h"
#include "emm_data.h"
#include "esm_data.h"
#include "hashtable.h"
#include "intertask_interface_types.h"
#include "itti_types.h"
#include "mme_api.h"
#include "mme_app_state.h"
#include "mme_app_messages_types.h"
#include "s11_messages_types.h"
#include "s1ap_messages_types.h"
#include "nas/securityDef.h"
#include "service303.h"
#include "sgs_messages_types.h"
#include "secu_defs.h"

#if EMBEDDED_SGW
#define TASK_SPGW TASK_SPGW_APP
#else
#define TASK_SPGW TASK_S11
#endif

int send_modify_bearer_req(mme_ue_s1ap_id_t ue_id, ebi_t ebi)
{
  OAILOG_FUNC_IN(LOG_MME_APP);

  uint8_t item = 0; // This function call is used for default bearer only
  ue_mm_context_t* ue_context_p = NULL;

  ue_context_p = mme_ue_context_exists_mme_ue_s1ap_id(ue_id);
  if (ue_context_p == NULL) {
    OAILOG_ERROR(
      LOG_MME_APP,
      "ue_context_p is NULL, did not send S11_MODIFY_BEARER_REQUEST\n");
    OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNerror);
  }

  bearer_context_t* bearer_cntxt =
    mme_app_get_bearer_context(ue_context_p, ebi);
  if (bearer_cntxt == NULL) {
    OAILOG_ERROR(
      LOG_MME_APP,
      "Bearer context is null, did not send S11_MODIFY_BEARER_REQUEST for ebi"
      "%u\n",
      ebi);
    OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNerror);
  }

  pdn_cid_t cid = ue_context_p->bearer_contexts[EBI_TO_INDEX(ebi)]->pdn_cx_id;
  pdn_context_t* pdn_context_p = ue_context_p->pdn_contexts[cid];
  if (pdn_context_p == NULL) {
    OAILOG_ERROR(LOG_MME_APP, "Did not find PDN context for ebi %u\n", ebi);
    OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNerror);
  }

  MessageDef* message_p =
    itti_alloc_new_message(TASK_MME_APP, S11_MODIFY_BEARER_REQUEST);
  if (message_p == NULL) {
    OAILOG_ERROR(
      LOG_MME_APP, "Cannot allocate memory to S11_MODIFY_BEARER_REQUEST\n");
    OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNerror);
  }

  itti_s11_modify_bearer_request_t* s11_modify_bearer_request =
    &message_p->ittiMsg.s11_modify_bearer_request;
  s11_modify_bearer_request->local_teid = ue_context_p->mme_teid_s11;
  /*
   * Delay Value in integer multiples of 50 millisecs, or zero
   */
  s11_modify_bearer_request->delay_dl_packet_notif_req = 0;
  s11_modify_bearer_request->bearer_contexts_to_be_modified
    .bearer_contexts[item]
    .eps_bearer_id = ebi;
  s11_modify_bearer_request->bearer_contexts_to_be_modified
    .bearer_contexts[item]
    .s1_eNB_fteid.teid = bearer_cntxt->enb_fteid_s1u.teid;
  s11_modify_bearer_request->bearer_contexts_to_be_modified
    .bearer_contexts[item]
    .s1_eNB_fteid.interface_type = S1_U_ENODEB_GTP_U;

//memcpy(
  //    (void *)&s11_modify_bearer_request->edns_peer_ip,
    //  (struct sockaddr *)&pdn_context_p->s_gw_address_s11_s4,
      //((struct sockaddr *)&pdn_context_p->s_gw_address_s11_s4)->sa_family == AF_INET
        //  ? sizeof(struct sockaddr_in)
          //: sizeof(struct sockaddr_in6));
  
  //ip_address_t s_gw_address_s11_s4;

  s11_modify_bearer_request->edns_peer_ip.addr_v4_.sin_addr =
  pdn_context_p->s_gw_address_s11_s4.address.ipv4_address;

  //OAILOG_DEBUG (LOG_MME_APP, "SGW  IP information %p\n", &s_gw_address_s11_s4.address.ipv4_address);  
  s11_modify_bearer_request->teid = pdn_context_p->s_gw_teid_s11_s4;

  if (bearer_cntxt->enb_fteid_s1u.ipv4) {
    s11_modify_bearer_request->bearer_contexts_to_be_modified
      .bearer_contexts[item]
      .s1_eNB_fteid.ipv4 = 1;
    memcpy(
      &s11_modify_bearer_request->bearer_contexts_to_be_modified
         .bearer_contexts[item]
         .s1_eNB_fteid.ipv4_address,
      &bearer_cntxt->enb_fteid_s1u.ipv4_address,
      sizeof(bearer_cntxt->enb_fteid_s1u.ipv4_address));
  } else if (bearer_cntxt->enb_fteid_s1u.ipv6) {
    s11_modify_bearer_request->bearer_contexts_to_be_modified
      .bearer_contexts[item]
      .s1_eNB_fteid.ipv6 = 1;
    memcpy(
      &s11_modify_bearer_request->bearer_contexts_to_be_modified
         .bearer_contexts[item]
         .s1_eNB_fteid.ipv6_address,
      &bearer_cntxt->enb_fteid_s1u.ipv6_address,
      sizeof(bearer_cntxt->enb_fteid_s1u.ipv6_address));
  } else {
    OAILOG_ERROR(LOG_MME_APP, "Unknown IP address\n");
  }

  // Only one bearer context to be sent for secondary PDN
  s11_modify_bearer_request->bearer_contexts_to_be_modified.num_bearer_context =
    1;
  s11_modify_bearer_request->bearer_contexts_to_be_removed.num_bearer_context =
    0;
  s11_modify_bearer_request->mme_fq_csid.node_id_type = GLOBAL_UNICAST_IPv4;
  s11_modify_bearer_request->mme_fq_csid.csid = 0;
  memset(
    &s11_modify_bearer_request->indication_flags,
    0,
    sizeof(s11_modify_bearer_request->indication_flags));
  s11_modify_bearer_request->rat_type = RAT_EUTRAN;

  message_p->ittiMsgHeader.imsi = ue_context_p->emm_context._imsi64;

  OAILOG_INFO(
    LOG_MME_APP,
    "Sending S11_MODIFY_BEARER_REQUEST to SGW for ue" MME_UE_S1AP_ID_FMT "\n",
    ue_id);
  itti_send_msg_to_task(TASK_SPGW, INSTANCE_DEFAULT, message_p);
  OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNok);
}


void print_bearer_ids_helper(const ebi_t *ebi, uint32_t no_of_bearers)
{
  OAILOG_FUNC_IN(LOG_MME_APP);

  char buf[128], *pos = buf;
  for (int i = 0 ; i != no_of_bearers ; i++) {
    if (i) {
      pos += sprintf(pos, ", ");
    }
    pos += sprintf(pos, "%d", ebi[i]);
  }
  OAILOG_INFO(LOG_MME_APP, " EBIs in the list %s\n", buf);
  OAILOG_FUNC_OUT(LOG_MME_APP);
}

//------------------------------------------------------------------------------
int _send_pcrf_bearer_actv_rsp(
  struct ue_mm_context_s *ue_context_p,
  ebi_t ebi,
  gtpv2c_cause_value_t cause)
{
  OAILOG_FUNC_IN(LOG_MME_APP);
  MessageDef *message_p =
    itti_alloc_new_message(TASK_MME_APP, S11_NW_INITIATED_ACTIVATE_BEARER_RESP);
  if (message_p == NULL) {
    OAILOG_ERROR(
      LOG_MME_APP,
      "Cannot allocte memory to S11_NW_INITIATED_BEARER_ACTV_RSP\n");
    OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNerror);
  }
  itti_s11_nw_init_actv_bearer_rsp_t *s11_nw_init_actv_bearer_rsp =
    &message_p->ittiMsg.s11_nw_init_actv_bearer_rsp;

  //Fetch PDN context
  pdn_cid_t cid = ue_context_p->bearer_contexts[EBI_TO_INDEX(ebi)]->pdn_cx_id;
  pdn_context_t *pdn_context = ue_context_p->pdn_contexts[cid];
  //Fill SGW S11 CP TEID
  s11_nw_init_actv_bearer_rsp->sgw_s11_teid = pdn_context->s_gw_teid_s11_s4;
  int msg_bearer_index = 0;

  bearer_context_t *bc = mme_app_get_bearer_context(ue_context_p, ebi);
  s11_nw_init_actv_bearer_rsp->cause.cause_value = cause;
  s11_nw_init_actv_bearer_rsp->bearer_contexts.bearer_contexts[msg_bearer_index]
    .eps_bearer_id = ebi;
  s11_nw_init_actv_bearer_rsp->bearer_contexts.bearer_contexts[msg_bearer_index]
    .cause.cause_value = REQUEST_ACCEPTED;
  //  FTEID eNB
  s11_nw_init_actv_bearer_rsp->bearer_contexts.bearer_contexts[msg_bearer_index]
    .s1u_enb_fteid = bc->enb_fteid_s1u;

  /* FTEID SGW S1U
   * This IE shall be sent on the S11 interface.
  It shall be used to fetch context*/
  s11_nw_init_actv_bearer_rsp->bearer_contexts.bearer_contexts[msg_bearer_index]
    .s1u_sgw_fteid = bc->s_gw_fteid_s1u;
  s11_nw_init_actv_bearer_rsp->bearer_contexts.num_bearer_context++;

  message_p->ittiMsgHeader.imsi = ue_context_p->emm_context._imsi64;

  OAILOG_INFO(
    LOG_MME_APP,
    "Sending create_dedicated_bearer_rsp to SGW with EBI %u s1u teid %u\n",
    ebi,
    bc->s_gw_fteid_s1u.teid);
  itti_send_msg_to_task(TASK_SPGW, INSTANCE_DEFAULT, message_p);
  OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNok);
}

//---------------------------------------------------------------------------
static bool mme_app_construct_guti(
  const plmn_t *const plmn_p,
  const s_tmsi_t *const s_tmsi_p,
  guti_t *const guti_p);
static void notify_s1ap_new_ue_mme_s1ap_id_association(
  struct ue_mm_context_s *ue_context_p);

//------------------------------------------------------------------------------
void mme_app_handle_conn_est_cnf(nas_establish_rsp_t* const nas_conn_est_cnf_p)
{
  OAILOG_FUNC_IN(LOG_MME_APP);
  struct ue_mm_context_s* ue_context_p = NULL;
  emm_context_t emm_context = {0};
  MessageDef* message_p = NULL;
  itti_mme_app_connection_establishment_cnf_t* establishment_cnf_p = NULL;
  int rc = RETURNok;

  OAILOG_DEBUG(
    LOG_MME_APP,
    "Handle MME_APP_CONNECTION_ESTABLISHMENT_CNF for ue-id " MME_UE_S1AP_ID_FMT
    "\n",
    nas_conn_est_cnf_p->ue_id);
  ue_context_p = mme_ue_context_exists_mme_ue_s1ap_id(nas_conn_est_cnf_p->ue_id);

  if (!ue_context_p) {
    OAILOG_ERROR(
      LOG_MME_APP,
      "UE context doesn't exist for UE " MME_UE_S1AP_ID_FMT "\n",
      nas_conn_est_cnf_p->ue_id);
    // memory leak
    bdestroy_wrapper(&nas_conn_est_cnf_p->nas_msg);
    OAILOG_FUNC_OUT(LOG_MME_APP);
  }
  emm_context = ue_context_p->emm_context;
  /* Check that if Service Request is recieved in response to SGS Paging for MT SMS */
  if (ue_context_p->sgs_context) {
    /*
     * Move the UE to ECM Connected State.
     */
    /*
    * Check that if SGS paging is recieved without LAI then
    * send IMSI Detach towads UE to re-attach for non-eps services
    * otherwise send itti SGS Service request message to SGS
    */
  OAILOG_DEBUG(
    LOG_MME_APP,
    "CSFB Service Type = (%d) for (ue_id = " MME_UE_S1AP_ID_FMT ")\n",
    ue_context_p->sgs_context->csfb_service_type,
    nas_conn_est_cnf_p->ue_id);
    if (ue_context_p->sgs_context->csfb_service_type == CSFB_SERVICE_MT_SMS) {
      /* send SGS SERVICE request message to SGS */
      if (
        RETURNok !=
        (rc = mme_app_send_sgsap_service_request(
           ue_context_p->sgs_context->service_indicator, ue_context_p))) {
        OAILOG_ERROR(
          LOG_MME_APP,
          "Failed to send CS-Service Request to SGS-Task for (ue_id = %u) \n",
          ue_context_p->mme_ue_s1ap_id);
      }
    } else if (
      ue_context_p->sgs_context->csfb_service_type ==
      CSFB_SERVICE_MT_CALL_OR_SMS_WITHOUT_LAI) {
      // Inform NAS module to send network initiated IMSI detach request to UE
      OAILOG_DEBUG(
        LOG_MME_APP,
        "Send SGS intiated Detach request to NAS module for ue_id = "
        MME_UE_S1AP_ID_FMT "\n"
        "csfb service type = CSFB_SERVICE_MT_CALL_OR_SMS_WITHOUT_LAI\n",
        ue_context_p->mme_ue_s1ap_id);

      mme_app_handle_nw_initiated_detach_request(
        ue_context_p->mme_ue_s1ap_id, SGS_INITIATED_IMSI_DETACH);
      ue_context_p->sgs_context->csfb_service_type = CSFB_SERVICE_NONE;
      OAILOG_FUNC_OUT(LOG_MME_APP);
    }
  }

  if (
    (((nas_conn_est_cnf_p->presencemask) & SERVICE_TYPE_PRESENT) ==
     SERVICE_TYPE_PRESENT) &&
   ((nas_conn_est_cnf_p->service_type == MO_CS_FB) ||
    (nas_conn_est_cnf_p->service_type == MO_CS_FB_EMRGNCY_CALL))) {
    if (ue_context_p->sgs_context != NULL) {
      ue_context_p->sgs_context->csfb_service_type = CSFB_SERVICE_MO_CALL;
    } else {
      OAILOG_ERROR(
        LOG_MME_APP,
        "SGS context doesn't exist for UE" MME_UE_S1AP_ID_FMT "\n",
        nas_conn_est_cnf_p->ue_id);
      mme_app_notify_service_reject_to_nas(
        ue_context_p->mme_ue_s1ap_id,
        EMM_CAUSE_CONGESTION,
        INTIAL_CONTEXT_SETUP_PROCEDURE_FAILED);
      OAILOG_FUNC_OUT(LOG_MME_APP);
    }
    if(nas_conn_est_cnf_p->service_type == MO_CS_FB_EMRGNCY_CALL){
     ue_context_p->sgs_context->is_emergency_call = true;
    }
  }
  if (
    (ue_context_p->sgs_context) &&
    (ue_context_p->sgs_context->csfb_service_type == CSFB_SERVICE_MT_CALL)) {
    if (nas_conn_est_cnf_p->csfb_response == CSFB_REJECTED_BY_UE) {
      /* CSFB MT calll rejected by user, send sgsap-paging reject to VLR */
      if (
        (rc = mme_app_send_sgsap_paging_reject(
           ue_context_p,
           emm_context._imsi64,
           emm_context._imsi.length,
           SGS_CAUSE_MT_CSFB_CALL_REJECTED_BY_USER)) != RETURNok) {
        OAILOG_WARNING(
          LOG_MME_APP,
          "Failed to send SGSAP-Paging Reject for imsi with reject cause:"
          "SGS_CAUSE_MT_CSFB_CALL_REJECTED_BY_USER" IMSI_64_FMT "\n",
          emm_context._imsi64);
      }
      OAILOG_FUNC_OUT(LOG_MME_APP);
    }
  }
  message_p =
    itti_alloc_new_message(TASK_MME_APP, MME_APP_CONNECTION_ESTABLISHMENT_CNF);
  establishment_cnf_p =
    &message_p->ittiMsg.mme_app_connection_establishment_cnf;

  establishment_cnf_p->ue_id = nas_conn_est_cnf_p->ue_id;

  if (
    (ue_context_p->sgs_context != NULL) &&
    ((ue_context_p->sgs_context->csfb_service_type == CSFB_SERVICE_MT_CALL) ||
     (ue_context_p->sgs_context->csfb_service_type == CSFB_SERVICE_MO_CALL))) {
    establishment_cnf_p->presencemask |= S1AP_CSFB_INDICATOR_PRESENT;
    if (ue_context_p->sgs_context->is_emergency_call == true){
      establishment_cnf_p->cs_fallback_indicator   = CSFB_HIGH_PRIORITY;
      ue_context_p->sgs_context->is_emergency_call = false;
    } else {
      establishment_cnf_p->cs_fallback_indicator      = CSFB_REQUIRED;
    }
  }
  OAILOG_DEBUG(
    LOG_MME_APP,
    "CSFB Fallback indicator = (%d)\n",establishment_cnf_p->cs_fallback_indicator );
  // Copy UE radio capabilities into message if it exists
  OAILOG_DEBUG(
    LOG_MME_APP,
    "UE radio context already cached: %s\n",
    ue_context_p->ue_radio_capability ? "yes" : "no");
  if (ue_context_p->ue_radio_capability) {
    establishment_cnf_p->ue_radio_capability =
      bstrcpy(ue_context_p->ue_radio_capability);
  }

  int j = 0;
  for (int i = 0; i < BEARERS_PER_UE; i++) {
    bearer_context_t *bc = ue_context_p->bearer_contexts[i];
    if (bc) {
      if (BEARER_STATE_SGW_CREATED & bc->bearer_state) {
        establishment_cnf_p->e_rab_id[j] =
          bc->ebi; //+ EPS_BEARER_IDENTITY_FIRST;
        establishment_cnf_p->e_rab_level_qos_qci[j] = bc->qci;
        establishment_cnf_p->e_rab_level_qos_priority_level[j] =
          bc->priority_level;
        establishment_cnf_p->e_rab_level_qos_preemption_capability[j] =
          bc->preemption_capability;
        establishment_cnf_p->e_rab_level_qos_preemption_vulnerability[j] =
          bc->preemption_vulnerability;
        establishment_cnf_p->transport_layer_address[j] =
          fteid_ip_address_to_bstring(&bc->s_gw_fteid_s1u);
        establishment_cnf_p->gtp_teid[j] = bc->s_gw_fteid_s1u.teid;
        if (!j) {
          establishment_cnf_p->nas_pdu[j] = nas_conn_est_cnf_p->nas_msg;
          nas_conn_est_cnf_p->nas_msg = NULL;
#if DEBUG_IS_ON
          if (!establishment_cnf_p->nas_pdu[j]) {
            OAILOG_ERROR(
              LOG_MME_APP,
              "No NAS PDU found ue " MME_UE_S1AP_ID_FMT "\n",
              nas_conn_est_cnf_p->ue_id);
          }
#endif
        }
        j = j + 1;
      }
    }
  }
  establishment_cnf_p->no_of_e_rabs = j;

  //#pragma message  "Check ue_context_p ambr"
  establishment_cnf_p->ue_ambr.br_ul = ue_context_p->subscribed_ue_ambr.br_ul;
  establishment_cnf_p->ue_ambr.br_dl = ue_context_p->subscribed_ue_ambr.br_dl;
  establishment_cnf_p->ue_security_capabilities_encryption_algorithms =
    ((uint16_t) emm_context._ue_network_capability.eea &
    ~(1 << 7)) << 1;

  establishment_cnf_p->ue_security_capabilities_integrity_algorithms =
    ((uint16_t) emm_context._ue_network_capability.eia &
    ~(1 << 7)) << 1;

  if (!((0 <= emm_context._security.vector_index) &&
    (MAX_EPS_AUTH_VECTORS > emm_context._security.vector_index))) {
    OAILOG_ERROR(
      LOG_MME_APP,
      "Invalid security vector index %d",
      emm_context._security.vector_index);
    OAILOG_FUNC_OUT(LOG_MME_APP);
  }

  derive_keNB(
    emm_context._vector[emm_context._security.vector_index].kasme,
    emm_context._security.kenb_ul_count.seq_num |
    (emm_context._security.kenb_ul_count.overflow << 8),
    establishment_cnf_p->kenb);

  /* Genarate Next HOP key parameter */
  derive_NH(
    emm_context._vector[emm_context._security.vector_index].kasme,
    establishment_cnf_p->kenb,
    emm_context._security.next_hop,
    &emm_context._security.next_hop_chaining_count);

  OAILOG_DEBUG(
    LOG_MME_APP,
    "security_capabilities_encryption_algorithms 0x%04X\n",
    establishment_cnf_p->ue_security_capabilities_encryption_algorithms);
  OAILOG_DEBUG(
    LOG_MME_APP,
    "security_capabilities_integrity_algorithms  0x%04X\n",
    establishment_cnf_p->ue_security_capabilities_integrity_algorithms);

  message_p->ittiMsgHeader.imsi = ue_context_p->emm_context._imsi64;
  itti_send_msg_to_task(TASK_S1AP, INSTANCE_DEFAULT, message_p);

  /*
   * Move the UE to ECM Connected State.However if S1-U bearer establishment fails then we need to move the UE to idle.
   * S1 Signaling connection gets established via first DL NAS Trasnport message in some scenarios so check the state
   * first
   */
  if (ue_context_p->ecm_state != ECM_CONNECTED) {
    mme_app_desc_t* mme_app_desc_p = get_mme_nas_state(false);
    mme_ue_context_update_ue_sig_connection_state(
      &mme_app_desc_p->mme_ue_contexts, ue_context_p, ECM_CONNECTED);

    if (
      (ue_context_p->sgs_context) &&
      (ue_context_p->sgs_context->csfb_service_type == CSFB_SERVICE_MT_CALL)) {
      /* send sgsap-Service Request to VLR */
      if (
        RETURNok !=
        (rc = mme_app_send_sgsap_service_request(
           ue_context_p->sgs_context->service_indicator, ue_context_p))) {
        OAILOG_ERROR(
          LOG_MME_APP,
          "Failed to send CS-Service Request to SGS-Task for ue-id:"
           MME_UE_S1AP_ID_FMT "\n", ue_context_p->mme_ue_s1ap_id);
      }
    }
  }

  /* Start timer to wait for Initial UE Context Response from eNB
   * If timer expires treat this as failure of ongoing procedure and abort corresponding NAS procedure such as ATTACH
   * or SERVICE REQUEST. Send UE context release command to eNB
   */
  nas_itti_timer_arg_t timer_callback_fun = {0};
  timer_callback_fun.nas_timer_callback =
    mme_app_handle_initial_context_setup_rsp_timer_expiry;
  timer_callback_fun.nas_timer_callback_arg =
    (void *) &(ue_context_p->mme_ue_s1ap_id);
  if (timer_setup(
    ue_context_p->initial_context_setup_rsp_timer.sec,
    0,
    TASK_MME_APP,
    INSTANCE_DEFAULT,
    TIMER_ONE_SHOT,
    &timer_callback_fun,
    sizeof(timer_callback_fun),
    &(ue_context_p->initial_context_setup_rsp_timer.id)) < 0) {
    OAILOG_ERROR(
      LOG_MME_APP,
      "Failed to start initial context setup response timer for UE id  %d \n",
      ue_context_p->mme_ue_s1ap_id);
    ue_context_p->initial_context_setup_rsp_timer.id =
      MME_APP_TIMER_INACTIVE_ID;
  } else {
    OAILOG_INFO(
      LOG_MME_APP,
      "MME APP : Sent Initial context Setup Request and Started guard timer "
      "for UE id  %d timer_id :%lx \n",
      ue_context_p->mme_ue_s1ap_id, (long)ue_context_p->initial_context_setup_rsp_timer.id);
  }
  OAILOG_FUNC_OUT(LOG_MME_APP);
}

// sent by S1AP
//------------------------------------------------------------------------------
void mme_app_handle_initial_ue_message(mme_app_desc_t *mme_app_desc_p,
  itti_s1ap_initial_ue_message_t *const initial_pP)
{
  OAILOG_FUNC_IN(LOG_MME_APP);
  struct ue_mm_context_s *ue_context_p = NULL;
  bool is_guti_valid = false;
  bool is_mm_ctx_new = false;
  enb_s1ap_id_key_t enb_s1ap_id_key = INVALID_ENB_UE_S1AP_ID_KEY;

  OAILOG_INFO(LOG_MME_APP, "Received MME_APP_INITIAL_UE_MESSAGE from S1AP\n");

  if (initial_pP->mme_ue_s1ap_id != INVALID_MME_UE_S1AP_ID) {
    OAILOG_ERROR(
      LOG_MME_APP,
      "MME UE S1AP Id (" MME_UE_S1AP_ID_FMT ") is already assigned\n",
      initial_pP->mme_ue_s1ap_id);
    OAILOG_FUNC_OUT(LOG_MME_APP);
  }

  // Check if there is any existing UE context using S-TMSI/GUTI
  if (initial_pP->is_s_tmsi_valid) {
    OAILOG_DEBUG(
      LOG_MME_APP,
      "INITIAL UE Message: Valid mme_code %u and S-TMSI %u received from "
      "eNB.\n",
      initial_pP->opt_s_tmsi.mme_code,
      initial_pP->opt_s_tmsi.m_tmsi);
    guti_t guti = {.gummei.plmn = {0},
                   .gummei.mme_gid = 0,
                   .gummei.mme_code = 0,
                   .m_tmsi = INVALID_M_TMSI};
    plmn_t plmn = {.mcc_digit1 = initial_pP->tai.mcc_digit1,
                   .mcc_digit2 = initial_pP->tai.mcc_digit2,
                   .mcc_digit3 = initial_pP->tai.mcc_digit3,
                   .mnc_digit1 = initial_pP->tai.mnc_digit1,
                   .mnc_digit2 = initial_pP->tai.mnc_digit2,
                   .mnc_digit3 = initial_pP->tai.mnc_digit3};
    is_guti_valid =
      mme_app_construct_guti(&plmn, &(initial_pP->opt_s_tmsi), &guti);
    if (is_guti_valid) {
      ue_context_p =
        mme_ue_context_exists_guti(&mme_app_desc_p->mme_ue_contexts, &guti);
      if (ue_context_p) {
        initial_pP->mme_ue_s1ap_id = ue_context_p->mme_ue_s1ap_id;
        if (ue_context_p->enb_s1ap_id_key != INVALID_ENB_UE_S1AP_ID_KEY) {
          /*
           * Ideally this should never happen. When UE moves to IDLE,
           * this key is set to INVALID.
           * Note - This can happen if eNB detects RLF late and by that time
           * UE sends Initial NAS message via new RRC connection.
           * However if this key is valid, remove the key from the hashtable.
           */

          OAILOG_ERROR(
            LOG_MME_APP,
            "MME_APP_INITAIL_UE_MESSAGE: enb_s1ap_id_key %ld has "
            "valid value \n",
            ue_context_p->enb_s1ap_id_key);
          // Inform s1ap for local cleanup of enb_ue_s1ap_id from ue context
          ue_context_p->ue_context_rel_cause = S1AP_INVALID_ENB_ID;
          OAILOG_ERROR(
            LOG_MME_APP,
            " Sending UE Context Release to S1AP for ue_id =(%u)\n",
            ue_context_p->mme_ue_s1ap_id);
          mme_app_itti_ue_context_release(
            ue_context_p, ue_context_p->ue_context_rel_cause);
          hashtable_uint64_ts_remove(
            mme_app_desc_p->mme_ue_contexts.enb_ue_s1ap_id_ue_context_htbl,
            (const hash_key_t) ue_context_p->enb_s1ap_id_key);
          ue_context_p->enb_s1ap_id_key = INVALID_ENB_UE_S1AP_ID_KEY;
          ue_context_p->ue_context_rel_cause = S1AP_INVALID_CAUSE;
        }
        // Update MME UE context with new enb_ue_s1ap_id
        ue_context_p->enb_ue_s1ap_id = initial_pP->enb_ue_s1ap_id;
        // regenerate the enb_s1ap_id_key as enb_ue_s1ap_id is changed.
        MME_APP_ENB_S1AP_ID_KEY(
          enb_s1ap_id_key, initial_pP->enb_id, initial_pP->enb_ue_s1ap_id);
        // Update enb_s1ap_id_key in hashtable
        mme_ue_context_update_coll_keys(
          &mme_app_desc_p->mme_ue_contexts,
          ue_context_p,
          enb_s1ap_id_key,
          ue_context_p->mme_ue_s1ap_id,
          ue_context_p->emm_context._imsi64,
          ue_context_p->mme_teid_s11,
          &guti);
        // Check if paging timer exists for UE and remove
        if (
          ue_context_p->paging_response_timer.id !=
          MME_APP_TIMER_INACTIVE_ID) {
          if (timer_remove(ue_context_p->paging_response_timer.id, NULL)) {
            OAILOG_ERROR(
              LOG_MME_APP,
              "Failed to stop paging response timer for UE id %d\n",
              ue_context_p->mme_ue_s1ap_id);
          }
          ue_context_p->paging_response_timer.id = MME_APP_TIMER_INACTIVE_ID;
        }
      } else {
        OAILOG_DEBUG(
          LOG_MME_APP,
          "No UE context found for MME code %u and S-TMSI %u\n",
          initial_pP->opt_s_tmsi.mme_code,
          initial_pP->opt_s_tmsi.m_tmsi);
      }
    } else {
      OAILOG_DEBUG(
        LOG_MME_APP,
        "No MME is configured with MME code %u received in S-TMSI %u from "
        "UE.\n",
        initial_pP->opt_s_tmsi.mme_code,
        initial_pP->opt_s_tmsi.m_tmsi);
    }
  } else {
    OAILOG_DEBUG(
      LOG_MME_APP, "MME_APP_INITIAL_UE_MESSAGE from S1AP,without S-TMSI. \n");
  }
  // create a new ue context if nothing is found
  if (!(ue_context_p)) {
    OAILOG_DEBUG(LOG_MME_APP, "UE context doesn't exist -> create one\n");
    if (!(ue_context_p = mme_create_new_ue_context())) {
      /*
       * Error during ue context malloc
       */
      OAILOG_ERROR(LOG_MME_APP, "Failed to create new ue context \n");
      OAILOG_FUNC_OUT(LOG_MME_APP);
    }
    is_mm_ctx_new = true;
    // Allocate new mme_ue_s1ap_id
    ue_context_p->mme_ue_s1ap_id = mme_app_ctx_get_new_ue_id();
    if (ue_context_p->mme_ue_s1ap_id == INVALID_MME_UE_S1AP_ID) {
      OAILOG_CRITICAL(
        LOG_MME_APP,
        "MME_APP_INITIAL_UE_MESSAGE. MME_UE_S1AP_ID allocation Failed.\n");
      mme_remove_ue_context(&mme_app_desc_p->mme_ue_contexts, ue_context_p);
      OAILOG_FUNC_OUT(LOG_MME_APP);
    }
    OAILOG_DEBUG(
      LOG_MME_APP,
      "Allocated new MME UE context and new "
      "(mme_ue_s1ap_id = %d)\n",
      ue_context_p->mme_ue_s1ap_id);
    ue_context_p->enb_ue_s1ap_id = initial_pP->enb_ue_s1ap_id;
    MME_APP_ENB_S1AP_ID_KEY(
      ue_context_p->enb_s1ap_id_key,
      initial_pP->enb_id,
      initial_pP->enb_ue_s1ap_id);

    if (
      mme_insert_ue_context(&mme_app_desc_p->mme_ue_contexts, ue_context_p) !=
      RETURNok) {
      OAILOG_ERROR(
        LOG_MME_APP,
        "Failed to insert UE contxt, MME UE S1AP Id: " MME_UE_S1AP_ID_FMT "\n",
        ue_context_p->mme_ue_s1ap_id);
      OAILOG_FUNC_OUT(LOG_MME_APP);
    }
  }
  ue_context_p->sctp_assoc_id_key = initial_pP->sctp_assoc_id;
  ue_context_p->e_utran_cgi = initial_pP->ecgi;
  // Notify S1AP about the mapping between mme_ue_s1ap_id and
  // sctp assoc id + enb_ue_s1ap_id
  notify_s1ap_new_ue_mme_s1ap_id_association(ue_context_p);
  s_tmsi_t s_tmsi = {0};
  if (initial_pP->is_s_tmsi_valid) {
    s_tmsi = initial_pP->opt_s_tmsi;
  } else {
    s_tmsi.mme_code = 0;
    s_tmsi.m_tmsi = INVALID_M_TMSI;
  }
  OAILOG_INFO(
    LOG_MME_APP,
    "INITIAL_UE_MESSAGE RCVD \n"
    "mme_ue_s1ap_id  = %d\n"
    "enb_ue_s1ap_id  = %d\n",
    ue_context_p->mme_ue_s1ap_id,
    ue_context_p->enb_ue_s1ap_id);
  OAILOG_DEBUG(
    LOG_MME_APP, "Is S-TMSI Valid - (%d)\n",initial_pP->is_s_tmsi_valid);

  OAILOG_INFO(LOG_MME_APP, "Sending NAS Establishment Indication to NAS for ue_id = (%d)\n",
    ue_context_p->mme_ue_s1ap_id );
  nas_proc_establish_ind(
    ue_context_p->mme_ue_s1ap_id,
    is_mm_ctx_new,
    initial_pP->tai,
    initial_pP->ecgi,
    initial_pP->rrc_establishment_cause,
    s_tmsi,
    &initial_pP->nas);
  // s1ap_initial_ue_message_t transparent; may be needed :
  // OLD CODE memcpy (
  //   &message_p->ittiMsg.nas_initial_ue_message.transparent,
  //   (const void*)&initial_pP->transparent,
  //   sizeof (message_p->ittiMsg.nas_initial_ue_message.transparent));

  initial_pP->nas = NULL;

  OAILOG_FUNC_OUT(LOG_MME_APP);
}

//------------------------------------------------------------------------------
void mme_app_handle_erab_setup_req(
  const mme_ue_s1ap_id_t ue_id,
  const ebi_t ebi,
  const bitrate_t mbr_dl,
  const bitrate_t mbr_ul,
  const bitrate_t gbr_dl,
  const bitrate_t gbr_ul,
  bstring nas_msg)
{
  OAILOG_FUNC_IN(LOG_MME_APP);

  OAILOG_DEBUG(
    LOG_MME_APP,
    "Handle mme app e-rab setup request for ue-id " MME_UE_S1AP_ID_FMT "\n",
    ue_id);

  ue_mm_context_t* ue_context_p = mme_ue_context_exists_mme_ue_s1ap_id(ue_id);

  if (!ue_context_p) {
    OAILOG_ERROR(
      LOG_MME_APP,
      "UE context doesn't exist for ue_id " MME_UE_S1AP_ID_FMT "\n",
      ue_id);
    bdestroy_wrapper(&nas_msg);
    OAILOG_FUNC_OUT(LOG_MME_APP);
  }

  bearer_context_t* bearer_context =
    mme_app_get_bearer_context(ue_context_p, ebi);

  if (bearer_context) {
    MessageDef* message_p =
      itti_alloc_new_message(TASK_MME_APP, S1AP_E_RAB_SETUP_REQ);
    if (message_p == NULL) {
      OAILOG_WARNING(
        LOG_MME_APP,
        "Failed to allocate the memory for s1ap erab set request message\n");
      bdestroy_wrapper(&nas_msg);
      OAILOG_FUNC_OUT(LOG_MME_APP);
    }

    itti_s1ap_e_rab_setup_req_t* s1ap_e_rab_setup_req =
      &message_p->ittiMsg.s1ap_e_rab_setup_req;

    s1ap_e_rab_setup_req->mme_ue_s1ap_id = ue_context_p->mme_ue_s1ap_id;
    s1ap_e_rab_setup_req->enb_ue_s1ap_id = ue_context_p->enb_ue_s1ap_id;

    // E-RAB to Be Setup List
    s1ap_e_rab_setup_req->e_rab_to_be_setup_list.no_of_items = 1;
    s1ap_e_rab_setup_req->e_rab_to_be_setup_list.item[0].e_rab_id =
      bearer_context->ebi;
    s1ap_e_rab_setup_req->e_rab_to_be_setup_list.item[0]
      .e_rab_level_qos_parameters.allocation_and_retention_priority
      .pre_emption_capability = bearer_context->preemption_capability;
    s1ap_e_rab_setup_req->e_rab_to_be_setup_list.item[0]
      .e_rab_level_qos_parameters.allocation_and_retention_priority
      .pre_emption_vulnerability = bearer_context->preemption_vulnerability;
    s1ap_e_rab_setup_req->e_rab_to_be_setup_list.item[0]
      .e_rab_level_qos_parameters.allocation_and_retention_priority
      .priority_level = bearer_context->priority_level;
    s1ap_e_rab_setup_req->e_rab_to_be_setup_list.item[0]
      .e_rab_level_qos_parameters.gbr_qos_information
      .e_rab_maximum_bit_rate_downlink = mbr_dl;
    s1ap_e_rab_setup_req->e_rab_to_be_setup_list.item[0]
      .e_rab_level_qos_parameters.gbr_qos_information
      .e_rab_maximum_bit_rate_uplink = mbr_ul;
    s1ap_e_rab_setup_req->e_rab_to_be_setup_list.item[0]
      .e_rab_level_qos_parameters.gbr_qos_information
      .e_rab_guaranteed_bit_rate_downlink = gbr_dl;
    s1ap_e_rab_setup_req->e_rab_to_be_setup_list.item[0]
      .e_rab_level_qos_parameters.gbr_qos_information
      .e_rab_guaranteed_bit_rate_uplink = gbr_ul;
    s1ap_e_rab_setup_req->e_rab_to_be_setup_list.item[0]
      .e_rab_level_qos_parameters.qci = bearer_context->qci;

    s1ap_e_rab_setup_req->e_rab_to_be_setup_list.item[0].gtp_teid =
      bearer_context->s_gw_fteid_s1u.teid;
    s1ap_e_rab_setup_req->e_rab_to_be_setup_list.item[0]
      .transport_layer_address =
      fteid_ip_address_to_bstring(&bearer_context->s_gw_fteid_s1u);

    s1ap_e_rab_setup_req->e_rab_to_be_setup_list.item[0].nas_pdu = nas_msg;

    message_p->ittiMsgHeader.imsi = ue_context_p->emm_context._imsi64;
    itti_send_msg_to_task(TASK_S1AP, INSTANCE_DEFAULT, message_p);
  } else {
    OAILOG_DEBUG(
      LOG_MME_APP,
      "No bearer context found for ue-id " MME_UE_S1AP_ID_FMT " ebi %u\n",
      ue_id,
      ebi);
    bdestroy_wrapper(&nas_msg);
  }
  OAILOG_FUNC_OUT(LOG_MME_APP);
}

//------------------------------------------------------------------------------
void mme_app_handle_delete_session_rsp(
  mme_app_desc_t* mme_app_desc_p,
  const itti_s11_delete_session_response_t* const delete_sess_resp_pP)
//------------------------------------------------------------------------------
{
  struct ue_mm_context_s* ue_context_p = NULL;
  emm_cn_pdn_disconnect_rsp_t pdn_disconnect_rsp = {0};

  OAILOG_FUNC_IN(LOG_MME_APP);
  if (!delete_sess_resp_pP) {
    OAILOG_WARNING(
      LOG_MME_APP,
      "message, itti_s11_delete_session_response_t received"
      " from SGW is NULL \n");
    OAILOG_FUNC_OUT(LOG_MME_APP);
  }
  OAILOG_DEBUG(
    LOG_MME_APP,
    "Received S11_DELETE_SESSION_RESPONSE from S+P-GW with teid " TEID_FMT
    "\n ",
    delete_sess_resp_pP->teid);
  ue_context_p = mme_ue_context_exists_s11_teid(
    &mme_app_desc_p->mme_ue_contexts, delete_sess_resp_pP->teid);

  if (!ue_context_p) {
    OAILOG_WARNING(
      LOG_MME_APP,
      "We didn't find this teid in list of UE: %08x\n",
      delete_sess_resp_pP->teid);
    OAILOG_FUNC_OUT(LOG_MME_APP);
  }
  if (delete_sess_resp_pP->cause.cause_value != REQUEST_ACCEPTED) {
    OAILOG_WARNING(
      LOG_MME_APP,
      "***WARNING****S11 Delete Session Rsp: NACK received from SPGW : %08x\n",
      delete_sess_resp_pP->teid);
    increment_counter("mme_spgw_delete_session_rsp", 1, 1, "result", "failure");
  }
  increment_counter("mme_spgw_delete_session_rsp", 1, 1, "result", "success");
  /*
   * Updating statistics
   */
  update_mme_app_stats_s1u_bearer_sub();
  update_mme_app_stats_default_bearer_sub();

  /* If VoLTE is enabled and UE has sent PDN Disconnect
   * send pdn disconnect response to NAS.
   * NAS will trigger deactivate Bearer Context Req to UE
   */
  if (ue_context_p->emm_context.esm_ctx.is_pdn_disconnect) {
    pdn_disconnect_rsp.ue_id = ue_context_p->mme_ue_s1ap_id;
    pdn_disconnect_rsp.lbi = delete_sess_resp_pP->lbi;
    if ((nas_proc_pdn_disconnect_rsp(&pdn_disconnect_rsp)) != RETURNok) {
      OAILOG_ERROR(
        LOG_MME_APP,
        "Failed to handle PDN Disconnect Response at NAS module for "
        "ue_id " MME_UE_S1AP_ID_FMT " and lbi:%u \n",
        ue_context_p->mme_ue_s1ap_id,
        pdn_disconnect_rsp.lbi);
    }
    OAILOG_FUNC_OUT(LOG_MME_APP);
  }

  hashtable_uint64_ts_remove(
    mme_app_desc_p->mme_ue_contexts.tun11_ue_context_htbl,
    (const hash_key_t) ue_context_p->mme_teid_s11);
  ue_context_p->mme_teid_s11 = 0;

  /*
   * If UE is already in idle state, skip asking eNB to release UE context and just clean up locally.
   * This can happen during implicit detach and UE initiated detach when UE sends detach req (type = switch off)
   */
  if (ECM_IDLE == ue_context_p->ecm_state) {
    ue_context_p->ue_context_rel_cause = S1AP_IMPLICIT_CONTEXT_RELEASE;
    // Notify S1AP to release S1AP UE context locally.
    mme_app_itti_ue_context_release(
      ue_context_p, ue_context_p->ue_context_rel_cause);
    // Free MME UE Context
    mme_notify_ue_context_released(&mme_app_desc_p->mme_ue_contexts,
        ue_context_p);
    // Send PUR,before removal of ue contexts
    if (
      (ue_context_p->send_ue_purge_request == true) &&
      (ue_context_p->hss_initiated_detach == false)) {
      mme_app_send_s6a_purge_ue_req(mme_app_desc_p, ue_context_p);
    }
    OAILOG_DEBUG(
      LOG_MME_APP,
      "Deleting UE context associated in MME for "
      "mme_ue_s1ap_id " MME_UE_S1AP_ID_FMT "\n ",
      ue_context_p->mme_ue_s1ap_id);
    mme_remove_ue_context(&mme_app_desc_p->mme_ue_contexts, ue_context_p);
    OAILOG_FUNC_OUT(LOG_MME_APP);
  } else {
    if (ue_context_p->ue_context_rel_cause == S1AP_INVALID_CAUSE) {
      ue_context_p->ue_context_rel_cause = S1AP_NAS_DETACH;
    }
    /* In case of Ue initiated explicit IMSI Detach or Combined EPS/IMSI detach
       Do not send UE Context Release Command to eNB before receiving SGs IMSI Detach Ack from MSC/VLR */
    if (ue_context_p->sgs_context != NULL) {
      if (
        ((ue_context_p->sgs_detach_type !=
          SGS_EXPLICIT_UE_INITIATED_IMSI_DETACH_FROM_NONEPS) ||
         (ue_context_p->sgs_detach_type !=
          SGS_COMBINED_UE_INITIATED_IMSI_DETACH_FROM_EPS_N_NONEPS)) &&
        (ue_context_p->sgs_context->ts9_timer.id ==
         MME_APP_TIMER_INACTIVE_ID)) {
        /* Notify S1AP to send UE Context Release Command to eNB or free
         * s1 context locally.
         */
        mme_app_itti_ue_context_release(
          ue_context_p, ue_context_p->ue_context_rel_cause);
        ue_context_p->ue_context_rel_cause = S1AP_INVALID_CAUSE;
      }
    } else {
      // Notify S1AP to send UE Context Release Command to eNB or free s1 context locally.
      mme_app_itti_ue_context_release(
        ue_context_p, ue_context_p->ue_context_rel_cause);
      ue_context_p->ue_context_rel_cause = S1AP_INVALID_CAUSE;
    }
  }
  OAILOG_FUNC_OUT(LOG_MME_APP);
}

//------------------------------------------------------------------------------
int mme_app_handle_create_sess_resp(
  mme_app_desc_t* mme_app_desc_p,
  itti_s11_create_session_response_t* const create_sess_resp_pP)
{
  OAILOG_FUNC_IN(LOG_MME_APP);
  struct ue_mm_context_s* ue_context_p = NULL;
  bearer_context_t* current_bearer_p = NULL;
  ebi_t bearer_id = 0;
  int rc = RETURNok;

  OAILOG_INFO(
    LOG_MME_APP, "Received S11_CREATE_SESSION_RESPONSE from S+P-GW\n");

  if (create_sess_resp_pP == NULL) {
    OAILOG_ERROR(
      LOG_MME_APP, "Invalid Create Session Response object received\n");
    OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNerror);
  }

  ue_context_p = mme_ue_context_exists_s11_teid(
    &mme_app_desc_p->mme_ue_contexts, create_sess_resp_pP->teid);

  if (ue_context_p == NULL) {
    OAILOG_ERROR(
      LOG_MME_APP,
      "We didn't find this teid in list of UE: %08x\n",
      create_sess_resp_pP->teid);
    OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNerror);
  }

  OAILOG_DEBUG(
    LOG_MME_APP,
    "MME S11 teid = %u, cause = %d, ue_id = %u\n",
    create_sess_resp_pP->teid,
    create_sess_resp_pP->cause.cause_value,
    ue_context_p->mme_ue_s1ap_id);

  proc_tid_t transaction_identifier = 0;
  pdn_cid_t pdn_cx_id = 0;

  /* Whether SGW has created the session (IP address allocation, local GTP-U
   * end point creation etc.) successfully or not, it is indicated by cause
   * value in create session response message.
   * If cause value is not equal to "REQUEST_ACCEPTED" then this implies that
   * SGW could not allocate the resources for the requested session. In this
   * case, MME-APP sends PDN Connectivity fail message to NAS-ESM with the
   * "cause" received in S11 Create Session Response message.
   * NAS-ESM maps this "S11 cause" to "ESM cause" and sends it in PDN
   * Connectivity Reject message to the UE.
   */

  emm_cn_ula_or_csrsp_fail_t create_session_response_fail = {0};
  if (create_sess_resp_pP->cause.cause_value != REQUEST_ACCEPTED) {
    // Send PDN CONNECTIVITY FAIL message  to NAS layer
    OAILOG_DEBUG(
      LOG_MME_APP,
      "Create Session Response Cause value = (%d) for ue_id =(%u)\n",
      create_sess_resp_pP->cause.cause_value,
      ue_context_p->mme_ue_s1ap_id);
    create_session_response_fail.cause =
      (pdn_conn_rsp_cause_t)(create_sess_resp_pP->cause.cause_value);
    goto error_handling_csr_failure;
  }
  increment_counter("mme_spgw_create_session_rsp", 1, 1, "result", "success");
  //---------------------------------------------------------
  // Process itti_sgw_create_session_response_t.bearer_context_created
  //---------------------------------------------------------
  int num_successful_bearers = 0;
  for (int i = 0;
       i < create_sess_resp_pP->bearer_contexts_created.num_bearer_context;
       i++) {
    bearer_id = create_sess_resp_pP->bearer_contexts_created.bearer_contexts[i]
                  .eps_bearer_id /* - 5 */;
    /*
     * Depending on s11 result we have to send reject or accept for bearers
     */
    if (ue_context_p->emm_context.esm_ctx.n_active_ebrs > BEARERS_PER_UE) {
      OAILOG_ERROR(
        LOG_MME_APP,
        "The total number of active EPS bearers has exceeded %d\n",
        ue_context_p->emm_context.esm_ctx.n_active_ebrs);
      OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNerror);
    }
    if (
      create_sess_resp_pP->bearer_contexts_created.bearer_contexts[i]
        .cause.cause_value != REQUEST_ACCEPTED) {
      OAILOG_ERROR(
        LOG_MME_APP,
        "Cases where bearer cause != REQUEST_ACCEPTED are not handled\n");
      continue;
    }
    if (
      create_sess_resp_pP->bearer_contexts_created.bearer_contexts[i]
        .s1u_sgw_fteid.interface_type != S1_U_SGW_GTP_U) {
      OAILOG_ERROR(
        LOG_MME_APP,
        "Invalid S1U SGW GTP F-TEID interface type: %d (Expected: %d)\n",
        create_sess_resp_pP->bearer_contexts_created.bearer_contexts[i]
          .s1u_sgw_fteid.interface_type,
        S1_U_SGW_GTP_U);
      continue;
    }

    current_bearer_p = mme_app_get_bearer_context(ue_context_p, bearer_id);
    if (current_bearer_p == NULL) {
      OAILOG_ERROR(
        LOG_MME_APP,
        "Failed to get bearer context for bearer Id (%d)\n",
        bearer_id);
      continue;
    }

    update_mme_app_stats_default_bearer_add();

    current_bearer_p->bearer_state |= BEARER_STATE_SGW_CREATED;
    if (!i) {
      pdn_cx_id = current_bearer_p->pdn_cx_id;
      /*
       * Store the S-GW teid
       */
      if ((pdn_cx_id < 0) || (pdn_cx_id >= MAX_APN_PER_UE)) {
        OAILOG_ERROR(LOG_MME_APP, "Bad pdn id (%d) for bearer\n", pdn_cx_id);
        continue;
      }
      ue_context_p->pdn_contexts[pdn_cx_id]->s_gw_teid_s11_s4 =
        create_sess_resp_pP->s11_sgw_fteid.teid;
      transaction_identifier = current_bearer_p->transaction_identifier;
    }

    current_bearer_p->s_gw_fteid_s1u =
      create_sess_resp_pP->bearer_contexts_created.bearer_contexts[i]
        .s1u_sgw_fteid;
    current_bearer_p->p_gw_fteid_s5_s8_up =
      create_sess_resp_pP->bearer_contexts_created.bearer_contexts[i]
        .s5_s8_u_pgw_fteid;
    OAILOG_DEBUG(
      LOG_MME_APP,
      "S1U S-GW teid   = (%u)\n"
      "S5/S8U PGW teid = (%u)\n",
      current_bearer_p->s_gw_fteid_s1u.teid,
      current_bearer_p->p_gw_fteid_s5_s8_up.teid);

    // if modified by pgw
    if (create_sess_resp_pP->bearer_contexts_created.bearer_contexts[i]
          .bearer_level_qos) {
      current_bearer_p->qci =
        create_sess_resp_pP->bearer_contexts_created.bearer_contexts[i]
          .bearer_level_qos->qci;
      current_bearer_p->priority_level =
        create_sess_resp_pP->bearer_contexts_created.bearer_contexts[i]
          .bearer_level_qos->pl;
      current_bearer_p->preemption_vulnerability =
        create_sess_resp_pP->bearer_contexts_created.bearer_contexts[i]
          .bearer_level_qos->pvi;
      current_bearer_p->preemption_capability =
        create_sess_resp_pP->bearer_contexts_created.bearer_contexts[i]
          .bearer_level_qos->pci;

      //TODO should be set in NAS_PDN_CONNECTIVITY_RSP message
      current_bearer_p->esm_ebr_context.gbr_dl =
        create_sess_resp_pP->bearer_contexts_created.bearer_contexts[i]
          .bearer_level_qos->gbr.br_dl;
      current_bearer_p->esm_ebr_context.gbr_ul =
        create_sess_resp_pP->bearer_contexts_created.bearer_contexts[i]
          .bearer_level_qos->gbr.br_ul;
      current_bearer_p->esm_ebr_context.mbr_dl =
        create_sess_resp_pP->bearer_contexts_created.bearer_contexts[i]
          .bearer_level_qos->mbr.br_dl;
      current_bearer_p->esm_ebr_context.mbr_ul =
        create_sess_resp_pP->bearer_contexts_created.bearer_contexts[i]
          .bearer_level_qos->mbr.br_ul;
      OAILOG_DEBUG(
        LOG_MME_APP,
        "Set qci %u in bearer %u\n",
        current_bearer_p->qci,
        bearer_id);
    } else {
      OAILOG_DEBUG(
        LOG_MME_APP,
        "Set qci %u in bearer %u (qos not modified by P-GW)\n",
        current_bearer_p->qci,
        bearer_id);
    }
    ++num_successful_bearers;
  }
  if (num_successful_bearers == 0) {
    // Send PDN CONNECTIVITY FAIL message to NAS layer, if none of the bearer
    // could be allocated
    create_session_response_fail.cause = CAUSE_NO_RESOURCES_AVAILABLE;
    goto error_handling_csr_failure;
  }
  /* Send Create Session Response to NAS module */
  emm_cn_cs_response_success_t nas_pdn_cs_respose_success = {0};
  nas_pdn_cs_respose_success.pdn_cid = pdn_cx_id;
  nas_pdn_cs_respose_success.pti = transaction_identifier; // NAS internal ref
  nas_pdn_cs_respose_success.pdn_addr =
    paa_to_bstring(&create_sess_resp_pP->paa);
  nas_pdn_cs_respose_success.pdn_type = create_sess_resp_pP->paa.pdn_type;

  // ASSUME NO HO now, so assume 1 bearer only and is default bearer

  nas_pdn_cs_respose_success.ue_id = ue_context_p->mme_ue_s1ap_id;
  nas_pdn_cs_respose_success.ebi = bearer_id;
  nas_pdn_cs_respose_success.qci = current_bearer_p->qci;
  nas_pdn_cs_respose_success.prio_level = current_bearer_p->priority_level;
  nas_pdn_cs_respose_success.pre_emp_vulnerability =
    current_bearer_p->preemption_vulnerability;
  nas_pdn_cs_respose_success.pre_emp_capability =
    current_bearer_p->preemption_capability;
  nas_pdn_cs_respose_success.sgw_s1u_fteid = current_bearer_p->s_gw_fteid_s1u;
  // optional IE
  nas_pdn_cs_respose_success.ambr.br_ul =
    ue_context_p->subscribed_ue_ambr.br_ul;
  nas_pdn_cs_respose_success.ambr.br_dl =
    ue_context_p->subscribed_ue_ambr.br_dl;

  // This IE is not applicable for TAU/RAU/Handover.
  // If PGW decides to return PCO to the UE, PGW shall send PCO to
  // SGW. If SGW receives the PCO IE, SGW shall forward it to MME/SGSN.
  if (create_sess_resp_pP->pco.num_protocol_or_container_id) {
    copy_protocol_configuration_options(
      &nas_pdn_cs_respose_success.pco, &create_sess_resp_pP->pco);
    clear_protocol_configuration_options(&create_sess_resp_pP->pco);
  }

  nas_proc_cs_respose_success(&nas_pdn_cs_respose_success);
  OAILOG_FUNC_RETURN(LOG_MME_APP, rc);

error_handling_csr_failure:
  increment_counter("mme_spgw_create_session_rsp", 1, 1, "result", "failure");
  bearer_id = create_sess_resp_pP->bearer_contexts_created.bearer_contexts[0]
                .eps_bearer_id /* - 5 */;
  current_bearer_p = mme_app_get_bearer_context(ue_context_p, bearer_id);
  if (current_bearer_p) {
    transaction_identifier = current_bearer_p->transaction_identifier;
  }
  create_session_response_fail.pti = transaction_identifier;
  create_session_response_fail.ue_id = ue_context_p->mme_ue_s1ap_id;
  OAILOG_ERROR(
    LOG_MME_APP,
    "Handling Create Session Response failure for ue_id = (%u), "
    "bearer id = (%d), pti = (%d)\n",
    ue_context_p->mme_ue_s1ap_id,
    bearer_id,
    transaction_identifier);
  rc = nas_proc_ula_or_csrsp_fail(&create_session_response_fail);
  OAILOG_FUNC_RETURN(LOG_MME_APP, rc);
}

//------------------------------------------------------------------------------
void mme_app_handle_initial_context_setup_rsp(
  itti_mme_app_initial_context_setup_rsp_t *const initial_ctxt_setup_rsp_pP)
{
  OAILOG_FUNC_IN(LOG_MME_APP);
  struct ue_mm_context_s *ue_context_p = NULL;
  MessageDef *message_p = NULL;

  OAILOG_INFO(
    LOG_MME_APP, "Received MME_APP_INITIAL_CONTEXT_SETUP_RSP from S1AP for ue_id = (%u)\n",
    initial_ctxt_setup_rsp_pP->ue_id);
  ue_context_p =
    mme_ue_context_exists_mme_ue_s1ap_id(initial_ctxt_setup_rsp_pP->ue_id);

  if (ue_context_p == NULL) {
    OAILOG_DEBUG(
      LOG_MME_APP,
      " We didn't find this mme_ue_s1ap_id in list of UE: " MME_UE_S1AP_ID_FMT
      "\n UE Context NULL...\n",
      initial_ctxt_setup_rsp_pP->ue_id);
    OAILOG_FUNC_OUT(LOG_MME_APP);
  }

  // Stop Initial context setup process guard timer,if running
  if (
    ue_context_p->initial_context_setup_rsp_timer.id !=
    MME_APP_TIMER_INACTIVE_ID) {
    if (timer_remove(ue_context_p->initial_context_setup_rsp_timer.id, NULL)) {
      OAILOG_ERROR(
        LOG_MME_APP,
        "Failed to stop Initial Context Setup Rsp timer for UE id  %d \n",
        ue_context_p->mme_ue_s1ap_id);
    }
    ue_context_p->initial_context_setup_rsp_timer.id =
      MME_APP_TIMER_INACTIVE_ID;
  }

  message_p = itti_alloc_new_message(TASK_MME_APP, S11_MODIFY_BEARER_REQUEST);
  if (message_p == NULL) {
    OAILOG_ERROR(
      LOG_MME_APP,
      "Failed to allocate new ITTI message for S11 Modify Bearer Request "
      "for MME UE S1AP Id: " MME_UE_S1AP_ID_FMT "\n",
      ue_context_p->mme_ue_s1ap_id);
    OAILOG_FUNC_OUT(LOG_MME_APP);
  }
  itti_s11_modify_bearer_request_t *s11_modify_bearer_request =
    &message_p->ittiMsg.s11_modify_bearer_request;
  s11_modify_bearer_request->local_teid = ue_context_p->mme_teid_s11;
  /*
   * Delay Value in integer multiples of 50 millisecs, or zero
   */
  s11_modify_bearer_request->delay_dl_packet_notif_req = 0; // TODO

  for (int item = 0; item < initial_ctxt_setup_rsp_pP->no_of_e_rabs; item++) {
    s11_modify_bearer_request->bearer_contexts_to_be_modified
      .bearer_contexts[item]
      .eps_bearer_id = initial_ctxt_setup_rsp_pP->e_rab_id[item];
    s11_modify_bearer_request->bearer_contexts_to_be_modified
      .bearer_contexts[item]
      .s1_eNB_fteid.teid = initial_ctxt_setup_rsp_pP->gtp_teid[item];
    s11_modify_bearer_request->bearer_contexts_to_be_modified
      .bearer_contexts[item]
      .s1_eNB_fteid.interface_type = S1_U_ENODEB_GTP_U;


    if (!item) {
      ebi_t ebi = initial_ctxt_setup_rsp_pP->e_rab_id[item];
      pdn_cid_t cid =
        ue_context_p->bearer_contexts[EBI_TO_INDEX(ebi)]->pdn_cx_id;
      pdn_context_t *pdn_context = ue_context_p->pdn_contexts[cid];
     pdn_context->s_gw_address_s11_s4.address.ipv4_address.s_addr = 
     mme_config.e_dns_emulation.sgw_ip_addr[0].s_addr;

     s11_modify_bearer_request->edns_peer_ip.addr_v4_.sin_addr.s_addr =
     pdn_context->s_gw_address_s11_s4.address.ipv4_address.s_addr;
     s11_modify_bearer_request->edns_peer_ip.addr_v4_.sin_family = AF_INET;    
            
      s11_modify_bearer_request->teid = pdn_context->s_gw_teid_s11_s4;
    }
    if (
      4 == blength(initial_ctxt_setup_rsp_pP->transport_layer_address[item])) {
      s11_modify_bearer_request->bearer_contexts_to_be_modified
        .bearer_contexts[item]
        .s1_eNB_fteid.ipv4 = 1;
      memcpy(
        &s11_modify_bearer_request->bearer_contexts_to_be_modified
           .bearer_contexts[item]
           .s1_eNB_fteid.ipv4_address,
        initial_ctxt_setup_rsp_pP->transport_layer_address[item]->data,
        blength(initial_ctxt_setup_rsp_pP->transport_layer_address[item]));
    } else if (
      16 == blength(initial_ctxt_setup_rsp_pP->transport_layer_address[item])) {
      s11_modify_bearer_request->bearer_contexts_to_be_modified
        .bearer_contexts[item]
        .s1_eNB_fteid.ipv6 = 1;
      memcpy(
        &s11_modify_bearer_request->bearer_contexts_to_be_modified
           .bearer_contexts[item]
           .s1_eNB_fteid.ipv6_address,
        initial_ctxt_setup_rsp_pP->transport_layer_address[item]->data,
        blength(initial_ctxt_setup_rsp_pP->transport_layer_address[item]));
    } else {
      OAILOG_ERROR(
        LOG_MME_APP,
        "Invalid IP address of %d bytes found for MME UE S1AP "
        "Id: " MME_UE_S1AP_ID_FMT " (4 or 16 bytes was expected)\n",
        blength(initial_ctxt_setup_rsp_pP->transport_layer_address[item]),
        ue_context_p->mme_ue_s1ap_id);
      OAILOG_FUNC_OUT(LOG_MME_APP);
    }
    bdestroy_wrapper(&initial_ctxt_setup_rsp_pP->transport_layer_address[item]);
  }
  s11_modify_bearer_request->bearer_contexts_to_be_modified.num_bearer_context =
    initial_ctxt_setup_rsp_pP->no_of_e_rabs;

  s11_modify_bearer_request->bearer_contexts_to_be_removed.num_bearer_context =
    0;

  s11_modify_bearer_request->mme_fq_csid.node_id_type =
    GLOBAL_UNICAST_IPv4;                           // TODO
  s11_modify_bearer_request->mme_fq_csid.csid = 0; // TODO ...
  memset(
    &s11_modify_bearer_request->indication_flags,
    0,
    sizeof(s11_modify_bearer_request->indication_flags)); // TODO
  s11_modify_bearer_request->rat_type = RAT_EUTRAN;
  /*
   * S11 stack specific parameter. Not used in standalone epc mode
   */
  s11_modify_bearer_request->trxn = NULL;

  message_p->ittiMsgHeader.imsi = ue_context_p->emm_context._imsi64;

  OAILOG_INFO(
    LOG_MME_APP,
    "Sending S11 MODIFY BEARER REQ to SPGW for ue_id = (%d), teid = (%u)\n",
    initial_ctxt_setup_rsp_pP->ue_id,
    s11_modify_bearer_request->teid);
  itti_send_msg_to_task(TASK_SPGW, INSTANCE_DEFAULT, message_p);
  /*
   * During Service request procedure,after initial context setup response
   * Send ULR, when UE moved from Idle to Connected and
   * flag location_info_confirmed_in_hss set to true during hss reset.
   */
  if (ue_context_p->location_info_confirmed_in_hss == true) {
    mme_app_send_s6a_update_location_req(ue_context_p);
  }
  if (ue_context_p->sgs_context) {
    ue_context_p->sgs_context->csfb_service_type = CSFB_SERVICE_NONE;
    //Reset mt_call_in_progress flag
    if (ue_context_p->sgs_context->mt_call_in_progress) {
      ue_context_p->sgs_context->mt_call_in_progress = false;
    }
  }
  OAILOG_FUNC_OUT(LOG_MME_APP);
}

//------------------------------------------------------------------------------
void mme_app_handle_release_access_bearers_resp(mme_app_desc_t *mme_app_desc_p,
  const itti_s11_release_access_bearers_response_t
    *const rel_access_bearers_rsp_pP)
{
  OAILOG_FUNC_IN(LOG_MME_APP);
  struct ue_mm_context_s *ue_context_p = NULL;

  ue_context_p = mme_ue_context_exists_s11_teid(
    &mme_app_desc_p->mme_ue_contexts, rel_access_bearers_rsp_pP->teid);

  if (ue_context_p == NULL) {
    OAILOG_DEBUG(
      LOG_MME_APP,
      "We didn't find this teid in list of UE: %" PRIX32 "\n",
      rel_access_bearers_rsp_pP->teid);
    OAILOG_FUNC_OUT(LOG_MME_APP);
  }
  /*
   * Updating statistics
   */
  update_mme_app_stats_s1u_bearer_sub();

  // Send UE Context Release Command
  mme_app_itti_ue_context_release(
    ue_context_p, ue_context_p->ue_context_rel_cause);
  if (
    ue_context_p->ue_context_rel_cause == S1AP_SCTP_SHUTDOWN_OR_RESET ||
    ue_context_p->ue_context_rel_cause ==
      S1AP_INITIAL_CONTEXT_SETUP_TMR_EXPRD) {
    // Just cleanup the MME APP state associated with s1.
    mme_ue_context_update_ue_sig_connection_state(
      &mme_app_desc_p->mme_ue_contexts, ue_context_p, ECM_IDLE);
    ue_context_p->ue_context_rel_cause = S1AP_INVALID_CAUSE;
  }
  OAILOG_FUNC_OUT(LOG_MME_APP);
}

//------------------------------------------------------------------------------
void mme_app_handle_s11_create_bearer_req(
  mme_app_desc_t* mme_app_desc_p,
  const itti_s11_create_bearer_request_t* const create_bearer_request_pP)
{
  OAILOG_FUNC_IN(LOG_MME_APP);
  struct ue_mm_context_s* ue_context_p = NULL;
  emm_cn_activate_dedicated_bearer_req_t activate_ded_bearer_req = {0};

  ue_context_p = mme_ue_context_exists_s11_teid(
    &mme_app_desc_p->mme_ue_contexts, create_bearer_request_pP->teid);

  if (ue_context_p == NULL) {
    OAILOG_DEBUG(
      LOG_MME_APP,
      "We didn't find this teid in list of UE: %" PRIX32 "\n",
      create_bearer_request_pP->teid);
    OAILOG_FUNC_OUT(LOG_MME_APP);
  }

  // check if default bearer already created
  ebi_t linked_eps_bearer_id = create_bearer_request_pP->linked_eps_bearer_id;
  bearer_context_t* linked_bc =
    mme_app_get_bearer_context(ue_context_p, linked_eps_bearer_id);
  if (!linked_bc) {
    // May create default EPS bearer ?
    OAILOG_DEBUG(
      LOG_MME_APP,
      "We didn't find the default bearer context for linked bearer id %" PRIu8
      " of ue_id: " MME_UE_S1AP_ID_FMT "\n",
      linked_eps_bearer_id,
      ue_context_p->mme_ue_s1ap_id);
    OAILOG_FUNC_OUT(LOG_MME_APP);
  }

  pdn_cid_t cid = linked_bc->pdn_cx_id;

  mme_app_s11_proc_create_bearer_t* s11_proc_create_bearer =
    mme_app_create_s11_procedure_create_bearer(ue_context_p);
  s11_proc_create_bearer->proc.s11_trxn =
    (uintptr_t) create_bearer_request_pP->trxn;

  for (int i = 0;
       i < create_bearer_request_pP->bearer_contexts.num_bearer_context;
       i++) {
    //!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
    // TODO THINK OF BEARER AGGREGATING SEVERAL SDFs, 1 bearer <-> (QCI, ARP)
    // TODO DELEGATE TO NAS THE CREATION OF THE BEARER
    //!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
    const bearer_context_within_create_bearer_request_t* msg_bc =
      &create_bearer_request_pP->bearer_contexts.bearer_contexts[i];
    bearer_context_t* dedicated_bc = mme_app_create_bearer_context(
      ue_context_p, cid, msg_bc->eps_bearer_id, false);

    s11_proc_create_bearer->num_bearers++;
    s11_proc_create_bearer->bearer_status[EBI_TO_INDEX(dedicated_bc->ebi)] =
      S11_PROC_BEARER_PENDING;

    dedicated_bc->bearer_state |= BEARER_STATE_SGW_CREATED;
    dedicated_bc->bearer_state |= BEARER_STATE_MME_CREATED;

    dedicated_bc->s_gw_fteid_s1u = msg_bc->s1u_sgw_fteid;
    dedicated_bc->p_gw_fteid_s5_s8_up = msg_bc->s5_s8_u_pgw_fteid;

    dedicated_bc->qci = msg_bc->bearer_level_qos.qci;
    dedicated_bc->priority_level = msg_bc->bearer_level_qos.pl;
    dedicated_bc->preemption_vulnerability = msg_bc->bearer_level_qos.pvi;
    dedicated_bc->preemption_capability = msg_bc->bearer_level_qos.pci;

    // forward request to NAS
    activate_ded_bearer_req.ue_id = ue_context_p->mme_ue_s1ap_id;
    activate_ded_bearer_req.cid = cid;
    activate_ded_bearer_req.ebi = dedicated_bc->ebi;
    activate_ded_bearer_req.linked_ebi =
      ue_context_p->pdn_contexts[cid]->default_ebi;
    activate_ded_bearer_req.bearer_qos = msg_bc->bearer_level_qos;
    if (msg_bc->tft.numberofpacketfilters) {
      activate_ded_bearer_req.tft = calloc(1, sizeof(traffic_flow_template_t));
      copy_traffic_flow_template(activate_ded_bearer_req.tft, &msg_bc->tft);
    }
    if (msg_bc->pco.num_protocol_or_container_id) {
      activate_ded_bearer_req.pco =
        calloc(1, sizeof(protocol_configuration_options_t));
      copy_protocol_configuration_options(
        activate_ded_bearer_req.pco, &msg_bc->pco);
    }
    if (
      (nas_proc_create_dedicated_bearer(&activate_ded_bearer_req)) !=
      RETURNok) {
      OAILOG_ERROR(
        LOG_MME_APP,
        "Failed to handle bearer activation at NAS module for "
        "ue_id " MME_UE_S1AP_ID_FMT "\n",
        ue_context_p->mme_ue_s1ap_id);
    }
    free_wrapper((void **) &activate_ded_bearer_req.tft);
    free_wrapper((void **) &activate_ded_bearer_req.pco);
  }
  OAILOG_FUNC_OUT(LOG_MME_APP);
}

void mme_app_handle_e_rab_setup_rsp(
  itti_s1ap_e_rab_setup_rsp_t* const e_rab_setup_rsp)
{
  OAILOG_FUNC_IN(LOG_MME_APP);
  struct ue_mm_context_s* ue_context_p = NULL;
  bool send_s11_response = false;
  OAILOG_INFO(
    LOG_MME_APP,
    "Received S1AP_E_RAB_SETUP_RSP from S1AP for ue_id:" MME_UE_S1AP_ID_FMT
    "\n",
    e_rab_setup_rsp->mme_ue_s1ap_id);

  ue_context_p =
    mme_ue_context_exists_mme_ue_s1ap_id(e_rab_setup_rsp->mme_ue_s1ap_id);

  if (ue_context_p == NULL) {
    OAILOG_ERROR(
      LOG_MME_APP,
      "We didn't find this mme_ue_s1ap_id in list of UE: " MME_UE_S1AP_ID_FMT
      "\n",
      e_rab_setup_rsp->mme_ue_s1ap_id);
    OAILOG_FUNC_OUT(LOG_MME_APP);
  }

  for (int i = 0; i < e_rab_setup_rsp->e_rab_setup_list.no_of_items; i++) {
    e_rab_id_t e_rab_id = e_rab_setup_rsp->e_rab_setup_list.item[i].e_rab_id;
    bearer_context_t* bc =
      mme_app_get_bearer_context(ue_context_p, (ebi_t) e_rab_id);
    if (bc->bearer_state & BEARER_STATE_SGW_CREATED) {
      bc->enb_fteid_s1u.teid =
        e_rab_setup_rsp->e_rab_setup_list.item[i].gtp_teid;
      // Do not process transport_layer_address now
      //bstring e_rab_setup_rsp->e_rab_setup_list.item[i].transport_layer_address;
      ip_address_t enb_ip_address = {0};
      bstring_to_ip_address(
        e_rab_setup_rsp->e_rab_setup_list.item[i].transport_layer_address,
        &enb_ip_address);

      bc->enb_fteid_s1u.interface_type = S1_U_ENODEB_GTP_U;
      // TODO better than that later
      switch (enb_ip_address.pdn_type) {
        case IPv4:
          bc->enb_fteid_s1u.ipv4 = 1;
          bc->enb_fteid_s1u.ipv4_address = enb_ip_address.address.ipv4_address;
          break;
        case IPv6:
          bc->enb_fteid_s1u.ipv6 = 1;
          memcpy(
            &bc->enb_fteid_s1u.ipv6_address,
            &enb_ip_address.address.ipv6_address,
            sizeof(enb_ip_address.address.ipv6_address));
          break;
        default:
          OAILOG_ERROR(
            LOG_MME_APP,
            "Invalid eNB IP address PDN type received for MME UE S1AP "
            "Id: " MME_UE_S1AP_ID_FMT "\n",
            e_rab_setup_rsp->mme_ue_s1ap_id);
          OAILOG_FUNC_OUT(LOG_MME_APP);
      }
      bdestroy_wrapper(
        &e_rab_setup_rsp->e_rab_setup_list.item[i].transport_layer_address);

      if (!(bc->bearer_state & BEARER_STATE_MME_CREATED)) {
        OAILOG_ERROR(
          LOG_MME_APP,
          "Bearer not created in MME for bearer Id: %u\n",
          e_rab_id);
        OAILOG_FUNC_OUT(LOG_MME_APP);
      }
      bc->bearer_state |= BEARER_STATE_ENB_CREATED;
    }
  }
  for (int i = 0; i < e_rab_setup_rsp->e_rab_failed_to_setup_list.no_of_items;
       i++) {
    e_rab_id_t e_rab_id =
      e_rab_setup_rsp->e_rab_failed_to_setup_list.item[i].e_rab_id;
    bearer_context_t* bc =
      mme_app_get_bearer_context(ue_context_p, (ebi_t) e_rab_id);
    if (bc->bearer_state & BEARER_STATE_SGW_CREATED) {
      if (!(bc->bearer_state & BEARER_STATE_MME_CREATED)) {
        OAILOG_ERROR(
          LOG_MME_APP,
          "Bearer creation failed in eNB, but successfully created in "
          "MME/SGW for bearer Id: %u\n",
          e_rab_id);
        OAILOG_FUNC_OUT(LOG_MME_APP);
      }
      bc->bearer_state &= (~BEARER_STATE_ENB_CREATED);
      bc->bearer_state &= (~BEARER_STATE_MME_CREATED);
    }
  }

  // check if UE already responded with NAS (may depend on eNB implementation?) -> send response to SGW
  if (send_s11_response) {
    MessageDef *message_p =
      itti_alloc_new_message(TASK_MME_APP, S11_CREATE_BEARER_RESPONSE);
    if (message_p == NULL) {
      OAILOG_ERROR(
        LOG_MME_APP,
        "Failed to allocate new ITTI message for S11 Create Bearer Response "
        "for MME UE S1AP Id: " MME_UE_S1AP_ID_FMT "\n",
        e_rab_setup_rsp->mme_ue_s1ap_id);
      OAILOG_FUNC_OUT(LOG_MME_APP);
    }
    itti_s11_create_bearer_response_t *s11_create_bearer_response =
      &message_p->ittiMsg.s11_create_bearer_response;
    s11_create_bearer_response->local_teid = ue_context_p->mme_teid_s11;
    s11_create_bearer_response->trxn = NULL;
    s11_create_bearer_response->cause.cause_value = 0;
    int msg_bearer_index = 0;

    for (int i = 0; i < e_rab_setup_rsp->e_rab_setup_list.no_of_items; i++) {
      e_rab_id_t e_rab_id = e_rab_setup_rsp->e_rab_setup_list.item[i].e_rab_id;
      bearer_context_t *bc =
        mme_app_get_bearer_context(ue_context_p, (ebi_t) e_rab_id);
      if (bc->bearer_state & BEARER_STATE_ENB_CREATED) {
        s11_create_bearer_response->cause.cause_value = REQUEST_ACCEPTED;
        s11_create_bearer_response->bearer_contexts
          .bearer_contexts[msg_bearer_index]
          .eps_bearer_id = e_rab_id;
        s11_create_bearer_response->bearer_contexts
          .bearer_contexts[msg_bearer_index]
          .cause.cause_value = REQUEST_ACCEPTED;
        //  FTEID eNB
        s11_create_bearer_response->bearer_contexts
          .bearer_contexts[msg_bearer_index]
          .s1u_enb_fteid = bc->enb_fteid_s1u;

        // FTEID SGW S1U
        s11_create_bearer_response->bearer_contexts
          .bearer_contexts[msg_bearer_index]
          .s1u_sgw_fteid =
          bc->s_gw_fteid_s1u; ///< This IE shall be sent on the S11 interface. It shall be used
        s11_create_bearer_response->bearer_contexts.num_bearer_context++;
      }
    }

    for (int i = 0; i < e_rab_setup_rsp->e_rab_setup_list.no_of_items; i++) {
      e_rab_id_t e_rab_id = e_rab_setup_rsp->e_rab_setup_list.item[i].e_rab_id;
      bearer_context_t *bc =
        mme_app_get_bearer_context(ue_context_p, (ebi_t) e_rab_id);
      if (bc->bearer_state & BEARER_STATE_MME_CREATED) {
        if (REQUEST_ACCEPTED == s11_create_bearer_response->cause.cause_value) {
          s11_create_bearer_response->cause.cause_value =
            REQUEST_ACCEPTED_PARTIALLY;
        } else {
          s11_create_bearer_response->cause.cause_value = REQUEST_REJECTED;
        }
        s11_create_bearer_response->bearer_contexts
          .bearer_contexts[msg_bearer_index]
          .eps_bearer_id = e_rab_id;
        s11_create_bearer_response->bearer_contexts
          .bearer_contexts[msg_bearer_index]
          .cause.cause_value =
          REQUEST_REJECTED; // TODO translation of S1AP cause to SGW cause
        s11_create_bearer_response->bearer_contexts.num_bearer_context++;
        bc->bearer_state = BEARER_STATE_NULL;
      }
    }

    message_p->ittiMsgHeader.imsi = ue_context_p->emm_context._imsi64;
    itti_send_msg_to_task(TASK_S11, INSTANCE_DEFAULT, message_p);
  } else {
    // not send S11 response
    // TODO create a procedure with bearers to receive a response from NAS
  }
  OAILOG_FUNC_OUT(LOG_MME_APP);
}

//------------------------------------------------------------------------------
void mme_app_handle_mobile_reachability_timer_expiry(void* args)
{
  OAILOG_FUNC_IN(LOG_MME_APP);
  mme_ue_s1ap_id_t mme_ue_s1ap_id = *((mme_ue_s1ap_id_t *) (args));
  struct ue_mm_context_s* ue_context_p = mme_app_get_ue_context_for_timer(
    mme_ue_s1ap_id,
    "Mobile reachability timer");
  if (ue_context_p == NULL) {
    OAILOG_ERROR(
      LOG_MME_APP,
      "Invalid UE context received, MME UE S1AP Id: " MME_UE_S1AP_ID_FMT "\n",
      mme_ue_s1ap_id);
    OAILOG_FUNC_OUT(LOG_MME_APP);
  }
  ue_context_p->mobile_reachability_timer.id = MME_APP_TIMER_INACTIVE_ID;

  // Start Implicit Detach timer
  nas_itti_timer_arg_t timer_callback_fun = {0};
  timer_callback_fun.nas_timer_callback =
    mme_app_handle_implicit_detach_timer_expiry;
  timer_callback_fun.nas_timer_callback_arg =
    (void *) &(ue_context_p->mme_ue_s1ap_id);
  if (timer_setup(
    ue_context_p->implicit_detach_timer.sec,
    0,
    TASK_MME_APP,
    INSTANCE_DEFAULT,
    TIMER_ONE_SHOT,
    &timer_callback_fun,
    sizeof(timer_callback_fun),
    &(ue_context_p->implicit_detach_timer.id)) < 0) {
    OAILOG_ERROR(
      LOG_MME_APP,
      "Failed to start Implicit Detach timer for UE id: " MME_UE_S1AP_ID_FMT
      "\n",
      ue_context_p->mme_ue_s1ap_id);
    ue_context_p->implicit_detach_timer.id = MME_APP_TIMER_INACTIVE_ID;
  } else {
    OAILOG_DEBUG(
      LOG_MME_APP,
      "Started Implicit Detach timer for UE id: " MME_UE_S1AP_ID_FMT  "\n",
      ue_context_p->mme_ue_s1ap_id);
  }
  /* PPF is set to false due to "Inactivity of UE including non reception of periodic TAU
   * If CS paging is received for MT call, MME shall indicate to VLR that UE is unreachable
   */
  ue_context_p->ppf = false;
  OAILOG_FUNC_OUT(LOG_MME_APP);
}
//------------------------------------------------------------------------------
void mme_app_handle_implicit_detach_timer_expiry(void* args)
{
  OAILOG_FUNC_IN(LOG_MME_APP);
  mme_ue_s1ap_id_t mme_ue_s1ap_id = *((mme_ue_s1ap_id_t*) (args));
  struct ue_mm_context_s* ue_context_p = mme_app_get_ue_context_for_timer(
    mme_ue_s1ap_id,
    "Implicit detach timer");
  if (ue_context_p == NULL) {
    OAILOG_ERROR(
      LOG_MME_APP,
      "Invalid UE context received, MME UE S1AP Id: " MME_UE_S1AP_ID_FMT "\n",
      mme_ue_s1ap_id);
    OAILOG_FUNC_OUT(LOG_MME_APP);
  }
  ue_context_p->implicit_detach_timer.id = MME_APP_TIMER_INACTIVE_ID;
  // Initiate Implicit Detach for the UE
  nas_proc_implicit_detach_ue_ind(mme_ue_s1ap_id);
  OAILOG_FUNC_OUT(LOG_MME_APP);
}

//------------------------------------------------------------------------------
void mme_app_handle_initial_context_setup_rsp_timer_expiry(void* args)
{
  OAILOG_FUNC_IN(LOG_MME_APP);
  mme_ue_s1ap_id_t mme_ue_s1ap_id = *((mme_ue_s1ap_id_t*) (args));
  struct ue_mm_context_s* ue_context_p = mme_app_get_ue_context_for_timer(
    mme_ue_s1ap_id,
    "Initial context setup response timer");
  if (ue_context_p == NULL) {
    OAILOG_ERROR(
      LOG_MME_APP,
      "Invalid UE context received, MME UE S1AP Id: " MME_UE_S1AP_ID_FMT "\n",
      mme_ue_s1ap_id);
    OAILOG_FUNC_OUT(LOG_MME_APP);
  }
  ue_context_p->initial_context_setup_rsp_timer.id = MME_APP_TIMER_INACTIVE_ID;
  /* *********Abort the ongoing procedure*********
   * Check if UE is registered already that implies service request procedure is active. If so then release the S1AP
   * context and move the UE back to idle mode. Otherwise if UE is not yet registered that implies attach procedure is
   * active. If so,then abort the attach procedure and release the UE context.
   */
  ue_context_p->ue_context_rel_cause = S1AP_INITIAL_CONTEXT_SETUP_TMR_EXPRD;
  if (ue_context_p->mm_state == UE_UNREGISTERED) {
    // Initiate Implicit Detach for the UE
    nas_proc_implicit_detach_ue_ind(mme_ue_s1ap_id);
    increment_counter(
      "ue_attach",
      1,
      2,
      "result",
      "failure",
      "cause",
      "no_context_setup_rsp_from_enb");
    increment_counter("ue_attach", 1, 1, "action", "attach_abort");
  } else {
    // Release S1-U bearer and move the UE to idle mode
    for (pdn_cid_t i = 0; i < MAX_APN_PER_UE; i++) {
      if (ue_context_p->pdn_contexts[i]) {
        mme_app_send_s11_release_access_bearers_req(ue_context_p, i);
      }
    }
    /* Handles CSFB failure */
    if (ue_context_p->sgs_context != NULL) {
      handle_csfb_s1ap_procedure_failure(
        ue_context_p,
        "initial_context_setup_timer_expired",
        INTIAL_CONTEXT_SETUP_PROCEDURE_FAILED);
    }
  }
  OAILOG_FUNC_OUT(LOG_MME_APP);
}
//------------------------------------------------------------------------------
void mme_app_handle_initial_context_setup_failure(
  const itti_mme_app_initial_context_setup_failure_t
    *const initial_ctxt_setup_failure_pP)
{
  struct ue_mm_context_s *ue_context_p = NULL;

  OAILOG_FUNC_IN(LOG_MME_APP);
  OAILOG_INFO(
    LOG_MME_APP, "Received MME_APP_INITIAL_CONTEXT_SETUP_FAILURE from S1AP\n");
  ue_context_p = mme_ue_context_exists_mme_ue_s1ap_id(
    initial_ctxt_setup_failure_pP->mme_ue_s1ap_id);

  if (ue_context_p == NULL) {
    OAILOG_DEBUG(
      LOG_MME_APP,
      "We didn't find this mme_ue_s1ap_id in list of UE: %d \n",
      initial_ctxt_setup_failure_pP->mme_ue_s1ap_id);
    OAILOG_FUNC_OUT(LOG_MME_APP);
  }
  increment_counter("initial_context_setup_failure_received", 1, NO_LABELS);
  // Stop Initial context setup process guard timer,if running
  if (
    ue_context_p->initial_context_setup_rsp_timer.id !=
    MME_APP_TIMER_INACTIVE_ID) {
    if (timer_remove(ue_context_p->initial_context_setup_rsp_timer.id, NULL)) {
      OAILOG_ERROR(
        LOG_MME_APP,
        "Failed to stop Initial Context Setup Rsp timer for UE id  %d \n",
        ue_context_p->mme_ue_s1ap_id);
    }
    ue_context_p->initial_context_setup_rsp_timer.id =
      MME_APP_TIMER_INACTIVE_ID;
  }
  /* *********Abort the ongoing procedure*********
   * Check if UE is registered already that implies service request procedure is active. If so then release the S1AP
   * context and move the UE back to idle mode. Otherwise if UE is not yet registered that implies attach procedure is
   * active. If so,then abort the attach procedure and release the UE context.
   */
  ue_context_p->ue_context_rel_cause = S1AP_INITIAL_CONTEXT_SETUP_FAILED;
  if (ue_context_p->mm_state == UE_UNREGISTERED) {
    // Initiate Implicit Detach for the UE
    nas_proc_implicit_detach_ue_ind(ue_context_p->mme_ue_s1ap_id);
    increment_counter(
      "ue_attach",
      1,
      2,
      "result",
      "failure",
      "cause",
      "initial_context_setup_failure_rcvd");
    increment_counter("ue_attach", 1, 1, "action", "attach_abort");
  } else {
    // Release S1-U bearer and move the UE to idle mode

    for (pdn_cid_t i = 0; i < MAX_APN_PER_UE; i++) {
      if (ue_context_p->pdn_contexts[i]) {
        mme_app_send_s11_release_access_bearers_req(ue_context_p, i);
      }
    }
    /* Handles CSFB failure */
    if (ue_context_p->sgs_context != NULL) {
      handle_csfb_s1ap_procedure_failure(
        ue_context_p,
        "initial_context_setup_failed",
        INTIAL_CONTEXT_SETUP_PROCEDURE_FAILED);
    }
  }
  OAILOG_FUNC_OUT(LOG_MME_APP);
}
//------------------------------------------------------------------------------
static bool mme_app_construct_guti(
  const plmn_t *const plmn_p,
  const s_tmsi_t *const s_tmsi_p,
  guti_t *const guti_p)
{
  /*
   * This is a helper function to construct GUTI from S-TMSI. It uses PLMN id and MME Group Id of the serving MME for
   * this purpose.
   *
   */

  bool is_guti_valid =
    false; // Set to true if serving MME is found and GUTI is constructed
  uint8_t num_mme = 0; // Number of configured MME in the MME pool
  guti_p->m_tmsi = s_tmsi_p->m_tmsi;
  guti_p->gummei.mme_code = s_tmsi_p->mme_code;
  // Create GUTI by using PLMN Id and MME-Group Id of serving MME
  OAILOG_DEBUG(
    LOG_MME_APP,
    "Construct GUTI using S-TMSI received form UE and MME Group Id and PLMN id "
    "from MME Conf: %u, %u \n",
    s_tmsi_p->m_tmsi,
    s_tmsi_p->mme_code);
  mme_config_read_lock(&mme_config);
  /*
   * Check number of MMEs in the pool.
   * At present it is assumed that one MME is supported in MME pool but in case there are more
   * than one MME configured then search the serving MME using MME code.
   * Assumption is that within one PLMN only one pool of MME will be configured
   */
  if (mme_config.gummei.nb > 1) {
    OAILOG_DEBUG(LOG_MME_APP, "More than one MMEs are configured.");
  }
  for (num_mme = 0; num_mme < mme_config.gummei.nb; num_mme++) {
    /*Verify that the MME code within S-TMSI is same as what is configured in MME conf*/
    if (
      (plmn_p->mcc_digit2 ==
       mme_config.gummei.gummei[num_mme].plmn.mcc_digit2) &&
      (plmn_p->mcc_digit1 ==
       mme_config.gummei.gummei[num_mme].plmn.mcc_digit1) &&
      (plmn_p->mnc_digit3 ==
       mme_config.gummei.gummei[num_mme].plmn.mnc_digit3) &&
      (plmn_p->mcc_digit3 ==
       mme_config.gummei.gummei[num_mme].plmn.mcc_digit3) &&
      (plmn_p->mnc_digit2 ==
       mme_config.gummei.gummei[num_mme].plmn.mnc_digit2) &&
      (plmn_p->mnc_digit1 ==
       mme_config.gummei.gummei[num_mme].plmn.mnc_digit1) &&
      (guti_p->gummei.mme_code == mme_config.gummei.gummei[num_mme].mme_code)) {
      break;
    }
  }
  if (num_mme >= mme_config.gummei.nb) {
    OAILOG_DEBUG(LOG_MME_APP, "No MME serves this UE");
  } else {
    guti_p->gummei.plmn = mme_config.gummei.gummei[num_mme].plmn;
    guti_p->gummei.mme_gid = mme_config.gummei.gummei[num_mme].mme_gid;
    is_guti_valid = true;
  }
  mme_config_unlock(&mme_config);
  return is_guti_valid;
}

//------------------------------------------------------------------------------
static void notify_s1ap_new_ue_mme_s1ap_id_association(
  struct ue_mm_context_s *ue_context_p)
{
  MessageDef *message_p = NULL;
  itti_mme_app_s1ap_mme_ue_id_notification_t *notification_p = NULL;

  OAILOG_FUNC_IN(LOG_MME_APP);
  if (ue_context_p == NULL) {
    OAILOG_ERROR(LOG_MME_APP, " NULL UE context pointer!\n");
    OAILOG_FUNC_OUT(LOG_MME_APP);
  }
  message_p =
    itti_alloc_new_message(TASK_MME_APP, MME_APP_S1AP_MME_UE_ID_NOTIFICATION);
  notification_p = &message_p->ittiMsg.mme_app_s1ap_mme_ue_id_notification;
  memset(notification_p, 0, sizeof(itti_mme_app_s1ap_mme_ue_id_notification_t));
  notification_p->enb_ue_s1ap_id = ue_context_p->enb_ue_s1ap_id;
  notification_p->mme_ue_s1ap_id = ue_context_p->mme_ue_s1ap_id;
  notification_p->sctp_assoc_id = ue_context_p->sctp_assoc_id_key;

  itti_send_msg_to_task(TASK_S1AP, INSTANCE_DEFAULT, message_p);
  OAILOG_DEBUG(
    LOG_MME_APP,
    " Sent MME_APP_S1AP_MME_UE_ID_NOTIFICATION to S1AP for (ue_id = %u)\n",
    notification_p->mme_ue_s1ap_id);
  OAILOG_FUNC_OUT(LOG_MME_APP);
}

/**
 * Helper function to send a paging request to S1AP in either the initial case
 * or the retransmission case.
 *
 * @param ue_context_p - Pointer to UE context
 * @param set_timer - set true if this is the first attempt at paging and false
 *                    if this is the retransmission
 * @param paging_id_stmsi- paging ID, either to page with IMSI or STMSI
 * @param domain_indicator- Informs paging initiated for CS/PS
 */
int mme_app_paging_request_helper(
  ue_mm_context_t *ue_context_p,
  bool set_timer,
  uint8_t paging_id_stmsi,
  s1ap_cn_domain_t domain_indicator)
{
  MessageDef *message_p = NULL;
  int rc = RETURNok;
  OAILOG_FUNC_IN(LOG_MME_APP);
  // First, check if the UE is already connected. If so, stop
  if (ue_context_p->ecm_state == ECM_CONNECTED) {
    OAILOG_ERROR(
      LOG_MME_APP,
      "Paging process attempted for connected UE with id %d\n",
      ue_context_p->mme_ue_s1ap_id);
    OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNerror);
  }
  message_p = itti_alloc_new_message(TASK_MME_APP, S1AP_PAGING_REQUEST);
  if (message_p == NULL) {
    OAILOG_ERROR(
      LOG_MME_APP,
      "Failed to allocate the memory for paging request message\n");
    OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNerror);
  }
  itti_s1ap_paging_request_t *paging_request =
    &message_p->ittiMsg.s1ap_paging_request;
  memset(paging_request, 0, sizeof(itti_s1ap_paging_request_t));

  // @TODO Check
  IMSI64_TO_STRING(
    ue_context_p->emm_context._imsi64,
    (char*) paging_request->imsi,
    ue_context_p->emm_context._imsi.length);
  paging_request->imsi_length = ue_context_p->emm_context._imsi.length;
  paging_request->mme_code = ue_context_p->emm_context._guti.gummei.mme_code;
  paging_request->m_tmsi = ue_context_p->emm_context._guti.m_tmsi;
  // TODO Pass enb ids based on TAIs
  paging_request->sctp_assoc_id = ue_context_p->sctp_assoc_id_key;
  if (paging_id_stmsi) {
    paging_request->paging_id = S1AP_PAGING_ID_STMSI;
  } else {
    paging_request->paging_id = S1AP_PAGING_ID_IMSI;
  }
  paging_request->domain_indicator = domain_indicator;

  // Send TAI List
  paging_request->tai_list_count =
    ue_context_p->emm_context._tai_list.numberoflists;
  tai_list_t* tai_list = &ue_context_p->emm_context._tai_list;
  paging_tai_list_t* p_tai_list = NULL;
  for (int tai_list_idx = 0; tai_list_idx < paging_request->tai_list_count;
       tai_list_idx++) {
    p_tai_list = &paging_request->paging_tai_list[tai_list_idx];
    mme_app_update_paging_tai_list(
      p_tai_list,
      &tai_list->partial_tai_list[tai_list_idx],
      tai_list->partial_tai_list[tai_list_idx].numberofelements);
  }
  message_p->ittiMsgHeader.imsi = ue_context_p->emm_context._imsi64;
  rc = itti_send_msg_to_task(TASK_S1AP, INSTANCE_DEFAULT, message_p);

  if (!set_timer) {
    OAILOG_FUNC_RETURN(LOG_MME_APP, rc);
  }
  nas_itti_timer_arg_t timer_callback_fun = {0};
  timer_callback_fun.nas_timer_callback =
    (void *) mme_app_handle_paging_timer_expiry;
  timer_callback_fun.nas_timer_callback_arg =
    (void *) &(ue_context_p->mme_ue_s1ap_id);
  int timer_rc = timer_setup(
    ue_context_p->paging_response_timer.sec,
    0,
    TASK_MME_APP,
    INSTANCE_DEFAULT,
    TIMER_ONE_SHOT,
    &timer_callback_fun,
    sizeof(timer_callback_fun),
    &(ue_context_p->paging_response_timer.id));
  if (timer_rc < 0) {
    OAILOG_ERROR(
      LOG_MME_APP,
      "Failed to start paging timer for ue %d\n",
      ue_context_p->mme_ue_s1ap_id);
  }
  OAILOG_FUNC_RETURN(LOG_MME_APP, timer_rc);
}

int mme_app_handle_initial_paging_request(mme_app_desc_t *mme_app_desc_p,
    const char *imsi)
{
  imsi64_t imsi64 = INVALID_IMSI64;
  ue_mm_context_t *ue_context_p = NULL;

  IMSI_STRING_TO_IMSI64(imsi, &imsi64);
  ue_context_p =
    mme_ue_context_exists_imsi(&mme_app_desc_p->mme_ue_contexts, imsi64);
  if (ue_context_p == NULL) {
    OAILOG_ERROR(
      LOG_MME_APP, "Unknown IMSI%s, could not initiate paging\n", imsi);
    mme_ue_context_dump_coll_keys(&mme_app_desc_p->mme_ue_contexts);
    OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNerror);
  }
  return mme_app_paging_request_helper(
    ue_context_p, true, true /* s-tmsi */, CN_DOMAIN_PS);
}

void mme_app_handle_paging_timer_expiry(void* args)
{
  OAILOG_FUNC_IN(LOG_MME_APP);
  mme_ue_s1ap_id_t mme_ue_s1ap_id = *((mme_ue_s1ap_id_t*) (args));
  struct ue_mm_context_s* ue_context_p = mme_app_get_ue_context_for_timer(
    mme_ue_s1ap_id,
    "Paging timer");

  if (ue_context_p == NULL) {
    OAILOG_FUNC_OUT(LOG_MME_APP);
  }
  ue_context_p->paging_response_timer.id = MME_APP_TIMER_INACTIVE_ID;
  if ((mme_app_paging_request_helper(
    ue_context_p, false, true /* s-tmsi */, CN_DOMAIN_PS)) != RETURNok) {
    OAILOG_ERROR(
      LOG_MME_APP,
      "Failed to send Paging Message for ue_id " MME_UE_S1AP_ID_FMT "\n",
      mme_ue_s1ap_id);
  }
  OAILOG_FUNC_OUT(LOG_MME_APP);
}

void mme_app_handle_ulr_timer_expiry(void* args)
{
  OAILOG_FUNC_IN(LOG_MME_APP);
  mme_ue_s1ap_id_t mme_ue_s1ap_id = *((mme_ue_s1ap_id_t*) (args));
  struct ue_mm_context_s* ue_context_p = mme_app_get_ue_context_for_timer(
    mme_ue_s1ap_id,
    "Update location timer");
  if (ue_context_p == NULL) {
    OAILOG_FUNC_OUT(LOG_MME_APP);
  }
  ue_context_p->ulr_response_timer.id = MME_APP_TIMER_INACTIVE_ID;

  // Send PDN CONNECTIVITY FAIL message  to NAS layer
  increment_counter("mme_s6a_update_location_ans", 1, 1, "result", "failure");
  emm_cn_ula_or_csrsp_fail_t cn_ula_fail = {0};
  cn_ula_fail.ue_id = ue_context_p->mme_ue_s1ap_id;
  cn_ula_fail.cause = CAUSE_SYSTEM_FAILURE;
  for (pdn_cid_t i = 0; i < MAX_APN_PER_UE; i++) {
    if (ue_context_p->pdn_contexts[i]) {
      bearer_context_t* bearer_context = mme_app_get_bearer_context(
        ue_context_p, ue_context_p->pdn_contexts[i]->default_ebi);
      cn_ula_fail.pti = bearer_context->transaction_identifier;
      break;
    }
  }
  nas_proc_ula_or_csrsp_fail(&cn_ula_fail);
  OAILOG_FUNC_OUT(LOG_MME_APP);
}

/**
 * Send Suspend Notification to inform SPGW that UE is not available for PS handover
 * and discard the DL data received for this UE
 *
 * */
int mme_app_send_s11_suspend_notification(
  struct ue_mm_context_s *const ue_context_pP,
  const pdn_cid_t pdn_index)
{
  MessageDef *message_p = NULL;
  itti_s11_suspend_notification_t *suspend_notification_p = NULL;
  int rc = RETURNok;

  OAILOG_FUNC_IN(LOG_MME_APP);
  if (ue_context_pP == NULL) {
    OAILOG_ERROR(LOG_MME_APP, "Invalid UE context received\n");
    OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNerror);
  }
  OAILOG_DEBUG(
    LOG_MME_APP,
    "Preparing to send Suspend Notification for imsi " IMSI_64_FMT "\n",
    ue_context_pP->emm_context._imsi64);

  message_p = itti_alloc_new_message(TASK_MME_APP, S11_SUSPEND_NOTIFICATION);
  if (message_p == NULL) {
    OAILOG_ERROR(
      LOG_MME_APP,
      "Failed to allocate new ITTI message for S11 Suspend Notification "
      "for IMSI: " IMSI_64_FMT "\n",
      ue_context_pP->emm_context._imsi64);
    OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNerror);
  }

  suspend_notification_p = &message_p->ittiMsg.s11_suspend_notification;
  memset(suspend_notification_p, 0, sizeof(itti_s11_suspend_notification_t));

  pdn_context_t* pdn_connection = ue_context_pP->pdn_contexts[pdn_index];
  suspend_notification_p->teid = pdn_connection->s_gw_teid_s11_s4;

  IMSI64_TO_STRING(
    ue_context_pP->emm_context._imsi64,
    (char*) suspend_notification_p->imsi.digit,
    ue_context_pP->emm_context._imsi.length);
  suspend_notification_p->imsi.length =
    (uint8_t) strlen((const char*) suspend_notification_p->imsi.digit);

  /* lbi: currently one default bearer, fill lbi from UE context
   * TODO for multiple PDN support, get lbi from PDN context
  */
  suspend_notification_p->lbi = ue_context_pP->pdn_contexts[pdn_index]->default_ebi;

  message_p->ittiMsgHeader.imsi = ue_context_pP->emm_context._imsi64;

  OAILOG_INFO(
    LOG_MME_APP,
    "Send Suspend Notification for IMSI = " IMSI_64_FMT "\n",
    ue_context_pP->emm_context._imsi64);
  rc = itti_send_msg_to_task(TASK_SPGW, INSTANCE_DEFAULT, message_p);

  OAILOG_FUNC_RETURN(LOG_MME_APP, rc);
}

/*
 * Handle Suspend Acknowledge from SPGW
 *
 */
void mme_app_handle_suspend_acknowledge(mme_app_desc_t *mme_app_desc_p,
  const itti_s11_suspend_acknowledge_t *const suspend_acknowledge_pP)
{
  struct ue_mm_context_s *ue_context_p = NULL;

  OAILOG_FUNC_IN(LOG_MME_APP);
  OAILOG_INFO(
    LOG_MME_APP,
    "Rx Suspend Acknowledge with MME_S11_TEID :%d \n",
    suspend_acknowledge_pP->teid);

  ue_context_p = mme_ue_context_exists_s11_teid(
    &mme_app_desc_p->mme_ue_contexts, suspend_acknowledge_pP->teid);
  if (ue_context_p == NULL) {
    OAILOG_ERROR(
      LOG_MME_APP,
      "We didn't find this teid in list of UE: %" PRIX32 "\n",
      suspend_acknowledge_pP->teid);
    OAILOG_FUNC_OUT(LOG_MME_APP);
  }
  OAILOG_DEBUG(
    LOG_MME_APP,
    " Rx Suspend Acknowledge with MME_S11_TEID " TEID_FMT " IMSI " IMSI_64_FMT
    " \n",
    suspend_acknowledge_pP->teid,
    ue_context_p->emm_context._imsi64);
  /*
   * Updating statistics
   */
  update_mme_app_stats_s1u_bearer_sub();

  // Send UE Context Release Command
  mme_app_itti_ue_context_release(
    ue_context_p, ue_context_p->ue_context_rel_cause);
  OAILOG_FUNC_OUT(LOG_MME_APP);
}

//------------------------------------------------------------------------------
int mme_app_handle_nas_extended_service_req(
  const mme_ue_s1ap_id_t ue_id,
  const uint8_t service_type,
  uint8_t csfb_response)
{
  struct ue_mm_context_s *ue_context_p = NULL;
  int rc = RETURNok;

  OAILOG_FUNC_IN(LOG_MME_APP);

  if (ue_id == INVALID_MME_UE_S1AP_ID) {
    OAILOG_ERROR(
      LOG_MME_APP,
      "ERROR***** Invalid UE Id received in Extended Service Request \n");
    OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNerror);
  }
  ue_context_p = mme_ue_context_exists_mme_ue_s1ap_id(ue_id);
  if (ue_context_p) {
    if (ue_id != ue_context_p->mme_ue_s1ap_id) {
      OAILOG_ERROR(
        LOG_MME_APP,
        "ERROR***** Abnormal case: ue_id does not match with ue_id in "
        "ue_context" MME_UE_S1AP_ID_FMT "," MME_UE_S1AP_ID_FMT "\n",
        ue_id,
        ue_context_p->mme_ue_s1ap_id);
      OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNerror);
    }
  } else {
    OAILOG_ERROR(
      LOG_MME_APP,
      "ERROR***** Invalid UE Id received from NAS in Extended Service Request "
      "%d\n",
      ue_id);
    OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNerror);
  }

  switch (service_type) {
    /* Extended Service request received for CSFB */
    case MO_CS_FB1:
    case MT_CS_FB1:
    case MO_CS_FB:
      if (ue_context_p->sgs_context != NULL) {
        ue_context_p->sgs_context->csfb_service_type = CSFB_SERVICE_MO_CALL;
        /* If call_cancelled is set to TRUE when MO call is triggered.
         * Set call_cancelled to false
         */
        if (ue_context_p->sgs_context->call_cancelled) {
          ue_context_p->sgs_context->call_cancelled = false;
        }
        mme_app_itti_ue_context_mod_for_csfb(ue_context_p);
      } else {
        OAILOG_ERROR(
          LOG_MME_APP,
          "SGS context is NULL for ue_id:" MME_UE_S1AP_ID_FMT
          "So send Service Reject to UE \n",
          ue_context_p->mme_ue_s1ap_id);
        /* send Service Reject to UE */
        mme_app_notify_service_reject_to_nas(
          ue_context_p->mme_ue_s1ap_id,
          EMM_CAUSE_CONGESTION,
          UE_CONTEXT_MODIFICATION_PROCEDURE_FAILED);
      }
      break;
    case MT_CS_FB:
      if (csfb_response == CSFB_REJECTED_BY_UE) {
        if (ue_context_p->sgs_context) {
          /* If call_cancelled is set to TRUE and
           * receive EXT Service Request with csfb_response
           * set to call_rejected. Set call_cancelled to false
           */
          if (ue_context_p->sgs_context->call_cancelled) {
            ue_context_p->sgs_context->call_cancelled = false;
          }
          rc = mme_app_send_sgsap_paging_reject(
               ue_context_p,
               ue_context_p->emm_context._imsi64,
               ue_context_p->emm_context._imsi.length,
               SGS_CAUSE_MT_CSFB_CALL_REJECTED_BY_USER);
          if (rc != RETURNok) {
            OAILOG_WARNING(
              LOG_MME_APP,
              "Failed to send SGSAP-Paging Reject for imsi with reject cause:"
              "SGS_CAUSE_MT_CSFB_CALL_REJECTED_BY_USER" IMSI_64_FMT "\n",
              ue_context_p->emm_context._imsi64);
          }
          increment_counter(
            "sgsap_paging_reject", 1, 1, "cause", "call_rejected_by_user");
        } else {
          OAILOG_ERROR(
            LOG_MME_APP,
            "sgs_context is null for IMSI" IMSI_64_FMT "\n",
            ue_context_p->emm_context._imsi64);
           OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNerror);
        }
      } else if (
        csfb_response == CSFB_ACCEPTED_BY_UE) {
        if (!ue_context_p->sgs_context) {
          OAILOG_ERROR(
            LOG_MME_APP,
            "sgs_context is null for IMSI" IMSI_64_FMT "\n",
            ue_context_p->emm_context._imsi64);
           OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNerror);
        }
        /* Set mt_call_in_progress flag as UE accepted the MT Call.
         * This will be used to decide whether to abort the on going MT call or
         * not when SERVICE ABORT request is received from MSC/VLR
         */
        ue_context_p->sgs_context->mt_call_in_progress = true;
        /* If call_cancelled is set, send Service Reject to UE as MSC/VLR
         * has triggered SGSAP SERVICE ABORT procedure
         */
        if (ue_context_p->sgs_context->call_cancelled) {
          /* If UE's ECM state is IDLE send
           * service_reject in Establish cnf else send in DL NAS Transport
           */
          if (ue_context_p->ecm_state == ECM_IDLE) {
            OAILOG_ERROR(
              LOG_MME_APP,
              "MT CS call is accepted by UE in idle mode for ue_id:"
              MME_UE_S1AP_ID_FMT " But MT_CALL_CANCEL is set by MSC,"
             " so sending service reject to UE \n",
              ue_id);
            mme_app_notify_service_reject_to_nas(
              ue_id,
              EMM_CAUSE_CS_SERVICE_NOT_AVAILABLE,
              INTIAL_CONTEXT_SETUP_PROCEDURE_FAILED);
          } else if (ue_context_p->ecm_state == ECM_CONNECTED) {
            OAILOG_ERROR(
              LOG_MME_APP,
              "MT CS call is accepted by UE in connected mode for ue_id:"
              MME_UE_S1AP_ID_FMT " But MT_CALL_CANCEL is set by MSC,"
             " so sending service reject to UE \n",
              ue_id);
            mme_app_notify_service_reject_to_nas(
              ue_id,
              EMM_CAUSE_CS_SERVICE_NOT_AVAILABLE,
              UE_CONTEXT_MODIFICATION_PROCEDURE_FAILED);
          }
          // Reset call_cancelled flag
          ue_context_p->sgs_context->call_cancelled = false;
          OAILOG_WARNING(
            LOG_MME_APP,
            "Sending Service Reject to NAS module as MSC has triggered SGS "
            "SERVICE ABORT Request for ue_id: " MME_UE_S1AP_ID_FMT "\n",
            ue_id);
        } else {
          mme_app_itti_ue_context_mod_for_csfb(ue_context_p);
        }
      } else {
        OAILOG_WARNING(
          LOG_MME_APP,
          "Invalid csfb_response for service type :%d and ue_id: "
          MME_UE_S1AP_ID_FMT "\n",
          service_type,
          ue_id);
      }
      break;
    case MO_CS_FB_EMRGNCY_CALL:
      if (ue_context_p->sgs_context != NULL) {
        ue_context_p->sgs_context->csfb_service_type = CSFB_SERVICE_MO_CALL;
        ue_context_p->sgs_context->is_emergency_call = true;
        mme_app_itti_ue_context_mod_for_csfb(ue_context_p);
      } else {
        // Notify NAS module to send Service Reject message to UE
        OAILOG_ERROR(
          LOG_MME_APP,
          "For MO_CS_FB_EMRGNCY_CALL, SGS context is not found for ue_id:"
          MME_UE_S1AP_ID_FMT " MME shall send Service Reject to ue",
          ue_context_p->mme_ue_s1ap_id);
        mme_app_notify_service_reject_to_nas(
          ue_context_p->mme_ue_s1ap_id,
          EMM_CAUSE_CONGESTION,
          UE_CONTEXT_MODIFICATION_PROCEDURE_FAILED);
      }
      break;
    /* packet service via s1 */
    case PKT_SRV_VIA_S1:
    case PKT_SRV_VIA_S1_1:
    case PKT_SRV_VIA_S1_2:
    case PKT_SRV_VIA_S1_3:
      /*TODO */
      break;
    default:
      OAILOG_ERROR(
        LOG_MME_APP,
        "ERROR***** Invalid Service Type Received %d\n",
        service_type);
  }
  OAILOG_FUNC_RETURN(LOG_MME_APP, rc);
}

//------------------------------------------------------------------------------
void mme_app_handle_ue_context_modification_timer_expiry(void* args)
{
  OAILOG_FUNC_IN(LOG_MME_APP);
  mme_ue_s1ap_id_t mme_ue_s1ap_id = *((mme_ue_s1ap_id_t *) (args));
  struct ue_mm_context_s* ue_context_p = mme_app_get_ue_context_for_timer(
    mme_ue_s1ap_id,
    "UE context modification timer");
  if (ue_context_p == NULL) {
    OAILOG_ERROR(
      LOG_MME_APP,
      "Invalid UE context received, MME UE S1AP Id: " MME_UE_S1AP_ID_FMT "\n",
      mme_ue_s1ap_id);
    OAILOG_FUNC_OUT(LOG_MME_APP);
  }
  ue_context_p->ue_context_modification_timer.id = MME_APP_TIMER_INACTIVE_ID;

  if (ue_context_p->sgs_context != NULL) {
    handle_csfb_s1ap_procedure_failure(
      ue_context_p,
      "ue_context_modification_timer_expired",
      UE_CONTEXT_MODIFICATION_PROCEDURE_FAILED);
  }
  OAILOG_FUNC_OUT(LOG_MME_APP);
}

/* Description: CSFB procedure to handle S1ap procedure failure,
 * In case of MT CS call, send SGSAP Paging reject to MSC/VLR
 * And Send Service Reject to UE
 * In case of of MO CS call, send Service Reject to UE
 */
int handle_csfb_s1ap_procedure_failure(
  ue_mm_context_t *ue_context_p,
  char *failed_statement,
  uint8_t failed_procedure)
{
  OAILOG_FUNC_IN(LOG_MME_APP);

  if (!ue_context_p) {
    OAILOG_ERROR(LOG_MME_APP, "Failed to find UE context \n");
    OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNerror);
  }

  OAILOG_ERROR(
    LOG_MME_APP,
    "Handle handle_csfb_s1ap_procedure_failure for ue_id"
    MME_UE_S1AP_ID_FMT "\n", ue_context_p->mme_ue_s1ap_id);
  /* If ICS procedure is initiated due to CS-Paging in UE idle mode
   * On ICS failure, send sgsap-Paging Reject to VLR
   */
  if (ue_context_p->sgs_context) {
    //Reset mt_call_in_progress flag
    if (ue_context_p->sgs_context->mt_call_in_progress) {
      ue_context_p->sgs_context->mt_call_in_progress = false;
    }
    if (ue_context_p->sgs_context->csfb_service_type == CSFB_SERVICE_MT_CALL) {
      /* send sgsap-Paging Reject to VLR */
      if (
        (mme_app_send_sgsap_paging_reject(
          ue_context_p,
          ue_context_p->emm_context._imsi64,
          ue_context_p->emm_context._imsi.length,
          SGS_CAUSE_MT_CSFB_CALL_REJECTED_BY_USER)) != RETURNok) {
        OAILOG_WARNING(
          LOG_MME_APP,
          "Failed to send SGSAP-Paging Reject for imsi with reject cause:"
          "SGS_CAUSE_MT_CSFB_CALL_REJECTED_BY_USER" IMSI_64_FMT "\n",
          ue_context_p->emm_context._imsi64);
      }
      if (failed_statement) {
        increment_counter(
          "sgsap_paging_reject", 1, 1, "cause", failed_statement);
      }
    }
    // send Service Reject to UE
    mme_app_notify_service_reject_to_nas(
      ue_context_p->mme_ue_s1ap_id,
      EMM_CAUSE_CONGESTION,
      failed_procedure);
    ue_context_p->sgs_context->csfb_service_type = CSFB_SERVICE_NONE;
    if (failed_statement) {
      increment_counter("nas service reject", 1, 1, "cause", failed_statement);
    }
  }
  OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNok);
}

/****************************************************************************
 **                                                                        **
 ** Name:    mme_app_notify_service_reject_to_nas()                        **
 **                                                                        **
 ** Description: As part of handling CSFB procedure, if ICS or UE context  **
 **      modification failed, indicate to NAS to send Service Reject to UE **
 **                                                                        **
 ** Inputs:  ue_id: UE identifier                                          **
 **          emm_casue: failed cause                                       **
 **          Failed_procedure: ICS/UE context modification                 **
 **                                                                        **
 ***************************************************************************/
void mme_app_notify_service_reject_to_nas(
  mme_ue_s1ap_id_t ue_id,
  uint8_t emm_cause,
  uint8_t failed_procedure)
{
  OAILOG_FUNC_IN(LOG_MME_APP);
  OAILOG_INFO(
    LOG_MME_APP,
    " Ongoing Service request procedure failed,"
    "send Notify Service Reject to NAS module for ue_id :"
    MME_UE_S1AP_ID_FMT" \n",
    ue_id);
  switch (failed_procedure) {
    case INTIAL_CONTEXT_SETUP_PROCEDURE_FAILED: {
      if ((emm_proc_service_reject(ue_id, emm_cause)) != RETURNok) {
        OAILOG_ERROR(
          LOG_MME_APP,
          "emm_proc_service_reject() failed for ue_id "MME_UE_S1AP_ID_FMT "\n",
          ue_id);
      }
      break;
    }
    case UE_CONTEXT_MODIFICATION_PROCEDURE_FAILED: {
      if ((emm_send_service_reject_in_dl_nas(ue_id, emm_cause)) != RETURNok) {
        OAILOG_ERROR(
          LOG_MME_APP,
          "emm_send_service_reject_in_dl_nas() failed for ue_id "
          MME_UE_S1AP_ID_FMT "\n",
          ue_id);
        }
      break;
    }
    default: {
      OAILOG_ERROR(LOG_MME_APP, "Invalid failed procedure for ue-id"
        MME_UE_S1AP_ID_FMT "\n", ue_id);
      break;
    }
  }
  OAILOG_FUNC_OUT(LOG_MME_APP);
}
//------------------------------------------------------------------------------
void mme_app_handle_create_dedicated_bearer_rsp(
  ue_mm_context_t* ue_context_p,
  ebi_t ebi)
{
  OAILOG_FUNC_IN(LOG_MME_APP);

#if EMBEDDED_SGW
  OAILOG_INFO(
    LOG_MME_APP,
    "Sending Activate Dedicated Bearer Response to SPGW for "
    "ue-id: " MME_UE_S1AP_ID_FMT "\n",
    ue_context_p->mme_ue_s1ap_id);
  _send_pcrf_bearer_actv_rsp(ue_context_p, ebi, REQUEST_ACCEPTED);
  OAILOG_FUNC_OUT(LOG_MME_APP);
#endif
  // TODO:
  /* Actually do it simple, because it appear we have to wait for NAS procedure
   * reworking (work in progress on another branch)
   * for responding to S11 without mistakes (may be the create bearer procedure
   * can be impacted by a S1 ue context release or
   * a UE originating  NAS procedure)
   */
  mme_app_s11_proc_create_bearer_t* s11_proc_create =
    mme_app_get_s11_procedure_create_bearer(ue_context_p);
  if (s11_proc_create) {
    s11_proc_create->num_status_received++;
    s11_proc_create->bearer_status[EBI_TO_INDEX(ebi)] = S11_PROC_BEARER_SUCCESS;
    // if received all bearers creation results
    if (s11_proc_create->num_status_received == s11_proc_create->num_bearers) {
      //Send Rsp to SGW if SPGW is embedded
      bearer_context_t* bc = mme_app_get_bearer_context(ue_context_p, ebi);
      if (bc == NULL) {
        OAILOG_ERROR(
          LOG_MME_APP, "Could not get bearer context for EBI:%d\n", ebi);
        OAILOG_FUNC_OUT(LOG_MME_APP);
      }
      mme_app_s11_procedure_create_bearer_send_response(
        ue_context_p, s11_proc_create);
      mme_app_delete_s11_procedure_create_bearer(ue_context_p);
    }
  }
  OAILOG_FUNC_OUT(LOG_MME_APP);
}

//------------------------------------------------------------------------------
void mme_app_handle_create_dedicated_bearer_rej(
  ue_mm_context_t* ue_context_p,
  ebi_t ebi)
{
  OAILOG_FUNC_IN(LOG_MME_APP);

  if (ue_context_p == NULL) {
    OAILOG_ERROR(
      LOG_MME_APP,
      "Failed to find UE context for ue-id :" MME_UE_S1AP_ID_FMT "\n",
      ue_context_p->mme_ue_s1ap_id);
    OAILOG_FUNC_OUT(LOG_MME_APP);
  }

#if EMBEDDED_SGW
  OAILOG_INFO(
    LOG_MME_APP,
    "Sending Activate Dedicated bearer Reject to SPGW: " MME_UE_S1AP_ID_FMT
    "\n",
    ue_context_p->mme_ue_s1ap_id);
  _send_pcrf_bearer_actv_rsp(ue_context_p, ebi, REQUEST_REJECTED);
  OAILOG_FUNC_OUT(LOG_MME_APP);
#endif

  // TODO:
  /* Actually do it simple, because it appear we have to wait for NAS procedure
   * reworking (work in progress on another branch)
   * for responding to S11 without mistakes (may be the create bearer procedure
   * can be impacted by a S1 ue context release or
   * a UE originating  NAS procedure)
   */
  mme_app_s11_proc_create_bearer_t *s11_proc_create =
    mme_app_get_s11_procedure_create_bearer(ue_context_p);
  if (s11_proc_create) {
    s11_proc_create->num_status_received++;
    s11_proc_create->bearer_status[EBI_TO_INDEX(ebi)] = S11_PROC_BEARER_FAILED;
    // if received all bearers creation results
    if (s11_proc_create->num_status_received == s11_proc_create->num_bearers) {
      mme_app_s11_procedure_create_bearer_send_response(
        ue_context_p, s11_proc_create);
      mme_app_delete_s11_procedure_create_bearer(ue_context_p);
    }
  }
  OAILOG_FUNC_OUT(LOG_MME_APP);
}

//------------------------------------------------------------------------------
// See 3GPP TS 23.401 version 10.13.0 Release 10: 5.4.4.2 MME Initiated Dedicated Bearer Deactivation
void mme_app_trigger_mme_initiated_dedicated_bearer_deactivation_procedure(
    mme_app_desc_t *mme_app_desc_p,
    ue_mm_context_t *const ue_context,
    const pdn_cid_t cid)
{
  OAILOG_DEBUG(LOG_MME_APP, "TODO \n");
}

/**
 * This Function checks for ue context based on given teid,
 * if present send ue context modification request to S1AP
 * otherwise drop the message
 */
void mme_app_handle_modify_ue_ambr_request(mme_app_desc_t *mme_app_desc_p,
  const itti_s11_modify_ue_ambr_request_t *const modify_ue_ambr_request_p)
{
  MessageDef *message_p;
  ue_mm_context_t *ue_context_p = NULL;
  OAILOG_FUNC_IN(LOG_MME_APP);

  ue_context_p = mme_ue_context_exists_s11_teid(
    &mme_app_desc_p->mme_ue_contexts, modify_ue_ambr_request_p->teid);

  if (ue_context_p == NULL) {
    OAILOG_WARNING(
      LOG_MME_APP,
      "We didn't find this teid in list of UE: \
        %08x\n, Dropping MODIFY_UE_AMBR_REQUEST",
      modify_ue_ambr_request_p->teid);
    OAILOG_FUNC_OUT(LOG_MME_APP);
  } else {
    message_p = itti_alloc_new_message(
      TASK_MME_APP, S1AP_UE_CONTEXT_MODIFICATION_REQUEST);
    if (message_p == NULL) {
      OAILOG_ERROR(
        LOG_MME_APP,
        "Failed to allocate new ITTI message for S1AP UE Context Modification "
        "Request for MME UE S1AP Id: " MME_UE_S1AP_ID_FMT "\n",
        ue_context_p->mme_ue_s1ap_id);
      OAILOG_FUNC_OUT(LOG_MME_APP);
    }
    memset(
      (void *) &message_p->ittiMsg.s1ap_ue_context_mod_request,
      0,
      sizeof(itti_s1ap_ue_context_mod_req_t));
    S1AP_UE_CONTEXT_MODIFICATION_REQUEST(message_p).mme_ue_s1ap_id =
      ue_context_p->mme_ue_s1ap_id;
    S1AP_UE_CONTEXT_MODIFICATION_REQUEST(message_p).enb_ue_s1ap_id =
      ue_context_p->enb_ue_s1ap_id;
    S1AP_UE_CONTEXT_MODIFICATION_REQUEST(message_p).presencemask =
      S1AP_UE_CONTEXT_MOD_UE_AMBR_INDICATOR_PRESENT;
    S1AP_UE_CONTEXT_MODIFICATION_REQUEST(message_p).ue_ambr.br_ul =
      modify_ue_ambr_request_p->ue_ambr.br_ul;
    S1AP_UE_CONTEXT_MODIFICATION_REQUEST(message_p).ue_ambr.br_dl =
      modify_ue_ambr_request_p->ue_ambr.br_dl;

    message_p->ittiMsgHeader.imsi = ue_context_p->emm_context._imsi64;
    itti_send_msg_to_task(TASK_S1AP, INSTANCE_DEFAULT, message_p);
    OAILOG_DEBUG(
      LOG_MME_APP,
      "MME APP :Sent UE context modification request \
        for UE id %d\n",
      ue_context_p->mme_ue_s1ap_id);
  }
  OAILOG_FUNC_OUT(LOG_MME_APP);
}

/**
 * This Function handles NW initiated
 * Dedicated bearer activation Request message from SGW
 */
void mme_app_handle_nw_init_ded_bearer_actv_req(
  mme_app_desc_t* mme_app_desc_p,
  const itti_s11_nw_init_actv_bearer_request_t* const nw_init_bearer_actv_req_p)
{
  ue_mm_context_t* ue_context_p = NULL;
  emm_cn_activate_dedicated_bearer_req_t activate_ded_bearer_req = {0};
  OAILOG_FUNC_IN(LOG_MME_APP);

  ebi_t linked_eps_bearer_id = nw_init_bearer_actv_req_p->lbi;
  ue_context_p = mme_ue_context_exists_s11_teid(
    &mme_app_desc_p->mme_ue_contexts, nw_init_bearer_actv_req_p->s11_mme_teid);

  if (ue_context_p == NULL) {
    OAILOG_ERROR(
      LOG_MME_APP,
      "Failed to find UE context from S11-Teid: for lbi %08x %u\n",
      nw_init_bearer_actv_req_p->s11_mme_teid,
      nw_init_bearer_actv_req_p->lbi);
    OAILOG_FUNC_OUT(LOG_MME_APP);
  }

  OAILOG_INFO(
    LOG_MME_APP,
    "Received Dedicated bearer activation Request from SGW for "
    "ue-id " MME_UE_S1AP_ID_FMT " with LBI %u\n",
    ue_context_p->mme_ue_s1ap_id,
    nw_init_bearer_actv_req_p->lbi);

  bearer_context_t* linked_bc =
    mme_app_get_bearer_context(ue_context_p, linked_eps_bearer_id);
  if (!linked_bc) {
    OAILOG_ERROR(
      LOG_MME_APP,
      "Failed to find the default bearer context from linked bearer id %" PRIu8
      " of ue_id: " MME_UE_S1AP_ID_FMT "\n",
      linked_eps_bearer_id,
      ue_context_p->mme_ue_s1ap_id);
    OAILOG_FUNC_OUT(LOG_MME_APP);
  }
  pdn_cid_t cid = linked_bc->pdn_cx_id;
  activate_ded_bearer_req.ue_id = ue_context_p->mme_ue_s1ap_id;
  activate_ded_bearer_req.cid = cid;
  //EBI Will be assigned by NAS module
  activate_ded_bearer_req.ebi = 0;
  activate_ded_bearer_req.linked_ebi =
    ue_context_p->pdn_contexts[cid]->default_ebi;
  activate_ded_bearer_req.bearer_qos =
    nw_init_bearer_actv_req_p->eps_bearer_qos;
  memcpy(
    &activate_ded_bearer_req.sgw_fteid,
    &nw_init_bearer_actv_req_p->s1_u_sgw_fteid,
    sizeof(fteid_t));

  if (nw_init_bearer_actv_req_p->tft.numberofpacketfilters) {
    activate_ded_bearer_req.tft = calloc(1, sizeof(traffic_flow_template_t));
    copy_traffic_flow_template(
      activate_ded_bearer_req.tft, &nw_init_bearer_actv_req_p->tft);
  }
  if (nw_init_bearer_actv_req_p->pco.num_protocol_or_container_id) {
    activate_ded_bearer_req.pco =
      calloc(1, sizeof(protocol_configuration_options_t));
    copy_protocol_configuration_options(
      activate_ded_bearer_req.pco, &nw_init_bearer_actv_req_p->pco);
  }

  if (
    (nas_proc_create_dedicated_bearer(&activate_ded_bearer_req)) != RETURNok) {
    OAILOG_ERROR(
      LOG_MME_APP,
      "Failed to handle bearer activation at NAS module for "
      "ue_id " MME_UE_S1AP_ID_FMT "\n",
      ue_context_p->mme_ue_s1ap_id);
  }
  free_wrapper((void **) &activate_ded_bearer_req.tft);
  free_wrapper((void **) &activate_ded_bearer_req.pco);

  OAILOG_FUNC_OUT(LOG_MME_APP);
}

void send_delete_dedicated_bearer_rsp(
  struct ue_mm_context_s* ue_context_p,
  bool delete_default_bearer,
  ebi_t ebi[],
  uint32_t num_bearer_context,
  teid_t s_gw_teid_s11_s4,
  gtpv2c_cause_value_t cause)
{
  itti_s11_nw_init_deactv_bearer_rsp_t* s11_deact_ded_bearer_rsp = NULL;
  MessageDef* message_p = NULL;
  uint32_t i = 0;

  OAILOG_FUNC_IN(LOG_MME_APP);
  if (ue_context_p == NULL) {
    OAILOG_ERROR(LOG_MME_APP, " NULL UE context ptr\n");
    OAILOG_FUNC_OUT(LOG_MME_APP);
  }
  message_p = itti_alloc_new_message(
    TASK_MME_APP, S11_NW_INITIATED_DEACTIVATE_BEARER_RESP);
  s11_deact_ded_bearer_rsp = &message_p->ittiMsg.s11_nw_init_deactv_bearer_rsp;

  if (message_p == NULL) {
    OAILOG_ERROR(
      LOG_MME_APP,
      "itti_alloc_new_message failed for"
      "S11_NW_INITIATED_DEACTIVATE_BEARER_RESP\n");
    OAILOG_FUNC_OUT(LOG_MME_APP);
  }
  memset(
    s11_deact_ded_bearer_rsp, 0, sizeof(itti_s11_nw_init_deactv_bearer_rsp_t));

  s11_deact_ded_bearer_rsp->delete_default_bearer = delete_default_bearer;

  if (delete_default_bearer) {
    s11_deact_ded_bearer_rsp->lbi = calloc(1, sizeof(ebi_t));
    *s11_deact_ded_bearer_rsp->lbi = ebi[0];
    s11_deact_ded_bearer_rsp->bearer_contexts.bearer_contexts[0].
      cause.cause_value = cause;
  } else {
    for (i = 0; i < num_bearer_context; i++) {
      s11_deact_ded_bearer_rsp->bearer_contexts.bearer_contexts[i].
        eps_bearer_id = ebi[i];
      s11_deact_ded_bearer_rsp->bearer_contexts.bearer_contexts[i].
        cause.cause_value = cause;
    }
  }
  /*Print bearer ids to be sent in nw_initiated_deactv_bearer_rsp*/
  print_bearer_ids_helper(ebi, num_bearer_context);
  s11_deact_ded_bearer_rsp->bearer_contexts.num_bearer_context =
    num_bearer_context;
  s11_deact_ded_bearer_rsp->imsi = ue_context_p->emm_context._imsi64;
  s11_deact_ded_bearer_rsp->s_gw_teid_s11_s4 = s_gw_teid_s11_s4;

  message_p->ittiMsgHeader.imsi = ue_context_p->emm_context._imsi64;

  OAILOG_INFO(
    LOG_MME_APP,
    " Sending nw_initiated_deactv_bearer_rsp to SGW with %d bearers\n",
    num_bearer_context);
  itti_send_msg_to_task(TASK_SPGW, INSTANCE_DEFAULT, message_p);

  OAILOG_FUNC_OUT(LOG_MME_APP);
}


/**
 * This Function handles NW-initiated
 * dedicated bearer deactivation request message from SGW
 */
void mme_app_handle_nw_init_bearer_deactv_req(
  mme_app_desc_t* mme_app_desc_p,
  itti_s11_nw_init_deactv_bearer_request_t* const nw_init_bearer_deactv_req_p)
{
  ue_mm_context_t* ue_context_p = NULL;
  uint32_t i = 0;
  OAILOG_FUNC_IN(LOG_MME_APP);
  ebi_t ebi[BEARERS_PER_UE];
  uint32_t num_bearers_deleted = 0;
  emm_cn_deactivate_dedicated_bearer_req_t deactivate_ded_bearer_req = {0};

  OAILOG_INFO(
    LOG_MME_APP,
    "Received nw_initiated_deactv_bearer_req from SGW for S11 teid" TEID_FMT
    "\n",
    nw_init_bearer_deactv_req_p->s11_mme_teid);

  /*Print bearer ids received in the message*/
  print_bearer_ids_helper(
    nw_init_bearer_deactv_req_p->ebi,
    nw_init_bearer_deactv_req_p->no_of_bearers);

  ue_context_p = mme_ue_context_exists_s11_teid(
    &mme_app_desc_p->mme_ue_contexts,
    nw_init_bearer_deactv_req_p->s11_mme_teid);
  if (ue_context_p == NULL) {
    OAILOG_ERROR(
      LOG_MME_APP,
      "Failed to find UE context for S11 Teid" TEID_FMT "\n",
      nw_init_bearer_deactv_req_p->s11_mme_teid);
    OAILOG_FUNC_OUT(LOG_MME_APP);
  }

  //Fetch PDN context
  pdn_cid_t cid =
    ue_context_p
      ->bearer_contexts[EBI_TO_INDEX(nw_init_bearer_deactv_req_p->ebi[0])]
      ->pdn_cx_id;
  pdn_context_t* pdn_context = ue_context_p->pdn_contexts[cid];

  /* If delete_default_bearer is set and this is the only active PDN,
  *  Send Detach Request to UE
  */
  if (
    (nw_init_bearer_deactv_req_p->delete_default_bearer) &&
    (ue_context_p->nb_active_pdn_contexts == 1)) {
    OAILOG_INFO(
      LOG_MME_APP,
      "Send MME initiated Detach Req to NAS module for EBI %u"
      " as delete_default_bearer is true\n",
      nw_init_bearer_deactv_req_p->ebi[0]);
    //Inform MME initiated Deatch Request to NAS module
    if (ue_context_p->ecm_state == ECM_CONNECTED) {
      mme_app_handle_nw_initiated_detach_request(
        ue_context_p->mme_ue_s1ap_id, MME_INITIATED_EPS_DETACH);
    } else {
      //If UE is in IDLE state send Paging Req
      mme_app_paging_request_helper(
        ue_context_p, true, false /* s-tmsi */, CN_DOMAIN_PS);
      // Set the flag and send detach to UE after receiving service req
      ue_context_p->emm_context.nw_init_bearer_deactv = true;
    }
  } else {
    /* If UE is in connected state, MME shall send Deactivate Bearer Req
     * in S1ap ERAB Rel Cmd
     */
    if (ue_context_p->ecm_state == ECM_CONNECTED) {
      deactivate_ded_bearer_req.ue_id = ue_context_p->mme_ue_s1ap_id;
      deactivate_ded_bearer_req.no_of_bearers =
        nw_init_bearer_deactv_req_p->no_of_bearers;
      memcpy(
        deactivate_ded_bearer_req.ebi,
        nw_init_bearer_deactv_req_p->ebi,
        ((sizeof(ebi_t)) * deactivate_ded_bearer_req.no_of_bearers));
      if (
        (nas_proc_delete_dedicated_bearer(&deactivate_ded_bearer_req)) !=
        RETURNok) {
        OAILOG_ERROR(
          LOG_MME_APP,
          "Failed to handle bearer deactivation at NAS module for "
          "ue_id " MME_UE_S1AP_ID_FMT "\n",
          ue_context_p->mme_ue_s1ap_id);
      }
    } else {
      /* If UE is in IDLE state remove bearer context
       * send delete dedicated bearer rsp to SPGW
       */
      for (i = 0; i < nw_init_bearer_deactv_req_p->no_of_bearers; i++) {
        /*Fetch bearer context*/
        bearer_context_t* bearer_context = mme_app_get_bearer_context(
          ue_context_p, nw_init_bearer_deactv_req_p->ebi[i]);
        if (bearer_context) {
          mme_app_free_bearer_context(&bearer_context);
          num_bearers_deleted++;
          ebi[i] = nw_init_bearer_deactv_req_p->ebi[i];
        } else {
          OAILOG_ERROR(
            LOG_MME_APP,
            "Bearer context does not exist for ebi %d\n",
            nw_init_bearer_deactv_req_p->ebi[i]);
        }
      }
      //Send delete_dedicated_bearer_rsp to SPGW
      send_delete_dedicated_bearer_rsp(
        ue_context_p,
        nw_init_bearer_deactv_req_p->delete_default_bearer,
        ebi,
        num_bearers_deleted,
        pdn_context->s_gw_teid_s11_s4,
        REQUEST_ACCEPTED);
    }
  }
  OAILOG_FUNC_OUT(LOG_MME_APP);
}

void mme_app_handle_path_switch_request(mme_app_desc_t *mme_app_desc_p,
  itti_s1ap_path_switch_request_t *const path_switch_req_p)
{
  OAILOG_FUNC_IN(LOG_MME_APP);
  struct ue_mm_context_s *ue_context_p = NULL;
  ue_network_capability_t ue_network_capability;
  enb_s1ap_id_key_t enb_s1ap_id_key = INVALID_ENB_UE_S1AP_ID_KEY;
  e_rab_to_be_switched_in_downlink_list_t e_rab_to_be_switched_dl_list =
    path_switch_req_p->e_rab_to_be_switched_dl_list;
  bearer_context_t *current_bearer_p = NULL;
  ebi_t bearer_id = 0;
  pdn_cid_t cid = 0;
  int idx = 0;
  pdn_context_t *pdn_context = NULL;
  MessageDef *message_p = NULL;

  OAILOG_DEBUG(LOG_MME_APP, "Received PATH_SWITCH_REQUEST from S1AP\n");

  ue_context_p = mme_ue_context_exists_mme_ue_s1ap_id(
    path_switch_req_p->mme_ue_s1ap_id);
  if (!ue_context_p) {
    OAILOG_ERROR(
      LOG_MME_APP,
      " PATH_SWITCH_REQUEST RECEIVED, Failed to find UE context for "
      "mme_ue_s1ap_id 0x%06" PRIX32 " \n",
      path_switch_req_p->mme_ue_s1ap_id);
    OAILOG_FUNC_OUT(LOG_MME_APP);
  }
  if (ue_context_p->enb_s1ap_id_key != INVALID_ENB_UE_S1AP_ID_KEY) {
    /* Remove existing enb_s1ap_id_key which is mapped with suorce eNB  */
    hashtable_uint64_ts_remove(
          mme_app_desc_p->mme_ue_contexts.enb_ue_s1ap_id_ue_context_htbl,
          (const hash_key_t) ue_context_p->enb_s1ap_id_key);
    ue_context_p->enb_s1ap_id_key = INVALID_ENB_UE_S1AP_ID_KEY;
  }
  // Update MME UE context with new enb_ue_s1ap_id
  ue_context_p->enb_ue_s1ap_id = path_switch_req_p->enb_ue_s1ap_id;
  // regenerate the enb_s1ap_id_key as enb_ue_s1ap_id is changed.
  MME_APP_ENB_S1AP_ID_KEY(enb_s1ap_id_key,
        path_switch_req_p->enb_id, path_switch_req_p->enb_ue_s1ap_id);
  // Update enb_s1ap_id_key in hashtable
  if (!IS_EMM_CTXT_PRESENT_GUTI(&(ue_context_p->emm_context))) {
    mme_ue_context_update_coll_keys(
      &mme_app_desc_p->mme_ue_contexts,
      ue_context_p,
      enb_s1ap_id_key,
      ue_context_p->mme_ue_s1ap_id,
      ue_context_p->emm_context._imsi64,
      ue_context_p->mme_teid_s11,
      &ue_context_p->emm_context._guti);
  }
  ue_context_p->sctp_assoc_id_key = path_switch_req_p->sctp_assoc_id;
  ue_context_p->e_utran_cgi = path_switch_req_p->ecgi;
  ue_network_capability.eea =
    path_switch_req_p->encryption_algorithm_capabilities;
  ue_network_capability.eia =
    path_switch_req_p->integrity_algorithm_capabilities;
  if ((ue_network_capability.eea != ue_context_p->emm_context
    ._ue_network_capability.eea) || (ue_network_capability.eia !=
    ue_context_p->emm_context._ue_network_capability.eia)) {
    /* clear ue security capabilities and store security capabilities
     * recieved in PATH_SWITCH REQUEST */
    emm_ctx_clear_ue_nw_cap(&ue_context_p->emm_context);
    emm_ctx_set_valid_ue_nw_cap(&ue_context_p->emm_context,
      &ue_network_capability);
  }
  /*Build and send Modify Bearer Request*/
  message_p = itti_alloc_new_message(TASK_MME_APP, S11_MODIFY_BEARER_REQUEST);
  if (message_p == NULL) {
    OAILOG_ERROR(
      LOG_MME_APP,
      "Failed to allocate new ITTI message for S11 Modify Bearer Request "
      "for MME UE S1AP Id: " MME_UE_S1AP_ID_FMT "\n",
      path_switch_req_p->mme_ue_s1ap_id);
    OAILOG_FUNC_OUT(LOG_MME_APP);
  }
  itti_s11_modify_bearer_request_t *s11_modify_bearer_request =
    &message_p->ittiMsg.s11_modify_bearer_request;
  s11_modify_bearer_request->local_teid = ue_context_p->mme_teid_s11;

  for (idx = 0; idx < e_rab_to_be_switched_dl_list.no_of_items; idx++) {
    bearer_id = e_rab_to_be_switched_dl_list.item[idx].e_rab_id;
    if ((current_bearer_p =
      mme_app_get_bearer_context(ue_context_p, bearer_id)) == NULL) {
      OAILOG_ERROR(
        LOG_MME_APP,
        "Bearer Contex for bearer_id %d does not exist for ue_id %d\n",
        bearer_id, ue_context_p->mme_ue_s1ap_id);
    } else {
      s11_modify_bearer_request->bearer_contexts_to_be_modified
      .bearer_contexts[idx].eps_bearer_id = e_rab_to_be_switched_dl_list
      .item[idx].e_rab_id;
      s11_modify_bearer_request->bearer_contexts_to_be_modified
      .bearer_contexts[idx].s1_eNB_fteid.teid = e_rab_to_be_switched_dl_list
      .item[idx].gtp_teid;
      s11_modify_bearer_request->bearer_contexts_to_be_modified
      .bearer_contexts[idx].s1_eNB_fteid.interface_type = S1_U_ENODEB_GTP_U;
      if (4 == blength(e_rab_to_be_switched_dl_list.item[idx]
        .transport_layer_address)) {
        s11_modify_bearer_request->bearer_contexts_to_be_modified
        .bearer_contexts[idx].s1_eNB_fteid.ipv4 = 1;
        memcpy(&s11_modify_bearer_request->bearer_contexts_to_be_modified
          .bearer_contexts[idx].s1_eNB_fteid.ipv4_address,
          e_rab_to_be_switched_dl_list.item[idx].transport_layer_address
          ->data, blength(e_rab_to_be_switched_dl_list.item[idx]
          .transport_layer_address));
      } else if (16 == blength(e_rab_to_be_switched_dl_list.item[idx]
               .transport_layer_address)) {
        s11_modify_bearer_request->bearer_contexts_to_be_modified
        .bearer_contexts[idx].s1_eNB_fteid.ipv6 = 1;
        memcpy(&s11_modify_bearer_request->bearer_contexts_to_be_modified
          .bearer_contexts[idx].s1_eNB_fteid.ipv6_address,
          e_rab_to_be_switched_dl_list.item[idx].transport_layer_address
          ->data, blength(e_rab_to_be_switched_dl_list.item[idx]
          .transport_layer_address));
      } else {
        OAILOG_ERROR(
          LOG_MME_APP,
          "Invalid IP address of %d bytes found for MME UE S1AP "
          "Id: " MME_UE_S1AP_ID_FMT " (4 or 16 bytes was expected)\n",
          blength(
            e_rab_to_be_switched_dl_list.item[idx].transport_layer_address),
          path_switch_req_p->mme_ue_s1ap_id);
        OAILOG_FUNC_OUT(LOG_MME_APP);
      }
      s11_modify_bearer_request->bearer_contexts_to_be_modified
        .num_bearer_context++;

      OAILOG_DEBUG(
        LOG_MME_APP,
        "Build MBR for ue_id %d\t bearer_id %d\t enb_teid %u\t sgw_teid %u\n",
        ue_context_p->mme_ue_s1ap_id, bearer_id, s11_modify_bearer_request
        ->bearer_contexts_to_be_modified.bearer_contexts[idx].s1_eNB_fteid.teid, current_bearer_p->s_gw_fteid_s1u.teid);
    }

    if (!idx) {
      cid = ue_context_p->bearer_contexts[EBI_TO_INDEX(bearer_id)]->pdn_cx_id;
      pdn_context = ue_context_p->pdn_contexts[cid];
      s11_modify_bearer_request->edns_peer_ip.addr_v4_.sin_addr =
        pdn_context->s_gw_address_s11_s4.address.ipv4_address;
      s11_modify_bearer_request->teid = pdn_context->s_gw_teid_s11_s4;
    }
  }
  if (pdn_context->esm_data.n_bearers == e_rab_to_be_switched_dl_list
    .no_of_items) {
    s11_modify_bearer_request->bearer_contexts_to_be_removed
      .num_bearer_context = 0;
  } else {
    /* find the bearer which are present in current UE context and not present
     * in Path Switch Request, add them to bearer_contexts_to_be_removed list
     * */
    for (idx = 0; idx < pdn_context->esm_data.n_bearers; idx++) {
      bearer_id = ue_context_p->bearer_contexts[idx]->ebi;
      if (is_e_rab_id_present(
        e_rab_to_be_switched_dl_list, bearer_id) == true) {
        continue;
      } else {
        s11_modify_bearer_request->bearer_contexts_to_be_removed
        .bearer_contexts[idx].eps_bearer_id = bearer_id;
        s11_modify_bearer_request->bearer_contexts_to_be_removed
        .bearer_contexts[idx].s4u_sgsn_fteid.teid =
        ue_context_p->bearer_contexts[idx]->enb_fteid_s1u.teid;
        s11_modify_bearer_request->bearer_contexts_to_be_removed
        .bearer_contexts[idx].s4u_sgsn_fteid.interface_type =
        ue_context_p->bearer_contexts[idx]->enb_fteid_s1u.interface_type;
        if (ue_context_p->bearer_contexts[idx]->enb_fteid_s1u.ipv4) {
          s11_modify_bearer_request->bearer_contexts_to_be_removed
          .bearer_contexts[idx].s4u_sgsn_fteid.ipv4 = 1;
          s11_modify_bearer_request->bearer_contexts_to_be_removed
          .bearer_contexts[idx].s4u_sgsn_fteid.ipv4_address =
          ue_context_p->bearer_contexts[idx]->enb_fteid_s1u.ipv4_address;
        } else if (ue_context_p->bearer_contexts[idx]->enb_fteid_s1u.ipv6) {
          s11_modify_bearer_request->bearer_contexts_to_be_removed
          .bearer_contexts[idx].s4u_sgsn_fteid.ipv6 = 1;
          s11_modify_bearer_request->bearer_contexts_to_be_removed
          .bearer_contexts[idx].s4u_sgsn_fteid.ipv6_address =
          ue_context_p->bearer_contexts[idx]->enb_fteid_s1u.ipv6_address;
        }
        s11_modify_bearer_request->bearer_contexts_to_be_removed
          .num_bearer_context++;
      }
    }
  }
  /*
   * S11 stack specific parameter. Not used in standalone epc mode
   */
  s11_modify_bearer_request->trxn = NULL;

  message_p->ittiMsgHeader.imsi = ue_context_p->emm_context._imsi64;

  OAILOG_DEBUG(
    LOG_MME_APP,
    "MME_APP send S11_MODIFY_BEARER_REQUEST to teid %u \n",
    s11_modify_bearer_request->teid);
  itti_send_msg_to_task(TASK_SPGW, INSTANCE_DEFAULT, message_p);
  ue_context_p->path_switch_req = true;

  OAILOG_FUNC_OUT(LOG_MME_APP);
}

//------------------------------------------------------------------------------
void mme_app_handle_erab_rel_cmd(
  const mme_ue_s1ap_id_t ue_id,
  const ebi_t ebi,
  bstring nas_msg)
{
  OAILOG_FUNC_IN(LOG_MME_APP);
  MessageDef* message_p = NULL;
  struct ue_mm_context_s* ue_context_p = NULL;

  ue_context_p = mme_ue_context_exists_mme_ue_s1ap_id(ue_id);

  if (!ue_context_p) {
    OAILOG_ERROR(
      LOG_MME_APP,
      "UE context doesn't exist for ue_id " MME_UE_S1AP_ID_FMT "\n",
      ue_id);
    bdestroy_wrapper(&nas_msg);
    OAILOG_FUNC_OUT(LOG_MME_APP);
  }

  bearer_context_t* bearer_context =
    mme_app_get_bearer_context(ue_context_p, ebi);
  if (!bearer_context) {
    OAILOG_ERROR(
      LOG_MME_APP,
      "No bearer context found ue_id " MME_UE_S1AP_ID_FMT " ebi %u\n",
      ue_id,
      ebi);
    bdestroy_wrapper(&nas_msg);
    OAILOG_FUNC_OUT(LOG_MME_APP);
  }

  message_p = itti_alloc_new_message(TASK_MME_APP, S1AP_E_RAB_REL_CMD);
  if (message_p == NULL) {
    OAILOG_ERROR(LOG_MME_APP, "Cannot allocte memory to S1AP_E_RAB_REL_CMD \n");
    bdestroy_wrapper(&nas_msg);
    OAILOG_FUNC_OUT(LOG_MME_APP);
  }
  itti_s1ap_e_rab_rel_cmd_t* s1ap_e_rab_rel_cmd =
    &message_p->ittiMsg.s1ap_e_rab_rel_cmd;

  s1ap_e_rab_rel_cmd->mme_ue_s1ap_id = ue_context_p->mme_ue_s1ap_id;
  s1ap_e_rab_rel_cmd->enb_ue_s1ap_id = ue_context_p->enb_ue_s1ap_id;

  if (ue_context_p->emm_context.esm_ctx.is_pdn_disconnect) {
    pdn_cid_t cid = ue_context_p->bearer_contexts[EBI_TO_INDEX(ebi)]->pdn_cx_id;
    pdn_context_t* pdn_context_p = ue_context_p->pdn_contexts[cid];

    // Fill bearers_to_be_rel to be sent in ERAB_REL_CMD
    s1ap_e_rab_rel_cmd->e_rab_to_be_rel_list.no_of_items =
      pdn_context_p->esm_data.n_bearers;
    uint8_t rel_index = 0;
    for (uint8_t idx = 0;
         ((idx < BEARERS_PER_UE) &&
          (rel_index < pdn_context_p->esm_data.n_bearers));
         idx++) {
      int8_t bearer_index = pdn_context_p->bearer_contexts[idx];
      if (
        (bearer_index != INVALID_BEARER_INDEX) &&
        (ue_context_p->bearer_contexts[bearer_index])) {
        s1ap_e_rab_rel_cmd->e_rab_to_be_rel_list.item[rel_index].e_rab_id =
          ue_context_p->bearer_contexts[bearer_index]->ebi;
        rel_index++;
      }
    }
  } else {
    s1ap_e_rab_rel_cmd->e_rab_to_be_rel_list.no_of_items = 1;
    s1ap_e_rab_rel_cmd->e_rab_to_be_rel_list.item[0].e_rab_id =
      bearer_context->ebi;
  }
  /* TODO Pruthvi, To fill cause for all bearers that are to be released
   * s1ap_e_rab_rel_cmd->e_rab_to_be_rel_list.item[0].cause = 0;
   */
  s1ap_e_rab_rel_cmd->nas_pdu = nas_msg;

  OAILOG_INFO(
    LOG_MME_APP,
    "Sending ERAB REL CMD to S1AP with ue_id: " MME_UE_S1AP_ID_FMT
    "and EBI %u \n",
    ue_id,
    ebi);

  message_p->ittiMsgHeader.imsi = ue_context_p->emm_context._imsi64;
  itti_send_msg_to_task(TASK_S1AP, INSTANCE_DEFAULT, message_p);

  OAILOG_FUNC_OUT(LOG_MME_APP);
}

//------------------------------------------------------------------------------
void mme_app_handle_e_rab_rel_rsp(
  itti_s1ap_e_rab_rel_rsp_t* const e_rab_rel_rsp)
{
  OAILOG_FUNC_IN(LOG_MME_APP);

  for (int i = 0; i < e_rab_rel_rsp->e_rab_rel_list.no_of_items; i++) {
    e_rab_id_t e_rab_id = e_rab_rel_rsp->e_rab_rel_list.item[i].e_rab_id;
    OAILOG_DEBUG(
      LOG_MME_APP,
      "ERAB released successfully at UE with ERAB-ID:%u for "
      "ue_id" MME_UE_S1AP_ID_FMT "\n",
      e_rab_id,
      e_rab_rel_rsp->mme_ue_s1ap_id);
  }

  for (int i = 0;
    i < e_rab_rel_rsp->e_rab_failed_to_rel_list.no_of_items; i++) {
    e_rab_id_t e_rab_id =
      e_rab_rel_rsp->e_rab_failed_to_rel_list.item[i].e_rab_id;
    OAILOG_DEBUG(
      LOG_MME_APP,
      "Failed to release ERAB with ERAB ID %u at UE for "
      "ue_id" MME_UE_S1AP_ID_FMT "\n",
      e_rab_id,
      e_rab_rel_rsp->mme_ue_s1ap_id);
  }
  OAILOG_FUNC_OUT(LOG_MME_APP);
}

//------------------------------------------------------------------------------
bool is_e_rab_id_present(
  e_rab_to_be_switched_in_downlink_list_t e_rab_to_be_switched_dl_list,
  ebi_t bearer_id)
{
  OAILOG_FUNC_IN(LOG_MME_APP);
  uint8_t idx = 0;
  bool rc = false;

  for (; idx < e_rab_to_be_switched_dl_list.no_of_items; ++idx) {
    if (bearer_id != e_rab_to_be_switched_dl_list.item[idx].e_rab_id) {
      continue;
    } else {
      rc = true;
      break;
    }
  }
  OAILOG_FUNC_RETURN(LOG_MME_APP, rc);
}
//------------------------------------------------------------------------------
void mme_app_handle_path_switch_req_ack(
  itti_s11_modify_bearer_response_t  *const s11_modify_bearer_response,
  ue_mm_context_t *ue_context_p)
{
  OAILOG_FUNC_IN(LOG_MME_APP);
  emm_context_t *emm_ctx = &ue_context_p->emm_context;
  MessageDef *message_p = NULL;

  if (s11_modify_bearer_response->bearer_contexts_modified.num_bearer_context
    == 0) {
    mme_app_handle_path_switch_req_failure(ue_context_p);
    OAILOG_FUNC_OUT(LOG_MME_APP);
  }
  OAILOG_DEBUG(LOG_MME_APP, "Build PATH_SWITCH_REQUEST_ACK for ue_id %d\n",
    ue_context_p->mme_ue_s1ap_id);
  message_p = itti_alloc_new_message(
    TASK_MME_APP, S1AP_PATH_SWITCH_REQUEST_ACK);
  if (message_p == NULL) {
    OAILOG_ERROR(
      LOG_MME_APP,
      "Failed to allocate new ITTI message for S1AP Path Switch Request Ack "
      "for MME UE S1AP Id: " MME_UE_S1AP_ID_FMT "\n",
      ue_context_p->mme_ue_s1ap_id);
    OAILOG_FUNC_OUT(LOG_MME_APP);
  }
  itti_s1ap_path_switch_request_ack_t *s1ap_path_switch_req_ack =
    &message_p->ittiMsg.s1ap_path_switch_request_ack;

  s1ap_path_switch_req_ack->sctp_assoc_id = ue_context_p->sctp_assoc_id_key;
  s1ap_path_switch_req_ack->enb_ue_s1ap_id = ue_context_p->enb_ue_s1ap_id;
  s1ap_path_switch_req_ack->mme_ue_s1ap_id = ue_context_p->mme_ue_s1ap_id;
  memcpy(s1ap_path_switch_req_ack->NH, emm_ctx->_security.next_hop,
    AUTH_NEXT_HOP_SIZE);
  s1ap_path_switch_req_ack->NCC = emm_ctx->_security.next_hop_chaining_count;
  /* Generate NH key parameter */
  if (emm_ctx->_security.vector_index != 0) {
    OAILOG_DEBUG(
    LOG_MME_APP,
    "Invalid Vector index %d for ue_id %d \n",
    emm_ctx->_security.vector_index, ue_context_p->mme_ue_s1ap_id);
  }
  derive_NH(emm_ctx->_vector[emm_ctx->_security.vector_index].kasme,
    emm_ctx->_security.next_hop,
    emm_ctx->_security.next_hop,
    &emm_ctx->_security.next_hop_chaining_count);

  OAILOG_DEBUG(
    LOG_MME_APP,
    "MME_APP send PATH_SWITCH_REQUEST_ACK to S1AP for ue_id %d \n",
    ue_context_p->mme_ue_s1ap_id);

  message_p->ittiMsgHeader.imsi = ue_context_p->emm_context._imsi64;
  itti_send_msg_to_task(TASK_S1AP, INSTANCE_DEFAULT, message_p);

  OAILOG_FUNC_OUT(LOG_MME_APP);
}
//------------------------------------------------------------------------------
void mme_app_handle_path_switch_req_failure(
  ue_mm_context_t *ue_context_p)
{
  OAILOG_FUNC_IN(LOG_MME_APP);
  MessageDef *message_p = NULL;

  OAILOG_DEBUG(LOG_MME_APP, "Build PATH_SWITCH_REQUEST_FAILURE for ue_id %d\n",
    ue_context_p->mme_ue_s1ap_id);
  message_p = itti_alloc_new_message(
    TASK_MME_APP, S1AP_PATH_SWITCH_REQUEST_FAILURE);
  if (message_p == NULL) {
    OAILOG_ERROR(
      LOG_MME_APP,
      "Failed to allocate new ITTI message for S1AP Path Switch Request "
      "Failure for MME UE S1AP Id: " MME_UE_S1AP_ID_FMT "\n",
      ue_context_p->mme_ue_s1ap_id);
    OAILOG_FUNC_OUT(LOG_MME_APP);
  }
  itti_s1ap_path_switch_request_failure_t *s1ap_path_switch_req_failure =
    &message_p->ittiMsg.s1ap_path_switch_request_failure;

  s1ap_path_switch_req_failure->sctp_assoc_id = ue_context_p->sctp_assoc_id_key;
  s1ap_path_switch_req_failure->enb_ue_s1ap_id = ue_context_p->enb_ue_s1ap_id;
  s1ap_path_switch_req_failure->mme_ue_s1ap_id = ue_context_p->mme_ue_s1ap_id;

  OAILOG_DEBUG(
    LOG_MME_APP,
    "MME_APP send PATH_SWITCH_REQUEST_FAILURE to S1AP for ue_id %d \n",
    ue_context_p->mme_ue_s1ap_id);
  message_p->ittiMsgHeader.imsi = ue_context_p->emm_context._imsi64;
  itti_send_msg_to_task(TASK_S1AP, INSTANCE_DEFAULT, message_p);

  OAILOG_FUNC_OUT(LOG_MME_APP);
}

void mme_app_update_paging_tai_list(
  paging_tai_list_t* p_tai_list,
  partial_tai_list_t* tai_list,
  uint8_t num_of_tac)
{
  OAILOG_FUNC_IN(LOG_MME_APP);
  OAILOG_DEBUG(LOG_MME_APP, "Updating TAI list\n");

  p_tai_list->numoftac = num_of_tac;
  switch (tai_list->typeoflist) {
    case TRACKING_AREA_IDENTITY_LIST_ONE_PLMN_NON_CONSECUTIVE_TACS:
      for (int idx = 0; idx < (num_of_tac + 1); idx++) {
        p_tai_list->tai_list[idx].mcc_digit1 =
          tai_list->u.tai_one_plmn_non_consecutive_tacs.mcc_digit1;
        p_tai_list->tai_list[idx].mcc_digit2 =
          tai_list->u.tai_one_plmn_non_consecutive_tacs.mcc_digit2;
        p_tai_list->tai_list[idx].mcc_digit3 =
          tai_list->u.tai_one_plmn_non_consecutive_tacs.mcc_digit3;
        p_tai_list->tai_list[idx].mnc_digit1 =
          tai_list->u.tai_one_plmn_non_consecutive_tacs.mnc_digit1;
        p_tai_list->tai_list[idx].mnc_digit2 =
          tai_list->u.tai_one_plmn_non_consecutive_tacs.mnc_digit2;
        p_tai_list->tai_list[idx].mnc_digit3 =
          tai_list->u.tai_one_plmn_non_consecutive_tacs.mnc_digit3;
        p_tai_list->tai_list[idx].tac =
          tai_list->u.tai_one_plmn_non_consecutive_tacs.tac[idx];
      }
      break;

    case TRACKING_AREA_IDENTITY_LIST_ONE_PLMN_CONSECUTIVE_TACS:
      for (int idx = 0; idx < (num_of_tac + 1); idx++) {
        p_tai_list->tai_list[idx].mcc_digit1 =
          tai_list->u.tai_one_plmn_consecutive_tacs.mcc_digit1;
        p_tai_list->tai_list[idx].mcc_digit2 =
          tai_list->u.tai_one_plmn_consecutive_tacs.mcc_digit2;
        p_tai_list->tai_list[idx].mcc_digit3 =
          tai_list->u.tai_one_plmn_consecutive_tacs.mcc_digit3;
        p_tai_list->tai_list[idx].mnc_digit1 =
          tai_list->u.tai_one_plmn_consecutive_tacs.mnc_digit1;
        p_tai_list->tai_list[idx].mnc_digit2 =
          tai_list->u.tai_one_plmn_consecutive_tacs.mnc_digit2;
        p_tai_list->tai_list[idx].mnc_digit3 =
          tai_list->u.tai_one_plmn_consecutive_tacs.mnc_digit3;

        p_tai_list->tai_list[idx].tac =
          tai_list->u.tai_one_plmn_consecutive_tacs.tac + idx;
      }
      break;

    case TRACKING_AREA_IDENTITY_LIST_MANY_PLMNS:
      for (int idx = 0; idx < (num_of_tac + 1); idx++) {
        p_tai_list->tai_list[idx].mcc_digit1 =
          tai_list->u.tai_many_plmn[idx].mcc_digit1;
        p_tai_list->tai_list[idx].mcc_digit2 =
          tai_list->u.tai_many_plmn[idx].mcc_digit2;
        p_tai_list->tai_list[idx].mcc_digit3 =
          tai_list->u.tai_many_plmn[idx].mcc_digit3;
        p_tai_list->tai_list[idx].mnc_digit1 =
          tai_list->u.tai_many_plmn[idx].mnc_digit1;
        p_tai_list->tai_list[idx].mnc_digit2 =
          tai_list->u.tai_many_plmn[idx].mnc_digit2;
        p_tai_list->tai_list[idx].mnc_digit3 =
          tai_list->u.tai_many_plmn[idx].mnc_digit3;

        p_tai_list->tai_list[idx].tac = tai_list->u.tai_many_plmn[idx].tac;
      }
      break;

    default:
      OAILOG_ERROR(
        LOG_MME_APP,
        "BAD TAI list configuration, unknown TAI list type %u",
        tai_list->typeoflist);
      break;
  }
  OAILOG_FUNC_OUT(LOG_MME_APP);
}

// Fetch UE context based on mme_ue_s1ap_id and return pointer to UE context
ue_mm_context_t* mme_app_get_ue_context_for_timer(
  mme_ue_s1ap_id_t mme_ue_s1ap_id, char* timer_name)
{
  OAILOG_FUNC_IN(LOG_MME_APP);
  OAILOG_INFO(LOG_MME_APP, "Expired- %s for ue_id " MME_UE_S1AP_ID_FMT
    "\n",
    timer_name,
    mme_ue_s1ap_id);

  ue_mm_context_t* ue_context_p =
    mme_ue_context_exists_mme_ue_s1ap_id(mme_ue_s1ap_id);
  if (ue_context_p == NULL) {
    OAILOG_ERROR(
      LOG_MME_APP,
      "Failed to get ue context while handling %s for ue_id " MME_UE_S1AP_ID_FMT
      "\n",
      timer_name,
      mme_ue_s1ap_id);
    return NULL;
  }
  return ue_context_p;
}

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

/*! \file mme_app_itti_messaging.c
  \brief
  \author Sebastien ROUX, Lionel Gauthier
  \company Eurecom
  \email: lionel.gauthier@eurecom.fr
*/

#include <stdio.h>
#include <string.h>
#include <stdint.h>
#include <netinet/in.h>

#include "bstrlib.h"
#include "log.h"
#include "conversions.h"
#include "common_types.h"
#include "intertask_interface.h"
#include "gcc_diag.h"
#include "mme_config.h"
#include "mme_app_ue_context.h"
#include "mme_app_apn_selection.h"
#include "mme_app_bearer_context.h"
#include "sgw_ie_defs.h"
#include "common_defs.h"
#include "mme_app_itti_messaging.h"
#include "mme_app_sgw_selection.h"
#include "3gpp_23.003.h"
#include "3gpp_24.008.h"
#include "3gpp_29.274.h"
#include "emm_data.h"
#include "esm_data.h"
#include "mme_app_desc.h"
#include "s11_messages_types.h"
#include "common_utility_funs.h"

#if EMBEDDED_SGW
#define TASK_SPGW TASK_SPGW_APP
#else
#define TASK_SPGW TASK_S11
#endif

extern task_zmq_ctx_t mme_app_task_zmq_ctx;

/****************************************************************************
 **                                                                        **
 ** name:    mme_app_itti_ue_context_release()                             **
 **                                                                        **
 ** description: Send itti mesage to S1ap task to send UE Context Release  **
 **              Request                                                   **
 **                                                                        **
 ** inputs:  ue_context_p: Pointer to UE context                           **
 **          emm_casue: failed cause                                       **
 **                                                                        **
 ***************************************************************************/
void mme_app_itti_ue_context_release(
    struct ue_mm_context_s* ue_context_p, enum s1cause cause) {
  MessageDef* message_p;

  OAILOG_FUNC_IN(LOG_MME_APP);
  message_p =
      itti_alloc_new_message(TASK_MME_APP, S1AP_UE_CONTEXT_RELEASE_COMMAND);
  if (message_p == NULL) {
    OAILOG_ERROR_UE(
        LOG_MME_APP, ue_context_p->emm_context._imsi64,
        "Failed to allocate memory for S1AP_UE_CONTEXT_RELEASE_COMMAND \n");
    OAILOG_FUNC_OUT(LOG_MME_APP);
  }

  OAILOG_INFO_UE(
      LOG_MME_APP, ue_context_p->emm_context._imsi64,
      "Sending UE Context Release Cmd to S1ap for (ue_id = " MME_UE_S1AP_ID_FMT
      ")\n"
      "UE Context Release Cause = (%d)\n",
      ue_context_p->mme_ue_s1ap_id, cause);

  S1AP_UE_CONTEXT_RELEASE_COMMAND(message_p).mme_ue_s1ap_id =
      ue_context_p->mme_ue_s1ap_id;
  S1AP_UE_CONTEXT_RELEASE_COMMAND(message_p).enb_ue_s1ap_id =
      ue_context_p->enb_ue_s1ap_id;
  S1AP_UE_CONTEXT_RELEASE_COMMAND(message_p).cause = cause;

  message_p->ittiMsgHeader.imsi = ue_context_p->emm_context._imsi64;
  send_msg_to_task(&mme_app_task_zmq_ctx, TASK_S1AP, message_p);
  OAILOG_FUNC_OUT(LOG_MME_APP);
}

/****************************************************************************
 **                                                                        **
 ** name:    mme_app_send_s11_release_access_bearers_req                   **
 **                                                                        **
 ** description: Send itti mesage to SPGW task to send Release Access      **
 **             Bearer Request (RAB)                                       **
 **                                                                        **
 ** inputs:  ue_context_p: Pointer to UE context                           **
 **          pdn_index: PDN index for which RAB is initiated               **
 **                                                                        **
 ** outputs:                                                               **
 **      Return:    RETURNok, RETURNerror                                  **
 **                                                                        **
 ***************************************************************************/
int mme_app_send_s11_release_access_bearers_req(
    struct ue_mm_context_s* const ue_mm_context, const pdn_cid_t pdn_index) {
  OAILOG_FUNC_IN(LOG_MME_APP);
  /*
   * Keep the identifier to the default APN
   */
  MessageDef* message_p = NULL;
  itti_s11_release_access_bearers_request_t* release_access_bearers_request_p =
      NULL;

  if (ue_mm_context == NULL) {
    OAILOG_ERROR(LOG_MME_APP, "Invalid UE MM context received\n");
    OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNerror);
  }
  message_p =
      itti_alloc_new_message(TASK_MME_APP, S11_RELEASE_ACCESS_BEARERS_REQUEST);
  if (message_p == NULL) {
    OAILOG_ERROR_UE(
        LOG_MME_APP, ue_mm_context->emm_context._imsi64,
        "Failed to allocate memory for S11_RELEASE_ACCESS_BEARERS_REQUEST  for "
        "ue id " MME_UE_S1AP_ID_FMT "\n",
        ue_mm_context->mme_ue_s1ap_id);
    OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNerror);
  }
  release_access_bearers_request_p =
      &message_p->ittiMsg.s11_release_access_bearers_request;
  release_access_bearers_request_p->local_teid = ue_mm_context->mme_teid_s11;
  pdn_context_t* pdn_connection = ue_mm_context->pdn_contexts[pdn_index];
  release_access_bearers_request_p->teid = pdn_connection->s_gw_teid_s11_s4;
  release_access_bearers_request_p->edns_peer_ip.addr_v4.sin_addr =
      pdn_connection->s_gw_address_s11_s4.address.ipv4_address;
  release_access_bearers_request_p->edns_peer_ip.addr_v4.sin_family = AF_INET;
  release_access_bearers_request_p->originating_node = NODE_TYPE_MME;

  message_p->ittiMsgHeader.imsi = ue_mm_context->emm_context._imsi64;
  if (pdn_connection->route_s11_messages_to_s8_task) {
    OAILOG_INFO_UE(
        LOG_MME_APP, ue_mm_context->emm_context._imsi64,
        "Send Release Access Bearer Req for teid to sgw_s8 task " TEID_FMT
        " for ue id " MME_UE_S1AP_ID_FMT "\n",
        ue_mm_context->mme_teid_s11, ue_mm_context->mme_ue_s1ap_id);
    send_msg_to_task(&mme_app_task_zmq_ctx, TASK_SGW_S8, message_p);
  } else {
    OAILOG_INFO_UE(
        LOG_MME_APP, ue_mm_context->emm_context._imsi64,
        "Send Release Access Bearer Req for teid to spgw task " TEID_FMT
        " for ue id " MME_UE_S1AP_ID_FMT "\n",
        ue_mm_context->mme_teid_s11, ue_mm_context->mme_ue_s1ap_id);
    send_msg_to_task(&mme_app_task_zmq_ctx, TASK_SPGW, message_p);
  }
  OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNok);
}

/****************************************************************************
 **                                                                        **
 ** name:    mme_app_send_s11_create_session_req                           **
 **                                                                        **
 ** description: Send itti mesage to SPGW task to send Create Session      **
 **              Request (CSR)                                             **
 **                                                                        **
 ** inputs:  mme_app_desc_p: Pointer to structure, mme_app_desc_t          **
 **          ue_mm_context: Pointer to ue_mm_context_s                     **
 **          pdn_index: PDN index for which CSR is initiated               **
 **                                                                        **
 ** outputs:                                                               **
 **      Return:    RETURNok, RETURNerror                                  **
 **                                                                        **
 ***************************************************************************/
int mme_app_send_s11_create_session_req(
    mme_app_desc_t* mme_app_desc_p, struct ue_mm_context_s* const ue_mm_context,
    const pdn_cid_t pdn_cid) {
  OAILOG_FUNC_IN(LOG_MME_APP);

  /*
   * Keep the identifier to the default APN
   */
  MessageDef* message_p                                = NULL;
  itti_s11_create_session_request_t* session_request_p = NULL;

  if (ue_mm_context == NULL) {
    OAILOG_ERROR(LOG_MME_APP, "Invalid UE MM context received\n");
    OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNerror);
  }
  if (ue_mm_context->subscriber_status != SS_SERVICE_GRANTED) {
    /*
     * HSS rejected the bearer creation or roaming is not allowed for this
     * UE. This result will trigger an ESM Failure message sent to UE.
     */
    OAILOG_ERROR_UE(
        LOG_MME_APP, ue_mm_context->emm_context._imsi64,
        "Not implemented: ACCESS NOT GRANTED, send ESM Failure to NAS\n");
    OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNerror);
  }

  message_p = itti_alloc_new_message(TASK_MME_APP, S11_CREATE_SESSION_REQUEST);
  /*
   * WARNING:
   * Some parameters should be provided by NAS Layer:
   * - ue_time_zone
   * - mei
   * - uli
   * - uci
   * Some parameters should be provided by HSS:
   * - PGW address for CP
   * - paa
   * - ambr
   * - charging characteristics
   * and by MME Application layer:
   * - selection_mode
   * Set these parameters with random values for now.
   */
  session_request_p = &message_p->ittiMsg.s11_create_session_request;
  /*
   * As the create session request is the first exchanged message and as
   * no tunnel had been previously setup, the distant teid is set to 0.
   * The remote teid will be provided in the response message.
   */
  session_request_p->teid = 0;
  IMSI64_TO_STRING(
      ue_mm_context->emm_context._imsi64,
      (char*) (&session_request_p->imsi.digit),
      ue_mm_context->emm_context._imsi.length);
  session_request_p->imsi.length = ue_mm_context->emm_context._imsi.length;

  message_p->ittiMsgHeader.imsi = ue_mm_context->emm_context._imsi64;

  /*
   * Copy the MSISDN
   */
  if (ue_mm_context->msisdn) {
    memcpy(
        session_request_p->msisdn.digit, ue_mm_context->msisdn->data,
        ue_mm_context->msisdn->slen);
    session_request_p->msisdn.length = ue_mm_context->msisdn->slen;
  } else {
    session_request_p->msisdn.length = 0;
  }
  session_request_p->mei.present       = MEI_IMEISV;
  session_request_p->mei.choice.imeisv = ue_mm_context->emm_context._imeisv;
  // Fill User Location Information
  session_request_p->uli.present = 0;  // initialize the presencemask
  mme_app_get_user_location_information(&session_request_p->uli, ue_mm_context);

  session_request_p->rat_type = RAT_EUTRAN;

  // default bearer already created by NAS
  bearer_context_t* bc = mme_app_get_bearer_context(
      ue_mm_context, ue_mm_context->pdn_contexts[pdn_cid]->default_ebi);
  if (bc == NULL) {
    OAILOG_ERROR_UE(
        LOG_MME_APP, ue_mm_context->emm_context._imsi64,
        "Failed to send create session req to SPGW as the bearer context is "
        "NULL for "
        "MME UE S1AP Id: " MME_UE_S1AP_ID_FMT " for bearer %u\n",
        ue_mm_context->mme_ue_s1ap_id,
        ue_mm_context->pdn_contexts[pdn_cid]->default_ebi);
    OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNerror);
  }
  // Zero because default bearer (see 29.274)
  session_request_p->bearer_contexts_to_be_created.bearer_contexts[0]
      .bearer_level_qos.gbr.br_ul = bc->esm_ebr_context.gbr_ul;
  session_request_p->bearer_contexts_to_be_created.bearer_contexts[0]
      .bearer_level_qos.gbr.br_dl = bc->esm_ebr_context.gbr_dl;
  session_request_p->bearer_contexts_to_be_created.bearer_contexts[0]
      .bearer_level_qos.mbr.br_ul = bc->esm_ebr_context.mbr_ul;
  session_request_p->bearer_contexts_to_be_created.bearer_contexts[0]
      .bearer_level_qos.mbr.br_dl = bc->esm_ebr_context.mbr_dl;
  session_request_p->bearer_contexts_to_be_created.bearer_contexts[0]
      .bearer_level_qos.qci = bc->qci;
  session_request_p->bearer_contexts_to_be_created.bearer_contexts[0]
      .bearer_level_qos.pvi = bc->preemption_vulnerability;
  session_request_p->bearer_contexts_to_be_created.bearer_contexts[0]
      .bearer_level_qos.pci = bc->preemption_capability;
  session_request_p->bearer_contexts_to_be_created.bearer_contexts[0]
      .bearer_level_qos.pl = bc->priority_level;
  session_request_p->bearer_contexts_to_be_created.bearer_contexts[0]
      .eps_bearer_id                                                  = bc->ebi;
  session_request_p->bearer_contexts_to_be_created.num_bearer_context = 1;
  session_request_p->sender_fteid_for_cp.teid =
      (teid_t) ue_mm_context->mme_ue_s1ap_id;
  session_request_p->sender_fteid_for_cp.interface_type = S11_MME_GTP_C;
  mme_config_read_lock(&mme_config);
  session_request_p->sender_fteid_for_cp.ipv4_address.s_addr =
      mme_config.ip.s11_mme_v4.s_addr;
  mme_config_unlock(&mme_config);
  session_request_p->sender_fteid_for_cp.ipv4 = 1;

  // ue_mm_context->mme_teid_s11 = session_request_p->sender_fteid_for_cp.teid;
  ue_mm_context->pdn_contexts[pdn_cid]->s_gw_teid_s11_s4 = 0;
  mme_ue_context_update_coll_keys(
      &mme_app_desc_p->mme_ue_contexts, ue_mm_context,
      ue_mm_context->enb_s1ap_id_key, ue_mm_context->mme_ue_s1ap_id,
      ue_mm_context->emm_context._imsi64,
      session_request_p->sender_fteid_for_cp.teid,  // mme_teid_s11 is new
      &ue_mm_context->emm_context._guti);
  struct apn_configuration_s* selected_apn_config_p = mme_app_get_apn_config(
      ue_mm_context, ue_mm_context->pdn_contexts[pdn_cid]->context_identifier);

  memcpy(
      session_request_p->apn, selected_apn_config_p->service_selection,
      selected_apn_config_p->service_selection_length);

  /*
   * Copy the APN AMBR to the sgw create session request message
   */
  memcpy(
      &session_request_p->ambr, &selected_apn_config_p->ambr, sizeof(ambr_t));
  /*
   * Set PDN type for pdn_type and PAA even if this IE is redundant
   */
  OAILOG_DEBUG_UE(
      LOG_MME_APP, ue_mm_context->emm_context._imsi64,
      "selected apn config PDN Type = %d for (ue_id = %u)\n",
      selected_apn_config_p->pdn_type, ue_mm_context->mme_ue_s1ap_id);
  session_request_p->pdn_type     = selected_apn_config_p->pdn_type;
  session_request_p->paa.pdn_type = selected_apn_config_p->pdn_type;

  if (selected_apn_config_p->nb_ip_address == 0) {
    /*
     * UE DHCPv4 allocated ip address
     */
    session_request_p->paa.ipv4_address.s_addr = INADDR_ANY;
    session_request_p->paa.ipv6_address        = in6addr_any;
  } else {
    uint8_t j;

    for (j = 0; j < selected_apn_config_p->nb_ip_address; j++) {
      ip_address_t* ip_address = &selected_apn_config_p->ip_address[j];

      if (ip_address->pdn_type == IPv4) {
        session_request_p->paa.ipv4_address.s_addr =
            ip_address->address.ipv4_address.s_addr;
      } else if (ip_address->pdn_type == IPv6) {
        memcpy(
            &session_request_p->paa.ipv6_address,
            &ip_address->address.ipv6_address,
            sizeof(session_request_p->paa.ipv6_address));
      }
    }
  }

  // Add Charging Characteristics
  // If per-APN characteristics is specified, pass it. Otherwise, pass the
  // default value. The length values should be set to 0 if there is no value
  // specified.
  if (selected_apn_config_p->charging_characteristics.length > 0) {
    memcpy(
        &session_request_p->charging_characteristics,
        &selected_apn_config_p->charging_characteristics,
        sizeof(charging_characteristics_t));
  } else {
    memcpy(
        &session_request_p->charging_characteristics,
        &ue_mm_context->default_charging_characteristics,
        sizeof(charging_characteristics_t));
  }

  if (ue_mm_context->pdn_contexts[pdn_cid]->pco) {
    copy_protocol_configuration_options(
        &session_request_p->pco, ue_mm_context->pdn_contexts[pdn_cid]->pco);
  }

  // TODO perform SGW selection
  // Actually, since S and P GW are bundled together, there is no PGW selection
  // (based on PGW id in ULA, or DNS query based on FQDN)
  if (1) {
    mme_app_select_sgw(
        &ue_mm_context->emm_context.originating_tai,
        (struct sockaddr* const) & session_request_p->edns_peer_ip);
  }
  COPY_PLMN_IN_ARRAY_FMT(
      (session_request_p->serving_network), (ue_mm_context->e_utran_cgi.plmn));
  session_request_p->selection_mode = MS_O_N_P_APN_S_V;
  int mode =
      match_fed_mode_map((char*) session_request_p->imsi.digit, LOG_MME_APP);
  if (mode == S8_SUBSCRIBER) {
    OAILOG_INFO_UE(
        LOG_MME_APP, ue_mm_context->emm_context._imsi64,
        "Sending s11 create session req message to SGW_s8 task for "
        "ue_id " MME_UE_S1AP_ID_FMT "\n",
        ue_mm_context->mme_ue_s1ap_id);
    send_msg_to_task(&mme_app_task_zmq_ctx, TASK_SGW_S8, message_p);
    ue_mm_context->pdn_contexts[pdn_cid]->route_s11_messages_to_s8_task = true;
  } else {
    OAILOG_INFO_UE(
        LOG_MME_APP, ue_mm_context->emm_context._imsi64,
        "Sending s11 create session req message to SPGW task for "
        "ue_id " MME_UE_S1AP_ID_FMT "\n",
        ue_mm_context->mme_ue_s1ap_id);
    send_msg_to_task(&mme_app_task_zmq_ctx, TASK_SPGW, message_p);
  }
  OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNok);
}

//------------------------------------------------------------------------------
/**
 * Send an S1AP E-RAB Modification confirm to the S1AP layer.
 * More IEs will be added soon.
 */
void mme_app_send_s1ap_e_rab_modification_confirm(
    const mme_ue_s1ap_id_t mme_ue_s1ap_id,
    const enb_ue_s1ap_id_t enb_ue_s1ap_id,
    const mme_app_s1ap_proc_modify_bearer_ind_t* const proc) {
  OAILOG_FUNC_IN(LOG_MME_APP);
  /** Send a S1AP E-RAB MODIFICATION CONFIRM TO THE ENB. */
  MessageDef* message_p =
      itti_alloc_new_message(TASK_MME_APP, S1AP_E_RAB_MODIFICATION_CNF);
  if (message_p == NULL) {
    OAILOG_ERROR(LOG_MME_APP, "itti_alloc_new_message Failed\n");
  }

  itti_s1ap_e_rab_modification_cnf_t* s1ap_e_rab_modification_cnf_p =
      &message_p->ittiMsg.s1ap_e_rab_modification_cnf;

  /** Set the identifiers. */
  s1ap_e_rab_modification_cnf_p->mme_ue_s1ap_id = mme_ue_s1ap_id;
  s1ap_e_rab_modification_cnf_p->enb_ue_s1ap_id = enb_ue_s1ap_id;

  for (int i = 0; i < proc->e_rab_modified_list.no_of_items; ++i) {
    s1ap_e_rab_modification_cnf_p->e_rab_modify_list.e_rab_id[i] =
        proc->e_rab_modified_list.e_rab_id[i];
  }
  s1ap_e_rab_modification_cnf_p->e_rab_modify_list.no_of_items =
      proc->e_rab_modified_list.no_of_items;

  for (int i = 0; i < proc->e_rab_failed_to_be_modified_list.no_of_items; ++i) {
    s1ap_e_rab_modification_cnf_p->e_rab_failed_to_modify_list.item[i]
        .e_rab_id = proc->e_rab_failed_to_be_modified_list.item[i].e_rab_id;
    s1ap_e_rab_modification_cnf_p->e_rab_failed_to_modify_list.item[i].cause =
        proc->e_rab_failed_to_be_modified_list.item[i].cause;
  }
  s1ap_e_rab_modification_cnf_p->e_rab_failed_to_modify_list.no_of_items =
      proc->e_rab_failed_to_be_modified_list.no_of_items;

  /** Sending a message to S1AP. */
  send_msg_to_task(&mme_app_task_zmq_ctx, TASK_S1AP, message_p);
  OAILOG_FUNC_OUT(LOG_MME_APP);
}

/****************************************************************************
 **                                                                        **
 ** name:    nas_itti_sgsap_uplink_unitdata                                **
 **                                                                        **
 ** description: Send itti mesage to SGS or SMS_ORC8R task to send NAS     **
 **              message.                                                  **
 **                                                                        **
 ** inputs:  imsi : IMSI of UE                                             **
 **          imsi_len : Length of IMSI                                     **
 **          nas_msg: NAS message                                          **
 **          imeisv_pP: IMEISV of UE                                       **
 **          mobilestationclassmark2_pP: Mobile station classmark-2 of UE  **
 **          tai_pP: TAI of UE                                             **
 **          ecgi_pP: ecgi of UE                                           **
 **                                                                        **
 ***************************************************************************/
void nas_itti_sgsap_uplink_unitdata(
    const char* const imsi, uint8_t imsi_len, bstring nas_msg,
    imeisv_t* imeisv_pP, MobileStationClassmark2* mobilestationclassmark2_pP,
    tai_t* tai_pP, ecgi_t* ecgi_pP, bool sms_orc8r_enabled) {
  OAILOG_FUNC_IN(LOG_MME_APP);
  MessageDef* message_p = NULL;
  int uetimezone        = 0;

  message_p = itti_alloc_new_message(TASK_MME_APP, SGSAP_UPLINK_UNITDATA);
  if (message_p == NULL) {
    OAILOG_ERROR(
        LOG_MME_APP, "Failed to allocate memory for SGSAP_UPLINK_UNITDATA \n");
    OAILOG_FUNC_OUT(LOG_MME_APP);
  }
  memset(
      &message_p->ittiMsg.sgsap_uplink_unitdata, 0,
      sizeof(itti_sgsap_uplink_unitdata_t));
  memcpy(SGSAP_UPLINK_UNITDATA(message_p).imsi, imsi, imsi_len);
  SGSAP_UPLINK_UNITDATA(message_p).imsi[imsi_len]    = '\0';
  SGSAP_UPLINK_UNITDATA(message_p).imsi_length       = imsi_len;
  SGSAP_UPLINK_UNITDATA(message_p).nas_msg_container = nas_msg;
  nas_msg                                            = NULL;
  /*
   * optional - UE Time Zone
   * update the ue time zone presence bitmask
   */
  if ((uetimezone = get_time_zone()) != RETURNerror) {
    SGSAP_UPLINK_UNITDATA(message_p).opt_ue_time_zone = timezone;
    SGSAP_UPLINK_UNITDATA(message_p).presencemask =
        UPLINK_UNITDATA_UE_TIMEZONE_PARAMETER_PRESENT;
  }
  /*
   * optional - IMEISV
   * update the imeisv presence bitmask
   */
  if (imeisv_pP) {
    hexa_to_ascii(
        (uint8_t*) imeisv_pP->u.value,
        SGSAP_UPLINK_UNITDATA(message_p).opt_imeisv, 8);
    SGSAP_UPLINK_UNITDATA(message_p).opt_imeisv[imeisv_pP->length] = '\0';
    SGSAP_UPLINK_UNITDATA(message_p).opt_imeisv_length = imeisv_pP->length;
    SGSAP_UPLINK_UNITDATA(message_p).presencemask |=
        UPLINK_UNITDATA_IMEISV_PARAMETER_PRESENT;
  }
  /*
   * optional - mobile station classmark2
   * update the mobile station classmark2 presence bitmask.
   */
  if (mobilestationclassmark2_pP) {
    SGSAP_UPLINK_UNITDATA(message_p).opt_mobilestationclassmark2 =
        *((MobileStationClassmark2_t*) mobilestationclassmark2_pP);
    SGSAP_UPLINK_UNITDATA(message_p).presencemask |=
        UPLINK_UNITDATA_MOBILE_STATION_CLASSMARK_2_PARAMETER_PRESENT;
  }
  /*
   * optional - tai
   * update the tai presence bitmask.
   */
  if (tai_pP) {
    SGSAP_UPLINK_UNITDATA(message_p).opt_tai = *((tai_t*) tai_pP);
    SGSAP_UPLINK_UNITDATA(message_p).presencemask |=
        UPLINK_UNITDATA_TAI_PARAMETER_PRESENT;
  }
  /*
   * optional - ecgi
   * update the ecgi presence bitmask.
   */
  if (ecgi_pP) {
    SGSAP_UPLINK_UNITDATA(message_p).opt_ecgi = *ecgi_pP;
    SGSAP_UPLINK_UNITDATA(message_p).presencemask |=
        UPLINK_UNITDATA_ECGI_PARAMETER_PRESENT;
  }

  IMSI_STRING_TO_IMSI64(imsi, &message_p->ittiMsgHeader.imsi);
  imsi64_t imsi64 = message_p->ittiMsgHeader.imsi;

  // Check if we're in SMS_ORC8R and send to the appropriate task if so
  if (sms_orc8r_enabled) {
    if (send_msg_to_task(&mme_app_task_zmq_ctx, TASK_SMS_ORC8R, message_p) !=
        RETURNok) {
      OAILOG_ERROR_UE(
          LOG_MME_APP, imsi64,
          "Failed to send SGSAP Uplink Unitdata to SMS_ORC8R task\n");
    } else {
      OAILOG_DEBUG_UE(
          LOG_MME_APP, imsi64,
          "Sent SGSAP Uplink Unitdata to SMS_ORC8R task\n");
    }
  } else {
    if (send_msg_to_task(&mme_app_task_zmq_ctx, TASK_SGS, message_p) !=
        RETURNok) {
      OAILOG_ERROR_UE(
          LOG_MME_APP, imsi64,
          "Failed to send SGSAP Uplink Unitdata to SGS task\n");
    } else {
      OAILOG_DEBUG_UE(
          LOG_MME_APP, imsi64, "Sent SGSAP Uplink Unitdata to SGS task\n");
    }
  }

  OAILOG_FUNC_OUT(LOG_MME_APP);
}

/****************************************************************************
 **                                                                        **
 ** name:    mme_app_itti_sgsap_tmsi_reallocation_comp **
 **                                                                        **
 ** description: Send itti mesage, TMSI Reallocation Complete message      **
 **             to SGS task                                                **
 ** inputs:  imsi : IMSI of UE                                             **
 **          imsi_len : Length of IMSI                                     **
 **                                                                        **
 ***************************************************************************/
void mme_app_itti_sgsap_tmsi_reallocation_comp(
    const char* imsi, const unsigned int imsi_len) {
  OAILOG_FUNC_IN(LOG_MME_APP);
  MessageDef* message_p = NULL;

  message_p = itti_alloc_new_message(TASK_MME_APP, SGSAP_TMSI_REALLOC_COMP);
  if (message_p == NULL) {
    OAILOG_ERROR(
        LOG_MME_APP,
        "Failed to allocate memory for SGSAP_TMSI_REALLOC_COMP \n");
    OAILOG_FUNC_OUT(LOG_MME_APP);
  }
  memset(
      &message_p->ittiMsg.sgsap_tmsi_realloc_comp, 0,
      sizeof(itti_sgsap_tmsi_reallocation_comp_t));
  memcpy(SGSAP_TMSI_REALLOC_COMP(message_p).imsi, imsi, imsi_len);
  SGSAP_TMSI_REALLOC_COMP(message_p).imsi[imsi_len] = '\0';
  SGSAP_TMSI_REALLOC_COMP(message_p).imsi_length    = imsi_len;

  IMSI_STRING_TO_IMSI64(imsi, &message_p->ittiMsgHeader.imsi);
  imsi64_t imsi64 = message_p->ittiMsgHeader.imsi;
  if (send_msg_to_task(&mme_app_task_zmq_ctx, TASK_SGS, message_p) !=
      RETURNok) {
    OAILOG_ERROR_UE(
        LOG_MME_APP, imsi64,
        "Failed to send SGSAP Tmsi Reallocation Complete to SGS task\n");
  } else {
    OAILOG_DEBUG_UE(
        LOG_MME_APP, imsi64,
        "Sent SGSAP Tmsi Reallocation Complete to SGS task\n");
  }
  OAILOG_FUNC_OUT(LOG_MME_APP);
}

/****************************************************************************
 **                                                                        **
 ** name:    mme_app_itti_sgsap_ue_activity_ind                            **
 **                                                                        **
 ** description: Send itti mesage, UE Activity Indication message          **
 **             to SGS task                                                **
 ** inputs:  imsi : IMSI of UE                                             **
 **          imsi_len : Length of IMSI                                     **
 **                                                                        **
 ***************************************************************************/
void mme_app_itti_sgsap_ue_activity_ind(
    const char* imsi, const unsigned int imsi_len) {
  OAILOG_FUNC_IN(LOG_MME_APP);
  MessageDef* message_p = NULL;

  message_p = itti_alloc_new_message(TASK_MME_APP, SGSAP_UE_ACTIVITY_IND);
  if (message_p == NULL) {
    OAILOG_ERROR(
        LOG_MME_APP, "Failed to allocate memory for SGSAP_UE_ACTIVITY_IND \n");
    OAILOG_FUNC_OUT(LOG_MME_APP);
  }
  memset(
      &message_p->ittiMsg.sgsap_ue_activity_ind, 0,
      sizeof(itti_sgsap_ue_activity_ind_t));
  memcpy(SGSAP_UE_ACTIVITY_IND(message_p).imsi, imsi, imsi_len);
  SGSAP_UE_ACTIVITY_IND(message_p).imsi[imsi_len] = '\0';
  SGSAP_UE_ACTIVITY_IND(message_p).imsi_length    = imsi_len;

  IMSI_STRING_TO_IMSI64(imsi, &message_p->ittiMsgHeader.imsi);
  imsi64_t imsi64 = message_p->ittiMsgHeader.imsi;
  if (send_msg_to_task(&mme_app_task_zmq_ctx, TASK_SGS, message_p) !=
      RETURNok) {
    OAILOG_ERROR_UE(
        LOG_MME_APP, imsi64,
        "Failed to send SGSAP UE ACTIVITY IND to SGS task for Imsi : %s \n",
        imsi);
  } else {
    OAILOG_DEBUG_UE(
        LOG_MME_APP, imsi64,
        "Sent SGSAP UE ACTIVITY IND to SGS task for Imsi :%s \n", imsi);
  }
  OAILOG_FUNC_OUT(LOG_MME_APP);
}

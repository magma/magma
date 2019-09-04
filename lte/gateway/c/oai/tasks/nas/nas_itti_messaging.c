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

/*! \file nas_itti_messaging.c
   \brief
   \author  Sebastien ROUX, Lionel GAUTHIER
   \date
   \email: lionel.gauthier@eurecom.fr
*/

#include <ctype.h>
#include <stdio.h>
#include <string.h>
#include <stdbool.h>
#include <stdint.h>
#include <time.h>

#include "bstrlib.h"
#include "log.h"
#include "assertions.h"
#include "conversions.h"
#include "intertask_interface.h"
#include "common_defs.h"
#include "secu_defs.h"
#include "mme_app_ue_context.h"
#include "esm_proc.h"
#include "nas_itti_messaging.h"
#include "nas_proc.h"
#include "emm_proc.h"
#include "3gpp_24.008.h"
#include "3gpp_24.301.h"
#include "3gpp_29.274.h"
#include "3gpp_33.401.h"
#include "EpsAttachType.h"
#include "common_ies.h"
#include "emm_data.h"
#include "intertask_interface_types.h"
#include "itti_types.h"
#include "mme_app_desc.h"
#include "mme_app_messages_types.h"
#include "nas_messages_types.h"
#include "nas_procedures.h"
#include "nas_timer.h"
#include "s6a_messages_types.h"
#include "nas/securityDef.h"
#include "sgs_messages_types.h"

#define TASK_ORIGIN TASK_NAS_MME

//------------------------------------------------------------------------------
int nas_itti_dl_data_req(
  const mme_ue_s1ap_id_t ue_id,
  bstring nas_msg,
  nas_error_code_t transaction_status)
{
  MessageDef *message_p =
    itti_alloc_new_message(TASK_NAS_MME, NAS_DOWNLINK_DATA_REQ);
  NAS_DL_DATA_REQ(message_p).ue_id = ue_id;
  NAS_DL_DATA_REQ(message_p).nas_msg = nas_msg;
  nas_msg = NULL;
  NAS_DL_DATA_REQ(message_p).transaction_status = transaction_status;
  // make a long way by MME_APP instead of S1AP to retrieve the sctp_association_id key.
  return itti_send_msg_to_task(TASK_MME_APP, INSTANCE_DEFAULT, message_p);
}

//------------------------------------------------------------------------------
int nas_itti_erab_setup_req(
  const mme_ue_s1ap_id_t ue_id,
  const ebi_t ebi,
  const bitrate_t mbr_dl,
  const bitrate_t mbr_ul,
  const bitrate_t gbr_dl,
  const bitrate_t gbr_ul,
  bstring nas_msg)
{
  MessageDef *message_p =
    itti_alloc_new_message(TASK_NAS_MME, NAS_ERAB_SETUP_REQ);
  NAS_ERAB_SETUP_REQ(message_p).ue_id = ue_id;
  NAS_ERAB_SETUP_REQ(message_p).ebi = ebi;
  NAS_ERAB_SETUP_REQ(message_p).mbr_dl = mbr_dl;
  NAS_ERAB_SETUP_REQ(message_p).mbr_ul = mbr_ul;
  NAS_ERAB_SETUP_REQ(message_p).gbr_dl = gbr_dl;
  NAS_ERAB_SETUP_REQ(message_p).gbr_ul = gbr_ul;
  NAS_ERAB_SETUP_REQ(message_p).nas_msg = nas_msg;
  nas_msg = NULL;
  // make a long way by MME_APP instead of S1AP to retrieve the sctp_association_id key.
  return itti_send_msg_to_task(TASK_MME_APP, INSTANCE_DEFAULT, message_p);
}

//------------------------------------------------------------------------------
int nas_itti_erab_rel_cmd(
  const mme_ue_s1ap_id_t ue_id,
  const ebi_t ebi,
  bstring nas_msg)
{
  MessageDef *message_p =
    itti_alloc_new_message(TASK_NAS_MME, NAS_ERAB_REL_CMD);
  NAS_ERAB_REL_CMD(message_p).ue_id = ue_id;
  NAS_ERAB_REL_CMD(message_p).ebi = ebi;
  NAS_ERAB_REL_CMD(message_p).nas_msg = nas_msg;
  nas_msg = NULL;
  return itti_send_msg_to_task(TASK_MME_APP, INSTANCE_DEFAULT, message_p);
}


//------------------------------------------------------------------------------
void nas_itti_dedicated_eps_bearer_complete(
  const mme_ue_s1ap_id_t ue_idP,
  const ebi_t ebiP)
{
  OAILOG_FUNC_IN(LOG_NAS);
  MessageDef *message_p =
    itti_alloc_new_message(TASK_NAS_MME, MME_APP_CREATE_DEDICATED_BEARER_RSP);
  MME_APP_CREATE_DEDICATED_BEARER_RSP(message_p).ue_id = ue_idP;
  MME_APP_CREATE_DEDICATED_BEARER_RSP(message_p).ebi = ebiP;
  itti_send_msg_to_task(TASK_MME_APP, INSTANCE_DEFAULT, message_p);
  OAILOG_FUNC_OUT(LOG_NAS);
}

//------------------------------------------------------------------------------
void nas_itti_dedicated_eps_bearer_reject(
  const mme_ue_s1ap_id_t ue_idP,
  const ebi_t ebiP)
{
  OAILOG_FUNC_IN(LOG_NAS);
  MessageDef *message_p =
    itti_alloc_new_message(TASK_NAS_MME, MME_APP_CREATE_DEDICATED_BEARER_REJ);
  MME_APP_CREATE_DEDICATED_BEARER_REJ(message_p).ue_id = ue_idP;
  MME_APP_CREATE_DEDICATED_BEARER_REJ(message_p).ebi = ebiP;
  itti_send_msg_to_task(TASK_MME_APP, INSTANCE_DEFAULT, message_p);
  OAILOG_FUNC_OUT(LOG_NAS);
}

//------------------------------------------------------------------------------
void nas_itti_pdn_config_req(
  int ptiP,
  unsigned int ue_idP,
  const imsi_t *const imsi_pP,
  esm_proc_data_t *proc_data_pP,
  esm_proc_pdn_request_t request_typeP)
{
  OAILOG_FUNC_IN(LOG_NAS);
  MessageDef *message_p = NULL;

  AssertFatal(imsi_pP != NULL, "imsi_pP param is NULL");
  AssertFatal(proc_data_pP != NULL, "proc_data_pP param is NULL");

  OAILOG_INFO(
    LOG_NAS_EMM, "Sending PDN CONFIG REQ to MME_APP , (ue_idP = %u) \n",
    ue_idP);

  message_p = itti_alloc_new_message(TASK_NAS_MME, NAS_PDN_CONFIG_REQ);

  hexa_to_ascii(
    (uint8_t *) imsi_pP->u.value,
    NAS_PDN_CONFIG_REQ(message_p).imsi,
    (imsi_pP->length + 1) / 2);
  NAS_PDN_CONFIG_REQ(message_p).imsi_length = imsi_pP->length;

  NAS_PDN_CONFIG_REQ(message_p).ue_id = ue_idP;

  bassign(NAS_PDN_CONFIG_REQ(message_p).apn, proc_data_pP->apn);
  bassign(NAS_PDN_CONFIG_REQ(message_p).pdn_addr, proc_data_pP->pdn_addr);

  OAILOG_DEBUG(
    LOG_NAS_ESM,
    "PDN Type = (%d) for (ue_id = %u)\n ",
    proc_data_pP->pdn_type,
    ue_idP);
  switch (proc_data_pP->pdn_type) {
    case ESM_PDN_TYPE_IPV4:
      NAS_PDN_CONFIG_REQ(message_p).pdn_type = IPv4;
      break;

    case ESM_PDN_TYPE_IPV6:
      NAS_PDN_CONFIG_REQ(message_p).pdn_type = IPv6;
      break;

    case ESM_PDN_TYPE_IPV4V6:
      NAS_PDN_CONFIG_REQ(message_p).pdn_type = IPv4_AND_v6;
      break;

    default: NAS_PDN_CONFIG_REQ(message_p).pdn_type = IPv4; break;
  }

  NAS_PDN_CONFIG_REQ(message_p).request_type = request_typeP;

  itti_send_msg_to_task(TASK_MME_APP, INSTANCE_DEFAULT, message_p);
  OAILOG_FUNC_OUT(LOG_NAS);
}

//------------------------------------------------------------------------------
void nas_itti_pdn_connectivity_req(
  int ptiP,
  mme_ue_s1ap_id_t ue_idP,
  pdn_cid_t pdn_cidP,
  const imsi_t *const imsi_pP,
  imeisv_t imeisv,
  esm_proc_data_t *proc_data_pP,
  esm_proc_pdn_request_t request_typeP)
{
  OAILOG_FUNC_IN(LOG_NAS);
  MessageDef *message_p = NULL;

  AssertFatal(imsi_pP != NULL, "imsi_pP param is NULL");
  AssertFatal(proc_data_pP != NULL, "proc_data_pP param is NULL");

  message_p = itti_alloc_new_message(TASK_NAS_MME, NAS_PDN_CONNECTIVITY_REQ);

  hexa_to_ascii(
    (uint8_t *) imsi_pP->u.value, NAS_PDN_CONNECTIVITY_REQ(message_p).imsi, 8);

  NAS_PDN_CONNECTIVITY_REQ(message_p).pdn_cid = pdn_cidP;
  NAS_PDN_CONNECTIVITY_REQ(message_p).pti = ptiP;
  NAS_PDN_CONNECTIVITY_REQ(message_p).ue_id = ue_idP;
  NAS_PDN_CONNECTIVITY_REQ(message_p).imsi[imsi_pP->length] = '\0';
  NAS_PDN_CONNECTIVITY_REQ(message_p).imsi_length = imsi_pP->length;

  // Send IMEISV received in Security mode complete message
  NAS_PDN_CONNECTIVITY_REQ(message_p).presencemask |= NAS_PRESENT_IMEI_SV;
  NAS_PDN_CONNECTIVITY_REQ(message_p).imeisv = imeisv;

  bassign(NAS_PDN_CONNECTIVITY_REQ(message_p).apn, proc_data_pP->apn);
  bassign(NAS_PDN_CONNECTIVITY_REQ(message_p).pdn_addr, proc_data_pP->pdn_addr);

  switch (proc_data_pP->pdn_type) {
    case ESM_PDN_TYPE_IPV4:
      NAS_PDN_CONNECTIVITY_REQ(message_p).pdn_type = IPv4;
      break;

    case ESM_PDN_TYPE_IPV6:
      NAS_PDN_CONNECTIVITY_REQ(message_p).pdn_type = IPv6;
      break;

    case ESM_PDN_TYPE_IPV4V6:
      NAS_PDN_CONNECTIVITY_REQ(message_p).pdn_type = IPv4_AND_v6;
      break;

    default: NAS_PDN_CONNECTIVITY_REQ(message_p).pdn_type = IPv4; break;
  }

  // not efficient but be careful about "typedef network_qos_t esm_proc_qos_t;"
  memcpy(
    &NAS_PDN_CONNECTIVITY_REQ(message_p).bearer_qos,
    &proc_data_pP->bearer_qos,
    sizeof(proc_data_pP->bearer_qos));

  NAS_PDN_CONNECTIVITY_REQ(message_p).request_type = request_typeP;

  copy_protocol_configuration_options(
    &NAS_PDN_CONNECTIVITY_REQ(message_p).pco, &proc_data_pP->pco);

  OAILOG_INFO(
    LOG_NAS_ESM,
    "Sending PDN CONNECTIVITY REQ to MME_APP for ue id " MME_UE_S1AP_ID_FMT " \n",
    ue_idP);
  itti_send_msg_to_task(TASK_MME_APP, INSTANCE_DEFAULT, message_p);

  OAILOG_FUNC_OUT(LOG_NAS);
}

//------------------------------------------------------------------------------
void nas_itti_auth_info_req(
  const mme_ue_s1ap_id_t ue_idP,
  const imsi_t *const imsiP,
  const bool is_initial_reqP,
  plmn_t *const visited_plmnP,
  const uint8_t num_vectorsP,
  const_bstring const auts_pP)
{
  OAILOG_FUNC_IN(LOG_NAS);
  MessageDef *message_p = NULL;
  s6a_auth_info_req_t *auth_info_req = NULL;

  OAILOG_INFO(
    LOG_NAS_EMM, " Sending Authentication Information Request message to S6A for ue_id = (%u) \n",
    ue_idP);

  message_p = itti_alloc_new_message(TASK_NAS_MME, S6A_AUTH_INFO_REQ);
  auth_info_req = &message_p->ittiMsg.s6a_auth_info_req;
  memset(auth_info_req, 0, sizeof(s6a_auth_info_req_t));

  IMSI_TO_STRING(imsiP, auth_info_req->imsi, IMSI_BCD_DIGITS_MAX + 1);
  auth_info_req->imsi_length = (uint8_t) strlen(auth_info_req->imsi);

  AssertFatal(
    (auth_info_req->imsi_length > 5) && (auth_info_req->imsi_length < 16),
    "Bad IMSI length %d",
    auth_info_req->imsi_length);

  auth_info_req->visited_plmn = *visited_plmnP;
  auth_info_req->nb_of_vectors = num_vectorsP;

  if (is_initial_reqP) {
    auth_info_req->re_synchronization = 0;
    memset(auth_info_req->resync_param, 0, sizeof auth_info_req->resync_param);
  } else {
    AssertFatal(auts_pP != NULL, "Autn Null during resynchronization");
    auth_info_req->re_synchronization = 1;
    memcpy(
      auth_info_req->resync_param,
      auts_pP->data,
      sizeof auth_info_req->resync_param);
  }

  itti_send_msg_to_task(TASK_S6A, INSTANCE_DEFAULT, message_p);

  OAILOG_FUNC_OUT(LOG_NAS);
}

//------------------------------------------------------------------------------
void nas_itti_establish_rej(
  const mme_ue_s1ap_id_t ue_idP,
  const imsi_t *const imsi_pP,
  uint8_t initial_reqP)
{
  OAILOG_FUNC_IN(LOG_NAS);
  MessageDef *message_p;

  message_p =
    itti_alloc_new_message(TASK_NAS_MME, NAS_AUTHENTICATION_PARAM_REQ);

  hexa_to_ascii(
    (uint8_t *) imsi_pP->u.value,
    NAS_AUTHENTICATION_PARAM_REQ(message_p).imsi,
    8);

  NAS_AUTHENTICATION_PARAM_REQ(message_p).imsi[15] = '\0';

  if (isdigit(NAS_AUTHENTICATION_PARAM_REQ(message_p).imsi[14])) {
    NAS_AUTHENTICATION_PARAM_REQ(message_p).imsi_length = 15;
  } else {
    NAS_AUTHENTICATION_PARAM_REQ(message_p).imsi_length = 14;
    NAS_AUTHENTICATION_PARAM_REQ(message_p).imsi[14] = '\0';
  }

  NAS_AUTHENTICATION_PARAM_REQ(message_p).initial_req = initial_reqP;
  NAS_AUTHENTICATION_PARAM_REQ(message_p).ue_id = ue_idP;

  itti_send_msg_to_task(TASK_MME_APP, INSTANCE_DEFAULT, message_p);
  OAILOG_FUNC_OUT(LOG_NAS);
}

//------------------------------------------------------------------------------
void nas_itti_establish_cnf(
  const mme_ue_s1ap_id_t ue_idP,
  const nas_error_code_t error_codeP,
  bstring msgP,
  const uint16_t selected_encryption_algorithmP,
  const uint16_t selected_integrity_algorithmP,
  const uint8_t csfb_response,
  const uint8_t presencemask,
  const uint8_t service_type)
{
  OAILOG_FUNC_IN(LOG_NAS);
  MessageDef *message_p = NULL;
  ue_mm_context_t *ue_mm_context =
    mme_ue_context_exists_mme_ue_s1ap_id(&mme_app_desc.mme_ue_contexts, ue_idP);
  emm_context_t *emm_ctx = NULL;

  if (ue_mm_context) {
    emm_ctx = &ue_mm_context->emm_context;

    message_p =
      itti_alloc_new_message(TASK_NAS_MME, NAS_CONNECTION_ESTABLISHMENT_CNF);

    NAS_CONNECTION_ESTABLISHMENT_CNF(message_p).ue_id = ue_idP;
    NAS_CONNECTION_ESTABLISHMENT_CNF(message_p).err_code = error_codeP;
    NAS_CONNECTION_ESTABLISHMENT_CNF(message_p).nas_msg = msgP;
    msgP = NULL;
    NAS_CONNECTION_ESTABLISHMENT_CNF(message_p).csfb_response = csfb_response;
    NAS_CONNECTION_ESTABLISHMENT_CNF(message_p).presencemask = presencemask;
    NAS_CONNECTION_ESTABLISHMENT_CNF(message_p).service_type = service_type;

    // According to 3GPP 9.2.1.40, the UE security capabilities are 16-bit
    // strings, EEA0 is inherently supported, so its support is not tracked in
    // the bit string. However, emm_ctx->eea is an 8-bit string with the highest
    // order bit representing EEA0 support, so we need to trim it. The same goes
    // for integrity.
    //
    // TODO: change the way the EEA and EIA are translated into the packets.
    //       Currently, the 16-bit string is 8-bit rotated to produce the string
    //       sent in the packets, which is why we're using bits 8-10 to
    //       represent EEA1/2/3 (and EIA1/2/3) support here.
    NAS_CONNECTION_ESTABLISHMENT_CNF(message_p)
      .encryption_algorithm_capabilities =
      ((uint16_t) emm_ctx->_ue_network_capability.eea & ~(1 << 7)) << 1;
    NAS_CONNECTION_ESTABLISHMENT_CNF(message_p)
      .integrity_algorithm_capabilities =
      ((uint16_t) emm_ctx->_ue_network_capability.eia & ~(1 << 7)) << 1;

    AssertFatal(
      (0 <= emm_ctx->_security.vector_index) &&
        (MAX_EPS_AUTH_VECTORS > emm_ctx->_security.vector_index),
      "Invalid vector index %d",
      emm_ctx->_security.vector_index);

    derive_keNB(
      emm_ctx->_vector[emm_ctx->_security.vector_index].kasme,
      emm_ctx->_security.kenb_ul_count.seq_num |
        (emm_ctx->_security.kenb_ul_count.overflow << 8),
      NAS_CONNECTION_ESTABLISHMENT_CNF(message_p).kenb);
    /* Genarate Next HOP key parameter */
    derive_NH(
      emm_ctx->_vector[emm_ctx->_security.vector_index].kasme,
      NAS_CONNECTION_ESTABLISHMENT_CNF(message_p).kenb,
      emm_ctx->_security.next_hop,
      &emm_ctx->_security.next_hop_chaining_count);

    unlock_ue_contexts(ue_mm_context);
    OAILOG_INFO(LOG_NAS_EMM, "Sending NAS Connection Establishment confirm for ue_id "MME_UE_S1AP_ID_FMT"\n",
      ue_idP);
    itti_send_msg_to_task(TASK_MME_APP, INSTANCE_DEFAULT, message_p);
  } else {
    OAILOG_WARNING(
      LOG_NAS_EMM,
      "UE MM context NULL! for ue_id = (%u)\n",
      ue_idP);
  }

  OAILOG_FUNC_OUT(LOG_NAS);
}

//------------------------------------------------------------------------------
void nas_itti_detach_req(const mme_ue_s1ap_id_t ue_idP)
{
  OAILOG_FUNC_IN(LOG_NAS);
  MessageDef *message_p;

  message_p = itti_alloc_new_message(TASK_NAS_MME, NAS_DETACH_REQ);

  NAS_DETACH_REQ(message_p).ue_id = ue_idP;

  itti_send_msg_to_task(TASK_MME_APP, INSTANCE_DEFAULT, message_p);
  OAILOG_FUNC_OUT(LOG_NAS);
}

//------------------------------------------------------------------------------
void nas_itti_sgs_detach_req(const uint32_t ue_idP, const uint8_t detach_type)
{
  OAILOG_FUNC_IN(LOG_NAS);
  MessageDef *message_p;

  OAILOG_INFO(
    LOG_MME_APP,
    "Send SGS Detach Request to MME for ue_id = %u\n",
    ue_idP);
  message_p = itti_alloc_new_message(TASK_NAS_MME, NAS_SGS_DETACH_REQ);
  memset(
    &message_p->ittiMsg.nas_sgs_detach_req,
    0,
    sizeof(itti_nas_sgs_detach_req_t));

  NAS_SGS_DETACH_REQ(message_p).ue_id = ue_idP;
  NAS_SGS_DETACH_REQ(message_p).detach_type = detach_type;

  itti_send_msg_to_task(TASK_MME_APP, INSTANCE_DEFAULT, message_p);
  OAILOG_FUNC_OUT(LOG_NAS);
}

//***************************************************************************
void s6a_auth_info_rsp_timer_expiry_handler(void *args)
{
  OAILOG_FUNC_IN(LOG_NAS_EMM);

  emm_context_t *emm_ctx = (emm_context_t *) (args);

  if (emm_ctx) {
    nas_auth_info_proc_t *auth_info_proc =
      get_nas_cn_procedure_auth_info(emm_ctx);
    if (!auth_info_proc) {
      OAILOG_FUNC_OUT(LOG_NAS_EMM);
    }

    void *timer_callback_args = NULL;
    nas_stop_Ts6a_auth_info(
      auth_info_proc->ue_id, &auth_info_proc->timer_s6a, timer_callback_args);

    auth_info_proc->timer_s6a.id = NAS_TIMER_INACTIVE_ID;
    if (auth_info_proc->resync) {
      OAILOG_ERROR(
        LOG_NAS_EMM,
        "EMM-PROC  - Timer timer_s6_auth_info_rsp expired. Resync auth "
        "procedure was in progress. Aborting attach procedure. UE "
        "id " MME_UE_S1AP_ID_FMT "\n",
        auth_info_proc->ue_id);
    } else {
      OAILOG_ERROR(
        LOG_NAS_EMM,
        "EMM-PROC  - Timer timer_s6_auth_info_rsp expired. Initial auth "
        "procedure was in progress. Aborting attach procedure. UE "
        "id " MME_UE_S1AP_ID_FMT "\n",
        auth_info_proc->ue_id);
    }

    // Send Attach Reject with cause NETWORK FAILURE and delete UE context
    nas_proc_auth_param_fail(auth_info_proc->ue_id, NAS_CAUSE_NETWORK_FAILURE);
  } else {
    OAILOG_ERROR(
      LOG_NAS_EMM,
      "EMM-PROC  - Timer timer_s6_auth_info_rsp expired. Null EMM Context for "
      "UE \n");
  }

  OAILOG_FUNC_OUT(LOG_NAS_EMM);
}

void nas_itti_extended_service_req(
  const mme_ue_s1ap_id_t ue_id,
  const uint8_t servicetype,
  uint8_t csfb_response)
{
  OAILOG_FUNC_IN(LOG_NAS);
  MessageDef *message_p = NULL;

  message_p = itti_alloc_new_message(TASK_NAS_MME, NAS_EXTENDED_SERVICE_REQ);
  memset(
    &message_p->ittiMsg.nas_extended_service_req,
    0,
    sizeof(itti_nas_extended_service_req_t));
  NAS_EXTENDED_SERVICE_REQ(message_p).ue_id = ue_id;
  NAS_EXTENDED_SERVICE_REQ(message_p).servType = servicetype;
  NAS_EXTENDED_SERVICE_REQ(message_p).csfb_response = csfb_response;

  OAILOG_INFO(
    LOG_MME_APP,
    "Send NAS_EXTENDED_SERVICE_REQ from Nas to Mme-app for ue_id :%u\n",
    ue_id);
  itti_send_msg_to_task(TASK_MME_APP, INSTANCE_DEFAULT, message_p);

  OAILOG_FUNC_OUT(LOG_NAS);
}

void nas_itti_sgsap_uplink_unitdata(
  const char *const imsi,
  uint8_t imsi_len,
  bstring nas_msg,
  imeisv_t *imeisv_pP,
  MobileStationClassmark2 *mobilestationclassmark2_pP,
  tai_t *tai_pP,
  ecgi_t *ecgi_pP)
{
  OAILOG_FUNC_IN(LOG_NAS);
  MessageDef *message_p = NULL;
  int uetimezone = 0;

  message_p = itti_alloc_new_message(TASK_NAS_MME, SGSAP_UPLINK_UNITDATA);
  AssertFatal(message_p, "itti_alloc_new_message Failed");
  memset(
    &message_p->ittiMsg.sgsap_uplink_unitdata,
    0,
    sizeof(itti_sgsap_uplink_unitdata_t));
  memcpy(SGSAP_UPLINK_UNITDATA(message_p).imsi, imsi, imsi_len);
  SGSAP_UPLINK_UNITDATA(message_p).imsi[imsi_len] = '\0';
  SGSAP_UPLINK_UNITDATA(message_p).imsi_length = imsi_len;
  SGSAP_UPLINK_UNITDATA(message_p).nas_msg_container = nas_msg;
  nas_msg = NULL;
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
      (uint8_t *) imeisv_pP->u.value,
      SGSAP_UPLINK_UNITDATA(message_p).opt_imeisv,
      8);
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
      *((MobileStationClassmark2_t *) mobilestationclassmark2_pP);
    SGSAP_UPLINK_UNITDATA(message_p).presencemask |=
      UPLINK_UNITDATA_MOBILE_STATION_CLASSMARK_2_PARAMETER_PRESENT;
  }
  /*
   * optional - tai
   * update the tai presence bitmask.
   */
  if (tai_pP) {
    SGSAP_UPLINK_UNITDATA(message_p).opt_tai = *((tai_t *) tai_pP);
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

  itti_send_msg_to_task(TASK_SGS, INSTANCE_DEFAULT, message_p);

  OAILOG_FUNC_OUT(LOG_NAS);
}

void nas_itti_sgsap_tmsi_reallocation_comp(
  const char *imsi,
  const unsigned int imsi_len)
{
  OAILOG_FUNC_IN(LOG_NAS);
  MessageDef *message_p = NULL;

  message_p = itti_alloc_new_message(TASK_NAS_MME, SGSAP_TMSI_REALLOC_COMP);
  memset(
    &message_p->ittiMsg.sgsap_tmsi_realloc_comp,
    0,
    sizeof(itti_sgsap_tmsi_reallocation_comp_t));
  memcpy(SGSAP_TMSI_REALLOC_COMP(message_p).imsi, imsi, imsi_len);
  SGSAP_TMSI_REALLOC_COMP(message_p).imsi[imsi_len] = '\0';
  SGSAP_TMSI_REALLOC_COMP(message_p).imsi_length = imsi_len;
  itti_send_msg_to_task(TASK_SGS, INSTANCE_DEFAULT, message_p);

  OAILOG_FUNC_OUT(LOG_NAS);
}

//------------------------------------------------------------------------------
//Mapping between EMM Attach Type and EPS Attach Type
uint8_t _get_eps_attach_type(uint8_t emm_attach_type)
{
  OAILOG_FUNC_IN(LOG_NAS);
  uint8_t eps_attach_type = 0;

  switch (emm_attach_type) {
    case EMM_ATTACH_TYPE_EPS: eps_attach_type = EPS_ATTACH_TYPE_EPS; break;
    case EMM_ATTACH_TYPE_COMBINED_EPS_IMSI:
      eps_attach_type = EPS_ATTACH_TYPE_COMBINED_EPS_IMSI;
      break;
    case EMM_ATTACH_TYPE_EMERGENCY:
      eps_attach_type = EPS_ATTACH_TYPE_EMERGENCY;
      break;
    default:
      OAILOG_WARNING(LOG_NAS_EMM, " No Matching EPS Atttach type");
      break;
  }

  return eps_attach_type;
}
//------------------------------------------------------------------------------
/*SGS Location Update Request message to be sent to MME APP*/
void nas_itti_cs_domain_location_update_req(
  const uint32_t ue_idP,
  uint8_t msg_type)
{
  OAILOG_FUNC_IN(LOG_NAS);
  MessageDef *message_p = NULL;

  emm_context_t *emm_ctx = emm_context_get(&_emm_data, ue_idP);

  DevAssert(emm_ctx);
  message_p =
    itti_alloc_new_message(TASK_NAS_MME, NAS_CS_DOMAIN_LOCATION_UPDATE_REQ);
  memset(
    &message_p->ittiMsg.nas_cs_domain_location_update_req,
    0,
    sizeof(itti_nas_cs_domain_location_update_req_t));
  DevAssert(message_p);

  NAS_CS_DOMAIN_LOCATION_UPDATE_REQ(message_p).ue_id = ue_idP;

  if (msg_type == ATTACH_REQUEST) {
    NAS_CS_DOMAIN_LOCATION_UPDATE_REQ(message_p).attach_type =
      _get_eps_attach_type(emm_ctx->attach_type);
    ;
    NAS_CS_DOMAIN_LOCATION_UPDATE_REQ(message_p).msg_type |= ATTACH_REQUEST;
  } else if (msg_type == TRACKING_AREA_UPDATE_REQUEST) {
    NAS_CS_DOMAIN_LOCATION_UPDATE_REQ(message_p).tau_updt_type =
      emm_ctx->tau_updt_type;
    NAS_CS_DOMAIN_LOCATION_UPDATE_REQ(message_p).msg_type |= TAU_REQUEST;
  }
  NAS_CS_DOMAIN_LOCATION_UPDATE_REQ(message_p).add_updt_type =
    emm_ctx->additional_update_type;

  emm_context_unlock(emm_ctx);
  itti_send_msg_to_task(TASK_MME_APP, INSTANCE_DEFAULT, message_p);
  OAILOG_INFO(
    LOG_NAS_EMM, " Sent CS Domain Location Update Request to MME APP\n");

  OAILOG_FUNC_OUT(LOG_NAS);
}

/*TAU Complete message to be sent to MME APP*/
void nas_itti_tau_complete(unsigned int ue_idP)
{
  OAILOG_FUNC_IN(LOG_NAS);
  MessageDef *message_p = NULL;

  message_p = itti_alloc_new_message(TASK_NAS_MME, NAS_TAU_COMPLETE);
  memset(
    &message_p->ittiMsg.nas_tau_complete, 0, sizeof(itti_nas_tau_complete_t));

  NAS_TAU_COMPLETE(message_p).ue_id = ue_idP;

  itti_send_msg_to_task(TASK_MME_APP, INSTANCE_DEFAULT, message_p);

  OAILOG_FUNC_OUT(LOG_NAS);
}

void nas_itti_sgsap_ue_activity_ind(
  const char *imsi,
  const unsigned int imsi_len)
{
  OAILOG_FUNC_IN(LOG_NAS);
  MessageDef *message_p = NULL;

  message_p = itti_alloc_new_message(TASK_NAS_MME, SGSAP_UE_ACTIVITY_IND);
  memset(
    &message_p->ittiMsg.sgsap_ue_activity_ind,
    0,
    sizeof(itti_sgsap_ue_activity_ind_t));
  memcpy(SGSAP_UE_ACTIVITY_IND(message_p).imsi, imsi, imsi_len);
  SGSAP_UE_ACTIVITY_IND(message_p).imsi[imsi_len] = '\0';
  SGSAP_UE_ACTIVITY_IND(message_p).imsi_length = imsi_len;
  itti_send_msg_to_task(TASK_SGS, INSTANCE_DEFAULT, message_p);
  OAILOG_DEBUG(
    LOG_NAS,
    " Sending NAS ITTI SGSAP UE ACTIVITY IND to SGS task for Imsi : %s \n",
    imsi);

  OAILOG_FUNC_OUT(LOG_NAS);
}

//------------------------------------------------------------------------------
void nas_itti_deactivate_eps_bearer_context(
  const mme_ue_s1ap_id_t ue_idP,
  const ebi_t ebiP,
  bool delete_default_bearer,
  teid_t s_gw_teid_s11_s4)
{
  OAILOG_FUNC_IN(LOG_NAS);
  MessageDef *message_p =
    itti_alloc_new_message(TASK_NAS_MME, MME_APP_DELETE_DEDICATED_BEARER_RSP);
  MME_APP_DELETE_DEDICATED_BEARER_RSP(message_p).ue_id = ue_idP;
  MME_APP_DELETE_DEDICATED_BEARER_RSP(message_p).ebi[0] = ebiP;
  MME_APP_DELETE_DEDICATED_BEARER_RSP(message_p).delete_default_bearer =
    delete_default_bearer;
  MME_APP_DELETE_DEDICATED_BEARER_RSP(message_p).s_gw_teid_s11_s4 =
    s_gw_teid_s11_s4;
  MME_APP_DELETE_DEDICATED_BEARER_RSP(message_p).no_of_bearers = 1;
  itti_send_msg_to_task(TASK_MME_APP, INSTANCE_DEFAULT, message_p);
  OAILOG_FUNC_OUT(LOG_NAS);
}

//------------------------------------------------------------------------------
void nas_itti_dedicated_eps_bearer_deactivation_reject(
  const mme_ue_s1ap_id_t ue_idP,
  const ebi_t ebiP,
  bool delete_default_bearer,
  teid_t s_gw_teid_s11_s4)
{
  OAILOG_FUNC_IN(LOG_NAS);
  MessageDef *message_p =
    itti_alloc_new_message(TASK_NAS_MME, MME_APP_DELETE_DEDICATED_BEARER_REJ);
  MME_APP_DELETE_DEDICATED_BEARER_REJ(message_p).ue_id = ue_idP;
  MME_APP_DELETE_DEDICATED_BEARER_REJ(message_p).no_of_bearers = 1;
  MME_APP_DELETE_DEDICATED_BEARER_REJ(message_p).ebi[0] = ebiP;
  MME_APP_DELETE_DEDICATED_BEARER_REJ(message_p).delete_default_bearer =
    delete_default_bearer;
  MME_APP_DELETE_DEDICATED_BEARER_REJ(message_p).s_gw_teid_s11_s4 =
    s_gw_teid_s11_s4;
  itti_send_msg_to_task(TASK_MME_APP, INSTANCE_DEFAULT, message_p);
  OAILOG_FUNC_OUT(LOG_NAS);
}


//***************************************************************************

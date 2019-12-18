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
#include "mme_app_state.h"
#include "mme_app_messages_types.h"
#include "nas_messages_types.h"
#include "nas_procedures.h"
#include "nas_timer.h"
#include "s6a_messages_types.h"
#include "nas/securityDef.h"
#include "sgs_messages_types.h"

#define TASK_ORIGIN TASK_NAS_MME

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

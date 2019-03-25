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

/*! \file mme_app_sgsap_location_update.c
*/

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <stdbool.h>
#include <stdint.h>

#include "assertions.h"
#include "common_types.h"
#include "conversions.h"
#include "log.h"
#include "intertask_interface.h"
#include "mme_app_ue_context.h"
#include "mme_app_defs.h"
#include "mme_config.h"
#include "timer.h"
#include "3gpp_24.301.h"
#include "3gpp_24.008.h"
#include "mme_app_sgs_fsm.h"
#include "mme_app_itti_messaging.h"
#include "mme_app_sgs_messages.h"
#include "3gpp_23.003.h"
#include "3gpp_36.401.h"
#include "EpsAttachType.h"
#include "EpsUpdateType.h"
#include "bstrlib.h"
#include "common_defs.h"
#include "common_ies.h"
#include "intertask_interface_types.h"
#include "itti_types.h"
#include "mme_api.h"
#include "mme_app_desc.h"
#include "nas_messages_types.h"
#include "s1ap_messages_types.h"
#include "sgs_messages_types.h"

/********************************************************************************
 **                                                                            **
 ** Name:               mme_app_send_itti_sgsap_ue_activity_ind()              **
 ** Description         Send UE Activity Indication Message to SGS Task        **
 **                                                                            **
 ** Inputs:              Mobile Id                                             **
 **                                                                            **
********************************************************************************/

void _mme_app_send_itti_sgsap_ue_activity_ind(
  const char *imsi, const unsigned int imsi_len)
{
  OAILOG_FUNC_IN(LOG_NAS);
  MessageDef *message_p = NULL;

  message_p = itti_alloc_new_message(TASK_MME_APP, SGSAP_UE_ACTIVITY_IND);
  memset(&message_p->ittiMsg.sgsap_ue_activity_ind, 0,
         sizeof(itti_sgsap_ue_activity_ind_t));
  memcpy(SGSAP_UE_ACTIVITY_IND (message_p).imsi, imsi, imsi_len);
  OAILOG_DEBUG(LOG_NAS," Imsi : %s %d \n", imsi,imsi_len);
  SGSAP_UE_ACTIVITY_IND (message_p).imsi[imsi_len] = '\0';
  SGSAP_UE_ACTIVITY_IND (message_p).imsi_length = imsi_len;
  itti_send_msg_to_task(TASK_SGS, INSTANCE_DEFAULT, message_p);
  OAILOG_DEBUG(LOG_NAS,
     "Sending NAS ITTI SGSAP UE ACTIVITY IND to SGS task for Imsi : 
     %s %d \n", imsi,imsi_len);

  OAILOG_FUNC_OUT(LOG_NAS);
}


/**********************************************************************************
 **                                                                              **
 ** Name:               _copy_mobile_identity_helper()                           **
 ** Description          Helper function to copy Mobile Identity                 **
 **                                                                              **
 ** Inputs:              Mobile Id                                               **
 **                                                                              **
***********************************************************************************/
static int _copy_mobile_identity_helper(
  MobileIdentity_t *mobileid_dest,
  MobileIdentity_t *mobileid_src)
{
  OAILOG_FUNC_IN(LOG_MME_APP);

  if (mobileid_src->typeofidentity == MOBILE_IDENTITY_IMSI) {
    memcpy(mobileid_dest->u.imsi, mobileid_src->u.imsi, mobileid_src->length);
  } else if (mobileid_src->typeofidentity == MOBILE_IDENTITY_TMSI) {
    memcpy(mobileid_dest->u.tmsi, mobileid_src->u.tmsi, mobileid_src->length);
  }
  mobileid_dest->length = mobileid_src->length;
  mobileid_dest->typeofidentity = mobileid_src->typeofidentity;

  OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNok);
}

/**********************************************************************************
 **                                                                              **
 ** Name:                _build_sgs_status()                                    **
 ** Description          encode SGS-STATUS message                               **
 **                      state                                                   **
 **                                                                              **
 ** Inputs:              imsi,imsi length,LAI,Mobile Identity,msg ide            **
 **                                                                              **
***********************************************************************************/

static int _build_sgs_status(
  char *imsi,
  uint8_t imsi_length,
  lai_t laicsfb,
  MobileIdentity_t *mobileid,
  uint8_t msg_id)
{
  int rc = RETURNok;

  OAILOG_FUNC_IN(LOG_MME_APP);

  MessageDef *message_p = NULL;
  message_p = itti_alloc_new_message(TASK_MME_APP, SGSAP_STATUS);
  itti_sgsap_status_t *sgsap_status = &message_p->ittiMsg.sgsap_status;
  memset((void *) sgsap_status, 0, sizeof(itti_sgsap_status_t));

  //Encode IMSI
  sgsap_status->presencemask = SGSAP_IMSI;
  memcpy(&sgsap_status->imsi, imsi, imsi_length);
  sgsap_status->imsi_length = imsi_length;

  //Encode Cause
  sgsap_status->cause = SGS_MSG_TYPE_NOT_COMPATIBLE_WITH_PROTOCOL_STATE;

  //Erroneous message
  sgsap_status->error_msg.msg_type = msg_id;

  //Encode IMSI
  memcpy(
    &sgsap_status->error_msg.u.sgsap_location_update_acc.imsi,
    imsi,
    imsi_length);
  sgsap_status->error_msg.u.sgsap_location_update_acc.imsi_length = imsi_length;
  sgsap_status->presencemask |= SGSAP_IMSI;

  //Encode LAI
  sgsap_status->error_msg.u.sgsap_location_update_acc.laicsfb.mccdigit2 =
    laicsfb.mccdigit2;
  sgsap_status->error_msg.u.sgsap_location_update_acc.laicsfb.mccdigit1 =
    laicsfb.mccdigit1;
  sgsap_status->error_msg.u.sgsap_location_update_acc.laicsfb.mncdigit3 =
    laicsfb.mncdigit3;
  sgsap_status->error_msg.u.sgsap_location_update_acc.laicsfb.mccdigit3 =
    laicsfb.mccdigit3;
  sgsap_status->error_msg.u.sgsap_location_update_acc.laicsfb.mncdigit2 =
    laicsfb.mncdigit2;
  sgsap_status->error_msg.u.sgsap_location_update_acc.laicsfb.mncdigit1 =
    laicsfb.mncdigit1;
  sgsap_status->error_msg.u.sgsap_location_update_acc.laicsfb.lac = laicsfb.lac;

  //Encode Mobile Identity
  if (mobileid) {
    _copy_mobile_identity_helper(
      &sgsap_status->error_msg.u.sgsap_location_update_acc.mobileid, mobileid);
    sgsap_status->error_msg.u.sgsap_location_update_acc.presencemask |=
      MOBILE_IDENTITY;
  }
  //Send STATUS message to SGS task
  rc = itti_send_msg_to_task(TASK_SGS, INSTANCE_DEFAULT, message_p);

  OAILOG_FUNC_RETURN(LOG_MME_APP, rc);
}

/**********************************************************************************
 **
 ** Name:                send_loc_updt_acc_to_nas()                              **
 ** Description          Upon receiving SGS_LOCATION_UPDATE_ACC                  **
 **                      to NAS                                                  **
 **                                                                              **
 ** Inputs:              itti_sgsap_location_update_acc                          **
 **
***********************************************************************************/
static int _send_cs_domain_loc_updt_acc_to_nas(
  itti_sgsap_location_update_acc_t *const itti_sgsap_location_update_acc,
  struct ue_mm_context_s *ue_context_p,
  bool is_sgs_assoc_exists)
{
  MessageDef *message_p = NULL;
  itti_nas_cs_domain_location_update_acc_t *itti_nas_location_update_acc_p;
  int rc = RETURNok;

  message_p =
    itti_alloc_new_message(TASK_MME_APP, NAS_CS_DOMAIN_LOCATION_UPDATE_ACC);
  itti_nas_location_update_acc_p =
    &message_p->ittiMsg.nas_cs_domain_location_update_acc;
  itti_nas_location_update_acc_p->ue_id = ue_context_p->mme_ue_s1ap_id;
  itti_nas_location_update_acc_p->is_sgs_assoc_exists = is_sgs_assoc_exists;

  /*This is to handle cases where we do not send SGS Location Update Request towards MSC/VLR
    as the association already exists*/
  if (NULL == itti_sgsap_location_update_acc) {
    if (is_sgs_assoc_exists == SGS_ASSOC_ACTIVE) {
      if (ue_context_p->granted_service == GRANTED_SERVICE_SMS_ONLY) {
        itti_nas_location_update_acc_p->add_updt_res =
          ADDITONAL_UPDT_RES_SMS_ONLY;
        itti_nas_location_update_acc_p->presencemask |= ADD_UPDT_TYPE;
      }
      rc = itti_send_msg_to_task(TASK_NAS_MME, INSTANCE_DEFAULT, message_p);
      OAILOG_FUNC_RETURN(LOG_MME_APP, rc);
    } else if (is_sgs_assoc_exists == SGS_ASSOC_INACTIVE) {
      OAILOG_ERROR(
        LOG_MME_APP, "Failed send CS domain Location Update Accept to NAS \n");
      OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNerror);
    }
  }

  //LAI
  itti_nas_location_update_acc_p->laicsfb.mccdigit2 =
    itti_sgsap_location_update_acc->laicsfb.mccdigit2;
  itti_nas_location_update_acc_p->laicsfb.mccdigit1 =
    itti_sgsap_location_update_acc->laicsfb.mccdigit1;
  itti_nas_location_update_acc_p->laicsfb.mncdigit3 =
    itti_sgsap_location_update_acc->laicsfb.mncdigit3;
  itti_nas_location_update_acc_p->laicsfb.mccdigit3 =
    itti_sgsap_location_update_acc->laicsfb.mccdigit3;
  itti_nas_location_update_acc_p->laicsfb.mncdigit2 =
    itti_sgsap_location_update_acc->laicsfb.mncdigit2;
  itti_nas_location_update_acc_p->laicsfb.mncdigit1 =
    itti_sgsap_location_update_acc->laicsfb.mncdigit1;
  itti_nas_location_update_acc_p->laicsfb.lac =
    itti_sgsap_location_update_acc->laicsfb.lac;

  // Mobile Identity
  if (itti_sgsap_location_update_acc->presencemask & SGSAP_MOBILE_IDENTITY) {
    itti_nas_location_update_acc_p->presencemask |= MOBILE_IDENTITY;
    if (
      itti_sgsap_location_update_acc->mobileid.typeofidentity ==
      MOBILE_IDENTITY_IMSI) {
      //Convert char* IMSI/TMSI to digit format to be sent in NAS message
      imsi_mobile_identity_t *imsi_p =
        &itti_nas_location_update_acc_p->mobileid.imsi;
      MOBILE_ID_CHAR_TO_MOBILE_ID_IMSI_NAS(
        itti_sgsap_location_update_acc->mobileid.u.imsi,
        imsi_p,
        itti_sgsap_location_update_acc->mobileid.length);
    } else if (
      itti_sgsap_location_update_acc->mobileid.typeofidentity ==
      MOBILE_IDENTITY_TMSI) {
      tmsi_mobile_identity_t *tmsi_p =
        &itti_nas_location_update_acc_p->mobileid.tmsi;
      MOBILE_ID_CHAR_TO_MOBILE_ID_TMSI_NAS(
        itti_sgsap_location_update_acc->mobileid.u.tmsi,
        tmsi_p,
        itti_sgsap_location_update_acc->mobileid.length);
    }
    itti_nas_location_update_acc_p->mobileid.imsi.typeofidentity =
      itti_sgsap_location_update_acc->mobileid.typeofidentity;
  }

  //Additional Update type
  if (ue_context_p->granted_service == GRANTED_SERVICE_SMS_ONLY) {
    itti_nas_location_update_acc_p->add_updt_res = ADDITONAL_UPDT_RES_SMS_ONLY;
    itti_nas_location_update_acc_p->presencemask |= ADD_UPDT_TYPE;
  }
  rc = itti_send_msg_to_task(TASK_NAS_MME, INSTANCE_DEFAULT, message_p);

  OAILOG_FUNC_RETURN(LOG_MME_APP, rc);
}

/**********************************************************************************
 **                                                                              **
 ** Name:                _is_combined_tau()                                      **
 ** Description          Handling combined tau                                   **
 **                                                                              **
 ** Inputs:              itti_nas_cs_domain_location_update_req_t msg            **
 **                      ue_context: UE context                                  **
***********************************************************************************/

static int _is_combined_tau(
  ue_mm_context_t *ue_context,
  itti_nas_cs_domain_location_update_req_t *const itti_nas_location_update_req)
{
  int rc = RETURNerror;

  DevAssert(ue_context != NULL);
  DevAssert(ue_context->sgs_context != NULL);

  ue_context->tau_updt_type = itti_nas_location_update_req->tau_updt_type;
  ue_context->sgs_context->ongoing_procedure = COMBINED_TAU;

  //Store granted service type based on TAU Update type & additional updt type
  //  if (ue_context->tau_updt_type == EPS_UPDATE_TYPE_COMBINED_TA_LA_UPDATING_WITH_IMSI_ATTACH) {
  if (
    (itti_nas_location_update_req->add_updt_type != MME_APP_SMS_ONLY) &&
    !(strcmp(
      (const char *) mme_config.non_eps_service_control->data, "CSFB_SMS"))) {
    ue_context->granted_service = GRANTED_SERVICE_CSFB_SMS;
    OAILOG_INFO(LOG_MME_APP, "Granted service is GRANTED_SERVICE_CSFB_SMS\n");
  } else if (itti_nas_location_update_req->add_updt_type == MME_APP_SMS_ONLY) {
    ue_context->granted_service = GRANTED_SERVICE_SMS_ONLY;
    OAILOG_INFO(LOG_MME_APP, "Granted service is GRANTED_SERVICE_SMS_ONLY\n");
  } else if (
    (itti_nas_location_update_req->add_updt_type != MME_APP_SMS_ONLY) &&
    !(strcmp((const char *) mme_config.non_eps_service_control->data, "SMS"))) {
    ue_context->granted_service = GRANTED_SERVICE_SMS_ONLY;
    OAILOG_INFO(LOG_MME_APP, "Granted service is GRANTED_SERVICE_SMS_ONLY\n");
  } else {
    ue_context->granted_service = GRANTED_SERVICE_EPS_ONLY;
    OAILOG_INFO(LOG_MME_APP, "Granted service is GRANTED_SERVICE_EPS_ONLY\n");
  }
  //  }

  /*
   * As per 29.118 section 4.3.3, if TAU is received with EPS UPdate as IMSI attach & the
   * SGS state is ASSOCIATED, send SGS Location Update Request to MSC/VLR & change the state
   * to SGS_LA_UPDATE_REQUESTED. If periodic TAU is received on ASSOCIATED state and
   * VLR reliable is et to false, send SGS Location Update Request to MSC/VLR & change the
   * state to SGS_LA_UPDATE_REQUESTED
   */
  if (
    (ue_context->sgs_context->sgs_state == SGS_NULL) &&
    (ue_context->sgs_context->vlr_reliable == false)) {
    OAILOG_INFO(LOG_MME_APP, "Either SGS is NULL or vlr_reliable is false\n");
    rc = RETURNok;
  } else if (ue_context->sgs_context->sgs_state == SGS_ASSOCIATED) {
    if (
      (ue_context->tau_updt_type ==
       EPS_UPDATE_TYPE_COMBINED_TA_LA_UPDATING_WITH_IMSI_ATTACH) ||
      ((ue_context->tau_updt_type == EPS_UPDATE_TYPE_PERIODIC_UPDATING) &&
       (ue_context->sgs_context->vlr_reliable == false))) {
      OAILOG_INFO(
        LOG_MME_APP,
        "In SGS_ASSOCIATED, tau_updt_type %d\n",
        ue_context->tau_updt_type);
      //ue_context->sgs_context->sgs_state = SGS_LA_UPDATE_REQUESTED;
      rc = RETURNok;
    }
    if (
      (ue_context->tau_updt_type == EPS_UPDATE_TYPE_COMBINED_TA_LA_UPDATING) ||
      (ue_context->tau_updt_type == EPS_UPDATE_TYPE_PERIODIC_UPDATING)) {
      if (ue_context->sgs_context->vlr_reliable == true) {
        OAILOG_INFO(
          LOG_MME_APP, "Did not send Location Update Request to MSC\n");
        /*No need to send Location Update Request as we are in associated state
         Send SGS_ASSOC_ACTIVE to NAS so that TAU accept is sent to UE*/
        if (
          RETURNerror == (_send_cs_domain_loc_updt_acc_to_nas(
                           NULL, ue_context, SGS_ASSOC_ACTIVE))) {
          OAILOG_ERROR(
            LOG_MME_APP,
            "Failed to send SGS Location update accept to NAS for "
            "UE" IMSI_64_FMT "\n",
            ue_context->imsi);
        }
        if ((mme_ue_context_get_ue_sgs_neaf(
           itti_nas_location_update_req->ue_id) == true)) {
          OAILOG_INFO(
            LOG_MME_APP,
            "Sending UE Activity Ind to MSC for UE ID %d\n",
            itti_nas_location_update_req->ue_id);
           /* neaf flag is true*/
           /* send the SGSAP Ue activity indication to MSC/VLR */
           char imsi_str[IMSI_BCD_DIGITS_MAX + 1];
           IMSI64_TO_STRING(ue_context->imsi, imsi_str,
                ue_context->imsi_len);
           _mme_app_send_itti_sgsap_ue_activity_ind(imsi_str,
                                                    strlen(imsi_str));
           mme_ue_context_update_ue_sgs_neaf(
              itti_nas_location_update_req->ue_id, false);
         }
      } else {
        rc = RETURNok;
      }
    }
  }
  OAILOG_FUNC_RETURN(LOG_MME_APP, rc);
}

/**********************************************************************************
 **
 ** Name:                mme_app_handle_sgs_location_update_req()                **
 ** Description          Upon receiving SGS_LOCATION_UPDATE_REQ                  **
 **                      create sgs context                                      **
 **                                                                              **
 ** Inputs:              nas_sgs_location_update_req                             **
 **
***********************************************************************************/
int mme_app_handle_nas_cs_domain_location_update_req(
  itti_nas_cs_domain_location_update_req_t *const itti_nas_location_update_req)
{
  OAILOG_FUNC_IN(LOG_MME_APP);
  OAILOG_INFO(
    LOG_MME_APP,
    "Creating SGS Context for UE %d\n",
    itti_nas_location_update_req->ue_id);

  /*Fetch UE context*/
  ue_mm_context_t *ue_context = mme_ue_context_exists_mme_ue_s1ap_id(
    &mme_app_desc.mme_ue_contexts, itti_nas_location_update_req->ue_id);
  /*Store Attch type in UE context*/
  ue_context->attach_type = itti_nas_location_update_req->attach_type;

  //Create SGS context
  if (ue_context->sgs_context == NULL) {
    ue_context->sgs_context = calloc(1, sizeof(sgs_context_t));
    if (ue_context->sgs_context == NULL) {
      OAILOG_ERROR(
        LOG_MME_APP,
        "Cannot create SGS Context for UE ID %d ",
        itti_nas_location_update_req->ue_id);
      OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNerror);
    }
    /*Initialize SGS context to default values*/
    ue_context->sgs_context->sgs_state = SGS_NULL;
    ue_context->sgs_context->vlr_reliable = false;
    ue_context->sgs_context->neaf = false;
    ue_context->sgs_context->ts6_1_timer.id = SGS_TIMER_INACTIVE_ID;
    ue_context->sgs_context->ts8_timer.id = SGS_TIMER_INACTIVE_ID;
    ue_context->sgs_context->ts9_timer.id = SGS_TIMER_INACTIVE_ID;
    ue_context->sgs_context->ts10_timer.id = SGS_TIMER_INACTIVE_ID;
    ue_context->sgs_context->ts13_timer.id = SGS_TIMER_INACTIVE_ID;
    ue_context->sgs_context->ts6_1_timer.sec = mme_config.sgs_config.ts6_1_sec;
    ue_context->sgs_context->ts8_timer.id = SGS_TIMER_INACTIVE_ID;
    ue_context->sgs_context->ts8_timer.sec = mme_config.sgs_config.ts8_sec;
    ue_context->sgs_context->ts8_retransmission_count = 0;
    ue_context->sgs_context->ts9_timer.id = SGS_TIMER_INACTIVE_ID;
    ue_context->sgs_context->ts9_timer.sec = mme_config.sgs_config.ts9_sec;
    ue_context->sgs_context->ts9_retransmission_count = 0;
    ue_context->sgs_context->ts10_timer.id = SGS_TIMER_INACTIVE_ID;
    ue_context->sgs_context->ts10_timer.sec = mme_config.sgs_config.ts10_sec;
    ue_context->sgs_context->ts10_retransmission_count = 0;
    ue_context->sgs_context->ts13_timer.id = SGS_TIMER_INACTIVE_ID;
    ue_context->sgs_context->ts13_timer.sec = mme_config.sgs_config.ts13_sec;
    ue_context->sgs_context->ts13_retransmission_count = 0;
    ue_context->sgs_context->call_cancelled = false;
  }

  //Store granted service type based on attach type & addition updt type
  mme_config_read_lock(&mme_config);
  if (itti_nas_location_update_req->msg_type == ATTACH_REQUEST) {
    if (ue_context->attach_type == EPS_ATTACH_TYPE_COMBINED_EPS_IMSI) {
      if (
        (itti_nas_location_update_req->add_updt_type != MME_APP_SMS_ONLY) &&
        !(strcmp(
          (const char *) mme_config.non_eps_service_control->data,
          "CSFB_SMS"))) {
        OAILOG_INFO(
          LOG_MME_APP, "Granted service is GRANTED_SERVICE_CSFB_SMS\n");
        ue_context->granted_service = GRANTED_SERVICE_CSFB_SMS;
      } else if (
        itti_nas_location_update_req->add_updt_type == MME_APP_SMS_ONLY) {
        OAILOG_INFO(
          LOG_MME_APP, "Granted service is  GRANTED_SERVICE_SMS_ONLY\n");
        ue_context->granted_service = GRANTED_SERVICE_SMS_ONLY;
      } else if (
        (itti_nas_location_update_req->add_updt_type != MME_APP_SMS_ONLY) &&
        !(strcmp(
          (const char *) mme_config.non_eps_service_control->data, "SMS"))) {
        OAILOG_INFO(
          LOG_MME_APP, "Granted service is  GRANTED_SERVICE_SMS_ONLY\n");
        ue_context->granted_service = GRANTED_SERVICE_SMS_ONLY;
      }
    } else {
      OAILOG_INFO(LOG_MME_APP, "Granted service is GRANTED_SERVICE_EPS_ONLY\n");
      ue_context->granted_service = GRANTED_SERVICE_EPS_ONLY;
    }
  }
  mme_config_unlock(&mme_config);

/*TODO CSFB, currently from HSS access_mode is rceived as PACKET_ONLY
 * For testing purpose we are commenting, later we will modify code as below
 */
#if 0
  if((ue_context->access_mode == NAM_PACKET_AND_CIRCUIT) &&
    (ue_context->sgs_context->ts6_1_timer.id == MME_APP_TIMER_INACTIVE_ID)) {
    /*Send SGSAP Location Update Request message to SGS task*/
    send_itti_sgsap_location_update_req(ue_context);
    OAILOG_DEBUG (LOG_MME_APP, "Sending Location Update message to SGS task");
  }else if(ue_context->sgs_context->ts6_1_timer.id != MME_APP_TIMER_INACTIVE_ID) {
    //Ignore the the messae as Location Update procedure is already triggered
    OAILOG_DEBUG (LOG_MME_APP, "Dropping the message as Location Update procedure is already triggered for UE %d\n",
    itti_nas_location_update_req->ue_id);
  }

#endif
  if (/*(ue_context->access_mode == NAM_PACKET_AND_CIRCUIT) && */
      (ue_context->sgs_context->ts6_1_timer.id == MME_APP_TIMER_INACTIVE_ID)) {
    // If we received combined TAU,set granted service and check if we have to send Location Update Request
    if (itti_nas_location_update_req->msg_type == TAU_REQUEST) {
      if (
        (_is_combined_tau(ue_context, itti_nas_location_update_req)) ==
        RETURNerror) {
        unlock_ue_contexts(ue_context);
        OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNok);
      }
    }
    /*Send SGSAP Location Update Request message to SGS task*/
    send_itti_sgsap_location_update_req(ue_context);
    OAILOG_DEBUG(
      LOG_MME_APP,
      "Sending Location Update message to SGS task with IMSI" IMSI_64_FMT "\n",
      ue_context->imsi);
  } else {
    //Ignore the the messae as Location Update procedure is already triggered
    OAILOG_WARNING(
      LOG_MME_APP,
      "Dropping the message as Location Update procedure is already triggered "
      "for UE %d\n",
      itti_nas_location_update_req->ue_id);
  }

  unlock_ue_contexts(ue_context);
  OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNok);
}

/**********************************************************************************
 **
 ** Name:                send_itti_sgsap_location_update_req()                   **
 ** Description          Send SGSAP_LOCATION_UPDATE_REQ to SGS task              **
 **                                                                              **
 ** Inputs:                                                                      **
 **
***********************************************************************************/
int send_itti_sgsap_location_update_req(ue_mm_context_t *ue_context)
{
  OAILOG_FUNC_IN(LOG_MME_APP);

  MessageDef *message_p = NULL;
  int rc = RETURNok;

  message_p = itti_alloc_new_message(TASK_MME_APP, SGSAP_LOCATION_UPDATE_REQ);
  itti_sgsap_location_update_req_t *sgsap_location_update_req =
    &message_p->ittiMsg.sgsap_location_update_req;
  memset(
    (void *) sgsap_location_update_req,
    0,
    sizeof(itti_sgsap_location_update_req_t));

  DevAssert(ue_context != NULL);
  DevAssert(ue_context->sgs_context != NULL);

  //IMSI
  IMSI64_TO_STRING(
    ue_context->imsi, sgsap_location_update_req->imsi, ue_context->imsi_len);
  sgsap_location_update_req->imsi_length = ue_context->imsi_len;

  //EPS Location update type
  //If Combined attach is received, set Location Update type as IMSI_ATTACH
  if (ue_context->sgs_context->ongoing_procedure == COMBINED_ATTACH) {
    sgsap_location_update_req->locationupdatetype = IMSI_ATTACH;
  }
  //If Combined TAU is received, set Location Update type based the tau_updt_type
  else if (ue_context->sgs_context->ongoing_procedure == COMBINED_TAU) {
    if (
      (ue_context->tau_updt_type == EPS_UPDATE_TYPE_COMBINED_TA_LA_UPDATING) ||
      (ue_context->tau_updt_type == EPS_UPDATE_TYPE_PERIODIC_UPDATING)) {
      sgsap_location_update_req->locationupdatetype = NORMAL_LOCATION_UPDATE;
    } else if (
      ue_context->tau_updt_type ==
      EPS_UPDATE_TYPE_COMBINED_TA_LA_UPDATING_WITH_IMSI_ATTACH) {
      sgsap_location_update_req->locationupdatetype = IMSI_ATTACH;
    }
  }
  //New LAI - Retrieve from conf
  mme_config_read_lock(&mme_config);
  sgsap_location_update_req->newlaicsfb.mccdigit2 = mme_config.lai.mccdigit2;
  sgsap_location_update_req->newlaicsfb.mccdigit1 = mme_config.lai.mccdigit1;
  sgsap_location_update_req->newlaicsfb.mncdigit3 = mme_config.lai.mncdigit3;
  sgsap_location_update_req->newlaicsfb.mccdigit3 = mme_config.lai.mccdigit3;
  sgsap_location_update_req->newlaicsfb.mncdigit2 = mme_config.lai.mncdigit2;
  sgsap_location_update_req->newlaicsfb.mncdigit1 = mme_config.lai.mncdigit1;
  sgsap_location_update_req->newlaicsfb.lac = mme_config.lai.lac;
  mme_config_unlock(&mme_config);

  //IMEISV
  sgsap_location_update_req->presencemask |= SGSAP_IMEISV;
  imeisv_t *imeisv = &ue_context->imeisv;
  IMEISV_TO_STRING(imeisv, sgsap_location_update_req->imeisv, MAX_IMEISV_SIZE);

  //TAI - TAI List currently not available in MME APP UE Context

  //ECGI
  sgsap_location_update_req->presencemask |= SGSAP_E_CGI;
  sgsap_location_update_req->ecgi.plmn.mcc_digit2 =
    ue_context->e_utran_cgi.plmn.mcc_digit2;
  sgsap_location_update_req->ecgi.plmn.mcc_digit1 =
    ue_context->e_utran_cgi.plmn.mcc_digit1;
  sgsap_location_update_req->ecgi.plmn.mnc_digit3 =
    ue_context->e_utran_cgi.plmn.mnc_digit3;
  sgsap_location_update_req->ecgi.plmn.mcc_digit3 =
    ue_context->e_utran_cgi.plmn.mcc_digit3;
  sgsap_location_update_req->ecgi.plmn.mnc_digit2 =
    ue_context->e_utran_cgi.plmn.mnc_digit2;
  sgsap_location_update_req->ecgi.plmn.mnc_digit1 =
    ue_context->e_utran_cgi.plmn.mnc_digit1;

  sgsap_location_update_req->ecgi.cell_identity.enb_id =
    ue_context->e_utran_cgi.cell_identity.enb_id;
  sgsap_location_update_req->ecgi.cell_identity.cell_id =
    ue_context->e_utran_cgi.cell_identity.cell_id;

  /*Send SGSAP Location Update Request to SGS task*/
  rc = itti_send_msg_to_task(TASK_SGS, INSTANCE_DEFAULT, message_p);

  /* update the neaf flag to false after sending the Location Update Request message to SGS */
  mme_ue_context_update_ue_sgs_neaf(ue_context->mme_ue_s1ap_id, false);

  if (ue_context->sgs_context == NULL) {
    OAILOG_ERROR(
      LOG_MME_APP,
      "SGS Context is NULL for UE ID %d ",
      ue_context->mme_ue_s1ap_id);
    OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNerror);
  }
  /* Start Ts6-1 timer and change SGS state to LA_UPDATE_REQUESTED*/
  sgs_fsm_set_status(
    ue_context->mme_ue_s1ap_id,
    ue_context->sgs_context,
    SGS_LA_UPDATE_REQUESTED);
  if (
    timer_setup(
      ue_context->sgs_context->ts6_1_timer.sec,
      0,
      TASK_MME_APP,
      INSTANCE_DEFAULT,
      TIMER_ONE_SHOT,
      (void *) &(ue_context->mme_ue_s1ap_id),
      sizeof(mme_ue_s1ap_id_t),
      &(ue_context->sgs_context->ts6_1_timer.id)) < 0) {
    OAILOG_ERROR(
      LOG_MME_APP,
      "Failed to start Ts6-1 timer for UE id  %d \n",
      ue_context->mme_ue_s1ap_id);
    ue_context->sgs_context->ts6_1_timer.id = MME_APP_TIMER_INACTIVE_ID;
  } else {
    OAILOG_DEBUG(
      LOG_MME_APP,
      "MME APP : Sent SGsAP Location Update Request and Started Ts6-1 timer "
      "for UE id  %d \n",
      ue_context->mme_ue_s1ap_id);
  }
  OAILOG_FUNC_RETURN(LOG_MME_APP, rc);
}

/**********************************************************************************
 **
 ** Name:                map_sgs_emm_cause()                                     **
 ** Description          Maps SGS Reject cause to EMM cause                      **
 **                                                                              **
 ** Inputs:              SGS Reject Cause                                        **
 **
***********************************************************************************/

int map_sgs_emm_cause(SgsRejectCause_t sgs_cause)
{
  int emm_cause;
  switch (sgs_cause) {
    case SGS_IMSI_UNKNOWN_IN_HLR: {
      emm_cause = EMM_CAUSE_IMSI_UNKNOWN_IN_HSS;
    } break;
    case SGS_ILLEGAL_MS: {
      emm_cause = EMM_CAUSE_ILLEGAL_UE;
    } break;
    case SGS_IMSI_UNKNOWN_IN_VLR: {
      emm_cause = EMM_CAUSE_IMSI_UNKNOWN_IN_HSS;
    } break;
    case SGS_IMEI_NOT_ACCEPTED: {
      emm_cause = EMM_CAUSE_IMEI_NOT_ACCEPTED;
    } break;
    case SGS_ILLEGAL_UE: {
      emm_cause = EMM_CAUSE_ILLEGAL_UE;
    } break;
    case SGS_PLMN_NOT_ALLOWED: {
      emm_cause = EMM_CAUSE_PLMN_NOT_ALLOWED;
    } break;
    case SGS_LOCATION_AREA_NOT_ALLOWED: {
      emm_cause = EMM_CAUSE_TA_NOT_ALLOWED;
    } break;
    case SGS_ROAMING_NOT_ALLOWED_IN_THIS_LOCATION_AREA: {
      emm_cause = EMM_CAUSE_ROAMING_NOT_ALLOWED;
    } break;
    case SGS_NO_SUITABLE_CELLS_IN_LOCATION_AREA: {
      emm_cause = EMM_CAUSE_NO_SUITABLE_CELLS;
    } break;
    case SGS_NETWORK_FAILURE: {
      emm_cause = EMM_CAUSE_NETWORK_FAILURE;
    } break;
    case SGS_MAC_FAILURE: {
      emm_cause = EMM_CAUSE_MAC_FAILURE;
    } break;
    case SGS_SYNCH_FAILURE: {
      emm_cause = EMM_CAUSE_SYNCH_FAILURE;
    } break;
    case SGS_CONGESTION: {
      emm_cause = EMM_CAUSE_CONGESTION;
    } break;
    case SGS_GSM_AUTHENTICATION_UNACCEPTABLE: {
      emm_cause = EMM_CAUSE_NON_EPS_AUTH_UNACCEPTABLE;
    } break;
    case SGS_NOT_AUTHORIZED_FOR_THIS_CSG: {
      emm_cause = EMM_CAUSE_CSG_NOT_AUTHORIZED;
    } break;
    case SGS_SERVICE_OPTION_NOT_SUPPORTED: {
      emm_cause = EMM_CAUSE_CS_DOMAIN_NOT_AVAILABLE;
    } break;
    case SGS_REQUESTED_SERVICE_OPTION_NOT_SUBSCRIBED: {
      emm_cause = EMM_CAUSE_CS_DOMAIN_NOT_AVAILABLE;
    } break;
    case SGS_SERVICE_OPTION_TEMPORARILY_OUT_OF_ORDER: {
      emm_cause = EMM_CAUSE_CS_DOMAIN_NOT_AVAILABLE;
    } break;
    case SGS_CALL_CANNOT_BE_IDENTIFIED: {
      emm_cause = EMM_CAUSE_UE_IDENTITY_CANT_BE_DERIVED_BY_NW;
    } break;
    case SGS_RETRY_UPON_ENTRY_INTO_NEW_CELL: { /*TODO : Need to map appropriate cause*/
      emm_cause = 0;
    } break;
    case SGS_SEMANTICALLY_INCORRECT_MESSAGE: {
      emm_cause = EMM_CAUSE_SEMANTICALLY_INCORRECT;
    } break;
    case SGS_INVALID_MANDATORY_INFORMATION: {
      emm_cause = EMM_CAUSE_INVALID_MANDATORY_INFO;
    } break;
    case SGS_MSG_TYPE_NON_EXISTENT_NOT_IMPLEMENTED: {
      emm_cause = EMM_CAUSE_MESSAGE_TYPE_NOT_IMPLEMENTED;
    } break;
    case SGS_MSG_TYPE_NOT_COMPATIBLE_WITH_PROTOCOL_STATE: {
      emm_cause = EMM_CAUSE_MESSAGE_TYPE_NOT_COMPATIBLE;
    } break;
    case SGS_INFORMATION_ELEMENT_NON_EXISTENT_NOT_IMPLEMENTED: {
      emm_cause = EMM_CAUSE_IE_NOT_IMPLEMENTED;
    } break;
    case SGS_CONDITIONAL_IE_ERROR: {
      emm_cause = EMM_CAUSE_CONDITIONAL_IE_ERROR;
    } break;
    case SGS_MSG_NOT_COMPATIBLE_WITH_PROTOCOL_STATE: {
      emm_cause = EMM_CAUSE_MESSAGE_NOT_COMPATIBLE;
    } break;
    case SGS_PROTOCOL_ERROR_UNSPECIFIED: {
      emm_cause = EMM_CAUSE_PROTOCOL_ERROR;
    } break;
    case SGS_MSC_NOT_REACHABLE: {
      emm_cause = EMM_CAUSE_MSC_NOT_REACHABLE;
    } break;
    default:
      OAILOG_INFO(LOG_NAS_EMM, "Invalid SGS Reject cause\n");
      emm_cause = EMM_CAUSE_CS_DOMAIN_NOT_AVAILABLE;
  }
  return emm_cause;
}

/**********************************************************************************
 **
 ** Name:                send_loc_updt_fail_to_nas()                             **
 ** Description          Upon receiving SGS_LOCATION_UPDATE_FAIL                 **
 **                      to NAS                                                  **
 **                                                                              **
 ** Inputs:              itti_sgsap_location_update_fail                         **
 **
***********************************************************************************/
int send_cs_domain_loc_updt_fail_to_nas(
  SgsRejectCause_t cause,
  lai_t *lai,
  mme_ue_s1ap_id_t mme_ue_s1ap_id)
{
  MessageDef *message_p = NULL;
  itti_nas_cs_domain_location_update_fail_t *itti_nas_location_update_fail_p;
  int rc = RETURNok;

  message_p =
    itti_alloc_new_message(TASK_MME_APP, NAS_CS_DOMAIN_LOCATION_UPDATE_FAIL);
  itti_nas_location_update_fail_p =
    &message_p->ittiMsg.nas_cs_domain_location_update_fail;
  itti_nas_location_update_fail_p->ue_id = mme_ue_s1ap_id;

  //LAI
  if (lai) {
    itti_nas_location_update_fail_p->laicsfb.mccdigit2 = lai->mccdigit2;
    itti_nas_location_update_fail_p->laicsfb.mccdigit1 = lai->mccdigit1;
    itti_nas_location_update_fail_p->laicsfb.mncdigit3 = lai->mncdigit3;
    itti_nas_location_update_fail_p->laicsfb.mccdigit3 = lai->mccdigit3;
    itti_nas_location_update_fail_p->laicsfb.mncdigit2 = lai->mncdigit2;
    itti_nas_location_update_fail_p->laicsfb.mncdigit1 = lai->mncdigit1;
    itti_nas_location_update_fail_p->laicsfb.lac = lai->lac;
    itti_nas_location_update_fail_p->presencemask = LAI;
  }
  //SGS Reject Cause
  itti_nas_location_update_fail_p->reject_cause = map_sgs_emm_cause(cause);

  rc = itti_send_msg_to_task(TASK_NAS_MME, INSTANCE_DEFAULT, message_p);

  OAILOG_FUNC_RETURN(LOG_MME_APP, rc);
}

/**********************************************************************************
 **
 ** Name:                mme_app_handle_sgs_location_update_acc()                **
 ** Description          Upon receiving SGS_LOCATION_UPDATE_ACC                  **
 **                      send itti_nas_location_update_acc_p to na               **
 **                                                                              **
 ** Inputs:              nas_sgs_location_update_acc                             **
 **
***********************************************************************************/
int mme_app_handle_sgsap_location_update_acc(
  itti_sgsap_location_update_acc_t *const itti_sgsap_location_update_acc)
{
  imsi64_t imsi64 = INVALID_IMSI64;
  struct ue_mm_context_s *ue_context_p = NULL;
  int rc = RETURNok;
  sgs_fsm_t sgs_fsm;

  OAILOG_FUNC_IN(LOG_MME_APP);

  IMSI_STRING_TO_IMSI64(itti_sgsap_location_update_acc->imsi, &imsi64);
  ue_context_p =
    mme_ue_context_exists_imsi(&mme_app_desc.mme_ue_contexts, imsi64);
  if (ue_context_p == NULL) {
    OAILOG_ERROR(
      LOG_MME_APP,
      "Unknown IMSI in mme_app_handle_sgsap_location_update_acc\n");
    OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNerror);
  }
  if (ue_context_p->sgs_context == NULL) {
    OAILOG_ERROR(
      LOG_MME_APP,
      "SGS context not found in mme_app_handle_sgsap_location_update_acc\n");
    OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNerror);
  }

  sgs_fsm.primitive = _SGS_LOCATION_UPDATE_ACCEPT;
  sgs_fsm.ue_id = ue_context_p->mme_ue_s1ap_id;
  sgs_fsm.ctx = ue_context_p->sgs_context;
  ((sgs_context_t *) sgs_fsm.ctx)->sgsap_msg =
    (void *) itti_sgsap_location_update_acc;

  unlock_ue_contexts(ue_context_p);
  if (sgs_fsm_process(&sgs_fsm) != RETURNok) {
    OAILOG_ERROR(
      LOG_MME_APP,
      "Error in invoking FSM handler for primitive %d \n",
      sgs_fsm.primitive);
    OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNerror);
  }

  OAILOG_FUNC_RETURN(LOG_MME_APP, rc);
}

/**********************************************************************************
 **
 ** Name:                mme_app_handle_sgs_location_update_rej()                **
 ** Description          Upon receiving SGS_LOCATION_UPDATE_REJ                  **
 **                      send itti_location_update_fail to nas                   **
 **                                                                              **
 ** Inputs:              nas_sgs_location_update_rej                             **
 **
***********************************************************************************/
int mme_app_handle_sgsap_location_update_rej(
  itti_sgsap_location_update_rej_t *const itti_sgsap_location_update_rej)
{
  imsi64_t imsi64 = INVALID_IMSI64;
  int rc = RETURNok;
  struct ue_mm_context_s *ue_context_p = NULL;
  sgs_fsm_t sgs_fsm;

  OAILOG_FUNC_IN(LOG_MME_APP);
  OAILOG_INFO(LOG_MME_APP, "Received SGSAP LOCATION UPDATE REJECT \n");

  /*Fetch UE context*/
  IMSI_STRING_TO_IMSI64(itti_sgsap_location_update_rej->imsi, &imsi64);
  ue_context_p =
    mme_ue_context_exists_imsi(&mme_app_desc.mme_ue_contexts, imsi64);
  if (ue_context_p == NULL) {
    OAILOG_ERROR(
      LOG_MME_APP,
      "Unknown IMSI in mme_app_handle_sgsap_location_update_rej\n");
    mme_ue_context_dump_coll_keys();
    OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNerror);
  }

  if (ue_context_p->sgs_context == NULL) {
    OAILOG_ERROR(
      LOG_MME_APP,
      "SGS context not found in mme_app_handle_sgsap_location_update_rej\n");
    OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNerror);
  }
  sgs_fsm.primitive = _SGS_LOCATION_UPDATE_REJECT;
  sgs_fsm.ue_id = ue_context_p->mme_ue_s1ap_id;
  sgs_fsm.ctx = (void *) ue_context_p->sgs_context;
  ((sgs_context_t *) sgs_fsm.ctx)->sgsap_msg =
    (void *) itti_sgsap_location_update_rej;

  unlock_ue_contexts(ue_context_p);
  if (sgs_fsm_process(&sgs_fsm) != RETURNok) {
    OAILOG_ERROR(
      LOG_MME_APP,
      "Error in invoking FSM handler for primitive %d \n",
      sgs_fsm.primitive);
    OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNerror);
  }

  OAILOG_FUNC_RETURN(LOG_MME_APP, rc);
}

/**********************************************************************************
 **
 ** Name:                sgs_fsm_null_loc_updt_acc()                             **
 ** Description          Handling of SGS_LOCATION UPDATE ACCEPT in NULL          **
 **                      state                                                   **
 **                                                                              **
 ** Inputs:              sgs_fsm_t                                               **
 **                                                                              **
***********************************************************************************/
int sgs_fsm_null_loc_updt_acc(const sgs_fsm_t *fsm_evt)
{
  int rc = RETURNok;
  itti_sgsap_location_update_acc_t *itti_sgsap_location_update_acc_p = NULL;
  MobileIdentity_t *mobileid = NULL;

  OAILOG_FUNC_IN(LOG_MME_APP);
  sgs_context_t *sgs_context = (sgs_context_t *) fsm_evt->ctx;
  itti_sgsap_location_update_acc_p =
    (itti_sgsap_location_update_acc_t *) sgs_context->sgsap_msg;

  if (sgs_context == NULL) {
    OAILOG_ERROR(
      LOG_MME_APP, "SGS Context is NULL for UE ID %d ", fsm_evt->ue_id);
    OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNerror);
  }

  if (sgs_context->ts6_1_timer.id != SGS_TIMER_INACTIVE_ID) {
    OAILOG_DEBUG(
      LOG_MME_APP,
      "Dropping Location Update Accept as Ts6-1 timer is running\n");
  }
  // If we received Location Updt Accept from MSC/VLR and ts6_1_timer is not running
  else if (sgs_context->ts6_1_timer.id == SGS_TIMER_INACTIVE_ID) {
    /* If Ts8/Ts9 timer is running, drop the message
    *  If Ts8/Ts9 timer is not running and SGs state is SGs-ASSOCIATED, drop the message
    *  If Ts8/Ts9 timer is not running and the SGs state is SGs-NULL or LA-UPDATE_REQUESTED
    *  send SGs-Status message
    */
    OAILOG_DEBUG(
      LOG_MME_APP,
      "Received Location Update Accept when Ts6-1 timer is not running \n");
    if (
      (sgs_context->ts8_timer.id != MME_APP_TIMER_INACTIVE_ID) ||
      (sgs_context->ts9_timer.id != MME_APP_TIMER_INACTIVE_ID)) {
      OAILOG_DEBUG(
        LOG_MME_APP,
        "Dropping Location Update Accept as Ts8/Ts9 timer is running\n");
    } else if (
      (sgs_context->ts8_timer.id == MME_APP_TIMER_INACTIVE_ID) &&
      (sgs_context->ts9_timer.id == MME_APP_TIMER_INACTIVE_ID)) {
      OAILOG_DEBUG(LOG_MME_APP, "Send SGS-STATUS message\n");
      // Send SGS-STATUS message to SGS task
      if (
        itti_sgsap_location_update_acc_p->presencemask &
        SGSAP_MOBILE_IDENTITY) {
        mobileid = &itti_sgsap_location_update_acc_p->mobileid;
      }
      if (
        _build_sgs_status(
          itti_sgsap_location_update_acc_p->imsi,
          itti_sgsap_location_update_acc_p->imsi_length,
          itti_sgsap_location_update_acc_p->laicsfb,
          mobileid,
          SGsAP_LOCATION_UPDATE_ACCEPT) == RETURNok) {
        OAILOG_DEBUG(LOG_MME_APP, "SGS-STATUS message sent to SGS task\n");
      }
    }
  }
  OAILOG_FUNC_RETURN(LOG_MME_APP, rc);
}

/**********************************************************************************
 **
 ** Name:                sgs_fsm_associated_loc_updt_acc()                             **
 ** Description          Handling of SGS_LOCATION UPDATE ACCEPT in Associated    **
 **                      state                                                   **
 **                                                                              **
 ** Inputs:              sgs_fsm_t                                               **
 **                                                                              **
***********************************************************************************/
int sgs_fsm_associated_loc_updt_acc(const sgs_fsm_t *fsm_evt)
{
  int rc = RETURNok;

  OAILOG_FUNC_IN(LOG_MME_APP);

  sgs_context_t *sgs_context = (sgs_context_t *) fsm_evt->ctx;
  if (sgs_context == NULL) {
    OAILOG_ERROR(LOG_MME_APP, "Unknown UE ID %d ", fsm_evt->ue_id);
    OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNerror);
  }
  // If we received Location Updt Accept from MSC/VLR and ts6_1_timer is not running
  if (sgs_context->ts6_1_timer.id == SGS_TIMER_INACTIVE_ID) {
    /* If Ts8/Ts9 timer is running, drop the message
    *  If Ts8/Ts9 timer is not running and SGs state is SGs-ASSOCIATED, drop the message
    */
    OAILOG_DEBUG(
      LOG_MME_APP,
      "Received Location Update Accept when Ts6-1 timer is not running \n");
    OAILOG_DEBUG(
      LOG_MME_APP,
      "Dropping Location Update Accept message as it is received in SGS\
                                 Associated state\n");
  }

  OAILOG_FUNC_RETURN(LOG_MME_APP, rc);
}

/**********************************************************************************
 **
 ** Name:                sgs_fsm_la_updt_req_loc_updt_acc()                             **
 ** Description          Handling of SGS_LOCATION UPDATE ACCEPT in LA Update     **
 **                      Requested state                                         **
 **                                                                              **
 ** Inputs:              sgs_fsm_t                                               **
 **                                                                              **
***********************************************************************************/
int sgs_fsm_la_updt_req_loc_updt_acc(const sgs_fsm_t *fsm_evt)
{
  int rc = RETURNok;
  itti_sgsap_location_update_acc_t *itti_sgsap_location_update_acc_p = NULL;
  struct ue_mm_context_s *ue_context_p = NULL;
  MobileIdentity_t *mobileid = NULL;

  OAILOG_FUNC_IN(LOG_MME_APP);

  ue_context_p = mme_ue_context_exists_mme_ue_s1ap_id(
    &mme_app_desc.mme_ue_contexts, fsm_evt->ue_id);
  if (ue_context_p == NULL) {
    OAILOG_ERROR(LOG_MME_APP, "Unknown UE ID %d ", fsm_evt->ue_id);
    OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNerror);
  }
  sgs_context_t *sgs_context = (sgs_context_t *) fsm_evt->ctx;
  if (sgs_context == NULL) {
    OAILOG_ERROR(
      LOG_MME_APP, "SGS context not found for UE ID %d ", fsm_evt->ue_id);
    OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNerror);
  }
  itti_sgsap_location_update_acc_p =
    (itti_sgsap_location_update_acc_t *) sgs_context->sgsap_msg;
  if (sgs_context->ts6_1_timer.id != MME_APP_TIMER_INACTIVE_ID) {
    /*Update SGS context and set VLR Reliable to true*/
    sgs_context->sgs_state = SGS_ASSOCIATED;
    sgs_context->vlr_reliable = true;

    /*Stop Ts6-1 timer*/
    if (timer_remove(ue_context_p->sgs_context->ts6_1_timer.id, NULL)) {
      OAILOG_ERROR(LOG_MME_APP, "Failed to stop Ts6_1 timer \n");
    }
    sgs_context->ts6_1_timer.id = MME_APP_TIMER_INACTIVE_ID;
    /* Send Location Update Acc to NAS*/
    if (
      RETURNerror ==
      (_send_cs_domain_loc_updt_acc_to_nas(
        itti_sgsap_location_update_acc_p, ue_context_p, SGS_ASSOC_INACTIVE))) {
      OAILOG_DEBUG(
        LOG_MME_APP,
        "Failed to send SGS Location update accept to NAS for UE " IMSI_64_FMT
        "\n",
        ue_context_p->imsi);
      rc = RETURNerror;
    }
    unlock_ue_contexts(ue_context_p);
  }

  // If we received Location Updt Accept from MSC/VLR and ts6_1_timer is not running
  else if (sgs_context->ts6_1_timer.id == SGS_TIMER_INACTIVE_ID) {
    /* If Ts8/Ts9 timer is running, drop the message
    *  If Ts8/Ts9 timer is not running and SGs state is SGs-ASSOCIATED, drop the message
    *  If Ts8/Ts9 timer is not running and the SGs state is SGs-NULL or LA-UPDATE_REQUESTED
    *  send SGs-Status message
    */
    OAILOG_DEBUG(
      LOG_MME_APP,
      "Received Location Update Accept when Ts6-1 timer is not running \n");
    if (
      (sgs_context->ts8_timer.id != MME_APP_TIMER_INACTIVE_ID) ||
      (sgs_context->ts9_timer.id != MME_APP_TIMER_INACTIVE_ID)) {
      OAILOG_DEBUG(
        LOG_MME_APP,
        "Dropping Location Update Accept as Ts8/Ts9 timer is running\n");
    } else if (
      (sgs_context->ts8_timer.id == MME_APP_TIMER_INACTIVE_ID) &&
      (sgs_context->ts8_timer.id == MME_APP_TIMER_INACTIVE_ID)) {
      OAILOG_DEBUG(LOG_MME_APP, "Send SGS-STATUS message\n");
      if (
        itti_sgsap_location_update_acc_p->presencemask &
        SGSAP_MOBILE_IDENTITY) {
        mobileid = &itti_sgsap_location_update_acc_p->mobileid;
      }
      // Send SGS-STATUS message to SGS task
      if (
        _build_sgs_status(
          itti_sgsap_location_update_acc_p->imsi,
          itti_sgsap_location_update_acc_p->imsi_length,
          itti_sgsap_location_update_acc_p->laicsfb,
          mobileid,
          SGsAP_LOCATION_UPDATE_ACCEPT) == RETURNok) {
        OAILOG_DEBUG(LOG_MME_APP, "SGS-STATUS message sent to SGS task\n");
      }
    }
  }
  OAILOG_FUNC_RETURN(LOG_MME_APP, rc);
}

/**********************************************************************************
 **
 ** Name:                sgs_fsm_null_loc_updt_rej()                             **
 ** Description          Handling of SGS_LOCATION UPDATE REJECT in NULL          **
 **                      state                                                   **
 **                                                                              **
 ** Inputs:              sgs_fsm_t                                               **
 **                                                                              **
***********************************************************************************/
int sgs_fsm_null_loc_updt_rej(const sgs_fsm_t *fsm_evt)
{
  int rc = RETURNok;
  OAILOG_FUNC_IN(LOG_MME_APP);
  OAILOG_ERROR(
    LOG_MME_APP,
    "Dropping message as it is received in NULL state for UE %d",
    fsm_evt->ue_id);

  OAILOG_FUNC_RETURN(LOG_MME_APP, rc);
}

/****************************************************************************************
 **
 ** Name:                sgs_fsm_la_updt_req_loc_updt_rej                              **
 ** Description          Handling of SGS_LOCATION UPDATE REJECT in LA UPDT REQUESTED   **
 **                      state                                                         **
 **                                                                                    **
 ** Inputs:              sgs_fsm_t                                                     **
 **                                                                                    **
*****************************************************************************************/
int sgs_fsm_la_updt_req_loc_updt_rej(const sgs_fsm_t *fsm_evt)
{
  int rc = RETURNok;
  struct ue_mm_context_s *ue_context_p = NULL;
  itti_sgsap_location_update_rej_t *itti_sgsap_location_update_rej_p = NULL;
  imsi64_t imsi64 = INVALID_IMSI64;
  lai_t *lai = NULL;

  OAILOG_FUNC_IN(LOG_MME_APP);

  sgs_context_t *sgs_context = (sgs_context_t *) fsm_evt->ctx;
  itti_sgsap_location_update_rej_p =
    (itti_sgsap_location_update_rej_t *) sgs_context->sgsap_msg;
  IMSI_STRING_TO_IMSI64(itti_sgsap_location_update_rej_p->imsi, &imsi64);
  ue_context_p = mme_ue_context_exists_mme_ue_s1ap_id(
    &mme_app_desc.mme_ue_contexts, fsm_evt->ue_id);
  if (ue_context_p == NULL) {
    OAILOG_ERROR(LOG_MME_APP, "Unknown UE ID %d ", fsm_evt->ue_id);
    mme_ue_context_dump_coll_keys();
    OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNerror);
  }

  if (sgs_context == NULL) {
    OAILOG_ERROR(
      LOG_MME_APP, "SGS Context is NULL for UE ID %d ", fsm_evt->ue_id);
    OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNerror);
  }
  // Change SGS state to NULL
  sgs_context->sgs_state = SGS_NULL;
  /*Stop Ts6-1 timer*/

  if (sgs_context->ts6_1_timer.id != SGS_TIMER_INACTIVE_ID) {
    if (timer_remove(sgs_context->ts6_1_timer.id, NULL)) {
      OAILOG_ERROR(LOG_MME_APP, "Failed to stop Ts6_1 timer \n");
    }
    sgs_context->ts6_1_timer.id = SGS_TIMER_INACTIVE_ID;
  }

  if (itti_sgsap_location_update_rej_p->presencemask & SGSAP_LAI) {
    lai = &itti_sgsap_location_update_rej_p->laicsfb;
  }
  unlock_ue_contexts(ue_context_p);
  /* Send Location Update Failure to NAS*/
  send_cs_domain_loc_updt_fail_to_nas(
    itti_sgsap_location_update_rej_p->cause, lai, ue_context_p->mme_ue_s1ap_id);

  OAILOG_FUNC_RETURN(LOG_MME_APP, rc);
}

/**********************************************************************************
 **
 ** Name:                mme_app_handle_ts6_1_timer_expiry()                     **
 ** Description          Ts6_1 timer expiry handler                              **
 **                                                                              **
 ** Inputs:              ue_mm_context_s                                            **
 **                                                                              **
***********************************************************************************/

void mme_app_handle_ts6_1_timer_expiry(struct ue_mm_context_s *ue_context_p)
{
  OAILOG_FUNC_IN(LOG_MME_APP);
  DevAssert(ue_context_p != NULL);
  DevAssert(ue_context_p->sgs_context != NULL);
  OAILOG_WARNING(
    LOG_MME_APP,
    "Expired- Ts6-1 timer for UE id  %d \n",
    ue_context_p->mme_ue_s1ap_id);
  ue_context_p->sgs_context->ts6_1_timer.id = MME_APP_TIMER_INACTIVE_ID;
  ue_context_p->sgs_context->sgs_state = SGS_NULL;

  /* Send Location Update Failure to NAS*/
  send_cs_domain_loc_updt_fail_to_nas(
    SGS_MSC_NOT_REACHABLE, NULL, ue_context_p->mme_ue_s1ap_id);

  OAILOG_FUNC_OUT(LOG_MME_APP);
}

/**********************************************************************************
 **
 ** Name:                sgs_fsm_associated_loc_updt_rej()                       **
 ** Description          Handling of SGS_LOCATION UPDATE REJ in Associated       **
 **                      state                                                   **
 **                                                                              **
 ** Inputs:              sgs_fsm_t                                               **
 **                                                                              **
***********************************************************************************/
int sgs_fsm_associated_loc_updt_rej(const sgs_fsm_t *fsm_evt)
{
  int rc = RETURNok;
  struct ue_mm_context_s *ue_context_p = NULL;
  itti_sgsap_location_update_rej_t *itti_sgsap_location_update_rej_p = NULL;
  imsi64_t imsi64 = INVALID_IMSI64;

  OAILOG_FUNC_IN(LOG_MME_APP);
  sgs_context_t *sgs_context = (sgs_context_t *) fsm_evt->ctx;
  itti_sgsap_location_update_rej_p =
    (itti_sgsap_location_update_rej_t *) sgs_context->sgsap_msg;
  /*Fetch UE context*/
  IMSI_STRING_TO_IMSI64(itti_sgsap_location_update_rej_p->imsi, &imsi64);
  ue_context_p = mme_ue_context_exists_mme_ue_s1ap_id(
    &mme_app_desc.mme_ue_contexts, fsm_evt->ue_id);
  if (ue_context_p == NULL) {
    OAILOG_ERROR(LOG_MME_APP, "Unknown UE ID %d ", fsm_evt->ue_id);
    OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNerror);
  }
  if (sgs_context == NULL) {
    OAILOG_ERROR(
      LOG_MME_APP, "SGS context not found for UE ID %d ", fsm_evt->ue_id);
    OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNerror);
  }
  OAILOG_DEBUG(
    LOG_MME_APP,
    "Dropping Location Update Reject message as it is received in SGS\
                                 Associated state\n");

  OAILOG_FUNC_RETURN(LOG_MME_APP, rc);
}

/**********************************************************************************
 **
 ** Name:                mme_app_handle_tau_complete()                           **
 ** Description          Handling of TAU complete                                **
 **                                                                              **
 ** Inputs:              itti_nas_tau_complete_t msg                             **
 **                                                                              **
***********************************************************************************/
void mme_app_handle_nas_tau_complete(
  itti_nas_tau_complete_t *itti_nas_tau_complete_p)
{
  struct ue_mm_context_s *ue_context_p = NULL;
  mme_ue_s1ap_id_t ue_id = INVALID_MME_UE_S1AP_ID;

  OAILOG_FUNC_IN(LOG_MME_APP);
  DevAssert(itti_nas_tau_complete_p);

  ue_id = itti_nas_tau_complete_p->ue_id;
  if (ue_id == INVALID_MME_UE_S1AP_ID) {
    OAILOG_ERROR(
      LOG_MME_APP,
      "ERROR***** Invalid UE Id received from NAS in TAU Complete\n");
  }
  ue_context_p =
    mme_ue_context_exists_mme_ue_s1ap_id(&mme_app_desc.mme_ue_contexts, ue_id);
  if (ue_context_p) {
    if (ue_id != ue_context_p->mme_ue_s1ap_id) {
      OAILOG_ERROR(
        LOG_MME_APP,
        "ERROR***** Abnormal case: ue_id does not match with ue_id in "
        "ue_context %d, %d\n",
        ue_id,
        ue_context_p->mme_ue_s1ap_id);
      OAILOG_FUNC_OUT(LOG_MME_APP);
    }
  } else {
    OAILOG_ERROR(
      LOG_MME_APP,
      "ERROR***** Invalid UE Id received from NAS in TAU Complete %d\n",
      ue_id);
  }
  ue_context_p->ue_context_rel_cause = S1AP_NAS_NORMAL_RELEASE;
  // Notify S1AP to send UE Context Release Command to eNB.
  mme_app_itti_ue_context_release(
    ue_context_p, ue_context_p->ue_context_rel_cause);

  OAILOG_FUNC_OUT(LOG_MME_APP);
}

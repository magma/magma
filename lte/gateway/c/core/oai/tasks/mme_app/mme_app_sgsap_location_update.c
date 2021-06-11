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

/*! \file mme_app_sgsap_location_update.c
 */

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <stdbool.h>
#include <stdint.h>

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
#include "mme_app_state.h"
#include "s1ap_messages_types.h"
#include "sgs_messages_types.h"
#include "nas_proc.h"
#include "dynamic_memory_check.h"
#include "mme_app_timer.h"

/*******************************************************************************
 **                                                                           **
 ** Name:                _mme_app_update_granted_service_for_ue()             **
 ** Description          Based on supported features configured at MME and UE **
 **                      request services, sets the granted service for UE    **
 **                                                                           **
 ** Inputs:              Pointer to UE context                                **
 **                                                                           **
 ********************************************************************************/
void mme_app_update_granted_service_for_ue(ue_mm_context_t* ue_context) {
  OAILOG_FUNC_IN(LOG_MME_APP);
  additional_updt_t additional_update_type =
      (additional_updt_t) ue_context->emm_context.additional_update_type;
  mme_config_read_lock(&mme_config);

  if ((additional_update_type != MME_APP_SMS_ONLY) &&
      !(strcmp(
          (const char*) mme_config.non_eps_service_control->data,
          "CSFB_SMS"))) {
    ue_context->granted_service = GRANTED_SERVICE_CSFB_SMS;
    OAILOG_INFO(LOG_MME_APP, "Granted service is GRANTED_SERVICE_CSFB_SMS\n");
  } else if (additional_update_type == MME_APP_SMS_ONLY) {
    ue_context->granted_service = GRANTED_SERVICE_SMS_ONLY;
    OAILOG_INFO(LOG_MME_APP, "Granted service is GRANTED_SERVICE_SMS_ONLY\n");
  } else if (
      (additional_update_type != MME_APP_SMS_ONLY) &&
      (!(strcmp(
           (const char*) mme_config.non_eps_service_control->data, "SMS")) ||
       !(strcmp(
           (const char*) mme_config.non_eps_service_control->data,
           "SMS_ORC8R")))) {
    ue_context->granted_service = GRANTED_SERVICE_SMS_ONLY;
    OAILOG_INFO(LOG_MME_APP, "Granted service is  GRANTED_SERVICE_SMS_ONLY\n");
  } else {
    ue_context->granted_service = GRANTED_SERVICE_EPS_ONLY;
    OAILOG_INFO(LOG_MME_APP, "Granted service is GRANTED_SERVICE_EPS_ONLY\n");
  }
  if (ue_context->granted_service == GRANTED_SERVICE_SMS_ONLY) {
    ue_context->emm_context.csfbparams.additional_updt_res =
        ADDITONAL_UPDT_RES_SMS_ONLY;
    ue_context->emm_context.csfbparams.presencemask |= ADD_UPDATE_TYPE;
  }
  mme_config_unlock(&mme_config);
  OAILOG_FUNC_OUT(LOG_MME_APP);
}

/*******************************************************************************
 ** Name:                _get_eps_attach_type()                               **
 ** Description          Maps EMM attach type to EPS attach type              **
 **                                                                           **
 ** Inputs:              emm_attach_type: EMM attach type                     **
 ** Returns:             Mapped EPS attach type                               **
 **                                                                           **
 ********************************************************************************/
uint8_t get_eps_attach_type(uint8_t emm_attach_type) {
  OAILOG_FUNC_IN(LOG_MME_APP);
  uint8_t eps_attach_type = 0;

  switch (emm_attach_type) {
    case EMM_ATTACH_TYPE_EPS:
      eps_attach_type = EPS_ATTACH_TYPE_EPS;
      break;
    case EMM_ATTACH_TYPE_COMBINED_EPS_IMSI:
      eps_attach_type = EPS_ATTACH_TYPE_COMBINED_EPS_IMSI;
      break;
    case EMM_ATTACH_TYPE_EMERGENCY:
      eps_attach_type = EPS_ATTACH_TYPE_EMERGENCY;
      break;
    default:
      OAILOG_WARNING(LOG_MME_APP, " No Matching EPS Atttach type");
      break;
  }
  return eps_attach_type;
}

/******************************************************************************
 **                                                                          **
 ** Name:               mme_app_send_itti_sgsap_ue_activity_ind()            **
 ** Description         Send UE Activity Indication Message to SGS Task      **
 **                                                                          **
 ** Inputs:              Mobile Id                                           **
 **                                                                          **
 ******************************************************************************/
void mme_app_send_itti_sgsap_ue_activity_ind(
    const char* imsi, const unsigned int imsi_len) {
  OAILOG_FUNC_IN(LOG_MME_APP);
  MessageDef* message_p = NULL;

  message_p = itti_alloc_new_message(TASK_MME_APP, SGSAP_UE_ACTIVITY_IND);
  if (!message_p) {
    OAILOG_ERROR(
        LOG_MME_APP,
        "Failed to allocate memory for SGSAP UE ACTIVITY IND for Imsi: "
        "%s %d \n",
        imsi, imsi_len);
    OAILOG_FUNC_OUT(LOG_MME_APP);
  }
  memset(
      &message_p->ittiMsg.sgsap_ue_activity_ind, 0,
      sizeof(itti_sgsap_ue_activity_ind_t));
  memcpy(SGSAP_UE_ACTIVITY_IND(message_p).imsi, imsi, imsi_len);
  OAILOG_DEBUG(LOG_MME_APP, " Imsi : %s %d \n", imsi, imsi_len);
  SGSAP_UE_ACTIVITY_IND(message_p).imsi[imsi_len] = '\0';
  SGSAP_UE_ACTIVITY_IND(message_p).imsi_length    = imsi_len;
  if ((send_msg_to_task(&mme_app_task_zmq_ctx, TASK_SGS, message_p)) ==
      RETURNok) {
    OAILOG_DEBUG(
        LOG_MME_APP,
        "Sending ITTI SGSAP UE ACTIVITY IND to SGS task for Imsi: %s"
        " imsi_len: %d \n",
        imsi, imsi_len);
  } else {
    OAILOG_ERROR(
        LOG_MME_APP,
        "Failed to send ITTI SGSAP UE ACTIVITY IND to SGS task for Imsi: %s"
        " imsi_len: %d \n",
        imsi, imsi_len);
  }
  OAILOG_FUNC_OUT(LOG_MME_APP);
}

/**********************************************************************************
 ** **
 ** Name:               _copy_mobile_identity_helper() **
 ** Description          Helper function to copy Mobile Identity **
 ** **
 ** Inputs:              Mobile Id **
 ** **
 ***********************************************************************************/
static int copy_mobile_identity_helper(
    MobileIdentity_t* mobileid_dest, MobileIdentity_t* mobileid_src) {
  OAILOG_FUNC_IN(LOG_MME_APP);

  if (mobileid_src->typeofidentity == MOBILE_IDENTITY_IMSI) {
    memcpy(mobileid_dest->u.imsi, mobileid_src->u.imsi, mobileid_src->length);
  } else if (mobileid_src->typeofidentity == MOBILE_IDENTITY_TMSI) {
    memcpy(mobileid_dest->u.tmsi, mobileid_src->u.tmsi, mobileid_src->length);
  }
  mobileid_dest->length         = mobileid_src->length;
  mobileid_dest->typeofidentity = mobileid_src->typeofidentity;

  OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNok);
}

/**********************************************************************************
 ** **
 ** Name:                _build_sgs_status() **
 ** Description          encode SGS-STATUS message **
 **                      state **
 ** **
 ** Inputs:              imsi,imsi length,LAI,Mobile Identity,msg ide **
 ** **
 ***********************************************************************************/

static int build_sgs_status(
    char* imsi, uint8_t imsi_length, lai_t laicsfb, MobileIdentity_t* mobileid,
    uint8_t msg_id) {
  int rc = RETURNok;

  OAILOG_FUNC_IN(LOG_MME_APP);

  MessageDef* message_p = NULL;
  message_p             = itti_alloc_new_message(TASK_MME_APP, SGSAP_STATUS);
  itti_sgsap_status_t* sgsap_status = &message_p->ittiMsg.sgsap_status;
  memset((void*) sgsap_status, 0, sizeof(itti_sgsap_status_t));

  // Encode IMSI
  sgsap_status->presencemask = SGSAP_IMSI;
  memcpy(&sgsap_status->imsi, imsi, imsi_length);
  sgsap_status->imsi_length = imsi_length;

  // Encode Cause
  sgsap_status->cause = SGS_MSG_TYPE_NOT_COMPATIBLE_WITH_PROTOCOL_STATE;

  // Erroneous message
  sgsap_status->error_msg.msg_type = msg_id;

  // Encode IMSI
  memcpy(
      &sgsap_status->error_msg.u.sgsap_location_update_acc.imsi, imsi,
      imsi_length);
  sgsap_status->error_msg.u.sgsap_location_update_acc.imsi_length = imsi_length;
  sgsap_status->presencemask |= SGSAP_IMSI;

  // Encode LAI
  sgsap_status->error_msg.u.sgsap_location_update_acc.laicsfb.mccdigit1 =
      laicsfb.mccdigit1;
  sgsap_status->error_msg.u.sgsap_location_update_acc.laicsfb.mccdigit2 =
      laicsfb.mccdigit2;
  sgsap_status->error_msg.u.sgsap_location_update_acc.laicsfb.mccdigit3 =
      laicsfb.mccdigit3;
  sgsap_status->error_msg.u.sgsap_location_update_acc.laicsfb.mncdigit1 =
      laicsfb.mncdigit1;
  sgsap_status->error_msg.u.sgsap_location_update_acc.laicsfb.mncdigit2 =
      laicsfb.mncdigit2;
  sgsap_status->error_msg.u.sgsap_location_update_acc.laicsfb.mncdigit3 =
      laicsfb.mncdigit3;
  sgsap_status->error_msg.u.sgsap_location_update_acc.laicsfb.lac = laicsfb.lac;

  // Encode Mobile Identity
  if (mobileid) {
    copy_mobile_identity_helper(
        &sgsap_status->error_msg.u.sgsap_location_update_acc.mobileid,
        mobileid);
    sgsap_status->error_msg.u.sgsap_location_update_acc.presencemask |=
        MOBILE_IDENTITY;
  }
  // Send STATUS message to SGS task
  rc = send_msg_to_task(&mme_app_task_zmq_ctx, TASK_SGS, message_p);

  OAILOG_FUNC_RETURN(LOG_MME_APP, rc);
}
/*******************************************************************************
 **                                                                            *
 ** Name:                _handle_cs_domain_loc_updt_acc()                      *
 ** Description          Upon receiving SGS_LOCATION_UPDATE_ACC                *
 **                      Update csfb specific params to UE context and send    *
 **                      Attach Accept or TAU accept                           *
 **                                                                            *
 ** Inputs:              Received ITTI message, itti_sgsap_location_update_acc *
 **                      Pointer to UE context                                 *
 ** Outputs:                                                                   *
 **      Return:    RETURNok, RETURNerror                                      *
 *******************************************************************************/
static int handle_cs_domain_loc_updt_acc(
    itti_sgsap_location_update_acc_t* const itti_sgsap_location_update_acc,
    struct ue_mm_context_s* ue_context_p) {
  OAILOG_FUNC_IN(LOG_MME_APP);
  int rc                          = RETURNok;
  struct emm_context_s* emm_ctx_p = &ue_context_p->emm_context;

  if (!emm_ctx_p) {
    OAILOG_ERROR(LOG_MME_APP, "Invalid emm context \n");
    OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNerror);
  }
  // Store LAI to be sent in Attach Accept/TAU Accept
  emm_ctx_p->csfbparams.presencemask |= LAI_CSFB;
  emm_ctx_p->csfbparams.lai.mccdigit1 =
      itti_sgsap_location_update_acc->laicsfb.mccdigit1;
  emm_ctx_p->csfbparams.lai.mccdigit2 =
      itti_sgsap_location_update_acc->laicsfb.mccdigit2;
  emm_ctx_p->csfbparams.lai.mccdigit3 =
      itti_sgsap_location_update_acc->laicsfb.mccdigit3;
  emm_ctx_p->csfbparams.lai.mncdigit1 =
      itti_sgsap_location_update_acc->laicsfb.mncdigit1;
  emm_ctx_p->csfbparams.lai.mncdigit2 =
      itti_sgsap_location_update_acc->laicsfb.mncdigit2;
  emm_ctx_p->csfbparams.lai.mncdigit3 =
      itti_sgsap_location_update_acc->laicsfb.mncdigit3;
  emm_ctx_p->csfbparams.lai.lac = itti_sgsap_location_update_acc->laicsfb.lac;

  OAILOG_DEBUG(
      LOG_MME_APP, "MME-APP - Mobile Identity presence mask %u \n",
      itti_sgsap_location_update_acc->presencemask);

  // Store Mobile Identity to be sent in Attach Accept/TAU Accept
  if (itti_sgsap_location_update_acc->presencemask & SGSAP_MOBILE_IDENTITY) {
    emm_ctx_p->csfbparams.presencemask |= MOBILE_IDENTITY;
    if (itti_sgsap_location_update_acc->mobileid.typeofidentity ==
        MOBILE_IDENTITY_IMSI) {
      // Convert char* IMSI/TMSI to digit format to be sent in NAS message
      imsi_mobile_identity_t* imsi_p = &emm_ctx_p->csfbparams.mobileid.imsi;
      MOBILE_ID_CHAR_TO_MOBILE_ID_IMSI_NAS(
          itti_sgsap_location_update_acc->mobileid.u.imsi, imsi_p,
          itti_sgsap_location_update_acc->mobileid.length);
    } else if (
        itti_sgsap_location_update_acc->mobileid.typeofidentity ==
        MOBILE_IDENTITY_TMSI) {
      tmsi_mobile_identity_t received_tmsi = {0};
      MOBILE_ID_CHAR_TO_MOBILE_ID_TMSI_NAS(
          itti_sgsap_location_update_acc->mobileid.u.tmsi, (&received_tmsi),
          itti_sgsap_location_update_acc->mobileid.length);
      /* If the rcvd TMSI is different from the stored TMSI,
       * store the new TMSI and set flag
       */
      if (MME_APP_COMPARE_TMSI(
              emm_ctx_p->csfbparams.mobileid.tmsi, received_tmsi) ==
          RETURNerror) {
        OAILOG_INFO(LOG_MME_APP, "MME-APP - New TMSI Allocated\n");
        memcpy(
            &emm_ctx_p->csfbparams.mobileid.tmsi, (&received_tmsi),
            sizeof(itti_sgsap_location_update_acc->mobileid.u.tmsi));
        emm_ctx_p->csfbparams.newTmsiAllocated = true;
      }
    }
  }
  /* Store the status of Location Update procedure(success/failure) to send
   * appropriate cause in Attach Accept/TAU Accept
   */
  emm_ctx_p->csfbparams.sgs_loc_updt_status = SUCCESS;

  // Additional Update type
  if (ue_context_p->granted_service == GRANTED_SERVICE_SMS_ONLY) {
    emm_ctx_p->csfbparams.presencemask |= ADD_UPDATE_TYPE;
    emm_ctx_p->csfbparams.additional_updt_res = ADDITONAL_UPDT_RES_SMS_ONLY;
  }
  // Send Attach Accept/TAU Accept
  rc = emm_send_cs_domain_attach_or_tau_accept(ue_context_p);

  OAILOG_FUNC_RETURN(LOG_MME_APP, rc);
}

/*****************************************************************************
 **                                                                         **
 ** Name:                mme_app_handle_nas_cs_domain_location_update_req   **
 ** Description          Upon receiving SGS_LOCATION_UPDATE_REQ             **
 **                      create sgs context                                 **
 **                                                                         **
 ** Inputs:              Pointer to UE context                              **
 **                      Type of message:Attach Request or TAU Request      **
 **                                                                         **
 ******************************************************************************/
int mme_app_handle_nas_cs_domain_location_update_req(
    ue_mm_context_t* ue_context_p, uint8_t msg_type) {
  OAILOG_FUNC_IN(LOG_MME_APP);
  OAILOG_INFO(
      LOG_MME_APP,
      "Handling cs_domain_location_update_req for UE-ID" MME_UE_S1AP_ID_FMT
      "\n",
      ue_context_p->mme_ue_s1ap_id);

  // Create SGS context
  if (ue_context_p->sgs_context == NULL) {
    if ((mme_app_create_sgs_context(ue_context_p)) != RETURNok) {
      OAILOG_CRITICAL(
          LOG_MME_APP,
          "Failed to create SGS context for ue_id " MME_UE_S1AP_ID_FMT "\n",
          ue_context_p->mme_ue_s1ap_id);
      OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNerror);
    }
  }

  // Store granted service type based on attach type & addition updt type
  if (msg_type == ATTACH_REQUEST) {
    ue_context_p->attach_type =
        get_eps_attach_type(ue_context_p->emm_context.attach_type);
    ue_context_p->sgs_context->ongoing_procedure = COMBINED_ATTACH;
    if (ue_context_p->attach_type == EPS_ATTACH_TYPE_COMBINED_EPS_IMSI) {
      mme_app_update_granted_service_for_ue(ue_context_p);
    }
  }
  if ((ue_context_p->network_access_mode == NAM_PACKET_AND_CIRCUIT) &&
      (ue_context_p->sgs_context->ts6_1_timer.id ==
       MME_APP_TIMER_INACTIVE_ID)) {
    if (msg_type == TAU_REQUEST) {
      ue_context_p->sgs_context->ongoing_procedure = COMBINED_TAU;
      mme_app_update_granted_service_for_ue(ue_context_p);
    }
    OAILOG_INFO(
        LOG_MME_APP,
        "Sending Location Update message to SGS task with IMSI" IMSI_64_FMT
        "\n",
        ue_context_p->emm_context._imsi64);
    // Send SGSAP Location Update Request message to SGS task
    send_itti_sgsap_location_update_req(ue_context_p);
  } else if (
      ue_context_p->sgs_context->ts6_1_timer.id != MME_APP_TIMER_INACTIVE_ID) {
    // Ignore the the messae as Location Update procedure is already triggered
    OAILOG_WARNING(
        LOG_MME_APP,
        "Dropping the message as Location Update Req is already triggered"
        "for UE-ID" MME_UE_S1AP_ID_FMT "\n",
        ue_context_p->mme_ue_s1ap_id);
  }
  OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNok);
}

/**********************************************************************************
 **
 ** Name:                send_itti_sgsap_location_update_req() **
 ** Description          Send SGSAP_LOCATION_UPDATE_REQ to SGS task **
 ** **
 ** Inputs: **
 **
 ***********************************************************************************/
int send_itti_sgsap_location_update_req(ue_mm_context_t* ue_context_p) {
  OAILOG_FUNC_IN(LOG_MME_APP);

  MessageDef* message_p = NULL;
  int rc                = RETURNok;
  uint8_t tau_updt_type = -1;

  message_p = itti_alloc_new_message(TASK_MME_APP, SGSAP_LOCATION_UPDATE_REQ);
  itti_sgsap_location_update_req_t* sgsap_location_update_req =
      &message_p->ittiMsg.sgsap_location_update_req;
  memset(
      (void*) sgsap_location_update_req, 0,
      sizeof(itti_sgsap_location_update_req_t));

  // IMSI
  IMSI64_TO_STRING(
      ue_context_p->emm_context._imsi64, sgsap_location_update_req->imsi,
      ue_context_p->emm_context._imsi.length);
  sgsap_location_update_req->imsi_length =
      ue_context_p->emm_context._imsi.length;

  tau_updt_type = ue_context_p->emm_context.tau_updt_type;
  // EPS Location update type
  // If Combined attach is received, set Location Update type as IMSI_ATTACH
  if (ue_context_p->sgs_context->ongoing_procedure == COMBINED_ATTACH) {
    sgsap_location_update_req->locationupdatetype = IMSI_ATTACH;
  }
  // If Combined TAU is received, set Location Update type based the
  // tau_updt_type
  else if (ue_context_p->sgs_context->ongoing_procedure == COMBINED_TAU) {
    if ((tau_updt_type == EPS_UPDATE_TYPE_COMBINED_TA_LA_UPDATING) ||
        (tau_updt_type == EPS_UPDATE_TYPE_PERIODIC_UPDATING)) {
      sgsap_location_update_req->locationupdatetype = NORMAL_LOCATION_UPDATE;
    } else if (
        tau_updt_type ==
        EPS_UPDATE_TYPE_COMBINED_TA_LA_UPDATING_WITH_IMSI_ATTACH) {
      sgsap_location_update_req->locationupdatetype = IMSI_ATTACH;
    }
  }
  // New LAI - Retrieve from conf
  mme_config_read_lock(&mme_config);
  sgsap_location_update_req->newlaicsfb.mccdigit1 = mme_config.lai.mccdigit1;
  sgsap_location_update_req->newlaicsfb.mccdigit2 = mme_config.lai.mccdigit2;
  sgsap_location_update_req->newlaicsfb.mccdigit3 = mme_config.lai.mccdigit3;
  sgsap_location_update_req->newlaicsfb.mncdigit1 = mme_config.lai.mncdigit1;
  sgsap_location_update_req->newlaicsfb.mncdigit2 = mme_config.lai.mncdigit2;
  sgsap_location_update_req->newlaicsfb.mncdigit3 = mme_config.lai.mncdigit3;
  sgsap_location_update_req->newlaicsfb.lac       = mme_config.lai.lac;
  mme_config_unlock(&mme_config);

  // IMEISV
  sgsap_location_update_req->presencemask |= SGSAP_IMEISV;
  imeisv_t* imeisv = &ue_context_p->emm_context._imeisv;
  IMEISV_TO_STRING(imeisv, sgsap_location_update_req->imeisv, MAX_IMEISV_SIZE);

  // TAI - TAI List currently not available in MME APP UE Context

  // ECGI
  sgsap_location_update_req->presencemask |= SGSAP_E_CGI;
  sgsap_location_update_req->ecgi.plmn.mcc_digit1 =
      ue_context_p->e_utran_cgi.plmn.mcc_digit1;
  sgsap_location_update_req->ecgi.plmn.mcc_digit2 =
      ue_context_p->e_utran_cgi.plmn.mcc_digit2;
  sgsap_location_update_req->ecgi.plmn.mcc_digit3 =
      ue_context_p->e_utran_cgi.plmn.mcc_digit3;
  sgsap_location_update_req->ecgi.plmn.mnc_digit1 =
      ue_context_p->e_utran_cgi.plmn.mnc_digit1;
  sgsap_location_update_req->ecgi.plmn.mnc_digit2 =
      ue_context_p->e_utran_cgi.plmn.mnc_digit2;
  sgsap_location_update_req->ecgi.plmn.mnc_digit3 =
      ue_context_p->e_utran_cgi.plmn.mnc_digit3;

  sgsap_location_update_req->ecgi.cell_identity.enb_id =
      ue_context_p->e_utran_cgi.cell_identity.enb_id;
  sgsap_location_update_req->ecgi.cell_identity.cell_id =
      ue_context_p->e_utran_cgi.cell_identity.cell_id;

  // Send SGSAP Location Update Request to SGS task
  if ((send_msg_to_task(&mme_app_task_zmq_ctx, TASK_SGS, message_p)) !=
      RETURNok) {
    OAILOG_ERROR(
        LOG_MME_APP,
        "Failed to send SGS-Location Update Request for UE "
        "ID" MME_UE_S1AP_ID_FMT "\n",
        ue_context_p->mme_ue_s1ap_id);
    OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNerror);
  }

  /* update the neaf flag to false after sending the Location Update
   * Request message to SGS
   */
  mme_ue_context_update_ue_sgs_neaf(ue_context_p->mme_ue_s1ap_id, false);

  if (ue_context_p->sgs_context == NULL) {
    OAILOG_ERROR(
        LOG_MME_APP, "SGS Context is NULL for UE ID %d ",
        ue_context_p->mme_ue_s1ap_id);
    OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNerror);
  }
  /* Start Ts6-1 timer and change SGS state to LA_UPDATE_REQUESTED */
  sgs_fsm_set_status(
      ue_context_p->mme_ue_s1ap_id, ue_context_p->sgs_context,
      SGS_LA_UPDATE_REQUESTED);
  if ((ue_context_p->sgs_context->ts6_1_timer.id = mme_app_start_timer(
           ue_context_p->sgs_context->ts6_1_timer.sec * 1000, TIMER_REPEAT_ONCE,
           mme_app_handle_ts6_1_timer_expiry, ue_context_p->mme_ue_s1ap_id)) ==
      -1) {
    OAILOG_ERROR(
        LOG_MME_APP, "Failed to start Ts6-1 timer for UE id  %d \n",
        ue_context_p->mme_ue_s1ap_id);
    ue_context_p->sgs_context->ts6_1_timer.id = MME_APP_TIMER_INACTIVE_ID;
  } else {
    OAILOG_DEBUG(
        LOG_MME_APP,
        "MME APP : Sent SGsAP Location Update Request and Started Ts6-1 timer "
        "for UE id: " MME_UE_S1AP_ID_FMT "\n",
        ue_context_p->mme_ue_s1ap_id);
  }
  OAILOG_FUNC_RETURN(LOG_MME_APP, rc);
}

/*****************************************************************************
 **
 ** Name:                mme_app_handle_sgsap_location_update_acc()         **
 ** Description          Upon receiving SGS_LOCATION_UPDATE_ACC             **
 **                      Based on the state, invoke state machine handlers  **
 **                                                                         **
 ** Inputs:              nas_sgs_location_update_acc                        **
 **
 ******************************************************************************/
int mme_app_handle_sgsap_location_update_acc(
    mme_app_desc_t* mme_app_desc_p,
    itti_sgsap_location_update_acc_t* const itti_sgsap_location_update_acc) {
  imsi64_t imsi64                      = INVALID_IMSI64;
  struct ue_mm_context_s* ue_context_p = NULL;
  int rc                               = RETURNok;
  sgs_fsm_t sgs_fsm;

  OAILOG_FUNC_IN(LOG_MME_APP);

  IMSI_STRING_TO_IMSI64(itti_sgsap_location_update_acc->imsi, &imsi64);
  ue_context_p =
      mme_ue_context_exists_imsi(&mme_app_desc_p->mme_ue_contexts, imsi64);
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
  sgs_fsm.ue_id     = ue_context_p->mme_ue_s1ap_id;
  sgs_fsm.ctx       = ue_context_p->sgs_context;
  ((sgs_context_t*) sgs_fsm.ctx)->sgsap_msg =
      (void*) itti_sgsap_location_update_acc;

  if (sgs_fsm_process(&sgs_fsm) != RETURNok) {
    OAILOG_ERROR(
        LOG_MME_APP, "Error in invoking FSM handler for primitive %d \n",
        sgs_fsm.primitive);
    OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNerror);
  }

  OAILOG_FUNC_RETURN(LOG_MME_APP, rc);
}

/******************************************************************************
 **
 ** Name:                mme_app_handle_sgs_location_update_rej()            **
 ** Description          Upon receiving SGS_LOCATION_UPDATE_REJ              **
 **                      Based on the state, invoke state machine handlers   **
 **                                                                          **
 ** Inputs:              nas_sgs_location_update_rej                         **
 **
 *******************************************************************************/
int mme_app_handle_sgsap_location_update_rej(
    mme_app_desc_t* mme_app_desc_p,
    itti_sgsap_location_update_rej_t* const itti_sgsap_location_update_rej) {
  imsi64_t imsi64                      = INVALID_IMSI64;
  int rc                               = RETURNok;
  struct ue_mm_context_s* ue_context_p = NULL;
  sgs_fsm_t sgs_fsm;

  OAILOG_FUNC_IN(LOG_MME_APP);
  OAILOG_INFO(LOG_MME_APP, "Received SGSAP LOCATION UPDATE REJECT \n");

  /*Fetch UE context*/
  IMSI_STRING_TO_IMSI64(itti_sgsap_location_update_rej->imsi, &imsi64);
  ue_context_p =
      mme_ue_context_exists_imsi(&mme_app_desc_p->mme_ue_contexts, imsi64);
  if (ue_context_p == NULL) {
    OAILOG_ERROR(
        LOG_MME_APP,
        "Unknown IMSI in mme_app_handle_sgsap_location_update_rej\n");
    mme_ue_context_dump_coll_keys(&mme_app_desc_p->mme_ue_contexts);
    OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNerror);
  }

  if (ue_context_p->sgs_context == NULL) {
    OAILOG_ERROR(
        LOG_MME_APP,
        "SGS context not found in mme_app_handle_sgsap_location_update_rej\n");
    OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNerror);
  }
  sgs_fsm.primitive = _SGS_LOCATION_UPDATE_REJECT;
  sgs_fsm.ue_id     = ue_context_p->mme_ue_s1ap_id;
  sgs_fsm.ctx       = (void*) ue_context_p->sgs_context;
  ((sgs_context_t*) sgs_fsm.ctx)->sgsap_msg =
      (void*) itti_sgsap_location_update_rej;

  if (sgs_fsm_process(&sgs_fsm) != RETURNok) {
    OAILOG_ERROR(
        LOG_MME_APP, "Error in invoking FSM handler for primitive %d \n",
        sgs_fsm.primitive);
    OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNerror);
  }

  OAILOG_FUNC_RETURN(LOG_MME_APP, rc);
}

/******************************************************************************
 **
 ** Name:                sgs_fsm_null_loc_updt_acc()                         **
 ** Description          Handling of SGS_LOCATION UPDATE ACCEPT in NULL      **
 **                      state                                               **
 **                                                                          **
 ** Inputs:              sgs_fsm_t                                           **
 **                                                                          **
 *******************************************************************************/
int sgs_fsm_null_loc_updt_acc(const sgs_fsm_t* fsm_evt) {
  int rc                                                             = RETURNok;
  itti_sgsap_location_update_acc_t* itti_sgsap_location_update_acc_p = NULL;
  MobileIdentity_t* mobileid                                         = NULL;

  OAILOG_FUNC_IN(LOG_MME_APP);
  sgs_context_t* sgs_context = (sgs_context_t*) fsm_evt->ctx;
  itti_sgsap_location_update_acc_p =
      (itti_sgsap_location_update_acc_t*) sgs_context->sgsap_msg;

  if (sgs_context == NULL) {
    OAILOG_ERROR(
        LOG_MME_APP, "SGS Context is NULL for UE ID %d ", fsm_evt->ue_id);
    OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNerror);
  }

  if (sgs_context->ts6_1_timer.id != MME_APP_TIMER_INACTIVE_ID) {
    OAILOG_DEBUG(
        LOG_MME_APP,
        "Dropping Location Update Accept as Ts6-1 timer is running\n");
  }
  // If we received Location Updt Accept from MSC/VLR and ts6_1_timer is not
  // running
  else if (sgs_context->ts6_1_timer.id == MME_APP_TIMER_INACTIVE_ID) {
    /* If Ts8/Ts9 timer is running, drop the message
     *  If Ts8/Ts9 timer is not running and SGs state is SGs-ASSOCIATED, drop
     * the message If Ts8/Ts9 timer is not running and the SGs state is SGs-NULL
     * or LA-UPDATE_REQUESTED send SGs-Status message
     */
    OAILOG_DEBUG(
        LOG_MME_APP,
        "Received Location Update Accept when Ts6-1 timer is not running \n");
    if ((sgs_context->ts8_timer.id != MME_APP_TIMER_INACTIVE_ID) ||
        (sgs_context->ts9_timer.id != MME_APP_TIMER_INACTIVE_ID)) {
      OAILOG_DEBUG(
          LOG_MME_APP,
          "Dropping Location Update Accept as Ts8/Ts9 timer is running\n");
    } else if (
        (sgs_context->ts8_timer.id == MME_APP_TIMER_INACTIVE_ID) &&
        (sgs_context->ts9_timer.id == MME_APP_TIMER_INACTIVE_ID)) {
      OAILOG_DEBUG(LOG_MME_APP, "Send SGS-STATUS message\n");
      // Send SGS-STATUS message to SGS task
      if (itti_sgsap_location_update_acc_p->presencemask &
          SGSAP_MOBILE_IDENTITY) {
        mobileid = &itti_sgsap_location_update_acc_p->mobileid;
      }
      if (build_sgs_status(
              itti_sgsap_location_update_acc_p->imsi,
              itti_sgsap_location_update_acc_p->imsi_length,
              itti_sgsap_location_update_acc_p->laicsfb, mobileid,
              SGsAP_LOCATION_UPDATE_ACCEPT) == RETURNok) {
        OAILOG_DEBUG(LOG_MME_APP, "SGS-STATUS message sent to SGS task\n");
      }
    }
  }
  OAILOG_FUNC_RETURN(LOG_MME_APP, rc);
}

/**********************************************************************************
 **
 ** Name:                sgs_fsm_associated_loc_updt_acc() **
 ** Description          Handling of SGS_LOCATION UPDATE ACCEPT in Associated **
 **                      state **
 ** **
 ** Inputs:              sgs_fsm_t **
 ** **
 ***********************************************************************************/
int sgs_fsm_associated_loc_updt_acc(const sgs_fsm_t* fsm_evt) {
  int rc = RETURNok;

  OAILOG_FUNC_IN(LOG_MME_APP);

  sgs_context_t* sgs_context = (sgs_context_t*) fsm_evt->ctx;
  if (sgs_context == NULL) {
    OAILOG_ERROR(LOG_MME_APP, "Unknown UE ID %d ", fsm_evt->ue_id);
    OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNerror);
  }
  // If we received Location Updt Accept from MSC/VLR and ts6_1_timer is not
  // running
  if (sgs_context->ts6_1_timer.id == MME_APP_TIMER_INACTIVE_ID) {
    /* If Ts8/Ts9 timer is running, drop the message
     *  If Ts8/Ts9 timer is not running and SGs state is SGs-ASSOCIATED, drop
     * the message
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

/******************************************************************************
 **
 ** Name:                sgs_fsm_la_updt_req_loc_updt_acc()                  **
 ** Description          Handling of SGS_LOCATION UPDATE ACCEPT in LA Update **
 **                      Requested state                                     **
 **                                                                          **
 ** Inputs:              sgs_fsm_t                                           **
 **                                                                          **
 *******************************************************************************/
int sgs_fsm_la_updt_req_loc_updt_acc(const sgs_fsm_t* fsm_evt) {
  int rc                                                             = RETURNok;
  itti_sgsap_location_update_acc_t* itti_sgsap_location_update_acc_p = NULL;
  struct ue_mm_context_s* ue_context_p                               = NULL;
  MobileIdentity_t* mobileid                                         = NULL;

  OAILOG_FUNC_IN(LOG_MME_APP);
  ue_context_p = mme_ue_context_exists_mme_ue_s1ap_id(fsm_evt->ue_id);
  if (ue_context_p == NULL) {
    OAILOG_ERROR(LOG_MME_APP, "Unknown UE ID %d ", fsm_evt->ue_id);
    OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNerror);
  }
  sgs_context_t* sgs_context = (sgs_context_t*) fsm_evt->ctx;
  if (sgs_context == NULL) {
    OAILOG_ERROR(
        LOG_MME_APP, "SGS context not found for UE ID %d ", fsm_evt->ue_id);
    OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNerror);
  }
  itti_sgsap_location_update_acc_p =
      (itti_sgsap_location_update_acc_t*) sgs_context->sgsap_msg;
  if (sgs_context->ts6_1_timer.id != MME_APP_TIMER_INACTIVE_ID) {
    /*Update SGS context and set VLR Reliable to true*/
    sgs_context->sgs_state    = SGS_ASSOCIATED;
    sgs_context->vlr_reliable = true;

    /*Stop Ts6-1 timer*/
    mme_app_stop_timer(ue_context_p->sgs_context->ts6_1_timer.id);

    sgs_context->ts6_1_timer.id = MME_APP_TIMER_INACTIVE_ID;
    if ((handle_cs_domain_loc_updt_acc(
            itti_sgsap_location_update_acc_p, ue_context_p)) == RETURNerror) {
      OAILOG_DEBUG(
          LOG_MME_APP,
          "Failed to update CSFB params received from MSC/VLR for "
          "UE " IMSI_64_FMT "\n",
          ue_context_p->emm_context._imsi64);
      rc = RETURNerror;
    }
  } else if (sgs_context->ts6_1_timer.id == MME_APP_TIMER_INACTIVE_ID) {
    /* If we received Location Updt Accept from MSC/VLR and
     * ts6_1_timer is not running
     */

    /* If Ts8/Ts9 timer is running, drop the message
     *  If Ts8/Ts9 timer is not running and SGs state is SGs-ASSOCIATED, drop
     * the message If Ts8/Ts9 timer is not running and the SGs state is SGs-NULL
     * or LA-UPDATE_REQUESTED send SGs-Status message
     */
    OAILOG_DEBUG(
        LOG_MME_APP,
        "Received Location Update Accept when Ts6-1 timer is not running \n");
    if ((sgs_context->ts8_timer.id != MME_APP_TIMER_INACTIVE_ID) ||
        (sgs_context->ts9_timer.id != MME_APP_TIMER_INACTIVE_ID)) {
      OAILOG_DEBUG(
          LOG_MME_APP,
          "Dropping Location Update Accept as Ts8/Ts9 timer is running\n");
    } else if (
        (sgs_context->ts8_timer.id == MME_APP_TIMER_INACTIVE_ID) &&
        (sgs_context->ts8_timer.id == MME_APP_TIMER_INACTIVE_ID)) {
      OAILOG_DEBUG(LOG_MME_APP, "Send SGS-STATUS message\n");
      if (itti_sgsap_location_update_acc_p->presencemask &
          SGSAP_MOBILE_IDENTITY) {
        mobileid = &itti_sgsap_location_update_acc_p->mobileid;
      }
      // Send SGS-STATUS message to SGS task
      if (build_sgs_status(
              itti_sgsap_location_update_acc_p->imsi,
              itti_sgsap_location_update_acc_p->imsi_length,
              itti_sgsap_location_update_acc_p->laicsfb, mobileid,
              SGsAP_LOCATION_UPDATE_ACCEPT) == RETURNok) {
        OAILOG_DEBUG(LOG_MME_APP, "SGS-STATUS message sent to SGS task\n");
      }
    }
  }
  OAILOG_FUNC_RETURN(LOG_MME_APP, rc);
}

/**********************************************************************************
 **
 ** Name:                sgs_fsm_null_loc_updt_rej() **
 ** Description          Handling of SGS_LOCATION UPDATE REJECT in NULL **
 **                      state **
 ** **
 ** Inputs:              sgs_fsm_t **
 ** **
 ***********************************************************************************/
int sgs_fsm_null_loc_updt_rej(const sgs_fsm_t* fsm_evt) {
  int rc = RETURNok;
  OAILOG_FUNC_IN(LOG_MME_APP);
  OAILOG_ERROR(
      LOG_MME_APP, "Dropping message as it is received in NULL state for UE %d",
      fsm_evt->ue_id);

  OAILOG_FUNC_RETURN(LOG_MME_APP, rc);
}

/****************************************************************************************
 **
 ** Name:                sgs_fsm_la_updt_req_loc_updt_rej **
 ** Description          Handling of SGS_LOCATION UPDATE REJECT in LA UPDT
 *REQUESTED   **
 **                      state **
 ** **
 ** Inputs:              sgs_fsm_t **
 ** **
 *****************************************************************************************/
int sgs_fsm_la_updt_req_loc_updt_rej(const sgs_fsm_t* fsm_evt) {
  int rc                                                             = RETURNok;
  struct ue_mm_context_s* ue_context_p                               = NULL;
  itti_sgsap_location_update_rej_t* itti_sgsap_location_update_rej_p = NULL;
  imsi64_t imsi64 = INVALID_IMSI64;
  lai_t* lai      = NULL;

  OAILOG_FUNC_IN(LOG_MME_APP);

  sgs_context_t* sgs_context = (sgs_context_t*) fsm_evt->ctx;
  itti_sgsap_location_update_rej_p =
      (itti_sgsap_location_update_rej_t*) sgs_context->sgsap_msg;
  IMSI_STRING_TO_IMSI64(itti_sgsap_location_update_rej_p->imsi, &imsi64);
  ue_context_p = mme_ue_context_exists_mme_ue_s1ap_id(fsm_evt->ue_id);
  if (ue_context_p == NULL) {
    mme_app_desc_t* mme_app_desc_p = get_mme_nas_state(false);
    OAILOG_ERROR(LOG_MME_APP, "Unknown UE ID %d ", fsm_evt->ue_id);
    mme_ue_context_dump_coll_keys(&mme_app_desc_p->mme_ue_contexts);
    OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNerror);
  }

  if (sgs_context == NULL) {
    OAILOG_ERROR(
        LOG_MME_APP, "SGS Context is NULL for UE ID %d ", fsm_evt->ue_id);
    OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNerror);
  }
  // Change SGS state to NULL
  sgs_context->sgs_state = SGS_NULL;
  // Stop Ts6-1 timer
  if (sgs_context->ts6_1_timer.id != MME_APP_TIMER_INACTIVE_ID) {
    mme_app_stop_timer(ue_context_p->sgs_context->ts6_1_timer.id);
    sgs_context->ts6_1_timer.id = MME_APP_TIMER_INACTIVE_ID;
  }

  if (itti_sgsap_location_update_rej_p->presencemask & SGSAP_LAI) {
    lai = &itti_sgsap_location_update_rej_p->laicsfb;
  }
  // Handle SGS Location Update Failure
  nas_proc_cs_domain_location_updt_fail(
      itti_sgsap_location_update_rej_p->cause, lai,
      ue_context_p->mme_ue_s1ap_id);

  OAILOG_FUNC_RETURN(LOG_MME_APP, rc);
}

/**********************************************************************************
 **
 ** Name:                mme_app_handle_ts6_1_timer_expiry() **
 ** Description          Ts6_1 timer expiry handler **
 ** **
 ** Inputs:              ue_mm_context_s **
 ** **
 ***********************************************************************************/
int mme_app_handle_ts6_1_timer_expiry(zloop_t* loop, int timer_id, void* args) {
  OAILOG_FUNC_IN(LOG_MME_APP);
  mme_ue_s1ap_id_t mme_ue_s1ap_id = 0;
  if (!mme_app_get_timer_arg(timer_id, &mme_ue_s1ap_id)) {
    OAILOG_WARNING(
        LOG_MME_APP, "Invalid Timer Id expiration, Timer Id: %u\n", timer_id);
    OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNerror);
  }
  struct ue_mm_context_s* ue_context_p =
      mme_app_get_ue_context_for_timer(mme_ue_s1ap_id, "sgs ts6_1 timer");
  if (ue_context_p == NULL) {
    OAILOG_ERROR(
        LOG_MME_APP,
        "Invalid UE context received, MME UE S1AP Id: " MME_UE_S1AP_ID_FMT "\n",
        mme_ue_s1ap_id);
    OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNerror);
  }
  if (ue_context_p->sgs_context == NULL) {
    OAILOG_ERROR(
        LOG_MME_APP,
        "Ts6-1  Timer expired, but sgs context is NULL for "
        "ue-id " MME_UE_S1AP_ID_FMT "\n",
        mme_ue_s1ap_id);
    OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNerror);
  }

  ue_context_p->sgs_context->ts6_1_timer.id = MME_APP_TIMER_INACTIVE_ID;
  ue_context_p->sgs_context->sgs_state      = SGS_NULL;

  // Handle SGS Location Update Failure
  nas_proc_cs_domain_location_updt_fail(
      SGS_MSC_NOT_REACHABLE, NULL, mme_ue_s1ap_id);

  OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNok);
}

/**********************************************************************************
 **
 ** Name:                sgs_fsm_associated_loc_updt_rej() **
 ** Description          Handling of SGS_LOCATION UPDATE REJ in Associated **
 **                      state **
 ** **
 ** Inputs:              sgs_fsm_t **
 ** **
 ***********************************************************************************/
int sgs_fsm_associated_loc_updt_rej(const sgs_fsm_t* fsm_evt) {
  int rc                                                             = RETURNok;
  struct ue_mm_context_s* ue_context_p                               = NULL;
  itti_sgsap_location_update_rej_t* itti_sgsap_location_update_rej_p = NULL;
  imsi64_t imsi64 = INVALID_IMSI64;

  OAILOG_FUNC_IN(LOG_MME_APP);
  sgs_context_t* sgs_context = (sgs_context_t*) fsm_evt->ctx;
  itti_sgsap_location_update_rej_p =
      (itti_sgsap_location_update_rej_t*) sgs_context->sgsap_msg;
  IMSI_STRING_TO_IMSI64(itti_sgsap_location_update_rej_p->imsi, &imsi64);
  ue_context_p = mme_ue_context_exists_mme_ue_s1ap_id(fsm_evt->ue_id);
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
/*******************************************************************************
 **                                                                           **
 ** Name:                mme_app_create_sgs_context()                         **
 ** Description          create sgs context                                   **
 **                                                                           **
 ** Inputs:              Pointer to UE context                               **
 ** Returns:             RETURNok: On successfull sgs context creation        **
 **                      RETURNerror: On failure                              **
 **                                                                           **
 ********************************************************************************/
int mme_app_create_sgs_context(ue_mm_context_t* ue_context_p) {
  OAILOG_FUNC_IN(LOG_MME_APP);
  if (ue_context_p == NULL) {
    OAILOG_ERROR(LOG_MME_APP, "Invalid UE context \n");
    OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNerror);
  }

  ue_context_p->sgs_context = calloc(1, sizeof(sgs_context_t));
  if (!ue_context_p->sgs_context) {
    OAILOG_CRITICAL(
        LOG_MME_APP,
        "Cannot create SGS Context for UE-ID " MME_UE_S1AP_ID_FMT " \n",
        ue_context_p->mme_ue_s1ap_id);
    OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNerror);
  }
  // Initialize SGS context to default values
  ue_context_p->sgs_context->sgs_state       = SGS_NULL;
  ue_context_p->sgs_context->vlr_reliable    = false;
  ue_context_p->sgs_context->neaf            = false;
  ue_context_p->sgs_context->ts6_1_timer.id  = MME_APP_TIMER_INACTIVE_ID;
  ue_context_p->sgs_context->ts6_1_timer.sec = mme_config.sgs_config.ts6_1_sec;
  ue_context_p->sgs_context->ts8_timer.id    = MME_APP_TIMER_INACTIVE_ID;
  ue_context_p->sgs_context->ts8_timer.sec   = mme_config.sgs_config.ts8_sec;
  ue_context_p->sgs_context->ts8_retransmission_count = 0;
  ue_context_p->sgs_context->ts9_timer.id  = MME_APP_TIMER_INACTIVE_ID;
  ue_context_p->sgs_context->ts9_timer.sec = mme_config.sgs_config.ts9_sec;
  ue_context_p->sgs_context->ts9_retransmission_count = 0;
  ue_context_p->sgs_context->ts10_timer.id  = MME_APP_TIMER_INACTIVE_ID;
  ue_context_p->sgs_context->ts10_timer.sec = mme_config.sgs_config.ts10_sec;
  ue_context_p->sgs_context->ts10_retransmission_count = 0;
  ue_context_p->sgs_context->ts13_timer.id  = MME_APP_TIMER_INACTIVE_ID;
  ue_context_p->sgs_context->ts13_timer.sec = mme_config.sgs_config.ts13_sec;
  ue_context_p->sgs_context->ts13_retransmission_count = 0;
  ue_context_p->sgs_context->call_cancelled            = false;
  OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNok);
}

/******************************************************************************
 **                                                                          **
 ** Name:                map_sgs_emm_cause()                                **
 ** Description          Maps SGS Reject cause to EMM cause                  **
 **                                                                          **
 ** Inputs:              SGS Reject Cause                                    **
 ** Outputs:                                                                 **
 **      Return:    Selected emm cause                                       **
 *******************************************************************************/
int map_sgs_emm_cause(SgsRejectCause_t sgs_cause) {
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
    // TODO : Need to map appropriate cause
    case SGS_RETRY_UPON_ENTRY_INTO_NEW_CELL: {
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
      OAILOG_INFO(LOG_MME_APP, "Invalid SGS Reject cause\n");
      emm_cause = EMM_CAUSE_CS_DOMAIN_NOT_AVAILABLE;
  }
  return emm_cause;
}

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
 *-----------------------------------------------------------------------------
 * For more information about the OpenAirInterface (OAI) Software Alliance:
 *      contact@openairinterface.org
 */

/*! \file mme_app_sgs_status.c
   \brief Handles  SGSAP Status message
   \author
   \version
   \company
   \email:
*/

#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#include "conversions.h"
#include "log.h"
#include "service303.h"
#include "intertask_interface.h"
#include "mme_app_defs.h"
#include "mme_app_ue_context.h"
#include "3gpp_29.018.h"
#include "mme_app_itti_messaging.h"
#include "nas_proc.h"
#include "mme_app_timer.h"

static void mme_app_handle_sgs_status_for_imsi_detach_ind(
    ue_mm_context_t* ue_context_p);

static void mme_app_handle_sgs_status_for_eps_detach_ind(
    ue_mm_context_t* ue_context_p);

static void mme_app_handle_sgs_status_for_loc_upd_req(
    ue_mm_context_t* ue_context_p);

/****************************************************************************
 **                                                                        **
 ** Name: mme_app_handle_sgs_status_message()                              **
 **                                                                        **
 ** Description: Processes the SGSAP Status message received               **
 **              from the SGS task                                         **
 **                                                                        **
 ** Inputs: itti_sgsap_status_t: SGSAP Status message                      **
 **                                                                        **
 ** Outputs:                                                               **
 **      Return:    RETURNok, RETURNerror                                  **
 **                                                                        **
 ***************************************************************************/

int mme_app_handle_sgs_status_message(
    mme_app_desc_t* mme_app_desc_p,
    itti_sgsap_status_t* const sgsap_status_pP) {
  struct ue_mm_context_s* ue_context_p = NULL;
  uint8_t message_type;
  imsi64_t imsi64 = INVALID_IMSI64;

  OAILOG_FUNC_IN(LOG_MME_APP);
  if (!sgsap_status_pP) {
    OAILOG_ERROR(
        LOG_MME_APP, "Received invalid sgsap status message :%p \n",
        sgsap_status_pP);
    OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNerror);
  }
  IMSI_STRING_TO_IMSI64(sgsap_status_pP->imsi, &imsi64);
  OAILOG_INFO(
      LOG_MME_APP, "Received SGS-Status message for IMSI " IMSI_64_FMT "\n",
      imsi64);

  if ((ue_context_p = mme_ue_context_exists_imsi(
           &mme_app_desc_p->mme_ue_contexts, imsi64)) == NULL) {
    OAILOG_ERROR(
        LOG_MME_APP,
        "SGS-Status message: Failed to find UE context"
        " for IMSI " IMSI_64_FMT "\n",
        imsi64);
    increment_counter("sgsap_status", 1, 1, "cause", "imsi_unknown");
    OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNerror);
  }
  if (ue_context_p->sgs_context == NULL) {
    OAILOG_ERROR(
        LOG_MME_APP, "SGS context not created for IMSI " IMSI_64_FMT "\n",
        imsi64);
    increment_counter("sgsap_status", 1, 1, "cause", "SGS context not created");
    OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNerror);
  }
#define MESSAGE_TYPE_POSITION 0
  message_type = sgsap_status_pP->error_msg_rcvd[MESSAGE_TYPE_POSITION];
  switch (message_type) {
    case SGS_LOCATION_UPDATE_REQUEST: {
      OAILOG_ERROR(
          LOG_MME_APP,
          "Received SGS Error Status message for"
          "LOCATION_UPDATE_REQUEST \n");
      mme_app_handle_sgs_status_for_loc_upd_req(ue_context_p);
      break;
    }
    case SGS_TMSI_REALLOCATION_COMPLETE: {
      OAILOG_ERROR(
          LOG_MME_APP,
          "Received SGS Error Status message for"
          "TMSI_REALLOCATION_COMPLETE \n");
      break;
    }
    case SGS_EPS_DETACH_INDICATION: {
      OAILOG_ERROR(
          LOG_MME_APP,
          "Received SGS Error Status message for"
          "EPS_DETACH_INDICATION \n");
      mme_app_handle_sgs_status_for_eps_detach_ind(ue_context_p);
      break;
    }
    case SGS_IMSI_DETACH_INDICATION: {
      OAILOG_ERROR(
          LOG_MME_APP,
          "Received SGS Error Status message for"
          "IMSI_DETACH_INDICATION \n");
      mme_app_handle_sgs_status_for_imsi_detach_ind(ue_context_p);
      break;
    }
    case SGS_RESET_ACK: {
      OAILOG_ERROR(
          LOG_MME_APP,
          "Received SGS Error Status message for"
          "RESET_ACK \n");
      break;
    }
    case SGS_UE_ACTIVITY_INDICATION: {
      OAILOG_ERROR(
          LOG_MME_APP,
          "Received SGS Error Status message for"
          "UE_ACTIVITY_INDICATION \n");
      ue_context_p->sgs_context->neaf = SET_NEAF;
      break;
    }
    case SGS_ALERT_ACK: {
      OAILOG_ERROR(
          LOG_MME_APP,
          "Received SGS Error Status message for"
          "ALERT_ACK \n");
      ue_context_p->sgs_context->neaf = SET_NEAF;
      break;
    }
    case SGS_ALERT_REJECT: {
      OAILOG_ERROR(
          LOG_MME_APP,
          "Received SGS Error Status message for"
          "ALERT_REJECT \n");
      break;
    }
    case SGS_SERVICE_REQUEST: {
      if (ue_context_p->sgs_context->csfb_service_type ==
          CSFB_SERVICE_MT_CALL) {
        OAILOG_ERROR(
            LOG_MME_APP,
            "Received SGS Error Status message for"
            "SERVICE_REQUEST for MT CS call \n");
        if (!ue_context_p->sgs_context->mt_call_in_progress) {
          ue_context_p->sgs_context->call_cancelled = true;
          OAILOG_DEBUG(LOG_MME_APP, "Setting Call Cancelled flag to true \n");
        } else {
          OAILOG_INFO(
              LOG_MME_APP,
              "Can not abort MT call, as MT call is"
              "already accepted by the user\n");
        }
      } else if (
          ue_context_p->sgs_context->csfb_service_type == CSFB_SERVICE_MT_SMS) {
        OAILOG_ERROR(
            LOG_MME_APP,
            "Received SGS Error Status message for"
            "SERVICE_REQUEST sent for SMS services \n");
      }
      break;
    }
    case SGS_UE_UNREACHABLE: {
      OAILOG_ERROR(
          LOG_MME_APP,
          "Received SGS Error Status message for"
          "UE_UN_REACHABLE \n");
      break;
    }
    case SGS_UPLINK_UNIT_DATA: {
      OAILOG_ERROR(
          LOG_MME_APP,
          "Received SGS Error Status message for"
          "SGS_UPLINK_UNIT_DATA \n");
      break;
    }
    default: {
      OAILOG_ERROR(
          LOG_MME_APP,
          "Received unknown messag type in SGS Error"
          "Status message, the received message is :%d \n",
          message_type);
      break;
    }
  }
  OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNok);
}

/****************************************************************************
 **                                                                        **
 ** name: _mme_app_handle_sgs_status_for_imsi_detach_ind()                 **
 **                                                                        **
 ** description: processes the sgsap status message received               **
 **              for IMSI Detach Ind from the sgs task                     **
 **                                                                        **
 ** inputs: ue_mm_context_t : pointer to UE's MM context                   **
 **                                                                        **
 ** outputs:                                                               **
 **      return: None                                                      **
 **                                                                        **
 ***************************************************************************/

static void mme_app_handle_sgs_status_for_imsi_detach_ind(
    ue_mm_context_t* ue_context_p) {
  OAILOG_FUNC_IN(LOG_MME_APP);

  if (ue_context_p->sgs_context) {
    // Send the S1AP NAS DL DATA REQ in case of IMSI or combined EPS/IMSI detach
    if ((ue_context_p->sgs_detach_type ==
         SGS_EXPLICIT_UE_INITIATED_IMSI_DETACH_FROM_NONEPS) ||
        (ue_context_p->sgs_detach_type ==
         SGS_COMBINED_UE_INITIATED_IMSI_DETACH_FROM_EPS_N_NONEPS)) {
      send_msg_to_task(
          &mme_app_task_zmq_ctx, TASK_S1AP,
          ue_context_p->sgs_context->message_p);
      ue_context_p->sgs_context->message_p = NULL;
      /*
       Notify S1AP to send UE Context Release Command to eNB or
       free s1 context locally,if the ue requested for combined EPS/IMSI detach
       if the ue is in idle state and requested for IMSI detach
      */
      if ((ue_context_p->sgs_detach_type ==
           SGS_COMBINED_UE_INITIATED_IMSI_DETACH_FROM_EPS_N_NONEPS) ||
          ((ue_context_p->sgs_detach_type ==
            SGS_EXPLICIT_UE_INITIATED_IMSI_DETACH_FROM_NONEPS) &&
           (ue_context_p->ue_context_rel_cause ==
            S1AP_RADIO_EUTRAN_GENERATED_REASON))) {
        mme_app_itti_ue_context_release(
            ue_context_p, ue_context_p->ue_context_rel_cause);
        ue_context_p->ue_context_rel_cause = S1AP_INVALID_CAUSE;
      }
    }
    // Free the UE SGS context
    mme_app_ue_sgs_context_free_content(
        ue_context_p->sgs_context, ue_context_p->emm_context._imsi64);
  }
  OAILOG_FUNC_OUT(LOG_MME_APP);
}

/****************************************************************************
 **                                                                        **
 ** name: _mme_app_handle_sgs_status_for_eps_detach_ind()                  **
 **                                                                        **
 ** description: processes the sgsap status message received               **
 **              for EPS Detach Ind from the sgs task                      **
 **                                                                        **
 ** inputs: ue_mm_context_t : pointer to UE's MM context                   **
 **                                                                        **
 ** outputs:                                                               **
 **      return: None                                                      **
 **                                                                        **
 ***************************************************************************/

static void mme_app_handle_sgs_status_for_eps_detach_ind(
    ue_mm_context_t* ue_context_p) {
  OAILOG_FUNC_IN(LOG_MME_APP);

  if (ue_context_p->sgs_context) {
    /* Stop SGS EPS Detach timer, after recieving the SGS Status message for
       EPS Detach Ind,if running
    */
    if (ue_context_p->sgs_context->ts8_timer.id != MME_APP_TIMER_INACTIVE_ID) {
      mme_app_stop_timer(ue_context_p->sgs_context->ts8_timer.id);
      ue_context_p->sgs_context->ts8_timer.id = MME_APP_TIMER_INACTIVE_ID;
    }
  }
  OAILOG_FUNC_OUT(LOG_MME_APP);
}

/****************************************************************************
 **                                                                        **
 ** name: _mme_app_handle_sgs_status_for_loc_upd_req()                     **
 **                                                                        **
 ** description: processes the sgsap status message received               **
 **              for Update Location Request from the sgs task             **
 **                                                                        **
 ** inputs: ue_mm_context_t : pointer to UE's MM context                   **
 **                                                                        **
 ** outputs:                                                               **
 **      return: None                                                      **
 **                                                                        **
 ***************************************************************************/

static void mme_app_handle_sgs_status_for_loc_upd_req(
    ue_mm_context_t* ue_context_p) {
  OAILOG_FUNC_IN(LOG_MME_APP);
  lai_t* lai = NULL;

  mme_app_ue_sgs_context_free_content(
      ue_context_p->sgs_context, ue_context_p->emm_context._imsi64);
  nas_proc_cs_domain_location_updt_fail(
      SGS_PROTOCOL_ERROR_UNSPECIFIED, lai, ue_context_p->mme_ue_s1ap_id);
  OAILOG_FUNC_OUT(LOG_MME_APP);
}

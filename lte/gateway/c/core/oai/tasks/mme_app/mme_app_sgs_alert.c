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

/*****************************************************************************

  Source      mme_app_sgs_alert.c

  Version

  Date

  Product    MME app

  Subsystem  SGS (an interface between MME and MSC/VLR) message handling for
             non-eps alert procedure

  Author

  Description Handles non-eps procedure

*****************************************************************************/

#include <stdbool.h>
#include <stdint.h>
#include <stdlib.h>
#include <string.h>

#include "log.h"
#include "service303.h"
#include "mme_config.h"
#include "intertask_interface.h"
#include "mme_app_sgs_fsm.h"
#include "conversions.h"
#include "mme_app_defs.h"
#include "common_defs.h"
#include "common_types.h"
#include "intertask_interface_types.h"
#include "itti_types.h"
#include "mme_api.h"
#include "mme_app_desc.h"
#include "mme_app_ue_context.h"
#include "sgs_messages_types.h"

static int mme_app_send_sgsap_alert_reject(
    itti_sgsap_alert_request_t* const sgsap_alert_req_pP, SgsCause_t sgs_cause,
    uint64_t imsi64);

static int mme_app_send_sgsap_alert_ack(
    itti_sgsap_alert_request_t* const sgsap_alert_req_pP, uint64_t imsi64);

/****************************************************************************
 **                                                                        **
 ** Name:    mme_app_handle_sgsap_alert_request()                          **
 **                                                                        **
 ** Description: Processes the SGSAP alert Request message                 **
 **      received from the SGS task                                        **
 **                                                                        **
 ** Inputs:  itti_sgsap_alert_request_t: SGSAP alert Request message       **
 **                                                                        **
 ** Outputs:                                                               **
 **      Return:    RETURNok, RETURNerror                                  **
 **                                                                        **
 ***************************************************************************/
int mme_app_handle_sgsap_alert_request(
    mme_app_desc_t* mme_app_desc_p,
    itti_sgsap_alert_request_t* const sgsap_alert_req_pP) {
  uint64_t imsi64                      = 0;
  struct ue_mm_context_s* ue_context_p = NULL;

  OAILOG_FUNC_IN(LOG_MME_APP);
  if (!sgsap_alert_req_pP) {
    OAILOG_ERROR(
        LOG_MME_APP, "Invalid SGSAP Alert Request ITTI message received\n");
    OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNerror);
  }

  IMSI_STRING_TO_IMSI64(sgsap_alert_req_pP->imsi, &imsi64);

  OAILOG_INFO(
      LOG_MME_APP, "Received SGS-ALERT REQUEST for IMSI " IMSI_64_FMT "\n",
      imsi64);
  if ((ue_context_p = mme_ue_context_exists_imsi(
           &mme_app_desc_p->mme_ue_contexts, imsi64)) == NULL) {
    OAILOG_ERROR(
        LOG_MME_APP,
        "SGS-ALERT REQUEST: Failed to find UE context for IMSI " IMSI_64_FMT
        "\n",
        imsi64);
    mme_app_send_sgsap_alert_reject(
        sgsap_alert_req_pP, SGS_CAUSE_IMSI_UNKNOWN, imsi64);
    increment_counter("sgsap_alert_reject", 1, 1, "cause", "imsi_unknown");
    OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNerror);
  }
  if (ue_context_p->mm_state == UE_UNREGISTERED) {
    OAILOG_INFO(
        LOG_MME_APP,
        "SGS-ALERT REQUEST: UE is currently not attached to EPS service and "
        "send Alert Reject to MSC/VLR for UE:" IMSI_64_FMT " \n",
        imsi64);
    mme_app_send_sgsap_alert_reject(
        sgsap_alert_req_pP, SGS_CAUSE_IMSI_DETACHED_FOR_EPS_SERVICE, imsi64);
    increment_counter(
        "sgsap_alert_reject", 1, 1, "cause", "ue_is_not_registered_to_eps");
    OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNerror);
  }
  if (ue_context_p->sgs_context == NULL) {
    if ((mme_app_create_sgs_context(ue_context_p)) != RETURNok) {
      OAILOG_CRITICAL(
          LOG_MME_APP,
          "Failed to create SGS context for ue_id " MME_UE_S1AP_ID_FMT "\n",
          ue_context_p->mme_ue_s1ap_id);
      OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNerror);
    }
  }
  ue_context_p->sgs_context->neaf = SET_NEAF;
  /* send Alert Ack */
  mme_app_send_sgsap_alert_ack(sgsap_alert_req_pP, imsi64);
  OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNok);
}

/**********************************************************************************
 ** **
 ** Name:    _mme_app_send_sgsap_alert_reject() **
 ** Description   Build and send Alert reject **
 ** Inputs: **
 **          sgsap_alert_req_pP: Received Alert Request message **
 **          sgs_cause         : alert reject cause **
 **          imsi              : imsi **
 ** Outputs: **
 **          Return:    RETURNok, RETURNerror **
 **
 ***********************************************************************************/
static int mme_app_send_sgsap_alert_reject(
    itti_sgsap_alert_request_t* const sgsap_alert_req_pP, SgsCause_t sgs_cause,
    uint64_t imsi64) {
  OAILOG_FUNC_IN(LOG_MME_APP);
  int rc                                           = RETURNerror;
  MessageDef* message_p                            = NULL;
  itti_sgsap_alert_reject_t* sgsap_alert_reject_pP = NULL;

  message_p = itti_alloc_new_message(TASK_MME_APP, SGSAP_ALERT_REJECT);
  if (!message_p) {
    OAILOG_ERROR(
        LOG_MME_APP,
        "Failed to allocate memory: SGSAP_ALERT_REJECT, IMSI: " IMSI_64_FMT
        "\n",
        imsi64);
    OAILOG_FUNC_RETURN(LOG_MME_APP, rc);
  }
  sgsap_alert_reject_pP = &message_p->ittiMsg.sgsap_alert_reject;
  memset((void*) sgsap_alert_reject_pP, 0, sizeof(itti_sgsap_alert_reject_t));

  memcpy(
      (void*) sgsap_alert_reject_pP->imsi,
      (const void*) sgsap_alert_req_pP->imsi, sgsap_alert_req_pP->imsi_length);
  sgsap_alert_reject_pP->imsi_length = sgsap_alert_req_pP->imsi_length;
  sgsap_alert_reject_pP->sgs_cause   = sgs_cause;

  OAILOG_INFO(
      LOG_MME_APP,
      "Send SGSAP-Alert Reject for IMSI" IMSI_64_FMT " with sgs-cause :%d \n",
      imsi64, (int) sgs_cause);
  rc = send_msg_to_task(&mme_app_task_zmq_ctx, TASK_SGS, message_p);
  OAILOG_FUNC_RETURN(LOG_MME_APP, rc);
}

/**********************************************************************************
 ** **
 ** Name:    _mme_app_send_sgsap_alert_ack() **
 ** Description   Build and send Alert ack **
 ** Inputs: **
 **          sgsap_alert_req_pP: Received Alert ack message **
 **          imsi64        : imsi **
 ** Outputs: **
 **          Return:    RETURNok, RETURNerror **
 **
 ***********************************************************************************/
static int mme_app_send_sgsap_alert_ack(
    itti_sgsap_alert_request_t* const sgsap_alert_req_pP, uint64_t imsi64) {
  int rc                                     = RETURNerror;
  MessageDef* message_p                      = NULL;
  itti_sgsap_alert_ack_t* sgsap_alert_ack_pP = NULL;
  OAILOG_FUNC_IN(LOG_MME_APP);

  message_p = itti_alloc_new_message(TASK_MME_APP, SGSAP_ALERT_ACK);
  if (!message_p) {
    OAILOG_ERROR(
        LOG_MME_APP,
        "Failed to allocate memory for SGSAP_ALERT_ACK, IMSI: " IMSI_64_FMT
        "\n",
        imsi64);
    OAILOG_FUNC_RETURN(LOG_MME_APP, rc);
  }
  sgsap_alert_ack_pP = &message_p->ittiMsg.sgsap_alert_ack;
  memset((void*) sgsap_alert_ack_pP, 0, sizeof(itti_sgsap_alert_ack_t));

  memcpy(
      (void*) sgsap_alert_ack_pP->imsi, (const void*) sgsap_alert_req_pP->imsi,
      sgsap_alert_req_pP->imsi_length);
  sgsap_alert_ack_pP->imsi_length = sgsap_alert_req_pP->imsi_length;

  OAILOG_INFO(
      LOG_MME_APP, "Send SGSAP-Alert Reject for IMSI" IMSI_64_FMT " \n",
      imsi64);
  rc = send_msg_to_task(&mme_app_task_zmq_ctx, TASK_SGS, message_p);
  OAILOG_FUNC_RETURN(LOG_MME_APP, rc);
}

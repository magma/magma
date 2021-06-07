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

/*! \file mme_app_sgsap_service_abort.c
 */

#include <stdio.h>
#include <stdbool.h>

#include "common_types.h"
#include "conversions.h"
#include "log.h"
#include "mme_app_ue_context.h"
#include "mme_app_defs.h"
#include "mme_app_sgs_fsm.h"
#include "common_defs.h"
#include "mme_app_desc.h"
#include "sgs_messages_types.h"

/**********************************************************************************
 ** **
 ** Name:                mme_app_handle_sgsap_service_abort_request() **
 ** Description          Handling of SGS_SERVICE_ABORT_REQUEST **
 **                      at MME App, invoke FSM **
 ** **
 ** Inputs:              itti_sgsap_service_abort_req_p - Received SGS Service
 ***
 **                      Abort Request message **
 ** Outputs:             Return:    RETURNok, RETURNerror **
 ***********************************************************************************/
int mme_app_handle_sgsap_service_abort_request(
    mme_app_desc_t* mme_app_desc_p,
    itti_sgsap_service_abort_req_t* const itti_sgsap_service_abort_req_p) {
  imsi64_t imsi64                      = INVALID_IMSI64;
  struct ue_mm_context_s* ue_context_p = NULL;
  sgs_fsm_t sgs_fsm;

  OAILOG_FUNC_IN(LOG_MME_APP);
  OAILOG_INFO(
      LOG_MME_APP, "Received SGSAP SERVICE ABORT REQUEST with IMSI %s\n",
      itti_sgsap_service_abort_req_p->imsi);

  if (NULL == itti_sgsap_service_abort_req_p) {
    OAILOG_ERROR(
        LOG_MME_APP,
        "Received itti_sgsap_service_abort_req_p is NULL in "
        "mme_app_handle_sgsap_service_abort_req\n");
    OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNerror);
  }
  IMSI_STRING_TO_IMSI64(itti_sgsap_service_abort_req_p->imsi, &imsi64);
  ue_context_p =
      mme_ue_context_exists_imsi(&mme_app_desc_p->mme_ue_contexts, imsi64);
  if (ue_context_p == NULL) {
    OAILOG_ERROR(
        LOG_MME_APP,
        "Unknown IMSI in mme_app_handle_sgsap_service_abort_req %s\n",
        itti_sgsap_service_abort_req_p->imsi);
    OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNerror);
  }
  if (ue_context_p->sgs_context == NULL) {
    OAILOG_ERROR(
        LOG_MME_APP,
        "SGS context not found in mme_app_handle_sgsap_service_abort_req for "
        "IMSI %s\n",
        itti_sgsap_service_abort_req_p->imsi);
    OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNerror);
  }

  sgs_fsm.primitive = _SGS_SERVICE_ABORT_REQUEST;
  sgs_fsm.ue_id     = ue_context_p->mme_ue_s1ap_id;
  sgs_fsm.ctx       = ue_context_p->sgs_context;
  ((sgs_context_t*) sgs_fsm.ctx)->sgsap_msg =
      (void*) itti_sgsap_service_abort_req_p;

  if (sgs_fsm_process(&sgs_fsm) != RETURNok) {
    OAILOG_ERROR(
        LOG_MME_APP, "Error in invoking FSM handler for primitive %d \n",
        sgs_fsm.primitive);
    OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNerror);
  }

  OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNok);
}

/**********************************************************************************
 **
 ** Name:                sgs_fsm_associated_service_abort_request()        **
 ** Description          Upon receiving SGS_SERVICE_ABORT_REQUEST **
 **                      set Call cancelled flag to true **
 ** **
 ** Inputs:              sgs_fsm_t **
 ** Outputs:             Return:    RETURNok, RETURNerror **
 ***********************************************************************************/
int sgs_fsm_associated_service_abort_request(const sgs_fsm_t* fsm_evt) {
  OAILOG_FUNC_IN(LOG_MME_APP);
  sgs_context_t* sgs_context = (sgs_context_t*) fsm_evt->ctx;

  if (sgs_context == NULL) {
    OAILOG_ERROR(
        LOG_MME_APP, "SGS context not found for UE ID %d \n", fsm_evt->ue_id);
    OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNerror);
  }
  /* As per 29.118 ,if we receive EXT Service Request with csfb_response set to
   * call_accepted drop SERVICE_ABORT message and proceed wit MT call
   */
  if (!sgs_context->mt_call_in_progress) {
    sgs_context->call_cancelled = true;
    OAILOG_DEBUG(LOG_MME_APP, "Setting Call Cancelled flag to true\n");
  } else {
    OAILOG_INFO(
        LOG_MME_APP,
        "Dropping SGS SERVICE ABORT REQUEST message as MT call is already "
        "accepted by the user for UE ID %d\n",
        fsm_evt->ue_id);
  }

  OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNok);
}

/**********************************************************************************
 **
 ** Name:                sgs_fsm_null_service_abort_request() **
 ** Description          Upon receiving SGS_SERVICE_ABORT_REQUEST in null state
 ***
 **                      log error **
 ** Inputs:              sgs_fsm_t **
 ** Outputs:             Return:    RETURNok, RETURNerror **
 ***********************************************************************************/
int sgs_fsm_null_service_abort_request(const sgs_fsm_t* fsm_evt) {
  OAILOG_FUNC_IN(LOG_MME_APP);
  OAILOG_ERROR(
      LOG_MME_APP,
      "Dropping SGS SERVICE ABORT REQUEST message as it is received in NULL "
      "state for UE ID %d\n",
      fsm_evt->ue_id);

  OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNok);
}

/**********************************************************************************
 ** **
 ** Name:                sgs_fsm_la_update_req_service_abort_request() **
 ** Description          Upon receiving SGS_SERVICE_ABORT_REQUEST in **
                         LA_UPDATE_REQUESTED state log error **
 ** Inputs:              sgs_fsm_t **
 ** Outputs:             Return:    RETURNok, RETURNerror **
***********************************************************************************/
int sgs_fsm_la_update_req_service_abort_request(const sgs_fsm_t* fsm_evt) {
  OAILOG_FUNC_IN(LOG_MME_APP);
  OAILOG_ERROR(
      LOG_MME_APP,
      "Dropping SGS SERVICE ABORT REQUEST message as it is received in "
      "LA-UPDATE-REQUESTED stae for UE ID %d\n",
      fsm_evt->ue_id);

  OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNok);
}

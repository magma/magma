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

/*! \file mme_app_sgsap.c
   \brief
   \author
   \version 1.0
   \company
   \email:
*/

#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#include "assertions.h"
#include "conversions.h"
#include "msc.h"
#include "log.h"
#include "intertask_interface.h"
#include "mme_app_ue_context.h"
#include "mme_app_defs.h"
#include "mme_app_sgs_fsm.h"
#include "service303.h"

/****************************************************************************
 **                                                                        **
 ** Name:    mme_app_handle_sgsap_paging_request()                         **
 **                                                                        **
 ** Description: Processes the SGSAP Paging Request message re-            **
 **      ceived from the SGS task and invokes FSM handler based on state   **
 **                                                                        **
 ** Inputs:  itti_sgsap_paging_request_t: SGSAP Paging Request message     **
 **                                                                        **
 ** Outputs:                                                               **
 **      Return:    RETURNok, RETURNerror                                  **
 **                                                                        **
 ***************************************************************************/
int mme_app_handle_sgsap_paging_request(
  itti_sgsap_paging_request_t *const sgsap_paging_req_pP)
{
  struct ue_mm_context_s *ue_context_p = NULL;
  int rc = RETURNok;
  sgs_fsm_t sgs_fsm;
  imsi64_t imsi64 = INVALID_IMSI64;

  OAILOG_FUNC_IN(LOG_MME_APP);
  DevAssert(sgsap_paging_req_pP);

  IMSI_STRING_TO_IMSI64(sgsap_paging_req_pP->imsi, &imsi64);

  OAILOG_INFO(
    LOG_MME_APP,
    "Received SGS-PAGING REQUEST for IMSI " IMSI_64_FMT "\n",
    imsi64);
  if (
    (ue_context_p = mme_ue_context_exists_imsi(
       &mme_app_desc.mme_ue_contexts, imsi64)) == NULL) {
    OAILOG_ERROR(
      LOG_MME_APP,
      "SGS-PAGING REQUEST: Failed to find UE context for IMSI " IMSI_64_FMT
      "\n",
      imsi64);
    mme_app_send_sgsap_paging_reject(
      NULL, imsi64, sgsap_paging_req_pP->imsi_length, SGS_CAUSE_IMSI_UNKNOWN);
    increment_counter("sgsap_paging_reject", 1, 1, "cause", "imsi_unknown");
    OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNerror);
  }
  if (ue_context_p->sgs_context == NULL) {
    OAILOG_ERROR(
      LOG_MME_APP,
      "SGS context not created for IMSI " IMSI_64_FMT "\n",
      imsi64);
    mme_app_send_sgsap_paging_reject(
      NULL,
      imsi64,
      sgsap_paging_req_pP->imsi_length,
      SGS_CAUSE_IMSI_DETACHED_FOR_NONEPS_SERVICE);
    increment_counter(
      "sgsap_paging_reject", 1, 1, "cause", "SGS context not created");
    unlock_ue_contexts(ue_context_p);
    OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNerror);
  }
  ue_context_p->sgs_context->sgsap_msg = (void *) sgsap_paging_req_pP;
  sgs_fsm.primitive = _SGS_PAGING_REQUEST;
  sgs_fsm.ue_id = ue_context_p->mme_ue_s1ap_id;
  sgs_fsm.ctx = (void *) ue_context_p->sgs_context;

  /* Invoke SGS FSM */
  if (RETURNok != (rc = sgs_fsm_process(&sgs_fsm))) {
    OAILOG_WARNING(
      LOG_MME_APP,
      "Failed  to execute SGS State machine for ue_id :%u \n",
      ue_context_p->mme_ue_s1ap_id);
  }
  ue_context_p->sgs_context->sgsap_msg = NULL;
  unlock_ue_contexts(ue_context_p);
  OAILOG_FUNC_RETURN(LOG_MME_APP, rc);
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
 ** Outputs:                                                               **
 **      Return:    RETURNok, RETURNerror                                  **
 **                                                                        **
 ***************************************************************************/

int mme_app_notify_service_reject_to_nas(
  mme_ue_s1ap_id_t ue_id,
  uint8_t emm_cause,
  uint8_t failed_procedure)
{
  int rc = RETURNok;
  MessageDef *message_p = NULL;
  itti_nas_notify_service_reject_t *itti_nas_notify_service_reject_p = NULL;
  OAILOG_FUNC_IN(LOG_MME_APP);
  OAILOG_INFO(
    LOG_MME_APP,
    " Ongoing Service request procedure failed,"
    "send Notify Service Reject to NAS for ue_id :%u \n",
    ue_id);
  message_p = itti_alloc_new_message(TASK_MME_APP, NAS_NOTIFY_SERVICE_REJECT);
  itti_nas_notify_service_reject_p =
    &message_p->ittiMsg.nas_notify_service_reject;
  memset(
    (void *) itti_nas_notify_service_reject_p,
    0,
    sizeof(itti_nas_extended_service_req_t));

  itti_nas_notify_service_reject_p->ue_id = ue_id;
  itti_nas_notify_service_reject_p->emm_cause = emm_cause;
  itti_nas_notify_service_reject_p->failed_procedure = failed_procedure;

  OAILOG_INFO(
    LOG_MME_APP, " Send Notify service reject to NAS for UE-id :%u \n", ue_id);
  rc = itti_send_msg_to_task(TASK_NAS_MME, INSTANCE_DEFAULT, message_p);

  OAILOG_FUNC_RETURN(LOG_MME_APP, rc);
}

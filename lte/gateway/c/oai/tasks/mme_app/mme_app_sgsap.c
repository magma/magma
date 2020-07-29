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

/*! \file mme_app_sgsap.c
   \brief
   \author
   \version 1.0
   \company
   \email:
*/

#include <stdio.h>
#include <string.h>
#include <stdint.h>

#include "conversions.h"
#include "log.h"
#include "intertask_interface.h"
#include "mme_app_ue_context.h"
#include "mme_app_defs.h"
#include "mme_app_sgs_fsm.h"
#include "service303.h"
#include "3gpp_36.401.h"
#include "common_defs.h"
#include "common_types.h"
#include "intertask_interface_types.h"
#include "itti_types.h"
#include "mme_app_desc.h"
#include "sgs_messages_types.h"

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
    mme_app_desc_t* mme_app_desc_p,
    itti_sgsap_paging_request_t* const sgsap_paging_req_pP) {
  struct ue_mm_context_s* ue_context_p = NULL;
  int rc                               = RETURNok;
  sgs_fsm_t sgs_fsm;
  imsi64_t imsi64 = INVALID_IMSI64;

  OAILOG_FUNC_IN(LOG_MME_APP);
  if (sgsap_paging_req_pP == NULL) {
    OAILOG_ERROR(
        LOG_MME_APP, "Invalid SGSAP Paging Request ITTI message received\n");
    OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNerror);
  }

  IMSI_STRING_TO_IMSI64(sgsap_paging_req_pP->imsi, &imsi64);

  OAILOG_INFO(
      LOG_MME_APP, "Received SGS-PAGING REQUEST for IMSI " IMSI_64_FMT "\n",
      imsi64);
  if ((ue_context_p = mme_ue_context_exists_imsi(
           &mme_app_desc_p->mme_ue_contexts, imsi64)) == NULL) {
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
        LOG_MME_APP, "SGS context not created for IMSI " IMSI_64_FMT "\n",
        imsi64);
    mme_app_send_sgsap_paging_reject(
        NULL, imsi64, sgsap_paging_req_pP->imsi_length,
        SGS_CAUSE_IMSI_DETACHED_FOR_NONEPS_SERVICE);
    increment_counter(
        "sgsap_paging_reject", 1, 1, "cause", "SGS context not created");
    OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNerror);
  }
  ue_context_p->sgs_context->sgsap_msg = (void*) sgsap_paging_req_pP;
  sgs_fsm.primitive                    = _SGS_PAGING_REQUEST;
  sgs_fsm.ue_id                        = ue_context_p->mme_ue_s1ap_id;
  sgs_fsm.ctx                          = (void*) ue_context_p->sgs_context;

  // Invoke SGS FSM
  rc = sgs_fsm_process(&sgs_fsm);
  if (rc != RETURNok) {
    OAILOG_WARNING(
        LOG_MME_APP, "Failed  to execute SGS State machine for ue_id :%u \n",
        ue_context_p->mme_ue_s1ap_id);
  }
  ue_context_p->sgs_context->sgsap_msg = NULL;
  OAILOG_FUNC_RETURN(LOG_MME_APP, rc);
}

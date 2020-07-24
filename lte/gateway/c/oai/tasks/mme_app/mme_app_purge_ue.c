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

/*! \file mme_app_purge_ue.c
   \brief
   \author Sebastien ROUX, Lionel GAUTHIER
   \version 1.0
   \company Eurecom
   \email: lionel.gauthier@eurecom.fr
*/

#include <stdio.h>
#include <string.h>
#include <stdint.h>

#include "common_types.h"
#include "common_defs.h"
#include "conversions.h"
#include "log.h"
#include "intertask_interface.h"
#include "mme_app_ue_context.h"
#include "mme_app_defs.h"
#include "intertask_interface_types.h"
#include "itti_types.h"
#include "mme_app_desc.h"
#include "s6a_messages_types.h"

int mme_app_send_s6a_purge_ue_req(
    mme_app_desc_t* mme_app_desc_p,
    struct ue_mm_context_s* const ue_context_pP) {
  struct ue_mm_context_s* ue_context_p = NULL;
  uint64_t imsi                        = 0;
  MessageDef* message_p                = NULL;
  s6a_purge_ue_req_t* s6a_pur_p        = NULL;
  int rc                               = RETURNok;

  OAILOG_FUNC_IN(LOG_MME_APP);
  imsi = ue_context_pP->emm_context._imsi64;
  OAILOG_DEBUG(LOG_MME_APP, "Handling imsi " IMSI_64_FMT "\n", imsi);

  if ((ue_context_p = mme_ue_context_exists_imsi(
           &mme_app_desc_p->mme_ue_contexts, imsi)) == NULL) {
    OAILOG_ERROR(
        LOG_MME_APP, "That's embarrassing as we don't know this IMSI\n");
    OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNerror);
  }

  message_p = itti_alloc_new_message(TASK_MME_APP, S6A_PURGE_UE_REQ);
  if (message_p == NULL) {
    OAILOG_WARNING(
        LOG_MME_APP,
        "Failed to allocate memory for S6A_PURGE_UE_REQ, IMSI:" IMSI_64_FMT
        "\n",
        imsi);
    OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNerror);
  }

  s6a_pur_p = &message_p->ittiMsg.s6a_purge_ue_req;
  memset((void*) s6a_pur_p, 0, sizeof(s6a_purge_ue_req_t));
  IMSI64_TO_STRING(
      imsi, s6a_pur_p->imsi, ue_context_p->emm_context._imsi.length);
  s6a_pur_p->imsi_length = strlen(s6a_pur_p->imsi);
  OAILOG_INFO(
      LOG_MME_APP, "Sent PUR to S6a TASK for IMSI " IMSI_64_FMT "\n", imsi);

  rc = send_msg_to_task(&mme_app_task_zmq_ctx, TASK_S6A, message_p);

  OAILOG_FUNC_RETURN(LOG_MME_APP, rc);
}

int mme_app_handle_s6a_purge_ue_ans(const s6a_purge_ue_ans_t* const pua_pP) {
  uint64_t imsi = 0;

  OAILOG_FUNC_IN(LOG_MME_APP);
  if (pua_pP == NULL) {
    OAILOG_ERROR(
        LOG_MME_APP, "Invalid S6a Purge UE Answer ITTI message received\n");
    OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNerror);
  }
  IMSI_STRING_TO_IMSI64((char*) pua_pP->imsi, &imsi);
  OAILOG_INFO(LOG_MME_APP, "Received PUA for imsi " IMSI_64_FMT "\n", imsi);

  if (pua_pP->result.present == S6A_RESULT_BASE) {
    if (pua_pP->result.choice.base != DIAMETER_SUCCESS) {
      OAILOG_WARNING(
          LOG_MME_APP,
          "PUR/PUA procedure returned non success "
          "(PUA.result.choice.base=%d)\n",
          pua_pP->result.choice.base);
    } else {
      OAILOG_INFO(
          LOG_MME_APP, "Received PUA Success for imsi " IMSI_64_FMT "\n", imsi);
    }
  } else {
    /*
     * The Purge Ue procedure has failed.
     */
    OAILOG_WARNING(
        LOG_MME_APP,
        "PUR/PUA procedure returned non success (ULA.result.present=%d)\n",
        pua_pP->result.present);
  }
  OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNok);
}

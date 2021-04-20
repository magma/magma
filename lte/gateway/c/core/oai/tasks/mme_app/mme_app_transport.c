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

/*! \file mme_app_transport.c
  \brief
  \author Lionel Gauthier
  \company Eurecom
  \email: lionel.gauthier@eurecom.fr
*/
#include <stdio.h>

#include "bstrlib.h"
#include "dynamic_memory_check.h"
#include "log.h"
#include "intertask_interface.h"
#include "mme_app_ue_context.h"
#include "mme_app_defs.h"
#include "common_defs.h"
#include "mme_app_itti_messaging.h"
#include "mme_app_state.h"
#include "3gpp_36.401.h"
#include "nas/as_message.h"
#include "common_types.h"
#include "intertask_interface_types.h"
#include "itti_types.h"
#include "mme_app_desc.h"
#include "s1ap_messages_types.h"
#include "sgs_messages_types.h"

//------------------------------------------------------------------------------
int mme_app_handle_nas_dl_req(
    const mme_ue_s1ap_id_t ue_id, bstring nas_msg,
    nas_error_code_t transaction_status)
//------------------------------------------------------------------------------
{
  OAILOG_FUNC_IN(LOG_MME_APP);
  MessageDef* message_p           = NULL;
  int rc                          = RETURNok;
  enb_ue_s1ap_id_t enb_ue_s1ap_id = 0;

  message_p = itti_alloc_new_message(TASK_MME_APP, S1AP_NAS_DL_DATA_REQ);

  mme_app_desc_t* mme_app_desc_p = get_mme_nas_state(false);
  if (!mme_app_desc_p) {
    OAILOG_CRITICAL(
        LOG_MME_APP,
        "DOWNLINK NAS TRANSPORT. Failed to get global mme_app_desc context \n");
    OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNerror);
  }
  ue_mm_context_t* ue_context = mme_ue_context_exists_mme_ue_s1ap_id(ue_id);
  if (ue_context) {
    enb_ue_s1ap_id = ue_context->enb_ue_s1ap_id;
  } else {
    OAILOG_WARNING(
        LOG_MME_APP,
        " DOWNLINK NAS TRANSPORT. Null UE Context for "
        "mme_ue_s1ap_id " MME_UE_S1AP_ID_FMT "\n",
        ue_id);
    OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNerror);
  }

  S1AP_NAS_DL_DATA_REQ(message_p).enb_ue_s1ap_id = enb_ue_s1ap_id;
  S1AP_NAS_DL_DATA_REQ(message_p).mme_ue_s1ap_id = ue_id;
  S1AP_NAS_DL_DATA_REQ(message_p).nas_msg        = bstrcpy(nas_msg);
  bdestroy_wrapper(&nas_msg);

  message_p->ittiMsgHeader.imsi = ue_context->emm_context._imsi64;
  /*
   * Store the S1AP NAS DL DATA REQ in case of IMSI or combined EPS/IMSI detach
   * in sgs context and send it after recieving the SGS IMSI Detach Ack from SGS
   * task.
   */
  if (ue_context->sgs_context != NULL) {
    if (((ue_context->sgs_detach_type ==
          SGS_EXPLICIT_UE_INITIATED_IMSI_DETACH_FROM_NONEPS) ||
         (ue_context->sgs_detach_type ==
          SGS_COMBINED_UE_INITIATED_IMSI_DETACH_FROM_EPS_N_NONEPS)) &&
        (ue_context->sgs_context->ts9_timer.id != MME_APP_TIMER_INACTIVE_ID)) {
      ue_context->sgs_context->message_p = message_p;
    } else { /* Send the S1AP NAS DL DATA REQ to S1AP */
      rc = send_msg_to_task(&mme_app_task_zmq_ctx, TASK_S1AP, message_p);
    }
  } else {
    rc = send_msg_to_task(&mme_app_task_zmq_ctx, TASK_S1AP, message_p);
  }

  /*
   * Move the UE to ECM Connected State ,if not in connected state already
   * S1 Signaling connection gets established via first DL NAS Trasnport
   * message in some scenarios so check the state
   * first
   */
  if (ue_context->ecm_state != ECM_CONNECTED) {
    OAILOG_DEBUG(
        LOG_MME_APP,
        "MME_APP:DOWNLINK NAS TRANSPORT. Establishing S1 sig connection. "
        "mme_ue_s1ap_id = " MME_UE_S1AP_ID_FMT
        "\t "
        "enb_ue_s1ap_id = " ENB_UE_S1AP_ID_FMT " \n",
        ue_id, enb_ue_s1ap_id);
    mme_ue_context_update_ue_sig_connection_state(
        &mme_app_desc_p->mme_ue_contexts, ue_context, ECM_CONNECTED);
  }

  // Check the transaction status. And trigger the UE context release command
  // accrordingly.
  if (transaction_status != AS_SUCCESS) {
    ue_context->ue_context_rel_cause = S1AP_NAS_NORMAL_RELEASE;
    // Notify S1AP to send UE Context Release Command to eNB.
    mme_app_itti_ue_context_release(
        ue_context, ue_context->ue_context_rel_cause);
  }
  OAILOG_FUNC_RETURN(LOG_MME_APP, rc);
}

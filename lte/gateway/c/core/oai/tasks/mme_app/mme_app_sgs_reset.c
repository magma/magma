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

  Source      mme_app_sgs_reset.c

  Version

  Date

  Product    MME app

  Subsystem  SGSAP Reset message handling
  Author

*****************************************************************************/

#include <stdbool.h>
#include <stddef.h>
#include <mme_app_state.h>

#include "log.h"
#include "mme_app_sgs_fsm.h"
#include "mme_app_defs.h"
#include "common_defs.h"
#include "common_types.h"
#include "hashtable.h"
#include "mme_api.h"
#include "mme_app_desc.h"
#include "mme_app_ue_context.h"
#include "sgs_messages_types.h"

/****************************************************************************/
/****************  E X T E R N A L    D E F I N I T I O N S  ****************/
/****************************************************************************/

/****************************************************************************
 **                                                                        **
 ** Name:    mme_app_handle_sgsap_reset_indication()                       **
 **                                                                        **
 ** Description: Handles the SGSAP Reset Indication from VLR               **
 **                                                                        **
 ** Inputs:  reset_indication_pP:   The received SGS Reset Indication      **
 **                                                                        **
 ** Outputs:                                                               **
 **          Return:    RETURNok, RETURNerror                              **
 **                                                                        **
 ***************************************************************************/
int mme_app_handle_sgsap_reset_indication(
    itti_sgsap_vlr_reset_indication_t* const reset_indication_pP) {
  int rc = RETURNerror;
  OAILOG_FUNC_IN(LOG_MME_APP);
  OAILOG_INFO(
      LOG_MME_APP, " Received SGSAP-Reset Indication from VLR :%s \n",
      reset_indication_pP->vlr_name);

  /* Handle VLR Reset for each SGS associated UE */
  hash_table_ts_t* mme_state_imsi_ht = get_mme_ue_state();
  hashtable_ts_apply_callback_on_elements(
      mme_state_imsi_ht, mme_app_handle_reset_indication, NULL, NULL);
  OAILOG_FUNC_RETURN(LOG_MME_APP, rc);
}

/****************************************************************************
 **                                                                        **
 ** Name:    mme_app_handle_reset_indication()                             **
 **                                                                        **
 ** Description: Handles SGSAP Reset Indication message from VLR           **
 **              This message is sent from the VLR to the MME to indicate  **
 **              that a failure in the VLR has occurred and all the SGs    **
 **              associations to the VLR are to be marked as invalid.      **
 **                                                                        **
 ** Inputs:  keyP: Hash key                                                **
 **          ue_context_pP: Pointer to UE context                          **
 **          unused_param_pP: Unused param list                            **
 **          unused_result_pP: Unused result                               **
 ** Outputs:                                                               **
 **          Return:    RETURNok, RETURNerror                              **
 **                                                                        **
 ***************************************************************************/
bool mme_app_handle_reset_indication(
    const hash_key_t keyP, void* const ue_context_pP, void* unused_param_pP,
    void** unused_result_pP) {
  int rc = RETURNerror;
  sgs_fsm_t sgs_fsm;
  OAILOG_FUNC_IN(LOG_MME_APP);

  struct ue_mm_context_s* const ue_context_p =
      (struct ue_mm_context_s*) ue_context_pP;
  if (ue_context_p == NULL) {
    OAILOG_WARNING(LOG_MME_APP, "UE context not found \n");
    OAILOG_FUNC_RETURN(LOG_MME_APP, rc);
  }
  if (ue_context_p->mm_state == UE_UNREGISTERED) {
    OAILOG_ERROR(
        LOG_MME_APP, "UE is not registered for ue_id:" MME_UE_S1AP_ID_FMT "\n",
        ue_context_p->mme_ue_s1ap_id);
    OAILOG_FUNC_RETURN(LOG_MME_APP, rc);
  }
  if (ue_context_p->sgs_context == NULL) {
    OAILOG_ERROR(
        LOG_MME_APP,
        "SGS context not created for ue_id:" MME_UE_S1AP_ID_FMT "\n",
        ue_context_p->mme_ue_s1ap_id);
    OAILOG_FUNC_RETURN(LOG_MME_APP, rc);
  }

  ue_context_p->sgs_context->sgsap_msg = NULL; /* sgs message */
  sgs_fsm.primitive                    = _SGS_RESET_INDICATION;
  sgs_fsm.ue_id                        = ue_context_p->mme_ue_s1ap_id;
  sgs_fsm.ctx                          = (void*) ue_context_p->sgs_context;

  /* Invoke SGS FSM */
  if (RETURNok != (rc = sgs_fsm_process(&sgs_fsm))) {
    OAILOG_WARNING(
        LOG_MME_APP, "Failed  to execute SGS State machine for ue_id :%u \n",
        ue_context_p->mme_ue_s1ap_id);
  }
  OAILOG_FUNC_RETURN(LOG_MME_APP, rc);
}

/****************************************************************************
 **                                                                        **
 ** Name:    sgs_fsm_associated_reset_indication()                         **
 **                                                                        **
 ** Description: Handles SGSAP Reset Indication message from VLR           **
 **              While SGS context is in Assocaited state                  **
 **                                                                        **
 ** Inputs:  None                                                          **
 **                                                                        **
 ** Outputs:                                                               **
 **          Return:    RETURNok, RETURNerror                              **
 **                                                                        **
 ***************************************************************************/
int sgs_fsm_associated_reset_indication(const sgs_fsm_t* fsm_evt) {
  OAILOG_FUNC_IN(LOG_MME_APP);
  if (fsm_evt == NULL) {
    OAILOG_ERROR(LOG_MME_APP, "Invalid SGS FSM Event object received\n");
    OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNerror);
  }
  OAILOG_DEBUG(
      LOG_MME_APP,
      "Handle Reset Indication in Associated state for ue-id :%u \n",
      fsm_evt->ue_id);
  sgs_context_t* sgs_context = (sgs_context_t*) fsm_evt->ctx;
  if (sgs_context == NULL) {
    OAILOG_WARNING(
        LOG_MME_APP, " Strange sgs context is NULL for ue_id :%u \n",
        fsm_evt->ue_id);
  }
  sgs_context->vlr_reliable = false;
  OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNok);
}

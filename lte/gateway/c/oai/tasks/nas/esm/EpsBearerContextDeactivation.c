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

#include <stdbool.h>
#include <stdlib.h>

#include "bstrlib.h"
#include "log.h"
#include "dynamic_memory_check.h"
#include "common_types.h"
#include "3gpp_24.007.h"
#include "mme_app_ue_context.h"
#include "mme_app_defs.h"
#include "esm_proc.h"
#include "emm_data.h"
#include "esm_data.h"
#include "esm_cause.h"
#include "esm_ebr.h"
#include "esm_ebr_context.h"
#include "emm_sap.h"
#include "mme_config.h"
#include "3gpp_24.301.h"
#include "3gpp_36.401.h"
#include "EsmCause.h"
#include "common_defs.h"
#include "emm_esmDef.h"
#include "esm_sapDef.h"
#include "esm_pt.h"
#include "mme_app_state.h"

/****************************************************************************/
/****************  E X T E R N A L    D E F I N I T I O N S  ****************/
/****************************************************************************/

/****************************************************************************/
/*******************  L O C A L    D E F I N I T I O N S  *******************/
/****************************************************************************/

/*
   --------------------------------------------------------------------------
   Internal data handled by the EPS bearer context deactivation procedure
   in the MME
   --------------------------------------------------------------------------
*/
/*
   Timer handlers
*/
/* Maximum value of the deactivate EPS bearer context request
   retransmission counter */
#define EPS_BEARER_DEACTIVATE_COUNTER_MAX 5

static int eps_bearer_deactivate(
    emm_context_t* emm_context_p, ebi_t ebi, STOLEN_REF bstring* msg);

/****************************************************************************/
/******************  E X P O R T E D    F U N C T I O N S  ******************/
/****************************************************************************/
extern int pdn_connectivity_delete(emm_context_t* emm_context, pdn_cid_t pid);

/*
   --------------------------------------------------------------------------
    EPS bearer context deactivation procedure executed by the MME
   --------------------------------------------------------------------------
*/
/****************************************************************************
 **                                                                        **
 ** Name:    esm_proc_eps_bearer_context_deactivate()                  **
 **                                                                        **
 ** Description: Locally releases the EPS bearer context identified by the **
 **      given EPS bearer identity, without peer-to-peer signal-   **
 **      ling between the UE and the MME, or checks whether an EPS **
 **      bearer context with specified EPS bearer identity has     **
 **      been activated for the given UE.                          **
 **                                                                        **
 ** Inputs:  ue_id:      UE lower layer identifier                  **
 **      is local:  true if the EPS bearer context has to be   **
 **             locally released without peer-to-peer si-  **
 **             gnalling between the UE and the MME        **
 **      ebi:       EPS bearer identity of the EPS bearer con- **
 **             text to be deactivated                     **
 **      Others:    _esm_data                                  **
 **                                                                        **
 ** Outputs:     pid:       Identifier of the PDN connection the EPS   **
 **             bearer belongs to                          **
 **      bid:       Identifier of the released EPS bearer con- **
 **             text entry                                 **
 **      esm_cause: Cause code returned upon ESM procedure     **
 **             failure                                    **
 **      Return:    RETURNok, RETURNerror                      **
 **      Others:    None                                       **
 **                                                                        **
 ***************************************************************************/
int esm_proc_eps_bearer_context_deactivate(
    emm_context_t* const emm_context_p, const bool is_local, const ebi_t ebi,
    pdn_cid_t* pid, int* const bidx, esm_cause_t* const esm_cause) {
  OAILOG_FUNC_IN(LOG_NAS_ESM);
  int rc = RETURNerror;
  ue_mm_context_t* ue_mm_context =
      PARENT_STRUCT(emm_context_p, struct ue_mm_context_s, emm_context);

  if (is_local) {
    if (ebi != ESM_SAP_ALL_EBI) {
      /*
       * Locally release the specified EPS bearer context
       */
      rc = eps_bearer_release(emm_context_p, ebi, pid, bidx);
    } else if (emm_context_p) {
      /*
       * Locally release all the EPS bearer contexts
       */
      for (int bix = 0; bix < BEARERS_PER_UE; bix++) {
        if (ue_mm_context->bearer_contexts[bix]) {
          *pid = ue_mm_context->bearer_contexts[bix]->pdn_cx_id;
          rc   = eps_bearer_release(
              emm_context_p, ue_mm_context->bearer_contexts[bix]->ebi, pid,
              bidx);

          if (rc != RETURNok) {
            OAILOG_FUNC_RETURN(LOG_NAS_ESM, rc);
          }
        }
      }
      rc = RETURNok;
    }

    OAILOG_FUNC_RETURN(LOG_NAS_ESM, rc);
  }

  OAILOG_INFO(
      LOG_NAS_ESM,
      "ESM-PROC  - EPS bearer context deactivation "
      "(ue_id=" MME_UE_S1AP_ID_FMT ", ebi=%d)\n",
      ue_mm_context->mme_ue_s1ap_id, ebi);

  if ((ue_mm_context) && (*pid < MAX_APN_PER_UE)) {
    if (ue_mm_context->pdn_contexts[*pid] == NULL) {
      OAILOG_ERROR(
          LOG_NAS_ESM,
          "ESM-PROC  - PDN connection %d has not been "
          "allocated for ue id " MME_UE_S1AP_ID_FMT "\n",
          *pid, ue_mm_context->mme_ue_s1ap_id);
      *esm_cause = ESM_CAUSE_PROTOCOL_ERROR;
    } else {
      int i;

      *esm_cause = ESM_CAUSE_INVALID_EPS_BEARER_IDENTITY;

      for (i = 0; i < BEARERS_PER_UE; i++) {
        if ((ue_mm_context->pdn_contexts[*pid]->bearer_contexts[i] <= 0) ||
            (ue_mm_context->bearer_contexts[i]->pdn_cx_id != *pid)) {
          continue;
        }

        if (ebi != ESM_SAP_ALL_EBI) {
          if (ue_mm_context->bearer_contexts[i]->ebi != ebi) {
            continue;
          }
        }
        /*
         * The EPS bearer context to be released is valid
         */
        *esm_cause = ESM_CAUSE_SUCCESS;
        rc         = RETURNok;
      }
    }
  }

  OAILOG_FUNC_RETURN(LOG_NAS_ESM, rc);
}

/****************************************************************************
 **                                                                        **
 ** Name:    esm_proc_eps_bearer_context_deactivate_request()          **
 **                                                                        **
 ** Description: Initiates the EPS bearer context deactivation procedure   **
 **                                                                        **
 **      3GPP TS 24.301, section 6.4.4.2                           **
 **      If a NAS signalling connection exists, the MME initiates  **
 **      the EPS bearer context deactivation procedure by sending  **
 **      a DEACTIVATE EPS BEARER CONTEXT REQUEST message to the    **
 **      UE, starting timer T3495 and entering state BEARER CON-   **
 **      TEXT INACTIVE PENDING.                                    **
 **                                                                        **
 ** Inputs:  is_standalone: Not used - Always true                     **
 **      ue_id:      UE lower layer identifier                  **
 **      ebi:       EPS bearer identity                        **
 **      msg:       Encoded ESM message to be sent             **
 **      ue_triggered:  true if the EPS bearer context procedure   **
 **             was triggered by the UE (not used)         **
 **                                                                        **
 ** Outputs:     None                                                      **
 **      Return:    RETURNok, RETURNerror                      **
 **      Others:    None                                       **
 **                                                                        **
 ***************************************************************************/
int esm_proc_eps_bearer_context_deactivate_request(
    const bool is_standalone, emm_context_t* const emm_context_p,
    const ebi_t ebi, STOLEN_REF bstring* msg, const bool ue_triggered) {
  OAILOG_FUNC_IN(LOG_NAS_ESM);
  int rc;
  mme_ue_s1ap_id_t ue_id =
      PARENT_STRUCT(emm_context_p, struct ue_mm_context_s, emm_context)
          ->mme_ue_s1ap_id;

  OAILOG_INFO(
      LOG_NAS_ESM,
      "ESM-PROC  - Initiate EPS bearer context deactivation "
      "(ue_id=" MME_UE_S1AP_ID_FMT ", ebi=%d)\n",
      ue_id, ebi);
  /*
   * Send deactivate EPS bearer context request message and
   * * * * start timer T3495
   */
  /*Currently we only support single bearear deactivation at NAS*/
  rc  = eps_bearer_deactivate(emm_context_p, ebi, msg);
  msg = NULL;

  if (rc != RETURNerror) {
    /*
     * Set the EPS bearer context state to ACTIVE PENDING
     */
    rc = esm_ebr_set_status(
        emm_context_p, ebi, ESM_EBR_INACTIVE_PENDING, ue_triggered);

    if (rc != RETURNok) {
      /*
       * The EPS bearer context was already in ACTIVE state
       */
      OAILOG_WARNING(
          LOG_NAS_ESM,
          "ESM-PROC  - EBI %d was already INACTIVE PENDING for ue "
          "id " MME_UE_S1AP_ID_FMT "\n",
          ebi, ue_id);
    }
  }
  OAILOG_FUNC_RETURN(LOG_NAS_ESM, rc);
}

/****************************************************************************
 **                                                                        **
 ** Name:    esm_proc_eps_bearer_context_deactivate_accept()           **
 **                                                                        **
 ** Description: Performs EPS bearer context deactivation procedure accep- **
 **      ted by the UE.                                            **
 **                                                                        **
 **      3GPP TS 24.301, section 6.4.4.3                           **
 **      Upon receipt of the DEACTIVATE EPS BEARER CONTEXT ACCEPT  **
 **      message, the MME shall enter the state BEARER CONTEXT     **
 **      INACTIVE and stop the timer T3495.                        **
 **                                                                        **
 ** Inputs:  ue_id:      UE local identifier                        **
 **      ebi:       EPS bearer identity                        **
 **      Others:    None                                       **
 **                                                                        **
 ** Outputs:     esm_cause: Cause code returned upon ESM procedure     **
 **             failure                                    **
 **      Return:    The identifier of the PDN connection to be **
 **             released, if it exists;                    **
 **             RETURNerror otherwise.                     **
 **      Others:    T3495                                      **
 **                                                                        **
 ***************************************************************************/
pdn_cid_t esm_proc_eps_bearer_context_deactivate_accept(
    emm_context_t* emm_context_p, ebi_t ebi, esm_cause_t* esm_cause) {
  OAILOG_FUNC_IN(LOG_NAS_ESM);
  int rc                        = RETURNerror;
  pdn_cid_t pid                 = MAX_APN_PER_UE;
  ue_mm_context_t* ue_context_p = NULL;
  bool delete_default_bearer    = false;
  int bid                       = BEARERS_PER_UE;
  teid_t s_gw_teid_s11_s4       = 0;

  ue_context_p =
      PARENT_STRUCT(emm_context_p, struct ue_mm_context_s, emm_context);
  OAILOG_INFO(
      LOG_NAS_ESM,
      "ESM-PROC  - EPS bearer context deactivation "
      "accepted by the UE (ue_id=" MME_UE_S1AP_ID_FMT ", ebi=%d)\n",
      ue_context_p->mme_ue_s1ap_id, ebi);
  /*
   * Stop T3495 timer if running
   */
  rc = esm_ebr_stop_timer(emm_context_p, ebi);

  if (rc != RETURNerror) {
    /*
     * Release the EPS bearer context
     */
    rc = eps_bearer_release(emm_context_p, ebi, &pid, &bid);

    if (rc != RETURNok) {
      /*
       * Failed to release the EPS bearer context
       */
      *esm_cause = ESM_CAUSE_PROTOCOL_ERROR;
      pid        = RETURNerror;
      OAILOG_FUNC_RETURN(LOG_NAS_ESM, pid);
    }
  }

  if (ue_context_p->pdn_contexts[pid] == NULL) {
    OAILOG_ERROR_UE(
        LOG_MME_APP, ue_context_p->emm_context._imsi64,
        "pdn_contexts is NULL for "
        "MME UE S1AP Id: " MME_UE_S1AP_ID_FMT "ebi-%u\n",
        ue_context_p->mme_ue_s1ap_id, ebi);
    OAILOG_FUNC_RETURN(LOG_NAS_ESM, RETURNerror);
  }

  s_gw_teid_s11_s4 = ue_context_p->pdn_contexts[pid]->s_gw_teid_s11_s4;

  // If bearer id == 0, default bearer is deleted
  if (ue_context_p->pdn_contexts[pid]->default_ebi == ebi) {
    delete_default_bearer = true;
    // Release the default bearer
    rc = mme_api_unsubscribe(NULL);

    if (rc != RETURNerror) {
      /*
       * Delete the PDN connection entry
       */
      pdn_connectivity_delete(emm_context_p, pid);
      // Free PDN context
      free_wrapper((void**) &ue_context_p->pdn_contexts[pid]);
      // Free bearer context entry
      if (ue_context_p->bearer_contexts[bid]) {
        free_wrapper((void**) &ue_context_p->bearer_contexts[bid]);
      }
    }
  } else {
    OAILOG_INFO(
        LOG_NAS_ESM,
        "ESM-PROC  - Removing dedicated bearer context "
        "for UE (ue_id=" MME_UE_S1AP_ID_FMT ", ebi=%d)\n",
        ue_context_p->mme_ue_s1ap_id, ebi);
    // Remove dedicated bearer context
    free_wrapper((void**) &ue_context_p->bearer_contexts[bid]);
  }
  /* In case of PDN disconnect, no need to inform MME/SPGW as the session would
   * have been already released
   */
  if (!emm_context_p->esm_ctx.is_pdn_disconnect) {
    // Send delete dedicated bearer response to SPGW
    send_delete_dedicated_bearer_rsp(
        ue_context_p, delete_default_bearer, &ebi, 1, s_gw_teid_s11_s4,
        REQUEST_ACCEPTED);
  }

  // Reset is_pdn_disconnect flag
  if (emm_context_p->esm_ctx.is_pdn_disconnect) {
    emm_context_p->esm_ctx.is_pdn_disconnect = false;
  }

  OAILOG_FUNC_RETURN(LOG_NAS_ESM, pid);
}

/****************************************************************************/
/*********************  L O C A L    F U N C T I O N S  *********************/
/****************************************************************************/
/*
   --------------------------------------------------------------------------
                Timer handlers
   --------------------------------------------------------------------------
*/
/****************************************************************************
 **                                                                        **
 ** Name:    eps_bearer_deactivate_t3495_handler()                    **
 **                                                                        **
 ** Description: T3495 timeout handler                                     **
 **                                                                        **
 **              3GPP TS 24.301, section 6.4.4.5, case a                   **
 **      On the first expiry of the timer T3495, the MME shall re- **
 **      send the DEACTIVATE EPS BEARER CONTEXT REQUEST and shall  **
 **      reset and restart timer T3495. This retransmission is     **
 **      repeated four times, i.e. on the fifth expiry of timer    **
 **      T3495, the MME shall abort the procedure and deactivate   **
 **      the EPS bearer context locally.                           **
 **                                                                        **
 ** Inputs:  args:      handler parameters                         **
 **      Others:    None                                       **
 **                                                                        **
 ** Outputs:     None                                                      **
 **      Return:    None                                       **
 **      Others:    None                                       **
 **                                                                        **
 ***************************************************************************/
void eps_bearer_deactivate_t3495_handler(void* args, imsi64_t* imsi64) {
  OAILOG_FUNC_IN(LOG_NAS_ESM);
  int rc;
  bool delete_default_bearer = false;
  ebi_t ebi                  = 0;
  mme_ue_s1ap_id_t ue_id     = 0;
  int bid                    = BEARERS_PER_UE;

  /*
   * Get retransmission timer parameters data
   */
  esm_ebr_timer_data_t* esm_ebr_timer_data = (esm_ebr_timer_data_t*) (args);

  if (esm_ebr_timer_data) {
    /*
     * Increment the retransmission counter
     */
    esm_ebr_timer_data->count += 1;
    OAILOG_WARNING(
        LOG_NAS_ESM,
        "ESM-PROC  - T3495 timer expired (ue_id=" MME_UE_S1AP_ID_FMT
        ", ebi=%d), "
        "retransmission counter = %d\n",
        esm_ebr_timer_data->ue_id, esm_ebr_timer_data->ebi,
        esm_ebr_timer_data->count);

    *imsi64 = esm_ebr_timer_data->ctx->_imsi64;
    if (esm_ebr_timer_data->count < EPS_BEARER_DEACTIVATE_COUNTER_MAX) {
      /*
       * Re-send deactivate EPS bearer context request message to the UE
       */
      bstring b = bstrcpy(esm_ebr_timer_data->msg);
      rc        = eps_bearer_deactivate(
          esm_ebr_timer_data->ctx, esm_ebr_timer_data->ebi, &b);
      bdestroy_wrapper(&b);
    } else {
      /*
       * The maximum number of deactivate EPS bearer context request
       * message retransmission has exceed
       */
      ue_mm_context_t* ue_mm_context =
          mme_ue_context_exists_mme_ue_s1ap_id(esm_ebr_timer_data->ue_id);

      ebi   = esm_ebr_timer_data->ebi;
      ue_id = esm_ebr_timer_data->ue_id;
      // Find the index in which the bearer id is stored
      for (bid = 0; bid < BEARERS_PER_UE; bid++) {
        if (ue_mm_context->bearer_contexts[bid]) {
          if (ue_mm_context->bearer_contexts[bid]->ebi == ebi) {
            break;
          }
        }
      }

      if (bid >= BEARERS_PER_UE) {
        OAILOG_WARNING(
            LOG_NAS_ESM,
            "ESM-PROC  - Did not find bearer context for "
            "(ue_id=" MME_UE_S1AP_ID_FMT ", ebi=%d), \n",
            ue_id, ebi);
        OAILOG_FUNC_OUT(LOG_NAS_ESM);
      }
      // Fetch pdn id using bearer index
      pdn_cid_t pdn_id = ue_mm_context->bearer_contexts[bid]->pdn_cx_id;

      if (!ue_mm_context->pdn_contexts[pdn_id]) {
        OAILOG_ERROR(
            LOG_NAS_ESM,
            "eps_bearer_deactivate_t3495_handler pid context NULL for ue "
            "id " MME_UE_S1AP_ID_FMT
            "\n"
            "pdn_id=%d\n",
            ue_id, pdn_id);
        OAILOG_FUNC_OUT(LOG_NAS_ESM);
      }
      // Send bearer_deactivation_reject to MME
      teid_t s_gw_teid_s11_s4 =
          ue_mm_context->pdn_contexts[pdn_id]->s_gw_teid_s11_s4;

      if (ue_mm_context->pdn_contexts[pdn_id]->default_ebi == ebi) {
        delete_default_bearer = true;
        // Release the default bearer
        /*
         * Delete the PDN connection entry
         */
        pdn_connectivity_delete(esm_ebr_timer_data->ctx, pdn_id);
      }
      /* In case of PDN disconnect, no need to inform MME/SPGW as the session
       * would have been already released
       */
      if (!ue_mm_context->emm_context.esm_ctx.is_pdn_disconnect) {
        // Send delete_dedicated_bearer_rsp to SPGW
        send_delete_dedicated_bearer_rsp(
            ue_mm_context, delete_default_bearer, &ebi, 1, s_gw_teid_s11_s4,
            UE_NOT_RESPONDING);
      }
      // Reset is_pdn_disconnect flag
      if (ue_mm_context->emm_context.esm_ctx.is_pdn_disconnect) {
        ue_mm_context->emm_context.esm_ctx.is_pdn_disconnect = false;
      }

      /*
       * Deactivate the EPS bearer context locally without peer-to-peer
       * * * * signalling between the UE and the MME
       */
      rc = eps_bearer_release(esm_ebr_timer_data->ctx, ebi, &pdn_id, &bid);

      if (rc == RETURNerror) {
        OAILOG_WARNING(
            LOG_NAS_ESM,
            "ESM-PROC  - Could not release bearer(ue_id=" MME_UE_S1AP_ID_FMT
            ", ebi=%d), \n",
            ue_id, ebi);
      }
    }
  }

  OAILOG_FUNC_OUT(LOG_NAS_ESM);
}
/*
   --------------------------------------------------------------------------
                MME specific local functions
   --------------------------------------------------------------------------
*/

/****************************************************************************
 **                                                                        **
 ** Name:    _eps_bearer_deactivate()                                  **
 **                                                                        **
 ** Description: Sends DEACTIVATE EPS BEREAR CONTEXT REQUEST message and   **
 **      starts timer T3495                                        **
 **                                                                        **
 ** Inputs:  ue_id:      UE local identifier                        **
 **      ebi:       EPS bearer identity                        **
 **      msg:       Encoded ESM message to be sent             **
 **      Others:    None                                       **
 **                                                                        **
 ** Outputs:     None                                                      **
 **      Return:    RETURNok, RETURNerror                      **
 **      Others:    T3495                                      **
 **                                                                        **
 ***************************************************************************/
static int eps_bearer_deactivate(
    emm_context_t* emm_context_p, ebi_t ebi, STOLEN_REF bstring* msg) {
  OAILOG_FUNC_IN(LOG_NAS_ESM);
  emm_sap_t emm_sap = {0};
  int rc;
  mme_ue_s1ap_id_t ue_id =
      PARENT_STRUCT(emm_context_p, struct ue_mm_context_s, emm_context)
          ->mme_ue_s1ap_id;

  /*
   * Notify EMM that a deactivate EPS bearer context request message
   * has to be sent to the UE
   */

  emm_sap.primitive                         = EMMESM_DEACTIVATE_BEARER_REQ;
  emm_sap.u.emm_esm.ue_id                   = ue_id;
  emm_sap.u.emm_esm.ctx                     = emm_context_p;
  emm_sap.u.emm_esm.u.deactivate_bearer.ebi = ebi;
  emm_sap.u.emm_esm.u.deactivate_bearer.msg = *msg;
  bstring msg_dup                           = bstrcpy(*msg);
  rc                                        = emm_sap_send(&emm_sap);

  if (rc != RETURNerror) {
    /*
     * Start T3495 retransmission timer
     */
    rc = esm_ebr_start_timer(
        emm_context_p, ebi, msg_dup, mme_config.nas_config.t3495_sec,
        eps_bearer_deactivate_t3495_handler);
  }
  bdestroy_wrapper(&msg_dup);
  OAILOG_FUNC_RETURN(LOG_NAS_ESM, rc);
}

/****************************************************************************
 **                                                                        **
 ** Name:    eps_bearer_release()                                     **
 **                                                                        **
 ** Description: Releases the EPS bearer context identified by the given   **
 **      EPS bearer identity and enters state INACTIVE.            **
 **                                                                        **
 ** Inputs:  ue_id:      UE local identifier                        **
 **      ebi:       EPS bearer identity                        **
 **      Others:    None                                       **
 **                                                                        **
 ** Outputs:     pid:       Identifier of the PDN connection the EPS   **
 **             bearer belongs to                          **
 **      bid:       Identifier of the released EPS bearer con- **
 **             text entry                                 **
 **      Return:    RETURNok, RETURNerror                      **
 **      Others:    None                                       **
 **                                                                        **
 ***************************************************************************/
int eps_bearer_release(
    emm_context_t* emm_context_p, ebi_t ebi, pdn_cid_t* pid, int* bidx) {
  OAILOG_FUNC_IN(LOG_NAS_ESM);
  int rc = RETURNerror;

  /*
   * Release the EPS bearer context entry
   */
  ebi = esm_ebr_context_release(emm_context_p, ebi, pid, bidx);

  if (ebi == ESM_EBI_UNASSIGNED) {
    OAILOG_WARNING_UE(
        LOG_NAS_ESM, emm_context_p->_imsi64,
        "ESM-PROC  - Failed to release EPS bearer context\n");
  } else {
    /*
     * Set the EPS bearer context state to INACTIVE
     */
    rc = esm_ebr_set_status(emm_context_p, ebi, ESM_EBR_INACTIVE, false);

    if (rc != RETURNok) {
      /*
       * The EPS bearer context was already in INACTIVE state
       */
      OAILOG_WARNING_UE(
          LOG_NAS_ESM, emm_context_p->_imsi64,
          "ESM-PROC  - EBI %d was already INACTIVE\n", ebi);
    }
    /*
     * Release EPS bearer data
     */
    rc = esm_ebr_release(emm_context_p, ebi);

    if (rc != RETURNok) {
      OAILOG_WARNING_UE(
          LOG_NAS_ESM, emm_context_p->_imsi64,
          "ESM-PROC  - Failed to release EPS bearer data\n");
    }
  }

  OAILOG_FUNC_RETURN(LOG_NAS_ESM, rc);
}

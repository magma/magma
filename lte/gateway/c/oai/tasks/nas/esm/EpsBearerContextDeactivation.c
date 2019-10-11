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

#include <stdbool.h>
#include <stdlib.h>

#include "bstrlib.h"
#include "log.h"
#include "dynamic_memory_check.h"
#include "common_types.h"
#include "3gpp_24.007.h"
#include "mme_app_ue_context.h"
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
#include "nas_itti_messaging.h"
#include "esm_pt.h"

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
static void _eps_bearer_deactivate_t3495_handler(void *);

/* Maximum value of the deactivate EPS bearer context request
   retransmission counter */
#define EPS_BEARER_DEACTIVATE_COUNTER_MAX 5

static int _eps_bearer_deactivate(
  emm_context_t *ue_context,
  ebi_t ebi,
  STOLEN_REF bstring *msg);
static int _eps_bearer_release(
  emm_context_t *ue_context,
  ebi_t ebi,
  pdn_cid_t *pid,
  int *bidx);

/****************************************************************************/
/******************  E X P O R T E D    F U N C T I O N S  ******************/
/****************************************************************************/
extern int _pdn_connectivity_delete(emm_context_t *emm_context, pdn_cid_t pid);

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
  emm_context_t *const ue_context,
  const bool is_local,
  const ebi_t ebi,
  pdn_cid_t *pid,
  int *const bidx,
  esm_cause_t *const esm_cause)
{
  OAILOG_FUNC_IN(LOG_NAS_ESM);
  int rc = RETURNerror;
  ue_mm_context_t *ue_mm_context =
    PARENT_STRUCT(ue_context, struct ue_mm_context_s, emm_context);

  if (is_local) {
    if (ebi != ESM_SAP_ALL_EBI) {
      /*
       * Locally release the specified EPS bearer context
       */
      rc = _eps_bearer_release(ue_context, ebi, pid, bidx);
    } else if (ue_context) {
      /*
       * Locally release all the EPS bearer contexts
       */
      for (int bix = 0; bix < BEARERS_PER_UE; bix++) {
        if (ue_mm_context->bearer_contexts[bix]) {
          *pid = ue_mm_context->bearer_contexts[bix]->pdn_cx_id;
          rc = _eps_bearer_release(
            ue_context, ue_mm_context->bearer_contexts[bix]->ebi, pid, bidx);

          if (rc != RETURNok) {
            break;
          }
        }
      }
    }

    OAILOG_FUNC_RETURN(LOG_NAS_ESM, rc);
  }

  OAILOG_INFO(
    LOG_NAS_ESM,
    "ESM-PROC  - EPS bearer context deactivation "
    "(ue_id=" MME_UE_S1AP_ID_FMT ", ebi=%d)\n",
    ue_mm_context->mme_ue_s1ap_id,
    ebi);

  if ((ue_mm_context) && (*pid < MAX_APN_PER_UE)) {
    if (ue_mm_context->pdn_contexts[*pid] == NULL) {
      OAILOG_ERROR(
        LOG_NAS_ESM,
        "ESM-PROC  - PDN connection %d has not been "
        "allocated\n",
        *pid);
      *esm_cause = ESM_CAUSE_PROTOCOL_ERROR;
    } else {
      int i;

      *esm_cause = ESM_CAUSE_INVALID_EPS_BEARER_IDENTITY;

      for (i = 0; i < BEARERS_PER_UE; i++) {
        if (
          (ue_mm_context->pdn_contexts[*pid]->bearer_contexts[i] <= 0) ||
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
        rc = RETURNok;
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
  const bool is_standalone,
  emm_context_t *const ue_context,
  const ebi_t ebi,
  STOLEN_REF bstring *msg,
  const bool ue_triggered)
{
  OAILOG_FUNC_IN(LOG_NAS_ESM);
  int rc;
  mme_ue_s1ap_id_t ue_id =
    PARENT_STRUCT(ue_context, struct ue_mm_context_s, emm_context)
      ->mme_ue_s1ap_id;

  OAILOG_INFO(
    LOG_NAS_ESM,
    "ESM-PROC  - Initiate EPS bearer context deactivation "
    "(ue_id=" MME_UE_S1AP_ID_FMT ", ebi=%d)\n",
    ue_id,
    ebi);
  /*
   * Send deactivate EPS bearer context request message and
   * * * * start timer T3495
   */
  /*Currently we only support single bearear deactivation at NAS*/
  rc = _eps_bearer_deactivate(ue_context, ebi, msg);
  msg = NULL;

  if (rc != RETURNerror) {
    /*
     * Set the EPS bearer context state to ACTIVE PENDING
     */
    rc = esm_ebr_set_status(
      ue_context, ebi, ESM_EBR_INACTIVE_PENDING, ue_triggered);

    if (rc != RETURNok) {
      /*
       * The EPS bearer context was already in ACTIVE state
       */
      OAILOG_WARNING(
        LOG_NAS_ESM, "ESM-PROC  - EBI %d was already INACTIVE PENDING\n", ebi);
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
  emm_context_t *ue_context,
  ebi_t ebi,
  esm_cause_t *esm_cause)
{
  OAILOG_FUNC_IN(LOG_NAS_ESM);
  int rc = RETURNerror;
  pdn_cid_t pid = MAX_APN_PER_UE;
  mme_ue_s1ap_id_t ue_id =
    PARENT_STRUCT(ue_context, struct ue_mm_context_s, emm_context)
      ->mme_ue_s1ap_id;
  bool delete_default_bearer = false;
  int bid = BEARERS_PER_UE;
  teid_t s_gw_teid_s11_s4 = 0;

  OAILOG_INFO(
    LOG_NAS_ESM,
    "ESM-PROC  - EPS bearer context deactivation "
    "accepted by the UE (ue_id=" MME_UE_S1AP_ID_FMT ", ebi=%d)\n",
    ue_id,
    ebi);
  /*
   * Stop T3495 timer if running
   */
  rc = esm_ebr_stop_timer(ue_context, ebi);

  if (rc != RETURNerror) {

    /*
     * Release the EPS bearer context
     */
    rc = _eps_bearer_release(ue_context, ebi, &pid, &bid);

    if (rc != RETURNok) {
      /*
       * Failed to release the EPS bearer context
       */
      *esm_cause = ESM_CAUSE_PROTOCOL_ERROR;
      pid = RETURNerror;
    }
  }

  s_gw_teid_s11_s4 =
    PARENT_STRUCT(ue_context, struct ue_mm_context_s, emm_context)
    ->pdn_contexts[pid]->s_gw_teid_s11_s4;

  //If bearer id == 0, default bearer is deleted
  if (PARENT_STRUCT(ue_context, struct ue_mm_context_s, emm_context)
    ->pdn_contexts[pid]->default_ebi == ebi) {
    delete_default_bearer = true;
    //Release the default bearer
    rc= mme_api_unsubscribe(NULL);

    if (rc != RETURNerror) {
      /*
       * Delete the PDN connection entry
       */
      _pdn_connectivity_delete(ue_context, pid);
    }
  } else {
    OAILOG_INFO(
      LOG_NAS_ESM,
      "ESM-PROC  - Removing dedicated bearer context "
      "for UE (ue_id=" MME_UE_S1AP_ID_FMT ", ebi=%d)\n",
      ue_id,
      ebi);

    ue_mm_context_t *ue_mm_context =
      PARENT_STRUCT(ue_context, struct ue_mm_context_s, emm_context);
    //Remove dedicated bearer context
    free_wrapper ((void**)&ue_mm_context->bearer_contexts[bid]);
  }
  //Send deactivate_eps_bearer_context to MME APP
  nas_itti_deactivate_eps_bearer_context(
    ue_id,
    ebi,
    delete_default_bearer,
    s_gw_teid_s11_s4);

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
 ** Name:    _eps_bearer_deactivate_t3495_handler()                    **
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
static void _eps_bearer_deactivate_t3495_handler(void *args)
{
  OAILOG_FUNC_IN(LOG_NAS_ESM);
  int rc;
  bool delete_default_bearer = false;
  /*
   * Get retransmission timer parameters data
   */
  esm_ebr_timer_data_t *esm_ebr_timer_data = (esm_ebr_timer_data_t *) (args);

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
      esm_ebr_timer_data->ue_id,
      esm_ebr_timer_data->ebi,
      esm_ebr_timer_data->count);

    if (esm_ebr_timer_data->count < EPS_BEARER_DEACTIVATE_COUNTER_MAX) {
      /*
       * Re-send deactivate EPS bearer context request message to the UE
       */
      bstring b = bstrcpy(esm_ebr_timer_data->msg);
      rc = _eps_bearer_deactivate(
        esm_ebr_timer_data->ctx, esm_ebr_timer_data->ebi, &b);
    } else {
      /*
       * The maximum number of deactivate EPS bearer context request
       * message retransmission has exceed
       */
      pdn_cid_t pid = MAX_APN_PER_UE;
      int bid = BEARERS_PER_UE;

      /*
       * Deactivate the EPS bearer context locally without peer-to-peer
       * * * * signalling between the UE and the MME
       */
      rc = _eps_bearer_release(
        esm_ebr_timer_data->ctx, esm_ebr_timer_data->ebi, &pid, &bid);

      if (rc != RETURNerror) {
        /*
         * Stop timer T3495
         */
        rc =
          esm_ebr_stop_timer(esm_ebr_timer_data->ctx, esm_ebr_timer_data->ebi);
      }

      //Send bearer_deactivation_reject to MME
      teid_t s_gw_teid_s11_s4 =
        PARENT_STRUCT(esm_ebr_timer_data->ctx,
          struct ue_mm_context_s, emm_context)
        ->pdn_contexts[pid]->s_gw_teid_s11_s4;

      if (PARENT_STRUCT(esm_ebr_timer_data->ctx,
        struct ue_mm_context_s, emm_context)
        ->pdn_contexts[pid]->default_ebi == esm_ebr_timer_data->ebi) {
        delete_default_bearer = true;
        //Release the default bearer
        /*
         * Delete the PDN connection entry
         */
        _pdn_connectivity_delete(esm_ebr_timer_data->ctx, pid);
      }
      nas_itti_dedicated_eps_bearer_deactivation_reject(
        esm_ebr_timer_data->ue_id,
        esm_ebr_timer_data->ebi,
        delete_default_bearer,
        s_gw_teid_s11_s4);
    }
    if (esm_ebr_timer_data->msg) {
      bdestroy_wrapper(&esm_ebr_timer_data->msg);
    }
    free_wrapper((void **) &esm_ebr_timer_data);
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
static int _eps_bearer_deactivate(
  emm_context_t *ue_context,
  ebi_t ebi,
  STOLEN_REF bstring *msg)
{
  OAILOG_FUNC_IN(LOG_NAS_ESM);
  emm_sap_t emm_sap = {0};
  int rc;
  mme_ue_s1ap_id_t ue_id =
    PARENT_STRUCT(ue_context, struct ue_mm_context_s, emm_context)
      ->mme_ue_s1ap_id;

  /*
   * Notify EMM that a deactivate EPS bearer context request message
   * has to be sent to the UE
   */

  emm_sap.primitive = EMMESM_DEACTIVATE_BEARER_REQ;
  emm_sap.u.emm_esm.ue_id = ue_id;
  emm_sap.u.emm_esm.ctx = ue_context;
  emm_sap.u.emm_esm.u.deactivate_bearer.ebi = ebi;
  emm_sap.u.emm_esm.u.deactivate_bearer.msg = *msg;
  bstring msg_dup = bstrcpy(*msg);
  *msg = NULL;
  rc = emm_sap_send(&emm_sap);

  if (rc != RETURNerror) {
    /*
     * Start T3495 retransmission timer
     */
    rc = esm_ebr_start_timer(
      ue_context,
      ebi,
      msg_dup,
      mme_config.nas_config.t3495_sec,
      _eps_bearer_deactivate_t3495_handler);
  } else {
    bdestroy_wrapper(&msg_dup);
  }
  OAILOG_FUNC_RETURN(LOG_NAS_ESM, rc);
}

/****************************************************************************
 **                                                                        **
 ** Name:    _eps_bearer_release()                                     **
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
static int _eps_bearer_release(
  emm_context_t *ue_context,
  ebi_t ebi,
  pdn_cid_t *pid,
  int *bidx)
{
  OAILOG_FUNC_IN(LOG_NAS_ESM);
  int rc = RETURNerror;

  /*
   * Release the EPS bearer context entry
   */
  ebi = esm_ebr_context_release(ue_context, ebi, pid, bidx);

  if (ebi == ESM_EBI_UNASSIGNED) {
    OAILOG_WARNING(
      LOG_NAS_ESM, "ESM-PROC  - Failed to release EPS bearer context\n");
  } else {
    /*
     * Set the EPS bearer context state to INACTIVE
     */
    rc = esm_ebr_set_status(ue_context, ebi, ESM_EBR_INACTIVE, false);

    if (rc != RETURNok) {
      /*
       * The EPS bearer context was already in INACTIVE state
       */
      OAILOG_WARNING(
        LOG_NAS_ESM, "ESM-PROC  - EBI %d was already INACTIVE\n", ebi);
    }
    /*
     * Release EPS bearer data
     */
    rc = esm_ebr_release(ue_context, ebi);

    if (rc != RETURNok) {
      OAILOG_WARNING(
        LOG_NAS_ESM, "ESM-PROC  - Failed to release EPS bearer data\n");
    }
  }

  OAILOG_FUNC_RETURN(LOG_NAS_ESM, rc);
}

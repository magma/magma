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
#include "3gpp_24.008.h"
#include "emm_data.h"
#include "mme_app_ue_context.h"
#include "esm_proc.h"
#include "esm_ebr.h"
#include "esm_ebr_context.h"
#include "emm_sap.h"
#include "mme_config.h"
#include "3gpp_24.301.h"
#include "3gpp_36.401.h"
#include "EsmCause.h"
#include "common_defs.h"
#include "emm_esmDef.h"
#include "esm_data.h"

/****************************************************************************/
/****************  E X T E R N A L    D E F I N I T I O N S  ****************/
/****************************************************************************/

/****************************************************************************/
/*******************  L O C A L    D E F I N I T I O N S  *******************/
/****************************************************************************/

/*
   --------------------------------------------------------------------------
   Internal data handled by the default EPS bearer context activation
   procedure in the MME
   --------------------------------------------------------------------------
*/
/*
   Timer handlers
*/

/* Maximum value of the activate default EPS bearer context request
   retransmission counter */
#define DEFAULT_EPS_BEARER_ACTIVATE_COUNTER_MAX 5

static int default_eps_bearer_activate(
    emm_context_t* emm_context, ebi_t ebi, STOLEN_REF bstring* msg);

static int default_eps_bearer_activate_in_bearer_setup_req(
    emm_context_t* emm_context, ebi_t ebi, STOLEN_REF bstring* msg);

/****************************************************************************/
/******************  E X P O R T E D    F U N C T I O N S  ******************/
/****************************************************************************/

/*
   --------------------------------------------------------------------------
      Default EPS bearer context activation procedure executed by the MME
   --------------------------------------------------------------------------
*/
/****************************************************************************
 **                                                                        **
 ** Name:    esm_proc_default_eps_bearer_context()                     **
 **                                                                        **
 ** Description: Allocates resources required for activation of a default  **
 **      EPS bearer context.                                       **
 **                                                                        **
 ** Inputs:  ue_id:      UE local identifier                        **
 **          pid:       PDN connection identifier                  **
 **      qos:       EPS bearer level QoS parameters            **
 **      Others:    None                                       **
 **                                                                        **
 ** Outputs:     ebi:       EPS bearer identity assigned to the de-    **
 **             fault EPS bearer context                   **
 **      esm_cause: Cause code returned upon ESM procedure     **
 **             failure                                    **
 **      Return:    RETURNok, RETURNerror                      **
 **      Others:    None                                       **
 **                                                                        **
 ***************************************************************************/
int esm_proc_default_eps_bearer_context(
    emm_context_t* emm_context, const proc_tid_t pti, pdn_cid_t pid, ebi_t* ebi,
    const qci_t qci, esm_cause_t* esm_cause) {
  OAILOG_FUNC_IN(LOG_NAS_ESM);
  mme_ue_s1ap_id_t ue_id =
      PARENT_STRUCT(emm_context, struct ue_mm_context_s, emm_context)
          ->mme_ue_s1ap_id;
  OAILOG_INFO(
      LOG_NAS_ESM,
      "ESM-PROC  - Default EPS bearer context activation "
      "(ue_id=" MME_UE_S1AP_ID_FMT ")\n",
      ue_id);
  OAILOG_DEBUG(LOG_NAS_ESM, "(pid=%d,  QCI %u)\n", pid, qci);
  /*
   * Assign new EPS bearer context
   */
  *ebi = esm_ebr_assign(emm_context);

  if (*ebi != ESM_EBI_UNASSIGNED) {
    /*
     * Create default EPS bearer context
     */
    *ebi = esm_ebr_context_create(
        emm_context, pti, pid, *ebi, IS_DEFAULT_BEARER_YES, qci, 0, 0, 0, 0,
        (traffic_flow_template_t*) NULL,
        (protocol_configuration_options_t*) NULL, NULL);

    if (*ebi == ESM_EBI_UNASSIGNED) {
      /*
       * No resource available
       */
      OAILOG_ERROR(
          LOG_NAS_ESM,
          "ESM-PROC  - Failed to create new default EPS "
          "bearer context (ebi=%d) for ue_id " MME_UE_S1AP_ID_FMT "\n",
          *ebi, ue_id);
      *esm_cause = ESM_CAUSE_INSUFFICIENT_RESOURCES;
      OAILOG_FUNC_RETURN(LOG_NAS_ESM, RETURNerror);
    }

    OAILOG_FUNC_RETURN(LOG_NAS_ESM, RETURNok);
  }

  OAILOG_ERROR(
      LOG_NAS_ESM,
      "ESM-PROC  - Failed to assign new EPS bearer context for "
      "ue_id " MME_UE_S1AP_ID_FMT "\n",
      ue_id);
  *esm_cause = ESM_CAUSE_INSUFFICIENT_RESOURCES;
  OAILOG_FUNC_RETURN(LOG_NAS_ESM, RETURNerror);
}

/****************************************************************************
 **                                                                        **
 ** Name:    esm_proc_default_eps_bearer_context_request()             **
 **                                                                        **
 ** Description: Initiates the default EPS bearer context activation pro-  **
 **      cedure                                                    **
 **                                                                        **
 **      3GPP TS 24.301, section 6.4.1.2                           **
 **      The MME initiates the default EPS bearer context activa-  **
 **      tion procedure by sending an ACTIVATE DEFAULT EPS BEARER  **
 **      CONTEXT REQUEST message, starting timer T3485 and ente-   **
 **      ring state BEARER CONTEXT ACTIVE PENDING.                 **
 **                                                                        **
 ** Inputs:  is_standalone: Indicate whether the default bearer is     **
 **             activated as part of the attach procedure  **
 **             or as the response to a stand-alone PDN    **
 **             CONNECTIVITY REQUEST message               **
 **      ue_id:      UE lower layer identifier                  **
 **      ebi:       EPS bearer identity                        **
 **      msg:       Encoded ESM message to be sent             **
 **      ue_triggered:  true if the EPS bearer context procedure   **
 **             was triggered by the UE                    **
 **      Others:    None                                       **
 **                                                                        **
 ** Outputs:     None                                                      **
 **      Return:    RETURNok, RETURNerror                      **
 **      Others:    None                                       **
 **                                                                        **
 ***************************************************************************/
int esm_proc_default_eps_bearer_context_request(
    bool is_standalone, emm_context_t* emm_context, ebi_t ebi,
    STOLEN_REF bstring* msg, bool ue_triggered) {
  OAILOG_FUNC_IN(LOG_NAS_ESM);
  int rc = RETURNok;

  mme_ue_s1ap_id_t ue_id =
      PARENT_STRUCT(emm_context, struct ue_mm_context_s, emm_context)
          ->mme_ue_s1ap_id;

  if (is_standalone) {
    /*
     * Send activate default EPS bearer context request message and
     * * * * start timer T3485
     */
    OAILOG_INFO(
        LOG_NAS_ESM,
        "ESM-PROC  - Initiate standalone default EPS bearer context activation "
        "(ue_id=" MME_UE_S1AP_ID_FMT ", ebi=%d)\n",
        ue_id, ebi);

    /* Send ACTIVATE DEFAULT EPS BEARER CONTEXT REQUEST
     * in ERAB SETUP REQ mesage
     */
    rc = default_eps_bearer_activate_in_bearer_setup_req(emm_context, ebi, msg);
  } else {
    OAILOG_INFO(
        LOG_NAS_ESM,
        "ESM-PROC  - Initiate non standalone default EPS bearer context "
        "activation "
        "(ue_id=" MME_UE_S1AP_ID_FMT ", ebi=%d)\n",
        ue_id, ebi);
  }

  if (rc != RETURNerror) {
    /*
     * Set the EPS bearer context state to ACTIVE PENDING
     */
    rc = esm_ebr_set_status(
        emm_context, ebi, ESM_EBR_ACTIVE_PENDING, ue_triggered);

    if (rc != RETURNok) {
      /*
       * The EPS bearer context was already in ACTIVE PENDING state
       */
      OAILOG_WARNING(
          LOG_NAS_ESM,
          "ESM-PROC  - EBI %d was already ACTIVE PENDING for ue "
          "id " MME_UE_S1AP_ID_FMT "\n",
          ebi, ue_id);
    }
  }

  OAILOG_FUNC_RETURN(LOG_NAS_ESM, rc);
}

/****************************************************************************
 **                                                                        **
 ** Name:    esm_proc_default_eps_bearer_context_accept()              **
 **                                                                        **
 ** Description: Performs default EPS bearer context activation procedure  **
 **      accepted by the UE.                                       **
 **                                                                        **
 **      3GPP TS 24.301, section 6.4.1.3                           **
 **      Upon receipt of the ACTIVATE DEFAULT EPS BEARER CONTEXT   **
 **      ACCEPT message, the MME shall enter the state BEARER CON- **
 **      TEXT ACTIVE and stop the timer T3485, if it is running.   **
 **                                                                        **
 ** Inputs:  ue_id:      UE local identifier                        **
 **      ebi:       EPS bearer identity                        **
 **      Others:    None                                       **
 **                                                                        **
 ** Outputs:     esm_cause: Cause code returned upon ESM procedure     **
 **             failure                                    **
 **      Return:    RETURNok, RETURNerror                      **
 **      Others:    None                                       **
 **                                                                        **
 ***************************************************************************/
int esm_proc_default_eps_bearer_context_accept(
    emm_context_t* emm_context, ebi_t ebi, esm_cause_t* esm_cause) {
  OAILOG_FUNC_IN(LOG_NAS_ESM);
  int rc;
  mme_ue_s1ap_id_t ue_id =
      PARENT_STRUCT(emm_context, struct ue_mm_context_s, emm_context)
          ->mme_ue_s1ap_id;

  OAILOG_INFO(
      LOG_NAS_ESM,
      "ESM-PROC  - Default EPS bearer context activation "
      "accepted by the UE (ue_id=" MME_UE_S1AP_ID_FMT ", ebi=%d)\n",
      ue_id, ebi);
  /*
   * Stop T3485 timer if running
   */
  rc = esm_ebr_stop_timer(emm_context, ebi);

  if (rc != RETURNerror) {
    /*
     * Set the EPS bearer context state to ACTIVE
     */
    rc = esm_ebr_set_status(emm_context, ebi, ESM_EBR_ACTIVE, false);

    if (rc != RETURNok) {
      /*
       * The EPS bearer context was already in ACTIVE state
       */
      OAILOG_WARNING(
          LOG_NAS_ESM,
          "ESM-PROC  - EBI %d was already ACTIVE for ue id " MME_UE_S1AP_ID_FMT
          "\n",
          ebi, ue_id);
      *esm_cause = ESM_CAUSE_PROTOCOL_ERROR;
    }
  }

  OAILOG_FUNC_RETURN(LOG_NAS_ESM, rc);
}

/****************************************************************************
 **                                                                        **
 ** Name:    esm_proc_default_eps_bearer_context_reject()              **
 **                                                                        **
 ** Description: Performs default EPS bearer context activation procedure  **
 **      not accepted by the UE.                                   **
 **                                                                        **
 **      3GPP TS 24.301, section 6.4.1.4                           **
 **      Upon receipt of the ACTIVATE DEFAULT EPS BEARER CONTEXT   **
 **      REJECT message, the MME shall enter the state BEARER CON- **
 **      TEXT INACTIVE and stop the timer T3485, if it is running. **
 **                                                                        **
 ** Inputs:  ue_id:      UE local identifier                        **
 **      ebi:       EPS bearer identity                        **
 **      Others:    None                                       **
 **                                                                        **
 ** Outputs:     esm_cause: Cause code returned upon ESM procedure     **
 **             failure                                    **
 **      Return:    RETURNok, RETURNerror                      **
 **      Others:    None                                       **
 **                                                                        **
 ***************************************************************************/
int esm_proc_default_eps_bearer_context_reject(
    emm_context_t* emm_context, ebi_t ebi, esm_cause_t* esm_cause) {
  OAILOG_FUNC_IN(LOG_NAS_ESM);
  int rc;
  mme_ue_s1ap_id_t ue_id =
      PARENT_STRUCT(emm_context, struct ue_mm_context_s, emm_context)
          ->mme_ue_s1ap_id;

  OAILOG_WARNING(
      LOG_NAS_ESM,
      "ESM-PROC  - Default EPS bearer context activation "
      "not accepted by the UE (ue_id=" MME_UE_S1AP_ID_FMT ", ebi=%d)\n",
      ue_id, ebi);
  /*
   * Stop T3485 timer if running
   */
  rc = esm_ebr_stop_timer(emm_context, ebi);

  if (rc != RETURNerror) {
    pdn_cid_t pid = MAX_APN_PER_UE;
    int bid       = BEARERS_PER_UE;

    /*
     * Release the default EPS bearer context and enter state INACTIVE
     */
    rc = esm_proc_eps_bearer_context_deactivate(
        emm_context, true, ebi, &pid, &bid, NULL);

    if (rc != RETURNok) {
      /*
       * Failed to release the default EPS bearer context
       */
      *esm_cause = ESM_CAUSE_PROTOCOL_ERROR;
    }

#if EMBEDDED_SGW
    ue_mm_context_t* ue_context_p =
        PARENT_STRUCT(emm_context, struct ue_mm_context_s, emm_context);
    if (!ue_context_p) {
      OAILOG_ERROR(
          LOG_NAS_ESM,
          "ESM-PROC  - ue_context_p is NULL for ue_id" MME_UE_S1AP_ID_FMT
          ", ebi %u\n",
          ue_id, ebi);
      OAILOG_FUNC_RETURN(LOG_NAS_ESM, RETURNerror);
    }
    // Send delete session req to spgw
    mme_app_send_delete_session_request(ue_context_p, ebi, pid);
    /* Set ue_rej_act_def_ber_req flag in order to delete the PDN session
     * after receiving delete session rsp from spgw
     */
    if (ue_context_p->pdn_contexts[pid]) {
      ue_context_p->pdn_contexts[pid]->ue_rej_act_def_ber_req = true;
    }
#endif
  }
  OAILOG_FUNC_RETURN(LOG_NAS_ESM, rc);
}

/****************************************************************************
 **                                                                        **
 ** Name:    esm_proc_default_eps_bearer_context_failure()             **
 **                                                                        **
 ** Description: Performs default EPS bearer context activation procedure  **
 **      upon receiving notification from the EPS Mobility Manage- **
 **      ment sublayer that EMM procedure that initiated EPS de-   **
 **      fault bearer context activation locally failed.           **
 **                                                                        **
 **      The MME releases the default EPS bearer context previous- **
 **      ly allocated when ACTIVATE DEFAULT EPS BEARER CONTEXT RE- **
 **      QUEST message was received.                               **
 **                                                                        **
 ** Inputs:  ue_id:      UE local identifier                        **
 **      Others:    None                                       **
 **                                                                        **
 ** Outputs:     None                                                      **
 **      Return:    The identifier of the PDN connection the   **
 **             default EPS bearer context belongs to if   **
 **             successfully released;                     **
 **             RETURNerror  otherwise.                    **
 **      Others:    None                                       **
 **                                                                        **
 ***************************************************************************/
int esm_proc_default_eps_bearer_context_failure(
    emm_context_t* emm_context, pdn_cid_t* const pid) {
  OAILOG_FUNC_IN(LOG_NAS_ESM);
  int rc = RETURNerror;
  mme_ue_s1ap_id_t ue_id =
      PARENT_STRUCT(emm_context, struct ue_mm_context_s, emm_context)
          ->mme_ue_s1ap_id;

  *pid = MAX_APN_PER_UE;

  if (emm_context) {
    OAILOG_WARNING(
        LOG_NAS_ESM,
        "ESM-PROC  - Default EPS bearer context activation "
        "failure (ue_id=" MME_UE_S1AP_ID_FMT ")\n",
        ue_id);
  } else {
    OAILOG_WARNING(
        LOG_NAS_ESM,
        "ESM-PROC  - Default EPS bearer context activation "
        "failure (context is NULL) for ue id " MME_UE_S1AP_ID_FMT "\n",
        ue_id);
    OAILOG_FUNC_RETURN(LOG_NAS_ESM, RETURNerror);
  }
  /*
   * Get the EPS bearer identity of the EPS bearer context which is still
   * * * * pending in the active pending state
   */
  ebi_t ebi = esm_ebr_get_pending_ebi(emm_context, ESM_EBR_ACTIVE_PENDING);

  if (ebi != ESM_EBI_UNASSIGNED) {
    int bid = BEARERS_PER_UE;

    /*
     * Release the default EPS bearer context and enter state INACTIVE
     */
    rc = esm_proc_eps_bearer_context_deactivate(
        emm_context, true, ebi, pid, &bid, NULL);
  }

  OAILOG_FUNC_RETURN(LOG_NAS_ESM, rc);
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
 ** Name:    default_eps_bearer_activate_t3485_handler()              **
 **                                                                        **
 ** Description: T3485 timeout handler                                     **
 **                                                                        **
 **              3GPP TS 24.301, section 6.4.1.6, case a                   **
 **      On the first expiry of the timer T3485, the MME shall re- **
 **      send the ACTIVATE DEFAULT EPS BEARER CONTEXT REQUEST and  **
 **      shall reset and restart timer T3485. This retransmission  **
 **      is repeated four times, i.e. on the fifth expiry of timer **
 **      T3485, the MME shall release possibly allocated resources **
 **      for this activation and shall abort the procedure.        **
 **                                                                        **
 ** Inputs:  args:      handler parameters                         **
 **      Others:    None                                       **
 **                                                                        **
 ** Outputs:     None                                                      **
 **      Return:    None                                       **
 **      Others:    None                                       **
 **                                                                        **
 ***************************************************************************/
void default_eps_bearer_activate_t3485_handler(void* args, imsi64_t* imsi64) {
  OAILOG_FUNC_IN(LOG_NAS_ESM);

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
        "ESM-PROC  - T3485 timer expired (ue_id=" MME_UE_S1AP_ID_FMT
        ", ebi=%d), "
        "retransmission counter = %d\n",
        esm_ebr_timer_data->ue_id, esm_ebr_timer_data->ebi,
        esm_ebr_timer_data->count);

    *imsi64 = esm_ebr_timer_data->ctx->_imsi64;
    if (esm_ebr_timer_data->count < DEFAULT_EPS_BEARER_ACTIVATE_COUNTER_MAX) {
      /*
       * Re-send activate default EPS bearer context request message
       * * * * to the UE
       */
      bstring b = bstrcpy(esm_ebr_timer_data->msg);

      /* If standalone PDN connectivity request is received send activate
       * default eps bearer req message in erab setup req s1ap message.
       * If PDN connectivity is received along with attach req send
       * activate default eps bearer req message in ICS req
       */
      ue_mm_context_t* ue_context_p = PARENT_STRUCT(
          ((emm_context_t*) esm_ebr_timer_data->ctx), struct ue_mm_context_s,
          emm_context);
      bearer_context_t* bc =
          mme_app_get_bearer_context(ue_context_p, esm_ebr_timer_data->ebi);
      if (((emm_context_t*) esm_ebr_timer_data->ctx)
              ->esm_ctx.pending_standalone &&
          (!(bc->enb_fteid_s1u.teid))) {
        default_eps_bearer_activate_in_bearer_setup_req(
            esm_ebr_timer_data->ctx, esm_ebr_timer_data->ebi, &b);
      } else {
        default_eps_bearer_activate(
            esm_ebr_timer_data->ctx, esm_ebr_timer_data->ebi, &b);
      }
      bdestroy_wrapper(&b);
    } else {
      /*
       * The maximum number of activate default EPS bearer context request
       * message retransmission has exceed
       */
      pdn_cid_t pid = MAX_APN_PER_UE;
      int bidx      = BEARERS_PER_UE;

      if (((emm_context_t*) esm_ebr_timer_data->ctx)
              ->esm_ctx.pending_standalone > 0) {
        ((emm_context_t*) esm_ebr_timer_data->ctx)
            ->esm_ctx.pending_standalone -= 1;
      }
      /*
       * Release the default EPS bearer context and enter state INACTIVE
       */
      esm_proc_eps_bearer_context_deactivate(
          esm_ebr_timer_data->ctx, true, esm_ebr_timer_data->ebi, &pid, &bidx,
          NULL);
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
 ** Name:    _default_eps_bearer_activate()                            **
 **                                                                        **
 ** Description: Sends ACTIVATE DEFAULT EPS BEREAR CONTEXT REQUEST message **
 **      and starts timer T3485                                    **
 **                                                                        **
 ** Inputs:  ue_id:      UE local identifier                        **
 **      ebi:       EPS bearer identity                        **
 **      msg:       Encoded ESM message to be sent             **
 **      Others:    None                                       **
 **                                                                        **
 ** Outputs:     None                                                      **
 **      Return:    RETURNok, RETURNerror                      **
 **      Others:    T3485                                      **
 **                                                                        **
 ***************************************************************************/
static int default_eps_bearer_activate(
    emm_context_t* emm_context, ebi_t ebi, STOLEN_REF bstring* msg) {
  OAILOG_FUNC_IN(LOG_NAS_ESM);
  emm_sap_t emm_sap = {0};
  int rc;
  mme_ue_s1ap_id_t ue_id =
      PARENT_STRUCT(emm_context, struct ue_mm_context_s, emm_context)
          ->mme_ue_s1ap_id;

  /*
   * Notify EMM that an activate default EPS bearer context request message
   * has to be sent to the UE
   */
  emm_esm_data_t* emm_esm = &emm_sap.u.emm_esm.u.data;

  emm_sap.primitive       = EMMESM_UNITDATA_REQ;
  emm_sap.u.emm_esm.ue_id = ue_id;
  emm_sap.u.emm_esm.ctx   = emm_context;
  emm_esm->msg            = *msg;
  rc                      = emm_sap_send(&emm_sap);

  if (rc != RETURNerror) {
    /*
     * Start T3485 retransmission timer
     */
    rc = esm_ebr_start_timer(
        emm_context, ebi, *msg, mme_config.nas_config.t3485_sec,
        default_eps_bearer_activate_t3485_handler);
    if (rc != RETURNerror) {
      OAILOG_DEBUG(
          LOG_NAS_ESM,
          "ESM-PROC  - Started t3485 for ue_id=" MME_UE_S1AP_ID_FMT "\n",
          ue_id);
    }
  }

  OAILOG_FUNC_RETURN(LOG_NAS_ESM, rc);
}

/****************************************************************************
 **                                                                        **
 ** Name:    default_eps_bearer_activate_in_bearer_setup_req()           **
 **                                                                        **
 ** Description: Sends ACTIVATE DEFAULT EPS BEREAR CONTEXT REQUEST message **
 ** in ERAB_REQ message and starts timer T3485                             **
 **                                                                        **
 ** Inputs: ebi:    EPS bearer identity                                    **
 **      msg:       Encoded ESM message to be sent                         **
 **      Others:    None                                                   **
 **                                                                        **
 ** Outputs:     None                                                      **
 **      Return:    RETURNok, RETURNerror                                  **
 **      Others:    T3485                                                  **
 **                                                                        **
 ***************************************************************************/
static int default_eps_bearer_activate_in_bearer_setup_req(
    emm_context_t* emm_context, ebi_t ebi, STOLEN_REF bstring* msg) {
  OAILOG_FUNC_IN(LOG_NAS_ESM);
  emm_sap_t emm_sap = {0};
  int rc;
  mme_ue_s1ap_id_t ue_id =
      PARENT_STRUCT(emm_context, struct ue_mm_context_s, emm_context)
          ->mme_ue_s1ap_id;

  bearer_context_t* bearer_context = mme_app_get_bearer_context(
      PARENT_STRUCT(emm_context, struct ue_mm_context_s, emm_context), ebi);
  /*
   * Notify EMM that an activate default EPS bearer context request message
   * has to be sent to the UE
   */
  emm_esm_activate_bearer_req_t* emm_esm_activate =
      &emm_sap.u.emm_esm.u.activate_bearer;

  emm_sap.primitive       = EMMESM_ACTIVATE_BEARER_REQ;
  emm_sap.u.emm_esm.ue_id = ue_id;
  emm_sap.u.emm_esm.ctx   = emm_context;
  emm_esm_activate->msg   = *msg;
  emm_esm_activate->ebi   = ebi;

  emm_esm_activate->mbr_dl = bearer_context->esm_ebr_context.mbr_dl;
  emm_esm_activate->mbr_ul = bearer_context->esm_ebr_context.mbr_ul;
  emm_esm_activate->gbr_dl = bearer_context->esm_ebr_context.gbr_dl;
  emm_esm_activate->gbr_ul = bearer_context->esm_ebr_context.gbr_ul;

  bstring msg_dup = bstrcpy(*msg);

  rc = emm_sap_send(&emm_sap);

  if (rc != RETURNerror) {
    /*
     * Start T3485 retransmission timer
     */
    rc = esm_ebr_start_timer(
        emm_context, ebi, msg_dup, mme_config.nas_config.t3485_sec,
        default_eps_bearer_activate_t3485_handler);
  }

  bdestroy_wrapper(&msg_dup);
  OAILOG_FUNC_RETURN(LOG_NAS_ESM, rc);
}

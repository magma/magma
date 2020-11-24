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
#include "mme_app_defs.h"

/****************************************************************************/
/****************  E X T E R N A L    D E F I N I T I O N S  ****************/
/****************************************************************************/

/****************************************************************************/
/*******************  L O C A L    D E F I N I T I O N S  *******************/
/****************************************************************************/

/*
   --------------------------------------------------------------------------
   Internal data handled by the dedicated EPS bearer context activation
   procedure in the MME
   --------------------------------------------------------------------------
*/
/*
   Timer handlers
*/
static void _dedicated_eps_bearer_activate_t3485_handler(
    void*, imsi64_t* imsi64);

/* Maximum value of the activate dedicated EPS bearer context request
   retransmission counter */
#define DEDICATED_EPS_BEARER_ACTIVATE_COUNTER_MAX 5

static int _dedicated_eps_bearer_activate(
    emm_context_t* emm_context, ebi_t ebi, STOLEN_REF bstring* msg);

static void _erab_setup_rsp_tmr_exp_ded_bearer_handler(
    void* args, imsi64_t* imsi64);

/****************************************************************************/
/******************  E X P O R T E D    F U N C T I O N S  ******************/
/****************************************************************************/

/*
   --------------------------------------------------------------------------
      Dedicated EPS bearer context activation procedure executed by the MME
   --------------------------------------------------------------------------
*/
/****************************************************************************
 **                                                                        **
 ** Name:    esm_proc_dedicated_eps_bearer_context()                       **
 **                                                                        **
 ** Description: Allocates resources required for activation of a dedica-  **
 **      ted EPS bearer context.                                           **
 **                                                                        **
 ** Inputs:  ue_id:      UE local identifier                               **
 **          pid:       PDN connection identifier                          **
 **      esm_qos:   EPS bearer level QoS parameters                        **
 **      tft:       Traffic flow template parameters                       **
 **      Others:    None                                                   **
 **                                                                        **
 ** Outputs:     ebi:       EPS bearer identity assigned to the new        **
 **             dedicated bearer context                                   **
 **      default_ebi:   EPS bearer identity of the associated de-          **
 **             fault EPS bearer context                                   **
 **      esm_cause: Cause code returned upon ESM procedure                 **
 **             failure                                                    **
 **      Return:    RETURNok, RETURNerror                                  **
 **      Others:    None                                                   **
 **                                                                        **
 ***************************************************************************/
int esm_proc_dedicated_eps_bearer_context(
    emm_context_t* emm_context, const proc_tid_t pti, pdn_cid_t pid, ebi_t* ebi,
    ebi_t* default_ebi, const qci_t qci, const bitrate_t gbr_dl,
    const bitrate_t gbr_ul, const bitrate_t mbr_dl, const bitrate_t mbr_ul,
    traffic_flow_template_t* tft, protocol_configuration_options_t* pco,
    fteid_t* sgw_fteid, esm_cause_t* esm_cause) {
  OAILOG_FUNC_IN(LOG_NAS_ESM);
  mme_ue_s1ap_id_t ue_id =
      PARENT_STRUCT(emm_context, struct ue_mm_context_s, emm_context)
          ->mme_ue_s1ap_id;
  OAILOG_INFO(
      LOG_NAS_ESM,
      "ESM-PROC  - Dedicated EPS bearer context activation "
      "(ue_id=" MME_UE_S1AP_ID_FMT ", pid=%d)\n",
      ue_id, pid);
  /*
   * Assign new EPS bearer context
   */
  if (*ebi == ESM_EBI_UNASSIGNED) {
    *ebi = esm_ebr_assign(emm_context);
  }

  if (*ebi != ESM_EBI_UNASSIGNED) {
    /*
     * Create dedicated EPS bearer context
     */
    *default_ebi = esm_ebr_context_create(
        emm_context, pti, pid, *ebi, IS_DEFAULT_BEARER_NO, qci, gbr_dl, gbr_ul,
        mbr_dl, mbr_ul, tft, pco, sgw_fteid);

    if (*default_ebi == ESM_EBI_UNASSIGNED) {
      /*
       * No resource available
       */
      OAILOG_WARNING(
          LOG_NAS_ESM,
          "ESM-PROC  - Failed to create dedicated EPS "
          "bearer context (ebi=%d)\n",
          *ebi);
      *esm_cause = ESM_CAUSE_INSUFFICIENT_RESOURCES;
      OAILOG_FUNC_RETURN(LOG_NAS_ESM, RETURNerror);
    }

    OAILOG_FUNC_RETURN(LOG_NAS_ESM, RETURNok);
  }

  OAILOG_WARNING(
      LOG_NAS_ESM, "ESM-PROC  - Failed to assign new EPS bearer context\n");
  *esm_cause = ESM_CAUSE_INSUFFICIENT_RESOURCES;
  OAILOG_FUNC_RETURN(LOG_NAS_ESM, RETURNerror);
}

/****************************************************************************
 **                                                                        **
 ** Name:    esm_proc_dedicated_eps_bearer_context_request()               **
 **                                                                        **
 ** Description: Initiates the dedicated EPS bearer context activation pro-**
 **      cedure                                                            **
 **                                                                        **
 **      3GPP TS 24.301, section 6.4.2.2                                   **
 **      The MME initiates the dedicated EPS bearer context activa-        **
 **      tion procedure by sending an ACTIVATE DEDICATED EPS BEA-          **
 **      RER CONTEXT REQUEST message, starting timer T3485 and en-         **
 **      tering state BEARER CONTEXT ACTIVE PENDING.                       **
 **                                                                        **
 ** Inputs:  is_standalone: Not used (always true)                         **
 **      ue_id:     UE lower layer identifier                              **
 **      ebi:       EPS bearer identity                                    **
 **      msg:       Encoded ESM message to be sent                         **
 **      ue_triggered:  true if the EPS bearer context procedure           **
 **             was triggered by the UE                                    **
 **      Others:    None                                                   **
 **                                                                        **
 ** Outputs:     None                                                      **
 **      Return:    RETURNok, RETURNerror                                  **
 **      Others:    None                                                   **
 **                                                                        **
 ***************************************************************************/
int esm_proc_dedicated_eps_bearer_context_request(
    bool is_standalone, emm_context_t* emm_context, ebi_t ebi,
    STOLEN_REF bstring* msg, bool ue_triggered) {
  OAILOG_FUNC_IN(LOG_NAS_ESM);
  int rc = RETURNok;
  mme_ue_s1ap_id_t ue_id =
      PARENT_STRUCT(emm_context, struct ue_mm_context_s, emm_context)
          ->mme_ue_s1ap_id;

  OAILOG_INFO(
      LOG_NAS_ESM,
      "ESM-PROC  - Initiate dedicated EPS bearer context "
      "activation (ue_id=" MME_UE_S1AP_ID_FMT ", ebi=%d)\n",
      ue_id, ebi);
  /*
   * Send activate dedicated EPS bearer context request message and
   * * * * start timer T3485
   */
  rc = _dedicated_eps_bearer_activate(emm_context, ebi, msg);

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
          LOG_NAS_ESM, "ESM-PROC  - EBI %d was already ACTIVE PENDING\n", ebi);
    }
  }

  OAILOG_FUNC_RETURN(LOG_NAS_ESM, rc);
}

/****************************************************************************
 **                                                                        **
 ** Name:    esm_proc_dedicated_eps_bearer_context_accept()                **
 **                                                                        **
 ** Description: Performs dedicated EPS bearer context activation procedu- **
 **      re accepted by the UE.                                            **
 **                                                                        **
 **      3GPP TS 24.301, section 6.4.2.3                                   **
 **      Upon receipt of the ACTIVATE DEDICATED EPS BEARER CONTEXT         **
 **      ACCEPT message, the MME shall stop the timer T3485 and            **
 **      enter the state BEARER CONTEXT ACTIVE.                            **
 **                                                                        **
 ** Inputs:  ue_id:      UE local identifier                               **
 **      ebi:       EPS bearer identity                                    **
 **      Others:    None                                                   **
 **                                                                        **
 ** Outputs:     esm_cause: Cause code returned upon ESM procedure         **
 **             failure                                                    **
 **      Return:    RETURNok, RETURNerror                                  **
 **      Others:    None                                                   **
 **                                                                        **
 ***************************************************************************/
int esm_proc_dedicated_eps_bearer_context_accept(
    emm_context_t* emm_context, ebi_t ebi, esm_cause_t* esm_cause) {
  OAILOG_FUNC_IN(LOG_NAS_ESM);
  int rc                        = RETURNerror;
  ue_mm_context_t* ue_context_p = NULL;

  ue_context_p =
      PARENT_STRUCT(emm_context, struct ue_mm_context_s, emm_context);
  if (!ue_context_p) {
    OAILOG_ERROR(LOG_NAS_ESM, "Failed to find ue context from emm_context \n");
    OAILOG_FUNC_RETURN(LOG_NAS_ESM, RETURNerror);
  }
  OAILOG_INFO(
      LOG_NAS_ESM,
      "ESM-PROC  - Dedicated EPS bearer context activation "
      "accepted by the UE (ue_id=" MME_UE_S1AP_ID_FMT ", ebi=%u)\n",
      ue_context_p->mme_ue_s1ap_id, ebi);
  /*
   * Stop T3485 timer
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
          LOG_NAS_ESM, "ESM-PROC  - EBI %u was already ACTIVE\n", ebi);
      *esm_cause = ESM_CAUSE_PROTOCOL_ERROR;
    }
    bearer_context_t* bearer_ctx =
        mme_app_get_bearer_context(ue_context_p, ebi);
    /* Send MBR only after receiving ERAB_SETUP_RSP.
     * bearer_ctx->enb_fteid_s1u.teid gets updated after receiving
     * ERAB_SETUP_RSP.*/
    if (bearer_ctx->enb_fteid_s1u.teid) {
      mme_app_handle_create_dedicated_bearer_rsp(ue_context_p, ebi);
    } else {
      rc = esm_ebr_start_timer(
          emm_context, ebi, NULL, ERAB_SETUP_RSP_TMR,
          _erab_setup_rsp_tmr_exp_ded_bearer_handler);
      if (rc != RETURNerror) {
        OAILOG_DEBUG(
            LOG_NAS_ESM,
            "ESM-PROC  - Started ERAB_SETUP_RSP_TMR for "
            "ue_id=" MME_UE_S1AP_ID_FMT "ebi (%u)",
            ue_context_p->mme_ue_s1ap_id, ebi);
      }
    }
  }

  OAILOG_FUNC_RETURN(LOG_NAS_ESM, rc);
}

/****************************************************************************
 **                                                                        **
 ** Name:    esm_proc_dedicated_eps_bearer_context_reject()                **
 **                                                                        **
 ** Description: Performs dedicated EPS bearer context activation procedu- **
 **      re not accepted by the UE.                                        **
 **                                                                        **
 **      3GPP TS 24.301, section 6.4.2.4                                   **
 **      Upon receipt of the ACTIVATE DEDICATED EPS BEARER CONTEXT         **
 **      REJECT message, the MME shall stop the timer T3485, enter         **
 **      the state BEARER CONTEXT INACTIVE and abort the dedicated         **
 **      EPS bearer context activation procedure.                          **
 **      The MME also requests the lower layer to release the ra-          **
 **      dio resources that were established during the dedicated          **
 **      EPS bearer context activation.                                    **
 **                                                                        **
 ** Inputs:  ue_id:      UE local identifier                               **
 **      ebi:       EPS bearer identity                                    **
 **      Others:    None                                                   **
 **                                                                        **
 ** Outputs:     esm_cause: Cause code returned upon ESM procedure         **
 **             failure                                                    **
 **      Return:    RETURNok, RETURNerror                                  **
 **      Others:    None                                                   **
 **                                                                        **
 ***************************************************************************/
int esm_proc_dedicated_eps_bearer_context_reject(
    emm_context_t* emm_context, ebi_t ebi) {
  OAILOG_FUNC_IN(LOG_NAS_ESM);
  int rc;
  ue_mm_context_t* ue_context_p = NULL;

  ue_context_p =
      PARENT_STRUCT(emm_context, struct ue_mm_context_s, emm_context);

  if (!ue_context_p) {
    OAILOG_ERROR(LOG_NAS_ESM, "Failed to find ue context from emm_context \n");
    OAILOG_FUNC_RETURN(LOG_NAS_ESM, RETURNerror);
  }
  OAILOG_INFO(
      LOG_NAS_ESM,
      "ESM-PROC  - Dedicated EPS bearer context activation "
      "not accepted by the UE for ue_id=" MME_UE_S1AP_ID_FMT ", ebi=%u\n",
      ue_context_p->mme_ue_s1ap_id, ebi);
  /*
   * Stop T3485 timer if running
   */
  rc = esm_ebr_stop_timer(emm_context, ebi);

  if (rc != RETURNerror) {
    pdn_cid_t pid = MAX_APN_PER_UE;
    int bid       = BEARERS_PER_UE;

    /*
     * Release the dedicated EPS bearer context and enter state INACTIVE
     */
    rc = esm_proc_eps_bearer_context_deactivate(
        emm_context, true, ebi, &pid, &bid, NULL);

    if (rc != RETURNok) {
      OAILOG_INFO(
          LOG_NAS_ESM,
          "Failed to release the dedicated EPS bearer context for ebi:%u\n",
          ebi);
    }
    mme_app_handle_create_dedicated_bearer_rej(ue_context_p, ebi);
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
 ** Name:    _dedicated_eps_bearer_activate_t3485_handler()                **
 **                                                                        **
 ** Description: T3485 timeout handler                                     **
 **                                                                        **
 **              3GPP TS 24.301, section 6.4.2.6, case a                   **
 **      On the first expiry of the timer T3485, the MME shall re-         **
 **      send the ACTIVATE DEDICATED EPS BEARER CONTEXT REQUEST            **
 **      and shall reset and restart timer T3485. This retransmis-         **
 **      sion is repeated four times, i.e. on the fifth expiry of          **
 **      timer T3485, the MME shall abort the procedure, release           **
 **      any resources allocated for this activation and enter the         **
 **      state BEARER CONTEXT INACTIVE.                                    **
 **                                                                        **
 ** Inputs:  args:      handler parameters                                 **
 **      Others:    None                                                   **
 **                                                                        **
 ** Outputs:     None                                                      **
 **      Return:    None                                                   **
 **      Others:    None                                                   **
 **                                                                        **
 ***************************************************************************/
static void _dedicated_eps_bearer_activate_t3485_handler(
    void* args, imsi64_t* imsi64) {
  OAILOG_FUNC_IN(LOG_NAS_ESM);
  int rc;

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
    if (esm_ebr_timer_data->count < DEDICATED_EPS_BEARER_ACTIVATE_COUNTER_MAX) {
      /*
       * Re-send activate dedicated EPS bearer context request message
       * * * * to the UE
       */
      bstring b = bstrcpy(esm_ebr_timer_data->msg);
      rc        = _dedicated_eps_bearer_activate(
          esm_ebr_timer_data->ctx, esm_ebr_timer_data->ebi, &b);
      bdestroy_wrapper(&b);
    } else {
      /*
       * The maximum number of activate dedicated EPS bearer context request
       * message retransmission has exceed
       */

      /* Store ebi and ue_id as esm_ebr_timer_data gets freed in
       * esm_proc_eps_bearer_context_deactivate().
       */
      pdn_cid_t pid                 = MAX_APN_PER_UE;
      int bid                       = BEARERS_PER_UE;
      ebi_t ebi                     = esm_ebr_timer_data->ebi;
      ue_mm_context_t* ue_context_p = PARENT_STRUCT(
          esm_ebr_timer_data->ctx, struct ue_mm_context_s, emm_context);

      /*
       * Release the dedicated EPS bearer context, enter state INACTIVE and
       * stop T3485 timer. Timer is stopped inside
       * esm_proc_eps_bearer_context_deactivate()
       */
      rc = esm_proc_eps_bearer_context_deactivate(
          esm_ebr_timer_data->ctx, true, esm_ebr_timer_data->ebi, &pid, &bid,
          NULL);

      // Send dedicated_eps_bearer_reject to MME APP
      if ((rc != RETURNerror) && (ue_context_p)) {
        mme_app_handle_create_dedicated_bearer_rej(ue_context_p, ebi);
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
 ** Name:    _dedicated_eps_bearer_activate()                              **
 **                                                                        **
 ** Description: Sends ACTIVATE DEDICATED EPS BEREAR CONTEXT REQUEST mes-  **
 **      sage and starts timer T3485                                       **
 **                                                                        **
 ** Inputs:  ue_id:      UE local identifier                               **
 **      ebi:       EPS bearer identity                                    **
 **      msg:       Encoded ESM message to be sent                         **
 **      Others:    None                                                   **
 **                                                                        **
 ** Outputs:     None                                                      **
 **      Return:    RETURNok, RETURNerror                                  **
 **      Others:    T3485                                                  **
 **                                                                        **
 ***************************************************************************/
static int _dedicated_eps_bearer_activate(
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
   * Notify EMM that an activate dedicated EPS bearer context request
   * message has to be sent to the UE
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
  rc              = emm_sap_send(&emm_sap);

  if (rc != RETURNerror) {
    /*
     * Start T3485 retransmission timer
     */
    rc = esm_ebr_start_timer(
        emm_context, ebi, msg_dup, mme_config.nas_config.t3485_sec,
        _dedicated_eps_bearer_activate_t3485_handler);
  }
  bdestroy_wrapper(&msg_dup);
  OAILOG_FUNC_RETURN(LOG_NAS_ESM, rc);
}

/****************************************************************************
 **                                                                        **
 ** Name:    _erab_setup_rsp_tmr_exp_ded_bearer_handler()                  **
 **                                                                        **
 ** Description: Handles Erab setup rsp timer expiry                       **
 **                                                                        **
 ** Inputs:                                                                **
 **      imsi64:     IMSI                                                  **
 **      args:       timer data                                            **
 **                                                                        **
 ** Outputs:     None                                                      **
 **      Return:  None                                                     **
 **      Others:  None                                                     **
 **                                                                        **
 ***************************************************************************/

static void _erab_setup_rsp_tmr_exp_ded_bearer_handler(
    void* args, imsi64_t* imsi64) {
  OAILOG_FUNC_IN(LOG_NAS_ESM);
  int rc;

  // Get retransmission timer parameters data
  esm_ebr_timer_data_t* esm_ebr_timer_data = (esm_ebr_timer_data_t*) (args);

  if (esm_ebr_timer_data) {
    // Increment the retransmission counter
    esm_ebr_timer_data->count += 1;
    OAILOG_WARNING(
        LOG_NAS_ESM,
        "ESM-PROC  - erab_setup_rsp timer expired (ue_id=" MME_UE_S1AP_ID_FMT
        ", ebi=%d), "
        "retransmission counter = %d\n",
        esm_ebr_timer_data->ue_id, esm_ebr_timer_data->ebi,
        esm_ebr_timer_data->count);

    *imsi64 = esm_ebr_timer_data->ctx->_imsi64;
    ue_mm_context_t* ue_mm_context =
        mme_ue_context_exists_mme_ue_s1ap_id(esm_ebr_timer_data->ue_id);

    bearer_context_t* bearer_ctx =
        mme_app_get_bearer_context(ue_mm_context, esm_ebr_timer_data->ebi);
    if (!bearer_ctx) {
      OAILOG_ERROR(
          LOG_NAS_ESM,
          "Bearer context is NULL for (ebi=%u)"
          "\n",
          esm_ebr_timer_data->ebi);
      OAILOG_FUNC_OUT(LOG_NAS_ESM);
    }
    if (!bearer_ctx->enb_fteid_s1u.teid) {
      if (esm_ebr_timer_data->count < ERAB_SETUP_RSP_COUNTER_MAX) {
        // Restart the timer
        rc = esm_ebr_start_timer(
            esm_ebr_timer_data->ctx, esm_ebr_timer_data->ebi, NULL,
            ERAB_SETUP_RSP_TMR, _erab_setup_rsp_tmr_exp_ded_bearer_handler);
        if (rc != RETURNerror) {
          OAILOG_INFO(
              LOG_NAS_ESM,
              "ESM-PROC  - Started ERAB_SETUP_RSP_TMR for "
              "ue_id=" MME_UE_S1AP_ID_FMT
              "ebi (%u)"
              "\n",
              esm_ebr_timer_data->ue_id, esm_ebr_timer_data->ebi);
        }
      } else {
        // Dedicated bearers on S1 will not be set up.
        // UE will need to release the session or perform attach/detach
        // to recover. Or network side should release the bearers by disabling
        // policy rules for the subscriber.
        OAILOG_WARNING(
            LOG_NAS_ESM,
            "ESM-PROC  - ERAB_SETUP_RSP_COUNTER_MAX reached for ERAB_SETUP_RSP "
            "ue_id= " MME_UE_S1AP_ID_FMT
            " ebi (%u)"
            "\n",
            esm_ebr_timer_data->ue_id, esm_ebr_timer_data->ebi);
        if (bearer_ctx->esm_ebr_context.timer.id != NAS_TIMER_INACTIVE_ID) {
          bearer_ctx->esm_ebr_context.timer.id = NAS_TIMER_INACTIVE_ID;
        }
        if (esm_ebr_timer_data) {
          free_wrapper((void**) &esm_ebr_timer_data);
        }
      }
    } else {
      mme_app_handle_create_dedicated_bearer_rsp(
          ue_mm_context, esm_ebr_timer_data->ebi);
      if (bearer_ctx->esm_ebr_context.timer.id != NAS_TIMER_INACTIVE_ID) {
        bearer_ctx->esm_ebr_context.timer.id = NAS_TIMER_INACTIVE_ID;
      }
      if (esm_ebr_timer_data) {
        free_wrapper((void**) &esm_ebr_timer_data);
      }
    }
  }
  OAILOG_FUNC_OUT(LOG_NAS_ESM);
}

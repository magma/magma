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
#include <string.h>
#include <stdlib.h>

#include "log.h"
#include "common_types.h"
#include "3gpp_24.007.h"
#include "3gpp_24.008.h"
#include "mme_app_ue_context.h"
#include "mme_app_bearer_context.h"
#include "common_defs.h"
#include "emm_data.h"
#include "esm_ebr.h"
#include "esm_ebr_context.h"
#include "nas_timer.h"
#include "esm_cause.h"
#include "esm_proc.h"
#include "dynamic_memory_check.h"
#include "mme_config.h"

/****************************************************************************/
/****************  E X T E R N A L    D E F I N I T I O N S  ****************/
/****************************************************************************/

/****************************************************************************/
/*******************  L O C A L    D E F I N I T I O N S  *******************/
/****************************************************************************/

/****************************************************************************/
/******************  E X P O R T E D    F U N C T I O N S  ******************/
/****************************************************************************/

/****************************************************************************
 **                                                                        **
 ** Name:    esm_ebr_context_create()                                  **
 **                                                                        **
 ** Description: Creates a new EPS bearer context to the PDN with the spe- **
 **      cified PDN connection identifier                          **
 **                                                                        **
 ** Inputs:  ue_id:      UE identifier                              **
 **      pid:       PDN connection identifier                  **
 **      ebi:       EPS bearer identity                        **
 **      is_default:    true if the new bearer is a default EPS    **
 **             bearer context                             **
 **      esm_qos:   EPS bearer level QoS parameters            **
 **      tft:       Traffic flow template parameters           **
 **      Others:    _esm_data                                  **
 **                                                                        **
 ** Outputs:     None                                                      **
 **      Return:    The EPS bearer identity of the default EPS **
 **             bearer associated to the new EPS bearer    **
 **             context if successfully created;           **
 **             UNASSIGN EPS bearer value otherwise.       **
 **      Others:    _esm_data                                  **
 **                                                                        **
 ***************************************************************************/
ebi_t esm_ebr_context_create(
    emm_context_t* emm_context, const proc_tid_t pti, pdn_cid_t pid, ebi_t ebi,
    bool is_default, const qci_t qci, const bitrate_t gbr_dl,
    const bitrate_t gbr_ul, const bitrate_t mbr_dl, const bitrate_t mbr_ul,
    traffic_flow_template_t* tft, protocol_configuration_options_t* pco,
    fteid_t* sgw_fteid) {
  int bidx               = 0;
  esm_context_t* esm_ctx = NULL;
  esm_pdn_t* pdn         = NULL;

  OAILOG_FUNC_IN(LOG_NAS_ESM);
  esm_ctx = &emm_context->esm_ctx;
  ue_mm_context_t* ue_mm_context =
      PARENT_STRUCT(emm_context, struct ue_mm_context_s, emm_context);
  if (ue_mm_context == NULL) {
    OAILOG_ERROR(LOG_NAS_ESM, "ESM-PROC  - ue_mm_context null\n ");
    OAILOG_FUNC_RETURN(LOG_NAS_ESM, RETURNerror);
  }
  OAILOG_INFO(
      LOG_NAS_ESM,
      "ESM-PROC  - Create new %s EPS bearer context (ebi=%d) "
      "for PDN connection (pid=%d) for ue_id " MME_UE_S1AP_ID_FMT "\n",
      (is_default) ? "default" : "dedicated", ebi, pid,
      ue_mm_context->mme_ue_s1ap_id);

  if (pid < MAX_APN_PER_UE) {
    if (ue_mm_context->pdn_contexts[pid] == NULL) {
      OAILOG_ERROR(
          LOG_NAS_ESM,
          "ESM-PROC  - PDN connection %d has not been "
          "allocated for (ue_id = " MME_UE_S1AP_ID_FMT ")\n",
          pid, ue_mm_context->mme_ue_s1ap_id);
    }
    /*
     * Check the total number of active EPS bearers
     */
    else if (esm_ctx->n_active_ebrs > BEARERS_PER_UE) {
      OAILOG_WARNING(
          LOG_NAS_ESM,
          "ESM-PROC  - The total number of active EPS"
          "bearers is exceeded for ue id " MME_UE_S1AP_ID_FMT "\n",
          ue_mm_context->mme_ue_s1ap_id);
      OAILOG_FUNC_RETURN(LOG_NAS_ESM, ESM_EBI_UNASSIGNED);
    } else {
      bidx = EBI_TO_INDEX(ebi);
      /*
       * Get the PDN connection entry
       */
      pdn = &ue_mm_context->pdn_contexts[pid]->esm_data;

      if ((ue_mm_context->bearer_contexts[bidx]) &&
          (ESM_EBR_INACTIVE !=
           ue_mm_context->bearer_contexts[bidx]->esm_ebr_context.status)) {
        OAILOG_WARNING(
            LOG_NAS_ESM,
            "ESM-PROC  - A EPS bearer context "
            "is already allocated\n");
        OAILOG_FUNC_RETURN(LOG_NAS_ESM, ESM_EBI_UNASSIGNED);
      }
    }
  }

  if (bidx < BEARERS_PER_UE) {
    /*
     * Create new EPS bearer context
     */
    bearer_context_t* bearer_context = NULL;
    if (ue_mm_context->bearer_contexts[bidx]) {
      bearer_context = ue_mm_context->bearer_contexts[bidx];
    } else {
      bearer_context =
          mme_app_create_bearer_context(ue_mm_context, pid, ebi, is_default);
    }

    if (bearer_context) {
      bearer_context->transaction_identifier = pti;
      /*
       * Increment the total number of active EPS bearers
       */
      esm_ctx->n_active_ebrs += 1;
      /*
       * Increment the number of EPS bearer for this PDN connection
       */
      pdn->n_bearers += 1;
      /*
       * Setup the EPS bearer data
       */

      bearer_context->qci                    = qci;
      bearer_context->esm_ebr_context.gbr_dl = gbr_dl;
      bearer_context->esm_ebr_context.gbr_ul = gbr_ul;
      bearer_context->esm_ebr_context.mbr_dl = mbr_dl;
      bearer_context->esm_ebr_context.mbr_ul = mbr_ul;

      if (bearer_context->esm_ebr_context.tft) {
        free_traffic_flow_template(&bearer_context->esm_ebr_context.tft);
      }
      bearer_context->esm_ebr_context.tft = tft;

      if (bearer_context->esm_ebr_context.pco) {
        free_protocol_configuration_options(
            &bearer_context->esm_ebr_context.pco);
      }
      bearer_context->esm_ebr_context.pco = pco;
      if (sgw_fteid) {
        memcpy(&bearer_context->s_gw_fteid_s1u, sgw_fteid, sizeof(fteid_t));
      }

      if (is_default) {
        /*
         * Set the PDN connection activation indicator
         */
        ue_mm_context->pdn_contexts[pid]->is_active = true;

        ue_mm_context->pdn_contexts[pid]->default_ebi = ebi;
        /*
         * Update the emergency bearer services indicator
         */
        if (pdn->is_emergency) {
          esm_ctx->is_emergency = true;
        }
      }

      /*
       * Return the EPS bearer identity of the default EPS bearer
       * * * * associated to the new EPS bearer context
       */
      OAILOG_FUNC_RETURN(
          LOG_NAS_ESM, ue_mm_context->pdn_contexts[pid]->default_ebi);
    }

    OAILOG_WARNING(
        LOG_NAS_ESM,
        "ESM-PROC  - Failed to create new EPS bearer "
        "context (ebi=%d) for ue id " MME_UE_S1AP_ID_FMT "\n",
        ebi, ue_mm_context->mme_ue_s1ap_id);
  }

  OAILOG_FUNC_RETURN(LOG_NAS_ESM, ESM_EBI_UNASSIGNED);
}

//------------------------------------------------------------------------------
void esm_ebr_context_init(esm_ebr_context_t* esm_ebr_context) {
  OAILOG_FUNC_IN(LOG_NAS_ESM);
  if (esm_ebr_context) {
    memset(esm_ebr_context, 0, sizeof(*esm_ebr_context));
    /*
     * Set the EPS bearer context status to INACTIVE
     */
    esm_ebr_context->status = ESM_EBR_INACTIVE;
    /*
     * Disable the retransmission timer
     */
    esm_ebr_context->timer.id = NAS_TIMER_INACTIVE_ID;
  }
  OAILOG_FUNC_OUT(LOG_NAS_ESM);
}

/****************************************************************************
 **                                                                        **
 ** Name:    esm_ebr_context_release()                                 **
 **                                                                        **
 ** Description: Releases EPS bearer context entry previously allocated    **
 **      to the EPS bearer with the specified EPS bearer identity  **
 **                                                                        **
 ** Inputs:  ue_id:      UE identifier                              **
 **      ebi:       EPS bearer identity                        **
 **      Others:    _esm_data                                  **
 **                                                                        **
 ** Outputs:     pid:       Identifier of the PDN connection entry the **
 **             EPS bearer context belongs to              **
 **      bid:       Identifier of the released EPS bearer con- **
 **             text entry                                 **
 **      Return:    The EPS bearer identity associated to the  **
 **             EPS bearer context if successfully relea-  **
 **             sed; UNASSIGN EPS bearer value otherwise.  **
 **      Others:    _esm_data                                  **
 **                                                                        **
 ***************************************************************************/
ebi_t esm_ebr_context_release(
    emm_context_t* emm_context, ebi_t ebi, pdn_cid_t* pid, int* bid) {
  OAILOG_FUNC_IN(LOG_NAS_ESM);
  int found      = false;
  esm_pdn_t* pdn = NULL;
  // esm_cause_t esm_cause = ESM_CAUSE_SUCCESS;

  ue_mm_context_t* ue_mm_context =
      PARENT_STRUCT(emm_context, struct ue_mm_context_s, emm_context);

  if (ebi != ESM_EBI_UNASSIGNED) {
    /*
     * The identity of the EPS bearer to released is given;
     * Release the EPS bearer context entry that match the specified EPS
     * bearer identity
     */

    for (*bid = 0; *bid < BEARERS_PER_UE; (*bid)++) {
      if (ue_mm_context->bearer_contexts[*bid]) {
        if (ue_mm_context->bearer_contexts[*bid]->ebi != ebi) {
          continue;
        }

        /*
         * The EPS bearer context entry is found
         */
        *pid = ue_mm_context->bearer_contexts[*bid]->pdn_cx_id;
        if (ue_mm_context->pdn_contexts[*pid]) {
          found = true;
          pdn   = &ue_mm_context->pdn_contexts[*pid]->esm_data;
          break;
        }
      }
    }
  } else {
    /*
     * The identity of the EPS bearer to released is not given;
     * Release the EPS bearer context entry allocated with the EPS
     * bearer context identifier (bid) to establish connectivity to
     * the PDN identified by the PDN connection identifier (pid).
     * Default EPS bearer to a given PDN is always identified by the
     * first EPS bearer context entry at index bid = 0
     */
    if (*pid < MAX_APN_PER_UE) {
      if (!ue_mm_context->pdn_contexts[*pid]) {
        OAILOG_ERROR(
            LOG_NAS_ESM,
            "ESM-PROC  - PDN connection identifier %d "
            "is not valid for ue id " MME_UE_S1AP_ID_FMT "\n",
            *pid, ue_mm_context->mme_ue_s1ap_id);
      } else {
        if (ue_mm_context->pdn_contexts[*pid]->bearer_contexts[*bid]) {
          pdn   = &ue_mm_context->pdn_contexts[*pid]->esm_data;
          ebi   = ue_mm_context->bearer_contexts[*bid]->ebi;
          found = true;
        }
      }
    }
  }

  if (found) {
    int i;

    /*
     * Delete the specified EPS bearer context entry
     */

    if (!ue_mm_context->pdn_contexts[*pid]) {
      OAILOG_ERROR_UE(
          LOG_NAS_ESM, ue_mm_context->emm_context._imsi64,
          "ESM-PROC  - PDN context does not exist for bearer id %u,"
          "for ue id " MME_UE_S1AP_ID_FMT "\n",
          *bid, ue_mm_context->mme_ue_s1ap_id);
      OAILOG_FUNC_RETURN(LOG_NAS_ESM, ESM_EBI_UNASSIGNED);
    }
    if (ue_mm_context->pdn_contexts[*pid]->bearer_contexts[*bid] != *bid) {
      OAILOG_ERROR_UE(
          LOG_NAS_ESM, ue_mm_context->emm_context._imsi64,
          "ESM-PROC  - EPS bearer identifier %d is "
          "not valid for ue id " MME_UE_S1AP_ID_FMT "\n",
          *bid, ue_mm_context->mme_ue_s1ap_id);
      OAILOG_FUNC_RETURN(LOG_NAS_ESM, ESM_EBI_UNASSIGNED);
    }

    OAILOG_INFO(
        LOG_NAS_ESM,
        "ESM-PROC  - Release EPS bearer context "
        "(ebi=%d) for ue id " MME_UE_S1AP_ID_FMT "\n",
        ebi, ue_mm_context->mme_ue_s1ap_id);

    /*
     * Delete the TFT
     */
    if (ue_mm_context->bearer_contexts[*bid]->esm_ebr_context.tft) {
      free_traffic_flow_template(
          &ue_mm_context->bearer_contexts[*bid]->esm_ebr_context.tft);
    }

    /*
     * Delete PCO
     */
    if (ue_mm_context->bearer_contexts[*bid]->esm_ebr_context.pco) {
      free_protocol_configuration_options(
          &ue_mm_context->bearer_contexts[*bid]->esm_ebr_context.pco);
    }
    /*
     * Release the specified EPS bearer data
     */
    //  "free pdn->bearer"
    /*
     * Decrement the number of EPS bearer context allocated
     * * * * to the PDN connection
     */
    pdn->n_bearers -= 1;
    // Decrement the number of active bearers in ESM context
    emm_context->esm_ctx.n_active_ebrs -= 1;

    if (ue_mm_context->pdn_contexts[*pid]->default_ebi == ebi) {
      /*
       * 3GPP TS 24.301, section 6.4.4.3, 6.4.4.6
       * * * * If the EPS bearer identity is that of the default bearer to a
       * * * * PDN, the UE shall delete all EPS bearer contexts associated to
       * * * * that PDN connection.
       */
      OAILOG_INFO(
          LOG_NAS_ESM,
          "ESM-PROC  - Release default EPS bearer context "
          "(ebi=%d) for ue id " MME_UE_S1AP_ID_FMT "\n",
          ebi, ue_mm_context->mme_ue_s1ap_id);

      for (i = 1; pdn->n_bearers > 0; i++) {
        int idx = ue_mm_context->pdn_contexts[*pid]->bearer_contexts[i];
        if ((idx >= 0) && (idx < BEARERS_PER_UE)) {
          /* Delete only dedicated bearer. Default bearer will be deleted
           * outside this function
           */
          /* If ue_mm_context->bearer_contexts[idx] is NULL, move to
           * the next index
           */
          if (!ue_mm_context->bearer_contexts[idx]) {
            continue;
          }

          if (ue_mm_context->bearer_contexts[idx]->ebi ==
              ue_mm_context->pdn_contexts[*pid]->default_ebi) {
            continue;
          }
          OAILOG_WARNING(
              LOG_NAS_ESM,
              "ESM-PROC  - Release EPS bearer context "
              "(ebi=%d) " MME_UE_S1AP_ID_FMT "\n",
              ue_mm_context->bearer_contexts[idx]->ebi,
              ue_mm_context->mme_ue_s1ap_id);

          /*
           * Delete the TFT
           */
          if (ue_mm_context->bearer_contexts[idx]->esm_ebr_context.tft) {
            free_traffic_flow_template(
                &ue_mm_context->bearer_contexts[idx]->esm_ebr_context.tft);
          }

          /*
           * Delete PCO
           */
          if (ue_mm_context->bearer_contexts[idx]->esm_ebr_context.pco) {
            free_protocol_configuration_options(
                &ue_mm_context->bearer_contexts[idx]->esm_ebr_context.pco);
          }
          /*
           * Set the EPS bearer context state to INACTIVE
           */
          (void) esm_ebr_set_status(
              emm_context, ue_mm_context->bearer_contexts[idx]->ebi,
              ESM_EBR_INACTIVE, true);
          /*
           * Release EPS bearer data
           */
          (void) esm_ebr_release(
              emm_context, ue_mm_context->bearer_contexts[idx]->ebi);
          // esm_ebr_release()
          /*
           * Release dedicated EPS bearer data
           */
          free_wrapper((void**) &ue_mm_context->bearer_contexts[idx]);
          /*
           * Decrement the number of EPS bearer context allocated
           * * * * to the PDN connection
           */
          pdn->n_bearers -= 1;
        }
      }

      /*
       * Reset the PDN connection activation indicator
       */
      // TODO Look at "Reset the PDN connection activation indicator"
      // .is_active = false;

      /*
       * Update the emergency bearer services indicator
       */
      if (pdn->is_emergency) {
        pdn->is_emergency = false;
      }

      ue_mm_context->pdn_contexts[*pid]->is_active = false;
    } else {
      OAILOG_INFO(
          LOG_NAS_ESM,
          "ESM-PROC  - Release dedicated EPS bearer context "
          "(ebi=%d) for ue id " MME_UE_S1AP_ID_FMT "\n",
          ebi, ue_mm_context->mme_ue_s1ap_id);
    }

    // if (pdn->n_bearers == 0) {
    /*
     * : Release the PDN connection and marked the UE as inactive
     * * * * in the network for EPS services (is_attached = false)
     */
    //}

    OAILOG_FUNC_RETURN(LOG_NAS_ESM, ebi);
  }

  OAILOG_FUNC_RETURN(LOG_NAS_ESM, ESM_EBI_UNASSIGNED);
}

/****************************************************************************/
/*********************  L O C A L    F U N C T I O N S  *********************/
/****************************************************************************/

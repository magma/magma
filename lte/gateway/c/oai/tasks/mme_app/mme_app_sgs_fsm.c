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

#include <stddef.h>

/*****************************************************************************

  Source      mme_app_sgs_fsm.c

  Version

  Date

  Product

  Subsystem

  Author

  Description Defines the SGS State Machine handling

*****************************************************************************/
#include "log.h"
#include "mme_app_sgs_fsm.h"
#include "common_defs.h"
#include "common_types.h"
#include "mme_app_ue_context.h"

/****************************************************************************/
/****************  E X T E R N A L    D E F I N I T I O N S  ****************/
/****************************************************************************/

/*
   -----------------------------------------------------------------------------
            Data used for trace logging
   -----------------------------------------------------------------------------
*/

/* String representation of SGS events */
static const char* sgs_fsm_event_str[] = {
    "_SGS_LOCATION_UPDATE_ACCEPT", "_SGS_LOCATION_UPDATE_REJECT",
    "_SGS_PAGING_REQUEST",         "_SGS_SERVICE_ABORT_REQUEST",
    "_SGS_EPS_DETACH_IND",         "_SGS_IMSI_DETACH_IND",
    "_SGS_RESET_INDICATION",
};

/* String representation of SGS state */
static const char* sgs_fsm_state_str[SGS_STATE_MAX] = {
    "_SGS_INVALID",
    "_SGS_NULL",
    "_SGS_LA-UPDATE-REQUESTED",
    "_SGS_ASSOCIATED",
};

/*
   -----------------------------------------------------------------------------
        EPS Mobility Management state machine handlers
   -----------------------------------------------------------------------------
*/

/* Type of the SGS state machine handler */
typedef int (*sgs_fsm_handler_t)(const sgs_fsm_t*);

int sgs_null_handler(const sgs_fsm_t*);
int sgs_la_update_requested_handler(const sgs_fsm_t*);
int sgs_associated_handler(const sgs_fsm_t*);

/* SGS state machine handlers */
static const sgs_fsm_handler_t sgs_fsm_handlers[SGS_STATE_MAX] = {
    NULL,
    sgs_null_handler,
    sgs_la_update_requested_handler,
    sgs_associated_handler,
};

/****************************************************************************/
/******************  E X P O R T E D    F U N C T I O N S  ******************/
/****************************************************************************/

/****************************************************************************
 **                                                                        **
 ** Name:    sgs_fsm_initialize()                                          **
 **                                                                        **
 ** Description: Initializes the SGS state machine                         **
 **                                                                        **
 ** Inputs:  None                                                          **
 **      Others:    None                                                   **
 **                                                                        **
 ** Outputs:     None                                                      **
 **      Return:    None                                                   **
 **                                                                        **
 ***************************************************************************/
void sgs_fsm_initialize(void) {
  OAILOG_FUNC_IN(LOG_MME_APP);
  OAILOG_FUNC_OUT(LOG_MME_APP);
}

/****************************************************************************
 **                                                                        **
 ** Name:    sgs_fsm_process()                                             **
 **                                                                        **
 ** Description: Executes the SGS state machine                            **
 **                                                                        **
 ** Inputs:  evt:   Received SGSAP messsage to process                     **
 **                                                                        **
 ** Outputs:                                                               **
 **      Return:    RETURNok, RETURNerror                                  **
 **                                                                        **
 ***************************************************************************/
int sgs_fsm_process(const sgs_fsm_t* sgs_evt) {
  int rc = RETURNerror;
  sgs_fsm_state_t state;
  sgs_primitive_t primitive;

  OAILOG_FUNC_IN(LOG_MME_APP);

  primitive              = sgs_evt->primitive;
  sgs_context_t* sgs_ctx = (sgs_context_t*) (sgs_evt->ctx);

  if (sgs_ctx) {
    state = sgs_fsm_get_status(sgs_evt->ue_id, sgs_ctx);
    OAILOG_INFO(
        LOG_MME_APP, "SGS-FSM   - Received sgs-event %s (%d) in state %s\n",
        sgs_fsm_event_str[primitive], primitive, sgs_fsm_state_str[state]);
    /*
     * Execute the SGS state machine
     */
    rc = (sgs_fsm_handlers[state])(sgs_evt);
  } else {
    OAILOG_WARNING(
        LOG_MME_APP,
        "SGS-FSM   - Received event %s (%d) but no SGS context context found "
        "\n",
        sgs_fsm_event_str[primitive], primitive);
  }
  OAILOG_FUNC_RETURN(LOG_MME_APP, rc);
}

/****************************************************************************
 **                                                                        **
 ** Name:    sgs_fsm_set_status()                                      **
 **                                                                        **
 ** Description: Set the SGS fsm status to the given state **
 **                                                                        **
 ** Inputs:  ue_id:      Lower layers UE identifier                 **
 **      status:    The new SGS status                         **
 **      Others:    None                                       **
 **                                                                        **
 ** Outputs:     None                                                      **
 **      Return:    RETURNok, RETURNerror                      **
 **      Others:    _sgs_fsm_state                            **
 **                                                                        **
 ***************************************************************************/
int sgs_fsm_set_status(
    mme_ue_s1ap_id_t ue_id, void* ctx, sgs_fsm_state_t state) {
  OAILOG_FUNC_IN(LOG_MME_APP);
  sgs_context_t* sgs_ctx = (sgs_context_t*) ctx;

  if (sgs_ctx == NULL) {
    OAILOG_ERROR(
        LOG_MME_APP,
        "Invalid SGS context received for UE Id: " MME_UE_S1AP_ID_FMT "\n",
        ue_id);
    OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNerror);
  }

  if (state < SGS_STATE_MAX) {
    if (state != sgs_ctx->sgs_state) {
      OAILOG_INFO(
          LOG_MME_APP,
          "UE " MME_UE_S1AP_ID_FMT " SGS-FSM   - State changed: %s ===> %s\n",
          ue_id, sgs_fsm_state_str[sgs_ctx->sgs_state],
          sgs_fsm_state_str[state]);
      sgs_ctx->sgs_state = state;
    }

    OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNok);
  }
  OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNerror);
}

/****************************************************************************
 **                                                                        **
 ** Name:    sgs_fsm_get_status()                                      **
 **                                                                        **
 ** Description: Get the current value of the SGS fsm status
 **                                                                        **
 ** Inputs:  ue_id:      Lower layers UE identifier                 **
 **      Others:    _sgs_fsm_state                            **
 **                                                                        **
 ** Outputs:     None                                                      **
 **      Return:    The current value of the SGS status        **
 **      Others:    None                                       **
 **                                                                        **
 ***************************************************************************/
sgs_fsm_state_t sgs_fsm_get_status(mme_ue_s1ap_id_t ue_id, void* ctx) {
  sgs_context_t* sgs_ctx = (sgs_context_t*) ctx;

  if (sgs_ctx) {
    if ((sgs_ctx->sgs_state >= SGS_STATE_MAX) ||
        (sgs_ctx->sgs_state <= SGS_STATE_MIN)) {
      OAILOG_ERROR(
          LOG_MME_APP, "BAD SGS state (%d) for UE Id: " MME_UE_S1AP_ID_FMT "\n",
          sgs_ctx->sgs_state, ue_id);
      return SGS_INVALID;
    }
    return sgs_ctx->sgs_state;
  }
  return SGS_INVALID;
}

/****************************************************************************/
/*********************  L O C A L    F U N C T I O N S  *********************/
/****************************************************************************/

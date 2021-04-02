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

#include "emm_fsm.h"
#include "log.h"
#include "common_defs.h"
#include "mme_app_ue_context.h"
#include "mme_api.h"
#include "emm_data.h"
#include "assertions.h"
#include "common_types.h"
#include "emm_regDef.h"

/****************************************************************************/
/****************  E X T E R N A L    D E F I N I T I O N S  ****************/
/****************************************************************************/

#define EMM_FSM_NB_UE_MAX (MME_API_NB_UE_MAX + 1)

/*
   -----------------------------------------------------------------------------
            Data used for trace logging
   -----------------------------------------------------------------------------
*/

/* String representation of EMM events */
static const char* emm_fsm_event_str[] = {
    "COMMON_PROC_REQ",
    "COMMON_PROC_CNF",
    "COMMON_PROC_REJ",
    "COMMON_PROC_ABORT",
    "ATTACH_CNF",
    "ATTACH_REJ",
    "ATTACH_ABORT",
    "DETACH_INIT",
    "DETACH_REQ",
    "DETACH_FAILED",
    "DETACH_CNF",
    "TAU_REQ",
    "TAU_CNF",
    "TAU_REJ",
    "SERVICE_REQ",
    "SERVICE_CNF",
    "SERVICE_REJ",
    "LOWERLAYER_SUCCESS",
    "LOWERLAYER_FAILURE",
    "LOWERLAYER_RELEASE",
    "LOWERLAYER_NON_DELIVERY",
};

/* String representation of EMM state */
static const char* const emm_fsm_status_str[EMM_STATE_MAX] = {
    "INVALID",
    "EMM-DEREGISTERED",
    "EMM-REGISTERED",
    "EMM-DEREGISTERED-INITIATED",
    "EMM-COMMON-PROCEDURE-INITIATED",
};

/*
   -----------------------------------------------------------------------------
        EPS Mobility Management state machine handlers
   -----------------------------------------------------------------------------
*/

/* Type of the EPS Mobility Management state machine handler */
typedef int (*emm_fsm_handler_t)(emm_reg_t* const);

int EmmDeregistered(emm_reg_t* const);
int EmmRegistered(emm_reg_t* const);
int EmmDeregisteredInitiated(emm_reg_t* const);
int EmmCommonProcedureInitiated(emm_reg_t* const);

/* EMM state machine handlers */
static const emm_fsm_handler_t emm_fsm_handlers[EMM_STATE_MAX] = {
    NULL,
    EmmDeregistered,
    EmmRegistered,
    EmmDeregisteredInitiated,
    EmmCommonProcedureInitiated,
};

/*
   -----------------------------------------------------------------------------
            Current EPS Mobility Management state
   -----------------------------------------------------------------------------
*/

/****************************************************************************/
/******************  E X P O R T E D    F U N C T I O N S  ******************/
/****************************************************************************/

/****************************************************************************
 **                                                                        **
 ** Name:    emm_fsm_initialize()                                      **
 **                                                                        **
 ** Description: Initializes the EMM state machine                         **
 **                                                                        **
 ** Inputs:  None                                                      **
 **      Others:    None                                       **
 **                                                                        **
 ** Outputs:     None                                                      **
 **      Return:    None                                       **
 **      Others:    _emm_fsm_status                            **
 **                                                                        **
 ***************************************************************************/
void emm_fsm_initialize(void) {
  OAILOG_FUNC_IN(LOG_NAS_EMM);
  OAILOG_FUNC_OUT(LOG_NAS_EMM);
}

/****************************************************************************
 **                                                                        **
 ** Name:    emm_fsm_set_status()                                      **
 **                                                                        **
 ** Description: Set the EPS Mobility Management state to the given state **
 **                                                                        **
 ** Inputs:  ue_id:      Lower layers UE identifier                 **
 **      state:    The new EMM state                         **
 **      Others:    None                                       **
 **                                                                        **
 ** Outputs:     None                                                      **
 **      Return:    RETURNok, RETURNerror                      **
 **      Others:    _emm_fsm_status                            **
 **                                                                        **
 ***************************************************************************/
int emm_fsm_set_state(
    const mme_ue_s1ap_id_t ue_id, struct emm_context_s* const emm_context,
    const emm_fsm_state_t state) {
  OAILOG_FUNC_IN(LOG_NAS_EMM);

  DevAssert(emm_context);
  if (state < EMM_STATE_MAX) {
    if (state != emm_context->_emm_fsm_state) {
      OAILOG_INFO(
          LOG_NAS_EMM,
          "UE " MME_UE_S1AP_ID_FMT " EMM-FSM   - Status changed: %s ===> %s\n",
          ue_id, emm_fsm_status_str[emm_context->_emm_fsm_state],
          emm_fsm_status_str[state]);
      emm_context->_emm_fsm_state   = state;
      emm_fsm_state_t new_emm_state = UE_UNREGISTERED;
      if (state == EMM_REGISTERED) {
        new_emm_state = UE_REGISTERED;
      } else if (state == EMM_DEREGISTERED) {
        new_emm_state = UE_UNREGISTERED;
      }
      // Update mme_ue_context's emm_state and overall stats
      mme_ue_context_update_ue_emm_state(ue_id, new_emm_state);
    }

    OAILOG_FUNC_RETURN(LOG_NAS_EMM, RETURNok);
  }
  OAILOG_FUNC_RETURN(LOG_NAS_EMM, RETURNerror);
}

/****************************************************************************
 **                                                                        **
 ** Name:    emm_fsm_get_state()                                      **
 **                                                                        **
 ** Description: Get the current value of the EPS Mobility Management      **
 **      state                                                    **
 **                                                                        **
 ** Inputs:  ue_id:      Lower layers UE identifier                 **
 **      Others:    _emm_fsm_status                            **
 **                                                                        **
 ** Outputs:     None                                                      **
 **      Return:    The current value of the EMM state        **
 **      Others:    None                                       **
 **                                                                        **
 ***************************************************************************/
emm_fsm_state_t emm_fsm_get_state(
    const struct emm_context_s* const emm_context) {
  if (emm_context) {
    AssertFatal(
        (emm_context->_emm_fsm_state < EMM_STATE_MAX) &&
            (emm_context->_emm_fsm_state >= EMM_STATE_MIN),
        "ue_id " MME_UE_S1AP_ID_FMT " BAD EMM state %d",
        PARENT_STRUCT(emm_context, struct ue_mm_context_s, emm_context)
            ->mme_ue_s1ap_id,
        emm_context->_emm_fsm_state);
    return emm_context->_emm_fsm_state;
  }
  return EMM_INVALID;
}

//------------------------------------------------------------------------------
const char* emm_fsm_get_state_str(
    const struct emm_context_s* const emm_context) {
  if (emm_context) {
    emm_fsm_state_t state = emm_fsm_get_state(emm_context);
    return emm_fsm_status_str[state];
  }
  return emm_fsm_status_str[EMM_INVALID];
}

/****************************************************************************
 **                                                                        **
 ** Name:    emm_fsm_process()                                         **
 **                                                                        **
 ** Description: Executes the EMM state machine                            **
 **                                                                        **
 ** Inputs:  evt:       The EMMREG-SAP event to process            **
 **      Others:    _emm_fsm_status                            **
 **                                                                        **
 ** Outputs:     None                                                      **
 **      Return:    RETURNok, RETURNerror                      **
 **      Others:    None                                       **
 **                                                                        **
 ***************************************************************************/
int emm_fsm_process(struct emm_reg_s* const evt) {
  int rc = RETURNerror;
  emm_fsm_state_t state;
  emm_reg_primitive_t primitive;

  OAILOG_FUNC_IN(LOG_NAS_EMM);
  primitive              = evt->primitive;
  emm_context_t* emm_ctx = (emm_context_t*) evt->ctx;

  if (emm_ctx) {
    state = emm_fsm_get_state(emm_ctx);
    OAILOG_INFO(
        LOG_NAS_EMM, "EMM-FSM   - Received event %s (%d) in state %s\n",
        emm_fsm_event_str[primitive - _EMMREG_START - 1], primitive,
        emm_fsm_status_str[state]);
    /*
     * Execute the EMM state machine
     */
    if (emm_ctx->is_imsi_only_detach == false) {
      rc = (emm_fsm_handlers[state])(evt);
    }
  } else {
    OAILOG_WARNING(
        LOG_NAS_EMM,
        "EMM-FSM   - Received event %s (%d) but no EMM data context provided\n",
        emm_fsm_event_str[primitive - _EMMREG_START - 1], primitive);
  }
  OAILOG_FUNC_RETURN(LOG_NAS_EMM, rc);
}

/****************************************************************************/
/*********************  L O C A L    F U N C T I O N S  *********************/
/****************************************************************************/

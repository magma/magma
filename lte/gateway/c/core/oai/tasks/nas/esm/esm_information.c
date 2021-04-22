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

#include <stdint.h>
#include <stdlib.h>

#include "bstrlib.h"
#include "log.h"
#include "dynamic_memory_check.h"
#include "common_types.h"
#include "3gpp_24.007.h"
#include "3gpp_24.008.h"
#include "mme_app_ue_context.h"
#include "esm_proc.h"
#include "emm_data.h"
#include "esm_data.h"
#include "esm_cause.h"
#include "emm_sap.h"
#include "esm_send.h"
#include "3gpp_36.401.h"
#include "EsmCause.h"
#include "common_defs.h"
#include "emm_esmDef.h"
#include "esm_msg.h"
#include "nas_timer.h"
/****************************************************************************/
/****************  E X T E R N A L    D E F I N I T I O N S  ****************/
/****************************************************************************/

/****************************************************************************/
/*******************  L O C A L    D E F I N I T I O N S  *******************/
/****************************************************************************/

/*
   Timer handlers
*/
static void esm_information_t3489_handler(void*, imsi64_t* imsi64);

/* Maximum value of the deactivate EPS bearer context request
   retransmission counter */
#define ESM_INFORMATION_COUNTER_MAX 3

static int esm_information(
    emm_context_t* emm_context_p, ebi_t ebi, esm_ebr_timer_data_t* const data);

/****************************************************************************/
/******************  E X P O R T E D    F U N C T I O N S  ******************/
/****************************************************************************/

//------------------------------------------------------------------------------
int esm_proc_esm_information_request(
    emm_context_t* const emm_context_p, const pti_t pti) {
  OAILOG_FUNC_IN(LOG_NAS_ESM);
  int rc;
  mme_ue_s1ap_id_t ue_id =
      PARENT_STRUCT(emm_context_p, struct ue_mm_context_s, emm_context)
          ->mme_ue_s1ap_id;

  OAILOG_INFO(
      LOG_NAS_ESM,
      "ESM-PROC  - Initiate ESM information ue_id=(" MME_UE_S1AP_ID_FMT ")\n",
      ue_id);

  ESM_msg esm_msg = {.header = {0}};
  rc              = esm_send_esm_information_request(
      pti, EPS_BEARER_IDENTITY_UNASSIGNED, &esm_msg.esm_information_request);

  if (rc != RETURNerror) {
    /*
     * Encode the returned ESM response message
     */
    char emm_sap_buffer[16];  // very short msg
    int size        = esm_msg_encode(&esm_msg, (uint8_t*) emm_sap_buffer, 16);
    bstring msg_req = NULL;
    OAILOG_INFO(LOG_NAS_EMM, "ESM encoded MSG size %d\n", size);
    if (size > 0) {
      msg_req = blk2bstr(emm_sap_buffer, size);
      /*
       * Send esm information request message and
       * start timer T3489
       */
      esm_ebr_timer_data_t* data =
          (esm_ebr_timer_data_t*) calloc(1, sizeof(*data));
      data->ctx   = emm_context_p;
      data->ebi   = EPS_BEARER_IDENTITY_UNASSIGNED;
      data->msg   = msg_req;
      data->ue_id = ue_id;
      rc          = esm_information(emm_context_p, pti, data);
    }
  }
  OAILOG_FUNC_RETURN(LOG_NAS_ESM, rc);
}

//------------------------------------------------------------------------------
int esm_proc_esm_information_response(
    emm_context_t* emm_context_p, pti_t pti, const_bstring const apn,
    const protocol_configuration_options_t* const pco,
    esm_cause_t* const esm_cause) {
  OAILOG_FUNC_IN(LOG_NAS_ESM);
  int rc = RETURNok;

  /*
   * Stop T3489 timer if running
   */
  nas_stop_T3489(&emm_context_p->esm_ctx);

  if (apn && (apn->slen > 0)) {
    if (emm_context_p->esm_ctx.esm_proc_data->apn) {
      bdestroy_wrapper(&emm_context_p->esm_ctx.esm_proc_data->apn);
    }
    emm_context_p->esm_ctx.esm_proc_data->apn = bstrcpy(apn);
  }

  if ((pco) && (pco->num_protocol_or_container_id)) {
    if (emm_context_p->esm_ctx.esm_proc_data->pco
            .num_protocol_or_container_id) {
      clear_protocol_configuration_options(
          &emm_context_p->esm_ctx.esm_proc_data->pco);
    }
    copy_protocol_configuration_options(
        &emm_context_p->esm_ctx.esm_proc_data->pco, pco);
  }

  *esm_cause = ESM_CAUSE_SUCCESS;

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
 ** Name:    _esm_information_t3489_handler()                    **
 **                                                                        **
 ** Description: T3489 timeout handler                                     **
 **                                                                        **
 **              3GPP TS 24.301, section 6.4.4.5, case a                   **
 **      On the first expiry of the timer T3489, the MME shall re- **
 **      send the DEACTIVATE EPS BEARER CONTEXT REQUEST and shall  **
 **      reset and restart timer T3489. This retransmission is     **
 **      repeated four times, i.e. on the fifth expiry of timer    **
 **      T3489, the MME shall abort the procedure and deactivate   **
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
static void esm_information_t3489_handler(void* args, imsi64_t* imsi64) {
  OAILOG_FUNC_IN(LOG_NAS_ESM);

  /*
   * Get retransmission timer parameters data
   */
  esm_ebr_timer_data_t* esm_ebr_timer_data = (esm_ebr_timer_data_t*) (args);

  if (esm_ebr_timer_data && esm_ebr_timer_data->ctx) {
    /*
     * Increment the retransmission counter
     */
    esm_ebr_timer_data->count += 1;
    OAILOG_WARNING(
        LOG_NAS_ESM,
        "ESM-PROC  - T3489 timer expired (ue_id=" MME_UE_S1AP_ID_FMT
        "), "
        "retransmission counter = %d\n",
        esm_ebr_timer_data->ue_id, esm_ebr_timer_data->count);

    *imsi64 = esm_ebr_timer_data->ctx->_imsi64;
    if (esm_ebr_timer_data->count < ESM_INFORMATION_COUNTER_MAX) {
      // Unset the timer id maintained in the esm_ctx, as the timer is no
      // longer valid.
      esm_ebr_timer_data->ctx->esm_ctx.T3489.id = NAS_TIMER_INACTIVE_ID;
      /*
       * Re-send deactivate EPS bearer context request message to the UE
       */
      esm_information(
          esm_ebr_timer_data->ctx, esm_ebr_timer_data->ebi, esm_ebr_timer_data);
    } else {
      /*
       * The maximum number of deactivate EPS bearer context request
       * message retransmission has exceed
       *
       * TODO call something like esm_send_pdn_connectivity_reject
       * #ESM information not received
       *
       * Stop timer T3489
       */
      esm_ebr_timer_data->ctx->esm_ctx.T3489.id = NAS_TIMER_INACTIVE_ID;
      /*
       * Re-start T3489 timer
       */
      bdestroy_wrapper(&esm_ebr_timer_data->msg);
      free_wrapper((void**) &esm_ebr_timer_data);
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
 ** Name:    _esm_information()                                            **
 **                                                                        **
 ** Description: Sends DEACTIVATE EPS BEREAR CONTEXT REQUEST message and   **
 **      starts timer T3489.                                               **
 **      Function also clearns out any existing T3489 timers referenced    **
 **      by the esm_ctx datastructure.                                     **
 **                                                                        **
 ** Inputs:  ue_id:      UE local identifier                               **
 **      ebi:       EPS bearer identity                                    **
 **      msg:       Encoded ESM message to be sent                         **
 **      Others:    None                                                   **
 **                                                                        **
 ** Outputs:     None                                                      **
 **      Return:    RETURNok, RETURNerror                                  **
 **      Others:    T3489                                                  **
 **                                                                        **
 ***************************************************************************/
static int esm_information(
    emm_context_t* emm_context_p, ebi_t ebi, esm_ebr_timer_data_t* const data) {
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
  emm_esm_data_t* emm_esm = &emm_sap.u.emm_esm.u.data;

  emm_sap.primitive       = EMMESM_UNITDATA_REQ;
  emm_sap.u.emm_esm.ue_id = ue_id;
  emm_sap.u.emm_esm.ctx   = emm_context_p;
  emm_esm->msg            = bstrcpy(data->msg);

  rc = emm_sap_send(&emm_sap);

  if (rc != RETURNerror) {
    nas_stop_T3489(&emm_context_p->esm_ctx);
    /*
     * Start T3489 timer
     */
    emm_context_p->esm_ctx.T3489.id = nas_timer_start(
        emm_context_p->esm_ctx.T3489.sec, 0 /*usec*/,
        esm_information_t3489_handler, data);

    OAILOG_INFO(
        LOG_NAS_EMM,
        "UE " MME_UE_S1AP_ID_FMT "Timer T3489 (%lx) expires in %d seconds\n",
        ue_id, emm_context_p->esm_ctx.T3489.id,
        emm_context_p->esm_ctx.T3489.sec);
  } else {
    bdestroy_wrapper(&data->msg);
    free_wrapper((void**) &data);
  }
  bdestroy_wrapper(&emm_esm->msg);
  OAILOG_FUNC_RETURN(LOG_NAS_ESM, rc);
}

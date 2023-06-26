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

#ifdef __cplusplus
extern "C" {
#endif
#include "lte/gateway/c/core/common/common_defs.h"
#include "lte/gateway/c/core/oai/common/common_types.h"
#include "lte/gateway/c/core/oai/common/log.h"
#include "lte/gateway/c/core/oai/lib/bstr/bstrlib.h"
#ifdef __cplusplus
}
#endif

#include "lte/gateway/c/core/common/dynamic_memory_check.h"
#include "lte/gateway/c/core/oai/include/mme_app_ue_context.hpp"
#include "lte/gateway/c/core/oai/lib/3gpp/3gpp_24.007.h"
#include "lte/gateway/c/core/oai/lib/3gpp/3gpp_24.008.h"
#include "lte/gateway/c/core/oai/lib/3gpp/3gpp_36.401.h"
#include "lte/gateway/c/core/oai/tasks/mme_app/mme_app_defs.hpp"
#include "lte/gateway/c/core/oai/tasks/mme_app/mme_app_timer.hpp"
#include "lte/gateway/c/core/oai/tasks/nas/emm/emm_data.hpp"
#include "lte/gateway/c/core/oai/tasks/nas/emm/sap/emm_esmDef.hpp"
#include "lte/gateway/c/core/oai/tasks/nas/emm/sap/emm_sap.hpp"
#include "lte/gateway/c/core/oai/tasks/nas/esm/esm_data.hpp"
#include "lte/gateway/c/core/oai/tasks/nas/esm/esm_proc.hpp"
#include "lte/gateway/c/core/oai/tasks/nas/esm/msg/esm_cause.hpp"
#include "lte/gateway/c/core/oai/tasks/nas/esm/msg/esm_msg.hpp"
#include "lte/gateway/c/core/oai/tasks/nas/esm/sap/esm_send.hpp"
#include "lte/gateway/c/core/oai/tasks/nas/ies/EsmCause.hpp"
#include "lte/gateway/c/core/oai/tasks/nas/util/nas_timer.hpp"

/****************************************************************************/
/****************  E X T E R N A L    D E F I N I T I O N S  ****************/
/****************************************************************************/

/****************************************************************************/
/*******************  L O C A L    D E F I N I T I O N S  *******************/
/****************************************************************************/

/* Maximum value of the deactivate EPS bearer context request
   retransmission counter */
#define ESM_INFORMATION_COUNTER_MAX 3

static status_code_e esm_information(emm_context_t* emm_context_p, ebi_t ebi,
                                     esm_ebr_timer_data_t* const data);

/****************************************************************************/
/******************  E X P O R T E D    F U N C T I O N S  ******************/
/****************************************************************************/

//------------------------------------------------------------------------------
status_code_e esm_proc_esm_information_request(
    emm_context_t* const emm_context_p, const pti_t pti) {
  OAILOG_FUNC_IN(LOG_NAS_ESM);
  status_code_e rc = RETURNok;
  mme_ue_s1ap_id_t ue_id =
      PARENT_STRUCT(emm_context_p, struct ue_mm_context_s, emm_context)
          ->mme_ue_s1ap_id;

  OAILOG_INFO_UE(
      LOG_NAS_ESM, emm_context_p->_imsi64,
      "ESM-PROC  - Initiate ESM information ue_id=(" MME_UE_S1AP_ID_FMT ")\n",
      ue_id);

  ESM_msg esm_msg = {.header = {0}};
  rc = esm_send_esm_information_request(pti, EPS_BEARER_IDENTITY_UNASSIGNED,
                                        &esm_msg.esm_information_request);

  if (rc != RETURNerror) {
    /*
     * Encode the returned ESM response message
     */
    char emm_sap_buffer[16];  // very short msg
    int size = esm_msg_encode(&esm_msg, (uint8_t*)emm_sap_buffer, 16);
    bstring msg_req = NULL;
    OAILOG_INFO_UE(LOG_NAS_EMM, emm_context_p->_imsi64,
                   "ESM encoded MSG size %d\n", size);
    if (size > 0) {
      msg_req = blk2bstr(emm_sap_buffer, size);
      /*
       * Send esm information request message and
       * start timer T3489
       */
      esm_ebr_timer_data_t* data =
          (esm_ebr_timer_data_t*)calloc(1, sizeof(*data));
      data->ctx = emm_context_p;
      data->ebi = EPS_BEARER_IDENTITY_UNASSIGNED;
      data->msg = msg_req;
      data->ue_id = ue_id;
      rc = esm_information(emm_context_p, pti, data);
    }
  }
  OAILOG_FUNC_RETURN(LOG_NAS_ESM, rc);
}

//------------------------------------------------------------------------------
status_code_e esm_proc_esm_information_response(
    emm_context_t* emm_context_p, pti_t pti, const_bstring const apn,
    const protocol_configuration_options_t* const pco,
    esm_cause_t* const esm_cause) {
  OAILOG_FUNC_IN(LOG_NAS_ESM);
  status_code_e rc = RETURNok;

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
 ** Name:    mme_app_handle_esm_information_t3489_expiry                   **
 **                                                                        **
 ** Description: T3489 timeout handler                                     **
 **                                                                        **
 **              3GPP TS 24.301, section 6.4.4.5, case a                   **
 **      On the first expiry of the timer T3489, the MME shall re-         **
 **      send the DEACTIVATE EPS BEARER CONTEXT REQUEST and shall          **
 **      reset and restart timer T3489. This retransmission is             **
 **      repeated four times, i.e. on the fifth expiry of timer            **
 **      T3489, the MME shall abort the procedure and deactivate           **
 **      the EPS bearer context locally.                                   **
 **                                                                        **
 ** Inputs:  args:      handler parameters                                 **
 **      Others:    None                                                   **
 **                                                                        **
 ** Outputs:     None                                                      **
 **      Return:    None                                                   **
 **      Others:    None                                                   **
 **                                                                        **
 ***************************************************************************/
status_code_e mme_app_handle_esm_information_t3489_expiry(zloop_t* loop,
                                                          int timer_id,
                                                          void* args) {
  OAILOG_FUNC_IN(LOG_NAS_ESM);

  mme_ue_s1ap_id_t mme_ue_s1ap_id = 0;
  if (!mme_pop_timer_arg_ue_id(timer_id, &mme_ue_s1ap_id)) {
    OAILOG_WARNING(LOG_NAS_EMM, "Invalid Timer Id expiration, Timer Id: %u\n",
                   timer_id);
    OAILOG_FUNC_RETURN(LOG_NAS_ESM, RETURNok);
  }

  struct ue_mm_context_s* ue_context_p = mme_app_get_ue_context_for_timer(
      mme_ue_s1ap_id,
      const_cast<char*>("Deactivate EPS Bearer Ctx Request T3489 Timer"));
  if (ue_context_p == NULL) {
    OAILOG_ERROR(
        LOG_MME_APP,
        "Invalid UE context received, MME UE S1AP Id: " MME_UE_S1AP_ID_FMT "\n",
        mme_ue_s1ap_id);
    OAILOG_FUNC_RETURN(LOG_NAS_ESM, RETURNok);
  }

  emm_context_t* emm_context = &ue_context_p->emm_context;
  /*
   * Get retransmission timer parameters data
   */
  esm_ebr_timer_data_t* esm_ebr_timer_data =
      (esm_ebr_timer_data_t*)(emm_context->esm_ctx.t3489_arg);

  if (esm_ebr_timer_data) {
    /*
     * Increment the retransmission counter
     */
    esm_ebr_timer_data->count += 1;
    OAILOG_WARNING_UE(
        LOG_NAS_ESM, emm_context->_imsi64,
        "ESM-PROC  - T3489 timer expired (ue_id=" MME_UE_S1AP_ID_FMT
        "), "
        "retransmission counter = %d\n",
        esm_ebr_timer_data->ue_id, esm_ebr_timer_data->count);

    if (esm_ebr_timer_data->count < ESM_INFORMATION_COUNTER_MAX) {
      // Unset the timer id maintained in the esm_ctx, as the timer is no
      // longer valid.
      emm_context->esm_ctx.T3489.id = NAS_TIMER_INACTIVE_ID;
      /*
       * Re-send deactivate EPS bearer context request message to the UE
       */
      esm_information(emm_context, esm_ebr_timer_data->ebi, esm_ebr_timer_data);
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
      emm_context->esm_ctx.T3489.id = NAS_TIMER_INACTIVE_ID;
      bdestroy_wrapper(&esm_ebr_timer_data->msg);
      free_wrapper((void**)&esm_ebr_timer_data);
    }
  }

  OAILOG_FUNC_RETURN(LOG_NAS_ESM, RETURNok);
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
 ** Description: Sends DEACTIVATE EPS BEARER CONTEXT REQUEST message and   **
 **      starts timer T3489.                                               **
 **      Function also cleans out any existing T3489 timers referenced     **
 **      by the esm_ctx data structure.                                    **
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
static status_code_e esm_information(emm_context_t* emm_context_p, ebi_t ebi,
                                     esm_ebr_timer_data_t* const data) {
  OAILOG_FUNC_IN(LOG_NAS_ESM);
  emm_sap_t emm_sap = {};
  status_code_e rc = RETURNerror;
  mme_ue_s1ap_id_t ue_id =
      PARENT_STRUCT(emm_context_p, struct ue_mm_context_s, emm_context)
          ->mme_ue_s1ap_id;

  /*
   * Notify EMM that a deactivate EPS bearer context request message
   * has to be sent to the UE
   */
  emm_esm_data_t* emm_esm = &emm_sap.u.emm_esm.u.data;

  emm_sap.primitive = EMMESM_UNITDATA_REQ;
  emm_sap.u.emm_esm.ue_id = ue_id;
  emm_sap.u.emm_esm.ctx = emm_context_p;
  emm_esm->msg = bstrcpy(data->msg);

  rc = emm_sap_send(&emm_sap);

  if (rc != RETURNerror) {
    nas_stop_T3489(&emm_context_p->esm_ctx);
    /*
     * Start T3489 timer
     */
    nas_start_T3489(ue_id, &(emm_context_p->esm_ctx.T3489),
                    (time_out_t)mme_app_handle_esm_information_t3489_expiry);
  }
  if (NAS_TIMER_INACTIVE_ID != emm_context_p->esm_ctx.T3489.id) {
    emm_context_p->esm_ctx.t3489_arg = (void*)data;

    OAILOG_INFO_UE(LOG_NAS_EMM, emm_context_p->_imsi64,
                   "UE " MME_UE_S1AP_ID_FMT
                   "Timer T3489 (%lx) expires in %d seconds\n",
                   ue_id, emm_context_p->esm_ctx.T3489.id,
                   emm_context_p->esm_ctx.T3489.msec);
  } else {
    bdestroy_wrapper(&data->msg);
    free_wrapper((void**)&data);
    rc = RETURNerror;
  }
  bdestroy_wrapper(&emm_esm->msg);
  OAILOG_FUNC_RETURN(LOG_NAS_ESM, rc);
}

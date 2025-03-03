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
#include <stdbool.h>
#include <stdlib.h>

#ifdef __cplusplus
extern "C" {
#endif
#include "lte/gateway/c/core/common/assertions.h"
#include "lte/gateway/c/core/common/common_defs.h"
#include "lte/gateway/c/core/oai/common/common_types.h"
#include "lte/gateway/c/core/oai/common/conversions.h"
#include "lte/gateway/c/core/oai/common/log.h"
#ifdef __cplusplus
}
#endif

#include "lte/gateway/c/core/oai/include/3gpp_requirements_24.301.h"
#include "lte/gateway/c/core/oai/include/mme_app_state.hpp"
#include "lte/gateway/c/core/oai/include/mme_app_ue_context.hpp"
#include "lte/gateway/c/core/oai/lib/3gpp/3gpp_23.003.h"
#include "lte/gateway/c/core/oai/lib/3gpp/3gpp_24.008.h"
#include "lte/gateway/c/core/oai/lib/3gpp/3gpp_36.401.h"
#include "lte/gateway/c/core/oai/tasks/mme_app/mme_app_defs.hpp"
#include "lte/gateway/c/core/oai/tasks/mme_app/mme_app_timer.hpp"
#include "lte/gateway/c/core/oai/tasks/nas/emm/EmmCommon.hpp"
#include "lte/gateway/c/core/oai/tasks/nas/emm/emm_data.hpp"
#include "lte/gateway/c/core/oai/tasks/nas/emm/emm_proc.hpp"
#include "lte/gateway/c/core/oai/tasks/nas/emm/msg/emm_cause.hpp"
#include "lte/gateway/c/core/oai/tasks/nas/emm/sap/emm_asDef.hpp"
#include "lte/gateway/c/core/oai/tasks/nas/emm/sap/emm_cnDef.hpp"
#include "lte/gateway/c/core/oai/tasks/nas/emm/sap/emm_fsm.hpp"
#include "lte/gateway/c/core/oai/tasks/nas/emm/sap/emm_regDef.hpp"
#include "lte/gateway/c/core/oai/tasks/nas/emm/sap/emm_sap.hpp"
#include "lte/gateway/c/core/oai/tasks/nas/nas_proc.hpp"
#include "lte/gateway/c/core/oai/tasks/nas/nas_procedures.hpp"
#include "lte/gateway/c/core/oai/tasks/nas/util/nas_timer.hpp"
#include "orc8r/gateway/c/common/service303/MetricsHelpers.hpp"

/****************************************************************************/
/****************  E X T E R N A L    D E F I N I T I O N S  ****************/
/****************************************************************************/
extern long mme_app_last_msg_latency;
extern long pre_mme_task_msg_latency;
extern bool mme_congestion_control_enabled;
extern mme_congestion_params_t mme_congestion_params;

extern int check_plmn_restriction(imsi_t imsi);
extern int validate_imei(imeisv_t* imeisv);
/****************************************************************************/
/*******************  L O C A L    D E F I N I T I O N S  *******************/
/****************************************************************************/

/* String representation of the requested identity type */
static const char* emm_identity_type_str[] = {"NOT AVAILABLE", "IMSI", "IMEI",
                                              "IMEISV", "TMSI"};

// callbacks for identification procedure
static int identification_ll_failure(struct emm_context_s* emm_context,
                                     struct nas_emm_proc_s* emm_proc);
static int identification_non_delivered_ho(struct emm_context_s* emm_context,
                                           struct nas_emm_proc_s* emm_proc);
static int identification_abort(struct emm_context_s* emm_context,
                                struct nas_base_proc_s* base_proc);

static status_code_e identification_request(nas_emm_ident_proc_t* const proc);

/****************************************************************************/
/******************  E X P O R T E D    F U N C T I O N S  ******************/
/****************************************************************************/

/*
   --------------------------------------------------------------------------
        Identification procedure executed by the MME
   --------------------------------------------------------------------------
*/
/********************************************************************
 **                                                                **
 ** Name:    emm_proc_identification()                             **
 **                                                                **
 ** Description: Initiates an identification procedure.            **
 **                                                                **
 **              3GPP TS 24.301, section 5.4.4.2                   **
 **      The network initiates the identification procedure by     **
 **      sending an IDENTITY REQUEST message to the UE and star-   **
 **      ting the timer T3470. The IDENTITY REQUEST message speci- **
 **      fies the requested identification parameters in the Iden- **
 **      tity type information element.                            **
 **                                                                **
 ** Inputs:  ue_id:      UE lower layer identifier                  **
 **      type:      Type of the requested identity                 **
 **      success:   Callback function executed when the identi-    **
 **             fication procedure successfully completes          **
 **      reject:    Callback function executed when the identi-    **
 **             fication procedure fails or is rejected            **
 **      failure:   Callback function executed whener a lower      **
 **             layer failure occured before the identifi-         **
 **             cation procedure completes                         **
 **      Others:    None                                           **
 **                                                                **
 ** Outputs:     None                                              **
 **      Return:    RETURNok, RETURNerror                          **
 **      Others:    _emm_data                                      **
 **                                                                **
 ********************************************************************/
status_code_e emm_proc_identification(struct emm_context_s* const emm_context,
                                      nas_emm_proc_t* const emm_proc,
                                      const identity_type2_t type,
                                      success_cb_t success,
                                      failure_cb_t failure) {
  OAILOG_FUNC_IN(LOG_NAS_EMM);
  status_code_e rc = RETURNerror;

  if ((emm_context) && ((EMM_DEREGISTERED == emm_context->_emm_fsm_state) ||
                        (EMM_REGISTERED == emm_context->_emm_fsm_state))) {
    REQUIREMENT_3GPP_24_301(R10_5_4_4_1);
    mme_ue_s1ap_id_t ue_id =
        PARENT_STRUCT(emm_context, struct ue_mm_context_s, emm_context)
            ->mme_ue_s1ap_id;

    OAILOG_INFO(
        LOG_NAS_EMM,
        "EMM-PROC  - Initiate identification type = %s (%d), ctx = %p for ue "
        "id " MME_UE_S1AP_ID_FMT "\n",
        emm_identity_type_str[type], type, emm_context, ue_id);

    nas_emm_ident_proc_t* ident_proc =
        nas_new_identification_procedure(emm_context);
    if (ident_proc) {
      if (emm_proc) {
        if ((NAS_EMM_PROC_TYPE_SPECIFIC == emm_proc->type) &&
            (EMM_SPEC_PROC_TYPE_ATTACH ==
             ((nas_emm_specific_proc_t*)emm_proc)->type)) {
          ident_proc->is_cause_is_attach = true;
        }
      }
      ident_proc->identity_type = type;
      ident_proc->retransmission_count = 0;
      ident_proc->ue_id = ue_id;
      ((nas_base_proc_t*)ident_proc)->parent = (nas_base_proc_t*)emm_proc;
      ident_proc->emm_com_proc.emm_proc.delivered = NULL;
      ident_proc->emm_com_proc.emm_proc.previous_emm_fsm_state =
          emm_fsm_get_state(emm_context);
      ident_proc->emm_com_proc.emm_proc.not_delivered =
          (sdu_out_not_delivered_t)identification_ll_failure;
      ident_proc->emm_com_proc.emm_proc.not_delivered_ho =
          (sdu_out_not_delivered_ho_t)identification_non_delivered_ho;
      ident_proc->emm_com_proc.emm_proc.base_proc.success_notif = success;
      ident_proc->emm_com_proc.emm_proc.base_proc.failure_notif = failure;
      ident_proc->emm_com_proc.emm_proc.base_proc.abort =
          (proc_abort_t)identification_abort;
      ident_proc->emm_com_proc.emm_proc.base_proc.fail_in =
          NULL;  // only response
      ident_proc->emm_com_proc.emm_proc.base_proc.time_out =
          (time_out_t)mme_app_handle_identification_t3470_expiry;
    }

    rc = identification_request(ident_proc);

    if (rc != RETURNerror) {
      /*
       * Notify EMM that common procedure has been initiated
       */
      emm_sap_t emm_sap = {};

      emm_sap.primitive = EMMREG_COMMON_PROC_REQ;
      emm_sap.u.emm_reg.ue_id = ue_id;
      emm_sap.u.emm_reg.ctx = emm_context;
      rc = emm_sap_send(&emm_sap);
    }
  }

  OAILOG_FUNC_RETURN(LOG_NAS_EMM, rc);
}

/****************************************************************************
 **                                                                        **
 ** Name:    emm_proc_identification_complete()                            **
 **                                                                        **
 ** Description: Performs the identification completion procedure executed **
 **      by the network.                                                   **
 **                                                                        **
 **              3GPP TS 24.301, section 5.4.4.4                           **
 **      Upon receiving the IDENTITY RESPONSE message, the MME             **
 **      shall stop timer T3470.                                           **
 **                                                                        **
 ** Inputs:  ue_id:      UE lower layer identifier                          **
 **      imsi:      The IMSI received from the UE                          **
 **      imei:      The IMEI received from the UE                          **
 **      tmsi:      The TMSI received from the UE                          **
 **      Others:    None                                                   **
 **                                                                        **
 ** Outputs:     None                                                      **
 **      Return:    RETURNok, RETURNerror                                  **
 **      Others:    _emm_data, T3470                                       **
 **                                                                        **
 ***************************************************************************/
status_code_e emm_proc_identification_complete(const mme_ue_s1ap_id_t ue_id,
                                               imsi_t* const imsi,
                                               imei_t* const imei,
                                               imeisv_t* const imeisv,
                                               uint32_t* const tmsi) {
  OAILOG_FUNC_IN(LOG_NAS_EMM);
  status_code_e rc = RETURNerror;
  emm_sap_t emm_sap = {};
  emm_context_t* emm_ctx = NULL;
  bool notify = true;

  OAILOG_INFO(LOG_NAS_EMM,
              "EMM-PROC  - Identification complete (ue_id=" MME_UE_S1AP_ID_FMT
              ")\n",
              ue_id);

  // Get the UE context
  ue_mm_context_t* ue_mm_context = mme_ue_context_exists_mme_ue_s1ap_id(ue_id);
  if (ue_mm_context) {
    emm_ctx = &ue_mm_context->emm_context;
    nas_emm_ident_proc_t* ident_proc =
        get_nas_common_procedure_identification(emm_ctx);

    if (ident_proc) {
      /* Process identification complete msg only if T3470 timer is running.
       * If it is not running it means that response was already received for
       * an earlier attempt.
       */
      if (ident_proc->T3470.id == NAS_TIMER_INACTIVE_ID) {
        OAILOG_WARNING_UE(
            LOG_NAS_EMM, emm_ctx->_imsi64,
            "Discarding identification complete as T3470 timer is not active "
            "for ueid " MME_UE_S1AP_ID_FMT "\n",
            ue_id);
        OAILOG_FUNC_RETURN(LOG_NAS_EMM, RETURNok);
      }

      /* If spent too much in ZMQ, then discard the packet.
       * MME is congested and this would create some relief in processing.
       */
      if (mme_congestion_control_enabled &&
          (mme_app_last_msg_latency + pre_mme_task_msg_latency >
           MME_APP_ZMQ_LATENCY_IDENT_TH)) {
        OAILOG_WARNING_UE(
            LOG_NAS_EMM, emm_ctx->_imsi64,
            "Discarding identification complete as cumulative ZMQ latency "
            "( %ld + %ld ) for ueid " MME_UE_S1AP_ID_FMT
            " is higher than the threshold.",
            mme_app_last_msg_latency, pre_mme_task_msg_latency, ue_id);
        OAILOG_FUNC_RETURN(LOG_NAS_EMM, RETURNok);
      }

      REQUIREMENT_3GPP_24_301(R10_5_4_4_4);
      /*
       * Stop timer T3470
       */
      nas_stop_T3470(ue_id, &ident_proc->T3470);

      if (imsi) {
        imsi64_t imsi64 = imsi_to_imsi64(imsi);
        // If context already exists for this IMSI, perform implicit detach
        mme_app_desc_t* mme_app_desc_p = get_mme_nas_state(false);
        ue_mm_context_t* old_imsi_ue_mm_ctx = mme_ue_context_exists_imsi(
            &mme_app_desc_p->mme_ue_contexts, imsi64);
        if ((emm_ctx->emm_context_state == UNKNOWN_GUTI) &&
            old_imsi_ue_mm_ctx) {
          OAILOG_INFO_UE(LOG_NAS_EMM, imsi64,
                         "EMMAS-SAP - UE context already exists for for ue_id "
                         "=." MME_UE_S1AP_ID_FMT
                         " Triggering implicit detach\n",
                         ue_id);
          nas_emm_attach_proc_t* attach_proc =
              get_nas_specific_procedure_attach(emm_ctx);
          if (!attach_proc) {
            OAILOG_ERROR_UE(
                LOG_NAS_EMM, imsi64,
                "EMMAS-SAP - Attach procedure does not exist for ue_id "
                "=" MME_UE_S1AP_ID_FMT "\n",
                ue_id);
            OAILOG_FUNC_RETURN(LOG_NAS_EMM, RETURNerror);
          }
          create_new_attach_info(&old_imsi_ue_mm_ctx->emm_context,
                                 ue_mm_context->mme_ue_s1ap_id,
                                 STOLEN_REF attach_proc->ies, true);
          emm_ctx->emm_context_state = NEW_EMM_CONTEXT_CREATED;
          nas_proc_implicit_detach_ue_ind(old_imsi_ue_mm_ctx->mme_ue_s1ap_id);
          notify = false;
        }
        int emm_cause = check_plmn_restriction(*imsi);
        if (emm_cause != EMM_CAUSE_SUCCESS) {
          OAILOG_ERROR_UE(
              LOG_NAS_EMM, imsi64,
              "EMMAS-SAP - Sending Attach Reject for ue_id =" MME_UE_S1AP_ID_FMT
              ", emm_cause (%d)\n",
              ue_id, emm_cause);
          rc = emm_proc_attach_reject(ue_id, emm_cause);
          OAILOG_FUNC_RETURN(LOG_NAS_EMM, rc);
        }

        /*
         * Update the IMSI
         */
        emm_ctx_set_valid_imsi(emm_ctx, imsi, imsi64);
        emm_context_upsert_imsi(&_emm_data, emm_ctx);
      } else if (imei) {
        /*
         * Update the IMEI
         */
        emm_ctx_set_valid_imei(emm_ctx, imei);
      } else if (imeisv) {
        // Validate IMEI
        int emm_cause = validate_imei(imeisv);
        if (emm_cause != EMM_CAUSE_SUCCESS) {
          OAILOG_ERROR(
              LOG_NAS_EMM,
              "EMMAS-SAP - Sending Attach Reject for ue_id =" MME_UE_S1AP_ID_FMT
              " , emm_cause =(%d)\n",
              ue_id, emm_cause);
          rc = emm_proc_attach_reject(ue_id, emm_cause);
          OAILOG_FUNC_RETURN(LOG_NAS_EMM, rc);
        }
        // Update the IMEISV
        emm_ctx_set_valid_imeisv(emm_ctx, imeisv);
      } else if (tmsi) {
        /*
         * Update the GUTI
         */
        OAILOG_ERROR(LOG_NAS_EMM,
                     "EMMAS-SAP - Received TMSI in Identication rsp for ue_id "
                     "=" MME_UE_S1AP_ID_FMT ", This case is not handled!\n",
                     ue_id);
      }

      // Helper ident proc ptr to avoid double free from unknown GUTI attach
      // processing.
      nas_emm_ident_proc_t* ident_proc_p =
          reinterpret_cast<nas_emm_ident_proc_t*>(
              calloc(1, sizeof(nas_emm_ident_proc_t)));
      memcpy(ident_proc_p, ident_proc, sizeof(nas_emm_ident_proc_t));

      /*
       * Notify EMM that the identification procedure successfully completed
       */
      emm_sap.primitive = EMMREG_COMMON_PROC_CNF;
      emm_sap.u.emm_reg.ue_id = ue_id;
      emm_sap.u.emm_reg.ctx = emm_ctx;
      emm_sap.u.emm_reg.notify = notify;
      emm_sap.u.emm_reg.free_proc = true;
      emm_sap.u.emm_reg.u.common.common_proc = &ident_proc_p->emm_com_proc;
      emm_sap.u.emm_reg.u.common.previous_emm_fsm_state =
          ident_proc_p->emm_com_proc.emm_proc.previous_emm_fsm_state;
      rc = emm_sap_send(&emm_sap);

    }  // else ignore the response if procedure not found
  }  // else ignore the response if ue context not found

  OAILOG_FUNC_RETURN(LOG_NAS_EMM, rc);
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
 ** Name:    mme_app_handle_identification_t3470_expiry() **
 **                                                                        **
 ** Description: T3470 timeout handler                                     **
 **      Upon T3470 timer expiration, the identification request   **
 **      message is retransmitted and the timer restarted. When    **
 **      retransmission counter is exceed, the MME shall abort the **
 **      identification procedure and any ongoing EMM procedure.   **
 **                                                                        **
 **              3GPP TS 24.301, section 5.4.4.6, case b                   **
 **                                                                        **
 ** Inputs:  args:      handler parameters                         **
 **      Others:    None                                       **
 **                                                                        **
 ** Outputs:     None                                                      **
 **      Return:    None                                       **
 **      Others:    None                                       **
 **                                                                        **
 ***************************************************************************/
status_code_e mme_app_handle_identification_t3470_expiry(zloop_t* loop,
                                                         int timer_id,
                                                         void* args) {
  OAILOG_FUNC_IN(LOG_NAS_EMM);
  mme_ue_s1ap_id_t mme_ue_s1ap_id = 0;
  if (!mme_pop_timer_arg_ue_id(timer_id, &mme_ue_s1ap_id)) {
    OAILOG_WARNING(LOG_NAS_EMM, "Invalid Timer Id expiration, Timer Id: %u\n",
                   timer_id);
    OAILOG_FUNC_RETURN(LOG_NAS_EMM, RETURNok);
  }

  struct ue_mm_context_s* ue_context_p = mme_app_get_ue_context_for_timer(
      mme_ue_s1ap_id, const_cast<char*>("Identification T3470 Timer"));
  if (ue_context_p == NULL) {
    OAILOG_ERROR(
        LOG_MME_APP,
        "Invalid UE context received, MME UE S1AP Id: " MME_UE_S1AP_ID_FMT "\n",
        mme_ue_s1ap_id);
    OAILOG_FUNC_RETURN(LOG_NAS_EMM, RETURNok);
  }

  emm_context_t* emm_ctx = &ue_context_p->emm_context;

  if (!(emm_ctx)) {
    OAILOG_ERROR(LOG_NAS_EMM, "T3470 timer expired No EMM context\n");
    OAILOG_FUNC_RETURN(LOG_NAS_EMM, RETURNok);
  }
  nas_emm_ident_proc_t* ident_proc =
      get_nas_common_procedure_identification(emm_ctx);

  if (ident_proc) {
    OAILOG_WARNING(LOG_NAS_EMM,
                   "T3470 timer (%lx) expired ue id " MME_UE_S1AP_ID_FMT " \n",
                   ident_proc->T3470.id, ident_proc->ue_id);
    ident_proc->T3470.id = NAS_TIMER_INACTIVE_ID;
    /*
     * Increment the retransmission counter
     */
    ident_proc->retransmission_count += 1;
    OAILOG_WARNING(LOG_NAS_EMM,
                   "EMM-PROC  - T3470 (%lx) retransmission counter = %d ue "
                   "id " MME_UE_S1AP_ID_FMT " \n",
                   ident_proc->T3470.id, ident_proc->retransmission_count,
                   ident_proc->ue_id);

    if (ident_proc->retransmission_count < IDENTIFICATION_COUNTER_MAX) {
      REQUIREMENT_3GPP_24_301(R10_5_4_4_6_b__1);
      /*
       * Send identity request message to the UE
       */
      identification_request(ident_proc);
    } else {
      /*
       * Abort the identification procedure
       */
      mme_ue_s1ap_id_t ue_id = ident_proc->ue_id;
      REQUIREMENT_3GPP_24_301(R10_5_4_4_6_b__2);
      emm_sap_t emm_sap = {};
      emm_sap.primitive = EMMREG_COMMON_PROC_ABORT;
      emm_sap.u.emm_reg.ue_id = ident_proc->ue_id;
      emm_sap.u.emm_reg.ctx = emm_ctx;
      emm_sap.u.emm_reg.notify = false;
      emm_sap.u.emm_reg.free_proc = true;
      emm_sap.u.emm_reg.u.common.common_proc =
          (nas_emm_common_proc_t*)(&ident_proc->emm_com_proc);
      emm_sap.u.emm_reg.u.common.previous_emm_fsm_state =
          ((nas_emm_proc_t*)ident_proc)->previous_emm_fsm_state;
      emm_sap_send(&emm_sap);
      nas_delete_all_emm_procedures(emm_ctx);
      /* clear emm_common_data_ctx */
      emm_common_cleanup_by_ueid(ue_id);
      memset((void*)&emm_sap, 0, sizeof(emm_sap));
      emm_sap.primitive = EMMCN_IMPLICIT_DETACH_UE;
      emm_sap.u.emm_cn.u.emm_cn_implicit_detach.ue_id = ue_id;
      emm_sap_send(&emm_sap);
      increment_counter("ue_attach", 1, 1, "action", "attach_abort");
    }
  } else {
    OAILOG_ERROR(LOG_NAS_EMM,
                 "T3470 timer expired, No Identification procedure found\n");
  }

  OAILOG_FUNC_RETURN(LOG_NAS_EMM, RETURNok);
}

/*
   --------------------------------------------------------------------------
                MME specific local functions
   --------------------------------------------------------------------------
*/

/*
 * Description: Sends IDENTITY REQUEST message and start timer T3470.
 *
 * Inputs:  args:      handler parameters
 *      Others:    None
 *
 * Outputs:     None
 *      Return:    RETURNok, RETURNerror
 *      Others:    T3470
 */
static status_code_e identification_request(nas_emm_ident_proc_t* const proc) {
  OAILOG_FUNC_IN(LOG_NAS_EMM);
  emm_sap_t emm_sap = {};
  status_code_e rc = RETURNok;
  struct emm_context_s* emm_ctx = NULL;

  ue_mm_context_t* ue_mm_context =
      mme_ue_context_exists_mme_ue_s1ap_id(proc->ue_id);
  if (ue_mm_context) {
    emm_ctx = &ue_mm_context->emm_context;
  } else {
    OAILOG_FUNC_RETURN(LOG_NAS_EMM, RETURNerror);
  }
  /*
   * Notify EMM-AS SAP that Identity Request message has to be sent
   * to the UE
   */
  emm_sap.primitive = EMMAS_SECURITY_REQ;
  emm_sap.u.emm_as.u.security.puid =
      proc->emm_com_proc.emm_proc.base_proc.nas_puid;
  emm_sap.u.emm_as.u.security.guti = NULL;
  emm_sap.u.emm_as.u.security.ue_id = proc->ue_id;
  emm_sap.u.emm_as.u.security.msg_type = EMM_AS_MSG_TYPE_IDENT;
  emm_sap.u.emm_as.u.security.ident_type = proc->identity_type;

  /*
   * Setup EPS NAS security data
   */
  emm_as_set_security_data(&emm_sap.u.emm_as.u.security.sctx,
                           &emm_ctx->_security, false, true);
  rc = emm_sap_send(&emm_sap);

  if (rc != RETURNerror) {
    REQUIREMENT_3GPP_24_301(R10_5_4_4_2);
    /*
     * Start T3470 timer
     */
    nas_start_T3470(proc->ue_id, &proc->T3470,
                    proc->emm_com_proc.emm_proc.base_proc.time_out);
  }

  OAILOG_FUNC_RETURN(LOG_NAS_EMM, rc);
}

//------------------------------------------------------------------------------
static int identification_ll_failure(struct emm_context_s* emm_ctx,
                                     struct nas_emm_proc_s* emm_proc) {
  OAILOG_FUNC_IN(LOG_NAS_EMM);
  int rc = RETURNerror;
  if ((emm_ctx) && (emm_proc)) {
    nas_emm_ident_proc_t* ident_proc = (nas_emm_ident_proc_t*)emm_proc;
    REQUIREMENT_3GPP_24_301(R10_5_4_4_6_a);
    emm_sap_t emm_sap = {};

    emm_sap.primitive = EMMREG_COMMON_PROC_ABORT;
    emm_sap.u.emm_reg.ue_id = ident_proc->ue_id;
    emm_sap.u.emm_reg.ctx = emm_ctx;
    rc = emm_sap_send(&emm_sap);
    nas_delete_all_emm_procedures(emm_ctx);
  }
  OAILOG_FUNC_RETURN(LOG_NAS_EMM, rc);
}

//------------------------------------------------------------------------------
static int identification_non_delivered_ho(struct emm_context_s* emm_ctx,
                                           struct nas_emm_proc_s* emm_proc) {
  OAILOG_FUNC_IN(LOG_NAS_EMM);
  int rc = RETURNerror;
  if (emm_proc) {
    nas_emm_ident_proc_t* ident_proc = (nas_emm_ident_proc_t*)emm_proc;
    /************************README*******************************************
  ** NAS Non Delivery indication during HO handling will be added when HO is
  ** supported In non hand-over case if MME receives NAS Non Delivery
  ** indication message that implies eNB and UE has lost radio connection.
  ** In this case aborting the Identification and Attach Procedure.
  *****************************************************************************
  REQUIREMENT_3GPP_24_301(R10_5_4_2_7_j);
  ******************************************************************************/
    if (emm_ctx) {
      REQUIREMENT_3GPP_24_301(R10_5_4_2_7_j);
      /*
       * Stop timer T3470
       */
      if (ident_proc->T3470.id != NAS_TIMER_INACTIVE_ID) {
        OAILOG_INFO(
            LOG_NAS_EMM,
            "EMM-PROC  - Stop timer T3460 (%ld) for ue id " MME_UE_S1AP_ID_FMT
            "\n",
            ident_proc->T3470.id, ident_proc->ue_id);
        nas_stop_T3470(ident_proc->ue_id, &ident_proc->T3470);
      }
      /*
       * Abort identification and attach procedure
       */
      emm_sap_t emm_sap = {};
      emm_sap.primitive = EMMREG_COMMON_PROC_ABORT;
      emm_sap.u.emm_reg.ue_id = ident_proc->ue_id;
      emm_sap.u.emm_reg.ctx = emm_ctx;
      emm_sap_send(&emm_sap);
      /* clear emm_common_data_ctx */
      emm_common_cleanup_by_ueid(ident_proc->ue_id);
      // Clean up MME APP UE context
      emm_sap.primitive = EMMCN_IMPLICIT_DETACH_UE;
      emm_sap.u.emm_cn.u.emm_cn_implicit_detach.ue_id = ident_proc->ue_id;
      rc = emm_sap_send(&emm_sap);
    }
  }
  OAILOG_FUNC_RETURN(LOG_NAS_EMM, rc);
}

//------------------------------------------------------------------------------
/*
 * Description: Aborts the identification procedure currently in progress
 *
 * Inputs:  args:      Identification data to be released
 *      Others:    None
 *
 * Outputs:     None
 *      Return:    None
 *      Others:    T3470
 */
static int identification_abort(struct emm_context_s* emm_context,
                                struct nas_base_proc_s* base_proc) {
  OAILOG_FUNC_IN(LOG_NAS_EMM);
  int rc = RETURNerror;

  if ((emm_context) && (base_proc)) {
    nas_emm_ident_proc_t* ident_proc = (nas_emm_ident_proc_t*)base_proc;
    AssertFatal(
        (NAS_PROC_TYPE_EMM == base_proc->type) &&
            (NAS_EMM_PROC_TYPE_COMMON == ((nas_emm_proc_t*)base_proc)->type) &&
            (EMM_COMM_PROC_IDENT == ((nas_emm_common_proc_t*)base_proc)->type),
        "Mismatch in procedure type");

    OAILOG_INFO(LOG_NAS_EMM,
                "EMM-PROC  - Abort identification procedure "
                "(ue_id=" MME_UE_S1AP_ID_FMT ")\n",
                ident_proc->ue_id);

    /*
     * Stop timer T3470
     */
    nas_stop_T3470(ident_proc->ue_id, &ident_proc->T3470);
  }
  OAILOG_FUNC_RETURN(LOG_NAS_EMM, rc);
}

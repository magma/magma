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

/*****************************************************************************
  Source      SecurityModeControl.c

  Version     0.1

  Date        2013/04/22

  Product     NAS stack

  Subsystem   Template body file

  Author      Frederic Maurel

  Description Defines the security mode control EMM procedure executed by the
        Non-Access Stratum.

        The purpose of the NAS security mode control procedure is to
        take an EPS security context into use, and initialise and start
        NAS signalling security between the UE and the MME with the
        corresponding EPS NAS keys and EPS security algorithms.

        Furthermore, the network may also initiate a SECURITY MODE COM-
        MAND in order to change the NAS security algorithms for a cur-
        rent EPS security context already in use.

*****************************************************************************/

#include <stdint.h>
#include <stdbool.h>
#include <stdlib.h>

#include "log.h"
#include "assertions.h"
#include "3gpp_requirements_24.301.h"
#include "common_types.h"
#include "3gpp_24.008.h"
#include "mme_app_ue_context.h"
#include "emm_proc.h"
#include "common_defs.h"
#include "nas_timer.h"
#include "emm_data.h"
#include "emm_sap.h"
#include "secu_defs.h"
#include "nas_itti_messaging.h"
#include "service303.h"
#include "EmmCommon.h"
#include "3gpp_23.003.h"
#include "3gpp_24.301.h"
#include "3gpp_33.401.h"
#include "3gpp_36.401.h"
#include "NasSecurityAlgorithms.h"
#include "emm_asDef.h"
#include "emm_cnDef.h"
#include "emm_fsm.h"
#include "emm_regDef.h"
#include "mme_api.h"
#include "mme_app_state.h"
#include "nas_procedures.h"
#include "nas/securityDef.h"
#include "security_types.h"
#include "mme_app_defs.h"

/****************************************************************************/
/****************  E X T E R N A L    D E F I N I T I O N S  ****************/
/****************************************************************************/

/****************************************************************************/
/*******************  L O C A L    D E F I N I T I O N S  *******************/
/****************************************************************************/

/*
   --------------------------------------------------------------------------
    Internal data handled by the security mode control procedure in the UE
   --------------------------------------------------------------------------
*/

/*
   --------------------------------------------------------------------------
    Internal data handled by the security mode control procedure in the MME
   --------------------------------------------------------------------------
*/
/*
   Timer handlers
*/
static void _security_t3460_handler(void *);
static int _security_ll_failure(
  emm_context_t *emm_context,
  struct nas_emm_proc_s *nas_emm_proc);
static int _security_non_delivered_ho(
  emm_context_t *emm_context,
  struct nas_emm_proc_s *nas_emm_proc);

/*
   Function executed whenever the ongoing EMM procedure that initiated
   the security mode control procedure is aborted or the maximum value of the
   retransmission timer counter is exceed
*/
static int _security_abort(
  emm_context_t *emm_context,
  struct nas_base_proc_s *base_proc);
static int _security_select_algorithms(
  const int ue_eiaP,
  const int ue_eeaP,
  int *const mme_eiaP,
  int *const mme_eeaP);

static int _security_request(nas_emm_smc_proc_t *const smc_proc);

/****************************************************************************/
/******************  E X P O R T E D    F U N C T I O N S  ******************/
/****************************************************************************/

/*
   --------------------------------------------------------------------------
        Security mode control procedure executed by the MME
   --------------------------------------------------------------------------
*/

/****************************************************************************
 **                                                                        **
 ** Name:    emm_proc_security_mode_control()                          **
 **                                                                        **
 ** Description: Initiates the security mode control procedure.            **
 **                                                                        **
 **              3GPP TS 24.301, section 5.4.3.2                           **
 **      The MME initiates the NAS security mode control procedure **
 **      by sending a SECURITY MODE COMMAND message to the UE and  **
 **      starting timer T3460. The message shall be sent unciphe-  **
 **      red but shall be integrity protected using the NAS inte-  **
 **      grity key based on KASME.                                 **
 **                                                                        **
 ** Inputs:  ue_id:      UE lower layer identifier                  **
 **      ksi:       NAS key set identifier                     **
 **      eea:       Replayed EPS encryption algorithms         **
 **      eia:       Replayed EPS integrity algorithms          **
 **      success:   Callback function executed when the secu-  **
 **             rity mode control procedure successfully   **
 **             completes                                  **
 **      reject:    Callback function executed when the secu-  **
 **             rity mode control procedure fails or is    **
 **             rejected                                   **
 **      failure:   Callback function executed whener a lower  **
 **             layer failure occured before the security  **
 **             mode control procedure completes          **
 **      Others:    None                                       **
 **                                                                        **
 ** Outputs:     None                                                      **
 **      Return:    RETURNok, RETURNerror                      **
 **      Others:    None                                       **
 **                                                                        **
 ***************************************************************************/
int emm_proc_security_mode_control(
  struct emm_context_s *emm_ctx,
  nas_emm_specific_proc_t *const emm_specific_proc,
  ksi_t ksi,
  success_cb_t success,
  failure_cb_t failure)
{
  OAILOG_FUNC_IN(LOG_NAS_EMM);
  int rc = RETURNerror;
  bool security_context_is_new = false;
  int mme_eea = NAS_SECURITY_ALGORITHMS_EEA0;
  int mme_eia = NAS_SECURITY_ALGORITHMS_EIA0;
  /*
   * Get the UE context
   */

  OAILOG_INFO(
    LOG_NAS_EMM,
    "EMM-PROC  - Initiate security mode control procedure, "
    "KSI = %d\n",
    ksi);

  if (!(emm_ctx)) {
    OAILOG_ERROR(
      LOG_NAS_EMM,
      "Emm Context NULL!\n");
    OAILOG_FUNC_RETURN(LOG_NAS_EMM, RETURNerror);
  }

  //TODO better than that (quick fixes)
  if (KSI_NO_KEY_AVAILABLE == ksi) {
    ksi = 0;
  }
  if (EMM_SECURITY_VECTOR_INDEX_INVALID == emm_ctx->_security.vector_index) {
    emm_ctx_set_security_vector_index(emm_ctx, 0);
  }
  /*
   * Allocate parameters of the retransmission timer callback
   */
  mme_ue_s1ap_id_t ue_id =
    PARENT_STRUCT(emm_ctx, struct ue_mm_context_s, emm_context)->mme_ue_s1ap_id;
  nas_emm_smc_proc_t *smc_proc = get_nas_common_procedure_smc(emm_ctx);
  if (!smc_proc) {
    smc_proc = nas_new_smc_procedure(emm_ctx);
    if (smc_proc) {
      // TODO check for removing test (emm_ctx->_security.sc_type == SECURITY_CTX_TYPE_NOT_AVAILABLE)
      //if ((emm_ctx->_security.sc_type == SECURITY_CTX_TYPE_NOT_AVAILABLE) &&

      smc_proc->saved_selected_eea =
        emm_ctx->_security.selected_algorithms.encryption;
      smc_proc->saved_selected_eia =
        emm_ctx->_security.selected_algorithms.integrity;
      smc_proc->saved_eksi = emm_ctx->_security.eksi;
      smc_proc->saved_overflow = emm_ctx->_security.dl_count.overflow;
      smc_proc->saved_seq_num = emm_ctx->_security.dl_count.seq_num;
      smc_proc->saved_sc_type = emm_ctx->_security.sc_type;
      /*
       * The security mode control procedure is initiated to take into use
       * * * * the EPS security context created after a successful execution of
       * * * * the EPS authentication procedure
       */
      //emm_ctx->_security.sc_type = SECURITY_CTX_TYPE_PARTIAL_NATIVE;
      emm_ctx_set_security_eksi(emm_ctx, ksi);
      REQUIREMENT_3GPP_24_301(R10_5_4_3_2__2);
      emm_ctx->_security.dl_count.overflow = 0;
      emm_ctx->_security.dl_count.seq_num = 0;

      /*
       *  Compute NAS cyphering and integrity keys
       */

      rc = _security_select_algorithms(
        emm_ctx->_ue_network_capability.eia,
        emm_ctx->_ue_network_capability.eea,
        &mme_eia,
        &mme_eea);
      emm_ctx->_security.selected_algorithms.encryption = mme_eea;
      emm_ctx->_security.selected_algorithms.integrity = mme_eia;

      if (rc == RETURNerror) {
        OAILOG_WARNING(
          LOG_NAS_EMM, "EMM-PROC  - Failed to select security algorithms\n");
        OAILOG_FUNC_RETURN(LOG_NAS_EMM, RETURNerror);
      }

      emm_ctx_set_security_type(emm_ctx, SECURITY_CTX_TYPE_FULL_NATIVE);
      AssertFatal(
        KSI_NO_KEY_AVAILABLE > emm_ctx->_security.eksi, "eksi not valid");
      derive_key_nas(
        NAS_INT_ALG,
        emm_ctx->_security.selected_algorithms.integrity,
        emm_ctx->_vector[emm_ctx->_security.eksi % MAX_EPS_AUTH_VECTORS].kasme,
        emm_ctx->_security.knas_int);
      derive_key_nas(
        NAS_ENC_ALG,
        emm_ctx->_security.selected_algorithms.encryption,
        emm_ctx->_vector[emm_ctx->_security.eksi % MAX_EPS_AUTH_VECTORS].kasme,
        emm_ctx->_security.knas_enc);
      /*
       * Set new security context indicator
       */
      security_context_is_new = true;
      emm_ctx_set_attribute_present(emm_ctx, EMM_CTXT_MEMBER_SECURITY);
    }
  } else {
    OAILOG_ERROR(LOG_NAS_EMM, "EMM-PROC  - No EPS security context exists\n");
    OAILOG_FUNC_RETURN(LOG_NAS_EMM, RETURNerror);
  }

  if (smc_proc) {
    /*
     * Setup ongoing EMM procedure callback functions
     */
    ((nas_base_proc_t *) smc_proc)->parent =
      (nas_base_proc_t *) emm_specific_proc;
    smc_proc->emm_com_proc.emm_proc.delivered = NULL;
    smc_proc->emm_com_proc.emm_proc.previous_emm_fsm_state =
      emm_fsm_get_state(emm_ctx);
    smc_proc->emm_com_proc.emm_proc.not_delivered = _security_ll_failure;
    smc_proc->emm_com_proc.emm_proc.not_delivered_ho =
      _security_non_delivered_ho;
    smc_proc->emm_com_proc.emm_proc.base_proc.success_notif = success;
    smc_proc->emm_com_proc.emm_proc.base_proc.failure_notif = failure;
    smc_proc->emm_com_proc.emm_proc.base_proc.abort = _security_abort;
    smc_proc->emm_com_proc.emm_proc.base_proc.fail_in = NULL; // only response
    smc_proc->emm_com_proc.emm_proc.base_proc.fail_out = NULL;
    smc_proc->emm_com_proc.emm_proc.base_proc.time_out =
      _security_t3460_handler;

    /*
     * Set the UE identifier
     */
    smc_proc->ue_id = ue_id;
    /*
     * Reset the retransmission counter
     */
    smc_proc->retransmission_count = 0;
    /*
     * Set the key set identifier
     */
    smc_proc->ksi = ksi;
    /*
     * Set the EPS encryption algorithms to be replayed to the UE
     */
    smc_proc->eea = emm_ctx->_ue_network_capability.eea;
    /*
     * Set the EPS integrity algorithms to be replayed to the UE
     */
    smc_proc->eia = emm_ctx->_ue_network_capability.eia;
    smc_proc->ucs2 = emm_ctx->_ue_network_capability.ucs2;
    /*
     * Set the UMTS encryption algorithms to be replayed to the UE
     */
    smc_proc->uea = emm_ctx->_ue_network_capability.uea;
    /*
     * Set the UMTS integrity algorithms to be replayed to the UE
     */
    smc_proc->uia = emm_ctx->_ue_network_capability.uia;
    /*
     * Set the GPRS integrity algorithms to be replayed to the UE
     */
    uint8_t gea = emm_ctx->_ms_network_capability.gea1;
    if (gea) {
      gea = (gea << 6) | emm_ctx->_ms_network_capability.egea;
    }
    smc_proc->gea = gea;
    smc_proc->umts_present = emm_ctx->_ue_network_capability.umts_present;
    smc_proc->gprs_present = (gea >= (MS_NETWORK_CAPABILITY_GEA1 >> 1));
    /*
     * Set the EPS encryption algorithms selected to the UE
     */
    smc_proc->selected_eea = emm_ctx->_security.selected_algorithms.encryption;
    OAILOG_DEBUG(LOG_NAS_EMM, "EPS encryption algorithm selected is (%d) for ue_id (%u)\n",
      smc_proc->selected_eea, ue_id);
    /*
     * Set the EPS integrity algorithms selected to the UE
     */
    smc_proc->selected_eia = emm_ctx->_security.selected_algorithms.integrity;
    OAILOG_DEBUG(LOG_NAS_EMM, "EPS integrity algorithm selected is (%d) for ue_id (%u)\n",
      smc_proc->selected_eia, ue_id);

    smc_proc->is_new = security_context_is_new;

    // always ask for IMEISV (Do it simple now)
    smc_proc->imeisv_request = true;
    //smc_proc->imeisv_request = (IS_EMM_CTXT_PRESENT_IMEISV(emm_ctx)) ? false:true;

    /*
     * Send security mode command message to the UE
     */
    rc = _security_request(smc_proc);

    if (rc != RETURNerror) {
      /*
       * Notify EMM that common procedure has been initiated
       */
      emm_sap_t emm_sap = {0};

      emm_sap.primitive = EMMREG_COMMON_PROC_REQ;
      emm_sap.u.emm_reg.ue_id = ue_id;
      emm_sap.u.emm_reg.ctx = emm_ctx;
      emm_sap.u.emm_reg.u.common.common_proc = &smc_proc->emm_com_proc;
      emm_sap.u.emm_reg.u.common.previous_emm_fsm_state =
        smc_proc->emm_com_proc.emm_proc.previous_emm_fsm_state;
      rc = emm_sap_send(&emm_sap);
    }
  }

  OAILOG_FUNC_RETURN(LOG_NAS_EMM, rc);
}

/****************************************************************************
 **                                                                        **
 ** Name:    emm_proc_security_mode_complete()                         **
 **                                                                        **
 ** Description: Performs the security mode control completion procedure   **
 **      executed by the network.                                  **
 **                                                                        **
 **              3GPP TS 24.301, section 5.4.3.4                           **
 **      Upon receiving the SECURITY MODE COMPLETE message, the    **
 **      MME shall stop timer T3460.                               **
 **      From this time onward the MME shall integrity protect and **
 **      encipher all signalling messages with the selected NAS    **
 **      integrity and ciphering algorithms.                       **
 **                                                                        **
 ** Inputs:  ue_id:      UE lower layer identifier                  **
 **      Others:    None                                       **
 **                                                                        **
 ** Outputs:     None                                                      **
 **      Return:    RETURNok, RETURNerror                      **
 **      Others:    None                                       **
 **                                                                        **
 ***************************************************************************/
int emm_proc_security_mode_complete(
  mme_ue_s1ap_id_t ue_id,
  const imeisv_mobile_identity_t *const imeisvmob)
{
  OAILOG_FUNC_IN(LOG_NAS_EMM);
  ue_mm_context_t *ue_mm_context = NULL;
  emm_context_t *emm_ctx = NULL;
  int rc = RETURNerror;

  OAILOG_INFO(
    LOG_NAS_EMM,
    "EMM-PROC  - Security mode complete (ue_id=" MME_UE_S1AP_ID_FMT ")\n",
    ue_id);
  /*
   * Get the UE context
   */
  mme_app_desc_t *mme_app_desc_p = get_mme_nas_state(false);
  ue_mm_context = mme_ue_context_exists_mme_ue_s1ap_id(
    &mme_app_desc_p->mme_ue_contexts, ue_id);
  if (ue_mm_context) {
    emm_ctx = &ue_mm_context->emm_context;
  } else {
    OAILOG_FUNC_RETURN(LOG_NAS_EMM, RETURNerror);
  }

  nas_emm_smc_proc_t *smc_proc = get_nas_common_procedure_smc(emm_ctx);

  if (smc_proc) {
    /*
     * Stop timer T3460
     */
    REQUIREMENT_3GPP_24_301(R10_5_4_3_4__1);

    void *timer_callback_arg = NULL;
    nas_stop_T3460(ue_id, &smc_proc->T3460, timer_callback_arg);

    if (imeisvmob) {
      imeisv_t imeisv = {0};
      imeisv.u.num.tac1 = imeisvmob->tac1;
      imeisv.u.num.tac2 = imeisvmob->tac2;
      imeisv.u.num.tac3 = imeisvmob->tac3;
      imeisv.u.num.tac4 = imeisvmob->tac4;
      imeisv.u.num.tac5 = imeisvmob->tac5;
      imeisv.u.num.tac6 = imeisvmob->tac6;
      imeisv.u.num.tac7 = imeisvmob->tac7;
      imeisv.u.num.tac8 = imeisvmob->tac8;
      imeisv.u.num.snr1 = imeisvmob->snr1;
      imeisv.u.num.snr2 = imeisvmob->snr2;
      imeisv.u.num.snr3 = imeisvmob->snr3;
      imeisv.u.num.snr4 = imeisvmob->snr4;
      imeisv.u.num.snr5 = imeisvmob->snr5;
      imeisv.u.num.snr6 = imeisvmob->snr6;
      imeisv.u.num.svn1 = imeisvmob->svn1;
      imeisv.u.num.svn2 = imeisvmob->svn2;
      imeisv.u.num.parity = imeisvmob->oddeven;
      emm_ctx_set_valid_imeisv(emm_ctx, &imeisv);
    }

    /*
     * Release retransmission timer parameters
     */

    if (emm_ctx && IS_EMM_CTXT_PRESENT_SECURITY(emm_ctx)) {
      /*
       * Notify EMM that the authentication procedure successfully completed
       */
      emm_sap_t emm_sap = {0};
      emm_sap.primitive = EMMREG_COMMON_PROC_CNF;
      emm_sap.u.emm_reg.ue_id = ue_id;
      emm_sap.u.emm_reg.ctx = emm_ctx;
      emm_sap.u.emm_reg.notify = true;
      emm_sap.u.emm_reg.free_proc = true;
      emm_sap.u.emm_reg.u.common.common_proc = &smc_proc->emm_com_proc;
      emm_sap.u.emm_reg.u.common.previous_emm_fsm_state =
        smc_proc->emm_com_proc.emm_proc.previous_emm_fsm_state;
      REQUIREMENT_3GPP_24_301(R10_5_4_3_4__2);

      emm_ctx->_security.kenb_ul_count = emm_ctx->_security.ul_count;
      emm_ctx_set_attribute_valid(emm_ctx, EMM_CTXT_MEMBER_SECURITY);
      rc = emm_sap_send(&emm_sap);
    }
    unlock_ue_contexts(ue_mm_context);
    OAILOG_FUNC_RETURN(LOG_NAS_EMM, rc);
  } else {
    OAILOG_ERROR(
      LOG_NAS_EMM,
      "EMM-PROC  - No EPS security context exists. Ignoring the Security Mode "
      "Complete message\n");
    rc = RETURNerror;
  }

  unlock_ue_contexts(ue_mm_context);
  OAILOG_FUNC_RETURN(LOG_NAS_EMM, rc);
}

/****************************************************************************
 **                                                                        **
 ** Name:    emm_proc_security_mode_reject()                           **
 **                                                                        **
 ** Description: Performs the security mode control not accepted by the UE **
 **                                                                        **
 **              3GPP TS 24.301, section 5.4.3.5                           **
 **      Upon receiving the SECURITY MODE REJECT message, the MME  **
 **      shall stop timer T3460 and abort the ongoing procedure    **
 **      that triggered the initiation of the NAS security mode    **
 **      control procedure.                                        **
 **      The MME shall apply the EPS security context in use befo- **
 **      re the initiation of the security mode control procedure, **
 **      if any, to protect any subsequent messages.               **
 **                                                                        **
 **                                                                        **
 ** Inputs:  ue_id:      UE lower layer identifier                  **
 **      Others:    None                                       **
 **                                                                        **
 ** Outputs:     None                                                      **
 **      Return:    RETURNok, RETURNerror                      **
 **      Others:    None                                       **
 **                                                                        **
 ***************************************************************************/
int emm_proc_security_mode_reject(mme_ue_s1ap_id_t ue_id)
{
  OAILOG_FUNC_IN(LOG_NAS_EMM);
  ue_mm_context_t *ue_mm_context = NULL;
  emm_context_t *emm_ctx = NULL;
  int rc = RETURNerror;

  OAILOG_WARNING(
    LOG_NAS_EMM,
    "EMM-PROC  - Security mode command not accepted by the UE"
    "(ue_id=" MME_UE_S1AP_ID_FMT ")\n",
    ue_id);
  /*
   * Get the UE context
   */

  mme_app_desc_t *mme_app_desc_p = get_mme_nas_state(false);
  ue_mm_context = mme_ue_context_exists_mme_ue_s1ap_id(
    &mme_app_desc_p->mme_ue_contexts, ue_id);
  if (ue_mm_context) {
    emm_ctx = &ue_mm_context->emm_context;
  } else {
    OAILOG_FUNC_RETURN(LOG_NAS_EMM, RETURNerror);
  }

  nas_emm_smc_proc_t *smc_proc = get_nas_common_procedure_smc(emm_ctx);

  if (smc_proc) {
    /*
     * Stop timer T3460
     */
    REQUIREMENT_3GPP_24_301(R10_5_4_3_5__2);
    void *timer_callback_arg = NULL;
    nas_stop_T3460(ue_id, &smc_proc->T3460, timer_callback_arg);

    // restore previous values
    REQUIREMENT_3GPP_24_301(R10_5_4_3_5__3);
    emm_ctx->_security.selected_algorithms.encryption =
      smc_proc->saved_selected_eea;
    emm_ctx->_security.selected_algorithms.integrity =
      smc_proc->saved_selected_eia;
    emm_ctx_set_security_eksi(emm_ctx, smc_proc->saved_eksi);
    emm_ctx->_security.dl_count.overflow = smc_proc->saved_overflow;
    emm_ctx->_security.dl_count.seq_num = smc_proc->saved_seq_num;
    emm_ctx_set_security_type(emm_ctx, smc_proc->saved_sc_type);

    /*
     * Notify EMM that the security mode procedure failed
     */
    emm_sap_t emm_sap = {0};

    REQUIREMENT_3GPP_24_301(R10_5_4_3_5__2);
    emm_sap.primitive = EMMREG_COMMON_PROC_REJ;
    emm_sap.u.emm_reg.ue_id = ue_id;
    emm_sap.u.emm_reg.ctx = emm_ctx;
    emm_sap.u.emm_reg.notify = true;
    emm_sap.u.emm_reg.free_proc = false;
    emm_sap.u.emm_reg.u.common.common_proc = &smc_proc->emm_com_proc;
    emm_sap.u.emm_reg.u.common.previous_emm_fsm_state =
      smc_proc->emm_com_proc.emm_proc.previous_emm_fsm_state;
    rc = emm_sap_send(&emm_sap);
  }
  unlock_ue_contexts(ue_mm_context);
  mme_app_handle_detach_req(ue_id);
  OAILOG_FUNC_RETURN(LOG_NAS_EMM, rc);
}

/**
 * When the NAS security procedures are restored from data store, the
 * references to callback functions need to be re-populated with the local
 * function pointers. The functions below populate these callbacks from
 * security mode controle procedure.
 * The memory for each procedure is allocated by the caller
 */

void set_callbacks_for_smc_proc(nas_emm_smc_proc_t *smc_proc)
{
  smc_proc->emm_com_proc.emm_proc.not_delivered =
    _security_ll_failure;
  smc_proc->emm_com_proc.emm_proc.not_delivered_ho =
    _security_non_delivered_ho;
  smc_proc->emm_com_proc.emm_proc.base_proc.abort = _security_abort;
  smc_proc->emm_com_proc.emm_proc.base_proc.fail_in =
    NULL;
  smc_proc->emm_com_proc.emm_proc.base_proc.fail_out =
    NULL;
  smc_proc->emm_com_proc.emm_proc.base_proc.time_out =
    _security_t3460_handler;
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
 ** Name:    _security_t3460_handler()                                 **
 **                                                                        **
 ** Description: T3460 timeout handler                                     **
 **      Upon T3460 timer expiration, the security mode command    **
 **      message is retransmitted and the timer restarted. When    **
 **      retransmission counter is exceed, the MME shall abort the **
 **      security mode control procedure.                          **
 **                                                                        **
 **              3GPP TS 24.301, section 5.4.3.7, case b                   **
 **                                                                        **
 ** Inputs:  args:      handler parameters                         **
 **      Others:    None                                       **
 **                                                                        **
 ** Outputs:     None                                                      **
 **      Return:    None                                       **
 **      Others:    None                                       **
 **                                                                        **
 ***************************************************************************/
static void _security_t3460_handler(void *args)
{
  OAILOG_FUNC_IN(LOG_NAS_EMM);
  emm_context_t *emm_ctx = (emm_context_t *) (args);

  if (!(emm_ctx)) {
    OAILOG_ERROR(LOG_NAS_EMM, "T3460 timer expired No EMM context\n");
    OAILOG_FUNC_OUT(LOG_NAS_EMM);
  }
  nas_emm_smc_proc_t *smc_proc = get_nas_common_procedure_smc(emm_ctx);

  if (smc_proc) {
    /*
     * Increment the retransmission counter
     */
    smc_proc->retransmission_count += 1;
    OAILOG_WARNING(
      LOG_NAS_EMM,
      "EMM-PROC  - T3460 timer expired, retransmission "
      "counter = %d\n",
      smc_proc->retransmission_count);
    if (SECURITY_COUNTER_MAX > smc_proc->retransmission_count) {
      REQUIREMENT_3GPP_24_301(R10_5_4_3_7_b__1);
      /*
       * Send security mode command message to the UE
       */
      _security_request(smc_proc);
    } else {
      REQUIREMENT_3GPP_24_301(R10_5_4_3_7_b__2);
      /*
     * Abort the security mode control and attach procedure
     */
      increment_counter(
        "nas_security_mode_command_timer_expired", 1, NO_LABELS);
      increment_counter(
        "ue_attach",
        1,
        2,
        "result",
        "failure",
        "cause",
        "no_response_for_security_mode_command");
      _security_abort(emm_ctx, (struct nas_base_proc_s *) smc_proc);
      emm_common_cleanup_by_ueid(smc_proc->ue_id);
      emm_sap_t emm_sap = {0};
      emm_sap.primitive = EMMCN_IMPLICIT_DETACH_UE;
      emm_sap.u.emm_cn.u.emm_cn_implicit_detach.ue_id = smc_proc->ue_id;
      emm_sap_send(&emm_sap);
      increment_counter("ue_attach", 1, 1, "action", "attach_abort");
    }
  }
  OAILOG_FUNC_OUT(LOG_NAS_EMM);
}

/*
   --------------------------------------------------------------------------
                MME specific local functions
   --------------------------------------------------------------------------
*/

/****************************************************************************
 **                                                                        **
 ** Name:    _security_request()                                       **
 **                                                                        **
 ** Description: Sends SECURITY MODE COMMAND message and start timer T3460 **
 **                                                                        **
 ** Inputs:  data:      Security mode control internal data        **
 **      is_new:    Indicates whether a new security context   **
 **             has just been taken into use               **
 **      Others:    None                                       **
 **                                                                        **
 ** Outputs:     None                                                      **
 **      Return:    RETURNok, RETURNerror                      **
 **      Others:    T3460                                      **
 **                                                                        **
 ***************************************************************************/
static int _security_request(nas_emm_smc_proc_t *const smc_proc)
{
  OAILOG_FUNC_IN(LOG_NAS_EMM);
  ue_mm_context_t *ue_mm_context = NULL;
  struct emm_context_s *emm_ctx = NULL;
  emm_sap_t emm_sap = {0};
  int rc = RETURNerror;

  if (smc_proc) {
    /*
     * Notify EMM-AS SAP that Security Mode Command message has to be sent
     * to the UE
     */
    REQUIREMENT_3GPP_24_301(R10_5_4_3_2__14);
    emm_sap.primitive = EMMAS_SECURITY_REQ;
    emm_sap.u.emm_as.u.security.puid =
      smc_proc->emm_com_proc.emm_proc.base_proc.nas_puid;
    emm_sap.u.emm_as.u.security.guti = NULL;
    emm_sap.u.emm_as.u.security.ue_id = smc_proc->ue_id;
    emm_sap.u.emm_as.u.security.msg_type = EMM_AS_MSG_TYPE_SMC;
    emm_sap.u.emm_as.u.security.ksi = smc_proc->ksi;
    emm_sap.u.emm_as.u.security.eea = smc_proc->eea;
    emm_sap.u.emm_as.u.security.eia = smc_proc->eia;
    emm_sap.u.emm_as.u.security.ucs2 = smc_proc->ucs2;
    emm_sap.u.emm_as.u.security.uea = smc_proc->uea;
    emm_sap.u.emm_as.u.security.uia = smc_proc->uia;
    emm_sap.u.emm_as.u.security.gea = smc_proc->gea;
    emm_sap.u.emm_as.u.security.umts_present = smc_proc->umts_present;
    emm_sap.u.emm_as.u.security.gprs_present = smc_proc->gprs_present;
    emm_sap.u.emm_as.u.security.selected_eea = smc_proc->selected_eea;
    emm_sap.u.emm_as.u.security.selected_eia = smc_proc->selected_eia;
    emm_sap.u.emm_as.u.security.imeisv_request = smc_proc->imeisv_request;

    mme_app_desc_t *mme_app_desc_p = get_mme_nas_state(false);
    ue_mm_context = mme_ue_context_exists_mme_ue_s1ap_id(
      &mme_app_desc_p->mme_ue_contexts, smc_proc->ue_id);
    if (ue_mm_context) {
      emm_ctx = &ue_mm_context->emm_context;
    } else {
      OAILOG_ERROR(
        LOG_NAS_EMM,
        "UE MM Context NULL! for ue_id = (%u)\n",
        smc_proc->ue_id);
      OAILOG_FUNC_RETURN(LOG_NAS_EMM, RETURNerror);
    }

    /*
     * Request for IMEISV from ue, if imeisv_request_enabled is enabled
     */
    emm_sap.u.emm_as.u.security.imeisv_request_enabled = EMM_IMEISV_REQUESTED;

    /*
   * Setup EPS NAS security data
   */
    emm_as_set_security_data(
      &emm_sap.u.emm_as.u.security.sctx,
      &emm_ctx->_security,
      smc_proc->is_new,
      false);
    rc = emm_sap_send(&emm_sap);

    if (rc != RETURNerror) {
      REQUIREMENT_3GPP_24_301(R10_5_4_3_2__1);
      void *timer_callback_args = NULL;
      nas_stop_T3460(smc_proc->ue_id, &smc_proc->T3460, timer_callback_args);
      /*
       * Start T3460 timer
       */
      nas_start_T3460(
        smc_proc->ue_id,
        &smc_proc->T3460,
        smc_proc->emm_com_proc.emm_proc.base_proc.time_out,
        emm_ctx);
    }
  }
  unlock_ue_contexts(ue_mm_context);
  OAILOG_FUNC_RETURN(LOG_NAS_EMM, rc);
}

//------------------------------------------------------------------------------
static int _security_ll_failure(
  emm_context_t *emm_context,
  struct nas_emm_proc_s *nas_emm_proc)
{
  OAILOG_FUNC_IN(LOG_NAS_EMM);
  int rc = RETURNerror;
  if (nas_emm_proc) {
    nas_emm_smc_proc_t *smc_proc = (nas_emm_smc_proc_t *) nas_emm_proc;
    REQUIREMENT_3GPP_24_301(R10_5_4_3_7_a);
    emm_sap_t emm_sap = {0};
    mme_ue_s1ap_id_t ue_id =
      PARENT_STRUCT(emm_context, struct ue_mm_context_s, emm_context)
        ->mme_ue_s1ap_id;

    emm_sap.primitive = EMMREG_COMMON_PROC_ABORT;
    emm_sap.u.emm_reg.ue_id = ue_id;
    emm_sap.u.emm_reg.ctx = emm_context;
    emm_sap.u.emm_reg.notify = true;
    emm_sap.u.emm_reg.free_proc = true;
    emm_sap.u.emm_reg.u.common.common_proc = &smc_proc->emm_com_proc;
    emm_sap.u.emm_reg.u.common.previous_emm_fsm_state =
      smc_proc->emm_com_proc.emm_proc.previous_emm_fsm_state;
    rc = emm_sap_send(&emm_sap);
  }
  OAILOG_FUNC_RETURN(LOG_NAS_EMM, rc);
}

//------------------------------------------------------------------------------
static int _security_non_delivered_ho(
  emm_context_t *emm_ctx,
  struct nas_emm_proc_s *nas_emm_proc)
{
  OAILOG_FUNC_IN(LOG_NAS_EMM);
  int rc = RETURNerror;

  if (emm_ctx && nas_emm_proc) {
    /************************README**********************************************
  ** NAS Non Delivery indication during HO handling will be added when HO is
  ** supported.
  ** In non hand-over case if MME receives NAS Non Delivery indication message
  ** that implies eNB and UE has lost radio connection. In this case aborting
  ** the SMC and Attach Procedure.
  *****************************************************************************
  REQUIREMENT_3GPP_24_301(R10_5_4_3_7_e);
  ****************************************************************************/
    /*
     * Abort the security mode control and attach procedure
     */
    nas_emm_smc_proc_t *smc_proc = (nas_emm_smc_proc_t *) nas_emm_proc;
    smc_proc->is_new = false;
    _security_abort(emm_ctx, (struct nas_base_proc_s *) smc_proc);
    emm_common_cleanup_by_ueid(smc_proc->ue_id);
    // Clean up MME APP UE context
    emm_sap_t emm_sap = {0};
    emm_sap.primitive = EMMCN_IMPLICIT_DETACH_UE;
    emm_sap.u.emm_cn.u.emm_cn_implicit_detach.ue_id = smc_proc->ue_id;
    emm_sap_send(&emm_sap);
    increment_counter("ue_attach", 1, 1, "action", "attach_abort");
  }
  OAILOG_FUNC_RETURN(LOG_NAS_EMM, rc);
}

/****************************************************************************
 **                                                                        **
 ** Name:    _security_abort()                                             **
 **                                                                        **
 ** Description: Aborts the security mode control procedure currently in   **
 **      progress                                                          **
 **                                                                        **
 ** Inputs:  args:      Security mode control data to be released          **
 **      Others:    None                                                   **
 **                                                                        **
 ** Outputs:     None                                                      **
 **      Return:    RETURNok, RETURNerror                                  **
 **                                                                        **
 ***************************************************************************/
static int _security_abort(
  emm_context_t *emm_ctx,
  struct nas_base_proc_s *base_proc)
{
  OAILOG_FUNC_IN(LOG_NAS_EMM);
  int rc = RETURNerror;
  unsigned int ue_id;

  if (emm_ctx && base_proc) {
    nas_emm_smc_proc_t *smc_proc = (nas_emm_smc_proc_t *) base_proc;
    ue_id = smc_proc->ue_id;
    OAILOG_WARNING(
      LOG_NAS_EMM,
      "EMM-PROC - Abort security mode control\
                    procedure "
      "(ue_id=" MME_UE_S1AP_ID_FMT ")\n",
      ue_id);
    /*
       * Stop timer T3460
       */
    if (smc_proc->T3460.id != NAS_TIMER_INACTIVE_ID) {
      OAILOG_INFO(
        LOG_NAS_EMM,
        "EMM-PROC  - Stop timer T3460 (%ld)\n",
        smc_proc->T3460.id);
      nas_stop_T3460(ue_id, &smc_proc->T3460, NULL);
    }
    /*
   * Release retransmission timer parameters
   * Do it after emm_sap_send
   */
    emm_proc_common_clear_args(ue_id);
  }
  OAILOG_FUNC_RETURN(LOG_NAS_EMM, rc);
}

/****************************************************************************
 **                                                                        **
 ** Name:    _security_select_algorithms()                                 **
 **                                                                        **
 ** Description: Select int and enc algorithms based on UE capabilities and**
 **      MME capabilities and MME preferences                              **
 **                                                                        **
 ** Inputs:  ue_eia:      integrity algorithms supported by UE             **
 **          ue_eea:      ciphering algorithms supported by UE             **
 **                                                                        **
 ** Outputs: mme_eia:     integrity algorithms supported by MME            **
 **          mme_eea:     ciphering algorithms supported by MME            **
 **                                                                        **
 **      Return:    RETURNok, RETURNerror                                  **
 **      Others:    None                                                   **
 **                                                                        **
 ***************************************************************************/
static int _security_select_algorithms(
  const int ue_eiaP,
  const int ue_eeaP,
  int *const mme_eiaP,
  int *const mme_eeaP)
{
  OAILOG_FUNC_IN(LOG_NAS_EMM);
  int preference_index;

  *mme_eiaP = NAS_SECURITY_ALGORITHMS_EIA0;
  *mme_eeaP = NAS_SECURITY_ALGORITHMS_EEA0;

  for (preference_index = 0; preference_index < 8; preference_index++) {
    if (
      ue_eiaP &
      (0x80 >> _emm_data.conf.prefered_integrity_algorithm[preference_index])) {
      OAILOG_DEBUG(
        LOG_NAS_EMM,
        "Selected  NAS_SECURITY_ALGORITHMS_EIA%d (choice num %d)\n",
        _emm_data.conf.prefered_integrity_algorithm[preference_index],
        preference_index);
      *mme_eiaP = _emm_data.conf.prefered_integrity_algorithm[preference_index];
      break;
    }
  }

  for (preference_index = 0; preference_index < 8; preference_index++) {
    if (
      ue_eeaP &
      (0x80 >> _emm_data.conf.prefered_ciphering_algorithm[preference_index])) {
      OAILOG_DEBUG(
        LOG_NAS_EMM,
        "Selected  NAS_SECURITY_ALGORITHMS_EEA%d (choice num %d)\n",
        _emm_data.conf.prefered_ciphering_algorithm[preference_index],
        preference_index);
      *mme_eeaP = _emm_data.conf.prefered_ciphering_algorithm[preference_index];
      break;
    }
  }

  OAILOG_FUNC_RETURN(LOG_NAS_EMM, RETURNok);
}

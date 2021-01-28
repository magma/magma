/**
 * Copyright 2020 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
/*****************************************************************************

  Source      amf_security_mode_control.cpp

  Version     0.1

  Date        2020/11/05

  Product     NAS stack

  Subsystem   Access and Mobility Management Function

  Author      Sandeep Kumar Mall

  Description Defines Access and Mobility Management Messages

*****************************************************************************/
#include <sstream>
#ifdef __cplusplus
extern "C" {
#endif
#include "log.h"
#include "secu_defs.h"
#ifdef __cplusplus
}
#endif
//#include "secu_defs.h"
#include "amf_asDefs.h"
#include "amf_fsm.h"
#include "amf_sap.h"
#include "amf_config.h"
#include "amf_app_ue_context_and_proc.h"
#include "nas5g_network.h"
using namespace std;

namespace magma5g {
extern ue_m5gmm_context_s
    ue_m5gmm_global_context;  // TODO AMF_TEST global var to temporarily store
                              // context inserted to ht
nas5g_config_t amf_data;
amf_sap_c amf_sap_sec;
nas_amf_smc_proc_t smc_ctrl;
amf_as_data_t amf_data_sec_ctrl;
nas_network nas_networks_smc;

//-----------------------------------------------------------------------------
nas_amf_smc_proc_t* nas5g_new_smc_procedure(amf_context_t* const amf_context) {
  if (!(amf_context->amf_procedures)) {
    amf_context->amf_procedures = _nas_new_amf_procedures(amf_context);
  }

  nas_amf_smc_proc_t* smc_proc = new (nas_amf_smc_proc_t);

  smc_proc->amf_com_proc.amf_proc.base_proc.nas_puid =
      __sync_fetch_and_add(&nas_puid, 1);
  smc_proc->amf_com_proc.amf_proc.base_proc.type = NAS_PROC_TYPE_AMF;
  smc_proc->amf_com_proc.amf_proc.type           = NAS_AMF_PROC_TYPE_COMMON;
  smc_proc->amf_com_proc.type                    = AMF_COMM_PROC_SMC;

  // smc_proc->T3460.sec = mme_config.nas_config.t3460_sec;
  // smc_proc->T3460.id  = NAS5G_TIMER_INACTIVE_ID;

  // nas_amf_common_procedure_t* wrapper = calloc(1, sizeof(*wrapper));
  nas_amf_common_procedure_t* wrapper = new nas_amf_common_procedure_t;
  if (wrapper) {
    wrapper->proc = &smc_proc->amf_com_proc;
    LIST_INSERT_HEAD(
        &amf_context->amf_procedures->amf_common_procs, wrapper, entries);
    OAILOG_TRACE(LOG_NAS_EMM, "New EMM_COMM_PROC_AUTH\n");
    return smc_proc;
  } else {
    nas_networks_smc.free_wrapper((void**) &smc_proc);
  }
  return NULL;
}

/****************************************************************************
 **                                                                        **
 ** Name:    amf_security_select_algorithms()                                 **
 **                                                                        **
 ** Description: Select int and enc algorithms based on UE capabilities and**
 **      5G MM capabilities and AMF preferences                              **
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
static int amf_security_select_algorithms(
    const int ue_eiaP, const int ue_eeaP, int* const amf_eiaP,
    int* const amf_eeaP) {
  OAILOG_FUNC_IN(LOG_NAS_AMF);
  int preference_index;

  *amf_eiaP = 0;
  *amf_eeaP = 0;

  for (preference_index = 0; preference_index < 8; preference_index++) {
    if (ue_eiaP &
        (0x80 >> amf_data.prefered_integrity_algorithm[preference_index])) {
      OAILOG_DEBUG(
          LOG_NAS_EMM,
          "Selected  NAS_SECURITY_ALGORITHMS_EIA%d (choice num %d)\n",
          amf_data.prefered_integrity_algorithm[preference_index],
          preference_index);
      *amf_eiaP = amf_data.prefered_integrity_algorithm[preference_index];
      break;
    }
  }

  for (preference_index = 0; preference_index < 8; preference_index++) {
    if (ue_eeaP &
        (0x80 >> amf_data.prefered_ciphering_algorithm[preference_index])) {
      OAILOG_DEBUG(
          LOG_NAS_AMF,
          "Selected  NAS_SECURITY_ALGORITHMS_EEA%d (choice num %d)\n",
          amf_data.prefered_ciphering_algorithm[preference_index],
          preference_index);
      *amf_eeaP = amf_data.prefered_ciphering_algorithm[preference_index];
      break;
    }
  }

  OAILOG_FUNC_RETURN(LOG_NAS_EMM, RETURNok);
}

//------------------------------------------------------------------------------
static int amf_security_ll_failure(
    amf_context_t* amf_context, nas_amf_proc_t* nas_amf_proc) {
  OAILOG_FUNC_IN(LOG_NAS_AMF);
  int rc = RETURNerror;
  if (nas_amf_proc) {
    // nas_amf_smc_proc_t* smc_proc = (nas_amf_smc_proc_t*)
    // nas_amf_proc;//TODO-RECHECK

    amf_sap_t amf_sap_seq;
    amf_ue_ngap_id_t ue_id =
        PARENT_STRUCT(amf_context, ue_m5gmm_context_s, amf_context)
            ->amf_ue_ngap_id;

    amf_sap_seq.primitive           = AMFREG_COMMON_PROC_ABORT;
    amf_sap_seq.u.amf_reg.ue_id     = ue_id;
    amf_sap_seq.u.amf_reg.ctx       = amf_context;
    amf_sap_seq.u.amf_reg.notify    = true;
    amf_sap_seq.u.amf_reg.free_proc = true;
    // amf_sap_seq.u.amf_reg.u.common.common_proc = &smc_proc->amf_com_proc;
    // amf_sap_seq.u.amf_reg.u.common.previous_amf_fsm_state =
    //    smc_proc->amf_com_proc.amf_proc.previous_amf_fsm_state;
    rc = amf_sap_sec.amf_sap_send(&amf_sap_seq);
  }
  OAILOG_FUNC_RETURN(LOG_NAS_AMF, rc);
}

/****************************************************************************
 **                                                                        **
 ** Name:    amf_security_request()                                       **
 **                                                                        **
 ** Description: Sends SECURITY MODE COMMAND message and start timer T3560 **
 **                                                                        **
 ** Inputs:  data:      Security mode control internal data        **
 **      is_new:    Indicates whether a new security context   **
 **             has just been taken into use               **
 **      Others:    None                                       **
 **                                                                        **
 ** Outputs:     None                                                      **
 **      Return:    RETURNok, RETURNerror                      **
 **      Others:    T3560                                      **
 **                                                                        **
 ***************************************************************************/
#if 1
static int amf_security_request(nas_amf_smc_proc_t* const smc_proc) {
  OAILOG_FUNC_IN(LOG_NAS_AMF);
  ue_m5gmm_context_s* ue_mm_context = NULL;
  amf_context_t* amf_ctx            = NULL;
  amf_sap_t amf_sap;
  int rc = RETURNerror;

  if (smc_proc) {
    /*
     * Notify AMF-AS SAP that Security Mode Command message has to be sent
     * to the UE
     */

    amf_sap.primitive = AMFAS_SECURITY_REQ;
    amf_sap.u.amf_as.u.security.puid =
        smc_proc->amf_com_proc.amf_proc.base_proc.nas_puid;
    amf_sap.u.amf_as.u.security.guti     = {0};
    amf_sap.u.amf_as.u.security.ue_id    = smc_proc->ue_id;
    amf_sap.u.amf_as.u.security.msg_type = AMF_AS_MSG_TYPE_SMC;
    amf_sap.u.amf_as.u.security.ksi      = smc_proc->ksi;
    amf_sap.u.amf_as.u.security.eea      = smc_proc->eea;
    amf_sap.u.amf_as.u.security.eia      = smc_proc->eia;
    amf_sap.u.amf_as.u.security.ucs2     = smc_proc->ucs2;
    // amf_sap.u.amf_as.u.security.uea            = smc_proc->uea;
    amf_sap.u.amf_as.u.security.selected_eea   = smc_proc->selected_eea;
    amf_sap.u.amf_as.u.security.selected_eia   = smc_proc->selected_eia;
    amf_sap.u.amf_as.u.security.imeisv_request = smc_proc->imeisv_request;

    //    ue_mm_context = amf_ue_context_exists_amf_ue_ngap_id(smc_proc->ue_id);
    ue_mm_context =
        &ue_m5gmm_global_context;  // TODO AMF_TEST global var to temporarily
                                   // store context inserted to ht
    if (ue_mm_context) {
      amf_ctx = &ue_mm_context->amf_context;
    } else {
      OAILOG_ERROR(
          LOG_NAS_AMF, "UE 5G-MM Context NULL! for ue_id = (%u)\n",
          smc_proc->ue_id);
      OAILOG_FUNC_RETURN(LOG_NAS_AMF, RETURNerror);
    }

    /*
     * Request for IMEISV from ue, if imeisv_request_enabled is enabled
     */
    amf_sap.u.amf_as.u.security.imeisv_request_enabled = AMF_IMEISV_REQUESTED;

    /*
     * Setup 5GCN NAS security data
     */
    amf_data_sec_ctrl.amf_as_set_security_data(
        &amf_sap.u.amf_as.u.security.sctx, &amf_ctx->_security,
        smc_proc->is_new, false);
    rc = amf_sap_sec.amf_sap_send(&amf_sap);
#if 0
    if (rc != RETURNerror) {
      void* timer_callback_args = NULL;
      nas_stop_T3560(smc_proc->ue_id, &smc_proc->T3560, timer_callback_args);
      /*
       * Start T3560 timer
       */
      nas_start_T3560(
          smc_proc->ue_id, &smc_proc->T3560,
          smc_proc->amf_com_proc.amf_proc.base_proc.time_out, amf_ctx);
    }
#endif
  }
  OAILOG_FUNC_RETURN(LOG_NAS_AMF, rc);
}
#endif

/*
--------------------------------------------------------------------------
        Security mode control procedure executed by the AMF
--------------------------------------------------------------------------
*/

/****************************************************************************
 **                                                                        **
 ** Name:    amf_proc_security_mode_control()                          **
 **                                                                        **
 ** Description: Initiates the security mode control procedure.            **
 **                                                                        **
 **              3GPP TS 24.501, section 8.2.25                          **
 **      The AMF initiates the NAS security mode control procedure **
 **      by sending a SECURITY MODE COMMAND message to the UE and  **
 **      starting timer T3560. The message shall be sent unciphe-  **
 **      red but shall be integrity protected using the NAS inte-  **
 **      grity key based on KASME.                                 **
 **                                                                        **
 ** Inputs:  ue_id:      UE lower layer identifier                  **
 **      ksi:       NAS key set identifier                     **
 **      eea:       Replayed 5G encryption algorithms         **
 **      eia:       Replayed 5G integrity algorithms          **
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
int amf_proc_security_mode_control(
    amf_context_t* amf_ctx, nas_amf_specific_proc_t* amf_specific_proc,
    ksi_t ksi, success_cb_t success, failure_cb_t failure) {
  OAILOG_FUNC_IN(LOG_NAS_AMF);
  int rc                       = RETURNerror;
  bool security_context_is_new = false;
  int amf_eea                  = 0;  // TODO
  int amf_eia                  = 0;  // TODO
  /*
   * Get the UE context
   */

  OAILOG_INFO(
      LOG_NAS_AMF,
      "AMF_TEST: Initiating security mode control procedure, "
      "KSI = %d\n",
      ksi);

  if (!(amf_ctx)) {
    OAILOG_ERROR(LOG_NAS_AMF, "Amf Context NULL!\n");
    OAILOG_FUNC_RETURN(LOG_NAS_AMF, RETURNerror);
  }

  // TODO better than that (quick fixes)
  if (KSI_NO_KEY_AVAILABLE == ksi) {
    ksi = 0;
  }
  if (AMF_SECURITY_VECTOR_INDEX_INVALID == amf_ctx->_security.vector_index) {
    // amf_ctx_set_security_vector_index(amf_ctx, 0);  // sandeep
    amf_ctx->_security.vector_index = 0;
  }

  amf_ue_ngap_id_t ue_id =
      PARENT_STRUCT(amf_ctx, ue_m5gmm_context_s, amf_context)->amf_ue_ngap_id;
  nas_amf_smc_proc_t* smc_proc = get_nas5g_common_procedure_smc(amf_ctx);
  if (!smc_proc) {
    smc_proc = nas5g_new_smc_procedure(amf_ctx);
    if (smc_proc) {
      smc_proc->saved_selected_eea =
          amf_ctx->_security.selected_algorithms.encryption;
      smc_proc->saved_selected_eia =
          amf_ctx->_security.selected_algorithms.integrity;
      smc_proc->saved_eksi     = amf_ctx->_security.eksi;
      smc_proc->saved_overflow = amf_ctx->_security.dl_count.overflow;
      smc_proc->saved_seq_num  = amf_ctx->_security.dl_count.seq_num;
      smc_proc->saved_sc_type  = amf_ctx->_security.sc_type;
      /*
       * The security mode control procedure is initiated to take into use
       * * * * the 5GS  security context created after a successful execution of
       * * * * the 5GS authentication procedure
       */

      smc_ctrl.amf_ctx_set_security_eksi(amf_ctx, ksi);
      amf_ctx->_security.dl_count.overflow = 0;
      amf_ctx->_security.dl_count.seq_num  = 0;

      /*
       *  Compute NAS cyphering and integrity keys
       */
      // need to wright m5g_ue_network_capability class in 24.501 Sandeep
      /* rc = amf_security_select_algorithms(
           amf_ctx->m5g_ue_network_capability.eia,
           amf_ctx->m5g_ue_network_capability.eea, &amf_eia, &amf_eea); */
      // need to wright m5g_ue_network_capability class in 24.501 Sandeep
      rc                                                = RETURNok;  // AMF_TEST
      amf_ctx->_security.selected_algorithms.encryption = amf_eea;
      amf_ctx->_security.selected_algorithms.integrity  = amf_eia;
      if (rc == RETURNerror) {
        // OAILOG_WARNING(
        //    LOG_NAS_AMF, "AMF-PROC  - Failed to select security
        //    algorithms\n");
        OAILOG_FUNC_RETURN(LOG_NAS_AMF, RETURNerror);
      }

      smc_ctrl.amf_ctx_set_security_type(
          amf_ctx, SECURITY_CTX_TYPE_FULL_NATIVE);
      // TODO AssertFatal(KSI_NO_KEY_AVAILABLE > amf_ctx->_security.eksi, "eksi
      // not valid");
      derive_key_nas(
          NAS_INT_ALG, amf_ctx->_security.selected_algorithms.integrity,
          amf_ctx->_vector[amf_ctx->_security.eksi % MAX_EPS_AUTH_VECTORS]
              .kasme,
          amf_ctx->_security.knas_int);
      derive_key_nas(
          NAS_ENC_ALG, amf_ctx->_security.selected_algorithms.encryption,
          amf_ctx->_vector[amf_ctx->_security.eksi % MAX_EPS_AUTH_VECTORS]
              .kasme,
          amf_ctx->_security.knas_enc);
      /*
       * Set new security context indicator
       */
      security_context_is_new = true;
      amf_ctx_set_attribute_present(amf_ctx, AMF_CTXT_MEMBER_SECURITY);
    }
  } else {
    OAILOG_ERROR(LOG_NAS_AMF, "AMF-PROC  - No 5G CN security context exists\n");
    OAILOG_FUNC_RETURN(LOG_NAS_AMF, RETURNerror);
  }

  if (smc_proc) {
    /*
     * Setup ongoing AMF procedure callback functions
     */
    //((nas5g_base_proc_t*) smc_proc)->parent =    (nas_base_proc_t*)
    // amf_specific_proc;
    smc_proc->amf_com_proc.amf_proc.delivered = NULL;
    // smc_proc->amf_com_proc.amf_proc.previous_amf_fsm_state =
    //    amf_fsm_get_state(amf_ctx);
    // smc_proc->amf_com_proc.amf_proc.not_delivered = amf_security_ll_failure;
    // smc_proc->amf_com_proc.amf_proc.not_delivered_ho =
    // _security_non_delivered_ho;
    smc_proc->amf_com_proc.amf_proc.base_proc.success_notif = success;
    smc_proc->amf_com_proc.amf_proc.base_proc.failure_notif = failure;
    // smc_proc->amf_com_proc.amf_proc.base_proc.abort         =
    // _security_abort;
    smc_proc->amf_com_proc.amf_proc.base_proc.fail_in  = NULL;  // only response
    smc_proc->amf_com_proc.amf_proc.base_proc.fail_out = NULL;
    // smc_proc->amf_com_proc.amf_proc.base_proc.time_out =
    // _security_t3560_handler;

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
    // smc_proc->eea = amf_ctx->m5g_ue_network_capability.eea;
    /*
     * Set the EPS integrity algorithms to be replayed to the UE
     */
    // smc_proc->eia  = amf_ctx->m5g_ue_network_capability.eia;
    // smc_proc->ucs2 = amf_ctx->m5g_ue_network_capability.ucs2;

    /*
     * Set the M5G encryption algorithms selected to the UE
     */
    smc_proc->selected_eea = amf_ctx->_security.selected_algorithms.encryption;
    OAILOG_DEBUG(
        LOG_NAS_AMF,
        "5G CN encryption algorithm selected is (%d) for ue_id (%u)\n",
        smc_proc->selected_eea, ue_id);
    /*
     * Set the 5G CN integrity algorithms selected to the UE
     */
    smc_proc->selected_eia = amf_ctx->_security.selected_algorithms.integrity;
    OAILOG_DEBUG(
        LOG_NAS_AMF,
        "5G CN integrity algorithm selected is (%d) for ue_id (%u)\n",
        smc_proc->selected_eia, ue_id);

    smc_proc->is_new = security_context_is_new;

    // always ask for IMEISV (Do it simple now)
    smc_proc->imeisv_request = true;
    // smc_proc->imeisv_request = (IS_AMF_CTXT_PRESENT_IMEISV(amf_ctx)) ?
    // false:true;

    /*
     * Send security mode command message to the UE
     */
    rc = amf_security_request(smc_proc);
  }

  OAILOG_FUNC_RETURN(LOG_NAS_AMF, rc);
}

}  // namespace magma5g

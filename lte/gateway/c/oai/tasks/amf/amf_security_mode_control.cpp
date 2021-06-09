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

#include <sstream>
#ifdef __cplusplus
extern "C" {
#endif
#include "log.h"
#include "secu_defs.h"
#include "intertask_interface_types.h"
#include "intertask_interface.h"
#include "dynamic_memory_check.h"
#ifdef __cplusplus
}
#endif
#include "common_defs.h"
#include "amf_asDefs.h"
#include "amf_data.h"
#include "amf_fsm.h"
#include "amf_sap.h"
#include "amf_config.h"
#include "amf_app_ue_context_and_proc.h"
#include "amf_identity.h"
#include "conversions.h"

#define IMSI64_TO_IMSI15(iMsI64_t, imsi15)                                     \
  {                                                                            \
    if ((iMsI64_t / 100000000000000) != 0) {                                   \
      imsi15[0]  = iMsI64_t / 100000000000000;                                 \
      iMsI64_t   = iMsI64_t % 100000000000000;                                 \
      imsi15[1]  = iMsI64_t / 10000000000000;                                  \
      iMsI64_t   = iMsI64_t % 10000000000000;                                  \
      imsi15[2]  = iMsI64_t / 1000000000000;                                   \
      iMsI64_t   = iMsI64_t % 1000000000000;                                   \
      imsi15[3]  = iMsI64_t / 100000000000;                                    \
      iMsI64_t   = iMsI64_t % 100000000000;                                    \
      imsi15[4]  = iMsI64_t / 10000000000;                                     \
      iMsI64_t   = iMsI64_t % 10000000000;                                     \
      imsi15[5]  = iMsI64_t / 1000000000;                                      \
      iMsI64_t   = iMsI64_t % 1000000000;                                      \
      imsi15[6]  = iMsI64_t / 100000000;                                       \
      iMsI64_t   = iMsI64_t % 100000000;                                       \
      imsi15[7]  = iMsI64_t / 10000000;                                        \
      iMsI64_t   = iMsI64_t % 10000000;                                        \
      imsi15[8]  = iMsI64_t / 1000000;                                         \
      iMsI64_t   = iMsI64_t % 1000000;                                         \
      imsi15[9]  = iMsI64_t / 100000;                                          \
      iMsI64_t   = iMsI64_t % 100000;                                          \
      imsi15[10] = iMsI64_t / 10000;                                           \
      iMsI64_t   = iMsI64_t % 10000;                                           \
      imsi15[11] = iMsI64_t / 1000;                                            \
      iMsI64_t   = iMsI64_t % 1000;                                            \
      imsi15[12] = iMsI64_t / 100;                                             \
      iMsI64_t   = iMsI64_t % 100;                                             \
      imsi15[13] = iMsI64_t / 10;                                              \
      iMsI64_t   = iMsI64_t % 10;                                              \
      imsi15[14] = iMsI64_t / 1;                                               \
    }                                                                          \
  }

namespace magma5g {

extern task_zmq_ctx_s amf_app_task_zmq_ctx;
nas5g_config_t amf_data;
nas_amf_smc_proc_t smc_ctrl;
amf_as_data_t amf_data_sec_ctrl;

//-----------------------------------------------------------------------------

void format_plmn(amf_plmn_t* plmn) {
  int loop       = 0;
  uint8_t* octet = (uint8_t*) plmn;
  /*TODO handle this better; for 2 digit mnc, the mnc_digit3 will be coming in
   * as 0xf. This has to be changed to 0x0 before being used to create SNNI.
   * When the PLMN value is used to form SNNI, the value is shifted such that
   * the mnc_digit3, which was made 0, has to become mncdigit_1. For example a 5
   * digit plmn such as 20895 will be 20895f in NAS. This will be expanded to
   * 208 095 in respective mcc mnc values
   */
  bool format_flag = false;
  for (loop = 0; loop < 3; loop++) {
    uint8_t d2 = octet[loop];
    uint8_t d1 = (d2 & 0xf0) >> 4;
    d2         = d2 & 0x0f;
    if (d2 >= 10) {
      octet[loop] = octet[loop] & 0xf0;
      format_flag = true;
    }
    if (d1 >= 10) {
      octet[loop] = octet[loop] & 0x0f;
      format_flag = true;
    }
  }
  if (format_flag) {
    amf_plmn_t temp_plmn;
    memcpy(&temp_plmn, plmn, 3);
    plmn->mnc_digit1 = 0;
    plmn->mnc_digit2 = temp_plmn.mnc_digit1;
    plmn->mnc_digit3 = temp_plmn.mnc_digit2;
  }
}

static int security_mode_t3560_handler(zloop_t* loop, int timer_id, void* arg);

/****************************************************************************
 **                                                                        **
 ** Name:    nas5g_new_smc_procedure()                                     **
 **                                                                        **
 ** Description: NAS5g smc procedure Creation                              **
 **                                                                        **
 **                                                                        **
 ***************************************************************************/
nas_amf_smc_proc_t* nas5g_new_smc_procedure(amf_context_t* const amf_context) {
  if (!(amf_context->amf_procedures)) {
    amf_context->amf_procedures = nas_new_amf_procedures(amf_context);
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
    free_wrapper((void**) &smc_proc);
  }
  return NULL;
}

/****************************************************************************
 **                                                                        **
 ** Name:    amf_security_request()                                        **
 **                                                                        **
 ** Description: Sends SECURITY MODE COMMAND message                       **
 **                                                                        **
 **                                                                        **
 ***************************************************************************/
static int amf_security_request(nas_amf_smc_proc_t* const smc_proc) {
  OAILOG_FUNC_IN(LOG_NAS_AMF);
  ue_m5gmm_context_s* ue_mm_context = NULL;
  amf_context_t* amf_ctx            = NULL;
  amf_sap_t amf_sap;
  int rc             = RETURNerror;
  smc_proc->T3560.id = NAS5G_TIMER_INACTIVE_ID;

  if (smc_proc) {
    /*
     * Notify AMF-AS SAP that Security Mode Command message has to be sent
     * to the UE
     */
    amf_sap.primitive = AMFAS_SECURITY_REQ;
    amf_sap.u.amf_as.u.security.puid =
        smc_proc->amf_com_proc.amf_proc.base_proc.nas_puid;
    amf_sap.u.amf_as.u.security.guti           = {0};
    amf_sap.u.amf_as.u.security.ue_id          = smc_proc->ue_id;
    amf_sap.u.amf_as.u.security.msg_type       = AMF_AS_MSG_TYPE_SMC;
    amf_sap.u.amf_as.u.security.ksi            = smc_proc->ksi;
    amf_sap.u.amf_as.u.security.eea            = smc_proc->eea;
    amf_sap.u.amf_as.u.security.eia            = smc_proc->eia;
    amf_sap.u.amf_as.u.security.ucs2           = smc_proc->ucs2;
    amf_sap.u.amf_as.u.security.selected_eea   = smc_proc->selected_eea;
    amf_sap.u.amf_as.u.security.selected_eia   = smc_proc->selected_eia;
    amf_sap.u.amf_as.u.security.imeisv_request = smc_proc->imeisv_request;
    ue_mm_context = amf_ue_context_exists_amf_ue_ngap_id(smc_proc->ue_id);

    if (ue_mm_context) {
      amf_ctx = &ue_mm_context->amf_context;
    } else {
      OAILOG_ERROR(
          LOG_NAS_AMF, "UE 5G-MM context NULL! for ue_id = (%u)\n",
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
    rc = amf_sap_send(&amf_sap);
    if (rc != RETURNerror) {
      OAILOG_INFO(
          LOG_AMF_APP, "Timer: Security Mode Calling start_timer_T3560 \n");
      smc_proc->T3560.id = start_timer(
          &amf_app_task_zmq_ctx, SECURITY_MODE_TIMER_EXPIRY_MSECS,
          TIMER_REPEAT_ONCE, security_mode_t3560_handler,
          (void*) smc_proc->ue_id);
      OAILOG_INFO(
          LOG_AMF_APP,
          "Timer:  After starting SECURITY_MODE_TIMER timer T3560 "
          "with id %d\n",
          smc_proc->T3560.id);
    }
  }
  OAILOG_FUNC_RETURN(LOG_NAS_AMF, rc);
}

/* Timer Expiry Handler for SECURITY COMMAND MODE Timer 3560 */
static int security_mode_t3560_handler(zloop_t* loop, int timer_id, void* arg) {
  OAILOG_INFO(LOG_AMF_APP, "Timer: In security_mode_t3560_handler - T3560\n");
  amf_context_t* amf_ctx = NULL;
  amf_ue_ngap_id_t ue_id = 0;
  ue_id                  = *((amf_ue_ngap_id_t*) (arg));

  ue_m5gmm_context_s* ue_context = amf_ue_context_exists_amf_ue_ngap_id(ue_id);

  if (ue_context == NULL) {
    OAILOG_INFO(LOG_AMF_APP, "AMF_TEST: ue_context is NULL\n");
    return -1;
  }

  OAILOG_INFO(
      LOG_AMF_APP,
      "Timer: Created ue_mm_context from global context - T3560\n");
  amf_ctx = &ue_context->amf_context;
  OAILOG_INFO(LOG_AMF_APP, "Timer: got amf ctx and calling common procedure\n");
  if (!(amf_ctx)) {
    OAILOG_ERROR(LOG_AMF_APP, "T3560 timer expired No AMF context\n");
    return 1;
    // OAILOG_FUNC_OUT(LOG_AMF_APP);
  }

  nas_amf_smc_proc_t* smc_proc = get_nas5g_common_procedure_smc(amf_ctx);

  OAILOG_ERROR(LOG_AMF_APP, "Timer:In Identity Expiration Handler ZMQ TIMER\n");

  if (smc_proc) {
    OAILOG_WARNING(
        LOG_AMF_APP, "T3560 timer   timer id %d ue id %d\n", smc_proc->T3560.id,
        smc_proc->ue_id);
    smc_proc->T3560.id = -1;
  }

  /*
   * Increment the retransmission counter
   */
  smc_proc->retransmission_count += 1;
  OAILOG_ERROR(
      LOG_AMF_APP, "Timer: Incrementing retransmission_count to %d\n",
      smc_proc->retransmission_count);

  if (smc_proc->retransmission_count < SECURITY_COUNTER_MAX) {
    /*
     * Send identity request message to the UE
     */
    OAILOG_ERROR(
        LOG_AMF_APP,
        "Timer: timer has expired Sending Security Command Mode request "
        "again\n");
    amf_security_request(smc_proc);
  } else {
    /*
     * Abort the smc procedure
     */
    OAILOG_ERROR(
        LOG_AMF_APP,
        "Timer: Maximum retires done hence Abort the smc procedure\n");
    return -1;
  }

  return 0;
}

/*
--------------------------------------------------------------------------
        Security mode control procedure executed by the AMF
--------------------------------------------------------------------------
*/

/****************************************************************************
 **                                                                        **
 ** Name:    amf_proc_security_mode_control()                              **
 **                                                                        **
 ** Description: Initiates the security mode control procedure.            **
 **                                                                        **
 **              3GPP TS 24.501, section 8.2.25                            **
 **      The AMF initiates the NAS security mode control procedure         **
 **      by sending a SECURITY MODE COMMAND message to the UE and          **
 **      starting timer T3560. The message shall be sent unciphe-          **
 **      red but shall be integrity protected using the NAS inte-          **
 **      grity key based on KASME.                                         **
 **                                                                        **
 ** Inputs:  amf_ctx:      amf context received                            **
 **      ksi:       NAS key set identifier                                 **
 **      amf_specific_proc:      AMF specific procedure                    **
 **      success:   Callback function executed when the secu-              **
 **             rity mode control procedure successfully                   **
 **             completes                                                  **
 **      failure:   Callback function executed whener a lower              **
 **             layer failure occured before the security                  **
 **             mode control procedure completes                           **
 **      Others:    None                                                   **
 **                                                                        **
 ** Outputs:     None                                                      **
 **      Return:    RETURNok, RETURNerror                                  **
 **      Others:    None                                                   **
 **                                                                        **
 ***************************************************************************/
int amf_proc_security_mode_control(
    amf_context_t* amf_ctx, nas_amf_specific_proc_t* amf_specific_proc,
    ksi_t ksi, success_cb_t success, failure_cb_t failure) {
  OAILOG_FUNC_IN(LOG_NAS_AMF);
  int rc                       = RETURNerror;
  bool security_context_is_new = false;
  // TODO: Hardcoded values Will be taken care in upcoming PR
  int amf_eea = 0;
  int amf_eia = 0;  // Integrity Algorithm 2
  amf_plmn_t plmn;
  uint8_t snni[32]  = {0};
  uint8_t kausf[32] = {0};
  uint8_t kseaf[32] = {0};
  uint8_t ck_ik[32] = {0};
  uint8_t ak_sqn[6] = {0};
  /*
   * Get the UE context
   */
  OAILOG_DEBUG(
      LOG_NAS_AMF,
      "Initiating security mode control procedure, "
      "KSI = %d\n",
      ksi);
  if (!(amf_ctx)) {
    OAILOG_ERROR(LOG_NAS_AMF, "Amf Context NULL!\n");
    OAILOG_FUNC_RETURN(LOG_NAS_AMF, RETURNerror);
  }

  // If ksi not available passing NULL key set identifier
  if (KSI_NO_KEY_AVAILABLE == ksi) {
    ksi = 0;
  }
  if (AMF_SECURITY_VECTOR_INDEX_INVALID == amf_ctx->_security.vector_index) {
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
       * the 5GS  security context created after a successful execution of
       * the 5GS authentication procedure
       */
      smc_ctrl.amf_ctx_set_security_eksi(amf_ctx, ksi);
      amf_ctx->_security.dl_count.overflow = 0;
      amf_ctx->_security.dl_count.seq_num  = 0;

      rc                                                = RETURNok;
      amf_ctx->_security.selected_algorithms.encryption = amf_eea;
      amf_ctx->_security.selected_algorithms.integrity  = amf_eia;
      if (rc == RETURNerror) {
        OAILOG_FUNC_RETURN(LOG_NAS_AMF, RETURNerror);
      }
      smc_ctrl.amf_ctx_set_security_type(
          amf_ctx, SECURITY_CTX_TYPE_FULL_NATIVE);

      // NAS Integrity key is calculated as specified in TS 33501, Annex A
      uint8_t imsi[15];
      IMSI64_TO_IMSI15(amf_ctx->imsi64, imsi);
      memcpy(&plmn, amf_ctx->imsi.u.value, 3);
      format_plmn(&plmn);

      /* Building 32 bytes of string with serving network SN
       * SN value 5G:mnc095.mcc208.3gppnetwork.org
       * mcc and mnc retrive saved _imsi from amf_context
       */
      uint32_t mcc              = 0;
      uint32_t mnc              = 0;
      uint32_t mnc_digit_length = 0;

      PLMN_T_TO_MCC_MNC(plmn, mcc, mnc, mnc_digit_length);
      uint32_t snni_buf_len =
          sprintf((char*) snni, "5G:mnc%03d.mcc%03d.3gppnetwork.org", mnc, mcc);
      if (snni_buf_len != 32) {
        OAILOG_ERROR(LOG_NAS_AMF, "Failed to create SNNI String\n");
        OAILOG_FUNC_RETURN(LOG_NAS_AMF, RETURNerror);
      }

      memcpy(
          ak_sqn,
          amf_ctx->_vector[amf_ctx->_security.eksi % MAX_EPS_AUTH_VECTORS].autn,
          6);
      memcpy(
          ck_ik,
          amf_ctx->_vector[amf_ctx->_security.eksi % MAX_EPS_AUTH_VECTORS].ck,
          16);
      memcpy(
          &ck_ik[16],
          amf_ctx->_vector[amf_ctx->_security.eksi % MAX_EPS_AUTH_VECTORS].ik,
          16);

      derive_5gkey_ausf(ck_ik, snni, ak_sqn, kausf);
      derive_5gkey_seaf(kausf, snni, kseaf);
      derive_5gkey_amf(imsi, 15, kseaf, amf_ctx->_security.kamf);

      derive_5gkey_nas(
          NAS_INT_ALG, 2, amf_ctx->_security.kamf, amf_ctx->_security.knas_int);

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
    // Setup ongoing AMF procedure callback functions
    smc_proc->amf_com_proc.amf_proc.delivered               = NULL;
    smc_proc->amf_com_proc.amf_proc.base_proc.success_notif = success;
    smc_proc->amf_com_proc.amf_proc.base_proc.failure_notif = failure;
    smc_proc->amf_com_proc.amf_proc.base_proc.fail_in       = NULL;
    smc_proc->amf_com_proc.amf_proc.base_proc.fail_out      = NULL;
    smc_proc->ue_id                                         = ue_id;
    smc_proc->retransmission_count                          = 0;
    smc_proc->ksi                                           = ksi;
    smc_proc->selected_eea = amf_ctx->_security.selected_algorithms.encryption;
    OAILOG_DEBUG(
        LOG_NAS_AMF,
        "5G CN encryption algorithm selected is (%d) for ue_id (%u)\n",
        smc_proc->selected_eea, ue_id);
    smc_proc->selected_eia = amf_ctx->_security.selected_algorithms.integrity;
    OAILOG_DEBUG(
        LOG_NAS_AMF,
        "5G CN integrity algorithm selected is (%d) for ue_id (%u)\n",
        smc_proc->selected_eia, ue_id);
    smc_proc->is_new = security_context_is_new;

    // always ask for IMEISV --TODO will be taken care in upcoming PRs
    smc_proc->imeisv_request = true;

    // Send security mode command message to the UE
    rc = amf_security_request(smc_proc);
  }
  OAILOG_FUNC_RETURN(LOG_NAS_AMF, rc);
}

}  // namespace magma5g

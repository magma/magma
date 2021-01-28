/**
 *copyright 2020 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 **/
/*****************************************************************************

  Source      amf_indentity.cpp

  Version     0.1

  Date        2020/07/28

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
#include "assertions.h"
#include "conversions.h"
#include "amf_config.h"
#ifdef __cplusplus
}
#endif
#include <unordered_map>
#include "amf_fsm.h"
#include "amf_identity.h"
#include "amf_asDefs.h"
#include "amf_sap.h"
#include "M5GSMobileIdentity.h"
#include "amf_app_ue_context_and_proc.h"
#include "amf_recv.h"
using namespace std;
extern amf_config_t amf_config;

namespace magma5g {
extern ue_m5gmm_context_s
    ue_m5gmm_global_context;  // TODO AMF-TEST global var to temporarily store
                              // context inserted to ht
// extern amf_config_t amf_config_handler;
// Global map of supi to guti along with amf_ue_ngap_id
std::unordered_map<imsi64_t, guti_and_amf_id_t> amf_supi_guti_map;
// std::unordered_map<supi_as_imsi_t, guti_and_amf_id_t> amf_supi_guti_map;

nas_proc nas_proc_indt;
/****************************************************************************
**                                                                        **
** Name:    amf_cn_identity_res()                                         **
**                                                                        **
** Description: Processes Identity Response message                       **
**                                                                        **
**      Inputs:  ue_id:      UE lower layer identifier                    **
**      msg:       The received EMM message                               **
**      Others:    None                                                   **
**                                                                        **
** Outputs:     amf_cause: AMF cause code                                 **
** Return:      RETURNok, RETURNerror                                     **
** Others:      None                                                      **
**                                                                        **
***************************************************************************/
int amf_identity_msg::amf_cn_identity_res(
    amf_ue_ngap_id_t ue_id, M5GSMobileIdentityMsg* msg, int* amf_cause,
    const amf_nas_message_decode_status_t* status) {
  OAILOG_FUNC_IN(LOG_NAS_AMF);
  int rc = RETURNok;

  OAILOG_INFO(LOG_NAS_AMF, "AMFAS-SAP - Received Identity Response message\n");
  /*
   * Message processing
   */
  /*
   * Get the mobile identity
   */
  imsi_t imsi = {0}, *p_imsi = NULL;
  imei_t imei = {0}, *p_imei = NULL;
  imeisv_t imeisv = {0}, *p_imeisv = NULL;
  tmsi_t tmsi = 0, *p_tmsi = NULL;

  if (msg->mobile_identity.imsi.type_of_identity == MOBILE_IDENTITY_IMSI) {
    /*
     * Get the IMSI
     */
    p_imsi             = &imsi;
    imsi.u.num.digit1  = msg->mobile_identity.imsi.mcc_digit1;
    imsi.u.num.digit2  = msg->mobile_identity.imsi.mcc_digit2;
    imsi.u.num.digit3  = msg->mobile_identity.imsi.mcc_digit3;
    imsi.u.num.digit4  = msg->mobile_identity.imsi.mnc_digit1;
    imsi.u.num.digit5  = msg->mobile_identity.imsi.mnc_digit2;
    imsi.u.num.digit6  = msg->mobile_identity.imsi.mnc_digit3;
    imsi.u.num.digit7  = msg->mobile_identity.imsi.rout_ind_digit_1;
    imsi.u.num.digit8  = msg->mobile_identity.imsi.rout_ind_digit_2;
    imsi.u.num.digit9  = msg->mobile_identity.imsi.rout_ind_digit_3;
    imsi.u.num.digit10 = msg->mobile_identity.imsi.rout_ind_digit_4;
    imsi.u.num.digit11 = msg->mobile_identity.imsi.msin_digit1;
    imsi.u.num.digit12 = msg->mobile_identity.imsi.msin_digit2;
    imsi.u.num.digit13 = msg->mobile_identity.imsi.msin_digit3;
    imsi.u.num.digit14 = msg->mobile_identity.imsi.msin_digit4;
    imsi.u.num.digit15 = msg->mobile_identity.imsi.msin_digit5;
    imsi.u.num.parity  = 0x0f;
    imsi.length        = msg->mobile_identity.imsi.numOfValidImsiDigits;

  } else if (
      msg->mobile_identity.imei.type_of_identity == MOBILE_IDENTITY_IMEI) {
    /*
     * Get the IMEI
     */
    p_imei            = &imei;
    imei.u.num.tac1   = msg->mobile_identity.imei.identity_digit1;
    imei.u.num.tac2   = msg->mobile_identity.imei.identity_digit2;
    imei.u.num.tac3   = msg->mobile_identity.imei.identity_digit3;
    imei.u.num.tac4   = msg->mobile_identity.imei.identity_digit4;
    imei.u.num.tac5   = msg->mobile_identity.imei.identity_digit5;
    imei.u.num.tac6   = msg->mobile_identity.imei.identity_digit6;
    imei.u.num.tac7   = msg->mobile_identity.imei.identity_digit7;
    imei.u.num.tac8   = msg->mobile_identity.imei.identity_digit8;
    imei.u.num.snr1   = msg->mobile_identity.imei.identity_digit9;
    imei.u.num.snr2   = msg->mobile_identity.imei.identity_digit10;
    imei.u.num.snr3   = msg->mobile_identity.imei.identity_digit11;
    imei.u.num.snr4   = msg->mobile_identity.imei.identity_digit12;
    imei.u.num.snr5   = msg->mobile_identity.imei.identity_digit13;
    imei.u.num.snr6   = msg->mobile_identity.imei.identity_digit14;
    imei.u.num.cdsd   = msg->mobile_identity.imei.identity_digit15;
    imei.u.num.parity = msg->mobile_identity.imei.odd_even;
  }
#if 0  
  else if ( msg->mobile_identity.guti.type_of_identity == MOBILE_IDENTITY_IMEISV) {
    /*
     * Get the IMEISV
     */
    p_imeisv            = &imeisv;
    imeisv.u.num.tac1   = msg->mobile_identity.imeisv.tac1;
    imeisv.u.num.tac2   = msg->mobile_identity.imeisv.tac2;
    imeisv.u.num.tac3   = msg->mobile_identity.imeisv.tac3;
    imeisv.u.num.tac4   = msg->mobile_identity.imeisv.tac4;
    imeisv.u.num.tac5   = msg->mobile_identity.imeisv.tac5;
    imeisv.u.num.tac6   = msg->mobile_identity.imeisv.tac6;
    imeisv.u.num.tac7   = msg->mobile_identity.imeisv.tac7;
    imeisv.u.num.tac8   = msg->mobile_identity.imeisv.tac8;
    imeisv.u.num.snr1   = msg->mobile_identity.imeisv.snr1;
    imeisv.u.num.snr2   = msg->mobile_identity.imeisv.snr2;
    imeisv.u.num.snr3   = msg->mobile_identity.imeisv.snr3;
    imeisv.u.num.snr4   = msg->mobile_identity.imeisv.snr4;
    imeisv.u.num.snr5   = msg->mobile_identity.imeisv.snr5;
    imeisv.u.num.snr6   = msg->mobile_identity.imeisv.snr6;
    imeisv.u.num.svn1   = msg->mobile_identity.imeisv.svn1;
    imeisv.u.num.svn2   = msg->mobile_identity.imeisv.svn2;
    imeisv.u.num.parity = msg->mobile_identity.imeisv.oddeven;
  }
#endif
  else if (msg->mobile_identity.tmsi.type_of_identity == MOBILE_IDENTITY_TMSI) {
    /*
     * Get the TMSI
     */
    p_tmsi = &tmsi;
    tmsi   = ((tmsi_t) msg->mobile_identity.tmsi.m5g_tmsi_1) << 24;
    tmsi |= (((tmsi_t) msg->mobile_identity.tmsi.m5g_tmsi_2) << 16);
    tmsi |= (((tmsi_t) msg->mobile_identity.tmsi.m5g_tmsi_3) << 8);
    tmsi |= ((tmsi_t) msg->mobile_identity.tmsi.m5g_tmsi_4);
  }

  /*
   * Execute the identification completion procedure
   */
  // rc = amf_proc_identification_complete(
  //    ue_id, p_imsi, p_imei, p_imeisv, (uint32_t*) (p_tmsi));
  OAILOG_FUNC_RETURN(LOG_NAS_AMF, rc);
}
//--------------------------------------------------------------------------
void amf_ctx_set_attribute_valid(
    amf_context_t* ctxt, const uint32_t attribute_bit_pos) {
  ctxt->member_present_mask |= attribute_bit_pos;
  ctxt->member_valid_mask |= attribute_bit_pos;
}
//--------------------------------------------------------------------------
void amf_ctx_set_attribute_present(
    amf_context_t* ctxt, const int attribute_bit_pos) {
  ctxt->member_present_mask |= attribute_bit_pos;
}
//------------------------------------------------------------------------------
nas_amf_ident_proc_t* get_5g_nas_common_procedure_identification(
    const amf_context_t* ctxt) {
  return (nas_amf_ident_proc_t*) nas_proc_indt.get_nas5g_common_procedure(
      ctxt,
      AMF_COMM_PROC_IDENT);  // TODO-RECHECK
}
//----------------------------------------------------------------------------
/* Set IMEI, mark it as valid */
void amf_ctx_set_valid_imei(amf_context_t* const ctxt, imei_t* imei) {
  ctxt->_imei = *imei;
  amf_ctx_set_attribute_valid(ctxt, AMF_CTXT_MEMBER_IMEI);
}

//------------------------------------------------------------------------------

/****************************************************************************
 **                                                                        **
 ** Name:    amf_proc_identification_complete()                            **
 **                                                                        **
 ** Description: Performs the identification completion procedure executed **
 **      by the network.                                                   **
 **                                                                        **
 **              3GPP TS 24.501, section 5.4.3.4                           **
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
 **      Others:    amf_data, T3570                                       **
 **                                                                        **
 ***************************************************************************/
int amf_identity_msg::amf_proc_identification_complete(
    const amf_ue_ngap_id_t ue_id, imsi_t* const imsi, imei_t* const imei,
    imeisv_t* const imeisv, uint32_t* const tmsi, guti_m5_t* amf_ctx_guti) {
  // imeisv_t* const imeisv, uint32_t* const tmsi) {
  OAILOG_FUNC_IN(LOG_NAS_AMF);
  int rc = RETURNerror;
  amf_sap_t amf_sap;
  amf_sap_c amf_sap_identy;
  amf_context_t* amf_ctx = NULL;

  OAILOG_INFO(
      LOG_NAS_AMF,
      "AMF-TEST: Identification procedure complete for "
      "(ue_id=" AMF_UE_NGAP_ID_FMT ")\n",
      ue_id);

  // Get the UE context
  // ue_m5gmm_context_s* ue_mm_context =
  //    amf_ue_context_exists_amf_ue_ngap_id(ue_id);
  ue_m5gmm_context_s* ue_mm_context =
      &ue_m5gmm_global_context;  // TODO AMF-TEST:global var to temporarily
                                 // store context inserted to ht

  if (ue_mm_context) {
    OAILOG_INFO(LOG_AMF_APP, "AMF-TEST: , from %s\n", __FUNCTION__);
    amf_ctx = &ue_mm_context->amf_context;
    OAILOG_INFO(
        LOG_AMF_APP, "AMF-TEST:amf_procedures:%p\n", amf_ctx->amf_procedures);
    // nas_amf_ident_proc_t* ident_proc =
    //     get_5g_nas_common_procedure_identification(amf_ctx);

    // if (ident_proc) {
    /*
     * Stop timer T3570
     */
    void* timer_callback_args = NULL;
    // nas_stop_T3570(ue_id, &ident_proc->T3570, timer_callback_args);  //
    // TODO RECHECK

    if (imsi) {
      OAILOG_INFO(LOG_AMF_APP, "AMF-TEST: , from %s\n", __FUNCTION__);
      /*
       * Update the IMSI
       */
      // imsi64_t imsi64 = imsi_to_imsi64(imsi);
      imsi64_t imsi64 = amf_imsi_to_imsi64(imsi);
      amf_ctx_set_valid_imsi(amf_ctx, imsi, imsi64);
      amf_context_upsert_imsi(amf_ctx);
      amf_ctx->_imsi64 = imsi64;  // TODO AMF_TEST global var to temporarily
                                  // store context inserted to ht  //pdu_change
      amf_ctx->_imsi.length = 8;
      // 15;  // TODO AMF_TEST global var to temporarily store context inserted
      // to ht //context set for IMSI64_TO_STRING() used in
      // amf_smf_send;
      // assign amf_ctx_guti to amf_context guti for quick access
      amf_ctx->_m5_guti = *amf_ctx_guti;
    } else if (imei) {
      /*
       * Update the IMEI
       */
      amf_ctx_set_valid_imei(amf_ctx, imei);  // TODO
    } else if (tmsi) {
      /*
       * Update the GUTI TODO later
       */
      AssertFatal(
          false,
          "TODO, should not happen because this type of identity is not "
          "requested by AMF");
    }

    /*
     * Notify EMM that the identification procedure successfully completed
     */
    // amf_sap.primitive               = AMFREG_COMMON_PROC_CNF;
#if 0
    amf_sap.primitive           = (amf_primitive_t) AMFCN_IDENTITY_PARAM_RES;
    amf_sap.u.amf_reg.ue_id     = ue_id;
    amf_sap.u.amf_reg.ctx       = amf_ctx;
    amf_sap.u.amf_reg.notify    = true;
    amf_sap.u.amf_reg.free_proc = true;
    // amf_sap.u.amf_reg.u.common_proc = &ident_proc->amf_com_proc;
    rc = amf_sap_identy.amf_sap_send(&amf_sap);

    // } TODO AMF-TEST commented from above  if (ident_proc)
    // else ignore the response if procedure not found
    // In success case of identity response, invoke Authentication Request
#endif  // TODO AMF-TEST, why is amf_sap.u.amf_reg set for
        // AMFCN_IDENTITY_PARAM_RES
    OAILOG_INFO(
        LOG_AMF_APP,
        "AMF-TEST: , calling amf_registration_success_identification_cb from "
        "%s\n",
        __FUNCTION__);
    amf_registration_procedure::amf_registration_success_identification_cb(
        amf_ctx);
  }  // else ignore the response if ue context not found

  OAILOG_FUNC_RETURN(LOG_NAS_AMF, rc);
}
// Generating GUTI based on SUPI/IMSI received from identity message.
void amf_app_generate_guti_on_supi(
    amf_guti_m5g_t* amf_guti, supi_as_imsi_t* supi_imsi)
// amf_plmn_t supi_imsi_plmn)
{
  /* Generate GUTI with 5g-tmsi as rand value */
  amf_guti->guamfi.plmn.mcc_digit1 = supi_imsi->plmn.mcc_digit1;
  amf_guti->guamfi.plmn.mcc_digit2 = supi_imsi->plmn.mcc_digit2;
  amf_guti->guamfi.plmn.mcc_digit3 = supi_imsi->plmn.mcc_digit3;
  amf_guti->guamfi.plmn.mnc_digit1 = supi_imsi->plmn.mnc_digit1;
  amf_guti->guamfi.plmn.mnc_digit2 = supi_imsi->plmn.mnc_digit2;
  amf_guti->guamfi.plmn.mnc_digit3 = supi_imsi->plmn.mnc_digit3;

#if 0
     amf_guti->guamfi.plmn.mcc_digit1 = supi_imsi_plmn.mcc_digit1;
     amf_guti->guamfi.plmn.mcc_digit2 = supi_imsi_plmn.mcc_digit2;
     amf_guti->guamfi.plmn.mcc_digit3 = supi_imsi_plmn.mcc_digit3;
     amf_guti->guamfi.plmn.mnc_digit1 = supi_imsi_plmn.mnc_digit1;
     amf_guti->guamfi.plmn.mnc_digit2 = supi_imsi_plmn.mnc_digit2;
     amf_guti->guamfi.plmn.mnc_digit3 = supi_imsi_plmn.mnc_digit3;
#endif
  // tmsi value is 4 octet random value.
  amf_guti->m_tmsi = (uint32_t) rand();

  // Filling data from amf_config file considering only one gNB
  amf_config_read_lock(&amf_config);
  amf_guti->guamfi.amf_regionid = amf_config.guamfi.guamfi[0].amf_code;
  amf_guti->guamfi.amf_set_id   = amf_config.guamfi.guamfi[0].amf_gid;
  amf_guti->guamfi.amf_pointer  = amf_config.guamfi.guamfi[0].amf_Pointer;

  amf_config_unlock(&amf_config);
  return;
}

}  // namespace magma5g

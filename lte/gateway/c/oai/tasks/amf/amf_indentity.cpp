/**
 *copyright 2020 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
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
#ifdef __cplusplus
}
#endif
#include "amf_app_ue_context_and_proc.h"
namespace magma5g {
// namespace NR_amf_msg
{
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
  int amf_cn_identity_res(
      amf_ue_ngap_id_t ue_id, identity_response_msg * msg, int* amf_cause,
      const amf_nas_message_decode_status_t* status) {
    OAILOG_FUNC_IN(LOG_NAS_AMF);
    int rc = RETURNok;

    OAILOG_INFO(
        LOG_NAS_AMF, "AMFAS-SAP - Received Identity Response message\n");
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

    if (msg->mobileidentity.imsi.typeofidentity == MOBILE_IDENTITY_IMSI) {
      /*
       * Get the IMSI
       */
      p_imsi             = &imsi;
      imsi.u.num.digit1  = msg->mobileidentity.imsi.digit1;
      imsi.u.num.digit2  = msg->mobileidentity.imsi.digit2;
      imsi.u.num.digit3  = msg->mobileidentity.imsi.digit3;
      imsi.u.num.digit4  = msg->mobileidentity.imsi.digit4;
      imsi.u.num.digit5  = msg->mobileidentity.imsi.digit5;
      imsi.u.num.digit6  = msg->mobileidentity.imsi.digit6;
      imsi.u.num.digit7  = msg->mobileidentity.imsi.digit7;
      imsi.u.num.digit8  = msg->mobileidentity.imsi.digit8;
      imsi.u.num.digit9  = msg->mobileidentity.imsi.digit9;
      imsi.u.num.digit10 = msg->mobileidentity.imsi.digit10;
      imsi.u.num.digit11 = msg->mobileidentity.imsi.digit11;
      imsi.u.num.digit12 = msg->mobileidentity.imsi.digit12;
      imsi.u.num.digit13 = msg->mobileidentity.imsi.digit13;
      imsi.u.num.digit14 = msg->mobileidentity.imsi.digit14;
      imsi.u.num.digit15 = msg->mobileidentity.imsi.digit15;
      imsi.u.num.parity  = 0x0f;
      imsi.length        = msg->mobileidentity.imsi.numOfValidImsiDigits;

    } else if (
        msg->mobileidentity.imei.typeofidentity == MOBILE_IDENTITY_IMEI) {
      /*
       * Get the IMEI
       */
      p_imei            = &imei;
      imei.u.num.tac1   = msg->mobileidentity.imei.tac1;
      imei.u.num.tac2   = msg->mobileidentity.imei.tac2;
      imei.u.num.tac3   = msg->mobileidentity.imei.tac3;
      imei.u.num.tac4   = msg->mobileidentity.imei.tac4;
      imei.u.num.tac5   = msg->mobileidentity.imei.tac5;
      imei.u.num.tac6   = msg->mobileidentity.imei.tac6;
      imei.u.num.tac7   = msg->mobileidentity.imei.tac7;
      imei.u.num.tac8   = msg->mobileidentity.imei.tac8;
      imei.u.num.snr1   = msg->mobileidentity.imei.snr1;
      imei.u.num.snr2   = msg->mobileidentity.imei.snr2;
      imei.u.num.snr3   = msg->mobileidentity.imei.snr3;
      imei.u.num.snr4   = msg->mobileidentity.imei.snr4;
      imei.u.num.snr5   = msg->mobileidentity.imei.snr5;
      imei.u.num.snr6   = msg->mobileidentity.imei.snr6;
      imei.u.num.cdsd   = msg->mobileidentity.imei.cdsd;
      imei.u.num.parity = msg->mobileidentity.imei.oddeven;
    } else if (
        msg->mobileidentity.imeisv.typeofidentity == MOBILE_IDENTITY_IMEISV) {
      /*
       * Get the IMEISV
       */
      p_imeisv            = &imeisv;
      imeisv.u.num.tac1   = msg->mobileidentity.imeisv.tac1;
      imeisv.u.num.tac2   = msg->mobileidentity.imeisv.tac2;
      imeisv.u.num.tac3   = msg->mobileidentity.imeisv.tac3;
      imeisv.u.num.tac4   = msg->mobileidentity.imeisv.tac4;
      imeisv.u.num.tac5   = msg->mobileidentity.imeisv.tac5;
      imeisv.u.num.tac6   = msg->mobileidentity.imeisv.tac6;
      imeisv.u.num.tac7   = msg->mobileidentity.imeisv.tac7;
      imeisv.u.num.tac8   = msg->mobileidentity.imeisv.tac8;
      imeisv.u.num.snr1   = msg->mobileidentity.imeisv.snr1;
      imeisv.u.num.snr2   = msg->mobileidentity.imeisv.snr2;
      imeisv.u.num.snr3   = msg->mobileidentity.imeisv.snr3;
      imeisv.u.num.snr4   = msg->mobileidentity.imeisv.snr4;
      imeisv.u.num.snr5   = msg->mobileidentity.imeisv.snr5;
      imeisv.u.num.snr6   = msg->mobileidentity.imeisv.snr6;
      imeisv.u.num.svn1   = msg->mobileidentity.imeisv.svn1;
      imeisv.u.num.svn2   = msg->mobileidentity.imeisv.svn2;
      imeisv.u.num.parity = msg->mobileidentity.imeisv.oddeven;
    } else if (
        msg->mobileidentity.tmsi.typeofidentity == MOBILE_IDENTITY_TMSI) {
      /*
       * Get the TMSI
       */
      p_tmsi = &tmsi;
      tmsi   = ((tmsi_t) msg->mobileidentity.tmsi.tmsi[0]) << 24;
      tmsi |= (((tmsi_t) msg->mobileidentity.tmsi.tmsi[1]) << 16);
      tmsi |= (((tmsi_t) msg->mobileidentity.tmsi.tmsi[2]) << 8);
      tmsi |= ((tmsi_t) msg->mobileidentity.tmsi.tmsi[3]);
    }

    /*
     * Execute the identification completion procedure
     */
    rc = amf_proc_identification_complete(
        ue_id, p_imsi, p_imei, p_imeisv, (uint32_t*) (p_tmsi));
    OAILOG_FUNC_RETURN(LOG_NAS_AMF, rc);
  }

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
  int amf_proc_identification_complete(
      const amf_ue_ngap_id_t ue_id, imsi_t* const imsi, imei_t* const imei,
      imeisv_t* const imeisv, uint32_t* const tmsi) {
    OAILOG_FUNC_IN(LOG_NAS_AMF);
    int rc = RETURNerror;
    amf_sap_t amf_sap;
    amf_context_t* amf_ctx = NULL;

    OAILOG_INFO(
        LOG_NAS_AMF,
        "AMF-PROC  - Identification complete (ue_id=" AMF_UE_NGAP_ID_FMT ")\n",
        ue_id);

    // Get the UE context
    ue_m5gmm_context_s* ue_mm_context =
        amf_ue_context_exists_amf_ue_ngap_id(ue_id);  // TODO
    if (ue_mm_context) {
      amf_ctx = &ue_mm_context->amf_context;
      nas_amf_ident_proc_t* ident_proc =
          get_5g_nas_common_procedure_identification(amf_ctx);  // TODO

      if (ident_proc) {
        REQUIREMENT_3GPP_24_501(R15_5_4_3_4);  // TODO
        /*
         * Stop timer T3570
         */
        void* timer_callback_args = NULL;
        nas_stop_T3570(ue_id, &ident_proc->T3570, timer_callback_args);  // TODO

        if (imsi) {
          /*
           * Update the IMSI
           */
          imsi64_t imsi64 = imsi_to_imsi64(imsi);
          amf_ctx_set_valid_imsi(amf_ctx, imsi, imsi64);  // TODO
          amf_context_upsert_imsi(&_amf_data, amf_ctx);   // TODO
        } else if (imei) {
          /*
           * Update the IMEI
           */
          amf_ctx_set_valid_imei(amf_ctx, imei);  // TODO
        } else if (imeisv) {
          /*
           * Update the IMEISV
           */
          amf_ctx_set_valid_imeisv(amf_ctx, imeisv);  // TODO
        } else if (tmsi) {
          /*
           * Update the GUTI
           */
          AssertFatal(
              false,
              "TODO, should not happen because this type of identity is not "
              "requested by AMF");
        }

        /*
         * Notify EMM that the identification procedure successfully completed
         */
        amf_sap.primitive                      = AMFREG_COMMON_PROC_CNF;
        amf_sap.u.amf_reg.ue_id                = ue_id;
        amf_sap.u.amf_reg.ctx                  = amf_ctx;
        amf_sap.u.amf_reg.notify               = true;
        amf_sap.u.mf_reg.free_proc             = true;
        amf_sap.u.amf_reg.u.common.common_proc = &ident_proc->amf_com_proc;
        amf_sap.u.amf_reg.u.common.previous_amf_fsm_state =
            ident_proc->amf_com_proc.amf_proc.previous_amf_fsm_state;
        rc = amf_sap_send(&amf_sap);

      }  // else ignore the response if procedure not found
    }    // else ignore the response if ue context not found

    OAILOG_FUNC_RETURN(LOG_NAS_AMF, rc);
  }
}  // magma5g

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

  Source      amf_Security_Mode.cpp

  Version     0.1

  Date        2020/07/28

  Product     AMF stack

  Subsystem   from AMF to NGAP

  Author      Sandeep Kumar Mall

  Description Defines Access and Mobility Management Messages

*****************************************************************************/
#include <sstream>
#ifdef __cplusplus
extern "C" {
#endif
#include "log.h"
#include "3gpp_24.501.h"
#include "conversions.h"
#ifdef __cplusplus
}
#endif

#include <thread>
#include "amf_fsm.h"
#include "amf_recv.h"
#include "amf_sap.h"
using namespace std;

namespace magma5g {
extern ue_m5gmm_context_s
    ue_m5gmm_global_context;  // TODO AMF-TEST global var to temporarily store
                              // context inserted to ht
amf_sap_c amf_sap_seq;
nas_proc nas_proc_seq;

int amf_procedure_handler::amf_handle_securitycomplete_response(
    amf_ue_ngap_id_t ue_id, amf_nas_message_decode_status_t decode_status) {
  OAILOG_FUNC_IN(LOG_NAS_AMF);
  ue_m5gmm_context_s* ue_mm_context = NULL;
  amf_context_t* amf_ctx            = NULL;
  int rc                            = RETURNerror;

  OAILOG_INFO(
      LOG_NAS_AMF,
      "AMF_TEST: Security mode procedures complete for "
      "(ue_id=" AMF_UE_NGAP_ID_FMT ")\n",
      ue_id);
  /*
   * Get the UE context
   */
  //  ue_mm_context = amf_ue_context_exists_amf_ue_ngap_id(ue_id);
  ue_mm_context =
      &ue_m5gmm_global_context;  // TODO AMF-TEST global var to temporarily
                                 // store context inserted to ht
  if (ue_mm_context) {
    amf_ctx = &ue_mm_context->amf_context;
  } else {
    OAILOG_FUNC_RETURN(LOG_NAS_AMF, RETURNerror);
  }

  nas_amf_smc_proc_t* smc_proc = get_nas5g_common_procedure_smc(amf_ctx);

  if (smc_proc) {
    /*
     * Stop timer T3560
     */
    // void* timer_callback_arg = NULL;
    // nas_stop_T3560(ue_id, &smc_proc->T3560, timer_callback_arg);

    /*
     * Release retransmission timer parameters
     */

    if (amf_ctx && IS_AMF_CTXT_PRESENT_SECURITY(amf_ctx)) {
      /*
       * Notify AMF that the authentication procedure successfully completed
       */
      amf_sap_t amf_sap;
      amf_sap.primitive = AMFCN_CS_RESPONSE;
      // amf_sap.primitive               = AMFREG_COMMON_PROC_CNF;
      amf_sap.u.amf_reg.ue_id         = ue_id;
      amf_sap.u.amf_reg.ctx           = amf_ctx;
      amf_sap.u.amf_reg.notify        = true;
      amf_sap.u.amf_reg.free_proc     = true;
      amf_sap.u.amf_reg.u.common_proc = &smc_proc->amf_com_proc;

      amf_ctx->_security.kenb_ul_count = amf_ctx->_security.ul_count;
      amf_ctx_set_attribute_valid(amf_ctx, AMF_CTXT_MEMBER_SECURITY);
      rc = amf_sap_seq.amf_sap_send(&amf_sap);
    }
    /* Nothing to do in
     * Calling SMC response success and triggering registration accept message*/
    amf_registration_procedure::amf_registration_success_security_cb(amf_ctx);
    // amf_registration_success_security_cb(amf_ctx);

  } else {
    OAILOG_ERROR(
        LOG_NAS_AMF,
        "AMF-PROC  - No 5GCN security context exists. Ignoring the Security "
        "Mode "
        "Complete message\n");
    rc = RETURNerror;
  }

  OAILOG_FUNC_RETURN(LOG_NAS_AMF, rc);
}
}  // namespace magma5g

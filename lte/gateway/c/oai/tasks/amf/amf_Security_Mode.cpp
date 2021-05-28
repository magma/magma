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
#include "3gpp_24.501.h"
#include "conversions.h"
#include "intertask_interface.h"
#include "intertask_interface_types.h"
#ifdef __cplusplus
}
#endif
#include "common_defs.h"
#include <thread>
#include "amf_fsm.h"
#include "amf_recv.h"
#include "amf_sap.h"

namespace magma5g {
extern task_zmq_ctx_s amf_app_task_zmq_ctx;

/****************************************************************************
 **                                                                        **
 ** Name:    amf_handle_security_complete_response()                        **
 **                                                                        **
 ** Description: Procedure to indicate Security mode procedure completed   **
 **                                                                        **
 **                                                                        **
 ***************************************************************************/
int amf_handle_security_complete_response(
    amf_ue_ngap_id_t ue_id, amf_nas_message_decode_status_t decode_status) {
  OAILOG_FUNC_IN(LOG_NAS_AMF);
  ue_m5gmm_context_s* ue_mm_context = NULL;
  amf_context_t* amf_ctx            = NULL;
  int rc                            = RETURNerror;
  OAILOG_INFO(
      LOG_NAS_AMF,
      "Security mode procedures complete for "
      "(ue_id=" AMF_UE_NGAP_ID_FMT ")\n",
      ue_id);
  /*
   * Get the UE context
   */
  ue_mm_context = amf_ue_context_exists_amf_ue_ngap_id(ue_id);

  if (ue_mm_context) {
    amf_ctx = &ue_mm_context->amf_context;
  } else {
    OAILOG_ERROR(LOG_AMF_APP, "ue context not found for the ue_id=%u\n", ue_id);
    OAILOG_FUNC_RETURN(LOG_NAS_AMF, RETURNerror);
  }
  nas_amf_smc_proc_t* smc_proc = get_nas5g_common_procedure_smc(amf_ctx);
  if (smc_proc) {
    /*
     * TODO:Stop timer T3560 This to be taken care in upcoming PR
     */

    stop_timer(&amf_app_task_zmq_ctx, smc_proc->T3560.id);
    OAILOG_DEBUG(LOG_AMF_APP, "Timer: After stopping SMC MODE timer \n");
    smc_proc->T3560.id = NAS5G_TIMER_INACTIVE_ID;

    OAILOG_DEBUG(
          LOG_AMF_APP, "ue_context_request : %d", ue_mm_context->ue_context_request);
    if (amf_ctx && IS_AMF_CTXT_PRESENT_SECURITY(amf_ctx)) {
      if (M5G_UEContextRequest_requested != ue_mm_context->ue_context_request) {
        /*
         * Notify AMF that the authentication procedure successfully completed
         */
        amf_sap_t amf_sap;
        amf_sap.primitive                = AMFCN_CS_RESPONSE;
        amf_sap.u.amf_reg.ue_id          = ue_id;
        amf_sap.u.amf_reg.ctx            = amf_ctx;
        amf_sap.u.amf_reg.notify         = true;
        amf_sap.u.amf_reg.free_proc      = true;
        amf_sap.u.amf_reg.u.common_proc  = &smc_proc->amf_com_proc;
        amf_ctx->_security.kenb_ul_count = amf_ctx->_security.ul_count;
        amf_ctx_set_attribute_valid(amf_ctx, AMF_CTXT_MEMBER_SECURITY);
        rc = amf_sap_send(&amf_sap);
      } else {
        // Send the intial context setup request
        amf_registration_success_security_cb(amf_ctx);
      }
    }

    OAILOG_INFO(LOG_AMF_APP, " mm_state %d", ue_mm_context->mm_state);
    OAILOG_INFO(LOG_AMF_APP, "ue_m5gmm_context %p\n", ue_mm_context);
    ue_state_handle_message_initial(
        // ue_mm_context->mm_state, STATE_EVENT_SEC_MODE_COMPLETE, SESSION_NULL,
        (magma5g::m5gmm_state_t) 5, STATE_EVENT_SEC_MODE_COMPLETE, SESSION_NULL,
        ue_mm_context, amf_ctx);

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

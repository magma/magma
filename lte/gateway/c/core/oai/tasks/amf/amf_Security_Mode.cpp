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
#include "lte/gateway/c/core/oai/common/log.h"
#include "lte/gateway/c/core/oai/lib/3gpp/3gpp_24.501.h"
#include "lte/gateway/c/core/oai/common/conversions.h"
#include "lte/gateway/c/core/oai/lib/itti/intertask_interface.h"
#include "lte/gateway/c/core/oai/lib/itti/intertask_interface_types.h"
#ifdef __cplusplus
}
#endif
#include "lte/gateway/c/core/common/common_defs.h"
#include <thread>
#include "lte/gateway/c/core/oai/tasks/amf/amf_fsm.hpp"
#include "lte/gateway/c/core/oai/tasks/amf/amf_recv.hpp"
#include "lte/gateway/c/core/oai/tasks/amf/amf_sap.hpp"
#include "lte/gateway/c/core/oai/tasks/amf/amf_app_timer_management.hpp"

namespace magma5g {

/****************************************************************************
 **                                                                        **
 ** Name:    amf_handle_security_complete_response()                        **
 **                                                                        **
 ** Description: Procedure to indicate Security mode procedure completed   **
 **                                                                        **
 **                                                                        **
 ***************************************************************************/
status_code_e amf_handle_security_complete_response(
    amf_ue_ngap_id_t ue_id, amf_nas_message_decode_status_t decode_status) {
  OAILOG_FUNC_IN(LOG_NAS_AMF);
  ue_m5gmm_context_s* ue_mm_context = NULL;
  amf_context_t* amf_ctx = NULL;
  status_code_e rc = RETURNok;
  OAILOG_DEBUG(LOG_NAS_AMF,
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
    OAILOG_ERROR(LOG_AMF_APP,
                 "ue context not found for the UE ID " AMF_UE_NGAP_ID_FMT,
                 ue_id);
    OAILOG_FUNC_RETURN(LOG_NAS_AMF, RETURNerror);
  }
  nas_amf_smc_proc_t* smc_proc = get_nas5g_common_procedure_smc(amf_ctx);
  if (smc_proc) {
    amf_app_stop_timer(smc_proc->T3560.id);
    OAILOG_DEBUG(LOG_AMF_APP,
                 "Timer: After stopping timer T3560 for securiy mode command"
                 " with id: %lu and UE ID: " AMF_UE_NGAP_ID_FMT,
                 smc_proc->T3560.id, ue_id);
    smc_proc->T3560.id = NAS5G_TIMER_INACTIVE_ID;

    // Send s6a update location request
    if (amf_send_n11_update_location_req(ue_mm_context->amf_ue_ngap_id) ==
        RETURNerror) {
      OAILOG_ERROR(LOG_AMF_APP,
                   "update location request failed for amf_ue_ngap_id "
                   ": " AMF_UE_NGAP_ID_FMT,
                   ue_mm_context->amf_ue_ngap_id);
    }

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

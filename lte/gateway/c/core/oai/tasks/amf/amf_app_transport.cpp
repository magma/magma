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
#include <thread>
#ifdef __cplusplus
extern "C" {
#endif
#include "lte/gateway/c/core/common/dynamic_memory_check.h"
#include "lte/gateway/c/core/oai/common/log.h"
#include "lte/gateway/c/core/oai/lib/itti/intertask_interface.h"
#include "lte/gateway/c/core/oai/lib/itti/intertask_interface_types.h"
#ifdef __cplusplus
}
#endif
#include "lte/gateway/c/core/common/common_defs.h"
#include "lte/gateway/c/core/oai/tasks/amf/amf_app_state_manager.hpp"

namespace magma5g {

extern task_zmq_ctx_t amf_app_task_zmq_ctx;

/* Handles NAS encoded message and sends it to NGAP task */
status_code_e amf_app_handle_nas_dl_req(const amf_ue_ngap_id_t ue_id,
                                        bstring nas_msg,
                                        nas5g_error_code_t transaction_status) {
  OAILOG_FUNC_IN(LOG_AMF_APP);
  MessageDef* message_p = NULL;
  status_code_e rc = RETURNok;
  gnb_ue_ngap_id_t gnb_ue_ngap_id = 0;
  message_p = itti_alloc_new_message(TASK_AMF_APP, NGAP_NAS_DL_DATA_REQ);
  amf_app_desc_t* amf_app_desc_p = get_amf_nas_state(false);

  if (!amf_app_desc_p) {
    OAILOG_ERROR(
        LOG_AMF_APP,
        "Downlink NAS transport: Failed to get global amf_app_desc context \n");
    OAILOG_FUNC_RETURN(LOG_AMF_APP, RETURNerror);
  }
  ue_m5gmm_context_s* ue_context = amf_ue_context_exists_amf_ue_ngap_id(ue_id);
  if (ue_context) {
    gnb_ue_ngap_id = ue_context->gnb_ue_ngap_id;
  } else {
    OAILOG_ERROR(LOG_AMF_APP,
                 "ue context not found for the UE ID = " AMF_UE_NGAP_ID_FMT,
                 ue_id);
    OAILOG_FUNC_RETURN(LOG_AMF_APP, rc);
  }

  NGAP_NAS_DL_DATA_REQ(message_p).gnb_ue_ngap_id = gnb_ue_ngap_id;
  NGAP_NAS_DL_DATA_REQ(message_p).amf_ue_ngap_id = ue_id;
  NGAP_NAS_DL_DATA_REQ(message_p).nas_msg = bstrcpy(nas_msg);
  bdestroy_wrapper(&nas_msg);
  message_p->ittiMsgHeader.imsi = ue_context->amf_context.imsi64;

  rc = amf_send_msg_to_task(&amf_app_task_zmq_ctx, TASK_NGAP, message_p);

  if (transaction_status != M5G_AS_SUCCESS) {
    if (M5G_AS_TERMINATED_NAS == transaction_status) {
      ue_context->ue_context_rel_cause = NGAP_IMPLICIT_CONTEXT_RELEASE;
    }
    amf_ctx_release_ue_context(ue_context, NGAP_NAS_AUTHENTICATION_FAILURE);
  }
  OAILOG_FUNC_RETURN(LOG_AMF_APP, rc);
}
}  // namespace magma5g

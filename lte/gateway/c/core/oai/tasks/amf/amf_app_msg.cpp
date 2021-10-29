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

#ifdef __cplusplus
extern "C" {
#endif
#include "lte/gateway/c/core/oai/lib/itti/intertask_interface.h"
#include "lte/gateway/c/core/oai/common/log.h"
#ifdef __cplusplus
}
#endif
#include "lte/gateway/c/core/oai/tasks/amf/amf_app_ue_context_and_proc.h"

namespace magma5g {
extern task_zmq_ctx_t amf_app_task_zmq_ctx;

/****************************************************************************
 **                                                                        **
 ** name:    amf_app_ue_context_release()                             **
 **                                                                        **
 ** description: Send itti mesage to ngap task to send UE Context Release  **
 **              Request                                                   **
 **                                                                        **
 ** inputs:  ue_context_p: Pointer to UE context amf_cause: failed cause   **
 **                                                                        **
 ***************************************************************************/
void amf_app_ue_context_release(
    ue_m5gmm_context_s* ue_context_p, ngap_Cause_t cause) {
  MessageDef* message_p;
  OAILOG_FUNC_IN(LOG_AMF_APP);
  message_p =
      itti_alloc_new_message(TASK_AMF_APP, NGAP_UE_CONTEXT_RELEASE_COMMAND);

  if (message_p == NULL) {
    OAILOG_ERROR(LOG_AMF_APP, "message is null");
    OAILOG_FUNC_OUT(LOG_AMF_APP);
  }
  memset(
      &NGAP_UE_CONTEXT_RELEASE_COMMAND(message_p), 0,
      sizeof(itti_ngap_ue_context_release_command_t));
  NGAP_UE_CONTEXT_RELEASE_COMMAND(message_p).amf_ue_ngap_id =
      ue_context_p->amf_ue_ngap_id;
  NGAP_UE_CONTEXT_RELEASE_COMMAND(message_p).gnb_ue_ngap_id =
      ue_context_p->gnb_ue_ngap_id;
  NGAP_UE_CONTEXT_RELEASE_COMMAND(message_p).cause =
      (Ngcause) cause.ngapCause_u.nas;
  message_p->ittiMsgHeader.imsi = ue_context_p->amf_context.imsi64;
  amf_send_msg_to_task(&amf_app_task_zmq_ctx, TASK_NGAP, message_p);
  OAILOG_FUNC_OUT(LOG_AMF_APP);
}
}  // namespace magma5g

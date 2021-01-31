/****************************************************************************
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
 ****************************************************************************/
/*****************************************************************************

  Source      amf_app_msg.cpp

  Version     0.1

  Date        2020/09/28

  Product     NAS stack

  Subsystem   Access and Mobility Management Function

  Author      Sandeep Kumar Mall

  Description Defines Access and Mobility Management Messages

*****************************************************************************/
#include "common_types.h"
//#include "amf_nas5g_proc.h"
#include "amf_fsm.h"
#include "amf_app_ue_context_and_proc.h"
//#include "amf_message.h"
#ifdef __cplusplus
extern "C" {
#endif
#include "intertask_interface.h"
#include "log.h"
#ifdef __cplusplus
}
#endif
//#include "amf_app_msg.h"

using namespace std;
namespace magma5g {
extern task_zmq_ctx_t amf_app_task_zmq_ctx;
/****************************************************************************
 **                                                                        **
 ** name:    amf_app_ue_context_release()                             **
 **                                                                        **
 ** description: Send itti mesage to ngap task to send UE Context Release  **
 **              Request                                                   **
 **                                                                        **
 ** inputs:  ue_context_p: Pointer to UE context amf_casue: failed cause   **
 **                                                                        **
 ***************************************************************************/
// void amf_app_ue_context_release(ue_m5gmm_context_s* ue_context_p,
// ngap_Cause_t cause)
void amf_app_msg::amf_app_ue_context_release(
    ue_m5gmm_context_s* ue_context_p, ngap_Cause_t cause) {
  MessageDef* message_p;

  OAILOG_FUNC_IN(LOG_AMF_APP);
  message_p = itti_alloc_new_message(
      TASK_AMF_APP, NGAP_UE_CONTEXT_RELEASE_COMMAND);  // ngap_message_def
  if (message_p == NULL) {
    // OAILOG_ERROR_UE(LOG_AMF_APP, ue_context_p->amf_context._imsi64,
    //    "Failed to allocate memory for NGAP_UE_CONTEXT_RELEASE_COMMAND \n");
    OAILOG_FUNC_OUT(LOG_AMF_APP);
  }

  // OAILOG_INFO_UE(LOG_AMF_APP, ue_context_p->amf_context._imsi64,
  //    "Sending UE Context Release Cmd to ngap for (ue_id = %u)\n"
  //    "UE Context Release Cause = (%d)\n",ue_context_p->amf_ue_ngap_id,
  //    cause);

  NGAP_UE_CONTEXT_RELEASE_COMMAND(message_p).amf_ue_ngap_id =
      ue_context_p->amf_ue_ngap_id;
  NGAP_UE_CONTEXT_RELEASE_COMMAND(message_p).gnb_ue_ngap_id =
      ue_context_p->gnb_ue_ngap_id;
  NGAP_UE_CONTEXT_RELEASE_COMMAND(message_p).cause =
      (Ngcause) cause.ngapCause_u.nas;

  message_p->ittiMsgHeader.imsi = ue_context_p->amf_context._imsi64;
  send_msg_to_task(&amf_app_task_zmq_ctx, TASK_NGAP, message_p);
  OAILOG_FUNC_OUT(LOG_AMF_APP);
}

}  // namespace magma5g

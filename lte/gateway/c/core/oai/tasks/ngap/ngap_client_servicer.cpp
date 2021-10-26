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

#include "include/ngap_client_servicer.h"
#include <memory>
extern "C" {
#include "common_defs.h"
#include "ngap_common.h"
#include "dynamic_memory_check.h"
}

namespace magma5g {

NGAPClientServicer::NGAPClientServicer() {}

NGAPClientServicer& NGAPClientServicer::getInstance() {
  static NGAPClientServicer instance;

  return instance;
}

status_code_e NGAPClientServicer::send_message_to_amf(
    task_zmq_ctx_t* task_zmq_ctx_p, task_id_t destination_task_id,
    MessageDef* message) {
  status_code_e ret = RETURNok;

#if !MME_UNIT_TEST
  ret = send_msg_to_task(task_zmq_ctx_p, destination_task_id, message);
#else  /* !MME_UNIT_TEST */
  OAILOG_DEBUG(LOG_NGAP, " Mock is Enabled \n");
  if (message->ittiMsgHeader.messageId == NGAP_INITIAL_UE_MESSAGE) {
    bdestroy_wrapper(&message->ittiMsg.ngap_initial_ue_message.nas);
  } else if (message->ittiMsgHeader.messageId == AMF_APP_UPLINK_DATA_IND) {
    bdestroy_wrapper(&message->ittiMsg.amf_app_ul_data_ind.nas_msg);
  } else if (message->ittiMsgHeader.messageId == SCTP_DATA_REQ) {
    bdestroy_wrapper(&message->ittiMsg.sctp_data_req.payload);
  }

  free(message);
#endif /* !MME_UNIT_TEST */

  return (ret);
}

}  // namespace magma5g

/****************************************************************************
 **                                                                        **
 ** Name:    ngap_send_msg_to_task()                                        **
 **                                                                        **
 ** Description:  wrapper api for itti send                                **
 **                                                                        **
 **                                                                        **
 ***************************************************************************/
status_code_e ngap_send_msg_to_task(
    task_zmq_ctx_t* task_zmq_ctx_p, task_id_t destination_task_id,
    MessageDef* message) {
  return (magma5g::NGAPClientServicer::getInstance().send_message_to_amf(
      task_zmq_ctx_p, destination_task_id, message));
}

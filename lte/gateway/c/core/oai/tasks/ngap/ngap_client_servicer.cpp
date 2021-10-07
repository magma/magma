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
#if !MME_UNIT_TEST
  return (send_msg_to_task(task_zmq_ctx_p, destination_task_id, message));
#else  /* !MME_UNIT_TEST */
  OAILOG_DEBUG(LOG_NGAP, " Mock is Enabled \n");
  if (message->ittiMsgHeader.messageId == NGAP_INITIAL_UE_MESSAGE) {
    bdestroy(NGAP_INITIAL_UE_MESSAGE(message).nas);
  }
  free(message);
  return (RETURNok);
#endif /* !MME_UNIT_TEST */
}

}  // namespace magma5g

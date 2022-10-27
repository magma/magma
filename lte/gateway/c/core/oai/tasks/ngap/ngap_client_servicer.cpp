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

#include "lte/gateway/c/core/oai/tasks/ngap/include/ngap_client_servicer.hpp"
#include <memory>
extern "C" {
#include "lte/gateway/c/core/common/common_defs.h"
#include "lte/gateway/c/core/common/dynamic_memory_check.h"
#include "lte/gateway/c/core/oai/tasks/ngap/ngap_common.h"
}

namespace magma5g {

NGAPClientServicer::NGAPClientServicer() {}

NGAPClientServicer& NGAPClientServicer::getInstance() {
  OAILOG_FUNC_IN(LOG_NGAP);
  static NGAPClientServicer instance;
  OAILOG_FUNC_RETURN(LOG_NGAP, instance);
}

status_code_e NGAPClientServicer::send_message_to_amf(
    task_zmq_ctx_t* task_zmq_ctx_p, task_id_t destination_task_id,
    MessageDef* message) {
  status_code_e ret = RETURNok;

  OAILOG_FUNC_IN(LOG_NGAP);
#if !MME_UNIT_TEST
  ret = send_msg_to_task(task_zmq_ctx_p, destination_task_id, message);
#else  /* !MME_UNIT_TEST */
  OAILOG_DEBUG(LOG_NGAP, " Mock is Enabled \n");
  msgtype_stack.push_back(ITTI_MSG_ID(message));
  itti_free_msg_content(message);
  free(message);
#endif /* !MME_UNIT_TEST */

  OAILOG_FUNC_RETURN(LOG_NGAP, ret);
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
status_code_e ngap_send_msg_to_task(task_zmq_ctx_t* task_zmq_ctx_p,
                                    task_id_t destination_task_id,
                                    MessageDef* message) {
  OAILOG_INFO(LOG_NGAP, "Sending msg to :[%s] id: [%d]-[%s]\n",
              itti_get_task_name(destination_task_id), ITTI_MSG_ID(message),
              ITTI_MSG_NAME(message));

  return (magma5g::NGAPClientServicer::getInstance().send_message_to_amf(
      task_zmq_ctx_p, destination_task_id, message));
}

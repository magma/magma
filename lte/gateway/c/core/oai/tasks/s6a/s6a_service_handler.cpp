/*
 * Licensed to the OpenAirInterface (OAI) Software Alliance under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The OpenAirInterface Software Alliance licenses this file to You under
 * the terms found in the LICENSE file in the root of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *-------------------------------------------------------------------------------
 * For more information about the OpenAirInterface (OAI) Software Alliance:
 *      contact@openairinterface.org
 */
#include <string.h>
#include <sys/types.h>

#ifdef __cplusplus
extern "C" {
#endif
#include "lte/gateway/c/core/oai/lib/itti/intertask_interface.h"
#include "lte/gateway/c/core/oai/common/common_types.h"
#include "lte/gateway/c/core/oai/lib/itti/intertask_interface_types.h"
#include "lte/gateway/c/core/oai/lib/itti/itti_types.h"
#ifdef __cplusplus
}
#endif

#include "lte/gateway/c/core/oai/tasks/s6a/s6a_defs.hpp"
#include "lte/gateway/c/core/oai/include/s6a_messages_types.hpp"

int delete_subscriber_request(const char* imsi, const uint imsi_len) {
  // send it to MME module for further processing
  MessageDef* message_p = NULL;
  s6a_cancel_location_req_t* s6a_cancel_location_req_p = NULL;
  message_p =
      DEPRECATEDitti_alloc_new_message_fatal(TASK_S6A, S6A_CANCEL_LOCATION_REQ);
  s6a_cancel_location_req_p = &message_p->ittiMsg.s6a_cancel_location_req;
  memcpy(s6a_cancel_location_req_p->imsi, imsi, imsi_len);
  s6a_cancel_location_req_p->imsi[imsi_len] = '\0';
  s6a_cancel_location_req_p->imsi_length = imsi_len;
  s6a_cancel_location_req_p->cancellation_type = SUBSCRIPTION_WITHDRAWL;
  int ret = 0;
  ret = send_msg_to_task(&s6a_task_zmq_ctx, TASK_MME_APP, message_p);
  return ret;
}

void handle_reset_request(void) {
  // send it to MME module for further processing
  MessageDef* message_p = NULL;
  message_p = DEPRECATEDitti_alloc_new_message_fatal(TASK_S6A, S6A_RESET_REQ);
  // TBD - To add support for partial reset
  send_msg_to_task(&s6a_task_zmq_ctx, TASK_MME_APP, message_p);
  return;
}

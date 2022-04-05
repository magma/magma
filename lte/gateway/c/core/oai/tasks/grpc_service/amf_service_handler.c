/*
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

#include "lte/gateway/c/core/common/common_defs.h"
#include "lte/gateway/c/core/oai/common/common_types.h"
#include "lte/gateway/c/core/oai/common/log.h"
#include "lte/gateway/c/core/oai/include/n11_messages_types.h"
#include "lte/gateway/c/core/oai/lib/itti/intertask_interface.h"
#include "lte/gateway/c/core/oai/lib/itti/intertask_interface_types.h"
#include "lte/gateway/c/core/oai/lib/itti/itti_types.h"

extern task_zmq_ctx_t grpc_service_task_zmq_ctx;

status_code_e send_n11_create_pdu_session_resp_itti(
    itti_n11_create_pdu_session_response_t* itti_msg) {
  OAILOG_DEBUG(LOG_UTIL,
               "Sending itti_n11_create_pdu_session_response to AMF \n");
  MessageDef* message_p = itti_alloc_new_message(
      TASK_GRPC_SERVICE, N11_CREATE_PDU_SESSION_RESPONSE);
  if (message_p == NULL) {
    OAILOG_ERROR(
        LOG_UTIL,
        "Failed to allocate memory for N11_CREATE_PDU_SESSION_RESPONSE\n");
    return RETURNerror;
  }
  message_p->ittiMsg.n11_create_pdu_session_response = *itti_msg;
  return send_msg_to_task(&grpc_service_task_zmq_ctx, TASK_AMF_APP, message_p);
}

int send_n11_notification_received_itti(
    itti_n11_received_notification_t* itti_msg) {
  OAILOG_INFO(LOG_UTIL,
              "Sending itti_n11_create_pdu_session_response to AMF \n");
  MessageDef* message_p =
      itti_alloc_new_message(TASK_GRPC_SERVICE, N11_NOTIFICATION_RECEIVED);
  message_p->ittiMsg.n11_notification_received = *itti_msg;
  return send_msg_to_task(&grpc_service_task_zmq_ctx, TASK_AMF_APP, message_p);
}

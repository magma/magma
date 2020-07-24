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
#include <conversions.h>

#include "common_types.h"
#include "intertask_interface.h"
#include "intertask_interface_types.h"
#include "itti_types.h"
#include "log.h"
#include "gx_messages_types.h"

extern task_zmq_ctx_t grpc_service_task_zmq_ctx;

int send_activate_bearer_request_itti(
    itti_gx_nw_init_actv_bearer_request_t* itti_msg) {
  OAILOG_DEBUG(LOG_SPGW_APP, "Sending nw_init_actv_bearer_request to SPGW \n");
  MessageDef* message_p = itti_alloc_new_message(
      TASK_GRPC_SERVICE, GX_NW_INITIATED_ACTIVATE_BEARER_REQ);
  message_p->ittiMsg.gx_nw_init_actv_bearer_request = *itti_msg;

  IMSI_STRING_TO_IMSI64((char*) itti_msg->imsi, &message_p->ittiMsgHeader.imsi);

  return send_msg_to_task(&grpc_service_task_zmq_ctx, TASK_SPGW_APP, message_p);
}

int send_deactivate_bearer_request_itti(
    itti_gx_nw_init_deactv_bearer_request_t* itti_msg) {
  OAILOG_DEBUG(LOG_SPGW_APP, "Sending spgw_nw_init_deactv_bearer_request\n");
  MessageDef* message_p = itti_alloc_new_message(
      TASK_GRPC_SERVICE, GX_NW_INITIATED_DEACTIVATE_BEARER_REQ);
  message_p->ittiMsg.gx_nw_init_deactv_bearer_request = *itti_msg;

  IMSI_STRING_TO_IMSI64((char*) itti_msg->imsi, &message_p->ittiMsgHeader.imsi);

  return send_msg_to_task(&grpc_service_task_zmq_ctx, TASK_SPGW_APP, message_p);
}

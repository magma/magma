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
 *-----------------------------------------------------------------------------
 * For more information about the OpenAirInterface (OAI) Software Alliance:
 *      contact@openairinterface.org
 */

extern "C" {
#include "log.h"
#include "mme_app_ha.h"
#include "ha_messages_types.h"
#include "conversions.h"
#include "intertask_interface.h"
#include "common_types.h"
#include "intertask_interface_types.h"
#include "itti_types.h"
}

extern task_zmq_ctx_t mme_app_task_zmq_ctx;

void mme_app_handle_ue_offload(ue_mm_context_t* ue_context_p) {
  MessageDef* message_p = itti_alloc_new_message(TASK_MME_APP, AGW_OFFLOAD_REQ);

  AGW_OFFLOAD_REQ(message_p).imsi_length =
      ue_context_p->emm_context._imsi.length;
  IMSI64_TO_STRING(
      ue_context_p->emm_context._imsi64,
      (char*) AGW_OFFLOAD_REQ(message_p).imsi,
      ue_context_p->emm_context._imsi.length);
  AGW_OFFLOAD_REQ(message_p).enb_offload_type = ANY;

  message_p->ittiMsgHeader.imsi = ue_context_p->emm_context._imsi64;
  send_msg_to_task(&mme_app_task_zmq_ctx, TASK_HA, message_p);
  return;
}

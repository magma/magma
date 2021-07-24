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

#include "service303_message_utils.h"

#include <stddef.h>

#include "assertions.h"
#include "intertask_interface.h"
#include "itti_types.h"

int send_app_health_to_service303(
    task_zmq_ctx_t* task_zmq_ctx_p, task_id_t origin_id, bool healthy) {
  MessageDef* message_p;
  if (healthy) {
    message_p = itti_alloc_new_message(origin_id, APPLICATION_HEALTHY_MSG);
  } else {
    message_p = itti_alloc_new_message(origin_id, APPLICATION_UNHEALTHY_MSG);
  }
  return send_msg_to_task(task_zmq_ctx_p, TASK_SERVICE303, message_p);
}

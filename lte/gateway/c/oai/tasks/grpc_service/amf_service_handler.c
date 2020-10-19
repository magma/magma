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

#include <string.h>
#include <sys/types.h>
#include <conversions.h>

#include "common_types.h"
#include "intertask_interface.h"
#include "intertask_interface_types.h"
#include "itti_types.h"
#include "log.h"
#include "amf_app_messages_types.h"

extern task_zmq_ctx_t grpc_service_task_zmq_ctx;

int send_smf_response_itti(itti_smf_response_t* itti_msg) {
  OAILOG_INFO(LOG_UTIL, "Sending itti_smf_response to AMF \n");
  MessageDef* message_p =
      itti_alloc_new_message(TASK_GRPC_SERVICE, SMF_RESPONSE);
  message_p->ittiMsg.smf_response = *itti_msg;
  return send_msg_to_task(&grpc_service_task_zmq_ctx, TASK_AMF_APP, message_p);
}

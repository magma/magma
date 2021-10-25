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

#include <arpa/inet.h>
#include <string.h>
#include <netinet/in.h>
#include <stdlib.h>
#include <sys/socket.h>
#include <lte/gateway/c/core/oai/common/conversions.h>

#include "lte/gateway/c/core/oai/lib/itti/intertask_interface.h"
#include "lte/gateway/c/core/oai/common/log.h"
#include "lte/gateway/c/core/oai/lib/mobility_client/MobilityClientAPI.h"
#include "lte/gateway/c/core/oai/tasks/sgw/sgw_defs.h"
#include "lte/gateway/c/core/oai/tasks/sgw/sgw_paging.h"
#include "lte/gateway/c/core/oai/lib/itti/intertask_interface_types.h"
#include "lte/gateway/c/core/oai/lib/itti/itti_types.h"
#include "lte/gateway/c/core/oai/include/s11_messages_types.h"

void sgw_send_paging_request(const struct in_addr* dest_ip) {
  OAILOG_DEBUG(
      TASK_SPGW_APP, "Paging procedure initiated for ue_ipv4: %x\n",
      dest_ip->s_addr);
  MessageDef* message_p                       = NULL;
  itti_s11_paging_request_t* paging_request_p = NULL;

  message_p =
      DEPRECATEDitti_alloc_new_message_fatal(TASK_SPGW_APP, S11_PAGING_REQUEST);
  paging_request_p = &message_p->ittiMsg.s11_paging_request;
  memset((void*) paging_request_p, 0, sizeof(itti_s11_paging_request_t));
  paging_request_p->ipv4_addr = *dest_ip;

  send_msg_to_task(&spgw_app_task_zmq_ctx, TASK_MME_APP, message_p);
  return;
}

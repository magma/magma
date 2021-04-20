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
#include <conversions.h>

#include "intertask_interface.h"
#include "log.h"
#include "MobilityClientAPI.h"
#include "sgw_defs.h"
#include "sgw_paging.h"
#include "intertask_interface_types.h"
#include "itti_types.h"
#include "s11_messages_types.h"

void sgw_send_paging_request(const struct in_addr* dest_ip) {
  OAILOG_DEBUG(
      TASK_SPGW_APP, "Paging procedure initiated for ue_ipv4: %x\n",
      dest_ip->s_addr);
  MessageDef* message_p                       = NULL;
  itti_s11_paging_request_t* paging_request_p = NULL;

  message_p        = itti_alloc_new_message(TASK_SPGW_APP, S11_PAGING_REQUEST);
  paging_request_p = &message_p->ittiMsg.s11_paging_request;
  memset((void*) paging_request_p, 0, sizeof(itti_s11_paging_request_t));
  paging_request_p->ipv4_addr = *dest_ip;

  send_msg_to_task(&spgw_app_task_zmq_ctx, TASK_MME_APP, message_p);
  return;
}

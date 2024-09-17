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

#include "lte/gateway/c/core/oai/tasks/sgw/sgw_paging.hpp"

#include <arpa/inet.h>
#include <string.h>
#include <netinet/in.h>
#include <stdlib.h>
#include <sys/socket.h>

#ifdef __cplusplus
extern "C" {
#endif
#include "lte/gateway/c/core/oai/common/conversions.h"
#include "lte/gateway/c/core/oai/common/log.h"
#include "lte/gateway/c/core/oai/lib/itti/intertask_interface.h"
#include "lte/gateway/c/core/oai/lib/itti/intertask_interface_types.h"
#include "lte/gateway/c/core/oai/lib/itti/itti_types.h"
#ifdef __cplusplus
}
#endif

#include "lte/gateway/c/core/oai/include/s11_messages_types.hpp"
#include "lte/gateway/c/core/oai/lib/mobility_client/MobilityClientAPI.hpp"
#include "lte/gateway/c/core/oai/tasks/sgw/sgw_handlers.hpp"

void sgw_send_paging_request(const struct in_addr* dest_ipv4,
                             const struct in6_addr* dest_ipv6) {
  MessageDef* message_p = NULL;
  itti_s11_paging_request_t* paging_request_p = NULL;

  message_p =
      DEPRECATEDitti_alloc_new_message_fatal(TASK_SPGW_APP, S11_PAGING_REQUEST);
  paging_request_p = &message_p->ittiMsg.s11_paging_request;
  memset((void*)paging_request_p, 0, sizeof(itti_s11_paging_request_t));

  if (dest_ipv6) {
    char ip6_str[INET6_ADDRSTRLEN];
    inet_ntop(AF_INET6, dest_ipv6, ip6_str, INET6_ADDRSTRLEN);
    OAILOG_DEBUG(LOG_SPGW_APP, "Paging procedure initiated for ue_ipv6: %s\n",
                 ip6_str);
    // Copy ipv6 address
    memset(paging_request_p->address.ipv6_addr.sin6_addr.s6_addr, 0,
           sizeof(paging_request_p->address.ipv6_addr.sin6_addr.s6_addr));
    memcpy(&paging_request_p->address.ipv6_addr.sin6_addr, dest_ipv6,
           sizeof(struct in6_addr));
    paging_request_p->ip_addr_type = IPV6_ADDR_TYPE;
  } else if (dest_ipv4) {
    char ip4_str[INET_ADDRSTRLEN];
    inet_ntop(AF_INET, dest_ipv4, ip4_str, sizeof(ip4_str));
    OAILOG_DEBUG(LOG_SPGW_APP, "Paging procedure initiated for ue_ipv4: %s\n",
                 ip4_str);
    memcpy(&paging_request_p->address.ipv4_addr.sin_addr, dest_ipv4,
           sizeof(const struct in_addr));
    paging_request_p->ip_addr_type = IPV4_ADDR_TYPE;
  } else {
    OAILOG_ERROR(LOG_SPGW_APP, "Both ipv4 and ipv6 addresses are NULL\n");
    return;
  }
  send_msg_to_task(&spgw_app_task_zmq_ctx, TASK_MME_APP, message_p);
  return;
}

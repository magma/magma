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

#include "pgw_ue_ip_address_alloc.h"

#include "log.h"
#include "MobilityClientAPI.h"
#include "service303.h"

struct in_addr;

int release_ue_ipv4_address(
    const char* imsi, const char* apn, struct in_addr* addr) {
  increment_counter(
      "ue_pdn_connection", 1, 2, "pdn_type", "ipv4", "result",
      "ip_address_released");
  // Release IP address back to PGW IP Address allocator
  return release_ipv4_address(imsi, apn, addr);
}

int release_ue_ipv6_address(
    const char* imsi, const char* apn, struct in6_addr* addr) {
  increment_counter(
      "ue_pdn_connection", 1, 2, "pdn_type", "ipv6", "result",
      "ip_address_released");
  // Release IP address back to PGW IP Address allocator
  return release_ipv6_address(imsi, apn, addr);
}

int get_ip_block(struct in_addr* netaddr, uint32_t* netmask) {
  int rv;

  rv = get_assigned_ipv4_block(0, netaddr, netmask);
  if (rv != 0) {
    OAILOG_CRITICAL(
        LOG_GTPV1U, "ERROR in getting assigned IP block from mobilityd\n");
    return -1;
  }
  return rv;
}

/*
 * Licensed to the OpenAirInterface (OAI) Software Alliance under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The OpenAirInterface Software Alliance licenses this file to You under
 * the Apache License, Version 2.0  (the "License"); you may not use this file
 * except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
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
#include "RpcClient.h"
#include "service303.h"

struct in_addr;

int allocate_ue_ipv4_address(const char *imsi, const char *apn,
                             struct in_addr *addr)
{
  // Call PGW IP Address allocator
  int ip_alloc_status = RPC_STATUS_OK;
  ip_alloc_status = allocate_ipv4_address(imsi, apn, addr);
  if (ip_alloc_status == RPC_STATUS_ALREADY_EXISTS) {
    increment_counter(
      "ue_pdn_connection",
      1,
      2,
      "pdn_type",
      "ipv4",
      "result",
      "ip_address_already_allocated");
    /*
     * This implies that UE session was not release properly.
     * Release the IP address so that subsequent attempt is successfull
     */
    release_ipv4_address(imsi, apn, addr);
    // TODO - Release the GTP-tunnel corresponding to this IP address
  }

  if (ip_alloc_status != RPC_STATUS_OK) {
    OAILOG_ERROR(
      LOG_SPGW_APP,
      "Failed to allocate IPv4 PAA for PDN type IPv4. IP alloc status = %d "
      "imsi<%s> apn<%s>\n", ip_alloc_status, imsi, apn);
  }
  return ip_alloc_status;
}

int release_ue_ipv4_address(const char *imsi, const char *apn,
                            struct in_addr *addr)
{
  increment_counter(
    "ue_pdn_connection",
    1,
    2,
    "pdn_type",
    "ipv4",
    "result",
    "ip_address_released");
  // Release IP address back to PGW IP Address allocator
  return release_ipv4_address(imsi, apn, addr);
}

int get_ip_block(struct in_addr *netaddr, uint32_t *netmask)
{
  int rv;

  rv = get_assigned_ipv4_block(0, netaddr, netmask);
  if (rv != 0) {
    OAILOG_CRITICAL(
      LOG_GTPV1U, "ERROR in getting assigned IP block from mobilityd\n");
    return -1;
  }
  return rv;
}

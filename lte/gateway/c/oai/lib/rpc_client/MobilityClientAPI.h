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

#ifndef RPC_CLIENT_H
#define RPC_CLIENT_H

#include "stdint.h"

#ifdef __cplusplus
extern "C" {
#endif
#include "log.h"
#include "ip_forward_messages_types.h"
#include "intertask_interface.h"
#include "spgw_state.h"

// Status codes from gRPC
#define RPC_STATUS_OK 0
#define RPC_STATUS_CANCELLED 1
#define RPC_STATUS_UNKNOWN 2
#define RPC_STATUS_INVALID_ARGUMENT 3
#define RPC_STATUS_DEADLINE_EXCEEDED 4
#define RPC_STATUS_ALREADY_EXISTS 6
#define RPC_STATUS_PERMISSION_DENIED 7
#define RPC_STATUS_UNAUTHENTICATED 16
#define RPC_STATUS_RESOURCE_EXHAUSTED 8
#define RPC_STATUS_FAILED_PRECONDITION 9
#define RPC_STATUS_ABORTED 10
#define RPC_STATUS_OUT_OF_RANGE 11
#define RPC_STATUS_UNIMPLEMENTED 12
#define RPC_STATUS_INTERNAL 13
#define RPC_STATUS_UNAVAILABLE 14
#define RPC_STATUS_DATA_LOSS 15

/*
 * Get the address and netmask of an assigned IPv4 block
 *
 * @param index (in): index of the IP block requested, currently only ONE
 * IP block (index=0) is supported
 * @param netaddr (out): network address in "network byte order"
 * @param netmask (out): network address mask
 * @return 0 on success
 * @return -RPC_STATUS_INVALID_ARGUMENT if IPBlock is invalid
 * @return -RPC_STATUS_FAILED_PRECONDITION if IPBlock overlaps
 */
int get_assigned_ipv4_block(
  int index,
  struct in_addr* netaddr,
  uint32_t* netmask);

/**
 * Allocate IP address from MobilityServiceClient over gRPC (non-blocking),
 * and handle response for PGW handler.
 * @param subscriber_id: subscriber id string, i.e. IMSI
 * @param apn: access point name string, e.g., "ims", "internet", etc.
 * @param addr: contains the IP address allocated upon returning in
 * "host byte order"
 * @param sgi_create_endpoint_resp itti message for sgi_create_endpoint_resp
 * @param pdn_type str for PDN type (ipv4, ipv6...)
 * @param context_teid tunnel id
 * @param eps_bearer_id bearer id
 * @param spgw_state spgw_state_t struct
 * @param new_bearer_ctxt_info_p SPGW ue context struct
 * @param s5_response itti message for s5_create_session response
 * @return status of gRPC call
 */
int pgw_handle_allocate_ipv4_address(
  const char* subscriber_id,
  const char* apn,
  struct in_addr *addr,
  itti_sgi_create_end_point_response_t sgi_create_endpoint_resp,
  const char* pdn_type,
  teid_t context_teid,
  ebi_t eps_bearer_id,
  spgw_state_t* spgw_state,
  s_plus_p_gw_eps_bearer_context_information_t* new_bearer_ctxt_info_p,
  s5_create_session_response_t s5_response);

/**
 * Allocate IP address from MobilityServiceClient over gRPC (non-blocking),
 * and handle response for PGW handler.
 * @param subscriber_id: subscriber id string, i.e. IMSI
 * @param apn: access point name string, e.g., "ims", "internet", etc.
 * @param addr: contains the IP address allocated upon returning in
 * "host byte order"
 * @param sgi_create_endpoint_resp itti message for sgi_create_endpoint_resp
 * @param pdn_type str for PDN type (ipv6...)
 * @param context_teid tunnel id
 * @param eps_bearer_id bearer id
 * @param spgw_state spgw_state_t struct
 * @param new_bearer_ctxt_info_p SPGW ue context struct
 * @param s5_response itti message for s5_create_session response
 * @return status of gRPC call
 */

int pgw_handle_allocate_ipv6_address(
  const char* subscriber_id,
  const char* apn,
  itti_sgi_create_end_point_response_t sgi_create_endpoint_resp,
  const char* pdn_type,
  teid_t context_teid,
  ebi_t eps_bearer_id,
  spgw_state_t* spgw_state,
  s_plus_p_gw_eps_bearer_context_information_t* new_bearer_ctxt_info_p,
  s5_create_session_response_t s5_response); 

/**
* Allocate IP address from MobilityServiceClient over gRPC (non-blocking),
* and handle response for SGW handler.
* @param subscriber_id: subscriber id string, i.e. IMSI
* @param apn: access point name string, e.g., "ims", "internet", etc.
* @param addr: contains the IP address allocated upon returning in
* "host byte order"
* @param sgi_create_endpoint_resp itti message for sgi_create_endpoint_resp
* @param pdn_type str for PDN type (ipv4, ipv6...)
* @param spgw_state spgw_state_t struct
* @param new_bearer_ctxt_info_p SPGW ue context struct
 * @return status of gRPC call
*/
int sgw_handle_allocate_ipv4_address(
  const char* subscriber_id,
  const char* apn,
  struct in_addr* addr,
  itti_sgi_create_end_point_response_t sgi_create_endpoint_resp,
  const char* pdn_type,
  spgw_state_t* spgw_state,
  s_plus_p_gw_eps_bearer_context_information_t* new_bearer_ctxt_info_p);

/*
 * Release an allocated IP address.
 *
 * The released IP address is put into a tombstone state, and recycled
 * periodically.
 *
 * @param subscriber_id: subscriber id string, i.e. IMSI
 * @param addr: IP address to release in "host byte order"
 * @return 0 on success
 * @return -RPC_STATUS_NOT_FOUND if the requested (SID, IP) pair is not found
 */
int release_ipv4_address(
  const char* subscriber_id,
  const char* apn,
  const struct in_addr* addr);

/*
 * Get the allocated IPv4 address for a subscriber
 * @param subscriber_id: IMSI string
 * @param addr (out): contains the allocated IPv4 address for the subscriber
 * @return 0 on success
 * @return -RPC_STATUS_NOT_FOUND if the SID is not found
 */
int get_ipv4_address_for_subscriber(
  const char* subscriber_id,
  const char* apn,
  struct in_addr* addr);

/*
 * Get the subscriber id given its allocated IPv4 address. If the address
 * isn't associated with a subscriber, then it returns an error
 * @param addr: ipv4 address of subscriber
 * @param subscriber_id (out): contains the imsi of the associated subscriber if
 *                             it exists
 * @return 0 on success
 * @return -RPC_STATUS_NOT_FOUND if IPv4 address is not found
 */
int get_subscriber_id_from_ipv4(
  const struct in_addr* addr,
  char** subscriber_id);

#ifdef __cplusplus
}
#endif

#endif // RPC_CLIENT_H

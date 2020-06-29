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

#include "MobilityClientAPI.h"

#include <grpcpp/security/credentials.h>
#include <cstdint>
#include <cstring>
#include <string>

#include "conversions.h"
#include "common_defs.h"
#include "pcef_handlers.h"
#include "service303.h"
#include "spgw_types.h"
#include "MobilityServiceClient.h"

using grpc::Channel;
using grpc::Status;
using grpc::ChannelCredentials;
using grpc::CreateChannel;
using grpc::InsecureChannelCredentials;
using magma::lte::IPAddress;
using magma::lte::MobilityServiceClient;

static itti_sgi_create_end_point_response_t handle_allocate_ipv4_address_status(
  const grpc::Status& status,
  struct in_addr inaddr,
  const char* imsi,
  const char* apn,
  const char* pdn_type,
  itti_sgi_create_end_point_response_t sgi_create_endpoint_resp);

static itti_sgi_create_end_point_response_t handle_allocate_ipv6_address_status(
    struct in6_addr addr, const char* imsi, const char* apn,
    const char* pdn_type,
    itti_sgi_create_end_point_response_t sgi_create_endpoint_resp,
    uint8_t ipv6_prefix_len);

static itti_sgi_create_end_point_response_t
handle_allocate_ipv4v6_address_status(
    const grpc::Status& status, struct in_addr ip4_addr,
    struct in6_addr ip6_addr, const char* imsi, const char* apn,
    const char* pdn_type,
    itti_sgi_create_end_point_response_t sgi_create_endpoint_resp,
    uint8_t ipv6_prefix_len);

int get_assigned_ipv4_block(
  int index,
  struct in_addr* netaddr,
  uint32_t *netmask)
{
  int status = MobilityServiceClient::getInstance().GetAssignedIPv4Block(
    index, netaddr, netmask);
  return status;
}

int pgw_handle_allocate_ipv4_address(
    const char* subscriber_id, const char* apn, struct in_addr* addr,
    itti_sgi_create_end_point_response_t sgi_create_endpoint_resp,
    const char* pdn_type, spgw_state_t* spgw_state,
    s_plus_p_gw_eps_bearer_context_information_t* new_bearer_ctxt_info_p,
    s5_create_session_response_t s5_response) {
  MobilityServiceClient::getInstance().AllocateIPv4AddressAsync(
      subscriber_id, apn,
      [=, &s5_response](const Status& status, IPAddress ip_msg) {
        memcpy(addr, ip_msg.mutable_address()->c_str(), sizeof(in_addr));

        auto sgi_resp = handle_allocate_ipv4_address_status(
            status, *addr, subscriber_id, apn, pdn_type,
            sgi_create_endpoint_resp);

        if (sgi_resp.status == SGI_STATUS_OK) {
          // create session in PCEF and return
          s5_create_session_request_t session_req = {0};
          session_req.context_teid  = sgi_create_endpoint_resp.context_teid;
          session_req.eps_bearer_id = sgi_create_endpoint_resp.eps_bearer_id;
          char ip_str[INET_ADDRSTRLEN];
          inet_ntop(AF_INET, &(addr->s_addr), ip_str, INET_ADDRSTRLEN);
          struct pcef_create_session_data session_data;
          get_session_req_data(
              spgw_state,
              &new_bearer_ctxt_info_p->sgw_eps_bearer_context_information
                   .saved_message,
              &session_data);
          pcef_create_session(
              spgw_state, subscriber_id, ip_str, NULL, &session_data, sgi_resp,
              session_req, new_bearer_ctxt_info_p);
          OAILOG_FUNC_OUT(LOG_PGW_APP);
        }

        s5_response.eps_bearer_id = sgi_create_endpoint_resp.eps_bearer_id;
        s5_response.context_teid  = sgi_create_endpoint_resp.context_teid;
        handle_s5_create_session_response(
            spgw_state, new_bearer_ctxt_info_p, s5_response);
        OAILOG_FUNC_OUT(LOG_PGW_APP);
      });
  return 0;
}


static itti_sgi_create_end_point_response_t handle_allocate_ipv4_address_status(
  const Status& status,
  struct in_addr inaddr,
  const char* imsi,
  const char* apn,
  const char* pdn_type,
  itti_sgi_create_end_point_response_t sgi_create_endpoint_resp)
{
  if (status.ok()) {
    increment_counter(
      "ue_pdn_connection",
      1,
      2,
      "pdn_type",
      pdn_type,
      "result",
      "success");
    sgi_create_endpoint_resp.paa.ipv4_address = inaddr;
    sgi_create_endpoint_resp.paa.pdn_type = IPv4;
    OAILOG_DEBUG(
      LOG_UTIL,
      "Allocated IPv4 address for imsi <%s>, apn <%s>\n",
      imsi,
      apn);
    sgi_create_endpoint_resp.status = SGI_STATUS_OK;
  } else {
    if (status.error_code() == RPC_STATUS_ALREADY_EXISTS) {
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
      release_ipv4_address(imsi, apn, &inaddr);
      // TODO - Release the GTP-tunnel corresponding to this IP address
      sgi_create_endpoint_resp.status = SGI_STATUS_ERROR_SYSTEM_FAILURE;
    } else {
      increment_counter(
        "ue_pdn_connection",
        1,
        2,
        "pdn_type",
        pdn_type,
        "result",
        "failure");
      OAILOG_ERROR(
        LOG_UTIL,
        "Failed to allocate IPv4 PAA for PDN type IPv4 for "
        "imsi <%s> and apn <%s>\n",
        imsi,
        apn);
      sgi_create_endpoint_resp.status =
        SGI_STATUS_ERROR_ALL_DYNAMIC_ADDRESSES_OCCUPIED;
    }
  }
  return sgi_create_endpoint_resp;
}

int release_ipv4_address(const char *subscriber_id, const char *apn,
                         const struct in_addr* addr)
{
  int status = MobilityServiceClient::getInstance().ReleaseIPv4Address(
    subscriber_id, apn, *addr);
  return status;
}

int get_ipv4_address_for_subscriber(
  const char* subscriber_id,
  const char* apn,
  struct in_addr* addr)
{
  int status = MobilityServiceClient::getInstance().GetIPv4AddressForSubscriber(
    subscriber_id, apn, addr);
  return status;
}

int get_subscriber_id_from_ipv4(
  const struct in_addr* addr,
  char** subscriber_id)
{
  std::string subscriber_id_str;
  int status = MobilityServiceClient::getInstance().GetSubscriberIDFromIPv4(
    *addr, &subscriber_id_str);
  if (!subscriber_id_str.empty()) {
    *subscriber_id = strdup(subscriber_id_str.c_str());
  }
  return status;
}

/* Temporary code to generate IPv6 Network Interface Identifier. To be removed
 * once Mobilityd is enhanced to support IPv6 address
 */
struct in6_addr* generate_random_ip6_interface_id(
    struct in6_addr config_ipv6_prefix)
{
  char str_ip6_addr[INET6_ADDRSTRLEN];
  char *temp_prefix[4], *temp_prefix_free[4];
  char* buf_ipv6            = str_ip6_addr;
  struct in6_addr* ip6_addr = (struct in6_addr*) calloc(1, INET6_ADDRSTRLEN);
  unsigned int random[4]    = {0};
  int itrn                  = 0;

  srand(time(0));
  // Fetch IPv6 prefix from the config
  inet_ntop(AF_INET6, &config_ipv6_prefix, buf_ipv6, INET6_ADDRSTRLEN);
  for (itrn = 0; itrn < 4; itrn++) {
    temp_prefix[itrn] = (char*) calloc(1, 4);
    /* Take a copy of temp_prefix to be freed later because strsep function
     * updates the pointer and points right after the token it found
     */
    temp_prefix_free[itrn] = temp_prefix[itrn];
    temp_prefix[itrn]      = strsep(&buf_ipv6, ":");
  }
  // Generate Random Interface Identifier
  for (itrn = 0; itrn < 4; itrn++) {
    random[itrn] = rand() % 0xffff;
  }
  sprintf(
      str_ip6_addr, "%s:%s:%s:%s:%x:%x:%x:%x", temp_prefix[0], temp_prefix[1],
      temp_prefix[2], temp_prefix[3], random[0], random[1], random[2],
      random[3]);

  // Convert the IPv6 address into in6_addr format
  inet_pton(AF_INET6, str_ip6_addr, ip6_addr);

  for (itrn = 0; itrn < 4; itrn++) {
    free_wrapper((void**) &temp_prefix_free[itrn]);
  }
  return ip6_addr;
}

int pgw_handle_allocate_ipv6_address(
    const char* subscriber_id, const char* apn, struct in6_addr* ip6_prefix,
    itti_sgi_create_end_point_response_t sgi_create_endpoint_resp,
    const char* pdn_type, spgw_state_t* spgw_state,
    s_plus_p_gw_eps_bearer_context_information_t* new_bearer_ctxt_info_p,
    s5_create_session_response_t s5_response,
    struct in6_addr config_ipv6_prefix, uint8_t ipv6_prefix_len)
{
  // TODO Make an RPC call to Mobilityd

  // TODO Temporary code to be removed once Mobilityd supports ipv6
  ip6_prefix    = generate_random_ip6_interface_id(config_ipv6_prefix);
  auto sgi_resp = handle_allocate_ipv6_address_status(
      *ip6_prefix, subscriber_id, apn, pdn_type, sgi_create_endpoint_resp,
      ipv6_prefix_len);

  char ip6_str[INET6_ADDRSTRLEN];
  inet_ntop(AF_INET6, ip6_prefix, ip6_str, INET6_ADDRSTRLEN);

  OAILOG_INFO(
      LOG_UTIL,
      "Allocated IPv6 Address <%s>, PDN Type <%s>,"
      " for IMSI <%s> and APN <%s>\n",
      ip6_str, pdn_type, subscriber_id, apn);

  // create session in PCEF and return
  s5_create_session_request_t session_req = {0};
  session_req.context_teid  = sgi_create_endpoint_resp.context_teid;
  session_req.eps_bearer_id = sgi_create_endpoint_resp.eps_bearer_id;
  struct pcef_create_session_data session_data;
  get_session_req_data(
      spgw_state,
      &new_bearer_ctxt_info_p->sgw_eps_bearer_context_information.saved_message,
      &session_data);
  pcef_create_session(
      spgw_state, subscriber_id, NULL, ip6_str, &session_data, sgi_resp,
      session_req, new_bearer_ctxt_info_p);
  increment_counter("spgw_create_session", 1, 1, "result", "success");
  free_wrapper((void**) &ip6_prefix);
  OAILOG_FUNC_RETURN(LOG_SPGW_APP, RETURNok);
}

static itti_sgi_create_end_point_response_t handle_allocate_ipv6_address_status(
    struct in6_addr addr, const char* imsi, const char* apn,
    const char* pdn_type,
    itti_sgi_create_end_point_response_t sgi_create_endpoint_resp,
    uint8_t ipv6_prefix_len)
{
  increment_counter(
      "ue_pdn_connection", 1, 2, "pdn_type", pdn_type, "result", "success");
  sgi_create_endpoint_resp.paa.ipv6_address       = addr;
  sgi_create_endpoint_resp.paa.ipv6_prefix_length = ipv6_prefix_len;
  sgi_create_endpoint_resp.paa.pdn_type           = IPv6;
  sgi_create_endpoint_resp.status                 = SGI_STATUS_OK;
  return sgi_create_endpoint_resp;
}

int pgw_handle_allocate_ipv4v6_address(
    const char* subscriber_id, const char* apn, struct in_addr* ip4_addr,
    struct in6_addr* ip6_prefix,
    itti_sgi_create_end_point_response_t sgi_create_endpoint_resp,
    const char* pdn_type, spgw_state_t* spgw_state,
    s_plus_p_gw_eps_bearer_context_information_t* new_bearer_ctxt_info_p,
    s5_create_session_response_t s5_response,
    struct in6_addr config_ipv6_prefix, uint8_t ipv6_prefix_len)
{
  // Get IPv4 address
  MobilityServiceClient::getInstance().AllocateIPv4AddressAsync(
      subscriber_id, apn,
      [=, &s5_response](const Status& status, IPAddress ip_msg) {
        memcpy(ip4_addr, ip_msg.mutable_address()->c_str(), sizeof(in_addr));

        // TODO Make an RPC call to Mobilityd to get IPv6 address

        // TODO Remove the below temporary code once Mobilityd supports ipv6
        struct in6_addr* ip6_prefix_temp =
            generate_random_ip6_interface_id(config_ipv6_prefix);
        memcpy(ip6_prefix, ip6_prefix_temp, sizeof(struct in6_addr));

        auto sgi_resp = handle_allocate_ipv4v6_address_status(
            status, *ip4_addr, *ip6_prefix, subscriber_id, apn, pdn_type,
            sgi_create_endpoint_resp, ipv6_prefix_len);

        // create session in PCEF and return
        if (sgi_resp.status == SGI_STATUS_OK) {
          s5_create_session_request_t session_req = {0};
          session_req.context_teid  = sgi_create_endpoint_resp.context_teid;
          session_req.eps_bearer_id = sgi_create_endpoint_resp.eps_bearer_id;
          char ip4_str[INET_ADDRSTRLEN];
          inet_ntop(AF_INET, &(ip4_addr->s_addr), ip4_str, INET_ADDRSTRLEN);
          char ip6_str[INET6_ADDRSTRLEN];
          inet_ntop(AF_INET6, ip6_prefix, ip6_str, INET6_ADDRSTRLEN);
          OAILOG_INFO(
              LOG_UTIL,
              "Allocated IPv4 Address <%s>, IPv6 Address <%s>, PDN Type <%s>,"
              " for IMSI <%s> and APN <%s>\n",
              ip4_str, ip6_str, pdn_type, subscriber_id, apn);

          struct pcef_create_session_data session_data;
          get_session_req_data(
              spgw_state,
              &new_bearer_ctxt_info_p->sgw_eps_bearer_context_information
                   .saved_message,
              &session_data);
          pcef_create_session(
              spgw_state, subscriber_id, ip4_str, ip6_str, &session_data,
              sgi_resp, session_req, new_bearer_ctxt_info_p);
          s5_response.eps_bearer_id = sgi_create_endpoint_resp.eps_bearer_id;
          s5_response.context_teid  = sgi_create_endpoint_resp.context_teid;
          free_wrapper((void**) &ip6_prefix_temp);
          OAILOG_FUNC_OUT(LOG_PGW_APP);
        }
        // If status != SGI_STATUS_OK
        s5_response.eps_bearer_id = sgi_create_endpoint_resp.eps_bearer_id;
        s5_response.context_teid  = sgi_create_endpoint_resp.context_teid;
        handle_s5_create_session_response(
            spgw_state, new_bearer_ctxt_info_p, s5_response);
        OAILOG_FUNC_OUT(LOG_PGW_APP);
      });
  OAILOG_FUNC_RETURN(LOG_SPGW_APP, RETURNok);
}

static itti_sgi_create_end_point_response_t
handle_allocate_ipv4v6_address_status(
    const Status& status, struct in_addr ip4_addr, struct in6_addr ip6_addr,
    const char* imsi, const char* apn, const char* pdn_type,
    itti_sgi_create_end_point_response_t sgi_create_endpoint_resp,
    uint8_t ipv6_prefix_len)
{
  if (status.ok()) {
    increment_counter(
        "ue_pdn_connection", 1, 2, "pdn_type", pdn_type, "result", "success");
    sgi_create_endpoint_resp.paa.ipv4_address = ip4_addr;
    sgi_create_endpoint_resp.status           = SGI_STATUS_OK;
  } else {
    if (status.error_code() == RPC_STATUS_ALREADY_EXISTS) {
      increment_counter(
          "ue_pdn_connection", 1, 2, "pdn_type", "ipv4v6", "result",
          "ip_address_already_allocated");
      /*
       * This implies that UE session was not release properly.
       * Release the IP address so that subsequent attempt is successfull
       */
      release_ipv4_address(imsi, apn, &ip4_addr);
    } else {
      increment_counter(
          "ue_pdn_connection", 1, 2, "pdn_type", pdn_type, "result", "failure");
      OAILOG_ERROR(
          LOG_UTIL,
          "Failed to allocate IPv4 PAA for PDN type IPv4v6 for "
          "imsi <%s> and apn <%s>\n",
          imsi, apn);
    }
  }
  // Assign IPv6 address
  sgi_create_endpoint_resp.paa.ipv6_address       = ip6_addr;
  sgi_create_endpoint_resp.paa.ipv6_prefix_length = ipv6_prefix_len;
  // Set PDN Type
  if (sgi_create_endpoint_resp.status == SGI_STATUS_OK) {
    sgi_create_endpoint_resp.paa.pdn_type = IPv4_AND_v6;
  } else {
    sgi_create_endpoint_resp.paa.pdn_type = IPv6;
  }
  return sgi_create_endpoint_resp;
}

int release_ipv6_address(
    const char* subscriber_id, const char* apn, const struct in6_addr* addr)
{
  int status = 0;
  // TODO- Uncomment once IPv6 is implemented at Mobilityd
  /*status = MobilityServiceClient::getInstance().ReleaseIPv6Address(
    subscriber_id, apn, *addr);*/
  return status;
}

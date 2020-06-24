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
  const char* subscriber_id,
  const char* apn,
  struct in_addr* addr,
  itti_sgi_create_end_point_response_t sgi_create_endpoint_resp,
  const char* pdn_type,
  teid_t context_teid,
  ebi_t eps_bearer_id,
  spgw_state_t* spgw_state,
  s_plus_p_gw_eps_bearer_context_information_t* new_bearer_ctxt_info_p,
  s5_create_session_response_t s5_response)
{
  MobilityServiceClient::getInstance().AllocateIPv4AddressAsync(
    subscriber_id,
    apn,
    [=, &s5_response](const Status& status, IPAddress ip_msg) {
      memcpy(addr, ip_msg.mutable_address()->c_str(), sizeof(in_addr));

      auto sgi_resp = handle_allocate_ipv4_address_status(
        status, *addr, subscriber_id, apn, pdn_type, sgi_create_endpoint_resp);

      if (sgi_resp.status == SGI_STATUS_OK) {
        // create session in PCEF and return
        s5_create_session_request_t session_req = {0};
        session_req.context_teid = context_teid;
        session_req.eps_bearer_id = eps_bearer_id;
        char ip_str[INET_ADDRSTRLEN];
        inet_ntop(AF_INET, &(addr->s_addr), ip_str, INET_ADDRSTRLEN);
        struct pcef_create_session_data session_data;
        get_session_req_data(
          spgw_state,
          &new_bearer_ctxt_info_p->sgw_eps_bearer_context_information.saved_message,
          &session_data);
        pcef_create_session(
          spgw_state,
          subscriber_id,
          ip_str,
          &session_data,
          sgi_resp,
          session_req,
          new_bearer_ctxt_info_p);
        OAILOG_FUNC_OUT(LOG_PGW_APP);
      }

      s5_response.eps_bearer_id = eps_bearer_id;
      s5_response.context_teid = context_teid;
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

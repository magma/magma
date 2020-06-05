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

// TODO: MobilityService IP:port config (t14002037)
#define MOBILITYD_ENDPOINT "localhost:60051"

using grpc::Channel;
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
  const char* addr,
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

int sgw_send_s11_create_session_response(
  const s_plus_p_gw_eps_bearer_context_information_t* new_bearer_ctxt_info_p,
  const gtpv2c_cause_value_t cause,
  imsi64_t imsi64)
{
  int rv = RETURNok;
  MessageDef* message_p = nullptr;
  message_p =
    itti_alloc_new_message(TASK_SPGW_APP, S11_CREATE_SESSION_RESPONSE);
  if (!message_p) {
    OAILOG_ERROR(
      LOG_SPGW_APP, "Message Create Session Response allocation failed\n");
    return RETURNerror;
  }
  itti_s11_create_session_response_t* create_session_response_p =
    &message_p->ittiMsg.s11_create_session_response;
  create_session_response_p->cause.cause_value = cause;
  create_session_response_p->bearer_contexts_created.bearer_contexts[0]
    .cause.cause_value = cause;
  create_session_response_p->bearer_contexts_created.num_bearer_context += 1;
  if (new_bearer_ctxt_info_p) {
    create_session_response_p->teid =
      new_bearer_ctxt_info_p->sgw_eps_bearer_context_information
        .mme_teid_S11;
    create_session_response_p->trxn =
      new_bearer_ctxt_info_p->sgw_eps_bearer_context_information.trxn;
  }
  message_p->ittiMsgHeader.imsi = imsi64;
  rv = itti_send_msg_to_task(TASK_MME_APP, INSTANCE_DEFAULT, message_p);
  return rv;
}

int sgw_handle_allocate_ipv4_address(
  const char* subscriber_id,
  const char* apn,
  struct in_addr* addr,
  itti_sgi_create_end_point_response_t sgi_create_endpoint_resp,
  const char* pdn_type,
  spgw_state_t* spgw_state,
  s_plus_p_gw_eps_bearer_context_information_t* new_bearer_ctxt_info_p)
{
  MobilityServiceClient::getInstance().AllocateIPv4AddressAsync(
    subscriber_id,
    apn,
    [=, &sgi_create_endpoint_resp](const Status& status, IPAddress ip_msg) {
      gtpv2c_cause_value_t cause = REQUEST_ACCEPTED;
      imsi64_t imsi64 =
        new_bearer_ctxt_info_p->sgw_eps_bearer_context_information.imsi64;

      memcpy(addr, ip_msg.mutable_address()->c_str(), sizeof(in_addr));
      auto sgi_resp = handle_allocate_ipv4_address_status(
        status, *addr, subscriber_id, apn, pdn_type, sgi_create_endpoint_resp);

      switch (sgi_create_endpoint_resp.status) {
        case SGI_STATUS_OK:
          // Send Create Session Response with ack
          sgw_handle_sgi_endpoint_created(spgw_state, &sgi_resp, imsi64);
          increment_counter("spgw_create_session", 1, 1, "result", "success");
          OAILOG_FUNC_OUT(LOG_SPGW_APP);
          break;

        case SGI_STATUS_ERROR_ALL_DYNAMIC_ADDRESSES_OCCUPIED:
          increment_counter(
            "spgw_create_session",
            1,
            1,
            "result",
            "failure",
            "cause",
            "resource_not_available");
          cause = ALL_DYNAMIC_ADDRESSES_ARE_OCCUPIED;
          break;

        default:
          cause = REQUEST_REJECTED; // Unspecified reason
          break;
      }

      // Send Create Session Response with Nack
      sgw_send_s11_create_session_response(new_bearer_ctxt_info_p, cause, imsi64);
      OAILOG_FUNC_OUT(LOG_SPGW_APP);
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
int pgw_handle_allocate_ipv6_address(
  const char* subscriber_id,
  const char* apn,
  itti_sgi_create_end_point_response_t sgi_create_endpoint_resp,
  const char* pdn_type,
  teid_t context_teid,
  ebi_t eps_bearer_id,
  spgw_state_t* spgw_state,
  s_plus_p_gw_eps_bearer_context_information_t* new_bearer_ctxt_info_p,
  s5_create_session_response_t s5_response)
{
  // TODO Pruthvi Make an RPC call to Mobilityd
  char addr_temp[INET6_ADDRSTRLEN] = "::1234:5678:abcd:e123";
  auto sgi_resp = handle_allocate_ipv6_address_status(
	addr_temp, subscriber_id, apn, pdn_type, sgi_create_endpoint_resp);

  // create session in PCEF and return
  s5_create_session_request_t session_req = {0};
  session_req.context_teid = context_teid;
  session_req.eps_bearer_id = eps_bearer_id;
  struct pcef_create_session_data session_data;
  get_session_req_data(
    spgw_state,
    &new_bearer_ctxt_info_p->sgw_eps_bearer_context_information.saved_message,
    &session_data);
  pcef_create_session(
    spgw_state,
    subscriber_id,
    addr_temp,
    &session_data,
    sgi_resp,
    session_req,
    new_bearer_ctxt_info_p);
  s5_response.eps_bearer_id = eps_bearer_id;
  s5_response.context_teid = context_teid;

  return 0;
}
static itti_sgi_create_end_point_response_t handle_allocate_ipv6_address_status(
  const char* addr,
  const char* imsi,
  const char* apn,
  const char* pdn_type,
  itti_sgi_create_end_point_response_t sgi_create_endpoint_resp)
{
    increment_counter(
      "ue_pdn_connection",
      1,
      2,
      "pdn_type",
      pdn_type,
      "result",
      "success");
    unsigned char  temp_addr[sizeof (struct in6_addr)] = {0};
    if(inet_pton(AF_INET6, addr, temp_addr)== 1) {
      // Copy only the interface identifier i.e the last 64 bits
      memcpy(&sgi_create_endpoint_resp.paa.ipv6_address,&temp_addr[8], 8);
      sgi_create_endpoint_resp.paa.ipv6_prefix_length = 64;
      sgi_create_endpoint_resp.paa.pdn_type = IPv6;
      OAILOG_DEBUG(
        LOG_UTIL,
        "Allocated IPv6 Interface identifier for imsi <%s>, apn <%s>\n",
        imsi,
        apn);
      sgi_create_endpoint_resp.status = SGI_STATUS_OK;
    }
  return sgi_create_endpoint_resp;
}


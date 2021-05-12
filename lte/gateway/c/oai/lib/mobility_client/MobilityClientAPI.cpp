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

#include "MobilityClientAPI.h"

#include <grpcpp/security/credentials.h>

#include <cstdint>
#include <cstring>
#include <string>

#include "common_defs.h"
#include "common_types.h"
#include "conversions.h"
#include "intertask_interface.h"
#include "log.h"
#include "MobilityServiceClient.h"
#include "service303.h"
#include "spgw_types.h"

using grpc::Channel;
using grpc::ChannelCredentials;
using grpc::CreateChannel;
using grpc::InsecureChannelCredentials;
using grpc::Status;
using magma::lte::AllocateIPAddressResponse;
using magma::lte::IPAddress;
using magma::lte::MobilityServiceClient;

extern task_zmq_ctx_t grpc_service_task_zmq_ctx;

static void handle_allocate_ipv4_address_status(
    const grpc::Status& status, struct in_addr inaddr, int vlan,
    const char* imsi, const char* apn, const char* pdn_type,
    teid_t context_teid, ebi_t eps_bearer_id);

static void handle_allocate_ipv6_address_status(
    const grpc::Status& status, struct in6_addr addr, int vlan,
    const char* imsi, const char* apn, const char* pdn_type,
    teid_t context_teid, ebi_t eps_bearer_id);

static void handle_allocate_ipv4v6_address_status(
    const grpc::Status& status, struct in_addr ip4_addr,
    struct in6_addr ip6_addr, int vlan, const char* imsi, const char* apn,
    const char* pdn_type, teid_t context_teid, ebi_t eps_bearer_id);

int get_assigned_ipv4_block(
    int index, struct in_addr* netaddr, uint32_t* netmask) {
  int status = MobilityServiceClient::getInstance().GetAssignedIPv4Block(
      index, netaddr, netmask);
  return status;
}

int pgw_handle_allocate_ipv4_address(
    const char* subscriber_id, const char* apn, const char* pdn_type,
    teid_t context_teid, ebi_t eps_bearer_id) {
  auto subscriber_id_str = std::string(subscriber_id);
  auto apn_str           = std::string(apn);
  auto pdn_type_str      = std::string(pdn_type);
  MobilityServiceClient::getInstance().AllocateIPv4AddressAsync(
      subscriber_id_str, apn_str,
      [subscriber_id_str, apn_str, pdn_type_str, context_teid, eps_bearer_id](
          const Status& status, const AllocateIPAddressResponse& ip_msg) {
        struct in_addr addr;
        std::string ipv4_addr_str;
        if (ip_msg.ip_list_size() > 0) {
          ipv4_addr_str = ip_msg.ip_list(0).address();
        }
        memcpy(&addr, ipv4_addr_str.c_str(), sizeof(in_addr));
        int vlan = atoi(ip_msg.vlan().c_str());
        handle_allocate_ipv4_address_status(
            status, addr, vlan, subscriber_id_str.c_str(), apn_str.c_str(),
            pdn_type_str.c_str(), context_teid, eps_bearer_id);
      });
  return RETURNok;
}

static void handle_allocate_ipv4_address_status(
    const Status& status, struct in_addr inaddr, int vlan, const char* imsi,
    const char* apn, const char* pdn_type, teid_t context_teid,
    ebi_t eps_bearer_id) {
  MessageDef* message_p;
  message_p = itti_alloc_new_message(TASK_GRPC_SERVICE, IP_ALLOCATION_RESPONSE);
  if (!message_p) {
    OAILOG_ERROR(
        LOG_UTIL, "Message IP Allocation Response allocation failed\n");
    return;
  }

  itti_ip_allocation_response_t* ip_allocation_response_p;
  ip_allocation_response_p = &message_p->ittiMsg.ip_allocation_response;
  memset(ip_allocation_response_p, 0, sizeof(itti_ip_allocation_response_t));

  ip_allocation_response_p->context_teid  = context_teid;
  ip_allocation_response_p->eps_bearer_id = eps_bearer_id;

  if (status.ok()) {
    increment_counter(
        "ue_pdn_connection", 1, 2, "pdn_type", pdn_type, "result", "success");
    ip_allocation_response_p->paa.ipv4_address = inaddr;
    ip_allocation_response_p->paa.pdn_type     = IPv4;
    ip_allocation_response_p->paa.vlan         = vlan;
    ip_allocation_response_p->status           = SGI_STATUS_OK;

    OAILOG_DEBUG(
        LOG_UTIL, "Allocated IPv4 address for imsi <%s>, apn <%s> vlan %d\n",
        imsi, apn, vlan);
  } else {
    if (status.error_code() == RPC_STATUS_ALREADY_EXISTS) {
      increment_counter(
          "ue_pdn_connection", 1, 2, "pdn_type", "ipv4", "result",
          "ip_address_already_allocated");

      ip_allocation_response_p->status = SGI_STATUS_ERROR_SYSTEM_FAILURE;
    } else {
      increment_counter(
          "ue_pdn_connection", 1, 2, "pdn_type", pdn_type, "result", "failure");
      OAILOG_ERROR(
          LOG_UTIL,
          "Failed to allocate IPv4 PAA for PDN type IPv4 for "
          "imsi <%s> and apn <%s>\n",
          imsi, apn);
      ip_allocation_response_p->status =
          SGI_STATUS_ERROR_ALL_DYNAMIC_ADDRESSES_OCCUPIED;
    }
  }

  IMSI_STRING_TO_IMSI64(imsi, &message_p->ittiMsgHeader.imsi);
  OAILOG_DEBUG_UE(
      LOG_UTIL, message_p->ittiMsgHeader.imsi,
      "Sending IP allocation response message with cause: %u\n",
      ip_allocation_response_p->status);
  send_msg_to_task(&grpc_service_task_zmq_ctx, TASK_SPGW_APP, message_p);
}

int release_ipv4_address(
    const char* subscriber_id, const char* apn, const struct in_addr* addr) {
  int status = MobilityServiceClient::getInstance().ReleaseIPv4Address(
      subscriber_id, apn, *addr);
  return status;
}

int get_ipv4_address_for_subscriber(
    const char* subscriber_id, const char* apn, struct in_addr* addr) {
  int status = MobilityServiceClient::getInstance().GetIPv4AddressForSubscriber(
      subscriber_id, apn, addr);
  return status;
}

int get_subscriber_id_from_ipv4(
    const struct in_addr* addr, char** subscriber_id) {
  std::string subscriber_id_str;
  int status = MobilityServiceClient::getInstance().GetSubscriberIDFromIPv4(
      *addr, &subscriber_id_str);
  if (!subscriber_id_str.empty()) {
    *subscriber_id = strdup(subscriber_id_str.c_str());
  }
  return status;
}

int pgw_handle_allocate_ipv6_address(
    const char* subscriber_id, const char* apn, const char* pdn_type,
    teid_t context_teid, ebi_t eps_bearer_id) {
  auto subscriber_id_str = std::string(subscriber_id);
  auto apn_str           = std::string(apn);
  auto pdn_type_str      = std::string(pdn_type);
  // Make an RPC call to Mobilityd
  MobilityServiceClient::getInstance().AllocateIPv6AddressAsync(
      subscriber_id_str, apn_str,
      [subscriber_id_str, apn_str, pdn_type_str, context_teid, eps_bearer_id](
          const Status& status, const AllocateIPAddressResponse& ip_msg) {
        struct in6_addr ip6_addr;
        std::string ipv6_addr_str;
        if (ip_msg.ip_list_size() > 0) {
          ipv6_addr_str = ip_msg.ip_list(0).address();
        } else {
          OAILOG_ERROR(
              LOG_UTIL,
              " Error in allocating ipv6 address for IMSI <%s> apn <%s>\n",
              subscriber_id_str.c_str(), apn_str.c_str());
        }

        memcpy(&ip6_addr.s6_addr, ipv6_addr_str.c_str(), sizeof(in6_addr));
        int vlan = atoi(ip_msg.vlan().c_str());

        handle_allocate_ipv6_address_status(
            status, ip6_addr, vlan, subscriber_id_str.c_str(), apn_str.c_str(),
            pdn_type_str.c_str(), context_teid, eps_bearer_id);
      });
  return RETURNok;
}

static void handle_allocate_ipv6_address_status(
    const Status& status, struct in6_addr addr, int vlan, const char* imsi,
    const char* apn, const char* pdn_type, teid_t context_teid,
    ebi_t eps_bearer_id) {
  MessageDef* message_p;
  message_p = itti_alloc_new_message(TASK_GRPC_SERVICE, IP_ALLOCATION_RESPONSE);
  if (!message_p) {
    OAILOG_ERROR(
        LOG_UTIL, "Message IP Allocation Response allocation failed\n");
    return;
  }

  itti_ip_allocation_response_t* ip_allocation_response_p;
  ip_allocation_response_p = &message_p->ittiMsg.ip_allocation_response;
  memset(ip_allocation_response_p, 0, sizeof(itti_ip_allocation_response_t));

  ip_allocation_response_p->context_teid  = context_teid;
  ip_allocation_response_p->eps_bearer_id = eps_bearer_id;

  if (status.ok()) {
    increment_counter(
        "ue_pdn_connection", 1, 2, "pdn_type", pdn_type, "result", "success");
    ip_allocation_response_p->paa.ipv6_address       = addr;
    ip_allocation_response_p->paa.ipv6_prefix_length = IPV6_PREFIX_LEN;
    ip_allocation_response_p->paa.pdn_type           = IPv6;
    ip_allocation_response_p->paa.vlan               = vlan;
    ip_allocation_response_p->status                 = SGI_STATUS_OK;

    OAILOG_DEBUG(
        LOG_UTIL, "Allocated IPv6 address for imsi <%s>, apn <%s> vlan %d\n",
        imsi, apn, vlan);
  } else {
    if (status.error_code() == RPC_STATUS_ALREADY_EXISTS) {
      increment_counter(
          "ue_pdn_connection", 1, 2, "pdn_type", pdn_type, "result",
          "ip_address_already_allocated");
      ip_allocation_response_p->status = SGI_STATUS_ERROR_SYSTEM_FAILURE;
    } else {
      increment_counter(
          "ue_pdn_connection", 1, 2, "pdn_type", pdn_type, "result", "failure");
      OAILOG_ERROR(
          LOG_UTIL,
          "Failed to allocate IPv6 PAA for PDN type IPv6 for "
          "imsi <%s> and apn <%s>\n",
          imsi, apn);
      ip_allocation_response_p->status =
          SGI_STATUS_ERROR_ALL_DYNAMIC_ADDRESSES_OCCUPIED;
    }
  }

  IMSI_STRING_TO_IMSI64(imsi, &message_p->ittiMsgHeader.imsi);
  OAILOG_DEBUG_UE(
      LOG_UTIL, message_p->ittiMsgHeader.imsi,
      "Sending IP allocation response message with cause: %u\n",
      ip_allocation_response_p->status);
  send_msg_to_task(&grpc_service_task_zmq_ctx, TASK_SPGW_APP, message_p);
}

int pgw_handle_allocate_ipv4v6_address(
    const char* subscriber_id, const char* apn, const char* pdn_type,
    teid_t context_teid, ebi_t eps_bearer_id) {
  auto subscriber_id_str = std::string(subscriber_id);
  auto apn_str           = std::string(apn);
  auto pdn_type_str      = std::string(pdn_type);
  // Get IPv4v6 address
  MobilityServiceClient::getInstance().AllocateIPv4v6AddressAsync(
      subscriber_id_str, apn_str,
      [subscriber_id_str, apn_str, pdn_type_str, context_teid, eps_bearer_id](
          const Status& status, const AllocateIPAddressResponse& ip_msg) {
        struct in_addr ip4_addr;
        struct in6_addr ip6_addr;
        std::string ipv4_addr_str;
        std::string ipv6_addr_str;
        if (ip_msg.ip_list_size() == 2) {
          ipv4_addr_str = ip_msg.ip_list(0).address();
          ipv6_addr_str = ip_msg.ip_list(1).address();
          OAILOG_INFO(
              LOG_UTIL,
              "Allocated IPv4 Address <%s>, IPv6 Address <%s>, PDN Type <%s>,"
              " for IMSI <%s> and APN <%s>\n",
              ipv4_addr_str.c_str(), ipv6_addr_str.c_str(),
              pdn_type_str.c_str(), subscriber_id_str.c_str(), apn_str.c_str());
        } else {
          OAILOG_ERROR(
              LOG_UTIL,
              " Error in allocating IPv4 and IPv6 addresses for IMSI <%s> apn "
              "<%s>\n",
              subscriber_id_str.c_str(), apn_str.c_str());
        }
        memcpy(&ip4_addr, ipv4_addr_str.c_str(), sizeof(in_addr));
        memcpy(&ip6_addr, ipv6_addr_str.c_str(), sizeof(in6_addr));
        int vlan = atoi(ip_msg.vlan().c_str());
        handle_allocate_ipv4v6_address_status(
            status, ip4_addr, ip6_addr, vlan, subscriber_id_str.c_str(),
            apn_str.c_str(), pdn_type_str.c_str(), context_teid, eps_bearer_id);
      });
  return RETURNok;
}

static void handle_allocate_ipv4v6_address_status(
    const Status& status, struct in_addr ip4_addr, struct in6_addr ip6_addr,
    int vlan, const char* imsi, const char* apn, const char* pdn_type,
    teid_t context_teid, ebi_t eps_bearer_id) {
  MessageDef* message_p;
  message_p = itti_alloc_new_message(TASK_GRPC_SERVICE, IP_ALLOCATION_RESPONSE);
  if (!message_p) {
    OAILOG_ERROR(
        LOG_UTIL, "Message IP Allocation Response allocation failed\n");
    return;
  }

  itti_ip_allocation_response_t* ip_allocation_response_p;
  ip_allocation_response_p = &message_p->ittiMsg.ip_allocation_response;
  memset(ip_allocation_response_p, 0, sizeof(itti_ip_allocation_response_t));

  ip_allocation_response_p->context_teid  = context_teid;
  ip_allocation_response_p->eps_bearer_id = eps_bearer_id;

  if (status.ok()) {
    increment_counter(
        "ue_pdn_connection", 1, 2, "pdn_type", pdn_type, "result", "success");
    ip_allocation_response_p->paa.ipv4_address       = ip4_addr;
    ip_allocation_response_p->paa.ipv6_address       = ip6_addr;
    ip_allocation_response_p->paa.ipv6_prefix_length = IPV6_PREFIX_LEN;
    ip_allocation_response_p->paa.pdn_type           = IPv4_AND_v6;
    ip_allocation_response_p->paa.vlan               = vlan;
    ip_allocation_response_p->status                 = SGI_STATUS_OK;

    OAILOG_DEBUG(
        LOG_UTIL, "Allocated IPv4v6 address for imsi <%s>, apn <%s> vlan %d\n",
        imsi, apn, vlan);
  } else {
    if (status.error_code() == RPC_STATUS_ALREADY_EXISTS) {
      increment_counter(
          "ue_pdn_connection", 1, 2, "pdn_type", pdn_type, "result",
          "ip_address_already_allocated");
      ip_allocation_response_p->status = SGI_STATUS_ERROR_SYSTEM_FAILURE;
    } else {
      increment_counter(
          "ue_pdn_connection", 1, 2, "pdn_type", pdn_type, "result", "failure");
      OAILOG_ERROR(
          LOG_UTIL,
          "Failed to allocate IPv4v6 PAA for PDN type IPv4v6 for "
          "imsi <%s> and apn <%s>\n",
          imsi, apn);
      ip_allocation_response_p->status =
          SGI_STATUS_ERROR_ALL_DYNAMIC_ADDRESSES_OCCUPIED;
    }
  }

  IMSI_STRING_TO_IMSI64(imsi, &message_p->ittiMsgHeader.imsi);
  OAILOG_DEBUG_UE(
      LOG_UTIL, message_p->ittiMsgHeader.imsi,
      "Sending IP allocation response message with cause: %u\n",
      ip_allocation_response_p->status);
  send_msg_to_task(&grpc_service_task_zmq_ctx, TASK_SPGW_APP, message_p);
}

int release_ipv6_address(
    const char* subscriber_id, const char* apn, const struct in6_addr* addr) {
  int status = MobilityServiceClient::getInstance().ReleaseIPv6Address(
      subscriber_id, apn, *addr);
  return status;
}

int release_ipv4v6_address(
    const char* subscriber_id, const char* apn, const struct in_addr* ipv4_addr,
    const struct in6_addr* ipv6_addr) {
  int status = MobilityServiceClient::getInstance().ReleaseIPv4v6Address(
      subscriber_id, apn, *ipv4_addr, *ipv6_addr);
  return status;
}

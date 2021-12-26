/**
 * Copyright 2020 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

#include <grpcpp/security/credentials.h>
#include <cstdint>
#include <cstring>
#include <string>

#ifdef __cplusplus
extern "C" {
#endif
#include "lte/gateway/c/core/oai/lib/itti/intertask_interface.h"
#ifdef __cplusplus
}
#endif

#include "lte/gateway/c/core/oai/common/conversions.h"
#include "lte/gateway/c/core/oai/common/common_defs.h"
#include "lte/gateway/c/core/oai/include/service303.h"
#include "lte/gateway/c/core/oai/common/common_types.h"
#include "lte/gateway/c/core/oai/lib/mobility_client/MobilityServiceClient.h"
#include "lte/gateway/c/core/oai/lib/n11/M5GMobilityServiceClient.h"

using grpc::Status;
using magma::lte::AllocateIPAddressResponse;
using magma::lte::IPAddress;
using magma::lte::MobilityServiceClient;

extern task_zmq_ctx_t grpc_service_task_zmq_ctx;

static void handle_allocate_ipv4_address_status(
    const grpc::Status& status, struct in_addr in_ip4_addr, int vlan,
    const char* imsi, const char* apn, uint32_t pdu_session_id, uint8_t pti,
    uint32_t pdu_session_type, uint32_t gnb_gtp_teid,
    uint8_t* gnb_gtp_teid_ip_addr, uint8_t gnb_gtp_teid_ip_addr_len) {
  MessageDef* message_p;
  message_p =
      itti_alloc_new_message(TASK_GRPC_SERVICE, AMF_IP_ALLOCATION_RESPONSE);

  itti_amf_ip_allocation_response_t* amf_ip_allocation_response_p;
  amf_ip_allocation_response_p = &message_p->ittiMsg.amf_ip_allocation_response;
  memset(
      amf_ip_allocation_response_p, 0,
      sizeof(itti_amf_ip_allocation_response_t));

  memcpy(amf_ip_allocation_response_p->imsi, imsi, IMSI_BCD_DIGITS_MAX);
  amf_ip_allocation_response_p->imsi_length      = IMSI_BCD_DIGITS_MAX;
  amf_ip_allocation_response_p->pdu_session_id   = pdu_session_id;
  amf_ip_allocation_response_p->pti              = pti;
  amf_ip_allocation_response_p->pdu_session_type = pdu_session_type;
  amf_ip_allocation_response_p->paa.ipv4_address = in_ip4_addr;
  amf_ip_allocation_response_p->paa.pdn_type     = IPv4;
  amf_ip_allocation_response_p->paa.vlan         = vlan;

  amf_ip_allocation_response_p->gnb_gtp_teid = gnb_gtp_teid;

  memcpy(
      amf_ip_allocation_response_p->gnb_gtp_teid_ip_addr, gnb_gtp_teid_ip_addr,
      4);

  memcpy(amf_ip_allocation_response_p->apn, apn, strlen(apn) + 1);

  if (status.ok()) {
    amf_ip_allocation_response_p->result = 0;
  } else {
    amf_ip_allocation_response_p->result = -1;
  }

  send_msg_to_task(&grpc_service_task_zmq_ctx, TASK_AMF_APP, message_p);
}

namespace magma5g {

int AsyncM5GMobilityServiceClient::allocate_ipv4_address(
    const char* subscriber_id, const char* apn, uint32_t pdu_session_id,
    uint8_t pti, uint32_t pdu_session_type, uint32_t gnb_gtp_teid,
    uint8_t* gnb_gtp_teid_ip_addr, uint8_t gnb_gtp_teid_ip_addr_len) {
  auto subscriber_id_str = std::string(subscriber_id);
  auto apn_str           = std::string(apn);
  MobilityServiceClient::getInstance().AllocateIPv4AddressAsync(
      subscriber_id_str, apn,
      [subscriber_id_str, apn, pdu_session_id, pti, pdu_session_type,
       gnb_gtp_teid, gnb_gtp_teid_ip_addr, gnb_gtp_teid_ip_addr_len](
          const Status& status, const AllocateIPAddressResponse& ip_msg) {
        struct in_addr addr;
        std::string ipv4_addr_str;

        if (ip_msg.ip_list_size() > 0) {
          ipv4_addr_str = ip_msg.ip_list(0).address();
        }
        memcpy(&addr, ipv4_addr_str.c_str(), sizeof(in_addr));
        int vlan = atoi(ip_msg.vlan().c_str());

        handle_allocate_ipv4_address_status(
            status, addr, vlan, subscriber_id_str.c_str(), apn, pdu_session_id,
            pti, pdu_session_type, gnb_gtp_teid, gnb_gtp_teid_ip_addr,
            gnb_gtp_teid_ip_addr_len);
      });
  return RETURNok;
}

int AsyncM5GMobilityServiceClient::release_ipv4_address(
    const char* subscriber_id, const char* apn, const struct in_addr* addr) {
  MobilityServiceClient::getInstance().ReleaseIPv4Address(
      subscriber_id, apn, *addr);
  return RETURNok;
}

AsyncM5GMobilityServiceClient::AsyncM5GMobilityServiceClient() {}

AsyncM5GMobilityServiceClient& AsyncM5GMobilityServiceClient::getInstance() {
  static AsyncM5GMobilityServiceClient instance;
  return instance;
}

}  // namespace magma5g

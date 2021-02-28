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

#pragma once

#include <arpa/inet.h>
#include <grpc++/grpc++.h>
#include <stdint.h>
#include <functional>
#include <memory>
#include <string>

#include "PipelinedClientAPI.h"
#include "GRPCReceiver.h"
#include "lte/protos/pipelined.grpc.pb.h"

namespace grpc {
class Channel;
class ClientContext;
class Status;
}  // namespace grpc

namespace magma {
namespace orc8r {
class Void;
}  // namespace orc8r
}  // namespace magma

using grpc::Channel;
using grpc::ClientContext;
using grpc::Status;

namespace magma {
namespace lte {

using namespace orc8r;
/*
 * gRPC client for PipelineService
 */
class PipelinedServiceClient : public GRPCReceiver {
 public:
  IPFlowDL set_ue_ip_flow_dl(struct ip_flow_dl flow_dl);

  // APIs for adding tunnels
  static int UpdateUEIPv4SessionSet(
      const struct in_addr& ue_ipv4_addr, int vlan,
      struct in_addr& enb_ipv4_addr, uint32_t in_teid, uint32_t out_teid,
      const std::string& imsi, uint32_t flow_precedence, const std::string& apn,
      uint32_t state,
      std::function<void(grpc::Status, UESessionContextResponse)> callback);

  static int UpdateUEIPv4SessionSetWithFlowdl(
      const struct in_addr& ue_ipv4_addr, int vlan,
      struct in_addr& enb_ipv4_addr, uint32_t in_teid, uint32_t out_teid,
      const std::string& imsi, const struct ip_flow_dl& flow_dl,
      uint32_t flow_precedence, const std::string& apn, uint32_t state,
      std::function<void(Status, UESessionContextResponse)> callback);

  static int UpdateUEIPv4v6SessionSet(
      const struct in_addr& ue_ipv4_addr, struct in6_addr& ue_ipv6_addr,
      int vlan, struct in_addr& enb_ipv4_addr, uint32_t in_teid,
      uint32_t out_teid, const std::string& imsi, uint32_t flow_precedence,
      const std::string& apn, uint32_t state,
      std::function<void(grpc::Status, UESessionContextResponse)> callback);

  static int UpdateUEIPv4v6SessionSetWithFlowdl(
      const struct in_addr& ue_ipv4_addr, struct in6_addr& ue_ipv6_addr,
      int vlan, struct in_addr& enb_ipv4_addr, uint32_t in_teid,
      uint32_t out_teid, const std::string& imsi,
      const struct ip_flow_dl& flow_dl, uint32_t flow_precedence,
      const std::string& apn, uint32_t state,
      std::function<void(Status, UESessionContextResponse)> callback);

  // APIs for deleting the tunnel
  static int UpdateUEIPv4SessionSet(
      struct in_addr& enb_ipv4_addr, const struct in_addr& ue_ipv4_addr,
      uint32_t in_teid, uint32_t out_teid, uint32_t state,
      std::function<void(grpc::Status, UESessionContextResponse)> callback);

  static int UpdateUEIPv4SessionSetWithFlowdl(
      struct in_addr& enb_ipv4_addr, const struct in_addr& ue_ipv4_addr,
      uint32_t in_teid, uint32_t out_teid, const struct ip_flow_dl& flow_dl,
      uint32_t state,
      std::function<void(Status, UESessionContextResponse)> callback);

  static int UpdateUEIPv4v6SessionSet(
      struct in_addr& enb_ipv4_addr, const struct in_addr& ue_ipv4_addr,
      struct in6_addr& ue_ipv6_addr, uint32_t in_teid, uint32_t out_teid,
      uint32_t state,
      std::function<void(grpc::Status, UESessionContextResponse)> callback);

  static int UpdateUEIPv4v6SessionSetWithFlowdl(
      struct in_addr& enb_ipv4_addr, const struct in_addr& ue_ipv4_addr,
      struct in6_addr& ue_ipv6_addr, uint32_t in_teid, uint32_t out_teid,
      const struct ip_flow_dl& flow_dl, uint32_t state,
      std::function<void(Status, UESessionContextResponse)> callback);

  // APIs for discarding data on tunnels
  static int UpdateUEIPv4SessionSet(
      const struct in_addr& ue_ipv4_addr, uint32_t in_teid, uint32_t state,
      std::function<void(grpc::Status, UESessionContextResponse)> callback);

  static int UpdateUEIPv4SessionSetWithFlowdl(
      const struct in_addr& ue_ipv4_addr, uint32_t in_teid,
      const struct ip_flow_dl& flow_dl, uint32_t state,
      std::function<void(Status, UESessionContextResponse)> callback);

  static int UpdateUEIPv4v6SessionSet(
      const struct in_addr& ue_ipv4_addr, struct in6_addr& ue_ipv6_addr,
      uint32_t in_teid, uint32_t state,
      std::function<void(grpc::Status, UESessionContextResponse)> callback);

  static int UpdateUEIPv4v6SessionSetWithFlowdl(
      const struct in_addr& ue_ipv4_addr, struct in6_addr& ue_ipv6_addr,
      uint32_t in_teid, const struct ip_flow_dl& flow_dl, uint32_t state,
      std::function<void(Status, UESessionContextResponse)> callback);

  // APIs for forwarding data on the tunnel
  static int UpdateUEIPv4SessionSet(
      const struct in_addr& ue_ipv4_addr, uint32_t in_teid,
      uint32_t flow_precedence, uint32_t state,
      std::function<void(grpc::Status, UESessionContextResponse)> callback);

  static int UpdateUEIPv4SessionSetWithFlowdl(
      const struct in_addr& ue_ipv4_addr, uint32_t in_teid,
      const struct ip_flow_dl& flow_dl, uint32_t flow_precedence,
      uint32_t state,
      std::function<void(Status, UESessionContextResponse)> callback);

  static int UpdateUEIPv4v6SessionSet(
      const struct in_addr& ue_ipv4_addr, struct in6_addr& ue_ipv6_addr,
      uint32_t in_teid, uint32_t flow_precedence, uint32_t state,
      std::function<void(grpc::Status, UESessionContextResponse)> callback);

  static int UpdateUEIPv4v6SessionSetWithFlowdl(
      const struct in_addr& ue_ipv4_addr, struct in6_addr& ue_ipv6_addr,
      uint32_t in_teid, const struct ip_flow_dl& flow_dl,
      uint32_t flow_precedence, uint32_t state,
      std::function<void(Status, UESessionContextResponse)> callback);

  // APIs for paging IDLE -> PAGING, ACTIVATE -> Delete PAGING
  static int UpdateUEIPv4SessionSet(
      const struct in_addr& ue_ipv4_addr, uint32_t state,
      std::function<void(grpc::Status, UESessionContextResponse)> callback);

 public:
  static PipelinedServiceClient& get_instance();
  PipelinedServiceClient(PipelinedServiceClient const&) = delete;
  void operator=(PipelinedServiceClient const&) = delete;

  // Set the UE IPv4 address
  static void ue_set_ipv4_addr(
      const struct in_addr& ue_ipv4_addr, UESessionSet& request) {
    IPAddress* encode_ue_ipv4_addr = request.mutable_ue_ipv4_address();
    encode_ue_ipv4_addr->set_version(IPAddress::IPV4);
    encode_ue_ipv4_addr->set_address(&ue_ipv4_addr, sizeof(struct in_addr));
  }

  // Set the UE IPv6 address
  static void ue_set_ipv6_addr(
      const struct in6_addr& ue_ipv6_addr, UESessionSet& request) {
    IPAddress* encode_ue_ipv6_addr = request.mutable_ue_ipv6_address();
    encode_ue_ipv6_addr->set_version(IPAddress::IPV6);
    encode_ue_ipv6_addr->set_address(&ue_ipv6_addr, sizeof(struct in6_addr));
  }

  // Set the GNB IPv4 address
  static void gnb_set_ipv4_addr(
      const struct in_addr& gnb_ipv4_addr, UESessionSet& request) {
    IPAddress* encode_gnb_ipv4_addr = request.mutable_enb_ip_address();
    encode_gnb_ipv4_addr->set_version(IPAddress::IPV4);
    encode_gnb_ipv4_addr->set_address(&gnb_ipv4_addr, sizeof(struct in_addr));
  }

  // Set the Session Config State
  static void config_ue_session_state(
      uint32_t& ue_state, UESessionSet& request) {
    UESessionState* ue_session_state = request.mutable_ue_session_state();

    if (UE_SESSION_ACTIVE_STATE == ue_state) {
      ue_session_state->set_ue_config_state(UESessionState::ACTIVE);
    } else if (UE_SESSION_UNREGISTERED_STATE == ue_state) {
      ue_session_state->set_ue_config_state(UESessionState::UNREGISTERED);
    } else if (UE_SESSION_INSTALL_IDLE_STATE == ue_state) {
      ue_session_state->set_ue_config_state(UESessionState::INSTALL_IDLE);
    } else if (UE_SESSION_UNINSTALL_IDLE_STATE == ue_state) {
      ue_session_state->set_ue_config_state(UESessionState::UNINSTALL_IDLE);
    } else if (UE_SESSION_SUSPENDED_DATA_STATE == ue_state) {
      ue_session_state->set_ue_config_state(UESessionState::SUSPENDED_DATA);
    } else if (UE_SESSION_RESUME_DATA_STATE == ue_state) {
      ue_session_state->set_ue_config_state(UESessionState::RESUME_DATA);
    }
  }

 private:
  PipelinedServiceClient();
  std::unique_ptr<Pipelined::Stub> stub_{};
  static const uint32_t RESPONSE_TIMEOUT = 3;  // seconds
};

}  // namespace lte
}  // namespace magma

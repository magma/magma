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
#include "includes/GRPCReceiver.h"
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

 private:
  PipelinedServiceClient();
  std::unique_ptr<Pipelined::Stub> stub_{};
  static const uint32_t RESPONSE_TIMEOUT = 3;  // seconds
};

}  // namespace lte
}  // namespace magma

/*
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

#include <grpc++/grpc++.h>
#include <grpcpp/impl/codegen/status.h>

#include "lte/protos/session_manager.grpc.pb.h"
#include "lte/protos/policydb.pb.h"

namespace grpc {
class ServerContext;
}  // namespace grpc
namespace magma {
namespace lte {
class SetSMSessionContextAccess;
class SetSmNotificationContext;
class SmContextVoid;
}  // namespace lte
}  // namespace magma

using grpc::ServerContext;
using magma::lte::SetSmNotificationContext;
using magma::lte::SetSMSessionContextAccess;
using magma::lte::SmContextVoid;
using magma::lte::SmfPduSessionSmContext;

namespace magma {
using namespace lte;

typedef struct ipv4_networks_s {
  uint32_t addr_hbo;
  int mask_len;
  bool success;
} ipv4_networks_t;

// SessionD to AMF server
class AmfServiceImpl final : public SmfPduSessionSmContext::Service {
 public:
  AmfServiceImpl();

  grpc::Status SetAmfNotification(ServerContext* context,
                                  const SetSmNotificationContext* notif,
                                  SmContextVoid* response) override;

  grpc::Status SetSmfSessionContext(ServerContext* context,
                                    const SetSMSessionContextAccess* request,
                                    SmContextVoid* response) override;
  bool SetSmfSessionContext_itti(
      const SetSMSessionContextAccess* request,
      itti_n11_create_pdu_session_response_t* itti_msg_p);
  grpc::Status SetAmfNotification_itti(
      const SetSmNotificationContext* notif,
      itti_n11_received_notification_t* itti_msg);
  bool fillUpPacketFilterContents(packet_filter_contents_t* pf_content,
                                  const FlowMatch* flow_match_rule);
  bool fillIpv6(packet_filter_contents_t* pf_content,
                const std::string ipv6network_str);
  bool fillIpv4(packet_filter_contents_t* pf_content,
                const std::string& ipv4network_str);
  ipv4_networks_t parseIpv4Network(const std::string& ipv4network_str);
};

}  // namespace magma

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

#include <grpc++/grpc++.h>
#include <lte/protos/session_manager.grpc.pb.h>
#include <lte/protos/session_manager.pb.h>
#include <lte/protos/subscriberdb.pb.h>
#include "orc8r/gateway/c/common/async_grpc/GRPCReceiver.hpp"
#include <stdint.h>
#include <functional>
#include <memory>
#include "lte/gateway/c/core/oai/common/common_types.h"
#include "lte/gateway/c/core/oai/lib/s6a_proxy/s6a_client_api.hpp"

using grpc::Status;
using magma::lte::SetSmNotificationContext;
using magma::lte::SetSMSessionContext;
using magma::lte::SmContextVoid;

namespace magma5g {

SetSMSessionContext create_sm_pdu_session(
    std::string&, uint8_t* apn, uint32_t pdu_session_id,
    uint32_t pdu_session_type, uint32_t gnb_gtp_teid, uint8_t pti,
    uint8_t* gnb_gtp_teid_ip_addr, std::string& ipv4_addr,
    std::string& ipv6_addr, const ambr_t& state_ambr, uint32_t version,
    const eps_subscribed_qos_profile_t& qos_profile);

class SmfServiceClient {
 public:
  virtual ~SmfServiceClient() {}
  virtual bool set_smf_session(SetSMSessionContext& request) = 0;
  virtual bool set_smf_notification(const SetSmNotificationContext& notify) = 0;
};

/**
 * AsyncSmfServiceClient implements SmfServiceClient but sends calls
 * asynchronously to sessiond.
 */
class AsyncSmfServiceClient : public magma::GRPCReceiver,
                              public SmfServiceClient {
 private:
  AsyncSmfServiceClient();
  static const uint32_t RESPONSE_TIMEOUT = 10;  // seconds
  std::unique_ptr<magma::lte::AmfPduSessionSmContext::Stub> stub_{};

  void SetSMFSessionRPC(
      SetSMSessionContext& request,
      const std::function<void(Status, SmContextVoid)>& callback);

  void SetSMFNotificationRPC(
      const SetSmNotificationContext& notify,
      const std::function<void(Status, SmContextVoid)>& callback);

 public:
  static AsyncSmfServiceClient& getInstance();

  AsyncSmfServiceClient(AsyncSmfServiceClient const&) = delete;
  void operator=(AsyncSmfServiceClient const&) = delete;

  int amf_smf_create_pdu_session(
      char* imsi, uint8_t* apn, uint32_t pdu_session_id,
      uint32_t pdu_session_type, uint32_t gnb_gtp_teid, uint8_t pti,
      uint8_t* gnb_gtp_teid_ip_addr, char* ue_ipv4_addr, char* ue_ipv6_addr,
      const ambr_t& state_ambr, uint32_t version,
      const eps_subscribed_qos_profile_t& qos_profile);

  bool set_smf_session(SetSMSessionContext& request);

  bool set_smf_notification(const SetSmNotificationContext& notify);

  bool n11_update_location_req(const s6a_update_location_req_t* const ulr_p);
};

}  // namespace magma5g

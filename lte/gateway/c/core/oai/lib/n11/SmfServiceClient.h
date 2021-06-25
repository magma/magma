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
#include "includes/GRPCReceiver.h"
#include <stdint.h>
#include <functional>
#include <memory>

using grpc::Status;
using magma::lte::SetSmNotificationContext;
using magma::lte::SetSMSessionContext;
using magma::lte::SmContextVoid;

namespace magma5g {

SetSMSessionContext create_sm_pdu_session_v4(
    char* imsi, uint8_t* apn, uint32_t pdu_session_id,
    uint32_t pdu_session_type, uint8_t* gnb_gtp_teid, uint8_t pti,
    uint8_t* gnb_gtp_teid_ip_addr, char* ipv4_addr, uint32_t version);

class SmfServiceClient {
 public:
  virtual ~SmfServiceClient() {}
  virtual bool set_smf_session(SetSMSessionContext& request)          = 0;
  virtual bool set_smf_notification(SetSmNotificationContext& notify) = 0;
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
      SetSmNotificationContext& notify,
      const std::function<void(Status, SmContextVoid)>& callback);

 public:
  static AsyncSmfServiceClient& getInstance();

  AsyncSmfServiceClient(AsyncSmfServiceClient const&) = delete;
  void operator=(AsyncSmfServiceClient const&) = delete;

  int amf_smf_create_pdu_session_ipv4(
      char* imsi, uint8_t* apn, uint32_t pdu_session_id,
      uint32_t pdu_session_type, uint8_t* gnb_gtp_teid, uint8_t pti,
      uint8_t* gnb_gtp_teid_ip_addr, char* ipv4_addr, uint32_t version);

  bool set_smf_session(SetSMSessionContext& request);

  bool set_smf_notification(SetSmNotificationContext& notify);
};

}  // namespace magma5g

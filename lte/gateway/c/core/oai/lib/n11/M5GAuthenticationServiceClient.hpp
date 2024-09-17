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

#include <stdint.h>
#include <functional>
#include <memory>

#include "lte/gateway/c/core/oai/lib/3gpp/3gpp_38.413.h"

#include <grpc++/grpc++.h>
#include "lte/protos/subscriberauth.grpc.pb.h"
#include "lte/protos/subscriberauth.pb.h"
#include "orc8r/gateway/c/common/async_grpc/GRPCReceiver.hpp"

using grpc::Status;
using magma::GRPCReceiver;
using magma::lte::M5GAuthenticationInformationAnswer;
using magma::lte::M5GAuthenticationInformationRequest;
using magma::lte::M5GSubscriberAuthentication;

namespace magma5g {
M5GAuthenticationInformationRequest create_subs_auth_request(
    const std::string& imsi, const std::string& snni);

class M5GAuthenticationServiceClient {
 public:
  virtual ~M5GAuthenticationServiceClient() {}
  virtual bool get_subs_auth_info(const std::string& imsi, uint8_t imsi_length,
                                  const char* snni, amf_ue_ngap_id_t ue_id) = 0;

  virtual bool get_subs_auth_info_resync(const std::string& imsi,
                                         uint8_t imsi_length, const char* snni,
                                         const void* resync_info,
                                         uint8_t resync_info_len,
                                         amf_ue_ngap_id_t ue_id) = 0;
};

/**
 * AsyncM5GAuthenticationServiceClient implements M5GAuthenticationServiceClient
 * but sends calls asynchronously to subscriberdb.
 */
class AsyncM5GAuthenticationServiceClient
    : public GRPCReceiver,
      public M5GAuthenticationServiceClient {
 public:
  bool get_subs_auth_info(const std::string& imsi, uint8_t imsi_length,
                          const char* snni, amf_ue_ngap_id_t ue_id);

  bool get_subs_auth_info_resync(const std::string& imsi, uint8_t imsi_length,
                                 const char* snni, const void* resync_info,
                                 uint8_t resync_info_len,
                                 amf_ue_ngap_id_t ue_id);

  static AsyncM5GAuthenticationServiceClient& getInstance();

  AsyncM5GAuthenticationServiceClient(
      AsyncM5GAuthenticationServiceClient const&) = delete;
  void operator=(AsyncM5GAuthenticationServiceClient const&) = delete;

 private:
  AsyncM5GAuthenticationServiceClient();
  static const uint32_t RESPONSE_TIMEOUT = 10;  // seconds
  std::unique_ptr<M5GSubscriberAuthentication::Stub> stub_{};

  void GetSubscriberAuthInfoRPC(
      M5GAuthenticationInformationRequest& request,
      const std::function<void(grpc::Status,
                               M5GAuthenticationInformationAnswer)>& callback);
};
}  // namespace magma5g

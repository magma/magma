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
#include "lte/protos/subscriberdb.grpc.pb.h"
#include "lte/protos/subscriberdb.pb.h"
#include "orc8r/gateway/c/common/async_grpc/includes/GRPCReceiver.h"

using grpc::Status;
using magma::GRPCReceiver;
using magma::lte::M5GSUCIRegistration;
using magma::lte::M5GSUCIRegistrationAnswer;
using magma::lte::M5GSUCIRegistrationRequest;

namespace magma5g {
M5GSUCIRegistrationRequest create_decrypt_imsi_request(
    const uint8_t ue_pubkey_identifier, const std::string& ue_pubkey,
    const std::string& ciphertext, const std::string& mac_tag);

class M5GSUCIRegistrationServiceClient {
 public:
  virtual ~M5GSUCIRegistrationServiceClient() {}
  virtual bool get_decrypt_imsi_info(
      const uint8_t ue_pubkey_identifier, const std::string& ue_pubkey,
      const std::string& ciphertext, const std::string& mac_tag,
      amf_ue_ngap_id_t ue_id) = 0;
};

class AsyncM5GSUCIRegistrationServiceClient
    : public GRPCReceiver,
      public M5GSUCIRegistrationServiceClient {
 public:
  bool get_decrypt_imsi_info(
      const uint8_t ue_pubkey_identifier, const std::string& ue_pubkey,
      const std::string& ciphertext, const std::string& mac_tag,
      amf_ue_ngap_id_t ue_id);

  static AsyncM5GSUCIRegistrationServiceClient& getInstance();

  AsyncM5GSUCIRegistrationServiceClient(
      AsyncM5GSUCIRegistrationServiceClient const&) = delete;
  void operator=(AsyncM5GSUCIRegistrationServiceClient const&) = delete;

 private:
  AsyncM5GSUCIRegistrationServiceClient();
  static const uint32_t RESPONSE_TIMEOUT = 10;  // seconds
  std::unique_ptr<M5GSUCIRegistration::Stub> stub_{};

  void GetSuciInfoRPC(
      const M5GSUCIRegistrationRequest& request,
      const std::function<void(Status, M5GSUCIRegistrationAnswer)>& callback);
};
}  // namespace magma5g

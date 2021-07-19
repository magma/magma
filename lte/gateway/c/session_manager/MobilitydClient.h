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

#include <mutex>
#include <unordered_map>
#include <memory>
#include <lte/protos/mobilityd.grpc.pb.h>
#include "Types.h"
#include "includes/GRPCReceiver.h"
#include "SessionState.h"

using grpc::Status;

namespace magma {
using namespace lte;
using std::experimental::optional;

/**
 * MobilitydClient is the base class for managing interactions with MobilityD.
 */
class MobilitydClient {
 public:
  virtual ~MobilitydClient() = default;

  /**
   * Get SubscriberID for correspoding of UE_IP
   */
  virtual void get_subscriberid_from_ipv4(
      const IPAddress& ue_ip_addr,
      std::function<void(Status status, SubscriberID)> callback) = 0;
};

/**
 * AsyncMobilitydClient sends asynchronous calls to Mobilityd to retrieve
 * UE information.
 */
class AsyncMobilitydClient : public GRPCReceiver, public MobilitydClient {
 public:
  AsyncMobilitydClient();

  explicit AsyncMobilitydClient(
      std::shared_ptr<grpc::Channel> mobilityd_channel);

  /**
   * Get SubscriberID for correspoding of UE_IP
   */
  void get_subscriberid_from_ipv4(
      const IPAddress& ue_ip_addr,
      std::function<void(Status status, SubscriberID)> callback);

 private:
  static const uint32_t RESPONSE_TIMEOUT = 6;  // seconds
  std::unique_ptr<MobilityService::Stub> stub_;
};
}  // namespace magma

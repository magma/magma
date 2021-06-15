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
#include <arpa/inet.h>
#include <memory>
#include <string>

#include "lte/protos/mobilityd.pb.h"
#include "lte/protos/mobilityd.grpc.pb.h"
#include "includes/GRPCReceiver.h"

using grpc::Status;

namespace magma {
namespace lte {

class MobilitydClient {
 public:
  virtual ~MobilitydClient() = default;

  /*
   * Get the subscriber id given its allocated IPv4 address. If the address
   * isn't associated with a subscriber, then it returns an error
   * @param addr: ipv4 address of subscriber
   * @param imsi (out): contains the imsi of the associated subscriber if it
   *                    exists
   * @return void
   */
  virtual void get_subscriber_id_from_ip(
      const struct in_addr& ip,
      std::function<void(Status status, SubscriberID)> callback) = 0;
};

/**
 * AsyncMobilitydClient sends asynchronous calls to mobilityd to retrieve
 * UE information.
 */
class AsyncMobilitydClient : public GRPCReceiver, public MobilitydClient {
 public:
  AsyncMobilitydClient();
  explicit AsyncMobilitydClient(
      std::shared_ptr<grpc::Channel> mobilityd_channel);

  void get_subscriber_id_from_ip(
      const struct in_addr& ip,
      std::function<void(Status status, SubscriberID)> callback);

 private:
  static const uint32_t RESPONSE_TIMEOUT_SECONDS = 6;
  std::unique_ptr<MobilityService::Stub> stub_;

 private:
  void get_subscriber_id_from_ip_rpc(
      const IPAddress& request,
      std::function<void(Status, SubscriberID)> callback);
};

}  // namespace lte
}  // namespace magma

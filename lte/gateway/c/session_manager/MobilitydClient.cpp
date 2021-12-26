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

#include <grpcpp/channel.h>
#include <memory>
#include <utility>

#include "MobilitydClient.h"
#include "includes/ServiceRegistrySingleton.h"
#include "lte/protos/mobilityd.grpc.pb.h"
#include "lte/protos/subscriberdb.pb.h"

namespace grpc {
class Status;
}  // namespace grpc
namespace magma {
namespace lte {
class IPAddress;
}  // namespace lte
}  // namespace magma

using grpc::Status;

namespace magma {

AsyncMobilitydClient::AsyncMobilitydClient(
    std::shared_ptr<grpc::Channel> channel)
    : stub_(MobilityService::NewStub(channel)) {}

AsyncMobilitydClient::AsyncMobilitydClient()
    : AsyncMobilitydClient(ServiceRegistrySingleton::Instance()->GetGrpcChannel(
          "mobilityd", ServiceRegistrySingleton::LOCAL)) {}

void AsyncMobilitydClient::get_subscriberid_from_ipv4(
    const IPAddress& ue_ip_addr,
    std::function<void(Status status, SubscriberID)> callback) {
  auto local_resp = new AsyncLocalResponse<SubscriberID>(std::move(callback),
                                                         RESPONSE_TIMEOUT);
  local_resp->set_response_reader(stub_->AsyncGetSubscriberIDFromIP(
      local_resp->get_context(), ue_ip_addr, &queue_));
}
}  // namespace magma

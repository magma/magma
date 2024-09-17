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

#include "lte/gateway/c/li_agent/src/MobilitydClient.hpp"

#include <grpcpp/channel.h>
#include <lte/protos/mobilityd.grpc.pb.h>
#include <lte/protos/mobilityd.pb.h>
#include <lte/protos/subscriberdb.pb.h>
#include <netinet/in.h>
#include <utility>

#include "orc8r/gateway/c/common/service_registry/ServiceRegistrySingleton.hpp"

namespace grpc {
class Status;
}  // namespace grpc

namespace magma {
namespace lte {

static magma::lte::IPAddress create_get_subscriber_id_from_ip_req(
    const struct in_addr& addr) {
  IPAddress req = IPAddress();
  req.set_version(IPAddress::IPV4);
  req.set_address(&addr, sizeof(struct in_addr));
  return req;
}

AsyncMobilitydClient::AsyncMobilitydClient(
    std::shared_ptr<grpc::Channel> channel)
    : stub_(MobilityService::NewStub(channel)) {}

AsyncMobilitydClient::AsyncMobilitydClient()
    : AsyncMobilitydClient(ServiceRegistrySingleton::Instance()->GetGrpcChannel(
          "mobilityd", ServiceRegistrySingleton::LOCAL)) {}

void AsyncMobilitydClient::get_subscriber_id_from_ip(
    const struct in_addr& ip,
    std::function<void(grpc::Status, SubscriberID)> callback) {
  IPAddress req = create_get_subscriber_id_from_ip_req(ip);
  get_subscriber_id_from_ip_rpc(req, callback);
}

void AsyncMobilitydClient::get_subscriber_id_from_ip_rpc(
    const IPAddress& request,
    std::function<void(grpc::Status, SubscriberID)> callback) {
  auto local_resp = new AsyncLocalResponse<SubscriberID>(
      std::move(callback), RESPONSE_TIMEOUT_SECONDS);
  local_resp->set_response_reader(stub_->AsyncGetSubscriberIDFromIP(
      local_resp->get_context(), request, &queue_));
}

}  // namespace lte
}  // namespace magma

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
#include "includes/EventdClient.h"
#include <grpcpp/channel.h>                        // for Channel
#include <grpcpp/impl/codegen/async_unary_call.h>  // for default_delete
#include <utility>                                 // for move
#include "includes/ServiceRegistrySingleton.h"     // for ServiceRegistrySin...
#include "orc8r/protos/common.pb.h"                // for Void
#include "orc8r/protos/eventd.grpc.pb.h"           // for EventService::Stub
namespace grpc {
class ClientContext;
}
namespace grpc {
class Status;
}
namespace magma {
namespace orc8r {
class Event;
}
}  // namespace magma

namespace magma {

using orc8r::Event;
using orc8r::EventService;
using orc8r::Void;

AsyncEventdClient& AsyncEventdClient::getInstance() {
  static AsyncEventdClient instance;
  return instance;
}

AsyncEventdClient::AsyncEventdClient() {
  std::shared_ptr<Channel> channel;
  channel = ServiceRegistrySingleton::Instance()->GetGrpcChannel(
      "eventd", ServiceRegistrySingleton::LOCAL);
  stub_ = EventService::NewStub(channel);
}

void AsyncEventdClient::log_event(
    const Event& request, std::function<void(Status status, Void)> callback) {
  auto local_response =
      new AsyncLocalResponse<Void>(std::move(callback), RESPONSE_TIMEOUT_SEC);
  local_response->set_response_reader(std::move(
      stub_->AsyncLogEvent(local_response->get_context(), request, &queue_)));
}

}  // namespace magma

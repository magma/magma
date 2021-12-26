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

#include <orc8r/protos/eventd.grpc.pb.h>  // for EventService::Stub, EventSe...
#include <stdint.h>                       // for uint32_t
#include <functional>                     // for function
#include <memory>                         // for unique_ptr
#include "orc8r/gateway/c/common/async_grpc/includes/GRPCReceiver.h"  // for GRPCReceiver
namespace grpc {
class Status;
}
namespace magma {
namespace orc8r {
class Event;
}
}  // namespace magma
namespace magma {
namespace orc8r {
class Void;
}
}  // namespace magma

using grpc::Status;

namespace magma {

/**
 * Base class for interfacing with EventD
 */
class EventdClient {
 public:
  virtual ~EventdClient() = default;
  virtual void log_event(
      const orc8r::Event& request,
      std::function<void(Status status, orc8r::Void)> callback) = 0;
};

/**
 * AsyncEventdClient sends asynchronous calls to EventD
 * to log events
 */
class AsyncEventdClient : public GRPCReceiver, public EventdClient {
 public:
  AsyncEventdClient(AsyncEventdClient const&) = delete;
  void operator=(AsyncEventdClient const&) = delete;

  static AsyncEventdClient& getInstance();

  void log_event(const orc8r::Event& request,
                 std::function<void(Status status, orc8r::Void)> callback);

 private:
  AsyncEventdClient();
  static const uint32_t RESPONSE_TIMEOUT_SEC = 6;
  std::unique_ptr<orc8r::EventService::Stub> stub_{};
};

}  // namespace magma

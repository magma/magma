/*
 * Licensed to the OpenAirInterface (OAI) Software Alliance under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The OpenAirInterface Software Alliance licenses this file to You under
 * the terms found in the LICENSE file in the root of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *-------------------------------------------------------------------------------
 * For more information about the OpenAirInterface (OAI) Software Alliance:
 *      contact@openairinterface.org
 */

#pragma once

#include <stdint.h>                              // for uint32_t
#include <functional>                            // for function
#include <memory>                                // for unique_ptr
#include "includes/GRPCReceiver.h"               // for GRPCReceiver
#include "lte/protos/apn.pb.h"                   // for lte
#include "lte/protos/session_manager.grpc.pb.h"  // for LocalSessionManager:...
namespace grpc {
class Status;
}  // namespace grpc
namespace magma {
namespace lte {
class LocalCreateSessionRequest;
}
}  // namespace magma
namespace magma {
namespace lte {
class LocalCreateSessionResponse;
}
}  // namespace magma
namespace magma {
namespace lte {
class LocalEndSessionRequest;
}
}  // namespace magma
namespace magma {
namespace lte {
class LocalEndSessionResponse;
}
}  // namespace magma
namespace magma {
namespace lte {
class PolicyBearerBindingRequest;
}
}  // namespace magma
namespace magma {
namespace lte {
class PolicyBearerBindingResponse;
}
}  // namespace magma
namespace magma {
namespace lte {
class UpdateTunnelIdsRequest;
}
}  // namespace magma
namespace magma {
namespace lte {
class UpdateTunnelIdsResponse;
}
}  // namespace magma

using grpc::Status;

namespace magma {
using namespace lte;

/**
 * PCEFClient is the main asynchronous client for interacting with sessiond.
 * Responses will come in a queue and call the callback passed
 * To start the client and make sure it receives calls, one must call the
 * rpc_response_loop method defined in the GRPCReceiver base class
 */
class PCEFClient : public GRPCReceiver {
 public:
  /**
   * Proxy a CreateSession gRPC call to sessiond
   */
  static void create_session(
      const LocalCreateSessionRequest& request,
      std::function<void(Status, LocalCreateSessionResponse)> callback);

  /**
   * Proxy an EndSession gRPC call to sessiond
   */
  static void end_session(
      const LocalEndSessionRequest& request,
      std::function<void(Status, LocalEndSessionResponse)> callback);

  /**
   * Proxy a BindPolicy2Bearer gRPC call to sessiond
   */
  static void bind_policy2bearer(
      const PolicyBearerBindingRequest& request,
      std::function<void(Status, PolicyBearerBindingResponse)> callback);

  /**
   * Proxy a UpdateTunnelIdsRequest gRPC call to sessiond
   */
  static void update_teids(
      const UpdateTunnelIdsRequest& request,
      std::function<void(Status, UpdateTunnelIdsResponse)> callback);

 public:
  PCEFClient(PCEFClient const&) = delete;
  void operator=(PCEFClient const&) = delete;

 private:
  PCEFClient();
  static PCEFClient& get_instance();
  std::unique_ptr<LocalSessionManager::Stub> stub_;
  static const uint32_t RESPONSE_TIMEOUT = 10;  // seconds
};

}  // namespace magma

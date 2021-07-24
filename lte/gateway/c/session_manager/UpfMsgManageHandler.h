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

#include <functional>
#include <memory>

#include "SessionStateEnforcer.h"
#include "SessionID.h"
#include "SessionReporter.h"
#include "SessionStore.h"

using grpc::Server;
using grpc::ServerContext;
using grpc::Status;

namespace magma {
using namespace orc8r;

/* SetUpfMsgManagHandler processes gRPC requests for the sessionD
 * This composites the user plane (pipelined)  sent messages over
 * sesiond channel
 */
class UpfMsgHandler {
 public:
  virtual ~UpfMsgHandler() {}

  /**
   * Node level GRPC message received from UPF
   * during startup
   */
  virtual void SetUPFNodeState(
      ServerContext* context, const UPFNodeState* node_request,
      std::function<void(Status, SmContextVoid)> response_callback) = 0;
};

/**
 * UPf Msg Handler processes gRPC requests to the session
 * manager. The handler uses a StateEnforcer as common synchronizing point for
 * any session state/rule changes either through UPF or access thread
 *
 */
class UpfMsgManageHandler : public UpfMsgHandler {
 public:
  UpfMsgManageHandler(
      std::shared_ptr<SessionStateEnforcer> enf, SessionStore& session_store);

  ~UpfMsgManageHandler() {}
  /**
   * Node level GRPC message received from UPF
   * during startup
   */

  virtual void SetUPFNodeState(
      ServerContext* context, const UPFNodeState* node_request,
      std::function<void(Status, SmContextVoid)> response_callback);

 private:
  SessionStore& session_store_;
  std::shared_ptr<SessionStateEnforcer> conv_enforcer_;
};

}  // namespace magma

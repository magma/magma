/*
 * Copyright 2020 The Magma Authors.
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
#include "UpfMsgManageHandler.h"
#include <google/protobuf/util/time_util.h>
#include <chrono>
#include <thread>
#include <memory>
#include <string>
#include "magma_logging.h"
#include "GrpcMagmaUtils.h"
#include "lte/protos/session_manager.pb.h"

using grpc::Status;

namespace magma {
/**
 * SetInterfaceForUserPlaneHandler processes gRPC requests for the sessionD
 * This composites the all the request that comes from UPF
 */

UpfMsgManageHandler::UpfMsgManageHandler(
    std::shared_ptr<SessionStateEnforcer> enforcer, SessionStore& session_store)
    : session_store_(session_store), conv_enforcer_(enforcer) {}

/**
 * Node level GRPC message received from UPF
 * during startup
 */

void UpfMsgManageHandler::SetUPFNodeState(
    ServerContext* context, const UPFNodeState* node_request,
    std::function<void(Status, SmContextVoid)> response_callback) {
  auto& request = *node_request;
  // Print the message from UPF
  PrintGrpcMessage(static_cast<const google::protobuf::Message&>(request));
  Status status;
  conv_enforcer_->get_event_base().runInEventBaseThread([this,
                                                         response_callback,
                                                         request]() {
    switch (request.upf_node_messages_case()) {
      case UPFNodeState::kAssociatonState: {
        std::string upf_id            = request.upf_id();
        UPFAssociationState Assostate = request.associaton_state();
        auto recovery_time            = Assostate.recovery_time_stamp();
        auto feature_set              = Assostate.feature_set();
        // For now get User Plan IPv4 resource at index '0' only
        std::string ipv4_addr = Assostate.ip_resource_schema(0).ipv4_address();
        // Set the UPF address
        conv_enforcer_->set_upf_node(upf_id, ipv4_addr);
        // send the same UPF Node association response back to UPF
        Status status(grpc::OK, "UPF Node message supported");
        break;
      }
      default:
        MLOG(MDEBUG) << "UPF Node message of type "
                     << request.upf_node_messages_case()
                     << "Not Being Handled in SMF";
        Status status(grpc::UNIMPLEMENTED, "UPF Node message not supported");
    }
  });
  response_callback(status, SmContextVoid());
  return;
}

}  // end namespace magma

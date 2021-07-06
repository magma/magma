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
#include "SessionStateEnforcer.h"
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

/**
 * Periodic messages about UPF session config
 *
 */
void UpfMsgManageHandler::SetUPFSessionsConfig(
    ServerContext* context, const UPFSessionConfigState* sess_config,
    std::function<void(Status, SmContextVoid)> response_callback) {
  auto& ses_config = *sess_config;
  int32_t count    = 0;
  conv_enforcer_->get_event_base().runInEventBaseThread([this, &count,
                                                         ses_config]() {
    for (auto& upf_session : ses_config.upf_session_state()) {
      // Deleting the IMSI prefix from imsi
      std::string imsi_upf = upf_session.subscriber_id();
      std::string imsi     = imsi_upf.substr(4, imsi_upf.length() - 4);
      uint32_t version     = upf_session.session_version();
      uint32_t teid        = upf_session.local_f_teid();
      auto session_map     = session_store_.read_sessions({imsi});
      /* Search with session search criteria of IMSI and session_id and
       * find  respective sesion to operate
       */
      SessionSearchCriteria criteria(imsi, IMSI_AND_TEID, teid);
      auto session_it = session_store_.find_session(session_map, criteria);
      if (!session_it) {
        MLOG(MERROR) << "No session found in SessionMap for IMSI " << imsi
                     << " with teid " << teid;
        continue;
      }
      auto& session    = **session_it;
      auto cur_version = session->get_current_version();
      if (version < cur_version) {
        MLOG(MINFO) << "UPF verions of session imsi " << imsi << " of  teid "
                    << teid << " recevied version " << version
                    << " SMF latest version: " << cur_version << " Resending";
        if (conv_enforcer_->is_incremented_rtx_counter_within_max(session)) {
          conv_enforcer_->m5g_send_session_request_to_upf(session);
        }
      } else {
        count++;
      }
    }
#if 0
    if (ses_config.upf_session_state_size() != count) {
      MLOG(MINFO) << "UPF periodic report config missmatch session:"
                  << (ses_config.upf_session_state_size() - count);
    }
#endif
  });
  response_callback(Status::OK, SmContextVoid());
  return;
}

}  // end namespace magma

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

#include <arpa/inet.h>
#include <folly/io/async/EventBase.h>
#include <glog/logging.h>
#include <grpcpp/impl/codegen/status.h>
#include <grpcpp/impl/codegen/status_code_enum.h>
#include <netinet/in.h>
#include <experimental/optional>
#include <memory>
#include <ostream>
#include <string>
#include <vector>

#include "GrpcMagmaUtils.h"
#include "MobilitydClient.h"
#include "SessionState.h"
#include "SessionStateEnforcer.h"
#include "SessionStore.h"
#include "Types.h"
#include "lte/protos/mobilityd.pb.h"
#include "lte/protos/session_manager.pb.h"
#include "lte/protos/subscriberdb.pb.h"
#include "magma_logging.h"

namespace google {
namespace protobuf {
class Message;
}  // namespace protobuf
}  // namespace google
namespace grpc {
class ServerContext;
}  // namespace grpc

using grpc::Status;

namespace magma {
/**
 * SetInterfaceForUserPlaneHandler processes gRPC requests for the sessionD
 * This composites the all the request that comes from UPF
 */

UpfMsgManageHandler::UpfMsgManageHandler(
    std::shared_ptr<SessionStateEnforcer> enforcer,
    std::shared_ptr<MobilitydClient> mobilityd_client,
    SessionStore& session_store)
    : session_store_(session_store),
      conv_enforcer_(enforcer),
      mobilityd_client_(mobilityd_client) {}

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
        std::string upf_id = request.upf_id();
        UPFAssociationState Assostate = request.associaton_state();
        auto recovery_time = Assostate.recovery_time_stamp();
        auto feature_set = Assostate.feature_set();
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
  int32_t count = 0;
  conv_enforcer_->get_event_base().runInEventBaseThread([this, &count,
                                                         ses_config]() {
    for (auto& upf_session : ses_config.upf_session_state()) {
      // Deleting the IMSI prefix from imsi
      std::string imsi_upf = upf_session.subscriber_id();
      std::string imsi = imsi_upf.substr(4, imsi_upf.length() - 4);
      uint32_t version = upf_session.session_version();
      uint32_t teid = upf_session.local_f_teid();
      auto session_map = session_store_.read_sessions({imsi});
      /* Search with session search criteria of IMSI and session_id and
       * find  respective session to operate
       */
      SessionSearchCriteria criteria(imsi, IMSI_AND_TEID, teid);
      auto session_it = session_store_.find_session(session_map, criteria);
      if (!session_it) {
        MLOG(MERROR) << "No session found in SessionMap for IMSI " << imsi
                     << " with teid " << teid;
        continue;
      }
      auto& session = **session_it;
      auto cur_version = session->get_current_version();
      if (version < cur_version) {
        MLOG(MINFO) << "UPF verions of session imsi " << imsi << " of  teid "
                    << teid << " recevied version " << version
                    << " SMF latest version: " << cur_version << " Resending";
        if (conv_enforcer_->is_incremented_rtx_counter_within_max(session)) {
          RulesToProcess pending_activation, pending_deactivation;
          const CreateSessionResponse& csr =
              session->get_create_session_response();
          std::vector<StaticRuleInstall> static_rule_installs =
              conv_enforcer_->to_vec(csr.static_rules());
          std::vector<DynamicRuleInstall> dynamic_rule_installs =
              conv_enforcer_->to_vec(csr.dynamic_rules());

          session->process_get_5g_rule_installs(
              static_rule_installs, dynamic_rule_installs, &pending_activation,
              &pending_deactivation);
          conv_enforcer_->m5g_send_session_request_to_upf(
              session, pending_activation, pending_deactivation);
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

void UpfMsgManageHandler::SendPagingRequest(
    ServerContext* context, const UPFPagingInfo* page_request,
    std::function<void(Status, SmContextVoid)> response_callback) {
  auto& pag_req = *page_request;

  uint32_t fte_id = pag_req.local_f_teid();
  std::string ip_addr = pag_req.ue_ip_addr();
  struct in_addr ue_ip;
  IPAddress req = IPAddress();

  inet_aton(ip_addr.c_str(), &ue_ip);
  req.set_version(IPAddress::IPV4);
  req.set_address(&ue_ip, sizeof(struct in_addr));

  mobilityd_client_->get_subscriberid_from_ipv4(
      req, [this, fte_id, response_callback](Status status,
                                             const SubscriberID& sid) {
        if (!status.ok()) {
          MLOG(MERROR) << "Subscriber could not be found for ip ";
        }
        const std::string& imsi = sid.id();
        get_session_from_imsi(imsi, fte_id, response_callback);
        return;
      });
}

void UpfMsgManageHandler::get_session_from_imsi(
    const std::string& imsi, uint32_t te_id,
    std::function<void(Status, SmContextVoid)> response_callback) {
  conv_enforcer_->get_event_base().runInEventBaseThread([this, imsi, te_id,
                                                         response_callback]() {
    if (!imsi.length()) {
      MLOG(MERROR) << "get_subscriberid_from_ipv4 for IP"
                      "returned an empty subscriber ID";
      Status status(grpc::NOT_FOUND,
                    "Session Not found because"
                    "subscriber ID not found for IP");
      response_callback(status, SmContextVoid());
      return;
    }

    // retrieve session_map entry
    auto session_map = session_store_.read_sessions({imsi});
    /* Search with session search criteria of IMSI and session_id and
     * find  respective session to operate
     */
    SessionSearchCriteria criteria(imsi, IMSI_AND_TEID, te_id);

    auto session_it = session_store_.find_session(session_map, criteria);
    if (!session_it) {
      MLOG(MERROR) << "No session found in SessionMap for IMSI " << imsi
                   << " with teid " << te_id;
      Status status(grpc::NOT_FOUND,
                    "Session was not found for IMSI with teid");
      response_callback(status, SmContextVoid());
      return;
    }

    auto& session = **session_it;
    MLOG(MINFO) << "IDLE_MODE::: Session found in SendingPaging "
                   "Request of imsi: "
                << imsi << "  session_id: " << session->get_session_id();
    /* Generate Paging notification to AMF, only if session is in INACTIVE
     * state.
     */
    if (session->get_state() == INACTIVE) {
      conv_enforcer_->handle_state_update_to_amf(
          *session, magma::lte::M5GSMCause::OPERATION_SUCCESS,
          UE_PAGING_NOTIFY);
      MLOG(MDEBUG) << "UPF Paging notification forwarded to AMF of imsi:"
                   << imsi;
      response_callback(Status::OK, SmContextVoid());
    } else {
      MLOG(MDEBUG) << "Can not Trigger Paging notification to AMF, as session "
                      "is not an INACTIVE state.";
      return;
    }
  });
}
}  // end namespace magma

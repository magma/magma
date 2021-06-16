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
#include <sys/socket.h>
#include <netinet/in.h>
#include <arpa/inet.h>
#include "magma_logging.h"
#include "GrpcMagmaUtils.h"
#include "lte/protos/session_manager.pb.h"
#include "lte/protos/mobilityd.grpc.pb.h"
#include "lte/protos/mobilityd.pb.h"
#include <MobilityServiceClient.h>

using grpc::Status;
using magma::lte::MobilityServiceClient;

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
  MLOG(MDEBUG) << "Node UPF details :";
  conv_enforcer_->get_event_base().runInEventBaseThread([this, request]() {
    UPFNodeState::UpfNodeMessagesCase msgtype =
        request.upf_node_messages_case();
    if (msgtype == UPFNodeState::kAssociatonState) {
      std::string upf_id            = request.upf_id();
      UPFAssociationState Assostate = request.associaton_state();
      auto recovery_time            = Assostate.recovery_time_stamp();
      auto feature_set              = Assostate.feature_set();
      // For now get User Plan IPv4 resource at index '0' only
      std::string ipv4_addr = Assostate.ip_resource_schema(0).ipv4_address();
      // Set the UPF address
      conv_enforcer_->set_upf_node(upf_id, ipv4_addr);
      // send the same UPF Node association response back to UPF
    }
  });
  response_callback(Status::OK, SmContextVoid());
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
    for (int i = 0; i < ses_config.upf_session_state_size(); i++) {
      // Deleting the IMSI prefix from imsi
      std::string imsi_upf = ses_config.upf_session_state(i).subscriber_id();
      std::string imsi     = imsi_upf.substr(4, imsi_upf.length() - 4);
      uint32_t version     = ses_config.upf_session_state(i).session_version();
      uint32_t teid        = ses_config.upf_session_state(i).local_f_teid();
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
        if (session->inc_rtx_counter()) {
          conv_enforcer_->m5g_send_session_request_to_upf(imsi, session);
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
/**
 * Upf Paging request
 *
 */
void UpfMsgManageHandler::SendPagingRequest(
    ServerContext* context, const UPFPagingInfo* page_request,
    std::function<void(Status, SmContextVoid)> response_callback) {
  auto& pag_req = *page_request;

  std::string imsi;
  uint32_t fte_id     = pag_req.local_f_teid();
  std::string ip_addr = pag_req.ue_ip_addr();
  struct in_addr ue_ip;

  inet_aton(ip_addr.c_str(), &ue_ip);
  int ret = MobilityServiceClient::getInstance().GetSubscriberIDFromIPv4(
      ue_ip, &imsi);
  if (ret > 0) {
    MLOG(MERROR) << "Subscriber could not be found for ip " << ip_addr;
    return;
  }

  MLOG(MINFO) << "IDLE_MODE::: SendingPagingRequest for subscriber" << imsi;
  conv_enforcer_->get_event_base().runInEventBaseThread([this, imsi, fte_id,
                                                         response_callback]() {
    // Deleting the IMSI prefix from imsi
    std::string imsi_temp = imsi.substr(4, imsi.length() - 4);
    // retreive session_map entry
    auto session_map = session_store_.read_sessions({imsi_temp});
    /* Search with session search criteria of IMSI and session_id and
     * find  respective sesion to operate
     */
    SessionSearchCriteria criteria(imsi_temp, IMSI_AND_TEID, fte_id);

    auto session_it = session_store_.find_session(session_map, criteria);
    if (!session_it) {
      MLOG(MERROR) << "No session found in SessionMap for IMSI " << imsi_temp
                   << " with teid " << fte_id;
      Status status(grpc::NOT_FOUND, "Sesion Not found");
      response_callback(status, SmContextVoid());
      return;
    }

    MLOG(MINFO)
        << "IDLE_MODE::: Session found in SendingPaging Request of imsi:"
        << imsi_temp;
    auto& session = **session_it;
    // Generate Paging trigget to AMF.
    conv_enforcer_->handle_state_update_to_amf(
        *session, magma::lte::M5GSMCause::OPERATION_SUCCESS, UE_PAGING_NOTIFY);
  });
  MLOG(MINFO) << "UPF Paging notificaiton forwarded to AMF of imsi:" << imsi;
  response_callback(Status::OK, SmContextVoid());
  return;
}
}  // end namespace magma

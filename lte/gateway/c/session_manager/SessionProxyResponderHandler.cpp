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
#include "lte/gateway/c/session_manager/SessionProxyResponderHandler.hpp"

#include <folly/io/async/EventBase.h>
#include <glog/logging.h>
#include <grpcpp/impl/codegen/status.h>
#include <grpcpp/impl/codegen/status_code_enum.h>
#include <lte/protos/abort_session.pb.h>
#include <lte/protos/session_manager.pb.h>
#include <ostream>
#include <string>

#include "lte/gateway/c/session_manager/EnumToString.hpp"
#include "lte/gateway/c/session_manager/GrpcMagmaUtils.hpp"
#include "lte/gateway/c/session_manager/LocalEnforcer.hpp"
#include "lte/gateway/c/session_manager/SessionStore.hpp"
#include "orc8r/gateway/c/common/logging/magma_logging.hpp"

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
SessionProxyResponderHandlerImpl::SessionProxyResponderHandlerImpl(
    std::shared_ptr<LocalEnforcer> enforcer, SessionStore& session_store)
    : enforcer_(enforcer), session_store_(session_store) {}

SessionProxyResponderHandlerImpl::SessionProxyResponderHandlerImpl(
    std::shared_ptr<LocalEnforcer> enforcer,
    std::shared_ptr<SessionStateEnforcer> m5genforcer,
    SessionStore& session_store)
    : enforcer_(enforcer),
      m5genforcer_(m5genforcer),
      session_store_(session_store) {}

void SessionProxyResponderHandlerImpl::ChargingReAuth(
    ServerContext* context, const ChargingReAuthRequest* request,
    std::function<void(Status, ChargingReAuthAnswer)> response_callback) {
  auto& request_cpy = *request;
  PrintGrpcMessage(static_cast<const google::protobuf::Message&>(request_cpy));
  MLOG(MDEBUG) << "Received a Gy (Charging) ReAuthRequest for "
               << request->session_id() << " and charging_key "
               << request->charging_key();
  enforcer_->get_event_base().runInEventBaseThread([this, request_cpy,
                                                    response_callback]() {
    auto session_map = session_store_.read_sessions({request_cpy.sid()});
    SessionUpdate update =
        SessionStore::get_default_session_update(session_map);
    auto result =
        enforcer_->init_charging_reauth(session_map, request_cpy, update);
    MLOG(MDEBUG) << "Result of Gy (Charging) ReAuthRequest "
                 << raa_result_to_str(result);
    ChargingReAuthAnswer ans;
    Status status;
    ans.set_result(result);

    bool update_success = session_store_.update_sessions(update);
    if (update_success) {
      PrintGrpcMessage(static_cast<const google::protobuf::Message&>(ans));
      status = Status::OK;
    } else {
      // Todo If update fails, we should rollback changes from the request
      MLOG(MERROR) << "Failed to update Gy (Charging) ReAuthRequest changes...";
      auto status =
          Status(grpc::ABORTED,
                 "ChargingReAuth no longer valid due to another update that "
                 "updated the session first.");
      PrintGrpcMessage(static_cast<const google::protobuf::Message&>(ans));
    }
    response_callback(status, ans);
    MLOG(MDEBUG) << "Sent RAA response " << result << " for Gy ReAuth "
                 << request_cpy.session_id();
  });
}

void SessionProxyResponderHandlerImpl::PolicyReAuth(
    ServerContext* context, const PolicyReAuthRequest* request,
    std::function<void(Status, PolicyReAuthAnswer)> response_callback) {
  auto& request_cpy = *request;
  PrintGrpcMessage(static_cast<const google::protobuf::Message&>(request_cpy));
  MLOG(MDEBUG) << "Received a Gx (Policy) ReAuthRequest for session_id "
               << request->session_id();
  enforcer_->get_event_base().runInEventBaseThread([this, request_cpy,
                                                    response_callback]() {
    PolicyReAuthAnswer ans;
    const std::string& imsi = request_cpy.imsi();
    const std::string& session_id = request_cpy.session_id();
    SessionSearchCriteria criteria(imsi, IMSI_AND_SESSION_ID, session_id);
    auto session_map = session_store_.read_sessions({imsi});
    auto session_it = session_store_.find_session(session_map, criteria);
    if (!session_it) {
      MLOG(MERROR) << "No session found in SessionMap for IMSI " << imsi
                   << " with session_id " << session_id;
      PrintGrpcMessage(static_cast<const google::protobuf::Message&>(ans));
      ans.set_result(ReAuthResult::SESSION_NOT_FOUND);
      Status status(grpc::NOT_FOUND, "Session not found");
      response_callback(status, ans);
      return;
    }

    auto& session = **session_it;
    const auto& config = session->get_config();
    const auto& common_context = config.common_context;

    SessionUpdate update =
        SessionStore::get_default_session_update(session_map);
    const auto rat_type = common_context.rat_type();
    const std::string& interface = (rat_type == TGPP_NR) ? "5G N7" : "4G Gx";
    if (rat_type == TGPP_NR) {
      m5genforcer_->init_policy_reauth(imsi, session_map, request_cpy, ans,
                                       update);
    } else {
      enforcer_->init_policy_reauth(session_map, request_cpy, ans, update);
    }
    MLOG(MDEBUG) << "Result of " << interface << " (Policy) ReAuthRequest "
                 << raa_result_to_str(ans.result());
    bool update_success = session_store_.update_sessions(update);
    if (update_success) {
      MLOG(MDEBUG) << "Sending RAA response for " << interface << " ReAuth "
                   << request_cpy.session_id();
      PrintGrpcMessage(static_cast<const google::protobuf::Message&>(ans));
      response_callback(Status::OK, ans);
    } else {
      // Todo If update fails, we should rollback changes from the request
      MLOG(MERROR) << "Failed to update "
                   << interface << " (Policy) ReAuthRequest changes...";
      auto status =
          Status(grpc::ABORTED,
                 "PolicyReAuth no longer valid due to another update that "
                 "updated the session first.");
      PrintGrpcMessage(static_cast<const google::protobuf::Message&>(ans));
      response_callback(status, ans);
    }
    MLOG(MDEBUG) << "Sent RAA response for " << interface << " ReAuth "
                 << request_cpy.session_id();
  });
}

void SessionProxyResponderHandlerImpl::AbortSession(
    ServerContext* context, const AbortSessionRequest* request,
    std::function<void(Status, AbortSessionResult)> response_callback) {
  PrintGrpcMessage(static_cast<const google::protobuf::Message&>(*request));
  auto imsi = request->user_name();
  // SessionD currently stores IMSIs with the 'IMSI' prefix so append if it is
  // not there already
  if (imsi.find("IMSI") == std::string::npos) {
    imsi = "IMSI" + imsi;
    MLOG(MDEBUG) << "Appended 'IMSI' to ASR user_name: " << imsi;
  }
  const auto session_id = request->session_id();
  MLOG(MINFO) << "Received an ASR for session_id " << session_id;
  enforcer_->get_event_base().runInEventBaseThread([this, imsi, session_id,
                                                    response_callback]() {
    grpc::Status status = Status::OK;
    AbortSessionResult answer;
    auto session_map = session_store_.read_sessions({imsi});
    SessionUpdate updates =
        SessionStore::get_default_session_update(session_map);
    auto found =
        enforcer_->handle_abort_session(session_map, imsi, session_id, updates);
    if (found) {
      answer.set_code(AbortSessionResult_Code_SESSION_REMOVED);
    } else {
      answer.set_code(AbortSessionResult_Code_SESSION_NOT_FOUND);
    }
    bool update_success = session_store_.update_sessions(updates);
    if (!update_success) {
      MLOG(MERROR) << "SessionStore update failed when processing ASR for "
                   << session_id;
      status =
          Status(grpc::ABORTED, "ASR failed due to internal datastore error");
      answer.set_code(AbortSessionResult_Code_SESSION_NOT_FOUND);
    }
    PrintGrpcMessage(static_cast<const google::protobuf::Message&>(answer));
    response_callback(status, answer);
    MLOG(MINFO) << "Sent ASA for " << session_id << " with code "
                << asr_result_to_str(answer.code());
  });
}
}  // namespace magma

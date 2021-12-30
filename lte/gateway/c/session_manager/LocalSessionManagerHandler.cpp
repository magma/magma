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
#include <folly/io/async/EventBase.h>
#include <glog/logging.h>
#include <grpcpp/impl/codegen/status_code_enum.h>
#include <chrono>
#include <exception>
#include <memory>
#include <ostream>
#include <string>
#include <typeinfo>
#include <unordered_map>
#include <utility>
#include <vector>

#include "DirectorydClient.h"
#include "GrpcMagmaUtils.h"
#include "LocalEnforcer.h"
#include "LocalSessionManagerHandler.h"
#include "ServiceAction.h"
#include "SessionCredit.h"
#include "SessionEvents.h"
#include "SessionReporter.h"
#include "SessionState.h"
#include "StoredState.h"
#include "Types.h"
#include "Utilities.h"
#include "includes/SentryWrapper.h"
#include "lte/protos/pipelined.pb.h"
#include "lte/protos/session_manager.pb.h"
#include "lte/protos/subscriberdb.pb.h"
#include "magma_logging.h"
#include "orc8r/protos/directoryd.pb.h"

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

LocalSessionManagerHandlerImpl::LocalSessionManagerHandlerImpl(
    std::shared_ptr<LocalEnforcer> enforcer, SessionReporter* reporter,
    std::shared_ptr<DirectorydClient> directoryd_client,
    std::shared_ptr<EventsReporter> events_reporter,
    SessionStore& session_store)
    : session_store_(session_store),
      enforcer_(enforcer),
      reporter_(reporter),
      directoryd_client_(directoryd_client),
      events_reporter_(events_reporter),
      current_epoch_(0),
      reported_epoch_(0),
      retry_timeout_ms_(std::chrono::milliseconds{5000}),
      pipelined_state_(PipelineDState::NOT_READY) {}

void LocalSessionManagerHandlerImpl::ReportRuleStats(
    ServerContext* context, const RuleRecordTable* request,
    std::function<void(Status, Void)> response_callback) {
  set_sentry_transaction("ReportRuleStats");
  auto& request_cpy = *request;
  if (request_cpy.records_size() > 0) {
    PrintGrpcMessage(
        static_cast<const google::protobuf::Message&>(request_cpy));
  }
  enforcer_->get_event_base().runInEventBaseThread([this, request_cpy]() {
    if (!session_store_.is_ready()) {
      // Since PipelineD reports a delta value for usage, this could lead to
      // SessionD missing some usage if Redis becomes unavailable. However,
      // since we usually only see this case on service restarts, we'll let this
      // slide for now. Once we move to flat-rate reporting from PipelineD this
      // will no longer be an issue.
      MLOG(MINFO) << "SessionStore client is not yet ready... Ignoring this "
                     "RuleRecordTable";
      return;
    }
    auto session_map = session_store_.read_all_sessions();
    SessionUpdate update =
        SessionStore::get_default_session_update(session_map);
    MLOG(MDEBUG) << "Aggregating " << request_cpy.records_size() << " records";
    enforcer_->aggregate_records(session_map, request_cpy, update);
    check_usage_for_reporting(std::move(session_map), update);
  });

  reported_epoch_ = request_cpy.epoch();
  if (is_pipelined_restarted()) {
    MLOG(MINFO) << "Pipelined has been restarted, attempting to sync flows,"
                << " old epoch = " << current_epoch_
                << ", new epoch = " << reported_epoch_;

    enforcer_->get_event_base().runInEventBaseThread(
        [this, epoch = reported_epoch_,
         update_rule_versions = request_cpy.update_rule_versions()]() {
          call_setup_pipelined(epoch, update_rule_versions);
        });
    // Set the current epoch right away to prevent double setup call requests
    current_epoch_ = reported_epoch_;
  }
  response_callback(Status::OK, Void());
}

void LocalSessionManagerHandlerImpl::check_usage_for_reporting(
    SessionMap session_map, SessionUpdate& session_uc) {
  std::vector<std::unique_ptr<ServiceAction>> actions;
  auto request = enforcer_->collect_updates(session_map, actions, session_uc);
  enforcer_->execute_actions(session_map, actions, session_uc);
  if (request.updates_size() == 0 && request.usage_monitors_size() == 0) {
    auto update_success = session_store_.update_sessions(session_uc);
    if (update_success) {
      MLOG(MDEBUG) << "Succeeded in updating session after no reporting";
    } else {
      MLOG(MERROR) << "Failed in updating session after no reporting";
    }
    return;  // nothing to report
  }

  MLOG(MINFO) << "Sending " << request.updates_size()
              << " charging updates and " << request.usage_monitors_size()
              << " monitor updates to OCS and PCRF";

  // set reporting flag for those sessions reporting
  session_store_.set_and_save_reporting_flag(true, request, session_uc);

  // Before reporting and returning control to the event loop, increment the
  // request numbers stored for the sessions in SessionStore
  session_store_.sync_request_numbers(session_uc);

  // report to cloud
  // NOTE: It is not possible to construct a std::function from a move-only type
  //       So because of this, we can't directly move session_map into the
  //       value-capture of the callback. As a workaround, a shared_ptr to
  //       the session_map is used.
  //       Check
  //       https://stackoverflow.com/questions/25421346/how-to-create-an-stdfunction-from-a-move-capturing-lambda-expression
  reporter_->report_updates(
      request,
      [this, request, session_uc,
       session_map_ptr = std::make_shared<SessionMap>(std::move(session_map))](
          Status status, UpdateSessionResponse response) mutable {
        enforcer_->handle_session_update_response(request, session_map_ptr,
                                                  session_uc, status, response);
      });
}

bool LocalSessionManagerHandlerImpl::is_pipelined_restarted() {
  // If 0 also setup pipelined because it always waits for setup instructions
  return (current_epoch_ == 0 || current_epoch_ != reported_epoch_);
}

void LocalSessionManagerHandlerImpl::handle_setup_callback(
    const std::uint64_t& epoch, Status status, SetupFlowsResult resp) {
  // Run everything in the event base thread since we asynchronously
  // read/modify pipelined_state_
  enforcer_->get_event_base().runInEventBaseThread([=] {
    if (status.ok() && resp.result() == resp.SUCCESS) {
      MLOG(MINFO) << "Successfully setup PipelineD with epoch: " << epoch;
      pipelined_state_ = PipelineDState::READY;
      return;
    }
    pipelined_state_ = PipelineDState::NOT_READY;
    if (current_epoch_ != epoch) {
      // This means that PipelineD has restarted since the initial Setup call
      // was called
      MLOG(MDEBUG) << "Received stale PipelineD setup callback for epoch: "
                   << epoch << ", current epoch: " << current_epoch_;
      return;
    }
    if (status.ok() && resp.result() == resp.OUTDATED_EPOCH) {
      MLOG(MWARNING) << "PipelineD setup call has outdated epoch, abandoning.";
      return;
    }

    // Cases for which we re-try the Setup call
    if (!status.ok()) {
      MLOG(MERROR) << "Could not setup PipelineD, rpc failed with: "
                   << status.error_message() << ", retrying PipelineD setup "
                   << "for epoch: " << epoch;
    } else if (resp.result() == resp.FAILURE) {
      MLOG(MWARNING) << "PipelineD setup failed, retrying PipelineD setup "
                     << "after delay, for epoch: " << epoch;
    }

    enforcer_->get_event_base().runAfterDelay(
        [=] { call_setup_pipelined(epoch, false); }, retry_timeout_ms_.count());
  });
}

void LocalSessionManagerHandlerImpl::call_setup_pipelined(
    const std::uint64_t& epoch, const bool update_rule_versions) {
  using namespace std::placeholders;
  if (pipelined_state_ == PipelineDState::SETTING_UP) {
    // Return if there is already a Setup call in progress
    return;
  }
  if (current_epoch_ != epoch) {
    // This means that PipelineD has restarted since the this call was
    // scheduled
    return;
  }
  pipelined_state_ = PipelineDState::SETTING_UP;

  auto session_map = session_store_.read_all_sessions();
  if (update_rule_versions) {
    MLOG(MINFO) << "Incrementing all policy version for pipelined setup";
    enforcer_->increment_all_policy_versions(session_map);
  }

  MLOG(MINFO) << "Sending a setup call to PipelineD with epoch: " << epoch;
  enforcer_->setup(
      session_map, epoch,
      std::bind(&LocalSessionManagerHandlerImpl::handle_setup_callback, this,
                epoch, _1, _2));
  return;
}

static CreateSessionRequest make_create_session_request(
    const SessionConfig& cfg, const std::string& session_id,
    const std::unique_ptr<Timezone>& access_timezone) {
  CreateSessionRequest create_request;
  create_request.set_session_id(session_id);
  create_request.mutable_common_context()->CopyFrom(cfg.common_context);
  create_request.mutable_rat_specific_context()->CopyFrom(
      cfg.rat_specific_context);

  if (access_timezone != nullptr) {
    create_request.mutable_access_timezone()->CopyFrom(*access_timezone);
  }

  const RequestedUnits requestedUnits =
      SessionCredit::get_initial_requested_credits_units();
  create_request.mutable_requested_units()->CopyFrom(requestedUnits);

  return create_request;
}

grpc::Status LocalSessionManagerHandlerImpl::check_sessiond_is_ready() {
  if (pipelined_state_ != READY) {
    MLOG(MINFO) << "Rejecting requests since PipelineD is still setting up";
    return Status(grpc::UNAVAILABLE, "PipelineD is not ready");
  }
  if (!session_store_.is_ready()) {
    MLOG(MINFO)
        << "Rejecting requests since SessionStore (Redis) is unavailable";
    return Status(grpc::UNAVAILABLE, "Storage backend is not available");
  }
  return Status::OK;
}

grpc::Status LocalSessionManagerHandlerImpl::validate_create_session_request(
    const SessionConfig cfg) {
  const auto rat_type = cfg.common_context.rat_type();
  const CommonSessionContext common = cfg.common_context;
  if (rat_type != TGPP_WLAN && rat_type != TGPP_LTE) {
    // We don't support outside of WLAN / LTE
    std::ostringstream failure_stream;
    failure_stream << "Received an invalid RAT type " << rat_type;
    std::string failure_msg = failure_stream.str();
    MLOG(MERROR) << failure_msg;
    events_reporter_->session_create_failure(cfg, failure_msg);
    return Status(grpc::FAILED_PRECONDITION, failure_msg);
  }
  return check_sessiond_is_ready();
}

void LocalSessionManagerHandlerImpl::CreateSession(
    ServerContext* context, const LocalCreateSessionRequest* request,
    std::function<void(Status, LocalCreateSessionResponse)> response_callback) {
  set_sentry_transaction("CreateSession");
  auto& request_cpy = *request;
  PrintGrpcMessage(static_cast<const google::protobuf::Message&>(request_cpy));
  enforcer_->get_event_base().runInEventBaseThread([this, response_callback,
                                                    request_cpy]() {
    SessionConfig cfg(request_cpy);
    const std::string& imsi = cfg.get_imsi();
    const CommonSessionContext common = cfg.common_context;

    log_create_session(cfg);
    auto status = validate_create_session_request(cfg);
    if (!status.ok()) {
      send_local_create_session_response(status, "", response_callback);
      return;
    }

    const auto& session_id = id_gen_.gen_session_id(imsi);
    auto session_map = session_store_.read_sessions({imsi});
    SessionActionOrStatus action;
    switch (common.rat_type()) {
      case TGPP_WLAN:
        action = handle_create_session_cwf(session_map, session_id, cfg);
        break;
      case TGPP_LTE:
        action = handle_create_session_lte(session_map, session_id, cfg);
        break;
      default:
        // this should be handled above
        MLOG(MERROR) << "RAT type " << common.rat_type() << " is not supported";
        return;
    }
    if (action.end_existing_session) {
      // Assumption here is that sessions are unique by IMSI+APN
      end_session(
          session_map, common.sid(), common.apn(),
          [&](grpc::Status status, LocalEndSessionResponse response) {});
    }
    if (action.create_new_session) {
      bool success =
          initialize_session(session_map, action.session_id_to_send_back, cfg);
      if (success) {
        send_create_session(session_map, action.session_id_to_send_back, cfg,
                            response_callback);
      } else {
        // abort and send back a failure response to access
        action.set_status(
            Status(grpc::ABORTED, "Failed to update SessionStore"));
      }
    }
    if (action.status_back_to_access) {
      send_local_create_session_response(*(action.status_back_to_access),
                                         action.session_id_to_send_back,
                                         response_callback);
    }
  });
}

bool LocalSessionManagerHandlerImpl::initialize_session(
    SessionMap& session_map, const std::string& session_id,
    const SessionConfig& cfg) {
  MLOG(MINFO) << "Initializing " << session_id
              << " in SessionStore before sending a "
                 "CreateSessionRequest";
  const std::string& imsi = cfg.get_imsi();

  session_map[imsi].push_back(
      std::move(enforcer_->create_initializing_session(session_id, cfg)));
  return session_store_.create_sessions(imsi, std::move(session_map[imsi]));
}

// This function must be called in the event base thread
void LocalSessionManagerHandlerImpl::send_create_session(
    SessionMap& session_map, const std::string& session_id,
    const SessionConfig& cfg,
    std::function<void(grpc::Status, LocalCreateSessionResponse)> cb) {
  auto create_req = make_create_session_request(
      cfg, session_id, enforcer_->get_access_timezone());

  MLOG(MINFO) << "Sending a CreateSessionRequest to fetch policies for "
              << session_id;
  reporter_->report_create_session(
      create_req,
      [this, session_id, cfg, cb,
       session_map_ptr = std::make_shared<SessionMap>(std::move(session_map))](
          Status status, CreateSessionResponse response) mutable {
        PrintGrpcMessage(
            static_cast<const google::protobuf::Message&>(response));
        const std::string& imsi = cfg.get_imsi();
        MLOG(MINFO) << "Processing CreateSessionResponse for " << session_id;

        auto session_map = session_store_.read_sessions({imsi});
        SessionUpdate update =
            SessionStore::get_default_session_update(session_map);
        SessionSearchCriteria criteria(imsi, IMSI_AND_SESSION_ID, session_id);
        auto session_it = session_store_.find_session(session_map, criteria);
        if (!session_it) {
          MLOG(MWARNING) << "Could not find session " << session_id;
          status = Status(grpc::ABORTED, "Session not found");
          send_local_create_session_response(status, session_id, cb);
          return;
        }

        auto& session = **session_it;
        SessionStateUpdateCriteria* session_uc = &update[imsi][session_id];

        if (status.ok()) {
          enforcer_->update_session_with_policy_response(session, response,
                                                         session_uc);
        } else {
          std::ostringstream failure_stream;
          MLOG(MERROR)
              << "Failed to initialize session in SessionProxy/PolicyDB for "
              << imsi << " APN " << cfg.common_context.apn() << ": "
              << status.error_message();
          events_reporter_->session_create_failure(
              cfg, "Failed to initialize session in SessionProxy/PolicyDB");
          session_uc->is_session_ended = true;
        }

        bool write_success = session_store_.update_sessions(update);
        if (write_success) {
          MLOG(MINFO) << "Successfully initialized " << session_id
                      << " in SessionD after talking to PolicyDB/SessionProxy";
          add_session_to_directory_record(imsi, session_id,
                                          cfg.common_context.msisdn());
        } else {
          MLOG(MINFO) << "Failed to initialize new session " << session_id
                      << " in SessionD for subscriber " << imsi
                      << " due to failure writing to SessionStore."
                      << " An earlier update may have invalidated it.";
          status = Status(grpc::ABORTED,
                          "Failed to write session to SessionD storage");
        }
        send_local_create_session_response(status, session_id, cb);
      });
}

SessionActionOrStatus LocalSessionManagerHandlerImpl::handle_create_session_cwf(
    SessionMap& session_map, const std::string& session_id,
    const SessionConfig& cfg) const {
  auto it = session_map.find(cfg.get_imsi());
  if (it != session_map.end()) {
    for (auto& session : it->second) {
      if (session->is_active()) {
        return recycle_cwf_session(session, cfg, session_map);
      }
    }
  }
  return SessionActionOrStatus::create_new_session_action(session_id);
}

SessionActionOrStatus LocalSessionManagerHandlerImpl::recycle_cwf_session(
    std::unique_ptr<SessionState>& session, const SessionConfig& cfg,
    SessionMap& session_map) const {
  const std::string imsi = cfg.get_imsi();
  const std::string session_id = session->get_session_id();
  // To recycle the session, it has to be active (i.e., not in
  // transition for termination).
  MLOG(MINFO) << "Found an active CWF session with the same IMSI " << imsi
              << "... Recycling the existing session.";
  // Since the new session context could be different from the current one,
  // update the config
  SessionUpdate session_update =
      SessionStore::get_default_session_update(session_map);
  SessionStateUpdateCriteria* session_uc = &session_update[imsi][session_id];
  enforcer_->handle_cwf_roaming(session, cfg, session_uc);
  bool success = session_store_.update_sessions(session_update);
  if (!success) {
    MLOG(MINFO) << "Failed to update session config for " << session_id
                << " in SessionD due to failure writing to SessionStore."
                << " An earlier update may have invalidated it.";
    return SessionActionOrStatus::status(
        Status(grpc::ABORTED, "Failed to update SessionStore"));
  }
  MLOG(MINFO) << "Successfully recycled an existing CWF session " << session_id;
  return SessionActionOrStatus::OK(session->get_session_id());
}

SessionActionOrStatus LocalSessionManagerHandlerImpl::handle_create_session_lte(
    SessionMap& session_map, const std::string& session_id,
    const SessionConfig& cfg) {
  SessionActionOrStatus action;
  const std::string imsi = cfg.get_imsi();
  const std::string apn = cfg.common_context.apn();
  const std::string msisdn = cfg.common_context.msisdn();

  // If there are no existing sessions for the IMSI+APN, just create a new one
  SessionSearchCriteria criteria(imsi, IMSI_AND_APN, apn);
  auto session_it = session_store_.find_session(session_map, criteria);
  if (!session_it) {
    return SessionActionOrStatus::create_new_session_action(session_id);
  }
  auto& session = **session_it;
  const std::string existing_session_id = session->get_session_id();

  if (cfg == session->get_config() && session->is_active()) {
    // To recycle the session, it has to be active (i.e., not in
    // transition for termination) AND it should have the exact same
    // configuration.
    MLOG(MINFO) << "Found an active completely duplicated session with IMSI "
                << imsi << " and APN " << cfg.common_context.apn()
                << ", and same configuration. Recycling the existing session "
                << existing_session_id;
    add_session_to_directory_record(imsi, existing_session_id, msisdn);
    return SessionActionOrStatus::OK(existing_session_id);
  }
  // At this point, we have session with same IMSI, but not identical config
  if (session->get_config().common_context.apn() == apn &&
      session->is_active()) {
    // If we have found an active session with the same IMSI+APN, but NOT
    // identical context, we should terminate the existing session.
    MLOG(MINFO) << "Found an active session " << existing_session_id
                << " with same IMSI " << imsi << " and APN " << apn
                << ", but different configuration. Ending the existing "
                << "session, and requesting a new session";
    action.set_end_existing_session();
  }
  action.set_create_new_session(session_id);
  return action;  // Return early
}

void LocalSessionManagerHandlerImpl::send_local_create_session_response(
    Status status, const std::string& session_id,
    std::function<void(Status, LocalCreateSessionResponse)> response_callback) {
  LocalCreateSessionResponse resp;
  resp.set_session_id(session_id);
  PrintGrpcMessage(static_cast<const google::protobuf::Message&>(resp));
  try {
    response_callback(status, resp);
  } catch (...) {
    std::exception_ptr ep = std::current_exception();
    MLOG(MERROR) << "CreateSession response_callback exception: "
                 << (ep ? ep.__cxa_exception_type()->name() : "<unknown");
  }
}

void LocalSessionManagerHandlerImpl::add_session_to_directory_record(
    const std::string& imsi, const std::string& session_id,
    const std::string& msisdn) {
  UpdateRecordRequest request;
  request.set_id(imsi);
  auto update_fields = request.mutable_fields();
  std::string session_id_key = "session_id";
  update_fields->insert({session_id_key, session_id});
  std::string msisdn_id_key = "msisdn";
  update_fields->insert({msisdn_id_key, msisdn});
  MLOG(MINFO) << "Sending a request to DirectoryD to register sessiond ID: "
              << session_id << " msisdn: " << msisdn;
  directoryd_client_->update_directoryd_record(
      request, [imsi](Status status, Void) {
        if (!status.ok()) {
          MLOG(MERROR) << "Could not add session_id to directory record for "
                          "subscriber "
                       << imsi << "; " << status.error_message();
        }
      });
}

/**
 * EndSession completes the entire termination procedure with the OCS & PCRF.
 * The process for session termination is as follows:
 * 1) Start termination process. Enforcer sends delete flow request to Pipelined
 * 2) Enforcer continues to collect usages until its flows are no longer
 *    included in the report (flow deleted in Pipelined) or a specified timeout
 * 3) Asynchronously report usages to cloud in termination requests to
 *    OCS & PCRF
 * 4) Remove the terminated session from being tracked locally, no matter cloud
 *    termination succeeds or not
 */
void LocalSessionManagerHandlerImpl::EndSession(
    ServerContext* context, const LocalEndSessionRequest* request,
    std::function<void(Status, LocalEndSessionResponse)> response_callback) {
  set_sentry_transaction("EndSession");
  auto& request_cpy = *request;
  auto& sid = request->sid();
  auto& apn = request->apn();
  PrintGrpcMessage(static_cast<const google::protobuf::Message&>(request_cpy));
  enforcer_->get_event_base().runInEventBaseThread(
      [this, sid, apn, response_callback]() {
        auto session_map = session_store_.read_sessions({sid.id()});
        MLOG(MINFO) << "Received a termination request from Access for "
                    << sid.id() << " apn " << apn;
        end_session(session_map, sid, apn, response_callback);
      });
}

void LocalSessionManagerHandlerImpl::end_session(
    SessionMap& session_map, const SubscriberID& sid, const std::string& apn,
    std::function<void(Status, LocalEndSessionResponse)> response_callback) {
  auto update = SessionStore::get_default_session_update(session_map);
  auto found = enforcer_->handle_termination_from_access(session_map, sid.id(),
                                                         apn, update);
  if (!found) {
    MLOG(MERROR) << "Failed to find session to terminate for subscriber "
                 << sid.id() << " apn " << apn;
    Status status(grpc::FAILED_PRECONDITION, "Session not found");
    response_callback(status, LocalEndSessionResponse());
    return;
  }
  bool update_success = session_store_.update_sessions(update);
  if (!update_success) {
    auto status =
        Status(grpc::ABORTED,
               "EndSession no longer valid due to another update that "
               "occurred to the session first.");
    response_callback(status, LocalEndSessionResponse());
    return;
  }
  // Success
  response_callback(Status::OK, LocalEndSessionResponse());
}

void LocalSessionManagerHandlerImpl::BindPolicy2Bearer(
    ServerContext* context, const PolicyBearerBindingRequest* request,
    std::function<void(Status, PolicyBearerBindingResponse)>
        response_callback) {
  set_sentry_transaction("BindPolicy2Bearer");
  auto& request_cpy = *request;
  PrintGrpcMessage(static_cast<const google::protobuf::Message&>(request_cpy));
  MLOG(MINFO) << "Received a BindPolicy2Bearer request for "
              << request->sid().id()
              << " with default bearerID: " << request->linked_bearer_id()
              << " policyID: " << request->policy_rule_id()
              << " created dedicated bearerID: " << request->bearer_id()
              << " agw TEID: " << request->teids().agw_teid()
              << " eNB TEID: " << request->teids().enb_teid();
  enforcer_->get_event_base().runInEventBaseThread([this, request_cpy]() {
    auto session_map = session_store_.read_sessions({request_cpy.sid().id()});
    SessionUpdate update =
        SessionStore::get_default_session_update(session_map);
    auto success =
        enforcer_->bind_policy_to_bearer(session_map, request_cpy, update);
    if (!success) {
      MLOG(MDEBUG) << "Failed to process policy -> bearer binding for "
                   << request_cpy.policy_rule_id();
    }
    auto update_success = session_store_.update_sessions(update);
    if (update_success) {
      MLOG(MDEBUG) << "Succeeded in updating SessionStore after processing "
                      "policy->bearer mapping";
    } else {
      MLOG(MERROR) << "Failed in updating SessionStore after processing "
                      "policy->bearer mapping";
    }
  });
  response_callback(Status::OK, PolicyBearerBindingResponse());
}

void LocalSessionManagerHandlerImpl::UpdateTunnelIds(
    ServerContext* context, UpdateTunnelIdsRequest* request,
    std::function<void(Status, UpdateTunnelIdsResponse)> response_callback) {
  set_sentry_transaction("UpdateTunnelIds");
  auto& request_cpy = *request;
  auto imsi = request->sid().id();
  PrintGrpcMessage(static_cast<const google::protobuf::Message&>(request_cpy));
  MLOG(MINFO) << "Received a UpdateTunnelIds request for " << imsi
              << " with default bearer id: " << request->bearer_id()
              << " enb_teid= " << request->enb_teid()
              << " agw_teid= " << request->agw_teid();
  enforcer_->get_event_base().runInEventBaseThread([this, request_cpy, imsi,
                                                    response_callback]() {
    auto session_map = session_store_.read_sessions({imsi});
    auto success = enforcer_->update_tunnel_ids(session_map, request_cpy);
    if (!success) {
      MLOG(MDEBUG) << "Failed to UpdateTunnelIds for imsi " << imsi
                   << " and bearer " << request_cpy.bearer_id();
      auto err_status = Status(grpc::ABORTED, "Failed to Update tunnels Ids");
      response_callback(err_status, UpdateTunnelIdsResponse());
      return;
    }
    bool update_success =
        session_store_.raw_write_sessions(std::move(session_map));
    if (!update_success) {
      MLOG(MERROR) << "Failed in updating SessionStore after processing "
                      "UpdateTunnelIds";
      auto err_status = Status(grpc::ABORTED, "Failed to store tunnels Ids");
      response_callback(err_status, UpdateTunnelIdsResponse());
      return;
    }
    response_callback(Status::OK, UpdateTunnelIdsResponse());
  });
}

void LocalSessionManagerHandlerImpl::SetSessionRules(
    ServerContext* context, const SessionRules* request,
    std::function<void(Status, Void)> response_callback) {
  set_sentry_transaction("SetSessionRules");
  auto& request_cpy = *request;
  PrintGrpcMessage(static_cast<const google::protobuf::Message&>(request_cpy));
  MLOG(MDEBUG) << "Received session <-> rule associations";

  enforcer_->get_event_base().runInEventBaseThread([this, request_cpy]() {
    SessionRead req = {};
    for (const auto& rule_sets : request_cpy.rules_per_subscriber()) {
      req.insert(rule_sets.imsi());
    }
    auto session_map = session_store_.read_sessions(req);
    SessionUpdate update =
        SessionStore::get_default_session_update(session_map);
    enforcer_->handle_set_session_rules(session_map, request_cpy, update);
    auto update_success = session_store_.update_sessions(update);
    if (update_success) {
      MLOG(MDEBUG) << "Succeeded in updating SessionStore after processing "
                      "session rules set";
    } else {
      MLOG(MERROR) << "Failed in updating SessionStore after processing "
                      "session rules set";
    }
  });
  response_callback(Status::OK, Void());
}

void LocalSessionManagerHandlerImpl::log_create_session(
    const SessionConfig& cfg) {
  const std::string& imsi = cfg.get_imsi();
  const auto& apn = cfg.common_context.apn();
  std::string create_message =
      "Received a LocalCreateSessionRequest for " + imsi + " with APN:" + apn;
  if (cfg.rat_specific_context.has_lte_context()) {
    const auto& lte = cfg.rat_specific_context.lte_context();
    create_message +=
        ", default bearer ID:" + std::to_string(lte.bearer_id()) +
        ", PLMN ID:" + lte.plmn_id() + ", IMSI PLMN ID:" + lte.imsi_plmn_id() +
        ", User location:" + magma::bytes_to_hex(lte.user_location());
  } else if (cfg.rat_specific_context.has_wlan_context()) {
    create_message +=
        ", MAC addr:" + cfg.rat_specific_context.wlan_context().mac_addr();
  }
  MLOG(MINFO) << create_message;
}

}  // namespace magma

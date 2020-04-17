/**
 * Copyright (c) 2016-present, Facebook, Inc.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */
#include <chrono>
#include <thread>

#include <google/protobuf/util/time_util.h>

#include "LocalSessionManagerHandler.h"
#include "magma_logging.h"

using grpc::Status;

namespace magma {

const std::string LocalSessionManagerHandlerImpl::hex_digit_ =
    "0123456789abcdef";

LocalSessionManagerHandlerImpl::LocalSessionManagerHandlerImpl(
    std::shared_ptr<LocalEnforcer> enforcer, SessionReporter* reporter,
    std::shared_ptr<AsyncDirectorydClient> directoryd_client,
    SessionStore& session_store)
    : enforcer_(enforcer),
      session_store_(session_store),
      reporter_(reporter),
      directoryd_client_(directoryd_client),

      current_epoch_(0),
      reported_epoch_(0),
      retry_timeout_(1) {}

void LocalSessionManagerHandlerImpl::ReportRuleStats(
    ServerContext* context, const RuleRecordTable* request,
    std::function<void(Status, Void)> response_callback) {
  auto& request_cpy = *request;
  if (request_cpy.records_size() > 0) {
    MLOG(MDEBUG) << "Aggregating " << request_cpy.records_size() << " records";

    enforcer_->get_event_base().runInEventBaseThread([this, request_cpy]() {
      auto session_map = get_sessions_for_reporting(request_cpy);
      SessionUpdate update =
          SessionStore::get_default_session_update(session_map);
      enforcer_->aggregate_records(session_map, request_cpy, update);
      check_usage_for_reporting(std::move(session_map), update);
    });
  }
  reported_epoch_ = request_cpy.epoch();
  if (is_pipelined_restarted()) {
    MLOG(MINFO) << "Pipelined has been restarted, attempting to sync flows";
    restart_pipelined(reported_epoch_);
    // Set the current epoch right away to prevent double setup call requests
    current_epoch_ = reported_epoch_;
  }
  response_callback(Status::OK, Void());
}

void LocalSessionManagerHandlerImpl::check_usage_for_reporting(
    SessionMap session_map, SessionUpdate& session_update) {
  std::vector<std::unique_ptr<ServiceAction>> actions;
  auto request =
      enforcer_->collect_updates(session_map, actions, session_update);
  enforcer_->execute_actions(session_map, actions, session_update);
  if (request.updates_size() == 0 && request.usage_monitors_size() == 0) {
    auto update_success = session_store_.update_sessions(session_update);
    if (update_success) {
      MLOG(MERROR) << "Succeeded in updating session after no reporting";
    } else {
      MLOG(MERROR) << "Failed in updating session after no reporting";
    }
    return;  // nothing to report
  }
  MLOG(MINFO) << "Sending " << request.updates_size()
              << " charging updates and " << request.usage_monitors_size()
              << " monitor updates to OCS and PCRF";

  // report to cloud
  // NOTE: It is not possible to construct a std::function from a move-only type
  //       So because of this, we can't directly move session_map into the
  //       value-capture of the callback. As a workaround, a shared_ptr to
  //       the session_map is used.
  //       Check
  //       https://stackoverflow.com/questions/25421346/how-to-create-an-stdfunction-from-a-move-capturing-lambda-expression
  reporter_->report_updates(
      request,
      [this, request,
       session_map_ptr = std::make_shared<SessionMap>(std::move(session_map)),
       session_update  = std::move(session_update)](
          Status status, UpdateSessionResponse response) mutable {
        if (!status.ok()) {
          MLOG(MERROR) << "Update of size " << request.updates_size()
                       << " to OCS failed entirely: " << status.error_message();
        } else {
          enforcer_->update_session_credits_and_rules(
              *session_map_ptr, response, session_update);
          session_store_.update_sessions(session_update);
        }
      });
}

bool LocalSessionManagerHandlerImpl::is_pipelined_restarted() {
  // If 0 also setup pipelined because it always waits for setup instructions
  if (current_epoch_ == 0 || current_epoch_ != reported_epoch_) {
    return true;
  }
  return false;
}

void LocalSessionManagerHandlerImpl::handle_setup_callback(
    const std::uint64_t& epoch, Status status, SetupFlowsResult resp) {
  using namespace std::placeholders;
  if (!status.ok()) {
    MLOG(MERROR) << "Could not setup pipelined, rpc failed with: "
                 << status.error_message() << ", retrying pipelined setup.";

    enforcer_->get_event_base().runInEventBaseThread([=] {
      enforcer_->get_event_base().timer().scheduleTimeoutFn(
          std::move([=] {
            auto session_map = session_store_.read_all_sessions();
            enforcer_->setup(
                session_map, epoch,
                std::bind(
                    &LocalSessionManagerHandlerImpl::handle_setup_callback,
                    this, epoch, _1, _2));
          }),
          retry_timeout_);
    });
  }

  if (resp.result() == resp.OUTDATED_EPOCH) {
    MLOG(MWARNING) << "Pipelined setup call has outdated epoch, abandoning.";
  } else if (resp.result() == resp.FAILURE) {
    MLOG(MWARNING) << "Pipelined setup failed, retrying pipelined setup "
                      "for epoch "
                   << epoch;
    enforcer_->get_event_base().runInEventBaseThread([=] {
      enforcer_->get_event_base().timer().scheduleTimeoutFn(
          std::move([=] {
            auto session_map = session_store_.read_all_sessions();
            enforcer_->setup(
                session_map, epoch,
                std::bind(
                    &LocalSessionManagerHandlerImpl::handle_setup_callback,
                    this, epoch, _1, _2));
          }),
          retry_timeout_);
    });
  } else {
    MLOG(MDEBUG) << "Successfully setup pipelined.";
  }
}

bool LocalSessionManagerHandlerImpl::restart_pipelined(
    const std::uint64_t& epoch) {
  using namespace std::placeholders;
  enforcer_->get_event_base().runInEventBaseThread([this, epoch]() {
    auto session_map = session_store_.read_all_sessions();
    enforcer_->setup(
        session_map, epoch,
        std::bind(
            &LocalSessionManagerHandlerImpl::handle_setup_callback, this, epoch,
            _1, _2));
  });
  return true;
}

static CreateSessionRequest copy_session_info2create_req(
    const LocalCreateSessionRequest& request, const std::string& sid) {
  CreateSessionRequest create_request;

  create_request.mutable_subscriber()->CopyFrom(request.sid());
  create_request.set_session_id(sid);
  create_request.set_ue_ipv4(request.ue_ipv4());
  create_request.set_spgw_ipv4(request.spgw_ipv4());
  create_request.set_apn(request.apn());
  create_request.set_msisdn(request.msisdn());
  create_request.set_imei(request.imei());
  create_request.set_plmn_id(request.plmn_id());
  create_request.set_imsi_plmn_id(request.imsi_plmn_id());
  create_request.set_user_location(request.user_location());
  create_request.set_hardware_addr(request.hardware_addr());
  create_request.set_rat_type(request.rat_type());
  if (request.has_qos_info()) {
    create_request.mutable_qos_info()->CopyFrom(request.qos_info());
  }

  return create_request;
}

void LocalSessionManagerHandlerImpl::CreateSession(
    ServerContext* context, const LocalCreateSessionRequest* request,
    std::function<void(Status, LocalCreateSessionResponse)> response_callback) {
  auto& request_cpy = *request;
  enforcer_->get_event_base().runInEventBaseThread([this, context,
                                                    response_callback,
                                                    request_cpy]() {
    auto imsi     = request_cpy.sid().id();
    auto sid      = id_gen_.gen_session_id(imsi);
    auto mac_addr = convert_mac_addr_to_str(request_cpy.hardware_addr());
    MLOG(MDEBUG) << "PLMN_ID: " << request_cpy.plmn_id()
                 << " IMSI_PLMN_ID: " << request_cpy.imsi_plmn_id();

    SessionConfig cfg = {.ue_ipv4           = request_cpy.ue_ipv4(),
                         .spgw_ipv4         = request_cpy.spgw_ipv4(),
                         .msisdn            = request_cpy.msisdn(),
                         .apn               = request_cpy.apn(),
                         .imei              = request_cpy.imei(),
                         .plmn_id           = request_cpy.plmn_id(),
                         .imsi_plmn_id      = request_cpy.imsi_plmn_id(),
                         .user_location     = request_cpy.user_location(),
                         .rat_type          = request_cpy.rat_type(),
                         .mac_addr          = mac_addr,
                         .hardware_addr     = request_cpy.hardware_addr(),
                         .radius_session_id = request_cpy.radius_session_id(),
                         .bearer_id         = request_cpy.bearer_id()};

    QoSInfo qos_info = {.enabled = request_cpy.has_qos_info()};
    if (request_cpy.has_qos_info()) {
      qos_info.qci = request_cpy.qos_info().qos_class_id();
    }
    cfg.qos_info = qos_info;

    auto session_map = get_sessions_for_creation(request_cpy);
    if (enforcer_->session_with_imsi_exists(session_map, imsi)) {
      std::string core_sid;

      // For LTE case, load session if and only if the configuration exactly
      // matches. For CWF use case, we can recycle any active session
      bool same_config = false;
      bool is_active   = false;
      bool is_wifi     = request_cpy.rat_type() == RATType::TGPP_WLAN;
      if (is_wifi) {
        is_active = enforcer_->has_active_session(session_map, imsi, &core_sid);
      } else {
        same_config = enforcer_->session_with_same_config_exists(
            session_map, imsi, cfg, &core_sid);
        is_active = enforcer_->is_session_active(session_map, imsi, core_sid);
      }
      // To recycle the session, it has to be active (i.e., not in transition
      // for termination), it should have the exact same configuration or it
      // should be CWF use case.
      if ((same_config || is_wifi) && is_active) {
        Status status;
        if (is_wifi) {
          MLOG(MINFO) << "Found an active session with the same IMSI " << imsi
                      << " and RAT Type is WLAN, not creating a new session";
          // Wifi only supports one session per subscriber, so update the config
          // here
          SessionUpdate session_update =
              SessionStore::get_default_session_update(session_map);
          enforcer_->handle_cwf_roaming(session_map, imsi, cfg, session_update);
          if (session_store_.update_sessions(session_update)) {
            MLOG(MINFO) << "Successfully updated session " << sid
                        << " in sessiond for subscriber " << imsi;
            status = grpc::Status::OK;
          } else {
            MLOG(MINFO) << "Failed to initialize new session " << sid
                        << " in sessiond for subscriber " << imsi
                        << " with default bearer id " << cfg.bearer_id
                        << " due to failure writing to SessionStore."
                        << " An earlier update may have invalidated it.";
            status = Status(grpc::ABORTED, "Failed to update SessionStore");
          }
        } else {
          MLOG(MINFO) << "Found completely duplicated session with IMSI "
                      << imsi << " and APN " << request_cpy.apn()
                      << ", not creating session";
          status = grpc::Status::OK;
        }
        try {
          LocalCreateSessionResponse resp;
          resp.set_session_id(core_sid);
          response_callback(status, resp);
        } catch (...) {
          std::exception_ptr ep = std::current_exception();
          MLOG(MERROR) << "CreateSession response_callback exception: "
                       << (ep ? ep.__cxa_exception_type()->name() : "<unknown");
        }
        // No new session created
        return;
      }

      if (!enforcer_->session_with_apn_exists(
              session_map, imsi, request_cpy.apn())) {
        MLOG(MINFO) << "Found session with the same IMSI " << imsi
                    << " but different APN " << request_cpy.apn()
                    << ", will request a new session from PCRF/PCF";
      } else if (is_active) {
        MLOG(MINFO) << "Found an active session with the same IMSI " << imsi
                    << " and APN " << request_cpy.apn()
                    << ", but different configuration."
                    << " Ending the existing session, "
                    << "will request a new session from PCRF/PCF";
        LocalEndSessionRequest end_session_req;
        end_session_req.mutable_sid()->CopyFrom(request_cpy.sid());
        end_session_req.set_apn(request_cpy.apn());
        end_session(
            session_map, end_session_req,
            [&](grpc::Status status, LocalEndSessionResponse response) {});
      } else {
        MLOG(MINFO) << "Found a session in termination with the same IMSI "
                    << imsi << " and same APN " << request_cpy.apn()
                    << ", will request a new session from PCRF/PCF";
      }
    }
    send_create_session(
        session_map, copy_session_info2create_req(request_cpy, sid), imsi, sid,
        cfg, response_callback);
  });
}

void LocalSessionManagerHandlerImpl::send_create_session(
    SessionMap& session_map, const CreateSessionRequest& request,
    const std::string& imsi, const std::string& sid, const SessionConfig& cfg,
    std::function<void(grpc::Status, LocalCreateSessionResponse)>
        response_callback) {
  reporter_->report_create_session(
      request,
      [this, imsi, sid, cfg, response_callback,
       session_map_ptr = std::make_shared<SessionMap>(std::move(session_map))](
          Status status, CreateSessionResponse response) mutable {
        if (status.ok()) {
          bool success = enforcer_->init_session_credit(
              *session_map_ptr, imsi, sid, cfg, response);
          if (!success) {
            MLOG(MERROR) << "Failed to initialize session for IMSI " << imsi;
            status = Status(
                grpc::FAILED_PRECONDITION, "Failed to initialize session");
          } else {
            bool write_success = session_store_.create_sessions(
                imsi, std::move((*session_map_ptr)[imsi]));
            if (write_success) {
              MLOG(MINFO) << "Successfully initialized new session " << sid
                          << " in sessiond for subscriber " << imsi
                          << " with default bearer id " << cfg.bearer_id;
              add_session_to_directory_record(imsi, sid);
            } else {
              MLOG(MINFO) << "Failed to initialize new session " << sid
                          << " in sessiond for subscriber " << imsi
                          << " with default bearer id " << cfg.bearer_id
                          << " due to failure writing to SessionStore."
                          << " An earlier update may have invalidated it.";
              status = Status(
                  grpc::ABORTED,
                  "Failed to write session to sessiond storage.");
            }
          }
        } else {
          MLOG(MERROR) << "Failed to initialize session in SessionProxy "
                       << "for IMSI " << imsi << ": " << status.error_message();
        }
        LocalCreateSessionResponse resp;
        resp.set_session_id(response.session_id());
        response_callback(status, resp);
      });
}

void LocalSessionManagerHandlerImpl::add_session_to_directory_record(
    const std::string& imsi, const std::string& session_id) {
  UpdateRecordRequest request;
  request.set_id(imsi);
  auto update_fields         = request.mutable_fields();
  std::string session_id_key = "session_id";
  update_fields->insert({session_id_key, session_id});
  directoryd_client_->update_directoryd_record(
      request, [this, imsi](Status status, Void) {
        if (!status.ok()) {
          MLOG(MERROR) << "Could not add session_id to directory record for "
                          "subscriber "
                       << imsi << "; " << status.error_message();
        }
      });
}

std::string LocalSessionManagerHandlerImpl::convert_mac_addr_to_str(
    const std::string& mac_addr) {
  std::string res;
  auto l = mac_addr.length();
  if (l == 0) {
    return res;
  }
  res.reserve(l * 3 - 1);
  for (int i = 0; i < l; i++) {
    if (i > 0) {
      res.push_back(':');
    }
    unsigned char c = mac_addr[i];
    res.push_back(hex_digit_[c >> 4]);
    res.push_back(hex_digit_[c & 0x0F]);
  }
  return res;
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
  auto& request_cpy = *request;
  enforcer_->get_event_base().runInEventBaseThread(
      [this, request_cpy = std::move(request_cpy), response_callback]() {
        auto session_map = get_sessions_for_deletion(request_cpy);
        end_session(session_map, request_cpy, response_callback);
      });
}

void LocalSessionManagerHandlerImpl::end_session(
    SessionMap& session_map, const LocalEndSessionRequest& request,
    std::function<void(Status, LocalEndSessionResponse)> response_callback) {
  try {
    auto reporter = reporter_;
    auto update   = SessionStore::get_default_session_update(session_map);
    enforcer_->terminate_subscriber(
        session_map, request.sid().id(), request.apn(), update);

    bool update_success = session_store_.update_sessions(update);
    if (update_success) {
      response_callback(Status::OK, LocalEndSessionResponse());
    } else {
      auto status = Status(
          grpc::ABORTED,
          "EndSession no longer valid due to another update that "
          "occurred to the session first.");
      response_callback(status, LocalEndSessionResponse());
    }
  } catch (const SessionNotFound& ex) {
    MLOG(MERROR) << "Failed to find session to terminate for subscriber "
                 << request.sid().id();
    Status status(grpc::FAILED_PRECONDITION, "Session not found");
    response_callback(status, LocalEndSessionResponse());
  }
}

SessionMap LocalSessionManagerHandlerImpl::get_sessions_for_creation(
    const LocalCreateSessionRequest& request) {
  SessionRead req = {request.sid().id()};
  return session_store_.read_sessions(req);
}

SessionMap LocalSessionManagerHandlerImpl::get_sessions_for_reporting(
    const RuleRecordTable& records) {
  SessionRead req = {};
  for (const RuleRecord& record : records.records()) {
    req.insert(record.sid());
  }
  return session_store_.read_sessions_for_reporting(req);
}

SessionMap LocalSessionManagerHandlerImpl::get_sessions_for_deletion(
    const LocalEndSessionRequest& request) {
  // Request numbers are not incremented here, but rather in the 5s delayed
  // callback when reporting is done.
  SessionRead req = {request.sid().id()};
  return session_store_.read_sessions(req);
}

}  // namespace magma

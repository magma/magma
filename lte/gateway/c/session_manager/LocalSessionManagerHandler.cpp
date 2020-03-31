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
  std::shared_ptr<LocalEnforcer> enforcer,
  SessionReporter* reporter,
  std::shared_ptr<AsyncDirectorydClient> directoryd_client,
  SessionMap & session_map):
  enforcer_(enforcer),
  session_map_(session_map),
  reporter_(reporter),
  directoryd_client_(directoryd_client),

  current_epoch_(0),
  reported_epoch_(0),
  retry_timeout_(1)
{
}

void LocalSessionManagerHandlerImpl::ReportRuleStats(
  ServerContext* context,
  const RuleRecordTable* request,
  std::function<void(Status, Void)> response_callback)
{
  auto &request_cpy = *request;
  MLOG(MDEBUG) << "Aggregating " << request_cpy.records_size() << " records";
  enforcer_->get_event_base().runInEventBaseThread([this, request_cpy]() {
    SessionUpdate update = SessionStore::get_default_session_update(session_map_);
    enforcer_->aggregate_records(session_map_, request_cpy, update);
    check_usage_for_reporting(update);
  });
  reported_epoch_ = request_cpy.epoch();
  if (is_pipelined_restarted()) {
    MLOG(MDEBUG) << "Pipelined has been restarted, attempting to sync flows";
    restart_pipelined(reported_epoch_);
    // Set the current epoch right away to prevent double setup call requests
    current_epoch_ = reported_epoch_;
  }
  response_callback(Status::OK, Void());
}

void LocalSessionManagerHandlerImpl::check_usage_for_reporting(SessionUpdate& session_update)
{
  std::vector<std::unique_ptr<ServiceAction>> actions;
  auto request = enforcer_->collect_updates(session_map_, actions, session_update);
  enforcer_->execute_actions(session_map_, actions, session_update);
  if (request.updates_size() == 0 && request.usage_monitors_size() == 0) {
    // TODO: Save updates back into the SessionStore
    return; // nothing to report
  }
  MLOG(MDEBUG) << "Sending " << request.updates_size()
               << " charging updates and " << request.usage_monitors_size()
               << " monitor updates to OCS and PCRF";

  // report to cloud
  reporter_->report_updates(
    request, [this, request, &session_update](Status status, UpdateSessionResponse response) {
      if (!status.ok()) {
        enforcer_->reset_updates(session_map_, request);
        MLOG(MERROR) << "Update of size " << request.updates_size()
                     << " to OCS failed entirely: " << status.error_message();
      } else {
        MLOG(MDEBUG) << "Received updated responses from OCS and PCRF";
        enforcer_->update_session_credits_and_rules(session_map_, response, session_update);
        // Check if we need to report more updates
        check_usage_for_reporting(session_update);
      }
    });
}

bool LocalSessionManagerHandlerImpl::is_pipelined_restarted()
{
  // If 0 also setup pipelined because it always waits for setup instructions
  if (current_epoch_ == 0 || current_epoch_ != reported_epoch_) {
    return true;
  }
  return false;
}

void LocalSessionManagerHandlerImpl::handle_setup_callback(
  const std::uint64_t& epoch,
  Status status,
  SetupFlowsResult resp)
{
  using namespace std::placeholders;
  if (!status.ok()) {
    MLOG(MERROR) << "Could not setup pipelined, rpc failed with: "
                 << status.error_message() << ", retrying pipelined setup.";

    enforcer_->get_event_base().runInEventBaseThread([=] {
      enforcer_->get_event_base().timer().scheduleTimeoutFn(
        std::move([=] {
          enforcer_->setup(
            session_map_,
            epoch,
            std::bind(
              &LocalSessionManagerHandlerImpl::handle_setup_callback,
              this,
              epoch,
              _1,
              _2));
        }),
        retry_timeout_);
    });
  }

  if (resp.result() == resp.OUTDATED_EPOCH) {
    MLOG(MWARNING) << "Pipelined setup call has outdated epoch, abandoning.";
  } else if (resp.result() == resp.FAILURE) {
    MLOG(MWARNING) << "Pipelined setup failed, retrying pipelined setup "
                      "for epoch " << epoch;
    enforcer_->get_event_base().runInEventBaseThread([=] {
      enforcer_->get_event_base().timer().scheduleTimeoutFn(
        std::move([=] {
          enforcer_->setup(
            session_map_,
            epoch,
            std::bind(
              &LocalSessionManagerHandlerImpl::handle_setup_callback,
              this,
              epoch,
              _1,
              _2));
        }),
        retry_timeout_);
    });
  } else {
    MLOG(MDEBUG) << "Successfully setup pipelined.";
  }
}

bool LocalSessionManagerHandlerImpl::restart_pipelined(
  const std::uint64_t& epoch)
{
  using namespace std::placeholders;
  enforcer_->get_event_base().runInEventBaseThread([this, epoch]() {
    enforcer_->setup(
      session_map_,
      epoch,
      std::bind(
        &LocalSessionManagerHandlerImpl::handle_setup_callback,
        this,
        epoch,
        _1,
        _2));
  });
  return true;
}

static CreateSessionRequest copy_session_info2create_req(
  const LocalCreateSessionRequest& request,
  const std::string& sid)
{
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
  ServerContext* context,
  const LocalCreateSessionRequest* request,
  std::function<void(Status, LocalCreateSessionResponse)> response_callback)
{
  auto imsi = request->sid().id();
  auto sid = id_gen_.gen_session_id(imsi);
  auto mac_addr = convert_mac_addr_to_str(request->hardware_addr());
  MLOG(MDEBUG) << "PLMN_ID: " << request->plmn_id()
              << " IMSI_PLMN_ID: " << request->imsi_plmn_id();

  SessionState::Config cfg = {.ue_ipv4 = request->ue_ipv4(),
                              .spgw_ipv4 = request->spgw_ipv4(),
                              .msisdn = request->msisdn(),
                              .apn = request->apn(),
                              .imei = request->imei(),
                              .plmn_id = request->plmn_id(),
                              .imsi_plmn_id = request->imsi_plmn_id(),
                              .user_location = request->user_location(),
                              .rat_type = request->rat_type(),
                              .mac_addr = mac_addr,
                              .hardware_addr = request->hardware_addr(),
                              .radius_session_id = request->radius_session_id(),
                              .bearer_id = request->bearer_id()};

  SessionState::QoSInfo qos_info = {.enabled = request->has_qos_info()};
  if (request->has_qos_info()) {
    qos_info.qci = request->qos_info().qos_class_id();
  }
  cfg.qos_info = qos_info;

  if (enforcer_->session_with_imsi_exists(session_map_, imsi)) {
    std::string core_sid;
    bool same_config = enforcer_->session_with_same_config_exists(
      session_map_, imsi, cfg, &core_sid);
    bool is_wifi = request->rat_type() == RATType::TGPP_WLAN;
    if (same_config || is_wifi){
      if (is_wifi) {
        MLOG(MINFO) << "Found a session with the same IMSI " << imsi
                    << " and RAT Type is WLAN, not creating a new session";
      } else {
        MLOG(MINFO) << "Found completely duplicated session with IMSI " << imsi
                    << " and APN " << request->apn()
                    << ", not creating session";
      }
      enforcer_->get_event_base().runInEventBaseThread(
          [response_callback, core_sid]() {
            try {
              LocalCreateSessionResponse resp;
              resp.set_session_id(core_sid);
              response_callback(grpc::Status::OK, resp);
            } catch (...) {
                std::exception_ptr ep = std::current_exception();
                MLOG(MERROR) << "CreateSession response_callback exception: "
                             << (ep ? ep.__cxa_exception_type()->name()
                                  : "<unknown");
            }
          }
      );
      // No new session created
      return;
    }
    if (enforcer_->session_with_apn_exists(
          session_map_, imsi, request->apn())) {
      MLOG(MINFO) << "Found session with the same IMSI " << imsi
                  << " and APN " << request->apn()
                  << ", but different configuration."
                  << " Ending the existing session";
      LocalEndSessionRequest end_session_req;
      end_session_req.mutable_sid()->CopyFrom(request->sid());
      end_session_req.set_apn(request->apn());
      EndSession(
        context,
        &end_session_req,
        [&](grpc::Status status, LocalEndSessionResponse response) { return; });
    } else {
      MLOG(MINFO) << "Found session with the same IMSI " << imsi
                  << " but different APN " << request->apn()
                  << ", will request a new session from PCRF/PCF";
    }
  }
  send_create_session(
    copy_session_info2create_req(*request, sid),
    imsi, sid, cfg, response_callback);
}

void LocalSessionManagerHandlerImpl::send_create_session(
  const CreateSessionRequest& request,
  const std::string& imsi,
  const std::string& sid,
  const SessionState::Config& cfg,
  std::function<void(grpc::Status, LocalCreateSessionResponse)> response_callback)
{
  reporter_->report_create_session(
    request,
    [this, imsi, sid, cfg, response_callback](
      Status status, CreateSessionResponse response) {
      if (status.ok()) {
        bool success = enforcer_->init_session_credit(
          session_map_, imsi, sid, cfg, response);
        if (!success) {
          MLOG(MERROR) << "Failed to init session in for IMSI " << imsi;
          status =
            Status(
              grpc::FAILED_PRECONDITION, "Failed to initialize session");
        } else {
          MLOG(MINFO) << "Successfully initialized new session " << sid
                      << " in sessiond for subscriber " << imsi
                      << " with default bearer id " << cfg.bearer_id;
          add_session_to_directory_record(imsi, sid);
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
  const std::string& imsi,
  const std::string& session_id)
{
  UpdateRecordRequest request;
  request.set_id(imsi);
  auto update_fields = request.mutable_fields();
  std::string session_id_key = "session_id";
  update_fields->insert({session_id_key, session_id});
  directoryd_client_->update_directoryd_record(request,
    [this, imsi] (Status status, Void) {
    if (!status.ok()) {
      MLOG(MERROR) << "Could not add session_id to directory record for "
      "subscriber " << imsi << "; " << status.error_message();
    }
  });
}

std::string LocalSessionManagerHandlerImpl::convert_mac_addr_to_str(
        const std::string& mac_addr)
{
  std::string res;
  auto l = mac_addr.length();
  if (l == 0) {
      return res;
  }
  res.reserve(l*3-1);
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
  ServerContext* context,
  const LocalEndSessionRequest* request,
  std::function<void(Status, LocalEndSessionResponse)> response_callback)
{
  auto &request_cpy = *request;
  enforcer_->get_event_base().runInEventBaseThread(
    [this, request_cpy, response_callback]() {
      try {
        auto reporter = reporter_;
        SessionStateUpdateCriteria update_criteria = get_default_update_criteria();
        enforcer_->terminate_subscriber(
          session_map_,
          request_cpy.sid().id(),
          request_cpy.apn(),
          [reporter](SessionTerminateRequest term_req) {
            // report to cloud
            auto logging_cb =
              SessionReporter::get_terminate_logging_cb(term_req);
            reporter->report_terminate_session(term_req, logging_cb);
          },
          update_criteria);
        // TODO: Write the delete back into the SessionStore
        response_callback(grpc::Status::OK, LocalEndSessionResponse());
      } catch (const SessionNotFound &ex) {
        MLOG(MERROR) << "Failed to find session to terminate for subscriber "
                     << request_cpy.sid().id();
        Status status(grpc::FAILED_PRECONDITION, "Session not found");
        response_callback(status, LocalEndSessionResponse());
      }
    });
}

} // namespace magma

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

#include "LocalSessionManagerHandler.h"
#include "magma_logging.h"

using grpc::Status;

namespace magma {

LocalSessionManagerHandlerImpl::LocalSessionManagerHandlerImpl(
  LocalEnforcer *enforcer,
  SessionCloudReporter *reporter):
  enforcer_(enforcer),
  reporter_(reporter)
{
}

void LocalSessionManagerHandlerImpl::ReportRuleStats(
  ServerContext *context,
  const RuleRecordTable *request,
  std::function<void(Status, Void)> response_callback)
{
  auto &request_cpy = *request;
  MLOG(MDEBUG) << "Aggregating " << request_cpy.records_size() << " records";
  enforcer_->get_event_base().runInEventBaseThread([this, request_cpy]() {
    enforcer_->aggregate_records(request_cpy);
    check_usage_for_reporting();
  });
  response_callback(Status::OK, Void());
}

void LocalSessionManagerHandlerImpl::check_usage_for_reporting()
{
  auto request = enforcer_->collect_updates();
  if (request.updates_size() == 0 && request.usage_monitors_size() == 0) {
    return; // nothing to report
  }
  MLOG(MDEBUG) << "Sending " << request.updates_size()
               << " charging updates and " << request.usage_monitors_size()
               << " monitor updates to OCS and PCRF";

  // report to cloud
  reporter_->report_updates(
    request, [this, request](Status status, UpdateSessionResponse response) {
      if (!status.ok()) {
        enforcer_->reset_updates(request);
        MLOG(MERROR) << "Update of size " << request.updates_size()
                     << " to OCS failed entirely: " << status.error_message();
      } else {
        MLOG(MDEBUG) << "Received updated responses from OCS and PCRF";
        enforcer_->update_session_credit(response);
        // Check if we need to report more updates
        check_usage_for_reporting();
      }
    });
}

static CreateSessionRequest copy_session_info2create_req(
  const LocalCreateSessionRequest *request,
  const std::string &sid)
{
  CreateSessionRequest create_request;

  create_request.mutable_subscriber()->CopyFrom(request->sid());
  create_request.set_session_id(sid);
  create_request.set_ue_ipv4(request->ue_ipv4());
  create_request.set_spgw_ipv4(request->spgw_ipv4());
  create_request.set_apn(request->apn());
  create_request.set_msisdn(request->msisdn());
  create_request.set_imei(request->imei());
  create_request.set_plmn_id(request->plmn_id());
  create_request.set_imsi_plmn_id(request->imsi_plmn_id());
  create_request.set_user_location(request->user_location());
  create_request.mutable_qos_info()->CopyFrom(request->qos_info());

  return create_request;
}

void LocalSessionManagerHandlerImpl::CreateSession(
  ServerContext *context,
  const LocalCreateSessionRequest *request,
  std::function<void(Status, LocalCreateSessionResponse)> response_callback)
{
  auto imsi = request->sid().id();
  auto sid = id_gen_.gen_session_id(imsi);
  SessionState::Config cfg = {.ue_ipv4 = request->ue_ipv4(),
                              .spgw_ipv4 = request->spgw_ipv4(),
                              .msisdn = request->msisdn(),
                              .apn = request->apn(),
                              .imei = request->imei(),
                              .plmn_id = request->plmn_id(),
                              .imsi_plmn_id = request->imsi_plmn_id(),
                              .user_location = request->user_location()};
  reporter_->report_create_session(
    copy_session_info2create_req(request, sid),
    [this, imsi, sid, cfg, response_callback](
      Status status, CreateSessionResponse response) {
      if (status.ok()) {
        bool success = enforcer_->init_session_credit(imsi, sid, cfg, response);
        if (!success) {
          MLOG(MERROR) << "Failed to init session in Usage Monitor for IMSI "
                       << imsi;
          status =
            Status(grpc::FAILED_PRECONDITION, "Failed to initialize session");
        } else {
          MLOG(MINFO) << "Successfully initialized new session in sessiond "
                      << "for subscriber " << imsi;
        }
      } else {
        MLOG(MERROR) << "Failed to initialize session in OCS for IMSI " << imsi
                     << ": " << status.error_message();
      }
      response_callback(status, LocalCreateSessionResponse());
    });
}

static void report_termination(
  LocalEnforcer &enforcer,
  SessionCloudReporter &reporter,
  const SessionTerminateRequest &term_req,
  std::function<void(Status, LocalEndSessionResponse)> response_callback)
{
  reporter.report_terminate_session(
    term_req,
    [&enforcer, &reporter, term_req, response_callback](
      Status status, SessionTerminateResponse response) {
      if (!status.ok()) {
        MLOG(MERROR) << "Failed to terminate session in controller for "
                        "subscriber "
                     << term_req.sid() << ": " << status.error_message();
      } else {
        MLOG(MDEBUG) << "Termination successful in controller for "
                        "subscriber "
                     << term_req.sid();
      }
      // No matter what, end session locally
      enforcer.complete_termination(term_req.sid(), term_req.session_id());
      response_callback(status, LocalEndSessionResponse());
    });
}

/**
 * EndSession completes the entire termination procedure with the OCS & PCRF.
 * Instead of waiting for the last usage updates, termination is reported
 * immediately. The process for session termination is as follows:
 * 1) Collect usages in current state in session state
 * 2) Report usages to cloud in termination requests to OCS & PCRF
 * 3) Wait for response
 * 4) Remove the terminated session from being tracked
 */
void LocalSessionManagerHandlerImpl::EndSession(
  ServerContext *context,
  const SubscriberID *request,
  std::function<void(Status, LocalEndSessionResponse)> response_callback)
{
  auto &request_cpy = *request;
  enforcer_->get_event_base().runInEventBaseThread(
    [this, request_cpy, response_callback]() {
      try {
        auto term_req = enforcer_->terminate_subscriber(request_cpy.id());
        // report to cloud
        report_termination(*enforcer_, *reporter_, term_req, response_callback);
      } catch (const SessionNotFound &ex) {
        MLOG(MERROR) << "Failed to find session to terminate for subscriber "
                     << request_cpy.id();
        Status status(grpc::FAILED_PRECONDITION, "Session not found");
        response_callback(status, LocalEndSessionResponse());
      }
    });
}

} // namespace magma

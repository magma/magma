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

#include "SessionProxyResponderHandler.h"
#include "magma_logging.h"

using grpc::Status;

namespace magma {
std::string raa_result_to_str(ReAuthResult res);

SessionProxyResponderHandlerImpl::SessionProxyResponderHandlerImpl(
    std::shared_ptr<LocalEnforcer> enforcer, SessionStore& session_store)
    : enforcer_(enforcer), session_store_(session_store) {}

void SessionProxyResponderHandlerImpl::ChargingReAuth(
    ServerContext* context, const ChargingReAuthRequest* request,
    std::function<void(Status, ChargingReAuthAnswer)> response_callback) {
  auto& request_cpy = *request;
  MLOG(MDEBUG) << "Received a Gy (Charging) ReAuthRequest for "
               << request->session_id() << " and charging_key "
               << request->charging_key();
  enforcer_->get_event_base().runInEventBaseThread(
      [this, request_cpy, response_callback]() {
        auto session_map = get_sessions_for_charging(request_cpy);
        SessionUpdate update =
            SessionStore::get_default_session_update(session_map);
        auto result =
            enforcer_->init_charging_reauth(session_map, request_cpy, update);
        MLOG(MDEBUG) << "Result of Gy (Charging) ReAuthRequest "
                     << raa_result_to_str(result);
        ChargingReAuthAnswer ans;
        ans.set_result(result);

        bool update_success = session_store_.update_sessions(update);
        if (update_success) {
          response_callback(Status::OK, ans);
        } else {
          // Todo If update fails, we should rollback changes from the request
          MLOG(MERROR) << "Failed to update Gy (Charging) ReAuthRequest changes...";
          auto status = Status(
              grpc::ABORTED,
              "ChargingReAuth no longer valid due to another update that "
              "updated the session first.");
          response_callback(status, ans);
        }
      });
}

void SessionProxyResponderHandlerImpl::PolicyReAuth(
    ServerContext* context, const PolicyReAuthRequest* request,
    std::function<void(Status, PolicyReAuthAnswer)> response_callback) {
  auto& request_cpy = *request;
  MLOG(MDEBUG) << "Received a Gx (Policy) ReAuthRequest for session_id "
               << request->session_id();
  enforcer_->get_event_base().runInEventBaseThread(
      [this, request_cpy, response_callback]() {
        PolicyReAuthAnswer ans;
        auto session_map = get_sessions_for_policy(request_cpy);
        SessionUpdate update =
            SessionStore::get_default_session_update(session_map);
        enforcer_->init_policy_reauth(session_map, request_cpy, ans, update);
        MLOG(MDEBUG) << "Result of Gx (Policy) ReAuthRequest " << ans.result();
        bool update_success = session_store_.update_sessions(update);
        if (update_success) {
          response_callback(Status::OK, ans);
        } else {
        // Todo If update fails, we should rollback changes from the request
          MLOG(MERROR) << "Failed to update Gx (Policy) ReAuthRequest changes...";
          auto status = Status(
              grpc::ABORTED,
              "PolicyReAuth no longer valid due to another update that "
              "updated the session first.");
          response_callback(status, ans);
        }
      });
}

SessionMap SessionProxyResponderHandlerImpl::get_sessions_for_charging(
    const ChargingReAuthRequest& request) {
  SessionRead req = {request.sid()};
  return session_store_.read_sessions(req);
}

SessionMap SessionProxyResponderHandlerImpl::get_sessions_for_policy(
    const PolicyReAuthRequest& request) {
  SessionRead req = {request.imsi()};
  return session_store_.read_sessions(req);
}

std::string raa_result_to_str(ReAuthResult res) {
  switch (res) {
    case UPDATE_INITIATED:
      return "UPDATE_INITIATED";
    case UPDATE_NOT_NEEDED:
      return "UPDATE_NOT_NEEDED";
    case SESSION_NOT_FOUND:
      return "SESSION_NOT_FOUND";
    case OTHER_FAILURE:
      return "OTHER_FAILURE";
    default:
      return "UNKNOWN_RESULT";
  }
}
}  // namespace magma

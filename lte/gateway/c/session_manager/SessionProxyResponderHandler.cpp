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

SessionProxyResponderHandlerImpl::SessionProxyResponderHandlerImpl(
    std::shared_ptr<LocalEnforcer> enforcer, SessionStore& session_store)
    : enforcer_(enforcer), session_store_(session_store) {}

void SessionProxyResponderHandlerImpl::ChargingReAuth(
    ServerContext* context, const ChargingReAuthRequest* request,
    std::function<void(Status, ChargingReAuthAnswer)> response_callback) {
  auto& request_cpy = *request;
  enforcer_->get_event_base().runInEventBaseThread(
      [this, request_cpy, response_callback]() {
        auto session_map = get_sessions_for_charging(request_cpy);
        SessionUpdate update =
            SessionStore::get_default_session_update(session_map);
        auto result =
            enforcer_->init_charging_reauth(session_map, request_cpy, update);
        ChargingReAuthAnswer ans;
        ans.set_result(result);

        bool update_success = session_store_.update_sessions(update);
        if (update_success) {
          response_callback(Status::OK, ans);
        } else {
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
  enforcer_->get_event_base().runInEventBaseThread(
      [this, request_cpy, response_callback]() {
        PolicyReAuthAnswer ans;
        auto session_map = get_sessions_for_policy(request_cpy);
        SessionUpdate update =
            SessionStore::get_default_session_update(session_map);
        enforcer_->init_policy_reauth(session_map, request_cpy, ans, update);

        bool update_success = session_store_.update_sessions(update);
        if (update_success) {
          response_callback(Status::OK, ans);
        } else {
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
}  // namespace magma

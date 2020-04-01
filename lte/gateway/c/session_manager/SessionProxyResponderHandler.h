/**
 * Copyright (c) 2016-present, Facebook, Inc.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */
#pragma once

#include <functional>

#include <grpc++/grpc++.h>
#include <lte/protos/session_manager.grpc.pb.h>

#include "LocalEnforcer.h"
#include "SessionStore.h"

using grpc::Server;
using grpc::ServerContext;
using grpc::Status;

namespace magma {

class SessionProxyResponderHandler {
 public:
  virtual ~SessionProxyResponderHandler() {}

  /**
   * Reengage a subscriber service, usually after new credit is added to the
   * account
   */
  virtual void ChargingReAuth(
    ServerContext *context,
    const ChargingReAuthRequest *request,
    std::function<void(Status, ChargingReAuthAnswer)> response_callback) = 0;

  virtual void PolicyReAuth(
    ServerContext *context,
    const PolicyReAuthRequest *request,
    std::function<void(Status, PolicyReAuthAnswer)> response_callback) = 0;
};

/**
 * SessionProxyResponderHandlerImpl responds to requests coming from the
 * federated gateway, such as Re-Auth
 */
class SessionProxyResponderHandlerImpl : public SessionProxyResponderHandler {
 public:
  SessionProxyResponderHandlerImpl(
    std::shared_ptr<LocalEnforcer> monitor,
    SessionMap& session_map,
    SessionStore& session_store);

  ~SessionProxyResponderHandlerImpl() {}

  /**
   * Reengage a subscriber service, usually after new credit is added to the
   * account
   */
  void ChargingReAuth(
    ServerContext *context,
    const ChargingReAuthRequest *request,
    std::function<void(Status, ChargingReAuthAnswer)> response_callback);

  /**
   * Install/uninstall rules for an existing session
   */
  void PolicyReAuth(
    ServerContext *context,
    const PolicyReAuthRequest *request,
    std::function<void(Status, PolicyReAuthAnswer)> response_callback);

 private:
   SessionMap& session_map_;
   SessionStore& session_store_;
   std::shared_ptr<LocalEnforcer> enforcer_;

 private:
  /**
   * Get the most recently written state of the session to be updated for
   * charging reauth.
   * Does not get any other sessions.
   *
   * NOTE: Call only from the main EventBase thread, otherwise there will
   *       be undefined behavior.
   */
    SessionMap get_sessions_for_charging(const ChargingReAuthRequest& request);

  /**
   * Get the most recently written state of the session to be updated for
   * policy reauth.
   * Does not get any other sessions.
   *
   * NOTE: Call only from the main EventBase thread, otherwise there will
   *       be undefined behavior.
   */
    SessionMap get_sessions_for_policy(const PolicyReAuthRequest& request);
};

} // namespace magma

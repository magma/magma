/**
 * Copyright (c) 2016-present, Facebook, Inc.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */
#pragma once

#include "GRPCReceiver.h"
#include "SessionState.h"
#include "SessionStore.h"
#include <feg/gateway/services/aaa/protos/accounting.grpc.pb.h>

using grpc::Status;

namespace aaa {
using namespace protos;

/**
 * AAAClient is the base class for interacting with AAA service
 */
class AAAClient {
 public:
  virtual bool terminate_session(
    const std::string &radius_session_id,
    const std::string &imsi) = 0;

  virtual bool add_sessions(magma::lte::SessionMap &session_map) = 0;
};

/**
 * AsyncAAAClient implements AAAClient and sends call
 * asynchronously to AAA service.
 */
class AsyncAAAClient : public magma::GRPCReceiver, public AAAClient {
 public:
  AsyncAAAClient();

  AsyncAAAClient(std::shared_ptr<grpc::Channel> aaa_channel);

  bool terminate_session(
    const std::string &radius_session_id,
    const std::string &imsi);

  bool add_sessions(magma::lte::SessionMap &session_map);

 private:
  static const uint32_t RESPONSE_TIMEOUT = 6; // seconds
  std::unique_ptr<accounting::Stub> stub_;

 private:
  void terminate_session_rpc(
    const terminate_session_request &request,
    std::function<void(Status, acct_resp)> callback);
  void add_sessions_rpc(
    const add_sessions_request &request,
    std::function<void(Status, acct_resp)> callback);
};

} // namespace aaa

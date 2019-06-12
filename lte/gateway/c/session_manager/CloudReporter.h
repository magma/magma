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
#include <memory>

#include <folly/io/async/EventBase.h>
#include <grpc++/grpc++.h>
#include <lte/protos/session_manager.grpc.pb.h>

#include "GRPCReceiver.h"

namespace magma {
using namespace lte;
/**
 * AsyncEvbResponse is used to call a callback in a particular event loop when
 * a response is received. This is defined here to limit the dependency on folly
 * in the common library.
 */
template<typename ResponseType>
class AsyncEvbResponse : public AsyncGRPCResponse<ResponseType> {
 public:
  AsyncEvbResponse(
    folly::EventBase *base,
    std::function<void(grpc::Status, ResponseType)> callback,
    uint32_t timeout_sec);

  void handle_response() override;

 private:
  folly::EventBase *base_;
};

class SessionCloudReporter{
 public:
  /**
   * Proxy an UpdateSessionRequest gRPC call to the cloud
   */
  virtual void report_updates(
    const UpdateSessionRequest &request,
    std::function<void(grpc::Status, UpdateSessionResponse)> callback) = 0;

  /**
   * Proxy a CreateSessionRequest gRPC call to the cloud
   */
  virtual void report_create_session(
    const CreateSessionRequest &request,
    std::function<void(grpc::Status, CreateSessionResponse)> callback) = 0;

  /**
   * Proxy a SessionTerminateRequest gRPC call to the cloud
   */
  virtual void report_terminate_session(
    const SessionTerminateRequest &request,
    std::function<void(grpc::Status, SessionTerminateResponse)> callback) = 0;
};

class SessionCloudReporterImpl : public GRPCReceiver, public SessionCloudReporter {
 public:
  SessionCloudReporterImpl(
    folly::EventBase *base,
    std::shared_ptr<grpc::Channel> channel);

  void report_updates(
    const UpdateSessionRequest &request,
    std::function<void(grpc::Status, UpdateSessionResponse)> callback);

  void report_create_session(
    const CreateSessionRequest &request,
    std::function<void(grpc::Status, CreateSessionResponse)> callback);

  void report_terminate_session(
    const SessionTerminateRequest &request,
    std::function<void(grpc::Status, SessionTerminateResponse)> callback);

 private:
  folly::EventBase *base_;
  std::unique_ptr<CentralSessionController::Stub> stub_;
  static const uint32_t RESPONSE_TIMEOUT = 6; // seconds
};

} // namespace magma

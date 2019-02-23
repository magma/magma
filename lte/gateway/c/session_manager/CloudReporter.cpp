/**
 * Copyright (c) 2016-present, Facebook, Inc.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */
#include <iostream>
#include <glog/logging.h>
#include "CloudReporter.h"

namespace magma {

template<class ResponseType>
AsyncEvbResponse<ResponseType>::AsyncEvbResponse(
  folly::EventBase* base,
  std::function<void(grpc::Status, ResponseType)> callback,
  uint32_t timeout_sec)
  : base_(base),
    AsyncGRPCResponse<ResponseType>(callback, timeout_sec) {}

template<class ResponseType>
void AsyncEvbResponse<ResponseType>::handle_response() {
  base_->runInEventBaseThread([this]() {
    this->callback_(this->status_, this->response_);
    delete this;
  });
}

SessionCloudReporter::SessionCloudReporter(
  folly::EventBase* base,
  std::shared_ptr<grpc::Channel> channel)
  : base_(base),
    stub_(CentralSessionController::NewStub(channel)) {}

void SessionCloudReporter::report_updates(
    const UpdateSessionRequest& request,
    std::function<void(grpc::Status, UpdateSessionResponse)> callback) {
  auto cloud_response = new AsyncEvbResponse<UpdateSessionResponse>(
    base_, callback, RESPONSE_TIMEOUT);
  cloud_response->set_response_reader(std::move(stub_->AsyncUpdateSession(
    cloud_response->get_context(), request, &queue_)));
}

void SessionCloudReporter::report_create_session(
    const CreateSessionRequest& request,
    std::function<void(grpc::Status, CreateSessionResponse)> callback) {
  auto cloud_response = new AsyncEvbResponse<CreateSessionResponse>(
    base_, callback, RESPONSE_TIMEOUT);
  cloud_response->set_response_reader(std::move(stub_->AsyncCreateSession(
    cloud_response->get_context(), request, &queue_)));
}

void SessionCloudReporter::report_terminate_session(
    const SessionTerminateRequest& request,
    std::function<void(grpc::Status, SessionTerminateResponse)> callback) {
  auto cloud_response = new AsyncEvbResponse<SessionTerminateResponse>(
    base_, callback, RESPONSE_TIMEOUT);
  cloud_response->set_response_reader(std::move(stub_->AsyncTerminateSession(
    cloud_response->get_context(), request, &queue_)));
}


}

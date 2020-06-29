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
#include "ServiceRegistrySingleton.h"
#include "SessionReporter.h"
#include "magma_logging.h"

namespace magma {

template<class ResponseType>
AsyncEvbResponse<ResponseType>::AsyncEvbResponse(
  folly::EventBase* base,
  std::function<void(grpc::Status, ResponseType)> callback,
  uint32_t timeout_sec):
  AsyncGRPCResponse<ResponseType>(callback, timeout_sec),
  base_(base)
{
}

template<class ResponseType>
void AsyncEvbResponse<ResponseType>::handle_response()
{
  base_->runInEventBaseThread([this]() {
    this->callback_(this->status_, this->response_);
    delete this;
  });
}

ReporterCallbackFn<SessionTerminateResponse>
  SessionReporter::get_terminate_logging_cb(
    const SessionTerminateRequest& request)
{
    return [request](grpc::Status status, SessionTerminateResponse response) {
      if (!status.ok()) {
        MLOG(MERROR) << "Failed to terminate session in controller for "
                        "subscriber "
                     << request.sid() << ": " << status.error_message();
      } else {
        MLOG(MDEBUG) << "Termination successful in controller for "
                        "subscriber "
                     << request.sid();
      }
    };
}

SessionReporterImpl::SessionReporterImpl(
  folly::EventBase* base,
  std::shared_ptr<grpc::Channel> channel):
  base_(base),
  stub_(CentralSessionController::NewStub(channel))
{
}

void SessionReporterImpl::report_updates(
  const UpdateSessionRequest& request,
  ReporterCallbackFn<UpdateSessionResponse> callback)
{
  auto controller_response = new AsyncEvbResponse<UpdateSessionResponse>(
    base_, callback, RESPONSE_TIMEOUT);
  controller_response->set_response_reader(std::move(stub_->AsyncUpdateSession(
    controller_response->get_context(), request, &queue_)));
}

void SessionReporterImpl::report_create_session(
  const CreateSessionRequest& request,
  ReporterCallbackFn<CreateSessionResponse> callback)
{
  auto controller_response = new AsyncEvbResponse<CreateSessionResponse>(
    base_, callback, RESPONSE_TIMEOUT);
  controller_response->set_response_reader(std::move(stub_->AsyncCreateSession(
    controller_response->get_context(), request, &queue_)));
}

void SessionReporterImpl::report_terminate_session(
  const SessionTerminateRequest& request,
  ReporterCallbackFn<SessionTerminateResponse> callback)
{
  auto controller_response = new AsyncEvbResponse<SessionTerminateResponse>(
    base_, callback, RESPONSE_TIMEOUT);
  controller_response->set_response_reader(
    std::move(stub_->AsyncTerminateSession(controller_response->get_context(),
      request, &queue_)));
}

} // namespace magma

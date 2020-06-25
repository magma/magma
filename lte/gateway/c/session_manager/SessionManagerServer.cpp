/**
 * Copyright (c) 2016-present, Facebook, Inc.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */
#include "SessionManagerServer.h"
#include "magma_logging.h"

using grpc::ServerContext;
using grpc::Status;

namespace magma {

AsyncService::AsyncService(std::unique_ptr<ServerCompletionQueue> cq):
  cq_(std::move(cq))
{
}

void AsyncService::wait_for_requests()
{
  init_call_data();
  void *tag;
  bool ok;
  running_ = true;
  while (running_) {
    if (cq_ == nullptr || !cq_->Next(&tag, &ok)) {
      // queue shut down
      MLOG(MINFO) << "sessiond request completion queue shutting down";
      return;
    }
    if (!ok) {
      MLOG(MINFO)
        << "sessiond server encountered error while processing request";
      continue;
    }
    static_cast<CallData *>(tag)->proceed();
  }
}

void AsyncService::stop()
{
  running_ = false;
  cq_->Shutdown();
}

LocalSessionManagerAsyncService::LocalSessionManagerAsyncService(
  std::unique_ptr<ServerCompletionQueue> cq,
  std::unique_ptr<LocalSessionManagerHandler> handler):
  AsyncService(std::move(cq)),
  handler_(std::move(handler))
{
}

void LocalSessionManagerAsyncService::init_call_data()
{
  new ReportRuleStatsCallData(cq_.get(), *this, *handler_);
  new CreateSessionCallData(cq_.get(), *this, *handler_);
  new EndSessionCallData(cq_.get(), *this, *handler_);
}

SessionProxyResponderAsyncService::SessionProxyResponderAsyncService(
  std::unique_ptr<ServerCompletionQueue> cq,
  std::unique_ptr<SessionProxyResponderHandler> handler):
  AsyncService(std::move(cq)),
  handler_(std::move(handler))
{
}

void SessionProxyResponderAsyncService::init_call_data()
{
  new ChargingReAuthCallData(cq_.get(), *this, *handler_);
  new PolicyReAuthCallData(cq_.get(), *this, *handler_);
}

template<class GRPCService, class RequestType, class ResponseType>
AsyncGRPCRequest<GRPCService, RequestType, ResponseType>::AsyncGRPCRequest(
  ServerCompletionQueue *cq,
  GRPCService &service):
  cq_(cq),
  status_(PROCESS),
  responder_(&ctx_),
  service_(service)
{
}

template<class GRPCService, class RequestType, class ResponseType>
void AsyncGRPCRequest<GRPCService, RequestType, ResponseType>::proceed()
{
  if (status_ == PROCESS) {
    clone();
    process();
    status_ = FINISH;
  } else {
    GPR_ASSERT(status_ == FINISH);
    delete this;
  }
}

template<class GRPCService, class RequestType, class ResponseType>
std::function<void(Status, ResponseType)>
AsyncGRPCRequest<GRPCService, RequestType, ResponseType>::get_finish_callback()
{
  return [this](Status status, ResponseType response) {
    responder_.Finish(response, status, (void *) this);
  };
}

} // namespace magma

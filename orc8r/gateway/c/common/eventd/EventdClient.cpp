/**
 * Copyright (c) 2016-present, Facebook, Inc.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */
#include <memory>
#include <utility>

#include <orc8r/protos/eventd.pb.h>
#include <orc8r/protos/eventd.grpc.pb.h>

#include "EventdClient.h"
#include "GRPCReceiver.h"
#include "ServiceRegistrySingleton.h"

using grpc::Status;
using grpc::ClientContext;

namespace magma {

AsyncEventdClient::AsyncEventdClient(std::shared_ptr<grpc::Channel> channel):
    stub_(EventService::NewStub(channel)) {
}

AsyncEventdClient::AsyncEventdClient():
  AsyncEventdClient(ServiceRegistrySingleton::Instance()->GetGrpcChannel(
    "eventd",
    ServiceRegistrySingleton::LOCAL)) {
}

void AsyncEventdClient::log_event(
    const Event& request,
    std::function<void(Status status, Void)> callback) {
  auto local_response =
    new AsyncLocalResponse<Void>(std::move(callback), RESPONSE_TIMEOUT);
  local_response->set_response_reader(std::move(
        stub_->AsyncLogEvent(local_response->get_context(), request, &queue_)));
}

}  // namespace magma

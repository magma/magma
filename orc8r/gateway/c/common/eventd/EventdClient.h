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

#include "GRPCReceiver.h"

using grpc::Status;

namespace magma {

/**
 * AsyncEventdClient sends asynchronous calls to eventd
 * to log events
 */
class AsyncEventdClient : public GRPCReceiver {
 public:
  AsyncEventdClient(AsyncEventdClient const &) = delete;
  void operator=(AsyncEventdClient const &) = delete;

  static AsyncEventdClient &getInstance();

  // Logs an event
  void log_event(
      const orc8r::Event& request,
      std::function<void(Status status, orc8r::Void)> callback);

 private:
  AsyncEventdClient();
  static const uint32_t RESPONSE_TIMEOUT = 6;  // seconds
  std::unique_ptr<orc8r::EventService::Stub> stub_{};
};

}  // namespace magma

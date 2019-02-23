/**
 * Copyright (c) 2016-present, Facebook, Inc.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */
#include "GRPCReceiver.h"
#include "magma_logging.h"

namespace magma {

void GRPCReceiver::rpc_response_loop() {
  running_ = true;
  void* tag;
  bool ok = false;
  while (running_) {
    if (!queue_.Next(&tag, &ok)) {
      return;
    }
    if (!ok) {
      MLOG(MINFO) << "gRPC receiver encountered error while processing request";
      continue;
    }
    static_cast<AsyncResponse*>(tag)->handle_response();
  }
}

void GRPCReceiver::stop() {
  running_ = false;
  queue_.Shutdown();
}

}

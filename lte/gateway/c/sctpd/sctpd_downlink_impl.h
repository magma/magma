/**
 * Copyright (c) 2016-present, Facebook, Inc.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#pragma once

#include <memory>
#include <thread>

#include <grpc/grpc.h>
#include <grpcpp/server_context.h>

#include <lte/protos/sctpd.grpc.pb.h>

#include "sctp_connection.h"

namespace magma {
namespace sctpd {

using grpc::ServerContext;
using grpc::Status;

// Implements the sctpd downlink server
class SctpdDownlinkImpl final : public SctpdDownlink::Service {
 public:
  // Construct a new SctpdDownlinkImpl service
  SctpdDownlinkImpl(SctpEventHandler &uplink_handler);

  // Implementation of SctpdDownlink.Init method (see sctpd.proto for more info)
  Status Init(ServerContext *context, const InitReq *request, InitRes *response)
    override;

  // Implementation of SctpdDownlink.SendDl method (see sctpd.proto for more info)
  Status SendDl(
    ServerContext *context,
    const SendDlReq *request,
    SendDlRes *response) override;

  // Close SCTP connection for this SctpdDownlink.
  void stop();

 private:
  SctpEventHandler &_uplink_handler;
  std::unique_ptr<SctpConnection> _sctp_connection;
};

} // namespace sctpd
} // namespace magma

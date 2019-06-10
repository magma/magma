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

#include <grpcpp/grpcpp.h>

#include <lte/protos/sctpd.grpc.pb.h>

namespace magma {
namespace sctpd {

using grpc::Channel;

// Grpc uplink client to allow sctpd to signal MME
class SctpdUplinkClient {
 public:
  // Construct SctpdUplinkClient with the specified channel
  explicit SctpdUplinkClient(std::shared_ptr<Channel> channel);

  // Send an uplink packet to MME (see sctpd.proto for more info)
  int sendUl(const SendUlReq &req, SendUlRes *res);
  // Notify MME of new association (see sctpd.proto for more info)
  int newAssoc(const NewAssocReq &req, NewAssocRes *res);
  // Notify MME of closing/reseting association (see sctpd.proto for more info)
  int closeAssoc(const CloseAssocReq &req, CloseAssocRes *res);

 private:
  // Stub used for client to communicate with server
  std::unique_ptr<SctpdUplink::Stub> _stub;
};

} // namespace sctpd
} // namespace magma

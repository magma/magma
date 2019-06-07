/**
 * Copyright (c) 2016-present, Facebook, Inc.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#pragma once

#include "sctp_connection.h"

#include "sctpd_uplink_client.h"

namespace magma {
namespace sctpd {

// Sctp handler that relays events to MME over GRPC
class SctpdEventHandler : public SctpEventHandler {
 public:
  // Construct SctpdEventHandler that communicates to MME over client
  explicit SctpdEventHandler(SctpdUplinkClient &client);

  // Relay new assocation to MME over GRPC
  void HandleNewAssoc(
    uint32_t assoc_id,
    uint32_t instreams,
    uint32_t outstreams) override;

  // Relay close assocation to MME over GRPC
  void HandleCloseAssoc(uint32_t assoc_id, bool reset) override;

  // Relay new message to MME over GRPC
  void HandleRecv(
    uint32_t assoc_id,
    uint32_t stream,
    const std::string &payload) override;

 private:
  SctpdUplinkClient &_client;
};

} // namespace sctpd
} // namespace magma

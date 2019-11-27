/**
 * Copyright (c) 2016-present, Facebook, Inc.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */
#pragma once

#include <mutex>
#include <unordered_map>

#include <orc8r/protos/directoryd.pb.h>
#include <orc8r/protos/directoryd.grpc.pb.h>

#include "GRPCReceiver.h"
#include "SessionState.h"

namespace magma {
using namespace orc8r;

/**
 * AsyncDirectorydClient sends asynchronous calls to directoryd to retrieve
 * UE information.
 */
class AsyncDirectorydClient : public GRPCReceiver {
 public:
  AsyncDirectorydClient();

  AsyncDirectorydClient(std::shared_ptr<grpc::Channel> directoryd_channel);

  /**
   * Gets the directoryd imsi's 'ip' field
   * @param imsi - UE to query
   * @return true if the operation was successful
   */
  bool get_directoryd_ip_field(
    const std::string& imsi,
    std::function<void(Status status, DirectoryField)> callback);

 private:
  static const uint32_t RESPONSE_TIMEOUT = 6; // seconds
  std::unique_ptr<GatewayDirectoryService::Stub> stub_;

 private:
  void get_directoryd_ip_field_rpc(
    const GetDirectoryFieldRequest &request,
    std::function<void(Status, DirectoryField)> callback);
};

} // namespace magma

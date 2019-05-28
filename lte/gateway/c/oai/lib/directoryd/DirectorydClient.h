/*
 * Licensed to the OpenAirInterface (OAI) Software Alliance under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The OpenAirInterface Software Alliance licenses this file to You under
 * the Apache License, Version 2.0  (the "License"); you may not use this file
 * except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *-------------------------------------------------------------------------------
 * For more information about the OpenAirInterface (OAI) Software Alliance:
 *      contact@openairinterface.org
 */

#pragma once

#include <grpc++/grpc++.h>
#include <stdint.h>
#include <functional>
#include <memory>
#include <string>

#include "orc8r/protos/directoryd.grpc.pb.h"
#include "GRPCReceiver.h"
#include "orc8r/protos/directoryd.pb.h"

namespace grpc {
class Channel;
class ClientContext;
class Status;
} // namespace grpc
namespace magma {
namespace orc8r {
class Void;
} // namespace orc8r
} // namespace magma

using grpc::Channel;
using grpc::ClientContext;
using grpc::Status;
using magma::orc8r::DirectoryService;

namespace magma {
using namespace orc8r;
/*
 * gRPC client for DirectoryService
 */
class DirectoryServiceClient : public GRPCReceiver {
 public:
  static bool UpdateLocation(
    TableID table,
    const std::string &id,
    const std::string &location,
    std::function<void(Status, Void)> callback);

  static bool DeleteLocation(
    TableID table,
    const std::string &id,
    std::function<void(Status, Void)> callback);

 public:
  DirectoryServiceClient(DirectoryServiceClient const &) = delete;
  void operator=(DirectoryServiceClient const &) = delete;

 private:
  DirectoryServiceClient();
  static DirectoryServiceClient &get_instance();
  std::shared_ptr<DirectoryService::Stub> stub_;
  static const uint32_t RESPONSE_TIMEOUT = 30; // seconds
};

} // namespace magma

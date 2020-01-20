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

#include "lte/protos/session_manager.grpc.pb.h"
#include "GRPCReceiver.h"

namespace grpc {
class Status;
}  // namespace grpc
namespace magma {
namespace lte {
class LocalCreateSessionRequest;
class LocalCreateSessionResponse;
class LocalEndSessionResponse;
class SubscriberID;
}  // namespace lte
}  // namespace magma

using grpc::Status;

namespace magma {
using namespace lte;

/**
 * PCEFClient is the main asynchronous client for interacting with sessiond.
 * Responses will come in a queue and call the callback passed
 * To start the client and make sure it receives calls, one must call the
 * rpc_response_loop method defined in the GRPCReceiver base class
 */
class PCEFClient : public GRPCReceiver {
 public:
  /**
   * Proxy a CreateSession gRPC call to sessiond
   */
  static void create_session(
    const LocalCreateSessionRequest &request,
    std::function<void(Status, LocalCreateSessionResponse)> callback);

  /**
   * Proxy a EndSession gRPC call to sessiond
   */
  static void end_session(
    const LocalEndSessionRequest &request,
    std::function<void(Status, LocalEndSessionResponse)> callback);

 public:
  PCEFClient(PCEFClient const &) = delete;
  void operator=(PCEFClient const &) = delete;

 private:
  PCEFClient();
  static PCEFClient &get_instance();
  std::unique_ptr<LocalSessionManager::Stub> stub_;
  static const uint32_t RESPONSE_TIMEOUT = 10; // seconds
};

} // namespace magma

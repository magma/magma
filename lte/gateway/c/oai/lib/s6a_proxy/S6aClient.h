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

#include <functional>
#include <memory>

#include <gmp.h>
#include <grpc++/grpc++.h>
#include "lte/protos/s6a_proxy.grpc.pb.h"

#include "GRPCReceiver.h"

extern "C" {
#include "intertask_interface.h"
}

using grpc::Status;

namespace magma {

/**
 * S6aClient is the main asynchronous client for interacting with s6a_proxy.
 * Responses will come in a queue and call the callback passed
 * To start the client and make sure it receives calls, one must call the
 * rpc_response_loop method defined in the GRPCReceiver base class
 */
class S6aClient : public GRPCReceiver {
 public:
  /**
   * Proxy a purge gRPC call to s6a_proxy
   */
  static void purge_ue(
    const char *imsi,
    std::function<void(Status, lte::PurgeUEAnswer)> callback);

  /**
   * Proxy a purge gRPC call to s6a_proxy
   */
  static void authentication_info_req(
    const s6a_auth_info_req_t *const msg,
    std::function<void(Status, lte::AuthenticationInformationAnswer)> callbk);

  /**
   * Proxy a purge gRPC call to s6a_proxy
   */
  static void update_location_request(
    const s6a_update_location_req_t *const msg,
    std::function<void(Status, lte::UpdateLocationAnswer)> callback);

 public:
  S6aClient(S6aClient const &) = delete;
  void operator=(S6aClient const &) = delete;

 private:
  S6aClient();
  static S6aClient &get_instance();
  std::unique_ptr<lte::S6aProxy::Stub> stub_;
  static const uint32_t RESPONSE_TIMEOUT = 3; // seconds
};

bool get_s6a_relay_enabled(void);

} // namespace magma

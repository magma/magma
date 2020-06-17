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

#include "EventClientAPI.h"

#include <iostream>
#include <thread>
#include <grpcpp/support/status.h>
#include <orc8r/protos/common.pb.h>

#include "EventdClient.h"

using grpc::Status;
using magma::AsyncEventdClient;
using magma::orc8r::Event;
using magma::orc8r::Void;

namespace magma {
namespace lte {

void init_eventd_client() {
  auto& client = AsyncEventdClient::getInstance();
  std::thread resp_loop_thread([&]() { client.rpc_response_loop(); });
  resp_loop_thread.detach();
}

int log_event(const Event& event) {
  AsyncEventdClient::getInstance().log_event(event, [=](Status status, Void v) {
    if (status.ok()) {
      std::cout << "[DEBUG] Success logging event: " << event.event_type()
                << std::endl;
    } else {
      std::cout << "[ERROR] Failed to log event: " << event.event_type()
                << "; Status: " << status.error_message() << std::endl;
    }
  });
  return 0;
}

}  // namespace lte
}  // namespace magma

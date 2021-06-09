/*
 * Licensed to the OpenAirInterface (OAI) Software Alliance under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The OpenAirInterface Software Alliance licenses this file to You under
 * the terms found in the LICENSE file in the root of this source tree.
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

#include "includes/EventdClient.h"

using grpc::Status;
using grpc::StatusCode::DEADLINE_EXCEEDED;
using grpc::StatusCode::UNAVAILABLE;
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
      return 0;
    }
    if (status.error_code() == DEADLINE_EXCEEDED ||
        status.error_code() == UNAVAILABLE) {
      return 0;  // Suppress error logs if EventD is unavailable
    }
    std::cout << "[ERROR] Failed to log event: " << event.event_type()
              << "; Status: " << status.error_message() << std::endl;
    return int(status.error_code());
  });
  return 0;
}

}  // namespace lte
}  // namespace magma

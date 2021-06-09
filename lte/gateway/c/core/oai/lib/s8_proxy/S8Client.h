/*
Copyright 2020 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

#pragma once

#include <gmp.h>
#include <grpc++/grpc++.h>
#include <stdint.h>
#include <string>
#include <functional>
#include <memory>

#include "includes/GRPCReceiver.h"
#include "feg/protos/s8_proxy.grpc.pb.h"

extern "C" {
#include "intertask_interface.h"

namespace grpc {
class Status;
}
}

namespace magma {
using namespace feg;

/**
 * S8Client is the main asynchronous client for interacting with FedGW.
 * Responses will come in a queue and call the callback passed.
 * To start the client and make sure it receives calls, one must call the
 * rpc_response_loop method defined in the GRPCReceiver base class
 */
class S8Client : public GRPCReceiver {
 public:
  // Send Create Session Request
  static void s8_create_session_request(
      const CreateSessionRequestPgw& csr_req,
      std::function<void(grpc::Status, CreateSessionResponsePgw)> callback);

  // Send Delete Session Request
  static void s8_delete_session_request(
      const DeleteSessionRequestPgw& dsr_req,
      std::function<void(grpc::Status, DeleteSessionResponsePgw)> callback);

 public:
  S8Client(S8Client const&) = delete;
  void operator=(S8Client const&) = delete;

 private:
  S8Client();
  static S8Client& get_instance();
  std::unique_ptr<S8Proxy::Stub> stub_;
  static const uint32_t RESPONSE_TIMEOUT = 10;  // seconds
};

}  // namespace magma

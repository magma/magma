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

#include "GRPCReceiver.h"
#include "feg/protos/s8_proxy.grpc.pb.h"

extern "C" {
#include "intertask_interface.h"

namespace grpc {
class Status;
}  // namespace grpc
namespace magma {
namespace orc8r {
class Void;
}  // namespace orc8r
}  // namespace magma
}

namespace magma {
using namespace orc8r;
using namespace feg;

/**
 * S8 is the main client for sending S8 messages to FeG
 * FeG will forward the message to Roaming network's PGW then respond instantly
 * with Void
 */
class S8Client : public GRPCReceiver {
 public:
  // Send Create Session Request
  static void s8_create_session_request(
      const CreateSessionRequestPgw& csr_req,
      std::function<void(grpc::Status, CreateSessionResponsePgw)> callback);

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

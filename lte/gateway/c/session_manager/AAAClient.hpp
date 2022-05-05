/**
 * Copyright 2020 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
#pragma once

// IWYU pragma: no_include <context.pb.h>
#include <feg/gateway/services/aaa/protos/accounting.grpc.pb.h>
#include <stdint.h>
#include <functional>
#include <memory>
#include <string>

#include "lte/gateway/c/session_manager/SessionState.hpp"
#include "lte/gateway/c/session_manager/SessionStore.hpp"
#include "lte/gateway/c/session_manager/StoreClient.hpp"
#include "orc8r/gateway/c/common/async_grpc/GRPCReceiver.hpp"

namespace aaa {
namespace protos {
class acct_resp;
class add_sessions_request;
class terminate_session_request;
}  // namespace protos
}  // namespace aaa
namespace grpc {
class Channel;
class Status;
}  // namespace grpc

using grpc::Status;

namespace aaa {
using namespace protos;

/**
 * AAAClient is the base class for interacting with AAA service
 */
class AAAClient {
 public:
  virtual ~AAAClient() = default;

  virtual bool terminate_session(const std::string& radius_session_id,
                                 const std::string& imsi) = 0;

  virtual bool add_sessions(const magma::lte::SessionMap& session_map) = 0;
};

/**
 * AsyncAAAClient implements AAAClient and sends call
 * asynchronously to AAA service.
 */
class AsyncAAAClient : public magma::GRPCReceiver, public AAAClient {
 public:
  AsyncAAAClient();

  explicit AsyncAAAClient(std::shared_ptr<grpc::Channel> aaa_channel);

  bool terminate_session(const std::string& radius_session_id,
                         const std::string& imsi);

  bool add_sessions(const magma::lte::SessionMap& session_map);

 private:
  static const uint32_t RESPONSE_TIMEOUT = 6;  // seconds
  std::unique_ptr<accounting::Stub> stub_;

 private:
  void terminate_session_rpc(const terminate_session_request& request,
                             std::function<void(Status, acct_resp)> callback);
  void add_sessions_rpc(const add_sessions_request& request,
                        std::function<void(Status, acct_resp)> callback);
};

}  // namespace aaa

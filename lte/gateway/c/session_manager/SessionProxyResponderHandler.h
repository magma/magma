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

#include <grpc++/grpc++.h>
#include <lte/protos/abort_session.grpc.pb.h>
#include <lte/protos/session_manager.grpc.pb.h>
#include <functional>
#include <memory>

#include "LocalEnforcer.h"
#include "SessionStore.h"

namespace grpc {
class ServerContext;
class Status;
}  // namespace grpc
namespace magma {
class LocalEnforcer;
namespace lte {
class AbortSessionRequest;
class AbortSessionResult;
class ChargingReAuthAnswer;
class ChargingReAuthRequest;
class PolicyReAuthAnswer;
class PolicyReAuthRequest;
class SessionStore;
}  // namespace lte
}  // namespace magma

using grpc::ServerContext;
using grpc::Status;

namespace magma {

class SessionProxyResponderHandler {
 public:
  virtual ~SessionProxyResponderHandler() {}

  /**
   * Reengage a subscriber service, usually after new credit is added to the
   * account
   */
  virtual void ChargingReAuth(
      ServerContext* context, const ChargingReAuthRequest* request,
      std::function<void(Status, ChargingReAuthAnswer)> response_callback) = 0;

  virtual void PolicyReAuth(
      ServerContext* context, const PolicyReAuthRequest* request,
      std::function<void(Status, PolicyReAuthAnswer)> response_callback) = 0;

  /**
   * Support for network initiated ungraceful termination
   */
  virtual void AbortSession(
      ServerContext* context, const AbortSessionRequest* request,
      std::function<void(Status, AbortSessionResult)> response_callback) = 0;
};

/**
 * SessionProxyResponderHandlerImpl responds to requests coming from the
 * federated gateway, such as Re-Auth
 */
class SessionProxyResponderHandlerImpl : public SessionProxyResponderHandler {
 public:
  SessionProxyResponderHandlerImpl(std::shared_ptr<LocalEnforcer> monitor,
                                   SessionStore& session_store);

  ~SessionProxyResponderHandlerImpl() {}

  /**
   * Reengage a subscriber service, usually after new credit is added to the
   * account
   */
  void ChargingReAuth(
      ServerContext* context, const ChargingReAuthRequest* request,
      std::function<void(Status, ChargingReAuthAnswer)> response_callback);

  /**
   * Install/uninstall rules for an existing session
   */
  void PolicyReAuth(
      ServerContext* context, const PolicyReAuthRequest* request,
      std::function<void(Status, PolicyReAuthAnswer)> response_callback);

  /**
   * Support for network initiated ungraceful termination
   */
  void AbortSession(
      ServerContext* context, const AbortSessionRequest* request,
      std::function<void(Status, AbortSessionResult)> response_callback);

 private:
  std::shared_ptr<LocalEnforcer> enforcer_;
  SessionStore& session_store_;
};

}  // namespace magma

/*
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
#include <grpcpp/impl/codegen/status.h>

#include "lte/protos/session_manager.grpc.pb.h"
#include "lte/protos/policydb.pb.h"

namespace grpc {
class ServerContext;
}  // namespace grpc
namespace magma {
namespace lte {
class SetSMSessionContextAccess;
class SetSmNotificationContext;
class SmContextVoid;
}  // namespace lte
}  // namespace magma

using grpc::ServerContext;
using magma::lte::SetSmNotificationContext;
using magma::lte::SetSMSessionContextAccess;
using magma::lte::SmContextVoid;
using magma::lte::SmfPduSessionSmContext;

namespace magma {
using namespace lte;

// SessionD to AMF server
class AmfServiceImpl final : public SmfPduSessionSmContext::Service {
 public:
  AmfServiceImpl();

  grpc::Status SetAmfNotification(
      ServerContext* context, const SetSmNotificationContext* notif,
      SmContextVoid* response) override;

  grpc::Status SetSmfSessionContext(
      ServerContext* context, const SetSMSessionContextAccess* request,
      SmContextVoid* response) override;
};

}  // namespace magma

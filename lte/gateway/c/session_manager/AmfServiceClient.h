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

#include <mutex>

#include <grpc++/grpc++.h>
#include <lte/protos/session_manager.pb.h>
#include <lte/protos/session_manager.grpc.pb.h>

#include "includes/GRPCReceiver.h"

using grpc::Status;

namespace magma {
using namespace lte;

/**
 * AmfServiceClient is the base class for responding on establish
 * modification and release request message to AMF
 */
class AmfServiceClient {
 public:
  virtual ~AmfServiceClient() {}
  virtual bool handle_response_to_access(
      const magma::SetSMSessionContextAccess& response) = 0;
  virtual bool handle_notification_to_access(
      const magma::SetSmNotificationContext& req) = 0;
};  // end of abstract class

class AsyncAmfServiceClient : public GRPCReceiver, public AmfServiceClient {
 public:
  AsyncAmfServiceClient();
  AsyncAmfServiceClient(std::shared_ptr<grpc::Channel> amf_srv_channel);

  /* This will send response back to AMF for all three request messages
   * i.e. establish, modification and release messages
   */
  bool handle_response_to_access(
      const magma::SetSMSessionContextAccess& response);

  bool handle_notification_to_access(
      const magma::SetSmNotificationContext& req);

 private:
  static const uint32_t RESPONSE_TIMEOUT = 6;  // seconds
  std::unique_ptr<SmfPduSessionSmContext::Stub> stub_;
};

}  // namespace magma

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
#include <lte/protos/session_manager.grpc.pb.h>
#include <lte/protos/session_manager.pb.h>
#include <lte/protos/subscriberdb.pb.h>
#include "includes/GRPCReceiver.h"
#include <stdint.h>
#include <functional>
#include <memory>

using grpc::Status;

namespace magma5g {

class SmfServiceClient {
 public:
  virtual ~SmfServiceClient() {}
  virtual bool set_smf_session(
      const magma::lte::SetSMSessionContext& request) = 0;
};

/**
 * AsyncSmfServiceClient implements SmfServiceClient but sends calls
 * asynchronously to sessiond.
 */
class AsyncSmfServiceClient : public magma::GRPCReceiver,
                              public SmfServiceClient {
 private:
  static const uint32_t RESPONSE_TIMEOUT = 6;  // seconds
  std::unique_ptr<magma::lte::AmfPduSessionSmContext::Stub> stub_;
  void set_smf_session_rpc(
      const magma::lte::SetSMSessionContext& request,
      std::function<void(Status, magma::lte::SmContextVoid)> callback);

 public:
  AsyncSmfServiceClient();
  AsyncSmfServiceClient(std::shared_ptr<grpc::Channel> smf_srv_channel);
  bool set_smf_session(const magma::lte::SetSMSessionContext& request);
};

}  // namespace magma5g

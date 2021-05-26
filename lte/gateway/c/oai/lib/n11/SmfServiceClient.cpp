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

#include "SmfServiceClient.h"
#include "includes/ServiceRegistrySingleton.h"
#include <google/protobuf/util/time_util.h>

using grpc::Status;

namespace {

void SetAmfSessionContextRpcCallback(
    grpc::Status status, magma::lte::SmContextVoid response) {
  if (!status.ok()) {
    std::cout << "AsyncSetAmfSessionContext fails with code "
              << status.error_code() << ", msg: " << status.error_message()
              << std::endl;
  }
}

}  // namespace

using namespace magma::lte;
namespace magma5g {

AsyncSmfServiceClient::AsyncSmfServiceClient(
    std::shared_ptr<grpc::Channel> smf_srv_channel)
    : stub_(AmfPduSessionSmContext::NewStub(smf_srv_channel)) {}

AsyncSmfServiceClient::AsyncSmfServiceClient()
    : AsyncSmfServiceClient(
          magma::ServiceRegistrySingleton::Instance()->GetGrpcChannel(
              "sessiond", magma::ServiceRegistrySingleton::LOCAL)) {}

bool AsyncSmfServiceClient::set_smf_session(
    const SetSMSessionContext& request) {
  set_smf_session_rpc(request, SetAmfSessionContextRpcCallback);
  return true;
}

void AsyncSmfServiceClient::set_smf_session_rpc(
    const SetSMSessionContext& request,
    std::function<void(Status, SmContextVoid)> callback) {
  auto local_resp = new magma::AsyncLocalResponse<SmContextVoid>(
      std::move(callback), RESPONSE_TIMEOUT);
  local_resp->set_response_reader(std::move(stub_->AsyncSetAmfSessionContext(
      local_resp->get_context(), request, &queue_)));
}

}  // namespace magma5g

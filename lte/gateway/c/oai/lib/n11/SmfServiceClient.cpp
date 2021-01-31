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

#include <string>
#include "SmfServiceClient.h"
#include "ServiceRegistrySingleton.h"
#include <google/protobuf/util/time_util.h>

using grpc::Status;
namespace {

void callback_void_smf(
    grpc::Status status, magma::lte::SmContextVoid response) {
  if (!status.ok()) {
    //	OAILOG_ERROR( LOG_UTIL, "GRPC Failure Message: %s Status Error Code:
    //%d", status.error_message().c_str(), status.error_code());
    std::cout << "GRPC Failure Message: " << status.error_message().c_str()
              << " Status Error Code: " << to_string(status.error_code());
  }
}

}  // namespace

namespace magma5g {
using namespace magma::lte;

AsyncSmfServiceClient::AsyncSmfServiceClient(
    std::shared_ptr<grpc::Channel> smf_srv_channel)
    : stub_(AmfPduSessionSmContext::NewStub(smf_srv_channel)) {}

AsyncSmfServiceClient::AsyncSmfServiceClient()
    : AsyncSmfServiceClient(
          magma::ServiceRegistrySingleton::Instance()->GetGrpcChannel(
              "sessiond", magma::ServiceRegistrySingleton::LOCAL)) {}

bool AsyncSmfServiceClient::set_smf_session(
    const SetSMSessionContext& request) {
  set_smf_session_rpc(request, callback_void_smf);
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

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
#include <google/protobuf/util/time_util.h>
#include <cassert>
#include <grpcpp/impl/codegen/client_context.h>
#include <grpcpp/impl/codegen/status.h>
#include <cstring>
#include <iostream>
#include <memory>
#include <string>
#include <thread>
#include <google/protobuf/util/time_util.h>
#include <lte/protos/session_manager.grpc.pb.h>
#include <lte/protos/session_manager.pb.h>

#include "ServiceRegistrySingleton.h"
#include "SmfServiceClient.h"
using grpc::Status;
using magma::AsyncLocalResponse;
using magma::ServiceRegistrySingleton;

void handle_session_context_response(
    grpc::Status status, magma::lte::SmContextVoid response) {
  if (!status.ok()) {
    std::cout << "AsyncSetAmfSessionContext fails with code "
              << status.error_code() << ", msg: " << status.error_message()
              << std::endl;
  }
}

using namespace magma::lte;
namespace magma5g {

bool AsyncSmfServiceClient::set_smf_session(SetSMSessionContext& request) {
  SetSMFSessionRPC(
      request, [](const Status& status, const SmContextVoid& response) {
        handle_session_context_response(status, response);
      });

  return true;
}

void AsyncSmfServiceClient::SetSMFSessionRPC(
    SetSMSessionContext& request,
    const std::function<void(Status, SmContextVoid)>& callback) {
  auto localResp = new AsyncLocalResponse<SmContextVoid>(
      std::move(callback), RESPONSE_TIMEOUT);

  localResp->set_response_reader(std::move(stub_->AsyncSetAmfSessionContext(
      localResp->get_context(), request, &queue_)));
}

bool AsyncSmfServiceClient::set_smf_notification(
    SetSmNotificationContext& notify) {
  SetSMFNotificationRPC(
      notify, [](const Status& status, const SmContextVoid& response) {
        handle_session_context_response(status, response);
      });

  return true;
}

void AsyncSmfServiceClient::SetSMFNotificationRPC(
    SetSmNotificationContext& notify,
    const std::function<void(Status, SmContextVoid)>& callback) {
  auto localResp = new AsyncLocalResponse<SmContextVoid>(
      std::move(callback), RESPONSE_TIMEOUT);

  localResp->set_response_reader(std::move(stub_->AsyncSetSmfNotification(
      localResp->get_context(), notify, &queue_)));
}

AsyncSmfServiceClient::AsyncSmfServiceClient() {
  auto channel = ServiceRegistrySingleton::Instance()->GetGrpcChannel(
      "sessiond", ServiceRegistrySingleton::LOCAL);
  stub_ = AmfPduSessionSmContext::NewStub(channel);
  std::thread resp_loop_thread([&]() { rpc_response_loop(); });
  resp_loop_thread.detach();
}

AsyncSmfServiceClient& AsyncSmfServiceClient::getInstance() {
  static AsyncSmfServiceClient instance;
  return instance;
}

}  // namespace magma5g

/*
 * Licensed to the OpenAirInterface (OAI) Software Alliance under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The OpenAirInterface Software Alliance licenses this file to You under
 * the terms found in the LICENSE file in the root of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *-------------------------------------------------------------------------------
 * For more information about the OpenAirInterface (OAI) Software Alliance:
 *      contact@openairinterface.org
 */

#include <grpcpp/channel.h>
#include <grpcpp/impl/codegen/async_unary_call.h>
#include <thread>
#include <iostream>
#include <string>
#include <utility>

#include "orc8r/protos/mconfig/mconfigs.pb.h"
#include "PCEFClient.h"
#include "includes/ServiceRegistrySingleton.h"
#include "lte/protos/session_manager.pb.h"

namespace grpc {
class Status;
}  // namespace grpc
namespace magma {
namespace lte {
class SubscriberID;
}  // namespace lte
}  // namespace magma

#define MAGMAD_SERVICE "magmad"

using grpc::Status;

namespace magma {

PCEFClient& PCEFClient::get_instance() {
  static PCEFClient client_instance;
  return client_instance;
}

PCEFClient::PCEFClient() {
  // Create channel
  std::shared_ptr<Channel> channel;
  channel = ServiceRegistrySingleton::Instance()->GetGrpcChannel(
      "sessiond", ServiceRegistrySingleton::LOCAL);
  // Create stub for LocalSessionManager gRPC service
  stub_ = LocalSessionManager::NewStub(channel);
  std::thread resp_loop_thread([&]() { rpc_response_loop(); });
  std::cout << "PCEF Client thread id " << resp_loop_thread.native_handle();
  resp_loop_thread.detach();
}

void PCEFClient::create_session(
    const LocalCreateSessionRequest& request,
    std::function<void(Status, LocalCreateSessionResponse)> callback) {
  PCEFClient& client = get_instance();
  // Create a raw response pointer that stores a callback to be called when the
  // gRPC call is answered
  auto local_response = new AsyncLocalResponse<LocalCreateSessionResponse>(
      std::move(callback), RESPONSE_TIMEOUT);
  // Create a response reader for the `CreateSession` RPC call. This reader
  // stores the client context, the request to pass in, and the queue to add
  // the response to when done
  auto response_reader = client.stub_->AsyncCreateSession(
      local_response->get_context(), request, &client.queue_);
  // Set the reader for the local response. This executes the `CreateSession`
  // response using the response reader. When it is done, the callback stored in
  // `local_response` will be called
  local_response->set_response_reader(std::move(response_reader));
}

void PCEFClient::end_session(
    const LocalEndSessionRequest& request,
    std::function<void(Status, LocalEndSessionResponse)> callback) {
  PCEFClient& client  = get_instance();
  auto local_response = new AsyncLocalResponse<LocalEndSessionResponse>(
      std::move(callback), RESPONSE_TIMEOUT);
  auto response_reader = client.stub_->AsyncEndSession(
      local_response->get_context(), request, &client.queue_);
  local_response->set_response_reader(std::move(response_reader));
}

void PCEFClient::bind_policy2bearer(
    const PolicyBearerBindingRequest& request,
    std::function<void(Status, PolicyBearerBindingResponse)> callback) {
  PCEFClient& client  = get_instance();
  auto local_response = new AsyncLocalResponse<PolicyBearerBindingResponse>(
      std::move(callback), RESPONSE_TIMEOUT);
  auto response_reader = client.stub_->AsyncBindPolicy2Bearer(
      local_response->get_context(), request, &client.queue_);
  local_response->set_response_reader(std::move(response_reader));
}

void PCEFClient::update_teids(
    const UpdateTunnelIdsRequest& request,
    std::function<void(Status, UpdateTunnelIdsResponse)> callback) {
  PCEFClient& client  = get_instance();
  auto local_response = new AsyncLocalResponse<UpdateTunnelIdsResponse>(
      std::move(callback), RESPONSE_TIMEOUT);
  auto response_reader = client.stub_->AsyncUpdateTunnelIds(
      local_response->get_context(), request, &client.queue_);
  local_response->set_response_reader(std::move(response_reader));
}

}  // namespace magma

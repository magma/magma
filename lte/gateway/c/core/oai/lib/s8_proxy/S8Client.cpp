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

#include <grpcpp/impl/codegen/async_unary_call.h>
#include <thread>  // std::thread
#include <utility>

#include "S8Client.h"
#include "includes/ServiceRegistrySingleton.h"
#include "feg/protos/s8_proxy.pb.h"
#include "orc8r/protos/common.pb.h"

namespace grpc {
class Status;
}  // namespace grpc

namespace magma {

S8Client& S8Client::get_instance() {
  static S8Client client_instance;
  return client_instance;
}

S8Client::S8Client() {
  // Create channel
  auto channel = ServiceRegistrySingleton::Instance()->GetGrpcChannel(
      "s8_proxy", ServiceRegistrySingleton::CLOUD);
  // Create stub for s8_proxy gRPC service
  stub_ = S8Proxy::NewStub(channel);
  std::thread resp_loop_thread([&]() { rpc_response_loop(); });
  resp_loop_thread.detach();
}

void S8Client::s8_create_session_request(
    const CreateSessionRequestPgw& csr_req,
    std::function<void(grpc::Status, CreateSessionResponsePgw)> callback) {
  S8Client& client = get_instance();
  // Create a raw response pointer that stores a callback to be called when the
  // gRPC call is answered
  auto response = new AsyncLocalResponse<CreateSessionResponsePgw>(
      std::move(callback), RESPONSE_TIMEOUT);
  // Create a response reader for the `CreateSession` RPC call. This reader
  // stores the client context, the request to pass in, and the queue to add
  // the response to when done
  auto response_reader = client.stub_->AsyncCreateSession(
      response->get_context(), csr_req, &client.queue_);
  // Set the reader for the local response. This executes the `CreateSession`
  // response using the response reader. When it is done, the callback stored in
  // `local_response` will be called
  response->set_response_reader(std::move(response_reader));
}

void S8Client::s8_delete_session_request(
    const DeleteSessionRequestPgw& dsr_req,
    std::function<void(grpc::Status, DeleteSessionResponsePgw)> callback) {
  S8Client& client = get_instance();
  // Create a raw response pointer that stores a callback to be called when the
  // gRPC call is answered
  auto response = new AsyncLocalResponse<DeleteSessionResponsePgw>(
      std::move(callback), RESPONSE_TIMEOUT);
  // Create a response reader for the `DeleteSession` RPC call. This reader
  // stores the client context, the request to pass in, and the queue to add
  // the response to when done
  auto response_reader = client.stub_->AsyncDeleteSession(
      response->get_context(), dsr_req, &client.queue_);
  // Set the reader for the local response. This executes the `DeleteSession`
  // response using the response reader. When it is done, the callback stored in
  // `local_response` will be called
  response->set_response_reader(std::move(response_reader));
}

}  // namespace magma

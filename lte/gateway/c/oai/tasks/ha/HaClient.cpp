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
#include <iostream>
#include <utility>

#include "HaClient.h"
#include "includes/ServiceRegistrySingleton.h"
#include "lte/protos/ha_orc8r.pb.h"

namespace magma {

HaClient& HaClient::get_instance() {
  static HaClient client_instance;
  return client_instance;
}

HaClient::HaClient() {
  auto channel = ServiceRegistrySingleton::Instance()->GetGrpcChannel(
      "ha", ServiceRegistrySingleton::CLOUD);
  // Create stub for HaProxy gRPC service
  stub_ = lte::Ha::NewStub(channel);

  std::thread resp_loop_thread([&]() { rpc_response_loop(); });
  resp_loop_thread.detach();
}

void HaClient::get_eNB_offload_state(
    std::function<void(grpc::Status, lte::GetEnodebOffloadStateResponse)>
        callback) {
  HaClient& client = get_instance();

  // Create a raw response pointer that stores a callback to be called when the
  // gRPC call is answered
  auto resp = new AsyncLocalResponse<lte::GetEnodebOffloadStateResponse>(
      std::move(callback), RESPONSE_TIMEOUT);

  // Create a response reader for the `GetEnodebOffloadStateRequest` RPC call.
  // This reader stores the client context, the request to pass in,
  // and the queue to add the response to when done
  lte::GetEnodebOffloadStateRequest request;

  auto resp_reader = client.stub_->AsyncGetEnodebOffloadState(
      resp->get_context(), request, &client.queue_);

  // Set the reader for the response. This executes the `GetEnodebOffloadState`
  // response using the response reader. When it is done, the callback stored
  // in `resp` will be called
  resp->set_response_reader(std::move(resp_reader));
}

}  // namespace magma

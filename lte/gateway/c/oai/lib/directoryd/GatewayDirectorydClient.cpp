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
#include "GatewayDirectorydClient.h"

#include <memory>
#include <thread>
#include <utility>

#include <grpcpp/impl/codegen/async_unary_call.h>
#include <google/protobuf/map.h>

#include "orc8r/protos/common.pb.h"
#include "includes/ServiceRegistrySingleton.h"

namespace grpc {
class Channel;
class ClientContext;
class Status;
}  // namespace grpc

using google::protobuf::Map;
using grpc::Channel;
using grpc::ClientContext;
using grpc::Status;
using magma::GatewayDirectoryService;
using magma::GatewayDirectoryServiceClient;
using magma::UpdateRecordRequest;
using magma::orc8r::Void;

GatewayDirectoryServiceClient& GatewayDirectoryServiceClient::get_instance() {
  static GatewayDirectoryServiceClient client_instance;
  return client_instance;
}

GatewayDirectoryServiceClient::GatewayDirectoryServiceClient() {
  auto channel = ServiceRegistrySingleton::Instance()->GetGrpcChannel(
      "directoryd", ServiceRegistrySingleton::LOCAL);
  stub_ = GatewayDirectoryService::NewStub(channel);
  std::thread resp_loop_thread([&]() { rpc_response_loop(); });
  resp_loop_thread.detach();
}

bool GatewayDirectoryServiceClient::UpdateRecord(
    const std::string& id, const std::string& location,
    std::function<void(Status, Void)> callback) {
  UpdateRecordRequest request;
  request.set_id(id);
  request.set_location(location);
  return GatewayDirectoryServiceClient::updateRecordImpl(request, callback);
}

bool GatewayDirectoryServiceClient::UpdateRecordField(
    const std::string& id, const std::string& field_key,
    const std::string& field_value,
    std::function<void(Status, Void)> callback) {
  UpdateRecordRequest request;
  request.set_id(id);
  auto update_fields = request.mutable_fields();
  update_fields->insert({field_key, field_value});
  return GatewayDirectoryServiceClient::updateRecordImpl(request, callback);
}

bool GatewayDirectoryServiceClient::updateRecordImpl(
    UpdateRecordRequest& request, std::function<void(Status, Void)> callback) {
  GatewayDirectoryServiceClient& client = get_instance();
  // Create a raw response pointer that stores a callback to be called when the
  // gRPC call is answered
  auto local_response =
      new AsyncLocalResponse<Void>(std::move(callback), RESPONSE_TIMEOUT);
  // Create a response reader for the `UpdateRecord` RPC call. This reader
  // stores the client context, the request to pass in, and the queue to add
  // the response to when done
  auto response_reader = client.stub_->AsyncUpdateRecord(
      local_response->get_context(), request, &client.queue_);
  // Set the reader for the local response. This executes the `UpdateRecord`
  // response using the response reader. When it is done, the callback stored in
  // `local_response` will be called
  local_response->set_response_reader(std::move(response_reader));
  return true;
}

bool GatewayDirectoryServiceClient::DeleteRecord(
    const std::string& id, std::function<void(Status, Void)> callback) {
  DeleteRecordRequest request;
  Void response;

  request.set_id(id);

  auto local_response =
      new AsyncLocalResponse<Void>(std::move(callback), RESPONSE_TIMEOUT);
  GatewayDirectoryServiceClient& client = get_instance();
  auto response_reader                  = client.stub_->AsyncDeleteRecord(
      local_response->get_context(), request, &client.queue_);
  local_response->set_response_reader(std::move(response_reader));
  return true;
}

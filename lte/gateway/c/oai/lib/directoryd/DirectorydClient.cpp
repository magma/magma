/*
 * Licensed to the OpenAirInterface (OAI) Software Alliance under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The OpenAirInterface Software Alliance licenses this file to You under
 * the Apache License, Version 2.0  (the "License"); you may not use this file
 * except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
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

#include <grpcpp/impl/codegen/async_unary_call.h>
#include <memory>
#include <thread>
#include <utility>

#include "orc8r/protos/common.pb.h"
#include "DirectorydClient.h"
#include "ServiceRegistrySingleton.h"

namespace grpc {
class Channel;
class ClientContext;
class Status;
}  // namespace grpc

using grpc::Channel;
using grpc::ClientContext;
using grpc::Status;
using magma::DirectoryService;
using magma::DirectoryServiceClient;
using magma::UpdateDirectoryLocationRequest;
using magma::orc8r::Void;

DirectoryServiceClient &DirectoryServiceClient::get_instance()
{
  static DirectoryServiceClient client_instance;
  return client_instance;
}

DirectoryServiceClient::DirectoryServiceClient()
{
  auto channel = ServiceRegistrySingleton::Instance()->GetGrpcChannel(
    "directoryd", ServiceRegistrySingleton::LOCAL);
  stub_ = DirectoryService::NewStub(channel);
  std::thread resp_loop_thread([&]() { rpc_response_loop(); });
  resp_loop_thread.detach();
}

bool DirectoryServiceClient::UpdateLocation(
  TableID table,
  const std::string &id,
  const std::string &location,
  std::function<void(Status, Void)> callback)
{
  DirectoryServiceClient &client = get_instance();

  UpdateDirectoryLocationRequest request;
  Void response;

  request.set_table(table);
  request.set_id(id);
  request.mutable_record()->set_location(location);

  // Create a raw response pointer that stores a callback to be called when the
  // gRPC call is answered
  auto local_response =
    new AsyncLocalResponse<Void>(std::move(callback), RESPONSE_TIMEOUT);
  // Create a response reader for the `UpdateLocation` RPC call. This reader
  // stores the client context, the request to pass in, and the queue to add
  // the response to when done
  auto response_reader = client.stub_->AsyncUpdateLocation(
    local_response->get_context(), request, &client.queue_);
  // Set the reader for the local response. This executes the `UpdateLocation`
  // response using the response reader. When it is done, the callback stored in
  // `local_response` will be called
  local_response->set_response_reader(std::move(response_reader));
  return true;
}

bool DirectoryServiceClient::DeleteLocation(
  TableID table,
  const std::string &id,
  std::function<void(Status, Void)> callback)
{
  DeleteLocationRequest request;
  Void response;

  request.set_table(table);
  request.set_id(id);

  auto local_response =
    new AsyncLocalResponse<Void>(std::move(callback), RESPONSE_TIMEOUT);
  DirectoryServiceClient &client = get_instance();
  auto response_reader = client.stub_->AsyncDeleteLocation(
    local_response->get_context(), request, &client.queue_);
  local_response->set_response_reader(std::move(response_reader));
  return true;
}

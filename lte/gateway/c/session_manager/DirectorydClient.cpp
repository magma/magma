/**
 * Copyright (c) 2016-present, Facebook, Inc.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include "DirectorydClient.h"

#include "ServiceRegistrySingleton.h"
#include "magma_logging.h"

using grpc::Status;

namespace { // anonymous

magma::GetDirectoryFieldRequest create_directory_field_req(
  const std::string& imsi
  )
{
  magma::GetDirectoryFieldRequest req;
  req.set_id(imsi);
  req.set_field_key("ipv4_addr");
  return req;
}

} // namespace

namespace magma {

AsyncDirectorydClient::AsyncDirectorydClient(
  std::shared_ptr<grpc::Channel> channel):
  stub_(GatewayDirectoryService::NewStub(channel))
{
}

AsyncDirectorydClient::AsyncDirectorydClient():
  AsyncDirectorydClient(ServiceRegistrySingleton::Instance()->GetGrpcChannel(
    "directoryd",
    ServiceRegistrySingleton::LOCAL))
{
}

bool AsyncDirectorydClient::get_directoryd_ip_field(
  const std::string& imsi,
  std::function<void(Status status, DirectoryField)> callback)
{
  auto req = create_directory_field_req(imsi);
  get_directoryd_ip_field_rpc(req, callback);
  return true;
}

void AsyncDirectorydClient::update_directoryd_record(
  const UpdateRecordRequest& request,
  std::function<void(Status status, Void)> callback) {
  auto local_response =
    new AsyncLocalResponse<Void>(std::move(callback), RESPONSE_TIMEOUT);
  local_response->set_response_reader(std::move(
        stub_->AsyncUpdateRecord(local_response->get_context(),
                                      request, &queue_)));
}

void AsyncDirectorydClient::get_directoryd_ip_field_rpc(
  const GetDirectoryFieldRequest& request,
  std::function<void(Status, DirectoryField)> callback)
{
  auto local_resp = new AsyncLocalResponse<DirectoryField>(
    std::move(callback), RESPONSE_TIMEOUT);
  local_resp->set_response_reader(std::move(
    stub_->AsyncGetDirectoryField(local_resp->get_context(),
      request, &queue_)));
}

bool AsyncDirectorydClient::delete_directoryd_record(
  const DeleteRecordRequest& request,
  std::function<void(Status status, Void)> callback) {
  auto local_response =
    new AsyncLocalResponse<Void>(std::move(callback), RESPONSE_TIMEOUT);
  local_response->set_response_reader(std::move(
        stub_->AsyncDeleteRecord(local_response->get_context(),
                                      request, &queue_)));
  return true;
}


void AsyncDirectorydClient::get_all_directoryd_records(
  std::function<void(Status status, AllDirectoryRecords)> callback)
{
  magma::Void request;
  auto local_resp = new AsyncLocalResponse<AllDirectoryRecords>(
    std::move(callback), RESPONSE_TIMEOUT);
  local_resp->set_response_reader(std::move(
    stub_->AsyncGetAllDirectoryRecords(local_resp->get_context(),
      request, &queue_)));
}
} // namespace magma

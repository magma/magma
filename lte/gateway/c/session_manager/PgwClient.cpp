/**
 * Copyright (c) 2016-present, Facebook, Inc.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include "PgwClient.h"
#include "ServiceRegistrySingleton.h"
#include "magma_logging.h"

// using google::protobuf::RepeatedPtrField;
using grpc::Status;

namespace { // anonymous

magma::DeleteBearerRequest create_delete_bearer_req(
  const std::string &imsi,
  const std::string &apn_ip_addr,
  const uint32_t link_bearer_id,
  const std::vector<magma::PolicyRule> &flows)
{
  magma::DeleteBearerRequest req;
  req.mutable_sid()->set_id(imsi);
  req.set_ip_addr(apn_ip_addr);
  req.set_link_bearer_id(link_bearer_id);

  auto req_policy_rules = req.mutable_policy_rules();
  for (const auto &flow : flows) {
    req_policy_rules->Add()->CopyFrom(flow);
  }

  return req;
}

magma::CreateBearerRequest create_add_bearer_req(
  const std::string &imsi,
  const std::string &apn_ip_addr,
  const uint32_t link_bearer_id,
  const std::vector<magma::PolicyRule> &flows)
{
  magma::CreateBearerRequest req;
  req.mutable_sid()->set_id(imsi);
  req.set_ip_addr(apn_ip_addr);
  req.set_link_bearer_id(link_bearer_id);

  auto req_policy_rules = req.mutable_policy_rules();
  for (const auto &flow : flows) {
    req_policy_rules->Add()->CopyFrom(flow);
  }

  return req;
}

} // namespace

namespace magma {

AsyncPgwClient::AsyncPgwClient(
  std::shared_ptr<grpc::Channel> channel):
  stub_(Pgw::NewStub(channel))
{
}

AsyncPgwClient::AsyncPgwClient():
  AsyncPgwClient(ServiceRegistrySingleton::Instance()->GetGrpcChannel(
    "pgw",
    ServiceRegistrySingleton::LOCAL))
{
}

bool AsyncPgwClient::delete_dedicated_bearer(
  const std::string &imsi,
  const std::string &apn_ip_addr,
  const uint32_t link_bearer_id,
  const std::vector<PolicyRule> &flows)
{
  auto req = create_delete_bearer_req(imsi, apn_ip_addr, link_bearer_id, flows);
  MLOG(MDEBUG) << "deleting dedicated bearer "
               << imsi << apn_ip_addr;
  delete_dedicated_bearer_rpc(
      req, [imsi, apn_ip_addr](Status status, DeleteBearerResult resp) {
    if (!status.ok()) {
      MLOG(MERROR) << "Could not delete dedicated bearer" << imsi << apn_ip_addr
                   << ": " << status.error_message();
    }
  });
  return true;
}

bool AsyncPgwClient::create_dedicated_bearer(
  const std::string &imsi,
  const std::string &apn_ip_addr,
  const uint32_t link_bearer_id,
  const std::vector<PolicyRule> &flows)
{
  auto req = create_add_bearer_req(imsi, apn_ip_addr, link_bearer_id, flows);
  MLOG(MDEBUG) << "creating dedicated bearer "
               << imsi << apn_ip_addr;
  create_dedicated_bearer_rpc(
      req, [imsi, apn_ip_addr](Status status, CreateBearerResult resp) {
    if (!status.ok()) {
      MLOG(MERROR) << "Could not create dedicated bearer" << imsi << apn_ip_addr
                   << ": " << status.error_message();
    }
  });
  return true;
}

void AsyncPgwClient::delete_dedicated_bearer_rpc(
  const DeleteBearerRequest &request,
  std::function<void(Status, DeleteBearerResult)> callback)
{
  auto local_resp = new AsyncLocalResponse<DeleteBearerResult>(
    std::move(callback), RESPONSE_TIMEOUT);
  local_resp->set_response_reader(std::move(
    stub_->AsyncDeleteBearer(local_resp->get_context(), request, &queue_)));
}

void AsyncPgwClient::create_dedicated_bearer_rpc(
  const CreateBearerRequest &request,
  std::function<void(Status, CreateBearerResult)> callback)
{
  auto local_resp = new AsyncLocalResponse<CreateBearerResult>(
    std::move(callback), RESPONSE_TIMEOUT);
  local_resp->set_response_reader(std::move(
    stub_->AsyncCreateBearer(local_resp->get_context(), request, &queue_)));
}

} // namespace magma

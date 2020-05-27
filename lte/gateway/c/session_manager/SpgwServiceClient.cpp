/**
 * Copyright (c) 2016-present, Facebook, Inc.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include "SpgwServiceClient.h"
#include "ServiceRegistrySingleton.h"
#include "magma_logging.h"

// using google::protobuf::RepeatedPtrField;
using grpc::Status;

namespace { // anonymous

magma::DeleteBearerRequest
create_delete_bearer_req(const std::string &imsi,
                         const std::string &apn_ip_addr,
                         const uint32_t linked_bearer_id,
                         const std::vector<uint32_t> &eps_bearer_ids) {
  magma::DeleteBearerRequest req;
  req.mutable_sid()->set_id(imsi);
  req.set_ip_addr(apn_ip_addr);
  req.set_link_bearer_id(linked_bearer_id);

  auto ebis = req.mutable_eps_bearer_ids();
  for (const auto &eps_bearer_id : eps_bearer_ids) {
    ebis->Add(eps_bearer_id);
  }

  return req;
}

magma::CreateBearerRequest
create_add_bearer_req(const std::string &imsi, const std::string &apn_ip_addr,
                      const uint32_t linked_bearer_id,
                      const std::vector<magma::PolicyRule> &flows) {
  magma::CreateBearerRequest req;
  req.mutable_sid()->set_id(imsi);
  req.set_ip_addr(apn_ip_addr);
  req.set_link_bearer_id(linked_bearer_id);

  auto req_policy_rules = req.mutable_policy_rules();
  for (const auto &flow : flows) {
    req_policy_rules->Add()->CopyFrom(flow);
  }

  return req;
}

} // namespace

namespace magma {

AsyncSpgwServiceClient::AsyncSpgwServiceClient(
    std::shared_ptr<grpc::Channel> channel)
    : stub_(SpgwService::NewStub(channel)) {}

AsyncSpgwServiceClient::AsyncSpgwServiceClient()
    : AsyncSpgwServiceClient(
          ServiceRegistrySingleton::Instance()->GetGrpcChannel(
              "spgw_service", ServiceRegistrySingleton::LOCAL)) {}

bool AsyncSpgwServiceClient::delete_default_bearer(
    const std::string &imsi, const std::string &apn_ip_addr,
    const uint32_t linked_bearer_id) {
  MLOG(MINFO) << "Deleting default bearer and corresponding PDN session for"
              << " IMSI: " << imsi << " APN IP addr " << apn_ip_addr
              << " Bearer ID " << linked_bearer_id;
  std::vector<uint32_t> eps_bearer_ids = {linked_bearer_id};
  return delete_bearer(imsi, apn_ip_addr, linked_bearer_id, eps_bearer_ids);
}

bool AsyncSpgwServiceClient::delete_dedicated_bearer(
    const std::string &imsi, const std::string &apn_ip_addr,
    const uint32_t linked_bearer_id,
    const std::vector<uint32_t> &eps_bearer_ids) {
  MLOG(MINFO) << "Deleting dedicated bearer IMSI: " << imsi << " APN IP addr "
              << apn_ip_addr << " Bearer ID " << linked_bearer_id;
  return delete_bearer(imsi, apn_ip_addr, linked_bearer_id, eps_bearer_ids);
}

bool AsyncSpgwServiceClient::create_dedicated_bearer(
    const std::string &imsi, const std::string &apn_ip_addr,
    const uint32_t linked_bearer_id, const std::vector<PolicyRule> &flows) {
  auto req = create_add_bearer_req(imsi, apn_ip_addr, linked_bearer_id, flows);
  MLOG(MINFO) << "creating dedicated bearer " << imsi << apn_ip_addr;
  create_dedicated_bearer_rpc(
      req, [imsi, apn_ip_addr](Status status, CreateBearerResult resp) {
        if (!status.ok()) {
          MLOG(MERROR) << "Could not create dedicated bearer" << imsi
                       << apn_ip_addr << ": " << status.error_message();
        }
      });
  return true;
}

// delete_bearer creates the DeleteBearerRequest and logs the error
bool AsyncSpgwServiceClient::delete_bearer(
    const std::string &imsi, const std::string &apn_ip_addr,
    const uint32_t linked_bearer_id,
    const std::vector<uint32_t> &eps_bearer_ids) {
  auto req = create_delete_bearer_req(imsi, apn_ip_addr, linked_bearer_id,
                                      eps_bearer_ids);
  delete_bearer_rpc(
      req, [imsi, apn_ip_addr](Status status, DeleteBearerResult resp) {
        if (!status.ok()) {
          // only log error for now
          MLOG(MERROR) << "Could not delete bearer" << imsi << apn_ip_addr
                       << ": " << status.error_message();
        }
      });
  return true;
}

void AsyncSpgwServiceClient::delete_bearer_rpc(
    const DeleteBearerRequest &request,
    std::function<void(Status, DeleteBearerResult)> callback) {
  auto local_resp = new AsyncLocalResponse<DeleteBearerResult>(
      std::move(callback), RESPONSE_TIMEOUT);
  local_resp->set_response_reader(std::move(
      stub_->AsyncDeleteBearer(local_resp->get_context(), request, &queue_)));
}

void AsyncSpgwServiceClient::create_dedicated_bearer_rpc(
    const CreateBearerRequest &request,
    std::function<void(Status, CreateBearerResult)> callback) {
  auto local_resp = new AsyncLocalResponse<CreateBearerResult>(
      std::move(callback), RESPONSE_TIMEOUT);
  local_resp->set_response_reader(std::move(
      stub_->AsyncCreateBearer(local_resp->get_context(), request, &queue_)));
}

} // namespace magma

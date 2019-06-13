/**
 * Copyright (c) 2016-present, Facebook, Inc.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include "PipelinedClient.h"

#include "ServiceRegistrySingleton.h"
#include "magma_logging.h"

using grpc::Status;

namespace { // anonymous

magma::DeactivateFlowsRequest create_deactivate_req(
  const std::string &imsi,
  const std::vector<std::string> &rule_ids,
  const std::vector<magma::PolicyRule> &dynamic_rules)
{
  magma::DeactivateFlowsRequest req;
  req.mutable_sid()->set_id(imsi);
  auto ids = req.mutable_rule_ids();
  for (const auto &id : rule_ids) {
    ids->Add()->assign(id);
  }
  for (const auto &rule : dynamic_rules) {
    ids->Add()->assign(rule.id());
  }
  return req;
}

magma::ActivateFlowsRequest create_activate_req(
  const std::string &imsi,
  const std::string &ip_addr,
  const std::vector<std::string> &static_rules,
  const std::vector<magma::PolicyRule> &dynamic_rules)
{
  magma::ActivateFlowsRequest req;
  req.mutable_sid()->set_id(imsi);
  req.set_ip_addr(ip_addr);
  auto ids = req.mutable_rule_ids();
  for (const auto &id : static_rules) {
    ids->Add()->assign(id);
  }
  auto mut_dyn_rules = req.mutable_dynamic_rules();
  for (const auto &dyn_rule : dynamic_rules) {
    mut_dyn_rules->Add()->CopyFrom(dyn_rule);
  }
  return req;
}

} // namespace

namespace magma {

AsyncPipelinedClient::AsyncPipelinedClient(
  std::shared_ptr<grpc::Channel> channel):
  stub_(Pipelined::NewStub(channel))
{
}

AsyncPipelinedClient::AsyncPipelinedClient():
  AsyncPipelinedClient(ServiceRegistrySingleton::Instance()->GetGrpcChannel(
    "pipelined",
    ServiceRegistrySingleton::LOCAL))
{
}

bool AsyncPipelinedClient::deactivate_all_flows(const std::string &imsi)
{
  DeactivateFlowsRequest req;
  req.mutable_sid()->set_id(imsi);
  MLOG(MDEBUG) << "Deactivating all flows for subscriber " << imsi;
  deactivate_flows_rpc(req, [imsi](Status status, DeactivateFlowsResult resp) {
    if (!status.ok()) {
      MLOG(MERROR) << "Could not deactivate flows for subscriber " << imsi
                   << ": " << status.error_message();
    }
  });
  return true;
}

bool AsyncPipelinedClient::deactivate_flows_for_rules(
  const std::string &imsi,
  const std::vector<std::string> &rule_ids,
  const std::vector<PolicyRule> &dynamic_rules)
{
  auto req = create_deactivate_req(imsi, rule_ids, dynamic_rules);
  MLOG(MDEBUG) << "Deactivating " << rule_ids.size() << " static rules and "
               << dynamic_rules.size() << " dynamic rules for subscriber "
               << imsi;
  deactivate_flows_rpc(req, [imsi](Status status, DeactivateFlowsResult resp) {
    if (!status.ok()) {
      MLOG(MERROR) << "Could not deactivate flows for subscriber " << imsi
                   << ": " << status.error_message();
    }
  });
  return true;
}

bool AsyncPipelinedClient::activate_flows_for_rules(
  const std::string &imsi,
  const std::string &ip_addr,
  const std::vector<std::string> &static_rules,
  const std::vector<PolicyRule> &dynamic_rules)
{
  auto req = create_activate_req(imsi, ip_addr, static_rules, dynamic_rules);
  MLOG(MDEBUG) << "Activating " << static_rules.size() << " static rules and "
               << dynamic_rules.size() << " dynamic rules for subscriber "
               << imsi;
  activate_flows_rpc(req, [imsi](Status status, ActivateFlowsResult resp) {
    if (!status.ok()) {
      MLOG(MERROR) << "Could not activate flows through pipelined for UE "
                   << imsi << ": " << status.error_message();
    }
  });
  return true;
}

void AsyncPipelinedClient::deactivate_flows_rpc(
  const DeactivateFlowsRequest &request,
  std::function<void(Status, DeactivateFlowsResult)> callback)
{
  auto local_resp = new AsyncLocalResponse<DeactivateFlowsResult>(
    std::move(callback), RESPONSE_TIMEOUT);
  local_resp->set_response_reader(std::move(
    stub_->AsyncDeactivateFlows(local_resp->get_context(), request, &queue_)));
}

void AsyncPipelinedClient::activate_flows_rpc(
  const ActivateFlowsRequest &request,
  std::function<void(Status, ActivateFlowsResult)> callback)
{
  auto local_resp = new AsyncLocalResponse<ActivateFlowsResult>(
    std::move(callback), RESPONSE_TIMEOUT);
  local_resp->set_response_reader(std::move(
    stub_->AsyncActivateFlows(local_resp->get_context(), request, &queue_)));
}

} // namespace magma

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

magma::UEMacFlowRequest create_add_ue_mac_flow_req(
  const magma::SubscriberID &sid,
  const std::string &mac_addr)
{
  magma::UEMacFlowRequest req;
  req.mutable_sid()->CopyFrom(sid);
  req.set_mac_addr(mac_addr);
  return req;
}

magma::SetupFlowsRequest create_setup_flows_req(
  const std::vector<magma::SessionState::SessionInfo> &infos,
  const std::uint64_t &epoch)
{
  magma::SetupFlowsRequest req;
  std::vector<magma::ActivateFlowsRequest> activation_reqs;
  for(auto it = infos.begin(); it != infos.end(); it++ )
  {
    auto activate_req = create_activate_req(it->imsi, it->ip_addr,
      it->static_rules, it->dynamic_rules);
    activation_reqs.push_back(activate_req);
  }
  auto mut_requests = req.mutable_requests();
  for (const auto &act_req : activation_reqs) {
    mut_requests->Add()->CopyFrom(act_req);
  }
  req.set_epoch(epoch);
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

bool AsyncPipelinedClient::setup(
   const std::vector<SessionState::SessionInfo> &infos,
   const std::uint64_t &epoch,
   std::function<void(Status status, SetupFlowsResult)> callback)
{
  SetupFlowsRequest setup_req = create_setup_flows_req(infos, epoch);
  setup_flows_rpc(setup_req, callback);
  return true;
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
  MLOG(MDEBUG) << "Activating " << static_rules.size() << " static rules and "
               << dynamic_rules.size() << " dynamic rules for subscriber "
               << imsi;
  // Activate static rules and dynamic rules separately until bug is fixed in
  // pipelined which crashes if activated at the same time
  auto static_req = create_activate_req(
    imsi, ip_addr, static_rules, std::vector<PolicyRule>());
  activate_flows_rpc(static_req,
    [imsi](Status status, ActivateFlowsResult resp) {
      if (!status.ok()) {
        MLOG(MERROR) << "Could not activate flows through pipelined for UE "
                     << imsi << ": " << status.error_message();
      }
  });
  auto dynamic_req = create_activate_req(
    imsi, ip_addr,std::vector<std::string>(), dynamic_rules);
  activate_flows_rpc(dynamic_req,
    [imsi](Status status, ActivateFlowsResult resp) {
      if (!status.ok()) {
        MLOG(MERROR) << "Could not activate flows through pipelined for UE "
                     << imsi << ": " << status.error_message();
      }
  });
  return true;
}

bool AsyncPipelinedClient::add_ue_mac_flow(
    const SubscriberID &sid,
    const std::string &mac_addr)
{
  auto req = create_add_ue_mac_flow_req(sid, mac_addr);
  add_ue_mac_flow_rpc(req, [mac_addr](Status status, FlowResponse resp) {
    if (!status.ok()) {
      MLOG(MERROR) << "Could not add flow for subscriber with UE MAC"
                   << mac_addr << ": " << status.error_message();
    }
  });
  return true;
}

void AsyncPipelinedClient::setup_flows_rpc(
  const SetupFlowsRequest &request,
  std::function<void(Status, SetupFlowsResult)> callback)
{
  auto local_resp = new AsyncLocalResponse<SetupFlowsResult>(
    std::move(callback), RESPONSE_TIMEOUT);
  local_resp->set_response_reader(std::move(
    stub_->AsyncSetupFlows(local_resp->get_context(), request, &queue_)));
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

void AsyncPipelinedClient::add_ue_mac_flow_rpc(
    const UEMacFlowRequest &request,
    std::function<void(Status, FlowResponse)> callback)
{
  auto local_resp = new AsyncLocalResponse<FlowResponse>(
    std::move(callback), RESPONSE_TIMEOUT);
  local_resp->set_response_reader(std::move(
    stub_->AsyncAddUEMacFlow(local_resp->get_context(), request, &queue_)));
}

} // namespace magma

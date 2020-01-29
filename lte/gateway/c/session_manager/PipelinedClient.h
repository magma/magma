/**
 * Copyright (c) 2016-present, Facebook, Inc.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */
#pragma once

#include <mutex>
#include <unordered_map>

#include <lte/protos/policydb.pb.h>
#include <lte/protos/pipelined.grpc.pb.h>
#include <lte/protos/subscriberdb.pb.h>

#include "GRPCReceiver.h"
#include "SessionState.h"

using grpc::Status;

namespace magma {
using namespace lte;

/**
 * PipelinedClient is the base class for managing rules and their activations.
 * The class is intended on interfacing with the data pipeline to enforce rules.
 */
class PipelinedClient {
 public:
  /**
   * Activates all rules for provided SessionInfos
   * @param infos - list of SessionInfos to setup flows for
   * @return true if the operation was successful
   */
  virtual bool setup(
    const std::vector<SessionState::SessionInfo> &infos,
    const std::uint64_t &epoch,
    std::function<void(Status status, SetupFlowsResult)> callback) = 0;

  /**
   * Deactivate all flows for a subscriber's session
   * @param imsi - UE to delete all policy flows for
   * @return true if the operation was successful
   */
  virtual bool deactivate_all_flows(const std::string &imsi) = 0;

  /**
   * Deactivate all flows for the specified rules
   * @param imsi - UE to delete flows for
   * @param rule_ids - rules to deactivate
   * @return true if the operation was successful
   */
  virtual bool deactivate_flows_for_rules(
    const std::string &imsi,
    const std::vector<std::string> &rule_ids,
    const std::vector<PolicyRule> &dynamic_rules) = 0;

  /**
   * Activate all rules for the specified rules, using a normal vector
   */
  virtual bool activate_flows_for_rules(
    const std::string &imsi,
    const std::string &ip_addr,
    const std::vector<std::string> &static_rules,
    const std::vector<PolicyRule> &dynamic_rules) = 0;

  /**
   * Send the MAC address of UE and the subscriberID
   * for pipelined to add a flow for the subscriber by matching the MAC
   */
  virtual bool add_ue_mac_flow(
    const SubscriberID &sid,
    const std::string &ue_mac_addr,
    const std::string &msisdn,
    const std::string &ap_mac_addr,
    const std::string &ap_name) = 0;
};

/**
 * AsyncPipelinedClient implements PipelinedClient but sends calls
 * asynchronously to pipelined.
 */
class AsyncPipelinedClient : public GRPCReceiver, public PipelinedClient {
 public:
  AsyncPipelinedClient();

  AsyncPipelinedClient(std::shared_ptr<grpc::Channel> pipelined_channel);

  /**
   * Activates all rules for provided SessionInfos
   * @param infos - list of SessionInfos to setup flows for
   * @return true if the operation was successful
   */
  bool setup(
    const std::vector<SessionState::SessionInfo> &infos,
    const std::uint64_t &epoch,
    std::function<void(Status status, SetupFlowsResult)> callback);

  /**
   * Deactivate all flows for a subscriber's session
   * @param imsi - UE to delete all policy flows for
   * @return true if the operation was successful
   */
  bool deactivate_all_flows(const std::string &imsi);

  /**
   * Deactivate all flows related to a specific charging key
   * @param imsi - UE to delete flows for
   * @param charging_key - key to deactivate
   * @return true if the operation was successful
   */
  bool deactivate_flows_for_rules(
    const std::string &imsi,
    const std::vector<std::string> &rule_ids,
    const std::vector<PolicyRule> &dynamic_rules);

  /**
   * Activate all rules for the specified rules, using a normal vector
   */
  bool activate_flows_for_rules(
    const std::string &imsi,
    const std::string &ip_addr,
    const std::vector<std::string> &static_rules,
    const std::vector<PolicyRule> &dynamic_rules);

  bool add_ue_mac_flow(
    const SubscriberID &sid,
    const std::string &ue_mac_addr,
    const std::string &msisdn,
    const std::string &ap_mac_addr,
    const std::string &ap_name);

 private:
  static const uint32_t RESPONSE_TIMEOUT = 6; // seconds
  std::unique_ptr<Pipelined::Stub> stub_;

 private:
  void setup_flows_rpc(
    const SetupFlowsRequest &request,
    std::function<void(Status, SetupFlowsResult)> callback);

  void deactivate_flows_rpc(
    const DeactivateFlowsRequest &request,
    std::function<void(Status, DeactivateFlowsResult)> callback);

  void activate_flows_rpc(
    const ActivateFlowsRequest &request,
    std::function<void(Status, ActivateFlowsResult)> callback);

  void add_ue_mac_flow_rpc(
    const UEMacFlowRequest &request,
    std::function<void(Status, FlowResponse)> callback);
};

} // namespace magma

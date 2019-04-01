/**
 * Copyright (c) 2016-present, Facebook, Inc.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */
#pragma once

#include <lte/protos/session_manager.grpc.pb.h>
#include <folly/io/async/EventBaseManager.h>

#include "RuleStore.h"
#include "PipelinedClient.h"
#include "SessionState.h"

namespace magma {
using namespace orc8r;

class SessionNotFound : public std::exception {
 public:
  SessionNotFound() = default;
};

/**
 * LocalEnforcer can register traffic records and credits to track when a flow
 * has run out of credit
 */
class LocalEnforcer {
 public:
  LocalEnforcer();

  LocalEnforcer(
    std::shared_ptr<StaticRuleStore> rule_store,
    std::shared_ptr<PipelinedClient> pipelined_client);

  void attachEventBase(folly::EventBase *evb);

  // blocks
  void start();

  void startOnce(folly::EventBase *evb);

  void stop();

  folly::EventBase &get_event_base();

  /**
   * Insert a group of rule usage into the monitor and update credit manager
   * Assumes records are aggregates, as in the usages sent are cumulative and
   * not differences.
   *
   * @param records - a RuleRecordTable protobuf with a vector of RuleRecords
   */
  void aggregate_records(const RuleRecordTable &records);

  /**
   * reset_updates resets all of the charging keys being updated in
   * failed_request. This should only be called if the *entire* request fails
   * (i.e. the entire request to the cloud timed out). Individual failures
   * are handled when update_session_credit is called.
   *
   * @param failed_request - UpdateSessionRequest that couldn't be sent to the
   *                         cloud for whatever reason
   */
  void reset_updates(const UpdateSessionRequest &failed_request);

  /**
   * Collect any credit keys that are either exhausted, timed out, or terminated
   * and apply actions to the services if need be
   * @param updates_out (out) - vector to add usage updates to, if they exist
   */
  UpdateSessionRequest collect_updates();

  /**
   * Initialize credit received from the cloud in the system. This adds all the
   * charging keys to the credit manager for tracking
   * @param credit_response - message from cloud containing initial credits
   * @return true if init was successful
   */
  bool init_session_credit(
    const std::string &imsi,
    const std::string &session_id,
    const SessionState::Config &cfg,
    const CreateSessionResponse &response);

  /**
   * Update allowed credit from the cloud in the system
   * @param credit_response - message from cloud containing new credits
   */
  void update_session_credit(const UpdateSessionResponse &response);

  SessionTerminateRequest terminate_subscriber(const std::string &imsi);

  void complete_termination(
    const std::string &imsi,
    const std::string &session_id);

  uint64_t get_charging_credit(
    const std::string &imsi,
    uint32_t charging_key,
    Bucket bucket) const;

  uint64_t get_monitor_credit(
    const std::string &imsi,
    const std::string &mkey,
    Bucket bucket) const;

  /**
   * Initialize reauth for a subscriber service. If the subscriber cannot be
   * found, the method returns SESSION_NOT_FOUND
   */
  ChargingReAuthAnswer::Result init_charging_reauth(
    ChargingReAuthRequest request);

  void init_policy_reauth(
    PolicyReAuthRequest request,
    PolicyReAuthAnswer &answer_out);

 private:
  struct RulesToProcess {
    std::vector<std::string> static_rules;
    std::vector<PolicyRule> dynamic_rules;
  };
  std::shared_ptr<StaticRuleStore> rule_store_;
  std::shared_ptr<PipelinedClient> pipelined_client_;
  std::unordered_map<std::string, std::unique_ptr<SessionState>> session_map_;
  folly::EventBase *evb_;

 private:
  void new_report();

  /**
   * Process the create session response to get rules to activate/deactivate
   * instantly and schedule rules with activation/deactivation time info
   * to activate/deactivate later.
   */
  void process_create_session_response(
    const CreateSessionResponse &response,
    const std::unordered_set<uint32_t> &successful_credits,
    const std::string &imsi,
    const std::string &ip_addr,
    RulesToProcess *rules_to_activate,
    RulesToProcess *rules_to_deactivate);

  /**
   * Process the policy reauth request to get rules to activate/deactivate
   * instantly and schedule rules with activation/deactivation time info
   * to activate/deactivate later.
   * Policy reauth request also specifies a flat list of rule IDs to remove.
   * Rules need to be deactivated are categorized as either staic or dynamic
   * rule and put in the vector.
   */
  void process_policy_reauth_request(
    const PolicyReAuthRequest &request,
    const std::unique_ptr<SessionState> &session,
    RulesToProcess *rules_to_activate,
    RulesToProcess *rules_to_deactivate);

  void schedule_static_rule_activation(
    const std::string &imsi,
    const std::string &ip_addr,
    const StaticRuleInstall &static_rule);

  void schedule_dynamic_rule_activation(
    const std::string &imsi,
    const std::string &ip_addr,
    const DynamicRuleInstall &dynamic_rule);

  void schedule_static_rule_deactivation(
    const std::string &imsi,
    const StaticRuleInstall &static_rule);

  void schedule_dynamic_rule_deactivation(
    const std::string &imsi,
    const DynamicRuleInstall &dynamic_rule);
};

} // namespace magma

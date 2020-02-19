/**
 * Copyright (c) 2016-present, Facebook, Inc.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */
#pragma once

#include <unordered_map>
#include <unordered_set>
#include <vector>

#include <lte/protos/session_manager.grpc.pb.h>
#include <orc8r/protos/directoryd.pb.h>
#include <folly/io/async/EventBaseManager.h>

#include "AAAClient.h"
#include "SessionReporter.h"
#include "PipelinedClient.h"
#include "DirectorydClient.h"
#include "RuleStore.h"
#include "SessionState.h"
#include "SpgwServiceClient.h"

namespace magma {

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
    std::shared_ptr<SessionReporter> reporter,
    std::shared_ptr<StaticRuleStore> rule_store,
    std::shared_ptr<PipelinedClient> pipelined_client,
    std::shared_ptr<AsyncDirectorydClient> directoryd_client,
    std::shared_ptr<SpgwServiceClient> spgw_client,
    std::shared_ptr<aaa::AAAClient> aaa_client,
    long session_force_termination_timeout_ms,
    long quota_exhaustion_termination_on_init_ms);

  void attachEventBase(folly::EventBase *evb);

  // blocks
  void start();

  void stop();

  folly::EventBase& get_event_base();

  /**
   * Setup rules for all sessions in pipelined, used whenever pipelined
   * restarts and needs to recover state
   */
  bool setup(
    const std::uint64_t& epoch,
    std::function<void(Status status, SetupFlowsResult)> callback);

  /**
   * Insert a group of rule usage into the monitor and update credit manager
   * Assumes records are aggregates, as in the usages sent are cumulative and
   * not differences.
   *
   * @param records - a RuleRecordTable protobuf with a vector of RuleRecords
   */
  void aggregate_records(const RuleRecordTable& records);

  /**
   * reset_updates resets all of the charging keys being updated in
   * failed_request. This should only be called if the *entire* request fails
   * (i.e. the entire request to the cloud timed out). Individual failures
   * are handled when update_session_credits_and_rules is called.
   *
   * @param failed_request - UpdateSessionRequest that couldn't be sent to the
   *                         cloud for whatever reason
   */
  void reset_updates(const UpdateSessionRequest& failed_request);

  /**
   * Collect any credit keys that are either exhausted, timed out, or terminated
   * and apply actions to the services if need be
   * @param updates_out (out) - vector to add usage updates to, if they exist
   */
  UpdateSessionRequest collect_updates(std::vector<std::unique_ptr<ServiceAction>>& actions) const;

  /**
   * Perform any rule installs/removals that need to be executed given a
   * CreateSessionResponse.
   */
  bool handle_session_init_rule_updates(
    const std::string& imsi,
    SessionState& session_state,
    const CreateSessionResponse& response,
    std::unordered_set<uint32_t>& charging_credits_received);

  /**
   * Initialize credit received from the cloud in the system. This adds all the
   * charging keys to the credit manager for tracking
   * @param credit_response - message from cloud containing initial credits
   * @return true if init was successful
   */
  bool init_session_credit(
    const std::string& imsi,
    const std::string& session_id,
    const SessionState::Config& cfg,
    const CreateSessionResponse& response);

  /**
   * Process the update response from the reporter and update the
   * monitoring/charging credits and attached rules.
   * @param credit_response - message from cloud containing new credits
   */
  void update_session_credits_and_rules(const UpdateSessionResponse& response);

  /**
   * Starts the termination process for the session. When termination completes,
   * the call back function is executed.
   * @param imsi - imsi of the subscirber
   * @param on_termination_callback - callback function to be executed after
   * termination
   */
  void terminate_subscriber(
    const std::string& imsi,
    const std::string& apn,
    std::function<void(SessionTerminateRequest)> on_termination_callback);

  uint64_t get_charging_credit(
    const std::string& imsi,
    const CreditKey& charging_key,
    Bucket bucket) const;

  uint64_t get_monitor_credit(
    const std::string& imsi,
    const std::string& mkey,
    Bucket bucket) const;

  /**
   * Initialize reauth for a subscriber service. If the subscriber cannot be
   * found, the method returns SESSION_NOT_FOUND
   */
  ChargingReAuthAnswer::Result init_charging_reauth(
    ChargingReAuthRequest request);

  /**
   * Handles the equivalent of a RAR.
   * For the matching session ID, activate and/or deactivate the specified
   * rules.
   * Afterwards, a bearer is created.
   * If a session is CWF and out of monitoring quota, it will trigger a session
   * terminate
   *
   * NOTE: If an empty session ID is specified, apply changes to all matching
   * sessions with the specified IMSI.
   */
  void init_policy_reauth(
    PolicyReAuthRequest request,
    PolicyReAuthAnswer& answer_out);

  bool is_imsi_duplicate(const std::string& imsi) const;
  bool is_apn_duplicate(const std::string& imsi, const std::string& apn) const;

  std::string* duplicate_session_id(
    const std::string& imsi, const magma::SessionState::Config& config) const;

  /**
   * Execute actions on subscriber's service, eg. terminate, redirect data, or
   * just continue
   */
  void execute_actions(
    const std::vector<std::unique_ptr<ServiceAction>>& actions);

  static uint32_t REDIRECT_FLOW_PRIORITY;

 private:
  struct RulesToProcess {
    std::vector<std::string> static_rules;
    std::vector<PolicyRule> dynamic_rules;
  };
  std::shared_ptr<SessionReporter> reporter_;
  std::shared_ptr<StaticRuleStore> rule_store_;
  std::shared_ptr<PipelinedClient> pipelined_client_;
  std::shared_ptr<AsyncDirectorydClient> directoryd_client_;
  std::shared_ptr<SpgwServiceClient> spgw_client_;
  std::shared_ptr<aaa::AAAClient> aaa_client_;
  std::unordered_map<std::string,
                     std::vector<std::unique_ptr<SessionState>>> session_map_;
  folly::EventBase* evb_;
  long session_force_termination_timeout_ms_;
  // [CWF-ONLY] This configures how long we should wait before terminating a
  // session after it is created without any monitoring quota
  long quota_exhaustion_termination_on_init_ms_;

 private:
  /**
   * new_report notifies all sessions that a new usage report is going to be
   * aggregated.
   */
  void new_report();

  /**
   * finish_report notifies all sessions that the aggregation of the usage
   * report is finished. For sessions that are terminating, complete the
   * termination if the session is not included in the report.
   */
  void finish_report();

  /**
   * Process the create session response to get rules to activate/deactivate
   * instantly and schedule rules with activation/deactivation time info
   * to activate/deactivate later. No state change is made.
   */
  void process_create_session_response(
    const CreateSessionResponse& response,
    const std::unordered_set<uint32_t>& successful_credits,
    const std::string& imsi,
    const std::string& ip_addr,
    RulesToProcess& rules_to_activate,
    RulesToProcess& rules_to_deactivate);

  /**
   * Processes the charging component of UpdateSessionResponse.
   * Updates charging credits according to the response.
   */
  void update_charging_credits(
    const UpdateSessionResponse& response,
    std::unordered_set<std::string>& subscribers_to_terminate);

  /**
   * Processes the monitoring component of UpdateSessionResponse.
   * Updates moniroting credits according to the response and updates rules
   * that are installed for this session.
   * If a session is CWF and out of monitoring quota, it will trigger a session
   * terminate
   */
  void update_monitoring_credits_and_rules(
    const UpdateSessionResponse& response,
    std::unordered_set<std::string>& subscribers_to_terminate);

  /**
   * Process the list of rule names given and fill in rules_to_deactivate by
   * determining whether each one is dynamic or static. Modifies session state.
   * TODO separate out logic that modifies state vs logic that does not.
   */
  void process_rules_to_remove(
    const std::string& imsi,
    const std::unique_ptr<SessionState>& session,
    const google::protobuf::RepeatedPtrField<std::basic_string<char>>
      rules_to_remove,
    RulesToProcess& rules_to_deactivate);

  /**
   * Populate existing rules from a specific session;
   * used to delete flow rules for a PDN session,
   * distinct APNs are assumed to have mutually exclusive
   * rules.
   */
  void populate_rules_from_session_to_remove(
    const std::string& imsi,
    const std::unique_ptr<SessionState>& session,
    RulesToProcess& rules_to_deactivate);

  /**
   * Process protobuf StaticRuleInstalls and DynamicRuleInstalls to fill in
   * rules_to_activate and rules_to_deactivate. Modifies session state.
   * TODO separate out logic that modifies state vs logic that does not.
   */
  void process_rules_to_install(
    const std::string& imsi,
    const std::unique_ptr<SessionState>& session,
    const google::protobuf::RepeatedPtrField<magma::lte::StaticRuleInstall>
      static_rules_to_install,
    const google::protobuf::RepeatedPtrField<magma::lte::DynamicRuleInstall>
      dynamic_rules_to_install,
    RulesToProcess& rules_to_activate,
    RulesToProcess& rules_to_deactivate);

  /**
   * For the matching session ID, activate and/or deactivate the specified
   * rules.
   * Also create a bearer for the session.
   */
  void init_policy_reauth_for_session(
    const PolicyReAuthRequest& request,
    const std::unique_ptr<SessionState>& session,
    bool& activate_success,
    bool& deactivate_success);

  /**
   * Completes the session termination and executes the callback function
   * registered in terminate_subscriber.
   * complete_termination is called some time after terminate_subscriber
   * when the flows no longer appear in the usage report, meaning that they have
   * been deleted.
   * It is also called after a timeout to perform forced termination.
   * If the session cannot be found, either because it has already terminated,
   * or a new session for the subscriber has been created, then it will do
   * nothing.
   */
  void complete_termination(
    const std::string& imsi,
    const std::string& session_id);

  void schedule_static_rule_activation(
    const std::string& imsi,
    const std::string& ip_addr,
    const StaticRuleInstall& static_rule);

  void schedule_dynamic_rule_activation(
    const std::string& imsi,
    const std::string& ip_addr,
    const DynamicRuleInstall& dynamic_rule);

  void schedule_static_rule_deactivation(
    const std::string& imsi,
    const StaticRuleInstall& static_rule);

  void schedule_dynamic_rule_deactivation(
    const std::string& imsi,
    const DynamicRuleInstall &dynamic_rule);

  /**
   * Get the monitoring credits from PolicyReAuthRequest (RAR) message
   * and add the credits to UsageMonitoringCreditPool of the session
   */
  void receive_monitoring_credit_from_rar(
    const PolicyReAuthRequest& request,
    const std::unique_ptr<SessionState>& session);

  /**
   * Send bearer creation request through the PGW client if rules were
   * activated successfully in pipelined
   */
  void create_bearer(
    const bool activate_success,
    const std::unique_ptr<SessionState>& session,
    const PolicyReAuthRequest& request,
    const std::vector<PolicyRule>& dynamic_rules);

  /**
   * Check if REVALIDATION_TIMEOUT is one of the event triggers
   */
  bool revalidation_required(
    const google::protobuf::RepeatedField<int>& event_triggers);

  void schedule_revalidation(
    const google::protobuf::Timestamp& revalidation_time);

  void check_usage_for_reporting();

  /**
    * Deactivate rules for certain IMSI.
    * Notify AAA service if the session is a CWF session.
    */
  void terminate_service(
    const std::string& imsi,
    const std::vector<std::string>& rule_ids,
    const std::vector<PolicyRule>& dynamic_rules);


  /**
    * Deactivate rules for multiple IMSIs.
    * Notify AAA service if the session is a CWF session.
    */
  void terminate_multiple_services(
    const std::unordered_set<std::string>& imsis);

  /**
    * Install flow for redirection through pipelined
    */
  void install_redirect_flow(const std::unique_ptr<ServiceAction>& action);

  bool rules_to_process_is_not_empty(const RulesToProcess& rules_to_process);

  void report_subscriber_state_to_pipelined(
    const std::string& imsi,
    const std::string& ue_mac_addr,
    const SubscriberQuotaUpdate_Type state);

  /**
   * [CWF-ONLY]
   * If the session has active monitored rules attached to it, then propagate
   * to pipelined that the subscriber has valid quota.
   * Otherwise, mark the subscriber as out of quota to pipelined, and schedule
   * the session to be terminated in a configured amount of time.
   */
  void handle_session_init_subscriber_quota_state(
    const std::string& imsi,
    SessionState& session_state);
};

} // namespace magma

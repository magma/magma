/**
 * Copyright (c) 2016-present, Facebook, Inc.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */
#pragma once

#include <functional>
#include <utility>

#include <lte/protos/session_manager.grpc.pb.h>

#include "RuleStore.h"
#include "SessionReporter.h"
#include "StoredState.h"
#include "CreditKey.h"
#include "SessionCredit.h"
#include "Monitor.h"
#include "ChargingGrant.h"

namespace magma {
typedef std::unordered_map<CreditKey, std::unique_ptr<ChargingGrant>,
                 decltype(&ccHash), decltype(&ccEqual)> CreditMap;
typedef std::unordered_map<std::string, std::unique_ptr<Monitor>> MonitorMap;
static SessionStateUpdateCriteria UNUSED_UPDATE_CRITERIA;
/**
 * SessionState keeps track of a current UE session in the PCEF, recording
 * usage and allowance for all charging keys
 */
class SessionState {
 public:
  struct SessionInfo {
    std::string imsi;
    std::string ip_addr;
    std::vector<std::string> static_rules;
    std::vector<PolicyRule> dynamic_rules;
    std::vector<PolicyRule> gy_dynamic_rules;
  };
  struct TotalCreditUsage {
    uint64_t monitoring_tx;
    uint64_t monitoring_rx;
    uint64_t charging_tx;
    uint64_t charging_rx;
  };

 public:
  SessionState(
      const std::string& imsi, const std::string& session_id,
      const std::string& core_session_id, const SessionConfig& cfg,
      StaticRuleStore& rule_store, const magma::lte::TgppContext& tgpp_context);

  SessionState(
      const StoredSessionState& marshaled, StaticRuleStore& rule_store);

  static std::unique_ptr<SessionState> unmarshal(
      const StoredSessionState& marshaled, StaticRuleStore& rule_store);

  StoredSessionState marshal();

  /**
   * Updates rules to be scheduled, active, or removed, depending on the
   * specified time.
   *
   * NOTE: This function has undefined behavior if attempting to go backwards
   *       in time.
   */
  void sync_rules_to_time(
      std::time_t current_time, SessionStateUpdateCriteria& update_criteria);

  /**
   * notify_new_report_for_sessions sets the state of terminating session to
   * aggregating, to tell if
   * flows for the terminating session is in the latest report.
   * Should be called before add_rule_usage.
   */
  void new_report(SessionStateUpdateCriteria& update_criteria);

  /**
   * notify_finish_report_for_sessions updates the state of aggregating session
   * not included report
   * to specify its flows are deleted and termination can be completed.
   * Should be called after notify_new_report_for_sessions and add_rule_usage.
   */
  void finish_report(SessionStateUpdateCriteria& update_criteria);

  /**
   * add_rule_usage adds used TX/RX bytes to a particular rule
   */
  void add_rule_usage( const std::string& rule_id, uint64_t used_tx,
    uint64_t used_rx, SessionStateUpdateCriteria& update_criteria);

  /**
   * get_updates collects updates and adds them to a UpdateSessionRequest
   * for reporting.
   * Only updates request number
   * @param update_request (out) - request to add new updates to
   * @param actions (out) - actions to take on services
   */
  void get_updates(
      UpdateSessionRequest& update_request_out,
      std::vector<std::unique_ptr<ServiceAction>>* actions_out,
      SessionStateUpdateCriteria& update_criteria);

  /**
   * start_termination starts the termination process for the session.
   * The session state transitions from SESSION_ACTIVE to
   * SESSION_TERMINATING_FLOW_ACTIVE.
   * When termination completes, the call back function is executed.
   *
   * @param on_termination_callback - call back function to be executed after
   * termination
   */
  void start_termination(SessionStateUpdateCriteria& update_criteria);

  /**
   * mark_as_awaiting_termination transitions the session state from
   * SESSION_ACTIVE to SESSION_TERMINATION_SCHEDULED
   */
  void mark_as_awaiting_termination(
      SessionStateUpdateCriteria& update_criteria);

  bool is_terminating();

  /**
   * can_complete_termination returns whether the termination for the session
   * can be completed.
   * For this to be true, start_termination needs to be called for the session,
   * and the flows for the session needs to be deleted.
   */
  bool can_complete_termination() const;

  /**
   * complete_termination collects final usages for all credits into a
   * SessionTerminateRequest and calls the on termination callback with the
   * request.
   * Note that complete_termination will forcefully complete the termination
   * no matter the current state of the session. To properly complete the
   * termination, this function should only be called when
   * can_complete_termination returns true.
   */
  void complete_termination(
      SessionReporter& reporter, SessionStateUpdateCriteria& update_criteria);

  bool reset_reporting_charging_credit(const CreditKey &key,
                              SessionStateUpdateCriteria &update_criteria);

  bool receive_charging_credit(const CreditUpdateResponse &update,
                               SessionStateUpdateCriteria &update_criteria);

  uint64_t get_charging_credit(const CreditKey &key, Bucket bucket) const;

  ReAuthResult reauth_key(const CreditKey &charging_key,
                               SessionStateUpdateCriteria &update_criteria);

  ReAuthResult reauth_all(SessionStateUpdateCriteria &update_criteria);

  void merge_charging_credit_update(
    const CreditKey &key, SessionCreditUpdateCriteria &credit_update);

  void set_charging_credit(
    const CreditKey &key, SessionCredit credit, SessionStateUpdateCriteria &uc);

  /**
   * get_total_credit_usage returns the tx and rx of the session,
   * accounting for all unique keys (charging and monitoring) used by all
   * rules (static and dynamic)
   * Should be called after complete_termination.
   */
  TotalCreditUsage get_total_credit_usage();

  std::string get_session_id() const;

  std::string get_core_session_id() const { return core_session_id_; };

  SubscriberQuotaUpdate_Type get_subscriber_quota_state() const;

  bool is_radius_cwf_session() const;

  bool is_same_config(const SessionConfig& new_config) const;

  void get_session_info(SessionState::SessionInfo& info);

  void set_tgpp_context(
      const magma::lte::TgppContext& tgpp_context,
      SessionStateUpdateCriteria& update_criteria);

  void set_config(const SessionConfig& config);

  SessionConfig get_config() const;

  void set_subscriber_quota_state(
      const magma::lte::SubscriberQuotaUpdate_Type state,
      SessionStateUpdateCriteria& update_criteria);

  bool active_monitored_rules_exist();

  uint32_t get_request_number();

  void increment_request_number(uint32_t incr);

  // Methods related to the session's static and dynamic rules
  bool is_dynamic_rule_installed(const std::string& rule_id);

  bool is_gy_dynamic_rule_installed(const std::string& rule_id);

  bool is_static_rule_installed(const std::string& rule_id);

  bool is_dynamic_rule_scheduled(const std::string& rule_id);

  bool is_static_rule_scheduled(const std::string& rule_id);

  /**
   * Add a dynamic rule to the session which is currently active.
   */
  void insert_dynamic_rule(
      const PolicyRule& rule, RuleLifetime& lifetime,
      SessionStateUpdateCriteria& update_criteria);

  /**
   * Add a static rule to the session which is currently active.
   */
  void activate_static_rule(
      const std::string& rule_id, RuleLifetime& lifetime,
      SessionStateUpdateCriteria& update_criteria);

  void insert_gy_dynamic_rule(
      const PolicyRule& rule, RuleLifetime& lifetime,
      SessionStateUpdateCriteria& update_criteria);

  /**
   * Remove a currently active dynamic rule to mark it as deactivated.
   *
   * @param rule_id ID of the rule to be removed.
   * @param rule_out Will point to the removed rule.
   * @param update_criteria Tracks updates to the session. To be passed back to
   *                        the SessionStore to resolve issues of concurrent
   *                        updates to a session.
   * @return True if successfully removed.
   */
  bool remove_dynamic_rule(
      const std::string& rule_id, PolicyRule* rule_out,
      SessionStateUpdateCriteria& update_criteria);

  bool remove_scheduled_dynamic_rule(
      const std::string& rule_id, PolicyRule* rule_out,
      SessionStateUpdateCriteria& update_criteria);

  /**
   * Remove a currently active static rule to mark it as deactivated.
   *
   * @param rule_id ID of the rule to be removed.
   * @param update_criteria Tracks updates to the session. To be passed back to
   *                        the SessionStore to resolve issues of concurrent
   *                        updates to a session.
   * @return True if successfully removed.
   */
  bool deactivate_static_rule(
      const std::string& rule_id, SessionStateUpdateCriteria& update_criteria);

  bool remove_gy_dynamic_rule(
      const std::string& rule_id, PolicyRule *rule_out,
      SessionStateUpdateCriteria& update_criteria);

  bool deactivate_scheduled_static_rule(
      const std::string& rule_id, SessionStateUpdateCriteria& update_criteria);

  std::vector<std::string>& get_static_rules();

  std::set<std::string>& get_scheduled_static_rules();

  DynamicRuleStore& get_dynamic_rules();

  DynamicRuleStore& get_scheduled_dynamic_rules();

  /**
   * Schedule a dynamic rule for activation in the future.
   */
  void schedule_dynamic_rule(
      const PolicyRule& rule, RuleLifetime& lifetime,
      SessionStateUpdateCriteria& update_criteria);

  /**
   * Schedule a static rule for activation in the future.
   */
  void schedule_static_rule(
      const std::string& rule_id, RuleLifetime& lifetime,
      SessionStateUpdateCriteria& update_criteria);

  /**
   * Mark a scheduled dynamic rule as activated.
   */
  void install_scheduled_dynamic_rule(
      const std::string& rule_id, SessionStateUpdateCriteria& update_criteria);

  /**
   * Mark a scheduled static rule as activated.
   */
  void install_scheduled_static_rule(
      const std::string& rule_id, SessionStateUpdateCriteria& update_criteria);

  RuleLifetime& get_rule_lifetime(const std::string& rule_id);

  DynamicRuleStore& get_gy_dynamic_rules();

  uint32_t total_monitored_rules_count();

  bool is_active();

  uint32_t get_credit_key_count();

  void set_fsm_state(
    SessionFsmState new_state, SessionStateUpdateCriteria& uc);

  StaticRuleInstall get_static_rule_install(
    const std::string& rule_id, const RuleLifetime& lifetime);

  DynamicRuleInstall get_dynamic_rule_install(
    const std::string& rule_id, const RuleLifetime& lifetime);

  SessionFsmState get_state();

  // Event Triggers
  void add_new_event_trigger(
  magma::lte::EventTrigger trigger,
  SessionStateUpdateCriteria& update_criteria);

  void mark_event_trigger_as_triggered(
    magma::lte::EventTrigger trigger,
    SessionStateUpdateCriteria& update_criteria);

  void set_event_trigger(
    magma::lte::EventTrigger trigger, const EventTriggerState value,
    SessionStateUpdateCriteria& update_criteria);

  void remove_event_trigger(magma::lte::EventTrigger trigger,
    SessionStateUpdateCriteria& update_criteria);

  void set_revalidation_time(const google::protobuf::Timestamp& time,
                             SessionStateUpdateCriteria& update_criteria);

  google::protobuf::Timestamp get_revalidation_time() {return revalidation_time_;}

  EventTriggerStatus get_event_triggers() {return pending_event_triggers_;}

  bool is_credit_state_redirected(const CreditKey &charging_key) const;

  // Monitors
  bool receive_monitor(const UsageMonitoringUpdateResponse &update,
                       SessionStateUpdateCriteria &uc);

  void merge_monitor_updates(
    const std::string &key, SessionCreditUpdateCriteria &update);

  uint64_t get_monitor(const std::string &key, Bucket bucket) const;

  bool add_to_monitor(const std::string &key, uint64_t used_tx,
                      uint64_t used_rx, SessionStateUpdateCriteria &uc);

  void set_monitor(
    const std::string &key, std::unique_ptr<Monitor> monitor,
    SessionStateUpdateCriteria &uc);

  bool reset_reporting_monitor(
    const std::string &key, SessionStateUpdateCriteria &uc);

  std::unique_ptr<std::string> get_session_level_key() const;
 private:
  std::string imsi_;
  std::string session_id_;
  std::string core_session_id_;
  uint32_t request_number_;
  SessionFsmState curr_state_;
  SessionConfig config_;
  // Used to keep track of whether the subscriber has valid quota.
  // (only used for CWF at the moment)
  magma::lte::SubscriberQuotaUpdate_Type subscriber_quota_state_;
  magma::lte::TgppContext tgpp_context_;
  std::function<void(SessionTerminateRequest)> on_termination_callback_;

  // All static rules synced from policy DB
  StaticRuleStore& static_rules_;
  // Static rules that are currently installed for the session
  std::vector<std::string> active_static_rules_;
  // Dynamic GX rules that are currently installed for the session
  DynamicRuleStore dynamic_rules_;
  // Dynamic GY rules that are currently installed for the session
  DynamicRuleStore gy_dynamic_rules_;

  // Static rules that are scheduled for installation for the session
  std::set<std::string> scheduled_static_rules_;
  // Dynamic rules that are scheduled for installation for the session
  DynamicRuleStore scheduled_dynamic_rules_;
  // Activation & deactivation times for each rule that is either currently
  // installed, or scheduled for installation for this session
  std::unordered_map<std::string, RuleLifetime> rule_lifetimes_;

  // map of Gx event_triggers that are pending and its status (bool)
  // If the value is true, that means an update request for that event trigger
  // should be sent
  EventTriggerStatus pending_event_triggers_;
  // todo for stateless we will have to store a bit more information so we can
  // reschedule triggers
  google::protobuf::Timestamp revalidation_time_;
  /*
  * ChargingCreditPool manages a pool of credits for OCS-based charging. It is
  * keyed by rating groups & service Identity (uint32, [uint32]) and receives
  * CreditUpdateResponses to update
  * credit
  */
  CreditMap credit_map_;
  MonitorMap monitor_map_;
  std::unique_ptr<std::string> session_level_key_;

 private:
  /**
   * For this session, add the CreditUsageUpdate to the UpdateSessionRequest.
   * Also
   *
   * @param update_request_out Modified with added CreditUsageUpdate
   * @param actions_out Modified with additional actions to take on session
   */
  void get_charging_updates(
      UpdateSessionRequest& update_request_out,
      std::vector<std::unique_ptr<ServiceAction>>* actions_out,
      SessionStateUpdateCriteria& uc);

  bool init_charging_credit(
    const CreditUpdateResponse &update, SessionStateUpdateCriteria &uc);

  /**
   * For this session, add the UsageMonitoringUpdateRequest to the
   * UpdateSessionRequest.
   *
   * @param update_request_out Modified with added UsdageMonitoringUpdateRequest
   * @param actions_out Modified with additional actions to take on session.
   */
  void get_monitor_updates(
      UpdateSessionRequest& update_request_out,
      std::vector<std::unique_ptr<ServiceAction>>* actions_out,
      SessionStateUpdateCriteria& update_criteria);

  void add_common_fields_to_usage_monitor_update(UsageMonitoringUpdateRequest* req);

  SessionTerminateRequest make_termination_request(
    SessionStateUpdateCriteria& update_criteria);

  /**
   * Returns true if the specified rule should be active at that time
   */
  bool should_rule_be_active(const std::string& rule_id, std::time_t time);

  /**
   * Returns true if the specified rule should be deactivated by that time
   */
  bool should_rule_be_deactivated(const std::string& rule_id, std::time_t time);

  SessionCreditUpdateCriteria* get_credit_uc(
    const CreditKey &key, SessionStateUpdateCriteria &uc);

  CreditUsageUpdate make_credit_usage_update_req(CreditUsage& usage) const;

  bool init_new_monitor(
    const UsageMonitoringUpdateResponse &update,
    SessionStateUpdateCriteria &update_criteria);

  void update_session_level_key(
    const UsageMonitoringUpdateResponse &update,
    SessionStateUpdateCriteria &update_criteria);

  SessionCreditUpdateCriteria* get_monitor_uc(
    const std::string &key, SessionStateUpdateCriteria &uc);

  void fill_protos_tgpp_context(magma::lte::TgppContext* tgpp_context) const;
};

}  // namespace magma

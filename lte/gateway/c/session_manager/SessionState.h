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
#include "StoredState.h"
#include "CreditPool.h"

namespace magma {

/**
 * SessionState keeps track of a current UE session in the PCEF, recording
 * usage and allowance for all charging keys
 */
class SessionState {
 public:
  static SessionStateUpdateCriteria UNUSED_UPDATE_CRITERIA;

  struct QoSInfo {
    bool enabled;
    uint32_t qci;
  };

  struct Config {
    std::string ue_ipv4;
    std::string spgw_ipv4;
    std::string msisdn;
    std::string apn;
    std::string imei;
    std::string plmn_id;
    std::string imsi_plmn_id;
    std::string user_location;
    RATType rat_type;
    std::string mac_addr; // MAC Address for WLAN
    std::string hardware_addr; // MAC Address for WLAN (binary)
    std::string radius_session_id;
    uint32_t bearer_id;
    QoSInfo qos_info;
  };
  struct SessionInfo {
    std::string imsi;
    std::string ip_addr;
    std::vector<std::string> static_rules;
    std::vector<PolicyRule> dynamic_rules;
  };
  struct TotalCreditUsage {
    uint64_t monitoring_tx;
    uint64_t monitoring_rx;
    uint64_t charging_tx;
    uint64_t charging_rx;
  };


 public:
  SessionState(
    const std::string& imsi,
    const std::string& session_id,
    const std::string& core_session_id,
    const SessionState::Config& cfg,
    StaticRuleStore& rule_store,
    const magma::lte::TgppContext& tgpp_context);

  SessionState(
    const StoredSessionState &marshaled,
    StaticRuleStore &rule_store);

  static std::unique_ptr<SessionState> unmarshal(
    const StoredSessionState &marshaled,
    StaticRuleStore &rule_store);

  StoredSessionState marshal();

  /**
   * notify_new_report_for_sessions sets the state of terminating session to aggregating, to tell if
   * flows for the terminating session is in the latest report.
   * Should be called before add_used_credit.
   */
  void new_report();

  /**
   * notify_finish_report_for_sessions updates the state of aggregating session not included report
   * to specify its flows are deleted and termination can be completed.
   * Should be called after notify_new_report_for_sessions and add_used_credit.
   */
  void finish_report();

  /**
   * add_used_credit adds used TX/RX bytes to a particular charging key
   */
  void add_used_credit(
    const std::string& rule_id,
    uint64_t used_tx,
    uint64_t used_rx,
    SessionStateUpdateCriteria& update_criteria = UNUSED_UPDATE_CRITERIA);

  /**
   * get_updates collects updates and adds them to a UpdateSessionRequest
   * for reporting.
   * Only updates request number
   * @param update_request (out) - request to add new updates to
   * @param actions (out) - actions to take on services
   * @param force_update force updates if revalidation timer expires
   */
  void get_updates(
    UpdateSessionRequest& update_request_out,
    std::vector<std::unique_ptr<ServiceAction>>* actions_out,
    SessionStateUpdateCriteria& update_criteria = UNUSED_UPDATE_CRITERIA,
    const bool force_update = false);

  /**
   * start_termination starts the termination process for the session.
   * The session state transitions from SESSION_ACTIVE to
   * SESSION_TERMINATING_FLOW_ACTIVE.
   * When termination completes, the call back function is executed.
   * @param on_termination_callback - call back function to be executed after
   * termination
   */
  void start_termination(
    std::function<void(SessionTerminateRequest)> on_termination_callback,
    SessionStateUpdateCriteria& update_criteria = UNUSED_UPDATE_CRITERIA);

  /**
   * mark_as_awaiting_termination transitions the session state from
   * SESSION_ACTIVE to SESSION_TERMINATION_SCHEDULED
   */
  void mark_as_awaiting_termination(
    SessionStateUpdateCriteria& update_criteria = UNUSED_UPDATE_CRITERIA);

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
    SessionStateUpdateCriteria& update_criteria = UNUSED_UPDATE_CRITERIA);

  ChargingCreditPool& get_charging_pool();

  UsageMonitoringCreditPool& get_monitor_pool();

  /**
   * get_total_credit_usage returns the tx and rx of the session,
   * accounting for all unique keys (charging and monitoring) used by all
   * rules (static and dynamic)
   * Should be called after complete_termination.
   */
  TotalCreditUsage get_total_credit_usage();

  std::string get_session_id() const;

  std::string get_subscriber_ip_addr() const;

  std::string get_mac_addr() const;

  std::string get_msisdn() const;

  std::string get_hardware_addr() const { return config_.hardware_addr; }

  std::string get_radius_session_id() const;

  std::string get_apn() const;

  std::string get_core_session_id() const { return core_session_id_; };

  uint32_t get_bearer_id() const;

  uint32_t get_qci() const;

  SubscriberQuotaUpdate_Type get_subscriber_quota_state() const;

  bool is_radius_cwf_session() const;

  bool is_same_config(const Config& new_config) const;

  void get_session_info(SessionState::SessionInfo& info);

  bool qos_enabled() const;

  void set_tgpp_context(
    const magma::lte::TgppContext& tgpp_context,
    SessionStateUpdateCriteria& update_criteria = UNUSED_UPDATE_CRITERIA);

  void fill_protos_tgpp_context(magma::lte::TgppContext* tgpp_context) const;

  void set_subscriber_quota_state(
    const magma::lte::SubscriberQuotaUpdate_Type state,
    SessionStateUpdateCriteria& update_criteria = UNUSED_UPDATE_CRITERIA);

  bool active_monitored_rules_exist();

  uint32_t get_request_number();

  void increment_request_number(uint32_t incr);

  // Methods related to the session's static and dynamic rules
  bool get_charging_key_for_rule_id(
    const std::string& rule_id,
    CreditKey* charging_key);

  bool get_monitoring_key_for_rule_id(
    const std::string& rule_id,
    std::string* monitoring_key);

  bool is_dynamic_rule_installed(const std::string& rule_id);

  bool is_static_rule_installed(const std::string& rule_id);

  void insert_dynamic_rule(
    const PolicyRule& rule,
    SessionStateUpdateCriteria& update_criteria = UNUSED_UPDATE_CRITERIA);

  void activate_static_rule(
    const std::string& rule_id,
    SessionStateUpdateCriteria& update_criteria = UNUSED_UPDATE_CRITERIA);

  bool remove_dynamic_rule(
    const std::string& rule_id,
    PolicyRule *rule_out,
    SessionStateUpdateCriteria& update_criteria = UNUSED_UPDATE_CRITERIA);

  bool deactivate_static_rule(
    const std::string& rule_id,
    SessionStateUpdateCriteria& update_criteria = UNUSED_UPDATE_CRITERIA);

  DynamicRuleStore& get_dynamic_rules();

  uint32_t total_monitored_rules_count();

  uint32_t get_credit_key_count();

 private:
  /**
   * State transitions of a session:
   * SESSION_ACTIVE  ---------
   *       |                  \
   *       |                   \
   *       |                    \
   *       |                     \
   *       | (start_termination)  SESSION_TERMINATION_SCHEDULED
   *       |                      /
   *       |                     /
   *       |                    /
   *       V                   V
   * SESSION_TERMINATING_FLOW_ACTIVE <----------
   *       |                                   |
   *       | (notify_new_report_for_sessions)  | (add_used_credit)
   *       V                                   |
   * SESSION_TERMINATING_AGGREGATING_STATS -----
   *       |
   *       | (notify_finish_report_for_sessions)
   *       V
   * SESSION_TERMINATING_FLOW_DELETED
   *       |
   *       | (complete_termination)
   *       V
   * SESSION_TERMINATED
   */
  enum State {
    SESSION_ACTIVE = 0,
    SESSION_TERMINATING_FLOW_ACTIVE = 1,
    SESSION_TERMINATING_AGGREGATING_STATS = 2,
    SESSION_TERMINATING_FLOW_DELETED = 3,
    SESSION_TERMINATED = 4,
    // TODO All sessions in this state should be terminated on sessiond restart
    SESSION_TERMINATION_SCHEDULED = 5
  };

  std::string imsi_;
  std::string session_id_;
  std::string core_session_id_;
  uint32_t request_number_;
  ChargingCreditPool charging_pool_;
  UsageMonitoringCreditPool monitor_pool_;
  SessionState::State curr_state_;
  SessionState::Config config_;
  // Used to keep track of whether the subscriber has valid quota.
  // (only used for CWF at the moment)
  magma::lte::SubscriberQuotaUpdate_Type subscriber_quota_state_;
  magma::lte::TgppContext tgpp_context_;
  std::function<void(SessionTerminateRequest)> on_termination_callback_;

  // All static rules synced from policy DB
  StaticRuleStore& static_rules_;
  // Static rules that are currently installed for the session
  std::vector<std::string> active_static_rules_;
  // Dynamic rules that are currently installed for the session
  DynamicRuleStore dynamic_rules_;

 private:
  /**
   * For this session, add the CreditUsageUpdate to the UpdateSessionRequest.
   * Also
   *
   * @param update_request_out Modified with added CreditUsageUpdate
   * @param actions_out Modified with additional actions to take on session
   * @param force_update force updates if revalidation timer expires
   */
  void get_updates_from_charging_pool(
    UpdateSessionRequest& update_request_out,
    std::vector<std::unique_ptr<ServiceAction>>* actions_out,
    SessionStateUpdateCriteria& update_criteria = UNUSED_UPDATE_CRITERIA,
    const bool force_update = false);

  /**
   * For this session, add the UsageMonitoringUpdateRequest to the
   * UpdateSessionRequest.
   *
   * @param update_request_out Modified with added UsdageMonitoringUpdateRequest
   * @param actions_out Modified with additional actions to take on session.
   * @param force_update force updates if revalidation timer expires
   */
  void get_updates_from_monitor_pool(
    UpdateSessionRequest& update_request_out,
    std::vector<std::unique_ptr<ServiceAction>>* actions_out,
    SessionStateUpdateCriteria& update_criteria = UNUSED_UPDATE_CRITERIA,
    const bool force_update = false);
};

} // namespace magma

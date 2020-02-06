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

#include <lte/protos/session_manager.grpc.pb.h>

#include "RuleStore.h"
#include "SessionRules.h"
#include "CreditPool.h"

namespace magma {

/**
 * SessionState keeps track of a current UE session in the PCEF, recording
 * usage and allowance for all charging keys
 */
class SessionState {
 public:
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

 public:
  SessionState(
    const std::string& imsi,
    const std::string& session_id,
    const std::string& core_session_id,
    const SessionState::Config& cfg,
    StaticRuleStore& rule_store,
    const magma::lte::TgppContext& tgpp_context);

  /**
   * new_report sets the state of terminating session to aggregating, to tell if
   * flows for the terminating session is in the latest report.
   * Should be called before add_used_credit.
   */
  void new_report();

  /**
   * finish_report updates the state of aggregating session not included report
   * to specify its flows are deleted and termination can be completed.
   * Should be called after new_report and add_used_credit.
   */
  void finish_report();

  /**
   * add_used_credit adds used TX/RX bytes to a particular charging key
   */
  void add_used_credit(
    const std::string& rule_id,
    uint64_t used_tx,
    uint64_t used_rx);

  /**
   * get_updates collects updates and adds them to a UpdateSessionRequest
   * for reporting.
   * @param update_request (out) - request to add new updates to
   * @param actions (out) - actions to take on services
   */
  void get_updates(
    UpdateSessionRequest& update_request_out,
    std::vector<std::unique_ptr<ServiceAction>>* actions_out);

  /**
   * start_termination starts the termination process for the session.
   * The session state transitions from SESSION_ACTIVE to
   * SESSION_TERMINATING_FLOW_ACTIVE.
   * When termination completes, the call back function is executed.
   * @param on_termination_callback - call back function to be executed after
   * termination
   */
  void start_termination(
    std::function<void(SessionTerminateRequest)> on_termination_callback);

  /**
   * can_complete_termination returns whether the termination for the session
   * can be completed.
   * For this to be true, start_termination needs to be called for the session,
   * and the flows for the session needs to be deleted.
   */
  bool can_complete_termination();

  /**
   * complete_termination collects final usages for all credits into a
   * SessionTerminateRequest and calls the on termination callback with the
   * request.
   * Note that complete_termination will forcefully complete the termination
   * no matter the current state of the session. To properly complete the
   * termination, this function should only be called when
   * can_complete_termination returns true.
   */
  void complete_termination();

  void insert_dynamic_rule(const PolicyRule& dynamic_rule);

  void activate_static_rule(const std::string& rule_id);

  bool remove_dynamic_rule(const std::string& rule_id, PolicyRule* rule_out);

  bool deactivate_static_rule(const std::string& rule_id);

  ChargingCreditPool& get_charging_pool();

  UsageMonitoringCreditPool& get_monitor_pool();

  std::string get_session_id();

  std::string get_subscriber_ip_addr();

  std::string get_mac_addr();

  std::string get_hardware_addr() { return config_.hardware_addr; }

  std::string get_radius_session_id();

  std::string get_apn();

  std::string get_core_session_id() const { return core_session_id_; };

  uint32_t get_bearer_id();

  uint32_t get_qci();

  bool is_radius_cwf_session();

  bool is_same_config(const Config& new_config);

  void get_session_info(SessionState::SessionInfo& info);

  bool qos_enabled();

  void set_tgpp_context(const magma::lte::TgppContext& tgpp_context);

  void fill_protos_tgpp_context(magma::lte::TgppContext* tgpp_context);

 private:
  /**
   * State transitions of a session:
   * SESSION_ACTIVE
   *       |
   *       | (start_termination)
   *       V
   * SESSION_TERMINATING_FLOW_ACTIVE <----------
   *       |                                   |
   *       | (new_report)                      | (add_used_credit)
   *       V                                   |
   * SESSION_TERMINATING_AGGREGATING_STATS -----
   *       |
   *       | (finish_report)
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
    SESSION_TERMINATED = 4
  };

  std::string imsi_;
  std::string session_id_;
  std::string core_session_id_;
  uint32_t request_number_;
  ChargingCreditPool charging_pool_;
  UsageMonitoringCreditPool monitor_pool_;
  SessionRules session_rules_;
  SessionState::State curr_state_;
  SessionState::Config config_;
  magma::lte::TgppContext tgpp_context_;
  std::function<void(SessionTerminateRequest)> on_termination_callback_;

 private:
  /**
   * For this session, add the CreditUsageUpdate to the UpdateSessionRequest.
   * Also
   *
   * @param update_request_out Modified with added CreditUsageUpdate
   * @param actions_out Modified with additional actions to take on session
   */
  void get_updates_from_charging_pool(
    UpdateSessionRequest& update_request_out,
    std::vector<std::unique_ptr<ServiceAction>>* actions_out);

  /**
   * For this session, add the UsageMonitoringUpdateRequest to the
   * UpdateSessionRequest.
   *
   * @param update_request_out Modified with added UsdageMonitoringUpdateRequest
   * @param actions_out Modified with additional actions to take on session.
   */
  void get_updates_from_monitor_pool(
    UpdateSessionRequest& update_request_out,
    std::vector<std::unique_ptr<ServiceAction>>* actions_out);
};

} // namespace magma

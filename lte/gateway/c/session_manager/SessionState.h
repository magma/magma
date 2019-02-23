/**
 * Copyright (c) 2016-present, Facebook, Inc.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */
#pragma once

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
  struct Config {
    std::string ue_ipv4;
    std::string spgw_ipv4;
    std::string msisdn;
    std::string apn;
    std::string imei;
    std::string plmn_id;
    std::string imsi_plmn_id;
    std::string user_location;
  };

public:
  SessionState(
    const std::string& imsi,
    const std::string& session_id,
    const SessionState::Config& cfg,
    StaticRuleStore& rule_store);

  /**
   * new_report unmarks all credits before an update, to tell if any credits
   * were not in the latest report
   */
  void new_report();

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
    UpdateSessionRequest* update_request_out,
    std::vector<std::unique_ptr<ServiceAction>>* actions_out);

  /**
   * Terminate collects final usages for all credits and returns them in a
   * SessionTerminateRequest
   */
  SessionTerminateRequest terminate();

  void insert_dynamic_rule(const PolicyRule& dynamic_rule);

  bool remove_dynamic_rule(const std::string& rule_id, PolicyRule* rule_out);

  ChargingCreditPool& get_charging_pool();

  UsageMonitoringCreditPool& get_monitor_pool();

  std::string get_session_id();

  std::string get_subscriber_ip_addr();


private:
  enum State {
    SESSION_ACTIVE = 0,
    SESSION_TERMINATING = 1,
    SESSION_TERMINATING_REPORTING = 2,
  };

  std::string imsi_;
  std::string session_id_;
  uint32_t request_number_;
  ChargingCreditPool charging_pool_;
  UsageMonitoringCreditPool monitor_pool_;
  SessionRules session_rules_;
  SessionState::State curr_state_;
  SessionState::Config config_;
private:
  void get_updates_from_charging_pool(
    UpdateSessionRequest* update_request_out,
    std::vector<std::unique_ptr<ServiceAction>>* actions_out);

  void get_updates_from_monitor_pool(
    UpdateSessionRequest* update_request_out,
    std::vector<std::unique_ptr<ServiceAction>>* actions_out);
};

}

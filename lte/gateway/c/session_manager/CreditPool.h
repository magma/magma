/**
 * Copyright (c) 2016-present, Facebook, Inc.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#pragma once

#include "CreditKey.h"
#include "RuleStore.h"
#include "SessionCredit.h"
#include "StoredState.h"
#include <unordered_map>

namespace magma {
/**
 * UsageMonitoringCreditPool manages a pool of credits for PCRF-based usage
 * monitoring. It is keyed by monitoring keys (string) and receives
 * UsageMonitoringUpdateResponse to update credit
 */
class UsageMonitoringCreditPool {
public:
  struct Monitor {
    SessionCredit credit;
    MonitoringLevel level;

    Monitor() : credit(CreditType::MONITORING) {}
  };

  static StoredMonitor
  marshal_monitor(std::unique_ptr<UsageMonitoringCreditPool::Monitor> &monitor);

  static std::unique_ptr<Monitor>
  unmarshal_monitor(const StoredMonitor &marshaled);

  static std::unique_ptr<UsageMonitoringCreditPool>
  unmarshal(const StoredUsageMonitoringCreditPool &marshaled);

  StoredUsageMonitoringCreditPool marshal();

  UsageMonitoringCreditPool(const std::string &imsi);

  bool add_used_credit(const std::string &key, uint64_t used_tx,
                       uint64_t used_rx,
                       SessionStateUpdateCriteria &update_criteria);

  bool
  reset_reporting_credit(const std::string &key,
                         SessionStateUpdateCriteria &update_criteria);

  void get_updates(std::string imsi, std::string ip_addr,
                   StaticRuleStore &static_rules,
                   DynamicRuleStore *dynamic_rules,
                   std::vector<UsageMonitorUpdate> *updates_out,
                   std::vector<std::unique_ptr<ServiceAction>> *actions_out,
                   SessionStateUpdateCriteria &update_criteria);

  bool
  get_termination_updates(SessionTerminateRequest *termination_out,
                          SessionStateUpdateCriteria &update_criteria);

  bool receive_credit(const UsageMonitoringUpdateResponse &update,
                      SessionStateUpdateCriteria &update_criteria);

  uint64_t get_credit(const std::string &key, Bucket bucket) const;

  void add_monitor(const std::string &key,
                   std::unique_ptr<UsageMonitoringCreditPool::Monitor> monitor,
                   SessionStateUpdateCriteria &update_criteria);

  SessionCreditUpdateCriteria *
  get_credit_update(const std::string &key,
                    SessionStateUpdateCriteria &update_criteria);

  void merge_credit_update(const std::string &key,
                           SessionCreditUpdateCriteria &credit_update);

  uint32_t get_credit_key_count() const;

  std::unique_ptr<std::string> get_session_level_key();

private:
  std::unordered_map<std::string, std::unique_ptr<Monitor>> monitor_map_;
  std::string imsi_;
  std::unique_ptr<std::string> session_level_key_;

private:
  void update_session_level_key(const UsageMonitoringUpdateResponse &update,
                                SessionStateUpdateCriteria &update_criteria);
  bool init_new_credit(const UsageMonitoringUpdateResponse &update,
                       SessionStateUpdateCriteria &update_criteria);
  void populate_output_actions(
      std::string imsi, std::string ip_addr, std::string key,
      StaticRuleStore &static_rules, DynamicRuleStore *dynamic_rules,
      std::unique_ptr<ServiceAction> &action,
      std::vector<std::unique_ptr<ServiceAction>> *actions_out) const;
};

} // namespace magma

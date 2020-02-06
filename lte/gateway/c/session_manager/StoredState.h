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

struct StoredQoSInfo {
  bool enabled;
  uint32_t qci;
};

struct StoredSessionConfig {
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
  StoredQoSInfo qos_info;
};

// Session Credit

struct StoredFinalActionInfo {
  ChargingCredit_FinalAction final_action;
  RedirectServer redirect_server;
};

struct StoredSessionCredit {
  bool reporting;
  bool is_final;
  bool unlimited_quota;
  StoredFinalActionInfo final_action_info;
  ReAuthState reauth_state;
  ServiceState service_state;
  std::time_t  expiry_time;
  std::unordered_map<Bucket, uint64_t> buckets;
  uint64_t usage_reporting_limit;
};

struct StoredMonitor {
  StoredSessionCredit credit;
  MonitoringLevel level;
};

struct StoredChargingCreditPool {
  std::string imsi;
  std::unordered_map<
    CreditKey, StoredSessionCredit,
    decltype(&ccHash), decltype(&ccEqual)> credit_map;
};

struct StoredUsageMonitoringCreditPool {
  std::string imsi;
  std::string session_level_key; // "" maps to nullptr
  std::unordered_map<std::string, StoredMonitor> monitor_map;
};

// Installed session rules
struct StoredSessionRules {
  std::vector<std::string> static_rule_ids;
  std::vector<PolicyRule> dynamic_rules;
};

struct StoredSessionState {
  StoredSessionConfig config;
  StoredSessionRules rules;
  StoredChargingCreditPool charging_pool;
  StoredUsageMonitoringCreditPool monitor_pool;
  std::string imsi;
  std::string session_id;
  std::string core_session_id;
};

}; // namespace magma

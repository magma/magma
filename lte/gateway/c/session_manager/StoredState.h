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

#include <folly/Format.h>
#include <folly/dynamic.h>
#include <folly/json.h>

#include <lte/protos/pipelined.grpc.pb.h>
#include <lte/protos/session_manager.grpc.pb.h>
#include <lte/protos/session_manager.grpc.pb.h>

#include "CreditKey.h"

namespace magma {
struct QoSInfo {
  bool enabled;
  uint32_t qci;
};

struct SessionConfig {
  std::string ue_ipv4;
  std::string spgw_ipv4;
  std::string msisdn;
  std::string apn;
  std::string imei;
  std::string plmn_id;
  std::string imsi_plmn_id;
  std::string user_location;
  RATType rat_type;
  std::string mac_addr;      // MAC Address for WLAN
  std::string hardware_addr; // MAC Address for WLAN (binary)
  std::string radius_session_id;
  uint32_t bearer_id;
  QoSInfo qos_info;
};

// Session Credit

struct StoredRedirectServer {
  RedirectServer_RedirectAddressType redirect_address_type;
  std::string redirect_server_address;
};

struct FinalActionInfo {
  ChargingCredit_FinalAction final_action;
  RedirectServer redirect_server;
};

enum EventTriggerState {
  PENDING   = 0, // trigger installed
  READY     = 1, // ready to be reported on
  CLEARED   = 2, // successfully reported
};
typedef std::unordered_map<magma::lte::EventTrigger, EventTriggerState> EventTriggerStatus;

/**
 * A bucket is a counter used for tracking credit volume across sessiond.
 * These are independently incremented and reset
 * Each value is in terms of a volume unit - either bytes or seconds
 */
enum Bucket {
  USED_TX = 0,
  USED_RX = 1,
  ALLOWED_TOTAL = 2,
  ALLOWED_TX = 3,
  ALLOWED_RX = 4,
  REPORTING_TX = 5,
  REPORTING_RX = 6,
  REPORTED_TX = 7,
  REPORTED_RX = 8,
  MAX_VALUES = 9,
};

enum ReAuthState {
  REAUTH_NOT_NEEDED = 0,
  REAUTH_REQUIRED = 1,
  REAUTH_PROCESSING = 2,
};

enum ServiceState {
  SERVICE_ENABLED = 0,
  SERVICE_NEEDS_DEACTIVATION = 1,
  SERVICE_DISABLED = 2,
  SERVICE_NEEDS_ACTIVATION = 3,
  SERVICE_REDIRECTED = 4,
  SERVICE_RESTRICTED = 5,
};

enum GrantTrackingType {
  TOTAL_ONLY = 0,
  TX_ONLY = 1,
  RX_ONLY = 2,
  TX_AND_RX = 3,
  NO_TRACKING = 4,
};

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
enum SessionFsmState {
  SESSION_ACTIVE                        = 0,
  SESSION_TERMINATING_FLOW_ACTIVE       = 1,
  SESSION_TERMINATING_AGGREGATING_STATS = 2,
  SESSION_TERMINATING_FLOW_DELETED      = 3,
  SESSION_TERMINATED                    = 4,
  SESSION_TERMINATION_SCHEDULED         = 5
};

struct StoredSessionCredit {
  bool reporting;
  bool is_final;
  CreditLimitType credit_limit_type;
  FinalActionInfo final_action_info;
  ReAuthState reauth_state;
  ServiceState service_state;
  std::time_t expiry_time;
  std::unordered_map<Bucket, uint64_t> buckets;
  GrantTrackingType grant_tracking_type;
};

struct StoredMonitor {
  StoredSessionCredit credit;
  MonitoringLevel level;
};

struct RuleLifetime {
  std::time_t activation_time; // Unix timestamp
  std::time_t deactivation_time; // Unix timestamp
};

typedef std::unordered_map<std::string, StoredMonitor> StoredMonitorMap;
typedef std::unordered_map<CreditKey, StoredSessionCredit, decltype(&ccHash),
                     decltype(&ccEqual)> StoredChargingCreditMap;

struct StoredSessionState {
  SessionFsmState fsm_state;
  SessionConfig config;
  StoredChargingCreditMap credit_map;
  StoredMonitorMap monitor_map;
  std::string session_level_key; // "" maps to nullptr
  std::string imsi;
  std::string session_id;
  std::string core_session_id;
  magma::lte::SubscriberQuotaUpdate_Type subscriber_quota_state;
  magma::lte::TgppContext tgpp_context;
  std::vector<std::string> static_rule_ids;
  std::vector<PolicyRule> dynamic_rules;
  std::vector<PolicyRule> gy_dynamic_rules;
  std::set<std::string> scheduled_static_rules;
  std::vector<PolicyRule> scheduled_dynamic_rules;
  std::unordered_map<std::string, RuleLifetime> rule_lifetimes;
  uint32_t request_number;
  EventTriggerStatus pending_event_triggers;
  google::protobuf::Timestamp revalidation_time;
};

// Update Criteria
struct SessionCreditUpdateCriteria {
  bool is_final;
  FinalActionInfo final_action_info;
  bool reporting;
  ReAuthState reauth_state;
  ServiceState service_state;
  std::time_t expiry_time;
  GrantTrackingType grant_tracking_type;
  // Do not mark REPORTING buckets, but do mark REPORTED
  std::unordered_map<Bucket, uint64_t> bucket_deltas;
};

struct SessionStateUpdateCriteria {
  bool is_session_ended;
  bool is_config_updated;
  SessionConfig updated_config;
  bool is_fsm_updated;
  SessionFsmState updated_fsm_state;
  // true if any of the event trigger state is updated
  bool is_pending_event_triggers_updated;
  EventTriggerStatus pending_event_triggers;
  // this value is only valid if one of the updated event trigger is
  // revalidation time
  google::protobuf::Timestamp revalidation_time;
  uint32_t request_number_increment;

  std::set<std::string> static_rules_to_install;
  std::set<std::string> static_rules_to_uninstall;
  std::set<std::string> new_scheduled_static_rules;
  std::vector<PolicyRule> dynamic_rules_to_install;
  std::vector<PolicyRule> gy_dynamic_rules_to_install;
  std::set<std::string> dynamic_rules_to_uninstall;
  std::set<std::string> gy_dynamic_rules_to_uninstall;
  std::vector<PolicyRule> new_scheduled_dynamic_rules;
  std::unordered_map<std::string, RuleLifetime> new_rule_lifetimes;
  StoredChargingCreditMap charging_credit_to_install;
  std::unordered_map<CreditKey, SessionCreditUpdateCriteria, decltype(&ccHash),
                     decltype(&ccEqual)>
      charging_credit_map;
  StoredMonitorMap monitor_credit_to_install;
  std::unordered_map<std::string, SessionCreditUpdateCriteria>
      monitor_credit_map;
  TgppContext updated_tgpp_context;
  magma::lte::SubscriberQuotaUpdate_Type updated_subscriber_quota_state;
};

SessionStateUpdateCriteria get_default_update_criteria();

std::string serialize_stored_qos_info(const QoSInfo &stored);

QoSInfo deserialize_stored_qos_info(const std::string &serialized);

std::string serialize_stored_session_config(const SessionConfig &stored);

SessionConfig deserialize_stored_session_config(const std::string &serialized);

std::string
serialize_stored_redirect_server(const StoredRedirectServer &stored);

StoredRedirectServer
deserialize_stored_redirect_server(const std::string &serialized);

std::string serialize_stored_final_action_info(const FinalActionInfo &stored);

FinalActionInfo
deserialize_stored_final_action_info(const std::string &serialized);

std::string serialize_stored_session_credit(StoredSessionCredit &stored);

StoredSessionCredit
deserialize_stored_session_credit(const std::string &serialized);

std::string serialize_stored_monitor(StoredMonitor &stored);

StoredMonitor deserialize_stored_monitor(const std::string &serialized);

std::string
serialize_stored_charging_credit_map(StoredChargingCreditMap &stored);

StoredChargingCreditMap
deserialize_stored_charging_credit_map(std::string &serialized) ;

std::string
serialize_stored_usage_monitor_map(StoredMonitorMap &stored);

StoredMonitorMap
deserialize_stored_usage_monitor_map(std::string &serialized);

std::string serialize_stored_session(StoredSessionState &stored);

StoredSessionState deserialize_stored_session(std::string &serialized);
} // namespace magma

/**
 * Copyright 2020 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
#pragma once

#include <folly/Format.h>
#include <folly/dynamic.h>
#include <folly/json.h>
#include <google/protobuf/timestamp.pb.h>
#include <lte/protos/pipelined.grpc.pb.h>
#include <lte/protos/session_manager.grpc.pb.h>
#include <sys/types.h>
#include <cstdint>
#include <experimental/optional>
#include <functional>
#include <set>
#include <string>
#include <unordered_map>
#include <vector>

#include "CreditKey.h"
#include "Types.h"
#include "lte/protos/pipelined.pb.h"
#include "lte/protos/policydb.pb.h"
#include "lte/protos/session_manager.pb.h"

namespace magma {
using std::experimental::optional;

struct StoredSessionCredit {
  bool reporting;
  CreditLimitType credit_limit_type;
  std::unordered_map<Bucket, uint64_t> buckets;
  GrantTrackingType grant_tracking_type;
  GrantedUnits received_granted_units;
  bool report_last_credit;
  uint64_t time_of_first_usage;
  uint64_t time_of_last_usage;
};

struct StoredMonitor {
  StoredSessionCredit credit;
  MonitoringLevel level;
};

struct StoredChargingGrant {
  StoredSessionCredit credit;
  bool is_final;
  FinalActionInfo final_action_info;
  ReAuthState reauth_state;
  ServiceState service_state;
  std::time_t expiry_time;
  bool suspended;
};

typedef std::unordered_map<std::string, StoredMonitor> StoredMonitorMap;
typedef std::unordered_map<CreditKey, StoredChargingGrant, decltype(&ccHash),
                           decltype(&ccEqual)>
    StoredChargingCreditMap;

struct StoredSessionState {
  SessionFsmState fsm_state;
  SessionConfig config;
  StoredChargingCreditMap credit_map;
  StoredMonitorMap monitor_map;
  std::string session_level_key;  // "" maps to nullptr
  std::string imsi;
  uint16_t shard_id;
  std::string session_id;
  uint64_t pdp_start_time;
  uint64_t pdp_end_time;
  // will store the response from the core between Create and Activate Session
  CreateSessionResponse create_session_response;
  // 5G session version handling
  uint32_t current_version;
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
  BearerIDByPolicyID bearer_id_by_policy;
  // converged core PDR rules
  std::vector<SetGroupPDR> pdr_list;
  PolicyStatsMap policy_version_and_stats;
};

// Update Criteria
struct SessionCreditUpdateCriteria {
  // Maintained by ChargingGrant
  bool is_final;
  FinalActionInfo final_action_info;
  ReAuthState reauth_state;
  ServiceState service_state;
  std::time_t expiry_time;

  // Maintained by SessionCredit
  bool reporting;
  GrantTrackingType grant_tracking_type;
  GrantedUnits received_granted_units;

  // Do not mark REPORTING buckets, but do mark REPORTED
  std::unordered_map<Bucket, uint64_t> bucket_deltas;

  bool deleted;
  bool report_last_credit;

  uint64_t time_of_first_usage;
  uint64_t time_of_last_usage;

  bool suspended;
};

struct SessionStateUpdateCriteria {
  bool is_session_ended;
  bool is_config_updated;
  SessionConfig updated_config;
  bool is_fsm_updated;
  SessionFsmState updated_fsm_state;
  // TODO keeping this structure updated for future use.
  bool is_current_version_updated;
  uint32_t updated_current_version;
  // true if any of the event trigger state is updated
  bool is_pending_event_triggers_updated;
  EventTriggerStatus pending_event_triggers;
  // this value is only valid if one of the updated event trigger is
  // revalidation time
  google::protobuf::Timestamp revalidation_time;
  uint32_t request_number_increment;
  uint64_t updated_pdp_end_time;

  // Map to maintain per-policy versions. Contains all values, not delta.
  optional<PolicyStatsMap> policy_version_and_stats;

  std::set<std::string> static_rules_to_install;
  std::set<std::string> static_rules_to_uninstall;
  std::set<std::string> new_scheduled_static_rules;
  std::vector<PolicyRule> dynamic_rules_to_install;
  std::vector<PolicyRule> gy_dynamic_rules_to_install;
  std::set<std::string> dynamic_rules_to_uninstall;
  std::set<std::string> gy_dynamic_rules_to_uninstall;
  std::vector<PolicyRule> new_scheduled_dynamic_rules;
  // Converged rules part of 5G rules
  bool clear_pdr_list;
  std::vector<SetGroupPDR> pdrs_to_install;
  std::unordered_map<std::string, RuleLifetime> new_rule_lifetimes;
  StoredChargingCreditMap charging_credit_to_install;
  std::unordered_map<CreditKey, SessionCreditUpdateCriteria, decltype(&ccHash),
                     decltype(&ccEqual)>
      charging_credit_map;
  bool is_session_level_key_updated;
  std::string updated_session_level_key;
  StoredMonitorMap monitor_credit_to_install;
  std::unordered_map<std::string, SessionCreditUpdateCriteria>
      monitor_credit_map;
  TgppContext updated_tgpp_context;
  // The value should be set only when there is an update
  optional<CreateSessionResponse> create_session_response;
  magma::lte::SubscriberQuotaUpdate_Type updated_subscriber_quota_state;

  bool is_bearer_mapping_updated;
  // Only valid if is_bearer_mapping_updated is true
  BearerIDByPolicyID bearer_id_by_policy;
  Teids teids;
};

SessionStateUpdateCriteria get_default_update_criteria();

std::string serialize_stored_session_config(const SessionConfig& stored);

SessionConfig deserialize_stored_session_config(const std::string& serialized);

std::string serialize_stored_final_action_info(const FinalActionInfo& stored);

FinalActionInfo deserialize_stored_final_action_info(
    const std::string& serialized);

std::string serialize_stored_session_credit(StoredSessionCredit& stored);

StoredSessionCredit deserialize_stored_session_credit(
    const std::string& serialized);

std::string serialize_stored_charging_grant(StoredChargingGrant& stored);

std::string serialize_stored_monitor(StoredMonitor& stored);

StoredMonitor deserialize_stored_monitor(const std::string& serialized);

std::string serialize_stored_charging_credit_map(
    StoredChargingCreditMap& stored);

StoredChargingCreditMap deserialize_stored_charging_credit_map(
    std::string& serialized);

std::string serialize_stored_usage_monitor_map(StoredMonitorMap& stored);

StoredMonitorMap deserialize_stored_usage_monitor_map(std::string& serialized);

BearerIDByPolicyID deserialize_bearer_id_by_policy(std::string& serialized);

std::string serialize_bearer_id_by_policy(BearerIDByPolicyID bearer_map);

std::string serialize_stored_session(StoredSessionState& stored);

StoredSessionState deserialize_stored_session(std::string& serialized);

std::string serialize_policy_stats_map(PolicyStatsMap stats_map);

PolicyStatsMap deserialize_policy_stats_map(std::string& serialized);
}  // namespace magma

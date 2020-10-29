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

#include <functional>
#include <experimental/optional>

#include <folly/Format.h>
#include <folly/dynamic.h>
#include <folly/json.h>

#include <lte/protos/pipelined.grpc.pb.h>
#include <lte/protos/session_manager.grpc.pb.h>
#include <lte/protos/session_manager.grpc.pb.h>

#include "CreditKey.h"

namespace magma {
struct SessionConfig {
  CommonSessionContext common_context;
  RatSpecificContext rat_specific_context;

  SessionConfig(){};
  SessionConfig(const LocalCreateSessionRequest& request);
  bool operator==(const SessionConfig& config) const;
  std::experimental::optional<AggregatedMaximumBitrate> get_apn_ambr() const;
};

// Session Credit
struct FinalActionInfo {
  ChargingCredit_FinalAction final_action;
  RedirectServer redirect_server;
  std::vector<std::string> restrict_rules;
};

enum EventTriggerState {
  PENDING = 0,  // trigger installed
  READY   = 1,  // ready to be reported on
  CLEARED = 2,  // successfully reported
};
typedef std::unordered_map<magma::lte::EventTrigger, EventTriggerState>
    EventTriggerStatus;

/**
 * A bucket is a counter used for tracking credit volume across sessiond.
 * These are independently incremented and reset
 * Each value is in terms of a volume unit - either bytes or seconds
 */
enum Bucket {
  // USED: the actual used quota by the UE.
  // USED = REPORTED + REPORTING
  USED_TX = 0,
  USED_RX = 1,
  // ALLOWED: the granted units received
  ALLOWED_TOTAL = 2,
  ALLOWED_TX    = 3,
  ALLOWED_RX    = 4,
  // REPORTING: quota that is in transit to be acknowledged by OCS/PCRF
  REPORTING_TX = 5,
  REPORTING_RX = 6,
  // REPORTED: quota that has been acknowledged by OCS/PCRF
  REPORTED_TX = 7,
  REPORTED_RX = 8,
  // ALLOWED_FLOOR: saves the previous ALLOWED value after a new grant is
  // received
  // last_valid_nonzero_received_grant = ALLOWED - ALLOWED_FLOOR
  ALLOWED_FLOOR_TOTAL = 9,
  ALLOWED_FLOOR_TX    = 10,
  ALLOWED_FLOOR_RX    = 11,

  // delimiter to iterate enum
  MAX_VALUES = 12,
};

enum ReAuthState {
  REAUTH_NOT_NEEDED = 0,
  REAUTH_REQUIRED   = 1,
  REAUTH_PROCESSING = 2,
};

enum ServiceState {
  SERVICE_ENABLED            = 0,
  SERVICE_NEEDS_DEACTIVATION = 1,
  SERVICE_DISABLED           = 2,
  SERVICE_NEEDS_ACTIVATION   = 3,
  SERVICE_REDIRECTED         = 4,
  SERVICE_RESTRICTED         = 5,
};

enum GrantTrackingType {
  TRACKING_UNSET  = -1,
  TOTAL_ONLY      = 0,
  TX_ONLY         = 1,
  RX_ONLY         = 2,
  TX_AND_RX       = 3,
  ALL_TOTAL_TX_RX = 4,
};

/**
 * State transitions of a session:
 * SESSION_ACTIVE
 *       |
 *       |
 *       |
 *       V
 * SESSION_RELEASED
 *       |
 *       | (PipelineD enforcement flows get deleted OR forced timeout)
 *       |      -> complete_termination
 *       V
 * SESSION_TERMINATED
 */
enum SessionFsmState {
  SESSION_ACTIVE     = 0,
  SESSION_TERMINATED = 4,
  SESSION_RELEASED   = 6,
  CREATING           = 7,
  CREATED            = 8,
  ACTIVE             = 9,
  INACTIVE           = 10,
  RELEASE            = 11,
};

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
};

struct RuleLifetime {
  std::time_t activation_time;    // Unix timestamp
  std::time_t deactivation_time;  // Unix timestamp
};

// QoS Management
enum PolicyType {
  STATIC  = 1,
  DYNAMIC = 2,
};

struct PolicyID {
  PolicyType policy_type;
  std::string rule_id;

  PolicyID(PolicyType p_type, std::string r_id) {
    policy_type = p_type;
    rule_id     = r_id;
  }

  bool operator==(const PolicyID& id) const {
    return rule_id == id.rule_id && policy_type == id.policy_type;
  }
};

// Custom hash for PolicyID
struct PolicyIDHash {
  std::size_t operator()(const PolicyID& id) const {
    std::size_t h1 = std::hash<std::string>{}(id.rule_id);
    std::size_t h2 = std::hash<int>{}(int(id.policy_type));
    return h1 ^ (h2 << 1);
  }
};

typedef std::unordered_map<std::string, StoredMonitor> StoredMonitorMap;
typedef std::unordered_map<
    CreditKey, StoredChargingGrant, decltype(&ccHash), decltype(&ccEqual)>
    StoredChargingCreditMap;
typedef std::unordered_map<PolicyID, uint32_t, PolicyIDHash> BearerIDByPolicyID;

struct StoredSessionState {
  SessionFsmState fsm_state;
  SessionConfig config;
  StoredChargingCreditMap credit_map;
  StoredMonitorMap monitor_map;
  std::string session_level_key;  // "" maps to nullptr
  std::string imsi;
  std::string session_id;
  uint64_t pdp_start_time;
  uint64_t pdp_end_time;
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
  uint64_t updated_pdp_end_time;

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
  std::unordered_map<
      CreditKey, SessionCreditUpdateCriteria, decltype(&ccHash),
      decltype(&ccEqual)>
      charging_credit_map;
  bool is_session_level_key_updated;
  std::string updated_session_level_key;
  StoredMonitorMap monitor_credit_to_install;
  std::unordered_map<std::string, SessionCreditUpdateCriteria>
      monitor_credit_map;
  TgppContext updated_tgpp_context;
  magma::lte::SubscriberQuotaUpdate_Type updated_subscriber_quota_state;

  bool is_bearer_mapping_updated;
  // Only valid if is_bearer_mapping_updated is true
  BearerIDByPolicyID bearer_id_by_policy;
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
}  // namespace magma

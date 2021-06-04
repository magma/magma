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

#include <experimental/optional>
#include <folly/dynamic.h>
#include <folly/Format.h>
#include <folly/json.h>
#include <lte/protos/pipelined.grpc.pb.h>
#include <lte/protos/session_manager.grpc.pb.h>

#include <functional>
#include <string>
#include <unordered_map>
#include <vector>

#include "CreditKey.h"

// NOTE:
// This file is intended for declaring types that are shared across classes.
// If a type has a clear owner, do NOT put in this file

namespace magma {
struct SessionConfig {
  CommonSessionContext common_context;
  RatSpecificContext rat_specific_context;

  SessionConfig() {}
  explicit SessionConfig(const LocalCreateSessionRequest& request);
  bool operator==(const SessionConfig& config) const;
  std::experimental::optional<AggregatedMaximumBitrate> get_apn_ambr() const;
  AggregatedMaximumBitrate_BitrateUnitsAMBR get_apn_ambr_units(
      QosInformationRequest_BitrateUnitsAMBR units) const;
  std::string get_imsi() const { return common_context.sid().id(); }
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

struct Usage {
  uint64_t bytes_tx;
  uint64_t bytes_rx;
};

struct TotalCreditUsage {
  uint64_t monitoring_tx;
  uint64_t monitoring_rx;
  uint64_t charging_tx;
  uint64_t charging_rx;
};

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
  SERVICE_NEEDS_SUSPENSION   = 2,
  SERVICE_DISABLED           = 3,
  SERVICE_NEEDS_ACTIVATION   = 4,
  SERVICE_REDIRECTED         = 5,
  SERVICE_RESTRICTED         = 6,
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
  SESSION_ACTIVE                = 0,
  SESSION_TERMINATED            = 4,
  SESSION_TERMINATION_SCHEDULED = 5,
  SESSION_RELEASED              = 6,
  CREATING                      = 7,
  CREATED                       = 8,
  ACTIVE                        = 9,
  INACTIVE                      = 10,
  RELEASE                       = 11,
};

struct RuleLifetime {
  std::time_t activation_time;    // Unix timestamp
  std::time_t deactivation_time;  // Unix timestamp
  RuleLifetime() : activation_time(0), deactivation_time(0) {}
  RuleLifetime(const time_t activation, const time_t deactivation)
      : activation_time(activation), deactivation_time(deactivation) {}
  explicit RuleLifetime(const StaticRuleInstall& rule_install);
  explicit RuleLifetime(const DynamicRuleInstall& rule_install);
  bool is_within_lifetime(std::time_t time);
  bool exceeded_lifetime(std::time_t time);
  bool before_lifetime(std::time_t time);
  bool should_schedule_deactivation(std::time_t time);
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

bool operator==(const Teids& lhs, const Teids& rhs);

struct BearerIDAndTeid {
  uint32_t bearer_id;
  Teids teids;

  bool operator==(const BearerIDAndTeid& id) const {
    return bearer_id == id.bearer_id && teids == id.teids;
  }
};

typedef std::unordered_map<PolicyID, BearerIDAndTeid, PolicyIDHash>
    BearerIDByPolicyID;

struct RuleToProcess {
  PolicyRule rule;
  uint32_t version;
  Teids teids;
};

struct RuleStats {
  uint64_t tx;
  uint64_t rx;
  uint64_t dropped_tx;
  uint64_t dropped_rx;
  RuleStats() : tx(0), rx(0), dropped_tx(0), dropped_rx(0) {}
  RuleStats(uint64_t tx, uint64_t rx, uint64_t dropped_tx, uint64_t dropped_rx)
      : tx(tx), rx(rx), dropped_tx(dropped_tx), dropped_rx(dropped_rx) {}
};

typedef std::vector<RuleToProcess> RulesToProcess;

enum PolicyAction {
  ACTIVATE   = 0,
  DEACTIVATE = 1,
};

struct RuleToSchedule {
  PolicyType p_type;
  std::string rule_id;
  PolicyAction p_action;
  std::time_t scheduled_time;
  RuleToSchedule() {}
  RuleToSchedule(
      PolicyType _p_type, std::string _rule_id, PolicyAction _p_action,
      std::time_t _time)
      : p_type(_p_type),
        rule_id(_rule_id),
        p_action(_p_action),
        scheduled_time(_time) {}
};

typedef std::vector<RuleToSchedule> RulesToSchedule;

struct StatsPerPolicy {
  // The version maintained by SessionD for this rule
  uint32_t current_version;
  // The last reported version from PipelineD
  uint32_t last_reported_version;

  std::unordered_map<int, RuleStats> stats_map;
  StatsPerPolicy() {
    current_version       = 0;
    last_reported_version = 0;
    RuleStats s           = {0, 0, 0, 0};
    stats_map             = {{0, s}};
  }
};
typedef std::unordered_map<std::string, StatsPerPolicy> PolicyStatsMap;

struct TeidHash {
  std::size_t operator()(const Teids& teid) const {
    std::size_t h1 = std::hash<uint32_t>{}(teid.enb_teid());
    std::size_t h2 = std::hash<uint32_t>{}(teid.agw_teid());
    return h1 ^ h2;
  }
};
struct TeidEqual {
  bool operator()(const Teids& lhs, const Teids& rhs) const {
    return lhs.agw_teid() == rhs.agw_teid() && lhs.enb_teid() == rhs.enb_teid();
  }
};
typedef std::unordered_map<Teids, ActivateFlowsRequest, TeidHash, TeidEqual>
    ActivateReqByTeids;
typedef std::unordered_map<Teids, DeactivateFlowsRequest, TeidHash, TeidEqual>
    DeactivateReqByTeids;

}  // namespace magma

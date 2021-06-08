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

#include <google/protobuf/timestamp.pb.h>
#include <google/protobuf/util/time_util.h>
#include <lte/protos/pipelined.grpc.pb.h>
#include <lte/protos/session_manager.grpc.pb.h>

#include <string>
#include <unordered_map>

#include "CreditKey.h"
#include "magma_logging.h"
#include "StoredState.h"

namespace magma {
using google::protobuf::util::TimeUtil;

SessionConfig::SessionConfig(const LocalCreateSessionRequest& request) {
  common_context       = request.common_context();
  rat_specific_context = request.rat_specific_context();
}

bool SessionConfig::operator==(const SessionConfig& config) const {
  auto common1 = common_context.SerializeAsString();
  auto common2 = config.common_context.SerializeAsString();
  if (common1 != common2) {
    return false;
  }
  std::string current_rat_specific = rat_specific_context.SerializeAsString();
  std::string new_rat_specific =
      config.rat_specific_context.SerializeAsString();
  return current_rat_specific == new_rat_specific;
}

optional<AggregatedMaximumBitrate> SessionConfig::get_apn_ambr() const {
  if (rat_specific_context.has_lte_context() &&
      rat_specific_context.lte_context().has_qos_info()) {
    AggregatedMaximumBitrate max_bitrate;

    const auto& qos_info = rat_specific_context.lte_context().qos_info();
    max_bitrate.set_max_bandwidth_ul(qos_info.apn_ambr_ul());
    max_bitrate.set_max_bandwidth_dl(qos_info.apn_ambr_dl());
    max_bitrate.set_br_unit(get_apn_ambr_units(qos_info.br_unit()));
    return max_bitrate;
  }
  return {};
}

AggregatedMaximumBitrate_BitrateUnitsAMBR SessionConfig::get_apn_ambr_units(
    QosInformationRequest_BitrateUnitsAMBR units) const {
  switch (units) {
    case QosInformationRequest_BitrateUnitsAMBR_BPS:
      return AggregatedMaximumBitrate_BitrateUnitsAMBR_BPS;
    case QosInformationRequest_BitrateUnitsAMBR_KBPS:
      return AggregatedMaximumBitrate_BitrateUnitsAMBR_KBPS;
    default:
      MLOG(MERROR) << "QOS bitrate unit not implemented (" << units
                   << "). Setting BPS";
      return AggregatedMaximumBitrate_BitrateUnitsAMBR_BPS;
  }
}

SessionStateUpdateCriteria get_default_update_criteria() {
  SessionStateUpdateCriteria uc{};
  uc.is_fsm_updated             = false;
  uc.is_config_updated          = false;
  uc.request_number_increment   = 0;
  uc.updated_pdp_end_time       = 0;
  uc.charging_credit_to_install = StoredChargingCreditMap(4, &ccHash, &ccEqual);
  uc.charging_credit_map        = std::unordered_map<
      CreditKey, SessionCreditUpdateCriteria, decltype(&ccHash),
      decltype(&ccEqual)>(4, &ccHash, &ccEqual);
  uc.is_session_level_key_updated = false;
  uc.is_bearer_mapping_updated    = false;
  uc.policy_version_and_stats     = {};
  return uc;
}

std::string serialize_stored_session_config(const SessionConfig& stored) {
  folly::dynamic marshaled    = folly::dynamic::object;
  marshaled["common_context"] = stored.common_context.SerializeAsString();
  marshaled["rat_specific_context"] =
      stored.rat_specific_context.SerializeAsString();

  std::string serialized = folly::toJson(marshaled);
  return serialized;
}

SessionConfig deserialize_stored_session_config(const std::string& serialized) {
  auto folly_serialized    = folly::StringPiece(serialized);
  folly::dynamic marshaled = folly::parseJson(folly_serialized);

  auto stored = SessionConfig{};
  magma::lte::CommonSessionContext common_context;
  common_context.ParseFromString(marshaled["common_context"].getString());
  stored.common_context = common_context;

  magma::lte::RatSpecificContext rat_specific_context;
  rat_specific_context.ParseFromString(
      marshaled["rat_specific_context"].getString());
  stored.rat_specific_context = rat_specific_context;

  return stored;
}

std::string serialize_stored_final_action_info(const FinalActionInfo& stored) {
  folly::dynamic marshaled  = folly::dynamic::object;
  marshaled["final_action"] = static_cast<int>(stored.final_action);

  std::string redirect_server;
  stored.redirect_server.SerializeToString(&redirect_server);
  marshaled["redirect_server"] = redirect_server;

  folly::dynamic restrict_rules = folly::dynamic::array;
  for (const auto& rule_id : stored.restrict_rules) {
    restrict_rules.push_back(rule_id);
  }
  marshaled["restrict_rules"] = restrict_rules;

  std::string serialized = folly::toJson(marshaled);
  return serialized;
}

FinalActionInfo deserialize_stored_final_action_info(
    const std::string& serialized) {
  auto folly_serialized    = folly::StringPiece(serialized);
  folly::dynamic marshaled = folly::parseJson(folly_serialized);

  auto stored         = FinalActionInfo{};
  stored.final_action = static_cast<ChargingCredit_FinalAction>(
      marshaled["final_action"].getInt());

  magma::lte::RedirectServer redirect_server;
  redirect_server.ParseFromString(marshaled["redirect_server"].getString());
  stored.redirect_server = redirect_server;

  for (auto& rule_id : marshaled["restrict_rules"]) {
    stored.restrict_rules.push_back(rule_id.getString());
  }
  return stored;
}

std::string serialize_stored_charging_grant(StoredChargingGrant& stored) {
  folly::dynamic marshaled = folly::dynamic::object;
  marshaled["is_final"]    = stored.is_final;
  marshaled["final_action_info"] =
      serialize_stored_final_action_info(stored.final_action_info);
  marshaled["expiry_time"]   = std::to_string(stored.expiry_time);
  marshaled["reauth_state"]  = static_cast<int>(stored.reauth_state);
  marshaled["service_state"] = static_cast<int>(stored.service_state);
  marshaled["credit"]        = serialize_stored_session_credit(stored.credit);
  marshaled["suspended"]     = stored.suspended;

  std::string serialized = folly::toJson(marshaled);
  return serialized;
}

StoredChargingGrant deserialize_stored_charging_grant(
    const std::string& serialized) {
  auto folly_serialized    = folly::StringPiece(serialized);
  folly::dynamic marshaled = folly::parseJson(folly_serialized);

  auto stored              = StoredChargingGrant{};
  stored.is_final          = marshaled["is_final"].getBool();
  stored.final_action_info = deserialize_stored_final_action_info(
      marshaled["final_action_info"].getString());
  stored.reauth_state =
      static_cast<ReAuthState>(marshaled["reauth_state"].getInt());
  stored.service_state =
      static_cast<ServiceState>(marshaled["service_state"].getInt());
  stored.expiry_time = static_cast<std::time_t>(
      std::stoul(marshaled["expiry_time"].getString()));
  stored.credit =
      deserialize_stored_session_credit(marshaled["credit"].getString());
  stored.suspended = marshaled["suspended"].getBool();

  return stored;
}

std::string serialize_stored_session_credit(StoredSessionCredit& stored) {
  folly::dynamic marshaled       = folly::dynamic::object;
  marshaled["reporting"]         = stored.reporting;
  marshaled["credit_limit_type"] = static_cast<int>(stored.credit_limit_type);
  marshaled["buckets"]           = folly::dynamic::object();
  marshaled["grant_tracking_type"] =
      static_cast<int>(stored.grant_tracking_type);
  marshaled["received_granted_units"] =
      stored.received_granted_units.SerializeAsString();
  marshaled["report_last_credit"]  = stored.report_last_credit;
  marshaled["time_of_first_usage"] = std::to_string(stored.time_of_first_usage);
  marshaled["time_of_last_usage"]  = std::to_string(stored.time_of_last_usage);

  for (int bucket_int = USED_TX; bucket_int != MAX_VALUES; bucket_int++) {
    Bucket bucket = static_cast<Bucket>(bucket_int);
    marshaled["buckets"][std::to_string(bucket_int)] =
        std::to_string(stored.buckets[bucket]);
  }

  std::string serialized = folly::toJson(marshaled);
  return serialized;
}

StoredSessionCredit deserialize_stored_session_credit(
    const std::string& serialized) {
  auto folly_serialized    = folly::StringPiece(serialized);
  folly::dynamic marshaled = folly::parseJson(folly_serialized);

  auto stored      = StoredSessionCredit{};
  stored.reporting = marshaled["reporting"].getBool();
  stored.credit_limit_type =
      static_cast<CreditLimitType>(marshaled["credit_limit_type"].getInt());
  stored.grant_tracking_type =
      static_cast<GrantTrackingType>(marshaled["grant_tracking_type"].getInt());

  GrantedUnits received_granted_units;
  received_granted_units.ParseFromString(
      marshaled["received_granted_units"].getString());
  stored.received_granted_units = received_granted_units;
  stored.report_last_credit     = marshaled["report_last_credit"].getBool();
  stored.time_of_first_usage    = static_cast<uint64_t>(
      std::stoul(marshaled["time_of_first_usage"].getString()));
  stored.time_of_last_usage = static_cast<uint64_t>(
      std::stoul(marshaled["time_of_last_usage"].getString()));

  for (int bucket_int = USED_TX; bucket_int != MAX_VALUES; bucket_int++) {
    Bucket bucket          = static_cast<Bucket>(bucket_int);
    stored.buckets[bucket] = static_cast<uint64_t>(std::stoul(
        marshaled["buckets"][std::to_string(bucket_int)].getString()));
  }

  return stored;
}

std::string serialize_stored_monitor(StoredMonitor& stored) {
  folly::dynamic marshaled = folly::dynamic::object;

  marshaled["credit"] = serialize_stored_session_credit(stored.credit);
  marshaled["level"]  = static_cast<int>(stored.level);

  std::string serialized = folly::toJson(marshaled);
  return serialized;
}

StoredMonitor deserialize_stored_monitor(const std::string& serialized) {
  auto folly_serialized    = folly::StringPiece(serialized);
  folly::dynamic marshaled = folly::parseJson(folly_serialized);

  auto stored = StoredMonitor{};
  stored.credit =
      deserialize_stored_session_credit(marshaled["credit"].getString());
  stored.level = static_cast<MonitoringLevel>(marshaled["level"].getInt());

  return stored;
}

std::string serialize_stored_charging_credit_map(
    StoredChargingCreditMap& stored) {
  folly::dynamic marshaled = folly::dynamic::object;

  folly::dynamic credit_keys = folly::dynamic::array;
  folly::dynamic credit_map  = folly::dynamic::object;
  for (auto& credit_pair : stored) {
    CreditKey credit_key       = credit_pair.first;
    folly::dynamic key2        = folly::dynamic::object;
    key2["rating_group"]       = std::to_string(credit_key.rating_group);
    key2["service_identifier"] = std::to_string(credit_key.service_identifier);
    credit_keys.push_back(key2);

    std::string key = std::to_string(credit_key.rating_group) +
                      std::to_string(credit_key.service_identifier);
    credit_map[key] = serialize_stored_charging_grant(credit_pair.second);
  }
  marshaled["credit_keys"] = credit_keys;
  marshaled["credit_map"]  = credit_map;

  std::string serialized = folly::toJson(marshaled);
  return serialized;
}

StoredChargingCreditMap deserialize_stored_charging_credit_map(
    std::string& serialized) {
  auto folly_serialized    = folly::StringPiece(serialized);
  folly::dynamic marshaled = folly::parseJson(folly_serialized);

  auto stored = StoredChargingCreditMap(4, &ccHash, &ccEqual);

  for (auto& key : marshaled["credit_keys"]) {
    auto credit_key = CreditKey(
        static_cast<uint32_t>(std::stoul(key["rating_group"].getString())),
        static_cast<uint32_t>(
            std::stoul(key["service_identifier"].getString())));

    std::string key2 =
        key["rating_group"].getString() + key["service_identifier"].getString();

    stored[credit_key] = deserialize_stored_charging_grant(
        marshaled["credit_map"][key2].getString());
  }
  return stored;
}

std::string serialize_stored_usage_monitor_map(StoredMonitorMap& stored) {
  folly::dynamic marshaled = folly::dynamic::object;

  folly::dynamic monitor_keys = folly::dynamic::array;
  folly::dynamic monitor_map  = folly::dynamic::object;
  for (auto& monitor_pair : stored) {
    monitor_keys.push_back(monitor_pair.first);
    monitor_map[monitor_pair.first] =
        serialize_stored_monitor(monitor_pair.second);
  }
  marshaled["monitor_keys"] = monitor_keys;
  marshaled["monitor_map"]  = monitor_map;

  std::string serialized = folly::toJson(marshaled);
  return serialized;
}

StoredMonitorMap deserialize_stored_usage_monitor_map(std::string& serialized) {
  auto folly_serialized    = folly::StringPiece(serialized);
  folly::dynamic marshaled = folly::parseJson(folly_serialized);
  auto stored              = StoredMonitorMap{};
  for (auto& key : marshaled["monitor_keys"]) {
    std::string monitor_key = key.getString();
    stored[monitor_key] =
        deserialize_stored_monitor(marshaled["monitor_map"][key].getString());
  }

  return stored;
}

EventTriggerStatus deserialize_pending_event_triggers(std::string& serialized) {
  auto folly_serialized    = folly::StringPiece(serialized);
  folly::dynamic marshaled = folly::parseJson(folly_serialized);

  auto stored = EventTriggerStatus{};
  for (auto& key : marshaled["event_trigger_keys"]) {
    auto map = marshaled["event_trigger_map"];
    magma::lte::EventTrigger eventKey;
    try {
      eventKey = magma::lte::EventTrigger(std::stoi(key.getString()));
    } catch (std::invalid_argument const& e) {
      MLOG(MWARNING) << "Could not deserialize event triggers";
      continue;
    }
    stored[eventKey] = EventTriggerState(map[key].getInt());
  }

  return stored;
}

std::string serialize_pending_event_triggers(
    EventTriggerStatus event_triggers) {
  folly::dynamic marshaled = folly::dynamic::object;

  folly::dynamic keys = folly::dynamic::array;
  folly::dynamic map  = folly::dynamic::object;
  for (auto& trigger_pair : event_triggers) {
    auto key = std::to_string(int(trigger_pair.first));
    keys.push_back(key);
    map[key] = int(trigger_pair.second);
  }
  marshaled["event_trigger_keys"] = keys;
  marshaled["event_trigger_map"]  = map;

  std::string serialized = folly::toJson(marshaled);
  return serialized;
}

RuleStats deserialize_rule_stats(std::string& serialized) {
  auto folly_serialized    = folly::StringPiece(serialized);
  folly::dynamic marshaled = folly::parseJson(folly_serialized);
  auto stored              = RuleStats{};

  stored.tx = static_cast<uint64_t>(std::stoul(marshaled["tx"].getString()));
  stored.rx = static_cast<uint64_t>(std::stoul(marshaled["rx"].getString()));
  stored.dropped_rx =
      static_cast<uint64_t>(std::stoul(marshaled["dropped_rx"].getString()));
  stored.dropped_tx =
      static_cast<uint64_t>(std::stoul(marshaled["dropped_tx"].getString()));

  return stored;
}

PolicyStatsMap deserialize_policy_stats_map(std::string& serialized) {
  auto folly_serialized    = folly::StringPiece(serialized);
  folly::dynamic marshaled = folly::parseJson(folly_serialized);

  auto stored = PolicyStatsMap{};
  auto map    = marshaled["policy_stats_map"];
  for (auto& key : marshaled["policy_stats_keys"]) {
    StatsPerPolicy stats;
    stats.current_version =
        static_cast<uint32_t>(map[key]["current_version"].getInt());
    stats.last_reported_version =
        static_cast<uint32_t>(map[key]["last_reported_version"].getInt());
    for (auto& key2 : map[key]["stats_keys"]) {
      int stat_key = static_cast<uint64_t>(std::stoul(key2.getString()));
      stats.stats_map[stat_key] =
          deserialize_rule_stats(map[key]["stats_map"][key2].getString());
    }
    stored[key.getString()] = stats;
  }
  return stored;
}

std::string serialize_policy_stats_map(PolicyStatsMap stats_map) {
  folly::dynamic marshaled = folly::dynamic::object;

  folly::dynamic keys = folly::dynamic::array;
  folly::dynamic map  = folly::dynamic::object;

  for (auto& policy_pair : stats_map) {
    auto key = policy_pair.first;
    keys.push_back(key);
    folly::dynamic usage = folly::dynamic::object;
    usage["current_version"] =
        static_cast<int>(policy_pair.second.current_version);
    usage["last_reported_version"] =
        static_cast<int>(policy_pair.second.last_reported_version);

    folly::dynamic stats_keys = folly::dynamic::array;
    folly::dynamic stats_map  = folly::dynamic::object;
    for (auto& stat : policy_pair.second.stats_map) {
      std::string version_str = std::to_string(stat.first);
      stats_keys.push_back(version_str);
      folly::dynamic stats   = folly::dynamic::object;
      stats["tx"]            = std::to_string(stat.second.tx);
      stats["rx"]            = std::to_string(stat.second.rx);
      stats["dropped_tx"]    = std::to_string(stat.second.dropped_tx);
      stats["dropped_rx"]    = std::to_string(stat.second.dropped_rx);
      stats_map[version_str] = folly::toJson(stats);
    }
    usage["stats_keys"] = stats_keys;
    usage["stats_map"]  = stats_map;

    map[key] = usage;
  }

  marshaled["policy_stats_keys"] = keys;
  marshaled["policy_stats_map"]  = map;

  std::string serialized = folly::toJson(marshaled);
  return serialized;
}

BearerIDByPolicyID deserialize_bearer_id_by_policy(std::string& serialized) {
  auto folly_serialized    = folly::StringPiece(serialized);
  folly::dynamic marshaled = folly::parseJson(folly_serialized);

  auto stored = BearerIDByPolicyID{};
  for (auto& bearer_id_by_policy : marshaled) {
    PolicyType policy_type = PolicyType(bearer_id_by_policy["type"].getInt());
    std::string rule_id    = bearer_id_by_policy["rule_id"].getString();
    auto policy_id         = PolicyID(policy_type, rule_id);
    stored[policy_id]      = BearerIDAndTeid();
    stored[policy_id].bearer_id =
        static_cast<uint32_t>(bearer_id_by_policy["bearer_id"].getInt());
    stored[policy_id].teids.set_agw_teid(
        static_cast<uint32_t>(bearer_id_by_policy["agw_teid"].getInt()));
    stored[policy_id].teids.set_enb_teid(
        static_cast<uint32_t>(bearer_id_by_policy["enb_teid"].getInt()));
  }
  return stored;
}

std::string serialize_bearer_id_by_policy(BearerIDByPolicyID bearer_map) {
  folly::dynamic marshaled = folly::dynamic::array;

  for (auto& pair : bearer_map) {
    folly::dynamic bearer_id_by_policy = folly::dynamic::object;
    bearer_id_by_policy["type"]      = static_cast<int>(pair.first.policy_type);
    bearer_id_by_policy["rule_id"]   = pair.first.rule_id;
    bearer_id_by_policy["bearer_id"] = static_cast<int>(pair.second.bearer_id);
    bearer_id_by_policy["agw_teid"] =
        static_cast<int>(pair.second.teids.agw_teid());
    bearer_id_by_policy["enb_teid"] =
        static_cast<int>(pair.second.teids.enb_teid());
    marshaled.push_back(bearer_id_by_policy);
  }
  std::string serialized = folly::toJson(marshaled);
  return serialized;
}

std::string serialize_stored_session(StoredSessionState& stored) {
  folly::dynamic marshaled = folly::dynamic::object;
  marshaled["fsm_state"]   = static_cast<int>(stored.fsm_state);
  marshaled["config"]      = serialize_stored_session_config(stored.config);
  marshaled["charging_pool"] =
      serialize_stored_charging_credit_map(stored.credit_map);
  marshaled["monitor_map"] =
      serialize_stored_usage_monitor_map(stored.monitor_map);
  marshaled["session_level_key"] = stored.session_level_key;
  marshaled["imsi"]              = stored.imsi;
  marshaled["session_id"]        = stored.session_id;
  marshaled["local_teid"]        = std::to_string(stored.local_teid);
  marshaled["subscriber_quota_state"] =
      static_cast<int>(stored.subscriber_quota_state);
  marshaled["create_session_response"] =
      stored.create_session_response.SerializeAsString();

  marshaled["tgpp_context"]   = stored.tgpp_context.SerializeAsString();
  marshaled["pdp_start_time"] = std::to_string(stored.pdp_start_time);
  marshaled["pdp_end_time"]   = std::to_string(stored.pdp_end_time);

  marshaled["pending_event_triggers"] =
      serialize_pending_event_triggers(stored.pending_event_triggers);
  std::string revalidation_time;
  stored.revalidation_time.SerializeToString(&revalidation_time);
  marshaled["revalidation_time"] = revalidation_time;

  marshaled["bearer_id_by_policy"] =
      serialize_bearer_id_by_policy(stored.bearer_id_by_policy);

  marshaled["policy_version_and_stats"] =
      serialize_policy_stats_map(stored.policy_version_and_stats);

  folly::dynamic static_rule_ids = folly::dynamic::array;
  for (const auto& rule_id : stored.static_rule_ids) {
    static_rule_ids.push_back(rule_id);
  }
  marshaled["static_rule_ids"] = static_rule_ids;

  folly::dynamic dynamic_rules = folly::dynamic::array;
  for (const auto& rule : stored.dynamic_rules) {
    std::string dynamic_rule;
    rule.SerializeToString(&dynamic_rule);
    dynamic_rules.push_back(dynamic_rule);
  }
  marshaled["dynamic_rules"] = dynamic_rules;

  folly::dynamic gy_dynamic_rules = folly::dynamic::array;
  for (const auto& rule : stored.gy_dynamic_rules) {
    std::string gy_dynamic_rule;
    rule.SerializeToString(&gy_dynamic_rule);
    gy_dynamic_rules.push_back(gy_dynamic_rule);
  }
  marshaled["gy_dynamic_rules"] = gy_dynamic_rules;

  marshaled["request_number"] = std::to_string(stored.request_number);

  std::string serialized = folly::toJson(marshaled);
  return serialized;
}

StoredSessionState deserialize_stored_session(std::string& serialized) {
  auto folly_serialized    = folly::StringPiece(serialized);
  folly::dynamic marshaled = folly::parseJson(folly_serialized);

  auto stored      = StoredSessionState{};
  stored.fsm_state = SessionFsmState(marshaled["fsm_state"].getInt());
  stored.config =
      deserialize_stored_session_config(marshaled["config"].getString());
  stored.credit_map = deserialize_stored_charging_credit_map(
      marshaled["charging_pool"].getString());
  stored.monitor_map = deserialize_stored_usage_monitor_map(
      marshaled["monitor_map"].getString());
  stored.session_level_key = marshaled["session_level_key"].getString();
  stored.imsi              = marshaled["imsi"].getString();
  stored.session_id        = marshaled["session_id"].getString();
  stored.local_teid =
      static_cast<uint32_t>(std::stoul(marshaled["local_teid"].getString()));
  stored.subscriber_quota_state =
      static_cast<magma::lte::SubscriberQuotaUpdate_Type>(
          marshaled["subscriber_quota_state"].getInt());

  google::protobuf::Timestamp revalidation_time;
  revalidation_time.ParseFromString(marshaled["revalidation_time"].getString());
  stored.revalidation_time      = revalidation_time;
  stored.pending_event_triggers = deserialize_pending_event_triggers(
      marshaled["pending_event_triggers"].getString());

  stored.bearer_id_by_policy = deserialize_bearer_id_by_policy(
      marshaled["bearer_id_by_policy"].getString());

  stored.policy_version_and_stats = deserialize_policy_stats_map(
      marshaled["policy_version_and_stats"].getString());

  CreateSessionResponse csr;
  csr.ParseFromString(marshaled["create_session_response"].getString());
  stored.create_session_response = csr;

  magma::lte::TgppContext tgpp_context;
  tgpp_context.ParseFromString(marshaled["tgpp_context"].getString());
  stored.tgpp_context = tgpp_context;

  for (auto& rule_id : marshaled["static_rule_ids"]) {
    stored.static_rule_ids.push_back(rule_id.getString());
  }

  for (auto& policy : marshaled["dynamic_rules"]) {
    PolicyRule policy_rule;
    policy_rule.ParseFromString(policy.getString());
    stored.dynamic_rules.push_back(policy_rule);
  }

  for (auto& policy : marshaled["gy_dynamic_rules"]) {
    PolicyRule policy_rule;
    policy_rule.ParseFromString(policy.getString());
    stored.gy_dynamic_rules.push_back(policy_rule);
  }

  stored.request_number = static_cast<uint32_t>(
      std::stoul(marshaled["request_number"].getString()));

  stored.pdp_start_time = static_cast<uint64_t>(
      std::stoul(marshaled["pdp_start_time"].getString()));
  stored.pdp_end_time =
      static_cast<uint64_t>(std::stoul(marshaled["pdp_end_time"].getString()));

  return stored;
}

RuleLifetime::RuleLifetime(const StaticRuleInstall& rule_install) {
  activation_time =
      std::time_t(TimeUtil::TimestampToSeconds(rule_install.activation_time()));
  deactivation_time = std::time_t(
      TimeUtil::TimestampToSeconds(rule_install.deactivation_time()));
}

RuleLifetime::RuleLifetime(const DynamicRuleInstall& rule_install) {
  activation_time =
      std::time_t(TimeUtil::TimestampToSeconds(rule_install.activation_time()));
  deactivation_time = std::time_t(
      TimeUtil::TimestampToSeconds(rule_install.deactivation_time()));
}

bool RuleLifetime::is_within_lifetime(std::time_t time) {
  auto past_activation_time = activation_time <= time;
  auto before_deactivation_time =
      (deactivation_time == 0) || (time < deactivation_time);
  return past_activation_time && before_deactivation_time;
}

bool RuleLifetime::exceeded_lifetime(std::time_t time) {
  return deactivation_time != 0 && deactivation_time <= time;
}

bool RuleLifetime::before_lifetime(std::time_t time) {
  return time < activation_time;
}

bool RuleLifetime::should_schedule_deactivation(std::time_t time) {
  return deactivation_time != 0 && time <= deactivation_time;
}

};  // namespace magma

/**
 * Copyright (c) 2016-present, Facebook, Inc.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include "StoredState.h"
#include "CreditKey.h"
#include "magma_logging.h"

namespace magma {

SessionStateUpdateCriteria get_default_update_criteria() {
  SessionStateUpdateCriteria uc{};
  uc.is_fsm_updated = false;
  uc.is_config_updated = false;
  uc.request_number_increment = 0;
  uc.charging_credit_to_install = StoredChargingCreditMap(4, &ccHash, &ccEqual);
  uc.charging_credit_map =
      std::unordered_map<CreditKey, SessionCreditUpdateCriteria,
                         decltype(&ccHash), decltype(&ccEqual)>(4, &ccHash,
                                                                &ccEqual);
  return uc;
}

std::string serialize_stored_qos_info(const QoSInfo &stored) {
  folly::dynamic marshaled = folly::dynamic::object;
  marshaled["enabled"] = stored.enabled;
  marshaled["qci"] = std::to_string(stored.qci);

  std::string serialized = folly::toJson(marshaled);
  return serialized;
}

QoSInfo deserialize_stored_qos_info(const std::string &serialized) {
  auto folly_serialized = folly::StringPiece(serialized);
  folly::dynamic marshaled = folly::parseJson(folly_serialized);

  auto stored = QoSInfo{};
  stored.enabled = marshaled["enabled"].getBool();
  stored.qci = static_cast<uint32_t>(std::stoul(marshaled["qci"].getString()));

  return stored;
}

std::string serialize_stored_session_config(const SessionConfig &stored) {
  folly::dynamic marshaled = folly::dynamic::object;
  marshaled["ue_ipv4"] = stored.ue_ipv4;
  marshaled["spgw_ipv4"] = stored.spgw_ipv4;
  marshaled["msisdn"] = stored.msisdn;
  marshaled["apn"] = stored.apn;
  marshaled["imei"] = stored.imei;
  marshaled["plmn_id"] = stored.plmn_id;
  marshaled["imsi_plmn_id"] = stored.imsi_plmn_id;
  marshaled["user_location"] = stored.user_location;
  marshaled["rat_type"] = static_cast<int>(stored.rat_type);
  marshaled["mac_addr"] = stored.mac_addr;
  marshaled["hardware_addr"] = stored.hardware_addr;
  marshaled["radius_session_id"] = stored.radius_session_id;
  marshaled["bearer_id"] = std::to_string(stored.bearer_id);
  marshaled["qos_info"] = serialize_stored_qos_info(stored.qos_info);

  std::string serialized = folly::toJson(marshaled);
  return serialized;
}

SessionConfig deserialize_stored_session_config(const std::string &serialized) {
  auto folly_serialized = folly::StringPiece(serialized);
  folly::dynamic marshaled = folly::parseJson(folly_serialized);

  auto stored = SessionConfig{};
  stored.ue_ipv4 = marshaled["ue_ipv4"].getString();
  stored.spgw_ipv4 = marshaled["spgw_ipv4"].getString();
  stored.msisdn = marshaled["msisdn"].getString();
  stored.apn = marshaled["apn"].getString();
  stored.imei = marshaled["imei"].getString();
  stored.plmn_id = marshaled["plmn_id"].getString();
  stored.imsi_plmn_id = marshaled["imsi_plmn_id"].getString();
  stored.user_location = marshaled["user_location"].getString();
  stored.rat_type = static_cast<RATType>(marshaled["rat_type"].getInt());
  stored.mac_addr = marshaled["mac_addr"].getString();
  stored.hardware_addr = marshaled["hardware_addr"].getString();
  stored.radius_session_id = marshaled["radius_session_id"].getString();
  stored.bearer_id =
      static_cast<uint32_t>(std::stoul(marshaled["bearer_id"].getString()));
  stored.qos_info =
      deserialize_stored_qos_info(marshaled["qos_info"].getString());

  return stored;
}

std::string serialize_stored_final_action_info(const FinalActionInfo &stored) {
  folly::dynamic marshaled = folly::dynamic::object;
  marshaled["final_action"] = static_cast<int>(stored.final_action);

  std::string redirect_server;
  stored.redirect_server.SerializeToString(&redirect_server);
  marshaled["redirect_server"] = redirect_server;

  std::string serialized = folly::toJson(marshaled);
  return serialized;
}

FinalActionInfo
deserialize_stored_final_action_info(const std::string &serialized) {
  auto folly_serialized = folly::StringPiece(serialized);
  folly::dynamic marshaled = folly::parseJson(folly_serialized);

  auto stored = FinalActionInfo{};
  stored.final_action = static_cast<ChargingCredit_FinalAction>(
      marshaled["final_action"].getInt());

  magma::lte::RedirectServer redirect_server;
  redirect_server.ParseFromString(marshaled["redirect_server"].getString());
  stored.redirect_server = redirect_server;

  return stored;
}

std::string serialize_stored_session_credit(StoredSessionCredit &stored) {
  folly::dynamic marshaled = folly::dynamic::object;
  marshaled["reporting"] = stored.reporting;
  marshaled["is_final"] = stored.is_final;
  marshaled["credit_limit_type"] = static_cast<int>(stored.credit_limit_type);
  marshaled["final_action_info"] =
      serialize_stored_final_action_info(stored.final_action_info);
  marshaled["reauth_state"] = static_cast<int>(stored.reauth_state);
  marshaled["service_state"] = static_cast<int>(stored.service_state);
  marshaled["expiry_time"] = std::to_string(stored.expiry_time);
  marshaled["buckets"] = folly::dynamic::object();
  marshaled["grant_tracking_type"] =
    static_cast<int>(stored.grant_tracking_type);
  for (int bucket_int = USED_TX; bucket_int != MAX_VALUES; bucket_int++) {
    Bucket bucket = static_cast<Bucket>(bucket_int);
    marshaled["buckets"][std::to_string(bucket_int)] =
        std::to_string(stored.buckets[bucket]);
  }

  std::string serialized = folly::toJson(marshaled);
  return serialized;
}

StoredSessionCredit
deserialize_stored_session_credit(const std::string &serialized) {
  auto folly_serialized = folly::StringPiece(serialized);
  folly::dynamic marshaled = folly::parseJson(folly_serialized);

  auto stored = StoredSessionCredit{};
  stored.reporting = marshaled["reporting"].getBool();
  stored.is_final = marshaled["is_final"].getBool();
  stored.credit_limit_type =
    static_cast<CreditLimitType>(marshaled["credit_limit_type"].getInt());
  stored.final_action_info = deserialize_stored_final_action_info(
      marshaled["final_action_info"].getString());
  stored.reauth_state =
      static_cast<ReAuthState>(marshaled["reauth_state"].getInt());
  stored.service_state =
      static_cast<ServiceState>(marshaled["service_state"].getInt());
  stored.expiry_time = static_cast<std::time_t>(
      std::stoul(marshaled["expiry_time"].getString()));
  stored.grant_tracking_type = static_cast<GrantTrackingType>(
    marshaled["grant_tracking_type"].getInt());
  for (int bucket_int = USED_TX; bucket_int != MAX_VALUES; bucket_int++) {
    Bucket bucket = static_cast<Bucket>(bucket_int);
    stored.buckets[bucket] = static_cast<uint64_t>(std::stoul(
        marshaled["buckets"][std::to_string(bucket_int)].getString()));
  }

  return stored;
}

std::string serialize_stored_monitor(StoredMonitor &stored) {
  folly::dynamic marshaled = folly::dynamic::object;

  marshaled["credit"] = serialize_stored_session_credit(stored.credit);
  marshaled["level"] = static_cast<int>(stored.level);

  std::string serialized = folly::toJson(marshaled);
  return serialized;
}

StoredMonitor deserialize_stored_monitor(const std::string &serialized) {
  auto folly_serialized = folly::StringPiece(serialized);
  folly::dynamic marshaled = folly::parseJson(folly_serialized);

  auto stored = StoredMonitor{};
  stored.credit =
      deserialize_stored_session_credit(marshaled["credit"].getString());
  stored.level = static_cast<MonitoringLevel>(marshaled["level"].getInt());

  return stored;
}

std::string
serialize_stored_charging_credit_map(StoredChargingCreditMap &stored) {
  folly::dynamic marshaled = folly::dynamic::object;

  folly::dynamic credit_keys = folly::dynamic::array;
  folly::dynamic credit_map = folly::dynamic::object;
  for (auto &credit_pair : stored) {
    CreditKey credit_key = credit_pair.first;
    folly::dynamic key2 = folly::dynamic::object;
    key2["rating_group"] = std::to_string(credit_key.rating_group);
    key2["service_identifier"] = std::to_string(credit_key.service_identifier);
    credit_keys.push_back(key2);

    std::string key = std::to_string(credit_key.rating_group) +
                      std::to_string(credit_key.service_identifier);
    credit_map[key] = serialize_stored_session_credit(credit_pair.second);
  }
  marshaled["credit_keys"] = credit_keys;
  marshaled["credit_map"] = credit_map;

  std::string serialized = folly::toJson(marshaled);
  return serialized;
}

StoredChargingCreditMap
deserialize_stored_charging_credit_map(std::string &serialized) {
  auto folly_serialized = folly::StringPiece(serialized);
  folly::dynamic marshaled = folly::parseJson(folly_serialized);

  auto stored = StoredChargingCreditMap(4, &ccHash, &ccEqual);

  for (auto &key : marshaled["credit_keys"]) {
    auto credit_key = CreditKey(
        static_cast<uint32_t>(std::stoul(key["rating_group"].getString())),
        static_cast<uint32_t>(
            std::stoul(key["service_identifier"].getString())));

    std::string key2 =
        key["rating_group"].getString() + key["service_identifier"].getString();

    stored[credit_key] = deserialize_stored_session_credit(
        marshaled["credit_map"][key2].getString());
  }
  return stored;
}

std::string serialize_stored_usage_monitor_map(
    StoredMonitorMap &stored) {
  folly::dynamic marshaled = folly::dynamic::object;

  folly::dynamic monitor_keys = folly::dynamic::array;
  folly::dynamic monitor_map = folly::dynamic::object;
  for (auto &monitor_pair : stored) {
    monitor_keys.push_back(monitor_pair.first);
    monitor_map[monitor_pair.first] =
        serialize_stored_monitor(monitor_pair.second);
  }
  marshaled["monitor_keys"] = monitor_keys;
  marshaled["monitor_map"] = monitor_map;

  std::string serialized = folly::toJson(marshaled);
  return serialized;
}

StoredMonitorMap
deserialize_stored_usage_monitor_map(std::string &serialized) {
  auto folly_serialized = folly::StringPiece(serialized);
  folly::dynamic marshaled = folly::parseJson(folly_serialized);
  auto stored = StoredMonitorMap{};
  for (auto &key : marshaled["monitor_keys"]) {
    std::string monitor_key = key.getString();
    stored[monitor_key] =
        deserialize_stored_monitor(marshaled["monitor_map"][key].getString());
  }

  return stored;
}

EventTriggerStatus deserialize_pending_event_triggers(std::string& serialized) {
  auto folly_serialized = folly::StringPiece(serialized);
  folly::dynamic marshaled = folly::parseJson(folly_serialized);

  auto stored = EventTriggerStatus{};
  for (auto &key : marshaled["event_trigger_keys"]) {
    auto map = marshaled["event_trigger_map"];
    magma::lte::EventTrigger eventKey;
    try {
      eventKey = magma::lte::EventTrigger(std::stoi(key.getString()));
    }
    catch (std::invalid_argument const &e){
      MLOG(MWARNING) << "Could not deserialize event triggers";
      continue;
    }
    stored[eventKey] = EventTriggerState(map[key].getInt());
  }

  return stored;
}

std::string serialize_pending_event_triggers(EventTriggerStatus event_triggers) {
  folly::dynamic marshaled = folly::dynamic::object;

  folly::dynamic keys = folly::dynamic::array;
  folly::dynamic map = folly::dynamic::object;
  for (auto &trigger_pair : event_triggers) {
    auto key = std::to_string(int(trigger_pair.first));
    keys.push_back(key);
    map[key] = int(trigger_pair.second);
  }
  marshaled["event_trigger_keys"] = keys;
  marshaled["event_trigger_map"] = map;

  std::string serialized = folly::toJson(marshaled);
  return serialized;
}

std::string serialize_stored_session(StoredSessionState &stored) {
  folly::dynamic marshaled = folly::dynamic::object;
  marshaled["fsm_state"] = static_cast<int>(stored.fsm_state);
  marshaled["config"] = serialize_stored_session_config(stored.config);
  marshaled["charging_pool"] =
      serialize_stored_charging_credit_map(stored.credit_map);
  marshaled["monitor_map"] =
      serialize_stored_usage_monitor_map(stored.monitor_map);
  marshaled["session_level_key"] = stored.session_level_key;
  marshaled["imsi"] = stored.imsi;
  marshaled["session_id"] = stored.session_id;
  marshaled["core_session_id"] = stored.core_session_id;
  marshaled["subscriber_quota_state"] =
      static_cast<int>(stored.subscriber_quota_state);

  std::string tgpp_context;
  stored.tgpp_context.SerializeToString(&tgpp_context);
  marshaled["tgpp_context"] = tgpp_context;

  marshaled["pending_event_triggers"] =
    serialize_pending_event_triggers(stored.pending_event_triggers);
  std::string revalidation_time;
  stored.revalidation_time.SerializeToString(&revalidation_time);
  marshaled["revalidation_time"] = revalidation_time;

  folly::dynamic static_rule_ids = folly::dynamic::array;
  for (const auto &rule_id : stored.static_rule_ids) {
    static_rule_ids.push_back(rule_id);
  }
  marshaled["static_rule_ids"] = static_rule_ids;

  folly::dynamic dynamic_rules = folly::dynamic::array;
  for (const auto &rule : stored.dynamic_rules) {
    std::string dynamic_rule;
    rule.SerializeToString(&dynamic_rule);
    dynamic_rules.push_back(dynamic_rule);
  }
  marshaled["dynamic_rules"] = dynamic_rules;

  folly::dynamic gy_dynamic_rules = folly::dynamic::array;
  for (const auto &rule : stored.gy_dynamic_rules) {
    std::string gy_dynamic_rule;
    rule.SerializeToString(&gy_dynamic_rule);
    gy_dynamic_rules.push_back(gy_dynamic_rule);
  }
  marshaled["gy_dynamic_rules"] = gy_dynamic_rules;

  marshaled["request_number"] = std::to_string(stored.request_number);

  std::string serialized = folly::toJson(marshaled);
  return serialized;
}

StoredSessionState deserialize_stored_session(std::string &serialized) {
  auto folly_serialized = folly::StringPiece(serialized);
  folly::dynamic marshaled = folly::parseJson(folly_serialized);

  auto stored = StoredSessionState{};
  stored.fsm_state = SessionFsmState(marshaled["fsm_state"].getInt());
  stored.config =
      deserialize_stored_session_config(marshaled["config"].getString());
  stored.credit_map = deserialize_stored_charging_credit_map(
      marshaled["charging_pool"].getString());
  stored.monitor_map = deserialize_stored_usage_monitor_map(
      marshaled["monitor_map"].getString());
  stored.session_level_key = marshaled["session_level_key"].getString();
  stored.imsi = marshaled["imsi"].getString();
  stored.session_id = marshaled["session_id"].getString();
  stored.core_session_id = marshaled["core_session_id"].getString();
  stored.subscriber_quota_state =
      static_cast<magma::lte::SubscriberQuotaUpdate_Type>(
          marshaled["subscriber_quota_state"].getInt());

  google::protobuf::Timestamp revalidation_time;
  revalidation_time.ParseFromString(marshaled["revalidation_time"].getString());
  stored.revalidation_time = revalidation_time;
  stored.pending_event_triggers =
    deserialize_pending_event_triggers(marshaled["pending_event_triggers"].getString());

  magma::lte::TgppContext tgpp_context;
  tgpp_context.ParseFromString(marshaled["tgpp_context"].getString());
  stored.tgpp_context = tgpp_context;

  for (auto &rule_id : marshaled["static_rule_ids"]) {
    stored.static_rule_ids.push_back(rule_id.getString());
  }

  for (auto &policy : marshaled["dynamic_rules"]) {
    PolicyRule policy_rule;
    policy_rule.ParseFromString(policy.getString());
    stored.dynamic_rules.push_back(policy_rule);
  }

  for (auto &policy : marshaled["gy_dynamic_rules"]) {
    PolicyRule policy_rule;
    policy_rule.ParseFromString(policy.getString());
    stored.gy_dynamic_rules.push_back(policy_rule);
  }

  stored.request_number = static_cast<uint32_t>(
      std::stoul(marshaled["request_number"].getString()));

  return stored;
}

}; // namespace magma

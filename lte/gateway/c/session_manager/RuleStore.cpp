/**
 * Copyright (c) 2016-present, Facebook, Inc.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */
#include <glog/logging.h>

#include "RuleStore.h"
#include "ServiceRegistrySingleton.h"

using grpc::Status;

namespace magma {

template<typename KeyType, typename hash, typename equal>
void PoliciesByKeyMap<KeyType, hash, equal>::insert(
  const KeyType& key,
  std::shared_ptr<PolicyRule> rule_p)
{
  auto iter = rules_by_key_.find(key);
  if (iter == rules_by_key_.end()) {
    rules_by_key_[key] = {rule_p};
    return;
  }
  iter->second.push_back(rule_p);
}

template<typename KeyType, typename hash, typename equal>
void PoliciesByKeyMap<KeyType, hash, equal>::remove(
  const KeyType& key,
  std::shared_ptr<PolicyRule> rule_p)
{
  auto iter = rules_by_key_.find(key);
  if (iter == rules_by_key_.end()) {
    return;
  }

  auto rules = iter->second;
  auto found = std::find(rules.begin(), rules.end(), rule_p);
  if (found == rules.end()) {
    return;
  }
  rules.erase(found);
}

template<typename KeyType, typename hash, typename equal>
bool PoliciesByKeyMap<KeyType, hash, equal>::get_rule_ids_for_key(
  const KeyType& key,
  std::vector<std::string>& rules_out)
{
  auto iter = rules_by_key_.find(key);
  if (iter == rules_by_key_.end()) {
    return false;
  }

  for (const auto& rule : iter->second) {
    rules_out.push_back(rule->id());
  }
  return true;
}

template<typename KeyType, typename hash, typename equal>
bool PoliciesByKeyMap<KeyType, hash, equal>::get_rule_definitions_for_key(
  const KeyType& key,
  std::vector<PolicyRule>& rules_out)
{
  auto iter = rules_by_key_.find(key);
  if (iter == rules_by_key_.end()) {
    return false;
  }

  for (const auto& rule : iter->second) {
    rules_out.push_back(*rule);
  }
  return true;
}

static bool should_track_charging_key(PolicyRule::TrackingType tracking_type)
{
  return tracking_type == PolicyRule::ONLY_OCS ||
         tracking_type == PolicyRule::OCS_AND_PCRF;
}

static bool should_track_monitoring_key(PolicyRule::TrackingType tracking_type)
{
  return tracking_type == PolicyRule::ONLY_PCRF ||
         tracking_type == PolicyRule::OCS_AND_PCRF;
}

void PolicyRuleBiMap::sync_rules(const std::vector<PolicyRule>& rules)
{
  std::lock_guard<std::mutex> lock(map_mutex_);
  rules_by_rule_id_.clear();
  rules_by_charging_key_ =
    PoliciesByKeyMap<CreditKey, decltype(&ccHash), decltype(&ccEqual)>(
      &ccHash, &ccEqual);
  rules_by_monitoring_key_ = PoliciesByKeyMap<std::string>();
  for (const auto& rule : rules) {
    auto rule_p = std::make_shared<PolicyRule>(rule);
    rules_by_rule_id_[rule.id()] = rule_p;
    if (should_track_charging_key(rule.tracking_type())) {
      rules_by_charging_key_.insert(CreditKey(rule), rule_p);
    }
    if (should_track_monitoring_key(rule.tracking_type())) {
      rules_by_monitoring_key_.insert(rule.monitoring_key(), rule_p);
    }
  }
}

void PolicyRuleBiMap::insert_rule(const PolicyRule& rule)
{
  auto rule_p = std::make_shared<PolicyRule>(rule);
  std::lock_guard<std::mutex> lock(map_mutex_);
  rules_by_rule_id_[rule.id()] = rule_p;
  if (should_track_charging_key(rule.tracking_type())) {
    rules_by_charging_key_.insert(CreditKey(rule), rule_p);
  }
  if (should_track_monitoring_key(rule.tracking_type())) {
    rules_by_monitoring_key_.insert(rule.monitoring_key(), rule_p);
  }
}

bool PolicyRuleBiMap::get_rule(const std::string& rule_id, PolicyRule* rule)
{
  std::lock_guard<std::mutex> lock(map_mutex_);
  auto it = rules_by_rule_id_.find(rule_id);
  if (it == rules_by_rule_id_.end()) {
    return false;
  }
  rule->CopyFrom(*it->second);
  return true;
}

bool PolicyRuleBiMap::remove_rule(
  const std::string& rule_id,
  PolicyRule* rule_out)
{
  std::lock_guard<std::mutex> lock(map_mutex_);
  auto it = rules_by_rule_id_.find(rule_id);
  if (it == rules_by_rule_id_.end()) {
    return false;
  }

  auto rule_ptr = it->second;
  rule_out->CopyFrom(*rule_ptr);

  // Remove the rule from all mappings
  rules_by_rule_id_.erase(it);
  if (should_track_charging_key(rule_ptr->tracking_type())) {
    rules_by_charging_key_.remove(CreditKey(rule_ptr.get()), rule_ptr);
  }
  if (should_track_monitoring_key(rule_ptr->tracking_type())) {
    rules_by_monitoring_key_.remove(rule_ptr->monitoring_key(), rule_ptr);
  }

  return true;
}

bool PolicyRuleBiMap::get_charging_key_for_rule_id(
  const std::string& rule_id,
  CreditKey* charging_key)
{
  std::lock_guard<std::mutex> lock(map_mutex_);
  auto it = rules_by_rule_id_.find(rule_id);
  if (it == rules_by_rule_id_.end()) {
    return false;
  }
  if (should_track_charging_key(it->second->tracking_type())) {
    charging_key->set(it->second.get());
    return true;
  }
  return false;
}

bool PolicyRuleBiMap::get_monitoring_key_for_rule_id(
  const std::string& rule_id,
  std::string* monitoring_key)
{
  std::lock_guard<std::mutex> lock(map_mutex_);
  auto it = rules_by_rule_id_.find(rule_id);
  if (it == rules_by_rule_id_.end()) {
    return false;
  }
  if (should_track_monitoring_key(it->second->tracking_type())) {
    monitoring_key->assign(it->second->monitoring_key());
    return true;
  }
  return false;
}

bool PolicyRuleBiMap::get_rule_ids_for_charging_key(
  const CreditKey& charging_key,
  std::vector<std::string>& rules_out)
{
  std::lock_guard<std::mutex> lock(map_mutex_);
  bool success =
    rules_by_charging_key_.get_rule_ids_for_key(charging_key, rules_out);
  return success;
}

bool PolicyRuleBiMap::get_rule_definitions_for_charging_key(
  const CreditKey& charging_key,
  std::vector<PolicyRule>& rules_out)
{
  std::lock_guard<std::mutex> lock(map_mutex_);
  bool success = rules_by_charging_key_.get_rule_definitions_for_key(
    charging_key, rules_out);
  return success;
}

bool PolicyRuleBiMap::get_rule_ids_for_monitoring_key(
  const std::string& monitoring_key,
  std::vector<std::string>& rules_out)
{
  std::lock_guard<std::mutex> lock(map_mutex_);
  bool success =
    rules_by_monitoring_key_.get_rule_ids_for_key(monitoring_key, rules_out);
  return success;
}

bool PolicyRuleBiMap::get_rule_definitions_for_monitoring_key(
  const std::string& monitoring_key,
  std::vector<PolicyRule>& rules_out)
{
  std::lock_guard<std::mutex> lock(map_mutex_);
  bool success = rules_by_monitoring_key_.get_rule_definitions_for_key(
    monitoring_key, rules_out);
  return success;
}

bool PolicyRuleBiMap::get_rule_ids(
  std::vector<std::string>& rules_ids_out)
{
  std::lock_guard<std::mutex> lock(map_mutex_);
  for(auto kv : rules_by_rule_id_) {
    rules_ids_out.push_back(kv.first);
  }
  return true;
}

bool PolicyRuleBiMap::get_rules(
  std::vector<PolicyRule>& rules_out)
{
  std::lock_guard<std::mutex> lock(map_mutex_);
  for(auto kv : rules_by_rule_id_) {
    rules_out.push_back(*kv.second);
  }
  return true;
}

} // namespace magma

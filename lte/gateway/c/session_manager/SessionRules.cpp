/**
 * Copyright (c) 2016-present, Facebook, Inc.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */
#include "SessionRules.h"
namespace magma {

std::unique_ptr<SessionRules> SessionRules::unmarshal(
  const StoredSessionRules& marshaled,
  StaticRuleStore& static_rule_ref)
{
  return std::make_unique<SessionRules>(marshaled, static_rule_ref);
}

SessionRules::SessionRules(
  const StoredSessionRules& marshaled,
  StaticRuleStore& static_rule_ref):
  static_rules_(static_rule_ref)
{
  for (const std::string& rule_id : marshaled.static_rule_ids)
  {
    active_static_rules_.push_back(rule_id);
  }
  for (auto& rule : marshaled.dynamic_rules)
  {
    dynamic_rules_.insert_rule(rule);
  }
}

StoredSessionRules SessionRules::marshal()
{
  StoredSessionRules stored_rules;
  for (auto& rule_id : active_static_rules_)
  {
    stored_rules.static_rule_ids.push_back(rule_id);
  }
  std::vector<PolicyRule> dynamic_rules;
  dynamic_rules_.get_rules(dynamic_rules);
  stored_rules.dynamic_rules = std::move(dynamic_rules);
  return stored_rules;
}

SessionRules::SessionRules(StaticRuleStore& static_rule_ref):
  static_rules_(static_rule_ref)
{
}

bool SessionRules::get_charging_key_for_rule_id(
  const std::string& rule_id,
  CreditKey* charging_key)
{
  // first check dynamic rules and then static rules
  if (dynamic_rules_.get_charging_key_for_rule_id(rule_id, charging_key)) {
    return true;
  }
  if (static_rules_.get_charging_key_for_rule_id(rule_id, charging_key)) {
    return true;
  }
  return false;
}

bool SessionRules::get_monitoring_key_for_rule_id(
  const std::string& rule_id,
  std::string* monitoring_key)
{
  // first check dynamic rules and then static rules
  if (dynamic_rules_.get_monitoring_key_for_rule_id(rule_id, monitoring_key)) {
    return true;
  }
  if (static_rules_.get_monitoring_key_for_rule_id(rule_id, monitoring_key)) {
    return true;
  }
  return false;
}

bool SessionRules::is_dynamic_rule_installed(const std::string& rule_id)
{
  auto _ = new PolicyRule();
  return dynamic_rules_.get_rule(rule_id, _);
}

bool SessionRules::is_static_rule_installed(const std::string& rule_id)
{
  return std::find(
    active_static_rules_.begin(),
    active_static_rules_.end(),
    rule_id) != active_static_rules_.end();
}

void SessionRules::insert_dynamic_rule(const PolicyRule& rule)
{
  dynamic_rules_.insert_rule(rule);
}

void SessionRules::activate_static_rule(const std::string& rule_id)
{
  active_static_rules_.push_back(rule_id);
}

bool SessionRules::remove_dynamic_rule(
  const std::string& rule_id,
  PolicyRule *rule_out)
{
  return dynamic_rules_.remove_rule(rule_id, rule_out);
}

bool SessionRules::deactivate_static_rule(const std::string& rule_id)
{
  auto it = std::find(active_static_rules_.begin(), active_static_rules_.end(),
                      rule_id);
  if (it == active_static_rules_.end()) {
    return false;
  }
  active_static_rules_.erase(it);
  return true;
}
/**
 * For the charging key, get any applicable rules from the static rule set
 * and the dynamic rule set
 */
void SessionRules::add_rules_to_action(
  ServiceAction& action,
  const CreditKey& charging_key)
{
  static_rules_.get_rule_ids_for_charging_key(
    charging_key, *action.get_mutable_rule_ids());
  dynamic_rules_.get_rule_definitions_for_charging_key(
    charging_key, *action.get_mutable_rule_definitions());
}

void SessionRules::add_rules_to_action(
  ServiceAction& action,
  std::string monitoring_key)
{
  static_rules_.get_rule_ids_for_monitoring_key(
    monitoring_key, *action.get_mutable_rule_ids());
  dynamic_rules_.get_rule_definitions_for_monitoring_key(
    monitoring_key, *action.get_mutable_rule_definitions());
}

std::vector<std::string>& SessionRules::get_static_rule_ids()
{
  return active_static_rules_;
}

DynamicRuleStore& SessionRules::get_dynamic_rules()
{
  return dynamic_rules_;
}

uint32_t SessionRules::total_monitored_rules_count()
{
  uint32_t monitored_dynamic_rules = dynamic_rules_.monitored_rules_count();
  uint32_t monitored_static_rules = 0;
  for (auto& rule_id : active_static_rules_)
  {
    std::string mkey; // ignore value
    auto is_monitored = static_rules_.get_monitoring_key_for_rule_id(
      rule_id,& mkey);
    if (is_monitored) {
      monitored_static_rules++;
    }
  }
  return monitored_dynamic_rules + monitored_static_rules;
}

} // namespace magma

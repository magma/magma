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
#include "ServiceAction.h"
#include "StoredState.h"

namespace magma {

/**
 * SessionRules maintains the dynamic and static rules for a subscriber session
 */
class SessionRules {
 public:
  static std::unique_ptr<SessionRules> unmarshal(
    const StoredSessionRules& marshaled,
    StaticRuleStore &static_rule_ref);

  StoredSessionRules marshal();

  SessionRules(StaticRuleStore &static_rule_ref);

  SessionRules(
    const StoredSessionRules& marshaled,
    StaticRuleStore &static_rule_ref);

  bool get_charging_key_for_rule_id(
    const std::string &rule_id,
    CreditKey *charging_key);

  bool get_monitoring_key_for_rule_id(
    const std::string &rule_id,
    std::string *monitoring_key);

  void insert_dynamic_rule(const PolicyRule &rule);

  void activate_static_rule(const std::string &rule_id);

  bool remove_dynamic_rule(const std::string &rule_id, PolicyRule *rule_out);

  bool deactivate_static_rule(const std::string &rule_id);

  void add_rules_to_action(ServiceAction &action, const CreditKey &charging_key);
  void add_rules_to_action(ServiceAction &action, std::string monitoring_key);

  std::vector<std::string> &get_static_rule_ids();
  DynamicRuleStore &get_dynamic_rules();

 private:
  StaticRuleStore &static_rules_;
  std::vector<std::string> active_static_rules_;
  DynamicRuleStore dynamic_rules_;
};

} // namespace magma

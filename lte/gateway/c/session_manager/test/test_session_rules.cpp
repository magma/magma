/**
 * Copyright (c) 2016-present, Facebook, Inc.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */
#include <memory>
#include <future>

#include <glog/logging.h>
#include <gtest/gtest.h>

#include "SessionRules.h"
#include "magma_logging.h"

using ::testing::Test;

namespace magma {

class SessionRulesTest : public ::testing::Test {
 protected:
 protected:
  virtual void SetUp()
  {
    rule_store = std::make_shared<StaticRuleStore>();
    session_rules = std::make_shared<SessionRules>(*rule_store);
  }

  PolicyRule get_rule(
    uint32_t rating_group,
    const std::string &m_key,
    const std::string &rule_id)
  {
    PolicyRule rule;
    rule.set_id(rule_id);
    rule.set_rating_group(rating_group);
    rule.set_monitoring_key(m_key);
    if (rating_group == 0 && m_key.length() > 0) {
      rule.set_tracking_type(PolicyRule::ONLY_PCRF);
    } else if (rating_group > 0 && m_key.length() == 0) {
      rule.set_tracking_type(PolicyRule::ONLY_OCS);
    } else if (rating_group > 0 && m_key.length() > 0) {
      rule.set_tracking_type(PolicyRule::OCS_AND_PCRF);
    } else {
      rule.set_tracking_type(PolicyRule::NO_TRACKING);
    }
    return rule;
  }

  void activate_rule(
    uint32_t rating_group,
    const std::string &m_key,
    const std::string &rule_id,
    bool is_static)
  {
    PolicyRule rule = get_rule(rating_group, m_key, rule_id);
    if (is_static) {
      rule_store->insert_rule(rule);
      session_rules->activate_static_rule(rule_id);
    } else {
      session_rules->insert_dynamic_rule(rule);
    }
  }

 protected:
  std::shared_ptr<StaticRuleStore> rule_store;
  std::shared_ptr<SessionRules> session_rules;
};

TEST_F(SessionRulesTest, test_marshal_unmarshal)
{
  // Activate a dynamic rule
  activate_rule(1, "m1", "rule1", false);

  // Activate static rules
  activate_rule(2, "m2", "rule2", true);

  std::vector<std::string> rules_out{};
  std::vector<std::string>& rules_out_ptr = rules_out;

  session_rules->get_dynamic_rules().get_rule_ids(rules_out_ptr);
  EXPECT_EQ(rules_out_ptr.size(), 1);
  EXPECT_EQ(rules_out_ptr[0], "rule1");

  auto static_rules = session_rules->get_static_rule_ids();
  EXPECT_EQ(static_rules.size(), 1);
  EXPECT_EQ(static_rules[0], "rule2");

  // Check that after marshaling/unmarshaling that the fields are still the
  // same.
  auto marshaled = (*session_rules).marshal();
  auto session_rules_2 = SessionRules::unmarshal(marshaled, *rule_store);

  rules_out = {};
  session_rules_2->get_dynamic_rules().get_rule_ids(rules_out_ptr);
  EXPECT_EQ(rules_out_ptr.size(), 1);
  EXPECT_EQ(rules_out_ptr[0], "rule1");

  static_rules = session_rules_2->get_static_rule_ids();
  EXPECT_EQ(static_rules.size(), 1);
  EXPECT_EQ(static_rules[0], "rule2");
}

int main(int argc, char **argv)
{
  ::testing::InitGoogleTest(&argc, argv);
  FLAGS_logtostderr = 1;
  FLAGS_v = 10;
  return RUN_ALL_TESTS();
}

} // namespace magma

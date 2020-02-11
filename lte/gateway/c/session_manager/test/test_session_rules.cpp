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

  enum RuleType {
    STATIC = 0,
    DYNAMIC = 1,
  };

  // TODO take these into a test common file
  void activate_rule(
    uint32_t rating_group,
    const std::string &m_key,
    const std::string &rule_id,
    RuleType rule_type)
  {
    PolicyRule rule = get_rule(rating_group, m_key, rule_id);
    switch (rule_type) {
      case STATIC:
        rule_store->insert_rule(rule);
        session_rules->activate_static_rule(rule_id);
        break;
      case DYNAMIC:
        session_rules->insert_dynamic_rule(rule);
        break;
    }
  }

 protected:
  std::shared_ptr<StaticRuleStore> rule_store;
  std::shared_ptr<SessionRules> session_rules;
};

TEST_F(SessionRulesTest, test_marshal_unmarshal)
{
  activate_rule(1, "m1", "rule1", DYNAMIC);
  activate_rule(2, "m2", "rule2", STATIC);

  std::vector<std::string> rules_out{};
  std::vector<std::string>& rules_out_ptr = rules_out;

  session_rules->get_dynamic_rules().get_rule_ids(rules_out_ptr);
  EXPECT_EQ(rules_out_ptr.size(), 1);
  EXPECT_EQ(rules_out_ptr[0], "rule1");

  auto static_rules = session_rules->get_static_rule_ids();
  EXPECT_EQ(static_rules.size(), 1);
  EXPECT_EQ(static_rules[0], "rule2");

  // Check that after marshaling/un-marshaling that the fields are still the
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

TEST_F(SessionRulesTest, test_insert_remove)
{
  activate_rule(1, "m1", "rule1", DYNAMIC);
  EXPECT_EQ(1, session_rules->total_monitored_rules_count());
  activate_rule(2, "m2", "rule2", STATIC);
  EXPECT_EQ(2, session_rules->total_monitored_rules_count());
  // add a OCS-ONLY static rule
  activate_rule(3, "", "rule3", STATIC);
  EXPECT_EQ(2, session_rules->total_monitored_rules_count());

  std::vector<std::string> rules_out{};
  std::vector<std::string>& rules_out_ptr = rules_out;

  session_rules->get_dynamic_rules().get_rule_ids(rules_out_ptr);
  EXPECT_EQ(rules_out_ptr.size(), 1);
  EXPECT_EQ(rules_out_ptr[0], "rule1");

  auto static_rules = session_rules->get_static_rule_ids();
  EXPECT_EQ(static_rules.size(), 2);
  EXPECT_EQ(static_rules[0], "rule2");
  EXPECT_EQ(static_rules[1], "rule3");

  // Test rule removals
  PolicyRule rule_out;
  session_rules->deactivate_static_rule("rule2");
  EXPECT_EQ(1, session_rules->total_monitored_rules_count());
  EXPECT_EQ(true, session_rules->remove_dynamic_rule("rule1", &rule_out));
  EXPECT_EQ("m1", rule_out.monitoring_key());
  EXPECT_EQ(0, session_rules->total_monitored_rules_count());

  // basic sanity checks to see it's properly deleted
  rules_out = {};
  session_rules->get_dynamic_rules().get_rule_ids(rules_out_ptr);
  EXPECT_EQ(rules_out_ptr.size(), 0);

  rules_out = {};
  session_rules->get_dynamic_rules()
    .get_rule_ids_for_monitoring_key("m1", rules_out);
  EXPECT_EQ(0, rules_out.size());

  std::string mkey;
  // searching for non-existent rule should fail
  EXPECT_EQ(false, session_rules->get_dynamic_rules()
    .get_monitoring_key_for_rule_id("rule1", &mkey));
  // deleting an already deleted rule should fail
  EXPECT_EQ(false, session_rules->get_dynamic_rules()
    .remove_rule("rule1", &rule_out));
}

int main(int argc, char **argv)
{
  ::testing::InitGoogleTest(&argc, argv);
  FLAGS_logtostderr = 1;
  FLAGS_v = 10;
  return RUN_ALL_TESTS();
}

} // namespace magma

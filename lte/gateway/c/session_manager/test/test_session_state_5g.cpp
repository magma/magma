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
#include <glog/logging.h>
#include <gtest/gtest.h>

#include <future>
#include <memory>
#include <utility>

#include "Consts.h"
#include "magma_logging.h"
#include "ProtobufCreators.h"
#include "SessiondMocks.h"
#include "SessionState.h"
#include "SessionStateTester5g.h"

using ::testing::Test;

namespace magma {

TEST_F(SessionStateTest5G, test_get_session_rules) {
  session_state->set_config(cfg, nullptr);
  EXPECT_TRUE(session_state->is_5g_session());
  // populate rule store with 2 static and 2 dynamic rules
  activate_rule(1, "", "rule-static-1", STATIC, 0, 0);
  activate_rule(2, "m1", "rule-static-2", STATIC, 0, 0);
  activate_rule(1, "", "rule-dynamic-1", DYNAMIC, 0, 0);
  activate_rule(2, "m1", "rule-dynamic-2", DYNAMIC, 0, 0);

  EXPECT_TRUE(session_state->is_static_rule_installed("rule-static-1"));
  EXPECT_TRUE(session_state->is_static_rule_installed("rule-static-2"));
  EXPECT_TRUE(session_state->is_dynamic_rule_installed("rule-dynamic-1"));
  EXPECT_TRUE(session_state->is_dynamic_rule_installed("rule-dynamic-2"));

  insert_static_rule_into_store(3, "m2", "rule-static-3");

  PolicyRule dynamic_2 = create_policy_rule("rule-dynamic-2", "m1", 2);
  PolicyRule dynamic_3 = create_policy_rule("rule-dynamic-3", "m1", 3);

  // Create a StaticRuleInstall with rules above
  std::vector<StaticRuleInstall> static_rule_installs{
      create_static_rule_install("rule-static-2"),
      create_static_rule_install("rule-static-3"),
  };

  // Create a DynamicRuleInstall
  std::vector<DynamicRuleInstall> dynamic_rule_installs{
      create_dynamic_rule_install(dynamic_2),
      create_dynamic_rule_install(dynamic_3),
  };

  RulesToProcess pending_activation, pending_deactivation;
  session_state->process_get_5g_rule_installs(
      static_rule_installs, dynamic_rule_installs, &pending_activation,
      &pending_deactivation);

  // First check the active rules in session
  EXPECT_TRUE(session_state->is_static_rule_installed("rule-static-1"));
  EXPECT_TRUE(session_state->is_static_rule_installed("rule-static-2"));
  EXPECT_TRUE(!session_state->is_static_rule_installed("rule-static-3"));
  EXPECT_TRUE(session_state->is_dynamic_rule_installed("rule-dynamic-1"));
  EXPECT_TRUE(session_state->is_dynamic_rule_installed("rule-dynamic-2"));
  EXPECT_TRUE(!session_state->is_dynamic_rule_installed("rule-dynamic-3"));

  // Check the RulesToProcess is properly filled out
  EXPECT_EQ(pending_activation.size(), 4);
  const std::string activate_rule1 = pending_activation[0].rule.id();
  const std::string activate_rule2 = pending_activation[1].rule.id();
  const std::string deactivate_rule1 = pending_deactivation[0].rule.id();
  const std::string deactivate_rule2 = pending_deactivation[1].rule.id();
  EXPECT_TRUE(activate_rule1 == "rule-static-2" ||
              activate_rule1 == "rule-dynamic-2");
  EXPECT_TRUE(activate_rule2 == "rule-static-3" ||
              activate_rule2 == "rule-dynamic-3");
  EXPECT_TRUE(deactivate_rule1 == "rule-static-2" ||
              deactivate_rule1 == "rule-dynamic-2");
  EXPECT_TRUE(deactivate_rule2 == "rule-static-3" ||
              deactivate_rule2 == "rule-dynamic-3");
}

TEST_F(SessionStateTest5G, test_process_static_rule_installs) {
  session_state->set_config(cfg, nullptr);
  // Insert 2 static rules without qos into static rule store
  insert_static_rule_into_store(0, "mkey1", "static-1");
  insert_static_rule_into_store(0, "mkey1", "static-2");
  // Insert 2 static rules with qos into static rule store
  insert_static_rule_with_qos_into_store(0, "mkey1", 1, "static-qos-3");
  insert_static_rule_with_qos_into_store(0, "mkey1", 2, "static-qos-4");

  // activate static-1 and static-qos-3 in advance
  RuleLifetime lifetime;
  session_state->activate_static_rule("static-1", lifetime, nullptr);
  session_state->activate_static_rule("static-qos-3", lifetime, nullptr);

  // Create a StaticRuleInstall with all four rules above
  std::vector<StaticRuleInstall> rule_installs{
      // should be ignored as it is already active
      create_static_rule_install("static-1"),
      // new non-qos rule
      create_static_rule_install("static-2"),
      // should be ignored as it is already active
      create_static_rule_install("static-qos-3"),
      // new qos rule
      create_static_rule_install("static-qos-4"),
  };
  EXPECT_EQ(4, rule_installs.size());
  RulesToProcess pending_activation, pending_deactivation, pending_bearer_setup;
  RulesToSchedule pending_scheduling;
  // For 5g verification
  EXPECT_TRUE(session_state->is_5g_session());
  session_state->process_static_rule_installs(
      rule_installs, &pending_activation, &pending_deactivation,
      &pending_bearer_setup, &pending_scheduling, &update_criteria);
  EXPECT_EQ(2, pending_activation.size());
  EXPECT_EQ("static-2", pending_activation[0].rule.id());
  EXPECT_EQ(1, pending_activation[0].version);
  EXPECT_EQ("static-qos-4", pending_activation[1].rule.id());
  EXPECT_EQ(1, pending_activation[1].version);
  // For 5g verification
  EXPECT_EQ(0, pending_activation[0].teids.agw_teid());
  EXPECT_EQ(0, pending_activation[0].teids.enb_teid());

  EXPECT_EQ(update_criteria.static_rules_to_install.size(), 2);
  EXPECT_TRUE(update_criteria.static_rules_to_install.count("static-2"));
  EXPECT_TRUE(update_criteria.static_rules_to_install.count("static-qos-4"));

  // For 5g verification
  EXPECT_EQ(0, pending_bearer_setup.size());
}

TEST_F(SessionStateTest5G, test_process_dynamic_rule_installs) {
  session_state->set_config(cfg, nullptr);
  PolicyRule dynamic_1 = create_policy_rule("dynamic-1", "", 0);
  PolicyRule dynamic_2 = create_policy_rule("dynamic-2", "", 0);
  PolicyRule dynamic_qos_3 =
      create_policy_rule_with_qos("dynamic-qos-3", "", 0, 1);
  PolicyRule dynamic_qos_4 =
      create_policy_rule_with_qos("dynamic-qos-4", "", 0, 2);

  // Install dynamic rules for dynamic-1 and dynamic-qos-3
  // Then version should be increased.
  RuleLifetime lifetime;
  session_state->insert_dynamic_rule(dynamic_1, lifetime, nullptr);
  session_state->insert_dynamic_rule(dynamic_qos_3, lifetime, nullptr);

  // Create a StaticRuleInstall with all four rules above
  std::vector<DynamicRuleInstall> rule_installs{
      // should be installed even though it is already active
      create_dynamic_rule_install(dynamic_1),
      // new non-qos rule
      create_dynamic_rule_install(dynamic_2),
      // should be installed even though it is already active
      create_dynamic_rule_install(dynamic_qos_3),
      // new qos rule
      create_dynamic_rule_install(dynamic_qos_4),
  };
  EXPECT_EQ(4, rule_installs.size());
  RulesToProcess pending_activation, pending_deactivation, pending_bearer_setup;
  RulesToSchedule pending_scheduling;
  // For 5g verification
  EXPECT_TRUE(session_state->is_5g_session());
  session_state->process_dynamic_rule_installs(
      rule_installs, &pending_activation, &pending_deactivation,
      &pending_bearer_setup, &pending_scheduling, &update_criteria);
  EXPECT_EQ(4, pending_activation.size());
  EXPECT_EQ("dynamic-1", pending_activation[0].rule.id());
  EXPECT_EQ(2, pending_activation[0].version);
  EXPECT_EQ("dynamic-2", pending_activation[1].rule.id());
  EXPECT_EQ(1, pending_activation[1].version);
  EXPECT_EQ("dynamic-qos-3", pending_activation[2].rule.id());
  EXPECT_EQ(2, pending_activation[2].version);
  EXPECT_EQ("dynamic-qos-4", pending_activation[3].rule.id());
  EXPECT_EQ(1, pending_activation[3].version);

  // For 5g verification
  EXPECT_EQ(0, pending_activation[0].teids.agw_teid());
  EXPECT_EQ(0, pending_activation[0].teids.enb_teid());

  EXPECT_EQ(update_criteria.dynamic_rules_to_install.size(), 4);
  EXPECT_EQ("dynamic-1", update_criteria.dynamic_rules_to_install[0].id());
  EXPECT_EQ("dynamic-2", update_criteria.dynamic_rules_to_install[1].id());
  EXPECT_EQ("dynamic-qos-3", update_criteria.dynamic_rules_to_install[2].id());
  EXPECT_EQ("dynamic-qos-4", update_criteria.dynamic_rules_to_install[3].id());

  // For 5g verification
  EXPECT_EQ(0, pending_bearer_setup.size());
}

TEST_F(SessionStateTest5G, test_remove_all_session_rules) {
  session_state->set_config(cfg, nullptr);
  // Insert 1 static rules without qos into static rule store
  insert_static_rule_into_store(0, "mkey1", "static-1");
  // Insert 1 static rules with qos into static rule store
  insert_static_rule_with_qos_into_store(0, "mkey1", 1, "static-2");

  // Create a StaticRuleInstall all rules
  std::vector<StaticRuleInstall> static_rule_installs{
      create_static_rule_install("static-1"),
      create_static_rule_install("static-2"),
  };
  EXPECT_EQ(2, static_rule_installs.size());

  PolicyRule dynamic_1 = create_policy_rule("dynamic-1", "", 0);
  PolicyRule dynamic_2 = create_policy_rule_with_qos("dynamic-2", "", 0, 1);

  // Create a DynamicRuleInstall rules
  std::vector<DynamicRuleInstall> dynamic_rule_installs{
      create_dynamic_rule_install(dynamic_1),
      create_dynamic_rule_install(dynamic_2),
  };
  EXPECT_EQ(2, dynamic_rule_installs.size());

  RulesToProcess pending_activation, pending_deactivation, pending_bearer_setup;
  RulesToSchedule pending_scheduling;
  // For 5g verification
  EXPECT_TRUE(session_state->is_5g_session());
  session_state->process_rules_to_install(
      static_rule_installs, dynamic_rule_installs, &pending_activation,
      &pending_deactivation, &pending_bearer_setup, &pending_scheduling,
      &update_criteria);

  EXPECT_EQ(4, pending_activation.size());

  const std::string activate_rule1 = pending_activation[0].rule.id();
  const std::string activate_rule2 = pending_activation[1].rule.id();

  EXPECT_TRUE(activate_rule1 == "static-1" || activate_rule1 == "dynamic-1");
  EXPECT_TRUE(activate_rule2 == "static-2" || activate_rule2 == "dynamic-2");

  EXPECT_EQ(update_criteria.static_rules_to_install.size(), 2);
  EXPECT_TRUE(update_criteria.static_rules_to_install.count("static-1"));
  EXPECT_TRUE(update_criteria.static_rules_to_install.count("static-2"));

  EXPECT_TRUE(session_state->is_static_rule_installed("static-1"));
  EXPECT_TRUE(session_state->is_static_rule_installed("static-2"));

  EXPECT_EQ(update_criteria.dynamic_rules_to_install.size(), 2);
  EXPECT_EQ("dynamic-1", update_criteria.dynamic_rules_to_install[0].id());
  EXPECT_EQ("dynamic-2", update_criteria.dynamic_rules_to_install[1].id());

  // For 5g verification
  EXPECT_EQ(0, pending_activation[0].teids.agw_teid());
  EXPECT_EQ(0, pending_activation[0].teids.enb_teid());

  // For 5g verification
  EXPECT_EQ(0, pending_bearer_setup.size());
  session_state->remove_all_rules_for_termination(&update_criteria);

  EXPECT_EQ(update_criteria.static_rules_to_uninstall.size(), 1);
  EXPECT_EQ(update_criteria.dynamic_rules_to_uninstall.size(), 2);

  std::vector<std::string> rules_out{};
  std::vector<std::string>& rules_out_ptr = rules_out;

  session_state->get_dynamic_rules().get_rule_ids(rules_out_ptr);
  EXPECT_EQ(rules_out_ptr.size(), 0);

  EXPECT_FALSE(session_state->is_dynamic_rule_installed("dynamic-1"));
  EXPECT_FALSE(session_state->is_dynamic_rule_installed("dynamic-2"));

  std::string mkey;
  // searching for non-existent rule should fail
  EXPECT_FALSE(
      session_state->get_dynamic_rules().get_monitoring_key_for_rule_id(
          "dynamic-1", &mkey));
  EXPECT_FALSE(
      session_state->get_dynamic_rules().get_monitoring_key_for_rule_id(
          "dynamic-2", &mkey));
  // deleting an already deleted rule should fail
  EXPECT_FALSE(
      session_state->get_dynamic_rules().remove_rule("dynamic-1", &dynamic_1));
  EXPECT_FALSE(
      session_state->get_dynamic_rules().remove_rule("dynamic-2", &dynamic_2));
}

TEST_F(SessionStateTest5G, test_marshal_unmarshal) {
  session_state->set_config(cfg, nullptr);
  // For 5g verification
  EXPECT_TRUE(session_state->is_5g_session());
  EXPECT_EQ(update_criteria.static_rules_to_install.size(), 0);
  activate_rule(1, "m1", "rule1", STATIC, 0, 0);
  EXPECT_EQ(session_state->is_static_rule_installed("rule1"), true);
  EXPECT_EQ(true, session_state->active_monitored_rules_exist());
  EXPECT_EQ(update_criteria.static_rules_to_install.size(), 1);

  std::time_t activation_time =
      static_cast<std::time_t>(std::stoul("2020:04:15 09:10:11"));
  std::time_t deactivation_time =
      static_cast<std::time_t>(std::stoul("2020:04:15 09:10:12"));

  EXPECT_EQ(update_criteria.new_rule_lifetimes.size(), 1);
  schedule_rule(1, "m1", "rule2", DYNAMIC, activation_time, deactivation_time);
  EXPECT_EQ(session_state->is_dynamic_rule_installed("rule2"), false);
  EXPECT_EQ(update_criteria.static_rules_to_install.size(), 1);

  session_state->increment_rule_stats("rule1", nullptr);
  auto marshaled = session_state->marshal();
  auto unmarshaled = SessionState::unmarshal(marshaled, *rule_store);
  EXPECT_EQ(unmarshaled->is_static_rule_installed("rule1"), true);
  EXPECT_EQ(unmarshaled->is_dynamic_rule_installed("rule2"), false);
}

int main(int argc, char** argv) {
  ::testing::InitGoogleTest(&argc, argv);
  FLAGS_logtostderr = 1;
  FLAGS_v = 10;
  return RUN_ALL_TESTS();
}

}  // namespace magma

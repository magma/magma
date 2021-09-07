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
#include "SessionStateTester.h"

using ::testing::Test;

namespace magma {

class SessionStateTest5G : public ::testing::Test {
 protected:
  virtual void SetUp() {
    Teids teids;
    rule_store = std::make_shared<StaticRuleStore>();
    cfg        = build_sm_context(IMSI1, "10.20.30.40", 5);
    session_state =
        std::make_shared<SessionState>(IMSI1, SESSION_ID_1, cfg, *rule_store);
    session_state->set_fsm_state(SESSION_ACTIVE, nullptr);
    session_state->set_create_session_response(
        CreateSessionResponse(), nullptr);
    update_criteria = get_default_update_criteria();
  }

  void insert_static_rule_into_store(
      uint32_t rating_group, const std::string& m_key,
      const std::string& rule_id) {
    rule_store->insert_rule(create_policy_rule(rule_id, m_key, rating_group));
  }

  void insert_static_rule_with_qos_into_store(
      uint32_t rating_group, const std::string& m_key, const int qci,
      const std::string& rule_id) {
    PolicyRule rule =
        create_policy_rule_with_qos(rule_id, m_key, rating_group, qci);
    rule_store->insert_rule(rule);
  }

  uint32_t get_monitored_rule_count(const std::string& mkey) {
    std::vector<PolicyRule> rules;
    EXPECT_TRUE(session_state->get_dynamic_rules().get_rules(rules));
    uint32_t count = 0;
    for (PolicyRule& rule : rules) {
      if (rule.monitoring_key() == mkey) {
        count++;
      }
    }
    return count;
  }

  SessionConfig build_sm_context(
      const std::string& imsi,  // assumes IMSI prefix
      const std::string& ue_ipv4, uint32_t pdu_id) {
    SetSMSessionContext request;
    auto* req =
        request.mutable_rat_specific_context()->mutable_m5gsm_session_context();
    auto* reqcmn = request.mutable_common_context();
    req->set_pdu_session_id(pdu_id);
    req->set_request_type(magma::RequestType::INITIAL_REQUEST);
    req->mutable_pdu_address()->set_redirect_address_type(
        magma::RedirectServer::IPV4);
    req->mutable_pdu_address()->set_redirect_server_address(ue_ipv4);
    req->set_priority_access(magma::priorityaccess::High);
    req->set_imei("123456789012345");
    req->set_gpsi("9876543210");
    req->set_pcf_id("1357924680123456");

    reqcmn->mutable_sid()->set_id(imsi);
    reqcmn->set_sm_session_state(magma::SMSessionFSMState::CREATING_0);

    SessionConfig cfg;
    cfg.common_context       = request.common_context();
    cfg.rat_specific_context = request.rat_specific_context();
    cfg.rat_specific_context.mutable_m5gsm_session_context()->set_ssc_mode(
        SSC_MODE_3);
    return cfg;
  }

  uint32_t activate_5g_rule(
      uint32_t rating_group, const std::string& m_key,
      const std::string& rule_id, PolicyType rule_type,
      std::time_t activation_time, std::time_t deactivation_time) {
    PolicyRule rule = create_policy_rule(rule_id, m_key, rating_group);
    RuleLifetime lifetime(activation_time, deactivation_time);
    switch (rule_type) {
      case STATIC:
        rule_store->insert_rule(rule);
        return session_state
            ->activate_static_5g_rule(rule_id, lifetime, &update_criteria)
            .version;
        break;
      case DYNAMIC:
        return session_state
            ->insert_dynamic_5g_rule(rule, lifetime, &update_criteria)
            .version;
        break;
      default:
        break;
    }
    return 0;
  }

 protected:
  std::shared_ptr<StaticRuleStore> rule_store;
  std::shared_ptr<SessionState> session_state;
  SessionStateUpdateCriteria update_criteria;
  SessionConfig cfg;
};

TEST_F(SessionStateTest5G, test_session_rules) {
  activate_5g_rule(1, "m1", "rule1", DYNAMIC, 0, 0);
  EXPECT_EQ(1, session_state->total_monitored_rules_count());
  activate_5g_rule(2, "m2", "rule2", STATIC, 0, 0);
  EXPECT_EQ(2, session_state->total_monitored_rules_count());
  // add a OCS-ONLY static rule
  activate_5g_rule(3, "", "rule3", STATIC, 0, 0);
  EXPECT_EQ(2, session_state->total_monitored_rules_count());

  std::vector<std::string> rules_out{};
  std::vector<std::string>& rules_out_ptr = rules_out;

  session_state->get_dynamic_rules().get_rule_ids(rules_out_ptr);
  EXPECT_EQ(rules_out_ptr.size(), 1);
  EXPECT_EQ(rules_out_ptr[0], "rule1");

  EXPECT_EQ(session_state->is_static_rule_installed("rule2"), true);
  EXPECT_EQ(session_state->is_static_rule_installed("rule3"), true);
  EXPECT_EQ(session_state->is_static_rule_installed("rule_DNE"), false);

  EXPECT_EQ(session_state->get_current_rule_version("rule2"), 1);
  EXPECT_EQ(session_state->get_current_rule_version("rule3"), 1);

  // Test rule removals
  PolicyRule rule_out;
  session_state->deactivate_static_5g_rule("rule2", &update_criteria);
  EXPECT_EQ(1, session_state->total_monitored_rules_count());
  EXPECT_TRUE(session_state->remove_dynamic_5g_rule(
      "rule1", &rule_out, &update_criteria));
  EXPECT_EQ("m1", rule_out.monitoring_key());
  EXPECT_EQ(0, session_state->total_monitored_rules_count());

  // basic sanity checks to see it's properly deleted
  rules_out = {};
  session_state->get_dynamic_rules().get_rule_ids(rules_out_ptr);
  EXPECT_EQ(rules_out_ptr.size(), 0);

  EXPECT_EQ(0, get_monitored_rule_count("m1"));

  std::string mkey;
  // searching for non-existent rule should fail
  EXPECT_EQ(
      false, session_state->get_dynamic_rules().get_monitoring_key_for_rule_id(
                 "rule1", &mkey));
  // deleting an already deleted rule should fail
  EXPECT_EQ(
      false,
      session_state->get_dynamic_rules().remove_rule("rule1", &rule_out));
}

TEST_F(SessionStateTest5G, test_get_session_rules) {
  // populate rule store with 2 static and 2 dynamic rules
  activate_5g_rule(1, "", "rule-static-1", STATIC, 0, 0);
  activate_5g_rule(2, "m1", "rule-static-2", STATIC, 0, 0);
  activate_5g_rule(1, "", "rule-dynamic-1", DYNAMIC, 0, 0);
  activate_5g_rule(2, "m1", "rule-dynamic-2", DYNAMIC, 0, 0);

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
  const std::string activate_rule1   = pending_activation[0].rule.id();
  const std::string activate_rule2   = pending_activation[1].rule.id();
  const std::string deactivate_rule1 = pending_deactivation[0].rule.id();
  const std::string deactivate_rule2 = pending_deactivation[1].rule.id();
  EXPECT_TRUE(
      activate_rule1 == "rule-static-2" || activate_rule1 == "rule-dynamic-2");
  EXPECT_TRUE(
      activate_rule2 == "rule-static-3" || activate_rule2 == "rule-dynamic-3");
  EXPECT_TRUE(
      deactivate_rule1 == "rule-static-2" ||
      deactivate_rule1 == "rule-dynamic-2");
  EXPECT_TRUE(
      deactivate_rule2 == "rule-static-3" ||
      deactivate_rule2 == "rule-dynamic-3");
}

TEST_F(SessionStateTest5G, test_process_static_rule_installs) {
  // Insert 2 static rules without qos into static rule store
  insert_static_rule_into_store(0, "mkey1", "static-1");
  insert_static_rule_into_store(0, "mkey1", "static-2");
  // Insert 2 static rules with qos into static rule store
  insert_static_rule_with_qos_into_store(0, "mkey1", 1, "static-qos-3");
  insert_static_rule_with_qos_into_store(0, "mkey1", 2, "static-qos-4");

  // activate static-1 and static-qos-3 in advance
  RuleLifetime lifetime;
  session_state->activate_static_5g_rule("static-1", lifetime, nullptr);
  session_state->activate_static_5g_rule("static-qos-3", lifetime, nullptr);

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
  RulesToProcess pending_activation, pending_deactivation;
  session_state->process_static_5g_rule_installs(
      rule_installs, &pending_activation, &pending_deactivation,
      &update_criteria);
  EXPECT_EQ(2, pending_activation.size());
  EXPECT_EQ("static-2", pending_activation[0].rule.id());
  EXPECT_EQ(1, pending_activation[0].version);
  EXPECT_EQ("static-qos-4", pending_activation[1].rule.id());
  EXPECT_EQ(1, pending_activation[1].version);

  EXPECT_EQ(update_criteria.static_rules_to_install.size(), 2);
  EXPECT_TRUE(update_criteria.static_rules_to_install.count("static-2"));
  EXPECT_TRUE(update_criteria.static_rules_to_install.count("static-qos-4"));
}

TEST_F(SessionStateTest5G, test_process_dynamic_rule_installs) {
  PolicyRule dynamic_1 = create_policy_rule("dynamic-1", "", 0);
  PolicyRule dynamic_2 = create_policy_rule("dynamic-2", "", 0);
  PolicyRule dynamic_qos_3 =
      create_policy_rule_with_qos("dynamic-qos-3", "", 0, 1);
  PolicyRule dynamic_qos_4 =
      create_policy_rule_with_qos("dynamic-qos-4", "", 0, 2);

  // Install dynamic rules for dynamic-1 and dynamic-qos-3
  // Then version should be increased.
  RuleLifetime lifetime;
  session_state->insert_dynamic_5g_rule(dynamic_1, lifetime, nullptr);
  session_state->insert_dynamic_5g_rule(dynamic_qos_3, lifetime, nullptr);

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
  RulesToProcess pending_activation, pending_deactivation;
  session_state->process_dynamic_5g_rule_installs(
      rule_installs, &pending_activation, &pending_deactivation,
      &update_criteria);
  EXPECT_EQ(4, pending_activation.size());
  EXPECT_EQ("dynamic-1", pending_activation[0].rule.id());
  EXPECT_EQ(2, pending_activation[0].version);
  EXPECT_EQ("dynamic-2", pending_activation[1].rule.id());
  EXPECT_EQ(1, pending_activation[1].version);
  EXPECT_EQ("dynamic-qos-3", pending_activation[2].rule.id());
  EXPECT_EQ(2, pending_activation[2].version);
  EXPECT_EQ("dynamic-qos-4", pending_activation[3].rule.id());
  EXPECT_EQ(1, pending_activation[3].version);

  EXPECT_EQ(update_criteria.dynamic_rules_to_install.size(), 4);
  EXPECT_EQ("dynamic-1", update_criteria.dynamic_rules_to_install[0].id());
  EXPECT_EQ("dynamic-2", update_criteria.dynamic_rules_to_install[1].id());
  EXPECT_EQ("dynamic-qos-3", update_criteria.dynamic_rules_to_install[2].id());
  EXPECT_EQ("dynamic-qos-4", update_criteria.dynamic_rules_to_install[3].id());
}

TEST_F(SessionStateTest5G, test_remove_all_session_rules) {
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

  RulesToProcess pending_activation, pending_deactivation;
  session_state->process_5g_rules_to_install(
      static_rule_installs, dynamic_rule_installs, &pending_activation,
      &pending_deactivation, &update_criteria);

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

  session_state->remove_all_5g_rules_for_termination(&update_criteria);

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

int main(int argc, char** argv) {
  ::testing::InitGoogleTest(&argc, argv);
  FLAGS_logtostderr = 1;
  FLAGS_v           = 10;
  return RUN_ALL_TESTS();
}

}  // namespace magma

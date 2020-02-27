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

#include "ProtobufCreators.h"
#include "SessionState.h"
#include "magma_logging.h"

using ::testing::Test;

namespace magma {
const SessionState::Config test_sstate_cfg = {.ue_ipv4 = "127.0.0.1",
                                              .spgw_ipv4 = "128.0.0.1"};

class SessionStateTest : public ::testing::Test {
 protected:
 protected:
  virtual void SetUp()
  {
    auto tgpp_ctx = TgppContext();
    create_tgpp_context("gx.dest.com", "gy.dest.com", &tgpp_ctx);
    rule_store = std::make_shared<StaticRuleStore>();
    session_state = std::make_shared<SessionState>(
      "imsi", "session", "", test_sstate_cfg, *rule_store, tgpp_ctx);
  }
  enum RuleType {
    STATIC = 0,
    DYNAMIC = 1,
  };

  void insert_rule(
    uint32_t rating_group,
    const std::string& m_key,
    const std::string& rule_id,
    RuleType rule_type)
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
    switch (rule_type)
    {
      case STATIC:
        // insert into list of existing rules
        rule_store->insert_rule(rule);
        // mark the rule as active in session
        session_state->activate_static_rule(rule_id);
        break;
      case DYNAMIC:
        session_state->insert_dynamic_rule(rule);
        break;
    }
  }

  void receive_credit_from_ocs(uint32_t rating_group, uint64_t volume)
  {
    CreditUpdateResponse charge_resp;
    create_credit_update_response("IMSI1", rating_group, volume, &charge_resp);
    session_state->get_charging_pool().receive_credit(charge_resp);
  }

  void receive_credit_from_pcrf(
    const std::string& mkey,
    uint64_t volume,
    MonitoringLevel level)
  {
    UsageMonitoringUpdateResponse monitor_resp;
    create_monitor_update_response("IMSI1", mkey, level, volume, &monitor_resp);
    session_state->get_monitor_pool().receive_credit(monitor_resp);
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
    RuleType rule_type)
  {
    PolicyRule rule = get_rule(rating_group, m_key, rule_id);
    switch (rule_type) {
      case STATIC:
        rule_store->insert_rule(rule);
        session_state->activate_static_rule(rule_id);
        break;
      case DYNAMIC:
        session_state->insert_dynamic_rule(rule);
        break;
    }
  }

 protected:
  std::shared_ptr<StaticRuleStore> rule_store;
  std::shared_ptr<SessionState> session_state;
};

TEST_F(SessionStateTest, test_session_rules)
{
  activate_rule(1, "m1", "rule1", DYNAMIC);
  EXPECT_EQ(1, session_state->total_monitored_rules_count());
  activate_rule(2, "m2", "rule2", STATIC);
  EXPECT_EQ(2, session_state->total_monitored_rules_count());
  // add a OCS-ONLY static rule
  activate_rule(3, "", "rule3", STATIC);
  EXPECT_EQ(2, session_state->total_monitored_rules_count());

  std::vector<std::string> rules_out{};
  std::vector<std::string>& rules_out_ptr = rules_out;

  session_state->get_dynamic_rules().get_rule_ids(rules_out_ptr);
  EXPECT_EQ(rules_out_ptr.size(), 1);
  EXPECT_EQ(rules_out_ptr[0], "rule1");

  EXPECT_EQ(session_state->is_static_rule_installed("rule2"), true);
  EXPECT_EQ(session_state->is_static_rule_installed("rule3"), true);
  EXPECT_EQ(session_state->is_static_rule_installed("rule_DNE"), false);

  // Test rule removals
  PolicyRule rule_out;
  session_state->deactivate_static_rule("rule2");
  EXPECT_EQ(1, session_state->total_monitored_rules_count());
  EXPECT_EQ(true, session_state->remove_dynamic_rule("rule1", &rule_out));
  EXPECT_EQ("m1", rule_out.monitoring_key());
  EXPECT_EQ(0, session_state->total_monitored_rules_count());

  // basic sanity checks to see it's properly deleted
  rules_out = {};
  session_state->get_dynamic_rules().get_rule_ids(rules_out_ptr);
  EXPECT_EQ(rules_out_ptr.size(), 0);

  rules_out = {};
  session_state->get_dynamic_rules()
    .get_rule_ids_for_monitoring_key("m1", rules_out);
  EXPECT_EQ(0, rules_out.size());

  std::string mkey;
  // searching for non-existent rule should fail
  EXPECT_EQ(false, session_state->get_dynamic_rules()
    .get_monitoring_key_for_rule_id("rule1", &mkey));
  // deleting an already deleted rule should fail
  EXPECT_EQ(false, session_state->get_dynamic_rules()
    .remove_rule("rule1", &rule_out));
}

TEST_F(SessionStateTest, test_marshal_unmarshal)
{
  insert_rule(1, "m1", "rule1", STATIC);
  EXPECT_EQ(session_state->is_static_rule_installed("rule1"), true);
  EXPECT_EQ(true, session_state->active_monitored_rules_exist());

  receive_credit_from_ocs(1, 1024);
  EXPECT_EQ(
    session_state->get_charging_pool().get_credit(1, ALLOWED_TOTAL), 1024);

  receive_credit_from_pcrf("m1", 1024, MonitoringLevel::PCC_RULE_LEVEL);
  EXPECT_EQ(
    session_state->get_monitor_pool().get_credit("m1", ALLOWED_TOTAL), 1024);

  auto marshaled = session_state->marshal();
  auto unmarshaled = SessionState::unmarshal(marshaled, *rule_store);
  EXPECT_EQ(
    unmarshaled->get_charging_pool().get_credit(1, ALLOWED_TOTAL), 1024);
  EXPECT_EQ(
    unmarshaled->get_monitor_pool().get_credit("m1", ALLOWED_TOTAL), 1024);
  EXPECT_EQ(unmarshaled->is_static_rule_installed("rule1"), true);
}

TEST_F(SessionStateTest, test_insert_credit)
{
  insert_rule(1, "m1", "rule1", STATIC);
  EXPECT_EQ(true, session_state->active_monitored_rules_exist());

  receive_credit_from_ocs(1, 1024);
  EXPECT_EQ(
    session_state->get_charging_pool().get_credit(1, ALLOWED_TOTAL), 1024);

  receive_credit_from_pcrf("m1", 1024, MonitoringLevel::PCC_RULE_LEVEL);
  EXPECT_EQ(
    session_state->get_monitor_pool().get_credit("m1", ALLOWED_TOTAL), 1024);
}

TEST_F(SessionStateTest, test_termination)
{
  std::promise<void> termination_promise;
  session_state->start_termination(
    [&termination_promise](SessionTerminateRequest term_req) {
      termination_promise.set_value();
    });
  session_state->complete_termination();
  auto status =
    termination_promise.get_future().wait_for(std::chrono::seconds(0));
  EXPECT_EQ(status, std::future_status::ready);
}

TEST_F(SessionStateTest, test_can_complete_termination)
{
  insert_rule(1, "m1", "rule1", STATIC);
  EXPECT_EQ(true, session_state->active_monitored_rules_exist());

  EXPECT_EQ(session_state->can_complete_termination(), false);

  session_state->start_termination([](SessionTerminateRequest term_req) {});
  EXPECT_EQ(session_state->can_complete_termination(), false);

  // If the rule is still being reported, termination should not be completed.
  session_state->new_report();
  EXPECT_EQ(session_state->can_complete_termination(), false);
  session_state->add_used_credit("rule1", 100, 100);
  EXPECT_EQ(session_state->can_complete_termination(), false);
  session_state->finish_report();
  EXPECT_EQ(session_state->can_complete_termination(), false);

  // The rule is not reported, termination can be completed.
  session_state->new_report();
  EXPECT_EQ(session_state->can_complete_termination(), false);
  session_state->finish_report();
  EXPECT_EQ(session_state->can_complete_termination(), true);

  // Termination should only be completed once.
  session_state->complete_termination();
  EXPECT_EQ(session_state->can_complete_termination(), false);
}

TEST_F(SessionStateTest, test_add_used_credit)
{
  insert_rule(1, "m1", "rule1", STATIC);
  insert_rule(2, "m2", "dyn_rule1", DYNAMIC);
  EXPECT_EQ(true, session_state->active_monitored_rules_exist());

  receive_credit_from_ocs(1, 3000);
  receive_credit_from_ocs(2, 6000);

  receive_credit_from_pcrf("m1", 3000, MonitoringLevel::PCC_RULE_LEVEL);
  receive_credit_from_pcrf("m2", 6000, MonitoringLevel::PCC_RULE_LEVEL);

  session_state->add_used_credit("rule1", 2000, 1000);
  EXPECT_EQ(session_state->get_charging_pool().get_credit(1, USED_TX), 2000);
  EXPECT_EQ(session_state->get_charging_pool().get_credit(1, USED_RX), 1000);
  EXPECT_EQ(session_state->get_monitor_pool().get_credit("m1", USED_TX), 2000);
  EXPECT_EQ(session_state->get_monitor_pool().get_credit("m1", USED_RX), 1000);

  session_state->add_used_credit("dyn_rule1", 4000, 2000);
  EXPECT_EQ(session_state->get_charging_pool().get_credit(2, USED_TX), 4000);
  EXPECT_EQ(session_state->get_charging_pool().get_credit(2, USED_RX), 2000);
  EXPECT_EQ(session_state->get_monitor_pool().get_credit("m2", USED_TX), 4000);
  EXPECT_EQ(session_state->get_monitor_pool().get_credit("m2", USED_RX), 2000);

  UpdateSessionRequest update;
  std::vector<std::unique_ptr<ServiceAction>> actions;
  session_state->get_updates(update, &actions);
  EXPECT_EQ(actions.size(), 0);
  EXPECT_EQ(update.updates_size(), 2);
  EXPECT_EQ(update.usage_monitors_size(), 2);

  PolicyRule policy_out;
  EXPECT_EQ(true, session_state->
    remove_dynamic_rule("dyn_rule1", &policy_out));
  EXPECT_EQ(true, session_state->deactivate_static_rule("rule1"));
  EXPECT_EQ(false, session_state->active_monitored_rules_exist());
}

TEST_F(SessionStateTest, test_mixed_tracking_rules)
{
  insert_rule(0, "m1", "dyn_rule1", DYNAMIC);
  insert_rule(2, "", "dyn_rule2", DYNAMIC);
  insert_rule(3, "m3", "dyn_rule3", DYNAMIC);
  EXPECT_EQ(true, session_state->active_monitored_rules_exist());

  receive_credit_from_ocs(2, 6000);
  receive_credit_from_ocs(3, 8000);

  receive_credit_from_pcrf("m1", 3000, MonitoringLevel::PCC_RULE_LEVEL);
  receive_credit_from_pcrf("m3", 8000, MonitoringLevel::PCC_RULE_LEVEL);

  session_state->add_used_credit("dyn_rule1", 2000, 1000);
  EXPECT_EQ(session_state->get_monitor_pool().get_credit("m1", USED_TX), 2000);
  EXPECT_EQ(session_state->get_monitor_pool().get_credit("m1", USED_RX), 1000);

  session_state->add_used_credit("dyn_rule2", 4000, 2000);
  EXPECT_EQ(session_state->get_charging_pool().get_credit(2, USED_TX), 4000);
  EXPECT_EQ(session_state->get_charging_pool().get_credit(2, USED_RX), 2000);
  session_state->add_used_credit("dyn_rule3", 5000, 3000);
  EXPECT_EQ(session_state->get_charging_pool().get_credit(3, USED_TX), 5000);
  EXPECT_EQ(session_state->get_charging_pool().get_credit(3, USED_RX), 3000);
  EXPECT_EQ(session_state->get_monitor_pool().get_credit("m3", USED_TX), 5000);
  EXPECT_EQ(session_state->get_monitor_pool().get_credit("m3", USED_RX), 3000);

  UpdateSessionRequest update;
  std::vector<std::unique_ptr<ServiceAction>> actions;
  session_state->get_updates(update, &actions);
  EXPECT_EQ(actions.size(), 0);
  EXPECT_EQ(update.updates_size(), 2);
  EXPECT_EQ(update.usage_monitors_size(), 2);
}

TEST_F(SessionStateTest, test_session_level_key)
{
  EXPECT_EQ(nullptr, session_state->get_monitor_pool().get_session_level_key());

  receive_credit_from_pcrf("m1", 8000, MonitoringLevel::SESSION_LEVEL);
  EXPECT_EQ("m1", *session_state->get_monitor_pool().get_session_level_key());
  EXPECT_EQ(
    session_state->get_monitor_pool().get_credit("m1", ALLOWED_TOTAL), 8000);

  session_state->add_used_credit("rule1", 5000, 3000);
  EXPECT_EQ(session_state->get_monitor_pool().get_credit("m1", USED_TX), 5000);
  EXPECT_EQ(session_state->get_monitor_pool().get_credit("m1", USED_RX), 3000);

  UpdateSessionRequest update;
  std::vector<std::unique_ptr<ServiceAction>> actions;
  session_state->get_updates(update, &actions);
  EXPECT_EQ(actions.size(), 0);
  EXPECT_EQ(update.usage_monitors_size(), 1);
  auto& single_update = update.usage_monitors(0).update();
  EXPECT_EQ(single_update.level(), MonitoringLevel::SESSION_LEVEL);
  EXPECT_EQ(single_update.bytes_rx(), 3000);
  EXPECT_EQ(single_update.bytes_tx(), 5000);
}

TEST_F(SessionStateTest, test_reauth_key)
{
  insert_rule(1, "", "rule1", STATIC);

  receive_credit_from_ocs(1, 1500);

  session_state->add_used_credit("rule1", 1000, 500);

  UpdateSessionRequest update;
  std::vector<std::unique_ptr<ServiceAction>> actions;
  session_state->get_updates(update, &actions);
  EXPECT_EQ(update.updates_size(), 1);
  EXPECT_EQ(
    session_state->get_charging_pool().get_credit(1, REPORTING_TX), 1000);
  EXPECT_EQ(
    session_state->get_charging_pool().get_credit(1, REPORTING_RX), 500);
  // credit is already reporting, no update needed
  auto reauth_res = session_state->get_charging_pool().reauth_key(1);
  EXPECT_EQ(reauth_res, ChargingReAuthAnswer::UPDATE_NOT_NEEDED);
  receive_credit_from_ocs(1, 1024);
  EXPECT_EQ(session_state->get_charging_pool().get_credit(1, REPORTING_TX), 0);
  EXPECT_EQ(session_state->get_charging_pool().get_credit(1, REPORTING_RX), 0);
  reauth_res = session_state->get_charging_pool().reauth_key(1);
  EXPECT_EQ(reauth_res, ChargingReAuthAnswer::UPDATE_INITIATED);

  session_state->add_used_credit("rule1", 2, 1);
  UpdateSessionRequest reauth_update;
  session_state->get_updates(reauth_update, &actions);
  EXPECT_EQ(reauth_update.updates_size(), 1);
  auto& usage = reauth_update.updates(0).usage();
  EXPECT_EQ(usage.bytes_tx(), 2);
  EXPECT_EQ(usage.bytes_rx(), 1);
}

TEST_F(SessionStateTest, test_reauth_new_key)
{
  // credit is already reporting, no update needed
  auto reauth_res = session_state->get_charging_pool().reauth_key(1);
  EXPECT_EQ(reauth_res, ChargingReAuthAnswer::UPDATE_INITIATED);

  UpdateSessionRequest reauth_update;
  std::vector<std::unique_ptr<ServiceAction>> actions;
  session_state->get_updates(reauth_update, &actions);
  EXPECT_EQ(reauth_update.updates_size(), 1);
  auto& usage = reauth_update.updates(0).usage();
  EXPECT_EQ(usage.charging_key(), 1);
  EXPECT_EQ(usage.bytes_tx(), 0);
  EXPECT_EQ(usage.bytes_rx(), 0);

  receive_credit_from_ocs(1, 1024);
  EXPECT_EQ(
    session_state->get_charging_pool().get_credit(1, ALLOWED_TOTAL), 1024);
}

TEST_F(SessionStateTest, test_reauth_all)
{
  insert_rule(1, "", "rule1", STATIC);
  insert_rule(2, "", "dyn_rule1", DYNAMIC);
  EXPECT_EQ(false, session_state->active_monitored_rules_exist());

  receive_credit_from_ocs(1, 1024);
  receive_credit_from_ocs(2, 1024);

  session_state->add_used_credit("rule1", 10, 20);
  session_state->add_used_credit("dyn_rule1", 30, 40);
  // If any charging key isn't reporting, an update is needed
  auto reauth_res = session_state->get_charging_pool().reauth_all();
  EXPECT_EQ(reauth_res, ChargingReAuthAnswer::UPDATE_INITIATED);

  UpdateSessionRequest reauth_update;
  std::vector<std::unique_ptr<ServiceAction>> actions;
  session_state->get_updates(reauth_update, &actions);
  EXPECT_EQ(reauth_update.updates_size(), 2);

  // All charging keys are reporting, no update needed
  reauth_res = session_state->get_charging_pool().reauth_all();
  EXPECT_EQ(reauth_res, ChargingReAuthAnswer::UPDATE_NOT_NEEDED);
}

TEST_F(SessionStateTest, test_tgpp_context_is_set_on_update)
{
  receive_credit_from_pcrf("m1", 1024, MonitoringLevel::PCC_RULE_LEVEL);
  receive_credit_from_ocs(1, 1024);
  insert_rule(1, "m1", "rule1", STATIC);
  session_state->add_used_credit("rule1", 1024, 0);
  EXPECT_EQ(true, session_state->active_monitored_rules_exist());

  EXPECT_EQ(
    session_state->get_monitor_pool().get_credit("m1", ALLOWED_TOTAL), 1024);
  EXPECT_EQ(
    session_state->get_charging_pool().get_credit(1, ALLOWED_TOTAL), 1024);

  UpdateSessionRequest update;
  std::vector<std::unique_ptr<ServiceAction>> actions;
  session_state->get_updates(update,& actions);
  EXPECT_EQ(actions.size(), 0);
  EXPECT_EQ(update.updates_size(), 1);
  EXPECT_EQ(update.updates().Get(0).tgpp_ctx().gx_dest_host(), "gx.dest.com");
  EXPECT_EQ(update.updates().Get(0).tgpp_ctx().gy_dest_host(), "gy.dest.com");
  EXPECT_EQ(update.usage_monitors_size(), 1);
  EXPECT_EQ(update.usage_monitors().Get(0).tgpp_ctx().gx_dest_host(), "gx.dest.com");
  EXPECT_EQ(update.usage_monitors().Get(0).tgpp_ctx().gy_dest_host(), "gy.dest.com");
}

int main(int argc, char **argv)
{
  ::testing::InitGoogleTest(&argc, argv);
  FLAGS_logtostderr = 1;
  FLAGS_v = 10;
  return RUN_ALL_TESTS();
}

} // namespace magma

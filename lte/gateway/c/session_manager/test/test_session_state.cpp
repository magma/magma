/**
 * Copyright (c) 2016-present, Facebook, Inc.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */
#include <future>
#include <memory>
#include <utility>

#include <glog/logging.h>
#include <gtest/gtest.h>

#include "ProtobufCreators.h"
#include "SessionState.h"
#include "magma_logging.h"

using ::testing::Test;

namespace magma {
const SessionConfig test_sstate_cfg = {.ue_ipv4 = "127.0.0.1",
                                       .spgw_ipv4 = "128.0.0.1"};

class SessionStateTest : public ::testing::Test {
protected:
protected:
  virtual void SetUp() {
    auto tgpp_ctx = TgppContext();
    create_tgpp_context("gx.dest.com", "gy.dest.com", &tgpp_ctx);
    rule_store = std::make_shared<StaticRuleStore>();
    session_state = std::make_shared<SessionState>(
        "imsi", "session", "", test_sstate_cfg, *rule_store, tgpp_ctx);
    update_criteria = get_default_update_criteria();
  }
  enum RuleType {
    STATIC = 0,
    DYNAMIC = 1,
  };

  PolicyRule build_rule(uint32_t rating_group, const std::string &m_key,
                   const std::string &rule_id) {
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

  void insert_rule(uint32_t rating_group, const std::string &m_key,
                   const std::string &rule_id, RuleType rule_type,
                   std::time_t activation_time, std::time_t deactivation_time) {
    PolicyRule rule = build_rule(rating_group, m_key, rule_id);
    RuleLifetime lifetime{
        .activation_time = activation_time,
        .deactivation_time = deactivation_time,
    };
    switch (rule_type) {
    case STATIC:
      // insert into list of existing rules
      rule_store->insert_rule(rule);
      // mark the rule as active in session
      session_state->activate_static_rule(rule_id, lifetime, update_criteria);
      break;
    case DYNAMIC:
      session_state->insert_dynamic_rule(rule, lifetime, update_criteria);
      break;
    }
  }

  void schedule_rule(uint32_t rating_group, const std::string &m_key,
                   const std::string &rule_id, RuleType rule_type,
                   std::time_t activation_time, std::time_t deactivation_time) {
    PolicyRule rule = build_rule(rating_group, m_key, rule_id);
    RuleLifetime lifetime{
      .activation_time = activation_time,
      .deactivation_time = deactivation_time,
    };
    switch (rule_type) {
      case STATIC:
        // insert into list of existing rules
        rule_store->insert_rule(rule);
        // mark the rule as scheduled in the session
        session_state->schedule_static_rule(rule_id, lifetime, update_criteria);
        break;
      case DYNAMIC:
        session_state->schedule_dynamic_rule(rule, lifetime, update_criteria);
        break;
    }
  }

  void receive_credit_from_ocs(uint32_t rating_group, uint64_t volume) {
    CreditUpdateResponse charge_resp;
    create_credit_update_response("IMSI1", rating_group, volume, &charge_resp);
    session_state->get_charging_pool().receive_credit(charge_resp,
                                                      update_criteria);
  }

  void receive_credit_from_pcrf(const std::string &mkey, uint64_t volume,
                                MonitoringLevel level) {
    UsageMonitoringUpdateResponse monitor_resp;
    create_monitor_update_response("IMSI1", mkey, level, volume, &monitor_resp);
    session_state->get_monitor_pool().receive_credit(monitor_resp,
                                                     update_criteria);
  }

  PolicyRule get_rule(uint32_t rating_group, const std::string &m_key,
                      const std::string &rule_id) {
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

  void activate_rule(uint32_t rating_group, const std::string &m_key,
                     const std::string &rule_id, RuleType rule_type,
                     std::time_t activation_time, std::time_t deactivation_time) {
      PolicyRule rule = get_rule(rating_group, m_key, rule_id);
    RuleLifetime lifetime{
        .activation_time = activation_time,
        .deactivation_time = deactivation_time,
    };
    switch (rule_type) {
    case STATIC:
      rule_store->insert_rule(rule);
      session_state->activate_static_rule(rule_id, lifetime, update_criteria);
      break;
    case DYNAMIC:
      session_state->insert_dynamic_rule(rule, lifetime, update_criteria);
      break;
    }
  }

protected:
  std::shared_ptr<StaticRuleStore> rule_store;
  std::shared_ptr<SessionState> session_state;
  SessionStateUpdateCriteria update_criteria;
};

TEST_F(SessionStateTest, test_session_rules) {
  activate_rule(1, "m1", "rule1", DYNAMIC, 0, 0);
  EXPECT_EQ(1, session_state->total_monitored_rules_count());
  activate_rule(2, "m2", "rule2", STATIC, 0, 0);
  EXPECT_EQ(2, session_state->total_monitored_rules_count());
  // add a OCS-ONLY static rule
  activate_rule(3, "", "rule3", STATIC, 0, 0);
  EXPECT_EQ(2, session_state->total_monitored_rules_count());

  std::vector<std::string> rules_out{};
  std::vector<std::string> &rules_out_ptr = rules_out;

  session_state->get_dynamic_rules().get_rule_ids(rules_out_ptr);
  EXPECT_EQ(rules_out_ptr.size(), 1);
  EXPECT_EQ(rules_out_ptr[0], "rule1");

  EXPECT_EQ(session_state->is_static_rule_installed("rule2"), true);
  EXPECT_EQ(session_state->is_static_rule_installed("rule3"), true);
  EXPECT_EQ(session_state->is_static_rule_installed("rule_DNE"), false);

  // Test rule removals
  PolicyRule rule_out;
  session_state->deactivate_static_rule("rule2", update_criteria);
  EXPECT_EQ(1, session_state->total_monitored_rules_count());
  EXPECT_EQ(true, session_state->remove_dynamic_rule("rule1", &rule_out,
                                                     update_criteria));
  EXPECT_EQ("m1", rule_out.monitoring_key());
  EXPECT_EQ(0, session_state->total_monitored_rules_count());

  // basic sanity checks to see it's properly deleted
  rules_out = {};
  session_state->get_dynamic_rules().get_rule_ids(rules_out_ptr);
  EXPECT_EQ(rules_out_ptr.size(), 0);

  rules_out = {};
  session_state->get_dynamic_rules().get_rule_ids_for_monitoring_key("m1",
                                                                     rules_out);
  EXPECT_EQ(0, rules_out.size());

  std::string mkey;
  // searching for non-existent rule should fail
  EXPECT_EQ(false,
            session_state->get_dynamic_rules().get_monitoring_key_for_rule_id(
                "rule1", &mkey));
  // deleting an already deleted rule should fail
  EXPECT_EQ(false,
            session_state->get_dynamic_rules().remove_rule("rule1", &rule_out));
}

/**
 * Check that rule scheduling and installation works from the perspective of
 * tracking in SessionState
 */
TEST_F(SessionStateTest, test_rule_scheduling) {
  auto _uc = get_default_update_criteria(); // unused

  // First schedule a dynamic and static rule. They are treated as inactive.
  schedule_rule(1, "m1", "rule1", DYNAMIC, 0, 0);
  EXPECT_EQ(0, session_state->total_monitored_rules_count());
  EXPECT_FALSE(session_state->is_dynamic_rule_installed("rule1"));

  schedule_rule(2, "m2", "rule2", STATIC, 0, 0);
  EXPECT_EQ(0, session_state->total_monitored_rules_count());
  EXPECT_FALSE(session_state->is_static_rule_installed("rule2"));

  // Now suppose some time has passed, and it's time to mark scheduled rules
  // as active. The responsibility is given to the session owner to make
  // these calls
  session_state->install_scheduled_dynamic_rule("rule1", _uc);
  EXPECT_EQ(1, session_state->total_monitored_rules_count());
  EXPECT_TRUE(session_state->is_dynamic_rule_installed("rule1"));

  session_state->install_scheduled_static_rule("rule2", _uc);
  EXPECT_EQ(2, session_state->total_monitored_rules_count());
  EXPECT_TRUE(session_state->is_static_rule_installed("rule2"));
}

TEST_F(SessionStateTest, test_marshal_unmarshal) {
  EXPECT_EQ(update_criteria.static_rules_to_install.size(), 0);
  insert_rule(1, "m1", "rule1", STATIC, 0, 0);
  EXPECT_EQ(session_state->is_static_rule_installed("rule1"), true);
  EXPECT_EQ(true, session_state->active_monitored_rules_exist());
  EXPECT_EQ(update_criteria.static_rules_to_install.size(), 1);

  std::time_t activation_time = static_cast<std::time_t>(std::stoul("2020:04:15 09:10:11"));
  std::time_t deactivation_time = static_cast<std::time_t>(std::stoul("2020:04:15 09:10:12"));

  EXPECT_EQ(update_criteria.new_rule_lifetimes.size(), 1);
  schedule_rule(1, "m1", "rule2", DYNAMIC, activation_time, deactivation_time);
  EXPECT_EQ(session_state->is_dynamic_rule_installed("rule2"), false);
  EXPECT_EQ(update_criteria.static_rules_to_install.size(), 1);

  EXPECT_EQ(update_criteria.charging_credit_to_install.size(), 0);
  receive_credit_from_ocs(1, 1024);
  EXPECT_EQ(update_criteria.charging_credit_to_install.size(), 1);
  EXPECT_EQ(session_state->get_charging_pool().get_credit(1, ALLOWED_TOTAL),
            1024);

  EXPECT_EQ(update_criteria.monitor_credit_to_install.size(), 0);
  receive_credit_from_pcrf("m1", 1024, MonitoringLevel::PCC_RULE_LEVEL);
  EXPECT_EQ(session_state->get_monitor_pool().get_credit("m1", ALLOWED_TOTAL),
            1024);
  EXPECT_EQ(update_criteria.monitor_credit_to_install.size(), 1);

  auto marshaled = session_state->marshal();
  auto unmarshaled = SessionState::unmarshal(marshaled, *rule_store);
  EXPECT_EQ(unmarshaled->get_charging_pool().get_credit(1, ALLOWED_TOTAL),
            1024);
  EXPECT_EQ(unmarshaled->get_monitor_pool().get_credit("m1", ALLOWED_TOTAL),
            1024);
  EXPECT_EQ(unmarshaled->is_static_rule_installed("rule1"), true);
  EXPECT_EQ(session_state->is_dynamic_rule_installed("rule2"), false);
}

TEST_F(SessionStateTest, test_insert_credit) {
  EXPECT_EQ(update_criteria.static_rules_to_install.size(), 0);
  insert_rule(1, "m1", "rule1", STATIC, 0, 0);
  EXPECT_EQ(true, session_state->active_monitored_rules_exist());
  EXPECT_TRUE(std::find(update_criteria.static_rules_to_install.begin(),
                        update_criteria.static_rules_to_install.end(),
                        "rule1") !=
              update_criteria.static_rules_to_install.end());

  receive_credit_from_ocs(1, 1024);
  EXPECT_EQ(session_state->get_charging_pool().get_credit(1, ALLOWED_TOTAL),
            1024);
  EXPECT_EQ(update_criteria.charging_credit_to_install[CreditKey(1)]
                .buckets[ALLOWED_TOTAL],
            1024);

  receive_credit_from_pcrf("m1", 1024, MonitoringLevel::PCC_RULE_LEVEL);
  EXPECT_EQ(session_state->get_monitor_pool().get_credit("m1", ALLOWED_TOTAL),
            1024);
  EXPECT_EQ(update_criteria.monitor_credit_to_install["m1"]
                .credit.buckets[ALLOWED_TOTAL],
            1024);
}

TEST_F(SessionStateTest, test_termination) {
  std::promise<void> termination_promise;
  session_state->start_termination(update_criteria);
  session_state->set_termination_callback([&termination_promise](
      SessionTerminateRequest term_req) { termination_promise.set_value(); });
  session_state->complete_termination(update_criteria);
  auto status =
      termination_promise.get_future().wait_for(std::chrono::seconds(0));
  EXPECT_EQ(status, std::future_status::ready);
}

TEST_F(SessionStateTest, test_can_complete_termination) {
  insert_rule(1, "m1", "rule1", STATIC, 0, 0);
  EXPECT_EQ(true, session_state->active_monitored_rules_exist());
  EXPECT_TRUE(std::find(update_criteria.static_rules_to_install.begin(),
                        update_criteria.static_rules_to_install.end(),
                        "rule1") !=
              update_criteria.static_rules_to_install.end());
  // Have not received credit
  EXPECT_EQ(update_criteria.monitor_credit_map.size(), 0);

  EXPECT_EQ(session_state->can_complete_termination(), false);

  session_state->start_termination(update_criteria);
  EXPECT_EQ(session_state->can_complete_termination(), false);

  // If the rule is still being reported, termination should not be completed.
  session_state->new_report();
  EXPECT_EQ(session_state->can_complete_termination(), false);
  session_state->add_used_credit("rule1", 100, 100, update_criteria);
  EXPECT_EQ(session_state->can_complete_termination(), false);
  EXPECT_EQ(update_criteria.monitor_credit_map.size(), 0);
  session_state->finish_report();
  EXPECT_EQ(session_state->can_complete_termination(), false);

  // The rule is not reported, termination can be completed.
  session_state->new_report();
  EXPECT_EQ(session_state->can_complete_termination(), false);
  session_state->finish_report();
  EXPECT_EQ(session_state->can_complete_termination(), true);

  // Termination should only be completed once.
  session_state->complete_termination(update_criteria);
  EXPECT_EQ(session_state->can_complete_termination(), false);
}

TEST_F(SessionStateTest, test_add_used_credit) {
  insert_rule(1, "m1", "rule1", STATIC, 0, 0);
  insert_rule(2, "m2", "dyn_rule1", DYNAMIC, 0, 0);
  EXPECT_EQ(true, session_state->active_monitored_rules_exist());
  EXPECT_TRUE(std::find(update_criteria.static_rules_to_install.begin(),
                        update_criteria.static_rules_to_install.end(),
                        "rule1") !=
              update_criteria.static_rules_to_install.end());
  EXPECT_EQ(update_criteria.dynamic_rules_to_install.size(), 1);

  receive_credit_from_ocs(1, 3000);
  receive_credit_from_ocs(2, 6000);
  EXPECT_EQ(update_criteria.charging_credit_to_install.size(), 2);
  EXPECT_EQ(update_criteria.charging_credit_to_install[CreditKey(1)]
                .buckets[ALLOWED_TOTAL],
            3000);

  receive_credit_from_pcrf("m1", 3000, MonitoringLevel::PCC_RULE_LEVEL);
  receive_credit_from_pcrf("m2", 6000, MonitoringLevel::PCC_RULE_LEVEL);
  EXPECT_EQ(update_criteria.monitor_credit_to_install.size(), 2);
  EXPECT_EQ(update_criteria.monitor_credit_to_install["m1"]
                .credit.buckets[ALLOWED_TOTAL],
            3000);

  session_state->add_used_credit("rule1", 2000, 1000, update_criteria);
  EXPECT_EQ(session_state->get_charging_pool().get_credit(1, USED_TX), 2000);
  EXPECT_EQ(session_state->get_charging_pool().get_credit(1, USED_RX), 1000);
  EXPECT_EQ(session_state->get_monitor_pool().get_credit("m1", USED_TX), 2000);
  EXPECT_EQ(session_state->get_monitor_pool().get_credit("m1", USED_RX), 1000);
  EXPECT_EQ(
      update_criteria.charging_credit_map[CreditKey(1)].bucket_deltas[USED_TX],
      2000);
  EXPECT_EQ(update_criteria.monitor_credit_map["m1"].bucket_deltas[USED_RX],
            1000);

  session_state->add_used_credit("dyn_rule1", 4000, 2000, update_criteria);
  EXPECT_EQ(session_state->get_charging_pool().get_credit(2, USED_TX), 4000);
  EXPECT_EQ(session_state->get_charging_pool().get_credit(2, USED_RX), 2000);
  EXPECT_EQ(session_state->get_monitor_pool().get_credit("m2", USED_TX), 4000);
  EXPECT_EQ(session_state->get_monitor_pool().get_credit("m2", USED_RX), 2000);
  EXPECT_EQ(
      update_criteria.charging_credit_map[CreditKey(2)].bucket_deltas[USED_TX],
      4000);
  EXPECT_EQ(update_criteria.monitor_credit_map["m2"].bucket_deltas[USED_RX],
            2000);

  UpdateSessionRequest update;
  std::vector<std::unique_ptr<ServiceAction>> actions;
  session_state->get_updates(update, &actions, update_criteria);
  EXPECT_EQ(actions.size(), 0);
  EXPECT_EQ(update.updates_size(), 2);
  EXPECT_EQ(update.usage_monitors_size(), 2);

  PolicyRule policy_out;
  EXPECT_EQ(true, session_state->remove_dynamic_rule("dyn_rule1", &policy_out,
                                                     update_criteria));
  EXPECT_EQ(true,
            session_state->deactivate_static_rule("rule1", update_criteria));
  EXPECT_EQ(false, session_state->active_monitored_rules_exist());
  EXPECT_TRUE(std::find(update_criteria.dynamic_rules_to_uninstall.begin(),
                        update_criteria.dynamic_rules_to_uninstall.end(),
                        "dyn_rule1") !=
              update_criteria.dynamic_rules_to_uninstall.end());
}

TEST_F(SessionStateTest, test_mixed_tracking_rules) {
  insert_rule(0, "m1", "dyn_rule1", DYNAMIC, 0, 0);
  insert_rule(2, "", "dyn_rule2", DYNAMIC, 0, 0);
  insert_rule(3, "m3", "dyn_rule3", DYNAMIC, 0, 0);
  EXPECT_EQ(true, session_state->active_monitored_rules_exist());
  // Installing a rule doesn't install credit
  EXPECT_EQ(update_criteria.charging_credit_to_install.size(), 0);
  EXPECT_EQ(update_criteria.dynamic_rules_to_install.size(), 3);

  receive_credit_from_ocs(2, 6000);
  receive_credit_from_ocs(3, 8000);
  EXPECT_EQ(update_criteria.charging_credit_to_install.size(), 2);
  EXPECT_EQ(update_criteria.charging_credit_to_install[CreditKey(2)]
                .buckets[ALLOWED_TOTAL],
            6000);

  receive_credit_from_pcrf("m1", 3000, MonitoringLevel::PCC_RULE_LEVEL);
  receive_credit_from_pcrf("m3", 8000, MonitoringLevel::PCC_RULE_LEVEL);
  EXPECT_EQ(update_criteria.monitor_credit_to_install["m1"]
                .credit.buckets[ALLOWED_TOTAL],
            3000);

  session_state->add_used_credit("dyn_rule1", 2000, 1000, update_criteria);
  EXPECT_EQ(session_state->get_monitor_pool().get_credit("m1", USED_TX), 2000);
  EXPECT_EQ(session_state->get_monitor_pool().get_credit("m1", USED_RX), 1000);
  EXPECT_EQ(update_criteria.monitor_credit_map["m1"].bucket_deltas[USED_TX],
            2000);
  EXPECT_EQ(update_criteria.monitor_credit_map["m1"].bucket_deltas[USED_RX],
            1000);

  session_state->add_used_credit("dyn_rule2", 4000, 2000, update_criteria);
  EXPECT_EQ(session_state->get_charging_pool().get_credit(2, USED_TX), 4000);
  EXPECT_EQ(session_state->get_charging_pool().get_credit(2, USED_RX), 2000);
  EXPECT_EQ(
      update_criteria.charging_credit_map[CreditKey(2)].bucket_deltas[USED_TX],
      4000);
  session_state->add_used_credit("dyn_rule3", 5000, 3000, update_criteria);
  EXPECT_EQ(session_state->get_charging_pool().get_credit(3, USED_TX), 5000);
  EXPECT_EQ(session_state->get_charging_pool().get_credit(3, USED_RX), 3000);
  EXPECT_EQ(
      update_criteria.charging_credit_map[CreditKey(3)].bucket_deltas[USED_TX],
      5000);
  EXPECT_EQ(session_state->get_monitor_pool().get_credit("m3", USED_TX), 5000);
  EXPECT_EQ(session_state->get_monitor_pool().get_credit("m3", USED_RX), 3000);
  EXPECT_EQ(update_criteria.monitor_credit_map["m3"].bucket_deltas[USED_TX],
            5000);

  UpdateSessionRequest update;
  std::vector<std::unique_ptr<ServiceAction>> actions;
  session_state->get_updates(update, &actions, update_criteria);
  EXPECT_EQ(actions.size(), 0);
  EXPECT_EQ(update.updates_size(), 2);
  EXPECT_EQ(update.usage_monitors_size(), 2);
}

TEST_F(SessionStateTest, test_session_level_key) {
  EXPECT_EQ(nullptr, session_state->get_monitor_pool().get_session_level_key());

  receive_credit_from_pcrf("m1", 8000, MonitoringLevel::SESSION_LEVEL);
  EXPECT_EQ("m1", *session_state->get_monitor_pool().get_session_level_key());
  EXPECT_EQ(session_state->get_monitor_pool().get_credit("m1", ALLOWED_TOTAL),
            8000);
  EXPECT_EQ(update_criteria.monitor_credit_to_install["m1"]
                .credit.buckets[ALLOWED_TOTAL],
            8000);

  session_state->add_used_credit("rule1", 5000, 3000, update_criteria);
  EXPECT_EQ(session_state->get_monitor_pool().get_credit("m1", USED_TX), 5000);
  EXPECT_EQ(session_state->get_monitor_pool().get_credit("m1", USED_RX), 3000);
  EXPECT_EQ(update_criteria.monitor_credit_map["m1"].bucket_deltas[USED_TX],
            5000);
  EXPECT_EQ(update_criteria.monitor_credit_map["m1"].bucket_deltas[USED_RX],
            3000);

  UpdateSessionRequest update;
  std::vector<std::unique_ptr<ServiceAction>> actions;
  session_state->get_updates(update, &actions, update_criteria);
  EXPECT_EQ(actions.size(), 0);
  EXPECT_EQ(update.usage_monitors_size(), 1);
  auto &single_update = update.usage_monitors(0).update();
  EXPECT_EQ(single_update.level(), MonitoringLevel::SESSION_LEVEL);
  EXPECT_EQ(single_update.bytes_rx(), 3000);
  EXPECT_EQ(single_update.bytes_tx(), 5000);
}

TEST_F(SessionStateTest, test_reauth_key) {
  insert_rule(1, "", "rule1", STATIC, 0, 0);

  receive_credit_from_ocs(1, 1500);

  session_state->add_used_credit("rule1", 1000, 500, update_criteria);

  UpdateSessionRequest update;
  std::vector<std::unique_ptr<ServiceAction>> actions;
  session_state->get_updates(update, &actions, update_criteria);
  EXPECT_EQ(update.updates_size(), 1);
  EXPECT_EQ(session_state->get_charging_pool().get_credit(1, REPORTING_TX),
            1000);
  EXPECT_EQ(session_state->get_charging_pool().get_credit(1, REPORTING_RX),
            500);
  // Reporting value is not tracked by UpdateCriteria
  EXPECT_EQ(update_criteria.charging_credit_map[CreditKey(1)]
                .bucket_deltas[REPORTING_TX],
            0);
  EXPECT_EQ(update_criteria.charging_credit_map[CreditKey(1)]
                .bucket_deltas[REPORTING_RX],
            0);
  // credit is already reporting, no update needed
  auto uc = get_default_update_criteria();
  auto reauth_res = session_state->get_charging_pool().reauth_key(1, uc);
  EXPECT_EQ(reauth_res, ChargingReAuthAnswer::UPDATE_NOT_NEEDED);
  receive_credit_from_ocs(1, 1024);
  EXPECT_EQ(session_state->get_charging_pool().get_credit(1, REPORTING_TX), 0);
  EXPECT_EQ(session_state->get_charging_pool().get_credit(1, REPORTING_RX), 0);
  reauth_res = session_state->get_charging_pool().reauth_key(1, uc);
  EXPECT_EQ(reauth_res, ChargingReAuthAnswer::UPDATE_INITIATED);

  session_state->add_used_credit("rule1", 2, 1, update_criteria);
  UpdateSessionRequest reauth_update;
  session_state->get_updates(reauth_update, &actions, update_criteria);
  EXPECT_EQ(reauth_update.updates_size(), 1);
  auto &usage = reauth_update.updates(0).usage();
  EXPECT_EQ(usage.bytes_tx(), 2);
  EXPECT_EQ(usage.bytes_rx(), 1);
}

TEST_F(SessionStateTest, test_reauth_new_key) {
  // credit is already reporting, no update needed
  auto reauth_res =
      session_state->get_charging_pool().reauth_key(1, update_criteria);
  EXPECT_EQ(reauth_res, ChargingReAuthAnswer::UPDATE_INITIATED);

  UpdateSessionRequest reauth_update;
  std::vector<std::unique_ptr<ServiceAction>> actions;
  session_state->get_updates(reauth_update, &actions, update_criteria);
  EXPECT_EQ(reauth_update.updates_size(), 1);
  auto &usage = reauth_update.updates(0).usage();
  EXPECT_EQ(usage.charging_key(), 1);
  EXPECT_EQ(usage.bytes_tx(), 0);
  EXPECT_EQ(usage.bytes_rx(), 0);

  receive_credit_from_ocs(1, 1024);
  EXPECT_EQ(session_state->get_charging_pool().get_credit(1, ALLOWED_TOTAL),
            1024);
  EXPECT_EQ(update_criteria.charging_credit_map[CreditKey(1)]
                .bucket_deltas[ALLOWED_TOTAL],
            1024);
}

TEST_F(SessionStateTest, test_reauth_all) {
  insert_rule(1, "", "rule1", STATIC, 0, 0);
  insert_rule(2, "", "dyn_rule1", DYNAMIC, 0, 0);
  EXPECT_EQ(false, session_state->active_monitored_rules_exist());
  EXPECT_TRUE(std::find(update_criteria.static_rules_to_install.begin(),
                        update_criteria.static_rules_to_install.end(),
                        "rule1") !=
              update_criteria.static_rules_to_install.end());
  EXPECT_EQ(update_criteria.dynamic_rules_to_install.size(), 1);

  receive_credit_from_ocs(1, 1024);
  receive_credit_from_ocs(2, 1024);

  session_state->add_used_credit("rule1", 10, 20, update_criteria);
  session_state->add_used_credit("dyn_rule1", 30, 40, update_criteria);
  // If any charging key isn't reporting, an update is needed
  auto uc = get_default_update_criteria();
  auto reauth_res = session_state->get_charging_pool().reauth_all(uc);
  EXPECT_EQ(reauth_res, ChargingReAuthAnswer::UPDATE_INITIATED);

  UpdateSessionRequest reauth_update;
  std::vector<std::unique_ptr<ServiceAction>> actions;
  session_state->get_updates(reauth_update, &actions, update_criteria);
  EXPECT_EQ(reauth_update.updates_size(), 2);

  // All charging keys are reporting, no update needed
  reauth_res = session_state->get_charging_pool().reauth_all(uc);
  EXPECT_EQ(reauth_res, ChargingReAuthAnswer::UPDATE_NOT_NEEDED);
}

TEST_F(SessionStateTest, test_tgpp_context_is_set_on_update) {
  receive_credit_from_pcrf("m1", 1024, MonitoringLevel::PCC_RULE_LEVEL);
  receive_credit_from_ocs(1, 1024);
  insert_rule(1, "m1", "rule1", STATIC, 0, 0);
  session_state->add_used_credit("rule1", 1024, 0, update_criteria);
  EXPECT_EQ(true, session_state->active_monitored_rules_exist());

  EXPECT_EQ(session_state->get_monitor_pool().get_credit("m1", ALLOWED_TOTAL),
            1024);
  EXPECT_EQ(session_state->get_charging_pool().get_credit(1, ALLOWED_TOTAL),
            1024);
  EXPECT_EQ(update_criteria.charging_credit_to_install[CreditKey(1)]
                .buckets[ALLOWED_TOTAL],
            1024);
  EXPECT_EQ(update_criteria.monitor_credit_to_install["m1"]
                .credit.buckets[ALLOWED_TOTAL],
            1024);

  UpdateSessionRequest update;
  std::vector<std::unique_ptr<ServiceAction>> actions;
  session_state->get_updates(update, &actions, update_criteria);
  EXPECT_EQ(actions.size(), 0);
  EXPECT_EQ(update.updates_size(), 1);
  EXPECT_EQ(update.updates().Get(0).tgpp_ctx().gx_dest_host(), "gx.dest.com");
  EXPECT_EQ(update.updates().Get(0).tgpp_ctx().gy_dest_host(), "gy.dest.com");
  EXPECT_EQ(update.usage_monitors_size(), 1);
  EXPECT_EQ(update.usage_monitors().Get(0).tgpp_ctx().gx_dest_host(),
            "gx.dest.com");
  EXPECT_EQ(update.usage_monitors().Get(0).tgpp_ctx().gy_dest_host(),
            "gy.dest.com");
}

TEST_F(SessionStateTest, test_get_total_credit_usage_single_rule_no_key) {
  insert_rule(0, "", "rule1", STATIC, 0, 0);
  session_state->add_used_credit("rule1", 2000, 1000, update_criteria);
  SessionState::TotalCreditUsage actual =
      session_state->get_total_credit_usage();
  EXPECT_EQ(actual.monitoring_tx, 0);
  EXPECT_EQ(actual.monitoring_rx, 0);
  EXPECT_EQ(actual.charging_tx, 0);
  EXPECT_EQ(actual.charging_rx, 0);
}

TEST_F(SessionStateTest, test_get_total_credit_usage_single_rule_single_key) {
  insert_rule(1, "", "rule1", STATIC, 0, 0);
  receive_credit_from_ocs(1, 3000);
  session_state->add_used_credit("rule1", 2000, 1000, update_criteria);
  SessionState::TotalCreditUsage actual =
      session_state->get_total_credit_usage();
  EXPECT_EQ(actual.monitoring_tx, 0);
  EXPECT_EQ(actual.monitoring_rx, 0);
  EXPECT_EQ(actual.charging_tx, 2000);
  EXPECT_EQ(actual.charging_rx, 1000);
}

TEST_F(SessionStateTest, test_get_total_credit_usage_single_rule_multiple_key) {
  insert_rule(1, "m1", "rule1", STATIC, 0, 0);
  receive_credit_from_ocs(1, 3000);
  receive_credit_from_pcrf("m1", 3000, MonitoringLevel::PCC_RULE_LEVEL);
  session_state->add_used_credit("rule1", 2000, 1000, update_criteria);
  SessionState::TotalCreditUsage actual =
      session_state->get_total_credit_usage();
  EXPECT_EQ(actual.monitoring_tx, 2000);
  EXPECT_EQ(actual.monitoring_rx, 1000);
  EXPECT_EQ(actual.charging_tx, 2000);
  EXPECT_EQ(actual.charging_rx, 1000);
}

TEST_F(SessionStateTest, test_get_total_credit_usage_multiple_rule_shared_key) {
  // Shared monitoring key
  // One rule is dynamic
  insert_rule(1, "m1", "rule1", STATIC, 0, 0);
  insert_rule(0, "m1", "rule2", DYNAMIC, 0, 0);
  receive_credit_from_ocs(1, 3000);
  receive_credit_from_pcrf("m1", 3000, MonitoringLevel::PCC_RULE_LEVEL);
  session_state->add_used_credit("rule1", 1000, 10, update_criteria);
  session_state->add_used_credit("rule2", 500, 5, update_criteria);
  SessionState::TotalCreditUsage actual =
      session_state->get_total_credit_usage();
  EXPECT_EQ(actual.monitoring_tx, 1500);
  EXPECT_EQ(actual.monitoring_rx, 15);
  EXPECT_EQ(actual.charging_tx, 1000);
  EXPECT_EQ(actual.charging_rx, 10);
}

int main(int argc, char **argv) {
  ::testing::InitGoogleTest(&argc, argv);
  FLAGS_logtostderr = 1;
  FLAGS_v = 10;
  return RUN_ALL_TESTS();
}

} // namespace magma

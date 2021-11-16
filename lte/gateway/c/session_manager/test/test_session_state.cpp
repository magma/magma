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

TEST_F(SessionStateTest, test_session_rules) {
  activate_rule(1, "m1", "rule1", DYNAMIC, 0, 0);
  EXPECT_EQ(1, session_state->total_monitored_rules_count());
  activate_rule(2, "m2", "rule2", STATIC, 0, 0);
  EXPECT_EQ(2, session_state->total_monitored_rules_count());
  // add a OCS-ONLY static rule
  activate_rule(3, "", "rule3", STATIC, 0, 0);
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
  session_state->deactivate_static_rule("rule2", &update_criteria);
  EXPECT_EQ(1, session_state->total_monitored_rules_count());
  EXPECT_TRUE(
      session_state->remove_dynamic_rule("rule1", &rule_out, &update_criteria));
  EXPECT_EQ("m1", rule_out.monitoring_key());
  EXPECT_EQ(0, session_state->total_monitored_rules_count());

  // basic sanity checks to see it's properly deleted
  rules_out = {};
  session_state->get_dynamic_rules().get_rule_ids(rules_out_ptr);
  EXPECT_EQ(rules_out_ptr.size(), 0);

  EXPECT_EQ(0, get_monitored_rule_count("m1"));

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
  PolicyRule rule;
  session_state->remove_scheduled_dynamic_rule("rule1", &rule, nullptr);
  session_state->insert_dynamic_rule(
      rule, session_state->get_rule_lifetime("rule1"), nullptr);
  EXPECT_EQ(1, session_state->total_monitored_rules_count());
  EXPECT_TRUE(session_state->is_dynamic_rule_installed("rule1"));

  session_state->activate_static_rule(
      "rule2", session_state->get_rule_lifetime("rule2"), nullptr);
  EXPECT_EQ(2, session_state->total_monitored_rules_count());
  EXPECT_TRUE(session_state->is_static_rule_installed("rule2"));

  EXPECT_EQ(session_state->get_current_rule_version("rule1"), 1);
  EXPECT_EQ(session_state->get_current_rule_version("rule2"), 1);
}

/**
 * Check that on restart, sessions can be updated to match the current time
 */
TEST_F(SessionStateTest, test_rule_time_sync) {
  auto uc = get_default_update_criteria();  // unused

  // These should be active after sync
  schedule_rule(1, "m1", "d1", DYNAMIC, 5, 15);
  schedule_rule(1, "m1", "s1", STATIC, 5, 15);

  // These should still be scheduled
  schedule_rule(1, "m1", "d2", DYNAMIC, 15, 20);
  schedule_rule(1, "m1", "s2", STATIC, 15, 20);

  // These should be expired afterwards
  schedule_rule(2, "m2", "d3", DYNAMIC, 2, 4);
  schedule_rule(2, "m2", "s3", STATIC, 2, 4);

  EXPECT_FALSE(session_state->is_dynamic_rule_installed("d1"));
  EXPECT_FALSE(session_state->is_dynamic_rule_installed("d2"));
  EXPECT_FALSE(session_state->is_dynamic_rule_installed("d3"));

  EXPECT_FALSE(session_state->is_static_rule_installed("s1"));
  EXPECT_FALSE(session_state->is_static_rule_installed("s2"));
  EXPECT_FALSE(session_state->is_static_rule_installed("s3"));

  // Update the time, and sync the rule states, then check our expectations
  std::time_t test_time(10);
  session_state->sync_rules_to_time(test_time, &uc);

  EXPECT_TRUE(session_state->is_dynamic_rule_installed("d1"));
  EXPECT_FALSE(session_state->is_dynamic_rule_installed("d2"));
  EXPECT_FALSE(session_state->is_dynamic_rule_installed("d3"));

  EXPECT_TRUE(session_state->is_static_rule_installed("s1"));
  EXPECT_FALSE(session_state->is_static_rule_installed("s2"));
  EXPECT_FALSE(session_state->is_static_rule_installed("s3"));

  EXPECT_EQ(uc.dynamic_rules_to_install.size(), 1);
  EXPECT_EQ(uc.dynamic_rules_to_install.front().id(), "d1");
  EXPECT_TRUE(uc.dynamic_rules_to_uninstall.count("d3"));

  EXPECT_TRUE(uc.static_rules_to_install.count("s1"));

  // Update the time once more, sync again, and check expectations
  test_time = std::time_t(16);
  uc = get_default_update_criteria();
  session_state->sync_rules_to_time(test_time, &uc);

  EXPECT_FALSE(session_state->is_dynamic_rule_installed("d1"));
  EXPECT_TRUE(session_state->is_dynamic_rule_installed("d2"));
  EXPECT_FALSE(session_state->is_dynamic_rule_installed("d3"));

  EXPECT_FALSE(session_state->is_static_rule_installed("s1"));
  EXPECT_TRUE(session_state->is_static_rule_installed("s2"));
  EXPECT_FALSE(session_state->is_static_rule_installed("s3"));

  EXPECT_EQ(uc.dynamic_rules_to_install.size(), 1);
  EXPECT_EQ(uc.dynamic_rules_to_install.front().id(), "d2");
  EXPECT_TRUE(uc.dynamic_rules_to_uninstall.count("d1"));

  EXPECT_TRUE(uc.static_rules_to_install.count("s2"));
  EXPECT_TRUE(uc.static_rules_to_uninstall.count("s1"));
}

TEST_F(SessionStateTest, test_marshal_unmarshal) {
  EXPECT_EQ(update_criteria.static_rules_to_install.size(), 0);
  insert_rule(1, "m1", "rule1", STATIC, 0, 0);
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

  EXPECT_EQ(update_criteria.charging_credit_to_install.size(), 0);
  receive_credit_from_ocs(1, 1024);
  EXPECT_EQ(update_criteria.charging_credit_to_install.size(), 1);
  EXPECT_EQ(session_state->get_charging_credit(1, ALLOWED_TOTAL), 1024);

  EXPECT_EQ(update_criteria.monitor_credit_to_install.size(), 0);
  receive_credit_from_pcrf("m1", 1024, MonitoringLevel::PCC_RULE_LEVEL);
  EXPECT_EQ(session_state->get_monitor("m1", ALLOWED_TOTAL), 1024);
  EXPECT_EQ(update_criteria.monitor_credit_to_install.size(), 1);

  session_state->add_rule_usage("rule1", 1, 2000, 1000, 0, 0, nullptr);
  EXPECT_EQ(session_state->get_charging_credit(1, USED_TX), 2000);
  EXPECT_EQ(session_state->get_charging_credit(1, USED_RX), 1000);
  EXPECT_EQ(session_state->get_monitor("m1", USED_TX), 2000);
  EXPECT_EQ(session_state->get_monitor("m1", USED_RX), 1000);
  session_state->increment_rule_stats("rule1", nullptr);
  session_state->add_rule_usage("rule1", 2, 1000, 500, 10, 20, nullptr);
  EXPECT_EQ(session_state->get_charging_credit(1, USED_TX), 3000);
  EXPECT_EQ(session_state->get_charging_credit(1, USED_RX), 1500);
  EXPECT_EQ(session_state->get_monitor("m1", USED_TX), 3000);
  EXPECT_EQ(session_state->get_monitor("m1", USED_RX), 1500);

  auto marshaled = session_state->marshal();
  auto unmarshaled = SessionState::unmarshal(marshaled, *rule_store);
  EXPECT_EQ(unmarshaled->get_charging_credit(1, ALLOWED_TOTAL), 1024);
  EXPECT_EQ(unmarshaled->get_monitor("m1", ALLOWED_TOTAL), 1024);
  EXPECT_EQ(unmarshaled->is_static_rule_installed("rule1"), true);
  EXPECT_EQ(unmarshaled->is_dynamic_rule_installed("rule2"), false);
  EXPECT_EQ(unmarshaled->get_policy_stats("rule1").last_reported_version, 2);
  EXPECT_EQ(unmarshaled->get_policy_stats("rule1").stats_map[1].tx, 2000);
  EXPECT_EQ(unmarshaled->get_policy_stats("rule1").stats_map[2].tx, 1000);
  EXPECT_EQ(unmarshaled->get_policy_stats("rule1").stats_map[2].dropped_rx, 20);
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
  EXPECT_EQ(session_state->get_charging_credit(1, ALLOWED_TOTAL), 1024);
  EXPECT_EQ(update_criteria.charging_credit_to_install[CreditKey(1)]
                .credit.buckets[ALLOWED_TOTAL],
            1024);

  receive_credit_from_pcrf("m1", 1024, MonitoringLevel::PCC_RULE_LEVEL);
  EXPECT_EQ(session_state->get_monitor("m1", ALLOWED_TOTAL), 1024);
  EXPECT_EQ(update_criteria.monitor_credit_to_install["m1"]
                .credit.buckets[ALLOWED_TOTAL],
            1024);
}

TEST_F(SessionStateTest, test_add_rule_usage) {
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
                .credit.buckets[ALLOWED_TOTAL],
            3000);

  receive_credit_from_pcrf("m1", 3000, MonitoringLevel::PCC_RULE_LEVEL);
  receive_credit_from_pcrf("m2", 6000, MonitoringLevel::PCC_RULE_LEVEL);
  EXPECT_EQ(update_criteria.monitor_credit_to_install.size(), 2);
  EXPECT_EQ(update_criteria.monitor_credit_to_install["m1"]
                .credit.buckets[ALLOWED_TOTAL],
            3000);

  session_state->add_rule_usage("rule1", 1, 2000, 1000, 0, 0, &update_criteria);
  EXPECT_EQ(session_state->get_charging_credit(1, USED_TX), 2000);
  EXPECT_EQ(session_state->get_charging_credit(1, USED_RX), 1000);
  EXPECT_EQ(session_state->get_monitor("m1", USED_TX), 2000);
  EXPECT_EQ(session_state->get_monitor("m1", USED_RX), 1000);
  EXPECT_EQ(
      update_criteria.charging_credit_map[CreditKey(1)].bucket_deltas[USED_TX],
      2000);
  EXPECT_EQ(update_criteria.monitor_credit_map["m1"].bucket_deltas[USED_RX],
            1000);

  session_state->add_rule_usage("dyn_rule1", 1, 4000, 2000, 0, 0,
                                &update_criteria);
  EXPECT_EQ(session_state->get_charging_credit(2, USED_TX), 4000);
  EXPECT_EQ(session_state->get_charging_credit(2, USED_RX), 2000);
  EXPECT_EQ(session_state->get_monitor("m2", USED_TX), 4000);
  EXPECT_EQ(session_state->get_monitor("m2", USED_RX), 2000);
  EXPECT_EQ(
      update_criteria.charging_credit_map[CreditKey(2)].bucket_deltas[USED_TX],
      4000);
  EXPECT_EQ(update_criteria.monitor_credit_map["m2"].bucket_deltas[USED_RX],
            2000);

  UpdateSessionRequest update;
  std::vector<std::unique_ptr<ServiceAction>> actions;
  session_state->get_updates(&update, &actions, &update_criteria);
  EXPECT_EQ(actions.size(), 0);
  EXPECT_EQ(update.updates_size(), 2);
  EXPECT_EQ(update.usage_monitors_size(), 2);

  PolicyRule policy_out;
  EXPECT_TRUE(session_state->remove_dynamic_rule("dyn_rule1", &policy_out,
                                                 &update_criteria));
  EXPECT_TRUE(session_state->deactivate_static_rule("rule1", &update_criteria));
  EXPECT_FALSE(session_state->active_monitored_rules_exist());
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
                .credit.buckets[ALLOWED_TOTAL],
            6000);

  receive_credit_from_pcrf("m1", 3000, MonitoringLevel::PCC_RULE_LEVEL);
  receive_credit_from_pcrf("m3", 8000, MonitoringLevel::PCC_RULE_LEVEL);
  EXPECT_EQ(update_criteria.monitor_credit_to_install["m1"]
                .credit.buckets[ALLOWED_TOTAL],
            3000);

  session_state->add_rule_usage("dyn_rule1", 1, 2000, 1000, 0, 0,
                                &update_criteria);
  EXPECT_EQ(session_state->get_monitor("m1", USED_TX), 2000);
  EXPECT_EQ(session_state->get_monitor("m1", USED_RX), 1000);

  EXPECT_EQ(update_criteria.monitor_credit_map["m1"].bucket_deltas[USED_TX],
            2000);
  EXPECT_EQ(update_criteria.monitor_credit_map["m1"].bucket_deltas[USED_RX],
            1000);

  session_state->add_rule_usage("dyn_rule2", 1, 4000, 2000, 0, 0,
                                &update_criteria);
  EXPECT_EQ(session_state->get_charging_credit(2, USED_TX), 4000);
  EXPECT_EQ(session_state->get_charging_credit(2, USED_RX), 2000);
  EXPECT_EQ(
      update_criteria.charging_credit_map[CreditKey(2)].bucket_deltas[USED_TX],
      4000);
  session_state->add_rule_usage("dyn_rule3", 1, 5000, 3000, 0, 0,
                                &update_criteria);
  EXPECT_EQ(session_state->get_charging_credit(3, USED_TX), 5000);
  EXPECT_EQ(session_state->get_charging_credit(3, USED_RX), 3000);
  EXPECT_EQ(
      update_criteria.charging_credit_map[CreditKey(3)].bucket_deltas[USED_TX],
      5000);
  EXPECT_EQ(session_state->get_monitor("m3", USED_TX), 5000);
  EXPECT_EQ(session_state->get_monitor("m3", USED_RX), 3000);
  EXPECT_EQ(update_criteria.monitor_credit_map["m3"].bucket_deltas[USED_TX],
            5000);

  UpdateSessionRequest update;
  std::vector<std::unique_ptr<ServiceAction>> actions;
  session_state->get_updates(&update, &actions, nullptr);
  EXPECT_EQ(actions.size(), 0);
  EXPECT_EQ(update.updates_size(), 2);
  EXPECT_EQ(update.usage_monitors_size(), 2);
}

TEST_F(SessionStateTest, test_session_level_key) {
  insert_rule(1, "m1", "rule1", DYNAMIC, 0, 0);
  receive_credit_from_pcrf("m1", 8000, MonitoringLevel::SESSION_LEVEL);
  EXPECT_EQ(session_state->get_monitor("m1", ALLOWED_TOTAL), 8000);
  EXPECT_EQ(update_criteria.monitor_credit_to_install["m1"]
                .credit.buckets[ALLOWED_TOTAL],
            8000);
  EXPECT_TRUE(update_criteria.is_session_level_key_updated);
  EXPECT_EQ(update_criteria.updated_session_level_key, "m1");

  // add usage to go over quota
  session_state->add_rule_usage("rule1", 1, 5000, 2000, 0, 0, &update_criteria);
  EXPECT_EQ(session_state->get_policy_stats("rule1").current_version, 1);
  EXPECT_EQ(session_state->get_policy_stats("rule1").last_reported_version, 1);
  EXPECT_EQ(session_state->get_monitor("m1", USED_TX), 5000);
  EXPECT_EQ(session_state->get_monitor("m1", USED_RX), 2000);
  EXPECT_EQ(update_criteria.monitor_credit_map["m1"].bucket_deltas[USED_TX],
            5000);
  EXPECT_EQ(update_criteria.monitor_credit_map["m1"].bucket_deltas[USED_RX],
            2000);

  // check one updates will be sent
  UpdateSessionRequest update;
  std::vector<std::unique_ptr<ServiceAction>> actions;
  session_state->get_updates(&update, &actions, &update_criteria);
  EXPECT_EQ(actions.size(), 0);
  EXPECT_EQ(update.usage_monitors_size(), 1);
  auto& single_update = update.usage_monitors(0).update();
  EXPECT_EQ(single_update.level(), MonitoringLevel::SESSION_LEVEL);
  EXPECT_EQ(single_update.bytes_tx(), 5000);
  EXPECT_EQ(single_update.bytes_rx(), 2000);

  // Send 0 value traffic which will indicate monitor must be disabled after
  // going out of quota
  receive_credit_from_pcrf("m1", 0, MonitoringLevel::SESSION_LEVEL);

  // add usage to go over quota
  session_state->add_rule_usage("rule1", 1, 6001, 2001, 0, 0, &update_criteria);
  EXPECT_EQ(session_state->get_monitor("m1", USED_TX), 6001);
  EXPECT_EQ(session_state->get_monitor("m1", USED_RX), 2001);
  EXPECT_TRUE(update_criteria.monitor_credit_map["m1"].report_last_credit);

  // check final update will be sent
  UpdateSessionRequest update_2;
  std::vector<std::unique_ptr<ServiceAction>> actions_2;
  session_state->get_updates(&update_2, &actions_2, &update_criteria);
  // TODO: session level seemsd to be adding total values, no deltas
  // EXPECT_EQ(
  //   update_criteria.monitor_credit_map["m1"].bucket_deltas[USED_TX], 1001);
  // EXPECT_EQ(
  //    update_criteria.monitor_credit_map["m1"].bucket_deltas[USED_RX], 1);
  EXPECT_EQ(actions_2.size(), 0);
  EXPECT_EQ(update_2.updates_size(), 0);
  EXPECT_EQ(update_2.usage_monitors_size(), 1);
  auto& single_update_2 = update_2.usage_monitors(0).update();
  EXPECT_EQ(single_update_2.level(), MonitoringLevel::SESSION_LEVEL);
  // Substract the values from first usage report
  EXPECT_EQ(single_update_2.bytes_rx(), 1);
  EXPECT_EQ(single_update_2.bytes_tx(), 1001);

  // apply updates (prepare the session to be merged into storage)
  // and check monitor has been deleted (=0)
  session_state->apply_update_criteria(update_criteria);
  EXPECT_EQ(session_state->get_monitor("m1", USED_TX), 0);
  EXPECT_EQ(session_state->get_monitor("m1", USED_RX), 0);
  EXPECT_EQ(session_state->get_session_level_key(), "");
  EXPECT_TRUE(update_criteria.monitor_credit_map["m1"].deleted);
}

TEST_F(SessionStateTest, test_reauth_key) {
  insert_rule(1, "", "rule1", STATIC, 0, 0);

  receive_credit_from_ocs(1, 1500);

  session_state->add_rule_usage("rule1", 1, 1000, 500, 0, 0, &update_criteria);

  UpdateSessionRequest update;
  std::vector<std::unique_ptr<ServiceAction>> actions;
  session_state->get_updates(&update, &actions, &update_criteria);
  EXPECT_EQ(update.updates_size(), 1);
  EXPECT_EQ(session_state->get_charging_credit(1, REPORTING_TX), 1000);
  EXPECT_EQ(session_state->get_charging_credit(1, REPORTING_RX), 500);
  // Reporting value is not tracked by UpdateCriteria
  EXPECT_EQ(update_criteria.charging_credit_map[CreditKey(1)]
                .bucket_deltas[REPORTING_TX],
            0);
  EXPECT_EQ(update_criteria.charging_credit_map[CreditKey(1)]
                .bucket_deltas[REPORTING_RX],
            0);
  // credit is already reporting, no update needed
  auto uc = get_default_update_criteria();
  auto reauth_res = session_state->reauth_key(1, &uc);
  EXPECT_EQ(reauth_res, ReAuthResult::UPDATE_NOT_NEEDED);
  receive_credit_from_ocs(1, 1024);
  EXPECT_EQ(session_state->get_charging_credit(1, REPORTING_TX), 0);
  EXPECT_EQ(session_state->get_charging_credit(1, REPORTING_RX), 0);
  reauth_res = session_state->reauth_key(1, &uc);
  EXPECT_EQ(reauth_res, ReAuthResult::UPDATE_INITIATED);

  session_state->add_rule_usage("rule1", 1, 1002, 501, 0, 0, &update_criteria);
  UpdateSessionRequest reauth_update;
  session_state->get_updates(&reauth_update, &actions, nullptr);
  EXPECT_EQ(reauth_update.updates_size(), 1);
  auto& usage = reauth_update.updates(0).usage();
  EXPECT_EQ(usage.bytes_tx(), 2);
  EXPECT_EQ(usage.bytes_rx(), 1);
}

TEST_F(SessionStateTest, test_reauth_new_key) {
  // credit is already reporting, no update needed
  auto reauth_res = session_state->reauth_key(1, &update_criteria);
  EXPECT_EQ(reauth_res, ReAuthResult::UPDATE_INITIATED);

  // assert stored charging grant fields are updated to reflect reauth state
  EXPECT_EQ(update_criteria.charging_credit_to_install.size(), 1);
  auto new_charging_credits = update_criteria.charging_credit_to_install;
  EXPECT_EQ(new_charging_credits[1].reauth_state, REAUTH_REQUIRED);

  UpdateSessionRequest reauth_update;
  std::vector<std::unique_ptr<ServiceAction>> actions;
  session_state->get_updates(&reauth_update, &actions, &update_criteria);
  EXPECT_EQ(reauth_update.updates_size(), 1);
  auto& usage = reauth_update.updates(0).usage();
  EXPECT_EQ(usage.charging_key(), 1);
  EXPECT_EQ(usage.bytes_tx(), 0);
  EXPECT_EQ(usage.bytes_rx(), 0);

  // assert stored charging grant fields are updated to reflect reauth state
  EXPECT_EQ(update_criteria.charging_credit_map.size(), 1);
  auto credit_uc = update_criteria.charging_credit_map[1];
  EXPECT_EQ(credit_uc.reauth_state, REAUTH_PROCESSING);

  receive_credit_from_ocs(1, 1024);
  EXPECT_EQ(session_state->get_charging_credit(1, ALLOWED_TOTAL), 1024);
  EXPECT_EQ(update_criteria.charging_credit_map[CreditKey(1)]
                .bucket_deltas[ALLOWED_TOTAL],
            1024);

  // assert stored charging grant fields are updated to reflect reauth state
  EXPECT_EQ(update_criteria.charging_credit_map.size(), 1);
  credit_uc = update_criteria.charging_credit_map[1];
  EXPECT_EQ(credit_uc.reauth_state, REAUTH_NOT_NEEDED);
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

  session_state->add_rule_usage("rule1", 1, 10, 20, 0, 0, &update_criteria);
  session_state->add_rule_usage("dyn_rule1", 1, 30, 40, 0, 0, &update_criteria);
  // If any charging key isn't reporting, an update is needed
  auto uc = get_default_update_criteria();
  auto reauth_res = session_state->reauth_all(&uc);
  EXPECT_EQ(reauth_res, ReAuthResult::UPDATE_INITIATED);

  // assert stored charging grant fields are updated to reflect reauth state
  EXPECT_EQ(uc.charging_credit_map.size(), 2);
  auto credit_uc_1 = uc.charging_credit_map[1];
  auto credit_uc_2 = uc.charging_credit_map[2];
  EXPECT_EQ(credit_uc_1.reauth_state, REAUTH_REQUIRED);
  EXPECT_EQ(credit_uc_2.reauth_state, REAUTH_REQUIRED);

  UpdateSessionRequest reauth_update;
  std::vector<std::unique_ptr<ServiceAction>> actions;
  session_state->get_updates(&reauth_update, &actions, &uc);
  EXPECT_EQ(reauth_update.updates_size(), 2);

  EXPECT_EQ(uc.charging_credit_map.size(), 2);
  credit_uc_1 = uc.charging_credit_map[1];
  credit_uc_2 = uc.charging_credit_map[2];
  EXPECT_EQ(credit_uc_1.reauth_state, REAUTH_PROCESSING);
  EXPECT_EQ(credit_uc_2.reauth_state, REAUTH_PROCESSING);

  // All charging keys are reporting, no update needed
  reauth_res = session_state->reauth_all(&uc);
  EXPECT_EQ(reauth_res, ReAuthResult::UPDATE_NOT_NEEDED);

  EXPECT_EQ(uc.charging_credit_map.size(), 2);
  credit_uc_1 = uc.charging_credit_map[1];
  credit_uc_2 = uc.charging_credit_map[2];
  EXPECT_EQ(credit_uc_1.reauth_state, REAUTH_PROCESSING);
  EXPECT_EQ(credit_uc_2.reauth_state, REAUTH_PROCESSING);
}

TEST_F(SessionStateTest, test_tgpp_context_is_set_on_update) {
  receive_credit_from_pcrf("m1", 1024, MonitoringLevel::PCC_RULE_LEVEL);
  receive_credit_from_ocs(1, 1024);
  insert_rule(1, "m1", "rule1", STATIC, 0, 0);
  session_state->add_rule_usage("rule1", 1, 1024, 0, 0, 0, &update_criteria);
  EXPECT_EQ(true, session_state->active_monitored_rules_exist());

  EXPECT_EQ(session_state->get_monitor("m1", ALLOWED_TOTAL), 1024);
  EXPECT_EQ(session_state->get_charging_credit(1, ALLOWED_TOTAL), 1024);
  EXPECT_EQ(update_criteria.charging_credit_to_install[CreditKey(1)]
                .credit.buckets[ALLOWED_TOTAL],
            1024);
  EXPECT_EQ(update_criteria.monitor_credit_to_install["m1"]
                .credit.buckets[ALLOWED_TOTAL],
            1024);

  UpdateSessionRequest update;
  std::vector<std::unique_ptr<ServiceAction>> actions;
  session_state->get_updates(&update, &actions, nullptr);
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
  session_state->add_rule_usage("rule1", 1, 2000, 1000, 0, 0, &update_criteria);
  TotalCreditUsage actual = session_state->get_total_credit_usage();
  EXPECT_EQ(actual.monitoring_tx, 0);
  EXPECT_EQ(actual.monitoring_rx, 0);
  EXPECT_EQ(actual.charging_tx, 0);
  EXPECT_EQ(actual.charging_rx, 0);
}

TEST_F(SessionStateTest, test_get_total_credit_usage_single_rule_single_key) {
  insert_rule(1, "", "rule1", STATIC, 0, 0);
  receive_credit_from_ocs(1, 3000);
  session_state->add_rule_usage("rule1", 1, 2000, 1000, 0, 0, &update_criteria);
  TotalCreditUsage actual = session_state->get_total_credit_usage();
  EXPECT_EQ(actual.monitoring_tx, 0);
  EXPECT_EQ(actual.monitoring_rx, 0);
  EXPECT_EQ(actual.charging_tx, 2000);
  EXPECT_EQ(actual.charging_rx, 1000);
}

TEST_F(SessionStateTest, test_get_total_credit_usage_single_rule_multiple_key) {
  insert_rule(1, "m1", "rule1", STATIC, 0, 0);
  receive_credit_from_ocs(1, 3000);
  receive_credit_from_pcrf("m1", 3000, MonitoringLevel::PCC_RULE_LEVEL);
  session_state->add_rule_usage("rule1", 1, 2000, 1000, 0, 0, &update_criteria);
  TotalCreditUsage actual = session_state->get_total_credit_usage();
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
  session_state->add_rule_usage("rule1", 1, 1000, 10, 0, 0, nullptr);
  session_state->add_rule_usage("rule2", 1, 500, 5, 0, 0, nullptr);
  TotalCreditUsage actual = session_state->get_total_credit_usage();
  EXPECT_EQ(actual.monitoring_tx, 1500);
  EXPECT_EQ(actual.monitoring_rx, 15);
  EXPECT_EQ(actual.charging_tx, 1000);
  EXPECT_EQ(actual.charging_rx, 10);
}

TEST_F(SessionStateTest, test_install_gy_rules) {
  uint32_t version = insert_gy_redirection_rule("redirect");
  EXPECT_EQ(1, version);

  std::vector<std::string> rules_out{};
  std::vector<std::string>& rules_out_ptr = rules_out;

  session_state->get_gy_dynamic_rules().get_rule_ids(rules_out_ptr);
  EXPECT_EQ(rules_out_ptr.size(), 1);
  EXPECT_EQ(rules_out_ptr[0], "redirect");

  EXPECT_TRUE(session_state->is_gy_dynamic_rule_installed("redirect"));
  EXPECT_EQ(update_criteria.gy_dynamic_rules_to_install.size(), 1);

  PolicyRule rule_out;
  optional<RuleToProcess> op_to_process =
      session_state->remove_gy_rule("redirect", &rule_out, &update_criteria);
  EXPECT_TRUE(op_to_process);
  EXPECT_EQ(op_to_process->version, 2);

  // basic sanity checks to see it's properly deleted
  rules_out = {};
  session_state->get_gy_dynamic_rules().get_rule_ids(rules_out_ptr);
  EXPECT_EQ(rules_out_ptr.size(), 0);

  EXPECT_FALSE(session_state->is_gy_dynamic_rule_installed("redirect"));
  EXPECT_EQ(update_criteria.gy_dynamic_rules_to_uninstall.size(), 1);

  EXPECT_EQ(0, get_monitored_rule_count("m1"));

  std::string mkey;
  // searching for non-existent rule should fail
  EXPECT_FALSE(
      session_state->get_dynamic_rules().get_monitoring_key_for_rule_id(
          "redirect", &mkey));
  // deleting an already deleted rule should fail
  EXPECT_FALSE(
      session_state->get_dynamic_rules().remove_rule("redirect", &rule_out));
}

TEST_F(SessionStateTest, test_final_credit_redirect_install) {
  insert_rule(1, "m1", "rule1", STATIC, 0, 0);
  CreditUpdateResponse charge_resp;
  charge_resp.set_success(true);
  charge_resp.set_sid("IMSI1");
  charge_resp.set_charging_key(1);

  bool is_final = true;
  auto p_credit = charge_resp.mutable_credit();
  create_charging_credit(1024, is_final, p_credit);
  auto redirect = p_credit->mutable_redirect_server();
  redirect->set_redirect_server_address("google.com");
  redirect->set_redirect_address_type(RedirectServer_RedirectAddressType_URL);
  p_credit->set_final_action(ChargingCredit_FinalAction_REDIRECT);

  session_state->receive_charging_credit(charge_resp, &update_criteria);

  // Test that the update criteria is filled out properly
  EXPECT_EQ(update_criteria.charging_credit_to_install.size(), 1);
  auto u_credit = update_criteria.charging_credit_to_install[1];
  EXPECT_TRUE(u_credit.is_final);
  auto fa = u_credit.final_action_info;
  EXPECT_EQ(fa.final_action, ChargingCredit_FinalAction_REDIRECT);
  EXPECT_EQ(fa.redirect_server.redirect_server_address(), "google.com");
  EXPECT_EQ(fa.redirect_server.redirect_address_type(),
            RedirectServer_RedirectAddressType_URL);
}

TEST_F(SessionStateTest, test_final_restrict_credit_install) {
  insert_rule(1, "m1", "rule1", STATIC, 0, 0);
  CreditUpdateResponse charge_resp;
  charge_resp.set_success(true);
  charge_resp.set_sid("IMSI1");
  charge_resp.set_charging_key(1);

  bool is_final = true;
  auto p_credit = charge_resp.mutable_credit();
  create_charging_credit(1024, is_final, p_credit);
  // auto restrict_rules = p_credit->restrict_rules();
  p_credit->add_restrict_rules("restrict-rule");
  p_credit->set_final_action(ChargingCredit_FinalAction_RESTRICT_ACCESS);

  session_state->receive_charging_credit(charge_resp, &update_criteria);

  // Test that the update criteria is filled out properly
  EXPECT_EQ(update_criteria.charging_credit_to_install.size(), 1);
  auto u_credit = update_criteria.charging_credit_to_install[1];
  EXPECT_TRUE(u_credit.is_final);
  auto fa = u_credit.final_action_info;
  EXPECT_EQ(fa.final_action, ChargingCredit_FinalAction_RESTRICT_ACCESS);
  EXPECT_EQ(fa.restrict_rules[0], "restrict-rule");
}

// Test the case where the GSU is empty. (All credit has is_valid=false). We
// treat this as an invalid credit and reject it.
TEST_F(SessionStateTest, test_empty_gsu_credit_grant) {
  insert_rule(1, "m1", "rule1", STATIC, 0, 0);
  CreditUpdateResponse charge_resp;
  charge_resp.set_success(true);
  charge_resp.set_sid(IMSI1);
  charge_resp.set_charging_key(1);
  // A ChargingCredit with no GSU but FinalAction
  charge_resp.mutable_credit()->set_type(ChargingCredit::BYTES);
  // Should return false to indicate credit installation did not go through
  EXPECT_FALSE(
      session_state->receive_charging_credit(charge_resp, &update_criteria));
  // Test that the update criteria is untouched
  EXPECT_EQ(update_criteria.charging_credit_to_install.size(), 0);
}

// We want to test a case where we receive a GSU with credit 0, but we
// receive a final_action on credit exhaust.
TEST_F(SessionStateTest, test_zero_gsu_credit_grant) {
  insert_rule(1, "m1", "rule1", STATIC, 0, 0);
  CreditUpdateResponse charge_resp;
  charge_resp.set_success(true);
  charge_resp.set_sid(IMSI1);
  charge_resp.set_charging_key(1);

  // A ChargingCredit with 0 GSU and FinalAction
  auto p_credit = charge_resp.mutable_credit();
  uint64_t zero = 0;
  p_credit->set_type(ChargingCredit::BYTES);
  p_credit->set_is_final(true);
  p_credit->set_final_action(ChargingCredit_FinalAction_TERMINATE);
  create_granted_units(&zero, &zero, &zero, p_credit->mutable_granted_units());

  EXPECT_TRUE(
      session_state->receive_charging_credit(charge_resp, &update_criteria));

  // Test that the update criteria is filled out properly
  EXPECT_EQ(update_criteria.charging_credit_to_install.size(), 1);
  auto u_credit = update_criteria.charging_credit_to_install[1];
  EXPECT_TRUE(u_credit.is_final);
  EXPECT_EQ(u_credit.final_action_info.final_action,
            ChargingCredit_FinalAction_TERMINATE);

  // At this point, the charging credit for RG=1 should have no available quota
  // and the tracking should be TRACKING_UNSET
  EXPECT_EQ(u_credit.credit.grant_tracking_type, TRACKING_UNSET);
  EXPECT_EQ(u_credit.credit.buckets[ALLOWED_TOTAL], 0);
  EXPECT_EQ(u_credit.credit.buckets[ALLOWED_TX], 0);
  EXPECT_EQ(u_credit.credit.buckets[ALLOWED_RX], 0);

  // Report some rule usage, and ensure the final action gets triggered
  session_state->add_rule_usage("rule1", 1, 100, 100, 0, 0, &update_criteria);
  auto credit_uc = update_criteria.charging_credit_map[1];
  EXPECT_EQ(credit_uc.service_state, SERVICE_NEEDS_DEACTIVATION);
}

TEST_F(SessionStateTest, test_multiple_final_action_empty_grant) {
  // add one rule with credits
  insert_rule(1, "", "rule1", STATIC, 0, 0);
  EXPECT_TRUE(std::find(update_criteria.static_rules_to_install.begin(),
                        update_criteria.static_rules_to_install.end(),
                        "rule1") !=
              update_criteria.static_rules_to_install.end());

  receive_credit_from_ocs(1, 3000, 2000, 2000, false);
  EXPECT_EQ(update_criteria.charging_credit_to_install.size(), 1);
  auto credit = update_criteria.charging_credit_to_install[CreditKey(1)].credit;
  EXPECT_EQ(credit.buckets[ALLOWED_TOTAL], 3000);
  EXPECT_EQ(credit.buckets[ALLOWED_TX], 2000);
  EXPECT_EQ(credit.buckets[ALLOWED_RX], 2000);
  // received granted units
  EXPECT_EQ(credit.received_granted_units.total().volume(), 3000);
  EXPECT_EQ(credit.received_granted_units.tx().volume(), 2000);
  EXPECT_EQ(credit.received_granted_units.rx().volume(), 2000);

  // add usage for 2 times to go over quota
  session_state->add_rule_usage("rule1", 1, 2000, 1000, 0, 0, &update_criteria);
  EXPECT_EQ(session_state->get_charging_credit(1, USED_TX), 2000);
  EXPECT_EQ(session_state->get_charging_credit(1, USED_RX), 1000);

  session_state->add_rule_usage("rule1", 1, 4000, 2000, 0, 0, &update_criteria);
  EXPECT_EQ(session_state->get_charging_credit(1, USED_TX), 4000);
  EXPECT_EQ(session_state->get_charging_credit(1, USED_RX), 2000);

  // check if we need to report the usage
  UpdateSessionRequest update;
  std::vector<std::unique_ptr<ServiceAction>> actions;
  session_state->get_updates(&update, &actions, &update_criteria);
  EXPECT_EQ(actions.size(), 0);
  EXPECT_EQ(update.updates_size(), 1);
  auto credit_uc = update_criteria.charging_credit_map[CreditKey(1)];
  EXPECT_EQ(credit_uc.bucket_deltas[USED_TX], 4000);
  EXPECT_EQ(credit_uc.bucket_deltas[USED_RX], 2000);
  EXPECT_EQ(credit_uc.service_state, SERVICE_ENABLED);
  EXPECT_FALSE(update_criteria.charging_credit_map[CreditKey(1)].is_final);
  EXPECT_TRUE(update_criteria.charging_credit_map[CreditKey(1)].reporting);

  // recive final unit without grant
  receive_credit_from_ocs(1, 0, 0, 0, true);
  EXPECT_EQ(update_criteria.charging_credit_to_install.size(), 1);
  credit_uc = update_criteria.charging_credit_map[CreditKey(1)];
  EXPECT_EQ(credit_uc.bucket_deltas[REPORTED_TX], 4000);
  EXPECT_EQ(credit_uc.bucket_deltas[REPORTED_RX], 2000);
  EXPECT_TRUE(credit_uc.is_final);
  EXPECT_EQ(credit_uc.service_state, SERVICE_NEEDS_DEACTIVATION);
  EXPECT_FALSE(credit_uc.reporting);
  // received granted units
  EXPECT_EQ(credit_uc.received_granted_units.total().volume(), 0);
  EXPECT_EQ(credit_uc.received_granted_units.tx().volume(), 0);
  EXPECT_EQ(credit_uc.received_granted_units.rx().volume(), 0);
}

TEST_F(SessionStateTest, test_apply_session_rule_set) {
  // populate rule store with 2 static and 2 dynamic rules
  insert_rule(1, "", "rule-static-1", STATIC, 0, 0);
  insert_rule(2, "m1", "rule-static-2", STATIC, 0, 0);
  insert_rule(1, "", "rule-dynamic-1", DYNAMIC, 0, 0);
  insert_rule(2, "m1", "rule-dynamic-2", DYNAMIC, 0, 0);

  EXPECT_TRUE(session_state->is_static_rule_installed("rule-static-1"));
  EXPECT_TRUE(session_state->is_static_rule_installed("rule-static-2"));
  EXPECT_TRUE(session_state->is_dynamic_rule_installed("rule-dynamic-1"));
  EXPECT_TRUE(session_state->is_dynamic_rule_installed("rule-dynamic-2"));

  // Send a set rule update with
  // 1 static rule addition: rule-static-3, 1 static rule removal: rule-static-1
  // 1 dynamic rule removal: rule-dynamic-3, 1 static rule removal:
  // rule-dynamic-1
  insert_static_rule_into_store(3, "m2", "rule-static-3");
  // Should contain all ACTIVE rules, not additional/removal
  RuleSetToApply rules_to_apply;
  rules_to_apply.static_rules.insert("rule-static-2");
  rules_to_apply.static_rules.insert("rule-static-3");

  PolicyRule dynamic_2 = create_policy_rule("rule-dynamic-2", "m1", 2);
  PolicyRule dynamic_3 = create_policy_rule("rule-dynamic-3", "m1", 3);
  rules_to_apply.dynamic_rules["rule-dynamic-2"] = dynamic_2;
  rules_to_apply.dynamic_rules["rule-dynamic-3"] = dynamic_3;

  SessionStateUpdateCriteria uc;
  RulesToProcess pending_activation, pending_deactivation, pending_bearer_setup;
  session_state->apply_session_rule_set(rules_to_apply, &pending_activation,
                                        &pending_deactivation,
                                        &pending_bearer_setup, &uc);

  // First check the active rules in session
  EXPECT_TRUE(!session_state->is_static_rule_installed("rule-static-1"));
  EXPECT_TRUE(session_state->is_static_rule_installed("rule-static-2"));
  EXPECT_TRUE(session_state->is_static_rule_installed("rule-static-3"));
  EXPECT_TRUE(!session_state->is_dynamic_rule_installed("rule-dynamic-1"));
  EXPECT_TRUE(session_state->is_dynamic_rule_installed("rule-dynamic-2"));
  EXPECT_TRUE(session_state->is_dynamic_rule_installed("rule-dynamic-3"));

  // Check the RulesToProcess is properly filled out
  EXPECT_EQ(pending_activation.size(), 2);
  const std::string activate_rule1 = pending_activation[0].rule.id();
  const std::string activate_rule2 = pending_activation[1].rule.id();
  const std::string deactivate_rule1 = pending_deactivation[0].rule.id();
  const std::string deactivate_rule2 = pending_deactivation[1].rule.id();
  EXPECT_TRUE(activate_rule1 == "rule-static-3" ||
              activate_rule1 == "rule-dynamic-3");
  EXPECT_TRUE(activate_rule2 == "rule-static-3" ||
              activate_rule2 == "rule-dynamic-3");
  EXPECT_TRUE(deactivate_rule1 == "rule-static-1" ||
              deactivate_rule1 == "rule-dynamic-1");
  EXPECT_TRUE(deactivate_rule2 == "rule-static-1" ||
              deactivate_rule2 == "rule-dynamic-1");

  // Finally assert the changes get applied to the update criteria
  EXPECT_EQ(uc.static_rules_to_install.size(), 1);
  EXPECT_EQ(uc.static_rules_to_uninstall.size(), 1);
  EXPECT_EQ(uc.dynamic_rules_to_install.size(), 1);
  EXPECT_EQ(uc.dynamic_rules_to_uninstall.size(), 1);
}

TEST_F(SessionStateTest, test_monitor_cycle) {
  // add one rule with credits
  insert_rule(1, "m1", "rule1", STATIC, 0, 0);
  EXPECT_TRUE(std::find(update_criteria.static_rules_to_install.begin(),
                        update_criteria.static_rules_to_install.end(),
                        "rule1") !=
              update_criteria.static_rules_to_install.end());

  // clear rules installed
  update_criteria = get_default_update_criteria();

  // get credit
  receive_credit_from_pcrf("m1", 3000, 2000, 2000,
                           MonitoringLevel::PCC_RULE_LEVEL);
  EXPECT_EQ(update_criteria.monitor_credit_to_install.size(), 1);
  EXPECT_EQ(update_criteria.monitor_credit_to_install["m1"]
                .credit.buckets[ALLOWED_TOTAL],
            3000);
  EXPECT_EQ(update_criteria.monitor_credit_to_install["m1"]
                .credit.buckets[ALLOWED_TX],
            2000);
  EXPECT_EQ(update_criteria.monitor_credit_to_install["m1"]
                .credit.buckets[ALLOWED_RX],
            2000);
  // received granted units
  EXPECT_EQ(update_criteria.monitor_credit_to_install["m1"]
                .credit.received_granted_units.total()
                .volume(),
            3000);
  EXPECT_EQ(update_criteria.monitor_credit_to_install["m1"]
                .credit.received_granted_units.tx()
                .volume(),
            2000);
  EXPECT_EQ(update_criteria.monitor_credit_to_install["m1"]
                .credit.received_granted_units.rx()
                .volume(),
            2000);

  // reset update_criteria (before any add_rule_usage)
  // update_criteria = get_default_update_criteria();
  // add usage for 2 times to go over quota
  session_state->add_rule_usage("rule1", 1, 2000, 1000, 0, 0, &update_criteria);
  EXPECT_EQ(session_state->get_monitor("m1", USED_TX), 2000);
  EXPECT_EQ(session_state->get_monitor("m1", USED_RX), 1000);
  EXPECT_FALSE(update_criteria.monitor_credit_map["m1"].report_last_credit);
  EXPECT_FALSE(update_criteria.monitor_credit_map["m1"].deleted);

  // receive a grant with total = 0 (meaning there is no more cuota left and
  // monitor needs to be removed once quota is exhausted
  receive_credit_from_pcrf("m1", 0, 100, 200, MonitoringLevel::PCC_RULE_LEVEL);
  // reset update_criteria (before any add_rule_usage)
  // update_criteria = get_default_update_criteria();
  session_state->add_rule_usage("rule1", 1, 4000, 2000, 0, 0, &update_criteria);
  EXPECT_EQ(session_state->get_monitor("m1", USED_TX), 4000);
  EXPECT_EQ(session_state->get_monitor("m1", USED_RX), 2000);
  EXPECT_TRUE(update_criteria.monitor_credit_map["m1"].report_last_credit);
  EXPECT_FALSE(update_criteria.monitor_credit_map["m1"].deleted);

  // Get the updates that will be sent to core
  UpdateSessionRequest update;
  std::vector<std::unique_ptr<ServiceAction>> actions;
  session_state->get_updates(&update, &actions, &update_criteria);
  EXPECT_EQ(actions.size(), 0);
  EXPECT_EQ(update.usage_monitors_size(), 1);
  EXPECT_EQ(update_criteria.monitor_credit_map["m1"].bucket_deltas[USED_TX],
            4000);
  EXPECT_EQ(update_criteria.monitor_credit_map["m1"].bucket_deltas[USED_RX],
            2000);
  EXPECT_EQ(update_criteria.monitor_credit_map["m1"].service_state,
            SERVICE_ENABLED);
  EXPECT_TRUE(update_criteria.monitor_credit_map["m1"].report_last_credit);
  EXPECT_TRUE(update_criteria.monitor_credit_map["m1"].deleted);
  EXPECT_TRUE(update_criteria.monitor_credit_map["m1"].reporting);

  // check that the monitor is actually deleted
  bool success = session_state->apply_update_criteria(update_criteria);
  EXPECT_TRUE(success);
  EXPECT_EQ(session_state->get_monitor("m1", USED_TX), 0);
  EXPECT_EQ(session_state->get_monitor("m1", USED_RX), 0);
}

TEST_F(SessionStateTest, test_get_charging_credit_summaries) {
  insert_rule(1, "m1", "rule1", STATIC, 0, 0);
  insert_rule(2, "m2", "dyn_rule1", DYNAMIC, 0, 0);

  receive_credit_from_ocs(1, 3000);
  receive_credit_from_ocs(2, 6000);
  EXPECT_EQ(update_criteria.charging_credit_to_install.size(), 2);
  EXPECT_EQ(update_criteria.charging_credit_to_install[CreditKey(1)]
                .credit.buckets[ALLOWED_TOTAL],
            3000);

  session_state->add_rule_usage("rule1", 1, 2000, 1000, 0, 0, nullptr);
  EXPECT_EQ(session_state->get_charging_credit(1, USED_TX), 2000);
  EXPECT_EQ(session_state->get_charging_credit(1, USED_RX), 1000);

  session_state->add_rule_usage("dyn_rule1", 1, 4000, 2000, 0, 0, nullptr);
  EXPECT_EQ(session_state->get_charging_credit(2, USED_TX), 4000);
  EXPECT_EQ(session_state->get_charging_credit(2, USED_RX), 2000);

  auto summaries = session_state->get_charging_credit_summaries();
  EXPECT_EQ(2, summaries.size());
  const auto& summary_rg1 = summaries[CreditKey(1)];
  const auto& summary_rg2 = summaries[CreditKey(2)];
  EXPECT_EQ(1000, summary_rg1.usage.bytes_rx);
  EXPECT_EQ(2000, summary_rg1.usage.bytes_tx);
  EXPECT_EQ(2000, summary_rg2.usage.bytes_rx);
  EXPECT_EQ(4000, summary_rg2.usage.bytes_tx);
}

TEST_F(SessionStateTest, test_process_static_rule_installs) {
  initialize_session_with_qos();

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
  RulesToProcess pending_activation, pending_deactivation, pending_bearer_setup;
  RulesToSchedule pending_scheduling;
  session_state->process_static_rule_installs(
      rule_installs, &pending_activation, &pending_deactivation,
      &pending_bearer_setup, &pending_scheduling, nullptr);
  EXPECT_EQ(1, pending_activation.size());
  EXPECT_EQ("static-2", pending_activation[0].rule.id());
  EXPECT_EQ(1, pending_bearer_setup.size());
  EXPECT_EQ("static-qos-4", pending_bearer_setup[0].rule.id());
}

TEST_F(SessionStateTest, test_process_dynamic_rule_installs) {
  initialize_session_with_qos();

  PolicyRule dynamic_1 = create_policy_rule("dynamic-1", "", 0);
  PolicyRule dynamic_2 = create_policy_rule("dynamic-2", "", 0);
  PolicyRule dynamic_qos_3 =
      create_policy_rule_with_qos("dynamic-qos-3", "", 0, 1);
  PolicyRule dynamic_qos_4 =
      create_policy_rule_with_qos("dynamic-qos-4", "", 0, 2);

  // Install dynamic rules for dynamic-1 and dynamic-qos-3
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
  RulesToProcess pending_activation, pending_deactivation, pending_bearer_setup;
  RulesToSchedule pending_scheduling;
  session_state->process_dynamic_rule_installs(
      rule_installs, &pending_activation, &pending_deactivation,
      &pending_bearer_setup, &pending_scheduling, nullptr);
  EXPECT_EQ(2, pending_activation.size());
  EXPECT_EQ(2, pending_bearer_setup.size());
}

int main(int argc, char** argv) {
  ::testing::InitGoogleTest(&argc, argv);
  FLAGS_logtostderr = 1;
  FLAGS_v = 10;
  return RUN_ALL_TESTS();
}

}  // namespace magma

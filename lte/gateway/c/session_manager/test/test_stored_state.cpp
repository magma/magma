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

#include <memory>

#include "magma_logging.h"
#include "ProtobufCreators.h"
#include "StoredState.h"

using ::testing::Test;

namespace magma {

class StoredStateTest : public ::testing::Test {
 protected:
  SessionConfig get_stored_session_config() {
    SessionConfig stored;
    Teids teids;
    teids.set_agw_teid(1);
    teids.set_enb_teid(2);
    stored.common_context = build_common_context(
        "IMSI1", "ue_ipv4", "ue_ipv6", teids, "apn", "msisdn", TGPP_WLAN);
    const auto& lte_context =
        build_lte_context("192.168.0.2", "imei", "plmn_id", "imsi_plmn_id",
                          "user_location", 321, nullptr);
    stored.rat_specific_context.mutable_lte_context()->CopyFrom(lte_context);
    return stored;
  }

  FinalActionInfo get_stored_redirect_final_action_info() {
    FinalActionInfo stored;
    stored.final_action =
        ChargingCredit_FinalAction::ChargingCredit_FinalAction_REDIRECT;
    stored.redirect_server.set_redirect_address_type(
        RedirectServer_RedirectAddressType::
            RedirectServer_RedirectAddressType_IPV6);
    stored.redirect_server.set_redirect_server_address(
        "redirect_server_address");
    return stored;
  }

  FinalActionInfo get_stored_restrict_final_action_info() {
    FinalActionInfo stored;
    stored.final_action =
        ChargingCredit_FinalAction::ChargingCredit_FinalAction_RESTRICT_ACCESS;
    stored.restrict_rules.push_back("restrict_rule");
    return stored;
  }

  StoredSessionCredit get_stored_session_credit() {
    StoredSessionCredit stored;

    stored.reporting = true;
    stored.credit_limit_type = INFINITE_METERED;

    stored.buckets[USED_TX] = 12345;
    stored.buckets[ALLOWED_TOTAL] = 54321;

    stored.grant_tracking_type = TX_ONLY;
    return stored;
  };

  StoredChargingGrant get_stored_charging_grant() {
    StoredChargingGrant stored;

    stored.is_final = true;

    stored.final_action_info.final_action =
        ChargingCredit_FinalAction::ChargingCredit_FinalAction_REDIRECT;
    stored.final_action_info.redirect_server.set_redirect_address_type(
        RedirectServer_RedirectAddressType::
            RedirectServer_RedirectAddressType_IPV6);
    stored.final_action_info.redirect_server.set_redirect_server_address(
        "redirect_server_address");

    stored.service_state = SERVICE_NEEDS_ACTIVATION;
    stored.reauth_state = REAUTH_REQUIRED;

    stored.expiry_time = 32;
    stored.credit = get_stored_session_credit();
    return stored;
  };

  StoredMonitor get_stored_monitor() {
    StoredMonitor stored;
    stored.credit = get_stored_session_credit();
    stored.level = MonitoringLevel::PCC_RULE_LEVEL;
    return stored;
  }

  StoredChargingCreditMap get_stored_charging_credit_map() {
    StoredChargingCreditMap stored(4, &ccHash, &ccEqual);
    stored[CreditKey(1, 2)] = get_stored_charging_grant();
    return stored;
  }

  StoredMonitorMap get_stored_monitor_map() {
    StoredMonitorMap stored;
    stored["mk1"] = get_stored_monitor();
    return stored;
  }

  PolicyStatsMap get_policy_stats_map() {
    PolicyStatsMap stored;
    stored["rule1"] = StatsPerPolicy();
    stored["rule1"].last_reported_version = 1;
    stored["rule1"].stats_map[1] = RuleStats{1, 2, 3, 4};
    stored["rule2"] = StatsPerPolicy();
    stored["rule2"].last_reported_version = 2;
    stored["rule2"].stats_map[1] = RuleStats{5, 2, 3, 4};
    stored["rule2"].stats_map[2] = RuleStats{50, 20, 30, 40};
    return stored;
  }

  BearerIDByPolicyID get_bearer_id_by_policy() {
    BearerIDByPolicyID stored;
    stored[PolicyID(DYNAMIC, "rule1")].bearer_id = 32;
    stored[PolicyID(DYNAMIC, "rule1")].teids.set_agw_teid(1);
    stored[PolicyID(DYNAMIC, "rule1")].teids.set_enb_teid(2);
    stored[PolicyID(STATIC, "rule1")].bearer_id = 64;
    stored[PolicyID(STATIC, "rule1")].teids.set_agw_teid(3);
    stored[PolicyID(STATIC, "rule1")].teids.set_enb_teid(4);
    return stored;
  }

  StoredSessionState get_stored_session() {
    StoredSessionState stored;

    stored.config = get_stored_session_config();
    stored.credit_map = get_stored_charging_credit_map();
    stored.monitor_map = get_stored_monitor_map();
    stored.session_level_key = "session_level_key";
    stored.imsi = "IMSI1";
    stored.session_id = "session_id";
    stored.subscriber_quota_state = SubscriberQuotaUpdate_Type_VALID_QUOTA;
    stored.fsm_state = SESSION_RELEASED;

    magma::lte::TgppContext tgpp_context;
    tgpp_context.set_gx_dest_host("gx");
    tgpp_context.set_gy_dest_host("gy");
    stored.tgpp_context = tgpp_context;
    stored.pdp_start_time = 112233;
    stored.pdp_end_time = 332211;

    stored.pending_event_triggers[REVALIDATION_TIMEOUT] = READY;
    stored.revalidation_time.set_seconds(32);

    stored.bearer_id_by_policy = get_bearer_id_by_policy();

    stored.request_number = 1;

    stored.policy_version_and_stats = get_policy_stats_map();

    return stored;
  }
};

TEST_F(StoredStateTest, test_stored_session_config) {
  auto stored = get_stored_session_config();

  std::string serialized = serialize_stored_session_config(stored);
  SessionConfig deserialized = deserialize_stored_session_config(serialized);

  // compare the serialized objects
  auto original_common = stored.common_context.SerializeAsString();
  auto original_rat_specific = stored.rat_specific_context.SerializeAsString();
  auto recovered_common = deserialized.common_context.SerializeAsString();
  auto recovered_rat_specific =
      deserialized.rat_specific_context.SerializeAsString();
  EXPECT_EQ(original_common, recovered_common);
  EXPECT_EQ(original_rat_specific, recovered_rat_specific);
}

TEST_F(StoredStateTest, test_stored_redirect_final_action_info) {
  auto stored = get_stored_redirect_final_action_info();

  auto serialized = serialize_stored_final_action_info(stored);
  auto deserialized = deserialize_stored_final_action_info(serialized);

  EXPECT_EQ(deserialized.final_action,
            ChargingCredit_FinalAction::ChargingCredit_FinalAction_REDIRECT);
  EXPECT_EQ(deserialized.redirect_server.redirect_address_type(),
            RedirectServer_RedirectAddressType::
                RedirectServer_RedirectAddressType_IPV6);
  EXPECT_EQ(deserialized.redirect_server.redirect_server_address(),
            "redirect_server_address");
}

TEST_F(StoredStateTest, test_stored_restrict_final_action_info) {
  auto stored = get_stored_restrict_final_action_info();

  auto serialized = serialize_stored_final_action_info(stored);
  auto deserialized = deserialize_stored_final_action_info(serialized);

  EXPECT_EQ(
      deserialized.final_action,
      ChargingCredit_FinalAction::ChargingCredit_FinalAction_RESTRICT_ACCESS);

  std::vector<std::string> restrict_rules = {"restrict_rule"};
  EXPECT_EQ(deserialized.restrict_rules, restrict_rules);
}

TEST_F(StoredStateTest, test_stored_bearer_id_by_policy) {
  auto stored = get_bearer_id_by_policy();
  auto serialized = serialize_bearer_id_by_policy(stored);
  auto deserialized = deserialize_bearer_id_by_policy(serialized);
  EXPECT_EQ(stored[PolicyID(DYNAMIC, "rule1")].bearer_id, 32);
  EXPECT_EQ(stored.size(), deserialized.size());
  EXPECT_EQ(stored[PolicyID(DYNAMIC, "rule1")],
            deserialized[PolicyID(DYNAMIC, "rule1")]);
  EXPECT_EQ(stored[PolicyID(STATIC, "rule1")],
            deserialized[PolicyID(STATIC, "rule1")]);
}

TEST_F(StoredStateTest, test_stored_session_credit) {
  auto stored = get_stored_session_credit();

  auto serialized = serialize_stored_session_credit(stored);
  auto deserialized = deserialize_stored_session_credit(serialized);

  EXPECT_EQ(deserialized.reporting, true);
  EXPECT_EQ(deserialized.credit_limit_type, INFINITE_METERED);

  EXPECT_EQ(deserialized.buckets[USED_TX], 12345);
  EXPECT_EQ(deserialized.buckets[ALLOWED_TOTAL], 54321);

  EXPECT_EQ(deserialized.grant_tracking_type, TX_ONLY);
}

TEST_F(StoredStateTest, test_stored_monitor) {
  auto stored = get_stored_monitor();

  auto serialized = serialize_stored_monitor(stored);
  auto deserialized = deserialize_stored_monitor(serialized);

  EXPECT_EQ(deserialized.credit.reporting, true);
  EXPECT_EQ(deserialized.credit.credit_limit_type, INFINITE_METERED);
  EXPECT_EQ(deserialized.credit.buckets[USED_TX], 12345);
  EXPECT_EQ(deserialized.credit.buckets[ALLOWED_TOTAL], 54321);
  EXPECT_EQ(deserialized.level, MonitoringLevel::PCC_RULE_LEVEL);
}

TEST_F(StoredStateTest, test_stored_charging_credit_map) {
  auto stored = get_stored_charging_credit_map();

  auto serialized = serialize_stored_charging_credit_map(stored);
  auto deserialized = deserialize_stored_charging_credit_map(serialized);

  auto stored_charging_credit = deserialized[CreditKey(1, 2)];
  // test charging grant fields
  EXPECT_EQ(stored_charging_credit.is_final, true);
  EXPECT_EQ(stored_charging_credit.final_action_info.final_action,
            ChargingCredit_FinalAction::ChargingCredit_FinalAction_REDIRECT);
  EXPECT_EQ(stored_charging_credit.final_action_info.redirect_server
                .redirect_address_type(),
            RedirectServer_RedirectAddressType::
                RedirectServer_RedirectAddressType_IPV6);
  EXPECT_EQ(stored_charging_credit.final_action_info.redirect_server
                .redirect_server_address(),
            "redirect_server_address");
  EXPECT_EQ(stored_charging_credit.reauth_state, REAUTH_REQUIRED);
  EXPECT_EQ(stored_charging_credit.service_state, SERVICE_NEEDS_ACTIVATION);
  EXPECT_EQ(stored_charging_credit.expiry_time, 32);

  // test session credit fields
  auto credit = stored_charging_credit.credit;
  EXPECT_EQ(credit.reporting, true);
  EXPECT_EQ(credit.credit_limit_type, INFINITE_METERED);
  EXPECT_EQ(credit.buckets[USED_TX], 12345);
  EXPECT_EQ(credit.buckets[ALLOWED_TOTAL], 54321);
}

TEST_F(StoredStateTest, test_stored_monitor_map) {
  auto stored = get_stored_monitor_map();

  auto serialized = serialize_stored_usage_monitor_map(stored);
  auto deserialized = deserialize_stored_usage_monitor_map(serialized);

  auto stored_monitor = deserialized["mk1"];
  EXPECT_EQ(stored_monitor.credit.reporting, true);
  EXPECT_EQ(stored_monitor.credit.credit_limit_type, INFINITE_METERED);
  EXPECT_EQ(stored_monitor.credit.buckets[USED_TX], 12345);
  EXPECT_EQ(stored_monitor.credit.buckets[ALLOWED_TOTAL], 54321);
  EXPECT_EQ(stored_monitor.level, MonitoringLevel::PCC_RULE_LEVEL);
}

TEST_F(StoredStateTest, test_stored_session) {
  auto stored = get_stored_session();

  auto serialized = serialize_stored_session(stored);
  auto deserialized = deserialize_stored_session(serialized);

  auto stored_charging_credit = deserialized.credit_map[CreditKey(1, 2)];
  // test charging grant fields
  EXPECT_EQ(stored_charging_credit.is_final, true);
  EXPECT_EQ(stored_charging_credit.final_action_info.final_action,
            ChargingCredit_FinalAction::ChargingCredit_FinalAction_REDIRECT);
  EXPECT_EQ(stored_charging_credit.final_action_info.redirect_server
                .redirect_address_type(),
            RedirectServer_RedirectAddressType::
                RedirectServer_RedirectAddressType_IPV6);
  EXPECT_EQ(stored_charging_credit.final_action_info.redirect_server
                .redirect_server_address(),
            "redirect_server_address");
  EXPECT_EQ(stored_charging_credit.reauth_state, REAUTH_REQUIRED);
  EXPECT_EQ(stored_charging_credit.service_state, SERVICE_NEEDS_ACTIVATION);
  EXPECT_EQ(stored_charging_credit.expiry_time, 32);

  // test session credit fields
  auto credit = stored_charging_credit.credit;
  EXPECT_EQ(credit.reporting, true);
  EXPECT_EQ(credit.buckets[USED_TX], 12345);
  EXPECT_EQ(credit.buckets[ALLOWED_TOTAL], 54321);
  EXPECT_EQ(credit.credit_limit_type, INFINITE_METERED);

  EXPECT_EQ(deserialized.session_level_key, "session_level_key");

  auto stored_monitor = deserialized.monitor_map["mk1"];
  EXPECT_EQ(stored_monitor.credit.reporting, true);
  EXPECT_EQ(stored_monitor.credit.credit_limit_type, INFINITE_METERED);
  EXPECT_EQ(stored_monitor.credit.buckets[USED_TX], 12345);
  EXPECT_EQ(stored_monitor.credit.buckets[ALLOWED_TOTAL], 54321);
  EXPECT_EQ(stored_monitor.level, MonitoringLevel::PCC_RULE_LEVEL);

  EXPECT_EQ(deserialized.imsi, "IMSI1");
  EXPECT_EQ(deserialized.session_id, "session_id");
  EXPECT_EQ(deserialized.subscriber_quota_state,
            SubscriberQuotaUpdate_Type_VALID_QUOTA);
  EXPECT_EQ(deserialized.fsm_state, SESSION_RELEASED);

  EXPECT_EQ(deserialized.tgpp_context.gx_dest_host(), "gx");
  EXPECT_EQ(deserialized.tgpp_context.gy_dest_host(), "gy");

  EXPECT_EQ(deserialized.pending_event_triggers.size(), 1);
  EXPECT_EQ(deserialized.pending_event_triggers[REVALIDATION_TIMEOUT], READY);
  EXPECT_EQ(deserialized.revalidation_time.seconds(), 32);

  EXPECT_EQ(deserialized.bearer_id_by_policy.size(), 2);

  EXPECT_EQ(deserialized.request_number, 1);
  EXPECT_EQ(deserialized.pdp_start_time, 112233);
  EXPECT_EQ(deserialized.pdp_end_time, 332211);

  EXPECT_EQ(
      deserialized.policy_version_and_stats["rule1"].last_reported_version, 1);
  EXPECT_EQ(deserialized.policy_version_and_stats["rule1"].stats_map[1].tx, 1);
  EXPECT_EQ(deserialized.policy_version_and_stats["rule1"].stats_map[1].rx, 2);
  EXPECT_EQ(
      deserialized.policy_version_and_stats["rule2"].stats_map[1].dropped_tx,
      3);
  EXPECT_EQ(
      deserialized.policy_version_and_stats["rule2"].stats_map[2].dropped_rx,
      40);
}

TEST_F(StoredStateTest, test_policy_stats_map) {
  PolicyStatsMap original;
  StatsPerPolicy og_stats1, og_stats2;
  const std::string rule1 = "rule1";
  const std::string rule2 = "rule2";

  og_stats1.current_version = 2;
  og_stats1.last_reported_version = 1;
  original[rule1] = og_stats1;

  og_stats2.current_version = 4;
  og_stats2.last_reported_version = 3;
  original[rule2] = og_stats2;

  std::string serialized = serialize_policy_stats_map(original);
  PolicyStatsMap deserialized = deserialize_policy_stats_map(serialized);

  EXPECT_EQ(2, deserialized.size());

  StatsPerPolicy deserialized_stats1 = deserialized[rule1];
  StatsPerPolicy deserialized_stats2 = deserialized[rule2];

  EXPECT_EQ(og_stats1.current_version, deserialized_stats1.current_version);
  EXPECT_EQ(og_stats1.last_reported_version,
            deserialized_stats1.last_reported_version);

  EXPECT_EQ(og_stats2.current_version, deserialized_stats2.current_version);
  EXPECT_EQ(og_stats2.last_reported_version,
            deserialized_stats2.last_reported_version);

  // Check that the value is empty by default
  EXPECT_FALSE(get_default_update_criteria().policy_version_and_stats);
}

int main(int argc, char** argv) {
  ::testing::InitGoogleTest(&argc, argv);
  return RUN_ALL_TESTS();
}

}  // namespace magma

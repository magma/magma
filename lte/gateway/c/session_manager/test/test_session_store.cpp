/**
 * Copyright (c) 2016-present, Facebook, Inc.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include <memory>

#include <glog/logging.h>
#include <gtest/gtest.h>

#include "RuleStore.h"
#include "SessionID.h"
#include "SessionState.h"
#include "SessionStore.h"
#include "StoredState.h"
#include "MagmaService.h"
#include "magma_logging.h"

using magma::orc8r::MetricsContainer;
using ::testing::Test;

namespace magma {

class SessionStoreTest : public ::testing::Test {
 protected:
  SessionIDGenerator id_gen_;

 protected:
  virtual void SetUp()
  {
    imsi = "IMSI1";
    imsi2 = "IMSI2";
    sid = id_gen_.gen_session_id(imsi);
    sid2 = id_gen_.gen_session_id(imsi2);
    sid3 = id_gen_.gen_session_id(imsi2);
    monitoring_key = "mk1";
    monitoring_key2 = "mk2";
    rule_id_1 = "test_rule_1";
    rule_id_2 = "test_rule_2";
    rule_id_3 = "test_rule_3";
    dynamic_rule_id_1 = "dynamic_rule_1";
    dynamic_rule_id_2 = "dynamic_rule_2";
  }

  PolicyRule get_dynamic_rule()
  {
    auto policy = new PolicyRule();
    policy->set_id(dynamic_rule_id_1);
    policy->set_priority(10);
    policy->set_tracking_type(PolicyRule::ONLY_OCS);
    return *policy;
  }

  std::unique_ptr<SessionState> get_session(
    std::string session_id,
    std::shared_ptr<StaticRuleStore> rule_store)
  {
    std::string hardware_addr_bytes = {0x0f, 0x10, 0x2e, 0x12, 0x3a, 0x55};
    std::string msisdn = "5100001234";
    std::string radius_session_id =
      "AA-AA-AA-AA-AA-AA:TESTAP__"
      "0F-10-2E-12-3A-55";
    std::string core_session_id = "asdf";
    SessionConfig cfg = {.ue_ipv4 = "",
      .spgw_ipv4 = "",
      .msisdn = msisdn,
      .apn = "",
      .imei = "",
      .plmn_id = "",
      .imsi_plmn_id = "",
      .user_location = "",
      .rat_type = RATType::TGPP_WLAN,
      .mac_addr = "0f:10:2e:12:3a:55",
      .hardware_addr = hardware_addr_bytes,
      .radius_session_id = radius_session_id};
    auto tgpp_context = TgppContext{};
    auto session = std::make_unique<SessionState>(
      imsi, session_id, core_session_id, cfg, *rule_store, tgpp_context);
    return std::move(session);
  }

  UsageMonitoringUpdateResponse* get_monitoring_update()
  {
    auto units = new GrantedUnits();
    auto total = new CreditUnit();
    total->set_is_valid(true);
    total->set_volume(1000);
    auto tx = new CreditUnit();
    tx->set_is_valid(true);
    tx->set_volume(1000);
    auto rx = new CreditUnit();
    rx->set_is_valid(true);
    rx->set_volume(1000);
    units->set_allocated_total(total);
    units->set_allocated_tx(tx);
    units->set_allocated_rx(rx);

    auto monitoring_credit = new UsageMonitoringCredit();
    monitoring_credit->set_action(UsageMonitoringCredit_Action_CONTINUE);
    monitoring_credit->set_monitoring_key(monitoring_key);
    monitoring_credit->set_level(SESSION_LEVEL);
    monitoring_credit->set_allocated_granted_units(units);

    auto credit_update = new UsageMonitoringUpdateResponse();
    credit_update->set_allocated_credit(monitoring_credit);
    credit_update->set_session_id("sid1");
    credit_update->set_success(true);
    // Don't set event triggers
    // Don't set result code since the response is already successful
    // Don't set any rule installation/uninstallation
    // Don't set the TgppContext, assume relay disabled
    return credit_update;
  }

  SessionStateUpdateCriteria get_update_criteria()
  {
    auto update_criteria = SessionStateUpdateCriteria{};

    // Rule installation
    update_criteria.static_rules_to_install = std::set<std::string>{};
    update_criteria.static_rules_to_install.insert(rule_id_1);
    update_criteria.dynamic_rules_to_install = std::vector<PolicyRule>{};
    update_criteria.dynamic_rules_to_install.push_back(get_dynamic_rule());
    RuleLifetime lifetime{
      .activation_time = std::time_t(0),
      .deactivation_time = std::time_t(0),
    };
    update_criteria.new_rule_lifetimes[rule_id_1] = lifetime;
    update_criteria.new_rule_lifetimes[dynamic_rule_id_1] = lifetime;

    // Monitoring credit installation
    update_criteria.monitor_credit_to_install =
      std::unordered_map<std::string, StoredMonitor>{};
    auto monitor2 = StoredMonitor{};
    auto credit2 = StoredSessionCredit{};
    credit2.reporting = false;
    credit2.is_final = false;
    credit2.unlimited_quota = false;
    credit2.service_state = SERVICE_ENABLED;
    credit2.expiry_time = 0;
    credit2.buckets = std::unordered_map<Bucket, uint64_t>{};
    credit2.buckets[USED_TX] = 100;
    credit2.buckets[USED_RX] = 200;
    credit2.buckets[ALLOWED_TOTAL] = 2;
    credit2.buckets[ALLOWED_TX] = 3;
    credit2.buckets[ALLOWED_RX] = 4;
    credit2.buckets[REPORTING_TX] = 5;
    credit2.buckets[REPORTING_RX] = 6;
    credit2.buckets[REPORTED_TX] = 7;
    credit2.buckets[REPORTED_RX] = 8;
    credit2.usage_reporting_limit = 12345;
    monitor2.level = SESSION_LEVEL;
    monitor2.credit = credit2;
    update_criteria.monitor_credit_to_install[monitoring_key2] = monitor2;

    // Monitoring credit updates
    SessionCreditUpdateCriteria monitoring_update{};
    monitoring_update.is_final = false;
    monitoring_update.reauth_state = REAUTH_NOT_NEEDED;
    monitoring_update.service_state = SERVICE_ENABLED;
    monitoring_update.expiry_time = 0;
    auto bucket_deltas = std::unordered_map<Bucket, uint64_t>{};
    bucket_deltas[USED_TX] = 111;
    bucket_deltas[USED_RX] = 333;
    bucket_deltas[ALLOWED_TOTAL] = 2;
    bucket_deltas[ALLOWED_TX] = 3;
    bucket_deltas[ALLOWED_RX] = 4;
    bucket_deltas[REPORTING_TX] = 5;
    bucket_deltas[REPORTING_RX] = 6;
    bucket_deltas[REPORTED_TX] = 7;
    bucket_deltas[REPORTED_RX] = 8;
    monitoring_update.bucket_deltas = bucket_deltas;

    update_criteria.monitor_credit_map =
      std::unordered_map<std::string, SessionCreditUpdateCriteria>{};
    update_criteria.monitor_credit_map[monitoring_key] = monitoring_update;

    return update_criteria;
  }

  bool is_equal(
      io::prometheus::client::LabelPair label_pair, const char*& name,
      const char*& value) {
    return label_pair.name().compare(name) == 0 &&
           label_pair.value().compare(value) == 0;
  }

 protected:
  std::string imsi;
  std::string imsi2;
  std::string sid;
  std::string sid2;
  std::string sid3;
  std::string monitoring_key;
  std::string monitoring_key2;
  std::string rule_id_1;
  std::string rule_id_2;
  std::string rule_id_3;
  std::string dynamic_rule_id_1;
  std::string dynamic_rule_id_2;
};

TEST_F(SessionStoreTest, test_metering_reporting)
{
  // 1) Create SessionStore
  auto rule_store = std::make_shared<StaticRuleStore>();
  auto session_store = new SessionStore(rule_store);

  // 2) Create a single session and write it into the store
  auto session1 = get_session(sid, rule_store);
  auto session_vec = std::vector<std::unique_ptr<SessionState>>{};
  session_vec.push_back(std::move(session1));
  session_store->create_sessions(imsi, std::move(session_vec));

  // 3) Try to update the session in SessionStore with a rule installation
  auto session_map = session_store->read_sessions(SessionRead{imsi});
  auto session_update = SessionStore::get_default_session_update(session_map);

  auto uc = get_default_update_criteria();
  uc.static_rules_to_install.insert("RULE_asdf");
  RuleLifetime lifetime{
      .activation_time = std::time_t(0),
      .deactivation_time = std::time_t(0),
  };
  uc.new_rule_lifetimes["RULE_asdf"] = lifetime;

  // Record some credit usage
  auto DIRECTION_LABEL   = "direction";
  auto DIRECTION_UP   = "up";
  auto DIRECTION_DOWN = "down";
  auto UPLOADED_BYTES   = 5;
  auto DOWNLOADED_BYTES = 7;
  SessionCreditUpdateCriteria credit_uc{};
  credit_uc.bucket_deltas[USED_TX]      = UPLOADED_BYTES;
  credit_uc.bucket_deltas[USED_RX]      = DOWNLOADED_BYTES;
  uc.monitor_credit_map[monitoring_key] = credit_uc;

  session_update[imsi][sid] = uc;

  auto update_success = session_store->update_sessions(session_update);
  EXPECT_TRUE(update_success);

  // verify if UE traffic metrics are recorded properly
  auto resp = new MetricsContainer();
  auto magma_service =
      std::make_shared<service303::MagmaService>("test_service", "1.0");
  magma_service->GetMetrics(nullptr, nullptr, resp);
  auto reported_metrics = 0;
  for (auto const& fam : resp->family()) {
    if (fam.name().compare("ue_traffic") == 0) {
      for (auto const& m : fam.metric()) {
        for (auto const& l : m.label()) {
          if (is_equal(l, DIRECTION_LABEL, DIRECTION_UP)) {
            EXPECT_EQ(m.counter().value(), UPLOADED_BYTES);
            reported_metrics += 1;
          } else if (is_equal(l, DIRECTION_LABEL, DIRECTION_DOWN)) {
            EXPECT_EQ(m.counter().value(), DOWNLOADED_BYTES);
            reported_metrics += 1;
          }
        }
      }
      break;
    }
  }
  EXPECT_EQ(reported_metrics, 2);
}

/**
 * End to end test of the SessionStore.
 * 1) Create SessionStore
 * 2) Create bare-bones session for IMSI1
 * 3) Commit session for IMSI1 into SessionStore
 * 4) Read session for IMSI1 from SessionStore
 * 5) Verify that state was written for IMSI1 and has been retrieved.
 * 6) Make updates to session
 * 7) Commit updates to SessionStore
 * 8) Read in session for IMSI1 again, and check that the update was successful
 * 9) Check request numbers again
 * 10) Update request numbers again and check to see that they're updated
 *     correctly still for multiple monitoring keys
 * 11) Delete the session for IMSI1
 * 12) Verify IMSI1 no longer has any sessions
 */
TEST_F(SessionStoreTest, test_read_and_write)
{
  // 1) Create SessionStore
  auto rule_store = std::make_shared<StaticRuleStore>();
  auto session_store = new SessionStore(rule_store);

  // 2) Create bare-bones session for IMSI1
  auto session = get_session(sid, rule_store);
  auto uc = get_default_update_criteria();
  RuleLifetime lifetime{
    .activation_time = std::time_t(0),
    .deactivation_time = std::time_t(0),
  };
  session->activate_static_rule(rule_id_3, lifetime, uc);
  EXPECT_EQ(session->get_session_id(), sid);
  EXPECT_EQ(session->get_request_number(), 1);
  EXPECT_EQ(session->is_static_rule_installed(rule_id_3),true);

  auto credit_update = get_monitoring_update();
  UsageMonitoringUpdateResponse& credit_update_ref = *credit_update;
  session->get_monitor_pool().receive_credit(credit_update_ref, uc);

  // Add some used credit
  session->get_monitor_pool().add_used_credit(monitoring_key, uint64_t(111), uint64_t(333), uc);
  EXPECT_EQ(session->get_monitor_pool().get_credit(monitoring_key, USED_TX), 111);
  EXPECT_EQ(session->get_monitor_pool().get_credit(monitoring_key, USED_RX), 333);

  // 3) Commit session for IMSI1 into SessionStore
  auto sessions = std::vector<std::unique_ptr<SessionState>>{};
  EXPECT_EQ(sessions.size(), 0);
  sessions.push_back(std::move(session));
  EXPECT_EQ(sessions.size(), 1);
  session_store->create_sessions(imsi, std::move(sessions));

  // 4) Read session for IMSI1 from SessionStore
  SessionRead read_req = {};
  read_req.insert(imsi);
  auto session_map = session_store->read_sessions_for_reporting(read_req);

  // 5) Verify that state was written for IMSI1 and has been retrieved.
  EXPECT_EQ(session_map.size(), 1);
  EXPECT_EQ(session_map[imsi].size(), 1);
  EXPECT_EQ(session_map[imsi].front()->get_request_number(), 1);

  // 6) Make updates to session via SessionUpdateCriteria
  auto update_req = SessionUpdate{};
  update_req[imsi] = std::unordered_map<std::string,
                                        SessionStateUpdateCriteria>{};
  auto update_criteria = get_update_criteria();
  update_req[imsi][sid] = update_criteria;

  // 7) Commit updates to SessionStore
  auto success = session_store->update_sessions(update_req);
  EXPECT_TRUE(success);

  // 8) Read in session for IMSI1 again to check that the update was successful
  session_map = session_store->read_sessions(read_req);
  EXPECT_EQ(session_map.size(), 1);
  EXPECT_EQ(session_map[imsi].size(), 1);
  EXPECT_EQ(session_map[imsi].front()->get_session_id(), sid);
  // Check installed rules
  EXPECT_EQ(session_map[imsi].front()->is_static_rule_installed(rule_id_1),
    true);
  EXPECT_EQ(session_map[imsi].front()->is_static_rule_installed(rule_id_3),
            true);
  EXPECT_EQ(session_map[imsi].front()->is_static_rule_installed(rule_id_2),
            false);
  EXPECT_EQ(session_map[imsi].front()->is_dynamic_rule_installed(dynamic_rule_id_1),
            true);
  EXPECT_EQ(session_map[imsi].front()->is_dynamic_rule_installed(dynamic_rule_id_2),
            false);

  // Check for installation of new monitoring credit
  session_map[imsi].front()->get_monitor_pool().add_monitor(monitoring_key2,
    UsageMonitoringCreditPool::unmarshal_monitor(update_criteria.monitor_credit_to_install[monitoring_key2]), uc);
  EXPECT_EQ(session_map[imsi].front()->get_monitor_pool().get_credit(monitoring_key2, USED_TX), 100);
  EXPECT_EQ(session_map[imsi].front()->get_monitor_pool().get_credit(monitoring_key2, USED_RX), 200);

  // Check monitoring credit usage
  EXPECT_EQ(session_map[imsi].front()->get_monitor_pool().get_credit(monitoring_key, USED_TX), 222);
  EXPECT_EQ(session_map[imsi].front()->get_monitor_pool().get_credit(monitoring_key, USED_RX), 666);
  EXPECT_EQ(session_map[imsi].front()->get_monitor_pool().get_credit(monitoring_key, ALLOWED_TOTAL), 1002);
  EXPECT_EQ(session_map[imsi].front()->get_monitor_pool().get_credit(monitoring_key, ALLOWED_TX), 1003);
  EXPECT_EQ(session_map[imsi].front()->get_monitor_pool().get_credit(monitoring_key, ALLOWED_RX), 1004);
  EXPECT_EQ(session_map[imsi].front()->get_monitor_pool().get_credit(monitoring_key, REPORTING_TX), 5);
  EXPECT_EQ(session_map[imsi].front()->get_monitor_pool().get_credit(monitoring_key, REPORTING_RX), 6);
  EXPECT_EQ(session_map[imsi].front()->get_monitor_pool().get_credit(monitoring_key, REPORTED_TX), 7);
  EXPECT_EQ(session_map[imsi].front()->get_monitor_pool().get_credit(monitoring_key, REPORTED_RX), 8);

  // 9) Check request numbers again
  // This request number should increment in storage every time a read is done.
  // The incremented value is set by the read request to the storage interface.
  EXPECT_EQ(session_map[imsi].front()->get_request_number(), 2);

  // 10) Read sessions for reporting to update request numbers for the session
  // The request number should be incremented by 2 for the session, 1 for
  // each monitoring key and charging key associated to it.
  session_map = session_store->read_sessions_for_reporting(read_req);
  EXPECT_EQ(session_map.size(), 1);
  EXPECT_EQ(session_map[imsi].size(), 1);

  session_map = session_store->read_sessions(read_req);
  EXPECT_EQ(session_map.size(), 1);
  EXPECT_EQ(session_map[imsi].front()->get_request_number(), 4);

  // 11) Delete sessions for IMSI1
  update_req = SessionUpdate{};
  update_criteria = SessionStateUpdateCriteria{};
  update_criteria.is_session_ended = true;
  update_req[imsi][sid] = update_criteria;
  session_store->update_sessions(update_req);

  // 12) Verify that IMSI1 no longer has a session
  session_map = session_store->read_sessions_for_reporting(read_req);
  EXPECT_EQ(session_map.size(), 1);
  EXPECT_EQ(session_map[imsi].size(), 0);
}

TEST_F(SessionStoreTest, test_get_default_session_update)
{
  // 1) Create a SessionMap with a few sessions
  auto rule_store = std::make_shared<StaticRuleStore>();
  SessionMap session_map = {};
  auto session1 = get_session(sid, rule_store);
  auto session2 = get_session(sid2, rule_store);
  auto session3 = get_session(sid3, rule_store);

  session_map[imsi] = std::vector<std::unique_ptr<SessionState>>{};
  session_map[imsi2] = std::vector<std::unique_ptr<SessionState>>{};

  session_map[imsi].push_back(std::move(session1));
  session_map[imsi2].push_back(std::move(session2));
  session_map[imsi2].push_back(std::move(session3));

  // 2) Build SessionUpdate
  auto update = SessionStore::get_default_session_update(session_map);
  EXPECT_EQ(update.size(), 2);
  EXPECT_EQ(update[imsi].size(), 1);
  EXPECT_EQ(update[imsi2].size(), 2);
}

TEST_F(SessionStoreTest, test_update_session_rules)
{
  // 1) Create SessionStore
  auto rule_store = std::make_shared<StaticRuleStore>();
  auto session_store = new SessionStore(rule_store);

  // 2) Create a single session and write it into the store
  auto session1 = get_session(sid, rule_store);
  auto session_vec = std::vector<std::unique_ptr<SessionState>>{};
  session_vec.push_back(std::move(session1));
  session_store->create_sessions(imsi, std::move(session_vec));

  // 3) Try to update the session in SessionStore with a rule installation
  auto session_map = session_store->read_sessions(SessionRead{imsi});
  auto session_update = SessionStore::get_default_session_update(session_map);

  auto uc = get_default_update_criteria();
  uc.static_rules_to_install.insert("RULE_asdf");
  RuleLifetime lifetime{
    .activation_time = std::time_t(0),
    .deactivation_time = std::time_t(0),
  };
  uc.new_rule_lifetimes["RULE_asdf"] = lifetime;
  session_update[imsi][sid] = uc;

  auto update_success = session_store->update_sessions(session_update);
  EXPECT_TRUE(update_success);
}

int main(int argc, char **argv)
{
  ::testing::InitGoogleTest(&argc, argv);
  return RUN_ALL_TESTS();
}

} // namespace magma

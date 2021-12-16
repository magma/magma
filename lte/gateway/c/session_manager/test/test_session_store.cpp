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

#include "Consts.h"
#include "magma_logging.h"
#include "includes/MagmaService.h"
#include "ProtobufCreators.h"
#include "RuleStore.h"
#include "SessionID.h"
#include "SessionState.h"
#include "SessionStore.h"
#include "StoredState.h"

using magma::orc8r::MetricsContainer;
using ::testing::Test;

namespace magma {

class SessionStoreTest : public ::testing::Test {
 protected:
  SessionIDGenerator id_gen_;

 protected:
  virtual void SetUp() {
    rule_store = std::make_shared<StaticRuleStore>();
    session_store = std::make_unique<SessionStore>(
        rule_store, std::make_shared<MeteringReporter>());

    session_id_3 = id_gen_.gen_session_id(IMSI2);
    monitoring_key = "mk1";
    monitoring_key2 = "mk2";
    rule_id_1 = "test_rule_1";
    rule_id_2 = "test_rule_2";
    rule_id_3 = "test_rule_3";
    dynamic_rule_id_1 = "dynamic_rule_1";
    dynamic_rule_id_2 = "dynamic_rule_2";

    auto credits = response1.mutable_credits();
    create_credit_update_response(IMSI2, session_id_3, 1, 1000, credits->Add());
  }

  PolicyRule get_dynamic_rule() {
    PolicyRule policy;
    policy.set_id(dynamic_rule_id_1);
    policy.set_priority(10);
    policy.set_tracking_type(PolicyRule::ONLY_OCS);
    return policy;
  }

  std::unique_ptr<SessionState> get_session(const std::string& imsi,
                                            std::string session_id) {
    Teids teid2;
    teid2.set_enb_teid(TEID_2_DL);
    teid2.set_agw_teid(TEID_2_UL);
    return get_session(imsi, session_id, IP2, IPv6_2, teid2, "APN");
  }

  std::unique_ptr<SessionState> get_session(const std::string& imsi,
                                            std::string session_id,
                                            std::string ip_addr,
                                            std::string ipv6_addr, Teids teids,
                                            const std::string& apn) {
    std::string hardware_addr_bytes = {0x0f, 0x10, 0x2e, 0x12, 0x3a, 0x55};
    SessionConfig cfg;
    cfg.common_context = build_common_context(imsi, ip_addr, ipv6_addr, teids,
                                              apn, MSISDN, TGPP_WLAN);
    const auto& wlan_context = build_wlan_context(MAC_ADDR, RADIUS_SESSION_ID);
    cfg.rat_specific_context.mutable_wlan_context()->CopyFrom(wlan_context);
    auto tgpp_context = TgppContext{};
    auto pdp_start_time = 12345;
    auto session = std::make_unique<SessionState>(session_id, cfg, *rule_store,
                                                  pdp_start_time);
    session->set_tgpp_context(tgpp_context, nullptr);
    session->set_fsm_state(SESSION_ACTIVE, nullptr);
    session->set_create_session_response(response1, nullptr);
    return session;
  }

  std::unique_ptr<SessionState> get_lte_session(const std::string& imsi,
                                                std::string session_id) {
    Teids teid;
    teid.set_enb_teid(TEID_1_DL);
    teid.set_agw_teid(TEID_1_UL);
    return get_lte_session(imsi, session_id, IP2, IPv6_1, teid, "APN");
  }

  std::unique_ptr<SessionState> get_lte_session(
      const std::string& imsi, std::string session_id, std::string ip_addr,
      std::string ipv6_addr, Teids teids, const std::string& apn) {
    SessionConfig cfg;
    cfg.common_context = build_common_context(imsi, ip_addr, ipv6_addr, teids,
                                              apn, MSISDN, TGPP_LTE);
    QosInformationRequest qos_info;
    qos_info.set_apn_ambr_dl(32);
    qos_info.set_apn_ambr_dl(64);
    const auto& lte_context =
        build_lte_context(imsi, "", "", "", "", 0, &qos_info);
    cfg.rat_specific_context.mutable_lte_context()->CopyFrom(lte_context);
    auto tgpp_context = TgppContext{};
    auto pdp_start_time = 12345;
    auto session = std::make_unique<SessionState>(session_id, cfg, *rule_store,
                                                  pdp_start_time);
    session->set_tgpp_context(tgpp_context, nullptr);
    session->set_fsm_state(SESSION_ACTIVE, nullptr);
    session->set_create_session_response(response1, nullptr);
    return session;
  }
  UsageMonitoringUpdateResponse get_monitoring_update() {
    UsageMonitoringUpdateResponse response;
    response.set_session_id("sid1");
    response.set_success(true);

    auto monitoring_credit = response.mutable_credit();
    monitoring_credit->set_action(UsageMonitoringCredit_Action_CONTINUE);
    monitoring_credit->set_monitoring_key(monitoring_key);
    monitoring_credit->set_level(SESSION_LEVEL);

    auto units = monitoring_credit->mutable_granted_units();
    auto total = units->mutable_total();
    auto tx = units->mutable_tx();
    auto rx = units->mutable_rx();

    total->set_is_valid(true);
    total->set_volume(1000);
    tx->set_is_valid(true);
    tx->set_volume(1000);
    rx->set_is_valid(true);
    rx->set_volume(1000);

    // Don't set event triggers
    // Don't set result code since the response is already successful
    // Don't set any rule installation/uninstallation
    // Don't set the TgppContext, assume gx_gy_relay disabled
    return response;
  }

  SessionStateUpdateCriteria get_update_criteria() {
    auto update_criteria = SessionStateUpdateCriteria{};

    // Rule installation
    update_criteria.static_rules_to_install = std::set<std::string>{};
    update_criteria.static_rules_to_install.insert(rule_id_1);
    update_criteria.dynamic_rules_to_install = std::vector<PolicyRule>{};
    PolicyRule policy = get_dynamic_rule();
    update_criteria.dynamic_rules_to_install.push_back(policy);
    RuleLifetime lifetime;
    update_criteria.new_rule_lifetimes[rule_id_1] = lifetime;
    update_criteria.new_rule_lifetimes[dynamic_rule_id_1] = lifetime;

    // Monitoring credit installation
    update_criteria.monitor_credit_to_install = StoredMonitorMap{};
    auto monitor2 = StoredMonitor{};
    auto credit2 = StoredSessionCredit{};
    credit2.reporting = false;
    credit2.credit_limit_type = INFINITE_METERED;
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
    credit2.buckets[ALLOWED_FLOOR_TOTAL] = 9;
    credit2.buckets[ALLOWED_FLOOR_TX] = 10;
    credit2.buckets[ALLOWED_FLOOR_RX] = 11;
    monitor2.level = SESSION_LEVEL;
    monitor2.credit = credit2;
    update_criteria.monitor_credit_to_install[monitoring_key2] = monitor2;

    // Monitoring credit updates
    SessionCreditUpdateCriteria monitoring_update{};
    monitoring_update.reauth_state = REAUTH_NOT_NEEDED;
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
    bucket_deltas[ALLOWED_FLOOR_TOTAL] = 9;
    bucket_deltas[ALLOWED_FLOOR_TX] = 10;
    bucket_deltas[ALLOWED_FLOOR_RX] = 11;
    monitoring_update.bucket_deltas = bucket_deltas;

    update_criteria.monitor_credit_map =
        std::unordered_map<std::string, SessionCreditUpdateCriteria>{};
    update_criteria.monitor_credit_map[monitoring_key] = monitoring_update;

    return update_criteria;
  }

  bool is_equal(io::prometheus::client::LabelPair label_pair, const char*& name,
                const char*& value) {
    return label_pair.name().compare(name) == 0 &&
           label_pair.value().compare(value) == 0;
  }

 protected:
  std::string session_id_3;
  std::string monitoring_key;
  std::string monitoring_key2;
  std::string rule_id_1;
  std::string rule_id_2;
  std::string rule_id_3;
  std::string dynamic_rule_id_1;
  std::string dynamic_rule_id_2;
  CreateSessionResponse response1;
  std::unique_ptr<SessionStore> session_store;
  std::shared_ptr<StaticRuleStore> rule_store;
};

TEST_F(SessionStoreTest, test_metering_reporting) {
  // 2) Create a single session and write it into the store
  auto session1 = get_session(IMSI1, SESSION_ID_1);
  auto session_vec = SessionVector{};
  session_vec.push_back(std::move(session1));
  session_store->create_sessions(IMSI1, std::move(session_vec));

  // 3) Try to update the session in SessionStore with a rule installation
  auto session_map = session_store->read_sessions(SessionRead{IMSI1});
  auto session_update = SessionStore::get_default_session_update(session_map);

  auto uc = get_default_update_criteria();
  uc.static_rules_to_install.insert("RULE_asdf");
  RuleLifetime lifetime;
  uc.new_rule_lifetimes["RULE_asdf"] = lifetime;

  // Record some credit usage
  auto DIRECTION_LABEL = "direction";
  auto DIRECTION_UP = "up";
  auto DIRECTION_DOWN = "down";
  auto UPLOADED_BYTES = 5;
  auto DOWNLOADED_BYTES = 7;
  SessionCreditUpdateCriteria credit_uc{};
  credit_uc.bucket_deltas[USED_TX] = UPLOADED_BYTES;
  credit_uc.bucket_deltas[USED_RX] = DOWNLOADED_BYTES;
  uc.monitor_credit_map[monitoring_key] = credit_uc;

  session_update[IMSI1][SESSION_ID_1] = uc;

  auto update_success = session_store->update_sessions(session_update);
  EXPECT_TRUE(update_success);

  // verify if UE traffic metrics are recorded properly
  MetricsContainer resp;
  auto magma_service =
      std::make_shared<service303::MagmaService>("test_service", "1.0");
  magma_service->GetMetrics(nullptr, nullptr, &resp);
  auto reported_metrics = 0;
  for (auto const& fam : resp.family()) {
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
TEST_F(SessionStoreTest, test_read_and_write) {
  // 2) Create bare-bones session for IMSI1
  auto session = get_session(IMSI1, SESSION_ID_1);

  auto uc = get_default_update_criteria();
  RuleLifetime lifetime;
  session->activate_static_rule(rule_id_3, lifetime, nullptr);
  EXPECT_EQ(session->get_session_id(), SESSION_ID_1);
  EXPECT_EQ(session->get_request_number(), 1);
  EXPECT_EQ(session->is_static_rule_installed(rule_id_3), true);
  EXPECT_EQ(session->get_create_session_response().DebugString(),
            response1.DebugString());

  auto monitor_update = get_monitoring_update();
  session->receive_monitor(monitor_update, &uc);

  // Add some used credit
  session->add_to_monitor(monitoring_key, uint64_t(111), uint64_t(333), &uc);
  EXPECT_EQ(session->get_monitor(monitoring_key, USED_TX), 111);
  EXPECT_EQ(session->get_monitor(monitoring_key, USED_RX), 333);

  // 2.1) create an extra session for the same IMSI
  auto session2 = get_session(IMSI1, SESSION_ID_2);

  // 3) Commit session for IMSI1 into SessionStore
  auto sessions = SessionVector{};
  EXPECT_EQ(sessions.size(), 0);
  sessions.push_back(std::move(session));
  sessions.push_back(std::move(session2));
  EXPECT_EQ(sessions.size(), 2);
  session_store->create_sessions(IMSI1, std::move(sessions));

  // 4) Read session for IMSI1 from SessionStore
  SessionRead read_req = {};
  read_req.insert(IMSI1);
  auto session_map = session_store->read_sessions(read_req);

  // 5) Verify that state was written for IMSI1 and has been retrieved.
  EXPECT_EQ(session_map.size(), 1);
  EXPECT_EQ(session_map[IMSI1].size(), 2);
  EXPECT_EQ(session_map[IMSI1].front()->get_request_number(), 1);
  EXPECT_EQ(
      session_map[IMSI1].front()->get_create_session_response().DebugString(),
      response1.DebugString());

  // 6) Make updates to session via SessionUpdateCriteria
  auto update_req = SessionUpdate{};
  update_req[IMSI1] =
      std::unordered_map<std::string, SessionStateUpdateCriteria>{};
  auto update_criteria = get_update_criteria();
  update_criteria.updated_pdp_end_time = 156789;
  update_req[IMSI1][SESSION_ID_1] = update_criteria;

  // 7) Commit updates to SessionStore
  auto success = session_store->update_sessions(update_req);
  EXPECT_TRUE(success);

  // 8) Read in session for IMSI1 again to check that the update was successful
  session_map = session_store->read_sessions(read_req);
  EXPECT_EQ(session_map.size(), 1);
  EXPECT_EQ(session_map[IMSI1].size(), 2);
  EXPECT_EQ(session_map[IMSI1].front()->get_session_id(), SESSION_ID_1);
  // Check installed rules
  EXPECT_EQ(session_map[IMSI1].front()->is_static_rule_installed(rule_id_1),
            true);
  EXPECT_EQ(session_map[IMSI1].front()->is_static_rule_installed(rule_id_3),
            true);
  EXPECT_EQ(session_map[IMSI1].front()->is_static_rule_installed(rule_id_2),
            false);
  EXPECT_EQ(
      session_map[IMSI1].front()->is_dynamic_rule_installed(dynamic_rule_id_1),
      true);
  EXPECT_EQ(
      session_map[IMSI1].front()->is_dynamic_rule_installed(dynamic_rule_id_2),
      false);

  // Check for installation of new monitoring credit
  session_map[IMSI1].front()->set_monitor(
      monitoring_key2,
      Monitor(update_criteria.monitor_credit_to_install[monitoring_key2]),
      nullptr);
  EXPECT_EQ(session_map[IMSI1].front()->get_monitor(monitoring_key2, USED_TX),
            100);
  EXPECT_EQ(session_map[IMSI1].front()->get_monitor(monitoring_key2, USED_RX),
            200);

  // Check monitoring credit usage
  EXPECT_EQ(session_map[IMSI1].front()->get_monitor(monitoring_key, USED_TX),
            222);
  EXPECT_EQ(session_map[IMSI1].front()->get_monitor(monitoring_key, USED_RX),
            666);
  EXPECT_EQ(
      session_map[IMSI1].front()->get_monitor(monitoring_key, ALLOWED_TOTAL),
      1002);
  EXPECT_EQ(session_map[IMSI1].front()->get_monitor(monitoring_key, ALLOWED_TX),
            1003);
  EXPECT_EQ(session_map[IMSI1].front()->get_monitor(monitoring_key, ALLOWED_RX),
            1004);
  EXPECT_EQ(
      session_map[IMSI1].front()->get_monitor(monitoring_key, REPORTING_TX), 5);
  EXPECT_EQ(
      session_map[IMSI1].front()->get_monitor(monitoring_key, REPORTING_RX), 6);
  EXPECT_EQ(
      session_map[IMSI1].front()->get_monitor(monitoring_key, REPORTED_TX), 7);
  EXPECT_EQ(
      session_map[IMSI1].front()->get_monitor(monitoring_key, REPORTED_RX), 8);

  // Check pdp end time update
  EXPECT_EQ(session_map[IMSI1].front()->get_pdp_end_time(), 156789);

  // 11) Delete session 1 for IMSI1
  update_req = SessionUpdate{};
  update_criteria = SessionStateUpdateCriteria{};
  update_criteria.is_session_ended = true;
  update_req[IMSI1][SESSION_ID_1] = update_criteria;
  session_store->update_sessions(update_req);

  // 12) Verify that session 1 on IMSI 1 is gone. Only session 2 is there
  session_map = session_store->read_sessions(read_req);
  EXPECT_EQ(session_map.size(), 1);
  EXPECT_EQ(session_map[IMSI1].size(), 1);
  EXPECT_EQ(session_map[IMSI1][0]->get_session_id(), SESSION_ID_2);

  // 13) Delete session 2 for IMSI1
  update_req = SessionUpdate{};
  update_criteria = SessionStateUpdateCriteria{};
  update_criteria.is_session_ended = true;
  update_req[IMSI1][SESSION_ID_2] = update_criteria;
  session_store->update_sessions(update_req);

  // 12) Verify that the IMSI is gone since has no sessions left
  session_map = session_store->read_all_sessions();
  EXPECT_EQ(session_map.size(), 0);
}

TEST_F(SessionStoreTest, test_sync_request_numbers) {
  // 2) Create bare-bones session for IMSI1
  auto session = get_session(IMSI1, SESSION_ID_1);
  auto uc = get_default_update_criteria();

  // 3) Commit session for IMSI1 into SessionStore
  auto sessions = SessionVector{};
  EXPECT_EQ(sessions.size(), 0);
  sessions.push_back(std::move(session));
  EXPECT_EQ(sessions.size(), 1);
  session_store->create_sessions(IMSI1, std::move(sessions));

  // 4) Read session for IMSI1 from SessionStore
  SessionRead read_req = {};
  read_req.insert(IMSI1);
  auto session_map = session_store->read_sessions(read_req);

  // 5) Verify that state was written for IMSI1 and has been retrieved.
  EXPECT_EQ(session_map.size(), 1);
  EXPECT_EQ(session_map[IMSI1].size(), 1);
  EXPECT_EQ(session_map[IMSI1].front()->get_request_number(), 1);

  // 6) Make updates to session via SessionUpdateCriteria
  auto update_req = SessionUpdate{};
  update_req[IMSI1] =
      std::unordered_map<std::string, SessionStateUpdateCriteria>{};
  auto update_criteria = get_update_criteria();
  update_criteria.request_number_increment = 3;
  update_req[IMSI1][SESSION_ID_1] = update_criteria;

  // 7) Sync updated request_numbers to SessionStore
  session_store->sync_request_numbers(update_req);

  // And then here a gRPC request would be made to another service.
  // The callback would be scheduled onto the event loop, and in the
  // interim, other callbacks can run and make reads to the SessionStore

  // 8) Read in session for IMSI1 again to check that the update was successful
  auto session_map_2 = session_store->read_sessions(read_req);
  EXPECT_EQ(session_map_2.size(), 1);
  EXPECT_EQ(session_map_2[IMSI1].size(), 1);
  EXPECT_EQ(session_map_2[IMSI1].front()->get_session_id(), SESSION_ID_1);
  EXPECT_EQ(session_map_2[IMSI1].front()->get_request_number(), 4);
}

TEST_F(SessionStoreTest, test_get_default_session_update) {
  // 1) Create a SessionMap with a few sessions
  SessionMap session_map = {};
  auto session1 = get_session(IMSI1, SESSION_ID_1);
  auto session2 = get_session(IMSI2, SESSION_ID_2);
  auto session3 = get_session(IMSI2, session_id_3);

  session_map[IMSI1] = SessionVector{};
  session_map[IMSI2] = SessionVector{};

  session_map[IMSI1].push_back(std::move(session1));
  session_map[IMSI2].push_back(std::move(session2));
  session_map[IMSI2].push_back(std::move(session3));

  // 2) Build SessionUpdate
  auto update = SessionStore::get_default_session_update(session_map);
  EXPECT_EQ(update.size(), 2);
  EXPECT_EQ(update[IMSI1].size(), 1);
  EXPECT_EQ(update[IMSI2].size(), 2);
}

TEST_F(SessionStoreTest, test_update_session_rules) {
  // 2) Create a single session and write it into the store
  auto session1 = get_session(IMSI1, SESSION_ID_1);
  auto session_vec = SessionVector{};
  session_vec.push_back(std::move(session1));
  session_store->create_sessions(IMSI1, std::move(session_vec));

  // 3) Try to update the session in SessionStore with a rule installation
  auto session_map = session_store->read_sessions(SessionRead{IMSI1});
  auto session_update = SessionStore::get_default_session_update(session_map);

  auto uc = get_default_update_criteria();
  uc.static_rules_to_install.insert("RULE_asdf");
  RuleLifetime lifetime;

  uc.new_rule_lifetimes["RULE_asdf"] = lifetime;
  session_update[IMSI1][SESSION_ID_1] = uc;

  auto update_success = session_store->update_sessions(session_update);
  EXPECT_TRUE(update_success);
}

TEST_F(SessionStoreTest, test_get_session) {
  // 1) Create a SessionMap with a few sessions
  SessionMap session_map = {};
  // cwag teid
  Teids teid1;
  teid1.set_enb_teid(0);
  teid1.set_agw_teid(0);
  Teids teid2;
  teid2.set_enb_teid(0);
  teid2.set_agw_teid(0);
  // lte teid
  Teids teid3;
  teid3.set_enb_teid(TEID_3_DL);
  teid3.set_agw_teid(TEID_3_UL);
  Teids teid4;
  teid4.set_enb_teid(TEID_4_DL);
  teid4.set_agw_teid(TEID_4_UL);

  // cwag sessions (1 and 2)
  auto session1 = get_session(IMSI1, SESSION_ID_1, IP1, IPv6_1, teid1, "APN1");
  auto session2 = get_session(IMSI1, SESSION_ID_2, IP2, IPv6_2, teid2, "APN2");

  // lte sessions (3 and 4)
  auto session3 =
      get_lte_session(IMSI3, SESSION_ID_3, IP3, IPv6_3, teid3, "APN2");
  auto session4 =
      get_lte_session(IMSI3, SESSION_ID_4, IP4, IPv6_4, teid4, "APN2");

  session_map[IMSI1] = SessionVector{};
  session_map[IMSI1].push_back(std::move(session1));
  session_map[IMSI1].push_back(std::move(session2));
  session_map[IMSI3].push_back(std::move(session3));
  session_map[IMSI3].push_back(std::move(session4));

  // Non-existing subscriber: IMSI4
  SessionSearchCriteria id1_fail1(IMSI4, IMSI_AND_SESSION_ID, SESSION_ID_1);
  SessionSearchCriteria id1_fail2(IMSI4, IMSI_AND_APN, "NON-EXISTING");
  EXPECT_FALSE(session_store->find_session(session_map, id1_fail1));
  EXPECT_FALSE(session_store->find_session(session_map, id1_fail2));

  // Existing subscriber, but non-existing APN/SESSION_ID
  SessionSearchCriteria id1_fail3(IMSI1, IMSI_AND_SESSION_ID, "NON-EXISTING");
  SessionSearchCriteria id1_fail4(IMSI1, IMSI_AND_APN, "NON-EXISTING");
  EXPECT_FALSE(session_store->find_session(session_map, id1_fail3));
  EXPECT_FALSE(session_store->find_session(session_map, id1_fail4));

  // Happy Path! IMSI+SessionID
  SessionSearchCriteria id1_success_sid(IMSI1, IMSI_AND_SESSION_ID,
                                        SESSION_ID_1);
  auto optional_it1 = session_store->find_session(session_map, id1_success_sid);
  EXPECT_TRUE(optional_it1);
  auto& found_session1 = **optional_it1;
  EXPECT_EQ(found_session1->get_session_id(), SESSION_ID_1);

  // Happy Path! IMSI+APN
  SessionSearchCriteria id1_success_apn(IMSI1, IMSI_AND_APN, "APN2");
  auto optional_it2 = session_store->find_session(session_map, id1_success_apn);
  EXPECT_TRUE(optional_it2);
  auto& found_session2 = **optional_it2;
  EXPECT_EQ(found_session2->get_config().common_context.apn(), "APN2");

  // Happy Path! IMSI+UE IPv4
  SessionSearchCriteria id1_success_ipv4(IMSI1, IMSI_AND_UE_IPV4, IP2);
  auto optional_it3 =
      session_store->find_session(session_map, id1_success_ipv4);
  EXPECT_TRUE(optional_it3);
  auto& found_session3 = **optional_it3;
  EXPECT_EQ(found_session3->get_config().common_context.ue_ipv4(), IP2);

  // Happy Path! LTE IMSI+UE IPv4 or IPv6
  SessionSearchCriteria id1_success_ipv46(IMSI3, IMSI_AND_UE_IPV4_OR_IPV6, IP3);
  auto optional_it46 =
      session_store->find_session(session_map, id1_success_ipv46);
  EXPECT_TRUE(optional_it46);
  auto& found_session4 = **optional_it46;
  EXPECT_EQ(found_session4->get_config().common_context.ue_ipv4(), IP3);
  SessionSearchCriteria id1_success_ipv46b(IMSI3, IMSI_AND_UE_IPV4_OR_IPV6,
                                           IPv6_3);
  auto optional_it46b =
      session_store->find_session(session_map, id1_success_ipv46b);
  EXPECT_TRUE(optional_it46b);
  auto& found_session46b = **optional_it46b;
  EXPECT_EQ(found_session46b->get_config().common_context.ue_ipv6(), IPv6_3);

  // Happy Path! cwag IMSI+UE IPv4 or IPv6
  SessionSearchCriteria id1_success_cwag1(IMSI1, IMSI_AND_UE_IPV4_OR_IPV6, "");
  auto optional_it_cwg1 =
      session_store->find_session(session_map, id1_success_cwag1);
  EXPECT_TRUE(optional_it_cwg1);
  auto& found_session_cwag1 = **optional_it_cwg1;
  EXPECT_EQ(found_session_cwag1->get_config().common_context.apn(), "APN1");

  // Happy Path! IMSI+TEID LTE
  SessionSearchCriteria id6_success_sid(IMSI3, IMSI_AND_TEID, TEID_3_DL);
  auto optional_it6 = session_store->find_session(session_map, id6_success_sid);
  EXPECT_TRUE(optional_it6);
  auto& found_session6 = **optional_it6;
  EXPECT_EQ(found_session6->get_session_id(), SESSION_ID_3);

  // Not found IMSI and TEID
  SessionSearchCriteria id7_success_sid(IMSI3, IMSI_AND_TEID, 99);
  auto optional_it7 = session_store->find_session(session_map, id7_success_sid);
  EXPECT_FALSE(optional_it7);
}

int main(int argc, char** argv) {
  ::testing::InitGoogleTest(&argc, argv);
  return RUN_ALL_TESTS();
}

}  // namespace magma

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
#include <folly/io/async/EventBaseManager.h>
#include <gtest/gtest.h>
#include <lte/protos/session_manager.grpc.pb.h>
#include <string.h>
#include <time.h>

#include <chrono>
#include <future>
#include <memory>

#include "Consts.h"
#include "DiameterCodes.h"
#include "LocalEnforcer.h"
#include "magma_logging.h"
#include "MagmaService.h"
#include "Matchers.h"
#include "ProtobufCreators.h"
#include "ServiceRegistrySingleton.h"
#include "SessiondMocks.h"
#include "SessionStore.h"

#define SECONDS_A_DAY 86400

using grpc::ServerContext;
using grpc::Status;
using ::testing::InSequence;
using ::testing::Test;

namespace magma {

Teids teids0;
Teids teids1;
Teids teids2;

class LocalEnforcerTest : public ::testing::Test {
 protected:
  virtual void SetUp() {
    reporter      = std::make_shared<MockSessionReporter>();
    rule_store    = std::make_shared<StaticRuleStore>();
    session_store = std::make_shared<SessionStore>(
        rule_store, std::make_shared<MeteringReporter>());
    pipelined_client     = std::make_shared<MockPipelinedClient>();
    spgw_client          = std::make_shared<MockSpgwServiceClient>();
    aaa_client           = std::make_shared<MockAAAClient>();
    events_reporter      = std::make_shared<MockEventsReporter>();
    auto default_mconfig = get_default_mconfig();
    local_enforcer       = std::make_unique<LocalEnforcer>(
        reporter, rule_store, *session_store, pipelined_client, events_reporter,
        spgw_client, aaa_client, 0, 0, default_mconfig);
    evb = folly::EventBaseManager::get()->getEventBase();
    local_enforcer->attachEventBase(evb);
    session_map   = SessionMap{};
    test_cfg_     = get_default_config("");
    default_cfg_1 = get_default_config(IMSI1);
    default_cfg_2 = get_default_config(IMSI2);

    teids0.set_agw_teid(0);
    teids0.set_enb_teid(0);
    teids1.set_agw_teid(TEID_1_UL);
    teids1.set_enb_teid(TEID_1_DL);
    teids2.set_agw_teid(TEID_2_UL);
    teids2.set_enb_teid(TEID_2_DL);
  }

  virtual void TearDown() { folly::EventBaseManager::get()->clearEventBase(); }

  void run_evb() {
    evb->runAfterDelay([this]() { local_enforcer->stop(); }, 100);
    local_enforcer->start();
  }

  SessionConfig get_default_config(const std::string& imsi) {
    SessionConfig cfg;
    cfg.common_context =
        build_common_context(imsi, IP1, IPv6_1, teids1, APN1, MSISDN, TGPP_LTE);
    QosInformationRequest qos_info;
    qos_info.set_apn_ambr_dl(32);
    qos_info.set_apn_ambr_dl(64);
    const auto& lte_context =
        build_lte_context(IP2, "", "", "", "", BEARER_ID_1, &qos_info);
    cfg.rat_specific_context.mutable_lte_context()->CopyFrom(lte_context);
    return cfg;
  }

  void insert_static_rule(
      uint32_t rating_group, const std::string& m_key,
      const std::string& rule_id) {
    rule_store->insert_rule(create_policy_rule(rule_id, m_key, rating_group));
  }

  void insert_static_rule_with_qos(
      uint32_t rating_group, const std::string& m_key,
      const std::string& rule_id, const int qci) {
    PolicyRule rule = create_policy_rule(rule_id, m_key, rating_group);
    rule.mutable_qos()->set_qci(static_cast<magma::lte::FlowQos_Qci>(qci));
    rule_store->insert_rule(rule);
  }

  void assert_charging_credit(
      SessionMap& session_map, const std::string& imsi,
      const std::string& session_id, Bucket bucket,
      const std::vector<std::pair<uint32_t, uint64_t>>& volumes) {
    EXPECT_TRUE(session_map.find(imsi) != session_map.end());
    bool found = false;
    for (const auto& session : session_map.find(imsi)->second) {
      if (session->get_session_id() == session_id) {
        found = true;
        for (auto& volume_pair : volumes) {
          EXPECT_EQ(
              session->get_charging_credit(volume_pair.first, bucket),
              volume_pair.second);
        }
      }
    }
    EXPECT_TRUE(found);
  }

  void assert_monitor_credit(
      SessionMap& session_map, const std::string& imsi,
      const std::string& session_id, Bucket bucket,
      const std::vector<std::pair<std::string, uint64_t>>& volumes) {
    EXPECT_TRUE(session_map.find(imsi) != session_map.end());
    bool found = false;
    for (const auto& session : session_map.find(imsi)->second) {
      if (session->get_session_id() == session_id) {
        found = true;
        for (auto& volume_pair : volumes) {
          EXPECT_EQ(
              session->get_monitor(volume_pair.first, bucket),
              volume_pair.second);
        }
      }
    }
    EXPECT_TRUE(found);
  }

  void assert_session_is_in_final_state(
      SessionMap& session_map, const std::string& imsi,
      const std::string& session_id, const CreditKey& charging_key,
      bool is_final) {
    EXPECT_TRUE(session_map.find(imsi) != session_map.end());
    for (const auto& session : session_map.find(imsi)->second) {
      if (session->get_session_id() == session_id) {
        optional<FinalActionInfo> fai =
            session->get_final_action_if_final_unit_state(charging_key);
        if (is_final) {
          EXPECT_TRUE(fai);
        } else {
          EXPECT_FALSE(fai);
        }
      }
    }
  }

  RuleToProcess make_rule_to_process(
      const std::string rule_id, const uint32_t agw_teid,
      const uint32_t enb_teid) {
    RuleToProcess to_process;
    to_process.rule.set_id(rule_id);
    to_process.teids.set_enb_teid(enb_teid);
    to_process.teids.set_agw_teid(agw_teid);
    return to_process;
  }

 protected:
  std::shared_ptr<MockSessionReporter> reporter;
  std::shared_ptr<StaticRuleStore> rule_store;
  std::shared_ptr<SessionStore> session_store;
  std::unique_ptr<LocalEnforcer> local_enforcer;
  std::shared_ptr<MockPipelinedClient> pipelined_client;
  std::shared_ptr<MockSpgwServiceClient> spgw_client;
  std::shared_ptr<MockAAAClient> aaa_client;
  std::shared_ptr<MockEventsReporter> events_reporter;
  SessionMap session_map;
  SessionConfig test_cfg_;
  SessionConfig default_cfg_1;
  SessionConfig default_cfg_2;
  folly::EventBase* evb;
};

TEST_F(LocalEnforcerTest, test_init_cwf_session_credit) {
  insert_static_rule(1, "", "rule1");

  SessionConfig test_cwf_cfg;
  Teids teids;
  test_cwf_cfg.common_context =
      build_common_context(IMSI1, "", "", teids, "", "", TGPP_WLAN);
  const auto& wlan = build_wlan_context(MAC_ADDR, RADIUS_SESSION_ID);
  test_cwf_cfg.rat_specific_context.mutable_wlan_context()->CopyFrom(wlan);

  CreateSessionResponse response;
  auto credits = response.mutable_credits();
  create_credit_update_response(IMSI1, SESSION_ID_1, 1, 1024, credits->Add());
  EXPECT_CALL(
      *pipelined_client, activate_flows_for_rules(
                             IMSI1, testing::_, testing::_, testing::_,
                             test_cwf_cfg.common_context.msisdn(), testing::_,
                             CheckRuleCount(0), testing::_))
      .Times(1);

  EXPECT_CALL(
      *pipelined_client, update_ipfix_flow(
                             testing::_, testing::_, testing::_, testing::_,
                             testing::_, testing::_))
      .Times(1);

  local_enforcer->init_session(
      session_map, IMSI1, SESSION_ID_1, test_cwf_cfg, response);

  local_enforcer->update_tunnel_ids(
      session_map,
      create_update_tunnel_ids_request(IMSI1, BEARER_ID_1, teids0));

  assert_charging_credit(
      session_map, IMSI1, SESSION_ID_1, ALLOWED_TOTAL, {{1, 1024}});
}

TEST_F(LocalEnforcerTest, test_init_infinite_metered_credit) {
  insert_static_rule(1, "", "rule1");

  CreateSessionResponse response;
  auto credits = response.mutable_credits();
  create_credit_update_response(
      IMSI1, SESSION_ID_1, 1, INFINITE_METERED, credits->Add());

  StaticRuleInstall rule1;
  rule1.set_rule_id("rule1");
  auto rules_to_install = response.mutable_static_rules();
  rules_to_install->Add()->CopyFrom(rule1);

  // Expect rule1 to be activated
  test_cfg_.common_context.mutable_sid()->set_id(IMSI1);
  EXPECT_CALL(
      *pipelined_client, activate_flows_for_rules(
                             IMSI1, IP1, IPv6_1, CheckTeids(teids1),
                             test_cfg_.common_context.msisdn(), testing::_,
                             CheckRuleCount(1), testing::_))
      .Times(1);
  local_enforcer->init_session(
      session_map, IMSI1, SESSION_ID_1, test_cfg_, response);

  local_enforcer->update_tunnel_ids(
      session_map,
      create_update_tunnel_ids_request(IMSI1, BEARER_ID_1, teids1));

  assert_charging_credit(
      session_map, IMSI1, SESSION_ID_1, ALLOWED_TOTAL, {{1, 0}});
}

TEST_F(LocalEnforcerTest, test_init_no_credit) {
  insert_static_rule(1, "", "rule1");

  CreateSessionResponse response;
  create_credit_update_response(
      IMSI1, SESSION_ID_1, 1, 0, response.mutable_credits()->Add());

  StaticRuleInstall rule1;
  rule1.set_rule_id("rule1");
  auto rules_to_install = response.mutable_static_rules();
  rules_to_install->Add()->CopyFrom(rule1);

  // Expect rule1 to be activated even if the GSU is all 0s
  test_cfg_.common_context.mutable_sid()->set_id(IMSI1);
  EXPECT_CALL(
      *pipelined_client,
      activate_flows_for_rules(
          IMSI1, test_cfg_.common_context.ue_ipv4(),
          test_cfg_.common_context.ue_ipv6(), CheckTeids(teids1),
          test_cfg_.common_context.msisdn(), testing::_, CheckRuleCount(1),
          testing::_))
      .Times(1);
  local_enforcer->init_session(
      session_map, IMSI1, SESSION_ID_1, test_cfg_, response);

  local_enforcer->update_tunnel_ids(
      session_map,
      create_update_tunnel_ids_request(IMSI1, BEARER_ID_1, teids1));

  assert_charging_credit(
      session_map, IMSI1, SESSION_ID_1, ALLOWED_TOTAL, {{1, 0}});
}

TEST_F(LocalEnforcerTest, test_init_session_credit) {
  insert_static_rule(1, "", "rule1");

  CreateSessionResponse response;
  auto credits = response.mutable_credits();
  create_credit_update_response(IMSI1, SESSION_ID_1, 1, 1024, credits->Add());

  EXPECT_CALL(
      *pipelined_client,
      activate_flows_for_rules(
          testing::_, testing::_, testing::_, CheckTeids(teids1),
          test_cfg_.common_context.msisdn(), testing::_, testing::_,
          testing::_))
      .Times(1);
  ;
  local_enforcer->init_session(
      session_map, IMSI1, SESSION_ID_1, test_cfg_, response);

  local_enforcer->update_tunnel_ids(
      session_map,
      create_update_tunnel_ids_request(IMSI1, BEARER_ID_1, teids1));

  assert_charging_credit(
      session_map, IMSI1, SESSION_ID_1, ALLOWED_TOTAL, {{1, 1024}});
}

TEST_F(LocalEnforcerTest, test_single_record) {
  // insert initial session credit
  CreateSessionResponse response;
  create_credit_update_response(
      IMSI1, SESSION_ID_1, 1, 1024, response.mutable_credits()->Add());
  local_enforcer->init_session(
      session_map, IMSI1, SESSION_ID_1, test_cfg_, response);
  local_enforcer->update_tunnel_ids(
      session_map,
      create_update_tunnel_ids_request(IMSI1, BEARER_ID_1, teids1));

  insert_static_rule(1, "", "rule1");
  RuleRecordTable table;
  auto record_list = table.mutable_records();
  create_rule_record(
      IMSI1, teids1.agw_teid(), "rule1", 16, 32, record_list->Add());

  auto update = SessionStore::get_default_session_update(session_map);
  local_enforcer->aggregate_records(session_map, table, update);
  assert_charging_credit(session_map, IMSI1, SESSION_ID_1, USED_RX, {{1, 16}});
  assert_charging_credit(session_map, IMSI1, SESSION_ID_1, USED_TX, {{1, 32}});
  assert_charging_credit(
      session_map, IMSI1, SESSION_ID_1, ALLOWED_TOTAL, {{1, 1024}});

  EXPECT_EQ(update.size(), 1);
  EXPECT_EQ(update[IMSI1][SESSION_ID_1].charging_credit_map.size(), 1);
  EXPECT_EQ(
      update[IMSI1][SESSION_ID_1].charging_credit_map[1].bucket_deltas[USED_RX],
      16);
  EXPECT_EQ(
      update[IMSI1][SESSION_ID_1].charging_credit_map[1].bucket_deltas[USED_TX],
      32);
}

TEST_F(LocalEnforcerTest, test_aggregate_records_mixed_ips) {
  CreateSessionResponse response;
  create_credit_update_response(
      IMSI1, SESSION_ID_1, 1, 1024, response.mutable_credits()->Add());
  create_credit_update_response(
      IMSI1, SESSION_ID_1, 2, 1024, response.mutable_credits()->Add());
  local_enforcer->init_session(
      session_map, IMSI1, SESSION_ID_1, test_cfg_, response);
  local_enforcer->update_tunnel_ids(
      session_map,
      create_update_tunnel_ids_request(IMSI1, BEARER_ID_1, teids1));

  insert_static_rule(1, "", "rule1");
  insert_static_rule(1, "", "rule2");
  insert_static_rule(2, "", "rule3");
  RuleRecordTable table;
  auto record_list = table.mutable_records();
  // ipv4 usage
  create_rule_record(
      IMSI1, teids1.agw_teid(), "rule1", 10, 20, record_list->Add());
  // ipv6 usage for the same charging key and subscriber
  create_rule_record(
      IMSI1, teids1.agw_teid(), "rule2", 5, 15, record_list->Add());
  create_rule_record(
      IMSI1, teids1.agw_teid(), "rule3", 100, 150, record_list->Add());

  auto update = SessionStore::get_default_session_update(session_map);
  local_enforcer->aggregate_records(session_map, table, update);

  assert_charging_credit(
      session_map, IMSI1, SESSION_ID_1, USED_RX, {{1, 15}, {2, 100}});
  assert_charging_credit(
      session_map, IMSI1, SESSION_ID_1, USED_TX, {{1, 35}, {2, 150}});

  EXPECT_EQ(update[IMSI1][SESSION_ID_1].charging_credit_map.size(), 2);
  EXPECT_EQ(
      update[IMSI1][SESSION_ID_1].charging_credit_map[1].bucket_deltas[USED_RX],
      15);
  EXPECT_EQ(
      update[IMSI1][SESSION_ID_1].charging_credit_map[1].bucket_deltas[USED_TX],
      35);
  EXPECT_EQ(
      update[IMSI1][SESSION_ID_1].charging_credit_map[2].bucket_deltas[USED_RX],
      100);
  EXPECT_EQ(
      update[IMSI1][SESSION_ID_1].charging_credit_map[2].bucket_deltas[USED_TX],
      150);
}

TEST_F(LocalEnforcerTest, test_aggregate_records_for_termination) {
  CreateSessionResponse response;
  create_credit_update_response(
      IMSI1, SESSION_ID_1, 1, 1024, response.mutable_credits()->Add());
  create_credit_update_response(
      IMSI1, SESSION_ID_1, 2, 1024, response.mutable_credits()->Add());
  local_enforcer->init_session(
      session_map, IMSI1, SESSION_ID_1, get_default_config(IMSI1), response);
  local_enforcer->update_tunnel_ids(
      session_map,
      create_update_tunnel_ids_request(IMSI1, BEARER_ID_1, teids1));

  insert_static_rule(1, "", "rule1");
  insert_static_rule(1, "", "rule2");
  insert_static_rule(2, "", "rule3");

  auto update = SessionStore::get_default_session_update(session_map);
  local_enforcer->handle_termination_from_access(
      session_map, IMSI1, APN1, update);

  RuleRecordTable table;
  auto record_list = table.mutable_records();
  create_rule_record(
      IMSI1, teids1.agw_teid(), "rule1", 10, 20, record_list->Add());
  create_rule_record(
      IMSI1, teids1.agw_teid(), "rule2", 5, 15, record_list->Add());
  create_rule_record(
      IMSI1, teids1.agw_teid(), "rule3", 100, 150, record_list->Add());

  EXPECT_CALL(
      *reporter,
      report_terminate_session(CheckTerminateRequestCount(IMSI1, 0, 2), _))
      .Times(1);
  local_enforcer->aggregate_records(session_map, table, update);

  RuleRecordTable empty_table;
  local_enforcer->aggregate_records(session_map, empty_table, update);
}

TEST_F(LocalEnforcerTest, test_collect_updates) {
  CreateSessionResponse response;
  create_credit_update_response(
      IMSI1, SESSION_ID_1, 1, 3072, response.mutable_credits()->Add());
  local_enforcer->init_session(
      session_map, IMSI1, SESSION_ID_1, test_cfg_, response);
  local_enforcer->update_tunnel_ids(
      session_map,
      create_update_tunnel_ids_request(IMSI1, BEARER_ID_1, teids1));

  insert_static_rule(1, "", "rule1");

  std::vector<std::unique_ptr<ServiceAction>> actions;
  auto update = SessionStore::get_default_session_update(session_map);
  auto empty_update =
      local_enforcer->collect_updates(session_map, actions, update);
  local_enforcer->execute_actions(session_map, actions, update);
  EXPECT_EQ(empty_update.updates_size(), 0);

  RuleRecordTable table;
  auto record_list = table.mutable_records();
  create_rule_record(
      IMSI1, teids1.agw_teid(), "rule1", 1024, 2048, record_list->Add());

  local_enforcer->aggregate_records(session_map, table, update);
  actions.clear();
  auto session_update =
      local_enforcer->collect_updates(session_map, actions, update);
  local_enforcer->execute_actions(session_map, actions, update);
  EXPECT_EQ(session_update.updates_size(), 1);
  assert_charging_credit(
      session_map, IMSI1, SESSION_ID_1, REPORTING_RX, {{1, 1024}});
  assert_charging_credit(
      session_map, IMSI1, SESSION_ID_1, REPORTING_TX, {{1, 2048}});
  EXPECT_EQ(update[IMSI1][SESSION_ID_1].charging_credit_map.size(), 1);
  // UpdateCriteria does not store REPORTING_RX / REPORTING_TX
  EXPECT_EQ(
      update[IMSI1][SESSION_ID_1]
          .charging_credit_map[1]
          .bucket_deltas[REPORTING_RX],
      0);
  EXPECT_EQ(
      update[IMSI1][SESSION_ID_1]
          .charging_credit_map[1]
          .bucket_deltas[REPORTING_TX],
      0);
}

TEST_F(LocalEnforcerTest, test_update_session_credits_and_rules) {
  insert_static_rule(1, "", "rule1");

  CreateSessionResponse response;
  auto credits = response.mutable_credits();
  create_credit_update_response(IMSI1, SESSION_ID_1, 1, 2048, credits->Add());
  local_enforcer->init_session(
      session_map, IMSI1, SESSION_ID_1, test_cfg_, response);
  local_enforcer->update_tunnel_ids(
      session_map,
      create_update_tunnel_ids_request(IMSI1, BEARER_ID_1, teids1));

  assert_charging_credit(
      session_map, IMSI1, SESSION_ID_1, ALLOWED_TOTAL, {{1, 2048}});

  insert_static_rule(1, "1", "rule1");

  RuleRecordTable table;
  auto record_list = table.mutable_records();
  create_rule_record(
      IMSI1, teids1.agw_teid(), "rule1", 1024, 1024, record_list->Add());
  auto update = SessionStore::get_default_session_update(session_map);
  local_enforcer->aggregate_records(session_map, table, update);

  session_store->create_sessions(IMSI1, std::move(session_map[IMSI1]));

  UpdateSessionResponse update_response;
  auto credit_updates_response = update_response.mutable_responses();
  create_credit_update_response(
      IMSI1, SESSION_ID_1, 1, 24, credit_updates_response->Add());

  auto monitor_updates_response =
      update_response.mutable_usage_monitor_responses();
  create_monitor_update_response(
      IMSI1, SESSION_ID_1, "1", MonitoringLevel::PCC_RULE_LEVEL, 2048,
      monitor_updates_response->Add());

  // add a credit and monitor update for a different session (this should not
  // impact session 1234. This helps to test feature to prevent adding credits
  // for a specific session with same RG/mKey to other sessions of the same IMSI
  // with the same RG/mKey
  create_credit_update_response(
      IMSI1, SESSION_ID_2, 1, 30000, credit_updates_response->Add());
  create_monitor_update_response(
      IMSI1, SESSION_ID_2, "1", MonitoringLevel::PCC_RULE_LEVEL, 40000,
      monitor_updates_response->Add());

  session_map = session_store->read_sessions(SessionRead{IMSI1});
  local_enforcer->update_session_credits_and_rules(
      session_map, update_response, update);

  assert_charging_credit(
      session_map, IMSI1, SESSION_ID_1, ALLOWED_TOTAL, {{1, 2072}});
  assert_monitor_credit(
      session_map, IMSI1, SESSION_ID_1, ALLOWED_TOTAL, {{"1", 2048}});
}

TEST_F(LocalEnforcerTest, test_update_session_credits_and_rules_with_failure) {
  insert_static_rule(0, "1", "rule1");

  CreateSessionResponse response;
  auto rules = response.mutable_static_rules()->Add();
  rules->mutable_rule_id()->assign("rule1");
  rules->mutable_activation_time()->set_seconds(0);
  rules->mutable_deactivation_time()->set_seconds(0);

  auto monitor_updates = response.mutable_usage_monitors();
  create_monitor_update_response(
      IMSI1, SESSION_ID_1, "1", MonitoringLevel::PCC_RULE_LEVEL, 1024,
      monitor_updates->Add());
  local_enforcer->init_session(
      session_map, IMSI1, SESSION_ID_1, default_cfg_1, response);
  local_enforcer->update_tunnel_ids(
      session_map,
      create_update_tunnel_ids_request(IMSI1, BEARER_ID_1, teids1));

  assert_monitor_credit(
      session_map, IMSI1, SESSION_ID_1, ALLOWED_TOTAL, {{"1", 1024}});
  assert_charging_credit(session_map, IMSI1, SESSION_ID_1, ALLOWED_TOTAL, {});

  // receive usages from pipelined
  RuleRecordTable table;
  auto record_list = table.mutable_records();
  create_rule_record(
      IMSI1, teids1.agw_teid(), "rule1", 10, 20, record_list->Add());
  auto update = SessionStore::get_default_session_update(session_map);
  local_enforcer->aggregate_records(session_map, table, update);
  assert_monitor_credit(session_map, IMSI1, SESSION_ID_1, USED_RX, {{"1", 10}});
  assert_monitor_credit(session_map, IMSI1, SESSION_ID_1, USED_TX, {{"1", 20}});

  UpdateSessionResponse update_response;
  auto monitor_updates_responses =
      update_response.mutable_usage_monitor_responses();
  auto monitor_response = monitor_updates_responses->Add();
  create_monitor_update_response(
      IMSI1, SESSION_ID_1, "1", MonitoringLevel::PCC_RULE_LEVEL, 2048,
      monitor_response);
  monitor_response->set_success(false);
  monitor_response->set_result_code(
      DIAMETER_USER_UNKNOWN);  // USER_UNKNOWN permanent failure

  // the request should has no rules so PipelineD deletes all rules
  EXPECT_CALL(
      *pipelined_client,
      deactivate_flows_for_rules_for_termination(
          IMSI1, testing::_, testing::_, testing::_, testing::_))
      .Times(1);
  local_enforcer->update_session_credits_and_rules(
      session_map, update_response, update);

  // expect no update to credit
  assert_monitor_credit(
      session_map, IMSI1, SESSION_ID_1, ALLOWED_TOTAL, {{"1", 1024}});
  assert_charging_credit(session_map, IMSI1, SESSION_ID_1, ALLOWED_TOTAL, {});
}

TEST_F(LocalEnforcerTest, test_terminate_credit) {
  CreateSessionResponse response;
  create_credit_update_response(
      IMSI1, SESSION_ID_1, 1, 1024, response.mutable_credits()->Add());
  create_credit_update_response(
      IMSI1, SESSION_ID_1, 2, 2048, response.mutable_credits()->Add());
  CreateSessionResponse response2;
  create_credit_update_response(
      IMSI2, SESSION_ID_2, 1, 4096, response2.mutable_credits()->Add());

  session_map = session_store->read_sessions(SessionRead{IMSI1});
  local_enforcer->init_session(
      session_map, IMSI1, SESSION_ID_1, get_default_config(IMSI1), response);
  local_enforcer->update_tunnel_ids(
      session_map,
      create_update_tunnel_ids_request(IMSI1, BEARER_ID_1, teids1));
  session_store->create_sessions(IMSI1, std::move(session_map[IMSI1]));

  session_map = session_store->read_sessions(SessionRead{IMSI2});
  local_enforcer->init_session(
      session_map, IMSI2, SESSION_ID_2, get_default_config(IMSI2), response2);
  local_enforcer->update_tunnel_ids(
      session_map,
      create_update_tunnel_ids_request(IMSI2, BEARER_ID_2, teids2));
  session_store->create_sessions(IMSI2, std::move(session_map[IMSI2]));

  session_map = session_store->read_sessions(SessionRead{IMSI1});
  auto update = SessionStore::get_default_session_update(session_map);

  EXPECT_CALL(
      *reporter,
      report_terminate_session(CheckTerminateRequestCount(IMSI1, 0, 2), _))
      .Times(1);
  local_enforcer->handle_termination_from_access(
      session_map, IMSI1, test_cfg_.common_context.apn(), update);

  // pipelined still reports default drop flow rule when all flows are removed
  RuleRecordTable only_drop_rule_table;
  auto record_list = only_drop_rule_table.mutable_records();
  create_rule_record(
      IMSI1, teids1.agw_teid(), "internal_default_drop_flow_rule", 0, 0,
      record_list->Add());

  local_enforcer->aggregate_records(session_map, only_drop_rule_table, update);
  run_evb();
  run_evb();
  bool success = session_store->update_sessions(update);
  EXPECT_TRUE(success);

  // No longer in system
  session_map = session_store->read_sessions(SessionRead{IMSI1});
  EXPECT_EQ(session_map[IMSI1].size(), 0);
}

TEST_F(LocalEnforcerTest, test_terminate_credit_during_reporting) {
  CreateSessionResponse response;
  create_credit_update_response(
      IMSI1, SESSION_ID_1, 1, 3072, response.mutable_credits()->Add());
  create_credit_update_response(
      IMSI1, SESSION_ID_1, 2, 2048, response.mutable_credits()->Add());
  create_monitor_update_response(
      IMSI1, SESSION_ID_1, "m1", MonitoringLevel::PCC_RULE_LEVEL, 1024,
      response.mutable_usage_monitors()->Add());
  local_enforcer->init_session(
      session_map, IMSI1, SESSION_ID_1, default_cfg_1, response);
  local_enforcer->update_tunnel_ids(
      session_map,
      create_update_tunnel_ids_request(IMSI1, BEARER_ID_1, teids1));
  insert_static_rule(1, "", "rule1");
  insert_static_rule(2, "", "rule2");

  // Insert record for key 1
  RuleRecordTable table;
  auto record_list = table.mutable_records();
  create_rule_record(
      IMSI1, teids1.agw_teid(), "rule1", 1024, 2048, record_list->Add());
  auto update = SessionStore::get_default_session_update(session_map);
  local_enforcer->aggregate_records(session_map, table, update);

  // Collect updates to put key 1 into reporting state
  std::vector<std::unique_ptr<ServiceAction>> actions;
  auto usage_updates =
      local_enforcer->collect_updates(session_map, actions, update);

  local_enforcer->execute_actions(session_map, actions, update);
  assert_charging_credit(
      session_map, IMSI1, SESSION_ID_1, REPORTING_RX, {{1, 1024}});

  session_store->create_sessions(IMSI1, std::move(session_map[IMSI1]));

  session_map = session_store->read_sessions(SessionRead{IMSI1});
  // Collecting terminations should key 1 anyways during reporting
  local_enforcer->handle_termination_from_access(
      session_map, IMSI1, APN1, update);

  EXPECT_CALL(
      *reporter,
      report_terminate_session(CheckTerminateRequestCount(IMSI1, 1, 2), _))
      .Times(1);

  RuleRecordTable empty_table;
  local_enforcer->aggregate_records(session_map, empty_table, update);
  run_evb();

  // pipelined still reports default drop rule and any packet dropped comming
  // from the previously unistalled rule. If we see only dropped traffic is
  // observed and session is scheduled for termination, session should be
  // terminated
  RuleRecordTable only_drop_rule_table;
  record_list = only_drop_rule_table.mutable_records();
  create_rule_record(
      IMSI1, teids1.agw_teid(), "rule1", 0, 0, 1000, 2000, record_list->Add());
  create_rule_record(
      IMSI1, teids1.agw_teid(), "internal_default_drop_flow_rule", 0, 0,
      record_list->Add());

  local_enforcer->aggregate_records(session_map, only_drop_rule_table, update);
  run_evb();
  run_evb();
  bool success = session_store->update_sessions(update);
  EXPECT_TRUE(success);

  // No longer in system
  session_map = session_store->read_sessions(SessionRead{IMSI1});
  EXPECT_EQ(session_map[IMSI1].size(), 0);
}

TEST_F(LocalEnforcerTest, test_sync_sessions_on_restart) {
  CreateSessionResponse response;
  create_credit_update_response(
      IMSI1, SESSION_ID_1, 1, 1024, true, response.mutable_credits()->Add());
  local_enforcer->init_session(
      session_map, IMSI1, SESSION_ID_1, test_cfg_, response);
  local_enforcer->update_tunnel_ids(
      session_map,
      create_update_tunnel_ids_request(IMSI1, BEARER_ID_1, teids1));

  insert_static_rule(1, "", "rule1");
  insert_static_rule(1, "", "rule2");
  insert_static_rule(1, "", "rule3");
  insert_static_rule(1, "", "rule4");

  EXPECT_EQ(session_map[IMSI1].size(), 1);
  bool success =
      session_store->create_sessions(IMSI1, std::move(session_map[IMSI1]));
  EXPECT_TRUE(success);

  auto session_map_2 = session_store->read_sessions(SessionRead{IMSI1});
  auto session_update =
      session_store->get_default_session_update(session_map_2);
  EXPECT_EQ(session_map_2[IMSI1].size(), 1);

  RuleLifetime lifetime1 = {
      .activation_time   = std::time_t(0),
      .deactivation_time = std::time_t(5),
  };
  RuleLifetime lifetime2 = {
      .activation_time   = std::time_t(5),
      .deactivation_time = std::time_t(10),
  };
  RuleLifetime lifetime3 = {
      .activation_time   = std::time_t(10),
      .deactivation_time = std::time_t(15),
  };
  RuleLifetime lifetime4 = {
      .activation_time   = std::time_t(15),
      .deactivation_time = std::time_t(20),
  };
  auto& uc    = session_update[IMSI1][SESSION_ID_1];
  uint32_t v1 = session_map_2[IMSI1]
                    .front()
                    ->activate_static_rule("rule1", lifetime1, uc)
                    .version;
  session_map_2[IMSI1].front()->schedule_static_rule("rule2", lifetime2, uc);
  session_map_2[IMSI1].front()->schedule_static_rule("rule3", lifetime3, uc);
  session_map_2[IMSI1].front()->schedule_static_rule("rule4", lifetime4, uc);

  EXPECT_EQ(v1, 1);

  EXPECT_TRUE(uc.policy_version_and_stats);
  EXPECT_EQ((*uc.policy_version_and_stats)["rule1"].current_version, 1);

  EXPECT_EQ(uc.static_rules_to_install.count("rule1"), 1);
  EXPECT_EQ(uc.new_scheduled_static_rules.count("rule2"), 1);
  EXPECT_EQ(uc.new_scheduled_static_rules.count("rule3"), 1);
  EXPECT_EQ(uc.new_scheduled_static_rules.count("rule4"), 1);

  PolicyRule d1, d2, d3, d4;
  d1.set_id("dynamic_rule1");
  d2.set_id("dynamic_rule2");
  d3.set_id("dynamic_rule3");
  d4.set_id("dynamic_rule4");

  session_map_2[IMSI1].front()->insert_dynamic_rule(d1, lifetime1, uc);
  session_map_2[IMSI1].front()->schedule_dynamic_rule(d2, lifetime2, uc);
  session_map_2[IMSI1].front()->schedule_dynamic_rule(d3, lifetime3, uc);
  session_map_2[IMSI1].front()->schedule_dynamic_rule(d4, lifetime4, uc);

  EXPECT_EQ(uc.dynamic_rules_to_install.size(), 1);
  EXPECT_EQ(uc.new_scheduled_dynamic_rules.size(), 3);

  success = session_store->update_sessions(session_update);
  EXPECT_TRUE(success);

  local_enforcer->sync_sessions_on_restart(std::time_t(12));

  session_map_2  = session_store->read_sessions(SessionRead{IMSI1});
  session_update = session_store->get_default_session_update(session_map_2);
  EXPECT_EQ(session_map_2[IMSI1].size(), 1);

  auto& session = session_map_2[IMSI1].front();
  EXPECT_FALSE(session->is_static_rule_installed("rule1"));
  EXPECT_FALSE(session->is_static_rule_installed("rule2"));
  EXPECT_TRUE(session->is_static_rule_installed("rule3"));
  EXPECT_FALSE(session->is_static_rule_installed("rule4"));

  EXPECT_FALSE(session->is_dynamic_rule_installed("dynamic_rule1"));
  EXPECT_FALSE(session->is_dynamic_rule_installed("dynamic_rule2"));
  EXPECT_TRUE(session->is_dynamic_rule_installed("dynamic_rule3"));
  EXPECT_FALSE(session->is_dynamic_rule_installed("dynamic_rule4"));
}

TEST_F(LocalEnforcerTest, test_sync_sessions_on_restart_revalidation_timer) {
  auto pdp_start_time = 12345;
  magma::lte::TgppContext tgpp_ctx;
  CreateSessionResponse response;
  create_credit_update_response(
      IMSI1, SESSION_ID_1, 1, 1024, true, response.mutable_credits()->Add());
  auto session_state = std::make_unique<SessionState>(
      IMSI1, SESSION_ID_1, default_cfg_1, *rule_store, tgpp_ctx, pdp_start_time,
      response);

  // manually place revalidation timer
  SessionStateUpdateCriteria uc;
  session_state->add_new_event_trigger(REVALIDATION_TIMEOUT, uc);
  EXPECT_EQ(uc.is_pending_event_triggers_updated, true);
  EXPECT_EQ(uc.pending_event_triggers[REVALIDATION_TIMEOUT], false);
  google::protobuf::Timestamp time;
  time.set_seconds(0);
  session_state->set_revalidation_time(time, uc);

  session_map[IMSI1] = SessionVector();
  session_map[IMSI1].push_back(std::move(session_state));

  EXPECT_EQ(session_map[IMSI1].size(), 1);
  bool success =
      session_store->create_sessions(IMSI1, std::move(session_map[IMSI1]));
  EXPECT_TRUE(success);

  local_enforcer->sync_sessions_on_restart(std::time_t(0));
  // sync_sessions_on_restart should call schedule_revalidation which puts
  // two things on the event loop
  evb->loopOnce();
  evb->loopOnce();
  // We expect that the event trigger will now be marked as ready to be acted on
  auto session_map_2 = session_store->read_sessions(SessionRead{IMSI1});
  EXPECT_EQ(session_map_2[IMSI1].size(), 1);

  auto& session = session_map_2[IMSI1].front();
  auto events   = session->get_event_triggers();
  EXPECT_EQ(events[REVALIDATION_TIMEOUT], true);
}

TEST_F(LocalEnforcerTest, test_final_unit_handling) {
  CreateSessionResponse response;
  create_credit_update_response(
      IMSI1, SESSION_ID_1, 1, 1024, true, response.mutable_credits()->Add());
  local_enforcer->init_session(
      session_map, IMSI1, SESSION_ID_1, default_cfg_1, response);
  local_enforcer->update_tunnel_ids(
      session_map,
      create_update_tunnel_ids_request(IMSI1, BEARER_ID_1, teids1));
  insert_static_rule(1, "", "rule1");
  insert_static_rule(1, "", "rule2");

  // Insert record for key 1
  RuleRecordTable table;
  auto record_list = table.mutable_records();
  create_rule_record(
      IMSI1, teids1.agw_teid(), "rule1", 1024, 2048, record_list->Add());
  create_rule_record(
      IMSI1, teids1.agw_teid(), "rule2", 1024, 2048, record_list->Add());
  auto update = SessionStore::get_default_session_update(session_map);
  local_enforcer->aggregate_records(session_map, table, update);

  // the request should has no rules so PipelineD deletes all rules
  EXPECT_CALL(
      *pipelined_client,
      deactivate_flows_for_rules_for_termination(
          testing::_, testing::_, testing::_, testing::_, testing::_))
      .Times(1);
  // Since this is a termination triggered by SessionD/Core (quota exhaustion
  // + FUA-Terminate), we expect MME to be notified to delete the bearer
  // created on session creation
  EXPECT_CALL(
      *spgw_client, delete_default_bearer(IMSI1, testing::_, testing::_));
  // call collect_updates to trigger actions
  std::vector<std::unique_ptr<ServiceAction>> actions;
  auto usage_updates =
      local_enforcer->collect_updates(session_map, actions, update);
  local_enforcer->execute_actions(session_map, actions, update);
}

TEST_F(LocalEnforcerTest, test_cwf_final_unit_handling) {
  CreateSessionResponse response;
  create_credit_update_response(
      IMSI1, SESSION_ID_1, 1, 1024, true, response.mutable_credits()->Add());
  auto monitors = response.mutable_usage_monitors();
  auto monitor  = monitors->Add();
  create_monitor_update_response(
      IMSI1, SESSION_ID_1, "m1", MonitoringLevel::PCC_RULE_LEVEL, 1024,
      monitor);
  StaticRuleInstall static_rule_install;
  static_rule_install.set_rule_id("rule3");
  response.mutable_static_rules()->Add()->CopyFrom(static_rule_install);

  insert_static_rule(0, "m1", "rule3");

  SessionConfig test_cwf_cfg;
  test_cwf_cfg.common_context =
      build_common_context(IMSI1, "", "", teids0, "", "", TGPP_WLAN);
  const auto& wlan = build_wlan_context(MAC_ADDR, RADIUS_SESSION_ID);
  test_cwf_cfg.common_context.set_rat_type(TGPP_WLAN);
  test_cwf_cfg.rat_specific_context.mutable_wlan_context()->CopyFrom(wlan);

  local_enforcer->init_session(
      session_map, IMSI1, SESSION_ID_1, test_cwf_cfg, response);
  local_enforcer->update_tunnel_ids(
      session_map,
      create_update_tunnel_ids_request(IMSI1, BEARER_ID_1, teids1));
  insert_static_rule(1, "", "rule1");
  insert_static_rule(1, "", "rule2");

  // Insert record for key 1
  RuleRecordTable table;
  auto record_list = table.mutable_records();
  create_rule_record(
      IMSI1, teids1.agw_teid(), "rule1", 1024, 2048, record_list->Add());
  create_rule_record(
      IMSI1, teids1.agw_teid(), "rule2", 1024, 2048, record_list->Add());
  auto update = SessionStore::get_default_session_update(session_map);
  local_enforcer->aggregate_records(session_map, table, update);

  // the request should has no rules so PipelineD deletes all rules
  EXPECT_CALL(
      *pipelined_client,
      deactivate_flows_for_rules_for_termination(
          testing::_, testing::_, testing::_, testing::_, testing::_))
      .Times(1);

  EXPECT_CALL(*aaa_client, terminate_session(testing::_, testing::_))
      .Times(1)
      .WillOnce(testing::Return(true));
  // call collect_updates to trigger actions
  std::vector<std::unique_ptr<ServiceAction>> actions;
  auto usage_updates =
      local_enforcer->collect_updates(session_map, actions, update);
  local_enforcer->execute_actions(session_map, actions, update);
}

TEST_F(LocalEnforcerTest, test_all) {
  // insert key rule mapping
  insert_static_rule(1, "", "rule1");
  insert_static_rule(1, "", "rule2");
  insert_static_rule(2, "", "rule3");

  // insert initial session credit
  auto cfg1 = get_default_config(IMSI1);
  auto cfg2 = get_default_config(IMSI2);
  cfg2.common_context.mutable_teids()->CopyFrom(teids2);
  cfg2.rat_specific_context.mutable_lte_context()->set_bearer_id(BEARER_ID_2);

  CreateSessionResponse response;
  create_credit_update_response(
      IMSI1, SESSION_ID_1, 1, 1024, response.mutable_credits()->Add());
  local_enforcer->init_session(
      session_map, IMSI1, SESSION_ID_1, cfg1, response);
  local_enforcer->update_tunnel_ids(
      session_map,
      create_update_tunnel_ids_request(IMSI1, BEARER_ID_1, teids1));

  CreateSessionResponse response2;
  create_credit_update_response(
      IMSI2, SESSION_ID_2, 2, 2048, response2.mutable_credits()->Add());
  local_enforcer->init_session(
      session_map, IMSI2, SESSION_ID_2, cfg2, response2);
  local_enforcer->update_tunnel_ids(
      session_map,
      create_update_tunnel_ids_request(IMSI2, BEARER_ID_2, teids2));

  assert_charging_credit(
      session_map, IMSI1, SESSION_ID_1, ALLOWED_TOTAL, {{1, 1024}});
  assert_charging_credit(
      session_map, IMSI2, SESSION_ID_2, ALLOWED_TOTAL, {{2, 2048}});

  // receive usages from pipelined
  RuleRecordTable table;
  auto record_list = table.mutable_records();
  create_rule_record(
      IMSI1, teids1.agw_teid(), "rule1", 10, 20, record_list->Add());
  create_rule_record(
      IMSI1, teids1.agw_teid(), "rule2", 5, 15, record_list->Add());
  create_rule_record(
      IMSI2, teids2.agw_teid(), "rule3", 1024, 1024, record_list->Add());
  auto update = SessionStore::get_default_session_update(session_map);
  local_enforcer->aggregate_records(session_map, table, update);

  assert_charging_credit(session_map, IMSI1, SESSION_ID_1, USED_RX, {{1, 15}});
  assert_charging_credit(session_map, IMSI1, SESSION_ID_1, USED_TX, {{1, 35}});
  assert_charging_credit(
      session_map, IMSI2, SESSION_ID_2, USED_RX, {{2, 1024}});
  assert_charging_credit(
      session_map, IMSI2, SESSION_ID_2, USED_TX, {{2, 1024}});

  EXPECT_EQ(update.size(), 2);

  EXPECT_EQ(update[IMSI1][SESSION_ID_1].charging_credit_map.size(), 1);
  // UpdateCriteria does not store REPORTING_RX / REPORTING_TX
  EXPECT_EQ(
      update[IMSI1][SESSION_ID_1].charging_credit_map[1].bucket_deltas[USED_RX],
      15);
  EXPECT_EQ(
      update[IMSI1][SESSION_ID_1].charging_credit_map[1].bucket_deltas[USED_TX],
      35);

  EXPECT_EQ(update[IMSI2][SESSION_ID_2].charging_credit_map.size(), 1);
  // UpdateCriteria does not store REPORTING_RX / REPORTING_TX
  EXPECT_EQ(
      update[IMSI2][SESSION_ID_2].charging_credit_map[2].bucket_deltas[USED_RX],
      1024);
  EXPECT_EQ(
      update[IMSI2][SESSION_ID_2].charging_credit_map[2].bucket_deltas[USED_TX],
      1024);

  // Collect updates for reporting
  std::vector<std::unique_ptr<ServiceAction>> actions;
  auto session_update =
      local_enforcer->collect_updates(session_map, actions, update);
  local_enforcer->execute_actions(session_map, actions, update);
  EXPECT_EQ(session_update.updates_size(), 1);

  assert_charging_credit(
      session_map, IMSI2, SESSION_ID_2, REPORTING_RX, {{2, 1024}});
  assert_charging_credit(
      session_map, IMSI2, SESSION_ID_2, REPORTING_TX, {{2, 1024}});
  assert_charging_credit(
      session_map, IMSI2, SESSION_ID_2, REPORTED_RX, {{2, 0}});
  assert_charging_credit(
      session_map, IMSI2, SESSION_ID_2, REPORTED_TX, {{2, 0}});

  // Add updated credit from cloud
  UpdateSessionResponse update_response;
  auto updates = update_response.mutable_responses();
  create_credit_update_response(IMSI2, SESSION_ID_2, 2, 4096, updates->Add());
  local_enforcer->update_session_credits_and_rules(
      session_map, update_response, update);

  assert_charging_credit(
      session_map, IMSI2, SESSION_ID_2, ALLOWED_TOTAL, {{2, 6144}});
  assert_charging_credit(
      session_map, IMSI2, SESSION_ID_2, REPORTING_RX, {{2, 0}});
  assert_charging_credit(
      session_map, IMSI2, SESSION_ID_2, REPORTING_TX, {{2, 0}});
  assert_charging_credit(
      session_map, IMSI2, SESSION_ID_2, REPORTED_RX, {{2, 1024}});
  assert_charging_credit(
      session_map, IMSI2, SESSION_ID_2, REPORTED_TX, {{2, 1024}});

  EXPECT_EQ(
      update[IMSI2][SESSION_ID_2]
          .charging_credit_map[2]
          .bucket_deltas[REPORTED_TX],
      1024);
  EXPECT_EQ(
      update[IMSI2][SESSION_ID_2]
          .charging_credit_map[2]
          .bucket_deltas[REPORTED_RX],
      1024);

  // Terminate IMSI1
  local_enforcer->handle_termination_from_access(
      session_map, IMSI1, APN1, update);

  EXPECT_CALL(
      *reporter,
      report_terminate_session(CheckTerminateRequestCount(IMSI1, 0, 1), _))
      .Times(1);
  run_evb();

  RuleRecordTable empty_table;
  local_enforcer->aggregate_records(session_map, empty_table, update);
}

TEST_F(LocalEnforcerTest, test_credit_init_with_transient_error_redirect) {
  // insert key rule mapping
  insert_static_rule(1, "", "rule1");

  // insert initial session credit
  CreateSessionResponse response;
  auto credits = response.mutable_credits();
  response.mutable_static_rules()->Add()->set_rule_id("rule1");
  create_credit_update_response_with_error(
      IMSI1, SESSION_ID_1, 1, false, DIAMETER_CREDIT_LIMIT_REACHED,
      ChargingCredit_FinalAction_REDIRECT, "12.7.7.4", "", credits->Add());

  EXPECT_CALL(
      *pipelined_client, deactivate_flows_for_rules(
                             testing::_, testing::_, testing::_, testing::_,
                             CheckRuleCount(1), testing::_))
      .Times(1);
  local_enforcer->init_session(
      session_map, IMSI1, SESSION_ID_1, default_cfg_1, response);
  local_enforcer->update_tunnel_ids(
      session_map,
      create_update_tunnel_ids_request(IMSI1, BEARER_ID_1, teids1));

  // Write + Read in/from SessionStore so all paramters during init are saved
  // to the session
  bool write_success =
      session_store->create_sessions(IMSI1, std::move(session_map[IMSI1]));
  EXPECT_TRUE(write_success);

  evb->loopOnce();

  session_map = session_store->read_sessions({IMSI1});
  auto update = SessionStore::get_default_session_update(session_map);

  // receive usages from pipeline and check we still collect them but don't
  // terminate session (this step helps to check the credit is in the suspended
  // state
  RuleRecordTable table;
  auto record_list = table.mutable_records();
  create_rule_record(
      IMSI1, teids1.agw_teid(), "rule1", 10, 20, record_list->Add());
  local_enforcer->aggregate_records(session_map, table, update);

  assert_charging_credit(session_map, IMSI1, SESSION_ID_1, USED_RX, {{1, 10}});
  assert_charging_credit(session_map, IMSI1, SESSION_ID_1, USED_TX, {{1, 20}});
  EXPECT_EQ(update.size(), 1);
  EXPECT_EQ(update[IMSI1][SESSION_ID_1].charging_credit_map.size(), 1);
  EXPECT_TRUE(update[IMSI1][SESSION_ID_1]
                  .charging_credit_map.find(1)
                  ->second.suspended);
  EXPECT_EQ(
      update[IMSI1][SESSION_ID_1]
          .charging_credit_map.find(1)
          ->second.service_state,
      SERVICE_NEEDS_SUSPENSION);

  // Collect updates for reporting and check Redirect action needs to be done
  std::vector<std::unique_ptr<ServiceAction>> actions;
  auto session_update =
      local_enforcer->collect_updates(session_map, actions, update);
  EXPECT_EQ(actions.size(), 1);

  // execute actions (from collect_updates)
  EXPECT_CALL(
      *pipelined_client,
      add_gy_final_action_flow(
          IMSI1, default_cfg_1.common_context.ue_ipv4(),
          default_cfg_1.common_context.ue_ipv6(), CheckTeids(teids1),
          default_cfg_1.common_context.msisdn(), CheckRuleCount(1)))
      .Times(1);
  local_enforcer->execute_actions(session_map, actions, update);
  EXPECT_EQ(session_update.updates_size(), 0);
  EXPECT_EQ(
      update[IMSI1][SESSION_ID_1]
          .charging_credit_map.find(1)
          ->second.service_state,
      SERVICE_REDIRECTED);
  EXPECT_TRUE(update[IMSI1][SESSION_ID_1]
                  .charging_credit_map.find(1)
                  ->second.suspended);
}

TEST_F(LocalEnforcerTest, test_update_with_transient_error) {
  // insert key rule mapping
  insert_static_rule(1, "", "rule1");
  insert_static_rule(1, "", "rule2");
  insert_static_rule(2, "", "rule3");

  // insert initial session credit + rules
  CreateSessionResponse response;
  create_credit_update_response(
      IMSI1, SESSION_ID_1, 1, 1024, response.mutable_credits()->Add());
  create_credit_update_response(
      IMSI1, SESSION_ID_1, 2, 1024, response.mutable_credits()->Add());
  response.mutable_static_rules()->Add()->set_rule_id("rule1");
  response.mutable_static_rules()->Add()->set_rule_id("rule2");
  response.mutable_static_rules()->Add()->set_rule_id("rule3");

  local_enforcer->init_session(
      session_map, IMSI1, SESSION_ID_1, test_cfg_, response);
  local_enforcer->update_tunnel_ids(
      session_map,
      create_update_tunnel_ids_request(IMSI1, BEARER_ID_1, teids1));

  // Add updated credit from cloud
  auto session_uc = SessionStore::get_default_session_update(session_map);
  UpdateSessionResponse update_response;
  auto updates = update_response.mutable_responses();
  create_credit_update_response_with_error(
      IMSI1, SESSION_ID_1, 1, false, DIAMETER_CREDIT_LIMIT_REACHED,
      updates->Add());
  create_credit_update_response(
      IMSI1, SESSION_ID_1, 2, 1024, response.mutable_credits()->Add());
  EXPECT_CALL(
      *pipelined_client, deactivate_flows_for_rules(
                             testing::_, testing::_, testing::_, testing::_,
                             CheckRuleCount(2), testing::_))
      .Times(1);

  local_enforcer->update_session_credits_and_rules(
      session_map, update_response, session_uc);
  EXPECT_FALSE(session_uc[IMSI1][SESSION_ID_1].updated_fsm_state);
  EXPECT_TRUE(session_uc[IMSI1][SESSION_ID_1].charging_credit_map[1].suspended);
  EXPECT_EQ(
      session_uc[IMSI1][SESSION_ID_1].charging_credit_map[1].service_state,
      SESSION_ACTIVE);
  EXPECT_FALSE(
      session_uc[IMSI1][SESSION_ID_1].charging_credit_map[2].suspended);
  EXPECT_EQ(
      session_uc[IMSI1][SESSION_ID_1].charging_credit_map[2].service_state,
      SESSION_ACTIVE);
}

// gy_rar
TEST_F(LocalEnforcerTest, test_reauth_with_redirected_suspended_credit) {
  // insert key rule mapping
  insert_static_rule(1, "", "rule1");

  // 1- INITIAL SET UP TO CREATE A REDIRECTED due to SUSPENDED CREDIT
  // insert initial suspended and redirected credit
  CreateSessionResponse response;
  response.mutable_static_rules()->Add()->set_rule_id("rule1");
  auto credits = response.mutable_credits();
  test_cfg_.common_context.mutable_sid()->set_id(IMSI1);
  create_credit_update_response_with_error(
      IMSI1, SESSION_ID_1, 1, false, DIAMETER_CREDIT_LIMIT_REACHED,
      ChargingCredit_FinalAction_REDIRECT, "12.7.7.4", "", credits->Add());
  local_enforcer->init_session(
      session_map, IMSI1, SESSION_ID_1, test_cfg_, response);
  local_enforcer->update_tunnel_ids(
      session_map,
      create_update_tunnel_ids_request(IMSI1, BEARER_ID_1, teids1));

  // execute redirection
  auto update = SessionStore::get_default_session_update(session_map);
  std::vector<std::unique_ptr<ServiceAction>> actions;
  auto session_update =
      local_enforcer->collect_updates(session_map, actions, update);
  local_enforcer->execute_actions(session_map, actions, update);
  EXPECT_EQ(actions.size(), 1);
  EXPECT_EQ(session_update.updates_size(), 0);
  EXPECT_EQ(
      update[IMSI1][SESSION_ID_1].charging_credit_map[1].service_state,
      SERVICE_REDIRECTED);
  EXPECT_TRUE(update[IMSI1][SESSION_ID_1].charging_credit_map[1].suspended);

  // 2- SETUP REAUTH ON A SUSPENDED CREDIT
  // reset update and actions
  update = SessionStore::get_default_session_update(session_map);
  actions.clear();

  ChargingReAuthRequest reauth;
  reauth.set_sid(IMSI1);
  reauth.set_session_id(SESSION_ID_1);
  reauth.set_type(ChargingReAuthRequest::ENTIRE_SESSION);
  auto result =
      local_enforcer->init_charging_reauth(session_map, reauth, update);
  EXPECT_EQ(result, ReAuthResult::UPDATE_INITIATED);

  // check suspended credit is requested
  auto update_req =
      local_enforcer->collect_updates(session_map, actions, update);
  local_enforcer->execute_actions(session_map, actions, update);
  EXPECT_EQ(update_req.updates_size(), 1);
  EXPECT_EQ(update_req.updates(0).common_context().sid().id(), IMSI1);
  EXPECT_EQ(update_req.updates(0).usage().type(), CreditUsage::REAUTH_REQUIRED);

  // 3- SUSPENSION IS CLEARED
  // receive credit from OCS as a reauth answer
  UpdateSessionResponse update_response;
  auto updates = update_response.mutable_responses();
  create_credit_update_response(IMSI1, SESSION_ID_1, 1, 4096, updates->Add());
  std::vector<std::string> pending_activation = {"rule1"};
  EXPECT_CALL(
      *pipelined_client, activate_flows_for_rules(
                             testing::_, testing::_, testing::_, testing::_,
                             test_cfg_.common_context.msisdn(), testing::_,
                             CheckRuleNames(pending_activation), testing::_))
      .Times(1);
  local_enforcer->update_session_credits_and_rules(
      session_map, update_response, update);
  EXPECT_FALSE(update[IMSI1][SESSION_ID_1].charging_credit_map[1].suspended);
}

TEST_F(LocalEnforcerTest, test_re_auth) {
  insert_static_rule(1, "", "rule1");
  CreateSessionResponse response;
  local_enforcer->init_session(
      session_map, IMSI1, SESSION_ID_1, get_default_config(IMSI1), response);
  local_enforcer->update_tunnel_ids(
      session_map,
      create_update_tunnel_ids_request(IMSI1, BEARER_ID_1, teids1));

  ChargingReAuthRequest reauth;
  reauth.set_sid(IMSI1);
  reauth.set_session_id(SESSION_ID_1);
  reauth.set_charging_key(1);
  reauth.set_type(ChargingReAuthRequest::SINGLE_SERVICE);
  auto update = SessionStore::get_default_session_update(session_map);
  auto result =
      local_enforcer->init_charging_reauth(session_map, reauth, update);
  EXPECT_EQ(result, ReAuthResult::UPDATE_INITIATED);

  std::vector<std::unique_ptr<ServiceAction>> actions;
  auto update_req =
      local_enforcer->collect_updates(session_map, actions, update);
  local_enforcer->execute_actions(session_map, actions, update);
  EXPECT_EQ(update_req.updates_size(), 1);
  EXPECT_EQ(update_req.updates(0).common_context().sid().id(), IMSI1);
  EXPECT_EQ(update_req.updates(0).usage().type(), CreditUsage::REAUTH_REQUIRED);

  // Give credit after re-auth
  UpdateSessionResponse update_response;
  auto updates = update_response.mutable_responses();
  create_credit_update_response(IMSI1, SESSION_ID_1, 1, 4096, updates->Add());
  local_enforcer->update_session_credits_and_rules(
      session_map, update_response, update);

  // when next update is collected, this should trigger an action to activate
  // the flow in pipelined
  EXPECT_CALL(
      *pipelined_client, activate_flows_for_rules(
                             testing::_, testing::_, testing::_, testing::_,
                             test_cfg_.common_context.msisdn(), testing::_,
                             testing::_, testing::_))
      .Times(1);
  actions.clear();
  local_enforcer->collect_updates(session_map, actions, update);
  local_enforcer->execute_actions(session_map, actions, update);
}

TEST_F(LocalEnforcerTest, test_dynamic_rules) {
  CreateSessionResponse response;
  create_credit_update_response(
      IMSI1, SESSION_ID_1, 1, 1024, response.mutable_credits()->Add());
  auto dynamic_rule = response.mutable_dynamic_rules()->Add();
  auto policy_rule  = dynamic_rule->mutable_policy_rule();
  policy_rule->set_id("rule1");
  policy_rule->set_rating_group(1);
  policy_rule->set_tracking_type(PolicyRule::ONLY_OCS);
  local_enforcer->init_session(
      session_map, IMSI1, SESSION_ID_1, test_cfg_, response);
  local_enforcer->update_tunnel_ids(
      session_map,
      create_update_tunnel_ids_request(IMSI1, BEARER_ID_1, teids1));

  insert_static_rule(1, "", "rule2");
  RuleRecordTable table;
  auto record_list = table.mutable_records();
  create_rule_record(
      IMSI1, teids1.agw_teid(), "rule1", 16, 32, record_list->Add());
  create_rule_record(
      IMSI1, teids1.agw_teid(), "rule2", 8, 8, record_list->Add());

  auto update = SessionStore::get_default_session_update(session_map);
  local_enforcer->aggregate_records(session_map, table, update);

  assert_charging_credit(session_map, IMSI1, SESSION_ID_1, USED_RX, {{1, 24}});
  assert_charging_credit(session_map, IMSI1, SESSION_ID_1, USED_TX, {{1, 40}});
  assert_charging_credit(
      session_map, IMSI1, SESSION_ID_1, ALLOWED_TOTAL, {{1, 1024}});

  EXPECT_EQ(
      update[IMSI1][SESSION_ID_1].charging_credit_map[1].bucket_deltas[USED_RX],
      24);
  EXPECT_EQ(
      update[IMSI1][SESSION_ID_1].charging_credit_map[1].bucket_deltas[USED_TX],
      40);
}

TEST_F(LocalEnforcerTest, test_dynamic_rule_actions) {
  CreateSessionResponse response;
  // with final action = terminate
  create_credit_update_response(
      IMSI1, SESSION_ID_1, 1, 1024, true, response.mutable_credits()->Add());
  auto dynamic_rule = response.mutable_dynamic_rules()->Add();
  auto policy_rule  = dynamic_rule->mutable_policy_rule();
  policy_rule->set_id("rule2");
  policy_rule->set_rating_group(1);
  policy_rule->set_tracking_type(PolicyRule::ONLY_OCS);
  auto static_rule = response.mutable_static_rules()->Add();
  static_rule->set_rule_id("rule1");
  static_rule = response.mutable_static_rules()->Add();
  static_rule->set_rule_id("rule3");
  insert_static_rule(1, "", "rule1");
  insert_static_rule(1, "", "rule3");

  // The activation for the static rules (rule1,rule3) and dynamic rule (rule2)
  EXPECT_CALL(
      *pipelined_client, activate_flows_for_rules(
                             testing::_, testing::_, testing::_, testing::_,
                             default_cfg_1.common_context.msisdn(), testing::_,
                             CheckRuleCount(3), testing::_))
      .Times(1);

  local_enforcer->init_session(
      session_map, IMSI1, SESSION_ID_1, default_cfg_1, response);
  local_enforcer->update_tunnel_ids(
      session_map,
      create_update_tunnel_ids_request(IMSI1, BEARER_ID_1, teids1));

  RuleRecordTable table;
  auto record_list = table.mutable_records();
  create_rule_record(
      IMSI1, teids1.agw_teid(), "rule1", 1024, 2048, record_list->Add());
  create_rule_record(
      IMSI1, teids1.agw_teid(), "rule2", 1024, 2048, record_list->Add());
  auto update = SessionStore::get_default_session_update(session_map);
  local_enforcer->aggregate_records(session_map, table, update);

  // the request should has no rules so PipelineD deletes all rules
  EXPECT_CALL(
      *pipelined_client,
      deactivate_flows_for_rules_for_termination(
          testing::_, testing::_, testing::_, testing::_, testing::_))
      .Times(1);
  std::vector<std::unique_ptr<ServiceAction>> actions;
  auto usage_updates =
      local_enforcer->collect_updates(session_map, actions, update);
  local_enforcer->execute_actions(session_map, actions, update);
}

TEST_F(LocalEnforcerTest, test_installing_rules_with_activation_time) {
  CreateSessionResponse response;
  create_credit_update_response(
      IMSI1, SESSION_ID_1, 1, 1024, true, response.mutable_credits()->Add());
  auto now = time(nullptr);

  // add a dynamic rule without activation time
  auto dynamic_rule = response.mutable_dynamic_rules()->Add();
  auto policy_rule  = dynamic_rule->mutable_policy_rule();
  policy_rule->set_id("rule1");
  policy_rule->set_rating_group(1);
  policy_rule->set_tracking_type(PolicyRule::ONLY_OCS);

  // add a dynamic rule with activation time in the future
  dynamic_rule = response.mutable_dynamic_rules()->Add();
  policy_rule  = dynamic_rule->mutable_policy_rule();
  dynamic_rule->mutable_activation_time()->set_seconds(now + SECONDS_A_DAY);
  policy_rule->set_id("rule2");
  policy_rule->set_rating_group(1);
  policy_rule->set_tracking_type(PolicyRule::ONLY_OCS);

  // add a dynamic rule with activation time in the past
  dynamic_rule = response.mutable_dynamic_rules()->Add();
  policy_rule  = dynamic_rule->mutable_policy_rule();
  dynamic_rule->mutable_activation_time()->set_seconds(now - SECONDS_A_DAY);
  policy_rule->set_id("rule3");
  policy_rule->set_rating_group(1);
  policy_rule->set_tracking_type(PolicyRule::ONLY_OCS);

  // add a static rule without activation time
  insert_static_rule(1, "", "rule4");
  auto static_rule = response.mutable_static_rules()->Add();
  static_rule->set_rule_id("rule4");

  // add a static rule with activation time in the future
  insert_static_rule(1, "", "rule5");
  static_rule = response.mutable_static_rules()->Add();
  static_rule->mutable_activation_time()->set_seconds(now + SECONDS_A_DAY);
  static_rule->set_rule_id("rule5");

  // add a static rule with activation time in the past
  insert_static_rule(1, "", "rule6");
  static_rule = response.mutable_static_rules()->Add();
  static_rule->mutable_activation_time()->set_seconds(now - SECONDS_A_DAY);
  static_rule->set_rule_id("rule6");

  // expect calling activate_flows_for_rules for activating rules instantly
  // dynamic rules: rule1, rule3
  // static rules: rule4, rule6
  test_cfg_.common_context.set_ue_ipv4(IP3);
  test_cfg_.common_context.mutable_sid()->set_id(IMSI1);
  std::string ip_addr   = test_cfg_.common_context.ue_ipv4();
  std::string ipv6_addr = test_cfg_.common_context.ue_ipv6();
  Teids teids           = teids1;
  EXPECT_CALL(
      *pipelined_client, activate_flows_for_rules(
                             IMSI1, ip_addr, ipv6_addr, CheckTeids(teids),
                             test_cfg_.common_context.msisdn(), testing::_,
                             CheckRuleCount(4), testing::_))
      .Times(1);

  // We do not expect rule5 and rule2 to be activated since they are scheduled a
  // day away

  local_enforcer->init_session(
      session_map, IMSI1, SESSION_ID_1, test_cfg_, response);
  local_enforcer->update_tunnel_ids(
      session_map,
      create_update_tunnel_ids_request(IMSI1, BEARER_ID_1, teids1));
  bool success =
      session_store->create_sessions(IMSI1, std::move(session_map[IMSI1]));
  EXPECT_TRUE(success);
}

TEST_F(LocalEnforcerTest, test_usage_monitors) {
  // insert key rule mapping
  insert_static_rule(1, "1", "both_rule");
  insert_static_rule(2, "", "ocs_rule");
  insert_static_rule(0, "3", "pcrf_only");
  insert_static_rule(0, "1", "pcrf_split");  // same mkey as both_rule
  // session level rule "4"

  // insert initial session credit
  CreateSessionResponse response;
  create_credit_update_response(
      IMSI1, SESSION_ID_1, 1, 1024, response.mutable_credits()->Add());
  create_credit_update_response(
      IMSI1, SESSION_ID_1, 2, 1024, response.mutable_credits()->Add());
  create_monitor_update_response(
      IMSI1, SESSION_ID_1, "1", MonitoringLevel::PCC_RULE_LEVEL, 1024,
      response.mutable_usage_monitors()->Add());
  create_monitor_update_response(
      IMSI1, SESSION_ID_1, "3", MonitoringLevel::PCC_RULE_LEVEL, 2048,
      response.mutable_usage_monitors()->Add());
  create_monitor_update_response(
      IMSI1, SESSION_ID_1, "4", MonitoringLevel::SESSION_LEVEL, 2128,
      response.mutable_usage_monitors()->Add());
  response.mutable_static_rules()->Add()->set_rule_id("both_rule");
  response.mutable_static_rules()->Add()->set_rule_id("ocs_rule");
  response.mutable_static_rules()->Add()->set_rule_id("pcrf_only");
  response.mutable_static_rules()->Add()->set_rule_id("pcrf_split");
  test_cfg_.common_context.mutable_sid()->set_id(IMSI1);

  local_enforcer->init_session(
      session_map, IMSI1, SESSION_ID_1, test_cfg_, response);
  local_enforcer->update_tunnel_ids(
      session_map,
      create_update_tunnel_ids_request(IMSI1, BEARER_ID_1, teids1));
  assert_charging_credit(
      session_map, IMSI1, SESSION_ID_1, ALLOWED_TOTAL, {{1, 1024}, {2, 1024}});
  assert_monitor_credit(
      session_map, IMSI1, SESSION_ID_1, ALLOWED_TOTAL,
      {{"1", 1024}, {"3", 2048}, {"4", 2128}});

  // receive usages from pipelined
  RuleRecordTable table;
  auto record_list = table.mutable_records();
  create_rule_record(
      IMSI1, teids1.agw_teid(), "both_rule", 10, 20, record_list->Add());
  create_rule_record(
      IMSI1, teids1.agw_teid(), "ocs_rule", 5, 15, record_list->Add());
  create_rule_record(
      IMSI1, teids1.agw_teid(), "pcrf_only", 1024, 1024, record_list->Add());
  create_rule_record(
      IMSI1, teids1.agw_teid(), "pcrf_split", 10, 20, record_list->Add());
  auto update = SessionStore::get_default_session_update(session_map);
  local_enforcer->aggregate_records(session_map, table, update);

  assert_charging_credit(
      session_map, IMSI1, SESSION_ID_1, USED_RX, {{1, 10}, {2, 5}});
  assert_charging_credit(
      session_map, IMSI1, SESSION_ID_1, USED_TX, {{1, 20}, {2, 15}});
  assert_monitor_credit(
      session_map, IMSI1, SESSION_ID_1, USED_RX,
      {{"1", 20}, {"3", 1024}, {"4", 1049}});
  assert_monitor_credit(
      session_map, IMSI1, SESSION_ID_1, USED_TX,
      {{"1", 40}, {"3", 1024}, {"4", 1079}});

  // Collect updates, should only have mkeys 3 and 4
  std::vector<std::unique_ptr<ServiceAction>> actions;
  auto session_update =
      local_enforcer->collect_updates(session_map, actions, update);
  local_enforcer->execute_actions(session_map, actions, update);
  EXPECT_EQ(session_update.usage_monitors_size(), 2);
  for (const auto& monitor : session_update.usage_monitors()) {
    EXPECT_EQ(monitor.sid(), IMSI1);
    if (monitor.update().monitoring_key() == "3") {
      EXPECT_EQ(monitor.update().level(), MonitoringLevel::PCC_RULE_LEVEL);
      EXPECT_EQ(monitor.update().bytes_rx(), 1024);
      EXPECT_EQ(monitor.update().bytes_tx(), 1024);
    } else if (monitor.update().monitoring_key() == "4") {
      EXPECT_EQ(monitor.update().level(), MonitoringLevel::SESSION_LEVEL);
      EXPECT_EQ(monitor.update().bytes_rx(), 1049);
      EXPECT_EQ(monitor.update().bytes_tx(), 1079);
    } else {
      EXPECT_TRUE(false);
    }
  }

  assert_charging_credit(
      session_map, IMSI1, SESSION_ID_1, REPORTING_RX, {{1, 0}, {2, 0}});
  assert_charging_credit(
      session_map, IMSI1, SESSION_ID_1, REPORTING_TX, {{1, 0}, {2, 0}});
  assert_monitor_credit(
      session_map, IMSI1, SESSION_ID_1, REPORTING_RX,
      {{"1", 0}, {"3", 1024}, {"4", 1049}});
  assert_monitor_credit(
      session_map, IMSI1, SESSION_ID_1, REPORTING_TX,
      {{"1", 0}, {"3", 1024}, {"4", 1079}});

  UpdateSessionResponse update_response;
  auto monitor_updates = update_response.mutable_usage_monitor_responses();
  create_monitor_update_response(
      IMSI1, SESSION_ID_1, "3", MonitoringLevel::PCC_RULE_LEVEL, 2048,
      monitor_updates->Add());
  create_monitor_update_response(
      IMSI1, SESSION_ID_1, "4", MonitoringLevel::SESSION_LEVEL, 2048,
      monitor_updates->Add());
  local_enforcer->update_session_credits_and_rules(
      session_map, update_response, update);
  assert_monitor_credit(
      session_map, IMSI1, SESSION_ID_1, REPORTING_RX, {{"3", 0}, {"4", 0}});
  assert_monitor_credit(
      session_map, IMSI1, SESSION_ID_1, REPORTING_TX, {{"3", 0}, {"4", 0}});
  assert_monitor_credit(
      session_map, IMSI1, SESSION_ID_1, REPORTED_RX,
      {{"3", 1024}, {"4", 1049}});
  assert_monitor_credit(
      session_map, IMSI1, SESSION_ID_1, REPORTED_TX,
      {{"3", 1024}, {"4", 1079}});
  assert_monitor_credit(
      session_map, IMSI1, SESSION_ID_1, ALLOWED_TOTAL,
      {{"3", 4096}, {"4", 4176}});

  // Test rule removal in usage monitor response for CCA-Update
  update_response.Clear();
  monitor_updates = update_response.mutable_usage_monitor_responses();
  auto monitor_updates_response = monitor_updates->Add();
  create_monitor_update_response(
      IMSI1, SESSION_ID_1, "3", MonitoringLevel::PCC_RULE_LEVEL, 0,
      monitor_updates_response);
  monitor_updates_response->add_rules_to_remove("pcrf_only");

  EXPECT_CALL(
      *pipelined_client, deactivate_flows_for_rules(
                             IMSI1, testing::_, testing::_, testing::_,
                             CheckPolicyID("pcrf_only"), RequestOriginType::GX))
      .Times(1);
  local_enforcer->update_session_credits_and_rules(
      session_map, update_response, update);

  // Test rule installation in usage monitor response for CCA-Update
  update_response.Clear();
  monitor_updates          = update_response.mutable_usage_monitor_responses();
  monitor_updates_response = monitor_updates->Add();
  create_monitor_update_response(
      IMSI1, SESSION_ID_1, "3", MonitoringLevel::PCC_RULE_LEVEL, 0,
      monitor_updates_response);

  StaticRuleInstall static_rule_install;
  static_rule_install.set_rule_id("pcrf_only");
  auto res_rules_to_install =
      monitor_updates_response->add_static_rules_to_install();
  res_rules_to_install->CopyFrom(static_rule_install);

  EXPECT_CALL(
      *pipelined_client, activate_flows_for_rules(
                             IMSI1, testing::_, testing::_, testing::_,
                             test_cfg_.common_context.msisdn(), testing::_,
                             CheckRuleCount(1), testing::_))
      .Times(1);
  local_enforcer->update_session_credits_and_rules(
      session_map, update_response, update);
}

// Test an insertion of a usage monitor, both session_level and rule level,
// and then a deletion. Additionally, test
// that a rule update from PipelineD for a deleted usage monitor should NOT
// trigger an update request,
// Note the actual deletion ofo a monitor will not happen until last update
// is sent. So we need to store and read the session to trigger that deletion
TEST_F(LocalEnforcerTest, test_usage_monitor_disable) {
  // insert key rule mapping
  insert_static_rule(0, "1", "pcrf_only_active");
  insert_static_rule(0, "3", "pcrf_only_to_be_disabled");

  // Monitor credit addition #1
  // insert initial session credit
  CreateSessionResponse response_1;
  create_monitor_update_response(
      IMSI1, SESSION_ID_1, "1", MonitoringLevel::PCC_RULE_LEVEL, 1024,
      response_1.mutable_usage_monitors()->Add());
  create_monitor_update_response(
      IMSI1, SESSION_ID_1, "2", MonitoringLevel::SESSION_LEVEL, 1024,
      response_1.mutable_usage_monitors()->Add());
  create_monitor_update_response(
      IMSI1, SESSION_ID_1, "3", MonitoringLevel::PCC_RULE_LEVEL, 1024,
      response_1.mutable_usage_monitors()->Add());
  local_enforcer->init_session(
      session_map, IMSI1, SESSION_ID_1, test_cfg_, response_1);
  local_enforcer->update_tunnel_ids(
      session_map,
      create_update_tunnel_ids_request(IMSI1, BEARER_ID_1, teids1));
  assert_monitor_credit(
      session_map, IMSI1, SESSION_ID_1, ALLOWED_TOTAL,
      {{"1", 1024}, {"2", 1024}, {"3", 1024}});

  // IMPORTANT: save the updates into store and reload
  bool success =
      session_store->create_sessions(IMSI1, std::move(session_map[IMSI1]));
  EXPECT_TRUE(success);
  session_map   = session_store->read_sessions(SessionRead{IMSI1});
  auto update_1 = SessionStore::get_default_session_update(session_map);

  // Monitor credit usage #1
  // Use the quota to exhaust monitor 2 and 3 and assert that we send usage
  // reports for all monitors receive usages from pipelined
  RuleRecordTable table_1;
  auto record_list_1 = table_1.mutable_records();
  create_rule_record(
      IMSI1, teids1.agw_teid(), "pcrf_only_active", 2000, 0,
      record_list_1->Add());
  create_rule_record(
      IMSI1, teids1.agw_teid(), "pcrf_only_to_be_disabled", 2000, 0,
      record_list_1->Add());
  local_enforcer->aggregate_records(session_map, table_1, update_1);

  // Collect updates, should have updates since all monitors got 80% exhausted
  std::vector<std::unique_ptr<ServiceAction>> actions_1;
  auto update_request_1 =
      local_enforcer->collect_updates(session_map, actions_1, update_1);
  EXPECT_EQ(update_request_1.updates_size(), 0);
  EXPECT_EQ(update_request_1.usage_monitors_size(), 3);

  // IMPORTANT: save the updates into store and reload to trigger
  // apply_monitor_updates
  success = session_store->update_sessions(update_1);
  EXPECT_TRUE(success);
  session_map = session_store->read_sessions(SessionRead{IMSI1});
  update_1    = SessionStore::get_default_session_update(session_map);
  EXPECT_EQ(session_map.size(), 1);

  // Monitor credit addition #2
  // Receive an update with zero grant for mkey=3 & mkey=2, but with
  // credit for mkey=1. That means 3 and 2 will have to stop reporting when
  // their quotas are exhausted
  UpdateSessionResponse update_response_2;
  auto monitors_2 = update_response_2.mutable_usage_monitor_responses();
  create_monitor_update_response(
      IMSI1, SESSION_ID_1, "1", MonitoringLevel::PCC_RULE_LEVEL, 1024,
      monitors_2->Add());
  create_monitor_update_response(
      IMSI1, SESSION_ID_1, "2", MonitoringLevel::SESSION_LEVEL, 0,
      monitors_2->Add());
  create_monitor_update_response(
      IMSI1, SESSION_ID_1, "3", MonitoringLevel::PCC_RULE_LEVEL, 0,
      monitors_2->Add());
  // Apply the updates
  auto update_2 = SessionStore::get_default_session_update(session_map);
  local_enforcer->update_session_credits_and_rules(
      session_map, update_response_2, update_2);

  auto monitor_updates_2 = update_2[IMSI1][SESSION_ID_1].monitor_credit_map;
  EXPECT_FALSE(monitor_updates_2["1"].deleted);
  EXPECT_FALSE(monitor_updates_2["2"].deleted);
  EXPECT_FALSE(monitor_updates_2["3"].deleted);
  // Check updates for disabling session level monitoring key
  EXPECT_TRUE(update_2[IMSI1][SESSION_ID_1].is_session_level_key_updated);
  EXPECT_EQ(update_2[IMSI1][SESSION_ID_1].updated_session_level_key, "2");
  // note that we have gone over allowed, so we will reset ALLOWED_TOTAL
  // 3024 because used 2000 before (out of the initial 1024) plues we added 1024
  // 4000 because session level has used 2000 + 2000 from mkey 1 and 3
  // 2000 because mkey 3 have used 2000 but we haven't added any extra
  assert_monitor_credit(
      session_map, IMSI1, SESSION_ID_1, ALLOWED_TOTAL,
      {{"1", 3024}, {"2", 4000}, {"3", 2000}});

  // Monitor credit usage #2
  // Generate more traffic to see monitors 2 and 3 are not triggering update
  RuleRecordTable table_2;
  auto record_list2 = table_2.mutable_records();
  create_rule_record(
      IMSI1, teids1.agw_teid(), "pcrf_only_active", 2000, 0,
      record_list2->Add());
  create_rule_record(
      IMSI1, teids1.agw_teid(), "pcrf_only_to_be_disabled", 2000, 0,
      record_list2->Add());

  update_2 = SessionStore::get_default_session_update(session_map);
  local_enforcer->aggregate_records(session_map, table_2, update_2);
  monitor_updates_2 = update_2[IMSI1][SESSION_ID_1].monitor_credit_map;
  EXPECT_FALSE(monitor_updates_2["1"].report_last_credit);
  EXPECT_TRUE(monitor_updates_2["2"].report_last_credit);
  EXPECT_TRUE(monitor_updates_2["3"].report_last_credit);
  // check session level key will be removed
  EXPECT_EQ(update_2[IMSI1][SESSION_ID_1].updated_session_level_key, "");

  // Collect updates, should receive three updatesResponses (two of them of the
  // last updates sent for  2 and 3.
  // Also it should mark 2 and 3 for deletion
  std::vector<std::unique_ptr<ServiceAction>> actions_2;
  auto update_request_2 =
      local_enforcer->collect_updates(session_map, actions_2, update_2);
  EXPECT_EQ(update_request_2.updates_size(), 0);
  EXPECT_EQ(update_request_2.usage_monitors_size(), 3);
  monitor_updates_2 = update_2[IMSI1][SESSION_ID_1].monitor_credit_map;
  EXPECT_FALSE(monitor_updates_2["1"].deleted);
  EXPECT_TRUE(monitor_updates_2["2"].deleted);
  EXPECT_TRUE(monitor_updates_2["3"].deleted);

  // Check deletion of monitors
  // Actual deletion will not happen until we iterate one more timea
  // local_enforcer->up
  // IMPORTANT: save the updates into store and reload to trigger
  success = session_store->update_sessions(update_2);
  EXPECT_TRUE(success);
  session_map = session_store->read_sessions(SessionRead{IMSI1});
  EXPECT_EQ(session_map.size(), 1);
  EXPECT_NE(session_map[IMSI1][0]->get_monitor("1", ALLOWED_TOTAL), 0);
  EXPECT_EQ(session_map[IMSI1][0]->get_monitor("2", ALLOWED_TOTAL), 0);
  EXPECT_EQ(session_map[IMSI1][0]->get_monitor("3", ALLOWED_TOTAL), 0);
}

TEST_F(LocalEnforcerTest, test_rar_create_dedicated_bearer) {
  QosInformationRequest test_qos_info;
  test_qos_info.set_qos_class_id(0);

  SessionConfig test_volte_cfg;
  test_volte_cfg.common_context =
      build_common_context("", IP1, IPv6_1, teids1, "", APN1, TGPP_LTE);
  const auto& lte_context =
      build_lte_context("", "", "", "", "", 1, &test_qos_info);
  test_volte_cfg.rat_specific_context.mutable_lte_context()->CopyFrom(
      lte_context);

  CreateSessionResponse response;
  local_enforcer->init_session(
      session_map, IMSI1, SESSION_ID_1, test_volte_cfg, response);
  local_enforcer->update_tunnel_ids(
      session_map,
      create_update_tunnel_ids_request(IMSI1, BEARER_ID_1, teids1));

  PolicyReAuthRequest rar;
  std::vector<std::string> rules_to_remove;
  std::vector<StaticRuleInstall> rules_to_install;
  std::vector<DynamicRuleInstall> dynamic_rules_to_install;
  std::vector<EventTrigger> event_triggers;
  std::vector<UsageMonitoringCredit> usage_monitoring_credits;
  create_policy_reauth_request(
      SESSION_ID_1, IMSI1, rules_to_remove, rules_to_install,
      dynamic_rules_to_install, event_triggers, time(nullptr),
      usage_monitoring_credits, &rar);
  auto rar_qos_info = rar.mutable_qos_info();
  rar_qos_info->set_qci(QCI_1);

  EXPECT_CALL(*spgw_client, create_dedicated_bearer(testing::_))
      .Times(1)
      .WillOnce(testing::Return(true));

  PolicyReAuthAnswer raa;
  auto update = SessionStore::get_default_session_update(session_map);
  local_enforcer->init_policy_reauth(session_map, rar, raa, update);
  EXPECT_EQ(raa.result(), ReAuthResult::UPDATE_INITIATED);
}

// This test covers some edge cases for the dedicated bearer creation scheduling
// for a new session. For session creation, we delay the actual call to create
// bearer by a few seconds. (BEARER_CREATION_DELAY_ON_SESSION_INIT)
// But since we could have a case where the rule state is modified during the
// scheduling and the actual call, we want to ensure that we do not create
// bearers if:
// 1. Policy is no longer installed.
// 2. Policy no longer needs a bearer for some reason. One case of this could be
// that it is already tied to a bearer.
// The success case where the creation takes place is covered by another test
// "test_dedicated_bearer_lifecycle"
TEST_F(LocalEnforcerTest, test_dedicated_bearer_creation_on_session_init) {
  CreateSessionResponse response1, response2;
  const uint32_t default_bearer_id = 5;
  const uint32_t bearer_1          = 6;
  insert_static_rule_with_qos(0, "m1", "rule1", 1);             // QCI=1
  LocalEnforcer::BEARER_CREATION_DELAY_ON_SESSION_INIT = 1000;  // 1 sec

  // test_cfg_ is initialized with QoSInfo field w/ QCI 5
  test_cfg_.common_context.mutable_sid()->set_id(IMSI1);
  test_cfg_.common_context.set_apn("apn1");
  auto lte_context = test_cfg_.rat_specific_context.mutable_lte_context();
  lte_context->mutable_qos_info()->set_qos_class_id(5);
  lte_context->set_bearer_id(default_bearer_id);  // linked_bearer_id

  // CASE 1: policy has a bearer by the time the scheduled function is executed
  response1.mutable_static_rules()->Add()->set_rule_id("rule1");
  // Schedules default bearer install in 1 sec
  local_enforcer->init_session(
      session_map, IMSI1, SESSION_ID_1, test_cfg_, response1);
  local_enforcer->update_tunnel_ids(
      session_map,
      create_update_tunnel_ids_request(IMSI1, BEARER_ID_1, teids1));
  // Write + Read in/from SessionStore
  bool write_success =
      session_store->create_sessions(IMSI1, std::move(session_map[IMSI1]));
  EXPECT_TRUE(write_success);

  // Before we hit the 1 second, simulate a case where another call triggered
  // the policy to be tied a bearer already
  session_map = session_store->read_sessions({IMSI1});
  auto update = SessionStore::get_default_session_update(session_map);
  PolicyBearerBindingRequest bearer_bind_req_success1 =
      create_policy_bearer_bind_req(
          IMSI1, default_bearer_id, "rule1", bearer_1, 1, 2);
  local_enforcer->bind_policy_to_bearer(
      session_map, bearer_bind_req_success1, update);
  // Write + Read in/from SessionStore
  write_success = session_store->update_sessions(update);
  EXPECT_TRUE(write_success);
  // Since the policy is already tied to a bearer, we should not see an
  // additional create bearer request when the scheduled creation happens
  EXPECT_CALL(
      *spgw_client, create_dedicated_bearer(CheckCreateBearerReq(IMSI1, 1)))
      .Times(0);
  evb->loopOnce();

  // CASE 2: policy has been removed
  session_map = session_store->read_all_sessions();
  response2.mutable_static_rules()->Add()->set_rule_id("rule1");

  // Schedules default bearer install in 1 sec
  auto test_cfg_2 = test_cfg_;
  test_cfg_2.common_context.mutable_teids()->CopyFrom(teids2);
  test_cfg_2.rat_specific_context.mutable_lte_context()->set_bearer_id(
      BEARER_ID_2);

  local_enforcer->init_session(
      session_map, IMSI2, SESSION_ID_2, test_cfg_2, response2);
  local_enforcer->update_tunnel_ids(
      session_map,
      create_update_tunnel_ids_request(IMSI2, BEARER_ID_2, teids2));

  // Write + Read in/from SessionStore
  write_success =
      session_store->create_sessions(IMSI2, std::move(session_map[IMSI2]));
  EXPECT_TRUE(write_success);
  // Before we hit the 1 second, remove the policy
  session_map = session_store->read_sessions({IMSI2});
  update      = SessionStore::get_default_session_update(session_map);
  SessionRules session_rules;
  auto rule_set_per_sub = session_rules.mutable_rules_per_subscriber()->Add();
  rule_set_per_sub->set_imsi(IMSI2);
  rule_set_per_sub->mutable_rule_set()->Add()->CopyFrom(
      create_rule_set(false, "apn1", {}, {}));  // Remove all rules
  local_enforcer->handle_set_session_rules(session_map, session_rules, update);
  // Write + Read in/from SessionStore
  write_success = session_store->update_sessions(update);
  EXPECT_TRUE(write_success);
  // Rule is removed, expect no bearer creation
  EXPECT_CALL(
      *spgw_client, create_dedicated_bearer(CheckCreateBearerReq(IMSI1, 1)))
      .Times(0);
  evb->loopOnce();
}

// Create multiple rules with QoS and assert dedicated bearers are created
// Simulate a policy->bearer mapping from MME, both success + failure.
// Simulate additional rule updates to trigger dedicated bearer deletions
TEST_F(LocalEnforcerTest, test_dedicated_bearer_lifecycle) {
  const uint32_t default_bearer_id = BEARER_ID_1;
  const uint32_t bearer_1          = BEARER_ID_2;
  const uint32_t bearer_2          = BEARER_ID_3;

  // Three rules with QoS & one without
  insert_static_rule_with_qos(0, "m1", "rule1", 1);  // QCI=1
  insert_static_rule_with_qos(0, "m1", "rule2", 2);  // QCI=2
  insert_static_rule_with_qos(0, "m1", "rule3", 3);  // QCI=3
  insert_static_rule(0, "m1", "rule4");

  // test_cfg_ is initialized with QoSInfo field w/ QCI 5
  test_cfg_.common_context.mutable_sid()->set_id(IMSI1);
  test_cfg_.common_context.set_apn("apn1");
  auto lte_context = test_cfg_.rat_specific_context.mutable_lte_context();
  lte_context->mutable_qos_info()->set_qos_class_id(5);
  lte_context->set_bearer_id(default_bearer_id);  // linked_bearer_id

  CreateSessionResponse response;
  response.mutable_static_rules()->Add()->set_rule_id("rule1");
  response.mutable_static_rules()->Add()->set_rule_id("rule2");
  response.mutable_static_rules()->Add()->set_rule_id("rule3");
  response.mutable_static_rules()->Add()->set_rule_id("rule4");

  // We only expect non-QoS policies to be installed immediately
  RuleToProcess expected_rule4 =
      make_rule_to_process("rule4", teids1.agw_teid(), teids1.enb_teid());
  EXPECT_CALL(
      *pipelined_client,
      activate_flows_for_rules(
          IMSI1, IP1, testing::_, CheckTeids(test_cfg_.common_context.teids()),
          test_cfg_.common_context.msisdn(), testing::_,
          CheckRulesToProcess(RulesToProcess{expected_rule4}), testing::_))
      .Times(1);

  std::unordered_set<std::string> no_install_ids({"rule1", "rule2", "rule3"});
  // Expect NO call to PipelineD for rule1,rule2,rule3 since they need bearer
  // activation
  EXPECT_CALL(
      *pipelined_client,
      activate_flows_for_rules(
          IMSI1, IP1, testing::_, CheckTeids(test_cfg_.common_context.teids()),
          test_cfg_.common_context.msisdn(), testing::_,
          CheckSubset(no_install_ids), testing::_))
      .Times(0);
  // expect only 1 rule in the request since only rules with a QoS field
  // should be mapped to a bearer
  EXPECT_CALL(
      *spgw_client, create_dedicated_bearer(CheckCreateBearerReq(IMSI1, 3)))
      .Times(1)
      .WillOnce(testing::Return(true));
  // For testing change the delay to 0 ms.
  LocalEnforcer::BEARER_CREATION_DELAY_ON_SESSION_INIT = 0;
  local_enforcer->init_session(
      session_map, IMSI1, SESSION_ID_1, test_cfg_, response);
  local_enforcer->update_tunnel_ids(
      session_map,
      create_update_tunnel_ids_request(IMSI1, default_bearer_id, teids1));
  // Write + Read in/from SessionStore
  bool write_success =
      session_store->create_sessions(IMSI1, std::move(session_map[IMSI1]));
  EXPECT_TRUE(write_success);
  session_map = session_store->read_sessions({IMSI1});
  auto update = SessionStore::get_default_session_update(session_map);
  // Progress the loop to run the scheduled bearer creation request
  evb->loopOnce();

  // Test successful creation of dedicated bearer for rule1 + rule2
  PolicyBearerBindingRequest bearer_bind_req_success1 =
      create_policy_bearer_bind_req(
          IMSI1, default_bearer_id, "rule1", bearer_1, 1, 2);
  PolicyBearerBindingRequest bearer_bind_req_success2 =
      create_policy_bearer_bind_req(
          IMSI1, default_bearer_id, "rule2", bearer_2, 3, 4);
  std::unordered_set<std::string> rule_ids({"rule1", "rule2"});
  // Expect NO call to PipelineD for rule1
  EXPECT_CALL(
      *pipelined_client, deactivate_flows_for_rules(
                             IMSI1, testing::_, testing::_, testing::_,
                             CheckSubset(rule_ids), testing::_))
      .Times(0);
  // All the policies that now have dedicated bearers can be activated now with
  // dedicated teids
  RuleToProcess expected_rule1 = make_rule_to_process("rule1", 1, 2);
  RuleToProcess expected_rule2 = make_rule_to_process("rule2", 3, 4);
  EXPECT_CALL(
      *pipelined_client,
      activate_flows_for_rules(
          IMSI1, IP1, testing::_, CheckTeids(test_cfg_.common_context.teids()),
          test_cfg_.common_context.msisdn(), testing::_,
          CheckRulesToProcess(RulesToProcess{expected_rule1}), testing::_))
      .Times(1);
  EXPECT_CALL(
      *pipelined_client,
      activate_flows_for_rules(
          IMSI1, IP1, testing::_, CheckTeids(test_cfg_.common_context.teids()),
          test_cfg_.common_context.msisdn(), testing::_,
          CheckRulesToProcess(RulesToProcess{expected_rule2}), testing::_))
      .Times(1);

  local_enforcer->bind_policy_to_bearer(
      session_map, bearer_bind_req_success1, update);
  local_enforcer->bind_policy_to_bearer(
      session_map, bearer_bind_req_success2, update);

  std::vector<Teids> existing_teids = session_map[IMSI1][0]->get_active_teids();
  EXPECT_EQ(3, existing_teids.size());  // default bearer + 2 dedicated bearers

  // Test unsuccessful creation of dedicated bearer for rule3 (bearer_id = 0)
  PolicyBearerBindingRequest bearer_bind_req_fail =
      create_policy_bearer_bind_req(IMSI1, default_bearer_id, "rule3", 0, 0, 0);
  // Since the dedicated bearer binding failed for rule3, we expect deactivate
  // request but not activate request
  no_install_ids = std::unordered_set<std::string>{"rule3"};
  EXPECT_CALL(
      *pipelined_client,
      activate_flows_for_rules(
          IMSI1, IP1, testing::_, CheckTeids(test_cfg_.common_context.teids()),
          test_cfg_.common_context.msisdn(), testing::_,
          CheckSubset(no_install_ids), testing::_))
      .Times(0);
  EXPECT_CALL(
      *pipelined_client, deactivate_flows_for_rules(
                             IMSI1, testing::_, testing::_, testing::_,
                             CheckPolicyID("rule3"), testing::_))
      .Times(1);
  local_enforcer->bind_policy_to_bearer(
      session_map, bearer_bind_req_fail, update);
  // Check update criteria has changes
  EXPECT_TRUE(update[IMSI1][SESSION_ID_1].is_bearer_mapping_updated);
  EXPECT_EQ(update[IMSI1][SESSION_ID_1].bearer_id_by_policy.size(), 2);
  // Write + Read in/from SessionStore
  write_success = session_store->update_sessions(update);
  EXPECT_TRUE(write_success);
  session_map = session_store->read_sessions({IMSI1});
  update      = SessionStore::get_default_session_update(session_map);

  // At this point we have rule1 -> bearer1, rule2 -> bearer2 rule3 -> deleted,
  // and rule4 -> no bearer. When we remove rule1, we expect to see a delete
  // dedicated bearer request. Use the set rule interface to remove rule1.
  SessionRules session_rules;
  auto rule_set_per_sub = session_rules.mutable_rules_per_subscriber()->Add();
  rule_set_per_sub->set_imsi(IMSI1);
  rule_set_per_sub->mutable_rule_set()->Add()->CopyFrom(
      create_rule_set(false, "apn1", {"rule2", "rule4"}, {}));
  EXPECT_CALL(
      *pipelined_client, deactivate_flows_for_rules(
                             IMSI1, testing::_, testing::_, testing::_,
                             CheckPolicyID("rule1"), testing::_))
      .Times(1);
  EXPECT_CALL(
      *spgw_client, delete_dedicated_bearer(CheckDeleteOneBearerReq(
                        IMSI1, default_bearer_id, bearer_1)))
      .Times(1);
  update = SessionStore::get_default_session_update(session_map);
  local_enforcer->handle_set_session_rules(session_map, session_rules, update);

  // Finally remove rule2 via an update response
  UpdateSessionResponse update_response;
  auto monitor_update =
      update_response.mutable_usage_monitor_responses()->Add();
  create_monitor_update_response(
      IMSI1, SESSION_ID_1, "m1", MonitoringLevel::PCC_RULE_LEVEL, 1024,
      monitor_update);
  monitor_update->add_rules_to_remove("rule2");

  EXPECT_CALL(
      *pipelined_client, deactivate_flows_for_rules(
                             IMSI1, testing::_, testing::_, testing::_,
                             CheckPolicyID("rule2"), testing::_))
      .Times(1);
  EXPECT_CALL(
      *spgw_client, delete_dedicated_bearer(CheckDeleteOneBearerReq(
                        IMSI1, default_bearer_id, bearer_2)))
      .Times(1);
  local_enforcer->update_session_credits_and_rules(
      session_map, update_response, update);
  // Check update criteria has changes + no bearer left!
  EXPECT_TRUE(update[IMSI1][SESSION_ID_1].is_bearer_mapping_updated);
  EXPECT_EQ(update[IMSI1][SESSION_ID_1].bearer_id_by_policy.size(), 0);
}

// Test the handle_set_session_rules function to apply rules to sessions
// We will test a case where a subscriber has 2 separate sessions
TEST_F(LocalEnforcerTest, test_set_session_rules) {
  SessionConfig config1, config2;

  // Set Session1 + Session2 to be LTE session with default QoS
  QosInformationRequest qos_info;
  qos_info.set_apn_ambr_dl(32);
  qos_info.set_apn_ambr_dl(64);
  qos_info.set_qos_class_id(1);
  const auto& lte_context1 =
      build_lte_context(IP2, "", "", "", "", BEARER_ID_1, &qos_info);
  config1.common_context =
      build_common_context(IMSI1, IP1, "", teids1, "apn1", "msisdn1", TGPP_LTE);
  config1.rat_specific_context.mutable_lte_context()->CopyFrom(lte_context1);

  const auto& lte_context2 =
      build_lte_context(IP2, "", "", "", "", BEARER_ID_2, &qos_info);
  config2.common_context =
      build_common_context(IMSI1, IP2, "", teids2, "apn2", "msisdn1", TGPP_LTE);
  config2.rat_specific_context.mutable_lte_context()->CopyFrom(lte_context2);

  // Initialize 3 static rules in RuleStore and create 2 dynamic rules
  insert_static_rule_with_qos(0, "m1", "static1", 2);
  insert_static_rule(0, "m1", "static2");
  insert_static_rule_with_qos(0, "m1", "static3", 3);
  PolicyRule dynamic_1 = create_policy_rule("dynamic1", "m1", 0);
  PolicyRule dynamic_2 = create_policy_rule("dynamic2", "m1", 0);

  // Create a session with static1/static2/dynamic1
  CreateSessionResponse response;
  response.mutable_static_rules()->Add()->set_rule_id("static1");
  response.mutable_static_rules()->Add()->set_rule_id("static2");
  response.mutable_dynamic_rules()->Add()->mutable_policy_rule()->CopyFrom(
      dynamic_1);

  local_enforcer->init_session(
      session_map, IMSI1, SESSION_ID_1, config1, response);
  local_enforcer->update_tunnel_ids(
      session_map,
      create_update_tunnel_ids_request(IMSI1, BEARER_ID_1, teids1));
  local_enforcer->init_session(
      session_map, IMSI1, SESSION_ID_2, config2, response);
  local_enforcer->update_tunnel_ids(
      session_map,
      create_update_tunnel_ids_request(IMSI1, BEARER_ID_2, teids2));
  evb->loopOnce();
  evb->loopOnce();
  // Assert the rules exist
  EXPECT_TRUE(session_map[IMSI1][0]->is_static_rule_installed("static1"));
  EXPECT_TRUE(session_map[IMSI1][0]->is_static_rule_installed("static2"));
  EXPECT_TRUE(session_map[IMSI1][0]->is_dynamic_rule_installed("dynamic1"));

  bool success =
      session_store->create_sessions(IMSI1, std::move(session_map[IMSI1]));
  EXPECT_TRUE(success);

  // Apply a set rule of
  // apn1 -> static1,static3,dynamic2
  // apn2 -> static2
  // subscriber_wide -> dynamic2
  // This should lead to the following actions:
  // apn1 ->
  //          add          : dynamic_2
  //          remove       : static2, dynamic1
  //          create bearer: static3
  // apn2 ->
  //          add          : dynamic_2
  //          remove       : static2, dynamic1
  //          create bearer:
  SessionRules session_rules;
  auto rule_set_per_sub = session_rules.mutable_rules_per_subscriber()->Add();
  rule_set_per_sub->set_imsi(IMSI1);
  rule_set_per_sub->mutable_rule_set()->Add()->CopyFrom(
      create_rule_set(false, "apn1", {"static1", "static3"}, {dynamic_2}));
  rule_set_per_sub->mutable_rule_set()->Add()->CopyFrom(
      create_rule_set(false, "apn2", {"static1"}, {}));
  rule_set_per_sub->mutable_rule_set()->Add()->CopyFrom(
      create_rule_set(true, "", {}, {dynamic_2}));

  // PipelineD expectations for Session1
  EXPECT_CALL(
      *pipelined_client,
      activate_flows_for_rules(
          IMSI1, IP1, testing::_, CheckTeids(config1.common_context.teids()),
          config1.common_context.msisdn(), testing::_,
          CheckRuleNames(std::vector<std::string>{"dynamic2"}), testing::_))
      .Times(1);
  // PipelineD expectations for Session2
  EXPECT_CALL(
      *pipelined_client,
      activate_flows_for_rules(
          IMSI1, IP2, testing::_, CheckTeids(config2.common_context.teids()),
          config2.common_context.msisdn(), testing::_,
          CheckRuleNames(std::vector<std::string>{"dynamic2"}), testing::_))
      .Times(1);
  // For both Session1 + Session2
  std::unordered_set<std::string> deactivate_ids = {"dynamic1", "static2"};
  EXPECT_CALL(
      *pipelined_client, deactivate_flows_for_rules(
                             IMSI1, testing::_, testing::_, testing::_,
                             CheckPolicyIDs(2, deactivate_ids), testing::_))
      .Times(2);

  // Since static3 is also a QoS rule with a new QCI (not equal to default), we
  // should also expect a create bearer request here
  EXPECT_CALL(
      *spgw_client, create_dedicated_bearer(CheckCreateBearerReq(IMSI1, 1)))
      .Times(1)
      .WillOnce(testing::Return(true));

  session_map = session_store->read_sessions(SessionRead{IMSI1});
  auto update = SessionStore::get_default_session_update(session_map);
  local_enforcer->handle_set_session_rules(session_map, session_rules, update);
}

TEST_F(LocalEnforcerTest, test_rar_session_not_found) {
  // verify session validity by passing in an invalid IMSI
  PolicyReAuthRequest rar;
  create_policy_reauth_request(
      "session1", IMSI1, {}, {}, {}, {}, time(nullptr), {}, &rar);
  PolicyReAuthAnswer raa;
  auto update = SessionStore::get_default_session_update(session_map);
  local_enforcer->init_policy_reauth(session_map, rar, raa, update);
  EXPECT_EQ(raa.result(), ReAuthResult::SESSION_NOT_FOUND);

  // verify session validity passing in a valid IMSI (IMSI1)
  // and an invalid session-id (session1)
  CreateSessionResponse response;
  local_enforcer->init_session(
      session_map, IMSI1, "session0", test_cfg_, response);
  local_enforcer->update_tunnel_ids(
      session_map,
      create_update_tunnel_ids_request(IMSI1, BEARER_ID_1, teids1));
  local_enforcer->init_policy_reauth(session_map, rar, raa, update);
  EXPECT_EQ(raa.result(), ReAuthResult::SESSION_NOT_FOUND);
}

TEST_F(LocalEnforcerTest, test_revalidation_timer_on_init) {
  const std::string mkey = "m1";
  insert_static_rule(1, mkey, "rule1");

  // Create a CreateSessionResponse with one Gx monitor, PCC rule, and an event
  // trigger
  CreateSessionResponse response;
  auto monitor = response.mutable_usage_monitors()->Add();
  std::vector<EventTrigger> event_triggers{EventTrigger::REVALIDATION_TIMEOUT};
  create_monitor_update_response(
      IMSI1, SESSION_ID_1, mkey, MonitoringLevel::PCC_RULE_LEVEL, 1024,
      monitor);

  response.add_event_triggers(EventTrigger::REVALIDATION_TIMEOUT);
  response.mutable_revalidation_time()->set_seconds(time(nullptr));

  StaticRuleInstall static_rule_install;
  static_rule_install.set_rule_id("rule1");
  response.mutable_static_rules()->Add()->CopyFrom(static_rule_install);

  local_enforcer->init_session(
      session_map, IMSI1, SESSION_ID_1, default_cfg_1, response);
  local_enforcer->update_tunnel_ids(
      session_map,
      create_update_tunnel_ids_request(IMSI1, BEARER_ID_1, teids1));
  bool success =
      session_store->create_sessions(IMSI1, std::move(session_map[IMSI1]));
  EXPECT_TRUE(success);
  // schedule_revalidation puts two things on the event loop
  // updates the session's state with revalidation_requested_ set to true
  evb->loopOnce();
  evb->loopOnce();
  auto update = SessionStore::get_default_session_update(session_map);
  session_map = session_store->read_sessions(SessionRead{IMSI1});
  std::vector<std::unique_ptr<ServiceAction>> actions;
  auto updates = local_enforcer->collect_updates(session_map, actions, update);
  EXPECT_EQ(updates.usage_monitors_size(), 1);
}

TEST_F(LocalEnforcerTest, test_revalidation_timer_on_rar) {
  CreateSessionResponse response;
  PolicyReAuthRequest rar;
  PolicyReAuthAnswer raa;
  std::vector<std::string> rules_to_install;
  std::vector<EventTrigger> event_triggers{EventTrigger::REVALIDATION_TIMEOUT};

  const std::string mkey = "m1";
  rules_to_install.push_back("rule1");
  insert_static_rule(1, mkey, "rule1");

  // Create a CreateSessionResponse with one Gx monitor, PCC rule
  create_session_create_response(
      IMSI1, SESSION_ID_1, mkey, rules_to_install, &response);

  local_enforcer->init_session(
      session_map, IMSI1, SESSION_ID_1, default_cfg_1, response);
  local_enforcer->update_tunnel_ids(
      session_map,
      create_update_tunnel_ids_request(IMSI1, BEARER_ID_1, teids1));
  EXPECT_EQ(session_map[IMSI1].size(), 1);

  // Write and read into session store, assert success
  bool success =
      session_store->create_sessions(IMSI1, std::move(session_map[IMSI1]));
  EXPECT_TRUE(success);
  session_map = session_store->read_sessions(SessionRead{IMSI1});

  // Create a RaR with a REVALIDATION event trigger
  create_policy_reauth_request(
      SESSION_ID_1, IMSI1, {}, {}, {}, event_triggers, time(nullptr), {}, &rar);

  auto update = SessionStore::get_default_session_update(session_map);
  // This should trigger a revalidation to be scheduled
  local_enforcer->init_policy_reauth(session_map, rar, raa, update);
  EXPECT_EQ(raa.result(), ReAuthResult::UPDATE_INITIATED);
  // Propagate the change to store
  success = session_store->update_sessions(update);
  EXPECT_TRUE(success);

  // schedule_revalidation puts two things on the event loop
  // updates the session's state with revalidation_requested_ set to true
  evb->loopOnce();
  evb->loopOnce();
  session_map = session_store->read_sessions(SessionRead{IMSI1});
  std::vector<std::unique_ptr<ServiceAction>> actions;
  auto updates = local_enforcer->collect_updates(session_map, actions, update);
  EXPECT_EQ(updates.usage_monitors_size(), 1);
}

TEST_F(LocalEnforcerTest, test_revalidation_timer_on_update) {
  CreateSessionResponse create_response;
  UpdateSessionResponse update_response;
  std::vector<std::string> rules_to_install;
  std::vector<EventTrigger> event_triggers{EventTrigger::REVALIDATION_TIMEOUT};
  const std::string mkey1 = "m1";
  const std::string mkey2 = "m2";
  rules_to_install.push_back("rule1");
  insert_static_rule(1, mkey1, "rule1");

  // Create a CreateSessionResponse with one Gx monitor, PCC rule
  create_session_create_response(
      IMSI1, SESSION_ID_1, mkey1, rules_to_install, &create_response);
  local_enforcer->init_session(
      session_map, IMSI1, SESSION_ID_1, default_cfg_1, create_response);
  local_enforcer->update_tunnel_ids(
      session_map,
      create_update_tunnel_ids_request(IMSI1, BEARER_ID_1, teids1));

  create_response.Clear();
  create_session_create_response(
      IMSI2, SESSION_ID_2, mkey1, rules_to_install, &create_response);
  auto test_cfg_2 = test_cfg_;
  test_cfg_2.common_context.mutable_teids()->CopyFrom(teids2);
  test_cfg_2.rat_specific_context.mutable_lte_context()->set_bearer_id(
      BEARER_ID_2);
  local_enforcer->init_session(
      session_map, IMSI2, SESSION_ID_2, test_cfg_2, create_response);
  local_enforcer->update_tunnel_ids(
      session_map,
      create_update_tunnel_ids_request(IMSI2, BEARER_ID_2, teids2));

  // Write and read into session store, assert success
  bool success =
      session_store->create_sessions(IMSI1, std::move(session_map[IMSI1]));
  EXPECT_TRUE(success);
  success =
      session_store->create_sessions(IMSI2, std::move(session_map[IMSI2]));
  EXPECT_TRUE(success);
  session_map = session_store->read_sessions(SessionRead{IMSI1, IMSI2});
  EXPECT_EQ(session_map[IMSI1].size(), 1);
  EXPECT_EQ(session_map[IMSI2].size(), 1);

  auto revalidation_timer = time(nullptr);
  // IMSI1 has two separate monitors with the same revalidation timer
  // IMSI2 does not have a revalidation timer
  auto monitor = update_response.mutable_usage_monitor_responses()->Add();
  create_monitor_update_response(
      IMSI1, SESSION_ID_1, mkey1, MonitoringLevel::PCC_RULE_LEVEL, 1024,
      event_triggers, revalidation_timer, monitor);
  monitor = update_response.mutable_usage_monitor_responses()->Add();
  create_monitor_update_response(
      IMSI1, SESSION_ID_1, mkey2, MonitoringLevel::PCC_RULE_LEVEL, 1024,
      event_triggers, revalidation_timer, monitor);
  monitor = update_response.mutable_usage_monitor_responses()->Add();
  create_monitor_update_response(
      IMSI1, SESSION_ID_1, mkey1, MonitoringLevel::PCC_RULE_LEVEL, 1024,
      monitor);
  auto update = SessionStore::get_default_session_update(session_map);
  // This should trigger a revalidation to be scheduled
  local_enforcer->update_session_credits_and_rules(
      session_map, update_response, update);
  // Propagate the change to store
  success = session_store->update_sessions(update);
  EXPECT_TRUE(success);

  // a single schedule_revalidation puts two things on the event loop
  // updates the session's state with revalidation_requested_ set to true
  evb->loopOnce();
  evb->loopOnce();

  session_map = session_store->read_sessions(SessionRead{IMSI1});
  std::vector<std::unique_ptr<ServiceAction>> actions;
  auto updates = local_enforcer->collect_updates(session_map, actions, update);
  EXPECT_EQ(updates.usage_monitors_size(), 1);
}

TEST_F(LocalEnforcerTest, test_revalidation_timer_on_update_no_monitor) {
  CreateSessionResponse create_response;
  UpdateSessionResponse update_response;
  StaticRuleInstall rule_install;

  rule_install.set_rule_id("rule1");
  insert_static_rule(1, "m1", "rule1");

  // Create a CreateSessionResponse with Revalidation Timeout, PCC rule
  auto res = create_response.mutable_usage_monitors()->Add();
  res->set_success(true);
  res->set_sid(IMSI1);
  res->set_session_id(SESSION_ID_1);
  create_response.mutable_static_rules()->Add()->CopyFrom(rule_install);
  local_enforcer->init_session(
      session_map, IMSI1, SESSION_ID_1, default_cfg_1, create_response);
  local_enforcer->update_tunnel_ids(
      session_map,
      create_update_tunnel_ids_request(IMSI1, BEARER_ID_1, teids1));

  // Write and read into session store, assert success
  bool success =
      session_store->create_sessions(IMSI1, std::move(session_map[IMSI1]));
  EXPECT_TRUE(success);

  session_map = session_store->read_sessions(SessionRead{IMSI1});
  EXPECT_EQ(session_map[IMSI1].size(), 1);

  // IMSI1 has no monitor credit and a revalidation timeout
  auto monitor = update_response.mutable_usage_monitor_responses()->Add();
  monitor->set_success(true);
  monitor->set_sid(IMSI1);
  monitor->set_session_id(SESSION_ID_1);
  monitor->add_event_triggers(EventTrigger::REVALIDATION_TIMEOUT);
  monitor->mutable_revalidation_time()->set_seconds(time(nullptr));
  auto update = SessionStore::get_default_session_update(session_map);
  // This should trigger a revalidation to be scheduled
  local_enforcer->update_session_credits_and_rules(
      session_map, update_response, update);
  // Propagate the change to store
  success = session_store->update_sessions(update);
  EXPECT_TRUE(success);

  // a single schedule_revalidation puts two things on the event loop
  // updates the session's state with revalidation_requested_ set to true
  evb->loopOnce();
  evb->loopOnce();

  session_map = session_store->read_sessions(SessionRead{IMSI1});
  std::vector<std::unique_ptr<ServiceAction>> actions;
  auto updates = local_enforcer->collect_updates(session_map, actions, update);
  EXPECT_EQ(updates.usage_monitors_size(), 1);
}

TEST_F(LocalEnforcerTest, test_pipelined_cwf_setup) {
  // insert into rule store first so init_session_credit can find the rule
  insert_static_rule(1, "", "rule2");

  CreateSessionResponse response;
  create_credit_update_response(
      IMSI1, SESSION_ID_1, 1, 1024, response.mutable_credits()->Add());
  auto epoch        = 145;
  auto dynamic_rule = response.mutable_dynamic_rules()->Add();
  auto policy_rule  = dynamic_rule->mutable_policy_rule();
  policy_rule->set_id("rule1");
  policy_rule->set_rating_group(1);
  policy_rule->set_tracking_type(PolicyRule::ONLY_OCS);
  auto static_rule = response.mutable_static_rules()->Add();
  static_rule->set_rule_id("rule2");
  SessionConfig test_cwf_cfg1;
  test_cwf_cfg1.common_context = build_common_context(
      IMSI1, IP1, "", teids0, "01-a1-20-c2-0f-bb:CWC_OFFLOAD", "msisdn1",
      TGPP_WLAN);
  const auto& wlan = build_wlan_context("11:22:00:00:22:11", "5555");
  test_cwf_cfg1.rat_specific_context.mutable_wlan_context()->CopyFrom(wlan);
  local_enforcer->init_session(
      session_map, IMSI1, SESSION_ID_1, test_cwf_cfg1, response);
  local_enforcer->update_tunnel_ids(
      session_map, create_update_tunnel_ids_request(IMSI1, 0, teids0));

  CreateSessionResponse response2;
  create_credit_update_response(
      IMSI2, SESSION_ID_2, 1, 2048, response2.mutable_credits()->Add());
  auto dynamic_rule2 = response2.mutable_dynamic_rules()->Add();
  auto policy_rule2  = dynamic_rule2->mutable_policy_rule();
  policy_rule2->set_id("rule22");
  policy_rule2->set_rating_group(1);
  policy_rule2->set_tracking_type(PolicyRule::ONLY_OCS);
  SessionConfig test_cwf_cfg2;
  test_cwf_cfg2.common_context = build_common_context(
      IMSI2, IP1, "", teids0, "03-21-00-02-00-20:Magma", "msisdn2", TGPP_WLAN);
  const auto& wlan2 = build_wlan_context("00:00:00:00:00:02", "5555");
  test_cwf_cfg2.rat_specific_context.mutable_wlan_context()->CopyFrom(wlan2);
  local_enforcer->init_session(
      session_map, IMSI2, SESSION_ID_2, test_cwf_cfg2, response2);
  local_enforcer->update_tunnel_ids(
      session_map, create_update_tunnel_ids_request(IMSI2, 0, teids0));

  std::vector<std::string> imsi_list              = {IMSI2, IMSI1};
  std::vector<std::string> ip_address_list        = {IP1, IP1};
  std::vector<std::string> ipv6_address_list      = {"", ""};
  std::vector<std::vector<std::string>> rule_list = {{"rule22"},
                                                     {"rule1", "rule2"}};
  std::vector<std::vector<uint32_t>> version_list = {{1}, {1, 1}};

  std::vector<std::string> ue_mac_addrs  = {"00:00:00:00:00:02",
                                           "11:22:00:00:22:11"};
  std::vector<std::string> msisdns       = {"msisdn2", "msisdn1"};
  std::vector<std::string> apn_mac_addrs = {"03-21-00-02-00-20",
                                            "01-a1-20-c2-0f-bb"};
  std::vector<std::string> apn_names     = {"Magma", "CWC_OFFLOAD"};
  EXPECT_CALL(
      *pipelined_client, setup_cwf(
                             CheckSessionInfos(
                                 imsi_list, ip_address_list, ipv6_address_list,
                                 test_cwf_cfg2, rule_list, version_list),
                             testing::_, ue_mac_addrs, msisdns, apn_mac_addrs,
                             apn_names, testing::_, testing::_, testing::_))
      .Times(1);

  local_enforcer->setup(
      session_map, epoch, [](Status status, SetupFlowsResult resp) {});
}

TEST_F(LocalEnforcerTest, test_pipelined_lte_setup) {
  // insert into rule store first so init_session_credit can find the rule
  insert_static_rule(1, "", "rule2");

  CreateSessionResponse response;
  create_credit_update_response(
      IMSI1, SESSION_ID_1, 1, 1024, response.mutable_credits()->Add());
  auto epoch        = 145;
  auto dynamic_rule = response.mutable_dynamic_rules()->Add();
  auto policy_rule  = dynamic_rule->mutable_policy_rule();
  policy_rule->set_id("rule1");
  policy_rule->set_rating_group(1);
  policy_rule->set_tracking_type(PolicyRule::ONLY_OCS);
  auto static_rule = response.mutable_static_rules()->Add();
  static_rule->set_rule_id("rule2");

  default_cfg_1.common_context.set_ue_ipv6(IPv6_1);
  local_enforcer->init_session(
      session_map, IMSI1, SESSION_ID_1, default_cfg_1, response);
  local_enforcer->update_tunnel_ids(
      session_map,
      create_update_tunnel_ids_request(IMSI1, BEARER_ID_1, teids1));

  CreateSessionResponse response2;
  create_credit_update_response(
      IMSI2, SESSION_ID_2, 1, 2048, response2.mutable_credits()->Add());
  auto dynamic_rule2 = response2.mutable_dynamic_rules()->Add();
  auto policy_rule2  = dynamic_rule2->mutable_policy_rule();
  policy_rule2->set_id("rule22");
  policy_rule2->set_rating_group(1);
  policy_rule2->set_tracking_type(PolicyRule::ONLY_OCS);

  default_cfg_2.common_context.mutable_teids()->CopyFrom(teids2);
  default_cfg_2.rat_specific_context.mutable_lte_context()->set_bearer_id(
      BEARER_ID_2);
  local_enforcer->init_session(
      session_map, IMSI2, SESSION_ID_2, default_cfg_2, response2);
  local_enforcer->update_tunnel_ids(
      session_map,
      create_update_tunnel_ids_request(IMSI2, BEARER_ID_2, teids2));

  std::vector<std::string> imsi_list              = {IMSI2, IMSI1};
  std::vector<std::string> ip_address_list        = {IP1, IP1};
  std::vector<std::string> ipv6_address_list      = {IPv6_1, ""};
  std::vector<std::vector<std::string>> rule_list = {{"rule22"},
                                                     {"rule1", "rule2"}};
  std::vector<std::vector<uint32_t>> version_list = {{1}, {1, 1}};

  std::vector<std::string> ue_mac_addrs  = {"00:00:00:00:00:02",
                                           "11:22:00:00:22:11"};
  std::vector<std::string> msisdns       = {"msisdn2", "msisdn1"};
  std::vector<std::string> apn_mac_addrs = {"03-21-00-02-00-20",
                                            "01-a1-20-c2-0f-bb"};
  std::vector<std::string> apn_names     = {"Magma", "CWC_OFFLOAD"};
  EXPECT_CALL(
      *pipelined_client, setup_lte(
                             CheckSessionInfos(
                                 imsi_list, ip_address_list, ipv6_address_list,
                                 test_cfg_, rule_list, version_list),
                             testing::_, testing::_))
      .Times(1);

  local_enforcer->setup(
      session_map, epoch, [](Status status, SetupFlowsResult resp) {});
}

TEST_F(LocalEnforcerTest, test_valid_apn_parsing) {
  insert_static_rule(1, "", "rule1");
  int epoch = 145;
  CreateSessionResponse response;
  SessionUpdate session_update =
      SessionStore::get_default_session_update(session_map);

  auto credits = response.mutable_credits();
  create_credit_update_response(IMSI1, SESSION_ID_1, 1, 1024, credits->Add());

  auto apn = "03-21-00-02-00-20:Magma";
  SessionConfig test_cwf_cfg;
  test_cwf_cfg.common_context =
      build_common_context(IMSI1, "", "", teids0, apn, MSISDN, TGPP_WLAN);
  const auto& wlan = build_wlan_context(MAC_ADDR, RADIUS_SESSION_ID);
  test_cwf_cfg.common_context.set_rat_type(TGPP_WLAN);
  test_cwf_cfg.rat_specific_context.mutable_wlan_context()->CopyFrom(wlan);

  local_enforcer->init_session(
      session_map, IMSI1, SESSION_ID_1, test_cwf_cfg, response);
  local_enforcer->update_tunnel_ids(
      session_map, create_update_tunnel_ids_request(IMSI1, 0, teids0));

  std::vector<std::string> ue_mac_addrs  = {MAC_ADDR};
  std::vector<std::string> msisdns       = {MSISDN};
  std::vector<std::string> apn_mac_addrs = {"03-21-00-02-00-20"};
  std::vector<std::string> apn_names     = {"Magma"};

  EXPECT_CALL(
      *pipelined_client,
      setup_cwf(
          testing::_, testing::_, ue_mac_addrs, msisdns, apn_mac_addrs,
          apn_names, testing::_, epoch, testing::_))
      .Times(1);

  local_enforcer->setup(
      session_map, epoch, [](Status status, SetupFlowsResult resp) {});
}

TEST_F(LocalEnforcerTest, test_invalid_apn_parsing) {
  insert_static_rule(1, "", "rule1");
  int epoch = 145;
  CreateSessionResponse response;
  SessionUpdate session_update =
      SessionStore::get_default_session_update(session_map);

  auto credits = response.mutable_credits();
  create_credit_update_response(IMSI1, SESSION_ID_1, 1, 1024, credits->Add());

  auto apn = "03-0BLAHBLAH0-00-02-00-20:ThisIsNotOkay";
  SessionConfig test_cwf_cfg;
  test_cwf_cfg.common_context =
      build_common_context(IMSI1, IP1, "", teids1, apn, MSISDN, TGPP_WLAN);
  const auto& wlan = build_wlan_context(MAC_ADDR, RADIUS_SESSION_ID);
  test_cwf_cfg.common_context.set_rat_type(TGPP_WLAN);
  test_cwf_cfg.rat_specific_context.mutable_wlan_context()->CopyFrom(wlan);

  local_enforcer->init_session(
      session_map, IMSI1, SESSION_ID_1, test_cwf_cfg, response);
  local_enforcer->update_tunnel_ids(
      session_map,
      create_update_tunnel_ids_request(IMSI1, BEARER_ID_1, teids1));

  std::vector<std::string> ue_mac_addrs  = {MAC_ADDR};
  std::vector<std::string> msisdns       = {MSISDN};
  std::vector<std::string> apn_mac_addrs = {""};
  std::vector<std::string> apn_names     = {
      "03-0BLAHBLAH0-00-02-00-20:ThisIsNotOkay"};

  EXPECT_CALL(
      *pipelined_client,
      setup_cwf(
          testing::_, testing::_, ue_mac_addrs, msisdns, apn_mac_addrs,
          apn_names, testing::_, epoch, testing::_))
      .Times(1);

  local_enforcer->setup(
      session_map, epoch, [](Status status, SetupFlowsResult resp) {});
}

TEST_F(LocalEnforcerTest, test_final_unit_redirect_activation_and_termination) {
  CreateSessionResponse response;
  test_cfg_.common_context.mutable_sid()->set_id(IMSI1);
  create_credit_update_response(
      IMSI1, SESSION_ID_1, 1, 1024, ChargingCredit_FinalAction_REDIRECT,
      "12.7.7.4", "", response.mutable_credits()->Add());

  auto static_rule = response.mutable_static_rules()->Add();
  static_rule->set_rule_id("static_1");

  insert_static_rule(1, "", "static_1");
  auto& ip_addr   = test_cfg_.common_context.ue_ipv4();
  auto& ipv6_addr = test_cfg_.common_context.ue_ipv6();
  auto& teids     = teids1;
  // The activation for the static rules (rule1,rule3) and dynamic rule (rule2)
  EXPECT_CALL(
      *pipelined_client, activate_flows_for_rules(
                             IMSI1, ip_addr, ipv6_addr, CheckTeids(teids),
                             test_cfg_.common_context.msisdn(), testing::_,
                             CheckRuleCount(1), testing::_))
      .Times(1);
  local_enforcer->init_session(
      session_map, IMSI1, SESSION_ID_1, test_cfg_, response);
  local_enforcer->update_tunnel_ids(
      session_map,
      create_update_tunnel_ids_request(IMSI1, BEARER_ID_1, teids1));

  // Insert record and aggregate over them
  RuleRecordTable table;
  auto record_list = table.mutable_records();
  create_rule_record(
      IMSI1, teids1.agw_teid(), "static_1", 1024, 2048, record_list->Add());
  auto update = SessionStore::get_default_session_update(session_map);
  local_enforcer->aggregate_records(session_map, table, update);

  // Collect actions and verify that restrict action is in the list
  std::vector<std::unique_ptr<ServiceAction>> actions;
  auto usage_updates =
      local_enforcer->collect_updates(session_map, actions, update);
  EXPECT_EQ(actions.size(), 1);
  EXPECT_EQ(actions[0]->get_type(), REDIRECT);
  PolicyRule redirect_rule = actions[0]->get_gy_rules_to_install()[0].rule;
  EXPECT_EQ(redirect_rule.redirect().server_address(), "12.7.7.4");

  EXPECT_CALL(
      *pipelined_client,
      add_gy_final_action_flow(
          IMSI1, ip_addr, ipv6_addr, CheckTeids(teids),
          test_cfg_.common_context.msisdn(), CheckRuleCount(1)))
      .Times(1);
  // Execute actions and asset final action state
  local_enforcer->execute_actions(session_map, actions, update);
  const CreditKey& credit_key(1);
  assert_session_is_in_final_state(
      session_map, IMSI1, SESSION_ID_1, credit_key, true);

  // the request should has no rules so PipelineD deletes all rules
  std::vector<Teids> expected_teids_vec{teids};
  EXPECT_CALL(
      *pipelined_client,
      deactivate_flows_for_rules_for_termination(
          IMSI1, ip_addr, ipv6_addr, CheckTeidVector(expected_teids_vec),
          RequestOriginType::WILDCARD));
  local_enforcer->handle_termination_from_access(
      session_map, IMSI1, APN1, update);
}

TEST_F(LocalEnforcerTest, test_final_unit_activation_and_canceling) {
  CreateSessionResponse response;
  create_credit_update_response(
      IMSI1, SESSION_ID_1, 1, 1024, ChargingCredit_FinalAction_RESTRICT_ACCESS,
      "", "rule1", response.mutable_credits()->Add());

  auto dynamic_rule = response.mutable_dynamic_rules()->Add();
  auto policy_rule  = dynamic_rule->mutable_policy_rule();
  policy_rule->set_id("rule2");
  policy_rule->set_rating_group(1);
  policy_rule->set_tracking_type(PolicyRule::ONLY_OCS);
  auto static_rule = response.mutable_static_rules()->Add();
  static_rule->set_rule_id("rule1");
  static_rule = response.mutable_static_rules()->Add();
  static_rule->set_rule_id("rule3");

  insert_static_rule(1, "", "rule1");
  insert_static_rule(1, "", "rule3");
  auto& ip_addr      = default_cfg_1.common_context.ue_ipv4();
  auto& ipv6_addr    = default_cfg_1.common_context.ue_ipv6();
  auto teids         = default_cfg_1.common_context.teids();
  const auto& msisdn = default_cfg_1.common_context.msisdn();
  // The activation for the static rules (rule1,rule3) and dynamic rule (rule2)
  EXPECT_CALL(
      *pipelined_client, activate_flows_for_rules(
                             IMSI1, ip_addr, ipv6_addr, CheckTeids(teids),
                             msisdn, testing::_, CheckRuleCount(3), testing::_))
      .Times(1);

  local_enforcer->init_session(
      session_map, IMSI1, SESSION_ID_1, default_cfg_1, response);
  local_enforcer->update_tunnel_ids(
      session_map,
      create_update_tunnel_ids_request(IMSI1, BEARER_ID_1, teids1));

  // Insert record and aggregate over them
  RuleRecordTable table;
  auto record_list = table.mutable_records();
  create_rule_record(
      IMSI1, teids.agw_teid(), "rule1", 1024, 2048, record_list->Add());
  create_rule_record(
      IMSI1, teids.agw_teid(), "rule2", 1024, 2048, record_list->Add());
  create_rule_record(
      IMSI1, teids.agw_teid(), "rule3", 1024, 2048, record_list->Add());
  auto update = SessionStore::get_default_session_update(session_map);
  local_enforcer->aggregate_records(session_map, table, update);

  // Collect actions and verify that restrict action is in the list
  std::vector<std::unique_ptr<ServiceAction>> actions;
  auto usage_updates =
      local_enforcer->collect_updates(session_map, actions, update);
  EXPECT_EQ(actions.size(), 1);
  EXPECT_EQ(actions[0]->get_type(), RESTRICT_ACCESS);
  EXPECT_EQ(actions[0]->get_gy_rules_to_install()[0].rule.id(), "rule1");

  EXPECT_CALL(
      *pipelined_client, add_gy_final_action_flow(
                             IMSI1, ip_addr, ipv6_addr, CheckTeids(teids),
                             msisdn, CheckRuleCount(1)))
      .Times(1);
  // Execute actions and asset final action state
  local_enforcer->execute_actions(session_map, actions, update);

  const CreditKey& credit_key(1);
  assert_session_is_in_final_state(
      session_map, IMSI1, SESSION_ID_1, credit_key, true);

  // Send a ReAuth request
  ChargingReAuthRequest reauth;
  reauth.set_sid(IMSI1);
  reauth.set_session_id(SESSION_ID_1);
  reauth.set_charging_key(1);
  reauth.set_type(ChargingReAuthRequest::SINGLE_SERVICE);
  update = SessionStore::get_default_session_update(session_map);
  auto result =
      local_enforcer->init_charging_reauth(session_map, reauth, update);
  EXPECT_EQ(result, ReAuthResult::UPDATE_INITIATED);

  actions.clear();
  auto update_req =
      local_enforcer->collect_updates(session_map, actions, update);
  local_enforcer->execute_actions(session_map, actions, update);
  EXPECT_EQ(update_req.updates_size(), 1);
  EXPECT_EQ(update_req.updates(0).common_context().sid().id(), IMSI1);
  EXPECT_EQ(update_req.updates(0).usage().type(), CreditUsage::REAUTH_REQUIRED);

  // Give credit after ReAuth
  UpdateSessionResponse update_response;
  auto updates = update_response.mutable_responses();
  create_credit_update_response(IMSI1, SESSION_ID_1, 1, 4096, updates->Add());
  local_enforcer->update_session_credits_and_rules(
      session_map, update_response, update);

  // when next update is collected, this should trigger an action to activate
  // the flow in pipelined
  EXPECT_CALL(
      *pipelined_client, activate_flows_for_rules(
                             IMSI1, ip_addr, ipv6_addr, CheckTeids(teids),
                             msisdn, testing::_, testing::_, testing::_))
      .Times(1);
  actions.clear();
  local_enforcer->collect_updates(session_map, actions, update);
  local_enforcer->execute_actions(session_map, actions, update);

  // Assert that we exited final action state
  assert_session_is_in_final_state(
      session_map, IMSI1, SESSION_ID_1, credit_key, false);
}

// If a credit is a final credit, we should not send any updates unless a ReAuth
// is pending
TEST_F(LocalEnforcerTest, test_final_unit_action_no_update) {
  CreateSessionResponse response;
  test_cfg_.common_context.mutable_sid()->set_id(IMSI1);
  create_credit_update_response(
      IMSI1, SESSION_ID_1, 1, 1024, ChargingCredit_FinalAction_RESTRICT_ACCESS,
      "", "restrict_rule", response.mutable_credits()->Add());

  auto static_rule = response.mutable_static_rules()->Add();
  static_rule->set_rule_id("static_rule1");

  insert_static_rule(1, "", "static_rule1");
  insert_static_rule(1, "", "restrict_rule");
  auto& ip_addr     = test_cfg_.common_context.ue_ipv4();
  auto& ipv6_addr   = test_cfg_.common_context.ue_ipv6();
  auto teids        = teids1;
  const auto msisdn = test_cfg_.common_context.msisdn();
  // The activation for the static rule (static_rule1) and no dynamic
  EXPECT_CALL(
      *pipelined_client, activate_flows_for_rules(
                             IMSI1, ip_addr, ipv6_addr, CheckTeids(teids),
                             msisdn, testing::_, CheckRuleCount(1), testing::_))
      .Times(1);

  local_enforcer->init_session(
      session_map, IMSI1, SESSION_ID_1, test_cfg_, response);
  local_enforcer->update_tunnel_ids(
      session_map,
      create_update_tunnel_ids_request(IMSI1, BEARER_ID_1, teids1));

  // Insert record just over the reporting threshold and aggregate over them
  RuleRecordTable table;
  auto record_list = table.mutable_records();
  create_rule_record(
      IMSI1, teids.agw_teid(), "static_rule1", 1023, 0, record_list->Add());
  auto update = SessionStore::get_default_session_update(session_map);
  local_enforcer->aggregate_records(session_map, table, update);

  std::vector<std::unique_ptr<ServiceAction>> actions;
  auto usage_updates =
      local_enforcer->collect_updates(session_map, actions, update);
  // No update should be seen
  EXPECT_EQ(usage_updates.updates_size(), 0);

  // Now exceed 100% quota
  table.Clear();
  record_list = table.mutable_records();
  create_rule_record(
      IMSI1, teids1.agw_teid(), "static_rule1", 1024, 0, record_list->Add());
  update = SessionStore::get_default_session_update(session_map);
  local_enforcer->aggregate_records(session_map, table, update);

  // Collect actions and verify that restrict action is in the list
  actions.clear();
  usage_updates = local_enforcer->collect_updates(session_map, actions, update);
  EXPECT_EQ(actions.size(), 1);
  EXPECT_EQ(actions[0]->get_type(), RESTRICT_ACCESS);
  EXPECT_EQ(
      actions[0]->get_gy_rules_to_install()[0].rule.id(), "restrict_rule");

  EXPECT_CALL(
      *pipelined_client,
      add_gy_final_action_flow(
          IMSI1, ip_addr, ipv6_addr, CheckTeids(teids),
          test_cfg_.common_context.msisdn(), CheckRuleCount(1)))
      .Times(1);
  // Execute actions and asset final action state
  local_enforcer->execute_actions(session_map, actions, update);

  const CreditKey& credit_key(1);
  assert_session_is_in_final_state(
      session_map, IMSI1, SESSION_ID_1, credit_key, true);
}

// Test how we handle UpdateSessionResponse with dynamic rule modification
// 1. Start with a case where dynamic rule X is installed.
// 2. Receive an UpdateSessionResponse with removal instruction for the
// rule as well as an install for a dynamic rule with the same name X.
// 3. Assert we have rule X installed with modified entry
TEST_F(LocalEnforcerTest, test_rar_dynamic_rule_modification) {
  CreateSessionResponse response;
  test_cfg_.common_context.mutable_sid()->set_id(IMSI1);
  create_credit_update_response(
      IMSI1, SESSION_ID_1, 1, 1024, response.mutable_credits()->Add());
  auto dynamic_rule =
      response.mutable_dynamic_rules()->Add()->mutable_policy_rule();
  dynamic_rule->set_id("d-rule1");
  dynamic_rule->set_rating_group(1);
  dynamic_rule->set_tracking_type(PolicyRule::ONLY_OCS);
  // The activation for no static rules and 1 dynamic rule (d-rule1)
  EXPECT_CALL(
      *pipelined_client, activate_flows_for_rules(
                             IMSI1, testing::_, testing::_, testing::_,
                             test_cfg_.common_context.msisdn(), testing::_,
                             CheckRuleCount(1), testing::_))
      .Times(1);

  local_enforcer->init_session(
      session_map, IMSI1, SESSION_ID_1, test_cfg_, response);
  local_enforcer->update_tunnel_ids(
      session_map,
      create_update_tunnel_ids_request(IMSI1, BEARER_ID_1, teids1));

  session_store->create_sessions(IMSI1, std::move(session_map[IMSI1]));

  // session store create session
  session_map      = session_store->read_sessions(SessionRead{IMSI1});
  auto session_ucs = SessionStore::get_default_session_update(session_map);

  PolicyReAuthRequest rar;
  PolicyReAuthAnswer raa;
  rar.set_imsi(IMSI1);
  rar.set_session_id(SESSION_ID_1);
  rar.add_rules_to_remove("d-rule1");
  auto dynamic_install = rar.mutable_dynamic_rules_to_install()->Add();
  dynamic_install->mutable_policy_rule()->set_id("d-rule1");
  dynamic_install->mutable_policy_rule()->set_rating_group(2);
  dynamic_install->mutable_policy_rule()->set_tracking_type(
      PolicyRule::ONLY_OCS);
  {
    InSequence s;

    EXPECT_CALL(
        *pipelined_client, deactivate_flows_for_rules(
                               IMSI1, testing::_, testing::_, testing::_,
                               CheckRuleCount(1), testing::_))
        .Times(1);
    EXPECT_CALL(
        *pipelined_client, activate_flows_for_rules(
                               IMSI1, testing::_, testing::_, testing::_,
                               test_cfg_.common_context.msisdn(), testing::_,
                               CheckRuleCount(1), testing::_))
        .Times(1);
  }
  local_enforcer->init_policy_reauth(session_map, rar, raa, session_ucs);
  auto& dynamic_rules = session_map[IMSI1][0]->get_dynamic_rules();
  PolicyRule policy_out;
  EXPECT_TRUE(dynamic_rules.get_rule("d-rule1", &policy_out));
  EXPECT_EQ(2, policy_out.rating_group());
  EXPECT_TRUE(session_store->update_sessions(session_ucs));
}

// Test the case where PipelineD sends a data usage report for a session that
// does not exist anymore. We expect SessionD to send a deactivate flows request
// to PipelineD.
TEST_F(LocalEnforcerTest, test_dead_session_in_usage_report) {
  uint32_t teid = 32;
  Teids expected_teids;
  expected_teids.set_agw_teid(teid);
  expected_teids.set_enb_teid(0);  // we don't care about this one
  std::vector<Teids> expected_teids_vec{expected_teids};
  // no sessions exist at this point
  // We expect to empty calls for both Gx + Gy
  EXPECT_CALL(
      *pipelined_client,
      deactivate_flows_for_rules_for_termination(
          IMSI1, testing::_, testing::_, CheckTeidVector(expected_teids_vec),
          RequestOriginType::WILDCARD))
      .Times(1);

  RuleRecordTable table;
  create_rule_record(
      IMSI1, teid, "rule1", 16, 32, table.mutable_records()->Add());
  auto update = SessionStore::get_default_session_update(session_map);
  local_enforcer->aggregate_records(session_map, table, update);
}

int main(int argc, char** argv) {
  ::testing::InitGoogleTest(&argc, argv);
  FLAGS_logtostderr = 1;
  FLAGS_v           = 10;
  return RUN_ALL_TESTS();
}

}  // namespace magma

/**
 * Copyright (c) 2016-present, Facebook, Inc.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */
#include <chrono>
#include <future>
#include <memory>
#include <string.h>
#include <time.h>

#include <folly/io/async/EventBaseManager.h>
#include <gtest/gtest.h>
#include <lte/protos/session_manager.grpc.pb.h>

#include "LocalEnforcer.h"
#include "MagmaService.h"
#include "ProtobufCreators.h"
#include "ServiceRegistrySingleton.h"
#include "SessionStore.h"
#include "SessiondMocks.h"
#include "magma_logging.h"

#define SECONDS_A_DAY 86400

using grpc::ServerContext;
using grpc::Status;
using ::testing::Test;


namespace magma {

class LocalEnforcerTest : public ::testing::Test {
protected:
protected:
  void SetUpWithMConfig(magma::mconfig::SessionD mconfig) {
    reporter = std::make_shared<MockSessionReporter>();
    rule_store = std::make_shared<StaticRuleStore>();
    session_store = std::make_shared<SessionStore>(rule_store);
    pipelined_client = std::make_shared<MockPipelinedClient>();
    directoryd_client = std::make_shared<MockDirectorydClient>();
    spgw_client = std::make_shared<MockSpgwServiceClient>();
    aaa_client = std::make_shared<MockAAAClient>();
    local_enforcer = std::make_unique<LocalEnforcer>(
        reporter, rule_store, *session_store, pipelined_client,
        directoryd_client, MockEventdClient::getInstance(), spgw_client,
        aaa_client, 0, 0, mconfig);
    evb = folly::EventBaseManager::get()->getEventBase();
    local_enforcer->attachEventBase(evb);
    session_map = SessionMap{};
  }

  virtual void SetUp() {}

  virtual void TearDown() { folly::EventBaseManager::get()->clearEventBase(); }

  void run_evb() {
    evb->runAfterDelay([this]() { local_enforcer->stop(); }, 100);
    local_enforcer->start();
  }

  magma::mconfig::SessionD get_mconfig_gx_rule_wallet_exhaust() {
    magma::mconfig::SessionD mconfig;
    mconfig.set_log_level(magma::orc8r::LogLevel::INFO);
    mconfig.set_relay_enabled(false);
    auto wallet_config = mconfig.mutable_wallet_exhaust_detection();
    wallet_config->set_terminate_on_exhaust(true);
    wallet_config->set_method(magma::mconfig::WalletExhaustDetection_Method_GxTrackedRules);
    return mconfig;
  }

  void insert_static_rule(uint32_t rating_group, const std::string &m_key,
                          const std::string &rule_id) {
    PolicyRule rule;
    create_policy_rule(rule_id, m_key, rating_group, &rule);
    rule_store->insert_rule(rule);
  }

protected:
  std::shared_ptr<MockSessionReporter> reporter;
  std::shared_ptr<StaticRuleStore> rule_store;
  std::shared_ptr<SessionStore> session_store;
  std::unique_ptr<LocalEnforcer> local_enforcer;
  std::shared_ptr<MockPipelinedClient> pipelined_client;
  std::shared_ptr<MockDirectorydClient> directoryd_client;
  std::shared_ptr<MockSpgwServiceClient> spgw_client;
  std::shared_ptr<MockAAAClient> aaa_client;
  SessionMap session_map;
  folly::EventBase *evb;
};

MATCHER_P2(CheckQuotaUpdateState, size, expected_states, "") {
  auto updates = static_cast<const std::vector<SubscriberQuotaUpdate>>(arg);
  int updates_size = updates.size();
  if (updates_size != size) {
    return false;
  }
  for (int i = 0; i < updates_size; i++) {
    if (updates[i].update_type() != expected_states[i]) {
      return false;
    }
  }
  return true;
}

TEST_F(LocalEnforcerTest, test_cwf_quota_exhaustion_on_init_has_quota) {
  SetUpWithMConfig(get_mconfig_gx_rule_wallet_exhaust());
  insert_static_rule(0, "m1", "static_1");

  std::vector<std::string> static_rules{"static_1"};
  SessionConfig test_cwf_cfg;
  test_cwf_cfg.rat_type = RATType::TGPP_WLAN;
  CreateSessionResponse response;
  create_session_create_response("IMSI1", "m1", static_rules, &response);

  StaticRuleInstall static_rule_install;
  static_rule_install.set_rule_id("static_1");
  auto res_rules_to_install = response.mutable_static_rules()->Add();
  res_rules_to_install->CopyFrom(static_rule_install);

  std::vector<SubscriberQuotaUpdate_Type> expected_states{
      SubscriberQuotaUpdate_Type_VALID_QUOTA};
  EXPECT_CALL(*pipelined_client, update_subscriber_quota_state(
                                     CheckQuotaUpdateState(1, expected_states)))
      .Times(1)
      .WillOnce(testing::Return(true));
  local_enforcer->init_session_credit(session_map, "IMSI1", "1234",
                                      test_cwf_cfg, response);
}

TEST_F(LocalEnforcerTest, test_cwf_quota_exhaustion_on_init_no_quota) {
  SetUpWithMConfig(get_mconfig_gx_rule_wallet_exhaust());
  insert_static_rule(1, "m1", "static_1");

  std::vector<std::string> static_rules{}; // no rule installs
  SessionConfig test_cwf_cfg;
  test_cwf_cfg.rat_type = RATType::TGPP_WLAN;
  CreateSessionResponse response;
  create_session_create_response("IMSI1", "m1", static_rules, &response);

  std::vector<SubscriberQuotaUpdate_Type> expected_states{
      SubscriberQuotaUpdate_Type_NO_QUOTA};
  EXPECT_CALL(*pipelined_client, update_subscriber_quota_state(
                                     CheckQuotaUpdateState(1, expected_states)))
      .Times(1)
      .WillOnce(testing::Return(true));
  local_enforcer->init_session_credit(session_map, "IMSI1", "1234",
                                      test_cwf_cfg, response);
}

TEST_F(LocalEnforcerTest, test_cwf_quota_exhaustion_on_rar) {
  SetUpWithMConfig(get_mconfig_gx_rule_wallet_exhaust());
  // setup : successful session creation with valid monitoring quota
  insert_static_rule(0, "m1", "static_1");

  std::vector<std::string> static_rules{"static_1"};
  SessionConfig test_cwf_cfg;
  test_cwf_cfg.rat_type = RATType::TGPP_WLAN;
  CreateSessionResponse response;
  create_session_create_response("IMSI1", "m1", static_rules, &response);
  local_enforcer->init_session_credit(session_map, "IMSI1", "1234",
                                      test_cwf_cfg, response);

  // send a policy reauth request with rule removals for "static_1" to indicate
  // total monitoring quota exhaustion
  PolicyReAuthRequest request;
  request.set_session_id("");
  request.set_imsi("IMSI1");
  request.add_rules_to_remove("static_1");

  std::vector<SubscriberQuotaUpdate_Type> expected_states{
      SubscriberQuotaUpdate_Type_TERMINATE};
  EXPECT_CALL(*pipelined_client, update_subscriber_quota_state(
                                     CheckQuotaUpdateState(1, expected_states)))
      .Times(1)
      .WillOnce(testing::Return(true));

  PolicyReAuthAnswer answer;
  auto update = SessionStore::get_default_session_update(session_map);
  local_enforcer->init_policy_reauth(session_map, request, answer, update);
}

TEST_F(LocalEnforcerTest, test_cwf_quota_exhaustion_on_update) {
  SetUpWithMConfig(get_mconfig_gx_rule_wallet_exhaust());
  // setup : successful session creation with valid monitoring quota
  insert_static_rule(0, "m1", "static_1");
  insert_static_rule(0, "m1", "static_2");

  std::vector<std::string> static_rules{"static_1", "static_2"};
  SessionConfig test_cwf_cfg;
  test_cwf_cfg.rat_type = RATType::TGPP_WLAN;
  CreateSessionResponse response;
  create_session_create_response("IMSI1", "m1", static_rules, &response);
  local_enforcer->init_session_credit(session_map, "IMSI1", "1234",
                                      test_cwf_cfg, response);

  // remove only static_2, should not change anything in terms of quota since
  // static_1 is still active
  UpdateSessionResponse update_response;
  auto monitor = update_response.mutable_usage_monitor_responses()->Add();
  create_monitor_update_response("IMSI1", "m1", MonitoringLevel::PCC_RULE_LEVEL,
                                 2048, monitor);
  monitor->add_rules_to_remove("static_2");
  auto update = SessionStore::get_default_session_update(session_map);
  local_enforcer->update_session_credits_and_rules(session_map, update_response,
                                                   update);

  // send an update response with rule removals for "static_1" to indicate
  // total monitoring quota exhaustion
  update_response.clear_usage_monitor_responses();
  monitor = update_response.mutable_usage_monitor_responses()->Add();
  create_monitor_update_response("IMSI1", "m1", MonitoringLevel::PCC_RULE_LEVEL,
                                 0, monitor);
  monitor->add_rules_to_remove("static_1");

  std::vector<SubscriberQuotaUpdate_Type> expected_states = {
      SubscriberQuotaUpdate_Type_TERMINATE};
  EXPECT_CALL(*pipelined_client, update_subscriber_quota_state(
                                     CheckQuotaUpdateState(1, expected_states)))
      .Times(1)
      .WillOnce(testing::Return(true));

  local_enforcer->update_session_credits_and_rules(session_map, update_response,
                                                   update);
}

int main(int argc, char **argv) {
  ::testing::InitGoogleTest(&argc, argv);
  FLAGS_logtostderr = 1;
  FLAGS_v = 10;
  return RUN_ALL_TESTS();
}

} // namespace magma

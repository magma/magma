/**
 * Copyright (c) 2016-present, Facebook, Inc.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */
#include <chrono>
#include <memory>
#include <string.h>
#include <time.h>
#include <future>

#include <folly/io/async/EventBaseManager.h>
#include <gtest/gtest.h>
#include <lte/protos/session_manager.grpc.pb.h>

#include "MagmaService.h"
#include "ProtobufCreators.h"
#include "ServiceRegistrySingleton.h"
#include "SessiondMocks.h"
#include "LocalEnforcer.h"
#include "magma_logging.h"

#define SECONDS_A_DAY 86400

using grpc::ServerContext;
using grpc::Status;
using ::testing::Test;

namespace magma {

const SessionState::Config test_cfg = {.ue_ipv4 = "127.0.0.1",
                                       .spgw_ipv4 = "128.0.0.1",
                                       .msisdn = "",
                                       .apn = "IMS"};

class LocalEnforcerTest : public ::testing::Test {
 protected:
 protected:
  virtual void SetUp()
  {
    reporter = std::make_shared<MockSessionReporter>();
    rule_store = std::make_shared<StaticRuleStore>();
    pipelined_client = std::make_shared<MockPipelinedClient>();
    directoryd_client = std::make_shared<MockDirectorydClient>();
    spgw_client = std::make_shared<MockSpgwServiceClient>();
    aaa_client = std::make_shared<MockAAAClient>();
    local_enforcer = std::make_unique<LocalEnforcer>(
      reporter, rule_store, pipelined_client, directoryd_client, spgw_client,
      aaa_client, 0);
    evb = folly::EventBaseManager::get()->getEventBase();
    local_enforcer->attachEventBase(evb);
  }

  virtual void TearDown()
  {
    folly::EventBaseManager::get()->clearEventBase();
  }

  void run_evb()
  {
    evb->runAfterDelay([this]() { local_enforcer->stop(); }, 100);
    local_enforcer->start();
  }

  void insert_static_rule(
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
    rule_store->insert_rule(rule);
  }

  void assert_charging_credit(
    const std::string &imsi,
    Bucket bucket,
    const std::vector<std::pair<uint32_t, uint64_t>> &volumes)
  {
    for (auto &volume_pair : volumes) {
      auto volume_out =
        local_enforcer->get_charging_credit(imsi, volume_pair.first, bucket);
      EXPECT_EQ(volume_out, volume_pair.second);
    }
  }

  void assert_monitor_credit(
    const std::string &imsi,
    Bucket bucket,
    const std::vector<std::pair<std::string, uint64_t>> &volumes)
  {
    for (auto &volume_pair : volumes) {
      auto volume_out =
        local_enforcer->get_monitor_credit(imsi, volume_pair.first, bucket);
      EXPECT_EQ(volume_out, volume_pair.second);
    }
  }

 protected:
  std::shared_ptr<MockSessionReporter> reporter;
  std::shared_ptr<StaticRuleStore> rule_store;
  std::unique_ptr<LocalEnforcer> local_enforcer;
  std::shared_ptr<MockPipelinedClient> pipelined_client;
  std::shared_ptr<MockDirectorydClient> directoryd_client;
  std::shared_ptr<MockSpgwServiceClient> spgw_client;
  std::shared_ptr<MockAAAClient> aaa_client;
  folly::EventBase *evb;
};

MATCHER_P(CheckCount, count, "")
{
  return arg.size() == count;
}

MATCHER_P2(CheckActivateFlows, imsi, rule_count, "")
{
  auto request = static_cast<const ActivateFlowsRequest *>(arg);
  return request->sid().id() == imsi && request->rule_ids_size() == rule_count;
}

MATCHER_P4(CheckSessionInfos, imsi_list, ip_address_list, static_rule_lists,
           dynamic_rule_ids_lists, "")
{
  auto infos = static_cast<const std::vector<SessionState::SessionInfo>>(arg);

  if (infos.size() != imsi_list.size())
    return false;

  for (int i = 0; i < infos.size(); i++)
  {
    if (infos[i].imsi != imsi_list[i])
      return false;
    if (infos[i].ip_addr != ip_address_list[i])
      return false;
    if (infos[i].static_rules.size() != static_rule_lists[i].size())
      return false;
    if (infos[i].dynamic_rules.size() != dynamic_rule_ids_lists[i].size())
      return false;
    for (int r_index = 0; i < infos[i].static_rules.size(); i++)
    {
      if (infos[i].static_rules[r_index] != static_rule_lists[i][r_index])
        return false;
    }
    for (int r_index = 0; i < infos[i].dynamic_rules.size(); i++)
    {
      if (infos[i].dynamic_rules[r_index].id() !=
          dynamic_rule_ids_lists[i][r_index])
        return false;
    }
  }
  return true;
}

TEST_F(LocalEnforcerTest, test_init_cwf_session_credit)
{
  insert_static_rule(1, "", "rule1");

  CreateSessionResponse response;
  auto credits = response.mutable_credits();
  create_credit_update_response("IMSI1", 1, 1024, credits->Add());

  EXPECT_CALL(
    *pipelined_client,
    add_ue_mac_flow(
      testing::_, testing::_))
    .Times(1)
    .WillOnce(testing::Return(true));

  EXPECT_CALL(
    *pipelined_client,
    activate_flows_for_rules(
      testing::_, testing::_, CheckCount(0), CheckCount(0)))
    .Times(1)
    .WillOnce(testing::Return(true));

  SessionState::Config test_cwf_cfg;
  test_cwf_cfg.rat_type = RATType::TGPP_WLAN;
  test_cwf_cfg.mac_addr = "00:00:00:00:00:00";
  test_cwf_cfg.radius_session_id = "1234567";

  local_enforcer->init_session_credit("IMSI1", "1234", test_cwf_cfg, response);

  EXPECT_EQ(
    local_enforcer->get_charging_credit("IMSI1", 1, ALLOWED_TOTAL), 1024);
}

TEST_F(LocalEnforcerTest, test_init_session_credit)
{
  insert_static_rule(1, "", "rule1");

  CreateSessionResponse response;
  auto credits = response.mutable_credits();
  create_credit_update_response("IMSI1", 1, 1024, credits->Add());

  EXPECT_CALL(
    *pipelined_client,
    activate_flows_for_rules(
      testing::_, testing::_, CheckCount(0), CheckCount(0)))
    .Times(1)
    .WillOnce(testing::Return(true));
  local_enforcer->init_session_credit("IMSI1", "1234", test_cfg, response);

  EXPECT_EQ(
    local_enforcer->get_charging_credit("IMSI1", 1, ALLOWED_TOTAL), 1024);
}

TEST_F(LocalEnforcerTest, test_single_record)
{
  // insert initial session credit
  CreateSessionResponse response;
  create_credit_update_response(
    "IMSI1", 1, 1024, response.mutable_credits()->Add());
  local_enforcer->init_session_credit("IMSI1", "1234", test_cfg, response);

  insert_static_rule(1, "", "rule1");
  RuleRecordTable table;
  auto record_list = table.mutable_records();
  create_rule_record("IMSI1", "rule1", 16, 32, record_list->Add());

  local_enforcer->aggregate_records(table);

  EXPECT_EQ(local_enforcer->get_charging_credit("IMSI1", 1, USED_RX), 16);
  EXPECT_EQ(local_enforcer->get_charging_credit("IMSI1", 1, USED_TX), 32);
  EXPECT_EQ(
    local_enforcer->get_charging_credit("IMSI1", 1, ALLOWED_TOTAL), 1024);
}

TEST_F(LocalEnforcerTest, test_aggregate_records)
{
  CreateSessionResponse response;
  create_credit_update_response(
    "IMSI1", 1, 1024, response.mutable_credits()->Add());
  create_credit_update_response(
    "IMSI1", 2, 1024, response.mutable_credits()->Add());
  local_enforcer->init_session_credit("IMSI1", "1234", test_cfg, response);

  insert_static_rule(1, "", "rule1");
  insert_static_rule(1, "", "rule2");
  insert_static_rule(2, "", "rule3");
  RuleRecordTable table;
  auto record_list = table.mutable_records();
  create_rule_record("IMSI1", "rule1", 10, 20, record_list->Add());
  create_rule_record("IMSI1", "rule2", 5, 15, record_list->Add());
  create_rule_record("IMSI1", "rule3", 100, 150, record_list->Add());

  local_enforcer->aggregate_records(table);

  EXPECT_EQ(local_enforcer->get_charging_credit("IMSI1", 1, USED_RX), 15);
  EXPECT_EQ(local_enforcer->get_charging_credit("IMSI1", 1, USED_TX), 35);
  EXPECT_EQ(local_enforcer->get_charging_credit("IMSI1", 2, USED_RX), 100);
  EXPECT_EQ(local_enforcer->get_charging_credit("IMSI1", 2, USED_TX), 150);
}

TEST_F(LocalEnforcerTest, test_aggregate_records_for_termination)
{
  CreateSessionResponse response;
  create_credit_update_response(
    "IMSI1", 1, 1024, response.mutable_credits()->Add());
  create_credit_update_response(
    "IMSI1", 2, 1024, response.mutable_credits()->Add());
  local_enforcer->init_session_credit("IMSI1", "1234", test_cfg, response);

  insert_static_rule(1, "", "rule1");
  insert_static_rule(1, "", "rule2");
  insert_static_rule(2, "", "rule3");

  std::promise<void> termination_promise;
  auto future = termination_promise.get_future();
  local_enforcer->terminate_subscriber(
    "IMSI1", "IMS", [&termination_promise](SessionTerminateRequest req) {
      termination_promise.set_value();

      EXPECT_EQ(req.credit_usages_size(), 2);
      for (const auto &usage : req.credit_usages()) {
        EXPECT_EQ(usage.type(), CreditUsage::TERMINATED);
      }
    });

  RuleRecordTable table;
  auto record_list = table.mutable_records();
  create_rule_record("IMSI1", "rule1", 10, 20, record_list->Add());
  create_rule_record("IMSI1", "rule2", 5, 15, record_list->Add());
  create_rule_record("IMSI1", "rule3", 100, 150, record_list->Add());

  local_enforcer->aggregate_records(table);

  // Termination should not have been completed since we are still aggregating
  // the records.
  auto status = future.wait_for(std::chrono::seconds(0));
  EXPECT_EQ(status, std::future_status::timeout);

  RuleRecordTable empty_table;

  local_enforcer->aggregate_records(empty_table);

  status = future.wait_for(std::chrono::seconds(0));
  EXPECT_EQ(status, std::future_status::ready);
}

TEST_F(LocalEnforcerTest, test_collect_updates)
{
  CreateSessionResponse response;
  create_credit_update_response(
    "IMSI1", 1, 3072, response.mutable_credits()->Add());
  local_enforcer->init_session_credit("IMSI1", "1234", test_cfg, response);
  insert_static_rule(1, "", "rule1");

  auto empty_update = local_enforcer->collect_updates();
  EXPECT_EQ(empty_update.updates_size(), 0);

  RuleRecordTable table;
  auto record_list = table.mutable_records();
  create_rule_record("IMSI1", "rule1", 1024, 2048, record_list->Add());

  local_enforcer->aggregate_records(table);
  auto session_update = local_enforcer->collect_updates();
  EXPECT_EQ(session_update.updates_size(), 1);
  EXPECT_EQ(
    local_enforcer->get_charging_credit("IMSI1", 1, REPORTING_RX), 1024);
  EXPECT_EQ(
    local_enforcer->get_charging_credit("IMSI1", 1, REPORTING_TX), 2048);
}

TEST_F(LocalEnforcerTest, test_update_session_credit)
{
  insert_static_rule(1, "", "rule1");

  CreateSessionResponse response;
  auto credits = response.mutable_credits();
  create_credit_update_response("IMSI1", 1, 2048, credits->Add());
  local_enforcer->init_session_credit("IMSI1", "1234", test_cfg, response);

  EXPECT_EQ(
    local_enforcer->get_charging_credit("IMSI1", 1, ALLOWED_TOTAL), 2048);

  insert_static_rule(1, "1", "rule1");

  RuleRecordTable table;
  auto record_list = table.mutable_records();
  create_rule_record("IMSI1", "rule1", 1024, 1024, record_list->Add());
  local_enforcer->aggregate_records(table);

  UpdateSessionResponse update_response;
  auto credit_updates_response = update_response.mutable_responses();
  create_credit_update_response("IMSI1", 1, 24, credit_updates_response->Add());

  std::vector<EventTrigger> event_triggers {
    EventTrigger::TAI_CHANGE, EventTrigger::REVALIDATION_TIMEOUT};
  auto monitor_updates_response =
    update_response.mutable_usage_monitor_responses();
  create_monitor_update_response(
    "IMSI",
    "1",
    MonitoringLevel::PCC_RULE_LEVEL,
    2048,
    event_triggers,
    time(NULL),
    monitor_updates_response->Add());
  EXPECT_CALL(*reporter, report_updates(_, _)).Times(1);
  local_enforcer->update_session_credit(update_response);
  EXPECT_EQ(
    local_enforcer->get_charging_credit("IMSI1", 1, ALLOWED_TOTAL), 2072);
}

TEST_F(LocalEnforcerTest, test_terminate_credit)
{
  CreateSessionResponse response;
  create_credit_update_response(
    "IMSI1", 1, 1024, response.mutable_credits()->Add());
  create_credit_update_response(
    "IMSI1", 2, 2048, response.mutable_credits()->Add());
  CreateSessionResponse response2;
  create_credit_update_response(
    "IMSI2", 1, 4096, response2.mutable_credits()->Add());
  local_enforcer->init_session_credit("IMSI1", "1234", test_cfg, response);
  local_enforcer->init_session_credit("IMSI2", "4321", test_cfg, response2);

  std::promise<void> termination_promise;
  local_enforcer->terminate_subscriber(
    "IMSI1", "IMS", [&termination_promise](SessionTerminateRequest req) {
      termination_promise.set_value();

      EXPECT_EQ(req.credit_usages_size(), 2);
      for (const auto &usage : req.credit_usages()) {
        EXPECT_EQ(usage.type(), CreditUsage::TERMINATED);
      }
    });
  run_evb();
  auto status =
    termination_promise.get_future().wait_for(std::chrono::seconds(0));
  EXPECT_EQ(status, std::future_status::ready);

  // No longer in system
  EXPECT_EQ(local_enforcer->get_charging_credit("IMSI1", 1, ALLOWED_TOTAL), 0);
  EXPECT_EQ(local_enforcer->get_charging_credit("IMSI1", 2, ALLOWED_TOTAL), 0);
}

TEST_F(LocalEnforcerTest, test_terminate_credit_during_reporting)
{
  CreateSessionResponse response;
  create_credit_update_response(
    "IMSI1", 1, 3072, response.mutable_credits()->Add());
  create_credit_update_response(
    "IMSI1", 2, 2048, response.mutable_credits()->Add());
  local_enforcer->init_session_credit("IMSI1", "1234", test_cfg, response);
  insert_static_rule(1, "", "rule1");
  insert_static_rule(2, "", "rule2");

  // Insert record for key 1
  RuleRecordTable table;
  auto record_list = table.mutable_records();
  create_rule_record("IMSI1", "rule1", 1024, 2048, record_list->Add());
  local_enforcer->aggregate_records(table);

  // Collect updates to put key 1 into reporting state
  auto usage_updates = local_enforcer->collect_updates();
  EXPECT_EQ(
    local_enforcer->get_charging_credit("IMSI1", 1, REPORTING_RX), 1024);

  // Collecting terminations should key 1 anyways during reporting
  std::promise<void> termination_promise;
  local_enforcer->terminate_subscriber(
    "IMSI1", "IMS", [&termination_promise](SessionTerminateRequest term_req) {
      termination_promise.set_value();

      EXPECT_EQ(term_req.credit_usages_size(), 2);
    });
  run_evb();
  auto status =
    termination_promise.get_future().wait_for(std::chrono::seconds(0));
  EXPECT_EQ(status, std::future_status::ready);
}

MATCHER_P2(CheckDeactivateFlows, imsi, rule_count, "")
{
  auto request = static_cast<const DeactivateFlowsRequest *>(arg);
  return request->sid().id() == imsi && request->rule_ids_size() == rule_count;
}

TEST_F(LocalEnforcerTest, test_final_unit_handling)
{
  CreateSessionResponse response;
  create_credit_update_response(
    "IMSI1", 1, true, 1024, response.mutable_credits()->Add());
  local_enforcer->init_session_credit("IMSI1", "1234", test_cfg, response);
  insert_static_rule(1, "", "rule1");
  insert_static_rule(1, "", "rule2");

  // Insert record for key 1
  RuleRecordTable table;
  auto record_list = table.mutable_records();
  create_rule_record("IMSI1", "rule1", 1024, 2048, record_list->Add());
  create_rule_record("IMSI1", "rule2", 1024, 2048, record_list->Add());
  local_enforcer->aggregate_records(table);

  EXPECT_CALL(
    *pipelined_client,
    deactivate_flows_for_rules(testing::_, testing::_, testing::_))
    .Times(1)
    .WillOnce(testing::Return(true));
  // call collect_updates to trigger actions
  auto usage_updates = local_enforcer->collect_updates();
}

TEST_F(LocalEnforcerTest, test_cwf_final_unit_handling)
{
  CreateSessionResponse response;
  create_credit_update_response(
    "IMSI1", 1, true, 1024, response.mutable_credits()->Add());

  SessionState::Config test_cwf_cfg;
  test_cwf_cfg.rat_type = RATType::TGPP_WLAN;
  test_cwf_cfg.mac_addr = "00:00:00:00:00:00";
  test_cwf_cfg.radius_session_id = "1234567";

  local_enforcer->init_session_credit("IMSI1", "1234", test_cwf_cfg, response);
  insert_static_rule(1, "", "rule1");
  insert_static_rule(1, "", "rule2");

  // Insert record for key 1
  RuleRecordTable table;
  auto record_list = table.mutable_records();
  create_rule_record("IMSI1", "rule1", 1024, 2048, record_list->Add());
  create_rule_record("IMSI1", "rule2", 1024, 2048, record_list->Add());
  local_enforcer->aggregate_records(table);

  EXPECT_CALL(
    *pipelined_client,
    deactivate_flows_for_rules(testing::_, testing::_, testing::_))
    .Times(1)
    .WillOnce(testing::Return(true));

  EXPECT_CALL(
    *aaa_client,
    terminate_session(testing::_, testing::_))
    .Times(1)
    .WillOnce(testing::Return(true));
  // call collect_updates to trigger actions
  auto usage_updates = local_enforcer->collect_updates();
}

TEST_F(LocalEnforcerTest, test_all)
{
  // insert key rule mapping
  insert_static_rule(1, "", "rule1");
  insert_static_rule(1, "", "rule2");
  insert_static_rule(2, "", "rule3");

  // insert initial session credit
  CreateSessionResponse response;
  create_credit_update_response(
    "IMSI1", 1, 1024, response.mutable_credits()->Add());
  local_enforcer->init_session_credit("IMSI1", "1234", test_cfg, response);
  CreateSessionResponse response2;
  create_credit_update_response(
    "IMSI2", 2, 2048, response2.mutable_credits()->Add());
  local_enforcer->init_session_credit("IMSI2", "4321", test_cfg, response2);

  EXPECT_EQ(
    local_enforcer->get_charging_credit("IMSI1", 1, ALLOWED_TOTAL), 1024);
  EXPECT_EQ(
    local_enforcer->get_charging_credit("IMSI2", 2, ALLOWED_TOTAL), 2048);

  // receive usages from pipelined
  RuleRecordTable table;
  auto record_list = table.mutable_records();
  create_rule_record("IMSI1", "rule1", 10, 20, record_list->Add());
  create_rule_record("IMSI1", "rule2", 5, 15, record_list->Add());
  create_rule_record("IMSI2", "rule3", 1024, 1024, record_list->Add());
  local_enforcer->aggregate_records(table);

  EXPECT_EQ(local_enforcer->get_charging_credit("IMSI1", 1, USED_RX), 15);
  EXPECT_EQ(local_enforcer->get_charging_credit("IMSI1", 1, USED_TX), 35);
  EXPECT_EQ(local_enforcer->get_charging_credit("IMSI2", 2, USED_RX), 1024);
  EXPECT_EQ(local_enforcer->get_charging_credit("IMSI2", 2, USED_TX), 1024);

  // Collect updates for reporting
  auto session_update = local_enforcer->collect_updates();
  EXPECT_EQ(session_update.updates_size(), 1);
  EXPECT_EQ(
    local_enforcer->get_charging_credit("IMSI2", 2, REPORTING_RX), 1024);
  EXPECT_EQ(
    local_enforcer->get_charging_credit("IMSI2", 2, REPORTING_TX), 1024);

  // Add updated credit from cloud
  UpdateSessionResponse update_response;
  auto updates = update_response.mutable_responses();
  create_credit_update_response("IMSI2", 2, 4096, updates->Add());
  local_enforcer->update_session_credit(update_response);

  EXPECT_EQ(
    local_enforcer->get_charging_credit("IMSI2", 2, ALLOWED_TOTAL), 6144);
  EXPECT_EQ(local_enforcer->get_charging_credit("IMSI2", 2, REPORTING_TX), 0);
  EXPECT_EQ(local_enforcer->get_charging_credit("IMSI2", 2, REPORTING_RX), 0);
  EXPECT_EQ(local_enforcer->get_charging_credit("IMSI2", 2, REPORTED_TX), 1024);
  EXPECT_EQ(local_enforcer->get_charging_credit("IMSI2", 2, REPORTED_RX), 1024);

  // Terminate IMSI1
  std::promise<void> termination_promise;
  local_enforcer->terminate_subscriber(
    "IMSI1", "IMS", [&termination_promise](SessionTerminateRequest term_req) {
      termination_promise.set_value();

      EXPECT_EQ(term_req.sid(), "IMSI1");
      EXPECT_EQ(term_req.credit_usages_size(), 1);
    });
  run_evb();
  auto status =
    termination_promise.get_future().wait_for(std::chrono::seconds(0));
  EXPECT_EQ(status, std::future_status::ready);
}

TEST_F(LocalEnforcerTest, test_re_auth)
{
  insert_static_rule(1, "", "rule1");
  CreateSessionResponse response;
  local_enforcer->init_session_credit("IMSI1", "1234", test_cfg, response);

  ChargingReAuthRequest reauth;
  reauth.set_sid("IMSI1");
  reauth.set_session_id("1234");
  reauth.set_charging_key(1);
  reauth.set_type(ChargingReAuthRequest::SINGLE_SERVICE);
  auto result = local_enforcer->init_charging_reauth(reauth);
  EXPECT_EQ(result, ChargingReAuthAnswer::UPDATE_INITIATED);

  auto update_req = local_enforcer->collect_updates();
  EXPECT_EQ(update_req.updates_size(), 1);
  EXPECT_EQ(update_req.updates(0).sid(), "IMSI1");
  EXPECT_EQ(update_req.updates(0).usage().type(), CreditUsage::REAUTH_REQUIRED);

  // Give credit after re-auth
  UpdateSessionResponse update_response;
  auto updates = update_response.mutable_responses();
  create_credit_update_response("IMSI1", 1, 4096, updates->Add());
  local_enforcer->update_session_credit(update_response);

  // when next update is collected, this should trigger an action to activate
  // the flow in pipelined
  EXPECT_CALL(
    *pipelined_client,
    activate_flows_for_rules(testing::_, testing::_, testing::_, testing::_))
    .Times(1)
    .WillOnce(testing::Return(true));
  local_enforcer->collect_updates();
}

TEST_F(LocalEnforcerTest, test_dynamic_rules)
{
  CreateSessionResponse response;
  create_credit_update_response(
    "IMSI1", 1, 1024, response.mutable_credits()->Add());
  auto dynamic_rule = response.mutable_dynamic_rules()->Add();
  auto policy_rule = dynamic_rule->mutable_policy_rule();
  policy_rule->set_id("rule1");
  policy_rule->set_rating_group(1);
  policy_rule->set_tracking_type(PolicyRule::ONLY_OCS);
  local_enforcer->init_session_credit("IMSI1", "1234", test_cfg, response);

  insert_static_rule(1, "", "rule2");
  RuleRecordTable table;
  auto record_list = table.mutable_records();
  create_rule_record("IMSI1", "rule1", 16, 32, record_list->Add());
  create_rule_record("IMSI1", "rule2", 8, 8, record_list->Add());

  local_enforcer->aggregate_records(table);

  EXPECT_EQ(local_enforcer->get_charging_credit("IMSI1", 1, USED_RX), 24);
  EXPECT_EQ(local_enforcer->get_charging_credit("IMSI1", 1, USED_TX), 40);
  EXPECT_EQ(
    local_enforcer->get_charging_credit("IMSI1", 1, ALLOWED_TOTAL), 1024);
}

TEST_F(LocalEnforcerTest, test_dynamic_rule_actions)
{
  CreateSessionResponse response;
  create_credit_update_response(
    "IMSI1", 1, true, 1024, response.mutable_credits()->Add());
  auto dynamic_rule = response.mutable_dynamic_rules()->Add();
  auto policy_rule = dynamic_rule->mutable_policy_rule();
  policy_rule->set_id("rule2");
  policy_rule->set_rating_group(1);
  policy_rule->set_tracking_type(PolicyRule::ONLY_OCS);
  insert_static_rule(1, "", "rule1");
  insert_static_rule(1, "", "rule3");

  EXPECT_CALL(
    *pipelined_client,
    activate_flows_for_rules(
      testing::_, testing::_, CheckCount(0), CheckCount(1)))
    .Times(1)
    .WillOnce(testing::Return(true));
  local_enforcer->init_session_credit("IMSI1", "1234", test_cfg, response);

  RuleRecordTable table;
  auto record_list = table.mutable_records();
  create_rule_record("IMSI1", "rule1", 1024, 2048, record_list->Add());
  create_rule_record("IMSI1", "rule2", 1024, 2048, record_list->Add());
  local_enforcer->aggregate_records(table);

  EXPECT_CALL(
    *pipelined_client,
    deactivate_flows_for_rules(testing::_, CheckCount(2), CheckCount(1)))
    .Times(1)
    .WillOnce(testing::Return(true));
  auto usage_updates = local_enforcer->collect_updates();
}

TEST_F(LocalEnforcerTest, test_installing_rules_with_activation_time)
{
  CreateSessionResponse response;
  create_credit_update_response(
    "IMSI1", 1, true, 1024, response.mutable_credits()->Add());

  // add a dynamic rule without activation time
  auto dynamic_rule = response.mutable_dynamic_rules()->Add();
  auto policy_rule = dynamic_rule->mutable_policy_rule();
  policy_rule->set_id("rule1");
  policy_rule->set_rating_group(1);
  policy_rule->set_tracking_type(PolicyRule::ONLY_OCS);

  // add a dynamic rule with activation time in the future
  dynamic_rule = response.mutable_dynamic_rules()->Add();
  policy_rule = dynamic_rule->mutable_policy_rule();
  auto activation_time = dynamic_rule->mutable_activation_time();
  activation_time->set_seconds(time(NULL) + SECONDS_A_DAY);
  policy_rule->set_id("rule2");
  policy_rule->set_rating_group(1);
  policy_rule->set_tracking_type(PolicyRule::ONLY_OCS);

  // add a dynamic rule with activation time in the past
  dynamic_rule = response.mutable_dynamic_rules()->Add();
  policy_rule = dynamic_rule->mutable_policy_rule();
  activation_time = dynamic_rule->mutable_activation_time();
  activation_time->set_seconds(time(NULL) - SECONDS_A_DAY);
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
  activation_time = static_rule->mutable_activation_time();
  activation_time->set_seconds(time(NULL) + SECONDS_A_DAY);
  static_rule->set_rule_id("rule5");

  // add a static rule with activation time in the past
  insert_static_rule(1, "", "rule6");
  static_rule = response.mutable_static_rules()->Add();
  activation_time = static_rule->mutable_activation_time();
  activation_time->set_seconds(time(NULL) - SECONDS_A_DAY);
  static_rule->set_rule_id("rule6");

  // expect calling activate_flows_for_rules for activating rules instantly
  // dynamic rules: rule1, rule3
  // static rules: rule4, rule6
  EXPECT_CALL(
    *pipelined_client,
    activate_flows_for_rules(
      testing::_, testing::_, CheckCount(2), CheckCount(2)))
    .Times(1)
    .WillOnce(testing::Return(true));
  // expect calling activate_flows_for_rules for activating a static rule later
  // static rules: rule5
  EXPECT_CALL(
    *pipelined_client,
    activate_flows_for_rules(
      testing::_, testing::_, CheckCount(1), CheckCount(0)))
    .Times(1)
    .WillOnce(testing::Return(true));
  // expect calling activate_flows_for_rules for activating a dynamic rule later
  // dynamic rules: rule2
  EXPECT_CALL(
    *pipelined_client,
    activate_flows_for_rules(
      testing::_, testing::_, CheckCount(0), CheckCount(1)))
    .Times(1)
    .WillOnce(testing::Return(true));
  local_enforcer->init_session_credit("IMSI1", "1234", test_cfg, response);
}

TEST_F(LocalEnforcerTest, test_usage_monitors)
{
  // insert key rule mapping
  insert_static_rule(1, "1", "both_rule");
  insert_static_rule(2, "", "ocs_rule");
  insert_static_rule(0, "3", "pcrf_only");
  insert_static_rule(0, "1", "pcrf_split"); // same mkey as both_rule
  // session level rule "4"

  // insert initial session credit
  CreateSessionResponse response;
  create_credit_update_response(
    "IMSI1", 1, 1024, response.mutable_credits()->Add());
  create_credit_update_response(
    "IMSI1", 2, 1024, response.mutable_credits()->Add());
  create_monitor_update_response(
    "IMSI1",
    "1",
    MonitoringLevel::PCC_RULE_LEVEL,
    1024,
    response.mutable_usage_monitors()->Add());
  create_monitor_update_response(
    "IMSI1",
    "3",
    MonitoringLevel::PCC_RULE_LEVEL,
    2048,
    response.mutable_usage_monitors()->Add());
  create_monitor_update_response(
    "IMSI1",
    "4",
    MonitoringLevel::SESSION_LEVEL,
    2128,
    response.mutable_usage_monitors()->Add());
  local_enforcer->init_session_credit("IMSI1", "1234", test_cfg, response);
  assert_charging_credit("IMSI1", ALLOWED_TOTAL, {{1, 1024}, {2, 1024}});
  assert_monitor_credit(
    "IMSI1", ALLOWED_TOTAL, {{"1", 1024}, {"3", 2048}, {"4", 2128}});

  // receive usages from pipelined
  RuleRecordTable table;
  auto record_list = table.mutable_records();
  create_rule_record("IMSI1", "both_rule", 10, 20, record_list->Add());
  create_rule_record("IMSI1", "ocs_rule", 5, 15, record_list->Add());
  create_rule_record("IMSI1", "pcrf_only", 1024, 1024, record_list->Add());
  create_rule_record("IMSI1", "pcrf_split", 10, 20, record_list->Add());
  local_enforcer->aggregate_records(table);

  assert_charging_credit("IMSI1", USED_RX, {{1, 10}, {2, 5}});
  assert_charging_credit("IMSI1", USED_TX, {{1, 20}, {2, 15}});
  assert_monitor_credit(
    "IMSI1", USED_RX, {{"1", 20}, {"3", 1024}, {"4", 1049}});
  assert_monitor_credit(
    "IMSI1", USED_TX, {{"1", 40}, {"3", 1024}, {"4", 1079}});

  // Collect updates, should only have mkeys 3 and 4
  auto session_update = local_enforcer->collect_updates();
  EXPECT_EQ(session_update.usage_monitors_size(), 2);
  for (const auto &monitor : session_update.usage_monitors()) {
    EXPECT_EQ(monitor.sid(), "IMSI1");
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

  assert_charging_credit("IMSI1", REPORTING_RX, {{1, 0}, {2, 0}});
  assert_charging_credit("IMSI1", REPORTING_TX, {{1, 0}, {2, 0}});
  assert_monitor_credit(
    "IMSI1", REPORTING_RX, {{"1", 0}, {"3", 1024}, {"4", 1049}});
  assert_monitor_credit(
    "IMSI1", REPORTING_TX, {{"1", 0}, {"3", 1024}, {"4", 1079}});

  UpdateSessionResponse update_response;
  auto monitor_updates = update_response.mutable_usage_monitor_responses();
  create_monitor_update_response(
    "IMSI1",
    "3",
    MonitoringLevel::PCC_RULE_LEVEL,
    2048,
    monitor_updates->Add());
  create_monitor_update_response(
    "IMSI1", "4", MonitoringLevel::SESSION_LEVEL, 2048, monitor_updates->Add());
  local_enforcer->update_session_credit(update_response);
  assert_monitor_credit("IMSI1", REPORTING_RX, {{"3", 0}, {"4", 0}});
  assert_monitor_credit("IMSI1", REPORTING_TX, {{"3", 0}, {"4", 0}});
  assert_monitor_credit("IMSI1", REPORTED_RX, {{"3", 1024}, {"4", 1049}});
  assert_monitor_credit("IMSI1", REPORTED_TX, {{"3", 1024}, {"4", 1079}});
  assert_monitor_credit("IMSI1", ALLOWED_TOTAL, {{"3", 4096}, {"4", 4176}});

  // Test rule removal in usage monitor response for CCA-Update
  update_response.Clear();
  monitor_updates = update_response.mutable_usage_monitor_responses();
  auto monitor_updates_response = monitor_updates->Add();
  create_monitor_update_response(
    "IMSI1",
    "3",
    MonitoringLevel::PCC_RULE_LEVEL,
    0,
    monitor_updates_response);
  monitor_updates_response->add_rules_to_remove("pcrf_only");

  EXPECT_CALL(
    *pipelined_client,
    deactivate_flows_for_rules("IMSI1",
      std::vector<std::string>{"pcrf_only"}, CheckCount(0)))
    .Times(1)
    .WillOnce(testing::Return(true));
  local_enforcer->update_session_credit(update_response);

  // Test rule installation in usage monitor response for CCA-Update
  update_response.Clear();
  monitor_updates = update_response.mutable_usage_monitor_responses();
  monitor_updates_response = monitor_updates->Add();
  create_monitor_update_response(
    "IMSI1",
    "3",
    MonitoringLevel::PCC_RULE_LEVEL,
    0,
    monitor_updates_response);

  StaticRuleInstall static_rule_install;
  static_rule_install.set_rule_id("pcrf_only");
  auto res_rules_to_install =
        monitor_updates_response->add_static_rules_to_install();
  res_rules_to_install->CopyFrom(static_rule_install);

  EXPECT_CALL(
    *pipelined_client,
    activate_flows_for_rules("IMSI1", testing::_,
      std::vector<std::string>{"pcrf_only"}, CheckCount(0)))
    .Times(1)
    .WillOnce(testing::Return(true));
  local_enforcer->update_session_credit(update_response);
}

TEST_F(LocalEnforcerTest, test_rar_create_dedicated_bearer)
{
  SessionState::QoSInfo test_qos_info;
  test_qos_info.enabled = true;
  test_qos_info.qci = 0;

  SessionState::Config test_volte_cfg;
  test_volte_cfg.ue_ipv4 = "127.0.0.1";
  test_volte_cfg.bearer_id = 1;
  test_volte_cfg.qos_info = test_qos_info;

  CreateSessionResponse response;
  local_enforcer->init_session_credit(
    "IMSI1", "1234", test_volte_cfg, response);

  PolicyReAuthRequest rar;
  std::vector<std::string> rules_to_remove;
  std::vector<StaticRuleInstall> rules_to_install;
  std::vector<DynamicRuleInstall> dynamic_rules_to_install;
  std::vector<EventTrigger> event_triggers;
  std::vector<UsageMonitoringCredit> usage_monitoring_credits;
  create_policy_reauth_request(
    "1234",
    "IMSI1",
    rules_to_remove,
    rules_to_install,
    dynamic_rules_to_install,
    event_triggers,
    time(NULL),
    usage_monitoring_credits,
    &rar);
  auto rar_qos_info = rar.mutable_qos_info();
  rar_qos_info->set_qci(QCI_1);

  EXPECT_CALL(
    *spgw_client,
    create_dedicated_bearer(testing::_, testing::_, testing::_, testing::_))
  .Times(1)
  .WillOnce(testing::Return(true));

  PolicyReAuthAnswer raa;
  local_enforcer->init_policy_reauth(rar, raa);
  EXPECT_EQ(raa.result(), ReAuthResult::UPDATE_INITIATED);
}

TEST_F(LocalEnforcerTest, test_rar_session_not_found)
{
  PolicyReAuthRequest rar;
  std::vector<std::string> rules_to_remove;
  std::vector<StaticRuleInstall> rules_to_install;
  std::vector<DynamicRuleInstall> dynamic_rules_to_install;
  std::vector<EventTrigger> event_triggers {EventTrigger::REVALIDATION_TIMEOUT};
  std::vector<UsageMonitoringCredit> usage_monitoring_credits;
  create_policy_reauth_request(
    "session1",
    "IMSI1",
    rules_to_remove,
    rules_to_install,
    dynamic_rules_to_install,
    event_triggers,
    time(NULL),
    usage_monitoring_credits,
    &rar);
  PolicyReAuthAnswer raa;
  local_enforcer->init_policy_reauth(rar, raa);
  EXPECT_EQ(raa.result(), ReAuthResult::SESSION_NOT_FOUND);
}

TEST_F(LocalEnforcerTest, test_rar_revalidation_timer)
{
  // init session first
  CreateSessionResponse response;
  create_credit_update_response(
    "IMSI1", 1, 3072, response.mutable_credits()->Add());
  local_enforcer->init_session_credit("IMSI1", "session1", test_cfg, response);
  insert_static_rule(1, "", "rule1");

  RuleRecordTable table;
  auto record_list = table.mutable_records();
  create_rule_record("IMSI1", "rule1", 1024, 2048, record_list->Add());
  local_enforcer->aggregate_records(table);

  // only RAR with REVALIDATION_TIMEOUT as one of the event triggers
  // will initiate an update
  EXPECT_CALL(*reporter, report_updates(_, _)).Times(1);

  PolicyReAuthRequest rar;
  std::vector<std::string> rules_to_remove;
  std::vector<StaticRuleInstall> rules_to_install;
  std::vector<DynamicRuleInstall> dynamic_rules_to_install;
  std::vector<EventTrigger> event_triggers {EventTrigger::TAI_CHANGE};
  std::vector<UsageMonitoringCredit> usage_monitoring_credits;
  create_policy_reauth_request(
    "session1",
    "IMSI1",
    rules_to_remove,
    rules_to_install,
    dynamic_rules_to_install,
    event_triggers,
    time(NULL),
    usage_monitoring_credits,
    &rar);
  PolicyReAuthAnswer raa;
  local_enforcer->init_policy_reauth(rar, raa);
  EXPECT_EQ(raa.result(), ReAuthResult::UPDATE_INITIATED);

  rar.add_event_triggers(EventTrigger::REVALIDATION_TIMEOUT);
  local_enforcer->init_policy_reauth(rar, raa);
  EXPECT_EQ(raa.result(), ReAuthResult::UPDATE_INITIATED);
}

TEST_F(LocalEnforcerTest, test_pipelined_setup)
{
  // insert into rule store first so init_session_credit can find the rule
  insert_static_rule(1, "", "rule2");

  CreateSessionResponse response;
  create_credit_update_response(
    "IMSI1", 1, 1024, response.mutable_credits()->Add());
  auto epoch = 145;
  auto dynamic_rule = response.mutable_dynamic_rules()->Add();
  auto policy_rule = dynamic_rule->mutable_policy_rule();
  policy_rule->set_id("rule1");
  policy_rule->set_rating_group(1);
  policy_rule->set_tracking_type(PolicyRule::ONLY_OCS);
  auto static_rule = response.mutable_static_rules()->Add();
  static_rule->set_rule_id("rule2");
  local_enforcer->init_session_credit("IMSI1", "1234", test_cfg, response);

  CreateSessionResponse response2;
  create_credit_update_response(
    "IMSI2", 1, 2048, response2.mutable_credits()->Add());
  auto dynamic_rule2 = response2.mutable_dynamic_rules()->Add();
  auto policy_rule2 = dynamic_rule2->mutable_policy_rule();
  policy_rule2->set_id("rule22");
  policy_rule2->set_rating_group(1);
  policy_rule2->set_tracking_type(PolicyRule::ONLY_OCS);
  local_enforcer->init_session_credit("IMSI2", "12345", test_cfg, response2);

  std::vector<std::string> imsi_list = {"IMSI2", "IMSI1"};
  std::vector<std::string> ip_address_list = {"127.0.0.1", "127.0.0.1"};
  std::vector<std::vector<std::string>> static_rule_list = {{}, {"rule2"}};
  std::vector<std::vector<std::string>> dynamic_rule_list = {{"rule22"},
                                                             {"rule1"}};

  EXPECT_CALL(
    *pipelined_client,
    setup(CheckSessionInfos(imsi_list, ip_address_list, static_rule_list,
                            dynamic_rule_list), testing::_, testing::_))
    .Times(1)
    .WillOnce(testing::Return(true));

    local_enforcer->setup(epoch,
      [](Status status, SetupFlowsResult resp) {});
}

int main(int argc, char **argv)
{
  ::testing::InitGoogleTest(&argc, argv);
  FLAGS_logtostderr = 1;
  FLAGS_v = 10;
  return RUN_ALL_TESTS();
}

} // namespace magma

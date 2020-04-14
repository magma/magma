/**
 * Copyright (c) 2016-present, Facebook, Inc.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include <memory>

#include <folly/io/async/EventBaseManager.h>
#include <glog/logging.h>
#include <gtest/gtest.h>

#include "LocalEnforcer.h"
#include "MagmaService.h"
#include "ProtobufCreators.h"
#include "RuleStore.h"
#include "ServiceRegistrySingleton.h"
#include "SessionID.h"
#include "SessionProxyResponderHandler.h"
#include "SessionState.h"
#include "SessionStore.h"
#include "SessiondMocks.h"
#include "StoredState.h"
#include "magma_logging.h"

using ::testing::Test;

namespace magma {

class SessionProxyResponderHandlerTest : public ::testing::Test {
protected:
  virtual void SetUp() {
    evb = new folly::EventBase();
    std::thread([&]() {
      std::cout << "Started event loop thread\n";
      folly::EventBaseManager::get()->setEventBase(evb, 0);
    }).detach();

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

    reporter = std::make_shared<MockSessionReporter>();
    auto rule_store = std::make_shared<StaticRuleStore>();
    session_store = std::make_shared<SessionStore>(rule_store);
    pipelined_client = std::make_shared<MockPipelinedClient>();
    auto directoryd_client = std::make_shared<MockDirectorydClient>();
    auto eventd_client = std::make_shared<MockEventdClient>();
    auto spgw_client = std::make_shared<MockSpgwServiceClient>();
    auto aaa_client = std::make_shared<MockAAAClient>();
    local_enforcer = std::make_shared<LocalEnforcer>(
        reporter, rule_store, *session_store, pipelined_client,
        directoryd_client, eventd_client, spgw_client, aaa_client, 0, 0);
    session_map = SessionMap{};

    proxy_responder = std::make_shared<SessionProxyResponderHandlerImpl>(
        local_enforcer, *session_store);

    local_enforcer->attachEventBase(evb);
  }

  std::unique_ptr<SessionState>
  get_session(std::string session_id,
              std::shared_ptr<StaticRuleStore> rule_store) {
    std::string hardware_addr_bytes = {0x0f, 0x10, 0x2e, 0x12, 0x3a, 0x55};
    std::string msisdn = "5100001234";
    std::string radius_session_id = "AA-AA-AA-AA-AA-AA:TESTAP__"
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

  UsageMonitoringUpdateResponse *get_monitoring_update() {
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

  PolicyReAuthRequest *get_policy_reauth_request() {
    //    message PolicyReAuthRequest {
    //      // NOTE: if no session_id is specified, apply to all sessions for
    //      the IMSI
    //      string session_id = 1;
    //      string imsi = 2;
    //      repeated string rules_to_remove = 3;
    //      repeated StaticRuleInstall rules_to_install = 6;
    //      repeated DynamicRuleInstall dynamic_rules_to_install = 7;
    //      repeated EventTrigger event_triggers = 8;
    //      google.protobuf.Timestamp revalidation_time = 9;
    //      repeated UsageMonitoringCredit usage_monitoring_credits = 10;
    //      QoSInformation qos_info = 11;
    //    }
    auto request = new PolicyReAuthRequest();
    request->set_session_id("");
    request->set_imsi("IMSI1");

    auto static_rule = new StaticRuleInstall();
    static_rule->set_rule_id("static_1");
    request->mutable_rules_to_install()->AddAllocated(static_rule);
    return request;
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

  std::shared_ptr<MockPipelinedClient> pipelined_client;
  std::shared_ptr<SessionProxyResponderHandlerImpl> proxy_responder;
  std::shared_ptr<MockSessionReporter> reporter;
  std::shared_ptr<LocalEnforcer> local_enforcer;
  SessionIDGenerator id_gen_;
  folly::EventBase *evb;
  SessionMap session_map;
  std::shared_ptr<SessionStore> session_store;
  std::shared_ptr<StaticRuleStore> rule_store;
};

TEST_F(SessionProxyResponderHandlerTest, test_policy_reauth) {
  // 1) Create SessionStore
  auto rule_store = std::make_shared<StaticRuleStore>();

  // 2) Create bare-bones session for IMSI1
  auto uc = get_default_update_criteria();
  auto session = get_session(sid, rule_store);
  session->activate_static_rule(rule_id_3, uc);
  EXPECT_EQ(session->get_session_id(), sid);
  EXPECT_EQ(session->get_request_number(), 2);
  EXPECT_EQ(session->is_static_rule_installed(rule_id_3), true);

  auto credit_update = get_monitoring_update();
  UsageMonitoringUpdateResponse &credit_update_ref = *credit_update;
  session->get_monitor_pool().receive_credit(credit_update_ref, uc);

  // Add some used credit
  session->get_monitor_pool().add_used_credit(monitoring_key, uint64_t(111),
                                              uint64_t(333), uc);
  EXPECT_EQ(session->get_monitor_pool().get_credit(monitoring_key, USED_TX),
            111);
  EXPECT_EQ(session->get_monitor_pool().get_credit(monitoring_key, USED_RX),
            333);

  // 3) Commit session for IMSI1 into SessionStore
  auto sessions = std::vector<std::unique_ptr<SessionState>>{};
  EXPECT_EQ(sessions.size(), 0);
  sessions.push_back(std::move(session));
  EXPECT_EQ(sessions.size(), 1);
  session_store->create_sessions(imsi, std::move(sessions));

  // Just verify some things about the session before doing PolicyReAuth
  SessionRead read_req = {};
  read_req.insert(imsi);
  auto session_map = session_store->read_sessions(read_req);
  EXPECT_EQ(session_map.size(), 1);
  EXPECT_EQ(session_map[imsi].size(), 1);
  EXPECT_EQ(session_map[imsi].front()->get_request_number(), 2);
  EXPECT_EQ(session_map[imsi].front()->is_static_rule_installed("static_1"),
            false);

  // 4) Now call PolicyReAuth
  auto request = get_policy_reauth_request();
  grpc::ServerContext create_context;
  EXPECT_CALL(*pipelined_client, activate_flows_for_rules(_, _, _, _)).Times(1);
  proxy_responder->PolicyReAuth(
      &create_context, request,
      [this](grpc::Status status, PolicyReAuthAnswer response_out) {});

  // run LocalEnforcer's init_policy_reauth which was scheduled by
  // proxy_responder
  evb->loopOnce();

  // 5) Read the session back from SessionStore and verify that the update
  //    was done correctly to what's stored.
  session_map = session_store->read_sessions(read_req);
  EXPECT_EQ(session_map.size(), 1);
  EXPECT_EQ(session_map[imsi].size(), 1);
  EXPECT_EQ(session_map[imsi].front()->get_request_number(), 2);
  EXPECT_EQ(session_map[imsi].front()->is_static_rule_installed("static_1"),
            true);
}

int main(int argc, char **argv) {
  ::testing::InitGoogleTest(&argc, argv);
  return RUN_ALL_TESTS();
}

} // namespace magma

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

#include <memory>

#include <folly/io/async/EventBaseManager.h>
#include <glog/logging.h>
#include <gtest/gtest.h>

#include "Consts.h"
#include "LocalEnforcer.h"
#include "MagmaService.h"
#include "Matchers.h"
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
    })
        .detach();

    monitoring_key = "mk1";
    rule_id        = "test_rule_1";

    reporter               = std::make_shared<MockSessionReporter>();
    auto rule_store        = std::make_shared<StaticRuleStore>();
    session_store          = std::make_shared<SessionStore>(rule_store);
    pipelined_client       = std::make_shared<MockPipelinedClient>();
    auto directoryd_client = std::make_shared<MockDirectorydClient>();
    auto spgw_client       = std::make_shared<MockSpgwServiceClient>();
    auto aaa_client        = std::make_shared<MockAAAClient>();
    auto events_reporter   = std::make_shared<MockEventsReporter>();
    auto default_mconfig   = get_default_mconfig();
    local_enforcer         = std::make_shared<LocalEnforcer>(
        reporter, rule_store, *session_store, pipelined_client,
        directoryd_client, events_reporter, spgw_client, aaa_client, 0, 0,
        default_mconfig);
    session_map = SessionMap{};

    proxy_responder = std::make_shared<SessionProxyResponderHandlerImpl>(
        local_enforcer, *session_store);

    local_enforcer->attachEventBase(evb);
  }

  std::unique_ptr<SessionState> get_session(
      std::shared_ptr<StaticRuleStore> rule_store) {
    SessionConfig cfg;
    cfg.common_context =
        build_common_context(IMSI1, IP1, "", APN1, MSISDN, TGPP_WLAN);
    const auto& wlan = build_wlan_context(MAC_ADDR, RADIUS_SESSION_ID);
    cfg.rat_specific_context.mutable_wlan_context()->CopyFrom(wlan);
    auto tgpp_context   = TgppContext{};
    auto pdp_start_time = 12345;
    return std::make_unique<SessionState>(
        IMSI1, SESSION_ID_1, cfg, *rule_store, tgpp_context, pdp_start_time);
  }

  UsageMonitoringUpdateResponse* get_monitoring_update() {
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
    // Don't set the TgppContext, assume gx_gy_relay disabled
    return credit_update;
  }

  PolicyReAuthRequest* get_policy_reauth_request() {
    auto request = new PolicyReAuthRequest();
    request->set_session_id("");
    request->set_imsi(IMSI1);

    auto static_rule = new StaticRuleInstall();
    static_rule->set_rule_id("static_1");

    // This should be a duplicate rule
    auto static_rule_2 = new StaticRuleInstall();
    static_rule_2->set_rule_id(rule_id);

    request->mutable_rules_to_install()->AddAllocated(static_rule);
    request->mutable_rules_to_install()->AddAllocated(static_rule_2);
    return request;
  }

 protected:
  std::string monitoring_key;
  std::string rule_id;

  std::shared_ptr<MockPipelinedClient> pipelined_client;
  std::shared_ptr<SessionProxyResponderHandlerImpl> proxy_responder;
  std::shared_ptr<MockSessionReporter> reporter;
  std::shared_ptr<LocalEnforcer> local_enforcer;
  SessionIDGenerator id_gen_;
  folly::EventBase* evb;
  SessionMap session_map;
  std::shared_ptr<SessionStore> session_store;
  std::shared_ptr<StaticRuleStore> rule_store;
};

TEST_F(SessionProxyResponderHandlerTest, test_policy_reauth) {
  // 1) Create SessionStore
  auto rule_store = std::make_shared<StaticRuleStore>();

  // 2) Create bare-bones session for IMSI1
  auto uc      = get_default_update_criteria();
  auto session = get_session(rule_store);
  RuleLifetime lifetime{
      .activation_time   = std::time_t(0),
      .deactivation_time = std::time_t(0),
  };
  session->activate_static_rule(rule_id, lifetime, uc);
  EXPECT_EQ(session->get_session_id(), SESSION_ID_1);
  EXPECT_EQ(session->get_request_number(), 1);
  EXPECT_EQ(session->is_static_rule_installed(rule_id), true);

  auto credit_update                               = get_monitoring_update();
  UsageMonitoringUpdateResponse& credit_update_ref = *credit_update;
  session->receive_monitor(credit_update_ref, uc);

  // Add some used credit
  session->add_to_monitor(monitoring_key, uint64_t(111), uint64_t(333), uc);
  EXPECT_EQ(session->get_monitor(monitoring_key, USED_TX), 111);
  EXPECT_EQ(session->get_monitor(monitoring_key, USED_RX), 333);

  // 3) Commit session for IMSI1 into SessionStore
  auto sessions = SessionVector{};
  EXPECT_EQ(sessions.size(), 0);
  sessions.push_back(std::move(session));
  EXPECT_EQ(sessions.size(), 1);
  session_store->create_sessions(IMSI1, std::move(sessions));

  // Just verify some things about the session before doing PolicyReAuth
  SessionRead read_req = {};
  read_req.insert(IMSI1);
  auto session_map = session_store->read_sessions(read_req);
  EXPECT_EQ(session_map.size(), 1);
  EXPECT_EQ(session_map[IMSI1].size(), 1);
  EXPECT_EQ(session_map[IMSI1].front()->get_request_number(), 1);
  EXPECT_EQ(
      session_map[IMSI1].front()->is_static_rule_installed("static_1"), false);

  // 4) Now call PolicyReAuth
  //    This is done with a duplicate install of rule_id. This checks that
  //    duplicate rule installs are ignored, as they are occasionally
  //    requested by the session proxy. If the duplicate rule install causes
  //    a failure, then the entire PolicyReAuth will not save to the
  //    SessionStore properly
  auto request = get_policy_reauth_request();
  grpc::ServerContext create_context;
  EXPECT_CALL(
      *pipelined_client,
      activate_flows_for_rules(IMSI1, _, _, _, _, CheckCount(1), _, _))
      .Times(1);
  proxy_responder->PolicyReAuth(
      &create_context, request,
      [this](grpc::Status status, PolicyReAuthAnswer response_out) {});

  // run LocalEnforcer's init_policy_reauth which was scheduled by
  // proxy_responder
  evb->loopOnce();

  // 5) Read the session back from SessionStore and verify that the update
  //    was done correctly to what's stored.
  //    If the PolicyReAuth failed, then rule static_1 will not have been
  //    installed.
  session_map = session_store->read_sessions(read_req);
  EXPECT_EQ(session_map.size(), 1);
  EXPECT_EQ(session_map[IMSI1].size(), 1);
  EXPECT_EQ(session_map[IMSI1].front()->get_request_number(), 1);
  EXPECT_EQ(
      session_map[IMSI1].front()->is_static_rule_installed("static_1"), true);
}

TEST_F(SessionProxyResponderHandlerTest, test_abort_session) {
  // 1) Create SessionStore
  auto rule_store = std::make_shared<StaticRuleStore>();

  // 2) Create bare-bones session for IMSI1
  auto uc      = get_default_update_criteria();
  auto session = get_session(rule_store);
  RuleLifetime lifetime{
      .activation_time   = std::time_t(0),
      .deactivation_time = std::time_t(0),
  };
  session->activate_static_rule(rule_id, lifetime, uc);
  EXPECT_EQ(session->get_session_id(), SESSION_ID_1);
  EXPECT_EQ(session->get_request_number(), 1);
  EXPECT_EQ(session->is_static_rule_installed(rule_id), true);

  // 3) Commit session for IMSI1 into SessionStore
  auto sessions = SessionVector{};
  EXPECT_EQ(sessions.size(), 0);
  sessions.push_back(std::move(session));
  EXPECT_EQ(sessions.size(), 1);
  session_store->create_sessions(IMSI1, std::move(sessions));

  // Just verify some things about the session before doing AbortSession
  SessionRead read_req = {};
  read_req.insert(IMSI1);
  auto session_map = session_store->read_sessions(read_req);
  EXPECT_EQ(session_map.size(), 1);
  EXPECT_EQ(session_map[IMSI1].size(), 1);
  EXPECT_EQ(session_map[IMSI1].front()->get_request_number(), 1);
  EXPECT_EQ(
      session_map[IMSI1].front()->is_static_rule_installed("static_1"), false);

  // 4) Now call AbortSession
  //    We should see a session deletion which would trigger deactivate_flows
  AbortSessionRequest request;
  request.set_user_name(IMSI1);
  request.set_session_id(SESSION_ID_1);
  grpc::ServerContext create_context;
  EXPECT_CALL(
      *pipelined_client,
      deactivate_flows_for_rules_for_termination(
          IMSI1, _, _, CheckCount(1), CheckCount(0), RequestOriginType::GX))
      .Times(1)
      .WillOnce(testing::Return(true));
  proxy_responder->AbortSession(
      &create_context, &request,
      [this](grpc::Status status, AbortSessionResult response_out) {});

  // The actual work of AbortSession gets pushed off to an event loop so run it
  // once
  evb->loopOnce();

  // 5) Read the session back from SessionStore and verify that the update
  //    was done correctly to what's stored.
  //    If the AbortSession failed, then the session should still exist
  session_map = session_store->read_sessions(read_req);
  EXPECT_EQ(session_map[IMSI1].size(), 0);
}

int main(int argc, char** argv) {
  ::testing::InitGoogleTest(&argc, argv);
  return RUN_ALL_TESTS();
}

}  // namespace magma

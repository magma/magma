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
#include <glog/logging.h>
#include <gtest/gtest.h>

#include <memory>

#include "Consts.h"
#include "LocalEnforcer.h"
#include "magma_logging.h"
#include "includes/MagmaService.h"
#include "Matchers.h"
#include "ProtobufCreators.h"
#include "RuleStore.h"
#include "includes/ServiceRegistrySingleton.h"
#include "SessiondMocks.h"
#include "SessionID.h"
#include "SessionProxyResponderHandler.h"
#include "SessionState.h"
#include "SessionStore.h"
#include "StoredState.h"

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
    rule_id_1      = "test_rule_1";
    rule_id_2      = "test_rule_2";

    reporter      = std::make_shared<MockSessionReporter>();
    rule_store    = std::make_shared<StaticRuleStore>();
    session_store = std::make_shared<SessionStore>(
        rule_store, std::make_shared<MeteringReporter>());
    pipelined_client     = std::make_shared<MockPipelinedClient>();
    auto spgw_client     = std::make_shared<MockSpgwServiceClient>();
    auto aaa_client      = std::make_shared<MockAAAClient>();
    auto events_reporter = std::make_shared<MockEventsReporter>();
    auto default_mconfig = get_default_mconfig();
    local_enforcer       = std::make_shared<LocalEnforcer>(
        reporter, rule_store, *session_store, pipelined_client, events_reporter,
        spgw_client, aaa_client, 0, 0, default_mconfig);
    session_map = SessionMap{};

    proxy_responder = std::make_shared<SessionProxyResponderHandlerImpl>(
        local_enforcer, *session_store);

    local_enforcer->attachEventBase(evb);
  }

  virtual void TearDown() { delete evb; }

  std::unique_ptr<SessionState> get_session(
      std::shared_ptr<StaticRuleStore> rule_store) {
    SessionConfig cfg;
    Teids teids;
    teids.set_agw_teid(TEID_1_UL);
    teids.set_enb_teid(TEID_1_DL);
    cfg.common_context =
        build_common_context(IMSI1, IP1, "", teids, APN1, MSISDN, TGPP_WLAN);
    const auto& wlan = build_wlan_context(MAC_ADDR, RADIUS_SESSION_ID);
    cfg.rat_specific_context.mutable_wlan_context()->CopyFrom(wlan);
    auto tgpp_context   = TgppContext{};
    auto pdp_start_time = 12345;
    return std::make_unique<SessionState>(
        IMSI1, SESSION_ID_1, cfg, *rule_store, tgpp_context, pdp_start_time,
        CreateSessionResponse{});
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
    auto tx    = units->mutable_tx();
    auto rx    = units->mutable_rx();

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

  PolicyReAuthRequest get_policy_reauth_request() {
    PolicyReAuthRequest request;
    request.set_session_id("");
    request.set_imsi(IMSI1);

    StaticRuleInstall static_rule_1;
    static_rule_1.set_rule_id(rule_id_2);

    // This should be a duplicate rule
    StaticRuleInstall static_rule_2;
    static_rule_2.set_rule_id(rule_id_1);

    request.mutable_rules_to_install()->Add()->CopyFrom(static_rule_1);
    request.mutable_rules_to_install()->Add()->CopyFrom(static_rule_2);
    return request;
  }

  void insert_static_rule(const std::string& rule_id) {
    rule_store->insert_rule(create_policy_rule(rule_id, "", 0));
  }

 protected:
  std::string monitoring_key;
  std::string rule_id_1;
  std::string rule_id_2;

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
  // 1) Initialize static rules
  insert_static_rule(rule_id_1);
  insert_static_rule(rule_id_2);

  // 2) Create bare-bones session for IMSI1
  auto uc      = get_default_update_criteria();
  auto session = get_session(rule_store);
  RuleLifetime lifetime;
  session->activate_static_rule(rule_id_1, lifetime, nullptr);
  EXPECT_EQ(session->get_session_id(), SESSION_ID_1);
  EXPECT_EQ(session->get_request_number(), 1);
  EXPECT_EQ(session->is_static_rule_installed(rule_id_1), true);

  auto monitor_update = get_monitoring_update();
  session->receive_monitor(monitor_update, &uc);

  // Add some used credit
  session->add_to_monitor(
      monitoring_key, uint64_t(111), uint64_t(333), nullptr);
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
      session_map[IMSI1].front()->is_static_rule_installed(rule_id_2), false);

  // 4) Now call PolicyReAuth
  //    This is done with a duplicate install of rule_id_1. This checks that
  //    duplicate rule installs are ignored, as they are occasionally
  //    requested by the session proxy. If the duplicate rule install causes
  //    a failure, then the entire PolicyReAuth will not save to the
  //    SessionStore properly
  auto request = get_policy_reauth_request();
  grpc::ServerContext create_context;
  EXPECT_CALL(
      *pipelined_client,
      activate_flows_for_rules(IMSI1, _, _, _, _, _, CheckRuleCount(1), _))
      .Times(1);
  proxy_responder->PolicyReAuth(
      &create_context, &request,
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
      session_map[IMSI1].front()->is_static_rule_installed(rule_id_2), true);
}

TEST_F(SessionProxyResponderHandlerTest, test_abort_session) {
  // 1) Initialize static rules
  insert_static_rule(rule_id_1);
  insert_static_rule(rule_id_2);

  // 2) Create bare-bones session for IMSI1
  auto uc      = get_default_update_criteria();
  auto session = get_session(rule_store);
  RuleLifetime lifetime;
  session->activate_static_rule(rule_id_1, lifetime, &uc);
  EXPECT_EQ(session->get_session_id(), SESSION_ID_1);
  EXPECT_EQ(session->get_request_number(), 1);
  EXPECT_EQ(session->is_static_rule_installed(rule_id_1), true);

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
      session_map[IMSI1].front()->is_static_rule_installed(rule_id_2), false);

  // 4) Now call AbortSession
  //    We should see a session deletion which would trigger deactivate_flows
  AbortSessionRequest request;
  request.set_user_name(IMSI1);
  request.set_session_id(SESSION_ID_1);
  grpc::ServerContext create_context;
  // the request should has no rules so PipelineD deletes all rules
  EXPECT_CALL(
      *pipelined_client, deactivate_flows_for_rules_for_termination(
                             IMSI1, _, _, _, RequestOriginType::WILDCARD))
      .Times(1);
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

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
#include <chrono>
#include <future>
#include <memory>
#include <string.h>
#include <time.h>

#include <folly/io/async/EventBaseManager.h>
#include <gtest/gtest.h>
#include <lte/protos/session_manager.grpc.pb.h>

#include "Consts.h"
#include "LocalEnforcer.h"
#include "MagmaService.h"
#include "Matchers.h"
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
  void SetUpWithMConfig(magma::mconfig::SessionD mconfig) {
    reporter          = std::make_shared<MockSessionReporter>();
    rule_store        = std::make_shared<StaticRuleStore>();
    session_store     = std::make_shared<SessionStore>(rule_store);
    pipelined_client  = std::make_shared<MockPipelinedClient>();
    directoryd_client = std::make_shared<MockDirectorydClient>();
    spgw_client       = std::make_shared<MockSpgwServiceClient>();
    aaa_client        = std::make_shared<MockAAAClient>();
    events_reporter   = std::make_shared<MockEventsReporter>();
    local_enforcer    = std::make_unique<LocalEnforcer>(
        reporter, rule_store, *session_store, pipelined_client,
        directoryd_client, events_reporter, spgw_client, aaa_client, 0, 0,
        mconfig);
    evb = folly::EventBaseManager::get()->getEventBase();
    local_enforcer->attachEventBase(evb);
    session_map = SessionMap{};
    cwf_session_config.common_context =
        build_common_context(IMSI1, "", "", "", "", TGPP_WLAN);
    wlan_context = build_wlan_context(MAC_ADDR, RADIUS_SESSION_ID);
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
    mconfig.set_gx_gy_relay_enabled(false);
    auto wallet_config = mconfig.mutable_wallet_exhaust_detection();
    wallet_config->set_terminate_on_exhaust(true);
    wallet_config->set_method(
        magma::mconfig::WalletExhaustDetection_Method_GxTrackedRules);
    return mconfig;
  }

  void insert_static_rule(
      uint32_t rating_group, const std::string& m_key,
      const std::string& rule_id) {
    PolicyRule rule;
    create_policy_rule(rule_id, m_key, rating_group, &rule);
    rule_store->insert_rule(rule);
  }

  void create_new_session(
      std::string imsi, std::string session_id, SessionConfig& cfg,
      std::vector<std::string> static_rules) {
    CreateSessionResponse response;

    // Initialize the session in session store
    auto session_map = session_store->read_sessions({imsi});
    local_enforcer->initialize_creating_session(
        session_map, imsi, session_id, cfg);
    EXPECT_EQ(session_map[IMSI1].size(), 1);
    EXPECT_EQ(session_map[IMSI1][0]->get_state(), CREATING);

    bool write_success =
        session_store->create_sessions(IMSI1, std::move(session_map[imsi]));
    EXPECT_TRUE(write_success);

    response.set_session_id(session_id);

    create_session_create_response(
        imsi, session_id, monitoring_key, static_rules, &response);
    create_credit_update_response(
        imsi, session_id, 1, 1025, response.mutable_credits()->Add());

    session_map = session_store->read_sessions({imsi});
    EXPECT_TRUE(session_map.find(imsi) != session_map.end());
    EXPECT_EQ(session_map[imsi].size(), 1);
    auto updates = SessionStore::get_default_session_update(session_map);
    local_enforcer->process_create_session_response(
        session_map[imsi][0], imsi, session_id, response,
        updates[imsi][SESSION_ID_1]);
    EXPECT_EQ(session_map[IMSI1].size(), 1);
    EXPECT_EQ(session_map[IMSI1][0]->get_state(), SESSION_ACTIVE);
    write_success = session_store->update_sessions(updates);
    EXPECT_TRUE(write_success);
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
  std::shared_ptr<MockEventsReporter> events_reporter;
  SessionMap session_map;
  SessionConfig cwf_session_config;
  folly::EventBase* evb;
  WLANSessionContext wlan_context;
  std::string monitoring_key = "m1";
};

// Make sure sessions that are scheduled to be terminated before sync are
// correctly scheduled to be terminated again.
TEST_F(LocalEnforcerTest, test_termination_scheduling_on_sync_sessions) {
  SetUpWithMConfig(get_mconfig_gx_rule_wallet_exhaust());
  insert_static_rule(0, monitoring_key, "static_1");

  EXPECT_CALL(
      *pipelined_client,
      update_subscriber_quota_state(
          CheckSubscriberQuotaUpdate(SubscriberQuotaUpdate_Type_VALID_QUOTA)));
  create_new_session(IMSI1, SESSION_ID_1, cwf_session_config, {"static_1"});

  auto session_map    = session_store->read_sessions(SessionRead{IMSI1});
  auto session_update = session_store->get_default_session_update(session_map);
  EXPECT_EQ(session_map[IMSI1].size(), 1);

  auto& uc = session_update[IMSI1][SESSION_ID_1];
  // remove all monitored policies to trigger a termination schedule
  uc.static_rules_to_uninstall = {"static_1"};
  auto success                 = session_store->update_sessions(session_update);
  EXPECT_TRUE(success);

  EXPECT_CALL(
      *pipelined_client,
      update_subscriber_quota_state(
          CheckSubscriberQuotaUpdate(SubscriberQuotaUpdate_Type_NO_QUOTA)));

  // Syncing will schedule a termination for this IMSI
  local_enforcer->sync_sessions_on_restart(std::time_t(0));

  // Terminate subscriber is the only thing on the event queue, and
  // quota_exhaust_termination_on_init_ms is set to 0
  // We expect the termination to take place once we run evb->loopOnce()
  EXPECT_CALL(*aaa_client, terminate_session(_, _)).Times(1);
  EXPECT_CALL(
      *pipelined_client,
      update_subscriber_quota_state(
          CheckSubscriberQuotaUpdate(SubscriberQuotaUpdate_Type_TERMINATE)));
  evb->loopOnce();

  // At this point, the state should have transitioned from
  // SESSION_ACTIVE -> SESSION_RELEASED
  session_map            = session_store->read_sessions(SessionRead{IMSI1});
  auto updated_fsm_state = session_map[IMSI1].front()->get_state();
  EXPECT_EQ(updated_fsm_state, SESSION_RELEASED);
}

TEST_F(LocalEnforcerTest, test_cwf_quota_exhaustion_on_init_has_quota) {
  SetUpWithMConfig(get_mconfig_gx_rule_wallet_exhaust());
  insert_static_rule(0, monitoring_key, "static_1");

  std::vector<SubscriberQuotaUpdate_Type> expected_states{
      SubscriberQuotaUpdate_Type_VALID_QUOTA};
  EXPECT_CALL(
      *pipelined_client,
      update_subscriber_quota_state(
          CheckSubscriberQuotaUpdate(SubscriberQuotaUpdate_Type_VALID_QUOTA)))
      .Times(1)
      .WillOnce(testing::Return(true));
  create_new_session(IMSI1, SESSION_ID_1, cwf_session_config, {"static_1"});
}

TEST_F(LocalEnforcerTest, test_cwf_quota_exhaustion_on_init_no_quota) {
  SetUpWithMConfig(get_mconfig_gx_rule_wallet_exhaust());
  EXPECT_CALL(
      *pipelined_client,
      update_subscriber_quota_state(
          CheckSubscriberQuotaUpdate(SubscriberQuotaUpdate_Type_NO_QUOTA)))
      .Times(1)
      .WillOnce(testing::Return(true));
  create_new_session(
      IMSI1, SESSION_ID_1, cwf_session_config, {});  // no rule installs
}

TEST_F(LocalEnforcerTest, test_cwf_quota_exhaustion_on_rar) {
  SetUpWithMConfig(get_mconfig_gx_rule_wallet_exhaust());
  // setup : successful session creation with valid monitoring quota
  insert_static_rule(0, monitoring_key, "static_1");
  EXPECT_CALL(
      *pipelined_client,
      update_subscriber_quota_state(
          CheckSubscriberQuotaUpdate(SubscriberQuotaUpdate_Type_VALID_QUOTA)))
      .Times(1)
      .WillOnce(testing::Return(true));
  create_new_session(IMSI1, SESSION_ID_1, cwf_session_config, {"static_1"});

  // send a policy reauth request with rule removals for "static_1" to indicate
  // total monitoring quota exhaustion
  PolicyReAuthRequest request;
  request.set_session_id(SESSION_ID_1);
  request.set_imsi(IMSI1);
  request.add_rules_to_remove("static_1");

  std::vector<SubscriberQuotaUpdate_Type> expected_states{
      SubscriberQuotaUpdate_Type_TERMINATE};
  EXPECT_CALL(
      *pipelined_client,
      update_subscriber_quota_state(
          CheckSubscriberQuotaUpdate(SubscriberQuotaUpdate_Type_TERMINATE)))
      .Times(1)
      .WillOnce(testing::Return(true));

  PolicyReAuthAnswer answer;
  auto session_map = session_store->read_sessions(SessionRead{IMSI1});
  auto update      = SessionStore::get_default_session_update(session_map);
  local_enforcer->init_policy_reauth(session_map, request, answer, update);
}

TEST_F(LocalEnforcerTest, test_cwf_quota_exhaustion_on_update) {
  SetUpWithMConfig(get_mconfig_gx_rule_wallet_exhaust());
  // setup : successful session creation with valid monitoring quota
  insert_static_rule(0, monitoring_key, "static_1");
  insert_static_rule(0, monitoring_key, "static_2");

  create_new_session(
      IMSI1, SESSION_ID_1, cwf_session_config, {"static_1", "static_2"});

  // remove only static_2, should not change anything in terms of quota since
  // static_1 is still active
  auto session_map = session_store->read_sessions(SessionRead{IMSI1});
  UpdateSessionResponse update_response;
  auto monitor = update_response.mutable_usage_monitor_responses()->Add();
  create_monitor_update_response(
      IMSI1, SESSION_ID_1, monitoring_key, MonitoringLevel::PCC_RULE_LEVEL, 2048,
      monitor);
  monitor->add_rules_to_remove("static_2");
  auto update = SessionStore::get_default_session_update(session_map);
  local_enforcer->update_session_credits_and_rules(
      session_map, update_response, update);

  // send an update response with rule removals for "static_1" to indicate
  // total monitoring quota exhaustion
  update_response.clear_usage_monitor_responses();
  monitor = update_response.mutable_usage_monitor_responses()->Add();
  create_monitor_update_response(
      IMSI1, SESSION_ID_1, "m1", MonitoringLevel::PCC_RULE_LEVEL, 0, monitor);
  monitor->add_rules_to_remove("static_1");

  std::vector<SubscriberQuotaUpdate_Type> expected_states = {
      SubscriberQuotaUpdate_Type_TERMINATE};
  EXPECT_CALL(
      *pipelined_client,
      update_subscriber_quota_state(
          CheckSubscriberQuotaUpdate(SubscriberQuotaUpdate_Type_TERMINATE)))
      .Times(1)
      .WillOnce(testing::Return(true));

  local_enforcer->update_session_credits_and_rules(
      session_map, update_response, update);
}

int main(int argc, char** argv) {
  ::testing::InitGoogleTest(&argc, argv);
  FLAGS_logtostderr = 1;
  FLAGS_v           = 10;
  return RUN_ALL_TESTS();
}

}  // namespace magma

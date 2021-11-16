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
#include "LocalEnforcer.h"
#include "magma_logging.h"
#include "includes/MagmaService.h"
#include "Matchers.h"
#include "ProtobufCreators.h"
#include "includes/ServiceRegistrySingleton.h"
#include "SessiondMocks.h"
#include "SessionStore.h"

#define SECONDS_A_DAY 86400

using grpc::ServerContext;
using grpc::Status;
using ::testing::Test;

namespace magma {

Teids teids0;

class LocalEnforcerTest : public ::testing::Test {
 protected:
  void SetUpWithMConfig(magma::mconfig::SessionD mconfig) {
    reporter = std::make_shared<MockSessionReporter>();
    rule_store = std::make_shared<StaticRuleStore>();
    session_store = std::make_shared<SessionStore>(
        rule_store, std::make_shared<MeteringReporter>());
    pipelined_client = std::make_shared<MockPipelinedClient>();
    spgw_client = std::make_shared<MockSpgwServiceClient>();
    aaa_client = std::make_shared<MockAAAClient>();
    events_reporter = std::make_shared<MockEventsReporter>();
    auto shard_tracker = std::make_shared<ShardTracker>();
    local_enforcer = std::make_unique<LocalEnforcer>(
        reporter, rule_store, *session_store, pipelined_client, events_reporter,
        spgw_client, aaa_client, shard_tracker, 0, 0, mconfig);
    evb = folly::EventBaseManager::get()->getEventBase();
    local_enforcer->attachEventBase(evb);
    session_map = SessionMap{};
    teids0.set_agw_teid(0);
    teids0.set_enb_teid(0);
    cwf_session_config.common_context =
        build_common_context(IMSI1, "", "", teids0, "", "", TGPP_WLAN);
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
    mconfig.set_gx_gy_relay_enabled(false);
    auto wallet_config = mconfig.mutable_wallet_exhaust_detection();
    wallet_config->set_terminate_on_exhaust(true);
    wallet_config->set_method(
        magma::mconfig::WalletExhaustDetection_Method_GxTrackedRules);
    return mconfig;
  }

  void insert_static_rule(uint32_t rating_group, const std::string& m_key,
                          const std::string& rule_id) {
    rule_store->insert_rule(create_policy_rule(rule_id, m_key, rating_group));
  }

  void initialize_session(SessionMap& session_map,
                          const std::string& session_id,
                          const SessionConfig& cfg,
                          const CreateSessionResponse& response) {
    const std::string imsi = cfg.get_imsi();
    auto session = local_enforcer->create_initializing_session(session_id, cfg);
    local_enforcer->update_session_with_policy_response(session, response,
                                                        nullptr);
    session_map[imsi].push_back(std::move(session));
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
  SessionConfig cwf_session_config;
  folly::EventBase* evb;
};

// Make sure sessions that are scheduled to be terminated before sync are
// correctly scheduled to be terminated again.
TEST_F(LocalEnforcerTest, test_termination_scheduling_on_sync_sessions) {
  SetUpWithMConfig(get_mconfig_gx_rule_wallet_exhaust());
  CreateSessionResponse response;

  std::vector<std::string> rules_to_install;
  rules_to_install.push_back("rule1");
  insert_static_rule(0, "m1", "rule1");

  // Create a CreateSessionResponse with one Gx monitor:m1 and one rule:rule1
  create_session_create_response(IMSI1, SESSION_ID_1, "m1", rules_to_install,
                                 &response);

  EXPECT_CALL(*pipelined_client,
              update_subscriber_quota_state(CheckSubscriberQuotaUpdate(
                  SubscriberQuotaUpdate_Type_VALID_QUOTA)));
  initialize_session(session_map, SESSION_ID_1, cwf_session_config, response);
  local_enforcer->update_tunnel_ids(
      session_map, create_update_tunnel_ids_request(IMSI1, 0, teids0));

  EXPECT_EQ(session_map[IMSI1].size(), 1);
  bool success =
      session_store->create_sessions(IMSI1, std::move(session_map[IMSI1]));
  EXPECT_TRUE(success);

  auto session_map = session_store->read_sessions(SessionRead{IMSI1});
  auto session_update = session_store->get_default_session_update(session_map);
  EXPECT_EQ(session_map[IMSI1].size(), 1);

  auto& uc = session_update[IMSI1][SESSION_ID_1];
  // remove all monitored policies to trigger a termination schedule
  uc.static_rules_to_uninstall = {"rule1"};
  success = session_store->update_sessions(session_update);
  EXPECT_TRUE(success);
  std::cerr << "\n\n going to call sync on restart \n\n\n";

  EXPECT_CALL(*pipelined_client,
              update_subscriber_quota_state(CheckSubscriberQuotaUpdate(
                  SubscriberQuotaUpdate_Type_NO_QUOTA)));

  // Syncing will schedule a termination for this IMSI
  local_enforcer->sync_sessions_on_restart(std::time_t(0));

  // Terminate subscriber is the only thing on the event queue, and
  // quota_exhaust_termination_on_init_ms is set to 0
  // We expect the termination to take place once we run evb->loopOnce()
  EXPECT_CALL(*aaa_client, terminate_session(_, _)).Times(1);
  EXPECT_CALL(*pipelined_client,
              update_subscriber_quota_state(CheckSubscriberQuotaUpdate(
                  SubscriberQuotaUpdate_Type_TERMINATE)));
  evb->loopOnce();

  // At this point, the state should have transitioned from
  // SESSION_ACTIVE -> SESSION_RELEASED
  session_map = session_store->read_sessions(SessionRead{IMSI1});
  auto updated_fsm_state = session_map[IMSI1].front()->get_state();
  EXPECT_EQ(updated_fsm_state, SESSION_RELEASED);
}

TEST_F(LocalEnforcerTest, test_cwf_quota_exhaustion_on_init_has_quota) {
  SetUpWithMConfig(get_mconfig_gx_rule_wallet_exhaust());
  insert_static_rule(0, "m1", "static_1");

  std::vector<std::string> static_rules{"static_1"};
  CreateSessionResponse response;
  create_session_create_response(IMSI1, SESSION_ID_1, "m1", static_rules,
                                 &response);

  StaticRuleInstall static_rule_install;
  static_rule_install.set_rule_id("static_1");
  auto res_rules_to_install = response.mutable_static_rules()->Add();
  res_rules_to_install->CopyFrom(static_rule_install);

  std::vector<SubscriberQuotaUpdate_Type> expected_states{
      SubscriberQuotaUpdate_Type_VALID_QUOTA};
  EXPECT_CALL(*pipelined_client,
              update_subscriber_quota_state(CheckSubscriberQuotaUpdate(
                  SubscriberQuotaUpdate_Type_VALID_QUOTA)))
      .Times(1);
  initialize_session(session_map, SESSION_ID_1, cwf_session_config, response);
  local_enforcer->update_tunnel_ids(
      session_map, create_update_tunnel_ids_request(IMSI1, 0, teids0));
}

TEST_F(LocalEnforcerTest, test_cwf_quota_exhaustion_on_init_no_quota) {
  SetUpWithMConfig(get_mconfig_gx_rule_wallet_exhaust());
  insert_static_rule(1, "m1", "static_1");

  std::vector<std::string> static_rules{};  // no rule installs
  CreateSessionResponse response;
  create_session_create_response(IMSI1, SESSION_ID_1, "m1", static_rules,
                                 &response);

  std::vector<SubscriberQuotaUpdate_Type> expected_states{
      SubscriberQuotaUpdate_Type_NO_QUOTA};
  EXPECT_CALL(*pipelined_client,
              update_subscriber_quota_state(CheckSubscriberQuotaUpdate(
                  SubscriberQuotaUpdate_Type_NO_QUOTA)))
      .Times(1);
  initialize_session(session_map, SESSION_ID_1, cwf_session_config, response);
  local_enforcer->update_tunnel_ids(
      session_map, create_update_tunnel_ids_request(IMSI1, 0, teids0));
}

TEST_F(LocalEnforcerTest, test_cwf_quota_exhaustion_on_rar) {
  SetUpWithMConfig(get_mconfig_gx_rule_wallet_exhaust());
  // setup : successful session creation with valid monitoring quota
  insert_static_rule(0, "m1", "static_1");

  std::vector<std::string> static_rules{"static_1"};
  CreateSessionResponse response;
  create_session_create_response(IMSI1, SESSION_ID_1, "m1", static_rules,
                                 &response);
  initialize_session(session_map, SESSION_ID_1, cwf_session_config, response);
  local_enforcer->update_tunnel_ids(
      session_map, create_update_tunnel_ids_request(IMSI1, 0, teids0));

  // send a policy reauth request with rule removals for "static_1" to indicate
  // total monitoring quota exhaustion
  PolicyReAuthRequest request;
  request.set_session_id("");
  request.set_imsi(IMSI1);
  request.add_rules_to_remove("static_1");

  std::vector<SubscriberQuotaUpdate_Type> expected_states{
      SubscriberQuotaUpdate_Type_TERMINATE};
  EXPECT_CALL(*pipelined_client,
              update_subscriber_quota_state(CheckSubscriberQuotaUpdate(
                  SubscriberQuotaUpdate_Type_TERMINATE)))
      .Times(1);

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
  CreateSessionResponse response;
  create_session_create_response(IMSI1, SESSION_ID_1, "m1", static_rules,
                                 &response);
  initialize_session(session_map, SESSION_ID_1, cwf_session_config, response);
  local_enforcer->update_tunnel_ids(
      session_map, create_update_tunnel_ids_request(IMSI1, 0, teids0));

  // remove only static_2, should not change anything in terms of quota since
  // static_1 is still active
  UpdateSessionResponse update_response;
  auto monitor = update_response.mutable_usage_monitor_responses()->Add();
  create_monitor_update_response(IMSI1, SESSION_ID_1, "m1",
                                 MonitoringLevel::PCC_RULE_LEVEL, 2048,
                                 monitor);
  monitor->add_rules_to_remove("static_2");
  auto update = SessionStore::get_default_session_update(session_map);
  local_enforcer->update_session_credits_and_rules(session_map, update_response,
                                                   update);

  // send an update response with rule removals for "static_1" to indicate
  // total monitoring quota exhaustion
  update_response.clear_usage_monitor_responses();
  monitor = update_response.mutable_usage_monitor_responses()->Add();
  create_monitor_update_response(IMSI1, SESSION_ID_1, "m1",
                                 MonitoringLevel::PCC_RULE_LEVEL, 0, monitor);
  monitor->add_rules_to_remove("static_1");

  std::vector<SubscriberQuotaUpdate_Type> expected_states = {
      SubscriberQuotaUpdate_Type_TERMINATE};
  EXPECT_CALL(*pipelined_client,
              update_subscriber_quota_state(CheckSubscriberQuotaUpdate(
                  SubscriberQuotaUpdate_Type_TERMINATE)))
      .Times(1);

  local_enforcer->update_session_credits_and_rules(session_map, update_response,
                                                   update);
}

int main(int argc, char** argv) {
  ::testing::InitGoogleTest(&argc, argv);
  FLAGS_logtostderr = 1;
  FLAGS_v = 10;
  return RUN_ALL_TESTS();
}

}  // namespace magma

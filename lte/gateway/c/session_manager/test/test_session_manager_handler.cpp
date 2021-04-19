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
#include "MagmaService.h"
#include "Matchers.h"
#include "MeteringReporter.h"
#include "ProtobufCreators.h"
#include "RuleStore.h"
#include "ServiceRegistrySingleton.h"
#include "SessiondMocks.h"
#include "SessionState.h"
#include "SessionStore.h"
#include "StoredState.h"

using ::testing::Test;
#define DEFAULT_PIPELINED_EPOCH 1

namespace magma {

Teids teids0;
Teids teids1;

ACTION_P(CallSetupCallback, result) {
  auto cb =
      static_cast<std::function<void(Status status, SetupFlowsResult)>>(arg2);
  SetupFlowsResult setup_result;
  setup_result.set_result(result);
  cb(grpc::Status::OK, setup_result);
}

class SessionManagerHandlerTest : public ::testing::Test {
 protected:
  virtual void SetUp() {
    monitoring_key = "mk1";

    reporter      = std::make_shared<MockSessionReporter>();
    rule_store    = std::make_shared<StaticRuleStore>();
    session_store = std::make_shared<SessionStore>(
        rule_store, std::make_shared<MeteringReporter>());
    pipelined_client       = std::make_shared<MockPipelinedClient>();
    auto directoryd_client = std::make_shared<MockDirectorydClient>();
    auto spgw_client       = std::make_shared<MockSpgwServiceClient>();
    auto aaa_client        = std::make_shared<MockAAAClient>();
    events_reporter        = std::make_shared<MockEventsReporter>();
    auto default_mconfig   = get_default_mconfig();
    local_enforcer         = std::make_shared<LocalEnforcer>(
        reporter, rule_store, *session_store, pipelined_client, events_reporter,
        spgw_client, aaa_client, 0, 0, default_mconfig);
    evb = new folly::EventBase();
    std::thread([&]() {
      std::cout << "Started event loop thread\n";
      folly::EventBaseManager::get()->setEventBase(evb, 0);
    })
        .detach();

    local_enforcer->attachEventBase(evb);
    session_map_ = SessionMap{};

    session_manager = std::make_shared<LocalSessionManagerHandlerImpl>(
        local_enforcer, reporter.get(), directoryd_client, events_reporter,
        *session_store);

    teids0.set_agw_teid(0);
    teids0.set_enb_teid(0);
    teids1.set_agw_teid(TEID_1_UL);
    teids1.set_enb_teid(TEID_1_DL);
  }

  virtual void TearDown() {
    local_enforcer->stop();
    delete evb;
  }

  void insert_static_rule(
      std::shared_ptr<StaticRuleStore> rule_store, const std::string& m_key,
      uint32_t charging_key, const std::string& rule_id) {
    PolicyRule rule;
    rule.set_id(rule_id);
    rule.set_rating_group(charging_key);
    rule.set_monitoring_key(m_key);
    rule.set_tracking_type(PolicyRule::OCS_AND_PCRF);
    rule_store->insert_rule(rule);
  }

  void send_empty_table() {
    RuleRecordTable empty_table;
    // epoch indicates the last PipelineD service start time
    empty_table.set_epoch(DEFAULT_PIPELINED_EPOCH);
    grpc::ServerContext context;
    session_manager->ReportRuleStats(
        &context, &empty_table,
        [this](grpc::Status status, Void response_out) {});
  }

  // This function should always be called at the beginning of the test to
  // prevent unexpected SessionD <-> PipelineD syncing logic mid-test.
  void send_empty_table_and_wait_for_successful_setup() {
    send_empty_table();
    EXPECT_CALL(
        *pipelined_client, setup_lte(testing::_, testing::_, testing::_))
        .Times(1)
        .WillOnce(CallSetupCallback(SetupFlowsResult_Result_SUCCESS));
    evb->loopOnce();
    evb->loopOnce();
  }

 protected:
  std::string monitoring_key;

  std::shared_ptr<SessionStore> session_store;
  std::shared_ptr<StaticRuleStore> rule_store;
  std::shared_ptr<LocalSessionManagerHandlerImpl> session_manager;
  std::shared_ptr<MockSessionReporter> reporter;
  std::shared_ptr<MockPipelinedClient> pipelined_client;
  std::shared_ptr<LocalEnforcer> local_enforcer;
  std::shared_ptr<MockEventsReporter> events_reporter;
  SessionIDGenerator id_gen_;
  folly::EventBase* evb;
  SessionMap session_map_;
};

TEST_F(SessionManagerHandlerTest, test_create_session_cfg) {
  send_empty_table_and_wait_for_successful_setup();
  // 1) Insert the entry for a rule
  insert_static_rule(rule_store, monitoring_key, 1, "rule1");
  std::vector<std::string> static_rules{"rule1"};

  LocalCreateSessionRequest request;
  CreateSessionResponse response;

  const std::string& hardware_addr_bytes = {0x0f, 0x10, 0x2e, 0x12, 0x3a, 0x55};
  SessionConfig cfg;
  cfg.common_context =
      build_common_context(IMSI1, "", "", teids0, "apn1", MSISDN, TGPP_WLAN);
  const auto& wlan = build_wlan_context(MAC_ADDR, RADIUS_SESSION_ID);
  cfg.rat_specific_context.mutable_wlan_context()->CopyFrom(wlan);

  response.set_session_id(SESSION_ID_1);
  // Only the active sessions are not recycled, to ensure that
  // this session is not automatically scheduled for termination
  // when RAT Type is WLAN, it needs monitoring keys...
  create_session_create_response(
      IMSI1, SESSION_ID_1, monitoring_key, static_rules, &response);
  response.mutable_static_rules()->Add()->mutable_rule_id()->assign("rule1");
  create_credit_update_response(
      IMSI1, SESSION_ID_1, 1, 1536, response.mutable_credits()->Add());

  auto session_map = session_store->read_sessions({IMSI1});
  local_enforcer->init_session(session_map, IMSI1, SESSION_ID_1, cfg, response);
  local_enforcer->update_tunnel_ids(
      session_map,
      create_update_tunnel_ids_request(IMSI1, BEARER_ID_1, teids0));

  bool write_success =
      session_store->create_sessions(IMSI1, std::move(session_map[IMSI1]));
  EXPECT_TRUE(write_success);
  session_map = session_store->read_sessions({IMSI1});
  auto it     = session_map.find(IMSI1);
  EXPECT_FALSE(it == session_map.end());
  EXPECT_EQ(session_map[IMSI1].size(), 1);
  auto& session = session_map[IMSI1][0];
  EXPECT_EQ(session->get_config().common_context.apn(), "apn1");

  grpc::ServerContext create_context;

  auto common =
      build_common_context(IMSI1, "", "", teids0, "apn2", MSISDN, TGPP_WLAN);
  request.mutable_common_context()->CopyFrom(common);
  request.mutable_rat_specific_context()->mutable_wlan_context()->CopyFrom(
      wlan);  // use same WLAN config as previous

  // Ensure session is not reported as its a duplicate
  EXPECT_CALL(*reporter, report_create_session(_, _)).Times(0);
  session_manager->CreateSession(
      &create_context, &request,
      [this](grpc::Status status, LocalCreateSessionResponse response_out) {});
  local_enforcer->update_tunnel_ids(
      session_map,
      create_update_tunnel_ids_request(IMSI1, BEARER_ID_1, teids0));

  // Run session creation in the EventBase loop
  // It needs to loop once here.
  evb->loopOnce();

  // Assert the internal session config is updated to the new one
  session_map = session_store->read_sessions({IMSI1});
  it          = session_map.find(IMSI1);
  EXPECT_FALSE(it == session_map.end());
  EXPECT_EQ(session_map[IMSI1].size(), 1);
  auto& session_apn2 = session_map[IMSI1][0];
  EXPECT_EQ(session_apn2->get_config().common_context.apn(), "apn2");
}

TEST_F(SessionManagerHandlerTest, test_session_recycling_lte) {
  send_empty_table_and_wait_for_successful_setup();
  // 1) Insert the entry for a rule
  insert_static_rule(rule_store, monitoring_key, 1, "rule1");
  std::vector<std::string> static_rules{"rule1"};

  CreateSessionResponse response;
  auto sid = id_gen_.gen_session_id(IMSI1);

  SessionConfig cfg;
  cfg.common_context =
      build_common_context(IMSI1, IP1, IPv6_1, teids1, APN1, MSISDN, TGPP_LTE);
  auto lte_context = build_lte_context(
      "spgw_ip", "imei", "plmn_id", "imsi_plmn_id", "user_loc", BEARER_ID_1,
      nullptr);
  cfg.rat_specific_context.mutable_lte_context()->CopyFrom(lte_context);

  response.set_session_id(sid);
  create_session_create_response(
      IMSI1, sid, monitoring_key, static_rules, &response);
  response.mutable_static_rules()->Add()->mutable_rule_id()->assign("rule1");
  create_credit_update_response(
      IMSI1, sid, 1, 1536, response.mutable_credits()->Add());

  auto session_map = session_store->read_sessions({IMSI1});

  local_enforcer->init_session(session_map, IMSI1, sid, cfg, response);
  local_enforcer->update_tunnel_ids(
      session_map,
      create_update_tunnel_ids_request(IMSI1, BEARER_ID_1, teids1));

  bool write_success =
      session_store->create_sessions(IMSI1, std::move(session_map[IMSI1]));
  EXPECT_TRUE(write_success);
  session_map = session_store->read_sessions({IMSI1});
  auto it     = session_map.find(IMSI1);
  EXPECT_FALSE(it == session_map.end());
  EXPECT_EQ(session_map[IMSI1].size(), 1);
  auto& session = session_map[IMSI1][0];
  EXPECT_EQ(session->get_config().common_context.apn(), APN1);

  // Only active, identical sessions can be recycled for LTE
  // The previously created session is active and this request has the same
  // context
  LocalCreateSessionRequest request;
  grpc::ServerContext create_context;
  auto common =
      build_common_context(IMSI1, IP1, IPv6_1, teids1, APN1, MSISDN, TGPP_LTE);
  request.mutable_common_context()->CopyFrom(common);
  lte_context = build_lte_context(
      "spgw_ip", "imei", "plmn_id", "imsi_plmn_id", "user_loc", BEARER_ID_1,
      nullptr);
  request.mutable_rat_specific_context()->mutable_lte_context()->CopyFrom(
      lte_context);

  // Ensure session is not reported as its a duplicate
  EXPECT_CALL(*reporter, report_create_session(_, _)).Times(0);

  // Run session creation in the EventBase loop
  // It needs to loop once here.
  evb->loopOnce();

  // Assert the internal session config is updated to the new one
  session_map = session_store->read_sessions({IMSI1});
  it          = session_map.find(IMSI1);
  EXPECT_FALSE(it == session_map.end());
  EXPECT_EQ(session_map[IMSI1].size(), 1);
  auto& session_apn2 = session_map[IMSI1][0];
  EXPECT_EQ(session_apn2->get_config().common_context.apn(), APN1);

  // Now make the config not identical but with the same APN=apn1, this should
  // trigger a terminate for the existing and a creation for the new session
  Teids teids2;
  teids2.set_enb_teid(TEID_2_UL);
  teids2.set_agw_teid(TEID_2_DL);
  LocalCreateSessionRequest request2;
  grpc::ServerContext create_context2;
  common = build_common_context(
      IMSI1, "", "", teids2, APN1, "different msisdn", TGPP_LTE);
  request2.mutable_common_context()->CopyFrom(common);
  lte_context = build_lte_context(
      "spgw_ip", "imei", "plmn_id", "imsi_plmn_id", "user_loc", BEARER_ID_1,
      nullptr);
  request2.mutable_rat_specific_context()->mutable_lte_context()->CopyFrom(
      lte_context);

  // Ensure a create session for the new session is sent, the old one is
  // terminated
  EXPECT_CALL(*reporter, report_create_session(_, _)).Times(1);

  session_manager->CreateSession(
      &create_context2, &request2,
      [this](grpc::Status status, LocalCreateSessionResponse response_out) {});

  // Run session creation in the EventBase loop
  // It needs to loop once here.
  evb->loopOnce();
  local_enforcer->update_tunnel_ids(
      session_map,
      create_update_tunnel_ids_request(IMSI1, BEARER_ID_1, teids2));
}

TEST_F(SessionManagerHandlerTest, test_create_session) {
  send_empty_table_and_wait_for_successful_setup();
  // 1) Create the session
  LocalCreateSessionRequest request;

  grpc::ServerContext server_context;
  request.mutable_common_context()->mutable_sid()->set_id(IMSI1);
  request.mutable_common_context()->set_rat_type(RATType::TGPP_LTE);
  request.mutable_common_context()->set_msisdn(MSISDN);

  CreateSessionResponse create_response;
  create_response.mutable_static_rules()->Add()->mutable_rule_id()->assign(
      "rule1");
  create_response.mutable_static_rules()->Add()->mutable_rule_id()->assign(
      "rule2");
  create_response.mutable_static_rules()->Add()->mutable_rule_id()->assign(
      "rule3");
  create_credit_update_response(
      IMSI1, "1234", 1, 1536, create_response.mutable_credits()->Add());
  create_credit_update_response(
      IMSI1, "1234", 2, 1024, create_response.mutable_credits()->Add());

  // create expected request for report_create_session call
  RequestedUnits expected_requestedUnits;
  expected_requestedUnits.set_total(SessionCredit::DEFAULT_REQUESTED_UNITS);
  expected_requestedUnits.set_rx(SessionCredit::DEFAULT_REQUESTED_UNITS);
  expected_requestedUnits.set_tx(SessionCredit::DEFAULT_REQUESTED_UNITS);
  CreateSessionRequest expected_request;
  expected_request.mutable_requested_units()->CopyFrom(expected_requestedUnits);
  expected_request.mutable_common_context()->CopyFrom(request.common_context());
  expected_request.mutable_rat_specific_context()->CopyFrom(
      request.rat_specific_context());

  EXPECT_CALL(
      *reporter, report_create_session(CheckCoreRequest(expected_request), _))
      .Times(1);

  // create session and expect one call
  session_manager->CreateSession(
      &server_context, &request,
      [this](grpc::Status status, LocalCreateSessionResponse response_out) {});

  // Run session creation in the EventBase loop
  evb->loopOnce();
  evb->loopOnce();
  evb->loopOnce();
}

TEST_F(SessionManagerHandlerTest, test_create_session_pipelined_unavailable) {
  send_empty_table();
  // On failure cases, LocalEnforcer will endlessly retry the setup call
  EXPECT_CALL(*pipelined_client, setup_lte(testing::_, testing::_, testing::_))
      .WillRepeatedly(CallSetupCallback(SetupFlowsResult_Result_FAILURE));
  evb->loopOnce();
  evb->loopOnce();
  // 1) Create the session
  LocalCreateSessionRequest request;

  grpc::ServerContext server_context;
  request.mutable_common_context()->mutable_sid()->set_id(IMSI1);
  request.mutable_common_context()->set_rat_type(RATType::TGPP_LTE);
  request.mutable_common_context()->set_msisdn(MSISDN);

  // create session and expect one call
  session_manager->CreateSession(
      &server_context, &request,
      [this](grpc::Status status, LocalCreateSessionResponse response_out) {
        EXPECT_FALSE(status.ok());
      });

  // Run session creation in the EventBase loop
  evb->loopOnce();
}

TEST_F(SessionManagerHandlerTest, test_report_rule_stats) {
  // 1) Insert the entry for a rule
  insert_static_rule(rule_store, monitoring_key, 1, "rule1");

  // 2) Create a session
  Teids teids;
  teids.set_enb_teid(TEID_1_UL);
  teids.set_agw_teid(TEID_1_DL);
  CreateSessionResponse response;
  response.mutable_static_rules()->Add()->mutable_rule_id()->assign("rule1");
  create_credit_update_response(
      IMSI1, SESSION_ID_1, 1, 1025, response.mutable_credits()->Add());
  SessionConfig cfg = {};
  cfg.common_context =
      build_common_context(IMSI1, IP1, IPv6_1, teids, "APN", MSISDN, TGPP_LTE);
  const auto& lte_context = build_lte_context(
      "127.0.0.1", "imei", "plmn_id", "imsi_plmn_id", "user_loc", BEARER_ID_1,
      nullptr);
  cfg.rat_specific_context.mutable_lte_context()->CopyFrom(lte_context);
  auto session_map = session_store->read_sessions({IMSI1});
  EXPECT_CALL(
      *events_reporter,
      session_created(IMSI1, SESSION_ID_1, testing::_, testing::_))
      .Times(1);

  local_enforcer->init_session(session_map, IMSI1, SESSION_ID_1, cfg, response);
  local_enforcer->update_tunnel_ids(
      session_map,
      create_update_tunnel_ids_request(IMSI1, BEARER_ID_1, teids0));

  bool write_success =
      session_store->create_sessions(IMSI1, std::move(session_map[IMSI1]));
  EXPECT_TRUE(write_success);

  // Check the request number
  auto session_map_2 = session_store->read_sessions(SessionRead{IMSI1});
  EXPECT_EQ(session_map_2[IMSI1].front()->get_request_number(), 1);
  // 2) ReportRuleStats
  grpc::ServerContext server_context;
  RuleRecordTable table;
  auto record_list = table.mutable_records();
  create_rule_record(
      IMSI1, teids0.agw_teid(), "rule1", 512, 512, record_list->Add());

  EXPECT_CALL(
      *reporter, report_updates(CheckUpdateRequestNumber(1), testing::_))
      .Times(1);
  session_manager->ReportRuleStats(
      &server_context, &table,
      [this](grpc::Status status, orc8r::Void response_out) {});
  evb->loopOnce();

  session_map_2 = session_store->read_sessions(SessionRead{IMSI1});
  EXPECT_EQ(session_map_2[IMSI1].front()->get_request_number(), 2);
  evb->loopOnce();
}

TEST_F(SessionManagerHandlerTest, test_end_session) {
  // 1) Insert the entry for a rule
  insert_static_rule(rule_store, monitoring_key, 1, "rule1");

  // 2) Create a session
  Teids teids;
  CreateSessionResponse response;
  response.mutable_static_rules()->Add()->mutable_rule_id()->assign("rule1");
  create_credit_update_response(
      IMSI1, SESSION_ID_1, 1, 1025, response.mutable_credits()->Add());
  const std::string& hardware_addr_bytes = {0x0f, 0x10, 0x2e, 0x12, 0x3a, 0x55};
  const std::string& apn                 = "apn1";
  SessionConfig cfg;
  cfg.common_context =
      build_common_context(IMSI1, "", "", teids, apn, MSISDN, TGPP_WLAN);
  const auto& wlan = build_wlan_context(MAC_ADDR, RADIUS_SESSION_ID);
  cfg.rat_specific_context.mutable_wlan_context()->CopyFrom(wlan);

  auto session_map = session_store->read_sessions({IMSI1});

  local_enforcer->init_session(session_map, IMSI1, SESSION_ID_1, cfg, response);
  local_enforcer->update_tunnel_ids(
      session_map, create_update_tunnel_ids_request(IMSI1, 0, teids0));

  bool write_success =
      session_store->create_sessions(IMSI1, std::move(session_map[IMSI1]));
  EXPECT_TRUE(write_success);

  // 3) EndSession
  session_map = session_store->read_sessions({IMSI1});
  EXPECT_EQ(session_map[IMSI1].size(), 1);
  LocalEndSessionRequest end_request;
  end_request.mutable_sid()->set_id(IMSI1);
  end_request.set_apn("apn1");
  grpc::ServerContext server_context;

  EXPECT_CALL(*reporter, report_terminate_session(_, _)).Times(1);
  session_manager->EndSession(
      &server_context, &end_request,
      [this](grpc::Status status, LocalEndSessionResponse response_out) {});
  evb->loopOnce();
  session_map = session_store->read_sessions({IMSI1});
  EXPECT_EQ(session_map[IMSI1].size(), 1);

  evb->loopOnce();

  session_map = session_store->read_sessions({IMSI1});
  EXPECT_EQ(session_map[IMSI1].size(), 0);
}

int main(int argc, char** argv) {
  ::testing::InitGoogleTest(&argc, argv);
  return RUN_ALL_TESTS();
}

}  // namespace magma

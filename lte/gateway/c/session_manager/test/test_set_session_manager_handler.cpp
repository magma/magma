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

/* session_manager.grpc.pb.h and SessionStateEnforcer.h
 * included in "SetMessageManagerHandler.h"
 */
#include "Matchers.h"
#include "SetMessageManagerHandler.h"
#include "SessionStateEnforcer.h"
#include "UpfMsgManageHandler.h"
#include "MobilitydClient.h"
#include "includes/MagmaService.h"
#include "ProtobufCreators.h"
#include "RuleStore.h"
#include "SessionState.h"
#include "SessionStore.h"
#include "includes/ServiceRegistrySingleton.h"
#include "SessiondMocks.h"
#include "StoredState.h"
#include "magma_logging.h"
#include "PipelinedClient.h"
#include "AmfServiceClient.h"
#include "DirectorydClient.h"
#include "Consts.h"
#include "EnumToString.h"

using grpc::ServerContext;
using grpc::Status;
using ::testing::InSequence;
using ::testing::Test;

#define session_force_termination_timeout_ms 5000
#define session_max_rtx_count 3
namespace magma {

class SessionManagerHandlerTest : public ::testing::Test {
 public:
  virtual void SetUp() {
    rule_store       = std::make_shared<StaticRuleStore>();
    reporter         = std::make_shared<MockSessionReporter>();
    mobilityd_client = std::make_shared<MockMobilitydClient>();
    session_store    = std::make_shared<SessionStore>(
        rule_store, std::make_shared<MeteringReporter>());
    reporter = std::make_shared<MockSessionReporter>();
    std::unordered_multimap<std::string, uint32_t> pdr_map;
    pipelined_client       = std::make_shared<MockPipelinedClient>();
    amf_srv_client         = std::make_shared<magma::MockAmfServiceClient>();
    events_reporter        = std::make_shared<MockEventsReporter>();
    auto shard_tracker     = std::make_shared<ShardTracker>();
    aaa_client             = std::make_shared<MockAAAClient>();
    spgw_client            = std::make_shared<MockSpgwServiceClient>();
    auto directoryd_client = std::make_shared<MockDirectorydClient>();
    magma::mconfig::SessionD mconfig;
    mconfig.set_log_level(magma::orc8r::LogLevel::INFO);
    auto default_mconfig = get_default_mconfig();

    pdr_map.insert(std::pair<std::string, uint32_t>(IMSI1, 1));
    pdr_map.insert(std::pair<std::string, uint32_t>(IMSI1, 2));

    session_enforcer = std::make_shared<magma::SessionStateEnforcer>(
        rule_store, *session_store, pdr_map, pipelined_client, amf_srv_client,
        reporter.get(), events_reporter, mconfig,
        session_force_termination_timeout_ms, session_max_rtx_count);

    evb = new folly::EventBase();
    std::thread([&]() {
      std::cout << "Started event loop thread\n";
      folly::EventBaseManager::get()->setEventBase(evb, 0);
    })
        .detach();
    session_enforcer->attachEventBase(evb);
    session_map_ = SessionMap{};
    // creating landing object and invoking contructor
    set_session_manager = std::make_shared<SetMessageManagerHandler>(
        session_enforcer, *session_store, reporter.get(), events_reporter);

    upf_session_manager = std::make_shared<UpfMsgManageHandler>(
        session_enforcer, mobilityd_client, *session_store);

    local_enforcer = std::make_unique<LocalEnforcer>(
        reporter, rule_store, *session_store, pipelined_client, events_reporter,
        spgw_client, aaa_client, shard_tracker, 0, 0, default_mconfig);
  }

  virtual void TearDown() { delete evb; }

  void insert_static_rule(
      uint32_t rating_group, const std::string& m_key,
      const std::string& rule_id) {
    rule_store->insert_rule(create_policy_rule(rule_id, m_key, rating_group));
  }

  void set_sm_session_context(magma::SetSMSessionContext* request) {
    auto* req = request->mutable_rat_specific_context()
                    ->mutable_m5gsm_session_context();
    auto* reqcmn = request->mutable_common_context();
    req->set_pdu_session_id({0x5});
    req->set_request_type(magma::RequestType::INITIAL_REQUEST);
    req->set_priority_access(magma::priorityaccess::High);
    req->set_access_type(magma::AccessType::M_3GPP_ACCESS_3GPP);
    req->set_imei("123456789012345");
    req->set_gpsi("9876543210");
    req->set_pcf_id("1357924680123456");

    reqcmn->mutable_sid()->set_id("IMSI000000000000001");
    reqcmn->set_sm_session_state(magma::SMSessionFSMState::CREATING_0);
    reqcmn->set_rat_type(magma::RATType::TGPP_NR);
    reqcmn->set_apn(APN1);
    reqcmn->set_msisdn(MSISDN);
    reqcmn->mutable_teids()->set_enb_teid(TEID_1_UL);
    reqcmn->mutable_teids()->set_agw_teid(TEID_1_DL);
    EXPECT_EQ(req->pdu_session_id(), 0x5);
    EXPECT_EQ(req->request_type(), magma::RequestType::INITIAL_REQUEST);
    EXPECT_EQ(req->priority_access(), magma::priorityaccess::High);
    EXPECT_EQ(req->access_type(), magma::AccessType::M_3GPP_ACCESS_3GPP);
    EXPECT_EQ(req->imei(), "123456789012345");
    EXPECT_EQ(req->gpsi(), "9876543210");
    EXPECT_EQ(req->pcf_id(), "1357924680123456");

    EXPECT_EQ(reqcmn->sid().id(), "IMSI000000000000001");
    EXPECT_EQ(reqcmn->sm_session_state(), magma::SMSessionFSMState::CREATING_0);
  }

  void set_sm_session_context_ipv4(magma::SetSMSessionContext* request) {
    auto* reqcmn = request->mutable_common_context();
    auto* req    = request->mutable_rat_specific_context()
                    ->mutable_m5gsm_session_context();
    req->set_pdu_session_type(magma::PduSessionType::IPV4);
    reqcmn->set_ue_ipv4(IPv4_1);
    set_sm_session_context(request);
  }

  void set_sm_session_context_ipv6(magma::SetSMSessionContext* request) {
    auto* reqcmn = request->mutable_common_context();
    auto* req    = request->mutable_rat_specific_context()
                    ->mutable_m5gsm_session_context();
    req->set_pdu_session_type(magma::PduSessionType::IPV6);
    reqcmn->set_ue_ipv6(IPv6_5);
    set_sm_session_context(request);
  }

  void set_sm_notif_context(
      magma::SetSmNotificationContext* request,
      magma::SetSMSessionContext* session_ctx_req) {
    auto* req    = request->mutable_rat_specific_notification();
    auto* reqcmn = request->mutable_common_context();
    req->set_pdu_session_id({0x5});
    req->set_request_type(magma::RequestType::INITIAL_REQUEST);
    req->set_access_type(magma::AccessType::M_3GPP_ACCESS_3GPP);
    req->set_notify_ue_event(
        magma::NotifyUeEvents::PDU_SESSION_INACTIVE_NOTIFY);

    reqcmn->mutable_sid()->set_id("IMSI000000000000001");
    reqcmn->set_apn("BLR");
    reqcmn->set_ue_ipv4("192.168.128.11");
    reqcmn->set_rat_type(magma::RATType::TGPP_NR);
    reqcmn->set_sm_session_state(magma::SMSessionFSMState::CREATING_0);

    auto imsi   = reqcmn->sid().id();
    auto apn    = reqcmn->apn();
    auto pdu_id = req->pdu_session_id();
    auto cfg = set_session_manager->m5g_build_session_config(*session_ctx_req);
    auto& noti    = *req;
    auto ue_event = noti.notify_ue_event();
    std::function<void(Status, SmContextVoid)> response_callback;

    EXPECT_EQ(pdu_id, 0x5);
    EXPECT_EQ(imsi, "IMSI000000000000001");
    EXPECT_EQ(apn, "BLR");
    EXPECT_EQ(req->request_type(), magma::RequestType::INITIAL_REQUEST);
    EXPECT_EQ(req->access_type(), magma::AccessType::M_3GPP_ACCESS_3GPP);

    EXPECT_EQ(reqcmn->ue_ipv4(), "192.168.128.11");
    EXPECT_EQ(reqcmn->rat_type(), magma::RATType::TGPP_NR);
    set_session_manager->send_create_session(session_map_, imsi, cfg, pdu_id);

    EXPECT_EQ(reqcmn->sm_session_state(), magma::SMSessionFSMState::CREATING_0);
    EXPECT_EQ(ue_event, magma::NotifyUeEvents::PDU_SESSION_INACTIVE_NOTIFY);

    req->set_notify_ue_event(magma::NotifyUeEvents::UE_IDLE_MODE_NOTIFY);
    EXPECT_EQ(
        noti.notify_ue_event(), magma::NotifyUeEvents::UE_IDLE_MODE_NOTIFY);

    set_session_manager->initiate_release_session(session_map_, pdu_id, imsi);
    reqcmn->set_sm_session_state(magma::SMSessionFSMState::RELEASED_4);
    EXPECT_EQ(reqcmn->sm_session_state(), magma::SMSessionFSMState::RELEASED_4);
    EXPECT_EQ(reqcmn->sm_session_version(), 0);
  }

 public:
  std::shared_ptr<SessionStore> session_store;
  std::shared_ptr<SetMessageManagerHandler> set_session_manager;
  std::shared_ptr<MockPipelinedClient> pipelined_client;
  std::shared_ptr<SessionStateEnforcer> session_enforcer;
  std::shared_ptr<MockAmfServiceClient> amf_srv_client;
  std::shared_ptr<StaticRuleStore> rule_store;
  SessionIDGenerator id_gen_;
  folly::EventBase* evb;
  SessionMap session_map_;
  std::unordered_multimap<std::string, uint32_t> pdr_map_;
  std::shared_ptr<MockSessionReporter> reporter;
  std::shared_ptr<MockEventsReporter> events_reporter;
  std::shared_ptr<UpfMsgManageHandler> upf_session_manager;
  std::shared_ptr<MockMobilitydClient> mobilityd_client;
  std::unique_ptr<LocalEnforcer> local_enforcer;
  std::shared_ptr<MockSpgwServiceClient> spgw_client;
  std::shared_ptr<MockAAAClient> aaa_client;
  std::shared_ptr<LocalSessionManagerHandlerImpl> session_manager;
};  // End of class

TEST_F(SessionManagerHandlerTest, test_SetAmfSessionContext) {
  magma::SetSMSessionContext request;
  set_sm_session_context_ipv4(&request);

  grpc::ServerContext server_context;

  set_session_manager->SetAmfSessionContext(
      &server_context, &request,
      [this](grpc::Status status, SmContextVoid Void) {});

  // Run session creation in the EventBase loop
  evb->loopOnce();

  auto session_map = session_store->read_sessions({IMSI1});
  auto it          = session_map.find(IMSI1);
  EXPECT_FALSE(it == session_map.end());
  EXPECT_EQ(session_map[IMSI1].size(), 1);

  auto& session_temp = session_map[IMSI1][0];
  EXPECT_EQ(session_temp->get_config().common_context.sid().id(), IMSI1);
  EXPECT_EQ(session_temp->get_config().common_context.ue_ipv4(), IPv4_1);
  EXPECT_EQ(session_temp->get_config().common_context.apn(), APN1);
  EXPECT_EQ(session_temp->get_config().common_context.msisdn(), MSISDN);
  EXPECT_EQ(
      session_temp->get_config().common_context.teids().enb_teid(), TEID_1_UL);
  EXPECT_EQ(
      session_temp->get_config().common_context.teids().agw_teid(), TEID_1_DL);
  EXPECT_EQ(
      session_temp->get_config().common_context.rat_type(),
      magma::RATType::TGPP_NR);
}

TEST_F(SessionManagerHandlerTest, test_SetSmfNotification) {
  magma::SetSMSessionContext session_ctx_req;
  set_sm_session_context_ipv4(&session_ctx_req);

  magma::SetSmNotificationContext request;
  set_sm_notif_context(&request, &session_ctx_req);

  grpc::ServerContext server_context;

  set_session_manager->SetSmfNotification(
      &server_context, &request,
      [this](grpc::Status status, SmContextVoid Void) {});

  EXPECT_EQ(session_ctx_req.mutable_common_context()->ue_ipv4(), IPv4_1);
  EXPECT_EQ(session_ctx_req.mutable_common_context()->sid().id(), IMSI1);
  EXPECT_EQ(session_ctx_req.mutable_common_context()->apn(), APN1);
  EXPECT_EQ(session_ctx_req.mutable_common_context()->msisdn(), MSISDN);
  EXPECT_EQ(
      session_ctx_req.mutable_common_context()->teids().enb_teid(), TEID_1_UL);
  EXPECT_EQ(
      session_ctx_req.mutable_common_context()->teids().agw_teid(), TEID_1_DL);
  EXPECT_EQ(
      session_ctx_req.mutable_common_context()->rat_type(),
      magma::RATType::TGPP_NR);
  /* Validating the funcationality when specific IMSi and its associated
   * sessions moved to idle mode
   */
  set_session_manager->idle_mode_change_sessions_handle(
      request, [](grpc::Status status, SmContextVoid Void) {});

  /* Validating the functionality when any specific  session moved to
   * idle state
   */
  set_session_manager->pdu_session_inactive(
      request, [](grpc::Status status, SmContextVoid Void) {});

  /* Validating the functionality when specific IMSI service
   * request is received from AMF.
   */
  set_session_manager->service_handle_request_on_paging(
      request, [](grpc::Status status, SmContextVoid Void) {});

  // Run session creation in the EventBase loop
  evb->loopOnce();
}

TEST_F(SessionManagerHandlerTest, test_InitSessionContext) {
  magma::SetSMSessionContext request;
  set_sm_session_context_ipv4(&request);

  grpc::ServerContext server_context;

  set_session_manager->SetAmfSessionContext(
      &server_context, &request,
      [this](grpc::Status status, SmContextVoid Void) {});

  // Run session creation in the EventBase loop
  evb->loopOnce();
  SessionConfig cfg;
  cfg.common_context       = request.common_context();
  cfg.rat_specific_context = request.rat_specific_context();
  cfg.rat_specific_context.mutable_m5gsm_session_context()->set_ssc_mode(
      SSC_MODE_3);

  auto session_map       = session_store->read_sessions({IMSI1});
  std::string session_id = id_gen_.gen_session_id(IMSI1);
  session_enforcer->m5g_init_session_credit(
      session_map, IMSI1, session_id, cfg);
  EXPECT_EQ(request.mutable_common_context()->ue_ipv4(), IPv4_1);
}

TEST_F(SessionManagerHandlerTest, test_UpdateSessionContext) {
  magma::SetSMSessionContext request;
  set_sm_session_context_ipv4(&request);

  grpc::ServerContext server_context;

  set_session_manager->SetAmfSessionContext(
      &server_context, &request,
      [this](grpc::Status status, SmContextVoid Void) {});

  // Run session creation in the EventBase loop
  evb->loopOnce();
  SessionConfig cfg;
  cfg.common_context       = request.common_context();
  cfg.rat_specific_context = request.rat_specific_context();
  cfg.rat_specific_context.mutable_m5gsm_session_context()->set_ssc_mode(
      SSC_MODE_3);

  auto session_map       = session_store->read_sessions({IMSI1});
  std::string session_id = id_gen_.gen_session_id(IMSI1);
  SessionUpdate update = SessionStore::get_default_session_update(session_map);
  uint32_t pdu_id      = 5;
  SessionSearchCriteria id1_success_sid(IMSI1, IMSI_AND_PDUID, pdu_id);
  auto session_it = session_store->find_session(session_map, id1_success_sid);
  auto& session   = **session_it;
  session->set_config(cfg, nullptr);
  session_enforcer->add_default_rules(session, IMSI1);
  session_enforcer->m5g_update_session_context(
      session_map, IMSI1, session, update);
  EXPECT_EQ(request.mutable_common_context()->ue_ipv4(), IPv4_1);
}

TEST_F(SessionManagerHandlerTest, test_SetPduSessionReleaseContext) {
  magma::SetSMSessionContext request;
  auto* req =
      request.mutable_rat_specific_context()->mutable_m5gsm_session_context();
  auto* reqcmn = request.mutable_common_context();
  req->set_pdu_session_id({0x5});
  req->set_request_type(magma::RequestType::INITIAL_REQUEST);
  req->set_pdu_session_type(magma::PduSessionType::IPV4);
  req->set_priority_access(magma::priorityaccess::High);
  req->set_imei("123456789012345");
  req->set_gpsi("9876543210");
  req->set_pcf_id("1357924680123456");

  reqcmn->mutable_sid()->set_id("IMSI000000000000001");
  reqcmn->set_sm_session_state(magma::SMSessionFSMState::CREATING_0);
  reqcmn->set_rat_type(magma::RATType::TGPP_NR);
  reqcmn->set_ue_ipv4(IPv4_1);

  grpc::ServerContext server_context;

  set_session_manager->SetAmfSessionContext(
      &server_context, &request,
      [this](grpc::Status status, SmContextVoid Void) {});

  // Run session creation in the EventBase loop
  evb->loopOnce();
  evb->loopOnce();

  auto session_map = session_store->read_sessions({IMSI1});
  auto it          = session_map.find(IMSI1);
  EXPECT_FALSE(it == session_map.end());
  EXPECT_EQ(session_map[IMSI1].size(), 1);

  auto& session_temp = session_map[IMSI1][0];
  EXPECT_EQ(session_temp->get_config().common_context.sid().id(), IMSI1);

  req->set_request_type(magma::RequestType::EXISTING_PDU_SESSION);
  reqcmn->set_sm_session_state(magma::SMSessionFSMState::RELEASED_4);
  reqcmn->set_sm_session_version(2);

  set_session_manager->SetAmfSessionContext(
      &server_context, &request,
      [this](grpc::Status status, SmContextVoid Void) {});

  // Run session creation in the EventBase loop
  evb->loopOnce();
}

TEST_F(SessionManagerHandlerTest, test_LocalReleaseSessionContext) {
  magma::SetSMSessionContext request;
  auto* req =
      request.mutable_rat_specific_context()->mutable_m5gsm_session_context();
  auto* reqcmn = request.mutable_common_context();
  req->set_pdu_session_id({0x5});
  req->set_request_type(magma::RequestType::INITIAL_REQUEST);
  req->set_pdu_session_type(magma::PduSessionType::IPV4);
  req->set_priority_access(magma::priorityaccess::High);
  req->set_imei("123456789012345");
  req->set_gpsi("9876543210");
  req->set_pcf_id("1357924680123456");

  reqcmn->mutable_sid()->set_id("IMSI000000000000001");
  reqcmn->set_sm_session_state(magma::SMSessionFSMState::CREATING_0);
  reqcmn->set_rat_type(magma::RATType::TGPP_NR);
  reqcmn->set_ue_ipv4(IPv4_1);

  grpc::ServerContext server_context;

  set_session_manager->SetAmfSessionContext(
      &server_context, &request,
      [this](grpc::Status status, SmContextVoid Void) {});

  // Run session creation in the EventBase loop
  evb->loopOnce();
  evb->loopOnce();

  auto session_map = session_store->read_sessions({IMSI1});
  auto it          = session_map.find(IMSI1);
  EXPECT_FALSE(it == session_map.end());
  EXPECT_EQ(session_map[IMSI1].size(), 1);

  SessionConfig cfg;
  cfg.common_context       = request.common_context();
  cfg.rat_specific_context = request.rat_specific_context();
  cfg.rat_specific_context.mutable_m5gsm_session_context()->set_ssc_mode(
      SSC_MODE_3);

  SessionUpdate session_update =
      SessionStore::get_default_session_update(session_map);
  uint32_t pdu_id = 5;

  session_enforcer->m5g_release_session(
      session_map, IMSI1, pdu_id, session_update);
  EXPECT_FALSE(it == session_map.end());
  EXPECT_EQ(session_map[IMSI1].size(), 1);
  auto& session_temp = session_map[IMSI1][0];
  EXPECT_EQ(session_temp->get_config().common_context.sid().id(), IMSI1);
}

TEST_F(SessionManagerHandlerTest, test_LocalSessionTerminationContext) {
  magma::SetSMSessionContext request;
  auto* req =
      request.mutable_rat_specific_context()->mutable_m5gsm_session_context();
  auto* reqcmn = request.mutable_common_context();
  req->set_pdu_session_id({0x5});
  req->set_request_type(magma::RequestType::INITIAL_REQUEST);
  req->set_pdu_session_type(magma::PduSessionType::IPV4);
  req->set_priority_access(magma::priorityaccess::High);
  req->set_imei("123456789012345");
  req->set_gpsi("9876543210");
  req->set_pcf_id("1357924680123456");

  reqcmn->mutable_sid()->set_id("IMSI000000000000002");
  reqcmn->set_sm_session_state(magma::SMSessionFSMState::CREATING_0);
  reqcmn->set_rat_type(magma::RATType::TGPP_NR);
  reqcmn->set_ue_ipv4(IPv4_1);

  grpc::ServerContext server_context;

  set_session_manager->SetAmfSessionContext(
      &server_context, &request,
      [this](grpc::Status status, SmContextVoid Void) {});

  // Run session creation in the EventBase loop
  evb->loopOnce();
  evb->loopOnce();

  auto session_map = session_store->read_sessions({IMSI2});
  auto it          = session_map.find(IMSI2);
  EXPECT_FALSE(it == session_map.end());
  EXPECT_EQ(session_map[IMSI2].size(), 1);

  SessionConfig cfg;
  cfg.common_context       = request.common_context();
  cfg.rat_specific_context = request.rat_specific_context();
  cfg.rat_specific_context.mutable_m5gsm_session_context()->set_ssc_mode(
      SSC_MODE_3);

  SessionUpdate session_update =
      SessionStore::get_default_session_update(session_map);
  uint32_t pdu_id = 5;
  SessionSearchCriteria id1_success_sid(IMSI2, IMSI_AND_PDUID, pdu_id);
  auto session_it = session_store->find_session(session_map, id1_success_sid);
  auto& session   = **session_it;
  auto session_id = session->get_session_id();
  SessionStateUpdateCriteria& session_uc = session_update[IMSI2][session_id];

  session_enforcer->m5g_start_session_termination(
      session_map, session, pdu_id, &session_uc);

  EXPECT_FALSE(it == session_map.end());
  EXPECT_EQ(session_map[IMSI2].size(), 1);
  auto& session_temp = session_map[IMSI2][0];
  EXPECT_EQ(session_temp->get_config().common_context.sid().id(), IMSI2);
}

TEST_F(SessionManagerHandlerTest, test_SessionCompleteTerminationContext) {
  magma::SetSMSessionContext request;
  auto* req =
      request.mutable_rat_specific_context()->mutable_m5gsm_session_context();
  auto* reqcmn = request.mutable_common_context();
  req->set_pdu_session_id({0x5});
  req->set_request_type(magma::RequestType::INITIAL_REQUEST);
  req->set_pdu_session_type(magma::PduSessionType::IPV4);
  req->set_priority_access(magma::priorityaccess::High);
  req->set_imei("123456789012345");
  req->set_gpsi("9876543210");
  req->set_pcf_id("1357924680123456");

  reqcmn->mutable_sid()->set_id("IMSI000000000000002");
  reqcmn->set_sm_session_state(magma::SMSessionFSMState::CREATING_0);
  reqcmn->set_rat_type(magma::RATType::TGPP_NR);
  reqcmn->set_ue_ipv4(IPv4_1);

  grpc::ServerContext server_context;
  set_session_manager->SetAmfSessionContext(
      &server_context, &request,
      [this](grpc::Status status, SmContextVoid Void) {});

  // Run session creation in the EventBase loop
  evb->loopOnce();
  evb->loopOnce();

  auto session_map = session_store->read_sessions({IMSI2});
  auto it          = session_map.find(IMSI2);
  EXPECT_FALSE(it == session_map.end());
  EXPECT_EQ(session_map[IMSI2].size(), 1);

  SessionConfig cfg;
  cfg.common_context       = request.common_context();
  cfg.rat_specific_context = request.rat_specific_context();
  cfg.rat_specific_context.mutable_m5gsm_session_context()->set_ssc_mode(
      SSC_MODE_3);

  SessionUpdate session_update =
      SessionStore::get_default_session_update(session_map);
  uint32_t pdu_id = 5;
  SessionSearchCriteria id1_success_sid(IMSI2, IMSI_AND_PDUID, pdu_id);
  auto session_it = session_store->find_session(session_map, id1_success_sid);
  auto& session   = **session_it;
  auto session_id = session->get_session_id();

  session_enforcer->m5g_complete_termination(
      session_map, IMSI2, session_id, session_update);

  EXPECT_EQ(session_map[IMSI2].size(), 0);
}

TEST_F(SessionManagerHandlerTest, test_IPv6PDUSessionEstablishmentContext) {
  magma::SetSMSessionContext request;
  set_sm_session_context_ipv6(&request);

  grpc::ServerContext server_context;

  // create expected request for report_create_session call
  CreateSessionRequest expected_request;
  expected_request.mutable_common_context()->CopyFrom(request.common_context());
  expected_request.mutable_rat_specific_context()->CopyFrom(
      request.rat_specific_context());

  // Create session request towards policydb
  EXPECT_CALL(
      *reporter, report_create_session(CheckSendRequest(expected_request), _))
      .Times(1);

  CreateSessionResponse response;
  set_session_manager->SetAmfSessionContext(
      &server_context, &request,
      [this](grpc::Status status, SmContextVoid Void) {});

  EXPECT_EQ(request.mutable_common_context()->sid().id(), IMSI1);

  // Run session creation in the EventBase loop
  evb->loopOnce();
  SessionConfig cfg;
  cfg.common_context       = request.common_context();
  cfg.rat_specific_context = request.rat_specific_context();
  cfg.rat_specific_context.mutable_m5gsm_session_context()->set_ssc_mode(
      SSC_MODE_3);

  auto session_map = session_store->read_sessions({IMSI1});
  uint32_t pdu_id  = request.mutable_rat_specific_context()
                        ->mutable_m5gsm_session_context()
                        ->pdu_session_id();
  set_session_manager->send_create_session(session_map, IMSI1, cfg, pdu_id);

  auto it = session_map.find(IMSI1);
  EXPECT_FALSE(it == session_map.end());
  EXPECT_EQ(session_map[IMSI1].size(), 1);

  auto& session_temp = session_map[IMSI1][0];
  EXPECT_EQ(session_temp->get_config().common_context.sid().id(), IMSI1);
  EXPECT_EQ(session_temp->get_config().common_context.ue_ipv6(), IPv6_5);
  EXPECT_EQ(session_temp->get_config().common_context.apn(), APN1);
  EXPECT_EQ(session_temp->get_config().common_context.msisdn(), MSISDN);
  EXPECT_EQ(
      session_temp->get_config().common_context.teids().enb_teid(), TEID_1_UL);
  EXPECT_EQ(
      session_temp->get_config().common_context.teids().agw_teid(), TEID_1_DL);
  EXPECT_EQ(
      session_temp->get_config().common_context.rat_type(),
      magma::RATType::TGPP_NR);

  std::string session_id = id_gen_.gen_session_id(IMSI1);
  auto credits           = response.mutable_credits();
  create_credit_update_response(IMSI1, session_id, 1, 1024, credits->Add());
  bool success = session_enforcer->m5g_init_session_credit(
      session_map, IMSI1, session_id, cfg);
  EXPECT_TRUE(success);
}

TEST_F(SessionManagerHandlerTest, test_PDUStateChangeHandling) {
  magma::SetSMSessionContext request;
  set_sm_session_context_ipv4(&request);
  request.mutable_common_context()->mutable_sid()->set_id(
      "IMSI000000000000002");

  grpc::ServerContext server_context;

  set_session_manager->SetAmfSessionContext(
      &server_context, &request,
      [this](grpc::Status status, SmContextVoid Void) {});

  // Run session creation in the EventBase loop
  evb->loopOnce();

  auto session_map = session_store->read_sessions({IMSI2});
  auto it          = session_map.find(IMSI2);
  EXPECT_FALSE(it == session_map.end());
  EXPECT_EQ(session_map[IMSI2].size(), 1);

  SessionConfig cfg;
  cfg.common_context       = request.common_context();
  cfg.rat_specific_context = request.rat_specific_context();
  cfg.rat_specific_context.mutable_m5gsm_session_context()->set_ssc_mode(
      SSC_MODE_3);

  SessionUpdate session_update =
      SessionStore::get_default_session_update(session_map);
  uint32_t pdu_id = 5;
  SessionSearchCriteria id1_success_sid(IMSI2, IMSI_AND_PDUID, pdu_id);
  auto session_it = session_store->find_session(session_map, id1_success_sid);
  auto& session   = **session_it;
  auto session_id = session->get_session_id();
  SessionStateUpdateCriteria& session_uc = session_update[IMSI2][session_id];
  SetSmNotificationContext notif;
  std::string imsi;

  /* pdu_session_inactive() and idle_mode_change_sessions_handle()
   * call flows
   */
  session_enforcer->m5g_move_to_inactive_state(
      imsi.assign(IMSI2), session, notif, &session_uc);

  session_enforcer->set_new_fsm_state_and_increment_version(
      session, INACTIVE, &session_uc);

  session_enforcer->m5g_pdr_rules_change_and_update_upf(
      session, magma::PdrState::IDLE);

  RulesToProcess pending_activation, pending_deactivation;
  session_enforcer->m5g_send_session_request_to_upf(
      session, pending_activation, pending_deactivation);

  /* service_handle_request_on_paging() call flows */
  session_enforcer->m5g_move_to_active_state(session, notif, &session_uc);
  bool gnb_teid_get = false;
  bool upf_teid_get = false;
  session_enforcer->update_session_rules(
      session, gnb_teid_get, upf_teid_get, &session_uc);

  ConvergedRuleStore GlobalRuleList;
  SetGroupPDR rule;
  auto& session_state = session_map[imsi];

  auto itp = pdr_map_.equal_range(imsi);
  for (auto itr = itp.first; itr != itp.second; itr++) {
    GlobalRuleList.get_rule(itr->second, &rule);
    // Get the UE ip address
    rule.mutable_pdi()->set_ue_ipv4(cfg.common_context.ue_ipv4());

    auto src_iface = rule.pdi().src_interface();
    EXPECT_EQ(src_iface, magma::SourceInterfaceType::ACCESS);
    uint32_t upf_teid = session_enforcer->insert_pdr_from_access(
        session_state[1], rule, &session_uc);
    EXPECT_EQ(upf_teid, session_enforcer->get_next_teid());

    EXPECT_EQ(src_iface, magma::SourceInterfaceType::CORE);
    EXPECT_TRUE(session_enforcer->insert_pdr_from_core(
        session_state[1], rule, &session_uc));
  }
  session_enforcer->m5g_pdr_rules_change_and_update_upf(
      session, magma::PdrState::INSTALL);

  session_enforcer->handle_state_update_to_amf(
      *session, magma::lte::M5GSMCause::OPERATION_SUCCESS,
      magma::NotifyUeEvents::UE_SERVICE_REQUEST_ON_PAGING);

  EXPECT_FALSE(it == session_map.end());
  EXPECT_EQ(session_map[IMSI2].size(), 1);
  auto& session_temp = session_map[IMSI2][0];
  EXPECT_EQ(session_temp->get_config().common_context.sid().id(), IMSI2);
  EXPECT_EQ(session_temp->get_config().common_context.ue_ipv4(), IPv4_1);

  session_enforcer->m5g_release_session(
      session_map, imsi, pdu_id, session_update);

  session_enforcer->m5g_start_session_termination(
      session_map, session, pdu_id, &session_uc);

  session_enforcer->m5g_pdr_rules_change_and_update_upf(
      session, magma::PdrState::REMOVE);
}

TEST_F(SessionManagerHandlerTest, test_SetAmfSessionAmbr) {
  magma::SetSMSessionContext request;
  set_sm_session_context_ipv4(&request);
  request.mutable_common_context()->mutable_sid()->set_id(
      "IMSI000000000000002");

  grpc::ServerContext server_context;

  set_session_manager->SetAmfSessionContext(
      &server_context, &request,
      [this](grpc::Status status, SmContextVoid Void) {});

  // Run session creation in the EventBase loop
  evb->loopOnce();

  auto session_map = session_store->read_sessions({IMSI2});
  auto it          = session_map.find(IMSI2);
  EXPECT_FALSE(it == session_map.end());
  EXPECT_EQ(session_map[IMSI2].size(), 1);

  SessionConfig cfg;
  cfg.common_context       = request.common_context();
  cfg.rat_specific_context = request.rat_specific_context();
  cfg.rat_specific_context.mutable_m5gsm_session_context()->set_ssc_mode(
      SSC_MODE_3);

  cfg.rat_specific_context.mutable_m5gsm_session_context()
      ->mutable_default_ambr()
      ->set_br_unit(AggregatedMaximumBitrate::KBPS);
  cfg.rat_specific_context.mutable_m5gsm_session_context()
      ->mutable_default_ambr()
      ->set_max_bandwidth_ul(1024);
  cfg.rat_specific_context.mutable_m5gsm_session_context()
      ->mutable_default_ambr()
      ->set_max_bandwidth_dl(1024);

  SessionUpdate session_update =
      SessionStore::get_default_session_update(session_map);
  uint32_t pdu_id = 5;
  SessionSearchCriteria id1_success_sid(IMSI2, IMSI_AND_PDUID, pdu_id);
  auto session_it = session_store->find_session(session_map, id1_success_sid);
  auto& session   = **session_it;
  auto session_id = session->get_session_id();
  SessionStateUpdateCriteria& session_uc = session_update[IMSI2][session_id];
  SetSmNotificationContext notif;
  std::string imsi;
  session->set_config(cfg, &session_uc);

  magma::SetSMSessionContextAccess expected_response;

  auto* rsp = expected_response.mutable_rat_specific_context()
                  ->mutable_m5g_session_context_rsp();

  rsp->mutable_session_ambr()->set_br_unit(AggregatedMaximumBitrate::KBPS);
  rsp->mutable_session_ambr()->set_max_bandwidth_ul(1024);
  rsp->mutable_session_ambr()->set_max_bandwidth_dl(1024);

  EXPECT_CALL(
      *amf_srv_client,
      handle_response_to_access(CheckSrvResponse(&expected_response)))
      .Times(1);

  session_enforcer->m5g_move_to_active_state(session, notif, &session_uc);
  EXPECT_EQ(request.mutable_common_context()->ue_ipv4(), IPv4_1);
}

TEST_F(SessionManagerHandlerTest, test_create_session_policy_report) {
  magma::SetSMSessionContext request;
  set_sm_session_context_ipv4(&request);

  grpc::ServerContext server_context;

  // create expected request for report_create_session call
  CreateSessionRequest expected_request;
  expected_request.mutable_common_context()->CopyFrom(request.common_context());
  expected_request.mutable_rat_specific_context()->CopyFrom(
      request.rat_specific_context());

  // Create session request towards policydb
  EXPECT_CALL(
      *reporter, report_create_session(CheckSendRequest(expected_request), _))
      .Times(1);

  // create session and expect one call
  set_session_manager->SetAmfSessionContext(
      &server_context, &request,
      [this](grpc::Status status, SmContextVoid Void) {});

  // Run session creation in the EventBase loop
  evb->loopOnce();
  evb->loopOnce();

  // Set the sessionconfig
  SessionConfig cfg;
  cfg.common_context       = request.common_context();
  cfg.rat_specific_context = request.rat_specific_context();
  cfg.rat_specific_context.mutable_m5gsm_session_context()->set_ssc_mode(
      SSC_MODE_3);

  auto session_map = session_store->read_sessions({IMSI1});
  auto it          = session_map.find(IMSI1);

  EXPECT_FALSE(it == session_map.end());
  EXPECT_EQ(session_map[IMSI1].size(), 1);

  auto& session_temp = session_map[IMSI1][0];
  EXPECT_EQ(session_temp->get_config().common_context.sid().id(), IMSI1);

  // Init the session
  session_enforcer->m5g_init_session_credit(
      session_map, IMSI1, SESSION_ID_1, cfg);

  CreateSessionResponse response;
  response.mutable_static_rules()->Add()->mutable_rule_id()->assign("rule1");
  create_credit_update_response(
      IMSI1, SESSION_ID_1, 1, 1025, response.mutable_credits()->Add());
  // Create the session state
  SessionUpdate update = SessionStore::get_default_session_update(session_map);
  SessionSearchCriteria criteria(IMSI1, IMSI_AND_SESSION_ID, SESSION_ID_1);
  auto session_it = session_store->find_session(session_map, criteria);
  auto& session_t = **session_it;
  SessionStateUpdateCriteria* session_uc = &update[IMSI1][SESSION_ID_1];
  session_t->set_config(cfg, session_uc);

  EXPECT_TRUE(session_t->is_5g_session());

  session_enforcer->update_session_with_policy(session_t, response, session_uc);
  session_t->set_upf_teid_endpoint("192.168.200.1", 2147483647, session_uc);
  // Fill the sess_info
  SessionState::SessionInfo sess_info;
  sess_info.local_f_teid  = 2147483647;
  sess_info.subscriber_id = IMSI1;

  // Event report for session creation
  EXPECT_CALL(
      *events_reporter,
      session_created(IMSI1, SESSION_ID_1, testing::_, testing::_))
      .Times(1);

  // Set upf session info
  EXPECT_CALL(
      *pipelined_client, set_upf_session(SessionCheck(sess_info), _, _, _))
      .Times(1);

  // Update the session
  session_enforcer->m5g_update_session_context(
      session_map, IMSI1, session_t, update);

  bool write_success =
      session_store->create_sessions(IMSI1, std::move(session_map[IMSI1]));
  EXPECT_TRUE(write_success);
  auto session_map_2 = session_store->read_sessions(SessionRead{IMSI1});
  EXPECT_EQ(session_map_2[IMSI1].front()->get_request_number(), 1);
}

TEST_F(SessionManagerHandlerTest, test_terminate_session_policy_report) {
  magma::SetSMSessionContext request;
  set_sm_session_context_ipv4(&request);

  grpc::ServerContext server_context;
  // create session and expect one call
  set_session_manager->SetAmfSessionContext(
      &server_context, &request,
      [this](grpc::Status status, SmContextVoid Void) {});

  // Run session creation in the EventBase loop
  evb->loopOnce();
  evb->loopOnce();
  // Set the session config
  SessionConfig cfg;
  cfg.common_context       = request.common_context();
  cfg.rat_specific_context = request.rat_specific_context();
  cfg.rat_specific_context.mutable_m5gsm_session_context()->set_ssc_mode(
      SSC_MODE_3);

  auto session_map = session_store->read_sessions({IMSI1});
  auto it          = session_map.find(IMSI1);
  EXPECT_FALSE(it == session_map.end());
  EXPECT_EQ(session_map[IMSI1].size(), 1);
  auto& session_temp = session_map[IMSI1][0];
  EXPECT_EQ(session_temp->get_config().common_context.sid().id(), IMSI1);
  // create session state
  SessionUpdate update = SessionStore::get_default_session_update(session_map);
  uint32_t pdu_id      = 5;
  SessionSearchCriteria criteria(IMSI1, IMSI_AND_PDUID, pdu_id);
  auto session_it = session_store->find_session(session_map, criteria);
  auto& session   = **session_it;
  auto session_i  = session->get_session_id();
  // SessionTerminateRequest towards policydb
  EXPECT_CALL(
      *reporter,
      report_terminate_session(CheckTerminateRequestCount(IMSI1, 0, 0), _))
      .Times(1);

  // Event report for session terminate
  EXPECT_CALL(*events_reporter, session_terminated(IMSI1, testing::_)).Times(1);

  // session complete terminate
  session_enforcer->m5g_complete_termination(
      session_map, IMSI1, session_i, update);

  EXPECT_EQ(session_map[IMSI1].size(), 0);
}

TEST_F(SessionManagerHandlerTest, test_single_record_5g) {
  magma::SetSMSessionContext request;
  set_sm_session_context_ipv4(&request);
  // Add static rule
  insert_static_rule(1, "", "rule1");

  // Make Session Response from polcydb
  CreateSessionResponse response;
  auto credits = response.mutable_credits();
  create_credit_update_response(IMSI1, SESSION_ID_1, 1, 1024, credits->Add());
  StaticRuleInstall rule1;
  rule1.set_rule_id("rule1");
  auto rules_to_install = response.mutable_static_rules();
  rules_to_install->Add()->CopyFrom(rule1);

  // Fill the RuleRecordTable
  RuleRecordTable table;
  auto record_list = table.mutable_records();
  create_rule_record(
      IMSI1, "192.168.128.11", "rule1", 16, 32, 2147483647, record_list->Add());
  // Set the session config
  SessionConfig cfg;
  cfg.common_context       = request.common_context();
  cfg.rat_specific_context = request.rat_specific_context();
  cfg.rat_specific_context.mutable_m5gsm_session_context()->set_ssc_mode(
      SSC_MODE_3);

  grpc::ServerContext server_context;
  // create session and expect one call
  set_session_manager->SetAmfSessionContext(
      &server_context, &request,
      [this](grpc::Status status, SmContextVoid Void) {});

  // Run session creation in the EventBase loop
  evb->loopOnce();

  auto session_map = session_store->read_sessions({IMSI1});
  // Init the session
  session_enforcer->m5g_init_session_credit(
      session_map, IMSI1, SESSION_ID_1, cfg);
  // Read the session and create session state
  SessionSearchCriteria criteria(IMSI1, IMSI_AND_SESSION_ID, SESSION_ID_1);
  auto session_it = session_store->find_session(session_map, criteria);
  auto& session_t = **session_it;
  auto update     = SessionStore::get_default_session_update(session_map);
  // Update the policy
  session_enforcer->update_session_with_policy(session_t, response, nullptr);
  // Update the session
  session_enforcer->m5g_update_session_context(
      session_map, IMSI1, session_t, update);

  // Report Traffic
  local_enforcer->aggregate_records(session_map, table, update);
  assert_charging_credit(
      session_map, IMSI1, SESSION_ID_1, ALLOWED_TOTAL, {{1, 1024}});
}

TEST_F(SessionManagerHandlerTest, test_update_session_credits_and_rules_5g) {
  magma::SetSMSessionContext request;
  set_sm_session_context_ipv4(&request);
  // Add static rule
  insert_static_rule(1, "", "rule1");
  // Make Session Response from polcydb
  CreateSessionResponse response;
  auto credits = response.mutable_credits();
  create_credit_update_response(IMSI1, SESSION_ID_1, 1, 4096, credits->Add());
  // Set session config
  SessionConfig cfg;
  cfg.common_context       = request.common_context();
  cfg.rat_specific_context = request.rat_specific_context();
  cfg.rat_specific_context.mutable_m5gsm_session_context()->set_ssc_mode(
      SSC_MODE_3);

  grpc::ServerContext server_context;
  // create session and expect one call
  set_session_manager->SetAmfSessionContext(
      &server_context, &request,
      [this](grpc::Status status, SmContextVoid Void) {});

  // Run session creation in the EventBase loop
  evb->loopOnce();
  // Read the session and create the session state
  auto session_map = session_store->read_sessions({IMSI1});
  session_enforcer->m5g_init_session_credit(
      session_map, IMSI1, SESSION_ID_1, cfg);
  SessionSearchCriteria criteria(IMSI1, IMSI_AND_SESSION_ID, SESSION_ID_1);
  auto session_it = session_store->find_session(session_map, criteria);
  auto& session_t = **session_it;
  auto update     = SessionStore::get_default_session_update(session_map);

  // Update the policy
  session_enforcer->update_session_with_policy(session_t, response, nullptr);
  // Update the session
  session_enforcer->m5g_update_session_context(
      session_map, IMSI1, session_t, update);

  assert_charging_credit(
      session_map, IMSI1, SESSION_ID_1, ALLOWED_TOTAL, {{1, 4096}});
  // Add the statis rule
  insert_static_rule(1, "1", "rule1");
  // Fill the RuleRecordTable
  RuleRecordTable table;
  auto record_list = table.mutable_records();
  create_rule_record(IMSI1, "rule1", 1024, 1024, record_list->Add());
  auto uc = get_default_update_criteria();
  session_map[IMSI1][0]->increment_rule_stats("rule1", &uc);

  // Report Traffic
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

  create_credit_update_response(
      IMSI1, SESSION_ID_2, 1, 400, credit_updates_response->Add());
  create_monitor_update_response(
      IMSI1, SESSION_ID_2, "1", MonitoringLevel::PCC_RULE_LEVEL, 500,
      monitor_updates_response->Add());

  session_map = session_store->read_sessions(SessionRead{IMSI1});
  // Update the session credits and rules
  local_enforcer->update_session_credits_and_rules(
      session_map, update_response, update);

  assert_charging_credit(
      session_map, IMSI1, SESSION_ID_1, ALLOWED_TOTAL, {{1, 4120}});
  assert_monitor_credit(
      session_map, IMSI1, SESSION_ID_1, ALLOWED_TOTAL, {{"1", 2048}});
}

int main(int argc, char** argv) {
  ::testing::InitGoogleTest(&argc, argv);
  return RUN_ALL_TESTS();
}
}  // namespace magma

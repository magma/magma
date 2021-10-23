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
#include "Consts.h"
#include "EnumToString.h"
#include "MobilitydClient.h"
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
    pipelined_client = std::make_shared<MockPipelinedClient>();
    amf_srv_client   = std::make_shared<magma::MockAmfServiceClient>();
    magma::mconfig::SessionD mconfig;
    mconfig.set_log_level(magma::orc8r::LogLevel::INFO);
    pdr_map.insert(std::pair<std::string, uint32_t>(IMSI1, 1));
    pdr_map.insert(std::pair<std::string, uint32_t>(IMSI1, 2));

    session_enforcer = std::make_shared<magma::SessionStateEnforcer>(
        rule_store, *session_store, pdr_map, pipelined_client, amf_srv_client,
        reporter.get(), mconfig, session_force_termination_timeout_ms,
        session_max_rtx_count);

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
        session_enforcer, *session_store, reporter.get());

    upf_session_manager = std::make_shared<UpfMsgManageHandler>(
        session_enforcer, mobilityd_client, *session_store);
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
    req->mutable_pdu_address()->set_redirect_address_type(
        magma::RedirectServer::IPV4);
    req->mutable_pdu_address()->set_redirect_server_address("10.20.30.40");
    req->set_priority_access(magma::priorityaccess::High);
    req->set_access_type(magma::AccessType::M_3GPP_ACCESS_3GPP);
    req->set_imei("123456789012345");
    req->set_gpsi("9876543210");
    req->set_pcf_id("1357924680123456");

    reqcmn->mutable_sid()->set_id("IMSI000000000000001");
    reqcmn->set_apn("BLR");
    reqcmn->set_ue_ipv4("192.168.128.11");
    reqcmn->set_rat_type(magma::RATType::TGPP_NR);
    reqcmn->set_sm_session_state(magma::SMSessionFSMState::CREATING_0);

    EXPECT_EQ(req->pdu_session_id(), 0x5);
    EXPECT_EQ(req->request_type(), magma::RequestType::INITIAL_REQUEST);
    EXPECT_EQ(
        req->pdu_address().redirect_address_type(),
        magma::RedirectServer::IPV4);
    EXPECT_EQ(req->pdu_address().redirect_server_address(), "10.20.30.40");
    EXPECT_EQ(req->priority_access(), magma::priorityaccess::High);
    EXPECT_EQ(req->access_type(), magma::AccessType::M_3GPP_ACCESS_3GPP);
    EXPECT_EQ(req->imei(), "123456789012345");
    EXPECT_EQ(req->gpsi(), "9876543210");
    EXPECT_EQ(req->pcf_id(), "1357924680123456");

    EXPECT_EQ(reqcmn->sid().id(), "IMSI000000000000001");
    EXPECT_EQ(reqcmn->sm_session_state(), magma::SMSessionFSMState::CREATING_0);
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
  std::shared_ptr<UpfMsgManageHandler> upf_session_manager;
  std::shared_ptr<MockMobilitydClient> mobilityd_client;
};  // End of class

TEST_F(SessionManagerHandlerTest, test_SetAmfSessionContext) {
  magma::SetSMSessionContext request;
  set_sm_session_context(&request);

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
}

TEST_F(SessionManagerHandlerTest, test_SetSmfNotification) {
  magma::SetSMSessionContext session_ctx_req;
  set_sm_session_context(&session_ctx_req);

  magma::SetSmNotificationContext request;
  set_sm_notif_context(&request, &session_ctx_req);

  grpc::ServerContext server_context;

  set_session_manager->SetSmfNotification(
      &server_context, &request,
      [this](grpc::Status status, SmContextVoid Void) {});

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
  set_sm_session_context(&request);

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
}

TEST_F(SessionManagerHandlerTest, test_UpdateSessionContext) {
  magma::SetSMSessionContext request;
  set_sm_session_context(&request);

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
}

TEST_F(SessionManagerHandlerTest, test_SetPduSessionReleaseContext) {
  magma::SetSMSessionContext request;
  auto* req =
      request.mutable_rat_specific_context()->mutable_m5gsm_session_context();
  auto* reqcmn = request.mutable_common_context();
  req->set_pdu_session_id({0x5});
  req->set_request_type(magma::RequestType::INITIAL_REQUEST);
  req->mutable_pdu_address()->set_redirect_address_type(
      magma::RedirectServer::IPV4);
  req->mutable_pdu_address()->set_redirect_server_address("10.20.30.40");
  req->set_priority_access(magma::priorityaccess::High);
  req->set_imei("123456789012345");
  req->set_gpsi("9876543210");
  req->set_pcf_id("1357924680123456");

  reqcmn->mutable_sid()->set_id("IMSI000000000000001");
  reqcmn->set_sm_session_state(magma::SMSessionFSMState::CREATING_0);

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
  req->mutable_pdu_address()->set_redirect_address_type(
      magma::RedirectServer::IPV4);
  req->mutable_pdu_address()->set_redirect_server_address("10.20.30.40");
  req->set_priority_access(magma::priorityaccess::High);
  req->set_imei("123456789012345");
  req->set_gpsi("9876543210");
  req->set_pcf_id("1357924680123456");

  reqcmn->mutable_sid()->set_id("IMSI000000000000001");
  reqcmn->set_sm_session_state(magma::SMSessionFSMState::CREATING_0);

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
  req->mutable_pdu_address()->set_redirect_address_type(
      magma::RedirectServer::IPV4);
  req->mutable_pdu_address()->set_redirect_server_address("10.20.30.40");
  req->set_priority_access(magma::priorityaccess::High);
  req->set_imei("123456789012345");
  req->set_gpsi("9876543210");
  req->set_pcf_id("1357924680123456");

  reqcmn->mutable_sid()->set_id("IMSI000000000000002");
  reqcmn->set_sm_session_state(magma::SMSessionFSMState::CREATING_0);

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
  req->mutable_pdu_address()->set_redirect_address_type(
      magma::RedirectServer::IPV4);
  req->mutable_pdu_address()->set_redirect_server_address("10.20.30.40");
  req->set_priority_access(magma::priorityaccess::High);
  req->set_imei("123456789012345");
  req->set_gpsi("9876543210");
  req->set_pcf_id("1357924680123456");

  reqcmn->mutable_sid()->set_id("IMSI000000000000002");
  reqcmn->set_sm_session_state(magma::SMSessionFSMState::CREATING_0);

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

TEST_F(SessionManagerHandlerTest, test_PDUStateChangeHandling) {
  magma::SetSMSessionContext request;
  set_sm_session_context(&request);
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
    rule.mutable_pdi()->set_ue_ip_adr(
        cfg.rat_specific_context.m5gsm_session_context()
            .pdu_address()
            .redirect_server_address());

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

  session_enforcer->m5g_release_session(
      session_map, imsi, pdu_id, session_update);

  session_enforcer->m5g_start_session_termination(
      session_map, session, pdu_id, &session_uc);

  session_enforcer->m5g_pdr_rules_change_and_update_upf(
      session, magma::PdrState::REMOVE);
}

TEST_F(SessionManagerHandlerTest, test_SetAmfSessionAmbr) {
  magma::SetSMSessionContext request;
  set_sm_session_context(&request);
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
}

TEST_F(SessionManagerHandlerTest, test_report_rule_stats) {
  insert_static_rule(1, "mk1", "rule1");

  magma::SetSMSessionContext request;
  auto* req =
      request.mutable_rat_specific_context()->mutable_m5gsm_session_context();
  auto* reqcmn = request.mutable_common_context();
  req->set_pdu_session_id({0x5});
  req->set_request_type(magma::RequestType::INITIAL_REQUEST);
  req->mutable_pdu_address()->set_redirect_address_type(
      magma::RedirectServer::IPV4);
  req->mutable_pdu_address()->set_redirect_server_address("10.20.30.40");
  req->set_priority_access(magma::priorityaccess::High);
  req->set_imei("123456789012345");
  req->set_gpsi("9876543210");
  req->set_pcf_id("1357924680123456");

  reqcmn->mutable_sid()->set_id(IMSI1);
  reqcmn->set_sm_session_state(magma::SMSessionFSMState::CREATING_0);

  grpc::ServerContext server_context;

  set_session_manager->SetAmfSessionContext(
      &server_context, &request,
      [this](grpc::Status status, SmContextVoid Void) {});

  // Run session creation in the EventBase loop
  evb->loopOnce();
  evb->loopOnce();

  SessionConfig cfg        = {};
  cfg.common_context       = request.common_context();
  cfg.rat_specific_context = request.rat_specific_context();
  cfg.rat_specific_context.mutable_m5gsm_session_context()->set_ssc_mode(
      SSC_MODE_3);

  auto session_map = session_store->read_sessions({IMSI1});

  session_enforcer->m5g_init_session_credit(
      session_map, IMSI1, SESSION_ID_1, cfg);
  bool write_success =
      session_store->create_sessions(IMSI1, std::move(session_map[IMSI1]));
  EXPECT_TRUE(write_success);
  SessionUpdate update = SessionStore::get_default_session_update(session_map);
  uint32_t pdu_id      = 5;
  SessionSearchCriteria id1_success_sid(IMSI1, IMSI_AND_PDUID, pdu_id);

  // Check the request number
  auto session_map_2 = session_store->read_sessions(SessionRead{IMSI1});
  EXPECT_EQ(session_map_2[IMSI1].front()->get_request_number(), 1);

  // ReportRuleStats
  RuleRecordTable table;
  auto record_list = table.mutable_records();
  create_rule_record(IMSI1, IP1, "rule1", 512, 512, record_list->Add());
  // TODO for will add EXPECT_CALL for usage report once full code is ready.

  upf_session_manager->SendReportRuleStats(
      &server_context, &table,
      [this](grpc::Status status, SmContextVoid Void) {});
  evb->loopOnce();

  session_map_2 = session_store->read_sessions(SessionRead{IMSI1});
  EXPECT_EQ(session_map_2[IMSI1].front()->get_request_number(), 1);
}

TEST_F(SessionManagerHandlerTest, test_collect_updates) {
  insert_static_rule(1, "mk1", "rule1");
  insert_static_rule(2, "mk2", "rule2");
  magma::SetSMSessionContext request;
  auto* req =
      request.mutable_rat_specific_context()->mutable_m5gsm_session_context();
  auto* reqcmn = request.mutable_common_context();
  req->set_pdu_session_id({0x5});
  req->set_request_type(magma::RequestType::INITIAL_REQUEST);
  req->mutable_pdu_address()->set_redirect_address_type(
      magma::RedirectServer::IPV4);
  req->mutable_pdu_address()->set_redirect_server_address("10.20.30.40");
  req->set_priority_access(magma::priorityaccess::High);
  req->set_imei("123456789012345");
  req->set_gpsi("9876543210");
  req->set_pcf_id("1357924680123456");

  reqcmn->mutable_sid()->set_id(IMSI1);
  reqcmn->set_sm_session_state(magma::SMSessionFSMState::CREATING_0);
  reqcmn->set_apn("BLR");
  reqcmn->set_ue_ipv4("192.168.128.11");
  reqcmn->set_rat_type(magma::RATType::TGPP_NR);

  grpc::ServerContext server_context;

  set_session_manager->SetAmfSessionContext(
      &server_context, &request,
      [this](grpc::Status status, SmContextVoid Void) {});

  // Run session creation in the EventBase loop
  evb->loopOnce();
  evb->loopOnce();

  SessionConfig cfg        = {};
  cfg.common_context       = request.common_context();
  cfg.rat_specific_context = request.rat_specific_context();
  cfg.rat_specific_context.mutable_m5gsm_session_context()->set_ssc_mode(
      SSC_MODE_3);
  auto session_map = session_store->read_sessions({IMSI1});

  session_enforcer->m5g_init_session_credit(
      session_map, IMSI1, SESSION_ID_1, cfg);
  bool write_success =
      session_store->create_sessions(IMSI1, std::move(session_map[IMSI1]));
  EXPECT_TRUE(write_success);

  std::vector<std::unique_ptr<ServiceAction>> actions;
  auto update = SessionStore::get_default_session_update(session_map);
  auto empty_update =
      session_enforcer->collect_updates(session_map, actions, update);
  EXPECT_EQ(empty_update.updates_size(), 0);

  RuleRecordTable table;
  auto record_list = table.mutable_records();
  create_rule_record(
      IMSI1, cfg.common_context.ue_ipv4(), "rule1", 1024, 2048,
      record_list->Add());

  create_rule_record(
      IMSI1, cfg.common_context.ue_ipv4(), "rule2", 10, 20, record_list->Add());

  auto update_1 = SessionStore::get_default_session_update(session_map);
  auto uc       = get_default_update_criteria();
  // session_map[IMSI1][0]->increment_rule_stats("rule1", &uc);
  session_enforcer->aggregate_records(session_map, table, update_1);
  EXPECT_EQ(table.records_size(), 2);
}

TEST_F(SessionManagerHandlerTest, test_create_session_report) {
  magma::SetSMSessionContext request;
  set_sm_session_context(&request);

  grpc::ServerContext server_context;

  CreateSessionResponse create_response;
  create_response.mutable_static_rules()->Add()->mutable_rule_id()->assign(
      "rule1");
  create_response.mutable_static_rules()->Add()->mutable_rule_id()->assign(
      "rule2");
  create_credit_update_response(
      IMSI1, SESSION_ID_1, 1, 1536, create_response.mutable_credits()->Add());
  create_credit_update_response(
      IMSI1, SESSION_ID_1, 2, 1024, create_response.mutable_credits()->Add());

  // create expected request for report_create_session call
  CreateSessionRequest expected_request;
  expected_request.mutable_common_context()->CopyFrom(request.common_context());
  expected_request.mutable_rat_specific_context()->CopyFrom(
      request.rat_specific_context());

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

  auto session_map = session_store->read_sessions({IMSI1});
  auto it          = session_map.find(IMSI1);
  EXPECT_FALSE(it == session_map.end());
  EXPECT_EQ(session_map[IMSI1].size(), 1);

  auto& session_temp = session_map[IMSI1][0];
  EXPECT_EQ(session_temp->get_config().common_context.sid().id(), IMSI1);
}

TEST_F(SessionManagerHandlerTest, test_terminate_session_report) {
  CreateSessionResponse response;
  create_credit_update_response(
      IMSI1, SESSION_ID_1, 1, 1024, response.mutable_credits()->Add());
  create_credit_update_response(
      IMSI1, SESSION_ID_1, 2, 2048, response.mutable_credits()->Add());
  magma::SetSMSessionContext request;
  set_sm_session_context(&request);

  grpc::ServerContext server_context;

  set_session_manager->SetAmfSessionContext(
      &server_context, &request,
      [this](grpc::Status status, SmContextVoid Void) {});

  // Run session creation in the EventBase loop
  evb->loopOnce();
  evb->loopOnce();
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
  SessionUpdate update = SessionStore::get_default_session_update(session_map);
  uint32_t pdu_id      = 5;
  SessionSearchCriteria id1_success_sid(IMSI1, IMSI_AND_PDUID, pdu_id);
  auto session_it = session_store->find_session(session_map, id1_success_sid);
  auto& session   = **session_it;
  auto session_i  = session->get_session_id();
  EXPECT_CALL(
      *reporter,
      report_terminate_session(CheckTerminateRequestCount(IMSI1, 0, 0), _))
      .Times(1);
  session_enforcer->m5g_complete_termination(
      session_map, IMSI1, session_i, update);

  EXPECT_EQ(session_map[IMSI1].size(), 0);
}

TEST_F(SessionManagerHandlerTest, test_cleanup_dangling_sessions) {
  magma::SetSMSessionContext request;
  set_sm_session_context(&request);

  grpc::ServerContext server_context;
  SessionState::SessionInfo sess_info;
  SetGroupPDR reqpdr_uplink;
  Action Value   = FORW;
  uint32_t count = DEFAULT_PDR_ID;

  sess_info.local_f_teid  = 10000;
  sess_info.subscriber_id = IMSI1;
  sess_info.ver_no        = 1;
  reqpdr_uplink.set_pdr_id(++count);
  reqpdr_uplink.set_precedence(32);
  reqpdr_uplink.set_pdr_version(1);
  reqpdr_uplink.set_pdr_state(PdrState::INSTALL);
  reqpdr_uplink.mutable_pdi()->set_src_interface(ACCESS);

  reqpdr_uplink.mutable_pdi()->set_net_instance("uplink");
  reqpdr_uplink.set_o_h_remo_desc(0);
  reqpdr_uplink.mutable_set_gr_far()->add_far_action_to_apply(Value);
  reqpdr_uplink.mutable_activate_flow_req()->mutable_request_origin()->set_type(
      RequestOriginType_OriginType_N4);
  reqpdr_uplink.mutable_pdi()->set_ue_ip_adr(IP1);
  reqpdr_uplink.mutable_pdi()->set_local_f_teid(10000);
  reqpdr_uplink.set_pdr_state(PdrState::REMOVE);
  sess_info.pdr_rules.push_back(reqpdr_uplink);

  // PDR 2 details
  SetGroupPDR reqpdr_downlink;
  reqpdr_downlink.set_pdr_id(++count);
  reqpdr_downlink.set_precedence(32);
  reqpdr_downlink.set_pdr_version(1);
  reqpdr_downlink.set_pdr_state(PdrState::INSTALL);
  reqpdr_downlink.mutable_pdi()->set_src_interface(CORE);
  reqpdr_downlink.mutable_set_gr_far()->add_far_action_to_apply(Value);

  // Filling qos params
  reqpdr_downlink.mutable_pdi()->set_net_instance("downlink");
  reqpdr_downlink.mutable_activate_flow_req()
      ->mutable_request_origin()
      ->set_type(RequestOriginType_OriginType_N4);
  reqpdr_downlink.mutable_pdi()->set_ue_ip_adr(IP1);
  reqpdr_downlink.set_pdr_state(PdrState::REMOVE);
  sess_info.pdr_rules.push_back(reqpdr_downlink);

  EXPECT_CALL(
      *pipelined_client,
      set_upf_session(SessionCleanupCheck(sess_info), _, _, _))
      .Times(1);

  session_enforcer->deactivate_flows_for_termination(IMSI1, IP1, "", 10000);
  EXPECT_EQ(sess_info.subscriber_id, IMSI1);
}

int main(int argc, char** argv) {
  ::testing::InitGoogleTest(&argc, argv);
  return RUN_ALL_TESTS();
}
}  // namespace magma

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
#include "SetMessageManagerHandler.h"
#include "ProtobufCreators.h"
#include "RuleStore.h"
#include "SessionState.h"
#include "SessionStore.h"
#include "SessiondMocks.h"
#include "StoredState.h"
#include "magma_logging.h"
#include "PipelinedClient.h"
#include "AmfServiceClient.h"
#include "Consts.h"

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
    rule_store    = std::make_shared<StaticRuleStore>();
    session_store = std::make_shared<SessionStore>(
        rule_store, std::make_shared<MeteringReporter>());
    std::unordered_multimap<std::string, uint32_t> pdr_map;
    pipelined_client = std::make_shared<MockPipelinedClient>();
    amf_srv_client   = std::make_shared<magma::AsyncAmfServiceClient>();
    magma::mconfig::SessionD mconfig;
    mconfig.set_log_level(magma::orc8r::LogLevel::INFO);

    pdr_map.insert(std::pair<std::string, uint32_t>(IMSI1, 1));
    pdr_map.insert(std::pair<std::string, uint32_t>(IMSI1, 2));

    session_enforcer = std::make_shared<magma::SessionStateEnforcer>(
        rule_store, *session_store, pdr_map, pipelined_client, amf_srv_client,
        mconfig, session_force_termination_timeout_ms, session_max_rtx_count);

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
        session_enforcer, *session_store);
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
  std::shared_ptr<AsyncAmfServiceClient> amf_srv_client;
  std::shared_ptr<StaticRuleStore> rule_store;
  SessionIDGenerator id_gen_;
  folly::EventBase* evb;
  SessionMap session_map_;
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
  session->set_config(cfg);
  session_enforcer->add_default_rules(session, IMSI1);
  session_enforcer->m5g_update_session_context(
      session_map, IMSI1, session, update);
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
  set_session_manager->idle_mode_change_sessions_handle(
      request, [](grpc::Status status, SmContextVoid Void) {});

  set_session_manager->pdu_session_inactive(
      request, [](grpc::Status status, SmContextVoid Void) {});

  // Run session creation in the EventBase loop
  evb->loopOnce();
}

int main(int argc, char** argv) {
  ::testing::InitGoogleTest(&argc, argv);
  return RUN_ALL_TESTS();
}

}  // namespace magma

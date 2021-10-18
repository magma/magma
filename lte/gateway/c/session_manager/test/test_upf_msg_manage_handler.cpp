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
#include <sys/socket.h>
#include <netinet/in.h>
#include <arpa/inet.h>
#include <gmock/gmock.h>
#include <gtest/gtest.h>

#include <memory>

#include "UpfMsgManageHandler.h"
#include "SessionStateEnforcer.h"
#include "includes/MagmaService.h"
#include "ProtobufCreators.h"
#include "RuleStore.h"
#include "includes/ServiceRegistrySingleton.h"
#include "SessionState.h"
#include "SessionStore.h"
#include "SessiondMocks.h"
#include "StoredState.h"
#include "magma_logging.h"
#include "PipelinedClient.h"
#include "AmfServiceClient.h"
#include "Consts.h"

#define SESSIOND_SERVICE "sessiond"
#define SESSIOND_VERSION "1.0"
#define session_max_rtx_count 3

using grpc::Status;
using ::testing::Test;

namespace magma {

class UpfMsgMessageHandlerTest : public ::testing::Test {
 public:
  virtual void SetUp() {
    rule_store    = std::make_shared<StaticRuleStore>();
    session_store = std::make_shared<SessionStore>(
        rule_store, std::make_shared<MeteringReporter>());
    std::unordered_multimap<std::string, uint32_t> pdr_map;
    auto pipelined_client = std::make_shared<magma::AsyncPipelinedClient>();
    amf_srv_client        = std::make_shared<magma::AsyncAmfServiceClient>();

    magma::mconfig::SessionD mconfig;
    mconfig.set_log_level(magma::orc8r::LogLevel::INFO);
    auto config =
        magma::ServiceConfigLoader{}.load_service_config(SESSIOND_SERVICE);
    magma::service303::MagmaService server(SESSIOND_SERVICE, SESSIOND_VERSION);

    session_enforcer = std::make_shared<magma::SessionStateEnforcer>(
        rule_store, *session_store, pdr_map, pipelined_client, amf_srv_client,
        mconfig, config["session_force_termination_timeout_ms"].as<int32_t>(),
        session_max_rtx_count);

    evb = new folly::EventBase();
    std::thread([&]() {
      std::cout << "Started event loop thread\n";
      folly::EventBaseManager::get()->setEventBase(evb, 0);
    })
        .detach();
    session_enforcer->attachEventBase(evb);

    set_interface_for_up_mock =
        std::make_shared<MockSetInterfaceForUserPlane>();

    mobilityd_client = std::make_shared<AsyncMobilitydClient>();
    node_request     = std::make_shared<UPFNodeState>();
    sess_config      = std::make_shared<UPFSessionConfigState>();
    page_request     = std::make_shared<UPFPagingInfo>();

    upf_msg_handler = std::make_shared<magma::UpfMsgManageHandler>(
        session_enforcer, mobilityd_client, *session_store);
  }
  virtual void TearDown() { delete evb; }

 public:
  std::shared_ptr<SessionStore> session_store;
  std::shared_ptr<StaticRuleStore> rule_store;
  std::shared_ptr<MockPipelinedClient> pipelined_client;
  std::shared_ptr<SessionStateEnforcer> session_enforcer;
  std::shared_ptr<AsyncAmfServiceClient> amf_srv_client;
  std::shared_ptr<MockSetInterfaceForUserPlane> set_interface_for_up_mock;
  std::shared_ptr<UPFSessionConfigState> sess_config;
  std::shared_ptr<AsyncMobilitydClient> mobilityd_client;
  std::shared_ptr<UPFPagingInfo> page_request;
  std::shared_ptr<UpfMsgManageHandler> upf_msg_handler;
  std::shared_ptr<UPFNodeState> node_request;
  folly::EventBase* evb;
};  // End of class

// Testing the functionality of SendPagingRequest
TEST_F(UpfMsgMessageHandlerTest, test_send_paging_request) {
  // Validate SendPagingRequest
  ON_CALL(*set_interface_for_up_mock, SendPagingRequest(_, _, _))
      .WillByDefault(Return());

  page_request->set_local_f_teid(11);
  page_request->set_ue_ip_addr("192.168.128.11");
  std::function<void(Status, SmContextVoid)> response_callback;

  auto& page_req      = *page_request;
  uint32_t fte_id     = page_req.local_f_teid();
  std::string ip_addr = page_req.ue_ip_addr();
  struct in_addr ue_ip;
  IPAddress req = IPAddress();

  inet_aton(ip_addr.c_str(), &ue_ip);
  req.set_version(IPAddress::IPV4);
  req.set_address(&ue_ip, sizeof(struct in_addr));

  grpc::ServerContext server_context;
  upf_msg_handler->SendPagingRequest(
      &server_context, &page_req,
      [this](grpc::Status status, SmContextVoid Void) {});
  mobilityd_client->get_subscriberid_from_ipv4(
      req, [this, fte_id, response_callback](
               Status status, const SubscriberID& sid) {
        // Validate the status
        EXPECT_TRUE(status.ok());

        const std::string& imsi = sid.id();
        // Validate the imsi length
        EXPECT_TRUE(imsi.length());
        EXPECT_EQ(imsi, IMSI1);

        auto session_map = session_store->read_sessions({imsi});
        SessionSearchCriteria criteria(imsi, IMSI_AND_TEID, fte_id);
        auto session_it = session_store->find_session(session_map, criteria);

        // Validate if we found session from session_store
        EXPECT_TRUE(session_it);

        auto& session = **session_it;

        EXPECT_EQ(fte_id, TEID_1_UL);

        // Validate the session state if it is INACTIVE or not
        EXPECT_TRUE(session->get_state() == INACTIVE);
        // Update the paging request to AMF
        session_enforcer->handle_state_update_to_amf(
            *session, magma::lte::M5GSMCause::OPERATION_SUCCESS,
            UE_PAGING_NOTIFY);

        // After Paging session should be ACTIVE
        EXPECT_TRUE(session->get_state() == ACTIVE);
        EXPECT_TRUE(magma::lte::M5GSMCause::OPERATION_SUCCESS);
      });
}

// Testing the functionality of SetUPFNodeState
TEST_F(UpfMsgMessageHandlerTest, test_upf_node_state) {
  // Validate SetUPFNodeState
  ON_CALL(*set_interface_for_up_mock, SetUPFNodeState(_, _, _))
      .WillByDefault(Return());

  auto& req = *node_request;
  node_request->set_upf_id("1");
  std::function<void(Status, SmContextVoid)> response_callback;

  std::string upf_node_id = "192.168.200.1";
  std::string upf_n3_addr = "192.168.60.112";

  grpc::ServerContext server_context;
  upf_msg_handler->SetUPFNodeState(
      &server_context, &req,
      [this](grpc::Status status, SmContextVoid Void) {});

  switch (req.upf_node_messages_case()) {
    case UPFNodeState::kAssociatonState: {
      UPFAssociationState Assostate = req.associaton_state();
      auto recovery_time            = Assostate.recovery_time_stamp();
      auto feature_set              = Assostate.feature_set();
      std::string local_ipv4_addr =
          Assostate.ip_resource_schema(0).ipv4_address();
      break;
    }
    case UPFNodeState::kNodeReport: {
      break;
    }
    case UPFNodeState::UPF_NODE_MESSAGES_NOT_SET: {
      break;
    }
    default: { break; }
  }
  session_enforcer->set_upf_node(upf_node_id, upf_n3_addr);
  std::string upf_id    = session_enforcer->get_upf_node_id();
  std::string ipv4_addr = session_enforcer->get_upf_n3_addr();
  // Validate the functionality of UPF Node ID
  EXPECT_TRUE(upf_node_id == upf_id);

  // Validate the functionality of UPF N3 interface address
  EXPECT_TRUE(upf_n3_addr == ipv4_addr);
  EXPECT_EQ(req.upf_id(), "1");
}

// Testing the functionality of SetUPFSessionConfig
TEST_F(UpfMsgMessageHandlerTest, test_upf_session_config) {
  // Validate SetUPFSessionConfig
  ON_CALL(*set_interface_for_up_mock, SetUPFSessionConfig(_, _, _))
      .WillByDefault(Return());

  auto& ses_config = *sess_config;
  int32_t count    = 0;
  grpc::ServerContext server_context;
  upf_msg_handler->SetUPFSessionsConfig(
      &server_context, &ses_config,
      [this](grpc::Status status, SmContextVoid Void) {});
  for (auto& upf_session : ses_config.upf_session_state()) {
    // Deleting the IMSI prefix from imsi
    std::string imsi_upf = upf_session.subscriber_id();
    std::string imsi     = imsi_upf.substr(4, imsi_upf.length() - 4);
    uint32_t version     = upf_session.session_version();
    uint32_t teid        = upf_session.local_f_teid();
    auto session_map     = session_store->read_sessions({imsi});
    /* Search with session search criteria of IMSI and session_id and
     * find  respective sesion to operate
     */
    SessionSearchCriteria criteria(imsi, IMSI_AND_TEID, teid);
    auto session_it = session_store->find_session(session_map, criteria);
    // If session not found with SearchCriteria continue with next imsi
    if (!session_it) {
      // Validating session not found with SearchCriteria
      EXPECT_FALSE(session_it);
      continue;
    }
    auto& session    = **session_it;
    auto cur_version = session->get_current_version();

    /* Validating UPF verions of session imsi of teid received version
     * with SMF latest version
     */
    EXPECT_LE(version, cur_version);

    EXPECT_TRUE(
        session_enforcer->is_incremented_rtx_counter_within_max(session));

    RulesToProcess pending_activation, pending_deactivation;
    const CreateSessionResponse& csr = session->get_create_session_response();
    std::vector<StaticRuleInstall> static_rule_installs =
        session_enforcer->to_vec(csr.static_rules());
    std::vector<DynamicRuleInstall> dynamic_rule_installs =
        session_enforcer->to_vec(csr.dynamic_rules());
    session->process_get_5g_rule_installs(
        static_rule_installs, dynamic_rule_installs, &pending_activation,
        &pending_deactivation);

    // Sending the session request to UPF
    session_enforcer->m5g_send_session_request_to_upf(
        session, pending_activation, pending_deactivation);

    // Validating count increment
    auto inc = count++;
    EXPECT_EQ(count, inc + 1);
  }
  // Validating UPF periodic report config missmatch session
  EXPECT_EQ(ses_config.upf_session_state_size(), count);
}

int main(int argc, char** argv) {
  ::testing::InitGoogleTest(&argc, argv);
  return RUN_ALL_TESTS();
}

}  // namespace magma

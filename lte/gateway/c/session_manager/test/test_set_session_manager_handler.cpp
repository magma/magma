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

/* session_manager.grpc.pb.h and SessionStateEnforcer.h
 * included in "SetMessageManagerHandler.h"
 */
#include "AmfServiceClient.h"
#include "magma_logging.h"
#include "includes/MagmaService.h"
#include "PipelinedClient.h"
#include "ProtobufCreators.h"
#include "RuleStore.h"
#include "includes/ServiceRegistrySingleton.h"
#include "SessiondMocks.h"
#include "SessionState.h"
#include "SessionStore.h"
#include "SetMessageManagerHandler.h"
#include "StoredState.h"

using ::testing::Test;

namespace magma {

class SessionManagerHandlerTest : public ::testing::Test {
 public:
  virtual void SetUp() {
    rule_store    = std::make_shared<StaticRuleStore>();
    session_store = std::make_shared<SessionStore>(
        rule_store, std::make_shared<MeteringReporter>());
    auto pipelined_client = std::make_shared<magma::AsyncPipelinedClient>();
    amf_srv_client        = std::make_shared<magma::AsyncAmfServiceClient>();

    magma::mconfig::SessionD mconfig;
    mconfig.set_log_level(magma::orc8r::LogLevel::INFO);

    auto session_enforcer = std::make_shared<magma::SessionStateEnforcer>(
        rule_store, *session_store, pipelined_client, amf_srv_client, mconfig);

    evb = new folly::EventBase();
    std::thread([&]() { folly::EventBaseManager::get()->setEventBase(evb, 0); })
        .detach();

    session_enforcer->attachEventBase(evb);
    session_map_ = SessionMap{};
    // creating landing object and invoking contructor
    set_session_manager = std::make_shared<SetMessageManagerHandler>(
        session_enforcer, *session_store);
  }
  virtual void TearDown() { delete evb; }

 public:
  std::shared_ptr<SessionStore> session_store;
  std::shared_ptr<StaticRuleStore> rule_store;
  std::shared_ptr<SetMessageManagerHandler> set_session_manager;
  std::shared_ptr<MockPipelinedClient> pipelined_client;
  std::shared_ptr<SessionStateEnforcer> session_enforcer;
  std::shared_ptr<AsyncAmfServiceClient> amf_srv_client;
  SessionIDGenerator id_gen_;
  folly::EventBase* evb;
  SessionMap session_map_;
};  // End of class

TEST_F(SessionManagerHandlerTest, test_SetAmfSessionContext) {
  magma::SetSMSessionContext request;
  auto* req =
      request.mutable_rat_specific_context()->mutable_m5gsm_session_context();
  auto* reqcmn = request.mutable_common_context();
  req->set_pdu_session_id({0x5});
  req->set_rquest_type(magma::RequestType::INITIAL_REQUEST);
  req->mutable_pdu_address()->set_redirect_address_type(
      magma::RedirectServer::IPV4);
  req->mutable_pdu_address()->set_redirect_server_address("10.20.30.40");
  req->set_priority_access(magma::priorityaccess::High);
  req->set_access_type(magma::AccessType::M_3GPP_ACCESS_3GPP);
  req->set_imei("123456789012345");
  req->set_gpsi("9876543210");
  req->set_pcf_id("1357924680123456");

  reqcmn->mutable_sid()->set_id("IMSI00000001");
  reqcmn->set_sm_session_state(magma::SMSessionFSMState::CREATING_0);

  grpc::ServerContext server_context;

  set_session_manager->SetAmfSessionContext(
      &server_context, &request,
      [this](grpc::Status status, SmContextVoid Void) {});

  // Run session creation in the EventBase loop
  evb->loopOnce();
}

int main(int argc, char** argv) {
  ::testing::InitGoogleTest(&argc, argv);
  return RUN_ALL_TESTS();
}

}  // namespace magma

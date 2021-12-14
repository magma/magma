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

#include <glog/logging.h>
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

#define SESSIOND_SERVICE "sessiond"
#define SESSIOND_VERSION "1.0"
#define session_max_rtx_count 3

using ::testing::Test;

namespace magma {

class SetUPFNodeState : public ::testing::Test {
 public:
  virtual void SetUp() {
    rule_store = std::make_shared<StaticRuleStore>();
    reporter = std::make_shared<MockSessionReporter>();
    session_store = std::make_shared<SessionStore>(
        rule_store, std::make_shared<MeteringReporter>());
    std::unordered_multimap<std::string, uint32_t> pdr_map;
    auto pipelined_client = std::make_shared<magma::AsyncPipelinedClient>();
    amf_srv_client = std::make_shared<magma::AsyncAmfServiceClient>();
    events_reporter = std::make_shared<MockEventsReporter>();

    magma::mconfig::SessionD mconfig;
    mconfig.set_log_level(magma::orc8r::LogLevel::INFO);
    auto config =
        magma::ServiceConfigLoader{}.load_service_config(SESSIOND_SERVICE);
    magma::service303::MagmaService server(SESSIOND_SERVICE, SESSIOND_VERSION);

    session_enforcer = std::make_shared<magma::SessionStateEnforcer>(
        rule_store, *session_store, pdr_map, pipelined_client, amf_srv_client,
        reporter.get(), events_reporter, mconfig,
        config["session_force_termination_timeout_ms"].as<int32_t>(),
        session_max_rtx_count);

    set_interface_for_up_mock =
        std::make_shared<MockSetInterfaceForUserPlane>();
    sess_config = std::make_shared<UPFSessionConfigState>();
  }
  virtual void TearDown() {}

 public:
  std::shared_ptr<SessionStore> session_store;
  std::shared_ptr<StaticRuleStore> rule_store;
  std::shared_ptr<MockPipelinedClient> pipelined_client;
  std::shared_ptr<SessionStateEnforcer> session_enforcer;
  std::shared_ptr<AsyncAmfServiceClient> amf_srv_client;
  std::shared_ptr<MockSetInterfaceForUserPlane> set_interface_for_up_mock;
  std::shared_ptr<UPFSessionConfigState> sess_config;
  std::shared_ptr<MockSessionReporter> reporter;
  std::shared_ptr<MockEventsReporter> events_reporter;
};  // End of class

// Testing the functionality of SetUPFNodeState
TEST_F(SetUPFNodeState, test_upf_node_state) {
  // Validate SetUPFNodeState
  ON_CALL(*set_interface_for_up_mock, SetUPFNodeState(_, _, _))
      .WillByDefault(Return(Status::OK));

  std::string upf_node_id = "192.168.200.1";
  std::string upf_n3_addr = "192.168.60.112";

  session_enforcer->set_upf_node(upf_node_id, upf_n3_addr);
  std::string upf_id = session_enforcer->get_upf_node_id();
  std::string ipv4_addr = session_enforcer->get_upf_n3_addr();

  // Validate the functionality of UPF Node ID
  EXPECT_TRUE(upf_node_id == upf_id);

  // Validate the functionality of UPF N3 interface address
  EXPECT_TRUE(upf_n3_addr == ipv4_addr);
}

// Testing the functionality of SetUPFSessionConfig
TEST_F(SetUPFNodeState, test_upf_session_config) {
  // Validate SetUPFSessionConfig
  ON_CALL(*set_interface_for_up_mock, SetUPFSessionConfig(_, _, _))
      .WillByDefault(Return(Status::OK));

  auto& ses_config = *sess_config;
  int32_t count = 0;
  for (auto& upf_session : ses_config.upf_session_state()) {
    // Deleting the IMSI prefix from imsi
    std::string imsi_upf = upf_session.subscriber_id();
    std::string imsi = imsi_upf.substr(4, imsi_upf.length() - 4);
    uint32_t version = upf_session.session_version();
    uint32_t teid = upf_session.local_f_teid();
    auto session_map = session_store->read_sessions({imsi});
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
    auto& session = **session_it;
    auto cur_version = session->get_current_version();

    /* Validating UPF verions of session imsi of teid received version
     * with SMF latest version
     */
    EXPECT_LE(version, cur_version);
    RulesToProcess pending_activation, pending_deactivation;
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

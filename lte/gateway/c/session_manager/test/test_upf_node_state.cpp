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

using ::testing::Test;

namespace magma {

class SetUPFNodeState : public ::testing::Test {
 public:
  virtual void SetUp() {
    rule_store    = std::make_shared<StaticRuleStore>();
    session_store = std::make_shared<SessionStore>(
        rule_store, std::make_shared<MeteringReporter>());
    auto pipelined_client = std::make_shared<magma::AsyncPipelinedClient>();
    amf_srv_client        = std::make_shared<magma::AsyncAmfServiceClient>();

    magma::mconfig::SessionD mconfig;
    mconfig.set_log_level(magma::orc8r::LogLevel::INFO);
    auto config =
        magma::ServiceConfigLoader{}.load_service_config(SESSIOND_SERVICE);
    magma::service303::MagmaService server(SESSIOND_SERVICE, SESSIOND_VERSION);

    session_enforcer = std::make_shared<magma::SessionStateEnforcer>(
        rule_store, *session_store, pipelined_client, amf_srv_client, mconfig,
        config["session_force_termination_timeout_ms"].as<int32_t>());

    set_interface_for_up_mock =
        std::make_shared<MockSetInterfaceForUserPlane>();
  }
  virtual void TearDown() {}

 public:
  std::shared_ptr<SessionStore> session_store;
  std::shared_ptr<StaticRuleStore> rule_store;
  std::shared_ptr<MockPipelinedClient> pipelined_client;
  std::shared_ptr<SessionStateEnforcer> session_enforcer;
  std::shared_ptr<AsyncAmfServiceClient> amf_srv_client;
  std::shared_ptr<MockSetInterfaceForUserPlane> set_interface_for_up_mock;
};  // End of class

// Testing the functionality of SetUPFNodeState
TEST_F(SetUPFNodeState, test_upf_node_state) {
  // Validate SetUPFNodeState
  ON_CALL(*set_interface_for_up_mock, SetUPFNodeState(_, _, _))
      .WillByDefault(Return(Status::OK));

  std::string upf_node_id = "192.168.200.1";
  std::string upf_n3_addr = "192.168.60.112";

  session_enforcer->set_upf_node(upf_node_id, upf_n3_addr);
  std::string upf_id    = session_enforcer->get_upf_node_id();
  std::string ipv4_addr = session_enforcer->get_upf_n3_addr();

  // Validate the functionality of UPF Node ID
  EXPECT_TRUE(upf_node_id == upf_id);

  // Validate the functionality of UPF N3 interface address
  EXPECT_TRUE(upf_n3_addr == ipv4_addr);
}

int main(int argc, char** argv) {
  ::testing::InitGoogleTest(&argc, argv);
  return RUN_ALL_TESTS();
}

}  // namespace magma

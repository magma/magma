/**
 * Copyright 2022 The Magma Authors.
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

#include <gtest/gtest.h>
#include <iostream>
#include <nlohmann/json.hpp>

#include "lte/gateway/c/session_manager/EnumToString.hpp"
#include "lte/gateway/c/session_manager/OperationalStatesHandler.hpp"
#include "lte/gateway/c/session_manager/RuleStore.hpp"
#include "lte/gateway/c/session_manager/SessionStore.hpp"
#include "lte/gateway/c/session_manager/Types.hpp"
#include "lte/gateway/c/session_manager/test/Consts.hpp"
#include "lte/gateway/c/session_manager/test/ProtobufCreators.hpp"

namespace magma {

class OperationalStatesHandlerTest : public ::testing::Test {
 protected:
  virtual void SetUp() {
    rule_store = std::make_shared<StaticRuleStore>();
    session_store = std::make_shared<SessionStore>(
        rule_store, std::make_shared<MeteringReporter>());
  }
  virtual void TearDown() {}
  SessionVector create_session_vec(std::string device_id, std::string ipv4,
                                   uint64_t session_start_time,
                                   std::string session_id) {
    Teids teids;
    SessionConfig cfg;
    cfg.common_context =
        build_common_context(device_id, ipv4, "", teids, APN1, "", TGPP_WLAN);
    auto session = std::make_unique<SessionState>(session_id, cfg, *rule_store,
                                                  session_start_time);
    session->set_fsm_state(SESSION_ACTIVE, nullptr);
    auto session_vec = SessionVector{};
    session_vec.push_back(std::move(session));
    return session_vec;
  }
  std::shared_ptr<SessionStore> session_store;
  std::shared_ptr<StaticRuleStore> rule_store;
};

TEST_F(OperationalStatesHandlerTest, test_get_operational_states) {
  std::list<std::map<std::string, std::string>> states =
      get_operational_states(session_store.get());
  EXPECT_EQ(states.size(), 1);
  std::map<std::string, std::string> state = states.front();
  auto it = state.find(TYPE);
  EXPECT_FALSE(it == state.end());
  EXPECT_EQ(it->second, GATEWAY_SUBSCRIBER_STATE_TYPE);
  it = state.find(VALUE);
  EXPECT_FALSE(it == state.end());
  nlohmann::json json_value = nlohmann::json::parse(it->second);
  nlohmann::json subscribers = json_value[SUBSCRIBERS];
  EXPECT_TRUE(subscribers.empty());
  auto session_vec =
      create_session_vec(IMSI1, IP1, SESSION_START_TIME_1, SESSION_ID_1);
  session_store->create_sessions(IMSI1, std::move(session_vec));
  session_vec =
      create_session_vec(IMSI2, IP2, SESSION_START_TIME_2, SESSION_ID_2);
  session_store->create_sessions(IMSI2, std::move(session_vec));
  states = get_operational_states(session_store.get());
  EXPECT_EQ(states.size(), 1);

  std::map<std::string, std::string> subscriber_state = states.front();
  std::map<std::string, std::string>::iterator subscriber_state_it =
      subscriber_state.find(TYPE);
  EXPECT_FALSE(subscriber_state_it == subscriber_state.end());
  EXPECT_EQ(subscriber_state_it->second, GATEWAY_SUBSCRIBER_STATE_TYPE);

  subscriber_state_it = subscriber_state.find(VALUE);
  EXPECT_FALSE(subscriber_state_it == subscriber_state.end());
  json_value = nlohmann::json::parse(subscriber_state_it->second);
  nlohmann::json gateway_subscribers = json_value[SUBSCRIBERS];
  nlohmann::json content = gateway_subscribers[IMSI1][APN1];
  EXPECT_EQ(content.size(), 1);
  content = content[0];
  EXPECT_EQ(content["lifecycle_state"],
            session_fsm_state_to_str(SESSION_ACTIVE));
  EXPECT_EQ(content["session_start_time"], SESSION_START_TIME_1);
  EXPECT_EQ(content["apn"], APN1);
  EXPECT_EQ(content["ipv4"], IP1);
  EXPECT_EQ(content["session_id"], SESSION_ID_1);

  content = gateway_subscribers[IMSI2][APN1];
  EXPECT_EQ(content.size(), 1);
  content = content[0];
  EXPECT_EQ(content["lifecycle_state"],
            session_fsm_state_to_str(SESSION_ACTIVE));
  EXPECT_EQ(content["session_start_time"], SESSION_START_TIME_2);
  EXPECT_EQ(content["apn"], APN1);
  EXPECT_EQ(content["ipv4"], IP2);
  EXPECT_EQ(content["session_id"], SESSION_ID_2);
}
}  // namespace magma

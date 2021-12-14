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

#include "Consts.h"
#include "magma_logging.h"
#include "MemoryStoreClient.h"
#include "ProtobufCreators.h"
#include "RuleStore.h"
#include "SessionID.h"
#include "SessionState.h"

using ::testing::Test;

namespace magma {

class StoreClientTest : public ::testing::Test {
 protected:
  SessionIDGenerator id_gen_;
};

/**
 * End to end test of the MemoryStoreClient.
 * 1) Create MemoryStoreClient
 * 2) Read in sessions for subscribers IMSI1 and IMSI2
 * 3) Create bare-bones session for IMSI1 and IMSI2
 * 4) Write and commit session state for IMSI1 and IMSI2
 * 5) Read for subscribers IMSI1 and IMSI2
 * 6) Verify that state was written for IMSI1/IMSI2 and has been retrieved.
 */
TEST_F(StoreClientTest, test_read_and_write) {
  std::string hardware_addr_bytes = {0x0f, 0x10, 0x2e, 0x12, 0x3a, 0x55};
  SessionConfig cfg1, cfg2, cfg3;
  Teids teids;
  teids.set_agw_teid(1);
  teids.set_enb_teid(2);
  cfg1.common_context =
      build_common_context(IMSI1, IP1, IPv6_1, teids, APN1, MSISDN, TGPP_WLAN);
  cfg2.common_context =
      build_common_context(IMSI2, IP1, IPv6_1, teids, APN1, MSISDN, TGPP_WLAN);
  cfg3.common_context =
      build_common_context(IMSI3, IP1, IPv6_1, teids, APN1, MSISDN, TGPP_WLAN);
  const auto& wlan = build_wlan_context("0f:10:2e:12:3a:55", RADIUS_SESSION_ID);
  cfg1.rat_specific_context.mutable_wlan_context()->CopyFrom(wlan);
  cfg2.rat_specific_context.mutable_wlan_context()->CopyFrom(wlan);
  cfg3.rat_specific_context.mutable_wlan_context()->CopyFrom(wlan);

  CreateSessionResponse response1, response2, response3;
  create_credit_update_response(IMSI1, SESSION_ID_1, 1, 1000,
                                response1.mutable_credits()->Add());
  create_credit_update_response(IMSI2, SESSION_ID_2, 2, 2000,
                                response2.mutable_credits()->Add());
  create_credit_update_response(IMSI3, SESSION_ID_3, 3, 3000,
                                response3.mutable_credits()->Add());

  // Emulate CreateSession, which needs to create a new session for a
  // subscriber
  std::set<std::string> requested_ids{IMSI1, IMSI2};
  auto rule_store = std::make_shared<StaticRuleStore>();
  auto tgpp_context = TgppContext{};
  auto pdp_start_time = 12345;

  auto store_client = MemoryStoreClient(rule_store);
  auto session_map = store_client.read_sessions(requested_ids);

  auto uc = get_default_update_criteria();

  auto session1 = std::make_unique<SessionState>(SESSION_ID_1, cfg1,
                                                 *rule_store, pdp_start_time);
  session1->set_tgpp_context(tgpp_context, nullptr);
  session1->set_fsm_state(SESSION_ACTIVE, nullptr);
  session1->set_create_session_response(response1, nullptr);

  auto session2 = std::make_unique<SessionState>(SESSION_ID_2, cfg2,
                                                 *rule_store, pdp_start_time);
  session2->set_tgpp_context(tgpp_context, nullptr);
  session2->set_fsm_state(SESSION_ACTIVE, nullptr);
  session2->set_create_session_response(response2, nullptr);

  auto session3 = std::make_unique<SessionState>(SESSION_ID_3, cfg3,
                                                 *rule_store, pdp_start_time);
  session3->set_tgpp_context(tgpp_context, nullptr);
  session3->set_fsm_state(SESSION_ACTIVE, nullptr);
  session3->set_create_session_response(response3, nullptr);

  EXPECT_EQ(session1->get_session_id(), SESSION_ID_1);
  EXPECT_EQ(session2->get_session_id(), SESSION_ID_2);

  RuleLifetime lifetime{};
  session1->activate_static_rule("rule1", lifetime, &uc);
  EXPECT_EQ(session1->is_static_rule_installed("rule1"), true);

  EXPECT_EQ(session_map.size(), 2);
  EXPECT_EQ(session_map[IMSI1].size(), 0);
  session_map[IMSI1].push_back(std::move(session1));
  EXPECT_EQ(session_map[IMSI1].size(), 1);

  // Since the grant was not given with R/W permission for subscriber IMSI2,
  // The session for IMSI2 should not be saved into the store
  session_map[IMSI2] = SessionVector();
  session_map[IMSI2].push_back(std::move(session2));
  EXPECT_EQ(session_map.size(), 2);
  EXPECT_EQ(session_map[IMSI2].size(), 1);

  // And now commit back to storage (memory actually, but later persistent)
  store_client.write_sessions(std::move(session_map));

  // Try to do a read to make sure that things are the same
  auto session_map_2 = store_client.read_sessions(requested_ids);
  EXPECT_EQ(session_map_2.size(), 2);
  EXPECT_EQ(session_map_2[IMSI1].size(), 1);
  EXPECT_EQ(session_map_2[IMSI1].front()->get_session_id(), SESSION_ID_1);
  EXPECT_EQ(session_map_2[IMSI1].front()->is_static_rule_installed("rule1"),
            true);
  EXPECT_EQ(session_map_2[IMSI1].front()->get_config(), cfg1);
  EXPECT_EQ(
      session_map_2[IMSI1].front()->get_create_session_response().DebugString(),
      response1.DebugString());

  // Now create a third session
  std::set<std::string> requested_imsi3{IMSI3};
  auto session_map_3 = store_client.read_sessions(requested_imsi3);
  EXPECT_EQ(session_map_3.size(), 1);
  session_map_3[IMSI3].push_back(std::move(session3));
  EXPECT_EQ(session_map_3[IMSI3].size(), 1);
  store_client.write_sessions(std::move(session_map_3));

  // Get all sessions
  auto all_sessions = store_client.read_all_sessions();
  EXPECT_EQ(all_sessions.size(), 3);
  EXPECT_EQ(all_sessions[IMSI1].size(), 1);
  EXPECT_EQ(all_sessions[IMSI1].front()->get_session_id(), SESSION_ID_1);
  EXPECT_EQ(all_sessions[IMSI3].size(), 1);
  EXPECT_EQ(all_sessions[IMSI3].front()->get_session_id(), SESSION_ID_3);
  EXPECT_EQ(
      all_sessions[IMSI3].front()->get_create_session_response().DebugString(),
      response3.DebugString());
}

TEST_F(StoreClientTest, test_lambdas) {
  auto sm = std::make_unique<int>(1);

  std::function<void(std::unique_ptr<int>&)> callback2 =
      [](std::unique_ptr<int>& inp) { EXPECT_EQ(*inp, 2); };

  std::function<void()> callback =
      [=, shared = std::make_shared<decltype(sm)>(std::move(sm))]() mutable {
        EXPECT_EQ(**shared, 1);
        **shared = 2;
        callback2(*shared);
      };
  callback();
}

int main(int argc, char** argv) {
  ::testing::InitGoogleTest(&argc, argv);
  return RUN_ALL_TESTS();
}

}  // namespace magma

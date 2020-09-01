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

#include <glog/logging.h>
#include <gtest/gtest.h>

#include "Consts.h"
#include "MemoryStoreClient.h"
#include "ProtobufCreators.h"
#include "RuleStore.h"
#include "SessionID.h"
#include "SessionState.h"
#include "magma_logging.h"

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
  auto sid                    = id_gen_.gen_session_id(IMSI1);
  auto sid2                   = id_gen_.gen_session_id(IMSI2);
  auto sid3                   = id_gen_.gen_session_id(IMSI3);
  SessionConfig cfg;
  cfg.common_context = build_common_context("", IP2, "APN", MSISDN, TGPP_WLAN);
  const auto& wlan = build_wlan_context("0f:10:2e:12:3a:55", RADIUS_SESSION_ID);
  cfg.rat_specific_context.mutable_wlan_context()->CopyFrom(wlan);
  auto rule_store   = std::make_shared<StaticRuleStore>();
  auto tgpp_context = TgppContext{};
  auto pdp_start_time = 12345;

  auto store_client = new MemoryStoreClient(rule_store);

  // Emulate CreateSession, which needs to create a new session for a subscriber
  auto session_map = store_client->read_sessions({IMSI1, IMSI2});

  auto uc = get_default_update_criteria();
  auto session =
      std::make_unique<SessionState>(IMSI1, sid, cfg, *rule_store, tgpp_context, pdp_start_time);
  auto session2 = std::make_unique<SessionState>(
      IMSI2, sid2, cfg, *rule_store, tgpp_context, pdp_start_time);
  auto session3 = std::make_unique<SessionState>(
      IMSI3, sid3, cfg, *rule_store, tgpp_context, pdp_start_time);
  EXPECT_EQ(session->get_session_id(), sid);
  EXPECT_EQ(session2->get_session_id(), sid2);

  RuleLifetime lifetime{};
  session->activate_static_rule("rule1", lifetime, uc);
  EXPECT_EQ(session->is_static_rule_installed("rule1"), true);

  EXPECT_EQ(session_map.size(), 2);
  EXPECT_EQ(session_map[IMSI1].size(), 0);
  session_map[IMSI1].push_back(std::move(session));
  EXPECT_EQ(session_map[IMSI1].size(), 1);

  // Since the grant was not given with R/W permission for subscriber IMSI2,
  // The session for IMSI2 should not be saved into the store
  session_map[IMSI2] = std::vector<std::unique_ptr<SessionState>>();
  session_map[IMSI2].push_back(std::move(session2));
  EXPECT_EQ(session_map.size(), 2);
  EXPECT_EQ(session_map[IMSI2].size(), 1);

  // And now commit back to storage (memory actually, but later persistent)
  store_client->write_sessions(std::move(session_map));

  // Try to do a read to make sure that things are the same
  auto session_map_2 = store_client->read_sessions({IMSI1, IMSI2});
  EXPECT_EQ(session_map_2.size(), 2);
  EXPECT_EQ(session_map_2[IMSI1].size(), 1);
  EXPECT_EQ(session_map_2[IMSI1].front()->get_session_id(), sid);
  EXPECT_EQ(
      session_map_2[IMSI1].front()->is_static_rule_installed("rule1"), true);

  // Now create a third session
  auto session_map_3 = store_client->read_sessions({IMSI3});
  EXPECT_EQ(session_map_3.size(), 1);
  session_map_3[IMSI3].push_back(std::move(session3));
  EXPECT_EQ(session_map_3[IMSI3].size(), 1);
  store_client->write_sessions(std::move(session_map_3));

  // Get all sessions
  auto all_sessions = store_client->read_all_sessions();
  EXPECT_EQ(all_sessions.size(), 3);
  EXPECT_EQ(all_sessions[IMSI1].size(), 1);
  EXPECT_EQ(all_sessions[IMSI1].front()->get_session_id(), sid);
  EXPECT_EQ(all_sessions[IMSI3].size(), 1);
  EXPECT_EQ(all_sessions[IMSI3].front()->get_session_id(), sid3);
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

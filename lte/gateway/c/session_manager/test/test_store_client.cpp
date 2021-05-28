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
  std::string imsi                = "IMSI1";
  std::string imsi2               = "IMSI2";
  std::string imsi3               = "IMSI3";
  std::string msisdn              = "5100001234";
  std::string radius_session_id =
      "AA-AA-AA-AA-AA-AA:TESTAP__"
      "0F-10-2E-12-3A-55";
  auto sid  = id_gen_.gen_session_id(imsi);
  auto sid2 = id_gen_.gen_session_id(imsi2);
  auto sid3 = id_gen_.gen_session_id(imsi3);
  Teids teids;
  teids.set_agw_teid(1);
  teids.set_enb_teid(2);
  SessionConfig cfg;
  cfg.common_context = build_common_context(
      "", "128.0.0.1", "2001:0db8:0a0b:12f0:0000:0000:0000:0001", teids, "APN",
      msisdn, TGPP_WLAN);
  const auto& wlan = build_wlan_context("0f:10:2e:12:3a:55", radius_session_id);
  cfg.rat_specific_context.mutable_wlan_context()->CopyFrom(wlan);
  auto rule_store     = std::make_shared<StaticRuleStore>();
  auto tgpp_context   = TgppContext{};
  auto pdp_start_time = 12345;

  auto store_client = MemoryStoreClient(rule_store);

  // Emulate CreateSession, which needs to create a new session for a subscriber
  std::set<std::string> requested_ids{imsi, imsi2};
  auto session_map = store_client.read_sessions(requested_ids);

  auto uc = get_default_update_criteria();

  CreateSessionResponse response1;
  auto credits = response1.mutable_credits();
  create_credit_update_response(imsi, sid, 1, 1000, credits->Add());
  auto session = std::make_unique<SessionState>(
      imsi, sid, cfg, *rule_store, tgpp_context, pdp_start_time, response1);

  CreateSessionResponse response2;
  credits = response2.mutable_credits();
  create_credit_update_response(imsi2, sid2, 2, 2000, credits->Add());
  auto session2 = std::make_unique<SessionState>(
      imsi2, sid2, cfg, *rule_store, tgpp_context, pdp_start_time, response2);

  CreateSessionResponse response3;
  credits = response3.mutable_credits();
  create_credit_update_response(imsi3, sid3, 3, 3000, credits->Add());
  auto session3 = std::make_unique<SessionState>(
      imsi3, sid3, cfg, *rule_store, tgpp_context, pdp_start_time, response3);

  EXPECT_EQ(session->get_session_id(), sid);
  EXPECT_EQ(session2->get_session_id(), sid2);

  RuleLifetime lifetime{};
  session->activate_static_rule("rule1", lifetime, &uc);
  EXPECT_EQ(session->is_static_rule_installed("rule1"), true);

  EXPECT_EQ(session_map.size(), 2);
  EXPECT_EQ(session_map[imsi].size(), 0);
  session_map[imsi].push_back(std::move(session));
  EXPECT_EQ(session_map[imsi].size(), 1);

  // Since the grant was not given with R/W permission for subscriber IMSI2,
  // The session for IMSI2 should not be saved into the store
  session_map[imsi2] = SessionVector();
  session_map[imsi2].push_back(std::move(session2));
  EXPECT_EQ(session_map.size(), 2);
  EXPECT_EQ(session_map[imsi2].size(), 1);

  // And now commit back to storage (memory actually, but later persistent)
  store_client.write_sessions(std::move(session_map));

  // Try to do a read to make sure that things are the same
  auto session_map_2 = store_client.read_sessions(requested_ids);
  EXPECT_EQ(session_map_2.size(), 2);
  EXPECT_EQ(session_map_2[imsi].size(), 1);
  EXPECT_EQ(session_map_2[imsi].front()->get_session_id(), sid);
  EXPECT_EQ(
      session_map_2[imsi].front()->is_static_rule_installed("rule1"), true);
  EXPECT_EQ(session_map_2[imsi].front()->get_config(), cfg);
  EXPECT_EQ(
      session_map_2[imsi].front()->get_create_session_response().DebugString(),
      response1.DebugString());

  // Now create a third session
  std::set<std::string> requested_imsi3{imsi3};
  auto session_map_3 = store_client.read_sessions(requested_imsi3);
  EXPECT_EQ(session_map_3.size(), 1);
  session_map_3[imsi3].push_back(std::move(session3));
  EXPECT_EQ(session_map_3[imsi3].size(), 1);
  store_client.write_sessions(std::move(session_map_3));

  // Get all sessions
  auto all_sessions = store_client.read_all_sessions();
  EXPECT_EQ(all_sessions.size(), 3);
  EXPECT_EQ(all_sessions[imsi].size(), 1);
  EXPECT_EQ(all_sessions[imsi].front()->get_session_id(), sid);
  EXPECT_EQ(all_sessions[imsi3].size(), 1);
  EXPECT_EQ(all_sessions[imsi3].front()->get_session_id(), sid3);
  EXPECT_EQ(
      all_sessions[imsi3].front()->get_create_session_response().DebugString(),
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

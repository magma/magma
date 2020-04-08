/**
 * Copyright (c) 2016-present, Facebook, Inc.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include <memory>

#include <glog/logging.h>
#include <gtest/gtest.h>

#include "RuleStore.h"
#include "SessionID.h"
#include "SessionState.h"
#include "MemoryStoreClient.h"
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
TEST_F(StoreClientTest, test_read_and_write)
{
  std::string hardware_addr_bytes = {0x0f, 0x10, 0x2e, 0x12, 0x3a, 0x55};
  std::string imsi = "IMSI1";
  std::string imsi2 = "IMSI2";
  std::string msisdn = "5100001234";
  std::string radius_session_id =
    "AA-AA-AA-AA-AA-AA:TESTAP__"
    "0F-10-2E-12-3A-55";
  auto sid = id_gen_.gen_session_id(imsi);
  auto sid2 = id_gen_.gen_session_id(imsi2);
  std::string core_session_id = "asdf";
  SessionState::Config cfg = {.ue_ipv4 = "",
    .spgw_ipv4 = "",
    .msisdn = msisdn,
    .apn = "",
    .imei = "",
    .plmn_id = "",
    .imsi_plmn_id = "",
    .user_location = "",
    .rat_type = RATType::TGPP_WLAN,
    .mac_addr = "0f:10:2e:12:3a:55",
    .hardware_addr = hardware_addr_bytes,
    .radius_session_id = radius_session_id};
  auto rule_store = std::make_shared<StaticRuleStore>();
  auto tgpp_context = TgppContext{};

  auto store_client = new MemoryStoreClient(rule_store);

  // Emulate CreateSession, which needs to create a new session for a subscriber
  std::set<std::string> requested_ids{imsi, imsi2};
  auto session_map = store_client->read_sessions(requested_ids);

  auto uc = get_default_update_criteria();
  auto session = std::make_unique<SessionState>(imsi, sid, core_session_id, cfg, *rule_store, tgpp_context);
  auto session2 = std::make_unique<SessionState>(imsi2, sid2, core_session_id, cfg, *rule_store, tgpp_context);
  EXPECT_EQ(session->get_session_id(), sid);
  EXPECT_EQ(session2->get_session_id(), sid2);

  session->activate_static_rule("rule1", uc);
  EXPECT_EQ(session->is_static_rule_installed("rule1"), true);

  EXPECT_EQ(session_map.size(), 2);
  EXPECT_EQ(session_map[imsi].size(), 0);
  session_map[imsi].push_back(std::move(session));
  EXPECT_EQ(session_map[imsi].size(), 1);

  // Since the grant was not given with R/W permission for subscriber IMSI2,
  // The session for IMSI2 should not be saved into the store
  session_map[imsi2] = std::vector<std::unique_ptr<SessionState>>();
  session_map[imsi2].push_back(std::move(session2));
  EXPECT_EQ(session_map.size(), 2);
  EXPECT_EQ(session_map[imsi2].size(), 1);

  // And now commit back to storage (memory actually, but later persistent)
  store_client->write_sessions(std::move(session_map));

  // Try to do a read to make sure that things are the same
  auto session_map_2 = store_client->read_sessions(requested_ids);
  EXPECT_EQ(session_map_2.size(), 2);
  EXPECT_EQ(session_map_2[imsi].size(), 1);
  EXPECT_EQ(session_map_2[imsi].front()->get_session_id(), sid);
  EXPECT_EQ(session_map_2[imsi].front()->is_static_rule_installed("rule1"), true);
}

TEST_F(StoreClientTest, test_lambdas)
{
  auto sm = std::make_unique<int>(1);

  std::function<void(std::unique_ptr<int>&)> callback2 = [](std::unique_ptr<int>& inp) {
    EXPECT_EQ(*inp, 2);
  };

  std::function<void()> callback = [=, shared = std::make_shared<decltype(sm)>(std::move(sm))]() mutable {
    EXPECT_EQ(**shared, 1);
    **shared = 2;
    callback2(*shared);
  };
  callback();
}

int main(int argc, char **argv)
{
  ::testing::InitGoogleTest(&argc, argv);
  return RUN_ALL_TESTS();
}

} // namespace magma

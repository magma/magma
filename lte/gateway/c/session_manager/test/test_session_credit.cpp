/**
 * Copyright (c) 2016-present, Facebook, Inc.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */
#include <chrono>
#include <thread>

#include "SessionCredit.h"
#include "ProtobufCreators.h"
#include <gtest/gtest.h>

using ::testing::Test;

namespace magma {
const FinalActionInfo default_final_action_info = {
    .final_action = ChargingCredit_FinalAction_TERMINATE};

TEST(test_marshal_unmarshal, test_session_credit) {
  SessionCredit credit;
  SessionCreditUpdateCriteria uc{};

  // Set some fields here to non default values. Credit is used.
  credit.add_used_credit((uint64_t) 39u, (uint64_t) 40u, uc);

  // Sanity check of credit usage. Test result after marshal/unmarshal should
  // match.
  EXPECT_EQ(credit.get_credit(USED_TX), (uint64_t) 39u);
  EXPECT_EQ(credit.get_credit(USED_RX), (uint64_t) 40u);

  // Check that the update criteria also includes the changes
  EXPECT_EQ(uc.bucket_deltas[USED_TX], (uint64_t) 39u);
  EXPECT_EQ(uc.bucket_deltas[USED_RX], (uint64_t) 40u);

  // Check that after marshaling/unmarshaling that the fields are still the
  // same.
  auto marshaled = credit.marshal();
  auto credit_2  = SessionCredit::unmarshal(marshaled);

  EXPECT_EQ(credit_2.get_credit(USED_TX), (uint64_t) 39u);
  EXPECT_EQ(credit_2.get_credit(USED_RX), (uint64_t) 40u);
}

TEST(test_track_credit, test_session_credit) {
  SessionCredit credit;
  SessionCreditUpdateCriteria uc{};
  GrantedUnits gsu;
  uint64_t grant = 1024;
  create_granted_units(&grant, NULL, NULL, &gsu);

  credit.receive_credit(gsu, &uc);

  EXPECT_EQ(1024, credit.get_credit(ALLOWED_TOTAL));
  EXPECT_EQ(0, credit.get_credit(USED_TX));

  // Check that the update criteria also includes the changes
  EXPECT_EQ(uc.bucket_deltas[ALLOWED_TOTAL], 1024);
  EXPECT_EQ(uc.bucket_deltas[USED_TX], 0);
}

TEST(test_add_received_credit, test_session_credit) {
  SessionCredit credit;
  SessionCreditUpdateCriteria uc{};
  GrantedUnits gsu;
  uint64_t grant = 1024;
  create_granted_units(&grant, NULL, NULL, &gsu);

  credit.receive_credit(gsu, &uc);
  credit.add_used_credit(10, 20, uc);
  EXPECT_EQ(credit.get_credit(USED_TX), 10);
  EXPECT_EQ(credit.get_credit(USED_RX), 20);
  EXPECT_EQ(uc.bucket_deltas[USED_TX], 10);
  EXPECT_EQ(uc.bucket_deltas[USED_RX], 20);
  credit.add_used_credit(30, 40, uc);
  EXPECT_EQ(credit.get_credit(USED_TX), 40);
  EXPECT_EQ(credit.get_credit(USED_RX), 60);
  EXPECT_EQ(uc.bucket_deltas[USED_TX], 40);
  EXPECT_EQ(uc.bucket_deltas[USED_RX], 60);
}

TEST(test_collect_updates, test_session_credit) {
  SessionCredit credit;
  SessionCreditUpdateCriteria uc{};
  GrantedUnits gsu;
  uint64_t grant = 1024;
  create_granted_units(&grant, NULL, NULL, &gsu);

  credit.receive_credit(gsu, &uc);
  credit.add_used_credit(500, 524, uc);
  EXPECT_TRUE(credit.is_quota_exhausted(0.8));
  auto update = credit.get_usage_for_reporting(uc);
  EXPECT_EQ(update.bytes_tx, 500);
  EXPECT_EQ(update.bytes_rx, 524);

  EXPECT_TRUE(credit.is_reporting());
  EXPECT_TRUE(uc.reporting);
  EXPECT_EQ(credit.get_credit(REPORTING_TX), 500);
  EXPECT_EQ(credit.get_credit(REPORTING_RX), 524);
  // Only track how much has been reported. The currently reporting amount
  // should be held in memory only.
  EXPECT_EQ(uc.bucket_deltas[REPORTING_TX], 0);
  EXPECT_EQ(uc.bucket_deltas[REPORTING_RX], 0);
}

// Default usage reporting threshold is 0.8, so session manager will report
// when quota is not completely used up.
TEST(test_collect_updates_when_nearly_exhausted, test_session_credit) {
  SessionCredit credit;
  SessionCreditUpdateCriteria uc{};
  GrantedUnits gsu;
  uint64_t grant = 1000;
  create_granted_units(&grant, NULL, NULL, &gsu);

  credit.receive_credit(gsu, &uc);
  credit.add_used_credit(300, 500, uc);
  EXPECT_TRUE(credit.is_quota_exhausted(0.8));
  auto update = credit.get_usage_for_reporting(uc);
  EXPECT_EQ(update.bytes_tx, 300);
  EXPECT_EQ(update.bytes_rx, 500);

  EXPECT_TRUE(credit.is_reporting());
  EXPECT_TRUE(uc.reporting);
  EXPECT_EQ(credit.get_credit(REPORTING_TX), 300);
  EXPECT_EQ(credit.get_credit(REPORTING_RX), 500);
  // Only track how much has been reported. The currently reporting amount
  // should be held in memory only.
  EXPECT_EQ(uc.bucket_deltas[REPORTING_TX], 0);
  EXPECT_EQ(uc.bucket_deltas[REPORTING_RX], 0);
}

TEST(test_collect_updates_none_available, test_session_credit) {
  SessionCredit credit;
  SessionCreditUpdateCriteria uc{};
  GrantedUnits gsu;
  uint64_t grant = 1000;
  create_granted_units(&grant, NULL, NULL, &gsu);

  credit.receive_credit(gsu, &uc);
  credit.add_used_credit(400, 399, uc);
  EXPECT_FALSE(credit.is_quota_exhausted(0.8));
}

// The maximum of reported usage is capped by what is granted even when an user
// overused.
TEST(test_collect_updates_when_overusing, test_session_credit) {
  SessionCredit credit;
  SessionCreditUpdateCriteria uc{};
  GrantedUnits gsu;
  uint64_t grant = 1000;
  create_granted_units(&grant, NULL, NULL, &gsu);

  credit.receive_credit(gsu, &uc);
  credit.add_used_credit(510, 500, uc);
  EXPECT_TRUE(credit.is_quota_exhausted(0.8));
  auto update = credit.get_usage_for_reporting(uc);
  EXPECT_EQ(update.bytes_tx, 510);
  EXPECT_EQ(update.bytes_rx, 490);

  EXPECT_TRUE(credit.is_reporting());
  EXPECT_TRUE(uc.reporting);
  EXPECT_EQ(credit.get_credit(REPORTING_TX), 510);
  EXPECT_EQ(credit.get_credit(REPORTING_RX), 490);
  // Only track how much has been reported. The currently reporting amount
  // should be held in memory only.
  EXPECT_EQ(uc.bucket_deltas[REPORTING_TX], 0);
  EXPECT_EQ(uc.bucket_deltas[REPORTING_RX], 0);
}

TEST(test_add_rx_tx_credit, test_session_credit) {
  SessionCredit credit;
  SessionCreditUpdateCriteria uc{};
  GrantedUnits gsu;
  uint64_t grant = 1000;
  create_granted_units(NULL, &grant, &grant, &gsu);

  // receive tx = 1000, rx = 1000
  credit.receive_credit(gsu, &uc);
  EXPECT_EQ(uc.grant_tracking_type, TX_AND_RX);
  // use tx = 1000
  credit.add_used_credit(1000, 0, uc);
  EXPECT_TRUE(credit.is_quota_exhausted(0.8));
  auto update = credit.get_usage_for_reporting(uc);
  EXPECT_EQ(update.bytes_tx, 1000);
  EXPECT_EQ(update.bytes_rx, 0);

  // receive another tx = 1000, rx = 1000, cumulative: tx = 2000, rx = 2000
  credit.receive_credit(gsu, &uc);
  EXPECT_EQ(uc.grant_tracking_type, TX_AND_RX);
  credit.add_used_credit(50, 2000, uc);
  EXPECT_TRUE(credit.is_quota_exhausted(0.8));
  auto update2 = credit.get_usage_for_reporting(uc);
  EXPECT_EQ(update2.bytes_tx, 50);
  EXPECT_EQ(update2.bytes_rx, 1000);

  // receive rx, tx, but no usage
  gsu.Clear();
  create_granted_units(NULL, &grant, &grant, &gsu);
  credit.receive_credit(gsu, &uc);
  EXPECT_FALSE(credit.is_quota_exhausted(0.8));
}

TEST(test_is_quota_exhausted_total_only, test_session_credit) {
  SessionCredit credit;
  SessionCreditUpdateCriteria uc{};

  GrantedUnits gsu;
  uint64_t total_grant = 1000;
  create_granted_units(&total_grant, NULL, NULL, &gsu);
  credit.receive_credit(gsu, &uc);
  EXPECT_EQ(uc.grant_tracking_type, TOTAL_ONLY);

  credit.add_used_credit(500, 0, uc);
  EXPECT_FALSE(credit.is_quota_exhausted(0.8));
  credit.mark_failure(0, &uc);

  credit.add_used_credit(500, 0, uc);
  EXPECT_TRUE(credit.is_quota_exhausted(1));
}

TEST(test_is_quota_exhausted_rx_only, test_session_credit) {
  SessionCredit credit;
  SessionCreditUpdateCriteria uc{};

  GrantedUnits gsu;
  uint64_t grant = 1000;
  create_granted_units(NULL, NULL, &grant, &gsu);
  credit.receive_credit(gsu, &uc);
  EXPECT_EQ(uc.grant_tracking_type, RX_ONLY);

  credit.add_used_credit(500, 500, uc);
  EXPECT_FALSE(credit.is_quota_exhausted(0.8));
  credit.mark_failure(0, &uc);

  credit.add_used_credit(500, 500, uc);
  EXPECT_TRUE(credit.is_quota_exhausted(1));
}

TEST(test_is_quota_exhausted_tx_only, test_session_credit) {
  SessionCredit credit;
  SessionCreditUpdateCriteria uc{};

  GrantedUnits gsu;
  uint64_t grant = 1000;
  create_granted_units(NULL, &grant, NULL, &gsu);
  credit.receive_credit(gsu, &uc);
  EXPECT_EQ(uc.grant_tracking_type, TX_ONLY);

  credit.add_used_credit(500, 500, uc);
  EXPECT_FALSE(credit.is_quota_exhausted(0.8));
  credit.mark_failure(0, &uc);

  credit.add_used_credit(500, 500, uc);
  EXPECT_TRUE(credit.is_quota_exhausted(1));
}

// Assert that receiving an invalid GSU does not change the way we track
// credits. If we request an additional quota and receive an empty GSU, the
// quota should still be exhausted.
TEST(test_is_quota_exhausted_after_empty_grant, test_session_credit) {
  SessionCredit credit;
  SessionCreditUpdateCriteria uc{};

  GrantedUnits gsu;
  uint64_t total_grant = 1000;
  create_granted_units(&total_grant, NULL, NULL, &gsu);
  credit.receive_credit(gsu, &uc);
  EXPECT_EQ(uc.grant_tracking_type, TOTAL_ONLY);

  // Add enough quota to hit is_quota_exhausted(0.8)
  credit.add_used_credit(900, 0, uc);
  EXPECT_TRUE(credit.is_quota_exhausted(0.8));
  credit.mark_failure(0, &uc);

  // Receive empty GSU, quota should still be exhausted
  gsu.Clear();
  credit.receive_credit(gsu, &uc);
  // assert uc grant_tracking type has not changed
  EXPECT_EQ(uc.grant_tracking_type, TOTAL_ONLY);
  EXPECT_TRUE(credit.is_quota_exhausted(0.8));
}

int main(int argc, char** argv) {
  ::testing::InitGoogleTest(&argc, argv);
  return RUN_ALL_TESTS();
}

}  // namespace magma

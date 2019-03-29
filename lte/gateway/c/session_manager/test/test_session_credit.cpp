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

#include <gtest/gtest.h>
#include "SessionCredit.h"

using ::testing::Test;

#define HIGH_CREDIT 1000000

namespace magma {

TEST(test_track_credit, test_session_credit)
{
  SessionCredit credit;
  credit.receive_credit(1024, HIGH_CREDIT, HIGH_CREDIT, 3600, false);

  EXPECT_EQ(1024, credit.get_credit(ALLOWED_TOTAL));
  EXPECT_EQ(0, credit.get_credit(USED_TX));
}

TEST(test_add_received_credit, test_session_credit)
{
  SessionCredit credit;
  credit.receive_credit(1024, HIGH_CREDIT, HIGH_CREDIT, 3600, false);
  credit.add_used_credit(10, 20);
  EXPECT_EQ(credit.get_credit(USED_TX), 10);
  EXPECT_EQ(credit.get_credit(USED_RX), 20);
  credit.add_used_credit(30, 40);
  EXPECT_EQ(credit.get_credit(USED_TX), 40);
  EXPECT_EQ(credit.get_credit(USED_RX), 60);
}

TEST(test_collect_updates, test_session_credit)
{
  SessionCredit credit;
  credit.receive_credit(1024, HIGH_CREDIT, HIGH_CREDIT, 3600, false);
  credit.add_used_credit(1024, 20);
  EXPECT_EQ(credit.get_update_type(), CREDIT_QUOTA_EXHAUSTED);
  auto update = credit.get_usage_for_reporting(false);
  EXPECT_EQ(update.bytes_tx, 1024);
  EXPECT_EQ(update.bytes_rx, 20);

  EXPECT_TRUE(credit.is_reporting());
  EXPECT_EQ(credit.get_credit(REPORTING_TX), 1024);
  EXPECT_EQ(credit.get_credit(REPORTING_RX), 20);
}

TEST(test_collect_updates_timer_expiries, test_credit_manager)
{
  SessionCredit credit;
  credit.receive_credit(1024, HIGH_CREDIT, HIGH_CREDIT, 1, false);
  credit.add_used_credit(20, 30);

  std::this_thread::sleep_for(std::chrono::milliseconds(1001));
  EXPECT_EQ(credit.get_update_type(), CREDIT_VALIDITY_TIMER_EXPIRED);
  auto update = credit.get_usage_for_reporting(false);
  EXPECT_EQ(update.bytes_tx, 20);
  EXPECT_EQ(update.bytes_rx, 30);
}

TEST(test_collect_updates_none_available, test_session_credit)
{
  SessionCredit credit;
  credit.receive_credit(1024, HIGH_CREDIT, HIGH_CREDIT, 3600, false);
  credit.add_used_credit(30, 20);
  EXPECT_EQ(credit.get_update_type(), CREDIT_NO_UPDATE);
}

TEST(test_get_action, test_session_credit)
{
  SessionCredit credit;
  credit.receive_credit(1024, HIGH_CREDIT, HIGH_CREDIT, 0, false);
  credit.add_used_credit(1025, 0);
  auto cont_action = credit.get_action();
  EXPECT_EQ(cont_action, CONTINUE_SERVICE);
  credit.receive_credit(1024, HIGH_CREDIT, HIGH_CREDIT, 0, true);
  credit.add_used_credit(2048, 0);
  credit.add_used_credit(30, 20);
  auto term_action = credit.get_action();
  EXPECT_EQ(term_action, TERMINATE_SERVICE);

  // Termination action only returned once
  auto repeated_action = credit.get_action();
  EXPECT_EQ(repeated_action, CONTINUE_SERVICE);
}

TEST(test_failures, test_session_credit)
{
  SessionCredit credit;
  credit.receive_credit(1024, HIGH_CREDIT, HIGH_CREDIT, 0, true);
  credit.add_used_credit(1025, 0);
  EXPECT_EQ(credit.get_action(), TERMINATE_SERVICE);

  SessionCredit credit2;
  credit2.receive_credit(1024, HIGH_CREDIT, HIGH_CREDIT, 0, false);
  credit2.add_used_credit(10, 0);
  EXPECT_EQ(credit2.get_action(), CONTINUE_SERVICE);
  credit2.mark_failure();
  EXPECT_EQ(credit2.get_action(), TERMINATE_SERVICE);
}

TEST(test_add_rx_tx_credit, test_session_credit)
{
  SessionCredit credit;
  credit.receive_credit(HIGH_CREDIT, 1000, 1000, 3600, false);
  credit.add_used_credit(1500, 0);
  EXPECT_EQ(credit.get_update_type(), CREDIT_QUOTA_EXHAUSTED);
  auto update = credit.get_usage_for_reporting(false);
  EXPECT_EQ(update.bytes_tx, 1500);
  EXPECT_EQ(update.bytes_rx, 0);

  // receive tx
  credit.receive_credit(0, 1000, 0, 3600, false);
  credit.add_used_credit(0, 1500);
  EXPECT_EQ(credit.get_update_type(), CREDIT_QUOTA_EXHAUSTED);
  auto update2 = credit.get_usage_for_reporting(false);
  EXPECT_EQ(update2.bytes_tx, 0);
  EXPECT_EQ(update2.bytes_rx, 1500);

  // receive rx
  credit.receive_credit(0, 0, 1000, 3600, false);
  credit.add_used_credit(1000, 1000);
  EXPECT_EQ(credit.get_update_type(), CREDIT_QUOTA_EXHAUSTED);
  auto update3 = credit.get_usage_for_reporting(false);
  EXPECT_EQ(update3.bytes_tx, 1000);
  EXPECT_EQ(update3.bytes_rx, 1000);

  // receive rx
  credit.receive_credit(0, 1000, 1000, 3600, false);
  EXPECT_EQ(credit.get_update_type(), CREDIT_NO_UPDATE);
}

int main(int argc, char **argv)
{
  ::testing::InitGoogleTest(&argc, argv);
  return RUN_ALL_TESTS();
}

} // namespace magma

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
const SessionCredit::FinalActionInfo default_final_action_info = {
  .final_action = ChargingCredit_FinalAction_TERMINATE};

class SessionCreditParameterizedTest :
  public ::testing::TestWithParam<CreditType> {};

TEST_P(SessionCreditParameterizedTest, test_marshal_unmarshal) {
  CreditType credit_type = GetParam();
  SessionCredit credit(credit_type);

  // Set some fields here to non default values. Credit is used.
  credit.add_used_credit((uint64_t) 39u, (uint64_t) 40u);

  // Sanity check of credit usage. Test result after marshal/unmarshal should
  // match.
  EXPECT_EQ(credit.get_credit(USED_TX), (uint64_t) 39u);
  EXPECT_EQ(credit.get_credit(USED_RX), (uint64_t) 40u);

  // Check that after marshaling/unmarshaling that the fields are still the
  // same.
  auto marshaled = credit.marshal();
  auto credit_2 = SessionCredit::unmarshal(marshaled, credit_type);

  EXPECT_EQ(credit_2->get_credit(USED_TX), (uint64_t) 39u);
  EXPECT_EQ(credit_2->get_credit(USED_RX), (uint64_t) 40u);
}

TEST_P(SessionCreditParameterizedTest, test_track_credit) {
  CreditType credit_type = GetParam();
  SessionCredit credit(credit_type);

  credit.receive_credit(1024, HIGH_CREDIT, HIGH_CREDIT, 3600, false,
    default_final_action_info);

  EXPECT_EQ(1024, credit.get_credit(ALLOWED_TOTAL));
  EXPECT_EQ(0, credit.get_credit(USED_TX));
}

TEST_P(SessionCreditParameterizedTest, test_add_received_credit)
{
  CreditType credit_type = GetParam();
  SessionCredit credit(credit_type);

  credit.receive_credit(1024, HIGH_CREDIT, HIGH_CREDIT, 3600, false,
    default_final_action_info);
  credit.add_used_credit(10, 20);
  EXPECT_EQ(credit.get_credit(USED_TX), 10);
  EXPECT_EQ(credit.get_credit(USED_RX), 20);
  credit.add_used_credit(30, 40);
  EXPECT_EQ(credit.get_credit(USED_TX), 40);
  EXPECT_EQ(credit.get_credit(USED_RX), 60);
}

TEST_P(SessionCreditParameterizedTest, test_collect_updates)
{
  CreditType credit_type = GetParam();
  SessionCredit credit(credit_type);

  credit.receive_credit(1024, HIGH_CREDIT, HIGH_CREDIT, 3600, false,
    default_final_action_info);
  credit.add_used_credit(500, 524);
  EXPECT_EQ(credit.get_update_type(), CREDIT_QUOTA_EXHAUSTED);
  auto update = credit.get_usage_for_reporting(false);
  EXPECT_EQ(update.bytes_tx, 500);
  EXPECT_EQ(update.bytes_rx, 524);

  EXPECT_TRUE(credit.is_reporting());
  EXPECT_EQ(credit.get_credit(REPORTING_TX), 500);
  EXPECT_EQ(credit.get_credit(REPORTING_RX), 524);
}

/*
 * Default usage reporting threshold is 0.8, so session manager will report
 * when quota is not completely used up.
 */
TEST_P(SessionCreditParameterizedTest,
       test_collect_updates_when_nearly_exhausted)
{
  CreditType credit_type = GetParam();
  SessionCredit credit(credit_type);

  credit.receive_credit(1000, HIGH_CREDIT, HIGH_CREDIT, 3600, false,
    default_final_action_info);
  credit.add_used_credit(300, 500);
  EXPECT_EQ(credit.get_update_type(), CREDIT_QUOTA_EXHAUSTED);
  auto update = credit.get_usage_for_reporting(false);
  EXPECT_EQ(update.bytes_tx, 300);
  EXPECT_EQ(update.bytes_rx, 500);

  EXPECT_TRUE(credit.is_reporting());
  EXPECT_EQ(credit.get_credit(REPORTING_TX), 300);
  EXPECT_EQ(credit.get_credit(REPORTING_RX), 500);
}

TEST_P(SessionCreditParameterizedTest, test_collect_updates_timer_expiries)
{
  CreditType credit_type = GetParam();
  SessionCredit credit(credit_type);

  credit.receive_credit(1024, HIGH_CREDIT, HIGH_CREDIT, 1, false,
    default_final_action_info);
  credit.add_used_credit(20, 30);

  std::this_thread::sleep_for(std::chrono::milliseconds(1001));
  EXPECT_EQ(credit.get_update_type(), CREDIT_VALIDITY_TIMER_EXPIRED);
  auto update = credit.get_usage_for_reporting(false);
  EXPECT_EQ(update.bytes_tx, 20);
  EXPECT_EQ(update.bytes_rx, 30);
}

TEST_P(SessionCreditParameterizedTest, test_collect_updates_none_available)
{
  CreditType credit_type = GetParam();
  SessionCredit credit(credit_type);

  credit.receive_credit(1000, HIGH_CREDIT, HIGH_CREDIT, 3600, false,
    default_final_action_info);
  credit.add_used_credit(400, 399);
  EXPECT_EQ(credit.get_update_type(), CREDIT_NO_UPDATE);
}

/*
 * The maximum of reported usage is capped by what is granted even when an user
 * overused.
 */
TEST_P(SessionCreditParameterizedTest, test_collect_updates_when_overusing)
{
  CreditType credit_type = GetParam();
  SessionCredit credit(credit_type);

  credit.receive_credit(1000, HIGH_CREDIT, HIGH_CREDIT, 3600, false,
    default_final_action_info);
  credit.add_used_credit(510, 500);
  EXPECT_EQ(credit.get_update_type(), CREDIT_QUOTA_EXHAUSTED);
  auto update = credit.get_usage_for_reporting(false);
  EXPECT_EQ(update.bytes_tx, 510);
  EXPECT_EQ(update.bytes_rx, 490);

  EXPECT_TRUE(credit.is_reporting());
  EXPECT_EQ(credit.get_credit(REPORTING_TX), 510);
  EXPECT_EQ(credit.get_credit(REPORTING_RX), 490);
}

TEST_P(SessionCreditParameterizedTest, test_add_rx_tx_credit)
{
  CreditType credit_type = GetParam();
  SessionCredit credit(credit_type);

  // receive tx
  credit.receive_credit(1000, 1000, 0, 3600, false,
    default_final_action_info);
  credit.add_used_credit(1000, 0);
  EXPECT_EQ(credit.get_update_type(), CREDIT_QUOTA_EXHAUSTED);
  auto update = credit.get_usage_for_reporting(false);
  EXPECT_EQ(update.bytes_tx, 1000);
  EXPECT_EQ(update.bytes_rx, 0);

  // receive rx
  credit.receive_credit(1000, 0, 1000, 3600, false,
    default_final_action_info);
  credit.add_used_credit(0, 1000);
  EXPECT_EQ(credit.get_update_type(), CREDIT_QUOTA_EXHAUSTED);
  auto update2 = credit.get_usage_for_reporting(false);
  EXPECT_EQ(update2.bytes_tx, 0);
  EXPECT_EQ(update2.bytes_rx, 1000);

  // receive rx, tx, but no usage
  credit.receive_credit(2000, 1000, 1000, 3600, false,
    default_final_action_info);
  EXPECT_EQ(credit.get_update_type(), CREDIT_NO_UPDATE);
}

INSTANTIATE_TEST_CASE_P(
  SessionCreditTests,
  SessionCreditParameterizedTest,
  ::testing::Values(MONITORING, CHARGING));

TEST(test_get_action_for_charging, test_session_credit)
{
  // Test Charging Credit
  SessionCredit charging_credit(CreditType::CHARGING);
  charging_credit.receive_credit(1024, HIGH_CREDIT, HIGH_CREDIT, 0, false,
    default_final_action_info);
  charging_credit.add_used_credit(1024, 0);
  auto cont_action = charging_credit.get_action();
  EXPECT_EQ(cont_action, CONTINUE_SERVICE);
  charging_credit.receive_credit(1024, HIGH_CREDIT, HIGH_CREDIT, 0, true,
    default_final_action_info);
  charging_credit.add_used_credit(2048, 0);
  charging_credit.add_used_credit(30, 20);
  auto term_action = charging_credit.get_action();
  EXPECT_EQ(term_action, TERMINATE_SERVICE);

  // Termination action only returned once
  auto repeated_action = charging_credit.get_action();
  EXPECT_EQ(repeated_action, CONTINUE_SERVICE);
}

TEST(test_get_action_for_monitoring, test_session_credit)
{
  // Monitoring Credit should never return TERMINATE_SERVICE
  SessionCredit monitoring_credit(CreditType::MONITORING);
  monitoring_credit.receive_credit(1024, HIGH_CREDIT, HIGH_CREDIT, 0, false,
    default_final_action_info);
  monitoring_credit.add_used_credit(1024, 0);
  auto cont_action = monitoring_credit.get_action();
  EXPECT_EQ(cont_action, CONTINUE_SERVICE);
  monitoring_credit.receive_credit(1024, HIGH_CREDIT, HIGH_CREDIT, 0, true,
    default_final_action_info);
  monitoring_credit.add_used_credit(2048, 0);
  monitoring_credit.add_used_credit(30, 20);
  auto term_action = monitoring_credit.get_action();
  EXPECT_EQ(term_action, CONTINUE_SERVICE);
}

TEST(test_last_grant_exhausted_for_charging, test_session_credit)
{
  SessionCredit charging_credit(CreditType::CHARGING);
  charging_credit.receive_credit(1024, HIGH_CREDIT, HIGH_CREDIT, 0, true,
    default_final_action_info);
  charging_credit.add_used_credit(1024, 0);
  EXPECT_EQ(charging_credit.get_action(), TERMINATE_SERVICE);
}

TEST(test_last_grant_exhausted_for_monitoring, test_session_credit)
{
  SessionCredit monitoring_credit(CreditType::MONITORING);
  monitoring_credit.receive_credit(1024, HIGH_CREDIT, HIGH_CREDIT, 0, true,
    default_final_action_info);
  monitoring_credit.add_used_credit(1024, 0);
  EXPECT_EQ(monitoring_credit.get_action(), CONTINUE_SERVICE);
}

TEST(test_final_unit_action_restrict_access, test_session_credit)
{
  SessionCredit::FinalActionInfo final_action_info;
  final_action_info.final_action = ChargingCredit_FinalAction_RESTRICT_ACCESS;

  SessionCredit charging_credit(CreditType::CHARGING);
  charging_credit.receive_credit(1024, HIGH_CREDIT, HIGH_CREDIT, 0, true,
    final_action_info);
  charging_credit.add_used_credit(1024, 0);
  EXPECT_EQ(charging_credit.get_action(), RESTRICT_ACCESS);
}

TEST(test_final_unit_action_redirect, test_session_credit)
{
  SessionCredit::FinalActionInfo final_action_info;
  final_action_info.final_action = ChargingCredit_FinalAction_REDIRECT;

  SessionCredit charging_credit(CreditType::CHARGING);
  charging_credit.receive_credit(1024, HIGH_CREDIT, HIGH_CREDIT, 0, true,
    final_action_info);
  charging_credit.add_used_credit(1024, 0);
  EXPECT_EQ(charging_credit.get_action(), REDIRECT);
}

TEST(test_tolerance_quota_exhausted, test_session_credit)
{
  SessionCredit credit(CreditType::CHARGING);
  credit.receive_credit(1024, HIGH_CREDIT, HIGH_CREDIT, 0, false,
    default_final_action_info);
  // continue the service when there was still available tolerance quota
  credit.add_used_credit(1024, 0);
  EXPECT_EQ(credit.get_action(), CONTINUE_SERVICE);
  // terminate the service when granted quota and tolerance quota are exhausted
  credit.add_used_credit(1024, 0);
  EXPECT_EQ(credit.get_action(), TERMINATE_SERVICE);
}

TEST(test_failures, test_session_credit)
{
  SessionCredit credit(CreditType::CHARGING);
  credit.receive_credit(1024, HIGH_CREDIT, HIGH_CREDIT, 0, false,
    default_final_action_info);
  credit.add_used_credit(1024, 0);
  EXPECT_EQ(credit.get_action(), CONTINUE_SERVICE);
  credit.mark_failure();
  EXPECT_EQ(credit.get_action(), CONTINUE_SERVICE);
  // extra tolerance quota are exhausted
  credit.add_used_credit(1024, 0);
  credit.mark_failure();
  EXPECT_EQ(credit.get_action(), TERMINATE_SERVICE);
}

int main(int argc, char **argv)
{
  ::testing::InitGoogleTest(&argc, argv);
  return RUN_ALL_TESTS();
}

} // namespace magma

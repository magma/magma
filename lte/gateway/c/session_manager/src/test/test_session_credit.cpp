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
#include <gtest/gtest.h>

#include <chrono>
#include <thread>

#include "ProtobufCreators.h"
#include "SessionCredit.h"

using ::testing::Test;

namespace magma {
const FinalActionInfo default_final_action_info = {
    .final_action = ChargingCredit_FinalAction_TERMINATE};

TEST(test_marshal_unmarshal, test_session_credit) {
  SessionCredit credit;
  SessionCreditUpdateCriteria uc{};

  // Set some fields here to non default values. Credit is used.
  credit.add_used_credit((uint64_t) 39u, (uint64_t) 40u, &uc);

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
  SessionCredit credit_2(marshaled);

  EXPECT_EQ(credit_2.get_credit(USED_TX), (uint64_t) 39u);
  EXPECT_EQ(credit_2.get_credit(USED_RX), (uint64_t) 40u);
}

TEST(test_track_credit, test_session_credit) {
  SessionCredit credit;
  SessionCreditUpdateCriteria uc{};
  GrantedUnits gsu;
  uint64_t total_grant = 300;
  uint64_t tx_grant    = 100;
  uint64_t rx_grant    = 200;
  create_granted_units(&total_grant, &tx_grant, &rx_grant, &gsu);

  credit.receive_credit(gsu, &uc);

  EXPECT_EQ(300, credit.get_credit(ALLOWED_TOTAL));
  EXPECT_EQ(100, credit.get_credit(ALLOWED_TX));
  EXPECT_EQ(200, credit.get_credit(ALLOWED_RX));
  EXPECT_EQ(0, credit.get_credit(ALLOWED_FLOOR_TOTAL));
  EXPECT_EQ(0, credit.get_credit(ALLOWED_FLOOR_TX));
  EXPECT_EQ(0, credit.get_credit(ALLOWED_FLOOR_RX));
  EXPECT_EQ(0, credit.get_credit(USED_TX));
  EXPECT_EQ(0, credit.get_credit(USED_RX));

  // Check that the update criteria also includes the changes
  EXPECT_EQ(uc.bucket_deltas[ALLOWED_TOTAL], 300);
  EXPECT_EQ(uc.bucket_deltas[ALLOWED_TX], 100);
  EXPECT_EQ(uc.bucket_deltas[ALLOWED_RX], 200);
  EXPECT_EQ(uc.bucket_deltas[ALLOWED_FLOOR_TOTAL], 0);
  EXPECT_EQ(uc.bucket_deltas[ALLOWED_FLOOR_TX], 0);
  EXPECT_EQ(uc.bucket_deltas[ALLOWED_FLOOR_RX], 0);
  EXPECT_EQ(uc.bucket_deltas[USED_TX], 0);
}

TEST(test_add_received_credit, test_session_credit) {
  SessionCredit credit;
  SessionCreditUpdateCriteria uc{};
  GrantedUnits gsu;
  uint64_t grant = 1024;
  create_granted_units(&grant, NULL, NULL, &gsu);

  credit.receive_credit(gsu, &uc);
  credit.add_used_credit(10, 20, &uc);
  EXPECT_EQ(credit.get_credit(USED_TX), 10);
  EXPECT_EQ(credit.get_credit(USED_RX), 20);
  EXPECT_EQ(uc.bucket_deltas[USED_TX], 10);
  EXPECT_EQ(uc.bucket_deltas[USED_RX], 20);
  credit.add_used_credit(30, 40, &uc);
  EXPECT_EQ(credit.get_credit(USED_TX), 40);
  EXPECT_EQ(credit.get_credit(USED_RX), 60);
  EXPECT_EQ(uc.bucket_deltas[USED_TX], 40);
  EXPECT_EQ(uc.bucket_deltas[USED_RX], 60);
  RequestedUnits ru = credit.get_requested_credits_units();
  EXPECT_EQ(ru.total(), 100);
}

TEST(test_collect_updates, test_session_credit) {
  SessionCredit credit;
  SessionCreditUpdateCriteria uc{};
  GrantedUnits gsu;
  uint64_t grant = 1024;
  create_granted_units(&grant, NULL, NULL, &gsu);

  credit.receive_credit(gsu, &uc);
  credit.add_used_credit(500, 524, &uc);
  EXPECT_TRUE(credit.is_quota_exhausted(0.8));
  RequestedUnits ru = credit.get_requested_credits_units();
  EXPECT_EQ(ru.total(), 1024);

  auto update = credit.get_usage_for_reporting(&uc);
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
  credit.add_used_credit(300, 500, &uc);
  EXPECT_TRUE(credit.is_quota_exhausted(0.8));
  RequestedUnits ru = credit.get_requested_credits_units();
  EXPECT_EQ(ru.total(), 800);

  auto update = credit.get_usage_for_reporting(&uc);
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
  credit.add_used_credit(400, 399, &uc);
  EXPECT_FALSE(credit.is_quota_exhausted(0.8));
  RequestedUnits ru = credit.get_requested_credits_units();
  EXPECT_EQ(ru.total(), 799);
}

// The maximum of reported usage is NOT capped by what is granted even when an
// user overused.
TEST(test_collect_updates_when_overusing, test_session_credit) {
  SessionCredit credit;
  SessionCreditUpdateCriteria uc{};
  GrantedUnits gsu;
  uint64_t grant = 1000;
  create_granted_units(&grant, NULL, NULL, &gsu);

  credit.receive_credit(gsu, &uc);
  credit.add_used_credit(510, 500, &uc);
  EXPECT_TRUE(credit.is_quota_exhausted(0.8));
  RequestedUnits ru = credit.get_requested_credits_units();
  EXPECT_EQ(ru.total(), 1000);
  auto update = credit.get_usage_for_reporting(&uc);
  EXPECT_EQ(update.bytes_tx, 510);
  EXPECT_EQ(update.bytes_rx, 500);

  EXPECT_TRUE(credit.is_reporting());
  EXPECT_TRUE(uc.reporting);
  EXPECT_EQ(credit.get_credit(REPORTING_TX), 510);
  EXPECT_EQ(credit.get_credit(REPORTING_RX), 500);
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
  credit.add_used_credit(1000, 0, &uc);
  EXPECT_TRUE(credit.is_quota_exhausted(0.8));
  auto update = credit.get_usage_for_reporting(&uc);
  EXPECT_EQ(update.bytes_tx, 1000);
  EXPECT_EQ(update.bytes_rx, 0);
  RequestedUnits ru = credit.get_requested_credits_units();
  EXPECT_EQ(ru.tx(), 1000);
  EXPECT_EQ(ru.rx(), 0);
  EXPECT_EQ(ru.total(), 0);

  // receive another tx = 1000, rx = 1000, cumulative: tx = 2000, rx = 2000
  credit.receive_credit(gsu, &uc);
  EXPECT_EQ(uc.grant_tracking_type, TX_AND_RX);
  credit.add_used_credit(50, 2000, &uc);
  EXPECT_TRUE(credit.is_quota_exhausted(0.8));
  ru = credit.get_requested_credits_units();
  EXPECT_EQ(ru.tx(), 50);
  EXPECT_EQ(ru.rx(), 1000);
  EXPECT_EQ(ru.total(), 0);

  auto update2 = credit.get_usage_for_reporting(&uc);
  EXPECT_EQ(update2.bytes_tx, 50);
  EXPECT_EQ(update2.bytes_rx, 2000);

  // receive rx, tx, but no usage
  gsu.Clear();
  create_granted_units(NULL, &grant, &grant, &gsu);
  credit.receive_credit(gsu, &uc);
  EXPECT_FALSE(credit.is_quota_exhausted(0.8));
}

TEST(test_counting_algorithm, test_session_credit) {
  SessionCredit credit;
  SessionCreditUpdateCriteria uc{};
  GrantedUnits gsu;
  uint64_t total_grant = 1000;
  uint64_t tx_grant    = 100;
  uint64_t rx_grant    = 200;
  create_granted_units(&total_grant, &tx_grant, &rx_grant, &gsu);

  // receive total = 300 tx = 100, rx = 200
  credit.receive_credit(gsu, &uc);
  EXPECT_EQ(uc.grant_tracking_type, ALL_TOTAL_TX_RX);
  EXPECT_EQ(0, credit.get_credit(ALLOWED_FLOOR_TOTAL));
  EXPECT_EQ(0, credit.get_credit(ALLOWED_FLOOR_TX));
  EXPECT_EQ(0, credit.get_credit(ALLOWED_FLOOR_RX));
  EXPECT_EQ(1000, credit.get_credit(ALLOWED_TOTAL));
  EXPECT_EQ(
      100,
      credit.get_credit(ALLOWED_TX));  // 250 because we overused so 150 + 100
  EXPECT_EQ(200, credit.get_credit(ALLOWED_RX));

  // use tx and rx = 99 + 150 = 249
  credit.add_used_credit(99, 150, &uc);
  EXPECT_TRUE(credit.is_quota_exhausted(0.8));
  EXPECT_FALSE(credit.is_quota_exhausted(1));
  auto update = credit.get_usage_for_reporting(&uc);
  EXPECT_EQ(update.bytes_tx, 99);
  EXPECT_EQ(update.bytes_rx, 150);

  // receive another total = 300 tx = 100, rx = 200,
  // cumulative: total = 600 tx = 200, rx = 400
  credit.receive_credit(gsu, &uc);
  EXPECT_EQ(uc.grant_tracking_type, ALL_TOTAL_TX_RX);
  EXPECT_EQ(2000, credit.get_credit(ALLOWED_TOTAL));
  EXPECT_EQ(200, credit.get_credit(ALLOWED_TX));
  EXPECT_EQ(400, credit.get_credit(ALLOWED_RX));
  EXPECT_EQ(1000, credit.get_credit(ALLOWED_FLOOR_TOTAL));
  EXPECT_EQ(
      100,
      credit.get_credit(ALLOWED_FLOOR_TX));  // 150 because we overused so 150
  EXPECT_EQ(200, credit.get_credit(ALLOWED_FLOOR_RX));
  EXPECT_EQ(99, credit.get_credit(USED_TX));
  EXPECT_EQ(150, credit.get_credit(USED_RX));
  // use tx and rx = 1 + 50
  // cumulative = 100 + 200 = 300
  credit.add_used_credit(1, 50, &uc);
  EXPECT_FALSE(credit.is_quota_exhausted(0.8));
  // right before triggering the update (just under80% of the grant)
  // use tx and rx = 79 + 100
  // cumulative = 179 + 300 = 479
  credit.add_used_credit(79, 100, &uc);
  EXPECT_FALSE(credit.is_quota_exhausted(0.8));
  // right after triggering the update (just got over 80% of the grant)
  // use tx and rx = 1 + 0
  // cumulative = 180 + 300 = 460
  credit.add_used_credit(1, 0, &uc);
  EXPECT_TRUE(credit.is_quota_exhausted(0.8));
  EXPECT_FALSE(credit.is_quota_exhausted(1));
  // exhausted by tx
  // use tx and rx = 20 + 100
  // cumulative = 200 + 400 = 600
  credit.add_used_credit(20, 100, &uc);
  EXPECT_TRUE(credit.is_quota_exhausted(0.8));
  EXPECT_TRUE(credit.is_quota_exhausted(1));

  // receive one grant more with 0 granted units
  gsu.Clear();
  total_grant = 0;
  tx_grant    = 0;
  rx_grant    = 0;
  create_granted_units(&total_grant, &tx_grant, &rx_grant, &gsu);
  credit.receive_credit(gsu, &uc);
  EXPECT_EQ(uc.grant_tracking_type, ALL_TOTAL_TX_RX);
  EXPECT_EQ(2000, credit.get_credit(ALLOWED_TOTAL));
  EXPECT_EQ(200, credit.get_credit(ALLOWED_TX));
  EXPECT_EQ(400, credit.get_credit(ALLOWED_RX));
  EXPECT_EQ(1000, credit.get_credit(ALLOWED_FLOOR_TOTAL));
  EXPECT_EQ(100, credit.get_credit(ALLOWED_FLOOR_TX));
  EXPECT_EQ(200, credit.get_credit(ALLOWED_FLOOR_RX));
  EXPECT_EQ(200, credit.get_credit(USED_TX));
  EXPECT_EQ(400, credit.get_credit(USED_RX));
  EXPECT_TRUE(credit.is_quota_exhausted(1));
}

TEST(test_is_quota_exhausted_total_only, test_session_credit) {
  SessionCredit credit;
  SessionCreditUpdateCriteria uc{};

  GrantedUnits gsu;
  uint64_t total_grant = 1000;
  create_granted_units(&total_grant, NULL, NULL, &gsu);
  credit.receive_credit(gsu, &uc);
  EXPECT_EQ(uc.grant_tracking_type, TOTAL_ONLY);

  credit.add_used_credit(500, 0, &uc);
  EXPECT_FALSE(credit.is_quota_exhausted(0.8));
  credit.mark_failure(0, &uc);

  credit.add_used_credit(500, 0, &uc);
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

  credit.add_used_credit(500, 500, &uc);
  EXPECT_FALSE(credit.is_quota_exhausted(0.8));
  credit.mark_failure(0, &uc);

  credit.add_used_credit(500, 500, &uc);
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

  credit.add_used_credit(500, 500, &uc);
  EXPECT_FALSE(credit.is_quota_exhausted(0.8));
  credit.mark_failure(0, &uc);

  credit.add_used_credit(500, 500, &uc);
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
  credit.add_used_credit(900, 0, &uc);
  EXPECT_TRUE(credit.is_quota_exhausted(0.8));
  credit.mark_failure(0, &uc);

  // Receive empty GSU, quota should still be exhausted
  gsu.Clear();
  credit.receive_credit(gsu, &uc);
  // assert uc grant_tracking type has not changed
  EXPECT_EQ(uc.grant_tracking_type, TOTAL_ONLY);
  EXPECT_TRUE(credit.is_quota_exhausted(0.8));
}

TEST(test_get_credit_summary, test_session_credit) {
  SessionCredit credit;
  SessionCreditUpdateCriteria uc{};
  auto summary = credit.get_credit_summary();
  EXPECT_EQ(summary.usage.bytes_tx, 0);
  EXPECT_EQ(summary.usage.bytes_rx, 0);
  EXPECT_EQ(summary.time_of_first_usage, 0);
  EXPECT_EQ(summary.time_of_last_usage, 0);

  // Expect first_usage > 0 && first_usage == last_usage
  credit.add_used_credit(64, 32, &uc);
  summary = credit.get_credit_summary();
  EXPECT_EQ(summary.usage.bytes_tx, 64);
  EXPECT_EQ(summary.usage.bytes_rx, 32);
  EXPECT_NE(summary.time_of_first_usage, 0);
  auto time_of_first_usage = summary.time_of_first_usage;
  EXPECT_EQ(summary.time_of_last_usage, time_of_first_usage);

  sleep(1);
  // Expect first_usage > 0 && first_usage < last_usage
  credit.add_used_credit(64, 32, &uc);
  summary = credit.get_credit_summary();
  EXPECT_EQ(summary.usage.bytes_tx, 128);
  EXPECT_EQ(summary.usage.bytes_rx, 64);
  EXPECT_EQ(summary.time_of_first_usage, time_of_first_usage);
  EXPECT_NE(summary.time_of_last_usage, time_of_first_usage);
  EXPECT_GT(summary.time_of_last_usage, summary.time_of_first_usage);
}

int main(int argc, char** argv) {
  ::testing::InitGoogleTest(&argc, argv);
  return RUN_ALL_TESTS();
}

}  // namespace magma

/*
 * Licensed to the OpenAirInterface (OAI) Software Alliance under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The OpenAirInterface Software Alliance licenses this file to You under
 * the terms found in the LICENSE file in the root of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *-------------------------------------------------------------------------------
 * For more information about the OpenAirInterface (OAI) Software Alliance:
 *      contact@openairinterface.org
 */

#include <gtest/gtest.h>
#include <glog/logging.h>

extern "C" {
#include "conversions.h"
#include "ExtendedApnAggregateMaximumBitRate.h"
#include "3gpp_23.003.h"
#include "log.h"
}

TEST(test_extended_bit_rate_value_test, extended_bit_rate_value_test) {
  ExtendedApnAggregateMaximumBitRate apn_testing;

  extended_bit_rate_value(&apn_testing, 65000000000, 65000000000);
  ASSERT_EQ(apn_testing.extendedapnambrfordownlinkunit, 3);
  ASSERT_EQ(apn_testing.extendedapnambrforuplinkunit, 3);
  ASSERT_EQ(apn_testing.extendedapnambrfordownlink, 122);
  ASSERT_EQ(apn_testing.extendedapnambrfordownlink_continued, 63);
  ASSERT_EQ(apn_testing.extendedapnambrforuplink, 122);
  ASSERT_EQ(apn_testing.extendedapnambrforuplink_continued, 63);

  extended_bit_rate_value(&apn_testing, 100000000000, 100000000000);
  ASSERT_EQ(apn_testing.extendedapnambrfordownlinkunit, 3);
  ASSERT_EQ(apn_testing.extendedapnambrforuplinkunit, 3);
  ASSERT_EQ(apn_testing.extendedapnambrfordownlink, 168);
  ASSERT_EQ(apn_testing.extendedapnambrforuplink, 168);

  extended_bit_rate_value(&apn_testing, 1000000000000, 1000000000000);
  ASSERT_EQ(apn_testing.extendedapnambrfordownlinkunit, 4);
  ASSERT_EQ(apn_testing.extendedapnambrforuplinkunit, 4);
  ASSERT_EQ(apn_testing.extendedapnambrfordownlink, 36);
  ASSERT_EQ(apn_testing.extendedapnambrforuplink, 36);

  extended_bit_rate_value(&apn_testing, 2500000000000, 2500000000000);
  ASSERT_EQ(apn_testing.extendedapnambrfordownlinkunit, 5);
  ASSERT_EQ(apn_testing.extendedapnambrforuplinkunit, 5);
  ASSERT_EQ(apn_testing.extendedapnambrfordownlink, 150);
  ASSERT_EQ(apn_testing.extendedapnambrforuplink, 150);

  extended_bit_rate_value(&apn_testing, 250000000000000, 250000000000000);
  ASSERT_EQ(apn_testing.extendedapnambrfordownlinkunit, 8);
  ASSERT_EQ(apn_testing.extendedapnambrforuplinkunit, 8);
  ASSERT_EQ(apn_testing.extendedapnambrfordownlink, 107);
  ASSERT_EQ(apn_testing.extendedapnambrforuplink, 107);
}

int main(int argc, char** argv) {
  ::testing::InitGoogleTest(&argc, argv);
  OAILOG_INIT("MME", OAILOG_LEVEL_DEBUG, MAX_LOG_PROTOS);
  return RUN_ALL_TESTS();
}

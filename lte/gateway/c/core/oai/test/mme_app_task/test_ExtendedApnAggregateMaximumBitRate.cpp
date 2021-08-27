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
  ASSERT_EQ(apn_testing.extendedapnambrfordownlink, 122);
  ASSERT_EQ(apn_testing.extendedapnambrfordownlink_continued, 63);
  ASSERT_EQ(apn_testing.extendedapnambrforuplinkunit, 3);
  ASSERT_EQ(apn_testing.extendedapnambrforuplink, 122);
  ASSERT_EQ(apn_testing.extendedapnambrforuplink_continued, 63);

  extended_bit_rate_value(&apn_testing, 100000000000, 100000000000);
  ASSERT_EQ(apn_testing.extendedapnambrfordownlinkunit, 3);
  ASSERT_EQ(apn_testing.extendedapnambrfordownlink, 168);
  ASSERT_EQ(apn_testing.extendedapnambrfordownlink_continued, 97);
  ASSERT_EQ(apn_testing.extendedapnambrforuplinkunit, 3);
  ASSERT_EQ(apn_testing.extendedapnambrforuplink, 168);
  ASSERT_EQ(apn_testing.extendedapnambrforuplink_continued, 97);

  extended_bit_rate_value(&apn_testing, 262140000000, 262140000000);
  ASSERT_EQ(apn_testing.extendedapnambrfordownlinkunit, 3);
  ASSERT_EQ(apn_testing.extendedapnambrfordownlink, 255);
  ASSERT_EQ(apn_testing.extendedapnambrfordownlink_continued, 255);
  ASSERT_EQ(apn_testing.extendedapnambrforuplinkunit, 3);
  ASSERT_EQ(apn_testing.extendedapnambrforuplink, 255);
  ASSERT_EQ(apn_testing.extendedapnambrforuplink_continued, 255);

  extended_bit_rate_value(&apn_testing, 1000000000000, 1000000000000);
  ASSERT_EQ(apn_testing.extendedapnambrfordownlinkunit, 4);
  ASSERT_EQ(apn_testing.extendedapnambrfordownlink, 36);
  ASSERT_EQ(apn_testing.extendedapnambrfordownlink_continued, 244);
  ASSERT_EQ(apn_testing.extendedapnambrforuplinkunit, 4);
  ASSERT_EQ(apn_testing.extendedapnambrforuplink, 36);
  ASSERT_EQ(apn_testing.extendedapnambrforuplink_continued, 244);

  extended_bit_rate_value(&apn_testing, 1048560000000, 1048560000000);
  ASSERT_EQ(apn_testing.extendedapnambrfordownlinkunit, 4);
  ASSERT_EQ(apn_testing.extendedapnambrfordownlink, 255);
  ASSERT_EQ(apn_testing.extendedapnambrfordownlink_continued, 255);
  ASSERT_EQ(apn_testing.extendedapnambrforuplinkunit, 4);
  ASSERT_EQ(apn_testing.extendedapnambrforuplink, 255);
  ASSERT_EQ(apn_testing.extendedapnambrforuplink_continued, 255);

  extended_bit_rate_value(&apn_testing, 2500000000000, 2500000000000);
  ASSERT_EQ(apn_testing.extendedapnambrfordownlinkunit, 5);
  ASSERT_EQ(apn_testing.extendedapnambrforuplinkunit, 5);
  ASSERT_EQ(apn_testing.extendedapnambrfordownlink, 150);
  ASSERT_EQ(apn_testing.extendedapnambrfordownlink_continued, 152);
  ASSERT_EQ(apn_testing.extendedapnambrforuplink, 150);
  ASSERT_EQ(apn_testing.extendedapnambrforuplink_continued, 152);

  extended_bit_rate_value(&apn_testing, 4194240000000, 4194240000000);
  ASSERT_EQ(apn_testing.extendedapnambrfordownlinkunit, 5);
  ASSERT_EQ(apn_testing.extendedapnambrfordownlink, 255);
  ASSERT_EQ(apn_testing.extendedapnambrfordownlink_continued, 255);
  ASSERT_EQ(apn_testing.extendedapnambrforuplinkunit, 5);
  ASSERT_EQ(apn_testing.extendedapnambrforuplink, 255);
  ASSERT_EQ(apn_testing.extendedapnambrforuplink_continued, 255);

  extended_bit_rate_value(&apn_testing, 4200000000000, 4200000000000);
  ASSERT_EQ(apn_testing.extendedapnambrfordownlinkunit, 6);
  ASSERT_EQ(apn_testing.extendedapnambrfordownlink, 22);
  ASSERT_EQ(apn_testing.extendedapnambrfordownlink_continued, 64);
  ASSERT_EQ(apn_testing.extendedapnambrforuplinkunit, 6);
  ASSERT_EQ(apn_testing.extendedapnambrforuplink, 22);
  ASSERT_EQ(apn_testing.extendedapnambrforuplink_continued, 64);

  extended_bit_rate_value(&apn_testing, 16776960000000, 16776960000000);
  ASSERT_EQ(apn_testing.extendedapnambrfordownlinkunit, 6);
  ASSERT_EQ(apn_testing.extendedapnambrfordownlink, 255);
  ASSERT_EQ(apn_testing.extendedapnambrfordownlink_continued, 255);
  ASSERT_EQ(apn_testing.extendedapnambrforuplinkunit, 6);
  ASSERT_EQ(apn_testing.extendedapnambrforuplink, 255);
  ASSERT_EQ(apn_testing.extendedapnambrforuplink_continued, 255);

  extended_bit_rate_value(&apn_testing, 20000000000000, 20000000000000);
  ASSERT_EQ(apn_testing.extendedapnambrfordownlinkunit, 7);
  ASSERT_EQ(apn_testing.extendedapnambrfordownlink, 75);
  ASSERT_EQ(apn_testing.extendedapnambrfordownlink_continued, 76);
  ASSERT_EQ(apn_testing.extendedapnambrforuplinkunit, 7);
  ASSERT_EQ(apn_testing.extendedapnambrforuplink, 75);
  ASSERT_EQ(apn_testing.extendedapnambrforuplink_continued, 76);

  extended_bit_rate_value(&apn_testing, 67107840000000, 67107840000000);
  ASSERT_EQ(apn_testing.extendedapnambrfordownlinkunit, 7);
  ASSERT_EQ(apn_testing.extendedapnambrfordownlink, 255);
  ASSERT_EQ(apn_testing.extendedapnambrfordownlink_continued, 255);
  ASSERT_EQ(apn_testing.extendedapnambrforuplinkunit, 7);
  ASSERT_EQ(apn_testing.extendedapnambrforuplink, 255);
  ASSERT_EQ(apn_testing.extendedapnambrforuplink_continued, 255);

  extended_bit_rate_value(&apn_testing, 250000000000000, 250000000000000);
  ASSERT_EQ(apn_testing.extendedapnambrfordownlinkunit, 8);
  ASSERT_EQ(apn_testing.extendedapnambrfordownlink, 107);
  ASSERT_EQ(apn_testing.extendedapnambrfordownlink_continued, 238);
  ASSERT_EQ(apn_testing.extendedapnambrforuplinkunit, 8);
  ASSERT_EQ(apn_testing.extendedapnambrforuplink, 107);
  ASSERT_EQ(apn_testing.extendedapnambrforuplink_continued, 238);

  extended_bit_rate_value(&apn_testing, 268431360000000, 268431360000000);
  ASSERT_EQ(apn_testing.extendedapnambrfordownlinkunit, 8);
  ASSERT_EQ(apn_testing.extendedapnambrfordownlink, 255);
  ASSERT_EQ(apn_testing.extendedapnambrfordownlink_continued, 255);
  ASSERT_EQ(apn_testing.extendedapnambrforuplinkunit, 8);
  ASSERT_EQ(apn_testing.extendedapnambrforuplink, 255);
  ASSERT_EQ(apn_testing.extendedapnambrforuplink_continued, 255);

  extended_bit_rate_value(&apn_testing, 268435456000000, 268435456000000);
  ASSERT_EQ(apn_testing.extendedapnambrfordownlinkunit, 9);
  ASSERT_EQ(apn_testing.extendedapnambrfordownlink, 0);
  ASSERT_EQ(apn_testing.extendedapnambrfordownlink_continued, 64);
  ASSERT_EQ(apn_testing.extendedapnambrforuplinkunit, 9);
  ASSERT_EQ(apn_testing.extendedapnambrforuplink, 0);
  ASSERT_EQ(apn_testing.extendedapnambrforuplink_continued, 64);

  extended_bit_rate_value(&apn_testing, 1073725440000000, 1073725440000000);
  ASSERT_EQ(apn_testing.extendedapnambrfordownlinkunit, 9);
  ASSERT_EQ(apn_testing.extendedapnambrfordownlink, 255);
  ASSERT_EQ(apn_testing.extendedapnambrfordownlink_continued, 255);
  ASSERT_EQ(apn_testing.extendedapnambrforuplinkunit, 9);
  ASSERT_EQ(apn_testing.extendedapnambrforuplink, 255);
  ASSERT_EQ(apn_testing.extendedapnambrforuplink_continued, 255);

  extended_bit_rate_value(&apn_testing, 3000000000000000, 3000000000000000);
  ASSERT_EQ(apn_testing.extendedapnambrfordownlinkunit, 10);
  ASSERT_EQ(apn_testing.extendedapnambrfordownlink, 208);
  ASSERT_EQ(apn_testing.extendedapnambrfordownlink_continued, 178);
  ASSERT_EQ(apn_testing.extendedapnambrforuplinkunit, 10);
  ASSERT_EQ(apn_testing.extendedapnambrforuplink, 208);
  ASSERT_EQ(apn_testing.extendedapnambrforuplink_continued, 178);

  extended_bit_rate_value(&apn_testing, 4294901760000000, 4294901760000000);
  ASSERT_EQ(apn_testing.extendedapnambrfordownlinkunit, 10);
  ASSERT_EQ(apn_testing.extendedapnambrfordownlink, 255);
  ASSERT_EQ(apn_testing.extendedapnambrfordownlink_continued, 255);
  ASSERT_EQ(apn_testing.extendedapnambrforuplinkunit, 10);
  ASSERT_EQ(apn_testing.extendedapnambrforuplink, 255);
  ASSERT_EQ(apn_testing.extendedapnambrforuplink_continued, 255);

  extended_bit_rate_value(&apn_testing, 6000000000000000, 6000000000000000);
  ASSERT_EQ(apn_testing.extendedapnambrfordownlinkunit, 11);
  ASSERT_EQ(apn_testing.extendedapnambrfordownlink, 104);
  ASSERT_EQ(apn_testing.extendedapnambrfordownlink_continued, 89);
  ASSERT_EQ(apn_testing.extendedapnambrforuplinkunit, 11);
  ASSERT_EQ(apn_testing.extendedapnambrforuplink, 104);
  ASSERT_EQ(apn_testing.extendedapnambrforuplink_continued, 89);

  extended_bit_rate_value(&apn_testing, 17179607040000000, 17179607040000000);
  ASSERT_EQ(apn_testing.extendedapnambrfordownlinkunit, 11);
  ASSERT_EQ(apn_testing.extendedapnambrfordownlink, 255);
  ASSERT_EQ(apn_testing.extendedapnambrfordownlink_continued, 255);
  ASSERT_EQ(apn_testing.extendedapnambrforuplinkunit, 11);
  ASSERT_EQ(apn_testing.extendedapnambrforuplink, 255);
  ASSERT_EQ(apn_testing.extendedapnambrforuplink_continued, 255);

  extended_bit_rate_value(&apn_testing, 19000000000000000, 19000000000000000);
  ASSERT_EQ(apn_testing.extendedapnambrfordownlinkunit, 12);
  ASSERT_EQ(apn_testing.extendedapnambrfordownlink, 199);
  ASSERT_EQ(apn_testing.extendedapnambrfordownlink_continued, 70);
  ASSERT_EQ(apn_testing.extendedapnambrforuplinkunit, 12);
  ASSERT_EQ(apn_testing.extendedapnambrforuplink, 199);
  ASSERT_EQ(apn_testing.extendedapnambrforuplink_continued, 70);

  extended_bit_rate_value(&apn_testing, 20000000000000000, 20000000000000000);
  ASSERT_EQ(apn_testing.extendedapnambrfordownlinkunit, 12);
  ASSERT_EQ(apn_testing.extendedapnambrfordownlink, 129);
  ASSERT_EQ(apn_testing.extendedapnambrfordownlink_continued, 74);
  ASSERT_EQ(apn_testing.extendedapnambrforuplinkunit, 12);
  ASSERT_EQ(apn_testing.extendedapnambrforuplink, 129);
  ASSERT_EQ(apn_testing.extendedapnambrforuplink_continued, 74);

  extended_bit_rate_value(&apn_testing, 68718428160000000, 68718428160000000);
  ASSERT_EQ(apn_testing.extendedapnambrfordownlinkunit, 12);
  ASSERT_EQ(apn_testing.extendedapnambrfordownlink, 255);
  ASSERT_EQ(apn_testing.extendedapnambrfordownlink_continued, 255);
  ASSERT_EQ(apn_testing.extendedapnambrforuplinkunit, 12);
  ASSERT_EQ(apn_testing.extendedapnambrforuplink, 255);
  ASSERT_EQ(apn_testing.extendedapnambrforuplink_continued, 255);

  extended_bit_rate_value(&apn_testing, 270000000000000000, 270000000000000000);
  ASSERT_EQ(apn_testing.extendedapnambrfordownlinkunit, 13);
  ASSERT_EQ(apn_testing.extendedapnambrfordownlink, 117);
  ASSERT_EQ(apn_testing.extendedapnambrfordownlink_continued, 251);
  ASSERT_EQ(apn_testing.extendedapnambrforuplinkunit, 13);
  ASSERT_EQ(apn_testing.extendedapnambrforuplink, 117);
  ASSERT_EQ(apn_testing.extendedapnambrforuplink_continued, 251);

  extended_bit_rate_value(&apn_testing, 274873712640000000, 274873712640000000);
  ASSERT_EQ(apn_testing.extendedapnambrfordownlinkunit, 13);
  ASSERT_EQ(apn_testing.extendedapnambrfordownlink, 255);
  ASSERT_EQ(apn_testing.extendedapnambrfordownlink_continued, 255);
  ASSERT_EQ(apn_testing.extendedapnambrforuplinkunit, 13);
  ASSERT_EQ(apn_testing.extendedapnambrforuplink, 255);
  ASSERT_EQ(apn_testing.extendedapnambrforuplink_continued, 255);

  extended_bit_rate_value(
      &apn_testing, 1000000000000000000, 1000000000000000000);
  ASSERT_EQ(apn_testing.extendedapnambrfordownlinkunit, 14);
  ASSERT_EQ(apn_testing.extendedapnambrfordownlink, 212);
  ASSERT_EQ(apn_testing.extendedapnambrfordownlink_continued, 232);
  ASSERT_EQ(apn_testing.extendedapnambrforuplinkunit, 14);
  ASSERT_EQ(apn_testing.extendedapnambrforuplink, 212);
  ASSERT_EQ(apn_testing.extendedapnambrforuplink_continued, 232);

  extended_bit_rate_value(
      &apn_testing, 1099494851000000000, 1099494851000000000);
  ASSERT_EQ(apn_testing.extendedapnambrfordownlinkunit, 14);
  ASSERT_EQ(apn_testing.extendedapnambrfordownlink, 255);
  ASSERT_EQ(apn_testing.extendedapnambrfordownlink_continued, 255);
  ASSERT_EQ(apn_testing.extendedapnambrforuplinkunit, 14);
  ASSERT_EQ(apn_testing.extendedapnambrforuplink, 255);
  ASSERT_EQ(apn_testing.extendedapnambrforuplink_continued, 255);

  extended_bit_rate_value(
      &apn_testing, 4000000000000000000, 4000000000000000000);
  ASSERT_EQ(apn_testing.extendedapnambrfordownlinkunit, 15);
  ASSERT_EQ(apn_testing.extendedapnambrfordownlink, 212);
  ASSERT_EQ(apn_testing.extendedapnambrfordownlink_continued, 232);
  ASSERT_EQ(apn_testing.extendedapnambrforuplinkunit, 15);
  ASSERT_EQ(apn_testing.extendedapnambrforuplink, 212);
  ASSERT_EQ(apn_testing.extendedapnambrforuplink_continued, 232);

  extended_bit_rate_value(
      &apn_testing, 4397979402000000000, 4397979402000000000);
  ASSERT_EQ(apn_testing.extendedapnambrfordownlinkunit, 15);
  ASSERT_EQ(apn_testing.extendedapnambrfordownlink, 254);
  ASSERT_EQ(apn_testing.extendedapnambrfordownlink_continued, 255);
  ASSERT_EQ(apn_testing.extendedapnambrforuplinkunit, 15);
  ASSERT_EQ(apn_testing.extendedapnambrforuplink, 254);
  ASSERT_EQ(apn_testing.extendedapnambrforuplink_continued, 255);

  extended_bit_rate_value(
      &apn_testing, 15000000000000000000, 15000000000000000000);
  ASSERT_EQ(apn_testing.extendedapnambrfordownlinkunit, 16);
  ASSERT_EQ(apn_testing.extendedapnambrfordownlink, 71);
  ASSERT_EQ(apn_testing.extendedapnambrfordownlink_continued, 218);
  ASSERT_EQ(apn_testing.extendedapnambrforuplinkunit, 16);
  ASSERT_EQ(apn_testing.extendedapnambrforuplink, 71);
  ASSERT_EQ(apn_testing.extendedapnambrforuplink_continued, 218);

  extended_bit_rate_value(
      &apn_testing, 17591917610000000000, 17591917610000000000);
  ASSERT_EQ(apn_testing.extendedapnambrfordownlinkunit, 16);
  ASSERT_EQ(apn_testing.extendedapnambrfordownlink, 255);
  ASSERT_EQ(apn_testing.extendedapnambrfordownlink_continued, 255);
  ASSERT_EQ(apn_testing.extendedapnambrforuplinkunit, 16);
  ASSERT_EQ(apn_testing.extendedapnambrforuplink, 255);
  ASSERT_EQ(apn_testing.extendedapnambrforuplink_continued, 255);
}

int main(int argc, char** argv) {
  ::testing::InitGoogleTest(&argc, argv);
  OAILOG_INIT("MME", OAILOG_LEVEL_DEBUG, MAX_LOG_PROTOS);
  return RUN_ALL_TESTS();
}

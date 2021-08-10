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
#include "ApnAggregateMaximumBitRate.h"
#include "3gpp_23.003.h"
#include "log.h"
}

TEST(
    test_bit_rate_value_to_eps_qos_for_apnambr_extended2_test,
    bit_rate_value_to_eps_qos_for_apnambr_extended2_test) {
  ApnAggregateMaximumBitRate apn_testing;

  bit_rate_value_to_eps_qos(&apn_testing, 0, 0);
  ASSERT_EQ(apn_testing.apnambrfordownlink, 255);
  ASSERT_EQ(apn_testing.apnambrforuplink, 255);
  ASSERT_EQ(apn_testing.apnambrforuplink_extended, 0);
  ASSERT_EQ(apn_testing.apnambrfordownlink_extended, 0);
  ASSERT_EQ(apn_testing.apnambrforuplink_extended2, 0);
  ASSERT_EQ(apn_testing.apnambrfordownlink_extended2, 0);

  bit_rate_value_to_eps_qos(&apn_testing, 1000, 1000);
  ASSERT_EQ(apn_testing.apnambrfordownlink, 1);
  ASSERT_EQ(apn_testing.apnambrforuplink, 1);
  ASSERT_EQ(apn_testing.apnambrforuplink_extended, 0);
  ASSERT_EQ(apn_testing.apnambrfordownlink_extended, 0);
  ASSERT_EQ(apn_testing.apnambrforuplink_extended2, 0);
  ASSERT_EQ(apn_testing.apnambrfordownlink_extended2, 0);

  bit_rate_value_to_eps_qos(&apn_testing, 63000, 63000);
  ASSERT_EQ(apn_testing.apnambrfordownlink, 63);
  ASSERT_EQ(apn_testing.apnambrforuplink, 63);
  ASSERT_EQ(apn_testing.apnambrforuplink_extended, 0);
  ASSERT_EQ(apn_testing.apnambrfordownlink_extended, 0);
  ASSERT_EQ(apn_testing.apnambrforuplink_extended2, 0);
  ASSERT_EQ(apn_testing.apnambrfordownlink_extended2, 0);

  bit_rate_value_to_eps_qos(&apn_testing, 64000, 64000);
  ASSERT_EQ(apn_testing.apnambrfordownlink, 64);
  ASSERT_EQ(apn_testing.apnambrforuplink, 64);
  ASSERT_EQ(apn_testing.apnambrforuplink_extended, 0);
  ASSERT_EQ(apn_testing.apnambrfordownlink_extended, 0);
  ASSERT_EQ(apn_testing.apnambrforuplink_extended2, 0);
  ASSERT_EQ(apn_testing.apnambrfordownlink_extended2, 0);

  bit_rate_value_to_eps_qos(&apn_testing, 72000, 72000);
  ASSERT_EQ(apn_testing.apnambrfordownlink, 65);
  ASSERT_EQ(apn_testing.apnambrforuplink, 65);
  ASSERT_EQ(apn_testing.apnambrforuplink_extended, 0);
  ASSERT_EQ(apn_testing.apnambrfordownlink_extended, 0);
  ASSERT_EQ(apn_testing.apnambrforuplink_extended2, 0);
  ASSERT_EQ(apn_testing.apnambrfordownlink_extended2, 0);

  bit_rate_value_to_eps_qos(&apn_testing, 568000, 568000);
  ASSERT_EQ(apn_testing.apnambrfordownlink, 127);
  ASSERT_EQ(apn_testing.apnambrforuplink, 127);
  ASSERT_EQ(apn_testing.apnambrforuplink_extended, 0);
  ASSERT_EQ(apn_testing.apnambrfordownlink_extended, 0);
  ASSERT_EQ(apn_testing.apnambrforuplink_extended2, 0);
  ASSERT_EQ(apn_testing.apnambrfordownlink_extended2, 0);

  bit_rate_value_to_eps_qos(&apn_testing, 576000, 576000);
  ASSERT_EQ(apn_testing.apnambrfordownlink, 128);
  ASSERT_EQ(apn_testing.apnambrforuplink, 128);
  ASSERT_EQ(apn_testing.apnambrforuplink_extended, 0);
  ASSERT_EQ(apn_testing.apnambrfordownlink_extended, 0);
  ASSERT_EQ(apn_testing.apnambrforuplink_extended2, 0);
  ASSERT_EQ(apn_testing.apnambrfordownlink_extended2, 0);

  bit_rate_value_to_eps_qos(&apn_testing, 1152000, 1152000);
  ASSERT_EQ(apn_testing.apnambrfordownlink, 137);
  ASSERT_EQ(apn_testing.apnambrforuplink, 137);
  ASSERT_EQ(apn_testing.apnambrforuplink_extended, 0);
  ASSERT_EQ(apn_testing.apnambrfordownlink_extended, 0);
  ASSERT_EQ(apn_testing.apnambrforuplink_extended2, 0);
  ASSERT_EQ(apn_testing.apnambrfordownlink_extended2, 0);

  bit_rate_value_to_eps_qos(&apn_testing, 8640000, 8640000);
  ASSERT_EQ(apn_testing.apnambrfordownlink, 254);
  ASSERT_EQ(apn_testing.apnambrforuplink, 254);
  ASSERT_EQ(apn_testing.apnambrforuplink_extended, 0);
  ASSERT_EQ(apn_testing.apnambrfordownlink_extended, 0);
  ASSERT_EQ(apn_testing.apnambrforuplink_extended2, 0);
  ASSERT_EQ(apn_testing.apnambrfordownlink_extended2, 0);

  bit_rate_value_to_eps_qos(&apn_testing, 8700000, 8700000);
  ASSERT_EQ(apn_testing.apnambrfordownlink, 255);
  ASSERT_EQ(apn_testing.apnambrforuplink, 255);
  ASSERT_EQ(apn_testing.apnambrforuplink_extended, 1);
  ASSERT_EQ(apn_testing.apnambrfordownlink_extended, 1);
  ASSERT_EQ(apn_testing.apnambrforuplink_extended2, 0);
  ASSERT_EQ(apn_testing.apnambrfordownlink_extended2, 0);

  bit_rate_value_to_eps_qos(&apn_testing, 9000000, 9000000);
  ASSERT_EQ(apn_testing.apnambrfordownlink, 255);
  ASSERT_EQ(apn_testing.apnambrforuplink, 255);
  ASSERT_EQ(apn_testing.apnambrforuplink_extended, 4);
  ASSERT_EQ(apn_testing.apnambrfordownlink_extended, 4);
  ASSERT_EQ(apn_testing.apnambrforuplink_extended2, 0);
  ASSERT_EQ(apn_testing.apnambrfordownlink_extended2, 0);

  bit_rate_value_to_eps_qos(&apn_testing, 16000000, 16000000);
  ASSERT_EQ(apn_testing.apnambrfordownlink, 255);
  ASSERT_EQ(apn_testing.apnambrforuplink, 255);
  ASSERT_EQ(apn_testing.apnambrforuplink_extended, 74);
  ASSERT_EQ(apn_testing.apnambrfordownlink_extended, 74);
  ASSERT_EQ(apn_testing.apnambrforuplink_extended2, 0);
  ASSERT_EQ(apn_testing.apnambrfordownlink_extended2, 0);

  bit_rate_value_to_eps_qos(&apn_testing, 17000000, 17000000);
  ASSERT_EQ(apn_testing.apnambrfordownlink, 255);
  ASSERT_EQ(apn_testing.apnambrforuplink, 255);
  ASSERT_EQ(apn_testing.apnambrforuplink_extended, 75);
  ASSERT_EQ(apn_testing.apnambrfordownlink_extended, 75);
  ASSERT_EQ(apn_testing.apnambrforuplink_extended2, 0);
  ASSERT_EQ(apn_testing.apnambrfordownlink_extended2, 0);

  bit_rate_value_to_eps_qos(&apn_testing, 20000000, 20000000);
  ASSERT_EQ(apn_testing.apnambrfordownlink, 255);
  ASSERT_EQ(apn_testing.apnambrforuplink, 255);
  ASSERT_EQ(apn_testing.apnambrforuplink_extended, 78);
  ASSERT_EQ(apn_testing.apnambrfordownlink_extended, 78);
  ASSERT_EQ(apn_testing.apnambrforuplink_extended2, 0);
  ASSERT_EQ(apn_testing.apnambrfordownlink_extended2, 0);

  bit_rate_value_to_eps_qos(&apn_testing, 128000000, 128000000);
  ASSERT_EQ(apn_testing.apnambrfordownlink, 255);
  ASSERT_EQ(apn_testing.apnambrforuplink, 255);
  ASSERT_EQ(apn_testing.apnambrforuplink_extended, 186);
  ASSERT_EQ(apn_testing.apnambrfordownlink_extended, 186);
  ASSERT_EQ(apn_testing.apnambrforuplink_extended2, 0);
  ASSERT_EQ(apn_testing.apnambrfordownlink_extended2, 0);

  bit_rate_value_to_eps_qos(&apn_testing, 130000000, 130000000);
  ASSERT_EQ(apn_testing.apnambrfordownlink, 255);
  ASSERT_EQ(apn_testing.apnambrforuplink, 255);
  ASSERT_EQ(apn_testing.apnambrforuplink_extended, 187);
  ASSERT_EQ(apn_testing.apnambrfordownlink_extended, 187);
  ASSERT_EQ(apn_testing.apnambrforuplink_extended2, 0);
  ASSERT_EQ(apn_testing.apnambrfordownlink_extended2, 0);

  bit_rate_value_to_eps_qos(&apn_testing, 200000000, 200000000);
  ASSERT_EQ(apn_testing.apnambrfordownlink, 255);
  ASSERT_EQ(apn_testing.apnambrforuplink, 255);
  ASSERT_EQ(apn_testing.apnambrforuplink_extended, 222);
  ASSERT_EQ(apn_testing.apnambrfordownlink_extended, 222);
  ASSERT_EQ(apn_testing.apnambrforuplink_extended2, 0);
  ASSERT_EQ(apn_testing.apnambrfordownlink_extended2, 0);

  bit_rate_value_to_eps_qos(&apn_testing, 256000000, 256000000);
  ASSERT_EQ(apn_testing.apnambrfordownlink, 255);
  ASSERT_EQ(apn_testing.apnambrforuplink, 255);
  ASSERT_EQ(apn_testing.apnambrforuplink_extended, 250);
  ASSERT_EQ(apn_testing.apnambrfordownlink_extended, 250);
  ASSERT_EQ(apn_testing.apnambrforuplink_extended2, 0);
  ASSERT_EQ(apn_testing.apnambrfordownlink_extended2, 0);

  bit_rate_value_to_eps_qos(&apn_testing, 300000000, 300000000);
  ASSERT_EQ(apn_testing.apnambrfordownlink, 255);
  ASSERT_EQ(apn_testing.apnambrforuplink, 255);
  ASSERT_EQ(apn_testing.apnambrforuplink_extended, 102);
  ASSERT_EQ(apn_testing.apnambrfordownlink_extended, 102);
  ASSERT_EQ(apn_testing.apnambrforuplink_extended2, 1);
  ASSERT_EQ(apn_testing.apnambrfordownlink_extended2, 1);

  bit_rate_value_to_eps_qos(&apn_testing, 310540000, 310540000);
  ASSERT_EQ(apn_testing.apnambrfordownlink, 123);
  ASSERT_EQ(apn_testing.apnambrforuplink, 123);
  ASSERT_EQ(apn_testing.apnambrforuplink_extended, 112);
  ASSERT_EQ(apn_testing.apnambrfordownlink_extended, 112);
  ASSERT_EQ(apn_testing.apnambrforuplink_extended2, 1);
  ASSERT_EQ(apn_testing.apnambrfordownlink_extended2, 1);

  bit_rate_value_to_eps_qos(&apn_testing, 522005000, 522005000);
  ASSERT_EQ(apn_testing.apnambrfordownlink, 5);
  ASSERT_EQ(apn_testing.apnambrforuplink, 5);
  ASSERT_EQ(apn_testing.apnambrforuplink_extended, 14);
  ASSERT_EQ(apn_testing.apnambrfordownlink_extended, 14);
  ASSERT_EQ(apn_testing.apnambrforuplink_extended2, 2);
  ASSERT_EQ(apn_testing.apnambrfordownlink_extended2, 2);

  bit_rate_value_to_eps_qos(&apn_testing, 65280000000, 65280000000);
  ASSERT_EQ(apn_testing.apnambrfordownlink, 255);
  ASSERT_EQ(apn_testing.apnambrforuplink, 255);
  ASSERT_EQ(apn_testing.apnambrforuplink_extended, 0);
  ASSERT_EQ(apn_testing.apnambrfordownlink_extended, 0);
  ASSERT_EQ(apn_testing.apnambrforuplink_extended2, 255);
  ASSERT_EQ(apn_testing.apnambrfordownlink_extended2, 255);
}

int main(int argc, char** argv) {
  ::testing::InitGoogleTest(&argc, argv);
  OAILOG_INIT("MME", OAILOG_LEVEL_DEBUG, MAX_LOG_PROTOS);
  return RUN_ALL_TESTS();
}

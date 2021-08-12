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

  extended_bit_rate_value(&apn_testing, 0, 0);
  // ASSERT_EQ(apn_testing.apnambrfordownlink, 255);
  // ASSERT_EQ(apn_testing.apnambrforuplink, 255);
  // ASSERT_EQ(apn_testing.apnambrforuplink_extended, 0);
  // ASSERT_EQ(apn_testing.apnambrfordownlink_extended, 0);
  // ASSERT_EQ(apn_testing.apnambrforuplink_extended2, 0);
  // ASSERT_EQ(apn_testing.apnambrfordownlink_extended2, 0);
}

int main(int argc, char** argv) {
  ::testing::InitGoogleTest(&argc, argv);
  OAILOG_INIT("MME", OAILOG_LEVEL_DEBUG, MAX_LOG_PROTOS);
  return RUN_ALL_TESTS();
}

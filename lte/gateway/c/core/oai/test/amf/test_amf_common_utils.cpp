/*
 * Copyright 2020 The Magma Authors.
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

#include <gtest/gtest.h>
#include "lte/gateway/c/core/oai/tasks/amf/amf_app_defs.hpp"

using ::testing::Test;

namespace magma5g {

TEST(test_ambr, test_ambr_convert) {
  uint32_t pdu_ambr_response_unit;
  uint32_t pdu_ambr_response_value;
  M5GSessionAmbrUnit ambr_unit;
  uint16_t ambr_value;

  pdu_ambr_response_unit = 0;     // BPS
  pdu_ambr_response_value = 500;  // 500 BPS

  convert_ambr(&pdu_ambr_response_unit, &pdu_ambr_response_value, &ambr_unit,
               &ambr_value);
  EXPECT_EQ(ambr_unit,
            M5GSessionAmbrUnit::MULTIPLES_1KBPS);  // Illegal ambr case
  EXPECT_EQ(ambr_value, 1);                        // 1 KBPS

  pdu_ambr_response_unit = 1;     // KBPS
  pdu_ambr_response_value = 500;  // 500 KBPS

  convert_ambr(&pdu_ambr_response_unit, &pdu_ambr_response_value, &ambr_unit,
               &ambr_value);
  EXPECT_EQ(ambr_unit, M5GSessionAmbrUnit::MULTIPLES_1KBPS);  // Kbps
  EXPECT_EQ(ambr_value, 500);                                 // 500 Kbps

  pdu_ambr_response_unit = 0;      // BPS
  pdu_ambr_response_value = 5000;  // 5000 BPS

  convert_ambr(&pdu_ambr_response_unit, &pdu_ambr_response_value, &ambr_unit,
               &ambr_value);
  EXPECT_EQ(ambr_unit, M5GSessionAmbrUnit::MULTIPLES_1KBPS);  // Kbps
  EXPECT_EQ(ambr_value, 5);                                   // 5 Kbps

  pdu_ambr_response_unit = 1;      // KBPS
  pdu_ambr_response_value = 5000;  // 5000 KBPS

  convert_ambr(&pdu_ambr_response_unit, &pdu_ambr_response_value, &ambr_unit,
               &ambr_value);
  EXPECT_EQ(ambr_unit, M5GSessionAmbrUnit::MULTIPLES_1KBPS);  // KBPS
  EXPECT_EQ(ambr_value, 5000);                                // 5000 KBPS

  pdu_ambr_response_unit = 0;           // BPS
  pdu_ambr_response_value = 536870000;  // 536870000 BPS

  convert_ambr(&pdu_ambr_response_unit, &pdu_ambr_response_value, &ambr_unit,
               &ambr_value);
  EXPECT_EQ(ambr_unit, M5GSessionAmbrUnit::MULTIPLES_1MBPS);  // Mbps
  EXPECT_EQ(ambr_value, 536);                                 // 536 Mbps
}

TEST(test_ambr, test_ambr_convert_pdu_aggregrate) {
  uint64_t dl_pdu_ambr;
  uint64_t ul_pdu_ambr;

  uint16_t dl_session_ambr = 5000;
  uint16_t ul_session_ambr = 6000;
  M5GSessionAmbrUnit dl_ambr_unit = M5GSessionAmbrUnit::MULTIPLES_1KBPS;
  M5GSessionAmbrUnit ul_ambr_unit = M5GSessionAmbrUnit::MULTIPLES_1KBPS;

  ambr_calculation_pdu_session(&dl_session_ambr, dl_ambr_unit, &ul_session_ambr,
                               ul_ambr_unit, &dl_pdu_ambr, &ul_pdu_ambr);
  EXPECT_EQ(dl_pdu_ambr, 5000000);
  EXPECT_EQ(ul_pdu_ambr, 6000000);

  dl_session_ambr = 5000;
  ul_session_ambr = 6000;
  dl_ambr_unit = M5GSessionAmbrUnit::MULTIPLES_1MBPS;
  ul_ambr_unit = M5GSessionAmbrUnit::MULTIPLES_1MBPS;

  ambr_calculation_pdu_session(&dl_session_ambr, dl_ambr_unit, &ul_session_ambr,
                               ul_ambr_unit, &dl_pdu_ambr, &ul_pdu_ambr);
  EXPECT_EQ(dl_pdu_ambr, 5000000000);
  EXPECT_EQ(ul_pdu_ambr, 6000000000);
}

}  // namespace magma5g

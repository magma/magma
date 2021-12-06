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
#include "lte/gateway/c/core/oai/tasks/amf/amf_app_defs.h"

using ::testing::Test;

namespace magma5g {

TEST(test_ambr, test_ambr_convert) {
  uint32_t pdu_ambr_response_unit;
  uint32_t pdu_ambr_response_value;
  uint8_t ambr_unit;
  uint16_t ambr_value;

  pdu_ambr_response_unit  = 0;    // BPS
  pdu_ambr_response_value = 500;  // 500 BPS

  convert_ambr(
      &pdu_ambr_response_unit, &pdu_ambr_response_value, &ambr_unit,
      &ambr_value);
  EXPECT_EQ(ambr_unit, 1);   // Illegal ambr case
  EXPECT_EQ(ambr_value, 1);  // 1 KBPS

  pdu_ambr_response_unit  = 1;    // KBPS
  pdu_ambr_response_value = 500;  // 500 KBPS

  convert_ambr(
      &pdu_ambr_response_unit, &pdu_ambr_response_value, &ambr_unit,
      &ambr_value);
  EXPECT_EQ(ambr_unit, 1);     // Kbps
  EXPECT_EQ(ambr_value, 500);  // 500 Kbps

  pdu_ambr_response_unit  = 0;     // BPS
  pdu_ambr_response_value = 5000;  // 5000 BPS

  convert_ambr(
      &pdu_ambr_response_unit, &pdu_ambr_response_value, &ambr_unit,
      &ambr_value);
  EXPECT_EQ(ambr_unit, 1);   // Kbps
  EXPECT_EQ(ambr_value, 5);  // 5 Kbps

  pdu_ambr_response_unit  = 1;     // KBPS
  pdu_ambr_response_value = 5000;  // 5000 KBPS

  convert_ambr(
      &pdu_ambr_response_unit, &pdu_ambr_response_value, &ambr_unit,
      &ambr_value);
  EXPECT_EQ(ambr_unit, 1);      // KBPS
  EXPECT_EQ(ambr_value, 5000);  // 5000 KBPS

  pdu_ambr_response_unit  = 0;          // BPS
  pdu_ambr_response_value = 536870000;  // 536870000 BPS

  convert_ambr(
      &pdu_ambr_response_unit, &pdu_ambr_response_value, &ambr_unit,
      &ambr_value);
  EXPECT_EQ(ambr_unit, 6);     // Mbps
  EXPECT_EQ(ambr_value, 536);  // 536 Mbps
}

}  // namespace magma5g

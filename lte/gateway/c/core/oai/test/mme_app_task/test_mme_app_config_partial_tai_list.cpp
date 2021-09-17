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

extern "C" {
#include "include/mme_config.h"
}

namespace magma {
namespace lte {

class MMEConfigTest : public ::testing::Test {
 protected:
  std::array<int, 25> ncon_tac;
  virtual void SetUp() {
    ncon_tac = {1,  2,  4,  6,  7,  8,  10, 11, 12, 14, 15,
                16, 19, 21, 23, 26, 28, 31, 33, 37, 39};
  }
  virtual void TearDown() {}
};

// Test partial list with 1 TAI
TEST_F(MMEConfigTest, TestOneTai) {
  mme_config_t config_pP        = {0};
  config_pP.served_tai.nb_tai   = 1;
  uint8_t itr                   = 0;
  uint16_t tac                  = 0;
  config_pP.served_tai.plmn_mcc = reinterpret_cast<uint16_t*>(
      calloc(config_pP.served_tai.nb_tai, sizeof(uint16_t)));
  config_pP.served_tai.plmn_mnc = reinterpret_cast<uint16_t*>(
      calloc(config_pP.served_tai.nb_tai, sizeof(uint16_t)));
  config_pP.served_tai.plmn_mnc_len = reinterpret_cast<uint16_t*>(
      calloc(config_pP.served_tai.nb_tai, sizeof(uint16_t)));
  config_pP.served_tai.tac = reinterpret_cast<uint16_t*>(
      calloc(config_pP.served_tai.nb_tai, sizeof(uint16_t)));
  config_pP.served_tai.plmn_mcc[itr]     = 1;
  config_pP.served_tai.plmn_mnc[itr]     = 1;
  config_pP.served_tai.plmn_mnc_len[itr] = 2;
  config_pP.served_tai.tac[itr]          = 1;
  create_partial_lists(&config_pP);
  EXPECT_FALSE(config_pP.partial_list == nullptr);
  // Check if consecutive tacs partial list is created
  ASSERT_EQ(
      config_pP.partial_list->list_type,
      TRACKING_AREA_IDENTITY_LIST_TYPE_ONE_PLMN_CONSECUTIVE_TACS);
  ASSERT_EQ(config_pP.num_par_lists, 1);
  ASSERT_EQ(config_pP.partial_list->nb_elem, config_pP.served_tai.nb_tai);

  EXPECT_FALSE(config_pP.partial_list[itr].plmn == nullptr);
  EXPECT_FALSE(config_pP.partial_list[itr].tac == nullptr);
  ASSERT_EQ(config_pP.partial_list->plmn[itr].mcc_digit1, 0);
  ASSERT_EQ(config_pP.partial_list->plmn[itr].mcc_digit2, 0);
  ASSERT_EQ(config_pP.partial_list->plmn[itr].mcc_digit3, 1);
  ASSERT_EQ(config_pP.partial_list->plmn[itr].mnc_digit1, 0);
  ASSERT_EQ(config_pP.partial_list->plmn[itr].mnc_digit2, 1);
  ASSERT_EQ(config_pP.partial_list->plmn[itr].mnc_digit3, 15);
  ASSERT_EQ(config_pP.partial_list->tac[itr], config_pP.served_tai.tac[itr]);

  free(config_pP.served_tai.plmn_mcc);
  free(config_pP.served_tai.plmn_mnc);
  free(config_pP.served_tai.plmn_mnc_len);
  free(config_pP.served_tai.tac);
  free(config_pP.partial_list[itr].plmn);
  free(config_pP.partial_list[itr].tac);
  free(config_pP.partial_list);
}

// Test 1 partial list with Consecutive Tacs
TEST_F(MMEConfigTest, TestParTaiListWithConsecutiveTacs) {
  mme_config_t config_pP        = {0};
  config_pP.served_tai.nb_tai   = 16;
  uint8_t itr                   = 0;
  uint16_t tac                  = 0;
  config_pP.served_tai.plmn_mcc = reinterpret_cast<uint16_t*>(
      calloc(config_pP.served_tai.nb_tai, sizeof(uint16_t)));
  config_pP.served_tai.plmn_mnc = reinterpret_cast<uint16_t*>(
      calloc(config_pP.served_tai.nb_tai, sizeof(uint16_t)));
  config_pP.served_tai.plmn_mnc_len = reinterpret_cast<uint16_t*>(
      calloc(config_pP.served_tai.nb_tai, sizeof(uint16_t)));
  config_pP.served_tai.tac = reinterpret_cast<uint16_t*>(
      calloc(config_pP.served_tai.nb_tai, sizeof(uint16_t)));
  for (itr = 0; itr < config_pP.served_tai.nb_tai; itr++) {
    config_pP.served_tai.plmn_mcc[itr]     = 1;
    config_pP.served_tai.plmn_mnc[itr]     = 1;
    config_pP.served_tai.plmn_mnc_len[itr] = 2;
  }
  // Sorted consecutive TACs
  for (itr = 0; itr < config_pP.served_tai.nb_tai; itr++) {
    config_pP.served_tai.tac[itr] = itr + 1;
  }
  // Check if consecutive tacs partial list is created
  create_partial_lists(&config_pP);
  EXPECT_FALSE(config_pP.partial_list == nullptr);
  ASSERT_EQ(
      config_pP.partial_list->list_type,
      TRACKING_AREA_IDENTITY_LIST_TYPE_ONE_PLMN_CONSECUTIVE_TACS);
  ASSERT_EQ(config_pP.num_par_lists, 1);
  ASSERT_EQ(config_pP.partial_list->nb_elem, config_pP.served_tai.nb_tai);

  itr = 0;
  EXPECT_FALSE(config_pP.partial_list[itr].plmn == nullptr);
  EXPECT_FALSE(config_pP.partial_list[itr].tac == nullptr);
  ASSERT_EQ(config_pP.partial_list->plmn[itr].mcc_digit1, 0);
  ASSERT_EQ(config_pP.partial_list->plmn[itr].mcc_digit2, 0);
  ASSERT_EQ(config_pP.partial_list->plmn[itr].mcc_digit3, 1);
  ASSERT_EQ(config_pP.partial_list->plmn[itr].mnc_digit1, 0);
  ASSERT_EQ(config_pP.partial_list->plmn[itr].mnc_digit2, 1);
  ASSERT_EQ(config_pP.partial_list->plmn[itr].mnc_digit3, 15);

  for (itr = 0; itr < config_pP.partial_list->nb_elem; itr++) {
    ASSERT_EQ(config_pP.partial_list->tac[itr], config_pP.served_tai.tac[itr]);
  }
  free(config_pP.served_tai.plmn_mcc);
  free(config_pP.served_tai.plmn_mnc);
  free(config_pP.served_tai.plmn_mnc_len);
  free(config_pP.served_tai.tac);
  free(config_pP.partial_list[0].plmn);
  free(config_pP.partial_list[0].tac);
  free(config_pP.partial_list);
}

// Test 2 partial lists with Consecutive Tacs
TEST_F(MMEConfigTest, TestTwoParTaiListsWithConsecutiveTacs) {
  mme_config_t config_pP        = {0};
  config_pP.served_tai.nb_tai   = 20;
  uint8_t itr                   = 0;
  config_pP.served_tai.plmn_mcc = reinterpret_cast<uint16_t*>(
      calloc(config_pP.served_tai.nb_tai, sizeof(uint16_t)));
  config_pP.served_tai.plmn_mnc = reinterpret_cast<uint16_t*>(
      calloc(config_pP.served_tai.nb_tai, sizeof(uint16_t)));
  config_pP.served_tai.plmn_mnc_len = reinterpret_cast<uint16_t*>(
      calloc(config_pP.served_tai.nb_tai, sizeof(uint16_t)));
  config_pP.served_tai.tac = reinterpret_cast<uint16_t*>(
      calloc(config_pP.served_tai.nb_tai, sizeof(uint16_t)));
  for (itr = 0; itr < config_pP.served_tai.nb_tai; itr++) {
    config_pP.served_tai.plmn_mcc[itr]     = 1;
    config_pP.served_tai.plmn_mnc[itr]     = 1;
    config_pP.served_tai.plmn_mnc_len[itr] = 2;
  }
  itr = 0;
  // Sorted consecutive TACs
  for (itr = 0; itr < config_pP.served_tai.nb_tai; itr++) {
    config_pP.served_tai.tac[itr] = itr + 1;
  }
  // Check if consecutive tacs partial list is created
  create_partial_lists(&config_pP);
  EXPECT_FALSE(config_pP.partial_list == nullptr);
  ASSERT_EQ(config_pP.num_par_lists, 2);
  ASSERT_EQ(config_pP.partial_list[0].nb_elem, 16);
  ASSERT_EQ(config_pP.partial_list[1].nb_elem, 4);

  for (itr = 0; itr < config_pP.num_par_lists; itr++) {
    ASSERT_EQ(
        config_pP.partial_list[itr].list_type,
        TRACKING_AREA_IDENTITY_LIST_TYPE_ONE_PLMN_CONSECUTIVE_TACS);
    EXPECT_FALSE(config_pP.partial_list[itr].plmn == nullptr);
    EXPECT_FALSE(config_pP.partial_list[itr].tac == nullptr);
    ASSERT_EQ(config_pP.partial_list[itr].plmn[0].mcc_digit1, 0);
    ASSERT_EQ(config_pP.partial_list[itr].plmn[0].mcc_digit2, 0);
    ASSERT_EQ(config_pP.partial_list[itr].plmn[0].mcc_digit3, 1);
    ASSERT_EQ(config_pP.partial_list[itr].plmn[0].mnc_digit1, 0);
    ASSERT_EQ(config_pP.partial_list[itr].plmn[0].mnc_digit2, 1);
    ASSERT_EQ(config_pP.partial_list[itr].plmn[0].mnc_digit3, 15);
  }
  itr = 0;
  for (uint8_t idx = 0; idx < config_pP.num_par_lists; idx++) {
    for (uint8_t idx2 = 0; idx2 < config_pP.partial_list[idx].nb_elem; idx2++) {
      ASSERT_EQ(
          config_pP.partial_list[idx].tac[idx2],
          config_pP.served_tai.tac[itr++]);
    }
  }
  free(config_pP.served_tai.plmn_mcc);
  free(config_pP.served_tai.plmn_mnc);
  free(config_pP.served_tai.plmn_mnc_len);
  free(config_pP.served_tai.tac);
  for (itr = 0; itr < config_pP.num_par_lists; itr++) {
    free(config_pP.partial_list[itr].plmn);
    free(config_pP.partial_list[itr].tac);
  }
  free(config_pP.partial_list);
}

// Test 1 partial list with Non-consecutive Tacs
TEST_F(MMEConfigTest, TestParTaiListWithNonConsecutiveTacs) {
  mme_config_t config_pP        = {0};
  config_pP.served_tai.nb_tai   = 16;
  uint8_t itr                   = 0;
  uint16_t tac                  = 0;
  config_pP.served_tai.plmn_mcc = reinterpret_cast<uint16_t*>(
      calloc(config_pP.served_tai.nb_tai, sizeof(uint16_t)));
  config_pP.served_tai.plmn_mnc = reinterpret_cast<uint16_t*>(
      calloc(config_pP.served_tai.nb_tai, sizeof(uint16_t)));
  config_pP.served_tai.plmn_mnc_len = reinterpret_cast<uint16_t*>(
      calloc(config_pP.served_tai.nb_tai, sizeof(uint16_t)));
  config_pP.served_tai.tac = reinterpret_cast<uint16_t*>(
      calloc(config_pP.served_tai.nb_tai, sizeof(uint16_t)));
  for (itr = 0; itr < config_pP.served_tai.nb_tai; itr++) {
    config_pP.served_tai.plmn_mcc[itr]     = 1;
    config_pP.served_tai.plmn_mnc[itr]     = 1;
    config_pP.served_tai.plmn_mnc_len[itr] = 2;
  }
  // Sorted non-consecutive TACs
  for (itr = 0; itr < config_pP.served_tai.nb_tai; itr++) {
    config_pP.served_tai.tac[itr] = ncon_tac[itr];
  }
  // Check if non-consecutive tacs partial list is created
  create_partial_lists(&config_pP);
  EXPECT_FALSE(config_pP.partial_list == nullptr);
  ASSERT_EQ(
      config_pP.partial_list->list_type,
      TRACKING_AREA_IDENTITY_LIST_TYPE_ONE_PLMN_NON_CONSECUTIVE_TACS);
  ASSERT_EQ(config_pP.num_par_lists, 1);
  ASSERT_EQ(config_pP.partial_list->nb_elem, config_pP.served_tai.nb_tai);

  itr = 0;
  EXPECT_FALSE(config_pP.partial_list[itr].plmn == nullptr);
  EXPECT_FALSE(config_pP.partial_list[itr].tac == nullptr);
  ASSERT_EQ(config_pP.partial_list->plmn[itr].mcc_digit1, 0);
  ASSERT_EQ(config_pP.partial_list->plmn[itr].mcc_digit2, 0);
  ASSERT_EQ(config_pP.partial_list->plmn[itr].mcc_digit3, 1);
  ASSERT_EQ(config_pP.partial_list->plmn[itr].mnc_digit1, 0);
  ASSERT_EQ(config_pP.partial_list->plmn[itr].mnc_digit2, 1);
  ASSERT_EQ(config_pP.partial_list->plmn[itr].mnc_digit3, 15);

  for (itr = 0; itr < config_pP.partial_list->nb_elem; itr++) {
    ASSERT_EQ(config_pP.partial_list->tac[itr], config_pP.served_tai.tac[itr]);
  }
  free(config_pP.served_tai.plmn_mcc);
  free(config_pP.served_tai.plmn_mnc);
  free(config_pP.served_tai.plmn_mnc_len);
  free(config_pP.served_tai.tac);
  free(config_pP.partial_list[0].plmn);
  free(config_pP.partial_list[0].tac);
  free(config_pP.partial_list);
}

// Test 2 partial lists with Non-consecutive Tacs
TEST_F(MMEConfigTest, TestTwoParTaiListsWithNonConsecutiveTacs) {
  mme_config_t config_pP        = {0};
  config_pP.served_tai.nb_tai   = 20;
  uint8_t itr                   = 0;
  config_pP.served_tai.plmn_mcc = reinterpret_cast<uint16_t*>(
      calloc(config_pP.served_tai.nb_tai, sizeof(uint16_t)));
  config_pP.served_tai.plmn_mnc = reinterpret_cast<uint16_t*>(
      calloc(config_pP.served_tai.nb_tai, sizeof(uint16_t)));
  config_pP.served_tai.plmn_mnc_len = reinterpret_cast<uint16_t*>(
      calloc(config_pP.served_tai.nb_tai, sizeof(uint16_t)));
  config_pP.served_tai.tac = reinterpret_cast<uint16_t*>(
      calloc(config_pP.served_tai.nb_tai, sizeof(uint16_t)));
  for (itr = 0; itr < config_pP.served_tai.nb_tai; itr++) {
    config_pP.served_tai.plmn_mcc[itr]     = 1;
    config_pP.served_tai.plmn_mnc[itr]     = 1;
    config_pP.served_tai.plmn_mnc_len[itr] = 2;
  }
  itr = 0;
  // Sorted non-consecutive TACs
  for (itr = 0; itr < config_pP.served_tai.nb_tai; itr++) {
    config_pP.served_tai.tac[itr] = ncon_tac[itr];
  }
  // Check if 2 non-consecutive tacs partial lists are created
  create_partial_lists(&config_pP);
  EXPECT_FALSE(config_pP.partial_list == nullptr);
  ASSERT_EQ(config_pP.num_par_lists, 2);
  ASSERT_EQ(config_pP.partial_list[0].nb_elem, 16);
  ASSERT_EQ(config_pP.partial_list[1].nb_elem, 4);

  for (itr = 0; itr < config_pP.num_par_lists; itr++) {
    ASSERT_EQ(
        config_pP.partial_list[itr].list_type,
        TRACKING_AREA_IDENTITY_LIST_TYPE_ONE_PLMN_NON_CONSECUTIVE_TACS);
    EXPECT_FALSE(config_pP.partial_list[itr].plmn == nullptr);
    EXPECT_FALSE(config_pP.partial_list[itr].tac == nullptr);
    ASSERT_EQ(config_pP.partial_list[itr].plmn[0].mcc_digit1, 0);
    ASSERT_EQ(config_pP.partial_list[itr].plmn[0].mcc_digit2, 0);
    ASSERT_EQ(config_pP.partial_list[itr].plmn[0].mcc_digit3, 1);
    ASSERT_EQ(config_pP.partial_list[itr].plmn[0].mnc_digit1, 0);
    ASSERT_EQ(config_pP.partial_list[itr].plmn[0].mnc_digit2, 1);
    ASSERT_EQ(config_pP.partial_list[itr].plmn[0].mnc_digit3, 15);
  }
  itr = 0;
  for (uint8_t idx = 0; idx < config_pP.num_par_lists; idx++) {
    for (uint8_t idx2 = 0; idx2 < config_pP.partial_list[idx].nb_elem; idx2++) {
      ASSERT_EQ(
          config_pP.partial_list[idx].tac[idx2],
          config_pP.served_tai.tac[itr++]);
    }
  }

  free(config_pP.served_tai.plmn_mcc);
  free(config_pP.served_tai.plmn_mnc);
  free(config_pP.served_tai.plmn_mnc_len);
  free(config_pP.served_tai.tac);
  for (itr = 0; itr < config_pP.num_par_lists; itr++) {
    free(config_pP.partial_list[itr].plmn);
    free(config_pP.partial_list[itr].tac);
  }
  free(config_pP.partial_list);
}

// Test 2 partial lists with 1-Consecutive tacs and 1-Non-consecutive Tacs
TEST_F(MMEConfigTest, TestTwoParTaiListsWithConsAndNonConsecutiveTacs) {
  mme_config_t config_pP        = {0};
  config_pP.served_tai.nb_tai   = 24;
  uint8_t itr                   = 0;
  config_pP.served_tai.plmn_mcc = reinterpret_cast<uint16_t*>(
      calloc(config_pP.served_tai.nb_tai, sizeof(uint16_t)));
  config_pP.served_tai.plmn_mnc = reinterpret_cast<uint16_t*>(
      calloc(config_pP.served_tai.nb_tai, sizeof(uint16_t)));
  config_pP.served_tai.plmn_mnc_len = reinterpret_cast<uint16_t*>(
      calloc(config_pP.served_tai.nb_tai, sizeof(uint16_t)));
  config_pP.served_tai.tac = reinterpret_cast<uint16_t*>(
      calloc(config_pP.served_tai.nb_tai, sizeof(uint16_t)));
  for (itr = 0; itr < config_pP.served_tai.nb_tai; itr++) {
    config_pP.served_tai.plmn_mcc[itr]     = 1;
    config_pP.served_tai.plmn_mnc[itr]     = 1;
    config_pP.served_tai.plmn_mnc_len[itr] = 2;
  }
  // Sorted consecutive TACs
  for (itr = 0; itr < 16; itr++) {
    config_pP.served_tai.tac[itr] = itr + 1;
  }
  // Sorted non-consecutive TACs
  config_pP.served_tai.tac[itr++] = 19;
  config_pP.served_tai.tac[itr++] = 21;
  config_pP.served_tai.tac[itr++] = 23;
  config_pP.served_tai.tac[itr++] = 26;
  config_pP.served_tai.tac[itr++] = 28;
  config_pP.served_tai.tac[itr++] = 31;
  config_pP.served_tai.tac[itr++] = 33;
  config_pP.served_tai.tac[itr++] = 35;

  /* Check if 1 consecutive tacs partial list and 1 non-consecutive
   * tacs partial lists are created
   */
  create_partial_lists(&config_pP);
  EXPECT_FALSE(config_pP.partial_list == nullptr);
  ASSERT_EQ(config_pP.num_par_lists, 2);
  ASSERT_EQ(config_pP.partial_list[0].nb_elem, 16);
  ASSERT_EQ(config_pP.partial_list[1].nb_elem, 8);
  ASSERT_EQ(
      config_pP.partial_list[0].list_type,
      TRACKING_AREA_IDENTITY_LIST_TYPE_ONE_PLMN_CONSECUTIVE_TACS);
  ASSERT_EQ(
      config_pP.partial_list[1].list_type,
      TRACKING_AREA_IDENTITY_LIST_TYPE_ONE_PLMN_NON_CONSECUTIVE_TACS);

  for (itr = 0; itr < config_pP.num_par_lists; itr++) {
    EXPECT_FALSE(config_pP.partial_list[itr].plmn == nullptr);
    EXPECT_FALSE(config_pP.partial_list[itr].tac == nullptr);
    ASSERT_EQ(config_pP.partial_list[itr].plmn[0].mcc_digit1, 0);
    ASSERT_EQ(config_pP.partial_list[itr].plmn[0].mcc_digit2, 0);
    ASSERT_EQ(config_pP.partial_list[itr].plmn[0].mcc_digit3, 1);
    ASSERT_EQ(config_pP.partial_list[itr].plmn[0].mnc_digit1, 0);
    ASSERT_EQ(config_pP.partial_list[itr].plmn[0].mnc_digit2, 1);
    ASSERT_EQ(config_pP.partial_list[itr].plmn[0].mnc_digit3, 15);
  }
  uint8_t idx2 = 0;
  for (itr = 0; itr < config_pP.num_par_lists; itr++) {
    for (uint8_t idx = 0; idx < config_pP.partial_list[itr].nb_elem; idx++) {
      ASSERT_EQ(
          config_pP.partial_list[itr].tac[idx],
          config_pP.served_tai.tac[idx2++]);
    }
  }
  free(config_pP.served_tai.plmn_mcc);
  free(config_pP.served_tai.plmn_mnc);
  free(config_pP.served_tai.plmn_mnc_len);
  free(config_pP.served_tai.tac);
  for (itr = 0; itr < config_pP.num_par_lists; itr++) {
    free(config_pP.partial_list[itr].plmn);
    free(config_pP.partial_list[itr].tac);
  }
  free(config_pP.partial_list);
}

// Test 1 partial list with many plmns
TEST_F(MMEConfigTest, TestParTaiListWithManyPlmns) {
  mme_config_t config_pP        = {0};
  config_pP.served_tai.nb_tai   = 6;
  uint8_t itr                   = 0;
  uint16_t tac                  = 0;
  config_pP.served_tai.plmn_mcc = reinterpret_cast<uint16_t*>(
      calloc(config_pP.served_tai.nb_tai, sizeof(uint16_t)));
  config_pP.served_tai.plmn_mnc = reinterpret_cast<uint16_t*>(
      calloc(config_pP.served_tai.nb_tai, sizeof(uint16_t)));
  config_pP.served_tai.plmn_mnc_len = reinterpret_cast<uint16_t*>(
      calloc(config_pP.served_tai.nb_tai, sizeof(uint16_t)));
  config_pP.served_tai.tac = reinterpret_cast<uint16_t*>(
      calloc(config_pP.served_tai.nb_tai, sizeof(uint16_t)));
  for (itr = 0; itr < config_pP.served_tai.nb_tai; itr++) {
    config_pP.served_tai.plmn_mcc[itr]     = 1;
    config_pP.served_tai.plmn_mnc[itr]     = itr + 1;
    config_pP.served_tai.plmn_mnc_len[itr] = 2;
    config_pP.served_tai.tac[itr]          = itr + 1;
  }
  create_partial_lists(&config_pP);
  EXPECT_FALSE(config_pP.partial_list == nullptr);
  // Check if partial list with many plmns is created
  ASSERT_EQ(
      config_pP.partial_list->list_type,
      TRACKING_AREA_IDENTITY_LIST_TYPE_MANY_PLMNS);
  ASSERT_EQ(config_pP.num_par_lists, 1);
  ASSERT_EQ(config_pP.partial_list->nb_elem, config_pP.served_tai.nb_tai);

  itr = 0;
  EXPECT_FALSE(config_pP.partial_list[itr].plmn == nullptr);
  EXPECT_FALSE(config_pP.partial_list[itr].tac == nullptr);
  for (itr = 0; itr < config_pP.partial_list->nb_elem; itr++) {
    ASSERT_EQ(config_pP.partial_list->plmn[itr].mcc_digit1, 0);
    ASSERT_EQ(config_pP.partial_list->plmn[itr].mcc_digit2, 0);
    ASSERT_EQ(config_pP.partial_list->plmn[itr].mcc_digit3, 1);
    ASSERT_EQ(config_pP.partial_list->plmn[itr].mnc_digit1, 0);
    ASSERT_EQ(config_pP.partial_list->plmn[itr].mnc_digit2, itr + 1);
    ASSERT_EQ(config_pP.partial_list->plmn[itr].mnc_digit3, 15);
    ASSERT_EQ(config_pP.partial_list->tac[itr], config_pP.served_tai.tac[itr]);
  }
  free(config_pP.served_tai.plmn_mcc);
  free(config_pP.served_tai.plmn_mnc);
  free(config_pP.served_tai.plmn_mnc_len);
  free(config_pP.served_tai.tac);
  free(config_pP.partial_list[0].plmn);
  free(config_pP.partial_list[0].tac);
  free(config_pP.partial_list);
}

// Test 3 partial lists, 1-consecutive tacs, 1-non consecutive tacs,1-many plmns
TEST_F(MMEConfigTest, TestMixedParTaiLists) {
  mme_config_t config_pP        = {0};
  config_pP.served_tai.nb_tai   = 35;
  uint8_t itr                   = 0;
  config_pP.served_tai.plmn_mcc = reinterpret_cast<uint16_t*>(
      calloc(config_pP.served_tai.nb_tai, sizeof(uint16_t)));
  config_pP.served_tai.plmn_mnc = reinterpret_cast<uint16_t*>(
      calloc(config_pP.served_tai.nb_tai, sizeof(uint16_t)));
  config_pP.served_tai.plmn_mnc_len = reinterpret_cast<uint16_t*>(
      calloc(config_pP.served_tai.nb_tai, sizeof(uint16_t)));
  config_pP.served_tai.tac = reinterpret_cast<uint16_t*>(
      calloc(config_pP.served_tai.nb_tai, sizeof(uint16_t)));
  // Fill the same PLMN for consecutive and non-consecutive tacs (16+16)
  for (itr = 0; itr < 32; itr++) {
    config_pP.served_tai.plmn_mcc[itr]     = 1;
    config_pP.served_tai.plmn_mnc[itr]     = 1;
    config_pP.served_tai.plmn_mnc_len[itr] = 2;
  }
  // First 16 sorted consecutive TACs
  for (itr = 0; itr < 16; itr++) {
    config_pP.served_tai.tac[itr] = itr + 1;
  }
  // Next 16 sorted non-consecutive TACs
  for (uint8_t idx = 0; itr < 32; itr++, idx++) {
    config_pP.served_tai.tac[itr] = ncon_tac[idx];
  }
  // Next 3 many plmns with tacs
  for (uint8_t idx = 0; itr < 35; itr++, idx++) {
    config_pP.served_tai.plmn_mcc[itr]     = 1;
    config_pP.served_tai.plmn_mnc[itr]     = idx + 1;
    config_pP.served_tai.plmn_mnc_len[itr] = 2;
    config_pP.served_tai.tac[itr]          = idx + 1;
  }

  // Check if 3 partial lists are created
  create_partial_lists(&config_pP);
  EXPECT_FALSE(config_pP.partial_list == nullptr);
  ASSERT_EQ(config_pP.num_par_lists, 3);
  ASSERT_EQ(config_pP.partial_list[0].nb_elem, 16);
  ASSERT_EQ(config_pP.partial_list[1].nb_elem, 16);
  ASSERT_EQ(config_pP.partial_list[2].nb_elem, 3);

  ASSERT_EQ(
      config_pP.partial_list[0].list_type,
      TRACKING_AREA_IDENTITY_LIST_TYPE_ONE_PLMN_CONSECUTIVE_TACS);
  ASSERT_EQ(
      config_pP.partial_list[1].list_type,
      TRACKING_AREA_IDENTITY_LIST_TYPE_ONE_PLMN_NON_CONSECUTIVE_TACS);
  ASSERT_EQ(
      config_pP.partial_list[2].list_type,
      TRACKING_AREA_IDENTITY_LIST_TYPE_MANY_PLMNS);
  // Verify plmn for consecutive and non-consecutive tacs
  for (itr = 0; itr < config_pP.num_par_lists - 1; itr++) {
    EXPECT_FALSE(config_pP.partial_list[itr].plmn == nullptr);
    ASSERT_EQ(config_pP.partial_list[itr].plmn[0].mcc_digit1, 0);
    ASSERT_EQ(config_pP.partial_list[itr].plmn[0].mcc_digit2, 0);
    ASSERT_EQ(config_pP.partial_list[itr].plmn[0].mcc_digit3, 1);
    ASSERT_EQ(config_pP.partial_list[itr].plmn[0].mnc_digit1, 0);
    ASSERT_EQ(config_pP.partial_list[itr].plmn[0].mnc_digit2, 1);
    ASSERT_EQ(config_pP.partial_list[itr].plmn[0].mnc_digit3, 15);
  }

  // Verify plmn for many plmns
  for (uint8_t idx = 0; idx < config_pP.partial_list[3].nb_elem; idx++) {
    EXPECT_FALSE(config_pP.partial_list[3].plmn == nullptr);
    ASSERT_EQ(config_pP.partial_list[3].plmn[idx].mcc_digit1, 0);
    ASSERT_EQ(config_pP.partial_list[3].plmn[idx].mcc_digit2, 0);
    ASSERT_EQ(config_pP.partial_list[3].plmn[idx].mcc_digit3, 1);
    ASSERT_EQ(config_pP.partial_list[3].plmn[idx].mnc_digit1, 0);
    ASSERT_EQ(config_pP.partial_list[3].plmn[idx].mnc_digit2, idx + 1);
    ASSERT_EQ(config_pP.partial_list[3].plmn[idx].mnc_digit3, 15);
  }
  // Verify TACs
  uint8_t idx2 = 0;
  for (itr = 0; itr < config_pP.num_par_lists; itr++) {
    EXPECT_FALSE(config_pP.partial_list[itr].tac == nullptr);
    for (uint8_t idx = 0; idx < config_pP.partial_list[itr].nb_elem; idx++) {
      ASSERT_EQ(
          config_pP.partial_list[itr].tac[idx],
          config_pP.served_tai.tac[idx2++]);
    }
  }
  free(config_pP.served_tai.plmn_mcc);
  free(config_pP.served_tai.plmn_mnc);
  free(config_pP.served_tai.plmn_mnc_len);
  free(config_pP.served_tai.tac);
  for (itr = 0; itr < config_pP.num_par_lists; itr++) {
    free(config_pP.partial_list[itr].plmn);
    free(config_pP.partial_list[itr].tac);
  }
  free(config_pP.partial_list);
}

}  // namespace lte
}  // namespace magma

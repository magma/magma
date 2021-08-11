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
#include "mme_config.h"
#include "dynamic_memory_check.h"
#include "log.h"
}

class MMEConfigTest : public ::testing::Test {
  virtual void SetUp() {}
  virtual void TearDown() {}
};

TEST_F(MMEConfigTest, TestOneTai) {
  mme_config_t config_pP = {0};
  config_pP.served_tai.nb_tai = 1;
  uint8_t itr=0;
  uint16_t tac = 0;
  config_pP.served_tai.plmn_mcc = (uint16_t*)calloc(config_pP.served_tai.nb_tai, sizeof(uint16_t));
  config_pP.served_tai.plmn_mnc = (uint16_t*)calloc(config_pP.served_tai.nb_tai, sizeof(uint16_t));
  config_pP.served_tai.plmn_mnc_len = (uint16_t*)calloc(config_pP.served_tai.nb_tai, sizeof(uint16_t));
  config_pP.served_tai.tac = (uint16_t*)calloc(config_pP.served_tai.nb_tai, sizeof(uint16_t));
  config_pP.served_tai.plmn_mcc[itr] = 1;
  config_pP.served_tai.plmn_mnc[itr] = 1;
  config_pP.served_tai.plmn_mnc_len[itr] = 2;
  config_pP.served_tai.tac[itr] = 1;
  /* Check if consecutive tacs partial list is created */
  create_partial_lists(&config_pP);
  EXPECT_FALSE(config_pP.partial_list == nullptr);
  ASSERT_EQ(config_pP.partial_list->list_type, TRACKING_AREA_IDENTITY_LIST_TYPE_ONE_PLMN_CONSECUTIVE_TACS);
  ASSERT_EQ(config_pP.num_par_lists, 1);
  ASSERT_EQ(config_pP.partial_list->nb_elem, config_pP.served_tai.nb_tai);

  EXPECT_FALSE(config_pP.partial_list[itr].plmn == nullptr);
  EXPECT_FALSE(config_pP.partial_list[itr].tac == nullptr);
  ASSERT_EQ(config_pP.partial_list->plmn[itr].mcc_digit1,0);
  ASSERT_EQ(config_pP.partial_list->plmn[itr].mcc_digit2,0);
  ASSERT_EQ(config_pP.partial_list->plmn[itr].mcc_digit3,1);
  ASSERT_EQ(config_pP.partial_list->plmn[itr].mnc_digit1,0);
  ASSERT_EQ(config_pP.partial_list->plmn[itr].mnc_digit2,1);
  ASSERT_EQ(config_pP.partial_list->plmn[itr].mnc_digit3,15);
  ASSERT_EQ(config_pP.partial_list->tac[itr],config_pP.served_tai.tac[itr]);

  free(config_pP.served_tai.plmn_mcc);
  free(config_pP.served_tai.plmn_mnc);
  free(config_pP.served_tai.plmn_mnc_len);
  free(config_pP.served_tai.tac);
  free(config_pP.partial_list[itr].plmn);
  free(config_pP.partial_list[itr].tac);
  free(config_pP.partial_list);
}

TEST_F(MMEConfigTest, TestParTaiListWithConsecutiveTacs) {
  mme_config_t config_pP = {0};
  config_pP.served_tai.nb_tai = 16;
  uint8_t itr=0;
  uint16_t tac = 0;
  config_pP.served_tai.plmn_mcc = (uint16_t*)calloc(config_pP.served_tai.nb_tai, sizeof(uint16_t));
  config_pP.served_tai.plmn_mnc = (uint16_t*)calloc(config_pP.served_tai.nb_tai, sizeof(uint16_t));
  config_pP.served_tai.plmn_mnc_len = (uint16_t*)calloc(config_pP.served_tai.nb_tai, sizeof(uint16_t));
  config_pP.served_tai.tac = (uint16_t*)calloc(config_pP.served_tai.nb_tai, sizeof(uint16_t));
  for (itr=0; itr<config_pP.served_tai.nb_tai; itr++) {
    config_pP.served_tai.plmn_mcc[itr] = 1;
    config_pP.served_tai.plmn_mnc[itr] = 1;
    config_pP.served_tai.plmn_mnc_len[itr] = 2;
  }
  itr = 0;
  // Sorted consecutive TACs
  for (itr=0; itr<config_pP.served_tai.nb_tai; itr++) {
    config_pP.served_tai.tac[itr] = itr+1;
  }
  /* Check if consecutive tacs partial list is created */
  create_partial_lists(&config_pP);
  EXPECT_FALSE(config_pP.partial_list == nullptr);
  ASSERT_EQ(config_pP.partial_list->list_type, TRACKING_AREA_IDENTITY_LIST_TYPE_ONE_PLMN_CONSECUTIVE_TACS);
  ASSERT_EQ(config_pP.num_par_lists, 1);
  ASSERT_EQ(config_pP.partial_list->nb_elem, config_pP.served_tai.nb_tai);

  itr=0;
  EXPECT_FALSE(config_pP.partial_list[itr].plmn == nullptr);
  EXPECT_FALSE(config_pP.partial_list[itr].tac == nullptr);
  ASSERT_EQ(config_pP.partial_list->plmn[itr].mcc_digit1,0);
  ASSERT_EQ(config_pP.partial_list->plmn[itr].mcc_digit2,0);
  ASSERT_EQ(config_pP.partial_list->plmn[itr].mcc_digit3,1);
  ASSERT_EQ(config_pP.partial_list->plmn[itr].mnc_digit1,0);
  ASSERT_EQ(config_pP.partial_list->plmn[itr].mnc_digit2,1);
  ASSERT_EQ(config_pP.partial_list->plmn[itr].mnc_digit3,15);

  for (itr=0;itr<config_pP.partial_list->nb_elem;itr++) {
    ASSERT_EQ(config_pP.partial_list->tac[itr],config_pP.served_tai.tac[itr]);
  }
  free(config_pP.served_tai.plmn_mcc);
  free(config_pP.served_tai.plmn_mnc);
  free(config_pP.served_tai.plmn_mnc_len);
  free(config_pP.served_tai.tac);
  free(config_pP.partial_list[0].plmn);
  free(config_pP.partial_list[0].tac);
  free(config_pP.partial_list);
}

TEST_F(MMEConfigTest, TestTwoParTaiListsWithConsecutiveTacs) {
  mme_config_t config_pP = {0};
  config_pP.served_tai.nb_tai = 20;
  uint8_t itr=0;
  config_pP.served_tai.plmn_mcc = (uint16_t*)calloc(config_pP.served_tai.nb_tai, sizeof(uint16_t));
  config_pP.served_tai.plmn_mnc = (uint16_t*)calloc(config_pP.served_tai.nb_tai, sizeof(uint16_t));
  config_pP.served_tai.plmn_mnc_len = (uint16_t*)calloc(config_pP.served_tai.nb_tai, sizeof(uint16_t));
  config_pP.served_tai.tac = (uint16_t*)calloc(config_pP.served_tai.nb_tai, sizeof(uint16_t));
  for (itr=0; itr<config_pP.served_tai.nb_tai; itr++) {
    config_pP.served_tai.plmn_mcc[itr] = 1;
    config_pP.served_tai.plmn_mnc[itr] = 1;
    config_pP.served_tai.plmn_mnc_len[itr] = 2;
  }
  itr = 0;
  // Sorted consecutive TACs
  for (itr=0; itr<config_pP.served_tai.nb_tai; itr++) {
    config_pP.served_tai.tac[itr] = itr+1;
  }
  /* Check if consecutive tacs partial list is created */
  create_partial_lists(&config_pP);
  EXPECT_FALSE(config_pP.partial_list == nullptr);
  ASSERT_EQ(config_pP.num_par_lists, 2);
  ASSERT_EQ(config_pP.partial_list[0].nb_elem, 16);
  ASSERT_EQ(config_pP.partial_list[1].nb_elem, 4);

  for (itr=0;itr<config_pP.num_par_lists;itr++) {
    ASSERT_EQ(config_pP.partial_list[itr].list_type, TRACKING_AREA_IDENTITY_LIST_TYPE_ONE_PLMN_CONSECUTIVE_TACS);
    EXPECT_FALSE(config_pP.partial_list[itr].plmn == nullptr);
    EXPECT_FALSE(config_pP.partial_list[itr].tac == nullptr);
    ASSERT_EQ(config_pP.partial_list[itr].plmn[0].mcc_digit1,0);
    ASSERT_EQ(config_pP.partial_list[itr].plmn[0].mcc_digit2,0);
    ASSERT_EQ(config_pP.partial_list[itr].plmn[0].mcc_digit3,1);
    ASSERT_EQ(config_pP.partial_list[itr].plmn[0].mnc_digit1,0);
    ASSERT_EQ(config_pP.partial_list[itr].plmn[0].mnc_digit2,1);
    ASSERT_EQ(config_pP.partial_list[itr].plmn[0].mnc_digit3,15);
  }
  uint8_t idx=0;
  for (idx=0;idx<config_pP.partial_list[0].nb_elem;idx++) {
    ASSERT_EQ(config_pP.partial_list[0].tac[idx],config_pP.served_tai.tac[idx]);
  }
  for (uint8_t idx2=0;idx2<config_pP.partial_list[1].nb_elem;idx2++) {
    ASSERT_EQ(config_pP.partial_list[1].tac[idx2],config_pP.served_tai.tac[idx]);
    idx++;
  }
  free(config_pP.served_tai.plmn_mcc);
  free(config_pP.served_tai.plmn_mnc);
  free(config_pP.served_tai.plmn_mnc_len);
  free(config_pP.served_tai.tac);
  for (itr=0;itr<config_pP.num_par_lists;itr++) {
    free(config_pP.partial_list[itr].plmn);
    free(config_pP.partial_list[itr].tac);
  }
  free(config_pP.partial_list);
}

TEST_F(MMEConfigTest, TestParTaiListWithNonConsecutiveTacs) {
  mme_config_t config_pP = {0};
  config_pP.served_tai.nb_tai = 16;
  uint8_t itr=0;
  uint16_t tac = 0;
  config_pP.served_tai.plmn_mcc = (uint16_t*)calloc(config_pP.served_tai.nb_tai, sizeof(uint16_t));
  config_pP.served_tai.plmn_mnc = (uint16_t*)calloc(config_pP.served_tai.nb_tai, sizeof(uint16_t));
  config_pP.served_tai.plmn_mnc_len = (uint16_t*)calloc(config_pP.served_tai.nb_tai, sizeof(uint16_t));
  config_pP.served_tai.tac = (uint16_t*)calloc(config_pP.served_tai.nb_tai, sizeof(uint16_t));
  for (itr=0; itr<config_pP.served_tai.nb_tai; itr++) {
    config_pP.served_tai.plmn_mcc[itr] = 1;
    config_pP.served_tai.plmn_mnc[itr] = 1;
    config_pP.served_tai.plmn_mnc_len[itr] = 2;
  }
  itr = 0;
  // Sorted non-consecutive TACs
  config_pP.served_tai.tac[itr++] = 1;
  config_pP.served_tai.tac[itr++] = 2;
  config_pP.served_tai.tac[itr++] = 4;
  config_pP.served_tai.tac[itr++] = 6;
  config_pP.served_tai.tac[itr++] = 7;
  config_pP.served_tai.tac[itr++] = 8;
  config_pP.served_tai.tac[itr++] = 10;
  config_pP.served_tai.tac[itr++] = 11;
  config_pP.served_tai.tac[itr++] = 12;
  config_pP.served_tai.tac[itr++] = 14;
  config_pP.served_tai.tac[itr++] = 15;
  config_pP.served_tai.tac[itr++] = 16;
  config_pP.served_tai.tac[itr++] = 19;
  config_pP.served_tai.tac[itr++] = 21;
  config_pP.served_tai.tac[itr++] = 23;
  config_pP.served_tai.tac[itr++] = 26;
  /* Check if non-consecutive tacs partial list is created */
  create_partial_lists(&config_pP);
  EXPECT_FALSE(config_pP.partial_list == nullptr);
  ASSERT_EQ(config_pP.partial_list->list_type, TRACKING_AREA_IDENTITY_LIST_TYPE_ONE_PLMN_NON_CONSECUTIVE_TACS);
  ASSERT_EQ(config_pP.num_par_lists, 1);
  ASSERT_EQ(config_pP.partial_list->nb_elem, config_pP.served_tai.nb_tai);

  itr=0;
  EXPECT_FALSE(config_pP.partial_list[itr].plmn == nullptr);
  EXPECT_FALSE(config_pP.partial_list[itr].tac == nullptr);
  ASSERT_EQ(config_pP.partial_list->plmn[itr].mcc_digit1,0);
  ASSERT_EQ(config_pP.partial_list->plmn[itr].mcc_digit2,0);
  ASSERT_EQ(config_pP.partial_list->plmn[itr].mcc_digit3,1);
  ASSERT_EQ(config_pP.partial_list->plmn[itr].mnc_digit1,0);
  ASSERT_EQ(config_pP.partial_list->plmn[itr].mnc_digit2,1);
  ASSERT_EQ(config_pP.partial_list->plmn[itr].mnc_digit3,15);

  for (itr=0;itr<config_pP.partial_list->nb_elem;itr++) {
    ASSERT_EQ(config_pP.partial_list->tac[itr],config_pP.served_tai.tac[itr]);
  }
  free(config_pP.served_tai.plmn_mcc);
  free(config_pP.served_tai.plmn_mnc);
  free(config_pP.served_tai.plmn_mnc_len);
  free(config_pP.served_tai.tac);
  free(config_pP.partial_list[0].plmn);
  free(config_pP.partial_list[0].tac);
  free(config_pP.partial_list);
}

TEST_F(MMEConfigTest, TestTwoParTaiListsWithNonConsecutiveTacs) {
  mme_config_t config_pP = {0};
  config_pP.served_tai.nb_tai = 20;
  uint8_t itr=0;
  config_pP.served_tai.plmn_mcc = (uint16_t*)calloc(config_pP.served_tai.nb_tai, sizeof(uint16_t));
  config_pP.served_tai.plmn_mnc = (uint16_t*)calloc(config_pP.served_tai.nb_tai, sizeof(uint16_t));
  config_pP.served_tai.plmn_mnc_len = (uint16_t*)calloc(config_pP.served_tai.nb_tai, sizeof(uint16_t));
  config_pP.served_tai.tac = (uint16_t*)calloc(config_pP.served_tai.nb_tai, sizeof(uint16_t));
  for (itr=0; itr<config_pP.served_tai.nb_tai; itr++) {
    config_pP.served_tai.plmn_mcc[itr] = 1;
    config_pP.served_tai.plmn_mnc[itr] = 1;
    config_pP.served_tai.plmn_mnc_len[itr] = 2;
  }
  itr = 0;
  // Sorted non-consecutive TACs
  config_pP.served_tai.tac[itr++] = 1;
  config_pP.served_tai.tac[itr++] = 2;
  config_pP.served_tai.tac[itr++] = 4;
  config_pP.served_tai.tac[itr++] = 6;
  config_pP.served_tai.tac[itr++] = 7;
  config_pP.served_tai.tac[itr++] = 8;
  config_pP.served_tai.tac[itr++] = 10;
  config_pP.served_tai.tac[itr++] = 11;
  config_pP.served_tai.tac[itr++] = 12;
  config_pP.served_tai.tac[itr++] = 14;
  config_pP.served_tai.tac[itr++] = 15;
  config_pP.served_tai.tac[itr++] = 16;
  config_pP.served_tai.tac[itr++] = 19;
  config_pP.served_tai.tac[itr++] = 21;
  config_pP.served_tai.tac[itr++] = 23;
  config_pP.served_tai.tac[itr++] = 26;
  config_pP.served_tai.tac[itr++] = 28;
  config_pP.served_tai.tac[itr++] = 31;
  config_pP.served_tai.tac[itr++] = 33;
  config_pP.served_tai.tac[itr++] = 35;

  /* Check if 2 non-consecutive tacs partial lists are created */
  create_partial_lists(&config_pP);
  EXPECT_FALSE(config_pP.partial_list == nullptr);
  ASSERT_EQ(config_pP.num_par_lists, 2);
  ASSERT_EQ(config_pP.partial_list[0].nb_elem, 16);
  ASSERT_EQ(config_pP.partial_list[1].nb_elem, 4);

  for (itr=0;itr<config_pP.num_par_lists;itr++) {
    ASSERT_EQ(config_pP.partial_list[itr].list_type, TRACKING_AREA_IDENTITY_LIST_TYPE_ONE_PLMN_NON_CONSECUTIVE_TACS);
    EXPECT_FALSE(config_pP.partial_list[itr].plmn == nullptr);
    EXPECT_FALSE(config_pP.partial_list[itr].tac == nullptr);
    ASSERT_EQ(config_pP.partial_list[itr].plmn[0].mcc_digit1,0);
    ASSERT_EQ(config_pP.partial_list[itr].plmn[0].mcc_digit2,0);
    ASSERT_EQ(config_pP.partial_list[itr].plmn[0].mcc_digit3,1);
    ASSERT_EQ(config_pP.partial_list[itr].plmn[0].mnc_digit1,0);
    ASSERT_EQ(config_pP.partial_list[itr].plmn[0].mnc_digit2,1);
    ASSERT_EQ(config_pP.partial_list[itr].plmn[0].mnc_digit3,15);
  }
  uint8_t idx=0;
  for (idx=0;idx<config_pP.partial_list[0].nb_elem;idx++) {
    ASSERT_EQ(config_pP.partial_list[0].tac[idx],config_pP.served_tai.tac[idx]);
  }
  for (uint8_t idx2=0;idx2<config_pP.partial_list[1].nb_elem;idx2++) {
    ASSERT_EQ(config_pP.partial_list[1].tac[idx2],config_pP.served_tai.tac[idx]);
    idx++;
  }
  free(config_pP.served_tai.plmn_mcc);
  free(config_pP.served_tai.plmn_mnc);
  free(config_pP.served_tai.plmn_mnc_len);
  free(config_pP.served_tai.tac);
  for (itr=0;itr<config_pP.num_par_lists;itr++) {
    free(config_pP.partial_list[itr].plmn);
    free(config_pP.partial_list[itr].tac);
  }
  free(config_pP.partial_list);
}

TEST_F(MMEConfigTest, TestTwoParTaiListsWithConsAndNonConsecutiveTacs) {
  mme_config_t config_pP = {0};
  config_pP.served_tai.nb_tai = 24;
  uint8_t itr=0;
  config_pP.served_tai.plmn_mcc = (uint16_t*)calloc(config_pP.served_tai.nb_tai, sizeof(uint16_t));
  config_pP.served_tai.plmn_mnc = (uint16_t*)calloc(config_pP.served_tai.nb_tai, sizeof(uint16_t));
  config_pP.served_tai.plmn_mnc_len = (uint16_t*)calloc(config_pP.served_tai.nb_tai, sizeof(uint16_t));
  config_pP.served_tai.tac = (uint16_t*)calloc(config_pP.served_tai.nb_tai, sizeof(uint16_t));
  for (itr=0; itr<config_pP.served_tai.nb_tai; itr++) {
    config_pP.served_tai.plmn_mcc[itr] = 1;
    config_pP.served_tai.plmn_mnc[itr] = 1;
    config_pP.served_tai.plmn_mnc_len[itr] = 2;
  }
  itr = 0;
  // Sorted consecutive TACs
  for (itr=0; itr<16; itr++) {
    config_pP.served_tai.tac[itr] = itr+1;
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
  ASSERT_EQ(config_pP.partial_list[0].list_type, TRACKING_AREA_IDENTITY_LIST_TYPE_ONE_PLMN_CONSECUTIVE_TACS);
  ASSERT_EQ(config_pP.partial_list[1].list_type, TRACKING_AREA_IDENTITY_LIST_TYPE_ONE_PLMN_NON_CONSECUTIVE_TACS);

  for (itr=0;itr<config_pP.num_par_lists;itr++) {
    EXPECT_FALSE(config_pP.partial_list[itr].plmn == nullptr);
    EXPECT_FALSE(config_pP.partial_list[itr].tac == nullptr);
    ASSERT_EQ(config_pP.partial_list[itr].plmn[0].mcc_digit1,0);
    ASSERT_EQ(config_pP.partial_list[itr].plmn[0].mcc_digit2,0);
    ASSERT_EQ(config_pP.partial_list[itr].plmn[0].mcc_digit3,1);
    ASSERT_EQ(config_pP.partial_list[itr].plmn[0].mnc_digit1,0);
    ASSERT_EQ(config_pP.partial_list[itr].plmn[0].mnc_digit2,1);
    ASSERT_EQ(config_pP.partial_list[itr].plmn[0].mnc_digit3,15);
  }
  uint8_t idx=0;
  for (idx=0;idx<config_pP.partial_list[0].nb_elem;idx++) {
    ASSERT_EQ(config_pP.partial_list[0].tac[idx],config_pP.served_tai.tac[idx]);
  }
  for (uint8_t idx2=0;idx2<config_pP.partial_list[1].nb_elem;idx2++) {
    ASSERT_EQ(config_pP.partial_list[1].tac[idx2],config_pP.served_tai.tac[idx]);
    idx++;
  }
  free(config_pP.served_tai.plmn_mcc);
  free(config_pP.served_tai.plmn_mnc);
  free(config_pP.served_tai.plmn_mnc_len);
  free(config_pP.served_tai.tac);
  for (itr=0;itr<config_pP.num_par_lists;itr++) {
    free(config_pP.partial_list[itr].plmn);
    free(config_pP.partial_list[itr].tac);
  }
  free(config_pP.partial_list);
}

int main(int argc, char** argv) {
  ::testing::InitGoogleTest(&argc, argv);
  OAILOG_INIT("MME", OAILOG_LEVEL_DEBUG, MAX_LOG_PROTOS);
  return RUN_ALL_TESTS();
}

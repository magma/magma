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
#include <stdlib.h>
#include <stdint.h>
#include <gtest/gtest.h>

extern "C" {
#include "lte/gateway/c/core/oai/common/conversions.h"
#include "lte/gateway/c/core/oai/lib/3gpp/3gpp_23.003.h"
}

#include "lte/gateway/c/core/oai/include/mme_app_ue_context.hpp"

#define TEST_CASE_COMMON_CONVERT_MAX 10

/**
 * @brief mme_app_convert_imsi_to_imsi_mme: converts the imsi_t struct to the
 * imsi mme struct
 * @param imsi_dst
 * @param imsi_src
 */
void mme_app_convert_imsi_to_imsi_mme(mme_app_imsi_t* imsi_dst,
                                      const imsi_t* imsi_src) {
  memset(imsi_dst->data, (uint8_t)'\0', sizeof(imsi_dst->data));
  IMSI_TO_STRING(imsi_src, imsi_dst->data, IMSI_BCD_DIGITS_MAX + 1);
  imsi_dst->length = strlen(imsi_dst->data);
}

TEST(imsi_empty_test, ok) {
  mme_app_imsi_t imsi_mme = {0};
  char imsi_str[] = "001012234567890";

  /* Check if imsi_mme is empty */
  ASSERT_EQ(mme_app_is_imsi_empty(&imsi_mme), true);

  /* Check if imsi_mme is not empty */
  mme_app_string_to_imsi((&imsi_mme), imsi_str);
  ASSERT_EQ(mme_app_is_imsi_empty(&imsi_mme), false);
}

TEST(imsi_convert_to_uint_test, ok) {
  mme_app_imsi_t imsi_mme_test;
  uint64_t imsi_uint64;

  char imsi_compare[TEST_CASE_COMMON_CONVERT_MAX][IMSI_BCD_DIGITS_MAX + 1] = {
      "001012234567890", "262011234567890", "262010000043448",
      "26201943210",     "262019876543210", "41004123456789",
      "310150123456789", "460001357924680", "520031234567890",
      "470010171566423"};

  uint64_t imsi_cmp[TEST_CASE_COMMON_CONVERT_MAX] = {
      1012234567890,   262011234567890, 262010000043448, 26201943210,
      262019876543210, 41004123456789,  310150123456789, 460001357924680,
      520031234567890, 470010171566423,
  };

  for (int i = 0; i < TEST_CASE_COMMON_CONVERT_MAX; i++) {
    mme_app_string_to_imsi(&imsi_mme_test, imsi_compare[i]);
    imsi_uint64 = mme_app_imsi_to_u64(imsi_mme_test);
    ASSERT_EQ(imsi_uint64, imsi_cmp[i]);
  }
}

TEST(imsi_convert_common_struct_test, ok) {
  mme_app_imsi_t imsi_mme;
  imsi_t imsi_structs[TEST_CASE_COMMON_CONVERT_MAX] = {0};
  // 001011234567890
  imsi_structs[0].u.num.digit1 = 0;
  imsi_structs[0].u.num.digit2 = 0;
  imsi_structs[0].u.num.digit3 = 1;
  imsi_structs[0].u.num.digit4 = 0;
  imsi_structs[0].u.num.digit5 = 1;
  imsi_structs[0].u.num.digit6 = 1;
  imsi_structs[0].u.num.digit7 = 2;
  imsi_structs[0].u.num.digit8 = 3;
  imsi_structs[0].u.num.digit9 = 4;
  imsi_structs[0].u.num.digit10 = 5;
  imsi_structs[0].u.num.digit11 = 6;
  imsi_structs[0].u.num.digit12 = 7;
  imsi_structs[0].u.num.digit13 = 8;
  imsi_structs[0].u.num.digit14 = 9;
  imsi_structs[0].u.num.digit15 = 0;
  // 262011234567890
  imsi_structs[1].u.num.digit1 = 2;
  imsi_structs[1].u.num.digit2 = 6;
  imsi_structs[1].u.num.digit3 = 2;
  imsi_structs[1].u.num.digit4 = 0;
  imsi_structs[1].u.num.digit5 = 1;
  imsi_structs[1].u.num.digit6 = 1;
  imsi_structs[1].u.num.digit7 = 2;
  imsi_structs[1].u.num.digit8 = 3;
  imsi_structs[1].u.num.digit9 = 4;
  imsi_structs[1].u.num.digit10 = 5;
  imsi_structs[1].u.num.digit11 = 6;
  imsi_structs[1].u.num.digit12 = 7;
  imsi_structs[1].u.num.digit13 = 8;
  imsi_structs[1].u.num.digit14 = 9;
  imsi_structs[1].u.num.digit15 = 0;
  // 310150123456789
  imsi_structs[2].u.num.digit1 = 3;
  imsi_structs[2].u.num.digit2 = 1;
  imsi_structs[2].u.num.digit3 = 0;
  imsi_structs[2].u.num.digit4 = 1;
  imsi_structs[2].u.num.digit5 = 5;
  imsi_structs[2].u.num.digit6 = 0;
  imsi_structs[2].u.num.digit7 = 1;
  imsi_structs[2].u.num.digit8 = 2;
  imsi_structs[2].u.num.digit9 = 3;
  imsi_structs[2].u.num.digit10 = 4;
  imsi_structs[2].u.num.digit11 = 5;
  imsi_structs[2].u.num.digit12 = 6;
  imsi_structs[2].u.num.digit13 = 7;
  imsi_structs[2].u.num.digit14 = 8;
  imsi_structs[2].u.num.digit15 = 9;
  // 460001357924680
  imsi_structs[3].u.num.digit1 = 4;
  imsi_structs[3].u.num.digit2 = 6;
  imsi_structs[3].u.num.digit3 = 0;
  imsi_structs[3].u.num.digit4 = 0;
  imsi_structs[3].u.num.digit5 = 0;
  imsi_structs[3].u.num.digit6 = 1;
  imsi_structs[3].u.num.digit7 = 3;
  imsi_structs[3].u.num.digit8 = 5;
  imsi_structs[3].u.num.digit9 = 7;
  imsi_structs[3].u.num.digit10 = 9;
  imsi_structs[3].u.num.digit11 = 2;
  imsi_structs[3].u.num.digit12 = 4;
  imsi_structs[3].u.num.digit13 = 6;
  imsi_structs[3].u.num.digit14 = 8;
  imsi_structs[3].u.num.digit15 = 0;
  // 520031234567890
  imsi_structs[4].u.num.digit1 = 5;
  imsi_structs[4].u.num.digit2 = 2;
  imsi_structs[4].u.num.digit3 = 0;
  imsi_structs[4].u.num.digit4 = 0;
  imsi_structs[4].u.num.digit5 = 3;
  imsi_structs[4].u.num.digit6 = 1;
  imsi_structs[4].u.num.digit7 = 2;
  imsi_structs[4].u.num.digit8 = 3;
  imsi_structs[4].u.num.digit9 = 4;
  imsi_structs[4].u.num.digit10 = 5;
  imsi_structs[4].u.num.digit11 = 6;
  imsi_structs[4].u.num.digit12 = 7;
  imsi_structs[4].u.num.digit13 = 8;
  imsi_structs[4].u.num.digit14 = 9;
  imsi_structs[4].u.num.digit15 = 0;
  // 470010171566423
  imsi_structs[5].u.num.digit1 = 4;
  imsi_structs[5].u.num.digit2 = 7;
  imsi_structs[5].u.num.digit3 = 0;
  imsi_structs[5].u.num.digit4 = 0;
  imsi_structs[5].u.num.digit5 = 1;
  imsi_structs[5].u.num.digit6 = 0;
  imsi_structs[5].u.num.digit7 = 1;
  imsi_structs[5].u.num.digit8 = 7;
  imsi_structs[5].u.num.digit9 = 1;
  imsi_structs[5].u.num.digit10 = 5;
  imsi_structs[5].u.num.digit11 = 6;
  imsi_structs[5].u.num.digit12 = 6;
  imsi_structs[5].u.num.digit13 = 4;
  imsi_structs[5].u.num.digit14 = 2;
  imsi_structs[5].u.num.digit15 = 3;
  // 41004123456789
  imsi_structs[6].u.num.digit1 = 4;
  imsi_structs[6].u.num.digit2 = 1;
  imsi_structs[6].u.num.digit3 = 0;
  imsi_structs[6].u.num.digit4 = 0;
  imsi_structs[6].u.num.digit5 = 4;
  imsi_structs[6].u.num.digit6 = 1;
  imsi_structs[6].u.num.digit7 = 2;
  imsi_structs[6].u.num.digit8 = 3;
  imsi_structs[6].u.num.digit9 = 4;
  imsi_structs[6].u.num.digit10 = 5;
  imsi_structs[6].u.num.digit11 = 6;
  imsi_structs[6].u.num.digit12 = 7;
  imsi_structs[6].u.num.digit13 = 8;
  imsi_structs[6].u.num.digit14 = 9;
  imsi_structs[6].u.num.digit15 = 0xf;
  // 4100412345678
  imsi_structs[7].u.num.digit1 = 4;
  imsi_structs[7].u.num.digit2 = 1;
  imsi_structs[7].u.num.digit3 = 0;
  imsi_structs[7].u.num.digit4 = 0;
  imsi_structs[7].u.num.digit5 = 4;
  imsi_structs[7].u.num.digit6 = 1;
  imsi_structs[7].u.num.digit7 = 2;
  imsi_structs[7].u.num.digit8 = 3;
  imsi_structs[7].u.num.digit9 = 4;
  imsi_structs[7].u.num.digit10 = 5;
  imsi_structs[7].u.num.digit11 = 6;
  imsi_structs[7].u.num.digit12 = 7;
  imsi_structs[7].u.num.digit13 = 8;
  imsi_structs[7].u.num.digit14 = 0xf;
  imsi_structs[7].u.num.digit15 = 0xf;
  // 410041234567
  imsi_structs[8].u.num.digit1 = 4;
  imsi_structs[8].u.num.digit2 = 1;
  imsi_structs[8].u.num.digit3 = 0;
  imsi_structs[8].u.num.digit4 = 0;
  imsi_structs[8].u.num.digit5 = 4;
  imsi_structs[8].u.num.digit6 = 1;
  imsi_structs[8].u.num.digit7 = 2;
  imsi_structs[8].u.num.digit8 = 3;
  imsi_structs[8].u.num.digit9 = 4;
  imsi_structs[8].u.num.digit10 = 5;
  imsi_structs[8].u.num.digit11 = 6;
  imsi_structs[8].u.num.digit12 = 7;
  imsi_structs[8].u.num.digit13 = 0xf;
  imsi_structs[8].u.num.digit14 = 0xf;
  imsi_structs[8].u.num.digit15 = 0xf;
  // 4100412
  imsi_structs[9].u.num.digit1 = 4;
  imsi_structs[9].u.num.digit2 = 1;
  imsi_structs[9].u.num.digit3 = 0;
  imsi_structs[9].u.num.digit4 = 0;
  imsi_structs[9].u.num.digit5 = 4;
  imsi_structs[9].u.num.digit6 = 1;
  imsi_structs[9].u.num.digit7 = 2;
  imsi_structs[9].u.num.digit8 = 0xf;
  imsi_structs[9].u.num.digit9 = 0xf;
  imsi_structs[9].u.num.digit10 = 0xf;
  imsi_structs[9].u.num.digit11 = 0xf;
  imsi_structs[9].u.num.digit12 = 0xf;
  imsi_structs[9].u.num.digit13 = 0xf;
  imsi_structs[9].u.num.digit14 = 0xf;
  imsi_structs[9].u.num.digit15 = 0xf;

  char* imsi_compare[TEST_CASE_COMMON_CONVERT_MAX] = {0};
  imsi_compare[0] = const_cast<char*>("001011234567890");
  imsi_compare[1] = const_cast<char*>("262011234567890");
  imsi_compare[2] = const_cast<char*>("310150123456789");
  imsi_compare[3] = const_cast<char*>("460001357924680");
  imsi_compare[4] = const_cast<char*>("520031234567890");
  imsi_compare[5] = const_cast<char*>("470010171566423");
  imsi_compare[6] = const_cast<char*>("41004123456789");
  imsi_compare[7] = const_cast<char*>("4100412345678");
  imsi_compare[8] = const_cast<char*>("410041234567");
  imsi_compare[9] = const_cast<char*>("4100412");

  for (int i = 0; i < TEST_CASE_COMMON_CONVERT_MAX; i++) {
    mme_app_convert_imsi_to_imsi_mme(&imsi_mme, &imsi_structs[i]);
    // string conversion is added to make sure that both end with a null
    // terminator character
    ASSERT_EQ(std::string(IMSI_DATA(imsi_mme)), std::string(imsi_compare[i]));
  }
}

TEST(imsi_equal_test, ok) {
  mme_app_imsi_t imsi_mme_a;
  mme_app_imsi_t imsi_mme_b;
  const char imsi_str_a[] = "001011234567890";
  const char imsi_str_b[] = "262011234567890";

  mme_app_string_to_imsi(&imsi_mme_a, imsi_str_a);
  mme_app_string_to_imsi(&imsi_mme_b, imsi_str_b);

  ASSERT_EQ(mme_app_imsi_compare(&imsi_mme_a, &imsi_mme_b), false);
  ASSERT_EQ(mme_app_imsi_compare(&imsi_mme_a, &imsi_mme_a), true);
  ASSERT_EQ(mme_app_imsi_compare(&imsi_mme_b, &imsi_mme_a), false);
  ASSERT_EQ(mme_app_imsi_compare(&imsi_mme_a, &imsi_mme_a), true);
}

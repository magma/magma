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
#include <glog/logging.h>

extern "C" {
#include "lte/gateway/c/core/common/dynamic_memory_check.h"
#include "lte/gateway/c/core/oai/common/log.h"
#include "lte/gateway/c/core/oai/include/EpsQualityOfService.h"
#include "lte/gateway/c/core/oai/lib/3gpp/3gpp_23.003.h"
}

#include "lte/gateway/c/core/oai/tasks/nas/ies/ApnAggregateMaximumBitRate.hpp"
#include "lte/gateway/c/core/oai/tasks/nas/emm/msg/AttachRequest.hpp"

class NASEncodeDecodeTest : public ::testing::Test {
  virtual void SetUp() {}
  virtual void TearDown() {}
};

TEST_F(NASEncodeDecodeTest, TestDecodeEncodeEPSQoS) {
  EpsQualityOfService eps_qos = {0};

  uint8_t eps_qos_buffersize1[] = {0x09, 0x09, 0x1c, 0x1c, 0x1c,
                                   0x1c, 0x00, 0x00, 0x00, 0x00};
  ASSERT_EQ(decode_eps_quality_of_service(&eps_qos, 0, eps_qos_buffersize1,
                                          sizeof(eps_qos_buffersize1)),
            sizeof(eps_qos_buffersize1));
  ASSERT_EQ((eps_qos.bitRates.maxBitRateForUL), 28);
  ASSERT_EQ((eps_qos.bitRates.maxBitRateForDL), 28);
  ASSERT_EQ((eps_qos.bitRatesExt.maxBitRateForUL), 0);
  ASSERT_EQ((eps_qos.bitRatesExt.maxBitRateForDL), 0);
  ASSERT_EQ((eps_qos.bitRatesExt2.maxBitRateForUL), 0);
  ASSERT_EQ((eps_qos.bitRatesExt2.maxBitRateForDL), 0);
  uint8_t eps_qos_buffersize_encoded1[10] = {0};
  ASSERT_EQ(
      encode_eps_quality_of_service(&eps_qos, 0, eps_qos_buffersize_encoded1,
                                    sizeof(eps_qos_buffersize_encoded1)),
      sizeof(eps_qos_buffersize_encoded1));

  ASSERT_EQ(memcmp(eps_qos_buffersize_encoded1, eps_qos_buffersize1,
                   sizeof(eps_qos_buffersize1)),
            0);

  uint8_t eps_qos_buffersize2[] = {0x09, 0x09, 0x3F, 0x3F, 0x3F,
                                   0x3F, 0x00, 0x00, 0x00, 0x00};

  ASSERT_EQ(decode_eps_quality_of_service(&eps_qos, 0, eps_qos_buffersize2,
                                          sizeof(eps_qos_buffersize2)),
            sizeof(eps_qos_buffersize2));
  ASSERT_EQ((eps_qos.qci), 9);
  ASSERT_EQ((eps_qos.bitRates.maxBitRateForUL), 63);
  ASSERT_EQ((eps_qos.bitRates.maxBitRateForDL), 63);
  ASSERT_EQ((eps_qos.bitRatesExt.maxBitRateForUL), 0);
  ASSERT_EQ((eps_qos.bitRatesExt.maxBitRateForDL), 0);
  ASSERT_EQ((eps_qos.bitRatesExt2.maxBitRateForUL), 0);
  ASSERT_EQ((eps_qos.bitRatesExt2.maxBitRateForDL), 0);
  uint8_t eps_qos_buffersize_encoded2[10] = {0};
  ASSERT_EQ(
      encode_eps_quality_of_service(&eps_qos, 0, eps_qos_buffersize_encoded2,
                                    sizeof(eps_qos_buffersize_encoded2)),
      sizeof(eps_qos_buffersize_encoded2));
  ASSERT_EQ(memcmp(eps_qos_buffersize_encoded2, eps_qos_buffersize2,
                   sizeof(eps_qos_buffersize2)),
            0);

  uint8_t eps_qos_buffersize3[] = {0x09, 0x09, 0x41, 0x41, 0x41,
                                   0x41, 0x00, 0x00, 0x00, 0x00};

  ASSERT_EQ(decode_eps_quality_of_service(&eps_qos, 0, eps_qos_buffersize3,
                                          sizeof(eps_qos_buffersize3)),
            sizeof(eps_qos_buffersize3));
  ASSERT_EQ((eps_qos.qci), 9);
  ASSERT_EQ((eps_qos.bitRates.maxBitRateForUL), 65);
  ASSERT_EQ((eps_qos.bitRates.maxBitRateForDL), 65);
  ASSERT_EQ((eps_qos.bitRatesExt.maxBitRateForUL), 0);
  ASSERT_EQ((eps_qos.bitRatesExt.maxBitRateForDL), 0);
  ASSERT_EQ((eps_qos.bitRatesExt2.maxBitRateForUL), 0);
  ASSERT_EQ((eps_qos.bitRatesExt2.maxBitRateForDL), 0);
  uint8_t eps_qos_buffersize_encoded3[10] = {0};
  ASSERT_EQ(
      encode_eps_quality_of_service(&eps_qos, 0, eps_qos_buffersize_encoded3,
                                    sizeof(eps_qos_buffersize_encoded3)),
      sizeof(eps_qos_buffersize_encoded3));
  ASSERT_EQ(memcmp(eps_qos_buffersize_encoded3, eps_qos_buffersize3,
                   sizeof(eps_qos_buffersize3)),
            0);

  uint8_t eps_qos_buffersize4[] = {0x09, 0x09, 0x7F, 0x7F, 0x7F,
                                   0x7F, 0x00, 0x00, 0x00, 0x00};

  ASSERT_EQ(decode_eps_quality_of_service(&eps_qos, 0, eps_qos_buffersize4,
                                          sizeof(eps_qos_buffersize4)),
            sizeof(eps_qos_buffersize4));
  ASSERT_EQ((eps_qos.qci), 9);
  ASSERT_EQ((eps_qos.bitRates.maxBitRateForUL), 127);
  ASSERT_EQ((eps_qos.bitRates.maxBitRateForDL), 127);
  ASSERT_EQ((eps_qos.bitRatesExt.maxBitRateForUL), 0);
  ASSERT_EQ((eps_qos.bitRatesExt.maxBitRateForDL), 0);
  ASSERT_EQ((eps_qos.bitRatesExt2.maxBitRateForUL), 0);
  ASSERT_EQ((eps_qos.bitRatesExt2.maxBitRateForDL), 0);
  uint8_t eps_qos_buffersize_encoded4[10] = {0};
  ASSERT_EQ(
      encode_eps_quality_of_service(&eps_qos, 0, eps_qos_buffersize_encoded4,
                                    sizeof(eps_qos_buffersize_encoded4)),
      sizeof(eps_qos_buffersize_encoded4));
  ASSERT_EQ(memcmp(eps_qos_buffersize_encoded4, eps_qos_buffersize4,
                   sizeof(eps_qos_buffersize4)),
            0);

  uint8_t eps_qos_buffersize5[] = {0x09, 0x09, 0x81, 0x81, 0x81,
                                   0x81, 0x00, 0x00, 0x00, 0x00};

  ASSERT_EQ(decode_eps_quality_of_service(&eps_qos, 0, eps_qos_buffersize5,
                                          sizeof(eps_qos_buffersize5)),
            sizeof(eps_qos_buffersize5));
  ASSERT_EQ((eps_qos.qci), 9);
  ASSERT_EQ((eps_qos.bitRates.maxBitRateForUL), 129);
  ASSERT_EQ((eps_qos.bitRates.maxBitRateForDL), 129);
  ASSERT_EQ((eps_qos.bitRatesExt.maxBitRateForUL), 0);
  ASSERT_EQ((eps_qos.bitRatesExt.maxBitRateForDL), 0);
  ASSERT_EQ((eps_qos.bitRatesExt2.maxBitRateForUL), 0);
  ASSERT_EQ((eps_qos.bitRatesExt2.maxBitRateForDL), 0);
  uint8_t eps_qos_buffersize_encoded5[10] = {0};
  ASSERT_EQ(
      encode_eps_quality_of_service(&eps_qos, 0, eps_qos_buffersize_encoded5,
                                    sizeof(eps_qos_buffersize_encoded5)),
      sizeof(eps_qos_buffersize_encoded5));
  ASSERT_EQ(memcmp(eps_qos_buffersize_encoded5, eps_qos_buffersize5,
                   sizeof(eps_qos_buffersize5)),
            0);

  uint8_t eps_qos_buffersize6[] = {0x09, 0x09, 0x82, 0x82, 0x82,
                                   0x82, 0x00, 0x00, 0x00, 0x00};

  ASSERT_EQ(decode_eps_quality_of_service(&eps_qos, 0, eps_qos_buffersize6,
                                          sizeof(eps_qos_buffersize6)),
            sizeof(eps_qos_buffersize6));
  ASSERT_EQ((eps_qos.qci), 9);
  ASSERT_EQ((eps_qos.bitRates.maxBitRateForUL), 130);
  ASSERT_EQ((eps_qos.bitRates.maxBitRateForDL), 130);
  ASSERT_EQ((eps_qos.bitRatesExt.maxBitRateForUL), 0);
  ASSERT_EQ((eps_qos.bitRatesExt.maxBitRateForDL), 0);
  ASSERT_EQ((eps_qos.bitRatesExt2.maxBitRateForUL), 0);
  ASSERT_EQ((eps_qos.bitRatesExt2.maxBitRateForDL), 0);
  uint8_t eps_qos_buffersize_encoded6[10] = {0};
  ASSERT_EQ(
      encode_eps_quality_of_service(&eps_qos, 0, eps_qos_buffersize_encoded6,
                                    sizeof(eps_qos_buffersize_encoded6)),
      sizeof(eps_qos_buffersize_encoded6));
  ASSERT_EQ(memcmp(eps_qos_buffersize_encoded6, eps_qos_buffersize6,
                   sizeof(eps_qos_buffersize6)),
            0);

  uint8_t eps_qos_buffersize7[] = {0x09, 0x09, 0x82, 0x82, 0x82,
                                   0x82, 0x00, 0x00, 0x00, 0x00};

  ASSERT_EQ(decode_eps_quality_of_service(&eps_qos, 0, eps_qos_buffersize7,
                                          sizeof(eps_qos_buffersize7)),
            sizeof(eps_qos_buffersize7));
  ASSERT_EQ((eps_qos.qci), 9);
  ASSERT_EQ((eps_qos.bitRates.maxBitRateForUL), 130);
  ASSERT_EQ((eps_qos.bitRates.maxBitRateForDL), 130);
  ASSERT_EQ((eps_qos.bitRatesExt.maxBitRateForUL), 0);
  ASSERT_EQ((eps_qos.bitRatesExt.maxBitRateForDL), 0);
  ASSERT_EQ((eps_qos.bitRatesExt2.maxBitRateForUL), 0);
  ASSERT_EQ((eps_qos.bitRatesExt2.maxBitRateForDL), 0);
  uint8_t eps_qos_buffersize_encoded7[10] = {0};
  ASSERT_EQ(
      encode_eps_quality_of_service(&eps_qos, 0, eps_qos_buffersize_encoded7,
                                    sizeof(eps_qos_buffersize_encoded7)),
      sizeof(eps_qos_buffersize_encoded7));
  ASSERT_EQ(memcmp(eps_qos_buffersize_encoded7, eps_qos_buffersize7,
                   sizeof(eps_qos_buffersize7)),
            0);

  uint8_t eps_qos_buffersize8[] = {0x09, 0x09, 0xFE, 0xFE, 0xFE,
                                   0xFE, 0x00, 0x00, 0x00, 0x00};

  ASSERT_EQ(decode_eps_quality_of_service(&eps_qos, 0, eps_qos_buffersize8,
                                          sizeof(eps_qos_buffersize8)),
            sizeof(eps_qos_buffersize8));
  ASSERT_EQ((eps_qos.qci), 9);
  ASSERT_EQ((eps_qos.bitRates.maxBitRateForUL), 254);
  ASSERT_EQ((eps_qos.bitRates.maxBitRateForDL), 254);
  ASSERT_EQ((eps_qos.bitRatesExt.maxBitRateForUL), 0);
  ASSERT_EQ((eps_qos.bitRatesExt.maxBitRateForDL), 0);
  ASSERT_EQ((eps_qos.bitRatesExt2.maxBitRateForUL), 0);
  ASSERT_EQ((eps_qos.bitRatesExt2.maxBitRateForDL), 0);
  uint8_t eps_qos_buffersize_encoded8[10] = {0};
  ASSERT_EQ(
      encode_eps_quality_of_service(&eps_qos, 0, eps_qos_buffersize_encoded8,
                                    sizeof(eps_qos_buffersize_encoded8)),
      sizeof(eps_qos_buffersize_encoded8));
  ASSERT_EQ(memcmp(eps_qos_buffersize_encoded8, eps_qos_buffersize8,
                   sizeof(eps_qos_buffersize8)),
            0);

  uint8_t eps_qos_buffersize_Ext1[] = {0x09, 0x09, 0xFE, 0xFE, 0xFE,
                                       0xFE, 0x02, 0x02, 0x02, 0x02};

  ASSERT_EQ(decode_eps_quality_of_service(&eps_qos, 0, eps_qos_buffersize_Ext1,
                                          sizeof(eps_qos_buffersize_Ext1)),
            sizeof(eps_qos_buffersize_Ext1));
  ASSERT_EQ((eps_qos.qci), 9);
  ASSERT_EQ((eps_qos.bitRates.maxBitRateForUL), 254);
  ASSERT_EQ((eps_qos.bitRates.maxBitRateForDL), 254);
  ASSERT_EQ((eps_qos.bitRatesExt.maxBitRateForUL), 2);
  ASSERT_EQ((eps_qos.bitRatesExt.maxBitRateForDL), 2);
  ASSERT_EQ((eps_qos.bitRatesExt2.maxBitRateForUL), 0);
  ASSERT_EQ((eps_qos.bitRatesExt2.maxBitRateForDL), 0);
  uint8_t eps_qos_buffersize_encoded_Ext1[10] = {0};
  ASSERT_EQ(encode_eps_quality_of_service(
                &eps_qos, 0, eps_qos_buffersize_encoded_Ext1,
                sizeof(eps_qos_buffersize_encoded_Ext1)),
            sizeof(eps_qos_buffersize_encoded_Ext1));
  ASSERT_EQ(memcmp(eps_qos_buffersize_encoded_Ext1, eps_qos_buffersize_Ext1,
                   sizeof(eps_qos_buffersize_Ext1)),
            0);

  uint8_t eps_qos_buffersize_Ext2[] = {0x09, 0x09, 0xFE, 0xFE, 0xFE,
                                       0xFE, 0x4A, 0x4A, 0x4A, 0x4A};

  ASSERT_EQ(decode_eps_quality_of_service(&eps_qos, 0, eps_qos_buffersize_Ext2,
                                          sizeof(eps_qos_buffersize_Ext2)),
            sizeof(eps_qos_buffersize_Ext2));
  ASSERT_EQ((eps_qos.qci), 9);
  ASSERT_EQ((eps_qos.bitRates.maxBitRateForUL), 254);
  ASSERT_EQ((eps_qos.bitRates.maxBitRateForDL), 254);
  ASSERT_EQ((eps_qos.bitRatesExt.maxBitRateForUL), 74);
  ASSERT_EQ((eps_qos.bitRatesExt.maxBitRateForDL), 74);
  ASSERT_EQ((eps_qos.bitRatesExt2.maxBitRateForUL), 0);
  ASSERT_EQ((eps_qos.bitRatesExt2.maxBitRateForDL), 0);
  uint8_t eps_qos_buffersize_encoded_Ext2[10] = {0};
  ASSERT_EQ(encode_eps_quality_of_service(
                &eps_qos, 0, eps_qos_buffersize_encoded_Ext2,
                sizeof(eps_qos_buffersize_encoded_Ext2)),
            sizeof(eps_qos_buffersize_encoded_Ext2));
  ASSERT_EQ(memcmp(eps_qos_buffersize_encoded_Ext2, eps_qos_buffersize_Ext2,
                   sizeof(eps_qos_buffersize_Ext2)),
            0);

  uint8_t eps_qos_buffersize_Ext3[] = {0x09, 0x09, 0xFE, 0xFE, 0xFE,
                                       0xFE, 0x4B, 0x4B, 0x4B, 0x4B};

  ASSERT_EQ(decode_eps_quality_of_service(&eps_qos, 0, eps_qos_buffersize_Ext3,
                                          sizeof(eps_qos_buffersize_Ext3)),
            sizeof(eps_qos_buffersize_Ext3));
  ASSERT_EQ((eps_qos.qci), 9);
  ASSERT_EQ((eps_qos.bitRates.maxBitRateForUL), 254);
  ASSERT_EQ((eps_qos.bitRates.maxBitRateForDL), 254);
  ASSERT_EQ((eps_qos.bitRatesExt.maxBitRateForUL), 75);
  ASSERT_EQ((eps_qos.bitRatesExt.maxBitRateForDL), 75);
  ASSERT_EQ((eps_qos.bitRatesExt2.maxBitRateForUL), 0);
  ASSERT_EQ((eps_qos.bitRatesExt2.maxBitRateForDL), 0);
  uint8_t eps_qos_buffersize_encoded_Ext3[10] = {0};
  ASSERT_EQ(encode_eps_quality_of_service(
                &eps_qos, 0, eps_qos_buffersize_encoded_Ext3,
                sizeof(eps_qos_buffersize_encoded_Ext3)),
            sizeof(eps_qos_buffersize_encoded_Ext3));
  ASSERT_EQ(memcmp(eps_qos_buffersize_encoded_Ext3, eps_qos_buffersize_Ext3,
                   sizeof(eps_qos_buffersize_Ext3)),
            0);

  uint8_t eps_qos_buffersize_Ext4[] = {0x09, 0x09, 0xFE, 0xFE, 0xFE,
                                       0xFE, 0x4E, 0x4E, 0x4E, 0x4E};

  ASSERT_EQ(decode_eps_quality_of_service(&eps_qos, 0, eps_qos_buffersize_Ext4,
                                          sizeof(eps_qos_buffersize_Ext4)),
            sizeof(eps_qos_buffersize_Ext4));
  ASSERT_EQ((eps_qos.qci), 9);
  ASSERT_EQ((eps_qos.bitRates.maxBitRateForUL), 254);
  ASSERT_EQ((eps_qos.bitRates.maxBitRateForDL), 254);
  ASSERT_EQ((eps_qos.bitRatesExt.maxBitRateForUL), 78);
  ASSERT_EQ((eps_qos.bitRatesExt.maxBitRateForDL), 78);
  ASSERT_EQ((eps_qos.bitRatesExt2.maxBitRateForUL), 0);
  ASSERT_EQ((eps_qos.bitRatesExt2.maxBitRateForDL), 0);
  uint8_t eps_qos_buffersize_encoded_Ext4[10] = {0};
  ASSERT_EQ(encode_eps_quality_of_service(
                &eps_qos, 0, eps_qos_buffersize_encoded_Ext4,
                sizeof(eps_qos_buffersize_encoded_Ext4)),
            sizeof(eps_qos_buffersize_encoded_Ext4));
  ASSERT_EQ(memcmp(eps_qos_buffersize_encoded_Ext4, eps_qos_buffersize_Ext4,
                   sizeof(eps_qos_buffersize_Ext4)),
            0);

  uint8_t eps_qos_buffersize_Ext5[] = {0x09, 0x09, 0xFE, 0xFE, 0xFE,
                                       0xFE, 0xBA, 0xBA, 0xBA, 0xBA};

  ASSERT_EQ(decode_eps_quality_of_service(&eps_qos, 0, eps_qos_buffersize_Ext5,
                                          sizeof(eps_qos_buffersize_Ext5)),
            sizeof(eps_qos_buffersize_Ext5));
  ASSERT_EQ((eps_qos.qci), 9);
  ASSERT_EQ((eps_qos.bitRates.maxBitRateForUL), 254);
  ASSERT_EQ((eps_qos.bitRates.maxBitRateForDL), 254);
  ASSERT_EQ((eps_qos.bitRatesExt.maxBitRateForUL), 186);
  ASSERT_EQ((eps_qos.bitRatesExt.maxBitRateForDL), 186);
  ASSERT_EQ((eps_qos.bitRatesExt2.maxBitRateForUL), 0);
  ASSERT_EQ((eps_qos.bitRatesExt2.maxBitRateForDL), 0);
  uint8_t eps_qos_buffersize_encoded_Ext5[10] = {0};
  ASSERT_EQ(encode_eps_quality_of_service(
                &eps_qos, 0, eps_qos_buffersize_encoded_Ext5,
                sizeof(eps_qos_buffersize_encoded_Ext5)),
            sizeof(eps_qos_buffersize_encoded_Ext5));
  ASSERT_EQ(memcmp(eps_qos_buffersize_encoded_Ext5, eps_qos_buffersize_Ext5,
                   sizeof(eps_qos_buffersize_Ext5)),
            0);

  uint8_t eps_qos_buffersize_Ext6[] = {0x09, 0x09, 0xFE, 0xFE, 0xFE,
                                       0xFE, 0xBB, 0xBB, 0xBB, 0xBB};

  ASSERT_EQ(decode_eps_quality_of_service(&eps_qos, 0, eps_qos_buffersize_Ext6,
                                          sizeof(eps_qos_buffersize_Ext6)),
            sizeof(eps_qos_buffersize_Ext6));
  ASSERT_EQ((eps_qos.qci), 9);
  ASSERT_EQ((eps_qos.bitRates.maxBitRateForUL), 254);
  ASSERT_EQ((eps_qos.bitRates.maxBitRateForDL), 254);
  ASSERT_EQ((eps_qos.bitRatesExt.maxBitRateForUL), 187);
  ASSERT_EQ((eps_qos.bitRatesExt.maxBitRateForDL), 187);
  ASSERT_EQ((eps_qos.bitRatesExt2.maxBitRateForUL), 0);
  ASSERT_EQ((eps_qos.bitRatesExt2.maxBitRateForDL), 0);
  uint8_t eps_qos_buffersize_encoded_Ext6[10] = {0};
  ASSERT_EQ(encode_eps_quality_of_service(
                &eps_qos, 0, eps_qos_buffersize_encoded_Ext6,
                sizeof(eps_qos_buffersize_encoded_Ext6)),
            sizeof(eps_qos_buffersize_encoded_Ext6));
  ASSERT_EQ(memcmp(eps_qos_buffersize_encoded_Ext6, eps_qos_buffersize_Ext6,
                   sizeof(eps_qos_buffersize_Ext6)),
            0);

  uint8_t eps_qos_buffersize_Ext7[] = {0x09, 0x09, 0xFE, 0xFE, 0xFE,
                                       0xFE, 0xC5, 0xC5, 0xC5, 0xC5};

  ASSERT_EQ(decode_eps_quality_of_service(&eps_qos, 0, eps_qos_buffersize_Ext7,
                                          sizeof(eps_qos_buffersize_Ext7)),
            sizeof(eps_qos_buffersize_Ext7));
  ASSERT_EQ((eps_qos.qci), 9);
  ASSERT_EQ((eps_qos.bitRates.maxBitRateForUL), 254);
  ASSERT_EQ((eps_qos.bitRates.maxBitRateForDL), 254);
  ASSERT_EQ((eps_qos.bitRatesExt.maxBitRateForUL), 197);
  ASSERT_EQ((eps_qos.bitRatesExt.maxBitRateForDL), 197);
  ASSERT_EQ((eps_qos.bitRatesExt2.maxBitRateForUL), 0);
  ASSERT_EQ((eps_qos.bitRatesExt2.maxBitRateForDL), 0);
  uint8_t eps_qos_buffersize_encoded_Ext7[10] = {0};
  ASSERT_EQ(encode_eps_quality_of_service(
                &eps_qos, 0, eps_qos_buffersize_encoded_Ext7,
                sizeof(eps_qos_buffersize_encoded_Ext7)),
            sizeof(eps_qos_buffersize_encoded_Ext7));
  ASSERT_EQ(memcmp(eps_qos_buffersize_encoded_Ext7, eps_qos_buffersize_Ext7,
                   sizeof(eps_qos_buffersize_Ext7)),
            0);

  uint8_t eps_qos_buffersize_Ext8[] = {0x09, 0x09, 0xFE, 0xFE, 0xFE,
                                       0xFE, 0xFA, 0xFA, 0xFA, 0xFA};

  ASSERT_EQ(decode_eps_quality_of_service(&eps_qos, 0, eps_qos_buffersize_Ext8,
                                          sizeof(eps_qos_buffersize_Ext8)),
            sizeof(eps_qos_buffersize_Ext8));
  ASSERT_EQ((eps_qos.qci), 9);
  ASSERT_EQ((eps_qos.bitRates.maxBitRateForUL), 254);
  ASSERT_EQ((eps_qos.bitRates.maxBitRateForDL), 254);
  ASSERT_EQ((eps_qos.bitRatesExt.maxBitRateForUL), 250);
  ASSERT_EQ((eps_qos.bitRatesExt.maxBitRateForDL), 250);
  ASSERT_EQ((eps_qos.bitRatesExt2.maxBitRateForUL), 0);
  ASSERT_EQ((eps_qos.bitRatesExt2.maxBitRateForDL), 0);
  uint8_t eps_qos_buffersize_encoded_Ext8[10] = {0};
  ASSERT_EQ(encode_eps_quality_of_service(
                &eps_qos, 0, eps_qos_buffersize_encoded_Ext8,
                sizeof(eps_qos_buffersize_encoded_Ext8)),
            sizeof(eps_qos_buffersize_encoded_Ext8));
  ASSERT_EQ(memcmp(eps_qos_buffersize_encoded_Ext8, eps_qos_buffersize_Ext8,
                   sizeof(eps_qos_buffersize_Ext8)),
            0);

  uint8_t eps_qos_buffersize_Extended21[] = {0x0D, 0x09, 0xFE, 0xFE, 0xFE,
                                             0xFE, 0xFA, 0xFA, 0xFA, 0xFA,
                                             0x01, 0x01, 0x01, 0x01};
  ASSERT_EQ(
      decode_eps_quality_of_service(&eps_qos, 0, eps_qos_buffersize_Extended21,
                                    sizeof(eps_qos_buffersize_Extended21)),
      sizeof(eps_qos_buffersize_Extended21));
  ASSERT_EQ((eps_qos.bitRates.maxBitRateForUL), 254);
  ASSERT_EQ((eps_qos.bitRates.maxBitRateForDL), 254);
  ASSERT_EQ((eps_qos.bitRatesExt.maxBitRateForUL), 250);
  ASSERT_EQ((eps_qos.bitRatesExt.maxBitRateForDL), 250);
  ASSERT_EQ((eps_qos.bitRatesExt2.maxBitRateForUL), 1);
  ASSERT_EQ((eps_qos.bitRatesExt2.maxBitRateForDL), 1);
  uint8_t eps_qos_buffersize_encoded_Extended21[14] = {0};
  ASSERT_EQ(encode_eps_quality_of_service(
                &eps_qos, 0, eps_qos_buffersize_encoded_Extended21,
                sizeof(eps_qos_buffersize_encoded_Extended21)),
            sizeof(eps_qos_buffersize_encoded_Extended21));
  ASSERT_EQ(memcmp(eps_qos_buffersize_encoded_Extended21,
                   eps_qos_buffersize_Extended21,
                   sizeof(eps_qos_buffersize_Extended21)),
            0);

  uint8_t eps_qos_buffersize_Extended22[] = {0x0D, 0x09, 0xFE, 0xFE, 0xFE,
                                             0xFE, 0xFA, 0xFA, 0xFA, 0xFA,
                                             0x0B, 0x0B, 0x0B, 0x0B};
  ASSERT_EQ(
      decode_eps_quality_of_service(&eps_qos, 0, eps_qos_buffersize_Extended22,
                                    sizeof(eps_qos_buffersize_Extended22)),
      sizeof(eps_qos_buffersize_Extended21));
  ASSERT_EQ((eps_qos.bitRates.maxBitRateForUL), 254);
  ASSERT_EQ((eps_qos.bitRates.maxBitRateForDL), 254);
  ASSERT_EQ((eps_qos.bitRatesExt.maxBitRateForUL), 250);
  ASSERT_EQ((eps_qos.bitRatesExt.maxBitRateForDL), 250);
  ASSERT_EQ((eps_qos.bitRatesExt2.maxBitRateForUL), 11);
  ASSERT_EQ((eps_qos.bitRatesExt2.maxBitRateForDL), 11);
  uint8_t eps_qos_buffersize_encoded_Extended22[14] = {0};
  ASSERT_EQ(encode_eps_quality_of_service(
                &eps_qos, 0, eps_qos_buffersize_encoded_Extended22,
                sizeof(eps_qos_buffersize_encoded_Extended22)),
            sizeof(eps_qos_buffersize_encoded_Extended22));
  ASSERT_EQ(memcmp(eps_qos_buffersize_encoded_Extended22,
                   eps_qos_buffersize_Extended22,
                   sizeof(eps_qos_buffersize_Extended22)),
            0);

  uint8_t eps_qos_buffersize_Extended23[] = {0x0D, 0x09, 0xFE, 0xFE, 0xFE,
                                             0xFE, 0xFA, 0xFA, 0xFA, 0xFA,
                                             0x3D, 0x3D, 0x3D, 0x3D};
  ASSERT_EQ(
      decode_eps_quality_of_service(&eps_qos, 0, eps_qos_buffersize_Extended23,
                                    sizeof(eps_qos_buffersize_Extended23)),
      sizeof(eps_qos_buffersize_Extended21));
  ASSERT_EQ((eps_qos.bitRates.maxBitRateForUL), 254);
  ASSERT_EQ((eps_qos.bitRates.maxBitRateForDL), 254);
  ASSERT_EQ((eps_qos.bitRatesExt.maxBitRateForUL), 250);
  ASSERT_EQ((eps_qos.bitRatesExt.maxBitRateForDL), 250);
  ASSERT_EQ((eps_qos.bitRatesExt2.maxBitRateForUL), 61);
  ASSERT_EQ((eps_qos.bitRatesExt2.maxBitRateForDL), 61);
  uint8_t eps_qos_buffersize_encoded_Extended23[14] = {0};
  ASSERT_EQ(encode_eps_quality_of_service(
                &eps_qos, 0, eps_qos_buffersize_encoded_Extended23,
                sizeof(eps_qos_buffersize_encoded_Extended23)),
            sizeof(eps_qos_buffersize_encoded_Extended23));
  ASSERT_EQ(memcmp(eps_qos_buffersize_encoded_Extended23,
                   eps_qos_buffersize_Extended23,
                   sizeof(eps_qos_buffersize_Extended23)),
            0);

  uint8_t eps_qos_buffersize_Extended24[] = {0x0D, 0x09, 0xFE, 0xFE, 0xFE,
                                             0xFE, 0xFA, 0xFA, 0xFA, 0xFA,
                                             0x3E, 0x3E, 0x3E, 0x3E};
  ASSERT_EQ(
      decode_eps_quality_of_service(&eps_qos, 0, eps_qos_buffersize_Extended24,
                                    sizeof(eps_qos_buffersize_Extended24)),
      sizeof(eps_qos_buffersize_Extended21));
  ASSERT_EQ((eps_qos.bitRates.maxBitRateForUL), 254);
  ASSERT_EQ((eps_qos.bitRates.maxBitRateForDL), 254);
  ASSERT_EQ((eps_qos.bitRatesExt.maxBitRateForUL), 250);
  ASSERT_EQ((eps_qos.bitRatesExt.maxBitRateForDL), 250);
  ASSERT_EQ((eps_qos.bitRatesExt2.maxBitRateForUL), 62);
  ASSERT_EQ((eps_qos.bitRatesExt2.maxBitRateForDL), 62);
  uint8_t eps_qos_buffersize_encoded_Extended24[14] = {0};
  ASSERT_EQ(encode_eps_quality_of_service(
                &eps_qos, 0, eps_qos_buffersize_encoded_Extended24,
                sizeof(eps_qos_buffersize_encoded_Extended24)),
            sizeof(eps_qos_buffersize_encoded_Extended24));
  ASSERT_EQ(memcmp(eps_qos_buffersize_encoded_Extended24,
                   eps_qos_buffersize_Extended24,
                   sizeof(eps_qos_buffersize_Extended24)),
            0);

  uint8_t eps_qos_buffersize_Extended25[] = {0x0D, 0x09, 0xFE, 0xFE, 0xFE,
                                             0xFE, 0xFA, 0xFA, 0xFA, 0xFA,
                                             0x47, 0x47, 0x47, 0x47};
  ASSERT_EQ(
      decode_eps_quality_of_service(&eps_qos, 0, eps_qos_buffersize_Extended25,
                                    sizeof(eps_qos_buffersize_Extended25)),
      sizeof(eps_qos_buffersize_Extended21));
  ASSERT_EQ((eps_qos.bitRates.maxBitRateForUL), 254);
  ASSERT_EQ((eps_qos.bitRates.maxBitRateForDL), 254);
  ASSERT_EQ((eps_qos.bitRatesExt.maxBitRateForUL), 250);
  ASSERT_EQ((eps_qos.bitRatesExt.maxBitRateForDL), 250);
  ASSERT_EQ((eps_qos.bitRatesExt2.maxBitRateForUL), 71);
  ASSERT_EQ((eps_qos.bitRatesExt2.maxBitRateForDL), 71);
  uint8_t eps_qos_buffersize_encoded_Extended25[14] = {0};
  ASSERT_EQ(encode_eps_quality_of_service(
                &eps_qos, 0, eps_qos_buffersize_encoded_Extended25,
                sizeof(eps_qos_buffersize_encoded_Extended25)),
            sizeof(eps_qos_buffersize_encoded_Extended25));
  ASSERT_EQ(memcmp(eps_qos_buffersize_encoded_Extended25,
                   eps_qos_buffersize_Extended25,
                   sizeof(eps_qos_buffersize_Extended25)),
            0);

  uint8_t eps_qos_buffersize_Extended26[] = {0x0D, 0x09, 0xFE, 0xFE, 0xFE,
                                             0xFE, 0xFA, 0xFA, 0xFA, 0xFA,
                                             0xA1, 0xA1, 0xA1, 0xA1};
  ASSERT_EQ(
      decode_eps_quality_of_service(&eps_qos, 0, eps_qos_buffersize_Extended26,
                                    sizeof(eps_qos_buffersize_Extended26)),
      sizeof(eps_qos_buffersize_Extended26));
  ASSERT_EQ((eps_qos.bitRates.maxBitRateForUL), 254);
  ASSERT_EQ((eps_qos.bitRates.maxBitRateForDL), 254);
  ASSERT_EQ((eps_qos.bitRatesExt.maxBitRateForUL), 250);
  ASSERT_EQ((eps_qos.bitRatesExt.maxBitRateForDL), 250);
  ASSERT_EQ((eps_qos.bitRatesExt2.maxBitRateForUL), 161);
  ASSERT_EQ((eps_qos.bitRatesExt2.maxBitRateForDL), 161);
  uint8_t eps_qos_buffersize_encoded_Extended26[14] = {0};
  ASSERT_EQ(encode_eps_quality_of_service(
                &eps_qos, 0, eps_qos_buffersize_encoded_Extended26,
                sizeof(eps_qos_buffersize_encoded_Extended26)),
            sizeof(eps_qos_buffersize_encoded_Extended26));
  ASSERT_EQ(memcmp(eps_qos_buffersize_encoded_Extended26,
                   eps_qos_buffersize_Extended26,
                   sizeof(eps_qos_buffersize_Extended26)),
            0);

  uint8_t eps_qos_buffersize_Extended27[] = {0x0D, 0x09, 0xFE, 0xFE, 0xFE,
                                             0xFE, 0xFA, 0xFA, 0xFA, 0xFA,
                                             0xA2, 0xA2, 0xA2, 0xA2};
  ASSERT_EQ(
      decode_eps_quality_of_service(&eps_qos, 0, eps_qos_buffersize_Extended27,
                                    sizeof(eps_qos_buffersize_Extended27)),
      sizeof(eps_qos_buffersize_Extended27));
  ASSERT_EQ((eps_qos.bitRates.maxBitRateForUL), 254);
  ASSERT_EQ((eps_qos.bitRates.maxBitRateForDL), 254);
  ASSERT_EQ((eps_qos.bitRatesExt.maxBitRateForUL), 250);
  ASSERT_EQ((eps_qos.bitRatesExt.maxBitRateForDL), 250);
  ASSERT_EQ((eps_qos.bitRatesExt2.maxBitRateForUL), 162);
  ASSERT_EQ((eps_qos.bitRatesExt2.maxBitRateForDL), 162);
  uint8_t eps_qos_buffersize_encoded_Extended27[14] = {0};
  ASSERT_EQ(encode_eps_quality_of_service(
                &eps_qos, 0, eps_qos_buffersize_encoded_Extended27,
                sizeof(eps_qos_buffersize_encoded_Extended27)),
            sizeof(eps_qos_buffersize_encoded_Extended27));
  ASSERT_EQ(memcmp(eps_qos_buffersize_encoded_Extended27,
                   eps_qos_buffersize_Extended27,
                   sizeof(eps_qos_buffersize_Extended27)),
            0);

  uint8_t eps_qos_buffersize_Extended28[] = {0x0D, 0x09, 0xFE, 0xFE, 0xFE,
                                             0xFE, 0xFA, 0xFA, 0xFA, 0xFA,
                                             0xA6, 0xA6, 0xA6, 0xA6};
  ASSERT_EQ(
      decode_eps_quality_of_service(&eps_qos, 0, eps_qos_buffersize_Extended28,
                                    sizeof(eps_qos_buffersize_Extended28)),
      sizeof(eps_qos_buffersize_Extended28));
  ASSERT_EQ((eps_qos.bitRates.maxBitRateForUL), 254);
  ASSERT_EQ((eps_qos.bitRates.maxBitRateForDL), 254);
  ASSERT_EQ((eps_qos.bitRatesExt.maxBitRateForUL), 250);
  ASSERT_EQ((eps_qos.bitRatesExt.maxBitRateForDL), 250);
  ASSERT_EQ((eps_qos.bitRatesExt2.maxBitRateForUL), 166);
  ASSERT_EQ((eps_qos.bitRatesExt2.maxBitRateForDL), 166);
  uint8_t eps_qos_buffersize_encoded_Extended28[14] = {0};
  ASSERT_EQ(encode_eps_quality_of_service(
                &eps_qos, 0, eps_qos_buffersize_encoded_Extended28,
                sizeof(eps_qos_buffersize_encoded_Extended28)),
            sizeof(eps_qos_buffersize_encoded_Extended28));
  ASSERT_EQ(memcmp(eps_qos_buffersize_encoded_Extended28,
                   eps_qos_buffersize_Extended28,
                   sizeof(eps_qos_buffersize_Extended28)),
            0);

  uint8_t eps_qos_buffersize_Extended29[] = {0x0D, 0x09, 0xFE, 0xFE, 0xFE,
                                             0xFE, 0xFA, 0xFA, 0xFA, 0xFA,
                                             0xF6, 0xF6, 0xF6, 0xF6};
  ASSERT_EQ(
      decode_eps_quality_of_service(&eps_qos, 0, eps_qos_buffersize_Extended29,
                                    sizeof(eps_qos_buffersize_Extended29)),
      sizeof(eps_qos_buffersize_Extended29));
  ASSERT_EQ((eps_qos.bitRates.maxBitRateForUL), 254);
  ASSERT_EQ((eps_qos.bitRates.maxBitRateForDL), 254);
  ASSERT_EQ((eps_qos.bitRatesExt.maxBitRateForUL), 250);
  ASSERT_EQ((eps_qos.bitRatesExt.maxBitRateForDL), 250);
  ASSERT_EQ((eps_qos.bitRatesExt2.maxBitRateForUL), 246);
  ASSERT_EQ((eps_qos.bitRatesExt2.maxBitRateForDL), 246);
  uint8_t eps_qos_buffersize_encoded_Extended29[14] = {0};
  ASSERT_EQ(encode_eps_quality_of_service(
                &eps_qos, 0, eps_qos_buffersize_encoded_Extended29,
                sizeof(eps_qos_buffersize_encoded_Extended29)),
            sizeof(eps_qos_buffersize_encoded_Extended29));
  ASSERT_EQ(memcmp(eps_qos_buffersize_encoded_Extended29,
                   eps_qos_buffersize_Extended29,
                   sizeof(eps_qos_buffersize_Extended29)),
            0);
}

TEST_F(NASEncodeDecodeTest, TestDecodeEncodeAPNAMBR) {
  ApnAggregateMaximumBitRate apn_ambr = {0};

  uint8_t apn_ambr_buffersize1[] = {0x06, 0xFF, 0xFF, 0x00, 0x00, 0x00, 0x00};
  ASSERT_EQ(
      decode_apn_aggregate_maximum_bit_rate(&apn_ambr, 0, apn_ambr_buffersize1,
                                            sizeof(apn_ambr_buffersize1)),
      sizeof(apn_ambr_buffersize1));
  ASSERT_EQ((apn_ambr.apnambrfordownlink), 255);
  ASSERT_EQ((apn_ambr.apnambrforuplink), 255);
  ASSERT_EQ((apn_ambr.apnambrfordownlink_extended), 0);
  ASSERT_EQ((apn_ambr.apnambrforuplink_extended), 0);

  uint8_t apn_ambr_encoded_buffersize1[7] = {0};
  ASSERT_EQ(encode_apn_aggregate_maximum_bit_rate(
                &apn_ambr, 0, apn_ambr_encoded_buffersize1,
                sizeof(apn_ambr_encoded_buffersize1)),
            sizeof(apn_ambr_encoded_buffersize1));
  ASSERT_EQ((memcmp(apn_ambr_encoded_buffersize1, apn_ambr_buffersize1,
                    sizeof(apn_ambr_buffersize1))),
            0);

  uint8_t apn_ambr_buffersize2[] = {0x06, 0x01, 0x01, 0x00, 0x00, 0x00, 0x00};
  ASSERT_EQ(
      decode_apn_aggregate_maximum_bit_rate(&apn_ambr, 0, apn_ambr_buffersize2,
                                            sizeof(apn_ambr_buffersize2)),
      sizeof(apn_ambr_buffersize2));
  ASSERT_EQ((apn_ambr.apnambrfordownlink), 1);
  ASSERT_EQ((apn_ambr.apnambrforuplink), 1);
  ASSERT_EQ((apn_ambr.apnambrfordownlink_extended), 0);
  ASSERT_EQ((apn_ambr.apnambrforuplink_extended), 0);
  ASSERT_EQ((apn_ambr.apnambrfordownlink_extended2), 0);
  ASSERT_EQ((apn_ambr.apnambrforuplink_extended2), 0);

  uint8_t apn_ambr_encoded_buffersize2[7] = {0};
  ASSERT_EQ(encode_apn_aggregate_maximum_bit_rate(
                &apn_ambr, 0, apn_ambr_encoded_buffersize2,
                sizeof(apn_ambr_encoded_buffersize2)),
            sizeof(apn_ambr_encoded_buffersize2));
  ASSERT_EQ((memcmp(apn_ambr_encoded_buffersize2, apn_ambr_buffersize2,
                    sizeof(apn_ambr_buffersize2))),
            0);

  uint8_t apn_ambr_buffersize3[] = {0x06, 0x3F, 0x3F, 0x00, 0x00, 0x00, 0x00};
  ASSERT_EQ(
      decode_apn_aggregate_maximum_bit_rate(&apn_ambr, 0, apn_ambr_buffersize3,
                                            sizeof(apn_ambr_buffersize3)),
      sizeof(apn_ambr_buffersize2));
  ASSERT_EQ((apn_ambr.apnambrfordownlink), 63);
  ASSERT_EQ((apn_ambr.apnambrforuplink), 63);
  ASSERT_EQ((apn_ambr.apnambrfordownlink_extended), 0);
  ASSERT_EQ((apn_ambr.apnambrforuplink_extended), 0);
  ASSERT_EQ((apn_ambr.apnambrfordownlink_extended2), 0);
  ASSERT_EQ((apn_ambr.apnambrforuplink_extended2), 0);

  uint8_t apn_ambr_encoded_buffersize3[7] = {0};
  ASSERT_EQ(encode_apn_aggregate_maximum_bit_rate(
                &apn_ambr, 0, apn_ambr_encoded_buffersize3,
                sizeof(apn_ambr_encoded_buffersize3)),
            sizeof(apn_ambr_encoded_buffersize3));
  ASSERT_EQ((memcmp(apn_ambr_encoded_buffersize3, apn_ambr_buffersize3,
                    sizeof(apn_ambr_buffersize3))),
            0);

  uint8_t apn_ambr_buffersize4[] = {0x06, 0x40, 0x40, 0x00, 0x00, 0x00, 0x00};
  ASSERT_EQ(
      decode_apn_aggregate_maximum_bit_rate(&apn_ambr, 0, apn_ambr_buffersize4,
                                            sizeof(apn_ambr_buffersize4)),
      sizeof(apn_ambr_buffersize4));
  ASSERT_EQ((apn_ambr.apnambrfordownlink), 64);
  ASSERT_EQ((apn_ambr.apnambrforuplink), 64);
  ASSERT_EQ((apn_ambr.apnambrfordownlink_extended), 0);
  ASSERT_EQ((apn_ambr.apnambrforuplink_extended), 0);
  ASSERT_EQ((apn_ambr.apnambrfordownlink_extended2), 0);
  ASSERT_EQ((apn_ambr.apnambrforuplink_extended2), 0);

  uint8_t apn_ambr_encoded_buffersize4[7] = {0};
  ASSERT_EQ(encode_apn_aggregate_maximum_bit_rate(
                &apn_ambr, 0, apn_ambr_encoded_buffersize4,
                sizeof(apn_ambr_encoded_buffersize4)),
            sizeof(apn_ambr_encoded_buffersize4));
  ASSERT_EQ((memcmp(apn_ambr_encoded_buffersize4, apn_ambr_buffersize4,
                    sizeof(apn_ambr_buffersize4))),
            0);

  uint8_t apn_ambr_buffersize5[] = {0x06, 0x41, 0x41, 0x00, 0x00, 0x00, 0x00};
  ASSERT_EQ(
      decode_apn_aggregate_maximum_bit_rate(&apn_ambr, 0, apn_ambr_buffersize5,
                                            sizeof(apn_ambr_buffersize5)),
      sizeof(apn_ambr_buffersize5));
  ASSERT_EQ((apn_ambr.apnambrfordownlink), 65);
  ASSERT_EQ((apn_ambr.apnambrforuplink), 65);
  ASSERT_EQ((apn_ambr.apnambrfordownlink_extended), 0);
  ASSERT_EQ((apn_ambr.apnambrforuplink_extended), 0);
  ASSERT_EQ((apn_ambr.apnambrfordownlink_extended2), 0);
  ASSERT_EQ((apn_ambr.apnambrforuplink_extended2), 0);

  uint8_t apn_ambr_encoded_buffersize5[7] = {0};
  ASSERT_EQ(encode_apn_aggregate_maximum_bit_rate(
                &apn_ambr, 0, apn_ambr_encoded_buffersize5,
                sizeof(apn_ambr_encoded_buffersize5)),
            sizeof(apn_ambr_encoded_buffersize5));
  ASSERT_EQ((memcmp(apn_ambr_encoded_buffersize5, apn_ambr_buffersize5,
                    sizeof(apn_ambr_buffersize5))),
            0);

  uint8_t apn_ambr_buffersize6[] = {0x06, 0x7F, 0x7F, 0x00, 0x00, 0x00, 0x00};
  ASSERT_EQ(
      decode_apn_aggregate_maximum_bit_rate(&apn_ambr, 0, apn_ambr_buffersize6,
                                            sizeof(apn_ambr_buffersize6)),
      sizeof(apn_ambr_buffersize6));
  ASSERT_EQ((apn_ambr.apnambrfordownlink), 127);
  ASSERT_EQ((apn_ambr.apnambrforuplink), 127);
  ASSERT_EQ((apn_ambr.apnambrfordownlink_extended), 0);
  ASSERT_EQ((apn_ambr.apnambrforuplink_extended), 0);
  ASSERT_EQ((apn_ambr.apnambrfordownlink_extended2), 0);
  ASSERT_EQ((apn_ambr.apnambrforuplink_extended2), 0);

  uint8_t apn_ambr_encoded_buffersize6[7] = {0};
  ASSERT_EQ(encode_apn_aggregate_maximum_bit_rate(
                &apn_ambr, 0, apn_ambr_encoded_buffersize6,
                sizeof(apn_ambr_encoded_buffersize6)),
            sizeof(apn_ambr_encoded_buffersize6));
  ASSERT_EQ((memcmp(apn_ambr_encoded_buffersize6, apn_ambr_buffersize6,
                    sizeof(apn_ambr_buffersize6))),
            0);

  uint8_t apn_ambr_buffersize7[] = {0x06, 0x80, 0x80, 0x00, 0x00, 0x00, 0x00};
  ASSERT_EQ(
      decode_apn_aggregate_maximum_bit_rate(&apn_ambr, 0, apn_ambr_buffersize7,
                                            sizeof(apn_ambr_buffersize7)),
      sizeof(apn_ambr_buffersize7));
  ASSERT_EQ((apn_ambr.apnambrfordownlink), 128);
  ASSERT_EQ((apn_ambr.apnambrforuplink), 128);
  ASSERT_EQ((apn_ambr.apnambrfordownlink_extended), 0);
  ASSERT_EQ((apn_ambr.apnambrforuplink_extended), 0);
  ASSERT_EQ((apn_ambr.apnambrfordownlink_extended2), 0);
  ASSERT_EQ((apn_ambr.apnambrforuplink_extended2), 0);

  uint8_t apn_ambr_encoded_buffersize7[7] = {0};
  ASSERT_EQ(encode_apn_aggregate_maximum_bit_rate(
                &apn_ambr, 0, apn_ambr_encoded_buffersize7,
                sizeof(apn_ambr_encoded_buffersize7)),
            sizeof(apn_ambr_encoded_buffersize7));
  ASSERT_EQ((memcmp(apn_ambr_encoded_buffersize7, apn_ambr_buffersize7,
                    sizeof(apn_ambr_buffersize7))),
            0);

  uint8_t apn_ambr_buffersize8[] = {0x06, 0x89, 0x89, 0x00, 0x00, 0x00, 0x00};
  ASSERT_EQ(
      decode_apn_aggregate_maximum_bit_rate(&apn_ambr, 0, apn_ambr_buffersize8,
                                            sizeof(apn_ambr_buffersize8)),
      sizeof(apn_ambr_buffersize8));
  ASSERT_EQ((apn_ambr.apnambrfordownlink), 137);
  ASSERT_EQ((apn_ambr.apnambrforuplink), 137);
  ASSERT_EQ((apn_ambr.apnambrfordownlink_extended), 0);
  ASSERT_EQ((apn_ambr.apnambrforuplink_extended), 0);
  ASSERT_EQ((apn_ambr.apnambrfordownlink_extended2), 0);
  ASSERT_EQ((apn_ambr.apnambrforuplink_extended2), 0);

  uint8_t apn_ambr_encoded_buffersize8[7] = {0};
  ASSERT_EQ(encode_apn_aggregate_maximum_bit_rate(
                &apn_ambr, 0, apn_ambr_encoded_buffersize8,
                sizeof(apn_ambr_encoded_buffersize8)),
            sizeof(apn_ambr_encoded_buffersize8));
  ASSERT_EQ((memcmp(apn_ambr_encoded_buffersize8, apn_ambr_buffersize8,
                    sizeof(apn_ambr_buffersize8))),
            0);

  uint8_t apn_ambr_buffersize9[] = {0x06, 0xFE, 0xFE, 0x00, 0x00, 0x00, 0x00};
  ASSERT_EQ(
      decode_apn_aggregate_maximum_bit_rate(&apn_ambr, 0, apn_ambr_buffersize9,
                                            sizeof(apn_ambr_buffersize9)),
      sizeof(apn_ambr_buffersize9));
  ASSERT_EQ((apn_ambr.apnambrfordownlink), 254);
  ASSERT_EQ((apn_ambr.apnambrforuplink), 254);
  ASSERT_EQ((apn_ambr.apnambrfordownlink_extended), 0);
  ASSERT_EQ((apn_ambr.apnambrforuplink_extended), 0);
  ASSERT_EQ((apn_ambr.apnambrfordownlink_extended2), 0);
  ASSERT_EQ((apn_ambr.apnambrforuplink_extended2), 0);

  uint8_t apn_ambr_encoded_buffersize9[7] = {0};
  ASSERT_EQ(encode_apn_aggregate_maximum_bit_rate(
                &apn_ambr, 0, apn_ambr_encoded_buffersize9,
                sizeof(apn_ambr_encoded_buffersize9)),
            sizeof(apn_ambr_encoded_buffersize9));
  ASSERT_EQ((memcmp(apn_ambr_encoded_buffersize9, apn_ambr_buffersize9,
                    sizeof(apn_ambr_buffersize9))),
            0);

  uint8_t apn_ambr_buffersize10[] = {0x06, 0xFF, 0xFF, 0x01, 0x01, 0x00, 0x00};
  ASSERT_EQ(
      decode_apn_aggregate_maximum_bit_rate(&apn_ambr, 0, apn_ambr_buffersize10,
                                            sizeof(apn_ambr_buffersize10)),
      sizeof(apn_ambr_buffersize10));
  ASSERT_EQ((apn_ambr.apnambrfordownlink), 255);
  ASSERT_EQ((apn_ambr.apnambrforuplink), 255);
  ASSERT_EQ((apn_ambr.apnambrfordownlink_extended), 1);
  ASSERT_EQ((apn_ambr.apnambrforuplink_extended), 1);
  ASSERT_EQ((apn_ambr.apnambrfordownlink_extended2), 0);
  ASSERT_EQ((apn_ambr.apnambrforuplink_extended2), 0);

  uint8_t apn_ambr_encoded_buffersize10[7] = {0};
  ASSERT_EQ(encode_apn_aggregate_maximum_bit_rate(
                &apn_ambr, 0, apn_ambr_encoded_buffersize10,
                sizeof(apn_ambr_encoded_buffersize10)),
            sizeof(apn_ambr_encoded_buffersize10));
  ASSERT_EQ((memcmp(apn_ambr_encoded_buffersize10, apn_ambr_buffersize10,
                    sizeof(apn_ambr_buffersize10))),
            0);

  uint8_t apn_ambr_buffersize11[] = {0x06, 0xFF, 0xFF, 0x04, 0x04, 0x00, 0x00};
  ASSERT_EQ(
      decode_apn_aggregate_maximum_bit_rate(&apn_ambr, 0, apn_ambr_buffersize11,
                                            sizeof(apn_ambr_buffersize11)),
      sizeof(apn_ambr_buffersize11));
  ASSERT_EQ((apn_ambr.apnambrfordownlink), 255);
  ASSERT_EQ((apn_ambr.apnambrforuplink), 255);
  ASSERT_EQ((apn_ambr.apnambrfordownlink_extended), 4);
  ASSERT_EQ((apn_ambr.apnambrforuplink_extended), 4);
  ASSERT_EQ((apn_ambr.apnambrfordownlink_extended2), 0);
  ASSERT_EQ((apn_ambr.apnambrforuplink_extended2), 0);

  uint8_t apn_ambr_encoded_buffersize11[7] = {0};
  ASSERT_EQ(encode_apn_aggregate_maximum_bit_rate(
                &apn_ambr, 0, apn_ambr_encoded_buffersize11,
                sizeof(apn_ambr_encoded_buffersize11)),
            sizeof(apn_ambr_encoded_buffersize11));
  ASSERT_EQ((memcmp(apn_ambr_encoded_buffersize11, apn_ambr_buffersize11,
                    sizeof(apn_ambr_buffersize11))),
            0);

  uint8_t apn_ambr_buffersize12[] = {0x06, 0xFF, 0xFF, 0x4A, 0x4A, 0x00, 0x00};
  ASSERT_EQ(
      decode_apn_aggregate_maximum_bit_rate(&apn_ambr, 0, apn_ambr_buffersize12,
                                            sizeof(apn_ambr_buffersize12)),
      sizeof(apn_ambr_buffersize12));
  ASSERT_EQ((apn_ambr.apnambrfordownlink), 255);
  ASSERT_EQ((apn_ambr.apnambrforuplink), 255);
  ASSERT_EQ((apn_ambr.apnambrfordownlink_extended), 74);
  ASSERT_EQ((apn_ambr.apnambrforuplink_extended), 74);
  ASSERT_EQ((apn_ambr.apnambrfordownlink_extended2), 0);
  ASSERT_EQ((apn_ambr.apnambrforuplink_extended2), 0);

  uint8_t apn_ambr_encoded_buffersize12[7] = {0};
  ASSERT_EQ(encode_apn_aggregate_maximum_bit_rate(
                &apn_ambr, 0, apn_ambr_encoded_buffersize12,
                sizeof(apn_ambr_encoded_buffersize12)),
            sizeof(apn_ambr_encoded_buffersize12));
  ASSERT_EQ((memcmp(apn_ambr_encoded_buffersize12, apn_ambr_buffersize12,
                    sizeof(apn_ambr_buffersize12))),
            0);

  uint8_t apn_ambr_buffersize13[] = {0x06, 0xFF, 0xFF, 0x4B, 0x4B, 0x00, 0x00};
  ASSERT_EQ(
      decode_apn_aggregate_maximum_bit_rate(&apn_ambr, 0, apn_ambr_buffersize13,
                                            sizeof(apn_ambr_buffersize13)),
      sizeof(apn_ambr_buffersize13));
  ASSERT_EQ((apn_ambr.apnambrfordownlink), 255);
  ASSERT_EQ((apn_ambr.apnambrforuplink), 255);
  ASSERT_EQ((apn_ambr.apnambrfordownlink_extended), 75);
  ASSERT_EQ((apn_ambr.apnambrforuplink_extended), 75);
  ASSERT_EQ((apn_ambr.apnambrfordownlink_extended2), 0);
  ASSERT_EQ((apn_ambr.apnambrforuplink_extended2), 0);

  uint8_t apn_ambr_encoded_buffersize13[7] = {0};
  ASSERT_EQ(encode_apn_aggregate_maximum_bit_rate(
                &apn_ambr, 0, apn_ambr_encoded_buffersize13,
                sizeof(apn_ambr_encoded_buffersize13)),
            sizeof(apn_ambr_encoded_buffersize13));
  ASSERT_EQ((memcmp(apn_ambr_encoded_buffersize13, apn_ambr_buffersize13,
                    sizeof(apn_ambr_buffersize13))),
            0);

  uint8_t apn_ambr_buffersize14[] = {0x06, 0xFF, 0xFF, 0xBA, 0xBA, 0x00, 0x00};
  ASSERT_EQ(
      decode_apn_aggregate_maximum_bit_rate(&apn_ambr, 0, apn_ambr_buffersize14,
                                            sizeof(apn_ambr_buffersize14)),
      sizeof(apn_ambr_buffersize14));
  ASSERT_EQ((apn_ambr.apnambrfordownlink), 255);
  ASSERT_EQ((apn_ambr.apnambrforuplink), 255);
  ASSERT_EQ((apn_ambr.apnambrfordownlink_extended), 186);
  ASSERT_EQ((apn_ambr.apnambrforuplink_extended), 186);
  ASSERT_EQ((apn_ambr.apnambrfordownlink_extended2), 0);
  ASSERT_EQ((apn_ambr.apnambrforuplink_extended2), 0);

  uint8_t apn_ambr_encoded_buffersize14[7] = {0};
  ASSERT_EQ(encode_apn_aggregate_maximum_bit_rate(
                &apn_ambr, 0, apn_ambr_encoded_buffersize14,
                sizeof(apn_ambr_encoded_buffersize14)),
            sizeof(apn_ambr_encoded_buffersize14));
  ASSERT_EQ((memcmp(apn_ambr_encoded_buffersize14, apn_ambr_buffersize14,
                    sizeof(apn_ambr_buffersize14))),
            0);

  uint8_t apn_ambr_buffersize15[] = {0x06, 0xFF, 0xFF, 0xBB, 0xBB, 0x00, 0x00};
  ASSERT_EQ(
      decode_apn_aggregate_maximum_bit_rate(&apn_ambr, 0, apn_ambr_buffersize15,
                                            sizeof(apn_ambr_buffersize15)),
      sizeof(apn_ambr_buffersize15));
  ASSERT_EQ((apn_ambr.apnambrfordownlink), 255);
  ASSERT_EQ((apn_ambr.apnambrforuplink), 255);
  ASSERT_EQ((apn_ambr.apnambrfordownlink_extended), 187);
  ASSERT_EQ((apn_ambr.apnambrforuplink_extended), 187);
  ASSERT_EQ((apn_ambr.apnambrfordownlink_extended2), 0);
  ASSERT_EQ((apn_ambr.apnambrforuplink_extended2), 0);

  uint8_t apn_ambr_encoded_buffersize15[7] = {0};
  ASSERT_EQ(encode_apn_aggregate_maximum_bit_rate(
                &apn_ambr, 0, apn_ambr_encoded_buffersize15,
                sizeof(apn_ambr_encoded_buffersize15)),
            sizeof(apn_ambr_encoded_buffersize15));
  ASSERT_EQ((memcmp(apn_ambr_encoded_buffersize15, apn_ambr_buffersize15,
                    sizeof(apn_ambr_buffersize15))),
            0);

  uint8_t apn_ambr_buffersize16[] = {0x06, 0xFF, 0xFF, 0xDE, 0xDE, 0x01, 0x01};
  ASSERT_EQ(
      decode_apn_aggregate_maximum_bit_rate(&apn_ambr, 0, apn_ambr_buffersize16,
                                            sizeof(apn_ambr_buffersize16)),
      sizeof(apn_ambr_buffersize16));
  ASSERT_EQ((apn_ambr.apnambrfordownlink), 255);
  ASSERT_EQ((apn_ambr.apnambrforuplink), 255);
  ASSERT_EQ((apn_ambr.apnambrfordownlink_extended), 222);
  ASSERT_EQ((apn_ambr.apnambrforuplink_extended), 222);
  ASSERT_EQ((apn_ambr.apnambrfordownlink_extended2), 1);
  ASSERT_EQ((apn_ambr.apnambrforuplink_extended2), 1);

  uint8_t apn_ambr_encoded_buffersize16[7] = {0};
  ASSERT_EQ(encode_apn_aggregate_maximum_bit_rate(
                &apn_ambr, 0, apn_ambr_encoded_buffersize16,
                sizeof(apn_ambr_encoded_buffersize16)),
            sizeof(apn_ambr_encoded_buffersize16));
  ASSERT_EQ((memcmp(apn_ambr_encoded_buffersize16, apn_ambr_buffersize16,
                    sizeof(apn_ambr_buffersize16))),
            0);

  uint8_t apn_ambr_buffersize17[] = {0x06, 0xFF, 0xFF, 0xFA, 0xFA, 0x00, 0x00};
  ASSERT_EQ(
      decode_apn_aggregate_maximum_bit_rate(&apn_ambr, 0, apn_ambr_buffersize17,
                                            sizeof(apn_ambr_buffersize17)),
      sizeof(apn_ambr_buffersize17));
  ASSERT_EQ((apn_ambr.apnambrfordownlink), 255);
  ASSERT_EQ((apn_ambr.apnambrforuplink), 255);
  ASSERT_EQ((apn_ambr.apnambrfordownlink_extended), 250);
  ASSERT_EQ((apn_ambr.apnambrforuplink_extended), 250);
  ASSERT_EQ((apn_ambr.apnambrfordownlink_extended2), 0);
  ASSERT_EQ((apn_ambr.apnambrforuplink_extended2), 0);

  uint8_t apn_ambr_encoded_buffersize17[7] = {0};
  ASSERT_EQ(encode_apn_aggregate_maximum_bit_rate(
                &apn_ambr, 0, apn_ambr_encoded_buffersize17,
                sizeof(apn_ambr_encoded_buffersize17)),
            sizeof(apn_ambr_encoded_buffersize17));
  ASSERT_EQ((memcmp(apn_ambr_encoded_buffersize17, apn_ambr_buffersize17,
                    sizeof(apn_ambr_buffersize17))),
            0);

  uint8_t apn_ambr_buffersize18[] = {0x06, 0xFF, 0xFF, 0xFA, 0xFA, 0x00, 0x00};
  ASSERT_EQ(
      decode_apn_aggregate_maximum_bit_rate(&apn_ambr, 0, apn_ambr_buffersize18,
                                            sizeof(apn_ambr_buffersize18)),
      sizeof(apn_ambr_buffersize18));
  ASSERT_EQ((apn_ambr.apnambrfordownlink), 255);
  ASSERT_EQ((apn_ambr.apnambrforuplink), 255);
  ASSERT_EQ((apn_ambr.apnambrfordownlink_extended), 250);
  ASSERT_EQ((apn_ambr.apnambrforuplink_extended), 250);
  ASSERT_EQ((apn_ambr.apnambrfordownlink_extended2), 0);
  ASSERT_EQ((apn_ambr.apnambrforuplink_extended2), 0);

  uint8_t apn_ambr_encoded_buffersize18[7] = {0};
  ASSERT_EQ(encode_apn_aggregate_maximum_bit_rate(
                &apn_ambr, 0, apn_ambr_encoded_buffersize18,
                sizeof(apn_ambr_encoded_buffersize18)),
            sizeof(apn_ambr_encoded_buffersize18));
  ASSERT_EQ((memcmp(apn_ambr_encoded_buffersize18, apn_ambr_buffersize18,
                    sizeof(apn_ambr_buffersize18))),
            0);

  uint8_t apn_ambr_buffersize19[] = {0x06, 0xFF, 0xFF, 0x66, 0x66, 0x01, 0x01};
  ASSERT_EQ(
      decode_apn_aggregate_maximum_bit_rate(&apn_ambr, 0, apn_ambr_buffersize19,
                                            sizeof(apn_ambr_buffersize19)),
      sizeof(apn_ambr_buffersize19));
  ASSERT_EQ((apn_ambr.apnambrfordownlink), 255);
  ASSERT_EQ((apn_ambr.apnambrforuplink), 255);
  ASSERT_EQ((apn_ambr.apnambrfordownlink_extended), 102);
  ASSERT_EQ((apn_ambr.apnambrforuplink_extended), 102);
  ASSERT_EQ((apn_ambr.apnambrfordownlink_extended2), 1);
  ASSERT_EQ((apn_ambr.apnambrforuplink_extended2), 1);

  uint8_t apn_ambr_encoded_buffersize19[7] = {0};
  ASSERT_EQ(encode_apn_aggregate_maximum_bit_rate(
                &apn_ambr, 0, apn_ambr_encoded_buffersize19,
                sizeof(apn_ambr_encoded_buffersize19)),
            sizeof(apn_ambr_encoded_buffersize19));
  ASSERT_EQ((memcmp(apn_ambr_encoded_buffersize19, apn_ambr_buffersize19,
                    sizeof(apn_ambr_buffersize19))),
            0);

  uint8_t apn_ambr_buffersize20[] = {0x06, 0x7B, 0x7B, 0x70, 0x70, 0x01, 0x01};
  ASSERT_EQ(
      decode_apn_aggregate_maximum_bit_rate(&apn_ambr, 0, apn_ambr_buffersize20,
                                            sizeof(apn_ambr_buffersize20)),
      sizeof(apn_ambr_buffersize20));
  ASSERT_EQ((apn_ambr.apnambrfordownlink), 123);
  ASSERT_EQ((apn_ambr.apnambrforuplink), 123);
  ASSERT_EQ((apn_ambr.apnambrfordownlink_extended), 112);
  ASSERT_EQ((apn_ambr.apnambrforuplink_extended), 112);
  ASSERT_EQ((apn_ambr.apnambrfordownlink_extended2), 1);
  ASSERT_EQ((apn_ambr.apnambrforuplink_extended2), 1);

  uint8_t apn_ambr_encoded_buffersize20[7] = {0};
  ASSERT_EQ(encode_apn_aggregate_maximum_bit_rate(
                &apn_ambr, 0, apn_ambr_encoded_buffersize20,
                sizeof(apn_ambr_encoded_buffersize20)),
            sizeof(apn_ambr_encoded_buffersize20));
  ASSERT_EQ((memcmp(apn_ambr_encoded_buffersize20, apn_ambr_buffersize20,
                    sizeof(apn_ambr_buffersize20))),
            0);

  uint8_t apn_ambr_buffersize21[] = {0x06, 0x05, 0x05, 0x0E, 0x0E, 0x02, 0x02};
  ASSERT_EQ(
      decode_apn_aggregate_maximum_bit_rate(&apn_ambr, 0, apn_ambr_buffersize21,
                                            sizeof(apn_ambr_buffersize21)),
      sizeof(apn_ambr_buffersize21));
  ASSERT_EQ((apn_ambr.apnambrfordownlink), 5);
  ASSERT_EQ((apn_ambr.apnambrforuplink), 5);
  ASSERT_EQ((apn_ambr.apnambrfordownlink_extended), 14);
  ASSERT_EQ((apn_ambr.apnambrforuplink_extended), 14);
  ASSERT_EQ((apn_ambr.apnambrfordownlink_extended2), 2);
  ASSERT_EQ((apn_ambr.apnambrforuplink_extended2), 2);

  uint8_t apn_ambr_encoded_buffersize21[7] = {0};
  ASSERT_EQ(encode_apn_aggregate_maximum_bit_rate(
                &apn_ambr, 0, apn_ambr_encoded_buffersize21,
                sizeof(apn_ambr_encoded_buffersize21)),
            sizeof(apn_ambr_encoded_buffersize21));
  ASSERT_EQ((memcmp(apn_ambr_encoded_buffersize21, apn_ambr_buffersize21,
                    sizeof(apn_ambr_buffersize21))),
            0);

  uint8_t apn_ambr_buffersize22[] = {0x06, 0xFF, 0xFF, 0x00, 0x00, 0xFF, 0xFF};
  ASSERT_EQ(
      decode_apn_aggregate_maximum_bit_rate(&apn_ambr, 0, apn_ambr_buffersize22,
                                            sizeof(apn_ambr_buffersize22)),
      sizeof(apn_ambr_buffersize22));
  ASSERT_EQ((apn_ambr.apnambrfordownlink), 255);
  ASSERT_EQ((apn_ambr.apnambrforuplink), 255);
  ASSERT_EQ((apn_ambr.apnambrfordownlink_extended), 0);
  ASSERT_EQ((apn_ambr.apnambrforuplink_extended), 0);
  ASSERT_EQ((apn_ambr.apnambrfordownlink_extended2), 255);
  ASSERT_EQ((apn_ambr.apnambrforuplink_extended2), 255);

  uint8_t apn_ambr_encoded_buffersize22[7] = {0};
  ASSERT_EQ(encode_apn_aggregate_maximum_bit_rate(
                &apn_ambr, 0, apn_ambr_encoded_buffersize22,
                sizeof(apn_ambr_encoded_buffersize22)),
            sizeof(apn_ambr_encoded_buffersize22));
  ASSERT_EQ((memcmp(apn_ambr_encoded_buffersize22, apn_ambr_buffersize22,
                    sizeof(apn_ambr_buffersize22))),
            0);
}

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
#include "AttachRequest.h"
#include "EpsQualityOfService.h"
#include "dynamic_memory_check.h"
#include "log.h"
}

class NASEncodeDecodeTest : public ::testing::Test {
  virtual void SetUp() {}
  virtual void TearDown() {}
};

TEST_F(NASEncodeDecodeTest, TestDecodeAttachRequest1) {
  //   Combined attach, NAS message generated from s1ap tester
  uint8_t buffer[] = {0x72, 0x08, 0x09, 0x10, 0x10, 0x00, 0x00, 0x00,
                      0x00, 0x10, 0x02, 0xe0, 0xe0, 0x00, 0x04, 0x02,
                      0x01, 0xd0, 0x11, 0x40, 0x08, 0x04, 0x02, 0x60,
                      0x04, 0x00, 0x02, 0x1c, 0x00};
  uint32_t len     = 29;
  attach_request_msg attach_request;

  int rc = decode_attach_request(&attach_request, buffer, len);
  ASSERT_EQ(rc, len);
  ASSERT_EQ(attach_request.epsattachtype, 2);
  ASSERT_EQ(attach_request.naskeysetidentifier.naskeysetidentifier, 7);
  ASSERT_EQ(attach_request.naskeysetidentifier.tsc, 0);

  bdestroy_wrapper(&attach_request.esmmessagecontainer);
}

TEST_F(NASEncodeDecodeTest, TestDecodeAttachRequest2) {
  //   Combined attach, NAS message generated from Pixel 4
  uint8_t buffer[] = {
      0x72, 0x08, 0x39, 0x51, 0x10, 0x00, 0x30, 0x09, 0x01, 0x07, 0x07, 0xf0,
      0x70, 0xc0, 0x40, 0x19, 0x00, 0x80, 0x00, 0x34, 0x02, 0x0c, 0xd0, 0x11,
      0xd1, 0x27, 0x2d, 0x80, 0x80, 0x21, 0x10, 0x01, 0x00, 0x00, 0x10, 0x81,
      0x06, 0x00, 0x00, 0x00, 0x00, 0x83, 0x06, 0x00, 0x00, 0x00, 0x00, 0x00,
      0x0d, 0x00, 0x00, 0x0a, 0x00, 0x00, 0x05, 0x00, 0x00, 0x10, 0x00, 0x00,
      0x11, 0x00, 0x00, 0x1a, 0x01, 0x01, 0x00, 0x23, 0x00, 0x00, 0x24, 0x00,
      0x5c, 0x0a, 0x01, 0x31, 0x04, 0x65, 0xe0, 0x3e, 0x00, 0x90, 0x11, 0x03,
      0x57, 0x58, 0xa6, 0x20, 0x0d, 0x60, 0x14, 0x04, 0xef, 0x65, 0x23, 0x3b,
      0x88, 0x00, 0x92, 0xf2, 0x00, 0x00, 0x40, 0x08, 0x04, 0x02, 0x60, 0x04,
      0x00, 0x02, 0x1f, 0x00, 0x5d, 0x01, 0x03, 0xc1};

  uint32_t len = 116;
  attach_request_msg attach_request;

  int rc = decode_attach_request(&attach_request, buffer, len);
  ASSERT_EQ(rc, len);
  ASSERT_EQ(attach_request.epsattachtype, 2);
  ASSERT_EQ(attach_request.naskeysetidentifier.naskeysetidentifier, 7);
  ASSERT_EQ(attach_request.naskeysetidentifier.tsc, 0);

  bdestroy_wrapper(&attach_request.esmmessagecontainer);
}

TEST_F(NASEncodeDecodeTest, TestDecodeEncodeEPSQoS) {
  EpsQualityOfService eps_qos = {0};

  uint8_t eps_qos_buffersize1[] = {0x09, 0x09, 0x1c, 0x1c, 0x1c,
                                   0x1c, 0x00, 0x00, 0x00, 0x00};
  ASSERT_EQ(
      decode_eps_quality_of_service(
          &eps_qos, 0, eps_qos_buffersize1, sizeof(eps_qos_buffersize1)),
      sizeof(eps_qos_buffersize1));
  ASSERT_EQ((eps_qos.bitRates.maxBitRateForUL), 28);
  ASSERT_EQ((eps_qos.bitRates.maxBitRateForDL), 28);
  ASSERT_EQ((eps_qos.bitRatesExt.maxBitRateForUL), 0);
  ASSERT_EQ((eps_qos.bitRatesExt.maxBitRateForDL), 0);
  ASSERT_EQ((eps_qos.bitRatesExt2.maxBitRateForUL), 0);
  ASSERT_EQ((eps_qos.bitRatesExt2.maxBitRateForDL), 0);
  uint8_t eps_qos_buffersize_encoded1[10] = {0};
  ASSERT_EQ(
      encode_eps_quality_of_service(
          &eps_qos, 0, eps_qos_buffersize1, sizeof(eps_qos_buffersize1)),
      sizeof(eps_qos_buffersize_encoded1));

  ASSERT_EQ(sizeof(eps_qos_buffersize1), sizeof(eps_qos_buffersize_encoded1));
  memcmp(
      eps_qos_buffersize_encoded1, eps_qos_buffersize1,
      sizeof(eps_qos_buffersize1));

  uint8_t eps_qos_buffersize2[] = {0x09, 0x09, 0x3F, 0x3F, 0x3F,
                                   0x3F, 0x00, 0x00, 0x00, 0x00};

  ASSERT_EQ(
      decode_eps_quality_of_service(
          &eps_qos, 0, eps_qos_buffersize2, sizeof(eps_qos_buffersize2)),
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
      encode_eps_quality_of_service(
          &eps_qos, 0, eps_qos_buffersize2, sizeof(eps_qos_buffersize2)),
      sizeof(eps_qos_buffersize_encoded2));
  ASSERT_EQ(sizeof(eps_qos_buffersize2), sizeof(eps_qos_buffersize_encoded2));
  memcmp(
      eps_qos_buffersize_encoded2, eps_qos_buffersize2,
      sizeof(eps_qos_buffersize2));

  uint8_t eps_qos_buffersize3[] = {0x09, 0x09, 0x41, 0x41, 0x41,
                                   0x41, 0x00, 0x00, 0x00, 0x00};

  ASSERT_EQ(
      decode_eps_quality_of_service(
          &eps_qos, 0, eps_qos_buffersize3, sizeof(eps_qos_buffersize3)),
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
      encode_eps_quality_of_service(
          &eps_qos, 0, eps_qos_buffersize3, sizeof(eps_qos_buffersize3)),
      sizeof(eps_qos_buffersize_encoded3));
  ASSERT_EQ(sizeof(eps_qos_buffersize3), sizeof(eps_qos_buffersize_encoded3));
  memcmp(
      eps_qos_buffersize_encoded3, eps_qos_buffersize3,
      sizeof(eps_qos_buffersize3));

  uint8_t eps_qos_buffersize4[] = {0x09, 0x09, 0x7F, 0x7F, 0x7F,
                                   0x7F, 0x00, 0x00, 0x00, 0x00};

  ASSERT_EQ(
      decode_eps_quality_of_service(
          &eps_qos, 0, eps_qos_buffersize4, sizeof(eps_qos_buffersize4)),
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
      encode_eps_quality_of_service(
          &eps_qos, 0, eps_qos_buffersize4, sizeof(eps_qos_buffersize4)),
      sizeof(eps_qos_buffersize_encoded4));
  ASSERT_EQ(sizeof(eps_qos_buffersize4), sizeof(eps_qos_buffersize_encoded4));
  memcmp(
      eps_qos_buffersize_encoded4, eps_qos_buffersize4,
      sizeof(eps_qos_buffersize4));

  uint8_t eps_qos_buffersize5[] = {0x09, 0x09, 0x81, 0x81, 0x81,
                                   0x81, 0x00, 0x00, 0x00, 0x00};

  ASSERT_EQ(
      decode_eps_quality_of_service(
          &eps_qos, 0, eps_qos_buffersize5, sizeof(eps_qos_buffersize5)),
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
      encode_eps_quality_of_service(
          &eps_qos, 0, eps_qos_buffersize5, sizeof(eps_qos_buffersize5)),
      sizeof(eps_qos_buffersize_encoded5));
  ASSERT_EQ(sizeof(eps_qos_buffersize5), sizeof(eps_qos_buffersize_encoded5));
  memcmp(
      eps_qos_buffersize_encoded5, eps_qos_buffersize5,
      sizeof(eps_qos_buffersize5));

  uint8_t eps_qos_buffersize6[] = {0x09, 0x09, 0x82, 0x82, 0x82,
                                   0x82, 0x00, 0x00, 0x00, 0x00};

  ASSERT_EQ(
      decode_eps_quality_of_service(
          &eps_qos, 0, eps_qos_buffersize6, sizeof(eps_qos_buffersize6)),
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
      encode_eps_quality_of_service(
          &eps_qos, 0, eps_qos_buffersize6, sizeof(eps_qos_buffersize6)),
      sizeof(eps_qos_buffersize_encoded6));
  ASSERT_EQ(sizeof(eps_qos_buffersize6), sizeof(eps_qos_buffersize_encoded6));
  memcmp(
      eps_qos_buffersize_encoded6, eps_qos_buffersize6,
      sizeof(eps_qos_buffersize6));

  uint8_t eps_qos_buffersize7[] = {0x09, 0x09, 0x82, 0x82, 0x82,
                                   0x82, 0x00, 0x00, 0x00, 0x00};

  ASSERT_EQ(
      decode_eps_quality_of_service(
          &eps_qos, 0, eps_qos_buffersize7, sizeof(eps_qos_buffersize7)),
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
      encode_eps_quality_of_service(
          &eps_qos, 0, eps_qos_buffersize7, sizeof(eps_qos_buffersize7)),
      sizeof(eps_qos_buffersize_encoded7));
  ASSERT_EQ(sizeof(eps_qos_buffersize7), sizeof(eps_qos_buffersize_encoded7));
  memcmp(
      eps_qos_buffersize_encoded7, eps_qos_buffersize7,
      sizeof(eps_qos_buffersize7));

  uint8_t eps_qos_buffersize8[] = {0x09, 0x09, 0xFE, 0xFE, 0xFE,
                                   0xFE, 0x00, 0x00, 0x00, 0x00};

  ASSERT_EQ(
      decode_eps_quality_of_service(
          &eps_qos, 0, eps_qos_buffersize8, sizeof(eps_qos_buffersize8)),
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
      encode_eps_quality_of_service(
          &eps_qos, 0, eps_qos_buffersize8, sizeof(eps_qos_buffersize8)),
      sizeof(eps_qos_buffersize_encoded8));
  ASSERT_EQ(sizeof(eps_qos_buffersize8), sizeof(eps_qos_buffersize_encoded8));
  memcmp(
      eps_qos_buffersize_encoded8, eps_qos_buffersize8,
      sizeof(eps_qos_buffersize8));

  uint8_t eps_qos_buffersize_Ext1[] = {0x09, 0x09, 0xFE, 0xFE, 0xFE,
                                       0xFE, 0x02, 0x02, 0x02, 0x02};

  ASSERT_EQ(
      decode_eps_quality_of_service(
          &eps_qos, 0, eps_qos_buffersize_Ext1,
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
  ASSERT_EQ(
      encode_eps_quality_of_service(
          &eps_qos, 0, eps_qos_buffersize_Ext1,
          sizeof(eps_qos_buffersize_Ext1)),
      sizeof(eps_qos_buffersize_encoded_Ext1));
  ASSERT_EQ(
      sizeof(eps_qos_buffersize_Ext1), sizeof(eps_qos_buffersize_encoded_Ext1));
  memcmp(
      eps_qos_buffersize_encoded_Ext1, eps_qos_buffersize_Ext1,
      sizeof(eps_qos_buffersize_Ext1));

  uint8_t eps_qos_buffersize_Ext2[] = {0x09, 0x09, 0xFE, 0xFE, 0xFE,
                                       0xFE, 0x4A, 0x4A, 0x4A, 0x4A};

  ASSERT_EQ(
      decode_eps_quality_of_service(
          &eps_qos, 0, eps_qos_buffersize_Ext2,
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
  ASSERT_EQ(
      encode_eps_quality_of_service(
          &eps_qos, 0, eps_qos_buffersize_Ext2,
          sizeof(eps_qos_buffersize_Ext2)),
      sizeof(eps_qos_buffersize_encoded_Ext2));
  ASSERT_EQ(
      sizeof(eps_qos_buffersize_Ext2), sizeof(eps_qos_buffersize_encoded_Ext2));
  memcmp(
      eps_qos_buffersize_encoded_Ext2, eps_qos_buffersize_Ext2,
      sizeof(eps_qos_buffersize_Ext2));

  uint8_t eps_qos_buffersize_Ext3[] = {0x09, 0x09, 0xFE, 0xFE, 0xFE,
                                       0xFE, 0x4B, 0x4B, 0x4B, 0x4B};

  ASSERT_EQ(
      decode_eps_quality_of_service(
          &eps_qos, 0, eps_qos_buffersize_Ext3,
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
  ASSERT_EQ(
      encode_eps_quality_of_service(
          &eps_qos, 0, eps_qos_buffersize_Ext3,
          sizeof(eps_qos_buffersize_Ext3)),
      sizeof(eps_qos_buffersize_encoded_Ext3));
  ASSERT_EQ(
      sizeof(eps_qos_buffersize_Ext3), sizeof(eps_qos_buffersize_encoded_Ext3));
  memcmp(
      eps_qos_buffersize_encoded_Ext3, eps_qos_buffersize_Ext3,
      sizeof(eps_qos_buffersize_Ext3));

  uint8_t eps_qos_buffersize_Ext4[] = {0x09, 0x09, 0xFE, 0xFE, 0xFE,
                                       0xFE, 0x4E, 0x4E, 0x4E, 0x4E};

  ASSERT_EQ(
      decode_eps_quality_of_service(
          &eps_qos, 0, eps_qos_buffersize_Ext4,
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
  ASSERT_EQ(
      encode_eps_quality_of_service(
          &eps_qos, 0, eps_qos_buffersize_Ext4,
          sizeof(eps_qos_buffersize_Ext4)),
      sizeof(eps_qos_buffersize_encoded_Ext4));
  ASSERT_EQ(
      sizeof(eps_qos_buffersize_Ext4), sizeof(eps_qos_buffersize_encoded_Ext4));
  memcmp(
      eps_qos_buffersize_encoded_Ext4, eps_qos_buffersize_Ext4,
      sizeof(eps_qos_buffersize_Ext4));

  uint8_t eps_qos_buffersize_Ext5[] = {0x09, 0x09, 0xFE, 0xFE, 0xFE,
                                       0xFE, 0xBA, 0xBA, 0xBA, 0xBA};

  ASSERT_EQ(
      decode_eps_quality_of_service(
          &eps_qos, 0, eps_qos_buffersize_Ext5,
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
  ASSERT_EQ(
      encode_eps_quality_of_service(
          &eps_qos, 0, eps_qos_buffersize_Ext5,
          sizeof(eps_qos_buffersize_Ext5)),
      sizeof(eps_qos_buffersize_encoded_Ext5));
  ASSERT_EQ(
      sizeof(eps_qos_buffersize_Ext5), sizeof(eps_qos_buffersize_encoded_Ext5));
  memcmp(
      eps_qos_buffersize_encoded_Ext5, eps_qos_buffersize_Ext5,
      sizeof(eps_qos_buffersize_Ext5));

  uint8_t eps_qos_buffersize_Ext6[] = {0x09, 0x09, 0xFE, 0xFE, 0xFE,
                                       0xFE, 0xBB, 0xBB, 0xBB, 0xBB};

  ASSERT_EQ(
      decode_eps_quality_of_service(
          &eps_qos, 0, eps_qos_buffersize_Ext6,
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
  ASSERT_EQ(
      encode_eps_quality_of_service(
          &eps_qos, 0, eps_qos_buffersize_Ext6,
          sizeof(eps_qos_buffersize_Ext6)),
      sizeof(eps_qos_buffersize_encoded_Ext6));
  ASSERT_EQ(
      sizeof(eps_qos_buffersize_Ext6), sizeof(eps_qos_buffersize_encoded_Ext6));
  memcmp(
      eps_qos_buffersize_encoded_Ext6, eps_qos_buffersize_Ext6,
      sizeof(eps_qos_buffersize_Ext6));

  uint8_t eps_qos_buffersize_Ext7[] = {0x09, 0x09, 0xFE, 0xFE, 0xFE,
                                       0xFE, 0xC5, 0xC5, 0xC5, 0xC5};

  ASSERT_EQ(
      decode_eps_quality_of_service(
          &eps_qos, 0, eps_qos_buffersize_Ext7,
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
  ASSERT_EQ(
      encode_eps_quality_of_service(
          &eps_qos, 0, eps_qos_buffersize_Ext7,
          sizeof(eps_qos_buffersize_Ext7)),
      sizeof(eps_qos_buffersize_encoded_Ext7));
  ASSERT_EQ(
      sizeof(eps_qos_buffersize_Ext7), sizeof(eps_qos_buffersize_encoded_Ext7));
  memcmp(
      eps_qos_buffersize_encoded_Ext7, eps_qos_buffersize_Ext7,
      sizeof(eps_qos_buffersize_Ext7));

  uint8_t eps_qos_buffersize_Ext8[] = {0x09, 0x09, 0xFE, 0xFE, 0xFE,
                                       0xFE, 0xFA, 0xFA, 0xFA, 0xFA};

  ASSERT_EQ(
      decode_eps_quality_of_service(
          &eps_qos, 0, eps_qos_buffersize_Ext8,
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
  ASSERT_EQ(
      encode_eps_quality_of_service(
          &eps_qos, 0, eps_qos_buffersize_Ext8,
          sizeof(eps_qos_buffersize_Ext8)),
      sizeof(eps_qos_buffersize_encoded_Ext8));
  ASSERT_EQ(
      sizeof(eps_qos_buffersize_Ext8), sizeof(eps_qos_buffersize_encoded_Ext8));
  memcmp(
      eps_qos_buffersize_encoded_Ext8, eps_qos_buffersize_Ext8,
      sizeof(eps_qos_buffersize_Ext8));

  uint8_t eps_qos_buffersize_Extended21[] = {0x0D, 0x09, 0xFE, 0xFE, 0xFE,
                                             0xFE, 0xFA, 0xFA, 0xFA, 0xFA,
                                             0x01, 0x01, 0x01, 0x01};
  ASSERT_EQ(
      decode_eps_quality_of_service(
          &eps_qos, 0, eps_qos_buffersize_Extended21,
          sizeof(eps_qos_buffersize_Extended21)),
      sizeof(eps_qos_buffersize_Extended21));
  ASSERT_EQ((eps_qos.bitRates.maxBitRateForUL), 254);
  ASSERT_EQ((eps_qos.bitRates.maxBitRateForDL), 254);
  ASSERT_EQ((eps_qos.bitRatesExt.maxBitRateForUL), 250);
  ASSERT_EQ((eps_qos.bitRatesExt.maxBitRateForDL), 250);
  ASSERT_EQ((eps_qos.bitRatesExt2.maxBitRateForUL), 1);
  ASSERT_EQ((eps_qos.bitRatesExt2.maxBitRateForDL), 1);
  uint8_t eps_qos_buffersize_encoded_Extended21[14] = {0};
  ASSERT_EQ(
      encode_eps_quality_of_service(
          &eps_qos, 0, eps_qos_buffersize_Extended21,
          sizeof(eps_qos_buffersize_Extended21)),
      sizeof(eps_qos_buffersize_encoded_Extended21));
  ASSERT_EQ(
      sizeof(eps_qos_buffersize_Extended21),
      sizeof(eps_qos_buffersize_encoded_Extended21));
  memcmp(
      eps_qos_buffersize_encoded_Extended21, eps_qos_buffersize_Extended21,
      sizeof(eps_qos_buffersize_Extended21));

  uint8_t eps_qos_buffersize_Extended22[] = {0x0d, 0x09, 0xFE, 0xFE, 0xFE,
                                             0xFE, 0xFA, 0xFA, 0xFA, 0xFA,
                                             0x0B, 0x0B, 0x0B, 0x0B};
  ASSERT_EQ(
      decode_eps_quality_of_service(
          &eps_qos, 0, eps_qos_buffersize_Extended22,
          sizeof(eps_qos_buffersize_Extended22)),
      sizeof(eps_qos_buffersize_Extended21));
  ASSERT_EQ((eps_qos.bitRates.maxBitRateForUL), 254);
  ASSERT_EQ((eps_qos.bitRates.maxBitRateForDL), 254);
  ASSERT_EQ((eps_qos.bitRatesExt.maxBitRateForUL), 250);
  ASSERT_EQ((eps_qos.bitRatesExt.maxBitRateForDL), 250);
  ASSERT_EQ((eps_qos.bitRatesExt2.maxBitRateForUL), 11);
  ASSERT_EQ((eps_qos.bitRatesExt2.maxBitRateForDL), 11);
  uint8_t eps_qos_buffersize_encoded_Extended22[14] = {0};
  ASSERT_EQ(
      encode_eps_quality_of_service(
          &eps_qos, 0, eps_qos_buffersize_Extended22,
          sizeof(eps_qos_buffersize_Extended22)),
      sizeof(eps_qos_buffersize_encoded_Extended22));
  ASSERT_EQ(
      sizeof(eps_qos_buffersize_Extended22),
      sizeof(eps_qos_buffersize_encoded_Extended22));
  memcmp(
      eps_qos_buffersize_encoded_Extended22, eps_qos_buffersize_Extended22,
      sizeof(eps_qos_buffersize_Extended21));

  uint8_t eps_qos_buffersize_Extended23[] = {0x0d, 0x09, 0xFE, 0xFE, 0xFE,
                                             0xFE, 0xFA, 0xFA, 0xFA, 0xFA,
                                             0x3D, 0x3D, 0x3D, 0x3D};
  ASSERT_EQ(
      decode_eps_quality_of_service(
          &eps_qos, 0, eps_qos_buffersize_Extended23,
          sizeof(eps_qos_buffersize_Extended23)),
      sizeof(eps_qos_buffersize_Extended21));
  ASSERT_EQ((eps_qos.bitRates.maxBitRateForUL), 254);
  ASSERT_EQ((eps_qos.bitRates.maxBitRateForDL), 254);
  ASSERT_EQ((eps_qos.bitRatesExt.maxBitRateForUL), 250);
  ASSERT_EQ((eps_qos.bitRatesExt.maxBitRateForDL), 250);
  ASSERT_EQ((eps_qos.bitRatesExt2.maxBitRateForUL), 61);
  ASSERT_EQ((eps_qos.bitRatesExt2.maxBitRateForDL), 61);
  uint8_t eps_qos_buffersize_encoded_Extended23[14] = {0};
  ASSERT_EQ(
      encode_eps_quality_of_service(
          &eps_qos, 0, eps_qos_buffersize_Extended23,
          sizeof(eps_qos_buffersize_Extended23)),
      sizeof(eps_qos_buffersize_encoded_Extended23));
  ASSERT_EQ(
      sizeof(eps_qos_buffersize_Extended23),
      sizeof(eps_qos_buffersize_encoded_Extended23));
  memcmp(
      eps_qos_buffersize_encoded_Extended23, eps_qos_buffersize_Extended23,
      sizeof(eps_qos_buffersize_Extended23));

  uint8_t eps_qos_buffersize_Extended24[] = {0x0d, 0x09, 0xFE, 0xFE, 0xFE,
                                             0xFE, 0xFA, 0xFA, 0xFA, 0xFA,
                                             0x3E, 0x3E, 0x3E, 0x3E};
  ASSERT_EQ(
      decode_eps_quality_of_service(
          &eps_qos, 0, eps_qos_buffersize_Extended24,
          sizeof(eps_qos_buffersize_Extended24)),
      sizeof(eps_qos_buffersize_Extended21));
  ASSERT_EQ((eps_qos.bitRates.maxBitRateForUL), 254);
  ASSERT_EQ((eps_qos.bitRates.maxBitRateForDL), 254);
  ASSERT_EQ((eps_qos.bitRatesExt.maxBitRateForUL), 250);
  ASSERT_EQ((eps_qos.bitRatesExt.maxBitRateForDL), 250);
  ASSERT_EQ((eps_qos.bitRatesExt2.maxBitRateForUL), 62);
  ASSERT_EQ((eps_qos.bitRatesExt2.maxBitRateForDL), 62);
  uint8_t eps_qos_buffersize_encoded_Extended24[14] = {0};
  ASSERT_EQ(
      encode_eps_quality_of_service(
          &eps_qos, 0, eps_qos_buffersize_Extended24,
          sizeof(eps_qos_buffersize_Extended24)),
      sizeof(eps_qos_buffersize_encoded_Extended24));
  ASSERT_EQ(
      sizeof(eps_qos_buffersize_Extended24),
      sizeof(eps_qos_buffersize_encoded_Extended24));
  memcmp(
      eps_qos_buffersize_encoded_Extended24, eps_qos_buffersize_Extended24,
      sizeof(eps_qos_buffersize_Extended24));

  uint8_t eps_qos_buffersize_Extended25[] = {0x0d, 0x09, 0xFE, 0xFE, 0xFE,
                                             0xFE, 0xFA, 0xFA, 0xFA, 0xFA,
                                             0x47, 0x47, 0x47, 0x47};
  ASSERT_EQ(
      decode_eps_quality_of_service(
          &eps_qos, 0, eps_qos_buffersize_Extended25,
          sizeof(eps_qos_buffersize_Extended25)),
      sizeof(eps_qos_buffersize_Extended21));
  ASSERT_EQ((eps_qos.bitRates.maxBitRateForUL), 254);
  ASSERT_EQ((eps_qos.bitRates.maxBitRateForDL), 254);
  ASSERT_EQ((eps_qos.bitRatesExt.maxBitRateForUL), 250);
  ASSERT_EQ((eps_qos.bitRatesExt.maxBitRateForDL), 250);
  ASSERT_EQ((eps_qos.bitRatesExt2.maxBitRateForUL), 71);
  ASSERT_EQ((eps_qos.bitRatesExt2.maxBitRateForDL), 71);
  uint8_t eps_qos_buffersize_encoded_Extended25[14] = {0};
  ASSERT_EQ(
      encode_eps_quality_of_service(
          &eps_qos, 0, eps_qos_buffersize_Extended25,
          sizeof(eps_qos_buffersize_Extended25)),
      sizeof(eps_qos_buffersize_encoded_Extended25));
  ASSERT_EQ(
      sizeof(eps_qos_buffersize_Extended25),
      sizeof(eps_qos_buffersize_encoded_Extended25));
  memcmp(
      eps_qos_buffersize_encoded_Extended25, eps_qos_buffersize_Extended25,
      sizeof(eps_qos_buffersize_Extended25));

  uint8_t eps_qos_buffersize_Extended26[] = {0x0d, 0x09, 0xFE, 0xFE, 0xFE,
                                             0xFE, 0xFA, 0xFA, 0xFA, 0xFA,
                                             0xA1, 0xA1, 0xA1, 0xA1};
  ASSERT_EQ(
      decode_eps_quality_of_service(
          &eps_qos, 0, eps_qos_buffersize_Extended26,
          sizeof(eps_qos_buffersize_Extended26)),
      sizeof(eps_qos_buffersize_Extended26));
  ASSERT_EQ((eps_qos.bitRates.maxBitRateForUL), 254);
  ASSERT_EQ((eps_qos.bitRates.maxBitRateForDL), 254);
  ASSERT_EQ((eps_qos.bitRatesExt.maxBitRateForUL), 250);
  ASSERT_EQ((eps_qos.bitRatesExt.maxBitRateForDL), 250);
  ASSERT_EQ((eps_qos.bitRatesExt2.maxBitRateForUL), 161);
  ASSERT_EQ((eps_qos.bitRatesExt2.maxBitRateForDL), 161);
  uint8_t eps_qos_buffersize_encoded_Extended26[14] = {0};
  ASSERT_EQ(
      encode_eps_quality_of_service(
          &eps_qos, 0, eps_qos_buffersize_Extended26,
          sizeof(eps_qos_buffersize_Extended26)),
      sizeof(eps_qos_buffersize_encoded_Extended26));
  ASSERT_EQ(
      sizeof(eps_qos_buffersize_Extended26),
      sizeof(eps_qos_buffersize_encoded_Extended26));
  memcmp(
      eps_qos_buffersize_encoded_Extended26, eps_qos_buffersize_Extended26,
      sizeof(eps_qos_buffersize_Extended26));

  uint8_t eps_qos_buffersize_Extended27[] = {0x0d, 0x09, 0xFE, 0xFE, 0xFE,
                                             0xFE, 0xFA, 0xFA, 0xFA, 0xFA,
                                             0xA2, 0xA2, 0xA2, 0xA2};
  ASSERT_EQ(
      decode_eps_quality_of_service(
          &eps_qos, 0, eps_qos_buffersize_Extended27,
          sizeof(eps_qos_buffersize_Extended27)),
      sizeof(eps_qos_buffersize_Extended27));
  ASSERT_EQ((eps_qos.bitRates.maxBitRateForUL), 254);
  ASSERT_EQ((eps_qos.bitRates.maxBitRateForDL), 254);
  ASSERT_EQ((eps_qos.bitRatesExt.maxBitRateForUL), 250);
  ASSERT_EQ((eps_qos.bitRatesExt.maxBitRateForDL), 250);
  ASSERT_EQ((eps_qos.bitRatesExt2.maxBitRateForUL), 162);
  ASSERT_EQ((eps_qos.bitRatesExt2.maxBitRateForDL), 162);
  uint8_t eps_qos_buffersize_encoded_Extended27[14] = {0};
  ASSERT_EQ(
      encode_eps_quality_of_service(
          &eps_qos, 0, eps_qos_buffersize_Extended27,
          sizeof(eps_qos_buffersize_Extended27)),
      sizeof(eps_qos_buffersize_encoded_Extended27));
  ASSERT_EQ(
      sizeof(eps_qos_buffersize_Extended27),
      sizeof(eps_qos_buffersize_encoded_Extended27));
  memcmp(
      eps_qos_buffersize_encoded_Extended27, eps_qos_buffersize_Extended27,
      sizeof(eps_qos_buffersize_Extended27));

  uint8_t eps_qos_buffersize_Extended28[] = {0x0d, 0x09, 0xFE, 0xFE, 0xFE,
                                             0xFE, 0xFA, 0xFA, 0xFA, 0xFA,
                                             0xA6, 0xA6, 0xA6, 0xA6};
  ASSERT_EQ(
      decode_eps_quality_of_service(
          &eps_qos, 0, eps_qos_buffersize_Extended28,
          sizeof(eps_qos_buffersize_Extended28)),
      sizeof(eps_qos_buffersize_Extended28));
  ASSERT_EQ((eps_qos.bitRates.maxBitRateForUL), 254);
  ASSERT_EQ((eps_qos.bitRates.maxBitRateForDL), 254);
  ASSERT_EQ((eps_qos.bitRatesExt.maxBitRateForUL), 250);
  ASSERT_EQ((eps_qos.bitRatesExt.maxBitRateForDL), 250);
  ASSERT_EQ((eps_qos.bitRatesExt2.maxBitRateForUL), 166);
  ASSERT_EQ((eps_qos.bitRatesExt2.maxBitRateForDL), 166);
  uint8_t eps_qos_buffersize_encoded_Extended28[14] = {0};
  ASSERT_EQ(
      encode_eps_quality_of_service(
          &eps_qos, 0, eps_qos_buffersize_Extended28,
          sizeof(eps_qos_buffersize_Extended28)),
      sizeof(eps_qos_buffersize_encoded_Extended28));
  ASSERT_EQ(
      sizeof(eps_qos_buffersize_Extended28),
      sizeof(eps_qos_buffersize_encoded_Extended28));
  memcmp(
      eps_qos_buffersize_encoded_Extended28, eps_qos_buffersize_Extended28,
      sizeof(eps_qos_buffersize_Extended28));

  uint8_t eps_qos_buffersize_Extended29[] = {0x0d, 0x09, 0xFE, 0xFE, 0xFE,
                                             0xFE, 0xFA, 0xFA, 0xFA, 0xFA,
                                             0xF6, 0xF6, 0xF6, 0xF6};
  ASSERT_EQ(
      decode_eps_quality_of_service(
          &eps_qos, 0, eps_qos_buffersize_Extended29,
          sizeof(eps_qos_buffersize_Extended29)),
      sizeof(eps_qos_buffersize_Extended29));
  ASSERT_EQ((eps_qos.bitRates.maxBitRateForUL), 254);
  ASSERT_EQ((eps_qos.bitRates.maxBitRateForDL), 254);
  ASSERT_EQ((eps_qos.bitRatesExt.maxBitRateForUL), 250);
  ASSERT_EQ((eps_qos.bitRatesExt.maxBitRateForDL), 250);
  ASSERT_EQ((eps_qos.bitRatesExt2.maxBitRateForUL), 246);
  ASSERT_EQ((eps_qos.bitRatesExt2.maxBitRateForDL), 246);
  uint8_t eps_qos_buffersize_encoded_Extended29[14] = {0};
  ASSERT_EQ(
      encode_eps_quality_of_service(
          &eps_qos, 0, eps_qos_buffersize_Extended29,
          sizeof(eps_qos_buffersize_Extended29)),
      sizeof(eps_qos_buffersize_encoded_Extended29));
  ASSERT_EQ(
      sizeof(eps_qos_buffersize_Extended29),
      sizeof(eps_qos_buffersize_encoded_Extended29));
  memcmp(
      eps_qos_buffersize_encoded_Extended29, eps_qos_buffersize_Extended29,
      sizeof(eps_qos_buffersize_Extended29));
}

int main(int argc, char** argv) {
  ::testing::InitGoogleTest(&argc, argv);
  OAILOG_INIT("MME", OAILOG_LEVEL_DEBUG, MAX_LOG_PROTOS);
  return RUN_ALL_TESTS();
}

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
#include "AttachRequest.h"
#include "ExtendedApnAggregateMaximumBitRate.h"
#include "dynamic_memory_check.h"
#include "log.h"
}

namespace magma {
namespace lte {

class EMMDecodeTest : public ::testing::Test {
  virtual void SetUp() {}
  virtual void TearDown() {}
};

TEST_F(EMMDecodeTest, TestDecodeAttachRequest1) {
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

TEST_F(EMMDecodeTest, TestDecodeAttachRequest2) {
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

TEST_F(EMMDecodeTest, TestDecodeEncodeExtendedAPNAMBR) {
  ExtendedApnAggregateMaximumBitRate extended_apn_ambr = {0};

  uint8_t extended_apn_ambr_buffersize1[] = {0x06, 0x03, 0x7A, 0x3F,
                                             0x03, 0x7A, 0x3F};
  ASSERT_EQ(
      decode_extended_apn_aggregate_maximum_bit_rate(
          &extended_apn_ambr, 0, extended_apn_ambr_buffersize1,
          sizeof(extended_apn_ambr_buffersize1)),
      sizeof(extended_apn_ambr_buffersize1));
  ASSERT_EQ((extended_apn_ambr.extendedapnambrfordownlinkunit), 3);
  ASSERT_EQ((extended_apn_ambr.extendedapnambrfordownlink), 122);
  ASSERT_EQ((extended_apn_ambr.extendedapnambrfordownlink_continued), 63);
  ASSERT_EQ((extended_apn_ambr.extendedapnambrforuplinkunit), 3);
  ASSERT_EQ((extended_apn_ambr.extendedapnambrforuplink), 122);
  ASSERT_EQ((extended_apn_ambr.extendedapnambrforuplink_continued), 63);

  uint8_t extended_apn_ambr_encoded_buffersize1[7] = {0};
  ASSERT_EQ(
      encode_extended_apn_aggregate_maximum_bit_rate(
          &extended_apn_ambr, 0, extended_apn_ambr_encoded_buffersize1,
          sizeof(extended_apn_ambr_encoded_buffersize1)),
      sizeof(extended_apn_ambr_encoded_buffersize1));
  ASSERT_EQ(
      (memcmp(
          extended_apn_ambr_encoded_buffersize1, extended_apn_ambr_buffersize1,
          sizeof(extended_apn_ambr_buffersize1))),
      0);

  uint8_t extended_apn_ambr_buffersize2[] = {0x06, 0x03, 0xA8, 0x61,
                                             0x03, 0xA8, 0x61};
  ASSERT_EQ(
      decode_extended_apn_aggregate_maximum_bit_rate(
          &extended_apn_ambr, 0, extended_apn_ambr_buffersize2,
          sizeof(extended_apn_ambr_buffersize2)),
      sizeof(extended_apn_ambr_buffersize2));
  ASSERT_EQ((extended_apn_ambr.extendedapnambrfordownlinkunit), 3);
  ASSERT_EQ((extended_apn_ambr.extendedapnambrfordownlink), 168);
  ASSERT_EQ((extended_apn_ambr.extendedapnambrfordownlink_continued), 97);
  ASSERT_EQ((extended_apn_ambr.extendedapnambrforuplinkunit), 3);
  ASSERT_EQ((extended_apn_ambr.extendedapnambrforuplink), 168);
  ASSERT_EQ((extended_apn_ambr.extendedapnambrforuplink_continued), 97);

  uint8_t extended_apn_ambr_encoded_buffersize2[7] = {0};
  ASSERT_EQ(
      encode_extended_apn_aggregate_maximum_bit_rate(
          &extended_apn_ambr, 0, extended_apn_ambr_encoded_buffersize2,
          sizeof(extended_apn_ambr_encoded_buffersize2)),
      sizeof(extended_apn_ambr_encoded_buffersize2));
  ASSERT_EQ(
      (memcmp(
          extended_apn_ambr_encoded_buffersize2, extended_apn_ambr_buffersize2,
          sizeof(extended_apn_ambr_buffersize2))),
      0);

  uint8_t extended_apn_ambr_buffersize3[] = {0x06, 0x03, 0xFF, 0xFF,
                                             0x03, 0xFF, 0xFF};
  ASSERT_EQ(
      decode_extended_apn_aggregate_maximum_bit_rate(
          &extended_apn_ambr, 0, extended_apn_ambr_buffersize3,
          sizeof(extended_apn_ambr_buffersize3)),
      sizeof(extended_apn_ambr_buffersize3));
  ASSERT_EQ((extended_apn_ambr.extendedapnambrfordownlinkunit), 3);
  ASSERT_EQ((extended_apn_ambr.extendedapnambrfordownlink), 255);
  ASSERT_EQ((extended_apn_ambr.extendedapnambrfordownlink_continued), 255);
  ASSERT_EQ((extended_apn_ambr.extendedapnambrforuplinkunit), 3);
  ASSERT_EQ((extended_apn_ambr.extendedapnambrforuplink), 255);
  ASSERT_EQ((extended_apn_ambr.extendedapnambrforuplink_continued), 255);

  uint8_t extended_apn_ambr_encoded_buffersize3[7] = {0};
  ASSERT_EQ(
      encode_extended_apn_aggregate_maximum_bit_rate(
          &extended_apn_ambr, 0, extended_apn_ambr_encoded_buffersize3,
          sizeof(extended_apn_ambr_encoded_buffersize3)),
      sizeof(extended_apn_ambr_encoded_buffersize3));
  ASSERT_EQ(
      (memcmp(
          extended_apn_ambr_encoded_buffersize3, extended_apn_ambr_buffersize3,
          sizeof(extended_apn_ambr_buffersize3))),
      0);

  uint8_t extended_apn_ambr_buffersize4[] = {0x06, 0x04, 0x24, 0xF4,
                                             0x04, 0x24, 0xF4};
  ASSERT_EQ(
      decode_extended_apn_aggregate_maximum_bit_rate(
          &extended_apn_ambr, 0, extended_apn_ambr_buffersize4,
          sizeof(extended_apn_ambr_buffersize4)),
      sizeof(extended_apn_ambr_buffersize4));
  ASSERT_EQ((extended_apn_ambr.extendedapnambrfordownlinkunit), 4);
  ASSERT_EQ((extended_apn_ambr.extendedapnambrfordownlink), 36);
  ASSERT_EQ((extended_apn_ambr.extendedapnambrfordownlink_continued), 244);
  ASSERT_EQ((extended_apn_ambr.extendedapnambrforuplinkunit), 4);
  ASSERT_EQ((extended_apn_ambr.extendedapnambrforuplink), 36);
  ASSERT_EQ((extended_apn_ambr.extendedapnambrforuplink_continued), 244);

  uint8_t extended_apn_ambr_encoded_buffersize4[7] = {0};
  ASSERT_EQ(
      encode_extended_apn_aggregate_maximum_bit_rate(
          &extended_apn_ambr, 0, extended_apn_ambr_encoded_buffersize4,
          sizeof(extended_apn_ambr_encoded_buffersize4)),
      sizeof(extended_apn_ambr_encoded_buffersize4));
  ASSERT_EQ(
      (memcmp(
          extended_apn_ambr_encoded_buffersize4, extended_apn_ambr_buffersize4,
          sizeof(extended_apn_ambr_buffersize4))),
      0);

  uint8_t extended_apn_ambr_buffersize5[] = {0x06, 0x04, 0xFF, 0xFF,
                                             0x04, 0xFF, 0xFF};
  ASSERT_EQ(
      decode_extended_apn_aggregate_maximum_bit_rate(
          &extended_apn_ambr, 0, extended_apn_ambr_buffersize5,
          sizeof(extended_apn_ambr_buffersize5)),
      sizeof(extended_apn_ambr_buffersize5));
  ASSERT_EQ((extended_apn_ambr.extendedapnambrfordownlinkunit), 4);
  ASSERT_EQ((extended_apn_ambr.extendedapnambrfordownlink), 255);
  ASSERT_EQ((extended_apn_ambr.extendedapnambrfordownlink_continued), 255);
  ASSERT_EQ((extended_apn_ambr.extendedapnambrforuplinkunit), 4);
  ASSERT_EQ((extended_apn_ambr.extendedapnambrforuplink), 255);
  ASSERT_EQ((extended_apn_ambr.extendedapnambrforuplink_continued), 255);

  uint8_t extended_apn_ambr_encoded_buffersize5[7] = {0};
  ASSERT_EQ(
      encode_extended_apn_aggregate_maximum_bit_rate(
          &extended_apn_ambr, 0, extended_apn_ambr_encoded_buffersize5,
          sizeof(extended_apn_ambr_encoded_buffersize5)),
      sizeof(extended_apn_ambr_encoded_buffersize5));
  ASSERT_EQ(
      (memcmp(
          extended_apn_ambr_encoded_buffersize5, extended_apn_ambr_buffersize5,
          sizeof(extended_apn_ambr_buffersize5))),
      0);

  uint8_t extended_apn_ambr_buffersize6[] = {0x06, 0x05, 0x96, 0x98,
                                             0x05, 0x96, 0x98};
  ASSERT_EQ(
      decode_extended_apn_aggregate_maximum_bit_rate(
          &extended_apn_ambr, 0, extended_apn_ambr_buffersize6,
          sizeof(extended_apn_ambr_buffersize6)),
      sizeof(extended_apn_ambr_buffersize6));
  ASSERT_EQ((extended_apn_ambr.extendedapnambrfordownlinkunit), 5);
  ASSERT_EQ((extended_apn_ambr.extendedapnambrfordownlink), 150);
  ASSERT_EQ((extended_apn_ambr.extendedapnambrfordownlink_continued), 152);
  ASSERT_EQ((extended_apn_ambr.extendedapnambrforuplinkunit), 5);
  ASSERT_EQ((extended_apn_ambr.extendedapnambrforuplink), 150);
  ASSERT_EQ((extended_apn_ambr.extendedapnambrforuplink_continued), 152);

  uint8_t extended_apn_ambr_encoded_buffersize6[7] = {0};
  ASSERT_EQ(
      encode_extended_apn_aggregate_maximum_bit_rate(
          &extended_apn_ambr, 0, extended_apn_ambr_encoded_buffersize6,
          sizeof(extended_apn_ambr_encoded_buffersize6)),
      sizeof(extended_apn_ambr_encoded_buffersize6));
  ASSERT_EQ(
      (memcmp(
          extended_apn_ambr_encoded_buffersize6, extended_apn_ambr_buffersize6,
          sizeof(extended_apn_ambr_buffersize6))),
      0);

  uint8_t extended_apn_ambr_buffersize7[] = {0x06, 0x05, 0xFF, 0xFF,
                                             0x05, 0xFF, 0xFF};
  ASSERT_EQ(
      decode_extended_apn_aggregate_maximum_bit_rate(
          &extended_apn_ambr, 0, extended_apn_ambr_buffersize7,
          sizeof(extended_apn_ambr_buffersize7)),
      sizeof(extended_apn_ambr_buffersize7));
  ASSERT_EQ((extended_apn_ambr.extendedapnambrfordownlinkunit), 5);
  ASSERT_EQ((extended_apn_ambr.extendedapnambrfordownlink), 255);
  ASSERT_EQ((extended_apn_ambr.extendedapnambrfordownlink_continued), 255);
  ASSERT_EQ((extended_apn_ambr.extendedapnambrforuplinkunit), 5);
  ASSERT_EQ((extended_apn_ambr.extendedapnambrforuplink), 255);
  ASSERT_EQ((extended_apn_ambr.extendedapnambrforuplink_continued), 255);

  uint8_t extended_apn_ambr_encoded_buffersize7[7] = {0};
  ASSERT_EQ(
      encode_extended_apn_aggregate_maximum_bit_rate(
          &extended_apn_ambr, 0, extended_apn_ambr_encoded_buffersize7,
          sizeof(extended_apn_ambr_encoded_buffersize7)),
      sizeof(extended_apn_ambr_encoded_buffersize7));
  ASSERT_EQ(
      (memcmp(
          extended_apn_ambr_encoded_buffersize7, extended_apn_ambr_buffersize7,
          sizeof(extended_apn_ambr_buffersize7))),
      0);

  uint8_t extended_apn_ambr_buffersize8[] = {0x06, 0x06, 0x16, 0x40,
                                             0x06, 0x16, 0x40};
  ASSERT_EQ(
      decode_extended_apn_aggregate_maximum_bit_rate(
          &extended_apn_ambr, 0, extended_apn_ambr_buffersize8,
          sizeof(extended_apn_ambr_buffersize8)),
      sizeof(extended_apn_ambr_buffersize8));
  ASSERT_EQ((extended_apn_ambr.extendedapnambrfordownlinkunit), 6);
  ASSERT_EQ((extended_apn_ambr.extendedapnambrfordownlink), 22);
  ASSERT_EQ((extended_apn_ambr.extendedapnambrfordownlink_continued), 64);
  ASSERT_EQ((extended_apn_ambr.extendedapnambrforuplinkunit), 6);
  ASSERT_EQ((extended_apn_ambr.extendedapnambrforuplink), 22);
  ASSERT_EQ((extended_apn_ambr.extendedapnambrforuplink_continued), 64);

  uint8_t extended_apn_ambr_encoded_buffersize8[7] = {0};
  ASSERT_EQ(
      encode_extended_apn_aggregate_maximum_bit_rate(
          &extended_apn_ambr, 0, extended_apn_ambr_encoded_buffersize8,
          sizeof(extended_apn_ambr_encoded_buffersize8)),
      sizeof(extended_apn_ambr_encoded_buffersize8));
  ASSERT_EQ(
      (memcmp(
          extended_apn_ambr_encoded_buffersize8, extended_apn_ambr_buffersize8,
          sizeof(extended_apn_ambr_buffersize8))),
      0);

  uint8_t extended_apn_ambr_buffersize9[] = {0x06, 0x06, 0xFF, 0xFF,
                                             0x06, 0xFF, 0xFF};
  ASSERT_EQ(
      decode_extended_apn_aggregate_maximum_bit_rate(
          &extended_apn_ambr, 0, extended_apn_ambr_buffersize9,
          sizeof(extended_apn_ambr_buffersize9)),
      sizeof(extended_apn_ambr_buffersize9));
  ASSERT_EQ((extended_apn_ambr.extendedapnambrfordownlinkunit), 6);
  ASSERT_EQ((extended_apn_ambr.extendedapnambrfordownlink), 255);
  ASSERT_EQ((extended_apn_ambr.extendedapnambrfordownlink_continued), 255);
  ASSERT_EQ((extended_apn_ambr.extendedapnambrforuplinkunit), 6);
  ASSERT_EQ((extended_apn_ambr.extendedapnambrforuplink), 255);
  ASSERT_EQ((extended_apn_ambr.extendedapnambrforuplink_continued), 255);

  uint8_t extended_apn_ambr_encoded_buffersize9[7] = {0};
  ASSERT_EQ(
      encode_extended_apn_aggregate_maximum_bit_rate(
          &extended_apn_ambr, 0, extended_apn_ambr_encoded_buffersize9,
          sizeof(extended_apn_ambr_encoded_buffersize9)),
      sizeof(extended_apn_ambr_encoded_buffersize9));
  ASSERT_EQ(
      (memcmp(
          extended_apn_ambr_encoded_buffersize9, extended_apn_ambr_buffersize9,
          sizeof(extended_apn_ambr_buffersize9))),
      0);

  uint8_t extended_apn_ambr_buffersize10[] = {0x06, 0x07, 0x4B, 0x4C,
                                              0x07, 0x4B, 0x4C};
  ASSERT_EQ(
      decode_extended_apn_aggregate_maximum_bit_rate(
          &extended_apn_ambr, 0, extended_apn_ambr_buffersize10,
          sizeof(extended_apn_ambr_buffersize10)),
      sizeof(extended_apn_ambr_buffersize10));
  ASSERT_EQ((extended_apn_ambr.extendedapnambrfordownlinkunit), 7);
  ASSERT_EQ((extended_apn_ambr.extendedapnambrfordownlink), 75);
  ASSERT_EQ((extended_apn_ambr.extendedapnambrfordownlink_continued), 76);
  ASSERT_EQ((extended_apn_ambr.extendedapnambrforuplinkunit), 7);
  ASSERT_EQ((extended_apn_ambr.extendedapnambrforuplink), 75);
  ASSERT_EQ((extended_apn_ambr.extendedapnambrforuplink_continued), 76);

  uint8_t extended_apn_ambr_encoded_buffersize10[7] = {0};
  ASSERT_EQ(
      encode_extended_apn_aggregate_maximum_bit_rate(
          &extended_apn_ambr, 0, extended_apn_ambr_encoded_buffersize10,
          sizeof(extended_apn_ambr_encoded_buffersize10)),
      sizeof(extended_apn_ambr_encoded_buffersize10));
  ASSERT_EQ(
      (memcmp(
          extended_apn_ambr_encoded_buffersize10,
          extended_apn_ambr_buffersize10,
          sizeof(extended_apn_ambr_buffersize10))),
      0);

  uint8_t extended_apn_ambr_buffersize11[] = {0x06, 0x07, 0xFF, 0xFF,
                                              0x07, 0xFF, 0xFF};
  ASSERT_EQ(
      decode_extended_apn_aggregate_maximum_bit_rate(
          &extended_apn_ambr, 0, extended_apn_ambr_buffersize11,
          sizeof(extended_apn_ambr_buffersize11)),
      sizeof(extended_apn_ambr_buffersize11));
  ASSERT_EQ((extended_apn_ambr.extendedapnambrfordownlinkunit), 7);
  ASSERT_EQ((extended_apn_ambr.extendedapnambrfordownlink), 255);
  ASSERT_EQ((extended_apn_ambr.extendedapnambrfordownlink_continued), 255);
  ASSERT_EQ((extended_apn_ambr.extendedapnambrforuplinkunit), 7);
  ASSERT_EQ((extended_apn_ambr.extendedapnambrforuplink), 255);
  ASSERT_EQ((extended_apn_ambr.extendedapnambrforuplink_continued), 255);

  uint8_t extended_apn_ambr_encoded_buffersize11[7] = {0};
  ASSERT_EQ(
      encode_extended_apn_aggregate_maximum_bit_rate(
          &extended_apn_ambr, 0, extended_apn_ambr_encoded_buffersize11,
          sizeof(extended_apn_ambr_encoded_buffersize11)),
      sizeof(extended_apn_ambr_encoded_buffersize11));
  ASSERT_EQ(
      (memcmp(
          extended_apn_ambr_encoded_buffersize11,
          extended_apn_ambr_buffersize11,
          sizeof(extended_apn_ambr_buffersize11))),
      0);

  uint8_t extended_apn_ambr_buffersize12[] = {0x06, 0x08, 0x6B, 0xEE,
                                              0x08, 0x6B, 0xEE};
  ASSERT_EQ(
      decode_extended_apn_aggregate_maximum_bit_rate(
          &extended_apn_ambr, 0, extended_apn_ambr_buffersize12,
          sizeof(extended_apn_ambr_buffersize12)),
      sizeof(extended_apn_ambr_buffersize12));
  ASSERT_EQ((extended_apn_ambr.extendedapnambrfordownlinkunit), 8);
  ASSERT_EQ((extended_apn_ambr.extendedapnambrfordownlink), 107);
  ASSERT_EQ((extended_apn_ambr.extendedapnambrfordownlink_continued), 238);
  ASSERT_EQ((extended_apn_ambr.extendedapnambrforuplinkunit), 8);
  ASSERT_EQ((extended_apn_ambr.extendedapnambrforuplink), 107);
  ASSERT_EQ((extended_apn_ambr.extendedapnambrforuplink_continued), 238);

  uint8_t extended_apn_ambr_encoded_buffersize12[7] = {0};
  ASSERT_EQ(
      encode_extended_apn_aggregate_maximum_bit_rate(
          &extended_apn_ambr, 0, extended_apn_ambr_encoded_buffersize12,
          sizeof(extended_apn_ambr_encoded_buffersize12)),
      sizeof(extended_apn_ambr_encoded_buffersize12));
  ASSERT_EQ(
      (memcmp(
          extended_apn_ambr_encoded_buffersize12,
          extended_apn_ambr_buffersize12,
          sizeof(extended_apn_ambr_buffersize12))),
      0);

  uint8_t extended_apn_ambr_buffersize13[] = {0x06, 0x08, 0xFF, 0xFF,
                                              0x08, 0xFF, 0xFF};
  ASSERT_EQ(
      decode_extended_apn_aggregate_maximum_bit_rate(
          &extended_apn_ambr, 0, extended_apn_ambr_buffersize13,
          sizeof(extended_apn_ambr_buffersize13)),
      sizeof(extended_apn_ambr_buffersize13));
  ASSERT_EQ((extended_apn_ambr.extendedapnambrfordownlinkunit), 8);
  ASSERT_EQ((extended_apn_ambr.extendedapnambrfordownlink), 255);
  ASSERT_EQ((extended_apn_ambr.extendedapnambrfordownlink_continued), 255);
  ASSERT_EQ((extended_apn_ambr.extendedapnambrforuplinkunit), 8);
  ASSERT_EQ((extended_apn_ambr.extendedapnambrforuplink), 255);
  ASSERT_EQ((extended_apn_ambr.extendedapnambrforuplink_continued), 255);

  uint8_t extended_apn_ambr_encoded_buffersize13[7] = {0};
  ASSERT_EQ(
      encode_extended_apn_aggregate_maximum_bit_rate(
          &extended_apn_ambr, 0, extended_apn_ambr_encoded_buffersize13,
          sizeof(extended_apn_ambr_encoded_buffersize13)),
      sizeof(extended_apn_ambr_encoded_buffersize13));
  ASSERT_EQ(
      (memcmp(
          extended_apn_ambr_encoded_buffersize13,
          extended_apn_ambr_buffersize13,
          sizeof(extended_apn_ambr_buffersize13))),
      0);

  uint8_t extended_apn_ambr_buffersize14[] = {0x06, 0x09, 0x00, 0x40,
                                              0x09, 0x00, 0x40};
  ASSERT_EQ(
      decode_extended_apn_aggregate_maximum_bit_rate(
          &extended_apn_ambr, 0, extended_apn_ambr_buffersize14,
          sizeof(extended_apn_ambr_buffersize14)),
      sizeof(extended_apn_ambr_buffersize14));
  ASSERT_EQ((extended_apn_ambr.extendedapnambrfordownlinkunit), 9);
  ASSERT_EQ((extended_apn_ambr.extendedapnambrfordownlink), 0);
  ASSERT_EQ((extended_apn_ambr.extendedapnambrfordownlink_continued), 64);
  ASSERT_EQ((extended_apn_ambr.extendedapnambrforuplinkunit), 9);
  ASSERT_EQ((extended_apn_ambr.extendedapnambrforuplink), 0);
  ASSERT_EQ((extended_apn_ambr.extendedapnambrforuplink_continued), 64);

  uint8_t extended_apn_ambr_encoded_buffersize14[7] = {0};
  ASSERT_EQ(
      encode_extended_apn_aggregate_maximum_bit_rate(
          &extended_apn_ambr, 0, extended_apn_ambr_encoded_buffersize14,
          sizeof(extended_apn_ambr_encoded_buffersize14)),
      sizeof(extended_apn_ambr_encoded_buffersize14));
  ASSERT_EQ(
      (memcmp(
          extended_apn_ambr_encoded_buffersize14,
          extended_apn_ambr_buffersize14,
          sizeof(extended_apn_ambr_buffersize14))),
      0);

  uint8_t extended_apn_ambr_buffersize15[] = {0x06, 0x09, 0xFF, 0xFF,
                                              0x09, 0xFF, 0xFF};
  ASSERT_EQ(
      decode_extended_apn_aggregate_maximum_bit_rate(
          &extended_apn_ambr, 0, extended_apn_ambr_buffersize15,
          sizeof(extended_apn_ambr_buffersize15)),
      sizeof(extended_apn_ambr_buffersize15));
  ASSERT_EQ((extended_apn_ambr.extendedapnambrfordownlinkunit), 9);
  ASSERT_EQ((extended_apn_ambr.extendedapnambrfordownlink), 255);
  ASSERT_EQ((extended_apn_ambr.extendedapnambrfordownlink_continued), 255);
  ASSERT_EQ((extended_apn_ambr.extendedapnambrforuplinkunit), 9);
  ASSERT_EQ((extended_apn_ambr.extendedapnambrforuplink), 255);
  ASSERT_EQ((extended_apn_ambr.extendedapnambrforuplink_continued), 255);

  uint8_t extended_apn_ambr_encoded_buffersize15[7] = {0};
  ASSERT_EQ(
      encode_extended_apn_aggregate_maximum_bit_rate(
          &extended_apn_ambr, 0, extended_apn_ambr_encoded_buffersize15,
          sizeof(extended_apn_ambr_encoded_buffersize15)),
      sizeof(extended_apn_ambr_encoded_buffersize15));
  ASSERT_EQ(
      (memcmp(
          extended_apn_ambr_encoded_buffersize15,
          extended_apn_ambr_buffersize15,
          sizeof(extended_apn_ambr_buffersize15))),
      0);

  uint8_t extended_apn_ambr_buffersize16[] = {0x06, 0x0A, 0xD0, 0xB2,
                                              0x0A, 0xD0, 0xB2};
  ASSERT_EQ(
      decode_extended_apn_aggregate_maximum_bit_rate(
          &extended_apn_ambr, 0, extended_apn_ambr_buffersize16,
          sizeof(extended_apn_ambr_buffersize16)),
      sizeof(extended_apn_ambr_buffersize16));
  ASSERT_EQ((extended_apn_ambr.extendedapnambrfordownlinkunit), 10);
  ASSERT_EQ((extended_apn_ambr.extendedapnambrfordownlink), 208);
  ASSERT_EQ((extended_apn_ambr.extendedapnambrfordownlink_continued), 178);
  ASSERT_EQ((extended_apn_ambr.extendedapnambrforuplinkunit), 10);
  ASSERT_EQ((extended_apn_ambr.extendedapnambrforuplink), 208);
  ASSERT_EQ((extended_apn_ambr.extendedapnambrforuplink_continued), 178);

  uint8_t extended_apn_ambr_encoded_buffersize16[7] = {0};
  ASSERT_EQ(
      encode_extended_apn_aggregate_maximum_bit_rate(
          &extended_apn_ambr, 0, extended_apn_ambr_encoded_buffersize16,
          sizeof(extended_apn_ambr_encoded_buffersize16)),
      sizeof(extended_apn_ambr_encoded_buffersize16));
  ASSERT_EQ(
      (memcmp(
          extended_apn_ambr_encoded_buffersize16,
          extended_apn_ambr_buffersize16,
          sizeof(extended_apn_ambr_buffersize16))),
      0);

  uint8_t extended_apn_ambr_buffersize17[] = {0x06, 0x0A, 0xFF, 0xFF,
                                              0x0A, 0xFF, 0xFF};
  ASSERT_EQ(
      decode_extended_apn_aggregate_maximum_bit_rate(
          &extended_apn_ambr, 0, extended_apn_ambr_buffersize17,
          sizeof(extended_apn_ambr_buffersize17)),
      sizeof(extended_apn_ambr_buffersize17));
  ASSERT_EQ((extended_apn_ambr.extendedapnambrfordownlinkunit), 10);
  ASSERT_EQ((extended_apn_ambr.extendedapnambrfordownlink), 255);
  ASSERT_EQ((extended_apn_ambr.extendedapnambrfordownlink_continued), 255);
  ASSERT_EQ((extended_apn_ambr.extendedapnambrforuplinkunit), 10);
  ASSERT_EQ((extended_apn_ambr.extendedapnambrforuplink), 255);
  ASSERT_EQ((extended_apn_ambr.extendedapnambrforuplink_continued), 255);

  uint8_t extended_apn_ambr_encoded_buffersize17[7] = {0};
  ASSERT_EQ(
      encode_extended_apn_aggregate_maximum_bit_rate(
          &extended_apn_ambr, 0, extended_apn_ambr_encoded_buffersize17,
          sizeof(extended_apn_ambr_encoded_buffersize17)),
      sizeof(extended_apn_ambr_encoded_buffersize17));
  ASSERT_EQ(
      (memcmp(
          extended_apn_ambr_encoded_buffersize17,
          extended_apn_ambr_buffersize17,
          sizeof(extended_apn_ambr_buffersize17))),
      0);

  uint8_t extended_apn_ambr_buffersize18[] = {0x06, 0x0B, 0x68, 0x59,
                                              0x0B, 0x68, 0x59};
  ASSERT_EQ(
      decode_extended_apn_aggregate_maximum_bit_rate(
          &extended_apn_ambr, 0, extended_apn_ambr_buffersize18,
          sizeof(extended_apn_ambr_buffersize18)),
      sizeof(extended_apn_ambr_buffersize18));
  ASSERT_EQ((extended_apn_ambr.extendedapnambrfordownlinkunit), 11);
  ASSERT_EQ((extended_apn_ambr.extendedapnambrfordownlink), 104);
  ASSERT_EQ((extended_apn_ambr.extendedapnambrfordownlink_continued), 89);
  ASSERT_EQ((extended_apn_ambr.extendedapnambrforuplinkunit), 11);
  ASSERT_EQ((extended_apn_ambr.extendedapnambrforuplink), 104);
  ASSERT_EQ((extended_apn_ambr.extendedapnambrforuplink_continued), 89);

  uint8_t extended_apn_ambr_encoded_buffersize18[7] = {0};
  ASSERT_EQ(
      encode_extended_apn_aggregate_maximum_bit_rate(
          &extended_apn_ambr, 0, extended_apn_ambr_encoded_buffersize18,
          sizeof(extended_apn_ambr_encoded_buffersize18)),
      sizeof(extended_apn_ambr_encoded_buffersize18));
  ASSERT_EQ(
      (memcmp(
          extended_apn_ambr_encoded_buffersize18,
          extended_apn_ambr_buffersize18,
          sizeof(extended_apn_ambr_buffersize18))),
      0);

  uint8_t extended_apn_ambr_buffersize19[] = {0x06, 0x0B, 0xFF, 0xFF,
                                              0x0B, 0xFF, 0xFF};
  ASSERT_EQ(
      decode_extended_apn_aggregate_maximum_bit_rate(
          &extended_apn_ambr, 0, extended_apn_ambr_buffersize19,
          sizeof(extended_apn_ambr_buffersize19)),
      sizeof(extended_apn_ambr_buffersize19));
  ASSERT_EQ((extended_apn_ambr.extendedapnambrfordownlinkunit), 11);
  ASSERT_EQ((extended_apn_ambr.extendedapnambrfordownlink), 255);
  ASSERT_EQ((extended_apn_ambr.extendedapnambrfordownlink_continued), 255);
  ASSERT_EQ((extended_apn_ambr.extendedapnambrforuplinkunit), 11);
  ASSERT_EQ((extended_apn_ambr.extendedapnambrforuplink), 255);
  ASSERT_EQ((extended_apn_ambr.extendedapnambrforuplink_continued), 255);

  uint8_t extended_apn_ambr_encoded_buffersize19[7] = {0};
  ASSERT_EQ(
      encode_extended_apn_aggregate_maximum_bit_rate(
          &extended_apn_ambr, 0, extended_apn_ambr_encoded_buffersize19,
          sizeof(extended_apn_ambr_encoded_buffersize19)),
      sizeof(extended_apn_ambr_encoded_buffersize19));
  ASSERT_EQ(
      (memcmp(
          extended_apn_ambr_encoded_buffersize19,
          extended_apn_ambr_buffersize19,
          sizeof(extended_apn_ambr_buffersize19))),
      0);

  uint8_t extended_apn_ambr_buffersize20[] = {0x06, 0x0C, 0xC7, 0x46,
                                              0x0C, 0xC7, 0x46};
  ASSERT_EQ(
      decode_extended_apn_aggregate_maximum_bit_rate(
          &extended_apn_ambr, 0, extended_apn_ambr_buffersize20,
          sizeof(extended_apn_ambr_buffersize20)),
      sizeof(extended_apn_ambr_buffersize20));
  ASSERT_EQ((extended_apn_ambr.extendedapnambrfordownlinkunit), 12);
  ASSERT_EQ((extended_apn_ambr.extendedapnambrfordownlink), 199);
  ASSERT_EQ((extended_apn_ambr.extendedapnambrfordownlink_continued), 70);
  ASSERT_EQ((extended_apn_ambr.extendedapnambrforuplinkunit), 12);
  ASSERT_EQ((extended_apn_ambr.extendedapnambrforuplink), 199);
  ASSERT_EQ((extended_apn_ambr.extendedapnambrforuplink_continued), 70);

  uint8_t extended_apn_ambr_encoded_buffersize20[7] = {0};
  ASSERT_EQ(
      encode_extended_apn_aggregate_maximum_bit_rate(
          &extended_apn_ambr, 0, extended_apn_ambr_encoded_buffersize20,
          sizeof(extended_apn_ambr_encoded_buffersize20)),
      sizeof(extended_apn_ambr_encoded_buffersize20));
  ASSERT_EQ(
      (memcmp(
          extended_apn_ambr_encoded_buffersize20,
          extended_apn_ambr_buffersize20,
          sizeof(extended_apn_ambr_buffersize20))),
      0);

  uint8_t extended_apn_ambr_buffersize21[] = {0x06, 0x0C, 0xFF, 0xFF,
                                              0x0C, 0xFF, 0xFF};
  ASSERT_EQ(
      decode_extended_apn_aggregate_maximum_bit_rate(
          &extended_apn_ambr, 0, extended_apn_ambr_buffersize21,
          sizeof(extended_apn_ambr_buffersize21)),
      sizeof(extended_apn_ambr_buffersize21));
  ASSERT_EQ((extended_apn_ambr.extendedapnambrfordownlinkunit), 12);
  ASSERT_EQ((extended_apn_ambr.extendedapnambrfordownlink), 255);
  ASSERT_EQ((extended_apn_ambr.extendedapnambrfordownlink_continued), 255);
  ASSERT_EQ((extended_apn_ambr.extendedapnambrforuplinkunit), 12);
  ASSERT_EQ((extended_apn_ambr.extendedapnambrforuplink), 255);
  ASSERT_EQ((extended_apn_ambr.extendedapnambrforuplink_continued), 255);

  uint8_t extended_apn_ambr_encoded_buffersize21[7] = {0};
  ASSERT_EQ(
      encode_extended_apn_aggregate_maximum_bit_rate(
          &extended_apn_ambr, 0, extended_apn_ambr_encoded_buffersize21,
          sizeof(extended_apn_ambr_encoded_buffersize21)),
      sizeof(extended_apn_ambr_encoded_buffersize21));
  ASSERT_EQ(
      (memcmp(
          extended_apn_ambr_encoded_buffersize21,
          extended_apn_ambr_buffersize21,
          sizeof(extended_apn_ambr_buffersize21))),
      0);

  uint8_t extended_apn_ambr_buffersize22[] = {0x06, 0x0D, 0x75, 0xFB,
                                              0x0D, 0x75, 0xFB};
  ASSERT_EQ(
      decode_extended_apn_aggregate_maximum_bit_rate(
          &extended_apn_ambr, 0, extended_apn_ambr_buffersize22,
          sizeof(extended_apn_ambr_buffersize22)),
      sizeof(extended_apn_ambr_buffersize22));
  ASSERT_EQ((extended_apn_ambr.extendedapnambrfordownlinkunit), 13);
  ASSERT_EQ((extended_apn_ambr.extendedapnambrfordownlink), 117);
  ASSERT_EQ((extended_apn_ambr.extendedapnambrfordownlink_continued), 251);
  ASSERT_EQ((extended_apn_ambr.extendedapnambrforuplinkunit), 13);
  ASSERT_EQ((extended_apn_ambr.extendedapnambrforuplink), 117);
  ASSERT_EQ((extended_apn_ambr.extendedapnambrforuplink_continued), 251);

  uint8_t extended_apn_ambr_encoded_buffersize22[7] = {0};
  ASSERT_EQ(
      encode_extended_apn_aggregate_maximum_bit_rate(
          &extended_apn_ambr, 0, extended_apn_ambr_encoded_buffersize22,
          sizeof(extended_apn_ambr_encoded_buffersize22)),
      sizeof(extended_apn_ambr_encoded_buffersize22));
  ASSERT_EQ(
      (memcmp(
          extended_apn_ambr_encoded_buffersize22,
          extended_apn_ambr_buffersize22,
          sizeof(extended_apn_ambr_buffersize22))),
      0);

  uint8_t extended_apn_ambr_buffersize23[] = {0x06, 0x0D, 0xFF, 0xFF,
                                              0x0D, 0xFF, 0xFF};
  ASSERT_EQ(
      decode_extended_apn_aggregate_maximum_bit_rate(
          &extended_apn_ambr, 0, extended_apn_ambr_buffersize23,
          sizeof(extended_apn_ambr_buffersize23)),
      sizeof(extended_apn_ambr_buffersize23));
  ASSERT_EQ((extended_apn_ambr.extendedapnambrfordownlinkunit), 13);
  ASSERT_EQ((extended_apn_ambr.extendedapnambrfordownlink), 255);
  ASSERT_EQ((extended_apn_ambr.extendedapnambrfordownlink_continued), 255);
  ASSERT_EQ((extended_apn_ambr.extendedapnambrforuplinkunit), 13);
  ASSERT_EQ((extended_apn_ambr.extendedapnambrforuplink), 255);
  ASSERT_EQ((extended_apn_ambr.extendedapnambrforuplink_continued), 255);

  uint8_t extended_apn_ambr_encoded_buffersize23[7] = {0};
  ASSERT_EQ(
      encode_extended_apn_aggregate_maximum_bit_rate(
          &extended_apn_ambr, 0, extended_apn_ambr_encoded_buffersize23,
          sizeof(extended_apn_ambr_encoded_buffersize23)),
      sizeof(extended_apn_ambr_encoded_buffersize23));
  ASSERT_EQ(
      (memcmp(
          extended_apn_ambr_encoded_buffersize23,
          extended_apn_ambr_buffersize23,
          sizeof(extended_apn_ambr_buffersize23))),
      0);

  uint8_t extended_apn_ambr_buffersize24[] = {0x06, 0x0E, 0xD4, 0xE8,
                                              0x0E, 0xD4, 0xE8};
  ASSERT_EQ(
      decode_extended_apn_aggregate_maximum_bit_rate(
          &extended_apn_ambr, 0, extended_apn_ambr_buffersize24,
          sizeof(extended_apn_ambr_buffersize24)),
      sizeof(extended_apn_ambr_buffersize24));
  ASSERT_EQ((extended_apn_ambr.extendedapnambrfordownlinkunit), 14);
  ASSERT_EQ((extended_apn_ambr.extendedapnambrfordownlink), 212);
  ASSERT_EQ((extended_apn_ambr.extendedapnambrfordownlink_continued), 232);
  ASSERT_EQ((extended_apn_ambr.extendedapnambrforuplinkunit), 14);
  ASSERT_EQ((extended_apn_ambr.extendedapnambrforuplink), 212);
  ASSERT_EQ((extended_apn_ambr.extendedapnambrforuplink_continued), 232);

  uint8_t extended_apn_ambr_encoded_buffersize24[7] = {0};
  ASSERT_EQ(
      encode_extended_apn_aggregate_maximum_bit_rate(
          &extended_apn_ambr, 0, extended_apn_ambr_encoded_buffersize24,
          sizeof(extended_apn_ambr_encoded_buffersize24)),
      sizeof(extended_apn_ambr_encoded_buffersize24));
  ASSERT_EQ(
      (memcmp(
          extended_apn_ambr_encoded_buffersize24,
          extended_apn_ambr_buffersize24,
          sizeof(extended_apn_ambr_buffersize24))),
      0);

  uint8_t extended_apn_ambr_buffersize25[] = {0x06, 0x0E, 0xFF, 0xFF,
                                              0x0E, 0xFF, 0xFF};
  ASSERT_EQ(
      decode_extended_apn_aggregate_maximum_bit_rate(
          &extended_apn_ambr, 0, extended_apn_ambr_buffersize25,
          sizeof(extended_apn_ambr_buffersize25)),
      sizeof(extended_apn_ambr_buffersize25));
  ASSERT_EQ((extended_apn_ambr.extendedapnambrfordownlinkunit), 14);
  ASSERT_EQ((extended_apn_ambr.extendedapnambrfordownlink), 255);
  ASSERT_EQ((extended_apn_ambr.extendedapnambrfordownlink_continued), 255);
  ASSERT_EQ((extended_apn_ambr.extendedapnambrforuplinkunit), 14);
  ASSERT_EQ((extended_apn_ambr.extendedapnambrforuplink), 255);
  ASSERT_EQ((extended_apn_ambr.extendedapnambrforuplink_continued), 255);

  uint8_t extended_apn_ambr_encoded_buffersize25[7] = {0};
  ASSERT_EQ(
      encode_extended_apn_aggregate_maximum_bit_rate(
          &extended_apn_ambr, 0, extended_apn_ambr_encoded_buffersize25,
          sizeof(extended_apn_ambr_encoded_buffersize25)),
      sizeof(extended_apn_ambr_encoded_buffersize25));
  ASSERT_EQ(
      (memcmp(
          extended_apn_ambr_encoded_buffersize25,
          extended_apn_ambr_buffersize25,
          sizeof(extended_apn_ambr_buffersize25))),
      0);

  uint8_t extended_apn_ambr_buffersize26[] = {0x06, 0x0F, 0xD4, 0xE8,
                                              0x0F, 0xD4, 0xE8};
  ASSERT_EQ(
      decode_extended_apn_aggregate_maximum_bit_rate(
          &extended_apn_ambr, 0, extended_apn_ambr_buffersize26,
          sizeof(extended_apn_ambr_buffersize26)),
      sizeof(extended_apn_ambr_buffersize26));
  ASSERT_EQ((extended_apn_ambr.extendedapnambrfordownlinkunit), 15);
  ASSERT_EQ((extended_apn_ambr.extendedapnambrfordownlink), 212);
  ASSERT_EQ((extended_apn_ambr.extendedapnambrfordownlink_continued), 232);
  ASSERT_EQ((extended_apn_ambr.extendedapnambrforuplinkunit), 15);
  ASSERT_EQ((extended_apn_ambr.extendedapnambrforuplink), 212);
  ASSERT_EQ((extended_apn_ambr.extendedapnambrforuplink_continued), 232);

  uint8_t extended_apn_ambr_encoded_buffersize26[7] = {0};
  ASSERT_EQ(
      encode_extended_apn_aggregate_maximum_bit_rate(
          &extended_apn_ambr, 0, extended_apn_ambr_encoded_buffersize26,
          sizeof(extended_apn_ambr_encoded_buffersize26)),
      sizeof(extended_apn_ambr_encoded_buffersize26));
  ASSERT_EQ(
      (memcmp(
          extended_apn_ambr_encoded_buffersize26,
          extended_apn_ambr_buffersize26,
          sizeof(extended_apn_ambr_buffersize26))),
      0);

  uint8_t extended_apn_ambr_buffersize27[] = {0x06, 0x0F, 0xFF, 0xFF,
                                              0x0F, 0xFF, 0xFF};
  ASSERT_EQ(
      decode_extended_apn_aggregate_maximum_bit_rate(
          &extended_apn_ambr, 0, extended_apn_ambr_buffersize27,
          sizeof(extended_apn_ambr_buffersize27)),
      sizeof(extended_apn_ambr_buffersize27));
  ASSERT_EQ((extended_apn_ambr.extendedapnambrfordownlinkunit), 15);
  ASSERT_EQ((extended_apn_ambr.extendedapnambrfordownlink), 255);
  ASSERT_EQ((extended_apn_ambr.extendedapnambrfordownlink_continued), 255);
  ASSERT_EQ((extended_apn_ambr.extendedapnambrforuplinkunit), 15);
  ASSERT_EQ((extended_apn_ambr.extendedapnambrforuplink), 255);
  ASSERT_EQ((extended_apn_ambr.extendedapnambrforuplink_continued), 255);

  uint8_t extended_apn_ambr_encoded_buffersize27[7] = {0};
  ASSERT_EQ(
      encode_extended_apn_aggregate_maximum_bit_rate(
          &extended_apn_ambr, 0, extended_apn_ambr_encoded_buffersize27,
          sizeof(extended_apn_ambr_encoded_buffersize27)),
      sizeof(extended_apn_ambr_encoded_buffersize27));
  ASSERT_EQ(
      (memcmp(
          extended_apn_ambr_encoded_buffersize27,
          extended_apn_ambr_buffersize27,
          sizeof(extended_apn_ambr_buffersize27))),
      0);

  uint8_t extended_apn_ambr_buffersize28[] = {0x06, 0x10, 0xDA, 0x47,
                                              0x10, 0xDA, 0x47};
  ASSERT_EQ(
      decode_extended_apn_aggregate_maximum_bit_rate(
          &extended_apn_ambr, 0, extended_apn_ambr_buffersize28,
          sizeof(extended_apn_ambr_buffersize28)),
      sizeof(extended_apn_ambr_buffersize28));
  ASSERT_EQ((extended_apn_ambr.extendedapnambrfordownlinkunit), 16);
  ASSERT_EQ((extended_apn_ambr.extendedapnambrfordownlink), 218);
  ASSERT_EQ((extended_apn_ambr.extendedapnambrfordownlink_continued), 71);
  ASSERT_EQ((extended_apn_ambr.extendedapnambrforuplinkunit), 16);
  ASSERT_EQ((extended_apn_ambr.extendedapnambrforuplink), 218);
  ASSERT_EQ((extended_apn_ambr.extendedapnambrforuplink_continued), 71);

  uint8_t extended_apn_ambr_encoded_buffersize28[7] = {0};
  ASSERT_EQ(
      encode_extended_apn_aggregate_maximum_bit_rate(
          &extended_apn_ambr, 0, extended_apn_ambr_encoded_buffersize28,
          sizeof(extended_apn_ambr_encoded_buffersize28)),
      sizeof(extended_apn_ambr_encoded_buffersize28));
  ASSERT_EQ(
      (memcmp(
          extended_apn_ambr_encoded_buffersize28,
          extended_apn_ambr_buffersize28,
          sizeof(extended_apn_ambr_buffersize28))),
      0);

  uint8_t extended_apn_ambr_buffersize29[] = {0x06, 0x10, 0xFF, 0xFF,
                                              0x10, 0xFF, 0xFF};
  ASSERT_EQ(
      decode_extended_apn_aggregate_maximum_bit_rate(
          &extended_apn_ambr, 0, extended_apn_ambr_buffersize29,
          sizeof(extended_apn_ambr_buffersize29)),
      sizeof(extended_apn_ambr_buffersize29));
  ASSERT_EQ((extended_apn_ambr.extendedapnambrfordownlinkunit), 16);
  ASSERT_EQ((extended_apn_ambr.extendedapnambrfordownlink), 255);
  ASSERT_EQ((extended_apn_ambr.extendedapnambrfordownlink_continued), 255);
  ASSERT_EQ((extended_apn_ambr.extendedapnambrforuplinkunit), 16);
  ASSERT_EQ((extended_apn_ambr.extendedapnambrforuplink), 255);
  ASSERT_EQ((extended_apn_ambr.extendedapnambrforuplink_continued), 255);

  uint8_t extended_apn_ambr_encoded_buffersize29[7] = {0};
  ASSERT_EQ(
      encode_extended_apn_aggregate_maximum_bit_rate(
          &extended_apn_ambr, 0, extended_apn_ambr_encoded_buffersize29,
          sizeof(extended_apn_ambr_encoded_buffersize29)),
      sizeof(extended_apn_ambr_encoded_buffersize29));
  ASSERT_EQ(
      (memcmp(
          extended_apn_ambr_encoded_buffersize29,
          extended_apn_ambr_buffersize29,
          sizeof(extended_apn_ambr_buffersize29))),
      0);
}

}  // namespace lte
}  // namespace magma

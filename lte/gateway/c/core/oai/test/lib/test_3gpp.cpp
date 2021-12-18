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
#include "lte/gateway/c/core/oai/lib/3gpp/3gpp_24.008.h"
#include "lte/gateway/c/core/oai/common/log.h"
#include "lte/gateway/c/core/oai/common/common_defs.h"
#include "lte/gateway/c/core/oai/tasks/nas/ies/MobileIdentity.h"
}

namespace magma {
namespace lte {

#define BUFFER_LEN 200

class m3GppTest : public ::testing::Test {
  virtual void SetUp() {}

  virtual void TearDown() {}

 protected:
  void initialize_imei() {
    imei.typeofidentity = MOBILE_IDENTITY_IMEI;
    imei.oddeven        = MOBILE_IDENTITY_EVEN;
    imei.cdsd           = 0xf;
    imei.tac1           = 1;
    imei.tac2           = 2;
    imei.tac3           = 7;
    imei.tac4           = 8;
    imei.tac5           = 9;
    imei.tac6           = 4;
    imei.tac7           = 5;
    imei.tac8           = 6;
    imei.snr1           = 3;
    imei.snr2           = 4;
    imei.snr3           = 5;
    imei.snr4           = 0;
    imei.snr5           = 6;
    imei.snr6           = 9;
  }

  void cmp_imei(
      const imei_mobile_identity_t* imei1,
      const imei_mobile_identity_t* imei2) {
    EXPECT_EQ(imei1->typeofidentity, imei2->typeofidentity);
    EXPECT_EQ(imei1->oddeven, imei2->oddeven);
    EXPECT_EQ(imei1->tac1, imei2->tac1);
    EXPECT_EQ(imei1->tac2, imei2->tac2);
    EXPECT_EQ(imei1->tac3, imei2->tac3);
    EXPECT_EQ(imei1->tac4, imei2->tac4);
    EXPECT_EQ(imei1->tac5, imei2->tac5);
    EXPECT_EQ(imei1->tac6, imei2->tac6);
    EXPECT_EQ(imei1->tac7, imei2->tac7);
    EXPECT_EQ(imei1->tac8, imei2->tac8);
    EXPECT_EQ(imei1->snr1, imei2->snr1);
    EXPECT_EQ(imei1->snr2, imei2->snr2);
    EXPECT_EQ(imei1->snr3, imei2->snr3);
    EXPECT_EQ(imei1->snr4, imei2->snr4);
    EXPECT_EQ(imei1->snr5, imei2->snr5);
    EXPECT_EQ(imei1->snr6, imei2->snr6);
    EXPECT_EQ(imei1->cdsd, imei2->cdsd);
  }

  void initialize_tmgi() {
    tmgi.typeofidentity          = MOBILE_IDENTITY_TMGI;
    tmgi.spare                   = 0;
    tmgi.mbmssessionidindication = 1;
    tmgi.mccmncindication        = 1;
    tmgi.oddeven                 = 1;
    tmgi.mccdigit1               = 9;
    tmgi.mccdigit2               = 5;
    tmgi.mccdigit3               = 3;
    tmgi.mncdigit1               = 4;
    tmgi.mncdigit2               = 8;
    tmgi.mncdigit3               = 7;
    tmgi.mbmsserviceid           = 0xCADB85;
    tmgi.mbmssessionid           = 12;
  }

  void cmp_tmgi(
      const tmgi_mobile_identity_t* tmgi1,
      const tmgi_mobile_identity_t* tmgi2) {
    EXPECT_EQ(tmgi1->typeofidentity, tmgi2->typeofidentity);
    EXPECT_EQ(tmgi1->spare, tmgi2->spare);
    EXPECT_EQ(tmgi1->mbmssessionidindication, tmgi2->mbmssessionidindication);
    EXPECT_EQ(tmgi1->mbmsserviceid, tmgi2->mbmsserviceid);
    EXPECT_EQ(tmgi1->mbmssessionid, tmgi2->mbmssessionid);
    EXPECT_EQ(tmgi1->oddeven, tmgi2->oddeven);
    EXPECT_EQ(tmgi1->mccdigit1, tmgi2->mccdigit1);
    EXPECT_EQ(tmgi1->mccdigit2, tmgi2->mccdigit2);
    EXPECT_EQ(tmgi1->mccdigit3, tmgi2->mccdigit3);
    EXPECT_EQ(tmgi1->mncdigit1, tmgi2->mncdigit1);
    EXPECT_EQ(tmgi1->mncdigit2, tmgi2->mncdigit2);
    EXPECT_EQ(tmgi1->mncdigit3, tmgi2->mncdigit3);
  }

  imei_mobile_identity_t imei;
  tmgi_mobile_identity_t tmgi;
  uint8_t buffer[BUFFER_LEN];
};

TEST_F(m3GppTest, TestImeiMobileIdentity) {
  imei_mobile_identity_t imei_decoded = {0};
  initialize_imei();

  int encoded = encode_imei_mobile_identity(&imei, buffer, BUFFER_LEN);
  int decoded = decode_imei_mobile_identity(&imei_decoded, buffer, encoded);

  EXPECT_EQ(encoded, decoded);
  cmp_imei(&imei, &imei_decoded);

  imei.oddeven = MOBILE_IDENTITY_ODD;
  imei.cdsd    = 9;

  encoded = encode_imei_mobile_identity(&imei, buffer, BUFFER_LEN);
  decoded = decode_imei_mobile_identity(&imei_decoded, buffer, encoded);

  EXPECT_EQ(encoded, decoded);
  cmp_imei(&imei, &imei_decoded);

  // Test for TLV_VALUE_DOESNT_MATCH
  imei.typeofidentity = 0;
  encoded             = encode_imei_mobile_identity(&imei, buffer, BUFFER_LEN);
  decoded = decode_imei_mobile_identity(&imei_decoded, buffer, encoded);

  EXPECT_EQ(decoded, TLV_VALUE_DOESNT_MATCH);

  // Test for cdsd default decoded to 0x0f then oddeven is set as even.
  imei.typeofidentity = MOBILE_IDENTITY_IMEI;
  imei.cdsd           = 1;
  imei.oddeven        = MOBILE_IDENTITY_EVEN;
  encoded             = encode_imei_mobile_identity(&imei, buffer, BUFFER_LEN);
  decoded = decode_imei_mobile_identity(&imei_decoded, buffer, encoded);

  EXPECT_EQ(encoded, decoded);
  EXPECT_EQ(imei_decoded.cdsd, 0x0f);
}

TEST_F(m3GppTest, TestTmgiMobileIdentity) {
  tmgi_mobile_identity_t tmgi_decoded = {0};
  initialize_tmgi();

  int encoded = encode_tmgi_mobile_identity(&tmgi, buffer, BUFFER_LEN);
  int decoded = decode_tmgi_mobile_identity(&tmgi_decoded, buffer, encoded);

  EXPECT_EQ(encoded, decoded);
  cmp_tmgi(&tmgi, &tmgi_decoded);
}

TEST_F(m3GppTest, TestNoMobileIdentity) {
  no_mobile_identity_t no_id         = {0};
  no_mobile_identity_t no_id_decoded = {0};
  no_id.typeofidentity               = MOBILE_IDENTITY_NOT_AVAILABLE;
  no_id.oddeven                      = 1;
  no_id.digit1                       = 9;
  no_id.digit2                       = 0;
  no_id.digit3                       = 0;
  no_id.digit4                       = 0;
  no_id.digit5                       = 0;

  int encoded = encode_no_mobile_identity(
      &no_id, buffer, MOBILE_IDENTITY_NOT_AVAILABLE_LTE_LENGTH);
  int decoded = decode_no_mobile_identity(&no_id_decoded, buffer, encoded);

  EXPECT_EQ(encoded, decoded);
  EXPECT_EQ(no_id.typeofidentity, no_id_decoded.typeofidentity);
  EXPECT_EQ(no_id.oddeven, no_id_decoded.oddeven);
  EXPECT_EQ(no_id.digit1, no_id_decoded.digit1);
  EXPECT_TRUE(!memcmp(&no_id, &no_id_decoded, sizeof(no_mobile_identity_t)));
}

}  // namespace lte
}  // namespace magma

int main(int argc, char** argv) {
  ::testing::InitGoogleTest(&argc, argv);
  OAILOG_INIT("MME", OAILOG_LEVEL_DEBUG, MAX_LOG_PROTOS);
  return RUN_ALL_TESTS();
}
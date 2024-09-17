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
#include <cstring>

extern "C" {
#include "lte/gateway/c/core/common/common_defs.h"
#include "lte/gateway/c/core/oai/common/log.h"
#include "lte/gateway/c/core/oai/lib/3gpp/3gpp_24.008.h"
}

#include "lte/gateway/c/core/oai/tasks/nas/ies/MobileIdentity.hpp"

namespace magma {
namespace lte {

#define BUFFER_LEN 200

class m3GppTest : public ::testing::Test {
  virtual void SetUp() {}

  virtual void TearDown() {}

 protected:
  void initialize_imei() {
    imei.typeofidentity = MOBILE_IDENTITY_IMEI;
    imei.oddeven = MOBILE_IDENTITY_EVEN;
    imei.cdsd = 0xf;
    imei.tac1 = 1;
    imei.tac2 = 2;
    imei.tac3 = 7;
    imei.tac4 = 8;
    imei.tac5 = 9;
    imei.tac6 = 4;
    imei.tac7 = 5;
    imei.tac8 = 6;
    imei.snr1 = 3;
    imei.snr2 = 4;
    imei.snr3 = 5;
    imei.snr4 = 0;
    imei.snr5 = 6;
    imei.snr6 = 9;
  }

  void cmp_imei(const imei_mobile_identity_t* imei1,
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
    tmgi.typeofidentity = MOBILE_IDENTITY_TMGI;
    tmgi.spare = 0;
    tmgi.mbmssessionidindication = 1;
    tmgi.mccmncindication = 1;
    tmgi.oddeven = 1;
    tmgi.mccdigit1 = 9;
    tmgi.mccdigit2 = 5;
    tmgi.mccdigit3 = 3;
    tmgi.mncdigit1 = 4;
    tmgi.mncdigit2 = 8;
    tmgi.mncdigit3 = 7;
    tmgi.mbmsserviceid = 0xCADB85;
    tmgi.mbmssessionid = 12;
  }

  void cmp_tmgi(const tmgi_mobile_identity_t* tmgi1,
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

  void initialize_imsi() {
    imsi.typeofidentity = MOBILE_IDENTITY_IMSI;
    imsi.digit1 = 3;
    imsi.digit2 = 5;
    imsi.digit3 = 6;
    imsi.digit4 = 8;
    imsi.digit5 = 8;
    imsi.digit6 = 5;
    imsi.digit7 = 1;
    imsi.digit8 = 9;
    imsi.digit9 = 7;
    imsi.digit10 = 2;
    imsi.digit11 = 1;
    imsi.digit12 = 4;
    imsi.digit13 = 0;
    imsi.digit14 = 3;
    imsi.digit15 = 5;
    imsi.numOfValidImsiDigits = 15;
    imsi.oddeven = MOBILE_IDENTITY_ODD;
  }

  void cmp_imsi(const imsi_mobile_identity_t* imsi1,
                const imsi_mobile_identity_t* imsi2) {
    EXPECT_EQ(imsi1->typeofidentity, imsi2->typeofidentity);
    EXPECT_EQ(imsi1->digit1, imsi2->digit1);
    EXPECT_EQ(imsi1->digit2, imsi2->digit2);
    EXPECT_EQ(imsi1->digit3, imsi2->digit3);
    EXPECT_EQ(imsi1->digit4, imsi2->digit4);
    EXPECT_EQ(imsi1->digit5, imsi2->digit5);
    EXPECT_EQ(imsi1->digit6, imsi2->digit6);
    EXPECT_EQ(imsi1->digit7, imsi2->digit7);
    EXPECT_EQ(imsi1->digit8, imsi2->digit8);
    EXPECT_EQ(imsi1->digit9, imsi2->digit9);
    EXPECT_EQ(imsi1->digit10, imsi2->digit10);
    EXPECT_EQ(imsi1->digit11, imsi2->digit11);
    EXPECT_EQ(imsi1->digit12, imsi2->digit12);
    EXPECT_EQ(imsi1->digit13, imsi2->digit13);
    EXPECT_EQ(imsi1->digit14, imsi2->digit14);
    EXPECT_EQ(imsi1->digit15, imsi2->digit15);
    EXPECT_EQ(imsi1->numOfValidImsiDigits, imsi2->numOfValidImsiDigits);
    EXPECT_EQ(imsi1->oddeven, imsi2->oddeven);
  }

  imei_mobile_identity_t imei;
  tmgi_mobile_identity_t tmgi;
  imsi_mobile_identity_t imsi;
  uint8_t buffer[BUFFER_LEN];
};

TEST_F(m3GppTest, TestImeiMobileIdentity) {
  imei_mobile_identity_t imei_decoded;
  initialize_imei();

  int encoded = encode_imei_mobile_identity(&imei, buffer, BUFFER_LEN);
  int decoded = decode_imei_mobile_identity(&imei_decoded, buffer, encoded);

  EXPECT_EQ(encoded, decoded);
  cmp_imei(&imei, &imei_decoded);

  imei.oddeven = MOBILE_IDENTITY_ODD;
  imei.cdsd = 9;

  encoded = encode_imei_mobile_identity(&imei, buffer, BUFFER_LEN);
  decoded = decode_imei_mobile_identity(&imei_decoded, buffer, encoded);

  EXPECT_EQ(encoded, decoded);
  cmp_imei(&imei, &imei_decoded);

  // Test for TLV_VALUE_DOESNT_MATCH
  imei.typeofidentity = 0;
  encoded = encode_imei_mobile_identity(&imei, buffer, BUFFER_LEN);
  decoded = decode_imei_mobile_identity(&imei_decoded, buffer, encoded);

  EXPECT_EQ(decoded, TLV_VALUE_DOESNT_MATCH);

  // Test for cdsd default decoded to 0x0f then oddeven is set as even.
  imei.typeofidentity = MOBILE_IDENTITY_IMEI;
  imei.cdsd = 1;
  imei.oddeven = MOBILE_IDENTITY_EVEN;
  encoded = encode_imei_mobile_identity(&imei, buffer, BUFFER_LEN);
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
  no_mobile_identity_t no_id = {0};
  no_mobile_identity_t no_id_decoded;
  no_id.typeofidentity = MOBILE_IDENTITY_NOT_AVAILABLE;
  no_id.oddeven = 1;
  no_id.digit1 = 9;
  no_id.digit2 = 0;
  no_id.digit3 = 0;
  no_id.digit4 = 0;
  no_id.digit5 = 0;

  int encoded = encode_no_mobile_identity(
      &no_id, buffer, MOBILE_IDENTITY_NOT_AVAILABLE_LTE_LENGTH);
  int decoded = decode_no_mobile_identity(&no_id_decoded, buffer, encoded);

  EXPECT_EQ(encoded, decoded);
  EXPECT_EQ(no_id.typeofidentity, no_id_decoded.typeofidentity);
  EXPECT_EQ(no_id.oddeven, no_id_decoded.oddeven);
  EXPECT_EQ(no_id.digit1, no_id_decoded.digit1);
  EXPECT_EQ(no_id.digit2, no_id_decoded.digit2);
  EXPECT_EQ(no_id.digit3, no_id_decoded.digit3);
  EXPECT_EQ(no_id.digit4, no_id_decoded.digit4);
  EXPECT_EQ(no_id.digit5, no_id_decoded.digit5);
  EXPECT_TRUE(!memcmp(&no_id, &no_id_decoded, sizeof(no_mobile_identity_t)));
}

TEST_F(m3GppTest, TestImsiMobileIdentity) {
  imsi_mobile_identity_t imsi_decoded;
  initialize_imsi();

  int encoded = encode_imsi_mobile_identity(&imsi, buffer, BUFFER_LEN);
  int decoded = decode_imsi_mobile_identity(&imsi_decoded, buffer, encoded);

  EXPECT_EQ(encoded, decoded);
  cmp_imsi(&imsi, &imsi_decoded);

  imsi.oddeven = MOBILE_IDENTITY_EVEN;
  imsi.digit15 = 0x0f;
  encoded = encode_imsi_mobile_identity(&imsi, buffer, BUFFER_LEN);
  decoded = decode_imsi_mobile_identity(&imsi_decoded, buffer, encoded);
  EXPECT_EQ(encoded, decoded);
  cmp_imsi(&imsi, &imsi_decoded);
}

TEST_F(m3GppTest, TestMobileStationClassmark2) {
  mobile_station_classmark2_t msclassmark2 = {0};
  mobile_station_classmark2_t msclassmark2_decoded = {0};

  msclassmark2.revisionlevel = 3;
  msclassmark2.esind = 1;
  msclassmark2.a51 = 1;
  msclassmark2.rfpowercapability = 6;
  msclassmark2.pscapability = 1;
  msclassmark2.ssscreenindicator = 2;
  msclassmark2.smcapability = 1;
  msclassmark2.vbs = 1;
  msclassmark2.vgcs = 1;
  msclassmark2.fc = 1;
  msclassmark2.cm3 = 1;
  msclassmark2.lcsvacap = 1;
  msclassmark2.ucs2 = 1;
  msclassmark2.solsa = 1;
  msclassmark2.cmsp = 1;
  msclassmark2.a52 = 1;
  msclassmark2.a53 = 1;

  // With iei present
  int encoded = encode_mobile_station_classmark_2_ie(&msclassmark2, true,
                                                     buffer, BUFFER_LEN);
  int decoded = decode_mobile_station_classmark_2_ie(&msclassmark2_decoded,
                                                     true, buffer, encoded);
  EXPECT_EQ(encoded, decoded);
  EXPECT_TRUE(
      !(memcmp((const void*)&msclassmark2, (const void*)&msclassmark2_decoded,
               sizeof(mobile_station_classmark2_t))));
  // Without iei present
  encoded = encode_mobile_station_classmark_2_ie(&msclassmark2, false, buffer,
                                                 BUFFER_LEN);
  decoded = decode_mobile_station_classmark_2_ie(&msclassmark2_decoded, false,
                                                 buffer, encoded);
  EXPECT_EQ(encoded, decoded);
  EXPECT_TRUE(
      !(memcmp((const void*)&msclassmark2, (const void*)&msclassmark2_decoded,
               sizeof(mobile_station_classmark2_t))));
}

TEST_F(m3GppTest, TestPlmnList) {
  plmn_list_t plmn_list;
  plmn_list_t plmn_list_decoded;

  plmn_list.num_plmn = PLMN_LIST_IE_MAX_PLMN;
  for (int i = 0; i < PLMN_LIST_IE_MAX_PLMN; ++i) {
    plmn_list.plmn[i].mcc_digit1 = 7;
    plmn_list.plmn[i].mcc_digit2 = 4;
    plmn_list.plmn[i].mcc_digit3 = 3;
    plmn_list.plmn[i].mnc_digit1 = 8;
    plmn_list.plmn[i].mnc_digit2 = 1;
    plmn_list.plmn[i].mnc_digit3 = 6;
  }

  // With iei present
  int encoded = encode_plmn_list_ie(&plmn_list, true, buffer, BUFFER_LEN);
  int decoded = decode_plmn_list_ie(&plmn_list_decoded, true, buffer, encoded);

  EXPECT_EQ(encoded, decoded);
  EXPECT_EQ(plmn_list.num_plmn, plmn_list_decoded.num_plmn);
  EXPECT_TRUE(!(memcmp((const void*)&plmn_list, (const void*)&plmn_list_decoded,
                       sizeof(plmn_list_t))));

  // Without iei present
  encoded = encode_plmn_list_ie(&plmn_list, false, buffer, BUFFER_LEN);
  decoded = decode_plmn_list_ie(&plmn_list_decoded, false, buffer, encoded);

  EXPECT_EQ(encoded, decoded);
  EXPECT_EQ(plmn_list.num_plmn, plmn_list_decoded.num_plmn);
  EXPECT_TRUE(!(memcmp((const void*)&plmn_list, (const void*)&plmn_list_decoded,
                       sizeof(plmn_list_t))));
}

}  // namespace lte
}  // namespace magma

// Note: This is necessary for setting up a log thread (Might be addressed by
// #11736)
int main(int argc, char** argv) {
  ::testing::InitGoogleTest(&argc, argv);
  OAILOG_INIT("MME", OAILOG_LEVEL_DEBUG, MAX_LOG_PROTOS);
  return RUN_ALL_TESTS();
}

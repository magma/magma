/*
Copyright 2020 The Magma Authors.
This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

#pragma once
#include <sstream>
#include <cstdint>
using namespace std;

namespace magma5g {
// 5GS mobile identity information element for type of identity "5G-GUTI" SPEC :
// TS-24501, SEC : 9.11.3.4
class GutiM5GSMobileIdentity {
 public:
  uint8_t iei;
  uint8_t spare : 4;
  uint8_t odd_even : 1;
  uint8_t type_of_identity : 3;
  uint8_t mcc_digit2 : 4;
  uint8_t mcc_digit1 : 4;
  uint8_t mnc_digit3 : 4;
  uint8_t mcc_digit3 : 4;
  uint8_t mnc_digit2 : 4;
  uint8_t mnc_digit1 : 4;
  uint8_t amf_regionid;
  uint8_t amf_setid;
  uint8_t amf_setid1;
  uint8_t amf_pointer : 6;
  uint8_t tmsi1;
  uint8_t tmsi2;
  uint8_t tmsi3;
  uint8_t tmsi4;
#define M5GS_MOBILE_IDENTITY_EVEN 0
#define M5GS_MOBILE_IDENTITY_ODD 1

  GutiM5GSMobileIdentity();
  ~GutiM5GSMobileIdentity();
};

// 5GS mobile identity information element for type of identity or "IMEI" SPEC :
// TS-24501, SEC : 9.11.3.4
class ImeiM5GSMobileIdentity {
 public:
  ImeiM5GSMobileIdentity();
  ~ImeiM5GSMobileIdentity();

  uint8_t identity_digit1 : 4;
  uint8_t odd_even : 1;
  uint8_t type_of_identity : 3;
  uint8_t identity_digit3 : 4;
  uint8_t identity_digit2 : 4;
  uint8_t identity_digit4 : 4;
  uint8_t identity_digit5 : 4;
  uint8_t identity_digit6 : 4;
  uint8_t identity_digit7 : 4;
  uint8_t identity_digit8 : 4;
  uint8_t identity_digit9 : 4;
  uint8_t identity_digit10 : 4;
  uint8_t identity_digit11 : 4;
  uint8_t identity_digit12 : 4;
  uint8_t identity_digit13 : 4;
  uint8_t identity_digit14 : 4;
  uint8_t identity_digit15 : 4;
};

// 5GS mobile identity information element for type of identity or "IMSI" SPEC :
// TS-24501, SEC : 9.11.3.4
class ImsiM5GSMobileIdentity {
 public:
  ImsiM5GSMobileIdentity();
  ~ImsiM5GSMobileIdentity();

  uint8_t spare2 : 1;
  uint8_t supi_format : 3;
  uint8_t spare1 : 1;
  uint8_t type_of_identity : 3;
  uint8_t mcc_digit2 : 4;
  uint8_t mcc_digit1 : 4;
  uint8_t mnc_digit3 : 4;
  uint8_t mcc_digit3 : 4;
  uint8_t mnc_digit2 : 4;
  uint8_t mnc_digit1 : 4;
  uint8_t rout_ind_digit_2 : 4;
  uint8_t rout_ind_digit_1 : 4;
  uint8_t rout_ind_digit_4 : 4;
  uint8_t rout_ind_digit_3 : 4;
  uint8_t spare6 : 1;
  uint8_t spare5 : 1;
  uint8_t spare4 : 1;
  uint8_t spare3 : 1;
  uint8_t protect_schm_id : 4;
  uint8_t home_nw_id;
  uint16_t scheme_len;
#define SCHEME_OUTPUT_MAX_LENGTH 65525
  uint8_t scheme_output[SCHEME_OUTPUT_MAX_LENGTH];
  uint8_t msin_digit1 : 4;
  uint8_t msin_digit2 : 4;
  uint8_t msin_digit3 : 4;
  uint8_t msin_digit4 : 4;
  uint8_t msin_digit5 : 4;
  uint8_t numOfValidImsiDigits : 4;
};

// 5GS mobile identity information element for type of identity or "SUCI" SPEC :
// TS-24501, SEC : 9.11.3.4
class SuciM5GSMobileIdentity {
 public:
  SuciM5GSMobileIdentity();
  ~SuciM5GSMobileIdentity();

  uint8_t spare2 : 1;
  uint8_t supi_format : 3;
  uint8_t spare1 : 1;
  uint8_t type_of_identity : 3;
  std::string suci_nai;  // till end of msg
};

// 5GS mobile identity information element for type of identity or "TMSI" SPEC :
// TS-24501, SEC : 9.11.3.4
class TmsiM5GSMobileIdentity {
 public:
  TmsiM5GSMobileIdentity();
  ~TmsiM5GSMobileIdentity();

  uint8_t spare : 4;
  uint8_t odd_even : 1;
  uint8_t type_of_identity : 3;
  uint8_t amf_setid;
  uint8_t amf_setid1 : 2;
  uint8_t amf_pointer : 6;
  uint8_t m5g_tmsi_1;
  uint8_t m5g_tmsi_2;
  uint8_t m5g_tmsi_3;
  uint8_t m5g_tmsi_4;
};

// M5GSMobileIdentityIe Type, SPEC : TS-24501, SEC : 9.11.3.4
union M5GSMobileIdentityIe {
  M5GSMobileIdentityIe();
  ~M5GSMobileIdentityIe();

  TmsiM5GSMobileIdentity tmsi;
  SuciM5GSMobileIdentity suci;
  ImsiM5GSMobileIdentity imsi;
  ImeiM5GSMobileIdentity imei;
  GutiM5GSMobileIdentity guti;
};

// M5GSMobileIdentityMsg Class
class M5GSMobileIdentityMsg {
 public:
  uint8_t iei;
  uint16_t len;
  uint8_t toi;
#define M5GSMobileIdentityMsg_SUCI 0x05
#define M5GSMobileIdentityMsg_GUTI 0x02
#define M5GSMobileIdentityMsg_IMEI 0x03
#define M5GSMobileIdentityMsg_TMSI 0x04
#define M5GSMobileIdentityMsg_IMSI 0x01
#define MOBILE_IDENTITY_MIN_LENGTH 3
#define MOBILE_IDENTITY_MAX_LENGTH 13
  M5GSMobileIdentityIe mobile_identity;

  M5GSMobileIdentityMsg();
  ~M5GSMobileIdentityMsg();
  int EncodeM5GSMobileIdentityMsg(
      M5GSMobileIdentityMsg* m5gs_mobile_identity, uint8_t iei, uint8_t* buffer,
      uint32_t len);
  int DecodeM5GSMobileIdentityMsg(
      M5GSMobileIdentityMsg* m5gs_mobile_identity, uint8_t iei, uint8_t* buffer,
      uint32_t len);
  int DecodeGutiMobileIdentityMsg(
      GutiM5GSMobileIdentity* guti, uint8_t* buffer, uint8_t ielen);
  int DecodeImeiMobileIdentityMsg(
      ImeiM5GSMobileIdentity* imei, uint8_t* buffer, uint8_t ielen);
  int DecodeImsiMobileIdentityMsg(
      ImsiM5GSMobileIdentity* imsi, uint8_t* buffer, uint8_t ielen);
  int DecodeSuciMobileIdentityMsg(
      SuciM5GSMobileIdentity* suci, uint8_t* buffer, uint8_t ielen);
  int DecodeTmsiMobileIdentityMsg(
      TmsiM5GSMobileIdentity* tmsi, uint8_t* buffer, uint8_t ielen);
  int EncodeGutiMobileIdentityMsg(
      GutiM5GSMobileIdentity* guti, uint8_t* buffer);
  int EncodeImeiMobileIdentityMsg(
      ImeiM5GSMobileIdentity* imei, uint8_t* buffer);
  int EncodeImsiMobileIdentityMsg(
      ImsiM5GSMobileIdentity* imsi, uint8_t* buffer);
  int EncodeSuciMobileIdentityMsg(
      SuciM5GSMobileIdentity* suci, uint8_t* buffer);
  int EncodeTmsiMobileIdentityMsg(
      TmsiM5GSMobileIdentity* tmsi, uint8_t* buffer);
};
}  // namespace magma5g

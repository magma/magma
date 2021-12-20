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

#include <iomanip>
#include <sstream>
#include <cstdint>
#include <cstring>
#include <array>
#include "lte/gateway/c/core/oai/tasks/nas5g/include/ies/M5GSMobileIdentity.h"
#include "lte/gateway/c/core/oai/tasks/nas5g/include/M5GCommonDefs.h"

namespace magma5g {
M5GSMobileIdentityMsg::M5GSMobileIdentityMsg(){};
M5GSMobileIdentityMsg::~M5GSMobileIdentityMsg(){};
GutiM5GSMobileIdentity::GutiM5GSMobileIdentity(){};
GutiM5GSMobileIdentity::~GutiM5GSMobileIdentity(){};
ImeiM5GSMobileIdentity::ImeiM5GSMobileIdentity(){};
ImeiM5GSMobileIdentity::~ImeiM5GSMobileIdentity(){};
ImsiM5GSMobileIdentity::ImsiM5GSMobileIdentity(){};
ImsiM5GSMobileIdentity::~ImsiM5GSMobileIdentity(){};
SuciM5GSMobileIdentity::SuciM5GSMobileIdentity(){};
SuciM5GSMobileIdentity::~SuciM5GSMobileIdentity(){};
TmsiM5GSMobileIdentity::TmsiM5GSMobileIdentity(){};
TmsiM5GSMobileIdentity::~TmsiM5GSMobileIdentity(){};
M5GSMobileIdentityIe::M5GSMobileIdentityIe(){};
M5GSMobileIdentityIe::~M5GSMobileIdentityIe(){};

// Decode GutiMobileIdentity IE Message
int M5GSMobileIdentityMsg::DecodeGutiMobileIdentityMsg(
    GutiM5GSMobileIdentity* guti, uint8_t* buffer, uint8_t ielen) {
  int decoded = 0;
  uint16_t setid;

  MLOG(MDEBUG) << " --- Guti Mobile Identity \n";
  guti->spare = (*(buffer + decoded) >> 4) & 0xf;

  // For the GUTI, bits 5 to 8 of octet 3 are coded as "1111"
  if (guti->spare != 0xf) {
    MLOG(MERROR) << "Error: " << std::dec << TLV_VALUE_DOESNT_MATCH;
    return (TLV_VALUE_DOESNT_MATCH);
  }

  guti->odd_even         = (*(buffer + decoded) >> 3) & 0x1;
  guti->type_of_identity = *(buffer + decoded) & 0x7;

  if (guti->type_of_identity != M5GSMobileIdentityMsg_GUTI) {
    MLOG(MERROR) << "Error: " << std::dec << TLV_VALUE_DOESNT_MATCH;
    return (TLV_VALUE_DOESNT_MATCH);
  }

  decoded++;
  guti->mcc_digit2 = (*(buffer + decoded) >> 4) & 0xf;
  guti->mcc_digit1 = *(buffer + decoded) & 0xf;
  decoded++;
  guti->mnc_digit3 = (*(buffer + decoded) >> 4) & 0xf;
  guti->mcc_digit3 = *(buffer + decoded) & 0xf;
  decoded++;
  guti->mnc_digit2 = (*(buffer + decoded) >> 4) & 0xf;
  guti->mnc_digit1 = *(buffer + decoded) & 0xf;
  decoded++;
  guti->amf_regionid = *(buffer + decoded);
  decoded++;
  setid = *(buffer + decoded);
  decoded++;
  guti->amf_setid =
      0x0000 | ((setid & 0xff) << 2) | ((*(buffer + decoded) >> 6) & 0x3);
  guti->amf_pointer = *(buffer + decoded) & 0x3f;
  decoded++;
  guti->tmsi1 = *(buffer + decoded);
  decoded++;
  guti->tmsi2 = *(buffer + decoded);
  decoded++;
  guti->tmsi3 = *(buffer + decoded);
  decoded++;
  guti->tmsi4 = *(buffer + decoded);
  decoded++;
  MLOG(MDEBUG) << "   Odd/Even Indication = " << std::dec << int(guti->odd_even)
               << "\n";
  MLOG(MDEBUG) << "   Mobile Country Code (MCC) = " << std::dec
               << int(guti->mcc_digit1) << std::dec << int(guti->mcc_digit2)
               << std::dec << int(guti->mcc_digit3) << "\n";
  MLOG(MDEBUG) << "   Mobile Network Code (MNC) = " << std::dec
               << int(guti->mnc_digit1) << std::dec << int(guti->mnc_digit2)
               << std::dec << int(guti->mnc_digit3) << "\n";
  MLOG(MDEBUG) << "   Amf Region ID = " << std::dec << int(guti->amf_regionid)
               << "\n";
  MLOG(MDEBUG) << "   Amf Set ID = " << std::dec << int(guti->amf_setid)
               << "\n";
  MLOG(MDEBUG) << "   Amf Pointer = " << std::dec << int(guti->amf_pointer)
               << "\n";
  MLOG(MDEUBG) << "   M5G-TMSI = "
               << "0x0" << std::hex << int(guti->tmsi1) << "0" << std::hex
               << int(guti->tmsi2) << "0" << std::hex << int(guti->tmsi3) << "0"
               << std::hex << int(guti->tmsi4) << "\n\n";
  return (decoded);
}

// Decode ImeiMobileIdentity IE
int M5GSMobileIdentityMsg::DecodeImeiMobileIdentityMsg(
    ImeiM5GSMobileIdentity* imei, uint8_t* buffer, uint8_t ielen) {
  int decoded = 0;

  MLOG(MDEBUG) << "         DecodeImeiMobileIdentityMsg : "
               << "\n";
  imei->identity_digit1 = (*(buffer + decoded) >> 4) & 0xf;

  if (imei->identity_digit1 != 0xf) {
    MLOG(MERROR) << "Error: " << std::hex << TLV_VALUE_DOESNT_MATCH;
    return (TLV_VALUE_DOESNT_MATCH);
  }

  imei->odd_even         = (*(buffer + decoded) >> 3) & 0x1;
  imei->type_of_identity = *(buffer + decoded) & 0x7;

  if (imei->type_of_identity != M5GSMobileIdentityMsg_IMEI) {
    MLOG(MERROR) << "Error: " << std::dec << TLV_VALUE_DOESNT_MATCH;
    return (TLV_VALUE_DOESNT_MATCH);
  }

  decoded++;
  imei->identity_digit3 = (*(buffer + decoded) >> 4) & 0xf;
  imei->identity_digit2 = *(buffer + decoded) & 0xf;
  decoded++;
  MLOG(MDEBUG) << "  odd_even = " << std::hex << int(imei->odd_even) << "\n";
  MLOG(MDEBUG) << "  digit1 = " << std::hex << int(imei->identity_digit1)
               << "\n";
  MLOG(MDEBUG) << "  digit2 = " << std::hex << int(imei->identity_digit2)
               << "\n";
  MLOG(MDEBUG) << "  digit3 = " << std::hex << int(imei->identity_digit3)
               << "\n";

  return (decoded);
};

// Decode ImsiMobileIdentity IE
int M5GSMobileIdentityMsg::DecodeImsiMobileIdentityMsg(
    ImsiM5GSMobileIdentity* imsi, uint8_t* buffer, uint8_t ielen) {
  int decoded = 0;

  MLOG(MDEBUG) << " DecodeImsiMobileIdentityMsg:";
  imsi->spare2           = (*(buffer + decoded) >> 7) & 0x1;
  imsi->supi_format      = (*(buffer + decoded) >> 4) & 0x7;
  imsi->spare1           = (*(buffer + decoded) >> 3) & 0x1;
  imsi->type_of_identity = *(buffer + decoded) & 0x7;

  if (imsi->type_of_identity != M5GSMobileIdentityMsg_SUCI_IMSI) {
    MLOG(MERROR) << "Error: " << std::hex << TLV_VALUE_DOESNT_MATCH;
    return (TLV_VALUE_DOESNT_MATCH);
  }

  decoded++;
  imsi->mcc_digit2 = (*(buffer + decoded) >> 4) & 0xf;
  imsi->mcc_digit1 = *(buffer + decoded) & 0xf;
  decoded++;
  imsi->mnc_digit3 = (*(buffer + decoded) >> 4) & 0xf;
  imsi->mcc_digit3 = *(buffer + decoded) & 0xf;
  decoded++;
  imsi->mnc_digit2 = (*(buffer + decoded) >> 4) & 0xf;
  imsi->mnc_digit1 = *(buffer + decoded) & 0xf;
  decoded++;
  imsi->rout_ind_digit_2 = (*(buffer + decoded) >> 4) & 0xf;
  imsi->rout_ind_digit_1 = *(buffer + decoded) & 0xf;
  decoded++;
  imsi->rout_ind_digit_4 = (*(buffer + decoded) >> 4) & 0xf;
  imsi->rout_ind_digit_3 = *(buffer + decoded) & 0xf;
  decoded++;
  // TODO
  /* Routing indicator (octets 8-9)
     Routing Indicator shall consist of 1 to 4 digits. The coding of this field
     is the responsibility of home network operator but BCD coding shall be
     used. If a network operator decides to assign less than 4 digits to Routing
     Indicator, the remaining digits shall be coded as "1111" and inserted at
     the left side to fill the 4 digits coding of Routing Indicator. If no
     Routing Indicator is configured in the USIM, the UE shall code the first
     digit of the Routing Indicator as "0000" and the remaining digits as
     “1111".
  */
  imsi->spare6          = (*(buffer + decoded) >> 7) & 0x1;
  imsi->spare5          = (*(buffer + decoded) >> 6) & 0x1;
  imsi->spare4          = (*(buffer + decoded) >> 5) & 0x1;
  imsi->spare3          = (*(buffer + decoded) >> 4) & 0x1;
  imsi->protect_schm_id = *(buffer + decoded) & 0xf;
  decoded++;
  imsi->home_nw_id = *(buffer + decoded);
  decoded++;

  MLOG(MDEBUG) << "Length :  " << int(ielen);
  memcpy(&imsi->scheme_output, buffer + decoded, ielen - decoded);

  if ((imsi->type_of_identity == M5GSMobileIdentityMsg_SUCI_IMSI) &&
      (imsi->protect_schm_id != 0)) {
    MLOG(MDEBUG) << "  SUCI Registration is enabled with protect_schm_id= "
                 << std::hex << int(imsi->protect_schm_id);
    if (imsi->protect_schm_id == PROFILE_A) {
      memcpy(
          &imsi->empheral_public_key, buffer + decoded,
          EPHEMERAL_PUBLIC_KEY_LENGTH);
      decoded += EPHEMERAL_PUBLIC_KEY_LENGTH;
      imsi->empheral_public_key[EPHEMERAL_PUBLIC_KEY_LENGTH] = '\0';
    } else {
      memcpy(
          &imsi->empheral_public_key, buffer + decoded,
          EPHEMERAL_PUBLIC_KEY_LENGTH + PROFILE_B_LEN);
      decoded += (EPHEMERAL_PUBLIC_KEY_LENGTH + PROFILE_B_LEN);
      imsi->empheral_public_key[EPHEMERAL_PUBLIC_KEY_LENGTH + PROFILE_B_LEN] =
          '\0';
    }

    memcpy(&imsi->ciphertext, buffer + decoded, CIPHERTEXT_LENGTH);
    decoded += CIPHERTEXT_LENGTH;
    imsi->ciphertext[CIPHERTEXT_LENGTH] = '\0';

    memcpy(&imsi->mac_tag, buffer + decoded, MAC_TAG_LENGTH);
    decoded += MAC_TAG_LENGTH;
    imsi->mac_tag[MAC_TAG_LENGTH] = '\0';
  }

  // AMF_TEST scheme output  nibbles needs to be reversed
  REV_NIBBLE(imsi->scheme_output, 5);

  // TODO
  /* Scheme output (octets 12 to x)
     The Scheme output field consists of a string of characters with a variable
     length or std::hexadecimal digits as specified in 3GPP TS 23.003 [4]. If
     Protection scheme identifier is set to "0000" (i.e. Null scheme), then the
     Scheme output consists of the MSIN and is coded using BCD coding with each
     digit of the MSIN coded over 4 bits. If the MSIN includes an odd number of
     digits, bits 5 to 8 of octet x shall be coded as "1111". If Protection
     scheme identifier is not "0000" (i.e. ECIES scheme profile A, ECIES scheme
     profile B or Operator-specific protection scheme), then Scheme output is
     coded as std::hexadecimal digits
  */

  int tmp = ielen - decoded;
  decoded = ielen;
  return (decoded);
};

// Will be supported POST MVC
// Decode SuciMobileIdentity IE
int M5GSMobileIdentityMsg::DecodeSuciMobileIdentityMsg(
    SuciM5GSMobileIdentity* suci, uint8_t* buffer, uint8_t ielen) {
  int decoded = 0;

  MLOG(MDEBUG) << "         DecodeSuciMobileIdentityMsg:"
               << "\n";
  suci->spare2           = (*(buffer + decoded) >> 7) & 0x1;
  suci->supi_format      = (*(buffer + decoded) >> 4) & 0x7;
  suci->spare1           = (*(buffer + decoded) >> 3) & 0x1;
  suci->type_of_identity = *(buffer + decoded) & 0x7;

  if (suci->type_of_identity != M5GSMobileIdentityMsg_IMEISV) {
    MLOG(MDEBUG) << "TLV_VALUE_DOESNT_MATCH error";
    return (TLV_VALUE_DOESNT_MATCH);
  }
  decoded++;

  // Will be supported POST MVC
  suci->suci_nai = *(buffer + decoded);
  decoded++;
  decoded++;

  return (decoded);
};

// Decode TmsiMobileIdentity IE
int M5GSMobileIdentityMsg::DecodeTmsiMobileIdentityMsg(
    TmsiM5GSMobileIdentity* tmsi, uint8_t* buffer, uint8_t ielen) {
  int decoded = 0;

  MLOG(MDEBUG) << "         DecodeTmsiMobileIdentityMsg:"
               << "\n";
  tmsi->spare = (*(buffer + decoded) >> 4) & 0xf;

  if (tmsi->spare != 0xf) {
    MLOG(MDEBUG) << "Error: " << int(TLV_VALUE_DOESNT_MATCH);
    return (TLV_VALUE_DOESNT_MATCH);
  }

  tmsi->odd_even         = (*(buffer + decoded) >> 3) & 0x1;
  tmsi->type_of_identity = *(buffer + decoded) & 0x7;

  if (tmsi->type_of_identity != M5GSMobileIdentityMsg_TMSI) {
    MLOG(MDEBUG) << "Error: " << int(TLV_VALUE_DOESNT_MATCH);
    return (TLV_VALUE_DOESNT_MATCH);
  }
  decoded++;
  uint8_t setid;
  setid = *(buffer + decoded);
  decoded++;
  tmsi->amf_setid =
      0x0000 | ((setid & 0xff) << 2) | ((*(buffer + decoded) >> 6) & 0x3);
  tmsi->amf_pointer = *(buffer + decoded) & 0x3f;
  decoded++;
  memcpy(&tmsi->m5g_tmsi, buffer + decoded, ielen - decoded);
#if 0
  tmsi->m5g_tmsi_1 = *(buffer + decoded);
  decoded++;
  tmsi->m5g_tmsi_2 = *(buffer + decoded);
  decoded++;
  tmsi->m5g_tmsi_3 = *(buffer + decoded);
  decoded++;
  tmsi->m5g_tmsi_4 = *(buffer + decoded);
  decoded++;
#endif
  int tmp = ielen - decoded;
  decoded = ielen;
  MLOG(MDEBUG) << "  spare2 = " << std::dec << int(tmsi->spare);
  MLOG(MDEBUG) << "  odd_even = " << std::dec << int(tmsi->odd_even);
  MLOG(MDEBUG) << "  type_of_identity = " << std::dec
               << int(tmsi->type_of_identity);
  MLOG(MDEBUG) << "  amf_setid = " << std::dec << int(tmsi->amf_setid);
  MLOG(MDEBUG) << "  amf_pointer = " << std::dec << int(tmsi->amf_pointer);
  MLOG(MDEBUG) << "  M5G TMSI = ";
  BUFFER_PRINT_LOG(tmsi->m5g_tmsi, tmp)
#if 0
  MLOG(MDEBUG) << "  m5g_tmsi_1 = " << std::dec << int(tmsi->m5g_tmsi_1);
  MLOG(MDEBUG) << "  m5g_tmsi_2 = " << std::dec << int(tmsi->m5g_tmsi_2);
  MLOG(MDEBUG) << "  m5g_tmsi_3 = " << std::dec << int(tmsi->m5g_tmsi_3);
  MLOG(MDEBUG) << "  m5g_tmsi_4 = " << std::dec << int(tmsi->m5g_tmsi_4);
#endif
  return (decoded);
};

// Decode M5GSMobileIdentity IE
int M5GSMobileIdentityMsg::DecodeM5GSMobileIdentityMsg(
    M5GSMobileIdentityMsg* mg5smobile_identity, uint8_t iei, uint8_t* buffer,
    uint32_t len) {
  int decoded_rc = TLV_VALUE_DOESNT_MATCH;
  int decoded    = 0;
  uint16_t ielen = 0;

  MLOG(MDEBUG) << "M5GS Mobile Identity : ";
  if (iei > 0) {
    CHECK_IEI_DECODER(iei, (unsigned char) *buffer);
    decoded++;
  }

  IES_DECODE_U16(buffer, decoded, ielen);
  CHECK_LENGTH_DECODER(len - decoded, ielen);
  unsigned char type_of_identity = *(buffer + decoded) & 0x7;
  MLOG(MDEBUG) << " Length = " << std::dec << int(ielen)
               << " Type of Identity = " << std::dec << int(type_of_identity);

  if (type_of_identity == M5GSMobileIdentityMsg_IMEISV) {
    MLOG(MDEBUG) << " Type suci";
    decoded_rc = DecodeSuciMobileIdentityMsg(
        &mg5smobile_identity->mobile_identity.suci, buffer, ielen);
  } else if (type_of_identity == M5GSMobileIdentityMsg_GUTI) {
    MLOG(MDEBUG) << " Type guti";
    decoded_rc = DecodeGutiMobileIdentityMsg(
        &mg5smobile_identity->mobile_identity.guti, buffer + decoded, ielen);
  } else if (type_of_identity == M5GSMobileIdentityMsg_IMEI) {
    MLOG(MDEBUG) << " Type imei";
    decoded_rc = DecodeImeiMobileIdentityMsg(
        &mg5smobile_identity->mobile_identity.imei, buffer + decoded, ielen);
  } else if (type_of_identity == M5GSMobileIdentityMsg_TMSI) {
    MLOG(MDEBUG) << " Type tmsi";
    decoded_rc = DecodeTmsiMobileIdentityMsg(
        &mg5smobile_identity->mobile_identity.tmsi, buffer + decoded, ielen);
  } else if (type_of_identity == M5GSMobileIdentityMsg_SUCI_IMSI) {
    MLOG(MDEBUG) << " Type imsi";
    decoded_rc = DecodeImsiMobileIdentityMsg(
        &mg5smobile_identity->mobile_identity.imsi, buffer + decoded, ielen);
  }

  if (decoded_rc < 0) {
    MLOG(MERROR) << "Decode Error";
    return decoded_rc;
  }
  return (decoded + decoded_rc);
};

// Encode GutiMobileIdentity IE
int M5GSMobileIdentityMsg::EncodeGutiMobileIdentityMsg(
    GutiM5GSMobileIdentity* guti, uint8_t* buffer) {
  uint32_t encoded = 0;

  MLOG(MDEBUG) << "EncodeGutiMobileIdentityMsg:";
  *(buffer + encoded) =
      0xf0 | ((guti->odd_even & 0x01) << 3) | (guti->type_of_identity & 0x7);
  MLOG(MDEBUG) << "odd_even type_of_identity = " << std::hex
               << int(*(buffer + encoded));
  encoded++;
  *(buffer + encoded) =
      0x00 | ((guti->mcc_digit2 & 0x0f) << 4) | (guti->mcc_digit1 & 0x0f);
  MLOG(MDEBUG) << "mcc_digit2 >mcc_digit1 type_of_identity = " << std::hex
               << int(*(buffer + encoded));
  encoded++;
  *(buffer + encoded) =
      0x00 | ((guti->mnc_digit3 & 0x0f) << 4) | (guti->mcc_digit3 & 0x0f);
  MLOG(MDEBUG) << "mnc_digit3 >mcc_digit3 type_of_identity = " << std::hex
               << int(*(buffer + encoded));
  encoded++;
  *(buffer + encoded) =
      0x00 | ((guti->mnc_digit2 & 0x0f) << 4) | (guti->mnc_digit1 & 0x0f);
  MLOG(MDEBUG) << "mnc_digit2 >mcc_digit1 type_of_identity = " << std::hex
               << int(*(buffer + encoded));
  encoded++;
  *(buffer + encoded) = 0x00 | guti->amf_regionid;
  MLOG(MDEBUG) << "amf_regionid = " << std::hex << int(*(buffer + encoded));
  encoded++;
  *(buffer + encoded) = 0x00 | ((guti->amf_setid >> 2) & 0xFF);
  MLOG(MDEBUG) << "amf_setid = " << std::hex << int(*(buffer + encoded));
  encoded++;
  *(buffer + encoded) =
      0x00 | ((guti->amf_setid & 0xF3) << 6) | (guti->amf_pointer & 0x3f);
  MLOG(MDEBUG) << "amf_setid amf_pointer = " << std::hex
               << int(*(buffer + encoded));
  encoded++;
  *(buffer + encoded) = 0x00 | guti->tmsi1;
  MLOG(MDEBUG) << "tmsi1 = " << std::hex << int(*(buffer + encoded));
  encoded++;
  *(buffer + encoded) = 0x00 | guti->tmsi2;
  MLOG(MDEBUG) << "tmsi2 = " << std::hex << int(*(buffer + encoded));
  encoded++;
  *(buffer + encoded) = 0x00 | guti->tmsi3;
  MLOG(MDEBUG) << "tmsi3 = " << std::hex << int(*(buffer + encoded));
  encoded++;
  *(buffer + encoded) = 0x00 | guti->tmsi4;
  MLOG(MDEBUG) << "tmsi4 = " << std::hex << int(*(buffer + encoded));
  encoded++;

  return encoded;
};

// Will be supported POST MVC
// Encode ImeiMobileIdentity IE
int M5GSMobileIdentityMsg::EncodeImeiMobileIdentityMsg(
    ImeiM5GSMobileIdentity* imei, uint8_t* buffer) {
  uint32_t encoded = 0;

  MLOG(MDEBUG) << "EncodeImeiMobileIdentityMsg:";
  *(buffer + encoded) = 0x00 | ((imei->identity_digit1 & 0xf0) << 4) |
                        ((imei->odd_even & 0x1) << 3) |
                        (imei->type_of_identity & 0x7);
  MLOG(MDEBUG) << "identity_digit1, odd_even, type_of_identity = " << std::hex
               << int(*(buffer + encoded));
  encoded++;
  *(buffer + encoded) = 0x00 | ((imei->identity_digit2 & 0xf0) << 4) |
                        (imei->identity_digit3 & 0x0f);
  MLOG(MDEBUG) << "identity_digit2,identity_digit3 = " << std::hex
               << int(*(buffer + encoded));
  encoded++;

  return encoded;
};

// Will be supported POST MVC
// Encode ImsiMobileIdentity IE
int M5GSMobileIdentityMsg::EncodeImsiMobileIdentityMsg(
    ImsiM5GSMobileIdentity* imsi, uint8_t* buffer) {
  uint32_t encoded = 0;
  MLOG(MDEBUG) << "EncodeImsiMobileIdentityMsg:";
  *(buffer + encoded) =
      0x00 | ((imsi->spare2 & 0x80) << 7) | ((imsi->supi_format & 0x07) << 4) |
      ((imsi->spare1 & 0x01) << 3) | (imsi->type_of_identity & 0x7);
  MLOG(MDEBUG) << "  Spare,supi_format,spare1,type_of_identity = " << std::hex
               << int(*(buffer + encoded));
  encoded++;
  *(buffer + encoded) =
      0x00 | ((imsi->mcc_digit2 & 0x0f) << 4) | (imsi->mcc_digit1 & 0x0f);
  MLOG(MDEBUG) << "  mcc_digit2,mcc_digit1 = " << std::hex
               << int(*(buffer + encoded));
  encoded++;
  *(buffer + encoded) =
      0x00 | ((imsi->mnc_digit3 & 0x0f) << 4) | (imsi->mcc_digit3 & 0x0f);
  MLOG(MDEBUG) << "  mnc_digit3,mcc_digit3 = " << std::hex
               << int(*(buffer + encoded));
  encoded++;
  *(buffer + encoded) =
      0x00 | ((imsi->mnc_digit2 & 0x0f) << 4) | (imsi->mnc_digit1 & 0x0f);
  MLOG(MDEBUG) << "  mnc_digit2, mnc_digit1 = " << std::hex
               << int(*(buffer + encoded));
  encoded++;
  *(buffer + encoded) =
      0x00 | ((imsi->rout_ind_digit_2) << 4) | (imsi->rout_ind_digit_1);
  MLOG(MDEBUG) << "  rout_ind_digit_2,rout_ind_digit_1 = " << std::hex
               << int(*(buffer + encoded));
  encoded++;
  *(buffer + encoded) =
      0x00 | ((imsi->rout_ind_digit_3) << 4) | (imsi->rout_ind_digit_4);
  MLOG(MDEBUG) << "  rout_ind_digit_3,rout_ind_digit_4 = " << std::hex
               << int(*(buffer + encoded));
  encoded++;
  *(buffer + encoded) =
      0x00 | ((imsi->spare6 & 0x01) << 7) | ((imsi->spare5 & 0x01) << 6) |
      ((imsi->spare4 & 0x01) << 5) | ((imsi->spare3 & 0x01) << 4) |
      (imsi->protect_schm_id & 0x0f);
  *(buffer + encoded) = imsi->home_nw_id;
  MLOG(MDEBUG) << "  spare,protect_schm_id = " << std::hex
               << int(*(buffer + encoded));
  encoded++;
  // Will be supported POST MVC

  memcpy(buffer + encoded, &imsi->scheme_output, imsi->scheme_len);
  MLOG(MDEBUG) << "  Scheme Output = ";
  BUFFER_PRINT_LOG(imsi->scheme_output, imsi->scheme_len);
  encoded = encoded + imsi->scheme_len;
  MLOG(MDEBUG) << std::endl;

  return encoded;
};

// Will be supported POST MVC
int M5GSMobileIdentityMsg::EncodeTmsiMobileIdentityMsg(
    TmsiM5GSMobileIdentity* tmsi, uint8_t* buffer) {
  uint32_t encoded = 0;

  MLOG(MDEBUG) << "EncodeTmsiMobileIdentityMsg:";
  *(buffer + encoded) = 0x00 | ((tmsi->spare & 0x0f) << 4) |
                        ((tmsi->odd_even & 0x01) << 3) |
                        (tmsi->type_of_identity & 0x7);
  encoded++;
  *(buffer + encoded) = 0x00 | tmsi->amf_setid;
  encoded++;
  *(buffer + encoded) = 0x00 | ((tmsi->amf_setid & 0xc0) << 6);
  *(buffer + encoded) = 0x00 | (tmsi->amf_pointer & 0x3f);
  encoded++;
#if 0
  *(buffer + encoded) = 0x00 | tmsi->m5g_tmsi_1;
  encoded++;
  *(buffer + encoded) = 0x00 | tmsi->m5g_tmsi_2;
  encoded++;
  *(buffer + encoded) = 0x00 | tmsi->m5g_tmsi_3;
  encoded++;
  *(buffer + encoded) = 0x00 | tmsi->m5g_tmsi_4;
  encoded++;
#endif
  return encoded;
};

// Encode SuciMobileIdentity IE
int M5GSMobileIdentityMsg::EncodeSuciMobileIdentityMsg(
    SuciM5GSMobileIdentity* suci, uint8_t* buffer) {
  uint32_t encoded = 0;

  MLOG(MDEBUG) << "EncodeSuciMobileIdentityMsg:";
  *(buffer + encoded) =
      0x00 | ((suci->spare2 & 0x80) << 7) | ((suci->supi_format & 0x07) << 4) |
      ((suci->spare1 & 0x01) << 3) | (suci->type_of_identity & 0x7);
  encoded++;
  suci->suci_nai.assign(
      (const char*) (buffer + encoded), suci->suci_nai.size());
  MLOG(MDEBUG) << "ielen = " << std::hex
               << (unsigned char) suci->suci_nai.size();
  MLOG(MDEBUG) << "contents";
  for (uint32_t i = 0; i < suci->suci_nai.size(); i++) {
    MLOG(MDEBUG) << std::hex << int(suci->suci_nai[i]);
  }
  MLOG(MDEBUG) << std::endl;

  return encoded;
};

// Encode M5GSMobileIdentity IE
int M5GSMobileIdentityMsg::EncodeM5GSMobileIdentityMsg(
    M5GSMobileIdentityMsg* m5gs_mobile_identity, uint8_t iei, uint8_t* buffer,
    uint32_t len) {
  uint16_t* lenPtr;
  int encoded_rc   = TLV_VALUE_DOESNT_MATCH;
  uint32_t encoded = 0;

  // Checking IEI and pointer
  CHECK_PDU_POINTER_AND_LENGTH_ENCODER(buffer, MOBILE_IDENTITY_MIN_LENGTH, len);

  if (m5gs_mobile_identity->iei > 0) {
    MLOG(MDEBUG) << "EncodeM5GSMobileIdentityMsg:";
    CHECK_IEI_ENCODER((unsigned char) iei, m5gs_mobile_identity->iei);
    *buffer = iei;
    MLOG(MDEBUG) << "iei" << std::hex << int(*buffer);
    encoded++;
  } else
    return 0;

  lenPtr = (uint16_t*) (buffer + encoded);
  encoded += 2;
  m5gs_mobile_identity->toi =
      m5gs_mobile_identity->mobile_identity.guti.type_of_identity;
  if (m5gs_mobile_identity->toi == M5GSMobileIdentityMsg_SUCI_IMSI) {
    MLOG(MDEBUG) << "Type imsi";
    encoded_rc = EncodeImsiMobileIdentityMsg(
        &m5gs_mobile_identity->mobile_identity.imsi, buffer + encoded);
  } else if (m5gs_mobile_identity->toi == M5GSMobileIdentityMsg_IMEI) {
    MLOG(MDEBUG) << "Type imei";
    encoded_rc = EncodeImeiMobileIdentityMsg(
        &m5gs_mobile_identity->mobile_identity.imei, buffer + encoded);
  } else if (m5gs_mobile_identity->toi == M5GSMobileIdentityMsg_GUTI) {
    MLOG(MDEBUG) << "Type guti";
    encoded_rc = EncodeGutiMobileIdentityMsg(
        &m5gs_mobile_identity->mobile_identity.guti, buffer + encoded);
  } else if (m5gs_mobile_identity->toi == M5GSMobileIdentityMsg_TMSI) {
    MLOG(MDEBUG) << "Type tmsi";
    encoded_rc = EncodeTmsiMobileIdentityMsg(
        &m5gs_mobile_identity->mobile_identity.tmsi, buffer + encoded);
  } else if (m5gs_mobile_identity->toi == M5GSMobileIdentityMsg_IMEISV) {
    MLOG(MDEBUG) << "Type suci";
    encoded_rc = EncodeSuciMobileIdentityMsg(
        &m5gs_mobile_identity->mobile_identity.suci, buffer + encoded);
  }

  if (encoded_rc < 0) {
    MLOG(MDEBUG) << "Encode error" << encoded_rc;
    return encoded_rc;
  }

  *lenPtr = htons(encoded + encoded_rc - 2 - ((iei > 0) ? 1 : 0));
  return (encoded + encoded_rc);
};
}  // namespace magma5g

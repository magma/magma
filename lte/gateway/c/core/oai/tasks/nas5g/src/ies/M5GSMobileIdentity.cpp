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
#ifdef __cplusplus
extern "C" {
#endif
#include "lte/gateway/c/core/oai/common/log.h"
#ifdef __cplusplus
}
#endif
#include "lte/gateway/c/core/oai/tasks/nas5g/include/ies/M5GSMobileIdentity.hpp"
#include "lte/gateway/c/core/oai/tasks/nas5g/include/M5GCommonDefs.h"

namespace magma5g {
M5GSMobileIdentityMsg::M5GSMobileIdentityMsg() {};
M5GSMobileIdentityMsg::~M5GSMobileIdentityMsg() {};
GutiM5GSMobileIdentity::GutiM5GSMobileIdentity() {};
GutiM5GSMobileIdentity::~GutiM5GSMobileIdentity() {};
ImeiM5GSMobileIdentity::ImeiM5GSMobileIdentity() {};
ImeiM5GSMobileIdentity::~ImeiM5GSMobileIdentity() {};
ImsiM5GSMobileIdentity::ImsiM5GSMobileIdentity() {};
ImsiM5GSMobileIdentity::~ImsiM5GSMobileIdentity() {};
SuciM5GSMobileIdentity::SuciM5GSMobileIdentity() {};
SuciM5GSMobileIdentity::~SuciM5GSMobileIdentity() {};
TmsiM5GSMobileIdentity::TmsiM5GSMobileIdentity() {};
TmsiM5GSMobileIdentity::~TmsiM5GSMobileIdentity() {};
M5GSMobileIdentityIe::M5GSMobileIdentityIe() {};
M5GSMobileIdentityIe::~M5GSMobileIdentityIe() {};

// Decode GutiMobileIdentity IE Message
int M5GSMobileIdentityMsg::DecodeGutiMobileIdentityMsg(
    GutiM5GSMobileIdentity* guti, uint8_t* buffer, uint8_t ielen) {
  int decoded = 0;
  uint16_t setid;

  guti->spare = (*(buffer + decoded) >> 4) & 0xf;

  // For the GUTI, bits 5 to 8 of octet 3 are coded as "1111"
  if (guti->spare != 0xf) {
    OAILOG_ERROR(LOG_NAS5G, "Error: %d", TLV_VALUE_DOESNT_MATCH);
    return (TLV_VALUE_DOESNT_MATCH);
  }

  guti->odd_even = (*(buffer + decoded) >> 3) & 0x1;
  guti->type_of_identity = *(buffer + decoded) & 0x7;

  if (guti->type_of_identity != M5GSMobileIdentityMsg_GUTI) {
    OAILOG_ERROR(LOG_NAS5G, "Error: %d", TLV_VALUE_DOESNT_MATCH);
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
  return (decoded);
}

// Decode ImeiMobileIdentity IE
int M5GSMobileIdentityMsg::DecodeImeiMobileIdentityMsg(
    ImeiM5GSMobileIdentity* imei, uint8_t* buffer, uint8_t ielen) {
  int decoded = 0;

  imei->identity_digit1 = (*(buffer + decoded) >> 4) & 0xf;

  if (imei->identity_digit1 != 0xf) {
    OAILOG_ERROR(LOG_NAS5G, "Error : %d", TLV_VALUE_DOESNT_MATCH);
    return (TLV_VALUE_DOESNT_MATCH);
  }

  imei->odd_even = (*(buffer + decoded) >> 3) & 0x1;
  imei->type_of_identity = *(buffer + decoded) & 0x7;

  if (imei->type_of_identity != M5GSMobileIdentityMsg_IMEI) {
    OAILOG_ERROR(LOG_NAS5G, "Error : %d", TLV_VALUE_DOESNT_MATCH);
    return (TLV_VALUE_DOESNT_MATCH);
  }

  decoded++;
  imei->identity_digit3 = (*(buffer + decoded) >> 4) & 0xf;
  imei->identity_digit2 = *(buffer + decoded) & 0xf;
  decoded++;

  return (decoded);
};

// Decode ImsiMobileIdentity IE
int M5GSMobileIdentityMsg::DecodeImsiMobileIdentityMsg(
    ImsiM5GSMobileIdentity* imsi, uint8_t* buffer, uint8_t ielen) {
  int decoded = 0;
  /* 5GS mobile identity comprises of
     1  byte  for spare bits, supi format and type of identity
     3  bytes for mcc and mnc length
     2  bytes for Routing indicator
     1  byte  for Protection scheme id
     1  byte  for Home network id
     32 bytes for EPHEMERAL PUBLIC KEY LENGTH for ProfileA
     or
     33 bytes for EPHEMERAL PUBLIC KEY LENGTH for ProfileB
     *  variable bytes for ciphertext
     8  bytes for MAC TAG LENGTH */

  int cipherTextLen = 0;
  imsi->spare2 = (*(buffer + decoded) >> 7) & 0x1;
  imsi->supi_format = (*(buffer + decoded) >> 4) & 0x7;
  imsi->spare1 = (*(buffer + decoded) >> 3) & 0x1;
  imsi->type_of_identity = *(buffer + decoded) & 0x7;

  if (imsi->type_of_identity != M5GSMobileIdentityMsg_SUCI_IMSI) {
    OAILOG_ERROR(LOG_NAS5G, "Error : %d", TLV_VALUE_DOESNT_MATCH);
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
  imsi->spare6 = (*(buffer + decoded) >> 7) & 0x1;
  imsi->spare5 = (*(buffer + decoded) >> 6) & 0x1;
  imsi->spare4 = (*(buffer + decoded) >> 5) & 0x1;
  imsi->spare3 = (*(buffer + decoded) >> 4) & 0x1;
  imsi->protect_schm_id = *(buffer + decoded) & 0xf;
  decoded++;
  imsi->home_nw_id = *(buffer + decoded);
  decoded++;

  memcpy(&imsi->scheme_output, buffer + decoded, ielen - decoded);

  if ((imsi->type_of_identity == M5GSMobileIdentityMsg_SUCI_IMSI) &&
      (imsi->protect_schm_id != 0)) {
    OAILOG_DEBUG(LOG_NAS5G,
                 "SUCI Registration is enabled with protect_schm_id : %X",
                 static_cast<int>(imsi->protect_schm_id));
    if (imsi->protect_schm_id == PROFILE_A) {
      memcpy(&imsi->empheral_public_key, buffer + decoded,
             EPHEMERAL_PUBLIC_KEY_LENGTH);
      decoded += EPHEMERAL_PUBLIC_KEY_LENGTH;
      imsi->empheral_public_key[EPHEMERAL_PUBLIC_KEY_LENGTH] = '\0';
      cipherTextLen = ielen - 48;
      OAILOG_DEBUG(LOG_NAS5G, "PROFILE-A ciphertext length: %d", cipherTextLen);

    } else {
      memcpy(&imsi->empheral_public_key, buffer + decoded,
             EPHEMERAL_PUBLIC_KEY_LENGTH + PROFILE_B_LEN);
      decoded += (EPHEMERAL_PUBLIC_KEY_LENGTH + PROFILE_B_LEN);
      imsi->empheral_public_key[EPHEMERAL_PUBLIC_KEY_LENGTH + PROFILE_B_LEN] =
          '\0';
      cipherTextLen = ielen - 48 - PROFILE_B_LEN;
      OAILOG_DEBUG(LOG_NAS5G, "PROFILE-B ciphertext length: %d", cipherTextLen);
    }

    imsi->ciphertext = blk2bstr(buffer + decoded, cipherTextLen);
    decoded += cipherTextLen;

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

  return ielen;
};

// Will be supported POST MVC
// Decode SuciMobileIdentity IE
int M5GSMobileIdentityMsg::DecodeSuciMobileIdentityMsg(
    SuciM5GSMobileIdentity* suci, uint8_t* buffer, uint8_t ielen) {
  int decoded = 0;

  suci->spare2 = (*(buffer + decoded) >> 7) & 0x1;
  suci->supi_format = (*(buffer + decoded) >> 4) & 0x7;
  suci->spare1 = (*(buffer + decoded) >> 3) & 0x1;
  suci->type_of_identity = *(buffer + decoded) & 0x7;

  if (suci->type_of_identity != M5GSMobileIdentityMsg_IMEISV) {
    OAILOG_ERROR(LOG_NAS5G, "TLV_VALUE_DOESNT_MATCH error");
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

  tmsi->spare = (*(buffer + decoded) >> 4) & 0xf;

  if (tmsi->spare != 0xf) {
    OAILOG_ERROR(LOG_NAS5G, "Error : %d",
                 static_cast<int>(TLV_VALUE_DOESNT_MATCH));
    return (TLV_VALUE_DOESNT_MATCH);
  }

  tmsi->odd_even = (*(buffer + decoded) >> 3) & 0x1;
  tmsi->type_of_identity = *(buffer + decoded) & 0x7;

  if (tmsi->type_of_identity != M5GSMobileIdentityMsg_TMSI) {
    OAILOG_ERROR(LOG_NAS5G, "Error : %d",
                 static_cast<int>(TLV_VALUE_DOESNT_MATCH));
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

  return ielen;
};

// Decode M5GSMobileIdentity IE
int M5GSMobileIdentityMsg::DecodeM5GSMobileIdentityMsg(
    M5GSMobileIdentityMsg* mg5smobile_identity, uint8_t iei, uint8_t* buffer,
    uint32_t len) {
  int decoded_rc = TLV_VALUE_DOESNT_MATCH;
  int decoded = 0;
  uint16_t ielen = 0;

  if (iei > 0) {
    CHECK_IEI_DECODER(iei, (unsigned char)*buffer);
    decoded++;
  }

  IES_DECODE_U16(buffer, decoded, ielen);
  CHECK_LENGTH_DECODER(len - decoded, ielen);
  unsigned char type_of_identity = *(buffer + decoded) & 0x7;

  if (type_of_identity == M5GSMobileIdentityMsg_IMEISV) {
    decoded_rc = DecodeSuciMobileIdentityMsg(
        &mg5smobile_identity->mobile_identity.suci, buffer, ielen);
  } else if (type_of_identity == M5GSMobileIdentityMsg_GUTI) {
    decoded_rc = DecodeGutiMobileIdentityMsg(
        &mg5smobile_identity->mobile_identity.guti, buffer + decoded, ielen);
  } else if (type_of_identity == M5GSMobileIdentityMsg_IMEI) {
    decoded_rc = DecodeImeiMobileIdentityMsg(
        &mg5smobile_identity->mobile_identity.imei, buffer + decoded, ielen);
  } else if (type_of_identity == M5GSMobileIdentityMsg_TMSI) {
    decoded_rc = DecodeTmsiMobileIdentityMsg(
        &mg5smobile_identity->mobile_identity.tmsi, buffer + decoded, ielen);
  } else if (type_of_identity == M5GSMobileIdentityMsg_SUCI_IMSI) {
    decoded_rc = DecodeImsiMobileIdentityMsg(
        &mg5smobile_identity->mobile_identity.imsi, buffer + decoded, ielen);
  } else if (type_of_identity == M5GSMobileIdentityMsg_NO_IDENTITY) {
    decoded_rc = 1;
  }
  if (decoded_rc < 0) {
    OAILOG_ERROR(LOG_NAS5G, "Decode Error");
    return decoded_rc;
  }
  return (decoded + decoded_rc);
};

// Encode GutiMobileIdentity IE
int M5GSMobileIdentityMsg::EncodeGutiMobileIdentityMsg(
    GutiM5GSMobileIdentity* guti, uint8_t* buffer) {
  uint32_t encoded = 0;

  *(buffer + encoded) =
      0xf0 | ((guti->odd_even & 0x01) << 3) | (guti->type_of_identity & 0x7);
  encoded++;
  *(buffer + encoded) =
      0x00 | ((guti->mcc_digit2 & 0x0f) << 4) | (guti->mcc_digit1 & 0x0f);
  encoded++;
  *(buffer + encoded) =
      0x00 | ((guti->mnc_digit3 & 0x0f) << 4) | (guti->mcc_digit3 & 0x0f);
  encoded++;
  *(buffer + encoded) =
      0x00 | ((guti->mnc_digit2 & 0x0f) << 4) | (guti->mnc_digit1 & 0x0f);
  encoded++;
  *(buffer + encoded) = 0x00 | guti->amf_regionid;
  encoded++;
  *(buffer + encoded) = 0x00 | ((guti->amf_setid >> 2) & 0xFF);
  encoded++;
  *(buffer + encoded) =
      0x00 | ((guti->amf_setid & 0xF3) << 6) | (guti->amf_pointer & 0x3f);
  encoded++;
  *(buffer + encoded) = 0x00 | guti->tmsi1;
  encoded++;
  *(buffer + encoded) = 0x00 | guti->tmsi2;
  encoded++;
  *(buffer + encoded) = 0x00 | guti->tmsi3;
  encoded++;
  *(buffer + encoded) = 0x00 | guti->tmsi4;
  encoded++;

  return encoded;
};

// Will be supported POST MVC
// Encode ImeiMobileIdentity IE
int M5GSMobileIdentityMsg::EncodeImeiMobileIdentityMsg(
    ImeiM5GSMobileIdentity* imei, uint8_t* buffer) {
  uint32_t encoded = 0;

  *(buffer + encoded) = 0x00 | ((imei->identity_digit1 & 0xf0) << 4) |
                        ((imei->odd_even & 0x1) << 3) |
                        (imei->type_of_identity & 0x7);
  encoded++;
  *(buffer + encoded) = 0x00 | ((imei->identity_digit2 & 0xf0) << 4) |
                        (imei->identity_digit3 & 0x0f);
  encoded++;

  return encoded;
};

// Will be supported POST MVC
// Encode ImsiMobileIdentity IE
int M5GSMobileIdentityMsg::EncodeImsiMobileIdentityMsg(
    ImsiM5GSMobileIdentity* imsi, uint8_t* buffer) {
  uint32_t encoded = 0;
  *(buffer + encoded) =
      0x00 | ((imsi->spare2 & 0x80) << 7) | ((imsi->supi_format & 0x07) << 4) |
      ((imsi->spare1 & 0x01) << 3) | (imsi->type_of_identity & 0x7);
  encoded++;
  *(buffer + encoded) =
      0x00 | ((imsi->mcc_digit2 & 0x0f) << 4) | (imsi->mcc_digit1 & 0x0f);
  encoded++;
  *(buffer + encoded) =
      0x00 | ((imsi->mnc_digit3 & 0x0f) << 4) | (imsi->mcc_digit3 & 0x0f);
  encoded++;
  *(buffer + encoded) =
      0x00 | ((imsi->mnc_digit2 & 0x0f) << 4) | (imsi->mnc_digit1 & 0x0f);
  encoded++;
  *(buffer + encoded) =
      0x00 | ((imsi->rout_ind_digit_2) << 4) | (imsi->rout_ind_digit_1);
  encoded++;
  *(buffer + encoded) =
      0x00 | ((imsi->rout_ind_digit_3) << 4) | (imsi->rout_ind_digit_4);
  encoded++;
  *(buffer + encoded) =
      0x00 | ((imsi->spare6 & 0x01) << 7) | ((imsi->spare5 & 0x01) << 6) |
      ((imsi->spare4 & 0x01) << 5) | ((imsi->spare3 & 0x01) << 4) |
      (imsi->protect_schm_id & 0x0f);
  *(buffer + encoded) = imsi->home_nw_id;
  encoded++;

  memcpy(buffer + encoded, &imsi->scheme_output, imsi->scheme_len);
  encoded = encoded + imsi->scheme_len;

  return encoded;
};

// Will be supported POST MVC
int M5GSMobileIdentityMsg::EncodeTmsiMobileIdentityMsg(
    TmsiM5GSMobileIdentity* tmsi, uint8_t* buffer) {
  uint32_t encoded = 0;

  *(buffer + encoded) = 0x00 | ((tmsi->spare & 0x0f) << 4) |
                        ((tmsi->odd_even & 0x01) << 3) |
                        (tmsi->type_of_identity & 0x7);
  encoded++;
  *(buffer + encoded) = 0x00 | tmsi->amf_setid;
  encoded++;
  *(buffer + encoded) = 0x00 | ((tmsi->amf_setid & 0xc0) << 6);
  *(buffer + encoded) = 0x00 | (tmsi->amf_pointer & 0x3f);
  encoded++;
  return encoded;
};

// Encode SuciMobileIdentity IE
int M5GSMobileIdentityMsg::EncodeSuciMobileIdentityMsg(
    SuciM5GSMobileIdentity* suci, uint8_t* buffer) {
  uint32_t encoded = 0;

  *(buffer + encoded) =
      0x00 | ((suci->spare2 & 0x80) << 7) | ((suci->supi_format & 0x07) << 4) |
      ((suci->spare1 & 0x01) << 3) | (suci->type_of_identity & 0x7);
  encoded++;
  suci->suci_nai.assign((const char*)(buffer + encoded), suci->suci_nai.size());

  return encoded;
};

// Encode M5GSMobileIdentity IE
int M5GSMobileIdentityMsg::EncodeM5GSMobileIdentityMsg(
    M5GSMobileIdentityMsg* m5gs_mobile_identity, uint8_t iei, uint8_t* buffer,
    uint32_t len) {
  uint16_t* lenPtr;
  int encoded_rc = TLV_VALUE_DOESNT_MATCH;
  uint32_t encoded = 0;

  // Checking IEI and pointer
  CHECK_PDU_POINTER_AND_LENGTH_ENCODER(buffer, MOBILE_IDENTITY_MIN_LENGTH, len);

  if (m5gs_mobile_identity->iei > 0) {
    CHECK_IEI_ENCODER((unsigned char)iei, m5gs_mobile_identity->iei);
    *buffer = iei;
    encoded++;
  } else
    return 0;

  lenPtr = (uint16_t*)(buffer + encoded);
  encoded += 2;
  m5gs_mobile_identity->toi =
      m5gs_mobile_identity->mobile_identity.guti.type_of_identity;
  if (m5gs_mobile_identity->toi == M5GSMobileIdentityMsg_SUCI_IMSI) {
    encoded_rc = EncodeImsiMobileIdentityMsg(
        &m5gs_mobile_identity->mobile_identity.imsi, buffer + encoded);
  } else if (m5gs_mobile_identity->toi == M5GSMobileIdentityMsg_IMEI) {
    encoded_rc = EncodeImeiMobileIdentityMsg(
        &m5gs_mobile_identity->mobile_identity.imei, buffer + encoded);
  } else if (m5gs_mobile_identity->toi == M5GSMobileIdentityMsg_GUTI) {
    encoded_rc = EncodeGutiMobileIdentityMsg(
        &m5gs_mobile_identity->mobile_identity.guti, buffer + encoded);
  } else if (m5gs_mobile_identity->toi == M5GSMobileIdentityMsg_TMSI) {
    encoded_rc = EncodeTmsiMobileIdentityMsg(
        &m5gs_mobile_identity->mobile_identity.tmsi, buffer + encoded);
  } else if (m5gs_mobile_identity->toi == M5GSMobileIdentityMsg_IMEISV) {
    encoded_rc = EncodeSuciMobileIdentityMsg(
        &m5gs_mobile_identity->mobile_identity.suci, buffer + encoded);
  }

  if (encoded_rc < 0) {
    OAILOG_ERROR(LOG_NAS5G, "Encode error");
    return encoded_rc;
  }

  *lenPtr = htons(encoded + encoded_rc - 2 - ((iei > 0) ? 1 : 0));
  return (encoded + encoded_rc);
};
}  // namespace magma5g

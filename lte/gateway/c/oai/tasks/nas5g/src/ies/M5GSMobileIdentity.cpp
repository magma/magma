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
#include "5GSMobileIdentity.h"
#include "CommonDefs.h"

using namespace std;
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

  guti->oddeven        = (*(buffer + decoded) >> 3) & 0x1;
  guti->typeofidentity = *(buffer + decoded) & 0x7;

  if (guti->typeofidentity != M5GSMobileIdentityMsg_GUTI) {
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
  guti->amfregionid = *(buffer + decoded);
  decoded++;
  setid = *(buffer + decoded);
  decoded++;
  guti->amfsetid =
      0x0000 | ((setid & 0xff) << 2) | ((*(buffer + decoded) >> 6) & 0x3);
  guti->amfpointer = *(buffer + decoded) & 0x3f;
  decoded++;
  guti->tmsi1 = *(buffer + decoded);
  decoded++;
  guti->tmsi2 = *(buffer + decoded);
  decoded++;
  guti->tmsi3 = *(buffer + decoded);
  decoded++;
  guti->tmsi4 = *(buffer + decoded);
  decoded++;
  MLOG(MDEBUG) << "   Odd/Even Indecation = " << dec << int(guti->oddeven)
               << "\n";
  MLOG(MDEBUG) << "   Mobile Country Code (MCC) = " << dec
               << int(guti->mcc_digit1) << dec << int(guti->mcc_digit2) << dec
               << int(guti->mcc_digit3) << "\n";
  MLOG(MDEBUG) << "   Mobile Network Code (MNC) = " << dec
               << int(guti->mnc_digit1) << dec << int(guti->mnc_digit2) << dec
               << int(guti->mnc_digit3) << "\n";
  MLOG(MDEBUG) << "   Amf Region ID = " << dec << int(guti->amfregionid)
               << "\n";
  MLOG(MDEBUG) << "   Amf Set ID = " << dec << int(guti->amfsetid) << "\n";
  MLOG(MDEBUG) << "   Amf Pointer = " << dec << int(guti->amfpointer) << "\n";
  MLOG(MDEUBG) << "   M5G-TMSI = "
               << "0x0" << hex << int(guti->tmsi1) << "0" << hex
               << int(guti->tmsi2) << "0" << hex << int(guti->tmsi3) << "0"
               << hex << int(guti->tmsi4) << "\n\n";
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

  imei->oddeven        = (*(buffer + decoded) >> 3) & 0x1;
  imei->typeofidentity = *(buffer + decoded) & 0x7;

  if (imei->typeofidentity != M5GSMobileIdentityMsg_IMEI) {
    MLOG(MERROR) << "Error: " << std::dec << TLV_VALUE_DOESNT_MATCH;
    return (TLV_VALUE_DOESNT_MATCH);
  }

  decoded++;
  imei->identity_digit3 = (*(buffer + decoded) >> 4) & 0xf;
  imei->identity_digit2 = *(buffer + decoded) & 0xf;
  decoded++;
  MLOG(MDEBUG) << "  oddeven = " << hex << int(imei->oddeven) << "\n";
  MLOG(MDEBUG) << "  digit1 = " << hex << int(imei->identity_digit1) << "\n";
  MLOG(MDEBUG) << "  digit2 = " << hex << int(imei->identity_digit2) << "\n";
  MLOG(MDEBUG) << "  digit3 = " << hex << int(imei->identity_digit3) << "\n";

  return (decoded);
};

// Decode ImsiMobileIdentity IE
int M5GSMobileIdentityMsg::DecodeImsiMobileIdentityMsg(
    ImsiM5GSMobileIdentity* imsi, uint8_t* buffer, uint8_t ielen) {
  int decoded = 0;

  MLOG(MDEBUG) << " DecodeImsiMobileIdentityMsg:";
  imsi->spare2         = (*(buffer + decoded) >> 7) & 0x1;
  imsi->supiformat     = (*(buffer + decoded) >> 4) & 0x7;
  imsi->spare1         = (*(buffer + decoded) >> 3) & 0x1;
  imsi->typeofidentity = *(buffer + decoded) & 0x7;

  if (imsi->typeofidentity != M5GSMobileIdentityMsg_IMSI) {
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
  imsi->routingindicatordigit2 = (*(buffer + decoded) >> 4) & 0xf;
  imsi->routingindicatordigit1 = *(buffer + decoded) & 0xf;
  decoded++;
  imsi->routingindicatordigit4 = (*(buffer + decoded) >> 4) & 0xf;
  imsi->routingindicatordigit3 = *(buffer + decoded) & 0xf;
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
  // TODO
  /* Scheme output (octets 12 to x)
     The Scheme output field consists of a string of characters with a variable
     length or hexadecimal digits as specified in 3GPP TS 23.003 [4]. If
     Protection scheme identifier is set to "0000" (i.e. Null scheme), then the
     Scheme output consists of the MSIN and is coded using BCD coding with each
     digit of the MSIN coded over 4 bits. If the MSIN includes an odd number of
     digits, bits 5 to 8 of octet x shall be coded as "1111". If Protection
     scheme identifier is not "0000" (i.e. ECIES scheme profile A, ECIES scheme
     profile B or Operator-specific protection scheme), then Scheme output is
     coded as hexadecimal digits
  */

  int tmp = ielen - decoded;
  decoded = ielen;

  MLOG(MDEBUG) << "  Spare = " << hex << int(imsi->spare2);
  MLOG(MDEBUG) << "  Supi Format = " << hex << int(imsi->supiformat);
  MLOG(MDEBUG) << "  Spare = " << hex << int(imsi->spare1);
  MLOG(MDEBUG) << "  Type of Identity = " << hex << int(imsi->typeofidentity);
  MLOG(MDEBUG) << "  Mobile Country Code (MCC) = " << dec
               << int(imsi->mcc_digit1) << dec << int(imsi->mcc_digit2) << dec
               << int(imsi->mcc_digit3);
  MLOG(MDEBUG) << "  Mobile Network Code (MNC) = " << dec
               << int(imsi->mnc_digit1) << dec << int(imsi->mnc_digit2) << dec
               << int(imsi->mnc_digit3);
  MLOG(MDEBUG) << "  Routing Indicator = " << hex
               << int(imsi->routingindicatordigit1);
  MLOG(MDEBUG) << "  Protection Scheme ID = " << hex
               << int(imsi->protect_schm_id);
  MLOG(MDEBUG) << "  Home Network Public Key Identifier = " << hex
               << int(imsi->home_nw_id);
  MLOG(MDEBUG) << "  Scheme Output = ";
  BUFFER_PRINT_LOG(imsi->scheme_output, tmp)

  return (decoded);
};

// Will be supported POST MVC
// Decode SuciMobileIdentity IE
int M5GSMobileIdentityMsg::DecodeSuciMobileIdentityMsg(
    SuciM5GSMobileIdentity* suci, uint8_t* buffer, uint8_t ielen) {
  int decoded = 0;

  MLOG(MDEBUG) << "         DecodeSuciMobileIdentityMsg:"
               << "\n";
  suci->spare2         = (*(buffer + decoded) >> 7) & 0x1;
  suci->supiformat     = (*(buffer + decoded) >> 4) & 0x7;
  suci->spare1         = (*(buffer + decoded) >> 3) & 0x1;
  suci->typeofidentity = *(buffer + decoded) & 0x7;

  if (suci->typeofidentity != M5GSMobileIdentityMsg_SUCI) {
    MLOG(MDEBUG) << "TLV_VALUE_DOESNT_MATCH error";
    return (TLV_VALUE_DOESNT_MATCH);
  }
  decoded++;

  // Will be supported POST MVC
  suci->sucinai = *(buffer + decoded);
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

  tmsi->oddeven        = (*(buffer + decoded) >> 3) & 0x1;
  tmsi->typeofidentity = *(buffer + decoded) & 0x7;

  if (tmsi->typeofidentity != M5GSMobileIdentityMsg_TMSI) {
    MLOG(MDEBUG) << "Error: " << int(TLV_VALUE_DOESNT_MATCH);
    return (TLV_VALUE_DOESNT_MATCH);
  }
  decoded++;
  tmsi->amfsetid = *(buffer + decoded);
  decoded++;
  tmsi->amfsetid1  = (*(buffer + decoded) >> 6) & 0x2;
  tmsi->amfpointer = *(buffer + decoded) & 0x3f;
  decoded++;
  tmsi->m5gtmsi1 = *(buffer + decoded);
  decoded++;
  tmsi->m5gtmsi2 = *(buffer + decoded);
  decoded++;
  tmsi->m5gtmsi3 = *(buffer + decoded);
  decoded++;
  tmsi->m5gtmsi4 = *(buffer + decoded);
  decoded++;

  MLOG(MDEBUG) << "  spare2 = " << hex << int(tmsi->spare);
  MLOG(MDEBUG) << "  oddeven = " << hex << int(tmsi->oddeven);
  MLOG(MDEBUG) << "  typeofidentity = " << hex << int(tmsi->typeofidentity);
  MLOG(MDEBUG) << "  amfsetid = " << hex << int(tmsi->amfsetid);
  MLOG(MDEBUG) << "  amfsetid1 = " << hex << int(tmsi->amfsetid1);
  MLOG(MDEBUG) << "  amfpointer = " << hex << int(tmsi->amfpointer);
  MLOG(MDEBUG) << "  m5gtmsi1 = " << hex << int(tmsi->m5gtmsi1);
  MLOG(MDEBUG) << "  m5gtmsi2 = " << hex << int(tmsi->m5gtmsi2);
  MLOG(MDEBUG) << "  m5gtmsi3 = " << hex << int(tmsi->m5gtmsi3);
  MLOG(MDEBUG) << "  m5gtmsi4 = " << hex << int(tmsi->m5gtmsi4);

  return (decoded);
};

// Decode M5GSMobileIdentity IE
int M5GSMobileIdentityMsg::DecodeM5GSMobileIdentityMsg(
    M5GSMobileIdentityMsg* mg5smobileidentity, uint8_t iei, uint8_t* buffer,
    uint32_t len) {
  int decoded_rc = TLV_VALUE_DOESNT_MATCH;
  int decoded    = 0;
  uint8_t ielen  = 0;

  MLOG(MDEBUG) << "M5GS Mobile Identity : ";
  if (iei > 0) {
    CHECK_IEI_DECODER(iei, (unsigned char) *buffer);
    decoded++;
  }

  IES_DECODE_U16(buffer, decoded, ielen);
  CHECK_LENGTH_DECODER(len - decoded, ielen);
  unsigned char typeofidentity = *(buffer + decoded) & 0x7;
  MLOG(MDEBUG) << " Length = " << dec << int(ielen)
               << " Type of Identity = " << dec << int(typeofidentity);

  if (typeofidentity == M5GSMobileIdentityMsg_SUCI) {
    MLOG(MDEBUG) << " Type suci";
    decoded_rc = DecodeSuciMobileIdentityMsg(
        &mg5smobileidentity->mobileidentity.suci, buffer, ielen);
  } else if (typeofidentity == M5GSMobileIdentityMsg_GUTI) {
    MLOG(MDEBUG) << " Type guti";
    decoded_rc = DecodeGutiMobileIdentityMsg(
        &mg5smobileidentity->mobileidentity.guti, buffer + decoded, ielen);
  } else if (typeofidentity == M5GSMobileIdentityMsg_IMEI) {
    MLOG(MDEBUG) << " Type imei";
    decoded_rc = DecodeImeiMobileIdentityMsg(
        &mg5smobileidentity->mobileidentity.imei, buffer + decoded, ielen);
  } else if (typeofidentity == M5GSMobileIdentityMsg_TMSI) {
    MLOG(MDEBUG) << " Type tmsi";
    decoded_rc = DecodeTmsiMobileIdentityMsg(
        &mg5smobileidentity->mobileidentity.tmsi, buffer + decoded, ielen);
  } else if (typeofidentity == M5GSMobileIdentityMsg_IMSI) {
    MLOG(MDEBUG) << " Type imsi";
    decoded_rc = DecodeImsiMobileIdentityMsg(
        &mg5smobileidentity->mobileidentity.imsi, buffer + decoded, ielen);
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
      0xf0 | ((guti->oddeven & 0x01) << 3) | (guti->typeofidentity & 0x7);
  MLOG(MDEBUG) << "oddeven typeofidentity = " << hex
               << int(*(buffer + encoded));
  encoded++;
  *(buffer + encoded) =
      0x00 | ((guti->mcc_digit2 & 0x0f) << 4) | (guti->mcc_digit1 & 0x0f);
  MLOG(MDEBUG) << "mcc_digit2 >mcc_digit1 typeofidentity = " << hex
               << int(*(buffer + encoded));
  encoded++;
  *(buffer + encoded) =
      0x00 | ((guti->mnc_digit3 & 0x0f) << 4) | (guti->mcc_digit3 & 0x0f);
  MLOG(MDEBUG) << "mnc_digit3 >mcc_digit3 typeofidentity = " << hex
               << int(*(buffer + encoded));
  encoded++;
  *(buffer + encoded) =
      0x00 | ((guti->mnc_digit2 & 0x0f) << 4) | (guti->mnc_digit1 & 0x0f);
  MLOG(MDEBUG) << "mnc_digit2 >mcc_digit1 typeofidentity = " << hex
               << int(*(buffer + encoded));
  encoded++;
  *(buffer + encoded) = 0x00 | guti->amfregionid;
  MLOG(MDEBUG) << "amfregionid = " << hex << int(*(buffer + encoded));
  encoded++;
  *(buffer + encoded) = 0x00 | guti->amfsetid;
  MLOG(MDEBUG) << "amfsetid = " << hex << int(*(buffer + encoded));
  encoded++;
  *(buffer + encoded) =
      0x00 | ((guti->amfsetid1 & 0x03) << 6) | (guti->amfpointer & 0x3f);
  MLOG(MDEBUG) << "amfsetid1 amfpointer = " << hex << int(*(buffer + encoded));
  encoded++;
  *(buffer + encoded) = 0x00 | guti->tmsi1;
  MLOG(MDEBUG) << "tmsi1 = " << hex << int(*(buffer + encoded));
  encoded++;
  *(buffer + encoded) = 0x00 | guti->tmsi2;
  MLOG(MDEBUG) << "tmsi2 = " << hex << int(*(buffer + encoded));
  encoded++;
  *(buffer + encoded) = 0x00 | guti->tmsi3;
  MLOG(MDEBUG) << "tmsi3 = " << hex << int(*(buffer + encoded));
  encoded++;
  *(buffer + encoded) = 0x00 | guti->tmsi4;
  MLOG(MDEBUG) << "tmsi4 = " << hex << int(*(buffer + encoded));
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
                        ((imei->oddeven & 0x1) << 3) |
                        (imei->typeofidentity & 0x7);
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
  MLOG(MDEBUG) << "EncodeImsiMobileIdentityMsg:";
  *(buffer + encoded) =
      0x00 | ((imsi->spare2 & 0x80) << 7) | ((imsi->supiformat & 0x07) << 4) |
      ((imsi->spare1 & 0x01) << 3) | (imsi->typeofidentity & 0x7);
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
  *(buffer + encoded) = 0x00 | ((imsi->routingindicatordigit2 & 0xf0) << 4) |
                        (imsi->routingindicatordigit1 & 0x0f);
  encoded++;
  *(buffer + encoded) = 0x00 | ((imsi->routingindicatordigit3 & 0xf0) << 4) |
                        (imsi->routingindicatordigit4 & 0x0f);
  encoded++;
  *(buffer + encoded) =
      0x00 | ((imsi->spare6 & 0x01) << 7) | ((imsi->spare5 & 0x01) << 6) |
      ((imsi->spare4 & 0x01) << 5) | ((imsi->spare3 & 0x01) << 4) |
      (imsi->protect_schm_id & 0x0f);
  *(buffer + encoded) = imsi->home_nw_id;
  encoded++;
  /*
  // Will be supported POST MVC
  imsi->scheme_output.assign((const char *)(buffer + encoded),
  imsi->scheme_output.size()); MLOG(MDEBUG) << "ielen = " << hex << (unsigned
  char)imsi->scheme_output.size() ; MLOG(MDEBUG) << "contents"; int i = 0;
  for(i; i < imsi->scheme_output.size(); i++) {
    MLOG(MDEBUG) << (uint8_t)(imsi->scheme_output[i]);
  }
  MLOG(MDEBUG) << endl;
  */
  return encoded;
};

// Will be supported POST MVC
int M5GSMobileIdentityMsg::EncodeTmsiMobileIdentityMsg(
    TmsiM5GSMobileIdentity* tmsi, uint8_t* buffer) {
  uint32_t encoded = 0;

  MLOG(MDEBUG) << "EncodeTmsiMobileIdentityMsg:";
  *(buffer + encoded) = 0x00 | ((tmsi->spare & 0x0f) << 4) |
                        ((tmsi->oddeven & 0x01) << 3) |
                        (tmsi->typeofidentity & 0x7);
  encoded++;
  *(buffer + encoded) = 0x00 | tmsi->amfsetid;
  encoded++;
  *(buffer + encoded) = 0x00 | ((tmsi->amfsetid1 & 0xc0) << 6);
  *(buffer + encoded) = 0x00 | (tmsi->amfpointer & 0x3f);
  encoded++;
  *(buffer + encoded) = 0x00 | tmsi->m5gtmsi1;
  encoded++;
  *(buffer + encoded) = 0x00 | tmsi->m5gtmsi2;
  encoded++;
  *(buffer + encoded) = 0x00 | tmsi->m5gtmsi3;
  encoded++;
  *(buffer + encoded) = 0x00 | tmsi->m5gtmsi4;
  encoded++;

  return encoded;
};

// Encode SuciMobileIdentity IE
int M5GSMobileIdentityMsg::EncodeSuciMobileIdentityMsg(
    SuciM5GSMobileIdentity* suci, uint8_t* buffer) {
  uint32_t encoded = 0;

  MLOG(MDEBUG) << "EncodeSuciMobileIdentityMsg:";
  *(buffer + encoded) =
      0x00 | ((suci->spare2 & 0x80) << 7) | ((suci->supiformat & 0x07) << 4) |
      ((suci->spare1 & 0x01) << 3) | (suci->typeofidentity & 0x7);
  encoded++;
  suci->sucinai.assign((const char*) (buffer + encoded), suci->sucinai.size());
  MLOG(MDEBUG) << "ielen = " << hex << (unsigned char) suci->sucinai.size();
  MLOG(MDEBUG) << "contents";
  int i = 0;
  for (i; i < suci->sucinai.size(); i++) {
    MLOG(MDEBUG) << hex << int(suci->sucinai[i]);
  }
  MLOG(MDEBUG) << endl;

  return encoded;
};

// Encode M5GSMobileIdentity IE
int M5GSMobileIdentityMsg::EncodeM5GSMobileIdentityMsg(
    M5GSMobileIdentityMsg* m5gsmobileidentity, uint8_t iei, uint8_t* buffer,
    uint32_t len) {
  uint16_t* lenPtr;
  int encoded_rc   = TLV_VALUE_DOESNT_MATCH;
  uint32_t encoded = 0;

  MLOG(MDEBUG) << "EncodeM5GSMobileIdentityMsg:";

  // Checking IEI and pointer
  CHECK_PDU_POINTER_AND_LENGTH_ENCODER(buffer, MOBILE_IDENTITY_MIN_LENGTH, len);

  if (iei > 0) {
    CHECK_IEI_ENCODER((unsigned char) iei, m5gsmobileidentity->iei);
    *buffer = iei;
    MLOG(MDEBUG) << "iei" << hex << int(*buffer);
    encoded++;
  }

  lenPtr = (uint16_t*) (buffer + encoded);
  encoded += 2;

  if (m5gsmobileidentity->mobileidentity.imsi.typeofidentity ==
      M5GSMobileIdentityMsg_IMSI) {
    MLOG(MDEBUG) << "Type imsi";
    encoded_rc = EncodeImsiMobileIdentityMsg(
        &m5gsmobileidentity->mobileidentity.imsi, buffer + encoded);
  } else if (
      m5gsmobileidentity->mobileidentity.imei.typeofidentity ==
      M5GSMobileIdentityMsg_IMEI) {
    MLOG(MDEBUG) << "Type imei";
    encoded_rc = EncodeImeiMobileIdentityMsg(
        &m5gsmobileidentity->mobileidentity.imei, buffer + encoded);
  } else if (
      m5gsmobileidentity->mobileidentity.guti.typeofidentity ==
      M5GSMobileIdentityMsg_GUTI) {
    MLOG(MDEBUG) << "Type guti";
    encoded_rc = EncodeGutiMobileIdentityMsg(
        &m5gsmobileidentity->mobileidentity.guti, buffer + encoded);
  } else if (
      m5gsmobileidentity->mobileidentity.tmsi.typeofidentity ==
      M5GSMobileIdentityMsg_TMSI) {
    MLOG(MDEBUG) << "Type tmsi";
    encoded_rc = EncodeTmsiMobileIdentityMsg(
        &m5gsmobileidentity->mobileidentity.tmsi, buffer + encoded);
  } else if (
      m5gsmobileidentity->mobileidentity.suci.typeofidentity ==
      M5GSMobileIdentityMsg_SUCI) {
    MLOG(MDEBUG) << "Type suci";
    encoded_rc = EncodeSuciMobileIdentityMsg(
        &m5gsmobileidentity->mobileidentity.suci, buffer + encoded);
  }

  if (encoded_rc < 0) {
    MLOG(MDEBUG) << "Encode error" << encoded_rc;
    return encoded_rc;
  }

  *lenPtr = htons(encoded + encoded_rc - 2 - ((iei > 0) ? 1 : 0));
  return (encoded + encoded_rc);
};
}  // namespace magma5g
